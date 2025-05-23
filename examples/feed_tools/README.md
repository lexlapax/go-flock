# Feed Tools Example

This example demonstrates how to use the RSS feed fetching tools provided by go-flock to retrieve and process RSS feeds from various sources.

## Overview

The example showcases three different use cases:

1. **Single Feed Fetching** - Basic usage of the `fetch_rss_feed` tool
2. **Concurrent Multi-Feed Fetching** - Fetching multiple feeds simultaneously
3. **Custom Parameters** - Using limit and timeout parameters

## Running the Example

```bash
# From the project root
make build-examples
./bin/feed_tools

# Or directly
cd examples/feed_tools
go run main.go
```

## Features Demonstrated

### 1. Single Feed Fetching
- Fetches a single RSS feed (Hacker News)
- Displays feed metadata (title, description, item count)
- Shows the first 3 items with title, link, and publication date

### 2. Concurrent Multi-Feed Fetching
- Fetches 5 popular RSS feeds concurrently:
  - Hacker News
  - BBC News - World
  - TechCrunch
  - The Verge
  - ArsTechnica
- Demonstrates concurrent execution with goroutines
- Shows fetch duration for performance monitoring
- Displays success/failure status for each feed

### 3. Parameter Customization
- Sets a custom timeout (5 seconds)
- Limits the number of items returned (3 items)
- Shows truncated descriptions for readability

## Tool Used

### `fetch_rss_feed`

Fetches and parses RSS feeds from the given URL.

**Parameters:**
- `url` (required): The URL of the RSS feed to fetch
- `limit` (optional): Maximum number of items to return
- `timeout` (optional): Timeout in seconds (default: 30)

**Returns:**
- Feed title, description, and link
- Array of feed items with title, link, description, published date, and GUID
- Timestamp of when the feed was fetched

## Example Output

```
go-flock RSS Feed Tools Example
===============================

Example 1: Fetching a single RSS feed
-------------------------------------
Feed Title: Hacker News: Front Page
Description: Hacker News RSS
Total Items: 30
Fetched At: 2024-12-15T10:30:00Z

First 3 items:

1. Show HN: I built a tool to analyze code complexity
   Link: https://news.ycombinator.com/item?id=12345
   Published: Sun, 15 Dec 2024 10:00:00 +0000

2. The Architecture of Modern Web Applications
   Link: https://example.com/article
   Published: Sun, 15 Dec 2024 09:30:00 +0000

...

Example 2: Fetching multiple RSS feeds concurrently
--------------------------------------------------
Fetching feeds concurrently...

✅ Hacker News: Success - 30 items (took 250ms)
   Latest: Show HN: I built a tool to analyze code complexity
✅ BBC News - World: Success - 50 items (took 380ms)
   Latest: Breaking: Major climate agreement reached
❌ TechCrunch: Error - timeout (took 5s)
✅ The Verge: Success - 25 items (took 420ms)
   Latest: Apple announces new AI features
✅ ArsTechnica: Success - 40 items (took 310ms)
   Latest: Scientists discover new exoplanet
```

## Error Handling

The example demonstrates proper error handling for:
- Network timeouts
- Invalid URLs
- Server errors
- Malformed RSS feeds

## Use Cases

This tool can be used for:
- News aggregation applications
- Content monitoring systems
- Feed readers
- Data collection pipelines
- Content curation tools

## Next Steps

- Add support for Atom feeds
- Implement feed caching
- Add feed validation
- Support for authenticated feeds
- RSS feed generation tools