package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
)

type BackupsHandler struct{ repo *repository.BackupRepo }

func (h *BackupsHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not list backups")
		return
	}
	jobs, _ := h.repo.ListJobs(r.Context())
	response.JSON(w, http.StatusOK, map[string]any{"backups": items, "jobs": jobs})
}
func (h *BackupsHandler) Queue(w http.ResponseWriter, r *http.Request) {
	actor, _ := r.Context().Value(models.CtxUserID).(string)
	var input struct {
		Action       string  `json:"action"`
		BackupID     *string `json:"backup_id"`
		Confirmation string  `json:"confirmation"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		response.Err(w, 400, "BAD_REQUEST", "Invalid request")
		return
	}
	kind := strings.ToUpper(input.Action)
	if kind != "BACKUP" && kind != "VERIFY" && kind != "RESTORE" {
		response.Err(w, 400, "BAD_REQUEST", "action must be BACKUP, VERIFY or RESTORE")
		return
	}
	if kind != "BACKUP" && (input.BackupID == nil || *input.BackupID == "") {
		response.Err(w, 400, "BAD_REQUEST", "backup_id is required")
		return
	}
	if kind == "RESTORE" {
		if input.Confirmation != "RESTORE VERIFIED BACKUP" {
			response.Err(w, 400, "CONFIRMATION_REQUIRED", "Type RESTORE VERIFIED BACKUP to confirm")
			return
		}
		record, err := h.repo.Get(r.Context(), *input.BackupID)
		if err != nil || record.Status != "VERIFIED" {
			response.Err(w, 400, "BACKUP_NOT_VERIFIED", "Only a verified backup can be restored")
			return
		}
		if !h.stateEnabled() {
			response.Err(w, 409, "MAINTENANCE_REQUIRED", "Enable maintenance mode before restoring")
			return
		}
	}
	job, err := h.repo.Queue(r.Context(), kind, input.BackupID, actor)
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not queue backup job")
		return
	}
	response.JSON(w, http.StatusAccepted, job)
}
func (h *BackupsHandler) stateEnabled() bool {
	enabled, _ := h.repo.MaintenanceEnabled(context.Background())
	return enabled
}
