package tools

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Sample RSS feed for testing
const sampleRSSFeed = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test News Feed</title>
    <link>https://example.com</link>
    <description>A test RSS feed for unit testing</description>
    <item>
      <title>Breaking News: Go Testing Made Easy</title>
      <link>https://example.com/news/1</link>
      <description>Learn how to write effective tests in Go</description>
      <pubDate>Mon, 15 Dec 2024 10:30:00 GMT</pubDate>
      <guid>https://example.com/news/1</guid>
    </item>
    <item>
      <title>Go-Flock Announces New Features</title>
      <link>https://example.com/news/2</link>
      <description>The latest features in go-flock library</description>
      <pubDate>Sun, 14 Dec 2024 15:45:00 GMT</pubDate>
      <guid>https://example.com/news/2</guid>
    </item>
    <item>
      <title>LLM Integration Best Practices</title>
      <link>https://example.com/news/3</link>
      <description>Tips for integrating LLMs into your applications</description>
      <pubDate>Sat, 13 Dec 2024 09:00:00 GMT</pubDate>
      <guid>https://example.com/news/3</guid>
    </item>
  </channel>
</rss>`

func TestNewFetchRSSFeedTool(t *testing.T) {
	tool := NewFetchRSSFeedTool()
	
	if tool == nil {
		t.Fatal("NewFetchRSSFeedTool returned nil")
	}
	
	if tool.Name() != "fetch_rss_feed" {
		t.Errorf("Expected tool name 'fetch_rss_feed', got %s", tool.Name())
	}
	
	if tool.Description() != "Fetches and parses an RSS feed from the given URL" {
		t.Errorf("Unexpected tool description: %s", tool.Description())
	}
	
	if tool.ParameterSchema() == nil {
		t.Error("Tool parameter schema is nil")
	}
}

func TestFetchRSSFeedHandler_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, sampleRSSFeed)
	}))
	defer server.Close()
	
	ctx := context.Background()
	params := FetchRSSFeedParams{
		URL:     server.URL,
		Timeout: 30,
	}
	
	result, err := fetchRSSFeedHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	// Verify result
	if result.Title != "Test News Feed" {
		t.Errorf("Expected title 'Test News Feed', got %s", result.Title)
	}
	
	if result.Description != "A test RSS feed for unit testing" {
		t.Errorf("Expected description 'A test RSS feed for unit testing', got %s", result.Description)
	}
	
	if result.Link != "https://example.com" {
		t.Errorf("Expected link 'https://example.com', got %s", result.Link)
	}
	
	if len(result.Items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(result.Items))
	}
	
	// Check first item
	if len(result.Items) > 0 {
		item := result.Items[0]
		if item.Title != "Breaking News: Go Testing Made Easy" {
			t.Errorf("Expected first item title 'Breaking News: Go Testing Made Easy', got %s", item.Title)
		}
		if item.Link != "https://example.com/news/1" {
			t.Errorf("Expected first item link 'https://example.com/news/1', got %s", item.Link)
		}
	}
}

func TestFetchRSSFeedHandler_WithLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		fmt.Fprint(w, sampleRSSFeed)
	}))
	defer server.Close()
	
	ctx := context.Background()
	params := FetchRSSFeedParams{
		URL:   server.URL,
		Limit: 2,
	}
	
	result, err := fetchRSSFeedHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if len(result.Items) != 2 {
		t.Errorf("Expected 2 items with limit, got %d", len(result.Items))
	}
}

func TestFetchRSSFeedHandler_InvalidURL(t *testing.T) {
	ctx := context.Background()
	params := FetchRSSFeedParams{
		URL: "not-a-valid-url",
	}
	
	_, err := fetchRSSFeedHandler(ctx, params)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestFetchRSSFeedHandler_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	
	ctx := context.Background()
	params := FetchRSSFeedParams{
		URL: server.URL,
	}
	
	_, err := fetchRSSFeedHandler(ctx, params)
	if err == nil {
		t.Error("Expected error for server error, got nil")
	}
}

func TestFetchRSSFeedHandler_InvalidXML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		fmt.Fprint(w, "This is not valid XML")
	}))
	defer server.Close()
	
	ctx := context.Background()
	params := FetchRSSFeedParams{
		URL: server.URL,
	}
	
	_, err := fetchRSSFeedHandler(ctx, params)
	if err == nil {
		t.Error("Expected error for invalid XML, got nil")
	}
}

func TestFetchRSSFeedHandler_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(2 * time.Second)
		fmt.Fprint(w, sampleRSSFeed)
	}))
	defer server.Close()
	
	ctx := context.Background()
	params := FetchRSSFeedParams{
		URL:     server.URL,
		Timeout: 1, // 1 second timeout
	}
	
	_, err := fetchRSSFeedHandler(ctx, params)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "deadline") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

func TestFetchRSSFeedHandler_EmptyFeed(t *testing.T) {
	emptyFeed := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Empty Feed</title>
    <link>https://example.com</link>
    <description>An empty RSS feed</description>
  </channel>
</rss>`
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		fmt.Fprint(w, emptyFeed)
	}))
	defer server.Close()
	
	ctx := context.Background()
	params := FetchRSSFeedParams{
		URL: server.URL,
	}
	
	result, err := fetchRSSFeedHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if len(result.Items) != 0 {
		t.Errorf("Expected 0 items for empty feed, got %d", len(result.Items))
	}
}