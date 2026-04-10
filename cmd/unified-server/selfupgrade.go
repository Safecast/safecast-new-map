package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"safecast-new-map/pkg/database"
	"safecast-new-map/pkg/selfupgrade"
)

// selfUpgradeFlagSet keeps the flag pointers optional so unsupported platforms
// never register a -selfupgrade flag, following the "Make the zero value useful"
// proverb. We also carry the default URLs so runtime decisions can fall back to
// platform-specific release assets without extra state.
type selfUpgradeFlagSet struct {
	enabled    *bool
	url        *string
	supported  bool
	defaultURL string
	duckDBURL  string
}

var selfUpgradeFlags = registerSelfUpgradeFlags()

// selfUpgradeDatabaseInfo recreates the DSN resolution logic so the deployment
// manager knows how to back up the live database. We intentionally mirror the
// database package defaults instead of importing internal helpers to avoid
// circular dependencies.
func selfUpgradeDatabaseInfo(cfg database.Config) (driver, dsn string) {
	driver = strings.ToLower(strings.TrimSpace(cfg.DBType))
	switch driver {
	case "sqlite", "chai":
		dsn = strings.TrimSpace(cfg.DBPath)
		if dsn == "" {
			dsn = fmt.Sprintf("database-%d.%s", cfg.Port, driver)
		}
	case "duckdb":
		dsn = strings.TrimSpace(cfg.DBPath)
		if dsn == "" {
			dsn = fmt.Sprintf("database-%d.duckdb", cfg.Port)
		}
	case "pgx":
		if strings.TrimSpace(cfg.DBConn) != "" {
			dsn = cfg.DBConn
		} else {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
				cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.PGSSLMode)
		}
	case "clickhouse":
		dsn = database.ClickHouseDSNFromConfig(cfg)
	}
	return driver, dsn
}

// startSelfUpgrade wires up the background deployment manager when operators
// pass -selfupgrade. We centralise the configuration assembly to keep main()
// readable while still reusing the same channels and goroutines that make the
// manager safe without mutexes, following "Don't communicate by sharing
// memory".
func startSelfUpgrade(ctx context.Context, dbCfg database.Config) context.CancelFunc {
	if ctx == nil {
		ctx = context.Background()
	}
	if !selfUpgradeFlags.supported {
		// Unsupported platforms skip registration entirely, staying silent because no
		// CLI flag is present and there is nothing to configure for operators.
		return nil
	}
	if selfUpgradeFlags.enabled == nil || !*selfUpgradeFlags.enabled {
		// Operators opted out, so we avoid spawning background work.
		return nil
	}

	exePath, err := os.Executable()
	if err != nil {
		// Without a deterministic binary path we cannot build rollback bundles, so we bail out early.
		log.Printf("selfupgrade disabled: cannot resolve executable path: %v", err)
		return nil
	}

	driverName, dsn := selfUpgradeDatabaseInfo(dbCfg)

	downloadURL := ""
	if selfUpgradeFlags.url != nil {
		downloadURL = strings.TrimSpace(*selfUpgradeFlags.url)
	}
	if downloadURL == "" {
		// We default to platform-specific artefacts and swap in DuckDB builds when requested.
		downloadURL = strings.TrimSpace(selfUpgradeFlags.defaultURL)
		if driverName == "duckdb" && strings.TrimSpace(selfUpgradeFlags.duckDBURL) != "" {
			downloadURL = strings.TrimSpace(selfUpgradeFlags.duckDBURL)
		}
	}
	if downloadURL == "" {
		log.Printf("selfupgrade disabled: download URL not configured for %s/%s", runtime.GOOS, runtime.GOARCH)
		return nil
	}

	const canaryPort = 9876
	workspace := filepath.Join(filepath.Dir(exePath), "selfupgrade-cache")
	backupsDir := filepath.Join(workspace, "db_backups")

	var dbController selfupgrade.DatabaseController
	switch driverName {
	case "sqlite", "chai", "duckdb":
		if dsn != "" {
			dbController = &selfupgrade.FileDatabaseController{
				Driver:       driverName,
				OriginalPath: dsn,
				BackupsDir:   backupsDir,
				Logf:         log.Printf,
			}
		}
	default:
		log.Printf("selfupgrade database backups not configured for driver %s", driverName)
	}

	cfgAuto := selfupgrade.Config{
		DownloadURL:    downloadURL,
		CurrentVersion: CompileVersion,
		BinaryPath:     exePath,
		DeployDir:      workspace,
		DBBackupsDir:   backupsDir,
		CanaryPort:     canaryPort,
		Logf:           log.Printf,
		Database:       dbController,
	}

	manager, err := selfupgrade.NewManager(cfgAuto)
	if err != nil {
		log.Printf("selfupgrade disabled: %v", err)
		return nil
	}

	ctxDeploy, cancel := context.WithCancel(ctx)
	if err := manager.Start(ctxDeploy); err != nil {
		log.Printf("selfupgrade start failed: %v", err)
		cancel()
		return nil
	}

	http.Handle("/selfupgrade/", manager.HTTPHandler())
	go func() {
		manager.Wait()
	}()
	log.Printf("selfupgrade manager polling %s", downloadURL)

	return cancel
}

// selfupgradeStartupDelay pauses the freshly spawned binary during a handoff so
// the previous process can close network listeners without racing. The old
// binary encodes the delay inside SELFUPGRADE_WAIT_SECONDS.
func selfupgradeStartupDelay(logf func(string, ...any)) {
	waitEnv := strings.TrimSpace(os.Getenv("SELFUPGRADE_WAIT_SECONDS"))
	if waitEnv == "" {
		return
	}
	defer os.Unsetenv("SELFUPGRADE_WAIT_SECONDS")
	seconds, err := strconv.ParseFloat(waitEnv, 64)
	if err != nil || seconds <= 0 {
		if logf != nil && err != nil {
			logf("selfupgrade: invalid wait seconds %q: %v", waitEnv, err)
		}
		return
	}
	delay := time.Duration(seconds * float64(time.Second))
	if logf != nil {
		logf("selfupgrade: waiting %s for predecessor shutdown", delay)
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	<-timer.C
}

func selfupgradeCleanEnv(env []string) []string {
	cleaned := make([]string, 0, len(env))
	for _, kv := range env {
		if strings.HasPrefix(kv, "SELFUPGRADE_") {
			continue
		}
		cleaned = append(cleaned, kv)
	}
	return cleaned
}

// selfupgradeRollback revives the last known good binary when the replacement
// process fails to bind its listeners. We relaunch the previous binary and let
// it resume service.
func selfupgradeRollback(logf func(string, ...any)) bool {
	lastGood := strings.TrimSpace(os.Getenv("SELFUPGRADE_LAST_GOOD"))
	if lastGood == "" {
		return false
	}
	exe, err := os.Executable()
	if err != nil {
		if logf != nil {
			logf("selfupgrade rollback skipped: executable unknown: %v", err)
		}
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := selfupgrade.RestoreBinary(ctx, lastGood, exe); err != nil {
		if logf != nil {
			logf("selfupgrade rollback restore failed: %v", err)
		}
		return false
	}
	env := selfupgradeCleanEnv(os.Environ())
	if logf != nil {
		logf("selfupgrade rollback: relaunching %s", lastGood)
	}
	if runtime.GOOS != "windows" {
		if err := syscall.Exec(exe, os.Args, env); err != nil {
			if logf != nil {
				logf("selfupgrade rollback exec failed: %v", err)
			}
			return false
		}
		return true
	}
	cmd := exec.Command(exe, os.Args[1:]...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		if logf != nil {
			logf("selfupgrade rollback spawn failed: %v", err)
		}
		return false
	}
	return true
}

// selfupgradeHandleServerError unifies server startup failures so the
// self-upgrade workflow can roll back to the previous binary when binding the
// listener fails.
func selfupgradeHandleServerError(err error, logf func(string, ...any)) {
	if err == nil {
		return
	}
	if errors.Is(err, http.ErrServerClosed) {
		return
	}
	if logf != nil {
		logf("HTTP server error: %v", err)
	}
	if selfupgradeRollback(logf) {
		os.Exit(0)
	}
}

// registerSelfUpgradeFlags inspects the compilation target and only registers
// self-upgrade flags for platforms where we publish release binaries. This keeps
// unsupported builds free of dead flags while still allowing callers to override
// the download location when needed.
func registerSelfUpgradeFlags() selfUpgradeFlagSet {
	const base = "https://github.com/Safecast/safecast-new-map/releases/download/latest/"

	switch runtime.GOOS {
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			return newSelfUpgradeFlagSet(base, "darwin/amd64", "safecast-new-map_darwin_amd64", "safecast-new-map_darwin_amd64_duckdb")
		case "arm64":
			return newSelfUpgradeFlagSet(base, "darwin/arm64", "safecast-new-map_darwin_arm64", "safecast-new-map_darwin_arm64_duckdb")
		}
	case "freebsd":
		if runtime.GOARCH == "amd64" {
			return newSelfUpgradeFlagSet(base, "freebsd/amd64", "safecast-new-map_freebsd_amd64", "")
		}
	case "linux":
		switch runtime.GOARCH {
		case "386":
			return newSelfUpgradeFlagSet(base, "linux/386", "safecast-new-map_linux_386", "")
		case "amd64":
			return newSelfUpgradeFlagSet(base, "linux/amd64", "safecast-new-map_linux_amd64", "")
		case "arm64":
			return newSelfUpgradeFlagSet(base, "linux/arm64", "safecast-new-map_linux_arm64", "")
		}
	case "netbsd":
		if runtime.GOARCH == "amd64" {
			return newSelfUpgradeFlagSet(base, "netbsd/amd64", "safecast-new-map_netbsd_amd64", "")
		}
	case "openbsd":
		if runtime.GOARCH == "amd64" {
			return newSelfUpgradeFlagSet(base, "openbsd/amd64", "safecast-new-map_openbsd_amd64", "")
		}
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			return newSelfUpgradeFlagSet(base, "windows/amd64", "safecast-new-map_windows_amd64.exe", "")
		case "arm64":
			return newSelfUpgradeFlagSet(base, "windows/arm64", "safecast-new-map_windows_arm64.exe", "")
		}
	}

	return selfUpgradeFlagSet{supported: false}
}

// newSelfUpgradeFlagSet keeps the flag wiring compact while documenting the
// release artefact associated with each platform. We bake descriptions into the
// help text to reduce guesswork for operators invoking `-help`.
func newSelfUpgradeFlagSet(base, platform, asset, duckAsset string) selfUpgradeFlagSet {
	defaultURL := base + asset
	duckURL := ""
	if strings.TrimSpace(duckAsset) != "" {
		duckURL = base + duckAsset
	}

	enabled := flag.Bool("selfupgrade", false, fmt.Sprintf("Enable the background auto-deployment manager on %s hosts", platform))
	url := flag.String("selfupgrade-url", defaultURL, fmt.Sprintf("Direct download URL for the %s binary", platform))

	return selfUpgradeFlagSet{
		enabled:    enabled,
		url:        url,
		supported:  true,
		defaultURL: defaultURL,
		duckDBURL:  duckURL,
	}
}
