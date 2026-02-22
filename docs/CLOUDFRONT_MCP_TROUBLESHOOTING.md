# CloudFront + MCP Server Troubleshooting Guide

## Problem: MCP Server Not Working After CloudFront Setup

### Symptoms
- MCP server accessible at `vps01.safecast.org` but not at `simplemap.safecast.org/mcp-http`
- Requests return map server HTML instead of MCP JSON responses
- CloudFront returns 502 Bad Gateway errors

### Root Causes & Solutions

#### Issue 1: CloudFront Connecting to Wrong Port

**Problem:** CloudFront was configured to connect directly to port 8765 (map server), bypassing nginx entirely.

**Diagnosis:**
```bash
aws cloudfront list-distributions --query 'DistributionList.Items[?Aliases.Items[?contains(@, `simplemap.safecast.org`)]].[Id,Origins.Items[0].DomainName,Origins.Items[0].CustomOriginConfig.HTTPPort]' --output table
```

**Solution:** Change CloudFront origin to use port 443 (HTTPS) so it goes through nginx:
```json
{
  "HTTPPort": 80,
  "HTTPSPort": 443,
  "OriginProtocolPolicy": "https-only"
}
```

#### Issue 2: SSL Certificate Mismatch

**Problem:** CloudFront connects to `origin-simplemap.safecast.org`, but SSL certificate only included `simplemap.safecast.org`.

**Diagnosis:**
```bash
openssl s_client -connect origin-simplemap.safecast.org:443 -servername origin-simplemap.safecast.org < /dev/null 2>&1 | grep -A 2 "Subject Alternative Name"
```

**Solution:** Add origin subdomain to SSL certificate:
```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "certbot certonly --nginx -d simplemap.safecast.org -d origin-simplemap.safecast.org --expand"
```

**Note:** If CloudFront is already intercepting requests, temporarily revert CloudFront to port 8765 first.

#### Issue 3: Nginx Not Matching Origin Domain

**Problem:** Nginx `server_name` only included `simplemap.safecast.org`, but CloudFront sends `Host: origin-simplemap.safecast.org`.

**Diagnosis:**
```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "grep server_name /etc/nginx/sites-enabled/simplemap.safecast.org"
```

**Solution:** Add origin subdomain to nginx server_name:
```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131 "sed -i 's/server_name simplemap.safecast.org;/server_name simplemap.safecast.org origin-simplemap.safecast.org;/' /etc/nginx/sites-enabled/simplemap.safecast.org && nginx -t && systemctl reload nginx"
```

#### Issue 4: CloudFront Custom Host Header

**Problem:** AWS CloudFront doesn't allow custom "Host" headers in origin settings.

**Error:**
```
An error occurred (InvalidArgument) when calling the UpdateDistribution operation:
The parameter HeaderName : Host is not allowed.
```

**Solution:** Remove custom headers and let CloudFront send the origin domain name as the Host header naturally:
```json
{
  "CustomHeaders": {
    "Quantity": 0,
    "Items": []
  }
}
```

## Step-by-Step Resolution Process

### 1. Verify Current Setup
```bash
# Check CloudFront origin config
aws cloudfront get-distribution-config --id E12FYIQ8RRXOJ1 | jq '.DistributionConfig.Origins.Items[0]'

# Check SSL certificate
ssh root@65.108.24.131 "openssl x509 -in /etc/letsencrypt/live/simplemap.safecast.org/fullchain.pem -text -noout | grep -A 2 'Subject Alternative Name'"

# Check nginx server_name
ssh root@65.108.24.131 "grep server_name /etc/nginx/sites-enabled/simplemap.safecast.org"
```

### 2. Add Origin Domain to SSL Certificate
```bash
# If CloudFront is already active, temporarily revert to port 8765
aws cloudfront get-distribution-config --id E12FYIQ8RRXOJ1 > /tmp/cf.json
# Edit config to use port 8765, http-only
# Apply and wait for deployment

# Add origin domain to certificate
ssh root@65.108.24.131 "certbot certonly --nginx -d simplemap.safecast.org -d origin-simplemap.safecast.org --expand"

# Reload nginx
ssh root@65.108.24.131 "systemctl reload nginx"
```

### 3. Update Nginx Configuration
```bash
ssh root@65.108.24.131 "sed -i 's/server_name simplemap.safecast.org;/server_name simplemap.safecast.org origin-simplemap.safecast.org;/' /etc/nginx/sites-enabled/simplemap.safecast.org && nginx -t && systemctl reload nginx"
```

### 4. Update CloudFront Origin Settings
```bash
aws cloudfront get-distribution-config --id E12FYIQ8RRXOJ1 > /tmp/cf-config.json

# Edit /tmp/cf-config.json:
# - CustomOriginConfig.HTTPPort = 80
# - CustomOriginConfig.HTTPSPort = 443
# - CustomOriginConfig.OriginProtocolPolicy = "https-only"
# - CustomHeaders.Quantity = 0, Items = []

ETAG=$(jq -r '.ETag' /tmp/cf-config.json)
jq '.DistributionConfig' /tmp/cf-config.json > /tmp/cf-dist.json

aws cloudfront update-distribution \
  --id E12FYIQ8RRXOJ1 \
  --distribution-config file:///tmp/cf-dist.json \
  --if-match "$ETAG"
```

### 5. Wait and Test
```bash
# Wait for deployment (5-15 minutes)
aws cloudfront get-distribution --id E12FYIQ8RRXOJ1 --query 'Distribution.Status'

# Test MCP server
curl -X POST https://simplemap.safecast.org/mcp-http \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | jq .

# Test map
curl -I https://simplemap.safecast.org/
```

## Verification Checklist

- [ ] SSL certificate includes both `simplemap.safecast.org` and `origin-simplemap.safecast.org`
- [ ] Nginx `server_name` includes both domains
- [ ] CloudFront origin uses HTTPS protocol policy
- [ ] CloudFront origin port is 443 (not 8765)
- [ ] CloudFront has no custom Host headers
- [ ] CloudFront deployment status is "Deployed"
- [ ] MCP initialize request returns success
- [ ] Map loads correctly in browser

## Common Errors

### "The plain HTTP request was sent to HTTPS port"
**Cause:** CloudFront sending HTTP to port 443
**Fix:** Set `OriginProtocolPolicy: "https-only"`

### "Invalid session ID" from MCP server
**Status:** ✅ This is normal! It means routing works.
**Explanation:** MCP requires session management. Claude Code/AI will handle this automatically.

### 502 Bad Gateway
**Causes:**
1. SSL certificate mismatch
2. Origin domain not resolving
3. Nginx not listening on port 443
4. Server firewall blocking CloudFront IPs

**Debug:**
```bash
# Test origin directly
curl -k -I https://origin-simplemap.safecast.org/

# Check nginx is listening
ssh root@65.108.24.131 "ss -tlnp | grep :443"

# Check nginx logs
ssh root@65.108.24.131 "tail -f /var/log/nginx/error.log"
```

## Final Working Configuration

### CloudFront
```json
{
  "DomainName": "origin-simplemap.safecast.org",
  "CustomOriginConfig": {
    "HTTPPort": 80,
    "HTTPSPort": 443,
    "OriginProtocolPolicy": "https-only"
  },
  "CustomHeaders": {
    "Quantity": 0,
    "Items": []
  }
}
```

### Nginx
```nginx
server {
    server_name simplemap.safecast.org origin-simplemap.safecast.org;

    location /mcp-http {
        proxy_pass http://localhost:3333;
        # ... proxy headers
    }

    location / {
        proxy_pass http://localhost:8765;
        # ... proxy headers
    }

    listen 443 ssl;
    ssl_certificate /etc/letsencrypt/live/simplemap.safecast.org/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/simplemap.safecast.org/privkey.pem;
}
```

### SSL Certificate
```
X509v3 Subject Alternative Name:
    DNS:origin-simplemap.safecast.org, DNS:simplemap.safecast.org
```

## Related Documentation

- [CloudFront + MCP Setup Guide](cloudfront-mcp-setup.md)
- [Main CloudFront Setup](cloudfront-setup.md)
- [MCP Server README](../../safecast-map-MCP/README.md)
