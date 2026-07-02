package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
	"github.com/atenimedia-llc/ncs-online/backend/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ApplicationsHandler struct {
	svc     *services.ApplicationService
	audit   *repository.AuditRepo
	storage *repository.CMSRepo
}

func (h *ApplicationsHandler) storageSettings(ctx context.Context) (storage.Settings, error) {
	if h.storage == nil {
		return storage.DefaultSettings(), nil
	}
	s, err := h.storage.GetSetting(ctx, storage.SettingsKey)
	if errors.Is(err, repository.ErrNotFound) {
		return storage.DefaultSettings(), nil
	}
	if err != nil {
		return storage.Settings{}, err
	}
	return storage.ParseSettings(s.Value)
}

// POST /api/v1/applications/{formType}/draft
func (h *ApplicationsHandler) SaveDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	formType := chi.URLParam(r, "formType")

	var req struct {
		FormData      json.RawMessage `json:"form_data"`
		LastSavedStep int             `json:"last_saved_step"`
		AppType       string          `json:"application_type"`
		OrgType       string          `json:"organisation_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	app, err := h.svc.SaveDraft(r.Context(), userID, formType, services.SaveDraftInput{
		FormData:      req.FormData,
		LastSavedStep: req.LastSavedStep,
		AppType:       req.AppType,
		OrgType:       req.OrgType,
	})
	if err != nil {
		response.Err(w, http.StatusBadRequest, "DRAFT_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, app)
}

// GET /api/v1/applications/{formType}/draft
func (h *ApplicationsHandler) GetDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	formType := chi.URLParam(r, "formType")
	app, err := h.svc.GetDraft(r.Context(), userID, formType)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NO_DRAFT", "No saved draft for this form type")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, app)
}

// DELETE /api/v1/applications/{formType}/draft
func (h *ApplicationsHandler) DeleteDraft(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	formType := chi.URLParam(r, "formType")
	if err := h.svc.DeleteDraft(r.Context(), userID, formType); err != nil {
		response.Err(w, http.StatusBadRequest, "DELETE_DRAFT_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Draft discarded")
}

// GET /api/v1/applications
func (h *ApplicationsHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	p := &models.PaginationParams{Page: page, PerPage: perPage}

	apps, total, err := h.svc.ListByUser(r.Context(), userID, p)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list applications")
		return
	}
	response.JSONPaged(w, http.StatusOK, apps, &response.Meta{Page: p.Page, PerPage: p.PerPage, Total: total})
}

// GET /api/v1/applications/{id}
func (h *ApplicationsHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	app, err := h.svc.GetByID(r.Context(), id, userID, false)
	if errors.Is(err, repository.ErrNotFound) || errors.Is(err, services.ErrNotOwner) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Application not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, app)
}

// POST /api/v1/applications/{id}/prefilled-pdf — stub
func (h *ApplicationsHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	response.Err(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "PDF generation is not yet implemented")
}

// POST /api/v1/applications/{id}/signed-form
func (h *ApplicationsHandler) UploadSignedForm(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	var req struct {
		SignedFormURL string `json:"signed_form_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SignedFormURL == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "signed_form_url is required")
		return
	}
	if err := h.svc.SetSignedForm(r.Context(), id, userID, req.SignedFormURL); err != nil {
		if errors.Is(err, services.ErrNotOwner) {
			response.Err(w, http.StatusForbidden, "FORBIDDEN", "Not your application")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Signed form uploaded. Application is now pending payment.")
}

// GET /api/v1/applications/{id}/signed-form
func (h *ApplicationsHandler) GetSignedForm(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	app, err := h.svc.GetByID(r.Context(), id, userID, false)
	if err != nil {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Application not found")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"signed_form_url": app.SignedFormURL})
}

// POST /api/v1/applications/{id}/payment/initiate — stub (awaiting gateway integration)
func (h *ApplicationsHandler) InitiatePayment(w http.ResponseWriter, r *http.Request) {
	response.Err(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Online payment gateway is not yet configured")
}

// POST /api/v1/applications/{id}/payment/proof
func (h *ApplicationsHandler) UploadPaymentProof(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	var req struct {
		PaymentDate      string  `json:"payment_date"`
		BankOrChannel    string  `json:"paying_bank_or_channel"`
		DepositorName    string  `json:"depositor_name"`
		PaymentReference string  `json:"payment_reference_number"`
		AmountUGX        float64 `json:"payment_amount_ugx"`
		ProofDocumentURL string  `json:"proof_document_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	errs := map[string]string{}
	if req.ProofDocumentURL == "" {
		errs["proof_document_url"] = "required"
	}
	if req.PaymentReference == "" {
		errs["payment_reference_number"] = "required"
	}
	if req.AmountUGX <= 0 {
		errs["payment_amount_ugx"] = "must be greater than 0"
	}
	if len(errs) > 0 {
		response.ValidationErr(w, errs)
		return
	}
	if err := h.svc.UploadPaymentProof(r.Context(), id, userID, services.PaymentProofInput{
		Method:    models.PaymentMethodProof,
		Reference: req.PaymentReference,
		Amount:    req.AmountUGX,
	}); err != nil {
		if errors.Is(err, services.ErrNotOwner) {
			response.Err(w, http.StatusForbidden, "FORBIDDEN", "Not your application")
			return
		}
		if errors.Is(err, services.ErrInvalidTransition) {
			response.Err(w, http.StatusBadRequest, "INVALID_STATUS", "Application is not in PENDING_PAYMENT status")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Payment proof uploaded. Awaiting reviewer verification.")
}

// GET /api/v1/applications/{id}/payment
func (h *ApplicationsHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	app, err := h.svc.GetByID(r.Context(), id, userID, false)
	if err != nil {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Application not found")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"payment_status":      app.PaymentStatus,
		"payment_method":      app.PaymentMethod,
		"payment_reference":   app.PaymentReference,
		"payment_amount_ugx":  app.PaymentAmountUGX,
		"payment_verified_at": app.PaymentVerifiedAt,
	})
}

// POST /api/v1/applications/{id}/submit
func (h *ApplicationsHandler) Submit(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	app, err := h.svc.Submit(r.Context(), id, userID)
	if errors.Is(err, services.ErrSubmitRequirements) {
		response.Err(w, http.StatusBadRequest, "SUBMIT_REQUIREMENTS", "Signed form and payment are required before submission")
		return
	}
	if errors.Is(err, services.ErrNotOwner) {
		response.Err(w, http.StatusForbidden, "FORBIDDEN", "Not your application")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, app)
}

// PATCH /api/v1/applications/{id}/respond
func (h *ApplicationsHandler) Respond(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	var req struct {
		FormData json.RawMessage `json:"form_data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if err := h.svc.Respond(r.Context(), id, userID, req.FormData); err != nil {
		response.Err(w, http.StatusBadRequest, "RESPOND_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Response submitted. Application resubmitted for review.")
}

// POST /api/v1/applications/{id}/attachments
func (h *ApplicationsHandler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	if _, err := h.svc.GetByID(r.Context(), id, userID, false); err != nil {
		if errors.Is(err, services.ErrNotOwner) {
			response.Err(w, http.StatusForbidden, "FORBIDDEN", "Not your application")
			return
		}
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Application not found")
		return
	}
	if err := r.ParseMultipartForm(24 << 20); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Could not parse form (max 24MB)")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "No file uploaded (field: 'file')")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{
		".pdf": true, ".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
		".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	}
	if !allowed[ext] {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "File type not allowed")
		return
	}
	cfg, err := h.storageSettings(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load storage settings")
		return
	}
	result, err := storage.NewUploader(cfg).Upload(r.Context(), storage.UploadInput{
		Scope:       storage.ScopeApplication,
		Reader:      io.LimitReader(file, 24<<20),
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
		Subdir:      id,
	})
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "UPLOAD_FAILED", err.Error())
		return
	}
	att := &models.Attachment{
		ID:        uuid.NewString(),
		FieldName: strings.TrimSpace(r.FormValue("field_name")),
		FileName:  filepath.Base(header.Filename),
		FileURL:   result.URL,
		FileSize:  header.Size,
		MimeType:  header.Header.Get("Content-Type"),
	}
	if att.FieldName == "" {
		att.FieldName = "attachment"
	}
	if err := h.svc.AddAttachment(r.Context(), id, userID, att); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save attachment record")
		return
	}
	response.JSON(w, http.StatusCreated, att)
}

// GET /api/v1/applications/{id}/attachments
func (h *ApplicationsHandler) ListAttachments(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	app, err := h.svc.GetByID(r.Context(), id, userID, false)
	if err != nil {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Application not found")
		return
	}
	response.JSON(w, http.StatusOK, app.Attachments)
}

// DELETE /api/v1/applications/{id}/attachments/{attachmentID}
func (h *ApplicationsHandler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	response.Err(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Secure attachment deletion is not yet configured")
}

// GET /api/v1/admin/applications
func (h *ApplicationsHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	p := &models.PaginationParams{Page: page, PerPage: perPage}
	apps, total, err := h.svc.AdminList(r.Context(), q.Get("status"), p)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list applications")
		return
	}
	response.JSONPaged(w, http.StatusOK, apps, &response.Meta{Page: p.Page, PerPage: p.PerPage, Total: total})
}

// GET /api/v1/admin/applications/{id}
func (h *ApplicationsHandler) AdminGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	app, err := h.svc.GetByID(r.Context(), id, "", true)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Application not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, app)
}

// POST /api/v1/admin/applications/{id}/approve
func (h *ApplicationsHandler) Approve(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	var req struct {
		Notes string `json:"notes"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.svc.Approve(r.Context(), id, reviewerID, req.Notes); err != nil {
		response.Err(w, http.StatusBadRequest, "APPROVE_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Application approved")
}

// POST /api/v1/admin/applications/{id}/reject
func (h *ApplicationsHandler) Reject(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	var req struct {
		Notes string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Notes == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "notes (rejection reason) is required")
		return
	}
	if err := h.svc.Reject(r.Context(), id, reviewerID, req.Notes); err != nil {
		response.Err(w, http.StatusBadRequest, "REJECT_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Application rejected")
}

// POST /api/v1/admin/applications/{id}/request-info
func (h *ApplicationsHandler) RequestInfo(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	var req struct {
		Notes string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Notes == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "notes is required")
		return
	}
	if err := h.svc.RequestInfo(r.Context(), id, reviewerID, req.Notes); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Application set to NEEDS_INFORMATION")
}

// POST /api/v1/admin/applications/{id}/verify-payment
func (h *ApplicationsHandler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := r.Context().Value(models.CtxUserID).(string)
	id := chi.URLParam(r, "id")
	if err := h.svc.VerifyPayment(r.Context(), id, reviewerID); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Payment verified")
}

// POST /api/v1/admin/applications/{id}/reject-payment
func (h *ApplicationsHandler) RejectPayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Notes string `json:"notes"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := h.svc.RejectPayment(r.Context(), id, req.Notes); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Payment proof rejected")
}

// GET /api/v1/transactions
func (h *ApplicationsHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	p := &models.PaginationParams{Page: page, PerPage: perPage}
	apps, total, err := h.svc.ListByUser(r.Context(), userID, p)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list transactions")
		return
	}
	// Return payment-focused projection
	type txn struct {
		ID               string   `json:"id"`
		Ref              string   `json:"application_reference"`
		FormType         string   `json:"form_type"`
		PaymentStatus    string   `json:"payment_status"`
		PaymentMethod    string   `json:"payment_method"`
		PaymentReference string   `json:"payment_reference"`
		AmountUGX        *float64 `json:"payment_amount_ugx"`
	}
	out := make([]txn, len(apps))
	for i, a := range apps {
		out[i] = txn{ID: a.ID, Ref: a.ApplicationReference, FormType: a.FormType,
			PaymentStatus: a.PaymentStatus, PaymentMethod: a.PaymentMethod,
			PaymentReference: a.PaymentReference, AmountUGX: a.PaymentAmountUGX}
	}
	response.JSONPaged(w, http.StatusOK, out, &response.Meta{Page: p.Page, PerPage: p.PerPage, Total: total})
}

// GET /api/v1/transactions/{reference}
func (h *ApplicationsHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	response.Err(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Lookup by reference coming soon")
}
