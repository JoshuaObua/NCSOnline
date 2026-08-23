package models

import "time"

type FederationLicense struct {
	ID                 string                  `json:"id"`
	FederationID       string                  `json:"federation_id"`
	FederationName     string                  `json:"federation_name,omitempty"`
	FederationAcronym  string                  `json:"federation_acronym,omitempty"`
	FederationRegNo    string                  `json:"federation_reg_no,omitempty"`
	FederationLogo     string                  `json:"federation_logo,omitempty"`
	PresidentName      string                  `json:"president_name,omitempty"`
	SecretaryName      string                  `json:"secretary_name,omitempty"`
	LicenseNumber      string                  `json:"license_number"`
	LicenseType        string                  `json:"license_type"` // FULL_RECOGNITION, PROVISIONAL, ANNUAL_COMPLIANCE, SPECIAL_CLEARANCE
	Category           string                  `json:"category"`     // Tier 1, Tier 2, etc.
	IssueDate          time.Time               `json:"issue_date"`
	ExpiryDate         time.Time               `json:"expiry_date"`
	Status             string                  `json:"status"` // ACTIVE, EXTENDED, EXPIRED, SUSPENDED, REVOKED
	Conditions         string                  `json:"conditions,omitempty"`
	DocumentURL        string                  `json:"document_url,omitempty"`
	IssuedBy           *string                 `json:"issued_by,omitempty"`
	IssuedByName       string                  `json:"issued_by_name,omitempty"`
	IssuedAt           time.Time               `json:"issued_at"`
	ExtendedBy         *string                 `json:"extended_by,omitempty"`
	ExtendedByName     string                  `json:"extended_by_name,omitempty"`
	ExtendedAt         *time.Time              `json:"extended_at,omitempty"`
	PreviousExpiryDate *time.Time              `json:"previous_expiry_date,omitempty"`
	ExtensionReason    string                  `json:"extension_reason,omitempty"`
	RevokedBy          *string                 `json:"revoked_by,omitempty"`
	RevokedByName      string                  `json:"revoked_by_name,omitempty"`
	RevokedAt          *time.Time              `json:"revoked_at,omitempty"`
	RevocationReason   string                  `json:"revocation_reason,omitempty"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
	Logs               []*FederationLicenseLog `json:"logs,omitempty"`
}

type FederationLicenseLog struct {
	ID              string     `json:"id"`
	LicenseID       string     `json:"license_id"`
	Action          string     `json:"action"` // CREATED, EXTENDED, REVOKED, REINSTATED, UPDATED
	PerformedBy     *string    `json:"performed_by,omitempty"`
	PerformedByName string     `json:"performed_by_name,omitempty"`
	Reason          string     `json:"reason,omitempty"`
	OldStatus       string     `json:"old_status,omitempty"`
	NewStatus       string     `json:"new_status,omitempty"`
	OldExpiryDate   *time.Time `json:"old_expiry_date,omitempty"`
	NewExpiryDate   *time.Time `json:"new_expiry_date,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type LicenseKPIs struct {
	TotalLicenses    int `json:"total_licenses"`
	ActiveLicenses   int `json:"active_licenses"`
	ExtendedLicenses int `json:"extended_licenses"`
	RevokedLicenses  int `json:"revoked_licenses"`
	ExpiringSoon     int `json:"expiring_soon"` // < 60 days
}
