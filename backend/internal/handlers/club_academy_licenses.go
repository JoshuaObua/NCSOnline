package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
)

type ClubAcademyLicenseHandler struct {
	svc  *services.ClubAcademyLicenseService
	user *repository.UserRepo
}

func NewClubAcademyLicenseHandler(svc *services.ClubAcademyLicenseService, user *repository.UserRepo) *ClubAcademyLicenseHandler {
	return &ClubAcademyLicenseHandler{svc: svc, user: user}
}

func (h *ClubAcademyLicenseHandler) resolveActor(r *http.Request) (string, string) {
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	actorName := "System Administrator"
	if actorID != "" && h.user != nil {
		if u, err := h.user.GetByID(r.Context(), actorID); err == nil && u != nil {
			name := (u.FirstName + " " + u.LastName)
			if name != " " && name != "" {
				actorName = name
			}
		}
	}
	return actorID, actorName
}

// GET /api/v1/admin/club-academy-licenses
func (h *ClubAcademyLicenseHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	p := &models.PaginationParams{Page: page, PerPage: perPage}

	filter := repository.ListClubLicensesFilter{
		ClubID:       q.Get("club_id"),
		FederationID: q.Get("federation_id"),
		Status:       q.Get("status"),
		LicenseType:  q.Get("license_type"),
		Category:     q.Get("category"),
		Search:       q.Get("search"),
		StartDate:    q.Get("start_date"),
		EndDate:      q.Get("end_date"),
		SortBy:       q.Get("sort_by"),
		SortOrder:    q.Get("sort_order"),
	}

	list, total, err := h.svc.ListLicenses(r.Context(), filter, p)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSONPaged(w, http.StatusOK, list, &response.Meta{Page: p.Page, PerPage: p.PerPage, Total: int64(total)})
}

// POST /api/v1/admin/club-academy-licenses
func (h *ClubAcademyLicenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID, actorName := h.resolveActor(r)

	var req models.CreateClubLicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Failed to parse request body")
		return
	}

	lic, err := h.svc.CreateLicense(r.Context(), req, actorID, actorName)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "CREATE_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, lic)
}

// GET /api/v1/admin/club-academy-licenses/kpis
func (h *ClubAcademyLicenseHandler) GetKPIs(w http.ResponseWriter, r *http.Request) {
	kpis, err := h.svc.GetKPIs(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, kpis)
}

// GET /api/v1/admin/club-academy-licenses/{id}
func (h *ClubAcademyLicenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	lic, err := h.svc.GetLicenseByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Club / Academy license not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, lic)
}

// POST /api/v1/admin/club-academy-licenses/{id}/extend
func (h *ClubAcademyLicenseHandler) Extend(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	actorID, actorName := h.resolveActor(r)

	var req models.ExtendClubLicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Failed to parse request body")
		return
	}

	lic, err := h.svc.ExtendLicense(r.Context(), id, req, actorID, actorName)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "EXTEND_FAILED", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, lic)
}

// POST /api/v1/admin/club-academy-licenses/{id}/revoke
func (h *ClubAcademyLicenseHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	actorID, actorName := h.resolveActor(r)

	var req models.RevokeClubLicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Failed to parse request body")
		return
	}

	lic, err := h.svc.RevokeLicense(r.Context(), id, req, actorID, actorName)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "REVOKE_FAILED", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, lic)
}

// POST /api/v1/admin/club-academy-licenses/{id}/reinstate
func (h *ClubAcademyLicenseHandler) Reinstate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	actorID, actorName := h.resolveActor(r)

	var req models.ReinstateClubLicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Failed to parse request body")
		return
	}

	lic, err := h.svc.ReinstateLicense(r.Context(), id, req, actorID, actorName)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "REINSTATE_FAILED", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, lic)
}

// GET /api/v1/clubs/{id}/license
func (h *ClubAcademyLicenseHandler) GetClubActiveLicense(w http.ResponseWriter, r *http.Request) {
	clubID := chi.URLParam(r, "id")
	lic, err := h.svc.GetActiveLicenseByClubID(r.Context(), clubID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Active license not found for this club/academy")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, lic)
}
