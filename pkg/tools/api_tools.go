// ABOUTME: API interaction tools for REST, GraphQL, and other API protocols
// ABOUTME: Provides tools to interact with various API services

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

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

// NewsAPI.org endpoint (can be overridden for testing)
var newsAPIURL = "https://newsapi.org/v2/everything"

func searchNewsAPIHandler(ctx context.Context, params SearchNewsAPIParams) (*SearchNewsAPIResult, error) {
	// Validate required parameters
	if strings.TrimSpace(params.Query) == "" {
		return nil, fmt.Errorf("query parameter is required")
	}

	// Use environment variable for API key if not provided
	apiKey := params.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("NEWS_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("API key is required (pass in params or set NEWS_API_KEY environment variable)")
		}
	}

	// Set defaults
	if params.PageSize == 0 {
		params.PageSize = 20
	}
	if params.Page == 0 {
		params.Page = 1
	}
	if params.SortBy == "" {
		params.SortBy = "publishedAt"
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("q", params.Query)
	queryParams.Set("apiKey", apiKey)
	queryParams.Set("pageSize", strconv.Itoa(params.PageSize))
	queryParams.Set("page", strconv.Itoa(params.Page))
	queryParams.Set("sortBy", params.SortBy)

	if params.Language != "" {
		queryParams.Set("language", params.Language)
	}
	if params.DateFrom != "" {
		queryParams.Set("from", params.DateFrom)
	}
	if params.DateTo != "" {
		queryParams.Set("to", params.DateTo)
	}
	if params.Domains != "" {
		queryParams.Set("domains", params.Domains)
	}
	if params.ExcludeDomains != "" {
		queryParams.Set("excludeDomains", params.ExcludeDomains)
	}

	// Build full URL
	fullURL := newsAPIURL + "?" + queryParams.Encode()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "go-flock/1.0 NewsSearcher")
	req.Header.Set("Accept", "application/json")

	// Make request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// Parse response
	var apiResponse newsAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	// Check for API errors
	if apiResponse.Status != "ok" {
		return nil, fmt.Errorf("news API error: %s (code: %s)", apiResponse.Message, apiResponse.Code)
	}

	// Convert to our result format
	result := &SearchNewsAPIResult{
		Status:       apiResponse.Status,
		TotalResults: apiResponse.TotalResults,
		Articles:     make([]NewsArticle, len(apiResponse.Articles)),
		FetchedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	for i, article := range apiResponse.Articles {
		result.Articles[i] = NewsArticle{
			Source: NewsSource{
				ID:   article.Source.ID,
				Name: article.Source.Name,
			},
			Author:      article.Author,
			Title:       article.Title,
			Description: article.Description,
			URL:         article.URL,
			PublishedAt: article.PublishedAt,
			Content:     article.Content,
		}
	}

	return result, nil
}

// Internal structures for NewsAPI response
type newsAPIResponse struct {
	Status       string           `json:"status"`
	TotalResults int              `json:"totalResults"`
	Articles     []newsAPIArticle `json:"articles"`
	Code         string           `json:"code,omitempty"`
	Message      string           `json:"message,omitempty"`
}

type newsAPIArticle struct {
	Source      newsAPISource `json:"source"`
	Author      string        `json:"author"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	URL         string        `json:"url"`
	URLToImage  string        `json:"urlToImage"`
	PublishedAt string        `json:"publishedAt"`
	Content     string        `json:"content"`
}

type newsAPISource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Generic REST API call tool could be added here in the future
// type CallRESTAPIParams struct { ... }
// func NewCallRESTAPITool() domain.Tool { ... }
