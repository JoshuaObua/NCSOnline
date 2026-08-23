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
	txs   *repository.TransactionRepo
	iotec *IoTecService
}

func NewFormService(forms *repository.FormRepo, depts *repository.DepartmentRepo, audit *repository.AuditRepo, txs *repository.TransactionRepo, iotec *IoTecService) *FormService {
	if iotec == nil {
		iotec = NewIoTecService()
	}
	return &FormService{forms: forms, depts: depts, audit: audit, txs: txs, iotec: iotec}
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
	DepartmentID          string
	Slug                  string
	Title                 string
	Description           string
	Sections              json.RawMessage
	BannerImageURL        string
	PriceUGX              float64
	AllowedPaymentMethods []string
	Status                string
	Fields                []*models.FormField
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
	var slug string
	if strings.TrimSpace(in.Slug) != "" {
		slug = slugify(in.Slug)
		// Check uniqueness
		other, err := s.forms.GetTemplateBySlug(ctx, slug)
		if err == nil && other != nil {
			return nil, fmt.Errorf("form slug %q is already in use by another form", slug)
		}
	} else {
		base := slugify(in.Title)
		slug = base + "-" + uuid.NewString()[:6]
	}

	payMethodsBytes := []byte(`["OVER_THE_COUNTER","MOBILE_MONEY"]`)
	if len(in.AllowedPaymentMethods) > 0 {
		if b, err := json.Marshal(in.AllowedPaymentMethods); err == nil {
			payMethodsBytes = b
		}
	}

	uid := userID
	t := &models.FormTemplate{
		DepartmentID:          in.DepartmentID,
		Slug:                  slug,
		Title:                 strings.TrimSpace(in.Title),
		Description:           in.Description,
		Sections:              sections,
		BannerImageURL:        in.BannerImageURL,
		PriceUGX:              in.PriceUGX,
		AllowedPaymentMethods: json.RawMessage(payMethodsBytes),
		Status:                st,
		CreatedBy:             &uid,
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
	if strings.TrimSpace(in.Slug) != "" {
		newSlug := slugify(in.Slug)
		if newSlug != existing.Slug {
			other, err := s.forms.GetTemplateBySlug(ctx, newSlug)
			if err == nil && other != nil && other.ID != existing.ID {
				return nil, fmt.Errorf("form slug %q is already in use by another form", newSlug)
			}
			existing.Slug = newSlug
		}
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
	if in.AllowedPaymentMethods != nil {
		if len(in.AllowedPaymentMethods) > 0 {
			if b, err := json.Marshal(in.AllowedPaymentMethods); err == nil {
				existing.AllowedPaymentMethods = json.RawMessage(b)
			}
		} else {
			existing.AllowedPaymentMethods = json.RawMessage(`["OVER_THE_COUNTER","MOBILE_MONEY"]`)
		}
	}
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
		if existing.Status != models.SubStatusDraft && existing.Status != models.SubStatusNeedsInfo {
			return nil, ErrDuplicatePending
		}
		if existing.Status == models.SubStatusNeedsInfo {
			err = s.forms.MakeEditable(ctx, existing.ID, in.Answers)
		} else {
			err = s.forms.UpdateAnswers(ctx, existing.ID, in.Answers)
		}
		if err != nil {
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
	status, err = normalizeSubmissionReviewStatus(status)
	if err != nil {
		return err
	}
	notes = strings.TrimSpace(notes)
	if (status == models.SubStatusNeedsInfo || status == models.SubStatusRejected) && notes == "" {
		return errors.New("review notes are required for queried or rejected applications")
	}
	if err := s.forms.Review(ctx, id, status, reviewerID, notes); err != nil {
		return err
	}
	updated, err := s.forms.GetSubmission(ctx, id)
	if err != nil {
		return err
	}
	return s.notifyReviewDecision(ctx, updated, reviewerID, status, notes)
}

func (s *FormService) VerifySubmissionPayment(ctx context.Context, id, reviewerID string, isSuperAdmin bool) error {
	return s.UpdateSubmissionPaymentStatus(ctx, id, "PAID", reviewerID, isSuperAdmin)
}

func (s *FormService) UpdateSubmissionPaymentStatus(ctx context.Context, id, status, reviewerID string, isSuperAdmin bool) error {
	if _, err := s.GetSubmissionForAdmin(ctx, id, reviewerID, isSuperAdmin); err != nil {
		return err
	}
	status = strings.ToUpper(strings.TrimSpace(status))
	if status == "" || status == "PENDING" {
		status = "PROOF_UPLOADED"
	}
	switch status {
	case "PROOF_UPLOADED", "VERIFICATION_FAILED", "PAID":
	default:
		return fmt.Errorf("invalid payment status: %s", status)
	}
	if err := s.forms.SetPaymentStatus(ctx, id, status, reviewerID); err != nil {
		return err
	}
	updated, err := s.forms.GetSubmission(ctx, id)
	if err != nil {
		return err
	}
	return s.notifyPaymentDecision(ctx, updated, reviewerID, status)
}

// ── helpers ──────────────────────────────────────────────────────

func normalizeSubmissionReviewStatus(status string) (string, error) {
	status = strings.ToUpper(strings.TrimSpace(status))
	switch status {
	case models.SubStatusUnderReview, models.SubStatusNeedsInfo, models.SubStatusComplete, models.SubStatusApproved, models.SubStatusRejected:
		return status, nil
	case "QUERIED":
		return models.SubStatusNeedsInfo, nil
	default:
		return "", fmt.Errorf("invalid status: %s", status)
	}
}

func (s *FormService) notifyReviewDecision(ctx context.Context, sub *models.FormSubmission, reviewerID, status, notes string) error {
	if sub == nil || sub.UserID == "" {
		return nil
	}
	reviewerName := sub.ReviewerName
	if reviewerName == "" {
		name, err := s.forms.UserDisplayName(ctx, reviewerID)
		if err != nil {
			return err
		}
		reviewerName = name
	}
	formName := sub.TemplateTitle
	if formName == "" {
		formName = "your application"
	}
	messageSuffix := ""
	if notes != "" {
		messageSuffix = " Reason: " + notes
	}
	switch status {
	case models.SubStatusNeedsInfo:
		return s.forms.CreateSubmissionNotification(ctx, sub.UserID, "application_queried", "Application queried", fmt.Sprintf("%s queried %s.%s You can edit and resubmit it from your dashboard.", reviewerName, formName, messageSuffix), "question-circle")
	case models.SubStatusRejected:
		return s.forms.CreateSubmissionNotification(ctx, sub.UserID, "application_rejected", "Application rejected", fmt.Sprintf("%s rejected %s.%s", reviewerName, formName, messageSuffix), "close-circled")
	case models.SubStatusApproved:
		return s.forms.CreateSubmissionNotification(ctx, sub.UserID, "application_approved", "Application approved", fmt.Sprintf("%s approved %s.%s", reviewerName, formName, messageSuffix), "check-circled")
	case models.SubStatusComplete:
		return s.forms.CreateSubmissionNotification(ctx, sub.UserID, "application_complete", "Application marked complete", fmt.Sprintf("%s marked %s as complete.%s", reviewerName, formName, messageSuffix), "check")
	default:
		return nil
	}
}

func (s *FormService) notifyPaymentDecision(ctx context.Context, sub *models.FormSubmission, reviewerID, status string) error {
	if sub == nil || sub.UserID == "" {
		return nil
	}
	reviewerName := sub.PaymentVerifierName
	if reviewerName == "" {
		name, err := s.forms.UserDisplayName(ctx, reviewerID)
		if err != nil {
			return err
		}
		reviewerName = name
	}
	formName := sub.TemplateTitle
	if formName == "" {
		formName = "your application"
	}
	switch status {
	case "PAID":
		return s.forms.CreateSubmissionNotification(ctx, sub.UserID, "payment_verified", "Payment verified", fmt.Sprintf("%s verified the payment for %s.", reviewerName, formName), "check-circled")
	case "VERIFICATION_FAILED":
		return s.forms.CreateSubmissionNotification(ctx, sub.UserID, "payment_failed", "Payment verification failed", fmt.Sprintf("%s marked the payment for %s as verification failed.", reviewerName, formName), "warning")
	default:
		return nil
	}
}

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

// ── Mobile Money Payments & ioTec Integration ────────────────────

func (s *FormService) InitiateMoMoPayment(ctx context.Context, submissionID, userID, phoneNumber string) (*models.PaymentTransaction, error) {
	sub, err := s.forms.GetSubmission(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if sub.UserID != userID {
		return nil, ErrNotOwner
	}
	if sub.PaymentStatus == "PAID" {
		return nil, errors.New("application payment has already been completed")
	}

	tmpl, err := s.forms.GetTemplate(ctx, sub.TemplateID, false)
	if err != nil {
		return nil, err
	}
	if tmpl.PriceUGX <= 0 {
		return nil, errors.New("this form is free of charge; no payment required")
	}

	cleanPhone := strings.TrimSpace(phoneNumber)
	cleanPhone = strings.ReplaceAll(cleanPhone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	if cleanPhone == "" {
		return nil, errors.New("phone number is required for mobile money payment")
	}

	txRef := fmt.Sprintf("NCS-TXN-%s-%s", time.Now().Format("20060102"), strings.ToUpper(uuid.NewString()[:8]))

	tx := &models.PaymentTransaction{
		ID:                   uuid.NewString(),
		TransactionReference: txRef,
		SubmissionID:         &sub.ID,
		TemplateID:           &tmpl.ID,
		UserID:               userID,
		PaymentMethod:        models.PaymentMethodMoMo,
		Provider:             "IOTEC",
		PhoneNumber:          cleanPhone,
		AmountUGX:            tmpl.PriceUGX,
		Currency:             "UGX",
		Status:               models.TxStatusPending,
		StatusMessage:        "Initiating Mobile Money collection request...",
		RawResponse:          json.RawMessage("{}"),
	}

	if err := s.txs.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to record payment transaction: %w", err)
	}

	// Initiate ioTec collection
	colReq := IoTecCollectionRequest{
		Category:   "MobileMoney",
		Currency:   "UGX",
		ExternalID: tx.TransactionReference,
		Payer:      cleanPhone,
		PayerName:  sub.ApplicantName,
		PayerNote:  fmt.Sprintf("Payment for %s", tmpl.Title),
		Amount:     tmpl.PriceUGX,
		PayeeNote:  fmt.Sprintf("Submission %s", sub.ID),
	}

	colResp, err := s.iotec.InitiateCollection(ctx, colReq)
	if err != nil {
		msg := fmt.Sprintf("Payment initiation failed: %v", err)
		_ = s.txs.UpdateStatus(ctx, tx.ID, models.TxStatusFailed, msg, nil)
		tx.Status = models.TxStatusFailed
		tx.StatusMessage = msg
		return tx, err
	}

	rawBytes, _ := json.Marshal(colResp)
	tx.ProviderRequestID = &colResp.ID
	tx.StatusMessage = colResp.StatusMessage
	if tx.StatusMessage == "" {
		tx.StatusMessage = "USSD push sent. Awaiting customer PIN authorization."
	}
	if strings.EqualFold(colResp.Status, "Success") {
		tx.Status = models.TxStatusSuccess
		_ = s.txs.UpdateStatusWithProvider(ctx, tx.ID, colResp.ID, models.TxStatusSuccess, tx.StatusMessage, rawBytes)
		_ = s.forms.SetPaymentPaidDirect(ctx, sub.ID, models.PaymentMethodMoMo, tx.TransactionReference, tx.AmountUGX)
		_ = s.forms.CreateSubmissionNotification(ctx, sub.UserID, "payment_success", "Payment successful", fmt.Sprintf("Your mobile money payment of UGX %.0f for %s has been confirmed.", tx.AmountUGX, tmpl.Title), "check-circled")
	} else if strings.EqualFold(colResp.Status, "Failed") {
		tx.Status = models.TxStatusFailed
		_ = s.txs.UpdateStatusWithProvider(ctx, tx.ID, colResp.ID, models.TxStatusFailed, tx.StatusMessage, rawBytes)
	} else {
		_ = s.txs.UpdateStatusWithProvider(ctx, tx.ID, colResp.ID, models.TxStatusPending, tx.StatusMessage, rawBytes)
	}

	return s.txs.GetByID(ctx, tx.ID)
}

func (s *FormService) CheckSubmissionPaymentStatus(ctx context.Context, submissionID, userID string) (*models.PaymentTransaction, error) {
	sub, err := s.forms.GetSubmission(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if sub.UserID != userID {
		return nil, ErrNotOwner
	}

	tx, err := s.txs.GetLatestBySubmissionID(ctx, submissionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("no payment transaction found for this submission")
		}
		return nil, err
	}

	if tx.Status == models.TxStatusPending && tx.Provider == "IOTEC" && tx.ProviderRequestID != nil && *tx.ProviderRequestID != "" {
		statusResp, err := s.iotec.GetStatus(ctx, *tx.ProviderRequestID, tx.PhoneNumber)
		if err == nil && statusResp != nil {
			rawBytes, _ := json.Marshal(statusResp)
			if strings.EqualFold(statusResp.Status, "Success") {
				_ = s.txs.UpdateStatus(ctx, tx.ID, models.TxStatusSuccess, statusResp.StatusMessage, rawBytes)
				_ = s.forms.SetPaymentPaidDirect(ctx, sub.ID, models.PaymentMethodMoMo, tx.TransactionReference, tx.AmountUGX)
				_ = s.forms.CreateSubmissionNotification(ctx, sub.UserID, "payment_success", "Payment confirmed", fmt.Sprintf("Your mobile money payment for %s was confirmed.", tx.TemplateTitle), "check-circled")
			} else if strings.EqualFold(statusResp.Status, "Failed") {
				_ = s.txs.UpdateStatus(ctx, tx.ID, models.TxStatusFailed, statusResp.StatusMessage, rawBytes)
			}
		}
	}

	return s.txs.GetByID(ctx, tx.ID)
}

func (s *FormService) SyncPendingTransactions(ctx context.Context) (int, error) {
	pending, err := s.txs.ListPendingForSync(ctx, 30)
	if err != nil {
		return 0, err
	}

	synced := 0
	for _, tx := range pending {
		if tx.ProviderRequestID == nil || *tx.ProviderRequestID == "" {
			continue
		}
		statusResp, err := s.iotec.GetStatus(ctx, *tx.ProviderRequestID, tx.PhoneNumber)
		if err != nil {
			continue
		}
		rawBytes, _ := json.Marshal(statusResp)
		if strings.EqualFold(statusResp.Status, "Success") {
			_ = s.txs.UpdateStatus(ctx, tx.ID, models.TxStatusSuccess, statusResp.StatusMessage, rawBytes)
			if tx.SubmissionID != nil && *tx.SubmissionID != "" {
				_ = s.forms.SetPaymentPaidDirect(ctx, *tx.SubmissionID, models.PaymentMethodMoMo, tx.TransactionReference, tx.AmountUGX)
				_ = s.forms.CreateSubmissionNotification(ctx, tx.UserID, "payment_success", "Payment confirmed", fmt.Sprintf("Your payment for transaction %s was confirmed.", tx.TransactionReference), "check-circled")
			}
			synced++
		} else if strings.EqualFold(statusResp.Status, "Failed") {
			_ = s.txs.UpdateStatus(ctx, tx.ID, models.TxStatusFailed, statusResp.StatusMessage, rawBytes)
			synced++
		}
	}
	return synced, nil
}

func (s *FormService) ListTransactions(ctx context.Context, f repository.ListTransactionsFilter, p *models.PaginationParams) ([]*models.PaymentTransaction, int, error) {
	return s.txs.List(ctx, f, p)
}

func (s *FormService) GetTransactionKPIs(ctx context.Context, userID string) (*repository.TransactionKPIs, error) {
	return s.txs.GetKPIs(ctx, userID)
}

