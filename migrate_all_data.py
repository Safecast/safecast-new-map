#!/usr/bin/env python3
"""
Migrate spectral data AND speed data from SQLite to PostgreSQL
This script properly handles BLOB data and efficiently updates millions of records
"""

import sqlite3
import psycopg2
import psycopg2.extras
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
BATCH_SIZE = 10000  # Process markers in batches for efficiency

def main():
    print("=" * 60)
    print("Safecast Data Migration: Spectral + Speed Data")
    print("SQLite → PostgreSQL")
    print("=" * 60)
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
    print("\nAnalyzing SQLite database...")
    
    sqlite_cursor.execute("SELECT COUNT(*) FROM spectra")
    sqlite_spectra_count = sqlite_cursor.fetchone()[0]
    
    sqlite_cursor.execute("SELECT COUNT(*) FROM markers WHERE has_spectrum = 1")
    sqlite_markers_spectrum = sqlite_cursor.fetchone()[0]
    
    sqlite_cursor.execute("SELECT COUNT(*) FROM markers WHERE speed IS NOT NULL AND speed > 0")
    sqlite_markers_speed = sqlite_cursor.fetchone()[0]
    
    sqlite_cursor.execute("SELECT COUNT(*) FROM markers")
    sqlite_total_markers = sqlite_cursor.fetchone()[0]

    print(f"  ✓ Total markers: {sqlite_total_markers:,}")
    print(f"  ✓ Spectral records: {sqlite_spectra_count:,}")
    print(f"  ✓ Markers with spectrum: {sqlite_markers_spectrum:,}")
    print(f"  ✓ Markers with speed data: {sqlite_markers_speed:,}")
    print()

    if sqlite_spectra_count == 0 and sqlite_markers_speed == 0:
        print("No spectral or speed data found in SQLite database. Nothing to migrate.")
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
    print("\nAnalyzing PostgreSQL database...")
    
    pg_cursor.execute("SELECT COUNT(*) FROM spectra")
    pg_spectra_count = pg_cursor.fetchone()[0]
    
    pg_cursor.execute("SELECT COUNT(*) FROM markers WHERE has_spectrum = true")
    pg_markers_spectrum = pg_cursor.fetchone()[0]
    
    pg_cursor.execute("SELECT COUNT(*) FROM markers WHERE speed IS NOT NULL AND speed > 0")
    pg_markers_speed = pg_cursor.fetchone()[0]
    
    pg_cursor.execute("SELECT COUNT(*) FROM markers")
    pg_total_markers = pg_cursor.fetchone()[0]

    print(f"  - Total markers: {pg_total_markers:,}")
    print(f"  - Spectral records: {pg_spectra_count:,}")
    print(f"  - Markers with spectrum: {pg_markers_spectrum:,}")
    print(f"  - Markers with speed data: {pg_markers_speed:,}")
    print()

    # Calculate what needs to be migrated
    spectra_to_migrate = sqlite_spectra_count - pg_spectra_count
    speed_to_migrate = sqlite_markers_speed - pg_markers_speed
    
    print("=" * 60)
    print("Migration Plan:")
    print("=" * 60)
    print(f"  1. Spectral records to add: {max(0, spectra_to_migrate):,}")
    print(f"  2. Marker spectrum flags to update: {max(0, sqlite_markers_spectrum - pg_markers_spectrum):,}")
    print(f"  3. Speed values to update: ~{max(0, speed_to_migrate):,}")
    print()
    
    if spectra_to_migrate <= 0 and speed_to_migrate <= 0:
        print("All data appears to be already migrated!")
        confirm = input("Do you want to re-check and update anyway? (yes/no): ")
        if confirm.lower() != 'yes':
            print("Migration cancelled.")
            sys.exit(0)
    
    print(f"Note: Speed data migration will process in batches of {BATCH_SIZE:,} records")
    print("      This may take several minutes for millions of records.")
    print()
    
    confirm = input("Continue with migration? (yes/no): ")
    
    if confirm.lower() != 'yes':
        print("Migration cancelled.")
        sys.exit(0)

    print()
    print("=" * 60)
    print("Starting Migration")
    print("=" * 60)

    # ========================================
    # PART 1: Migrate Spectral Data
    # ========================================
    if spectra_to_migrate > 0:
        print("\n[1/3] Migrating spectral data...")
        
        # Fetch all spectral data from SQLite
        sqlite_cursor.execute("""
            SELECT id, marker_id, channels, channel_count, energy_min_kev, energy_max_kev,
                   live_time_sec, real_time_sec, device_model, calibration,
                   source_format, filename, raw_data, created_at
            FROM spectra
            ORDER BY id
        """)
        
        spectra_rows = sqlite_cursor.fetchall()
        
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
                    print(f"  Inserted {inserted_count:,} spectral records...")
            
            pg_conn.commit()
            print(f"  ✓ Inserted {inserted_count:,} spectral records")
            if skipped_count > 0:
                print(f"    (Skipped {skipped_count:,} existing records)")
            
        except Exception as e:
            pg_conn.rollback()
            print(f"  ✗ Error during spectral migration: {e}")
            print("  Transaction rolled back.")
            sys.exit(1)
    else:
        print("\n[1/3] Spectral data already migrated, skipping...")

    # ========================================
    # PART 2: Update Marker Spectrum Flags
    # ========================================
    print("\n[2/3] Updating marker spectrum flags...")
    
    try:
        pg_cursor.execute("BEGIN")
        
        sqlite_cursor.execute("SELECT id FROM markers WHERE has_spectrum = 1")
        marker_ids = [row[0] for row in sqlite_cursor.fetchall()]
        
        updated_count = 0
        for marker_id in marker_ids:
            pg_cursor.execute(
                "UPDATE markers SET has_spectrum = true WHERE id = %s AND has_spectrum = false",
                (marker_id,)
            )
            if pg_cursor.rowcount > 0:
                updated_count += 1
            
            if updated_count % 1000 == 0 and updated_count > 0:
                print(f"  Updated {updated_count:,} marker flags...")
        
        pg_conn.commit()
        print(f"  ✓ Updated {updated_count:,} marker spectrum flags")
        
    except Exception as e:
        pg_conn.rollback()
        print(f"  ✗ Error updating marker flags: {e}")
        print("  Transaction rolled back.")
        sys.exit(1)

    # ========================================
    # PART 3: Migrate Speed Data
    # ========================================
    if speed_to_migrate > 0:
        print(f"\n[3/3] Migrating speed data (batch size: {BATCH_SIZE:,})...")
        
        # Get total count for progress tracking
        sqlite_cursor.execute("SELECT COUNT(*) FROM markers WHERE speed IS NOT NULL AND speed > 0")
        total_speed_records = sqlite_cursor.fetchone()[0]
        
        print(f"  Processing {total_speed_records:,} markers with speed data...")
        
        try:
            offset = 0
            total_updated = 0
            
            while True:
                # Fetch batch from SQLite
                sqlite_cursor.execute("""
                    SELECT id, speed
                    FROM markers
                    WHERE speed IS NOT NULL AND speed > 0
                    ORDER BY id
                    LIMIT ? OFFSET ?
                """, (BATCH_SIZE, offset))
                
                batch = sqlite_cursor.fetchall()
                if not batch:
                    break
                
                # Update PostgreSQL in batch using execute_values for efficiency
                pg_cursor.execute("BEGIN")
                
                # Build UPDATE query using a temporary table approach for better performance
                update_query = """
                    UPDATE markers AS m
                    SET speed = v.speed
                    FROM (VALUES %s) AS v(id, speed)
                    WHERE m.id = v.id AND (m.speed IS NULL OR m.speed = 0)
                """
                
                psycopg2.extras.execute_values(
                    pg_cursor,
                    update_query,
                    batch,
                    template="(%s, %s)"
                )
                
                pg_conn.commit()
                
                total_updated += len(batch)
                offset += BATCH_SIZE
                
                progress_pct = (total_updated / total_speed_records) * 100
                print(f"  Progress: {total_updated:,}/{total_speed_records:,} ({progress_pct:.1f}%)")
                
            print(f"  ✓ Processed {total_updated:,} speed values")
            
        except Exception as e:
            pg_conn.rollback()
            print(f"  ✗ Error during speed migration: {e}")
            print("  Transaction rolled back.")
            sys.exit(1)
    else:
        print("\n[3/3] Speed data already migrated, skipping...")

    # ========================================
    # Verification
    # ========================================
    print()
    print("=" * 60)
    print("Verifying Migration")
    print("=" * 60)
    
    pg_cursor.execute("SELECT COUNT(*) FROM spectra")
    final_spectra_count = pg_cursor.fetchone()[0]
    
    pg_cursor.execute("SELECT COUNT(*) FROM markers WHERE has_spectrum = true")
    final_markers_spectrum = pg_cursor.fetchone()[0]
    
    pg_cursor.execute("SELECT COUNT(*) FROM markers WHERE speed IS NOT NULL AND speed > 0")
    final_markers_speed = pg_cursor.fetchone()[0]

    print(f"  - Spectral records: {final_spectra_count:,}")
    print(f"  - Markers with spectrum: {final_markers_spectrum:,}")
    print(f"  - Markers with speed data: {final_markers_speed:,}")
    print()

    # Close connections
    sqlite_conn.close()
    pg_conn.close()

    # Final status
    success = True
    if final_spectra_count < sqlite_spectra_count:
        print(f"  ⚠ Warning: Spectral records mismatch (expected {sqlite_spectra_count:,}, got {final_spectra_count:,})")
        success = False
    
    if final_markers_speed < sqlite_markers_speed * 0.99:  # Allow 1% tolerance
        print(f"  ⚠ Warning: Speed data mismatch (expected ~{sqlite_markers_speed:,}, got {final_markers_speed:,})")
        success = False

    if success:
        print("=" * 60)
        print("✓ Migration Completed Successfully!")
        print("=" * 60)
        print()
        print("Next steps:")
        print("  1. Restart your safecast-new-map application")
        print("  2. Click on a marker with spectral data")
        print("  3. The spectrum graph should now display correctly")
        print("  4. Speed calculations are now complete - no more waiting!")
    else:
        print("=" * 60)
        print("⚠ Migration completed with warnings")
        print("=" * 60)
        print("Please review the warnings above.")

if __name__ == '__main__':
    main()
