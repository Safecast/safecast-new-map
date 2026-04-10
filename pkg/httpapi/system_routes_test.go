package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterSystemRoutesInventory(t *testing.T) {
	mux := http.NewServeMux()
	RegisterSystemRoutes(mux, SystemRoutesConfig{
		SelfUpgradeHandler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
	})

	req := httptest.NewRequest(http.MethodGet, "/selfupgrade/status", nil)
	_, pattern := mux.Handler(req)
	if pattern != "/selfupgrade/" {
		t.Fatalf("expected /selfupgrade/ pattern, got %q", pattern)
	}
}
