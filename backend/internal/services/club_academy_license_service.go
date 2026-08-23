package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
)

type ClubAcademyLicenseService struct {
	repo *repository.ClubAcademyLicenseRepo
}

func NewClubAcademyLicenseService(repo *repository.ClubAcademyLicenseRepo) *ClubAcademyLicenseService {
	return &ClubAcademyLicenseService{repo: repo}
}

func (s *ClubAcademyLicenseService) ListLicenses(ctx context.Context, f repository.ListClubLicensesFilter, p *models.PaginationParams) ([]*models.ClubAcademyLicense, int, error) {
	return s.repo.List(ctx, f, p)
}

func (s *ClubAcademyLicenseService) GetLicenseByID(ctx context.Context, id string) (*models.ClubAcademyLicense, error) {
	if id == "" {
		return nil, errors.New("license id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ClubAcademyLicenseService) GetActiveLicenseByClubID(ctx context.Context, clubID string) (*models.ClubAcademyLicense, error) {
	if clubID == "" {
		return nil, errors.New("club id is required")
	}
	return s.repo.GetActiveLicenseByClubID(ctx, clubID)
}

func (s *ClubAcademyLicenseService) GetKPIs(ctx context.Context) (*models.ClubLicenseKPIs, error) {
	return s.repo.GetKPIs(ctx)
}

func (s *ClubAcademyLicenseService) CreateLicense(ctx context.Context, req models.CreateClubLicenseRequest, actorID, actorName string) (*models.ClubAcademyLicense, error) {
	if req.ClubID == "" {
		return nil, errors.New("club_id is required")
	}
	if req.LicenseNumber == "" {
		return nil, errors.New("license_number is required")
	}
	if req.ExpiryDate == "" {
		return nil, errors.New("expiry_date is required")
	}

	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return nil, fmt.Errorf("invalid expiry_date format: %w", err)
	}

	var issueDate time.Time
	if req.IssueDate != "" {
		issueDate, err = time.Parse("2006-01-02", req.IssueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid issue_date format: %w", err)
		}
	} else {
		issueDate = time.Now()
	}

	licType := req.LicenseType
	if licType == "" {
		licType = "STATUTORY_RECOGNITION"
	}
	category := req.Category
	if category == "" {
		category = "SPORTS_ACADEMY"
	}

	lic := &models.ClubAcademyLicense{
		ClubID:        req.ClubID,
		LicenseNumber: req.LicenseNumber,
		LicenseType:   licType,
		Category:      category,
		IssueDate:     issueDate,
		ExpiryDate:    expiryDate,
		Status:        "ACTIVE",
		Conditions:    req.Conditions,
		DocumentURL:   req.DocumentURL,
	}

	if err := s.repo.Create(ctx, lic, actorID, actorName); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, lic.ID)
}

func (s *ClubAcademyLicenseService) ExtendLicense(ctx context.Context, id string, req models.ExtendClubLicenseRequest, actorID, actorName string) (*models.ClubAcademyLicense, error) {
	if id == "" {
		return nil, errors.New("license id is required")
	}
	if req.NewExpiryDate == "" {
		return nil, errors.New("new_expiry_date is required")
	}
	if req.Reason == "" {
		return nil, errors.New("reason for extension is required")
	}

	newExpiry, err := time.Parse("2006-01-02", req.NewExpiryDate)
	if err != nil {
		return nil, fmt.Errorf("invalid new_expiry_date format: %w", err)
	}

	return s.repo.Extend(ctx, id, newExpiry, req.Reason, req.Notes, actorID, actorName)
}

func (s *ClubAcademyLicenseService) RevokeLicense(ctx context.Context, id string, req models.RevokeClubLicenseRequest, actorID, actorName string) (*models.ClubAcademyLicense, error) {
	if id == "" {
		return nil, errors.New("license id is required")
	}
	if req.Reason == "" {
		return nil, errors.New("revocation reason is required")
	}

	return s.repo.Revoke(ctx, id, req.Reason, req.Notes, actorID, actorName)
}

func (s *ClubAcademyLicenseService) ReinstateLicense(ctx context.Context, id string, req models.ReinstateClubLicenseRequest, actorID, actorName string) (*models.ClubAcademyLicense, error) {
	if id == "" {
		return nil, errors.New("license id is required")
	}
	if req.Reason == "" {
		return nil, errors.New("reinstatement reason is required")
	}

	return s.repo.Reinstate(ctx, id, req.Reason, req.Notes, actorID, actorName)
}
