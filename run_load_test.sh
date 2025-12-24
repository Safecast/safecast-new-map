#!/bin/bash

# Phase 2 Load Test Runner - Benchmarks gzip, caching, and DB pool efficiency
# Usage: ./run_load_test.sh [concurrency] [duration]

CONCURRENCY=${1:-50}
DURATION=${2:-30}
BINARY="./load_test_bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Phase 2 Performance Optimization - Load Test ===${NC}"
echo ""

# Check if server is running
if ! curl -s http://localhost:8765/api/geoip > /dev/null 2>&1; then
    echo -e "${RED}✗ Server not running on http://localhost:8765${NC}"
    echo "Start the server first with: ./safecast-new-map -safecast-realtime -admin-password test123"
    exit 1
fi
echo -e "${GREEN}✓ Server is running${NC}"

# Check if load test binary exists, build if needed
if [ ! -f "$BINARY" ]; then
    echo "Building load test binary..."
    go build -o "$BINARY" ./cmd/load_test/main.go
    if [ $? -ne 0 ]; then
        echo -e "${RED}✗ Failed to build load test${NC}"
        exit 1
    fi
fi

# Check database
DB_COUNT=$(psql -U postgres -h 127.0.0.1 -d safecast -c "SELECT COUNT(*) FROM markers" 2>/dev/null | grep -v "count" | tail -1)
if [ -n "$DB_COUNT" ]; then
    echo -e "${GREEN}✓ Database connected ($(printf "%'d\n" $DB_COUNT) markers)${NC}"
else
    echo -e "${RED}✗ Database not reachable${NC}"
    exit 1
fi

echo ""
echo "Test Configuration:"
echo "  Concurrency: $CONCURRENCY workers"
echo "  Duration: $DURATION seconds"
echo ""

# Run the load test
echo -e "${YELLOW}Starting load test...${NC}"
echo ""
./$BINARY -concurrency $CONCURRENCY -duration $DURATION

echo ""
echo -e "${GREEN}Load test complete!${NC}"
echo ""
echo "Next Steps:"
echo "  • Check server logs: tail -f safecast.log | grep -E 'cache|pool|gzip'"
echo "  • Inspect pool usage: psql -U postgres -h 127.0.0.1 -d safecast -c \"SELECT count(*) FROM pg_stat_activity\""
echo "  • Run again with different concurrency: ./run_load_test.sh 100 60"
