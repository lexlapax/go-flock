// ABOUTME: This agent specializes in searching and gathering academic research papers from various sources.
// ABOUTME: It supports multiple output formats (markdown, JSON, text) and uses research-specific tools.

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

// NewResearchPapersAgent creates an agent specialized in academic research paper search
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
	
	logger.Debug(ctx, "Created ResearchPapersAgent with tools: %s, %s, %s", 
		researchTool.Name(), fetchTool.Name(), metadataTool.Name())

	// Set model if specified
	if options.Model != "" {
		agent.WithModel(options.Model)
		logger.Debug(ctx, "Set model to: %s", options.Model)
	}

	// Set system prompt based on output format
	prompt := getResearchPapersPrompt(options.OutputFormat)
	agent.SetSystemPrompt(prompt)
	logger.Debug(ctx, "Set output format to: %s", options.OutputFormat)

	return agent
}

// getResearchPapersPrompt returns the appropriate system prompt based on output format
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

// coreResearchPapersPrompt defines the agent's role and approach - shared across all formats
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

// formatInstructionsMarkdown specifies how to format output as Markdown
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

// formatInstructionsJSON specifies how to format output as JSON
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

// formatInstructionsText specifies how to format output as plain text
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
