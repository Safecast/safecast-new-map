#!/bin/bash
#
# reset_postgres_sequences.sh
# Resets PostgreSQL BIGSERIAL sequences to match the actual max IDs in tables
#
# This script fixes the common issue where sequences get out of sync after
# importing data with explicit IDs, which causes "duplicate key value violates
# unique constraint" errors.
#
# Usage:
#   export POSTGRES_URL='host=localhost port=5432 dbname=safecast user=safecast password=yourpassword sslmode=disable'
#   ./reset_postgres_sequences.sh
#
# Or with inline credentials:
#   PGPASSWORD=yourpassword ./reset_postgres_sequences.sh

set -e

# Parse POSTGRES_URL if set, otherwise use defaults
if [ -n "$POSTGRES_URL" ]; then
    # Extract components from POSTGRES_URL
    DB_HOST=$(echo "$POSTGRES_URL" | grep -oP 'host=\K[^ ]+' || echo "localhost")
    DB_PORT=$(echo "$POSTGRES_URL" | grep -oP 'port=\K[^ ]+' || echo "5432")
    DB_NAME=$(echo "$POSTGRES_URL" | grep -oP 'dbname=\K[^ ]+' || echo "safecast")
    DB_USER=$(echo "$POSTGRES_URL" | grep -oP 'user=\K[^ ]+' || echo "safecast")
    DB_PASS=$(echo "$POSTGRES_URL" | grep -oP 'password=\K[^ ]+' || echo "")
else
    DB_HOST="${DB_HOST:-localhost}"
    DB_PORT="${DB_PORT:-5432}"
    DB_NAME="${DB_NAME:-safecast}"
    DB_USER="${DB_USER:-safecast}"
    DB_PASS="${PGPASSWORD:-}"
fi

echo "🔄 Resetting PostgreSQL sequences..."
echo "   Database: $DB_NAME"
echo "   Host: $DB_HOST:$DB_PORT"
echo "   User: $DB_USER"
echo ""

# Function to reset a sequence
reset_sequence() {
    local table=$1
    local seq_name="${table}_id_seq"

    echo -n "   Resetting ${seq_name}... "

    local new_val=$(PGPASSWORD="$DB_PASS" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c \
        "SELECT setval('${seq_name}', (SELECT COALESCE(MAX(id), 1) FROM ${table}));" 2>&1)

    if [ $? -eq 0 ]; then
        echo "✓ reset to $new_val"
    else
        echo "✗ failed: $new_val"
        return 1
    fi
}

# Reset sequences for all tables with BIGSERIAL id columns
reset_sequence "markers"
reset_sequence "spectra"
reset_sequence "uploads"

echo ""
echo "✅ All sequences reset successfully!"
echo ""
echo "You can now insert new records without primary key conflicts."
