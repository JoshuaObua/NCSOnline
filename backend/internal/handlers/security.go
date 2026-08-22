package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
	"github.com/go-chi/chi/v5"
)

type SecurityHandler struct {
	svc    *services.SecurityService
	users  *repository.UserRepo
	audit  *repository.AuditRepo
	tokens *repository.TokenRepo
	cfg    *config.Config
}

// GET /api/v1/me/activities
func (h *SecurityHandler) MyActivities(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	if page < 1 {
		page = 1
	}
	if perPage <= 0 || perPage > 200 {
		perPage = 25
	}
	p := &models.PaginationParams{Page: page, PerPage: perPage}
	logs, total, err := h.audit.ListByUser(r.Context(), userID, p)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSONPaged(w, http.StatusOK, logs, &response.Meta{Page: p.Page, PerPage: p.PerPage, Total: total})
}

// GET /api/v1/me/security
// Returns the user's IP allowlist + 2FA state in one envelope.
func (h *SecurityHandler) GetMySecurity(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	wl, err := h.svc.ListWhitelist(r.Context(), userID)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	twofa, err := h.svc.GetTwoFAState(r.Context(), userID)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{
		"ip_whitelist": wl,
		"twofa":        twofa,
	})
}

// POST /api/v1/me/security/ip-whitelist  { ip_or_cidr, label }
func (h *SecurityHandler) AddIPWhitelist(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		IPOrCIDR string `json:"ip_or_cidr"`
		Label    string `json:"label"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	entry, err := h.svc.AddWhitelist(r.Context(), userID, req.IPOrCIDR, req.Label)
	if err != nil {
		if errors.Is(err, services.ErrInvalidIPEntry) {
			response.Err(w, http.StatusBadRequest, "BAD_IP", "Provide a valid IPv4/IPv6 address or CIDR block (e.g. 41.74.32.5 or 192.168.0.0/24)")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, entry)
}

// DELETE /api/v1/me/security/ip-whitelist/{id}
func (h *SecurityHandler) RemoveIPWhitelist(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	if err := h.svc.RemoveWhitelist(r.Context(), userID, id); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Removed")
}

func (h *SecurityHandler) send2FAEmail(recipientEmail, code string, purpose string) error {
	if h.cfg.SMTPHost == "" {
		slog.Info("DEVELOPMENT: SMTP_HOST not configured. 2FA activation code email bypassed.", "email", recipientEmail, "code", code, "purpose", purpose)
		return nil
	}
	subject := "NCS Uganda 2FA Code"
	var body string
	if purpose == "enroll" {
		subject = "NCS Uganda 2FA Activation Code"
		body = fmt.Sprintf("Use this code to activate two-factor authentication (2FA) for your account: %s\r\n\r\nIf you did not request this, please ignore this email.", code)
	} else {
		subject = "NCS Uganda 2FA Verification Code"
		body = fmt.Sprintf("Your NCS Uganda two-factor authentication verification code is: %s\r\n\r\nThis code will expire in 5 minutes.", code)
	}
	
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

// POST /api/v1/me/security/2fa/enroll
func (h *SecurityHandler) Enroll2FA(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	email, _ := r.Context().Value(models.CtxUserEmail).(string)
	if email == "" {
		u, err := h.users.GetByID(r.Context(), userID)
		if err == nil {
			email = u.Email
		}
	}
	
	code := cryptoRand6DigitCode()
	if err := h.svc.SetSecretDirectly(r.Context(), userID, code); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save verification secret")
		return
	}
	
	if err := h.send2FAEmail(email, code, "enroll"); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not send 2FA activation email: "+err.Error())
		return
	}
	
	response.JSON(w, http.StatusOK, map[string]string{"status": "pending_verification"})
}

// POST /api/v1/me/security/2fa/verify  { code }
func (h *SecurityHandler) Verify2FA(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	
	secret, _, err := h.svc.GetSecretDirectly(r.Context(), userID)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	if secret == "" {
		response.Err(w, http.StatusConflict, "TWOFA_NOT_SETUP", "Start enrollment first.")
		return
	}
	
	if secret != req.Code {
		response.Err(w, http.StatusBadRequest, "TWOFA_INVALID", "That code didn't match. Try again with the fresh code sent to your email.")
		return
	}
	
	if err := h.svc.EnableDirectly(r.Context(), userID); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	
	response.JSONMsg(w, http.StatusOK, "Two-factor authentication enabled")
}

// POST /api/v1/me/security/2fa/disable
func (h *SecurityHandler) Disable2FA(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	if err := h.svc.Disable(r.Context(), userID); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Two-factor authentication disabled")
}

// GET /api/v1/me/security/sessions
func (h *SecurityHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	if userID == "" {
		response.Err(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}
	sessions, err := h.tokens.ListActiveForUser(r.Context(), userID)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not retrieve active sessions")
		return
	}
	response.JSON(w, http.StatusOK, sessions)
}

// DELETE /api/v1/me/security/sessions/{id}
func (h *SecurityHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	sessionID := chi.URLParam(r, "id")
	if sessionID == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Session ID is required")
		return
	}
	if err := h.tokens.RevokeForUser(r.Context(), userID, sessionID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Session not found or already revoked")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not revoke session")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Session revoked successfully")
}

// POST /api/v1/me/security/sessions/revoke-others
func (h *SecurityHandler) RevokeOtherSessions(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	currentSessionID, _ := r.Context().Value(models.CtxSessionID).(string)
	if err := h.tokens.RevokeOthersForUser(r.Context(), userID, currentSessionID); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not revoke other sessions")
		return
	}
	response.JSONMsg(w, http.StatusOK, "All other sessions revoked successfully")
}

