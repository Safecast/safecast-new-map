-- Fix for: ERROR: duplicate key value violates unique constraint "markers_pkey"
-- This happens when PostgreSQL's sequence is out of sync with actual data

-- Check current sequence value
SELECT last_value, is_called FROM markers_id_seq;

-- Check actual MAX(id) in the table
SELECT MAX(id) FROM markers;

-- Fix: Reset the sequence to be after the maximum ID
-- This ensures the next generated ID won't conflict with existing data
SELECT setval('markers_id_seq', COALESCE((SELECT MAX(id) FROM markers), 1), true);

-- Verify the fix
SELECT last_value, is_called FROM markers_id_seq;
