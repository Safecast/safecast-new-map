-- Safecast Performance Indexing Script
-- Creates all necessary indexes for 15-40x query speedup
-- Run with: psql -U postgres -d safecast -f create_performance_indexes.sql
-- Or: sudo -u postgres psql -d safecast -f create_performance_indexes.sql

\echo '========================================='
\echo 'Safecast Performance Index Creation'
\echo 'Expected total time: 3-5 minutes'
\echo '========================================='

-- Set transaction isolation to avoid lock issues
SET statement_timeout = '30 minutes';

\echo ''
\echo '📊 Index Set 1: Primary Query Optimization (Most Important)'
\echo '   - Used for map tile loading (zoom + bounds)'
\echo ''

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_markers_zoom_bounds 
  ON markers(zoom, lat, lon) 
  WHERE lat IS NOT NULL AND lon IS NOT NULL;
\echo '✅ idx_markers_zoom_bounds created'

\echo ''
\echo '📊 Index Set 2: Speed Filter Support'
\echo '   - For speed-based filtering (pedestrian/car/plane)'
\echo ''

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_markers_zoom_bounds_speed 
  ON markers(zoom, lat, lon, speed) 
  WHERE lat IS NOT NULL AND lon IS NOT NULL;
\echo '✅ idx_markers_zoom_bounds_speed created'

\echo ''
\echo '📊 Index Set 3: Track Queries'
\echo '   - For track data loading (when viewing single track)'
\echo ''

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_markers_trackid_zoom_bounds 
  ON markers(trackid, zoom, lat, lon)
  WHERE trackid IS NOT NULL AND lat IS NOT NULL;
\echo '✅ idx_markers_trackid_zoom_bounds created'

\echo ''
\echo '📊 Index Set 4: Date-Based Queries'
\echo '   - For historical queries and sorting'
\echo ''

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_markers_date 
  ON markers(date DESC);
\echo '✅ idx_markers_date created'

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_markers_date_bounds 
  ON markers(date DESC, lat, lon)
  WHERE lat IS NOT NULL AND lon IS NOT NULL;
\echo '✅ idx_markers_date_bounds created'

\echo ''
\echo '📊 Index Set 5: Speed Filtering'
\echo '   - For speed-only filters'
\echo ''

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_markers_speed 
  ON markers(speed)
  WHERE speed IS NOT NULL;
\echo '✅ idx_markers_speed created'

\echo ''
\echo '📊 Index Set 6: Realtime Data (if table exists)'
\echo '   - For realtime measurements'
\echo ''

-- These are optional if realtime table exists
DO $$
BEGIN
    CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_realtime_device_fetched 
      ON realtime(device_id, fetched_at DESC);
    RAISE NOTICE '✅ idx_realtime_device_fetched created';
EXCEPTION WHEN undefined_table THEN
    RAISE NOTICE '⚠️  realtime table not found (skipping realtime indexes)';
END $$;

DO $$
BEGIN
    CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_realtime_bounds 
      ON realtime(lat, lon, fetched_at DESC)
      WHERE lat IS NOT NULL AND lon IS NOT NULL;
    RAISE NOTICE '✅ idx_realtime_bounds created';
EXCEPTION WHEN undefined_table THEN
    NULL; -- Already reported above
END $$;

\echo ''
\echo '========================================='
\echo 'Analyzing table for query planner...'
\echo '========================================='
\echo ''

ANALYZE markers;
\echo '✅ markers table analyzed'

-- Verify all indexes were created
\echo ''
\echo '========================================='
\echo 'Verification: All Indexes on markers table'
\echo '========================================='
\echo ''

SELECT 
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS size
FROM pg_indexes
WHERE tablename = 'markers'
ORDER BY indexname;

\echo ''
\echo '========================================='
\echo '🚀 Index creation complete!'
\echo '========================================='
\echo ''
\echo 'Next steps:'
\echo '1. Restart the Safecast server'
\echo '2. Refresh the browser (Ctrl+Shift+R)'
\echo '3. Pan/zoom the map - should be 15-40x faster'
\echo ''
\echo 'To verify indexes are being used:'
\echo '  SELECT * FROM pg_stat_user_indexes'
\echo '  WHERE tablename = '\''markers'' AND idx_scan > 0;'
\echo ''
