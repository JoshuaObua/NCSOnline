package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/middleware"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountDisabled    = errors.New("account is deactivated")
	ErrTokenInvalid       = errors.New("refresh token is invalid or expired")
)

type AuthService struct {
	users  *repository.UserRepo
	tokens *repository.TokenRepo
	cfg    *config.Config
}

func NewAuthService(users *repository.UserRepo, tokens *repository.TokenRepo, cfg *config.Config) *AuthService {
	return &AuthService{users: users, tokens: tokens, cfg: cfg}
}

type LoginResult struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in_seconds"`
	User         *models.User `json:"user"`
}

func (s *AuthService) Login(ctx context.Context, email, password, ip, ua string) (*LoginResult, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrAccountDisabled
	}

	roles, err := s.users.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get roles: %w", err)
	}
	user.Roles = roles

	roleNames := user.RoleNames()

	accessToken, err := middleware.GenerateAccessToken(s.cfg.JWTSecret, s.cfg.AccessTokenTTL, user.ID, user.Email, roleNames)
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

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    int(s.cfg.AccessTokenTTL.Seconds()),
		User:         user,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, rawToken, ip, ua string) (*LoginResult, error) {
	hash := hashToken(rawToken)
	rt, err := s.tokens.GetByHash(ctx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrTokenInvalid
	}
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	if rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
		return nil, ErrTokenInvalid
	}

	// Rotate: revoke old token
	if err := s.tokens.Revoke(ctx, rt.ID); err != nil {
		return nil, fmt.Errorf("revoke old token: %w", err)
	}

	user, err := s.users.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if !user.IsActive {
		return nil, ErrAccountDisabled
	}

	roles, _ := s.users.GetRoles(ctx, user.ID)
	user.Roles = roles

	accessToken, err := middleware.GenerateAccessToken(s.cfg.JWTSecret, s.cfg.AccessTokenTTL, user.ID, user.Email, user.RoleNames())
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	rawNew, newHash, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	newRT := &models.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: newHash,
		ExpiresAt: time.Now().Add(s.cfg.RefreshTokenTTL),
		IPAddress: ip,
		UserAgent: ua,
	}
	if err := s.tokens.StoreRefreshToken(ctx, newRT); err != nil {
		return nil, fmt.Errorf("store new refresh token: %w", err)
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: rawNew,
		ExpiresIn:    int(s.cfg.AccessTokenTTL.Seconds()),
		User:         user,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	hash := hashToken(rawToken)
	rt, err := s.tokens.GetByHash(ctx, hash)
	if errors.Is(err, repository.ErrNotFound) {
		return nil // already gone — treat as success
	}
	if err != nil {
		return err
	}
	return s.tokens.Revoke(ctx, rt.ID)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, currentPw, newPw string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPw)); err != nil {
		return ErrInvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.users.UpdatePassword(ctx, userID, string(hash)); err != nil {
		return err
	}
	// Revoke all refresh tokens after password change
	return s.tokens.RevokeAllForUser(ctx, userID)
}

func HashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func generateRefreshToken() (raw, hashed string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	raw = hex.EncodeToString(b)
	hashed = hashToken(raw)
	return
}

func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
