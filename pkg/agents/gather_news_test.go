// ABOUTME: Test file for the gather_news agent that searches current news and events.
// ABOUTME: Tests cover markdown, JSON, and text output formats as well as error handling.

package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	ldomain "github.com/lexlapax/go-llms/pkg/llm/domain"
	"github.com/lexlapax/go-llms/pkg/llm/provider"
)

func TestNewGatherNewsAgent(t *testing.T) {
	tests := []struct {
		name    string
		options AgentOptions
		wantErr bool
	}{
		{
			name:    "default options",
			options: DefaultAgentOptions(),
			wantErr: false,
		},
		{
			name: "json output format",
			options: AgentOptions{
				OutputFormat: OutputFormatJSON,
			},
			wantErr: false,
		},
		{
			name: "text output format",
			options: AgentOptions{
				OutputFormat: OutputFormatText,
			},
			wantErr: false,
		},
		{
			name: "with custom model",
			options: AgentOptions{
				OutputFormat: OutputFormatMarkdown,
				Model:        "gpt-4",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := provider.NewMockProvider()
			agent := NewGatherNewsAgent(mockProvider, tt.options)
			if agent == nil {
				t.Error("Expected agent to be created, got nil")
			}
		})
	}
}

func TestGatherNewsPromptGeneration(t *testing.T) {
	tests := []struct {
		name            string
		format          OutputFormat
		expectedParts   []string
		unexpectedParts []string
	}{
		{
			name:   "markdown format includes markdown instructions",
			format: OutputFormatMarkdown,
			expectedParts: []string{
				"You are a news research specialist",
				"search_news_api",
				"search_web_brave",
				"# News Analysis:",
				"## Executive Summary",
				"## Major Stories",
			},
			unexpectedParts: []string{
				"JSON structure",
				"TOPIC IN CAPS",
			},
		},
		{
			name:   "json format includes json schema",
			format: OutputFormatJSON,
			expectedParts: []string{
				"You are a news research specialist",
				"search_news_api",
				"JSON structure following this exact schema",
				`"articles": [`,
				`"trends": ["string"]`,
			},
			unexpectedParts: []string{
				"# News Analysis:",
				"TOPIC IN CAPS",
			},
		},
		{
			name:   "text format includes plain text instructions",
			format: OutputFormatText,
			expectedParts: []string{
				"You are a news research specialist",
				"search_news_api",
				"NEWS ANALYSIS: [TOPIC IN CAPS]",
				"EXECUTIVE SUMMARY",
				"MAJOR STORIES",
			},
			unexpectedParts: []string{
				"# News Analysis:",
				"JSON structure",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := getGatherNewsPrompt(tt.format)

			for _, expected := range tt.expectedParts {
				if !strings.Contains(prompt, expected) {
					t.Errorf("Prompt should contain '%s' but it doesn't", expected)
				}
			}

			for _, unexpected := range tt.unexpectedParts {
				if strings.Contains(prompt, unexpected) {
					t.Errorf("Prompt should not contain '%s' but it does", unexpected)
				}
			}
		})
	}
}

func TestGatherNewsAgentMarkdownOutput(t *testing.T) {
	// Create mock provider
	mockProvider := provider.NewMockProvider()
	mockProvider.WithGenerateMessageFunc(func(ctx context.Context, messages []ldomain.Message, options ...ldomain.Option) (ldomain.Response, error) {
		return ldomain.Response{
			Content: `# News Analysis: Climate Change

## Executive Summary
Recent developments show significant progress in climate policy.

## Major Stories
### New Climate Agreement Reached
- **Source**: Reuters
- **Date**: 2024-01-20
- **Summary**: World leaders agree on emissions targets
- **URL**: https://example.com/climate-agreement

## Emerging Trends
- Increased renewable energy investments
- Carbon capture technology advances`,
		}, nil
	})

	// Create agent with markdown output
	agent := NewGatherNewsAgent(mockProvider, AgentOptions{
		OutputFormat: OutputFormatMarkdown,
	})

	// Execute search
	ctx := context.Background()
	result, err := agent.Run(ctx, "climate change latest news")

	if err != nil {
		t.Fatalf("agent.Run returned error: %v", err)
	}

	// Check result is string
	output, ok := result.(string)
	if !ok {
		t.Fatalf("expected string result, got %T", result)
	}

	// Verify markdown structure
	if !strings.Contains(output, "# News Analysis:") {
		t.Error("output missing markdown header")
	}
	if !strings.Contains(output, "## Executive Summary") {
		t.Error("output missing executive summary")
	}
	if !strings.Contains(output, "## Major Stories") {
		t.Error("output missing major stories section")
	}
}

func TestGatherNewsAgentJSONOutput(t *testing.T) {
	// Create mock provider
	mockProvider := provider.NewMockProvider()
	mockProvider.WithGenerateMessageFunc(func(ctx context.Context, messages []ldomain.Message, options ...ldomain.Option) (ldomain.Response, error) {
		return ldomain.Response{
			Content: `{
				"topic": "climate change",
				"summary": "Recent climate developments",
				"articles": [
					{
						"title": "New Climate Agreement",
						"source": "Reuters",
						"published_date": "2024-01-20",
						"summary": "World leaders agree on targets",
						"url": "https://example.com/climate"
					}
				],
				"trends": ["renewable energy", "carbon capture"],
				"analysis_date": "2024-01-20"
			}`,
		}, nil
	})

	// Create agent with JSON output
	agent := NewGatherNewsAgent(mockProvider, AgentOptions{
		OutputFormat: OutputFormatJSON,
	})

	// Execute search
	ctx := context.Background()
	result, err := agent.Run(ctx, "climate change")

	if err != nil {
		t.Fatalf("agent.Run returned error: %v", err)
	}

	// Check result is string
	output, ok := result.(string)
	if !ok {
		t.Fatalf("expected string result, got %T", result)
	}

	// Verify it's valid JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &jsonData); err != nil {
		t.Fatalf("output should be valid JSON: %v", err)
	}

	// Verify JSON structure
	if _, ok := jsonData["topic"]; !ok {
		t.Error("JSON missing 'topic' field")
	}
	if _, ok := jsonData["articles"]; !ok {
		t.Error("JSON missing 'articles' field")
	}
}

func TestGatherNewsAgentTextOutput(t *testing.T) {
	// Create mock provider
	mockProvider := provider.NewMockProvider()
	mockProvider.WithGenerateMessageFunc(func(ctx context.Context, messages []ldomain.Message, options ...ldomain.Option) (ldomain.Response, error) {
		return ldomain.Response{
			Content: `NEWS ANALYSIS: CLIMATE CHANGE

EXECUTIVE SUMMARY
Recent developments in climate policy show progress.

MAJOR STORIES
1. New Climate Agreement
   Source: Reuters (2024-01-20)
   Summary: World leaders agree on targets
   Link: https://example.com/climate

EMERGING TRENDS
- Renewable energy growth
- Carbon capture advances`,
		}, nil
	})

	// Create agent with text output
	agent := NewGatherNewsAgent(mockProvider, AgentOptions{
		OutputFormat: OutputFormatText,
	})

	// Execute search
	ctx := context.Background()
	result, err := agent.Run(ctx, "climate change")

	if err != nil {
		t.Fatalf("agent.Run returned error: %v", err)
	}

	// Check result is string
	output, ok := result.(string)
	if !ok {
		t.Fatalf("expected string result, got %T", result)
	}

	// Verify text structure (no markdown or JSON formatting)
	if strings.Contains(output, "#") {
		t.Error("text output should not contain markdown headers")
	}
	if strings.Contains(output, "{") || strings.Contains(output, "}") {
		t.Error("text output should not contain JSON formatting")
	}
	if !strings.Contains(output, "NEWS ANALYSIS") {
		t.Error("output missing news analysis header")
	}
	if !strings.Contains(output, "MAJOR STORIES") {
		t.Error("output missing major stories section")
	}
}

func TestGatherNewsAgentErrorHandling(t *testing.T) {
	// Create mock provider that returns an error
	mockProvider := provider.NewMockProvider()
	mockProvider.WithGenerateMessageFunc(func(ctx context.Context, messages []ldomain.Message, options ...ldomain.Option) (ldomain.Response, error) {
		return ldomain.Response{}, fmt.Errorf("API rate limit exceeded")
	})

	agent := NewGatherNewsAgent(mockProvider, DefaultAgentOptions())

	ctx := context.Background()
	_, err := agent.Run(ctx, "test query")

	if err == nil {
		t.Error("expected error from provider")
	}
	if !strings.Contains(err.Error(), "API rate limit exceeded") {
		t.Errorf("expected error message to contain 'API rate limit exceeded', got: %v", err)
	}
}

func TestGatherNewsPromptCoreInstructions(t *testing.T) {
	prompt := getGatherNewsPrompt(OutputFormatMarkdown)

	criticalInstructions := []string{
		"DO NOT generate placeholder articles",
		"use ONLY real results from tool calls",
		"WAIT for tool results",
		"Execute tools, don't describe them",
		"IMMEDIATELY call BOTH search_news_api AND search_web_brave",
		"arguments\" field MUST be a JSON string",
	}

	for _, instruction := range criticalInstructions {
		if !strings.Contains(prompt, instruction) {
			t.Errorf("Missing critical instruction: %s", instruction)
		}
	}
}

func TestGatherNewsPromptTools(t *testing.T) {
	prompt := getGatherNewsPrompt(OutputFormatMarkdown)

	expectedTools := []string{
		"search_news_api",
		"search_web_brave",
		"fetch_webpage",
		"extract_metadata",
	}

	for _, tool := range expectedTools {
		if !strings.Contains(prompt, tool) {
			t.Errorf("Prompt should mention tool: %s", tool)
		}
	}
}

// Benchmark tests
func BenchmarkNewGatherNewsAgent(b *testing.B) {
	mockProvider := provider.NewMockProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewGatherNewsAgent(mockProvider)
	}
}

func BenchmarkGatherNewsPromptGeneration(b *testing.B) {
	formats := []OutputFormat{
		OutputFormatMarkdown,
		OutputFormatJSON,
		OutputFormatText,
	}

	for _, format := range formats {
		b.Run(string(format), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = getGatherNewsPrompt(format)
			}
		})
	}
}
