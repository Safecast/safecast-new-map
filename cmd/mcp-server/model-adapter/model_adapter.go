package modeladapter

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

// Adapter encapsulates the model hint logic and provides helpers to register
// tools with model-specific hints and prompts.
type Adapter struct {
	hints map[ModelName]ModelHintProvider
	defaultProvider ModelHintProvider
	hintsLoader *HintsLoader
}

// NewAdapter constructs a new Adapter with built-in provider set.
func NewAdapter() *Adapter {
	return &Adapter{
		hints: map[ModelName]ModelHintProvider{
			ModelClaude:  NewClaudeHintProvider(),
			ModelKimi:    NewKimiHintProvider(),
			ModelQwen:    NewQwenHintProvider(),
		},
		defaultProvider: NewDefaultHintProvider(),
	}
}

// SetHintsLoader sets the hints loader for JSON-based hints.
func (a *Adapter) SetHintsLoader(loader *HintsLoader) {
	a.hintsLoader = loader
}

// EnhanceTool returns a copy of the provided tool definition with hints
// injected for the model that has been detected in ctx.  If no model is
// available or a hint provider is not registered, the unmodified tool is
// returned.
func (a *Adapter) EnhanceTool(ctx context.Context, tool *mcp.Tool) *mcp.Tool {
	model := ModelFromContext(ctx)
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

	// create a shallow copy with overridden description or examples
	newTool := *tool
	if hint.Description != "" {
		newTool.Description = hint.Description
	}
	// The MCP library currently lacks an "examples" field; consider adding
	// via custom metadata or modifying the library when needed.

	return &newTool
}

// EnrichResult allows tools to modify their result based on model-specific
// formatting rules; placeholder for future work.
func (a *Adapter) EnrichResult(ctx context.Context, res *mcp.CallToolResult) *mcp.CallToolResult {
	// TODO: implement
	return res
}
