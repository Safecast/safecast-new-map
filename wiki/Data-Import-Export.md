# Data Import & Export

Comprehensive guide to importing and exporting radiation data.

[[Home|← Back to Home]]

---

## Supported Import Formats

The platform auto-detects and parses multiple file formats:

### Measurement Files

| Format | Extension | Description |
|--------|-----------|-------------|
| **KML/KMZ** | `.kml`, `.kmz` | Google Earth, Safecast bGeigie tracks |
| **JSON** | `.json` | Exported track format |
| **RCTRK** | `.rctrk` | RadiaCode devices (JSON with embedded spectra) |
| **CSV** | `.csv` | AtomFast and custom formats |
| **GPX** | `.gpx` | GPS tracks |
| **LOG** | `.log` | bGeigie Nano/Zen |
| **Radiacode CSV** | `.csv` | Tab-separated with DoseRate, CountRate, GPS columns |

### Spectrum Files

| Format | Extension | Description |
|--------|-----------|-------------|
| **Maestro** | `.spe` | Standard spectrum format |
| **ANSI N42.42** | `.n42` | Industry standard format |
| **RadiaCode** | `.rctrk`, `.rcxml`, `.xml` | RadiaCode spectrum files |

See [Spectral Analysis](Spectral-Analysis) for spectrum-specific documentation.

---

## File Upload via Web Interface

### Upload Process

1. Log in to your account (if authentication enabled)
2. Navigate to the map interface
3. Click the upload button
4. Select files from your device
5. Wait for processing completion

### Upload Limits

- Files are validated for format and data integrity
- Processing time varies by file size
- Large files may take several minutes

### Track Metadata

After upload, you can update metadata:
- Recording date
- Detector type
- Username
- Notes and comments

---

## Bulk Import

### Import from Remote Archive

Download and import a `.tgz` archive from a URL:

```bash
./safecast-new-map -import-tgz-url https://simplemap.safecast.org/api/json/weekly.tgz
```

This is the fastest way to get started with production data.

### Import from Local Archive

Import a local `.tgz` file:

```bash
./safecast-new-map -import-tgz-file /path/to/archive.tgz
```

### Archive Structure

The `.tgz` archive should contain:
- Individual JSON track files
- Organized by date or track ID
- Compatible with the platform's JSON export format

---

## Automated Data Sync

### Safecast API Fetcher

Automatically poll `api.safecast.org` for approved bGeigie imports:

```bash
./safecast-new-map \
  -safecast-fetcher \
  -safecast-fetcher-interval 5m \
  -safecast-fetcher-batch-size 10 \
  -safecast-fetcher-start-date 2024-01-01
```

#### Fetcher Options

| Flag | Default | Description |
|------|---------|-------------|
| `-safecast-fetcher` | false | Enable the fetcher |
| `-safecast-fetcher-interval` | 5m | Polling interval (e.g., 5m, 1h) |
| `-safecast-fetcher-batch-size` | 10 | Records per API call |
| `-safecast-fetcher-start-date` | - | Start date (YYYY-MM-DD) |
| `-safecast-fetcher-backfill` | false | Import ALL historical records |
| `-safecast-fetcher-newest-first` | false | Fetch newest records first |

#### Backfill Mode

For importing all historical data from a start date:

```bash
./safecast-new-map \
  -safecast-fetcher \
  -safecast-fetcher-backfill \
  -safecast-fetcher-start-date 2012-01-01
```

**Note:** Backfill mode ignores database state and imports all records from the start date.

#### Batch Import Script

Use the provided script for large imports:

```bash
# Complete backfill (~4,775 files)
./backfill_complete.sh

# Incremental backfill (~250 files per run)
./backfill_2015_2017.sh
```

See [BACKFILL_GUIDE.md](/BACKFILL_GUIDE.md) for details.

### Real-Time Sensor Polling

Poll live Safecast device data (Pointcast, Solarcast, bGeigieZen sensors):

```bash
./safecast-new-map -safecast-realtime
```

Sensors appear on the map in real-time with current readings.

**Features:**
- Automatic unit conversion to µSv/h
- Sensors displayed on map with current values
- History tracking for each sensor

---

## Data Export

### JSON Archive Export

Export data as compressed JSON archives:

```
GET /api/json/weekly.tgz
```

#### Archive Frequencies

Configure generation frequency with `-json-archive-frequency`:

- `daily` - Daily archives
- `weekly` - Weekly archives (default)
- `monthly` - Monthly archives
- `yearly` - Yearly archives

### Individual Track Export

Export specific tracks as JSON:

```
GET /api/track/{id}
```

Returns JSON with all measurement data.

### Spectrum Export

Export spectrum data in multiple formats:

```
GET /api/spectrum/{marker_id}/download?format=json
GET /api/spectrum/{marker_id}/download?format=csv
GET /api/spectrum/{marker_id}/download?format=n42
GET /api/spectrum/{marker_id}/download?format=spe
```

---

## Import via API

### Programmatic Upload

Upload files via the REST API:

```bash
curl -X POST http://localhost:8765/api/upload \
  -H "X-API-Key: your-api-key" \
  -F "file=@track.kml" \
  -F "source=manual"
```

### Import from Safecast API

Admin endpoint to import specific records:

```bash
curl -X POST "http://localhost:8765/api/admin/import-from-safecast?password=admin-password" \
  -H "Content-Type: application/json" \
  -d '{"safecast_id": 12345}'
```

### Import by ID Range

Import multiple records by ID:

```bash
curl -X POST "http://localhost:8765/api/admin/import-by-id?password=admin-password" \
  -H "Content-Type: application/json" \
  -d '{"start_id": 100000, "end_id": 110000}'
```

---

## Data Validation

### Auto-Detection

The platform automatically:
- Detects file format from extension and content
- Validates GPS coordinates
- Checks timestamp validity
- Parses radiation measurements
- Extracts spectrum data when available

### Error Handling

- Invalid files are rejected with error messages
- Partial imports log warnings
- Duplicate tracks are skipped
- Database constraints ensure data integrity

---

## Database Optimization

### Post-Import Maintenance

For DuckDB/SQLite databases, automatic maintenance runs after import:
- Optimize tables
- Checkpoint WAL files
- Vacuum deleted records

For PostgreSQL:
- Consider running `VACUUM ANALYZE` after large imports
- Update statistics for query planning

### Track Statistics

Refresh materialized views after bulk imports:

```bash
./tools/refresh_track_stats.sh
```

---

## Migration Tools

### Migrate from SQLite to PostgreSQL

```bash
go run ./cmd/tools/migrate-to-postgres \
  -source /path/to/data.db \
  -dest "postgres://user:pass@localhost/dbname"
```

See [MIGRATION_GUIDE.md](/MIGRATION_GUIDE.md) for details.

### Import API Keys

Import API keys from CSV:

```bash
go run ./cmd/tools/import-api-keys \
  -file api_keys.csv \
  -db-conn "postgres://..."
```

CSV format:
```csv
user_id,api_key
1,abcdefghijklmnopqrst
2,uvwxyz1234567890abcd
```

---

## Best Practices

### Large Imports

1. **Use PostgreSQL** for production with large datasets
2. **Batch imports** - Use fetcher with appropriate batch sizes
3. **Monitor disk space** - Large imports require significant storage
4. **Index after import** - Run database maintenance commands

### Data Quality

1. **Validate GPS coordinates** before upload
2. **Check timestamps** for accuracy
3. **Use consistent detector names** for better tracking
4. **Add metadata** (username, notes) for traceability

### Performance

1. **Use `-import-tgz-url`** for fastest initial data load
2. **Enable fetcher** for ongoing sync
3. **Schedule maintenance** during low-traffic periods
4. **Monitor import logs** for errors

---

## Troubleshooting

### Import Fails

**Check file format:**
```bash
# Verify file is valid JSON
jq . track.json

# Check KML structure
head -20 track.kml
```

**Check logs:**
```bash
# Look for import errors
grep "import" /var/log/safecast.log
```

### Duplicate Data

The platform skips duplicate tracks based on track ID. If you need to reimport:

1. Delete existing track via admin panel
2. Reimport the file

### Slow Import

**Optimize database:**
```sql
-- PostgreSQL
VACUUM ANALYZE markers;
VACUUM ANALYZE tracks;

-- DuckDB
CHECKPOINT;
OPTIMIZE;
```

**Increase batch size:**
```bash
./safecast-new-map -safecast-fetcher-batch-size 50
```

---

## See Also

- [Spectral Analysis](Spectral-Analysis) - Spectrum file formats and analysis
- [API Documentation](API-Documentation) - API endpoints for data access
- [Database Maintenance](Database-Maintenance) - Post-import optimization
- [BACKFILL_GUIDE.md](/BACKFILL_GUIDE.md) - Detailed backfill instructions
- [RADIACODE_CSV_FORMAT.md](/RADIACODE_CSV_FORMAT.md) - RadiaCode CSV format specification
