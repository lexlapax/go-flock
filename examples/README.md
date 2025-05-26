# Examples

This directory contains example usage patterns for the go-flock library.

## Running Examples

Each example is a standalone Go program that demonstrates specific functionality:

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

### Feed Tools (RSS/Atom Feed Processing)
```bash
cd examples/feed_tools
go run main.go
```

### Web Tools (Web Scraping & Link Extraction)
```bash
cd examples/web_tools
go run main.go
```

### News API (News Search via NewsAPI.org)
```bash
cd examples/news_api
go run main.go
```

### Brave Search (Web Search with AI Summaries)
```bash
cd examples/brave_search
go run main.go
```

## Example Categories

- **basic_agent/**: Shows the pattern for specialized agents  
- **basic_workflow/**: Illustrates workflow orchestration patterns
- **datetime_tools/**: Working example of using datetime tools from go-flock
- **feed_tools/**: Demonstrates RSS feed fetching and processing capabilities
- **web_tools/**: Shows web page fetching and link extraction tools
- **news_api/**: Demonstrates news search using NewsAPI.org integration
- **brave_search/**: Demonstrates web search with AI summaries using Brave Search API

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