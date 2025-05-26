# Agent Examples

This directory contains examples demonstrating how to create and use agents in go-flock.

## Overview

Agents in go-flock are intelligent units that combine LLM capabilities with tool usage to perform complex tasks. They extend the base `domain.Agent` interface from go-llms.

## Examples

### research_papers/
A complete, production-ready agent implementation that serves as the primary example for building agents:
- **Research Papers Agent** for academic paper discovery and analysis
- **Full CLI implementation** with comprehensive command-line options
- **Multiple output formats** (Markdown, JSON, Text) for different use cases
- **Tool integration** with ResearchPaperAPI, FetchWebPage, and ExtractMetadata
- **Provider flexibility** supporting OpenAI, Anthropic, and Gemini
- **Environment configuration** for API keys and settings
- **Error handling** and user feedback patterns
- **Production patterns** including logging, validation, and graceful degradation

This example demonstrates all the key concepts needed to build your own agents.

## Running the Example

```bash
# Navigate to the research papers agent
cd research_papers

# Run with a research query
go run main.go -query "your research topic"

# Run with additional options
go run main.go -query "deep learning" -format json -output results.json
```

## Key Concepts

1. **Agent Creation**: Agents are created with an LLM provider and configured with tools
2. **Tool Integration**: Agents can use multiple tools to accomplish tasks
3. **Output Formats**: Agents support configurable output formats (markdown, json, text)
4. **Coordination**: Agents can work together in multi-agent scenarios

## Next Steps

- Explore specific agent implementations in `pkg/agents/`
- See how agents are used in workflows in `examples/workflows/`
- Review the agent documentation in `docs/agents/`