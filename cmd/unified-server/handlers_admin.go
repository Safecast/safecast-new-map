package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/countryresolver"
	"safecast-new-map/pkg/database"
	safecastfetcher "safecast-new-map/pkg/safecast-fetcher"
)

// =====================
// API  — Spectrum Data (handlers in pkg/httpapi)
// =====================

// =====================
// ADMIN API
// =====================

// formatUploadRow formats a single upload record as an HTML table row.
func formatUploadRow(upload database.Upload, password string) string {
	uploadTime := time.Unix(upload.CreatedAt, 0).Format("2006-01-02 15:04:05")
	fileSize := formatFileSize(upload.FileSize)

	// Format recording date
	recordingDate := "-"
	recordingDateISO := ""
	if upload.RecordingDate > 0 {
		t := time.Unix(upload.RecordingDate, 0)
		recordingDate = t.Format("2006-01-02 15:04:05")
		recordingDateISO = t.UTC().Format("2006-01-02")
	}

	// Format detector display
	detectorDisplay := "-"
	if upload.Detector != "" {
		detectorDisplay = fmt.Sprintf(`<span style="background: #607D8B; color: white; padding: 2px 8px; border-radius: 3px; font-size: 0.85em;">%s</span>`, upload.Detector)
	}

	// Format source display and get numeric source ID for sorting
	sourceDisplay := "manual"
	sourceIDNumeric := "0" // Default for manual uploads
	if upload.Source != "" {
		if upload.SourceID != "" {
			sourceIDNumeric = upload.SourceID // Store for data attribute
			if upload.SourceURL != "" {
				sourceDisplay = fmt.Sprintf(`%s (<a href="%s" target="_blank">#%s</a>)`,
					upload.Source, upload.SourceURL, upload.SourceID)
			} else {
				sourceDisplay = fmt.Sprintf("%s (#%s)", upload.Source, upload.SourceID)
			}
		} else {
			sourceDisplay = upload.Source
		}
	}

	// Format user display (username + ID)
	userDisplay := "-"
	if upload.InternalUserID != "" {
		// Internal user (authenticated upload)
		userText := upload.Username
		if userText == "" {
			userText = upload.InternalUserID
		}
		userDisplay = fmt.Sprintf(`<a href="/api/admin/uploads?password=%s&user_id=%s">%s</a>`,
			password, upload.InternalUserID, userText)
	} else if upload.UserID != "" {
		// External user (Safecast API import)
		userText := upload.UserID
		if upload.Username != "" {
			userText = fmt.Sprintf("%s (%s)", upload.Username, upload.UserID)
		}
		userDisplay = fmt.Sprintf(`<a href="/api/admin/uploads?password=%s&user_id=%s">%s</a>`,
			password, upload.UserID, userText)
	}

	// HTML-escape data attributes to prevent breaking the attribute syntax
	escapedFilename := template.HTMLEscapeString(upload.Filename)
	escapedUsername := template.HTMLEscapeString(upload.Username)
	escapedDetector := template.HTMLEscapeString(upload.Detector)
	escapedNotes := template.HTMLEscapeString(upload.Notes)
	escapedComment := template.HTMLEscapeString(upload.Comment)

	// Build optional link to original Safecast import page
	safecastLink := ""
	if upload.Source == "safecast-api" && upload.SourceID != "" {
		safecastLink = fmt.Sprintf(` <a href="https://api.safecast.org/en-US/bgeigie_imports/%s" target="_blank" title="View on Safecast" style="color:var(--link-color);font-size:0.8em;">↗ Safecast</a>`, upload.SourceID)
	}

	return fmt.Sprintf(`
			<tr>
				<td class="checkbox-col"><input type="checkbox" class="track-checkbox" value="%s" onchange="updateDeleteButton()"></td>
				<td>%d</td>
				<td class="filename">%s</td>
				<td>%s</td>
				<td class="trackid"><a href="/trackid/%s">%s</a></td>
				<td>%s</td>
				<td class="datetime">%s</td>
				<td class="filesize">%s</td>
				<td class="source" data-source-id="%s">%s</td>
				<td>%s</td>
				<td>%s</td>
				<td class="datetime">%s</td>
				<td class="comment" title="%s">%s</td>
				<td>
					<button class="edit-btn" onclick="openEditUpload(%d,'%s','%s','%s','%s','%s','%s')">Edit</button>
					<button class="delete-btn" onclick="deleteTrack('%s')">Delete</button>%s
				</td>
			</tr>`,
		upload.TrackID,
		upload.ID,
		upload.Filename,
		upload.FileType,
		upload.TrackID, upload.TrackID,
		detectorDisplay,
		recordingDate,
		fileSize,
		sourceIDNumeric, sourceDisplay,
		userDisplay,
		upload.UploadIP,
		uploadTime,
		escapedComment, truncateString(upload.Comment, 40),
		upload.ID, escapedFilename, escapedUsername, escapedDetector, recordingDateISO, escapedNotes, escapedComment,
		upload.TrackID,
		safecastLink,
	)
}

// checkAdminAuth verifies admin access via URL password or session-based admin user.
// Returns (authorized bool, password string for backwards compat in templates).
// If not authorized, it writes an HTTP error response and returns false.
func checkAdminAuth(w http.ResponseWriter, r *http.Request) (bool, string) {
	// First check if user is authenticated via session and is admin
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		// For backwards compatibility, return an empty password value for templates
		// When using session auth, we don't need password in URLs
		return true, ""
	}

	// Fall back to URL password authentication
	if *adminPassword == "" {
		http.Error(w, "Admin endpoints are disabled - please login as admin", http.StatusForbidden)
		return false, ""
	}

	password := r.URL.Query().Get("password")
	if password != *adminPassword {
		http.Error(w, "Unauthorized - Invalid password or not logged in as admin", http.StatusUnauthorized)
		return false, ""
	}

	return true, password
}

// adminUploadsHandler lists all file uploads with metadata.
// GET /api/admin/uploads?password=xxx&limit=100
// adminUploadsHandler lists all file uploads with metadata and search functionality
//
// @Summary     Admin uploads dashboard data
// @Description Returns paginated upload records for admin users.
// @Tags        admin
// @Produce     html
// @Param       limit query int false "Page size"
// @Param       page query int false "Page number"
// @Param       user_id query string false "Filter by internal user ID"
// @Param       search query string false "Free-text search"
// @Success     200 {string} string "HTML admin page"
// @Failure     401 {string} string "Unauthorized"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/uploads [get]
func adminUploadsHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent CloudFront from caching this dynamic admin page
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	authorized, password := checkAdminAuth(w, r)
	if !authorized {
		return
	}
	_ = password // Used in HTML template generation below

	if db == nil || db.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	// Get limit parameter (page size)
	limit := 500 // Default to 500 per page
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Get page parameter
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get user_id filter parameter
	userID := r.URL.Query().Get("user_id")

	// Get search parameter
	search := r.URL.Query().Get("search")

	ctx := r.Context()

	// Get total count for pagination
	totalCount, err := db.CountUploads(ctx, userID, search)
	if err != nil {
		log.Printf("Error counting uploads: %v", err)
		http.Error(w, "Failed to count uploads", http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := (totalCount + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	// Fetch current page of uploads
	uploads, err := db.GetUploadsPaginated(ctx, limit, offset, userID, search)
	if err != nil {
		log.Printf("Error fetching uploads: %v", err)
		http.Error(w, "Failed to fetch uploads", http.StatusInternalServerError)
		return
	}

	// Return HTML table
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html>
<head>
	<title>Admin - File Uploads</title>
	<script>
	(function() {
		var saved = localStorage.getItem('safecastDocTheme');
		var prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		document.documentElement.setAttribute('data-theme', saved || (prefersDark ? 'dark' : 'light'));
	})();
	</script>

	<!-- favicon -->
	<link rel="apple-touch-icon" sizes="180x180" href="/static/images/apple-touch-icon.png">
	<link rel="icon" type="image/png" sizes="32x32" href="/static/images/favicon-32x32.png">
	<link rel="icon" type="image/png" sizes="16x16" href="/static/images/favicon-16x16.png">
	<link rel="manifest" href="/static/images/site.webmanifest">

	<style>
		:root {
			--bg-primary: #f5f5f5;
			--bg-card: white;
			--text-primary: #333;
			--text-secondary: #666;
			--text-muted: #999;
			--border-color: #ddd;
			--link-color: #0066cc;
			--shadow: 0 1px 3px rgba(0,0,0,0.1);
			--th-bg: #424242;
			--hover-bg: #f9f9f9;
		}
		@media (prefers-color-scheme: dark) {
			:root {
				--bg-primary: #1a1a1a;
				--bg-card: #2b2b2b;
				--text-primary: #eee;
				--text-secondary: #aaa;
				--text-muted: #777;
				--border-color: #444;
				--link-color: #90caf9;
				--shadow: 0 1px 3px rgba(255,255,255,0.1);
				--th-bg: #616161;
				--hover-bg: #333;
				color-scheme: dark;
			}
		}
		:root[data-theme='light'] {
			--bg-primary: #f5f5f5; --bg-card: #fff; --text-primary: #333;
			--text-secondary: #666; --text-muted: #999; --border-color: #ddd;
			--link-color: #0066cc; --shadow: 0 1px 3px rgba(0,0,0,0.1);
			--hover-bg: #f9f9f9; --th-bg: #424242; color-scheme: light;
		}
		:root[data-theme='dark'] {
			--bg-primary: #1a1a1a; --bg-card: #2b2b2b; --text-primary: #eee;
			--text-secondary: #aaa; --text-muted: #777; --border-color: #444;
			--link-color: #90caf9; --shadow: 0 1px 3px rgba(255,255,255,0.07);
			--hover-bg: #333; --th-bg: #616161; color-scheme: dark;
		}
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Arial, sans-serif; margin: 0; background: var(--bg-primary); color: var(--text-primary); transition: background 0.2s, color 0.2s; }
		h1 { color: var(--text-primary); }
		.top-nav { background: #1a3a5c; color: #fff; padding: 0 24px; display: flex; align-items: center; gap: 16px; height: 52px; box-shadow: 0 2px 6px rgba(0,0,0,0.3); position: sticky; top: 0; z-index: 100; }
		.nav-logo { display: flex; align-items: center; gap: 8px; text-decoration: none; color: #fff; font-weight: 700; font-size: 17px; white-space: nowrap; }
		.nav-logo img { height: 28px; width: 28px; object-fit: contain; }
		.nav-sep { color: rgba(255,255,255,0.3); font-size: 18px; }
		.nav-title { font-size: 15px; font-weight: 600; color: #d0e8ff; white-space: nowrap; }
		.nav-spacer { flex: 1; }
		.back-link { color: #afd4f5; text-decoration: none; font-size: 13px; white-space: nowrap; }
		.back-link:hover { color: #fff; }
		#theme-toggle { background: rgba(255,255,255,0.12); color: #fff; border: 1px solid rgba(255,255,255,0.25); border-radius: 6px; padding: 6px 14px; font-size: 13px; font-weight: 600; cursor: pointer; white-space: nowrap; transition: background 0.2s; }
		#theme-toggle:hover { background: rgba(255,255,255,0.22); }
		.page-content { padding: 20px 24px; }
		.nav { background: var(--bg-card); padding: 15px; margin-bottom: 20px; border-radius: 5px; box-shadow: var(--shadow); display: flex; align-items: center; justify-content: space-between; }
		.nav-left { display: flex; align-items: center; gap: 15px; }
		.nav a { color: var(--link-color); text-decoration: none; }
		.nav a:hover { text-decoration: underline; }
		.summary { background: var(--bg-card); padding: 15px; margin-bottom: 20px; border-radius: 5px; box-shadow: var(--shadow); }
		table { border-collapse: collapse; width: 100%; background: var(--bg-card); box-shadow: var(--shadow); table-layout: auto; }
		th { background: var(--th-bg); color: white; padding: 12px; text-align: left; font-weight: 600; white-space: nowrap; position: relative; overflow: hidden; }
		td { padding: 6px 8px; border-bottom: 1px solid var(--border-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
		.resize-handle { position: absolute; right: 0; top: 0; width: 5px; height: 100%; cursor: col-resize; background: transparent; z-index: 1; }
		.resize-handle:hover, .resize-handle.active { background: rgba(255,255,255,0.3); }
		tr:hover { background: var(--hover-bg); }
		.empty { text-align: center; padding: 40px; color: var(--text-muted); font-style: italic; }
		.trackid { font-family: monospace; color: var(--link-color); white-space: nowrap; }
		.filename { color: var(--text-primary); font-weight: 500; max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
		.filesize { color: var(--text-secondary); white-space: nowrap; }
		.datetime { color: var(--text-secondary); font-size: 0.9em; white-space: nowrap; }
		.comment { max-width: 220px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: var(--text-secondary); font-size: 0.9em; }
		.delete-btn { background: #f44336; color: white; border: none; padding: 5px 10px; border-radius: 3px; cursor: pointer; font-size: 0.85em; }
		.delete-btn:hover { background: #d32f2f; }
		.edit-btn { background: #FF9800; color: white; border: none; padding: 5px 10px; border-radius: 3px; cursor: pointer; font-size: 0.85em; margin-right: 4px; }
		.edit-btn:hover { background: #F57C00; }
		.modal-overlay { display: none; position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 1000; align-items: center; justify-content: center; }
		.modal-overlay.open { display: flex; }
		.modal { background: var(--bg-card); padding: 24px; border-radius: 8px; width: 480px; max-width: 95vw; box-shadow: 0 8px 32px rgba(0,0,0,0.3); }
		.modal h3 { margin: 0 0 16px; color: var(--text-primary); }
		.modal label { display: block; margin-bottom: 4px; color: var(--text-secondary); font-size: 0.9em; font-weight: 500; }
		.modal input, .modal textarea { width: 100%; padding: 8px 10px; border: 1px solid var(--border-color); border-radius: 4px; background: var(--bg-primary); color: var(--text-primary); font-size: 0.95em; box-sizing: border-box; margin-bottom: 12px; }
		.modal textarea { min-height: 80px; resize: vertical; font-family: inherit; }
		.modal-btns { display: flex; gap: 8px; justify-content: flex-end; margin-top: 4px; }
		.modal-save-btn { background: #4CAF50; color: white; border: none; padding: 8px 20px; border-radius: 4px; cursor: pointer; font-size: 0.95em; font-weight: 500; }
		.modal-save-btn:hover { background: #388E3C; }
		.modal-cancel-btn { background: var(--bg-primary); color: var(--text-primary); border: 1px solid var(--border-color); padding: 8px 20px; border-radius: 4px; cursor: pointer; font-size: 0.95em; }
		.modal-cancel-btn:hover { background: var(--hover-bg); }
		.delete-selected-btn { background: #f44336; color: white; border: none; padding: 10px 20px; border-radius: 3px; cursor: pointer; font-size: 1em; margin-left: 10px; }
		.delete-selected-btn:hover { background: #d32f2f; }
		.delete-selected-btn:disabled { background: #ccc; cursor: not-allowed; }
		.view-selected-btn { background: #2196F3; color: white; border: none; padding: 10px 20px; border-radius: 3px; cursor: pointer; font-size: 1em; margin-left: 10px; }
		.view-selected-btn:hover { background: #1976D2; }
		.view-selected-btn:disabled { background: #ccc; cursor: not-allowed; }
		.checkbox-col { width: 40px; text-align: center; }
		.sortable { cursor: pointer; user-select: none; position: relative; padding-right: 20px; }
		.sortable:hover { background: rgba(255,255,255,0.1); }
		.sortable::after { content: '⇅'; position: absolute; right: 8px; opacity: 0.5; }
		.sortable.asc::after { content: '▲'; opacity: 1; }
		.sortable.desc::after { content: '▼'; opacity: 1; }
		.filter-input { width: 100%; padding: 4px 8px; border: 1px solid var(--border-color); border-radius: 3px; background: var(--bg-card); color: var(--text-primary); font-size: 0.85em; box-sizing: border-box; }
		.filter-row th { background: var(--bg-card); padding: 8px 12px; }
		.import-form { background: var(--bg-card); padding: 20px; margin-bottom: 20px; border-radius: 5px; box-shadow: var(--shadow); }
		.import-form h3 { margin-top: 0; color: var(--text-primary); }
		.form-row { display: flex; gap: 15px; align-items: flex-end; margin-top: 15px; }
		.form-group { min-width: 200px; }
		.form-group label { display: block; margin-bottom: 5px; color: var(--text-secondary); font-size: 0.9em; font-weight: 500; }
		.form-group input[type="date"] { width: 200px; padding: 8px 12px; border: 1px solid var(--border-color); border-radius: 3px; background: var(--bg-card); color: var(--text-primary); font-size: 1em; }
		.import-btn { background: #2196F3; color: white; border: none; padding: 10px 24px; border-radius: 3px; cursor: pointer; font-size: 1em; font-weight: 500; }
		.import-btn:hover { background: #1976D2; }
		.import-btn:disabled { background: #ccc; cursor: not-allowed; }
		.import-status { margin-top: 15px; padding: 12px; border-radius: 3px; display: none; }
		.import-status.success { background: #4CAF50; color: white; }
		.import-status.error { background: #f44336; color: white; }
		.import-status.info { background: #2196F3; color: white; }
		/* Admin tab bar */
		.admin-tabs { display: flex; gap: 2px; margin-top: 10px; margin-bottom: 10px; background: var(--border-color); border-radius: 8px; overflow: hidden; }
		.admin-tabs a, .admin-tabs span { padding: 10px 20px; text-decoration: none; color: var(--text-secondary); background: var(--bg-card); font-weight: 500; font-size: 0.95em; transition: background 0.2s; }
		.admin-tabs a:hover { background: var(--hover-bg); color: var(--text-primary); }
		.admin-tabs a.active { background: #2196F3; color: white; }
		.admin-tabs span.disabled { color: var(--text-muted); cursor: not-allowed; font-style: italic; }
	</style>
</head>
<body>
<nav class="top-nav">
  <a href="/" class="nav-logo">
    <img src="/static/images/safecast-logo-squared.png" alt="Safecast">
    Safecast
  </a>
  <span class="nav-sep">|</span>
  <span class="nav-title">File Uploads Administration</span>
  <span class="nav-spacer"></span>
  <a href="/" class="back-link">&#8592; Back to Map</a>
  <button id="theme-toggle" onclick="toggleTheme()">&#127769; Dark Mode</button>
</nav>
<div class="page-content">
	<div class="admin-tabs">
		<a href="/admin/users` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">Users</a>
		<a href="/admin/uploads` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `" class="active">Uploads</a>
		<a href="/admin/mcp` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">MCP Analytics</a>
		<a href="/admin/realtime` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">Realtime</a>
		<a href="/admin/translations` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">Translations</a>
		<a href="/admin/tour` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">Tour Steps</a>
		<a href="/admin/ai-hints` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">AI Hints</a>
	</div>
	<div class="nav">
		<div class="nav-left">
			<button class="view-selected-btn" id="viewSelectedBtn" onclick="viewSelected()" disabled>View Selected on Map</button>
			<button class="delete-selected-btn" id="deleteSelectedBtn" onclick="deleteSelected()" disabled>Delete Selected</button>
			<button class="import-btn" onclick="importSafecastMeta()" id="importSafecastBtn" style="margin-left:10px;">Import Safecast API Metadata</button>
		</div>
		<div id="importMetaStatus" style="font-size:0.85em;color:var(--text-secondary);margin-top:4px;"></div>
	</div>
	<div class="import-form">
		<h3>Import from Safecast API</h3>
		<div class="form-row">
			<div class="form-group">
				<label for="startDate">Start Date</label>
				<input type="date" id="startDate" required>
			</div>
			<div class="form-group">
				<label for="endDate">End Date</label>
				<input type="date" id="endDate" required>
			</div>
			<button class="import-btn" onclick="importFromAPI()" id="importBtn">Import</button>
		</div>
		<div class="import-status" id="importStatus"></div>
	</div>
	<div class="summary">
		<strong>Total Uploads:</strong> ` + strconv.Itoa(totalCount) + ` files
		<span style="margin-left: 20px;">
			<strong>Page ` + strconv.Itoa(page) + ` of ` + strconv.Itoa(totalPages) + `</strong>
			(showing ` + strconv.Itoa(len(uploads)) + ` uploads)
		</span>
		<span style="margin-left: 20px;">
			<label for="limitSelect"><strong>Per page:</strong></label>
			<select id="limitSelect" onchange="changeLimit()" style="margin-left: 5px; padding: 4px 8px; border-radius: 4px; border: 1px solid var(--border-color); background: var(--bg-card); color: var(--text-primary);">
				<option value="100"` + func() string {
		if limit == 100 {
			return " selected"
		}
		return ""
	}() + `>100</option>
				<option value="250"` + func() string {
		if limit == 250 {
			return " selected"
		}
		return ""
	}() + `>250</option>
				<option value="500"` + func() string {
		if limit == 500 {
			return " selected"
		}
		return ""
	}() + `>500</option>
				<option value="1000"` + func() string {
		if limit == 1000 {
			return " selected"
		}
		return ""
	}() + `>1000</option>
			</select>
		</span>
		<span style="margin-left: 20px;">
			<label for="searchInput"><strong>Search:</strong></label>
			<input type="text" id="searchInput" value="` + search + `" placeholder="Search all fields..." autocomplete="off" style="margin-left: 5px; padding: 4px 8px; border-radius: 4px; border: 1px solid var(--border-color); background: var(--bg-card); color: var(--text-primary); width: 200px;" onkeypress="if(event.key === 'Enter') performSearch()">
			<button onclick="performSearch()" style="margin-left: 5px; padding: 4px 12px; border-radius: 4px; border: 1px solid var(--border-color); background: var(--link-color); color: white; cursor: pointer;">🔍</button>
			` + func() string {
		if search != "" {
			return `<button onclick="clearSearch()" style="margin-left: 5px; padding: 4px 12px; border-radius: 4px; border: 1px solid var(--border-color); background: var(--bg-card); color: var(--text-primary); cursor: pointer;">Clear</button>`
		}
		return ""
	}() + `
		</span>`

	if userID != "" {
		params := []string{}
		if password != "" {
			params = append(params, "password="+password)
		}
		if search != "" {
			params = append(params, "search="+url.QueryEscape(search))
		}
		clearFilterURL := "/api/admin/uploads"
		if len(params) > 0 {
			clearFilterURL += "?" + strings.Join(params, "&")
		}
		html += ` | <strong>Filtered by User ID:</strong> ` + userID + ` <a href="` + clearFilterURL + `">[Clear Filter]</a>`
	}

	if search != "" {
		html += ` | <strong>Search:</strong> "` + search + `"`
	}

	// Add pagination controls inline in the summary
	html += `<div style="margin-top: 10px;">`

	// Helper function to build query parameters
	buildURL := func(pageNum int) string {
		params := []string{}
		if password != "" {
			params = append(params, "password="+password)
		}
		params = append(params, "page="+strconv.Itoa(pageNum))
		params = append(params, "limit="+strconv.Itoa(limit))
		if userID != "" {
			params = append(params, "user_id="+userID)
		}
		if search != "" {
			params = append(params, "search="+url.QueryEscape(search))
		}
		return "?" + strings.Join(params, "&")
	}

	// Previous button
	if page > 1 {
		html += `<a href="` + buildURL(page-1) + `" class="page-btn">&laquo; Previous</a>`
	} else {
		html += `<span class="page-btn disabled">&laquo; Previous</span>`
	}

	// Page numbers
	startPage := page - 2
	if startPage < 1 {
		startPage = 1
	}
	endPage := startPage + 4
	if endPage > totalPages {
		endPage = totalPages
		startPage = endPage - 4
		if startPage < 1 {
			startPage = 1
		}
	}

	// First page
	if startPage > 1 {
		html += `<a href="` + buildURL(1) + `" class="page-btn">1</a>`
		if startPage > 2 {
			html += `<span class="page-btn disabled">...</span>`
		}
	}

	// Page range
	for i := startPage; i <= endPage; i++ {
		if i == page {
			html += `<span class="page-btn active">` + strconv.Itoa(i) + `</span>`
		} else {
			html += `<a href="` + buildURL(i) + `" class="page-btn">` + strconv.Itoa(i) + `</a>`
		}
	}

	// Last page
	if endPage < totalPages {
		if endPage < totalPages-1 {
			html += `<span class="page-btn disabled">...</span>`
		}
		html += `<a href="` + buildURL(totalPages) + `" class="page-btn">` + strconv.Itoa(totalPages) + `</a>`
	}

	// Next button
	if page < totalPages {
		html += `<a href="` + buildURL(page+1) + `" class="page-btn">Next &raquo;</a>`
	} else {
		html += `<span class="page-btn disabled">Next &raquo;</span>`
	}

	html += `</div>
	</div>`

	if len(uploads) == 0 {
		html += `<div class="empty">No uploads found. Upload a spectrum file (.n42 or .spe) to see it appear here.</div>`
	} else {
		html += `
	<div style="overflow-x: auto;">
	<table id="uploadsTable">
		<thead>
			<tr>
				<th class="checkbox-col"><input type="checkbox" id="selectAll" onchange="toggleSelectAll(this)"></th>
				<th class="sortable" onclick="sortTable(1)" data-type="number">ID<span class="resize-handle"></span></th>
				<th class="sortable" onclick="sortTable(2)" data-type="text">Filename<span class="resize-handle"></span></th>
				<th class="sortable" onclick="sortTable(3)" data-type="text">Type<span class="resize-handle"></span></th>
				<th class="sortable" onclick="sortTable(4)" data-type="text">Track ID<span class="resize-handle"></span></th>
				<th class="sortable" onclick="sortTable(5)" data-type="text">Detector<span class="resize-handle"></span></th>
				<th class="sortable" onclick="sortTable(6)" data-type="date">Recording Date<span class="resize-handle"></span></th>
				<th class="sortable" onclick="sortTable(7)" data-type="text">Size<span class="resize-handle"></span></th>
				<th class="sortable" onclick="sortTable(8)" data-type="text">Source<span class="resize-handle"></span></th>
				<th class="sortable" onclick="sortTable(9)" data-type="text">User<span class="resize-handle"></span></th>
				<th class="sortable" onclick="sortTable(10)" data-type="text">Upload IP<span class="resize-handle"></span></th>
				<th class="sortable" onclick="sortTable(11)" data-type="date">Upload Time<span class="resize-handle"></span></th>
				<th class="sortable" onclick="sortTable(12)" data-type="text">Comment<span class="resize-handle"></span></th>
				<th>Actions</th>
			</tr>
			<tr class="filter-row">
				<th></th>
				<th><input type="text" class="filter-input" placeholder="ID..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Filter filename..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Type..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Filter Track ID..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Filter Detector..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Filter date..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Size..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Source..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="User..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="IP..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Filter date..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Filter comment..." onkeyup="filterTable()"></th>
				<th></th>
			</tr>
		</thead>
		<tbody id="uploadsTableBody">`

		// Parallel processing for formatting upload records
		numWorkers := 8
		if len(uploads) < 50 {
			// For small datasets, use sequential processing (faster due to no overhead)
			numWorkers = 1
		}
		if len(uploads) < numWorkers {
			numWorkers = len(uploads)
		}

		// Results slice to store formatted HTML in order
		rows := make([]string, len(uploads))

		if numWorkers == 1 {
			// Sequential processing for small datasets
			for i, upload := range uploads {
				rows[i] = formatUploadRow(upload, password)
			}
		} else {
			// Parallel processing for large datasets
			var wg sync.WaitGroup
			batchSize := (len(uploads) + numWorkers - 1) / numWorkers

			for w := 0; w < numWorkers; w++ {
				wg.Add(1)
				start := w * batchSize
				end := start + batchSize
				if end > len(uploads) {
					end = len(uploads)
				}

				go func(start, end int) {
					defer wg.Done()
					for i := start; i < end; i++ {
						rows[i] = formatUploadRow(uploads[i], password)
					}
				}(start, end)
			}
			wg.Wait()
		}

		// Concatenate all rows
		for _, row := range rows {
			html += row
		}

		html += `
		</tbody>
	</table>
	</div>`
	}

	html += `
	<script>
		// Column resize handles
		(function() {
			let resizing = false;
			document.querySelectorAll('#uploadsTable .resize-handle').forEach(handle => {
				handle.addEventListener('mousedown', function(e) {
					e.preventDefault(); e.stopPropagation();
					resizing = true;
					const th = this.parentElement;
					const startX = e.pageX, startWidth = th.offsetWidth;
					this.classList.add('active');
					const onMove = e => { th.style.width = Math.max(40, startWidth + (e.pageX - startX)) + 'px'; };
					const onUp = () => {
						this.classList.remove('active');
						document.removeEventListener('mousemove', onMove);
						document.removeEventListener('mouseup', onUp);
						setTimeout(() => { resizing = false; }, 50);
					};
					document.addEventListener('mousemove', onMove);
					document.addEventListener('mouseup', onUp);
				});
			});
			document.querySelectorAll('#uploadsTable .sortable').forEach(th => {
				th.addEventListener('click', e => { if (resizing) e.stopImmediatePropagation(); });
			});
		})();

		// Apply theme from sessionStorage to match map preference
		(function() {
			const media = window.matchMedia('(prefers-color-scheme: dark)');
			const storedTheme = sessionStorage.getItem('themePreference');
			const theme = storedTheme ? storedTheme : (media.matches ? 'dark' : 'light');
			document.documentElement.dataset.theme = theme;
		})();

		function toggleSelectAll(checkbox) {
			const checkboxes = document.querySelectorAll('.track-checkbox');
			checkboxes.forEach(cb => cb.checked = checkbox.checked);
			updateDeleteButton();
		}

		function updateDeleteButton() {
			const checkboxes = document.querySelectorAll('.track-checkbox:checked');
			const deleteBtn = document.getElementById('deleteSelectedBtn');
			const viewBtn = document.getElementById('viewSelectedBtn');
			const count = checkboxes.length;
			
			deleteBtn.disabled = count === 0;
			deleteBtn.textContent = count > 0 ? 'Delete Selected (' + count + ')' : 'Delete Selected';
			
			viewBtn.disabled = count === 0;
			viewBtn.textContent = count > 0 ? 'View Selected on Map (' + count + ')' : 'View Selected on Map';
		}

		function viewSelected() {
			const checkboxes = document.querySelectorAll('.track-checkbox:checked');
			const trackIDs = Array.from(checkboxes).map(cb => cb.value);

			if (trackIDs.length === 0) return;

			// Navigate to map with multiple tracks
			const tracksParam = trackIDs.join(',');
			window.location.href = '/tracks/' + encodeURIComponent(tracksParam);
		}

		function deleteSelected() {
			const checkboxes = document.querySelectorAll('.track-checkbox:checked');
			const trackIDs = Array.from(checkboxes).map(cb => cb.value);

			if (trackIDs.length === 0) return;

			if (!confirm('Delete ' + trackIDs.length + ' track(s) and all associated data?')) {
				return;
			}

			const password = new URLSearchParams(window.location.search).get('password');

			fetch('/api/admin/delete-multiple', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ password: password, trackIDs: trackIDs })
			})
			.then(r => r.json())
			.then(data => {
				if (data.status === 'success') {
					alert('Successfully deleted ' + data.deleted + ' track(s)');
					window.location.reload();
				} else if (data.status === 'partial') {
					alert('Partially deleted ' + data.deleted + ' of ' + trackIDs.length + ' track(s). Some errors occurred.');
					window.location.reload();
				} else {
					alert('Error: ' + (data.error || 'Unknown error'));
				}
			})
			.catch(err => alert('Error: ' + err));
		}

		function deleteTrack(trackID) {
			if (!confirm('Delete track ' + trackID + ' and all associated data?')) {
				return;
			}

			const password = new URLSearchParams(window.location.search).get('password');

			fetch('/api/admin/delete', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ password: password, trackID: trackID })
			})
			.then(r => r.json())
			.then(data => {
				if (data.status === 'success') {
					alert('Track deleted successfully');
					window.location.reload();
				} else {
					alert('Error: ' + (data.error || 'Unknown error'));
				}
			})
			.catch(err => alert('Error: ' + err));
		}

		// Sorting functionality
		let sortDirection = {};
		function sortTable(columnIndex) {
			const table = document.getElementById('uploadsTable');
			const tbody = document.getElementById('uploadsTableBody');
			const rows = Array.from(tbody.querySelectorAll('tr'));
			const header = table.querySelector('thead tr:first-child th:nth-child(' + (columnIndex + 1) + ')');
			const dataType = header.getAttribute('data-type');

			// Toggle sort direction
			const currentDir = sortDirection[columnIndex] || 'none';
			sortDirection[columnIndex] = currentDir === 'asc' ? 'desc' : 'asc';

			// Remove sort classes from all headers
			table.querySelectorAll('.sortable').forEach(h => {
				h.classList.remove('asc', 'desc');
			});

			// Add sort class to current header
			header.classList.add(sortDirection[columnIndex]);

			// Sort rows
			rows.sort((a, b) => {
				let aVal = a.cells[columnIndex].textContent.trim();
				let bVal = b.cells[columnIndex].textContent.trim();

				// Special handling for Source column - sort by numeric import ID
				if (columnIndex === 7) {
					const aSourceID = a.cells[columnIndex].getAttribute('data-source-id');
					const bSourceID = b.cells[columnIndex].getAttribute('data-source-id');
					aVal = parseInt(aSourceID) || 0;
					bVal = parseInt(bSourceID) || 0;
					return sortDirection[columnIndex] === 'asc' ? aVal - bVal : bVal - aVal;
				}

				// Numeric comparison
				if (dataType === 'number') {
					aVal = parseInt(aVal) || 0;
					bVal = parseInt(bVal) || 0;
					return sortDirection[columnIndex] === 'asc' ? aVal - bVal : bVal - aVal;
				}

				// Date comparison
				if (dataType === 'date') {
					aVal = new Date(aVal).getTime();
					bVal = new Date(bVal).getTime();
					return sortDirection[columnIndex] === 'asc' ? aVal - bVal : bVal - aVal;
				}

				// Text comparison
				if (sortDirection[columnIndex] === 'asc') {
					return aVal.localeCompare(bVal);
				} else {
					return bVal.localeCompare(aVal);
				}
			});

			// Reappend sorted rows
			rows.forEach(row => tbody.appendChild(row));
		}

		// Filtering functionality
		function filterTable() {
			const table = document.getElementById('uploadsTable');
			const tbody = document.getElementById('uploadsTableBody');
			const filters = table.querySelectorAll('.filter-input');
			const rows = tbody.querySelectorAll('tr');

			rows.forEach(row => {
				let show = true;
				filters.forEach((filter, index) => {
					const filterValue = filter.value.toLowerCase();
					if (filterValue) {
						const cellIndex = index + 1; // +1 because first column is checkbox
						const cell = row.cells[cellIndex];
						if (cell) {
							let cellText = cell.textContent.toLowerCase();
							// For Source column, also check data-source-id attribute
							if (cellIndex === 7) {
								const sourceID = cell.getAttribute('data-source-id') || '';
								cellText = cellText + ' ' + sourceID;
							}
							if (!cellText.includes(filterValue)) {
								show = false;
							}
						}
					}
				});
				row.style.display = show ? '' : 'none';
			});

			// Update delete button after filtering
			updateDeleteButton();
		}

		// Change limit and reload page, reset to page 1
		function changeLimit() {
			const limit = document.getElementById('limitSelect').value;
			const url = new URL(window.location.href);
			url.searchParams.set('limit', limit);
			url.searchParams.set('page', '1');
			window.location.href = url.toString();
		}

		// Perform search
		function performSearch() {
			const searchValue = document.getElementById('searchInput').value;
			const url = new URL(window.location.href);
			if (searchValue.trim()) {
				url.searchParams.set('search', searchValue.trim());
			} else {
				url.searchParams.delete('search');
			}
			url.searchParams.set('page', '1'); // Reset to page 1 when searching
			window.location.href = url.toString();
		}

		// Clear search
		function clearSearch() {
			const url = new URL(window.location.href);
			url.searchParams.delete('search');
			url.searchParams.set('page', '1'); // Reset to page 1 when clearing
			window.location.href = url.toString();
		}

		// Import from Safecast API
		async function importFromAPI() {
			const startDate = document.getElementById('startDate').value;
			const endDate = document.getElementById('endDate').value;
			const status = document.getElementById('importStatus');
			const btn = document.getElementById('importBtn');

			if (!startDate || !endDate) {
				status.className = 'import-status error';
				status.style.display = 'block';
				status.textContent = 'Please select both start and end dates';
				return;
			}

			if (new Date(startDate) > new Date(endDate)) {
				status.className = 'import-status error';
				status.style.display = 'block';
				status.textContent = 'Start date must be before end date';
				return;
			}

			btn.disabled = true;
			status.className = 'import-status info';
			status.style.display = 'block';
			status.textContent = 'Starting import...';

			const password = new URLSearchParams(window.location.search).get('password');

			try {
				// Start the import with streaming response
				const response = await fetch('/api/admin/import-from-safecast?password=' + password, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ start_date: startDate, end_date: endDate })
				});

				if (!response.ok) {
					throw new Error('Server returned error: ' + response.status);
				}

				// Read the SSE stream
				const reader = response.body.getReader();
				const decoder = new TextDecoder();
				let buffer = '';

				while (true) {
					const { done, value } = await reader.read();
					if (done) break;

					// Decode the chunk and add to buffer
					buffer += decoder.decode(value, { stream: true });

					// Process complete SSE messages
					let lines = buffer.split('\n');
					buffer = lines.pop(); // Keep incomplete line in buffer

					let currentEvent = null;
					let currentData = '';

					for (let line of lines) {
						if (line.startsWith('event:')) {
							currentEvent = line.substring(6).trim();
						} else if (line.startsWith('data:')) {
							currentData = line.substring(5).trim();
						} else if (line === '') {
							// Empty line marks end of message
							if (currentData) {
								try {
									const data = JSON.parse(currentData);

									if (currentEvent === 'done') {
										// Import complete
										status.className = 'import-status success';
										status.textContent = 'Import complete! ' + data.imported + ' files imported, ' +
										                     data.skipped + ' skipped, ' + data.errors + ' errors';
										setTimeout(() => window.location.reload(), 2000);
									} else {
										// Progress update
										let message = data.message;
										if (data.total > 0) {
											message += ' (' + data.imported + ' imported, ' + data.skipped + ' skipped, ' + data.errors + ' errors)';
										}
										status.className = 'import-status info';
										status.textContent = message;
									}
								} catch (e) {
									console.error('Error parsing SSE data:', e, currentData);
								}

								currentEvent = null;
								currentData = '';
							}
						}
					}
				}
			} catch (err) {
				console.error('Import error:', err);
				status.className = 'import-status error';
				status.textContent = 'Error: ' + err.message;
				btn.disabled = false;
			}
		}

		let _editUploadID = 0;

		function openEditUpload(id, filename, username, detector, recordingDate, notes, comment) {
			_editUploadID = id;
			document.getElementById('editUploadFilename').value = filename;
			document.getElementById('editUploadUsername').value = username;
			document.getElementById('editUploadDetector').value = detector;
			document.getElementById('editUploadRecordingDate').value = recordingDate;
			document.getElementById('editUploadNotes').value = notes;
			document.getElementById('editUploadComment').value = comment || '';
			document.getElementById('editUploadModal').classList.add('open');
		}

		function closeEditUpload() {
			document.getElementById('editUploadModal').classList.remove('open');
		}

		function importSafecastMeta() {
			if (!confirm('Fetch track names and comments from the old Safecast API for all imported tracks? This may take a minute.')) return;
			const password = new URLSearchParams(window.location.search).get('password');
			const btn = document.getElementById('importSafecastBtn');
			const status = document.getElementById('importMetaStatus');
			btn.disabled = true;
			btn.textContent = 'Importing…';
			status.textContent = 'Fetching metadata from api.safecast.org…';
			fetch('/api/admin/tracks/import-safecast', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ password })
			})
			.then(r => r.json())
			.then(data => {
				btn.disabled = false;
				btn.textContent = 'Import Safecast API Metadata';
				if (data.ok) {
					status.textContent = 'Done: ' + data.updated + ' track(s) updated.';
					setTimeout(() => window.location.reload(), 1500);
				} else {
					status.textContent = 'Error: ' + (data.error || 'unknown');
				}
			})
			.catch(err => {
				btn.disabled = false;
				btn.textContent = 'Import Safecast API Metadata';
				status.textContent = 'Error: ' + err;
			});
		}

		async function saveEditUpload() {
			const password = new URLSearchParams(window.location.search).get('password');
			const body = {
				password: password,
				upload_id: _editUploadID,
				filename: document.getElementById('editUploadFilename').value,
				username: document.getElementById('editUploadUsername').value,
				detector: document.getElementById('editUploadDetector').value,
				recording_date: document.getElementById('editUploadRecordingDate').value,
				notes: document.getElementById('editUploadNotes').value,
				comment: document.getElementById('editUploadComment').value
			};
			try {
				const resp = await fetch('/api/admin/uploads/update', {
					method: 'PUT',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(body)
				});
				if (!resp.ok) {
					const text = await resp.text();
					alert('Save failed: ' + text);
					return;
				}
				closeEditUpload();
				window.location.reload();
			} catch (err) {
				alert('Error: ' + err);
			}
		}
	</script>

	<!-- Edit Upload Modal -->
	<div class="modal-overlay" id="editUploadModal" onclick="if(event.target===this)closeEditUpload()">
		<div class="modal">
			<h3>Edit Upload Metadata</h3>
			<label>Filename</label>
			<input type="text" id="editUploadFilename">
			<label>Username</label>
			<input type="text" id="editUploadUsername" placeholder="uploader username">
			<label>Detector</label>
			<input type="text" id="editUploadDetector" placeholder="e.g. LND7317">
			<label>Recording Date</label>
			<input type="date" id="editUploadRecordingDate">
			<label>Notes</label>
			<textarea id="editUploadNotes" placeholder="Admin notes..."></textarea>
			<label>Comment</label>
			<textarea id="editUploadComment" placeholder="User comment from Safecast API..."></textarea>
			<div class="modal-btns">
				<button class="modal-cancel-btn" onclick="closeEditUpload()">Cancel</button>
				<button class="modal-save-btn" onclick="saveEditUpload()">Save</button>
			</div>
		</div>
	</div>
</div><!-- .page-content -->
<script>
(function() {
	function applyLabel() {
		var btn = document.getElementById('theme-toggle');
		if (btn) btn.textContent = document.documentElement.getAttribute('data-theme') === 'dark' ? '\u2600\uFE0F Light Mode' : '\uD83C\uDF19 Dark Mode';
	}
	applyLabel();
	window.toggleTheme = function() {
		var next = document.documentElement.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
		document.documentElement.setAttribute('data-theme', next);
		localStorage.setItem('safecastDocTheme', next);
		applyLabel();
	};
})();
</script>
</body>
</html>`

	fmt.Fprint(w, html)
}

// truncateString shortens s to at most n runes, appending "…" if truncated.
func truncateString(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}

// formatFileSize converts bytes to human-readable format
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// adminDeleteTrackHandler deletes a track and all associated data.
// POST /api/admin/delete
// Body: {"password": "xxx", "trackID": "abc123"} (password optional if session-authenticated)
//
// @Summary     Admin delete single track
// @Description Deletes a track and associated marker data.
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param       body body object true "Delete request payload"
// @Success     200 {object} map[string]interface{} "Delete result"
// @Failure     400 {string} string "Invalid request"
// @Failure     401 {string} string "Unauthorized"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/delete [post]
func adminDeleteTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if db == nil || db.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	// Parse request
	var req struct {
		Password string `json:"password"`
		TrackID  string `json:"trackID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check for admin access: session-based or password-based
	isSessionAdmin := false
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		isSessionAdmin = true
	}

	if !isSessionAdmin {
		// Fall back to password authentication
		if *adminPassword == "" {
			http.Error(w, "Admin endpoints are disabled - please login as admin", http.StatusForbidden)
			return
		}
		if req.Password != *adminPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Validate trackID
	if req.TrackID == "" {
		http.Error(w, "trackID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := db.DeleteTrack(ctx, req.TrackID)
	if err != nil {
		log.Printf("Error deleting track %s: %v", req.TrackID, err)
		http.Error(w, "Failed to delete track", http.StatusInternalServerError)
		return
	}

	log.Printf("Admin deleted track: %s", req.TrackID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Track deleted successfully",
		"trackID": req.TrackID,
	})
}

// adminDeleteMultipleTracksHandler deletes multiple tracks and all associated data.
// POST /api/admin/delete-multiple
// Body: {"password": "xxx", "trackIDs": ["abc123", "def456"]} (password optional if session-authenticated)
//
// @Summary     Admin delete multiple tracks
// @Description Deletes multiple tracks in one request.
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param       body body object true "Bulk delete request payload"
// @Success     200 {object} map[string]interface{} "Delete result"
// @Failure     400 {string} string "Invalid request"
// @Failure     401 {string} string "Unauthorized"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/delete-multiple [post]
func adminDeleteMultipleTracksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if db == nil || db.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	// Parse request
	var req struct {
		Password string   `json:"password"`
		TrackIDs []string `json:"trackIDs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check for admin access: session-based or password-based
	isSessionAdmin := false
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		isSessionAdmin = true
	}

	if !isSessionAdmin {
		// Fall back to password authentication
		if *adminPassword == "" {
			http.Error(w, "Admin endpoints are disabled - please login as admin", http.StatusForbidden)
			return
		}
		if req.Password != *adminPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Validate trackIDs
	if len(req.TrackIDs) == 0 {
		http.Error(w, "trackIDs is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	deleted := 0
	var errors []string

	for _, trackID := range req.TrackIDs {
		err := db.DeleteTrack(ctx, trackID)
		if err != nil {
			log.Printf("Error deleting track %s: %v", trackID, err)
			errors = append(errors, fmt.Sprintf("%s: %v", trackID, err))
		} else {
			deleted++
		}
	}

	log.Printf("Admin deleted %d tracks (attempted %d)", deleted, len(req.TrackIDs))

	w.Header().Set("Content-Type", "application/json")
	if len(errors) > 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "partial",
			"message": fmt.Sprintf("Deleted %d of %d tracks", deleted, len(req.TrackIDs)),
			"deleted": deleted,
			"errors":  errors,
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "All tracks deleted successfully",
			"deleted": deleted,
		})
	}
}

// adminUpdateTrackHandler updates the editable metadata (name, username, notes) for a track.
// PUT /api/admin/tracks/update
// Body: {"password":"xxx","track_id":"abc","name":"...","username":"...","notes":"..."}
func adminUpdateTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if db == nil || db.DB == nil {
		http.Error(w, "database not available", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		Password string `json:"password"`
		TrackID  string `json:"track_id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Notes    string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TrackID == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	isSessionAdmin := false
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		isSessionAdmin = true
	}
	if !isSessionAdmin {
		if *adminPassword == "" {
			http.Error(w, "admin disabled", http.StatusForbidden)
			return
		}
		if req.Password != *adminPassword {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	_, err := db.DB.ExecContext(r.Context(),
		`UPDATE uploads SET name = $1, username = $2, notes = $3 WHERE track_id = $4`,
		req.Name, req.Username, req.Notes, req.TrackID,
	)
	if err != nil {
		log.Printf("adminUpdateTrack: %v", err)
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}

// adminUpdateUploadHandler updates editable metadata for a single upload record.
// PUT /api/admin/uploads/update
// Body: {"password":"xxx","upload_id":123,"filename":"...","username":"...","detector":"...","recording_date":"YYYY-MM-DD","notes":"..."}
func adminUpdateUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if db == nil || db.DB == nil {
		http.Error(w, "database not available", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		Password      string `json:"password"`
		UploadID      int64  `json:"upload_id"`
		Filename      string `json:"filename"`
		Username      string `json:"username"`
		Detector      string `json:"detector"`
		RecordingDate string `json:"recording_date"` // "YYYY-MM-DD" or ""
		Notes         string `json:"notes"`
		Comment       string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UploadID == 0 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	isSessionAdmin := false
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		isSessionAdmin = true
	}
	if !isSessionAdmin {
		if *adminPassword == "" {
			http.Error(w, "admin disabled", http.StatusForbidden)
			return
		}
		if req.Password != *adminPassword {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	var err error
	if req.RecordingDate != "" {
		t, parseErr := time.Parse("2006-01-02", req.RecordingDate)
		if parseErr != nil {
			http.Error(w, "invalid recording_date format", http.StatusBadRequest)
			return
		}
		_, err = db.DB.ExecContext(r.Context(),
			`UPDATE uploads SET filename = $1, username = $2, detector = $3, notes = $4, recording_date = $5, comment = $6 WHERE id = $7`,
			req.Filename, req.Username, req.Detector, req.Notes, t, req.Comment, req.UploadID,
		)
	} else {
		_, err = db.DB.ExecContext(r.Context(),
			`UPDATE uploads SET filename = $1, username = $2, detector = $3, notes = $4, comment = $5 WHERE id = $6`,
			req.Filename, req.Username, req.Detector, req.Notes, req.Comment, req.UploadID,
		)
	}
	if err != nil {
		log.Printf("adminUpdateUpload: %v", err)
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}

// adminImportSafecastMetadataHandler fetches name/metadata from the old Safecast API
// for all uploads with source='safecast-api' and a non-empty source_id, then saves
// the name field back to uploads.name.
// POST /api/admin/tracks/import-safecast
// Body: {"password":"xxx"}
func adminImportSafecastMetadataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if db == nil || db.DB == nil {
		http.Error(w, "database not available", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	isSessionAdmin := false
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		isSessionAdmin = true
	}
	if !isSessionAdmin {
		if *adminPassword == "" {
			http.Error(w, "admin disabled", http.StatusForbidden)
			return
		}
		if req.Password != *adminPassword {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Fetch all uploads from the Safecast API source that have a source_id
	rows, err := db.DB.QueryContext(r.Context(),
		`SELECT track_id, source_id FROM uploads WHERE source = 'safecast-api' AND source_id IS NOT NULL AND source_id != '' AND (name IS NULL OR name = '' OR name = filename OR comment IS NULL OR comment = '')`,
	)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type candidate struct {
		TrackID  string
		SourceID string
	}
	var candidates []candidate
	for rows.Next() {
		var c candidate
		if err := rows.Scan(&c.TrackID, &c.SourceID); err == nil {
			candidates = append(candidates, c)
		}
	}
	rows.Close()

	const numWorkers = 16
	type result struct{ updated bool }
	jobs := make(chan candidate, len(candidates))
	results := make(chan result, len(candidates))

	httpClient := &http.Client{Timeout: 10 * time.Second}
	for i := 0; i < numWorkers; i++ {
		go func() {
			for c := range jobs {
				apiURL := "https://api.safecast.org/bgeigie_imports/" + c.SourceID + ".json"
				resp, err := httpClient.Get(apiURL)
				if err != nil || resp.StatusCode != http.StatusOK {
					if resp != nil {
						resp.Body.Close()
					}
					results <- result{}
					continue
				}
				var meta struct {
					Name    string `json:"name"`
					Comment string `json:"comment"`
				}
				decodeErr := json.NewDecoder(resp.Body).Decode(&meta)
				resp.Body.Close()
				if decodeErr != nil || meta.Name == "" {
					results <- result{}
					continue
				}
				_, err = db.DB.ExecContext(r.Context(),
					`UPDATE uploads SET name = $1, comment = $2 WHERE track_id = $3`,
					meta.Name, meta.Comment, c.TrackID,
				)
				results <- result{updated: err == nil}
			}
		}()
	}

	for _, c := range candidates {
		jobs <- c
	}
	close(jobs)

	updated := 0
	for range candidates {
		if r := <-results; r.updated {
			updated++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"ok": true, "updated": updated, "total": len(candidates)})
}

// adminImportFromSafecastHandler manually imports files from Safecast API for a date range.
// POST /api/admin/import-from-safecast?password=xxx
// Body: {"start_date": "2025-01-01", "end_date": "2025-01-31"}
// Streams progress updates via Server-Sent Events
//
// @Summary     Admin import from Safecast API
// @Description Triggers a date-ranged Safecast import and streams progress updates.
// @Tags        admin
// @Accept      json
// @Produce     text/event-stream
// @Success     200 {string} string "Streaming progress"
// @Failure     401 {string} string "Unauthorized"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/import-from-safecast [post]
func adminImportFromSafecastHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for admin access: session-based or password-based
	isSessionAdmin := false
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		isSessionAdmin = true
	}

	if !isSessionAdmin {
		// Fall back to password authentication
		if *adminPassword == "" {
			http.Error(w, "Admin endpoints are disabled - please login as admin", http.StatusForbidden)
			return
		}
		password := r.URL.Query().Get("password")
		if password != *adminPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	if db == nil || db.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	// Parse request
	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate dates
	if req.StartDate == "" || req.EndDate == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	// Set up Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Helper function to send SSE progress updates
	sendProgress := func(message string, imported, skipped, errors, total int) {
		data := map[string]interface{}{
			"message":  message,
			"imported": imported,
			"skipped":  skipped,
			"errors":   errors,
			"total":    total,
		}
		jsonData, _ := json.Marshal(data)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
	}

	// Import files from Safecast API
	ctx := r.Context()
	client := safecastfetcher.NewClient()

	// Parse date range once before the loop
	startTime, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		sendProgress("Error: Invalid start_date format (use YYYY-MM-DD)", 0, 0, 1, 0)
		return
	}
	endTime, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		sendProgress("Error: Invalid end_date format (use YYYY-MM-DD)", 0, 0, 1, 0)
		return
	}
	endTime = endTime.Add(24 * time.Hour) // Include end date

	var allImports []safecastfetcher.SafecastImport
	imported := 0
	skipped := 0
	errors := 0

	// Send initial status
	sendProgress("Fetching imports from Safecast API...", 0, 0, 0, 0)

	// Fetch imports page by page
	// Note: API returns ~25-50 items per page. With 2042+ pages total, we need high limit.
	// Loop will stop early once we pass the start date, so this is just a safety ceiling.
	for page := 1; page <= 3000; page++ { // Safety limit (enough for all historical imports)
		imports, err := client.FetchApprovedImports(ctx, req.StartDate, page, false) // oldest first for batch import
		if err != nil {
			log.Printf("[admin-import] Error fetching page %d: %v", page, err)
			break
		}

		if len(imports) == 0 {
			break
		}

		// Filter by date range
		for _, imp := range imports {
			// Skip imports after end date (too new)
			if imp.CreatedAt.After(endTime) {
				continue
			}
			// Skip imports before start date (too old)
			if imp.CreatedAt.Before(startTime) {
				continue
			}
			// Within range - add it
			allImports = append(allImports, imp)
		}

		// Update progress
		sendProgress(fmt.Sprintf("Fetched page %d, found %d imports so far...", page, len(allImports)), 0, 0, 0, len(allImports))

		// Stop pagination if we've gone past the start date
		// API returns results in DESC order (newest first), so if the oldest import
		// on this page is before the start date, we've collected all imports in range
		if len(imports) > 0 && imports[len(imports)-1].CreatedAt.Before(startTime) {
			log.Printf("[admin-import] Reached imports before start date, stopping pagination")
			break
		}
	}

	log.Printf("[admin-import] Found %d imports in date range", len(allImports))
	sendProgress(fmt.Sprintf("Found %d imports in date range. Starting parallel import...", len(allImports)), 0, 0, 0, len(allImports))

	// Parallel import with worker pool
	numWorkers := 8 // Process 8 files concurrently
	if len(allImports) < numWorkers {
		numWorkers = len(allImports)
	}

	// Create channels for work distribution
	type importJob struct {
		index int
		imp   safecastfetcher.SafecastImport
	}
	jobs := make(chan importJob, len(allImports))

	// Mutex for protecting shared counters and progress updates
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Track number of processed imports for progress
	processed := 0

	// Start worker goroutines
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for job := range jobs {
				imp := job.imp

				// Check if already imported
				exists, err := db.CheckImportExists(ctx, safecastfetcher.SourceTypeSafecastAPI, imp.ID)
				if err != nil {
					log.Printf("[admin-import] worker %d: Error checking import #%d: %v", workerID, imp.ID, err)
					mu.Lock()
					errors++
					processed++
					sendProgress(fmt.Sprintf("Processing %d/%d: Error checking #%d", processed, len(allImports), imp.ID), imported, skipped, errors, len(allImports))
					mu.Unlock()
					continue
				}

				if exists {
					mu.Lock()
					skipped++
					processed++
					sendProgress(fmt.Sprintf("Processing %d/%d: Skipped #%d (already imported)", processed, len(allImports), imp.ID), imported, skipped, errors, len(allImports))
					mu.Unlock()
					continue
				}

				// Download and import
				mu.Lock()
				sendProgress(fmt.Sprintf("Processing %d/%d: Worker %d downloading #%d (%s)...", processed+1, len(allImports), workerID, imp.ID, imp.Name), imported, skipped, errors, len(allImports))
				mu.Unlock()

				content, filename, err := safecastfetcher.DownloadLogFile(ctx, imp.SourceURL)
				if err != nil {
					log.Printf("[admin-import] worker %d: import #%d: download failed: %v", workerID, imp.ID, err)
					mu.Lock()
					errors++
					processed++
					sendProgress(fmt.Sprintf("Processing %d/%d: Error downloading #%d", processed, len(allImports), imp.ID), imported, skipped, errors, len(allImports))
					mu.Unlock()
					continue
				}

				// Fetch username from API
				username := ""
				if imp.UserID > 0 {
					user, err := client.FetchUser(ctx, imp.UserID)
					if err != nil {
						log.Printf("[admin-import] worker %d: import #%d: warning: failed to fetch username for user %d: %v", workerID, imp.ID, imp.UserID, err)
					} else if user != nil {
						username = user.Name
					}
				}

				// Import using the importer function from the fetcher
				mu.Lock()
				sendProgress(fmt.Sprintf("Processing %d/%d: Worker %d importing #%d...", processed+1, len(allImports), workerID, imp.ID), imported, skipped, errors, len(allImports))
				mu.Unlock()

				// Use the efficient batch processing importer from safecast-fetcher package
				// This uses the same optimized code path as backfill command
				result, err := safecastfetcher.ImportSafecastFile(
					ctx,
					content,
					filename,
					int64(imp.ID),
					imp.SourceURL,
					fmt.Sprintf("%d", imp.UserID),
					username,
					imp.Comment,
					db,
					*dbType,
					nil, // Use default importer
				)

				if err != nil {
					log.Printf("[admin-import] worker %d: import #%d: import failed: %v", workerID, imp.ID, err)
					mu.Lock()
					errors++
					processed++
					sendProgress(fmt.Sprintf("Processing %d/%d: Error importing #%d", processed, len(allImports), imp.ID), imported, skipped, errors, len(allImports))
					mu.Unlock()
					continue
				}

				mu.Lock()
				imported++
				processed++
				log.Printf("[admin-import] worker %d: import #%d: success (track %s, %d markers)", workerID, imp.ID, result.TrackID, result.MarkerCount)
				sendProgress(fmt.Sprintf("Processing %d/%d: Worker %d imported #%d ✓ (%d markers)", processed, len(allImports), workerID, imp.ID, result.MarkerCount), imported, skipped, errors, len(allImports))
				mu.Unlock()
			}
		}(w)
	}

	// Send all jobs to workers
	for i, imp := range allImports {
		jobs <- importJob{index: i, imp: imp}
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()

	log.Printf("[admin-import] Complete: imported=%d skipped=%d errors=%d", imported, skipped, errors)

	// Send final status
	finalMsg := fmt.Sprintf("Complete! Imported %d files, skipped %d, errors %d", imported, skipped, errors)
	sendProgress(finalMsg, imported, skipped, errors, len(allImports))

	// Send done event
	fmt.Fprintf(w, "event: done\ndata: {\"status\": \"success\", \"imported\": %d, \"skipped\": %d, \"errors\": %d, \"total\": %d}\n\n", imported, skipped, errors, len(allImports))
	flusher.Flush()
}

// adminImportByIDHandler imports a single track from the Safecast API by its bgeigie import ID.
// POST /api/admin/import-by-id?password=xxx
// Body: {"id": 32558}
func adminImportByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	isSessionAdmin := false
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		isSessionAdmin = true
	}
	if !isSessionAdmin {
		if *adminPassword == "" {
			http.Error(w, "Admin endpoints are disabled - please login as admin", http.StatusForbidden)
			return
		}
		if r.URL.Query().Get("password") != *adminPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	if db == nil || db.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.ID <= 0 {
		http.Error(w, "id must be a positive integer", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	client := safecastfetcher.NewClient()

	imp, err := client.FetchImportByID(ctx, req.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("fetch import: %v", err), http.StatusBadGateway)
		return
	}
	if imp == nil {
		http.Error(w, fmt.Sprintf("import #%d not found", req.ID), http.StatusNotFound)
		return
	}

	// Check for duplicate
	exists, err := db.CheckImportExists(ctx, safecastfetcher.SourceTypeSafecastAPI, imp.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("check duplicate: %v", err), http.StatusInternalServerError)
		return
	}
	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"skipped": true,
			"message": fmt.Sprintf("import #%d already exists", imp.ID),
		})
		return
	}

	content, filename, err := safecastfetcher.DownloadLogFile(ctx, imp.SourceURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("download: %v", err), http.StatusBadGateway)
		return
	}

	username := ""
	if imp.UserID > 0 {
		user, err := client.FetchUser(ctx, imp.UserID)
		if err != nil {
			log.Printf("[admin-import-by-id] warning: failed to fetch username for user %d: %v", imp.UserID, err)
		} else if user != nil {
			username = user.Name
		}
	}

	result, err := safecastfetcher.ImportSafecastFile(
		ctx,
		content,
		filename,
		imp.ID,
		imp.SourceURL,
		fmt.Sprintf("%d", imp.UserID),
		username,
		imp.Comment,
		db,
		*dbType,
		nil,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("import failed: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("[admin-import-by-id] import #%d: success (track %s, %d markers)", imp.ID, result.TrackID, result.MarkerCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"imported":     true,
		"track_id":     result.TrackID,
		"marker_count": result.MarkerCount,
		"filename":     result.Filename,
		"comment":      imp.Comment,
	})
}

// adminTracksHandler lists all tracks in the system with statistics.
// GET /api/admin/tracks?password=xxx&limit=1000
//
// @Summary     Admin tracks dashboard data
// @Description Returns paginated track statistics for admin users.
// @Tags        admin
// @Produce     html
// @Param       limit query int false "Page size"
// @Param       page query int false "Page number"
// @Param       search query string false "Free-text search"
// @Param       detector query string false "Detector filter"
// @Success     200 {string} string "HTML admin page"
// @Failure     401 {string} string "Unauthorized"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/tracks [get]
func adminTracksHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent CloudFront from caching this dynamic admin page
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	authorized, password := checkAdminAuth(w, r)
	if !authorized {
		return
	}
	_ = password // Used in HTML template generation below

	if db == nil || db.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	// Get limit parameter (page size)
	limit := 500 // Default to 500 per page
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Get page parameter
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get search parameter
	search := r.URL.Query().Get("search")

	// Get detector filter parameter
	detectorFilter := r.URL.Query().Get("detector")

	ctx := r.Context()

	// Build WHERE conditions for search
	var whereConditions []string
	var countArgs []interface{}
	paramCount := 0

	// Always exclude live tracks (use ts. prefix since we now have a JOIN)
	paramCount++
	if *dbType == "pgx" {
		whereConditions = append(whereConditions, fmt.Sprintf("ts.trackID NOT LIKE $%d", paramCount))
	} else {
		whereConditions = append(whereConditions, "ts.trackID NOT LIKE ?")
	}
	countArgs = append(countArgs, "live:%")

	// Add detector filter condition
	if detectorFilter != "" {
		paramCount++
		if *dbType == "pgx" {
			whereConditions = append(whereConditions, fmt.Sprintf("ts.detector ILIKE $%d", paramCount))
			countArgs = append(countArgs, "%"+detectorFilter+"%")
		} else {
			whereConditions = append(whereConditions, "ts.detector LIKE ?")
			countArgs = append(countArgs, "%"+detectorFilter+"%")
		}
	}

	// Add search condition (track ID, detector, name, username)
	if search != "" {
		paramCount++
		if *dbType == "pgx" {
			whereConditions = append(whereConditions, fmt.Sprintf(
				"(ts.trackID ILIKE $%d OR COALESCE(ts.detector,'') ILIKE $%d)",
				paramCount, paramCount))
			countArgs = append(countArgs, "%"+search+"%")
		} else {
			whereConditions = append(whereConditions, "(ts.trackID LIKE ? OR COALESCE(ts.detector,'') LIKE ?)")
			searchPattern := "%" + search + "%"
			for i := 0; i < 2; i++ {
				countArgs = append(countArgs, searchPattern)
			}
		}
	}

	whereClause := "WHERE " + whereConditions[0]
	for i := 1; i < len(whereConditions); i++ {
		whereClause += " AND " + whereConditions[i]
	}

	// Get total count for pagination
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM track_statistics ts " + whereClause
	err := db.DB.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		log.Printf("Error counting tracks: %v", err)
		http.Error(w, "Failed to count tracks", http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := (totalCount + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	// Build args for main query
	args := make([]interface{}, len(countArgs))
	copy(args, countArgs)
	args = append(args, limit, offset)

	// Query tracks using materialized view for performance, joined with uploads for editable metadata.
	// DISTINCT ON ensures one row per track even if multiple upload records exist.
	var query string
	if *dbType == "pgx" {
		paramCount++
		limitParam := paramCount
		paramCount++
		offsetParam := paramCount
		query = `
			SELECT ts.trackID, ts.marker_count, ts.first_date, ts.last_date, ts.spectra_count,
			       COALESCE(ts.detector, '') as detector,
			       COALESCE(u.name, u.filename, '') as name,
			       COALESCE(u.notes, '') as notes,
			       COALESCE(u.username, '') as username,
			       COALESCE(u.source, '') as source,
			       COALESCE(u.source_id, '') as source_id
			FROM track_statistics ts
			LEFT JOIN LATERAL (
			    SELECT filename, name, notes, username, source, source_id
			    FROM uploads WHERE track_id = ts.trackID LIMIT 1
			) u ON true
			` + whereClause + `
			ORDER BY ts.last_date DESC
			LIMIT $` + strconv.Itoa(limitParam) + ` OFFSET $` + strconv.Itoa(offsetParam)
	} else {
		query = `
			SELECT ts.trackID, ts.marker_count, ts.first_date, ts.last_date, ts.spectra_count,
			       COALESCE(ts.detector, '') as detector,
			       COALESCE(u.name, u.filename, '') as name,
			       COALESCE(u.notes, '') as notes,
			       COALESCE(u.username, '') as username,
			       COALESCE(u.source, '') as source,
			       COALESCE(u.source_id, '') as source_id
			FROM track_statistics ts
			LEFT JOIN uploads u ON u.track_id = ts.trackID
			` + whereClause + `
			ORDER BY ts.last_date DESC
			LIMIT ? OFFSET ?`
	}

	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error fetching tracks: %v", err)
		http.Error(w, "Failed to fetch tracks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type TrackInfo struct {
		TrackID      string `json:"trackID"`
		MarkerCount  int64  `json:"markerCount"`
		FirstDate    int64  `json:"firstDate"`
		LastDate     int64  `json:"lastDate"`
		SpectraCount int64  `json:"spectraCount"`
		Detector     string `json:"detector"`
		Name         string `json:"name"`
		Notes        string `json:"notes"`
		Username     string `json:"username"`
		Source       string `json:"source"`
		SourceID     string `json:"sourceID"`
	}

	var tracks []TrackInfo
	for rows.Next() {
		var t TrackInfo
		if err := rows.Scan(&t.TrackID, &t.MarkerCount, &t.FirstDate, &t.LastDate, &t.SpectraCount, &t.Detector,
			&t.Name, &t.Notes, &t.Username, &t.Source, &t.SourceID); err != nil {
			log.Printf("Error scanning track row: %v", err)
			continue
		}
		tracks = append(tracks, t)
	}

	// Return HTML table
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html>
<head>
	<title>Admin - All Tracks</title>
	<script>
	(function() {
		var saved = localStorage.getItem('safecastDocTheme');
		var prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		document.documentElement.setAttribute('data-theme', saved || (prefersDark ? 'dark' : 'light'));
	})();
	</script>

	<!-- favicon -->
	<link rel="apple-touch-icon" sizes="180x180" href="/static/images/apple-touch-icon.png">
	<link rel="icon" type="image/png" sizes="32x32" href="/static/images/favicon-32x32.png">
	<link rel="icon" type="image/png" sizes="16x16" href="/static/images/favicon-16x16.png">
	<link rel="manifest" href="/static/images/site.webmanifest">

	<style>
		:root {
			--bg-primary: #f5f5f5;
			--bg-card: white;
			--text-primary: #333;
			--text-secondary: #666;
			--text-muted: #999;
			--border-color: #ddd;
			--link-color: #0066cc;
			--shadow: 0 1px 3px rgba(0,0,0,0.1);
			--th-bg: #424242;
			--hover-bg: #f9f9f9;
			--badge-bg: #e3f2fd;
			--badge-text: #1976d2;
			--badge-spectrum-bg: #e8f5e9;
			--badge-spectrum-text: #388e3c;
		}
		@media (prefers-color-scheme: dark) {
			:root {
				--bg-primary: #1a1a1a;
				--bg-card: #2b2b2b;
				--text-primary: #eee;
				--text-secondary: #aaa;
				--text-muted: #777;
				--border-color: #444;
				--link-color: #90caf9;
				--shadow: 0 1px 3px rgba(255,255,255,0.1);
				--th-bg: #616161;
				--hover-bg: #333;
				--badge-bg: rgba(144, 202, 249, 0.2);
				--badge-text: #90caf9;
				--badge-spectrum-bg: rgba(76, 175, 80, 0.2);
				--badge-spectrum-text: #81c784;
				color-scheme: dark;
			}
		}
		:root[data-theme='light'] {
			--bg-primary: #f5f5f5;
			--bg-card: white;
			--text-primary: #333;
			--text-secondary: #666;
			--text-muted: #999;
			--border-color: #ddd;
			--link-color: #0066cc;
			--shadow: 0 1px 3px rgba(0,0,0,0.1);
			--th-bg: #424242;
			--hover-bg: #f9f9f9;
			--badge-bg: #e3f2fd;
			--badge-text: #1976d2;
			--badge-spectrum-bg: #e8f5e9;
			--badge-spectrum-text: #388e3c;
			color-scheme: light;
		}
		:root[data-theme='dark'] {
			--bg-primary: #1a1a1a;
			--bg-card: #2b2b2b;
			--text-primary: #eee;
			--text-secondary: #aaa;
			--text-muted: #777;
			--border-color: #444;
			--link-color: #90caf9;
			--shadow: 0 1px 3px rgba(255,255,255,0.1);
			--th-bg: #616161;
			--hover-bg: #333;
			--badge-bg: rgba(144, 202, 249, 0.2);
			--badge-text: #90caf9;
			--badge-spectrum-bg: rgba(76, 175, 80, 0.2);
			--badge-spectrum-text: #81c784;
			color-scheme: dark;
		}
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Arial, sans-serif; margin: 0; background: var(--bg-primary); color: var(--text-primary); transition: background 0.2s, color 0.2s; }
		h1 { color: var(--text-primary); }
		.top-nav { background: #1a3a5c; color: #fff; padding: 0 24px; display: flex; align-items: center; gap: 16px; height: 52px; box-shadow: 0 2px 6px rgba(0,0,0,0.3); position: sticky; top: 0; z-index: 100; }
		.nav-logo { display: flex; align-items: center; gap: 8px; text-decoration: none; color: #fff; font-weight: 700; font-size: 17px; white-space: nowrap; }
		.nav-logo img { height: 28px; width: 28px; object-fit: contain; }
		.nav-sep { color: rgba(255,255,255,0.3); font-size: 18px; }
		.nav-title { font-size: 15px; font-weight: 600; color: #d0e8ff; white-space: nowrap; }
		.nav-spacer { flex: 1; }
		.back-link { color: #afd4f5; text-decoration: none; font-size: 13px; white-space: nowrap; }
		.back-link:hover { color: #fff; }
		#theme-toggle { background: rgba(255,255,255,0.12); color: #fff; border: 1px solid rgba(255,255,255,0.25); border-radius: 6px; padding: 6px 14px; font-size: 13px; font-weight: 600; cursor: pointer; white-space: nowrap; transition: background 0.2s; }
		#theme-toggle:hover { background: rgba(255,255,255,0.22); }
		.page-content { padding: 20px 24px; }
		.nav { background: var(--bg-card); padding: 15px; margin-bottom: 20px; border-radius: 5px; box-shadow: var(--shadow); display: flex; align-items: center; justify-content: space-between; }
		.nav-left { display: flex; align-items: center; gap: 15px; }
		.nav a { color: var(--link-color); text-decoration: none; }
		.nav a:hover { text-decoration: underline; }
		.summary { background: var(--bg-card); padding: 15px; margin-bottom: 20px; border-radius: 5px; box-shadow: var(--shadow); }
		table { border-collapse: collapse; width: 100%; background: var(--bg-card); box-shadow: var(--shadow); }
		th { background: var(--th-bg); color: white; padding: 12px; text-align: left; font-weight: 600; }
		td { padding: 6px 8px; border-bottom: 1px solid var(--border-color); }
		tr:hover { background: var(--hover-bg); }
		.empty { text-align: center; padding: 40px; color: var(--text-muted); font-style: italic; }
		.trackid { font-family: monospace; color: var(--link-color); }
		.detector { font-family: monospace; font-size: 0.9em; }
		.detector-link { color: var(--badge-text); background: var(--badge-bg); padding: 2px 8px; border-radius: 10px; text-decoration: none; }
		.detector-link:hover { text-decoration: underline; }
		.badge { display: inline-block; padding: 3px 8px; border-radius: 12px; font-size: 0.85em; background: var(--badge-bg); color: var(--badge-text); }
		.badge.spectrum { background: var(--badge-spectrum-bg); color: var(--badge-spectrum-text); }
		.datetime { color: var(--text-secondary); font-size: 0.9em; }
		.delete-btn { background: #f44336; color: white; border: none; padding: 5px 10px; border-radius: 3px; cursor: pointer; font-size: 0.85em; }
		.delete-btn:hover { background: #d32f2f; }
		.edit-btn { background: #2196F3; color: white; border: none; padding: 5px 10px; border-radius: 3px; cursor: pointer; font-size: 0.85em; margin-right: 4px; }
		.edit-btn:hover { background: #1976D2; }
		.import-btn { background: #4caf50; color: white; border: none; padding: 8px 16px; border-radius: 3px; cursor: pointer; font-size: 0.9em; }
		.import-btn:hover { background: #388e3c; }
		.import-btn:disabled { background: #888; cursor: default; }
		.modal-overlay { display:none; position:fixed; top:0; left:0; width:100%; height:100%; background:rgba(0,0,0,0.5); z-index:1000; align-items:center; justify-content:center; }
		.modal-overlay.open { display:flex; }
		.modal { background:var(--bg-card); border-radius:8px; padding:24px; width:480px; max-width:95vw; box-shadow:0 4px 24px rgba(0,0,0,0.3); }
		.modal h3 { margin:0 0 16px; color:var(--text-primary); }
		.modal label { display:block; margin-bottom:4px; font-size:0.85em; color:var(--text-secondary); }
		.modal input, .modal textarea { width:100%; padding:8px; border:1px solid var(--border-color); border-radius:4px; background:var(--bg-primary); color:var(--text-primary); font-size:0.9em; margin-bottom:12px; box-sizing:border-box; }
		.modal textarea { min-height:80px; resize:vertical; }
		.modal-actions { display:flex; gap:8px; justify-content:flex-end; margin-top:8px; }
		.modal-save { background:#2196F3; color:white; border:none; padding:8px 20px; border-radius:4px; cursor:pointer; }
		.modal-save:hover { background:#1976D2; }
		.modal-cancel { background:var(--border-color); color:var(--text-primary); border:none; padding:8px 16px; border-radius:4px; cursor:pointer; }
		#importStatus { margin-top:8px; font-size:0.85em; color:var(--text-secondary); }
		.delete-selected-btn { background: #f44336; color: white; border: none; padding: 10px 20px; border-radius: 3px; cursor: pointer; font-size: 1em; margin-left: 10px; }
		.delete-selected-btn:hover { background: #d32f2f; }
		.delete-selected-btn:disabled { background: #ccc; cursor: not-allowed; }
		.backfill-btn { background: #4CAF50; color: white; border: none; padding: 10px 20px; border-radius: 3px; cursor: pointer; font-size: 1em; margin-left: 10px; }
		.backfill-btn:hover { background: #45a049; }
		.view-selected-btn { background: #2196F3; color: white; border: none; padding: 10px 20px; border-radius: 3px; cursor: pointer; font-size: 1em; margin-left: 10px; }
		.view-selected-btn:hover { background: #1976D2; }
		.view-selected-btn:disabled { background: #ccc; cursor: not-allowed; }
		.checkbox-col { width: 40px; text-align: center; }
		.sortable { cursor: pointer; user-select: none; position: relative; padding-right: 20px; }
		.sortable:hover { background: rgba(255,255,255,0.1); }
		.sortable::after { content: '⇅'; position: absolute; right: 8px; opacity: 0.5; }
		.sortable.asc::after { content: '▲'; opacity: 1; }
		.sortable.desc::after { content: '▼'; opacity: 1; }
		.filter-input { width: 100%; padding: 4px 8px; border: 1px solid var(--border-color); border-radius: 3px; background: var(--bg-card); color: var(--text-primary); font-size: 0.85em; box-sizing: border-box; }
		.filter-row th { background: var(--bg-card); padding: 8px 12px; }
		/* Admin tab bar */
		.admin-tabs { display: flex; gap: 2px; margin-top: 10px; margin-bottom: 10px; background: var(--border-color); border-radius: 8px; overflow: hidden; }
		.admin-tabs a, .admin-tabs span { padding: 10px 20px; text-decoration: none; color: var(--text-secondary); background: var(--bg-card); font-weight: 500; font-size: 0.95em; transition: background 0.2s; }
		.admin-tabs a:hover { background: var(--hover-bg); color: var(--text-primary); }
		.admin-tabs a.active { background: #2196F3; color: white; }
		.admin-tabs span.disabled { color: var(--text-muted); cursor: not-allowed; font-style: italic; }
	</style>
</head>
<body>
<nav class="top-nav">
  <a href="/" class="nav-logo">
    <img src="/static/images/safecast-logo-squared.png" alt="Safecast">
    Safecast
  </a>
  <span class="nav-sep">|</span>
  <span class="nav-title">All Tracks Administration</span>
  <span class="nav-spacer"></span>
  <a href="/" class="back-link">&#8592; Back to Map</a>
  <button id="theme-toggle" onclick="toggleTheme()">&#127769; Dark Mode</button>
</nav>
<div class="page-content">
	<div class="admin-tabs">
		<a href="/admin/users` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">Users</a>
		<a href="/admin/uploads` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">Uploads</a>
		<a href="/admin/mcp` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">MCP Analytics</a>
		<a href="/admin/realtime` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">Realtime</a>
		<a href="/admin/translations` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">Translations</a>
		<a href="/admin/tour` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">Tour Steps</a>
		<a href="/admin/ai-hints` + func() string {
		if password != "" {
			return "?password=" + password
		}
		return ""
	}() + `">AI Hints</a>
	</div>
	<div class="nav">
		<div class="nav-left">
			<button class="backfill-btn" onclick="backfillUploads()">Backfill Upload Records</button>
			<button class="import-btn" onclick="importSafecastMeta()" id="importSafecastBtn">Import Safecast API Metadata</button>
			<button class="view-selected-btn" id="viewSelectedBtn" onclick="viewSelected()" disabled>View Selected on Map</button>
			<button class="delete-selected-btn" id="deleteSelectedBtn" onclick="deleteSelected()" disabled>Delete Selected</button>
		</div>
	</div>
	<div id="importStatus"></div>
	<!-- Edit Track Modal -->
	<div class="modal-overlay" id="editModal">
		<div class="modal">
			<h3>Edit Track Metadata</h3>
			<input type="hidden" id="editTrackID">
			<label>Name</label>
			<input type="text" id="editName" placeholder="Display name">
			<label>Uploader / Username</label>
			<input type="text" id="editUsername" placeholder="Username">
			<label>Admin Notes</label>
			<textarea id="editNotes" placeholder="Internal notes (not shown to public)"></textarea>
			<div class="modal-actions">
				<button class="modal-cancel" onclick="closeEdit()">Cancel</button>
				<button class="modal-save" onclick="saveEdit()">Save</button>
			</div>
		</div>
	</div>
	<div class="summary">
		<strong>Total Tracks:</strong> ` + strconv.Itoa(totalCount) + ` tracks
		<span style="margin-left: 20px;">
			<strong>Page ` + strconv.Itoa(page) + ` of ` + strconv.Itoa(totalPages) + `</strong>
			(showing ` + strconv.Itoa(len(tracks)) + ` tracks)
		</span>
		<span style="margin-left: 20px;">
			<label for="limitSelect"><strong>Per page:</strong></label>
			<select id="limitSelect" onchange="changeLimit()" style="margin-left: 5px; padding: 4px 8px; border-radius: 4px; border: 1px solid var(--border-color); background: var(--bg-card); color: var(--text-primary);">
				<option value="100"` + func() string {
		if limit == 100 {
			return " selected"
		}
		return ""
	}() + `>100</option>
				<option value="250"` + func() string {
		if limit == 250 {
			return " selected"
		}
		return ""
	}() + `>250</option>
				<option value="500"` + func() string {
		if limit == 500 {
			return " selected"
		}
		return ""
	}() + `>500</option>
				<option value="1000"` + func() string {
		if limit == 1000 {
			return " selected"
		}
		return ""
	}() + `>1000</option>
			</select>
		</span>
		<span style="margin-left: 20px;">
			<label for="searchInput"><strong>Search:</strong></label>
			<input type="text" id="searchInput" value="` + search + `" placeholder="Search all fields..." autocomplete="off" style="margin-left: 5px; padding: 4px 8px; border-radius: 4px; border: 1px solid var(--border-color); background: var(--bg-card); color: var(--text-primary); width: 200px;" onkeypress="if(event.key === 'Enter') performSearch()">
			<button onclick="performSearch()" style="margin-left: 5px; padding: 4px 12px; border-radius: 4px; border: 1px solid var(--border-color); background: var(--link-color); color: white; cursor: pointer;">🔍</button>
			` + func() string {
		if search != "" {
			return `<button onclick="clearSearch()" style="margin-left: 5px; padding: 4px 12px; border-radius: 4px; border: 1px solid var(--border-color); background: var(--bg-card); color: var(--text-primary); cursor: pointer;">Clear</button>`
		}
		return ""
	}() + `
		</span>`

	if search != "" {
		html += ` | <strong>Search:</strong> "` + search + `"`
	}
	if detectorFilter != "" {
		clearURL := "/api/admin/tracks"
		if password != "" {
			clearURL += "?password=" + password
		}
		html += ` | <strong>Detector:</strong> "` + detectorFilter + `" <a href="` + clearURL + `" style="color: var(--link-color);">(clear)</a>`
	}

	// Add pagination controls inline in the summary
	html += `<div style="margin-top: 10px;">`

	// Helper function to build query parameters
	buildURL := func(pageNum int) string {
		params := []string{}
		if password != "" {
			params = append(params, "password="+password)
		}
		params = append(params, "page="+strconv.Itoa(pageNum))
		params = append(params, "limit="+strconv.Itoa(limit))
		if search != "" {
			params = append(params, "search="+url.QueryEscape(search))
		}
		if detectorFilter != "" {
			params = append(params, "detector="+url.QueryEscape(detectorFilter))
		}
		return "?" + strings.Join(params, "&")
	}

	// Previous button
	if page > 1 {
		html += `<a href="` + buildURL(page-1) + `" class="page-btn">&laquo; Previous</a>`
	} else {
		html += `<span class="page-btn disabled">&laquo; Previous</span>`
	}

	// Page numbers
	startPage := page - 2
	if startPage < 1 {
		startPage = 1
	}
	endPage := startPage + 4
	if endPage > totalPages {
		endPage = totalPages
		startPage = endPage - 4
		if startPage < 1 {
			startPage = 1
		}
	}

	// First page
	if startPage > 1 {
		html += `<a href="` + buildURL(1) + `" class="page-btn">1</a>`
		if startPage > 2 {
			html += `<span class="page-btn disabled">...</span>`
		}
	}

	// Page range
	for i := startPage; i <= endPage; i++ {
		if i == page {
			html += `<span class="page-btn active">` + strconv.Itoa(i) + `</span>`
		} else {
			html += `<a href="` + buildURL(i) + `" class="page-btn">` + strconv.Itoa(i) + `</a>`
		}
	}

	// Last page
	if endPage < totalPages {
		if endPage < totalPages-1 {
			html += `<span class="page-btn disabled">...</span>`
		}
		html += `<a href="` + buildURL(totalPages) + `" class="page-btn">` + strconv.Itoa(totalPages) + `</a>`
	}

	// Next button
	if page < totalPages {
		html += `<a href="` + buildURL(page+1) + `" class="page-btn">Next &raquo;</a>`
	} else {
		html += `<span class="page-btn disabled">Next &raquo;</span>`
	}

	html += `</div>
	</div>`

	if len(tracks) == 0 {
		html += `<div class="empty">No tracks found in the database.</div>`
	} else {
		html += `
	<table id="tracksTable">
		<thead>
			<tr>
				<th class="checkbox-col"><input type="checkbox" id="selectAll" onchange="toggleSelectAll(this)"></th>
				<th class="sortable" onclick="sortTable(1)" data-type="text">Track ID</th>
				<th class="sortable" onclick="sortTable(2)" data-type="text">Name</th>
				<th class="sortable" onclick="sortTable(3)" data-type="text">Uploader</th>
				<th class="sortable" onclick="sortTable(4)" data-type="text">Detector</th>
				<th class="sortable" onclick="sortTable(5)" data-type="number">Markers</th>
				<th class="sortable" onclick="sortTable(6)" data-type="number">Spectra</th>
				<th class="sortable" onclick="sortTable(7)" data-type="date">Last Point</th>
				<th>Actions</th>
			</tr>
			<tr class="filter-row">
				<th></th>
				<th><input type="text" class="filter-input" placeholder="Filter Track ID..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Filter Name..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Filter Uploader..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Filter Detector..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Min..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Min..." onkeyup="filterTable()"></th>
				<th><input type="text" class="filter-input" placeholder="Filter date..." onkeyup="filterTable()"></th>
				<th></th>
			</tr>
		</thead>
		<tbody id="tracksTableBody">`

		for _, track := range tracks {
			lastDate := time.Unix(track.LastDate, 0).Format("2006-01-02 15:04")

			// Format detector with link for filtering
			detectorDisplay := "-"
			if track.Detector != "" {
				detectorDisplay = fmt.Sprintf(`<a href="/api/admin/tracks?password=%s&detector=%s" class="detector-link">%s</a>`,
					password, url.QueryEscape(track.Detector), track.Detector)
			}

			nameDisplay := track.Name
			if nameDisplay == "" {
				nameDisplay = "-"
			}
			usernameDisplay := track.Username
			if usernameDisplay == "" {
				usernameDisplay = "-"
			}

			// Encode fields for data attributes (JS edit modal)
			notesEsc := template.HTMLEscapeString(track.Notes)
			nameEsc := template.HTMLEscapeString(track.Name)
			usernameEsc := template.HTMLEscapeString(track.Username)

			html += fmt.Sprintf(`
			<tr>
				<td class="checkbox-col"><input type="checkbox" class="track-checkbox" value="%s" onchange="updateDeleteButton()"></td>
				<td class="trackid"><a href="/trackid/%s">%s</a></td>
				<td>%s</td>
				<td>%s</td>
				<td class="detector">%s</td>
				<td><span class="badge">%d points</span></td>
				<td><span class="badge spectrum">%d spectra</span></td>
				<td class="datetime">%s</td>
				<td>
					<button class="edit-btn" onclick="openEdit('%s','%s','%s','%s')">Edit</button>
					<button class="delete-btn" onclick="deleteTrack('%s')">Delete</button>
				</td>
			</tr>`,
				track.TrackID,
				track.TrackID, track.TrackID,
				nameDisplay,
				usernameDisplay,
				detectorDisplay,
				track.MarkerCount,
				track.SpectraCount,
				lastDate,
				track.TrackID, nameEsc, usernameEsc, notesEsc,
				track.TrackID,
			)
		}

		html += `
		</tbody>
	</table>`
	}

	html += `
	<script>
		// Apply theme from sessionStorage to match map preference
		(function() {
			const media = window.matchMedia('(prefers-color-scheme: dark)');
			const storedTheme = sessionStorage.getItem('themePreference');
			const theme = storedTheme ? storedTheme : (media.matches ? 'dark' : 'light');
			document.documentElement.dataset.theme = theme;
		})();

		function toggleSelectAll(checkbox) {
			const checkboxes = document.querySelectorAll('.track-checkbox');
			checkboxes.forEach(cb => cb.checked = checkbox.checked);
			updateDeleteButton();
		}

		function updateDeleteButton() {
			const checkboxes = document.querySelectorAll('.track-checkbox:checked');
			const btn = document.getElementById('deleteSelectedBtn');
			btn.disabled = checkboxes.length === 0;
			btn.textContent = checkboxes.length > 0 ? '🗑️ Delete Selected (' + checkboxes.length + ')' : '🗑️ Delete Selected';
		}

		function deleteSelected() {
			const checkboxes = document.querySelectorAll('.track-checkbox:checked');
			const trackIDs = Array.from(checkboxes).map(cb => cb.value);

			if (trackIDs.length === 0) return;

			if (!confirm('Delete ' + trackIDs.length + ' track(s) and all associated data?')) {
				return;
			}

			const password = new URLSearchParams(window.location.search).get('password');

			fetch('/api/admin/delete-multiple', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ password: password, trackIDs: trackIDs })
			})
			.then(r => r.json())
			.then(data => {
				if (data.status === 'success') {
					alert('Successfully deleted ' + data.deleted + ' track(s)');
					window.location.reload();
				} else {
					alert('Error: ' + (data.error || 'Unknown error'));
				}
			})
			.catch(err => alert('Error: ' + err));
		}

		function deleteTrack(trackID) {
			if (!confirm('Delete track ' + trackID + ' and all associated data?')) {
				return;
			}

			const password = new URLSearchParams(window.location.search).get('password');

			fetch('/api/admin/delete', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ password: password, trackID: trackID })
			})
			.then(r => r.json())
			.then(data => {
				if (data.status === 'success') {
					alert('Track deleted successfully');
					window.location.reload();
				} else {
					alert('Error: ' + (data.error || 'Unknown error'));
				}
			})
			.catch(err => alert('Error: ' + err));
		}

		// Change limit and reload page, reset to page 1
		function changeLimit() {
			const limit = document.getElementById('limitSelect').value;
			const url = new URL(window.location.href);
			url.searchParams.set('limit', limit);
			url.searchParams.set('page', '1');
			window.location.href = url.toString();
		}

		// Perform search
		function performSearch() {
			const searchValue = document.getElementById('searchInput').value;
			const url = new URL(window.location.href);
			if (searchValue.trim()) {
				url.searchParams.set('search', searchValue.trim());
			} else {
				url.searchParams.delete('search');
			}
			url.searchParams.set('page', '1'); // Reset to page 1 when searching
			window.location.href = url.toString();
		}

		// Clear search
		function clearSearch() {
			const url = new URL(window.location.href);
			url.searchParams.delete('search');
			url.searchParams.set('page', '1'); // Reset to page 1 when clearing
			window.location.href = url.toString();
		}

		function backfillUploads() {
			if (!confirm('Backfill the uploads table with existing spectrum data? This will create upload records for all spectra currently in the database.')) {
				return;
			}

			const password = new URLSearchParams(window.location.search).get('password');

			fetch('/api/admin/backfill?password=' + password, {
				method: 'POST'
			})
			.then(r => r.json())
			.then(data => {
				if (data.status === 'success') {
					alert('Backfill complete: ' + data.count + ' records created');
					window.location.href = '/api/admin/uploads?password=' + password;
				} else {
					alert('Error: ' + (data.error || 'Unknown error'));
				}
			})
			.catch(err => alert('Error: ' + err));
		}

		// ── Edit track metadata ──────────────────────────────────────
		function openEdit(trackID, name, username, notes) {
			document.getElementById('editTrackID').value  = trackID;
			document.getElementById('editName').value     = name;
			document.getElementById('editUsername').value = username;
			document.getElementById('editNotes').value    = notes;
			document.getElementById('editModal').classList.add('open');
		}

		function closeEdit() {
			document.getElementById('editModal').classList.remove('open');
		}

		function saveEdit() {
			const password  = new URLSearchParams(window.location.search).get('password');
			const trackID   = document.getElementById('editTrackID').value;
			const name      = document.getElementById('editName').value.trim();
			const username  = document.getElementById('editUsername').value.trim();
			const notes     = document.getElementById('editNotes').value.trim();

			fetch('/api/admin/tracks/update', {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ password, track_id: trackID, name, username, notes })
			})
			.then(r => r.json())
			.then(data => {
				if (data.ok) {
					closeEdit();
					window.location.reload();
				} else {
					alert('Save failed: ' + (data.error || 'unknown error'));
				}
			})
			.catch(err => alert('Error: ' + err));
		}

		// Close modal on overlay click
		document.getElementById('editModal').addEventListener('click', function(e) {
			if (e.target === this) closeEdit();
		});

		// ── Import Safecast API metadata ─────────────────────────────
		function importSafecastMeta() {
			if (!confirm('Fetch track names and metadata from the old Safecast API for all imported tracks? This may take a moment.')) return;
			const password = new URLSearchParams(window.location.search).get('password');
			const btn = document.getElementById('importSafecastBtn');
			const status = document.getElementById('importStatus');
			btn.disabled = true;
			btn.textContent = 'Importing…';
			status.textContent = 'Fetching metadata from api.safecast.org…';

			fetch('/api/admin/tracks/import-safecast', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ password })
			})
			.then(r => r.json())
			.then(data => {
				btn.disabled = false;
				btn.textContent = 'Import Safecast API Metadata';
				if (data.ok) {
					status.textContent = 'Done: ' + data.updated + ' track(s) updated.';
					if (data.updated > 0) window.location.reload();
				} else {
					status.textContent = 'Error: ' + (data.error || 'unknown');
				}
			})
			.catch(err => {
				btn.disabled = false;
				btn.textContent = 'Import Safecast API Metadata';
				status.textContent = 'Error: ' + err;
			});
		}

		// Sorting functionality
		let sortDirection = {};
		function sortTable(columnIndex) {
			const table = document.getElementById('tracksTable');
			const tbody = document.getElementById('tracksTableBody');
			const rows = Array.from(tbody.querySelectorAll('tr'));
			const header = table.querySelector('thead tr:first-child th:nth-child(' + (columnIndex + 1) + ')');
			const dataType = header.getAttribute('data-type');

			// Toggle sort direction
			const currentDir = sortDirection[columnIndex] || 'none';
			sortDirection[columnIndex] = currentDir === 'asc' ? 'desc' : 'asc';

			// Remove sort classes from all headers
			table.querySelectorAll('.sortable').forEach(h => {
				h.classList.remove('asc', 'desc');
			});

			// Add sort class to current header
			header.classList.add(sortDirection[columnIndex]);

			// Sort rows
			rows.sort((a, b) => {
				let aVal = a.cells[columnIndex].textContent.trim();
				let bVal = b.cells[columnIndex].textContent.trim();

				// Extract numeric values from badges
				if (dataType === 'number') {
					aVal = parseInt(aVal.match(/\d+/) || '0');
					bVal = parseInt(bVal.match(/\d+/) || '0');
					return sortDirection[columnIndex] === 'asc' ? aVal - bVal : bVal - aVal;
				}

				// Date comparison
				if (dataType === 'date') {
					aVal = new Date(aVal).getTime();
					bVal = new Date(bVal).getTime();
					return sortDirection[columnIndex] === 'asc' ? aVal - bVal : bVal - aVal;
				}

				// Text comparison
				if (sortDirection[columnIndex] === 'asc') {
					return aVal.localeCompare(bVal);
				} else {
					return bVal.localeCompare(aVal);
				}
			});

			// Reappend sorted rows
			rows.forEach(row => tbody.appendChild(row));
		}

		// Filtering functionality
		function filterTable() {
			const table = document.getElementById('tracksTable');
			const tbody = document.getElementById('tracksTableBody');
			const filters = table.querySelectorAll('.filter-input');
			const rows = tbody.querySelectorAll('tr');

			rows.forEach(row => {
				let show = true;
				filters.forEach((filter, index) => {
					const filterValue = filter.value.toLowerCase();
					if (filterValue) {
						const cellIndex = index + 1; // +1 because first column is checkbox
						const cell = row.cells[cellIndex];
						if (cell) {
							const cellText = cell.textContent.toLowerCase();
							if (!cellText.includes(filterValue)) {
								show = false;
							}
						}
					}
				});
				row.style.display = show ? '' : 'none';
			});

			// Update delete button after filtering
			updateDeleteButton();
		}
	</script>
</div><!-- .page-content -->
<script>
(function() {
	function applyLabel() {
		var btn = document.getElementById('theme-toggle');
		if (btn) btn.textContent = document.documentElement.getAttribute('data-theme') === 'dark' ? '\u2600\uFE0F Light Mode' : '\uD83C\uDF19 Dark Mode';
	}
	applyLabel();
	window.toggleTheme = function() {
		var next = document.documentElement.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
		document.documentElement.setAttribute('data-theme', next);
		localStorage.setItem('safecastDocTheme', next);
		applyLabel();
	};
})();
</script>
</body>
</html>`

	fmt.Fprint(w, html)
}

// adminBackfillHandler backfills the uploads table with existing spectrum data.
// POST /api/admin/backfill?password=xxx
//
// @Summary     Admin backfill uploads from spectra
// @Description Backfills upload records from existing spectra entries.
// @Tags        admin
// @Produce     json
// @Success     200 {object} map[string]interface{} "Backfill result"
// @Failure     401 {string} string "Unauthorized"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/backfill [post]
func adminBackfillHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for admin access: session-based or password-based
	isSessionAdmin := false
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		isSessionAdmin = true
	}

	if !isSessionAdmin {
		// Fall back to password authentication
		if *adminPassword == "" {
			http.Error(w, "Admin endpoints are disabled - please login as admin", http.StatusForbidden)
			return
		}
		password := r.URL.Query().Get("password")
		if password != *adminPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	if db == nil || db.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	ctx := r.Context()

	// Query all spectra (including those without filenames for legacy data)
	query := `
		SELECT
			s.id,
			COALESCE(s.filename, ''),
			s.source_format,
			m.trackID,
			s.created_at
		FROM spectra s
		JOIN markers m ON s.marker_id = m.id
	`

	rows, err := db.DB.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error querying spectra: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  fmt.Sprintf("Failed to query spectra: %v", err),
		})
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var spectrumID int64
		var filename, sourceFormat, trackID string
		var createdAtTime time.Time

		if err := rows.Scan(&spectrumID, &filename, &sourceFormat, &trackID, &createdAtTime); err != nil {
			log.Printf("Error scanning spectrum row: %v", err)
			continue
		}
		createdAt := createdAtTime.Unix()

		// Generate filename for legacy records that don't have one
		if filename == "" {
			ext := ".unknown"
			if sourceFormat == "n42" {
				ext = ".n42"
			} else if sourceFormat == "spe" {
				ext = ".spe"
			}
			filename = fmt.Sprintf("spectrum_%d%s", spectrumID, ext)
		}

		// Create upload record
		upload := database.Upload{
			Filename:  filename,
			FileType:  sourceFormat,
			TrackID:   trackID,
			FileSize:  0, // Unknown for backfilled records
			UploadIP:  "backfilled",
			CreatedAt: createdAt,
		}

		if _, uploadErr := db.InsertUpload(ctx, upload); uploadErr != nil {
			// Skip if already exists (duplicate key error)
			if !strings.Contains(uploadErr.Error(), "UNIQUE constraint") &&
				!strings.Contains(uploadErr.Error(), "duplicate key") {
				log.Printf("Warning: failed to backfill upload record for %s: %v", filename, uploadErr)
			}
			continue
		}

		count++
	}

	log.Printf("Admin backfilled %d upload records", count)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Upload records backfilled successfully",
		"count":   count,
	})
}

// countryBackfillRunning guards against concurrent backfill runs.
var countryBackfillRunning sync.Mutex

// adminBackfillCountriesHandler populates the country column for markers that
// have coordinates but no country assigned.  The work runs in the background;
// the endpoint returns immediately so the caller is not blocked.
// POST /api/admin/backfill-countries?password=xxx
//
// @Summary     Admin backfill marker countries
// @Description Starts background country backfill for markers with coordinates.
// @Tags        admin
// @Produce     json
// @Success     200 {object} map[string]interface{} "Backfill status"
// @Failure     401 {string} string "Unauthorized"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/backfill-countries [post]
func adminBackfillCountriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for admin access: session-based or password-based
	isSessionAdmin := false
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		isSessionAdmin = true
	}
	if !isSessionAdmin {
		if *adminPassword == "" {
			http.Error(w, "Admin endpoints are disabled - please login as admin", http.StatusForbidden)
			return
		}
		password := r.URL.Query().Get("password")
		if password != *adminPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	if db == nil || db.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	if !countryBackfillRunning.TryLock() {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "already_running",
			"message": "Country backfill is already in progress.",
		})
		return
	}

	go func() {
		defer countryBackfillRunning.Unlock()
		backfillCountries()
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "started",
		"message": "Country backfill started in background. Watch server logs for progress.",
	})
}

// adminCacheHandler provides administrative access to cache management functions
//
// @Summary     Admin cache operations
// @Description Performs cache operations such as clear/stats via action query parameter.
// @Tags        admin
// @Produce     json
// @Param       action query string false "Cache action (e.g. clear, stats)"
// @Success     200 {object} map[string]interface{} "Cache operation result"
// @Failure     401 {string} string "Unauthorized"
// @Router      /api/admin/cache [post]
func adminCacheHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for admin access: session-based or password-based
	isSessionAdmin := false
	if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
		isSessionAdmin = true
	}
	if !isSessionAdmin {
		if *adminPassword == "" {
			http.Error(w, "Admin endpoints are disabled - please login as admin", http.StatusForbidden)
			return
		}
		password := r.URL.Query().Get("password")
		if password != *adminPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	action := r.URL.Query().Get("action")
	switch action {
	case "clear":
		tileCacheMu.Lock()
		tileCache.Purge()
		tileCacheMu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Cache cleared successfully",
		})
	case "stats":
		tileCacheMu.RLock()
		size := tileCache.Len()
		capacity := 1000 // Our cache capacity
		tileCacheMu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "success",
			"size":     size,
			"capacity": capacity,
		})
	default:
		http.Error(w, "Invalid action. Use ?action=clear or ?action=stats", http.StatusBadRequest)
		return
	}
}

func backfillCountries() {
	const batchSize = 50000
	const updateChunk = 10000 // max IDs per UPDATE statement

	log.Printf("country backfill: starting (batch size %d)", batchSize)
	start := time.Now()
	var lastID int64
	var totalUpdated int64

	for {
		rows, err := db.DB.QueryContext(context.Background(),
			`SELECT id, lat, lon FROM markers
			 WHERE id > $1 AND (country IS NULL OR country = '')
			 ORDER BY id LIMIT $2`, lastID, batchSize)
		if err != nil {
			log.Printf("country backfill: query error: %v", err)
			return
		}

		type markerRow struct {
			id       int64
			lat, lon float64
		}
		var batch []markerRow
		for rows.Next() {
			var r markerRow
			if err := rows.Scan(&r.id, &r.lat, &r.lon); err != nil {
				log.Printf("country backfill: scan error: %v", err)
				rows.Close()
				return
			}
			batch = append(batch, r)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			log.Printf("country backfill: rows error: %v", err)
			return
		}

		if len(batch) == 0 {
			break
		}
		lastID = batch[len(batch)-1].id

		// Resolve countries in parallel across CPU cores.
		type resolved struct {
			id      int64
			country string
		}
		results := make([]resolved, len(batch))
		workers := runtime.GOMAXPROCS(0)
		var wg sync.WaitGroup
		chunkSize := (len(batch) + workers - 1) / workers
		for w := 0; w < workers; w++ {
			lo := w * chunkSize
			hi := lo + chunkSize
			if lo >= len(batch) {
				break
			}
			if hi > len(batch) {
				hi = len(batch)
			}
			wg.Add(1)
			go func(lo, hi int) {
				defer wg.Done()
				for i := lo; i < hi; i++ {
					code, _ := countryresolver.Resolve(batch[i].lat, batch[i].lon)
					if code == "" {
						code = "unknown"
					}
					results[i] = resolved{id: batch[i].id, country: code}
				}
			}(lo, hi)
		}
		wg.Wait()

		// Group by country code.
		groups := make(map[string][]int64)
		for _, r := range results {
			groups[r.country] = append(groups[r.country], r.id)
		}

		// One UPDATE per country code, chunked to keep query size reasonable.
		for code, ids := range groups {
			for i := 0; i < len(ids); i += updateChunk {
				end := i + updateChunk
				if end > len(ids) {
					end = len(ids)
				}
				chunk := ids[i:end]
				ph := make([]string, len(chunk))
				args := make([]interface{}, 1+len(chunk))
				args[0] = code
				for j, id := range chunk {
					ph[j] = fmt.Sprintf("$%d", j+2)
					args[j+1] = id
				}
				query := fmt.Sprintf("UPDATE markers SET country = $1 WHERE id IN (%s)", strings.Join(ph, ","))
				if _, err := db.DB.ExecContext(context.Background(), query, args...); err != nil {
					log.Printf("country backfill: update error for %s: %v", code, err)
					return
				}
			}
			totalUpdated += int64(len(ids))
		}

		rate := float64(totalUpdated) / time.Since(start).Seconds()
		log.Printf("country backfill: id=%d, updated %d total (%.0f rows/sec)", lastID, totalUpdated, rate)
	}

	log.Printf("country backfill: finished, updated %d markers in %s", totalUpdated, time.Since(start).Round(time.Second))
}

// =====================
// WEB  — страница трека
// =====================
// trackHandler serves the single-track map page.
//
// @Summary     Render map page for one track
// @Description Serves the HTML map page for a specific track ID.
// @Tags        web
// @Produce     html
// @Param       id path string true "Track ID"
// @Success     200 {string} string "HTML page"
// @Failure     400 {string} string "Track ID missing"
// @Router      /trackid/{id} [get]
