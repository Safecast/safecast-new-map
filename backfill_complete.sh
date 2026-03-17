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
go build -o safecast-new-map main.go

# ... (lines 34-60) ...

# Rebuild with original settings
echo "Rebuilding with original 10-page limit..."
go build -o safecast-new-map main.go

echo ""
echo "========================================="
echo "Backfill complete!"
echo "========================================="
echo ""
echo "Check your database to verify the imports:"
echo "  SELECT COUNT(*) FROM uploads WHERE source = 'safecast-api';"
echo ""
