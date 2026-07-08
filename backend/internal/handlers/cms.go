package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"html"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/atenimedia-llc/ncs-online/backend/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CMSHandler struct {
	repo *repository.CMSRepo

	// In-memory, single-use state tokens for the Google Drive OAuth connect
	// flow. They only need to survive the few seconds between redirecting the
	// admin to Google and them coming back, so a DB table would be overkill.
	oauthStateMu sync.Mutex
	oauthStates  map[string]time.Time
}

// nullableID returns nil for "" / nil / whitespace, otherwise the trimmed value.
// Used so that "" coming from the JSON payload clears the FK instead of
// triggering a 23503 fkey violation on insert.
func nullableID(s *string) *string {
	if s == nil {
		return nil
	}
	v := strings.TrimSpace(*s)
	if v == "" {
		return nil
	}
	return &v
}

// GET /api/v1/cms/departments — institutional departments with staff counts.
// Drives the dropdowns in the careers + team editors and the team
// directory's department grouping.
func (h *CMSHandler) ListInstitutionalDepartments(w http.ResponseWriter, r *http.Request) {
	out, err := h.repo.ListDepartmentsWithStaffCount(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, out)
}

func notificationTheme(kind string) (string, string) {
	switch kind {
	case "new_comment":
		return "chat", "New Comment"
	case "contact_form":
		return "mail", "Contact Message"
	case "investment_request":
		return "money", "Investment Request"
	case "new_sign_in", "unusual_activity":
		return "shield-alert", "Security Activity"
	case "password_reset":
		return "key", "Password Reset"
	default:
		return "check", "System Notification"
	}
}

func inboundPayload(values map[string]interface{}) json.RawMessage {
	raw, err := json.Marshal(values)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return raw
}

func (h *CMSHandler) dispatchInbound(ctx context.Context, sourceType, sourceID string, payload map[string]interface{}, statusState, workflowStatus string) {
	_ = h.repo.UpsertInboundSubmission(ctx, &models.InboundSubmission{
		ID:             "inbound_" + sourceType + "_" + sourceID,
		SourceType:     sourceType,
		SourceID:       sourceID,
		Payload:        inboundPayload(payload),
		StatusState:    firstNonEmpty(statusState, "unread"),
		WorkflowStatus: firstNonEmpty(workflowStatus, "pending_review"),
		InternalNotes:  json.RawMessage(`[]`),
	})
}

func (h *CMSHandler) notifyAdmins(ctx context.Context, kind, message string) {
	icon, title := notificationTheme(kind)
	_ = h.repo.CreateNotification(ctx, &models.Notification{
		ID: uuid.NewString(), Type: kind, Title: title, Message: message,
		Status: "unread", IconKey: icon,
	})
}

var googleAnalyticsIDPattern = regexp.MustCompile(`^G-[A-Za-z0-9]+$`)

type thirdPartySettings struct {
	GoogleAnalyticsEnabled bool   `json:"google_analytics_enabled"`
	GoogleAnalyticsID      string `json:"google_analytics_id"`
}

func normalizeThirdPartySettings(s thirdPartySettings) thirdPartySettings {
	s.GoogleAnalyticsID = strings.ToUpper(strings.TrimSpace(s.GoogleAnalyticsID))
	if !s.GoogleAnalyticsEnabled {
		s.GoogleAnalyticsID = ""
	}
	return s
}

func validateThirdPartySettings(s thirdPartySettings) (thirdPartySettings, map[string]string) {
	s = normalizeThirdPartySettings(s)
	errs := map[string]string{}
	if s.GoogleAnalyticsEnabled && !googleAnalyticsIDPattern.MatchString(s.GoogleAnalyticsID) {
		errs["google_analytics_id"] = "must start with G- and contain only letters or numbers"
	}
	return s, errs
}

type captchaSettings struct {
	Provider                string  `json:"captcha_provider"`
	CloudflareSiteKey       string  `json:"cloudflare_site_key"`
	CloudflareSecretKey     string  `json:"cloudflare_secret_key,omitempty"`
	CloudflareSecretSaved   bool    `json:"cloudflare_secret_saved,omitempty"`
	RecaptchaSiteKey        string  `json:"recaptcha_site_key"`
	RecaptchaSecretKey      string  `json:"recaptcha_secret_key,omitempty"`
	RecaptchaSecretSaved    bool    `json:"recaptcha_secret_saved,omitempty"`
	RecaptchaScoreThreshold float64 `json:"recaptcha_score_threshold"`
}

func defaultCaptchaSettings() captchaSettings {
	return captchaSettings{Provider: "none", RecaptchaScoreThreshold: 0.5}
}

func normalizeCaptchaSettings(s captchaSettings) captchaSettings {
	s.Provider = strings.TrimSpace(s.Provider)
	if s.Provider == "" {
		s.Provider = "none"
	}
	if s.Provider != "none" && s.Provider != "cloudflare_turnstile" && s.Provider != "google_recaptcha" {
		s.Provider = "none"
	}
	s.CloudflareSiteKey = strings.TrimSpace(s.CloudflareSiteKey)
	s.CloudflareSecretKey = strings.TrimSpace(s.CloudflareSecretKey)
	s.RecaptchaSiteKey = strings.TrimSpace(s.RecaptchaSiteKey)
	s.RecaptchaSecretKey = strings.TrimSpace(s.RecaptchaSecretKey)
	if s.RecaptchaScoreThreshold < 0.1 || s.RecaptchaScoreThreshold > 1 {
		s.RecaptchaScoreThreshold = 0.5
	}
	s.CloudflareSecretSaved = s.CloudflareSecretKey != ""
	s.RecaptchaSecretSaved = s.RecaptchaSecretKey != ""
	return s
}

func publicCaptchaSettings(s captchaSettings) captchaSettings {
	s.CloudflareSecretSaved = s.CloudflareSecretKey != ""
	s.RecaptchaSecretSaved = s.RecaptchaSecretKey != ""
	s.CloudflareSecretKey = ""
	s.RecaptchaSecretKey = ""
	return s
}

func (h *CMSHandler) captchaSettings(ctx context.Context) captchaSettings {
	settings := defaultCaptchaSettings()
	s, err := h.repo.GetSetting(ctx, "captcha")
	if err != nil {
		return settings
	}
	_ = json.Unmarshal(s.Value, &settings)
	return normalizeCaptchaSettings(settings)
}

func (h *CMSHandler) validateCaptchaSettings(ctx context.Context, incoming captchaSettings) (captchaSettings, map[string]string) {
	existing := h.captchaSettings(ctx)
	incoming.Provider = strings.TrimSpace(incoming.Provider)
	if incoming.Provider == "" {
		incoming.Provider = "none"
	}
	if strings.TrimSpace(incoming.CloudflareSecretKey) == "" || strings.TrimSpace(incoming.CloudflareSecretKey) == "********" {
		incoming.CloudflareSecretKey = existing.CloudflareSecretKey
	}
	if strings.TrimSpace(incoming.RecaptchaSecretKey) == "" || strings.TrimSpace(incoming.RecaptchaSecretKey) == "********" {
		incoming.RecaptchaSecretKey = existing.RecaptchaSecretKey
	}
	incoming = normalizeCaptchaSettings(incoming)
	errs := map[string]string{}
	switch incoming.Provider {
	case "cloudflare_turnstile":
		if incoming.CloudflareSiteKey == "" {
			errs["cloudflare_site_key"] = "required"
		}
		if incoming.CloudflareSecretKey == "" {
			errs["cloudflare_secret_key"] = "required"
		}
	case "google_recaptcha":
		if incoming.RecaptchaSiteKey == "" {
			errs["recaptcha_site_key"] = "required"
		}
		if incoming.RecaptchaSecretKey == "" {
			errs["recaptcha_secret_key"] = "required"
		}
	}
	return incoming, errs
}

func (h *CMSHandler) verifyCaptcha(ctx context.Context, r *http.Request, token, action string) error {
	settings := h.captchaSettings(ctx)
	if settings.Provider == "none" {
		return nil
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("captcha token required")
	}
	client := &http.Client{Timeout: 3 * time.Second}
	form := url.Values{"response": {token}, "remoteip": {clientIP(r)}}
	var endpoint string
	switch settings.Provider {
	case "cloudflare_turnstile":
		if settings.CloudflareSecretKey == "" {
			return errors.New("captcha is not configured")
		}
		endpoint = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
		form.Set("secret", settings.CloudflareSecretKey)
	case "google_recaptcha":
		if settings.RecaptchaSecretKey == "" {
			return errors.New("captcha is not configured")
		}
		endpoint = "https://www.google.com/recaptcha/api/siteverify"
		form.Set("secret", settings.RecaptchaSecretKey)
	default:
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var out struct {
		Success bool     `json:"success"`
		Score   float64  `json:"score"`
		Action  string   `json:"action"`
		Errors  []string `json:"error-codes"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return err
	}
	if !out.Success {
		return errors.New("captcha verification failed")
	}
	if settings.Provider == "google_recaptcha" {
		if out.Score < settings.RecaptchaScoreThreshold {
			return errors.New("captcha score too low")
		}
		if action != "" && out.Action != "" && out.Action != action {
			return errors.New("captcha action mismatch")
		}
	}
	return nil
}

func (h *CMSHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	limit, offset := paginate(r)
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	items, total, err := h.repo.ListNotifications(r.Context(), userID, r.URL.Query().Get("status"), limit, offset)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list notifications")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"total": total, "items": items})
}

func (h *CMSHandler) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Status != "read" && req.Status != "unread" && req.Status != "dismissed" {
		response.ValidationErr(w, map[string]string{"status": "must be read, unread, or dismissed"})
		return
	}
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	if err := h.repo.UpdateNotificationStatus(r.Context(), chi.URLParam(r, "id"), userID, req.Status); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Notification not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update notification")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Notification updated")
}

func (h *CMSHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	if err := h.repo.UpdateNotificationStatus(r.Context(), chi.URLParam(r, "id"), userID, "dismissed"); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Notification not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not dismiss notification")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Notification dismissed")
}

func (h *CMSHandler) ClearNotifications(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	if err := h.repo.ClearNotifications(r.Context(), userID); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not clear notifications")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Notifications cleared")
}

func (h *CMSHandler) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	if err := h.repo.MarkAllNotificationsRead(r.Context(), userID); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not mark notifications read")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Notifications marked read")
}

func (h *CMSHandler) ListInboundSubmissions(w http.ResponseWriter, r *http.Request) {
	limit, offset := paginate(r)
	items, total, err := h.repo.ListInboundSubmissions(r.Context(), r.URL.Query().Get("source_type"), r.URL.Query().Get("status"), limit, offset)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list inbound submissions")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"total": total, "items": items})
}

func (h *CMSHandler) InboundSubmissionCounts(w http.ResponseWriter, r *http.Request) {
	counts, err := h.repo.InboundSubmissionCounts(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load inbound counters")
		return
	}
	response.JSON(w, http.StatusOK, counts)
}

func (h *CMSHandler) UpdateInboundSubmission(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StatusState    string `json:"status_state"`
		WorkflowStatus string `json:"workflow_status"`
		AssignToMe     bool   `json:"assign_to_me"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.StatusState != "" && req.StatusState != "unread" && req.StatusState != "read" && req.StatusState != "replied" && req.StatusState != "archived" {
		response.ValidationErr(w, map[string]string{"status_state": "must be unread, read, replied, or archived"})
		return
	}
	if req.WorkflowStatus != "" && !validInboundWorkflow(req.WorkflowStatus) {
		response.ValidationErr(w, map[string]string{"workflow_status": "unsupported workflow status"})
		return
	}
	adminID := ""
	if req.AssignToMe {
		adminID, _ = r.Context().Value(models.CtxUserID).(string)
	}
	if err := h.repo.UpdateInboundSubmission(r.Context(), chi.URLParam(r, "id"), req.StatusState, req.WorkflowStatus, adminID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Inbound submission not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update inbound submission")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Inbound submission updated")
}

func (h *CMSHandler) AddInboundSubmissionNote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Note string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	note := sanitizePlain(req.Note, 2000)
	if note == "" {
		response.ValidationErr(w, map[string]string{"note": "required"})
		return
	}
	adminID, _ := r.Context().Value(models.CtxUserID).(string)
	if err := h.repo.AddInboundSubmissionNote(r.Context(), chi.URLParam(r, "id"), adminID, note); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Inbound submission not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save note")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Note saved")
}

func validInboundWorkflow(status string) bool {
	switch status {
	case "pending_review", "approved", "under_negotiation", "declined", "spam", "trash":
		return true
	default:
		return false
	}
}

func (h *CMSHandler) ListContactMessages(w http.ResponseWriter, r *http.Request) {
	limit, offset := paginate(r)
	items, total, err := h.repo.ListContactMessages(r.Context(), r.URL.Query().Get("status"), limit, offset)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list messages")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"total": total, "items": items})
}

func (h *CMSHandler) CreateContactMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		Subject      string `json:"subject"`
		Message      string `json:"message"`
		Website      string `json:"website"`
		CaptchaToken string `json:"captcha_token"`
		CaptchaAction string `json:"captcha_action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Website) != "" {
		response.JSONMsg(w, http.StatusCreated, "Message submitted")
		return
	}
	if err := h.verifyCaptcha(r.Context(), r, req.CaptchaToken, firstNonEmpty(req.CaptchaAction, "contact_form")); err != nil {
		response.Err(w, http.StatusUnprocessableEntity, "CAPTCHA_FAILED", "Bot verification failed. Please try again.")
		return
	}
	errs := map[string]string{}
	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "required"
	}
	if !strings.Contains(req.Email, "@") {
		errs["email"] = "valid email required"
	}
	if strings.TrimSpace(req.Message) == "" {
		errs["message"] = "required"
	}
	if len(errs) > 0 {
		response.ValidationErr(w, errs)
		return
	}
	msg := &models.ContactMessage{
		ID: uuid.NewString(), Name: sanitizePlain(req.Name, 160), Email: sanitizePlain(req.Email, 220),
		Subject: sanitizePlain(firstNonEmpty(req.Subject, "Contact Us Message"), 220),
		Message: sanitizePlain(req.Message, 5000), Status: "unread",
	}
	if err := h.repo.CreateContactMessage(r.Context(), msg); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not submit message")
		return
	}
	h.dispatchInbound(r.Context(), "contact_form", msg.ID, map[string]interface{}{
		"message_id": msg.ID,
		"name":       msg.Name,
		"email":      msg.Email,
		"subject":    msg.Subject,
		"message":    msg.Message,
	}, "unread", "pending_review")
	h.notifyAdmins(r.Context(), "contact_form", "New public contact message received.")
	response.JSON(w, http.StatusCreated, msg)
}

func (h *CMSHandler) CreateInvestmentRequest(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if website, _ := req["website"].(string); strings.TrimSpace(website) != "" {
		response.JSONMsg(w, http.StatusCreated, "Investment request submitted")
		return
	}
	if err := h.verifyCaptcha(r.Context(), r, stringFromMap(req, "captcha_token", "captchaToken"), firstNonEmpty(stringFromMap(req, "captcha_action", "captchaAction"), "investment_request")); err != nil {
		response.Err(w, http.StatusUnprocessableEntity, "CAPTCHA_FAILED", "Bot verification failed. Please try again.")
		return
	}
	name := sanitizePlain(stringFromMap(req, "name", "full_name", "investor_name"), 180)
	email := sanitizePlain(stringFromMap(req, "email", "investor_email"), 220)
	message := sanitizePlain(stringFromMap(req, "message", "proposal", "description"), 5000)
	if name == "" || !strings.Contains(email, "@") || message == "" {
		response.ValidationErr(w, map[string]string{"request": "name, valid email, and proposal/message are required"})
		return
	}
	id := uuid.NewString()
	payload := map[string]interface{}{}
	for key, value := range req {
		if key == "website" || key == "captcha_token" || key == "captchaToken" || key == "captcha_action" || key == "captchaAction" {
			continue
		}
		switch typed := value.(type) {
		case string:
			payload[key] = sanitizePlain(typed, 5000)
		default:
			payload[key] = typed
		}
	}
	payload["id"] = id
	payload["name"] = name
	payload["email"] = email
	h.dispatchInbound(r.Context(), "investment_request", id, payload, "unread", "pending_review")
	h.notifyAdmins(r.Context(), "investment_request", "New investment or funding request received.")
	response.JSON(w, http.StatusCreated, map[string]string{"id": id, "status": "submitted"})
}

func stringFromMap(values map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value, ok := values[key]; ok {
			if text := strings.TrimSpace(toString(value)); text != "" {
				return text
			}
		}
	}
	return ""
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strings.TrimRight(strings.TrimRight(strconv.FormatFloat(v, 'f', -1, 64), "0"), ".")
	default:
		raw, _ := json.Marshal(v)
		return string(raw)
	}
}

func (h *CMSHandler) UpdateContactMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Status != "read" && req.Status != "unread" && req.Status != "replied" {
		response.ValidationErr(w, map[string]string{"status": "must be read, unread, or replied"})
		return
	}
	if err := h.repo.UpdateContactMessageStatus(r.Context(), chi.URLParam(r, "id"), req.Status); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Message not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update message")
		return
	}
	inboundState := req.Status
	if inboundState == "replied" {
		inboundState = "replied"
	}
	_ = h.repo.UpdateInboundSubmission(r.Context(), "inbound_contact_form_"+chi.URLParam(r, "id"), inboundState, "", "")
	response.JSONMsg(w, http.StatusOK, "Message updated")
}

func (h *CMSHandler) DeleteContactMessage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeleteContactMessage(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Message not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete message")
		return
	}
	_ = h.repo.UpdateInboundSubmission(r.Context(), "inbound_contact_form_"+id, "archived", "trash", "")
	response.JSONMsg(w, http.StatusOK, "Message deleted")
}

func (h *CMSHandler) ClearContactMessages(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.ClearContactMessages(r.Context()); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not clear messages")
		return
	}
	_ = h.repo.ArchiveInboundBySource(r.Context(), "contact_form", "trash")
	response.JSONMsg(w, http.StatusOK, "Messages cleared")
}

func (h *CMSHandler) CreateInstitutionalDepartment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		response.ValidationErr(w, map[string]string{"name": "required"})
		return
	}
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	if code == "" {
		code = strings.ToUpper(strings.ReplaceAll(toSlug(firstNonEmpty(req.Slug, req.Name)), "-", "_"))
	}
	dept := &repository.DepartmentWithCount{ID: uuid.NewString(), Name: req.Name, Code: code, Description: req.Description}
	if err := h.repo.CreateDepartment(r.Context(), dept); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create department")
		return
	}
	response.JSON(w, http.StatusCreated, dept)
}

func (h *CMSHandler) UpdateInstitutionalDepartment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		response.ValidationErr(w, map[string]string{"name": "required"})
		return
	}
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	if code == "" {
		code = strings.ToUpper(strings.ReplaceAll(toSlug(firstNonEmpty(req.Slug, req.Name)), "-", "_"))
	}
	dept := &repository.DepartmentWithCount{ID: chi.URLParam(r, "id"), Name: req.Name, Code: code, Description: req.Description}
	if err := h.repo.UpdateDepartment(r.Context(), dept); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Department not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update department")
		return
	}
	response.JSON(w, http.StatusOK, dept)
}

func (h *CMSHandler) DeleteInstitutionalDepartment(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteDepartment(r.Context(), chi.URLParam(r, "id")); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Department not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete department")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Department deleted")
}

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)
var manualSlugRe = regexp.MustCompile(`^[a-z0-9-_]+$`)
var scriptTagRe = regexp.MustCompile(`(?is)<\s*(script|iframe|object|embed|style)[^>]*>.*?<\s*/\s*(script|iframe|object|embed|style)\s*>`)
var eventAttrRe = regexp.MustCompile(`(?i)\s+on[a-z]+\s*=\s*(".*?"|'.*?'|[^\s>]+)`)

// Go's RE2 engine has no backreferences, so this can't require the closing
// quote to match the opening one (\2 in PCRE/JS). Matching either quote style
// regardless of symmetry is still safe here: we only ever strip, never keep.
var unsafeHrefRe = regexp.MustCompile(`(?i)(href|src)\s*=\s*['"]?\s*javascript:[^'"\s>]*`)

func toSlug(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	return strings.Trim(slugRe.ReplaceAllString(b.String(), "-"), "-")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func sanitizePlain(s string, max int) string {
	s = strings.TrimSpace(html.EscapeString(s))
	if max > 0 && len(s) > max {
		return s[:max]
	}
	return s
}

func sanitizeRichText(s string) string {
	s = strings.TrimSpace(s)
	s = scriptTagRe.ReplaceAllString(s, "")
	s = eventAttrRe.ReplaceAllString(s, "")
	s = unsafeHrefRe.ReplaceAllString(s, "")
	return s
}

func normalizePostStatus(status string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "", "draft":
		return "draft", true
	case "approved":
		return "approved", true
	case "published":
		return "published", true
	default:
		return "", false
	}
}

func paginate(r *http.Request) (limit, offset int) {
	p := &models.PaginationParams{}
	if v := r.URL.Query().Get("page"); v != "" {
		for _, c := range v {
			if c >= '0' && c <= '9' {
				p.Page = p.Page*10 + int(c-'0')
			}
		}
	}
	if v := r.URL.Query().Get("per_page"); v != "" {
		for _, c := range v {
			if c >= '0' && c <= '9' {
				p.PerPage = p.PerPage*10 + int(c-'0')
			}
		}
	}
	p.Offset()
	if p.PerPage == 0 {
		p.PerPage = 20
	}
	return p.PerPage, p.Offset()
}

// ── Posts ─────────────────────────────────────────────────────────

// GET /api/v1/cms/posts?category=blog&status=published&page=1
func (h *CMSHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	status := r.URL.Query().Get("status")
	limit, offset := paginate(r)
	posts, total, err := h.repo.ListPosts(r.Context(), category, status, limit, offset)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list posts")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"total": total, "items": posts})
}

// GET /api/v1/cms/posts/{slug}
func (h *CMSHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	post, err := h.repo.GetPostBySlug(r.Context(), slug)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get post")
		return
	}
	response.JSON(w, http.StatusOK, post)
}

// POST /api/v1/admin/cms/posts
func (h *CMSHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	authorID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		Title              string          `json:"title"`
		Content            string          `json:"content"`
		Excerpt            string          `json:"excerpt"`
		Category           string          `json:"category"`
		CategoryTag        string          `json:"category_tag"`
		Status             string          `json:"status"`
		CoverImageURL      string          `json:"cover_image_url"`
		BreadcrumbImageURL string          `json:"breadcrumb_image_url"`
		PageBuilder        json.RawMessage `json:"page_builder"`
		Slug               *string         `json:"slug"`
		MetaTitle          string          `json:"meta_title"`
		MetaDescription    string          `json:"meta_description"`
		FocusKeywords      string          `json:"focus_keywords"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Title == "" {
		response.ValidationErr(w, map[string]string{"title": "required"})
		return
	}
	slug := req.Title
	if req.Slug != nil && *req.Slug != "" {
		slug = *req.Slug
	}
	if req.Category == "" {
		req.Category = "blog"
	}
	if len(req.PageBuilder) > 0 && json.Valid(req.PageBuilder) {
		req.Content = string(req.PageBuilder)
	}
	if req.BreadcrumbImageURL != "" {
		req.CoverImageURL = req.BreadcrumbImageURL
	}
	status, ok := normalizePostStatus(req.Status)
	if !ok {
		response.ValidationErr(w, map[string]string{"status": "must be draft, approved, or published"})
		return
	}

	now := time.Now()
	var publishedAt *time.Time
	var approvedAt *time.Time
	if status == "published" {
		publishedAt = &now
	}
	if status == "approved" || status == "published" {
		approvedAt = &now
	}
	if req.Slug != nil && *req.Slug != "" && !manualSlugRe.MatchString(*req.Slug) {
		response.ValidationErr(w, map[string]string{"slug": "must match ^[a-z0-9-_]+$"})
		return
	}

	post := &models.CMSPost{
		ID:              uuid.NewString(),
		Title:           sanitizePlain(req.Title, 180),
		Slug:            toSlug(slug),
		Content:         sanitizeRichText(req.Content),
		Excerpt:         sanitizePlain(req.Excerpt, 500),
		Category:        sanitizePlain(req.Category, 80),
		CategoryTag:     sanitizePlain(req.CategoryTag, 80),
		Status:          status,
		CoverImageURL:   sanitizePlain(req.CoverImageURL, 500),
		AuthorID:        &authorID,
		PublishedAt:     publishedAt,
		ApprovedAt:      approvedAt,
		MetaTitle:       sanitizePlain(req.MetaTitle, 180),
		MetaDescription: sanitizePlain(req.MetaDescription, 160),
		FocusKeywords:   sanitizePlain(req.FocusKeywords, 240),
	}
	if err := h.repo.CreatePost(r.Context(), post); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_SLUG", "A post with this slug already exists")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create post")
		return
	}
	response.JSON(w, http.StatusCreated, post)
}

// PUT /api/v1/admin/cms/posts/{id}
func (h *CMSHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := h.repo.GetPostByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get post")
		return
	}

	var req struct {
		Title              string          `json:"title"`
		Content            string          `json:"content"`
		Excerpt            string          `json:"excerpt"`
		Category           string          `json:"category"`
		CategoryTag        *string         `json:"category_tag"`
		Status             string          `json:"status"`
		CoverImageURL      string          `json:"cover_image_url"`
		BreadcrumbImageURL string          `json:"breadcrumb_image_url"`
		PageBuilder        json.RawMessage `json:"page_builder"`
		Slug               *string         `json:"slug"`
		MetaTitle          string          `json:"meta_title"`
		MetaDescription    string          `json:"meta_description"`
		FocusKeywords      string          `json:"focus_keywords"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}

	if req.Title != "" {
		post.Title = sanitizePlain(req.Title, 180)
	}
	if len(req.PageBuilder) > 0 && json.Valid(req.PageBuilder) {
		req.Content = string(req.PageBuilder)
	}
	if req.Content != "" {
		post.Content = sanitizeRichText(req.Content)
	}
	if req.Excerpt != "" {
		post.Excerpt = sanitizePlain(req.Excerpt, 500)
	}
	if req.Category != "" {
		post.Category = sanitizePlain(req.Category, 80)
	}
	if req.CategoryTag != nil {
		post.CategoryTag = sanitizePlain(*req.CategoryTag, 80)
	}
	if req.BreadcrumbImageURL != "" {
		req.CoverImageURL = req.BreadcrumbImageURL
	}
	post.CoverImageURL = sanitizePlain(req.CoverImageURL, 500)
	if req.Slug != nil && *req.Slug != "" {
		if !manualSlugRe.MatchString(*req.Slug) {
			response.ValidationErr(w, map[string]string{"slug": "must match ^[a-z0-9-_]+$"})
			return
		}
		post.Slug = toSlug(*req.Slug)
	}
	post.MetaTitle = sanitizePlain(req.MetaTitle, 180)
	post.MetaDescription = sanitizePlain(req.MetaDescription, 160)
	post.FocusKeywords = sanitizePlain(req.FocusKeywords, 240)
	if req.Status != "" {
		prevStatus := post.Status
		status, ok := normalizePostStatus(req.Status)
		if !ok {
			response.ValidationErr(w, map[string]string{"status": "must be draft, approved, or published"})
			return
		}
		post.Status = status
		if prevStatus != "approved" && (status == "approved" || status == "published") && post.ApprovedAt == nil {
			now := time.Now()
			post.ApprovedAt = &now
		}
		if prevStatus != "published" && status == "published" && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
	}

	if err := h.repo.UpdatePost(r.Context(), post); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_SLUG", "A post with this slug already exists")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update post")
		return
	}
	response.JSON(w, http.StatusOK, post)
}

// DELETE /api/v1/admin/cms/posts/{id}
func (h *CMSHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeletePost(r.Context(), id); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete post")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Post deleted")
}

func (h *CMSHandler) ListBlogCategories(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") != "false"
	contentType := strings.TrimSpace(r.URL.Query().Get("content_type"))
	if contentType == "" {
		contentType = "blog"
	}
	cats, err := h.repo.ListBlogCategories(r.Context(), activeOnly, contentType)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list blog categories")
		return
	}
	response.JSON(w, http.StatusOK, cats)
}

func (h *CMSHandler) CreateBlogCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		ContentType string `json:"content_type"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		response.ValidationErr(w, map[string]string{"name": "required"})
		return
	}
	slug := req.Slug
	if slug == "" {
		slug = toSlug(req.Name)
	}
	if !manualSlugRe.MatchString(slug) {
		response.ValidationErr(w, map[string]string{"slug": "must match ^[a-z0-9-_]+$"})
		return
	}
	contentType := strings.TrimSpace(req.ContentType)
	if contentType == "" {
		contentType = "blog"
	}
	cat := &models.BlogCategory{ID: uuid.NewString(), Name: sanitizePlain(req.Name, 120), Slug: slug, Description: sanitizePlain(req.Description, 500), ContentType: contentType, SortOrder: req.SortOrder, IsActive: req.IsActive}
	if err := h.repo.CreateBlogCategory(r.Context(), cat); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_SLUG", "A category with this slug already exists for this section")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create blog category")
		return
	}
	response.JSON(w, http.StatusCreated, cat)
}

func (h *CMSHandler) UpdateBlogCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		ContentType string `json:"content_type"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Name) == "" || !manualSlugRe.MatchString(req.Slug) {
		response.ValidationErr(w, map[string]string{"name": "required", "slug": "must match ^[a-z0-9-_]+$"})
		return
	}
	contentType := strings.TrimSpace(req.ContentType)
	if contentType == "" {
		contentType = "blog"
	}
	cat := &models.BlogCategory{ID: id, Name: sanitizePlain(req.Name, 120), Slug: req.Slug, Description: sanitizePlain(req.Description, 500), ContentType: contentType, SortOrder: req.SortOrder, IsActive: req.IsActive}
	if err := h.repo.UpdateBlogCategory(r.Context(), cat); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Category not found")
			return
		}
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_SLUG", "A category with this slug already exists for this section")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update blog category")
		return
	}
	response.JSON(w, http.StatusOK, cat)
}

func (h *CMSHandler) DeleteBlogCategory(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteBlogCategory(r.Context(), chi.URLParam(r, "id")); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Category not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete blog category")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Category deleted")
}

func (h *CMSHandler) ListPostComments(w http.ResponseWriter, r *http.Request) {
	post, err := h.repo.GetPostBySlug(r.Context(), chi.URLParam(r, "slug"))
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get post")
		return
	}
	comments, err := h.repo.ListApprovedCommentsForPost(r.Context(), post.ID)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list comments")
		return
	}
	response.JSON(w, http.StatusOK, nestComments(comments))
}

func (h *CMSHandler) SubmitPostComment(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	post, err := h.repo.GetPostBySlug(r.Context(), chi.URLParam(r, "slug"))
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Post not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get post")
		return
	}
	var req struct {
		Body          string  `json:"body"`
		ParentID      *string `json:"parent_id"`
		Website       string  `json:"website"`
		CaptchaToken  string  `json:"captcha_token"`
		CaptchaAction string  `json:"captcha_action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Website) != "" {
		response.JSONMsg(w, http.StatusCreated, "Comment submitted")
		return
	}
	if err := h.verifyCaptcha(r.Context(), r, req.CaptchaToken, firstNonEmpty(req.CaptchaAction, "blog_comment")); err != nil {
		response.Err(w, http.StatusUnprocessableEntity, "CAPTCHA_FAILED", "Bot verification failed. Please try again.")
		return
	}
	body := sanitizePlain(req.Body, 1500)
	if len(body) < 3 {
		response.ValidationErr(w, map[string]string{"body": "must be at least 3 characters"})
		return
	}
	depth := 0
	if req.ParentID != nil && strings.TrimSpace(*req.ParentID) != "" {
		parent, err := h.repo.GetCommentByID(r.Context(), *req.ParentID)
		if err != nil || parent.PostID != post.ID || parent.Depth >= 1 {
			response.ValidationErr(w, map[string]string{"parent_id": "replies are limited to two levels"})
			return
		}
		depth = parent.Depth + 1
	} else {
		req.ParentID = nil
	}
	comment := &models.BlogComment{ID: uuid.NewString(), PostID: post.ID, UserID: &userID, ParentID: req.ParentID, Body: body, Status: "pending", Depth: depth, IPAddress: r.RemoteAddr, UserAgent: sanitizePlain(r.UserAgent(), 500)}
	if err := h.repo.CreateComment(r.Context(), comment); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not submit comment")
		return
	}
	h.dispatchInbound(r.Context(), "blog_comment", comment.ID, map[string]interface{}{
		"comment_id":     comment.ID,
		"post_id":        post.ID,
		"post_title":     post.Title,
		"post_slug":      post.Slug,
		"comment_body":   comment.Body,
		"comment_status": comment.Status,
		"user_id":        userID,
		"ip_address":     comment.IPAddress,
		"user_agent":     comment.UserAgent,
	}, "unread", "pending_review")
	h.notifyAdmins(r.Context(), "new_comment", "A public blog comment is waiting for moderation.")
	response.JSON(w, http.StatusCreated, comment)
}

func (h *CMSHandler) ListCommentsModeration(w http.ResponseWriter, r *http.Request) {
	limit, offset := paginate(r)
	comments, total, err := h.repo.ListComments(r.Context(), r.URL.Query().Get("status"), limit, offset)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list comments")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"total": total, "items": comments})
}

func (h *CMSHandler) ApproveComment(w http.ResponseWriter, r *http.Request) {
	h.moderateComment(w, r, "approved")
}
func (h *CMSHandler) FlagComment(w http.ResponseWriter, r *http.Request) {
	h.moderateComment(w, r, "flagged")
}

func (h *CMSHandler) moderateComment(w http.ResponseWriter, r *http.Request, status string) {
	id := chi.URLParam(r, "id")
	if err := h.repo.ModerateComment(r.Context(), id, status); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Comment not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not moderate comment")
		return
	}
	workflowStatus := "approved"
	statusState := "read"
	if status == "flagged" {
		workflowStatus = "spam"
		statusState = "archived"
	}
	_ = h.repo.UpdateInboundSubmission(r.Context(), "inbound_blog_comment_"+id, statusState, workflowStatus, "")
	response.JSON(w, http.StatusOK, map[string]string{"status": status})
}

func (h *CMSHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeleteComment(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, http.StatusNotFound, "NOT_FOUND", "Comment not found")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete comment")
		return
	}
	_ = h.repo.UpdateInboundSubmission(r.Context(), "inbound_blog_comment_"+id, "archived", "trash", "")
	response.JSONMsg(w, http.StatusOK, "Comment deleted")
}

func nestComments(comments []*models.BlogComment) []*models.BlogComment {
	byID := map[string]*models.BlogComment{}
	roots := []*models.BlogComment{}
	for _, c := range comments {
		c.Replies = []*models.BlogComment{}
		byID[c.ID] = c
	}
	for _, c := range comments {
		if c.ParentID != nil {
			if parent, ok := byID[*c.ParentID]; ok {
				parent.Replies = append(parent.Replies, c)
				continue
			}
		}
		roots = append(roots, c)
	}
	return roots
}

// ── Events ────────────────────────────────────────────────────────

// GET /api/v1/cms/events?status=published
func (h *CMSHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limit, offset := paginate(r)
	events, total, err := h.repo.ListEvents(r.Context(), status, limit, offset)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list events")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"total": total, "items": events})
}

// GET /api/v1/cms/events/{slug}
func (h *CMSHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	event, err := h.repo.GetEventBySlug(r.Context(), slug)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Event not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get event")
		return
	}
	response.JSON(w, http.StatusOK, event)
}

// POST /api/v1/admin/cms/events
func (h *CMSHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	authorID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		Title         string     `json:"title"`
		Description   string     `json:"description"`
		Category      string     `json:"category"`
		Location      string     `json:"location"`
		EventDate     *time.Time `json:"event_date"`
		EndDate       *time.Time `json:"end_date"`
		CoverImageURL string     `json:"cover_image_url"`
		Status        string     `json:"status"`
		Slug          *string    `json:"slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Title == "" {
		response.ValidationErr(w, map[string]string{"title": "required"})
		return
	}
	slug := req.Title
	if req.Slug != nil && *req.Slug != "" {
		slug = *req.Slug
	}
	if req.Status == "" {
		req.Status = "draft"
	}
	event := &models.CMSEvent{
		ID:            uuid.NewString(),
		Title:         req.Title,
		Slug:          toSlug(slug),
		Description:   req.Description,
		Category:      req.Category,
		Location:      req.Location,
		EventDate:     req.EventDate,
		EndDate:       req.EndDate,
		CoverImageURL: req.CoverImageURL,
		Status:        req.Status,
		AuthorID:      &authorID,
	}
	if err := h.repo.CreateEvent(r.Context(), event); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_SLUG", "An event with this slug already exists")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create event")
		return
	}
	response.JSON(w, http.StatusCreated, event)
}

// PUT /api/v1/admin/cms/events/{id}
func (h *CMSHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	event, err := h.repo.GetEventByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Event not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get event")
		return
	}
	var req struct {
		Title         string     `json:"title"`
		Description   string     `json:"description"`
		Category      string     `json:"category"`
		Location      string     `json:"location"`
		EventDate     *time.Time `json:"event_date"`
		EndDate       *time.Time `json:"end_date"`
		CoverImageURL string     `json:"cover_image_url"`
		Status        string     `json:"status"`
		Slug          *string    `json:"slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.Category != "" {
		event.Category = req.Category
	}
	if req.Location != "" {
		event.Location = req.Location
	}
	if req.EventDate != nil {
		event.EventDate = req.EventDate
	}
	if req.EndDate != nil {
		event.EndDate = req.EndDate
	}
	if req.CoverImageURL != "" {
		event.CoverImageURL = req.CoverImageURL
	}
	if req.Status != "" {
		event.Status = req.Status
	}
	if req.Slug != nil && *req.Slug != "" {
		event.Slug = toSlug(*req.Slug)
	}
	if err := h.repo.UpdateEvent(r.Context(), event); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update event")
		return
	}
	response.JSON(w, http.StatusOK, event)
}

// DELETE /api/v1/admin/cms/events/{id}
func (h *CMSHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeleteEvent(r.Context(), id); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete event")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Event deleted")
}

// ── Careers ───────────────────────────────────────────────────────

// GET /api/v1/cms/careers?status=published&category=jobs
func (h *CMSHandler) ListCareers(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	category := r.URL.Query().Get("category")
	limit, offset := paginate(r)
	careers, total, err := h.repo.ListCareers(r.Context(), status, category, limit, offset)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list careers")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"total": total, "items": careers})
}

// GET /api/v1/cms/careers/{id}
func (h *CMSHandler) GetCareer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	career, err := h.repo.GetCareerByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Job posting not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get career")
		return
	}
	response.JSON(w, http.StatusOK, career)
}

// POST /api/v1/admin/cms/careers
func (h *CMSHandler) CreateCareer(w http.ResponseWriter, r *http.Request) {
	authorID, _ := r.Context().Value(models.CtxUserID).(string)
	var req struct {
		Title        string     `json:"title"`
		Department   string     `json:"department"`
		DepartmentID *string    `json:"department_id"`
		Location     string     `json:"location"`
		JobType      string     `json:"job_type"`
		Category     string     `json:"category"`
		Description  string     `json:"description"`
		Requirements string     `json:"requirements"`
		SalaryRange  string     `json:"salary_range"`
		Status       string     `json:"status"`
		DeadlineAt   *time.Time `json:"deadline_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Title == "" || req.Description == "" {
		errs := map[string]string{}
		if req.Title == "" {
			errs["title"] = "required"
		}
		if req.Description == "" {
			errs["description"] = "required"
		}
		response.ValidationErr(w, errs)
		return
	}
	if req.JobType == "" {
		req.JobType = "full_time"
	}
	if req.Category == "" {
		req.Category = "jobs"
	}
	if req.Status == "" {
		req.Status = "draft"
	}
	career := &models.CMSCareer{
		ID:           uuid.NewString(),
		Title:        req.Title,
		Department:   req.Department,
		DepartmentID: nullableID(req.DepartmentID),
		Location:     req.Location,
		JobType:      req.JobType,
		Category:     req.Category,
		Description:  req.Description,
		Requirements: req.Requirements,
		SalaryRange:  req.SalaryRange,
		Status:       req.Status,
		DeadlineAt:   req.DeadlineAt,
		AuthorID:     &authorID,
	}
	if err := h.repo.CreateCareer(r.Context(), career); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create career")
		return
	}
	response.JSON(w, http.StatusCreated, career)
}

// PUT /api/v1/admin/cms/careers/{id}
func (h *CMSHandler) UpdateCareer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	career, err := h.repo.GetCareerByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Career not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get career")
		return
	}
	var req struct {
		Title        string     `json:"title"`
		Department   string     `json:"department"`
		DepartmentID *string    `json:"department_id"`
		Location     string     `json:"location"`
		JobType      string     `json:"job_type"`
		Category     string     `json:"category"`
		Description  string     `json:"description"`
		Requirements string     `json:"requirements"`
		SalaryRange  string     `json:"salary_range"`
		Status       string     `json:"status"`
		DeadlineAt   *time.Time `json:"deadline_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if req.Title != "" {
		career.Title = req.Title
	}
	if req.Department != "" {
		career.Department = req.Department
	}
	if req.DepartmentID != nil {
		career.DepartmentID = nullableID(req.DepartmentID)
	}
	if req.Location != "" {
		career.Location = req.Location
	}
	if req.JobType != "" {
		career.JobType = req.JobType
	}
	if req.Category != "" {
		career.Category = req.Category
	}
	if req.Description != "" {
		career.Description = req.Description
	}
	if req.Requirements != "" {
		career.Requirements = req.Requirements
	}
	if req.SalaryRange != "" {
		career.SalaryRange = req.SalaryRange
	}
	if req.Status != "" {
		career.Status = req.Status
	}
	if req.DeadlineAt != nil {
		career.DeadlineAt = req.DeadlineAt
	}
	if err := h.repo.UpdateCareer(r.Context(), career); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update career")
		return
	}
	response.JSON(w, http.StatusOK, career)
}

// DELETE /api/v1/admin/cms/careers/{id}
func (h *CMSHandler) DeleteCareer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeleteCareer(r.Context(), id); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete career")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Career posting deleted")
}

// ── Slides ────────────────────────────────────────────────────────

// GET /api/v1/cms/slides
func (h *CMSHandler) ListSlides(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") != "false"
	slides, err := h.repo.ListSlides(r.Context(), activeOnly)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list slides")
		return
	}
	response.JSON(w, http.StatusOK, slides)
}

func (h *CMSHandler) GetSlideshow(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		slug = "homepage-hero"
	}
	show, err := h.repo.GetSlideshowBySlug(r.Context(), slug, r.URL.Query().Get("active") != "false")
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Slideshow not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get slideshow")
		return
	}
	response.JSON(w, http.StatusOK, show)
}

func (h *CMSHandler) UpdateSlideshow(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	show, err := h.repo.GetSlideshowBySlug(r.Context(), slug, false)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Slideshow not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get slideshow")
		return
	}
	var req slideshowReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		response.ValidationErr(w, map[string]string{"name": "required"})
		return
	}
	if req.TransitionEffect != "slide" {
		req.TransitionEffect = "fade"
	}
	if req.TransitionDuration < 150 || req.TransitionDuration > 5000 || req.AutoplaySpeed < 2000 || req.AutoplaySpeed > 30000 {
		response.ValidationErr(w, map[string]string{"timing": "duration must be 150-5000ms and autoplay 2000-30000ms"})
		return
	}
	updated := &models.CMSSlideshow{
		ID: show.ID, Slug: show.Slug, Name: sanitizePlain(req.Name, 120),
		TransitionEffect: req.TransitionEffect, TransitionDuration: req.TransitionDuration,
		AutoplaySpeed: req.AutoplaySpeed, PauseOnHover: req.PauseOnHover, IsActive: req.IsActive,
	}
	if err := h.repo.UpdateSlideshow(r.Context(), updated); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update slideshow")
		return
	}
	response.JSON(w, http.StatusOK, updated)
}

type slideReq struct {
	SlideshowID   string                   `json:"slideshow_id"`
	Title         string                   `json:"title"`
	Subtitle      string                   `json:"subtitle"`
	Description   string                   `json:"description"`
	ImageURL      string                   `json:"image_url"`
	MediaType     string                   `json:"media_type"`
	AnimationType string                   `json:"animation_type"`
	ButtonText    string                   `json:"button_text"`
	ButtonURL     string                   `json:"button_url"`
	Buttons       []*models.CMSSlideButton `json:"buttons"`
	SortOrder     int                      `json:"sort_order"`
	IsActive      bool                     `json:"is_active"`
}

var slideAnimationTypes = map[string]struct{}{
	"fade-in":     {},
	"slide-up":    {},
	"slide-left":  {},
	"slide-right": {},
	"zoom-in":     {},
	"zoom-out":    {},
	"flip-in":     {},
	"blur-in":     {},
	"bounce-in":   {},
	"ken-burns":   {},
}

type slideshowReq struct {
	Name               string `json:"name"`
	TransitionEffect   string `json:"transition_effect"`
	TransitionDuration int    `json:"transition_duration"`
	AutoplaySpeed      int    `json:"autoplay_speed"`
	PauseOnHover       bool   `json:"pause_on_hover"`
	IsActive           bool   `json:"is_active"`
}

func safeSlideURL(raw string) bool {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" || strings.HasPrefix(raw, "javascript:") || strings.HasPrefix(raw, "data:") {
		return false
	}
	return strings.HasPrefix(raw, "/") || strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "mailto:")
}

func normalizeSlidePayload(req *slideReq) (*models.CMSSlide, map[string]string) {
	errs := map[string]string{}
	if strings.TrimSpace(req.Title) == "" {
		errs["title"] = "required"
	}
	if req.SlideshowID == "" {
		req.SlideshowID = "homepage-hero"
	}
	if req.MediaType == "" {
		req.MediaType = "image"
	}
	if req.MediaType != "image" && req.MediaType != "video" {
		errs["media_type"] = "must be image or video"
	}
	if req.AnimationType == "" {
		req.AnimationType = "fade-in"
	}
	if _, ok := slideAnimationTypes[req.AnimationType]; !ok {
		errs["animation_type"] = "unsupported animation type"
	}
	buttons := req.Buttons
	if len(buttons) == 0 && (req.ButtonText != "" || req.ButtonURL != "") {
		buttons = []*models.CMSSlideButton{{Text: req.ButtonText, URL: req.ButtonURL, StyleClass: "primary", LinkTarget: "_self"}}
	}
	if len(buttons) > 2 {
		errs["buttons"] = "maximum of two buttons"
		buttons = buttons[:2]
	}
	for i, b := range buttons {
		b.Text = sanitizePlain(b.Text, 80)
		b.URL = strings.TrimSpace(b.URL)
		if b.Text == "" || !safeSlideURL(b.URL) {
			errs["buttons"] = "button text and safe URL are required"
		}
		if b.StyleClass != "primary" && b.StyleClass != "secondary" && b.StyleClass != "outline" {
			b.StyleClass = "primary"
		}
		if b.LinkTarget != "_blank" {
			b.LinkTarget = "_self"
		}
		b.SortOrder = i
	}
	s := &models.CMSSlide{
		SlideshowID:   sanitizePlain(req.SlideshowID, 80),
		Title:         sanitizePlain(req.Title, 180),
		Subtitle:      sanitizePlain(req.Subtitle, 180),
		Description:   sanitizePlain(req.Description, 500),
		ImageURL:      sanitizePlain(req.ImageURL, 600),
		MediaType:     req.MediaType,
		AnimationType: req.AnimationType,
		Buttons:       buttons,
		SortOrder:     req.SortOrder,
		IsActive:      req.IsActive,
	}
	if len(buttons) > 0 {
		s.ButtonText = buttons[0].Text
		s.ButtonURL = buttons[0].URL
	}
	return s, errs
}

// POST /api/v1/admin/cms/slides
func (h *CMSHandler) CreateSlide(w http.ResponseWriter, r *http.Request) {
	var req slideReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}
	s, errs := normalizeSlidePayload(&req)
	if len(errs) > 0 {
		response.ValidationErr(w, errs)
		return
	}
	s.ID = uuid.NewString()
	if err := h.repo.CreateSlide(r.Context(), s); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create slide")
		return
	}
	if err := h.repo.ReplaceSlideButtons(r.Context(), s.ID, s.Buttons); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save slide buttons")
		return
	}
	response.JSON(w, http.StatusCreated, s)
}

// PUT /api/v1/admin/cms/slides/{id}
func (h *CMSHandler) UpdateSlide(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req slideReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
		return
	}
	s, errs := normalizeSlidePayload(&req)
	if len(errs) > 0 {
		response.ValidationErr(w, errs)
		return
	}
	s.ID = id
	if err := h.repo.UpdateSlide(r.Context(), s); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update slide")
		return
	}
	if err := h.repo.ReplaceSlideButtons(r.Context(), s.ID, s.Buttons); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save slide buttons")
		return
	}
	response.JSON(w, http.StatusOK, s)
}

// DELETE /api/v1/admin/cms/slides/{id}
func (h *CMSHandler) DeleteSlide(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.DeleteSlide(r.Context(), id); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete slide")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Slide deleted")
}

// ── Menus ─────────────────────────────────────────────────────────

// GET /api/v1/cms/menus/{name}
func (h *CMSHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	menu, err := h.repo.GetMenu(r.Context(), name)
	if errors.Is(err, repository.ErrNotFound) {
		response.JSON(w, http.StatusOK, &models.CMSMenu{Name: name, Items: []models.CMSMenuItem{}})
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get menu")
		return
	}
	response.JSON(w, http.StatusOK, menu)
}

// PUT /api/v1/admin/cms/menus/{name}
func (h *CMSHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var items []models.CMSMenuItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid menu items")
		return
	}
	if err := h.repo.UpdateMenu(r.Context(), name, items); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update menu")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Menu updated")
}

// ── Settings (generic JSON key/value) ─────────────────────────────

// GET /api/v1/cms/settings/{key}
func (h *CMSHandler) GetSetting(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "captcha" {
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"key":   key,
			"value": publicCaptchaSettings(h.captchaSettings(r.Context())),
		})
		return
	}
	s, err := h.repo.GetSetting(r.Context(), key)
	if errors.Is(err, repository.ErrNotFound) {
		// Treat missing as empty value so the public site can render
		// defaults without per-request error handling.
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"key":   key,
			"value": map[string]interface{}{},
		})
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get setting")
		return
	}
	response.JSON(w, http.StatusOK, s)
}

// PUT /api/v1/admin/cms/settings/{key}
func (h *CMSHandler) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Body required")
		return
	}
	// Validate the body is well-formed JSON; store raw bytes.
	var probe interface{}
	if err := json.Unmarshal(body, &probe); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Body must be valid JSON")
		return
	}
	if key == "captcha" {
		var req captchaSettings
		if err := json.Unmarshal(body, &req); err != nil {
			response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid captcha settings")
			return
		}
		settings, errs := h.validateCaptchaSettings(r.Context(), req)
		if len(errs) > 0 {
			response.ValidationErr(w, errs)
			return
		}
		body, err = json.Marshal(settings)
		if err != nil {
			response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid captcha settings")
			return
		}
	}
	if key == "third_party" {
		var req thirdPartySettings
		if err := json.Unmarshal(body, &req); err != nil {
			response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid third-party settings")
			return
		}
		settings, errs := validateThirdPartySettings(req)
		if len(errs) > 0 {
			response.ValidationErr(w, errs)
			return
		}
		body, err = json.Marshal(settings)
		if err != nil {
			response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid third-party settings")
			return
		}
	}
	if err := h.repo.UpdateSetting(r.Context(), key, body); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save setting")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Setting saved")
}

// ── Media Upload ──────────────────────────────────────────────────

// GET /api/v1/admin/storage-settings
func (h *CMSHandler) GetStorageSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.storageSettings(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load storage settings")
		return
	}
	response.JSON(w, http.StatusOK, storageSettingsView(settings))
}

// PUT /api/v1/admin/storage-settings
func (h *CMSHandler) UpdateStorageSettings(w http.ResponseWriter, r *http.Request) {
	existing, err := h.storageSettings(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load existing storage settings")
		return
	}
	var req storage.Settings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	req = req.MergeSecrets(existing)
	req.Normalize()
	raw, err := json.Marshal(req)
	if err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid storage settings")
		return
	}
	if err := h.repo.UpdateSetting(r.Context(), storage.SettingsKey, raw); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save storage settings")
		return
	}
	response.JSON(w, http.StatusOK, storageSettingsView(req))
}

func (h *CMSHandler) storageSettings(ctx context.Context) (storage.Settings, error) {
	s, err := h.repo.GetSetting(ctx, storage.SettingsKey)
	if errors.Is(err, repository.ErrNotFound) {
		return storage.DefaultSettings(), nil
	}
	if err != nil {
		return storage.Settings{}, err
	}
	return storage.ParseSettings(s.Value)
}

type storageSettingsResponse struct {
	storage.Settings
	GoogleDriveConnected bool `json:"google_drive_connected"`
	SupabaseConfigured   bool `json:"supabase_configured"`
}

func storageSettingsView(s storage.Settings) storageSettingsResponse {
	return storageSettingsResponse{Settings: s.Redacted(), GoogleDriveConnected: s.GoogleDriveConnected(), SupabaseConfigured: s.SupabaseConfigured()}
}

// ── Google Drive OAuth connect flow ──────────────────────────────────────
//
// The admin only ever provides three things: the folder ID, the OAuth
// client ID, and the OAuth client secret. Google Drive uploads still need a
// refresh token under the hood, but instead of asking the admin to obtain
// one manually (e.g. via the OAuth playground), this flow gets it for them:
// they click "Connect", approve access on Google's consent screen, and
// Google redirects back into the SPA with a one-time code. The frontend
// posts that code (and the state we handed it) to the exchange endpoint
// below, which trades it for a refresh token and stores it server-side —
// it's never shown to or re-entered by the admin.

const googleDriveOAuthStateTTL = 10 * time.Minute

// POST /api/v1/admin/storage-settings/google-drive/connect
func (h *CMSHandler) ConnectGoogleDrive(w http.ResponseWriter, r *http.Request) {
	settings, err := h.storageSettings(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load storage settings")
		return
	}
	if strings.TrimSpace(settings.GoogleDriveClientID) == "" || strings.TrimSpace(settings.GoogleDriveClientSecret) == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Save the Google Drive client ID and client secret first")
		return
	}
	if strings.TrimSpace(settings.GoogleDriveFolderID) == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Save the Google Drive folder ID first")
		return
	}
	var body struct {
		RedirectURI string `json:"redirect_uri"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if strings.TrimSpace(body.RedirectURI) == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Missing redirect_uri")
		return
	}
	state, err := h.newOAuthState()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not start Google Drive connection")
		return
	}
	authURL := storage.GoogleDriveAuthURL(settings.GoogleDriveClientID, body.RedirectURI, state)
	response.JSON(w, http.StatusOK, map[string]string{"auth_url": authURL})
}

// POST /api/v1/admin/storage-settings/google-drive/exchange
// Called by the SPA after Google redirects the admin back with ?code=&state=.
func (h *CMSHandler) ExchangeGoogleDriveCode(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code        string `json:"code"`
		State       string `json:"state"`
		RedirectURI string `json:"redirect_uri"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if !h.consumeOAuthState(body.State) {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Connection request expired or already used, please try connecting again")
		return
	}
	if strings.TrimSpace(body.Code) == "" || strings.TrimSpace(body.RedirectURI) == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Missing code or redirect_uri")
		return
	}
	existing, err := h.storageSettings(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load storage settings")
		return
	}
	if strings.TrimSpace(existing.GoogleDriveClientID) == "" || strings.TrimSpace(existing.GoogleDriveClientSecret) == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Google Drive client ID and client secret are not configured")
		return
	}
	refreshToken, err := storage.ExchangeGoogleDriveCode(r.Context(), existing.GoogleDriveClientID, existing.GoogleDriveClientSecret, body.RedirectURI, body.Code)
	if err != nil {
		response.Err(w, http.StatusBadGateway, "GOOGLE_DRIVE_EXCHANGE_FAILED", err.Error())
		return
	}
	existing.GoogleDriveRefreshToken = refreshToken
	raw, err := json.Marshal(existing)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save storage settings")
		return
	}
	if err := h.repo.UpdateSetting(r.Context(), storage.SettingsKey, raw); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save storage settings")
		return
	}
	response.JSON(w, http.StatusOK, storageSettingsView(existing))
}

// POST /api/v1/admin/storage-settings/google-drive/disconnect
func (h *CMSHandler) DisconnectGoogleDrive(w http.ResponseWriter, r *http.Request) {
	existing, err := h.storageSettings(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load storage settings")
		return
	}
	existing.GoogleDriveRefreshToken = ""
	raw, err := json.Marshal(existing)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save storage settings")
		return
	}
	if err := h.repo.UpdateSetting(r.Context(), storage.SettingsKey, raw); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save storage settings")
		return
	}
	response.JSON(w, http.StatusOK, storageSettingsView(existing))
}

func (h *CMSHandler) newOAuthState() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	state := hex.EncodeToString(b)
	h.oauthStateMu.Lock()
	defer h.oauthStateMu.Unlock()
	if h.oauthStates == nil {
		h.oauthStates = map[string]time.Time{}
	}
	h.pruneOAuthStatesLocked()
	h.oauthStates[state] = time.Now().Add(googleDriveOAuthStateTTL)
	return state, nil
}

func (h *CMSHandler) consumeOAuthState(state string) bool {
	if state == "" {
		return false
	}
	h.oauthStateMu.Lock()
	defer h.oauthStateMu.Unlock()
	expiry, ok := h.oauthStates[state]
	delete(h.oauthStates, state)
	return ok && time.Now().Before(expiry)
}

func (h *CMSHandler) pruneOAuthStatesLocked() {
	now := time.Now()
	for k, v := range h.oauthStates {
		if now.After(v) {
			delete(h.oauthStates, k)
		}
	}
}

// POST /api/v1/admin/media/upload
func (h *CMSHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Could not parse form (max 32MB)")
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
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".webp": true, ".svg": true, ".pdf": true, ".mp4": true, ".webm": true,
	}
	if !allowed[ext] {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "File type not allowed. Allowed: jpg, jpeg, png, gif, webp, svg, pdf, mp4, webm")
		return
	}
	head := make([]byte, 512)
	n, _ := io.ReadFull(file, head)
	head = head[:n]
	detected := http.DetectContentType(head)
	if !allowedMagic(ext, detected, head) {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Uploaded file content does not match an allowed media type")
		return
	}

	subDir := "images"
	if ext == ".pdf" {
		subDir = "documents"
	} else if ext == ".mp4" || ext == ".webm" {
		subDir = "videos"
	}

	settings, err := h.storageSettings(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load storage settings")
		return
	}
	result, err := storage.NewUploader(settings).Upload(r.Context(), storage.UploadInput{
		Scope:       storage.ScopePublic,
		Reader:      io.MultiReader(bytes.NewReader(head), file),
		Filename:    header.Filename,
		ContentType: detected,
		Size:        header.Size,
		Subdir:      subDir,
	})
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "UPLOAD_FAILED", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

// ServeLocalUpload serves files written by the "local" storage provider.
// In production, nginx serves this same directory directly as a static
// alias (see nginx/nginx.conf's "location /uploads/" block) — this handler
// exists so local uploads are still reachable when running without nginx.
// It resolves the current path from live storage settings (not a fixed
// directory) so it keeps working if an admin changes the local upload path
// from the Storage Settings panel.
func (h *CMSHandler) ServeLocalUpload(w http.ResponseWriter, r *http.Request) {
	settings, err := h.storageSettings(r.Context())
	if err != nil {
		http.NotFound(w, r)
		return
	}
	rel := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, settings.LocalPublicURLPrefix), "/")
	if rel == "" {
		http.NotFound(w, r)
		return
	}
	absBase, err := filepath.Abs(settings.LocalPublicPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	target := filepath.Join(absBase, filepath.FromSlash(rel))
	if target != absBase && !strings.HasPrefix(target, absBase+string(filepath.Separator)) {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=2592000, immutable")
	http.ServeFile(w, r, target)
}

func allowedMagic(ext, detected string, head []byte) bool {
	switch ext {
	case ".jpg", ".jpeg":
		return detected == "image/jpeg"
	case ".png":
		return detected == "image/png"
	case ".gif":
		return detected == "image/gif"
	case ".webp":
		return len(head) >= 12 && string(head[0:4]) == "RIFF" && string(head[8:12]) == "WEBP"
	case ".svg":
		s := strings.TrimSpace(strings.ToLower(string(head)))
		return strings.HasPrefix(s, "<svg") || strings.Contains(s, "<svg")
	case ".pdf":
		return detected == "application/pdf"
	case ".mp4":
		return len(head) >= 12 && strings.Contains(string(head[4:12]), "ftyp")
	case ".webm":
		return len(head) >= 4 && head[0] == 0x1a && head[1] == 0x45 && head[2] == 0xdf && head[3] == 0xa3
	default:
		return false
	}
}

// ── Custom Fonts ──────────────────────────────────────────────────

// GET /api/v1/cms/fonts (public — the site needs the list to render @font-face)
// GET /api/v1/admin/cms/fonts (admin — same data, listed in the Font Manager)
func (h *CMSHandler) ListFonts(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.ListCustomFonts(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list fonts")
		return
	}
	response.JSON(w, http.StatusOK, items)
}

var allowedFontExt = map[string]string{
	".ttf": "ttf", ".otf": "otf", ".woff": "woff", ".woff2": "woff2",
}

// POST /api/v1/admin/cms/fonts
func (h *CMSHandler) UploadFont(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Could not parse form (max 32MB)")
		return
	}
	displayName := sanitizePlain(r.FormValue("display_name"), 120)
	if displayName == "" {
		response.ValidationErr(w, map[string]string{"display_name": "required"})
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "No file uploaded (field: 'file')")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	format, ok := allowedFontExt[ext]
	if !ok {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "File type not allowed. Allowed: ttf, otf, woff, woff2")
		return
	}
	head := make([]byte, 12)
	n, _ := io.ReadFull(file, head)
	head = head[:n]
	if !allowedFontMagic(format, head) {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Uploaded file content does not match an allowed font type")
		return
	}

	settings, err := h.storageSettings(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load storage settings")
		return
	}
	result, err := storage.NewUploader(settings).Upload(r.Context(), storage.UploadInput{
		Scope:       storage.ScopePublic,
		Reader:      io.MultiReader(bytes.NewReader(head), file),
		Filename:    header.Filename,
		ContentType: "font/" + format,
		Size:        header.Size,
		Subdir:      "fonts",
	})
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "UPLOAD_FAILED", err.Error())
		return
	}

	f := &models.CMSCustomFont{
		ID: uuid.NewString(), FontName: toSlug(displayName), DisplayName: displayName,
		FileURL: result.URL, FontFormat: format,
	}
	if err := h.repo.CreateCustomFont(r.Context(), f); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_NAME", "A font with this name already exists")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not save font")
		return
	}
	response.JSON(w, http.StatusCreated, f)
}

// DELETE /api/v1/admin/cms/fonts/{id}
func (h *CMSHandler) DeleteFont(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteCustomFont(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete font")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

func allowedFontMagic(format string, head []byte) bool {
	// .otf and .ttf both wrap an sfnt container and are validated the same
	// way: an .otf file isn't required to use "OTTO" (PostScript/CFF
	// outlines) — it's just as valid for it to use TrueType outlines
	// internally, which uses the same sfnt version tag as a .ttf file. The
	// file extension, not the internal magic, is what actually distinguishes
	// the two as far as this app cares.
	switch format {
	case "otf", "ttf":
		return len(head) >= 4 && (string(head[0:4]) == "OTTO" || string(head[0:4]) == "\x00\x01\x00\x00" || string(head[0:4]) == "true" || string(head[0:4]) == "ttcf")
	case "woff":
		return len(head) >= 4 && string(head[0:4]) == "wOFF"
	case "woff2":
		return len(head) >= 4 && string(head[0:4]) == "wOF2"
	default:
		return false
	}
}

// ── Fun Facts ─────────────────────────────────────────────────────

func (h *CMSHandler) ListFunFacts(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") != "false"
	items, err := h.repo.ListFunFacts(r.Context(), activeOnly)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list fun facts")
		return
	}
	if items == nil {
		items = []*models.CMSFunFact{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateFunFact(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Label     string `json:"label"`
		Value     string `json:"value"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Label == "" || req.Value == "" {
		response.ValidationErr(w, map[string]string{"label": "required", "value": "required"})
		return
	}
	f := &models.CMSFunFact{
		ID: uuid.NewString(), Label: req.Label, Value: req.Value,
		Icon: req.Icon, SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateFunFact(r.Context(), f); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create fun fact")
		return
	}
	response.JSON(w, http.StatusCreated, f)
}

func (h *CMSHandler) UpdateFunFact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Label     string `json:"label"`
		Value     string `json:"value"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	f := &models.CMSFunFact{
		ID: id, Label: req.Label, Value: req.Value,
		Icon: req.Icon, SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.UpdateFunFact(r.Context(), f); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update fun fact")
		return
	}
	response.JSON(w, http.StatusOK, f)
}

func (h *CMSHandler) DeleteFunFact(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteFunFact(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete fun fact")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── FAQs ──────────────────────────────────────────────────────────

// ── Newsletter ────────────────────────────────────────────────────

var newsletterEmailRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// SubscribeNewsletter handles the public footer sign-up.
// POST /api/v1/cms/newsletter/subscribe
func (h *CMSHandler) SubscribeNewsletter(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email   string `json:"email"`
		Source  string `json:"source"`
		Website string `json:"website"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if strings.TrimSpace(req.Website) != "" {
		response.JSONMsg(w, http.StatusCreated, "Subscribed")
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || len(email) > 254 || !newsletterEmailRe.MatchString(email) {
		response.ValidationErr(w, map[string]string{"email": "a valid email address is required"})
		return
	}
	source := sanitizePlain(req.Source, 60)
	if source == "" {
		source = "website"
	}
	sub, err := h.repo.SubscribeNewsletter(r.Context(), uuid.NewString(), email, source)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not subscribe")
		return
	}
	response.JSON(w, http.StatusCreated, sub)
}

// ListNewsletterSubscribers returns a paginated list for the admin dashboard.
// GET /api/v1/admin/cms/newsletter/subscribers
func (h *CMSHandler) ListNewsletterSubscribers(w http.ResponseWriter, r *http.Request) {
	limit, offset := paginate(r)
	items, total, err := h.repo.ListNewsletterSubscribers(r.Context(), limit, offset)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list subscribers")
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{"total": total, "items": items})
}

// ExportNewsletterSubscribers streams every subscriber as CSV (also used for
// the "Export Excel" button, which opens the CSV in Excel).
// GET /api/v1/admin/cms/newsletter/subscribers/export?format=csv|excel
func (h *CMSHandler) ExportNewsletterSubscribers(w http.ResponseWriter, r *http.Request) {
	items, _, err := h.repo.ListNewsletterSubscribers(r.Context(), 0, 0)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not export subscribers")
		return
	}
	var buf bytes.Buffer
	cw := csv.NewWriter(&buf)
	_ = cw.Write([]string{"Email", "Source", "Subscribed At", "Status"})
	for _, s := range items {
		status := "subscribed"
		if !s.IsActive {
			status = "unsubscribed"
		}
		_ = cw.Write([]string{s.Email, s.Source, s.CreatedAt.Format(time.RFC3339), status})
	}
	cw.Flush()

	contentType := "text/csv; charset=utf-8"
	filename := "newsletter-subscribers.csv"
	if strings.EqualFold(r.URL.Query().Get("format"), "excel") {
		contentType = "application/vnd.ms-excel"
		filename = "newsletter-subscribers.xls"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func (h *CMSHandler) ListFAQs(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	items, err := h.repo.ListFAQs(r.Context(), category)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list FAQs")
		return
	}
	if items == nil {
		items = []*models.CMSFAQ{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateFAQ(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		Category  string `json:"category"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Question == "" {
		response.ValidationErr(w, map[string]string{"question": "required"})
		return
	}
	if req.Category == "" {
		req.Category = "general"
	}
	f := &models.CMSFAQ{
		ID: uuid.NewString(), Question: req.Question, Answer: req.Answer,
		Category: req.Category, SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateFAQ(r.Context(), f); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create FAQ")
		return
	}
	response.JSON(w, http.StatusCreated, f)
}

func (h *CMSHandler) UpdateFAQ(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		Category  string `json:"category"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	f := &models.CMSFAQ{
		ID: id, Question: req.Question, Answer: req.Answer,
		Category: req.Category, SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.UpdateFAQ(r.Context(), f); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update FAQ")
		return
	}
	response.JSON(w, http.StatusOK, f)
}

func (h *CMSHandler) DeleteFAQ(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteFAQ(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete FAQ")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── Resources ─────────────────────────────────────────────────────

func (h *CMSHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	items, err := h.repo.ListResources(r.Context(), category)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list resources")
		return
	}
	if items == nil {
		items = []*models.CMSResource{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Category    string `json:"category"`
		FileURL     string `json:"file_url"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Title == "" {
		response.ValidationErr(w, map[string]string{"title": "required"})
		return
	}
	if req.Category == "" {
		req.Category = "general"
	}
	res := &models.CMSResource{
		ID: uuid.NewString(), Title: req.Title, Category: req.Category,
		FileURL: req.FileURL, Description: req.Description,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateResource(r.Context(), res); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create resource")
		return
	}
	response.JSON(w, http.StatusCreated, res)
}

func (h *CMSHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.repo.GetResourceByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Resource not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get resource")
		return
	}
	var req struct {
		Title       string `json:"title"`
		Category    string `json:"category"`
		FileURL     string `json:"file_url"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Category != "" {
		existing.Category = req.Category
	}
	if req.FileURL != "" {
		existing.FileURL = req.FileURL
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	existing.SortOrder = req.SortOrder
	existing.IsActive = req.IsActive
	if err := h.repo.UpdateResource(r.Context(), existing); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update resource")
		return
	}
	response.JSON(w, http.StatusOK, existing)
}

func (h *CMSHandler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteResource(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete resource")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── Documents (Sports Rules, Press Releases, Reports, Speeches) ────

func validDocType(t string) bool {
	switch t {
	case "sports_rule", "press_release", "report", "speech":
		return true
	default:
		return false
	}
}

// GET /api/v1/cms/documents?doc_type=sports_rule&category=...
func (h *CMSHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	docType := r.URL.Query().Get("doc_type")
	if !validDocType(docType) {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid or missing doc_type")
		return
	}
	category := r.URL.Query().Get("category")
	activeOnly := r.URL.Query().Get("active") != "false"
	items, err := h.repo.ListDocuments(r.Context(), docType, category, activeOnly)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list documents")
		return
	}
	if items == nil {
		items = []*models.CMSDocument{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DocType     string `json:"doc_type"`
		Title       string `json:"title"`
		Category    string `json:"category"`
		FileURL     string `json:"file_url"`
		VideoURL    string `json:"video_url"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if !validDocType(req.DocType) {
		response.ValidationErr(w, map[string]string{"doc_type": "invalid"})
		return
	}
	if req.Title == "" {
		response.ValidationErr(w, map[string]string{"title": "required"})
		return
	}
	if req.Category == "" {
		req.Category = "general"
	}
	d := &models.CMSDocument{
		ID: uuid.NewString(), DocType: req.DocType, Title: req.Title, Category: req.Category,
		FileURL: req.FileURL, VideoURL: req.VideoURL, Description: req.Description,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateDocument(r.Context(), d); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create document")
		return
	}
	response.JSON(w, http.StatusCreated, d)
}

func (h *CMSHandler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.repo.GetDocumentByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Document not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get document")
		return
	}
	var req struct {
		Title       string `json:"title"`
		Category    string `json:"category"`
		FileURL     string `json:"file_url"`
		VideoURL    string `json:"video_url"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Category != "" {
		existing.Category = req.Category
	}
	// Overwritten unconditionally (unlike Title/Category): press releases toggle
	// between a PDF upload and a YouTube URL, so clearing one when setting the
	// other must actually persist.
	existing.FileURL = req.FileURL
	existing.VideoURL = req.VideoURL
	if req.Description != "" {
		existing.Description = req.Description
	}
	existing.SortOrder = req.SortOrder
	existing.IsActive = req.IsActive
	if err := h.repo.UpdateDocument(r.Context(), existing); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update document")
		return
	}
	response.JSON(w, http.StatusOK, existing)
}

func (h *CMSHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteDocument(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete document")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── Facilities ────────────────────────────────────────────────────

func (h *CMSHandler) ListFacilities(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") != "false"
	items, err := h.repo.ListFacilities(r.Context(), activeOnly)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list facilities")
		return
	}
	if items == nil {
		items = []*models.CMSFacility{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateFacility(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Category    string `json:"category"`
		ImageURL    string `json:"image_url"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Name == "" {
		response.ValidationErr(w, map[string]string{"name": "required"})
		return
	}
	slug := req.Slug
	if slug == "" {
		slug = toSlug(req.Name)
	}
	f := &models.CMSFacility{
		ID: uuid.NewString(), Name: req.Name, Slug: slug,
		Description: req.Description, Category: req.Category, ImageURL: req.ImageURL,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateFacility(r.Context(), f); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_SLUG", "A facility with this slug already exists")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create facility")
		return
	}
	response.JSON(w, http.StatusCreated, f)
}

func (h *CMSHandler) UpdateFacility(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.repo.GetFacilityByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Facility not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get facility")
		return
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		ImageURL    string `json:"image_url"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Category != "" {
		existing.Category = req.Category
	}
	if req.ImageURL != "" {
		existing.ImageURL = req.ImageURL
	}
	existing.SortOrder = req.SortOrder
	existing.IsActive = req.IsActive
	if err := h.repo.UpdateFacility(r.Context(), existing); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update facility")
		return
	}
	response.JSON(w, http.StatusOK, existing)
}

func (h *CMSHandler) DeleteFacility(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteFacility(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete facility")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── Associations ──────────────────────────────────────────────────

func (h *CMSHandler) ListAssociations(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") != "false"
	items, err := h.repo.ListAssociations(r.Context(), activeOnly)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list associations")
		return
	}
	if items == nil {
		items = []*models.CMSAssociation{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateAssociation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name         string `json:"name"`
		Slug         string `json:"slug"`
		Abbreviation string `json:"abbreviation"`
		Description  string `json:"description"`
		LogoURL      string `json:"logo_url"`
		WebsiteURL   string `json:"website_url"`
		Category     string `json:"category"`
		President    string `json:"president"`
		Secretary    string `json:"secretary"`
		Address      string `json:"address"`
		Phone        string `json:"phone"`
		SortOrder    int    `json:"sort_order"`
		IsActive     bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Name == "" {
		response.ValidationErr(w, map[string]string{"name": "required"})
		return
	}
	slug := req.Slug
	if slug == "" {
		slug = toSlug(req.Name)
	}
	a := &models.CMSAssociation{
		ID: uuid.NewString(), Name: req.Name, Slug: slug, Abbreviation: req.Abbreviation,
		Description: req.Description, LogoURL: req.LogoURL, WebsiteURL: req.WebsiteURL,
		Category: req.Category, President: req.President, Secretary: req.Secretary, Address: req.Address, Phone: req.Phone,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateAssociation(r.Context(), a); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			response.Err(w, http.StatusConflict, "DUPLICATE_SLUG", "An association with this slug already exists")
			return
		}
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create association")
		return
	}
	response.JSON(w, http.StatusCreated, a)
}

func (h *CMSHandler) UpdateAssociation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.repo.GetAssociationByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Association not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get association")
		return
	}
	var req struct {
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
		Description  string `json:"description"`
		LogoURL      string `json:"logo_url"`
		WebsiteURL   string `json:"website_url"`
		Category     string `json:"category"`
		President    string `json:"president"`
		Secretary    string `json:"secretary"`
		Address      string `json:"address"`
		Phone        string `json:"phone"`
		SortOrder    int    `json:"sort_order"`
		IsActive     bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.LogoURL != "" {
		existing.LogoURL = req.LogoURL
	}
	if req.WebsiteURL != "" {
		existing.WebsiteURL = req.WebsiteURL
	}
	if req.Abbreviation != "" {
		existing.Abbreviation = req.Abbreviation
	}
	existing.Category = req.Category
	existing.President = req.President
	existing.Secretary = req.Secretary
	existing.Address = req.Address
	existing.Phone = req.Phone
	existing.SortOrder = req.SortOrder
	existing.IsActive = req.IsActive
	if err := h.repo.UpdateAssociation(r.Context(), existing); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update association")
		return
	}
	response.JSON(w, http.StatusOK, existing)
}

func (h *CMSHandler) DeleteAssociation(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteAssociation(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete association")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── Invest with Us ────────────────────────────────────────────────

func (h *CMSHandler) ListInvest(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.ListInvest(r.Context())
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list invest items")
		return
	}
	if items == nil {
		items = []*models.CMSInvest{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateInvest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title     string `json:"title"`
		Subtitle  string `json:"subtitle"`
		Content   string `json:"content"`
		ImageURL  string `json:"image_url"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Title == "" {
		response.ValidationErr(w, map[string]string{"title": "required"})
		return
	}
	inv := &models.CMSInvest{
		ID: uuid.NewString(), Title: req.Title, Subtitle: req.Subtitle,
		Content: req.Content, ImageURL: req.ImageURL,
		SortOrder: req.SortOrder, IsActive: req.IsActive,
	}
	if err := h.repo.CreateInvest(r.Context(), inv); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create invest item")
		return
	}
	response.JSON(w, http.StatusCreated, inv)
}

func (h *CMSHandler) UpdateInvest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.repo.GetInvestByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Invest item not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get invest item")
		return
	}
	var req struct {
		Title     string `json:"title"`
		Subtitle  string `json:"subtitle"`
		Content   string `json:"content"`
		ImageURL  string `json:"image_url"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Subtitle != "" {
		existing.Subtitle = req.Subtitle
	}
	if req.Content != "" {
		existing.Content = req.Content
	}
	if req.ImageURL != "" {
		existing.ImageURL = req.ImageURL
	}
	existing.SortOrder = req.SortOrder
	existing.IsActive = req.IsActive
	if err := h.repo.UpdateInvest(r.Context(), existing); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update invest item")
		return
	}
	response.JSON(w, http.StatusOK, existing)
}

func (h *CMSHandler) DeleteInvest(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteInvest(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete invest item")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}

// ── Team Members & Governing Council ────────────────────────────────
// Both share the cms_team_members table; member_group ('team' | 'council')
// discriminates which directory a profile belongs to, same pattern as
// cms_documents.doc_type.

func validMemberGroup(g string) bool {
	switch g {
	case "team", "council":
		return true
	default:
		return false
	}
}

func (h *CMSHandler) ListTeam(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	if group == "" {
		group = "team"
	} else if !validMemberGroup(group) {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid group")
		return
	}
	activeOnly := r.URL.Query().Get("active") != "false"
	items, err := h.repo.ListTeamMembers(r.Context(), group, activeOnly)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not list team members")
		return
	}
	if items == nil {
		items = []*models.CMSTeamMember{}
	}
	response.JSON(w, http.StatusOK, items)
}

func (h *CMSHandler) CreateTeamMember(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FullName     string  `json:"full_name"`
		Designation  string  `json:"designation"`
		ImageURL     string  `json:"image_url"`
		Bio          string  `json:"bio"`
		SortOrder    int     `json:"sort_order"`
		IsActive     bool    `json:"is_active"`
		MemberGroup  string  `json:"member_group"`
		DepartmentID *string `json:"department_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.FullName == "" {
		response.ValidationErr(w, map[string]string{"full_name": "required"})
		return
	}
	if req.MemberGroup == "" {
		req.MemberGroup = "team"
	} else if !validMemberGroup(req.MemberGroup) {
		response.ValidationErr(w, map[string]string{"member_group": "invalid"})
		return
	}
	m := &models.CMSTeamMember{
		ID: uuid.NewString(), FullName: req.FullName, Designation: req.Designation,
		ImageURL: req.ImageURL, Bio: req.Bio, SortOrder: req.SortOrder, IsActive: req.IsActive,
		MemberGroup: req.MemberGroup, DepartmentID: nullableID(req.DepartmentID),
	}
	if err := h.repo.CreateTeamMember(r.Context(), m); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create team member")
		return
	}
	response.JSON(w, http.StatusCreated, m)
}

func (h *CMSHandler) UpdateTeamMember(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.repo.GetTeamMemberByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Team member not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not get team member")
		return
	}
	var req struct {
		FullName     string  `json:"full_name"`
		Designation  string  `json:"designation"`
		ImageURL     string  `json:"image_url"`
		Bio          string  `json:"bio"`
		SortOrder    int     `json:"sort_order"`
		IsActive     bool    `json:"is_active"`
		DepartmentID *string `json:"department_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON")
		return
	}
	if req.DepartmentID != nil {
		existing.DepartmentID = nullableID(req.DepartmentID)
	}
	if req.FullName != "" {
		existing.FullName = req.FullName
	}
	existing.Designation = req.Designation
	existing.ImageURL = req.ImageURL
	existing.Bio = req.Bio
	existing.SortOrder = req.SortOrder
	existing.IsActive = req.IsActive
	if err := h.repo.UpdateTeamMember(r.Context(), existing); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update team member")
		return
	}
	response.JSON(w, http.StatusOK, existing)
}

func (h *CMSHandler) DeleteTeamMember(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteTeamMember(r.Context(), chi.URLParam(r, "id")); err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not delete team member")
		return
	}
	response.JSONMsg(w, http.StatusOK, "Deleted")
}
