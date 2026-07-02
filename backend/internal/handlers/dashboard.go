package handlers

import (
	"net/http"

	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
)

type DashboardHandler struct {
	users *repository.UserRepo
	apps  *repository.ApplicationRepo
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

	totalApps := int64(0)
	for _, n := range byStatus {
		totalApps += n
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"total_users":            totalUsers,
		"total_applications":     totalApps,
		"applications_by_status": byStatus,
		"pending_review":         byStatus["SUBMITTED"] + byStatus["UNDER_REVIEW"] + byStatus["RESUBMITTED"],
		"needs_attention":        byStatus["NEEDS_INFORMATION"] + byStatus["PROOF_UPLOADED"],
	})
}
