# Web Tools

The web tools package provides utilities for web scraping, content extraction, and URL analysis.

## Available Tools

### fetch_webpage

Fetches content from web pages with options for text extraction and custom headers.

**Function**: `NewFetchWebPageTool()`

**Parameters**:
- `url` (string, required): The URL of the web page to fetch
- `extract_text` (boolean, optional): Extract only text content, removing HTML tags (default: true)
- `timeout` (integer, optional): Timeout in seconds (default: 30, max: 300)
- `user_agent` (string, optional): Custom user agent string
- `headers` (object, optional): Additional HTTP headers as key-value pairs
- `follow_redirects` (boolean, optional): Follow HTTP redirects (default: true)

**Returns**:
```json
{
  "url": "https://example.com",
  "title": "Example Domain",
  "content": "Example Domain\nThis domain is for use in illustrative examples...",
  "content_type": "text/html; charset=UTF-8",
  "status_code": 200,
  "headers": {
    "Content-Type": "text/html; charset=UTF-8",
    "Server": "nginx/1.18.0"
  },
  "fetched_at": "2024-12-15T10:30:00Z"
}
```

**Example Usage**:
```go
tool := tools.NewFetchWebPageTool()
params := tools.FetchWebPageParams{
    URL:         "https://example.com",
    ExtractText: true,
    UserAgent:   "MyBot/1.0",
    Headers: map[string]string{
        "Accept-Language": "en-US",
    },
}
result, err := tool.Execute(ctx, params)
```

### extract_links

Extracts all links from a web page, categorized by type.

**Function**: `NewExtractLinksTool()`

**Parameters**:
- `url` (string, required): The URL of the web page to extract links from
- `include_external` (boolean, optional): Include external links (default: true)
- `include_internal` (boolean, optional): Include internal links (default: true)
- `include_email` (boolean, optional): Include email links (default: true)
- `include_media` (boolean, optional): Include media links like images (default: false)
- `timeout` (integer, optional): Timeout in seconds (default: 30)

**Returns**:
```json
{
  "url": "https://example.com",
  "total_links": 15,
  "internal_links": [
    {
      "url": "https://example.com/about",
      "text": "About Us",
      "title": "Learn more about our company"
    }
  ],
  "external_links": [
    {
      "url": "https://github.com/example",
      "text": "Our GitHub",
      "title": ""
    }
  ],
  "email_links": ["contact@example.com"],
  "media_links": [
    {
      "url": "https://example.com/logo.png",
      "type": "image",
      "alt": "Company Logo"
    }
  ],
  "fetched_at": "2024-12-15T10:30:00Z"
}
```

**Example Usage**:
```go
tool := tools.NewExtractLinksTool()
params := tools.ExtractLinksParams{
    URL:             "https://news.ycombinator.com",
    IncludeExternal: true,
    IncludeInternal: true,
    IncludeEmail:    false,
    IncludeMedia:    true,
}
result, err := tool.Execute(ctx, params)
```

### extract_metadata

Extracts metadata from web pages including Open Graph and Twitter Card data.

**Function**: `NewExtractMetadataTool()`

**Parameters**:
- `url` (string, required): The URL of the web page to extract metadata from
- `timeout` (integer, optional): Timeout in seconds (default: 30)

**Returns**:
```json
{
  "url": "https://example.com/article",
  "title": "Article Title",
  "description": "A brief description of the article",
  "author": "John Doe",
  "keywords": ["technology", "news", "innovation"],
  "open_graph": {
    "title": "Article Title - Example.com",
    "description": "A brief description for social sharing",
    "type": "article",
    "image": "https://example.com/article-image.jpg",
    "url": "https://example.com/article"
  },
  "twitter": {
    "card": "summary_large_image",
    "site": "@example",
    "creator": "@johndoe",
    "title": "Article Title",
    "image": "https://example.com/twitter-image.jpg"
  },
  "meta": {
    "viewport": "width=device-width, initial-scale=1.0",
    "robots": "index, follow",
    "theme-color": "#007bff"
  },
  "fetched_at": "2024-12-15T10:30:00Z"
}
```

**Example Usage**:
```go
tool := tools.NewExtractMetadataTool()
params := tools.ExtractMetadataParams{
    URL: "https://www.bbc.com/news/article-123",
}
result, err := tool.Execute(ctx, params)
```

### check_url_status

Checks if a URL is accessible and returns status information.

**Function**: `NewCheckURLStatusTool()`

**Parameters**:
- `url` (string, required): The URL to check
- `timeout` (integer, optional): Timeout in seconds (default: 10, max: 300)
- `follow_redirects` (boolean, optional): Follow redirects (default: true)
- `headers` (object, optional): Additional HTTP headers

**Returns**:
```json
{
  "url": "https://example.com",
  "status": "alive",  // "alive", "dead", or "redirect"
  "status_code": 200,
  "response_time_ms": 125,
  "content_type": "text/html; charset=UTF-8",
  "content_length": 1256,
  "final_url": "",  // Only if redirected
  "redirect_chain": [],  // List of redirect URLs
  "headers": {
    "Server": "nginx/1.18.0",
    "X-Frame-Options": "DENY"
  },
  "checked_at": "2024-12-15T10:30:00Z"
}
```

**Example Usage**:
```go
tool := tools.NewCheckURLStatusTool()
params := tools.CheckURLStatusParams{
    URL:             "https://github.com",
    FollowRedirects: true,
    Timeout:         5,
    Headers: map[string]string{
        "Accept": "text/html",
    },
}
result, err := tool.Execute(ctx, params)
```

## Use Cases

### News Aggregation
```go
// 1. Check if news source is available
statusTool := tools.NewCheckURLStatusTool()
status, _ := statusTool.Execute(ctx, tools.CheckURLStatusParams{
    URL: "https://techcrunch.com",
})

// 2. Fetch the page
fetchTool := tools.NewFetchWebPageTool()
page, _ := fetchTool.Execute(ctx, tools.FetchWebPageParams{
    URL: "https://techcrunch.com",
    ExtractText: false,  // Keep HTML for link extraction
})

// 3. Extract all article links
linkTool := tools.NewExtractLinksTool()
links, _ := linkTool.Execute(ctx, tools.ExtractLinksParams{
    URL: "https://techcrunch.com",
    IncludeInternal: true,
    IncludeExternal: false,
})

// 4. Get metadata for sharing
metaTool := tools.NewExtractMetadataTool()
meta, _ := metaTool.Execute(ctx, tools.ExtractMetadataParams{
    URL: articleURL,
})
```

### SEO Analysis
```go
// Extract metadata for SEO audit
metaTool := tools.NewExtractMetadataTool()
meta, _ := metaTool.Execute(ctx, tools.ExtractMetadataParams{
    URL: "https://example.com",
})

// Check for required meta tags
if meta.Description == "" {
    log.Println("Missing meta description")
}
if len(meta.OpenGraph) == 0 {
    log.Println("Missing Open Graph tags")
}
```

### Link Validation
```go
// Extract all links from a page
linkTool := tools.NewExtractLinksTool()
links, _ := linkTool.Execute(ctx, tools.ExtractLinksParams{
    URL: "https://example.com",
})

// Check each link's status
statusTool := tools.NewCheckURLStatusTool()
for _, link := range links.ExternalLinks {
    status, _ := statusTool.Execute(ctx, tools.CheckURLStatusParams{
        URL: link.URL,
        Timeout: 5,
    })
    if status.Status == "dead" {
        log.Printf("Broken link found: %s", link.URL)
    }
}
```

### Content Monitoring
```go
// Monitor a web page for changes
fetchTool := tools.NewFetchWebPageTool()

// Initial fetch
result1, _ := fetchTool.Execute(ctx, tools.FetchWebPageParams{
    URL: "https://example.com/status",
    ExtractText: true,
})
initialContent := result1.Content

// Later fetch
result2, _ := fetchTool.Execute(ctx, tools.FetchWebPageParams{
    URL: "https://example.com/status",
    ExtractText: true,
})

if result2.Content != initialContent {
    log.Println("Page content has changed")
}
```

## Best Practices

1. **Set Appropriate Timeouts**: Use shorter timeouts for status checks and longer ones for content fetching
2. **Handle Redirects Carefully**: Some sites use redirects for mobile versions or localization
3. **Respect robots.txt**: Check the site's robots.txt before scraping
4. **Use Custom User Agents**: Identify your bot properly
5. **Rate Limiting**: Don't overwhelm servers with requests
6. **Error Handling**: Always check for network errors and invalid responses
7. **Cache Results**: Cache extracted data to reduce unnecessary requests

## Error Handling

All tools return errors for:
- Network failures
- Invalid URLs
- Timeouts
- HTTP errors (4xx, 5xx)
- Malformed HTML (for parsing tools)

Example error handling:
```go
result, err := tool.Execute(ctx, params)
if err != nil {
    if strings.Contains(err.Error(), "timeout") {
        // Handle timeout
    } else if strings.Contains(err.Error(), "creating request") {
        // Handle invalid URL
    } else {
        // Handle other errors
    }
}
```

## Performance Tips

1. **Concurrent Requests**: Use goroutines for checking multiple URLs
2. **Connection Pooling**: Reuse HTTP clients when possible
3. **Selective Extraction**: Only extract what you need (e.g., skip media links if not needed)
4. **HEAD Requests**: Use check_url_status for quick availability checks instead of full page fetches
5. **Text Extraction**: Enable text extraction to reduce memory usage when HTML isn't needed

## Future Enhancements

Planned additions to the web tools package:
- `extract_tables`: Extract tabular data from HTML tables
- `submit_form`: Submit forms programmatically
- `extract_structured_data`: Use CSS selectors for precise extraction
- `monitor_changes`: Track specific page elements for changes
- `screenshot_webpage`: Capture visual representations