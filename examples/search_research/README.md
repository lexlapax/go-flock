# Search Research Tool Example

This example demonstrates how to use the SearchResearch tool to search for academic papers across multiple research databases.

## Overview

The SearchResearch tool provides parallel search capabilities across:
- **arXiv** - Open access repository for scientific papers
- **PubMed** - Biomedical literature database  
- **CORE** - Aggregator of open access research papers

## Features

- **Parallel Search**: Queries multiple databases simultaneously
- **Unified Results**: Combines papers from all sources into a single result set
- **Rich Metadata**: Returns title, authors, abstract, publication date, URLs, and identifiers
- **Provider Status**: Shows which providers returned results and response times
- **Error Handling**: Gracefully handles missing API keys or provider failures

## Running the Example

```bash
cd examples/search_research
go run main.go
```

## API Keys

The tool works with multiple search providers:

### Required for Full Functionality
- **CORE_API_KEY**: Get from https://core.ac.uk/services/api
  ```bash
  export CORE_API_KEY=your_core_api_key
  ```

### Optional (No Keys Required)
- **arXiv**: Free access, no API key needed
- **PubMed**: Free access, optional API key for higher rate limits
  ```bash
  export PUBMED_API_KEY=your_pubmed_api_key  # Optional
  ```

## Example Usage

The example demonstrates four use cases:

1. **Basic Research Search** - Simple query across all available providers
2. **Multi-Source Research** - Shows provider statistics and response times
3. **Domain-Specific Research** - Searches for papers in a specific field
4. **Time-Sensitive Research** - Analyzes temporal distribution of results

## Output Format

Each paper in the results includes:
- Title
- Authors (list)
- Abstract
- Publication date
- Source database
- URL to the paper
- PDF URL (when available)
- Identifiers (DOI, arXiv ID, PubMed ID)
- Journal name (when available)
- Relevance score

## Demo Mode

When no API keys are set, the example runs in demo mode with mock data showing the expected output format.

## Parameters

The SearchResearch tool accepts these parameters:
- `query` (required): Search query string
- `max_results`: Maximum results per provider (default: 10)
- `start_date`: Filter papers from this date (YYYY-MM-DD)
- `end_date`: Filter papers until this date (YYYY-MM-DD)
- `authors`: Filter by author names
- `categories`: Subject categories (cs, physics, medicine, etc.)
- `open_access`: Only return open access papers
- `sort_by`: Sort order (relevance, date, citations)
- `providers`: Specific providers to search