# Examples Organization

## Directory Structure

The examples have been reorganized into a clear hierarchy based on component type:

```
examples/
├── agents/          # Agent implementation examples
│   ├── README.md
│   └── search_research/ # Complete Search Research Agent
├── tools/           # Tool usage examples  
│   ├── README.md
│   ├── brave_search/    # Brave Search API
│   ├── datetime/        # Date/time tools
│   ├── feed/           # RSS feed tools
│   ├── news_api/       # News API search
│   ├── search_research/ # Academic research
│   └── web/            # Web scraping tools
└── workflows/       # Workflow examples
    ├── README.md
    └── basic/       # Basic workflow pattern
```

## Benefits

1. **Clear Organization**: Examples are grouped by their type (agent, tool, workflow)
2. **Easy Navigation**: Developers can quickly find relevant examples
3. **Scalable Structure**: Easy to add new examples in the appropriate category
4. **Documentation**: Each category has its own README explaining concepts

## Building Examples

The Makefile has been updated to handle the new structure:

```bash
# Build all examples
make build-examples

# Binaries are created with category prefixes:
# - agents-search_research
# - tools-datetime
# - tools-feed
# - workflows-basic
# etc.
```

## Running Examples

Examples can be run directly from their directories:

```bash
cd examples/tools/datetime
go run main.go
```

Or using the built binaries:

```bash
./bin/tools-datetime
```

## Documentation Updates

All documentation has been updated to reflect the new paths:
- Main README.md
- Examples README.md  
- Individual category READMEs