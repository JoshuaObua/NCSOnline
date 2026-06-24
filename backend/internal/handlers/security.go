package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
	"github.com/go-chi/chi/v5"
)

type SecurityHandler struct {
	svc   *services.SecurityService
	users *repository.UserRepo
	audit *repository.AuditRepo
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

// POST /api/v1/me/security/2fa/enroll
// Generates and persists a new TOTP secret in a pending state, returns the
// otpauth URI for QR rendering. The frontend renders that URI as a QR code
// via a public service (e.g. api.qrserver.com) or displays the secret as
// text for manual entry.
func (h *SecurityHandler) Enroll2FA(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	email, _ := r.Context().Value(models.CtxUserEmail).(string)
	if email == "" {
		// Fall back to db lookup if the JWT didn't carry the email.
		u, err := h.users.GetByID(r.Context(), userID)
		if err == nil {
			email = u.Email
		}
	}
	enr, err := h.svc.BeginEnrollment(r.Context(), userID, email)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, enr)
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
	if err := h.svc.VerifyAndEnable(r.Context(), userID, req.Code); err != nil {
		switch {
		case errors.Is(err, services.ErrTwoFAInvalid):
			response.Err(w, http.StatusBadRequest, "TWOFA_INVALID", "That code didn't match. Try again with a fresh code from your authenticator app.")
		case errors.Is(err, services.ErrTwoFANotSetup):
			response.Err(w, http.StatusConflict, "TWOFA_NOT_SETUP", "Start enrollment first.")
		default:
			response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		}
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
