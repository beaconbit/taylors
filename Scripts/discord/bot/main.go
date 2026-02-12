package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Task represents a task from the JSON files
type Task struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Status         string   `json:"status"`
	Priority       string   `json:"priority"`
	Created        string   `json:"created"`
	Description    string   `json:"description"`
	Assignee       string   `json:"assignee"`
	InvolvedPeople []string `json:"involved_people"`
}

var (
	DiscordToken string
	TasksDir     = "/home/ubuntu/Taylors"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Load configuration from environment variables
	DiscordToken = os.Getenv("DISCORD_TOKEN")
	if DiscordToken == "" {
		log.Fatal("ERROR: DISCORD_TOKEN not set. Please create a .env file with DISCORD_TOKEN=your_bot_token")
	}
}

func main() {
	// Create Discord session
	dg, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Register message handler
	dg.AddHandler(messageCreate)

	// We need message content intent
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	// Open connection
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer dg.Close()

	log.Println("Discord bot is now running. Press CTRL-C to exit.")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Shutting down bot...")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check for command prefix
	if !strings.HasPrefix(m.Content, "!task") {
		return
	}

	// Parse command
	parts := strings.Fields(m.Content)
	if len(parts) < 2 {
		sendHelp(s, m.ChannelID)
		return
	}

	command := parts[1]
	args := parts[2:]

	switch command {
	case "list":
		handleListCommand(s, m.ChannelID, args)
	case "search":
		handleSearchCommand(s, m.ChannelID, args)
	case "status":
		handleStatusCommand(s, m.ChannelID, args)
	case "help":
		sendHelp(s, m.ChannelID)
	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unknown command: %s. Use `!task help` for available commands.", command))
	}
}

func sendHelp(s *discordgo.Session, channelID string) {
	helpText := `**Task Management Bot Commands:**

!task list [status] [priority] - List tasks (optional filters)
  Examples:
  !task list
  !task list pending
  !task list pending high
  !task list completed

!task search <query> - Search tasks by title or description
  Example: !task search "conveyor belt"

!task status <task-id> - Get detailed status of a task
  Example: !task status discord-123456

!task help - Show this help message

**Note:** This bot reads tasks from the JSON files in the task directories.`

	s.ChannelMessageSend(channelID, helpText)
}

func handleListCommand(s *discordgo.Session, channelID string, args []string) {
	// Get all tasks
	tasks, err := getAllTasks()
	if err != nil {
		s.ChannelMessageSend(channelID, fmt.Sprintf("Error reading tasks: %v", err))
		return
	}

	// Apply filters
	var filteredTasks []Task
	for _, task := range tasks {
		if len(args) >= 1 && task.Status != args[0] {
			continue
		}
		if len(args) >= 2 && task.Priority != args[1] {
			continue
		}
		filteredTasks = append(filteredTasks, task)
	}

	if len(filteredTasks) == 0 {
		s.ChannelMessageSend(channelID, "No tasks found matching the criteria.")
		return
	}

	// Create response
	var response strings.Builder
	response.WriteString(fmt.Sprintf("**Found %d tasks:**\n\n", len(filteredTasks)))

	for i, task := range filteredTasks {
		if i >= 10 { // Limit to 10 tasks per message
			response.WriteString(fmt.Sprintf("\n... and %d more tasks", len(filteredTasks)-10))
			break
		}
		response.WriteString(fmt.Sprintf("• **%s** - %s (%s priority)\n", task.ID, task.Title, task.Priority))
		response.WriteString(fmt.Sprintf("  Status: %s, Assignee: %s\n\n", task.Status, task.Assignee))
	}

	s.ChannelMessageSend(channelID, response.String())
}

func handleSearchCommand(s *discordgo.Session, channelID string, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(channelID, "Please provide a search query. Example: `!task search \"conveyor belt\"`")
		return
	}

	query := strings.ToLower(strings.Join(args, " "))
	tasks, err := getAllTasks()
	if err != nil {
		s.ChannelMessageSend(channelID, fmt.Sprintf("Error reading tasks: %v", err))
		return
	}

	var matchingTasks []Task
	for _, task := range tasks {
		title := strings.ToLower(task.Title)
		desc := strings.ToLower(task.Description)

		if strings.Contains(title, query) || strings.Contains(desc, query) {
			matchingTasks = append(matchingTasks, task)
		}
	}

	if len(matchingTasks) == 0 {
		s.ChannelMessageSend(channelID, fmt.Sprintf("No tasks found matching: %s", query))
		return
	}

	// Create response
	var response strings.Builder
	response.WriteString(fmt.Sprintf("**Found %d tasks matching \"%s\":**\n\n", len(matchingTasks), query))

	for i, task := range matchingTasks {
		if i >= 5 { // Limit to 5 tasks per message
			response.WriteString(fmt.Sprintf("\n... and %d more matching tasks", len(matchingTasks)-5))
			break
		}
		response.WriteString(fmt.Sprintf("• **%s** - %s\n", task.ID, task.Title))
		response.WriteString(fmt.Sprintf("  Status: %s, Priority: %s\n", task.Status, task.Priority))
		response.WriteString(fmt.Sprintf("  Created: %s\n\n", formatTime(task.Created)))
	}

	s.ChannelMessageSend(channelID, response.String())
}

func handleStatusCommand(s *discordgo.Session, channelID string, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(channelID, "Please provide a task ID. Example: `!task status discord-123456`")
		return
	}

	taskID := args[0]
	task, err := findTaskByID(taskID)
	if err != nil {
		s.ChannelMessageSend(channelID, fmt.Sprintf("Error finding task: %v", err))
		return
	}

	if task == nil {
		s.ChannelMessageSend(channelID, fmt.Sprintf("Task not found: %s", taskID))
		return
	}

	// Create detailed response
	response := fmt.Sprintf(`**Task Details: %s**

**Title:** %s
**Status:** %s
**Priority:** %s
**Assignee:** %s
**Created:** %s
**Involved People:** %s

**Description:**
%s`,
		task.ID, task.Title, task.Status, task.Priority, task.Assignee,
		formatTime(task.Created), strings.Join(task.InvolvedPeople, ", "),
		task.Description)

	s.ChannelMessageSend(channelID, response)
}

func getAllTasks() ([]Task, error) {
	var tasks []Task

	// This is a simplified version - in a real implementation,
	// you would walk through all task directories and read JSON files
	// For now, we'll return an empty list as a placeholder

	return tasks, nil
}

func findTaskByID(taskID string) (*Task, error) {
	// This is a simplified version - in a real implementation,
	// you would search through all task files
	// For now, return nil (not found)
	return nil, nil
}

func formatTime(timeStr string) string {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr
	}
	return t.Format("2006-01-02 15:04")
}
