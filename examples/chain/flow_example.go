package main

import (
	"context"
	"log"

	"github.com/go-kratos/blades"
	"github.com/go-kratos/blades/contrib/gemini"
	"github.com/go-kratos/blades/flow"
)

func main() {
	// Load configuration from .env file or environment variables
	loadConfig()

	provider := gemini.NewChatProvider()

	// Create agents
	storyOutline := blades.NewAgent(
		"story_outline_agent",
		blades.WithModel("gemini-2.0-flash"),
		blades.WithProvider(provider),
		blades.WithInstructions("Generate a very short story outline based on the user's input."),
	)
	storyChecker := blades.NewAgent(
		"outline_checker_agent",
		blades.WithModel("gemini-2.0-flash"),
		blades.WithProvider(provider),
		blades.WithInstructions("Read the given story outline, and judge the quality. Also, determine if it is a scifi story."),
	)
	storyAgent := blades.NewAgent(
		"story_agent",
		blades.WithModel("gemini-2.0-flash"),
		blades.WithProvider(provider),
		blades.WithInstructions("Write a short story based on the given outline."),
	)

	// Create the chain using flow package
	chain := flow.NewChain(storyOutline, storyChecker, storyAgent)

	// Define the steps for visualization
	steps := []ChainStep{
		{
			Name:         "Story Outline Generator",
			Instructions: "Generate a very short story outline based on the user's input.",
			Agent:        storyOutline,
		},
		{
			Name:         "Quality Checker",
			Instructions: "Read the given story outline, and judge the quality. Also, determine if it is a scifi story.",
			Agent:        storyChecker,
		},
		{
			Name:         "Story Writer",
			Instructions: "Write a short story based on the given outline.",
			Agent:        storyAgent,
		},
	}

	// Create flow executor
	executor := NewFlowExecutor(chain, steps)

	// Initial prompt
	prompt := blades.NewPrompt(
		blades.UserMessage("A brave knight embarks on a quest to find a hidden treasure."),
	)

	// Execute the chain with beautiful visualization
	_, err := executor.ExecuteWithVisualization(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}
}
