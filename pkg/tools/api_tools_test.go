// ABOUTME: Unit tests for API interaction tools
// ABOUTME: Tests news API search functionality with mock HTTP responses

package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
