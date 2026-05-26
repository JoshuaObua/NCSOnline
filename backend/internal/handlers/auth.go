package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
)

type AuthHandler struct{ svc *services.AuthService }

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
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "An unexpected error occurred")
		return
	}
	response.JSON(w, http.StatusOK, result)
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "refresh_token is required")
		return
	}

	result, err := h.svc.RefreshToken(r.Context(), req.RefreshToken, r.RemoteAddr, r.UserAgent())
	if errors.Is(err, services.ErrTokenInvalid) {
		response.Err(w, http.StatusUnauthorized, "TOKEN_INVALID", "Refresh token is invalid or expired")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "An unexpected error occurred")
		return
	}
	response.JSON(w, http.StatusOK, result)
}

// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	_ = h.svc.Logout(r.Context(), req.RefreshToken)
	response.JSONMsg(w, http.StatusOK, "Logged out successfully")
}

// GET /api/v1/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	email, _ := r.Context().Value(models.CtxUserEmail).(string)
	roles, _ := r.Context().Value(models.CtxUserRoles).([]string)
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"email":   email,
		"roles":   roles,
	})
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

// POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: implement email-based reset token generation and dispatch
	response.JSONMsg(w, http.StatusOK, "If that email exists, a reset link has been sent")
}

// POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: validate reset token and update password
	response.JSONMsg(w, http.StatusOK, "Password has been reset")
}
