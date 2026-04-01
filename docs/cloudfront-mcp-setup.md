# CloudFront + MCP Server Setup

## Overview

All three services (map server, MCP server, web-chat assistant) live in the **same repository** (`safecast-new-map`) and are built and deployed together via a single GitHub Actions workflow. They run behind CloudFront, with nginx routing requests to the correct backend based on the URL path.

## Repository Structure

```
safecast-new-map/
├── safecast-new-map.go   # Map server (port 8765)
├── cmd/
│   ├── mcp-server/       # MCP + REST API server (port 3333)
│   └── web-chat/         # Claude AI chat UI (port 3334)
└── .github/workflows/
    └── deploy.yml        # Builds and deploys all 3 binaries
```

## Architecture

```
User Request
    ↓
CloudFront (simplemap.safecast.org)
    ↓ HTTPS (port 443)
origin-simplemap.safecast.org
    ↓
Nginx (routes by path)
    ├─ /mcp-http         → MCP Server (port 3333)
    ├─ /mcp              → MCP Server (port 3333)
    ├─ /mcp-api/         → MCP Server (port 3333) — Swagger UI
    ├─ /assistant/       → Web-chat (port 3334)
    ├─ /api/* (selected MCP REST routes) → MCP Server (port 3333)
    └─ /                 → Map Server (port 8765) — everything else
                           including /api/auth/, /api/user/,
                           /api/admin/, /api/spectrum/, /api/markers/
```

> **Important:** The nginx config uses specific location blocks for MCP endpoints rather than a broad `/api/` rule. This ensures that map-server API paths (`/api/auth/`, `/api/spectrum/`, etc.) are not accidentally routed to the MCP server.

## Key Configuration

### 1. SSL Certificate

The Let's Encrypt certificate includes both domains:
```
DNS:simplemap.safecast.org
DNS:origin-simplemap.safecast.org
```

Renew/expand certificate:
```bash
certbot certonly --nginx \
  -d simplemap.safecast.org \
  -d origin-simplemap.safecast.org \
  --expand
```

### 2. Nginx Configuration

**File:** `/etc/nginx/sites-available/origin-simplemap.safecast.org`

This is the config that CloudFront hits (via `origin-simplemap.safecast.org`). It uses **specific** location blocks for each MCP endpoint — a broad `/api/` rule would break map-server routes like `/api/auth/` and `/api/spectrum/`.

```nginx
server {
    listen 443 ssl http2;
    server_name origin-simplemap.safecast.org;

    # MCP protocol transports
    location /mcp-http { proxy_pass http://localhost:3333; ... }
    location /mcp      { proxy_pass http://localhost:3333; proxy_http_version 1.1; ... }

    # Swagger UI
    location /mcp-api/ { proxy_pass http://localhost:3333/mcp-api/; ... }

    # Web-chat assistant
    location /assistant/ { proxy_pass http://localhost:3334/; ... }

    # Map server auth/user/admin (must come before catch-all)
    location /api/auth/  { proxy_pass http://localhost:8765; ... }
    location /api/user/  { proxy_pass http://localhost:8765; ... }
    location /api/admin/ { proxy_pass http://localhost:8765; client_max_body_size 100M; ... }

    # MCP REST API — specific endpoints only
    location = /api/radiation { proxy_pass http://localhost:3333/api/radiation; ... }
    location = /api/area      { proxy_pass http://localhost:3333/api/area; ... }
    location = /api/sensors   { proxy_pass http://localhost:3333/api/sensors; ... }
    location /api/sensor/     { proxy_pass http://localhost:3333/api/sensor/; ... }
    location /api/device/     { proxy_pass http://localhost:3333/api/device/; ... }
    location = /api/spectra   { proxy_pass http://localhost:3333/api/spectra; ... }
    location = /api/stats     { proxy_pass http://localhost:3333/api/stats; ... }
    location = /api/extreme   { proxy_pass http://localhost:3333/api/extreme; ... }
    location /api/info/       { proxy_pass http://localhost:3333/api/info/; ... }
    location /api/gpt/        { proxy_pass http://localhost:3333/api/gpt/; ... }
    location /api/track/      { proxy_pass http://localhost:3333/api/track/; ... }

    # Map server — everything else (including /api/spectrum/, /api/markers/, etc.)
    location / {
        proxy_pass http://localhost:8765;
        client_max_body_size 100M;
        ...
    }
}
```

### 3. CloudFront Distribution

**Distribution ID:** E12FYIQ8RRXOJ1

**Origin Settings:**
- Origin Domain: `origin-simplemap.safecast.org`
- Protocol: HTTPS only
- HTTPS Port: 443
- HTTP Port: 80
- Origin Protocol Policy: `https-only`
- Custom Headers: None (uses default Host header)

**Cache Behavior:**
- Allowed Methods: GET, HEAD, OPTIONS, PUT, POST, PATCH, DELETE
- Viewer Protocol Policy: Redirect HTTP to HTTPS
- Origin Request Policy: AllViewer
- Compress: Yes

**Domain:**
- CNAME: `simplemap.safecast.org`
- SSL Certificate: ACM certificate in us-east-1

## Testing

Test the setup:
```bash
# MCP Server
curl -X POST https://simplemap.safecast.org/mcp-http \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}'

# Map Page
curl -I https://simplemap.safecast.org/

# Map API
curl https://simplemap.safecast.org/api/tracks?limit=1
```

## Connecting to MCP Server

### Claude Code CLI
```bash
claude mcp add --transport http safecast https://simplemap.safecast.org/mcp-http
```

### Claude.ai Web Interface
1. Settings → Integrations → Add custom integration
2. URL: `https://simplemap.safecast.org/mcp-http`

## Troubleshooting

### CloudFront returns 502 errors
**Cause:** SSL certificate doesn't include the origin domain

**Fix:** Ensure certificate includes both domains:
```bash
certbot certonly --nginx \
  -d simplemap.safecast.org \
  -d origin-simplemap.safecast.org \
  --expand
```

### MCP requests go to map server
**Cause:** Nginx not matching the correct server_name

**Fix:** Check nginx server_name includes both domains:
```bash
grep server_name /etc/nginx/sites-enabled/simplemap.safecast.org
# Should show: server_name simplemap.safecast.org origin-simplemap.safecast.org;
```

### Certificate renewal fails
**Cause:** CloudFront is intercepting Let's Encrypt validation

**Solution:** Temporarily change CloudFront origin to port 8765 (direct to map server):
```bash
# Update CloudFront
aws cloudfront get-distribution-config --id E12FYIQ8RRXOJ1 > /tmp/cf.json
# Edit: CustomOriginConfig.HTTPPort = 8765, OriginProtocolPolicy = "http-only"
# Apply changes, wait for deployment
# Run certbot
# Revert CloudFront back to port 443, https-only
```

## Important Notes

- **DNS Records:**
  - `simplemap.safecast.org` → CloudFront (ALIAS record)
  - `origin-simplemap.safecast.org` → 65.108.24.131 (A record)

- **SSH/Rsync:** Always use IP address (65.108.24.131), not domain name

- **CloudFront Invalidation:** After deploying new code:
  ```bash
  aws cloudfront create-invalidation --distribution-id E12FYIQ8RRXOJ1 --paths "/*"
  ```

- **Swagger Documentation (source of truth):** The MCP server's API docs at `/mcp-api/` are generated from Swaggo annotations in `cmd/mcp-server/rest.go` and `cmd/mcp-server/rest_*.go`. The map/unified API docs at `/map-api/` are generated from annotations in `cmd/unified-server/`, `pkg/httpapi/`, and `pkg/auth/`.
  - Host: `simplemap.safecast.org`
  - Base Path: `/api`

  If these values need updating, edit `cmd/mcp-server/rest.go` annotations, regenerate with `swag init -g rest.go` from inside `cmd/mcp-server/`, rebuild, and invalidate CloudFront cache for `/mcp-api/*`.

- **Nginx routing pitfall:** Do NOT use a broad `location /api/` rule pointing to port 3333. This breaks map-server routes (`/api/auth/`, `/api/spectrum/`, etc.). Always use specific `location = /api/endpoint` or `location /api/prefix/` rules for MCP endpoints.

## Related Documentation

- [CloudFront Setup Guide](cloudfront-setup.md)
- [SSH Subdomain Setup](ssh-subdomain-setup.md)
- [Deployment Guide](DEPLOYMENT.md)
