# Fix: Duplicate Key Error on markers_pkey

## Error Message
```
ERROR: duplicate key value violates unique constraint "markers_pkey" (SQLSTATE 23505)
```

## Root Cause

PostgreSQL uses a **sequence** (`markers_id_seq`) to auto-generate unique IDs for the `markers` table's `id` column (defined as `BIGSERIAL PRIMARY KEY`).

This error occurs when the sequence value falls behind the actual maximum `id` in the table. Common causes:

1. **Database restore from backup** - Data restored but sequence not reset
2. **Manual data import** - Direct INSERT with explicit IDs bypasses sequence
3. **Database migration** - Migration from SQLite/DuckDB didn't sync the sequence
4. **Transaction rollback** - Rare edge case where sequence advanced but transaction rolled back

## Quick Fix

### Option 1: Run the SQL script directly

```bash
# Connect to your PostgreSQL database
psql -U safecast -d safecast -f fix_markers_sequence.sql
```

Or manually:
```sql
-- Sync sequence with actual data
SELECT setval('markers_id_seq', COALESCE((SELECT MAX(id) FROM markers), 1), true);
```

### Option 2: Use the Go utility

```bash
# Set your database URL
export DATABASE_URL="postgres://user:password@host:port/database?sslmode=disable"

# Run the fix utility
go run tools/fix_sequence.go
```

## Verification

After applying the fix, verify:

```sql
-- Check sequence is ahead of MAX(id)
SELECT 
    (SELECT last_value FROM markers_id_seq) AS sequence_value,
    (SELECT MAX(id) FROM markers) AS max_id;
```

The `sequence_value` should be greater than `max_id`.

## Prevention

The application code already handles this correctly by:
- Not specifying `id` in PostgreSQL INSERT statements (lets BIGSERIAL auto-generate)
- Using `ON CONFLICT ON CONSTRAINT markers_unique DO NOTHING` for duplicate handling

However, if you're seeing this error repeatedly, check:
1. No manual INSERTs are specifying explicit IDs
2. Database migrations properly sync sequences
3. Backup/restore procedures include sequence reset

## Technical Details

The `BIGSERIAL` type in PostgreSQL:
- Creates a sequence named `<table>_<column>_seq`
- Sets the column default to `nextval('<sequence>')`
- Does NOT automatically sync with existing data

To manually sync:
```sql
SELECT setval('markers_id_seq', (SELECT MAX(id) FROM markers), true);
```

The `true` parameter sets `is_called = true`, so the next `nextval()` returns `MAX(id) + 1`.
