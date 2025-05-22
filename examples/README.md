# Examples

This directory contains example usage patterns for the go-flock library.

## Running Examples

Each example is a standalone Go program that demonstrates specific functionality:

### Basic Tool Usage
```bash
cd examples/basic_tool
go run main.go
```

### Basic Agent Usage  
```bash
cd examples/basic_agent
go run main.go
```

### Basic Workflow Usage
```bash
cd examples/basic_workflow  
go run main.go
```

## Example Categories

- **basic_tool/**: Demonstrates creating and executing individual tools
- **basic_agent/**: Shows how to create agents that combine tools and logic
- **basic_workflow/**: Illustrates workflow orchestration of multiple agents

## Next Steps

These examples provide the foundation for understanding go-flock concepts. As the library develops, more sophisticated examples will be added showing:

- Real LLM integration using go-llms
- Complex tool implementations
- Multi-agent coordination
- Advanced workflow patterns
- Custom tool development