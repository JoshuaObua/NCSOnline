package handlers

import (
	"net/http"

	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
)

type DashboardHandler struct {
	users         *repository.UserRepo
	apps          *repository.ApplicationRepo
	forms         *repository.FormRepo
	nsmis         *repository.NSMISRepo
	organisations *repository.OrganisationRepo
}

// GET /api/v1/admin/dashboard
func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	totalUsers, err := h.users.CountAll(ctx)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch user count")
		return
	}

	byStatus, err := h.apps.CountByStatus(ctx)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch application stats")
		return
	}

	dynamicByStatus, err := h.forms.CountSubmissionsByStatus(ctx)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch form submission stats")
		return
	}

	ordinaryUsers, err := h.users.CountOrdinary(ctx)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch ordinary user count")
		return
	}
	totalAthletes, err := h.nsmis.CountAthletes(ctx)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch athlete count")
		return
	}
	totalProfiles, err := h.organisations.CountAll(ctx)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch profile count")
		return
	}
	openForms, err := h.forms.CountOpenTemplates(ctx)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch open form count")
		return
	}

	legacyApps := int64(0)
	for _, n := range byStatus {
		legacyApps += n
	}
	dynamicApps := int64(0)
	combinedByStatus := map[string]int64{}
	for status, n := range byStatus {
		combinedByStatus[status] += n
	}
	for status, n := range dynamicByStatus {
		dynamicApps += n
		combinedByStatus[status] += n
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"total_users":                totalUsers,
		"ordinary_users":             ordinaryUsers,
		"total_athletes":             totalAthletes,
		"total_profiles":             totalProfiles,
		"open_forms":                 openForms,
		"total_applications":         legacyApps + dynamicApps,
		"legacy_applications":        legacyApps,
		"dynamic_applications":       dynamicApps,
		"applications_by_status":     combinedByStatus,
		"dynamic_submissions_status": dynamicByStatus,
		"pending_review":             combinedByStatus["SUBMITTED"] + combinedByStatus["UNDER_REVIEW"] + combinedByStatus["RESUBMITTED"],
		"needs_attention":            combinedByStatus["NEEDS_INFORMATION"] + combinedByStatus["PROOF_UPLOADED"] + combinedByStatus["PENDING_PAYMENT"],
	})
}
