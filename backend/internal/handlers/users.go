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

// GET /api/v1/admin/users
func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	p := &models.PaginationParams{Page: page, PerPage: perPage, Search: q.Get("search")}

	users, total, err := h.svc.List(r.Context(), p)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch users")
		return
	}
	response.JSONPaged(w, http.StatusOK, users, &response.Meta{Page: p.Page, PerPage: p.PerPage, Total: total})
}

// POST /api/v1/admin/users
func (h *UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
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

	user, err := h.svc.Create(r.Context(), services.CreateUserInput{
		Email: req.Email, Password: req.Password,
		FirstName: req.FirstName, LastName: req.LastName, Phone: req.Phone,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "EMAIL_TAKEN", "Email is already registered")
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
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	user, err := h.svc.Update(r.Context(), id, services.UpdateUserInput{
		FirstName: req.FirstName, LastName: req.LastName, Phone: req.Phone,
	})
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update user")
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// DELETE /api/v1/admin/users/{id}
func (h *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
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
	if err := h.svc.SetActive(r.Context(), id, true); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not activate user")
		return
	}
	response.JSONMsg(w, http.StatusOK, "User activated")
}

// POST /api/v1/admin/users/{id}/deactivate
func (h *UsersHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.SetActive(r.Context(), id, false); err != nil {
		if protectedGuard(w, err) {
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not deactivate user")
		return
	}
	response.JSONMsg(w, http.StatusOK, "User deactivated")
}

// POST /api/v1/admin/users/{id}/roles
func (h *UsersHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	assignerID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		RoleID string `json:"role_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RoleID == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "role_id is required")
		return
	}
	if err := h.svc.AssignRole(r.Context(), userID, req.RoleID, assignerID); err != nil {
		response.Err(w, http.StatusBadRequest, "ASSIGN_ROLE_FAILED", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Role assigned")
}

// DELETE /api/v1/admin/users/{id}/roles/{roleID}
func (h *UsersHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	roleID := chi.URLParam(r, "roleID")
	if err := h.svc.RemoveRole(r.Context(), userID, roleID); err != nil {
		if protectedGuard(w, err) {
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not remove role")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Role removed")
}
