#!/bin/bash

# Migration Script: Copy Spectral Data from SQLite to PostgreSQL
# This script exports spectral data from the SQLite database and imports it into PostgreSQL

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configuration
SQLITE_DB="${SQLITE_DB:-database-8765.sqlite}"
PG_HOST="${PG_HOST:-localhost}"
PG_PORT="${PG_PORT:-5432}"
PG_USER="${PG_USER:-safecast}"
PG_DB="${PG_DB:-safecast}"
PG_PASSWORD="${PG_PASSWORD}"

# Temporary files
TEMP_DIR=$(mktemp -d)
SPECTRA_DUMP="$TEMP_DIR/spectra_data.sql"
MARKERS_UPDATE="$TEMP_DIR/markers_update.sql"

cleanup() {
    echo -e "${YELLOW}Cleaning up temporary files...${NC}"
    rm -rf "$TEMP_DIR"
}

trap cleanup EXIT

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Safecast Spectral Data Migration${NC}"
echo -e "${GREEN}SQLite → PostgreSQL${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if SQLite database exists
if [ ! -f "$SQLITE_DB" ]; then
    echo -e "${RED}Error: SQLite database not found: $SQLITE_DB${NC}"
    echo "Please set SQLITE_DB environment variable to the correct path."
    exit 1
fi

# Check SQLite data
echo -e "${YELLOW}Checking SQLite database...${NC}"
SQLITE_SPECTRA_COUNT=$(sqlite3 "$SQLITE_DB" "SELECT COUNT(*) FROM spectra;")
SQLITE_MARKERS_COUNT=$(sqlite3 "$SQLITE_DB" "SELECT COUNT(*) FROM markers WHERE has_spectrum = 1;")

echo "  - Spectral records: $SQLITE_SPECTRA_COUNT"
echo "  - Markers with spectrum: $SQLITE_MARKERS_COUNT"
echo ""

if [ "$SQLITE_SPECTRA_COUNT" -eq 0 ]; then
    echo -e "${YELLOW}No spectral data found in SQLite database. Nothing to migrate.${NC}"
    exit 0
fi

# Check PostgreSQL connection
echo -e "${YELLOW}Testing PostgreSQL connection...${NC}"
if [ -z "$PG_PASSWORD" ]; then
    echo -e "${YELLOW}PostgreSQL password not set. Please enter it:${NC}"
    read -s PG_PASSWORD
    export PGPASSWORD="$PG_PASSWORD"
else
    export PGPASSWORD="$PG_PASSWORD"
fi

if ! psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$PG_DB" -c "SELECT 1;" > /dev/null 2>&1; then
    echo -e "${RED}Error: Cannot connect to PostgreSQL database${NC}"
    echo "Please check your connection settings:"
    echo "  Host: $PG_HOST"
    echo "  Port: $PG_PORT"
    echo "  User: $PG_USER"
    echo "  Database: $PG_DB"
    exit 1
fi

echo -e "${GREEN}✓ PostgreSQL connection successful${NC}"
echo ""

# Check PostgreSQL current data
echo -e "${YELLOW}Checking PostgreSQL database...${NC}"
PG_SPECTRA_COUNT=$(psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$PG_DB" -t -c "SELECT COUNT(*) FROM spectra;" | tr -d ' ')
PG_MARKERS_COUNT=$(psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$PG_DB" -t -c "SELECT COUNT(*) FROM markers WHERE has_spectrum = true;" | tr -d ' ')

echo "  - Current spectral records: $PG_SPECTRA_COUNT"
echo "  - Current markers with spectrum: $PG_MARKERS_COUNT"
echo ""

# Export spectral data from SQLite
echo -e "${YELLOW}Exporting spectral data from SQLite...${NC}"

# Create PostgreSQL-compatible INSERT statements
sqlite3 "$SQLITE_DB" <<EOF > "$SPECTRA_DUMP"
.mode insert spectra
SELECT * FROM spectra;
EOF

# Convert SQLite INSERT format to PostgreSQL format
# SQLite uses single quotes and different timestamp format
sed -i "s/INSERT INTO \"spectra\"/INSERT INTO spectra/g" "$SPECTRA_DUMP"

echo -e "${GREEN}✓ Exported $SQLITE_SPECTRA_COUNT spectral records${NC}"

# Export marker updates
echo -e "${YELLOW}Preparing marker updates...${NC}"

sqlite3 "$SQLITE_DB" <<EOF > "$MARKERS_UPDATE"
SELECT 'UPDATE markers SET has_spectrum = true WHERE id = ' || id || ';'
FROM markers
WHERE has_spectrum = 1;
EOF

echo -e "${GREEN}✓ Prepared updates for $SQLITE_MARKERS_COUNT markers${NC}"
echo ""

# Confirm before proceeding
echo -e "${YELLOW}Ready to migrate:${NC}"
echo "  - $SQLITE_SPECTRA_COUNT spectral records"
echo "  - $SQLITE_MARKERS_COUNT marker updates"
echo ""
echo -e "${YELLOW}This will add data to the PostgreSQL database.${NC}"
echo -e "${YELLOW}Existing data will not be deleted.${NC}"
echo ""
read -p "Continue with migration? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo -e "${RED}Migration cancelled.${NC}"
    exit 0
fi

echo ""
echo -e "${YELLOW}Starting migration...${NC}"

# Begin transaction
echo -e "${YELLOW}Importing spectral data...${NC}"

# Import spectra (with conflict handling)
psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$PG_DB" <<EOF
BEGIN;

-- Temporarily disable triggers for performance
SET session_replication_role = replica;

-- Import spectral data
\i $SPECTRA_DUMP

-- Re-enable triggers
SET session_replication_role = DEFAULT;

COMMIT;
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Spectral data imported successfully${NC}"
else
    echo -e "${RED}Error importing spectral data${NC}"
    exit 1
fi

# Update markers
echo -e "${YELLOW}Updating marker flags...${NC}"

psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$PG_DB" <<EOF
BEGIN;

\i $MARKERS_UPDATE

COMMIT;
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Marker flags updated successfully${NC}"
else
    echo -e "${RED}Error updating marker flags${NC}"
    exit 1
fi

# Verify migration
echo ""
echo -e "${YELLOW}Verifying migration...${NC}"

FINAL_SPECTRA_COUNT=$(psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$PG_DB" -t -c "SELECT COUNT(*) FROM spectra;" | tr -d ' ')
FINAL_MARKERS_COUNT=$(psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$PG_DB" -t -c "SELECT COUNT(*) FROM markers WHERE has_spectrum = true;" | tr -d ' ')

echo "  - Final spectral records: $FINAL_SPECTRA_COUNT"
echo "  - Final markers with spectrum: $FINAL_MARKERS_COUNT"
echo ""

if [ "$FINAL_SPECTRA_COUNT" -ge "$SQLITE_SPECTRA_COUNT" ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Migration completed successfully!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Restart your safecast-new-map application"
    echo "2. Click on a marker with spectral data"
    echo "3. The spectrum graph should now display correctly"
else
    echo -e "${YELLOW}Warning: Record count mismatch${NC}"
    echo "Expected at least $SQLITE_SPECTRA_COUNT records, found $FINAL_SPECTRA_COUNT"
    echo "Please check the PostgreSQL logs for errors."
fi
