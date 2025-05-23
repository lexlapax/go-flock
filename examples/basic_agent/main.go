// ABOUTME: This example demonstrates creating agents that will be used in workflows.
// ABOUTME: Shows the pattern for specialized agents that go-flock will provide.

package main

import (
	"fmt"
)

// This example shows the pattern that go-flock agents will follow.
// Each agent in pkg/agents/ will be a ready-to-use implementation
// that extends go-llms agents with specific capabilities.

func main() {
	fmt.Println("go-flock Agent Pattern Example")
	fmt.Println("==============================")

	// In actual usage, you would import an agent like:
	// import "github.com/your-org/go-flock/pkg/agents/research"
	// agent := research.NewWebResearchAgent(provider)

	// Or for code analysis:
	// import "github.com/your-org/go-flock/pkg/agents/code"
	// agent := code.NewCodeReviewAgent(provider)

	// The agents would extend go-llms Agent interface
	// with pre-configured tools and system prompts

	fmt.Println("\ngo-flock will provide specialized agents:")
	fmt.Println("- research: Web research, fact checking, summarization")
	fmt.Println("- code: Code review, refactoring, documentation")
	fmt.Println("- data: Data analysis, visualization, reporting")
	fmt.Println("- devops: CI/CD, deployment, monitoring")
	fmt.Println("- testing: Test generation, validation, QA")
	fmt.Println("\nEach agent will come with appropriate tools and prompts pre-configured.")
}
