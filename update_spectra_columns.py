#!/usr/bin/env python3
"""
Update existing spectral data in PostgreSQL with parsed channel data from SQLite
"""

import sqlite3
import psycopg2
import sys
import os
from getpass import getpass

# Configuration
SQLITE_DB = os.getenv('SQLITE_DB', 'database-8765.sqlite')
PG_HOST = os.getenv('PG_HOST', 'localhost')
PG_PORT = os.getenv('PG_PORT', '5432')
PG_USER = os.getenv('PG_USER', 'safecast')
PG_DB = os.getenv('PG_DB', 'safecast')
PG_PASSWORD = os.getenv('PG_PASSWORD', '')

def main():
    print("=" * 60)
    print("Update PostgreSQL Spectral Data with Channel Information")
    print("=" * 60)
    print()

    # Check SQLite database
    if not os.path.exists(SQLITE_DB):
        print(f"Error: SQLite database not found: {SQLITE_DB}")
        sys.exit(1)

    # Get password if not set
    pg_password = PG_PASSWORD
    if not pg_password:
        pg_password = getpass("PostgreSQL password: ")

    # Connect to SQLite
    print(f"Connecting to SQLite: {SQLITE_DB}")
    sqlite_conn = sqlite3.connect(SQLITE_DB)
    sqlite_cursor = sqlite_conn.cursor()

    # Connect to PostgreSQL
    print(f"Connecting to PostgreSQL: {PG_HOST}:{PG_PORT}/{PG_DB}")
    try:
        pg_conn = psycopg2.connect(
            host=PG_HOST,
            port=PG_PORT,
            user=PG_USER,
            password=pg_password,
            database=PG_DB
        )
        pg_cursor = pg_conn.cursor()
        print("✓ PostgreSQL connection successful")
    except Exception as e:
        print(f"Error connecting to PostgreSQL: {e}")
        sys.exit(1)

    # Check how many records need updating
    pg_cursor.execute("SELECT COUNT(*) FROM spectra WHERE channels IS NULL")
    null_count = pg_cursor.fetchone()[0]

    print(f"  - Records with NULL channels: {null_count}")
    print()

    if null_count == 0:
        print("No records need updating. All done!")
        sys.exit(0)

    # Confirm update
    confirm = input(f"Update {null_count} spectral records? (yes/no): ")
    if confirm.lower() != 'yes':
        print("Update cancelled.")
        sys.exit(0)

    print()
    print("Fetching spectral data from SQLite...")

    # Fetch all spectral data from SQLite with channel information
    sqlite_cursor.execute("""
        SELECT id, channels, channel_count, energy_min_kev, energy_max_kev,
               live_time_sec, real_time_sec
        FROM spectra
        WHERE channels IS NOT NULL
        ORDER BY id
    """)

    spectra_rows = sqlite_cursor.fetchall()
    print(f"✓ Fetched {len(spectra_rows)} records from SQLite")

    # Update PostgreSQL
    print("Updating PostgreSQL records...")

    try:
        pg_cursor.execute("BEGIN")

        updated_count = 0

        for row in spectra_rows:
            spec_id, channels_str, channel_count, energy_min, energy_max, live_time, real_time = row

            # Convert channels string to array
            # SQLite stores as string "[1,2,3,...]", PostgreSQL needs list
            import json
            try:
                if isinstance(channels_str, str):
                    channels = json.loads(channels_str)
                elif channels_str is None:
                    channels = None
                else:
                    channels = channels_str
            except:
                print(f"  Warning: Could not parse channels for spectrum {spec_id}, skipping...")
                continue

            # Update spectrum with channel data
            pg_cursor.execute("""
                UPDATE spectra
                SET channels = %s,
                    channel_count = %s,
                    energy_min_kev = %s,
                    energy_max_kev = %s,
                    live_time_sec = %s,
                    real_time_sec = %s
                WHERE id = %s
            """, (
                channels, channel_count, energy_min, energy_max,
                live_time, real_time, spec_id
            ))

            if pg_cursor.rowcount > 0:
                updated_count += 1

            if updated_count % 100 == 0:
                print(f"  Updated {updated_count} records...")

        # Commit transaction
        pg_conn.commit()
        print(f"✓ Updated {updated_count} spectral records")
        print("✓ Transaction committed")

    except Exception as e:
        pg_conn.rollback()
        print(f"Error during update: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

    # Verify update
    print()
    print("Verifying update...")
    pg_cursor.execute("SELECT COUNT(*) FROM spectra WHERE channels IS NULL")
    remaining_null = pg_cursor.fetchone()[0]

    pg_cursor.execute("SELECT COUNT(*) FROM spectra WHERE channels IS NOT NULL")
    populated_count = pg_cursor.fetchone()[0]

    print(f"  - Records with channels: {populated_count}")
    print(f"  - Records still NULL: {remaining_null}")

    print()
    print("=" * 60)
    print("Update completed successfully!")
    print("=" * 60)
    print()
    print("Next step: Restart your safecast-new-map application")

    sqlite_conn.close()
    pg_conn.close()

if __name__ == "__main__":
    main()
