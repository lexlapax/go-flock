# Prompt Refactoring Documentation

## Overview

The Research Papers Agent prompt system has been refactored to use a modular architecture that separates core functionality from format-specific instructions.

## Architecture

### Before Refactoring
Previously, the entire prompt was duplicated for each output format:
- `researchPapersPromptMarkdown` - Complete prompt for Markdown format
- `researchPapersPromptJSON` - Complete prompt for JSON format  
- `researchPapersPromptText` - Complete prompt for plain text format

This led to:
- Duplication of core instructions across three constants
- Difficulty maintaining consistency when updating agent behavior
- Risk of prompts diverging over time

### After Refactoring

The new architecture separates concerns:

```go
// Core prompt - shared across all formats
const coreResearchPapersPrompt = `...agent role, tools, and approach...`

// Format-specific instructions
const formatInstructionsMarkdown = `...markdown formatting rules...`
const formatInstructionsJSON = `...JSON schema and rules...`
const formatInstructionsText = `...plain text formatting rules...`

// Combine them dynamically
func getResearchPapersPrompt(format OutputFormat) string {
    formatInstructions := getFormatInstructions(format)
    return coreResearchPapersPrompt + "\n\n" + formatInstructions
}
```

## Benefits

1. **Single Source of Truth**: Core agent behavior defined once
2. **Easy Maintenance**: Update agent capabilities without touching format instructions
3. **Consistency**: All formats share the same core behavior
4. **Extensibility**: Add new output formats by creating new format instructions
5. **Clear Separation**: Tool usage and agent role separate from formatting

## Implementation Details

### Core Prompt Contents
The `coreResearchPapersPrompt` includes:
- Agent role and expertise
- Available tools and their descriptions
- Step-by-step instructions for using tools
- Quality standards and analysis approach

### Format Instructions
Each format constant (`formatInstructionsMarkdown`, `formatInstructionsJSON`, `formatInstructionsText`) contains only:
- Output structure requirements
- Format-specific syntax rules
- Section organization
- Styling guidelines

## Usage Examples

### Modifying Agent Behavior
To change how the agent analyzes papers, edit only `coreResearchPapersPrompt`:
```go
const coreResearchPapersPrompt = `You are a research specialist...
// Add new instruction here - applies to all formats
- Prioritize papers from the last 5 years
...`
```

### Adding a New Output Format
To add a new format (e.g., LaTeX):
```go
// 1. Define format instructions
const formatInstructionsLatex = `Provide your findings as a LaTeX document...`

// 2. Add case to getResearchPapersPrompt
case OutputFormatLatex:
    formatInstructions = formatInstructionsLatex
```

## Testing Considerations

The modular architecture makes testing easier:
- Test core prompt logic independently
- Verify format compliance separately
- Ensure consistent behavior across formats

## Future Enhancements

This pattern can be extended to:
- Support user-provided custom prompts
- Enable prompt versioning
- Allow format customization via configuration
- Create a prompt template system for all agents