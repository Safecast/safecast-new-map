-- Fix markers_pkey duplicate key errors caused by the BIGSERIAL sequence
-- falling behind the max id already in the table.
-- This happens after bulk imports that inserted rows with explicit IDs.
-- Safe to run at any time; does nothing if the sequence is already correct.

SELECT setval(
    pg_get_serial_sequence('markers', 'id'),
    GREATEST((SELECT MAX(id) FROM markers), 1)
);
