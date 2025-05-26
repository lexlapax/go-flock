# Research Workflow Agent Design

## Overview

The research workflow system uses multiple specialized agents working together to gather, synthesize, verify, and present information from various sources. Each agent is a standard go-llms Agent with specific tools and prompts.

## Agent Architecture

All agents use the standard go-llms `domain.Agent` interface. The workflow is responsible for:
- Creating agents with appropriate tools
- Managing agent execution order
- Passing context between agents
- Handling results and errors

### Information Gathering Agents

#### 1. research_papers (ResearchPapersAgent)
**Purpose**: Search and gather academic papers from research databases
**File**: `pkg/agents/research_papers.go`
**Tools**:
- ResearchPaperAPI (primary tool for arXiv, PubMed, CORE)
- FetchWebPage (for extracting full abstracts)
- ExtractMetadata (for paper metadata)

**System Prompt**: Specialized for academic research, focusing on peer-reviewed sources, methodology quality, and citation relevance.

**Output Structure**:
```go
type ResearchFindings struct {
    Papers []ResearchPaper `json:"papers"`
    Themes []string        `json:"themes"`
    KeyAuthors []string    `json:"key_authors"`
    Timeline string        `json:"timeline"`
}
```

#### 2. gather_news (GatherNewsAgent)
**Purpose**: Collect current news and events related to the research topic
**File**: `pkg/agents/gather_news.go`
**Tools**:
- NewsAPISearch (primary news source)
- BraveSearch (supplementary web search)
- FetchRSSFeed (for specific news feeds)

**System Prompt**: Focused on recency, source credibility, and relevance to the research topic.

**Output Structure**:
```go
type NewsFindings struct {
    Articles []NewsArticle `json:"articles"`
    Trends []string       `json:"trends"`
    Timeline []Event      `json:"timeline"`
    Sentiment string      `json:"overall_sentiment"`
}
```

#### 3. extract_web (ExtractWebAgent)
**Purpose**: Extract supplementary information from general web sources
**File**: `pkg/agents/extract_web.go`
**Tools**:
- BraveSearch (web search)
- FetchWebPage (content extraction)
- ExtractLinks (for following references)
- CheckURLStatus (for validation)

**System Prompt**: Emphasis on finding authoritative sources, technical documentation, and expert opinions.

**Output Structure**:
```go
type WebFindings struct {
    Sources []WebSource   `json:"sources"`
    Insights []string     `json:"key_insights"`
    References []string   `json:"external_references"`
}
```

### Processing Agents

#### 4. synthesize_content (SynthesizeContentAgent)
**Purpose**: Combine and organize information from all gathering agents
**File**: `pkg/agents/synthesize_content.go`
**Tools**: None (uses only LLM capabilities)

**System Prompt**: Expert at organizing information, identifying patterns, and creating coherent narratives from multiple sources.

**Input**: Combined outputs from all gathering agents
**Output Structure**:
```go
type SynthesizedContent struct {
    Outline []Section           `json:"outline"`
    MainFindings []string       `json:"main_findings"`
    Contradictions []string     `json:"contradictions"`
    GapsIdentified []string     `json:"gaps_identified"`
}
```

#### 5. verify_facts (VerifyFactsAgent)
**Purpose**: Cross-reference and verify claims from the synthesized content
**File**: `pkg/agents/verify_facts.go`
**Tools**:
- BraveSearch (for fact-checking)
- Original source access (via context)

**System Prompt**: Critical thinking focus, identifying unsubstantiated claims and verifying facts.

**Output Structure**:
```go
type VerificationResult struct {
    VerifiedClaims []Claim   `json:"verified_claims"`
    DisputedClaims []Claim   `json:"disputed_claims"`
    UnverifiedClaims []Claim `json:"unverified_claims"`
}
```

### Output Agents

#### 6. format_citations (FormatCitationsAgent)
**Purpose**: Create proper citations and bibliography
**File**: `pkg/agents/format_citations.go`
**Tools**: None (uses LLM for formatting)

**System Prompt**: Expert in various citation formats (APA, MLA, Chicago), ensuring proper attribution.

**Output Structure**:
```go
type CitationResult struct {
    InTextCitations map[string]string `json:"in_text_citations"`
    Bibliography []string            `json:"bibliography"`
    FootNotes []string              `json:"footnotes"`
}
```

#### 7. polish_output (PolishOutputAgent)
**Purpose**: Refine and polish the final document
**File**: `pkg/agents/polish_output.go`
**Tools**: None (uses LLM capabilities)

**System Prompt**: Professional editor focusing on clarity, tone consistency, and readability.

**Output**: Final polished document with consistent formatting and style

#### 8. create_summary (CreateSummaryAgent)
**Purpose**: Generate executive summaries and abstracts
**File**: `pkg/agents/create_summary.go`
**Tools**: None (uses LLM capabilities)

**System Prompt**: Skilled at distilling complex information into concise summaries.

**Output Structure**:
```go
type SummaryResult struct {
    ExecutiveSummary string   `json:"executive_summary"`
    Abstract string           `json:"abstract"`
    KeyPoints []string        `json:"key_points"`
    TLDR string              `json:"tldr"`
}
```

## Agent Implementation Pattern

Each agent follows this pattern with configurable output format:

```go
// AgentOptions configures agent behavior
type AgentOptions struct {
    OutputFormat OutputFormat // markdown (default), json, or text
    Model        string      // LLM model to use
}

// OutputFormat defines the output format for agent responses
type OutputFormat string

const (
    OutputFormatMarkdown OutputFormat = "markdown" // Default
    OutputFormatJSON     OutputFormat = "json"
    OutputFormatText     OutputFormat = "text"
)

// NewResearchPapersAgent creates an agent for research paper search
func NewResearchPapersAgent(provider ldomain.Provider, opts ...AgentOptions) domain.Agent {
    agent := workflow.NewAgent(provider)
    
    // Add required tools
    agent.AddTool(tools.NewResearchPaperAPITool())
    agent.AddTool(tools.NewFetchWebPageTool())
    agent.AddTool(tools.NewExtractMetadataTool())
    
    // Configure output format (default: markdown)
    format := OutputFormatMarkdown
    if len(opts) > 0 && opts[0].OutputFormat != "" {
        format = opts[0].OutputFormat
    }
    
    // Set specialized system prompt based on format
    agent.SetSystemPrompt(getResearchPaperAPIPrompt(format))
    
    return agent
}

func getResearchPapersPrompt(format OutputFormat) string {
    // Combine core prompt with format-specific instructions
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

// Core prompt defines agent role and capabilities - shared across all formats
const coreResearchPapersPrompt = `You are a research specialist focused on finding and analyzing academic papers. You have access to research databases through your tools.

Tools available to you:
- ResearchPaperAPI: Search for academic papers across arXiv, PubMed, and CORE databases
- FetchWebPage: Retrieve full content from paper URLs or research websites
- ExtractMetadata: Extract structured metadata from research sources

When given a research query:
1. Use the ResearchPaperAPI tool to find relevant academic papers from multiple databases
2. Use FetchWebPage if you need to get more details about specific papers or access full text
3. Use ExtractMetadata to get structured information from research sources

Your analysis should:
- Focus on peer-reviewed and reputable sources
- Identify key themes and trends in the research
- Highlight important researchers and institutions
- Provide a chronological view of how the research has evolved
- Include proper citations for all referenced work`

// Format-specific instructions appended to core prompt
const formatInstructionsMarkdown = `Provide your findings as a well-formatted Markdown report with these sections:

# Research Findings: [Topic]

## Executive Summary
A brief overview of the research landscape and key findings.

## Key Papers
For each relevant paper:
### [Paper Title]
- **Authors**: List of authors
- **Year**: Publication year
- **Summary**: Brief summary of the paper's contribution
- **URL**: Link to the paper

## Research Timeline
Chronological development of research in this area.

## Key Researchers
Notable researchers and their institutions.

## References
Formatted bibliography of all cited papers.

Use proper Markdown formatting with headers, lists, bold text, and links.`

const formatInstructionsJSON = `Provide your findings as a JSON structure following this exact schema:
{
  "papers": [
    {
      "title": "string",
      "authors": ["string"],
      "year": "string",
      "abstract": "string",
      "url": "string",
      "citations": number
    }
  ],
  "themes": ["string"],
  "key_authors": ["string"],
  "timeline": "string",
  "summary": "string"
}

Ensure the response is valid JSON with proper syntax. Include at least 2-5 papers if available.`

const formatInstructionsText = `Provide your findings as plain text with clear sections:

RESEARCH FINDINGS: [TOPIC IN CAPS]

SUMMARY
Write a paragraph summarizing the research landscape.

KEY PAPERS
List each paper with number, title, authors, year, summary, and link.
1. Paper Title
   Authors: Name1, Name2 (Year)
   Summary: Brief description
   Link: URL

RESEARCH THEMES
- Theme 1
- Theme 2
- Theme 3

KEY RESEARCHERS
- Researcher name at Institution

TIMELINE
- Year: Major development
- Year: Another development

Use simple formatting without markdown symbols or JSON syntax.`
```

### Output Format Usage

Agents can be configured for different output formats:

```go
// Default: Markdown output for human reading
agent := NewResearchPapersAgent(provider)

// JSON output for programmatic processing
agent := NewResearchPapersAgent(provider, AgentOptions{
    OutputFormat: OutputFormatJSON,
})

// Plain text for simple integration
agent := NewResearchPapersAgent(provider, AgentOptions{
    OutputFormat: OutputFormatText,
})
```

This ensures agents can be:
- Used standalone with readable Markdown reports (default)
- Integrated into workflows requiring structured JSON data
- Used in text-only environments
- Configured per use case without code changes

## Workflow Coordination

The ResearchWorkflow manages agent execution without requiring a custom agent interface:

### Workflow Structure

```go
type ResearchWorkflow struct {
    name string
    llmProvider ldomain.Provider
    agents map[string]domain.Agent
    results map[string]interface{}
}

func NewResearchWorkflow(provider ldomain.Provider) *ResearchWorkflow {
    w := &ResearchWorkflow{
        name: "research_workflow",
        llmProvider: provider,
        agents: make(map[string]domain.Agent),
        results: make(map[string]interface{}),
    }
    
    // Initialize all agents
    w.agents["research_papers"] = NewResearchPapersAgent(provider)
    w.agents["gather_news"] = NewGatherNewsAgent(provider)
    w.agents["extract_web"] = NewExtractWebAgent(provider)
    w.agents["synthesize_content"] = NewSynthesizeContentAgent(provider)
    w.agents["verify_facts"] = NewVerifyFactsAgent(provider)
    w.agents["format_citations"] = NewFormatCitationsAgent(provider)
    w.agents["polish_output"] = NewPolishOutputAgent(provider)
    w.agents["create_summary"] = NewCreateSummaryAgent(provider)
    
    return w
}
```

### Execution Flow

1. **Parallel Information Gathering**
   ```go
   // Execute gathering agents in parallel
   var wg sync.WaitGroup
   for _, agentName := range []string{"research_papers", "gather_news", "extract_web"} {
       wg.Add(1)
       go func(name string) {
           defer wg.Done()
           result, err := w.agents[name].Run(ctx, query)
           w.results[name] = result
       }(agentName)
   }
   wg.Wait()
   ```

2. **Sequential Processing**
   ```go
   // Synthesize gathered information
   synthesisInput := combineResults(w.results)
   synthesized, _ := w.agents["synthesize_content"].Run(ctx, synthesisInput)
   
   // Verify facts
   verified, _ := w.agents["verify_facts"].Run(ctx, synthesized)
   
   // Format citations
   withCitations, _ := w.agents["format_citations"].Run(ctx, verified)
   
   // Polish output
   polished, _ := w.agents["polish_output"].Run(ctx, withCitations)
   
   // Create summaries
   summary, _ := w.agents["create_summary"].Run(ctx, polished)
   ```

### Context Passing

Since agents are stateless, context is passed through structured prompts:

```go
func combineResults(results map[string]interface{}) string {
    // Format results into a structured prompt for the next agent
    prompt := fmt.Sprintf(`
Based on the following research findings, synthesize a comprehensive report:

ACADEMIC RESEARCH:
%s

CURRENT NEWS:
%s

WEB SOURCES:
%s

Please organize this information into a coherent structure...
    `, 
    formatJSON(results["research_papers"]),
    formatJSON(results["gather_news"]),
    formatJSON(results["extract_web"]))
    
    return prompt
}
```

## Example Usage

```go
func main() {
    // Initialize provider
    provider := provider.NewOpenAIProvider(apiKey, "gpt-4")
    
    // Create workflow
    workflow := NewResearchWorkflow(provider)
    
    // Execute workflow
    result, err := workflow.Execute(context.Background(), WorkflowInput{
        Query: "Impact of AI on healthcare diagnostics",
        Options: WorkflowOptions{
            MaxResults: 20,
            DateRange: "2020-2024",
            OutputFormat: "academic_paper",
        },
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Access results
    fmt.Println("Executive Summary:", result.Summary)
    fmt.Println("Full Report:", result.FullReport)
    fmt.Println("Citations:", result.Bibliography)
}
```

## Benefits of This Design

1. **Simplicity**: Uses standard go-llms Agent interface
2. **Modularity**: Each agent is independent and focused
3. **Flexibility**: Easy to add/remove agents or change execution order
4. **Testability**: Each agent can be tested independently
5. **Reusability**: Agents can be used in other workflows

## Implementation Order

1. Create base agent implementations (gathering agents first)
2. Implement the workflow orchestration
3. Add processing and output agents
4. Create comprehensive example
5. Add tests for each component