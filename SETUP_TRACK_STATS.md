# Track Statistics Setup Guide

This guide explains how to set up the track statistics materialized view for optimal admin page performance.

## Quick Start

For an existing database with data:

```bash
# 1. Create the materialized view
export POSTGRES_URL='host=localhost port=5432 dbname=safecast user=safecast password=safecast2025 sslmode=disable'
PGPASSWORD=safecast2025 psql -h localhost -U safecast -d safecast -f tools/create_track_stats_view.sql

# 2. Rebuild the application
go build -o safecast-new-map safecast-new-map.go

# 3. Restart the server
./safecast-new-map
```

## For New Database Setup

If setting up a fresh database:

```bash
# 1. Initialize the database (tables, indexes, etc.)
./safecast-new-map  # This creates tables automatically

# 2. Import your data
# (upload files, run migrations, etc.)

# 3. Create the track statistics view
PGPASSWORD=safecast2025 psql -h localhost -U safecast -d safecast -f tools/create_track_stats_view.sql
```

## Keeping Statistics Updated

The materialized view needs periodic refreshing to include new uploads.

### Option 1: Manual Refresh (Recommended for testing)

After uploading files or importing data:

```bash
./tools/refresh_track_stats.sh
```

### Option 2: Automatic Refresh via Cron (Recommended for production)

Edit your crontab:

```bash
crontab -e
```

Add this line to refresh every hour:

```
0 * * * * cd /home/rob/Documents/Safecast/safecast-new-map/tools && PGPASSWORD=safecast2025 ./refresh_track_stats.sh >> /var/log/track_stats_refresh.log 2>&1
```

Or refresh every 15 minutes for more up-to-date stats:

```
*/15 * * * * cd /home/rob/Documents/Safecast/safecast-new-map/tools && PGPASSWORD=safecast2025 ./refresh_track_stats.sh >> /var/log/track_stats_refresh.log 2>&1
```

### Option 3: Refresh After Data Import (Recommended for batch imports)

If you're importing large datasets, refresh after import completes:

```bash
# Import data
./import_my_data.sh

# Refresh stats
./tools/refresh_track_stats.sh
```

## Refresh Performance

The refresh time depends on database size:

| Markers | Tracks | Refresh Time (approx) |
|---------|--------|----------------------|
| 1M      | 100    | ~1-2 seconds         |
| 10M     | 1,000  | ~10-20 seconds       |
| 100M    | 10,000 | ~2-5 minutes         |

The refresh uses `CONCURRENTLY` which:
- ✅ Doesn't block queries (users can still access admin page)
- ✅ Doesn't lock the view
- ❌ Takes slightly longer than non-concurrent refresh

## Monitoring

Check the last refresh time:

```sql
SELECT
    schemaname,
    matviewname,
    last_refresh
FROM pg_matviews
WHERE matviewname = 'track_statistics';
```

Check view statistics:

```sql
SELECT
    COUNT(*) as total_tracks,
    SUM(marker_count) as total_markers,
    SUM(spectra_count) as total_spectra,
    MIN(first_date) as earliest_measurement,
    MAX(last_date) as latest_measurement
FROM track_statistics;
```

## Troubleshooting

### View doesn't exist error

If you get "relation track_statistics does not exist":

```bash
# Create it
PGPASSWORD=safecast2025 psql -h localhost -U safecast -d safecast -f tools/create_track_stats_view.sql
```

### Refresh takes too long

If refresh is taking too long:

1. Check if you have the required indexes on markers:
   ```sql
   SELECT indexname FROM pg_indexes WHERE tablename = 'markers' AND indexname LIKE '%trackid%';
   ```

2. You should see:
   - `idx_markers_trackid`
   - `idx_markers_trackid_date`
   - `idx_markers_trackid_id`

3. If missing, run the optimization SQL:
   ```bash
   psql -h localhost -U safecast -d safecast -f tools/add_performance_indexes.sql
   ```

### Admin page shows old data

The view needs to be refreshed:

```bash
./tools/refresh_track_stats.sh
```

### Stats don't match reality

If the stats seem wrong, force a complete refresh:

```sql
-- Drop and recreate
DROP MATERIALIZED VIEW IF EXISTS track_statistics;

-- Then run the create script
\i tools/create_track_stats_view.sql
```

## Performance Comparison

**Before (direct query):**
- Admin page load: 10-60+ seconds with 100M markers
- Database CPU: High during page load
- Blocks other queries while aggregating

**After (materialized view):**
- Admin page load: < 1 second (typically 0.3-0.5ms)
- Database CPU: Minimal
- No impact on other queries

## What Gets Tracked

The materialized view includes:

- `trackID` - Track identifier
- `marker_count` - Number of markers in track
- `first_date` - Earliest measurement timestamp
- `last_date` - Latest measurement timestamp
- `spectra_count` - Number of spectra in track
- `min_lat`, `max_lat` - Latitude bounds
- `min_lon`, `max_lon` - Longitude bounds
- `avg_doserate` - Average dose rate
- `max_doserate` - Maximum dose rate

## See Also

- [ADMIN_PERFORMANCE_FIX.md](ADMIN_PERFORMANCE_FIX.md) - Detailed technical explanation
- [tools/create_track_stats_view.sql](tools/create_track_stats_view.sql) - View creation SQL
- [tools/refresh_track_stats.sh](tools/refresh_track_stats.sh) - Refresh script
