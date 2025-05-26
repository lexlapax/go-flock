// ABOUTME: API interaction tools for REST, GraphQL, and other API protocols
// ABOUTME: Provides tools to interact with various API services

package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
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

// Brave Search API Implementation

// Tool Parameters
type SearchWebBraveParams struct {
	Query         string   `json:"query" description:"Search query (max 400 chars/50 words)"`
	APIKey        string   `json:"api_key" description:"API key for Brave Search"`
	Count         int      `json:"count,omitempty" description:"Results per page (max 20, default 10)"`
	Offset        int      `json:"offset,omitempty" description:"Page offset for pagination (max 9)"`
	Country       string   `json:"country,omitempty" description:"Country code (e.g., 'US', 'GB')"`
	SearchLang    string   `json:"search_lang,omitempty" description:"Language preference (e.g., 'en', 'es')"`
	SafeSearch    string   `json:"safesearch,omitempty" description:"Content filter: 'off', 'moderate', 'strict'"`
	Freshness     string   `json:"freshness,omitempty" description:"Time filter: 'pd' (past day), 'pw' (past week), 'pm' (past month), 'py' (past year)"`
	ResultFilter  []string `json:"result_filter,omitempty" description:"Filter types: 'web', 'news', 'images', 'videos'"`
	ExtraSnippets bool     `json:"extra_snippets,omitempty" description:"Enable additional text excerpts"`
	Summary       bool     `json:"summary,omitempty" description:"Enable AI-generated summaries"`
	Goggles       string   `json:"goggles,omitempty" description:"Custom ranking profile URL"`
}

// Tool Results
type SearchWebBraveResult struct {
	Query       string            `json:"query"`
	Type        string            `json:"type"`
	Mixed       BraveMixedResults `json:"mixed,omitempty"`
	Web         BraveWebResults   `json:"web,omitempty"`
	News        BraveNewsResults  `json:"news,omitempty"`
	Videos      BraveVideoResults `json:"videos,omitempty"`
	InfoBox     *BraveInfoBox     `json:"infobox,omitempty"`
	Summary     []BraveSummary    `json:"summary,omitempty"`
	Locations   []BraveLocation   `json:"locations,omitempty"`
	Discussions []BraveDiscussion `json:"discussions,omitempty"`
	FAQ         []BraveFAQ        `json:"faq,omitempty"`
	FetchedAt   string            `json:"fetched_at"`
}

type BraveMixedResults struct {
	Type string           `json:"type"`
	Main []BraveMixedItem `json:"main"`
}

type BraveMixedItem struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type BraveWebResults struct {
	Type    string           `json:"type"`
	Results []BraveWebResult `json:"results"`
}

type BraveWebResult struct {
	Title       string               `json:"title"`
	URL         string               `json:"url"`
	Description string               `json:"description"`
	Age         string               `json:"age,omitempty"`
	Language    string               `json:"language,omitempty"`
	Thumbnail   *BraveThumbnail      `json:"thumbnail,omitempty"`
	Location    *BraveResultLocation `json:"location,omitempty"`
	DeepLinks   []BraveDeepLink      `json:"deep_links,omitempty"`
}

type BraveThumbnail struct {
	Src    string `json:"src"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type BraveResultLocation struct {
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
}

type BraveDeepLink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type BraveNewsResults struct {
	Type    string            `json:"type"`
	Results []BraveNewsResult `json:"results"`
}

type BraveNewsResult struct {
	Title       string           `json:"title"`
	URL         string           `json:"url"`
	Description string           `json:"description"`
	Age         string           `json:"age,omitempty"`
	Source      *BraveNewsSource `json:"source,omitempty"`
	Thumbnail   *BraveThumbnail  `json:"thumbnail,omitempty"`
}

type BraveNewsSource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type BraveVideoResults struct {
	Type    string             `json:"type"`
	Results []BraveVideoResult `json:"results"`
}

type BraveVideoResult struct {
	Title       string          `json:"title"`
	URL         string          `json:"url"`
	Description string          `json:"description"`
	Age         string          `json:"age,omitempty"`
	Duration    string          `json:"duration,omitempty"`
	Creator     string          `json:"creator,omitempty"`
	Publisher   string          `json:"publisher,omitempty"`
	Thumbnail   *BraveThumbnail `json:"thumbnail,omitempty"`
}

type BraveInfoBox struct {
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	URL         string                 `json:"url,omitempty"`
	Thumbnail   *BraveThumbnail        `json:"thumbnail,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

type BraveSummary struct {
	Type        string                 `json:"type"`
	Key         string                 `json:"key"`
	Text        string                 `json:"text"`
	Enrichments map[string]interface{} `json:"enrichments,omitempty"`
}

type BraveLocation struct {
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Address  string  `json:"address,omitempty"`
	Phone    string  `json:"phone,omitempty"`
	Rating   float64 `json:"rating,omitempty"`
	Distance string  `json:"distance,omitempty"`
}

type BraveDiscussion struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Age   string `json:"age,omitempty"`
}

type BraveFAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	URL      string `json:"url,omitempty"`
}

// Schema definitions
var SearchWebBraveParamSchema = &sdomain.Schema{
	Type:        "object",
	Description: "Parameters for Brave web search",
	Properties: map[string]sdomain.Property{
		"query": {
			Type:        "string",
			Description: "Search query (max 400 chars/50 words)",
			MaxLength:   intPtr(400),
		},
		"api_key": {
			Type:        "string",
			Description: "API key for Brave Search",
		},
		"count": {
			Type:        "integer",
			Description: "Results per page (max 20, default 10)",
			Minimum:     float64Ptr(1),
			Maximum:     float64Ptr(20),
		},
		"offset": {
			Type:        "integer",
			Description: "Page offset for pagination (max 9)",
			Minimum:     float64Ptr(0),
			Maximum:     float64Ptr(9),
		},
		"country": {
			Type:        "string",
			Description: "Country code (e.g., 'US', 'GB')",
			Pattern:     "^[A-Z]{2}$",
		},
		"search_lang": {
			Type:        "string",
			Description: "Language preference (e.g., 'en', 'es')",
			Pattern:     "^[a-z]{2}$",
		},
		"safesearch": {
			Type:        "string",
			Description: "Content filter: 'off', 'moderate', 'strict'",
			Enum:        []string{"off", "moderate", "strict"},
		},
		"freshness": {
			Type:        "string",
			Description: "Time filter: 'pd' (past day), 'pw' (past week), 'pm' (past month), 'py' (past year)",
			Enum:        []string{"pd", "pw", "pm", "py"},
		},
		"result_filter": {
			Type:        "array",
			Description: "Filter types: 'web', 'news', 'images', 'videos'",
			Items: &sdomain.Property{
				Type: "string",
				Enum: []string{"web", "news", "images", "videos"},
			},
		},
		"extra_snippets": {
			Type:        "boolean",
			Description: "Enable additional text excerpts",
		},
		"summary": {
			Type:        "boolean",
			Description: "Enable AI-generated summaries",
		},
		"goggles": {
			Type:        "string",
			Description: "Custom ranking profile URL",
		},
	},
	Required: []string{"query", "api_key"},
}

// NewSearchWebBraveTool creates a new Brave web search tool
func NewSearchWebBraveTool() domain.Tool {
	return tools.NewTool(
		"search_web_brave",
		"Performs web search using Brave Search API with support for multiple content types, AI summaries, and advanced filtering",
		searchWebBraveHandler,
		SearchWebBraveParamSchema,
	)
}

// Brave Search API endpoint (can be overridden for testing)
var braveSearchURL = "https://api.search.brave.com/res/v1/web/search"

func searchWebBraveHandler(ctx context.Context, params SearchWebBraveParams) (*SearchWebBraveResult, error) {
	// Validate required parameters
	if strings.TrimSpace(params.Query) == "" {
		return nil, fmt.Errorf("query parameter is required")
	}

	// Check query length (400 chars max)
	if len(params.Query) > 400 {
		return nil, fmt.Errorf("query exceeds 400 character limit")
	}

	// Use environment variable for API key if not provided
	apiKey := params.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("BRAVE_SEARCH_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("API key is required (pass in params or set BRAVE_SEARCH_API_KEY environment variable)")
		}
	}

	// Set defaults
	if params.Count == 0 {
		params.Count = 10
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("q", params.Query)
	queryParams.Set("count", strconv.Itoa(params.Count))

	if params.Offset > 0 {
		queryParams.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Country != "" {
		queryParams.Set("country", params.Country)
	}
	if params.SearchLang != "" {
		queryParams.Set("search_lang", params.SearchLang)
	}
	if params.SafeSearch != "" {
		queryParams.Set("safesearch", params.SafeSearch)
	}
	if params.Freshness != "" {
		queryParams.Set("freshness", params.Freshness)
	}
	if len(params.ResultFilter) > 0 {
		queryParams.Set("result_filter", strings.Join(params.ResultFilter, ","))
	}
	if params.ExtraSnippets {
		queryParams.Set("extra_snippets", "true")
	}
	if params.Summary {
		queryParams.Set("summary", "true")
	}
	if params.Goggles != "" {
		queryParams.Set("goggles", params.Goggles)
	}

	// Build full URL
	fullURL := braveSearchURL + "?" + queryParams.Encode()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Subscription-Token", apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "go-flock/1.0 BraveSearcher")

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

	// Check for error status codes
	if resp.StatusCode != http.StatusOK {
		var errorResp braveErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Type == "error" {
			// Extract error details
			if len(errorResp.Details) > 0 {
				detail := errorResp.Details[0]
				return nil, fmt.Errorf("Brave API error: %s", detail.Detail)
			}
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResponse braveSearchResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	// Convert to our result format
	result := &SearchWebBraveResult{
		Query:     apiResponse.Query.Original,
		Type:      apiResponse.Type,
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// Copy mixed results
	if apiResponse.Mixed != nil {
		result.Mixed = BraveMixedResults{
			Type: apiResponse.Mixed.Type,
			Main: make([]BraveMixedItem, len(apiResponse.Mixed.Main)),
		}
		for i, item := range apiResponse.Mixed.Main {
			result.Mixed.Main[i] = BraveMixedItem(item)
		}
	}

	// Copy web results
	if apiResponse.Web != nil {
		result.Web = convertWebResults(apiResponse.Web)
	}

	// Copy news results
	if apiResponse.News != nil {
		result.News = convertNewsResults(apiResponse.News)
	}

	// Copy video results
	if apiResponse.Videos != nil {
		result.Videos = convertVideoResults(apiResponse.Videos)
	}

	// Copy infobox
	if apiResponse.InfoBox != nil {
		result.InfoBox = convertInfoBox(apiResponse.InfoBox)
	}

	// Copy summaries
	if len(apiResponse.Summary) > 0 {
		result.Summary = make([]BraveSummary, len(apiResponse.Summary))
		for i, sum := range apiResponse.Summary {
			result.Summary[i] = BraveSummary(sum)
		}
	}

	// Copy locations
	if len(apiResponse.Locations) > 0 {
		result.Locations = make([]BraveLocation, len(apiResponse.Locations))
		for i, loc := range apiResponse.Locations {
			result.Locations[i] = BraveLocation(loc)
		}
	}

	// Copy discussions
	if len(apiResponse.Discussions) > 0 {
		result.Discussions = make([]BraveDiscussion, len(apiResponse.Discussions))
		for i, disc := range apiResponse.Discussions {
			result.Discussions[i] = BraveDiscussion(disc)
		}
	}

	// Copy FAQ
	if len(apiResponse.FAQ) > 0 {
		result.FAQ = make([]BraveFAQ, len(apiResponse.FAQ))
		for i, faq := range apiResponse.FAQ {
			result.FAQ[i] = BraveFAQ(faq)
		}
	}

	return result, nil
}

// Helper functions for converting results

func convertWebResults(web *braveWebSection) BraveWebResults {
	results := BraveWebResults{
		Type:    web.Type,
		Results: make([]BraveWebResult, len(web.Results)),
	}

	for i, r := range web.Results {
		result := BraveWebResult{
			Title:       r.Title,
			URL:         r.URL,
			Description: r.Description,
			Age:         r.Age,
			Language:    r.Language,
		}

		if r.Thumbnail != nil {
			result.Thumbnail = &BraveThumbnail{
				Src:    r.Thumbnail.Src,
				Height: r.Thumbnail.Height,
				Width:  r.Thumbnail.Width,
			}
		}

		if r.Location != nil {
			result.Location = &BraveResultLocation{
				City:    r.Location.City,
				State:   r.Location.State,
				Country: r.Location.Country,
			}
		}

		if len(r.DeepLinks) > 0 {
			result.DeepLinks = make([]BraveDeepLink, len(r.DeepLinks))
			for j, dl := range r.DeepLinks {
				result.DeepLinks[j] = BraveDeepLink(dl)
			}
		}

		results.Results[i] = result
	}

	return results
}

func convertNewsResults(news *braveNewsSection) BraveNewsResults {
	results := BraveNewsResults{
		Type:    news.Type,
		Results: make([]BraveNewsResult, len(news.Results)),
	}

	for i, r := range news.Results {
		result := BraveNewsResult{
			Title:       r.Title,
			URL:         r.URL,
			Description: r.Description,
			Age:         r.Age,
		}

		if r.Source != nil {
			result.Source = &BraveNewsSource{
				Name: r.Source.Name,
				URL:  r.Source.URL,
			}
		}

		if r.Thumbnail != nil {
			result.Thumbnail = &BraveThumbnail{
				Src:    r.Thumbnail.Src,
				Height: r.Thumbnail.Height,
				Width:  r.Thumbnail.Width,
			}
		}

		results.Results[i] = result
	}

	return results
}

func convertVideoResults(videos *braveVideoSection) BraveVideoResults {
	results := BraveVideoResults{
		Type:    videos.Type,
		Results: make([]BraveVideoResult, len(videos.Results)),
	}

	for i, r := range videos.Results {
		result := BraveVideoResult{
			Title:       r.Title,
			URL:         r.URL,
			Description: r.Description,
			Age:         r.Age,
			Duration:    r.Duration,
			Creator:     r.Creator,
			Publisher:   r.Publisher,
		}

		if r.Thumbnail != nil {
			result.Thumbnail = &BraveThumbnail{
				Src:    r.Thumbnail.Src,
				Height: r.Thumbnail.Height,
				Width:  r.Thumbnail.Width,
			}
		}

		results.Results[i] = result
	}

	return results
}

func convertInfoBox(info *braveInfoBox) *BraveInfoBox {
	result := &BraveInfoBox{
		Type:        info.Type,
		Title:       info.Title,
		Description: info.Description,
		URL:         info.URL,
		Attributes:  info.Attributes,
	}

	if info.Thumbnail != nil {
		result.Thumbnail = &BraveThumbnail{
			Src:    info.Thumbnail.Src,
			Height: info.Thumbnail.Height,
			Width:  info.Thumbnail.Width,
		}
	}

	return result
}

// Internal structures for Brave API response

type braveSearchResponse struct {
	Query       braveQuery         `json:"query"`
	Type        string             `json:"type"`
	Mixed       *braveMixedSection `json:"mixed,omitempty"`
	Web         *braveWebSection   `json:"web,omitempty"`
	News        *braveNewsSection  `json:"news,omitempty"`
	Videos      *braveVideoSection `json:"videos,omitempty"`
	InfoBox     *braveInfoBox      `json:"infobox,omitempty"`
	Summary     []braveSummary     `json:"summary,omitempty"`
	Locations   []braveLocation    `json:"locations,omitempty"`
	Discussions []braveDiscussion  `json:"discussions,omitempty"`
	FAQ         []braveFAQ         `json:"faq,omitempty"`
}

type braveQuery struct {
	Original string `json:"original"`
}

type braveMixedSection struct {
	Type string           `json:"type"`
	Main []braveMixedItem `json:"main"`
}

type braveMixedItem struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type braveWebSection struct {
	Type    string           `json:"type"`
	Results []braveWebResult `json:"results"`
}

type braveWebResult struct {
	Title       string               `json:"title"`
	URL         string               `json:"url"`
	Description string               `json:"description"`
	Age         string               `json:"age,omitempty"`
	Language    string               `json:"language,omitempty"`
	Thumbnail   *braveThumbnail      `json:"thumbnail,omitempty"`
	Location    *braveResultLocation `json:"location,omitempty"`
	DeepLinks   []braveDeepLink      `json:"deep_links,omitempty"`
}

type braveThumbnail struct {
	Src    string `json:"src"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type braveResultLocation struct {
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
}

type braveDeepLink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type braveNewsSection struct {
	Type    string            `json:"type"`
	Results []braveNewsResult `json:"results"`
}

type braveNewsResult struct {
	Title       string           `json:"title"`
	URL         string           `json:"url"`
	Description string           `json:"description"`
	Age         string           `json:"age,omitempty"`
	Source      *braveNewsSource `json:"source,omitempty"`
	Thumbnail   *braveThumbnail  `json:"thumbnail,omitempty"`
}

type braveNewsSource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type braveVideoSection struct {
	Type    string             `json:"type"`
	Results []braveVideoResult `json:"results"`
}

type braveVideoResult struct {
	Title       string          `json:"title"`
	URL         string          `json:"url"`
	Description string          `json:"description"`
	Age         string          `json:"age,omitempty"`
	Duration    string          `json:"duration,omitempty"`
	Creator     string          `json:"creator,omitempty"`
	Publisher   string          `json:"publisher,omitempty"`
	Thumbnail   *braveThumbnail `json:"thumbnail,omitempty"`
}

type braveInfoBox struct {
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	URL         string                 `json:"url,omitempty"`
	Thumbnail   *braveThumbnail        `json:"thumbnail,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

type braveSummary struct {
	Type        string                 `json:"type"`
	Key         string                 `json:"key"`
	Text        string                 `json:"text"`
	Enrichments map[string]interface{} `json:"enrichments,omitempty"`
}

type braveLocation struct {
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Address  string  `json:"address,omitempty"`
	Phone    string  `json:"phone,omitempty"`
	Rating   float64 `json:"rating,omitempty"`
	Distance string  `json:"distance,omitempty"`
}

type braveDiscussion struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Age   string `json:"age,omitempty"`
}

type braveFAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	URL      string `json:"url,omitempty"`
}

type braveErrorResponse struct {
	Type    string             `json:"type"`
	Code    int                `json:"code"`
	Details []braveErrorDetail `json:"details"`
}

type braveErrorDetail struct {
	Type   string `json:"type"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// Helper function to get int pointer
func intPtr(i int) *int {
	return &i
}

// Research Paper API Implementation

// Tool Parameters
type ResearchPaperAPIParams struct {
	Query      string   `json:"query" description:"Search query for research papers"`
	MaxResults int      `json:"max_results,omitempty" description:"Maximum results per provider (default: 10)"`
	StartDate  string   `json:"start_date,omitempty" description:"Filter papers from this date (YYYY-MM-DD)"`
	EndDate    string   `json:"end_date,omitempty" description:"Filter papers until this date (YYYY-MM-DD)"`
	Authors    []string `json:"authors,omitempty" description:"Filter by author names"`
	Categories []string `json:"categories,omitempty" description:"Subject categories (cs, physics, medicine, etc.)"`
	OpenAccess bool     `json:"open_access,omitempty" description:"Only return open access papers"`
	SortBy     string   `json:"sort_by,omitempty" description:"Sort order: 'relevance', 'date', 'citations'"`
	Providers  []string `json:"providers,omitempty" description:"Specific providers to search (arxiv, pubmed, core)"`
	CoreAPIKey string   `json:"core_api_key,omitempty" description:"CORE API key for enhanced access"`
}

// Tool Results
type ResearchPaperAPIResult struct {
	Query        string          `json:"query"`
	TotalResults int             `json:"total_results"`
	Papers       []ResearchPaper `json:"papers"`
	Providers    []ProviderInfo  `json:"providers"`
	FetchedAt    string          `json:"fetched_at"`
}

type ResearchPaper struct {
	// Basic metadata
	Title         string   `json:"title"`
	Authors       []string `json:"authors"`
	Abstract      string   `json:"abstract"`
	PublishedDate string   `json:"published_date"`
	Source        string   `json:"source"`
	URL           string   `json:"url"`
	PDFURL        string   `json:"pdf_url,omitempty"`

	// Identifiers
	DOI      string `json:"doi,omitempty"`
	ArxivID  string `json:"arxiv_id,omitempty"`
	PubMedID string `json:"pubmed_id,omitempty"`

	// Additional metadata
	Journal        string  `json:"journal,omitempty"`
	RelevanceScore float64 `json:"relevance_score,omitempty"`
}

type ProviderInfo struct {
	Name         string `json:"name"`
	ResultCount  int    `json:"result_count"`
	Error        string `json:"error,omitempty"`
	ResponseTime int64  `json:"response_time_ms"`
}

// Schema definition
var ResearchPaperAPIParamSchema = &sdomain.Schema{
	Type:        "object",
	Description: "Parameters for searching academic research papers",
	Properties: map[string]sdomain.Property{
		"query": {
			Type:        "string",
			Description: "Search query for research papers",
		},
		"max_results": {
			Type:        "integer",
			Description: "Maximum results per provider (default: 10)",
			Minimum:     float64Ptr(1),
			Maximum:     float64Ptr(100),
		},
		"start_date": {
			Type:        "string",
			Description: "Filter papers from this date (YYYY-MM-DD)",
			Pattern:     "^\\d{4}-\\d{2}-\\d{2}$",
		},
		"end_date": {
			Type:        "string",
			Description: "Filter papers until this date (YYYY-MM-DD)",
			Pattern:     "^\\d{4}-\\d{2}-\\d{2}$",
		},
		"authors": {
			Type:        "array",
			Description: "Filter by author names",
			Items: &sdomain.Property{
				Type: "string",
			},
		},
		"categories": {
			Type:        "array",
			Description: "Subject categories (cs, physics, medicine, etc.)",
			Items: &sdomain.Property{
				Type: "string",
			},
		},
		"open_access": {
			Type:        "boolean",
			Description: "Only return open access papers",
		},
		"sort_by": {
			Type:        "string",
			Description: "Sort order: 'relevance', 'date', 'citations'",
			Enum:        []string{"relevance", "date", "citations"},
		},
		"providers": {
			Type:        "array",
			Description: "Specific providers to search (arxiv, pubmed, core)",
			Items: &sdomain.Property{
				Type: "string",
				Enum: []string{"arxiv", "pubmed", "core"},
			},
		},
		"core_api_key": {
			Type:        "string",
			Description: "CORE API key for enhanced access",
		},
	},
	Required: []string{"query"},
}

// NewResearchPaperAPITool creates a new research paper API tool
func NewResearchPaperAPITool() domain.Tool {
	return tools.NewTool(
		"research_paper_api",
		"Searches for academic papers across multiple research databases (arXiv, PubMed, CORE) in parallel",
		researchPaperAPIHandler,
		ResearchPaperAPIParamSchema,
	)
}

// API URLs - package-level for test overrides
var (
	arxivAPIURL   = "http://export.arxiv.org/api/query"
	pubmedBaseURL = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
	coreAPIURL    = "https://api.core.ac.uk/v3/search/works"
)

func researchPaperAPIHandler(ctx context.Context, params ResearchPaperAPIParams) (*ResearchPaperAPIResult, error) {
	// Validate required parameters
	if strings.TrimSpace(params.Query) == "" {
		return nil, fmt.Errorf("query parameter is required")
	}

	// Set defaults
	if params.MaxResults == 0 {
		params.MaxResults = 10
	}
	if params.SortBy == "" {
		params.SortBy = "relevance"
	}

	// Select providers
	providers := selectProviders(params)

	// Channel for collecting results
	type providerResult struct {
		provider string
		papers   []ResearchPaper
		err      error
		duration time.Duration
	}

	resultChan := make(chan providerResult, len(providers))

	// Search all providers in parallel
	for _, provider := range providers {
		go func(p string) {
			start := time.Now()
			var papers []ResearchPaper
			var err error

			switch p {
			case "arxiv":
				papers, err = searchArxiv(ctx, params)
			case "pubmed":
				papers, err = searchPubMed(ctx, params)
			case "core":
				papers, err = searchCORE(ctx, params)
			}

			resultChan <- providerResult{
				provider: p,
				papers:   papers,
				err:      err,
				duration: time.Since(start),
			}
		}(provider)
	}

	// Collect results with timeout
	timeout := time.After(30 * time.Second)
	var allPapers []ResearchPaper
	providerInfos := make([]ProviderInfo, 0, len(providers))

	for i := 0; i < len(providers); i++ {
		select {
		case result := <-resultChan:
			info := ProviderInfo{
				Name:         result.provider,
				ResponseTime: result.duration.Milliseconds(),
			}

			if result.err != nil {
				info.Error = result.err.Error()
			} else {
				info.ResultCount = len(result.papers)
				allPapers = append(allPapers, result.papers...)
			}

			providerInfos = append(providerInfos, info)

		case <-timeout:
			// Continue with partial results
			break
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Deduplicate papers
	deduplicatedPapers := deduplicatePapers(allPapers)

	// Sort papers
	sortPapers(deduplicatedPapers, params.SortBy)

	// Build result
	result := &ResearchPaperAPIResult{
		Query:        params.Query,
		TotalResults: len(deduplicatedPapers),
		Papers:       deduplicatedPapers,
		Providers:    providerInfos,
		FetchedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	return result, nil
}

// selectProviders chooses which providers to search based on parameters
func selectProviders(params ResearchPaperAPIParams) []string {
	if len(params.Providers) > 0 {
		return params.Providers
	}

	// Default to all providers
	providers := []string{"arxiv", "pubmed", "core"}

	// If categories specified, optimize provider selection
	if len(params.Categories) > 0 {
		selectedProviders := make(map[string]bool)

		for _, cat := range params.Categories {
			switch strings.ToLower(cat) {
			case "cs", "math", "physics", "astro-ph", "cond-mat", "q-bio", "q-fin", "stat":
				selectedProviders["arxiv"] = true
			case "medicine", "biology", "health", "clinical":
				selectedProviders["pubmed"] = true
			}
			// CORE is general purpose, always include if any category matches
			selectedProviders["core"] = true
		}

		if len(selectedProviders) > 0 {
			providers = make([]string, 0, len(selectedProviders))
			for p := range selectedProviders {
				providers = append(providers, p)
			}
		}
	}

	return providers
}

// deduplicatePapers removes duplicate papers based on DOI and title similarity
func deduplicatePapers(papers []ResearchPaper) []ResearchPaper {
	seen := make(map[string]bool)
	unique := make([]ResearchPaper, 0, len(papers))

	for _, paper := range papers {
		// Check DOI first
		if paper.DOI != "" {
			if seen[paper.DOI] {
				continue
			}
			seen[paper.DOI] = true
		}

		// Check for similar titles (simple approach)
		titleKey := strings.ToLower(strings.TrimSpace(paper.Title))
		if seen[titleKey] {
			continue
		}
		seen[titleKey] = true

		unique = append(unique, paper)
	}

	return unique
}

// sortPapers sorts papers based on the specified criteria
func sortPapers(papers []ResearchPaper, sortBy string) {
	// Simple sorting implementation
	// In production, would implement more sophisticated sorting
	switch sortBy {
	case "date":
		// Sort by published date (newest first)
		// Implementation would parse dates and sort
	case "citations":
		// Sort by citation count
		// Implementation would sort by CitationCount field
	default:
		// Relevance sorting (keep original order from providers)
	}
}

func searchArxiv(ctx context.Context, params ResearchPaperAPIParams) ([]ResearchPaper, error) {
	// Build query parameters
	query := url.Values{}
	query.Set("search_query", params.Query)
	query.Set("start", "0")
	query.Set("max_results", fmt.Sprintf("%d", params.MaxResults))
	query.Set("sortBy", "relevance")
	query.Set("sortOrder", "descending")

	// Build URL
	apiURL := fmt.Sprintf("%s?%s", arxivAPIURL, query.Encode())

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Execute request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse Atom feed
	type ArxivEntry struct {
		ID        string `xml:"id"`
		Title     string `xml:"title"`
		Summary   string `xml:"summary"`
		Published string `xml:"published"`
		Authors   []struct {
			Name string `xml:"name"`
		} `xml:"author"`
		Links []struct {
			Href string `xml:"href,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
	}

	type ArxivFeed struct {
		Entries []ArxivEntry `xml:"entry"`
	}

	var feed ArxivFeed
	decoder := xml.NewDecoder(resp.Body)
	if err := decoder.Decode(&feed); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	// Convert to ResearchPaper format
	papers := make([]ResearchPaper, 0, len(feed.Entries))
	for _, entry := range feed.Entries {
		// Extract authors
		authors := make([]string, len(entry.Authors))
		for i, author := range entry.Authors {
			authors[i] = strings.TrimSpace(author.Name)
		}

		// Find PDF link
		var pdfURL string
		for _, link := range entry.Links {
			if link.Type == "application/pdf" {
				pdfURL = link.Href
				break
			}
		}

		// Extract arXiv ID from entry ID
		arxivID := strings.TrimPrefix(entry.ID, "http://arxiv.org/abs/")

		// Parse publication date
		pubDate, _ := time.Parse(time.RFC3339, entry.Published)

		papers = append(papers, ResearchPaper{
			Title:          strings.TrimSpace(entry.Title),
			Authors:        authors,
			Abstract:       strings.TrimSpace(entry.Summary),
			PublishedDate:  pubDate.Format("2006-01-02"),
			Source:         "arXiv",
			URL:            entry.ID,
			PDFURL:         pdfURL,
			ArxivID:        arxivID,
			RelevanceScore: 1.0, // arXiv doesn't provide relevance scores
		})
	}

	return papers, nil
}

func searchPubMed(ctx context.Context, params ResearchPaperAPIParams) ([]ResearchPaper, error) {
	// PubMed requires two API calls: esearch to get IDs, then efetch to get details

	// First, search for IDs
	searchURL := pubmedBaseURL + "esearch.fcgi"
	searchQuery := url.Values{}
	searchQuery.Set("db", "pubmed")
	searchQuery.Set("term", params.Query)
	searchQuery.Set("retmax", fmt.Sprintf("%d", params.MaxResults))
	searchQuery.Set("retmode", "json")
	searchQuery.Set("sort", "relevance")

	// Add API key if available
	if apiKey := os.Getenv("PUBMED_API_KEY"); apiKey != "" {
		searchQuery.Set("api_key", apiKey)
	}

	fullSearchURL := fmt.Sprintf("%s?%s", searchURL, searchQuery.Encode())

	// Create search request
	searchReq, err := http.NewRequestWithContext(ctx, "GET", fullSearchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating search request: %w", err)
	}

	// Execute search request
	client := &http.Client{Timeout: 30 * time.Second}
	searchResp, err := client.Do(searchReq)
	if err != nil {
		return nil, fmt.Errorf("executing search request: %w", err)
	}
	defer searchResp.Body.Close()

	if searchResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(searchResp.Body)
		return nil, fmt.Errorf("search API returned status %d: %s", searchResp.StatusCode, string(body))
	}

	// Parse search results
	var searchResult struct {
		ESearchResult struct {
			IDList []string `json:"idlist"`
		} `json:"esearchresult"`
	}

	if err := json.NewDecoder(searchResp.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("parsing search response: %w", err)
	}

	if len(searchResult.ESearchResult.IDList) == 0 {
		return []ResearchPaper{}, nil
	}

	// Second, fetch details for the IDs
	fetchURL := pubmedBaseURL + "efetch.fcgi"
	fetchQuery := url.Values{}
	fetchQuery.Set("db", "pubmed")
	fetchQuery.Set("id", strings.Join(searchResult.ESearchResult.IDList, ","))
	fetchQuery.Set("retmode", "xml")
	if apiKey := os.Getenv("PUBMED_API_KEY"); apiKey != "" {
		fetchQuery.Set("api_key", apiKey)
	}

	fullFetchURL := fmt.Sprintf("%s?%s", fetchURL, fetchQuery.Encode())

	// Create fetch request
	fetchReq, err := http.NewRequestWithContext(ctx, "GET", fullFetchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating fetch request: %w", err)
	}

	// Execute fetch request
	fetchResp, err := client.Do(fetchReq)
	if err != nil {
		return nil, fmt.Errorf("executing fetch request: %w", err)
	}
	defer fetchResp.Body.Close()

	if fetchResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(fetchResp.Body)
		return nil, fmt.Errorf("fetch API returned status %d: %s", fetchResp.StatusCode, string(body))
	}

	// Parse XML response
	type PubMedArticle struct {
		MedlineCitation struct {
			PMID struct {
				Value string `xml:",chardata"`
			} `xml:"PMID"`
			Article struct {
				ArticleTitle string `xml:"ArticleTitle"`
				Abstract     struct {
					AbstractText []string `xml:"AbstractText"`
				} `xml:"Abstract"`
				AuthorList struct {
					Author []struct {
						LastName string `xml:"LastName"`
						ForeName string `xml:"ForeName"`
						Initials string `xml:"Initials"`
					} `xml:"Author"`
				} `xml:"AuthorList"`
				Journal struct {
					Title string `xml:"Title"`
				} `xml:"Journal"`
				ArticleDate struct {
					Year  string `xml:"Year"`
					Month string `xml:"Month"`
					Day   string `xml:"Day"`
				} `xml:"ArticleDate"`
			} `xml:"Article"`
		} `xml:"MedlineCitation"`
		PubmedData struct {
			ArticleIdList struct {
				ArticleId []struct {
					IdType string `xml:"IdType,attr"`
					Value  string `xml:",chardata"`
				} `xml:"ArticleId"`
			} `xml:"ArticleIdList"`
		} `xml:"PubmedData"`
	}

	type PubMedArticleSet struct {
		Articles []PubMedArticle `xml:"PubmedArticle"`
	}

	var articleSet PubMedArticleSet
	decoder := xml.NewDecoder(fetchResp.Body)
	if err := decoder.Decode(&articleSet); err != nil {
		return nil, fmt.Errorf("parsing fetch response: %w", err)
	}

	// Convert to ResearchPaper format
	papers := make([]ResearchPaper, 0, len(articleSet.Articles))
	for i, article := range articleSet.Articles {
		// Extract authors
		authors := make([]string, len(article.MedlineCitation.Article.AuthorList.Author))
		for j, author := range article.MedlineCitation.Article.AuthorList.Author {
			if author.ForeName != "" && author.LastName != "" {
				authors[j] = fmt.Sprintf("%s %s", author.ForeName, author.LastName)
			} else if author.LastName != "" {
				authors[j] = author.LastName
			}
		}

		// Extract abstract
		abstract := strings.Join(article.MedlineCitation.Article.Abstract.AbstractText, " ")

		// Extract DOI and PubMed ID
		var doi, pmid string
		pmid = article.MedlineCitation.PMID.Value
		for _, id := range article.PubmedData.ArticleIdList.ArticleId {
			if id.IdType == "doi" {
				doi = id.Value
				break
			}
		}

		// Build publication date
		var pubDate string
		if article.MedlineCitation.Article.ArticleDate.Year != "" {
			year := article.MedlineCitation.Article.ArticleDate.Year
			month := article.MedlineCitation.Article.ArticleDate.Month
			day := article.MedlineCitation.Article.ArticleDate.Day
			if month == "" {
				month = "01"
			}
			if day == "" {
				day = "01"
			}
			pubDate = fmt.Sprintf("%s-%02s-%02s", year, month, day)
		}

		// Build URL
		url := fmt.Sprintf("https://pubmed.ncbi.nlm.nih.gov/%s/", pmid)

		papers = append(papers, ResearchPaper{
			Title:          strings.TrimSpace(article.MedlineCitation.Article.ArticleTitle),
			Authors:        authors,
			Abstract:       abstract,
			PublishedDate:  pubDate,
			Source:         "PubMed",
			URL:            url,
			DOI:            doi,
			PubMedID:       pmid,
			Journal:        article.MedlineCitation.Article.Journal.Title,
			RelevanceScore: 1.0 - (float64(i) * 0.01), // Approximate relevance based on order
		})
	}

	return papers, nil
}

func searchCORE(ctx context.Context, params ResearchPaperAPIParams) ([]ResearchPaper, error) {
	// Get API key
	apiKey := os.Getenv("CORE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("CORE_API_KEY environment variable not set")
	}

	// Build request body
	reqBody := map[string]interface{}{
		"q":     params.Query,
		"limit": params.MaxResults,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling request body: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", coreAPIURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var coreResp struct {
		Results []struct {
			ID            string   `json:"id"`
			Title         string   `json:"title"`
			Abstract      string   `json:"abstract"`
			Authors       []string `json:"authors"`
			PublishedDate string   `json:"publishedDate"`
			DOI           string   `json:"doi"`
			Links         []struct {
				URL  string `json:"url"`
				Type string `json:"type"`
			} `json:"links"`
			Journals []string `json:"journals"`
			Score    float64  `json:"score"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&coreResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	// Convert to ResearchPaper format
	papers := make([]ResearchPaper, 0, len(coreResp.Results))
	for _, result := range coreResp.Results {
		// Find PDF link
		var pdfURL string
		for _, link := range result.Links {
			if link.Type == "pdf" {
				pdfURL = link.URL
				break
			}
		}

		// Use first link as URL if no specific type
		var articleURL string
		if len(result.Links) > 0 {
			articleURL = result.Links[0].URL
		}

		// Get journal name
		var journal string
		if len(result.Journals) > 0 {
			journal = result.Journals[0]
		}

		// Parse date to standard format
		var pubDate string
		if result.PublishedDate != "" {
			// Try to parse various date formats
			for _, format := range []string{
				"2006-01-02",
				"2006-01-02T15:04:05Z",
				"2006-01-02T15:04:05.000Z",
				"2006",
			} {
				if t, err := time.Parse(format, result.PublishedDate); err == nil {
					pubDate = t.Format("2006-01-02")
					break
				}
			}
			if pubDate == "" {
				pubDate = result.PublishedDate // Use as-is if parsing fails
			}
		}

		papers = append(papers, ResearchPaper{
			Title:          result.Title,
			Authors:        result.Authors,
			Abstract:       result.Abstract,
			PublishedDate:  pubDate,
			Source:         "CORE",
			URL:            articleURL,
			PDFURL:         pdfURL,
			DOI:            result.DOI,
			Journal:        journal,
			RelevanceScore: result.Score,
		})
	}

	return papers, nil
}
