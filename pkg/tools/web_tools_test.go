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

// Sample HTML for testing
const sampleHTML = `<!DOCTYPE html>
<html>
<head>
	<title>Test Page</title>
	<meta name="description" content="A test page for web tools">
	<meta property="og:title" content="Test Page OG Title">
	<meta property="og:description" content="Test page for Open Graph">
	<meta property="og:image" content="https://example.com/image.jpg">
	<meta name="twitter:card" content="summary_large_image">
	<meta name="twitter:title" content="Test Page Twitter Title">
</head>
<body>
	<h1>Welcome to the Test Page</h1>
	<p>This is a paragraph with <a href="https://example.com">a link</a>.</p>
	<p>Here's another <a href="/relative/path">relative link</a>.</p>
	<img src="https://example.com/image1.jpg" alt="Image 1">
	<img src="/images/local.png" alt="Local Image">
	<a href="mailto:test@example.com">Email link</a>
	<a href="https://external.com" target="_blank">External link</a>
</body>
</html>`

func TestNewFetchWebPageTool(t *testing.T) {
	tool := NewFetchWebPageTool()

	if tool == nil {
		t.Fatal("NewFetchWebPageTool returned nil")
	}

	if tool.Name() != "fetch_webpage" {
		t.Errorf("Expected tool name 'fetch_webpage', got %s", tool.Name())
	}

	if tool.ParameterSchema() == nil {
		t.Error("Tool parameter schema is nil")
	}
}

func TestFetchWebPageHandler_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, sampleHTML)
	}))
	defer server.Close()

	ctx := context.Background()
	params := FetchWebPageParams{
		URL:         server.URL,
		ExtractText: true,
	}

	result, err := fetchWebPageHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify result
	if result.URL != server.URL {
		t.Errorf("Expected URL %s, got %s", server.URL, result.URL)
	}

	if result.Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got %s", result.Title)
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}

	if !strings.Contains(result.Content, "Welcome to the Test Page") {
		t.Error("Expected content to contain 'Welcome to the Test Page'")
	}

	if !strings.Contains(result.ContentType, "text/html") {
		t.Errorf("Expected content type to contain 'text/html', got %s", result.ContentType)
	}
}

func TestFetchWebPageHandler_NoTextExtraction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, sampleHTML)
	}))
	defer server.Close()

	ctx := context.Background()
	params := FetchWebPageParams{
		URL:         server.URL,
		ExtractText: false,
	}

	result, err := fetchWebPageHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should contain full HTML
	if !strings.Contains(result.Content, "<!DOCTYPE html>") {
		t.Error("Expected content to contain full HTML")
	}

	if !strings.Contains(result.Content, "<body>") {
		t.Error("Expected content to contain HTML tags")
	}
}

func TestFetchWebPageHandler_CustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check custom headers
		if r.Header.Get("X-Custom-Header") != "test-value" {
			t.Errorf("Expected custom header 'test-value', got %s", r.Header.Get("X-Custom-Header"))
		}
		if r.Header.Get("User-Agent") != "TestBot/1.0" {
			t.Errorf("Expected User-Agent 'TestBot/1.0', got %s", r.Header.Get("User-Agent"))
		}

		fmt.Fprint(w, "<html><body>OK</body></html>")
	}))
	defer server.Close()

	ctx := context.Background()
	params := FetchWebPageParams{
		URL:       server.URL,
		UserAgent: "TestBot/1.0",
		Headers: map[string]string{
			"X-Custom-Header": "test-value",
		},
	}

	_, err := fetchWebPageHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestFetchWebPageHandler_Redirect(t *testing.T) {
	redirected := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/redirected", http.StatusMovedPermanently)
			return
		}
		redirected = true
		fmt.Fprint(w, "<html><body>Redirected</body></html>")
	}))
	defer server.Close()

	ctx := context.Background()

	// Test with follow redirects enabled (default)
	params := FetchWebPageParams{
		URL:             server.URL,
		FollowRedirects: true,
	}

	result, err := fetchWebPageHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !redirected {
		t.Error("Expected request to follow redirect")
	}

	if !strings.Contains(result.Content, "Redirected") {
		t.Error("Expected content from redirected page")
	}
}

func TestFetchWebPageHandler_NoFollowRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/redirected", http.StatusMovedPermanently)
	}))
	defer server.Close()

	ctx := context.Background()
	params := FetchWebPageParams{
		URL:             server.URL,
		FollowRedirects: false,
	}

	result, err := fetchWebPageHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should get redirect status code
	if result.StatusCode != 301 {
		t.Errorf("Expected status code 301, got %d", result.StatusCode)
	}
}

func TestFetchWebPageHandler_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprint(w, "Too late")
	}))
	defer server.Close()

	ctx := context.Background()
	params := FetchWebPageParams{
		URL:     server.URL,
		Timeout: 1, // 1 second timeout
	}

	_, err := fetchWebPageHandler(ctx, params)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "deadline") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

func TestFetchWebPageHandler_InvalidURL(t *testing.T) {
	ctx := context.Background()
	params := FetchWebPageParams{
		URL: "not-a-valid-url",
	}

	_, err := fetchWebPageHandler(ctx, params)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestFetchWebPageHandler_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
	}))
	defer server.Close()

	ctx := context.Background()
	params := FetchWebPageParams{
		URL: server.URL,
	}

	result, err := fetchWebPageHandler(ctx, params)
	// Should not error, but should capture the status code
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", result.StatusCode)
	}
}

func TestFetchWebPageHandler_NonHTMLContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"key": "value"}`)
	}))
	defer server.Close()

	ctx := context.Background()
	params := FetchWebPageParams{
		URL:         server.URL,
		ExtractText: true,
	}

	result, err := fetchWebPageHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should still return the content
	if !strings.Contains(result.Content, `{"key": "value"}`) {
		t.Error("Expected JSON content to be returned")
	}

	if result.ContentType != "application/json" {
		t.Errorf("Expected content type 'application/json', got %s", result.ContentType)
	}
}

// Tests for ExtractLinks tool

func TestNewExtractLinksTool(t *testing.T) {
	tool := NewExtractLinksTool()

	if tool == nil {
		t.Fatal("NewExtractLinksTool returned nil")
	}

	if tool.Name() != "extract_links" {
		t.Errorf("Expected tool name 'extract_links', got %s", tool.Name())
	}

	if tool.ParameterSchema() == nil {
		t.Error("Tool parameter schema is nil")
	}
}

func TestExtractLinksHandler_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, sampleHTML)
	}))
	defer server.Close()

	ctx := context.Background()
	params := ExtractLinksParams{
		URL:             server.URL,
		IncludeExternal: true,
		IncludeInternal: true,
		IncludeEmail:    true,
		IncludeMedia:    true,
	}

	result, err := extractLinksHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify basic results
	if result.URL != server.URL {
		t.Errorf("Expected URL %s, got %s", server.URL, result.URL)
	}

	// Check internal links
	if len(result.InternalLinks) != 1 {
		t.Errorf("Expected 1 internal link, got %d", len(result.InternalLinks))
	}

	// Check external links
	if len(result.ExternalLinks) != 2 {
		t.Errorf("Expected 2 external links, got %d", len(result.ExternalLinks))
	}

	// Check email links
	if len(result.EmailLinks) != 1 {
		t.Errorf("Expected 1 email link, got %d", len(result.EmailLinks))
	}

	// Check media links
	if len(result.MediaLinks) != 2 {
		t.Errorf("Expected 2 media links, got %d", len(result.MediaLinks))
	}

	// Verify total count
	expectedTotal := len(result.InternalLinks) + len(result.ExternalLinks) + len(result.EmailLinks)
	if result.TotalLinks != expectedTotal {
		t.Errorf("Expected total links %d, got %d", expectedTotal, result.TotalLinks)
	}
}

func TestExtractLinksHandler_FilterTypes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, sampleHTML)
	}))
	defer server.Close()

	ctx := context.Background()

	// Test with only external links
	params := ExtractLinksParams{
		URL:             server.URL,
		IncludeExternal: true,
		IncludeInternal: false,
		IncludeEmail:    false,
		IncludeMedia:    false,
	}

	result, err := extractLinksHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result.InternalLinks) != 0 {
		t.Errorf("Expected 0 internal links when disabled, got %d", len(result.InternalLinks))
	}

	if len(result.ExternalLinks) == 0 {
		t.Error("Expected external links to be included")
	}

	if len(result.EmailLinks) != 0 {
		t.Errorf("Expected 0 email links when disabled, got %d", len(result.EmailLinks))
	}

	if len(result.MediaLinks) != 0 {
		t.Errorf("Expected 0 media links when disabled, got %d", len(result.MediaLinks))
	}
}

func TestExtractLinksHandler_NoLinks(t *testing.T) {
	noLinksHTML := `<html><body><p>No links here</p></body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, noLinksHTML)
	}))
	defer server.Close()

	ctx := context.Background()
	params := ExtractLinksParams{
		URL: server.URL,
	}

	result, err := extractLinksHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.TotalLinks != 0 {
		t.Errorf("Expected 0 total links, got %d", result.TotalLinks)
	}
}

func TestExtractLinksHandler_RelativeURLResolution(t *testing.T) {
	htmlWithRelativeLinks := `<html><body>
		<a href="/page1">Page 1</a>
		<a href="page2">Page 2</a>
		<a href="../page3">Page 3</a>
		<a href="//example.com/page4">Page 4</a>
	</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, htmlWithRelativeLinks)
	}))
	defer server.Close()

	ctx := context.Background()
	params := ExtractLinksParams{
		URL: server.URL + "/subdir/page.html",
	}

	result, err := extractLinksHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// All relative links should be resolved to absolute URLs
	for _, link := range result.InternalLinks {
		if !strings.HasPrefix(link.URL, "http") {
			t.Errorf("Expected absolute URL, got relative: %s", link.URL)
		}
	}
}

// Tests for ExtractMetadata tool

func TestNewExtractMetadataTool(t *testing.T) {
	tool := NewExtractMetadataTool()

	if tool == nil {
		t.Fatal("NewExtractMetadataTool returned nil")
	}

	if tool.Name() != "extract_metadata" {
		t.Errorf("Expected tool name 'extract_metadata', got %s", tool.Name())
	}

	if tool.ParameterSchema() == nil {
		t.Error("Tool parameter schema is nil")
	}
}

func TestExtractMetadataHandler_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, sampleHTML)
	}))
	defer server.Close()

	ctx := context.Background()
	params := ExtractMetadataParams{
		URL: server.URL,
	}

	result, err := extractMetadataHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify basic metadata
	if result.Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got %s", result.Title)
	}

	if result.Description != "A test page for web tools" {
		t.Errorf("Expected description 'A test page for web tools', got %s", result.Description)
	}

	// Check Open Graph data
	if result.OpenGraph == nil {
		t.Fatal("Expected OpenGraph data, got nil")
	}

	if result.OpenGraph["title"] != "Test Page OG Title" {
		t.Errorf("Expected OG title 'Test Page OG Title', got %s", result.OpenGraph["title"])
	}

	if result.OpenGraph["image"] != "https://example.com/image.jpg" {
		t.Errorf("Expected OG image URL, got %s", result.OpenGraph["image"])
	}

	// Check Twitter Card data
	if result.Twitter == nil {
		t.Fatal("Expected Twitter data, got nil")
	}

	if result.Twitter["card"] != "summary_large_image" {
		t.Errorf("Expected Twitter card type, got %s", result.Twitter["card"])
	}
}

func TestExtractMetadataHandler_ComplexMetadata(t *testing.T) {
	complexHTML := `<!DOCTYPE html>
<html>
<head>
	<title>Complex Page</title>
	<meta name="description" content="A complex test page">
	<meta name="author" content="Test Author">
	<meta name="keywords" content="test, metadata, extraction">
	
	<!-- Open Graph -->
	<meta property="og:title" content="Complex OG Title">
	<meta property="og:description" content="Complex OG Description">
	<meta property="og:type" content="article">
	<meta property="og:url" content="https://example.com/complex">
	<meta property="og:image" content="https://example.com/og-image.jpg">
	<meta property="og:image:width" content="1200">
	<meta property="og:image:height" content="630">
	
	<!-- Twitter Card -->
	<meta name="twitter:card" content="summary_large_image">
	<meta name="twitter:site" content="@example">
	<meta name="twitter:creator" content="@author">
	<meta name="twitter:title" content="Complex Twitter Title">
	<meta name="twitter:description" content="Complex Twitter Description">
	<meta name="twitter:image" content="https://example.com/twitter-image.jpg">
	
	<!-- Other meta tags -->
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="robots" content="index, follow">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
</head>
<body>Content</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, complexHTML)
	}))
	defer server.Close()

	ctx := context.Background()
	params := ExtractMetadataParams{
		URL: server.URL,
	}

	result, err := extractMetadataHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check author
	if result.Author != "Test Author" {
		t.Errorf("Expected author 'Test Author', got %s", result.Author)
	}

	// Check keywords
	expectedKeywords := []string{"test", "metadata", "extraction"}
	if len(result.Keywords) != len(expectedKeywords) {
		t.Errorf("Expected %d keywords, got %d", len(expectedKeywords), len(result.Keywords))
	}

	// Check OpenGraph properties
	if result.OpenGraph["type"] != "article" {
		t.Errorf("Expected OG type 'article', got %s", result.OpenGraph["type"])
	}

	if result.OpenGraph["image:width"] != "1200" {
		t.Errorf("Expected OG image width '1200', got %s", result.OpenGraph["image:width"])
	}

	// Check Twitter properties
	if result.Twitter["site"] != "@example" {
		t.Errorf("Expected Twitter site '@example', got %s", result.Twitter["site"])
	}

	// Check other meta tags
	if result.Meta["viewport"] != "width=device-width, initial-scale=1.0" {
		t.Errorf("Expected viewport meta tag, got %s", result.Meta["viewport"])
	}
}

func TestExtractMetadataHandler_NoMetadata(t *testing.T) {
	minimalHTML := `<html><head><title>Minimal</title></head><body>No metadata</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, minimalHTML)
	}))
	defer server.Close()

	ctx := context.Background()
	params := ExtractMetadataParams{
		URL: server.URL,
	}

	result, err := extractMetadataHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Title != "Minimal" {
		t.Errorf("Expected title 'Minimal', got %s", result.Title)
	}

	if result.Description != "" {
		t.Errorf("Expected empty description, got %s", result.Description)
	}

	if len(result.Keywords) != 0 {
		t.Errorf("Expected no keywords, got %d", len(result.Keywords))
	}
}

// Tests for CheckURLStatus tool

func TestNewCheckURLStatusTool(t *testing.T) {
	tool := NewCheckURLStatusTool()

	if tool == nil {
		t.Fatal("NewCheckURLStatusTool returned nil")
	}

	if tool.Name() != "check_url_status" {
		t.Errorf("Expected tool name 'check_url_status', got %s", tool.Name())
	}

	if tool.ParameterSchema() == nil {
		t.Error("Tool parameter schema is nil")
	}
}

func TestCheckURLStatusHandler_Alive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", "13")
		fmt.Fprint(w, "Hello, World!")
	}))
	defer server.Close()

	ctx := context.Background()
	params := CheckURLStatusParams{
		URL: server.URL,
	}

	result, err := checkURLStatusHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Status != "alive" {
		t.Errorf("Expected status 'alive', got %s", result.Status)
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}

	if result.ContentType != "text/html" {
		t.Errorf("Expected content type 'text/html', got %s", result.ContentType)
	}

	if result.ContentLength != 13 {
		t.Errorf("Expected content length 13, got %d", result.ContentLength)
	}

	if result.ResponseTime < 0 {
		t.Errorf("Expected non-negative response time, got %d", result.ResponseTime)
	}
}

func TestCheckURLStatusHandler_Redirect(t *testing.T) {
	redirectCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			redirectCount++
			http.Redirect(w, r, "/redirected", http.StatusMovedPermanently)
			return
		}
		fmt.Fprint(w, "Redirected content")
	}))
	defer server.Close()

	ctx := context.Background()
	params := CheckURLStatusParams{
		URL:             server.URL,
		FollowRedirects: true,
	}

	result, err := checkURLStatusHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Status != "redirect" {
		t.Errorf("Expected status 'redirect', got %s", result.Status)
	}

	if result.FinalURL == "" || result.FinalURL == server.URL {
		t.Error("Expected final URL to be different from original")
	}

	if len(result.RedirectChain) == 0 {
		t.Error("Expected redirect chain to contain entries")
	}
}

func TestCheckURLStatusHandler_Dead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	ctx := context.Background()
	params := CheckURLStatusParams{
		URL: server.URL,
	}

	result, err := checkURLStatusHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Status != "dead" {
		t.Errorf("Expected status 'dead', got %s", result.Status)
	}

	if result.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", result.StatusCode)
	}
}

func TestCheckURLStatusHandler_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprint(w, "Too late")
	}))
	defer server.Close()

	ctx := context.Background()
	params := CheckURLStatusParams{
		URL:     server.URL,
		Timeout: 1,
	}

	_, err := checkURLStatusHandler(ctx, params)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "deadline") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

func TestCheckURLStatusHandler_NoFollowRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/redirected", http.StatusMovedPermanently)
	}))
	defer server.Close()

	ctx := context.Background()
	params := CheckURLStatusParams{
		URL:             server.URL,
		FollowRedirects: false,
	}

	result, err := checkURLStatusHandler(ctx, params)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Status != "redirect" {
		t.Errorf("Expected status 'redirect', got %s", result.Status)
	}

	if result.StatusCode != 301 {
		t.Errorf("Expected status code 301, got %d", result.StatusCode)
	}

	// Should not have followed the redirect
	if result.FinalURL != "" {
		t.Error("Expected no final URL when not following redirects")
	}
}
