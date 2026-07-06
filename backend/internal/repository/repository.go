package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("record not found")
var ErrDuplicate = errors.New("duplicate record")

type Repos struct {
	Users         *UserRepo
	Roles         *RoleRepo
	Applications  *ApplicationRepo
	Audit         *AuditRepo
	Tokens        *TokenRepo
	CMS           *CMSRepo
	NSMIS         *NSMISRepo
	USSD          *USSDRepo
	Organisations *OrganisationRepo
	Notifications *NotificationRepo
	Operator      *OperatorRepo
	Backups       *BackupRepo
	Forms         *FormRepo
	Departments   *DepartmentRepo
	Security      *SecurityRepo
	Analytics     *AnalyticsRepo
}

func New(db *pgxpool.Pool) *Repos {
	return &Repos{
		Users:         &UserRepo{db},
		Roles:         &RoleRepo{db},
		Applications:  &ApplicationRepo{db},
		Audit:         &AuditRepo{db},
		Tokens:        &TokenRepo{db},
		CMS:           &CMSRepo{db},
		NSMIS:         &NSMISRepo{db},
		USSD:          &USSDRepo{db: db},
		Organisations: &OrganisationRepo{db: db},
		Notifications: &NotificationRepo{db: db},
		Operator:      &OperatorRepo{db: db},
		Backups:       &BackupRepo{db: db},
		Forms:         &FormRepo{db: db},
		Departments:   &DepartmentRepo{db: db},
		Security:      &SecurityRepo{db: db},
		Analytics:     &AnalyticsRepo{db: db},
	}
}

// ── User Repository ───────────────────────────────────────────────

type UserRepo struct{ db *pgxpool.Pool }

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	const q = `INSERT INTO users (id, email, password_hash, first_name, last_name, phone, avatar_url, auth_provider, google_sub, is_email_verified, email_verified_at)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,CASE WHEN $10 THEN NOW() ELSE NULL END)
	           RETURNING created_at, updated_at`
	if u.AuthProvider == "" {
		u.AuthProvider = "password"
	}
	var googleSub interface{}
	if u.GoogleSub != "" {
		googleSub = u.GoogleSub
	}
	return r.db.QueryRow(ctx, q,
		u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Phone, u.AvatarURL, u.AuthProvider, googleSub, u.IsEmailVerified,
	).Scan(&u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	const q = `SELECT id, email, password_hash, first_name, last_name, COALESCE(phone,''),
	                  COALESCE(avatar_url,''), COALESCE(auth_provider,''), COALESCE(google_sub,''),
	                  COALESCE(pin_hash,''), COALESCE(pin_change_required, TRUE),
	                  is_active, account_status, status_reason, fraud_flag, fraud_reason,
	                  suspended_until, status_changed_at,
	                  is_email_verified, email_verified_at, last_login_at,
	                  auth_invalid_before,
	                  created_at, updated_at
	           FROM users WHERE id=$1 AND deleted_at IS NULL`
	u := &models.User{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Phone,
		&u.AvatarURL, &u.AuthProvider, &u.GoogleSub,
		&u.PinHash, &u.PinChangeRequired,
		&u.IsActive, &u.AccountStatus, &u.StatusReason, &u.FraudFlag, &u.FraudReason,
		&u.SuspendedUntil, &u.StatusChangedAt,
		&u.IsEmailVerified, &u.EmailVerifiedAt, &u.LastLoginAt,
		&u.AuthInvalidBefore,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const q = `SELECT id, email, password_hash, first_name, last_name, COALESCE(phone,''),
	                  COALESCE(avatar_url,''), COALESCE(auth_provider,''), COALESCE(google_sub,''),
	                  COALESCE(pin_hash,''), COALESCE(pin_change_required, TRUE),
	                  is_active, account_status, status_reason, fraud_flag, fraud_reason,
	                  suspended_until, status_changed_at,
	                  is_email_verified, email_verified_at, last_login_at,
	                  auth_invalid_before,
	                  created_at, updated_at
	           FROM users WHERE email=$1 AND deleted_at IS NULL`
	u := &models.User{}
	err := r.db.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Phone,
		&u.AvatarURL, &u.AuthProvider, &u.GoogleSub,
		&u.PinHash, &u.PinChangeRequired,
		&u.IsActive, &u.AccountStatus, &u.StatusReason, &u.FraudFlag, &u.FraudReason,
		&u.SuspendedUntil, &u.StatusChangedAt,
		&u.IsEmailVerified, &u.EmailVerifiedAt, &u.LastLoginAt,
		&u.AuthInvalidBefore,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

func (r *UserRepo) GetByGoogleSub(ctx context.Context, googleSub string) (*models.User, error) {
	const q = `SELECT id FROM users WHERE google_sub=$1 AND deleted_at IS NULL`
	var id string
	err := r.db.QueryRow(ctx, q, googleSub).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *UserRepo) LinkGoogleIdentity(ctx context.Context, userID, googleSub, avatarURL string) error {
	const q = `UPDATE users
	           SET google_sub=$2, avatar_url=$3, auth_provider=CASE WHEN auth_provider='' THEN 'google' ELSE auth_provider END,
	               is_email_verified=TRUE, email_verified_at=COALESCE(email_verified_at, NOW()), updated_at=NOW()
	           WHERE id=$1 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, q, userID, googleSub, avatarURL)
	return err
}

func (r *UserRepo) List(ctx context.Context, p *models.PaginationParams) ([]*models.User, int64, error) {
	const countQ = `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL
	                AND ($1='' OR first_name ILIKE $1 OR last_name ILIKE $1 OR email ILIKE $1)`
	const q = `SELECT id, email, first_name, last_name, COALESCE(phone,''),
	                  is_active, account_status, status_reason, fraud_flag, fraud_reason,
	                  suspended_until, status_changed_at, is_email_verified,
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

	users := []*models.User{}
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone,
			&u.IsActive, &u.AccountStatus, &u.StatusReason, &u.FraudFlag, &u.FraudReason,
			&u.SuspendedUntil, &u.StatusChangedAt, &u.IsEmailVerified,
			&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (r *UserRepo) Update(ctx context.Context, u *models.User) error {
	const q = `UPDATE users SET first_name=$2, last_name=$3, phone=$4, avatar_url=$5, updated_at=NOW()
	           WHERE id=$1 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, q, u.ID, u.FirstName, u.LastName, u.Phone, u.AvatarURL)
	return err
}

func (r *UserRepo) UpdatePassword(ctx context.Context, userID, hash string) error {
	const q = `UPDATE users SET password_hash=$2, auth_invalid_before=date_trunc('second', NOW()), updated_at=NOW()
	           WHERE id=$1 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, q, userID, hash)
	return err
}

func (r *UserRepo) SetActive(ctx context.Context, userID string, active bool) error {
	status := models.AccountStatusSuspended
	if active {
		status = models.AccountStatusActive
	}
	const q = `UPDATE users SET is_active=$2, account_status=$3,
	                  status_reason=CASE WHEN $2 THEN '' ELSE status_reason END,
	                  suspended_until=CASE WHEN $2 THEN NULL ELSE suspended_until END,
	                  auth_invalid_before=CASE WHEN $2=FALSE THEN date_trunc('second', NOW()) ELSE auth_invalid_before END,
	                  updated_at=NOW()
	           WHERE id=$1 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, q, userID, active, status)
	return err
}

func (r *UserRepo) SetAccountManagement(ctx context.Context, userID, actorID, status, reason string, fraud bool, fraudReason string, suspendedUntil *time.Time) error {
	const q = `UPDATE users
	           SET account_status=$3, is_active=($3='ACTIVE'), status_reason=$4,
	               fraud_flag=$5, fraud_reason=$6, suspended_until=$7,
	               status_changed_by=$2, status_changed_at=NOW(),
	               auth_invalid_before=CASE WHEN $3<>'ACTIVE' THEN date_trunc('second', NOW()) ELSE auth_invalid_before END,
	               updated_at=NOW()
	           WHERE id=$1 AND deleted_at IS NULL`
	result, err := r.db.Exec(ctx, q, userID, actorID, status, reason, fraud, fraudReason, suspendedUntil)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *UserRepo) SoftDelete(ctx context.Context, userID string) error {
	const q = `UPDATE users SET deleted_at=NOW(), auth_invalid_before=date_trunc('second', NOW()), updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, userID)
	return err
}

func (r *UserRepo) GetAuthorizationState(ctx context.Context, userID string) (string, bool, *time.Time, []string, error) {
	const q = `SELECT email, is_active, auth_invalid_before
	           FROM users WHERE id=$1 AND deleted_at IS NULL`
	var email string
	var active bool
	var invalidBefore *time.Time
	if err := r.db.QueryRow(ctx, q, userID).Scan(&email, &active, &invalidBefore); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil, nil, ErrNotFound
		}
		return "", false, nil, nil, err
	}
	roles, err := r.GetRoles(ctx, userID)
	if err != nil {
		return "", false, nil, nil, err
	}
	names := make([]string, 0, len(roles))
	for _, role := range roles {
		names = append(names, role.Name)
	}
	return email, active, invalidBefore, names, nil
}

func (r *UserRepo) UpdateLastLogin(ctx context.Context, userID string) error {
	const q = `UPDATE users SET last_login_at=NOW(), updated_at=NOW() WHERE id=$1`
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

func (r *UserRepo) CountSuperAdmins(ctx context.Context) (int64, error) {
	const q = `SELECT COUNT(*) FROM user_roles ur
	           JOIN roles ro ON ro.id = ur.role_id
	           WHERE ro.name = 'super_admin'`
	var n int64
	err := r.db.QueryRow(ctx, q).Scan(&n)
	return n, err
}

// PIN methods

func (r *UserRepo) SetPinHash(ctx context.Context, userID, hash string) error {
	const q = `UPDATE users SET pin_hash=$2, pin_change_required=FALSE, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, userID, hash)
	return err
}

func (r *UserRepo) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	const q = `SELECT id, name, description, is_system, created_at FROM roles WHERE name=$1`
	role := &models.Role{}
	err := r.db.QueryRow(ctx, q, name).Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return role, err
}

func (r *UserRepo) GetPinHash(ctx context.Context, userID string) (string, error) {
	var hash *string
	err := r.db.QueryRow(ctx, `SELECT pin_hash FROM users WHERE id=$1`, userID).Scan(&hash)
	if err != nil {
		return "", err
	}
	if hash == nil {
		return "", nil
	}
	return *hash, nil
}

func (r *UserRepo) SetPinChangeRequired(ctx context.Context, userID string, required bool) error {
	const q = `UPDATE users SET pin_change_required=$2, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, userID, required)
	return err
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
	const q = `
		SELECT r.id, r.name, r.description, r.is_system, r.created_at,
		       p.id, p.name, p.description, p.resource, p.action
		FROM roles r
		LEFT JOIN role_permissions rp ON rp.role_id = r.id
		LEFT JOIN permissions p ON p.id = rp.permission_id
		WHERE r.id = $1
		ORDER BY p.resource, p.action`
	rows, err := r.db.Query(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	role := &models.Role{}
	found := false
	for rows.Next() {
		found = true
		var pID, pName, pDesc, pResource, pAction *string
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt,
			&pID, &pName, &pDesc, &pResource, &pAction); err != nil {
			return nil, err
		}
		if pID != nil {
			role.Permissions = append(role.Permissions, models.Permission{
				ID: *pID, Name: *pName, Description: *pDesc, Resource: *pResource, Action: *pAction,
			})
		}
	}
	if !found {
		return nil, ErrNotFound
	}
	return role, rows.Err()
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
	const q = `SELECT id, COALESCE(application_reference,''), user_id, form_type,
	                  COALESCE(application_type,''), COALESCE(organisation_type,''),
	                  status, payment_status, COALESCE(payment_method,''),
	                  COALESCE(payment_reference,''), payment_amount_ugx, COALESCE(signed_form_url,''),
	                  signed_form_uploaded_at, submitted_at, reviewer_id, COALESCE(review_notes,''),
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
	const q = `SELECT id, COALESCE(application_reference,''), user_id, form_type,
	                  COALESCE(application_type,''), status, payment_status, submitted_at,
	                  created_at, updated_at
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
	var total int64
	if status != "" {
		if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM applications WHERE status=$1`, status).Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM applications`).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	const baseQ = `SELECT id, COALESCE(application_reference,''), user_id, form_type,
	                      COALESCE(application_type,''), status, payment_status, submitted_at,
	                      created_at, updated_at
	               FROM applications`
	var rows pgx.Rows
	var err error
	if status != "" {
		rows, err = r.db.Query(ctx, baseQ+` WHERE status=$3 ORDER BY updated_at DESC LIMIT $1 OFFSET $2`,
			p.PerPage, p.Offset(), status)
	} else {
		rows, err = r.db.Query(ctx, baseQ+` ORDER BY updated_at DESC LIMIT $1 OFFSET $2`,
			p.PerPage, p.Offset())
	}
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
	apps := []*models.Application{}
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

type AuditExportFilter struct {
	From           *time.Time
	To             *time.Time
	Severity       string
	UserID         string
	IP             string
	EventType      string
	FailedAuthOnly bool
	Limit          int
}

const auditInsertSQL = `INSERT INTO audit_logs
	           (id, user_id, action, resource, resource_id, old_values, new_values,
	            ip_address, user_agent, method, endpoint, response_code, response_time_ms, device_info,
	            event_type, event_status, severity_level, forwarded_ip,
	            geo_country, geo_city, geo_region, geo_latitude, geo_longitude, geo_timezone, geo_source,
	            platform, authenticated, vpn_detected,
	            browser, os_name, client_type,
	            threat_score, anomaly_detected, session_id, username, payload_excerpt)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,
	                   $15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36)`

func auditValues(a *models.AuditLog) []any {
	var responseCode *int
	if a.ResponseCode != 0 {
		responseCode = &a.ResponseCode
	}
	var responseTimeMs *int64
	if a.ResponseTimeMs != 0 {
		responseTimeMs = &a.ResponseTimeMs
	}
	return []any{
		a.ID, a.UserID, a.Action, a.Resource,
		nullableStr(a.ResourceID), a.OldValues, a.NewValues,
		nullableStr(a.IPAddress), nullableStr(a.UserAgent),
		nullableStr(a.Method), nullableStr(a.Endpoint),
		responseCode, responseTimeMs,
		nullableStr(a.DeviceInfo),
		nullableStr(a.EventType), nullableStr(a.EventStatus), nullableStr(a.SeverityLevel),
		nullableStr(a.ForwardedIP),
		nullableStr(a.GeoCountry), nullableStr(a.GeoCity), nullableStr(a.GeoRegion), a.GeoLatitude, a.GeoLongitude,
		nullableStr(a.GeoTimezone), nullableStr(a.GeoSource), nullableStr(a.Platform), a.Authenticated, a.VPNDetected,
		nullableStr(a.Browser), nullableStr(a.OSName), nullableStr(a.ClientType),
		a.ThreatScore, a.AnomalyDetected, nullableStr(a.SessionID), nullableStr(a.Username), a.PayloadExcerpt,
	}
}

func (r *AuditRepo) Log(ctx context.Context, a *models.AuditLog) error {
	return r.db.QueryRow(ctx, auditInsertSQL+` RETURNING created_at`, auditValues(a)...).Scan(&a.CreatedAt)
}

func (r *AuditRepo) LogBatch(ctx context.Context, entries []*models.AuditLog) error {
	if len(entries) == 0 {
		return nil
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for _, entry := range entries {
		if _, err = tx.Exec(ctx, auditInsertSQL, auditValues(entry)...); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *AuditRepo) ListByUser(ctx context.Context, userID string, p *models.PaginationParams) ([]*models.AuditLog, int64, error) {
	const countQ = `SELECT COUNT(*) FROM audit_logs WHERE user_id=$1`
	const q = `SELECT al.id, al.user_id, al.action, al.resource, COALESCE(al.resource_id,''),
	                  COALESCE(al.ip_address,''), COALESCE(al.method,''), COALESCE(al.endpoint,''),
	                  COALESCE(al.response_code,0), COALESCE(al.response_time_ms,0),
	                  COALESCE(al.device_info,''),
	                  COALESCE(al.event_type,''), COALESCE(al.event_status,''), COALESCE(al.severity_level,''),
	                  COALESCE(al.geo_city,''), COALESCE(al.geo_country,''),
	                  COALESCE(al.browser,''), COALESCE(al.os_name,''),
	                  al.created_at
	           FROM audit_logs al
	           WHERE al.user_id=$1
	           ORDER BY al.created_at DESC LIMIT $2 OFFSET $3`
	var total int64
	if err := r.db.QueryRow(ctx, countQ, userID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, q, userID, p.PerPage, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	logs := []*models.AuditLog{}
	for rows.Next() {
		l := &models.AuditLog{}
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.Action, &l.Resource, &l.ResourceID,
			&l.IPAddress, &l.Method, &l.Endpoint, &l.ResponseCode, &l.ResponseTimeMs,
			&l.DeviceInfo,
			&l.EventType, &l.EventStatus, &l.SeverityLevel,
			&l.GeoCity, &l.GeoCountry,
			&l.Browser, &l.OSName,
			&l.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, total, rows.Err()
}

func (r *AuditRepo) List(ctx context.Context, p *models.PaginationParams) ([]*models.AuditLog, int64, error) {
	const countQ = `SELECT COUNT(*) FROM audit_logs al
	               LEFT JOIN users u ON al.user_id = u.id
	               WHERE ($1='' OR al.action ILIKE $1 OR al.resource ILIKE $1 OR
	                      COALESCE(al.ip_address,'') ILIKE $1 OR COALESCE(al.username,'') ILIKE $1 OR
	                      COALESCE(u.first_name,'') ILIKE $1 OR COALESCE(u.last_name,'') ILIKE $1 OR
	                      COALESCE(al.geo_country,'') ILIKE $1 OR COALESCE(al.event_type,'') ILIKE $1)`
	const q = `SELECT al.id, al.user_id, al.action, al.resource, COALESCE(al.resource_id,''),
	                  COALESCE(al.ip_address,''), COALESCE(al.method,''), COALESCE(al.endpoint,''),
	                  COALESCE(al.response_code,0), COALESCE(al.response_time_ms,0),
	                  COALESCE(al.device_info,''),
	                  COALESCE(al.event_type,''), COALESCE(al.event_status,''), COALESCE(al.severity_level,''),
	                  COALESCE(al.forwarded_ip,''), COALESCE(al.geo_country,''), COALESCE(al.geo_city,''),
	                  COALESCE(al.geo_region,''), al.geo_latitude, al.geo_longitude, COALESCE(al.geo_timezone,''),
	                  COALESCE(al.geo_source,''), COALESCE(al.platform,''), COALESCE(al.authenticated,false),
	                  COALESCE(al.vpn_detected,false), COALESCE(al.browser,''), COALESCE(al.os_name,''),
	                  COALESCE(al.client_type,''), COALESCE(al.threat_score,0),
	                  COALESCE(al.anomaly_detected,false), COALESCE(al.session_id,''),
	                  COALESCE(al.username,''),
	                  COALESCE(al.previous_hash,''), COALESCE(al.entry_hash,''), COALESCE(al.chain_sequence,0),
	                  COALESCE(u.first_name,''), COALESCE(u.last_name,''),
	                  al.created_at
	           FROM audit_logs al
	           LEFT JOIN users u ON al.user_id = u.id
	           WHERE ($1='' OR al.action ILIKE $1 OR al.resource ILIKE $1 OR
	                  COALESCE(al.ip_address,'') ILIKE $1 OR COALESCE(al.username,'') ILIKE $1 OR
	                  COALESCE(u.first_name,'') ILIKE $1 OR COALESCE(u.last_name,'') ILIKE $1 OR
	                  COALESCE(al.geo_country,'') ILIKE $1 OR COALESCE(al.event_type,'') ILIKE $1)
	           ORDER BY al.created_at DESC LIMIT $2 OFFSET $3`
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
	logs := []*models.AuditLog{}
	for rows.Next() {
		l := &models.AuditLog{}
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.Action, &l.Resource, &l.ResourceID,
			&l.IPAddress, &l.Method, &l.Endpoint, &l.ResponseCode, &l.ResponseTimeMs,
			&l.DeviceInfo,
			&l.EventType, &l.EventStatus, &l.SeverityLevel,
			&l.ForwardedIP, &l.GeoCountry, &l.GeoCity,
			&l.GeoRegion, &l.GeoLatitude, &l.GeoLongitude, &l.GeoTimezone,
			&l.GeoSource, &l.Platform, &l.Authenticated,
			&l.VPNDetected, &l.Browser, &l.OSName,
			&l.ClientType, &l.ThreatScore,
			&l.AnomalyDetected, &l.SessionID, &l.Username,
			&l.PreviousHash, &l.EntryHash, &l.ChainSequence,
			&l.FirstName, &l.LastName,
			&l.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, total, rows.Err()
}

func (r *AuditRepo) GetByID(ctx context.Context, id string) (*models.AuditLog, error) {
	const q = `SELECT al.id, al.user_id, al.action, al.resource, COALESCE(al.resource_id,''),
	                  al.old_values, al.new_values, COALESCE(al.ip_address,''), COALESCE(al.user_agent,''),
	                  COALESCE(al.method,''), COALESCE(al.endpoint,''),
	                  COALESCE(al.response_code,0), COALESCE(al.response_time_ms,0),
	                  COALESCE(al.device_info,''),
	                  COALESCE(al.event_type,''), COALESCE(al.event_status,''), COALESCE(al.severity_level,''),
	                  COALESCE(al.forwarded_ip,''), COALESCE(al.geo_country,''), COALESCE(al.geo_city,''),
	                  COALESCE(al.geo_region,''), al.geo_latitude, al.geo_longitude, COALESCE(al.geo_timezone,''),
	                  COALESCE(al.geo_source,''), COALESCE(al.platform,''), COALESCE(al.authenticated,false),
	                  COALESCE(al.vpn_detected,false), COALESCE(al.browser,''), COALESCE(al.os_name,''),
	                  COALESCE(al.client_type,''), COALESCE(al.threat_score,0),
	                  COALESCE(al.anomaly_detected,false), COALESCE(al.session_id,''),
	                  COALESCE(al.username,''),
	                  al.payload_excerpt, COALESCE(al.previous_hash,''), COALESCE(al.entry_hash,''), COALESCE(al.chain_sequence,0),
	                  COALESCE(u.first_name,''), COALESCE(u.last_name,''),
	                  al.created_at
	           FROM audit_logs al
	           LEFT JOIN users u ON al.user_id = u.id
	           WHERE al.id=$1`
	l := &models.AuditLog{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&l.ID, &l.UserID, &l.Action, &l.Resource, &l.ResourceID,
		&l.OldValues, &l.NewValues, &l.IPAddress, &l.UserAgent,
		&l.Method, &l.Endpoint, &l.ResponseCode, &l.ResponseTimeMs,
		&l.DeviceInfo,
		&l.EventType, &l.EventStatus, &l.SeverityLevel,
		&l.ForwardedIP, &l.GeoCountry, &l.GeoCity,
		&l.GeoRegion, &l.GeoLatitude, &l.GeoLongitude, &l.GeoTimezone,
		&l.GeoSource, &l.Platform, &l.Authenticated,
		&l.VPNDetected, &l.Browser, &l.OSName,
		&l.ClientType, &l.ThreatScore,
		&l.AnomalyDetected, &l.SessionID, &l.Username,
		&l.PayloadExcerpt, &l.PreviousHash, &l.EntryHash, &l.ChainSequence,
		&l.FirstName, &l.LastName,
		&l.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return l, err
}

func (r *AuditRepo) Export(ctx context.Context, f AuditExportFilter) ([]*models.AuditLog, error) {
	if f.Limit < 1 || f.Limit > 10000 {
		f.Limit = 5000
	}
	const q = `SELECT id,user_id,action,resource,COALESCE(ip_address,''),COALESCE(method,''),COALESCE(endpoint,''),COALESCE(response_code,0),COALESCE(event_type,''),COALESCE(event_status,''),COALESCE(severity_level,''),COALESCE(geo_country,''),COALESCE(device_info,''),COALESCE(username,''),COALESCE(entry_hash,''),created_at FROM audit_logs WHERE ($1::timestamptz IS NULL OR created_at >= $1) AND ($2::timestamptz IS NULL OR created_at <= $2) AND ($3='' OR severity_level=$3) AND ($4='' OR user_id=$4) AND ($5='' OR ip_address=$5) AND ($6='' OR event_type=$6) AND ($7=FALSE OR (event_type LIKE 'AUTH_%' AND event_status='FAILURE')) ORDER BY created_at DESC LIMIT $8`
	rows, err := r.db.Query(ctx, q, f.From, f.To, f.Severity, f.UserID, f.IP, f.EventType, f.FailedAuthOnly, f.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.AuditLog{}
	for rows.Next() {
		l := &models.AuditLog{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.Resource, &l.IPAddress, &l.Method, &l.Endpoint, &l.ResponseCode, &l.EventType, &l.EventStatus, &l.SeverityLevel, &l.GeoCountry, &l.DeviceInfo, &l.Username, &l.EntryHash, &l.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, l)
	}
	return items, rows.Err()
}

// ── CMS Repository ────────────────────────────────────────────────

type CMSRepo struct{ db *pgxpool.Pool }

// Posts

func (r *CMSRepo) CreatePost(ctx context.Context, p *models.CMSPost) error {
	const q = `INSERT INTO cms_posts (id, title, slug, content, excerpt, category, category_tag, status, cover_image_url, author_id, published_at, meta_title, meta_description, focus_keywords, approved_at)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q,
		p.ID, p.Title, p.Slug, p.Content, p.Excerpt, p.Category, p.CategoryTag, p.Status,
		p.CoverImageURL, p.AuthorID, p.PublishedAt, p.MetaTitle, p.MetaDescription,
		p.FocusKeywords, p.ApprovedAt,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *CMSRepo) GetPostBySlug(ctx context.Context, slug string) (*models.CMSPost, error) {
	_, _ = r.db.Exec(ctx, `UPDATE cms_posts SET view_count = COALESCE(view_count,0)+1 WHERE slug=$1`, slug)
	const q = `SELECT p.id, p.title, p.slug, p.content, p.excerpt, p.category, COALESCE(p.category_tag,''), p.status,
	                  COALESCE(p.cover_image_url,''), p.author_id,
	                  COALESCE(u.first_name||' '||u.last_name,'') AS author_name,
	                  COALESCE(p.meta_title,''), COALESCE(p.meta_description,''),
	                  COALESCE(p.focus_keywords,''), COALESCE(p.view_count,0),
	                  p.approved_at, p.published_at, p.created_at, p.updated_at
	           FROM cms_posts p
	           LEFT JOIN users u ON u.id = p.author_id
	           WHERE p.slug=$1`
	p := &models.CMSPost{}
	err := r.db.QueryRow(ctx, q, slug).Scan(
		&p.ID, &p.Title, &p.Slug, &p.Content, &p.Excerpt, &p.Category, &p.CategoryTag, &p.Status,
		&p.CoverImageURL, &p.AuthorID, &p.AuthorName, &p.MetaTitle, &p.MetaDescription,
		&p.FocusKeywords, &p.ViewCount, &p.ApprovedAt, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

func (r *CMSRepo) GetPostByID(ctx context.Context, id string) (*models.CMSPost, error) {
	const q = `SELECT p.id, p.title, p.slug, p.content, p.excerpt, p.category, COALESCE(p.category_tag,''), p.status,
	                  COALESCE(p.cover_image_url,''), p.author_id,
	                  COALESCE(u.first_name||' '||u.last_name,'') AS author_name,
	                  COALESCE(p.meta_title,''), COALESCE(p.meta_description,''),
	                  COALESCE(p.focus_keywords,''), COALESCE(p.view_count,0),
	                  p.approved_at, p.published_at, p.created_at, p.updated_at
	           FROM cms_posts p
	           LEFT JOIN users u ON u.id = p.author_id
	           WHERE p.id=$1`
	p := &models.CMSPost{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&p.ID, &p.Title, &p.Slug, &p.Content, &p.Excerpt, &p.Category, &p.CategoryTag, &p.Status,
		&p.CoverImageURL, &p.AuthorID, &p.AuthorName, &p.MetaTitle, &p.MetaDescription,
		&p.FocusKeywords, &p.ViewCount, &p.ApprovedAt, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return p, err
}

func (r *CMSRepo) ListPosts(ctx context.Context, category, status string, limit, offset int) ([]*models.CMSPost, int64, error) {
	const countQ = `SELECT COUNT(*) FROM cms_posts p
	                WHERE ($1='' OR p.category=$1) AND ($2='' OR p.status=$2)`
	const q = `SELECT p.id, p.title, p.slug, p.excerpt, p.category, COALESCE(p.category_tag,''), p.status,
	                  COALESCE(p.cover_image_url,''), p.author_id,
	                  COALESCE(u.first_name||' '||u.last_name,'') AS author_name,
	                  COALESCE(p.meta_title,''), COALESCE(p.meta_description,''),
	                  COALESCE(p.focus_keywords,''), COALESCE(p.view_count,0),
	                  p.approved_at, p.published_at, p.created_at, p.updated_at
	           FROM cms_posts p
	           LEFT JOIN users u ON u.id = p.author_id
	           WHERE ($1='' OR p.category=$1) AND ($2='' OR p.status=$2)
	           ORDER BY COALESCE(p.published_at, p.created_at) DESC LIMIT $3 OFFSET $4`
	var total int64
	if err := r.db.QueryRow(ctx, countQ, category, status).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, q, category, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	posts := []*models.CMSPost{}
	for rows.Next() {
		p := &models.CMSPost{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Category, &p.CategoryTag, &p.Status,
			&p.CoverImageURL, &p.AuthorID, &p.AuthorName, &p.MetaTitle, &p.MetaDescription,
			&p.FocusKeywords, &p.ViewCount, &p.ApprovedAt, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}
	return posts, total, rows.Err()
}

func (r *CMSRepo) UpdatePost(ctx context.Context, p *models.CMSPost) error {
	const q = `UPDATE cms_posts SET title=$2, slug=$3, content=$4, excerpt=$5, category=$6, category_tag=$14,
	           status=$7, cover_image_url=$8, published_at=$9,
	           meta_title=$10, meta_description=$11, focus_keywords=$12, approved_at=$13, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, p.ID, p.Title, p.Slug, p.Content, p.Excerpt,
		p.Category, p.Status, p.CoverImageURL, p.PublishedAt, p.MetaTitle, p.MetaDescription, p.FocusKeywords, p.ApprovedAt, p.CategoryTag)
	return err
}

func (r *CMSRepo) DeletePost(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_posts WHERE id=$1`, id)
	return err
}

func (r *CMSRepo) ListBlogCategories(ctx context.Context, activeOnly bool, contentType string) ([]*models.BlogCategory, error) {
	q := `SELECT id, name, slug, COALESCE(description,''), content_type, sort_order, is_active, created_at, updated_at
	      FROM blog_categories WHERE ($1=FALSE OR is_active=TRUE) AND content_type=$2
	      ORDER BY sort_order, name`
	rows, err := r.db.Query(ctx, q, activeOnly, contentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*models.BlogCategory{}
	for rows.Next() {
		c := &models.BlogCategory{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ContentType, &c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *CMSRepo) CreateBlogCategory(ctx context.Context, c *models.BlogCategory) error {
	const q = `INSERT INTO blog_categories (id, name, slug, description, content_type, sort_order, is_active)
	           VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING created_at, updated_at`
	err := r.db.QueryRow(ctx, q, c.ID, c.Name, c.Slug, c.Description, c.ContentType, c.SortOrder, c.IsActive).Scan(&c.CreatedAt, &c.UpdatedAt)
	if err != nil && isDuplicate(err) {
		return ErrDuplicate
	}
	return err
}

func (r *CMSRepo) UpdateBlogCategory(ctx context.Context, c *models.BlogCategory) error {
	const q = `UPDATE blog_categories SET name=$2, slug=$3, description=$4, content_type=$5, sort_order=$6, is_active=$7, updated_at=NOW()
	           WHERE id=$1`
	tag, err := r.db.Exec(ctx, q, c.ID, c.Name, c.Slug, c.Description, c.ContentType, c.SortOrder, c.IsActive)
	if err != nil {
		if isDuplicate(err) {
			return ErrDuplicate
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CMSRepo) DeleteBlogCategory(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM blog_categories WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CMSRepo) GetCommentByID(ctx context.Context, id string) (*models.BlogComment, error) {
	const q = `SELECT id, post_id, user_id, parent_id, body, status, depth, COALESCE(ip_address,''), COALESCE(user_agent,''),
	                  created_at, updated_at, approved_at, flagged_at
	           FROM blog_comments WHERE id=$1`
	c := &models.BlogComment{}
	err := r.db.QueryRow(ctx, q, id).Scan(&c.ID, &c.PostID, &c.UserID, &c.ParentID, &c.Body, &c.Status, &c.Depth, &c.IPAddress, &c.UserAgent, &c.CreatedAt, &c.UpdatedAt, &c.ApprovedAt, &c.FlaggedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return c, err
}

func (r *CMSRepo) ListApprovedCommentsForPost(ctx context.Context, postID string) ([]*models.BlogComment, error) {
	const q = `SELECT c.id, c.post_id, c.user_id, COALESCE(u.first_name||' '||u.last_name,'Reader') AS user_name,
	                  c.parent_id, c.body, c.status, c.depth, c.created_at, c.updated_at, c.approved_at, c.flagged_at
	           FROM blog_comments c
	           LEFT JOIN users u ON u.id = c.user_id
	           WHERE c.post_id=$1 AND c.status='approved'
	           ORDER BY c.created_at ASC`
	rows, err := r.db.Query(ctx, q, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*models.BlogComment{}
	for rows.Next() {
		c := &models.BlogComment{}
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.UserName, &c.ParentID, &c.Body, &c.Status, &c.Depth, &c.CreatedAt, &c.UpdatedAt, &c.ApprovedAt, &c.FlaggedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *CMSRepo) CreateComment(ctx context.Context, c *models.BlogComment) error {
	const q = `INSERT INTO blog_comments (id, post_id, user_id, parent_id, body, status, depth, ip_address, user_agent)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q, c.ID, c.PostID, c.UserID, c.ParentID, c.Body, c.Status, c.Depth, c.IPAddress, c.UserAgent).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *CMSRepo) ListComments(ctx context.Context, status string, limit, offset int) ([]*models.BlogComment, int64, error) {
	const countQ = `SELECT COUNT(*) FROM blog_comments WHERE ($1='' OR status=$1)`
	const q = `SELECT c.id, c.post_id, COALESCE(p.title,''), COALESCE(p.slug,''), c.user_id,
	                  COALESCE(u.first_name||' '||u.last_name,'Reader'), c.parent_id, c.body,
	                  c.status, c.depth, COALESCE(c.ip_address,''), COALESCE(c.user_agent,''),
	                  c.created_at, c.updated_at, c.approved_at, c.flagged_at
	           FROM blog_comments c
	           LEFT JOIN cms_posts p ON p.id = c.post_id
	           LEFT JOIN users u ON u.id = c.user_id
	           WHERE ($1='' OR c.status=$1)
	           ORDER BY c.created_at DESC LIMIT $2 OFFSET $3`
	var total int64
	if err := r.db.QueryRow(ctx, countQ, status).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, q, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []*models.BlogComment{}
	for rows.Next() {
		c := &models.BlogComment{}
		if err := rows.Scan(&c.ID, &c.PostID, &c.PostTitle, &c.PostSlug, &c.UserID, &c.UserName, &c.ParentID, &c.Body, &c.Status, &c.Depth, &c.IPAddress, &c.UserAgent, &c.CreatedAt, &c.UpdatedAt, &c.ApprovedAt, &c.FlaggedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, c)
	}
	return out, total, rows.Err()
}

func (r *CMSRepo) ModerateComment(ctx context.Context, id, status string) error {
	const q = `UPDATE blog_comments
	           SET status=$2,
	               approved_at=CASE WHEN $2='approved' THEN NOW() ELSE approved_at END,
	               flagged_at=CASE WHEN $2='flagged' THEN NOW() ELSE flagged_at END,
	               updated_at=NOW()
	           WHERE id=$1`
	tag, err := r.db.Exec(ctx, q, id, status)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CMSRepo) DeleteComment(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM blog_comments WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Events

func (r *CMSRepo) CreateEvent(ctx context.Context, e *models.CMSEvent) error {
	const q = `INSERT INTO cms_events (id, title, slug, description, category, location, event_date, end_date, cover_image_url, status, author_id)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q,
		e.ID, e.Title, e.Slug, e.Description, e.Category, e.Location, e.EventDate, e.EndDate,
		e.CoverImageURL, e.Status, e.AuthorID,
	).Scan(&e.CreatedAt, &e.UpdatedAt)
}

func (r *CMSRepo) GetEventBySlug(ctx context.Context, slug string) (*models.CMSEvent, error) {
	const q = `SELECT id, title, slug, description, COALESCE(category,''), COALESCE(location,''), event_date, end_date,
	                  COALESCE(cover_image_url,''), status, author_id, created_at, updated_at
	           FROM cms_events WHERE slug=$1`
	e := &models.CMSEvent{}
	err := r.db.QueryRow(ctx, q, slug).Scan(
		&e.ID, &e.Title, &e.Slug, &e.Description, &e.Category, &e.Location, &e.EventDate, &e.EndDate,
		&e.CoverImageURL, &e.Status, &e.AuthorID, &e.CreatedAt, &e.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return e, err
}

func (r *CMSRepo) GetEventByID(ctx context.Context, id string) (*models.CMSEvent, error) {
	const q = `SELECT id, title, slug, description, COALESCE(category,''), COALESCE(location,''), event_date, end_date,
	                  COALESCE(cover_image_url,''), status, author_id, created_at, updated_at
	           FROM cms_events WHERE id=$1`
	e := &models.CMSEvent{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&e.ID, &e.Title, &e.Slug, &e.Description, &e.Category, &e.Location, &e.EventDate, &e.EndDate,
		&e.CoverImageURL, &e.Status, &e.AuthorID, &e.CreatedAt, &e.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return e, err
}

func (r *CMSRepo) ListEvents(ctx context.Context, status string, limit, offset int) ([]*models.CMSEvent, int64, error) {
	const countQ = `SELECT COUNT(*) FROM cms_events WHERE ($1='' OR status=$1)`
	const q = `SELECT id, title, slug, COALESCE(category,''), COALESCE(location,''), event_date, end_date,
	                  COALESCE(cover_image_url,''), status, author_id, created_at, updated_at
	           FROM cms_events WHERE ($1='' OR status=$1)
	           ORDER BY COALESCE(event_date, created_at) DESC LIMIT $2 OFFSET $3`
	var total int64
	if err := r.db.QueryRow(ctx, countQ, status).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, q, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	events := []*models.CMSEvent{}
	for rows.Next() {
		e := &models.CMSEvent{}
		if err := rows.Scan(&e.ID, &e.Title, &e.Slug, &e.Category, &e.Location, &e.EventDate, &e.EndDate,
			&e.CoverImageURL, &e.Status, &e.AuthorID, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, 0, err
		}
		events = append(events, e)
	}
	return events, total, rows.Err()
}

func (r *CMSRepo) UpdateEvent(ctx context.Context, e *models.CMSEvent) error {
	const q = `UPDATE cms_events SET title=$2, slug=$3, description=$4, category=$5, location=$6,
	           event_date=$7, end_date=$8, cover_image_url=$9, status=$10, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, e.ID, e.Title, e.Slug, e.Description, e.Category, e.Location,
		e.EventDate, e.EndDate, e.CoverImageURL, e.Status)
	return err
}

func (r *CMSRepo) DeleteEvent(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_events WHERE id=$1`, id)
	return err
}

// Careers

func (r *CMSRepo) CreateCareer(ctx context.Context, c *models.CMSCareer) error {
	const q = `INSERT INTO cms_careers (id, title, department, department_id, location, job_type, category, description, requirements, salary_range, status, deadline_at, author_id)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q,
		c.ID, c.Title, c.Department, c.DepartmentID, c.Location, c.JobType, c.Category, c.Description,
		c.Requirements, c.SalaryRange, c.Status, c.DeadlineAt, c.AuthorID,
	).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *CMSRepo) GetCareerByID(ctx context.Context, id string) (*models.CMSCareer, error) {
	const q = `SELECT c.id, c.title, COALESCE(c.department,''), c.department_id, COALESCE(d.name,''),
	                  COALESCE(c.location,''), c.job_type,
	                  COALESCE(c.category,'jobs'), c.description, COALESCE(c.requirements,''), COALESCE(c.salary_range,''),
	                  c.status, c.deadline_at, c.author_id, c.created_at, c.updated_at
	           FROM cms_careers c
	           LEFT JOIN departments d ON d.id = c.department_id
	           WHERE c.id=$1`
	c := &models.CMSCareer{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&c.ID, &c.Title, &c.Department, &c.DepartmentID, &c.DepartmentName,
		&c.Location, &c.JobType, &c.Category, &c.Description,
		&c.Requirements, &c.SalaryRange, &c.Status, &c.DeadlineAt, &c.AuthorID,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return c, err
}

func (r *CMSRepo) ListCareers(ctx context.Context, status, category string, limit, offset int) ([]*models.CMSCareer, int64, error) {
	const countQ = `SELECT COUNT(*) FROM cms_careers WHERE ($1='' OR status=$1) AND ($2='' OR COALESCE(category,'jobs')=$2)`
	const q = `SELECT c.id, c.title, COALESCE(c.department,''), c.department_id, COALESCE(d.name,''),
	                  COALESCE(c.location,''), c.job_type,
	                  COALESCE(c.category,'jobs'), COALESCE(c.salary_range,''), c.status, c.deadline_at, c.author_id, c.created_at, c.updated_at
	           FROM cms_careers c
	           LEFT JOIN departments d ON d.id = c.department_id
	           WHERE ($1='' OR c.status=$1) AND ($2='' OR COALESCE(c.category,'jobs')=$2)
	           ORDER BY c.created_at DESC LIMIT $3 OFFSET $4`
	var total int64
	if err := r.db.QueryRow(ctx, countQ, status, category).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, q, status, category, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	careers := []*models.CMSCareer{}
	for rows.Next() {
		c := &models.CMSCareer{}
		if err := rows.Scan(&c.ID, &c.Title, &c.Department, &c.DepartmentID, &c.DepartmentName,
			&c.Location, &c.JobType, &c.Category, &c.SalaryRange, &c.Status, &c.DeadlineAt, &c.AuthorID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		careers = append(careers, c)
	}
	return careers, total, rows.Err()
}

func (r *CMSRepo) UpdateCareer(ctx context.Context, c *models.CMSCareer) error {
	const q = `UPDATE cms_careers SET title=$2, department=$3, department_id=$4, location=$5, job_type=$6,
	           category=$7, description=$8, requirements=$9, salary_range=$10, status=$11,
	           deadline_at=$12, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, c.ID, c.Title, c.Department, c.DepartmentID, c.Location, c.JobType,
		c.Category, c.Description, c.Requirements, c.SalaryRange, c.Status, c.DeadlineAt)
	return err
}

func (r *CMSRepo) DeleteCareer(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_careers WHERE id=$1`, id)
	return err
}

// ── Slides ────────────────────────────────────────────────────────

func (r *CMSRepo) ListSlides(ctx context.Context, activeOnly bool) ([]*models.CMSSlide, error) {
	q := `SELECT id, slideshow_id, title, COALESCE(subtitle,''), COALESCE(description,''),
	             COALESCE(image_url,''), media_type, animation_type,
	             COALESCE(button_text,''), COALESCE(button_url,''),
	             sort_order, is_active, created_at, updated_at
	      FROM cms_slides WHERE slideshow_id='homepage-hero'`
	if activeOnly {
		q += ` AND is_active=true`
	}
	q += ` ORDER BY sort_order ASC, created_at ASC`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	slides := []*models.CMSSlide{}
	for rows.Next() {
		s := &models.CMSSlide{}
		if err := rows.Scan(&s.ID, &s.SlideshowID, &s.Title, &s.Subtitle, &s.Description,
			&s.ImageURL, &s.MediaType, &s.AnimationType, &s.ButtonText, &s.ButtonURL,
			&s.SortOrder, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		slides = append(slides, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return r.attachSlideButtons(ctx, slides)
}

func (r *CMSRepo) CreateSlide(ctx context.Context, s *models.CMSSlide) error {
	const q = `INSERT INTO cms_slides (id, slideshow_id, title, subtitle, description, image_url, media_type, animation_type, button_text, button_url, sort_order, is_active)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q,
		s.ID, s.SlideshowID, s.Title, s.Subtitle, s.Description, s.ImageURL, s.MediaType, s.AnimationType,
		s.ButtonText, s.ButtonURL, s.SortOrder, s.IsActive,
	).Scan(&s.CreatedAt, &s.UpdatedAt)
}

func (r *CMSRepo) UpdateSlide(ctx context.Context, s *models.CMSSlide) error {
	const q = `UPDATE cms_slides SET title=$2, subtitle=$3, description=$4, image_url=$5,
	           button_text=$6, button_url=$7, sort_order=$8, is_active=$9, media_type=$10, animation_type=$11, updated_at=NOW()
	           WHERE id=$1`
	_, err := r.db.Exec(ctx, q, s.ID, s.Title, s.Subtitle, s.Description, s.ImageURL,
		s.ButtonText, s.ButtonURL, s.SortOrder, s.IsActive, s.MediaType, s.AnimationType)
	return err
}

func (r *CMSRepo) DeleteSlide(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_slides WHERE id=$1`, id)
	return err
}

func (r *CMSRepo) GetSlideshowBySlug(ctx context.Context, slug string, activeOnly bool) (*models.CMSSlideshow, error) {
	const q = `SELECT id, name, slug, transition_effect, transition_duration, autoplay_speed, pause_on_hover, is_active, created_at, updated_at
	           FROM cms_slideshows WHERE slug=$1 AND ($2=FALSE OR is_active=TRUE)`
	show := &models.CMSSlideshow{}
	if err := r.db.QueryRow(ctx, q, slug, activeOnly).Scan(&show.ID, &show.Name, &show.Slug, &show.TransitionEffect, &show.TransitionDuration, &show.AutoplaySpeed, &show.PauseOnHover, &show.IsActive, &show.CreatedAt, &show.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	slides, err := r.ListSlidesBySlideshow(ctx, show.ID, activeOnly)
	if err != nil {
		return nil, err
	}
	show.Slides = slides
	return show, nil
}

func (r *CMSRepo) UpdateSlideshow(ctx context.Context, s *models.CMSSlideshow) error {
	const q = `UPDATE cms_slideshows SET name=$2, transition_effect=$3, transition_duration=$4,
	           autoplay_speed=$5, pause_on_hover=$6, is_active=$7, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, s.ID, s.Name, s.TransitionEffect, s.TransitionDuration, s.AutoplaySpeed, s.PauseOnHover, s.IsActive)
	return err
}

func (r *CMSRepo) ListSlidesBySlideshow(ctx context.Context, slideshowID string, activeOnly bool) ([]*models.CMSSlide, error) {
	q := `SELECT id, slideshow_id, title, COALESCE(subtitle,''), COALESCE(description,''),
	             COALESCE(image_url,''), media_type, animation_type,
	             COALESCE(button_text,''), COALESCE(button_url,''), sort_order, is_active, created_at, updated_at
	      FROM cms_slides WHERE slideshow_id=$1`
	if activeOnly {
		q += ` AND is_active=true`
	}
	q += ` ORDER BY sort_order ASC, created_at ASC`
	rows, err := r.db.Query(ctx, q, slideshowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	slides := []*models.CMSSlide{}
	for rows.Next() {
		s := &models.CMSSlide{}
		if err := rows.Scan(&s.ID, &s.SlideshowID, &s.Title, &s.Subtitle, &s.Description, &s.ImageURL, &s.MediaType, &s.AnimationType, &s.ButtonText, &s.ButtonURL, &s.SortOrder, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		slides = append(slides, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return r.attachSlideButtons(ctx, slides)
}

func (r *CMSRepo) ReplaceSlideButtons(ctx context.Context, slideID string, buttons []*models.CMSSlideButton) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `DELETE FROM cms_slide_buttons WHERE slide_id=$1`, slideID); err != nil {
		return err
	}
	for i, b := range buttons {
		if i >= 2 {
			break
		}
		if b.ID == "" {
			b.ID = slideID + "-button-" + string(rune('0'+i))
		}
		const q = `INSERT INTO cms_slide_buttons (id, slide_id, text, url, style_class, link_target, sort_order)
		           VALUES ($1,$2,$3,$4,$5,$6,$7)`
		if _, err := tx.Exec(ctx, q, b.ID, slideID, b.Text, b.URL, b.StyleClass, b.LinkTarget, i); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *CMSRepo) attachSlideButtons(ctx context.Context, slides []*models.CMSSlide) ([]*models.CMSSlide, error) {
	if len(slides) == 0 {
		return slides, nil
	}
	ids := make([]string, 0, len(slides))
	byID := map[string]*models.CMSSlide{}
	for _, s := range slides {
		ids = append(ids, s.ID)
		byID[s.ID] = s
	}
	rows, err := r.db.Query(ctx, `SELECT id, slide_id, text, url, style_class, link_target, sort_order, created_at, updated_at FROM cms_slide_buttons WHERE slide_id = ANY($1) ORDER BY sort_order`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		b := &models.CMSSlideButton{}
		if err := rows.Scan(&b.ID, &b.SlideID, &b.Text, &b.URL, &b.StyleClass, &b.LinkTarget, &b.SortOrder, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		if s := byID[b.SlideID]; s != nil {
			s.Buttons = append(s.Buttons, b)
		}
	}
	return slides, rows.Err()
}

// ── Menus ─────────────────────────────────────────────────────────

func (r *CMSRepo) GetMenu(ctx context.Context, name string) (*models.CMSMenu, error) {
	const q = `SELECT name, items, updated_at FROM cms_menus WHERE name=$1`
	m := &models.CMSMenu{}
	var itemsJSON []byte
	err := r.db.QueryRow(ctx, q, name).Scan(&m.Name, &itemsJSON, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(itemsJSON, &m.Items)
	return m, nil
}

func (r *CMSRepo) UpdateMenu(ctx context.Context, name string, items []models.CMSMenuItem) error {
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return err
	}
	const q = `INSERT INTO cms_menus (name, items, updated_at) VALUES ($1,$2,NOW())
	           ON CONFLICT (name) DO UPDATE SET items=$2, updated_at=NOW()`
	_, err = r.db.Exec(ctx, q, name, itemsJSON)
	return err
}

// ── Settings (generic key/value JSONB store) ──────────────────────

func (r *CMSRepo) GetSetting(ctx context.Context, key string) (*models.CMSSetting, error) {
	var s models.CMSSetting
	var raw []byte
	const q = `SELECT key, value, updated_at FROM cms_settings WHERE key=$1`
	err := r.db.QueryRow(ctx, q, key).Scan(&s.Key, &raw, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	s.Value = raw
	return &s, nil
}

func (r *CMSRepo) UpdateSetting(ctx context.Context, key string, value []byte) error {
	const q = `INSERT INTO cms_settings (key, value, updated_at) VALUES ($1,$2,NOW())
	           ON CONFLICT (key) DO UPDATE SET value=$2, updated_at=NOW()`
	_, err := r.db.Exec(ctx, q, key, value)
	return err
}

// ── Custom Fonts ────────────────────────────────────────────────────

func (r *CMSRepo) ListCustomFonts(ctx context.Context) ([]*models.CMSCustomFont, error) {
	const q = `SELECT id, font_name, display_name, file_url, font_format, created_at
	           FROM cms_custom_fonts ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.CMSCustomFont{}
	for rows.Next() {
		f := &models.CMSCustomFont{}
		if err := rows.Scan(&f.ID, &f.FontName, &f.DisplayName, &f.FileURL, &f.FontFormat, &f.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	return items, rows.Err()
}

func (r *CMSRepo) CreateCustomFont(ctx context.Context, f *models.CMSCustomFont) error {
	const q = `INSERT INTO cms_custom_fonts (id, font_name, display_name, file_url, font_format)
	           VALUES ($1,$2,$3,$4,$5) RETURNING created_at`
	err := r.db.QueryRow(ctx, q, f.ID, f.FontName, f.DisplayName, f.FileURL, f.FontFormat).Scan(&f.CreatedAt)
	if err != nil && isDuplicate(err) {
		return ErrDuplicate
	}
	return err
}

func (r *CMSRepo) DeleteCustomFont(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_custom_fonts WHERE id=$1`, id)
	return err
}

// ── Fun Facts ─────────────────────────────────────────────────────

func (r *CMSRepo) ListFunFacts(ctx context.Context, activeOnly bool) ([]*models.CMSFunFact, error) {
	q := `SELECT id, label, value, COALESCE(icon,''), sort_order, is_active, created_at, updated_at
	      FROM cms_fun_facts`
	if activeOnly {
		q += ` WHERE is_active = true`
	}
	q += ` ORDER BY sort_order ASC, created_at ASC`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.CMSFunFact{}
	for rows.Next() {
		f := &models.CMSFunFact{}
		if err := rows.Scan(&f.ID, &f.Label, &f.Value, &f.Icon, &f.SortOrder, &f.IsActive, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	return items, rows.Err()
}

func (r *CMSRepo) CreateFunFact(ctx context.Context, f *models.CMSFunFact) error {
	const q = `INSERT INTO cms_fun_facts (id, label, value, icon, sort_order, is_active)
	           VALUES ($1,$2,$3,$4,$5,$6) RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q, f.ID, f.Label, f.Value, f.Icon, f.SortOrder, f.IsActive).Scan(&f.CreatedAt, &f.UpdatedAt)
}

func (r *CMSRepo) UpdateFunFact(ctx context.Context, f *models.CMSFunFact) error {
	const q = `UPDATE cms_fun_facts SET label=$2, value=$3, icon=$4, sort_order=$5, is_active=$6, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, f.ID, f.Label, f.Value, f.Icon, f.SortOrder, f.IsActive)
	return err
}

func (r *CMSRepo) DeleteFunFact(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_fun_facts WHERE id=$1`, id)
	return err
}

// ── FAQs ──────────────────────────────────────────────────────────

func (r *CMSRepo) ListFAQs(ctx context.Context, category string) ([]*models.CMSFAQ, error) {
	const q = `SELECT id, question, answer, category, sort_order, is_active, created_at, updated_at
	           FROM cms_faqs WHERE ($1='' OR category=$1)
	           ORDER BY sort_order ASC, created_at ASC`
	rows, err := r.db.Query(ctx, q, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.CMSFAQ{}
	for rows.Next() {
		f := &models.CMSFAQ{}
		if err := rows.Scan(&f.ID, &f.Question, &f.Answer, &f.Category, &f.SortOrder, &f.IsActive, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	return items, rows.Err()
}

func (r *CMSRepo) CreateFAQ(ctx context.Context, f *models.CMSFAQ) error {
	const q = `INSERT INTO cms_faqs (id, question, answer, category, sort_order, is_active)
	           VALUES ($1,$2,$3,$4,$5,$6) RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q, f.ID, f.Question, f.Answer, f.Category, f.SortOrder, f.IsActive).Scan(&f.CreatedAt, &f.UpdatedAt)
}

func (r *CMSRepo) UpdateFAQ(ctx context.Context, f *models.CMSFAQ) error {
	const q = `UPDATE cms_faqs SET question=$2, answer=$3, category=$4, sort_order=$5, is_active=$6, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, f.ID, f.Question, f.Answer, f.Category, f.SortOrder, f.IsActive)
	return err
}

func (r *CMSRepo) DeleteFAQ(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_faqs WHERE id=$1`, id)
	return err
}

// ── Newsletter subscribers ────────────────────────────────────────

// SubscribeNewsletter upserts an email. A repeat sign-up re-activates a
// previously unsubscribed address rather than erroring on the UNIQUE(email).
func (r *CMSRepo) SubscribeNewsletter(ctx context.Context, id, email, source string) (*models.NewsletterSubscriber, error) {
	const q = `INSERT INTO newsletter_subscribers (id, email, source)
	           VALUES ($1,$2,$3)
	           ON CONFLICT (email) DO UPDATE
	             SET is_active = TRUE,
	                 unsubscribed_at = NULL,
	                 source = CASE WHEN EXCLUDED.source <> '' THEN EXCLUDED.source ELSE newsletter_subscribers.source END,
	                 updated_at = NOW()
	           RETURNING id, email, source, is_active, unsubscribed_at, created_at, updated_at`
	s := &models.NewsletterSubscriber{}
	err := r.db.QueryRow(ctx, q, id, email, source).Scan(
		&s.ID, &s.Email, &s.Source, &s.IsActive, &s.UnsubscribedAt, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

// ListNewsletterSubscribers returns subscribers newest-first plus the total
// count. A non-positive limit returns every row (used by the CSV export).
func (r *CMSRepo) ListNewsletterSubscribers(ctx context.Context, limit, offset int) ([]*models.NewsletterSubscriber, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM newsletter_subscribers`).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := `SELECT id, email, source, is_active, unsubscribed_at, created_at, updated_at
	      FROM newsletter_subscribers
	      ORDER BY created_at DESC`
	args := []interface{}{}
	if limit > 0 {
		q += ` LIMIT $1 OFFSET $2`
		args = append(args, limit, offset)
	}
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []*models.NewsletterSubscriber{}
	for rows.Next() {
		s := &models.NewsletterSubscriber{}
		if err := rows.Scan(&s.ID, &s.Email, &s.Source, &s.IsActive, &s.UnsubscribedAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, s)
	}
	return items, total, rows.Err()
}

// ── Resources ─────────────────────────────────────────────────────

func (r *CMSRepo) ListResources(ctx context.Context, category string) ([]*models.CMSResource, error) {
	const q = `SELECT id, title, category, COALESCE(file_url,''), COALESCE(description,''),
	                  sort_order, is_active, created_at, updated_at
	           FROM cms_resources WHERE ($1='' OR category=$1)
	           ORDER BY sort_order ASC, created_at DESC`
	rows, err := r.db.Query(ctx, q, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.CMSResource{}
	for rows.Next() {
		res := &models.CMSResource{}
		if err := rows.Scan(&res.ID, &res.Title, &res.Category, &res.FileURL, &res.Description,
			&res.SortOrder, &res.IsActive, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, res)
	}
	return items, rows.Err()
}

func (r *CMSRepo) CreateResource(ctx context.Context, res *models.CMSResource) error {
	const q = `INSERT INTO cms_resources (id, title, category, file_url, description, sort_order, is_active)
	           VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q, res.ID, res.Title, res.Category, res.FileURL, res.Description, res.SortOrder, res.IsActive).Scan(&res.CreatedAt, &res.UpdatedAt)
}

func (r *CMSRepo) GetResourceByID(ctx context.Context, id string) (*models.CMSResource, error) {
	const q = `SELECT id, title, category, COALESCE(file_url,''), COALESCE(description,''),
	                  sort_order, is_active, created_at, updated_at
	           FROM cms_resources WHERE id=$1`
	res := &models.CMSResource{}
	err := r.db.QueryRow(ctx, q, id).Scan(&res.ID, &res.Title, &res.Category, &res.FileURL, &res.Description,
		&res.SortOrder, &res.IsActive, &res.CreatedAt, &res.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return res, err
}

func (r *CMSRepo) UpdateResource(ctx context.Context, res *models.CMSResource) error {
	const q = `UPDATE cms_resources SET title=$2, category=$3, file_url=$4, description=$5,
	           sort_order=$6, is_active=$7, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, res.ID, res.Title, res.Category, res.FileURL, res.Description, res.SortOrder, res.IsActive)
	return err
}

func (r *CMSRepo) DeleteResource(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_resources WHERE id=$1`, id)
	return err
}

// ── Documents (Sports Rules, Press Releases, Reports, Speeches) ────

func (r *CMSRepo) ListDocuments(ctx context.Context, docType, category string, activeOnly bool) ([]*models.CMSDocument, error) {
	q := `SELECT id, doc_type, title, category, COALESCE(file_url,''), COALESCE(video_url,''), COALESCE(description,''),
	             sort_order, is_active, created_at, updated_at
	      FROM cms_documents WHERE doc_type=$1 AND ($2='' OR category=$2)`
	if activeOnly {
		q += ` AND is_active=true`
	}
	q += ` ORDER BY sort_order ASC, created_at DESC`
	rows, err := r.db.Query(ctx, q, docType, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.CMSDocument{}
	for rows.Next() {
		d := &models.CMSDocument{}
		if err := rows.Scan(&d.ID, &d.DocType, &d.Title, &d.Category, &d.FileURL, &d.VideoURL, &d.Description,
			&d.SortOrder, &d.IsActive, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, d)
	}
	return items, rows.Err()
}

func (r *CMSRepo) CreateDocument(ctx context.Context, d *models.CMSDocument) error {
	const q = `INSERT INTO cms_documents (id, doc_type, title, category, file_url, video_url, description, sort_order, is_active)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q, d.ID, d.DocType, d.Title, d.Category, d.FileURL, d.VideoURL, d.Description, d.SortOrder, d.IsActive).Scan(&d.CreatedAt, &d.UpdatedAt)
}

func (r *CMSRepo) GetDocumentByID(ctx context.Context, id string) (*models.CMSDocument, error) {
	const q = `SELECT id, doc_type, title, category, COALESCE(file_url,''), COALESCE(video_url,''), COALESCE(description,''),
	                  sort_order, is_active, created_at, updated_at
	           FROM cms_documents WHERE id=$1`
	d := &models.CMSDocument{}
	err := r.db.QueryRow(ctx, q, id).Scan(&d.ID, &d.DocType, &d.Title, &d.Category, &d.FileURL, &d.VideoURL, &d.Description,
		&d.SortOrder, &d.IsActive, &d.CreatedAt, &d.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return d, err
}

func (r *CMSRepo) UpdateDocument(ctx context.Context, d *models.CMSDocument) error {
	const q = `UPDATE cms_documents SET title=$2, category=$3, file_url=$4, video_url=$5, description=$6,
	           sort_order=$7, is_active=$8, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, d.ID, d.Title, d.Category, d.FileURL, d.VideoURL, d.Description, d.SortOrder, d.IsActive)
	return err
}

func (r *CMSRepo) DeleteDocument(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_documents WHERE id=$1`, id)
	return err
}

// ── Facilities ────────────────────────────────────────────────────

func (r *CMSRepo) ListFacilities(ctx context.Context, activeOnly bool) ([]*models.CMSFacility, error) {
	q := `SELECT id, name, slug, COALESCE(description,''), COALESCE(category,''), COALESCE(image_url,''),
	             sort_order, is_active, created_at, updated_at
	      FROM cms_facilities`
	if activeOnly {
		q += ` WHERE is_active=true`
	}
	q += ` ORDER BY sort_order ASC, name ASC`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.CMSFacility{}
	for rows.Next() {
		f := &models.CMSFacility{}
		if err := rows.Scan(&f.ID, &f.Name, &f.Slug, &f.Description, &f.Category, &f.ImageURL,
			&f.SortOrder, &f.IsActive, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	return items, rows.Err()
}

func (r *CMSRepo) GetFacilityByID(ctx context.Context, id string) (*models.CMSFacility, error) {
	const q = `SELECT id, name, slug, COALESCE(description,''), COALESCE(category,''), COALESCE(image_url,''),
	                  sort_order, is_active, created_at, updated_at
	           FROM cms_facilities WHERE id=$1`
	f := &models.CMSFacility{}
	err := r.db.QueryRow(ctx, q, id).Scan(&f.ID, &f.Name, &f.Slug, &f.Description, &f.Category, &f.ImageURL,
		&f.SortOrder, &f.IsActive, &f.CreatedAt, &f.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return f, err
}

func (r *CMSRepo) CreateFacility(ctx context.Context, f *models.CMSFacility) error {
	const q = `INSERT INTO cms_facilities (id, name, slug, description, category, image_url, sort_order, is_active)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING created_at, updated_at`
	err := r.db.QueryRow(ctx, q, f.ID, f.Name, f.Slug, f.Description, f.Category, f.ImageURL, f.SortOrder, f.IsActive).Scan(&f.CreatedAt, &f.UpdatedAt)
	if err != nil && isDuplicate(err) {
		return ErrDuplicate
	}
	return err
}

func (r *CMSRepo) UpdateFacility(ctx context.Context, f *models.CMSFacility) error {
	const q = `UPDATE cms_facilities SET name=$2, slug=$3, description=$4, category=$5, image_url=$6,
	           sort_order=$7, is_active=$8, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, f.ID, f.Name, f.Slug, f.Description, f.Category, f.ImageURL, f.SortOrder, f.IsActive)
	return err
}

func (r *CMSRepo) DeleteFacility(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_facilities WHERE id=$1`, id)
	return err
}

// ── Associations ──────────────────────────────────────────────────

func (r *CMSRepo) ListAssociations(ctx context.Context, activeOnly bool) ([]*models.CMSAssociation, error) {
	q := `SELECT id, name, slug, COALESCE(abbreviation,''), COALESCE(description,''), COALESCE(logo_url,''), COALESCE(website_url,''),
	             COALESCE(category,'Other'), COALESCE(president,''), COALESCE(secretary,''), COALESCE(address,''), COALESCE(phone,''),
	             sort_order, is_active, created_at, updated_at
	      FROM cms_associations`
	if activeOnly {
		q += ` WHERE is_active=true`
	}
	q += ` ORDER BY sort_order ASC, name ASC`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.CMSAssociation{}
	for rows.Next() {
		a := &models.CMSAssociation{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Slug, &a.Abbreviation, &a.Description, &a.LogoURL, &a.WebsiteURL,
			&a.Category, &a.President, &a.Secretary, &a.Address, &a.Phone,
			&a.SortOrder, &a.IsActive, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, rows.Err()
}

func (r *CMSRepo) GetAssociationByID(ctx context.Context, id string) (*models.CMSAssociation, error) {
	const q = `SELECT id, name, slug, COALESCE(abbreviation,''), COALESCE(description,''), COALESCE(logo_url,''), COALESCE(website_url,''),
	                  COALESCE(category,'Other'), COALESCE(president,''), COALESCE(secretary,''), COALESCE(address,''), COALESCE(phone,''),
	                  sort_order, is_active, created_at, updated_at
	           FROM cms_associations WHERE id=$1`
	a := &models.CMSAssociation{}
	err := r.db.QueryRow(ctx, q, id).Scan(&a.ID, &a.Name, &a.Slug, &a.Abbreviation, &a.Description, &a.LogoURL, &a.WebsiteURL,
		&a.Category, &a.President, &a.Secretary, &a.Address, &a.Phone,
		&a.SortOrder, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return a, err
}

func (r *CMSRepo) CreateAssociation(ctx context.Context, a *models.CMSAssociation) error {
	const q = `INSERT INTO cms_associations (id, name, slug, abbreviation, description, logo_url, website_url, category, president, secretary, address, phone, sort_order, is_active)
	           VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING created_at, updated_at`
	err := r.db.QueryRow(ctx, q, a.ID, a.Name, a.Slug, a.Abbreviation, a.Description, a.LogoURL, a.WebsiteURL, a.Category, a.President, a.Secretary, a.Address, a.Phone, a.SortOrder, a.IsActive).Scan(&a.CreatedAt, &a.UpdatedAt)
	if err != nil && isDuplicate(err) {
		return ErrDuplicate
	}
	return err
}

func (r *CMSRepo) UpdateAssociation(ctx context.Context, a *models.CMSAssociation) error {
	const q = `UPDATE cms_associations SET name=$2, slug=$3, abbreviation=$4, description=$5, logo_url=$6,
	           website_url=$7, category=$8, president=$9, secretary=$10, address=$11, phone=$12,
	           sort_order=$13, is_active=$14, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, a.ID, a.Name, a.Slug, a.Abbreviation, a.Description, a.LogoURL, a.WebsiteURL, a.Category, a.President, a.Secretary, a.Address, a.Phone, a.SortOrder, a.IsActive)
	return err
}

func (r *CMSRepo) DeleteAssociation(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_associations WHERE id=$1`, id)
	return err
}

// ── Invest with Us ────────────────────────────────────────────────

func (r *CMSRepo) ListInvest(ctx context.Context) ([]*models.CMSInvest, error) {
	const q = `SELECT id, title, COALESCE(subtitle,''), COALESCE(content,''), COALESCE(image_url,''),
	                  sort_order, is_active, created_at, updated_at
	           FROM cms_invest ORDER BY sort_order ASC, created_at ASC`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*models.CMSInvest{}
	for rows.Next() {
		inv := &models.CMSInvest{}
		if err := rows.Scan(&inv.ID, &inv.Title, &inv.Subtitle, &inv.Content, &inv.ImageURL,
			&inv.SortOrder, &inv.IsActive, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, inv)
	}
	return items, rows.Err()
}

func (r *CMSRepo) GetInvestByID(ctx context.Context, id string) (*models.CMSInvest, error) {
	const q = `SELECT id, title, COALESCE(subtitle,''), COALESCE(content,''), COALESCE(image_url,''),
	                  sort_order, is_active, created_at, updated_at
	           FROM cms_invest WHERE id=$1`
	inv := &models.CMSInvest{}
	err := r.db.QueryRow(ctx, q, id).Scan(&inv.ID, &inv.Title, &inv.Subtitle, &inv.Content, &inv.ImageURL,
		&inv.SortOrder, &inv.IsActive, &inv.CreatedAt, &inv.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return inv, err
}

func (r *CMSRepo) CreateInvest(ctx context.Context, inv *models.CMSInvest) error {
	const q = `INSERT INTO cms_invest (id, title, subtitle, content, image_url, sort_order, is_active)
	           VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q, inv.ID, inv.Title, inv.Subtitle, inv.Content, inv.ImageURL, inv.SortOrder, inv.IsActive).Scan(&inv.CreatedAt, &inv.UpdatedAt)
}

func (r *CMSRepo) UpdateInvest(ctx context.Context, inv *models.CMSInvest) error {
	const q = `UPDATE cms_invest SET title=$2, subtitle=$3, content=$4, image_url=$5,
	           sort_order=$6, is_active=$7, updated_at=NOW() WHERE id=$1`
	_, err := r.db.Exec(ctx, q, inv.ID, inv.Title, inv.Subtitle, inv.Content, inv.ImageURL, inv.SortOrder, inv.IsActive)
	return err
}

func (r *CMSRepo) DeleteInvest(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_invest WHERE id=$1`, id)
	return err
}

// ── Team Members ─────────────────────────────────────────────────

func (r *CMSRepo) ListTeamMembers(ctx context.Context, group string, activeOnly bool) ([]*models.CMSTeamMember, error) {
	q := `SELECT m.id, m.full_name, COALESCE(m.designation,''), COALESCE(m.image_url,''), COALESCE(m.bio,''),
	             m.sort_order, m.is_active, m.member_group, m.department_id, COALESCE(d.name,''),
	             m.created_at, m.updated_at
	      FROM cms_team_members m
	      LEFT JOIN departments d ON d.id = m.department_id
	      WHERE m.member_group=$1`
	args := []any{group}
	if activeOnly {
		q += ` AND m.is_active=TRUE`
	}
	q += ` ORDER BY m.sort_order ASC, m.created_at ASC`
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*models.CMSTeamMember{}
	for rows.Next() {
		m := &models.CMSTeamMember{}
		if err := rows.Scan(&m.ID, &m.FullName, &m.Designation, &m.ImageURL, &m.Bio,
			&m.SortOrder, &m.IsActive, &m.MemberGroup, &m.DepartmentID, &m.DepartmentName,
			&m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func (r *CMSRepo) GetTeamMemberByID(ctx context.Context, id string) (*models.CMSTeamMember, error) {
	m := &models.CMSTeamMember{}
	err := r.db.QueryRow(ctx, `SELECT m.id, m.full_name, COALESCE(m.designation,''), COALESCE(m.image_url,''), COALESCE(m.bio,''),
	                                  m.sort_order, m.is_active, m.member_group, m.department_id, COALESCE(d.name,''),
	                                  m.created_at, m.updated_at
	                           FROM cms_team_members m
	                           LEFT JOIN departments d ON d.id = m.department_id
	                           WHERE m.id=$1`, id).
		Scan(&m.ID, &m.FullName, &m.Designation, &m.ImageURL, &m.Bio,
			&m.SortOrder, &m.IsActive, &m.MemberGroup, &m.DepartmentID, &m.DepartmentName,
			&m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return m, err
}

func (r *CMSRepo) CreateTeamMember(ctx context.Context, m *models.CMSTeamMember) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO cms_team_members (id, full_name, designation, image_url, bio, sort_order, is_active, member_group, department_id)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING created_at, updated_at`,
		m.ID, m.FullName, m.Designation, m.ImageURL, m.Bio, m.SortOrder, m.IsActive, m.MemberGroup, m.DepartmentID,
	).Scan(&m.CreatedAt, &m.UpdatedAt)
}

func (r *CMSRepo) UpdateTeamMember(ctx context.Context, m *models.CMSTeamMember) error {
	_, err := r.db.Exec(ctx,
		`UPDATE cms_team_members SET full_name=$2, designation=$3, image_url=$4, bio=$5,
		 sort_order=$6, is_active=$7, department_id=$8, updated_at=NOW() WHERE id=$1`,
		m.ID, m.FullName, m.Designation, m.ImageURL, m.Bio, m.SortOrder, m.IsActive, m.DepartmentID)
	return err
}

// ListDepartmentsWithStaffCount returns the 10-tier institutional units
// along with how many active team members are assigned to each.
type DepartmentWithCount struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	StaffCount  int    `json:"staff_count"`
}

func (r *CMSRepo) ListDepartmentsWithStaffCount(ctx context.Context) ([]*DepartmentWithCount, error) {
	const q = `SELECT d.id, d.name, d.code, COALESCE(d.description,''),
	                  COALESCE(c.cnt, 0)
	           FROM departments d
	           LEFT JOIN (
	             SELECT department_id, COUNT(*) AS cnt
	             FROM cms_team_members WHERE is_active=TRUE AND department_id IS NOT NULL
	             GROUP BY department_id
	           ) c ON c.department_id = d.id
	           ORDER BY d.name`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*DepartmentWithCount{}
	for rows.Next() {
		d := &DepartmentWithCount{}
		if err := rows.Scan(&d.ID, &d.Name, &d.Code, &d.Description, &d.StaffCount); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *CMSRepo) CreateDepartment(ctx context.Context, d *DepartmentWithCount) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO departments (id, name, code, description) VALUES ($1,$2,$3,$4)`,
		d.ID, d.Name, d.Code, d.Description)
	return err
}

func (r *CMSRepo) UpdateDepartment(ctx context.Context, d *DepartmentWithCount) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE departments SET name=$2, code=$3, description=$4 WHERE id=$1`,
		d.ID, d.Name, d.Code, d.Description)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CMSRepo) DeleteDepartment(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM departments WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CMSRepo) ListNotifications(ctx context.Context, userID, status string, limit, offset int) ([]*models.Notification, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	const countQ = `SELECT COUNT(*) FROM notifications WHERE ($1='' OR user_id=$1 OR user_id IS NULL) AND ($2='' OR status=$2) AND status <> 'dismissed'`
	const q = `SELECT id, user_id, type, title, message, status, icon_key, created_at FROM notifications WHERE ($1='' OR user_id=$1 OR user_id IS NULL) AND ($2='' OR status=$2) AND status <> 'dismissed' ORDER BY created_at DESC LIMIT $3 OFFSET $4`
	var total int64
	if err := r.db.QueryRow(ctx, countQ, userID, status).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, q, userID, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []*models.Notification{}
	for rows.Next() {
		n := &models.Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Status, &n.IconKey, &n.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, n)
	}
	return out, total, rows.Err()
}

func (r *CMSRepo) CreateNotification(ctx context.Context, n *models.Notification) error {
	const q = `INSERT INTO notifications (id, user_id, type, title, message, status, icon_key) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING created_at`
	if err := r.db.QueryRow(ctx, q, n.ID, n.UserID, n.Type, n.Title, n.Message, n.Status, n.IconKey).Scan(&n.CreatedAt); err != nil {
		return err
	}
	_, _ = r.db.Exec(ctx, `INSERT INTO system_notifications (id,user_id,event_type,title,message,severity,payload,status,created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb,$8,$9)
		ON CONFLICT (id) DO NOTHING`,
		n.ID, n.UserID, n.Type, n.Title, n.Message, notificationSeverity(n.Type), `{}`, n.Status, n.CreatedAt)
	return nil
}

func notificationSeverity(kind string) string {
	switch kind {
	case "new_sign_in", "unusual_activity", "deploy_failed", "service_offline":
		return "warning"
	case "system_failure", "security_alert":
		return "critical"
	case "system_success", "backup_success":
		return "success"
	default:
		return "info"
	}
}

func (r *CMSRepo) UpdateNotificationStatus(ctx context.Context, id, userID, status string) error {
	ct, err := r.db.Exec(ctx, `UPDATE notifications SET status=$3 WHERE id=$1 AND ($2='' OR user_id=$2 OR user_id IS NULL)`, id, userID, status)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CMSRepo) ClearNotifications(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `UPDATE notifications SET status='dismissed' WHERE ($1='' OR user_id=$1 OR user_id IS NULL)`, userID)
	return err
}

func (r *CMSRepo) MarkAllNotificationsRead(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `UPDATE notifications SET status='read' WHERE status='unread' AND ($1='' OR user_id=$1 OR user_id IS NULL)`, userID)
	return err
}

func (r *CMSRepo) UpsertInboundSubmission(ctx context.Context, s *models.InboundSubmission) error {
	if s.ID == "" {
		return errors.New("inbound submission id required")
	}
	if s.Payload == nil {
		s.Payload = json.RawMessage(`{}`)
	}
	if s.InternalNotes == nil {
		s.InternalNotes = json.RawMessage(`[]`)
	}
	if s.StatusState == "" {
		s.StatusState = "unread"
	}
	if s.WorkflowStatus == "" {
		s.WorkflowStatus = "pending_review"
	}
	const q = `INSERT INTO inbound_submissions (id,source_type,source_id,payload,status_state,workflow_status,assigned_admin_id,internal_notes)
	           VALUES ($1,$2,$3,$4::jsonb,$5,$6,$7,$8::jsonb)
	           ON CONFLICT (source_type, source_id) DO UPDATE SET
	             payload=EXCLUDED.payload,
	             status_state=EXCLUDED.status_state,
	             workflow_status=EXCLUDED.workflow_status,
	             assigned_admin_id=COALESCE(EXCLUDED.assigned_admin_id, inbound_submissions.assigned_admin_id),
	             updated_at=NOW()
	           RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, q, s.ID, s.SourceType, nullableText(s.SourceID), s.Payload, s.StatusState, s.WorkflowStatus, s.AssignedAdminID, s.InternalNotes).Scan(&s.CreatedAt, &s.UpdatedAt)
}

func (r *CMSRepo) ListInboundSubmissions(ctx context.Context, sourceType, status string, limit, offset int) ([]*models.InboundSubmission, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	const countQ = `SELECT COUNT(*) FROM inbound_submissions WHERE ($1='' OR source_type=$1) AND ($2='' OR status_state=$2)`
	const q = `SELECT id,source_type,COALESCE(source_id,''),payload,status_state,workflow_status,assigned_admin_id,internal_notes,created_at,updated_at
	           FROM inbound_submissions
	           WHERE ($1='' OR source_type=$1) AND ($2='' OR status_state=$2)
	           ORDER BY created_at DESC LIMIT $3 OFFSET $4`
	var total int64
	if err := r.db.QueryRow(ctx, countQ, sourceType, status).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, q, sourceType, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []*models.InboundSubmission{}
	for rows.Next() {
		s := &models.InboundSubmission{}
		if err := rows.Scan(&s.ID, &s.SourceType, &s.SourceID, &s.Payload, &s.StatusState, &s.WorkflowStatus, &s.AssignedAdminID, &s.InternalNotes, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, s)
	}
	return out, total, rows.Err()
}

func (r *CMSRepo) InboundSubmissionCounts(ctx context.Context) (*models.InboundSubmissionCounts, error) {
	rows, err := r.db.Query(ctx, `SELECT source_type, COUNT(*) FROM inbound_submissions WHERE status_state='unread' GROUP BY source_type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := &models.InboundSubmissionCounts{}
	for rows.Next() {
		var source string
		var count int64
		if err := rows.Scan(&source, &count); err != nil {
			return nil, err
		}
		switch source {
		case "blog_comment":
			counts.BlogComments = count
		case "contact_form":
			counts.ContactForms = count
		case "investment_request":
			counts.InvestmentRequests = count
		}
		counts.Total += count
	}
	return counts, rows.Err()
}

func (r *CMSRepo) UpdateInboundSubmission(ctx context.Context, id, statusState, workflowStatus, adminID string) error {
	if statusState == "" && workflowStatus == "" && adminID == "" {
		return nil
	}
	ct, err := r.db.Exec(ctx, `UPDATE inbound_submissions
		SET status_state=CASE WHEN $2='' THEN status_state ELSE $2 END,
		    workflow_status=CASE WHEN $3='' THEN workflow_status ELSE $3 END,
		    assigned_admin_id=CASE WHEN $4='' THEN assigned_admin_id ELSE $4 END,
		    updated_at=NOW()
		WHERE id=$1`, id, statusState, workflowStatus, adminID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CMSRepo) ArchiveInboundBySource(ctx context.Context, sourceType, workflowStatus string) error {
	if workflowStatus == "" {
		workflowStatus = "archived"
	}
	_, err := r.db.Exec(ctx, `UPDATE inbound_submissions
		SET status_state='archived', workflow_status=$2, updated_at=NOW()
		WHERE source_type=$1`, sourceType, workflowStatus)
	return err
}

func (r *CMSRepo) AddInboundSubmissionNote(ctx context.Context, id, adminID, note string) error {
	payload, _ := json.Marshal(map[string]string{
		"id":         "note_" + time.Now().UTC().Format("20060102150405.000000000"),
		"admin_id":   adminID,
		"note":       note,
		"created_at": time.Now().UTC().Format(time.RFC3339),
	})
	ct, err := r.db.Exec(ctx, `UPDATE inbound_submissions
		SET internal_notes = COALESCE(internal_notes, '[]'::jsonb) || jsonb_build_array($2::jsonb),
		    updated_at=NOW()
		WHERE id=$1`, id, string(payload))
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CMSRepo) ListContactMessages(ctx context.Context, status string, limit, offset int) ([]*models.ContactMessage, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	const countQ = `SELECT COUNT(*) FROM contact_messages WHERE ($1='' OR status=$1)`
	const q = `SELECT id, name, email, subject, message, status, created_at FROM contact_messages WHERE ($1='' OR status=$1) ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	var total int64
	if err := r.db.QueryRow(ctx, countQ, status).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, q, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []*models.ContactMessage{}
	for rows.Next() {
		m := &models.ContactMessage{}
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Subject, &m.Message, &m.Status, &m.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, m)
	}
	return out, total, rows.Err()
}

func (r *CMSRepo) CreateContactMessage(ctx context.Context, m *models.ContactMessage) error {
	const q = `INSERT INTO contact_messages (id, name, email, subject, message, status) VALUES ($1,$2,$3,$4,$5,$6) RETURNING created_at`
	return r.db.QueryRow(ctx, q, m.ID, m.Name, m.Email, m.Subject, m.Message, m.Status).Scan(&m.CreatedAt)
}

func (r *CMSRepo) UpdateContactMessageStatus(ctx context.Context, id, status string) error {
	ct, err := r.db.Exec(ctx, `UPDATE contact_messages SET status=$2 WHERE id=$1`, id, status)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func nullableText(value string) interface{} {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return strings.TrimSpace(value)
}

func (r *CMSRepo) DeleteContactMessage(ctx context.Context, id string) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM contact_messages WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CMSRepo) ClearContactMessages(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `DELETE FROM contact_messages`)
	return err
}

func (r *CMSRepo) DeleteTeamMember(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cms_team_members WHERE id=$1`, id)
	return err
}

// ── Helpers ───────────────────────────────────────────────────────

func isDuplicate(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "unique constraint"))
}

func nullableStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
