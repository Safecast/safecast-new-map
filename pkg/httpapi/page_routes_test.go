package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"safecast-new-map/pkg/auth"
)

func TestRegisterPageRoutesPublicInventory(t *testing.T) {
	mux := http.NewServeMux()
	RegisterPageRoutes(mux, PageRoutesConfig{
		HomeHandler:        func(http.ResponseWriter, *http.Request) {},
		StoriesPageHandler: func(http.ResponseWriter, *http.Request) {},
		StoriesDataHandler: func(http.ResponseWriter, *http.Request) {},
		MapHandler:         func(http.ResponseWriter, *http.Request) {},
	})

	requiredRoutes := []string{
		"/home",
		"/stories.html",
		"/data/stories.json",
		"/",
	}

	for _, route := range requiredRoutes {
		t.Run(route, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, route, nil)
			_, pattern := mux.Handler(req)
			if pattern == "" {
				t.Fatalf("route %s was not registered", route)
			}
		})
	}
}

func TestRegisterPageRoutesAuthInventory(t *testing.T) {
	mux := http.NewServeMux()
	RegisterPageRoutes(mux, PageRoutesConfig{
		AuthManager:           &auth.Manager{},
		AdminPassword:         "secret",
		ProfileHandler:        func(http.ResponseWriter, *http.Request) {},
		ResetPasswordHandler:  func(http.ResponseWriter, *http.Request) {},
		AdminUsersPageHandler: func(http.ResponseWriter, *http.Request) {},
	})

	requiredRoutes := []string{
		"/profile",
		"/reset-password",
		"/admin/users",
	}

	for _, route := range requiredRoutes {
		t.Run(route, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, route, nil)
			_, pattern := mux.Handler(req)
			if pattern == "" {
				t.Fatalf("route %s was not registered", route)
			}
		})
	}
}

func TestRegisterPageRoutesAdminGuard(t *testing.T) {
	mux := http.NewServeMux()
	RegisterPageRoutes(mux, PageRoutesConfig{
		AuthManager:           &auth.Manager{},
		AdminPassword:         "secret",
		AdminUsersPageHandler: func(http.ResponseWriter, *http.Request) {},
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
