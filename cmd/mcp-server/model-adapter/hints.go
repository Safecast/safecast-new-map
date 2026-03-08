package modeladapter

// ModelName identifies an AI model/client.
type ModelName string

const (
	ModelClaude   ModelName = "claude"
	ModelKimi     ModelName = "kimi"
	ModelQwen     ModelName = "qwen"
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

// NewDefaultHintProvider returns a provider that yields generic (non-model-specific)
hints. It can be used as a fallback when the model is unknown.
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
