#!/bin/bash
# Local Development Server Configuration
# Based on production settings from simplemap.safecast.org
# Adapted for local development on localhost:8765

./safecast-new-map \
  -port 8765 \
  -safecast-fetcher \
  -db-type pgx \
  -db-conn "postgres://postgres:@127.0.0.1:5432/safecast?sslmode=prefer" \
  -safecast-realtime \
  -support-email rob@safecast.org \
  -session-secret "some-random-secret-string" \
  -admin-password "admin123" \
  -smtp-host "smtp.gmail.com" \
  -smtp-port 587 \
  -smtp-username "oudendijk.biz@gmail.com" \
  -smtp-password "uyev szqd wsit dfnh" \
  -smtp-from "oudendijk.biz@gmail.com" \
  -base-url "http://localhost:8765" \
  -default-lat=37.43336 \
  -default-lon=141.01244 \
  -default-zoom=7 \
  -default-layer="Google Satellite" \
  -allow-registration \
  -require-auth
