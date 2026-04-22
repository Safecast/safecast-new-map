# MCP Model Hints

This directory contains model-specific hints and prompts for the Safecast MCP server.

## Overview

The hints system allows the MCP server to provide tailored descriptions, examples, and formatting guidance to different AI models (Claude, Qwen, Kimi, GPT, etc.) that connect to the Safecast database.

## Source of Truth

> **Runtime source of truth is the `ai_hints` table in PostgreSQL**, not the JSON files in this directory.
>
> On startup the server:
> 1. Loads JSON files from this directory into memory (bootstrap/fallback).
> 2. Seeds any missing rows into `ai_hints` with `ON CONFLICT (model) DO NOTHING` — existing DB rows are never overwritten by the JSON files.
> 3. Reloads all non-deleted rows from `ai_hints` into memory and hands them to the MCP server.
>
> Admins edit hints through the **`/admin/ai-hints`** page (structured forms, history, soft delete, import/export). Changes in the DB take effect after clicking **Reload into Memory** (or restarting the service). The JSON files here are treated as seed templates and are **not** re-read at runtime unless the DB row is missing.
>
> To sync DB edits back into the repo, use **Export JSON** in the admin UI and commit the file.

## Files

| File | Model | Description |
|------|-------|-------------|
| `claude.json` | Claude 3.5 Sonnet | Detailed, verbose hints with comprehensive examples |
| `qwen.json` | Qwen 2.5 | Bilingual hints (Chinese/English) with structured formatting |
| `kimi.json` | Kimi | Medium-length hints with clear guidance |
| `gpt.json` | GPT-4 / ChatGPT | Concise, token-optimized hints |
| `default.json` | Unknown/Fallback | Generic hints for unknown models |

## JSON Structure

Each hints file follows this schema:

```json
{
  "model": "claude",
  "display_name": "Claude 3.5 Sonnet",
  "capabilities": ["long_context", "tool_use", "json_mode"],
  "system_prompt": "System prompt for this model...",
  "tools": {
    "tool_name": {
      "description": "Tool description override",
      "description_en": "English description (for bilingual hints)",
      "examples": [
        {
          "user_query": "Example user query",
          "tool_call": {
            "name": "tool_name",
            "arguments": {"arg1": "value1"}
          },
          "follow_up": {
            "name": "follow_up_tool",
            "arguments": {"arg": "value"}
          }
        }
      ],
      "output_hints": "Guidance for formatting output",
      "formatting_rules": {
        "rule_name": "rule value"
      },
      "warnings": ["Known gotchas"],
      "topics": ["Available topics for info tools"]
    }
  },
  "global_formatting_rules": {
    "tone": "Scientific, objective",
    "coordinates": "Map link format",
    "tables": "Markdown table usage"
  }
}
```

## Configuration

Set the hints directory via environment variable:

```bash
export MCP_HINTS_DIR=/path/to/hints
```

Default: `./hints/` relative to the binary location.

## Adding a New Model

The recommended path is the admin UI:

1. Open `/admin/ai-hints?password=...` (or log in as an admin user).
2. Click **+ Add bot**, fill in the display name (the `model` slug is derived automatically), then save.
3. Fill in capabilities, system prompt, global rules, and per-tool hints via the structured editor.
4. Click **Reload into Memory** so the MCP server picks up the change without a restart.
5. (Optional) Update `model_detection.go` so the server can auto-detect the new model from the `User-Agent` header.

If you prefer to seed a new bot via JSON (e.g. for a fresh environment):

1. Create a new JSON file: `hints/newmodel.json`.
2. Set `"model": "newmodel"` in the JSON — this is the slug that will be inserted into `ai_hints`.
3. Restart the server — the row will be seeded into the DB on the next startup.
4. Add detection logic to `model_detection.go` and (optionally) a constant to `hints.go`.

Example detection addition:

```go
case strings.Contains(ua, "newmodel"):
    return ModelNewModel
```

## Customizing Hints

### Tool Description

Override the default tool description for model-specific behavior:

```json
"query_radiation": {
  "description": "Your custom description here..."
}
```

### Examples

Provide usage examples that resonate with the model's training:

```json
"examples": [
  {
    "user_query": "What's the radiation level?",
    "tool_call": {
      "name": "query_radiation",
      "arguments": {"lat": 35.6762, "lon": 139.6503}
    }
  }
]
```

### Formatting Rules

Specify output formatting preferences:

```json
"formatting_rules": {
  "table_format": "Always use markdown tables",
  "map_links": "Required for all coordinates"
}
```

### Warnings

Document model-specific gotchas:

```json
"warnings": [
  "Never report 'no real-time data' without calling sensor_current first"
]
```

## Testing

Test model detection with curl:

```bash
# Test Claude detection
curl -H "User-Agent: Claude/3.5" \
     -H "X-AI-Model: claude" \
     https://simplemap.safecast.org/mcp-http \
     -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"ping"}}'

# Test Qwen detection
curl -H "User-Agent: Qwen/2.5" \
     https://simplemap.safecast.org/mcp-http \
     -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"ping"}}'

# Test unknown model (uses default)
curl -H "User-Agent: UnknownBot/1.0" \
     https://simplemap.safecast.org/mcp-http \
     -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"ping"}}'
```

## Runtime Behavior

1. **Startup (file bootstrap)**: JSON files in this directory are loaded and cached in memory.
2. **Startup (DB seed)**: Missing rows are inserted into `ai_hints` with `ON CONFLICT DO NOTHING`; existing DB rows are left alone.
3. **Startup (DB load)**: All rows where `deleted_at IS NULL` are read from `ai_hints` and swapped into memory — this becomes the authoritative in-memory state.
4. **Admin edit**: Writes go to `ai_hints` (with a history snapshot to `ai_hints_history`). The in-memory copy is refreshed on **Reload into Memory**.
5. **Request**: Model is detected from `User-Agent` or `X-AI-Model` header.
6. **Response**: Tool results include model-specific formatting hints as metadata.

## Performance

- Hints are loaded once at startup and on admin reload (no per-request file or DB I/O).
- Thread-safe access via read-write mutex.
- Atomic in-memory swap on reload (`SetHints`) — no half-loaded state.

## Troubleshooting

### Hints not loading

Check logs for:
```
Loaded hints for models: [claude qwen kimi gpt unknown]
```

If you see:
```
Warning: failed to load hints: ... (using default hints)
```

Verify:
1. Hints directory exists
2. JSON files are valid
3. `MCP_HINTS_DIR` is set correctly

### Model not detected

1. Check `User-Agent` header contains model name
2. Or set `X-AI-Model: <model-name>` header explicitly
3. Verify detection logic in `model_detection.go`

## License

Hints files are licensed under the same CC0 license as Safecast data.
