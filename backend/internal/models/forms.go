package models

import (
	"encoding/json"
	"time"
)

// ── Departments ──────────────────────────────────────────────────

type Department struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ── Form Template & Fields ───────────────────────────────────────

const (
	FormStatusDraft    = "DRAFT"
	FormStatusOpen     = "OPEN"
	FormStatusClosed   = "CLOSED"
	FormStatusArchived = "ARCHIVED"
)

type FormTemplate struct {
	ID             string          `json:"id"`
	DepartmentID   string          `json:"department_id"`
	DepartmentName string          `json:"department_name,omitempty"`
	Slug           string          `json:"slug"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Sections       json.RawMessage `json:"sections"`
	BannerImageURL string          `json:"banner_image_url"`
	PriceUGX       float64         `json:"price_ugx"`
	Status         string          `json:"status"`
	LegacyFormType *string         `json:"legacy_form_type,omitempty"`
	CreatedBy      *string         `json:"created_by,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Fields         []*FormField    `json:"fields,omitempty"`
}

type FormField struct {
	ID          string          `json:"id"`
	TemplateID  string          `json:"template_id"`
	FieldKey    string          `json:"field_key"`
	FieldType   string          `json:"field_type"`
	Label       string          `json:"label"`
	Placeholder string          `json:"placeholder"`
	HelpText    string          `json:"help_text"`
	IsRequired  bool            `json:"is_required"`
	OrderIndex  int             `json:"order_index"`
	Config      json.RawMessage `json:"config"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// ── Form Submissions ─────────────────────────────────────────────

const (
	SubStatusDraft       = "DRAFT"
	SubStatusPendingPay  = "PENDING_PAYMENT"
	SubStatusSubmitted   = "SUBMITTED"
	SubStatusUnderReview = "UNDER_REVIEW"
	SubStatusNeedsInfo   = "NEEDS_INFORMATION"
	SubStatusApproved    = "APPROVED"
	SubStatusRejected    = "REJECTED"
)

type FormSubmission struct {
	ID                  string          `json:"id"`
	TemplateID          string          `json:"template_id"`
	TemplateTitle       string          `json:"template_title,omitempty"`
	DepartmentID        string          `json:"department_id"`
	UserID              string          `json:"user_id"`
	ApplicantName       string          `json:"applicant_name,omitempty"`
	ApplicantEmail      string          `json:"applicant_email,omitempty"`
	SubmissionReference string          `json:"submission_reference"`
	Status              string          `json:"status"`
	PaymentStatus       string          `json:"payment_status"`
	PaymentReference    string          `json:"payment_reference,omitempty"`
	PaymentProofURL     string          `json:"payment_proof_url,omitempty"`
	PaymentAmountUGX    *float64        `json:"payment_amount_ugx,omitempty"`
	PaymentVerifiedAt   *time.Time      `json:"payment_verified_at,omitempty"`
	Answers             json.RawMessage `json:"answers"`
	ReviewerID          *string         `json:"reviewer_id,omitempty"`
	ReviewNotes         string          `json:"review_notes,omitempty"`
	SubmittedAt         *time.Time      `json:"submitted_at,omitempty"`
	ApprovedAt          *time.Time      `json:"approved_at,omitempty"`
	RejectedAt          *time.Time      `json:"rejected_at,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}
