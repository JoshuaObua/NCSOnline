package services

import (
	"context"
	"testing"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type fakeUSSDData struct {
	user    *models.User
	athlete *repository.USSDAthlete
	licence *repository.USSDLicence
}

func (f *fakeUSSDData) FindActiveUserByPhone(context.Context, string) (*models.User, error) {
	if f.user == nil {
		return nil, repository.ErrNotFound
	}
	return f.user, nil
}
func (f *fakeUSSDData) FindAthlete(context.Context, string) (*repository.USSDAthlete, error) {
	if f.athlete == nil {
		return nil, repository.ErrNotFound
	}
	return f.athlete, nil
}
func (f *fakeUSSDData) FindLicence(context.Context, string) (*repository.USSDLicence, error) {
	if f.licence == nil {
		return nil, repository.ErrNotFound
	}
	return f.licence, nil
}

type fakeUSSDApps struct {
	apps  []*models.Application
	calls int
}

func (f *fakeUSSDApps) ListByUser(context.Context, string, *models.PaginationParams) ([]*models.Application, int64, error) {
	f.calls++
	return f.apps, int64(len(f.apps)), nil
}

func TestUSSDPublicAthleteLookupDoesNotRequirePIN(t *testing.T) {
	verified := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	data := &fakeUSSDData{athlete: &repository.USSDAthlete{Number: "NCS-1", FullName: "Jane Doe", Status: "ACTIVE", Discipline: "Athletics", Federation: "Athletics Uganda", VerifiedAt: &verified}}
	out := NewUSSDService(data, &fakeUSSDApps{}).Handle(context.Background(), USSDInput{Text: "1*NCS-1"})
	if out.Continue || out.Message == "" {
		t.Fatalf("expected terminal public result: %#v", out)
	}
}

func TestUSSDProtectedApplicationsRequirePhoneAndPIN(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.MinCost)
	apps := &fakeUSSDApps{apps: []*models.Application{{ApplicationReference: "NCS-2026-F3-ABC", Status: "UNDER_REVIEW", PaymentStatus: "PAYMENT_VERIFIED"}}}
	data := &fakeUSSDData{user: &models.User{ID: "u1", PinHash: string(hash), IsActive: true, AccountStatus: "ACTIVE"}}
	service := NewUSSDService(data, apps)
	denied := service.Handle(context.Background(), USSDInput{PhoneNumber: "+256772123456", Text: "3*0772123456*9999"})
	if denied.Continue || apps.calls != 0 {
		t.Fatalf("wrong PIN must not access applications: %#v calls=%d", denied, apps.calls)
	}
	allowed := service.Handle(context.Background(), USSDInput{PhoneNumber: "+256772123456", Text: "3*0772123456*1234"})
	if allowed.Continue || apps.calls != 1 {
		t.Fatalf("valid phone and PIN should access applications: %#v calls=%d", allowed, apps.calls)
	}
	mismatch := service.Handle(context.Background(), USSDInput{PhoneNumber: "+256701000000", Text: "3*0772123456*1234"})
	if mismatch.Continue || apps.calls != 1 {
		t.Fatalf("caller MSISDN mismatch must be denied: %#v calls=%d", mismatch, apps.calls)
	}
}

func TestNormalizeUGPhone(t *testing.T) {
	for _, input := range []string{"0772 123456", "+256772123456", "256772123456"} {
		got, ok := normalizeUGPhone(input)
		if !ok || got != "+256772123456" {
			t.Fatalf("normalize %q = %q,%v", input, got, ok)
		}
	}
}
