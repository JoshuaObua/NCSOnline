package handlers

import (
	"errors"
	"net/http"

	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
)

type UpdatesHandler struct {
	svc      *services.UpdatesService
	deployer *services.Deployer
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

// GET /api/v1/admin/system/updates/deploy — last/current deploy progress.
func (h *UpdatesHandler) DeployStatus(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, h.deployer.Status())
}

// POST /api/v1/admin/system/updates/deploy — kick off `docker compose pull && up -d`.
func (h *UpdatesHandler) Deploy(w http.ResponseWriter, r *http.Request) {
	if err := h.deployer.Start(r.Context()); err != nil {
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
