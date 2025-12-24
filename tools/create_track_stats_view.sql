-- Create materialized view for track statistics
-- This pre-computes expensive aggregations for fast admin page loading

-- Drop existing view if it exists
DROP MATERIALIZED VIEW IF EXISTS track_statistics;

-- Create the materialized view
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

-- Create indexes on the materialized view for fast queries
CREATE INDEX idx_track_stats_trackid ON track_statistics(trackID);
CREATE INDEX idx_track_stats_last_date ON track_statistics(last_date DESC);
CREATE INDEX idx_track_stats_marker_count ON track_statistics(marker_count DESC);
CREATE INDEX idx_track_stats_spectra_count ON track_statistics(spectra_count DESC) WHERE spectra_count > 0;

-- Refresh the view immediately
REFRESH MATERIALIZED VIEW track_statistics;

-- Show summary
SELECT
    COUNT(*) as total_tracks,
    SUM(marker_count) as total_markers,
    SUM(spectra_count) as total_spectra
FROM track_statistics;
