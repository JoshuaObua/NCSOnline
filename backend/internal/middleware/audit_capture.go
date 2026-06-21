package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

var auditPayloadAllowlist = map[string]struct{}{"action": {}, "application_type": {}, "form_type": {}, "organisation_type": {}, "status": {}, "channel": {}}

func scrubPayload(r *http.Request) json.RawMessage {
	if r.Body == nil || !strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "application/json") {
		return nil
	}
	limited, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	r.Body = io.NopCloser(bytes.NewReader(limited))
	if len(limited) > 8192 {
		return nil
	}
	var raw map[string]any
	if json.Unmarshal(limited, &raw) != nil {
		return nil
	}
	safe := map[string]any{}
	for key, value := range raw {
		if _, ok := auditPayloadAllowlist[strings.ToLower(key)]; ok {
			safe[key] = value
		}
	}
	if len(safe) == 0 {
		return nil
	}
	encoded, _ := json.Marshal(safe)
	if len(encoded) > 2048 {
		return nil
	}
	return encoded
}

func fingerprintRequest(r *http.Request) (event string, anomaly bool) {
	ua := strings.ToLower(r.UserAgent())
	path := strings.ToLower(r.URL.Path)
	if strings.Contains(ua, "mobile") && strings.HasPrefix(path, "/api/") && r.Header.Get("X-NCS-SDK-Version") == "" {
		return "security.client.fingerprint", true
	}
	if strings.Contains(path, "/admin/") && strings.Contains(ua, "mozilla") && r.Referer() == "" {
		return "security.admin.no_referrer", true
	}
	return "", false
}

func HoneypotScanner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.ToLower(r.URL.Path)
		for _, probe := range []string{"/.env", "/wp-admin", "/phpmyadmin", "/.git/", "/server-status"} {
			if strings.HasPrefix(p, probe) {
				http.Redirect(w, r, "/blank.gif", http.StatusTemporaryRedirect)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
