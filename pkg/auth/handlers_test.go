package auth

// NOTE: Tests that call RegisterHandler involve bcrypt at cost 12, which is
// intentionally slow (~200ms per call). Keep the number of register calls low
// or parallelize tests to avoid a slow suite.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"safecast-new-map/pkg/database"
	_ "safecast-new-map/pkg/database/drivers"
)

// newTestManager creates a Manager backed by a temporary SQLite database.
// EmailSender is nil so no actual emails are sent during tests.
func newTestManager(t *testing.T) *Manager {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "auth_test.db")
	cfg := database.Config{DBType: "sqlite", DBPath: dbPath, Port: 1}
	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase: %v", err)
	}
	if err := db.InitSchema(cfg); err != nil {
		t.Fatalf("InitSchema: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return &Manager{
		DB:                db.DB,
		DBDriver:          "sqlite",
		SessionCookieName: "session",
		AllowRegistration: true,
		EmailSender:       nil, // no emails in tests
		BaseURL:           "http://localhost:8765",
	}
}

// registerTestUser is a helper that registers a user and returns the JSON response.
// bcrypt cost 12 makes this slow (~200ms); only call once per test group.
func registerTestUser(t *testing.T, m *Manager, email, username, password string) map[string]interface{} {
	t.Helper()
	body := `{"email":"` + email + `","username":"` + username + `","password":"` + password + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	m.RegisterHandler(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("registerTestUser: got status %d, want 201; body: %s", rec.Code, rec.Body.String())
	}
	var out map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatalf("registerTestUser decode: %v", err)
	}
	return out
}

// --- POST /api/auth/register ---

func TestRegisterHandler(t *testing.T) {
	m := newTestManager(t)

	t.Run("GET returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/auth/register", nil)
		rec := httptest.NewRecorder()
		m.RegisterHandler(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got %d, want 405", rec.Code)
		}
	})

	t.Run("invalid JSON returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader("not json"))
		rec := httptest.NewRecorder()
		m.RegisterHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("invalid email returns 400", func(t *testing.T) {
		body := `{"email":"notanemail","username":"alice","password":"Password123!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
		rec := httptest.NewRecorder()
		m.RegisterHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("weak password returns 400", func(t *testing.T) {
		body := `{"email":"test@example.com","username":"alice","password":"abc"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
		rec := httptest.NewRecorder()
		m.RegisterHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("successful registration returns 201", func(t *testing.T) {
		out := registerTestUser(t, m, "alice@example.com", "alice", "Password123!")
		if out["success"] != true {
			t.Errorf("success = %v, want true", out["success"])
		}
		if _, ok := out["userId"]; !ok {
			t.Error("response missing userId")
		}
	})

	t.Run("duplicate email returns 409", func(t *testing.T) {
		// alice was already registered above in the same manager
		body := `{"email":"alice@example.com","username":"alice2","password":"Password123!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
		rec := httptest.NewRecorder()
		m.RegisterHandler(rec, req)
		if rec.Code != http.StatusConflict {
			t.Errorf("got %d, want 409", rec.Code)
		}
	})

	t.Run("registration disabled returns 403", func(t *testing.T) {
		m2 := newTestManager(t)
		m2.AllowRegistration = false
		body := `{"email":"bob@example.com","username":"bob","password":"Password123!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body))
		rec := httptest.NewRecorder()
		m2.RegisterHandler(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Errorf("got %d, want 403", rec.Code)
		}
	})
}

// --- POST /api/auth/login ---

func TestLoginHandler(t *testing.T) {
	m := newTestManager(t)
	// Register once, reuse for all login sub-tests
	registerTestUser(t, m, "login@example.com", "loginuser", "Password123!")

	t.Run("GET returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
		rec := httptest.NewRecorder()
		m.LoginHandler(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got %d, want 405", rec.Code)
		}
	})

	t.Run("no credentials returns 400", func(t *testing.T) {
		body := `{"email":"login@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
		rec := httptest.NewRecorder()
		m.LoginHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	t.Run("wrong password returns 401", func(t *testing.T) {
		body := `{"email":"login@example.com","password":"wrongpassword"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
		rec := httptest.NewRecorder()
		m.LoginHandler(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", rec.Code)
		}
	})

	t.Run("unknown email returns 401", func(t *testing.T) {
		body := `{"email":"nobody@example.com","password":"Password123!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
		rec := httptest.NewRecorder()
		m.LoginHandler(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", rec.Code)
		}
	})

	t.Run("valid credentials returns 200 with session cookie", func(t *testing.T) {
		body := `{"email":"login@example.com","password":"Password123!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
		rec := httptest.NewRecorder()
		m.LoginHandler(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
		var out map[string]interface{}
		if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if out["success"] != true {
			t.Errorf("success = %v, want true", out["success"])
		}
		if _, ok := out["user"]; !ok {
			t.Error("response missing user")
		}
		// Session cookie should be set
		found := false
		for _, c := range rec.Result().Cookies() {
			if c.Name == m.SessionCookieName {
				found = true
				break
			}
		}
		if !found {
			t.Error("session cookie not set")
		}
	})
}

// --- POST /api/auth/logout ---

func TestLogoutHandler(t *testing.T) {
	m := newTestManager(t)

	t.Run("GET returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/auth/logout", nil)
		rec := httptest.NewRecorder()
		m.LogoutHandler(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got %d, want 405", rec.Code)
		}
	})

	t.Run("logout without cookie returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		rec := httptest.NewRecorder()
		m.LogoutHandler(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("logout with valid session clears cookie", func(t *testing.T) {
		// Register and login to get a session
		registerTestUser(t, m, "logout@example.com", "logoutuser", "Password123!")
		loginBody := `{"email":"logout@example.com","password":"Password123!"}`
		loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(loginBody))
		loginRec := httptest.NewRecorder()
		m.LoginHandler(loginRec, loginReq)
		if loginRec.Code != http.StatusOK {
			t.Fatalf("login failed: %d", loginRec.Code)
		}
		var sessionCookie *http.Cookie
		for _, c := range loginRec.Result().Cookies() {
			if c.Name == m.SessionCookieName {
				sessionCookie = c
				break
			}
		}
		if sessionCookie == nil {
			t.Fatal("no session cookie after login")
		}

		// Now logout
		logoutReq := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		logoutReq.AddCookie(sessionCookie)
		logoutRec := httptest.NewRecorder()
		m.LogoutHandler(logoutRec, logoutReq)
		if logoutRec.Code != http.StatusOK {
			t.Errorf("got %d, want 200", logoutRec.Code)
		}
	})
}

// --- POST /api/auth/forgot-password ---

func TestForgotPasswordHandler(t *testing.T) {
	m := newTestManager(t)
	registerTestUser(t, m, "forgot@example.com", "forgotuser", "Password123!")

	t.Run("GET returns 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/auth/forgot-password", nil)
		rec := httptest.NewRecorder()
		m.ForgotPasswordHandler(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("got %d, want 405", rec.Code)
		}
	})

	t.Run("unknown email still returns 200 (anti-enumeration)", func(t *testing.T) {
		body := `{"email":"nobody@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", strings.NewReader(body))
		rec := httptest.NewRecorder()
		m.ForgotPasswordHandler(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got %d, want 200 (anti-enumeration)", rec.Code)
		}
	})

	t.Run("known email returns 200", func(t *testing.T) {
		body := `{"email":"forgot@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", strings.NewReader(body))
		rec := httptest.NewRecorder()
		m.ForgotPasswordHandler(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got %d, want 200; body: %s", rec.Code, rec.Body.String())
		}
	})
}

// --- GET /api/user/profile ---
// ProfileHandler reads user from context set by RequireAuth middleware.
// Tests use m.RequireAuth(m.ProfileHandler) to match production wiring.

func TestProfileHandler(t *testing.T) {
	m := newTestManager(t)
	handler := m.RequireAuth(m.ProfileHandler)

	t.Run("unauthenticated returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
		rec := httptest.NewRecorder()
		handler(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", rec.Code)
		}
	})

	t.Run("authenticated user returns profile", func(t *testing.T) {
		registerTestUser(t, m, "profile@example.com", "profileuser", "Password123!")

		// Login to get session
		loginBody := `{"email":"profile@example.com","password":"Password123!"}`
		loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(loginBody))
		loginRec := httptest.NewRecorder()
		m.LoginHandler(loginRec, loginReq)
		if loginRec.Code != http.StatusOK {
			t.Fatalf("login failed: %d", loginRec.Code)
		}
		var sessionCookie *http.Cookie
		for _, c := range loginRec.Result().Cookies() {
			if c.Name == m.SessionCookieName {
				sessionCookie = c
				break
			}
		}
		if sessionCookie == nil {
			t.Fatal("no session cookie")
		}

		// Request profile via RequireAuth-wrapped handler
		profileReq := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
		profileReq.AddCookie(sessionCookie)
		profileRec := httptest.NewRecorder()
		handler(profileRec, profileReq)
		if profileRec.Code != http.StatusOK {
			t.Fatalf("got %d, want 200; body: %s", profileRec.Code, profileRec.Body.String())
		}
		var out map[string]interface{}
		if err := json.NewDecoder(profileRec.Body).Decode(&out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if out["email"] != "profile@example.com" {
			t.Errorf("email = %v, want profile@example.com", out["email"])
		}
		if _, ok := out["username"]; !ok {
			t.Error("response missing username")
		}
	})
}

// --- POST /api/user/change-password ---
// ChangePasswordHandler also reads user from context (RequireAuth required).

func TestChangePasswordHandler(t *testing.T) {
	m := newTestManager(t)
	handler := m.RequireAuth(m.ChangePasswordHandler)
	registerTestUser(t, m, "chpwd@example.com", "chpwduser", "Password123!")

	// Login to get session
	loginBody := `{"email":"chpwd@example.com","password":"Password123!"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(loginBody))
	loginRec := httptest.NewRecorder()
	m.LoginHandler(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login setup failed: %d", loginRec.Code)
	}
	var sessionCookie *http.Cookie
	for _, c := range loginRec.Result().Cookies() {
		if c.Name == m.SessionCookieName {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("no session cookie")
	}

	t.Run("unauthenticated returns 401", func(t *testing.T) {
		body := `{"current_password":"Password123!","new_password":"NewPass456!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/change-password", strings.NewReader(body))
		rec := httptest.NewRecorder()
		handler(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", rec.Code)
		}
	})

	t.Run("wrong current password returns 401", func(t *testing.T) {
		body := `{"current_password":"wrongpass","new_password":"NewPass456!"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user/change-password", strings.NewReader(body))
		req.AddCookie(sessionCookie)
		rec := httptest.NewRecorder()
		handler(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got %d, want 401", rec.Code)
		}
	})
}
