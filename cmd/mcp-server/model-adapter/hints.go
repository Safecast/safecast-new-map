package modeladapter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// ModelName identifies an AI model/client.
type ModelName string

const (
	ModelClaude  ModelName = "claude"
	ModelKimi    ModelName = "kimi"
	ModelQwen    ModelName = "qwen"
	ModelGPT     ModelName = "gpt"
	ModelUnknown ModelName = "unknown"
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

// HintsLoader manages loading and caching of model-specific hints from JSON files.
type HintsLoader struct {
	mu      sync.RWMutex
	hints   map[ModelName]*ModelHints
	baseDir string
}

// ModelHints represents the structure of a hints JSON file.
type ModelHints struct {
	Model              string                 `json:"model"`
	DisplayName        string                 `json:"display_name"`
	Capabilities       []string               `json:"capabilities"`
	SystemPrompt       string                 `json:"system_prompt"`
	Tools              map[string]*ToolHints  `json:"tools"`
	GlobalFormatting   map[string]string      `json:"global_formatting_rules"`
	ResponseGuidelines map[string]interface{} `json:"response_guidelines,omitempty"`
}

// ToolHints represents hints for a specific tool.
type ToolHints struct {
	Description     string            `json:"description"`
	DescriptionEn   string            `json:"description_en,omitempty"`
	Examples        []ToolExample     `json:"examples"`
	OutputHints     string            `json:"output_hints"`
	FormattingRules map[string]string `json:"formatting_rules,omitempty"`
	Warnings        []string          `json:"warnings,omitempty"`
	FollowUp        *FollowUpHint     `json:"follow_up,omitempty"`
	Topics          []string          `json:"topics,omitempty"`
}

// ToolExample represents a usage example for a tool.
type ToolExample struct {
	UserQuery string          `json:"user_query"`
	ToolCall  ExampleToolCall `json:"tool_call"`
	FollowUp  *ExampleToolCall `json:"follow_up,omitempty"`
}

// ExampleToolCall represents a tool call in an example.
type ExampleToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// FollowUpHint represents a suggested follow-up tool call.
type FollowUpHint struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// NewHintsLoader creates a new HintsLoader instance.
func NewHintsLoader(baseDir string) *HintsLoader {
	return &HintsLoader{
		hints:   make(map[ModelName]*ModelHints),
		baseDir: baseDir,
	}
}

// Load loads all hints from JSON files in the configured directory.
func (h *HintsLoader) Load() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Define expected hint files
	hintFiles := map[ModelName]string{
		ModelClaude:  "claude.json",
		ModelKimi:    "kimi.json",
		ModelQwen:    "qwen.json",
		ModelGPT:     "gpt.json",
		ModelUnknown: "default.json",
	}

	for modelName, fileName := range hintFiles {
		filePath := filepath.Join(h.baseDir, fileName)
		data, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				// Skip if file doesn't exist, will use fallback
				continue
			}
			return fmt.Errorf("failed to read hints file %s: %w", fileName, err)
		}

		var modelHints ModelHints
		if err := json.Unmarshal(data, &modelHints); err != nil {
			return fmt.Errorf("failed to parse hints file %s: %w", fileName, err)
		}

		h.hints[modelName] = &modelHints
	}

	return nil
}

// GetHints returns the hints for a specific model.
func (h *HintsLoader) GetHints(model ModelName) *ModelHints {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if hints, ok := h.hints[model]; ok {
		return hints
	}
	// Fallback to default/unknown
	if hints, ok := h.hints[ModelUnknown]; ok {
		return hints
	}
	return nil
}

// GetToolHint returns the hint for a specific tool and model.
func (h *HintsLoader) GetToolHint(model ModelName, toolName string) *ToolHints {
	hints := h.GetHints(model)
	if hints == nil || hints.Tools == nil {
		return nil
	}
	return hints.Tools[toolName]
}

// GetAllModels returns a list of all loaded model names.
func (h *HintsLoader) GetAllModels() []ModelName {
	h.mu.RLock()
	defer h.mu.RUnlock()

	models := make([]ModelName, 0, len(h.hints))
	for model := range h.hints {
		models = append(models, model)
	}
	return models
}

// NewDefaultHintProvider returns a provider that yields generic (non-model-specific)
// hints. It can be used as a fallback when the model is unknown.
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
type claudeHintProvider struct{}

func (c *claudeHintProvider) GetHint(toolName string) *ToolHint {
	// future: return per-tool details
	return nil
}

func (c *claudeHintProvider) GetSystemPrompt() string {
	return "You are connected to the Safecast MCP server. Use concise, fact-based answers."
}

func (c *claudeHintProvider) GetDefaultExample(toolName string) string {
	return "" // can be populated later
}

// kimiHintProvider provides hints for Kimi.
type kimiHintProvider struct{}

func (k *kimiHintProvider) GetHint(toolName string) *ToolHint {
	return nil
}

func (k *kimiHintProvider) GetSystemPrompt() string {
	return "You are Kimi, an AI assistant. Follow Kimi-specific formatting rules."
}

func (k *kimiHintProvider) GetDefaultExample(toolName string) string {
	return ""
}

// qwenHintProvider provides hints for Qwen.
type qwenHintProvider struct{}

func (q *qwenHintProvider) GetHint(toolName string) *ToolHint {
	return nil
}

func (q *qwenHintProvider) GetSystemPrompt() string {
	return "You are Qwen. Emphasize technical accuracy and short responses."
}

func (q *qwenHintProvider) GetDefaultExample(toolName string) string {
	return ""
}

// gptHintProvider provides hints for GPT.
type gptHintProvider struct{}

func (g *gptHintProvider) GetHint(toolName string) *ToolHint {
	return nil
}

func (g *gptHintProvider) GetSystemPrompt() string {
	return "You are an AI assistant with access to the Safecast database. Provide concise, factual responses."
}

func (g *gptHintProvider) GetDefaultExample(toolName string) string {
	return ""
}
