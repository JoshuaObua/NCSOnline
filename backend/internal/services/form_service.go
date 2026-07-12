package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrForbiddenDept    = errors.New("forbidden: outside your department")
	ErrTemplateClosed   = errors.New("form is not accepting submissions")
	ErrPaymentRequired  = errors.New("payment is required before submission")
	ErrDuplicatePending = errors.New("you already have a pending application for this form")
)

type FormService struct {
	forms *repository.FormRepo
	depts *repository.DepartmentRepo
	audit *repository.AuditRepo
}

func NewFormService(forms *repository.FormRepo, depts *repository.DepartmentRepo, audit *repository.AuditRepo) *FormService {
	return &FormService{forms: forms, depts: depts, audit: audit}
}

// ── Departments ──────────────────────────────────────────────────

func (s *FormService) ListDepartments(ctx context.Context) ([]*models.Department, error) {
	return s.depts.List(ctx, true)
}

func (s *FormService) UserDepartmentIDs(ctx context.Context, userID string) ([]string, error) {
	return s.depts.UserDepartmentIDs(ctx, userID)
}

// ── Template authoring (admin) ───────────────────────────────────

type SaveTemplateInput struct {
	DepartmentID   string
	Title          string
	Description    string
	Sections       json.RawMessage
	BannerImageURL string
	PriceUGX       float64
	Status         string
	Fields         []*models.FormField
}

var slugInvalid = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugInvalid.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "form"
	}
	if len(s) > 100 {
		s = s[:100]
	}
	return s
}

func validateStatus(st string) (string, error) {
	switch strings.ToUpper(st) {
	case "", models.FormStatusDraft:
		return models.FormStatusDraft, nil
	case models.FormStatusOpen, models.FormStatusClosed, models.FormStatusArchived:
		return strings.ToUpper(st), nil
	}
	return "", fmt.Errorf("invalid status: %s", st)
}

func normalizeSectionsJSON(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 || strings.TrimSpace(string(raw)) == "" || string(raw) == "null" {
		return json.RawMessage("[]"), nil
	}
	var sections []json.RawMessage
	if err := json.Unmarshal(raw, &sections); err != nil {
		return nil, fmt.Errorf("sections must be a JSON array")
	}
	return raw, nil
}

func (s *FormService) ensureCanManage(ctx context.Context, userID string, isSuperAdmin bool, templateDeptID string) error {
	if isSuperAdmin {
		return nil
	}
	userDepts, err := s.depts.UserDepartmentIDs(ctx, userID)
	if err != nil {
		return err
	}
	for _, d := range userDepts {
		if d == templateDeptID {
			return nil
		}
	}
	return ErrForbiddenDept
}

func (s *FormService) CreateTemplate(ctx context.Context, userID string, isSuperAdmin bool, in SaveTemplateInput) (*models.FormTemplate, error) {
	if in.DepartmentID == "" {
		return nil, errors.New("department_id is required")
	}
	if strings.TrimSpace(in.Title) == "" {
		return nil, errors.New("title is required")
	}
	if err := s.ensureCanManage(ctx, userID, isSuperAdmin, in.DepartmentID); err != nil {
		return nil, err
	}
	st, err := validateStatus(in.Status)
	if err != nil {
		return nil, err
	}
	sections, err := normalizeSectionsJSON(in.Sections)
	if err != nil {
		return nil, err
	}
	base := slugify(in.Title)
	slug := base + "-" + uuid.NewString()[:6]
	uid := userID
	t := &models.FormTemplate{
		DepartmentID:   in.DepartmentID,
		Slug:           slug,
		Title:          strings.TrimSpace(in.Title),
		Description:    in.Description,
		Sections:       sections,
		BannerImageURL: in.BannerImageURL,
		PriceUGX:       in.PriceUGX,
		Status:         st,
		CreatedBy:      &uid,
	}
	if err := s.forms.CreateTemplate(ctx, t); err != nil {
		return nil, err
	}
	if len(in.Fields) > 0 {
		if err := s.normalizeAndSaveFields(ctx, t.ID, in.Fields); err != nil {
			return nil, err
		}
	}
	return s.forms.GetTemplate(ctx, t.ID, true)
}

func (s *FormService) UpdateTemplate(ctx context.Context, id, userID string, isSuperAdmin bool, in SaveTemplateInput) (*models.FormTemplate, error) {
	existing, err := s.forms.GetTemplate(ctx, id, false)
	if err != nil {
		return nil, err
	}
	if err := s.ensureCanManage(ctx, userID, isSuperAdmin, existing.DepartmentID); err != nil {
		return nil, err
	}
	if in.DepartmentID != "" && in.DepartmentID != existing.DepartmentID {
		// reassigning department requires permission on the new one too
		if err := s.ensureCanManage(ctx, userID, isSuperAdmin, in.DepartmentID); err != nil {
			return nil, err
		}
		existing.DepartmentID = in.DepartmentID
	}
	if strings.TrimSpace(in.Title) != "" {
		existing.Title = strings.TrimSpace(in.Title)
	}
	existing.Description = in.Description
	if in.Sections != nil {
		sections, err := normalizeSectionsJSON(in.Sections)
		if err != nil {
			return nil, err
		}
		existing.Sections = sections
	} else if len(existing.Sections) == 0 {
		existing.Sections = json.RawMessage("[]")
	}
	existing.BannerImageURL = in.BannerImageURL
	existing.PriceUGX = in.PriceUGX
	if in.Status != "" {
		st, err := validateStatus(in.Status)
		if err != nil {
			return nil, err
		}
		existing.Status = st
	}
	if err := s.forms.UpdateTemplate(ctx, existing); err != nil {
		return nil, err
	}
	if in.Fields != nil {
		if err := s.normalizeAndSaveFields(ctx, existing.ID, in.Fields); err != nil {
			return nil, err
		}
	}
	return s.forms.GetTemplate(ctx, existing.ID, true)
}

func (s *FormService) DeleteTemplate(ctx context.Context, id, userID string, isSuperAdmin bool) error {
	existing, err := s.forms.GetTemplate(ctx, id, false)
	if err != nil {
		return err
	}
	if err := s.ensureCanManage(ctx, userID, isSuperAdmin, existing.DepartmentID); err != nil {
		return err
	}
	return s.forms.DeleteTemplate(ctx, id)
}

func (s *FormService) ListAdminTemplates(ctx context.Context, userID string, isSuperAdmin bool, status string) ([]*models.FormTemplate, error) {
	f := repository.ListFormTemplatesFilter{Status: status}
	if !isSuperAdmin {
		depts, err := s.depts.UserDepartmentIDs(ctx, userID)
		if err != nil {
			return nil, err
		}
		if len(depts) == 0 {
			return []*models.FormTemplate{}, nil
		}
		f.DepartmentIDs = depts
	}
	return s.forms.ListTemplates(ctx, f)
}

func (s *FormService) GetAdminTemplate(ctx context.Context, id, userID string, isSuperAdmin bool) (*models.FormTemplate, error) {
	t, err := s.forms.GetTemplate(ctx, id, true)
	if err != nil {
		return nil, err
	}
	if err := s.ensureCanManage(ctx, userID, isSuperAdmin, t.DepartmentID); err != nil {
		return nil, err
	}
	return t, nil
}

// ── Public portal ────────────────────────────────────────────────

func (s *FormService) ListPublicTemplates(ctx context.Context) ([]*models.FormTemplate, error) {
	return s.forms.ListTemplates(ctx, repository.ListFormTemplatesFilter{PublicOnly: true})
}

func (s *FormService) GetPublicTemplate(ctx context.Context, slug string) (*models.FormTemplate, error) {
	t, err := s.forms.GetTemplateBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if t.Status != models.FormStatusOpen {
		return nil, ErrTemplateClosed
	}
	return t, nil
}

// ── Submissions (applicant) ──────────────────────────────────────

type SaveDraftAnswersInput struct {
	TemplateID string
	Answers    json.RawMessage
}

func (s *FormService) SaveOrCreateDraft(ctx context.Context, userID string, in SaveDraftAnswersInput) (*models.FormSubmission, error) {
	t, err := s.forms.GetTemplate(ctx, in.TemplateID, false)
	if err != nil {
		return nil, err
	}
	if t.Status != models.FormStatusOpen {
		return nil, ErrTemplateClosed
	}
	existing, err := s.forms.GetActiveUserSubmissionForTemplate(ctx, userID, t.ID)
	if err == nil {
		if existing.Status != models.SubStatusDraft {
			return nil, ErrDuplicatePending
		}
		if err := s.forms.UpdateAnswers(ctx, existing.ID, in.Answers); err != nil {
			return nil, err
		}
		return s.forms.GetSubmission(ctx, existing.ID)
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	sub := &models.FormSubmission{
		TemplateID:    t.ID,
		DepartmentID:  t.DepartmentID,
		UserID:        userID,
		Status:        models.SubStatusDraft,
		PaymentStatus: "UNPAID",
		Answers:       in.Answers,
	}
	if err := s.forms.CreateSubmission(ctx, sub); err != nil {
		return nil, err
	}
	return s.forms.GetSubmission(ctx, sub.ID)
}

func (s *FormService) GetSubmissionForUser(ctx context.Context, id, userID string) (*models.FormSubmission, error) {
	sub, err := s.forms.GetSubmission(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub.UserID != userID {
		return nil, ErrNotOwner
	}
	return sub, nil
}

func (s *FormService) ListUserSubmissions(ctx context.Context, userID, status string, p *models.PaginationParams) ([]*models.FormSubmission, int64, error) {
	return s.forms.ListSubmissions(ctx, repository.ListSubmissionsFilter{
		UserID: userID,
		Status: status,
	}, p)
}

func (s *FormService) Submit(ctx context.Context, id, userID string) (*models.FormSubmission, error) {
	sub, err := s.forms.GetSubmission(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub.UserID != userID {
		return nil, ErrNotOwner
	}
	t, err := s.forms.GetTemplate(ctx, sub.TemplateID, false)
	if err != nil {
		return nil, err
	}
	// If paid template, payment must be uploaded/verified first.
	if t.PriceUGX > 0 && sub.PaymentStatus != "PAID" && sub.PaymentStatus != "PROOF_UPLOADED" {
		return nil, ErrPaymentRequired
	}
	ref := generateSubmissionReference(t.Slug)
	if err := s.forms.SetSubmitted(ctx, id, ref); err != nil {
		return nil, err
	}
	return s.forms.GetSubmission(ctx, id)
}

func (s *FormService) UploadPaymentProof(ctx context.Context, id, userID, reference, proofURL string, amount float64) error {
	sub, err := s.forms.GetSubmission(ctx, id)
	if err != nil {
		return err
	}
	if sub.UserID != userID {
		return ErrNotOwner
	}
	reference = strings.TrimSpace(reference)
	proofURL = strings.TrimSpace(proofURL)
	if (reference == "" && proofURL == "") || (reference != "" && proofURL != "") {
		return errors.New("provide either a PRN payment reference or an uploaded proof file, not both")
	}
	if amount <= 0 {
		return errors.New("payment_amount_ugx must be greater than 0")
	}
	return s.forms.SetPaymentProof(ctx, id, reference, proofURL, amount)
}

// ── Submissions (admin / dept-scoped) ────────────────────────────

func (s *FormService) ListSubmissions(ctx context.Context, userID string, isSuperAdmin bool, templateID, status string, p *models.PaginationParams) ([]*models.FormSubmission, int64, error) {
	f := repository.ListSubmissionsFilter{TemplateID: templateID, Status: status}
	if !isSuperAdmin {
		depts, err := s.depts.UserDepartmentIDs(ctx, userID)
		if err != nil {
			return nil, 0, err
		}
		if len(depts) == 0 {
			return []*models.FormSubmission{}, 0, nil
		}
		f.DepartmentIDs = depts
	}
	return s.forms.ListSubmissions(ctx, f, p)
}

func (s *FormService) GetSubmissionForAdmin(ctx context.Context, id, userID string, isSuperAdmin bool) (*models.FormSubmission, error) {
	sub, err := s.forms.GetSubmission(ctx, id)
	if err != nil {
		return nil, err
	}
	if !isSuperAdmin {
		depts, err := s.depts.UserDepartmentIDs(ctx, userID)
		if err != nil {
			return nil, err
		}
		ok := false
		for _, d := range depts {
			if d == sub.DepartmentID {
				ok = true
				break
			}
		}
		if !ok {
			return nil, ErrForbiddenDept
		}
	}
	return sub, nil
}

func (s *FormService) Review(ctx context.Context, id, status, reviewerID, notes string, isSuperAdmin bool) error {
	sub, err := s.GetSubmissionForAdmin(ctx, id, reviewerID, isSuperAdmin)
	if err != nil {
		return err
	}
	if sub.Status == models.SubStatusDraft {
		return ErrInvalidTransition
	}
	return s.forms.Review(ctx, id, status, reviewerID, notes)
}

func (s *FormService) VerifySubmissionPayment(ctx context.Context, id, reviewerID string, isSuperAdmin bool) error {
	if _, err := s.GetSubmissionForAdmin(ctx, id, reviewerID, isSuperAdmin); err != nil {
		return err
	}
	return s.forms.VerifyPayment(ctx, id, reviewerID)
}

// ── helpers ──────────────────────────────────────────────────────

func (s *FormService) normalizeAndSaveFields(ctx context.Context, templateID string, fields []*models.FormField) error {
	seenKeys := map[string]bool{}
	for i, f := range fields {
		if strings.TrimSpace(f.Label) == "" {
			return fmt.Errorf("field %d: label is required", i)
		}
		if f.FieldType == "" {
			return fmt.Errorf("field %d: field_type is required", i)
		}
		if f.FieldKey == "" {
			f.FieldKey = slugify(f.Label)
		}
		// Ensure uniqueness within the template.
		original := f.FieldKey
		n := 1
		for seenKeys[f.FieldKey] {
			n++
			f.FieldKey = fmt.Sprintf("%s-%d", original, n)
		}
		seenKeys[f.FieldKey] = true
		f.OrderIndex = i
		if len(f.Config) == 0 {
			f.Config = json.RawMessage("{}")
		}
	}
	return s.forms.ReplaceFields(ctx, templateID, fields)
}

func generateSubmissionReference(slug string) string {
	prefix := strings.ToUpper(strings.ReplaceAll(slug, "-", ""))
	if len(prefix) > 8 {
		prefix = prefix[:8]
	}
	return fmt.Sprintf("NCS-%d-%s-%s", time.Now().Year(), prefix, strings.ToUpper(uuid.NewString()[:6]))
}
