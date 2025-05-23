# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-flock is a Go library providing tools, agents, and workflows for LLM-powered automation tasks. The library extends [go-llms](https://github.com/lexlapax/go-llms) with multi-agent coordination capabilities and workflow orchestration.

**Current Status: Foundation Complete** - Core interfaces, build system, and examples are implemented and working.

## Development Commands

The project uses a comprehensive Makefile for all development tasks:

### Primary Development Workflow
```bash
make dev                    # Complete development cycle: clean, fmt, vet, test, build
make check                  # All quality checks: fmt, vet, lint, test
make ci                     # Full CI pipeline with coverage and benchmarks
```

### Building
```bash
make build                  # Build flock CLI binary to bin/flock
make build-all              # Build all packages
make examples               # Build example programs
make run-examples           # Build and run all examples
```

### Testing
```bash
make test                   # Run unit tests
make test-coverage          # Generate HTML coverage report
make test-integration       # Run integration tests
make bench                  # Run benchmarks
make bench-cpu              # Run benchmarks with CPU profiling
make bench-mem              # Run benchmarks with memory profiling
```

### Code Quality
```bash
make fmt                    # Format Go code
make vet                    # Run go vet
make lint                   # Run golangci-lint
make security               # Run security checks with gosec
```

### Dependencies
```bash
make deps                   # Download and tidy dependencies
make deps-update            # Update dependencies to latest versions
make deps-vendor            # Vendor dependencies
```

### Tools and Utilities
```bash
make tools                  # Install development tools
make clean                  # Clean build artifacts
make info                   # Show build information
make serve-docs             # Serve documentation locally
```

## Architecture

### Core Components

1. **Tools (`pkg/tools/`)** - Individual LLM-callable functions with categories, dependencies, and validation
2. **Agents (`pkg/agents/`)** - Intelligent units with multi-agent coordination capabilities
3. **Workflows (`pkg/workflows/`)** - Complex multi-step process orchestration
4. **Common (`pkg/common/`)** - Shared task management and coordination types

### Interface Architecture

go-flock extends go-llms base interfaces:

- **FlockTool** extends `agentDomain.Tool` with categories, dependencies, async execution
- **FlockAgent** extends `agentDomain.Agent` with coordination, task execution, workflows
- **Workflow** provides step-based orchestration with conditions and failure strategies
- **Coordinator** manages multi-agent coordination and task distribution

### Key Design Patterns

- **Domain-Driven Design**: Clean separation between domain interfaces and implementations
- **Fluent Interfaces**: Builder patterns for tool, agent, and workflow construction
- **Registry Pattern**: Centralized management of tools, agents, and workflows
- **Strategy Pattern**: Pluggable coordination and failure handling strategies

## Dependencies

- **go-llms v0.2.6** - Core LLM interface library (dependency + submodule reference)
- Go 1.24.3+ required

## Source Reference

The `go-llms/` directory contains the complete source code of the go-llms library as a git submodule for development reference.

## Testing Strategy

Following TDD practices:
- Write tests before implementations
- Cover all interface implementations
- Integration tests for multi-agent scenarios
- Benchmarks for performance validation

## Current Implementation Status

✅ **Foundation Complete**
- Core interfaces designed and implemented
- Build system with comprehensive Makefile
- Working examples for all major components
- Proper go-llms integration and extension

🚧 **Next Steps**
- Implement concrete tool, agent, and workflow types
- Add comprehensive test coverage
- Build integration tests with real LLM providers
- Add advanced coordination strategies

## Notes for Development

- Use `make dev` for standard development workflow
- All code follows Go best practices with ABOUTME comments
- Interfaces are designed for extensibility and testing
- Multi-agent coordination is the key differentiator from base go-llms