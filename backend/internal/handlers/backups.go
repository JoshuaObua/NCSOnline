package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
)

type BackupsHandler struct {
	repo *repository.BackupRepo
	cfg  *config.Config
}

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
func (h *BackupsHandler) Download(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	record, err := h.repo.Get(r.Context(), id)
	if err != nil {
		response.Err(w, 404, "NOT_FOUND", "Backup not found")
		return
	}
	path := h.backupPath(record.FileName)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filepath.Base(record.FileName)}))
	http.ServeFile(w, r, path)
}
func (h *BackupsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	record, err := h.repo.Get(r.Context(), id)
	if err != nil {
		response.Err(w, 404, "NOT_FOUND", "Backup not found")
		return
	}
	if err := h.repo.Delete(r.Context(), id); err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not delete backup record")
		return
	}
	_ = os.Remove(h.backupPath(record.FileName))
	_ = os.Remove(h.backupPath(record.FileName) + ".sha256")
	response.JSONMsg(w, http.StatusOK, "Backup deleted")
}

func (h *BackupsHandler) ExportJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.repo.ListJobs(r.Context())
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not export job logs")
		return
	}
	name := fmt.Sprintf("backup-jobs-%s.csv", time.Now().UTC().Format("20060102-150405"))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": name}))
	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"id", "job_type", "backup_id", "status", "progress", "message", "requested_by", "created_at", "started_at", "finished_at"})
	for _, job := range jobs {
		_ = cw.Write([]string{
			job.ID, job.JobType, strPtr(job.BackupID), job.Status, fmt.Sprint(job.Progress), job.Message,
			strPtr(job.RequestedBy), formatTime(job.CreatedAt), formatTimePtr(job.StartedAt), formatTimePtr(job.FinishedAt),
		})
	}
	cw.Flush()
}

func (h *BackupsHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "jobID")
	if err := h.repo.DeleteJob(r.Context(), id); err != nil {
		response.Err(w, 404, "NOT_FOUND", "Job log not found")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Job log deleted")
}

func (h *BackupsHandler) ClearJobs(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Confirmation string `json:"confirmation"`
	}
	_ = json.NewDecoder(r.Body).Decode(&input)
	if input.Confirmation != "CLEAR BACKUP JOB LOGS" {
		response.Err(w, 400, "CONFIRMATION_REQUIRED", "Type CLEAR BACKUP JOB LOGS to confirm")
		return
	}
	if err := h.repo.ClearJobs(r.Context()); err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not clear job logs")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Finished backup job logs cleared")
}
func (h *BackupsHandler) DownloadSchema(w http.ResponseWriter, r *http.Request) {
	name := fmt.Sprintf("ncs-schema-%s.sql", time.Now().UTC().Format("20060102-150405"))
	cmd := exec.CommandContext(r.Context(), "pg_dump", "--dbname", h.cfg.DatabaseURL, "--schema-only", "--no-owner", "--no-privileges")
	out, err := cmd.Output()
	if err != nil {
		response.Err(w, 500, "SCHEMA_EXPORT_FAILED", "Could not export database schema")
		return
	}
	w.Header().Set("Content-Type", "application/sql")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": name}))
	_, _ = w.Write(out)
}
func (h *BackupsHandler) ImportSchema(w http.ResponseWriter, r *http.Request) {
	if !h.stateEnabled() {
		response.Err(w, 409, "MAINTENANCE_REQUIRED", "Enable maintenance mode before importing a schema")
		return
	}
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		response.Err(w, 400, "BAD_REQUEST", "Could not parse schema file")
		return
	}
	if r.FormValue("confirmation") != "IMPORT SCHEMA" {
		response.Err(w, 400, "CONFIRMATION_REQUIRED", "Type IMPORT SCHEMA to confirm")
		return
	}
	file, header, err := r.FormFile("schema")
	if err != nil {
		response.Err(w, 400, "BAD_REQUEST", "schema file is required")
		return
	}
	defer file.Close()
	if strings.ToLower(filepath.Ext(header.Filename)) != ".sql" {
		response.Err(w, 400, "BAD_REQUEST", "Only .sql schema files are accepted")
		return
	}
	tmp, err := os.CreateTemp("", "ncs-schema-*.sql")
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not stage schema file")
		return
	}
	defer os.Remove(tmp.Name())
	if _, err := io.Copy(tmp, io.LimitReader(file, 8<<20)); err != nil {
		tmp.Close()
		response.Err(w, 500, "SERVER_ERROR", "Could not read schema file")
		return
	}
	_ = tmp.Close()
	out, err := exec.CommandContext(r.Context(), "psql", h.cfg.DatabaseURL, "-v", "ON_ERROR_STOP=1", "--single-transaction", "-f", tmp.Name()).CombinedOutput()
	if err != nil {
		response.Err(w, 500, "SCHEMA_IMPORT_FAILED", safeCommandOutput(out))
		return
	}
	response.JSONMsg(w, http.StatusOK, "Schema imported")
}
func (h *BackupsHandler) DeleteSchema(w http.ResponseWriter, r *http.Request) {
	if !h.stateEnabled() {
		response.Err(w, 409, "MAINTENANCE_REQUIRED", "Enable maintenance mode before deleting the database schema")
		return
	}
	var input struct {
		Confirmation string `json:"confirmation"`
	}
	_ = json.NewDecoder(r.Body).Decode(&input)
	if input.Confirmation != "DELETE DATABASE SCHEMA" {
		response.Err(w, 400, "CONFIRMATION_REQUIRED", "Type DELETE DATABASE SCHEMA to confirm")
		return
	}
	sql := `DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO public;`
	out, err := exec.CommandContext(r.Context(), "psql", h.cfg.DatabaseURL, "-v", "ON_ERROR_STOP=1", "-c", sql).CombinedOutput()
	if err != nil {
		response.Err(w, 500, "SCHEMA_DELETE_FAILED", safeCommandOutput(out))
		return
	}
	response.JSONMsg(w, http.StatusOK, "Database schema deleted and recreated")
}
func (h *BackupsHandler) stateEnabled() bool {
	enabled, _ := h.repo.MaintenanceEnabled(context.Background())
	return enabled
}
func (h *BackupsHandler) backupPath(name string) string {
	dir := os.Getenv("BACKUP_DIR")
	if dir == "" {
		dir = "/var/backups/ncs"
	}
	return filepath.Join(dir, filepath.Base(name))
}
func safeCommandOutput(out []byte) string {
	msg := strings.TrimSpace(string(out))
	if msg == "" {
		return "Command failed"
	}
	if len(msg) > 1000 {
		msg = msg[:1000]
	}
	return msg
}

func strPtr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func formatTime(v time.Time) string {
	return v.UTC().Format(time.RFC3339)
}

func formatTimePtr(v *time.Time) string {
	if v == nil {
		return ""
	}
	return formatTime(*v)
}
