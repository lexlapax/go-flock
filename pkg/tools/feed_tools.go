// ABOUTME: Feed processing tools for RSS, Atom, and other feed formats
// ABOUTME: Provides tools to fetch and parse various feed formats

package tools

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	domain "github.com/lexlapax/go-llms/pkg/agent/domain"
	"github.com/lexlapax/go-llms/pkg/agent/tools"
	sdomain "github.com/lexlapax/go-llms/pkg/schema/domain"
)

// RSS Feed structures
type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

// Tool Parameters
type FetchRSSFeedParams struct {
	URL     string `json:"url" description:"The URL of the RSS feed to fetch"`
	Limit   int    `json:"limit,omitempty" description:"Maximum number of items to return (default: all)"`
	Timeout int    `json:"timeout,omitempty" description:"Timeout in seconds (default: 30)"`
}

// Tool Results
type FetchRSSFeedResult struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Link        string     `json:"link"`
	Items       []FeedItem `json:"items"`
	FetchedAt   string     `json:"fetched_at"`
}

type FeedItem struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Published   string `json:"published"`
	GUID        string `json:"guid"`
}

// Schema definitions
var FetchRSSFeedParamSchema = &sdomain.Schema{
	Type:        "object",
	Description: "Parameters for fetching an RSS feed",
	Properties: map[string]sdomain.Property{
		"url": {
			Type:        "string",
			Description: "The URL of the RSS feed to fetch",
		},
		"limit": {
			Type:        "integer",
			Description: "Maximum number of items to return (default: all)",
			Minimum:     float64Ptr(1),
		},
		"timeout": {
			Type:        "integer",
			Description: "Timeout in seconds (default: 30)",
			Minimum:     float64Ptr(1),
			Maximum:     float64Ptr(300),
		},
	},
	Required: []string{"url"},
}

// NewFetchRSSFeedTool creates a new RSS feed fetching tool
func NewFetchRSSFeedTool() domain.Tool {
	return tools.NewTool(
		"fetch_rss_feed",
		"Fetches and parses an RSS feed from the given URL",
		fetchRSSFeedHandler,
		FetchRSSFeedParamSchema,
	)
}

func fetchRSSFeedHandler(ctx context.Context, params FetchRSSFeedParams) (*FetchRSSFeedResult, error) {
	// Set default timeout
	timeout := 30
	if params.Timeout > 0 {
		timeout = params.Timeout
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", params.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set User-Agent
	req.Header.Set("User-Agent", "go-flock/1.0 RSS Reader")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching RSS feed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// Parse RSS feed
	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("parsing RSS feed: %w", err)
	}

	// Convert to result format
	result := &FetchRSSFeedResult{
		Title:       feed.Channel.Title,
		Description: feed.Channel.Description,
		Link:        feed.Channel.Link,
		Items:       make([]FeedItem, 0, len(feed.Channel.Items)),
		FetchedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// Apply limit if specified
	itemCount := len(feed.Channel.Items)
	if params.Limit > 0 && params.Limit < itemCount {
		itemCount = params.Limit
	}

	// Convert items
	for i := range itemCount {
		item := feed.Channel.Items[i]
		result.Items = append(result.Items, FeedItem{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Published:   item.PubDate,
			GUID:        item.GUID,
		})
	}

	return result, nil
}

// Helper function for float64 pointer
func float64Ptr(f float64) *float64 {
	return &f
}
