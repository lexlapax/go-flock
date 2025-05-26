// ABOUTME: Unit tests for API interaction tools
// ABOUTME: Tests news API search functionality with mock HTTP responses

package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewSearchNewsAPITool(t *testing.T) {
	tool := NewSearchNewsAPITool()

	if tool.Name() != "search_news_api" {
		t.Errorf("Expected tool name 'search_news_api', got '%s'", tool.Name())
	}

	if tool.Description() != "Searches for news articles using a news API service" {
		t.Errorf("Unexpected tool description: %s", tool.Description())
	}

	// Tool interface doesn't expose schema directly in go-llms
}

func TestSearchNewsAPIHandler_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check query parameters
		q := r.URL.Query()
		if q.Get("q") != "artificial intelligence" {
			t.Errorf("Expected query 'artificial intelligence', got '%s'", q.Get("q"))
		}
		if q.Get("apiKey") != "test-api-key" {
			t.Errorf("Expected apiKey 'test-api-key', got '%s'", q.Get("apiKey"))
		}
		if q.Get("language") != "en" {
			t.Errorf("Expected language 'en', got '%s'", q.Get("language"))
		}
		if q.Get("sortBy") != "popularity" {
			t.Errorf("Expected sortBy 'popularity', got '%s'", q.Get("sortBy"))
		}
		if q.Get("pageSize") != "10" {
			t.Errorf("Expected pageSize '10', got '%s'", q.Get("pageSize"))
		}

		// Send mock response
		response := map[string]interface{}{
			"status":       "ok",
			"totalResults": 2,
			"articles": []map[string]interface{}{
				{
					"source": map[string]string{
						"id":   "techcrunch",
						"name": "TechCrunch",
					},
					"author":      "John Doe",
					"title":       "AI Breakthrough in Natural Language Processing",
					"description": "Researchers announce major advancement in AI language models",
					"url":         "https://example.com/ai-breakthrough",
					"publishedAt": "2024-01-15T10:00:00Z",
					"content":     "Full article content here...",
				},
				{
					"source": map[string]string{
						"id":   "wired",
						"name": "Wired",
					},
					"author":      "Jane Smith",
					"title":       "The Future of AI in Healthcare",
					"description": "How artificial intelligence is transforming medical diagnosis",
					"url":         "https://example.com/ai-healthcare",
					"publishedAt": "2024-01-14T15:30:00Z",
					"content":     "Healthcare AI article content...",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Override NewsAPI URL for testing
	originalURL := newsAPIURL
	newsAPIURL = server.URL + "/v2/everything"
	defer func() { newsAPIURL = originalURL }()

	params := SearchNewsAPIParams{
		Query:    "artificial intelligence",
		APIKey:   "test-api-key",
		Language: "en",
		SortBy:   "popularity",
		PageSize: 10,
	}

	ctx := context.Background()
	result, err := searchNewsAPIHandler(ctx, params)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify results
	if result.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result.Status)
	}
	if result.TotalResults != 2 {
		t.Errorf("Expected 2 total results, got %d", result.TotalResults)
	}
	if len(result.Articles) != 2 {
		t.Errorf("Expected 2 articles, got %d", len(result.Articles))
	}

	// Check first article
	if len(result.Articles) > 0 {
		article := result.Articles[0]
		if article.Title != "AI Breakthrough in Natural Language Processing" {
			t.Errorf("Unexpected article title: %s", article.Title)
		}
		if article.Source.Name != "TechCrunch" {
			t.Errorf("Expected source 'TechCrunch', got '%s'", article.Source.Name)
		}
		if article.URL != "https://example.com/ai-breakthrough" {
			t.Errorf("Unexpected article URL: %s", article.URL)
		}
	}
}

func TestSearchNewsAPIHandler_WithFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check date filters
		q := r.URL.Query()
		if q.Get("from") != "2024-01-01" {
			t.Errorf("Expected from date '2024-01-01', got '%s'", q.Get("from"))
		}
		if q.Get("to") != "2024-01-31" {
			t.Errorf("Expected to date '2024-01-31', got '%s'", q.Get("to"))
		}
		if q.Get("domains") != "bbc.com,cnn.com" {
			t.Errorf("Expected domains 'bbc.com,cnn.com', got '%s'", q.Get("domains"))
		}

		// Send response
		response := map[string]interface{}{
			"status":       "ok",
			"totalResults": 1,
			"articles": []map[string]interface{}{
				{
					"source": map[string]string{
						"id":   "bbc-news",
						"name": "BBC News",
					},
					"author":      "BBC Reporter",
					"title":       "Tech News Update",
					"description": "Latest technology news",
					"url":         "https://bbc.com/tech-news",
					"publishedAt": "2024-01-15T12:00:00Z",
					"content":     "Article content...",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	originalURL := newsAPIURL
	newsAPIURL = server.URL + "/v2/everything"
	defer func() { newsAPIURL = originalURL }()

	params := SearchNewsAPIParams{
		Query:    "technology",
		APIKey:   "test-api-key",
		DateFrom: "2024-01-01",
		DateTo:   "2024-01-31",
		Domains:  "bbc.com,cnn.com",
	}

	ctx := context.Background()
	result, err := searchNewsAPIHandler(ctx, params)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.TotalResults != 1 {
		t.Errorf("Expected 1 result, got %d", result.TotalResults)
	}
}

func TestSearchNewsAPIHandler_EmptyQuery(t *testing.T) {
	params := SearchNewsAPIParams{
		Query:  "",
		APIKey: "test-api-key",
	}

	ctx := context.Background()
	_, err := searchNewsAPIHandler(ctx, params)

	if err == nil {
		t.Error("Expected error for empty query")
	}
}

func TestSearchNewsAPIHandler_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send error response
		response := map[string]interface{}{
			"status":  "error",
			"code":    "apiKeyInvalid",
			"message": "Your API key is invalid",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	originalURL := newsAPIURL
	newsAPIURL = server.URL + "/v2/everything"
	defer func() { newsAPIURL = originalURL }()

	params := SearchNewsAPIParams{
		Query:  "test",
		APIKey: "invalid-key",
	}

	ctx := context.Background()
	_, err := searchNewsAPIHandler(ctx, params)

	if err == nil {
		t.Error("Expected error for API error response")
	}
	if err.Error() != "news API error: Your API key is invalid (code: apiKeyInvalid)" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestSearchNewsAPIHandler_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate timeout
		time.Sleep(3 * time.Second)
	}))
	defer server.Close()

	originalURL := newsAPIURL
	newsAPIURL = server.URL + "/v2/everything"
	defer func() { newsAPIURL = originalURL }()

	params := SearchNewsAPIParams{
		Query:  "test",
		APIKey: "test-key",
	}

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := searchNewsAPIHandler(ctx, params)

	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestSearchNewsAPIHandler_Pagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("page") != "2" {
			t.Errorf("Expected page '2', got '%s'", q.Get("page"))
		}
		if q.Get("pageSize") != "50" {
			t.Errorf("Expected pageSize '50', got '%s'", q.Get("pageSize"))
		}

		response := map[string]interface{}{
			"status":       "ok",
			"totalResults": 100,
			"articles":     []map[string]interface{}{}, // Empty for brevity
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	originalURL := newsAPIURL
	newsAPIURL = server.URL + "/v2/everything"
	defer func() { newsAPIURL = originalURL }()

	params := SearchNewsAPIParams{
		Query:    "test",
		APIKey:   "test-key",
		Page:     2,
		PageSize: 50,
	}

	ctx := context.Background()
	result, err := searchNewsAPIHandler(ctx, params)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.TotalResults != 100 {
		t.Errorf("Expected 100 total results, got %d", result.TotalResults)
	}
}

func TestSearchNewsAPIHandler_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte("invalid json")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	originalURL := newsAPIURL
	newsAPIURL = server.URL + "/v2/everything"
	defer func() { newsAPIURL = originalURL }()

	params := SearchNewsAPIParams{
		Query:  "test",
		APIKey: "test-key",
	}

	ctx := context.Background()
	_, err := searchNewsAPIHandler(ctx, params)

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

// Brave Search API Tests

func TestNewSearchWebBraveTool(t *testing.T) {
	tool := NewSearchWebBraveTool()

	if tool.Name() != "search_web_brave" {
		t.Errorf("Expected tool name 'search_web_brave', got '%s'", tool.Name())
	}

	if tool.Description() != "Performs web search using Brave Search API with support for multiple content types, AI summaries, and advanced filtering" {
		t.Errorf("Unexpected tool description: %s", tool.Description())
	}
}

func TestSearchWebBraveHandler_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check headers
		if r.Header.Get("X-Subscription-Token") != "test-api-key" {
			t.Errorf("Expected X-Subscription-Token header")
		}

		// Check query parameters
		q := r.URL.Query()
		if q.Get("q") != "artificial intelligence" {
			t.Errorf("Expected query 'artificial intelligence', got '%s'", q.Get("q"))
		}
		if q.Get("count") != "10" {
			t.Errorf("Expected count '10', got '%s'", q.Get("count"))
		}

		// Send mock response
		response := map[string]any{
			"query": map[string]any{
				"original": "artificial intelligence",
			},
			"type": "search",
			"mixed": map[string]any{
				"type": "mixed",
				"main": []map[string]any{
					{
						"type":  "web",
						"index": 0,
					},
					{
						"type":  "news",
						"index": 0,
					},
				},
			},
			"web": map[string]any{
				"type": "search",
				"results": []map[string]any{
					{
						"title":       "What is Artificial Intelligence?",
						"url":         "https://example.com/ai-intro",
						"description": "Learn about AI and machine learning",
						"age":         "2 days ago",
						"language":    "en",
					},
					{
						"title":       "AI Research Papers",
						"url":         "https://example.com/ai-research",
						"description": "Latest research in artificial intelligence",
						"thumbnail": map[string]any{
							"src":    "https://example.com/thumb.jpg",
							"height": 200,
							"width":  300,
						},
					},
				},
			},
			"news": map[string]any{
				"type": "news",
				"results": []map[string]any{
					{
						"title":       "AI Breakthrough Announced",
						"url":         "https://news.example.com/ai-breakthrough",
						"description": "Major advancement in AI technology",
						"age":         "3 hours ago",
						"source": map[string]any{
							"name": "Tech News",
							"url":  "https://technews.com",
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Override Brave API URL for testing
	originalURL := braveSearchURL
	braveSearchURL = server.URL + "/res/v1/web/search"
	defer func() { braveSearchURL = originalURL }()

	params := SearchWebBraveParams{
		Query:  "artificial intelligence",
		APIKey: "test-api-key",
		Count:  10,
	}

	ctx := context.Background()
	result, err := searchWebBraveHandler(ctx, params)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify results
	if result.Query != "artificial intelligence" {
		t.Errorf("Expected query 'artificial intelligence', got '%s'", result.Query)
	}
	if result.Type != "search" {
		t.Errorf("Expected type 'search', got '%s'", result.Type)
	}
	if len(result.Web.Results) != 2 {
		t.Errorf("Expected 2 web results, got %d", len(result.Web.Results))
	}
	if len(result.News.Results) != 1 {
		t.Errorf("Expected 1 news result, got %d", len(result.News.Results))
	}

	// Check web result details
	if len(result.Web.Results) > 0 {
		webResult := result.Web.Results[0]
		if webResult.Title != "What is Artificial Intelligence?" {
			t.Errorf("Unexpected web result title: %s", webResult.Title)
		}
		if webResult.URL != "https://example.com/ai-intro" {
			t.Errorf("Unexpected web result URL: %s", webResult.URL)
		}
	}
}

func TestSearchWebBraveHandler_WithFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check filter parameters
		q := r.URL.Query()
		if q.Get("safesearch") != "moderate" {
			t.Errorf("Expected safesearch 'moderate', got '%s'", q.Get("safesearch"))
		}
		if q.Get("freshness") != "pd" {
			t.Errorf("Expected freshness 'pd', got '%s'", q.Get("freshness"))
		}
		if q.Get("country") != "US" {
			t.Errorf("Expected country 'US', got '%s'", q.Get("country"))
		}
		if q.Get("search_lang") != "en" {
			t.Errorf("Expected search_lang 'en', got '%s'", q.Get("search_lang"))
		}
		if q.Get("result_filter") != "web,news" {
			t.Errorf("Expected result_filter 'web,news', got '%s'", q.Get("result_filter"))
		}

		// Send response
		response := map[string]any{
			"query": map[string]any{
				"original": "technology",
			},
			"type": "search",
			"web": map[string]any{
				"type": "search",
				"results": []map[string]any{
					{
						"title":       "Latest Tech News",
						"url":         "https://example.com/tech",
						"description": "Technology updates",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	originalURL := braveSearchURL
	braveSearchURL = server.URL + "/res/v1/web/search"
	defer func() { braveSearchURL = originalURL }()

	params := SearchWebBraveParams{
		Query:        "technology",
		APIKey:       "test-api-key",
		Country:      "US",
		SearchLang:   "en",
		SafeSearch:   "moderate",
		Freshness:    "pd",
		ResultFilter: []string{"web", "news"},
	}

	ctx := context.Background()
	result, err := searchWebBraveHandler(ctx, params)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result.Web.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result.Web.Results))
	}
}

func TestSearchWebBraveHandler_WithAISummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check AI summary parameter
		q := r.URL.Query()
		if q.Get("summary") != "true" {
			t.Errorf("Expected summary 'true', got '%s'", q.Get("summary"))
		}

		// Send response with AI summary
		response := map[string]any{
			"query": map[string]any{
				"original": "climate change",
			},
			"type": "search",
			"summary": []map[string]any{
				{
					"type": "summary",
					"key":  "brave_search_llm_summary",
					"text": "Climate change refers to long-term shifts in global temperatures and weather patterns.",
					"enrichments": map[string]any{
						"raw": true,
					},
				},
			},
			"web": map[string]any{
				"type": "search",
				"results": []map[string]any{
					{
						"title":       "Understanding Climate Change",
						"url":         "https://example.com/climate",
						"description": "Comprehensive guide to climate science",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	originalURL := braveSearchURL
	braveSearchURL = server.URL + "/res/v1/web/search"
	defer func() { braveSearchURL = originalURL }()

	params := SearchWebBraveParams{
		Query:   "climate change",
		APIKey:  "test-api-key",
		Summary: true,
	}

	ctx := context.Background()
	result, err := searchWebBraveHandler(ctx, params)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result.Summary) == 0 {
		t.Error("Expected AI summary in results")
	}

	if len(result.Summary) > 0 && result.Summary[0].Text != "Climate change refers to long-term shifts in global temperatures and weather patterns." {
		t.Errorf("Unexpected summary text: %s", result.Summary[0].Text)
	}
}

func TestSearchWebBraveHandler_EmptyQuery(t *testing.T) {
	params := SearchWebBraveParams{
		Query:  "",
		APIKey: "test-api-key",
	}

	ctx := context.Background()
	_, err := searchWebBraveHandler(ctx, params)

	if err == nil {
		t.Error("Expected error for empty query")
	}
}

func TestSearchWebBraveHandler_QueryTooLong(t *testing.T) {
	// Create a query longer than 400 characters
	longQuery := strings.Repeat("test ", 100)

	params := SearchWebBraveParams{
		Query:  longQuery,
		APIKey: "test-api-key",
	}

	ctx := context.Background()
	_, err := searchWebBraveHandler(ctx, params)

	if err == nil {
		t.Error("Expected error for query too long")
	}
	if !strings.Contains(err.Error(), "exceeds 400 character limit") {
		t.Errorf("Expected error about character limit, got: %v", err)
	}
}

func TestSearchWebBraveHandler_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send error response
		w.WriteHeader(http.StatusUnauthorized)
		response := map[string]any{
			"type": "error",
			"code": 401,
			"details": []map[string]any{
				{
					"type":   "error",
					"code":   "invalid_api_key",
					"detail": "Invalid API key",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	originalURL := braveSearchURL
	braveSearchURL = server.URL + "/res/v1/web/search"
	defer func() { braveSearchURL = originalURL }()

	params := SearchWebBraveParams{
		Query:  "test",
		APIKey: "invalid-key",
	}

	ctx := context.Background()
	_, err := searchWebBraveHandler(ctx, params)

	if err == nil {
		t.Error("Expected error for API error response")
	}
	if !strings.Contains(err.Error(), "Invalid API key") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestSearchWebBraveHandler_Pagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("offset") != "20" {
			t.Errorf("Expected offset '20', got '%s'", q.Get("offset"))
		}
		if q.Get("count") != "20" {
			t.Errorf("Expected count '20', got '%s'", q.Get("count"))
		}

		response := map[string]any{
			"query": map[string]any{
				"original": "test",
			},
			"type": "search",
			"web": map[string]any{
				"type":    "search",
				"results": []map[string]any{},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	originalURL := braveSearchURL
	braveSearchURL = server.URL + "/res/v1/web/search"
	defer func() { braveSearchURL = originalURL }()

	params := SearchWebBraveParams{
		Query:  "test",
		APIKey: "test-key",
		Count:  20,
		Offset: 20,
	}

	ctx := context.Background()
	result, err := searchWebBraveHandler(ctx, params)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Query != "test" {
		t.Errorf("Expected query 'test', got '%s'", result.Query)
	}
}

// Research Search API Tests

func TestNewSearchResearchTool(t *testing.T) {
	tool := NewSearchResearchTool()

	if tool.Name() != "search_research" {
		t.Errorf("Expected tool name 'search_research', got '%s'", tool.Name())
	}

	expectedDesc := "Searches for academic papers across multiple research databases (arXiv, PubMed, CORE) in parallel"
	if tool.Description() != expectedDesc {
		t.Errorf("Unexpected tool description: %s", tool.Description())
	}
}

// Research search tests are in api_tools_research_test.go

func TestResearchToolRegistration(t *testing.T) {
	tool := NewSearchResearchTool()

	if tool.Name() != "search_research" {
		t.Errorf("Expected tool name 'search_research', got '%s'", tool.Name())
	}

	if tool.Description() == "" {
		t.Error("Expected non-empty tool description")
	}
}
