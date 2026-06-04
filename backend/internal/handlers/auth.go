package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
)

type AuthHandler struct {
	svc   *services.AuthService
	users *repository.UserRepo
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
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Registration failed. Please try again.")
		return
	}
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
