package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestRegisterStaticRoutesInventory(t *testing.T) {
	mux := http.NewServeMux()
	RegisterStaticRoutes(mux, StaticRoutesConfig{
		StaticFS: fstest.MapFS{
			"index.txt": &fstest.MapFile{Data: []byte("ok")},
		},
		JSDir: ".",
	})

	requiredRoutes := []string{
		"/static/index.txt",
		"/js/app.js",
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
