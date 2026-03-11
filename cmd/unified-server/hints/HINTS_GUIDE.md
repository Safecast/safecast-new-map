# AI Hints Development Guide

## Overview

AI hints are JSON configuration files that provide model-specific guidance to AI assistants using the Safecast MCP server.

## File Structure

```
cmd/mcp-server/hints/
├── claude.json    - Claude (Anthropic)
├── gpt.json       - GPT-4 / ChatGPT (OpenAI)
├── qwen.json      - Qwen (Alibaba)
├── kimi.json      - Kimi (Moonshot AI)
└── default.json   - Fallback for unknown models
```

## JSON Schema

```json
{
  "model": "model_name",
  "display_name": "Human-readable name",
  "capabilities": ["tool_use", "json_mode"],
  "system_prompt": "Base system prompt",
  "tools": {
    "tool_name": {
      "description": "What + critical rules",
      "examples": [{"user_query": "...", "tool_call": {...}}],
      "output_hints": "How to format results",
      "warnings": ["Mistakes to avoid"],
      "formatting_rules": {"rule": "value"}
    }
  },
  "global_formatting_rules": {...}
}
```

## Best Practices

### System Prompt
- Define role: "You are..."
- Set context: what data available
- Behavior rules: "Always...", "Never..."
- Output requirements

### Tool Description
- What it does (1-2 sentences)
- Data source (mobile/fixed, historical/real-time)
- Critical workflow rules ("MUST call X after")
- Limitations

### Examples
- Realistic user queries
- Proper argument values
- Follow-up calls for multi-step workflows

### Warnings
- Safety-critical rules
- Common misconceptions
- Data limitations

### Output Hints
- Table structure
- Link formatting
- Context requirements
- Unit specifications

## Model-Specific Tips

| Model | Strengths | Hint Focus |
|-------|-----------|------------|
| Claude | Long context, reasoning | Detailed workflows, safety |
| GPT-4 | Function calling, JSON | Structured outputs |
| Qwen | Multilingual | Bilingual hints |
| Kimi | Long documents | Concise outputs |

## Testing

```bash
# Validate JSON
python -m json.tool cmd/mcp-server/hints/your-model.json

# Test loading
export MCP_HINTS_DIR=./cmd/mcp-server/hints
./bin/mcp-server 2>&1 | grep "Loaded hints"

# Run tests
go test ./cmd/mcp-server/model-adapter/... -v
```

## Common Mistakes

1. Vague descriptions
2. Missing follow-up calls
3. No warnings
4. Inconsistent formatting
5. Outdated examples
