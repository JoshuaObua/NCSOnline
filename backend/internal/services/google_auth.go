package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/middleware"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrGoogleAuthNotConfigured = errors.New("google auth is not configured")
	ErrGoogleTokenInvalid      = errors.New("google token is invalid")
)

type googleTokenInfo struct {
	Issuer        string `json:"iss"`
	Audience      string `json:"aud"`
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	ExpiresIn     string `json:"expires_in"`
}

func (s *AuthService) LoginWithGoogle(ctx context.Context, credential, ip, ua string) (*LoginResult, error) {
	profile, err := s.verifyGoogleCredential(ctx, credential)
	if err != nil {
		return nil, err
	}

	user, err := s.users.GetByGoogleSub(ctx, profile.Subject)
	if errors.Is(err, repository.ErrNotFound) {
		user, err = s.users.GetByEmail(ctx, profile.Email)
		if errors.Is(err, repository.ErrNotFound) {
			user, err = s.createGoogleSubscriber(ctx, profile)
		} else if err == nil {
			err = s.users.LinkGoogleIdentity(ctx, user.ID, profile.Subject, profile.Picture)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("resolve google user: %w", err)
	}
	if !user.IsActive {
		return nil, ErrAccountDisabled
	}

	roles, err := s.users.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get roles: %w", err)
	}
	if len(roles) == 0 {
		if role, roleErr := s.users.GetRoleByName(ctx, "subscriber"); roleErr == nil {
			_ = s.users.AssignRole(ctx, user.ID, role.ID, user.ID)
			roles, _ = s.users.GetRoles(ctx, user.ID)
		}
	}
	user.Roles = roles

	return s.issueSession(ctx, user, ip, ua)
}

func (s *AuthService) createGoogleSubscriber(ctx context.Context, profile *googleTokenInfo) (*models.User, error) {
	first, last := splitGoogleName(profile.FirstName, profile.LastName, profile.Name)
	hash, err := randomPasswordHash()
	if err != nil {
		return nil, err
	}
	u := &models.User{
		ID:              uuid.NewString(),
		Email:           profile.Email,
		PasswordHash:    hash,
		FirstName:       first,
		LastName:        last,
		AvatarURL:       profile.Picture,
		AuthProvider:    "google",
		GoogleSub:       profile.Subject,
		IsActive:        true,
		IsEmailVerified: true,
	}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	if role, err := s.users.GetRoleByName(ctx, "subscriber"); err == nil {
		_ = s.users.AssignRole(ctx, u.ID, role.ID, u.ID)
	}
	return u, nil
}

func (s *AuthService) issueSession(ctx context.Context, user *models.User, ip, ua string) (*LoginResult, error) {
	accessToken, err := middleware.GenerateAccessToken(s.cfg.JWTSecret, s.cfg.AccessTokenTTL, user.ID, user.Email, user.RoleNames())
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}
	rawRefresh, tokenHash, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}
	rt := &models.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.cfg.RefreshTokenTTL),
		IPAddress: ip,
		UserAgent: ua,
	}
	if err := s.tokens.StoreRefreshToken(ctx, rt); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}
	_ = s.users.UpdateLastLogin(ctx, user.ID)
	user.PasswordHash = ""
	user.PinHash = ""
	return &LoginResult{
		AccessToken:       accessToken,
		RefreshToken:      rawRefresh,
		ExpiresIn:         int(s.cfg.AccessTokenTTL.Seconds()),
		PinChangeRequired: user.PinChangeRequired,
		User:              user,
	}, nil
}

func (s *AuthService) verifyGoogleCredential(ctx context.Context, credential string) (*googleTokenInfo, error) {
	credential = strings.TrimSpace(credential)
	if s.cfg.GoogleClientID == "" {
		return nil, ErrGoogleAuthNotConfigured
	}
	if credential == "" || len(credential) > 8192 {
		return nil, ErrGoogleTokenInvalid
	}

	endpoint := s.cfg.GoogleTokenInfoURL
	if endpoint == "" {
		endpoint = "https://oauth2.googleapis.com/tokeninfo"
	}
	reqURL := endpoint + "?id_token=" + url.QueryEscape(credential)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("verify google token: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, ErrGoogleTokenInvalid
	}
	var info googleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, ErrGoogleTokenInvalid
	}
	if info.Audience != s.cfg.GoogleClientID ||
		(info.Issuer != "accounts.google.com" && info.Issuer != "https://accounts.google.com") ||
		info.Subject == "" ||
		info.Email == "" ||
		!strings.EqualFold(info.EmailVerified, "true") {
		return nil, ErrGoogleTokenInvalid
	}
	info.Email = strings.ToLower(strings.TrimSpace(info.Email))
	return &info, nil
}

func splitGoogleName(first, last, full string) (string, string) {
	first = strings.TrimSpace(first)
	last = strings.TrimSpace(last)
	if first != "" || last != "" {
		if first == "" {
			first = "Google"
		}
		if last == "" {
			last = "User"
		}
		return first, last
	}
	parts := strings.Fields(full)
	if len(parts) == 0 {
		return "Google", "User"
	}
	if len(parts) == 1 {
		return parts[0], "User"
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func randomPasswordHash() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return HashPassword("google:" + hex.EncodeToString(b))
}
