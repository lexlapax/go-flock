# Troubleshooting Guide

This guide helps resolve common issues when using go-flock tools, agents, and workflows.

## Tool Calling Issues

### Problem: Agent Shows Raw JSON Tool Calls Instead of Executing Them

**Symptoms:**
- Agent displays JSON like `{"tool_calls": [...]}` in the output
- Tools are not being executed
- No actual results from tool execution

**Common with:** Google Gemini provider

**Root Cause:** 
Some LLM providers (like Gemini) don't have native tool/function calling support in go-llms. The agent instructs them to return tool calls in JSON format, but the format must be exact for the extraction to work.

**Solutions:**

1. **Ensure Proper JSON Format**
   - The `arguments` field must be a JSON string, not an object
   - Correct: `"arguments": "{\"query\": \"test\", \"max_results\": 5}"`
   - Wrong: `"arguments": {"query": "test", "max_results": 5}`

2. **Enable Debug Logging**
   ```bash
   # Using environment variable
   FLOCK_DEBUG=true ./bin/research_papers -query "test" -provider gemini
   
   # Using CLI flag
   ./bin/research_papers -query "test" -provider gemini -debug
   ```

3. **Check Debug Output**
   Look for these key log entries:
   - `"Response content"` - Shows what the LLM returned
   - `"Calling tool"` - Indicates tool execution started
   - `"Tool executed successfully"` - Confirms tool ran properly

### Problem: Tool Not Found Error

**Symptoms:**
- Error message: `Tool 'tool_name' not found`
- Agent lists available tools in error message

**Solutions:**

1. **Check Tool Registration**
   ```go
   // Ensure tool is added to agent
   agent.AddTool(tools.NewResearchPaperAPITool())
   ```

2. **Verify Tool Name**
   - Tool names are case-sensitive
   - Use exact name as defined in tool creation
   - Common names: `research_paper_api`, `fetch_webpage`, `extract_metadata`

### Problem: API Key Missing or Invalid

**Symptoms:**
- Error: `API key is required`
- Error: `401 Unauthorized`
- Tool execution fails with authentication error

**Solutions:**

1. **Set Environment Variables**
   ```bash
   export NEWS_API_KEY=your_newsapi_key
   export BRAVE_SEARCH_API_KEY=your_brave_api_key
   export CORE_API_KEY=your_core_api_key
   ```

2. **Pass API Key in Parameters**
   ```go
   params := tools.SearchNewsAPIParams{
       Query:  "test",
       APIKey: "your-api-key", // Direct parameter
   }
   ```

## Agent Issues

### Problem: Agent Reaches Maximum Iterations

**Symptoms:**
- Message: `Agent reached maximum iterations without final result`
- Agent seems stuck in a loop

**Solutions:**

1. **Improve System Prompt**
   - Make instructions clearer and more specific
   - Add explicit termination conditions
   - Reduce ambiguity in task descriptions

2. **Check Tool Results**
   - Ensure tools return proper results
   - Verify tool error handling
   - Check for empty or malformed responses

### Problem: Agent Generates Placeholder Data

**Symptoms:**
- Agent returns example/mock data instead of real results
- Papers or data seem fabricated
- URLs are example.com or similar

**Solutions:**

1. **Explicit Instructions**
   The research papers agent already includes instructions to prevent this:
   ```
   DO NOT generate placeholder data or example papers - use ONLY real results from tool calls
   ```

2. **Verify Tool Execution**
   Enable debug logging to ensure tools are actually being called

## Logging and Debugging

### Using the Logger

go-flock uses `slog` for structured logging compatible with go-llms:

```go
import "github.com/lexlapax/go-flock/pkg/common"

// Initialize logger
common.InitLogger(true) // true for debug mode
logger := common.GetLogger()

// Log with structured data
logger.Debug(ctx, "Processing request", "query", query, "provider", provider)
logger.Error(ctx, "Tool execution failed", "error", err, "tool", toolName)
```

### Debug Output Interpretation

Key log entries to watch for:

1. **Agent Initialization**
   ```
   DEBUG msg="Creating research papers agent" format=markdown model=""
   DEBUG msg="Added debug logging hook to agent"
   ```

2. **Tool Registration**
   ```
   DEBUG msg="Created ResearchPapersAgent" tools="[research_paper_api fetch_webpage extract_metadata]"
   ```

3. **LLM Communication**
   ```
   INFO msg="Generating response" emoji=🤔
   DEBUG msg="Message details" index=0 role=system content="..."
   ```

4. **Tool Execution**
   ```
   INFO msg="Calling tool" tool=research_paper_api emoji=🔧
   DEBUG msg="Tool parameters" params="{...}"
   INFO msg="Tool executed successfully" tool=research_paper_api emoji=✅
   ```

## Performance Issues

### Problem: Slow Tool Execution

**Solutions:**

1. **Check API Rate Limits**
   - Most APIs have rate limits
   - Tool implements automatic rate limiting
   - Consider caching results

2. **Parallel Execution**
   - Research Paper API searches providers in parallel
   - Agent can execute multiple tools concurrently

3. **Optimize Parameters**
   - Reduce `max_results` for faster responses
   - Use specific date ranges
   - Filter by categories/providers

## Common Error Messages

### "content blocked"
**Cause:** LLM provider's safety filters triggered
**Solution:** Adjust query or use different provider

### "no candidates in response"
**Cause:** LLM couldn't generate a response
**Solution:** Simplify prompt or check for provider issues

### "response does not contain valid JSON"
**Cause:** Tool call extraction failed
**Solution:** Check LLM response format, enable debug logging

### "timeout"
**Cause:** API request took too long
**Solution:** Reduce query complexity, check network, try again

## Best Practices

1. **Always Enable Debug Logging** when troubleshooting
2. **Check Environment Variables** before running
3. **Test Tools Individually** before using in agents
4. **Monitor API Usage** to avoid rate limits
5. **Use Appropriate Providers** for your use case
6. **Keep Prompts Clear** and unambiguous

## Getting Help

If you continue to experience issues:

1. Check the [examples](../examples/) for working implementations
2. Review the [API documentation](tools/api.md)
3. Enable debug logging and analyze the output
4. Report issues at https://github.com/anthropics/claude-code/issues