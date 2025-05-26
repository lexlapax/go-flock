// ABOUTME: Tests for the research search API tool
// ABOUTME: Validates search_research tool functionality with proper mock responses

package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestResearchPaperAPIHandler_SuccessWithProperMocks(t *testing.T) {
	// Set up mock servers with correct response structures
	arxivServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock arXiv response (Atom feed)
		arxivResponse := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:opensearch="http://a9.com/-/spec/opensearch/1.1/">
	<totalResults xmlns:opensearch="http://a9.com/-/spec/opensearch/1.1/">1</totalResults>
	<entry>
		<id>http://arxiv.org/abs/2101.00001v1</id>
		<title>Deep Learning Advances in 2024</title>
		<summary>We present recent advances in deep learning architectures...</summary>
		<author><name>John Doe</name></author>
		<author><name>Jane Smith</name></author>
		<published>2024-01-15T00:00:00Z</published>
		<link href="http://arxiv.org/abs/2101.00001v1" rel="alternate" type="text/html"/>
		<link href="http://arxiv.org/pdf/2101.00001v1" rel="related" type="application/pdf"/>
	</entry>
</feed>`

		w.Header().Set("Content-Type", "application/atom+xml")
		_, _ = w.Write([]byte(arxivResponse))
	}))
	defer arxivServer.Close()

	pubmedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock PubMed response
		if r.URL.Path == "/entrez/eutils/esearch.fcgi" {
			// Search response (JSON format)
			searchResponse := `{
				"esearchresult": {
					"count": "1",
					"idlist": ["35000001"]
				}
			}`
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(searchResponse))
		} else if r.URL.Path == "/entrez/eutils/efetch.fcgi" {
			// Fetch response (XML format)
			fetchResponse := `<?xml version="1.0" encoding="UTF-8"?>
<PubmedArticleSet>
	<PubmedArticle>
		<MedlineCitation>
			<PMID>35000001</PMID>
			<Article>
				<ArticleTitle>Machine Learning in Medical Diagnosis</ArticleTitle>
				<Abstract>
					<AbstractText>This study explores machine learning applications...</AbstractText>
				</Abstract>
				<AuthorList>
					<Author>
						<LastName>Johnson</LastName>
						<ForeName>Alice</ForeName>
					</Author>
				</AuthorList>
				<Journal>
					<Title>Journal of Medical AI</Title>
				</Journal>
				<ArticleDate>
					<Year>2024</Year>
					<Month>01</Month>
					<Day>20</Day>
				</ArticleDate>
			</Article>
		</MedlineCitation>
		<PubmedData>
			<ArticleIdList>
				<ArticleId IdType="doi">10.1001/jama.2024.12345</ArticleId>
				<ArticleId IdType="pubmed">35000001</ArticleId>
			</ArticleIdList>
		</PubmedData>
	</PubmedArticle>
</PubmedArticleSet>`
			w.Header().Set("Content-Type", "text/xml")
			_, _ = w.Write([]byte(fetchResponse))
		}
	}))
	defer pubmedServer.Close()

	coreServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock CORE response
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-core-key" {
			t.Errorf("Expected Authorization header with Bearer token, got %s", authHeader)
		}

		coreResponse := map[string]any{
			"totalHits": 1,
			"results": []map[string]any{
				{
					"id":            "core-123456",
					"title":         "Machine Learning: A Comprehensive Survey",
					"abstract":      "This paper provides a comprehensive survey of ML techniques...",
					"authors":       []string{"Bob Wilson"},
					"publishedDate": "2024-01-10",
					"doi":           "10.1234/example.doi",
					"links": []map[string]any{
						{"url": "https://core.ac.uk/display/123456", "type": "display"},
						{"url": "https://core.ac.uk/download/pdf/123456.pdf", "type": "pdf"},
					},
					"journals": []string{"Computer Science Review"},
					"score":    0.95,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(coreResponse); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer coreServer.Close()

	// Override API URLs for testing
	arxivAPIURL = arxivServer.URL + "/api/query"
	pubmedBaseURL = pubmedServer.URL + "/entrez/eutils/"
	coreAPIURL = coreServer.URL + "/v3/search/works"

	// Set test CORE API key
	os.Setenv("CORE_API_KEY", "test-core-key")
	defer os.Unsetenv("CORE_API_KEY")

	// Create test parameters
	params := ResearchPaperAPIParams{
		Query:      "machine learning",
		MaxResults: 10,
	}

	// Execute handler
	ctx := context.Background()
	result, err := researchPaperAPIHandler(ctx, params)

	// Assertions
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.Query != "machine learning" {
		t.Errorf("Expected query 'machine learning', got %s", result.Query)
	}
	if result.FetchedAt == "" {
		t.Error("Expected non-empty FetchedAt")
	}

	// Should have results from all three providers
	if len(result.Papers) < 3 {
		t.Errorf("Expected at least 3 papers, got %d", len(result.Papers))
	}

	// Check provider info
	if len(result.Providers) != 3 {
		t.Errorf("Expected 3 providers, got %d", len(result.Providers))
	}

	// Verify we have papers from each provider
	providerCounts := make(map[string]int)
	for _, paper := range result.Papers {
		providerCounts[paper.Source]++
	}

	if providerCounts["arXiv"] == 0 {
		t.Error("Expected papers from arXiv")
	}
	if providerCounts["PubMed"] == 0 {
		t.Error("Expected papers from PubMed")
	}
	if providerCounts["CORE"] == 0 {
		t.Error("Expected papers from CORE")
	}

	// Verify paper details
	for _, paper := range result.Papers {
		if paper.Title == "" {
			t.Errorf("Paper from %s has empty title", paper.Source)
		}
		if len(paper.Authors) == 0 {
			t.Errorf("Paper from %s has no authors", paper.Source)
		}
		if paper.Abstract == "" {
			t.Errorf("Paper from %s has empty abstract", paper.Source)
		}
		if paper.Source == "" {
			t.Error("Paper has empty source")
		}
		if paper.URL == "" {
			t.Errorf("Paper from %s has empty URL", paper.Source)
		}

		// Check provider-specific fields
		switch paper.Source {
		case "arXiv":
			if paper.ArxivID == "" {
				t.Error("arXiv paper missing ArxivID")
			}
			if paper.PDFURL == "" {
				t.Error("arXiv paper missing PDFURL")
			}
		case "PubMed":
			if paper.PubMedID == "" {
				t.Error("PubMed paper missing PubMedID")
			}
			if paper.DOI == "" {
				t.Error("PubMed paper missing DOI")
			}
			if paper.Journal == "" {
				t.Error("PubMed paper missing Journal")
			}
		case "CORE":
			if paper.DOI == "" {
				t.Error("CORE paper missing DOI")
			}
			if paper.PDFURL == "" {
				t.Error("CORE paper missing PDFURL")
			}
			if paper.RelevanceScore <= 0.0 {
				t.Error("CORE paper has invalid relevance score")
			}
		}
	}
}

func TestResearchPaperAPIHandler_PartialFailure(t *testing.T) {
	// Set up only arXiv server (others will fail)
	arxivServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arxivResponse := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
	<entry>
		<id>http://arxiv.org/abs/2101.00001v1</id>
		<title>Test Paper</title>
		<summary>Abstract</summary>
		<author><name>Author</name></author>
		<published>2024-01-15T00:00:00Z</published>
	</entry>
</feed>`
		_, _ = w.Write([]byte(arxivResponse))
	}))
	defer arxivServer.Close()

	// Override only arXiv URL
	arxivAPIURL = arxivServer.URL + "/api/query"
	pubmedBaseURL = "http://invalid-url/"
	coreAPIURL = "http://invalid-url/"

	params := ResearchPaperAPIParams{
		Query:      "test",
		MaxResults: 5,
	}

	ctx := context.Background()
	result, err := researchPaperAPIHandler(ctx, params)

	// Should succeed with partial results
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Should have results from arXiv only
	if len(result.Papers) < 1 {
		t.Errorf("Expected at least 1 paper, got %d", len(result.Papers))
	}

	// Check provider info shows failures
	hasError := false
	for _, provider := range result.Providers {
		if provider.Error != "" {
			hasError = true
			break
		}
	}
	if !hasError {
		t.Error("Expected at least one provider to have an error")
	}
}

func TestResearchPaperAPIHandler_NoAPIKey(t *testing.T) {
	// Ensure no CORE API key is set
	os.Unsetenv("CORE_API_KEY")

	// Set up mock servers
	arxivServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><feed xmlns="http://www.w3.org/2005/Atom"></feed>`))
	}))
	defer arxivServer.Close()

	arxivAPIURL = arxivServer.URL

	params := ResearchPaperAPIParams{
		Query:      "test",
		MaxResults: 5,
		Providers:  []string{"core"}, // Only request CORE
	}

	ctx := context.Background()
	result, err := researchPaperAPIHandler(ctx, params)

	// Should not fail completely, but CORE should have an error
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Find CORE provider info
	foundCore := false
	for _, provider := range result.Providers {
		if provider.Name == "core" {
			foundCore = true
			if !strings.Contains(provider.Error, "CORE_API_KEY") {
				t.Errorf("Expected CORE error to mention CORE_API_KEY, got: %s", provider.Error)
			}
			break
		}
	}
	if !foundCore {
		t.Error("Expected to find CORE provider info")
	}
}
