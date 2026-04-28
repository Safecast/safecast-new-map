# Migration Plan: AWS CloudFront + Route53 → Cloudflare (Free)

## Goal

Replace the current AWS edge stack with Cloudflare Free:

- **DNS:** Cloudflare hosts the `safecast.org` zone; Route53 retired.
- **CDN/edge:** Cloudflare Free fronts `simplemap.safecast.org` → origin `origin-simplemap.safecast.org` (`65.108.24.131`).
- **SSL:** Cloudflare Universal SSL at the edge; existing Let's Encrypt cert at origin covers both names. Mode: **Full (strict)**.
- **CloudFront + ACM + Route53 zone:** decommissioned at the end.

> Prerequisite: registrar access for `safecast.org` (needed to change nameservers and DS records). The earlier `docs/cloudflare-cname-setup.md` was written under the assumption registrar access was unavailable, which is why it recommended sticking with CloudFront. That assumption no longer holds.

## Cost outcome

- **AWS today:** CloudFront egress (~$10–50/mo for Asia traffic) + Route53 hosted zone ($0.50/mo) + ACM (free).
- **Cloudflare after migration:** $0/mo on the Free plan.

Free is sufficient as long as three constraints are acceptable (see "Risks & constraints" below). If any one becomes a real problem, Pro ($25/mo) is a billing flip, not a re-migration.

## Phase 0 — Pre-flight (no user-visible changes)

1. **Inventory the Route53 zone.** Export every record so nothing is missed (especially MX/TXT/SPF/DKIM/DMARC, which can break email if dropped):
   ```
   aws route53 list-resource-record-sets --hosted-zone-id <ZONE_ID> > route53-export.json
   ```
2. **Inventory CloudFront behavior** that must be reproduced: custom `Host` header injection, cache policy, redirect HTTP→HTTPS, allowed methods (GET/HEAD/OPTIONS/PUT/POST/PATCH/DELETE), HTTP/2+3, IPv6.
3. **DNSSEC inventory.** Check whether DNSSEC is currently enabled on `safecast.org`:
   ```
   dig +short DS safecast.org @8.8.8.8
   dig +short DNSKEY safecast.org @8.8.8.8
   aws route53 get-dnssec --hosted-zone-id <ZONE_ID>
   ```
   If DS records exist at the registrar, DNSSEC is live. A botched migration will hard-fail resolution (SERVFAIL) for every validating resolver — not just slow, fully broken. Record the current DS values before changing anything.
4. **Confirm origin SSL cert** on `65.108.24.131` covers both `simplemap.safecast.org` AND `origin-simplemap.safecast.org` and is auto-renewing (so Cloudflare Full-strict works on day one).
5. **Lower Route53 TTLs** to 60–300s on `simplemap` and any other records you'll migrate. Wait at least one current-TTL window before cutover.
6. **Capture baseline behavior** so regressions are detectable post-cutover:
   ```
   curl -sI https://simplemap.safecast.org/api/radiation?...   # baseline latency
   curl -sI https://simplemap.safecast.org/api/track/<id>.csv  # confirm CSV not JSON
   curl -sI https://simplemap.safecast.org/mcp-http            # confirm 200/streaming
   ```
7. **Drop nginx `client_max_body_size`** from `100M` to `95M` on `/` and `/api/admin/` so origin still accepts everything Cloudflare lets through (Cloudflare counts boundaries/headers in the 100 MB cap).

## Phase 1 — Set up the Cloudflare zone (parallel, no traffic yet)

1. Create a free Cloudflare account; **Add Site → `safecast.org`** → Free plan.
2. Cloudflare auto-imports records via scan. **Diff against `route53-export.json`** and add anything missing (especially MX/SPF/DKIM/DMARC, `_acme-challenge` records, internal subdomains).
3. Recreate records:
   - `simplemap` → CNAME `origin-simplemap.safecast.org`, **Proxied (orange cloud)**.
   - `origin-simplemap` → A `65.108.24.131`, **DNS-only (grey cloud)** so SSH/rsync still work.
   - All other A/AAAA/MX/TXT records: copy verbatim, default DNS-only unless they're an HTTP host you want proxied.
4. **DNSSEC prep.** DNS → Settings → DNSSEC → Enable. Cloudflare shows a DS record (algorithm, key tag, digest). **Do not give it to the registrar yet** — that happens during cutover (Phase 3).
5. **SSL/TLS → Overview:** set to **Full (strict)**. Flexible would break HTTPS POSTs; Full-strict matches CloudFront's "HTTPS only" origin behavior.
6. **SSL/TLS → Edge Certificates:** Always Use HTTPS, HSTS (after verifying), TLS 1.2 min, HTTP/2 + HTTP/3, 0-RTT off (safer with POST-heavy app).
7. **Network:** WebSockets ON (MCP HTTP/SSE, chat). gRPC off unless needed.
8. **Rules → Cache Rules.** Bypass set, derived from the live origin nginx config:

   ```
   # MCP server (port 3333) — all dynamic, never cache
   /mcp
   /mcp-http
   /api/radiation
   /api/area
   /api/sensors
   /api/sensor/*
   /api/device/*
   /api/spectra
   /api/stats
   /api/extreme
   /api/info/*
   /api/gpt/*

   # Unified server (port 8765) — auth + admin + uploads
   /api/auth/*
   /api/user/*
   /api/admin/*
   /admin/*

   # Track downloads (CSV/XLSX/JSON, format-by-extension on port 8765)
   /api/track/*

   # Web-chat (port 3334)
   /assistant/*

   # Docs (port 8765, cheap to leave dynamic)
   /docs/*
   /mcp-api/*
   /map-api/*
   ```

   **Cache** for `/static/*`, `*.js`, `*.css`, `*.png`, `*.svg`, `*.woff2` and tile images: edge TTL ~1d, browser TTL ~1h.

9. **Rules → Configuration Rules:** if origin nginx requires a specific `Host` header on any path, set it here. Default behavior (pass through `simplemap.safecast.org`) should match the existing nginx server block, so likely no rule needed.
10. **Security → WAF:** keep managed rules off initially (Free has limited WAF anyway). Bot Fight Mode off until tested — it can break MCP/AI clients.

## Phase 2 — Test before cutover

1. Use Cloudflare's "Test" tab or `curl --resolve` to hit Cloudflare edge IPs while DNS still points at CloudFront:
   ```
   curl --resolve simplemap.safecast.org:443:<cloudflare-edge-ip> https://simplemap.safecast.org/
   ```
   Verify: home page, `/admin/users`, `/api/track/<id>.csv` (must be CSV not JSON), `/mcp-http`, file upload, `/assistant/`, `/docs/`.
2. Validate origin nginx accepts the request when Cloudflare passes `Host: simplemap.safecast.org`. (CloudFront injects this; Cloudflare passes it through by default since the CNAME target is `origin-simplemap`, but worth confirming.)
3. Optional: restrict origin firewall to **Cloudflare IP ranges only** (`https://www.cloudflare.com/ips/`) via iptables — but only after cutover is stable. Keep `origin-simplemap` reachable for SSH/admin.

## Phase 3 — Cutover (DNSSEC-safe order)

The unsafe path is "switch nameservers while old DS records still point at Route53 KSK" — validating resolvers will reject Cloudflare's signatures and SERVFAIL the domain.

If `safecast.org` is currently DNSSEC-signed:

1. **At the registrar: remove the existing DS records** for `safecast.org` (disables DNSSEC at the parent).
2. **Wait for the DS TTL to expire** at the parent zone (typically 24–48h for most TLDs; check with `dig DS safecast.org @<tld-nameserver>`). During this window the domain is unsigned but still resolves normally.
3. **In Route53: disable DNSSEC signing** (`aws route53 disable-hosted-zone-dnssec`) — only after the parent DS is gone.
4. **At the registrar: change nameservers** from Route53 (`ns-*.awsdns-*`) to the two Cloudflare nameservers shown in the dashboard.
5. Wait for Cloudflare to flip the zone to "Active" (minutes to a few hours; depends on TLD). Both providers can answer for a while — that's fine.
6. **At the registrar: add Cloudflare's DS record** (from Phase 1 step 4) to re-enable DNSSEC.
7. Verify:
   ```
   dig +short simplemap.safecast.org @1.1.1.1     # → Cloudflare proxy IP
   dig NS safecast.org @1.1.1.1                   # → Cloudflare NS
   dig +dnssec simplemap.safecast.org @1.1.1.1    # expect ad flag
   dig DS safecast.org @8.8.8.8                   # expect Cloudflare's DS
   curl -I https://simplemap.safecast.org         # expect: server: cloudflare, cf-ray header
   ```
   Plus a full chain validation at https://dnsviz.net/d/safecast.org/dnssec/.
8. Smoke-test the production site end-to-end: login, upload, map tiles, CSV download, MCP, admin pages.

If `safecast.org` is **not** currently DNSSEC-signed: skip steps 1–3, just change nameservers (step 4), then add Cloudflare's DS record at the registrar (step 6).

## Phase 4 — Decommission AWS

After ~48 hours of stable operation **and** clean DNSSEC validation on dnsviz:

1. Disable CloudFront distribution (wait for "Disabled" state, then delete).
2. Delete ACM certificate in `us-east-1`.
3. Delete the Route53 hosted zone for `safecast.org`. (Do this last — it's your rollback if the DS handoff goes wrong.)
4. Cancel/clean up any AWS IAM users/keys that existed only for CloudFront invalidation.
5. Repo cleanup:
   - Delete `docs/cloudfront-setup.md`, `cloudfront-fix-*.md`, `CLOUDFRONT_MCP_TROUBLESHOOTING.md`, `cloudflare-cname-setup.md`.
   - Add `docs/cloudflare-setup.md` documenting the new flow.
   - Update `MEMORY.md` "CloudFront + Nginx Routing" section → "Cloudflare + Nginx Routing".
6. **Collapse the two nginx server blocks** (`simplemap.safecast.org` and `origin-simplemap.safecast.org`) into one with `server_name simplemap.safecast.org origin-simplemap.safecast.org;`. Today they have different location blocks and will drift; post-migration both names hit the same nginx with the same routing needs.
7. Replace any `aws cloudfront create-invalidation` calls in deploy scripts / GitHub Actions with Cloudflare cache purge:
   ```
   curl -X POST https://api.cloudflare.com/client/v4/zones/<ZONE_ID>/purge_cache \
     -H "Authorization: Bearer <API_TOKEN>" \
     -H "Content-Type: application/json" \
     --data '{"purge_everything":true}'
   ```

## Risks & constraints

### High-risk

- **DNSSEC handoff.** A mismatched DS → SERVFAIL for every validating resolver (Google 8.8.8.8, Cloudflare 1.1.1.1, Quad9 — i.e. most users). Always remove old DS first, wait for TTL, then add new DS. Never overlap. This is the single most fragile step in the migration.
- **Email (MX/SPF/DKIM/DMARC).** Second biggest footgun in any DNS migration. Diff carefully against `route53-export.json`.

### Free-tier constraints (acceptable but real)

- **100 s origin response timeout.** Origin nginx sets `proxy_read_timeout 120;` on every MCP REST endpoint deliberately — those queries genuinely run >30s sometimes. Cloudflare Free terminates origin responses at 100 s. Mitigation, in order of preference:
  1. Optimize/stream slow MCP endpoints so first byte ships <100 s (best fix; helps direct users too).
  2. Move slow endpoints to a grey-cloud subdomain (DNS-only) to bypass the proxy — loses TLS/DDoS protection.
  3. Accept the regression and tune queries.

  Pro ($25/mo) does not lift this; only Enterprise does.

- **100 MB upload cap.** Hard ceiling on Free, Pro, and Business; only Enterprise raises it. Origin currently accepts up to 100 MB; after dropping nginx to 95 MB (Phase 0.7), uploads ≤95 MB pass cleanly. For larger imports route to `origin-simplemap.safecast.org` (grey-cloud, direct).

- **SSE/WebSocket idle timeout.** Same 100 s edge-to-origin timeout applies to idle SSE streams. If MCP HTTP transport doesn't already send a heartbeat (~30 s), idle MCP sessions will get cut. Add a comment-line heartbeat to the MCP HTTP handler if not present.

- **WAF / rate limiting.** Free includes Bot Fight Mode and 5 custom WAF rules. If managed OWASP rules or fine-grained rate limiting on `/api/admin/upload` become necessary, that's Pro ($25/mo).

### Lower-risk

- **Cache poisoning of admin pages.** The bypass rules in Phase 1.8 are the safety net — get those right before going live.
- **SSH access.** `origin-simplemap` stays grey-cloud (DNS-only); deploy scripts already use the bare IP `65.108.24.131`, so SSH/rsync are unaffected.
- **GitHub Actions deploy.** Uses IP, not the proxied hostname — unaffected.

## Routing reference (current origin nginx, source of truth)

```
# MCP (port 3333)
/mcp                       → 3333
/mcp-http                  → 3333
/api/radiation             → 3333
/api/area                  → 3333
/api/sensors               → 3333
/api/sensor/*              → 3333
/api/device/*              → 3333
/api/spectra               → 3333
/api/stats                 → 3333
/api/extreme               → 3333
/api/info/*                → 3333
/api/gpt/*                 → 3333

# Unified server (port 8765)
/api/auth/*                → 8765
/api/user/*                → 8765
/api/admin/*               → 8765 (uploads, client_max_body_size 100M)
/api/track/*/insights      → 8765 (regex match before /api/track/)
/api/track/*               → 8765 (CSV/XLSX/JSON by extension)
/mcp-api/*                 → 8765 (MCP API docs)
/docs/*                    → 8765 (combined API docs)
/                          → 8765 (everything else)

# Web-chat (port 3334)
/assistant/*               → 3334
```

`/api/track/*` MUST stay on 8765 — port 3333's `rest_tracks.go` does no extension detection and always returns JSON, breaking CSV/XLSX downloads.
