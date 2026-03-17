#!/bin/bash
#
# refresh_track_stats.sh
# Refreshes the track_statistics materialized view
#
# This should be run periodically (e.g., via cron) to keep track statistics up-to-date
# The refresh is concurrent, so it doesn't block queries to the existing view
#
# Usage:
#   export POSTGRES_URL='host=localhost port=5432 dbname=safecast user=safecast password=yourpassword sslmode=disable'
#   ./refresh_track_stats.sh
#
# Or with inline credentials:
#   PGPASSWORD=yourpassword ./refresh_track_stats.sh
#
# Add to crontab for automatic refresh (every hour):
#   0 * * * * cd /path/to/safecast-new-map/tools && PGPASSWORD=yourpass ./refresh_track_stats.sh >> /var/log/track_stats_refresh.log 2>&1

set -e

# Parse POSTGRES_URL if set, otherwise use defaults
if [ -n "$POSTGRES_URL" ]; then
    DB_HOST=$(echo "$POSTGRES_URL" | grep -oP 'host=\K[^ ]+' || echo "localhost")
    DB_PORT=$(echo "$POSTGRES_URL" | grep -oP 'port=\K[^ ]+' || echo "5432")
    DB_NAME=$(echo "$POSTGRES_URL" | grep -oP 'dbname=\K[^ ]+' || echo "safecast")
    DB_USER=$(echo "$POSTGRES_URL" | grep -oP 'user=\K[^ ]+' || echo "safecast")
    DB_PASS=$(echo "$POSTGRES_URL" | grep -oP 'password=\K[^ ]+' || echo "")
else
    DB_HOST="${DB_HOST:-localhost}"
    DB_PORT="${DB_PORT:-5432}"
    DB_NAME="${DB_NAME:-safecast}"
    DB_USER="${DB_USER:-safecast}"
    DB_PASS="${PGPASSWORD:-}"
fi

echo "[$(date '+%Y-%m-%d %H:%M:%S')] Starting track statistics refresh..."

# Use CONCURRENTLY to avoid locking the view during refresh
# This allows queries to continue using the old data while the refresh is in progress
START_TIME=$(date +%s)

PGPASSWORD="$DB_PASS" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
    -c "REFRESH MATERIALIZED VIEW CONCURRENTLY track_statistics;" 2>&1

if [ $? -eq 0 ]; then
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))

    # Get summary statistics
    STATS=$(PGPASSWORD="$DB_PASS" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t \
        -c "SELECT COUNT(*) || ' tracks, ' || SUM(marker_count) || ' markers, ' || SUM(spectra_count) || ' spectra' FROM track_statistics;")

    echo "[$(date '+%Y-%m-%d %H:%M:%S')] ✅ Refresh completed in ${DURATION}s: $STATS"
else
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] ❌ Refresh failed"
    exit 1
fi
