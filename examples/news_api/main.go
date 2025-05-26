// ABOUTME: Example demonstrating news API search functionality
// ABOUTME: Shows how to search for news articles using NewsAPI.org integration

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
	fmt.Println("go-flock News API Search Example")
	fmt.Println("================================")

	// Create news search tool
	searchTool := tools.NewSearchNewsAPITool()
	fmt.Printf("Created tool: %s\n", searchTool.Name())
	fmt.Printf("Description: %s\n\n", searchTool.Description())

	// Check for API key
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  No NEWS_API_KEY environment variable found")
		fmt.Println("You can get a free API key at https://newsapi.org")
		fmt.Println("Then set it: export NEWS_API_KEY=your_api_key")
		fmt.Println("\nUsing demo mode with mock data...")
		runDemoMode()
		return
	}

	// Example 1: Basic search
	fmt.Println("Example 1: Basic News Search")
	fmt.Println("----------------------------")

	searchParams := tools.SearchNewsAPIParams{
		Query:    "artificial intelligence",
		APIKey:   apiKey,
		Language: "en",
		PageSize: 5,
		SortBy:   "popularity",
	}

	ctx := context.Background()
	result, err := searchTool.Execute(ctx, searchParams)
	if err != nil {
		log.Printf("Error searching news: %v\n", err)
		return
	}

	newsResult := result.(*tools.SearchNewsAPIResult)
	fmt.Printf("Status: %s\n", newsResult.Status)
	fmt.Printf("Total Results: %d\n", newsResult.TotalResults)
	fmt.Printf("Articles Retrieved: %d\n\n", len(newsResult.Articles))

	// Display articles
	for i, article := range newsResult.Articles {
		fmt.Printf("%d. %s\n", i+1, article.Title)
		fmt.Printf("   Source: %s\n", article.Source.Name)
		fmt.Printf("   Published: %s\n", article.PublishedAt)
		if article.Description != "" {
			fmt.Printf("   Description: %s\n", article.Description)
		}
		fmt.Printf("   URL: %s\n\n", article.URL)
	}

	// Example 2: Search with date filters
	fmt.Println("\nExample 2: Search with Date Filters")
	fmt.Println("------------------------------------")

	// Search for articles from the last 7 days
	weekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")

	filteredParams := tools.SearchNewsAPIParams{
		Query:    "climate change",
		APIKey:   apiKey,
		Language: "en",
		DateFrom: weekAgo,
		DateTo:   today,
		PageSize: 3,
		SortBy:   "publishedAt",
	}

	result2, err := searchTool.Execute(ctx, filteredParams)
	if err != nil {
		log.Printf("Error with filtered search: %v\n", err)
		return
	}

	newsResult2 := result2.(*tools.SearchNewsAPIResult)
	fmt.Printf("Articles from %s to %s: %d\n\n", weekAgo, today, len(newsResult2.Articles))

	for i, article := range newsResult2.Articles {
		fmt.Printf("%d. %s (%s)\n", i+1, article.Title, article.PublishedAt)
	}

	// Example 3: Domain-specific search
	fmt.Println("\nExample 3: Domain-Specific Search")
	fmt.Println("---------------------------------")

	domainParams := tools.SearchNewsAPIParams{
		Query:    "technology",
		APIKey:   apiKey,
		Domains:  "techcrunch.com,wired.com,arstechnica.com",
		PageSize: 5,
	}

	result3, err := searchTool.Execute(ctx, domainParams)
	if err != nil {
		log.Printf("Error with domain search: %v\n", err)
		return
	}

	newsResult3 := result3.(*tools.SearchNewsAPIResult)
	fmt.Printf("Tech news articles: %d\n\n", len(newsResult3.Articles))

	for i, article := range newsResult3.Articles {
		fmt.Printf("%d. [%s] %s\n", i+1, article.Source.Name, article.Title)
	}

	// Example 4: Pagination
	fmt.Println("\nExample 4: Pagination")
	fmt.Println("---------------------")

	// Get first page
	page1Params := tools.SearchNewsAPIParams{
		Query:    "sports",
		APIKey:   apiKey,
		PageSize: 10,
		Page:     1,
	}

	result4, err := searchTool.Execute(ctx, page1Params)
	if err != nil {
		log.Printf("Error with pagination: %v\n", err)
		return
	}

	newsResult4 := result4.(*tools.SearchNewsAPIResult)
	fmt.Printf("Total sports articles available: %d\n", newsResult4.TotalResults)
	fmt.Printf("Page 1 articles: %d\n", len(newsResult4.Articles))

	// Get second page
	page2Params := page1Params
	page2Params.Page = 2

	result5, err := searchTool.Execute(ctx, page2Params)
	if err != nil {
		log.Printf("Error getting page 2: %v\n", err)
		return
	}

	newsResult5 := result5.(*tools.SearchNewsAPIResult)
	fmt.Printf("Page 2 articles: %d\n", len(newsResult5.Articles))
}

// Demo mode with mock data when no API key is available
func runDemoMode() {
	fmt.Println("Demo Mode - Mock News Results")
	fmt.Println("-----------------------------")

	mockResult := &tools.SearchNewsAPIResult{
		Status:       "ok",
		TotalResults: 3,
		Articles: []tools.NewsArticle{
			{
				Source: tools.NewsSource{
					ID:   "demo-tech",
					Name: "Demo Tech News",
				},
				Author:      "Demo Author",
				Title:       "Breaking: AI Achieves New Milestone in Natural Language Understanding",
				Description: "Researchers announce breakthrough in AI language models with improved comprehension",
				URL:         "https://example.com/ai-breakthrough",
				PublishedAt: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
				Content:     "Full article content would appear here...",
			},
			{
				Source: tools.NewsSource{
					ID:   "demo-science",
					Name: "Demo Science Daily",
				},
				Author:      "Science Reporter",
				Title:       "Climate Scientists Discover New Carbon Capture Method",
				Description: "Innovative technique could remove CO2 from atmosphere more efficiently",
				URL:         "https://example.com/climate-news",
				PublishedAt: time.Now().Add(-5 * time.Hour).Format(time.RFC3339),
				Content:     "Article about climate science...",
			},
			{
				Source: tools.NewsSource{
					ID:   "demo-business",
					Name: "Demo Business Wire",
				},
				Author:      "Business Editor",
				Title:       "Tech Stocks Rally on Strong Earnings Reports",
				Description: "Major technology companies exceed quarterly expectations",
				URL:         "https://example.com/business-news",
				PublishedAt: time.Now().Add(-8 * time.Hour).Format(time.RFC3339),
				Content:     "Financial news content...",
			},
		},
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
	}

	fmt.Printf("Status: %s\n", mockResult.Status)
	fmt.Printf("Total Results: %d\n", mockResult.TotalResults)
	fmt.Printf("Articles Retrieved: %d\n\n", len(mockResult.Articles))

	for i, article := range mockResult.Articles {
		fmt.Printf("%d. %s\n", i+1, article.Title)
		fmt.Printf("   Source: %s\n", article.Source.Name)
		fmt.Printf("   Author: %s\n", article.Author)
		fmt.Printf("   Published: %s\n", article.PublishedAt)
		fmt.Printf("   Description: %s\n", article.Description)
		fmt.Printf("   URL: %s\n\n", article.URL)
	}

	fmt.Println("\n💡 To use real news data:")
	fmt.Println("1. Sign up for a free API key at https://newsapi.org")
	fmt.Println("2. Export it: export NEWS_API_KEY=your_api_key")
	fmt.Println("3. Run this example again")
}
