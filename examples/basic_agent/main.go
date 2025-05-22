// ABOUTME: This example demonstrates basic agent usage in the go-flock library.
// ABOUTME: Shows how to create and execute agents that combine tools and logic.

package main

import (
	"context"
	"fmt"
	"log"
)

// ExampleAgent demonstrates a simple agent implementation
type ExampleAgent struct {
	tools []interface{} // Placeholder for tools
}

func (a *ExampleAgent) Name() string {
	return "example_agent"
}

func (a *ExampleAgent) Description() string {
	return "A simple example agent that processes text input"
}

func (a *ExampleAgent) Execute(ctx context.Context, input string) (string, error) {
	// Simple processing logic
	return fmt.Sprintf("Agent processed: %s", input), nil
}

func (a *ExampleAgent) Tools() []interface{} {
	return a.tools
}

func main() {
	fmt.Println("Basic Agent Example")
	fmt.Println("===================")

	// Create an agent instance
	agent := &ExampleAgent{
		tools: []interface{}{}, // Tools would be added here
	}

	// Execute the agent
	ctx := context.Background()
	input := "Process this text"

	result, err := agent.Execute(ctx, input)
	if err != nil {
		log.Fatalf("Agent execution failed: %v", err)
	}

	fmt.Printf("Agent: %s\n", agent.Name())
	fmt.Printf("Description: %s\n", agent.Description())
	fmt.Printf("Input: %s\n", input)
	fmt.Printf("Output: %s\n", result)
}
