# CloudFront + MCP Server Setup

## Overview

The MCP server and map server both run behind CloudFront, with nginx routing requests to the correct backend based on the URL path.

## Architecture

```
User Request
    ↓
CloudFront (simplemap.safecast.org)
    ↓ HTTPS (port 443)
origin-simplemap.safecast.org
    ↓
Nginx (routes by path)
    ├─ /mcp-http → MCP Server (port 3333)
    ├─ /mcp      → MCP Server (port 3333)
    ├─ /docs/    → MCP Server (port 3333)
    ├─ /api/mcp/ → MCP Server (port 3333)
    └─ /         → Map Server (port 8765)
```

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

**Server names:** `/etc/nginx/sites-enabled/simplemap.safecast.org`
```nginx
server {
    server_name simplemap.safecast.org origin-simplemap.safecast.org;

    # MCP Server endpoints
    location /mcp-http {
        proxy_pass http://localhost:3333;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /mcp {
        proxy_pass http://localhost:3333;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        # ... other headers
    }

    location /docs/ {
        proxy_pass http://localhost:3333/docs/;
        # ... headers
    }

    location /api/mcp/ {
        proxy_pass http://localhost:3333/api/;
        # ... headers
    }

    # Map server (catch-all)
    location / {
        proxy_pass http://localhost:8765;
        # ... headers
    }

    listen [::]:443 ssl ipv6only=on;
    listen 443 ssl;
    ssl_certificate /etc/letsencrypt/live/simplemap.safecast.org/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/simplemap.safecast.org/privkey.pem;
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

- **Swagger Documentation:** The MCP server's Swagger docs at `/docs/` are configured with:
  - Host: `simplemap.safecast.org`
  - Base Path: `/api/mcp`

  If these values need updating, edit `go/cmd/mcp-server/rest.go` annotations, regenerate with `swag init -g rest.go`, rebuild the binary, and invalidate CloudFront cache for `/docs/*`

## Related Documentation

- [CloudFront Setup Guide](cloudfront-setup.md)
- [SSH Subdomain Setup](ssh-subdomain-setup.md)
- [MCP Server README](../../safecast-map-MCP/README.md)
