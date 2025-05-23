// ABOUTME: Web scraping and content extraction tools
// ABOUTME: Provides tools to fetch and extract content from web pages

package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	domain "github.com/lexlapax/go-llms/pkg/agent/domain"
	"github.com/lexlapax/go-llms/pkg/agent/tools"
	sdomain "github.com/lexlapax/go-llms/pkg/schema/domain"
)

// Tool Parameters
type FetchWebPageParams struct {
	URL             string            `json:"url" description:"The URL of the web page to fetch"`
	Timeout         int               `json:"timeout,omitempty" description:"Timeout in seconds (default: 30)"`
	UserAgent       string            `json:"user_agent,omitempty" description:"User agent string for the request"`
	Headers         map[string]string `json:"headers,omitempty" description:"Additional HTTP headers"`
	ExtractText     bool              `json:"extract_text,omitempty" description:"Extract only text content (default: true)"`
	FollowRedirects bool              `json:"follow_redirects,omitempty" description:"Follow HTTP redirects (default: true)"`
}

// Tool Results
type FetchWebPageResult struct {
	URL         string            `json:"url"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	ContentType string            `json:"content_type"`
	StatusCode  int               `json:"status_code"`
	Headers     map[string]string `json:"headers"`
	FetchedAt   string            `json:"fetched_at"`
}

// Schema definitions
var FetchWebPageParamSchema = &sdomain.Schema{
	Type:        "object",
	Description: "Parameters for fetching a web page",
	Properties: map[string]sdomain.Property{
		"url": {
			Type:        "string",
			Description: "The URL of the web page to fetch",
		},
		"timeout": {
			Type:        "integer",
			Description: "Timeout in seconds (default: 30)",
			Minimum:     float64Ptr(1),
			Maximum:     float64Ptr(300),
		},
		"user_agent": {
			Type:        "string",
			Description: "User agent string for the request",
		},
		"headers": {
			Type:                 "object",
			Description:          "Additional HTTP headers",
			AdditionalProperties: func() *bool { b := true; return &b }(),
		},
		"extract_text": {
			Type:        "boolean",
			Description: "Extract only text content (default: true)",
		},
		"follow_redirects": {
			Type:        "boolean",
			Description: "Follow HTTP redirects (default: true)",
		},
	},
	Required: []string{"url"},
}

// NewFetchWebPageTool creates a new web page fetching tool
func NewFetchWebPageTool() domain.Tool {
	return tools.NewTool(
		"fetch_webpage",
		"Fetches content from a web page and optionally extracts text",
		fetchWebPageHandler,
		FetchWebPageParamSchema,
	)
}

func fetchWebPageHandler(ctx context.Context, params FetchWebPageParams) (*FetchWebPageResult, error) {
	// Set defaults
	timeout := 30
	if params.Timeout > 0 {
		timeout = params.Timeout
	}

	extractText := true
	if !params.ExtractText {
		extractText = false
	}

	followRedirects := true
	if !params.FollowRedirects {
		followRedirects = false
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Configure redirect policy
	if !followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", params.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	userAgent := "go-flock/1.0 WebFetcher"
	if params.UserAgent != "" {
		userAgent = params.UserAgent
	}
	req.Header.Set("User-Agent", userAgent)

	// Add custom headers
	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching web page: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// Process content
	content := string(body)
	title := ""

	if extractText && strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		// Extract title
		title = extractTitle(content)

		// Extract text content
		content = extractTextFromHTML(content)
	} else if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		// Still extract title even if not extracting text
		title = extractTitle(string(body))
	}

	// Build result
	result := &FetchWebPageResult{
		URL:         params.URL,
		Title:       title,
		Content:     content,
		ContentType: resp.Header.Get("Content-Type"),
		StatusCode:  resp.StatusCode,
		Headers:     make(map[string]string),
		FetchedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// Copy relevant headers
	for key, values := range resp.Header {
		if len(values) > 0 {
			result.Headers[key] = values[0]
		}
	}

	return result, nil
}

// extractTitle extracts the title from HTML content
func extractTitle(html string) string {
	// Try to find <title> tag
	titleRegex := regexp.MustCompile(`(?i)<title[^>]*>([^<]+)</title>`)
	matches := titleRegex.FindStringSubmatch(html)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// extractTextFromHTML strips HTML tags and returns text content
func extractTextFromHTML(html string) string {
	// Remove script and style elements
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	html = scriptRegex.ReplaceAllString(html, "")

	styleRegex := regexp.MustCompile(`(?i)<style[^>]*>.*?</style>`)
	html = styleRegex.ReplaceAllString(html, "")

	// Remove HTML comments
	commentRegex := regexp.MustCompile(`<!--.*?-->`)
	html = commentRegex.ReplaceAllString(html, "")

	// Replace common block elements with newlines
	blockRegex := regexp.MustCompile(`(?i)<(p|div|br|h[1-6]|li|tr)[^>]*>`)
	html = blockRegex.ReplaceAllString(html, "\n")

	// Remove all remaining HTML tags
	tagRegex := regexp.MustCompile(`<[^>]+>`)
	text := tagRegex.ReplaceAllString(html, "")

	// Decode HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")

	// Clean up whitespace
	lines := strings.Split(text, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// Extract Links Tool

type ExtractLinksParams struct {
	URL             string `json:"url" description:"The URL of the web page to extract links from"`
	IncludeExternal bool   `json:"include_external,omitempty" description:"Include external links (default: true)"`
	IncludeInternal bool   `json:"include_internal,omitempty" description:"Include internal links (default: true)"`
	IncludeEmail    bool   `json:"include_email,omitempty" description:"Include email links (default: true)"`
	IncludeMedia    bool   `json:"include_media,omitempty" description:"Include media links (images, videos) (default: false)"`
	Timeout         int    `json:"timeout,omitempty" description:"Timeout in seconds (default: 30)"`
}

type ExtractLinksResult struct {
	URL           string      `json:"url"`
	TotalLinks    int         `json:"total_links"`
	InternalLinks []LinkInfo  `json:"internal_links"`
	ExternalLinks []LinkInfo  `json:"external_links"`
	EmailLinks    []string    `json:"email_links"`
	MediaLinks    []MediaLink `json:"media_links,omitempty"`
	FetchedAt     string      `json:"fetched_at"`
}

type LinkInfo struct {
	URL   string `json:"url"`
	Text  string `json:"text"`
	Title string `json:"title,omitempty"`
}

type MediaLink struct {
	URL  string `json:"url"`
	Type string `json:"type"` // image, video, audio
	Alt  string `json:"alt,omitempty"`
}

var ExtractLinksParamSchema = &sdomain.Schema{
	Type:        "object",
	Description: "Parameters for extracting links from a web page",
	Properties: map[string]sdomain.Property{
		"url": {
			Type:        "string",
			Description: "The URL of the web page to extract links from",
		},
		"include_external": {
			Type:        "boolean",
			Description: "Include external links (default: true)",
		},
		"include_internal": {
			Type:        "boolean",
			Description: "Include internal links (default: true)",
		},
		"include_email": {
			Type:        "boolean",
			Description: "Include email links (default: true)",
		},
		"include_media": {
			Type:        "boolean",
			Description: "Include media links (images, videos) (default: false)",
		},
		"timeout": {
			Type:        "integer",
			Description: "Timeout in seconds (default: 30)",
			Minimum:     float64Ptr(1),
			Maximum:     float64Ptr(300),
		},
	},
	Required: []string{"url"},
}

func NewExtractLinksTool() domain.Tool {
	return tools.NewTool(
		"extract_links",
		"Extracts all links from a web page, categorized by type",
		extractLinksHandler,
		ExtractLinksParamSchema,
	)
}

func extractLinksHandler(ctx context.Context, params ExtractLinksParams) (*ExtractLinksResult, error) {
	// Set defaults
	includeExternal := true
	if !params.IncludeExternal {
		includeExternal = params.IncludeExternal
	}
	includeInternal := true
	if !params.IncludeInternal {
		includeInternal = params.IncludeInternal
	}
	includeEmail := true
	if !params.IncludeEmail {
		includeEmail = params.IncludeEmail
	}
	includeMedia := params.IncludeMedia // Default false

	timeout := 30
	if params.Timeout > 0 {
		timeout = params.Timeout
	}

	// Parse base URL
	baseURL, err := url.Parse(params.URL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}

	// Fetch the web page
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", params.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", "go-flock/1.0 LinkExtractor")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching web page: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	html := string(body)

	// Initialize result
	result := &ExtractLinksResult{
		URL:           params.URL,
		InternalLinks: []LinkInfo{},
		ExternalLinks: []LinkInfo{},
		EmailLinks:    []string{},
		MediaLinks:    []MediaLink{},
		FetchedAt:     time.Now().UTC().Format(time.RFC3339),
	}

	// Extract anchor links
	anchorRegex := regexp.MustCompile(`(?i)<a[^>]+href\s*=\s*["']([^"']+)["'][^>]*>([^<]*)</a>`)
	anchorMatches := anchorRegex.FindAllStringSubmatch(html, -1)

	for _, match := range anchorMatches {
		if len(match) < 3 {
			continue
		}

		href := strings.TrimSpace(match[1])
		text := strings.TrimSpace(match[2])

		// Skip empty hrefs
		if href == "" || href == "#" {
			continue
		}

		// Check for email links
		if strings.HasPrefix(href, "mailto:") {
			if includeEmail {
				email := strings.TrimPrefix(href, "mailto:")
				result.EmailLinks = append(result.EmailLinks, email)
			}
			continue
		}

		// Parse and resolve URL
		linkURL, err := resolveURL(baseURL, href)
		if err != nil {
			continue
		}

		linkInfo := LinkInfo{
			URL:  linkURL.String(),
			Text: text,
		}

		// Extract title attribute if present
		titleRegex := regexp.MustCompile(`title\s*=\s*["']([^"']+)["']`)
		if titleMatch := titleRegex.FindStringSubmatch(match[0]); len(titleMatch) > 1 {
			linkInfo.Title = titleMatch[1]
		}

		// Categorize link
		if isInternalLink(baseURL, linkURL) {
			if includeInternal {
				result.InternalLinks = append(result.InternalLinks, linkInfo)
			}
		} else {
			if includeExternal {
				result.ExternalLinks = append(result.ExternalLinks, linkInfo)
			}
		}
	}

	// Extract media links if requested
	if includeMedia {
		// Extract images
		imgRegex := regexp.MustCompile(`(?i)<img[^>]+src\s*=\s*["']([^"']+)["'][^>]*>`)
		imgMatches := imgRegex.FindAllStringSubmatch(html, -1)

		for _, match := range imgMatches {
			if len(match) < 2 {
				continue
			}

			src := strings.TrimSpace(match[1])
			if src == "" {
				continue
			}

			imgURL, err := resolveURL(baseURL, src)
			if err != nil {
				continue
			}

			mediaLink := MediaLink{
				URL:  imgURL.String(),
				Type: "image",
			}

			// Extract alt text
			altRegex := regexp.MustCompile(`alt\s*=\s*["']([^"']+)["']`)
			if altMatch := altRegex.FindStringSubmatch(match[0]); len(altMatch) > 1 {
				mediaLink.Alt = altMatch[1]
			}

			result.MediaLinks = append(result.MediaLinks, mediaLink)
		}

		// Could add video/audio extraction here in the future
	}

	// Calculate total links (excluding media)
	result.TotalLinks = len(result.InternalLinks) + len(result.ExternalLinks) + len(result.EmailLinks)

	return result, nil
}

// resolveURL resolves a potentially relative URL against a base URL
func resolveURL(base *url.URL, href string) (*url.URL, error) {
	// Parse the href
	hrefURL, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	// Resolve against base
	return base.ResolveReference(hrefURL), nil
}

// isInternalLink checks if a URL is internal relative to the base URL
func isInternalLink(base, link *url.URL) bool {
	return link.Host == base.Host || link.Host == ""
}

// Extract Metadata Tool

type ExtractMetadataParams struct {
	URL     string `json:"url" description:"The URL of the web page to extract metadata from"`
	Timeout int    `json:"timeout,omitempty" description:"Timeout in seconds (default: 30)"`
}

type ExtractMetadataResult struct {
	URL         string            `json:"url"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Author      string            `json:"author,omitempty"`
	Keywords    []string          `json:"keywords,omitempty"`
	OpenGraph   map[string]string `json:"open_graph,omitempty"`
	Twitter     map[string]string `json:"twitter,omitempty"`
	Meta        map[string]string `json:"meta"`
	FetchedAt   string            `json:"fetched_at"`
}

var ExtractMetadataParamSchema = &sdomain.Schema{
	Type:        "object",
	Description: "Parameters for extracting metadata from a web page",
	Properties: map[string]sdomain.Property{
		"url": {
			Type:        "string",
			Description: "The URL of the web page to extract metadata from",
		},
		"timeout": {
			Type:        "integer",
			Description: "Timeout in seconds (default: 30)",
			Minimum:     float64Ptr(1),
			Maximum:     float64Ptr(300),
		},
	},
	Required: []string{"url"},
}

func NewExtractMetadataTool() domain.Tool {
	return tools.NewTool(
		"extract_metadata",
		"Extracts metadata from a web page including Open Graph and Twitter Card data",
		extractMetadataHandler,
		ExtractMetadataParamSchema,
	)
}

func extractMetadataHandler(ctx context.Context, params ExtractMetadataParams) (*ExtractMetadataResult, error) {
	// Set default timeout
	timeout := 30
	if params.Timeout > 0 {
		timeout = params.Timeout
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", params.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", "go-flock/1.0 MetadataExtractor")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching web page: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	html := string(body)

	// Initialize result
	result := &ExtractMetadataResult{
		URL:       params.URL,
		OpenGraph: make(map[string]string),
		Twitter:   make(map[string]string),
		Meta:      make(map[string]string),
		Keywords:  []string{},
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// Extract title
	result.Title = extractTitle(html)

	// Extract all meta tags
	metaRegex := regexp.MustCompile(`(?i)<meta\s+([^>]+)>`)
	metaMatches := metaRegex.FindAllStringSubmatch(html, -1)

	for _, match := range metaMatches {
		if len(match) < 2 {
			continue
		}

		attrs := parseAttributes(match[1])

		// Handle standard meta tags
		if name, ok := attrs["name"]; ok {
			content := attrs["content"]

			switch strings.ToLower(name) {
			case "description":
				result.Description = content
			case "author":
				result.Author = content
			case "keywords":
				// Split keywords by comma
				keywords := strings.Split(content, ",")
				for _, kw := range keywords {
					kw = strings.TrimSpace(kw)
					if kw != "" {
						result.Keywords = append(result.Keywords, kw)
					}
				}
			default:
				// Store other meta tags
				if strings.HasPrefix(name, "twitter:") {
					// Twitter Card data
					twitterKey := strings.TrimPrefix(name, "twitter:")
					result.Twitter[twitterKey] = content
				} else {
					// General meta tags
					result.Meta[name] = content
				}
			}
		}

		// Handle Open Graph meta tags
		if property, ok := attrs["property"]; ok {
			if strings.HasPrefix(property, "og:") {
				ogKey := strings.TrimPrefix(property, "og:")
				result.OpenGraph[ogKey] = attrs["content"]
			}
		}

		// Handle http-equiv meta tags
		if httpEquiv, ok := attrs["http-equiv"]; ok {
			result.Meta[httpEquiv] = attrs["content"]
		}
	}

	return result, nil
}

// parseAttributes parses HTML attributes from a string
func parseAttributes(attrString string) map[string]string {
	attrs := make(map[string]string)

	// Match attribute="value" or attribute='value' patterns
	attrRegex := regexp.MustCompile(`(\w+(?:-\w+)*)\s*=\s*["']([^"']+)["']`)
	matches := attrRegex.FindAllStringSubmatch(attrString, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			attrs[strings.ToLower(match[1])] = match[2]
		}
	}

	return attrs
}

// Check URL Status Tool

type CheckURLStatusParams struct {
	URL             string            `json:"url" description:"The URL to check"`
	Timeout         int               `json:"timeout,omitempty" description:"Timeout in seconds (default: 10)"`
	FollowRedirects bool              `json:"follow_redirects,omitempty" description:"Follow redirects (default: true)"`
	Headers         map[string]string `json:"headers,omitempty" description:"Additional HTTP headers"`
}

type CheckURLStatusResult struct {
	URL           string            `json:"url"`
	Status        string            `json:"status"` // "alive", "dead", "redirect"
	StatusCode    int               `json:"status_code"`
	ResponseTime  int64             `json:"response_time_ms"`
	ContentType   string            `json:"content_type,omitempty"`
	ContentLength int64             `json:"content_length,omitempty"`
	FinalURL      string            `json:"final_url,omitempty"`
	RedirectChain []string          `json:"redirect_chain,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	CheckedAt     string            `json:"checked_at"`
}

var CheckURLStatusParamSchema = &sdomain.Schema{
	Type:        "object",
	Description: "Parameters for checking URL status",
	Properties: map[string]sdomain.Property{
		"url": {
			Type:        "string",
			Description: "The URL to check",
		},
		"timeout": {
			Type:        "integer",
			Description: "Timeout in seconds (default: 10)",
			Minimum:     float64Ptr(1),
			Maximum:     float64Ptr(300),
		},
		"follow_redirects": {
			Type:        "boolean",
			Description: "Follow redirects (default: true)",
		},
		"headers": {
			Type:                 "object",
			Description:          "Additional HTTP headers",
			AdditionalProperties: func() *bool { b := true; return &b }(),
		},
	},
	Required: []string{"url"},
}

func NewCheckURLStatusTool() domain.Tool {
	return tools.NewTool(
		"check_url_status",
		"Checks if a URL is accessible and returns status information",
		checkURLStatusHandler,
		CheckURLStatusParamSchema,
	)
}

func checkURLStatusHandler(ctx context.Context, params CheckURLStatusParams) (*CheckURLStatusResult, error) {
	// Set defaults
	timeout := 10
	if params.Timeout > 0 {
		timeout = params.Timeout
	}

	followRedirects := true
	if !params.FollowRedirects {
		followRedirects = false
	}

	// Track start time
	startTime := time.Now()

	// Create HTTP client
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Track redirect chain
	var redirectChain []string
	var finalURL string

	// Configure redirect policy
	if followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			// Track redirect chain
			redirectChain = append(redirectChain, req.URL.String())

			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		}
	} else {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "HEAD", params.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "go-flock/1.0 URLChecker")

	// Add custom headers
	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		// Check if it's a timeout
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline") {
			return nil, fmt.Errorf("request timeout: %w", err)
		}
		return nil, fmt.Errorf("checking URL: %w", err)
	}
	defer resp.Body.Close()

	// Calculate response time
	responseTime := time.Since(startTime).Milliseconds()

	// Determine status
	status := "alive"
	if resp.StatusCode >= 400 {
		status = "dead"
	} else if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		status = "redirect"
	}

	// Get final URL if redirects were followed
	if followRedirects && resp.Request != nil {
		finalURL = resp.Request.URL.String()
		if finalURL == params.URL {
			finalURL = ""
		}
		// If we followed redirects and ended up somewhere else, it's a redirect
		if finalURL != "" && len(redirectChain) > 0 {
			status = "redirect"
		}
	}

	// Get content length
	contentLength := resp.ContentLength
	if contentLength == -1 {
		contentLength = 0
	}

	// Build result
	result := &CheckURLStatusResult{
		URL:           params.URL,
		Status:        status,
		StatusCode:    resp.StatusCode,
		ResponseTime:  responseTime,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: contentLength,
		FinalURL:      finalURL,
		RedirectChain: redirectChain,
		Headers:       make(map[string]string),
		CheckedAt:     time.Now().UTC().Format(time.RFC3339),
	}

	// Copy relevant headers
	for key, values := range resp.Header {
		if len(values) > 0 {
			result.Headers[key] = values[0]
		}
	}

	return result, nil
}

// Future tools could include:
// - extract_structured_data: Extract specific data using CSS selectors
// - screenshot_webpage: Capture visual representation
// - submit_form: Submit forms on web pages
// - extract_tables: Extract tabular data from HTML tables
