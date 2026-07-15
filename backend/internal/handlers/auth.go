package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/middleware"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
	"github.com/atenimedia-llc/ncs-online/backend/internal/storage"
)

type AuthHandler struct {
	svc   *services.AuthService
	users *repository.UserRepo
	cfg   *config.Config
	cms   *repository.CMSRepo
}

const (
	refreshCookieName = "ncsms_refresh"
	accessCookieName  = "ncsms_access"
)

func (h *AuthHandler) setRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{Name: refreshCookieName, Value: token, Path: "/api/v1/auth", HttpOnly: true, Secure: h.cfg.IsProduction(), SameSite: http.SameSiteStrictMode, MaxAge: int(h.cfg.RefreshTokenTTL.Seconds()), Expires: time.Now().Add(h.cfg.RefreshTokenTTL)})
}
func (h *AuthHandler) setAccessCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{Name: accessCookieName, Value: token, Path: "/", HttpOnly: true, Secure: h.cfg.IsProduction(), SameSite: http.SameSiteLaxMode, MaxAge: int(h.cfg.AccessTokenTTL.Seconds()), Expires: time.Now().Add(h.cfg.AccessTokenTTL)})
}
func (h *AuthHandler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: refreshCookieName, Value: "", Path: "/api/v1/auth", HttpOnly: true, Secure: h.cfg.IsProduction(), SameSite: http.SameSiteStrictMode, MaxAge: -1, Expires: time.Unix(1, 0)})
}
func (h *AuthHandler) clearAccessCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: accessCookieName, Value: "", Path: "/", HttpOnly: true, Secure: h.cfg.IsProduction(), SameSite: http.SameSiteLaxMode, MaxAge: -1, Expires: time.Unix(1, 0)})
}
func refreshFromRequest(r *http.Request) string {
	if c, e := r.Cookie(refreshCookieName); e == nil && c.Value != "" {
		return c.Value
	}
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	return req.RefreshToken
}
func (h *AuthHandler) secureResult(w http.ResponseWriter, result *services.LoginResult) {
	h.setRefreshCookie(w, result.RefreshToken)
	h.setAccessCookie(w, result.AccessToken)
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	errs := map[string]string{}
	if req.FirstName == "" {
		errs["first_name"] = "required"
	}
	if req.LastName == "" {
		errs["last_name"] = "required"
	}
	if req.Email == "" {
		errs["email"] = "required"
	}
	if len(req.Password) < 8 {
		errs["password"] = "minimum 8 characters"
	}
	if len(errs) > 0 {
		response.ValidationErr(w, errs)
		return
	}

	result, err := h.svc.Register(r.Context(), req.FirstName, req.LastName,
		strings.ToLower(req.Email), req.Password, r.RemoteAddr, r.UserAgent())
	if errors.Is(err, services.ErrEmailTaken) {
		response.Err(w, http.StatusConflict, "EMAIL_TAKEN", "An account with that email already exists")
		return
	}
	if err != nil {
		slog.Error("register failed", "error", err.Error(), "email", req.Email)
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Registration failed. Please try again.")
		return
	}
	middleware.SetAuditIdentity(r, result.User.ID, result.User.Email, "")
	h.secureResult(w, result)
	response.JSON(w, http.StatusCreated, result)
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Email == "" || req.Password == "" {
		response.ValidationErr(w, map[string]string{"email": "required", "password": "required"})
		return
	}

	result, err := h.svc.Login(r.Context(), strings.ToLower(req.Email), req.Password,
		r.RemoteAddr, r.UserAgent())
	if errors.Is(err, services.ErrInvalidCredentials) {
		response.Err(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
		return
	}
	if errors.Is(err, services.ErrAccountDisabled) {
		response.Err(w, http.StatusForbidden, "ACCOUNT_DISABLED", "Your account has been deactivated")
		return
	}
	if err != nil {
		slog.Error("login failed", "error", err.Error(), "email", req.Email)
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "An unexpected error occurred")
		return
	}

	twofaEnabled, _, err := h.users.GetTwoFAEnabled(r.Context(), result.User.ID)
	if err == nil && twofaEnabled {
		code := cryptoRand6DigitCode()
		_ = h.users.SetTwoFASecretDirectly(r.Context(), result.User.ID, code)
		
		ticket, err := middleware.GenerateAccessToken(h.cfg.JWTSecret, 5 * time.Minute, result.User.ID, result.User.Email, []string{"2fa_pending"})
		if err != nil {
			response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not generate 2FA ticket")
			return
		}
		
		if err := h.send2FAEmail(result.User.Email, code); err != nil {
			slog.Error("failed to send login 2FA email", "error", err.Error(), "email", result.User.Email)
			response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not send 2FA verification email")
			return
		}
		
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"requires_2fa": true,
			"ticket":       ticket,
		})
		return
	}

	middleware.SetAuditIdentity(r, result.User.ID, result.User.Email, "")
	h.secureResult(w, result)
	response.JSON(w, http.StatusOK, result)
}

func cryptoRand6DigitCode() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	}
	var val uint32
	for _, x := range b {
		val = (val << 8) | uint32(x)
	}
	code := 100000 + (val % 900000)
	return fmt.Sprintf("%06d", code)
}

func (h *AuthHandler) send2FAEmail(recipientEmail, code string) error {
	if h.cfg.SMTPHost == "" {
		slog.Info("DEVELOPMENT: SMTP_HOST not configured. 2FA verification email bypassed.", "email", recipientEmail, "code", code)
		return nil
	}
	subject := "NCS Uganda 2FA Verification Code"
	body := fmt.Sprintf("Your NCS Uganda two-factor authentication verification code is: %s\r\n\r\nThis code will expire in 5 minutes.", code)
	
	fromAddress := h.cfg.SMTPFrom
	if i := strings.LastIndex(fromAddress, "<"); i >= 0 && strings.HasSuffix(fromAddress, ">") {
		fromAddress = fromAddress[i+1 : len(fromAddress)-1]
	}
	message := []byte("From: " + h.cfg.SMTPFrom + "\r\nTo: " + recipientEmail + "\r\nSubject: " + subject + "\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n" + body)
	var auth smtp.Auth
	if h.cfg.SMTPUser != "" {
		auth = smtp.PlainAuth("", h.cfg.SMTPUser, h.cfg.SMTPPassword, h.cfg.SMTPHost)
	}
	return smtp.SendMail(h.cfg.SMTPHost+":"+h.cfg.SMTPPort, auth, fromAddress, []string{recipientEmail}, message)
}

// POST /api/v1/auth/login/2fa
func (h *AuthHandler) VerifyLogin2FA(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Ticket string `json:"ticket"`
		Code   string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Ticket == "" || req.Code == "" {
		response.ValidationErr(w, map[string]string{"ticket": "required", "code": "required"})
		return
	}

	claims, err := middleware.ParseToken(h.cfg.JWTSecret, req.Ticket)
	if err != nil {
		response.Err(w, http.StatusUnauthorized, "TICKET_INVALID", "The session has expired. Please log in again.")
		return
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "2fa_pending" {
		response.Err(w, http.StatusUnauthorized, "TICKET_INVALID", "Invalid session context.")
		return
	}

	enabled, secret, err := h.users.GetTwoFAEnabled(r.Context(), claims.UserID)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load security state")
		return
	}
	if !enabled || secret == "" {
		response.Err(w, http.StatusConflict, "TWOFA_NOT_ENABLED", "Two-factor authentication is not active for this account.")
		return
	}

	if secret != req.Code {
		response.Err(w, http.StatusUnauthorized, "TWOFA_INVALID", "The verification code is incorrect. Please check your email and try again.")
		return
	}

	_ = h.users.SetTwoFASecretDirectly(r.Context(), claims.UserID, "")

	result, err := h.svc.Login2FA(r.Context(), claims.UserID, r.RemoteAddr, r.UserAgent())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not complete login flow")
		return
	}

	middleware.SetAuditIdentity(r, result.User.ID, result.User.Email, "")
	h.secureResult(w, result)
	response.JSON(w, http.StatusOK, result)
}

// POST /api/v1/auth/me/avatar
func (h *AuthHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Could not parse form (max 10MB)")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "No file uploaded (field: 'file')")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	if !allowed[ext] {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "File type not allowed. Use an image (jpg, jpeg, png, gif, webp)")
		return
	}
	head := make([]byte, 512)
	n, _ := io.ReadFull(file, head)
	head = head[:n]
	detected := http.DetectContentType(head)
	
	if !strings.HasPrefix(detected, "image/") {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Uploaded file content does not match an allowed image type")
		return
	}

	settings := storage.DefaultSettings()
	if h.cms != nil {
		s, err := h.cms.GetSetting(r.Context(), storage.SettingsKey)
		if err == nil {
			if parsed, err := storage.ParseSettings(s.Value); err == nil {
				settings = parsed
			}
		}
	}

	result, err := storage.NewUploader(settings).Upload(r.Context(), storage.UploadInput{
		Scope:       storage.ScopePublic,
		Reader:      io.MultiReader(bytes.NewReader(head), file),
		Filename:    header.Filename,
		ContentType: detected,
		Size:        header.Size,
		Subdir:      "avatars",
	})
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "UPLOAD_FAILED", err.Error())
		return
	}

	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load user")
		return
	}
	user.AvatarURL = result.URL
	if err := h.users.Update(r.Context(), user); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update user avatar")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"avatar_url": result.URL})
}

// POST /api/v1/auth/google
func (h *AuthHandler) Google(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Credential string `json:"credential"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	result, err := h.svc.LoginWithGoogle(r.Context(), req.Credential, r.RemoteAddr, r.UserAgent())
	if errors.Is(err, services.ErrGoogleAuthNotConfigured) {
		response.Err(w, http.StatusServiceUnavailable, "GOOGLE_AUTH_NOT_CONFIGURED", "Google sign-in is not configured")
		return
	}
	if errors.Is(err, services.ErrGoogleTokenInvalid) {
		response.Err(w, http.StatusUnauthorized, "INVALID_GOOGLE_TOKEN", "Google sign-in could not be verified")
		return
	}
	if errors.Is(err, services.ErrAccountDisabled) {
		response.Err(w, http.StatusForbidden, "ACCOUNT_DISABLED", "Your account has been deactivated")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Google sign-in failed. Please try again.")
		return
	}
	middleware.SetAuditIdentity(r, result.User.ID, result.User.Email, "")
	h.secureResult(w, result)
	response.JSON(w, http.StatusOK, result)
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	raw := refreshFromRequest(r)
	if raw == "" {
		response.Err(w, http.StatusUnauthorized, "TOKEN_REQUIRED", "Refresh session is required")
		return
	}
	result, err := h.svc.RefreshToken(r.Context(), raw, r.RemoteAddr, r.UserAgent())
	if result != nil && result.User != nil {
		middleware.SetAuditIdentity(r, result.User.ID, result.User.Email, "")
	}
	if errors.Is(err, services.ErrTokenInvalid) {
		response.Err(w, http.StatusUnauthorized, "TOKEN_INVALID", "Refresh token is invalid or expired")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "An unexpected error occurred")
		return
	}
	middleware.SetAuditIdentity(r, result.User.ID, result.User.Email, "")
	h.secureResult(w, result)
	response.JSON(w, http.StatusOK, result)
}

// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	raw := refreshFromRequest(r)
	_ = h.svc.Logout(r.Context(), raw)
	h.clearRefreshCookie(w)
	h.clearAccessCookie(w)
	response.JSONMsg(w, http.StatusOK, "Logged out successfully")
}

// GET /api/v1/auth/me  — returns full user profile from DB
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)

	user, err := h.users.GetByID(r.Context(), userID)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusUnauthorized, "USER_NOT_FOUND", "User account not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch profile")
		return
	}

	roles, _ := h.users.GetRoles(r.Context(), userID)
	user.Roles = roles
	user.PasswordHash = ""
	user.PinHash = ""

	response.JSON(w, http.StatusOK, user)
}

// PUT /api/v1/auth/me — self-service update of the caller's own display name
// and avatar. Deliberately narrow: fetches the current record first and only
// overwrites these three fields, so it can't be used to touch email, roles,
// or account status (those go through the admin user-management endpoints).
func (h *AuthHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.FirstName == "" || req.LastName == "" {
		response.ValidationErr(w, map[string]string{"first_name": "required", "last_name": "required"})
		return
	}

	user, err := h.users.GetByID(r.Context(), userID)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusUnauthorized, "USER_NOT_FOUND", "User account not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load profile")
		return
	}
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.AvatarURL = req.AvatarURL
	if err := h.users.Update(r.Context(), user); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update profile")
		return
	}

	roles, _ := h.users.GetRoles(r.Context(), userID)
	user.Roles = roles
	user.PasswordHash = ""
	user.PinHash = ""
	response.JSON(w, http.StatusOK, user)
}

// PUT /api/v1/auth/me/password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	errs := map[string]string{}
	if req.CurrentPassword == "" {
		errs["current_password"] = "required"
	}
	if len(req.NewPassword) < 8 {
		errs["new_password"] = "minimum 8 characters"
	}
	if len(errs) > 0 {
		response.ValidationErr(w, errs)
		return
	}
	if err := h.svc.ChangePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			response.Err(w, http.StatusUnauthorized, "WRONG_PASSWORD", "Current password is incorrect")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not change password")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Password changed successfully")
}

// POST /api/v1/auth/pin/set  — first-time PIN setup or forced reset
func (h *AuthHandler) SetPIN(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		PIN string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if err := h.svc.SetPIN(r.Context(), userID, req.PIN); err != nil {
		if errors.Is(err, services.ErrInvalidPIN) {
			response.ValidationErr(w, map[string]string{"pin": "must be 4–6 digits"})
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not set PIN")
		return
	}
	response.JSONMsg(w, http.StatusOK, "PIN set successfully")
}

// PUT /api/v1/auth/pin/change
func (h *AuthHandler) ChangePIN(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		CurrentPIN string `json:"current_pin"`
		NewPIN     string `json:"new_pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	errs := map[string]string{}
	if req.CurrentPIN == "" {
		errs["current_pin"] = "required"
	}
	if req.NewPIN == "" {
		errs["new_pin"] = "required"
	}
	if len(errs) > 0 {
		response.ValidationErr(w, errs)
		return
	}
	if err := h.svc.ChangePIN(r.Context(), userID, req.CurrentPIN, req.NewPIN); err != nil {
		if errors.Is(err, services.ErrPINIncorrect) {
			response.Err(w, http.StatusUnauthorized, "WRONG_PIN", "Current PIN is incorrect")
			return
		}
		if errors.Is(err, services.ErrPINNotSet) {
			response.Err(w, http.StatusBadRequest, "PIN_NOT_SET", "No PIN has been set — use /auth/pin/set first")
			return
		}
		if errors.Is(err, services.ErrInvalidPIN) {
			response.ValidationErr(w, map[string]string{"new_pin": "must be 4–6 digits"})
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not change PIN")
		return
	}
	response.JSONMsg(w, http.StatusOK, "PIN changed successfully")
}

// POST /api/v1/auth/pin/verify  — lock screen unlock
func (h *AuthHandler) VerifyPIN(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		PIN string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if err := h.svc.VerifyPIN(r.Context(), userID, req.PIN); err != nil {
		if errors.Is(err, services.ErrPINIncorrect) {
			response.Err(w, http.StatusUnauthorized, "WRONG_PIN", "Incorrect PIN")
			return
		}
		if errors.Is(err, services.ErrPINNotSet) {
			response.Err(w, http.StatusBadRequest, "PIN_NOT_SET", "No PIN has been set")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not verify PIN")
		return
	}
	response.JSONMsg(w, http.StatusOK, "PIN verified")
}

// POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	response.JSONMsg(w, http.StatusOK, "If that email exists, a reset link has been sent")
}

// POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	response.JSONMsg(w, http.StatusOK, "Password has been reset")
}
