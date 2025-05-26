// ABOUTME: Example application demonstrating the SearchResearch tool from go-flock
// ABOUTME: Shows how to perform research searches using multiple search providers

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lexlapax/go-flock/pkg/tools"
)

func main() {
	fmt.Println("===========================================")
	fmt.Println("   SearchResearch Tool Example")
	fmt.Println("===========================================")
	fmt.Println()

	// Create SearchResearch tool
	searchTool := tools.NewSearchResearchTool()

	// Display tool information
	fmt.Printf("Tool: %s\n", searchTool.Name())
	fmt.Printf("Description: %s\n", searchTool.Description())
	fmt.Println()

	// Check for API keys
	braveKey := os.Getenv("BRAVE_API_KEY")
	newsKey := os.Getenv("NEWS_API_KEY")

	if braveKey == "" && newsKey == "" {
		fmt.Println("⚠️  No API keys found. Running in demo mode...")
		fmt.Println()
		fmt.Println("To use real search functionality, set one or both:")
		fmt.Println("  export BRAVE_API_KEY=your_brave_api_key")
		fmt.Println("  export NEWS_API_KEY=your_newsapi_key")
		fmt.Println()
		fmt.Println("Get your API keys from:")
		fmt.Println("  Brave Search: https://brave.com/search/api/")
		fmt.Println("  News API: https://newsapi.org/")
		fmt.Println()
		runDemoMode()
		return
	}

	// Show which providers are available
	fmt.Println("Available search providers:")
	if braveKey != "" {
		fmt.Println("  ✓ Brave Search")
	}
	if newsKey != "" {
		fmt.Println("  ✓ News API")
	}
	fmt.Println()

	ctx := context.Background()

	// Example 1: Basic research search
	fmt.Println("Example 1: Basic Research Search")
	fmt.Println("--------------------------------")
	
	result1, err := searchTool.Execute(ctx, map[string]interface{}{
		"query": "artificial intelligence breakthroughs 2024",
	})
	if err != nil {
		log.Printf("Error in example 1: %v\n", err)
	} else {
		if result, ok := result1.(*tools.SearchResearchResult); ok {
			fmt.Printf("Found %d papers from %d providers\n\n", len(result.Papers), len(result.Providers))
			for i, paper := range result.Papers {
				if i >= 5 { // Show first 5 papers
					break
				}
				fmt.Printf("%d. %s\n", i+1, paper.Title)
				fmt.Printf("   Authors: %s\n", strings.Join(paper.Authors, ", "))
				fmt.Printf("   Source: %s\n", paper.Source)
				fmt.Printf("   URL: %s\n", paper.URL)
				if paper.PublishedDate != "" {
					fmt.Printf("   Published: %s\n", paper.PublishedDate)
				}
				fmt.Printf("   %s\n", paper.Abstract[:min(200, len(paper.Abstract))])
				fmt.Println()
			}
		}
	}

	// Example 2: Research with multiple sources
	fmt.Println("\nExample 2: Multi-Source Research")
	fmt.Println("--------------------------------")
	
	result2, err := searchTool.Execute(ctx, map[string]interface{}{
		"query": "climate change latest research",
		"max_results": 10,
	})
	if err != nil {
		log.Printf("Error in example 2: %v\n", err)
	} else {
		if result, ok := result2.(*tools.SearchResearchResult); ok {
			// Show provider info
			fmt.Println("Results by provider:")
			for _, provider := range result.Providers {
				fmt.Printf("  %s: %d papers", provider.Name, provider.ResultCount)
				if provider.Error != "" {
					fmt.Printf(" (Error: %s)", provider.Error)
				}
				fmt.Printf(" [%dms]\n", provider.ResponseTime)
			}
			fmt.Println()
			
			// Show recent papers
			fmt.Println("Most recent papers:")
			shown := 0
			for _, paper := range result.Papers {
				if shown >= 3 {
					break
				}
				if paper.PublishedDate != "" {
					fmt.Printf("- %s (%s)\n", paper.Title, paper.PublishedDate)
					fmt.Printf("  %s\n", paper.URL)
					shown++
				}
			}
		}
	}

	// Example 3: Domain-specific research
	fmt.Println("\n\nExample 3: Domain-Specific Research")
	fmt.Println("-----------------------------------")
	
	result3, err := searchTool.Execute(ctx, map[string]interface{}{
		"query": "machine learning healthcare applications",
		"max_results": 8,
	})
	if err != nil {
		log.Printf("Error in example 3: %v\n", err)
	} else {
		if result, ok := result3.(*tools.SearchResearchResult); ok {
			fmt.Printf("Found %d papers about ML in healthcare\n\n", len(result.Papers))
			
			// Categorize by source
			sourceCount := make(map[string]int)
			for _, paper := range result.Papers {
				sourceCount[paper.Source]++
			}
			
			fmt.Println("Papers by source:")
			for source, count := range sourceCount {
				fmt.Printf("  %s: %d papers\n", source, count)
			}
			fmt.Println()
			
			// Show sample papers
			fmt.Println("Sample papers:")
			for i, paper := range result.Papers {
				if i >= 4 {
					break
				}
				fmt.Printf("\n%d. %s\n", i+1, paper.Title)
				fmt.Printf("   Authors: %s\n", strings.Join(paper.Authors[:min(3, len(paper.Authors))], ", "))
				if paper.Journal != "" {
					fmt.Printf("   Journal: %s\n", paper.Journal)
				}
				fmt.Printf("   %s\n", paper.Abstract[:min(150, len(paper.Abstract))]+"...")
			}
		}
	}

	// Example 4: Time-sensitive research
	fmt.Println("\n\nExample 4: Time-Sensitive Research")
	fmt.Println("----------------------------------")
	
	result4, err := searchTool.Execute(ctx, map[string]interface{}{
		"query": "quantum computing news",
		"max_results": 15,
	})
	if err != nil {
		log.Printf("Error in example 4: %v\n", err)
	} else {
		if result, ok := result4.(*tools.SearchResearchResult); ok {
			fmt.Printf("Analyzing %d papers for recency...\n\n", len(result.Papers))
			
			// Separate by time period (simple year-based check)
			currentYear := time.Now().Year()
			thisYear := 0
			lastYear := 0
			older := 0
			
			for _, paper := range result.Papers {
				if paper.PublishedDate != "" {
					if contains(paper.PublishedDate, []string{fmt.Sprintf("%d", currentYear)}) {
						thisYear++
					} else if contains(paper.PublishedDate, []string{fmt.Sprintf("%d", currentYear-1)}) {
						lastYear++
					} else {
						older++
					}
				} else {
					older++
				}
			}
			
			fmt.Println("Temporal distribution:")
			fmt.Printf("  This year (%d): %d\n", currentYear, thisYear)
			fmt.Printf("  Last year (%d): %d\n", currentYear-1, lastYear)
			fmt.Printf("  Older/Undated: %d\n", older)
			fmt.Println()
			
			// Show most recent
			fmt.Println("Latest papers:")
			shown := 0
			for _, paper := range result.Papers {
				if shown >= 3 {
					break
				}
				fmt.Printf("- %s\n", paper.Title)
				if paper.PublishedDate != "" {
					fmt.Printf("  Date: %s\n", paper.PublishedDate)
				}
				fmt.Printf("  %s\n", paper.URL)
				shown++
				fmt.Println()
			}
		}
	}

	fmt.Println("\n===========================================")
	fmt.Println("Research search examples completed!")
	fmt.Println("===========================================")
}

func runDemoMode() {
	fmt.Println("===========================================")
	fmt.Println("         DEMO MODE")
	fmt.Println("===========================================")
	fmt.Println()

	// Simulate research results
	mockPapers := []tools.ResearchPaper{
		{
			Title:         "Deep Learning for Medical Image Analysis: A Comprehensive Survey",
			Authors:       []string{"Smith, J.", "Johnson, L.", "Chen, X."},
			Abstract:      "This comprehensive survey examines the latest advances in deep learning techniques for medical image analysis, covering applications in radiology, pathology, and surgical planning. We review over 300 papers published between 2020-2024...",
			PublishedDate: "2024-01-15",
			Source:        "arXiv",
			URL:           "https://arxiv.org/abs/2401.12345",
			PDFURL:        "https://arxiv.org/pdf/2401.12345.pdf",
			ArxivID:       "2401.12345",
			Journal:       "arXiv preprint",
		},
		{
			Title:         "Climate Change Impacts on Global Food Security: A Machine Learning Analysis",
			Authors:       []string{"Anderson, R.", "Lee, K.", "Martinez, P.", "Wang, H."},
			Abstract:      "Using advanced machine learning models, we analyze the projected impacts of climate change on global food security through 2050. Our analysis integrates data from 150 countries...",
			PublishedDate: "2024-01-14",
			Source:        "PubMed",
			URL:           "https://pubmed.ncbi.nlm.nih.gov/38234567/",
			PubMedID:      "38234567",
			Journal:       "Nature Climate Change",
		},
		{
			Title:         "Quantum Computing Algorithms for Drug Discovery",
			Authors:       []string{"Thompson, D.", "Patel, S."},
			Abstract:      "We present novel quantum computing algorithms specifically designed for molecular simulation in drug discovery. Our approach reduces computational complexity from exponential to polynomial time...",
			PublishedDate: "2024-01-13",
			Source:        "arXiv",
			URL:           "https://arxiv.org/abs/2401.11111",
			PDFURL:        "https://arxiv.org/pdf/2401.11111.pdf",
			ArxivID:       "2401.11111",
		},
		{
			Title:         "AI Ethics in Healthcare: A Systematic Review",
			Authors:       []string{"Williams, M.", "Davis, J.", "Garcia, A."},
			Abstract:      "This systematic review examines ethical considerations in the deployment of artificial intelligence systems in healthcare settings. We analyze 200 case studies and propose a comprehensive ethical framework...",
			PublishedDate: "2024-01-12",
			Source:        "PubMed",
			URL:           "https://pubmed.ncbi.nlm.nih.gov/38234568/",
			PubMedID:      "38234568",
			Journal:       "Journal of Medical Ethics",
		},
		{
			Title:         "Advances in Natural Language Processing for Scientific Literature",
			Authors:       []string{"Brown, T.", "Miller, E.", "Zhang, Y."},
			Abstract:      "We present state-of-the-art NLP techniques for automated analysis of scientific literature, achieving 95% accuracy in paper classification and 89% in key finding extraction...",
			PublishedDate: "2024-01-11",
			Source:        "CORE",
			URL:           "https://core.ac.uk/display/123456789",
			DOI:           "10.1234/example.2024.001",
		},
	}

	fmt.Println("Example: Multi-Source Research Results")
	fmt.Println("-------------------------------------")
	fmt.Printf("Query: \"artificial intelligence breakthroughs\"\n")
	fmt.Printf("Found %d papers from multiple sources\n\n", len(mockPapers))

	// Show papers grouped by source
	sourceCount := make(map[string]int)
	for _, paper := range mockPapers {
		sourceCount[paper.Source]++
	}

	fmt.Println("Papers by source:")
	for source, count := range sourceCount {
		fmt.Printf("  %s: %d papers\n", source, count)
	}
	fmt.Println()

	fmt.Println("Sample papers:")
	for i, paper := range mockPapers {
		if i >= 3 { // Show first 3
			break
		}
		fmt.Printf("\n%d. %s\n", i+1, paper.Title)
		fmt.Printf("   Authors: %s\n", strings.Join(paper.Authors, ", "))
		fmt.Printf("   Source: %s\n", paper.Source)
		fmt.Printf("   Published: %s\n", paper.PublishedDate)
		if paper.Journal != "" {
			fmt.Printf("   Journal: %s\n", paper.Journal)
		}
		fmt.Printf("   URL: %s\n", paper.URL)
		if paper.PDFURL != "" {
			fmt.Printf("   PDF: %s\n", paper.PDFURL)
		}
		fmt.Printf("   Abstract: %s...\n", paper.Abstract[:min(150, len(paper.Abstract))])
	}

	fmt.Println("\n\nDEMO: Provider Information")
	fmt.Println("-------------------------")
	fmt.Println("Simulated provider response times:")
	fmt.Println("  arXiv: 2 papers [320ms]")
	fmt.Println("  PubMed: 2 papers [450ms]")
	fmt.Println("  CORE: 1 paper [280ms]")
	
	fmt.Println("\nDEMO: Paper Timeline")
	fmt.Println("-------------------")
	fmt.Println("Papers by date:")
	for _, paper := range mockPapers {
		fmt.Printf("  %s - %s\n", paper.PublishedDate, paper.Title[:min(60, len(paper.Title))]+"...")
	}

	fmt.Println("\n===========================================")
	fmt.Println("Demo mode completed. Set API keys to search real data!")
	fmt.Println("===========================================")
}

// Helper function to check if string contains any of the substrings
func contains(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// Helper function for minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}