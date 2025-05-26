# Tool Examples

This directory contains examples demonstrating individual tools available in go-flock.

## Overview

Tools are individual functions that can be called by LLMs to perform specific tasks. Each tool has a clear purpose, defined parameters, and structured output.

## Available Examples

### brave_search/
Web search using the Brave Search API, providing:
- Web search with AI-powered summaries
- News and image search capabilities
- Result ranking and relevance

### datetime/
Date and time manipulation tools:
- GetCurrentDateTime - Current date/time in any timezone
- CalculateDuration - Time between two dates
- ConvertTimezone - Timezone conversion
- FormatDateTime - Custom date formatting

### feed/
RSS/Atom feed processing:
- FetchRSSFeed - Fetch and parse feed content
- Concurrent feed processing
- Entry filtering and sorting

### news_api/
News search using NewsAPI.org:
- NewsAPISearch - Search news articles
- Category and source filtering
- Date range queries

### research_paper_api/
Academic research paper search:
- ResearchPaperAPI - Multi-database search (arXiv, PubMed, CORE)
- Citation extraction
- Research synthesis

### web/
Web scraping and analysis:
- FetchWebPage - Extract page content
- ExtractLinks - Find all links on a page
- ExtractMetadata - Get page metadata
- CheckURLStatus - Verify URL availability

## Running Examples

Each example can be run independently:

```bash
# Choose an example directory
cd datetime
go run main.go

# Some examples require API keys
export BRAVE_API_KEY=your_key
cd brave_search
go run main.go
```

## Environment Variables

Some tools require API keys:
- `BRAVE_API_KEY` - For Brave Search
- `NEWS_API_KEY` - For NewsAPI.org
- `OPENAI_API_KEY` / `ANTHROPIC_API_KEY` / `GEMINI_API_KEY` - For LLM providers

## Tool Patterns

All tools follow consistent patterns:
1. **Naming**: `verb_object` format (e.g., `fetch_webpage`, `research_paper_api`)
2. **Parameters**: Validated using go-llms schema
3. **Error Handling**: Clear error messages with context
4. **Output**: Structured data that can be easily processed

## Next Steps

- Implement your own tools following the patterns shown
- Combine tools in agents for complex tasks
- See `docs/developer/creating-tools.md` for tool development guide