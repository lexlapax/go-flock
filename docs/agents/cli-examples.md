# Agent CLI Examples

## Overview

Agent command-line interfaces demonstrate how to make agents accessible and usable. The examples show different patterns for creating CLI tools that leverage agent capabilities.

## Search Research Agent CLI

Location: `examples/agents/search_research/`

This is a complete, production-ready CLI for the Search Research Agent that demonstrates:

### Features
- Command-line argument parsing with flags
- Multiple output format support (Markdown, JSON, Text)
- Provider selection (OpenAI, Anthropic, Gemini)
- File output options
- Comprehensive help documentation
- Environment variable configuration

### Usage Patterns

```bash
# Basic usage
go run main.go -query "topic"

# Advanced options
go run main.go \
  -query "deep learning medical imaging" \
  -format json \
  -output results.json \
  -provider openai \
  -model gpt-4
```

### Key Implementation Details

1. **Provider Factory Pattern**
```go
func createProvider(providerName string) (ldomain.Provider, error) {
    // Auto-detect from environment if not specified
    // Create appropriate provider based on selection
}
```

2. **Output Format Configuration**
```go
agentOpts := agents.AgentOptions{
    OutputFormat: outputFormat,
    Model:        *model,
}
```

3. **Error Handling**
- Validation of required parameters
- Clear error messages
- Proper exit codes

## Building Your Own Agent CLI

When creating CLI interfaces for agents, consider:

1. **User Experience**
   - Clear help messages
   - Sensible defaults
   - Progress indicators
   - Structured output options

2. **Configuration**
   - Environment variables for secrets
   - Command-line flags for options
   - Configuration files for complex setups

3. **Integration**
   - Pipe-friendly output
   - JSON format for automation
   - File output for persistence

4. **Error Handling**
   - Validate inputs early
   - Provide helpful error messages
   - Use appropriate exit codes

## Future CLI Examples

Planned agent CLI examples:
- News Analysis Agent CLI
- Code Review Agent CLI
- Data Analysis Agent CLI
- Multi-Agent Coordinator CLI