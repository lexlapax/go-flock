# Gather News Agent

## Overview

The Gather News Agent specializes in collecting current news and events related to specific topics. It searches multiple news sources simultaneously, analyzes coverage patterns, and presents findings in various formats. This agent is designed to provide comprehensive news coverage by using both dedicated news APIs and general web search with news filtering.

> **Developer Note**: This agent follows the same architectural patterns as the Research Papers Agent. See the [Creating Custom Agents](../developer/creating-agents.md) guide for implementation details.

## Features

- **Multi-source Search**: Queries NewsAPI and Brave Search simultaneously
- **Flexible Output**: Supports Markdown (default), JSON, and plain text formats
- **Comprehensive Analysis**: Extracts articles, trends, perspectives, and timelines
- **Tool Integration**: Uses search_news_api, search_web_brave, fetch_webpage, and extract_metadata tools
- **Configurable Models**: Works with any go-llms supported LLM provider
- **Debug Support**: Comprehensive logging for troubleshooting

## Usage

### As a Library

```go
import (
    "context"
    "github.com/lexlapax/go-flock/pkg/agents"
    "github.com/lexlapax/go-llms/pkg/llm/provider"
)

// Create provider
provider := provider.NewOpenAIProvider(apiKey, "gpt-4")

// Create agent with default markdown output
agent := agents.NewGatherNewsAgent(provider)

// Or with specific output format
agent := agents.NewGatherNewsAgent(provider, agents.AgentOptions{
    OutputFormat: agents.OutputFormatJSON,
    Model: "gpt-4",
})

// Execute search
ctx := context.Background()
result, err := agent.Run(ctx, "artificial intelligence breakthroughs")
```

### Command Line

The gather news agent includes a complete CLI example in `examples/agents/gather_news/`:

```bash
# Navigate to the example
cd examples/agents/gather_news

# Basic search with markdown output
go run main.go -query "artificial intelligence breakthroughs"

# JSON output saved to file
go run main.go -query "climate change policy" -format json -output news.json

# Using specific provider and model
go run main.go -query "tech industry layoffs" -provider openai -model gpt-4

# Plain text output
go run main.go -query "renewable energy developments" -format text

# With debug logging enabled
go run main.go -query "space exploration" -debug

# Or using environment variable
FLOCK_DEBUG=true go run main.go -query "quantum computing news"
```

## Output Formats

### Markdown (Default)

Well-formatted report with:
- Executive summary
- Major stories with details
- Emerging trends
- Different perspectives
- Timeline of events
- Source citations

Example:
```markdown
# News Analysis: Artificial Intelligence Breakthroughs

## Executive Summary
Recent weeks have seen significant AI developments...

## Major Stories

### OpenAI Announces New Model Architecture
- **Source**: TechCrunch
- **Date**: 2024-01-20
- **Summary**: OpenAI reveals breakthrough in efficiency...
- **Key Points**:
  - 50% reduction in computational requirements
  - Improved reasoning capabilities
- **URL**: https://techcrunch.com/...
```

### JSON

Structured data format:
```json
{
  "topic": "artificial intelligence breakthroughs",
  "summary": "Recent AI developments show...",
  "articles": [
    {
      "title": "OpenAI Announces New Model",
      "source": "TechCrunch",
      "author": "Jane Smith",
      "published_date": "2024-01-20",
      "summary": "OpenAI reveals breakthrough...",
      "key_points": [
        "50% efficiency improvement",
        "Enhanced reasoning"
      ],
      "url": "https://techcrunch.com/...",
      "sentiment": "positive"
    }
  ],
  "trends": ["efficiency improvements", "reasoning capabilities"],
  "perspectives": [
    {
      "viewpoint": "Industry optimism",
      "sources": ["TechCrunch", "Wired"],
      "summary": "Tech companies see potential..."
    }
  ],
  "timeline": [
    {
      "date": "2024-01-20",
      "event": "OpenAI announcement"
    }
  ],
  "analysis_date": "2024-01-21"
}
```

### Plain Text

Simple, readable format without markup:
```
NEWS ANALYSIS: ARTIFICIAL INTELLIGENCE BREAKTHROUGHS

EXECUTIVE SUMMARY
Recent weeks have seen significant AI developments...

MAJOR STORIES
1. OpenAI Announces New Model Architecture
   Source: TechCrunch (2024-01-20)
   Summary: OpenAI reveals breakthrough in efficiency...
   Key Points:
   - 50% reduction in computational requirements
   - Improved reasoning capabilities
   Link: https://techcrunch.com/...

EMERGING TRENDS
- Efficiency improvements across models
- Focus on reasoning capabilities

SOURCES USED
- TechCrunch
- Wired
- Reuters
```

## Configuration

### Environment Variables

- `OPENAI_API_KEY` - For OpenAI provider
- `ANTHROPIC_API_KEY` - For Anthropic provider  
- `GEMINI_API_KEY` - For Google Gemini provider
- `NEWS_API_KEY` - For NewsAPI.org integration
- `BRAVE_SEARCH_API_KEY` - For Brave Search integration
- `FLOCK_DEBUG` - Enable debug logging (set to 'true' or '1')

### Agent Options

```go
type AgentOptions struct {
    OutputFormat OutputFormat // markdown, json, or text
    Model        string      // LLM model name (optional)
}
```

## Tools Used

1. **search_news_api** - Primary tool for dedicated news search via NewsAPI
2. **search_web_brave** - Web search with news filtering capabilities
3. **fetch_webpage** - Retrieves full article content when needed
4. **extract_metadata** - Extracts publication metadata from web pages

## Prompt Architecture

The Gather News Agent uses the same modular prompt structure as other agents:

### Core Prompt
- Agent role as news research specialist
- Available tools and their purposes
- Critical instructions for tool usage
- Analysis guidelines

### Format-Specific Instructions
- **Markdown**: Section structure, formatting guidelines
- **JSON**: Schema definition, field requirements
- **Text**: Plain text formatting rules

This separation allows easy maintenance and extension of the agent's capabilities.

## Integration Examples

### News Monitoring Pipeline

```go
// Create news gathering agent
newsAgent := agents.NewGatherNewsAgent(provider)

// Get current news
newsData, _ := newsAgent.Run(ctx, "renewable energy policy")

// Pass to fact verification agent
verifyAgent := agents.NewVerifyFactsAgent(provider)
verified, _ := verifyAgent.Run(ctx, newsData)

// Generate summary
summaryAgent := agents.NewCreateSummaryAgent(provider)
summary, _ := summaryAgent.Run(ctx, verified)
```

### Scheduled News Updates

```go
ticker := time.NewTicker(6 * time.Hour)
defer ticker.Stop()

topics := []string{
    "artificial intelligence",
    "climate change",
    "space exploration",
}

for {
    select {
    case <-ticker.C:
        for _, topic := range topics {
            result, err := agent.Run(ctx, topic)
            if err != nil {
                log.Printf("Error gathering news for %s: %v", topic, err)
                continue
            }
            saveNewsUpdate(topic, result)
            notifySubscribers(topic, result)
        }
    }
}
```

### News API Service

```go
http.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    format := r.URL.Query().Get("format")
    
    opts := agents.AgentOptions{
        OutputFormat: agents.OutputFormatJSON,
    }
    if format != "" {
        opts.OutputFormat = agents.OutputFormat(format)
    }
    
    agent := agents.NewGatherNewsAgent(provider, opts)
    result, err := agent.Run(context.Background(), query)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(result.(string)))
})
```

## Best Practices

1. **API Keys**: Set both NEWS_API_KEY and BRAVE_SEARCH_API_KEY for comprehensive coverage
2. **Specific Queries**: More specific queries yield more relevant results
3. **Date Ranges**: Consider recent news vs. historical coverage needs
4. **Output Format**: Choose format based on downstream processing needs
5. **Rate Limiting**: Be aware of API rate limits for both news services
6. **Caching**: Consider caching results for frequently requested topics
7. **Error Handling**: Always check for API errors and rate limit issues

## Limitations

- Requires active internet connection
- Subject to API rate limits (NewsAPI: 100-500/day, Brave: 2000/month)
- News coverage depends on source availability
- Real-time news may have slight delays
- Some sources may require additional authentication

## Debugging

Enable debug logging to troubleshoot issues:

```bash
# Via CLI flag
go run main.go -query "your topic" -debug

# Via environment variable
FLOCK_DEBUG=true go run main.go -query "your topic"
```

Debug logging shows:
- Tool execution details
- API calls and responses
- LLM communication
- Error diagnostics

## Troubleshooting

### No Results Found
- Verify NEWS_API_KEY and/or BRAVE_SEARCH_API_KEY are set
- Check internet connectivity
- Try broader search terms
- Enable debug logging to see API responses

### API Key Issues
- Ensure API keys are valid and not expired
- Check rate limits haven't been exceeded
- Verify environment variables are properly set
- Consider using both NewsAPI and Brave for redundancy

### Tool Call Issues (Gemini)
- The agent includes Gemini-specific formatting instructions
- Tool arguments are properly formatted as JSON strings
- See [Troubleshooting Guide](../troubleshooting.md) for details

### Output Format Errors
- Ensure consistent format selection
- Check LLM model compatibility
- Review debug logs for parsing issues

## Comparison with Research Papers Agent

| Feature | Gather News | Research Papers |
|---------|-------------|-----------------|
| Content Type | Current news articles | Academic papers |
| Sources | NewsAPI, Brave Search | arXiv, PubMed, CORE |
| Timeliness | Real-time to recent | Historical to recent |
| Content Depth | News summaries | Full abstracts |
| Update Frequency | Hourly/Daily | Weekly/Monthly |
| Use Case | Current events tracking | Academic research |

## Future Enhancements

- Support for additional news APIs (Google News, Bing News)
- Sentiment analysis integration
- Real-time news monitoring with webhooks
- News categorization and tagging
- Multi-language support
- RSS feed generation