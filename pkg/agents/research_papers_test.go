// ABOUTME: Test file for the research_papers agent that searches academic papers and research.
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

func TestNewResearchPapersAgent(t *testing.T) {
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
			// Create mock provider
			mockProvider := provider.NewMockProvider()

			// Create agent
			agent := NewResearchPapersAgent(mockProvider, tt.options)

			if agent == nil {
				t.Error("NewResearchPapersAgent returned nil")
			}
		})
	}
}

func TestResearchPapersAgentMarkdownOutput(t *testing.T) {
	// Create mock provider with predefined response
	mockProvider := provider.NewMockProvider()
	mockProvider.WithGenerateMessageFunc(func(ctx context.Context, messages []ldomain.Message, options ...ldomain.Option) (ldomain.Response, error) {
		return ldomain.Response{
			Content: `# Research Findings: AI in Healthcare

## Executive Summary
Recent research shows significant advances in AI applications for healthcare diagnostics, with deep learning models achieving expert-level performance in several domains.

## Key Papers

### Deep Learning in Medical Imaging
- **Title**: "Deep Learning for Medical Image Analysis: A Comprehensive Survey"
- **Authors**: Smith, J., Johnson, L., Chen, X.
- **Year**: 2024
- **Summary**: Comprehensive survey of deep learning techniques in radiology, pathology, and surgical planning.
- **URL**: https://arxiv.org/abs/2401.12345

### AI for Drug Discovery
- **Title**: "Machine Learning Accelerates Drug Discovery Pipeline"
- **Authors**: Anderson, R., Lee, K.
- **Year**: 2024
- **Summary**: Novel ML approaches reduce drug discovery time by 40%.
- **URL**: https://pubmed.ncbi.nlm.nih.gov/38234567/

## Research Timeline
- 2020: Initial breakthroughs in medical image classification
- 2022: FDA approval of first AI diagnostic systems
- 2024: Widespread adoption in clinical settings

## Key Researchers
- Dr. Jane Smith (Stanford) - Medical imaging AI
- Prof. Robert Anderson (MIT) - Drug discovery algorithms

## References
1. Smith, J., et al. (2024). Deep Learning for Medical Image Analysis. arXiv:2401.12345
2. Anderson, R., Lee, K. (2024). Machine Learning Accelerates Drug Discovery. Nature Medicine.`,
		}, nil
	})

	// Create agent with markdown output
	agent := NewResearchPapersAgent(mockProvider, AgentOptions{
		OutputFormat: OutputFormatMarkdown,
	})

	// Execute search
	ctx := context.Background()
	result, err := agent.Run(ctx, "AI applications in healthcare diagnostics")

	if err != nil {
		t.Fatalf("agent.Run returned error: %v", err)
	}

	// Check result is string
	output, ok := result.(string)
	if !ok {
		t.Fatalf("expected string result, got %T", result)
	}

	// Verify markdown structure
	if !strings.Contains(output, "# Research Findings") {
		t.Error("output missing markdown header")
	}
	if !strings.Contains(output, "## Executive Summary") {
		t.Error("output missing executive summary section")
	}
	if !strings.Contains(output, "## Key Papers") {
		t.Error("output missing key papers section")
	}
	if !strings.Contains(output, "**Title**:") {
		t.Error("output missing paper titles")
	}
	if !strings.Contains(output, "## References") {
		t.Error("output missing references section")
	}
}

func TestResearchPapersAgentJSONOutput(t *testing.T) {
	// Create mock provider with JSON response
	mockProvider := provider.NewMockProvider()
	mockProvider.WithGenerateMessageFunc(func(ctx context.Context, messages []ldomain.Message, options ...ldomain.Option) (ldomain.Response, error) {
		return ldomain.Response{
			Content: `{
  "papers": [
    {
      "title": "Deep Learning for Medical Image Analysis: A Comprehensive Survey",
      "authors": ["Smith, J.", "Johnson, L.", "Chen, X."],
      "year": "2024",
      "abstract": "Comprehensive survey of deep learning techniques in radiology, pathology, and surgical planning.",
      "url": "https://arxiv.org/abs/2401.12345",
      "citations": 45
    },
    {
      "title": "Machine Learning Accelerates Drug Discovery Pipeline",
      "authors": ["Anderson, R.", "Lee, K."],
      "year": "2024",
      "abstract": "Novel ML approaches reduce drug discovery time by 40%.",
      "url": "https://pubmed.ncbi.nlm.nih.gov/38234567/",
      "citations": 23
    }
  ],
  "themes": [
    "Medical imaging analysis",
    "Drug discovery optimization",
    "Clinical decision support"
  ],
  "key_authors": [
    "Smith, J. (Stanford)",
    "Anderson, R. (MIT)"
  ],
  "timeline": "2020-2024: Rapid adoption of AI in clinical settings",
  "summary": "AI is transforming healthcare through improved diagnostics and accelerated drug discovery."
}`,
		}, nil
	})

	// Create agent with JSON output
	agent := NewResearchPapersAgent(mockProvider, AgentOptions{
		OutputFormat: OutputFormatJSON,
	})

	// Execute search
	ctx := context.Background()
	result, err := agent.Run(ctx, "AI applications in healthcare diagnostics")

	if err != nil {
		t.Fatalf("agent.Run returned error: %v", err)
	}

	// Check result is string
	output, ok := result.(string)
	if !ok {
		t.Fatalf("expected string result, got %T", result)
	}

	// Parse JSON to verify structure
	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(output), &jsonResult); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	// Verify JSON structure
	if _, ok := jsonResult["papers"]; !ok {
		t.Error("JSON missing 'papers' field")
	}
	if _, ok := jsonResult["themes"]; !ok {
		t.Error("JSON missing 'themes' field")
	}
	if _, ok := jsonResult["key_authors"]; !ok {
		t.Error("JSON missing 'key_authors' field")
	}
	if _, ok := jsonResult["timeline"]; !ok {
		t.Error("JSON missing 'timeline' field")
	}
	if _, ok := jsonResult["summary"]; !ok {
		t.Error("JSON missing 'summary' field")
	}

	// Verify papers array
	papers, ok := jsonResult["papers"].([]interface{})
	if !ok || len(papers) != 2 {
		t.Error("JSON papers field should be array with 2 items")
	}
}

func TestResearchPapersAgentTextOutput(t *testing.T) {
	// Create mock provider with plain text response
	mockProvider := provider.NewMockProvider()
	mockProvider.WithGenerateMessageFunc(func(ctx context.Context, messages []ldomain.Message, options ...ldomain.Option) (ldomain.Response, error) {
		return ldomain.Response{
			Content: `RESEARCH FINDINGS: AI IN HEALTHCARE

SUMMARY
Recent research shows significant advances in AI applications for healthcare diagnostics. Deep learning models are achieving expert-level performance in medical imaging and drug discovery.

KEY PAPERS
1. Deep Learning for Medical Image Analysis: A Comprehensive Survey
   Authors: Smith, J., Johnson, L., Chen, X. (2024)
   Summary: Survey of deep learning in radiology and pathology
   Link: https://arxiv.org/abs/2401.12345

2. Machine Learning Accelerates Drug Discovery Pipeline
   Authors: Anderson, R., Lee, K. (2024)
   Summary: ML reduces drug discovery time by 40%
   Link: https://pubmed.ncbi.nlm.nih.gov/38234567/

RESEARCH THEMES
- Medical imaging analysis
- Drug discovery optimization
- Clinical decision support

KEY RESEARCHERS
- Dr. Jane Smith at Stanford University
- Prof. Robert Anderson at MIT

TIMELINE
- 2020: Initial breakthroughs in medical image classification
- 2022: FDA approval of first AI diagnostic systems
- 2024: Widespread adoption in clinical settings`,
		}, nil
	})

	// Create agent with text output
	agent := NewResearchPapersAgent(mockProvider, AgentOptions{
		OutputFormat: OutputFormatText,
	})

	// Execute search
	ctx := context.Background()
	result, err := agent.Run(ctx, "AI applications in healthcare diagnostics")

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
	if !strings.Contains(output, "RESEARCH FINDINGS") {
		t.Error("output missing research findings header")
	}
	if !strings.Contains(output, "KEY PAPERS") {
		t.Error("output missing key papers section")
	}
}

func TestResearchPapersAgentWithTools(t *testing.T) {
	// This test verifies that the agent has the required tools configured
	mockProvider := provider.NewMockProvider()
	mockProvider.WithGenerateMessageFunc(func(ctx context.Context, messages []ldomain.Message, options ...ldomain.Option) (ldomain.Response, error) {
		return ldomain.Response{Content: "# Test Response"}, nil
	})

	agent := NewResearchPapersAgent(mockProvider, DefaultAgentOptions())

	// Since we can't directly inspect tools on the agent interface,
	// we'll verify through the system prompt that tools are mentioned
	// The actual tool testing will be done in integration tests

	// For now, just verify the agent can be created and run
	ctx := context.Background()
	_, err := agent.Run(ctx, "test query")

	if err != nil {
		t.Errorf("agent with tools should execute successfully: %v", err)
	}
}

func TestResearchPapersAgentErrorHandling(t *testing.T) {
	// Create mock provider that returns an error
	mockProvider := provider.NewMockProvider()
	mockProvider.WithGenerateMessageFunc(func(ctx context.Context, messages []ldomain.Message, options ...ldomain.Option) (ldomain.Response, error) {
		return ldomain.Response{}, fmt.Errorf("API rate limit exceeded")
	})

	agent := NewResearchPapersAgent(mockProvider, DefaultAgentOptions())

	ctx := context.Background()
	_, err := agent.Run(ctx, "test query")

	if err == nil {
		t.Error("expected error from provider")
	}
	if !strings.Contains(err.Error(), "API rate limit exceeded") {
		t.Errorf("expected error message to contain 'API rate limit exceeded', got: %v", err)
	}
}
