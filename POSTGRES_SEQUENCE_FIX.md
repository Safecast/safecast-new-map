# PostgreSQL Sequence Synchronization Fix

## Problem

When uploading SPE files, you may encounter this error:
```
ERROR: duplicate key value violates unique constraint "markers_pkey" (SQLSTATE 23505)
```

## Root Cause

PostgreSQL uses `BIGSERIAL` columns which auto-increment using sequences. When data is imported with explicit IDs (e.g., from SQLite/ClickHouse migrations), the sequence counter doesn't automatically update. This causes the sequence to fall behind the actual maximum ID in the table.

When new data is inserted, PostgreSQL tries to use the next sequence value, which may already exist in the table, causing a primary key violation.

## Quick Fix

Run the sequence reset script:

```bash
export POSTGRES_URL='host=localhost port=5432 dbname=safecast user=safecast password=yourpassword sslmode=disable'
./tools/reset_postgres_sequences.sh
```

Or manually reset sequences:

```sql
-- Reset markers sequence
SELECT setval('markers_id_seq', (SELECT COALESCE(MAX(id), 1) FROM markers));

-- Reset spectra sequence
SELECT setval('spectra_id_seq', (SELECT COALESCE(MAX(id), 1) FROM spectra));

-- Reset uploads sequence
SELECT setval('uploads_id_seq', (SELECT COALESCE(MAX(id), 1) FROM uploads));
```

## Prevention

The migration tool ([tools/migrate_to_postgres.go](tools/migrate_to_postgres.go)) has been updated to automatically reset sequences after migration. Future migrations will not have this issue.

## When to Use the Fix

Run the sequence reset script whenever:
- You encounter primary key violation errors on inserts
- After manually importing data with explicit IDs
- After restoring from a backup
- After running data migrations

## Technical Details

**Affected Tables:**
- `markers` (uses `markers_id_seq`)
- `spectra` (uses `spectra_id_seq`)
- `uploads` (uses `uploads_id_seq`)

**Schema Definition:**
```sql
CREATE TABLE markers (
  id BIGSERIAL PRIMARY KEY,  -- Auto-incrementing, uses markers_id_seq
  ...
);
```

**How the Fix Works:**
```sql
SELECT setval('markers_id_seq', (SELECT COALESCE(MAX(id), 1) FROM markers));
```

This sets the sequence to the maximum ID currently in the table, ensuring the next auto-generated ID won't conflict.

## Files Modified

1. [tools/migrate_to_postgres.go](tools/migrate_to_postgres.go) - Added `resetSequences()` function
2. [tools/reset_postgres_sequences.sh](tools/reset_postgres_sequences.sh) - New standalone script

## Resolution Log

**Date:** 2025-12-24

**Issue:** SPE upload failed with "duplicate key value violates unique constraint markers_pkey"

**Diagnosis:**
- Max ID in markers table: 115780507
- Sequence last_value: 115780504
- Sequence was 3 IDs behind actual data

**Action Taken:**
1. Reset sequence to match max ID
2. Updated migration tool to prevent future occurrences
3. Created standalone script for easy sequence resets

**Status:** ✅ Resolved
