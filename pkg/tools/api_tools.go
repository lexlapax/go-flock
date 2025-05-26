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
