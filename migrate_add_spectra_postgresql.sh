#!/bin/bash
#
# Migration Script: Add Spectral Data Support to PostgreSQL Database
# This script safely adds the spectra table and has_spectrum column
# to an existing safecast-new-map PostgreSQL database without data loss.
#
# Usage:
#   ./migrate_add_spectra_postgresql.sh
#
# Or with custom connection:
#   DB_HOST=localhost DB_USER=postgres DB_NAME=safecast ./migrate_add_spectra_postgresql.sh
#

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Database connection parameters (customize these or pass as environment variables)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-safecast}"
DB_PASSWORD="${DB_PASSWORD:-}"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Safecast Spectral Data Migration${NC}"
echo -e "${GREEN}PostgreSQL Database${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"
echo ""

# Check if psql is installed
if ! command -v psql &> /dev/null; then
    echo -e "${RED}ERROR: psql command not found. Please install PostgreSQL client.${NC}"
    exit 1
fi

# Function to run SQL
run_sql() {
    local sql="$1"
    local description="$2"

    echo -e "${YELLOW}${description}...${NC}"

    if [ -n "$DB_PASSWORD" ]; then
        PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$sql"
    else
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$sql"
    fi

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Done${NC}"
    else
        echo -e "${RED}✗ Failed${NC}"
        return 1
    fi
}

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

# Step 1: Check current state
echo -e "${YELLOW}Step 1: Checking current database state${NC}"
run_sql "SELECT COUNT(*) as marker_count FROM markers;" "Counting markers"
echo ""

# Step 2: Create backup recommendation
echo -e "${YELLOW}Step 2: Backup recommendation${NC}"
echo -e "${YELLOW}It's recommended to create a backup before proceeding.${NC}"
echo ""
echo "To create a backup, run (in another terminal):"
echo "  pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME > backup_$(date +%Y%m%d_%H%M%S).sql"
echo ""
read -p "Have you created a backup or want to skip? (yes/skip): " backup_confirm

if [ "$backup_confirm" != "yes" ] && [ "$backup_confirm" != "skip" ]; then
    echo "Please create a backup first."
    exit 0
fi
echo ""

# Step 3: Add has_spectrum column to markers table
echo -e "${YELLOW}Step 3: Adding has_spectrum column to markers table${NC}"

SQL_ADD_COLUMN="
DO \$\$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'markers' AND column_name = 'has_spectrum'
    ) THEN
        ALTER TABLE markers ADD COLUMN has_spectrum BOOLEAN DEFAULT FALSE;
        RAISE NOTICE 'Column has_spectrum added successfully';
    ELSE
        RAISE NOTICE 'Column has_spectrum already exists';
    END IF;
END \$\$;
"

if [ -n "$DB_PASSWORD" ]; then
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
$SQL_ADD_COLUMN
EOF
else
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
$SQL_ADD_COLUMN
EOF
fi

echo -e "${GREEN}✓ has_spectrum column added${NC}"
echo ""

# Step 4: Create spectra table
echo -e "${YELLOW}Step 4: Creating spectra table${NC}"

SQL_CREATE_TABLE="
CREATE TABLE IF NOT EXISTS spectra (
  id              BIGSERIAL PRIMARY KEY,
  marker_id       BIGINT NOT NULL REFERENCES markers(id) ON DELETE CASCADE,
  channels        TEXT,
  channel_count   INTEGER DEFAULT 1024,
  energy_min_kev  DOUBLE PRECISION,
  energy_max_kev  DOUBLE PRECISION,
  live_time_sec   DOUBLE PRECISION,
  real_time_sec   DOUBLE PRECISION,
  device_model    TEXT,
  calibration     TEXT,
  source_format   TEXT,
  filename        TEXT,
  raw_data        BYTEA,
  created_at      TIMESTAMPTZ DEFAULT NOW()
);
"

if [ -n "$DB_PASSWORD" ]; then
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_TABLE"
else
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_TABLE"
fi

echo -e "${GREEN}✓ spectra table created${NC}"
echo ""

# Step 5: Create index
echo -e "${YELLOW}Step 5: Creating index on spectra table${NC}"
run_sql "CREATE INDEX IF NOT EXISTS idx_spectra_marker_id ON spectra(marker_id);" "Creating index"
echo ""

# Step 6: Verify migration
echo -e "${YELLOW}Step 6: Verifying migration${NC}"

VERIFY_SQL="
SELECT
    'markers' as table_name,
    COUNT(*) as total_rows,
    SUM(CASE WHEN has_spectrum THEN 1 ELSE 0 END) as with_spectra
FROM markers
UNION ALL
SELECT
    'spectra' as table_name,
    COUNT(*) as total_rows,
    COUNT(*) as total_spectra
FROM spectra;
"

if [ -n "$DB_PASSWORD" ]; then
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$VERIFY_SQL"
else
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$VERIFY_SQL"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Migration completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Next steps:"
echo "  1. Restart your safecast-new-map application"
echo "  2. Upload spectrum files (.spe, .n42, .rctrk) via the web interface"
echo "  3. Markers with spectra will show a spectrum icon on the map"
echo ""
echo "To test spectral data import:"
echo "  go run scripts/insert_test_spectrum.go"
echo ""
