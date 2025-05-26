# go-flock Documentation

Welcome to the go-flock documentation. This directory contains comprehensive guides for using and extending go-flock.

## Documentation Structure

### [Tools Documentation](./tools/)
Detailed documentation for all available tools organized by category:
- [DateTime Tools](./tools/datetime.md) - Date and time operations
- [Feed Tools](./tools/feed.md) - RSS/Atom feed fetching and parsing
- [Web Tools](./tools/web.md) - Web scraping, link extraction, and metadata tools
- [API Tools](./tools/api.md) - External API integration (NewsAPI, REST endpoints)
- [File Tools](./tools/file.md) - File system operations (coming soon)
- [HTTP Tools](./tools/http.md) - Network operations (coming soon)
- [Data Tools](./tools/data.md) - Data transformation and processing (coming soon)

### [Agents Documentation](./agents/)
Documentation for agent implementations:
- [Research Papers Agent](./agents/research_papers.md) - Academic paper search and analysis
- [CLI Examples](./agents/cli-examples.md) - Patterns for building agent CLIs
- [Implementation Status](./agents/implementation-status.md) - Current agent development status
- [Research Workflow Design](./agents/research-workflow-design.md) - Multi-agent research system design

### [Workflows Documentation](./workflows/)
Documentation for workflow patterns:
- [Research Workflow](./workflows/research-workflow.md) - Comprehensive research workflow design

### [Developer Guide](./developer/)
- [Creating Custom Tools](./developer/creating-tools.md)
- [Creating Custom Agents](./developer/creating-agents.md)
- [Creating Custom Workflows](./developer/creating-workflows.md)
- [Contributing Guidelines](./developer/contributing.md)

### [Troubleshooting](./troubleshooting.md)
- Common issues and solutions
- Debug logging guide
- Tool calling issues with different providers
- API key configuration

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