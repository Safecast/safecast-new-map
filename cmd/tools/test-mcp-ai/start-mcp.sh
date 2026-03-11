#!/bin/bash
# Start MCP server with database connection

DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-safecast}"
DB_PASS="${DB_PASS:-}"

if [ -n "$DB_PASS" ]; then
    DATABASE_URL="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=prefer"
else
    DATABASE_URL="postgres://${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=prefer"
fi

export DATABASE_URL
export MCP_HINTS_DIR="$PWD/cmd/mcp-server/hints"

echo "🚀 Starting MCP server..."
echo "   Database: ${DB_HOST}:${DB_PORT}/${DB_NAME}"
echo "   Hints: $MCP_HINTS_DIR"
echo ""

# Kill existing
pkill -f "mcp-server" 2>/dev/null || true
sleep 1

./bin/mcp-server-test &
MCP_PID=$!

sleep 3

if kill -0 $MCP_PID 2>/dev/null; then
    echo "✅ MCP server running on :3333 (PID: $MCP_PID)"
    echo ""
    echo "Test with:"
    echo "  curl http://localhost:3333/api/sensors"
    echo "  ./bin/test-mcp-ai"
else
    echo "❌ Failed to start"
    exit 1
fi

wait $MCP_PID
