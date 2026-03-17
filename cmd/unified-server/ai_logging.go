package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

type aiLogEvent struct {
	UserID         string `json:"user_id"`
	UserEmail      string `json:"user_email"`
	SessionID      string `json:"session_id"`
	Timestamp      string `json:"timestamp"`
	ToolName       string `json:"tool_name"`
	GeneratedQuery string `json:"generated_query"`
	DurationMs     int64  `json:"duration_ms"`
	CommitHash     string `json:"commit_hash"`
	Error          string `json:"error"`
}

var (
	gitCommitOnce sync.Once
	gitCommitHash string
)

// logAISessionWithUser records a structured log entry for an AI tool execution (DuckDB-only).
// It is intentionally asynchronous and must not block the main tool handler.
func logAISessionWithUser(
	toolName string,
	query string,
	duration int64,
	err error,
	userID string,
	userEmail string,
) {

	go func() {

		event := aiLogEvent{

			UserID:    userID,
			UserEmail: userEmail,

			SessionID: newSessionID(),

			Timestamp: time.Now().
				UTC().
				Format(time.RFC3339),

			ToolName:       toolName,
			GeneratedQuery: sanitizeQuery(query),

			DurationMs: duration,

			CommitHash: getGitCommit(),

			Error: errString(err),
		}

		data, marshalErr := json.Marshal(event)

		if marshalErr != nil {

			log.Printf(
				"failed to marshal AI log event: %v",
				marshalErr,
			)

			return
		}

		log.Println(string(data))

		insertQueryLog(event)
	}()
}


// insertQueryLog writes one aiLogEvent to the DuckDB table mcp_ai_query_log using the shared duckDB connection.
// It is safe to call from the logging goroutine; errors are logged and never panic.
func insertQueryLog(event aiLogEvent) {

	if duckDB == nil {
		return
	}

	_, err := duckDB.Exec(`

		INSERT INTO mcp_ai_query_log (

			user_id,
			user_email,
			session_id,
			timestamp,
			tool_name,
			generated_query,
			duration_ms,
			commit_hash,
			error

		)

		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)

	`,
		event.UserID,
		event.UserEmail,
		event.SessionID,
		event.Timestamp,
		event.ToolName,
		event.GeneratedQuery,
		event.DurationMs,
		event.CommitHash,
		event.Error,
	)

	if err != nil {

		log.Printf(
			"failed to insert AI log event into DuckDB: %v",
			err,
		)
	}
}

// executeWithLogging logs at MCP runtime (tool execution) level.
// It must never panic and must not block tool execution (logging is asynchronous).
func executeWithLogging(
	toolName string,
	query string,
	fn func() (*sql.Rows, error),
) (*sql.Rows, error) {
	start := time.Now()
	rows, err := fn()
	duration := time.Since(start).Milliseconds()

	logAISessionWithUser(
		toolName,
		query,
		duration,
		err,
		"", // no user info available in this path
		"",
	)
		return rows, err
}

// getGitCommit returns the current git HEAD commit hash.
// It is cached after the first successful lookup.
func getGitCommit() string {
	gitCommitOnce.Do(func() {
		out, err := exec.Command("git", "rev-parse", "HEAD").Output()
		if err != nil {
			// Do not fail tool execution if git is unavailable.
			log.Printf("failed to read git commit hash: %v", err)
			return
		}
		gitCommitHash = strings.TrimSpace(string(out))
	})
	return gitCommitHash
}

var singleQuotedLiteral = regexp.MustCompile(`'[^']*'`)

// sanitizeQuery normalizes whitespace and scrubs obvious literal values.
// This helps avoid logging sensitive data while retaining query structure.
func sanitizeQuery(q string) string {
	if q == "" {
		return ""
	}

	// Collapse all whitespace (including newlines) to single spaces.
	q = strings.Join(strings.Fields(q), " ")

	// Replace string literals with a placeholder.
	q = singleQuotedLiteral.ReplaceAllString(q, "'?'")

	// Truncate overly long queries to avoid oversized log lines.
	const maxLen = 1000
	if len(q) > maxLen {
		q = q[:maxLen] + "..."
	}

	return q
}

// errString returns a safe string representation of an error.
func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// newSessionID generates a RFC4122-ish random UUID v4 string.
func newSessionID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// In the unlikely event of failure, return empty string rather than blocking.
		return ""
	}

	// Set version (4) and variant bits.
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return formatUUID(b)
}

func formatUUID(b [16]byte) string {
	return strings.ToLower(
		strings.Join([]string{
			sprintf8(b[0:4]),
			sprintf4(b[4:6]),
			sprintf4(b[6:8]),
			sprintf4(b[8:10]),
			sprintf8(b[10:14]) + sprintf2(b[14:16]),
		}, "-"),
	)
}

func sprintf2(b []byte) string {
	return sprintfN(b, 2)
}

func sprintf4(b []byte) string {
	return sprintfN(b, 4)
}

func sprintf8(b []byte) string {
	return sprintfN(b, 4) // 4 bytes → 8 hex chars
}

func sprintfN(b []byte, byteCount int) string {
	if len(b) < byteCount {
		byteCount = len(b)
	}
	val := uint64(0)
	for i := 0; i < byteCount; i++ {
		val = (val << 8) | uint64(b[i])
	}
	width := byteCount * 2
	buf := make([]byte, width)
	for i := width - 1; i >= 0; i-- {
		buf[i] = "0123456789abcdef"[val&0xf]
		val >>= 4
	}
	return string(buf)
}

