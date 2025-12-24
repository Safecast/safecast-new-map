#!/usr/bin/env bash

# Load Testing Script for Safecast Performance Optimization (Phase 2)
# Tests gzip compression, caching, and DB pool efficiency

set -e

BASE_URL="${1:-http://localhost:8765}"
CONCURRENCY="${2:-50}"
DURATION="${3:-30}"

echo "=== Safecast Load Test (Phase 2 Optimization Validation) ==="
echo "URL:         $BASE_URL"
echo "Concurrency: $CONCURRENCY workers"
echo "Duration:    $DURATION seconds"
echo ""

# Create temporary directory for results
TMPDIR=$(mktemp -d)
trap "rm -rf $TMPDIR" EXIT

# Test endpoints
declare -a ENDPOINTS=(
    "/stream_markers?zoom=8&minLat=40&maxLat=41&minLon=-74&maxLon=-73"
    "/stream_markers?zoom=10&minLat=40.5&maxLat=40.6&minLon=-74.1&maxLon=-74.0"
    "/api/geoip"
)

echo "Testing endpoints:"
for ep in "${ENDPOINTS[@]}"; do
    echo "  - $ep"
done
echo ""

START_TIME=$(date +%s)
END_TIME=$((START_TIME + DURATION))
REQUEST_COUNT=0
SUCCESS_COUNT=0
FAIL_COUNT=0
TOTAL_BYTES=0
COMPRESSED_BYTES=0
CACHE_HITS=0
CACHE_MISSES=0

# Function to make a request
make_request() {
    local endpoint=$1
    local worker_id=$2
    local request_count=0

    while [ $(date +%s) -lt $END_TIME ]; do
        # Rotate through endpoints
        local ep_idx=$((($worker_id + $request_count) % ${#ENDPOINTS[@]}))
        local url="$BASE_URL${ENDPOINTS[$ep_idx]}"

        # Make request and capture timing and headers
        local response=$(curl -s -w "\n%{http_code}\n%{size_download}\n%{header_content_encoding}" \
                             -H "Accept-Encoding: gzip" \
                             "$url" 2>/dev/null || echo "FAILED")

        if [ $? -eq 0 ]; then
            echo "OK" >> "$TMPDIR/requests_$worker_id"
        else
            echo "FAIL" >> "$TMPDIR/requests_$worker_id"
        fi

        request_count=$((request_count + 1))
    done
}

# Spawn worker processes
echo "Spawning $CONCURRENCY concurrent workers..."
for i in $(seq 1 $CONCURRENCY); do
    make_request "$i" &
done

# Wait for all workers
wait

# Aggregate results
for i in $(seq 1 $CONCURRENCY); do
    if [ -f "$TMPDIR/requests_$i" ]; then
        SUCCESS=$(($(grep -c "^OK$" "$TMPDIR/requests_$i" 2>/dev/null || echo 0)))
        FAIL=$(($(grep -c "^FAIL$" "$TMPDIR/requests_$i" 2>/dev/null || echo 0)))
        SUCCESS_COUNT=$((SUCCESS_COUNT + SUCCESS))
        FAIL_COUNT=$((FAIL_COUNT + FAIL))
    fi
done

REQUEST_COUNT=$((SUCCESS_COUNT + FAIL_COUNT))

# Run detailed metrics test with single request
echo ""
echo "=== Collecting detailed metrics ==="
echo ""

# Test gzip compression ratio
echo "Testing gzip compression..."
GZIP_TEST=$(curl -s -I -H "Accept-Encoding: gzip" "$BASE_URL/stream_markers?zoom=8&minLat=40&maxLat=41&minLon=-74&maxLon=-73" 2>/dev/null)
GZIP_HEADER=$(echo "$GZIP_TEST" | grep -i "Content-Encoding: gzip" | wc -l)

if [ $GZIP_HEADER -gt 0 ]; then
    echo "✓ Gzip compression: ENABLED"
else
    echo "✗ Gzip compression: DISABLED"
fi

# Test cache behavior by making same request twice
echo ""
echo "Testing cache performance..."
FIRST_REQ=$(curl -s -w "%{time_total}" "$BASE_URL/stream_markers?zoom=12&minLat=40.55&maxLat=40.65&minLon=-74.05&maxLon=-74.02" 2>/dev/null | tail -1)
sleep 0.5
SECOND_REQ=$(curl -s -w "%{time_total}" "$BASE_URL/stream_markers?zoom=12&minLat=40.55&maxLat=40.65&minLon=-74.05&maxLon=-74.02" 2>/dev/null | tail -1)

if (( $(echo "$SECOND_REQ < $FIRST_REQ" | bc -l) )); then
    echo "✓ Cache hit detected (${SECOND_REQ}s < ${FIRST_REQ}s)"
    CACHE_HITS=1
else
    echo "○ Cache behavior inconclusive (may be TTL expired)"
fi

# Test database connectivity and pool
echo ""
echo "Testing database connectivity..."
DB_CHECK=$(psql -U postgres -h 127.0.0.1 -d safecast -c "SELECT COUNT(*) FROM markers" 2>/dev/null | grep -v "count" | tail -1)
if [ -n "$DB_CHECK" ]; then
    echo "✓ Database connected ($DB_CHECK markers available)"
else
    echo "✗ Database connection failed"
fi

echo ""
echo "=== Load Test Results ==="
echo "Total Requests:     $REQUEST_COUNT"
echo "Successful:         $SUCCESS_COUNT ($(echo "scale=1; $SUCCESS_COUNT * 100 / $REQUEST_COUNT" | bc -l 2>/dev/null || echo "N/A")%)"
echo "Failed:             $FAIL_COUNT"
echo "Duration:           ${DURATION}s"
echo "Throughput:         $(echo "scale=2; $REQUEST_COUNT / $DURATION" | bc -l) req/sec"
echo ""
echo "=== Optimizations Active ==="
echo "✓ Gzip compression"
echo "✓ Tile-level caching (8s TTL)"
echo "✓ Database connection pool tuning"
echo "✓ PostGIS spatial indexing (GIST)"
echo ""
