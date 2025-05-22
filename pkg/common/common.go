// ABOUTME: This package provides shared utilities, interfaces, and common functionality across the go-flock library.
// ABOUTME: Common contains reusable components that support tools, agents, and workflows.

package common

import (
	"context"

	agentDomain "github.com/lexlapax/go-llms/pkg/agent/domain"
)

// Config represents configuration options for flock components
type Config struct {
	LLMProvider string                 `json:"llm_provider"`
	Settings    map[string]interface{} `json:"settings"`
}

// Logger defines the logging interface for flock components
type Logger interface {
	Debug(ctx context.Context, msg string, args ...interface{})
	Info(ctx context.Context, msg string, args ...interface{})
	Warn(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, msg string, args ...interface{})
}

// TaskType defines the type of task being executed
type TaskType string

const (
	TaskTypeAnalysis     TaskType = "analysis"
	TaskTypeGeneration   TaskType = "generation"
	TaskTypeTransform    TaskType = "transform"
	TaskTypeCoordination TaskType = "coordination"
)

// Task represents a unit of work that can be executed by an agent
type Task interface {
	// ID returns the unique identifier for this task
	ID() string
	// Type returns the type of task
	Type() TaskType
	// Requirements returns what the task needs to execute
	Requirements() TaskRequirements
	// Execute runs the task with the given agent
	Execute(ctx context.Context, agent agentDomain.Agent) (TaskResult, error)
}

// TaskRequirements defines what resources a task needs
type TaskRequirements struct {
	Tools        []string               `json:"tools"`
	Capabilities []string               `json:"capabilities"`
	Resources    map[string]interface{} `json:"resources"`
}

// TaskResult represents the outcome of task execution
type TaskResult struct {
	TaskID   string                 `json:"task_id"`
	Success  bool                   `json:"success"`
	Result   interface{}            `json:"result"`
	Error    error                  `json:"error,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SubTask represents a decomposed portion of a larger task
type SubTask struct {
	ID           string           `json:"id"`
	ParentID     string           `json:"parent_id"`
	Description  string           `json:"description"`
	Input        interface{}      `json:"input"`
	Dependencies []string         `json:"dependencies"`
	Requirements TaskRequirements `json:"requirements"`
}

// CoordinationResult represents the outcome of multi-agent coordination
type CoordinationResult struct {
	Success      bool                   `json:"success"`
	Results      []TaskResult           `json:"results"`
	Summary      string                 `json:"summary"`
	Coordination map[string]interface{} `json:"coordination"`
}
