package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/database"
	_ "safecast-new-map/pkg/database/drivers"
)

// fakeSafecastSubmitter is a test double for safecastSubmitter.
type fakeSafecastSubmitter struct {
	resolveUserID string
	resolveErr    error
	existsResult  bool
	existsErr     error
	submitID      string
	submitErr     error
	checkCalled   bool
	submitCalled  bool
	gotAPIKey     string
	gotUserID     string
	gotFilename   string
}

func (f *fakeSafecastSubmitter) ResolveUserID(ctx context.Context, apiKey string) (string, error) {
	return f.resolveUserID, f.resolveErr
}

func (f *fakeSafecastSubmitter) CheckExists(ctx context.Context, apiKey, safecastUserID, filename string) (bool, error) {
	f.checkCalled = true
	f.gotAPIKey = apiKey
	f.gotUserID = safecastUserID
	f.gotFilename = filename
	return f.existsResult, f.existsErr
}

func (f *fakeSafecastSubmitter) Submit(ctx context.Context, apiKey, filename string, content []byte) (string, error) {
	f.submitCalled = true
	return f.submitID, f.submitErr
}

func newSafecastSubmitTestDB(t *testing.T) *database.Database {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	cfg := database.Config{DBType: "sqlite", DBPath: dbPath, Port: 1}
	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase: %v", err)
	}
	if err := db.InitSchema(cfg); err != nil {
		t.Fatalf("InitSchema: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func insertTestUpload(t *testing.T, db *database.Database, trackID, filename string) int64 {
	t.Helper()
	id, err := db.InsertUpload(context.Background(), database.Upload{
		Filename: filename,
		TrackID:  trackID,
		Source:   "user-upload",
	})
	if err != nil {
		t.Fatalf("InsertUpload: %v", err)
	}
	return id
}

func getSafecastStatus(t *testing.T, db *database.Database, uploadID int64) (importID, submitErr string) {
	t.Helper()
	var i, e *string
	err := db.DB.QueryRowContext(context.Background(),
		`SELECT safecast_import_id, safecast_submit_error FROM uploads WHERE id = ?`, uploadID).Scan(&i, &e)
	if err != nil {
		t.Fatalf("query upload status: %v", err)
	}
	if i != nil {
		importID = *i
	}
	if e != nil {
		submitErr = *e
	}
	return
}

func TestSubmitToSafecastIfNeeded_SkipsWhenNoCredentials(t *testing.T) {
	db := newSafecastSubmitTestDB(t)
	fake := &fakeSafecastSubmitter{}
	orig := safecastSubmitClient
	safecastSubmitClient = fake
	defer func() { safecastSubmitClient = orig }()

	uploadID := insertTestUpload(t, db, "track-1", "log1.log")
	submitToSafecastIfNeeded("track-1", uploadID, "log1.log", []byte("data"), "", "", db)

	if fake.checkCalled || fake.submitCalled {
		t.Fatal("expected no upstream calls when credentials are empty")
	}
}

func TestSubmitToSafecastIfNeeded_SkipsWhenAlreadyExists(t *testing.T) {
	db := newSafecastSubmitTestDB(t)
	fake := &fakeSafecastSubmitter{existsResult: true}
	orig := safecastSubmitClient
	safecastSubmitClient = fake
	defer func() { safecastSubmitClient = orig }()

	uploadID := insertTestUpload(t, db, "track-2", "log2.log")
	submitToSafecastIfNeeded("track-2", uploadID, "log2.log", []byte("data"), "key", "7", db)

	if !fake.checkCalled {
		t.Fatal("expected CheckExists to be called")
	}
	if fake.submitCalled {
		t.Fatal("expected Submit not to be called when already exists")
	}
	importID, _ := getSafecastStatus(t, db, uploadID)
	if importID != "" {
		t.Errorf("expected no import id recorded, got %q", importID)
	}
}

func TestSubmitToSafecastIfNeeded_SubmitsAndRecordsSuccess(t *testing.T) {
	db := newSafecastSubmitTestDB(t)
	fake := &fakeSafecastSubmitter{existsResult: false, submitID: "999"}
	orig := safecastSubmitClient
	safecastSubmitClient = fake
	defer func() { safecastSubmitClient = orig }()

	uploadID := insertTestUpload(t, db, "track-3", "log3.log")
	submitToSafecastIfNeeded("track-3", uploadID, "log3.log", []byte("data"), "key", "7", db)

	if !fake.submitCalled {
		t.Fatal("expected Submit to be called")
	}
	if fake.gotAPIKey != "key" || fake.gotUserID != "7" || fake.gotFilename != "log3.log" {
		t.Errorf("unexpected CheckExists args: apiKey=%q userID=%q filename=%q", fake.gotAPIKey, fake.gotUserID, fake.gotFilename)
	}
	importID, submitErr := getSafecastStatus(t, db, uploadID)
	if importID != "999" {
		t.Errorf("got importID %q, want %q", importID, "999")
	}
	if submitErr != "" {
		t.Errorf("expected no error recorded, got %q", submitErr)
	}
}

func TestSubmitToSafecastIfNeeded_RecordsSubmitFailure(t *testing.T) {
	db := newSafecastSubmitTestDB(t)
	fake := &fakeSafecastSubmitter{existsResult: false, submitErr: errors.New("upstream rejected")}
	orig := safecastSubmitClient
	safecastSubmitClient = fake
	defer func() { safecastSubmitClient = orig }()

	uploadID := insertTestUpload(t, db, "track-4", "log4.log")
	submitToSafecastIfNeeded("track-4", uploadID, "log4.log", []byte("data"), "key", "7", db)

	importID, submitErr := getSafecastStatus(t, db, uploadID)
	if importID != "" {
		t.Errorf("expected no import id on failure, got %q", importID)
	}
	if submitErr == "" {
		t.Error("expected submit error to be recorded")
	}
}

func TestSubmitToSafecastIfNeeded_RecordsCheckExistsFailure(t *testing.T) {
	db := newSafecastSubmitTestDB(t)
	fake := &fakeSafecastSubmitter{existsErr: errors.New("network error")}
	orig := safecastSubmitClient
	safecastSubmitClient = fake
	defer func() { safecastSubmitClient = orig }()

	uploadID := insertTestUpload(t, db, "track-5", "log5.log")
	submitToSafecastIfNeeded("track-5", uploadID, "log5.log", []byte("data"), "key", "7", db)

	if fake.submitCalled {
		t.Fatal("expected Submit not to be called when CheckExists errors")
	}
	_, submitErr := getSafecastStatus(t, db, uploadID)
	if submitErr == "" {
		t.Error("expected check-exists error to be recorded")
	}
}

func TestSubmitToSafecastIfNeeded_KeyedByUploadIDNotTrackID(t *testing.T) {
	// Two files uploaded in one request share the same track_id (assigned once
	// per request). Their statuses must not clobber each other.
	db := newSafecastSubmitTestDB(t)
	uploadA := insertTestUpload(t, db, "shared-track", "a.log")
	uploadB := insertTestUpload(t, db, "shared-track", "b.log")

	fakeA := &fakeSafecastSubmitter{submitID: "111"}
	orig := safecastSubmitClient
	safecastSubmitClient = fakeA
	submitToSafecastIfNeeded("shared-track", uploadA, "a.log", []byte("data"), "key", "7", db)

	fakeB := &fakeSafecastSubmitter{submitID: "222"}
	safecastSubmitClient = fakeB
	submitToSafecastIfNeeded("shared-track", uploadB, "b.log", []byte("data"), "key", "7", db)
	safecastSubmitClient = orig

	importIDA, _ := getSafecastStatus(t, db, uploadA)
	importIDB, _ := getSafecastStatus(t, db, uploadB)
	if importIDA != "111" {
		t.Errorf("got upload A importID %q, want %q", importIDA, "111")
	}
	if importIDB != "222" {
		t.Errorf("got upload B importID %q, want %q", importIDB, "222")
	}
}

// withHandlerTestDB points the package-level db/dbType globals used by
// updateSafecastCredentialsHandler at a temporary sqlite DB for the duration
// of the test, restoring the previous globals afterward.
func withHandlerTestDB(t *testing.T) *database.Database {
	t.Helper()
	testDB := newSafecastSubmitTestDB(t)

	origDB := db
	origDBType := *dbType
	db = testDB
	*dbType = "sqlite"
	t.Cleanup(func() {
		db = origDB
		*dbType = origDBType
	})
	return testDB
}

func createTestUser(t *testing.T, testDB *database.Database, email string) *auth.User {
	t.Helper()
	user := &auth.User{Email: email, IsActive: true}
	id, err := auth.CreateUser(context.Background(), testDB.DB, "sqlite", user, "password123")
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	user.ID = id
	return user
}

func TestUpdateSafecastCredentialsHandler_SetsCredentials(t *testing.T) {
	testDB := withHandlerTestDB(t)
	user := createTestUser(t, testDB, "setter@example.com")

	fake := &fakeSafecastSubmitter{resolveUserID: "55"}
	origClient := safecastSubmitClient
	safecastSubmitClient = fake
	t.Cleanup(func() { safecastSubmitClient = origClient })

	body, _ := json.Marshal(map[string]string{"safecast_api_key": "my-safecast-key"})
	req := httptest.NewRequest(http.MethodPost, "/api/user/safecast-credentials", bytes.NewReader(body))
	req = req.WithContext(auth.WithAuthContext(req.Context(), user))
	w := httptest.NewRecorder()

	updateSafecastCredentialsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	updated, err := auth.GetUserByID(context.Background(), testDB.DB, "sqlite", user.ID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if updated.SafecastAPIKey != "my-safecast-key" {
		t.Errorf("got SafecastAPIKey %q, want %q", updated.SafecastAPIKey, "my-safecast-key")
	}
	if updated.SafecastUserID != "55" {
		t.Errorf("got SafecastUserID %q, want %q", updated.SafecastUserID, "55")
	}
}

func TestUpdateSafecastCredentialsHandler_ClearsCredentials(t *testing.T) {
	testDB := withHandlerTestDB(t)
	user := createTestUser(t, testDB, "clearer@example.com")
	if err := auth.UpdateUserSafecastCredentials(context.Background(), testDB.DB, "sqlite", user.ID, "old-key", "10"); err != nil {
		t.Fatalf("seed UpdateUserSafecastCredentials: %v", err)
	}

	body, _ := json.Marshal(map[string]string{"safecast_api_key": ""})
	req := httptest.NewRequest(http.MethodPost, "/api/user/safecast-credentials", bytes.NewReader(body))
	req = req.WithContext(auth.WithAuthContext(req.Context(), user))
	w := httptest.NewRecorder()

	updateSafecastCredentialsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	updated, err := auth.GetUserByID(context.Background(), testDB.DB, "sqlite", user.ID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if updated.SafecastAPIKey != "" || updated.SafecastUserID != "" {
		t.Errorf("expected cleared credentials, got apiKey=%q userID=%q", updated.SafecastAPIKey, updated.SafecastUserID)
	}
}

func TestUpdateSafecastCredentialsHandler_RejectsBadKey(t *testing.T) {
	testDB := withHandlerTestDB(t)
	user := createTestUser(t, testDB, "badkey@example.com")

	fake := &fakeSafecastSubmitter{resolveErr: errors.New("unauthorized")}
	origClient := safecastSubmitClient
	safecastSubmitClient = fake
	t.Cleanup(func() { safecastSubmitClient = origClient })

	body, _ := json.Marshal(map[string]string{"safecast_api_key": "bad-key"})
	req := httptest.NewRequest(http.MethodPost, "/api/user/safecast-credentials", bytes.NewReader(body))
	req = req.WithContext(auth.WithAuthContext(req.Context(), user))
	w := httptest.NewRecorder()

	updateSafecastCredentialsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}

	updated, err := auth.GetUserByID(context.Background(), testDB.DB, "sqlite", user.ID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if updated.SafecastAPIKey != "" {
		t.Errorf("expected no credentials saved on failure, got %q", updated.SafecastAPIKey)
	}
}

func TestUpdateSafecastCredentialsHandler_RequiresAuth(t *testing.T) {
	withHandlerTestDB(t)

	body, _ := json.Marshal(map[string]string{"safecast_api_key": "some-key"})
	req := httptest.NewRequest(http.MethodPost, "/api/user/safecast-credentials", bytes.NewReader(body))
	w := httptest.NewRecorder()

	updateSafecastCredentialsHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}
}
