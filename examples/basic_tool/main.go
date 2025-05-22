// ABOUTME: This example demonstrates basic tool usage in the go-flock library.
// ABOUTME: Shows how to create and execute individual tools independently.

package main

import (
	"context"
	"fmt"
	"log"
)

// ExampleTool demonstrates a simple tool implementation
type ExampleTool struct{}

func (t *ExampleTool) Name() string {
	return "example_tool"
}

func (t *ExampleTool) Description() string {
	return "A simple example tool that echoes input"
}

func (t *ExampleTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	input, ok := params["input"].(string)
	if !ok {
		return nil, fmt.Errorf("input parameter must be a string")
	}
	return fmt.Sprintf("Echo: %s", input), nil
}

func main() {
	fmt.Println("Basic Tool Example")
	fmt.Println("==================")

	// Create a tool instance
	tool := &ExampleTool{}

	// Execute the tool
	ctx := context.Background()
	params := map[string]interface{}{
		"input": "Hello, go-flock!",
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Fatalf("Tool execution failed: %v", err)
	}

	fmt.Printf("Tool: %s\n", tool.Name())
	fmt.Printf("Description: %s\n", tool.Description())
	fmt.Printf("Result: %v\n", result)
}
