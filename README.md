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
- User authentication with API key support
- Admin panel for content moderation and translation management
- Comprehensive authentication logging
- Short link generation for sharing
- Multi-language interface (29 languages) with PostgreSQL-backed translations and admin UI
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

**Note for CloudFront/CDN deployments:** If your domain uses CloudFront or another CDN for web traffic, SSH/rsync deployment must use the server's IP address directly, not the domain name. See [CloudFront Setup Guide](docs/cloudfront-setup.md) for details.

### Option 4: Docker

```bash
docker run -d -p 8765:8765 --name safecast-map safecastr/safecast-new-map:latest
```

Open [http://localhost:8765](http://localhost:8765)

---

## Configuration

### Database Options

**PostgreSQL (Recommended for Production)**

Requires PostGIS extension for spatial indexing:
```bash
# Ubuntu/Debian
sudo apt install postgresql-16-postgis-3

# RHEL/CentOS/Rocky
sudo dnf install postgis34_16
```

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

| Flag                      | Default       | Description                                          |
| ------------------------- | ------------- | ---------------------------------------------------- |
| `-port`                   | 8765          | HTTP server port                                     |
| `-domain`                 | -             | Domain for HTTPS (enables Let's Encrypt)             |
| `-db-type`                | pgx           | Database: pgx, duckdb, sqlite, chai, clickhouse      |
| `-db-path`                | .             | Path for file-based databases                        |
| `-db-conn`                | -             | Connection string for network databases              |
| `-default-lat`            | 44.08832      | Initial map latitude                                 |
| `-default-lon`            | 42.97577      | Initial map longitude                                |
| `-default-zoom`           | 11            | Initial map zoom level                               |
| `-default-layer`          | OpenStreetMap | Base map layer                                       |
| `-admin-password`         | -             | Enable admin panel (track management)                |
| `-allow-registration`     | false         | Enable user registration                             |
| `-require-auth`           | false         | Require authentication for uploads                   |
| `-smtp-host`              | -             | SMTP server for email (e.g., smtp.gmail.com)         |
| `-smtp-port`              | 587           | SMTP server port                                     |
| `-smtp-username`          | -             | SMTP authentication username                         |
| `-smtp-password`          | -             | SMTP authentication password                         |
| `-smtp-from`              | -             | Email "From" address                                 |
| `-session-secret`         | -             | Secret key for session encryption                    |
| `-base-url`               | -             | Base URL for email links (e.g., https://example.com) |
| `-safecast-realtime`      | false         | Poll live Safecast device data                       |
| `-safecast-fetcher`       | false         | Auto-sync approved bGeigie imports                   |
| `-json-archive-frequency` | weekly        | Archive generation: daily, weekly, monthly, yearly   |

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
- **Canonical docs (main app):** `/swagger/` (generated from Swaggo annotations)

### MCP Server & AI Integration

The unified server includes an MCP (Model Context Protocol) server with a Claude-powered web chat interface, all in a single binary.

| Component                  | Source                | Port | URL                             |
| -------------------------- | --------------------- | ---- | ------------------------------- |
| Unified Server (Map + MCP) | `cmd/unified-server/` | 8765 | `/`, `/mcp-http`, `/assistant/` |
| MCP Server                 | `cmd/unified-server/` | 8765 | `/mcp-http`, `/mcp/sse`         |
| Web Chat                   | `cmd/unified-server/` | 8765 | `/assistant/`                   |
| REST API + Swagger         | `cmd/unified-server/` | 8765 | `/api/`, `/docs/`               |

**Build:**
```bash
# Basic build (PostgreSQL only)
go build -o safecast-new-map ./cmd/unified-server

# With DuckDB analytics via DuckLake (requires CGO)
CGO_ENABLED=1 go build -tags duckdb -o safecast-new-map ./cmd/unified-server
```

**Run:**
```bash
# Map server with MCP (PostgreSQL required)
DATABASE_URL="postgres://user:pass@localhost/db" ./safecast-new-map

# With DuckLake analytics (in-memory DuckDB + PostgreSQL catalog + Parquet data)
DATABASE_URL="postgres://..." \
DUCKLAKE_PG_URL="dbname=ducklake_catalog host=localhost user=ducklake_rw" \
DUCKLAKE_DATA_PATH="/var/lib/safecast/ducklake/" \
./safecast-new-map

# With AI web chat (requires Anthropic API key)
DATABASE_URL="postgres://..." \
ANTHROPIC_API_KEY="your-key" \
CLAUDE_MODEL="claude-haiku-4-5-20251001" \
./safecast-new-map
```

**Connect Claude to the live Safecast data:**
```bash
# Claude Code CLI
claude mcp add --transport http safecast https://simplemap.safecast.org/mcp-http

# Claude.ai: Settings → Integrations → Add custom integration
# URL: https://simplemap.safecast.org/mcp-http
```

**Web Chat:** Open `http://localhost:8765/assistant/` in your browser.

**Swagger API docs:** [simplemap.safecast.org/docs/](https://simplemap.safecast.org/docs/)

Regenerate documentation after API changes:

```bash
# Main app API docs
swag init -g doc.go -o docs/api --parseDependency --parseInternal --parseDependencyLevel 2

# MCP server API docs
cd cmd/mcp-server && swag init -g rest.go
```

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

**Architecture diagram:** [docs/spectral-data-flow.mmd](docs/spectral-data-flow.mmd) — full pipeline from upload → parsers → DB storage → REST API → browser chart construction.

---

## Administration

Enable the admin panel with:
```bash
./safecast-new-map -admin-password your-secure-password
```

Access the admin panel at `/admin/users?password=your-secure-password` or log in as an admin user.

**Admin pages:**

| Page | URL | Description |
|------|-----|-------------|
| Users | `/admin/users` | Manage user accounts, roles, and API keys |
| Uploads | `/admin/uploads` | View uploads, import from Safecast API, delete tracks |
| MCP Analytics | `/admin/mcp` | Monitor MCP tool usage and AI query logs |
| Realtime | `/admin/realtime` | Manage real-time sensor device data |
| Translations | `/admin/translations` | Edit, add, and delete UI translations for all 29 languages |

### Translation Management

Translations are stored in PostgreSQL and loaded into memory at startup. The admin translations page provides:
- Filter by language (29 languages: ar, bg, cs, da, de, el, en, es, fa, fi, fr, he, hi, hu, id, it, ja, ko, ms, nl, no, pl, pt, ru, sv, th, tr, uk, vi, zh)
- Search across keys and values
- Inline editing with save
- Add new translation keys
- **Reload into Memory** button to apply changes without restarting the server

All UI components are translated: map legend, AI assistant widget, login/register modals, search bar, spectrum viewer, profile page, and coordinate input dialog.

**How it works:**
- **Language selection:** The `?lang=` URL parameter takes priority, then the browser's `Accept-Language` header, defaulting to English
- **Incremental seeding:** On startup, any new keys in the embedded `translations.json` are inserted into the DB (`ON CONFLICT DO NOTHING` preserves existing edits)
- **Performance:** Only the active language + English fallback are embedded in the page HTML (~30KB vs ~850KB for all 30 languages), keeping the AI widget within Claude's token limits
- **Branding:** "Safecast" must remain untranslated as a brand name in all languages

---

## User Authentication & API Keys

### User Registration & Login

Enable user registration and authentication with email support:

```bash
./safecast-new-map \
  -allow-registration \
  -smtp-host smtp.gmail.com \
  -smtp-port 587 \
  -smtp-username your-email@gmail.com \
  -smtp-password your-app-password \
  -smtp-from your-email@gmail.com \
  -session-secret "your-random-secret-key" \
  -base-url "https://your-domain.com"
```

**Features:**
- Email-based user registration with verification
- Session-based authentication (30-day sessions)
- Password reset via email
- User profile pages with upload history
- Upload tracking per user

### API Key Authentication

Every registered user receives a unique API key that can be used for:
- **Web login** - Alternative to password authentication
- **Programmatic uploads** - Automated data submission
- **API access** - Protected endpoint authentication

**Two ways to authenticate:**

1. **Password Login** (traditional):
   ```json
   POST /api/auth/login
   {
     "email": "user@example.com",
     "password": "your-password"
   }
   ```

2. **API Key Login** (new):
   ```json
   POST /api/auth/login
   {
     "email": "user@example.com",
     "api_key": "your-20-char-api-key"
   }
   ```

**Using API Keys with HTTP requests:**
   ```bash
   # Header method (recommended)
   curl -H "X-API-Key: your-api-key" https://your-domain.com/api/protected-endpoint

   # Query parameter method
   curl "https://your-domain.com/api/protected-endpoint?api_key=your-api-key"
   ```

### Finding Your API Key

1. Log in to your account
2. Click your username → "Profile"
3. Your API key is displayed in the "API Key" section
4. Click "Copy" to copy it to clipboard

**Security Note:** Treat your API key like a password. Don't share it or commit it to version control.

### Admin: API Key Management

Admins can regenerate API keys for any user:

```bash
# Regenerate API key for user ID 42
curl -X POST "https://your-domain.com/api/admin/users/42/regenerate-api-key?password=admin-password"
```

**Response:**
```json
{
  "message": "API key regenerated successfully",
  "api_key": "new-20-char-key",
  "user_id": 42
}
```

The user will receive an email with their new API key. The old key becomes invalid immediately.

### Authentication Logging

All authentication events are logged for security monitoring:

**Logged Events:**
- ✅ Successful logins (with method: password or api_key)
- ✅ Failed login attempts (with detailed failure reasons)
- ✅ User registrations
- ✅ Email verifications
- ✅ Password changes
- ✅ User logouts
- ✅ API key regenerations

**Log Format:**
```
AUTH: Successful login - user_id=42 email=user@example.com method=password ip=192.168.1.100 user_agent=Mozilla/5.0...
AUTH: Failed login attempt - email=user@example.com method=password reason=invalid_password ip=192.168.1.100
AUTH: New user registered - user_id=43 email=new@example.com username=john_doe ip=192.168.1.100
```

**Filtering logs:**
```bash
# View all authentication events
grep "AUTH:" /var/log/safecast.log

# View failed login attempts
grep "AUTH: Failed login" /var/log/safecast.log

# Monitor for brute force attacks
grep "AUTH: Failed login.*ip=192.168.1.100" /var/log/safecast.log | wc -l
```

### Database Migrations

Add user authentication to existing databases:

**PostgreSQL:**
```bash
psql -h 127.0.0.1 -U postgres -d safecast -f migrations/create_users_table.sql
psql -h 127.0.0.1 -U postgres -d safecast -f migrations/link_historical_uploads_to_users.sql
psql -h 127.0.0.1 -U postgres -d safecast -f migrations/add_ui_translations.sql
```

**Note:** The translations table is auto-created at startup and seeded from the embedded `translations.json` if no DB rows exist. The `add_ui_translations.sql` migration adds the extended UI translations (AI widget, auth modals, search, spectrum, profile page).

**Key columns in users table:**
- `id` - Unique user identifier
- `email` - User email (unique)
- `password_hash` - Bcrypt hashed password
- `api_key` - 20-character unique API key
- `username` - Display name
- `email_verified` - Email verification status
- `is_active` - Account status
- `is_admin` - Admin privileges
- `created_at`, `updated_at`, `last_login_at` - Timestamps

---

## URL Parameters

Customize map views with URL parameters:

| Parameter                              | Values                          | Description                                    |
| -------------------------------------- | ------------------------------- | ---------------------------------------------- |
| `place`                                | City, country, or location name | Navigate to a specific location by name        |
| `lat`                                  | Latitude (-90 to 90)            | Latitude coordinate for map center             |
| `lon`                                  | Longitude (-180 to 180)         | Longitude coordinate for map center            |
| `zoom`                                 | 1-18                            | Zoom level (used with lat/lon)                 |
| `minLat`, `minLon`, `maxLat`, `maxLon` | Coordinates                     | Define map bounds (legacy format)              |
| `coloring`                             | safecast, chicha                | Scientific gradient vs. safety bins            |
| `unit`                                 | uSv, uR                         | Display units (microsieverts or microroentgen) |
| `legend`                               | 1, 0                            | Show/hide legend                               |
| `lang`                                 | en, ja, de, fr, etc. (29 langs) | Interface language (server-side + client-side)  |
| `layer`                                | OpenStreetMap, Google Satellite | Base map                                       |
| `show`                                 | rt                              | Filter to show only realtime sensors           |

**Location Examples:**
- Direct to Tokyo: `/?place=Tokyo`
- Coordinates with zoom: `/?lat=48.8566&lon=2.3522&zoom=13`
- Share search result: `/?place=São%20Paulo` (auto-generated when searching)

**Display Examples:**
- Safety view: `/?coloring=chicha&unit=uR`
- Clean embed: `/?legend=0`
- Russian interface: `/?lang=ru`
- Realtime sensors only: `/?show=rt`

**Parameter Priority:** When loading URLs, `place` takes priority over `lat/lon`, which takes priority over bounds (`minLat/maxLat`).

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
go build -o safecast-new-map ./cmd/unified-server
./safecast-new-map
```

### Run Tests

```bash
go test ./...
```

### Utility Tools

The `cmd/tools/` directory contains standalone utilities for database maintenance and migration:

| Tool                   | Usage                             |
| ---------------------- | --------------------------------- |
| `fix-pg-sequence`      | Reset PostgreSQL markers sequence |
| `fix-sequence`         | Alternative sequence sync utility |
| `add-internal-user-id` | Add internal user IDs to uploads  |
| `cleanup-test-users`   | Remove test users from database   |
| `import-api-keys`      | Import API keys from CSV          |
| `migrate-to-postgres`  | Migrate from SQLite to PostgreSQL |
| `migrate-users`        | Migrate user data                 |

Run any tool with:
```bash
go run ./cmd/tools/fix-pg-sequence
export DATABASE_URL="postgres://user:pass@localhost/db" && go run ./cmd/tools/fix-pg-sequence
```

### Cross-Compile

```bash
# See scripts/crosscompile/crosscompile.go
go run scripts/crosscompile/crosscompile.go
```

### Production Deployment

For deploying to production servers with CloudFront/CDN setup, see:
- [Deployment Guide](docs/DEPLOYMENT.md) - Complete deployment procedures and troubleshooting
- [GitHub Actions Guide](GITHUB_ACTIONS_GUIDE.md) - Automated CI/CD setup
- [CloudFront Setup](docs/cloudfront-setup.md) - CDN configuration details
- [Architecture Diagram](docs/architecture-overview.mmd) - System overview (Mermaid)
- [Spectral Data Flow](docs/spectral-data-flow.mmd) - Spectral upload/parse/store/render pipeline (Mermaid)

**Important:** When using CloudFront, SSH and rsync must use the server's IP address directly, not the domain name.

---

## Performance Notes

**Analytics:** The platform uses DuckLake for analytics (in-memory DuckDB with a PostgreSQL catalog and Parquet data files), allowing multiple services to share analytics tables concurrently. Set `DUCKLAKE_PG_URL` and `DUCKLAKE_DATA_PATH` environment variables to configure.

**PostgreSQL Tuning:** Use connection pooling and appropriate `work_mem` settings for large batch imports.

**Rate Limiting:** API endpoints are rate-limited by default. Adjust in code if needed for high-traffic deployments.

---

## Database Maintenance

### PostgreSQL Sequence Reset

After data migration or restore, you may encounter "duplicate key value violates unique constraint" errors. This happens when PostgreSQL sequences are out of sync with existing data.

**Fix the markers sequence:**
```bash
psql -h 127.0.0.1 -U postgres -d safecast -c \
  "SELECT setval('markers_id_seq', (SELECT MAX(id) FROM markers) + 1);"
```

**Fix all sequences at once:**
```bash
./tools/reset_postgres_sequences.sh
```

### Adding Missing Columns

If you see errors like `column X does not exist`, the database schema may need updating:

```bash
psql -h 127.0.0.1 -U postgres -d safecast -c "
  ALTER TABLE uploads ADD COLUMN IF NOT EXISTS recording_date TIMESTAMPTZ;
  ALTER TABLE uploads ADD COLUMN IF NOT EXISTS detector TEXT;
  ALTER TABLE uploads ADD COLUMN IF NOT EXISTS username TEXT;
  ALTER TABLE uploads ADD COLUMN IF NOT EXISTS internal_user_id TEXT;
"
```

### Useful Maintenance Scripts

| Script                                    | Purpose                                        |
| ----------------------------------------- | ---------------------------------------------- |
| `tools/reset_postgres_sequences.sh`       | Reset all PostgreSQL sequences after migration |
| `tools/reset_postgres_marker_sequence.sh` | Reset only the markers sequence                |
| `tools/refresh_track_stats.sh`            | Refresh track statistics materialized view     |
| `tools/populate_usernames.sh`             | Fetch usernames from Safecast API for uploads  |

---

## Community & Support

This project is developed and maintained by the Safecast community with contributions from:
- [ Matvey Gladkikh](https://github.com/matveynator) (primary designer/developer) 
- Rob Oudendijk (developer)
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
