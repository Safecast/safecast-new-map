# Cloudflare CNAME Setup (Without Domain Registrar Access)

## This option requires Cloudflare Business plan ($200/month) or using Cloudflare Partner

Since you don't have domain registrar access, standard Cloudflare setup won't work.
You have two options:

### Option A: Cloudflare for SaaS (Requires Business Plan - $200/month)
Not recommended for cost reasons.

### Option B: Use Cloudflare Pages/Workers with Custom Domain

1. **Create Cloudflare Account**
2. **Set up Cloudflare Workers/Pages as a proxy**
3. **Point CNAME to Workers domain**

This is more complex and not worth it compared to CloudFront.

## Recommendation

**Use AWS CloudFront instead** (see cloudfront-setup.md) because:
- ✅ You already have AWS access
- ✅ Free SSL certificates
- ✅ No additional account needed
- ✅ Edge locations in Tokyo/Osaka for Japan traffic
- ✅ Can set up entirely in Route53 + CloudFront
- ✅ Similar performance to Cloudflare
- ❌ Cloudflare CNAME setup requires Business plan or complex workarounds
