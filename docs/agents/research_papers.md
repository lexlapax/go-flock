# Research Papers Agent

## Overview

The Research Papers Agent specializes in finding and analyzing academic papers from multiple research databases. It uses the ResearchPaperAPI tool to query arXiv, PubMed, and CORE databases simultaneously, providing comprehensive research findings.

## Features

- **Multi-source Search**: Queries arXiv, PubMed, and CORE databases
- **Flexible Output**: Supports Markdown (default), JSON, and plain text formats
- **Comprehensive Analysis**: Extracts papers, themes, key researchers, and timelines
- **Tool Integration**: Uses ResearchPaperAPI, FetchWebPage, and ExtractMetadata tools
- **Configurable Models**: Works with any go-llms supported LLM provider

## Usage

### As a Library

```go
import (
    "context"
    "github.com/lexlapax/go-flock/pkg/agents"
    "github.com/lexlapax/go-llms/pkg/llm/provider"
)

// Create provider
provider := provider.NewOpenAIProvider(apiKey, "gpt-4")

// Create agent with default markdown output
agent := agents.NewResearchPapersAgent(provider)

// Or with specific output format
agent := agents.NewResearchPapersAgent(provider, agents.AgentOptions{
    OutputFormat: agents.OutputFormatJSON,
    Model: "gpt-4",
})

// Execute search
ctx := context.Background()
result, err := agent.Run(ctx, "deep learning medical imaging")
```

### Command Line

The research papers agent includes a complete CLI example in `examples/agents/research_papers/`:

```bash
# Navigate to the example
cd examples/agents/research_papers

# Basic search with markdown output
go run main.go -query "deep learning medical imaging"

# JSON output saved to file
go run main.go -query "climate change impacts" -format json -output results.json

# Using specific provider and model
go run main.go -query "quantum computing" -provider openai -model gpt-4

# Plain text output
go run main.go -query "renewable energy" -format text
```

## Output Formats

### Markdown (Default)

Well-formatted report with:
- Executive summary
- Key papers with details
- Research timeline
- Notable researchers
- Formatted references

Example:
```markdown
# Research Findings: Deep Learning Medical Imaging

## Executive Summary
Recent advances in deep learning have revolutionized medical imaging...

## Key Papers

### Deep Learning for Medical Image Analysis
- **Authors**: Smith, J., Johnson, L.
- **Year**: 2024
- **Summary**: Comprehensive survey of DL techniques...
- **URL**: https://arxiv.org/abs/2401.12345
```

### JSON

Structured data format:
```json
{
  "papers": [
    {
      "title": "Deep Learning for Medical Image Analysis",
      "authors": ["Smith, J.", "Johnson, L."],
      "year": "2024",
      "abstract": "Comprehensive survey...",
      "url": "https://arxiv.org/abs/2401.12345",
      "citations": 45
    }
  ],
  "themes": ["Medical imaging", "Deep learning"],
  "key_authors": ["Smith, J. (Stanford)"],
  "timeline": "2020-2024: Rapid adoption",
  "summary": "Deep learning transforms medical imaging"
}
```

### Plain Text

Simple, readable format without markup:
```
RESEARCH FINDINGS: DEEP LEARNING MEDICAL IMAGING

SUMMARY
Recent advances in deep learning...

KEY PAPERS
1. Deep Learning for Medical Image Analysis
   Authors: Smith, J., Johnson, L. (2024)
   Summary: Comprehensive survey...
   Link: https://arxiv.org/abs/2401.12345
```

## Configuration

### Environment Variables

- `OPENAI_API_KEY` - For OpenAI provider
- `ANTHROPIC_API_KEY` - For Anthropic provider  
- `GEMINI_API_KEY` - For Google Gemini provider
- `BRAVE_API_KEY` - Required for ResearchPaperAPI tool

### Agent Options

```go
type AgentOptions struct {
    OutputFormat OutputFormat // markdown, json, or text
    Model        string      // LLM model name (optional)
}
```

## Tools Used

1. **ResearchPaperAPI** - Primary tool for querying academic databases
2. **FetchWebPage** - Retrieves full content from paper URLs
3. **ExtractMetadata** - Extracts structured metadata from sources

## Integration Examples

### Research Pipeline

```go
// Create research agent
researchAgent := agents.NewResearchPapersAgent(provider)

// Get research papers
papers, _ := researchAgent.Run(ctx, "machine learning healthcare")

// Pass to synthesis agent
synthesisAgent := agents.NewSynthesizeContentAgent(provider)
report, _ := synthesisAgent.Run(ctx, papers)
```

### Batch Processing

```go
queries := []string{
    "AI ethics healthcare",
    "federated learning medical data",
    "privacy preserving ML",
}

for _, query := range queries {
    result, err := agent.Run(ctx, query)
    if err != nil {
        log.Printf("Error for %s: %v", query, err)
        continue
    }
    saveResult(query, result)
}
```

### Web Service

```go
http.HandleFunc("/research", func(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    format := r.URL.Query().Get("format")
    
    opts := agents.AgentOptions{
        OutputFormat: agents.OutputFormatJSON,
    }
    if format != "" {
        opts.OutputFormat = agents.OutputFormat(format)
    }
    
    agent := agents.NewResearchPapersAgent(provider, opts)
    result, err := agent.Run(context.Background(), query)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(result.(string)))
})
```

## Prompt Architecture

The Research Papers Agent uses a modular prompt structure for maintainability:

### Core Prompt
The core prompt (shared across all formats) defines:
- Agent role and expertise
- Available tools and their purposes
- General analysis approach
- Quality standards for research

### Format-Specific Instructions
Each output format has specific formatting instructions that are appended to the core prompt:
- **Markdown**: Section structure, formatting guidelines
- **JSON**: Schema definition, field requirements
- **Text**: Plain text formatting rules

This separation allows you to:
- Modify the core behavior without changing format instructions
- Add new output formats easily
- Maintain consistency across formats
- Test prompts independently

## Best Practices

1. **API Keys**: Ensure BRAVE_API_KEY is set for the ResearchPaperAPI tool
2. **Specific Queries**: More specific queries yield better results
3. **Output Format**: Choose format based on downstream processing needs
4. **Error Handling**: Always check for errors, especially API rate limits
5. **Caching**: Consider caching results for repeated queries
6. **Prompt Tuning**: Modify the core prompt in `research_papers.go` to adjust agent behavior

## Limitations

- Requires active internet connection
- Subject to API rate limits
- Results depend on database coverage
- May not include very recent papers (database lag)

## Troubleshooting

### No Results Found
- Check if BRAVE_API_KEY is set
- Verify internet connectivity
- Try broader search terms

### API Errors
- Check API key validity
- Monitor rate limits
- Verify provider configuration

### Parsing Errors
- Ensure consistent output format
- Check for LLM model compatibility
- Review agent logs for details