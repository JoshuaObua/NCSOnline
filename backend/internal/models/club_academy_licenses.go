package models

import (
	"time"
)

type ClubAcademyLicense struct {
	ID                 string                    `json:"id"`
	ClubID             string                    `json:"club_id"`
	ClubName           string                    `json:"club_name,omitempty"`
	ClubNumber         string                    `json:"club_number,omitempty"`
	ClubAcronym        string                    `json:"club_acronym,omitempty"`
	ClubDistrict       string                    `json:"club_district,omitempty"`
	ClubRegion         string                    `json:"club_region,omitempty"`
	ClubEmail          string                    `json:"club_email,omitempty"`
	ClubPhone          string                    `json:"club_phone,omitempty"`
	ClubLogo           string                    `json:"club_logo,omitempty"`
	PresidentName      string                    `json:"president_name,omitempty"`
	SecretaryName      string                    `json:"secretary_name,omitempty"`
	FederationID       string                    `json:"federation_id,omitempty"`
	FederationName     string                    `json:"federation_name,omitempty"`
	LicenseNumber      string                    `json:"license_number"`
	LicenseType        string                    `json:"license_type"`
	Category           string                    `json:"category"`
	IssueDate          time.Time                 `json:"issue_date"`
	ExpiryDate         time.Time                 `json:"expiry_date"`
	Status             string                    `json:"status"` // ACTIVE, EXTENDED, EXPIRED, REVOKED, SUSPENDED
	Conditions         string                    `json:"conditions"`
	DocumentURL        string                    `json:"document_url,omitempty"`
	IssuedBy           *string                   `json:"issued_by,omitempty"`
	IssuedByName       string                    `json:"issued_by_name,omitempty"`
	IssuedAt           time.Time                 `json:"issued_at"`
	ExtendedBy         *string                   `json:"extended_by,omitempty"`
	ExtendedByName     string                    `json:"extended_by_name,omitempty"`
	ExtendedAt         *time.Time                `json:"extended_at,omitempty"`
	PreviousExpiryDate *time.Time                `json:"previous_expiry_date,omitempty"`
	ExtensionReason    string                    `json:"extension_reason,omitempty"`
	RevokedBy          *string                   `json:"revoked_by,omitempty"`
	RevokedByName      string                    `json:"revoked_by_name,omitempty"`
	RevokedAt          *time.Time                `json:"revoked_at,omitempty"`
	RevocationReason   string                    `json:"revocation_reason,omitempty"`
	CreatedAt          time.Time                 `json:"created_at"`
	UpdatedAt          time.Time                 `json:"updated_at"`
	Logs               []*ClubAcademyLicenseLog  `json:"logs,omitempty"`
}

type ClubAcademyLicenseLog struct {
	ID              string     `json:"id"`
	LicenseID       string     `json:"license_id"`
	Action          string     `json:"action"` // CREATED, EXTENDED, REVOKED, REINSTATED, UPDATED
	PerformedBy     *string    `json:"performed_by,omitempty"`
	PerformedByName string     `json:"performed_by_name"`
	Reason          string     `json:"reason"`
	OldStatus       *string    `json:"old_status,omitempty"`
	NewStatus       *string    `json:"new_status,omitempty"`
	OldExpiryDate   *time.Time `json:"old_expiry_date,omitempty"`
	NewExpiryDate   *time.Time `json:"new_expiry_date,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ClubLicenseKPIs struct {
	TotalLicenses    int `json:"total_licenses"`
	ActiveLicenses   int `json:"active_licenses"`
	ExtendedLicenses int `json:"extended_licenses"`
	RevokedLicenses  int `json:"revoked_licenses"`
	ExpiringSoon     int `json:"expiring_soon"`
}

type CreateClubLicenseRequest struct {
	ClubID        string `json:"club_id"`
	LicenseNumber string `json:"license_number"`
	LicenseType   string `json:"license_type"`
	Category      string `json:"category"`
	IssueDate     string `json:"issue_date"`  // YYYY-MM-DD
	ExpiryDate    string `json:"expiry_date"` // YYYY-MM-DD
	Conditions    string `json:"conditions"`
	DocumentURL   string `json:"document_url"`
}

type ExtendClubLicenseRequest struct {
	NewExpiryDate string `json:"new_expiry_date"` // YYYY-MM-DD
	Reason        string `json:"reason"`
	Notes         string `json:"notes"`
}

type RevokeClubLicenseRequest struct {
	Reason string `json:"reason"`
	Notes  string `json:"notes"`
}

type ReinstateClubLicenseRequest struct {
	Reason string `json:"reason"`
	Notes  string `json:"notes"`
}
