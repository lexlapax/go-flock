# Creating Custom Tools

This guide explains how to create custom tools for go-flock that work with go-llms agents.

## Overview

Tools in go-flock are functions that LLM agents can call to perform specific operations. All tools must implement the go-llms `domain.Tool` interface.

## Tool Interface

```go
type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, params interface{}) (interface{}, error)
    ParameterSchema() *domain.Schema
}
```

## Creating a Tool

### Step 1: Define Parameter Types

```go
// Define input parameters
type MyToolParams struct {
    Input    string `json:"input"`
    Option   string `json:"option,omitempty"`
    MaxItems int    `json:"max_items,omitempty"`
}

// Define result type (optional, can use map[string]interface{})
type MyToolResult struct {
    Output string `json:"output"`
    Count  int    `json:"count"`
}
```

### Step 2: Create the Tool Function

```go
func NewMyCustomTool() domain.Tool {
    return tools.NewTool(
        "my_custom_tool",                    // Tool name (verb_object format)
        "Does something useful with input",   // Description
        func(ctx context.Context, params MyToolParams) (*MyToolResult, error) {
            // Validate inputs
            if params.Input == "" {
                return nil, fmt.Errorf("input is required")
            }
            
            // Perform the operation
            result := processInput(params.Input)
            
            // Return the result
            return &MyToolResult{
                Output: result,
                Count:  len(result),
            }, nil
        },
        nil, // Parameter schema (optional for simple types)
    )
}
```

### Step 3: Add Parameter Schema (Optional)

For better validation and LLM understanding:

```go
var MyToolParamSchema = &sdomain.Schema{
    Type: "object",
    Properties: map[string]sdomain.Property{
        "input": {
            Type:        "string",
            Description: "The input to process",
            MinLength:   &[]int{1}[0],
        },
        "option": {
            Type:        "string",
            Description: "Processing option",
            Enum:        []string{"fast", "accurate", "balanced"},
        },
        "max_items": {
            Type:        "integer",
            Description: "Maximum items to return",
            Minimum:     &[]float64{1}[0],
            Maximum:     &[]float64{100}[0],
        },
    },
    Required: []string{"input"},
}
```

## Best Practices

### 1. Naming Convention

Follow the `verb_object` pattern:
- ✅ `get_user_profile`
- ✅ `create_document`
- ✅ `parse_json_data`
- ❌ `user_profile_get`
- ❌ `document_creator`

### 2. Error Handling

```go
func(ctx context.Context, params MyParams) (interface{}, error) {
    // Check context cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    // Validate inputs with descriptive errors
    if params.Input == "" {
        return nil, fmt.Errorf("input parameter is required")
    }
    
    // Wrap errors with context
    data, err := fetchData(params.Input)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch data: %w", err)
    }
    
    return processData(data)
}
```

### 3. Return Consistent Types

```go
// Good: Consistent structure
type ToolResult struct {
    Success bool                   `json:"success"`
    Data    interface{}            `json:"data,omitempty"`
    Error   string                 `json:"error,omitempty"`
    Meta    map[string]interface{} `json:"meta,omitempty"`
}

// Bad: Inconsistent returns
func execute() (interface{}, error) {
    if success {
        return "Success!", nil      // Sometimes string
    }
    return map[string]int{...}, nil // Sometimes map
}
```

### 4. Handle Complex Parameters

```go
// For nested structures
type ComplexParams struct {
    Query    string            `json:"query"`
    Filters  []Filter          `json:"filters"`
    Options  map[string]string `json:"options"`
}

type Filter struct {
    Field    string `json:"field"`
    Operator string `json:"operator"`
    Value    interface{} `json:"value"`
}

// Parse with validation
func parseComplexParams(params interface{}) (*ComplexParams, error) {
    // Handle both direct structs and maps
    switch p := params.(type) {
    case ComplexParams:
        return &p, nil
    case *ComplexParams:
        return p, nil
    case map[string]interface{}:
        // Convert map to struct
        data, _ := json.Marshal(p)
        var result ComplexParams
        if err := json.Unmarshal(data, &result); err != nil {
            return nil, fmt.Errorf("invalid parameters: %w", err)
        }
        return &result, nil
    default:
        return nil, fmt.Errorf("unexpected parameter type: %T", params)
    }
}
```

### 5. File Organization

Organize tools by category:

```
pkg/tools/
├── datetime_tools.go      # All datetime-related tools
├── datetime_tools_test.go
├── file_tools.go         # All file-related tools
├── file_tools_test.go
├── http_tools.go         # All HTTP-related tools
├── http_tools_test.go
└── data_tools.go         # All data processing tools
```

## Complete Example

Here's a complete example of a custom tool:

```go
package tools

import (
    "context"
    "fmt"
    "strings"
    
    "github.com/lexlapax/go-llms/pkg/agent/domain"
    "github.com/lexlapax/go-llms/pkg/agent/tools"
    sdomain "github.com/lexlapax/go-llms/pkg/schema/domain"
)

// TextAnalysisParams defines parameters for text analysis
type TextAnalysisParams struct {
    Text     string `json:"text"`
    Analysis string `json:"analysis"` // word_count, char_count, sentence_count
}

// TextAnalysisResult defines the analysis result
type TextAnalysisResult struct {
    Text     string                 `json:"text"`
    Analysis string                 `json:"analysis"`
    Result   int                    `json:"result"`
    Details  map[string]interface{} `json:"details"`
}

// TextAnalysisParamSchema defines parameter validation
var TextAnalysisParamSchema = &sdomain.Schema{
    Type: "object",
    Properties: map[string]sdomain.Property{
        "text": {
            Type:        "string",
            Description: "Text to analyze",
        },
        "analysis": {
            Type:        "string",
            Description: "Type of analysis to perform",
            Enum:        []string{"word_count", "char_count", "sentence_count"},
        },
    },
    Required: []string{"text", "analysis"},
}

// NewAnalyzeTextTool creates a tool for text analysis
func NewAnalyzeTextTool() domain.Tool {
    return tools.NewTool(
        "analyze_text",
        "Analyzes text and returns statistics",
        func(ctx context.Context, params TextAnalysisParams) (*TextAnalysisResult, error) {
            // Check context
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            default:
            }
            
            // Validate input
            if params.Text == "" {
                return nil, fmt.Errorf("text parameter is required")
            }
            
            result := &TextAnalysisResult{
                Text:     params.Text,
                Analysis: params.Analysis,
                Details:  make(map[string]interface{}),
            }
            
            // Perform analysis
            switch params.Analysis {
            case "word_count":
                words := strings.Fields(params.Text)
                result.Result = len(words)
                result.Details["words"] = words
                
            case "char_count":
                result.Result = len(params.Text)
                result.Details["without_spaces"] = len(strings.ReplaceAll(params.Text, " ", ""))
                
            case "sentence_count":
                // Simple sentence detection
                sentences := 0
                for _, r := range params.Text {
                    if r == '.' || r == '!' || r == '?' {
                        sentences++
                    }
                }
                result.Result = sentences
                
            default:
                return nil, fmt.Errorf("unknown analysis type: %s", params.Analysis)
            }
            
            return result, nil
        },
        TextAnalysisParamSchema,
    )
}
```

## Testing Your Tool

Always write comprehensive tests:

```go
func TestAnalyzeTextTool(t *testing.T) {
    tool := NewAnalyzeTextTool()
    ctx := context.Background()
    
    tests := []struct {
        name    string
        params  TextAnalysisParams
        want    int
        wantErr bool
    }{
        {
            name: "count words",
            params: TextAnalysisParams{
                Text:     "Hello world from go-flock",
                Analysis: "word_count",
            },
            want: 4,
        },
        {
            name: "empty text",
            params: TextAnalysisParams{
                Text:     "",
                Analysis: "word_count",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := tool.Execute(ctx, tt.params)
            if (err != nil) != tt.wantErr {
                t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if err == nil {
                res := result.(*TextAnalysisResult)
                if res.Result != tt.want {
                    t.Errorf("Execute() result = %v, want %v", res.Result, tt.want)
                }
            }
        })
    }
}
```

## Integration with Agents

Once your tool is created, it can be used with any go-llms agent:

```go
agent := workflow.NewAgent(provider)
agent.AddTool(NewAnalyzeTextTool())

// The agent can now use the tool
result, _ := agent.Run(ctx, "Count the words in 'Hello world from go-flock'")
```

## Next Steps

- Review existing tools in `pkg/tools/` for more examples
- Check the [datetime tools](../tools/datetime.md) for comprehensive examples
- See [Creating Custom Agents](./creating-agents.md) to build agents with your tools