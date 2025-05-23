# go-flock

A Go library providing a collection of tools, agents, and workflows for LLM-powered automation tasks. The name "go-flock" evokes the concept of "a flock of golems" or a flock of go-llms working together to accomplish complex automation.

## Overview

go-flock is built on top of the [go-llms](https://github.com/lexlapax/go-llms) interface library and provides a comprehensive framework for creating LLM-callable tools, intelligent agents, and complex workflows. Each component can be used independently or combined to create sophisticated automation solutions.

## Architecture

The library is organized around three core concepts:

- **Tools**: Individual functions that can be called by LLMs to perform specific tasks with categorization, dependency management, and validation
- **Agents**: Intelligent units that wrap functionality, prompts, and tool calls into cohesive working units with multi-agent coordination capabilities
- **Workflows**: Orchestrate multiple agents and tools to accomplish complex, multi-step tasks with conditional execution, retry logic, and failure handling

### Core Interfaces

go-flock extends go-llms with specialized interfaces:

- **FlockTool**: Extends `agentDomain.Tool` with categories, dependencies, async execution, and validation
- **FlockAgent**: Extends `agentDomain.Agent` with multi-agent coordination, task execution, and workflow capabilities
- **Workflow**: Comprehensive workflow orchestration with steps, conditions, and execution strategies
- **Coordinator**: Manages multi-agent coordination, task distribution, and result aggregation

## Features

- **Multi-Agent Coordination**: Agents can work together with role-based coordination and conflict resolution
- **Workflow Orchestration**: Complex multi-step processes with conditional execution and failure strategies
- **Tool Ecosystem**: Categorized tools with dependency management and async execution support
- **Schema-Driven**: Full integration with go-llms schema validation and structured output
- **Modular Design**: Each tool, agent, and workflow can be used independently
- **Command Line Interface**: Explore and test functionality through the included CLI utility
- **Build Automation**: Comprehensive Makefile with testing, benchmarks, and development workflows
- **Example Usage**: Working examples demonstrating different usage patterns

## Installation

```bash
go get github.com/lexlapax/go-flock
```

## Quick Start

### Using Tools with Categories

```go
package main

import (
    "context"
    "github.com/lexlapax/go-flock/pkg/tools"
)

func main() {
    // FlockTool extends go-llms Tool interface
    var tool tools.FlockTool
    
    // Tools have categories and dependencies
    category := tool.Category() // e.g., tools.CategoryFileSystem
    deps := tool.Dependencies() // Required tool dependencies
    
    // Validate before execution
    if err := tool.Validate(ctx, params); err != nil {
        // Handle validation error
    }
    
    // Execute with context
    result, err := tool.Execute(ctx, params)
}
```

### Creating Multi-Agent Systems

```go
package main

import (
    "context"
    "github.com/lexlapax/go-flock/pkg/agents"
    "github.com/lexlapax/go-flock/pkg/common"
)

func main() {
    // FlockAgent extends go-llms Agent interface
    var agent agents.FlockAgent
    
    // Add coordinator for multi-agent scenarios
    coordinator := &MyCoordinator{}
    agent.SetCoordinator(coordinator)
    
    // Execute structured tasks
    task := &MyTask{}
    result, err := agent.ExecuteTask(ctx, task)
    
    // Get agent capabilities
    capabilities := agent.GetCapabilities()
}
```

### Building Complex Workflows

```go
package main

import (
    "context"
    "github.com/lexlapax/go-flock/pkg/workflows"
    "github.com/lexlapax/go-flock/pkg/agents"
)

func main() {
    // Create workflow with steps and conditions
    workflow := &MyWorkflow{}
    
    // Configure workflow input
    input := workflows.WorkflowInput{
        Data: "process this",
        Options: workflows.WorkflowOptions{
            FailureStrategy: workflows.FailureRetry,
            Parallel: true,
        },
    }
    
    // Execute with flock agent
    var flockAgent agents.FlockAgent
    result, err := workflow.Execute(ctx, input, flockAgent)
}
```

## Command Line Interface

The library includes a CLI utility for exploring and testing functionality:

```bash
# Build the CLI
make build

# Or install directly
go install github.com/lexlapax/go-flock/cmd/flock

# Check version
./bin/flock version

# View available commands
./bin/flock --help
```

## Development

### Prerequisites

- Go 1.24.3 or later
- [go-llms](https://github.com/lexlapax/go-llms) v0.2.6

### Building with Make

The project includes a comprehensive Makefile for development:

```bash
# Development workflow
make dev              # Clean, format, vet, test, build

# Building
make build            # Build flock CLI binary
make build-all        # Build all packages
make examples         # Build example programs

# Testing
make test             # Run unit tests
make test-coverage    # Generate coverage report
make test-integration # Run integration tests
make bench            # Run benchmarks

# Code Quality
make check            # Format, vet, lint, and test
make fmt              # Format Go code
make vet              # Run go vet
make lint             # Run golangci-lint

# Dependencies
make deps             # Download and tidy dependencies
make deps-update      # Update to latest versions

# Tools
make tools            # Install development tools
make clean            # Clean build artifacts
```

### Manual Building

```bash
# Build all packages
go build ./...

# Build CLI
go build ./cmd/flock

# Run tests
go test ./...

# Format code
go fmt ./...
```

### Testing

The project follows Test-Driven Development (TDD) practices:

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run benchmarks
make bench

# Run specific package tests
go test ./pkg/tools/...
```

## Directory Structure

```
go-flock/
├── Makefile                # Build automation
├── cmd/                    # Command line applications
│   └── flock/             # Main CLI application
├── pkg/                   # Library packages
│   ├── tools/             # Extended tool interfaces with categories
│   ├── agents/            # Multi-agent coordination interfaces
│   ├── workflows/         # Workflow orchestration with conditions
│   └── common/            # Task management and shared types
├── examples/              # Working usage examples
│   ├── basic_agent/       # Agent usage example
│   ├── basic_workflow/    # Workflow orchestration example
│   ├── datetime_tools/    # Date/time manipulation tools
│   ├── feed_tools/        # RSS feed processing tools
│   └── web_tools/         # Web scraping and extraction tools
├── docs/                  # Documentation
├── test/                  # Integration tests
└── vendor/                # Vendored dependencies
```

## Interface Architecture

go-flock builds on go-llms with these key interfaces:

### Tools (`pkg/tools`)
- **FlockTool**: Extends `agentDomain.Tool` with categories, dependencies, validation
- **ToolRegistry**: Manages tool registration and discovery
- **ToolBuilder**: Fluent interface for tool construction

### Agents (`pkg/agents`)
- **FlockAgent**: Extends `agentDomain.Agent` with coordination capabilities
- **Coordinator**: Multi-agent coordination and task distribution
- **AgentRegistry**: Agent registration with role and capability filtering

### Workflows (`pkg/workflows`)
- **Workflow**: Multi-step process orchestration with conditions
- **WorkflowEngine**: Workflow registration and execution
- **WorkflowBuilder**: Fluent interface for workflow construction

### Common (`pkg/common`)
- **Task**: Structured task interface with requirements
- **TaskResult**: Standardized task execution results
- **CoordinationResult**: Multi-agent coordination outcomes

## Documentation

Comprehensive documentation is available in the `docs/` directory:

- **[Documentation Index](docs/README.md)** - Main documentation overview
- **[Tools Documentation](docs/tools/)** - Detailed documentation for all available tools
  - [Datetime Tools](docs/tools/datetime.md) - Date and time manipulation utilities
  - [Feed Tools](docs/tools/feed.md) - RSS/Atom feed fetching and processing
  - [Web Tools](docs/tools/web.md) - Web scraping, link extraction, and metadata tools
- **[Developer Guides](docs/developer/)** - Guides for extending go-flock
  - [Creating Custom Tools](docs/developer/creating-tools.md) - Step-by-step guide for building new tools

## Examples

The `examples/` directory contains working demonstrations:

```bash
# Build and run all examples
make run-examples

# Or run individually
./bin/basic_agent  
./bin/basic_workflow
./bin/datetime_tools
./bin/feed_tools
./bin/web_tools
```

Examples cover:
- **basic_agent/**: Agent usage with tool integration
- **basic_workflow/**: Multi-step workflow orchestration
- **datetime_tools/**: Complete datetime tools demonstration
- **feed_tools/**: RSS feed fetching with concurrent processing
- **web_tools/**: Web scraping, link extraction, and metadata extraction

## Contributing

1. Follow Go best practices and TDD approach
2. Use the Makefile for development tasks: `make dev`
3. Ensure all checks pass: `make check`
4. Write tests for new functionality
5. Update documentation for new features

## License

[Add your license information here]

## Related Projects

- [go-llms](https://github.com/lexlapax/go-llms) - The underlying LLM interface library providing the foundation for agent and tool interfaces