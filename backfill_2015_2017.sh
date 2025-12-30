#!/bin/bash
# Backfill Safecast bGeigie imports from November 2015 to November 2017
# 
# This script runs the Safecast fetcher in backfill mode, which will:
# - Ignore the current database state (lastID tracking)
# - Import ALL approved bGeigie imports from the specified date range
# - Process files in chronological order (oldest first)
#
# The 10-page limit means each run will fetch ~250 files max (25 per page).
# With 191 pages total for your date range (~4,775 files), you'll need to:
# - Run this multiple times
# - Or increase the batch size
# - Or remove the 10-page safety limit temporarily

# Configuration
START_DATE="2015-11-01"
BATCH_SIZE=0  # 0 = unlimited (will still stop at 10 pages due to safety limit)
INTERVAL="1h" # Set high interval so it only runs once per hour

echo "========================================="
echo "Safecast Backfill: $START_DATE onwards"
echo "========================================="
echo ""
echo "This will import historical bGeigie data in BACKFILL MODE."
echo "The fetcher will process ~250 files per poll cycle (10 pages × 25 files)."
echo ""
echo "Press Ctrl+C to stop when you've imported enough data."
echo ""

./safecast-new-map \
  -db-type pgx \
  -db-conn 'postgres://postgres:@127.0.0.1:5432/safecast?sslmode=prefer' \
  -safecast-fetcher \
  -safecast-fetcher-start-date "$START_DATE" \
  -safecast-fetcher-interval "$INTERVAL" \
  -safecast-fetcher-batch-size $BATCH_SIZE \
  -safecast-fetcher-backfill

echo ""
echo "Backfill stopped."
echo "Check the logs above to see how many imports were processed."
