package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ── Departments (sourced from the `roles` table) ────────────────
//
// Departments are no longer a separate concept — every role in the system
// is also a department on the dynamic form builder. Adding a new role
// therefore automatically adds a department option, and the role-id is
// stored as the form's department_id.

type DepartmentRepo struct{ db *pgxpool.Pool }

// List returns every role, excluding the bottom-tier "user" applicant role
// (which would never own a form) and "super_admin" (which sees everything).
func (r *DepartmentRepo) List(ctx context.Context, activeOnly bool) ([]*models.Department, error) {
	const q = `SELECT id, name, COALESCE(description,''), created_at, updated_at
	           FROM roles
	           WHERE name NOT IN ('user','super_admin')
	           ORDER BY name`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*models.Department{}
	for rows.Next() {
		d := &models.Department{IsActive: true}
		if err := rows.Scan(&d.ID, &d.Name, &d.Description, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		d.Code = d.Name
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *DepartmentRepo) GetByID(ctx context.Context, id string) (*models.Department, error) {
	const q = `SELECT id, name, COALESCE(description,''), created_at, updated_at
	           FROM roles WHERE id=$1`
	d := &models.Department{IsActive: true}
	err := r.db.QueryRow(ctx, q, id).Scan(&d.ID, &d.Name, &d.Description, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	d.Code = d.Name
	return d, err
}

// UserDepartmentIDs returns every role-id the user currently holds.
// The form-tenancy check uses these to decide whether the user can
// manage a form owned by a given role.
func (r *DepartmentRepo) UserDepartmentIDs(ctx context.Context, userID string) ([]string, error) {
	const q = `SELECT role_id FROM user_roles WHERE user_id = $1`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ── Form Templates ───────────────────────────────────────────────

type FormRepo struct{ db *pgxpool.Pool }

func (r *FormRepo) CountSubmissionsByStatus(ctx context.Context) (map[string]int64, error) {
	rows, err := r.db.Query(ctx, `SELECT status, COUNT(*) FROM form_submissions GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int64{}
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		out[status] = count
	}
	return out, rows.Err()
}

func (r *FormRepo) CountOpenTemplates(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM form_templates WHERE status='OPEN'`).Scan(&count)
	return count, err
}

type ListFormTemplatesFilter struct {
	DepartmentIDs []string // empty = no filter (super-admin)
	Status        string   // "" = all
	PublicOnly    bool     // true = only OPEN
}

func (r *FormRepo) ListTemplates(ctx context.Context, f ListFormTemplatesFilter) ([]*models.FormTemplate, error) {
	conds := []string{"1=1"}
	args := []any{}
	if len(f.DepartmentIDs) > 0 {
		placeholders := make([]string, len(f.DepartmentIDs))
		for i, id := range f.DepartmentIDs {
			args = append(args, id)
			placeholders[i] = fmt.Sprintf("$%d", len(args))
		}
		conds = append(conds, "t.department_id IN ("+strings.Join(placeholders, ",")+")")
	}
	if f.Status != "" {
		args = append(args, f.Status)
		conds = append(conds, fmt.Sprintf("t.status = $%d", len(args)))
	}
	if f.PublicOnly {
		conds = append(conds, "t.status = 'OPEN'")
	}
	q := `SELECT t.id, t.department_id, COALESCE(r.name,''), t.slug, t.title, t.description,
	             t.banner_image_url, t.price_ugx, t.status, t.legacy_form_type, t.created_by,
	             t.created_at, t.updated_at
	      FROM form_templates t
	      LEFT JOIN roles r ON r.id = t.department_id
	      WHERE ` + strings.Join(conds, " AND ") + `
	      ORDER BY t.updated_at DESC`
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*models.FormTemplate{}
	for rows.Next() {
		t := &models.FormTemplate{}
		if err := rows.Scan(&t.ID, &t.DepartmentID, &t.DepartmentName, &t.Slug, &t.Title, &t.Description,
			&t.BannerImageURL, &t.PriceUGX, &t.Status, &t.LegacyFormType, &t.CreatedBy,
			&t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *FormRepo) GetTemplate(ctx context.Context, id string, withFields bool) (*models.FormTemplate, error) {
	const q = `SELECT t.id, t.department_id, COALESCE(r.name,''), t.slug, t.title, t.description,
	                  t.banner_image_url, t.price_ugx, t.status, t.legacy_form_type, t.created_by,
	                  t.created_at, t.updated_at
	           FROM form_templates t
	           LEFT JOIN roles r ON r.id = t.department_id
	           WHERE t.id = $1`
	t := &models.FormTemplate{}
	err := r.db.QueryRow(ctx, q, id).Scan(&t.ID, &t.DepartmentID, &t.DepartmentName, &t.Slug, &t.Title, &t.Description,
		&t.BannerImageURL, &t.PriceUGX, &t.Status, &t.LegacyFormType, &t.CreatedBy,
		&t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if withFields {
		fields, err := r.ListFields(ctx, id)
		if err != nil {
			return nil, err
		}
		t.Fields = fields
	}
	return t, nil
}

func (r *FormRepo) GetTemplateBySlug(ctx context.Context, slug string) (*models.FormTemplate, error) {
	var id string
	err := r.db.QueryRow(ctx, `SELECT id FROM form_templates WHERE slug=$1`, slug).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.GetTemplate(ctx, id, true)
}

func (r *FormRepo) CreateTemplate(ctx context.Context, t *models.FormTemplate) error {
	const q = `INSERT INTO form_templates
	            (department_id, slug, title, description, banner_image_url, price_ugx, status, created_by)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	           RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, q,
		t.DepartmentID, t.Slug, t.Title, t.Description, t.BannerImageURL, t.PriceUGX, t.Status, t.CreatedBy,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *FormRepo) UpdateTemplate(ctx context.Context, t *models.FormTemplate) error {
	const q = `UPDATE form_templates
	           SET department_id=$2, title=$3, description=$4, banner_image_url=$5,
	               price_ugx=$6, status=$7, updated_at=NOW()
	           WHERE id=$1
	           RETURNING updated_at`
	return r.db.QueryRow(ctx, q,
		t.ID, t.DepartmentID, t.Title, t.Description, t.BannerImageURL, t.PriceUGX, t.Status,
	).Scan(&t.UpdatedAt)
}

func (r *FormRepo) DeleteTemplate(ctx context.Context, id string) error {
	// Soft-archive instead of destroy so historical submissions stay valid.
	_, err := r.db.Exec(ctx, `UPDATE form_templates SET status='ARCHIVED', updated_at=NOW() WHERE id=$1`, id)
	return err
}

// ── Form Fields ──────────────────────────────────────────────────

func (r *FormRepo) ListFields(ctx context.Context, templateID string) ([]*models.FormField, error) {
	const q = `SELECT id, template_id, field_key, field_type, label, placeholder, help_text,
	                  is_required, order_index, config, created_at, updated_at
	           FROM form_fields
	           WHERE template_id=$1 AND deleted_at IS NULL
	           ORDER BY order_index, created_at`
	rows, err := r.db.Query(ctx, q, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*models.FormField{}
	for rows.Next() {
		f := &models.FormField{}
		var cfg []byte
		if err := rows.Scan(&f.ID, &f.TemplateID, &f.FieldKey, &f.FieldType, &f.Label, &f.Placeholder, &f.HelpText,
			&f.IsRequired, &f.OrderIndex, &cfg, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		if len(cfg) > 0 {
			f.Config = json.RawMessage(cfg)
		} else {
			f.Config = json.RawMessage("{}")
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// ReplaceFields wipes the template's fields and re-inserts the provided list
// in order. Soft-deletes existing rows (so submissions referencing them stay
// readable) and creates new rows in a single transaction.
func (r *FormRepo) ReplaceFields(ctx context.Context, templateID string, fields []*models.FormField) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx,
		`UPDATE form_fields SET deleted_at=NOW(), updated_at=NOW()
		 WHERE template_id=$1 AND deleted_at IS NULL`, templateID); err != nil {
		return fmt.Errorf("soft-delete old fields: %w", err)
	}

	for i, f := range fields {
		cfg := f.Config
		if len(cfg) == 0 {
			cfg = json.RawMessage("{}")
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO form_fields
			  (template_id, field_key, field_type, label, placeholder, help_text,
			   is_required, order_index, config)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb)`,
			templateID, f.FieldKey, f.FieldType, f.Label, f.Placeholder, f.HelpText,
			f.IsRequired, i, string(cfg))
		if err != nil {
			return fmt.Errorf("insert field %s: %w", f.FieldKey, err)
		}
	}

	_, _ = tx.Exec(ctx, `UPDATE form_templates SET updated_at=NOW() WHERE id=$1`, templateID)
	return tx.Commit(ctx)
}

// ── Submissions ──────────────────────────────────────────────────

type ListSubmissionsFilter struct {
	TemplateID    string
	DepartmentIDs []string
	UserID        string
	Status        string
}

func (r *FormRepo) ListSubmissions(ctx context.Context, f ListSubmissionsFilter, p *models.PaginationParams) ([]*models.FormSubmission, int64, error) {
	if p == nil {
		p = &models.PaginationParams{}
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage <= 0 || p.PerPage > 200 {
		p.PerPage = 25
	}
	conds := []string{"1=1"}
	args := []any{}
	add := func(col, val string) {
		args = append(args, val)
		conds = append(conds, fmt.Sprintf("s.%s = $%d", col, len(args)))
	}
	if f.TemplateID != "" {
		add("template_id", f.TemplateID)
	}
	if len(f.DepartmentIDs) > 0 {
		placeholders := make([]string, len(f.DepartmentIDs))
		for i, id := range f.DepartmentIDs {
			args = append(args, id)
			placeholders[i] = fmt.Sprintf("$%d", len(args))
		}
		conds = append(conds, "s.department_id IN ("+strings.Join(placeholders, ",")+")")
	}
	if f.UserID != "" {
		add("user_id", f.UserID)
	}
	if f.Status != "" {
		add("status", f.Status)
	}

	where := strings.Join(conds, " AND ")

	var total int64
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM form_submissions s WHERE `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, p.PerPage, (p.Page-1)*p.PerPage)
	limitClause := fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	q := `SELECT s.id, s.template_id, COALESCE(t.title,''), s.department_id, s.user_id,
	             COALESCE(u.first_name || ' ' || u.last_name,''), COALESCE(u.email,''),
	             COALESCE(s.submission_reference,''), s.status, s.payment_status,
	             COALESCE(s.payment_reference,''), s.payment_amount_ugx, s.payment_verified_at,
	             s.answers, s.reviewer_id, COALESCE(s.review_notes,''),
	             s.submitted_at, s.approved_at, s.rejected_at, s.created_at, s.updated_at
	      FROM form_submissions s
	      LEFT JOIN form_templates t ON t.id = s.template_id
	      LEFT JOIN users u          ON u.id = s.user_id
	      WHERE ` + where + ` ORDER BY s.updated_at DESC` + limitClause

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []*models.FormSubmission{}
	for rows.Next() {
		s := &models.FormSubmission{}
		var answers []byte
		if err := rows.Scan(&s.ID, &s.TemplateID, &s.TemplateTitle, &s.DepartmentID, &s.UserID,
			&s.ApplicantName, &s.ApplicantEmail,
			&s.SubmissionReference, &s.Status, &s.PaymentStatus,
			&s.PaymentReference, &s.PaymentAmountUGX, &s.PaymentVerifiedAt,
			&answers, &s.ReviewerID, &s.ReviewNotes,
			&s.SubmittedAt, &s.ApprovedAt, &s.RejectedAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if len(answers) > 0 {
			s.Answers = json.RawMessage(answers)
		} else {
			s.Answers = json.RawMessage("{}")
		}
		out = append(out, s)
	}
	return out, total, rows.Err()
}

func (r *FormRepo) GetSubmission(ctx context.Context, id string) (*models.FormSubmission, error) {
	const q = `SELECT s.id, s.template_id, COALESCE(t.title,''), s.department_id, s.user_id,
	                  COALESCE(u.first_name || ' ' || u.last_name,''), COALESCE(u.email,''),
	                  COALESCE(s.submission_reference,''), s.status, s.payment_status,
	                  COALESCE(s.payment_reference,''), s.payment_amount_ugx, s.payment_verified_at,
	                  s.answers, s.reviewer_id, COALESCE(s.review_notes,''),
	                  s.submitted_at, s.approved_at, s.rejected_at, s.created_at, s.updated_at
	           FROM form_submissions s
	           LEFT JOIN form_templates t ON t.id = s.template_id
	           LEFT JOIN users u          ON u.id = s.user_id
	           WHERE s.id=$1`
	s := &models.FormSubmission{}
	var answers []byte
	err := r.db.QueryRow(ctx, q, id).Scan(&s.ID, &s.TemplateID, &s.TemplateTitle, &s.DepartmentID, &s.UserID,
		&s.ApplicantName, &s.ApplicantEmail,
		&s.SubmissionReference, &s.Status, &s.PaymentStatus,
		&s.PaymentReference, &s.PaymentAmountUGX, &s.PaymentVerifiedAt,
		&answers, &s.ReviewerID, &s.ReviewNotes,
		&s.SubmittedAt, &s.ApprovedAt, &s.RejectedAt, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if len(answers) > 0 {
		s.Answers = json.RawMessage(answers)
	} else {
		s.Answers = json.RawMessage("{}")
	}
	return s, nil
}

func (r *FormRepo) GetUserDraftForTemplate(ctx context.Context, userID, templateID string) (*models.FormSubmission, error) {
	var id string
	err := r.db.QueryRow(ctx,
		`SELECT id FROM form_submissions WHERE user_id=$1 AND template_id=$2 AND status='DRAFT'
		 ORDER BY updated_at DESC LIMIT 1`, userID, templateID).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.GetSubmission(ctx, id)
}

func (r *FormRepo) CreateSubmission(ctx context.Context, s *models.FormSubmission) error {
	answers := s.Answers
	if len(answers) == 0 {
		answers = json.RawMessage("{}")
	}
	const q = `INSERT INTO form_submissions
	            (template_id, department_id, user_id, status, payment_status, answers)
	           VALUES ($1,$2,$3,$4,$5,$6::jsonb)
	           RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, q,
		s.TemplateID, s.DepartmentID, s.UserID, s.Status, s.PaymentStatus, string(answers),
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *FormRepo) UpdateAnswers(ctx context.Context, id string, answers json.RawMessage) error {
	if len(answers) == 0 {
		answers = json.RawMessage("{}")
	}
	_, err := r.db.Exec(ctx,
		`UPDATE form_submissions SET answers=$2::jsonb, updated_at=NOW() WHERE id=$1`,
		id, string(answers))
	return err
}

func (r *FormRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE form_submissions SET status=$2, updated_at=NOW() WHERE id=$1`, id, status)
	return err
}

func (r *FormRepo) SetSubmitted(ctx context.Context, id, reference string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE form_submissions
		SET status='SUBMITTED', submission_reference=$2, submitted_at=NOW(), updated_at=NOW()
		WHERE id=$1`, id, reference)
	return err
}

func (r *FormRepo) SetPaymentProof(ctx context.Context, id, reference string, amount float64) error {
	_, err := r.db.Exec(ctx, `
		UPDATE form_submissions
		SET payment_status='PROOF_UPLOADED', payment_reference=$2,
		    payment_amount_ugx=$3, updated_at=NOW()
		WHERE id=$1`, id, reference, amount)
	return err
}

func (r *FormRepo) VerifyPayment(ctx context.Context, id, reviewerID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE form_submissions
		SET payment_status='PAID', payment_verified_by=$2, payment_verified_at=NOW(),
		    updated_at=NOW()
		WHERE id=$1`, id, reviewerID)
	return err
}

func (r *FormRepo) Review(ctx context.Context, id, status, reviewerID, notes string) error {
	col := ""
	switch status {
	case models.SubStatusApproved:
		col = "approved_at"
	case models.SubStatusRejected:
		col = "rejected_at"
	}
	if col == "" {
		_, err := r.db.Exec(ctx, `
			UPDATE form_submissions
			SET status=$2, reviewer_id=$3, review_notes=$4, updated_at=NOW()
			WHERE id=$1`, id, status, reviewerID, notes)
		return err
	}
	q := fmt.Sprintf(`
		UPDATE form_submissions
		SET status=$2, reviewer_id=$3, review_notes=$4, %s=NOW(), updated_at=NOW()
		WHERE id=$1`, col)
	_, err := r.db.Exec(ctx, q, id, status, reviewerID, notes)
	return err
}
