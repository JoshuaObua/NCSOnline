package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ── Audit Context ─────────────────────────────────────────────────
// Mutable struct passed by pointer through request context so that
// downstream middlewares (e.g. Authenticate) can populate user details
// that AuditLogger reads after the handler chain completes.

type AuditContext struct {
	UserID    string
	UserEmail string
	SessionID string
}

var locationService *LocationClient
var auditJWTSecret string
var geoAllowedCountries = map[string]struct{}{"UG": {}}
var geoFailClosed bool

func ConfigureLocationService(client *LocationClient, jwtSecret string, allowedCountries []string, failClosed bool) {
	locationService = client
	auditJWTSecret = jwtSecret
	geoAllowedCountries = make(map[string]struct{}, len(allowedCountries))
	for _, country := range allowedCountries {
		if country = strings.ToUpper(strings.TrimSpace(country)); country != "" {
			geoAllowedCountries[country] = struct{}{}
		}
	}
	if len(geoAllowedCountries) == 0 {
		geoAllowedCountries["UG"] = struct{}{}
	}
	geoFailClosed = failClosed
}

type auditCtxKey struct{}

func getAuditCtx(r *http.Request) *AuditContext {
	if ac, ok := r.Context().Value(auditCtxKey{}).(*AuditContext); ok {
		return ac
	}
	return nil
}

// SetAuditIdentity attributes public authentication endpoints to the user
// established by the handler after credentials or a refresh token are valid.
func SetAuditIdentity(r *http.Request, userID, email, sessionID string) {
	if ac := getAuditCtx(r); ac != nil {
		ac.UserID = strings.TrimSpace(userID)
		ac.UserEmail = strings.ToLower(strings.TrimSpace(email))
		ac.SessionID = strings.TrimSpace(sessionID)
	}
}

// ── JWT Claims ────────────────────────────────────────────────────

type Claims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(secret string, ttl time.Duration, userID, email string, roles []string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ncsms-api",
			Subject:   userID,
			ID:        uuid.NewString(), // JTI — unique session identifier
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(secret, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer("ncsms-api"), jwt.WithIssuedAt())
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.UserID == "" || claims.IssuedAt == nil {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// parseSignedTokenIdentity recovers identity from an expired token only after
// validating its signature, algorithm and issuer. It is used for audit
// attribution and never grants access.
func parseSignedTokenIdentity(secret, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.Issuer != "ncsms-api" || claims.UserID == "" {
		return nil, fmt.Errorf("invalid token identity")
	}
	return claims, nil
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func LimitRequestBody(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil && maxBytes > 0 {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RejectAmbiguousPaths(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.URL.EscapedPath()
		decoded, err := url.PathUnescape(raw)
		if err != nil || strings.Contains(decoded, "\\") || strings.Contains(decoded, ";") ||
			strings.Contains(decoded, "//") || path.Clean(decoded) != decoded {
			response.Err(w, http.StatusBadRequest, "INVALID_PATH", "Request path is not canonical")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ── Auth Middleware ───────────────────────────────────────────────

func Authenticate(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Err(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				response.Err(w, http.StatusUnauthorized, "INVALID_TOKEN", "Bearer token required")
				return
			}

			claims, err := ParseToken(jwtSecret, parts[1])
			if err != nil {
				if auditClaims, auditErr := parseSignedTokenIdentity(jwtSecret, parts[1]); auditErr == nil {
					SetAuditIdentity(r, auditClaims.UserID, auditClaims.Email, auditClaims.ID)
				}
				response.Err(w, http.StatusUnauthorized, "TOKEN_EXPIRED", "Token is invalid or expired")
				return
			}

			ctx := context.WithValue(r.Context(), models.CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, models.CtxUserEmail, claims.Email)
			ctx = context.WithValue(ctx, models.CtxUserRoles, claims.Roles)
			ctx = context.WithValue(ctx, models.CtxSessionID, claims.ID) // JTI as session ID
			ctx = context.WithValue(ctx, models.CtxTokenIssuedAt, claims.IssuedAt.Time)

			// Propagate to the shared audit context so AuditLogger captures
			// the authenticated user even though it runs at the outer scope.
			SetAuditIdentity(r, claims.UserID, claims.Email, claims.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ValidateAuthenticatedUser prevents deleted, disabled, password-reset, or
// stale-role sessions from reaching protected handlers.
type AuthorizationStore interface {
	GetAuthorizationState(context.Context, string) (string, bool, *time.Time, []string, error)
}

func ValidateAuthenticatedUser(users AuthorizationStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, _ := r.Context().Value(models.CtxUserID).(string)
			issuedAt, _ := r.Context().Value(models.CtxTokenIssuedAt).(time.Time)
			if userID == "" || issuedAt.IsZero() {
				response.Err(w, http.StatusUnauthorized, "INVALID_SESSION", "Authenticated session is invalid")
				return
			}
			email, active, invalidBefore, roles, err := users.GetAuthorizationState(r.Context(), userID)
			if err != nil || !active {
				response.Err(w, http.StatusUnauthorized, "ACCOUNT_UNAVAILABLE", "User account is unavailable or inactive")
				return
			}
			if invalidBefore != nil && issuedAt.Before(*invalidBefore) {
				response.Err(w, http.StatusUnauthorized, "SESSION_REVOKED", "Session has been revoked. Please sign in again")
				return
			}
			ctx := context.WithValue(r.Context(), models.CtxUserEmail, email)
			ctx = context.WithValue(ctx, models.CtxUserRoles, roles)
			SetAuditIdentity(r, userID, email, r.Context().Value(models.CtxSessionID).(string))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ── RBAC Middleware ───────────────────────────────────────────────

func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles, _ := r.Context().Value(models.CtxUserRoles).([]string)
			for _, role := range userRoles {
				if _, ok := allowed[role]; ok {
					next.ServeHTTP(w, r)
					return
				}
			}
			response.Err(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
		})
	}
}

// ── Logger Middleware ─────────────────────────────────────────────

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(ww, r)
		log.Printf("[%s] %s %s %d %s",
			r.Method, r.RequestURI,
			r.RemoteAddr, ww.status,
			time.Since(start),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.status != http.StatusOK {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// ── Rate Limiter (in-memory, per IP) ─────────────────────────────

type ipEntry struct {
	count    int
	windowAt time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*ipEntry
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*ipEntry),
		limit:   limit,
		window:  window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realIP(r)
		if !rl.allow(ip) {
			w.Header().Set("Retry-After", "60")
			response.Err(w, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests — please slow down")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	e, ok := rl.entries[ip]
	if !ok || now.After(e.windowAt.Add(rl.window)) {
		rl.entries[ip] = &ipEntry{count: 1, windowAt: now}
		return true
	}
	e.count++
	return e.count <= rl.limit
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		rl.mu.Lock()
		cutoff := time.Now().Add(-rl.window)
		for ip, e := range rl.entries {
			if e.windowAt.Before(cutoff) {
				delete(rl.entries, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.TrimSpace(strings.SplitN(fwd, ",", 2)[0])
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	if host == "" {
		return r.RemoteAddr
	}
	return host
}

// ── Geo-IP Service ────────────────────────────────────────────────
// Uses ip-api.com (free, no key required, 45 req/min limit).
// Results are cached for 24h to minimise external calls.

type geoEntry struct {
	country     string
	city        string
	countryCode string
	vpn         bool
	fetchedAt   time.Time
}

var geoCache = struct {
	sync.RWMutex
	m map[string]*geoEntry
}{m: make(map[string]*geoEntry)}

type ipAPIResp struct {
	Status      string `json:"status"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	City        string `json:"city"`
	Proxy       bool   `json:"proxy"`
	Hosting     bool   `json:"hosting"`
}

// isPrivateIP returns true for loopback or RFC-1918 addresses.
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(strings.TrimSpace(ipStr))
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() {
		return true
	}
	private := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}
	for _, cidr := range private {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func lookupGeo(ipStr string) (country, city, countryCode string, vpn bool) {
	// Private / internal IPs are treated as Uganda (internal deployment)
	if isPrivateIP(ipStr) {
		return "Uganda", "Internal", "UG", false
	}
	if locationService != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		data, err := locationService.Locate(ctx, ipStr, "", "")
		if err != nil {
			return "Unknown", "", "", false
		}
		return data.Country, data.City, data.CountryCode,
			data.IsProxy || data.IsVPN || data.IsTor || data.IsHosting
	}

	geoCache.RLock()
	if e, ok := geoCache.m[ipStr]; ok && time.Since(e.fetchedAt) < 24*time.Hour {
		geoCache.RUnlock()
		return e.country, e.city, e.countryCode, e.vpn
	}
	geoCache.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode,city,proxy,hosting", ipStr)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// On failure return unknown — don't block the request
		return "Unknown", "", "", false
	}
	defer resp.Body.Close()

	var data ipAPIResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || data.Status != "success" {
		return "Unknown", "", "", false
	}

	entry := &geoEntry{
		country:     data.Country,
		city:        data.City,
		countryCode: data.CountryCode,
		vpn:         data.Proxy || data.Hosting,
		fetchedAt:   time.Now(),
	}
	geoCache.Lock()
	geoCache.m[ipStr] = entry
	geoCache.Unlock()

	return entry.country, entry.city, entry.countryCode, entry.vpn
}

// ── Geo-Blocker Middleware ────────────────────────────────────────
// Blocks requests from IP addresses outside Uganda and VPN/proxy users.
// Apply only to admin and auth routes — public website stays accessible.

func GeoBlocker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realIP(r)
		_, _, countryCode, vpn := lookupGeo(ip)

		if vpn {
			response.Err(w, http.StatusForbidden, "GEO_BLOCKED",
				"Access denied: VPN and proxy connections are not permitted on this system")
			return
		}
		// Allow private/internal IPs (countryCode="UG" from isPrivateIP) and Uganda
		if countryCode == "" && geoFailClosed {
			response.Err(w, http.StatusServiceUnavailable, "LOCATION_UNAVAILABLE", "Location verification is temporarily unavailable")
			return
		}
		if countryCode != "" && countryCode != "Unknown" {
			if _, allowed := geoAllowedCountries[strings.ToUpper(countryCode)]; allowed {
				next.ServeHTTP(w, r)
				return
			}
			response.Err(w, http.StatusForbidden, "GEO_BLOCKED",
				"Access denied: this system is only accessible from Uganda")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ── Audit Logger Middleware ───────────────────────────────────────
// Logs every HTTP request with comprehensive security context.
// Fire-and-forget goroutine — never blocks the HTTP response.

func AuditLogger(repo *repository.AuditRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// Install shared audit context so downstream middlewares
			// (Authenticate) can populate user details on it.
			ac := &AuditContext{}
			if auth := strings.TrimSpace(r.Header.Get("Authorization")); auth != "" && auditJWTSecret != "" {
				parts := strings.SplitN(auth, " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
					if claims, err := parseSignedTokenIdentity(auditJWTSecret, parts[1]); err == nil {
						ac.UserID, ac.UserEmail, ac.SessionID = claims.UserID, claims.Email, claims.ID
					}
				}
			}
			ctx := context.WithValue(r.Context(), auditCtxKey{}, ac)
			r = r.WithContext(ctx)

			next.ServeHTTP(ww, r)
			elapsed := time.Since(start).Milliseconds()

			// Read user details that downstream middlewares populated
			userID := ac.UserID
			userEmail := ac.UserEmail
			sessionID := ac.SessionID
			method := r.Method
			endpoint := r.URL.Path
			ip := realIP(r)
			forwardedFor := r.Header.Get("X-Forwarded-For")
			ua := r.UserAgent()
			statusCode := ww.status

			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				geoCountry, geoCity, _, vpnDetected := lookupGeo(ip)
				var location *LocationResult
				if locationService != nil {
					location, _ = locationService.Locate(ctx, ip, ua, userID)
					if location != nil {
						geoCountry, geoCity = location.Country, location.City
						vpnDetected = location.IsProxy || location.IsVPN || location.IsTor || location.IsHosting
					}
				}

				// Classify the request
				eventType := classifyEventType(method, endpoint)
				eventStatus := classifyEventStatus(statusCode)
				threatScore := computeThreatScore(geoCountry, vpnDetected, statusCode)
				severityLevel := classifySeverity(eventStatus, threatScore, vpnDetected)
				anomaly := threatScore >= 60

				browser, osName := parseBrowserOS(ua)
				clientType := parseClientType(ua)
				deviceInfo := osName + " / " + browser

				var uid *string
				if userID != "" {
					uid = &userID
				}

				entry := &models.AuditLog{
					ID:              uuid.NewString(),
					UserID:          uid,
					Username:        userEmail,
					SessionID:       sessionID,
					Action:          method + " " + endpoint,
					Resource:        extractResource(endpoint),
					Method:          method,
					Endpoint:        endpoint,
					IPAddress:       ip,
					ForwardedIP:     forwardedFor,
					UserAgent:       ua,
					Browser:         browser,
					OSName:          osName,
					ClientType:      clientType,
					DeviceInfo:      deviceInfo,
					ResponseCode:    statusCode,
					ResponseTimeMs:  elapsed,
					EventType:       eventType,
					EventStatus:     eventStatus,
					SeverityLevel:   severityLevel,
					GeoCountry:      geoCountry,
					GeoCity:         geoCity,
					VPNDetected:     vpnDetected,
					ThreatScore:     threatScore,
					AnomalyDetected: anomaly,
				}
				if location != nil {
					entry.GeoRegion = location.Region
					entry.GeoLatitude = location.Latitude
					entry.GeoLongitude = location.Longitude
					entry.GeoTimezone = location.Timezone
					entry.GeoSource = location.Source
					entry.Platform = location.Platform
					entry.Authenticated = userID != ""
					if location.Browser != "" {
						entry.Browser = location.Browser
					}
					if location.DeviceType != "" {
						entry.ClientType = location.DeviceType
					}
				}
				if err := repo.Log(ctx, entry); err != nil {
					log.Printf("[audit] failed to log request: %v", err)
				}
			}()
		})
	}
}

// ── Classification Helpers ────────────────────────────────────────

func classifyEventType(method, path string) string {
	p := strings.ToLower(path)
	switch {
	case strings.Contains(p, "/auth/login"):
		return "AUTH_LOGIN"
	case strings.Contains(p, "/auth/logout"):
		return "AUTH_LOGOUT"
	case strings.Contains(p, "/auth/register"):
		return "AUTH_REGISTER"
	case strings.Contains(p, "/auth/forgot-password"):
		return "AUTH_FORGOT_PASSWORD"
	case strings.Contains(p, "/auth/reset-password"):
		return "AUTH_RESET_PASSWORD"
	case strings.Contains(p, "/auth/change-password"):
		return "AUTH_PASSWORD_CHANGE"
	case strings.Contains(p, "/auth/pin"):
		return "AUTH_PIN_CHANGE"
	case strings.Contains(p, "/auth/me"):
		return "AUTH_PROFILE_VIEW"
	case strings.Contains(p, "/auth/refresh"):
		return "AUTH_TOKEN_REFRESH"
	case strings.Contains(p, "/admin/users") && method == http.MethodGet:
		return "USER_LIST"
	case strings.Contains(p, "/admin/users") && method == http.MethodPost:
		if strings.Contains(p, "/reset-password") {
			return "USER_PASSWORD_RESET"
		}
		if strings.Contains(p, "/activate") {
			return "USER_ACTIVATE"
		}
		if strings.Contains(p, "/deactivate") {
			return "USER_DEACTIVATE"
		}
		if strings.Contains(p, "/roles") {
			return "ROLE_ASSIGN"
		}
		return "USER_CREATE"
	case strings.Contains(p, "/admin/users") && method == http.MethodPut:
		return "USER_UPDATE"
	case strings.Contains(p, "/admin/users") && method == http.MethodDelete:
		if strings.Contains(p, "/roles/") {
			return "ROLE_REMOVE"
		}
		return "USER_DELETE"
	case strings.Contains(p, "/admin/audit-logs"):
		return "AUDIT_VIEW"
	case strings.Contains(p, "/admin/applications") && method == http.MethodGet:
		return "APPLICATION_VIEW"
	case strings.Contains(p, "/approve"):
		return "APPLICATION_APPROVE"
	case strings.Contains(p, "/reject"):
		return "APPLICATION_REJECT"
	case strings.Contains(p, "/verify-payment"):
		return "PAYMENT_VERIFY"
	case strings.Contains(p, "/admin/cms") && method == http.MethodPost:
		return "CMS_CREATE"
	case strings.Contains(p, "/admin/cms") && method == http.MethodPut:
		return "CMS_UPDATE"
	case strings.Contains(p, "/admin/cms") && method == http.MethodDelete:
		return "CMS_DELETE"
	case strings.Contains(p, "/admin/roles") && method == http.MethodPost:
		return "ROLE_CREATE"
	case strings.Contains(p, "/admin/roles") && method == http.MethodPut:
		return "ROLE_UPDATE"
	case strings.Contains(p, "/admin/roles") && method == http.MethodDelete:
		return "ROLE_DELETE"
	case strings.Contains(p, "/admin/dashboard"):
		return "DASHBOARD_VIEW"
	default:
		switch method {
		case http.MethodGet:
			return "VIEW"
		case http.MethodPost:
			return "CREATE"
		case http.MethodPut, http.MethodPatch:
			return "UPDATE"
		case http.MethodDelete:
			return "DELETE"
		}
		return "ACTION"
	}
}

func classifyEventStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "SUCCESS"
	case code == 401:
		return "AUTH_FAILED"
	case code == 403:
		return "DENIED"
	case code >= 400 && code < 500:
		return "FAILED"
	case code >= 500:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func computeThreatScore(country string, vpn bool, statusCode int) int {
	score := 0
	if country != "" && country != "Uganda" && country != "Internal" && country != "Unknown" {
		score += 40
	}
	if vpn {
		score += 30
	}
	if statusCode == http.StatusUnauthorized {
		score += 20
	}
	if statusCode == http.StatusForbidden {
		score += 10
	}
	if statusCode >= 500 {
		score += 10
	}
	if score > 100 {
		score = 100
	}
	return score
}

func classifySeverity(eventStatus string, threatScore int, vpn bool) string {
	if threatScore >= 70 || (vpn && eventStatus != "SUCCESS") {
		return "CRITICAL"
	}
	if threatScore >= 40 || eventStatus == "DENIED" || eventStatus == "AUTH_FAILED" || vpn {
		return "WARNING"
	}
	if eventStatus == "ERROR" {
		return "ERROR"
	}
	return "INFO"
}

// parseBrowserOS returns (browser, osName) from a User-Agent string.
func parseBrowserOS(ua string) (browser, osName string) {
	lower := strings.ToLower(ua)

	// OS detection
	switch {
	case strings.Contains(lower, "windows nt 10") || strings.Contains(lower, "windows 10"):
		osName = "Windows 10"
	case strings.Contains(lower, "windows nt 11") || strings.Contains(lower, "windows 11"):
		osName = "Windows 11"
	case strings.Contains(lower, "windows"):
		osName = "Windows"
	case strings.Contains(lower, "android"):
		osName = "Android"
	case strings.Contains(lower, "iphone"):
		osName = "iOS (iPhone)"
	case strings.Contains(lower, "ipad"):
		osName = "iOS (iPad)"
	case strings.Contains(lower, "mac os x") || strings.Contains(lower, "macos"):
		osName = "macOS"
	case strings.Contains(lower, "linux"):
		osName = "Linux"
	case strings.Contains(lower, "ubuntu"):
		osName = "Ubuntu"
	default:
		osName = "Unknown"
	}

	// Browser/client detection
	switch {
	case strings.Contains(lower, "postman"):
		browser = "Postman"
	case strings.Contains(lower, "insomnia"):
		browser = "Insomnia"
	case strings.Contains(lower, "thunder-client") || strings.Contains(lower, "thunderclient"):
		browser = "Thunder Client"
	case strings.Contains(lower, "httpie"):
		browser = "HTTPie"
	case strings.Contains(lower, "python-requests"):
		browser = "Python/requests"
	case strings.Contains(lower, "go-http-client"):
		browser = "Go/HTTP"
	case strings.Contains(lower, "curl"):
		browser = "curl"
	case strings.Contains(lower, "edg/"):
		browser = "Edge"
	case strings.Contains(lower, "opr/") || strings.Contains(lower, "opera"):
		browser = "Opera"
	case strings.Contains(lower, "chrome") && !strings.Contains(lower, "chromium"):
		browser = "Chrome"
	case strings.Contains(lower, "chromium"):
		browser = "Chromium"
	case strings.Contains(lower, "firefox"):
		browser = "Firefox"
	case strings.Contains(lower, "safari") && !strings.Contains(lower, "chrome"):
		browser = "Safari"
	case strings.Contains(lower, "brave"):
		browser = "Brave"
	default:
		if ua == "" {
			browser = "Unknown"
		} else {
			browser = "Other"
		}
	}

	return browser, osName
}

// parseClientType classifies the client as Browser, API Tool, or programmatic client.
func parseClientType(ua string) string {
	lower := strings.ToLower(ua)
	switch {
	case strings.Contains(lower, "postman"):
		return "Postman"
	case strings.Contains(lower, "insomnia"):
		return "Insomnia"
	case strings.Contains(lower, "thunder-client") || strings.Contains(lower, "thunderclient"):
		return "Thunder Client"
	case strings.Contains(lower, "httpie"):
		return "HTTPie"
	case strings.Contains(lower, "python-requests"):
		return "Python/requests"
	case strings.Contains(lower, "go-http-client"):
		return "Go/HTTP"
	case strings.Contains(lower, "curl"):
		return "curl"
	case strings.Contains(lower, "mozilla") || strings.Contains(lower, "webkit"):
		return "Browser"
	case ua == "":
		return "Unknown"
	default:
		return "Other"
	}
}

// extractResource pulls the first meaningful path segment (e.g. "auth", "users", "cms").
func extractResource(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	for _, p := range parts {
		if p != "" && p != "api" && p != "v1" && p != "admin" {
			return p
		}
	}
	return "unknown"
}

// parseDeviceInfo is kept for backwards compatibility.
func parseDeviceInfo(ua string) string {
	browser, osName := parseBrowserOS(ua)
	return osName + " / " + browser
}
