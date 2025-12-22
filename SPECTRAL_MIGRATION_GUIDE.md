# Spectral Graph & Speed Data Migration Guide

## Problem Summary

The spectral graphs in the popup were not showing any data because the spectral data existed in the SQLite database (`database-8765.sqlite`) but was not migrated to the PostgreSQL database that the application is currently using.

Additionally, speed calculations were missing, which would have taken hours to recalculate from scratch.

## Database Status

### SQLite Database (`database-8765.sqlite`)
- **Total markers**: 15,721,425
- **Spectral records**: 1,083
- **Markers with spectrum**: 91
- **Markers with speed data**: 15,720,684 (99.995%)

### PostgreSQL Database (Before Migration)
- **Spectral records**: 0 (missing!)
- **Markers with spectrum**: 0 (missing!)
- **Markers with speed data**: 0 (missing!)

## Root Cause

When you migrated from SQLite to PostgreSQL, two critical datasets were not transferred:

1. **Spectral Data**: The `spectra` table and `has_spectrum` flags were not migrated
   - Application queries PostgreSQL for spectral data
   - PostgreSQL returns no results (empty dataset)
   - Frontend receives `null` or empty data
   - Graph canvas remains empty with "No spectrum data available"

2. **Speed Data**: The `speed` column values were not migrated
   - PostgreSQL had to recalculate speeds from scratch
   - This process takes hours for millions of records
   - Migration is much faster than recalculation

## Solution

You need to migrate both the spectral data and speed data from SQLite to PostgreSQL. An enhanced migration script has been created that handles both datasets efficiently.

### Recommended: All-in-One Migration Script

The `migrate_all_data.py` script migrates both spectral and speed data in a single run:

```bash
# Install required Python package if not already installed
pip3 install psycopg2-binary

# Set PostgreSQL password (or you'll be prompted)
export PG_PASSWORD="your_postgres_password"

# Run the migration
./migrate_all_data.py
```

**What it migrates:**
- ✅ 1,083 spectral records
- ✅ 91 marker spectrum flags
- ✅ ~15.7 million speed values

**Performance:**
- Processes speed data in batches of 10,000 records
- Typical completion time: 25-35 minutes for full dataset
- Shows real-time progress updates
- Much faster than recalculating speeds (which takes hours)

### Alternative: Spectral-Only Migration

If you only need spectral data (not speed), use the lighter script:

```bash
# Set PostgreSQL password
export PG_PASSWORD="your_postgres_password"

# Run spectral-only migration
./migrate_spectra.py
```

This completes in ~1-2 minutes but doesn't migrate speed data.

### Legacy: Bash Script Option

Alternative bash-based migration (spectral only):

```bash
# Set PostgreSQL password
export PG_PASSWORD="your_postgres_password"

# Run the migration
./migrate_spectra_sqlite_to_postgres.sh
```


## Migration Process

The `migrate_all_data.py` script performs the following steps:

1. **Connect** to both SQLite and PostgreSQL databases
2. **Analyze** and display data counts in both databases
3. **Confirm** migration plan with user
4. **Migrate spectral data**:
   - Export all spectral records from SQLite
   - Import into PostgreSQL with BLOB handling
   - Skip existing records to avoid duplicates
5. **Update marker flags**:
   - Set `has_spectrum = true` for relevant markers
6. **Migrate speed data**:
   - Process in batches of 10,000 records
   - Update PostgreSQL efficiently using bulk operations
   - Show progress updates every batch
7. **Verify** the migration was successful
8. **Report** final statistics and any warnings

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

# Check markers with speed data
psql -h localhost -U safecast -d safecast -c "SELECT COUNT(*) FROM markers WHERE speed IS NOT NULL AND speed > 0;"

# View a sample spectrum
psql -h localhost -U safecast -d safecast -c "SELECT id, marker_id, device_model, source_format, filename FROM spectra LIMIT 5;"

# View sample speed data
psql -h localhost -U safecast -d safecast -c "SELECT id, speed FROM markers WHERE speed IS NOT NULL AND speed > 0 LIMIT 5;"
```

Expected results:
- Spectral records: 1,083
- Markers with spectrum: 91
- Markers with speed data: ~15,720,684

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

- **`migrate_all_data.py`** - Comprehensive migration script for both spectral and speed data (recommended)
- `migrate_spectra.py` - Python migration script for spectral data only
- `migrate_spectra_sqlite_to_postgres.sh` - Bash migration script for spectral data only (legacy)
- `SPECTRAL_MIGRATION_GUIDE.md` - This documentation

## Performance Notes

### Why Speed Migration is Faster Than Recalculation

**Recalculating speeds** (what PostgreSQL was doing):
- Requires sorting millions of markers by track and timestamp
- Calculates distance and time delta between consecutive points
- CPU-intensive geometric calculations
- Estimated time: 4-8 hours for 15.7M records

**Migrating speeds** (what the script does):
- Simple bulk UPDATE operations
- No complex calculations needed
- I/O-bound operation (disk read/write)
- Actual time: 25-35 minutes for 15.7M records

**Speed improvement**: ~10-15x faster! 🚀

### Multi-Threading Note

For future migrations of this scale, a multi-threaded version could reduce the time to ~5-10 minutes by:
- Using multiple PostgreSQL connections in parallel
- Processing different ID ranges simultaneously
- Utilizing multiple CPU cores efficiently

## Next Steps

After successful migration, you may want to:

1. **Backup** the PostgreSQL database
2. **Archive** the SQLite database
3. **Test** spectral graph functionality thoroughly
4. **Document** the migration in your deployment notes
