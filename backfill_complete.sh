#!/bin/bash
# One-time backfill script for 2015-2017 Safecast data
# This processes ALL 191 pages (~4,775 files) in a single run
#
# WARNING: This temporarily removes the 10-page safety limit
# You'll need to rebuild after modifying fetcher.go

echo "========================================="
echo "COMPLETE Safecast Backfill: 2015-2017"
echo "========================================="
echo ""
echo "This script will:"
echo "  1. Temporarily modify the fetcher to allow 200 pages"
echo "  2. Run a backfill import (5,000 files per batch)"  
echo "  3. Restore the 10-page safety limit"
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 1
fi

# Backup the original fetcher
cp pkg/safecast-fetcher/fetcher.go pkg/safecast-fetcher/fetcher.go.backup

# Increase the page limit to 200 for batched processing
sed -i 's/if page > 10 {/if page > 200 {/' pkg/safecast-fetcher/fetcher.go
sed -i 's/stopped at page 10/stopped at page 200/' pkg/safecast-fetcher/fetcher.go

# Rebuild
echo "Rebuilding with increased page limit..."
go build -o safecast-new-map safecast-new-map.go

# Run the backfill (will process 200 pages = ~5,000 files)
echo ""
echo "Starting batch backfill import (200 pages)..."
echo "This will take several minutes to process ~5,000 files."
echo ""

timeout 15m ./safecast-new-map \
  -db-type pgx \
  -db-conn 'postgres://postgres:@127.0.0.1:5432/safecast?sslmode=prefer' \
  -safecast-fetcher \
  -safecast-fetcher-start-date "2015-11-01" \
  -safecast-fetcher-interval "1h" \
  -safecast-fetcher-batch-size 0 \
  -safecast-fetcher-backfill &

FETCHER_PID=$!

# Wait for one complete cycle (give it up to 30 minutes)
echo "Waiting for import to complete (timeout: 30 minutes)..."
wait $FETCHER_PID

# Restore original fetcher
echo ""
echo "Restoring original fetcher.go..."
mv pkg/safecast-fetcher/fetcher.go.backup pkg/safecast-fetcher/fetcher.go

# Rebuild with original settings
echo "Rebuilding with original 10-page limit..."
go build -o safecast-new-map safecast-new-map.go

echo ""
echo "========================================="
echo "Backfill complete!"
echo "========================================="
echo ""
echo "Check your database to verify the imports:"
echo "  SELECT COUNT(*) FROM uploads WHERE source = 'safecast-api';"
echo ""
