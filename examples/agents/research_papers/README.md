# Research Papers Agent Example

This example demonstrates a complete agent implementation - the Research Papers Agent that specializes in finding and analyzing academic papers from multiple research databases.

## Overview

The Research Papers Agent showcases:
- Agent creation with configurable output formats (Markdown, JSON, Text)
- Integration with multiple tools (ResearchPaperAPI, FetchWebPage, ExtractMetadata)
- Command-line interface for agent interaction
- Provider flexibility (OpenAI, Anthropic, Gemini)

## Running the Example

### Prerequisites

Set up your environment variables:
```bash
# LLM Provider (at least one required)
export OPENAI_API_KEY=your_openai_key
# OR
export ANTHROPIC_API_KEY=your_anthropic_key
# OR
export GEMINI_API_KEY=your_gemini_key

# Tool API key (required for ResearchPaperAPI tool)
export BRAVE_API_KEY=your_brave_key
```

### Basic Usage

```bash
# Run with default settings (Markdown output)
go run main.go -query "deep learning medical imaging"

# JSON output
go run main.go -query "climate change impacts" -format json

# Save to file
go run main.go -query "quantum computing" -output results.md

# Use specific provider and model
go run main.go -query "renewable energy" -provider openai -model gpt-4
```

### Command Line Options

- `-query` (required) - Research topic to search for
- `-format` - Output format: markdown (default), json, or text
- `-output` - Save results to file instead of stdout
- `-provider` - LLM provider: openai, anthropic, or gemini
- `-model` - Specific model to use (optional)
- `-help` - Show usage information

## Output Formats

### Markdown (Default)
```markdown
# Research Findings: [Topic]

## Executive Summary
Overview of the research landscape...

## Key Papers
### Paper Title
- **Authors**: Author list
- **Year**: 2024
- **Summary**: Brief summary
- **URL**: Link to paper
```

### JSON
```json
{
  "papers": [...],
  "themes": [...],
  "key_authors": [...],
  "timeline": "...",
  "summary": "..."
}
```

### Plain Text
```
RESEARCH FINDINGS: [TOPIC]

SUMMARY
Overview text...

KEY PAPERS
1. Paper Title
   Authors: Names (Year)
   Summary: Description
   Link: URL
```

## Code Structure

The example demonstrates:

1. **Provider Creation** - Flexible provider selection based on environment
2. **Agent Configuration** - Setting output format and model options
3. **Agent Execution** - Running the agent with user query
4. **Result Handling** - Processing and displaying/saving results

## Key Learnings

- Agents combine LLM capabilities with tool usage
- Output formats make agents flexible for different use cases
- Command-line interfaces make agents accessible
- Environment-based configuration simplifies deployment

## Next Steps

- Modify the system prompt for different research styles
- Add additional tools for enhanced capabilities
- Integrate with workflow systems for complex research tasks
- Create similar agents for other domains