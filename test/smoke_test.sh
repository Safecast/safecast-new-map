#!/usr/bin/env bash
# Smoke tests for safecast-new-map endpoints.
# Usage: ./test/smoke_test.sh [BASE_URL] [ADMIN_PASSWORD]
# Example: ./test/smoke_test.sh http://localhost:8765 test123
# Example (production): ./test/smoke_test.sh https://simplemap.safecast.org test123

set -euo pipefail

BASE="${1:-http://localhost:8765}"
ADMIN_PW="${2:-test123}"
PASS=0
FAIL=0

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
RESET='\033[0m'

# --------------------------------------------------------------------------
# Helpers
# --------------------------------------------------------------------------

check() {
  local description="$1"
  local expected_code="$2"
  local actual_code="$3"
  local body="$4"
  if [ "$actual_code" = "$expected_code" ]; then
    echo -e "${GREEN}PASS${RESET} [$actual_code] $description"
    PASS=$((PASS + 1))
  else
    echo -e "${RED}FAIL${RESET} [$actual_code expected $expected_code] $description"
    if [ -n "$body" ]; then
      echo "       Body: $(echo "$body" | head -c 200)"
    fi
    FAIL=$((FAIL + 1))
  fi
}

get() {
  local url="$BASE$1"
  local resp
  resp=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$url" "${@:2}")
  local body
  body=$(cat /tmp/smoke_body)
  echo "$resp $body"
}

post() {
  local path="$1"
  local data="$2"
  shift 2
  local url="$BASE$path"
  local resp
  resp=$(curl -s -o /tmp/smoke_body -w "%{http_code}" \
    -X POST -H "Content-Type: application/json" \
    -d "$data" "$url" "$@")
  local body
  body=$(cat /tmp/smoke_body)
  echo "$resp $body"
}

# Run a GET and extract status code
get_code()  { get "$1" | awk '{print $1}'; }
get_body()  { get "$1" | cut -d' ' -f2-; }
post_code() { post "$1" "$2" | awk '{print $1}'; }
post_body() { post "$1" "$2" | cut -d' ' -f2-; }

section() { echo -e "\n${YELLOW}=== $1 ===${RESET}"; }

# --------------------------------------------------------------------------
# Map server endpoints
# --------------------------------------------------------------------------

section "Page routes"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/")
check "GET / (map page)"        200 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/home")
check "GET /home"               200 "$CODE"

section "Static / utility"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/qrpng?u=https://safecast.org")
check "GET /qrpng returns PNG"  200 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/licenses/mit")
check "GET /licenses/mit"       200 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/licenses/cc0")
check "GET /licenses/cc0"       200 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/map-api/")
check "GET /map-api/"           200 "$CODE"

section "API overview & data"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api")
check "GET /api overview"       200 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/tracks")
check "GET /api/tracks"         200 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/countries")
check "GET /api/countries"      200 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/latest?lat=35.6&lon=139.7")
check "GET /api/latest (Tokyo)" 200 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/latest")
check "GET /api/latest (no params) -> 400" 400 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/geoip")
check "GET /api/geoip"          "200 204" "$CODE" || true  # 204 when auto-locate off
# Allow either 200 or 204
if [ "$CODE" = "200" ] || [ "$CODE" = "204" ]; then
  echo -e "${GREEN}PASS${RESET} [$CODE] GET /api/geoip (200 or 204)"
  PASS=$((PASS + 1))
else
  echo -e "${RED}FAIL${RESET} [$CODE expected 200 or 204] GET /api/geoip"
  FAIL=$((FAIL + 1))
fi

section "Shorten endpoint"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/shorten")
check "GET /api/shorten -> 405" 405 "$CODE"
BODY=$(curl -s -o /tmp/smoke_body -w "%{http_code}" \
  -X POST -H "Content-Type: application/json" \
  -d '{"url":"/map?lat=35.6&lon=139.7","commit":false}' \
  -H "Host: simplemap.safecast.org" \
  "$BASE/api/shorten")
CODE=$(echo "$BODY" | awk '{print $1}')
check "POST /api/shorten preview" 200 "$CODE"

section "Spectrum & track info"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/markers/spectra?minLat=-90&maxLat=90&minLon=-180&maxLon=180")
check "GET /api/markers/spectra" 200 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/track-info/nonexistent-track")
check "GET /api/track-info (unknown track)" 200 "$CODE"  # returns minimal JSON

section "Markers / streams"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/get_markers?minLat=34&maxLat=36&minLon=138&maxLon=141")
check "GET /get_markers (bounding box)" 200 "$CODE"

section "Upload endpoint"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/upload")
check "GET /upload -> 405"      405 "$CODE"

# --------------------------------------------------------------------------
# Auth endpoints
# --------------------------------------------------------------------------

section "Auth - register"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/auth/register")
check "GET /api/auth/register -> 405" 405 "$CODE"

# Use a unique email and randomly generated password to avoid conflicts on repeated runs
TS=$(date +%s)
REG_EMAIL="smoketest_${TS}@example.com"
REG_PW="$(openssl rand -base64 16 | tr -dc 'A-Za-z0-9' | head -c 16)Aa1!"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" \
  -X POST -H "Content-Type: application/json" \
  -d "{\"email\":\"$REG_EMAIL\",\"username\":\"smoketest_$TS\",\"password\":\"$REG_PW\"}" \
  "$BASE/api/auth/register")
check "POST /api/auth/register (new user)" 201 "$CODE"

# Duplicate registration
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" \
  -X POST -H "Content-Type: application/json" \
  -d "{\"email\":\"$REG_EMAIL\",\"username\":\"smoketest2_$TS\",\"password\":\"$REG_PW\"}" \
  "$BASE/api/auth/register")
check "POST /api/auth/register (duplicate email) -> 409" 409 "$CODE"

section "Auth - login"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/auth/login")
check "GET /api/auth/login -> 405" 405 "$CODE"

# Login and capture cookie
LOGIN_RESP=$(curl -s -c /tmp/smoke_cookies -o /tmp/smoke_body -w "%{http_code}" \
  -X POST -H "Content-Type: application/json" \
  -d "{\"email\":\"$REG_EMAIL\",\"password\":\"$REG_PW\"}" \
  "$BASE/api/auth/login")
check "POST /api/auth/login (valid credentials)" 200 "$LOGIN_RESP"

CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" \
  -X POST -H "Content-Type: application/json" \
  -d "{\"email\":\"$REG_EMAIL\",\"password\":\"wrongpassword\"}" \
  "$BASE/api/auth/login")
check "POST /api/auth/login (wrong password) -> 401" 401 "$CODE"

section "Auth - profile (requires session)"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/user/profile")
check "GET /api/user/profile (no auth) -> 401" 401 "$CODE"

CODE=$(curl -s -b /tmp/smoke_cookies -o /tmp/smoke_body -w "%{http_code}" "$BASE/api/user/profile")
check "GET /api/user/profile (with session)" 200 "$CODE"

section "Auth - logout"
CODE=$(curl -s -b /tmp/smoke_cookies -o /tmp/smoke_body -w "%{http_code}" \
  -X POST "$BASE/api/auth/logout")
check "POST /api/auth/logout" 200 "$CODE"

section "Auth - forgot password"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" \
  -X POST -H "Content-Type: application/json" \
  -d '{"email":"nobody@example.com"}' \
  "$BASE/api/auth/forgot-password")
check "POST /api/auth/forgot-password (unknown email, anti-enum)" 200 "$CODE"

# --------------------------------------------------------------------------
# Admin endpoints (password auth)
# --------------------------------------------------------------------------

section "Admin endpoints"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" \
  "$BASE/api/admin/uploads?password=$ADMIN_PW")
check "GET /api/admin/uploads (admin pw)" 200 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" \
  "$BASE/api/admin/tracks?password=$ADMIN_PW")
check "GET /api/admin/tracks (admin pw)" 200 "$CODE"
CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" \
  "$BASE/api/admin/users?password=$ADMIN_PW")
check "GET /api/admin/users (admin pw)" 200 "$CODE"

# --------------------------------------------------------------------------
# MCP / Swagger (if on same host; skip if separate port)
# --------------------------------------------------------------------------

section "Swagger docs (MCP server - may be on different port)"
MCP_BASE="${MCP_BASE:-}"
if [ -n "$MCP_BASE" ]; then
  CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$MCP_BASE/mcp-api/")
  check "GET /mcp-api/ (Swagger UI)" 200 "$CODE"
  CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$MCP_BASE/api/radiation?lat=35.6&lon=139.7")
  check "GET /api/radiation (MCP REST)" 200 "$CODE"
  CODE=$(curl -s -o /tmp/smoke_body -w "%{http_code}" "$MCP_BASE/api/stats")
  check "GET /api/stats (MCP REST)" 200 "$CODE"
else
  echo -e "${YELLOW}SKIP${RESET} MCP server endpoints (set MCP_BASE=http://host:3333 to test)"
fi

# --------------------------------------------------------------------------
# Summary
# --------------------------------------------------------------------------

echo ""
echo "========================================"
TOTAL=$((PASS + FAIL))
if [ $FAIL -eq 0 ]; then
  echo -e "${GREEN}All $TOTAL tests passed${RESET}"
else
  echo -e "${RED}$FAIL/$TOTAL tests FAILED${RESET}"
  exit 1
fi
