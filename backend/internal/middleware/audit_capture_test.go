package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScrubPayloadKeepsAllowlistAndRestoresBody(t *testing.T) {
	body := `{"action":"submit","password":"secret","form_type":"form_3"}`
	r := httptest.NewRequest("POST", "/api/v1/applications", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	safe := string(scrubPayload(r))
	if strings.Contains(safe, "secret") || !strings.Contains(safe, "form_3") {
		t.Fatalf("unexpected scrubbed payload: %s", safe)
	}
	restored, _ := io.ReadAll(r.Body)
	if string(restored) != body {
		t.Fatalf("request body was not restored")
	}
}

func TestHoneypotRedirectsScannerProbe(t *testing.T) {
	r := httptest.NewRequest("GET", "/.env", nil)
	w := httptest.NewRecorder()
	HoneypotScanner(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { t.Fatal("probe reached application") })).ServeHTTP(w, r)
	if w.Code != 307 || w.Header().Get("Location") != "/blank.gif" {
		t.Fatalf("unexpected response: %d %s", w.Code, w.Header().Get("Location"))
	}
}
