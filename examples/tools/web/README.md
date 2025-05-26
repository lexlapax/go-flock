# Web Tools Example

This example demonstrates how to use the web scraping and content extraction tools provided by go-flock.

## Overview

The example showcases seven different use cases:

1. **Fetch Web Page with Text Extraction** - Extract clean text content from HTML
2. **Fetch Raw HTML** - Get the full HTML source with custom headers
3. **Extract All Links** - Extract and categorize all links from a page
4. **Extract External Links Only** - Filter to only external links, grouped by domain
5. **Advanced Link Analysis** - Analyze a tech news site for content patterns
6. **Extract Metadata** - Extract Open Graph, Twitter Cards, and meta tags
7. **Check URL Status** - Check URL availability and response times

## Running the Example

```bash
# From the project root
make build-examples
./bin/web_tools

# Or directly
cd examples/web_tools
go run main.go
```

## Tools Demonstrated

### 1. `fetch_webpage`

Fetches content from web pages with options for text extraction, custom headers, and redirect handling.

**Key Features:**
- Text extraction (strips HTML tags)
- Custom user agent and headers
- Redirect control
- Timeout configuration
- Title extraction

**Parameters:**
- `url` (required): The URL to fetch
- `extract_text`: Extract text content (default: true)
- `timeout`: Request timeout in seconds
- `user_agent`: Custom user agent string
- `headers`: Additional HTTP headers
- `follow_redirects`: Follow HTTP redirects

### 2. `extract_links`

Extracts and categorizes all links from a web page.

**Key Features:**
- Categorizes links (internal, external, email, media)
- Resolves relative URLs to absolute
- Extracts link text and title attributes
- Optional filtering by link type

**Parameters:**
- `url` (required): The URL to extract links from
- `include_external`: Include external links (default: true)
- `include_internal`: Include internal links (default: true)
- `include_email`: Include email links (default: true)
- `include_media`: Include image/media links (default: false)
- `timeout`: Request timeout in seconds

### 3. `extract_metadata`

Extracts metadata from web pages including Open Graph and Twitter Card data.

**Key Features:**
- Standard meta tags (description, author, keywords)
- Open Graph protocol metadata
- Twitter Card metadata
- Custom meta tags parsing
- All metadata in structured format

**Parameters:**
- `url` (required): The URL to extract metadata from
- `timeout`: Request timeout in seconds

### 4. `check_url_status`

Quickly checks if URLs are accessible and measures response times.

**Key Features:**
- URL availability checking (alive/dead/redirect)
- Response time measurement
- Redirect chain tracking
- Content type and length info
- Custom header support

**Parameters:**
- `url` (required): The URL to check
- `timeout`: Request timeout in seconds (default: 10)
- `follow_redirects`: Follow HTTP redirects (default: true)
- `headers`: Additional HTTP headers

## Example Output

```
go-flock Web Tools Example
==========================

Example 1: Fetching a web page with text extraction
---------------------------------------------------
URL: https://example.com
Title: Example Domain
Status Code: 200
Content Type: text/html; charset=UTF-8
Fetched At: 2024-12-15T10:30:00Z
Text Content:
Example Domain
This domain is for use in illustrative examples in documents...

Example 3: Extracting all links from a web page
-----------------------------------------------
Extracted links from: https://news.ycombinator.com
Total links found: 245
Internal links: 203
External links: 42
Email links: 0

First 3 internal links:
  - https://news.ycombinator.com/news
    Text: Hacker News
  - https://news.ycombinator.com/newest
    Text: new
  - https://news.ycombinator.com/ask
    Text: ask

First 3 external links:
  - https://example-startup.com/blog/our-journey
    Text: Our Journey Building a YC Startup
  - https://github.com/awesome/project
    Text: Show HN: Awesome Project

Example 6: Extracting metadata (Open Graph, Twitter Cards)
----------------------------------------------------------
Page: https://www.bbc.com/news
Title: BBC News - Home
Description: Visit BBC News for up-to-the-minute news, breaking news...

Open Graph metadata:
  og:title = BBC News - Home
  og:description = Visit BBC News for up-to-the-minute news...
  og:image = https://www.bbc.com/news/special/2015/newsspec_10857/bbc_news_logo.png
  og:type = website

Twitter Card metadata:
  twitter:card = summary_large_image
  twitter:site = @BBCNews
  twitter:title = BBC News - Home

Example 7: Checking URL status and response times
-------------------------------------------------
Checking multiple URLs:

✅ https://www.google.com
   Status: alive (HTTP 200)
   Response time: 125ms
   Content-Type: text/html; charset=ISO-8859-1

❌ https://httpstat.us/404
   Status: dead (HTTP 404)
   Response time: 89ms

↪️ https://httpstat.us/301
   Status: redirect (HTTP 301)
   Response time: 67ms
   Final URL: https://httpstat.us
   Redirect chain: 1 hops
```

## Use Cases

These tools are perfect for:

1. **Content Aggregation**
   - Fetch articles from multiple sources
   - Extract clean text for processing
   - Collect links for further crawling

2. **SEO Analysis**
   - Check page titles and content
   - Analyze link structure
   - Find broken links

3. **Web Monitoring**
   - Track content changes
   - Monitor new links
   - Detect website updates

4. **Research & Data Collection**
   - Extract article text
   - Collect external references
   - Build link graphs

5. **News Aggregation**
   - Fetch news articles
   - Extract article links
   - Process media content

## Advanced Usage

### Custom Headers Example
```go
params := tools.FetchWebPageParams{
    URL: "https://api.example.com",
    Headers: map[string]string{
        "Authorization": "Bearer token",
        "Accept": "application/json",
    },
}
```

### Link Filtering Example
```go
// Only get external links, no media or emails
params := tools.ExtractLinksParams{
    URL:             "https://example.com",
    IncludeExternal: true,
    IncludeInternal: false,
    IncludeEmail:    false,
    IncludeMedia:    false,
}
```

## Future Enhancements

Additional tools that could be added:
- `extract_tables` - Extract tabular data from HTML tables
- `submit_form` - Form submission capabilities
- `screenshot_webpage` - Visual page capture
- `extract_structured_data` - Extract specific data using CSS selectors
- `monitor_changes` - Track changes in web pages over time