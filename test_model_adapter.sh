#!/bin/bash
# Test script for MCP model adapter and hints loading

set -e

echo "=== Testing MCP Model Adapter ==="
echo ""

# Build the MCP server
echo "1. Building MCP server..."
cd /home/rob/Documents/Safecast/safecast-new-map
go build -o bin/mcp-server-test ./cmd/mcp-server
echo "   ✓ Build successful"
echo ""

# Test 1: Check hints directory
echo "2. Checking hints directory..."
if [ -d "cmd/mcp-server/hints" ]; then
    echo "   ✓ Hints directory exists"
    ls -1 cmd/mcp-server/hints/*.json 2>/dev/null | while read f; do
        echo "     - $(basename $f)"
    done
else
    echo "   ✗ Hints directory not found!"
    exit 1
fi
echo ""

# Test 2: Run unit tests
echo "3. Running unit tests..."
go test ./cmd/mcp-server/model-adapter/... -v 2>&1 | grep -E "^(=== RUN|--- PASS|--- FAIL|PASS|FAIL|ok)" || true
echo ""

echo "=== All Tests Complete ==="
