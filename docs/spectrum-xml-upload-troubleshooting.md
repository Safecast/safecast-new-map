# Spectrum XML Upload Troubleshooting (RadiaCode `ResultDataFile`)

## Symptom

A logged-in user uploads a RadiaCode spectrum XML file (`<ResultDataFile>` /
`<EnergySpectrum>`, e.g. exported from the RadiaCode app) at
`https://simplemap.safecast.org`, and the upload fails. Depending on which
layer is broken, the failure looks different — see the three root causes
below, all of which hit the same file in practice and had to be fixed
together.

## Root cause #1 — WAF false-positive XSS block

**Symptom:** Upload fails immediately; browser Network tab shows the `/upload`
POST returning **403**. The frontend shows a misleading "🔒 Please log in or
register to upload files" message for *any* 401 or 403 response
(`public_html/map.html`, the `xhr.status === 401 || xhr.status === 403`
branch), so this looks like an auth bug even though the user is logged in.

**Diagnosis:** Query the WAF sampled requests for the ACL attached to the
CloudFront distribution:

```bash
aws wafv2 get-sampled-requests \
  --web-acl-arn "arn:aws:wafv2:us-east-1:985752656544:global/webacl/CreatedByCloudFront-92534502/49844f95-3749-41c1-ada2-3ace439c3fe4" \
  --rule-metric-name AWS-AWSManagedRulesCommonRuleSet \
  --scope CLOUDFRONT \
  --time-window StartTime=<ISO8601>,EndTime=<ISO8601> \
  --region us-east-1 \
  --query "SampledRequests[?Request.URI=='/upload']"
```

This showed `Action: BLOCK`, `RuleNameWithinRuleGroup:
CrossSiteScripting_BODY` for the exact request (matching Content-Length),
with a valid `safecast_session` cookie present — i.e. WAF blocked the request
*before* it ever reached the app, so the app's auth code never ran.

**Why XML trips this rule:** `AWSManagedRulesCommonRuleSet`'s
`CrossSiteScripting_BODY` heuristic flags tag-dense bodies
(`<ResultDataFile>`, `<EnergySpectrum>`, `<DataPoint>` repeated hundreds of
times) as a script-injection pattern. Plain-text bGeigie `.log` uploads don't
trip it; structured XML does.

**Fix:** Override `CrossSiteScripting_BODY` from Block to Count on
`AWS-AWSManagedRulesCommonRuleSet`, the same treatment already applied to
`SizeRestrictions_BODY` for oversized multipart bodies. Full details,
including why a path-scoped override (`/upload` only) isn't possible on this
WAF pricing plan, are in
[cloudfront-fix-waf-403.md](cloudfront-fix-waf-403.md).

## Root cause #2 — dose rate silently zeroed, marker filtered out

**Symptom:** Once WAF let the request through, the upload still failed with
"all markers filtered out" (or just silently produced no marker).

**Root cause:** Many RadiaCode XML exports include `<MeasurementTime>` but
omit `<LiveTime>` entirely. `pkg/spectrum/rcxml.go`'s dose/count-rate
calculation only ran `if liveTime > 0`, so the marker got `DoseRate = 0`.
`filterZeroMarkers` (`cmd/unified-server/markers.go`) then drops any marker
with `DoseRate == 0`, so the entire upload came back empty.

**Fix:** Fall back to `MeasurementTime` (real time) for the dose/count-rate
calculation when `LiveTime` is absent or zero. See
[`pkg/spectrum/rcxml.go`](../pkg/spectrum/rcxml.go) (`effectiveTime` logic
introduced in commit `25e1f41` / PR #150).

## Root cause #3 — coordinate-set endpoint was admin-only

**Symptom:** After root causes #1 and #2 were fixed, the upload succeeded and
showed the "Set Spectrum Location" dialog (since RadiaCode spectrum files
have no GPS data). Submitting a location triggered a **native browser Basic
Auth popup** ("simplemap.safecast.org is asking you to sign in") — not the
app's session login.

**Root cause:** `/api/update-coordinates`
(`pkg/httpapi/handlers_markers.go`) was gated entirely behind
`auth.IsStaticAdminAuthorized` / `auth.ChallengeStaticAdminBasic` — the
site's single static admin password over HTTP Basic Auth, unrelated to the
per-user session cookie. The upload flow's coordinate dialog
(`public_html/map.html`, `showCoordinateDialog`) is shown to any logged-in
user whose spectrum file lacks GPS data, but it called this same admin-only
endpoint, so ordinary users could never complete the flow.

**Fix:** `updateCoordinates` now authorizes either:
1. the static admin (unchanged — still gets the 401 + `WWW-Authenticate`
   Basic Auth challenge if unauthenticated), or
2. the logged-in session user who owns the track's upload record, checked
   via the new `Database.TrackBelongsToUser(trackID, internalUserID)`
   (`pkg/database/uploads.go`), which looks up `uploads.internal_user_id`.

A logged-in user who *doesn't* own the track gets a plain `403 Forbidden`
JSON response with no `WWW-Authenticate` header (no browser popup). See PR
#151 (`pkg/httpapi/server_web.go`, `pkg/httpapi/handlers_markers.go`,
`pkg/database/uploads.go`, `cmd/unified-server/main.go`).

## Summary table

| # | Layer | Symptom | Fix |
|---|-------|---------|-----|
| 1 | CloudFront WAF | 403, mislabeled "please log in" | `CrossSiteScripting_BODY` → Count on `AWS-AWSManagedRulesCommonRuleSet` |
| 2 | App: spectrum parsing | Upload succeeds at HTTP layer but marker is silently dropped | Fall back to `MeasurementTime` when `LiveTime` is absent ([rcxml.go](../pkg/spectrum/rcxml.go)) |
| 3 | App: authorization | Native Basic Auth popup when setting a GPS-less spectrum's location | Allow session-owner of the track, not just the static admin ([handlers_markers.go](../pkg/httpapi/handlers_markers.go)) |

## Related docs

- [cloudfront-fix-waf-403.md](cloudfront-fix-waf-403.md) — WAF override details and the pricing-plan limitation on scoped rules.
- [cloudfront-fix-upload-403.md](cloudfront-fix-upload-403.md) — a different, earlier 403 cause (CloudFront not forwarding session cookies).
- [api-endpoint-ownership.md](api-endpoint-ownership.md) — endpoint ownership/auth conventions, if present.
