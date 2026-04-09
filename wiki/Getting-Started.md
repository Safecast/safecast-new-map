# Getting Started

This guide will help you get Safecast New Map up and running in minutes.

[[Home|← Back to Home]]

---

## Quick Start Options

### Option 1: Download Binary (Fastest)

Download and run in seconds:

```bash
# Download from https://github.com/safecast/safecast-new-map/releases
# Choose your platform's binary

chmod +x ./safecast-new-map
./safecast-new-map
```

Open [http://localhost:8765](http://localhost:8765) in your browser.

### Option 2: With Production Data

Start with a complete dataset from simplemap.safecast.org:

```bash
./safecast-new-map -import-tgz-url https://simplemap.safecast.org/api/json/weekly.tgz
```

This imports all public tracks and starts the server. Future runs will use the cached data locally.

### Option 3: Docker

```bash
docker run -d -p 8765:8765 --name safecast-map safecastr/safecast-new-map:latest
```

Open [http://localhost:8765](http://localhost:8765)

### Option 4: Build from Source

```bash
git clone https://github.com/Safecast/safecast-new-map.git
cd safecast-new-map
go build -o safecast-new-map ./cmd/unified-server
./safecast-new-map
```

---

## First Configuration

### Minimal Setup

The server runs with sensible defaults. Just start it:

```bash
./safecast-new-map
```

Default settings:
- **Port:** 8765
- **Database:** SQLite (local file)
- **Location:** Default coordinates (44.08832, 42.97577)
- **Zoom:** 11
- **Base map:** OpenStreetMap

### Basic Configuration Flags

```bash
./safecast-new-map \
  -port 8765 \
  -db-type sqlite \
  -db-path /path/to/data \
  -default-lat 35.6762 \
  -default-lon 139.6503 \
  -default-zoom 12
```

### Enable Admin Panel

```bash
./safecast-new-map -admin-password your-secure-password
```

Access at: `/admin/users?password=your-secure-password`

---

## Database Setup

Choose a database backend based on your needs:

### SQLite (Simple, Single-User)

No setup required. Data stored in local file.

```bash
./safecast-new-map -db-type sqlite -db-path /path/to/data.db
```

**Best for:** Personal use, testing, small deployments

### DuckDB (Fast Local Storage)

```bash
./safecast-new-map -db-type duckdb -db-path /path/to/data
```

**Best for:** Fast local analytics, single-user scenarios

### PostgreSQL (Production Recommended)

Requires PostGIS extension:

```bash
# Ubuntu/Debian
sudo apt install postgresql-16-postgis-3

# RHEL/CentOS/Rocky
sudo dnf install postgis34_16
```

```bash
./safecast-new-map \
  -db-type pgx \
  -db-conn "postgres://user:pass@host:5432/dbname?sslmode=require"
```

**Best for:** Production, multi-user, large datasets

### ClickHouse (Large-Scale Analytics)

```bash
./safecast-new-map \
  -db-type clickhouse \
  -db-conn "clickhouse://user:pass@host:9000/dbname?secure=true"
```

**Best for:** Analytics workloads, large-scale deployments

See [Database Setup](Database-Setup) for detailed configuration.

---

## Enable Features

### User Authentication

```bash
./safecast-new-map \
  -allow-registration \
  -require-auth \
  -session-secret "your-random-secret-key"
```

### Email Notifications

```bash
./safecast-new-map \
  -smtp-host smtp.gmail.com \
  -smtp-port 587 \
  -smtp-username your-email@gmail.com \
  -smtp-password your-app-password \
  -smtp-from your-email@gmail.com \
  -base-url "https://your-domain.com"
```

### Automated Data Sync

**Real-time sensors:**
```bash
./safecast-new-map -safecast-realtime
```

**Safecast API fetcher:**
```bash
./safecast-new-map \
  -safecast-fetcher \
  -safecast-fetcher-interval 5m \
  -safecast-fetcher-batch-size 10 \
  -safecast-fetcher-start-date 2024-01-01
```

### MCP Server & AI Integration

```bash
# With DuckLake analytics
DATABASE_URL="postgres://user:pass@localhost/db" \
DUCKLAKE_PG_URL="dbname=ducklake_catalog host=localhost" \
DUCKLAKE_DATA_PATH="/var/lib/safecast/ducklake/" \
./safecast-new-map

# With AI web chat
DATABASE_URL="postgres://..." \
ANTHROPIC_API_KEY="your-key" \
CLAUDE_MODEL="claude-haiku-4-5-20251001" \
./safecast-new-map
```

Access AI chat at: `http://localhost:8765/assistant/`

---

## Production Deployment

### With HTTPS (Let's Encrypt)

```bash
./safecast-new-map \
  -domain maps.example.org \
  -db-type pgx \
  -db-conn "postgres://user:pass@localhost/safecast"
```

Requires ports 80 and 443 open for certificate validation.

### Docker Production

```bash
docker run -d \
  -p 8765:8765 \
  -v /path/to/data:/data \
  -e DATABASE_URL="postgres://..." \
  safecastr/safecast-new-map:latest
```

See [Deployment](Deployment) for CloudFront, AWS, and advanced setups.

---

## Verify Installation

### Check Server is Running

```bash
curl http://localhost:8765/api/stats
```

Should return JSON with database statistics.

### Check API Documentation

Open: [http://localhost:8765/map-api/](http://localhost:8765/map-api/)

### Check MCP API

Open: [http://localhost:8765/mcp-api/](http://localhost:8765/mcp-api/)

---

## Next Steps

- [📥 Data Import & Export](Data-Import-Export) - Import your radiation data
- [🔌 API Documentation](API-Documentation) - Use the REST API
- [🤖 MCP Server & AI Integration](MCP-Server-AI-Integration) - Enable AI features
- [🔐 User Authentication](User-Authentication) - Set up users and API keys
- [⚙️ Admin Panel](Admin-Panel) - Manage your deployment
- [🗺️ Map Features](Map-Features) - Explore map features and visualization

---

## Troubleshooting

### Port Already in Use

Change the port:
```bash
./safecast-new-map -port 9000
```

### Database Connection Failed

Verify connection string:
```bash
./safecast-new-map -db-type sqlite -db-path ./data.db
```

### PostGIS Not Available

Install PostGIS:
```bash
# Ubuntu/Debian
sudo apt install postgresql-16-postgis-3

# Or use SQLite for testing
./safecast-new-map -db-type sqlite
```
