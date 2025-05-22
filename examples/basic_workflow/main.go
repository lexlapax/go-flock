// ABOUTME: This example demonstrates basic workflow usage in the go-flock library.
// ABOUTME: Shows how to create and execute workflows that orchestrate multiple agents.

package main

import (
	"context"
	"fmt"
	"log"
)

// ExampleWorkflow demonstrates a simple workflow implementation
type ExampleWorkflow struct {
	agents []interface{} // Placeholder for agents
}

func (w *ExampleWorkflow) Name() string {
	return "example_workflow"
}

func (w *ExampleWorkflow) Description() string {
	return "A simple example workflow that coordinates multiple steps"
}

func (w *ExampleWorkflow) Execute(ctx context.Context, config map[string]interface{}) (interface{}, error) {
	steps := []string{"Step 1: Initialize", "Step 2: Process", "Step 3: Finalize"}
	results := make([]string, len(steps))

	for i, step := range steps {
		results[i] = fmt.Sprintf("Completed: %s", step)
	}

	return map[string]interface{}{
		"workflow": w.Name(),
		"steps":    results,
		"status":   "completed",
	}, nil
}

func (w *ExampleWorkflow) Agents() []interface{} {
	return w.agents
}

func main() {
	fmt.Println("Basic Workflow Example")
	fmt.Println("======================")

	// Create a workflow instance
	workflow := &ExampleWorkflow{
		agents: []interface{}{}, // Agents would be added here
	}

	// Execute the workflow
	ctx := context.Background()
	config := map[string]interface{}{
		"input": "Example workflow input",
	}

	result, err := workflow.Execute(ctx, config)
	if err != nil {
		log.Fatalf("Workflow execution failed: %v", err)
	}

	fmt.Printf("Workflow: %s\n", workflow.Name())
	fmt.Printf("Description: %s\n", workflow.Description())
	fmt.Printf("Result: %+v\n", result)
}
