package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
)

type OrganisationsHandler struct{ repo *repository.OrganisationRepo }

func (h *OrganisationsHandler) Contexts(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	items, err := h.repo.ListForUser(r.Context(), userID)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load account contexts")
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"individual": true, "organisations": items})
}

func (h *OrganisationsHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	item, err := h.repo.GetForUser(r.Context(), chi.URLParam(r, "organisationID"), userID)
	if err != nil {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Organisation profile not found or access denied")
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *OrganisationsHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	email, _ := r.Context().Value(models.CtxUserEmail).(string)
	var input struct {
		Token string `json:"token"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil || input.Token == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "token is required")
		return
	}
	if err := h.repo.AcceptInvitation(r.Context(), input.Token, userID, email); err != nil {
		response.Err(w, http.StatusBadRequest, "INVITATION_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Organisation invitation accepted")
}
