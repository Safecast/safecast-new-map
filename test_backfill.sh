#!/bin/bash
# Test the backfill mode with 2015-11-01 data (should fetch ~250 files)

echo "Testing Safecast backfill mode..."
echo "This will import one batch (~250 files from Nov 2015)"
echo ""

# Kill the running instance temporarily
pkill -f "safecast-new-map.*safecast-fetcher"
sleep 2

# Run a single-cycle test
timeout 2m ./safecast-new-map \
  -db-type pgx \
  -db-conn 'postgres://postgres:@127.0.0.1:5432/safecast?sslmode=prefer' \
  -safecast-fetcher \
  -safecast-fetcher-start-date "2015-11-01" \
  -safecast-fetcher-interval "10m" \
  -safecast-fetcher-batch-size 0 \
  -safecast-fetcher-backfill \
  -admin-password test123

echo ""
echo "Test complete. Check logs above to see import results."
echo "If successful, run ./backfill_complete.sh for full import."
