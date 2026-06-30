package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type USSDData interface {
	FindActiveUserByPhone(context.Context, string) (*models.User, error)
	FindAthlete(context.Context, string) (*repository.USSDAthlete, error)
	FindLicence(context.Context, string) (*repository.USSDLicence, error)
}

type USSDApplications interface {
	ListByUser(context.Context, string, *models.PaginationParams) ([]*models.Application, int64, error)
}

type USSDService struct {
	data         USSDData
	applications USSDApplications
}

func NewUSSDService(data USSDData, applications USSDApplications) *USSDService {
	return &USSDService{data: data, applications: applications}
}

type USSDInput struct {
	SessionID   string `json:"sessionId"`
	ServiceCode string `json:"serviceCode"`
	PhoneNumber string `json:"phoneNumber"`
	Text        string `json:"text"`
}
type USSDOutput struct {
	Continue bool
	Message  string
}

func con(message string) USSDOutput { return USSDOutput{Continue: true, Message: message} }
func end(message string) USSDOutput { return USSDOutput{Message: message} }

func (s *USSDService) Handle(ctx context.Context, input USSDInput) USSDOutput {
	parts := splitUSSD(input.Text)
	if len(parts) == 0 {
		return con("NCS Quick Services\n1. Check athlete registration\n2. Check licence validity\n3. My applications\n4. Renewal information\n5. Help")
	}
	switch parts[0] {
	case "1":
		return s.athlete(ctx, parts)
	case "2":
		return s.licence(ctx, parts)
	case "3":
		return s.applicationsFlow(ctx, parts, input.PhoneNumber)
	case "4":
		return s.renewal(ctx, parts, input.PhoneNumber)
	case "5":
		return end("NCS Help\nCall: +256 414 254477\nEmail: info@ncs.go.ug\nWeb: www.ncs.go.ug")
	default:
		return end("Invalid selection. Please dial the NCS code and try again.")
	}
}

func (s *USSDService) athlete(ctx context.Context, parts []string) USSDOutput {
	if len(parts) == 1 {
		return con("Enter the NCS athlete number:")
	}
	a, err := s.data.FindAthlete(ctx, parts[1])
	if errors.Is(err, repository.ErrNotFound) {
		return end("No matching athlete registration was found.")
	}
	if err != nil {
		return end("Service temporarily unavailable. Please try again.")
	}
	verified := "Not yet verified"
	if a.VerifiedAt != nil {
		verified = "Verified " + a.VerifiedAt.Format("02 Jan 2006")
	}
	return end(fmt.Sprintf("ATHLETE REGISTERED\n%s\n%s | %s\nFederation: %s\n%s", maskName(a.FullName), a.Discipline, a.Status, emptyAs(a.Federation, "Not recorded"), verified))
}
func (s *USSDService) licence(ctx context.Context, parts []string) USSDOutput {
	if len(parts) == 1 {
		return con("Enter the NCS licence or credential number:")
	}
	l, err := s.data.FindLicence(ctx, parts[1])
	if errors.Is(err, repository.ErrNotFound) {
		return end("No matching licence or credential was found.")
	}
	if err != nil {
		return end("Service temporarily unavailable. Please try again.")
	}
	expires := "No expiry recorded"
	if l.ExpiresOn != nil {
		expires = "Expires: " + l.ExpiresOn.Format("02 Jan 2006")
	}
	return end(fmt.Sprintf("%s\n%s\n%s\nHolder: %s\n%s", l.Status, l.Number, l.Type, maskName(l.HolderName), expires))
}
func (s *USSDService) applicationsFlow(ctx context.Context, parts []string, callerPhone string) USSDOutput {
	if len(parts) == 1 {
		return con("Enter your registered phone number:")
	}
	if len(parts) == 2 {
		return con("Enter your 4-6 digit NCS PIN:")
	}
	u, ok := s.authenticate(ctx, parts[1], parts[2], callerPhone)
	if !ok {
		return end("Phone number or PIN is incorrect. Access denied.")
	}
	apps, _, err := s.applications.ListByUser(ctx, u.ID, &models.PaginationParams{Page: 1, PerPage: 3})
	if err != nil {
		return end("Service temporarily unavailable. Please try again.")
	}
	if len(apps) == 0 {
		return end("You have no applications on your NCS account.")
	}
	lines := []string{"YOUR NCS APPLICATIONS"}
	for _, a := range apps {
		ref := a.ApplicationReference
		if ref == "" {
			ref = "Draft " + strings.ToUpper(a.FormType)
		}
		lines = append(lines, fmt.Sprintf("%s: %s (%s)", ref, a.Status, a.PaymentStatus))
	}
	return end(strings.Join(lines, "\n"))
}
func (s *USSDService) renewal(ctx context.Context, parts []string, callerPhone string) USSDOutput {
	if len(parts) == 1 {
		return con("Enter your registered phone number:")
	}
	if len(parts) == 2 {
		return con("Enter your 4-6 digit NCS PIN:")
	}
	_, ok := s.authenticate(ctx, parts[1], parts[2], callerPhone)
	if !ok {
		return end("Phone number or PIN is incorrect. Access denied.")
	}
	if len(parts) == 3 {
		return con("Enter the licence number to check renewal:")
	}
	l, err := s.data.FindLicence(ctx, parts[3])
	if errors.Is(err, repository.ErrNotFound) {
		return end("No matching licence was found on the NCS register.")
	}
	if err != nil {
		return end("Service temporarily unavailable. Please try again.")
	}
	if !l.Renewable || l.Status == "SUSPENDED" || l.Status == "REVOKED" || l.Status == "CANCELLED" || l.Status == "SUPERSEDED" {
		return end("This licence is not currently eligible for renewal. Contact NCS support.")
	}
	expiry := ""
	if l.ExpiresOn != nil {
		expiry = "\nExpiry: " + l.ExpiresOn.Format("02 Jan 2006")
	}
	return end("RENEWAL AVAILABLE\n" + l.Number + expiry + "\nSign in at the NCS website and select Renew, or contact NCS for assistance.")
}

func (s *USSDService) authenticate(ctx context.Context, phone, pin, callerPhone string) (*models.User, bool) {
	normalized, ok := normalizeUGPhone(phone)
	if !ok || !regexp.MustCompile(`^\d{4,6}$`).MatchString(pin) {
		return nil, false
	}
	if strings.TrimSpace(callerPhone) != "" {
		caller, valid := normalizeUGPhone(callerPhone)
		if !valid || caller != normalized {
			return nil, false
		}
	}
	u, err := s.data.FindActiveUserByPhone(ctx, normalized)
	if err != nil || u.PinHash == "" || u.PinChangeRequired {
		return nil, false
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PinHash), []byte(pin)) != nil {
		return nil, false
	}
	return u, true
}
func splitUSSD(text string) []string {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	raw := strings.Split(text, "*")
	out := make([]string, 0, len(raw))
	for _, part := range raw {
		out = append(out, strings.TrimSpace(part))
	}
	return out
}
func normalizeUGPhone(phone string) (string, bool) {
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	if strings.HasPrefix(digits, "0") && len(digits) == 10 {
		digits = "256" + digits[1:]
	}
	if len(digits) != 12 || !strings.HasPrefix(digits, "256") {
		return "", false
	}
	return "+" + digits, true
}
func maskName(name string) string {
	words := strings.Fields(name)
	for i, w := range words {
		r := []rune(w)
		if len(r) > 1 {
			words[i] = string(r[0]) + strings.Repeat("*", len(r)-1)
		}
	}
	return strings.Join(words, " ")
}
func emptyAs(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
