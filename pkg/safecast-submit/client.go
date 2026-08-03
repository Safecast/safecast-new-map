// Package safecastsubmit submits bGeigie log uploads to api.safecast.org on
// behalf of a user, using that user's own Safecast API key.
package safecastsubmit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

// Client handles outbound submission calls to api.safecast.org.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a client pointed at the real api.safecast.org.
func NewClient() *Client {
	return newClientWithBaseURL("https://api.safecast.org")
}

// NewClientWithBaseURL creates a client pointed at an arbitrary base URL, for tests.
func NewClientWithBaseURL(baseURL string) *Client {
	return newClientWithBaseURL(baseURL)
}

func newClientWithBaseURL(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// meResponse mirrors the relevant fields of UsersController#me's JSON response.
type meResponse struct {
	ID int64 `json:"id"`
}

// ResolveUserID looks up the numeric Safecast user id for the given API key via
// GET /users/me.json?api_key=<key>.
func (c *Client) ResolveUserID(ctx context.Context, apiKey string) (string, error) {
	u := fmt.Sprintf("%s/users/me.json?api_key=%s", c.baseURL, url.QueryEscape(apiKey))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var me meResponse
	if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if me.ID == 0 {
		return "", fmt.Errorf("resolve user id: empty id in response")
	}

	return fmt.Sprintf("%d", me.ID), nil
}

// bgeigieImportSummary is the subset of a bgeigie_imports.json list entry we need.
type bgeigieImportSummary struct {
	ID     int64  `json:"id"`
	Source string `json:"source"`
}

// CheckExists reports whether the given user already has a bgeigie import whose
// stored filename matches, via GET /bgeigie_imports.json?by_user_id=<id>&q=<filename>.
func (c *Client) CheckExists(ctx context.Context, apiKey, safecastUserID, filename string) (bool, error) {
	u, err := url.Parse(fmt.Sprintf("%s/bgeigie_imports.json", c.baseURL))
	if err != nil {
		return false, fmt.Errorf("parse base URL: %w", err)
	}
	query := u.Query()
	query.Set("api_key", apiKey)
	query.Set("by_user_id", safecastUserID)
	query.Set("q", filename)
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return false, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var imports []bgeigieImportSummary
	if err := json.NewDecoder(resp.Body).Decode(&imports); err != nil {
		return false, fmt.Errorf("decode response: %w", err)
	}

	for _, imp := range imports {
		if imp.Source == filename {
			return true, nil
		}
	}
	return false, nil
}

// submitResponse mirrors the relevant fields of BgeigieImportsController#create's
// JSON response.
type submitResponse struct {
	ID int64 `json:"id"`
}

// Submit uploads a bGeigie log file to api.safecast.org via
// POST /bgeigie_imports.json?api_key=<key>, multipart/form-data with
// bgeigie_import[description] and bgeigie_import[source] parts, matching the
// bGeigieZen firmware's upload_detail() implementation.
func (c *Client) Submit(ctx context.Context, apiKey, filename string, content []byte) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("bgeigie_import[description]", "Uploaded from simplemap"); err != nil {
		return "", fmt.Errorf("write description field: %w", err)
	}

	part, err := writer.CreateFormFile("bgeigie_import[source]", filename)
	if err != nil {
		return "", fmt.Errorf("create file part: %w", err)
	}
	if _, err := part.Write(content); err != nil {
		return "", fmt.Errorf("write file part: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}

	u := fmt.Sprintf("%s/bgeigie_imports.json?api_key=%s", c.baseURL, url.QueryEscape(apiKey))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, &body)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result submitResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if result.ID == 0 {
		return "", fmt.Errorf("submit: empty id in response")
	}

	return fmt.Sprintf("%d", result.ID), nil
}
