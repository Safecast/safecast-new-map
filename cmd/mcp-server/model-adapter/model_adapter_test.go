package modeladapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestDetectModel(t *testing.T) {
	tests := []struct {
		name       string
		userAgent  string
		aiModelHdr string
		want       ModelName
	}{
		{"Claude from User-Agent", "Claude/3.5", "", ModelClaude},
		{"Claude from header", "", "claude", ModelClaude},
		{"Claude mixed case", "Mozilla/5.0 Claude/3.5", "", ModelClaude},
		{"Qwen from User-Agent", "Qwen/2.5", "", ModelQwen},
		{"Qwen from header", "", "qwen", ModelQwen},
		{"Kimi from User-Agent", "Kimi/1.0", "", ModelKimi},
		{"Kimi from header", "", "kimi", ModelKimi},
		{"GPT from User-Agent", "GPT/4", "", ModelGPT},
		{"GPT from header", "", "gpt", ModelGPT},
		{"Unknown model", "UnknownBot/1.0", "", ModelUnknown},
		{"Empty headers", "", "", ModelUnknown},
		{"Header takes precedence", "Mozilla/5.0", "claude", ModelClaude},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/mcp", nil)
			if tt.userAgent != "" {
				req.Header.Set("User-Agent", tt.userAgent)
			}
			if tt.aiModelHdr != "" {
				req.Header.Set("X-AI-Model", tt.aiModelHdr)
			}

			got := detectModel(req)
			if got != tt.want {
				t.Errorf("detectModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelDetectionMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		model := ModelFromContext(r.Context())
		if model == ModelUnknown {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := ModelDetectionMiddleware(handler)

	req := httptest.NewRequest("GET", "/mcp", nil)
	req.Header.Set("User-Agent", "Claude/3.5")
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestHintsLoader_Load(t *testing.T) {
	// Create temp directory with test hints
	tmpDir := t.TempDir()

	// Create a minimal valid hints file
	hintsContent := `{
		"model": "test",
		"display_name": "Test Model",
		"capabilities": ["tool_use"],
		"system_prompt": "Test system prompt",
		"tools": {
			"ping": {
				"description": "Test ping description",
				"examples": [],
				"output_hints": "Test output hints"
			}
		},
		"global_formatting_rules": {
			"tone": "test"
		}
	}`

	if err := os.WriteFile(filepath.Join(tmpDir, "default.json"), []byte(hintsContent), 0644); err != nil {
		t.Fatalf("Failed to write test hints: %v", err)
	}

	loader := NewHintsLoader(tmpDir)
	if err := loader.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify hints were loaded
	hints := loader.GetHints(ModelUnknown)
	if hints == nil {
		t.Fatal("Expected hints to be loaded")
	}

	if hints.Model != "test" {
		t.Errorf("Expected model 'test', got '%s'", hints.Model)
	}

	if hints.DisplayName != "Test Model" {
		t.Errorf("Expected display_name 'Test Model', got '%s'", hints.DisplayName)
	}
}

func TestHintsLoader_GetToolHint(t *testing.T) {
	tmpDir := t.TempDir()

	hintsContent := `{
		"model": "test",
		"display_name": "Test",
		"capabilities": [],
		"system_prompt": "",
		"tools": {
			"query_radiation": {
				"description": "Custom radiation query description",
				"examples": [],
				"output_hints": "Format as table"
			}
		},
		"global_formatting_rules": {}
	}`

	os.WriteFile(filepath.Join(tmpDir, "default.json"), []byte(hintsContent), 0644)

	loader := NewHintsLoader(tmpDir)
	loader.Load()

	hint := loader.GetToolHint(ModelUnknown, "query_radiation")
	if hint == nil {
		t.Fatal("Expected tool hint for query_radiation")
	}

	if hint.Description != "Custom radiation query description" {
		t.Errorf("Unexpected description: %s", hint.Description)
	}
}

func TestAdapter_EnhanceTool(t *testing.T) {
	tmpDir := t.TempDir()

	hintsContent := `{
		"model": "claude",
		"display_name": "Claude",
		"capabilities": [],
		"system_prompt": "Test prompt",
		"tools": {
			"ping": {
				"description": "Enhanced ping description",
				"examples": [],
				"output_hints": ""
			}
		},
		"global_formatting_rules": {}
	}`

	os.WriteFile(filepath.Join(tmpDir, "claude.json"), []byte(hintsContent), 0644)

	loader := NewHintsLoader(tmpDir)
	loader.Load()

	adapter := NewAdapter()
	adapter.SetHintsLoader(loader)

	// Create a test tool
	originalTool := &mcp.Tool{
		Name:        "ping",
		Description: "Original description",
	}

	// Create context with Claude model
	ctx := contextWithModel(context.Background(), ModelClaude)

	enhancedTool := adapter.EnhanceTool(ctx, originalTool)

	if enhancedTool == nil {
		t.Fatal("Expected enhanced tool")
	}

	if enhancedTool.Description != "Enhanced ping description" {
		t.Errorf("Expected enhanced description, got: %s", enhancedTool.Description)
	}

	// Verify original tool is unchanged
	if originalTool.Description != "Original description" {
		t.Error("Original tool should not be modified")
	}
}

func TestAdapter_GetSystemPrompt(t *testing.T) {
	tmpDir := t.TempDir()

	hintsContent := `{
		"model": "claude",
		"display_name": "Claude",
		"capabilities": [],
		"system_prompt": "Custom Claude system prompt",
		"tools": {},
		"global_formatting_rules": {}
	}`

	os.WriteFile(filepath.Join(tmpDir, "claude.json"), []byte(hintsContent), 0644)

	loader := NewHintsLoader(tmpDir)
	loader.Load()

	adapter := NewAdapter()
	adapter.SetHintsLoader(loader)

	ctx := contextWithModel(context.Background(), ModelClaude)
	prompt := adapter.GetSystemPrompt(ctx)

	if prompt != "Custom Claude system prompt" {
		t.Errorf("Expected custom prompt, got: %s", prompt)
	}
}

func TestAdapter_GetModelDisplayName(t *testing.T) {
	tmpDir := t.TempDir()

	hintsContent := `{
		"model": "claude",
		"display_name": "Claude 3.5 Sonnet",
		"capabilities": [],
		"system_prompt": "",
		"tools": {},
		"global_formatting_rules": {}
	}`

	os.WriteFile(filepath.Join(tmpDir, "claude.json"), []byte(hintsContent), 0644)

	loader := NewHintsLoader(tmpDir)
	loader.Load()

	adapter := NewAdapter()
	adapter.SetHintsLoader(loader)

	ctx := contextWithModel(context.Background(), ModelClaude)
	name := adapter.GetModelDisplayName(ctx)

	if name != "Claude 3.5 Sonnet" {
		t.Errorf("Expected 'Claude 3.5 Sonnet', got: %s", name)
	}
}

func TestHintsLoader_FallbackToDefault(t *testing.T) {
	tmpDir := t.TempDir()

	// Only create default.json, no model-specific files
	hintsContent := `{
		"model": "default",
		"display_name": "Default",
		"capabilities": [],
		"system_prompt": "Default prompt",
		"tools": {},
		"global_formatting_rules": {}
	}`

	os.WriteFile(filepath.Join(tmpDir, "default.json"), []byte(hintsContent), 0644)

	loader := NewHintsLoader(tmpDir)
	loader.Load()

	// Request unknown model - should fall back to default
	hints := loader.GetHints(ModelUnknown)
	if hints == nil {
		t.Fatal("Expected fallback to default hints")
	}

	if hints.Model != "default" {
		t.Errorf("Expected fallback to 'default', got: %s", hints.Model)
	}
}
