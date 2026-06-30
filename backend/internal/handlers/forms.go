package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
	"github.com/go-chi/chi/v5"
)

type FormsHandler struct {
	svc *services.FormService
}

// isSuperAdmin checks the user's roles attached to the request context.
func isSuperAdmin(r *http.Request) bool {
	roles, _ := r.Context().Value(models.CtxUserRoles).([]string)
	for _, role := range roles {
		if role == "super_admin" {
			return true
		}
	}
	return false
}

// ── Departments ──────────────────────────────────────────────────

// GET /api/v1/departments — read-only list, available to all authed users
func (h *FormsHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	depts, err := h.svc.ListDepartments(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, depts)
}

// ── Admin: form templates ────────────────────────────────────────

type formTemplateRequest struct {
	DepartmentID   string                 `json:"department_id"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	BannerImageURL string                 `json:"banner_image_url"`
	PriceUGX       float64                `json:"price_ugx"`
	Status         string                 `json:"status"`
	Fields         []formTemplateReqField `json:"fields"`
}

type formTemplateReqField struct {
	ID          string          `json:"id"`
	FieldKey    string          `json:"field_key"`
	FieldType   string          `json:"field_type"`
	Label       string          `json:"label"`
	Placeholder string          `json:"placeholder"`
	HelpText    string          `json:"help_text"`
	IsRequired  bool            `json:"is_required"`
	Config      json.RawMessage `json:"config"`
}

func (req formTemplateRequest) toServiceInput(includeFields bool) services.SaveTemplateInput {
	in := services.SaveTemplateInput{
		DepartmentID:   req.DepartmentID,
		Title:          req.Title,
		Description:    req.Description,
		BannerImageURL: req.BannerImageURL,
		PriceUGX:       req.PriceUGX,
		Status:         req.Status,
	}
	if includeFields {
		in.Fields = make([]*models.FormField, len(req.Fields))
		for i, f := range req.Fields {
			in.Fields[i] = &models.FormField{
				FieldKey:    f.FieldKey,
				FieldType:   f.FieldType,
				Label:       f.Label,
				Placeholder: f.Placeholder,
				HelpText:    f.HelpText,
				IsRequired:  f.IsRequired,
				Config:      f.Config,
			}
		}
	}
	return in
}

// GET /api/v1/admin/forms
func (h *FormsHandler) AdminListTemplates(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	out, err := h.svc.ListAdminTemplates(r.Context(), userID, isSuperAdmin(r), r.URL.Query().Get("status"))
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, out)
}

// POST /api/v1/admin/forms
func (h *FormsHandler) AdminCreateTemplate(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	var req formTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	t, err := h.svc.CreateTemplate(r.Context(), userID, isSuperAdmin(r), req.toServiceInput(true))
	if err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, t)
}

// GET /api/v1/admin/forms/{id}
func (h *FormsHandler) AdminGetTemplate(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	t, err := h.svc.GetAdminTemplate(r.Context(), id, userID, isSuperAdmin(r))
	if err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, t)
}

// PUT /api/v1/admin/forms/{id}
func (h *FormsHandler) AdminUpdateTemplate(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	var req formTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	t, err := h.svc.UpdateTemplate(r.Context(), id, userID, isSuperAdmin(r), req.toServiceInput(true))
	if err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, t)
}

// DELETE /api/v1/admin/forms/{id}
func (h *FormsHandler) AdminDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteTemplate(r.Context(), id, userID, isSuperAdmin(r)); err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSONMsg(w, http.StatusOK, "Form template archived")
}

// ── Admin: submissions (department-scoped) ───────────────────────

// GET /api/v1/admin/forms/submissions
// Optional ?template_id=&status=
func (h *FormsHandler) AdminListSubmissions(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	p := &models.PaginationParams{Page: page, PerPage: perPage}
	out, total, err := h.svc.ListSubmissions(r.Context(), userID, isSuperAdmin(r),
		q.Get("template_id"), q.Get("status"), p)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSONPaged(w, http.StatusOK, out, &response.Meta{Page: p.Page, PerPage: p.PerPage, Total: total})
}

// GET /api/v1/admin/forms/submissions/{id}
func (h *FormsHandler) AdminGetSubmission(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	s, err := h.svc.GetSubmissionForAdmin(r.Context(), id, userID, isSuperAdmin(r))
	if err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, s)
}

// POST /api/v1/admin/forms/submissions/{id}/review  { status, notes }
func (h *FormsHandler) AdminReviewSubmission(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "status is required")
		return
	}
	if err := h.svc.Review(r.Context(), id, req.Status, reviewerID, req.Notes, isSuperAdmin(r)); err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSONMsg(w, http.StatusOK, "Submission reviewed")
}

// POST /api/v1/admin/forms/submissions/{id}/verify-payment
func (h *FormsHandler) AdminVerifySubmissionPayment(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	if err := h.svc.VerifySubmissionPayment(r.Context(), id, reviewerID, isSuperAdmin(r)); err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSONMsg(w, http.StatusOK, "Payment verified")
}

// ── Public portal ────────────────────────────────────────────────

// GET /api/v1/portal/forms/open
func (h *FormsHandler) PortalListOpen(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.ListPublicTemplates(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, out)
}

// GET /api/v1/portal/forms/{slug}
func (h *FormsHandler) PortalGetForm(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	t, err := h.svc.GetPublicTemplate(r.Context(), slug)
	if err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, t)
}

// POST /api/v1/portal/forms/{templateID}/draft  { answers }
func (h *FormsHandler) PortalSaveDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	templateID := chi.URLParam(r, "templateID")
	var req struct {
		Answers json.RawMessage `json:"answers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	sub, err := h.svc.SaveOrCreateDraft(r.Context(), userID, services.SaveDraftAnswersInput{
		TemplateID: templateID,
		Answers:    req.Answers,
	})
	if err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, sub)
}

// GET /api/v1/portal/submissions/{id}
func (h *FormsHandler) PortalGetSubmission(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	sub, err := h.svc.GetSubmissionForUser(r.Context(), id, userID)
	if err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, sub)
}

// POST /api/v1/portal/submissions/{id}/submit
func (h *FormsHandler) PortalSubmit(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	sub, err := h.svc.Submit(r.Context(), id, userID)
	if err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, sub)
}

// POST /api/v1/portal/submissions/{id}/payment-proof
func (h *FormsHandler) PortalUploadPaymentProof(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	var req struct {
		PaymentReference string  `json:"payment_reference"`
		AmountUGX        float64 `json:"payment_amount_ugx"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.PaymentReference == "" || req.AmountUGX <= 0 {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "payment_reference and payment_amount_ugx required")
		return
	}
	if err := h.svc.UploadPaymentProof(r.Context(), id, userID, req.PaymentReference, req.AmountUGX); err != nil {
		h.writeServiceErr(w, err)
		return
	}
	response.JSONMsg(w, http.StatusOK, "Payment proof recorded. Awaiting verification.")
}

// ── error mapper ─────────────────────────────────────────────────

func (h *FormsHandler) writeServiceErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		response.Err(w, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, services.ErrForbiddenDept):
		response.Err(w, http.StatusForbidden, "FORBIDDEN", err.Error())
	case errors.Is(err, services.ErrNotOwner):
		response.Err(w, http.StatusForbidden, "FORBIDDEN", "Not your submission")
	case errors.Is(err, services.ErrTemplateClosed):
		response.Err(w, http.StatusConflict, "TEMPLATE_CLOSED", err.Error())
	case errors.Is(err, services.ErrPaymentRequired):
		response.Err(w, http.StatusPaymentRequired, "PAYMENT_REQUIRED", err.Error())
	case errors.Is(err, services.ErrInvalidTransition):
		response.Err(w, http.StatusBadRequest, "INVALID_STATUS", err.Error())
	default:
		response.Err(w, http.StatusBadRequest, "FORM_ERROR", err.Error())
	}
}
