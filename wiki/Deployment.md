# Deployment Guide

Complete guide to deploying Safecast New Map in production.

[[Home|← Back to Home]]

---

## Overview

This guide covers deployment options from simple binary to production-grade infrastructure with HTTPS, CDN, and automated updates.

---

## Deployment Options

### Quick Comparison

| Method | Complexity | Best For | HTTPS | Auto-Update |
|--------|-----------|----------|-------|-------------|
| **Binary** | ⭐ | Testing, personal use | Manual | No |
| **Docker** | ⭐⭐ | Simple production | Via proxy | No |
| **Let's Encrypt** | ⭐⭐ | Production | ✅ Automatic | Optional |
| **CloudFront/CDN** | ⭐⭐⭐ | High-traffic production | ✅ Automatic | ✅ Yes |

---

## Binary Deployment

### Download Binary

```bash
# Download from releases
wget https://github.com/safecast/safecast-new-map/releases/latest/download/safecast-new-map-linux-amd64

chmod +x safecast-new-map-linux-amd64
```

### Run with Systemd

**Create service file:**
```bash
sudo tee /etc/systemd/system/safecast.service << 'EOF'
[Unit]
Description=Safecast New Map Server
After=network.target postgresql.service

[Service]
Type=simple
User=safecast
WorkingDirectory=/opt/safecast
ExecStart=/opt/safecast/safecast-new-map \
  -db-type pgx \
  -db-conn "postgres://user:pass@localhost/safecast" \
  -admin-password "your-admin-password"
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
```

**Enable and start:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable safecast
sudo systemctl start safecast
sudo systemctl status safecast
```

### Run with Supervisor

**Install supervisor:**
```bash
sudo apt install supervisor
```

**Create config:**
```bash
sudo tee /etc/supervisor/conf.d/safecast.conf << 'EOF'
[program:safecast]
command=/opt/safecast/safecast-new-map -db-type pgx -db-conn "postgres://user:pass@localhost/safecast"
directory=/opt/safecast
user=safecast
autostart=true
autorestart=true
stderr_logfile=/var/log/safecast/err.log
stdout_logfile=/var/log/safecast/out.log
EOF
```

**Start:**
```bash
sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl start safecast
```

---

## Docker Deployment

### Basic Docker Run

```bash
docker run -d \
  -p 8765:8765 \
  --name safecast-map \
  safecastr/safecast-new-map:latest
```

### Docker with Volumes

```bash
docker run -d \
  -p 8765:8765 \
  -v /path/to/data:/data \
  -v /path/to/config:/config \
  --name safecast-map \
  safecastr/safecast-new-map:latest \
  -db-type sqlite \
  -db-path /data/safecast.db
```

### Docker Compose

**Create `docker-compose.yml`:**
```yaml
version: '3.8'

services:
  safecast:
    image: safecastr/safecast-new-map:latest
    ports:
      - "8765:8765"
    volumes:
      - ./data:/data
    environment:
      - DATABASE_URL=postgres://user:pass@postgres:5432/safecast
      - SESSION_SECRET=your-secret-key
    depends_on:
      - postgres
    restart: unless-stopped

  postgres:
    image: postgres:16
    environment:
      - POSTGRES_USER=safecast_user
      - POSTGRES_PASSWORD=your-password
      - POSTGRES_DB=safecast
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    restart: unless-stopped

volumes:
  postgres_data:
```

**Initialize PostGIS:**
```sql
-- init.sql
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;
```

**Start:**
```bash
docker-compose up -d
```

### Docker with Environment Variables

```bash
docker run -d \
  -p 8765:8765 \
  -e DATABASE_URL="postgres://user:pass@localhost/safecast" \
  -e SESSION_SECRET="your-secret-key" \
  -e ANTHROPIC_API_KEY="your-api-key" \
  --name safecast-map \
  safecastr/safecast-new-map:latest
```

---

## HTTPS with Let's Encrypt

### Automatic HTTPS

```bash
./safecast-new-map \
  -domain maps.example.org \
  -db-type pgx \
  -db-conn "postgres://user:pass@localhost/safecast"
```

**Requirements:**
- Ports 80 and 443 must be open
- Domain must point to server IP
- DNS A record configured

### How It Works

1. Server starts on port 80 and 443
2. Let's Encrypt validates domain ownership via HTTP-01 challenge
3. Certificate automatically obtained and renewed
4. HTTP redirects to HTTPS

### Firewall Configuration

**Ubuntu/Debian (UFW):**
```bash
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8765/tcp  # If keeping HTTP
```

**RHEL/CentOS (firewalld):**
```bash
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### Certificate Storage

Certificates stored in:
- `./.autocert/` (default)
- Or custom path via configuration

---

## Reverse Proxy Setup

### Nginx Configuration

**Basic proxy:**
```nginx
server {
    listen 80;
    server_name maps.example.org;

    location / {
        proxy_pass http://localhost:8765;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support (if needed)
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

**With SSL:**
```nginx
server {
    listen 443 ssl http2;
    server_name maps.example.org;

    ssl_certificate /etc/letsencrypt/live/maps.example.org/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/maps.example.org/privkey.pem;

    location / {
        proxy_pass http://localhost:8765;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name maps.example.org;
    return 301 https://$host$request_uri;
}
```

### Apache Configuration

```apache
<VirtualHost *:443>
    ServerName maps.example.org

    SSLEngine on
    SSLCertificateFile /etc/letsencrypt/live/maps.example.org/fullchain.pem
    SSLCertificateKeyFile /etc/letsencrypt/live/maps.example.org/privkey.pem

    ProxyPreserveHost On
    ProxyPass / http://localhost:8765/
    ProxyPassReverse / http://localhost:8765/
</VirtualHost>

<VirtualHost *:80>
    ServerName maps.example.org
    Redirect permanent / https://maps.example.org/
</VirtualHost>
```

### Caddy Configuration

```caddy
maps.example.org {
    reverse_proxy localhost:8765
}
```

Caddy automatically obtains and renews Let's Encrypt certificates.

---

## CloudFront/CDN Deployment

### Architecture

```
Internet
    ↓
CloudFront Distribution
    ↓
AWS WAF (Web Application Firewall)
    ↓
Origin Server (Your EC2 instance)
    ↓
Safecast New Map (port 8765)
```

### CloudFront Setup

**1. Create CloudFront Distribution:**
- Origin: Your server IP or hostname
- Origin protocol: HTTP (if SSL terminated at CloudFront)
- Viewer protocol: HTTPS only
- SSL certificate: Use ACM certificate for your domain

**2. Configure Behaviors:**
- Cache static assets (JS, CSS, images)
- Forward API requests to origin
- Cache JSON archives
- Bypass cache for admin endpoints

**3. Configure WAF:**
- SQL injection protection
- XSS protection
- Rate limiting
- Body size limits (for uploads)

See [cloudfront-setup.md](/docs/cloudfront-setup.md) for detailed setup.

### Important: Deployment with CloudFront

**Use IP Address for SSH/Rsync:**

When CloudFront is in front, you must use the server's IP address directly:

```bash
# ✅ Correct - Use IP
ssh root@65.108.24.131
rsync -avz ./binary root@65.108.24.131:/opt/safecast/

# ❌ Wrong - Domain points to CloudFront
ssh root@maps.example.org
```

See [GITHUB_ACTIONS_GUIDE.md](/GITHUB_ACTIONS_GUIDE.md) for CI/CD setup.

### CloudFront Caching

**Cache invalidation after deployment:**
```bash
aws cloudfront create-invalidation \
  --distribution-id E12FYIQ8RRXOJ1 \
  --paths "/*"
```

**Configure cache behaviors:**
- Static assets: Long TTL (1 day)
- API responses: Short TTL or no cache
- Admin endpoints: No cache
- JSON archives: Medium TTL (1 hour)

See [CLOUDFRONT_MCP_TROUBLESHOOTING.md](/docs/CLOUDFRONT_MCP_TROUBLESHOOTING.md) for troubleshooting.

---

## Self-Upgrade System

The platform includes an automatic self-upgrade system.

### Enable Self-Upgrade

```bash
./safecast-new-map \
  -selfupgrade \
  -domain maps.example.org
```

### How It Works

1. **Check for updates** - Polls GitHub releases
2. **Download binary** - Downloads new release
3. **Backup database** - Automatic backup before upgrade
4. **Canary test** - Tests new binary on alternate port
5. **Switch binary** - Replaces old binary
6. **Restart service** - Restarts with new version

### Canary Testing

Before switching:
1. New binary starts on alternate port (e.g., 8766)
2. Health check performed
3. If healthy, replaces current binary
4. If unhealthy, rollback to previous version

### Rollback

If upgrade fails:
1. Old binary restored automatically
2. Database backup available
3. Service restarts with previous version

### Manual Upgrade

```bash
# Stop service
sudo systemctl stop safecast

# Backup
cp safecast-new-map safecast-new-map.backup

# Download new binary
wget https://github.com/safecast/safecast-new-map/releases/latest/download/safecast-new-map-linux-amd64
mv safecast-new-map-linux-amd64 safecast-new-map
chmod +x safecast-new-map

# Start service
sudo systemctl start safecast

# Verify
curl http://localhost:8765/api/stats
```

---

## Production Infrastructure Example

### Current Production Setup

**Server:**
- IP: 65.108.24.131
- Domain: simplemap.safecast.org
- OS: Linux (Ubuntu/Debian)

**Services:**
- Map server (port 8765)
- MCP server (port 3333)
- Web chat (port 3334)
- All served via unified server

**Infrastructure:**
- CloudFront distribution ID: E12FYIQ8RRXOJ1
- AWS WAF with SQL injection/XSS protection
- PostgreSQL database
- DuckLake for analytics

**Deployment:**
- GitHub Actions CI/CD
- SSH/rsync deployment (via IP)
- CloudFront cache invalidation after deploy

See [DEPLOYMENT.md](/docs/DEPLOYMENT.md) for complete deployment procedures.

---

## Database Backup

### PostgreSQL Backup Script

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backups/safecast"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/safecast_$DATE.sql.gz"

mkdir -p $BACKUP_DIR

# Backup database
pg_dump -h localhost -U safecast_user safecast | gzip > $BACKUP_FILE

# Keep 30 days
find $BACKUP_DIR -name "safecast_*.sql.gz" -mtime +30 -delete

echo "Backup completed: $BACKUP_FILE"
```

**Cron job:**
```bash
# Daily backup at 2 AM
0 2 * * * /opt/safecast/backup.sh >> /var/log/safecast-backup.log 2>&1
```

### SQLite/DuckDB Backup

```bash
# Copy database file
cp /data/safecast.db /backups/safecast_$(date +%Y%m%d).db

# Or use .backup command
sqlite3 /data/safecast.db ".backup '/backups/safecast_$(date +%Y%m%d).db'"
```

---

## Monitoring

### Health Check Endpoint

```bash
curl http://localhost:8765/api/stats
```

**Response:**
```json
{
  "totalMarkers": 518400000,
  "totalTracks": 45000,
  "status": "healthy"
}
```

### Systemd Monitoring

```bash
# Check status
systemctl status safecast

# View logs
journalctl -u safecast -f

# Check resource usage
systemd-cgtop
```

### Process Monitoring

**Supervisor:**
```bash
supervisorctl status safecast
supervisorctl tail safecast
```

**PM2 (Node.js style):**
```bash
pm2 start safecast-new-map --name safecast
pm2 monit
```

### Log Monitoring

```bash
# Follow logs
tail -f /var/log/safecast.log

# Filter errors
grep "ERROR" /var/log/safecast.log

# Filter auth events
grep "AUTH:" /var/log/safecast.log
```

### Metrics Collection

**Prometheus exporter:**
- Track request counts
- Response times
- Database query times
- Cache hit rates

**Grafana dashboards:**
- Server health
- API performance
- Database metrics
- User activity

---

## Security Hardening

### Firewall

```bash
# Allow only necessary ports
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP (Let's Encrypt)
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

### SSL/TLS Configuration

**Nginx:**
```nginx
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers HIGH:!aNULL:!MD5;
ssl_prefer_server_ciphers on;
```

### Rate Limiting

**Nginx:**
```nginx
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;

location /api/ {
    limit_req zone=api burst=20 nodelay;
    proxy_pass http://localhost:8765;
}
```

### Database Security

**PostgreSQL:**
```sql
-- Restrict connections
ALTER SYSTEM SET listen_addresses = 'localhost';

-- Use strong passwords
ALTER USER safecast_user WITH PASSWORD 'strong-password';

-- Restrict access via pg_hba.conf
```

### Environment Variables

Store secrets in environment variables:
```bash
export DATABASE_URL="postgres://user:pass@localhost/db"
export SESSION_SECRET="random-secret-key"
export ANTHROPIC_API_KEY="api-key"
```

Use `.env` file with proper permissions:
```bash
chmod 600 .env
```

---

## Performance Tuning

### PostgreSQL Tuning

```sql
-- Increase shared_buffers (25% of RAM)
ALTER SYSTEM SET shared_buffers = '4GB';

-- Increase work_mem for complex queries
ALTER SYSTEM SET work_mem = '256MB';

-- Increase maintenance_work_mem
ALTER SYSTEM SET maintenance_work_mem = '1GB';

-- Set max_connections
ALTER SYSTEM SET max_connections = 200;

-- Reload
SELECT pg_reload_conf();
```

### Connection Pooling

**PgBouncer:**
```bash
# Install
sudo apt install pgbouncer

# Configure
sudo tee /etc/pgbouncer/pgbouncer.ini << 'EOF'
[databases]
safecast = host=localhost dbname=safecast

[pgbouncer]
listen_port = 6432
listen_addr = localhost
auth_type = md5
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 20
EOF
```

### Application Tuning

**Increase cache size:**
```bash
# Set LRU cache size (in code or config)
```

**Optimize database queries:**
- Use indexes
- Avoid N+1 queries
- Batch operations

---

## High Availability

### Database Replication

**PostgreSQL streaming replication:**
- Primary database
- Standby replica(s)
- Automatic failover

**PgPool-II:**
- Connection pooling
- Load balancing
- Automatic failover

### Application Clustering

**Multiple instances:**
- Load balancer (Nginx, HAProxy)
- Shared database
- Shared cache (Redis)

**Sticky sessions:**
- Required for session-based auth
- Configure load balancer

### Backup Strategy

**Daily:**
- Full database backup
- Configuration backup

**Hourly:**
- WAL archiving (PostgreSQL)
- Incremental backup

**Offsite:**
- S3 backup
- Remote server backup

---

## Troubleshooting

### Service Won't Start

**Check logs:**
```bash
journalctl -u safecast -n 100
tail -100 /var/log/safecast.log
```

**Test binary:**
```bash
./safecast-new-map --help
```

**Check database connection:**
```bash
psql -h localhost -U safecast_user -d safecast -c "SELECT 1;"
```

### High Memory Usage

**Check process:**
```bash
ps aux | grep safecast
top -p $(pgrep safecast)
```

**Reduce connection pool:**
```bash
# Adjust in configuration
```

**Restart service:**
```bash
sudo systemctl restart safecast
```

### SSL Certificate Issues

**Check certificate:**
```bash
openssl s_client -connect maps.example.org:443 -servername maps.example.org
```

**Renew Let's Encrypt:**
```bash
sudo certbot renew
```

### CloudFront Issues

**Check distribution:**
```bash
aws cloudfront get-distribution --id E12FYIQ8RRXOJ1
```

**Invalidate cache:**
```bash
aws cloudfront create-invalidation --distribution-id E12FYIQ8RRXOJ1 --paths "/*"
```

See [cloudfront-fix-upload-403.md](/docs/cloudfront-fix-upload-403.md) for upload issues.
See [cloudfront-fix-waf-403.md](/docs/cloudfront-fix-waf-403.md) for WAF issues.

---

## See Also

- [GitHub Actions Guide](/GITHUB_ACTIONS_GUIDE.md) - CI/CD automation
- [Deployment Guide](/docs/DEPLOYMENT.md) - Detailed deployment procedures
- [CloudFront Setup](/docs/cloudfront-setup.md) - CDN configuration
- [CloudFront MCP Troubleshooting](/docs/CLOUDFRONT_MCP_TROUBLESHOOTING.md) - MCP with CloudFront
- [Configuration Reference](Configuration-Reference) - All configuration options
- [Database Setup](Database-Setup) - Database configuration and tuning
