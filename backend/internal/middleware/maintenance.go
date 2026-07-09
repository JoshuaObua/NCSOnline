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
			_, _ = fmt.Fprint(w, publicMaintenanceHTML(scoped))
		})
	}
}

func isMaintenanceAlwaysAllowed(path string) bool {
	switch path {
	case "/health", "/healthz", "/readyz", "/metrics", "/api/v1/auth/login", "/api/v1/auth/refresh", "/api/v1/system/maintenance-status":
		return true
	case "/favicon.ico", "/favicon.png", "/main-logo.png", "/main-logo-white.png", "/robots.txt", "/manifest.json", "/manifest.webmanifest", "/sw.js":
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

func publicMaintenanceHTML(scoped maintenance.ScopedSnapshot) string {
	title := scoped.DisplayMeta.CustomTitle
	if title == "" {
		title = "We'll be right back"
	}
	expected := "Shortly"
	if scoped.ExpectedEnd != nil {
		expected = scoped.ExpectedEnd.Format("Mon, 02 Jan at 15:04 MST")
	}
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>%s | NCS Uganda</title>
<style>
*{box-sizing:border-box}body{min-height:100vh;margin:0;display:flex;flex-direction:column;justify-content:center;overflow-x:hidden;background:linear-gradient(135deg,#102b4d 0%%,#1a365d 52%%,#0d223d 100%%);color:#fff;font-family:Arial,sans-serif;padding:clamp(20px,4vw,56px)}body:before{position:fixed;inset:0;content:"";pointer-events:none;background-image:linear-gradient(rgba(255,255,255,.045) 1px,transparent 1px),linear-gradient(90deg,rgba(255,255,255,.045) 1px,transparent 1px);background-size:52px 52px;mask-image:linear-gradient(to bottom right,#000,transparent 78%%)}.shell{position:relative;width:min(100%%,1160px);margin:auto;display:grid;grid-template-columns:minmax(0,1.08fr) minmax(320px,.72fr);align-items:center;gap:clamp(32px,7vw,104px)}.brand{display:flex;align-items:center;gap:14px;margin-bottom:clamp(36px,7vh,80px);color:rgba(255,255,255,.86);font-size:13px;font-weight:700;text-transform:uppercase}.brand img{width:auto;max-width:168px;height:60px;object-fit:contain;border-radius:8px;background:#fff;padding:6px 11px}.kicker{display:flex;align-items:center;gap:9px;margin:0 0 18px;color:#f5a623;font-size:12px;font-weight:800;text-transform:uppercase}.kicker:before{width:28px;height:2px;content:"";background:#f5a623}h1{max-width:720px;margin:0;font-size:clamp(42px,6vw,86px);line-height:1.02;overflow-wrap:anywhere}p.message{max-width:650px;margin:24px 0 0;color:rgba(255,255,255,.76);font-size:clamp(16px,1.8vw,19px);line-height:1.75}.actions{display:flex;flex-wrap:wrap;gap:12px;margin-top:32px}.actions a{min-height:44px;display:inline-flex;align-items:center;justify-content:center;border:1px solid #f5a623;border-radius:6px;background:#f5a623;color:#102b4d;padding:12px 18px;text-decoration:none;font-size:14px;font-weight:800}.actions a.alt{border-color:rgba(255,255,255,.34);background:transparent;color:#fff}.panel{border:1px solid rgba(255,255,255,.22);border-radius:8px;background:#fff;color:#1a365d;padding:clamp(24px,4vw,38px);box-shadow:0 28px 80px rgba(4,18,36,.34)}.icon{width:56px;height:56px;display:grid;place-items:center;border-radius:8px;background:#fff5e4;color:#e2920f;font-size:29px;font-weight:700}.label{margin:24px 0 8px;color:#e2920f;font-size:12px;font-weight:800;text-transform:uppercase}.panel h2{margin:0;font-size:clamp(22px,2.6vw,28px);line-height:1.25}.meta{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:1px;margin-top:24px;overflow:hidden;border:1px solid #e5e9ef;border-radius:6px;background:#e5e9ef}.meta div{min-width:0;background:#f8fafc;padding:16px}.meta span,.meta strong{display:block}.meta span{margin-bottom:6px;color:#6b7280;font-size:11px;font-weight:700;text-transform:uppercase}.meta strong{font-size:14px;line-height:1.45;overflow-wrap:anywhere}.steps{display:grid;gap:14px;margin:24px 0 0;padding:22px 0 0;border-top:1px solid #e5e9ef;list-style:none}.steps li{display:grid;grid-template-columns:12px minmax(0,1fr);align-items:center;gap:11px;color:#8a94a3;font-size:14px;font-weight:700}.steps i{width:10px;height:10px;border:2px solid #cbd5e1;border-radius:50%%}.steps .done{color:#237a4b}.steps .done i{border-color:#2fa96b;background:#2fa96b}.steps .active{color:#1a365d}.steps .active i{border-color:#f5a623;background:#f5a623;box-shadow:0 0 0 4px rgba(245,166,35,.18)}footer{position:relative;width:min(100%%,1160px);margin:clamp(32px,7vh,72px) auto 0;color:rgba(255,255,255,.48);font-size:12px}@media(max-width:820px){body{justify-content:flex-start}.shell{grid-template-columns:1fr;gap:36px}.brand{margin-bottom:40px}.panel{max-width:620px}}@media(max-width:480px){body{padding:18px}.brand{align-items:flex-start;flex-direction:column}.brand img{height:52px;max-width:150px}.actions{display:grid}.actions a{width:100%%}.meta{grid-template-columns:1fr}}
</style>
</head>
<body>
<main class="shell">
<section>
<div class="brand"><img src="/main-logo.png" alt="National Council of Sports logo"><span>NCS Uganda</span></div>
<p class="kicker">Public Portal Maintenance</p>
<h1>%s</h1>
<p class="message">%s</p>
<div class="actions"><a href="mailto:info@ncs.go.ug">Email NCS</a><a class="alt" href="tel:+256414254477">Call NCS</a></div>
</section>
<aside class="panel">
<div class="icon" aria-hidden="true">&#9881;</div>
<p class="label">Current Status</p>
<h2>Scheduled upgrades are in progress</h2>
<div class="meta"><div><span>Expected return</span><strong>%s</strong></div><div><span>System status</span><strong>In progress</strong></div></div>
<ol class="steps"><li class="done"><i></i>Updates started</li><li class="active"><i></i>Quality checks</li><li><i></i>Website restored</li></ol>
</aside>
</main>
<footer>&copy; 1964 - %d National Council of Sports. All Rights Reserved.</footer>
</body>
</html>`, htmlEscape(title), htmlEscape(title), htmlEscape(publicMaintenanceMessage(scoped)), htmlEscape(expected), time.Now().Year())
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
