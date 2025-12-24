# Admin Page Performance Optimization

## Problem

The admin tracks page ([/api/admin/tracks](http://localhost:8080/api/admin/tracks)) was extremely slow when loading, taking several seconds or even minutes with many tracks in the database.

## Root Cause

**Database Size:**
- 8,285 tracks
- 97,886,265 markers (nearly 98 million!)
- 105 spectra

**Query Performance Issue:**

The original query joined the `tracks` and `markers` tables, then aggregated all markers for each track on every page load:

```sql
SELECT
    t.trackID,
    COUNT(DISTINCT m.id) as marker_count,
    MIN(m.date) as first_date,
    MAX(m.date) as last_date,
    COALESCE(SUM(CASE WHEN m.has_spectrum = true THEN 1 ELSE 0 END), 0) as spectra_count
FROM tracks t
LEFT JOIN markers m ON t.trackID = m.trackID
WHERE t.trackID NOT LIKE 'live:%'
GROUP BY t.trackID
ORDER BY last_date DESC
LIMIT 1000
```

This query:
1. Scans nearly 98 million markers
2. Groups by 8,285 track IDs
3. Computes aggregates (COUNT, MIN, MAX, SUM) for each group
4. Sorts all results
5. Only then applies the LIMIT

Even with proper indexes (`idx_markers_trackid`), this is extremely expensive.

## Solution: Materialized View

Created a **materialized view** that pre-computes track statistics once, making queries instant.

### Files Created

1. [tools/create_track_stats_view.sql](tools/create_track_stats_view.sql) - SQL to create the materialized view
2. [tools/refresh_track_stats.sh](tools/refresh_track_stats.sh) - Script to refresh the view periodically

### Files Modified

1. [safecast-new-map.go:5817-5846](safecast-new-map.go#L5817-L5846) - Updated `adminTracksHandler` to use materialized view

### Materialized View Schema

```sql
CREATE MATERIALIZED VIEW track_statistics AS
SELECT
    t.trackID,
    COUNT(DISTINCT m.id) as marker_count,
    MIN(m.date) as first_date,
    MAX(m.date) as last_date,
    COALESCE(SUM(CASE WHEN m.has_spectrum = true THEN 1 ELSE 0 END), 0) as spectra_count,
    MIN(m.lat) as min_lat,
    MAX(m.lat) as max_lat,
    MIN(m.lon) as min_lon,
    MAX(m.lon) as max_lon,
    AVG(m.doserate) as avg_doserate,
    MAX(m.doserate) as max_doserate
FROM tracks t
LEFT JOIN markers m ON t.trackID = m.trackID
WHERE t.trackID NOT LIKE 'live:%'
GROUP BY t.trackID;
```

### Indexes on Materialized View

```sql
CREATE INDEX idx_track_stats_trackid ON track_statistics(trackID);
CREATE INDEX idx_track_stats_last_date ON track_statistics(last_date DESC);
CREATE INDEX idx_track_stats_marker_count ON track_statistics(marker_count DESC);
CREATE INDEX idx_track_stats_spectra_count ON track_statistics(spectra_count DESC) WHERE spectra_count > 0;
```

## Performance Improvement

**Before:**
- Query time: Several seconds to minutes
- Scanned: ~98 million markers on every page load

**After:**
- Query time: **0.378 milliseconds** (less than 1ms!)
- Scanned: Pre-computed view with 8,280 rows

**Improvement: ~10,000x faster or more!**

### Query Plan (After Optimization)

```
Limit  (cost=0.29..89.20 rows=1000 width=64) (actual time=0.020..0.346 rows=1000 loops=1)
  ->  Index Scan using idx_track_stats_last_date on track_statistics
Planning Time: 0.659 ms
Execution Time: 0.378 ms
```

## Maintenance

The materialized view needs periodic refreshing to include new data.

### Manual Refresh

```bash
export POSTGRES_URL='host=localhost port=5432 dbname=safecast user=safecast password=yourpassword sslmode=disable'
./tools/refresh_track_stats.sh
```

Or directly in PostgreSQL:

```sql
REFRESH MATERIALIZED VIEW CONCURRENTLY track_statistics;
```

### Automatic Refresh (Recommended)

Add to crontab for hourly refresh:

```bash
crontab -e
```

Add this line:

```
0 * * * * cd /home/rob/Documents/Safecast/safecast-new-map/tools && PGPASSWORD=safecast2025 ./refresh_track_stats.sh >> /var/log/track_stats_refresh.log 2>&1
```

### Refresh After Data Import

After importing large datasets, refresh the view:

```bash
PGPASSWORD=safecast2025 ./tools/refresh_track_stats.sh
```

## How It Works

1. **Initial Creation**: The materialized view computes expensive aggregations once
2. **Admin Query**: Reads from the pre-computed view (instant)
3. **Periodic Refresh**: Update the view with new data (using `CONCURRENTLY` to avoid locking)

The `CONCURRENTLY` option allows queries to continue using the old data while the refresh is in progress, ensuring zero downtime.

## Additional Features in View

The materialized view includes extra columns for future use:
- `min_lat`, `max_lat`, `min_lon`, `max_lon` - Track bounding box
- `avg_doserate`, `max_doserate` - Radiation statistics

These can be displayed in the admin UI or used for filtering.

## Resolution Log

**Date:** 2025-12-24

**Issue:** Admin tracks page extremely slow with many tracks

**Diagnosis:**
- 8,285 tracks, 97.9 million markers
- Query aggregating all markers on every page load
- Execution time: several seconds to minutes

**Action Taken:**
1. Created materialized view `track_statistics`
2. Added indexes for optimal query performance
3. Updated admin handler to use materialized view
4. Created refresh script for periodic updates

**Performance Result:**
- Query time reduced from seconds/minutes to **0.378ms**
- ~10,000x performance improvement

**Status:** ✅ Resolved
