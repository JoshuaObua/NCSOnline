package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrNotOwner           = errors.New("not the application owner")
	ErrInvalidTransition  = errors.New("invalid status transition")
	ErrSubmitRequirements = errors.New("signed form and payment required before submission")
)

// Valid form types matching the official NCS forms
var validFormTypes = map[string]bool{
	"form_1":  true, // Declaration of National Sport
	"form_3":  true, // Registration/Renewal NSA/NSF
	"form_5":  true, // Transformation NSA to NSF
	"form_7":  true, // Organise Sports Competition
	"form_8":  true, // Operate Sports Facility
	"form_10": true, // Community Sports Club Registration/Renewal
	"form_11": true, // Operate Sports Academy
}

type ApplicationService struct {
	apps          *repository.ApplicationRepo
	audit         *repository.AuditRepo
	organisations *repository.OrganisationRepo
}

func NewApplicationService(apps *repository.ApplicationRepo, audit *repository.AuditRepo, organisations ...*repository.OrganisationRepo) *ApplicationService {
	s := &ApplicationService{apps: apps, audit: audit}
	if len(organisations) > 0 {
		s.organisations = organisations[0]
	}
	return s
}

type SaveDraftInput struct {
	FormData      json.RawMessage
	LastSavedStep int
	AppType       string
	OrgType       string
}

func (s *ApplicationService) SaveDraft(ctx context.Context, userID, formType string, in SaveDraftInput) (*models.Application, error) {
	if !validFormTypes[formType] {
		return nil, fmt.Errorf("invalid form type: %s", formType)
	}

	existing, err := s.apps.GetDraftByUserAndType(ctx, userID, formType)
	if err == nil {
		// Update existing draft
		existing.FormData = in.FormData
		existing.LastSavedStep = in.LastSavedStep
		if err := s.apps.UpdateDraft(ctx, existing); err != nil {
			return nil, fmt.Errorf("update draft: %w", err)
		}
		return existing, nil
	}

	// Create new draft
	exp := time.Now().Add(90 * 24 * time.Hour)
	app := &models.Application{
		ID:               uuid.NewString(),
		UserID:           userID,
		FormType:         formType,
		ApplicationType:  in.AppType,
		OrganisationType: in.OrgType,
		Status:           models.StatusDraft,
		PaymentStatus:    models.PaymentUnpaid,
		FormData:         in.FormData,
		LastSavedStep:    in.LastSavedStep,
		DraftExpiresAt:   &exp,
	}
	if err := s.apps.Create(ctx, app); err != nil {
		return nil, fmt.Errorf("create draft: %w", err)
	}
	return app, nil
}

func (s *ApplicationService) GetDraft(ctx context.Context, userID, formType string) (*models.Application, error) {
	return s.apps.GetDraftByUserAndType(ctx, userID, formType)
}

func (s *ApplicationService) DeleteDraft(ctx context.Context, userID, formType string) error {
	app, err := s.apps.GetDraftByUserAndType(ctx, userID, formType)
	if err != nil {
		return err
	}
	if app.UserID != userID {
		return ErrNotOwner
	}
	return s.apps.UpdateStatus(ctx, app.ID, "DELETED")
}

func (s *ApplicationService) GetByID(ctx context.Context, id, userID string, isAdmin bool) (*models.Application, error) {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !isAdmin && app.UserID != userID {
		return nil, ErrNotOwner
	}
	atts, _ := s.apps.ListAttachments(ctx, id)
	app.Attachments = atts
	return app, nil
}

func (s *ApplicationService) ListByUser(ctx context.Context, userID string, p *models.PaginationParams) ([]*models.Application, int64, error) {
	return s.apps.ListByUser(ctx, userID, p)
}

func (s *ApplicationService) AdminList(ctx context.Context, status string, p *models.PaginationParams) ([]*models.Application, int64, error) {
	return s.apps.AdminList(ctx, status, p)
}

func (s *ApplicationService) AdvanceToReview(ctx context.Context, id string) error {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app.Status != models.StatusDraft && app.Status != models.StatusPendingSignature {
		return ErrInvalidTransition
	}
	return s.apps.UpdateStatus(ctx, id, models.StatusPendingSignature)
}

func (s *ApplicationService) SetSignedForm(ctx context.Context, id, userID, url string) error {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app.UserID != userID {
		return ErrNotOwner
	}
	return s.apps.SetSignedForm(ctx, id, url)
}

type PaymentProofInput struct {
	Method    string
	Reference string
	Amount    float64
}

func (s *ApplicationService) UploadPaymentProof(ctx context.Context, id, userID string, in PaymentProofInput) error {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app.UserID != userID {
		return ErrNotOwner
	}
	if app.Status != models.StatusPendingPayment {
		return ErrInvalidTransition
	}
	return s.apps.SetPaymentProof(ctx, id, in.Method, in.Reference, in.Amount)
}

func (s *ApplicationService) AddAttachment(ctx context.Context, id, userID string, att *models.Attachment) error {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app.UserID != userID {
		return ErrNotOwner
	}
	att.ApplicationID = id
	return s.apps.AddAttachment(ctx, att)
}

func (s *ApplicationService) Submit(ctx context.Context, id, userID string) (*models.Application, error) {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app.UserID != userID {
		return nil, ErrNotOwner
	}
	if app.SignedFormURL == "" || (app.PaymentStatus != models.PaymentPaid && app.PaymentStatus != models.PaymentProofUploaded) {
		return nil, ErrSubmitRequirements
	}
	ref := generateReference(app.FormType)
	if err := s.apps.Submit(ctx, id, ref); err != nil {
		return nil, err
	}
	return s.apps.GetByID(ctx, id)
}

func (s *ApplicationService) Approve(ctx context.Context, id, reviewerID, notes string) error {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app.Status == models.StatusApproved && (app.FormType == "form_3" || app.FormType == "form_10") && s.organisations != nil {
		_, err = s.organisations.ApproveAndProvision(ctx, app, reviewerID, notes)
		return err
	}
	if app.Status != models.StatusUnderReview && app.Status != models.StatusSubmitted && app.Status != models.StatusResubmitted {
		return ErrInvalidTransition
	}
	if (app.FormType == "form_3" || app.FormType == "form_10") && s.organisations != nil {
		_, err = s.organisations.ApproveAndProvision(ctx, app, reviewerID, notes)
		return err
	}
	return s.apps.Approve(ctx, id, reviewerID, notes)
}

func (s *ApplicationService) Reject(ctx context.Context, id, reviewerID, notes string) error {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app.Status != models.StatusUnderReview && app.Status != models.StatusSubmitted && app.Status != models.StatusResubmitted {
		return ErrInvalidTransition
	}
	return s.apps.Reject(ctx, id, reviewerID, notes)
}

func (s *ApplicationService) RequestInfo(ctx context.Context, id, reviewerID, notes string) error {
	return s.apps.RequestInfo(ctx, id, reviewerID, notes)
}

func (s *ApplicationService) VerifyPayment(ctx context.Context, id, reviewerID string) error {
	return s.apps.VerifyPayment(ctx, id, reviewerID)
}

func (s *ApplicationService) RejectPayment(ctx context.Context, id, notes string) error {
	return s.apps.RejectPayment(ctx, id, notes)
}

func (s *ApplicationService) Respond(ctx context.Context, id, userID string, formData json.RawMessage) error {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app.UserID != userID {
		return ErrNotOwner
	}
	if app.Status != models.StatusNeedsInformation {
		return ErrInvalidTransition
	}
	app.FormData = formData
	if err := s.apps.UpdateDraft(ctx, app); err != nil {
		return err
	}
	return s.apps.UpdateStatus(ctx, id, models.StatusResubmitted)
}

func (s *ApplicationService) CountByStatus(ctx context.Context) (map[string]int64, error) {
	return s.apps.CountByStatus(ctx)
}

func generateReference(formType string) string {
	formMap := map[string]string{
		"form_1": "F1", "form_3": "F3", "form_5": "F5",
		"form_7": "F7", "form_8": "F8", "form_10": "F10", "form_11": "F11",
	}
	code := formMap[formType]
	if code == "" {
		code = "FX"
	}
	year := time.Now().Year()
	return fmt.Sprintf("NCS-%d-%s-%s", year, code, uuid.NewString()[:6])
}
