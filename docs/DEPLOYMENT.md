# Production Deployment Guide

This guide covers deploying to the production server at simplemap.safecast.org.

## Infrastructure Overview

**Domain:** simplemap.safecast.org
**Server IP:** 65.108.24.131
**CDN:** AWS CloudFront (Distribution ID: E12FYIQ8RRXOJ1)

### Services

All three services are built from this repo and deployed by the same GitHub Actions workflow:

| Service | Binary | Port | Location on VPS |
|---------|--------|------|-----------------|
| Map server | `safecast-new-map` | 8765 | `/usr/local/bin/safecast-new-map` |
| MCP server | `safecast-mcp` | 3333 | `/root/safecast-mcp-server/safecast-mcp` |
| Web-chat | `safecast-web-chat` | 3334 | `/root/safecast-web-chat-server/safecast-web-chat` |

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
   go build -o safecast-new-map .
   go build -o safecast-mcp ./cmd/mcp-server/
   go build -o safecast-web-chat ./cmd/web-chat/
   ```
3. **Server access:** Ability to SSH to 65.108.24.131

### Deployment Steps — Map Server

```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl stop safecast-new-map"
rsync -avP -e "ssh -i ~/.ssh/safecast-deploy" ./safecast-new-map root@65.108.24.131:/usr/local/bin/
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl start safecast-new-map && systemctl status safecast-new-map"
```

### Deployment Steps — MCP Server

```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl stop safecast-mcp"
rsync -avP -e "ssh -i ~/.ssh/safecast-deploy" ./safecast-mcp root@65.108.24.131:/root/safecast-mcp-server/safecast-mcp
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl start safecast-mcp && systemctl status safecast-mcp"
```

### Deployment Steps — Web-chat

```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl stop safecast-web-chat"
rsync -avP -e "ssh -i ~/.ssh/safecast-deploy" ./safecast-web-chat root@65.108.24.131:/root/safecast-web-chat-server/safecast-web-chat
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl start safecast-web-chat && systemctl status safecast-web-chat"
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

## Automated Deployment (GitHub Actions)

**Current workflow:** `.github/workflows/deploy.yml`

### How It Works

1. **Trigger:** Automatic on push to `main` branch, or manual via workflow_dispatch
2. **Build:** Compiles Go binary on GitHub runners
3. **Deploy:** Uses SSH with stored private key to deploy to 65.108.24.131
4. **Invalidate:** Clears CloudFront cache to serve new version

### Required GitHub Secrets

| Secret | Purpose |
|--------|---------|
| `DEPLOY_SSH_KEY` | SSH key for deploying to 65.108.24.131 |
| `AWS_ACCESS_KEY_ID` | CloudFront cache invalidation |
| `AWS_SECRET_ACCESS_KEY` | CloudFront cache invalidation |
| `DATABASE_URL` | Postgres connection string for MCP server |
| `ANTHROPIC_API_KEY` | Claude API key for web-chat service |

### Workflow Steps

```
1. Build all 3 binaries (safecast-new-map, safecast-mcp, safecast-web-chat)
2. Setup SSH
3. Stop → rsync → start: safecast-new-map
4. Stop → rsync → start: safecast-mcp  (also writes /root/safecast-mcp-server/.env)
5. Stop → rsync → start: safecast-web-chat  (also writes /root/safecast-web-chat-server/.env)
6. Invalidate CloudFront cache (/*)
7. Cleanup SSH keys
```

**Note:** The workflow correctly uses the IP address (65.108.24.131) for all SSH operations.

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

## Server Configuration

### Service Details

| Service | systemd name | Binary | Config |
|---------|-------------|--------|--------|
| Map server | `safecast-new-map` | `/usr/local/bin/safecast-new-map` | flags in service file |
| MCP server | `safecast-mcp` | `/root/safecast-mcp-server/safecast-mcp` | `/root/safecast-mcp-server/.env` |
| Web-chat | `safecast-web-chat` | `/root/safecast-web-chat-server/safecast-web-chat` | `/root/safecast-web-chat-server/.env` |

### Useful Commands

```bash
# View service status (replace service name as needed)
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl status safecast-new-map"
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl status safecast-mcp"
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl status safecast-web-chat"

# View logs
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "journalctl -u safecast-new-map -f"
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "journalctl -u safecast-mcp -f"

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

## Related Documentation

- [CloudFront Setup Guide](cloudfront-setup.md) - Initial CloudFront configuration
- [Upload 403 Fix](cloudfront-fix-upload-403.md) - Cookie forwarding configuration
- [WAF 403 Fix](cloudfront-fix-waf-403.md) - Large file upload configuration
- [GitHub Actions Guide](../GITHUB_ACTIONS_GUIDE.md) - Workflow details
- [Memory (Project Notes)](~/.claude/projects/-home-rob-Documents-Safecast-safecast-new-map/memory/MEMORY.md)

## Quick Reference

**Server IP:** 65.108.24.131
**SSH Key:** `~/.ssh/safecast-deploy`
**CloudFront Distribution ID:** E12FYIQ8RRXOJ1
**Service Name:** `safecast-new-map`
**Binary Path:** `/usr/local/bin/safecast-new-map`

**One-Line Deploy:**
```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl stop safecast-new-map" && \
rsync -avP -e "ssh -i ~/.ssh/safecast-deploy" ./safecast-new-map root@65.108.24.131:/usr/local/bin/ && \
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "systemctl start safecast-new-map && systemctl status safecast-new-map"
```
