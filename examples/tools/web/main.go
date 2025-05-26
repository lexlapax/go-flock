// ABOUTME: Example demonstrating web scraping and content extraction tools
// ABOUTME: Shows how to fetch web pages and extract links using go-flock tools

package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/lexlapax/go-flock/pkg/tools"
)

func main() {
	fmt.Println("go-flock Web Tools Example")
	fmt.Println("==========================")
	fmt.Println()

	ctx := context.Background()

	// Example 1: Fetch a web page with text extraction
	fmt.Println("Example 1: Fetching a web page with text extraction")
	fmt.Println("---------------------------------------------------")
	fetchWebPageExample(ctx)
	fmt.Println()

	// Example 2: Fetch a web page without text extraction
	fmt.Println("Example 2: Fetching raw HTML content")
	fmt.Println("------------------------------------")
	fetchRawHTMLExample(ctx)
	fmt.Println()

	// Example 3: Extract links from a web page
	fmt.Println("Example 3: Extracting all links from a web page")
	fmt.Println("-----------------------------------------------")
	extractAllLinksExample(ctx)
	fmt.Println()

	// Example 4: Extract specific types of links
	fmt.Println("Example 4: Extracting only external links")
	fmt.Println("------------------------------------------")
	extractExternalLinksExample(ctx)
	fmt.Println()

	// Example 5: Advanced link analysis
	fmt.Println("Example 5: Link analysis for a tech news site")
	fmt.Println("---------------------------------------------")
	analyzeTechNewsLinks(ctx)
	fmt.Println()

	// Example 6: Extract metadata from a web page
	fmt.Println("Example 6: Extracting metadata (Open Graph, Twitter Cards)")
	fmt.Println("----------------------------------------------------------")
	extractMetadataExample(ctx)
	fmt.Println()

	// Example 7: Check URL status
	fmt.Println("Example 7: Checking URL status and response times")
	fmt.Println("-------------------------------------------------")
	checkURLStatusExample(ctx)
}

func fetchWebPageExample(ctx context.Context) {
	tool := tools.NewFetchWebPageTool()

	params := tools.FetchWebPageParams{
		URL:         "https://example.com",
		ExtractText: true,
		Timeout:     10,
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Printf("Error fetching web page: %v", err)
		return
	}

	if webResult, ok := result.(*tools.FetchWebPageResult); ok {
		fmt.Printf("URL: %s\n", webResult.URL)
		fmt.Printf("Title: %s\n", webResult.Title)
		fmt.Printf("Status Code: %d\n", webResult.StatusCode)
		fmt.Printf("Content Type: %s\n", webResult.ContentType)
		fmt.Printf("Fetched At: %s\n", webResult.FetchedAt)

		// Show first 200 characters of extracted text
		content := webResult.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}
		fmt.Printf("Text Content:\n%s\n", content)
	}
}

func fetchRawHTMLExample(ctx context.Context) {
	tool := tools.NewFetchWebPageTool()

	params := tools.FetchWebPageParams{
		URL:         "https://www.wikipedia.org",
		ExtractText: false, // Keep HTML tags
		UserAgent:   "go-flock-bot/1.0",
		Headers: map[string]string{
			"Accept-Language": "en-US,en;q=0.9",
		},
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Printf("Error fetching web page: %v", err)
		return
	}

	if webResult, ok := result.(*tools.FetchWebPageResult); ok {
		fmt.Printf("Fetched raw HTML from: %s\n", webResult.URL)
		fmt.Printf("Title extracted: %s\n", webResult.Title)

		// Show HTML snippet
		content := webResult.Content
		startIdx := strings.Index(content, "<body")
		if startIdx >= 0 && startIdx+100 < len(content) {
			fmt.Printf("HTML snippet: %s...\n", content[startIdx:startIdx+100])
		}

		// Show response headers
		fmt.Println("Response headers:")
		for key, value := range webResult.Headers {
			if strings.HasPrefix(key, "Content-") || key == "Server" {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
	}
}

func extractAllLinksExample(ctx context.Context) {
	tool := tools.NewExtractLinksTool()

	params := tools.ExtractLinksParams{
		URL:             "https://news.ycombinator.com",
		IncludeExternal: true,
		IncludeInternal: true,
		IncludeEmail:    true,
		IncludeMedia:    false,
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Printf("Error extracting links: %v", err)
		return
	}

	if linkResult, ok := result.(*tools.ExtractLinksResult); ok {
		fmt.Printf("Extracted links from: %s\n", linkResult.URL)
		fmt.Printf("Total links found: %d\n", linkResult.TotalLinks)
		fmt.Printf("Internal links: %d\n", len(linkResult.InternalLinks))
		fmt.Printf("External links: %d\n", len(linkResult.ExternalLinks))
		fmt.Printf("Email links: %d\n", len(linkResult.EmailLinks))

		// Show first few internal links
		fmt.Println("\nFirst 3 internal links:")
		for i, link := range linkResult.InternalLinks {
			if i >= 3 {
				break
			}
			fmt.Printf("  - %s\n", link.URL)
			if link.Text != "" {
				fmt.Printf("    Text: %s\n", link.Text)
			}
		}

		// Show first few external links
		fmt.Println("\nFirst 3 external links:")
		for i, link := range linkResult.ExternalLinks {
			if i >= 3 {
				break
			}
			fmt.Printf("  - %s\n", link.URL)
			if link.Text != "" {
				fmt.Printf("    Text: %s\n", link.Text)
			}
		}
	}
}

func extractExternalLinksExample(ctx context.Context) {
	tool := tools.NewExtractLinksTool()

	params := tools.ExtractLinksParams{
		URL:             "https://golang.org",
		IncludeExternal: true,
		IncludeInternal: false, // Only external links
		IncludeEmail:    false,
		IncludeMedia:    false,
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Printf("Error extracting links: %v", err)
		return
	}

	if linkResult, ok := result.(*tools.ExtractLinksResult); ok {
		fmt.Printf("External links from %s:\n", linkResult.URL)

		// Group external links by domain
		domainLinks := make(map[string][]tools.LinkInfo)
		for _, link := range linkResult.ExternalLinks {
			domain := extractDomain(link.URL)
			domainLinks[domain] = append(domainLinks[domain], link)
		}

		// Display grouped links
		for domain, links := range domainLinks {
			fmt.Printf("\n%s (%d links):\n", domain, len(links))
			for i, link := range links {
				if i >= 2 { // Show max 2 per domain
					fmt.Printf("  ... and %d more\n", len(links)-2)
					break
				}
				fmt.Printf("  - %s\n", link.URL)
			}
		}
	}
}

func analyzeTechNewsLinks(ctx context.Context) {
	// First, fetch the page to analyze
	fetchTool := tools.NewFetchWebPageTool()
	linkTool := tools.NewExtractLinksTool()

	// Fetch TechCrunch homepage
	fetchParams := tools.FetchWebPageParams{
		URL:         "https://techcrunch.com",
		ExtractText: true,
	}

	fetchResult, err := fetchTool.Execute(ctx, fetchParams)
	if err != nil {
		log.Printf("Error fetching page: %v", err)
		return
	}

	webResult, ok := fetchResult.(*tools.FetchWebPageResult)
	if !ok {
		return
	}

	fmt.Printf("Analyzing: %s\n", webResult.Title)
	fmt.Printf("Status: %d\n", webResult.StatusCode)

	// Extract links including media
	linkParams := tools.ExtractLinksParams{
		URL:             webResult.URL,
		IncludeExternal: true,
		IncludeInternal: true,
		IncludeEmail:    true,
		IncludeMedia:    true,
	}

	linkResult, err := linkTool.Execute(ctx, linkParams)
	if err != nil {
		log.Printf("Error extracting links: %v", err)
		return
	}

	links, ok := linkResult.(*tools.ExtractLinksResult)
	if !ok {
		return
	}

	// Analyze link patterns
	fmt.Printf("\nLink Analysis:\n")
	fmt.Printf("- Total links: %d\n", links.TotalLinks)
	fmt.Printf("- Internal navigation: %d\n", len(links.InternalLinks))
	fmt.Printf("- External references: %d\n", len(links.ExternalLinks))
	fmt.Printf("- Contact emails: %d\n", len(links.EmailLinks))
	fmt.Printf("- Images: %d\n", len(links.MediaLinks))

	// Find article links (usually contain dates in URL)
	articleCount := 0
	for _, link := range links.InternalLinks {
		if strings.Contains(link.URL, "/2024/") || strings.Contains(link.URL, "/2025/") {
			articleCount++
		}
	}
	fmt.Printf("- Potential articles: %d\n", articleCount)

	// Show some image URLs
	if len(links.MediaLinks) > 0 {
		fmt.Println("\nSample images found:")
		for i, media := range links.MediaLinks {
			if i >= 3 {
				break
			}
			fmt.Printf("  - %s\n", media.URL)
			if media.Alt != "" {
				fmt.Printf("    Alt: %s\n", media.Alt)
			}
		}
	}
}

// Helper function to extract domain from URL
func extractDomain(urlStr string) string {
	parts := strings.Split(urlStr, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return urlStr
}

func extractMetadataExample(ctx context.Context) {
	tool := tools.NewExtractMetadataTool()

	// Try a popular news site that usually has good metadata
	params := tools.ExtractMetadataParams{
		URL: "https://www.bbc.com/news",
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Printf("Error extracting metadata: %v", err)
		return
	}

	if metaResult, ok := result.(*tools.ExtractMetadataResult); ok {
		fmt.Printf("Page: %s\n", metaResult.URL)
		fmt.Printf("Title: %s\n", metaResult.Title)

		if metaResult.Description != "" {
			fmt.Printf("Description: %s\n", metaResult.Description)
		}

		if metaResult.Author != "" {
			fmt.Printf("Author: %s\n", metaResult.Author)
		}

		// Show Open Graph data
		if len(metaResult.OpenGraph) > 0 {
			fmt.Println("\nOpen Graph metadata:")
			for key, value := range metaResult.OpenGraph {
				fmt.Printf("  og:%s = %s\n", key, value)
			}
		}

		// Show Twitter Card data
		if len(metaResult.Twitter) > 0 {
			fmt.Println("\nTwitter Card metadata:")
			for key, value := range metaResult.Twitter {
				fmt.Printf("  twitter:%s = %s\n", key, value)
			}
		}

		// Show keywords if any
		if len(metaResult.Keywords) > 0 {
			fmt.Printf("\nKeywords: %s\n", strings.Join(metaResult.Keywords, ", "))
		}

		// Show other interesting meta tags
		if len(metaResult.Meta) > 0 {
			fmt.Println("\nOther metadata:")
			for key, value := range metaResult.Meta {
				if key == "viewport" || key == "robots" || key == "theme-color" {
					fmt.Printf("  %s = %s\n", key, value)
				}
			}
		}
	}
}

func checkURLStatusExample(ctx context.Context) {
	tool := tools.NewCheckURLStatusTool()

	// Check various URLs
	urls := []string{
		"https://www.google.com",
		"https://httpstat.us/404", // Always returns 404
		"https://httpstat.us/301", // Always redirects
		"https://github.com",
		"https://example.com/nonexistent-page",
	}

	fmt.Println("Checking multiple URLs:")
	fmt.Println()

	for _, url := range urls {
		params := tools.CheckURLStatusParams{
			URL:             url,
			FollowRedirects: true,
			Timeout:         5,
		}

		result, err := tool.Execute(ctx, params)
		if err != nil {
			fmt.Printf("❌ %s\n", url)
			fmt.Printf("   Error: %v\n", err)
			continue
		}

		if statusResult, ok := result.(*tools.CheckURLStatusResult); ok {
			// Choose icon based on status
			icon := "✅"
			if statusResult.Status == "dead" {
				icon = "❌"
			} else if statusResult.Status == "redirect" {
				icon = "↪️"
			}

			fmt.Printf("%s %s\n", icon, url)
			fmt.Printf("   Status: %s (HTTP %d)\n", statusResult.Status, statusResult.StatusCode)
			fmt.Printf("   Response time: %dms\n", statusResult.ResponseTime)

			if statusResult.ContentType != "" {
				fmt.Printf("   Content-Type: %s\n", statusResult.ContentType)
			}

			if statusResult.ContentLength > 0 {
				fmt.Printf("   Content-Length: %d bytes\n", statusResult.ContentLength)
			}

			if statusResult.FinalURL != "" {
				fmt.Printf("   Final URL: %s\n", statusResult.FinalURL)
			}

			if len(statusResult.RedirectChain) > 0 {
				fmt.Printf("   Redirect chain: %d hops\n", len(statusResult.RedirectChain))
			}

			fmt.Println()
		}
	}

	// Example with custom headers
	fmt.Println("Checking URL with custom headers:")

	params := tools.CheckURLStatusParams{
		URL: "https://api.github.com",
		Headers: map[string]string{
			"Accept": "application/vnd.github.v3+json",
		},
		FollowRedirects: false,
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Printf("Error checking API URL: %v", err)
		return
	}

	if statusResult, ok := result.(*tools.CheckURLStatusResult); ok {
		fmt.Printf("API endpoint: %s\n", statusResult.URL)
		fmt.Printf("Status: %s (HTTP %d)\n", statusResult.Status, statusResult.StatusCode)
		fmt.Printf("Response time: %dms\n", statusResult.ResponseTime)

		// Show some response headers
		if xRateLimit, ok := statusResult.Headers["X-Ratelimit-Limit"]; ok {
			fmt.Printf("Rate limit: %s\n", xRateLimit)
		}
	}
}
