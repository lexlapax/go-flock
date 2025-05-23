# Feed Tools

The feed tools package provides utilities for fetching and parsing RSS and Atom feeds.

## Available Tools

### fetch_rss_feed

Fetches and parses RSS feeds from the given URL.

**Function**: `NewFetchRSSFeedTool()`

**Parameters**:
- `url` (string, required): The URL of the RSS feed to fetch
- `limit` (integer, optional): Maximum number of items to return (default: all)
- `timeout` (integer, optional): Timeout in seconds (default: 30, max: 300)

**Returns**:
```json
{
  "title": "Hacker News: Front Page",
  "description": "Hacker News RSS",
  "link": "https://news.ycombinator.com/",
  "items": [
    {
      "title": "Show HN: I built a tool to analyze code complexity",
      "link": "https://news.ycombinator.com/item?id=12345",
      "description": "A new approach to measuring and visualizing code complexity...",
      "published": "Mon, 15 Dec 2024 10:00:00 +0000",
      "guid": "https://news.ycombinator.com/item?id=12345"
    }
  ],
  "fetched_at": "2024-12-15T10:30:00Z"
}
```

**Example Usage**:
```go
tool := tools.NewFetchRSSFeedTool()
params := tools.FetchRSSFeedParams{
    URL:     "https://hnrss.org/frontpage",
    Limit:   10,
    Timeout: 15,
}
result, err := tool.Execute(ctx, params)
```

## Supported Feed Formats

Currently supports:
- **RSS 2.0**: The most common RSS format
- **RSS 1.0**: Older RSS format (planned)
- **Atom**: Modern feed format (planned)

## Use Cases

### News Aggregation
```go
// Fetch news from multiple sources
feeds := []string{
    "https://hnrss.org/frontpage",
    "http://feeds.bbci.co.uk/news/world/rss.xml",
    "https://techcrunch.com/feed/",
}

tool := tools.NewFetchRSSFeedTool()
var allArticles []FeedItem

for _, feedURL := range feeds {
    result, err := tool.Execute(ctx, tools.FetchRSSFeedParams{
        URL:   feedURL,
        Limit: 5,  // Get top 5 from each
    })
    if err != nil {
        log.Printf("Error fetching %s: %v", feedURL, err)
        continue
    }
    
    feedResult := result.(*tools.FetchRSSFeedResult)
    allArticles = append(allArticles, feedResult.Items...)
}
```

### Concurrent Feed Fetching
```go
// Fetch multiple feeds concurrently
var wg sync.WaitGroup
results := make(chan *tools.FetchRSSFeedResult, len(feeds))

tool := tools.NewFetchRSSFeedTool()

for _, feedURL := range feeds {
    wg.Add(1)
    go func(url string) {
        defer wg.Done()
        
        result, err := tool.Execute(ctx, tools.FetchRSSFeedParams{
            URL:     url,
            Timeout: 10,
        })
        if err != nil {
            log.Printf("Error: %v", err)
            return
        }
        
        results <- result.(*tools.FetchRSSFeedResult)
    }(feedURL)
}

go func() {
    wg.Wait()
    close(results)
}()

// Process results as they arrive
for result := range results {
    fmt.Printf("Feed: %s (%d items)\n", result.Title, len(result.Items))
}
```

### Feed Monitoring
```go
// Monitor a feed for new items
tool := tools.NewFetchRSSFeedTool()
seenItems := make(map[string]bool)

// Initial fetch
result, _ := tool.Execute(ctx, tools.FetchRSSFeedParams{
    URL: "https://example.com/feed.xml",
})
feedResult := result.(*tools.FetchRSSFeedResult)

// Mark existing items as seen
for _, item := range feedResult.Items {
    seenItems[item.GUID] = true
}

// Periodic check for new items
ticker := time.NewTicker(5 * time.Minute)
for range ticker.C {
    result, err := tool.Execute(ctx, tools.FetchRSSFeedParams{
        URL: "https://example.com/feed.xml",
    })
    if err != nil {
        continue
    }
    
    feedResult := result.(*tools.FetchRSSFeedResult)
    for _, item := range feedResult.Items {
        if !seenItems[item.GUID] {
            fmt.Printf("New item: %s\n", item.Title)
            seenItems[item.GUID] = true
        }
    }
}
```

### Content Filtering
```go
// Filter feed items by keywords
tool := tools.NewFetchRSSFeedTool()
keywords := []string{"golang", "kubernetes", "cloud"}

result, _ := tool.Execute(ctx, tools.FetchRSSFeedParams{
    URL: "https://techcrunch.com/feed/",
})
feedResult := result.(*tools.FetchRSSFeedResult)

var relevantItems []tools.FeedItem
for _, item := range feedResult.Items {
    content := strings.ToLower(item.Title + " " + item.Description)
    for _, keyword := range keywords {
        if strings.Contains(content, keyword) {
            relevantItems = append(relevantItems, item)
            break
        }
    }
}
```

## Best Practices

1. **Set Reasonable Limits**: Use the `limit` parameter to avoid processing too many items
2. **Handle Timeouts**: RSS feeds can be slow; set appropriate timeouts
3. **Cache Results**: Store fetched feeds to reduce server load
4. **Respect Rate Limits**: Don't fetch the same feed too frequently
5. **Validate Feed URLs**: Check that URLs are valid before fetching
6. **Error Recovery**: Handle malformed feeds gracefully

## Error Handling

The tool handles various error conditions:
- Invalid URLs
- Network timeouts
- Server errors (404, 500, etc.)
- Malformed XML
- Empty feeds

Example error handling:
```go
result, err := tool.Execute(ctx, params)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "timeout"):
        log.Println("Feed fetch timed out")
    case strings.Contains(err.Error(), "parsing RSS feed"):
        log.Println("Invalid RSS format")
    case strings.Contains(err.Error(), "status code"):
        log.Println("Server returned error")
    default:
        log.Printf("Unexpected error: %v", err)
    }
}
```

## Performance Optimization

### Concurrent Processing
```go
// Process feed items concurrently
feedResult := result.(*tools.FetchRSSFeedResult)
itemChan := make(chan tools.FeedItem, len(feedResult.Items))

// Producer
go func() {
    for _, item := range feedResult.Items {
        itemChan <- item
    }
    close(itemChan)
}()

// Consumers
var wg sync.WaitGroup
for i := 0; i < 4; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        for item := range itemChan {
            // Process item
            processItem(item)
        }
    }()
}
wg.Wait()
```

### Caching Strategy
```go
type FeedCache struct {
    mu    sync.RWMutex
    cache map[string]*CachedFeed
}

type CachedFeed struct {
    Result    *tools.FetchRSSFeedResult
    FetchedAt time.Time
}

func (fc *FeedCache) Get(url string, maxAge time.Duration) (*tools.FetchRSSFeedResult, bool) {
    fc.mu.RLock()
    defer fc.mu.RUnlock()
    
    cached, ok := fc.cache[url]
    if !ok {
        return nil, false
    }
    
    if time.Since(cached.FetchedAt) > maxAge {
        return nil, false
    }
    
    return cached.Result, true
}
```

## Integration with Other Tools

### Combine with Web Tools
```go
// 1. Fetch RSS feed
feedTool := tools.NewFetchRSSFeedTool()
feedResult, _ := feedTool.Execute(ctx, tools.FetchRSSFeedParams{
    URL: "https://example.com/feed.xml",
})

// 2. Extract full article content
webTool := tools.NewFetchWebPageTool()
for _, item := range feedResult.Items[:3] {
    article, _ := webTool.Execute(ctx, tools.FetchWebPageParams{
        URL:         item.Link,
        ExtractText: true,
    })
    
    // Process full article text
    fmt.Printf("Article: %s\nContent: %s\n", item.Title, article.Content)
}
```

### Feed Discovery
```go
// Use web tools to discover RSS feeds on a website
linkTool := tools.NewExtractLinksTool()
links, _ := linkTool.Execute(ctx, tools.ExtractLinksParams{
    URL: "https://example.com",
})

// Look for RSS feed links
var feedURLs []string
for _, link := range links.InternalLinks {
    if strings.Contains(link.URL, "feed") || 
       strings.Contains(link.URL, "rss") ||
       strings.Contains(link.URL, ".xml") {
        feedURLs = append(feedURLs, link.URL)
    }
}
```

## Future Enhancements

Planned additions:
- **Atom Support**: Parse Atom 1.0 feeds
- **Feed Validation**: Validate feed structure before parsing
- **OPML Support**: Import/export OPML subscription lists
- **Feed Generation**: Create RSS feeds from data
- **Podcast Support**: Enhanced support for podcast RSS extensions
- **JSON Feed**: Support for JSON Feed format
- **Conditional GET**: Support for If-Modified-Since headers
- **Feed Autodiscovery**: Automatically find feeds on websites