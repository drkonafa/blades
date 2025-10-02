package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	"github.com/go-kratos/blades"
	"github.com/go-kratos/blades/contrib/zeus"
	"github.com/go-kratos/blades/flow"
)

// loadEnvFile loads environment variables from a .env file
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// loadConfig loads configuration from environment or .env file
func loadConfig() {
	// Try to load from .env file first
	if err := loadEnvFile(".env"); err != nil {
		// If .env file doesn't exist, that's okay - use system environment variables
	}

	// Also try to load from chain directory
	if err := loadEnvFile("chain/.env"); err != nil {
		// If chain/.env file doesn't exist, that's okay - use system environment variables
	}
}

func main() {
	// Load configuration from .env file or environment variables
	loadConfig()

	provider := zeus.NewChatProvider()

	// Create agents with proper chain flow
	// Step 1: Generate story outline
	storyOutline := blades.NewAgent(
		"story_outline_agent",
		blades.WithModel("llama-3.3-70b"),
		blades.WithProvider(provider),
		blades.WithInstructions("Generate a very short story outline based on the user's input. Keep it concise and clear."),
	)

	// Step 2: Analyze the outline
	storyAnalyzer := blades.NewAgent(
		"story_analyzer_agent",
		blades.WithModel("llama-3.3-70b"),
		blades.WithProvider(provider),
		blades.WithInstructions("Analyze the story outline provided. Judge its quality and determine if it's a sci-fi story. Provide a brief analysis."),
	)

	// Step 3: Write the story based on the ORIGINAL outline
	storyWriter := blades.NewAgent(
		"story_writer_agent",
		blades.WithModel("llama-3.3-70b"),
		blades.WithProvider(provider),
		blades.WithInstructions("Write a short story based on the story outline provided. Use the outline as your guide to create an engaging narrative."),
	)

	// Create the chain - it will automatically show beautiful visualization!
	chain := flow.NewChain(storyOutline, storyAnalyzer, storyWriter)

	// Initial prompt
	prompt := blades.NewPrompt(
		blades.UserMessage("A brave knight embarks on a quest to find a hidden treasure."),
	)

	// Run the chain - all visualization happens automatically!
	result, err := chain.Run(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}

	// The result is already printed by the chain, but you can access it if needed
	_ = result.Text()
}
