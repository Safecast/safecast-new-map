package database

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Config holds the configuration details for initializing the database.
type Config struct {
	DBType      string // The type of the database driver (e.g., "sqlite", "chai", or "pgx" (PostgreSQL))
	DBPath      string // The file path to the database file (for file-based databases)
	DBConn      string // Raw DSN for network drivers (pgx or clickhouse)
	DBHost      string // The host for PostgreSQL
	DBPort      int    // The port for PostgreSQL
	DBUser      string // The user for PostgreSQL
	DBPass      string // The password for PostgreSQL
	DBName      string // The name of the PostgreSQL database
	PGSSLMode   string // The SSL mode for PostgreSQL
	ClickSecure bool   // Enable TLS when connecting to ClickHouse over HTTP transport
	Port        int    // The port number (used in database file naming if needed)
}

// normalizeDBType trims and lowercases driver names and maps common aliases.
func normalizeDBType(dbType string) string {
	cleaned := strings.ToLower(strings.TrimSpace(dbType))
	switch cleaned {
	case "postgres", "postgresql", "pq", "postgres+psql", "postgresql+psql":
		return "pgx"
	default:
		return cleaned
	}
}

// ClickHouseDSNFromConfig assembles a DSN understood by the lightweight HTTP driver.
func ClickHouseDSNFromConfig(cfg Config) string {
	if trimmed := strings.TrimSpace(cfg.DBConn); trimmed != "" {
		return trimmed
	}

	host := strings.TrimSpace(cfg.DBHost)
	if host == "" {
		host = "127.0.0.1"
	}
	if _, _, err := net.SplitHostPort(host); err != nil {
		port := cfg.DBPort
		if port <= 0 {
			port = 9000
		}
		host = net.JoinHostPort(host, strconv.Itoa(port))
	}

	user := strings.TrimSpace(cfg.DBUser)
	pass := cfg.DBPass
	name := strings.Trim(strings.TrimSpace(cfg.DBName), "/")

	dsn := url.URL{Scheme: "clickhouse", Host: host}
	if user != "" {
		if strings.TrimSpace(pass) != "" {
			dsn.User = url.UserPassword(user, pass)
		} else {
			dsn.User = url.User(user)
		}
	}
	if name != "" {
		dsn.Path = "/" + name
	}

	params := url.Values{}
	if cfg.ClickSecure {
		params.Set("secure", "true")
	}
	dsn.RawQuery = params.Encode()
	return dsn.String()
}

// startIDGenerator launches a goroutine for generating unique IDs.
func startIDGenerator(initialID int64) chan int64 {
	idChannel := make(chan int64)
	go func(start int64) {
		currentID := start
		for {
			idChannel <- currentID
			currentID++
		}
	}(initialID)
	return idChannel
}
