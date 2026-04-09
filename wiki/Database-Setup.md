# Database Setup

Complete guide to configuring and optimizing database backends.

[[Home|← Back to Home]]

---

## Overview

Safecast New Map supports multiple database backends to suit different deployment scenarios:

| Database | Best For | Pros | Cons |
|----------|----------|------|------|
| **PostgreSQL** | Production, multi-user | Spatial indexing, scalability, reliability | Requires setup, PostGIS dependency |
| **DuckDB** | Fast local analytics | Extremely fast, no server needed | Single-writer, local only |
| **SQLite** | Simple, single-user | Zero configuration, portable | Limited concurrency |
| **ClickHouse** | Large-scale analytics | Columnar storage, analytics optimized | Complex setup |
| **Chai** | Alternative SQLite | Specialized features | Less common |

---

## PostgreSQL Setup

### Installation

#### Ubuntu/Debian

```bash
# Install PostgreSQL and PostGIS
sudo apt update
sudo apt install postgresql-16 postgresql-16-postgis-3

# Verify installation
psql --version
```

#### RHEL/CentOS/Rocky

```bash
# Install PostgreSQL and PostGIS
sudo dnf install postgresql16 postgresql16-server postgis34_16

# Initialize and start
sudo postgresql-setup --initdb
sudo systemctl enable postgresql
sudo systemctl start postgresql
```

#### macOS (Homebrew)

```bash
brew install postgresql postgis
brew services start postgresql
```

### Create Database

```bash
# Connect as postgres user
sudo -u postgres psql

# Create database and user
CREATE DATABASE safecast;
CREATE USER safecast_user WITH PASSWORD 'your-password';
GRANT ALL PRIVILEGES ON DATABASE safecast TO safecast_user;

# Connect to database
\c safecast

# Enable PostGIS extension
CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;

# Grant schema permissions
GRANT ALL ON SCHEMA public TO safecast_user;
```

### Configure Connection

```bash
./safecast-new-map \
  -db-type pgx \
  -db-conn "postgres://safecast_user:your-password@localhost:5432/safecast?sslmode=disable"
```

**Connection string parameters:**
- `sslmode` - `disable`, `require`, `verify-full`
- `pool_max_conns` - Maximum connections (default: auto-calculated)
- `pool_min_conns` - Minimum connections

### PostgreSQL-Specific Features

**Spatial Indexing:**
- Uses PostGIS for efficient geographic queries
- `ST_DWithin` for radius searches
- `ST_Expand` for bounding box queries
- GIST indexes on geometry columns

**Connection Pooling:**
- Auto-tuned based on CPU cores (4x cores, min 16)
- Connection limits set at startup
- Efficient connection reuse

**Performance Tuning:**
```sql
-- Increase work_mem for large imports
ALTER SYSTEM SET work_mem = '256MB';

-- Increase maintenance_work_mem for index creation
ALTER SYSTEM SET maintenance_work_mem = '1GB';

-- Set shared_buffers (typically 25% of RAM)
ALTER SYSTEM SET shared_buffers = '4GB';

-- Reload configuration
SELECT pg_reload_conf();
```

### Large Batch Imports

For optimal performance during bulk imports:

```sql
-- Temporarily disable autovacuum
ALTER TABLE markers SET (autovacuum_enabled = false);

-- After import, re-enable and vacuum
ALTER TABLE markers SET (autovacuum_enabled = true);
VACUUM ANALYZE markers;
```

---

## DuckDB Setup

### Installation

DuckDB is embedded in the application. No separate installation required.

### Basic Configuration

```bash
./safecast-new-map \
  -db-type duckdb \
  -db-path /path/to/data
```

### Advanced Configuration

**Memory Limits:**
DuckDB automatically detects and uses 20% of system RAM.

**Manual memory limit:**
```bash
export DUCKDB_MEMORY_LIMIT=4GB
./safecast-new-map -db-type duckdb -db-path /path/to/data
```

**Checkpoint Threshold:**
```bash
export DUCKDB_CHECKPOINT_THRESHOLD=100000
```

### DuckLake Analytics

DuckLake combines in-memory DuckDB with PostgreSQL catalog and Parquet data files:

```bash
DATABASE_URL="postgres://user:pass@localhost/safecast" \
DUCKLAKE_PG_URL="dbname=ducklake_catalog host=localhost user=ducklake_rw password=pass" \
DUCKLAKE_DATA_PATH="/var/lib/safecast/ducklake/" \
./safecast-new-map -db-type duckdb -db-path /path/to/data
```

**Shared Analytics Tables:**
- `chat_questions` - AI chat logs
- `mcp_query_log` - MCP tool usage
- `mcp_ai_query_log` - AI query execution

**Benefits:**
- Multiple services share analytics concurrently
- Parquet files for efficient storage
- PostgreSQL catalog for metadata

### Post-Import Maintenance

DuckDB automatically runs maintenance after imports:

```sql
-- Manual maintenance
CHECKPOINT;
OPTIMIZE;
VACUUM;
```

### Single-Writer Mode

DuckDB operates in serialized pipeline mode:
- Single writer at a time
- Automatic queue management
- No concurrent write conflicts

---

## SQLite Setup

### Basic Configuration

```bash
./safecast-new-map \
  -db-type sqlite \
  -db-path /path/to/data.db
```

### WAL Mode

SQLite automatically uses WAL (Write-Ahead Logging) mode for better concurrency:

```sql
-- Verify WAL mode
PRAGMA journal_mode;

-- Should return: wal
```

### Performance Tuning

```sql
-- Increase cache size
PRAGMA cache_size = -2000000;  -- 2GB

-- Increase mmap size for faster reads
PRAGMA mmap_size = 268435456;  -- 256MB

-- Set synchronous mode
PRAGMA synchronous = NORMAL;

-- Enable foreign keys
PRAGMA foreign_keys = ON;
```

### Single-Writer Mode

Like DuckDB, SQLite uses serialized pipeline:
- Automatic write serialization
- Queue-based access
- No concurrent write conflicts

---

## ClickHouse Setup

### Installation

#### Ubuntu/Debian

```bash
# Add ClickHouse repository
curl https://packages.clickhouse.com/rpm/lts/repodata/repomd.xml.key | sudo apt-key add -
sudo apt-add-repository "deb https://packages.clickhouse.com/deb stable main"
sudo apt update

# Install ClickHouse
sudo apt install clickhouse-server clickhouse-client
```

### Create Database

```bash
# Connect to ClickHouse
clickhouse-client

# Create database
CREATE DATABASE safecast;

# Create user
CREATE USER safecast_user IDENTIFIED BY 'your-password';
GRANT ALL ON safecast.* TO safecast_user;
```

### Configure Connection

```bash
./safecast-new-map \
  -db-type clickhouse \
  -db-conn "clickhouse://safecast_user:your-password@localhost:9000/safecast?secure=true"
```

### ClickHouse-Specific Features

**Columnar Storage:**
- Optimized for analytics queries
- Fast aggregations
- Efficient compression

**Connection Pooling:**
- 4x CPU cores, min 16 connections
- Optimized for analytical workloads

---

## Database Schema

### Core Tables

#### `markers`

Core radiation measurement data:

```sql
CREATE TABLE markers (
  id BIGSERIAL PRIMARY KEY,
  doseRate DOUBLE PRECISION,
  date TIMESTAMPTZ,
  lat DOUBLE PRECISION,
  lon DOUBLE PRECISION,
  countRate DOUBLE PRECISION,
  zoom INTEGER,
  speed DOUBLE PRECISION,
  trackID TEXT,
  altitude DOUBLE PRECISION,
  detector TEXT,
  radiation TEXT,
  temperature DOUBLE PRECISION,
  humidity DOUBLE PRECISION,
  deviceID TEXT,
  deviceName TEXT,
  transport TEXT,
  tube TEXT,
  country TEXT,
  liveExtra TEXT,
  hasSpectrum BOOLEAN
);
```

**Indexes:**
- Primary key on `id`
- GIST index on geometry (PostgreSQL)
- Composite index on `(lat, lon)`
- Index on `trackID`
- Index on `date`

#### `tracks`

Lightweight track registry:

```sql
CREATE TABLE tracks (
  trackID TEXT PRIMARY KEY,
  count INTEGER,
  start_date TIMESTAMPTZ,
  end_date TIMESTAMPTZ
);
```

#### `uploads`

File upload metadata:

```sql
CREATE TABLE uploads (
  id BIGSERIAL PRIMARY KEY,
  filename TEXT,
  source TEXT,
  uploaded_at TIMESTAMPTZ,
  recording_date TIMESTAMPTZ,
  detector TEXT,
  username TEXT,
  internal_user_id TEXT,
  name TEXT,
  notes TEXT,
  comment TEXT
);
```

#### `realtime_measurements`

Live sensor readings:

```sql
CREATE TABLE realtime_measurements (
  id BIGSERIAL PRIMARY KEY,
  sensor_id TEXT,
  doseRate DOUBLE PRECISION,
  timestamp TIMESTAMPTZ,
  lat DOUBLE PRECISION,
  lon DOUBLE PRECISION
);
```

#### `spectra`

Gamma spectrum data:

```sql
CREATE TABLE spectra (
  id BIGSERIAL PRIMARY KEY,
  marker_id BIGINT REFERENCES markers(id),
  channels JSON,
  channel_count INTEGER,
  energy_min_kev DOUBLE PRECISION,
  energy_max_kev DOUBLE PRECISION,
  live_time_sec DOUBLE PRECISION,
  real_time_sec DOUBLE PRECISION,
  device_model TEXT,
  calibration JSON,
  source_format TEXT,
  filename TEXT,
  raw_data BYTEA,
  created_at TIMESTAMPTZ
);
```

#### `users`

User accounts:

```sql
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email TEXT UNIQUE,
  password_hash TEXT,
  api_key TEXT UNIQUE,
  username TEXT,
  email_verified BOOLEAN,
  is_active BOOLEAN,
  is_admin BOOLEAN,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  last_login_at TIMESTAMPTZ
);
```

#### `translations`

UI translations (29 languages):

```sql
CREATE TABLE translations (
  id BIGSERIAL PRIMARY KEY,
  language_code TEXT,
  key TEXT,
  value TEXT,
  UNIQUE(language_code, key)
);
```

#### `short_links`

Short link mappings:

```sql
CREATE TABLE short_links (
  id BIGSERIAL PRIMARY KEY,
  code TEXT UNIQUE,
  target TEXT,
  created_at TIMESTAMPTZ
);
```

---

## Migrations

### Running Migrations

Migrations are in the `migrations/` directory:

```bash
# PostgreSQL
psql -d safecast -f migrations/add_comment_to_uploads.sql
psql -d safecast -f migrations/add_detector_index.sql
psql -d safecast -f migrations/add_translations_table.sql

# SQLite
sqlite3 data.db < migrations/add_comment_to_uploads.sql

# DuckDB
duckdb data.duckdb < migrations/add_comment_to_uploads.sql
```

### Available Migrations

| Migration | Purpose |
|-----------|---------|
| `add_comment_to_uploads.sql` | Add comment field to uploads |
| `add_detector_index.sql` | Add index on detector column |
| `add_is_admin_column.sql` | Add admin flag to users |
| `add_recording_date_column.sql` | Add recording date to uploads |
| `add_translations_table.sql` | Create translations table |
| `add_ui_translations.sql` | Seed UI translations |
| `add_upload_metadata.sql` | Add metadata fields to uploads |
| `add_user_id_column.sql` | Add user ID to uploads |
| `add_username_column.sql` | Add username column |
| `fix_markers_sequence.sql` | Fix PostgreSQL sequence |
| `fix_recording_dates.sql` | Fix recording dates |
| `link_historical_uploads_to_users.sql` | Link uploads to users |
| `seed_translations.sql` | Seed translation data |

### Spectrum Support Migration

Add spectrum support to existing database:

```bash
# PostgreSQL
psql -d your_database -f migrations/add_spectrum_support.sql

# SQLite
sqlite3 data.db < migrations/add_spectrum_support_sqlite.sql

# DuckDB
duckdb data.duckdb < migrations/add_spectrum_support_duckdb.sql
```

See [SPECTRAL_MIGRATION_GUIDE.md](/SPECTRAL_MIGRATION_GUIDE.md) for details.

---

## Performance Optimization

### Connection Pooling

Auto-tuned based on database type and CPU cores:

| Database | Max Connections | Min Connections |
|----------|----------------|-----------------|
| PostgreSQL | 4x CPU cores (min 16) | 4 |
| DuckDB | 1 | 1 |
| SQLite | 1 | 1 |
| ClickHouse | 4x CPU cores (min 16) | 4 |

### Index Strategy

**PostgreSQL:**
- GIST spatial indexes for geography
- B-tree indexes for common queries
- Composite indexes for multi-column queries

**SQLite/DuckDB:**
- Primary key indexes
- Secondary indexes on frequently queried columns
- Automatic index optimization

### Background Index Building

PostgreSQL builds indexes in background (non-blocking):
- Uses `CREATE INDEX CONCURRENTLY`
- No write locks during index creation
- Safe for production use

### Materialized Views

Track statistics materialized view:

```sql
-- Refresh materialized view
REFRESH MATERIALIZED VIEW track_stats;

-- Or use script
./tools/refresh_track_stats.sh
```

---

## Backup & Restore

### PostgreSQL Backup

```bash
# Full backup
pg_dump -h localhost -U safecast_user safecast > backup.sql

# Compressed backup
pg_dump -h localhost -U safecast_user safecast | gzip > backup.sql.gz

# Schema only
pg_dump -s -h localhost -U safecast_user safecast > schema.sql

# Data only
pg_dump -a -h localhost -U safecast_user safecast > data.sql
```

### PostgreSQL Restore

```bash
# Restore from backup
psql -h localhost -U safecast_user safecast < backup.sql

# Restore compressed backup
gunzip < backup.sql.gz | psql -h localhost -U safecast_user safecast
```

### SQLite/DuckDB Backup

```bash
# Copy database file
cp data.db data.db.backup

# Or use .backup command (SQLite)
sqlite3 data.db ".backup 'backup.db'"
```

### Automated Backups

**PostgreSQL with cron:**
```bash
# Daily backup at 2 AM
0 2 * * * pg_dump -h localhost -U safecast_user safecast | gzip > /backups/safecast_$(date +\%Y\%m\%d).sql.gz

# Keep 30 days
0 3 * * * find /backups -name "safecast_*.sql.gz" -mtime +30 -delete
```

---

## Monitoring

### PostgreSQL Monitoring

```sql
-- Database size
SELECT pg_size_pretty(pg_database_size('safecast'));

-- Table sizes
SELECT
  relname as table_name,
  pg_size_pretty(pg_total_relation_size(relid)) as total_size
FROM pg_catalog.pg_statio_user_tables
ORDER BY pg_total_relation_size(relid) DESC;

-- Active connections
SELECT count(*) FROM pg_stat_activity WHERE datname = 'safecast';

-- Slow queries
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### Connection Monitoring

```bash
# Monitor connection count
watch -n 1 "psql -c 'SELECT count(*) FROM pg_stat_activity;'"
```

### Performance Metrics

Track key metrics:
- Query execution time
- Connection pool utilization
- Cache hit ratio
- Index usage statistics

---

## Troubleshooting

### Connection Refused

**PostgreSQL:**
```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Check listening ports
sudo netstat -tlnp | grep 5432

# Check pg_hba.conf for authentication
sudo cat /etc/postgresql/16/main/pg_hba.conf
```

### Out of Memory

**PostgreSQL:**
```sql
-- Reduce work_mem
ALTER SYSTEM SET work_mem = '64MB';
SELECT pg_reload_conf();
```

**DuckDB:**
```bash
# Set memory limit
export DUCKDB_MEMORY_LIMIT=2GB
```

### Slow Queries

**Analyze query plan:**
```sql
EXPLAIN ANALYZE SELECT * FROM markers WHERE lat BETWEEN 35.6 AND 35.7;
```

**Update statistics:**
```sql
ANALYZE markers;
ANALYZE tracks;
```

### Sequence Out of Sync

After restore or migration:

```bash
# Fix markers sequence
psql -d safecast -c "SELECT setval('markers_id_seq', (SELECT MAX(id) FROM markers) + 1);"

# Or use tool
go run ./cmd/tools/fix-pg-sequence
```

---

## See Also

- [Database Maintenance](Database-Maintenance) - Maintenance tasks and utilities
- [Getting Started](Getting-Started) - Initial database setup
- [MIGRATION_GUIDE.md](/MIGRATION_GUIDE.md) - Migration between databases
- [FIX_DUPLICATE_KEY.md](/FIX_DUPLICATE_KEY.md) - Fix sequence issues
