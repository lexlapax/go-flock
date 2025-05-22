// ABOUTME: This package provides workflow orchestration for complex multi-step automation tasks.
// ABOUTME: Workflows coordinate multiple agents and tools to accomplish sophisticated automation scenarios.

package workflows

import (
	"context"
	"time"

	"github.com/lexlapax/go-flock/pkg/agents"
	schemaDomain "github.com/lexlapax/go-llms/pkg/schema/domain"
)

// Workflow represents a complex multi-step process that coordinates agents and tools
type Workflow interface {
	// Name returns the name of the workflow
	Name() string
	// Description returns a description of what the workflow accomplishes
	Description() string
	// Execute runs the workflow with the given input and agents
	Execute(ctx context.Context, input WorkflowInput, flockAgent agents.FlockAgent) (WorkflowResult, error)
	// GetSchema returns the schema for workflow input validation
	GetSchema() *schemaDomain.Schema
	// GetSteps returns the defined workflow steps
	GetSteps() []WorkflowStep
	// Validate checks if the workflow configuration is valid
	Validate() error
}

// WorkflowInput represents the input to a workflow execution
type WorkflowInput struct {
	Data       interface{}            `json:"data"`
	Parameters map[string]interface{} `json:"parameters"`
	Context    map[string]interface{} `json:"context"`
	Options    WorkflowOptions        `json:"options"`
}

// WorkflowOptions configure workflow execution behavior
type WorkflowOptions struct {
	Timeout         time.Duration   `json:"timeout"`
	MaxRetries      int             `json:"max_retries"`
	FailureStrategy FailureStrategy `json:"failure_strategy"`
	Parallel        bool            `json:"parallel"`
	StepTimeout     time.Duration   `json:"step_timeout"`
}

// WorkflowResult represents the outcome of workflow execution
type WorkflowResult struct {
	WorkflowName string                 `json:"workflow_name"`
	Success      bool                   `json:"success"`
	Result       interface{}            `json:"result"`
	Steps        []StepResult           `json:"steps"`
	Duration     time.Duration          `json:"duration"`
	Error        error                  `json:"error,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	AgentRole    agents.AgentRole         `json:"agent_role"`
	Capabilities []agents.AgentCapability `json:"capabilities"`
	Input        interface{}              `json:"input"`
	Dependencies []string                 `json:"dependencies"`
	Condition    StepCondition            `json:"condition,omitempty"`
	Retryable    bool                     `json:"retryable"`
	Timeout      time.Duration            `json:"timeout,omitempty"`
}

// StepResult represents the outcome of a workflow step
type StepResult struct {
	StepID   string        `json:"step_id"`
	Success  bool          `json:"success"`
	Result   interface{}   `json:"result"`
	Duration time.Duration `json:"duration"`
	Attempt  int           `json:"attempt"`
	Error    error         `json:"error,omitempty"`
	Skipped  bool          `json:"skipped,omitempty"`
}

// StepCondition defines when a step should execute
type StepCondition struct {
	Type       ConditionType          `json:"type"`
	Expression string                 `json:"expression,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// ConditionType defines types of step conditions
type ConditionType string

const (
	ConditionAlways    ConditionType = "always"
	ConditionOnSuccess ConditionType = "on_success"
	ConditionOnFailure ConditionType = "on_failure"
	ConditionCustom    ConditionType = "custom"
)

// FailureStrategy defines how to handle workflow failures
type FailureStrategy string

const (
	FailureStopImmediately FailureStrategy = "stop_immediately"
	FailureContinue        FailureStrategy = "continue"
	FailureRetry           FailureStrategy = "retry"
	FailureRollback        FailureStrategy = "rollback"
)

// WorkflowEngine coordinates workflow execution
type WorkflowEngine interface {
	// RegisterWorkflow adds a workflow to the engine
	RegisterWorkflow(workflow Workflow) error
	// ExecuteWorkflow runs a registered workflow
	ExecuteWorkflow(ctx context.Context, name string, input WorkflowInput, agent agents.FlockAgent) (WorkflowResult, error)
	// ListWorkflows returns all registered workflows
	ListWorkflows() []Workflow
	// GetWorkflow retrieves a workflow by name
	GetWorkflow(name string) (Workflow, error)
	// ValidateWorkflow checks if a workflow is properly configured
	ValidateWorkflow(name string) error
}

// WorkflowBuilder helps construct workflows with proper validation
type WorkflowBuilder interface {
	// WithName sets the workflow name
	WithName(name string) WorkflowBuilder
	// WithDescription sets the workflow description
	WithDescription(desc string) WorkflowBuilder
	// WithSchema sets the input schema
	WithSchema(schema *schemaDomain.Schema) WorkflowBuilder
	// AddStep adds a step to the workflow
	AddStep(step WorkflowStep) WorkflowBuilder
	// WithOptions sets workflow execution options
	WithOptions(options WorkflowOptions) WorkflowBuilder
	// Build creates the final workflow
	Build() (Workflow, error)
}

// WorkflowTemplate represents a reusable workflow pattern
type WorkflowTemplate interface {
	// Name returns the template name
	Name() string
	// Parameters returns required template parameters
	Parameters() []TemplateParameter
	// Generate creates a workflow from the template
	Generate(ctx context.Context, params map[string]interface{}) (Workflow, error)
}

// TemplateParameter defines a parameter for workflow templates
type TemplateParameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description"`
}
