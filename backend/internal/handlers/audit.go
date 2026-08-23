package handlers

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/middleware"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuditHandler struct {
	repo   *repository.AuditRepo
	writer *middleware.AuditWriter
}

func (h *AuditHandler) SetWriter(writer *middleware.AuditWriter) {
	h.writer = writer
}

// GET /api/v1/admin/audit-logs
func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))
	p := &models.PaginationParams{Page: page, PerPage: perPage, Search: q.Get("search")}

	logs, total, err := h.repo.List(r.Context(), p)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch audit logs")
		return
	}
	response.JSONPaged(w, http.StatusOK, logs, &response.Meta{Page: p.Page, PerPage: p.PerPage, Total: total})
}

// GET /api/v1/admin/audit-logs/{id}
func (h *AuditHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log, err := h.repo.GetByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Audit log entry not found")
		return
	}
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch audit log")
		return
	}
	response.JSON(w, http.StatusOK, log)
}

type frontendActivityRequest struct {
	Type       string         `json:"type"`
	Action     string         `json:"action"`
	FromPath   string         `json:"from_path"`
	ToPath     string         `json:"to_path"`
	Path       string         `json:"path"`
	PageName   string         `json:"page_name"`
	FromPage   string         `json:"from_page"`
	Label      string         `json:"label"`
	RouteName  string         `json:"route_name"`
	Section    string         `json:"section"`
	Resource   string         `json:"resource"`
	OccurredAt *time.Time     `json:"occurred_at"`
	Metadata   map[string]any `json:"metadata"`
}

type normalizedFrontendActivity struct {
	Action         string
	EventType      string
	Resource       string
	ResourceID     string
	Endpoint       string
	PageName       string
	PayloadExcerpt []byte
}

// POST /api/v1/account/activity-events
func (h *AuditHandler) RecordFrontendActivity(w http.ResponseWriter, r *http.Request) {
	if h.writer == nil {
		response.Err(w, http.StatusServiceUnavailable, "AUDIT_UNAVAILABLE", "Audit writer is not ready")
		return
	}
	var input frontendActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	normalized, err := normalizeFrontendActivity(input)
	if err != nil {
		response.Err(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	userEmail, _ := r.Context().Value(models.CtxUserEmail).(string)
	uid := userID
	ua := r.UserAgent()
	browser, osName := middleware.ParseBrowserOS(ua)
	deviceType := middleware.ParseDeviceType(ua)
	deviceInfo := middleware.FormatDeviceInfo(osName, deviceType, browser, "Browser")

	entry := &models.AuditLog{
		ID:              uuid.NewString(),
		UserID:          &uid,
		Username:        strings.ToLower(strings.TrimSpace(userEmail)),
		SessionID:       strings.TrimSpace(r.Header.Get("X-Request-ID")),
		Action:          normalized.Action,
		Resource:        normalized.Resource,
		ResourceID:      normalized.ResourceID,
		Method:          "UI",
		Endpoint:        normalized.Endpoint,
		IPAddress:       requestIP(r),
		ForwardedIP:     r.Header.Get("X-Forwarded-For"),
		UserAgent:       ua,
		Browser:         browser,
		OSName:          osName,
		ClientType:      "Browser",
		DeviceInfo:      deviceInfo,
		Platform:        osName,
		ResponseCode:    http.StatusOK,
		ResponseTimeMs:  0,
		EventType:       normalized.EventType,
		EventStatus:     "SUCCESS",
		SeverityLevel:   "INFO",
		Authenticated:   true,
		PayloadExcerpt:  normalized.PayloadExcerpt,
		ThreatScore:     0,
		AnomalyDetected: false,
	}
	h.writer.Enqueue(entry, userID)
	response.JSON(w, http.StatusAccepted, map[string]any{"recorded": true, "action": normalized.Action})
}

func normalizeFrontendActivity(input frontendActivityRequest) (*normalizedFrontendActivity, error) {
	kind := strings.ToLower(cleanAuditText(input.Type, 40))
	if kind == "" {
		kind = strings.ToLower(cleanAuditText(input.Action, 40))
	}
	if kind == "" {
		kind = "interaction"
	}

	page := cleanAuditText(input.PageName, 120)
	if page == "" {
		page = friendlyPathName(input.ToPath)
	}
	path := cleanAuditPath(auditFirstNonEmpty(input.ToPath, input.Path))
	fromPath := cleanAuditPath(input.FromPath)
	label := cleanAuditText(input.Label, 120)
	section := cleanAuditText(input.Section, 80)
	resource := cleanAuditText(input.Resource, 80)
	if resource == "" || resource == "Frontend Activity" {
		if section != "" {
			resource = strings.Title(strings.ReplaceAll(section, "-", " "))
		} else if page != "" {
			resource = page
		} else {
			resource = "Portal Activity"
		}
	}

	resourceID := cleanAuditText(auditFirstNonEmpty(label, section, input.RouteName, page), 120)
	if resourceID == "" {
		resourceID = kind
	}

	metadata := map[string]any{
		"type":       kind,
		"page":       page,
		"from_page":  cleanAuditText(input.FromPage, 120),
		"from_path":  fromPath,
		"to_path":    path,
		"route_name": cleanAuditText(input.RouteName, 80),
		"section":    section,
		"label":      label,
	}
	for key, value := range input.Metadata {
		key = cleanAuditText(key, 50)
		if key == "" || strings.Contains(strings.ToLower(key), "password") || strings.Contains(strings.ToLower(key), "token") {
			continue
		}
		metadata[key] = value
	}
	payload, _ := json.Marshal(metadata)

	eventType := "UI_" + strings.ToUpper(kind)
	action := cleanAuditText(input.Action, 120)

	// If action is generic, generate descriptive action
	if action == "" || action == "navigate" || action == "navigation" || action == "click" || action == "interaction" || action == "ui:navigate" || action == "ui:click" {
		switch kind {
		case "navigation", "navigate", "route", "view", "page_view":
			eventType = "UI_NAVIGATION"
			if page != "" {
				action = "Visited " + page
			} else if section != "" {
				action = "Visited " + strings.Title(strings.ReplaceAll(section, "-", " "))
			} else {
				action = "Visited Portal Page"
			}
		case "update_profile", "profile":
			eventType = "PROFILE_UPDATE"
			action = "Updated Profile Information"
		case "upload_avatar", "avatar":
			eventType = "PROFILE_AVATAR"
			action = "Uploaded Profile Photo"
		case "change_password", "password":
			eventType = "SECURITY_PASSWORD"
			action = "Changed Account Password"
		case "security_2fa", "2fa":
			eventType = "SECURITY_2FA"
			action = "Configured Two-Factor Authentication"
		case "update_preferences", "preferences":
			eventType = "PREFERENCES_UPDATE"
			action = "Saved Account Preferences"
		case "download_document", "download_file":
			eventType = "DOCUMENT_DOWNLOAD"
			if label != "" {
				action = "Downloaded " + label
			} else {
				action = "Downloaded Document"
			}
		case "save_draft":
			eventType = "APPLICATION_DRAFT"
			action = "Saved Application Draft"
		case "submit_application":
			eventType = "APPLICATION_SUBMIT"
			action = "Submitted Application Form"
		case "payment", "momo_payment":
			eventType = "PAYMENT_INITIATE"
			action = "Initiated Payment Request"
		default: // clicks & generic interactions
			eventType = "UI_INTERACTION"
			if label != "" {
				lowLabel := strings.ToLower(label)
				lowSec := strings.ToLower(section)
				lowPage := strings.ToLower(page)
				if strings.Contains(lowLabel, "save") && (strings.Contains(lowSec, "profile") || strings.Contains(lowPage, "profile")) {
					action = "Updated Profile Information"
				} else if strings.Contains(lowLabel, "password") || strings.Contains(lowLabel, "security") {
					action = "Updated Security Settings"
				} else if strings.Contains(lowLabel, "save") && (strings.Contains(lowSec, "settings") || strings.Contains(lowPage, "settings")) {
					action = "Updated Account Settings"
				} else if strings.Contains(lowLabel, "submit") && (strings.Contains(lowLabel, "application") || strings.Contains(lowLabel, "form")) {
					action = "Submitted Application Form"
				} else if strings.Contains(lowLabel, "draft") {
					action = "Saved Application Draft"
				} else if strings.Contains(lowLabel, "pay") || strings.Contains(lowLabel, "momo") || strings.Contains(lowLabel, "mobile money") {
					action = "Initiated Mobile Money Payment"
				} else if strings.Contains(lowLabel, "upload") {
					action = "Uploaded Document / File"
				} else if strings.Contains(lowLabel, "logout") || strings.Contains(lowLabel, "sign out") {
					action = "Logged Out of Portal"
				} else {
					action = "Clicked: " + label
				}
			} else if page != "" {
				action = "Interacted with " + page
			} else {
				action = "Portal Interaction"
			}
		}
	}

	return &normalizedFrontendActivity{
		Action:         action,
		EventType:      eventType,
		Resource:       resource,
		ResourceID:     resourceID,
		Endpoint:       path,
		PageName:       page,
		PayloadExcerpt: payload,
	}, nil
}

func cleanAuditText(value string, maxLen int) string {
	value = strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
	if maxLen > 0 && len(value) > maxLen {
		return value[:maxLen]
	}
	return value
}

func cleanAuditPath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		if idx := strings.Index(value[8:], "/"); idx >= 0 {
			value = value[idx+8:]
		}
	}
	if idx := strings.Index(value, "#"); idx >= 0 {
		value = value[:idx]
	}
	if len(value) > 240 {
		value = value[:240]
	}
	return value
}

func friendlyPathName(path string) string {
	path = cleanAuditPath(path)
	if path == "" || path == "/" {
		return "Portal"
	}
	if idx := strings.Index(path, "section="); idx >= 0 {
		section := path[idx+len("section="):]
		if end := strings.IndexAny(section, "&#"); end >= 0 {
			section = section[:end]
		}
		return strings.Title(strings.ReplaceAll(section, "-", " "))
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return "Portal"
	}
	return strings.Title(strings.ReplaceAll(parts[len(parts)-1], "-", " "))
}

func auditFirstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func requestIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	if real := strings.TrimSpace(r.Header.Get("X-Real-IP")); real != "" {
		return real
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
