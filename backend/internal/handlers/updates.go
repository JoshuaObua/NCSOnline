package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
)

type UpdatesHandler struct {
	svc      *services.UpdatesService
	deployer *services.Deployer
	backups  *repository.BackupRepo
}

// GET /api/v1/admin/system/updates — preflight + cached latest release.
func (h *UpdatesHandler) Status(w http.ResponseWriter, r *http.Request) {
	pf, err := h.svc.Preflight(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, pf)
}

// POST /api/v1/admin/system/updates/check — force an immediate poll.
func (h *UpdatesHandler) Check(w http.ResponseWriter, r *http.Request) {
	pf, err := h.svc.ForceRefresh(r.Context())
	if err != nil {
		response.Err(w, http.StatusBadGateway, "GITHUB_UNREACHABLE", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, pf)
}

func (h *UpdatesHandler) Settings(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, h.svc.Settings(r.Context()).Redacted())
}

func (h *UpdatesHandler) SaveSettings(w http.ResponseWriter, r *http.Request) {
	actor, _ := r.Context().Value(models.CtxUserID).(string)
	var input services.SmartUpdateSettings
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request")
		return
	}
	if strings.TrimSpace(input.GithubToken) == "********" {
		current := h.svc.Settings(r.Context())
		input.GithubToken = current.GithubToken
	}
	settings, err := h.svc.SaveSettings(r.Context(), input, actor)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save update settings")
		return
	}
	response.JSON(w, http.StatusOK, settings)
}

// GET /api/v1/admin/system/updates/deploy — last/current deploy progress.
func (h *UpdatesHandler) DeployStatus(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, h.deployer.Status())
}

// POST /api/v1/admin/system/updates/deploy — kick off `docker compose pull && up -d`.
func (h *UpdatesHandler) Deploy(w http.ResponseWriter, r *http.Request) {
	actor, _ := r.Context().Value(models.CtxUserID).(string)
	if h.deployer.Status().State == "running" {
		h.writeDeployErr(w, services.ErrDeployInProgress)
		return
	}
	settings := h.svc.Settings(r.Context())
	backupJob, err := h.backups.Queue(r.Context(), "BACKUP", nil, actor)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "BACKUP_QUEUE_FAILED", "Could not queue pre-deploy backup")
		return
	}
	if err := h.deployer.StartWithOptions(r.Context(), services.DeployOptions{
		ComposeProject: settings.ComposeProject,
		Services:       settings.DeployServices,
		Script:         settings.DeployScript,
		RepoSlug:       settings.RepoSlug,
		GithubToken:    settings.GithubToken,
		TargetBranch:   settings.TargetBranch,
		WorkTree:       settings.WorkTree,
		PreBackupJobID: backupJob.ID,
	}); err != nil {
		h.writeDeployErr(w, err)
		return
	}
	response.JSON(w, http.StatusAccepted, h.deployer.Status())
}

// POST /api/v1/admin/system/updates/rollback — re-tag :previous → :latest.
func (h *UpdatesHandler) Rollback(w http.ResponseWriter, r *http.Request) {
	if err := h.deployer.Rollback(r.Context()); err != nil {
		h.writeDeployErr(w, err)
		return
	}
	response.JSON(w, http.StatusAccepted, h.deployer.Status())
}

func (h *UpdatesHandler) writeDeployErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrDockerSocketMissing):
		response.Err(w, http.StatusFailedDependency, "DOCKER_SOCKET_MISSING", err.Error())
	case errors.Is(err, services.ErrDeployInProgress):
		response.Err(w, http.StatusConflict, "DEPLOY_IN_PROGRESS", err.Error())
	default:
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
	}
}
