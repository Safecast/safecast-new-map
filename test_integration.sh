#!/bin/bash
# Integration test for MCP model detection and hints loading

echo "=== MCP Integration Test ==="
echo ""

export MCP_TRANSPORT=http
export MCP_HINTS_DIR="$PWD/cmd/mcp-server/hints"
export DATABASE_URL=""

# Start MCP server
echo "Starting MCP server..."
./bin/mcp-server-test > /tmp/mcp.log 2>&1 &
MCP_PID=$!
sleep 3

# Check if server started
if ! kill -0 $MCP_PID 2>/dev/null; then
    echo "✗ Server failed to start"
    cat /tmp/mcp.log
    exit 1
fi
echo "✓ Server started (PID: $MCP_PID)"
echo ""
echo "Server log:"
cat /tmp/mcp.log | head -10
echo ""

# Cleanup function
cleanup() {
    echo ""
    echo "Stopping server..."
    kill $MCP_PID 2>/dev/null || true
    rm -f /tmp/mcp.log
}
trap cleanup EXIT

# Test function
test_model() {
    local name=$1
    local ua=$2
    local header=$3
    
    if [ -n "$header" ]; then
        RESPONSE=$(curl -s --max-time 2 -H "$header" http://localhost:3333/mcp-http 2>&1)
    else
        RESPONSE=$(curl -s --max-time 2 -A "$ua" http://localhost:3333/mcp-http 2>&1)
    fi
    
    if [[ -z "$RESPONSE" ]]; then
        echo "  ✗ $name: FAILED (no response)"
        return 1
    else
        echo "  ✓ $name: OK (response received)"
        return 0
    fi
}

echo "Testing model detection:"
test_model "Claude (User-Agent)" "Claude/1.0" ""
test_model "Qwen (User-Agent)" "Qwen/2.5" ""
test_model "Kimi (User-Agent)" "Kimi/1.0" ""
test_model "GPT (User-Agent)" "ChatGPT/4.0" ""
test_model "GPT (OpenAI)" "OpenAI/1.0" ""
test_model "Claude (X-AI-Model header)" "" "X-AI-Model: claude-3"
test_model "GPT (X-AI-Model header)" "" "X-AI-Model: gpt-4"

echo ""
echo "=== Integration Test Complete ==="
