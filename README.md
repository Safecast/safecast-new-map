# Safecast New Map

[![Build Status](https://github.com/Safecast/safecast-new-map/actions/workflows/release.yml/badge.svg)](https://github.com/Safecast/safecast-new-map/actions/workflows/release.yml)

A modern, self-hosted radiation monitoring platform that makes environmental data accessible to everyone. Built to help communities understand radiation levels in their environment through open data and transparent mapping.

**Live Demo:** [simplemap.safecast.org](https://simplemap.safecast.org/)

**Downloads:** [Latest releases for all platforms](https://github.com/safecast/safecast-new-map/releases)

## Language Support

[English](/README.md) | [Français](/doc/README_FR.md) | [日本語](/doc/README_JP.md) | [Deutsch (CH)](/doc/README_DE_CH.md) | [Italiano](/doc/README_IT.md) | [中文](/doc/README_ZH.md) | [हिन्दी](/doc/README_HI.md) | [فارسی](/doc/README_FA.md) | [Русский](/doc/README_RU.md) | [Монгол](/doc/README_MN.md) | [Қазақша](/doc/README_KK.md)

---

## About

This project provides a complete radiation mapping solution that runs on your own infrastructure. Whether you're monitoring environmental safety, conducting research, or providing public information, the platform handles everything from data collection to visualization.

Natural background radiation is typically low and safe. This map helps identify areas where levels rise above normal due to contamination, natural deposits, or other factors. Understanding these patterns helps communities make informed decisions about drinking water sources, agriculture, and land use.

### Key Features

**Data Collection & Import**
- Upload radiation measurements from multiple device formats (bGeigie, RadiaCode, AtomFast, and more)
- Automatic sync with Safecast's global database
- Import from files (.kml, .kmz, .json, .rctrk, .csv, .gpx) or URLs
- Real-time device monitoring (optional)

**Visualization & Analysis**
- Interactive map with multiple coloring schemes (scientific gradient or safety-focused)
- Speed-based layer separation (walking, driving, flying)
- Time-series analysis for specific locations
- Country-level statistics and reporting
- Print mode with QR codes for field marking

**Data Management**
- Multiple database backends (PostgreSQL, DuckDB, SQLite, ClickHouse)
- Automated JSON archive generation (daily/weekly/monthly)
- Track streaming for efficient large dataset handling
- RESTful API with rate limiting

**Advanced Capabilities**
- Gamma spectrum analysis (.spe, .n42 formats)
- Admin panel for content moderation
- Short link generation for sharing
- Multi-language interface
- Auto-update system

---

## Quick Start

### Option 1: Binary (Recommended)

Download and run in seconds:

```bash
# Download from https://github.com/safecast/safecast-new-map/releases
chmod +x ./safecast-new-map
./safecast-new-map
```

Open [http://localhost:8765](http://localhost:8765)

### Option 2: With Production Data

Start with a complete dataset from simplemap.safecast.org:

```bash
./safecast-new-map -import-tgz-url https://simplemap.safecast.org/api/json/weekly.tgz
```

This imports all public tracks and starts the server. Future runs will use the cached data.

### Option 3: Production Deployment

Deploy with HTTPS using Let's Encrypt:

```bash
./safecast-new-map -domain maps.example.org -db-type pgx -db-conn "postgres://user:pass@localhost/safecast"
```

Requires ports 80 and 443 open for certificate validation.

### Option 4: Docker

```bash
docker run -d -p 8765:8765 --name safecast-map safecastr/safecast-new-map:latest
```

Open [http://localhost:8765](http://localhost:8765)

---

## Configuration

### Database Options

**PostgreSQL (Recommended for Production)**
```bash
./safecast-new-map -db-type pgx -db-conn "postgres://user:pass@host:5432/dbname?sslmode=require"
```

**DuckDB (Fast Local Storage)**
```bash
./safecast-new-map -db-type duckdb -db-path /path/to/data
```

**SQLite (Simple Single-User)**
```bash
./safecast-new-map -db-type sqlite -db-path /path/to/data
```

**ClickHouse (Large-Scale Analytics)**
```bash
./safecast-new-map -db-type clickhouse -db-conn "clickhouse://user:pass@host:9000/dbname?secure=true"
```

### Common Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | 8765 | HTTP server port |
| `-domain` | - | Domain for HTTPS (enables Let's Encrypt) |
| `-db-type` | pgx | Database: pgx, duckdb, sqlite, chai, clickhouse |
| `-db-path` | . | Path for file-based databases |
| `-db-conn` | - | Connection string for network databases |
| `-default-lat` | 44.08832 | Initial map latitude |
| `-default-lon` | 42.97577 | Initial map longitude |
| `-default-zoom` | 11 | Initial map zoom level |
| `-default-layer` | OpenStreetMap | Base map layer |
| `-admin-password` | - | Enable admin panel (track management) |
| `-safecast-realtime` | false | Poll live Safecast device data |
| `-safecast-fetcher` | false | Auto-sync approved bGeigie imports |
| `-json-archive-frequency` | weekly | Archive generation: daily, weekly, monthly, yearly |

---

## Data Import & Export

### Supported Formats

**Import:**
- KML/KMZ (Google Earth, Safecast bGeigie)
- JSON (exported tracks)
- RCTRK (RadiaCode)
- CSV (AtomFast, custom formats)
- GPX (GPS tracks)
- LOG (bGeigie Nano/Zen)

**Export:**
- JSON archives (compressed .tgz)
- Individual track JSON
- Legacy CIM format

### Bulk Import

Import from remote archive:
```bash
./safecast-new-map -import-tgz-url https://simplemap.safecast.org/api/json/weekly.tgz
```

Import from local file:
```bash
./safecast-new-map -import-tgz-file /path/to/archive.tgz
```

### API Endpoints

- **Track data:** `/api/track/{id}.json`
- **Archives:** `/api/json/weekly.tgz` (also daily, monthly, yearly)
- **Track list:** `/api/tracks`
- **Statistics:** `/api/stats`
- **Countries:** `/api/countries`

---

## Spectral Data Analysis

Analyze gamma spectra for isotope identification:

### Supported Formats
- `.spe` - Maestro spectrum format
- `.n42` - ANSI N42.42 standard
- `.rctrk` - RadiaCode with embedded spectra

### Database Migration

Add spectrum support to existing databases:

```bash
# PostgreSQL
psql -d your_database -f migrations/add_spectrum_support.sql

# SQLite
sqlite3 data.db < migrations/add_spectrum_support_sqlite.sql

# DuckDB
duckdb data.duckdb < migrations/add_spectrum_support_duckdb.sql
```

See [SPECTRAL_MIGRATION_GUIDE.md](SPECTRAL_MIGRATION_GUIDE.md) for details.

### Features
- Automatic spectrum extraction during upload
- Energy calibration support
- Peak detection and isotope identification
- Visualization on map markers
- API access for external analysis

---

## Administration

Enable the admin panel with:
```bash
./safecast-new-map -admin-password your-secure-password
```

**Admin capabilities:**
- View all uploads and tracks
- Delete inappropriate content
- Monitor system statistics
- Manage user contributions

---

## URL Parameters

Customize map views with URL parameters:

| Parameter | Values | Description |
|-----------|--------|-------------|
| `coloring` | safecast, chicha | Scientific gradient vs. safety bins |
| `unit` | uSv, uR | Display units (microsieverts or microroentgen) |
| `legend` | 1, 0 | Show/hide legend |
| `lang` | en, ru, ja, etc. | Interface language |
| `layer` | OpenStreetMap, Google Satellite | Base map |

**Examples:**
- Safety view: `/?coloring=chicha&unit=uR`
- Clean embed: `/?legend=0`
- Russian interface: `/?lang=ru`

---

## Automated Data Sync

### Safecast Realtime Devices

Poll live sensor data:
```bash
./safecast-new-map -safecast-realtime
```

Sensors appear on the map in real-time with current readings.

### Safecast API Fetcher

Automatically import approved bGeigie measurements:
```bash
./safecast-new-map -safecast-fetcher \
  -safecast-fetcher-interval 5m \
  -safecast-fetcher-batch-size 10 \
  -safecast-fetcher-start-date 2024-01-01
```

---

## Development

### Build from Source

```bash
git clone https://github.com/Safecast/safecast-new-map.git
cd safecast-new-map
go build -o safecast-new-map
./safecast-new-map
```

### Run Tests

```bash
go test ./...
```

### Cross-Compile

```bash
# See scripts/crosscompile/crosscompile.go
go run scripts/crosscompile/crosscompile.go
```

---

## Performance Notes

**DuckDB Performance:** For large imports, see [doc/DUCKDB_PERFORMANCE.md](doc/DUCKDB_PERFORMANCE.md) for checkpoint and Parquet optimization.

**PostgreSQL Tuning:** Use connection pooling and appropriate `work_mem` settings for large batch imports.

**Rate Limiting:** API endpoints are rate-limited by default. Adjust in code if needed for high-traffic deployments.

---

## Community & Support

This project is developed and maintained by the Safecast community with contributions from:
- Rob Oudendijk (primary developer)
- Safecast volunteers worldwide
- AtomFast community
- RadiaCode community
- Open dosimetry researchers

**Contributing:** We welcome contributions! See issues for areas needing help.

**History:** Inspired by Dmitry Ignatenko's field research and built on Safecast's decade of open radiation monitoring. The goal remains simple: make environmental radiation data accessible to everyone who needs it.

---

## License

Code: [Apache 2.0](LICENSE)  
Data: [CC0 1.0 Universal](LICENSE.CC0)

The map data is provided as-is for public benefit. While we strive for accuracy, always consult professional advice for health and safety decisions.
