// ABOUTME: This file defines common types and options for all agents in the go-flock system.
// ABOUTME: It provides output format configuration and agent options used across all agent implementations.

package agents

// OutputFormat defines the output format for agent responses
type OutputFormat string

const (
	OutputFormatMarkdown OutputFormat = "markdown" // Default format
	OutputFormatJSON     OutputFormat = "json"
	OutputFormatText     OutputFormat = "text"
)

// AgentOptions configures agent behavior
type AgentOptions struct {
	OutputFormat OutputFormat // Output format for responses
	Model        string       // LLM model to use (optional, uses provider default if empty)
}

// DefaultAgentOptions returns the default options for agents
func DefaultAgentOptions() AgentOptions {
	return AgentOptions{
		OutputFormat: OutputFormatMarkdown,
	}
}
