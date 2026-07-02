package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SecurityRepo struct{ db *pgxpool.Pool }

// ── IP whitelist ─────────────────────────────────────────────────

type IPAllowEntry struct {
	ID        string    `json:"id"`
	IPOrCIDR  string    `json:"ip_or_cidr"`
	Label     string    `json:"label"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *SecurityRepo) ListWhitelist(ctx context.Context, userID string) ([]*IPAllowEntry, error) {
	const q = `SELECT id, ip_or_cidr, label, is_active, created_at
	           FROM user_ip_whitelist
	           WHERE user_id=$1
	           ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*IPAllowEntry{}
	for rows.Next() {
		e := &IPAllowEntry{}
		if err := rows.Scan(&e.ID, &e.IPOrCIDR, &e.Label, &e.IsActive, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *SecurityRepo) ActiveWhitelistFor(ctx context.Context, userID string) ([]string, error) {
	const q = `SELECT ip_or_cidr FROM user_ip_whitelist WHERE user_id=$1 AND is_active=TRUE`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r *SecurityRepo) AddWhitelist(ctx context.Context, userID, ipOrCIDR, label string) (*IPAllowEntry, error) {
	const q = `INSERT INTO user_ip_whitelist (user_id, ip_or_cidr, label)
	           VALUES ($1,$2,$3)
	           ON CONFLICT (user_id, ip_or_cidr) DO UPDATE
	             SET label=EXCLUDED.label, is_active=TRUE
	           RETURNING id, ip_or_cidr, label, is_active, created_at`
	e := &IPAllowEntry{}
	err := r.db.QueryRow(ctx, q, userID, ipOrCIDR, label).
		Scan(&e.ID, &e.IPOrCIDR, &e.Label, &e.IsActive, &e.CreatedAt)
	return e, err
}

func (r *SecurityRepo) DeleteWhitelist(ctx context.Context, userID, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM user_ip_whitelist WHERE user_id=$1 AND id=$2`, userID, id)
	return err
}

// ── 2FA ──────────────────────────────────────────────────────────

type TwoFAState struct {
	Enabled    bool       `json:"enabled"`
	EnabledAt  *time.Time `json:"enabled_at,omitempty"`
	HasPending bool       `json:"has_pending"` // secret set but not enabled
}

func (r *SecurityRepo) GetTwoFAState(ctx context.Context, userID string) (*TwoFAState, error) {
	var secret *string
	var enabled bool
	var enabledAt *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT twofa_secret, twofa_enabled, twofa_enabled_at FROM users WHERE id=$1`, userID).
		Scan(&secret, &enabled, &enabledAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &TwoFAState{
		Enabled:    enabled,
		EnabledAt:  enabledAt,
		HasPending: secret != nil && *secret != "" && !enabled,
	}, nil
}

func (r *SecurityRepo) GetTwoFASecret(ctx context.Context, userID string) (string, bool, error) {
	var secret *string
	var enabled bool
	err := r.db.QueryRow(ctx,
		`SELECT twofa_secret, twofa_enabled FROM users WHERE id=$1`, userID).
		Scan(&secret, &enabled)
	if err != nil {
		return "", false, err
	}
	if secret == nil {
		return "", enabled, nil
	}
	return *secret, enabled, nil
}

func (r *SecurityRepo) SetTwoFASecret(ctx context.Context, userID, secret string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET twofa_secret=$2, twofa_enabled=FALSE, twofa_enabled_at=NULL WHERE id=$1`,
		userID, secret)
	return err
}

func (r *SecurityRepo) EnableTwoFA(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET twofa_enabled=TRUE, twofa_enabled_at=NOW() WHERE id=$1`, userID)
	return err
}

func (r *SecurityRepo) DisableTwoFA(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET twofa_secret=NULL, twofa_enabled=FALSE, twofa_enabled_at=NULL WHERE id=$1`,
		userID)
	return err
}
