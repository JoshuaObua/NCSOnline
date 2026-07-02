package models

import (
	"encoding/json"
	"time"
)

// ── User ─────────────────────────────────────────────────────────

type User struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	PinHash           string     `json:"-"`
	PinChangeRequired bool       `json:"pin_change_required"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Phone             string     `json:"phone,omitempty"`
	AvatarURL         string     `json:"avatar_url,omitempty"`
	AuthProvider      string     `json:"auth_provider,omitempty"`
	GoogleSub         string     `json:"-"`
	IsActive          bool       `json:"is_active"`
	AccountStatus     string     `json:"account_status"`
	StatusReason      string     `json:"status_reason,omitempty"`
	FraudFlag         bool       `json:"fraud_flag"`
	FraudReason       string     `json:"fraud_reason,omitempty"`
	SuspendedUntil    *time.Time `json:"suspended_until,omitempty"`
	StatusChangedAt   *time.Time `json:"status_changed_at,omitempty"`
	IsEmailVerified   bool       `json:"is_email_verified"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	AuthInvalidBefore *time.Time `json:"-"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	Roles             []Role     `json:"roles,omitempty"`
}

// HasPIN reports whether the user has set a screen-lock PIN.
func (u *User) HasPIN() bool { return u.PinHash != "" }

// MarshalJSON adds the derived `has_pin` field so the frontend can
// toggle the Set PIN / Change PIN UI without exposing the hash itself.
func (u User) MarshalJSON() ([]byte, error) {
	type alias User
	return json.Marshal(&struct {
		alias
		HasPIN bool `json:"has_pin"`
	}{alias(u), u.PinHash != ""})
}

func (u *User) FullName() string { return u.FirstName + " " + u.LastName }

func (u *User) HasRole(name string) bool {
	for _, r := range u.Roles {
		if r.Name == name {
			return true
		}
	}
	return false
}

func (u *User) RoleNames() []string {
	names := make([]string, len(u.Roles))
	for i, r := range u.Roles {
		names[i] = r.Name
	}
	return names
}

// ── Role & Permission ─────────────────────────────────────────────

type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	IsSystem    bool         `json:"is_system"`
	CreatedAt   time.Time    `json:"created_at"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
}

type Notification struct {
	ID        string    `json:"id"`
	UserID    *string   `json:"user_id,omitempty"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	IconKey   string    `json:"icon_key"`
	CreatedAt time.Time `json:"created_at"`
}

type ContactMessage struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ── Auth Tokens ───────────────────────────────────────────────────

type RefreshToken struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	IPAddress string     `json:"ip_address,omitempty"`
	UserAgent string     `json:"user_agent,omitempty"`
}

type ActiveSession struct {
	ID        string    `json:"id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ── Application ───────────────────────────────────────────────────

type Application struct {
	ID                   string          `json:"id"`
	ApplicationReference string          `json:"application_reference,omitempty"`
	UserID               string          `json:"user_id"`
	FormType             string          `json:"form_type"`
	ApplicationType      string          `json:"application_type,omitempty"`
	OrganisationType     string          `json:"organisation_type,omitempty"`
	Status               string          `json:"status"`
	PaymentStatus        string          `json:"payment_status"`
	PaymentMethod        string          `json:"payment_method,omitempty"`
	PaymentReference     string          `json:"payment_reference,omitempty"`
	PaymentAmountUGX     *float64        `json:"payment_amount_ugx,omitempty"`
	PaymentVerifiedAt    *time.Time      `json:"payment_verified_at,omitempty"`
	PaymentVerifiedBy    *string         `json:"payment_verified_by,omitempty"`
	FormData             json.RawMessage `json:"form_data"`
	SignedFormURL        string          `json:"signed_form_url,omitempty"`
	SignedFormUploadedAt *time.Time      `json:"signed_form_uploaded_at,omitempty"`
	PDFGeneratedAt       *time.Time      `json:"pdf_generated_at,omitempty"`
	SubmittedAt          *time.Time      `json:"submitted_at,omitempty"`
	ReviewerID           *string         `json:"reviewer_id,omitempty"`
	ReviewNotes          string          `json:"review_notes,omitempty"`
	ApprovedAt           *time.Time      `json:"approved_at,omitempty"`
	RejectedAt           *time.Time      `json:"rejected_at,omitempty"`
	DraftExpiresAt       *time.Time      `json:"draft_expires_at,omitempty"`
	LastSavedStep        int             `json:"last_saved_step"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
	Attachments          []Attachment    `json:"attachments,omitempty"`
}

type Attachment struct {
	ID            string    `json:"id"`
	ApplicationID string    `json:"application_id"`
	FieldName     string    `json:"field_name"`
	FileName      string    `json:"file_name"`
	FileURL       string    `json:"file_url"`
	FileSize      int64     `json:"file_size,omitempty"`
	MimeType      string    `json:"mime_type,omitempty"`
	UploadedAt    time.Time `json:"uploaded_at"`
}

type Organisation struct {
	ID                 string          `json:"id"`
	ProfileReference   string          `json:"profile_reference"`
	OrganisationType   string          `json:"organisation_type"`
	LegalName          string          `json:"legal_name"`
	DisplayName        string          `json:"display_name"`
	OfficialEmail      string          `json:"official_email"`
	OfficialPhone      string          `json:"official_phone"`
	RegistrationNumber string          `json:"registration_number"`
	Status             string          `json:"status"`
	Role               string          `json:"role,omitempty"`
	ProfileData        json.RawMessage `json:"profile_data,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// ── Audit Log ─────────────────────────────────────────────────────

type AuditLog struct {
	ID             string          `json:"id"`
	UserID         *string         `json:"user_id,omitempty"`
	Action         string          `json:"action"`
	Resource       string          `json:"resource"`
	ResourceID     string          `json:"resource_id,omitempty"`
	OldValues      json.RawMessage `json:"old_values,omitempty"`
	NewValues      json.RawMessage `json:"new_values,omitempty"`
	IPAddress      string          `json:"ip_address,omitempty"`
	UserAgent      string          `json:"user_agent,omitempty"`
	Method         string          `json:"method,omitempty"`
	Endpoint       string          `json:"endpoint,omitempty"`
	ResponseCode   int             `json:"response_code"`
	ResponseTimeMs int64           `json:"response_time_ms"`
	DeviceInfo     string          `json:"device_info,omitempty"`
	// Enhanced audit fields
	EventType       string          `json:"event_type,omitempty"`
	EventStatus     string          `json:"event_status,omitempty"`
	SeverityLevel   string          `json:"severity_level,omitempty"`
	ForwardedIP     string          `json:"forwarded_ip,omitempty"`
	GeoCountry      string          `json:"geo_country,omitempty"`
	GeoCity         string          `json:"geo_city,omitempty"`
	GeoRegion       string          `json:"geo_region,omitempty"`
	GeoLatitude     *float64        `json:"geo_latitude,omitempty"`
	GeoLongitude    *float64        `json:"geo_longitude,omitempty"`
	GeoTimezone     string          `json:"geo_timezone,omitempty"`
	GeoSource       string          `json:"geo_source,omitempty"`
	Platform        string          `json:"platform,omitempty"`
	Authenticated   bool            `json:"authenticated"`
	PayloadExcerpt  json.RawMessage `json:"payload_excerpt,omitempty"`
	PreviousHash    string          `json:"previous_hash,omitempty"`
	EntryHash       string          `json:"entry_hash,omitempty"`
	ChainSequence   int64           `json:"chain_sequence,omitempty"`
	VPNDetected     bool            `json:"vpn_detected"`
	Browser         string          `json:"browser,omitempty"`
	OSName          string          `json:"os_name,omitempty"`
	ClientType      string          `json:"client_type,omitempty"`
	ThreatScore     int             `json:"threat_score"`
	AnomalyDetected bool            `json:"anomaly_detected"`
	SessionID       string          `json:"session_id,omitempty"`
	Username        string          `json:"username,omitempty"`
	FirstName       string          `json:"first_name,omitempty"`
	LastName        string          `json:"last_name,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// ── CMS ───────────────────────────────────────────────────────────

type CMSPost struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Content         string     `json:"content"`
	Excerpt         string     `json:"excerpt"`
	Category        string     `json:"category"`
	Status          string     `json:"status"`
	CoverImageURL   string     `json:"cover_image_url,omitempty"`
	AuthorID        *string    `json:"author_id,omitempty"`
	AuthorName      string     `json:"author_name,omitempty"`
	MetaTitle       string     `json:"meta_title,omitempty"`
	MetaDescription string     `json:"meta_description,omitempty"`
	FocusKeywords   string     `json:"focus_keywords,omitempty"`
	ViewCount       int        `json:"view_count"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	PublishedAt     *time.Time `json:"published_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type BlogCategory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BlogComment struct {
	ID         string         `json:"id"`
	PostID     string         `json:"post_id"`
	PostTitle  string         `json:"post_title,omitempty"`
	PostSlug   string         `json:"post_slug,omitempty"`
	UserID     *string        `json:"user_id,omitempty"`
	UserName   string         `json:"user_name,omitempty"`
	ParentID   *string        `json:"parent_id,omitempty"`
	Body       string         `json:"body"`
	Status     string         `json:"status"`
	Depth      int            `json:"depth"`
	IPAddress  string         `json:"ip_address,omitempty"`
	UserAgent  string         `json:"user_agent,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	ApprovedAt *time.Time     `json:"approved_at,omitempty"`
	FlaggedAt  *time.Time     `json:"flagged_at,omitempty"`
	Replies    []*BlogComment `json:"replies,omitempty"`
}

type CMSEvent struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Slug          string     `json:"slug"`
	Description   string     `json:"description"`
	Category      string     `json:"category,omitempty"`
	Location      string     `json:"location,omitempty"`
	EventDate     *time.Time `json:"event_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	CoverImageURL string     `json:"cover_image_url,omitempty"`
	Status        string     `json:"status"`
	AuthorID      *string    `json:"author_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CMSCareer struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Department     string     `json:"department,omitempty"`
	DepartmentID   *string    `json:"department_id,omitempty"`
	DepartmentName string     `json:"department_name,omitempty"`
	Location       string     `json:"location,omitempty"`
	JobType        string     `json:"job_type"`
	Category       string     `json:"category"`
	Description    string     `json:"description"`
	Requirements   string     `json:"requirements,omitempty"`
	SalaryRange    string     `json:"salary_range,omitempty"`
	Status         string     `json:"status"`
	DeadlineAt     *time.Time `json:"deadline_at,omitempty"`
	AuthorID       *string    `json:"author_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CMSSlide struct {
	ID            string            `json:"id"`
	SlideshowID   string            `json:"slideshow_id,omitempty"`
	Title         string            `json:"title"`
	Subtitle      string            `json:"subtitle,omitempty"`
	Description   string            `json:"description,omitempty"`
	Category      string            `json:"category,omitempty"`
	ImageURL      string            `json:"image_url,omitempty"`
	MediaType     string            `json:"media_type,omitempty"`
	AnimationType string            `json:"animation_type,omitempty"`
	ButtonText    string            `json:"button_text,omitempty"`
	ButtonURL     string            `json:"button_url,omitempty"`
	Buttons       []*CMSSlideButton `json:"buttons,omitempty"`
	SortOrder     int               `json:"sort_order"`
	IsActive      bool              `json:"is_active"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type CMSSlideshow struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	Slug               string      `json:"slug"`
	TransitionEffect   string      `json:"transition_effect"`
	TransitionDuration int         `json:"transition_duration"`
	AutoplaySpeed      int         `json:"autoplay_speed"`
	PauseOnHover       bool        `json:"pause_on_hover"`
	IsActive           bool        `json:"is_active"`
	Slides             []*CMSSlide `json:"slides,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

type CMSSlideButton struct {
	ID         string    `json:"id"`
	SlideID    string    `json:"slide_id,omitempty"`
	Text       string    `json:"text"`
	URL        string    `json:"url"`
	StyleClass string    `json:"style_class"`
	LinkTarget string    `json:"link_target"`
	SortOrder  int       `json:"sort_order"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type CMSMenuItem struct {
	ID        string        `json:"id"`
	Label     string        `json:"label"`
	URL       string        `json:"url"`
	Target    string        `json:"target,omitempty"`
	Hidden    bool          `json:"hidden,omitempty"`
	Icon      string        `json:"icon,omitempty"`
	Children  []CMSMenuItem `json:"children"`
	Mega      bool          `json:"mega,omitempty"`
	MegaItems []CMSMenuItem `json:"megaItems,omitempty"`
}

type CMSMenu struct {
	Name      string        `json:"name"`
	Items     []CMSMenuItem `json:"items"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// CMSSetting is a key/value record used for site-wide configuration:
// footer config, contact info, etc. Value is opaque JSON so the
// frontend can store whatever shape it needs without backend changes.
type CMSSetting struct {
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CMSFunFact struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	Value     string    `json:"value"`
	Icon      string    `json:"icon,omitempty"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CMSFAQ struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Category  string    `json:"category"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewsletterSubscriber struct {
	ID             string     `json:"id"`
	Email          string     `json:"email"`
	Source         string     `json:"source"`
	IsActive       bool       `json:"is_active"`
	UnsubscribedAt *time.Time `json:"unsubscribed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CMSResource struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	FileURL     string    `json:"file_url,omitempty"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CMSFacility struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	Category    string    `json:"category,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CMSAssociation struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	LogoURL     string    `json:"logo_url,omitempty"`
	WebsiteURL  string    `json:"website_url,omitempty"`
	Category    string    `json:"category,omitempty"`
	President   string    `json:"president,omitempty"`
	Secretary   string    `json:"secretary,omitempty"`
	Address     string    `json:"address,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CMSInvest struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle,omitempty"`
	Content   string    `json:"content,omitempty"`
	ImageURL  string    `json:"image_url,omitempty"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CMSTeamMember struct {
	ID             string    `json:"id"`
	FullName       string    `json:"full_name"`
	Designation    string    `json:"designation,omitempty"`
	ImageURL       string    `json:"image_url,omitempty"`
	Bio            string    `json:"bio,omitempty"`
	SortOrder      int       `json:"sort_order"`
	IsActive       bool      `json:"is_active"`
	DepartmentID   *string   `json:"department_id,omitempty"`
	DepartmentName string    `json:"department_name,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ── NSMIS federation reporting ───────────────────────────────────

type Federation struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Acronym               string    `json:"acronym"`
	NCSRegistrationNumber string    `json:"ncs_registration_number"`
	RecognitionStatus     string    `json:"recognition_status"`
	PhysicalAddress       string    `json:"physical_address"`
	Email                 string    `json:"email"`
	Website               string    `json:"website"`
	ContactPerson         string    `json:"contact_person"`
	IsActive              bool      `json:"is_active"`
	Version               int       `json:"version"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type ReportingPeriod struct {
	ID          string     `json:"id"`
	PeriodType  string     `json:"period_type"`
	Name        string     `json:"name"`
	StartsOn    time.Time  `json:"starts_on"`
	EndsOn      time.Time  `json:"ends_on"`
	DueOn       time.Time  `json:"due_on"`
	GraceEndsOn *time.Time `json:"grace_ends_on,omitempty"`
	Timezone    string     `json:"timezone"`
	IsOpen      bool       `json:"is_open"`
	CreatedAt   time.Time  `json:"created_at"`
}

type ReportObligation struct {
	ID                string    `json:"id"`
	FederationID      string    `json:"federation_id"`
	FederationName    string    `json:"federation_name"`
	ReportingPeriodID string    `json:"reporting_period_id"`
	PeriodName        string    `json:"period_name"`
	ReportType        string    `json:"report_type"`
	DueOn             time.Time `json:"due_on"`
	Status            string    `json:"status"`
	DaysOverdue       int       `json:"days_overdue"`
	ReportID          *string   `json:"report_id,omitempty"`
	ReportVersion     *int      `json:"report_version,omitempty"`
	ReportStatus      *string   `json:"report_status,omitempty"`
}

type GovernanceFederationRow struct {
	FederationID        string   `json:"federation_id"`
	FederationName      string   `json:"federation_name"`
	Acronym             string   `json:"acronym"`
	ComplianceScore     *float64 `json:"compliance_score,omitempty"`
	IsCompliant         *bool    `json:"is_compliant,omitempty"`
	MissingReports      int64    `json:"missing_reports"`
	ExpiredConstitution bool     `json:"expired_constitution"`
}

type GovernanceDashboard struct {
	ReportingPeriodID    string                    `json:"reporting_period_id,omitempty"`
	ReportingPeriodName  string                    `json:"reporting_period_name,omitempty"`
	TotalFederations     int64                     `json:"total_federations"`
	CompliantFederations int64                     `json:"compliant_federations"`
	NonCompliant         int64                     `json:"non_compliant_federations"`
	MissingReports       int64                     `json:"missing_reports"`
	ExpiredConstitutions int64                     `json:"expired_constitutions"`
	AverageScore         *float64                  `json:"average_compliance_score,omitempty"`
	Rows                 []GovernanceFederationRow `json:"federations"`
	AsOf                 time.Time                 `json:"as_of"`
	CalculationVersion   int                       `json:"calculation_version"`
}

// ── Pagination ────────────────────────────────────────────────────

type PaginationParams struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Search  string `json:"search,omitempty"`
	SortBy  string `json:"sort_by,omitempty"`
	Order   string `json:"order,omitempty"`
}

func (p *PaginationParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 || p.PerPage > 100 {
		p.PerPage = 20
	}
	return (p.Page - 1) * p.PerPage
}

// ── Application Status & Payment Status constants ─────────────────

const (
	AccountStatusActive    = "ACTIVE"
	AccountStatusSuspended = "SUSPENDED"
	AccountStatusBanned    = "BANNED"

	StatusDraft            = "DRAFT"
	StatusPendingSignature = "PENDING_SIGNATURE"
	StatusPendingPayment   = "PENDING_PAYMENT"
	StatusSubmitted        = "SUBMITTED"
	StatusUnderReview      = "UNDER_REVIEW"
	StatusNeedsInformation = "NEEDS_INFORMATION"
	StatusResubmitted      = "RESUBMITTED"
	StatusApproved         = "APPROVED"
	StatusRejected         = "REJECTED"

	PaymentUnpaid        = "UNPAID"
	PaymentInitiated     = "PAYMENT_INITIATED"
	PaymentPaid          = "PAID"
	PaymentProofUploaded = "PROOF_UPLOADED"
	PaymentVerified      = "PAYMENT_VERIFIED"
	PaymentRejected      = "PAYMENT_REJECTED"

	PaymentMethodOnline = "ONLINE"
	PaymentMethodProof  = "PROOF_UPLOAD"
)

// ── Context keys ──────────────────────────────────────────────────

type contextKey string

const (
	CtxUserID        contextKey = "ctx_user_id"
	CtxUserEmail     contextKey = "ctx_user_email"
	CtxUserRoles     contextKey = "ctx_user_roles"
	CtxSessionID     contextKey = "ctx_session_id"
	CtxTokenIssuedAt contextKey = "ctx_token_issued_at"
)
