# Workflow Examples

This directory contains examples demonstrating workflow orchestration in go-flock.

## Overview

Workflows coordinate multiple agents and tools to accomplish complex, multi-step tasks. They provide execution control, error handling, and result aggregation.

## Examples

### basic/
Demonstrates fundamental workflow patterns:
- Sequential task execution
- Parallel agent coordination
- Error handling and retry logic
- Result aggregation

## Running Examples

```bash
# Run the basic workflow example
cd basic
go run main.go
```

## Key Concepts

1. **Workflow Steps**: Define discrete units of work
2. **Dependencies**: Control execution order
3. **Conditions**: Execute steps based on previous results
4. **Failure Strategies**: Handle errors gracefully
5. **Parallel Execution**: Run independent steps concurrently

## Workflow Components

```go
// Workflow input configuration
type WorkflowInput struct {
    Data       interface{}
    Parameters map[string]interface{}
    Options    WorkflowOptions
}

// Execution options
type WorkflowOptions struct {
    Timeout         time.Duration
    MaxRetries      int
    FailureStrategy FailureStrategy
    Parallel        bool
}
```

## Next Steps

- Explore workflow implementations in `pkg/workflows/`
- See how agents work together in `docs/workflows/`
- Build custom workflows for your use cases