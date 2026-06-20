package handlers

import (
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
)

var allowedJobs = map[string]bool{"COMPLIANCE_REFRESH": true, "DEADLINE_REMINDERS": true, "CREDENTIAL_EXPIRY": true, "DATA_QUALITY_SCAN": true}

func (h *NSMISHandler) QueueJob(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "reports:review:any") {
		response.Err(w, 403, "FORBIDDEN", "Job administration requires NCS review permission")
		return
	}
	var req struct {
		JobType        string                 `json:"job_type"`
		IdempotencyKey string                 `json:"idempotency_key"`
		Payload        map[string]interface{} `json:"payload"`
	}
	if !decodeStrict(w, r, &req) {
		return
	}
	req.JobType = strings.ToUpper(strings.TrimSpace(req.JobType))
	if !allowedJobs[req.JobType] {
		response.ValidationErr(w, map[string]string{"job_type": "Unsupported job type"})
		return
	}
	uid, _ := r.Context().Value(models.CtxUserID).(string)
	id, err := h.repo.EnqueueJob(r.Context(), req.JobType, strings.TrimSpace(req.IdempotencyKey), req.Payload, uid)
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Job could not be queued")
		return
	}
	response.JSON(w, 202, map[string]string{"id": id, "status": "PENDING"})
}
func (h *NSMISHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "reports:review:any") {
		response.Err(w, 403, "FORBIDDEN", "Job administration requires NCS review permission")
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	items, err := h.repo.ListJobs(r.Context(), limit)
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load jobs")
		return
	}
	response.JSON(w, 200, items)
}
func (h *NSMISHandler) RetryJob(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "reports:review:any") {
		response.Err(w, 403, "FORBIDDEN", "Job administration requires NCS review permission")
		return
	}
	if err := h.repo.RetryJob(r.Context(), chi.URLParam(r, "jobID")); err == repository.ErrNotFound {
		response.Err(w, 404, "NOT_FOUND", "Failed job not found")
		return
	} else if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Job could not be retried")
		return
	}
	response.JSONMsg(w, 200, "Job queued for retry")
}
