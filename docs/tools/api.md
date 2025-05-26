# API Tools Documentation

This document describes the API interaction tools available in go-flock. These tools enable communication with external REST APIs, GraphQL endpoints, and other web services.

## Available Tools

### search_news_api

Searches for news articles using NewsAPI.org service.

**Tool Name:** `search_news_api`

**Description:** Searches for news articles using a news API service

**Parameters:**
- `query` (string, required): Search query for news articles
- `api_key` (string, required*): API key for the news service (*can use NEWS_API_KEY environment variable)
- `language` (string, optional): Language code (e.g., 'en', 'es') - must be 2 characters
- `sort_by` (string, optional): Sort order - one of: 'relevancy', 'popularity', 'publishedAt' (default: 'publishedAt')
- `page_size` (integer, optional): Number of results per page, 1-100 (default: 20)
- `page` (integer, optional): Page number for pagination, minimum 1 (default: 1)
- `date_from` (string, optional): Oldest article date in YYYY-MM-DD format
- `date_to` (string, optional): Newest article date in YYYY-MM-DD format
- `domains` (string, optional): Comma-separated list of domains to include
- `exclude_domains` (string, optional): Comma-separated list of domains to exclude

**Returns:**
```json
{
  "status": "ok",
  "total_results": 100,
  "articles": [
    {
      "source": {
        "id": "techcrunch",
        "name": "TechCrunch"
      },
      "author": "John Doe",
      "title": "Article Title",
      "description": "Brief description",
      "url": "https://example.com/article",
      "published_at": "2024-01-15T10:00:00Z",
      "content": "Full article content..."
    }
  ],
  "fetched_at": "2024-01-15T12:00:00Z"
}
```

## Usage Examples

### Basic News Search

```go
tool := tools.NewSearchNewsAPITool()

params := tools.SearchNewsAPIParams{
    Query:    "artificial intelligence",
    APIKey:   "your-api-key", // or set NEWS_API_KEY env var
    Language: "en",
    PageSize: 10,
}

result, err := tool.Execute(ctx, params)
if err != nil {
    log.Fatal(err)
}

newsResult := result.(*tools.SearchNewsAPIResult)
fmt.Printf("Found %d articles\n", newsResult.TotalResults)

for _, article := range newsResult.Articles {
    fmt.Printf("- %s (%s)\n", article.Title, article.Source.Name)
}
```

### Search with Date Filters

```go
params := tools.SearchNewsAPIParams{
    Query:    "climate change",
    APIKey:   apiKey,
    DateFrom: "2024-01-01",
    DateTo:   "2024-01-31",
    SortBy:   "popularity",
}
```

### Domain-Specific Search

```go
params := tools.SearchNewsAPIParams{
    Query:   "technology",
    APIKey:  apiKey,
    Domains: "techcrunch.com,wired.com,arstechnica.com",
}
```

### Pagination

```go
// Get page 2 of results
params := tools.SearchNewsAPIParams{
    Query:    "sports",
    APIKey:   apiKey,
    PageSize: 50,
    Page:     2,
}
```

## API Key Configuration

The `search_news_api` tool supports two methods for API key configuration:

1. **Direct Parameter**: Pass the API key in the `api_key` parameter
2. **Environment Variable**: Set the `NEWS_API_KEY` environment variable

```bash
export NEWS_API_KEY=your_api_key_here
```

## NewsAPI.org Integration

This tool integrates with [NewsAPI.org](https://newsapi.org), a popular news aggregation API that provides access to headlines and articles from over 80,000 news sources worldwide.

### Getting an API Key

1. Visit https://newsapi.org
2. Sign up for a free account
3. Get your API key from the dashboard
4. Free tier includes 100 requests per day

### API Limitations

- Free tier: 100 requests/day, 30 days historical data
- Developer tier: 500 requests/day, 30 days historical data
- Business tier: Unlimited requests, full historical access

### Supported Languages

Common language codes:
- `en` - English
- `es` - Spanish
- `fr` - French
- `de` - German
- `it` - Italian
- `pt` - Portuguese
- `ru` - Russian
- `ar` - Arabic
- `zh` - Chinese
- `ja` - Japanese

### Sort Options

- `relevancy` - Articles more closely related to query come first
- `popularity` - Articles from popular sources and publishers come first
- `publishedAt` - Newest articles come first (default)

## Error Handling

The tool provides detailed error messages for common issues:

```go
result, err := tool.Execute(ctx, params)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "query parameter is required"):
        // Handle missing query
    case strings.Contains(err.Error(), "API key is required"):
        // Handle missing API key
    case strings.Contains(err.Error(), "apiKeyInvalid"):
        // Handle invalid API key
    case strings.Contains(err.Error(), "timeout"):
        // Handle timeout
    default:
        // Handle other errors
    }
}
```

## Best Practices

1. **Rate Limiting**: Respect API rate limits based on your subscription tier
2. **Caching**: Consider caching results to reduce API calls
3. **Error Handling**: Always check for errors and handle them appropriately
4. **Pagination**: Use pagination for large result sets
5. **Date Ranges**: Keep date ranges reasonable to avoid timeouts
6. **Domain Filtering**: Use domain filtering to get more relevant results

## Integration with Workflows

The news search tool can be combined with other tools in workflows:

```go
// Search for news
newsResults := searchNewsAPI(query)

// Extract content from articles
for _, article := range newsResults.Articles {
    content := fetchWebpage(article.URL)
    
    // Analyze with LLM agent
    summary := summarizerAgent.Summarize(content)
    sentiment := analyzerAgent.AnalyzeSentiment(content)
}
```

## Future Enhancements

Planned improvements for API tools:

1. **Additional News APIs**: Support for Bing News, Google News, etc.
2. **Generic REST Tool**: Configurable tool for any REST API
3. **GraphQL Support**: Tool for GraphQL queries
4. **OAuth Support**: Built-in OAuth flow handling
5. **Webhook Tools**: Tools for webhook integration
6. **Rate Limiting**: Built-in rate limiting and retry logic