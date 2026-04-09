# MCP + AI Integration Test Tools

## Files

- `start-mcp.sh` - Start unified server with MCP endpoints
- `test-mcp-ai` - Main integration test (MCP + AI + Hints)
- `test-nvidia-models` - Test NVIDIA API models directly

## Quick Start

### 1. Start Unified Server

```bash
./start-mcp.sh
```

This starts the unified server on port 8765 with database connection.

### 2. Test AI Integration

```bash
export NVIDIA_API_KEY=nvapi-your-key
./test-mcp-ai
```

### 3. Test NVIDIA Models Only

```bash
export NVIDIA_API_KEY=nvapi-your-key
./test-nvidia-models
```

## Test Queries

When running `test-mcp-ai`, try these:

1. List gamma spectra
2. Current sensors in Japan
3. Radiation near Tokyo
4. Search area around Fukushima

## Requirements

- NVIDIA API key (get from https://build.nvidia.com)
- PostgreSQL database with Safecast data
- Unified server running on port 8765
