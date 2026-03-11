# NVIDIA Qwen3.5-122B-A10B Integration Guide

## Quick Start

### 1. Get NVIDIA API Key
1. Visit: https://build.nvidia.com/qwen/qwen3.5-122b-a10b
2. Click "Get API Key"
3. Copy your API key (starts with `nvapi-...`)

### 2. Configure (.env or systemd)
```bash
AI_PROVIDER=nvidia
NVIDIA_API_KEY=nvapi-your-key-here
QWEN_MODEL=qwen/qwen3.5-122b-a10b
MCP_URL=http://localhost:3333/mcp-http
PORT=3334
```

### 3. Restart
```bash
systemctl restart safecast-web-chat
```

## Model Comparison

| Feature | Claude Haiku | Qwen3.5-122B |
|---------|-------------|--------------|
| Reasoning | Good | 🔥 Excellent |
| Speed | ⚡ Very fast | Good |
| Multilingual | Good | 🔥 Excellent |
| Tool Calling | ✅ Very reliable | ✅ Good |
| Cost | $0.80-4/1M tokens | ~$0.80-2.40/1M tokens |

## Code Changes Needed

### 1. main.go - Add provider selection
```go
aiProvider := os.Getenv("AI_PROVIDER")  // "anthropic" or "nvidia"
if aiProvider == "nvidia" {
    apiKey = os.Getenv("NVIDIA_API_KEY")
    model = "qwen/qwen3.5-122b-a10b"
} else {
    apiKey = os.Getenv("ANTHROPIC_API_KEY")
    model = "claude-sonnet-4-5"
}
```

### 2. Add NVIDIA API client function
- Uses OpenAI-compatible API format
- Endpoint: `https://integrate.api.nvidia.com/v1/chat/completions`
- Auth: `Bearer {NVIDIA_API_KEY}`

## Testing

```bash
# Local test
export AI_PROVIDER=nvidia
export NVIDIA_API_KEY=nvapi-xxx
go run ./cmd/web-chat/

# Check logs
journalctl -u safecast-web-chat -f
# Should show: "Using NVIDIA Qwen API: model=qwen/qwen3.5-122b-a10b"
```

## Rollback

```bash
export AI_PROVIDER=anthropic
systemctl restart safecast-web-chat
```

## MCP Server Compatibility

✅ Your MCP server already supports Qwen via `qwen.json` hints file.
No MCP server changes needed!
