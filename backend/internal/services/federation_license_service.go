package services

import (
	"context"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
)

type FederationLicenseService struct {
	repo *repository.FederationLicenseRepo
}

func NewFederationLicenseService(repo *repository.FederationLicenseRepo) *FederationLicenseService {
	return &FederationLicenseService{repo: repo}
}

func (s *FederationLicenseService) CreateLicense(ctx context.Context, lic *models.FederationLicense, actorID, actorName string) error {
	return s.repo.Create(ctx, lic, actorID, actorName)
}

func (s *FederationLicenseService) GetLicenseByID(ctx context.Context, id string) (*models.FederationLicense, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FederationLicenseService) GetActiveLicenseByFederationID(ctx context.Context, fedID string) (*models.FederationLicense, error) {
	return s.repo.GetActiveLicenseByFederationID(ctx, fedID)
}

func (s *FederationLicenseService) ListLicenses(ctx context.Context, f repository.ListLicensesFilter, p *models.PaginationParams) ([]*models.FederationLicense, int, error) {
	return s.repo.List(ctx, f, p)
}

func (s *FederationLicenseService) ExtendLicense(ctx context.Context, id string, newExpiryDate time.Time, reason, actorID, actorName string) (*models.FederationLicense, error) {
	return s.repo.Extend(ctx, id, newExpiryDate, reason, actorID, actorName)
}

func (s *FederationLicenseService) RevokeLicense(ctx context.Context, id, reason, actorID, actorName string) (*models.FederationLicense, error) {
	return s.repo.Revoke(ctx, id, reason, actorID, actorName)
}

func (s *FederationLicenseService) ReinstateLicense(ctx context.Context, id, reason, actorID, actorName string) (*models.FederationLicense, error) {
	return s.repo.Reinstate(ctx, id, reason, actorID, actorName)
}

func (s *FederationLicenseService) GetKPIs(ctx context.Context) (*models.LicenseKPIs, error) {
	return s.repo.GetKPIs(ctx)
}
