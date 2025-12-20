#!/bin/bash
#
# Migration Script: Add Spectral Data Support to SQLite Database
# This script safely adds the spectra table and has_spectrum column
# to an existing safecast-new-map SQLite database without data loss.
#
# Usage:
#   ./migrate_add_spectra_sqlite.sh /path/to/database.db
#

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Database file path
DB_FILE="${1:-safecast.db}"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Safecast Spectral Data Migration${NC}"
echo -e "${GREEN}SQLite Database${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Database file: $DB_FILE"
echo ""

# Check if database file exists
if [ ! -f "$DB_FILE" ]; then
    echo -e "${RED}ERROR: Database file not found: $DB_FILE${NC}"
    echo "Usage: $0 /path/to/database.db"
    exit 1
fi

# Check if sqlite3 is installed
if ! command -v sqlite3 &> /dev/null; then
    echo -e "${RED}ERROR: sqlite3 command not found. Please install SQLite.${NC}"
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
MARKER_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM markers;")
echo "Total markers: $MARKER_COUNT"
echo ""

# Step 3: Add has_spectrum column to markers table
echo -e "${YELLOW}Step 3: Adding has_spectrum column to markers table${NC}"

# Check if column exists
COLUMN_EXISTS=$(sqlite3 "$DB_FILE" "PRAGMA table_info(markers);" | grep -c "has_spectrum" || true)

if [ "$COLUMN_EXISTS" -eq 0 ]; then
    sqlite3 "$DB_FILE" "ALTER TABLE markers ADD COLUMN has_spectrum INTEGER DEFAULT 0;"
    echo -e "${GREEN}✓ has_spectrum column added${NC}"
else
    echo -e "${YELLOW}⚠ has_spectrum column already exists${NC}"
fi
echo ""

# Step 4: Create spectra table
echo -e "${YELLOW}Step 4: Creating spectra table${NC}"

sqlite3 "$DB_FILE" <<'EOF'
CREATE TABLE IF NOT EXISTS spectra (
  id              INTEGER PRIMARY KEY,
  marker_id       BIGINT NOT NULL,
  channels        TEXT,
  channel_count   INTEGER DEFAULT 1024,
  energy_min_kev  REAL,
  energy_max_kev  REAL,
  live_time_sec   REAL,
  real_time_sec   REAL,
  device_model    TEXT,
  calibration     TEXT,
  source_format   TEXT,
  filename        TEXT,
  raw_data        BLOB,
  created_at      BIGINT NOT NULL,
  FOREIGN KEY (marker_id) REFERENCES markers(id) ON DELETE CASCADE
);
EOF

echo -e "${GREEN}✓ spectra table created${NC}"
echo ""

# Step 5: Create index
echo -e "${YELLOW}Step 5: Creating index on spectra table${NC}"
sqlite3 "$DB_FILE" "CREATE INDEX IF NOT EXISTS idx_spectra_marker_id ON spectra(marker_id);"
echo -e "${GREEN}✓ Index created${NC}"
echo ""

# Step 6: Verify migration
echo -e "${YELLOW}Step 6: Verifying migration${NC}"

echo ""
echo "Table: markers"
sqlite3 "$DB_FILE" "SELECT COUNT(*) as total_markers FROM markers;"
sqlite3 "$DB_FILE" "SELECT COUNT(*) as markers_with_spectra FROM markers WHERE has_spectrum = 1;"

echo ""
echo "Table: spectra"
sqlite3 "$DB_FILE" "SELECT COUNT(*) as total_spectra FROM spectra;"

echo ""
echo "Schema verification:"
sqlite3 "$DB_FILE" ".schema spectra" | head -5

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
