package modeladapter

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

// Adapter encapsulates the model hint logic and provides helpers to register
// tools with model-specific hints and prompts.
type Adapter struct {
	hints         map[ModelName]ModelHintProvider
	defaultProvider ModelHintProvider
	hintsLoader   *HintsLoader
}

// NewAdapter constructs a new Adapter with built-in provider set.
func NewAdapter() *Adapter {
	return &Adapter{
		hints: map[ModelName]ModelHintProvider{
			ModelClaude:  NewClaudeHintProvider(),
			ModelKimi:    NewKimiHintProvider(),
			ModelQwen:    NewQwenHintProvider(),
			ModelGPT:     NewGPTHintProvider(),
		},
		defaultProvider: NewDefaultHintProvider(),
		hintsLoader:     nil, // Set via SetHintsLoader if JSON hints are used
	}
}

// SetHintsLoader sets the HintsLoader for JSON-based hints.
func (a *Adapter) SetHintsLoader(loader *HintsLoader) {
	a.hintsLoader = loader
}

// EnhanceTool returns a copy of the provided tool definition with hints
// injected for the model that has been detected in ctx. If no model is
// available or a hint provider is not registered, the unmodified tool is
// returned.
func (a *Adapter) EnhanceTool(ctx context.Context, tool *mcp.Tool) *mcp.Tool {
	model := ModelFromContext(ctx)
	
	// Try JSON hints first if loader is available
	if a.hintsLoader != nil {
		if toolHints := a.hintsLoader.GetToolHint(model, tool.Name); toolHints != nil {
			if toolHints.Description != "" {
				newTool := *tool
				newTool.Description = toolHints.Description
				return &newTool
			}
		}
	}
	
	// Fallback to provider-based hints
	provider, ok := a.hints[model]
	if !ok {
		provider = a.defaultProvider
	}

	if tool == nil {
		return nil
	}

	hint := provider.GetHint(tool.Name)
	if hint == nil {
		return tool
	}

	// create a shallow copy with overridden description
	newTool := *tool
	if hint.Description != "" {
		newTool.Description = hint.Description
	}

	return &newTool
}

// EnrichResult modifies the tool result to include model-specific formatting
// hints and metadata.
func (a *Adapter) EnrichResult(ctx context.Context, res *mcp.CallToolResult) *mcp.CallToolResult {
	if res == nil || len(res.Content) == 0 {
		return res
	}

	model := ModelFromContext(ctx)
	
	// Get model-specific formatting hints
	var globalFormatting map[string]string
	
	if a.hintsLoader != nil {
		if hints := a.hintsLoader.GetHints(model); hints != nil {
			globalFormatting = hints.GlobalFormatting
		}
	}
	
	// Add formatting metadata to the result
	if len(globalFormatting) > 0 {
		// Find the text content and append formatting info
		for i, content := range res.Content {
			if tc, ok := content.(mcp.TextContent); ok {
				// Add formatting hints as a JSON comment at the end
				formattingJSON, _ := json.Marshal(map[string]interface{}{
					"_formatting_hints": globalFormatting,
					"_model":            model,
				})
				tc.Text = tc.Text + "\n\n<!-- " + string(formattingJSON) + " -->"
				res.Content[i] = tc
			}
		}
	}

	return res
}

// GetSystemPrompt returns the system prompt for the detected model.
func (a *Adapter) GetSystemPrompt(ctx context.Context) string {
	model := ModelFromContext(ctx)
	
	// Try JSON hints first
	if a.hintsLoader != nil {
		if hints := a.hintsLoader.GetHints(model); hints != nil && hints.SystemPrompt != "" {
			return hints.SystemPrompt
		}
	}
	
	// Fallback to provider
	provider, ok := a.hints[model]
	if !ok {
		provider = a.defaultProvider
	}
	
	return provider.GetSystemPrompt()
}

// GetModelDisplayName returns a human-readable name for the detected model.
func (a *Adapter) GetModelDisplayName(ctx context.Context) string {
	model := ModelFromContext(ctx)
	
	if a.hintsLoader != nil {
		if hints := a.hintsLoader.GetHints(model); hints != nil && hints.DisplayName != "" {
			return hints.DisplayName
		}
	}
	
	// Default display names
	switch model {
	case ModelClaude:
		return "Claude"
	case ModelKimi:
		return "Kimi"
	case ModelQwen:
		return "Qwen"
	case ModelGPT:
		return "GPT"
	default:
		return "Unknown Model"
	}
}
