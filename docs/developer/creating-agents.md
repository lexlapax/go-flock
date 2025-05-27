# Creating Custom Agents

This guide walks through creating a custom agent for go-flock, using the Research Papers Agent as a comprehensive example. We'll cover architecture, implementation, debugging, and best practices.

## Overview

Agents in go-flock are intelligent units that combine LLM capabilities with tool usage to accomplish specific tasks. They extend the go-llms `domain.Agent` interface and can coordinate with other agents in complex workflows.

## Case Study: Research Papers Agent

The Research Papers Agent demonstrates a complete production-ready implementation that searches academic databases and generates comprehensive research reports.

### 1. Planning the Agent

Before coding, define:

#### Purpose
- Search multiple academic databases (arXiv, PubMed, CORE)
- Analyze and synthesize research findings
- Generate reports in multiple formats (Markdown, JSON, Text)

#### Required Tools
- `research_paper_api` - Search academic databases
- `fetch_webpage` - Retrieve paper content
- `extract_metadata` - Extract structured information

#### User Interface
- Command-line interface with flags
- Support for different LLM providers
- Configurable output formats

### 2. Project Structure

```
pkg/agents/
└── research_papers.go      # Agent implementation
examples/agents/
└── research_papers/
    └── main.go            # CLI interface
```

### 3. Agent Implementation

#### Basic Structure

```go
// pkg/agents/research_papers.go
package agents

import (
    "context"
    "log/slog"
    "os"
    
    "github.com/lexlapax/go-flock/pkg/common"
    "github.com/lexlapax/go-flock/pkg/tools"
    "github.com/lexlapax/go-llms/pkg/agent/domain"
    "github.com/lexlapax/go-llms/pkg/agent/workflow"
    ldomain "github.com/lexlapax/go-llms/pkg/llm/domain"
)

// AgentOptions configures the agent behavior
type AgentOptions struct {
    Model        string
    OutputFormat OutputFormat
}

// OutputFormat defines how results are formatted
type OutputFormat string

const (
    OutputFormatMarkdown OutputFormat = "markdown"
    OutputFormatJSON     OutputFormat = "json"
    OutputFormatText     OutputFormat = "text"
)
```

#### Creating the Agent Function

```go
func NewResearchPapersAgent(provider ldomain.Provider, opts ...AgentOptions) domain.Agent {
    logger := common.GetLogger()
    ctx := context.Background()
    
    // Get options or use defaults
    options := DefaultAgentOptions()
    if len(opts) > 0 {
        options = opts[0]
    }

    // Create base agent
    agent := workflow.NewAgent(provider)
    
    // Add logging hook if debug mode is enabled
    if os.Getenv("FLOCK_DEBUG") == "true" || os.Getenv("FLOCK_DEBUG") == "1" {
        slogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
            Level: slog.LevelDebug,
        }))
        loggingHook := workflow.NewLoggingHook(slogger, workflow.LogLevelDebug)
        agent.WithHook(loggingHook)
        logger.Debug(ctx, "Added debug logging hook to agent")
    }
    
    // Add research-specific tools
    researchTool := tools.NewResearchPaperAPITool()
    fetchTool := tools.NewFetchWebPageTool()
    metadataTool := tools.NewExtractMetadataTool()
    
    agent.AddTool(researchTool)
    agent.AddTool(fetchTool)
    agent.AddTool(metadataTool)
    
    // Set model if specified
    if options.Model != "" {
        agent.WithModel(options.Model)
    }
    
    // Set system prompt based on output format
    prompt := getResearchPapersPrompt(options.OutputFormat)
    agent.SetSystemPrompt(prompt)
    
    return agent
}
```

### 4. System Prompt Design

The system prompt is crucial for agent behavior. We used a modular approach:

#### Core Prompt (Shared Across Formats)

```go
const coreResearchPapersPrompt = `You are a research specialist focused on finding and analyzing academic papers. You have access to research databases through your tools.

Tools available to you:
- ResearchPaperAPI: Search for academic papers across arXiv, PubMed, and CORE databases
- FetchWebPage: Retrieve full content from paper URLs or research websites
- ExtractMetadata: Extract structured metadata from research sources

CRITICAL INSTRUCTIONS:
1. When you receive a research query, your FIRST action MUST be to call the ResearchPaperAPI tool
2. DO NOT generate placeholder data or example papers - use ONLY real results from tool calls
3. DO NOT show tool call JSON in your response - execute tools and show their results
4. WAIT for tool results before continuing with your analysis
5. When calling tools, the "arguments" field MUST be a JSON string, not an object. Example:
   CORRECT: "arguments": "{\"query\": \"test\", \"max_results\": 5}"
   WRONG: "arguments": {"query": "test", "max_results": 5}

When given a research query:
1. IMMEDIATELY call ResearchPaperAPI with appropriate parameters
2. WAIT for the tool to return actual papers
3. Analyze ONLY the papers returned by the tool (do not invent papers)
4. Use FetchWebPage if you need more details about specific papers
5. Use ExtractMetadata for additional structured information

Your analysis should:
- Focus on peer-reviewed and reputable sources from the actual search results
- Identify key themes and trends in the papers found
- Highlight important researchers and institutions mentioned
- Provide a chronological view based on publication dates
- Include proper citations for all referenced papers

Remember: Execute tools, don't describe them. Show real results, not examples.`
```

#### Format-Specific Instructions

```go
// Modular prompt architecture
func getResearchPapersPrompt(format OutputFormat) string {
    formatInstructions := ""
    switch format {
    case OutputFormatJSON:
        formatInstructions = formatInstructionsJSON
    case OutputFormatText:
        formatInstructions = formatInstructionsText
    default:
        formatInstructions = formatInstructionsMarkdown
    }
    
    return coreResearchPapersPrompt + "\n\n" + formatInstructions
}
```

### 5. CLI Implementation

Create a user-friendly command-line interface:

```go
// examples/agents/research_papers/main.go
package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    
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
        model        = flag.String("model", "", "LLM model to use")
        providerName = flag.String("provider", "", "LLM provider: openai, anthropic, or gemini")
        output       = flag.String("output", "", "Output file (optional)")
        debug        = flag.Bool("debug", false, "Enable debug logging")
        help         = flag.Bool("help", false, "Show help message")
    )
    
    flag.Parse()
    
    // Initialize logging based on debug flag
    common.InitLogger(*debug)
    logger := common.GetLogger()
    
    // Validate inputs...
    
    // Create LLM provider
    ctx := context.Background()
    logger.Debug(ctx, "Creating LLM provider", "provider", *providerName)
    llmProvider, err := createProvider(*providerName)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error creating LLM provider: %v\n", err)
        os.Exit(1)
    }
    
    // Create the research papers agent
    logger.Debug(ctx, "Creating research papers agent", "format", outputFormat, "model", agentOpts.Model)
    agent := agents.NewResearchPapersAgent(llmProvider, agentOpts)
    
    // Execute the search
    logger.Debug(ctx, "Running agent", "query", *query)
    result, err := agent.Run(ctx, *query)
    if err != nil {
        logger.Error(ctx, "Agent execution failed", "error", err)
        fmt.Fprintf(os.Stderr, "Error running agent: %v\n", err)
        os.Exit(1)
    }
    
    // Output results...
}
```

### 6. Debugging and Troubleshooting

#### The Gemini Tool Calling Challenge

During development, we encountered an issue where Google Gemini would return raw JSON tool calls instead of executing them. This taught us important lessons:

**Problem**: Gemini returned:
```json
{
  "tool_calls": [
    {
      "function": {
        "name": "research_paper_api",
        "arguments": {
          "query": "quantum computing"
        }
      }
    }
  ]
}
```

**Root Cause**: 
- go-llms expects the `arguments` field to be a JSON string, not an object
- Gemini doesn't have native tool calling support in go-llms

**Solution**: Update the system prompt to be explicit:
```go
"5. When calling tools, the \"arguments\" field MUST be a JSON string, not an object. Example:\n" +
"   CORRECT: \"arguments\": \"{\\\"query\\\": \\\"test\\\", \\\"max_results\\\": 5}\"\n" +
"   WRONG: \"arguments\": {\"query\": \"test\", \"max_results\": 5}"
```

#### Implementing Debug Logging

We implemented comprehensive logging using slog (compatible with go-llms):

```go
// pkg/common/logger.go
type slogLogger struct {
    logger *slog.Logger
}

func InitLogger(debugMode bool) {
    loggerOnce.Do(func() {
        level := slog.LevelInfo
        if debugMode {
            level = slog.LevelDebug
        }
        
        handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
            Level: level,
        })
        globalLogger = &slogLogger{
            logger: slog.New(handler),
        }
    })
}
```

### 7. Best Practices Learned

#### 1. Modular Prompt Architecture
- Separate core functionality from format-specific instructions
- Makes prompts maintainable and reusable
- Easier to add new output formats

#### 2. Explicit Tool Calling Instructions
- Be very specific about expected formats
- Different LLM providers may interpret instructions differently
- Test with all target providers

#### 3. Comprehensive Logging
- Use structured logging (slog) for compatibility
- Add debug hooks to trace agent execution
- Log at appropriate levels (Debug, Info, Error)

#### 4. Error Handling
- Gracefully handle provider failures
- Provide clear error messages to users
- Use exit codes appropriately in CLI

#### 5. Testing Strategy
- Test individual tools first
- Test agent with mock provider
- Test with real providers
- Test edge cases (no results, API failures)

### 8. Complete Example Structure

Here's the final structure of a production-ready agent:

```
pkg/agents/research_papers.go
├── Package definition
├── Type definitions (AgentOptions, OutputFormat)
├── NewResearchPapersAgent function
│   ├── Logger initialization
│   ├── Options handling
│   ├── Agent creation
│   ├── Debug hook setup
│   ├── Tool registration
│   ├── Model configuration
│   └── System prompt setup
├── Prompt generation functions
│   ├── getResearchPapersPrompt
│   ├── Core prompt constant
│   └── Format-specific constants
└── Helper functions

examples/agents/research_papers/main.go
├── Package and imports
├── Flag definitions
├── Main function
│   ├── Flag parsing
│   ├── Logger initialization
│   ├── Input validation
│   ├── Provider creation
│   ├── Agent creation
│   ├── Query execution
│   └── Result output
└── Helper functions (createProvider)
```

## Creating Your Own Agent

Follow these steps to create a new agent:

1. **Define the Purpose**: What specific task will your agent accomplish?

2. **Identify Required Tools**: Which tools does your agent need?

3. **Design the Interface**: CLI? API? Library function?

4. **Implement the Agent**:
   ```go
   func NewMyAgent(provider ldomain.Provider, opts ...MyAgentOptions) domain.Agent {
       agent := workflow.NewAgent(provider)
       
       // Add tools
       agent.AddTool(tool1)
       agent.AddTool(tool2)
       
       // Set system prompt
       agent.SetSystemPrompt(myPrompt)
       
       return agent
   }
   ```

5. **Write Clear System Prompts**:
   - Define the agent's role
   - List available tools
   - Provide step-by-step instructions
   - Include examples if needed
   - Handle edge cases

6. **Add Debug Support**:
   - Use common.Logger for logging
   - Add workflow hooks for tracing
   - Support FLOCK_DEBUG environment variable

7. **Test Thoroughly**:
   - Unit tests for agent creation
   - Integration tests with mock provider
   - End-to-end tests with real providers
   - Edge case testing

8. **Document Your Agent**:
   - Purpose and capabilities
   - Required environment variables
   - Usage examples
   - Troubleshooting guide

## Advanced Topics

### Multi-Agent Coordination

Agents can work together in workflows:

```go
type ResearchCoordinator struct {
    paperAgent   domain.Agent
    summaryAgent domain.Agent
    reportAgent  domain.Agent
}

func (rc *ResearchCoordinator) Coordinate(ctx context.Context, query string) (string, error) {
    // Step 1: Find papers
    papers, err := rc.paperAgent.Run(ctx, query)
    
    // Step 2: Summarize each paper
    summaries := []string{}
    for _, paper := range papers {
        summary, err := rc.summaryAgent.Run(ctx, paper)
        summaries = append(summaries, summary)
    }
    
    // Step 3: Generate final report
    return rc.reportAgent.Run(ctx, summaries)
}
```

### Custom Hooks

Implement custom hooks for specialized monitoring:

```go
type MetricsHook struct {
    toolCalls map[string]int
    mu        sync.Mutex
}

func (h *MetricsHook) BeforeToolCall(ctx context.Context, toolName string, params map[string]interface{}) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.toolCalls[toolName]++
}
```

### Error Recovery

Implement retry logic and fallback strategies:

```go
func (agent *MyAgent) RunWithRetry(ctx context.Context, prompt string) (string, error) {
    var lastErr error
    for i := 0; i < 3; i++ {
        result, err := agent.Run(ctx, prompt)
        if err == nil {
            return result, nil
        }
        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return "", fmt.Errorf("failed after 3 attempts: %w", lastErr)
}
```

## Conclusion

Creating agents in go-flock involves:
1. Clear purpose definition
2. Thoughtful tool selection
3. Well-crafted system prompts
4. Robust error handling
5. Comprehensive debugging support
6. Thorough testing

The Research Papers Agent demonstrates these principles in a production-ready implementation that handles real-world challenges like provider compatibility and debugging.