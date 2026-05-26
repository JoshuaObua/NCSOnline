package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("record not found")
var ErrDuplicate = errors.New("duplicate record")

type Repos struct {
	Users        *UserRepo
	Roles        *RoleRepo
	Applications *ApplicationRepo
	Audit        *AuditRepo
	Tokens       *TokenRepo
}

func New(db *pgxpool.Pool) *Repos {
	return &Repos{
		Users:        &UserRepo{db},
		Roles:        &RoleRepo{db},
		Applications: &ApplicationRepo{db},
		Audit:        &AuditRepo{db},
		Tokens:       &TokenRepo{db},
	}
}

// ── User Repository ───────────────────────────────────────────────

type UserRepo struct{ db *pgxpool.Pool }

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	const q = `INSERT INTO users (id, email, password_hash, first_name, last_name, phone)
	           VALUES ($1,$2,$3,$4,$5,$6)
	           RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q,
		u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Phone,
	).Scan(&u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	const q = `SELECT id, email, password_hash, first_name, last_name, phone,
	                  is_active, is_email_verified, email_verified_at, last_login_at,
	                  created_at, updated_at
	           FROM users WHERE id=$1 AND deleted_at IS NULL`
	u := &models.User{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Phone,
		&u.IsActive, &u.IsEmailVerified, &u.EmailVerifiedAt, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const q = `SELECT id, email, password_hash, first_name, last_name, phone,
	                  is_active, is_email_verified, email_verified_at, last_login_at,
	                  created_at, updated_at
	           FROM users WHERE email=$1 AND deleted_at IS NULL`
	u := &models.User{}
	err := r.db.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Phone,
		&u.IsActive, &u.IsEmailVerified, &u.EmailVerifiedAt, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

func (r *UserRepo) List(ctx context.Context, p *models.PaginationParams) ([]*models.User, int64, error) {
	const countQ = `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL
	                AND ($1='' OR first_name ILIKE $1 OR last_name ILIKE $1 OR email ILIKE $1)`
	const q = `SELECT id, email, first_name, last_name, phone, is_active, is_email_verified,
	                  last_login_at, created_at, updated_at
	           FROM users WHERE deleted_at IS NULL
	           AND ($1='' OR first_name ILIKE $1 OR last_name ILIKE $1 OR email ILIKE $1)
	           ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	search := ""
	if p.Search != "" {
		search = "%" + p.Search + "%"
	}

	var total int64
	if err := r.db.QueryRow(ctx, countQ, search).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, q, search, p.PerPage, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone,
			&u.IsActive, &u.IsEmailVerified, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (r *UserRepo) Update(ctx context.Context, u *models.User) error {
	const q = `UPDATE users SET first_name=$2, last_name=$3, phone=$4, updated_at=NOW()
	           WHERE id=$1 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, q, u.ID, u.FirstName, u.LastName, u.Phone)
	return err
}

func (r *UserRepo) UpdatePassword(ctx context.Context, userID, hash string) error {
	const q = `UPDATE users SET password_hash=$2, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, userID, hash)
	return err
}

func (r *UserRepo) SetActive(ctx context.Context, userID string, active bool) error {
	const q = `UPDATE users SET is_active=$2, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, userID, active)
	return err
}

func (r *UserRepo) SoftDelete(ctx context.Context, userID string) error {
	const q = `UPDATE users SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, userID)
	return err
}

func (r *UserRepo) UpdateLastLogin(ctx context.Context, userID string) error {
	const q = `UPDATE users SET last_login_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, userID)
	return err
}

func (r *UserRepo) GetRoles(ctx context.Context, userID string) ([]models.Role, error) {
	const q = `SELECT r.id, r.name, r.description, r.is_system, r.created_at
	           FROM roles r
	           JOIN user_roles ur ON ur.role_id = r.id
	           WHERE ur.user_id = $1`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *UserRepo) AssignRole(ctx context.Context, userID, roleID, assignedBy string) error {
	const q = `INSERT INTO user_roles (user_id, role_id, assigned_by)
	           VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, q, userID, roleID, assignedBy)
	return err
}

func (r *UserRepo) RemoveRole(ctx context.Context, userID, roleID string) error {
	const q = `DELETE FROM user_roles WHERE user_id=$1 AND role_id=$2`
	_, err := r.db.Exec(ctx, q, userID, roleID)
	return err
}

func (r *UserRepo) CountAll(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&n)
	return n, err
}

// ── Role Repository ───────────────────────────────────────────────

type RoleRepo struct{ db *pgxpool.Pool }

func (r *RoleRepo) List(ctx context.Context) ([]models.Role, error) {
	const q = `SELECT id, name, description, is_system, created_at FROM roles ORDER BY name`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *RoleRepo) GetByID(ctx context.Context, id string) (*models.Role, error) {
	const q = `SELECT id, name, description, is_system, created_at FROM roles WHERE id=$1`
	role := &models.Role{}
	err := r.db.QueryRow(ctx, q, id).Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return role, err
}

func (r *RoleRepo) Create(ctx context.Context, role *models.Role) error {
	const q = `INSERT INTO roles (id, name, description) VALUES ($1,$2,$3) RETURNING created_at`
	err := r.db.QueryRow(ctx, q, role.ID, role.Name, role.Description).Scan(&role.CreatedAt)
	if err != nil && isDuplicate(err) {
		return ErrDuplicate
	}
	return err
}

func (r *RoleRepo) Update(ctx context.Context, role *models.Role) error {
	const q = `UPDATE roles SET name=$2, description=$3 WHERE id=$1 AND is_system=false`
	_, err := r.db.Exec(ctx, q, role.ID, role.Name, role.Description)
	return err
}

func (r *RoleRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM roles WHERE id=$1 AND is_system=false`
	_, err := r.db.Exec(ctx, q, id)
	return err
}

func (r *RoleRepo) ListPermissions(ctx context.Context) ([]models.Permission, error) {
	const q = `SELECT id, name, description, resource, action FROM permissions ORDER BY resource, action`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var perms []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Resource, &p.Action); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}

func (r *RoleRepo) AssignPermission(ctx context.Context, roleID, permID string) error {
	const q = `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, q, roleID, permID)
	return err
}

func (r *RoleRepo) RemovePermission(ctx context.Context, roleID, permID string) error {
	const q = `DELETE FROM role_permissions WHERE role_id=$1 AND permission_id=$2`
	_, err := r.db.Exec(ctx, q, roleID, permID)
	return err
}

// ── Token Repository ──────────────────────────────────────────────

type TokenRepo struct{ db *pgxpool.Pool }

func (r *TokenRepo) StoreRefreshToken(ctx context.Context, t *models.RefreshToken) error {
	const q = `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, ip_address, user_agent)
	           VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.Exec(ctx, q, t.ID, t.UserID, t.TokenHash, t.ExpiresAt, t.IPAddress, t.UserAgent)
	return err
}

func (r *TokenRepo) GetByHash(ctx context.Context, hash string) (*models.RefreshToken, error) {
	const q = `SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
	           FROM refresh_tokens WHERE token_hash=$1`
	t := &models.RefreshToken{}
	err := r.db.QueryRow(ctx, q, hash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.CreatedAt, &t.RevokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return t, err
}

func (r *TokenRepo) Revoke(ctx context.Context, id string) error {
	const q = `UPDATE refresh_tokens SET revoked_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id)
	return err
}

func (r *TokenRepo) RevokeAllForUser(ctx context.Context, userID string) error {
	const q = `UPDATE refresh_tokens SET revoked_at=NOW() WHERE user_id=$1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, q, userID)
	return err
}

// ── Application Repository ────────────────────────────────────────

type ApplicationRepo struct{ db *pgxpool.Pool }

func (r *ApplicationRepo) Create(ctx context.Context, a *models.Application) error {
	const q = `INSERT INTO applications
	           (id, user_id, form_type, application_type, organisation_type, status,
	            payment_status, form_data, last_saved_step, draft_expires_at)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	           RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q,
		a.ID, a.UserID, a.FormType, a.ApplicationType, a.OrganisationType,
		a.Status, a.PaymentStatus, a.FormData, a.LastSavedStep, a.DraftExpiresAt,
	).Scan(&a.CreatedAt, &a.UpdatedAt)
}

func (r *ApplicationRepo) GetByID(ctx context.Context, id string) (*models.Application, error) {
	const q = `SELECT id, application_reference, user_id, form_type, application_type,
	                  organisation_type, status, payment_status, payment_method,
	                  payment_reference, payment_amount_ugx, signed_form_url,
	                  signed_form_uploaded_at, submitted_at, reviewer_id, review_notes,
	                  approved_at, rejected_at, last_saved_step, form_data,
	                  created_at, updated_at
	           FROM applications WHERE id=$1`
	a := &models.Application{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&a.ID, &a.ApplicationReference, &a.UserID, &a.FormType, &a.ApplicationType,
		&a.OrganisationType, &a.Status, &a.PaymentStatus, &a.PaymentMethod,
		&a.PaymentReference, &a.PaymentAmountUGX, &a.SignedFormURL,
		&a.SignedFormUploadedAt, &a.SubmittedAt, &a.ReviewerID, &a.ReviewNotes,
		&a.ApprovedAt, &a.RejectedAt, &a.LastSavedStep, &a.FormData,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return a, err
}

func (r *ApplicationRepo) GetDraftByUserAndType(ctx context.Context, userID, formType string) (*models.Application, error) {
	const q = `SELECT id, user_id, form_type, status, form_data, last_saved_step, created_at, updated_at
	           FROM applications WHERE user_id=$1 AND form_type=$2 AND status=$3
	           ORDER BY updated_at DESC LIMIT 1`
	a := &models.Application{}
	err := r.db.QueryRow(ctx, q, userID, formType, models.StatusDraft).Scan(
		&a.ID, &a.UserID, &a.FormType, &a.Status, &a.FormData, &a.LastSavedStep,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return a, err
}

func (r *ApplicationRepo) UpdateDraft(ctx context.Context, a *models.Application) error {
	const q = `UPDATE applications SET form_data=$2, last_saved_step=$3, updated_at=NOW()
	           WHERE id=$1`
	_, err := r.db.Exec(ctx, q, a.ID, a.FormData, a.LastSavedStep)
	return err
}

func (r *ApplicationRepo) UpdateStatus(ctx context.Context, id, status string) error {
	const q = `UPDATE applications SET status=$2, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id, status)
	return err
}

func (r *ApplicationRepo) SetSignedForm(ctx context.Context, id, url string) error {
	const q = `UPDATE applications SET signed_form_url=$2, signed_form_uploaded_at=NOW(),
	           status=$3, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id, url, models.StatusPendingPayment)
	return err
}

func (r *ApplicationRepo) SetPaymentProof(ctx context.Context, id, method, ref string, amount float64) error {
	const q = `UPDATE applications SET payment_method=$2, payment_reference=$3,
	           payment_amount_ugx=$4, payment_status=$5, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id, method, ref, amount, models.PaymentProofUploaded)
	return err
}

func (r *ApplicationRepo) Submit(ctx context.Context, id, ref string) error {
	const q = `UPDATE applications SET application_reference=$2, status=$3,
	           submitted_at=NOW(), updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id, ref, models.StatusSubmitted)
	return err
}

func (r *ApplicationRepo) Approve(ctx context.Context, id, reviewerID, notes string) error {
	const q = `UPDATE applications SET status=$2, reviewer_id=$3, review_notes=$4,
	           approved_at=NOW(), updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id, models.StatusApproved, reviewerID, notes)
	return err
}

func (r *ApplicationRepo) Reject(ctx context.Context, id, reviewerID, notes string) error {
	const q = `UPDATE applications SET status=$2, reviewer_id=$3, review_notes=$4,
	           rejected_at=NOW(), updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id, models.StatusRejected, reviewerID, notes)
	return err
}

func (r *ApplicationRepo) RequestInfo(ctx context.Context, id, reviewerID, notes string) error {
	const q = `UPDATE applications SET status=$2, reviewer_id=$3, review_notes=$4,
	           updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id, models.StatusNeedsInformation, reviewerID, notes)
	return err
}

func (r *ApplicationRepo) VerifyPayment(ctx context.Context, id, verifiedBy string) error {
	const q = `UPDATE applications SET payment_status=$2, payment_verified_by=$3,
	           payment_verified_at=NOW(), updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id, models.PaymentVerified, verifiedBy)
	return err
}

func (r *ApplicationRepo) RejectPayment(ctx context.Context, id, notes string) error {
	const q = `UPDATE applications SET payment_status=$2, review_notes=$3, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id, models.PaymentRejected, notes)
	return err
}

func (r *ApplicationRepo) ListByUser(ctx context.Context, userID string, p *models.PaginationParams) ([]*models.Application, int64, error) {
	const countQ = `SELECT COUNT(*) FROM applications WHERE user_id=$1`
	const q = `SELECT id, application_reference, form_type, application_type, status,
	                  payment_status, submitted_at, created_at, updated_at
	           FROM applications WHERE user_id=$1
	           ORDER BY updated_at DESC LIMIT $2 OFFSET $3`
	var total int64
	if err := r.db.QueryRow(ctx, countQ, userID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, q, userID, p.PerPage, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return scanApplicationRows(rows, total)
}

func (r *ApplicationRepo) AdminList(ctx context.Context, status string, p *models.PaginationParams) ([]*models.Application, int64, error) {
	filter := ""
	args := []interface{}{p.PerPage, p.Offset()}
	if status != "" {
		filter = "WHERE status=$3"
		args = append(args, status)
	}
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM applications %s`, filter)
	listQ := fmt.Sprintf(`SELECT id, application_reference, user_id, form_type, application_type,
	                             status, payment_status, submitted_at, created_at, updated_at
	                      FROM applications %s ORDER BY updated_at DESC LIMIT $1 OFFSET $2`, filter)

	var total int64
	if err := r.db.QueryRow(ctx, countQ, args[2:]...).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, listQ, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return scanApplicationRows(rows, total)
}

func (r *ApplicationRepo) CountByStatus(ctx context.Context) (map[string]int64, error) {
	const q = `SELECT status, COUNT(*) FROM applications GROUP BY status`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[status] = count
	}
	return result, rows.Err()
}

func scanApplicationRows(rows pgx.Rows, total int64) ([]*models.Application, int64, error) {
	var apps []*models.Application
	for rows.Next() {
		a := &models.Application{}
		if err := rows.Scan(&a.ID, &a.ApplicationReference, &a.UserID, &a.FormType,
			&a.ApplicationType, &a.Status, &a.PaymentStatus, &a.SubmittedAt,
			&a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, err
		}
		apps = append(apps, a)
	}
	return apps, total, rows.Err()
}

// ── Attachment Repository ─────────────────────────────────────────

func (r *ApplicationRepo) AddAttachment(ctx context.Context, a *models.Attachment) error {
	const q = `INSERT INTO application_attachments (id, application_id, field_name, file_name, file_url, file_size, mime_type)
	           VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING uploaded_at`
	return r.db.QueryRow(ctx, q,
		a.ID, a.ApplicationID, a.FieldName, a.FileName, a.FileURL, a.FileSize, a.MimeType,
	).Scan(&a.UploadedAt)
}

func (r *ApplicationRepo) ListAttachments(ctx context.Context, applicationID string) ([]models.Attachment, error) {
	const q = `SELECT id, application_id, field_name, file_name, file_url, file_size, mime_type, uploaded_at
	           FROM application_attachments WHERE application_id=$1`
	rows, err := r.db.Query(ctx, q, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var atts []models.Attachment
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(&a.ID, &a.ApplicationID, &a.FieldName, &a.FileName, &a.FileURL,
			&a.FileSize, &a.MimeType, &a.UploadedAt); err != nil {
			return nil, err
		}
		atts = append(atts, a)
	}
	return atts, rows.Err()
}

func (r *ApplicationRepo) DeleteAttachment(ctx context.Context, id, applicationID string) error {
	const q = `DELETE FROM application_attachments WHERE id=$1 AND application_id=$2`
	_, err := r.db.Exec(ctx, q, id, applicationID)
	return err
}

// ── Audit Repository ──────────────────────────────────────────────

type AuditRepo struct{ db *pgxpool.Pool }

func (r *AuditRepo) Log(ctx context.Context, a *models.AuditLog) error {
	const q = `INSERT INTO audit_logs (id, user_id, action, resource, resource_id, old_values, new_values, ip_address, user_agent)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING created_at`
	return r.db.QueryRow(ctx, q,
		a.ID, a.UserID, a.Action, a.Resource, a.ResourceID,
		a.OldValues, a.NewValues, a.IPAddress, a.UserAgent,
	).Scan(&a.CreatedAt)
}

func (r *AuditRepo) List(ctx context.Context, p *models.PaginationParams) ([]*models.AuditLog, int64, error) {
	const countQ = `SELECT COUNT(*) FROM audit_logs WHERE ($1='' OR action ILIKE $1 OR resource ILIKE $1)`
	const q = `SELECT id, user_id, action, resource, resource_id, ip_address, created_at
	           FROM audit_logs WHERE ($1='' OR action ILIKE $1 OR resource ILIKE $1)
	           ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	search := ""
	if p.Search != "" {
		search = "%" + p.Search + "%"
	}
	var total int64
	if err := r.db.QueryRow(ctx, countQ, search).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, q, search, p.PerPage, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var logs []*models.AuditLog
	for rows.Next() {
		l := &models.AuditLog{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.Resource, &l.ResourceID, &l.IPAddress, &l.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, total, rows.Err()
}

func (r *AuditRepo) GetByID(ctx context.Context, id string) (*models.AuditLog, error) {
	const q = `SELECT id, user_id, action, resource, resource_id, old_values, new_values, ip_address, user_agent, created_at
	           FROM audit_logs WHERE id=$1`
	l := &models.AuditLog{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&l.ID, &l.UserID, &l.Action, &l.Resource, &l.ResourceID,
		&l.OldValues, &l.NewValues, &l.IPAddress, &l.UserAgent, &l.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return l, err
}

// ── Helpers ───────────────────────────────────────────────────────

func isDuplicate(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "unique constraint"))
}

func init() {
	// suppress unused import warning — time is used in token expiry
	_ = time.Now
}
