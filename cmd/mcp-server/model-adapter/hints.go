package modeladapter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ModelName identifies an AI model/client.
type ModelName string

const (
	ModelClaude   ModelName = "claude"
	ModelKimi     ModelName = "kimi"
	ModelQwen     ModelName = "qwen"
	ModelGPT      ModelName = "gpt"
	ModelUnknown  ModelName = "unknown"
)

// ToolHint contains model-specific guidance for a particular tool.
type ToolHint struct {
	Description     string            // Optional override for tool description.
	Examples        []string          // Model-specific usage examples.
	FormattingRules []string          // Output formatting hints or rules.
	Warnings        []string          // Known model-specific gotchas.
	Metadata        map[string]string // Arbitrary key/value metadata.
}

// ModelHintProvider returns hints and prompts tailored to a model.
type ModelHintProvider interface {
	// GetHint returns the hint for the named tool. May return nil if none.
	GetHint(toolName string) *ToolHint

	// GetSystemPrompt returns a model-specific system prompt for the MCP session.
	GetSystemPrompt() string

	// GetDefaultExample returns a single example for the named tool.
	GetDefaultExample(toolName string) string
}

// NewDefaultHintProvider returns a provider that yields generic (non-model-specific) hints.
// It can be used as a fallback when the model is unknown.
func NewDefaultHintProvider() ModelHintProvider {
	return &defaultHintProvider{}
}

// NewClaudeHintProvider returns a provider with hints tailored to Claude.
func NewClaudeHintProvider() ModelHintProvider {
	return &claudeHintProvider{}
}

// NewKimiHintProvider returns a provider with hints tailored to Kimi.
func NewKimiHintProvider() ModelHintProvider {
	return &kimiHintProvider{}
}

// NewQwenHintProvider returns a provider with hints tailored to Qwen.
func NewQwenHintProvider() ModelHintProvider {
	return &qwenHintProvider{}
}

// NewGPTHintProvider returns a provider with hints tailored to GPT.
func NewGPTHintProvider() ModelHintProvider {
	return &gptHintProvider{}
}

// defaultHintProvider is a trivial implementation.
type defaultHintProvider struct{}

func (d *defaultHintProvider) GetHint(toolName string) *ToolHint {
	return nil
}

func (d *defaultHintProvider) GetSystemPrompt() string {
	return "" // no special system prompt
}

func (d *defaultHintProvider) GetDefaultExample(toolName string) string {
	return ""
}

// claudeHintProvider provides hints for Claude-based clients.
type claudeHintProvider struct {
	hints *Hints
}

func (c *claudeHintProvider) GetHint(toolName string) *ToolHint {
	if c.hints == nil || c.hints.Tools == nil {
		return nil
	}
	tool, ok := c.hints.Tools[toolName]
	if !ok {
		return nil
	}
	return &ToolHint{
		Description: tool.Description,
	}
}

func (c *claudeHintProvider) GetSystemPrompt() string {
	if c.hints != nil && c.hints.SystemPrompt != "" {
		return c.hints.SystemPrompt
	}
	return "You are connected to the Safecast MCP server. Use concise, fact-based answers."
}

func (c *claudeHintProvider) GetDefaultExample(toolName string) string {
	return "" // can be populated later
}

// kimiHintProvider provides hints for Kimi.
type kimiHintProvider struct {
	hints *Hints
}

func (k *kimiHintProvider) GetHint(toolName string) *ToolHint {
	if k.hints == nil || k.hints.Tools == nil {
		return nil
	}
	tool, ok := k.hints.Tools[toolName]
	if !ok {
		return nil
	}
	return &ToolHint{
		Description: tool.Description,
	}
}

func (k *kimiHintProvider) GetSystemPrompt() string {
	if k.hints != nil && k.hints.SystemPrompt != "" {
		return k.hints.SystemPrompt
	}
	return "You are Kimi, an AI assistant. Follow Kimi-specific formatting rules."
}

func (k *kimiHintProvider) GetDefaultExample(toolName string) string {
	return ""
}

// qwenHintProvider provides hints for Qwen.
type qwenHintProvider struct {
	hints *Hints
}

func (q *qwenHintProvider) GetHint(toolName string) *ToolHint {
	if q.hints == nil || q.hints.Tools == nil {
		return nil
	}
	tool, ok := q.hints.Tools[toolName]
	if !ok {
		return nil
	}
	return &ToolHint{
		Description: tool.Description,
	}
}

func (q *qwenHintProvider) GetSystemPrompt() string {
	if q.hints != nil && q.hints.SystemPrompt != "" {
		return q.hints.SystemPrompt
	}
	return "You are Qwen. Emphasize technical accuracy and short responses."
}

func (q *qwenHintProvider) GetDefaultExample(toolName string) string {
	return ""
}

// gptHintProvider provides hints for GPT.
type gptHintProvider struct {
	hints *Hints
}

func (g *gptHintProvider) GetHint(toolName string) *ToolHint {
	if g.hints == nil || g.hints.Tools == nil {
		return nil
	}
	tool, ok := g.hints.Tools[toolName]
	if !ok {
		return nil
	}
	return &ToolHint{
		Description: tool.Description,
	}
}

func (g *gptHintProvider) GetSystemPrompt() string {
	if g.hints != nil && g.hints.SystemPrompt != "" {
		return g.hints.SystemPrompt
	}
	return "You are an AI assistant with access to the Safecast radiation monitoring database. Use tools to retrieve data. Provide concise, factual responses."
}

func (g *gptHintProvider) GetDefaultExample(toolName string) string {
	return ""
}

// Hints manages per-model hint configuration loaded from JSON files.
type Hints struct {
	Model              string                 `json:"model"`
	DisplayName        string                 `json:"display_name"`
	Capabilities       []string               `json:"capabilities"`
	SystemPrompt       string                 `json:"system_prompt"`
	Tools              map[string]ToolConfig  `json:"tools"`
	GlobalFormatting   map[string]interface{} `json:"global_formatting_rules"`
}

// ToolConfig describes a tool's hint configuration.
type ToolConfig struct {
	Description     string                   `json:"description"`
	Examples        []map[string]interface{} `json:"examples"`
	OutputHints     string                   `json:"output_hints"`
	FormattingRules map[string]interface{}   `json:"formatting_rules"`
	Warnings        []string                 `json:"warnings"`
}

// HintsLoader loads and manages hint configurations from JSON files.
type HintsLoader struct {
	hintsDir        string
	hints           map[ModelName]*Hints
	claudeProvider  *claudeHintProvider
	kimiProvider    *kimiHintProvider
	qwenProvider    *qwenHintProvider
	gptProvider     *gptHintProvider
}

// NewHintsLoader creates a new HintsLoader for the given directory.
func NewHintsLoader(hintsDir string) *HintsLoader {
	loader := &HintsLoader{
		hintsDir:       hintsDir,
		hints:          make(map[ModelName]*Hints),
		claudeProvider: &claudeHintProvider{},
		kimiProvider:   &kimiHintProvider{},
		qwenProvider:   &qwenHintProvider{},
		gptProvider:    &gptHintProvider{},
	}
	return loader
}

// Load reads all JSON hint files from the configured directory.
func (h *HintsLoader) Load() error {
	files := []struct {
		model ModelName
		file  string
	}{
		{ModelClaude, "claude.json"},
		{ModelKimi, "kimi.json"},
		{ModelQwen, "qwen.json"},
		{ModelGPT, "gpt.json"},
		{ModelUnknown, "default.json"},
	}

	for _, f := range files {
		path := filepath.Join(h.hintsDir, f.file)
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Skip if file doesn't exist
			}
			return fmt.Errorf("reading hint file %s: %w", path, err)
		}

		var hints Hints
		if err := json.Unmarshal(data, &hints); err != nil {
			return fmt.Errorf("parsing hint file %s: %w", path, err)
		}

		h.hints[f.model] = &hints

		// Update providers with loaded hints
		switch f.model {
		case ModelClaude:
			h.claudeProvider.hints = &hints
		case ModelKimi:
			h.kimiProvider.hints = &hints
		case ModelQwen:
			h.qwenProvider.hints = &hints
		case ModelGPT:
			h.gptProvider.hints = &hints
		}
	}

	return nil
}

// GetProvider returns the hint provider for a model.
func (h *HintsLoader) GetProvider(model ModelName) ModelHintProvider {
	switch model {
	case ModelClaude:
		return h.claudeProvider
	case ModelKimi:
		return h.kimiProvider
	case ModelQwen:
		return h.qwenProvider
	case ModelGPT:
		return h.gptProvider
	default:
		return &defaultHintProvider{}
	}
}

// GetHints returns the hints configuration for a model.
func (h *HintsLoader) GetHints(model ModelName) *Hints {
	return h.hints[model]
}

// GetToolHint returns the tool-specific hint for a model.
func (h *HintsLoader) GetToolHint(model ModelName, toolName string) *ToolConfig {
	hints := h.hints[model]
	if hints == nil {
		return nil
	}
	tool, ok := hints.Tools[toolName]
	if !ok {
		return nil
	}
	return &tool
}

// GetAllModels returns a list of all loaded model names.
func (h *HintsLoader) GetAllModels() []ModelName {
	models := make([]ModelName, 0, len(h.hints))
	for model := range h.hints {
		models = append(models, model)
	}
	return models
}
