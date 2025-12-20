# Spectral Data Migration Guide

This guide explains how to add spectral data support to an existing safecast-new-map database without losing any data.

## Overview

Spectral data support requires:
1. A new `spectra` table to store gamma spectrum data
2. A new `has_spectrum` column in the `markers` table to flag which markers have spectra

## Migration Scripts

Three migration scripts are provided, one for each database type:

### PostgreSQL (Recommended for Production)
```bash
./migrate_add_spectra_postgresql.sh
```

**Environment Variables:**
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_NAME` - Database name (default: safecast)
- `DB_PASSWORD` - Database password (optional)

**Example:**
```bash
DB_HOST=db.example.com DB_USER=safecast DB_NAME=radiation_data ./migrate_add_spectra_postgresql.sh
```

### SQLite
```bash
./migrate_add_spectra_sqlite.sh /path/to/database.db
```

**Example:**
```bash
./migrate_add_spectra_sqlite.sh /var/lib/safecast/safecast.db
```

### DuckDB
```bash
./migrate_add_spectra_duckdb.sh /path/to/database.duckdb
```

**Example:**
```bash
./migrate_add_spectra_duckdb.sh /var/lib/safecast/safecast.duckdb
```

## What the Migration Does

### 1. Adds `has_spectrum` Column
Adds a boolean/integer flag to the `markers` table:
- PostgreSQL/DuckDB: `BOOLEAN DEFAULT FALSE`
- SQLite: `INTEGER DEFAULT 0`

### 2. Creates `spectra` Table
Creates a new table to store gamma spectrum data with columns:
- `id` - Primary key
- `marker_id` - Foreign key to markers table
- `channels` - JSON array of channel counts
- `channel_count` - Number of channels (typically 1024)
- `energy_min_kev` - Minimum energy in keV
- `energy_max_kev` - Maximum energy in keV
- `live_time_sec` - Live time in seconds
- `real_time_sec` - Real time in seconds
- `device_model` - Detector model
- `calibration` - JSON energy calibration coefficients
- `source_format` - Original file format (spe, n42, rctrk)
- `filename` - Original filename
- `raw_data` - Raw spectrum file data (BLOB)
- `created_at` - Timestamp

### 3. Creates Index
Creates an index on `spectra.marker_id` for fast lookups.

## Safety Features

All migration scripts include:
- ✅ **Automatic backups** (SQLite/DuckDB) or backup recommendations (PostgreSQL)
- ✅ **IF NOT EXISTS** checks to prevent errors on re-run
- ✅ **No data deletion** - only adds new tables/columns
- ✅ **Confirmation prompts** before making changes
- ✅ **Verification steps** to confirm successful migration
- ✅ **Error handling** with clear error messages

## Before You Start

### 1. Check Your Database Type
```bash
# Look at your startup command or config
./safecast-new-map -db-type pgx    # PostgreSQL
./safecast-new-map -db-type sqlite # SQLite
./safecast-new-map -db-type duckdb # DuckDB
```

### 2. Create a Backup (Important!)

**PostgreSQL:**
```bash
pg_dump -h hostname -U username database_name > backup_$(date +%Y%m%d).sql
```

**SQLite:**
```bash
cp /path/to/database.db /path/to/database.db.backup
```

**DuckDB:**
```bash
cp /path/to/database.duckdb /path/to/database.duckdb.backup
```

### 3. Check Disk Space
Spectral data can be large. Ensure you have sufficient disk space:
```bash
df -h /path/to/database/directory
```

## Running the Migration

### Step 1: Make Scripts Executable
```bash
chmod +x migrate_add_spectra_*.sh
```

### Step 2: Run the Appropriate Script

**For PostgreSQL (500GB production database):**
```bash
# Set your database connection details
export DB_HOST=your-db-host.com
export DB_USER=safecast
export DB_NAME=safecast_production
export DB_PASSWORD=your_password

# Run migration
./migrate_add_spectra_postgresql.sh
```

**For SQLite:**
```bash
./migrate_add_spectra_sqlite.sh /path/to/your/database.db
```

**For DuckDB:**
```bash
./migrate_add_spectra_duckdb.sh /path/to/your/database.duckdb
```

### Step 3: Verify Migration

The script will automatically verify the migration, but you can also check manually:

**PostgreSQL:**
```sql
-- Check column exists
SELECT column_name, data_type
FROM information_schema.columns
WHERE table_name = 'markers' AND column_name = 'has_spectrum';

-- Check table exists
\d spectra
```

**SQLite:**
```bash
sqlite3 database.db "PRAGMA table_info(markers);" | grep has_spectrum
sqlite3 database.db ".schema spectra"
```

**DuckDB:**
```bash
duckdb database.duckdb "DESCRIBE markers;" | grep has_spectrum
duckdb database.duckdb "DESCRIBE spectra;"
```

## After Migration

### 1. Restart Application
```bash
# Stop current instance
pkill safecast-new-map

# Start with your normal configuration
./safecast-new-map -db-type pgx -db-conn "your_connection_string"
```

### 2. Test Spectral Data Upload

**Via Web Interface:**
1. Open http://localhost:8765
2. Click the green "Upload" button
3. Upload a spectrum file (.spe, .n42, or .rctrk)
4. The marker should now have a spectrum icon

**Via Test Script:**
```bash
go run scripts/insert_test_spectrum.go
```

### 3. Query Spectral Data

**PostgreSQL/DuckDB:**
```sql
-- Count markers with spectra
SELECT COUNT(*) FROM markers WHERE has_spectrum = TRUE;

-- View spectrum data
SELECT id, marker_id, device_model, source_format, created_at
FROM spectra
LIMIT 10;

-- Find markers with spectra in a region
SELECT m.*, s.device_model, s.source_format
FROM markers m
JOIN spectra s ON m.id = s.marker_id
WHERE m.lat BETWEEN 35.0 AND 38.0
  AND m.lon BETWEEN 139.0 AND 141.0;
```

**SQLite:**
```sql
-- Count markers with spectra
SELECT COUNT(*) FROM markers WHERE has_spectrum = 1;

-- Same queries as above
```

## Troubleshooting

### "Column already exists" Error
This is normal if you run the migration twice. The script detects existing columns and skips them.

### "Permission denied" Error (PostgreSQL)
Make sure your database user has ALTER TABLE privileges:
```sql
GRANT ALTER ON TABLE markers TO your_user;
```

### "Out of disk space" Error
Free up space or expand your disk before importing spectrum files.

### Migration Failed Partway
- PostgreSQL: Check the error message and fix the issue, then re-run (safe to re-run)
- SQLite/DuckDB: Restore from backup and try again

## Rolling Back

If you need to remove spectral data support:

**PostgreSQL/DuckDB:**
```sql
-- Remove spectra table
DROP TABLE IF EXISTS spectra CASCADE;

-- Remove has_spectrum column
ALTER TABLE markers DROP COLUMN IF EXISTS has_spectrum;
```

**SQLite:**
```sql
-- Remove spectra table
DROP TABLE IF EXISTS spectra;

-- Remove has_spectrum column (SQLite doesn't support DROP COLUMN easily)
-- You'll need to recreate the table without the column or restore from backup
```

## Performance Considerations

### Large Databases (500GB+)

1. **PostgreSQL:**
   - Adding a column with DEFAULT is fast (metadata change only)
   - No table rewrite occurs
   - Migration should complete in seconds

2. **Indexing:**
   - The `idx_spectra_marker_id` index is created empty
   - Index grows as you add spectra

3. **Spectrum Storage:**
   - Each spectrum is ~100KB-1MB depending on format
   - 1000 spectra ≈ 100MB-1GB
   - Plan disk space accordingly

## Support

If you encounter issues:
1. Check the error message in the script output
2. Verify database permissions
3. Check disk space
4. Review application logs
5. Create a GitHub issue with the error details

## Summary

✅ Safe migration with automatic backups
✅ No data loss - only adds tables/columns
✅ Can be run multiple times safely
✅ Minimal downtime
✅ Automatic verification

The migration is designed to be safe for production databases with 500GB+ of data.
