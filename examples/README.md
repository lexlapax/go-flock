# Examples

This directory contains example usage patterns for the go-flock library, organized by component type.

## Directory Structure

```
examples/
├── agents/      # Agent implementation examples
├── tools/       # Tool usage examples
└── workflows/   # Workflow orchestration examples
```

## Running Examples

Each example is a standalone Go program that demonstrates specific functionality.

### Agent Examples

Examples showing how to use and implement agents:

```bash
# Research Papers Agent (requires API keys)
cd examples/agents/research_papers
go run main.go -query "deep learning medical imaging"
```

### Tool Examples

Examples demonstrating individual tool usage:

```bash
# DateTime tools
cd examples/tools/datetime
go run main.go

# Feed processing tools
cd examples/tools/feed
go run main.go

# Web scraping tools
cd examples/tools/web
go run main.go

# News API search
cd examples/tools/news_api
go run main.go

# Brave web search
cd examples/tools/brave_search
go run main.go

# Academic research search
cd examples/tools/research_paper_api
go run main.go
```

### Workflow Examples

Examples showing multi-agent workflow orchestration:

```bash
# Basic workflow pattern
cd examples/workflows/basic
go run main.go
```

## Example Categories

### Agents (`agents/`)
- **research_papers/** - Complete production-ready Research Papers Agent with CLI interface showing:
  - How to build agent command-line tools
  - Integration with multiple tools
  - Configurable output formats
  - Provider selection and configuration
  - Best practices for agent implementation

### Tools (`tools/`)
- **datetime/** - Date/time manipulation tools (GetCurrentDateTime, CalculateDuration, etc.)
- **feed/** - RSS/Atom feed fetching and processing (FetchRSSFeed)
- **web/** - Web page fetching and analysis (FetchWebPage, ExtractLinks, ExtractMetadata)
- **news_api/** - News search using NewsAPI.org (NewsAPISearch)
- **brave_search/** - Web search with AI summaries using Brave Search API (BraveSearch)
- **research_paper_api/** - Multi-source academic research tool (ResearchPaperAPI)

### Workflows (`workflows/`)
- **basic/** - Simple workflow orchestration pattern

## Tool Naming Convention

go-flock follows a consistent `verb_object` naming pattern:
- `get_current_datetime` - Gets the current date/time
- `calculate_duration` - Calculates time between dates
- `fetch_rss_feed` - Fetches RSS/Atom feeds
- `extract_links` - Extracts links from web pages
- `research_paper_api` - Searches academic databases

## Environment Variables

Many examples require API keys. Set these before running:

```bash
# For LLM providers (when using real agents)
export OPENAI_API_KEY=your_key
export ANTHROPIC_API_KEY=your_key
export GEMINI_API_KEY=your_key

# For specific tools
export NEWS_API_KEY=your_newsapi_key
export BRAVE_API_KEY=your_brave_api_key
```

## Quick Start

1. **Explore Tools**: Start with tool examples to understand individual capabilities
2. **Learn Agents**: Review agent examples to see how tools are integrated
3. **Build Workflows**: Study workflow examples to orchestrate multiple agents

## Next Steps

As go-flock develops, more examples will be added:
- Production-ready agent implementations
- Complex multi-agent workflows
- Custom tool development patterns
- Integration with external systems
- Performance optimization techniques