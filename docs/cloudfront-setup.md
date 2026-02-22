# AWS CloudFront Setup for simplemap.safecast.org

## Prerequisites
- AWS Console access
- Route53 access

## Step 1: Create CloudFront Distribution

1. **Go to CloudFront in AWS Console**
   - Navigate to: https://console.aws.amazon.com/cloudfront/

2. **Create Distribution**
   - Click "Create Distribution"

3. **Origin Settings:**
   - **Origin Domain:** `simplemap.safecast.org` (or use the IP: `65.108.24.131`)
   - **Protocol:** HTTPS only (or Match viewer)
   - **Origin Path:** Leave blank
   - **Name:** `safecast-simplemap-origin`
   - **Add custom header (IMPORTANT):**
     - Header name: `Host`
     - Value: `simplemap.safecast.org`

4. **Default Cache Behavior:**
   - **Viewer Protocol Policy:** Redirect HTTP to HTTPS
   - **Allowed HTTP Methods:** GET, HEAD, OPTIONS, PUT, POST, PATCH, DELETE
   - **Cache Policy:** CachingOptimized (or create custom)
   - **Origin Request Policy:** AllViewer
   - **Compress Objects Automatically:** Yes

5. **Distribution Settings:**
   - **Alternate Domain Names (CNAMEs):** `simplemap.safecast.org`
   - **Custom SSL Certificate:** Request a certificate (see Step 2)
   - **Supported HTTP Versions:** HTTP/2, HTTP/3
   - **Default Root Object:** Leave blank (your app handles this)
   - **IPv6:** On

6. **Create Distribution** (but don't deploy yet - need SSL cert first)

## Step 2: Request SSL Certificate (ACM)

1. **Go to Certificate Manager**
   - Navigate to: https://console.aws.amazon.com/acm/
   - **IMPORTANT:** Switch region to **US East (N. Virginia) us-east-1** (CloudFront requires this)

2. **Request Certificate:**
   - Click "Request certificate"
   - Choose "Request a public certificate"
   - Domain names: `simplemap.safecast.org`
   - Validation method: **DNS validation**
   - Click "Request"

3. **Add CNAME Record to Route53:**
   - ACM will show a CNAME record to add for validation
   - Click "Create records in Route53" button (it will auto-add)
   - Wait 5-10 minutes for validation to complete
   - Status should change to "Issued"

4. **Go back to CloudFront Distribution:**
   - Edit the distribution
   - Under "Custom SSL Certificate", select the certificate you just created
   - Save changes

## Step 3: Update Route53 DNS

1. **Go to Route53 Hosted Zones**
   - Find `safecast.org` zone

2. **Update simplemap.safecast.org record:**
   - Find the existing `simplemap.safecast.org` A record (currently points to `65.108.24.131`)
   - **Delete the A record**
   - **Create new A record:**
     - Record name: `simplemap`
     - Record type: `A - IPv4 address`
     - **Alias:** Yes (toggle on)
     - **Route traffic to:** Alias to CloudFront distribution
     - **Choose distribution:** Select your CloudFront distribution (e.g., `d111111abcdef8.cloudfront.net`)
     - Routing policy: Simple
     - Create record

3. **Optional: Add AAAA record for IPv6:**
   - Same as above but type: `AAAA - IPv6 address`
   - Alias to same CloudFront distribution

## Step 4: Configure Cache Policies (Optional but Recommended)

### Create Custom Cache Policy for Better Performance:

1. **Go to CloudFront → Policies → Cache**
2. **Create Policy:**
   - Name: `safecast-cache-policy`
   - **Time to Live (TTL):**
     - Min: 1 second
     - Max: 31536000 (1 year)
     - Default: 86400 (1 day)
   - **Cache key settings:**
     - Headers: Include specific headers → Add: `Accept-Encoding`
     - Query strings: All
     - Cookies: None (or All if your app uses cookies)
   - **Compression:** Gzip, Brotli

3. **Update Distribution** to use this cache policy

## Step 5: Verify Setup

After DNS propagates (5-60 minutes):

```bash
# Check DNS points to CloudFront
dig simplemap.safecast.org

# Should see CloudFront domain like: d111111abcdef8.cloudfront.net

# Test with curl
curl -I https://simplemap.safecast.org

# Look for CloudFront headers:
# x-cache: Hit from cloudfront
# x-amz-cf-pop: NRT57-P1 (Tokyo edge location)
# via: 2.0 xxxxx.cloudfront.net (CloudFront)
```

## Step 6: Invalidate Cache When Deploying

When you deploy new code, create an invalidation:

```bash
# Get distribution ID
aws cloudfront list-distributions --query 'DistributionList.Items[*].[Id,Aliases.Items[0]]' --output table

# Create invalidation
aws cloudfront create-invalidation --distribution-id E1234567890ABC --paths "/*"
```

Or via AWS Console:
1. Go to CloudFront → Distributions
2. Select your distribution
3. Invalidations tab → Create Invalidation
4. Object paths: `/*`
5. Create

## Important: SSH/Rsync After CloudFront Setup

**⚠️ CRITICAL:** Once you point your domain to CloudFront, SSH and rsync will **NOT work** with the domain name.

### Why?
- CloudFront only handles HTTP/HTTPS traffic (ports 80/443)
- SSH operates on port 22, which CloudFront doesn't accept or forward
- When you SSH to the domain name, it tries to connect to CloudFront's IPs, not your server

### Solution: Always Use IP Address for Deployment

**✅ Correct - Use IP address:**
```bash
ssh -i ~/.ssh/safecast-deploy root@65.108.24.131
rsync -avP -e "ssh -i ~/.ssh/safecast-deploy" ./safecast-new-map root@65.108.24.131:/usr/local/bin/
```

**❌ Wrong - Domain won't work:**
```bash
ssh -i ~/.ssh/safecast-deploy root@simplemap.safecast.org  # Will fail!
rsync -avP -e "ssh -i ~/.ssh/safecast-deploy" ./safecast-new-map root@simplemap.safecast.org:/usr/local/bin/  # Will fail!
```

### Traffic Flow After CloudFront Setup

```
Web Traffic (HTTP/HTTPS):
User → simplemap.safecast.org (DNS) → CloudFront → Origin (65.108.24.131)

Deployment (SSH/Rsync):
Developer → 65.108.24.131 (Direct IP) → Server
```

## Expected Benefits

- **Reduced latency from Japan:** CloudFront has edge locations in Tokyo, Osaka
- **Improved performance:** Content cached closer to users
- **Better reliability:** DDoS protection, automatic failover
- **Cost:** ~$0.085/GB for Asia traffic (likely $10-50/month depending on usage)

## Monitoring

- CloudFront → Monitoring tab shows:
  - Requests
  - Bytes downloaded
  - Error rates
  - Popular objects
