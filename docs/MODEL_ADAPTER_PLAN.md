# Model Adapter Layer Implementation Plan

**Branch:** `model-adapter-layer`
**Date:** 2026-03-07
**Status:** 🗄️ DEPRECATED - Superseded by [MCP_MODEL_ADAPTER_PLAN.md](MCP_MODEL_ADAPTER_PLAN.md)

---

> **Note:** This document has been superseded by [MCP_MODEL_ADAPTER_PLAN.md](MCP_MODEL_ADAPTER_PLAN.md), which reflects the actual implementation in the unified server architecture. This document is kept for historical reference.

---

## Overview

Create a model adapter layer that detects which AI model is calling the MCP server and injects model-specific hints/examples into tool definitions and responses.

## Goals

1. **Model Detection**: Identify which AI model/client is making requests (Claude, Kimi, Qwen, etc.)
2. **Dynamic Hints**: Serve model-specific hints, examples, and formatting instructions
3. **Backward Compatible**: Existing clients continue working without changes
4. **Minimal Performance Impact**: No significant latency added to tool calls

## Current State Analysis

### Architecture
- **Language:** Go (not JavaScript as initially proposed)
- **Library:** `github.com/mark3labs/mcp-go`
- **Transport:** SSE (`/mcp/sse`, `/mcp/message`) and Streamable HTTP (`/mcp-http`)
- **Tools:** 17 MCP tools with hardcoded `_ai_hint` annotations

### Key Files
| File                              | Purpose                                |
| --------------------------------- | -------------------------------------- |
| `go/cmd/mcp-server/main.go`       | Server entry point, tool registration  |
| `go/cmd/mcp-server/tool_*.go`     | Individual tool definitions (17 files) |
| `go/cmd/mcp-server/api_client.go` | REST API client                        |
| `go/cmd/mcp-server/ai_logging.go` | AI session logging                     |

### Current Limitations
1. **No model detection** - Server hardcodes `"claude-client"` in logging
2. **Static hints** - `_ai_hint` fields are identical for all models
3. **MCP library doesn't expose HTTP headers** - Model info not available in `CallToolRequest`

## Proposed Architecture

```
go/cmd/mcp-server/
├── model-adapter/
│   ├── model_adapter.go      # Core adapter logic, hint selection
│   ├── model_detection.go    # Model detection from various signals
│   ├── hints/
│   │   ├── claude_hints.go   # Claude-specific hints
│   │   ├── kimi_hints.go     # Kimi-specific hints
│   │   └── qwen_hints.go     # Qwen-specific hints
│   └── hints.go              # Common hint types and interfaces
├── middleware/
│   └── model_middleware.go   # HTTP middleware for model detection
└── ...existing files...
```

## Model Detection Strategies

### Option 1: HTTP Middleware (Recommended)
**Mechanism:** Custom HTTP middleware wraps MCP server, extracts model from headers before passing to MCP handler.

**Detection Signals:**
- `User-Agent` header (e.g., `Claude-Bot/1.0`, `Kimi-Assistant/2.0`)
- Custom headers (e.g., `X-AI-Model: claude-sonnet-4-5-20250929`)
- TLS client certificate CN (if deployed with mTLS)
- Source IP ranges (per-client allowlists)

**Pros:**
- Clean separation of concerns
- Works with any MCP transport
- Can cache detection per session

**Cons:**
- Requires wrapping the MCP server
- May need upstream library changes

### Option 2: Query Parameter
**Mechanism:** Clients pass `?model=claude` in MCP_BASE_URL.

**Example:**
```bash
MCP_BASE_URL="http://localhost:3333?model=claude"
```

**Pros:**
- Simple to implement
- No middleware needed

**Cons:**
- Requires client configuration changes
- Less secure (visible in logs)

### Option 3: Environment Variable (Per-Deployment)
**Mechanism:** Each deployment sets `AI_MODEL_NAME=claude` in environment.

**Pros:**
- Trivial to implement
- No code changes for clients

**Cons:**
- One model per deployment
- Not flexible for multi-tenant servers

### Recommended Approach
**Option 1 (HTTP Middleware)** with **Option 3 (Environment Variable)** as fallback.

## Implementation Tasks

### Phase 1: Core Infrastructure
- [ ] Create `model-adapter/` directory structure
- [ ] Define `ModelHint` interface and data structures
- [ ] Implement `ModelDetector` with middleware approach
- [ ] Add session-based model caching (per `sessionId`)

### Phase 2: Model-Specific Hints
- [ ] Extract existing `_ai_hint` from all 17 tool files
- [ ] Create `hints/claude_hints.go` (baseline from current hints)
- [ ] Create `hints/kimi_hints.go` (adapted for Kimi's strengths)
- [ ] Create `hints/qwen_hints.go` (adapted for Qwen's strengths)
- [ ] Document hint differences and rationale

### Phase 3: Integration
- [ ] Modify `main.go` to use model adapter at tool registration
- [ ] Update `instrument()` wrapper to detect and log model
- [ ] Add middleware to HTTP server setup
- [ ] Ensure fallback to default hints if model unknown

### Phase 4: Testing
- [ ] Unit tests for model detection logic
- [ ] Integration tests with simulated clients
- [ ] Performance benchmarks (ensure <5ms overhead)
- [ ] Test backward compatibility (no model = default hints)

## Data Structures

```go
// model-adapter/hints.go
type ModelName string

const (
    ModelClaude ModelName = "claude"
    ModelKimi   ModelName = "kimi"
    ModelQwen   ModelName = "qwen"
    ModelUnknown ModelName = "unknown"
)

type ToolHint struct {
    Description     string            // Tool description override
    Examples        []string          // Model-specific examples
    FormattingRules []string          // Output formatting rules
    Warnings        []string          // Model-specific gotchas
    Metadata        map[string]string // Additional context
}

type ModelHintProvider interface {
    GetHint(toolName string) *ToolHint
    GetSystemPrompt() string
    GetDefaultExample(toolName string) string
}
```

## Example Usage

```go
// In main.go tool registration
adapter := modeladapter.NewAdapter()

mcpServer.AddTool(
    adapter.EnhanceTool(queryRadiationToolDef, "query_radiation"),
    instrument("query_radiation", handleQueryRadiation),
)

// In middleware
func ModelDetectionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        model := detectModel(r)  // Check headers, etc.
        ctx := context.WithValue(r.Context(), ModelContextKey, model)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## Migration Path

### Existing Deployments
- No breaking changes
- Unknown models receive default hints
- Gradual rollout via feature flag (`ENABLE_MODEL_ADAPTER=true`)

### New Clients
- Can advertise model via `User-Agent` or `X-AI-Model` header
- Receive optimized hints for their model
- Better response quality and formatting

## Success Metrics

| Metric                    | Target                     |
| ------------------------- | -------------------------- |
| Model detection accuracy  | >95%                       |
| Added latency per request | <5ms                       |
| Code coverage (new code)  | >80%                       |
| Backward compatibility    | 100% (no breaking changes) |

## Risks and Mitigations

| Risk                                 | Impact | Mitigation                                  |
| ------------------------------------ | ------ | ------------------------------------------- |
| MCP library changes break middleware | High   | Pin library version, add integration tests  |
| Model detection false positives      | Medium | Conservative detection, fallback to default |
| Performance regression               | Medium | Benchmark before/after, optimize hot paths  |
| Hint maintenance burden              | Low    | Document clearly, keep hints DRY            |

## Next Steps

1. **Create directory structure** and stub files
2. **Implement model detection middleware** (Option 1)
3. **Extract and refactor existing hints** from tool files
4. **Build model-specific hint providers**
5. **Test with real clients** (Claude.ai, etc.)

## References

- [MCP Go Library](https://github.com/mark3labs/mcp-go)
- [Current Tool Definitions](../go/cmd/mcp-server/tool_*.go)
- [Conversation Notes](../conversation-notes.md)
