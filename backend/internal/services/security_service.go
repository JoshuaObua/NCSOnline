package services

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
)

var (
	ErrTwoFAInvalid   = errors.New("invalid 2FA code")
	ErrTwoFANotSetup  = errors.New("2FA is not set up for this account")
	ErrInvalidIPEntry = errors.New("invalid ip address or CIDR")
)

type SecurityService struct {
	repo  *repository.SecurityRepo
	users *repository.UserRepo
}

func NewSecurityService(repo *repository.SecurityRepo, users *repository.UserRepo) *SecurityService {
	return &SecurityService{repo: repo, users: users}
}

// ── IP whitelist ─────────────────────────────────────────────────

func (s *SecurityService) ListWhitelist(ctx context.Context, userID string) ([]*repository.IPAllowEntry, error) {
	return s.repo.ListWhitelist(ctx, userID)
}

func (s *SecurityService) AddWhitelist(ctx context.Context, userID, ipOrCIDR, label string) (*repository.IPAllowEntry, error) {
	v := strings.TrimSpace(ipOrCIDR)
	if v == "" {
		return nil, ErrInvalidIPEntry
	}
	if strings.Contains(v, "/") {
		if _, _, err := net.ParseCIDR(v); err != nil {
			return nil, ErrInvalidIPEntry
		}
	} else if net.ParseIP(v) == nil {
		return nil, ErrInvalidIPEntry
	}
	return s.repo.AddWhitelist(ctx, userID, v, strings.TrimSpace(label))
}

func (s *SecurityService) RemoveWhitelist(ctx context.Context, userID, id string) error {
	return s.repo.DeleteWhitelist(ctx, userID, id)
}

// ── 2FA ──────────────────────────────────────────────────────────

type TwoFAEnrollment struct {
	Secret    string `json:"secret"`     // Base32, no padding — shown to user
	OTPAuth   string `json:"otpauth_url"` // otpauth:// URI suitable for QR encoding
	Issuer    string `json:"issuer"`
	AccountID string `json:"account_id"`
}

func (s *SecurityService) GetTwoFAState(ctx context.Context, userID string) (*repository.TwoFAState, error) {
	return s.repo.GetTwoFAState(ctx, userID)
}

// BeginEnrollment generates a new TOTP secret, stores it pending, and returns
// the secret + otpauth URI. The user scans the URI in their authenticator app
// and then calls Verify with a code to activate 2FA.
func (s *SecurityService) BeginEnrollment(ctx context.Context, userID, accountEmail string) (*TwoFAEnrollment, error) {
	secret, err := randomBase32(20) // 160-bit secret per RFC 6238 recommendation
	if err != nil {
		return nil, err
	}
	if err := s.repo.SetTwoFASecret(ctx, userID, secret); err != nil {
		return nil, err
	}
	issuer := "NCS Uganda"
	label := url.PathEscape(fmt.Sprintf("%s:%s", issuer, accountEmail))
	q := url.Values{}
	q.Set("secret", secret)
	q.Set("issuer", issuer)
	q.Set("algorithm", "SHA1")
	q.Set("digits", "6")
	q.Set("period", "30")
	return &TwoFAEnrollment{
		Secret:    secret,
		OTPAuth:   fmt.Sprintf("otpauth://totp/%s?%s", label, q.Encode()),
		Issuer:    issuer,
		AccountID: accountEmail,
	}, nil
}

func (s *SecurityService) VerifyAndEnable(ctx context.Context, userID, code string) error {
	secret, _, err := s.repo.GetTwoFASecret(ctx, userID)
	if err != nil {
		return err
	}
	if secret == "" {
		return ErrTwoFANotSetup
	}
	if !verifyTOTP(secret, code, time.Now().Unix()) {
		return ErrTwoFAInvalid
	}
	return s.repo.EnableTwoFA(ctx, userID)
}

func (s *SecurityService) Disable(ctx context.Context, userID string) error {
	return s.repo.DisableTwoFA(ctx, userID)
}

func (s *SecurityService) SetSecretDirectly(ctx context.Context, userID, secret string) error {
	return s.repo.SetTwoFASecret(ctx, userID, secret)
}

func (s *SecurityService) GetSecretDirectly(ctx context.Context, userID string) (string, bool, error) {
	return s.repo.GetTwoFASecret(ctx, userID)
}

func (s *SecurityService) EnableDirectly(ctx context.Context, userID string) error {
	return s.repo.EnableTwoFA(ctx, userID)
}

// VerifyCode returns whether the supplied code is valid for the user's
// enabled 2FA secret (used by the login flow once we wire it up).
func (s *SecurityService) VerifyCode(ctx context.Context, userID, code string) (bool, error) {
	secret, enabled, err := s.repo.GetTwoFASecret(ctx, userID)
	if err != nil || !enabled || secret == "" {
		return false, err
	}
	return verifyTOTP(secret, code, time.Now().Unix()), nil
}

// ── TOTP (RFC 6238, SHA1, 30s, 6 digits) ─────────────────────────

func randomBase32(nBytes int) (string, error) {
	buf := make([]byte, nBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	return enc.EncodeToString(buf), nil
}

func generateTOTP(secretBase32 string, unixTime int64) string {
	dec := base32.StdEncoding.WithPadding(base32.NoPadding)
	// Tolerate accidental padding/casing from the user.
	clean := strings.ToUpper(strings.TrimRight(secretBase32, "="))
	key, err := dec.DecodeString(clean)
	if err != nil {
		return ""
	}
	counter := uint64(unixTime / 30)
	cb := make([]byte, 8)
	binary.BigEndian.PutUint64(cb, counter)
	h := hmac.New(sha1.New, key)
	h.Write(cb)
	sum := h.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	bin := (uint32(sum[offset])&0x7f)<<24 |
		uint32(sum[offset+1])<<16 |
		uint32(sum[offset+2])<<8 |
		uint32(sum[offset+3])
	return fmt.Sprintf("%06d", bin%1_000_000)
}

func verifyTOTP(secret, code string, now int64) bool {
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return false
	}
	// Allow ±1 step (30s) drift to be friendly to clock skew.
	for _, drift := range []int64{-30, 0, 30} {
		if generateTOTP(secret, now+drift) == code {
			return true
		}
	}
	return false
}
