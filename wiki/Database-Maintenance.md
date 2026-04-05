# Database Maintenance

Guide to maintaining, optimizing, and troubleshooting your database.

[[Home|← Back to Home]]

---

## Overview

Regular database maintenance ensures optimal performance and data integrity. This guide covers:
- Sequence management
- Index optimization
- Statistics refresh
- Backup and restore
- Common troubleshooting

---

## PostgreSQL Maintenance

### Reset Sequences

After database restore or migration, sequences may be out of sync:

**Symptom:**
```
ERROR: duplicate key value violates unique constraint "markers_pkey"
```

**Fix markers sequence:**
```bash
psql -h localhost -U safecast_user -d safecast -c \
  "SELECT setval('markers_id_seq', (SELECT MAX(id) FROM markers) + 1);"
```

**Fix all sequences:**
```bash
./tools/reset_postgres_sequences.sh
```

**Fix marker sequence only:**
```bash
./tools/reset_postgres_marker_sequence.sh
```

**Using Go tool:**
```bash
export DATABASE_URL="postgres://user:pass@localhost/dbname"
go run ./cmd/tools/fix-pg-sequence
```

### Refresh Statistics

After bulk imports, refresh materialized views:

```bash
./tools/refresh_track_stats.sh
```

**Or manually:**
```sql
REFRESH MATERIALIZED VIEW track_stats;
```

### Vacuum and Analyze

**After large imports:**
```sql
VACUUM ANALYZE markers;
VACUUM ANALYZE tracks;
VACUUM ANALYZE uploads;
```

**Full vacuum (requires exclusive lock):**
```sql
VACUUM FULL markers;
```

**Note:** `VACUUM FULL` locks the table. Use during maintenance windows.

### Update Statistics

```sql
-- Update statistics for query planner
ANALYZE markers;
ANALYZE tracks;
ANALYZE uploads;
ANALYZE users;
```

### Index Maintenance

**Check index usage:**
```sql
SELECT
  schemaname,
  tablename,
  indexname,
  idx_scan,
  idx_tup_read,
  idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan ASC;
```

**Rebuild index:**
```sql
REINDEX INDEX markers_lat_lon_idx;

-- Or reindex entire table
REINDEX TABLE markers;
```

**Create missing indexes:**
```sql
-- Index on frequently queried columns
CREATE INDEX CONCURRENTLY idx_markers_date ON markers(date);
CREATE INDEX CONCURRENTLY idx_markers_trackid ON markers(trackID);
CREATE INDEX CONCURRENTLY idx_markers_detector ON markers(detector);
```

### Check Database Size

```sql
-- Database size
SELECT pg_size_pretty(pg_database_size('safecast'));

-- Table sizes
SELECT
  relname as table_name,
  pg_size_pretty(pg_total_relation_size(relid)) as total_size,
  pg_size_pretty(pg_relation_size(relid)) as data_size,
  pg_size_pretty(pg_total_relation_size(relid) - pg_relation_size(relid)) as index_size
FROM pg_catalog.pg_statio_user_tables
ORDER BY pg_total_relation_size(relid) DESC;

-- Index sizes
SELECT
  indexname,
  pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
ORDER BY pg_relation_size(indexrelid) DESC;
```

### Check Connections

```sql
-- Active connections
SELECT count(*) FROM pg_stat_activity WHERE datname = 'safecast';

-- Connections by state
SELECT state, count(*)
FROM pg_stat_activity
WHERE datname = 'safecast'
GROUP BY state;

-- Long-running queries
SELECT
  pid,
  now() - pg_stat_activity.query_start AS duration,
  query,
  state
FROM pg_stat_activity
WHERE datname = 'safecast'
  AND now() - pg_stat_activity.query_start > interval '5 minutes'
ORDER BY duration DESC;
```

### Kill Long-Running Queries

```sql
-- Terminate query by PID
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = 'safecast'
  AND now() - pg_stat_activity.query_start > interval '30 minutes';
```

---

## SQLite Maintenance

### Vacuum Database

```bash
sqlite3 data.db "VACUUM;"
```

### Analyze Tables

```bash
sqlite3 data.db "ANALYZE;"
```

### Check Database Integrity

```bash
sqlite3 data.db "PRAGMA integrity_check;"
```

### Optimize WAL Mode

```sql
-- Check journal mode
PRAGMA journal_mode;

-- Should return: wal

-- Set WAL mode if needed
PRAGMA journal_mode = WAL;

-- Set checkpoint threshold
PRAGMA wal_autocheckpoint = 1000;
```

### Manual Checkpoint

```sql
-- Force checkpoint
PRAGMA wal_checkpoint(TRUNCATE);
```

---

## DuckDB Maintenance

### Optimize Database

```sql
-- Optimize tables
OPTIMIZE;

-- Checkpoint WAL
CHECKPOINT;

-- Vacuum
VACUUM;
```

### Check Memory Usage

```sql
-- Check memory limit
SELECT current_setting('memory_limit');

-- Set memory limit
SET memory_limit = '4GB';
```

### Check Database Size

```sql
-- Check database file size
-- (Check file system)
ls -lh data.duckdb
```

---

## ClickHouse Maintenance

### Optimize Tables

```sql
-- Optimize table (merge parts)
OPTIMIZE TABLE markers FINAL;

-- Check table size
SELECT
  table,
  formatReadableSize(sum(bytes)) as size,
  sum(rows) as rows
FROM system.parts
WHERE database = 'safecast'
  AND active
GROUP BY table;
```

### Check Partitions

```sql
-- View partitions
SELECT
  partition,
  sum(rows) as rows,
  count() as parts
FROM system.parts
WHERE database = 'safecast'
  AND table = 'markers'
  AND active
GROUP BY partition
ORDER BY partition;
```

---

## Backup Strategies

### PostgreSQL Backup

**Full backup:**
```bash
pg_dump -h localhost -U safecast_user safecast > backup.sql
```

**Compressed backup:**
```bash
pg_dump -h localhost -U safecast_user safecast | gzip > backup.sql.gz
```

**Schema only:**
```bash
pg_dump -s -h localhost -U safecast_user safecast > schema.sql
```

**Data only:**
```bash
pg_dump -a -h localhost -U safecast_user safecast > data.sql
```

**Custom format (parallel):**
```bash
pg_dump -h localhost -U safecast_user -F c -j 4 safecast > backup.dump
```

**Restore custom format:**
```bash
pg_restore -h localhost -U safecast_user -d safecast backup.dump
```

### Automated PostgreSQL Backup

**Backup script:**
```bash
#!/bin/bash
# /opt/safecast/backup.sh

BACKUP_DIR="/backups/safecast"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/safecast_$DATE.sql.gz"
RETENTION_DAYS=30

mkdir -p $BACKUP_DIR

# Backup
pg_dump -h localhost -U safecast_user safecast | gzip > $BACKUP_FILE

# Verify backup
if [ -f "$BACKUP_FILE" ]; then
  echo "Backup successful: $BACKUP_FILE"
else
  echo "Backup failed!"
  exit 1
fi

# Cleanup old backups
find $BACKUP_DIR -name "safecast_*.sql.gz" -mtime +$RETENTION_DAYS -delete

echo "Cleaned up backups older than $RETENTION_DAYS days"
```

**Cron job:**
```bash
# Daily backup at 2 AM
0 2 * * * /opt/safecast/backup.sh >> /var/log/safecast-backup.log 2>&1
```

### SQLite/DuckDB Backup

```bash
# Copy database file
cp data.db data.db.backup.$(date +%Y%m%d)

# Or use .backup command
sqlite3 data.db ".backup 'data.backup.db'"
```

### Automated SQLite Backup

```bash
#!/bin/bash
# /opt/safecast/backup-sqlite.sh

DATA_DIR="/data"
BACKUP_DIR="/backups/safecast"
DATE=$(date +%Y%m%d)
RETENTION_DAYS=30

mkdir -p $BACKUP_DIR

# Backup
cp $DATA_DIR/data.db $BACKUP_DIR/safecast_$DATE.db

# Verify
if [ -f "$BACKUP_DIR/safecast_$DATE.db" ]; then
  echo "Backup successful"
else
  echo "Backup failed!"
  exit 1
fi

# Cleanup
find $BACKUP_DIR -name "safecast_*.db" -mtime +$RETENTION_DAYS -delete
```

---

## Migration Tools

### Fix PostgreSQL Sequence

```bash
export DATABASE_URL="postgres://user:pass@localhost/dbname"
go run ./cmd/tools/fix-pg-sequence
```

### Alternative Sequence Fix

```bash
go run ./cmd/tools/fix-sequence
```

### Add Internal User IDs

```bash
go run ./cmd/tools/add-internal-user-id
```

Links historical uploads to users.

### Cleanup Test Users

```bash
go run ./cmd/tools/cleanup-test-users
```

Removes test users from database.

### Import API Keys

```bash
go run ./cmd/tools/import-api-keys -file api_keys.csv
```

**CSV format:**
```csv
user_id,api_key
1,abcdefghijklmnopqrst
2,uvwxyz1234567890abcd
```

### Migrate to PostgreSQL

```bash
go run ./cmd/tools/migrate-to-postgres \
  -source /path/to/data.db \
  -dest "postgres://user:pass@localhost/dbname"
```

### Migrate Users

```bash
go run ./cmd/tools/migrate-users
```

---

## Populate Usernames

Fetch usernames from Safecast API for historical uploads:

```bash
./tools/populate_usernames.sh
```

**Manual SQL:**
```sql
UPDATE uploads u
SET username = (
  SELECT username FROM users
  WHERE id = u.internal_user_id
)
WHERE u.username IS NULL
  AND u.internal_user_id IS NOT NULL;
```

---

## Fix Recording Dates

Correct recording dates from uploads:

**PostgreSQL:**
```bash
psql -d safecast -f migrations/fix_recording_dates.sql
```

**SQLite:**
```bash
sqlite3 data.db < migrations/fix_recording_dates_sqlite.sql
```

---

## Performance Monitoring

### Slow Query Log

**PostgreSQL:**
```sql
-- Enable slow query log
ALTER SYSTEM SET log_min_duration_statement = 1000;  -- 1 second
SELECT pg_reload_conf();

-- View recent slow queries
SELECT query, calls, mean_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### Query Plan Analysis

```sql
-- Analyze query plan
EXPLAIN ANALYZE
SELECT * FROM markers
WHERE lat BETWEEN 35.6 AND 35.7
  AND lon BETWEEN 139.6 AND 139.7
ORDER BY date DESC
LIMIT 100;
```

**Look for:**
- Sequential scans (bad for large tables)
- Missing indexes
- High cost operations

### Cache Hit Ratio

```sql
-- Check cache hit ratio
SELECT
  sum(blks_hit) * 100.0 / sum(blks_hit + blks_read) as cache_hit_ratio
FROM pg_stat_database
WHERE datname = 'safecast';
```

**Good ratio:** > 95%

**Improve cache:**
```sql
-- Increase shared_buffers
ALTER SYSTEM SET shared_buffers = '4GB';
SELECT pg_reload_conf();
```

### Index Hit Ratio

```sql
-- Check index usage
SELECT
  relname,
  idx_scan,
  idx_tup_read,
  idx_tup_fetch
FROM pg_stat_user_tables
ORDER BY idx_scan ASC;
```

Tables with low `idx_scan` may need indexes.

---

## Troubleshooting

### Database Connection Failed

**Check if PostgreSQL is running:**
```bash
sudo systemctl status postgresql
```

**Check listening port:**
```bash
sudo netstat -tlnp | grep 5432
```

**Test connection:**
```bash
psql -h localhost -U safecast_user -d safecast -c "SELECT 1;"
```

**Check pg_hba.conf:**
```bash
sudo cat /etc/postgresql/16/main/pg_hba.conf
```

### Out of Disk Space

**Check disk usage:**
```bash
df -h
du -sh /var/lib/postgresql/16/main/
```

**Clean up:**
```bash
# Remove old backups
find /backups -name "*.sql.gz" -mtime +30 -delete

# Vacuum database
VACUUM FULL markers;
```

### High Memory Usage

**PostgreSQL:**
```sql
-- Reduce work_mem
ALTER SYSTEM SET work_mem = '64MB';
SELECT pg_reload_conf();

-- Reduce shared_buffers
ALTER SYSTEM SET shared_buffers = '2GB';
SELECT pg_reload_conf();
```

**DuckDB:**
```bash
export DUCKDB_MEMORY_LIMIT=2GB
```

### Lock Contention

**Check locks:**
```sql
SELECT
  blocked_locks.pid AS blocked_pid,
  blocked_activity.query AS blocked_query,
  blocking_locks.pid AS blocking_pid,
  blocking_activity.query AS blocking_query
FROM pg_catalog.pg_locks blocked_locks
JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted;
```

**Kill blocking query:**
```sql
SELECT pg_terminate_backend(blocking_pid);
```

### Replication Lag

**Check replication status:**
```sql
-- On primary
SELECT client_addr, state, sent_lsn, write_lsn, flush_lsn, replay_lsn
FROM pg_stat_replication;

-- Calculate lag
SELECT
  client_addr,
  pg_wal_lsn_diff(sent_lsn, replay_lsn) as replication_lag_bytes
FROM pg_stat_replication;
```

### Corruption Check

**PostgreSQL:**
```sql
-- Check table integrity
SELECT * FROM markers LIMIT 1;

-- Check indexes
REINDEX TABLE markers;
```

**SQLite:**
```bash
sqlite3 data.db "PRAGMA integrity_check;"
```

---

## Scheduled Maintenance

### Daily Tasks

- Monitor disk space
- Check error logs
- Verify backups completed

### Weekly Tasks

- Analyze tables
- Check index usage
- Review slow queries

### Monthly Tasks

- Full vacuum (during maintenance window)
- Review and cleanup old backups
- Update statistics
- Check replication health

### Quarterly Tasks

- Review and optimize queries
- Update PostgreSQL version
- Review and adjust configuration
- Disaster recovery test

---

## See Also

- [Database Setup](Database-Setup) - Database configuration
- [FIX_DUPLICATE_KEY.md](/FIX_DUPLICATE_KEY.md) - Fix sequence issues
- [MIGRATION_GUIDE.md](/MIGRATION_GUIDE.md) - Migration procedures
- [Deployment](Deployment) - Backup strategies
- [Configuration Reference](Configuration-Reference) - Database settings
