#!/bin/bash
#
# Migration Script: Add Spectral Data Support to DuckDB Database
# This script safely adds the spectra table and has_spectrum column
# to an existing safecast-new-map DuckDB database without data loss.
#
# Usage:
#   ./migrate_add_spectra_duckdb.sh /path/to/database.duckdb
#

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Database file path
DB_FILE="${1:-safecast.duckdb}"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Safecast Spectral Data Migration${NC}"
echo -e "${GREEN}DuckDB Database${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Database file: $DB_FILE"
echo ""

# Check if database file exists
if [ ! -f "$DB_FILE" ]; then
    echo -e "${RED}ERROR: Database file not found: $DB_FILE${NC}"
    echo "Usage: $0 /path/to/database.duckdb"
    exit 1
fi

# Check if duckdb is installed
if ! command -v duckdb &> /dev/null; then
    echo -e "${RED}ERROR: duckdb command not found. Please install DuckDB.${NC}"
    echo "Installation: https://duckdb.org/docs/installation/"
    exit 1
fi

# Confirmation prompt
echo -e "${YELLOW}This will add spectral data support to your database.${NC}"
echo -e "${YELLOW}The migration is safe and will not delete any existing data.${NC}"
echo ""
read -p "Continue? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "Migration cancelled."
    exit 0
fi

echo ""
echo -e "${GREEN}Starting migration...${NC}"
echo ""

# Step 1: Create backup
echo -e "${YELLOW}Step 1: Creating backup${NC}"
BACKUP_FILE="${DB_FILE}.backup_$(date +%Y%m%d_%H%M%S)"
cp "$DB_FILE" "$BACKUP_FILE"
echo -e "${GREEN}✓ Backup created: $BACKUP_FILE${NC}"
echo ""

# Step 2: Check current state
echo -e "${YELLOW}Step 2: Checking current database state${NC}"
duckdb "$DB_FILE" "SELECT COUNT(*) as marker_count FROM markers;"
echo ""

# Step 3: Add has_spectrum column to markers table
echo -e "${YELLOW}Step 3: Adding has_spectrum column to markers table${NC}"

duckdb "$DB_FILE" <<'EOF'
ALTER TABLE markers ADD COLUMN IF NOT EXISTS has_spectrum BOOLEAN DEFAULT FALSE;
EOF

echo -e "${GREEN}✓ has_spectrum column added${NC}"
echo ""

# Step 4: Create spectra table and sequence
echo -e "${YELLOW}Step 4: Creating spectra table${NC}"

duckdb "$DB_FILE" <<'EOF'
CREATE SEQUENCE IF NOT EXISTS spectra_id_seq START 1;

CREATE TABLE IF NOT EXISTS spectra (
  id              BIGINT PRIMARY KEY DEFAULT nextval('spectra_id_seq'),
  marker_id       BIGINT NOT NULL,
  channels        TEXT,
  channel_count   INTEGER DEFAULT 1024,
  energy_min_kev  DOUBLE,
  energy_max_kev  DOUBLE,
  live_time_sec   DOUBLE,
  real_time_sec   DOUBLE,
  device_model    TEXT,
  calibration     TEXT,
  source_format   TEXT,
  filename        TEXT,
  raw_data        BLOB,
  created_at      TIMESTAMP DEFAULT NOW()
);
EOF

echo -e "${GREEN}✓ spectra table created${NC}"
echo ""

# Step 5: Create index
echo -e "${YELLOW}Step 5: Creating index on spectra table${NC}"
duckdb "$DB_FILE" "CREATE INDEX IF NOT EXISTS idx_spectra_marker_id ON spectra(marker_id);"
echo -e "${GREEN}✓ Index created${NC}"
echo ""

# Step 6: Verify migration
echo -e "${YELLOW}Step 6: Verifying migration${NC}"

echo ""
echo "Table: markers"
duckdb "$DB_FILE" "SELECT COUNT(*) as total_markers FROM markers;"
duckdb "$DB_FILE" "SELECT COUNT(*) as markers_with_spectra FROM markers WHERE has_spectrum = TRUE;"

echo ""
echo "Table: spectra"
duckdb "$DB_FILE" "SELECT COUNT(*) as total_spectra FROM spectra;"

echo ""
echo "Schema verification:"
duckdb "$DB_FILE" "DESCRIBE spectra;" | head -10

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Migration completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Backup location: $BACKUP_FILE"
echo ""
echo "Next steps:"
echo "  1. Restart your safecast-new-map application"
echo "  2. Upload spectrum files (.spe, .n42, .rctrk) via the web interface"
echo "  3. Markers with spectra will show a spectrum icon on the map"
echo ""
