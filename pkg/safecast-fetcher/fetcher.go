package safecastfetcher

import (
	"context"
	"fmt"
	"time"

	"safecast-new-map/pkg/database"
)

const (
	// SourceTypeSafecastAPI is the source identifier for API-imported files
	SourceTypeSafecastAPI = "safecast-api"
)

// Fetcher handles periodic polling and importing from Safecast API
type Fetcher struct {
	client             *Client
	db                 *database.Database
	dbType             string
	batchSize          int
	startDate          string
	effectiveStartDate string // Tracks where to resume backfill (updated as we import)
	importer           ImporterFunc
	logf               func(string, ...any)
	backfillMode       bool // When true, imports all matching records regardless of lastID
	newestFirst        bool // When true, fetch newest imports first
}

// Config contains configuration for the fetcher
type Config struct {
	DB           *database.Database
	DBType       string
	Interval     time.Duration
	BatchSize    int
	StartDate    string
	Importer     ImporterFunc
	Logf         func(string, ...any)
	BackfillMode bool // When true, imports all matching records regardless of database state
	NewestFirst  bool // When true, fetch newest imports first instead of oldest
}

// Start launches the background polling service
func Start(ctx context.Context, cfg Config) {
	if cfg.Logf == nil {
		cfg.Logf = func(string, ...any) {}
	}

	cfg.Logf("[safecast-fetcher] start: interval=%s batch=%d start_date=%s backfill=%v newest_first=%v",
		cfg.Interval, cfg.BatchSize, cfg.StartDate, cfg.BackfillMode, cfg.NewestFirst)

	fetcher := &Fetcher{
		client:       NewClient(),
		db:           cfg.DB,
		dbType:       cfg.DBType,
		batchSize:    cfg.BatchSize,
		startDate:    cfg.StartDate,
		importer:     cfg.Importer,
		logf:         cfg.Logf,
		backfillMode: cfg.BackfillMode,
		newestFirst:  cfg.NewestFirst,
	}

	// Launch background polling goroutine
	go fetcher.run(ctx, cfg.Interval)
}

// run is the main polling loop
func (f *Fetcher) run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Poll immediately on startup
	if err := f.poll(ctx); err != nil {
		f.logf("[safecast-fetcher] initial poll error: %v", err)
	}

	// Then poll on each tick
	for {
		select {
		case <-ctx.Done():
			f.logf("[safecast-fetcher] stopped")
			return
		case <-ticker.C:
			if err := f.poll(ctx); err != nil {
				f.logf("[safecast-fetcher] poll error: %v", err)
			}
		}
	}
}

// poll fetches and imports new approved files
func (f *Fetcher) poll(ctx context.Context) error {
	// Get the last imported Safecast ID
	// In backfill mode, start from 0 to import all historical data
	var lastID int64
	var startPage int = 1
	if f.backfillMode {
		lastID = 0
		// If start_date is specified, check if we can resume from a later date
		if f.startDate != "" {
			// Check for existing imports to resume from where we left off
			latestDate, err := f.db.GetLatestImportDate(ctx, SourceTypeSafecastAPI)
			if err == nil && latestDate != "" && latestDate > f.startDate {
				// Resume from the day after the latest import
				f.effectiveStartDate = latestDate
				f.logf("[safecast-fetcher] poll: BACKFILL MODE - resuming from %s (latest import date)", f.effectiveStartDate)
			} else if f.effectiveStartDate == "" || f.effectiveStartDate < f.startDate {
				// First run or reset - use configured start date
				f.effectiveStartDate = f.startDate
				f.logf("[safecast-fetcher] poll: BACKFILL MODE - starting from %s", f.effectiveStartDate)
			} else {
				f.logf("[safecast-fetcher] poll: BACKFILL MODE - continuing from %s", f.effectiveStartDate)
			}
			startPage = 1
		} else {
			// Count how many records we've imported to calculate starting page
			count, err := f.db.CountImportsBySource(ctx, SourceTypeSafecastAPI)
			if err == nil && count > 0 {
				minID, _ := f.db.GetMinImportedSafecastID(ctx, SourceTypeSafecastAPI)
				maxID, _ := f.db.GetLastImportedSafecastID(ctx, SourceTypeSafecastAPI)
				// Calculate page based on how many records we've actually imported
				// 25 records per page, start a few pages earlier to be safe
				pagesImported := count / 25
				startPage = pagesImported - 5 // Start 5 pages earlier to catch gaps
				if startPage < 1 {
					startPage = 1
				}
				f.logf("[safecast-fetcher] poll: BACKFILL MODE - imported=%d, min_id=%d, max_id=%d, starting at page %d", count, minID, maxID, startPage)
			} else {
				f.logf("[safecast-fetcher] poll: BACKFILL MODE - importing all records from start")
			}
		}
	} else {
		var err error
		lastID, err = f.db.GetLastImportedSafecastID(ctx, SourceTypeSafecastAPI)
		if err != nil {
			return fmt.Errorf("get last imported ID: %w", err)
		}
		f.logf("[safecast-fetcher] poll: checking for imports after ID %d", lastID)
	}

	// Fetch approved imports from API with pagination
	allImports, err := f.fetchNewImports(ctx, lastID, startPage)
	if err != nil {
		return fmt.Errorf("fetch imports: %w", err)
	}

	f.logf("[safecast-fetcher] poll: found %d new approved imports", len(allImports))

	if len(allImports) == 0 {
		return nil
	}

	// Process each import
	imported := 0
	skipped := 0
	errors := 0

	for _, imp := range allImports {
		// Check context cancellation
		select {
		case <-ctx.Done():
			f.logf("[safecast-fetcher] poll cancelled: imported %d/%d, skipped %d, errors %d",
				imported, len(allImports), skipped, errors)
			return ctx.Err()
		default:
		}

		// Double-check if already imported
		exists, err := f.db.CheckImportExists(ctx, SourceTypeSafecastAPI, imp.ID)
		if err != nil {
			f.logf("[safecast-fetcher] import #%d: check failed: %v", imp.ID, err)
			errors++
			continue
		}
		if exists {
			f.logf("[safecast-fetcher] import #%d: already imported, skipping", imp.ID)
			skipped++
			continue
		}

		// Download log file
		// Download log file
		f.logf("[safecast-fetcher] import #%d: downloading %s", imp.ID, imp.Name)
		content, filename, err := DownloadLogFile(ctx, imp.SourceURL)
		if err != nil {
			f.logf("[safecast-fetcher] import #%d: download failed: %v", imp.ID, err)
			errors++
			continue
		}

		// Fetch username from API
		username := ""
		if imp.UserID > 0 {
			user, err := f.client.FetchUser(ctx, imp.UserID)
			if err != nil {
				f.logf("[safecast-fetcher] import #%d: warning: failed to fetch username for user %d: %v", imp.ID, imp.UserID, err)
			} else if user != nil {
				username = user.Name
			}
		}

		// Import the file
		result, err := ImportSafecastFile(ctx, content, filename, imp.ID, imp.SourceURL, fmt.Sprintf("%d", imp.UserID), username, imp.Comment, f.db, f.dbType, f.importer)
		if err != nil {
			f.logf("[safecast-fetcher] import #%d: import failed: %v", imp.ID, err)
			errors++
			continue
		}

		f.logf("[safecast-fetcher] import #%d: imported track %s with %d markers",
			imp.ID, result.TrackID, result.MarkerCount)
		imported++
	}

	f.logf("[safecast-fetcher] summary: imported %d/%d, skipped %d, errors %d",
		imported, len(allImports), skipped, errors)

	// In backfill mode, update effectiveStartDate based on latest import
	// so next poll cycle resumes from where we left off
	if f.backfillMode && imported > 0 {
		latestDate, err := f.db.GetLatestImportDate(ctx, SourceTypeSafecastAPI)
		if err == nil && latestDate != "" && latestDate > f.effectiveStartDate {
			f.effectiveStartDate = latestDate
			f.logf("[safecast-fetcher] backfill: updated resume point to %s", f.effectiveStartDate)
		}
	}

	return nil
}

// fetchNewImports retrieves all new imports from the API with pagination
func (f *Fetcher) fetchNewImports(ctx context.Context, lastID int64, startPage int) ([]SafecastImport, error) {
	var allImports []SafecastImport
	page := startPage
	consecutiveSkipped := 0   // Track how many consecutive pages have all-skipped records
	consecutiveErrors := 0    // Track consecutive page fetch errors
	var failedPages []int     // Track which pages failed for logging
	// Use effectiveStartDate in backfill mode (tracks where we left off)
	currentStartDate := f.startDate
	if f.backfillMode && f.effectiveStartDate != "" {
		currentStartDate = f.effectiveStartDate
	}
	var lastImportDate string        // Track last successful import date

	// In normal mode, always fetch newest first to find recent imports quickly
	// In backfill mode, respect the newestFirst flag
	fetchNewestFirst := f.newestFirst
	if !f.backfillMode {
		fetchNewestFirst = true
	}

	for {
		// Fetch page (pass newestFirst to control sort order)
		imports, err := f.client.FetchApprovedImports(ctx, currentStartDate, page, fetchNewestFirst)
		if err != nil {
			// Log error and skip this page in backfill mode
			f.logf("[safecast-fetcher] page %d: ERROR - %v (skipping page)", page, err)
			failedPages = append(failedPages, page)
			consecutiveErrors++

			// If we hit too many consecutive errors, something might be wrong
			if consecutiveErrors >= 10 {
				f.logf("[safecast-fetcher] ERROR: 10 consecutive page errors, stopping backfill")
				if len(failedPages) > 0 {
					f.logf("[safecast-fetcher] Failed pages: %v", failedPages)
				}
				return nil, fmt.Errorf("too many consecutive errors at page %d", page)
			}

			// Skip this page and continue
			page++
			if page > 3000 {
				break
			}
			continue
		}

		// Reset consecutive errors on success
		consecutiveErrors = 0

		// No more results
		if len(imports) == 0 {
			if f.backfillMode {
				// In backfill mode, treat empty pages like pages with no new imports
				consecutiveSkipped++
				f.logf("[safecast-fetcher] page %d: no results, continuing (%d consecutive empty)", page, consecutiveSkipped)
				// WORKAROUND: Detect stuck pagination - advance date filter after 20 empty pages
				if consecutiveSkipped >= 20 {
					if lastImportDate != "" && lastImportDate > currentStartDate {
						// Only advance if lastImportDate is NEWER (prevent backwards loop)
						f.logf("[safecast-fetcher] backfill: STUCK at page %d - advancing date filter from %s to %s",
							page, currentStartDate, lastImportDate)
						currentStartDate = lastImportDate
					} else if currentStartDate != "" && currentStartDate < time.Now().Format("2006-01-02") {
						// Advance by 7 days if no newer date available
						nextDate := advanceDate(currentStartDate, 7)
						f.logf("[safecast-fetcher] backfill: STUCK at page %d - advancing by 7 days: %s -> %s",
							page, currentStartDate, nextDate)
						currentStartDate = nextDate
					} else {
						// No date filter, can't advance - stop after 100 pages
						if consecutiveSkipped >= 100 {
							f.logf("[safecast-fetcher] backfill: 100 consecutive empty pages, assuming complete")
							break
						}
					}
					page = 1  // Reset to first page with new date filter
					consecutiveSkipped = 0
					continue
				}
				// Continue to next page
				page++
				if page > 3000 {
					break
				}
				continue
			}
			// In normal mode, stop at first empty page
			f.logf("[safecast-fetcher] page %d: no results, stopping", page)
			break
		}

		// Log what we got
		if len(imports) > 0 {
			f.logf("[safecast-fetcher] page %d: fetched %d imports (IDs %d-%d)",
				page, len(imports), imports[0].ID, imports[len(imports)-1].ID)
		}

		// Filter out already-processed imports
		newImports := 0
		for _, imp := range imports {
		// In normal mode, check all imports on first 5 pages (catches out-of-order approvals)
		// In backfill mode, only process imports with ID > lastID
		shouldProcess := false
		if f.backfillMode {
			shouldProcess = imp.ID > lastID
		} else {
			// Normal mode: check all imports on first 5 pages
			shouldProcess = true
		}

		if shouldProcess {
			// Check if already exists in DB
			exists, err := f.db.CheckImportExists(ctx, SourceTypeSafecastAPI, imp.ID)
			if err == nil && exists {
				continue // Skip already imported
			}
			allImports = append(allImports, imp)
			newImports++
			// Track the latest import date for pagination workaround
			if !imp.CreatedAt.IsZero() {
				lastImportDate = imp.CreatedAt.Format("2006-01-02") // YYYY-MM-DD
			}
		}
		}

		f.logf("[safecast-fetcher] page %d: found %d new imports", page, newImports)

		// In backfill mode, track consecutive pages with no new imports
		if f.backfillMode {
			if newImports == 0 {
				consecutiveSkipped++
				// Continue through a few skipped pages in case there are gaps
				if consecutiveSkipped >= 3 {
					f.logf("[safecast-fetcher] backfill: 3 consecutive pages with no new records, checking deeper...")
				}
				// WORKAROUND: Detect stuck pagination (Safecast API bug)
				// If we hit 20 consecutive skipped pages, advance the date filter
				if consecutiveSkipped >= 20 {
					// Stop if we have no date filter to advance (backfill without start_date)
					if f.startDate == "" && currentStartDate == "" {
						f.logf("[safecast-fetcher] backfill: reached pagination limit at page %d (no date filter), stopping", page)
						break
					}

					// Get today's date to prevent advancing into the future
					today := time.Now().Format("2006-01-02")

					if lastImportDate != "" && lastImportDate > currentStartDate && lastImportDate <= today {
						// Only advance if lastImportDate is NEWER but not in the future
						f.logf("[safecast-fetcher] backfill: STUCK at page %d - advancing date filter from %s to %s",
							page, currentStartDate, lastImportDate)
						currentStartDate = lastImportDate
					} else if currentStartDate != "" && currentStartDate < today {
						// Advance by 7 days if no newer date available, but don't go past today
						nextDate := advanceDate(currentStartDate, 7)
						if nextDate > today {
							nextDate = today
						}
						f.logf("[safecast-fetcher] backfill: STUCK at page %d - advancing by 7 days: %s -> %s",
							page, currentStartDate, nextDate)
						currentStartDate = nextDate

						// If we've reached today, stop
						if currentStartDate >= today {
							f.logf("[safecast-fetcher] backfill: reached today's date (%s), stopping", today)
							break
						}
					} else {
						// Can't advance further, stop
						f.logf("[safecast-fetcher] backfill: cannot advance date further, stopping at %s", currentStartDate)
						break
					}
					page = 1  // Reset to first page with new date filter
					consecutiveSkipped = 0
					continue
				}
			} else {
				consecutiveSkipped = 0
			}
		}

		// Check if we've hit our batch size limit
		if f.batchSize > 0 && len(allImports) >= f.batchSize {
			allImports = allImports[:f.batchSize]
			break
		}

		// Move to next page
		page++

		// In normal mode, only check first 5 pages for new uploads
		if !f.backfillMode && page > 5 {
			f.logf("[safecast-fetcher] normal mode: stopped after 5 pages")
			break
		}

		// Safety: don't fetch more than 3000 pages in one poll cycle (75,000 records)
		// In backfill mode with all records imported, stop after 100 consecutive empty pages (2500 records gap)
		if f.backfillMode && consecutiveSkipped >= 100 {
			f.logf("[safecast-fetcher] backfill: 100 consecutive pages with no new records, assuming complete")
			break
		}
		if page > 3000 {
			f.logf("[safecast-fetcher] warning: stopped at page 3000, consider increasing poll interval")
			break
		}
	}

	// Log any failed pages at the end
	if len(failedPages) > 0 {
		f.logf("[safecast-fetcher] WARNING: %d pages failed during fetch: %v", len(failedPages), failedPages)
		f.logf("[safecast-fetcher] You may want to retry these pages later or investigate the API issues")
	}

	return allImports, nil
}

// advanceDate adds days to a YYYY-MM-DD date string
func advanceDate(dateStr string, days int) string {
	if dateStr == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr // Return original on parse error
	}
	advanced := t.AddDate(0, 0, days)
	return advanced.Format("2006-01-02")
}
