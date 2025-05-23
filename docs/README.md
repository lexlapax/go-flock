# go-flock Documentation

Welcome to the go-flock documentation. This directory contains comprehensive guides for using and extending go-flock.

## Documentation Structure

### [Tools Documentation](./tools/)
Detailed documentation for all available tools organized by category:
- [DateTime Tools](./tools/datetime.md) - Date and time operations
- [File Tools](./tools/file.md) - File system operations (coming soon)
- [HTTP Tools](./tools/http.md) - Network and API operations (coming soon)
- [Data Tools](./tools/data.md) - Data transformation and processing (coming soon)

### [Agents Documentation](./agents/)
Documentation for pre-configured agents:
- [Research Agents](./agents/research.md) - Web research and fact-checking (coming soon)
- [Code Agents](./agents/code.md) - Code analysis and generation (coming soon)
- [Data Agents](./agents/data.md) - Data analysis and reporting (coming soon)

### [Workflows Documentation](./workflows/)
Documentation for workflow patterns:
- [Research Workflows](./workflows/research.md) - Multi-source research patterns (coming soon)
- [Deployment Workflows](./workflows/deployment.md) - CI/CD and deployment patterns (coming soon)
- [Analysis Workflows](./workflows/analysis.md) - Data and code analysis patterns (coming soon)

### [Developer Guide](./developer/)
- [Creating Custom Tools](./developer/creating-tools.md)
- [Creating Custom Agents](./developer/creating-agents.md)
- [Creating Custom Workflows](./developer/creating-workflows.md)
- [Contributing Guidelines](./developer/contributing.md)

## Quick Links

- [Getting Started](../README.md#getting-started)
- [Examples](../examples/)
- [API Reference](#) (coming soon)

## Tool Naming Convention

go-flock follows a consistent `verb_object` naming pattern for all tools:

- `get_*` - Retrieve information
- `create_*` - Create new resources
- `update_*` - Modify existing resources
- `delete_*` - Remove resources
- `list_*` - List multiple items
- `parse_*` - Parse data from one format
- `format_*` - Format data to a specific format
- `convert_*` - Convert between formats
- `calculate_*` - Perform calculations
- `validate_*` - Validate data

## Integration with go-llms

All go-flock components are designed to work seamlessly with go-llms:

1. **Tools** implement the `domain.Tool` interface
2. **Agents** extend the `domain.Agent` interface
3. **Workflows** orchestrate multiple agents and tools

See the [examples](../examples/) directory for practical demonstrations.