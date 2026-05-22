# Production Deployment Guide

This guide covers deploying to the production server at simplemap.safecast.org.

## Infrastructure Overview

**Domain:** simplemap.safecast.org
**Server IP:** 65.108.24.131
**CDN:** AWS CloudFront (Distribution ID: E12FYIQ8RRXOJ1)

### Services

The unified service is built from this repo and deployed by the **Codeberg Actions** workflow (`.forgejo/workflows/deploy.yml`). GitHub Actions runs tests and releases only.

| Service | Binary | Port | Location on VPS |
|---------|--------|------|-----------------|
| Unified server (Map + MCP + chat) | `safecast-new-map` | 8765 | `/usr/local/bin/safecast-new-map` |

### Traffic Flow

```
┌─────────────────────────────────────────────────────────┐
│                   Web Traffic (Ports 80/443)            │
└─────────────────────────────────────────────────────────┘
  User Browser
       ↓
  simplemap.safecast.org (DNS → CloudFront)
       ↓
  CloudFront Edge Locations (Global CDN)
       ↓
  Origin Server: 65.108.24.131


┌─────────────────────────────────────────────────────────┐
│              Deployment Traffic (Port 22)               │
└─────────────────────────────────────────────────────────┘
  Developer Machine
       ↓
  SSH/Rsync to 65.108.24.131 (Direct IP)
       ↓
  Server: 65.108.24.131
```

## ⚠️ Critical: SSH Must Use IP Address

**CloudFront only handles HTTP/HTTPS traffic.** SSH connections MUST go directly to the server IP.

**✅ Correct:**
```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131
```

**❌ Wrong (will fail):**
```bash
ssh -i ~/.ssh/safecast-deploy root@simplemap.safecast.org  # CloudFront can't handle SSH!
```

## Manual Deployment

### Prerequisites

1. **SSH Key:** `~/.ssh/safecast-deploy` (private key)
2. **Binaries built:**
   ```bash
   go build -o safecast-new-map ./cmd/unified-server
   ```
3. **Server access:** Ability to SSH to 65.108.24.131

### Deployment Steps — Map Server

```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl stop safecast-new-map"
rsync -avP -e "ssh -i ~/.ssh/safecast-deploy" ./safecast-new-map root@65.108.24.131:/usr/local/bin/
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl start safecast-new-map && systemctl status safecast-new-map"
```

### Invalidate CloudFront cache (optional)
```bash
aws cloudfront create-invalidation --distribution-id E12FYIQ8RRXOJ1 --paths "/*"
```

### Why This Order Matters

1. **Stop first:** Prevents file conflicts when replacing the binary
2. **Sync binary:** Upload new version while service is stopped
3. **Start service:** New version loads cleanly
4. **Verify:** Ensure service started successfully

## Automated Deployment (Codeberg Actions)

**Current workflow:** `.forgejo/workflows/deploy.yml` — triggers on push to `main`.

GitHub Actions (`.github/workflows/`) runs tests and multi-platform releases only. Deploy is Codeberg-only.

### CI Split

| Workflow | GitHub | Codeberg |
|----------|--------|----------|
| `test.yml` | ✅ PRs + main | ✅ PRs + main |
| `deploy.yml` | ❌ disabled | ✅ push to main |
| `release.yml` | ✅ multi-platform | not present |

### How It Works

1. **Trigger:** Automatic on push to `main`, or manual via workflow_dispatch
2. **Build:** Compiles Go binary on `codeberg-medium` runner
3. **Deploy:** Stop → rsync binary → start service via SSH to 65.108.24.131
4. **Health check:** Polls `http://localhost:8765/api/health` up to 5 times
5. **Invalidate:** Clears CloudFront cache

### Required Secrets (set in Codeberg repo Settings → Secrets)

| Secret | Purpose |
|--------|---------|
| `DEPLOY_SSH_KEY` | SSH key for deploying to 65.108.24.131 |
| `AWS_ACCESS_KEY_ID` | CloudFront cache invalidation |
| `AWS_SECRET_ACCESS_KEY` | CloudFront cache invalidation |
| `DATABASE_URL` | Postgres connection string for unified server |
| `DUCKLAKE_PG_URL` | DuckLake PostgreSQL catalog connection |
| `DUCKLAKE_DATA_PATH` | Path for DuckLake Parquet data files |
| `ANTHROPIC_API_KEY` | Claude API key for AI assistant |

**Note:** All SSH/rsync in the workflow uses the IP address (65.108.24.131), not the domain.

## CloudFront Considerations

### Cache Invalidation

After deploying new code, CloudFront may serve cached versions of static files. Create an invalidation to force cache refresh:

```bash
# Invalidate all cached content
aws cloudfront create-invalidation --distribution-id E12FYIQ8RRXOJ1 --paths "/*"

# Invalidate specific paths
aws cloudfront create-invalidation --distribution-id E12FYIQ8RRXOJ1 --paths "/index.html" "/static/*"
```

**Cost:** First 1,000 invalidation paths per month are free, then $0.005 per path.

### Cache Behavior

Current configuration:
- **API endpoints:** No caching (Cache-Control headers prevent it)
- **Static files:** Cached at edge locations
- **Uploads:** Not cached (Cache-Control: private)
- **Admin pages:** Not cached (Cache-Control: no-store)

### WAF Configuration

AWS WAF protects against:
- SQL injection
- Cross-site scripting (XSS)
- ~~Large request bodies~~ (Changed to "Count" mode for file uploads)

**Important:** `SizeRestrictions_BODY` rule is in "Count" mode to allow large file uploads. See [docs/cloudfront-fix-waf-403.md](cloudfront-fix-waf-403.md) for details.

## Analytics (DuckLake)

The unified server uses DuckLake for analytics (tool usage logs, chat questions). Architecture: in-memory DuckDB attaches a shared DuckLake catalog backed by PostgreSQL + Parquet files.

**Required env vars** (set in `.env` files or systemd service):
- `DUCKLAKE_PG_URL` — PostgreSQL connection for DuckLake catalog (e.g., `dbname=ducklake_catalog host=localhost user=ducklake_rw`)
- `DUCKLAKE_DATA_PATH` — Directory for Parquet data files (e.g., `/var/lib/safecast/ducklake/`)

**Shared tables:** `chat_questions`, `mcp_query_log`, `mcp_ai_query_log`,
`qa_embeddings` (semantic chat cache), `location_knowledge` (curated geo notes).

### Q&A Semantic Cache (May 2026)

The web-chat assistant caches answered questions in `qa_embeddings` and serves
paraphrased repeats without re-calling the LLM. Hit criteria: cosine similarity
≥ 0.85, `feedback_score > 0`, `status='active'`, matching `lang`.

**Admin UI:** `/admin/qa-embeddings` — searchable/sortable, promote / demote /
archive / restore actions per row.

**Schema:** base table is created at startup; curation columns (`used_count`,
`last_used_at`, `status`, `lang`) are applied in-place via idempotent
`ALTER TABLE ADD COLUMN IF NOT EXISTS` on every boot. Standalone migration
file: `migrations/add_qa_embeddings_columns.sql`.

**Manual inspection** (DuckDB CLI against the live catalog):
```sql
ATTACH 'ducklake:postgres:dbname=ducklake_catalog host=localhost user=ducklake_rw'
       AS analytics (DATA_PATH '/var/lib/safecast/ducklake/');
USE analytics;
SELECT id, chat_id, feedback_score, used_count, status, lang,
       LEFT(question, 80) AS q
FROM qa_embeddings
ORDER BY created_at DESC
LIMIT 20;
```

The cache only runs when the unified server is built with `-tags duckdb`.

## Translations (i18n)

Translations are stored in PostgreSQL (`translations` table) and loaded into memory at startup.

### How Seeding Works

On every startup, `seedTranslationsDB()` reads the embedded `translations.json` and inserts any missing keys using `ON CONFLICT DO NOTHING`. This means:
- New translation keys added to `translations.json` are automatically seeded on next deploy
- Existing DB values (including admin edits) are never overwritten
- No manual migration is needed when adding new keys

### Language Selection Priority

1. `?lang=` URL parameter (e.g., `/?lang=ja`) — checked first, server-side
2. `Accept-Language` HTTP header from the browser
3. Falls back to English (`en`)

### Performance: Filtered TranslationsJSON

Only the active language + English fallback are embedded in the page HTML (`TranslationsJSON`), reducing the payload from ~850KB (all 30 languages) to ~30KB. This keeps the AI assistant widget within Claude's 200K token context limit. Server-side template rendering (`{{translate "key"}}`) still uses the full translations map.

### Branding Rule

"Safecast" must remain untranslated as a brand name in all languages. Never translate it to local equivalents.

**Admin UI:** `/admin/translations` — edit translations live, then click "Reload into Memory" to apply without restart.

**Supported languages (29):** ar, bg, cs, da, de, el, en, es, fa, fi, fr, he, hi, hu, id, it, ja, ko, ms, nl, no, pl, pt, ru, sv, th, tr, uk, vi, zh

**Translated components:** Map legend, AI assistant widget, login/register/forgot-password modals, user menu, search bar, spectrum viewer, coordinate input dialog, profile page.

### Fixing Translations via SQL

To fix a translation directly in the production DB:
```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 \
  "psql -h 127.0.0.1 -U postgres -d safecast -c \"UPDATE translations SET value = 'New value' WHERE language_code='ja' AND key='title'\""
```
Then restart the service to reload: `systemctl restart safecast-new-map`

## Server Configuration

### Service Details

| Service | systemd name | Binary | Config |
|---------|-------------|--------|--------|
| Unified server | `safecast-new-map` | `/usr/local/bin/safecast-new-map` | flags + `Environment=` lines in service file |

### Required Environment Variables — safecast-new-map service

These must be set as `Environment=` lines in `/etc/systemd/system/safecast-new-map.service`. Without them, API documentation cross-links fall back to `http://localhost:8765`.

| Variable | Value | Purpose |
|----------|-------|---------|
| `MAP_BASE_URL` | `https://simplemap.safecast.org` | Base URL for the Map API docs "Switch to MCP API" button |
| `MCP_BASE_URL` | `https://simplemap.safecast.org` | Base URL for the MCP API docs "Switch to Map API" button |
| `DUCKLAKE_PG_URL` | `dbname=ducklake_catalog host=localhost user=ducklake_rw password=...` | DuckLake catalog |
| `DUCKLAKE_DATA_PATH` | `/var/lib/safecast/ducklake/` | DuckLake Parquet data path |
| `ANTHROPIC_API_KEY` | `sk-ant-...` | Claude API key |

**Example service file snippet:**
```ini
Environment=MAP_BASE_URL=https://simplemap.safecast.org
Environment=MCP_BASE_URL=https://simplemap.safecast.org
Environment=DUCKLAKE_PG_URL="dbname=ducklake_catalog host=localhost user=ducklake_rw password=SECRET"
Environment=DUCKLAKE_DATA_PATH=/var/lib/safecast/ducklake/
Environment=ANTHROPIC_API_KEY=sk-ant-...
```

After editing the service file, always reload and restart:
```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 \
  "systemctl daemon-reload && systemctl restart safecast-new-map"
```

### Useful Commands

```bash
# View service status (replace service name as needed)
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl status safecast-new-map"

# View logs
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "journalctl -u safecast-new-map -f"

# Restart service
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl restart safecast-new-map"

# Check disk space
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "df -h"

# Check memory usage
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "free -h"
```

## Troubleshooting

### SSH Connection Fails

**Symptom:** `ssh: connect to host simplemap.safecast.org port 22: Connection refused`

**Cause:** Trying to SSH to the domain name instead of IP address.

**Solution:** Use IP address: `65.108.24.131`

### Deployment Succeeds But Changes Not Visible

**Cause:** CloudFront is serving cached content.

**Solution:** Invalidate CloudFront cache (see above).

### Upload 403 Errors After Deployment

**Cause:** WAF blocking large uploads or missing cache headers.

**Solution:**
1. Verify WAF rules are in "Count" mode (see [cloudfront-fix-waf-403.md](cloudfront-fix-waf-403.md))
2. Check Cache-Control headers in code (safecast-new-map.go)

### Service Won't Start

```bash
# Check service status
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl status safecast-new-map"

# View recent logs
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "journalctl -u safecast-new-map -n 50"

# Check if port is already in use
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "netstat -tulpn | grep 8765"
```

## Rollback Procedure

If deployment fails:

```bash
# Stop broken version
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl stop safecast-new-map"

# Restore previous binary (if backed up)
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "cp /usr/local/bin/safecast-new-map.old /usr/local/bin/safecast-new-map"

# Start service
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl start safecast-new-map"
```

## Security Notes

### SSH Key Management

- **Private key:** Keep `~/.ssh/safecast-deploy` secure and never commit to git
- **GitHub Secret:** Stored encrypted in GitHub, only accessible to workflows
- **Server:** Public key in `/root/.ssh/authorized_keys` on 65.108.24.131

### CloudFront Security

- **HTTPS only:** HTTP requests redirected to HTTPS
- **WAF enabled:** Protects against common web attacks
- **DDoS protection:** AWS Shield Standard included with CloudFront
- **Origin protection:** Origin server (65.108.24.131) can be firewalled to only accept CloudFront IPs

### PostgreSQL Security

PostgreSQL (port 5432) **must never be exposed to the internet.**

**Configuration** — `/etc/postgresql/16/main/postgresql.conf`:
```
listen_addresses = 'localhost'
```

**Firewall rules** (persisted via `iptables-persistent`):
```bash
# Allow postgres only on loopback
iptables -A INPUT -i lo -p tcp --dport 5432 -j ACCEPT
# Drop all external access
iptables -A INPUT -p tcp --dport 5432 -j DROP
```

Rules are saved in `/etc/iptables/rules.v4` and restored automatically on reboot.

> **Background:** In March 2026 the BSI (via Hetzner abuse) flagged port 5432 as publicly accessible. The root cause was `listen_addresses` including the public IP. Both the config and firewall were fixed and the rules persisted.

## Database Migrations

### Running Migrations on Production

Migration SQL files are in `migrations/`. Run them directly via SSH:

```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 \
  "psql -h 127.0.0.1 -U postgres -d safecast -f -" < migrations/your_migration.sql
```

Or for inline changes:
```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 \
  "psql -h 127.0.0.1 -U postgres -d safecast -c 'ALTER TABLE uploads ADD COLUMN IF NOT EXISTS comment VARCHAR;'"
```

### Columns Added Since Initial Import (Apr 2026)

The following columns were added to the `uploads` table after the initial Safecast API import and must be migrated manually on any fresh production database:

| Column | Migration file | Notes |
|--------|---------------|-------|
| `name` | `add_upload_metadata.sql` | Display name; back-filled from `filename` |
| `notes` | `add_upload_metadata.sql` | Admin-only internal notes |
| `comment` | *(inline)* | User comment from old Safecast API |

After adding `name`/`comment`, back-fill from the old Safecast API (see below).

### Tour Steps Table (Apr 2026)

The map tutorial tour is now DB-backed with multilingual text and an admin UI at `/admin/tour`. A fresh database (or any production server deployed before this change) must run:

```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 \
  "sudo -u postgres psql -d safecast -f -" < migrations/add_tour_steps_table.sql
```

This creates the `tour_steps` table (step_key, sort_order, selector, center, enabled, tri-state conditions, viewport, first_time_only, updated_at) plus the `idx_tour_steps_order` index. **The server will NOT start seeding tour steps into a non-existent table** — if you restart the service before running the migration, `seedTourStepsDB()` will log errors and the `/admin/tour` page will 500 on list fetch. Run the migration first, then restart.

Per-step text lives in the existing `translations` table under keys `tour.<step_key>.text`. `seedTourStepsDB()` inserts the 11 default steps and English translations on every startup via `ON CONFLICT DO NOTHING`, so admin edits are preserved across deploys.

### Back-filling Metadata from the Old Safecast API

The admin Uploads page has an **"Import Safecast API Metadata"** button that fetches `name` and `comment` for all `safecast-api` tracks from `api.safecast.org`. For new imports this is automatic; the button is only needed as a one-time backfill for rows that existed before the columns were added.

#### ⚠️ CloudFront timeout

The button makes a synchronous browser request. CloudFront cuts connections after ~60 seconds, which only processes ~1,000 tracks before the request is killed. For a full backfill (~47k tracks, ~15 min), **run directly on the server via SSH** to bypass CloudFront:

```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 \
  "curl -s -X POST 'http://localhost:8765/api/admin/tracks/import-safecast' \
   -H 'Content-Type: application/json' \
   -d '{\"password\":\"ADMIN_PASSWORD\"}' \
   --max-time 1800"
```

Expected response: `{"ok":true,"total":46380,"updated":46380}`

The admin password is in the systemd service file:
```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl cat safecast-new-map | grep admin-password"
```

**Note:** Many tracks have no comment in the old API (the uploader never wrote one) — this is normal. Only tracks where the uploader provided a description will have a non-empty `comment`.

## API Documentation Pages

### Overview

Three Swagger UI pages are served by the unified server (`safecast-new-map`, port 8765):

| URL | Description |
|-----|-------------|
| `/docs/` | **Combined page** — Map API and MCP API in a tabbed interface |
| `/map-api/` | Map API only (standalone Swagger UI) |
| `/mcp-api/` | MCP API only (standalone Swagger UI) |

All three support dark mode. `/docs/` is the canonical entry point.

### Nginx routing

All three paths are proxied from both nginx configs to port 8765:

```nginx
location /docs/     { proxy_pass http://localhost:8765/docs/; ... }
location /mcp-api/  { proxy_pass http://localhost:8765/mcp-api/; ... }
location /map-api/  { ... }  # falls through to the default location / block → 8765
```

**Note:** `/mcp-api/` requires an explicit nginx `location` block.

### Combined docs page (`/docs/`)

Implemented in `cmd/unified-server/rest.go` — `serveAPIDocsPage()` handler.

- Admin-style CSS variables (`--bg-primary`, `--bg-card`, `--text-primary`, etc.) matching all admin pages
- Tab bar uses the same `.api-tabs` pattern as admin `.admin-tabs`
- Dark mode stored in `localStorage` key `safecastDocTheme`, applied via `data-theme` on `<html>`
- MCP API swagger initialized **lazily** — only on first tab click (avoids loading both at page load)
- Active tab persisted in `localStorage` key `safecastDocTab`
- Both swagger instances use `BaseLayout` (no duplicate swagger topbars)
- Swagger UI assets reused from `/map-api/swagger-ui-bundle.js`

### Code block visibility fix (light mode)

Swagger's microlight syntax highlighter was rendering code examples with near-invisible light text on a light background in light mode. Fixed in both theme CSS constants (`mapSwaggerThemeCSS`, `mcpSwaggerThemeCSS`) in `cmd/unified-server/rest.go`:

```css
.swagger-ui .microlight, .swagger-ui pre.microlight {
  color: #24292e !important;  /* dark charcoal — was unset, defaulting to near-white */
}
.swagger-ui .microlight span { color: inherit !important; }
```

### `swaggerFiles` singleton conflict fix

`github.com/swaggo/files` exposes a **package-level webdav singleton** (`swaggerFiles.Handler`). When two `httpSwagger.Handler` instances are registered on the same `http.ServeMux` (e.g. `/map-api/` and `/mcp-api/`), each closure sets `handler.Prefix` via its own `sync.Once` on first request. Whichever runs second overwrites the prefix, causing the first handler to return 404 for all its static assets (`swagger-ui-bundle.js`, `swagger-ui.css`, etc.).

**Fix (PR #63):** Only one `httpSwagger.Handler` is registered on port 8765 (`/map-api/`). The `/mcp-api/` path uses:

- `/mcp-api/doc.json` — served directly via `swag.ReadDoc("swagger")` (no webdav involved)
- `/mcp-api/` — custom HTML page (`serveMCPAPIPage` in `rest.go`) that loads all static assets from `/map-api/`

The combined `/docs/` page was also already loading assets from `/map-api/`, so it benefits from this fix too.

### Preamble dark mode fix (`/map-api/`)

The white preamble header box (title, description, nav buttons) on `/map-api/` was not responding to dark mode. Fixed by injecting a `<style>` element via the `onComplete` JS callback with `body.dark-mode #safecast-preamble` rules.

## Route Ownership

All routes live in the unified server binary (`cmd/unified-server/`, port 8765).

### HTTP route composers

| Composer | Owns |
|----------|------|
| `pkg/httpapi/register.go` | All `/api/*` map-server routes |
| `pkg/httpapi/handlers_core.go` | Core REST handlers |
| `pkg/auth/*` | Auth routes (`/api/auth/`, `/api/user/`) |
| `cmd/unified-server/mcp_register.go` | MCP surface wiring |

### MCP surface (same binary, same port 8765)

| Path | Description |
|------|-------------|
| `/mcp-http` | MCP HTTP transport |
| `/mcp/sse` | MCP SSE transport |
| `/mcp-api/*` | MCP Swagger UI |

Shared MCP composers: `pkg/mcpserver/tools.go`, `transports.go`, `routes.go`, `docs.go`

### Guardrails

- Do not add `/api/*` routes directly in `main.go` — use the registry composers.
- For MCP changes: update `pkg/mcpserver/*`, keep `cmd/unified-server/` as thin wiring.
- When adding an endpoint: implement in the owning package → register via composer → add route inventory test.

### Nginx routing (production)

All traffic hits nginx on port 443, which proxies to port 8765. Key path rules:

```
/mcp-http        → 8765 (MCP HTTP transport)
/mcp/sse         → 8765 (MCP SSE transport)
/mcp-api/        → 8765 (MCP Swagger UI)
/assistant/      → 8765 (AI assistant)
/docs/           → 8765 (combined API docs)
/map-api/        → 8765 (Map API Swagger)
/api/*           → 8765 (REST API)
/                → 8765 (map UI + everything else)
```

> **Note:** Use specific `location = /api/endpoint` rules for any future path-based routing splits — a broad `location /api/` rule will break things.

## Related Documentation

- [CloudFront Setup Guide](cloudfront-setup.md) - Initial CloudFront configuration
- [Upload 403 Fix](cloudfront-fix-upload-403.md) - Cookie forwarding configuration
- [WAF 403 Fix](cloudfront-fix-waf-403.md) - Large file upload configuration

## Quick Reference

**Server IP:** 65.108.24.131
**SSH Key:** `~/.ssh/safecast-deploy`
**CloudFront Distribution ID:** E12FYIQ8RRXOJ1
**Service Name:** `safecast-new-map`
**Binary Path:** `/usr/local/bin/safecast-new-map`

**API Documentation URLs:**
- **Combined docs (tabs):** `https://simplemap.safecast.org/docs/`
- Map API only: `https://simplemap.safecast.org/map-api/`
- MCP API only: `https://simplemap.safecast.org/mcp-api/`

**One-Line Deploy:**
```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl stop safecast-new-map" && \
rsync -avP -e "ssh -i ~/.ssh/safecast-deploy" ./safecast-new-map root@65.108.24.131:/usr/local/bin/ && \
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl start safecast-new-map && systemctl status safecast-new-map"
```
