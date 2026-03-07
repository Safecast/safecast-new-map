# MCP Model Adapter Plan

**Goal:** Add a model detection and hints injection layer to the Go-based MCP server that serves model-specific tool descriptions and examples based on the calling AI client.

**Repository:** https://github.com/Safecast/safecast-new-map

---

## Current Architecture

```
┌──────────────┐     ┌─────────────────────┐     ┌────────────┐
│  AI Client   │────▶│  Go MCP Server      │────▶│ PostgreSQL │
│  (Claude,    │◀────│  (cmd/mcp-server/)  │◀────│  Database  │
│   Qwen, etc) │     │  - main.go          │     │            │
└──────────────┘     │  - tool_*.go (18)   │     └────────────┘
                     │  - rest.go          │
                     │  - db_client.go     │
                     └─────────────────────┘
```

**Current state:**
- 18 MCP tools defined in `cmd/mcp-server/tool_*.go`
- Tool descriptions hardcoded in Go (e.g., `mcp.NewTool("query_radiation", mcp.WithDescription(...))`)
- SSE transport at `/mcp/sse` and Streamable HTTP at `/mcp-http`
- No model detection or dynamic hints

---

## Target Architecture

```
┌──────────────┐     ┌─────────────────────┐     ┌───────────────────┐     ┌────────────┐
│  AI Client   │────▶│  Model Adapter      │────▶│  Go MCP Server    │────▶│ PostgreSQL │
│  + headers   │◀────│  (Go, new layer)    │◀────│  (existing)       │◀────│  Database  │
└──────────────┘     │  - Detect model     │     │  - Tool execution │     └────────────┘
                     │  - Inject hints     │     └───────────────────┘
                     │  - Load JSON hints  │
                     └─────────────────────┘
                            │
                            ▼
                     ┌───────────────────┐
                     │  hints/           │
                     │  - claude.json    │
                     │  - qwen.json      │
                     │  - kimi.json      │
                     │  - gpt.json       │
                     │  - default.json   │
                     └───────────────────┘
```

---

## Implementation Phases

### Phase 1: Model Detection Layer

**File:** `cmd/mcp-server/model_adapter.go` (new)

**Responsibilities:**
1. Extract model name from HTTP request headers
2. Map model names to hint profiles
3. Provide model context to tool handlers

**Header sources to check:**
- `User-Agent` (e.g., "Claude/3.5", "Qwen/2.5")
- `X-Model-Name` (custom header)
- `X-Client-Name` (custom header)
- MCP request metadata (if available in context)

**Model detection logic:**
```go
type ModelProfile struct {
    Name        string   // "claude", "qwen", "kimi", "gpt"
    DisplayName string   // "Claude 3.5 Sonnet"
    Capabilities []string // ["long_context", "json_mode", "tool_use_v2"]
}

func DetectModel(r *http.Request) ModelProfile
func GetModelFromContext(ctx context.Context) ModelProfile
```

**Model mapping examples:**
| User-Agent pattern | Model Profile |
|-------------------|---------------|
| `Claude/*` | `claude` |
| `Qwen/*` | `qwen` |
| `Kimi/*` | `kimi` |
| `GPT/*` or `ChatGPT/*` | `gpt` |
| (unknown) | `default` |

---

### Phase 2: Hints JSON Structure

**Directory:** `cmd/mcp-server/hints/` (new)

**File structure:**
```
cmd/mcp-server/
├── hints/
│   ├── claude.json
│   ├── qwen.json
│   ├── kimi.json
│   ├── gpt.json
│   └── default.json
```

**JSON schema:**
```json
{
  "model": "claude",
  "display_name": "Claude 3.5 Sonnet",
  "capabilities": ["long_context", "tool_use", "json_mode"],
  "tools": {
    "query_radiation": {
      "description": "Optimized description for Claude...",
      "examples": [
        {
          "user_query": "What's the radiation level near Fukushima?",
          "tool_call": {
            "name": "query_radiation",
            "arguments": {
              "lat": 37.4267,
              "lon": 140.5367,
              "radius_m": 50000
            }
          }
        }
      ],
      "output_hints": "Present results in markdown tables with map links",
      "follow_up_suggestions": ["Call sensor_current to check for real-time sensors"]
    },
    "get_track": {
      "description": "...",
      "examples": [...]
    }
  },
  "system_prompt_additions": [
    "Always present radiation data in scientific, objective language",
    "Never use personal pronouns when reporting measurements"
  ],
  "formatting_rules": {
    "coordinates": "Always use [lat°N, lon°E](map_link) format",
    "units": "CPM = counts per minute, µSv/h = microsieverts per hour",
    "tables": "Use markdown tables for all measurement lists"
  }
}
```

**Key differences per model:**

| Model | Description Style | Example Count | Special Notes |
|-------|------------------|---------------|---------------|
| `claude` | Detailed, verbose | 2-3 per tool | Handles long context well |
| `qwen` | Concise, structured | 1-2 per tool | Prefers clear formatting |
| `kimi` | Medium length | 1-2 per tool | Good at following examples |
| `gpt` | Very concise | 1 per tool | Token-optimized |
| `default` | Generic | 1 per tool | Fallback for unknown models |

---

### Phase 3: Hints Loader

**File:** `cmd/mcp-server/hints_loader.go` (new)

**Responsibilities:**
1. Load all hint JSON files at startup
2. Cache in memory (map[string]ModelHints)
3. Provide thread-safe access via `GetHints(modelName string)`
4. Hot-reload support (optional, via file watcher)

**API:**
```go
type ModelHints struct {
    Model              string                    `json:"model"`
    DisplayName        string                    `json:"display_name"`
    Capabilities       []string                  `json:"capabilities"`
    Tools              map[string]ToolHints      `json:"tools"`
    SystemPromptAdditions []string               `json:"system_prompt_additions"`
    FormattingRules    map[string]string         `json:"formatting_rules"`
}

type ToolHints struct {
    Description          string      `json:"description"`
    Examples             []Example   `json:"examples"`
    OutputHints          string      `json:"output_hints"`
    FollowUpSuggestions  []string    `json:"follow_up_suggestions"`
}

type Example struct {
    UserQuery string      `json:"user_query"`
    ToolCall  ToolCallEx  `json:"tool_call"`
}

func LoadHints(hintsDir string) error
func GetHints(modelName string) ModelHints
func GetAllModels() []string
```

**Startup integration:**
```go
// In main.go
if err := LoadHints("cmd/mcp-server/hints"); err != nil {
    log.Printf("Warning: failed to load hints: %v (using defaults)", err)
}
```

---

### Phase 4: Dynamic Tool Registration

**Approach:** Modify tool registration in `main.go` to use model-specific descriptions

**Current:**
```go
var queryRadiationToolDef = mcp.NewTool("query_radiation",
    mcp.WithDescription("Find radiation measurements near..."),
    // ...
)
```

**New approach:**
```go
// tool_query_radiation.go
func getQueryRadiationTool(modelName string) mcp.Tool {
    hints := GetHints(modelName)
    desc := hints.Tools["query_radiation"].Description
    if desc == "" {
        desc = defaultQueryRadiationDesc // fallback
    }
    
    return mcp.NewTool("query_radiation",
        mcp.WithDescription(desc),
        // ... parameters unchanged
    )
}
```

**Challenge:** The MCP SDK's `AddTool()` is called once at startup, not per-request.

**Solution:** Use a **middleware approach** instead:

1. Keep base tool definitions static
2. Inject hints into the **response metadata** or **system prompt** dynamically
3. Use MCP's `ToolResult` content to include model-specific guidance

**Revised approach:**
- Tools remain static (MCP protocol limitation)
- Hints injected via:
  - `_ai_hint` fields in tool responses (already present in some tools)
  - Response metadata in `CallToolResult`
  - System prompt augmentation (if MCP supports it)

---

### Phase 5: HTTP Middleware for Model Detection

**File:** `cmd/mcp-server/model_middleware.go` (new)

**Wrap SSE and HTTP handlers to detect model:**

```go
func ModelDetectionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        model := DetectModel(r)
        
        // Store in context for downstream handlers
        ctx := context.WithValue(r.Context(), "model", model)
        r = r.WithContext(ctx)
        
        // Log for analytics
        log.Printf("MCP request from model=%s (%s)", model.Name, model.DisplayName)
        
        next.ServeHTTP(w, r)
    })
}
```

**Integration in main.go:**
```go
mux := http.NewServeMux()
mux.Handle("/mcp-http", ModelDetectionMiddleware(httpServer))
mux.Handle("/mcp/", ModelDetectionMiddleware(sseServer))
```

---

### Phase 6: Update Tool Handlers to Use Hints

**Modify existing tool handlers** to inject model-specific hints into responses.

**Example for `query_radiation`:**

```go
func handleQueryRadiation(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    model := GetModelFromContext(ctx)
    hints := GetHints(model.Name)
    
    // ... existing query logic ...
    
    result := map[string]any{
        "count": len(measurements),
        "measurements": measurements,
        "_ai_hint": hints.Tools["query_radiation"].OutputHints,
        "_formatting_rules": hints.FormattingRules,
        "_ai_generated_note": "...",
    }
    
    return jsonResult(result)
}
```

---

### Phase 7: Testing & Validation

**Test scenarios:**

1. **Model detection:**
   - Send requests with different User-Agent headers
   - Verify correct model profile is detected

2. **Hints injection:**
   - Call tools with different model contexts
   - Verify response includes model-specific hints

3. **Fallback behavior:**
   - Send request with unknown User-Agent
   - Verify `default` hints are used

4. **Performance:**
   - Measure latency impact of hints loading (should be negligible)
   - Verify hints are cached (no per-request file I/O)

**Test commands:**
```bash
# Test Claude detection
curl -H "User-Agent: Claude/3.5" https://simplemap.safecast.org/mcp-http -d '{...}'

# Test Qwen detection
curl -H "User-Agent: Qwen/2.5" https://simplemap.safecast.org/mcp-http -d '{...}'

# Test unknown model (should use default)
curl -H "User-Agent: UnknownBot/1.0" https://simplemap.safecast.org/mcp-http -d '{...}'
```

---

## File Changes Summary

| File | Action | Purpose |
|------|--------|---------|
| `cmd/mcp-server/model_adapter.go` | Create | Model detection logic |
| `cmd/mcp-server/hints_loader.go` | Create | JSON hints loading & caching |
| `cmd/mcp-server/model_middleware.go` | Create | HTTP middleware for model detection |
| `cmd/mcp-server/hints/claude.json` | Create | Claude-specific hints |
| `cmd/mcp-server/hints/qwen.json` | Create | Qwen-specific hints |
| `cmd/mcp-server/hints/kimi.json` | Create | Kimi-specific hints |
| `cmd/mcp-server/hints/gpt.json` | Create | GPT-specific hints |
| `cmd/mcp-server/hints/default.json` | Create | Default fallback hints |
| `cmd/mcp-server/main.go` | Modify | Integrate model detection, load hints |
| `cmd/mcp-server/tool_*.go` | Modify | Inject hints into responses (18 files) |

---

## Estimated Scope

| Component | Lines of Code | Complexity |
|-----------|--------------|------------|
| Model adapter (detection) | ~150 | Low |
| Hints loader | ~100 | Low |
| HTTP middleware | ~50 | Low |
| Hints JSON files (5) | ~500 | Medium |
| Tool handler updates (18) | ~360 | Low |
| **Total** | **~1160** | **Medium** |

---

## Timeline

| Phase | Estimated Time | Dependencies |
|-------|---------------|--------------|
| 1. Model detection | 1-2 hours | None |
| 2. Hints JSON design | 2-3 hours | Phase 1 |
| 3. Hints loader | 1-2 hours | Phase 2 |
| 4. HTTP middleware | 1 hour | Phase 1 |
| 5. Tool updates | 3-4 hours | Phase 3 |
| 6. Testing | 2 hours | All phases |
| **Total** | **10-14 hours** | |

---

## Future Enhancements

1. **Hot reload:** Watch hints JSON files for changes, reload without restart
2. **Admin API:** `/api/hints/{model}` endpoint to view/edit hints dynamically
3. **Analytics:** Track which models are using which tools most frequently
4. **A/B testing:** Serve different hint variants to optimize AI behavior
5. **Model capability detection:** Auto-detect features like context window size, JSON mode support

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| MCP protocol doesn't support dynamic tool descriptions | Medium | Use response metadata instead of tool definitions |
| Hints files become large/unmaintainable | Low | Keep hints focused, use shared formatting rules |
| Model detection fails for new clients | Low | Graceful fallback to `default` profile |
| Performance impact from JSON parsing | Low | Cache hints at startup, no per-request file I/O |

---

## Next Steps

1. **Create Phase 1 files** (model_adapter.go, hints_loader.go)
2. **Design initial hints JSON** for Claude and Qwen (priority models)
3. **Integrate middleware** into main.go
4. **Test with real AI clients** (Claude Desktop, Qwen CLI, etc.)
5. **Iterate on hints** based on actual AI behavior
