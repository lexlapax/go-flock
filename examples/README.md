# Examples

This directory contains example usage patterns for the go-flock library.

## Running Examples

Each example is a standalone Go program that demonstrates specific functionality:

### Basic Tool Pattern
```bash
cd examples/basic_tool
go run main.go
```

### Basic Agent Pattern
```bash
cd examples/basic_agent
go run main.go
```

### Basic Workflow Pattern
```bash
cd examples/basic_workflow  
go run main.go
```

### DateTime Tools (Working Example)
```bash
cd examples/datetime_tools
go run main.go
```

## Example Categories

- **basic_tool/**: Demonstrates the pattern for tools that go-flock provides
- **basic_agent/**: Shows the pattern for specialized agents  
- **basic_workflow/**: Illustrates workflow orchestration patterns
- **datetime_tools/**: Working example of using datetime tools from go-flock

## Tool Naming Convention

go-flock follows a consistent `verb_object` naming pattern:
- `get_current_datetime` - Gets the current date/time
- `calculate_duration` - Calculates time between dates
- `read_file` - Reads file contents
- `send_http_request` - Sends HTTP requests

Tools are grouped by functionality in files like:
- `datetime_tools.go` - All datetime-related tools
- `file_tools.go` - All file operation tools
- `http_tools.go` - All HTTP/network tools

## Next Steps

These examples provide the foundation for understanding go-flock concepts. As the library develops, more sophisticated examples will be added showing:

- Real LLM integration using go-llms providers
- Complex tool implementations
- Multi-agent coordination
- Advanced workflow patterns
- Custom tool development