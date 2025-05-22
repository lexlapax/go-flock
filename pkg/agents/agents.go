// ABOUTME: This package provides intelligent agents that combine tools, prompts, and logic into cohesive units.
// ABOUTME: Agents wrap functionality to create autonomous units capable of completing specific tasks.

package agents

import (
	"context"

	"github.com/lexlapax/go-flock/pkg/common"
	"github.com/lexlapax/go-flock/pkg/tools"
	agentDomain "github.com/lexlapax/go-llms/pkg/agent/domain"
)

// FlockAgent extends the base go-llms Agent interface with multi-agent coordination capabilities
type FlockAgent interface {
	agentDomain.Agent // Embed base Agent interface

	// AddFlockTool registers a flock-specific tool with the agent
	AddFlockTool(tool tools.FlockTool) FlockAgent
	// SetCoordinator configures how this agent coordinates with others
	SetCoordinator(coordinator Coordinator) FlockAgent
	// ExecuteTask runs a structured task
	ExecuteTask(ctx context.Context, task common.Task) (common.TaskResult, error)
	// ExecuteWorkflow runs a complex workflow (workflow interface defined in workflows package)
	ExecuteWorkflow(ctx context.Context, workflow interface{}) (interface{}, error)
	// GetCapabilities returns what this agent can do
	GetCapabilities() []string
}

// Coordinator manages multi-agent coordination and task distribution
type Coordinator interface {
	// CoordinateAgents orchestrates multiple agents working together
	CoordinateAgents(ctx context.Context, task common.Task, agents []FlockAgent) (common.CoordinationResult, error)
	// DistributeTasks breaks down a main task into subtasks for different agents
	DistributeTasks(ctx context.Context, mainTask common.Task) ([]common.SubTask, error)
	// AggregateResults combines results from multiple agents
	AggregateResults(ctx context.Context, results []common.TaskResult) (interface{}, error)
	// ResolveConflicts handles disagreements between agents
	ResolveConflicts(ctx context.Context, conflictingResults []common.TaskResult) (common.TaskResult, error)
}

// AgentRole defines the role of an agent in multi-agent scenarios
type AgentRole string

const (
	RoleCoordinator AgentRole = "coordinator"
	RoleSpecialist  AgentRole = "specialist"
	RoleValidator   AgentRole = "validator"
	RoleAggregator  AgentRole = "aggregator"
)

// AgentCapability defines what an agent can do
type AgentCapability string

const (
	CapabilityAnalysis     AgentCapability = "analysis"
	CapabilityGeneration   AgentCapability = "generation"
	CapabilityValidation   AgentCapability = "validation"
	CapabilityCoordination AgentCapability = "coordination"
	CapabilitySpecialized  AgentCapability = "specialized"
)

// AgentRegistry manages available agents
type AgentRegistry interface {
	// Register adds an agent to the registry
	Register(name string, agent FlockAgent) error
	// Get retrieves an agent by name
	Get(name string) (FlockAgent, error)
	// List returns all registered agents
	List() []FlockAgent
	// ListByCapability returns agents with specific capabilities
	ListByCapability(capability AgentCapability) []FlockAgent
	// ListByRole returns agents with specific roles
	ListByRole(role AgentRole) []FlockAgent
}

// AgentBuilder helps construct agents with proper configuration
type AgentBuilder interface {
	// WithName sets the agent name
	WithName(name string) AgentBuilder
	// WithRole sets the agent's role
	WithRole(role AgentRole) AgentBuilder
	// WithCapabilities sets what the agent can do
	WithCapabilities(capabilities []AgentCapability) AgentBuilder
	// WithTools adds tools to the agent
	WithTools(tools []tools.FlockTool) AgentBuilder
	// WithSystemPrompt sets the agent's system prompt
	WithSystemPrompt(prompt string) AgentBuilder
	// WithModel specifies which LLM model to use
	WithModel(modelName string) AgentBuilder
	// WithCoordinator sets the coordination strategy
	WithCoordinator(coordinator Coordinator) AgentBuilder
	// Build creates the final agent
	Build() (FlockAgent, error)
}
