package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Notification struct {
	ID               string
	TemplateCode     string
	RecipientAddress string
	Payload          json.RawMessage
	AttemptCount     int
}

type NotificationRepo struct{ db *pgxpool.Pool }

func (r *NotificationRepo) ClaimEmail(ctx context.Context) (*Notification, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var n Notification
	err = tx.QueryRow(ctx, `SELECT id,template_code,recipient_address,payload,attempt_count FROM notification_deliveries
		WHERE channel='EMAIL' AND status IN('PENDING','FAILED') AND next_attempt_at<=NOW() AND attempt_count<5
		ORDER BY created_at FOR UPDATE SKIP LOCKED LIMIT 1`).Scan(&n.ID, &n.TemplateCode, &n.RecipientAddress, &n.Payload, &n.AttemptCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if _, err = tx.Exec(ctx, `UPDATE notification_deliveries SET status='SENDING',attempt_count=attempt_count+1,updated_at=NOW() WHERE id=$1`, n.ID); err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NotificationRepo) Delivered(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `UPDATE notification_deliveries SET status='DELIVERED',delivered_at=NOW(),last_error='',updated_at=NOW() WHERE id=$1`, id)
	return err
}
func (r *NotificationRepo) Failed(ctx context.Context, id, message string, attempt int) error {
	status := "FAILED"
	if attempt+1 >= 5 {
		status = "DEAD_LETTER"
	}
	_, err := r.db.Exec(ctx, `UPDATE notification_deliveries SET status=$2,last_error=$3,next_attempt_at=NOW()+(LEAST(attempt_count,5)*INTERVAL '5 minutes'),updated_at=NOW() WHERE id=$1`, id, status, message)
	return err
}
