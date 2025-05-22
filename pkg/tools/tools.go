// ABOUTME: This package provides individual LLM-callable tools for various automation tasks.
// ABOUTME: Tools are atomic functions that can be invoked by LLMs to perform specific operations.

package tools

import (
	"context"

	agentDomain "github.com/lexlapax/go-llms/pkg/agent/domain"
	schemaDomain "github.com/lexlapax/go-llms/pkg/schema/domain"
)

// FlockTool extends the base go-llms Tool interface with additional flock-specific functionality
type FlockTool interface {
	agentDomain.Tool // Embed base Tool interface

	// Category returns the functional category of this tool
	Category() ToolCategory
	// Dependencies returns other tools this tool depends on
	Dependencies() []string
	// IsAsync returns true if this tool can run asynchronously
	IsAsync() bool
	// Validate validates parameters before execution
	Validate(ctx context.Context, params interface{}) error
}

// ToolCategory defines functional categories for tools
type ToolCategory string

const (
	CategoryFileSystem   ToolCategory = "filesystem"
	CategoryNetwork      ToolCategory = "network"
	CategoryData         ToolCategory = "data"
	CategoryComputation  ToolCategory = "computation"
	CategoryCoordination ToolCategory = "coordination"
	CategoryMonitoring   ToolCategory = "monitoring"
)

// ToolRegistry manages available tools
type ToolRegistry interface {
	// Register adds a tool to the registry
	Register(tool FlockTool) error
	// Get retrieves a tool by name
	Get(name string) (FlockTool, error)
	// List returns all registered tools
	List() []FlockTool
	// ListByCategory returns tools in a specific category
	ListByCategory(category ToolCategory) []FlockTool
	// Validate checks if all tool dependencies are available
	Validate() error
}

// ToolBuilder helps construct tools with validation
type ToolBuilder interface {
	// WithName sets the tool name
	WithName(name string) ToolBuilder
	// WithDescription sets the tool description
	WithDescription(desc string) ToolBuilder
	// WithCategory sets the tool category
	WithCategory(category ToolCategory) ToolBuilder
	// WithSchema sets the parameter schema
	WithSchema(schema *schemaDomain.Schema) ToolBuilder
	// WithExecutor sets the execution function
	WithExecutor(executor func(context.Context, interface{}) (interface{}, error)) ToolBuilder
	// Build creates the final tool
	Build() (FlockTool, error)
}
