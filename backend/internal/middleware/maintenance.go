package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/maintenance"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
)

func MaintenanceMode(state *maintenance.State, enabled bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !enabled || !state.Get().Enabled || strings.HasPrefix(r.URL.Path, "/api/v1/admin/") || r.URL.Path == "/api/v1/auth/login" || r.URL.Path == "/api/v1/auth/refresh" || r.URL.Path == "/health" || r.URL.Path == "/healthz" || r.URL.Path == "/readyz" || r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}
			if r.URL.Path == "/ncs-ussd" {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = fmt.Fprint(w, "END System is temporarily undergoing maintenance. Please try again later.")
				return
			}
			if strings.HasPrefix(r.URL.Path, "/api/") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				_ = json.NewEncoder(w).Encode(map[string]any{"status": 503, "message": "System undergoing scheduled maintenance"})
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = fmt.Fprint(w, `<!doctype html><html lang="en"><meta name="viewport" content="width=device-width"><title>NCS maintenance</title><body style="font-family:system-ui;background:#112b4e;color:white;display:grid;place-items:center;min-height:100vh;margin:0"><main style="max-width:38rem;text-align:center;padding:2rem"><h1>We’ll be right back</h1><p>The National Council of Sports platform is undergoing scheduled maintenance.</p></main></body></html>`)
		})
	}
}

func ValidateGlobalSession(state *maintenance.State) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			issued, _ := r.Context().Value(models.CtxTokenIssuedAt).(time.Time)
			cutoff := state.Get().GlobalAuthInvalidBefore
			if cutoff != nil && !issued.After(*cutoff) {
				response.Err(w, http.StatusUnauthorized, "SESSION_REVOKED", "Session has been revoked. Please sign in again")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func WithMaintenanceState(ctx context.Context, state *maintenance.State) context.Context {
	return context.WithValue(ctx, struct{ maintenance string }{"state"}, state)
}
