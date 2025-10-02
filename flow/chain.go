package flow

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/blades"
)

var (
	_ blades.Runner = (*Chain)(nil)
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

// Chain represents a sequence of Runnable runners that process input sequentially.
type Chain struct {
	runners []blades.Runner
	verbose bool
}

// NewChain creates a new Chain with the given runners.
func NewChain(runners ...blades.Runner) *Chain {
	return &Chain{
		runners: runners,
		verbose: true, // Enable verbose output by default
	}
}

// NewChainSilent creates a new Chain with verbose output disabled.
func NewChainSilent(runners ...blades.Runner) *Chain {
	return &Chain{
		runners: runners,
		verbose: false,
	}
}

// SetVerbose enables or disables verbose output.
func (c *Chain) SetVerbose(verbose bool) {
	c.verbose = verbose
}

// Run executes the chain of runners sequentially, passing the output of one as the input to the next.
func (c *Chain) Run(ctx context.Context, prompt *blades.Prompt, opts ...blades.ModelOption) (*blades.Generation, error) {
	if !c.verbose {
		return c.runSilent(ctx, prompt, opts...)
	}

	return c.runVerbose(ctx, prompt, opts...)
}

// runSilent executes the chain without verbose output.
func (c *Chain) runSilent(ctx context.Context, prompt *blades.Prompt, opts ...blades.ModelOption) (*blades.Generation, error) {
	var (
		err  error
		last *blades.Generation
	)
	for _, runner := range c.runners {
		last, err = runner.Run(ctx, prompt, opts...)
		if err != nil {
			return nil, err
		}
		prompt = blades.NewPrompt(last.Messages...)
	}
	return last, nil
}

// runVerbose executes the chain with beautiful visualization.
func (c *Chain) runVerbose(ctx context.Context, prompt *blades.Prompt, opts ...blades.ModelOption) (*blades.Generation, error) {
	totalSteps := len(c.runners)

	// Print header
	c.printHeader(totalSteps)

	// Print initial prompt
	fmt.Printf("\n%s%sINITIAL PROMPT%s\n", ColorBold, ColorCyan, ColorReset)
	c.printText(prompt.String(), ColorCyan)

	var currentPrompt = prompt
	var finalResult *blades.Generation

	// Execute each step
	for i, runner := range c.runners {
		stepNum := i + 1

		// Print progress bar
		c.printProgressBar(stepNum, totalSteps)

		// Get step info dynamically
		stepName, instructions := c.getStepInfo(runner, stepNum)

		// Print step header
		c.printStepHeader(stepNum, stepName, instructions)

		// Print input
		c.printInput(currentPrompt.String())

		// Execute step
		start := time.Now()
		result, err := runner.Run(ctx, currentPrompt, opts...)
		if err != nil {
			c.printError(err)
			return nil, err
		}
		duration := time.Since(start)

		// Print output
		c.printOutput(result.Text(), duration)

		// Update prompt for next step
		currentPrompt = blades.NewPrompt(result.Messages...)
		finalResult = result

		// Add separator between steps
		if i < totalSteps-1 {
			c.printSeparator()
		}
	}

	// Print final result
	c.printFinalResult(finalResult.Text())

	return finalResult, nil
}

// RunStream executes the chain of runners sequentially, streaming the output of the last runner.
func (c *Chain) RunStream(ctx context.Context, prompt *blades.Prompt, opts ...blades.ModelOption) (blades.Streamer[*blades.Generation], error) {
	pipe := blades.NewStreamPipe[*blades.Generation]()
	pipe.Go(func() error {
		for _, runner := range c.runners {
			last, err := runner.Run(ctx, prompt, opts...)
			if err != nil {
				return err
			}
			pipe.Send(last)
			prompt = blades.NewPrompt(last.Messages...)
		}
		return nil
	})
	return pipe, nil
}

// getStepInfo extracts step name and instructions from a runner (Agent)
func (c *Chain) getStepInfo(runner blades.Runner, stepNum int) (string, string) {
	// Try to get info from Agent if it's an Agent type
	if agent, ok := runner.(*blades.Agent); ok {
		name := agent.Name()
		instructions := agent.Instructions()

		// Use defaults if empty
		if name == "" {
			name = fmt.Sprintf("Agent %d", stepNum)
		}
		if instructions == "" {
			instructions = "Processing request..."
		}

		return name, instructions
	}

	// Fallback for other runner types
	return fmt.Sprintf("Step %d", stepNum), "Executing task..."
}

// printHeader prints the chain execution header
func (c *Chain) printHeader(totalSteps int) {
	fmt.Printf("\n%s%s╔════════════════════════════════════════════════════════════════════════════════╗%s\n", ColorBold, ColorBlue, ColorReset)
	fmt.Printf("%s%s║%s %sCHAIN EXECUTION STARTED%s %s│ Steps: %d%s %s║%s\n", ColorBold, ColorBlue, ColorReset, ColorBold, ColorWhite, ColorReset, ColorYellow, totalSteps, ColorBold, ColorBlue, ColorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════════════════════════╝%s\n\n", ColorBold, ColorBlue, ColorReset)
}

// printProgressBar prints a progress bar
func (c *Chain) printProgressBar(current, total int) {
	width := 50
	filled := int(float64(current) / float64(total) * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	percentage := int(float64(current) / float64(total) * 100)

	fmt.Printf("%s[%s%s%s] %d%% (%d/%d)%s\n",
		ColorYellow, bar, ColorReset, ColorYellow, percentage, current, total, ColorReset)
}

// printStepHeader prints the step header
func (c *Chain) printStepHeader(stepNum int, name, instructions string) {
	fmt.Printf("\n%s%s┌─ STEP %d: %s ─────────────────────────────────────────────────────────────┐%s\n",
		ColorBold, ColorGreen, stepNum, strings.ToUpper(name), ColorReset)
	fmt.Printf("%s%s│%s Instructions: %s%s%s\n",
		ColorBold, ColorGreen, ColorReset, ColorWhite, instructions, ColorReset)
	fmt.Printf("%s%s└─────────────────────────────────────────────────────────────────────────────┘%s\n",
		ColorBold, ColorGreen, ColorReset)
}

// printInput prints the input
func (c *Chain) printInput(input string) {
	fmt.Printf("\n%s%s📥 INPUT:%s\n", ColorBold, ColorBlue, ColorReset)
	c.printText(input, ColorBlue)
}

// printOutput prints the output
func (c *Chain) printOutput(output string, duration time.Duration) {
	fmt.Printf("\n%s%s📤 OUTPUT:%s %s(%.2fs)%s\n", ColorBold, ColorGreen, ColorReset, ColorYellow, duration.Seconds(), ColorReset)
	c.printText(output, ColorGreen)
}

// printText prints text with color
func (c *Chain) printText(text string, color string) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			fmt.Printf("%s%s%s\n", color, line, ColorReset)
		} else {
			fmt.Printf("\n")
		}
	}
}

// printSeparator prints a separator between steps
func (c *Chain) printSeparator() {
	fmt.Printf("\n%s%s─────────────────────────────────────────────────────────────────────────────%s\n", ColorPurple, ColorBold, ColorReset)
}

// printError prints an error
func (c *Chain) printError(err error) {
	fmt.Printf("\n%s%s❌ ERROR: %s%s\n", ColorBold, ColorRed, err.Error(), ColorReset)
}

// printFinalResult prints the final result
func (c *Chain) printFinalResult(result string) {
	fmt.Printf("\n%s%s╔════════════════════════════════════════════════════════════════════════════════╗%s\n", ColorBold, ColorGreen, ColorReset)
	fmt.Printf("%s%s║%s %s🎉 CHAIN EXECUTION COMPLETE! 🎉%s %s║%s\n", ColorBold, ColorGreen, ColorReset, ColorBold, ColorWhite, ColorBold, ColorGreen, ColorReset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════════════════════════╝%s\n", ColorBold, ColorGreen, ColorReset)

	fmt.Printf("\n%s%s📋 FINAL RESULT:%s\n", ColorBold, ColorCyan, ColorReset)
	c.printText(result, ColorCyan)
}
