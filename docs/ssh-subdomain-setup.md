# SSH Subdomain Setup for simplemap.safecast.org

## Overview

Create a separate SSH-specific subdomain that bypasses CloudFront and points directly to the server IP.

## Why?

- **CloudFront** handles web traffic (HTTP/HTTPS) on `simplemap.safecast.org`
- **SSH subdomain** handles SSH traffic (port 22) on `ssh.simplemap.safecast.org`
- Both point to the same server (65.108.24.131), but SSH bypasses CloudFront

## Setup in Route53

### Step 1: Add A Record

1. Go to Route53: https://console.aws.amazon.com/route53/
2. Click on the `safecast.org` hosted zone
3. Click "Create record"
4. Configure:
   - **Record name:** `ssh.simplemap`
   - **Record type:** `A - IPv4 address`
   - **Value:** `65.108.24.131`
   - **TTL:** `300` (5 minutes)
   - **Routing policy:** Simple routing
5. Click "Create records"

### Step 2: Verify DNS Propagation

Wait a few minutes, then test:

```bash
# Check DNS resolution
dig ssh.simplemap.safecast.org

# Should show:
# ssh.simplemap.safecast.org. 300 IN A 65.108.24.131

# Test SSH connection
ssh root@ssh.simplemap.safecast.org
```

## Usage

### For Humans (Manual SSH)

Instead of:
```bash
ssh root@65.108.24.131  # Hard to remember
```

Use:
```bash
ssh root@ssh.simplemap.safecast.org  # Easy to remember!
```

### For GitHub Actions

Update the `MAP_SERVER_HOST` secret to use the subdomain:

```yaml
# Option 1: Use subdomain (cleaner)
MAP_SERVER_HOST: ssh.simplemap.safecast.org

# Option 2: Keep using IP (more reliable)
MAP_SERVER_HOST: 65.108.24.131
```

**Recommendation:** Keep using the IP in GitHub Actions for maximum reliability. DNS resolution failures won't break deployments.

### For SSH Config

Add to `~/.ssh/config`:

```bash
Host simplemap
    HostName ssh.simplemap.safecast.org
    User root
    IdentityFile ~/.ssh/safecast-deploy

Host simplemap-mcp
    HostName ssh.simplemap.safecast.org
    User root
    IdentityFile ~/.ssh/safecast-mcp-deploy
```

Then:
```bash
ssh simplemap
scp file.txt simplemap:/path/
rsync -avz ./dir/ simplemap:/remote/
```

## Traffic Flow Comparison

### Before (Using IP Address)
```
Developer → 65.108.24.131:22 → SSH Server
```

### After (Using SSH Subdomain)
```
Developer → ssh.simplemap.safecast.org (DNS lookup)
         → 65.108.24.131:22 → SSH Server
```

### Web Traffic (Unchanged)
```
User → simplemap.safecast.org (DNS → CloudFront)
    → CloudFront → Apache → Map Server
```

## Benefits

- ✅ **Memorable hostname** for SSH access
- ✅ **Free** (no additional AWS costs)
- ✅ **Flexible** - Can change IP without updating everywhere
- ✅ **Separate concerns** - Web traffic vs SSH traffic use different hostnames
- ✅ **No CloudFront bypass needed** - Clean separation

## Considerations

### If Server IP Changes

With subdomain, you only update DNS once:
```bash
# Update the A record in Route53
ssh.simplemap.safecast.org → NEW_IP

# All users automatically use new IP (after TTL expires)
```

Without subdomain, you must update:
- GitHub Actions secrets (both repos)
- Local SSH configs (every developer)
- Documentation
- Any scripts

### DNS Dependency

- SSH subdomain requires DNS to be working
- If Route53 is down, SSH won't work
- IP address always works, regardless of DNS

**For critical deployments:** Keep using the IP in GitHub Actions for reliability.

## Rollback

If you want to remove the subdomain:

1. Delete the DNS record in Route53
2. Update any configs that use `ssh.simplemap.safecast.org` back to the IP

## Related Documentation

- [Main CloudFront Setup](cloudfront-setup.md)
- [Deployment Guide](DEPLOYMENT.md)
- [SSH Config Examples](https://man.openbsd.org/ssh_config)
