package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"discord-scraper/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

// Configuration - loaded from .env file
var (
	DiscordToken string
	ChannelID    string
	DatabasePath = "discord_messages.db"
)

// Message represents a Discord message for storage
type Message struct {
	DiscordMessageID string
	ChannelID        string
	AuthorID         string
	AuthorName       string
	Content          string
	Timestamp        time.Time
	AttachmentsCount int
}

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Load configuration from environment variables
	DiscordToken = os.Getenv("DISCORD_TOKEN")
	ChannelID = os.Getenv("DISCORD_CHANNEL_ID")

	// Allow DatabasePath to be overridden by environment variable
	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		DatabasePath = dbPath
	}
}

func main() {
	// Check if token is set
	if DiscordToken == "" {
		log.Fatal("ERROR: DISCORD_TOKEN not set. Please create a .env file with DISCORD_TOKEN=your_bot_token")
	}

	// Check if channel ID is set
	if ChannelID == "" {
		log.Fatal("ERROR: DISCORD_CHANNEL_ID not set. Please create a .env file with DISCORD_CHANNEL_ID=your_channel_id")
	}

	log.Printf("Configuration loaded: Channel=%s, Database=%s", ChannelID, DatabasePath)

	// Setup database
	db := setupDatabase()
	defer db.Close()

	// Create Discord session
	dg, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}
	defer dg.Close()

	// We only need message content intent
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	// Open connection
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	// Calculate time threshold (1 hour ago)
	oneHourAgo := time.Now().Add(-5 * time.Minute)
	log.Printf("Fetching messages since: %v", oneHourAgo.Format(time.RFC3339))

	// Fetch messages
	messages, err := fetchMessages(dg, ChannelID, oneHourAgo)
	if err != nil {
		log.Fatalf("Error fetching messages: %v", err)
	}

	// Save messages to database
	savedCount := saveMessages(db, messages)

	// pass the messages to be processed by opencode
	for _, msg := range messages {
		authorName := msg.Author.Username
		content := msg.Content
		if msg.Member != nil && msg.Member.Nick != "" {
			authorName = msg.Member.Nick
		}
		err := utils.QueryOpenCode(content, authorName)
		if err != nil {
			log.Println(err)
		}
	}
	log.Printf("Successfully saved %d messages to database", savedCount)
}

func setupDatabase() *sql.DB {
	// Open or create database
	db, err := sql.Open("sqlite3", DatabasePath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Create messages table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		discord_message_id TEXT UNIQUE,
		channel_id TEXT,
		author_id TEXT,
		author_name TEXT,
		content TEXT,
		timestamp DATETIME,
		attachments_count INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_timestamp ON messages(timestamp);
	CREATE INDEX IF NOT EXISTS idx_discord_message_id ON messages(discord_message_id);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	log.Printf("Database setup complete at %s", DatabasePath)
	return db
}

func fetchMessages(s *discordgo.Session, channelID string, since time.Time) ([]*discordgo.Message, error) {
	var allMessages []*discordgo.Message
	var lastMessageID string

	for {
		// Fetch messages (100 at a time, which is the max per request)
		messages, err := s.ChannelMessages(channelID, 100, lastMessageID, "", "")
		if err != nil {
			return nil, fmt.Errorf("error fetching messages: %v", err)
		}

		if len(messages) == 0 {
			break
		}

		// Filter messages that are newer than our threshold
		for _, msg := range messages {
			// msg.Timestamp is already a time.Time, no need to parse
			msgTime := msg.Timestamp

			// If message is older than our threshold, stop fetching
			if msgTime.Before(since) {
				return allMessages, nil
			}

			allMessages = append(allMessages, msg)
		}

		// Set last message ID for next pagination
		lastMessageID = messages[len(messages)-1].ID

		// Rate limiting: be nice to Discord API
		time.Sleep(100 * time.Millisecond)
	}

	return allMessages, nil
}

func saveMessages(db *sql.DB, discordMessages []*discordgo.Message) int {
	savedCount := 0

	// Prepare insert statement
	stmt, err := db.Prepare(`
		INSERT OR IGNORE INTO messages 
		(discord_message_id, channel_id, author_id, author_name, content, timestamp, attachments_count)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Begin transaction for better performance
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error beginning transaction: %v", err)
	}

	txStmt := tx.Stmt(stmt)

	for _, msg := range discordMessages {
		// msg.Timestamp is already a time.Time, no need to parse
		timestamp := msg.Timestamp

		// Get author name
		authorName := msg.Author.Username
		if msg.Member != nil && msg.Member.Nick != "" {
			authorName = msg.Member.Nick
		}

		// Prepare message content (truncate if too long for logging)
		content := msg.Content
		logContent := content
		if len(logContent) > 50 {
			logContent = logContent[:50] + "..."
		}

		// Execute insert
		_, err = txStmt.Exec(
			msg.ID,
			msg.ChannelID,
			msg.Author.ID,
			authorName,
			content,
			timestamp.Format(time.RFC3339),
			len(msg.Attachments),
		)

		if err != nil {
			// Check if it's a duplicate error (we use INSERT OR IGNORE, so this shouldn't happen often)
			if strings.Contains(err.Error(), "UNIQUE constraint") {
				log.Printf("Message %s already exists in database", msg.ID)
			} else {
				log.Printf("Error saving message %s: %v", msg.ID, err)
			}
			continue
		}

		savedCount++
		log.Printf("Saved message from %s: %s", authorName, logContent)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	return savedCount
}

// Helper function to query recent messages (optional, for testing)
func queryRecentMessages(db *sql.DB, limit int) {
	rows, err := db.Query(`
		SELECT timestamp, author_name, content, attachments_count
		FROM messages 
		ORDER BY timestamp DESC 
		LIMIT ?
	`, limit)
	if err != nil {
		log.Printf("Error querying messages: %v", err)
		return
	}
	defer rows.Close()

	log.Printf("\n=== Last %d messages ===", limit)
	for rows.Next() {
		var timestampStr, authorName, content string
		var attachmentsCount int

		err := rows.Scan(&timestampStr, &authorName, &content, &attachmentsCount)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		timestamp, _ := time.Parse(time.RFC3339, timestampStr)
		contentPreview := content
		if len(contentPreview) > 50 {
			contentPreview = contentPreview[:50] + "..."
		}

		log.Printf("[%s] %s: %s (attachments: %d)",
			timestamp.Format("2006-01-02 15:04:05"),
			authorName,
			contentPreview,
			attachmentsCount)
	}
}
