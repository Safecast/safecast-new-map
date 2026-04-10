# Safecast Backfill Guide - Importing Historical Data (2015-2017)

## Problem

The Safecast fetcher was designed for **incremental polling** (checking for new imports), but this prevented importing **historical data** because:

1. **Ordering issue**: API used `desc` order (newest first) + 10-page limit = old data unreachable
2. **ID filtering**: Fetcher only imported records with ID > highest ID in database

## Solution

Added **backfill mode** that:
- Uses `asc` order (oldest first) to fetch historical data
- Ignores database state (imports ALL matching records from start-date)
- Processes files chronologically

## Files Changed

- `pkg/safecast-fetcher/client.go` - Updated API URL and ordering
- `pkg/safecast-fetcher/fetcher.go` - Added backfill mode
- `safecast-new-map.go` - Added `-safecast-fetcher-backfill` flag

## Usage

### Option 1: Quick Test (Recommended First)

Test backfill mode with one batch (~250 files):

```bash
./test_backfill.sh
```

This will:
- Stop your running fetcher temporarily
- Import one batch of 2015 data (10 pages × ~25 files = 250 files)
- Show results in logs

**Check the output** for lines like:
```
[safecast-fetcher] poll: BACKFILL MODE - importing all records from startDate=2015-11-01
[safecast-fetcher] page 1: fetched 25 imports (IDs 20531-20552)
[safecast-fetcher] page 1: found 25 new imports (> 0)
[safecast-fetcher] importing: ID=20531 file=safecast-upload.log (1338 measurements)
```

### Option 2: Complete Backfill (~4,775 files)

After testing works, run the complete backfill:

```bash
./backfill_complete.sh
```

This will:
1. Temporarily increase the 10-page limit to 500
2. Import ALL ~4,775 files from Nov 2015 onwards
3. Restore the original 10-page safety limit
4. Takes ~10-30 minutes depending on network and processing speed

### Option 3: Manual Incremental Backfill

If you prefer multiple smaller runs:

```bash
./backfill_2015_2017.sh
```

- Imports ~250 files per hour
- Run multiple times until all data imported
- More conservative, easier to monitor

## Verification

After backfill, check your database:

```sql
-- Count total Safecast API imports
SELECT COUNT(*) FROM uploads WHERE source = 'safecast-api';

-- Check date range
SELECT 
  MIN(filename) as first_file,
  MAX(filename) as last_file,
  COUNT(*) as total_imports
FROM uploads 
WHERE source = 'safecast-api';

-- Check for 2015-2017 data specifically
SELECT COUNT(*) 
FROM uploads 
WHERE source = 'safecast-api' 
  AND uploaded_at BETWEEN '2015-11-01' AND '2017-11-01';
```

## Important Notes

### The 10-Page Safety Limit

The fetcher has a 10-page limit per poll cycle to prevent:
- API overload
- Excessive memory usage
- Long-running blocking operations

For backfill, you can either:
1. **Use backfill_complete.sh** - Automatically removes/restores the limit
2. **Run incrementally** - Multiple runs of 250 files each
3. **Manually edit** - Increase limit in `pkg/safecast-fetcher/fetcher.go` line 221

### Date Ranges

Your API query shows 191 pages for 2015-11-01 to 2016-11-01. Adjust the start date in scripts:

- `2015-11-01` - All data from Nov 2015 onwards
- `2015-11-01` to `2017-11-01` - Use API's uploaded_before filter (not currently supported in fetcher)

### Backfill Mode vs Normal Mode

**Backfill Mode** (`-safecast-fetcher-backfill`):
- Starts from ID=0
- Imports ALL matching records
- Use for historical data import

**Normal Mode** (default):
- Starts from highest ID in database
- Only imports NEW records
- Use for ongoing polling

## Troubleshooting

### "Found 0 new imports"

If you see this in backfill mode, the fetcher isn't in backfill mode. Check:
- Is `-safecast-fetcher-backfill` flag present?
- Does log say "BACKFILL MODE"?

### No imports processed

- Check API is accessible: `curl http://safecastapi-prd-010.baebmmfncu.us-west-2.elasticbeanstalk.com/`
- Verify database connection
- Check logs for error messages

### Want to import specific date range

Currently the fetcher supports:
- `uploaded_after` (start date)
- Not yet: `uploaded_before` (end date)

To import only 2015-2017:
1. Run backfill with `2015-11-01` start date
2. Stop when imports reach 2017
3. OR: Modify client.go to support uploaded_before parameter

## Resetting After Backfill

After completing the historical backfill, restart your normal fetcher:

```bash
# Stop backfill
pkill -f safecast-fetcher-backfill

# Restart normal polling
./safecast-new-map -safecast-fetcher -safecast-realtime -admin-password test123
```

The fetcher will now resume normal incremental polling from the highest ID.
