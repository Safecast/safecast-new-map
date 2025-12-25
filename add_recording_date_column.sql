-- Add recording_date column to uploads table and backfill with data

\echo 'Adding recording_date column to uploads table...'

-- Add the column if it doesn't exist
ALTER TABLE uploads ADD COLUMN IF NOT EXISTS recording_date TIMESTAMP;

\echo 'Backfilling recording_date for existing uploads...'

-- Update recording_date with the earliest marker date for each track
UPDATE uploads u 
SET recording_date = to_timestamp((SELECT MIN(date) FROM markers WHERE trackID = u.track_id))
WHERE recording_date IS NULL AND track_id IS NOT NULL;

\echo 'Migration complete!'
\echo 'Total uploads with recording_date:'
SELECT COUNT(*) as updated_count FROM uploads WHERE recording_date IS NOT NULL;
