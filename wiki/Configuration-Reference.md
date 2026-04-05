# Configuration Reference

Complete reference for all command-line flags and environment variables.

[[Home|← Back to Home]]

---

## Overview

This page provides a comprehensive reference for all configuration options available in Safecast New Map.

---

## Command-Line Flags

### Server Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | 8765 | HTTP server port |
| `-domain` | - | Domain name for HTTPS (enables Let's Encrypt) |
| `-base-url` | - | Base URL for email links (e.g., https://example.com) |

**Example:**
```bash
./safecast-new-map -port 9000 -domain maps.example.org -base-url "https://maps.example.org"
```

### Database Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-db-type` | pgx | Database driver: pgx, duckdb, sqlite, chai, clickhouse |
| `-db-path` | . | Path for file-based databases (SQLite, DuckDB) |
| `-db-conn` | - | Connection string for network databases (PostgreSQL, ClickHouse) |

**Examples:**
```bash
# PostgreSQL
./safecast-new-map -db-type pgx -db-conn "postgres://user:pass@localhost/dbname"

# SQLite
./safecast-new-map -db-type sqlite -db-path /path/to/data.db

# DuckDB
./safecast-new-map -db-type duckdb -db-path /path/to/data

# ClickHouse
./safecast-new-map -db-type clickhouse -db-conn "clickhouse://user:pass@localhost:9000/dbname"
```

### Map Defaults

| Flag | Default | Description |
|------|---------|-------------|
| `-default-lat` | 44.08832 | Initial map latitude |
| `-default-lon` | 42.97577 | Initial map longitude |
| `-default-zoom` | 11 | Initial map zoom level |
| `-default-layer` | OpenStreetMap | Default base map layer |

**Example:**
```bash
./safecast-new-map \
  -default-lat 35.6762 \
  -default-lon 139.6503 \
  -default-zoom 13 \
  -default-layer "Google Satellite"
```

### Authentication

| Flag | Default | Description |
|------|---------|-------------|
| `-admin-password` | - | Enable admin panel with password |
| `-allow-registration` | false | Enable user registration |
| `-require-auth` | false | Require authentication for uploads |
| `-session-secret` | - | Secret key for session encryption (required for auth) |

**Example:**
```bash
./safecast-new-map \
  -admin-password "your-secure-password" \
  -allow-registration \
  -require-auth \
  -session-secret "your-random-32-char-secret-key"
```

### Email Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-smtp-host` | - | SMTP server hostname (e.g., smtp.gmail.com) |
| `-smtp-port` | 587 | SMTP server port |
| `-smtp-username` | - | SMTP authentication username |
| `-smtp-password` | - | SMTP authentication password |
| `-smtp-from` | - | Email "From" address |

**Example:**
```bash
./safecast-new-map \
  -smtp-host smtp.gmail.com \
  -smtp-port 587 \
  -smtp-username your-email@gmail.com \
  -smtp-password your-app-password \
  -smtp-from your-email@gmail.com
```

### Data Sync

| Flag | Default | Description |
|------|---------|-------------|
| `-safecast-realtime` | false | Poll live Safecast device data |
| `-safecast-fetcher` | false | Auto-sync approved bGeigie imports from Safecast API |
| `-safecast-fetcher-interval` | 5m | Polling interval for fetcher |
| `-safecast-fetcher-batch-size` | 10 | Records per API call |
| `-safecast-fetcher-start-date` | - | Start date for fetcher (YYYY-MM-DD) |
| `-safecast-fetcher-backfill` | false | Import ALL historical records (ignores DB state) |
| `-safecast-fetcher-newest-first` | false | Fetch newest records first |

**Example:**
```bash
./safecast-new-map \
  -safecast-realtime \
  -safecast-fetcher \
  -safecast-fetcher-interval 10m \
  -safecast-fetcher-batch-size 20 \
  -safecast-fetcher-start-date 2024-01-01
```

### Data Import

| Flag | Default | Description |
|------|---------|-------------|
| `-import-tgz-url` | - | Download and import remote .tgz archive |
| `-import-tgz-file` | - | Import local .tgz archive |

**Example:**
```bash
# Import from URL (runs once then starts server)
./safecast-new-map -import-tgz-url https://simplemap.safecast.org/api/json/weekly.tgz

# Import from local file
./safecast-new-map -import-tgz-file /path/to/archive.tgz
```

### Data Export

| Flag | Default | Description |
|------|---------|-------------|
| `-json-archive-frequency` | weekly | Archive generation: daily, weekly, monthly, yearly |

**Example:**
```bash
./safecast-new-map -json-archive-frequency daily
```

### Self-Upgrade

| Flag | Default | Description |
|------|---------|-------------|
| `-selfupgrade` | false | Enable automatic self-upgrade from GitHub releases |

**Example:**
```bash
./safecast-new-map -selfupgrade -domain maps.example.org
```

---

## Environment Variables

### Database

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | Database connection string (alternative to -db-conn) |

**Example:**
```bash
export DATABASE_URL="postgres://user:pass@localhost/safecast"
./safecast-new-map
```

### DuckLake Analytics

| Variable | Description |
|----------|-------------|
| `DUCKLAKE_PG_URL` | PostgreSQL connection string for DuckLake catalog |
| `DUCKLAKE_DATA_PATH` | Path for DuckLake Parquet data files |

**Example:**
```bash
export DUCKLAKE_PG_URL="dbname=ducklake_catalog host=localhost user=ducklake_rw password=pass"
export DUCKLAKE_DATA_PATH="/var/lib/safecast/ducklake/"
./safecast-new-map
```

### DuckDB

| Variable | Default | Description |
|----------|---------|-------------|
| `DUCKDB_MEMORY_LIMIT` | auto (20% of RAM) | Memory limit for DuckDB |
| `DUCKDB_CHECKPOINT_THRESHOLD` | 100000 | Checkpoint threshold |

**Example:**
```bash
export DUCKDB_MEMORY_LIMIT=4GB
export DUCKDB_CHECKPOINT_THRESHOLD=50000
./safecast-new-map -db-type duckdb
```

### AI Integration

| Variable | Description |
|----------|-------------|
| `ANTHROPIC_API_KEY` | Anthropic API key for AI web chat |
| `CLAUDE_MODEL` | Claude model to use (e.g., claude-haiku-4-5-20251001) |

**Example:**
```bash
export ANTHROPIC_API_KEY="sk-ant-..."
export CLAUDE_MODEL="claude-haiku-4-5-20251001"
./safecast-new-map
```

### Logging

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | info | Log level: debug, info, warn, error |

**Example:**
```bash
export LOG_LEVEL=debug
./safecast-new-map
```

---

## Configuration Examples

### Minimal Development

```bash
./safecast-new-map
```

Uses SQLite with default settings.

### Local Development with PostgreSQL

```bash
./safecast-new-map \
  -db-type pgx \
  -db-conn "postgres://localhost/safecast_dev" \
  -admin-password "dev-password" \
  -allow-registration
```

### Production with HTTPS

```bash
./safecast-new-map \
  -domain maps.example.org \
  -db-type pgx \
  -db-conn "postgres://user:pass@localhost/safecast" \
  -admin-password "production-admin-password" \
  -allow-registration \
  -require-auth \
  -session-secret "$(openssl rand -base64 32)" \
  -smtp-host smtp.gmail.com \
  -smtp-port 587 \
  -smtp-username noreply@example.org \
  -smtp-password "smtp-password" \
  -smtp-from noreply@example.org \
  -base-url "https://maps.example.org"
```

### Production with Docker

```bash
docker run -d \
  -p 8765:8765 \
  -e DATABASE_URL="postgres://user:pass@postgres:5432/safecast" \
  -e SESSION_SECRET="random-secret-key" \
  -e ANTHROPIC_API_KEY="sk-ant-..." \
  -v /data:/data \
  safecastr/safecast-new-map:latest \
  -db-type pgx \
  -admin-password "admin-password" \
  -allow-registration
```

### MCP Server with DuckLake

```bash
DATABASE_URL="postgres://user:pass@localhost/safecast" \
DUCKLAKE_PG_URL="dbname=ducklake_catalog host=localhost user=ducklake_rw" \
DUCKLAKE_DATA_PATH="/var/lib/safecast/ducklake/" \
ANTHROPIC_API_KEY="sk-ant-..." \
./safecast-new-map \
  -db-type duckdb \
  -db-path /data/safecast.duckdb \
  -safecast-realtime
```

### Safecast Data Sync

```bash
./safecast-new-map \
  -db-type pgx \
  -db-conn "postgres://user:pass@localhost/safecast" \
  -safecast-fetcher \
  -safecast-fetcher-interval 5m \
  -safecast-fetcher-batch-size 10 \
  -safecast-fetcher-start-date 2024-01-01 \
  -safecast-realtime
```

### Auto-Upgrade Production

```bash
./safecast-new-map \
  -domain maps.example.org \
  -db-type pgx \
  -db-conn "postgres://user:pass@localhost/safecast" \
  -selfupgrade \
  -admin-password "admin-password"
```

---

## Flag Precedence

When multiple configuration methods are used:

1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **Default values** (lowest priority)

---

## Database Connection Strings

### PostgreSQL

```
postgres://user:password@host:port/database?sslmode=require
```

**Parameters:**
- `sslmode` - `disable`, `require`, `verify-full`
- `pool_max_conns` - Max connections (default: auto)
- `pool_min_conns` - Min connections

**Examples:**
```
postgres://safecast:password@localhost:5432/safecast?sslmode=disable
postgres://safecast:password@db.example.org:5432/safecast?sslmode=require&pool_max_conns=50
```

### ClickHouse

```
clickhouse://user:password@host:port/database?secure=true
```

**Parameters:**
- `secure` - Use TLS
- `timeout` - Connection timeout

**Examples:**
```
clickhouse://safecast:password@localhost:9000/safecast?secure=true
clickhouse://safecast:password@clickhouse.example.org:9440/safecast?secure=true&timeout=30s
```

---

## Configuration Validation

### Check Configuration

```bash
# Show help with all flags
./safecast-new-map --help

# Verify database connection
./safecast-new-map -db-type pgx -db-conn "postgres://..." api/stats
```

### Common Errors

**Missing session secret:**
```
Error: session secret is required when authentication is enabled
```

**Fix:**
```bash
./safecast-new-map -session-secret "your-random-secret"
```

**Invalid database connection:**
```
Error: failed to connect to database: connection refused
```

**Fix:**
```bash
# Verify database is running
psql -h localhost -U user -d dbname -c "SELECT 1;"
```

**Port already in use:**
```
Error: listen tcp :8765: bind: address already in use
```

**Fix:**
```bash
./safecast-new-map -port 9000
```

---

## See Also

- [Getting Started](Getting-Started) - Quick start guide
- [Database Setup](Database-Setup) - Database configuration details
- [Deployment](Deployment) - Production deployment
- [Development](Development) - Development configuration
