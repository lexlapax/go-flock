# Research Workflow

## Overview

The Research Workflow is a comprehensive information gathering and synthesis system that orchestrates multiple specialized agents to produce in-depth research reports. It combines academic research, current news, and web sources to create well-referenced, fact-checked documents.

## Workflow Architecture

### Core Components

1. **Information Gathering Phase** (Parallel Execution)
   - Academic research from scholarly databases
   - Current news and events
   - General web content and expert opinions

2. **Processing Phase** (Sequential Execution)
   - Content synthesis and organization
   - Fact verification and cross-referencing
   - Citation formatting and attribution

3. **Output Phase** (Sequential Execution)
   - Content polishing and editing
   - Summary generation for different audiences

### Workflow Types

#### 1. ComprehensiveResearchWorkflow
**Purpose**: Full academic-style research with extensive verification
**Duration**: 5-10 minutes
**Best For**: In-depth analysis, white papers, research reports

**Stages**:
1. Parallel gathering (3 agents)
2. Synthesis
3. Fact-checking
4. Citation formatting
5. Polish and edit
6. Generate summaries

#### 2. QuickResearchWorkflow
**Purpose**: Rapid research overview with key findings
**Duration**: 2-3 minutes
**Best For**: Executive briefings, quick market analysis

**Stages**:
1. Parallel gathering (limited to 10 sources each)
2. Synthesis
3. Basic fact-checking
4. Summary generation

#### 3. NewsAnalysisWorkflow
**Purpose**: Current events analysis with trend identification
**Duration**: 3-5 minutes
**Best For**: News digests, trend reports, event analysis

**Stages**:
1. News gathering (extended)
2. Web verification
3. Trend synthesis
4. Timeline creation
5. Executive summary

#### 4. AcademicReviewWorkflow
**Purpose**: Literature review focusing on scholarly sources
**Duration**: 5-8 minutes
**Best For**: Academic papers, systematic reviews

**Stages**:
1. Research paper gathering (extended)
2. Methodology analysis
3. Citation network mapping
4. Academic synthesis
5. Formal citation formatting

## Workflow Configuration

### Input Parameters

```go
type ResearchWorkflowInput struct {
    Query        string              `json:"query"`
    WorkflowType string              `json:"workflow_type"`
    Options      ResearchOptions     `json:"options"`
}

type ResearchOptions struct {
    // Scope parameters
    MaxSources      int      `json:"max_sources"`
    DateRange       string   `json:"date_range"`
    Domains         []string `json:"domains"`
    ExcludeDomains  []string `json:"exclude_domains"`
    
    // Output parameters
    OutputFormat    string   `json:"output_format"`    // "academic", "business", "technical"
    CitationStyle   string   `json:"citation_style"`   // "APA", "MLA", "Chicago"
    SummaryFormats  []string `json:"summary_formats"`  // ["executive", "abstract", "bullets"]
    
    // Quality parameters
    VerificationLevel string `json:"verification_level"` // "basic", "standard", "comprehensive"
    MinSourceQuality  float64 `json:"min_source_quality"`
}
```

### Output Structure

```go
type ResearchWorkflowResult struct {
    // Metadata
    WorkflowID   string    `json:"workflow_id"`
    Query        string    `json:"query"`
    Timestamp    time.Time `json:"timestamp"`
    Duration     string    `json:"duration"`
    
    // Main outputs
    FullReport   ResearchReport `json:"full_report"`
    Summaries    Summaries      `json:"summaries"`
    
    // Supporting data
    Sources      []Source       `json:"sources"`
    Citations    Bibliography   `json:"citations"`
    Verification FactCheckReport `json:"verification"`
    
    // Agent outputs (for debugging/transparency)
    AgentOutputs map[string]string `json:"agent_outputs"`
}

type ResearchReport struct {
    Title       string    `json:"title"`
    Content     string    `json:"content"`      // Markdown, JSON, or text based on output format
    Format      string    `json:"format"`       // "markdown", "json", or "text"
    Sections    []Section `json:"sections"`
    Confidence  float64   `json:"confidence_score"`
}

type Summaries struct {
    Executive   string   `json:"executive"`
    Abstract    string   `json:"abstract"`
    KeyPoints   []string `json:"key_points"`
    TLDR        string   `json:"tldr"`
}
```

## Workflow Implementation

### Base Workflow Structure

```go
type ResearchWorkflow struct {
    name        string
    description string
    llmProvider ldomain.Provider
    agents      map[string]domain.Agent
    config      WorkflowConfig
    results     sync.Map // Thread-safe results storage
}

type WorkflowConfig struct {
    MaxConcurrency int
    Timeout        time.Duration
    RetryPolicy    RetryPolicy
    OutputFormat   OutputFormat // Default output format for all agents
}
```

### Execution Flow

```go
func (w *ResearchWorkflow) Execute(ctx context.Context, input ResearchWorkflowInput) (ResearchWorkflowResult, error) {
    start := time.Now()
    
    // 1. Validate input
    if err := w.validateInput(input); err != nil {
        return ResearchWorkflowResult{}, err
    }
    
    // 2. Configure agents based on workflow type and output format
    outputFormat := w.config.OutputFormat
    if input.Options.OutputFormat != "" {
        outputFormat = OutputFormat(input.Options.OutputFormat)
    }
    w.configureAgents(input.WorkflowType, input.Options, outputFormat)
    
    // 3. Execute gathering phase
    gatherResults, err := w.executeGatheringPhase(ctx, input.Query)
    if err != nil {
        return ResearchWorkflowResult{}, err
    }
    
    // 4. Execute processing phase
    processedResults, err := w.executeProcessingPhase(ctx, gatherResults)
    if err != nil {
        return ResearchWorkflowResult{}, err
    }
    
    // 5. Execute output phase
    finalResults, err := w.executeOutputPhase(ctx, processedResults)
    if err != nil {
        return ResearchWorkflowResult{}, err
    }
    
    // 6. Compile final result
    return w.compileResults(input, finalResults, time.Since(start)), nil
}
```

### Error Handling and Recovery

```go
type RetryPolicy struct {
    MaxRetries     int
    BackoffFactor  float64
    MaxBackoff     time.Duration
}

func (w *ResearchWorkflow) executeWithRetry(ctx context.Context, agentName string, input string) (interface{}, error) {
    var lastErr error
    
    for attempt := 0; attempt <= w.config.RetryPolicy.MaxRetries; attempt++ {
        if attempt > 0 {
            backoff := time.Duration(float64(time.Second) * math.Pow(w.config.RetryPolicy.BackoffFactor, float64(attempt-1)))
            if backoff > w.config.RetryPolicy.MaxBackoff {
                backoff = w.config.RetryPolicy.MaxBackoff
            }
            time.Sleep(backoff)
        }
        
        result, err := w.agents[agentName].Run(ctx, input)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        
        // Check if error is retryable
        if !isRetryableError(err) {
            break
        }
    }
    
    return nil, fmt.Errorf("agent %s failed after %d attempts: %w", agentName, w.config.RetryPolicy.MaxRetries+1, lastErr)
}
```

## Usage Examples

### Example 1: Comprehensive Research

```go
func main() {
    // Initialize provider
    provider := provider.NewOpenAIProvider(apiKey, "gpt-4")
    
    // Create workflow
    workflow := workflows.NewComprehensiveResearchWorkflow(provider)
    
    // Execute research
    result, err := workflow.Execute(context.Background(), ResearchWorkflowInput{
        Query: "Impact of large language models on software development practices",
        WorkflowType: "comprehensive",
        Options: ResearchOptions{
            MaxSources: 30,
            DateRange: "2022-2024",
            OutputFormat: "technical",
            CitationStyle: "APA",
            SummaryFormats: []string{"executive", "bullets"},
            VerificationLevel: "comprehensive",
        },
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Save outputs
    os.WriteFile("research_report.md", []byte(result.FullReport.Markdown), 0644)
    os.WriteFile("executive_summary.md", []byte(result.Summaries.Executive), 0644)
    
    // Print statistics
    fmt.Printf("Research completed in %s\n", result.Duration)
    fmt.Printf("Sources analyzed: %d\n", len(result.Sources))
    fmt.Printf("Confidence score: %.2f\n", result.FullReport.Confidence)
}
```

### Example 2: Quick News Analysis

```go
func quickNewsAnalysis() {
    workflow := workflows.NewQuickResearchWorkflow(provider)
    
    result, err := workflow.Execute(context.Background(), ResearchWorkflowInput{
        Query: "Latest developments in renewable energy technology",
        WorkflowType: "quick",
        Options: ResearchOptions{
            MaxSources: 10,
            DateRange: "7d", // Last 7 days
            OutputFormat: "business",
            SummaryFormats: []string{"tldr", "bullets"},
        },
    })
    
    // Get quick insights
    fmt.Println("TL;DR:", result.Summaries.TLDR)
    fmt.Println("\nKey Points:")
    for _, point := range result.Summaries.KeyPoints {
        fmt.Println("-", point)
    }
}
```

### Example 3: Academic Literature Review

```go
func academicReview() {
    workflow := workflows.NewAcademicReviewWorkflow(provider)
    
    result, err := workflow.Execute(context.Background(), ResearchWorkflowInput{
        Query: "Machine learning applications in drug discovery",
        WorkflowType: "academic",
        Options: ResearchOptions{
            MaxSources: 50,
            DateRange: "2020-2024",
            Domains: []string{"arxiv.org", "pubmed.ncbi.nlm.nih.gov", "nature.com"},
            OutputFormat: "academic",
            CitationStyle: "APA",
            VerificationLevel: "comprehensive",
            MinSourceQuality: 0.8,
        },
    })
    
    // Generate LaTeX bibliography
    bibliography := result.Citations.ToLaTeX()
    os.WriteFile("references.bib", []byte(bibliography), 0644)
}
```

## Advanced Features

### 1. Progressive Enhancement
The workflow can start with quick results and progressively enhance them:

```go
func (w *ResearchWorkflow) ExecuteProgressive(ctx context.Context, input ResearchWorkflowInput, updates chan<- ProgressUpdate) {
    // Send initial quick results
    quickResults := w.executeQuickGathering(ctx, input.Query)
    updates <- ProgressUpdate{Stage: "initial", Results: quickResults}
    
    // Enhance with deeper research
    deepResults := w.executeDeepGathering(ctx, input.Query)
    updates <- ProgressUpdate{Stage: "enhanced", Results: deepResults}
    
    // Final verification and polish
    finalResults := w.executeFinalProcessing(ctx, deepResults)
    updates <- ProgressUpdate{Stage: "final", Results: finalResults}
}
```

### 2. Custom Agent Configuration
Workflows can be customized with different agent configurations:

```go
func NewCustomResearchWorkflow(provider ldomain.Provider, agentConfigs map[string]AgentConfig) *ResearchWorkflow {
    w := &ResearchWorkflow{
        name: "custom_research",
        llmProvider: provider,
        agents: make(map[string]domain.Agent),
    }
    
    // Configure each agent with custom settings
    for agentName, config := range agentConfigs {
        agentOpts := AgentOptions{
            OutputFormat: config.OutputFormat,
            Model: config.Model,
        }
        agent := createAgent(agentName, provider, agentOpts)
        if config.SystemPrompt != "" {
            agent.SetSystemPrompt(config.SystemPrompt)
        }
        w.agents[agentName] = agent
    }
    
    return w
}
```

### 3. Workflow Composition
Workflows can be composed from smaller workflows:

```go
type CompositeWorkflow struct {
    workflows []Workflow
}

func (c *CompositeWorkflow) Execute(ctx context.Context, input interface{}) (interface{}, error) {
    var result interface{}
    
    for _, workflow := range c.workflows {
        var err error
        result, err = workflow.Execute(ctx, result)
        if err != nil {
            return nil, err
        }
    }
    
    return result, nil
}
```

## Performance Considerations

1. **Parallel Execution**: Gathering agents run concurrently to minimize latency
2. **Result Caching**: Intermediate results are cached to enable retry without re-execution
3. **Streaming**: Large results can be streamed to avoid memory issues
4. **Timeout Management**: Each phase has configurable timeouts
5. **Resource Limits**: Maximum source counts prevent runaway costs

## Integration Points

The Research Workflow can be integrated with:

1. **Web APIs**: Expose as REST endpoint
2. **CLI Tools**: Command-line interface for batch processing
3. **Scheduled Jobs**: Periodic research updates
4. **Chat Interfaces**: Interactive research sessions
5. **Document Systems**: Direct integration with knowledge bases