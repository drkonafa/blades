package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/blades"
	"github.com/go-kratos/blades/flow"
)

// FlowExecutor wraps the flow.Chain with beautiful execution visualization
type FlowExecutor struct {
	chain *flow.Chain
	steps []ChainStep
}

// NewFlowExecutor creates a new flow executor that can work with any chain
func NewFlowExecutor(chain *flow.Chain, steps []ChainStep) *FlowExecutor {
	return &FlowExecutor{
		chain: chain,
		steps: steps,
	}
}

// ExecuteWithVisualization runs the chain with beautiful output
func (fe *FlowExecutor) ExecuteWithVisualization(ctx context.Context, prompt *blades.Prompt) (*blades.Generation, error) {
	totalSteps := len(fe.steps)
	
	// Print header
	fe.printHeader(totalSteps)
	
	// Print initial prompt
	fe.printSection("INITIAL PROMPT", ColorCyan, prompt.String())
	
	// Execute the chain step by step for visualization
	var currentPrompt = prompt
	var finalResult *blades.Generation
	
	for i, step := range fe.steps {
		stepNum := i + 1
		
		// Print progress bar
		fe.printProgressBar(stepNum, totalSteps)
		
		// Print step header
		fe.printStepHeader(stepNum, step.Name, step.Instructions)
		
		// Print input
		fe.printInput(currentPrompt.String())
		
		// Execute step
		start := time.Now()
		result, err := step.Agent.Run(ctx, currentPrompt)
		if err != nil {
			fe.printError(err)
			return nil, err
		}
		duration := time.Since(start)
		
		// Print output
		fe.printOutput(result.Text(), duration)
		
		// Update prompt for next step
		currentPrompt = blades.NewPrompt(result.Messages...)
		finalResult = result
		
		// Add separator between steps
		if i < totalSteps-1 {
			fe.printSeparator()
		}
	}
	
	// Print final result
	fe.printFinalResult(finalResult.Text())
	
	return finalResult, nil
}

// Execute runs the chain normally (without visualization)
func (fe *FlowExecutor) Execute(ctx context.Context, prompt *blades.Prompt) (*blades.Generation, error) {
	return fe.chain.Run(ctx, prompt)
}

func (fe *FlowExecutor) printHeader(totalSteps int) {
	fmt.Printf("\n%s%s╔════════════════════════════════════════════════════════════════════════════════╗%s\n", ColorBold, ColorBlue, ColorReset)
	fmt.Printf("%s%s║%s %sCHAIN EXECUTION STARTED%s %s│ Steps: %d%s %s║%s\n", ColorBold, ColorBlue, ColorReset, ColorBold, ColorWhite, ColorReset, ColorYellow, totalSteps, ColorBold, ColorBlue, ColorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════════════════════════╝%s\n\n", ColorBold, ColorBlue, ColorReset)
}

func (fe *FlowExecutor) printProgressBar(current, total int) {
	width := 50
	filled := int(float64(current) / float64(total) * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	percentage := int(float64(current) / float64(total) * 100)
	
	fmt.Printf("%s[%s%s%s] %d%% (%d/%d)%s\n", 
		ColorYellow, bar, ColorReset, ColorYellow, percentage, current, total, ColorReset)
}

func (fe *FlowExecutor) printStepHeader(stepNum int, name, instructions string) {
	fmt.Printf("\n%s%s┌─ STEP %d: %s ─────────────────────────────────────────────────────────────┐%s\n", 
		ColorBold, ColorGreen, stepNum, strings.ToUpper(name), ColorReset)
	fmt.Printf("%s%s│%s Instructions: %s%s%s\n", 
		ColorBold, ColorGreen, ColorReset, ColorWhite, instructions, ColorReset)
	fmt.Printf("%s%s└─────────────────────────────────────────────────────────────────────────────┘%s\n", 
		ColorBold, ColorGreen, ColorReset)
}

func (fe *FlowExecutor) printInput(input string) {
	fmt.Printf("\n%s%s📥 INPUT:%s\n", ColorBold, ColorBlue, ColorReset)
	fe.printText(input, ColorBlue)
}

func (fe *FlowExecutor) printOutput(output string, duration time.Duration) {
	fmt.Printf("\n%s%s📤 OUTPUT:%s %s(%.2fs)%s\n", ColorBold, ColorGreen, ColorReset, ColorYellow, duration.Seconds(), ColorReset)
	fe.printText(output, ColorGreen)
}

func (fe *FlowExecutor) printText(text string, color string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			fmt.Printf("%s%s%s\n", color, line, ColorReset)
		} else {
			fmt.Printf("\n")
		}
	}
}

func (fe *FlowExecutor) printSeparator() {
	fmt.Printf("\n%s%s─────────────────────────────────────────────────────────────────────────────%s\n", ColorPurple, ColorBold, ColorReset)
}

func (fe *FlowExecutor) printError(err error) {
	fmt.Printf("\n%s%s❌ ERROR: %s%s\n", ColorBold, ColorRed, err.Error(), ColorReset)
}

func (fe *FlowExecutor) printFinalResult(result string) {
	fmt.Printf("\n%s%s╔════════════════════════════════════════════════════════════════════════════════╗%s\n", ColorBold, ColorGreen, ColorReset)
	fmt.Printf("%s%s║%s %s🎉 CHAIN EXECUTION COMPLETE! 🎉%s %s║%s\n", ColorBold, ColorGreen, ColorReset, ColorBold, ColorWhite, ColorBold, ColorGreen, ColorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════════════════════════╝%s\n", ColorBold, ColorGreen, ColorReset)
	
	fmt.Printf("\n%s%s📋 FINAL RESULT:%s\n", ColorBold, ColorCyan, ColorReset)
	fe.printText(result, ColorCyan)
}

func (fe *FlowExecutor) printSection(title, color, content string) {
	fmt.Printf("\n%s%s%s%s%s\n", ColorBold, color, title, ColorReset, ColorBold)
	fe.printText(content, color)
}
