package main

import (
	"context"
	"log"

	"github.com/go-kratos/blades"
	"github.com/go-kratos/blades/contrib/gemini"
)

// Simple example showing how to use the ChainExecutor with any agents
func main() {
	// Load configuration from .env file or environment variables
	loadConfig()

	provider := gemini.NewChatProvider()

	// Create any agents you want
	agent1 := blades.NewAgent(
		"analyzer",
		blades.WithModel("gemini-2.0-flash"),
		blades.WithProvider(provider),
		blades.WithInstructions("Analyze the given text and provide insights."),
	)

	agent2 := blades.NewAgent(
		"summarizer",
		blades.WithModel("gemini-2.0-flash"),
		blades.WithProvider(provider),
		blades.WithInstructions("Summarize the analysis in 3 key points."),
	)

	agent3 := blades.NewAgent(
		"enhancer",
		blades.WithModel("gemini-2.0-flash"),
		blades.WithProvider(provider),
		blades.WithInstructions("Enhance the summary with actionable recommendations."),
	)

	// Create chain executor
	executor := NewChainExecutor()

	// Add steps - this is completely agnostic!
	executor.AddStep(
		"Text Analyzer",
		"Analyze the given text and provide insights.",
		agent1,
	)
	executor.AddStep(
		"Summary Generator",
		"Summarize the analysis in 3 key points.",
		agent2,
	)
	executor.AddStep(
		"Enhancement Engine",
		"Enhance the summary with actionable recommendations.",
		agent3,
	)

	// Initial prompt
	prompt := blades.NewPrompt(
		blades.UserMessage("The future of artificial intelligence in healthcare"),
	)

	// Execute with beautiful visualization
	_, err := executor.Execute(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}
}
