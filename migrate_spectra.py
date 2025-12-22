#!/usr/bin/env python3
"""
Migrate spectral data from SQLite to PostgreSQL
This script properly handles BLOB data and JSON fields
"""

import sqlite3
import psycopg2
import sys
import os
from datetime import datetime

# Configuration
SQLITE_DB = os.getenv('SQLITE_DB', 'database-8765.sqlite')
PG_HOST = os.getenv('PG_HOST', 'localhost')
PG_PORT = os.getenv('PG_PORT', '5432')
PG_USER = os.getenv('PG_USER', 'safecast')
PG_DB = os.getenv('PG_DB', 'safecast')
PG_PASSWORD = os.getenv('PG_PASSWORD', '')

def main():
    print("=" * 50)
    print("Safecast Spectral Data Migration")
    print("SQLite → PostgreSQL")
    print("=" * 50)
    print()

    # Check SQLite database
    if not os.path.exists(SQLITE_DB):
        print(f"Error: SQLite database not found: {SQLITE_DB}")
        sys.exit(1)

    # Get password if not set
    pg_password = PG_PASSWORD
    if not pg_password:
        import getpass
        pg_password = getpass.getpass("PostgreSQL password: ")

    # Connect to SQLite
    print(f"Connecting to SQLite: {SQLITE_DB}")
    sqlite_conn = sqlite3.connect(SQLITE_DB)
    sqlite_cursor = sqlite_conn.cursor()

    # Check SQLite data
    sqlite_cursor.execute("SELECT COUNT(*) FROM spectra")
    sqlite_spectra_count = sqlite_cursor.fetchone()[0]
    
    sqlite_cursor.execute("SELECT COUNT(*) FROM markers WHERE has_spectrum = 1")
    sqlite_markers_count = sqlite_cursor.fetchone()[0]

    print(f"  - Spectral records: {sqlite_spectra_count}")
    print(f"  - Markers with spectrum: {sqlite_markers_count}")
    print()

    if sqlite_spectra_count == 0:
        print("No spectral data found in SQLite database. Nothing to migrate.")
        sys.exit(0)

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

    # Check PostgreSQL current data
    pg_cursor.execute("SELECT COUNT(*) FROM spectra")
    pg_spectra_count = pg_cursor.fetchone()[0]
    
    pg_cursor.execute("SELECT COUNT(*) FROM markers WHERE has_spectrum = true")
    pg_markers_count = pg_cursor.fetchone()[0]

    print(f"  - Current spectral records: {pg_spectra_count}")
    print(f"  - Current markers with spectrum: {pg_markers_count}")
    print()

    # Confirm migration
    print(f"Ready to migrate:")
    print(f"  - {sqlite_spectra_count} spectral records")
    print(f"  - {sqlite_markers_count} marker updates")
    print()
    confirm = input("Continue with migration? (yes/no): ")
    
    if confirm.lower() != 'yes':
        print("Migration cancelled.")
        sys.exit(0)

    print()
    print("Starting migration...")

    # Fetch all spectral data from SQLite
    print("Fetching spectral data from SQLite...")
    sqlite_cursor.execute("""
        SELECT id, marker_id, channels, channel_count, energy_min_kev, energy_max_kev,
               live_time_sec, real_time_sec, device_model, calibration,
               source_format, filename, raw_data, created_at
        FROM spectra
        ORDER BY id
    """)
    
    spectra_rows = sqlite_cursor.fetchall()
    print(f"✓ Fetched {len(spectra_rows)} spectral records")

    # Insert into PostgreSQL
    print("Inserting spectral data into PostgreSQL...")
    
    try:
        pg_cursor.execute("BEGIN")
        
        inserted_count = 0
        skipped_count = 0
        
        for row in spectra_rows:
            spec_id, marker_id, channels, channel_count, energy_min, energy_max, \
            live_time, real_time, device_model, calibration, \
            source_format, filename, raw_data, created_at = row
            
            # Check if this spectrum already exists
            pg_cursor.execute("SELECT id FROM spectra WHERE id = %s", (spec_id,))
            if pg_cursor.fetchone():
                skipped_count += 1
                continue
            
            # Convert created_at to timestamp
            created_timestamp = datetime.fromtimestamp(created_at)
            
            # Insert spectrum
            pg_cursor.execute("""
                INSERT INTO spectra (
                    id, marker_id, channels, channel_count, energy_min_kev, energy_max_kev,
                    live_time_sec, real_time_sec, device_model, calibration,
                    source_format, filename, raw_data, created_at
                ) VALUES (
                    %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s
                )
            """, (
                spec_id, marker_id, channels, channel_count, energy_min, energy_max,
                live_time, real_time, device_model, calibration,
                source_format, filename, raw_data, created_timestamp
            ))
            
            inserted_count += 1
            
            if inserted_count % 100 == 0:
                print(f"  Inserted {inserted_count} records...")
        
        print(f"✓ Inserted {inserted_count} spectral records")
        if skipped_count > 0:
            print(f"  Skipped {skipped_count} existing records")
        
        # Update marker flags
        print("Updating marker flags...")
        sqlite_cursor.execute("SELECT id FROM markers WHERE has_spectrum = 1")
        marker_ids = [row[0] for row in sqlite_cursor.fetchall()]
        
        updated_count = 0
        for marker_id in marker_ids:
            pg_cursor.execute(
                "UPDATE markers SET has_spectrum = true WHERE id = %s",
                (marker_id,)
            )
            updated_count += 1
        
        print(f"✓ Updated {updated_count} marker flags")
        
        # Commit transaction
        pg_conn.commit()
        print("✓ Transaction committed")
        
    except Exception as e:
        pg_conn.rollback()
        print(f"Error during migration: {e}")
        print("Transaction rolled back.")
        sys.exit(1)

    # Verify migration
    print()
    print("Verifying migration...")
    
    pg_cursor.execute("SELECT COUNT(*) FROM spectra")
    final_spectra_count = pg_cursor.fetchone()[0]
    
    pg_cursor.execute("SELECT COUNT(*) FROM markers WHERE has_spectrum = true")
    final_markers_count = pg_cursor.fetchone()[0]

    print(f"  - Final spectral records: {final_spectra_count}")
    print(f"  - Final markers with spectrum: {final_markers_count}")
    print()

    # Close connections
    sqlite_conn.close()
    pg_conn.close()

    if final_spectra_count >= sqlite_spectra_count:
        print("=" * 50)
        print("Migration completed successfully!")
        print("=" * 50)
        print()
        print("Next steps:")
        print("1. Restart your safecast-new-map application")
        print("2. Click on a marker with spectral data")
        print("3. The spectrum graph should now display correctly")
    else:
        print("Warning: Record count mismatch")
        print(f"Expected at least {sqlite_spectra_count} records, found {final_spectra_count}")

if __name__ == '__main__':
    main()
