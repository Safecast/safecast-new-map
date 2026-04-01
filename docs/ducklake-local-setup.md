# DuckLake Local Development Setup

DuckLake is used for shared analytics tables (`chat_questions`, `mcp_ai_query_log`, `mcp_query_log`).
It uses PostgreSQL as the catalog and Parquet files for data storage.

## Prerequisites

- PostgreSQL 16 running locally
- `ducklake_catalog` database and `ducklake_rw` user already created (done during initial setup)

## One-time: Set ducklake_rw password

The `ducklake_rw` user requires a password for DuckDB to connect via the `postgres` extension:

```bash
psql -h 127.0.0.1 -U postgres -c "ALTER USER ducklake_rw PASSWORD 'your-password-here';"
```

Verify it works:

```bash
psql "postgresql://ducklake_rw:your-password-here@127.0.0.1:5432/ducklake_catalog" -c "SELECT 1;"
```

## Environment variables for local-server-config.sh

Add these to `local-server-config.sh` alongside `ANTHROPIC_API_KEY`:

```bash
export DUCKLAKE_PG_URL="postgresql://ducklake_rw:your-password-here@127.0.0.1:5432/ducklake_catalog"
export DUCKLAKE_DATA_PATH="/var/lib/safecast/ducklake/"
```

## Data directory

The Parquet data directory must exist and be writable:

```bash
sudo mkdir -p /var/lib/safecast/ducklake/
sudo chown $USER /var/lib/safecast/ducklake/
```

## What happens on startup

When `DUCKLAKE_PG_URL` is set and the connection succeeds, `RegisterMCP()` calls
`initDuckDBAnalytics()` which:

1. Opens an in-memory DuckDB instance
2. Attaches the DuckLake catalog via PostgreSQL
3. Creates the three analytics tables if they don't exist:
   - `chat_questions` — web-chat and map widget conversations
   - `mcp_ai_query_log` — MCP tool execution logs
   - `mcp_query_log` — MCP tool usage stats

If `DUCKLAKE_PG_URL` is missing or the connection fails, the server logs a warning and
continues without analytics (chat still works, logs are just dropped).

## Troubleshooting

**`Table with name mcp_ai_query_log does not exist`**
→ DuckLake attachment failed silently. Check server logs for `Warning: DuckDB initialization failed`.
→ Most likely cause: `ducklake_rw` password not set or `DUCKLAKE_PG_URL` not exported.

**`fe_sendauth: no password supplied`**
→ The `ducklake_rw` user has a password set but the connection string omits it.
→ Set `DUCKLAKE_PG_URL` with the password included (see above).
