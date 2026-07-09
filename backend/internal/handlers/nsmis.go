package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
)

type NSMISHandler struct {
	repo *repository.NSMISRepo
	cfg  *config.Config
}

func (h *NSMISHandler) allowed(r *http.Request, permissions ...string) bool {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	ok, err := h.repo.HasAnyPermission(r.Context(), userID, permissions...)
	return err == nil && ok
}
func (h *NSMISHandler) wide(r *http.Request, permission string) bool { return h.allowed(r, permission) }
func decodeStrict(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Request body is invalid or contains unknown fields")
		return false
	}
	return true
}

func (h *NSMISHandler) ListFederations(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "federations:read:any", "federations:read:own") {
		response.Err(w, http.StatusForbidden, "FORBIDDEN", "Federation access is not assigned to this account")
		return
	}
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	items, err := h.repo.ListFederations(r.Context(), userID, h.wide(r, "federations:read:any"))
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load federations")
		return
	}
	response.JSON(w, 200, items)
}

func (h *NSMISHandler) CreateFederation(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "federations:write:any") {
		response.Err(w, 403, "FORBIDDEN", "Only an NCS Administrator can create a federation")
		return
	}
	var f models.Federation
	if !decodeStrict(w, r, &f) {
		return
	}
	f.Name = strings.TrimSpace(f.Name)
	f.Acronym = strings.ToUpper(strings.TrimSpace(f.Acronym))
	f.NCSRegistrationNumber = strings.TrimSpace(f.NCSRegistrationNumber)
	if f.RecognitionStatus == "" {
		f.RecognitionStatus = "PENDING"
	}
	f.RecognitionStatus = strings.ToUpper(f.RecognitionStatus)
	if f.Name == "" || f.Acronym == "" || f.NCSRegistrationNumber == "" {
		response.ValidationErr(w, map[string]string{"name": "Name, acronym and NCS registration number are required"})
		return
	}
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	err := h.repo.CreateFederation(r.Context(), &f, actorID)
	if errors.Is(err, repository.ErrDuplicate) {
		response.Err(w, 409, "DUPLICATE", "A federation with this acronym or registration number already exists")
		return
	}
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not create federation")
		return
	}
	response.JSON(w, 201, &f)
}
func (h *NSMISHandler) UpdateFederation(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "federations:write:any", "federations:write:own") {
		response.Err(w, 403, "FORBIDDEN", "Federation editing is not assigned to this account")
		return
	}
	var f models.Federation
	if !decodeStrict(w, r, &f) {
		return
	}
	f.ID = chi.URLParam(r, "federationID")
	f.Name = strings.TrimSpace(f.Name)
	f.Acronym = strings.ToUpper(strings.TrimSpace(f.Acronym))
	f.RecognitionStatus = strings.ToUpper(strings.TrimSpace(f.RecognitionStatus))
	if f.Version < 1 || f.Name == "" || f.Acronym == "" || f.NCSRegistrationNumber == "" {
		response.ValidationErr(w, map[string]string{"federation": "Name, acronym, registration number and current version are required"})
		return
	}
	uid, _ := r.Context().Value(models.CtxUserID).(string)
	err := h.repo.UpdateFederation(r.Context(), &f, uid, h.wide(r, "federations:write:any"))
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, 409, "VERSION_CONFLICT", "Federation is outside your scope or was changed by another user")
		return
	}
	if errors.Is(err, repository.ErrDuplicate) {
		response.Err(w, 409, "DUPLICATE", "A federation with this acronym or registration number exists")
		return
	}
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Federation could not be updated")
		return
	}
	response.JSON(w, 200, &f)
}

func (h *NSMISHandler) ListPeriods(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "reports:read:any", "reports:read:own") {
		response.Err(w, http.StatusForbidden, "FORBIDDEN", "Reporting access is not assigned to this account")
		return
	}
	items, err := h.repo.ListPeriods(r.Context(), r.URL.Query().Get("open") != "false")
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load reporting periods")
		return
	}
	response.JSON(w, 200, items)
}

func (h *NSMISHandler) CreatePeriod(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "federations:write:any") {
		response.Err(w, 403, "FORBIDDEN", "Only an NCS Administrator can create reporting periods")
		return
	}
	var req struct {
		PeriodType  string `json:"period_type"`
		Name        string `json:"name"`
		StartsOn    string `json:"starts_on"`
		EndsOn      string `json:"ends_on"`
		DueOn       string `json:"due_on"`
		GraceEndsOn string `json:"grace_ends_on"`
		Timezone    string `json:"timezone"`
		IsOpen      *bool  `json:"is_open"`
	}
	if !decodeStrict(w, r, &req) {
		return
	}
	parse := func(v string) (time.Time, error) { return time.Parse("2006-01-02", v) }
	starts, e1 := parse(req.StartsOn)
	ends, e2 := parse(req.EndsOn)
	due, e3 := parse(req.DueOn)
	if req.Name == "" || e1 != nil || e2 != nil || e3 != nil || ends.Before(starts) || due.Before(ends) {
		response.ValidationErr(w, map[string]string{"period": "Name and valid YYYY-MM-DD start/end/due dates are required"})
		return
	}
	p := models.ReportingPeriod{PeriodType: strings.ToUpper(req.PeriodType), Name: strings.TrimSpace(req.Name), StartsOn: starts, EndsOn: ends, DueOn: due, Timezone: req.Timezone, IsOpen: true}
	if p.Timezone == "" {
		p.Timezone = "Africa/Kampala"
	}
	if req.IsOpen != nil {
		p.IsOpen = *req.IsOpen
	}
	if req.GraceEndsOn != "" {
		v, err := parse(req.GraceEndsOn)
		if err != nil {
			response.ValidationErr(w, map[string]string{"grace_ends_on": "Use YYYY-MM-DD"})
			return
		}
		p.GraceEndsOn = &v
	}
	err := h.repo.CreatePeriod(r.Context(), &p)
	if errors.Is(err, repository.ErrDuplicate) {
		response.Err(w, 409, "DUPLICATE", "This reporting period already exists")
		return
	}
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not create reporting period")
		return
	}
	response.JSON(w, 201, &p)
}

func (h *NSMISHandler) GenerateObligations(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "reports:review:any") {
		response.Err(w, 403, "FORBIDDEN", "Only an NCS Administrator can generate obligations")
		return
	}
	periodID := strings.TrimSpace(chi.URLParam(r, "periodID"))
	var req struct {
		ReportType string `json:"report_type"`
	}
	if !decodeStrict(w, r, &req) {
		return
	}
	req.ReportType = strings.ToUpper(req.ReportType)
	valid := map[string]bool{"GOVERNANCE": true, "ATHLETE": true, "COMPETITION": true, "MEDAL": true, "WORKFORCE": true, "TALENT": true, "SAFEGUARDING": true, "FINANCE": true, "EQUIPMENT": true}
	if !valid[req.ReportType] {
		response.ValidationErr(w, map[string]string{"report_type": "Unsupported report type"})
		return
	}
	n, err := h.repo.GenerateObligations(r.Context(), periodID, req.ReportType)
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not generate obligations")
		return
	}
	response.JSON(w, 201, map[string]int64{"created": n})
}

func (h *NSMISHandler) ListObligations(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "reports:read:any", "reports:read:own") {
		response.Err(w, http.StatusForbidden, "FORBIDDEN", "Reporting access is not assigned to this account")
		return
	}
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	items, err := h.repo.ListObligations(r.Context(), userID, r.URL.Query().Get("period_id"), h.wide(r, "reports:read:any"))
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load report obligations")
		return
	}
	response.JSON(w, 200, items)
}

func (h *NSMISHandler) GovernanceDashboard(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "dashboard:governance:read") {
		response.Err(w, http.StatusForbidden, "FORBIDDEN", "This role cannot view governance analytics")
		return
	}
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	data, err := h.repo.GovernanceDashboard(r.Context(), userID, r.URL.Query().Get("period_id"), h.wide(r, "federations:read:any"))
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load governance dashboard")
		return
	}
	response.JSON(w, 200, data)
}

func parseOptionalDate(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	v, err := time.Parse("2006-01-02", value)
	return &v, err
}

func (h *NSMISHandler) SaveGovernanceDraft(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "reports:write:own", "reports:review:any") {
		response.Err(w, 403, "FORBIDDEN", "This role cannot edit governance reports")
		return
	}
	var req struct {
		Version                  int    `json:"version"`
		ExecutiveMeetingHeld     bool   `json:"executive_meeting_held"`
		ExecutiveMeetingDate     string `json:"executive_meeting_date"`
		AGMConducted             bool   `json:"agm_conducted"`
		AGMDate                  string `json:"agm_date"`
		BoardMeetingHeld         bool   `json:"board_meeting_held"`
		BoardMeetingDate         string `json:"board_meeting_date"`
		ElectionsConducted       bool   `json:"elections_conducted"`
		ElectionDate             string `json:"election_date"`
		DisciplinaryCasesHandled int    `json:"disciplinary_cases_handled"`
	}
	if !decodeStrict(w, r, &req) {
		return
	}
	if req.DisciplinaryCasesHandled < 0 {
		response.ValidationErr(w, map[string]string{"disciplinary_cases_handled": "Cannot be negative"})
		return
	}
	execDate, e1 := parseOptionalDate(req.ExecutiveMeetingDate)
	agmDate, e2 := parseOptionalDate(req.AGMDate)
	boardDate, e3 := parseOptionalDate(req.BoardMeetingDate)
	electionDate, e4 := parseOptionalDate(req.ElectionDate)
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || (req.ExecutiveMeetingHeld && execDate == nil) || (req.AGMConducted && agmDate == nil) || (req.BoardMeetingHeld && boardDate == nil) || (req.ElectionsConducted && electionDate == nil) {
		response.ValidationErr(w, map[string]string{"dates": "Use YYYY-MM-DD and provide a date for every Yes response"})
		return
	}
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	id, version, err := h.repo.SaveGovernanceDraft(r.Context(), chi.URLParam(r, "obligationID"), actorID, h.wide(r, "reports:review:any"), req.Version, repository.GovernanceDraftInput{ExecutiveMeetingHeld: req.ExecutiveMeetingHeld, ExecutiveMeetingDate: execDate, AGMConducted: req.AGMConducted, AGMDate: agmDate, BoardMeetingHeld: req.BoardMeetingHeld, BoardMeetingDate: boardDate, ElectionsConducted: req.ElectionsConducted, ElectionDate: electionDate, DisciplinaryCasesHandled: req.DisciplinaryCasesHandled})
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, 404, "NOT_FOUND", "Governance obligation was not found in your assigned scope")
		return
	}
	if err != nil && strings.Contains(err.Error(), "version conflict") {
		response.Err(w, 409, "VERSION_CONFLICT", "This report was changed by another user; reload before saving")
		return
	}
	if err != nil {
		response.Err(w, 422, "REPORT_INVALID", err.Error())
		return
	}
	response.JSON(w, 200, map[string]interface{}{"id": id, "status": "DRAFT", "version": version})
}

func (h *NSMISHandler) TransitionReport(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ToStatus string `json:"to_status"`
		Reason   string `json:"reason"`
	}
	if !decodeStrict(w, r, &req) {
		return
	}
	req.ToStatus = strings.ToUpper(strings.TrimSpace(req.ToStatus))
	allowed := false
	wide := false
	switch req.ToStatus {
	case "SUBMITTED":
		allowed = h.allowed(r, "reports:write:own", "reports:review:any")
		wide = h.wide(r, "reports:review:any")
	case "PRESIDENT_APPROVED", "PRESIDENT_RETURNED":
		allowed = h.allowed(r, "reports:approve:own")
	case "NCS_UNDER_REVIEW", "NEEDS_CORRECTION", "APPROVED", "LOCKED":
		allowed = h.allowed(r, "reports:review:any")
		wide = allowed
	}
	if !allowed {
		response.Err(w, 403, "FORBIDDEN", "This role cannot perform the requested report transition")
		return
	}
	if (req.ToStatus == "PRESIDENT_RETURNED" || req.ToStatus == "NEEDS_CORRECTION") && strings.TrimSpace(req.Reason) == "" {
		response.ValidationErr(w, map[string]string{"reason": "A reason is required when returning a report"})
		return
	}
	actorID, _ := r.Context().Value(models.CtxUserID).(string)
	err := h.repo.TransitionReport(r.Context(), chi.URLParam(r, "reportID"), actorID, req.ToStatus, strings.TrimSpace(req.Reason), wide)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, 404, "NOT_FOUND", "Report was not found in your assigned scope")
		return
	}
	if err != nil {
		response.Err(w, 409, "INVALID_TRANSITION", err.Error())
		return
	}
	response.JSONMsg(w, 200, "Report status updated")
}

func (h *NSMISHandler) AthleteDashboard(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "dashboard:athletes:read") {
		response.Err(w, 403, "FORBIDDEN", "This role cannot view athlete analytics")
		return
	}
	data, err := h.repo.AthleteDashboard(r.Context())
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load athlete dashboard")
		return
	}
	response.JSON(w, 200, data)
}
func (h *NSMISHandler) PerformanceDashboard(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "dashboard:performance:read") {
		response.Err(w, 403, "FORBIDDEN", "This role cannot view performance analytics")
		return
	}
	data, err := h.repo.PerformanceDashboard(r.Context())
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load performance dashboard")
		return
	}
	response.JSON(w, 200, data)
}
func (h *NSMISHandler) FinanceDashboard(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "dashboard:finance:read") {
		response.Err(w, 403, "FORBIDDEN", "This role cannot view financial analytics")
		return
	}
	data, err := h.repo.FinanceDashboard(r.Context())
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load finance dashboard")
		return
	}
	response.JSON(w, 200, data)
}
func (h *NSMISHandler) TalentDashboard(w http.ResponseWriter, r *http.Request) {
	if !h.allowed(r, "dashboard:talent:read") {
		response.Err(w, 403, "FORBIDDEN", "This role cannot view talent analytics")
		return
	}
	data, err := h.repo.TalentDashboard(r.Context())
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load talent dashboard")
		return
	}
	response.JSON(w, 200, data)
}
