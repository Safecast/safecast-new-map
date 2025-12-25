-- Backfill recording_date for existing uploads
-- This updates the recording_date column with the earliest marker date for each track

\echo 'Starting recording_date backfill...'

UPDATE uploads u 
SET recording_date = to_timestamp((SELECT MIN(date) FROM markers WHERE trackID = u.track_id))
WHERE (recording_date IS NULL OR EXTRACT(EPOCH FROM recording_date) = 0) 
  AND track_id IS NOT NULL;

\echo 'Recording_date backfill complete!'
\echo 'Updated rows:'
SELECT COUNT(*) as updated_count FROM uploads WHERE recording_date IS NOT NULL;
