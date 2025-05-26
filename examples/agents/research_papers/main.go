// ABOUTME: Command-line interface for the research_papers agent that searches academic papers.
// ABOUTME: Supports multiple output formats and configurable LLM providers.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lexlapax/go-flock/pkg/agents"
	"github.com/lexlapax/go-flock/pkg/common"
	ldomain "github.com/lexlapax/go-llms/pkg/llm/domain"
	"github.com/lexlapax/go-llms/pkg/llm/provider"
)

func main() {
	// Define command-line flags
	var (
		query        = flag.String("query", "", "Research query (required)")
		format       = flag.String("format", "markdown", "Output format: markdown, json, or text")
		model        = flag.String("model", "", "LLM model to use (optional, uses provider default if not specified)")
		providerName = flag.String("provider", "", "LLM provider: openai, anthropic, or gemini (uses environment default if not specified)")
		output       = flag.String("output", "", "Output file (optional, prints to stdout if not specified)")
		debug        = flag.Bool("debug", false, "Enable debug logging")
		help         = flag.Bool("help", false, "Show help message")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Research Papers Agent - Academic Paper Search Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -query \"your research topic\" [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -query \"deep learning medical imaging\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -query \"climate change impacts\" -format json -output results.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -query \"quantum computing algorithms\" -provider openai -model gpt-4\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  OPENAI_API_KEY     - API key for OpenAI\n")
		fmt.Fprintf(os.Stderr, "  ANTHROPIC_API_KEY  - API key for Anthropic\n")
		fmt.Fprintf(os.Stderr, "  GEMINI_API_KEY     - API key for Google Gemini\n")
		fmt.Fprintf(os.Stderr, "  BRAVE_API_KEY      - API key for Brave Search (for ResearchPaperAPI tool)\n")
	}

	flag.Parse()

	// Initialize logging based on debug flag
	common.InitLogger(*debug)
	logger := common.GetLogger()

	// Show help if requested
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Validate required parameters
	if *query == "" {
		fmt.Fprintf(os.Stderr, "Error: -query is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Validate output format
	var outputFormat agents.OutputFormat
	switch strings.ToLower(*format) {
	case "markdown", "md":
		outputFormat = agents.OutputFormatMarkdown
	case "json":
		outputFormat = agents.OutputFormatJSON
	case "text", "txt":
		outputFormat = agents.OutputFormatText
	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid format '%s'. Must be markdown, json, or text\n", *format)
		os.Exit(1)
	}

	// Create LLM provider
	ctx := context.Background()
	logger.Debug(ctx, "Creating LLM provider: %s", *providerName)
	llmProvider, err := createProvider(*providerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating LLM provider: %v\n", err)
		os.Exit(1)
	}
	logger.Debug(ctx, "LLM provider created successfully")

	// Create agent options
	agentOpts := agents.AgentOptions{
		OutputFormat: outputFormat,
	}
	if *model != "" {
		agentOpts.Model = *model
	}

	// Create the research papers agent
	logger.Debug(ctx, "Creating research papers agent with format: %s, model: %s", outputFormat, agentOpts.Model)
	agent := agents.NewResearchPapersAgent(llmProvider, agentOpts)

	// Execute the search
	fmt.Fprintf(os.Stderr, "Searching for: %s\n", *query)
	fmt.Fprintf(os.Stderr, "Output format: %s\n", outputFormat)
	if *model != "" {
		fmt.Fprintf(os.Stderr, "Using model: %s\n", *model)
	}
	fmt.Fprintf(os.Stderr, "\nProcessing...\n\n")

	logger.Debug(ctx, "Running agent with query: %s", *query)
	result, err := agent.Run(ctx, *query)
	if err != nil {
		logger.Error(ctx, "Agent execution failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error running agent: %v\n", err)
		os.Exit(1)
	}
	logger.Debug(ctx, "Agent execution completed successfully")

	// Convert result to string
	output_content, ok := result.(string)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: Unexpected result type: %T\n", result)
		os.Exit(1)
	}

	// Output results
	if *output != "" {
		// Write to file
		err = os.WriteFile(*output, []byte(output_content), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Results saved to: %s\n", *output)
	} else {
		// Print to stdout
		fmt.Println(output_content)
	}
}

// createProvider creates an LLM provider based on the name or environment
func createProvider(providerName string) (ldomain.Provider, error) {
	// If no provider specified, try to detect from environment
	if providerName == "" {
		if os.Getenv("OPENAI_API_KEY") != "" {
			providerName = "openai"
		} else if os.Getenv("ANTHROPIC_API_KEY") != "" {
			providerName = "anthropic"
		} else if os.Getenv("GEMINI_API_KEY") != "" {
			providerName = "gemini"
		} else {
			return nil, fmt.Errorf("no LLM provider specified and no API keys found in environment")
		}
	}

	// Create provider based on name
	switch strings.ToLower(providerName) {
	case "openai":
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
		}
		return provider.NewOpenAIProvider(apiKey, "gpt-4"), nil

	case "anthropic":
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
		}
		return provider.NewAnthropicProvider(apiKey, "claude-3-5-sonnet-20241022"), nil

	case "gemini":
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
		}
		return provider.NewGeminiProvider(apiKey, "gemini-1.5-flash"), nil

	default:
		return nil, fmt.Errorf("unknown provider: %s (supported: openai, anthropic, gemini)", providerName)
	}
}
