package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/blades"
)

// Colors for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// ChainStep represents a single step in the chain
type ChainStep struct {
	Name         string
	Instructions string
	Agent        blades.Runner
}

// ChainExecutor handles the execution and visualization of chains
type ChainExecutor struct {
	steps []ChainStep
}

// NewChainExecutor creates a new chain executor
func NewChainExecutor() *ChainExecutor {
	return &ChainExecutor{
		steps: make([]ChainStep, 0),
	}
}

// AddStep adds a step to the chain
func (ce *ChainExecutor) AddStep(name, instructions string, agent blades.Runner) {
	ce.steps = append(ce.steps, ChainStep{
		Name:         name,
		Instructions: instructions,
		Agent:        agent,
	})
}

// Execute runs the chain with beautiful output
func (ce *ChainExecutor) Execute(ctx context.Context, initialPrompt *blades.Prompt) (*blades.Generation, error) {
	totalSteps := len(ce.steps)
	
	// Print header
	ce.printHeader(totalSteps)
	
	// Print initial prompt
	fmt.Printf("\n%s%sINITIAL PROMPT%s\n", ColorBold, ColorCyan, ColorReset)
	ce.printText(initialPrompt.String(), ColorCyan)
	
	var currentPrompt = initialPrompt
	var finalResult *blades.Generation
	
	// Execute each step
	for i, step := range ce.steps {
		stepNum := i + 1
		
		// Print progress bar
		ce.printProgressBar(stepNum, totalSteps)
		
		// Print step header
		ce.printStepHeader(stepNum, step.Name, step.Instructions)
		
		// Print input
		ce.printInput(currentPrompt.String())
		
		// Execute step
		start := time.Now()
		result, err := step.Agent.Run(ctx, currentPrompt)
		if err != nil {
			ce.printError(err)
			return nil, err
		}
		duration := time.Since(start)
		
		// Print output
		ce.printOutput(result.Text(), duration)
		
		// Update prompt for next step
		currentPrompt = blades.NewPrompt(result.Messages...)
		finalResult = result
		
		// Add separator between steps
		if i < totalSteps-1 {
			ce.printSeparator()
		}
	}
	
	// Print final result
	ce.printFinalResult(finalResult.Text())
	
	return finalResult, nil
}

func (ce *ChainExecutor) printHeader(totalSteps int) {
	fmt.Printf("\n%s%s╔════════════════════════════════════════════════════════════════════════════════╗%s\n", ColorBold, ColorBlue, ColorReset)
	fmt.Printf("%s%s║%s %sCHAIN EXECUTION STARTED%s %s│ Steps: %d%s %s║%s\n", ColorBold, ColorBlue, ColorReset, ColorBold, ColorWhite, ColorReset, ColorYellow, totalSteps, ColorBold, ColorBlue, ColorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════════════════════════╝%s\n\n", ColorBold, ColorBlue, ColorReset)
}

func (ce *ChainExecutor) printProgressBar(current, total int) {
	width := 50
	filled := int(float64(current) / float64(total) * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	percentage := int(float64(current) / float64(total) * 100)
	
	fmt.Printf("%s[%s%s%s] %d%% (%d/%d)%s\n", 
		ColorYellow, bar, ColorReset, ColorYellow, percentage, current, total, ColorReset)
}

func (ce *ChainExecutor) printStepHeader(stepNum int, name, instructions string) {
	fmt.Printf("\n%s%s┌─ STEP %d: %s ─────────────────────────────────────────────────────────────┐%s\n", 
		ColorBold, ColorGreen, stepNum, strings.ToUpper(name), ColorReset)
	fmt.Printf("%s%s│%s Instructions: %s%s%s\n", 
		ColorBold, ColorGreen, ColorReset, ColorWhite, instructions, ColorReset)
	fmt.Printf("%s%s└─────────────────────────────────────────────────────────────────────────────┘%s\n", 
		ColorBold, ColorGreen, ColorReset)
}

func (ce *ChainExecutor) printInput(input string) {
	fmt.Printf("\n%s%s📥 INPUT:%s\n", ColorBold, ColorBlue, ColorReset)
	ce.printText(input, ColorBlue)
}

func (ce *ChainExecutor) printOutput(output string, duration time.Duration) {
	fmt.Printf("\n%s%s📤 OUTPUT:%s %s(%.2fs)%s\n", ColorBold, ColorGreen, ColorReset, ColorYellow, duration.Seconds(), ColorReset)
	ce.printText(output, ColorGreen)
}

func (ce *ChainExecutor) printText(text string, color string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			fmt.Printf("%s%s%s\n", color, line, ColorReset)
		} else {
			fmt.Printf("\n")
		}
	}
}

func (ce *ChainExecutor) printSeparator() {
	fmt.Printf("\n%s%s─────────────────────────────────────────────────────────────────────────────%s\n", ColorPurple, ColorBold, ColorReset)
}

func (ce *ChainExecutor) printError(err error) {
	fmt.Printf("\n%s%s❌ ERROR: %s%s\n", ColorBold, ColorRed, err.Error(), ColorReset)
}

func (ce *ChainExecutor) printFinalResult(result string) {
	fmt.Printf("\n%s%s╔════════════════════════════════════════════════════════════════════════════════╗%s\n", ColorBold, ColorGreen, ColorReset)
	fmt.Printf("%s%s║%s %s🎉 CHAIN EXECUTION COMPLETE! 🎉%s %s║%s\n", ColorBold, ColorGreen, ColorReset, ColorBold, ColorWhite, ColorBold, ColorGreen, ColorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════════════════════════╝%s\n", ColorBold, ColorGreen, ColorReset)
	
	fmt.Printf("\n%s%s📋 FINAL RESULT:%s\n", ColorBold, ColorCyan, ColorReset)
	ce.printText(result, ColorCyan)
}
