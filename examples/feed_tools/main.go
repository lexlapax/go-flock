// ABOUTME: Example demonstrating RSS feed fetching tools
// ABOUTME: Shows how to fetch and process multiple RSS feeds concurrently

package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lexlapax/go-flock/pkg/tools"
)

// Popular RSS feeds for demonstration
var feeds = []struct {
	name string
	url  string
}{
	{"Hacker News", "https://hnrss.org/frontpage"},
	{"BBC News - World", "http://feeds.bbci.co.uk/news/world/rss.xml"},
	{"TechCrunch", "https://techcrunch.com/feed/"},
	{"The Verge", "https://www.theverge.com/rss/index.xml"},
	{"ArsTechnica", "http://feeds.arstechnica.com/arstechnica/index"},
}

func main() {
	fmt.Println("go-flock RSS Feed Tools Example")
	fmt.Println("===============================")
	fmt.Println()

	ctx := context.Background()

	// Example 1: Fetch a single RSS feed
	fmt.Println("Example 1: Fetching a single RSS feed")
	fmt.Println("-------------------------------------")
	fetchSingleFeed(ctx)
	fmt.Println()

	// Example 2: Fetch multiple feeds concurrently
	fmt.Println("Example 2: Fetching multiple RSS feeds concurrently")
	fmt.Println("--------------------------------------------------")
	fetchMultipleFeeds(ctx)
	fmt.Println()

	// Example 3: Fetch with custom parameters
	fmt.Println("Example 3: Fetching with custom parameters (limit and timeout)")
	fmt.Println("-------------------------------------------------------------")
	fetchWithParameters(ctx)
}

func fetchSingleFeed(ctx context.Context) {
	tool := tools.NewFetchRSSFeedTool()

	params := tools.FetchRSSFeedParams{
		URL: "https://hnrss.org/frontpage",
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Printf("Error fetching feed: %v", err)
		return
	}

	if feedResult, ok := result.(*tools.FetchRSSFeedResult); ok {
		fmt.Printf("Feed Title: %s\n", feedResult.Title)
		fmt.Printf("Description: %s\n", feedResult.Description)
		fmt.Printf("Total Items: %d\n", len(feedResult.Items))
		fmt.Printf("Fetched At: %s\n", feedResult.FetchedAt)

		if len(feedResult.Items) > 0 {
			fmt.Println("\nFirst 3 items:")
			for i, item := range feedResult.Items {
				if i >= 3 {
					break
				}
				fmt.Printf("\n%d. %s\n", i+1, item.Title)
				fmt.Printf("   Link: %s\n", item.Link)
				fmt.Printf("   Published: %s\n", item.Published)
			}
		}
	}
}

func fetchMultipleFeeds(ctx context.Context) {
	tool := tools.NewFetchRSSFeedTool()

	var wg sync.WaitGroup
	results := make(chan feedResult, len(feeds))

	// Fetch all feeds concurrently
	for _, feed := range feeds {
		wg.Add(1)
		go func(name, url string) {
			defer wg.Done()

			params := tools.FetchRSSFeedParams{
				URL:     url,
				Limit:   5,
				Timeout: 10,
			}

			start := time.Now()
			result, err := tool.Execute(ctx, params)
			duration := time.Since(start)

			if err != nil {
				results <- feedResult{
					name:     name,
					err:      err,
					duration: duration,
				}
				return
			}

			if feedData, ok := result.(*tools.FetchRSSFeedResult); ok {
				results <- feedResult{
					name:     name,
					feed:     feedData,
					duration: duration,
				}
			}
		}(feed.name, feed.url)
	}

	// Wait for all feeds to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Display results as they come in
	fmt.Println("Fetching feeds concurrently...")
	fmt.Println()

	for result := range results {
		if result.err != nil {
			fmt.Printf("❌ %s: Error - %v (took %v)\n", result.name, result.err, result.duration)
		} else {
			fmt.Printf("✅ %s: Success - %d items (took %v)\n",
				result.name, len(result.feed.Items), result.duration)

			// Show first item from each feed
			if len(result.feed.Items) > 0 {
				fmt.Printf("   Latest: %s\n", result.feed.Items[0].Title)
			}
		}
	}
}

func fetchWithParameters(ctx context.Context) {
	tool := tools.NewFetchRSSFeedTool()

	// Fetch with a 5-second timeout and limit to 3 items
	params := tools.FetchRSSFeedParams{
		URL:     "https://techcrunch.com/feed/",
		Limit:   3,
		Timeout: 5,
	}

	fmt.Println("Fetching TechCrunch feed with:")
	fmt.Printf("- Limit: %d items\n", params.Limit)
	fmt.Printf("- Timeout: %d seconds\n", params.Timeout)
	fmt.Println()

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Printf("Error fetching feed: %v", err)
		return
	}

	if feedResult, ok := result.(*tools.FetchRSSFeedResult); ok {
		fmt.Printf("Feed: %s\n", feedResult.Title)
		fmt.Printf("Items returned: %d\n", len(feedResult.Items))

		for i, item := range feedResult.Items {
			fmt.Printf("\n%d. %s\n", i+1, item.Title)
			fmt.Printf("   Published: %s\n", item.Published)

			// Show truncated description
			desc := item.Description
			if len(desc) > 150 {
				desc = desc[:150] + "..."
			}
			fmt.Printf("   Description: %s\n", desc)
		}
	}
}

// Helper struct for concurrent results
type feedResult struct {
	name     string
	feed     *tools.FetchRSSFeedResult
	err      error
	duration time.Duration
}
