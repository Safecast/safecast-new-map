# CloudFront Fix: Upload 403 Forbidden Error

## Problem

Users get **403 Forbidden** error when uploading multiple large files, even though they're logged in.
Single files and small multiple files work fine.

## Root Cause

CloudFront is **not forwarding session cookies** to the origin server. Without the session cookie,
the server's authentication middleware rejects the upload request with 403.

## Solution

Configure CloudFront to forward cookies properly.

### Option 1: Use Managed Origin Request Policy (Recommended)

1. Go to CloudFront Console: https://console.aws.amazon.com/cloudfront/
2. Select your distribution (simplemap.safecast.org)
3. Go to **Behaviors** tab
4. Edit the **Default (*)** behavior
5. Scroll to **Cache key and origin requests**
6. **Origin request policy:** Select **AllViewerExceptHostHeader** (Managed policy)
   - This forwards all headers, cookies, and query strings except Host
7. Click **Save changes**
8. Wait 5-10 minutes for deployment

### Option 2: Create Custom Origin Request Policy

If you want more control:

1. Go to CloudFront → **Policies** → **Origin request**
2. Click **Create origin request policy**
3. Configure:
   - **Name:** `safecast-origin-request-policy`
   - **Cookies:** Include **All cookies**
   - **Headers:** Include **All viewer headers except Host**
   - **Query strings:** Include **All query strings**
4. Click **Create**
5. Go back to your distribution → Behaviors → Edit default behavior
6. **Origin request policy:** Select `safecast-origin-request-policy`
7. Save changes

### Option 3: Whitelist Specific Session Cookie (Most Efficient)

If you know your session cookie name:

1. Create custom Origin Request Policy (as in Option 2)
2. But for **Cookies**, select **Include specific cookies**
3. Add cookie names:
   - `session` (or whatever your session cookie name is)
   - `safecast_session` (if different)
4. Continue with steps above

## Verification

After CloudFront deployment completes:

1. **Hard refresh** your browser: Ctrl+Shift+R
2. **Log in** to simplemap.safecast.org
3. **Upload multiple large files** (>10 KB each)
4. **Check Network tab** - POST to /upload should return **200 OK**, not 403

## Alternative: Bypass CloudFront for Uploads

If you can't update CloudFront configuration, you can bypass it for uploads:

Add this to your upload form JavaScript to POST directly to the server IP:

```javascript
// Upload directly to server, bypassing CloudFront
const uploadUrl = 'https://65.108.24.131/upload';
// Note: This requires valid SSL certificate on the server
```

However, this is **not recommended** because:
- Loses CloudFront DDoS protection
- Slower uploads from distant locations
- SSL certificate warnings if using self-signed cert

## Related Issues

This same cookie forwarding issue can affect:
- Session-based authentication
- Any user-specific API endpoints
- Form submissions with CSRF tokens

Make sure CloudFront is configured to forward cookies for all authenticated endpoints.
