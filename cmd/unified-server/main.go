//new: stream markers by track

package main

import (

	// http://localhost:8765/debug/pprof/profile?seconds=30
	// go tool pprof -http=:8080 Downloads/profile
	//_ "net/http/pprof"

	"context"
	"crypto/tls"
	"embed"

	"errors"
	"flag"
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"html"
	"html/template"
	"io/fs"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"

	_ "safecast-new-map/cmd/unified-server/docs/api"
	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/database"
	"safecast-new-map/pkg/database/drivers"
	"safecast-new-map/pkg/email"
	"safecast-new-map/pkg/httpapi"
	"safecast-new-map/pkg/jsonarchive"
	"safecast-new-map/pkg/logger"
	safecastfetcher "safecast-new-map/pkg/safecast-fetcher"
	safecastrealtime "safecast-new-map/pkg/safecast-realtime"
)

// content bundles the UI and the license texts so single-file binaries still
// expose the legal notice when served offline. Embedding keeps deployment
// simple and mirrors the "A little copying is better than a little dependency"
// proverb by avoiding extra runtime file IO.
//
//go:embed public_html/* LICENSE LICENSE.CC0
var content embed.FS

var doseData database.Data

var domain = flag.String("domain", "", "Serve HTTPS on 80/443 via Let's Encrypt when a domain is provided.")
var dbType = flag.String("db-type", "pgx", "Database driver: pgx (PostgreSQL, default), sqlite, chai, duckdb, or clickhouse")
var dbPath = flag.String("db-path", "", "Filesystem path for chai/sqlite/duckdb databases; defaults to the working directory.")
var dbConn = flag.String("db-conn", "", "Connection URI for network databases.\n  PostgreSQL: postgres://user:pass@host:5432/<database>?sslmode=verify-full\n  ClickHouse: clickhouse://user:pass@host:9000/<database>?secure=true")
var port = flag.Int("port", 8765, "Port for running the HTTP server when not using -domain.")
var version = flag.Bool("version", false, "Show the application version")
var defaultLat = flag.Float64("default-lat", 44.08832, "Default map latitude")
var defaultLon = flag.Float64("default-lon", 42.97577, "Default map longitude")
var defaultZoom = flag.Int("default-zoom", 11, "Default map zoom")
var defaultLayer = flag.String("default-layer", "OpenStreetMap", `Default base layer: "OpenStreetMap" or "Google Satellite"`)
var autoLocateDefault = flag.Bool("auto-locate-default", true, "Auto-center initial map view using browser or GeoIP fallbacks when no URL bounds are provided.")
var safecastRealtimeEnabled = flag.Bool("safecast-realtime", false, "Enable polling and display of Safecast realtime devices")
var jsonArchivePathFlag = flag.String("json-archive-path", "", "Filesystem destination for the generated JSON archive tgz bundle")
var jsonArchiveFrequencyFlag = flag.String("json-archive-frequency", "weekly", "How often to rebuild the JSON archive: daily, weekly, monthly, or yearly")
var importTGZURLFlag = flag.String("import-tgz-url", "", "Download and import a remote .tgz of exported JSON files, log progress, and exit once finished. Example: https://simplemap.safecast.org/api/json/weekly.tgz")
var importTGZFileFlag = flag.String("import-tgz-file", "", "Import a local .tgz of exported JSON files, log progress, and exit once finished.")
var supportEmail = flag.String("support-email", "", "Contact e-mail shown in the legal notice for feedback")
var debugIPsFlag = flag.String("debug", "", "Comma separated IP addresses allowed to view the debug overlay")
var adminPassword = flag.String("admin-password", "", "Password for admin endpoints (upload listing, track deletion). If not set, admin endpoints are disabled.")
var safecastFetcherEnabled = flag.Bool("safecast-fetcher", false, "Enable automatic fetching of approved bGeigie imports from api.safecast.org")
var safecastFetcherInterval = flag.Duration("safecast-fetcher-interval", 5*time.Minute, "How often to poll api.safecast.org for new approved imports")
var safecastFetcherBatchSize = flag.Int("safecast-fetcher-batch-size", 10, "Maximum number of files to import per polling cycle (0 = unlimited)")
var safecastFetcherStartDate = flag.String("safecast-fetcher-start-date", "", "Only import files uploaded after this date (YYYY-MM-DD format, empty = no filter)")
var safecastFetcherBackfill = flag.Bool("safecast-fetcher-backfill", false, "Backfill mode: import all matching records from start-date, ignoring what's already in database")
var safecastFetcherNewestFirst = flag.Bool("safecast-fetcher-newest-first", false, "Fetch newest imports first instead of oldest first")

// Authentication flags
var smtpHost = flag.String("smtp-host", "", "SMTP server hostname for sending emails")
var smtpPort = flag.Int("smtp-port", 587, "SMTP server port")
var smtpUsername = flag.String("smtp-username", "", "SMTP username")
var smtpPassword = flag.String("smtp-password", "", "SMTP password")
var smtpFrom = flag.String("smtp-from", "", "Email from address")
var smtpFromName = flag.String("smtp-from-name", "Safecast", "Email from name")
var sessionSecret = flag.String("session-secret", "", "Secret key for session encryption (required for authentication)")
var sessionCookieName = flag.String("session-cookie-name", "safecast_session", "Name of the session cookie")
var requireAuth = flag.Bool("require-auth", false, "Require authentication for uploads")
var allowRegistration = flag.Bool("allow-registration", true, "Allow new user registration")
var sessionDuration = flag.Duration("session-duration", 30*24*time.Hour, "Session duration (default: 30 days)")
var baseURL = flag.String("base-url", "http://localhost:8765", "Base URL for the application (used in emails)")

// debugIPAllowlist keeps a fast lookup of remote addresses that should see the
// technical overlay. We keep it as a map so lookups stay O(1) without extra
// synchronization, following "Clear is better than clever" by leaning on Go's
// built-in map semantics.
var debugIPAllowlist map[string]struct{}

// tileCache stores precomputed clustered markers to avoid expensive on-the-fly clustering
// for frequently requested map tiles. The cache uses a simple LRU strategy to limit memory usage.
var tileCache *lru.Cache[string, []database.Marker]
var tileCacheMu sync.RWMutex // Protects access to the tileCache

var CompileVersion = "dev"

var (
	apiDocsArchiveEnabled   bool
	apiDocsArchiveRoute     string
	apiDocsArchiveFrequency string
)

var db *database.Database

// Upload progress tracking
type UploadProgress struct {
	Total            int
	Current          int
	Complete         bool
	Error            string
	RedirectURL      string
	NeedsCoordinates bool
	TrackID          string
	FileName         string
	mu               sync.RWMutex
}

var uploadProgress = struct {
	sync.RWMutex
	tracks map[string]*UploadProgress
}{
	tracks: make(map[string]*UploadProgress),
}

func init() {
	// We trigger driver registration here so "go run safecast-new-map.go" keeps
	// working even when auxiliary files are skipped; relying on init avoids extra
	// coordination primitives and mirrors Go's preference for simplicity.
	drivers.Ready()

	// Initialize the tile cache to store precomputed clustered markers
	// Using a cache size of 1000 entries which should be sufficient for most use cases
	tileCache, _ = lru.New[string, []database.Marker](1000)

	// CLI usage grouping is also configured once during init so every entry point
	// inherits the readable help layout without repeating boilerplate.
	configureCLIUsage()
}

// resolveArchivePath decides where the JSON archive tgz should live.
// We prefer explicit destinations from flags, otherwise fall back to the user's
// home directory so long-running services do not clutter the repository tree.
// We log resolution failures so operators notice and can correct their setup.
// The defaultFile argument feeds through the configured cadence and domain so
// implicit directories still produce predictable filenames.
func resolveArchivePath(flagValue, defaultFile string, logf func(string, ...any)) string {
	cleaned := strings.TrimSpace(flagValue)
	fallback := strings.TrimSpace(defaultFile)
	if fallback == "" {
		fallback = "weekly-json.tgz"
	}
	if cleaned != "" {
		abs, err := filepath.Abs(cleaned)
		if err != nil {
			if logf != nil {
				logf("json archive path resolution fallback for %q: %v", cleaned, err)
			}
			return filepath.Clean(cleaned)
		}
		return abs
	}

	home, err := os.UserHomeDir()
	if err == nil && strings.TrimSpace(home) != "" {
		return filepath.Join(home, fallback)
	}

	// Falling back to the working directory keeps the archive predictable even
	// in minimal environments where HOME is undefined, trading cleverness for
	// clarity per the Go proverbs.
	wd, wdErr := os.Getwd()
	if wdErr == nil && strings.TrimSpace(wd) != "" {
		return filepath.Join(wd, fallback)
	}

	// As a last resort return a relative filename so the generator can still run.
	return fallback
}

// getEnvOrDefault returns the environment variable value or a default if not set.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault returns the environment variable value as int or a default if not set.
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// applyDBConnection parses a DSN passed via -db-conn and copies relevant fields into the
// database configuration. We normalise defaults for host, port, and SSL/TLS so operators can
// supply concise URLs while the rest of the application continues using structured settings.
func applyDBConnection(driverName, conn string, cfg *database.Config) error {
	if cfg == nil {
		return fmt.Errorf("db config is nil")
	}
	cleaned := strings.TrimSpace(conn)
	if cleaned == "" {
		return nil
	}

	parsed, err := url.Parse(cleaned)
	if err != nil {
		return fmt.Errorf("%s connection string: %w", driverName, err)
	}

	driver := strings.ToLower(strings.TrimSpace(driverName))
	switch driver {
	case "pgx":
		if parsed.Scheme == "" {
			parsed.Scheme = "postgres"
		}
	case "clickhouse":
		if parsed.Scheme == "" {
			parsed.Scheme = "clickhouse"
		}
	default:
		return fmt.Errorf("db-conn is only supported for pgx or clickhouse (got %q)", driverName)
	}

	host := parsed.Hostname()
	if host == "" {
		host = "127.0.0.1"
	}
	cfg.DBHost = host

	portValue := parsed.Port()
	var port int
	if portValue != "" {
		port, err = strconv.Atoi(portValue)
		if err != nil {
			return fmt.Errorf("%s connection string: invalid port %q", driverName, portValue)
		}
	} else {
		if driver == "pgx" {
			port = 5432
		} else {
			port = 9000
		}
	}
	cfg.DBPort = port

	if parsed.User != nil {
		if user := strings.TrimSpace(parsed.User.Username()); user != "" {
			cfg.DBUser = user
		}
		if pass, ok := parsed.User.Password(); ok {
			cfg.DBPass = pass
		}
	}

	name := strings.Trim(strings.TrimPrefix(parsed.Path, "/"), " ")
	if driver == "pgx" && name == "" {
		return fmt.Errorf("%s connection string must include a database name", driverName)
	}
	if name != "" || driver == "pgx" {
		cfg.DBName = name
	}

	query := parsed.Query()
	switch driver {
	case "pgx":
		sslMode := strings.TrimSpace(query.Get("sslmode"))
		if sslMode == "" {
			sslMode = "prefer"
			query.Set("sslmode", sslMode)
		}
		cfg.PGSSLMode = sslMode
	case "clickhouse":
		secureValue := strings.TrimSpace(query.Get("secure"))
		secure := false
		if secureValue != "" {
			secure = secureValue == "1" || strings.EqualFold(secureValue, "true") || strings.EqualFold(secureValue, "yes") || strings.EqualFold(secureValue, "on")
		} else if strings.EqualFold(parsed.Scheme, "https") {
			secure = true
			query.Set("secure", "true")
		}
		cfg.ClickSecure = secure
	}
	parsed.RawQuery = query.Encode()

	cfg.DBConn = parsed.String()
	return nil
}

var errNotSafecastTrackJSON = errors.New("not safecast track json payload")

// processBGeigieZenFile parses bGeigie Zen/Nano $BNRDD logs.
// Supports ISO8601 timestamps at field[2] and DMM coordinates with N/S/E/W.
func withServerHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "safecast-new-map/"+CompileVersion)

		if r.Method == http.MethodHead && r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// serveWithDomain запускает:
//   • :80  — ACME HTTP-01 + 301-redirect на https://<domain>/…
//   • :443 — HTTPS с автоматическими сертификатами Let’s Encrypt.
//
// Новое: если autocert не может выдать cert (любой host/SNI),
//        сервер всё-таки отдаёт ранее полученный fallback-cert,
//        тем самым устраняя «host not configured» в логах.
//
// Совместимость: TLS ≥ 1.0, ALPN h2/http1.1/http1.0.
// Все ошибки только логируются.

func serveWithDomain(domain string, handler http.Handler) {
	// ----------- ACME manager -----------
	certMgr := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("certs"),
		HostPolicy: func(ctx context.Context, host string) error {
			// Разрешаем голый и www.<domain>
			if host == domain || host == "www."+domain {
				return nil
			}
			// IP-адрес? — не блокируем, просто не пытаемся получить cert.
			if net.ParseIP(host) != nil {
				return nil
			}
			return errors.New("acme/autocert: host not configured")
		},
	}

	// ----------- :80 (challenge + redirect) -----------
	go func() {
		mux80 := http.NewServeMux()
		mux80.Handle("/.well-known/acme-challenge/", certMgr.HTTPHandler(nil))
		mux80.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + domain + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		})

		log.Printf("HTTP  server (ACME+redirect) ➜ :80")
		if err := (&http.Server{
			Addr:              ":80",
			Handler:           mux80,
			ReadHeaderTimeout: 10 * time.Second,
		}).ListenAndServe(); err != nil {
			selfupgradeHandleServerError(err, log.Printf)
		}
	}()

	// ----------- ежедневная проверка сертификата -----------
	go func() {
		t := time.NewTicker(24 * time.Hour)
		defer t.Stop()
		for range t.C {
			if _, err := certMgr.GetCertificate(&tls.ClientHelloInfo{ServerName: domain}); err != nil {
				log.Printf("autocert renewal check: %v", err)
			}
		}
	}()

	// ----------- :443 (HTTPS) -----------
	tlsCfg := certMgr.TLSConfig()
	tlsCfg.MinVersion = tls.VersionTLS10
	tlsCfg.NextProtos = append([]string{"http/1.0"}, tlsCfg.NextProtos...)

	// fallback-сертификат для IP / странных SNI
	var defaultCert *tls.Certificate
	go func() {
		for defaultCert == nil {
			if c, err := certMgr.GetCertificate(&tls.ClientHelloInfo{ServerName: domain}); err == nil {
				defaultCert = c
			}
			time.Sleep(time.Minute)
		}
	}()
	tlsCfg.GetCertificate = func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
		c, err := certMgr.GetCertificate(chi)
		if err == nil {
			return c, nil
		}
		// Любой сбой — пытаемся отдать fallback-cert (если уже есть)
		if defaultCert != nil {
			return defaultCert, nil
		}
		// Пока fallback нет — повторяем оригинальную ошибку
		return nil, err
	}

	log.Printf("HTTPS server for %s ➜ :443 (TLS ≥1.0, ALPN h2/http1.1/1.0)", domain)
	if err := (&http.Server{
		Addr:              ":443",
		Handler:           handler,
		TLSConfig:         tlsCfg,
		ReadHeaderTimeout: 10 * time.Second,
	}).ListenAndServeTLS("", ""); err != nil {
		selfupgradeHandleServerError(err, log.Printf)
	}
}

// logT формирует строку "[trackID][component] …" и передаёт её в пакет logger.
// logger сам решит: буферизовать или вывести сразу.
func logT(trackID, component, format string, v ...any) {
	line := fmt.Sprintf("[%-6s][%s] %s", trackID, component, fmt.Sprintf(format, v...))
	logger.Append(trackID, line)
}

// rxFind returns the first submatch (group #1) of pattern in s or an empty string.
// Entities are unescaped and result is TrimSpace-обработан.
func rxFind(s, pattern string) string {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(s)
	if len(m) > 1 {
		return strings.TrimSpace(html.UnescapeString(m[1]))
	}
	return ""
}

// isClientDisconnect returns true for network errors indicating that the client
// has gone away (e.g., browser navigated away or closed the tab) while we were
// writing the response. These are normal and should not be logged as errors.
func isClientDisconnect(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "broken pipe") || strings.Contains(msg, "connection reset by peer")
}

// GenerateSerialNumber генерирует TrackID
func GenerateSerialNumber() string {
	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const maxLength = 6

	timestamp := uint64(time.Now().UnixNano() / 1e6) // время в мс
	encoded := ""
	base := uint64(len(base62Chars))

	for timestamp > 0 && len(encoded) < maxLength {
		remainder := timestamp % base
		encoded = string(base62Chars[remainder]) + encoded
		timestamp = timestamp / base
	}

	rand.Seed(time.Now().UnixNano())
	for len(encoded) < maxLength {
		encoded += string(base62Chars[rand.Intn(len(base62Chars))])
	}

	return encoded
}

// enqueueArchiveImport writes the uploaded tgz to a temporary file and processes it in the background.
// Using a goroutine prevents the HTTP handler from blocking while still reusing the common parser.
func main() {
	// 1. Флаги и версии
	flag.Parse()
	debugIPAllowlist = parseDebugAllowlist(*debugIPsFlag)
	loadTranslationsFromFile(content, "public_html/translations.json")
	selfupgradeStartupDelay(log.Printf)

	archiveFrequency, freqErr := jsonarchive.ParseFrequency(*jsonArchiveFrequencyFlag)
	if freqErr != nil {
		log.Fatalf("json archive frequency: %v", freqErr)
	}
	archiveFileName := jsonarchive.FileName(*domain, archiveFrequency)

	if *version {
		fmt.Printf("safecast-new-map version %s\n", CompileVersion)
		return
	}

	// 2. Предупреждение о привилегиях (для :80 / :443)
	if *domain != "" && runtime.GOOS != "windows" && os.Geteuid() != 0 {
		log.Println("⚠  Binding to :80 / :443 requires super-user rights; run with sudo or as root.")
	}

	// 3. База данных
	driverName := strings.ToLower(strings.TrimSpace(*dbType))
	// Persist the normalized driver back into the flag so downstream helpers never
	// miss engine-specific branches because of incidental casing or whitespace.
	*dbType = driverName
	dbCfg := database.Config{
		DBType: driverName,
		DBPath: *dbPath,
		Port:   *port,
	}
	switch driverName {
	case "pgx":
		// Default PostgreSQL connection settings
		// Can be overridden via -db-conn flag or environment variables:
		// DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME, DB_SSLMODE
		dbCfg.DBHost = getEnvOrDefault("DB_HOST", "127.0.0.1")
		dbCfg.DBPort = getEnvIntOrDefault("DB_PORT", 5432)
		dbCfg.DBUser = getEnvOrDefault("DB_USER", "postgres")
		dbCfg.DBPass = getEnvOrDefault("DB_PASS", "")
		dbCfg.DBName = getEnvOrDefault("DB_NAME", "safecast")
		dbCfg.PGSSLMode = getEnvOrDefault("DB_SSLMODE", "prefer")
		// If -db-conn is provided, it overrides individual settings
		if err := applyDBConnection(driverName, *dbConn, &dbCfg); err != nil {
			log.Fatalf("DB config: %v", err)
		}
	case "clickhouse":
		dbCfg.DBHost = "127.0.0.1"
		dbCfg.DBPort = 9000
		dbCfg.DBName = "IsotopePathways"
		if err := applyDBConnection(driverName, *dbConn, &dbCfg); err != nil {
			log.Fatalf("DB config: %v", err)
		}
	default:
		if strings.TrimSpace(*dbConn) != "" {
			log.Fatalf("db-conn is only valid for pgx or clickhouse drivers (current: %s)", *dbType)
		}
	}
	var err error
	db, err = database.NewDatabase(dbCfg)
	if err != nil {
		log.Fatalf("DB init: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("DB close: %v", closeErr)
		}
	}()
	if err = db.InitSchema(dbCfg); err != nil {
		log.Fatalf("DB schema: %v", err)
	}
	if strings.EqualFold(driverName, "pgx") {
		db.SyncPostgreSQLSequence(context.Background(), log.Printf)
	}
	queueDuckDBMaintenanceAfterImport(driverName, db, log.Printf, "startup")

	// Seed translations into DB if empty. Tour step seeding also inserts
	// tour.<step_key>.text rows, so run it before loading translations so the
	// in-memory map picks up the new keys on first boot.
	seedTranslationsDB(content, "public_html/translations.json")
	seedTourStepsDB()
	loadTranslationsFromDB()

	// Initialize authentication system if configured
	var authManager *auth.Manager
	var emailSender *email.Sender
	if *sessionSecret != "" && *smtpHost != "" {
		// Initialize email sender
		emailConfig := email.SMTPConfig{
			Host:     *smtpHost,
			Port:     *smtpPort,
			Username: *smtpUsername,
			Password: *smtpPassword,
			From:     *smtpFrom,
			FromName: *smtpFromName,
			UseTLS:   true,
		}
		emailSender = email.NewSender(emailConfig)

		// Initialize auth manager
		authManager = &auth.Manager{
			DB:                db.DB,
			DBDriver:          driverName,
			SessionCookieName: *sessionCookieName,
			AllowRegistration: *allowRegistration,
			EmailSender:       emailSender,
			BaseURL:           *baseURL,
		}

		log.Printf("Authentication system enabled")
	} else if *requireAuth {
		log.Fatalf("Authentication required but not properly configured. Set -session-secret and -smtp-host")
	}

	remoteURL := strings.TrimSpace(*importTGZURLFlag)
	localArchive := strings.TrimSpace(*importTGZFileFlag)
	if remoteURL != "" && localArchive != "" {
		log.Fatalf("choose only one import flag: -import-tgz-url or -import-tgz-file")
	}

	var importDone <-chan struct{}

	if localArchive != "" {
		fallback := GenerateSerialNumber()
		importDone = startBackgroundArchiveImport(context.Background(), fmt.Sprintf("local file %s", localArchive), func(ctx context.Context) error {
			return importArchiveFromFile(ctx, localArchive, fallback, db, driverName, log.Printf)
		}, log.Printf)
	}

	if remoteURL != "" {
		fallback := GenerateSerialNumber()
		importDone = startBackgroundArchiveImport(context.Background(), fmt.Sprintf("remote url %s", remoteURL), func(ctx context.Context) error {
			return importArchiveFromURL(ctx, remoteURL, fallback, db, driverName, log.Printf)
		}, log.Printf)
	}

	if *safecastRealtimeEnabled {
		// Launch realtime Safecast polling under the dedicated flag so the
		// feature stays opt-in.
		database.SetRealtimeConverter(safecastrealtime.FromRealtime)
		ctxRT, cancelRT := context.WithCancel(context.Background())
		defer cancelRT()
		safecastrealtime.Start(ctxRT, db, *dbType, log.Printf)
	}

	if *safecastFetcherEnabled {
		// Launch Safecast API fetcher to automatically import approved bGeigie logs
		ctxFetcher, cancelFetcher := context.WithCancel(context.Background())
		defer cancelFetcher()

		// Create importer function that bridges to existing import logic
		importerFunc := func(
			ctx context.Context,
			fileContent []byte,
			filename string,
			safecastImportID int64,
			sourceURL string,
			userID string,
			username string,
			comment string,
			db *database.Database,
			dbType string,
		) (trackID string, markerCount int, err error) {
			trackID = GenerateSerialNumber()

			// Create BytesFile from content
			bytesFile := safecastfetcher.NewBytesFile(fileContent, filename)

			// Parse and import the file
			_, finalTrackID, err := processBGeigieZenFile(bytesFile, trackID, db, dbType)
			if err != nil {
				return "", 0, fmt.Errorf("process bGeigie file: %w", err)
			}

			// Count markers for this track by querying the database
			countQuery := "SELECT COUNT(*) FROM markers WHERE trackID = ?"
			if dbType == "pgx" || dbType == "postgres" || dbType == "duckdb" {
				countQuery = "SELECT COUNT(*) FROM markers WHERE trackID = $1"
			}
			err = db.DB.QueryRowContext(ctx, countQuery, finalTrackID).Scan(&markerCount)
			if err != nil {
				// If count fails, log but don't fail the import
				log.Printf("[safecast-fetcher] warning: failed to count markers for track %s: %v", finalTrackID, err)
				markerCount = 0
			}

			// Get earliest marker date for this track to populate RecordingDate
			var recordingDate int64
			minDateQuery := "SELECT MIN(date) FROM markers WHERE trackID = ?"
			if dbType == "pgx" {
				minDateQuery = "SELECT MIN(date) FROM markers WHERE trackID = $1"
			}
			_ = db.DB.QueryRowContext(ctx, minDateQuery, finalTrackID).Scan(&recordingDate)

			// Get detector info from markers table
			var detector string
			detectorQuery := "SELECT COALESCE(detector, '') FROM markers WHERE trackID = ? AND detector != '' LIMIT 1"
			if dbType == "pgx" {
				detectorQuery = "SELECT COALESCE(detector, '') FROM markers WHERE trackID = $1 AND detector != '' LIMIT 1"
			}
			_ = db.DB.QueryRowContext(ctx, detectorQuery, finalTrackID).Scan(&detector)

			// Record the upload with source tracking
			upload := database.Upload{
				Filename:      filename,
				FileType:      ".log",
				TrackID:       finalTrackID,
				FileSize:      int64(len(fileContent)),
				UploadIP:      "safecast-api",
				CreatedAt:     time.Now().Unix(),
				RecordingDate: recordingDate, // Earliest marker date from log file
				Source:        safecastfetcher.SourceTypeSafecastAPI,
				SourceID:      fmt.Sprintf("%d", safecastImportID),
				SourceURL:     sourceURL,
				UserID:        userID,
				Username:      username,
				Detector:      detector,
				Comment:       comment,
			}

			if _, err := db.InsertUpload(ctx, upload); err != nil {
				return finalTrackID, markerCount, fmt.Errorf("record upload: %w", err)
			}

			return finalTrackID, markerCount, nil
		}

		// Set the default importer so GUI batch imports can use the same efficient code path
		safecastfetcher.DefaultImporter = importerFunc

		// Start the fetcher
		safecastfetcher.Start(ctxFetcher, safecastfetcher.Config{
			DB:           db,
			DBType:       *dbType,
			Interval:     *safecastFetcherInterval,
			BatchSize:    *safecastFetcherBatchSize,
			StartDate:    *safecastFetcherStartDate,
			Importer:     importerFunc,
			Logf:         log.Printf,
			BackfillMode: *safecastFetcherBackfill,
			NewestFirst:  *safecastFetcherNewestFirst,
		})

		log.Printf("safecast API fetcher enabled: interval=%s batch=%d start_date=%s backfill=%v newest_first=%v",
			*safecastFetcherInterval, *safecastFetcherBatchSize, *safecastFetcherStartDate, *safecastFetcherBackfill, *safecastFetcherNewestFirst)
	}

	// Build a JSON archive tgz with all known exported tracks only when
	// operators explicitly opt in via -json-archive-path. The cadence flag
	// keeps IO predictable while letting deployments choose how fresh the
	// bundle should be.
	var (
		archiveGen     *jsonarchive.Generator
		archiveCancel  context.CancelFunc
		archivePath    string
		archiveEnabled = strings.TrimSpace(*jsonArchivePathFlag) != ""
	)
	if archiveEnabled {
		ctxArchive, cancelArchive := context.WithCancel(context.Background())
		archiveCancel = cancelArchive
		archivePath = resolveArchivePath(*jsonArchivePathFlag, archiveFileName, log.Printf)
		if abs, err := filepath.Abs(archivePath); err == nil {
			archivePath = abs
		}
		log.Printf("json archive destination resolved: %s", archivePath)
		archiveGen = jsonarchive.Start(ctxArchive, db, *dbType, archivePath, archiveFileName, archiveFrequency.Interval(), log.Printf)
	} else {
		log.Printf("json archive disabled: set -json-archive-path to enable tarball generation")
	}
	if archiveCancel != nil {
		defer archiveCancel()
	}

	apiDocsArchiveEnabled = archiveEnabled
	route := archiveFrequency.RoutePath()
	if strings.TrimSpace(route) == "" {
		route = "/api/json/weekly.tgz"
	}
	apiDocsArchiveRoute = route
	apiDocsArchiveFrequency = archiveFrequency.HumanInterval()

	// Web server: handlers in pkg/httpapi.
	webConfig := httpapi.WebConfig{
		Domain:                  *domain,
		Port:                    *port,
		DefaultLat:              *defaultLat,
		DefaultLon:              *defaultLon,
		DefaultZoom:             *defaultZoom,
		DefaultLayer:            *defaultLayer,
		AutoLocateDefault:       *autoLocateDefault,
		SupportEmail:            *supportEmail,
		CompileVersion:          CompileVersion,
		DBType:                  *dbType,
		AdminPassword:           *adminPassword,
		APIDocsArchiveEnabled:   apiDocsArchiveEnabled,
		APIDocsArchiveRoute:     apiDocsArchiveRoute,
		APIDocsArchiveFrequency: apiDocsArchiveFrequency,
		DebugIPAllowlist:        debugIPAllowlist,
	}
	webServer := httpapi.NewWebServer(db, content, webConfig, log.Printf)

	// 4. Маршруты и статика
	staticFS, err := fs.Sub(content, "public_html")
	if err != nil {
		log.Fatalf("static fs: %v", err)
	}

	// Static asset routes are registered through the static registrar to keep
	// main() focused on composition rather than per-path wiring.
	httpapi.RegisterStaticRoutes(http.DefaultServeMux, httpapi.StaticRoutesConfig{
		StaticFS: staticFS,
		JSDir:    "public_html/",
	})
	mcpPortForDocs := strings.TrimSpace(os.Getenv("MCP_PORT"))
	if mcpPortForDocs == "" {
		mcpPortForDocs = "3333"
	}
	mcpBaseForDocs := strings.TrimSpace(os.Getenv("MCP_BASE_URL"))
	if mcpBaseForDocs == "" {
		mcpBaseForDocs = fmt.Sprintf("http://localhost:%s", mcpPortForDocs)
	}
	mcpDocsURL := strings.TrimRight(mcpBaseForDocs, "/") + "/mcp-api/"

	registerMainAPIDocsRoutes(http.DefaultServeMux, mcpDocsURL)

	// Stories page and its data feed
	storiesPageHandler := func(w http.ResponseWriter, r *http.Request) {
		data, err := content.ReadFile("public_html/stories.html")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	}
	storiesDataHandler := func(w http.ResponseWriter, r *http.Request) {
		data, err := content.ReadFile("public_html/data/stories.json")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}

	var profilePageHandler http.HandlerFunc
	var resetPasswordPageHandler http.HandlerFunc
	var adminUsersPageHandler http.HandlerFunc
	var adminUploadsPageHandler http.HandlerFunc
	var adminMCPPageHandler http.HandlerFunc
	var adminRealtimePageHandler http.HandlerFunc
	var adminTranslationsPageHandler http.HandlerFunc
	var adminTourPageHandler http.HandlerFunc
	var adminAIHintsPageHandler http.HandlerFunc

	// Register authentication and admin page handlers only when auth is enabled.
	if authManager != nil {
		profilePageHandler = func(w http.ResponseWriter, r *http.Request) {
			// Prevent CloudFront from caching user-specific pages
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			lang := getPreferredLanguage(r)
			translationsMu.RLock()
			localTranslations := translations
			translationsMu.RUnlock()

			tmpl, err := template.New("profile.html").Funcs(template.FuncMap{
				"translate": func(key string) string {
					if val, ok := localTranslations[lang][key]; ok {
						return val
					}
					return localTranslations["en"][key]
				},
			}).ParseFS(content, "public_html/profile.html")
			if err != nil {
				http.Error(w, "Profile page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			tmpl.Execute(w, nil)
		}

		resetPasswordPageHandler = func(w http.ResponseWriter, r *http.Request) {
			data, err := content.ReadFile("public_html/reset-password.html")
			if err != nil {
				http.Error(w, "Reset password page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}

		adminUsersPageHandler = func(w http.ResponseWriter, r *http.Request) {
			// Prevent CloudFront from caching this dynamic admin page
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			data, err := content.ReadFile("public_html/admin-users.html")
			if err != nil {
				http.Error(w, "Admin page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}

		adminUploadsPageHandler = func(w http.ResponseWriter, r *http.Request) {
			// Forward to the API endpoint which handles the uploads listing.
			adminUploadsHandler(w, r)
		}

		adminMCPPageHandler = func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			data, err := content.ReadFile("public_html/admin-mcp.html")
			if err != nil {
				http.Error(w, "Admin MCP page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}

		adminRealtimePageHandler = func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			data, err := content.ReadFile("public_html/admin-realtime.html")
			if err != nil {
				http.Error(w, "Admin Realtime page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}

		adminTranslationsPageHandler = func(w http.ResponseWriter, r *http.Request) {
			data, err := content.ReadFile("public_html/admin-translations.html")
			if err != nil {
				http.Error(w, "Page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}

		adminTourPageHandler = func(w http.ResponseWriter, r *http.Request) {
			data, err := content.ReadFile("public_html/admin-tour.html")
			if err != nil {
				http.Error(w, "Page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}

		adminAIHintsPageHandler = func(w http.ResponseWriter, r *http.Request) {
			data, err := content.ReadFile("public_html/admin-ai-hints.html")
			if err != nil {
				http.Error(w, "Page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}
	}

	httpapi.RegisterPageRoutes(http.DefaultServeMux, httpapi.PageRoutesConfig{
		AuthManager:                  authManager,
		AdminPassword:                *adminPassword,
		HomeHandler:                  homeHandler,
		StoriesPageHandler:           storiesPageHandler,
		StoriesDataHandler:           storiesDataHandler,
		MapHandler:                   mapHandler,
		ProfileHandler:               profilePageHandler,
		ResetPasswordHandler:         resetPasswordPageHandler,
		AdminUsersPageHandler:        adminUsersPageHandler,
		AdminUploadsPageHandler:      adminUploadsPageHandler,
		AdminMCPPageHandler:          adminMCPPageHandler,
		AdminRealtimePageHandler:     adminRealtimePageHandler,
		AdminTranslationsPageHandler: adminTranslationsPageHandler,
		AdminTourPageHandler:         adminTourPageHandler,
		AdminAIHintsPageHandler:      adminAIHintsPageHandler,
	})

	// Legacy public endpoints (non-/api) live in one registrar to keep route
	// ownership explicit and avoid growth of direct HandleFunc wiring in main().
	httpapi.RegisterLegacyRoutes(http.DefaultServeMux, httpapi.LegacyRoutesConfig{
		AuthManager:            authManager,
		RequireAuth:            *requireAuth,
		UploadHandler:          uploadHandler,
		UploadProgressHandler:  progressHandler,
		GetMarkersHandler:      getMarkersHandler,
		StreamMarkersHandler:   streamMarkersHandler,
		RealtimeHistoryHandler: realtimeHistoryHandler,
		TrackByIDHandler:       trackHandler,
		TracksByPrefixHandler:  tracksHandler,
	})
	// api/docs, licenses/, api/geoip, s/, api/spectrum/, api/markers/spectra, api/tracks/bounds, api/track-info/, api/update-coordinates, qrpng — registered via webServer.Register above
	// API endpoints ship JSON/archives. Keeping registration close to other
	// routes avoids surprises for operators scanning main() for handlers.
	limiter := httpapi.NewRateLimiter(time.Minute)
	apiHandler := httpapi.NewHandler(db, *dbType, archiveGen, limiter, log.Printf, archiveFrequency)
	restHandler := &RESTHandler{}

	// Keep MCP/realtime/translations admin APIs aligned with the legacy behavior:
	// they are registered only when auth is configured.
	var adminMCPDataAPIHandler http.HandlerFunc
	var adminMCPExportAPIHandler http.HandlerFunc
	var adminMCPDeleteAPIHandler http.HandlerFunc
	var adminMCPUpdateAPIHandler http.HandlerFunc
	var adminRealtimeDataAPIHandler http.HandlerFunc
	var adminRealtimeExportAPIHandler http.HandlerFunc
	var adminRealtimeDeleteAPIHandler http.HandlerFunc
	var adminTranslationsReloadAPIHandler http.HandlerFunc
	var adminTranslationByIDAPIHandler http.HandlerFunc
	var adminTranslationsAPIHandler http.HandlerFunc
	var adminTourStepsReorderAPIHandler http.HandlerFunc
	var adminTourStepsByIDAPIHandler http.HandlerFunc
	var adminTourStepsAPIHandler http.HandlerFunc
	var adminAIHintsReloadAPIHandler http.HandlerFunc
	var adminAIHintsImportAPIHandler http.HandlerFunc
	var adminAIHintsByIDAPIHandler http.HandlerFunc
	var adminAIHintsAPIHandler http.HandlerFunc
	if authManager != nil {
		adminMCPDataAPIHandler = adminMCPDataHandler
		adminMCPExportAPIHandler = adminMCPExportHandler
		adminMCPDeleteAPIHandler = adminMCPDeleteHandler
		adminMCPUpdateAPIHandler = adminMCPUpdateHandler
		adminRealtimeDataAPIHandler = adminRealtimeDataHandler
		adminRealtimeExportAPIHandler = adminRealtimeExportHandler
		adminRealtimeDeleteAPIHandler = adminRealtimeDeleteHandler
		adminTranslationsReloadAPIHandler = adminTranslationsReloadHandler
		adminTranslationByIDAPIHandler = func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPut:
				adminTranslationUpdateHandler(w, r)
			case http.MethodDelete:
				adminTranslationDeleteHandler(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
		adminTranslationsAPIHandler = func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				adminTranslationsDataHandler(w, r)
			case http.MethodPost:
				adminTranslationCreateHandler(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
		adminTourStepsReorderAPIHandler = func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			adminTourStepsReorderHandler(w, r)
		}
		adminTourStepsByIDAPIHandler = adminTourStepsByIDHandler
		adminTourStepsAPIHandler = func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				adminTourStepsListHandler(w, r)
			case http.MethodPost:
				adminTourStepsCreateHandler(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}

		adminAIHintsReloadAPIHandler = func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			adminAIHintReloadHandler(w, r)
		}
		adminAIHintsImportAPIHandler = func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			adminAIHintImportHandler(w, r)
		}
		// Subtree dispatcher for /api/admin/ai-hints/{model}[/history[/{id}/restore]|/restore|/snapshot|/export].
		adminAIHintsByIDAPIHandler = func(w http.ResponseWriter, r *http.Request) {
			stripped := strings.TrimPrefix(r.URL.Path, "/api/admin/ai-hints/")
			stripped = strings.Trim(stripped, "/")
			parts := strings.Split(stripped, "/")
			switch {
			case len(parts) == 1:
				// /api/admin/ai-hints/{model}
				switch r.Method {
				case http.MethodGet:
					adminAIHintGetHandler(w, r)
				case http.MethodPut, http.MethodPatch:
					adminAIHintUpdateHandler(w, r)
				case http.MethodDelete:
					adminAIHintDeleteHandler(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			case len(parts) == 2 && parts[1] == "restore" && r.Method == http.MethodPost:
				adminAIHintRestoreHandler(w, r)
			case len(parts) == 2 && parts[1] == "snapshot" && r.Method == http.MethodPost:
				adminAIHintSnapshotHandler(w, r)
			case len(parts) == 2 && parts[1] == "export" && r.Method == http.MethodGet:
				adminAIHintExportHandler(w, r)
			case len(parts) == 2 && parts[1] == "history" && r.Method == http.MethodGet:
				adminAIHintHistoryListHandler(w, r)
			case len(parts) == 4 && parts[1] == "history" && parts[3] == "restore" && r.Method == http.MethodPost:
				adminAIHintHistoryRestoreHandler(w, r)
			default:
				http.Error(w, "Not found", http.StatusNotFound)
			}
		}
		adminAIHintsAPIHandler = func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				adminAIHintsListHandler(w, r)
			case http.MethodPost:
				adminAIHintCreateHandler(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	}

	httpapi.Register(http.DefaultServeMux, httpapi.RegisterConfig{
		WebServer:                        webServer,
		APIHandler:                       apiHandler,
		AuthManager:                      authManager,
		DB:                               db,
		AdminPassword:                    *adminPassword,
		Logf:                             log.Printf,
		AdminUploadsHandler:              adminUploadsHandler,
		AdminTracksHandler:               adminTracksHandler,
		AdminBackfillHandler:             adminBackfillHandler,
		AdminBackfillCountriesHandler:    adminBackfillCountriesHandler,
		AdminDeleteTrackHandler:          adminDeleteTrackHandler,
		AdminDeleteMultipleTracksHandler: adminDeleteMultipleTracksHandler,
		AdminImportFromSafecastHandler:   adminImportFromSafecastHandler,
		AdminImportByIDHandler:           adminImportByIDHandler,
		AdminUpdateTrackHandler:          adminUpdateTrackHandler,
		AdminUpdateUploadHandler:         adminUpdateUploadHandler,
		AdminImportSafecastMetaHandler:   adminImportSafecastMetadataHandler,
		AdminCacheHandler:                adminCacheHandler,
		AdminMCPDataHandler:              adminMCPDataAPIHandler,
		AdminMCPExportHandler:            adminMCPExportAPIHandler,
		AdminMCPDeleteHandler:            adminMCPDeleteAPIHandler,
		AdminMCPUpdateHandler:            adminMCPUpdateAPIHandler,
		AdminRealtimeDataHandler:         adminRealtimeDataAPIHandler,
		AdminRealtimeExportHandler:       adminRealtimeExportAPIHandler,
		AdminRealtimeDeleteHandler:       adminRealtimeDeleteAPIHandler,
		AdminTranslationsReloadHandler:   adminTranslationsReloadAPIHandler,
		AdminTranslationByIDHandler:      adminTranslationByIDAPIHandler,
		AdminTranslationsHandler:         adminTranslationsAPIHandler,
		APITourStepsHandler:              tourStepsPublicHandler,
		AdminTourStepsReorderHandler:     adminTourStepsReorderAPIHandler,
		AdminTourStepsByIDHandler:        adminTourStepsByIDAPIHandler,
		AdminTourStepsHandler:            adminTourStepsAPIHandler,
		AdminAIHintsReloadHandler:        adminAIHintsReloadAPIHandler,
		AdminAIHintsImportHandler:        adminAIHintsImportAPIHandler,
		AdminAIHintsByIDHandler:          adminAIHintsByIDAPIHandler,
		AdminAIHintsHandler:              adminAIHintsAPIHandler,
		APISensorsHandler:                restHandler.handleSensors,
		APISensorsExportHandler:          restHandler.handleSensorsExport,
		APISensorByIDHandler:             restHandler.handleSensor,
		APIFeedbackHandler:               handleFeedback(),
		APITrackInsightsHandler:          trackInsightsHandler,
	})

	// Geocoding proxy — lets the browser search for places without hitting
	// Nominatim directly (browsers cannot set User-Agent in fetch()).
	http.DefaultServeMux.HandleFunc("/api/geocode", handleGeocode)

	// Register MCP Server (AI assistant, REST API, Swagger) on port 3333
	// Uses existing PostgreSQL (db) and DuckDB (duckDB) connections
	RegisterMCP()

	// Selfupgrade runs in the background only when explicitly enabled so existing
	// installations keep their manual release cadence. We assemble the config
	// near main() so filesystem paths, database settings, and HTTP handlers stay
	// consistent with the rest of the binary.
	selfUpgradeCancel := startSelfUpgrade(context.Background(), dbCfg)
	if selfUpgradeCancel != nil {
		defer selfUpgradeCancel()
	}

	var rootHandler http.Handler = http.DefaultServeMux
	if shield := importShield(importDone, driverName, log.Printf); shield != nil {
		// Keep HTTP responsive while a single-user DB import runs by declining
		// DB-backed endpoints. The middleware only activates for file engines
		// so multi-user databases remain fully live during imports.
		rootHandler = shield(rootHandler)
	}
	rootHandler = withServerHeader(rootHandler)

	// 5. HTTP/HTTPS-серверы
	if *domain != "" {
		// Двойной сервер :80 + :443 с Let’s Encrypt
		go serveWithDomain(*domain, rootHandler)
	} else {
		// Обычный HTTP на порт из -port
		addr := fmt.Sprintf(":%d", *port)
		go func() {
			log.Printf("HTTP server ➜ http://localhost:%d", *port)
			if err := http.ListenAndServe(addr, rootHandler); err != nil {
				selfupgradeHandleServerError(err, log.Printf)
			}
		}()
	}

	// асинхронные индексы в бд без блокирования основного процесса начало
	ctxIdx, cancelIdx := context.WithCancel(context.Background())
	defer cancelIdx()
	// Пояснение в лог: что делаем и почему это не блокирует сервер
	log.Printf("⏳ background index build scheduled (engine=%s). Listeners are up; pages may be slower until indexes are ready.", dbCfg.DBType)
	// Запуск асинхронной индексации с прогрессом
	db.EnsureIndexesAsync(ctxIdx, dbCfg, func(format string, args ...any) {
		log.Printf(format, args...)
	})
	// асинхронные индексы в бд без блокирования основного процесса конец

	// 6. Держим main-goroutine живой
	select {}
}
