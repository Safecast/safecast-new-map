package modeladapter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
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

// HintsLoader manages loading model-specific hints from JSON files.
type HintsLoader struct {
	hintsDir string
	hints    map[string]*ModelHint // keyed by model name
	mu       sync.RWMutex
}

// ModelHint represents the JSON structure for model hints.
type ModelHint struct {
	Model               string                       `json:"model"`
	DisplayName         string                       `json:"display_name"`
	Capabilities        []string                     `json:"capabilities"`
	SystemPrompt        string                       `json:"system_prompt"`
	Tools               map[string]*ToolHintJSON     `json:"tools"`
	GlobalFormattingRules map[string]string          `json:"global_formatting_rules"`
}

// ToolHintJSON represents the JSON structure for a tool hint.
type ToolHintJSON struct {
	Description       string            `json:"description"`
	Examples          []json.RawMessage `json:"examples"`
	OutputHints       string            `json:"output_hints"`
	FormattingRules   map[string]string `json:"formatting_rules"`
	Warnings          []string          `json:"warnings"`
	Metadata          map[string]string `json:"metadata"`
}

// NewHintsLoader creates a new HintsLoader for the given directory.
func NewHintsLoader(hintsDir string) *HintsLoader {
	return &HintsLoader{
		hintsDir: hintsDir,
		hints:    make(map[string]*ModelHint),
	}
}

// Load reads all JSON hint files from the hints directory.
func (h *HintsLoader) Load() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	entries, err := os.ReadDir(h.hintsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(h.hintsDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var hint ModelHint
		if err := json.Unmarshal(data, &hint); err != nil {
			continue
		}

		// Map "default.json" to ModelUnknown for fallback
		modelKey := hint.Model
		if entry.Name() == "default.json" {
			modelKey = "unknown"
		}
		h.hints[modelKey] = &hint
	}

	return nil
}

// GetHints returns the hint for a specific model (for testing).
func (h *HintsLoader) GetHints(model ModelName) *ModelHint {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.hints[string(model)]
}

// GetAllModels returns a list of all loaded model names.
func (h *HintsLoader) GetAllModels() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	models := make([]string, 0, len(h.hints))
	for model := range h.hints {
		models = append(models, model)
	}
	return models
}

// GetHint returns the hint for a specific model.
func (h *HintsLoader) GetHint(model string) *ModelHint {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.hints[model]
}

// GetToolHint returns the hint for a specific tool and model.
func (h *HintsLoader) GetToolHint(model ModelName, toolName string) *ToolHintJSON {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if hint, ok := h.hints[string(model)]; ok {
		return hint.Tools[toolName]
	}
	return nil
}

// GetSystemPrompt returns the system prompt for a specific model.
func (h *HintsLoader) GetSystemPrompt(model ModelName) string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if hint, ok := h.hints[string(model)]; ok {
		return hint.SystemPrompt
	}
	return ""
}
