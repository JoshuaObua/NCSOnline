package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) *TransactionRepo {
	return &TransactionRepo{db: db}
}

type ListTransactionsFilter struct {
	UserID        string
	SubmissionID  string
	Status        string
	PaymentMethod string
	Search        string
	StartDate     string
	EndDate       string
}

type TransactionKPIs struct {
	TotalAmountUGX float64 `json:"total_amount_ugx"`
	TotalCount     int     `json:"total_count"`
	SuccessCount   int     `json:"success_count"`
	PendingCount   int     `json:"pending_count"`
	FailedCount    int     `json:"failed_count"`
}

func (r *TransactionRepo) Create(ctx context.Context, tx *models.PaymentTransaction) error {
	if tx.ID == "" {
		tx.ID = uuid.NewString()
	}
	if tx.Currency == "" {
		tx.Currency = "UGX"
	}
	if tx.Status == "" {
		tx.Status = models.TxStatusPending
	}
	if len(tx.RawResponse) == 0 {
		tx.RawResponse = json.RawMessage("{}")
	}

	const q = `
		INSERT INTO payment_transactions (
			id, transaction_reference, submission_id, template_id, user_id,
			payment_method, provider, provider_request_id, phone_number,
			amount_ugx, currency, status, status_message, raw_response,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13, $14::jsonb,
			NOW(), NOW()
		)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRow(ctx, q,
		tx.ID, tx.TransactionReference, tx.SubmissionID, tx.TemplateID, tx.UserID,
		tx.PaymentMethod, tx.Provider, tx.ProviderRequestID, tx.PhoneNumber,
		tx.AmountUGX, tx.Currency, tx.Status, tx.StatusMessage, string(tx.RawResponse),
	).Scan(&tx.CreatedAt, &tx.UpdatedAt)
}

func (r *TransactionRepo) UpdateStatus(ctx context.Context, id, status, statusMsg string, rawResp []byte) error {
	return r.UpdateStatusWithProvider(ctx, id, "", status, statusMsg, rawResp)
}

func (r *TransactionRepo) UpdateStatusWithProvider(ctx context.Context, id, providerReqID, status, statusMsg string, rawResp []byte) error {
	if len(rawResp) == 0 {
		rawResp = []byte("{}")
	}
	var completedAt *time.Time
	if status == models.TxStatusSuccess || status == models.TxStatusFailed || status == models.TxStatusCancelled {
		now := time.Now()
		completedAt = &now
	}

	const q = `
		UPDATE payment_transactions
		SET status = $2, status_message = $3, raw_response = $4::jsonb,
		    provider_request_id = COALESCE(NULLIF($5, ''), provider_request_id),
		    updated_at = NOW(), completed_at = COALESCE($6, completed_at)
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, q, id, status, statusMsg, string(rawResp), providerReqID, completedAt)
	return err
}

func (r *TransactionRepo) GetByID(ctx context.Context, id string) (*models.PaymentTransaction, error) {
	const q = `
		SELECT pt.id, pt.transaction_reference, pt.submission_id,
		       COALESCE(fs.submission_reference, ''), pt.template_id,
		       COALESCE(ft.title, ''), pt.user_id,
		       COALESCE(u.first_name || ' ' || u.last_name, ''),
		       COALESCE(u.email, ''), pt.payment_method, pt.provider,
		       pt.provider_request_id, COALESCE(pt.phone_number, ''),
		       pt.amount_ugx, pt.currency, pt.status, COALESCE(pt.status_message, ''),
		       pt.raw_response, pt.created_at, pt.updated_at, pt.completed_at
		FROM payment_transactions pt
		LEFT JOIN form_submissions fs ON fs.id = pt.submission_id
		LEFT JOIN form_templates ft ON ft.id = pt.template_id
		LEFT JOIN users u ON u.id = pt.user_id
		WHERE pt.id = $1
	`
	row := r.db.QueryRow(ctx, q, id)
	return r.scanTransaction(row)
}

func (r *TransactionRepo) GetLatestBySubmissionID(ctx context.Context, submissionID string) (*models.PaymentTransaction, error) {
	const q = `
		SELECT pt.id, pt.transaction_reference, pt.submission_id,
		       COALESCE(fs.submission_reference, ''), pt.template_id,
		       COALESCE(ft.title, ''), pt.user_id,
		       COALESCE(u.first_name || ' ' || u.last_name, ''),
		       COALESCE(u.email, ''), pt.payment_method, pt.provider,
		       pt.provider_request_id, COALESCE(pt.phone_number, ''),
		       pt.amount_ugx, pt.currency, pt.status, COALESCE(pt.status_message, ''),
		       pt.raw_response, pt.created_at, pt.updated_at, pt.completed_at
		FROM payment_transactions pt
		LEFT JOIN form_submissions fs ON fs.id = pt.submission_id
		LEFT JOIN form_templates ft ON ft.id = pt.template_id
		LEFT JOIN users u ON u.id = pt.user_id
		WHERE pt.submission_id = $1
		ORDER BY pt.created_at DESC
		LIMIT 1
	`
	row := r.db.QueryRow(ctx, q, submissionID)
	return r.scanTransaction(row)
}

func (r *TransactionRepo) List(ctx context.Context, f ListTransactionsFilter, p *models.PaginationParams) ([]*models.PaymentTransaction, int, error) {
	where := []string{"1=1"}
	args := []any{}
	idx := 1

	if f.UserID != "" {
		where = append(where, fmt.Sprintf("pt.user_id = $%d", idx))
		args = append(args, f.UserID)
		idx++
	}
	if f.SubmissionID != "" {
		where = append(where, fmt.Sprintf("pt.submission_id = $%d", idx))
		args = append(args, f.SubmissionID)
		idx++
	}
	if f.Status != "" && f.Status != "ALL" {
		where = append(where, fmt.Sprintf("UPPER(pt.status) = UPPER($%d)", idx))
		args = append(args, f.Status)
		idx++
	}
	if f.PaymentMethod != "" && f.PaymentMethod != "ALL" {
		where = append(where, fmt.Sprintf("UPPER(pt.payment_method) = UPPER($%d)", idx))
		args = append(args, f.PaymentMethod)
		idx++
	}
	if f.Search != "" {
		search := "%" + strings.ToLower(f.Search) + "%"
		where = append(where, fmt.Sprintf("(LOWER(pt.transaction_reference) LIKE $%d OR LOWER(COALESCE(fs.submission_reference, '')) LIKE $%d OR LOWER(COALESCE(ft.title, '')) LIKE $%d OR LOWER(COALESCE(pt.phone_number, '')) LIKE $%d OR LOWER(COALESCE(u.email, '')) LIKE $%d OR LOWER(COALESCE(u.first_name || ' ' || u.last_name, '')) LIKE $%d)", idx, idx, idx, idx, idx, idx))
		args = append(args, search)
		idx++
	}
	if f.StartDate != "" {
		where = append(where, fmt.Sprintf("pt.created_at >= $%d", idx))
		args = append(args, f.StartDate)
		idx++
	}
	if f.EndDate != "" {
		where = append(where, fmt.Sprintf("pt.created_at <= $%d", idx))
		args = append(args, f.EndDate)
		idx++
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM payment_transactions pt
		LEFT JOIN form_submissions fs ON fs.id = pt.submission_id
		LEFT JOIN form_templates ft ON ft.id = pt.template_id
		LEFT JOIN users u ON u.id = pt.user_id
		WHERE %s
	`, whereClause)

	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := 20
	offset := 0
	if p != nil {
		if p.PerPage > 0 {
			limit = p.PerPage
		}
		if p.Page > 0 {
			offset = (p.Page - 1) * limit
		}
	}

	args = append(args, limit, offset)
	q := fmt.Sprintf(`
		SELECT pt.id, pt.transaction_reference, pt.submission_id,
		       COALESCE(fs.submission_reference, ''), pt.template_id,
		       COALESCE(ft.title, ''), pt.user_id,
		       COALESCE(u.first_name || ' ' || u.last_name, ''),
		       COALESCE(u.email, ''), pt.payment_method, pt.provider,
		       pt.provider_request_id, COALESCE(pt.phone_number, ''),
		       pt.amount_ugx, pt.currency, pt.status, COALESCE(pt.status_message, ''),
		       pt.raw_response, pt.created_at, pt.updated_at, pt.completed_at
		FROM payment_transactions pt
		LEFT JOIN form_submissions fs ON fs.id = pt.submission_id
		LEFT JOIN form_templates ft ON ft.id = pt.template_id
		LEFT JOIN users u ON u.id = pt.user_id
		WHERE %s
		ORDER BY pt.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, idx, idx+1)

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*models.PaymentTransaction
	for rows.Next() {
		tx, err := r.scanTransaction(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, tx)
	}
	return list, total, nil
}

func (r *TransactionRepo) ListPendingForSync(ctx context.Context, limit int) ([]*models.PaymentTransaction, error) {
	if limit <= 0 {
		limit = 20
	}
	const q = `
		SELECT pt.id, pt.transaction_reference, pt.submission_id,
		       COALESCE(fs.submission_reference, ''), pt.template_id,
		       COALESCE(ft.title, ''), pt.user_id,
		       COALESCE(u.first_name || ' ' || u.last_name, ''),
		       COALESCE(u.email, ''), pt.payment_method, pt.provider,
		       pt.provider_request_id, COALESCE(pt.phone_number, ''),
		       pt.amount_ugx, pt.currency, pt.status, COALESCE(pt.status_message, ''),
		       pt.raw_response, pt.created_at, pt.updated_at, pt.completed_at
		FROM payment_transactions pt
		LEFT JOIN form_submissions fs ON fs.id = pt.submission_id
		LEFT JOIN form_templates ft ON ft.id = pt.template_id
		LEFT JOIN users u ON u.id = pt.user_id
		WHERE pt.status = 'PENDING'
		  AND pt.provider = 'IOTEC'
		  AND pt.provider_request_id IS NOT NULL
		  AND pt.created_at >= NOW() - INTERVAL '24 hours'
		ORDER BY pt.created_at ASC
		LIMIT $1
	`
	rows, err := r.db.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.PaymentTransaction
	for rows.Next() {
		tx, err := r.scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, tx)
	}
	return list, nil
}

func (r *TransactionRepo) GetKPIs(ctx context.Context, userID string) (*TransactionKPIs, error) {
	where := "1=1"
	args := []any{}
	if userID != "" {
		where = "user_id = $1"
		args = append(args, userID)
	}

	q := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(CASE WHEN status = 'SUCCESS' THEN amount_ugx ELSE 0 END), 0) AS total_amount,
			COUNT(*) AS total_count,
			COUNT(CASE WHEN status = 'SUCCESS' THEN 1 END) AS success_count,
			COUNT(CASE WHEN status = 'PENDING' THEN 1 END) AS pending_count,
			COUNT(CASE WHEN status = 'FAILED' THEN 1 END) AS failed_count
		FROM payment_transactions
		WHERE %s
	`, where)

	var kpi TransactionKPIs
	err := r.db.QueryRow(ctx, q, args...).Scan(
		&kpi.TotalAmountUGX, &kpi.TotalCount, &kpi.SuccessCount,
		&kpi.PendingCount, &kpi.FailedCount,
	)
	if err != nil {
		return nil, err
	}
	return &kpi, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func (r *TransactionRepo) scanTransaction(s rowScanner) (*models.PaymentTransaction, error) {
	var tx models.PaymentTransaction
	var rawResp []byte
	err := s.Scan(
		&tx.ID, &tx.TransactionReference, &tx.SubmissionID,
		&tx.SubmissionReference, &tx.TemplateID,
		&tx.TemplateTitle, &tx.UserID,
		&tx.ApplicantName, &tx.ApplicantEmail,
		&tx.PaymentMethod, &tx.Provider,
		&tx.ProviderRequestID, &tx.PhoneNumber,
		&tx.AmountUGX, &tx.Currency, &tx.Status, &tx.StatusMessage,
		&rawResp, &tx.CreatedAt, &tx.UpdatedAt, &tx.CompletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if len(rawResp) > 0 {
		tx.RawResponse = json.RawMessage(rawResp)
	} else {
		tx.RawResponse = json.RawMessage("{}")
	}
	return &tx, nil
}
