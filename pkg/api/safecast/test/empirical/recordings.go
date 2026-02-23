// Package empirical loads API recordings from embedded testdata (or an optional
// directory) and provides types and helpers for empirical response-format tests.
// It does not depend on the parent safecast_test package.
package empirical

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

//go:embed testdata/*.json
var testdataFS embed.FS

// Recording is the top-level structure of each JSON file in api_recordings/.
// Name is set by LoadRecordings to the base filename without .json (e.g. get_users).
type Recording struct {
	RecordedAt string           `json:"recorded_at"`
	Request    RecordedRequest  `json:"request"`
	Response   RecordedResponse `json:"response"`
	Name       string           // set by LoadRecordings, not from JSON
}

// RecordedRequest holds the request fields from a recording.
type RecordedRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// RecordedResponse holds the response fields from a recording.
type RecordedResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// LoadRecordings reads all *.json files from dir and returns one Recording per
// file. Name is set to the filename without .json. Non-JSON files or invalid
// JSON are skipped (no error). Returns nil slice if dir cannot be read.
func LoadRecordings(dir string) ([]Recording, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []Recording
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".json") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var rec Recording
		if err := json.Unmarshal(data, &rec); err != nil {
			continue
		}
		rec.Name = strings.TrimSuffix(e.Name(), ".json")
		out = append(out, rec)
	}
	return out, nil
}

// LoadRecordingsFromFS reads all *.json files from the given fs.FS (e.g. embedded
// testdata). Names with a directory prefix (e.g. "testdata/get_users.json") are
// stripped to the base name without .json for Recording.Name.
func LoadRecordingsFromFS(fsys fs.FS) ([]Recording, error) {
	entries, err := fs.Glob(fsys, "testdata/*.json")
	if err != nil {
		return nil, err
	}
	var out []Recording
	for _, path := range entries {
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			continue
		}
		var rec Recording
		if err := json.Unmarshal(data, &rec); err != nil {
			continue
		}
		name := path
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		rec.Name = strings.TrimSuffix(name, ".json")
		out = append(out, rec)
	}
	return out, nil
}

// PathFromRecordedURL returns the request path (and raw query) to use when
// replaying the recording. The new API exposes root only at /api/v1, so
// path "/" is mapped to "/api/v1". Other paths are used as-is (path + query).
func PathFromRecordedURL(rawURL string) (path, query string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}
	path = u.Path
	if path == "" {
		path = "/"
	}
	if path == "/" {
		path = "/api/v1"
	}
	return path, u.RawQuery, nil
}
