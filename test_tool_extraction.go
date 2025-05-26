package main

import (
	"fmt"

	"github.com/lexlapax/go-llms/pkg/agent/workflow"
	"github.com/lexlapax/go-llms/pkg/llm/provider"
)

func main() {
	// Test content that mimics Gemini's response - note the arguments field is a JSON string
	testContent := "```json\n{\n  \"tool_calls\": [\n    {\n      \"function\": {\n        \"name\": \"research_paper_api\",\n        \"arguments\": \"{\\\"query\\\": \\\"test\\\", \\\"max_results\\\": 5}\"\n      }\n    }\n  ]\n}\n```"

	fmt.Println("Test content:")
	fmt.Println(testContent)
	fmt.Println()

	// Create a mock provider and agent to test extraction
	mockProvider := provider.NewMockProvider()
	defaultAgent := workflow.NewAgent(mockProvider)

	// Test the extraction
	toolNames, params, found := defaultAgent.ExtractMultipleToolCalls(testContent)

	fmt.Printf("Found tool calls: %v\n", found)
	fmt.Printf("Tool names: %v\n", toolNames)
	fmt.Printf("Parameters: %v\n", params)
}
