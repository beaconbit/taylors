package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func Test() {
	log.Println("test works")
}

// ExampleQueryOpenCode calls the local OpenCode API with the given prompt
// It uses ~/Taylors/.opencode as the context path and includes instructions.md and config.json
func QueryOpenCode(prompt string, author string) error {
	// Get the context directory
	contextDir := "/home/ubuntu/Taylors/.opencode"

	// Check if the context directory exists
	if _, err := os.Stat(contextDir); os.IsNotExist(err) {
		return fmt.Errorf("context directory does not exist: %s", contextDir)
	}

	// Check for required files
	instructionsPath := filepath.Join(contextDir, "instructions.md")
	configPath := filepath.Join(contextDir, "config.json")

	if _, err := os.Stat(instructionsPath); os.IsNotExist(err) {
		log.Printf("Warning: instructions.md not found at %s", instructionsPath)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Warning: config.json not found at %s", configPath)
	}

	promptWithInstructions, err := WrapPromptWithInstructions(prompt, author)
	if err != nil {
		log.Fatalf("Warning: could not wrap prompt with instructions %s", prompt)
	}
	// Construct the full prompt with context
	fullPrompt := fmt.Sprintf(`Context from %s/.opencode directory:

Instructions: %s
Config: %s

User query: %s

Please provide a helpful response based on the task management context.`,
		"/home/ubuntu/Taylors",
		readFileIfExists(instructionsPath),
		readFileIfExists(configPath),
		promptWithInstructions)

	// create command
	cmd := exec.Command("/home/ubuntu/.opencode/bin/opencode", "run", fullPrompt)

	// Set working directory to the context directory so OpenCode can find the .opencode files
	cmd.Dir = "/home/ubuntu/Taylors"

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr


	// Run the command
	log.Printf("Calling OpenCode with prompt: %s...", truncateString(prompt, 100))
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("error running opencode: %v\nstderr: %s", err, stderr.String())
	}

	return nil
}

func WrapPromptWithInstructions(discordMessage string, author string) (string, error) {
	prompt := fmt.Sprintf(`Analyze this Discord message from %s 


and determine if it should become a new task or an update to an existing task,
if it should be a new task then create a new task as specified, if it is an update to 
an existing task then add an update object to the existing task. If it is an update to
an existing task that has a status "repeating" then update the "updated" and "next_occurance" fields.

if the message doesn't sound like it can be turned into a task, or isn't related to mechanical engineering in a laundry PLEASE IGNORE


Message: "%s"

Consider:
1. Does this describe work that needs to be done?
2. does it mention a priority (high/medium/low)? otherwise default to medium
4. What Category (All new tasks go in Maintenance, Updates to existing Tasks go in the existing task file which can be anywhere)
5. Suggest a concise title for the task.

`, author, discordMessage)

	return prompt, nil
}

// Helper function to read file if it exists
func readFileIfExists(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("[File not found: %s]", path)
	}
	return string(content)
}

// Helper function to truncate string for logging
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}
