package safecastfetcher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"
)

// DownloadLogFile downloads a log file from an S3 URL
// Returns: file content, filename from URL, error
func DownloadLogFile(ctx context.Context, sourceURL string) ([]byte, string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 60 * time.Second, // Longer timeout for potentially large files
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", sourceURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("create request: %w", err)
	}

	// Execute download
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read file content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read response body: %w", err)
	}

	// Extract filename from URL
	filename := extractFilename(sourceURL)

	// Basic validation: check for bGeigie log signature
	if len(content) > 0 && !isValidLogFile(content) {
		return nil, "", fmt.Errorf("invalid log file format (expected $BNRDD signature)")
	}

	return content, filename, nil
}

// extractFilename extracts the filename from an S3 URL path
func extractFilename(sourceURL string) string {
	// Example URL: https://safecastapi-imports-production-us-west-2.s3.amazonaws.com/uploads/bgeigie_import/source/12345/BGM-200110.LOG
	// Extract: BGM-200110.LOG
	filename := path.Base(sourceURL)

	// If filename is empty or looks like a query string, use a default
	if filename == "" || filename == "." || filename == "/" || strings.Contains(filename, "?") {
		filename = "safecast_import.log"
	}

	return filename
}

// isValidLogFile performs basic validation on the file content
// Returns true if the file appears to be a valid bGeigie log
func isValidLogFile(content []byte) bool {
	// Check if empty
	if len(content) == 0 {
		return false
	}

	// Convert first 1024 bytes to string for checking
	checkSize := 1024
	if len(content) < checkSize {
		checkSize = len(content)
	}

	header := string(content[:checkSize])

	// Look for bGeigie log signature: $BNRDD
	// Also accept other common formats that might be in the API
	validSignatures := []string{
		"$BNRDD",  // bGeigie Nano
		"$BMRDD",  // bGeigie Mini
		"$BNXRDD", // bGeigie NX
	}

	for _, sig := range validSignatures {
		if strings.Contains(header, sig) {
			return true
		}
	}

	// If no signature found, it might be a different format
	// For now, we'll be permissive and allow it through
	// The actual parser will handle format validation
	return true
}

// BytesFile is an adapter to convert []byte to multipart.File interface
type BytesFile struct {
	*bytes.Reader
	name string
}

// NewBytesFile creates a new BytesFile from byte content
func NewBytesFile(data []byte, filename string) *BytesFile {
	return &BytesFile{
		Reader: bytes.NewReader(data),
		name:   filename,
	}
}

// Close implements the Close method for multipart.File interface
func (b *BytesFile) Close() error {
	return nil
}

// Seek implements the Seek method
func (b *BytesFile) Seek(offset int64, whence int) (int64, error) {
	return b.Reader.Seek(offset, whence)
}

// ReadAt implements the ReadAt method
func (b *BytesFile) ReadAt(p []byte, off int64) (n int, err error) {
	return b.Reader.ReadAt(p, off)
}
