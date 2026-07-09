package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func TestParseTokenRequiresHS256AndIssuer(t *testing.T) {
	const secret = "a-test-secret-that-is-long-enough-for-unit-tests"
	valid, err := GenerateAccessToken(secret, time.Minute, "user-1", "user@example.com", []string{"admin"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseToken(secret, valid); err != nil {
		t.Fatalf("valid token rejected: %v", err)
	}

	claims := &Claims{UserID: "user-1", RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "wrong-issuer",
	}}
	wrongIssuer, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseToken(secret, wrongIssuer); err == nil {
		t.Fatal("token with wrong issuer was accepted")
	}
}

func TestSecurityHeaders(t *testing.T) {
	h := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/health", nil))

	for name, want := range map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Cache-Control":          "no-store",
	} {
		if got := rr.Header().Get(name); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}
}

func TestSecurityHeadersIncludeHSTSForForwardedHTTPS(t *testing.T) {
	h := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/cms/posts", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if got := rr.Header().Get("Strict-Transport-Security"); !strings.Contains(got, "max-age=31536000") {
		t.Fatalf("Strict-Transport-Security = %q", got)
	}
	if got := rr.Header().Get("Cross-Origin-Resource-Policy"); got != "same-site" {
		t.Fatalf("Cross-Origin-Resource-Policy = %q, want same-site", got)
	}
}

func TestSetAuditIdentity(t *testing.T) {
	ac := &AuditContext{}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	req = req.WithContext(context.WithValue(req.Context(), auditCtxKey{}, ac))

	SetAuditIdentity(req, " user-1 ", " ADMIN@EXAMPLE.COM ", " session-1 ")

	if ac.UserID != "user-1" || ac.UserEmail != "admin@example.com" || ac.SessionID != "session-1" {
		t.Fatalf("unexpected audit identity: %#v", ac)
	}
}

func TestSignedExpiredTokenIdentityIsAvailableForAuditOnly(t *testing.T) {
	const secret = "a-test-secret-that-is-long-enough-for-unit-tests"
	claims := &Claims{
		UserID: "user-1",
		Email:  "user@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
			Issuer:    "ncsms-api",
			ID:        "session-1",
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseToken(secret, token); err == nil {
		t.Fatal("expired token unexpectedly authorized")
	}
	auditClaims, err := parseSignedTokenIdentity(secret, token)
	if err != nil {
		t.Fatalf("signed expired token was not available for audit: %v", err)
	}
	if auditClaims.UserID != "user-1" || auditClaims.ID != "session-1" {
		t.Fatalf("unexpected audit claims: %#v", auditClaims)
	}
}

func TestRejectAmbiguousPaths(t *testing.T) {
	h := RejectAmbiguousPaths(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for _, target := range []string{"/api//v1/users", "/api/v1/admin;test", "/api/v1/../admin"} {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, target, nil))
		if rr.Code != http.StatusBadRequest {
			t.Errorf("%s returned %d, want 400", target, rr.Code)
		}
	}
}

func TestLimitRequestBody(t *testing.T) {
	h := LimitRequestBody(4)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.ReadAll(r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("12345")))
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rr.Code)
	}
}

func TestLimitRequestBodyByTypeUsesSmallerJSONBudget(t *testing.T) {
	h := LimitRequestBodyByType(4, 16)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
	}))

	jsonRequest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("12345"))
	jsonRequest.Header.Set("Content-Type", "application/json")
	jsonResponse := httptest.NewRecorder()
	h.ServeHTTP(jsonResponse, jsonRequest)
	if jsonResponse.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("JSON status = %d, want 413", jsonResponse.Code)
	}

	uploadRequest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("12345"))
	uploadRequest.Header.Set("Content-Type", "multipart/form-data; boundary=test")
	uploadResponse := httptest.NewRecorder()
	h.ServeHTTP(uploadResponse, uploadRequest)
	if uploadResponse.Code != http.StatusNoContent {
		t.Fatalf("upload status = %d, want 204", uploadResponse.Code)
	}
}

func TestClassifyRequestClientDetectsUSSDAndAPITools(t *testing.T) {
	ussd := httptest.NewRequest(http.MethodPost, "/ncs-ussd", nil)
	if got := classifyRequestClient(ussd); got != "USSD Gateway" {
		t.Fatalf("USSD client = %q", got)
	}

	postman := httptest.NewRequest(http.MethodGet, "/api/v1/cms/posts", nil)
	postman.Header.Set("User-Agent", "PostmanRuntime/7.43.0")
	if got := classifyRequestClient(postman); got != "Postman" {
		t.Fatalf("Postman client = %q", got)
	}

	integration := httptest.NewRequest(http.MethodPost, "/api/v1/applications", nil)
	integration.Header.Set("X-NCS-Channel", "integration")
	if got := classifyRequestClient(integration); got != "API Client" {
		t.Fatalf("integration client = %q", got)
	}
}

func TestDeviceFingerprintPreservesClientAndDeviceDetails(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 Chrome/126.0 Mobile Safari/537.36"
	browser, platform := parseBrowserOS(ua)
	info := formatDeviceInfo(platform, parseDeviceType(ua), browser, parseClientType(ua))
	for _, want := range []string{"Android", "Mobile", "Chrome", "Browser"} {
		if !strings.Contains(info, want) {
			t.Fatalf("device info %q is missing %q", info, want)
		}
	}
}

func TestRequireRoles(t *testing.T) {
	h := RequireRoles("admin")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/admin", nil)
	request = request.WithContext(context.WithValue(request.Context(), models.CtxUserRoles, []string{"applicant"}))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, request)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rr.Code)
	}
}

func TestAuthenticateRejectsMissingToken(t *testing.T) {
	h := Authenticate("a-test-secret-that-is-long-enough-for-unit-tests")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rr.Code)
	}
}

type fakeAuthorizationStore struct {
	email         string
	active        bool
	invalidBefore *time.Time
	roles         []string
	err           error
}

func (f fakeAuthorizationStore) GetAuthorizationState(context.Context, string) (string, bool, *time.Time, []string, error) {
	return f.email, f.active, f.invalidBefore, f.roles, f.err
}

func authenticatedRequest(roles []string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	ctx := context.WithValue(req.Context(), models.CtxUserID, "user-1")
	ctx = context.WithValue(ctx, models.CtxUserEmail, "old@example.com")
	ctx = context.WithValue(ctx, models.CtxUserRoles, roles)
	ctx = context.WithValue(ctx, models.CtxSessionID, "session-1")
	ctx = context.WithValue(ctx, models.CtxTokenIssuedAt, time.Now().Add(-time.Minute))
	return req.WithContext(ctx)
}

func TestValidateAuthenticatedUserRejectsInactiveAccount(t *testing.T) {
	h := ValidateAuthenticatedUser(fakeAuthorizationStore{active: false})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, authenticatedRequest([]string{"admin"}))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rr.Code)
	}
}

func TestValidateAuthenticatedUserReloadsCurrentRoles(t *testing.T) {
	h := ValidateAuthenticatedUser(fakeAuthorizationStore{
		email:  "user@example.com",
		active: true,
		roles:  []string{"applicant"},
	})(RequireRoles("admin")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, authenticatedRequest([]string{"admin"}))
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 after admin role removal", rr.Code)
	}
}

func TestValidateAuthenticatedUserRejectsRevokedAccessToken(t *testing.T) {
	invalidBefore := time.Now()
	h := ValidateAuthenticatedUser(fakeAuthorizationStore{
		email:         "user@example.com",
		active:        true,
		invalidBefore: &invalidBefore,
		roles:         []string{"admin"},
	})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, authenticatedRequest([]string{"admin"}))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401 for revoked session", rr.Code)
	}
}
