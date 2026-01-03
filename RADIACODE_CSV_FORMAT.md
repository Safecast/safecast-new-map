# Radiacode 103 CSV Format Support

## Overview
The Safecast New Map now supports uploading track data from **Radiacode 103** devices in CSV format.

## File Format
The Radiacode 103 exports data as a **tab-separated values (TSV)** file with the following columns:

| Column | Description | Example |
|--------|-------------|---------|
| Timestamp | High-precision timestamp (nanoseconds) | 134093590261930000 |
| Time | Human-readable datetime | 2025-12-04 21:57:06 |
| Latitude | GPS latitude in decimal degrees | 44.0014413 |
| Longitude | GPS longitude in decimal degrees | -79.4869216 |
| Accuracy | GPS accuracy in meters | 3.79 |
| DoseRate | Dose rate in µSv/h | 4.03 |
| CountRate | Count rate in CPS | 2.17 |
| Comment | Optional comment field | (empty or text) |

## Example Data
```
Timestamp	Time	Latitude	Longitude	Accuracy	DoseRate	CountRate	Comment
134093590261930000	2025-12-04 21:57:06	44.0014413	-79.4869216	3.79	4.03	2.17	 
134093590263080000	2025-12-04 21:57:06	44.0014754	-79.4872114	6.37	4.03	2.17	 
```

## How to Upload

1. Export your track data from your Radiacode 103 device in CSV format
2. Go to the Safecast map page
3. Click the green **Upload** button
4. Select your `.csv` file
5. The system will automatically detect the Radiacode format and process it

## Auto-Detection
The upload system automatically detects whether a CSV file is:
- **Radiacode 103 format**: Tab-separated with "Timestamp", "DoseRate", and "CountRate" columns
- **AtomSwift format**: Semicolon-separated with different structure

No manual selection is needed - just upload your file!

## Technical Details

### Parser Implementation
- **Function**: `parseRadiacodeCSV()`
- **Location**: `safecast-new-map.go`
- **Delimiter**: Tab character (`\t`)
- **Date Format**: `2006-01-02 15:04:05` (from Time column)
- **Units**:
  - Dose rate is stored as-is in µSv/h
  - Count rate is stored as-is in CPS (counts per second)

### Data Validation
The parser performs the following validations:
- Skips rows with invalid coordinates (0,0)
- Skips rows with negative dose or count rates
- Validates date/time format
- Requires minimum 7 columns per row

### Storage
Uploaded tracks are stored in the database with:
- Latitude/Longitude for map display
- Date/Time (converted to Unix timestamp)
- Dose rate in µSv/h
- Count rate in CPS
- Automatically generated Track ID

## Supported Radiacode Formats
- **Radiacode 103 CSV** ✅ (This format)
- **Radiacode 101 RCTRK** ✅ (Already supported via `.rctrk` files)

## Troubleshooting

### Upload fails with "no valid data rows found"
- Check that your file has the correct header row
- Verify coordinates are not all (0,0)
- Ensure dose/count rates are non-negative
- Confirm the file uses tab separators (not spaces or commas)

### Track doesn't appear on map
- Verify GPS coordinates are valid (not 0,0)
- Check the date format matches: `YYYY-MM-DD HH:MM:SS`
- Ensure at least one valid data row exists after the header

## Credits
Radiacode format support added to the Safecast Isotope Map project to enable easy upload and visualization of radiation measurement tracks from Radiacode 103 devices.
