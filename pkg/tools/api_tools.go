// ABOUTME: API interaction tools for REST, GraphQL, and other API protocols
// ABOUTME: Provides tools to interact with various API services

package tools

import (
	"context"
	"fmt"

	domain "github.com/lexlapax/go-llms/pkg/agent/domain"
	"github.com/lexlapax/go-llms/pkg/agent/tools"
	sdomain "github.com/lexlapax/go-llms/pkg/schema/domain"
)

// Tool Parameters
type SearchNewsAPIParams struct {
	Query          string `json:"query" description:"Search query for news articles"`
	APIKey         string `json:"api_key" description:"API key for the news service"`
	Language       string `json:"language,omitempty" description:"Language code (e.g., 'en', 'es')"`
	SortBy         string `json:"sort_by,omitempty" description:"Sort order: 'relevancy', 'popularity', 'publishedAt'"`
	PageSize       int    `json:"page_size,omitempty" description:"Number of results per page (default: 20)"`
	Page           int    `json:"page,omitempty" description:"Page number for pagination (default: 1)"`
	DateFrom       string `json:"date_from,omitempty" description:"Oldest article date (YYYY-MM-DD)"`
	DateTo         string `json:"date_to,omitempty" description:"Newest article date (YYYY-MM-DD)"`
	Domains        string `json:"domains,omitempty" description:"Comma-separated list of domains to include"`
	ExcludeDomains string `json:"exclude_domains,omitempty" description:"Comma-separated list of domains to exclude"`
}

// Tool Results
type SearchNewsAPIResult struct {
	Status       string        `json:"status"`
	TotalResults int           `json:"total_results"`
	Articles     []NewsArticle `json:"articles"`
	FetchedAt    string        `json:"fetched_at"`
}

type NewsArticle struct {
	Source      NewsSource `json:"source"`
	Author      string     `json:"author"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	URL         string     `json:"url"`
	PublishedAt string     `json:"published_at"`
	Content     string     `json:"content"`
}

type NewsSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Schema definitions
var SearchNewsAPIParamSchema = &sdomain.Schema{
	Type:        "object",
	Description: "Parameters for searching news articles via API",
	Properties: map[string]sdomain.Property{
		"query": {
			Type:        "string",
			Description: "Search query for news articles",
		},
		"api_key": {
			Type:        "string",
			Description: "API key for the news service",
		},
		"language": {
			Type:        "string",
			Description: "Language code (e.g., 'en', 'es')",
			Pattern:     "^[a-z]{2}$",
		},
		"sort_by": {
			Type:        "string",
			Description: "Sort order: 'relevancy', 'popularity', 'publishedAt'",
			Enum:        []string{"relevancy", "popularity", "publishedAt"},
		},
		"page_size": {
			Type:        "integer",
			Description: "Number of results per page (default: 20)",
			Minimum:     float64Ptr(1),
			Maximum:     float64Ptr(100),
		},
		"page": {
			Type:        "integer",
			Description: "Page number for pagination (default: 1)",
			Minimum:     float64Ptr(1),
		},
		"date_from": {
			Type:        "string",
			Description: "Oldest article date (YYYY-MM-DD)",
			Pattern:     "^\\d{4}-\\d{2}-\\d{2}$",
		},
		"date_to": {
			Type:        "string",
			Description: "Newest article date (YYYY-MM-DD)",
			Pattern:     "^\\d{4}-\\d{2}-\\d{2}$",
		},
		"domains": {
			Type:        "string",
			Description: "Comma-separated list of domains to include",
		},
		"exclude_domains": {
			Type:        "string",
			Description: "Comma-separated list of domains to exclude",
		},
	},
	Required: []string{"query", "api_key"},
}

// NewSearchNewsAPITool creates a new news API search tool
func NewSearchNewsAPITool() domain.Tool {
	return tools.NewTool(
		"search_news_api",
		"Searches for news articles using a news API service",
		searchNewsAPIHandler,
		SearchNewsAPIParamSchema,
	)
}

func searchNewsAPIHandler(ctx context.Context, params SearchNewsAPIParams) (*SearchNewsAPIResult, error) {
	// TODO: Implement news API search logic
	return nil, fmt.Errorf("not implemented")
}

// Generic REST API call tool could be added here in the future
// type CallRESTAPIParams struct { ... }
// func NewCallRESTAPITool() domain.Tool { ... }
