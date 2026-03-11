package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	_ "github.com/marcboeker/go-duckdb"
)

var duckDB *sql.DB

func initDuckDB() error {

	// 1. Resolve DuckDB path safely
	duckPath := os.Getenv("DUCKDB_PATH")
	if duckPath == "" {
		duckPath = "./analytics.duckdb"
	}

	dsn := duckPath + "?access_mode=READ_WRITE"

	var err error
	duckDB, err = sql.Open("duckdb", dsn)
	if err != nil {
		return fmt.Errorf("failed to open duckdb: %w", err)
	}

	// 2. Production-safe connection pool config
	duckDB.SetMaxOpenConns(1)  // DuckDB works best with single writer
	duckDB.SetMaxIdleConns(1)
	duckDB.SetConnMaxLifetime(0)

	if err := duckDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping duckdb: %w", err)
	}

	log.Printf("DuckDB initialized at %s", duckPath)

	// 3. Enable WAL checkpointing for durability
	duckDB.Exec("PRAGMA wal_autocheckpoint=1000;")

	// 4. Load postgres extension safely (non-fatal)
	if _, err := duckDB.Exec("INSTALL postgres;"); err != nil {
		log.Printf("Warning: postgres extension install failed: %v", err)
	}
	if _, err := duckDB.Exec("LOAD postgres;"); err != nil {
		log.Printf("Warning: postgres extension load failed: %v", err)
	}

	// 5. Attach Postgres if configured
	pgURL := os.Getenv("DATABASE_URL")
	if pgURL != "" {
		query := fmt.Sprintf(
			"ATTACH '%s' AS postgres_db (TYPE POSTGRES, READ_ONLY)",
			pgURL,
		)

		if _, err := duckDB.Exec(query); err != nil {
			log.Printf("Warning: failed to attach postgres: %v", err)
		} else {
			log.Println("PostgreSQL attached as postgres_db")
		}
	}

	// Create sequence first, then the audit log table that references it
	createSchemaQuery := `
	CREATE SEQUENCE IF NOT EXISTS seq_query_log;
	CREATE TABLE IF NOT EXISTS mcp_query_log (
		id            BIGINT DEFAULT nextval('seq_query_log'),
		tool_name     VARCHAR,
		params        JSON,
		result_count  INTEGER,
		duration_ms   DOUBLE,
		client_info   VARCHAR,
		created_at    TIMESTAMPTZ DEFAULT now()
	);
	`
	if _, err := duckDB.Exec(createSchemaQuery); err != nil {
		return fmt.Errorf("failed to create audit log schema: %w", err)
	}

	// 6. Schema version table
	_, err = duckDB.Exec(`
	CREATE TABLE IF NOT EXISTS schema_version (
		version INTEGER
	);
	`)
	if err != nil {
		return err
	}

	var version int
	err = duckDB.QueryRow(`
	SELECT COALESCE(MAX(version),0) FROM schema_version
	`).Scan(&version)

	if err != nil {
		return err
	}

	// 7. Migration to version 2 (adds user info)
	if version < 2 {

		log.Println("Running schema migration to v2")

		_, err = duckDB.Exec(`
		CREATE SEQUENCE IF NOT EXISTS seq_query_log;

		CREATE TABLE IF NOT EXISTS mcp_query_log (
			id BIGINT DEFAULT nextval('seq_query_log'),
			tool_name VARCHAR,
			params JSON,
			result_count INTEGER,
			duration_ms DOUBLE,
			client_info VARCHAR,
			created_at TIMESTAMPTZ DEFAULT now()
		);
		`)
		if err != nil {
			return err
		}

		// AI log table with user info
		_, err = duckDB.Exec(`
		CREATE TABLE IF NOT EXISTS mcp_ai_query_log (
			user_id TEXT,
			user_email TEXT,
			session_id TEXT,
			timestamp TIMESTAMP,
			tool_name TEXT,
			generated_query TEXT,
			duration_ms BIGINT,
			commit_hash TEXT,
			error TEXT
		);
		`)
		if err != nil {
			return err
		}

        indexes := []string{

            `CREATE INDEX IF NOT EXISTS idx_ai_timestamp
             ON mcp_ai_query_log(timestamp);`,
        
            `CREATE INDEX IF NOT EXISTS idx_ai_user
             ON mcp_ai_query_log(user_id);`,
        
            `CREATE INDEX IF NOT EXISTS idx_ai_user_email
             ON mcp_ai_query_log(user_email);`,
        
            `CREATE INDEX IF NOT EXISTS idx_ai_tool
             ON mcp_ai_query_log(tool_name);`,
        }        

		for _, idx := range indexes {
			duckDB.Exec(idx)
		}

		_, err = duckDB.Exec(`
		DELETE FROM schema_version;

        INSERT INTO schema_version(version) VALUES (2);

		`)
		if err != nil {
			return err
		}
	}

	log.Println("DuckDB schema ready")

	return nil
}

// LogQueryAsync logs a tool execution to DuckDB asynchronously.
func LogQueryAsync(toolName string, params map[string]any, resultCount int, duration time.Duration, clientInfo string) {
    if duckDB == nil {
        return
    }
    
    go func() {
        // Serialize params as proper JSON for DuckDB's JSON column type.
        paramsJSON, err := json.Marshal(params)
        if err != nil {
            log.Printf("Error marshaling params to JSON: %v", err)
            return
        }
        paramsStr := string(paramsJSON)

        _, execErr := duckDB.Exec(`
            INSERT INTO mcp_query_log (tool_name, params, result_count, duration_ms, client_info)
            VALUES (?, ?, ?, ?, ?)
        `, toolName, paramsStr, resultCount, float64(duration.Milliseconds()), clientInfo)

        if execErr != nil {
            log.Printf("Error logging query to DuckDB: %v", execErr)
        }
    }()
}

// Analytics Functions

// GetToolUsageStats returns usage statistics for tools.
func GetToolUsageStats() ([]map[string]any, error) {
    if duckDB == nil {
        return nil, fmt.Errorf("duckdb not initialized")
    }
    
    rows, err := duckDB.Query(`
        SELECT tool_name, COUNT(*) as calls, AVG(duration_ms) as avg_duration, MAX(duration_ms) as max_duration
        FROM mcp_query_log
        GROUP BY tool_name
        ORDER BY calls DESC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []map[string]any
    for rows.Next() {
        var toolName string
        var calls int64
        var avgDur, maxDur float64
        if err := rows.Scan(&toolName, &calls, &avgDur, &maxDur); err != nil {
            return nil, err
        }
        stats = append(stats, map[string]any{
            "tool": toolName,
            "calls": calls,
            "avg_ms": avgDur,
            "max_ms": maxDur,
        })
    }
    return stats, nil
}

// QueryPostgresAnalytics executes an arbitrary analytical query on the attached Postgres DB.
// This is the powerful "FAQ" enabler.
// WARNING: Logic constraints should be applied in a real production environment.
func QueryPostgresAnalytics(query string, args ...any) ([]map[string]any, error) {
    if duckDB == nil {
        return nil, fmt.Errorf("duckdb not initialized")
    }
    
    // We execute the query directly against DuckDB, which can reference postgres_db.tables
    rows, err := duckDB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Generic map scanning
    cols, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    
    var results []map[string]any
    for rows.Next() {
        // Create a slice of interface{} to hold values
        columns := make([]interface{}, len(cols))
        columnPointers := make([]interface{}, len(cols))
        for i := range columns {
            columnPointers[i] = &columns[i]
        }

        if err := rows.Scan(columnPointers...); err != nil {
            return nil, err
        }

        row := make(map[string]any)
        for i, colName := range cols {
            val := columns[i]
            // DuckDB driver might return specific types, handle basic conversion if needed
            // For now, pass through
            if b, ok := val.([]byte); ok {
                row[colName] = string(b)
            } else {
                row[colName] = val
            }
        }
        results = append(results, row)
    }
    return results, nil
}
