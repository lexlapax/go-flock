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

### search_web_brave

Performs comprehensive web search using Brave Search API with support for multiple content types, AI summaries, and advanced filtering.

**Tool Name:** `search_web_brave`

**Description:** Performs web search using Brave Search API with support for multiple content types, AI summaries, and advanced filtering

**Parameters:**
- `query` (string, required): Search query (max 400 chars/50 words)
- `api_key` (string, required*): API key for Brave Search (*can use BRAVE_SEARCH_API_KEY environment variable)
- `count` (integer, optional): Results per page, 1-20 (default: 10)
- `offset` (integer, optional): Page offset for pagination, 0-9 (default: 0)
- `country` (string, optional): Country code (e.g., 'US', 'GB') - must be 2 uppercase letters
- `search_lang` (string, optional): Language preference (e.g., 'en', 'es') - must be 2 lowercase letters
- `safesearch` (string, optional): Content filter - one of: 'off', 'moderate', 'strict'
- `freshness` (string, optional): Time filter - one of: 'pd' (past day), 'pw' (past week), 'pm' (past month), 'py' (past year)
- `result_filter` (array, optional): Filter types - array of: 'web', 'news', 'images', 'videos'
- `extra_snippets` (boolean, optional): Enable additional text excerpts
- `summary` (boolean, optional): Enable AI-generated summaries
- `goggles` (string, optional): Custom ranking profile URL

**Returns:**
```json
{
  "query": "artificial intelligence",
  "type": "search",
  "mixed": {
    "type": "mixed",
    "main": [
      {
        "type": "web",
        "index": 0
      }
    ]
  },
  "web": {
    "type": "search",
    "results": [
      {
        "title": "What is AI?",
        "url": "https://example.com/ai",
        "description": "Introduction to artificial intelligence",
        "age": "2 days ago",
        "language": "en",
        "thumbnail": {
          "src": "https://example.com/thumb.jpg",
          "height": 200,
          "width": 300
        }
      }
    ]
  },
  "news": {
    "type": "news",
    "results": [
      {
        "title": "AI Breakthrough",
        "url": "https://news.example.com/ai",
        "description": "Latest AI research",
        "age": "3 hours ago",
        "source": {
          "name": "Tech News",
          "url": "https://technews.com"
        }
      }
    ]
  },
  "summary": [
    {
      "type": "summary",
      "key": "brave_search_llm_summary",
      "text": "AI refers to computer systems that can perform tasks requiring human intelligence..."
    }
  ],
  "infobox": {
    "type": "infobox",
    "title": "Artificial Intelligence",
    "description": "Branch of computer science",
    "url": "https://en.wikipedia.org/wiki/AI"
  },
  "fetched_at": "2024-01-15T12:00:00Z"
}
```

### Brave Search Usage Examples

#### Basic Web Search

```go
tool := tools.NewSearchWebBraveTool()

params := tools.SearchWebBraveParams{
    Query:  "golang best practices",
    APIKey: "your-api-key", // or set BRAVE_SEARCH_API_KEY env var
    Count:  10,
}

result, err := tool.Execute(ctx, params)
braveResult := result.(*tools.SearchWebBraveResult)

for _, webResult := range braveResult.Web.Results {
    fmt.Printf("- %s (%s)\n", webResult.Title, webResult.URL)
}
```

#### Search with AI Summary

```go
params := tools.SearchWebBraveParams{
    Query:   "explain quantum computing",
    APIKey:  apiKey,
    Summary: true,
}

result, err := tool.Execute(ctx, params)
braveResult := result.(*tools.SearchWebBraveResult)

if len(braveResult.Summary) > 0 {
    fmt.Println("AI Summary:", braveResult.Summary[0].Text)
}
```

#### News-Focused Search

```go
params := tools.SearchWebBraveParams{
    Query:        "climate change",
    APIKey:       apiKey,
    Freshness:    "pd", // Past day
    ResultFilter: []string{"news"},
    Country:      "US",
}
```

#### Multi-Type Search

```go
params := tools.SearchWebBraveParams{
    Query:        "SpaceX launch",
    APIKey:       apiKey,
    ResultFilter: []string{"web", "news", "videos"},
}

result, err := tool.Execute(ctx, params)
braveResult := result.(*tools.SearchWebBraveResult)

// Access different content types
webResults := braveResult.Web.Results
newsResults := braveResult.News.Results
videoResults := braveResult.Videos.Results
```

#### Safe Search with Language

```go
params := tools.SearchWebBraveParams{
    Query:      "educational content",
    APIKey:     apiKey,
    SafeSearch: "strict",
    SearchLang: "en",
    Country:    "US",
}
```

### Brave Search API Integration

This tool integrates with [Brave Search API](https://api.search.brave.com), a privacy-focused search API that provides:

#### Getting an API Key

1. Visit https://api.search.brave.com
2. Sign up for an account
3. Get your API key from the dashboard
4. Free tier includes 2,000 queries per month

#### API Limitations

- Free tier: 2,000 queries/month, 1 query/second
- Basic tier: 20,000 queries/month, 5 queries/second
- Pro tier: 250,000+ queries/month, custom rate limits
- Max 20 results per request
- Max offset of 9 (pages 0-9)

#### Supported Countries

Common country codes:
- `US` - United States
- `GB` - United Kingdom
- `CA` - Canada
- `AU` - Australia
- `DE` - Germany
- `FR` - France
- `JP` - Japan
- `IN` - India

#### Freshness Options

- `pd` - Past day (24 hours)
- `pw` - Past week (7 days)
- `pm` - Past month (30 days)
- `py` - Past year (365 days)

#### Content Types

- `web` - Standard web pages
- `news` - News articles
- `images` - Image results
- `videos` - Video results

### Brave vs NewsAPI Comparison

| Feature | Brave Search | NewsAPI |
|---------|--------------|---------|
| Content Type | Web, News, Images, Videos | News only |
| AI Summaries | ✅ Yes | ❌ No |
| Privacy Focus | ✅ Strong | ➖ Standard |
| Free Tier | 2,000/month | 100/day |
| Search Scope | Entire web | News sources |
| InfoBox Data | ✅ Yes | ❌ No |
| FAQ Extraction | ✅ Yes | ❌ No |
| Location Results | ✅ Yes | ❌ No |

### Best Practices for Brave Search

1. **Query Optimization**: Keep queries under 50 words for best results
2. **AI Summaries**: Use for complex topics or quick overviews
3. **Result Filtering**: Use `result_filter` to reduce API calls
4. **Pagination**: Limited to 10 pages (offset 0-9)
5. **Language Matching**: Set both `country` and `search_lang` for best results
6. **Rate Limiting**: Implement client-side rate limiting

### Integration with News Search Workflow

Brave Search complements NewsAPI for comprehensive news gathering:

```go
// Use Brave for general web + news search
braveResults := searchWebBrave("AI breakthroughs", []string{"web", "news"})

// Use NewsAPI for deep news coverage
newsAPIResults := searchNewsAPI("AI breakthroughs", dateFrom: "2024-01-01")

// Combine and deduplicate results
allResults := combineResults(braveResults, newsAPIResults)
uniqueResults := deduplicateByURL(allResults)

// Use AI summary from Brave for overview
summary := braveResults.Summary[0].Text
```

### research_paper_api

Searches for academic papers across multiple research databases (arXiv, PubMed, CORE) in parallel.

**Tool Name:** `research_paper_api`

**Description:** Searches for academic papers across multiple research databases (arXiv, PubMed, CORE) in parallel

**Parameters:**
- `query` (string, required): Search query for research papers
- `max_results` (integer, optional): Maximum results per provider, 1-100 (default: 10)
- `start_date` (string, optional): Filter papers from this date (YYYY-MM-DD format)
- `end_date` (string, optional): Filter papers until this date (YYYY-MM-DD format)
- `authors` (array, optional): Filter by author names
- `categories` (array, optional): Subject categories (cs, physics, medicine, etc.)
- `open_access` (boolean, optional): Only return open access papers
- `sort_by` (string, optional): Sort order - one of: 'relevance', 'date', 'citations' (default: 'relevance')
- `providers` (array, optional): Specific providers to search - array of: 'arxiv', 'pubmed', 'core' (default: all)
- `core_api_key` (string, optional): CORE API key for enhanced access

**Returns:**
```json
{
  "query": "machine learning healthcare",
  "total_results": 45,
  "papers": [
    {
      "title": "Deep Learning for Medical Image Analysis",
      "authors": ["John Doe", "Jane Smith"],
      "abstract": "This paper presents a comprehensive survey...",
      "published_date": "2024-01-15",
      "source": "arXiv",
      "url": "http://arxiv.org/abs/2401.12345",
      "pdf_url": "http://arxiv.org/pdf/2401.12345",
      "doi": "10.1234/example.doi",
      "arxiv_id": "2401.12345",
      "pubmed_id": "",
      "journal": "",
      "relevance_score": 0.95
    }
  ],
  "providers": [
    {
      "name": "arxiv",
      "result_count": 20,
      "response_time_ms": 523
    },
    {
      "name": "pubmed",
      "result_count": 15,
      "response_time_ms": 867
    },
    {
      "name": "core",
      "result_count": 10,
      "response_time_ms": 445
    }
  ],
  "fetched_at": "2024-01-15T12:00:00Z"
}
```

### Research Search Usage Examples

#### Basic Research Search

```go
tool := tools.NewResearchPaperAPITool()

params := tools.ResearchPaperAPIParams{
    Query:      "machine learning healthcare",
    MaxResults: 20,
}

result, err := tool.Execute(ctx, params)
researchResult := result.(*tools.ResearchPaperAPIResult)

fmt.Printf("Found %d unique papers\n", researchResult.TotalResults)

for _, paper := range researchResult.Papers {
    fmt.Printf("- [%s] %s by %s\n", 
        paper.Source, 
        paper.Title, 
        strings.Join(paper.Authors, ", "))
}
```

#### Category-Specific Search

```go
params := tools.ResearchPaperAPIParams{
    Query:      "neural networks",
    Categories: []string{"cs", "math"},
    MaxResults: 30,
    SortBy:     "date",
}
```

#### Medical Research with Date Filter

```go
params := tools.ResearchPaperAPIParams{
    Query:      "COVID-19 vaccines",
    Providers:  []string{"pubmed"},
    StartDate:  "2023-01-01",
    EndDate:    "2024-01-01",
    OpenAccess: true,
}
```

#### Author-Specific Search

```go
params := tools.ResearchPaperAPIParams{
    Query:   "deep learning",
    Authors: []string{"Geoffrey Hinton", "Yann LeCun"},
}
```

#### CORE-Only Search with API Key

```go
params := tools.ResearchPaperAPIParams{
    Query:       "climate change mitigation",
    Providers:   []string{"core"},
    CoreAPIKey:  "your-core-api-key", // or set CORE_API_KEY env var
    MaxResults:  50,
}
```

### Provider Details

#### arXiv
- **Coverage**: Physics, Mathematics, Computer Science, Quantitative Biology, Quantitative Finance, Statistics, Electrical Engineering, Economics
- **Access**: Free, no API key required
- **Rate Limit**: 3 requests per second
- **Special Features**: Preprints, version history, PDF access
- **Categories**: cs (Computer Science), math (Mathematics), physics, q-bio (Quantitative Biology), etc.

#### PubMed
- **Coverage**: Biomedical and life sciences literature
- **Access**: Free, optional API key for higher rate limits
- **Rate Limit**: 3/second without key, 10/second with key
- **Special Features**: Peer-reviewed content, clinical trials, MeSH terms
- **Environment Variable**: `PUBMED_API_KEY`

#### CORE
- **Coverage**: Open access research from repositories worldwide
- **Access**: Requires free API key from https://core.ac.uk
- **Rate Limit**: 10 requests per second
- **Special Features**: Full-text search, diverse content, global coverage
- **Environment Variable**: `CORE_API_KEY`

### Environment Variables

```bash
# Optional: For higher PubMed rate limits
export PUBMED_API_KEY=your_pubmed_api_key

# Required: For CORE access
export CORE_API_KEY=your_core_api_key
```

### Advanced Features

#### Parallel Search
The tool searches all selected providers in parallel for optimal performance:

```go
// This searches arXiv, PubMed, and CORE simultaneously
params := tools.ResearchPaperAPIParams{
    Query:      "artificial intelligence",
    MaxResults: 100, // 100 per provider
}
```

#### Deduplication
Results are automatically deduplicated based on:
- DOI (Digital Object Identifier)
- Title similarity

#### Provider Selection
The tool intelligently selects providers based on categories:

```go
// Automatically uses arXiv for CS categories
params := tools.ResearchPaperAPIParams{
    Query:      "algorithms",
    Categories: []string{"cs"},
}

// Automatically uses PubMed for medical categories
params := tools.ResearchPaperAPIParams{
    Query:      "cancer treatment",
    Categories: []string{"medicine"},
}
```

#### Error Handling
Provider failures are handled gracefully:

```go
result, err := tool.Execute(ctx, params)
if err != nil {
    // Total failure - no providers returned results
    log.Fatal(err)
}

// Check individual provider status
for _, provider := range result.Providers {
    if provider.Error != "" {
        fmt.Printf("Provider %s failed: %s\n", provider.Name, provider.Error)
    }
}
```

### Research Tool vs Other Search Tools

| Feature | Research Search | News API | Brave Search |
|---------|----------------|----------|--------------|
| Content Type | Academic papers | News articles | Web pages |
| Peer Review | ✅ Yes (PubMed) | ❌ No | ❌ No |
| PDF Access | ✅ Yes | ❌ No | ➖ Limited |
| Citation Info | ✅ Yes | ❌ No | ❌ No |
| DOI Support | ✅ Yes | ❌ No | ❌ No |
| Preprints | ✅ Yes (arXiv) | ❌ No | ➖ Maybe |
| Medical Focus | ✅ Yes (PubMed) | ➖ General | ➖ General |

### Integration with Research Workflow

```go
// Search for papers
researchResults := researchPaperAPI("machine learning healthcare")

// Download PDFs for open access papers
for _, paper := range researchResults.Papers {
    if paper.PDFURL != "" {
        pdf := downloadPDF(paper.PDFURL)
        
        // Extract text for analysis
        text := extractTextFromPDF(pdf)
        
        // Analyze with LLM agent
        summary := summarizerAgent.Summarize(text)
        keyFindings := researchAgent.ExtractFindings(text)
    }
}

// Cross-reference with news
newsResults := searchNewsAPI(paper.Title)
mediaVoverage := analyzeMediaCoverage(paper, newsResults)
```

### Best Practices

1. **API Keys**: Set up environment variables for API keys
2. **Rate Limiting**: The tool implements automatic rate limiting
3. **Category Selection**: Use categories to optimize provider selection
4. **Date Ranges**: Use reasonable date ranges for faster results
5. **Deduplication**: Results are automatically deduplicated
6. **Error Handling**: Check provider-specific errors in results
7. **Parallel Queries**: Default behavior for best performance

## Future Enhancements

Planned improvements for API tools:

1. **Additional Search APIs**: Support for Bing Search, Google Custom Search, DuckDuckGo
2. **Generic REST Tool**: Configurable tool for any REST API
3. **GraphQL Support**: Tool for GraphQL queries
4. **OAuth Support**: Built-in OAuth flow handling
5. **Webhook Tools**: Tools for webhook integration
6. **Rate Limiting**: Built-in rate limiting and retry logic
7. **Result Caching**: Cache search results to reduce API calls
8. **Brave Local POI Search**: Integration with Brave's local search endpoints