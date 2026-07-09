package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
)

type handlerUSSDData struct{}

func (handlerUSSDData) FindActiveUserByPhone(context.Context, string) (*models.User, error) {
	return nil, repository.ErrNotFound
}
func (handlerUSSDData) FindAthlete(context.Context, string) (*repository.USSDAthlete, error) {
	return nil, repository.ErrNotFound
}
func (handlerUSSDData) FindLicence(context.Context, string) (*repository.USSDLicence, error) {
	return nil, repository.ErrNotFound
}

type handlerUSSDApps struct{}

func (handlerUSSDApps) ListByUser(context.Context, string, *models.PaginationParams) ([]*models.Application, int64, error) {
	return nil, 0, nil
}

func TestUSSDHandlerAuthenticatesCallbackAndRendersMenu(t *testing.T) {
	h := &USSDHandler{service: services.NewUSSDService(handlerUSSDData{}, handlerUSSDApps{}), callbackSecret: "secret"}
	unauthorized := httptest.NewRecorder()
	h.ServeHTTP(unauthorized, httptest.NewRequest(http.MethodPost, "/ncs-ussd", strings.NewReader("text=")))
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", unauthorized.Code)
	}
	req := httptest.NewRequest(http.MethodPost, "/ncs-ussd", strings.NewReader("sessionId=s1&serviceCode=%2A123%23&phoneNumber=%2B256772123456&text="))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-USSD-Token", "secret")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || !strings.HasPrefix(rr.Body.String(), "CON NCS Quick Services") {
		t.Fatalf("unexpected response %d %q", rr.Code, rr.Body.String())
	}
	if rr.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("USSD response must not be cached")
	}
}
