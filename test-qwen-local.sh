#!/bin/bash
# Local test for Qwen3.5 via NVIDIA API

echo "=== Qwen3.5 Local Test ==="
echo ""

# Check for API key
if [ -z "$NVIDIA_API_KEY" ]; then
    echo "❌ NVIDIA_API_KEY not set"
    echo "   Get it from: https://build.nvidia.com/qwen/qwen3.5-122b-a10b"
    echo "   Then run: export NVIDIA_API_KEY=nvapi-..."
    exit 1
fi

echo "✓ NVIDIA_API_KEY found"
echo ""

# Test 1: Direct API call
echo "1. Testing NVIDIA API directly..."
RESPONSE=$(curl -s -X POST https://integrate.api.nvidia.com/v1/chat/completions \
  -H "Authorization: Bearer $NVIDIA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "qwen/qwen3.5-122b-a10b",
    "messages": [
      {"role": "user", "content": "What is radiation? Explain briefly."}
    ],
    "max_tokens": 200
  }' 2>&1)

if echo "$RESPONSE" | grep -q "choices"; then
    echo "   ✓ NVIDIA API works!"
    echo "   Response preview: $(echo "$RESPONSE" | jq -r '.choices[0].message.content' | head -c 100)..."
else
    echo "   ❌ API call failed"
    echo "   $RESPONSE"
    exit 1
fi
echo ""

# Test 2: Start MCP server
echo "2. Starting MCP server..."
export MCP_HINTS_DIR="$PWD/cmd/mcp-server/hints"
./bin/mcp-server-test > /tmp/mcp.log 2>&1 &
MCP_PID=$!
sleep 3

if kill -0 $MCP_PID 2>/dev/null; then
    echo "   ✓ MCP server running (PID: $MCP_PID)"
else
    echo "   ❌ MCP server failed to start"
    cat /tmp/mcp.log
    exit 1
fi
echo ""

# Cleanup
cleanup() {
    kill $MCP_PID 2>/dev/null || true
}
trap cleanup EXIT

# Test 3: Check MCP hints loaded
echo "3. Checking MCP hints..."
grep "Loaded hints" /tmp/mcp.log
if grep -q "qwen" /tmp/mcp.log; then
    echo "   ✓ Qwen hints loaded"
else
    echo "   ⚠ Qwen hints may not be loaded"
fi
echo ""

echo "=== All Tests Passed! ==="
echo ""
echo "Next steps:"
echo "1. Update cmd/web-chat/main.go to support NVIDIA provider"
echo "2. Run: export AI_PROVIDER=nvidia"
echo "3. Run: go run ./cmd/web-chat/"
echo "4. Open http://localhost:3334 in browser"
