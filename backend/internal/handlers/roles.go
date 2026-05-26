package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"
)

type RolesHandler struct {
	roles *repository.RoleRepo
	audit *repository.AuditRepo
}

// GET /api/v1/admin/roles
func (h *RolesHandler) List(w http.ResponseWriter, r *http.Request) {
	roles, err := h.roles.List(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch roles")
		return
	}
	response.JSON(w, http.StatusOK, roles)
}

// POST /api/v1/admin/roles
func (h *RolesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "name is required")
		return
	}
	role := &models.Role{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
	}
	if err := h.roles.Create(r.Context(), role); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "ROLE_EXISTS", "A role with that name already exists")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create role")
		return
	}
	response.JSON(w, http.StatusCreated, role)
}

// GET /api/v1/admin/roles/{id}
func (h *RolesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	role, err := h.roles.GetByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Role not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch role")
		return
	}
	response.JSON(w, http.StatusOK, role)
}

// PUT /api/v1/admin/roles/{id}
func (h *RolesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "name is required")
		return
	}
	role := &models.Role{ID: id, Name: req.Name, Description: req.Description}
	if err := h.roles.Update(r.Context(), role); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update role")
		return
	}
	response.JSON(w, http.StatusOK, role)
}

// DELETE /api/v1/admin/roles/{id}
func (h *RolesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.roles.Delete(r.Context(), id); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete role")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Role deleted")
}

// GET /api/v1/admin/permissions
func (h *RolesHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	perms, err := h.roles.ListPermissions(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch permissions")
		return
	}
	response.JSON(w, http.StatusOK, perms)
}

// POST /api/v1/admin/roles/{id}/permissions
func (h *RolesHandler) AssignPermission(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	var req struct {
		PermissionID string `json:"permission_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PermissionID == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "permission_id is required")
		return
	}
	if err := h.roles.AssignPermission(r.Context(), roleID, req.PermissionID); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not assign permission")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Permission assigned")
}

// DELETE /api/v1/admin/roles/{id}/permissions/{permID}
func (h *RolesHandler) RemovePermission(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	permID := chi.URLParam(r, "permID")
	if err := h.roles.RemovePermission(r.Context(), roleID, permID); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not remove permission")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Permission removed")
}
