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
	httpSwagger "github.com/swaggo/http-swagger"
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
	"github.com/swaggo/swag"
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
	queueDuckDBMaintenanceAfterImport(driverName, db, log.Printf, "startup")

	// Seed translations into DB if empty, then reload from DB
	seedTranslationsDB(content, "public_html/translations.json")
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

	// Serve static files from embedded filesystem - this must come BEFORE the catch-all route
	// to avoid the map handler catching static file requests
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.FS(staticFS))))

	// Serve JS files from the physical directory as a workaround
	// This ensures the marker-worker.js file is accessible to the browser
	// Access files from public_html root and let StripPrefix handle the path
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("public_html/"))))
	mcpPortForDocs := strings.TrimSpace(os.Getenv("MCP_PORT"))
	if mcpPortForDocs == "" {
		mcpPortForDocs = "3333"
	}
	mcpBaseForDocs := strings.TrimSpace(os.Getenv("MCP_BASE_URL"))
	if mcpBaseForDocs == "" {
		mcpBaseForDocs = fmt.Sprintf("http://localhost:%s", mcpPortForDocs)
	}
	mcpDocsURL := strings.TrimRight(mcpBaseForDocs, "/") + "/mcp-api/"

	// Combined API docs page (Map API + MCP API tabs)
	http.HandleFunc("/docs/", serveAPIDocsPage)
	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	// Register Map API favicon and theme CSS endpoints
	http.HandleFunc("/map-api/favicon.ico", serveFavicon)
	http.HandleFunc("/map-api/favicon-16x16.png", serveFavicon16)
	http.HandleFunc("/map-api/favicon-32x32.png", serveFavicon32)
	http.HandleFunc("/map-api/swagger-theme.css", serveMapSwaggerTheme)

	mapAPINavScript := fmt.Sprintf(`function() {
				document.title = 'Safecast Map API Docs';
				const mcpDocsURL = %q;

				// ── Favicons ──
				const link16 = document.createElement('link');
				link16.rel = 'icon'; link16.type = 'image/png'; link16.sizes = '16x16';
				link16.href = '/map-api/favicon-16x16.png';
				document.head.appendChild(link16);
				const link32 = document.createElement('link');
				link32.rel = 'icon'; link32.type = 'image/png'; link32.sizes = '32x32';
				link32.href = '/map-api/favicon-32x32.png';
				document.head.appendChild(link32);
				const linkICO = document.createElement('link');
				linkICO.rel = 'shortcut icon'; linkICO.href = '/map-api/favicon.ico';
				document.head.appendChild(linkICO);

				// ── Theme CSS ──
				const style = document.createElement('link');
				style.rel = 'stylesheet'; style.href = '/map-api/swagger-theme.css';
				document.head.appendChild(style);

				// ── Remove Swagger logo ──
				const swaggerLogo = document.querySelector('.topbar-wrapper .link');
				if (swaggerLogo) swaggerLogo.remove();
				document.querySelectorAll('.topbar-wrapper img').forEach(img => img.remove());

				// ── Inject "Switch to MCP API" button into topbar ──
				const topbar = document.querySelector('.swagger-ui .topbar');
				if (topbar) {
					const existingBtn = document.getElementById('safecast-switch-btn');
					if (existingBtn) existingBtn.remove();
					const switchBtn = document.createElement('a');
					switchBtn.id = 'safecast-switch-btn';
					switchBtn.href = mcpDocsURL;
					switchBtn.textContent = 'Switch to MCP API \u2192';
					switchBtn.style.cssText = 'display:inline-block;margin-left:auto;margin-right:16px;padding:6px 14px;background:#0a4f8a;color:#fff;border-radius:6px;font:600 13px/1.4 -apple-system,BlinkMacSystemFont,sans-serif;text-decoration:none;border:1px solid rgba(255,255,255,0.25);white-space:nowrap;';
					switchBtn.onmouseover = function() { this.style.background = '#083d6e'; };
					switchBtn.onmouseout  = function() { this.style.background = '#0a4f8a'; };
					const wrapper = topbar.querySelector('.topbar-wrapper');
					if (wrapper) {
						wrapper.style.display = 'flex';
						wrapper.style.alignItems = 'center';
						wrapper.style.width = '100%%';
						wrapper.appendChild(switchBtn);
					} else {
						topbar.appendChild(switchBtn);
					}
				}

				// ── Dark mode toggle ──
				const btn = document.createElement('button');
				btn.id = 'dark-mode-toggle';
				btn.textContent = '\u{1F319} Dark Mode';
				const isDark = localStorage.getItem('safecastMapDarkMode') === 'true';
				if (isDark) { document.body.classList.add('dark-mode'); btn.textContent = '\u2600\uFE0F Light Mode'; }
				btn.onclick = function() {
					document.body.classList.toggle('dark-mode');
					const nowDark = document.body.classList.contains('dark-mode');
					btn.textContent = nowDark ? '\u2600\uFE0F Light Mode' : '\u{1F319} Dark Mode';
					localStorage.setItem('safecastMapDarkMode', nowDark);
				};
				document.body.appendChild(btn);

				// ── Dark-mode styles for preamble ──
				const dmStyle = document.createElement('style');
				dmStyle.textContent = [
					'body.dark-mode #safecast-preamble { background: #1a2535 !important; border-bottom-color: #0d9488 !important; }',
					'body.dark-mode #safecast-preamble h2 { color: #93c5fd !important; }',
					'body.dark-mode #safecast-preamble > div > p { color: #b0b8c8 !important; }',
					'body.dark-mode #safecast-preamble a[href*="creativecommons"] { color: #5eead4 !important; }',
					'body.dark-mode #safecast-preamble summary { color: #93c5fd !important; }',
					'body.dark-mode #safecast-preamble details > div > div { background: #0f1c2e !important; border-color: #2a3f5f !important; }',
					'body.dark-mode #safecast-preamble details > div > div strong { color: #93c5fd !important; }',
					'body.dark-mode #safecast-preamble details > div > div p { color: #8899aa !important; }',
					'body.dark-mode #safecast-preamble details > div > p { color: #6b7a8d !important; }',
					'body.dark-mode #safecast-preamble code { background: #0f1c2e !important; color: #5eead4 !important; }',
				].join('\n');
				document.head.appendChild(dmStyle);

				// ── Preamble section (inserted before #swagger-ui) ──
				const existing = document.getElementById('safecast-preamble');
				if (existing) existing.remove();
				const preamble = document.createElement('div');
				preamble.id = 'safecast-preamble';
				preamble.style.cssText = 'font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;background:#fff;border-bottom:3px solid #0066cc;padding:28px 32px 24px;line-height:1.6;';
				preamble.innerHTML = '<div style="max-width:900px;margin:0 auto;">' +
					'<h2 style="margin:0 0 6px;font-size:22px;color:#1a3a5c;">Safecast Map API</h2>' +
					'<p style="margin:0 0 16px;font-size:15px;color:#555;">' +
						'Safecast has collected over 200 million radiation measurements from sensors carried by volunteers ' +
						'and fixed monitoring stations around the world. This page is the technical interface that lets software ' +
						'applications query, filter and download all of that data. You do not need an account or API key &mdash; ' +
						'all data is <a href="https://creativecommons.org/publicdomain/zero/1.0/" target="_blank" style="color:#0066cc;">CC0 licensed</a> and freely accessible.' +
					'</p>' +
					'<details style="margin-bottom:16px;">' +
						'<summary style="cursor:pointer;font-weight:600;font-size:14px;color:#1a3a5c;user-select:none;">For developers &mdash; endpoint overview</summary>' +
						'<div style="margin-top:10px;display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:10px;">' +
							'<div style="background:#f0f6ff;border:1px solid #c8ddf8;border-radius:8px;padding:12px;">' +
								'<strong style="color:#1a3a5c;">Historical</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">bGeigie mobile track measurements from citizen scientists worldwide.</p>' +
							'</div>' +
							'<div style="background:#f0f6ff;border:1px solid #c8ddf8;border-radius:8px;padding:12px;">' +
								'<strong style="color:#1a3a5c;">Realtime Sensors</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">Fixed Pointcast / Solarcast station readings &mdash; current &amp; history.</p>' +
							'</div>' +
							'<div style="background:#f0f6ff;border:1px solid #c8ddf8;border-radius:8px;padding:12px;">' +
								'<strong style="color:#1a3a5c;">Spectroscopy</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">Gamma spectroscopy records linked to measurement markers.</p>' +
							'</div>' +
							'<div style="background:#f0f6ff;border:1px solid #c8ddf8;border-radius:8px;padding:12px;">' +
								'<strong style="color:#1a3a5c;">Stats &amp; Reference</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">Aggregate stats, extreme readings, and dataset metadata.</p>' +
							'</div>' +
						'</div>' +
						'<p style="margin:12px 0 0;font-size:13px;color:#777;">All endpoints return JSON. Base path: <code style="background:#f0f6ff;padding:1px 5px;border-radius:4px;">/api</code>. No authentication required.</p>' +
					'</details>' +
					'<div style="display:flex;gap:10px;flex-wrap:wrap;align-items:center;">' +
						'<a href="/" style="display:inline-block;padding:7px 16px;background:#1a3a5c;color:#fff;border-radius:6px;font:600 13px/1.4 sans-serif;text-decoration:none;">\u2190 Back to Map</a>' +
						'<a href="' + mcpDocsURL + '" style="display:inline-block;padding:7px 16px;background:#0d9488;color:#fff;border-radius:6px;font:600 13px/1.4 sans-serif;text-decoration:none;">Switch to MCP API \u2192</a>' +
					'</div>' +
				'</div>';
				const swaggerUIEl = document.getElementById('swagger-ui');
				if (swaggerUIEl) {
					document.body.insertBefore(preamble, swaggerUIEl);
				} else {
					document.body.prepend(preamble);
				}
			}`, mcpDocsURL)
	http.Handle("/map-api/", httpSwagger.Handler(
		httpSwagger.URL("/map-api/doc.json"),
		httpSwagger.InstanceName("unifiedapi"),
		httpSwagger.UIConfig(map[string]string{
			"onComplete": mapAPINavScript,
		}),
	))

	// Serve MCP API spec JSON directly (avoids sharing swaggerFiles.Handler singleton
	// with the /map-api/ swagger instance, which causes a prefix conflict).
	http.HandleFunc("/mcp-api/doc.json", func(w http.ResponseWriter, r *http.Request) {
		doc, err := swag.ReadDoc("swagger")
		if err != nil {
			http.Error(w, "swagger spec unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(doc))
	})
	// Standalone /mcp-api/ page: custom HTML that loads assets from /map-api/
	// (single swagger static-file handler on this port — avoids prefix conflict).
	http.HandleFunc("/mcp-api/favicon.ico", serveFavicon)
	http.HandleFunc("/mcp-api/favicon-16x16.png", serveFavicon16)
	http.HandleFunc("/mcp-api/favicon-32x32.png", serveFavicon32)
	http.HandleFunc("/mcp-api/swagger-theme.css", serveSwaggerTheme)
	http.HandleFunc("/mcp-api/", serveMCPAPIPage)

	http.HandleFunc("/home", homeHandler)

	// Stories page and its data feed
	http.HandleFunc("/stories.html", func(w http.ResponseWriter, r *http.Request) {
		data, err := content.ReadFile("public_html/stories.html")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})
	http.HandleFunc("/data/stories.json", func(w http.ResponseWriter, r *http.Request) {
		data, err := content.ReadFile("public_html/data/stories.json")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	http.HandleFunc("/", mapHandler)

	// Register authentication routes if auth system is enabled
	if authManager != nil {
		// Serve profile page
		http.HandleFunc("/profile", authManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
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
		}))

		// Serve reset-password page
		http.HandleFunc("/reset-password", func(w http.ResponseWriter, r *http.Request) {
			data, err := content.ReadFile("public_html/reset-password.html")
			if err != nil {
				http.Error(w, "Reset password page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		})
	}

	// User admin routes (supports both URL password and session-based admin auth)
	if authManager != nil {
		// Helper function to check admin access (session-based or password-based)
		checkAdminAccess := func(w http.ResponseWriter, r *http.Request) bool {
			// First check for session-based admin auth
			if user, ok := auth.GetUserFromContext(r.Context()); ok && user.IsAdmin {
				return true
			}
			// Fall back to URL password
			if *adminPassword != "" {
				password := r.URL.Query().Get("password")
				if password == *adminPassword {
					return true
				}
			}
			http.Error(w, "Unauthorized - Please login as admin or provide password", http.StatusUnauthorized)
			return false
		}

		// Serve admin users page
		http.HandleFunc("/admin/users", authManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			// Prevent CloudFront from caching this dynamic admin page
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			if !checkAdminAccess(w, r) {
				return
			}
			data, err := content.ReadFile("public_html/admin-users.html")
			if err != nil {
				http.Error(w, "Admin page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}))

		// Serve admin uploads page (wrapper for /api/admin/uploads)
		http.HandleFunc("/admin/uploads", authManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			if !checkAdminAccess(w, r) {
				return
			}
			// Forward to the API endpoint which handles the uploads listing
			adminUploadsHandler(w, r)
		}))

		// Serve admin MCP analytics page
		http.HandleFunc("/admin/mcp", authManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			if !checkAdminAccess(w, r) {
				return
			}
			data, err := content.ReadFile("public_html/admin-mcp.html")
			if err != nil {
				http.Error(w, "Admin MCP page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}))

		// MCP analytics API endpoints
		// Note: /api/admin/mcp/data, /mcp/export, /mcp/delete are registered via httpapi.Register
		http.HandleFunc("/api/admin/mcp/update", authManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			if !checkAdminAccess(w, r) {
				return
			}
			adminMCPUpdateHandler(w, r)
		}))
		// Serve admin Realtime page and API endpoints
		http.HandleFunc("/admin/realtime", authManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")

			if !checkAdminAccess(w, r) {
				return
			}
			data, err := content.ReadFile("public_html/admin-realtime.html")
			if err != nil {
				http.Error(w, "Admin Realtime page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}))

		// Admin translations page and API
		http.HandleFunc("/admin/translations", authManager.OptionalAuth(func(w http.ResponseWriter, r *http.Request) {
			if !checkAdminAccess(w, r) {
				return
			}
			data, err := content.ReadFile("public_html/admin-translations.html")
			if err != nil {
				http.Error(w, "Page not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		}))
	}

	// Upload endpoint - protected with auth if required
	if *requireAuth && authManager != nil {
		http.HandleFunc("/upload", authManager.RequireAuth(uploadHandler))
	} else {
		http.HandleFunc("/upload", uploadHandler)
	}
	http.HandleFunc("/upload/progress", progressHandler)
	http.HandleFunc("/get_markers", getMarkersHandler)
	// Note: /stream_markers is Server-Sent Events (streaming) so gzip is skipped.
	// Gzip doesn't work well with streaming responses due to buffering.
	http.HandleFunc("/stream_markers", streamMarkersHandler)
	http.HandleFunc("/realtime_history", realtimeHistoryHandler)
	http.HandleFunc("/trackid/", trackHandler)
	http.HandleFunc("/tracks/", tracksHandler)
	// api/docs, licenses/, api/geoip, s/, api/spectrum/, api/markers/spectra, api/tracks/bounds, api/track-info/, api/update-coordinates, qrpng — registered via webServer.Register above
	// API endpoints ship JSON/archives. Keeping registration close to other
	// routes avoids surprises for operators scanning main() for handlers.
	limiter := httpapi.NewRateLimiter(time.Minute)
	apiHandler := httpapi.NewHandler(db, *dbType, archiveGen, limiter, log.Printf, archiveFrequency)

	// Keep MCP/realtime/translations admin APIs aligned with the legacy behavior:
	// they are registered only when auth is configured.
	var adminMCPDataAPIHandler http.HandlerFunc
	var adminMCPExportAPIHandler http.HandlerFunc
	var adminMCPDeleteAPIHandler http.HandlerFunc
	var adminRealtimeDataAPIHandler http.HandlerFunc
	var adminRealtimeExportAPIHandler http.HandlerFunc
	var adminRealtimeDeleteAPIHandler http.HandlerFunc
	var adminTranslationsReloadAPIHandler http.HandlerFunc
	var adminTranslationByIDAPIHandler http.HandlerFunc
	var adminTranslationsAPIHandler http.HandlerFunc
	if authManager != nil {
		adminMCPDataAPIHandler = adminMCPDataHandler
		adminMCPExportAPIHandler = adminMCPExportHandler
		adminMCPDeleteAPIHandler = adminMCPDeleteHandler
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
		AdminRealtimeDataHandler:         adminRealtimeDataAPIHandler,
		AdminRealtimeExportHandler:       adminRealtimeExportAPIHandler,
		AdminRealtimeDeleteHandler:       adminRealtimeDeleteAPIHandler,
		AdminTranslationsReloadHandler:   adminTranslationsReloadAPIHandler,
		AdminTranslationByIDHandler:      adminTranslationByIDAPIHandler,
		AdminTranslationsHandler:         adminTranslationsAPIHandler,
	})

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
