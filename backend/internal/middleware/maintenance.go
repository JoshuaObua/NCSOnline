package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/maintenance"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
)

const maintenanceBypassCookie = "ncs_maintenance_bypass"

func MaintenanceMode(state *maintenance.State, enabled bool, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !enabled || isMaintenanceAlwaysAllowed(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			snapshot := state.Get()
			scope := maintenanceScopeForPath(r.URL.Path)
			scoped := snapshot.Scoped(scope)
			if !scoped.IsActiveAt(time.Now()) || hasMaintenanceBypass(w, r, scoped, jwtSecret) {
				next.ServeHTTP(w, r)
				return
			}

			if scope == maintenance.ScopeAdminDashboard {
				response.Err(w, http.StatusServiceUnavailable, "ADMIN_MAINTENANCE", publicMaintenanceMessage(scoped))
				return
			}
			if r.URL.Path == "/ncs-ussd" {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = fmt.Fprintf(w, "END %s", ussdMaintenanceMessage(scoped))
				return
			}
			if strings.HasPrefix(r.URL.Path, "/api/") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				_ = json.NewEncoder(w).Encode(map[string]any{"status": 503, "message": publicMaintenanceMessage(scoped), "scope": scoped.Scope})
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = fmt.Fprintf(w, `<!doctype html><html lang="en"><meta name="viewport" content="width=device-width"><title>NCS maintenance</title><body style="font-family:system-ui;background:#112b4e;color:white;display:grid;place-items:center;min-height:100vh;margin:0"><main style="max-width:38rem;text-align:center;padding:2rem"><h1>%s</h1><p>%s</p></main></body></html>`, htmlEscape(scoped.DisplayMeta.CustomTitle), htmlEscape(publicMaintenanceMessage(scoped)))
		})
	}
}

func isMaintenanceAlwaysAllowed(path string) bool {
	switch path {
	case "/health", "/healthz", "/readyz", "/metrics", "/api/v1/auth/login", "/api/v1/auth/refresh", "/api/v1/system/maintenance-status":
		return true
	case "/favicon.ico", "/favicon.png", "/robots.txt", "/manifest.json", "/manifest.webmanifest", "/sw.js":
		return true
	}
	// Static SPA bundles must never be gated — without them the SPA cannot
	// boot to read the bypass param/cookie for authorised users.
	if strings.HasPrefix(path, "/assets/") ||
		strings.HasPrefix(path, "/static/") ||
		strings.HasPrefix(path, "/fonts/") ||
		strings.HasPrefix(path, "/img/") ||
		strings.HasPrefix(path, "/images/") {
		return true
	}
	return false
}

func maintenanceScopeForPath(path string) string {
	if strings.HasPrefix(path, "/api/v1/admin/") || path == "/api/v1/me" || strings.HasPrefix(path, "/admin") || strings.HasPrefix(path, "/dashboard") {
		return maintenance.ScopeAdminDashboard
	}
	return maintenance.ScopePublicCMS
}

func publicMaintenanceMessage(scoped maintenance.ScopedSnapshot) string {
	if scoped.DisplayMeta.CustomMessage != "" {
		return scoped.DisplayMeta.CustomMessage
	}
	if scoped.Reason != "" {
		return scoped.Reason
	}
	return "System undergoing scheduled maintenance"
}

func ussdMaintenanceMessage(scoped maintenance.ScopedSnapshot) string {
	if scoped.Reason != "" {
		return scoped.Reason
	}
	return "System is temporarily undergoing maintenance. Please try again later."
}

func htmlEscape(v string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;")
	return replacer.Replace(v)
}

func hasMaintenanceBypass(w http.ResponseWriter, r *http.Request, scoped maintenance.ScopedSnapshot, jwtSecret string) bool {
	if scoped.Scope == maintenance.ScopeAdminDashboard && isBreakGlassAdminPath(r.URL.Path) {
		if claims := bearerClaims(r, jwtSecret); hasAnyRole(claims, []string{"super_admin"}) {
			return true
		}
	}
	if claims := bearerClaims(r, jwtSecret); hasAnyRole(claims, scoped.BypassRules.AllowedRoles) {
		return true
	}
	if claims := bearerClaims(r, jwtSecret); userAllowed(claims, scoped.BypassRules.AllowedUserIDs) {
		return true
	}
	if ipAllowed(realIP(r), scoped.BypassRules.AllowedIPRanges) {
		return true
	}
	token := strings.TrimSpace(scoped.BypassRules.SecretQueryParam)
	if token == "" {
		return false
	}
	if r.URL.Query().Get("maintenance_bypass") == token || r.URL.Query().Has(token) {
		setMaintenanceBypassCookie(w, scoped.Scope, token, jwtSecret)
		return true
	}
	return hasValidMaintenanceBypassCookie(r, scoped.Scope, token, jwtSecret)
}

func isBreakGlassAdminPath(path string) bool {
	return path == "/api/v1/admin/system/status" ||
		path == "/api/v1/admin/system/resources" ||
		path == "/api/v1/admin/system/maintenance" ||
		strings.HasPrefix(path, "/api/v1/admin/system/maintenance/")
}

func bearerClaims(r *http.Request, jwtSecret string) *Claims {
	if jwtSecret == "" {
		return nil
	}
	authz := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(strings.ToLower(authz), "bearer ") {
		return nil
	}
	claims, err := ParseToken(jwtSecret, strings.TrimSpace(authz[7:]))
	if err != nil {
		return nil
	}
	return claims
}

func hasAnyRole(claims *Claims, allowed []string) bool {
	if claims == nil {
		return false
	}
	for _, role := range claims.Roles {
		for _, allowedRole := range allowed {
			if strings.EqualFold(strings.TrimSpace(role), strings.TrimSpace(allowedRole)) {
				return true
			}
		}
	}
	return false
}

func userAllowed(claims *Claims, allowed []string) bool {
	if claims == nil || claims.UserID == "" {
		return false
	}
	for _, userID := range allowed {
		if strings.EqualFold(strings.TrimSpace(userID), claims.UserID) {
			return true
		}
	}
	return false
}

func ipAllowed(ipText string, ranges []string) bool {
	if len(ranges) == 0 {
		return false
	}
	ip := net.ParseIP(ipText)
	if ip == nil {
		return false
	}
	for _, item := range ranges {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if strings.Contains(item, "/") {
			if _, network, err := net.ParseCIDR(item); err == nil && network.Contains(ip) {
				return true
			}
			continue
		}
		if parsed := net.ParseIP(item); parsed != nil && parsed.Equal(ip) {
			return true
		}
	}
	return false
}

func setMaintenanceBypassCookie(w http.ResponseWriter, scope, token, secret string) {
	exp := time.Now().Add(4 * time.Hour).Unix()
	value := fmt.Sprintf("%s|%d|%s", scope, exp, signMaintenanceBypass(scope, exp, token, secret))
	http.SetCookie(w, &http.Cookie{Name: maintenanceBypassCookie, Value: value, Path: "/", Expires: time.Unix(exp, 0), HttpOnly: true, SameSite: http.SameSiteLaxMode})
}

func hasValidMaintenanceBypassCookie(r *http.Request, scope, token, secret string) bool {
	c, err := r.Cookie(maintenanceBypassCookie)
	if err != nil {
		return false
	}
	parts := strings.Split(c.Value, "|")
	if len(parts) != 3 || parts[0] != scope {
		return false
	}
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return false
	}
	expected := signMaintenanceBypass(scope, exp, token, secret)
	return hmac.Equal([]byte(expected), []byte(parts[2]))
}

func signMaintenanceBypass(scope string, exp int64, token, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = fmt.Fprintf(mac, "%s|%d|%s", scope, exp, token)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
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
