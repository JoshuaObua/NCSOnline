package repository

import (
	"context"
	"errors"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/maintenance"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OperatorRepo struct{ db *pgxpool.Pool }

const maintReturning = `RETURNING maintenance_mode,maintenance_reason,maintenance_scheduled_start,maintenance_expected_end,changed_at,COALESCE(changed_by,''),global_auth_invalid_before,cache_generation`

func scanSnapshot(row pgx.Row, s *maintenance.Snapshot) error {
	return row.Scan(&s.Enabled, &s.Reason, &s.ScheduledStart, &s.ExpectedEnd,
		&s.ChangedAt, &s.ChangedBy, &s.GlobalAuthInvalidBefore, &s.CacheGeneration)
}

func (r *OperatorRepo) LoadState(ctx context.Context) (maintenance.Snapshot, error) {
	var s maintenance.Snapshot
	err := scanSnapshot(r.db.QueryRow(ctx,
		`SELECT maintenance_mode,maintenance_reason,maintenance_scheduled_start,maintenance_expected_end,
		         changed_at,COALESCE(changed_by,''),global_auth_invalid_before,cache_generation
		 FROM system_control WHERE singleton=TRUE`), &s)
	if errors.Is(err, pgx.ErrNoRows) {
		return maintenance.Snapshot{ChangedAt: time.Now().UTC()}, nil
	}
	return s, err
}

func (r *OperatorRepo) SaveMaintenance(ctx context.Context, enabled bool, reason string, scheduledStart, expectedEnd *time.Time, actor string) (maintenance.Snapshot, error) {
	var s maintenance.Snapshot
	row := r.db.QueryRow(ctx, `
		INSERT INTO system_control
		  (singleton, maintenance_mode, maintenance_reason, maintenance_scheduled_start, maintenance_expected_end, changed_by, changed_at)
		VALUES (TRUE, $1, $2, $3, $4, $5, NOW())
		ON CONFLICT(singleton) DO UPDATE
		  SET maintenance_mode = EXCLUDED.maintenance_mode,
		      maintenance_reason = EXCLUDED.maintenance_reason,
		      maintenance_scheduled_start = EXCLUDED.maintenance_scheduled_start,
		      maintenance_expected_end = EXCLUDED.maintenance_expected_end,
		      changed_by = EXCLUDED.changed_by,
		      changed_at = NOW()
		` + maintReturning, enabled, reason, scheduledStart, expectedEnd, actor)
	err := scanSnapshot(row, &s)
	return s, err
}

// SetEnabled flips just the maintenance_mode bit (used by the schedule
// sentinel when start/end times cross). Preserves the reason and window.
func (r *OperatorRepo) SetEnabled(ctx context.Context, enabled bool, actor string) (maintenance.Snapshot, error) {
	var s maintenance.Snapshot
	err := scanSnapshot(r.db.QueryRow(ctx, `
		UPDATE system_control
		   SET maintenance_mode = $1, changed_by = $2, changed_at = NOW()
		 WHERE singleton = TRUE
		 ` + maintReturning, enabled, actor), &s)
	return s, err
}
func (r *OperatorRepo) RevokeAllSessions(ctx context.Context, actor string) (maintenance.Snapshot, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return maintenance.Snapshot{}, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `UPDATE refresh_tokens SET revoked_at=COALESCE(revoked_at,NOW()) WHERE revoked_at IS NULL`); err != nil {
		return maintenance.Snapshot{}, err
	}
	var s maintenance.Snapshot
	err = scanSnapshot(tx.QueryRow(ctx, `
		INSERT INTO system_control(singleton, global_auth_invalid_before, changed_by, changed_at)
		VALUES (TRUE, date_trunc('second',NOW()), $1, NOW())
		ON CONFLICT(singleton) DO UPDATE
		  SET global_auth_invalid_before = date_trunc('second',NOW()),
		      changed_by = $1,
		      changed_at = NOW()
		` + maintReturning, actor), &s)
	if err != nil {
		return maintenance.Snapshot{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return maintenance.Snapshot{}, err
	}
	return s, nil
}
func (r *OperatorRepo) SaveCacheGeneration(ctx context.Context, generation uint64, actor string) error {
	_, err := r.db.Exec(ctx, `UPDATE system_control SET cache_generation=$1,changed_by=$2,changed_at=NOW() WHERE singleton=TRUE`, generation, actor)
	return err
}

func (r *OperatorRepo) ListSessions(ctx context.Context, userID string) ([]models.ActiveSession, error) {
	rows, err := r.db.Query(ctx, `SELECT id,COALESCE(ip_address,''),COALESCE(user_agent,''),created_at,expires_at FROM refresh_tokens WHERE user_id=$1 AND revoked_at IS NULL AND expires_at>NOW() ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.ActiveSession{}
	for rows.Next() {
		var s models.ActiveSession
		if err := rows.Scan(&s.ID, &s.IPAddress, &s.UserAgent, &s.CreatedAt, &s.ExpiresAt); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, rows.Err()
}
func (r *OperatorRepo) RevokeSession(ctx context.Context, userID, sessionID string) error {
	tag, err := r.db.Exec(ctx, `UPDATE refresh_tokens SET revoked_at=NOW() WHERE id=$1 AND user_id=$2 AND revoked_at IS NULL`, sessionID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
func (r *OperatorRepo) RevokeOtherSessions(ctx context.Context, userID, currentSessionID string) error {
	_, err := r.db.Exec(ctx, `UPDATE refresh_tokens SET revoked_at=NOW() WHERE user_id=$1 AND id<>$2 AND revoked_at IS NULL`, userID, currentSessionID)
	return err
}
