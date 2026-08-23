package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
	"github.com/go-chi/chi/v5"
)

// protectedGuard maps ErrProtectedAccount to HTTP 403 and returns true if handled.
func protectedGuard(w http.ResponseWriter, err error) bool {
	if errors.Is(err, services.ErrProtectedAccount) {
		response.Err(w, http.StatusForbidden, "PROTECTED_ACCOUNT", "The super admin account cannot be modified or deleted")
		return true
	}
	return false
}

type UsersHandler struct {
	svc   *services.UserService
	audit *repository.AuditRepo
}

func (h *UsersHandler) checkAccess(w http.ResponseWriter, r *http.Request, targetUserID, requiredGlobalPerm, requiredScopedPerm string) (bool, []string) {
	actorID, _ := r.Context().Value(models.CtxUserID).(string)

	global, err := h.svc.HasPermission(r.Context(), actorID, requiredGlobalPerm)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not check permissions")
		return false, nil
	}
	if global {
		return true, nil
	}

	scoped, err := h.svc.HasPermission(r.Context(), actorID, requiredScopedPerm)
	if err != nil || !scoped {
		response.Err(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to perform this action")
		return false, nil
	}

	federations, err := h.svc.GetFederationIDs(r.Context(), actorID)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load federation scope")
		return false, nil
	}
	if len(federations) == 0 {
		response.Err(w, http.StatusForbidden, "FORBIDDEN", "You are not associated with any active federation")
		return false, nil
	}

	if targetUserID != "" {
		inFed, err := h.svc.IsUserInFederations(r.Context(), targetUserID, federations)
		if err != nil || !inFed {
			response.Err(w, http.StatusForbidden, "FORBIDDEN", "The target user does not belong to your federation")
			return false, nil
		}
	}

	return true, federations
}

// GET /api/v1/admin/users
func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	p := &models.PaginationParams{Page: page, PerPage: perPage, Search: q.Get("search")}

	ok, federations := h.checkAccess(w, r, "", "users:read", "users:read:own")
	if !ok {
		return
	}

	users, total, err := h.svc.List(r.Context(), p, federations)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch users")
		return
	}
	response.JSONPaged(w, http.StatusOK, users, &response.Meta{Page: p.Page, PerPage: p.PerPage, Total: total})
}

// POST /api/v1/admin/users
func (h *UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Phone        string `json:"phone"`
		NIN          string `json:"nin"`
		FederationID string `json:"federation_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}

	ok, federations := h.checkAccess(w, r, "", "users:write", "users:write:own")
	if !ok {
		return
	}

	targetFed := req.FederationID
	if len(federations) > 0 { // Actor is federation scoped
		if targetFed == "" {
			if len(federations) == 1 {
				targetFed = federations[0]
			} else {
				response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "federation_id is required since you belong to multiple federations")
				return
			}
		} else {
			// Verify they belong to targetFed
			valid := false
			for _, f := range federations {
				if f == targetFed {
					valid = true
					break
				}
			}
			if !valid {
				response.Err(w, http.StatusForbidden, "FORBIDDEN", "You do not belong to the specified federation")
				return
			}
		}
	}

	errs := map[string]string{}
	if req.Email == "" {
		errs["email"] = "required"
	}
	if len(req.Password) < 8 {
		errs["password"] = "minimum 8 characters"
	}
	if req.FirstName == "" {
		errs["first_name"] = "required"
	}
	if req.LastName == "" {
		errs["last_name"] = "required"
	}
	if len(errs) > 0 {
		response.ValidationErr(w, errs)
		return
	}

	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	user, err := h.svc.Create(r.Context(), services.CreateUserInput{
		Email: req.Email, Password: req.Password,
		FirstName: req.FirstName, LastName: req.LastName, Phone: req.Phone, NIN: req.NIN,
		FederationID: targetFed, ActorID: actorID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "EMAIL_TAKEN", "Email is already registered")
			return
		}
		// NIN duplicate surfaces as a plain error string from the service
		if strings.Contains(err.Error(), "NIN") && strings.Contains(err.Error(), "already registered") {
			response.Err(w, http.StatusConflict, "NIN_TAKEN", err.Error())
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, user)
}

// GET /api/v1/admin/users/{id}
func (h *UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ok, _ := h.checkAccess(w, r, id, "users:read", "users:read:own")
	if !ok {
		return
	}
	user, err := h.svc.GetByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch user")
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// PUT /api/v1/admin/users/{id}
func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ok, _ := h.checkAccess(w, r, id, "users:write", "users:write:own")
	if !ok {
		return
	}
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		NIN       string `json:"nin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	user, err := h.svc.Update(r.Context(), id, services.UpdateUserInput{
		FirstName: req.FirstName, LastName: req.LastName, Email: req.Email, Phone: req.Phone, NIN: req.NIN,
	})
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}
	if err != nil {
		if strings.Contains(err.Error(), "NIN") && strings.Contains(err.Error(), "already registered") {
			response.Err(w, http.StatusConflict, "NIN_TAKEN", err.Error())
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update user")
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// DELETE /api/v1/admin/users/{id}
func (h *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ok, _ := h.checkAccess(w, r, id, "users:delete", "users:write:own")
	if !ok {
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		if protectedGuard(w, err) {
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete user")
		return
	}
	response.JSONMsg(w, http.StatusOK, "User deleted")
}

// POST /api/v1/admin/users/{id}/activate
func (h *UsersHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ok, _ := h.checkAccess(w, r, id, "users:activate", "users:write:own")
	if !ok {
		return
	}
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	if err := h.svc.ApplyAccountAction(r.Context(), services.AccountActionInput{UserID: id, ActorID: actorID, Action: "REACTIVATE"}); err != nil {
		if protectedGuard(w, err) {
			return
		}
		response.Err(w, http.StatusBadRequest, "ACTIVATE_FAILED", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "User activated")
}

// POST /api/v1/admin/users/{id}/deactivate
func (h *UsersHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ok, _ := h.checkAccess(w, r, id, "users:activate", "users:write:own")
	if !ok {
		return
	}
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	if err := h.svc.ApplyAccountAction(r.Context(), services.AccountActionInput{
		UserID: id, ActorID: actorID, Action: "SUSPEND", Reason: "Deactivated by administrator",
	}); err != nil {
		if protectedGuard(w, err) {
			return
		}
		response.Err(w, http.StatusBadRequest, "DEACTIVATE_FAILED", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "User deactivated")
}

// POST /api/v1/admin/users/{id}/roles
func (h *UsersHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	assignerID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		RoleID   string `json:"role_id"`
		RoleName string `json:"role_name"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}
	targetRole := req.RoleID
	if targetRole == "" {
		targetRole = req.RoleName
	}
	if targetRole == "" {
		targetRole = req.Role
	}
	if targetRole == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "role_id is required")
		return
	}
	if err := h.svc.AssignRole(r.Context(), userID, targetRole, assignerID); err != nil {
		response.Err(w, http.StatusBadRequest, "ASSIGN_ROLE_FAILED", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Role assigned")
}

// DELETE /api/v1/admin/users/{id}/roles/{roleID}
func (h *UsersHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	roleID := chi.URLParam(r, "roleID")
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	if err := h.svc.RemoveRole(r.Context(), userID, roleID, actorID); err != nil {
		if protectedGuard(w, err) {
			return
		}
		response.Err(w, http.StatusBadRequest, "REMOVE_ROLE_FAILED", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Role removed")
}

// POST /api/v1/admin/users/{id}/reset-password
func (h *UsersHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	ok, _ := h.checkAccess(w, r, userID, "users:reset_password", "users:write:own")
	if !ok {
		return
	}
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	if actorID == userID {
		response.Err(w, http.StatusBadRequest, "USE_CHANGE_PASSWORD", "Use the profile password-change form for your own account")
		return
	}
	var req struct {
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if len(req.NewPassword) < 12 {
		response.ValidationErr(w, map[string]string{"new_password": "minimum 12 characters"})
		return
	}
	if err := h.svc.ResetPassword(r.Context(), userID, req.NewPassword); err != nil {
		if protectedGuard(w, err) {
			return
		}
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "PASSWORD_RESET_FAILED", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Password reset successfully. All existing sessions were revoked")
}

// POST /api/v1/admin/users/{id}/account-action
func (h *UsersHandler) AccountAction(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		Action         string     `json:"action"`
		Reason         string     `json:"reason"`
		SuspendedUntil *time.Time `json:"suspended_until"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	err := h.svc.ApplyAccountAction(r.Context(), services.AccountActionInput{
		UserID: userID, ActorID: actorID, Action: req.Action,
		Reason: req.Reason, SuspendedUntil: req.SuspendedUntil,
	})
	if err != nil {
		if protectedGuard(w, err) {
			return
		}
		switch {
		case errors.Is(err, services.ErrSelfAccountAction):
			response.Err(w, http.StatusBadRequest, "SELF_ACCOUNT_ACTION", err.Error())
		case errors.Is(err, services.ErrFraudFlagged):
			response.Err(w, http.StatusConflict, "FRAUD_FLAGGED", err.Error())
		case errors.Is(err, services.ErrInvalidAccountAction):
			response.Err(w, http.StatusBadRequest, "INVALID_ACCOUNT_ACTION", err.Error())
		case errors.Is(err, repository.ErrNotFound):
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		default:
			response.Err(w, http.StatusBadRequest, "ACCOUNT_ACTION_FAILED", err.Error())
		}
		return
	}
	response.JSONMsg(w, http.StatusOK, "Account action completed successfully")
}
