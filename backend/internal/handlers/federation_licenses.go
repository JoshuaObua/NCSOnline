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

type FederationLicenseHandler struct {
	svc  *services.FederationLicenseService
	user *repository.UserRepo
}

func NewFederationLicenseHandler(svc *services.FederationLicenseService, user *repository.UserRepo) *FederationLicenseHandler {
	return &FederationLicenseHandler{svc: svc, user: user}
}

// GET /api/v1/admin/federation-licenses
func (h *FederationLicenseHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	p := &models.PaginationParams{Page: page, PerPage: perPage}

	filter := repository.ListLicensesFilter{
		FederationID: q.Get("federation_id"),
		Status:       q.Get("status"),
		LicenseType:  q.Get("license_type"),
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

// POST /api/v1/admin/federation-licenses
func (h *FederationLicenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	actorName := "System Administrator"
	if actorID != "" {
		if u, err := h.user.GetByID(r.Context(), actorID); err == nil && u != nil {
			actorName = strings.TrimSpace(u.FirstName + " " + u.LastName)
		}
	}

	var req struct {
		FederationID  string `json:"federation_id"`
		LicenseNumber string `json:"license_number"`
		LicenseType   string `json:"license_type"`
		Category      string `json:"category"`
		IssueDate     string `json:"issue_date"`
		ExpiryDate    string `json:"expiry_date"`
		Status        string `json:"status"`
		Conditions    string `json:"conditions"`
		DocumentURL   string `json:"document_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request payload")
		return
	}

	if req.FederationID == "" {
		response.Err(w, http.StatusBadRequest, "VALIDATION_ERROR", "Federation is required")
		return
	}
	if req.LicenseNumber == "" {
		req.LicenseNumber = "NCS/FED-LIC/" + time.Now().Format("2006") + "/" + strings.ToUpper(time.Now().Format("021504"))
	}

	issueDate := time.Now()
	if req.IssueDate != "" {
		if t, err := time.Parse("2006-01-02", req.IssueDate); err == nil {
			issueDate = t
		}
	}

	expiryDate := issueDate.AddDate(1, 0, 0)
	if req.ExpiryDate != "" {
		if t, err := time.Parse("2006-01-02", req.ExpiryDate); err == nil {
			expiryDate = t
		}
	}

	lic := &models.FederationLicense{
		FederationID:  req.FederationID,
		LicenseNumber: strings.TrimSpace(req.LicenseNumber),
		LicenseType:   req.LicenseType,
		Category:      req.Category,
		IssueDate:     issueDate,
		ExpiryDate:    expiryDate,
		Status:        req.Status,
		Conditions:    req.Conditions,
		DocumentURL:   req.DocumentURL,
	}

	if err := h.svc.CreateLicense(r.Context(), lic, actorID, actorName); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, lic)
}

// GET /api/v1/admin/federation-licenses/{id}
func (h *FederationLicenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	lic, err := h.svc.GetLicenseByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Federation license not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, lic)
}

// POST /api/v1/admin/federation-licenses/{id}/extend
func (h *FederationLicenseHandler) Extend(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	actorName := "System Administrator"
	if actorID != "" {
		if u, err := h.user.GetByID(r.Context(), actorID); err == nil && u != nil {
			actorName = strings.TrimSpace(u.FirstName + " " + u.LastName)
		}
	}

	var req struct {
		NewExpiryDate string `json:"new_expiry_date"`
		Reason        string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}
	if req.NewExpiryDate == "" {
		response.Err(w, http.StatusBadRequest, "VALIDATION_ERROR", "New expiry date is required")
		return
	}

	newExpiry, err := time.Parse("2006-01-02", req.NewExpiryDate)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid date format, expected YYYY-MM-DD")
		return
	}

	lic, err := h.svc.ExtendLicense(r.Context(), id, newExpiry, req.Reason, actorID, actorName)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, lic)
}

// POST /api/v1/admin/federation-licenses/{id}/revoke
func (h *FederationLicenseHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	actorName := "System Administrator"
	if actorID != "" {
		if u, err := h.user.GetByID(r.Context(), actorID); err == nil && u != nil {
			actorName = strings.TrimSpace(u.FirstName + " " + u.LastName)
		}
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if strings.TrimSpace(req.Reason) == "" {
		req.Reason = "Revocation ordered by NCS Board & Regulatory Secretariat."
	}

	lic, err := h.svc.RevokeLicense(r.Context(), id, req.Reason, actorID, actorName)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, lic)
}

// POST /api/v1/admin/federation-licenses/{id}/reinstate
func (h *FederationLicenseHandler) Reinstate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	actorName := "System Administrator"
	if actorID != "" {
		if u, err := h.user.GetByID(r.Context(), actorID); err == nil && u != nil {
			actorName = strings.TrimSpace(u.FirstName + " " + u.LastName)
		}
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if strings.TrimSpace(req.Reason) == "" {
		req.Reason = "Compliance cleared; statutory license reinstated."
	}

	lic, err := h.svc.ReinstateLicense(r.Context(), id, req.Reason, actorID, actorName)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, lic)
}

// GET /api/v1/admin/federation-licenses/kpis
func (h *FederationLicenseHandler) GetKPIs(w http.ResponseWriter, r *http.Request) {
	kpis, err := h.svc.GetKPIs(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, kpis)
}

// GET /api/v1/federations/{id}/license
func (h *FederationLicenseHandler) GetFederationActiveLicense(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	lic, err := h.svc.GetActiveLicenseByFederationID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "No license found for this federation")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, lic)
}
