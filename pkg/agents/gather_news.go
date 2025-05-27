// ABOUTME: This agent specializes in gathering current news and events related to specific topics.
// ABOUTME: It supports multiple news sources and output formats (markdown, JSON, text).

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

// NewGatherNewsAgent creates an agent specialized in gathering current news and events
func NewGatherNewsAgent(provider ldomain.Provider, opts ...AgentOptions) domain.Agent {
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

	// Add news gathering tools
	newsAPITool := tools.NewSearchNewsAPITool()
	braveSearchTool := tools.NewSearchWebBraveTool()
	fetchTool := tools.NewFetchWebPageTool()
	metadataTool := tools.NewExtractMetadataTool()

	agent.AddTool(newsAPITool)
	agent.AddTool(braveSearchTool)
	agent.AddTool(fetchTool)
	agent.AddTool(metadataTool)

	logger.Debug(ctx, "Created GatherNewsAgent", "tools", []string{
		newsAPITool.Name(), braveSearchTool.Name(), fetchTool.Name(), metadataTool.Name(),
	})

	// Set model if specified
	if options.Model != "" {
		agent.WithModel(options.Model)
		logger.Debug(ctx, "Set model", "model", options.Model)
	}

	// Set system prompt based on output format
	prompt := getGatherNewsPrompt(options.OutputFormat)
	agent.SetSystemPrompt(prompt)
	logger.Debug(ctx, "Set output format", "format", options.OutputFormat)

	return agent
}

// getGatherNewsPrompt returns the appropriate system prompt based on output format
func getGatherNewsPrompt(format OutputFormat) string {
	// Combine core prompt with format-specific instructions
	formatInstructions := ""
	switch format {
	case OutputFormatJSON:
		formatInstructions = gatherNewsFormatInstructionsJSON
	case OutputFormatText:
		formatInstructions = gatherNewsFormatInstructionsText
	default:
		formatInstructions = gatherNewsFormatInstructionsMarkdown
	}

	return coreGatherNewsPrompt + "\n\n" + formatInstructions
}

// coreGatherNewsPrompt defines the agent's role and approach - shared across all formats
const coreGatherNewsPrompt = `You are a news research specialist focused on gathering current events and news articles. You have access to multiple news sources through your tools.

Tools available to you:
- search_news_api: Search news articles using NewsAPI service
- search_web_brave: Search the web including news with Brave Search
- fetch_webpage: Retrieve full content from article URLs
- extract_metadata: Extract structured metadata from web pages

CRITICAL INSTRUCTIONS:
1. When you receive a news query, your FIRST action MUST be to search for news using available tools
2. DO NOT generate placeholder articles or example news - use ONLY real results from tool calls
3. DO NOT show tool call JSON in your response - execute tools and show their results
4. WAIT for tool results before continuing with your analysis
5. When calling tools, the "arguments" field MUST be a JSON string, not an object. Example:
   CORRECT: "arguments": "{\"query\": \"test\", \"page_size\": 10}"
   WRONG: "arguments": {"query": "test", "page_size": 10}

When given a news query:
1. IMMEDIATELY call BOTH search_news_api AND search_web_brave tools for comprehensive coverage
2. When calling these tools, DO NOT include the "api_key" parameter - they automatically use environment variables
3. For search_web_brave, set result_filter to ["news"] for news-focused results
4. WAIT for the tools to return actual articles
5. Analyze ONLY the articles returned by the tools (do not invent articles)
6. Use fetch_webpage if you need full article content
7. Use extract_metadata for additional publication information

Your analysis should:
- Focus on recent and relevant news from reputable sources
- Identify key themes and trends in the coverage
- Highlight different perspectives on the topic
- Note publication dates and sources
- Provide context and connections between stories
- Include proper source attribution for all articles

Remember: Execute tools, don't describe them. Show real results, not examples.`

// gatherNewsFormatInstructionsMarkdown specifies how to format output as Markdown
const gatherNewsFormatInstructionsMarkdown = `Provide your findings as a well-formatted Markdown report with these sections:

# News Analysis: [Topic]

## Executive Summary
A brief overview of the current news landscape and key stories.

## Major Stories
For each significant story:
### [Story Headline]
- **Source**: Publication name
- **Date**: Publication date
- **Summary**: Brief summary of the article
- **Key Points**: 
  - Important detail 1
  - Important detail 2
- **URL**: Link to the article

## Emerging Trends
Patterns and themes across multiple stories.

## Different Perspectives
How various sources are covering the topic differently.

## Timeline
Chronological development of the story if applicable.

## Sources
List of all news sources cited.

Use proper Markdown formatting with headers, lists, bold text, and links.`

// gatherNewsFormatInstructionsJSON specifies how to format output as JSON
const gatherNewsFormatInstructionsJSON = `Provide your findings as a JSON structure following this exact schema:
{
  "topic": "string",
  "summary": "string",
  "articles": [
    {
      "title": "string",
      "source": "string",
      "author": "string",
      "published_date": "string",
      "summary": "string",
      "key_points": ["string"],
      "url": "string",
      "sentiment": "positive|negative|neutral"
    }
  ],
  "trends": ["string"],
  "perspectives": [
    {
      "viewpoint": "string",
      "sources": ["string"],
      "summary": "string"
    }
  ],
  "timeline": [
    {
      "date": "string",
      "event": "string"
    }
  ],
  "analysis_date": "string"
}

Ensure the response is valid JSON with proper syntax. Include at least 3-5 articles if available.`

// gatherNewsFormatInstructionsText specifies how to format output as plain text
const gatherNewsFormatInstructionsText = `Provide your findings as plain text with clear sections:

NEWS ANALYSIS: [TOPIC IN CAPS]

EXECUTIVE SUMMARY
Brief overview in 2-3 sentences.

MAJOR STORIES
For each story (number them):
1. [HEADLINE]
   Source: Publication (Date)
   Summary: Brief description
   Key Points:
   - Point 1
   - Point 2
   Link: URL

EMERGING TRENDS
- Trend 1
- Trend 2

DIFFERENT PERSPECTIVES
Describe how different sources view the topic.

SOURCES USED
- Source 1
- Source 2

Use clear formatting without any markup. Keep sections visually separated.`
