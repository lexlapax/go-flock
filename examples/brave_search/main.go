// ABOUTME: Example demonstrating Brave web search functionality
// ABOUTME: Shows how to search the web using Brave Search API with AI summaries

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lexlapax/go-flock/pkg/tools"
)

func main() {
	fmt.Println("go-flock Brave Web Search Example")
	fmt.Println("=================================")

	// Create Brave search tool
	searchTool := tools.NewSearchWebBraveTool()
	fmt.Printf("Created tool: %s\n", searchTool.Name())
	fmt.Printf("Description: %s\n\n", searchTool.Description())

	// Check for API key
	apiKey := os.Getenv("BRAVE_SEARCH_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  No BRAVE_SEARCH_API_KEY environment variable found")
		fmt.Println("You can get a free API key at https://api.search.brave.com")
		fmt.Println("Then set it: export BRAVE_SEARCH_API_KEY=your_api_key")
		fmt.Println("\nUsing demo mode with mock data...")
		runDemoMode()
		return
	}

	// Example 1: Basic web search
	fmt.Println("Example 1: Basic Web Search")
	fmt.Println("---------------------------")

	searchParams := tools.SearchWebBraveParams{
		Query:  "golang web frameworks 2024",
		APIKey: apiKey,
		Count:  5,
	}

	ctx := context.Background()
	result, err := searchTool.Execute(ctx, searchParams)
	if err != nil {
		log.Printf("Error searching web: %v\n", err)
		return
	}

	braveResult := result.(*tools.SearchWebBraveResult)
	fmt.Printf("Query: %s\n", braveResult.Query)
	fmt.Printf("Type: %s\n", braveResult.Type)
	fmt.Printf("Web Results: %d\n\n", len(braveResult.Web.Results))

	// Display web results
	for i, webResult := range braveResult.Web.Results {
		fmt.Printf("%d. %s\n", i+1, webResult.Title)
		fmt.Printf("   URL: %s\n", webResult.URL)
		if webResult.Age != "" {
			fmt.Printf("   Age: %s\n", webResult.Age)
		}
		if webResult.Description != "" {
			fmt.Printf("   Description: %s\n", webResult.Description)
		}
		fmt.Println()
	}

	// Example 2: Search with AI Summary
	fmt.Println("\nExample 2: Search with AI Summary")
	fmt.Println("---------------------------------")

	summaryParams := tools.SearchWebBraveParams{
		Query:   "what is quantum computing explained simply",
		APIKey:  apiKey,
		Count:   3,
		Summary: true,
	}

	result2, err := searchTool.Execute(ctx, summaryParams)
	if err != nil {
		log.Printf("Error with summary search: %v\n", err)
		return
	}

	braveResult2 := result2.(*tools.SearchWebBraveResult)

	// Display AI summary if available
	if len(braveResult2.Summary) > 0 {
		fmt.Println("AI Summary:")
		for _, summary := range braveResult2.Summary {
			fmt.Printf("  %s\n", summary.Text)
		}
		fmt.Println()
	}

	fmt.Printf("Supporting results: %d\n\n", len(braveResult2.Web.Results))
	for i, webResult := range braveResult2.Web.Results {
		fmt.Printf("%d. %s - %s\n", i+1, webResult.Title, webResult.URL)
	}

	// Example 3: News-focused search
	fmt.Println("\nExample 3: News-Focused Search")
	fmt.Println("------------------------------")

	newsParams := tools.SearchWebBraveParams{
		Query:        "artificial intelligence breakthroughs",
		APIKey:       apiKey,
		Count:        10,
		Freshness:    "pd", // Past day
		ResultFilter: []string{"news"},
		Country:      "US",
	}

	result3, err := searchTool.Execute(ctx, newsParams)
	if err != nil {
		log.Printf("Error with news search: %v\n", err)
		return
	}

	braveResult3 := result3.(*tools.SearchWebBraveResult)

	if len(braveResult3.News.Results) > 0 {
		fmt.Printf("News articles from past day: %d\n\n", len(braveResult3.News.Results))
		for i, newsResult := range braveResult3.News.Results {
			fmt.Printf("%d. %s\n", i+1, newsResult.Title)
			if newsResult.Source != nil {
				fmt.Printf("   Source: %s\n", newsResult.Source.Name)
			}
			if newsResult.Age != "" {
				fmt.Printf("   Age: %s\n", newsResult.Age)
			}
			fmt.Printf("   URL: %s\n\n", newsResult.URL)
		}
	} else {
		fmt.Println("No news results found in the response")
	}

	// Example 4: Multi-type search
	fmt.Println("\nExample 4: Multi-Type Search (Web + News + Videos)")
	fmt.Println("--------------------------------------------------")

	multiParams := tools.SearchWebBraveParams{
		Query:        "SpaceX starship launch",
		APIKey:       apiKey,
		Count:        10,
		ResultFilter: []string{"web", "news", "videos"},
	}

	result4, err := searchTool.Execute(ctx, multiParams)
	if err != nil {
		log.Printf("Error with multi-type search: %v\n", err)
		return
	}

	braveResult4 := result4.(*tools.SearchWebBraveResult)

	// Display mixed results
	if len(braveResult4.Mixed.Main) > 0 {
		fmt.Println("Mixed results order:")
		for _, item := range braveResult4.Mixed.Main {
			fmt.Printf("  - Type: %s, Index: %d\n", item.Type, item.Index)
		}
		fmt.Println()
	}

	fmt.Printf("Web results: %d\n", len(braveResult4.Web.Results))
	fmt.Printf("News results: %d\n", len(braveResult4.News.Results))
	fmt.Printf("Video results: %d\n", len(braveResult4.Videos.Results))

	// Show some video results if available
	if len(braveResult4.Videos.Results) > 0 {
		fmt.Println("\nVideo Results:")
		for i, video := range braveResult4.Videos.Results {
			if i >= 3 {
				break // Show only first 3
			}
			fmt.Printf("%d. %s\n", i+1, video.Title)
			if video.Duration != "" {
				fmt.Printf("   Duration: %s\n", video.Duration)
			}
			if video.Creator != "" {
				fmt.Printf("   Creator: %s\n", video.Creator)
			}
			fmt.Printf("   URL: %s\n\n", video.URL)
		}
	}

	// Example 5: Safe search with language preference
	fmt.Println("\nExample 5: Safe Search with Language Preference")
	fmt.Println("----------------------------------------------")

	safeParams := tools.SearchWebBraveParams{
		Query:      "educational content for kids",
		APIKey:     apiKey,
		Count:      5,
		SafeSearch: "strict",
		SearchLang: "en",
		Country:    "US",
	}

	result5, err := searchTool.Execute(ctx, safeParams)
	if err != nil {
		log.Printf("Error with safe search: %v\n", err)
		return
	}

	braveResult5 := result5.(*tools.SearchWebBraveResult)
	fmt.Printf("Safe search results: %d\n", len(braveResult5.Web.Results))

	// Show InfoBox if available
	if braveResult5.InfoBox != nil {
		fmt.Println("\nInfoBox:")
		fmt.Printf("  Title: %s\n", braveResult5.InfoBox.Title)
		if braveResult5.InfoBox.Description != "" {
			fmt.Printf("  Description: %s\n", braveResult5.InfoBox.Description)
		}
	}

	// Show FAQ if available
	if len(braveResult5.FAQ) > 0 {
		fmt.Println("\nFrequently Asked Questions:")
		for i, faq := range braveResult5.FAQ {
			fmt.Printf("%d. Q: %s\n", i+1, faq.Question)
			fmt.Printf("   A: %s\n", faq.Answer)
		}
	}
}

// Demo mode with mock data when no API key is available
func runDemoMode() {
	fmt.Println("\nDemo Mode - Mock Brave Search Results")
	fmt.Println("-------------------------------------")

	mockResult := &tools.SearchWebBraveResult{
		Query: "golang web frameworks",
		Type:  "search",
		Web: tools.BraveWebResults{
			Type: "search",
			Results: []tools.BraveWebResult{
				{
					Title:       "Gin Web Framework - Fast and Minimalist",
					URL:         "https://gin-gonic.com",
					Description: "Gin is a HTTP web framework written in Go. It features a Martini-like API with much better performance.",
					Age:         "2 days ago",
					Language:    "en",
				},
				{
					Title:       "Echo - High performance Go web framework",
					URL:         "https://echo.labstack.com",
					Description: "Echo is a high performance, extensible, minimalist Go web framework",
					Age:         "1 week ago",
					Thumbnail: &tools.BraveThumbnail{
						Src:    "https://example.com/echo-thumb.jpg",
						Height: 200,
						Width:  300,
					},
				},
				{
					Title:       "Fiber - Express inspired web framework",
					URL:         "https://gofiber.io",
					Description: "Fiber is an Express inspired web framework built on top of Fasthttp",
					Age:         "3 days ago",
				},
			},
		},
		News: tools.BraveNewsResults{
			Type: "news",
			Results: []tools.BraveNewsResult{
				{
					Title:       "Go 1.22 Released with Major Performance Improvements",
					URL:         "https://news.example.com/go-release",
					Description: "The latest Go release brings significant improvements to web framework performance",
					Age:         "1 day ago",
					Source: &tools.BraveNewsSource{
						Name: "Tech News Daily",
						URL:  "https://technewsdaily.com",
					},
				},
			},
		},
		Summary: []tools.BraveSummary{
			{
				Type: "summary",
				Key:  "brave_search_llm_summary",
				Text: "Go (Golang) offers several popular web frameworks. Gin is known for its speed and minimalist design, Echo provides high performance with extensibility, and Fiber offers an Express-like experience. Other notable frameworks include Beego, Revel, and Buffalo. Choose based on your project's needs for performance, features, and community support.",
			},
		},
		InfoBox: &tools.BraveInfoBox{
			Type:        "infobox",
			Title:       "Go Programming Language",
			Description: "Go is an open source programming language that makes it simple to build secure, scalable systems",
			URL:         "https://go.dev",
		},
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// Display mock results
	fmt.Printf("\nQuery: %s\n", mockResult.Query)
	fmt.Printf("Type: %s\n\n", mockResult.Type)

	// Show AI Summary
	if len(mockResult.Summary) > 0 {
		fmt.Println("AI Summary:")
		fmt.Printf("  %s\n\n", mockResult.Summary[0].Text)
	}

	// Show InfoBox
	if mockResult.InfoBox != nil {
		fmt.Println("InfoBox:")
		fmt.Printf("  %s - %s\n\n", mockResult.InfoBox.Title, mockResult.InfoBox.Description)
	}

	// Show Web Results
	fmt.Printf("Web Results (%d):\n", len(mockResult.Web.Results))
	for i, result := range mockResult.Web.Results {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   URL: %s\n", result.URL)
		fmt.Printf("   Description: %s\n", result.Description)
		fmt.Printf("   Age: %s\n\n", result.Age)
	}

	// Show News Results
	if len(mockResult.News.Results) > 0 {
		fmt.Printf("News Results (%d):\n", len(mockResult.News.Results))
		for i, result := range mockResult.News.Results {
			fmt.Printf("%d. %s\n", i+1, result.Title)
			if result.Source != nil {
				fmt.Printf("   Source: %s\n", result.Source.Name)
			}
			fmt.Printf("   Age: %s\n", result.Age)
			fmt.Printf("   Description: %s\n\n", result.Description)
		}
	}

	fmt.Println("\n💡 To use real Brave Search results:")
	fmt.Println("1. Sign up for a free API key at https://api.search.brave.com")
	fmt.Println("2. Export it: export BRAVE_SEARCH_API_KEY=your_api_key")
	fmt.Println("3. Run this example again")
	fmt.Println("\nBrave Search advantages:")
	fmt.Println("- Privacy-focused search")
	fmt.Println("- AI-powered summaries")
	fmt.Println("- Multi-content type results (web, news, images, videos)")
	fmt.Println("- No tracking or profiling")
}
