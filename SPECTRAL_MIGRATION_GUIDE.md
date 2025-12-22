# Spectral Graph Data Migration Guide

## Problem Summary

The spectral graphs in the popup are not showing any data because the spectral data exists in the SQLite database (`database-8765.sqlite`) but was not migrated to the PostgreSQL database that the application is currently using.

## Database Status

### SQLite Database (`database-8765.sqlite`)
- **Spectral records**: 1,083
- **Markers with spectrum**: 91

### PostgreSQL Database
- **Spectral records**: 0 (missing!)
- **Markers with spectrum**: 0 (missing!)

## Root Cause

When you migrated from SQLite to PostgreSQL, the spectral data in the `spectra` table was not transferred. This means:

1. The application is querying PostgreSQL for spectral data
2. PostgreSQL returns no results (empty dataset)
3. The frontend receives `null` or empty data
4. The graph canvas remains empty with the message "No spectrum data available"

## Solution

You need to migrate the spectral data from SQLite to PostgreSQL. I've created two migration scripts for you:

### Option 1: Python Script (Recommended)

The Python script properly handles BLOB data and provides better error handling:

```bash
# Install required Python package if not already installed
pip3 install psycopg2-binary

# Set PostgreSQL password (or you'll be prompted)
export PG_PASSWORD="your_postgres_password"

# Run the migration
./migrate_spectra.py
```

### Option 2: Bash Script

Alternative bash-based migration:

```bash
# Set PostgreSQL password
export PG_PASSWORD="your_postgres_password"

# Run the migration
./migrate_spectra_sqlite_to_postgres.sh
```

## Migration Process

The migration scripts will:

1. **Connect** to both SQLite and PostgreSQL databases
2. **Verify** the data counts in both databases
3. **Export** all spectral records from SQLite
4. **Import** spectral data into PostgreSQL
5. **Update** marker flags (`has_spectrum = true`)
6. **Verify** the migration was successful

## Configuration

Both scripts use environment variables for configuration:

```bash
export SQLITE_DB="database-8765.sqlite"  # Source SQLite database
export PG_HOST="localhost"                # PostgreSQL host
export PG_PORT="5432"                     # PostgreSQL port
export PG_USER="safecast"                 # PostgreSQL user
export PG_DB="safecast"                   # PostgreSQL database name
export PG_PASSWORD="your_password"        # PostgreSQL password
```

## After Migration

Once the migration is complete:

1. **Restart** the safecast-new-map application
2. **Open** the map in your browser
3. **Click** on a marker that has spectral data (indicated by a special icon)
4. **View** the spectrum modal - the graph should now display correctly

## Verification

To verify the migration was successful, check the PostgreSQL database:

```bash
# Check spectral records count
psql -h localhost -U safecast -d safecast -c "SELECT COUNT(*) FROM spectra;"

# Check markers with spectrum
psql -h localhost -U safecast -d safecast -c "SELECT COUNT(*) FROM markers WHERE has_spectrum = true;"

# View a sample spectrum
psql -h localhost -U safecast -d safecast -c "SELECT id, marker_id, device_model, source_format, filename FROM spectra LIMIT 5;"
```

Expected results:
- Spectral records: 1,083
- Markers with spectrum: 91

## Technical Details

### Spectral Data Structure

Each spectrum record contains:
- **channels**: JSON array of counts per energy channel (typically 1024-11000 channels)
- **calibration**: JSON object with energy calibration coefficients (a, b, c)
- **raw_data**: BLOB containing the original spectrum file
- **device_model**: Detector model (e.g., "OSASP", "RadiaCode-102")
- **source_format**: Original file format ("n42", "spe", "rctrk")

### Frontend Rendering

The frontend (`map.html`) uses the following flow:

1. User clicks on a marker with `hasSpectrum = true`
2. JavaScript calls `/api/spectrum/{markerID}`
3. Backend queries PostgreSQL `spectra` table
4. Returns JSON with channels, calibration, and metadata
5. `drawSpectrumChart()` renders the data on a canvas element

### Why the Graph Was Empty

The API endpoint was returning:
- **404 Not Found** - No spectrum found for this marker
- Or **null** data - Empty result set

This caused the frontend to display the fallback message instead of rendering the graph.

## Troubleshooting

### If migration fails:

1. **Check PostgreSQL credentials**:
   ```bash
   psql -h localhost -U safecast -d safecast -c "SELECT version();"
   ```

2. **Verify SQLite database path**:
   ```bash
   ls -lh database-8765.sqlite
   ```

3. **Check table schemas match**:
   ```bash
   # SQLite
   sqlite3 database-8765.sqlite ".schema spectra"
   
   # PostgreSQL
   psql -h localhost -U safecast -d safecast -c "\d spectra"
   ```

4. **Check for foreign key constraints**:
   - Ensure all `marker_id` values in `spectra` exist in the `markers` table

### If graphs still don't show after migration:

1. **Clear browser cache** and reload the page
2. **Check browser console** for JavaScript errors
3. **Verify API response**:
   ```bash
   curl http://localhost:8765/api/spectrum/1350
   ```
4. **Check application logs** for database errors

## Files Created

- `migrate_spectra.py` - Python migration script (recommended)
- `migrate_spectra_sqlite_to_postgres.sh` - Bash migration script (alternative)
- `SPECTRAL_MIGRATION_GUIDE.md` - This documentation

## Next Steps

After successful migration, you may want to:

1. **Backup** the PostgreSQL database
2. **Archive** the SQLite database
3. **Test** spectral graph functionality thoroughly
4. **Document** the migration in your deployment notes
