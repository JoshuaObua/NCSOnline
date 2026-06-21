package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/atenimedia-llc/ncs-online/backend/internal/maintenance"
)

func TestMaintenanceResponsesByChannel(t *testing.T) {
	state := maintenance.New(maintenance.Snapshot{Enabled: true})
	h := MaintenanceMode(state, true)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { t.Fatal("request bypassed maintenance") }))
	for _, tc := range []struct {
		path     string
		status   int
		contains string
	}{{"/api/v1/cms/posts", 503, "scheduled maintenance"}, {"/ncs-ussd", 503, "END System"}, {"/", 503, "right back"}} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
		if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.contains) {
			t.Fatalf("%s: %d %s", tc.path, w.Code, w.Body.String())
		}
	}
}

func TestMaintenanceAllowsHealthAndAdmin(t *testing.T) {
	state := maintenance.New(maintenance.Snapshot{Enabled: true})
	called := 0
	h := MaintenanceMode(state, true)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called++ }))
	for _, path := range []string{"/healthz", "/readyz", "/api/v1/admin/system/status", "/api/v1/auth/login"} {
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", path, nil))
	}
	if called != 4 {
		t.Fatalf("expected 4 bypasses, got %d", called)
	}
}
