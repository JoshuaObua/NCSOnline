package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
)

func domainPermission(resource string, write bool) []string {
	if strings.HasPrefix(resource, "safeguarding-") {
		return []string{"safeguarding:cases:manage"}
	}
	if strings.HasPrefix(resource, "federation-") {
		if write {
			return []string{"federations:write:own", "federations:write:any"}
		}
		return []string{"federations:read:own", "federations:read:any"}
	}
	if write {
		return []string{"reports:write:own", "reports:review:any"}
	}
	return []string{"reports:read:own", "reports:read:any"}
}

func (h *NSMISHandler) domainAccess(w http.ResponseWriter, r *http.Request, write bool) (string, []string, bool, bool) {
	resource := chi.URLParam(r, "resource")
	if !h.allowed(r, domainPermission(resource, write)...) {
		response.Err(w, 403, "FORBIDDEN", "This account cannot access the requested NSMIS records")
		return "", nil, false, false
	}
	uid, _ := r.Context().Value(models.CtxUserID).(string)
	wide := h.wide(r, "reports:review:any") || h.wide(r, "reports:read:any") || h.wide(r, "federations:read:any") || h.wide(r, "federations:write:any")
	scope, err := h.repo.FederationScope(r.Context(), uid, wide)
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not resolve federation scope")
		return "", nil, false, false
	}
	return uid, scope, wide, true
}

func (h *NSMISHandler) ListDomain(w http.ResponseWriter, r *http.Request) {
	uid, _, wide, ok := h.domainAccess(w, r, false)
	if !ok {
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	per, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if per < 1 {
		per = 25
	}
	if per > 100 {
		per = 100
	}
	items, total, err := h.repo.ListDomain(r.Context(), chi.URLParam(r, "resource"), uid, strings.TrimSpace(r.URL.Query().Get("search")), wide, per, (page-1)*per)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, 404, "NOT_FOUND", "Unknown NSMIS resource")
		return
	}
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load records")
		return
	}
	response.JSON(w, 200, map[string]interface{}{"items": items, "pagination": map[string]interface{}{"page": page, "per_page": per, "total": total}})
}

func (h *NSMISHandler) CreateDomain(w http.ResponseWriter, r *http.Request) {
	uid, scope, wide, ok := h.domainAccess(w, r, true)
	if !ok {
		return
	}
	var payload map[string]interface{}
	if !decodeStrict(w, r, &payload) {
		return
	}
	item, err := h.repo.SaveDomain(r.Context(), chi.URLParam(r, "resource"), "", uid, payload, scope, wide)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, 403, "OUT_OF_SCOPE", "The federation is outside your assigned scope")
		return
	}
	if errors.Is(err, repository.ErrDuplicate) {
		response.Err(w, 409, "DUPLICATE", "A matching record already exists")
		return
	}
	if err != nil {
		response.Err(w, 422, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(w, 201, item)
}

func (h *NSMISHandler) UpdateDomain(w http.ResponseWriter, r *http.Request) {
	uid, scope, wide, ok := h.domainAccess(w, r, true)
	if !ok {
		return
	}
	var payload map[string]interface{}
	if !decodeStrict(w, r, &payload) {
		return
	}
	item, err := h.repo.SaveDomain(r.Context(), chi.URLParam(r, "resource"), chi.URLParam(r, "id"), uid, payload, scope, wide)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, 404, "NOT_FOUND", "Record not found in your assigned scope")
		return
	}
	if errors.Is(err, repository.ErrDuplicate) {
		response.Err(w, 409, "DUPLICATE", "A matching record already exists")
		return
	}
	if err != nil {
		response.Err(w, 422, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(w, 200, item)
}

func (h *NSMISHandler) DeleteDomain(w http.ResponseWriter, r *http.Request) {
	_, scope, wide, ok := h.domainAccess(w, r, true)
	if !ok {
		return
	}
	err := h.repo.DeleteDomain(r.Context(), chi.URLParam(r, "resource"), chi.URLParam(r, "id"), scope, wide)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, 404, "NOT_FOUND", "Record not found in your assigned scope")
		return
	}
	if err != nil {
		response.Err(w, 409, "RECORD_IN_USE", "This record is referenced by history and cannot be deleted")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
