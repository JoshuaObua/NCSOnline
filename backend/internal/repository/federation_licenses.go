package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FederationLicenseRepo struct {
	db *pgxpool.Pool
}

func NewFederationLicenseRepo(db *pgxpool.Pool) *FederationLicenseRepo {
	return &FederationLicenseRepo{db: db}
}

type ListLicensesFilter struct {
	FederationID string
	Status       string
	LicenseType  string
	Search       string
	StartDate    string
	EndDate      string
	SortBy       string
	SortOrder    string
}

func (r *FederationLicenseRepo) Create(ctx context.Context, lic *models.FederationLicense, actorID, actorName string) error {
	if lic.ID == "" {
		lic.ID = uuid.NewString()
	}
	if lic.LicenseType == "" {
		lic.LicenseType = "FULL_RECOGNITION"
	}
	if lic.Category == "" {
		lic.Category = "Tier 1 National Sports Federation"
	}
	if lic.Status == "" {
		lic.Status = "ACTIVE"
	}
	if lic.IssueDate.IsZero() {
		lic.IssueDate = time.Now()
	}
	if lic.ExpiryDate.IsZero() {
		lic.ExpiryDate = lic.IssueDate.AddDate(1, 0, 0)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const q = `
		INSERT INTO federation_licenses (
			id, federation_id, license_number, license_type, category,
			issue_date, expiry_date, status, conditions, document_url,
			issued_by, issued_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, NOW(), NOW(), NOW()
		)
		RETURNING created_at, updated_at
	`
	err = tx.QueryRow(ctx, q,
		lic.ID, lic.FederationID, lic.LicenseNumber, lic.LicenseType, lic.Category,
		lic.IssueDate, lic.ExpiryDate, lic.Status, lic.Conditions, lic.DocumentURL,
		actorID,
	).Scan(&lic.CreatedAt, &lic.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert audit log
	const logQ = `
		INSERT INTO federation_license_logs (
			id, license_id, action, performed_by, performed_by_name, reason,
			old_status, new_status, new_expiry_date, notes, created_at
		) VALUES (
			$1, $2, 'CREATED', $3, $4, 'Initial License Issuance',
			NULL, $5, $6, $7, NOW()
		)
	`
	_, err = tx.Exec(ctx, logQ,
		uuid.NewString(), lic.ID, actorID, actorName,
		lic.Status, lic.ExpiryDate, lic.Conditions,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *FederationLicenseRepo) GetByID(ctx context.Context, id string) (*models.FederationLicense, error) {
	const q = `
		SELECT fl.id, fl.federation_id,
		       COALESCE(f.name, ''), COALESCE(f.acronym, ''), COALESCE(f.ncs_registration_number, ''),
		       COALESCE(f.logo_url, ''), COALESCE(f.president, ''), COALESCE(f.secretary, ''),
		       fl.license_number, fl.license_type, fl.category,
		       fl.issue_date, fl.expiry_date, fl.status,
		       COALESCE(fl.conditions, ''), COALESCE(fl.document_url, ''),
		       fl.issued_by, COALESCE(u_iss.first_name || ' ' || u_iss.last_name, ''), fl.issued_at,
		       fl.extended_by, COALESCE(u_ext.first_name || ' ' || u_ext.last_name, ''), fl.extended_at,
		       fl.previous_expiry_date, COALESCE(fl.extension_reason, ''),
		       fl.revoked_by, COALESCE(u_rev.first_name || ' ' || u_rev.last_name, ''), fl.revoked_at,
		       COALESCE(fl.revocation_reason, ''),
		       fl.created_at, fl.updated_at
		FROM federation_licenses fl
		LEFT JOIN federations f ON f.id = fl.federation_id
		LEFT JOIN users u_iss ON u_iss.id = fl.issued_by
		LEFT JOIN users u_ext ON u_ext.id = fl.extended_by
		LEFT JOIN users u_rev ON u_rev.id = fl.revoked_by
		WHERE fl.id = $1
	`
	row := r.db.QueryRow(ctx, q, id)
	lic, err := r.scanLicense(row)
	if err != nil {
		return nil, err
	}

	// Fetch logs
	logs, err := r.GetLicenseLogs(ctx, lic.ID)
	if err == nil {
		lic.Logs = logs
	}
	return lic, nil
}

func (r *FederationLicenseRepo) GetActiveLicenseByFederationID(ctx context.Context, fedID string) (*models.FederationLicense, error) {
	const q = `
		SELECT fl.id, fl.federation_id,
		       COALESCE(f.name, ''), COALESCE(f.acronym, ''), COALESCE(f.ncs_registration_number, ''),
		       COALESCE(f.logo_url, ''), COALESCE(f.president, ''), COALESCE(f.secretary, ''),
		       fl.license_number, fl.license_type, fl.category,
		       fl.issue_date, fl.expiry_date, fl.status,
		       COALESCE(fl.conditions, ''), COALESCE(fl.document_url, ''),
		       fl.issued_by, COALESCE(u_iss.first_name || ' ' || u_iss.last_name, ''), fl.issued_at,
		       fl.extended_by, COALESCE(u_ext.first_name || ' ' || u_ext.last_name, ''), fl.extended_at,
		       fl.previous_expiry_date, COALESCE(fl.extension_reason, ''),
		       fl.revoked_by, COALESCE(u_rev.first_name || ' ' || u_rev.last_name, ''), fl.revoked_at,
		       COALESCE(fl.revocation_reason, ''),
		       fl.created_at, fl.updated_at
		FROM federation_licenses fl
		LEFT JOIN federations f ON f.id = fl.federation_id
		LEFT JOIN users u_iss ON u_iss.id = fl.issued_by
		LEFT JOIN users u_ext ON u_ext.id = fl.extended_by
		LEFT JOIN users u_rev ON u_rev.id = fl.revoked_by
		WHERE fl.federation_id = $1
		ORDER BY fl.created_at DESC
		LIMIT 1
	`
	row := r.db.QueryRow(ctx, q, fedID)
	lic, err := r.scanLicense(row)
	if err != nil {
		return nil, err
	}
	logs, err := r.GetLicenseLogs(ctx, lic.ID)
	if err == nil {
		lic.Logs = logs
	}
	return lic, nil
}

func (r *FederationLicenseRepo) List(ctx context.Context, f ListLicensesFilter, p *models.PaginationParams) ([]*models.FederationLicense, int, error) {
	where := []string{"1=1"}
	args := []any{}
	idx := 1

	if f.FederationID != "" {
		where = append(where, fmt.Sprintf("fl.federation_id = $%d", idx))
		args = append(args, f.FederationID)
		idx++
	}
	if f.Status != "" && f.Status != "ALL" {
		where = append(where, fmt.Sprintf("UPPER(fl.status) = UPPER($%d)", idx))
		args = append(args, f.Status)
		idx++
	}
	if f.LicenseType != "" && f.LicenseType != "ALL" {
		where = append(where, fmt.Sprintf("UPPER(fl.license_type) = UPPER($%d)", idx))
		args = append(args, f.LicenseType)
		idx++
	}
	if f.Search != "" {
		s := "%" + strings.ToLower(f.Search) + "%"
		where = append(where, fmt.Sprintf("(LOWER(fl.license_number) LIKE $%d OR LOWER(f.name) LIKE $%d OR LOWER(COALESCE(f.acronym, '')) LIKE $%d OR LOWER(COALESCE(f.ncs_registration_number, '')) LIKE $%d OR LOWER(COALESCE(f.president, '')) LIKE $%d OR LOWER(COALESCE(f.secretary, '')) LIKE $%d)", idx, idx, idx, idx, idx, idx))
		args = append(args, s)
		idx++
	}
	if f.StartDate != "" {
		where = append(where, fmt.Sprintf("fl.issue_date >= $%d::date", idx))
		args = append(args, f.StartDate)
		idx++
	}
	if f.EndDate != "" {
		where = append(where, fmt.Sprintf("fl.expiry_date <= $%d::date", idx))
		args = append(args, f.EndDate)
		idx++
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQ := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM federation_licenses fl
		LEFT JOIN federations f ON f.id = fl.federation_id
		WHERE %s
	`, whereClause)
	if err := r.db.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
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

	orderCol := "fl.created_at"
	switch strings.ToLower(f.SortBy) {
	case "license_number", "number":
		orderCol = "fl.license_number"
	case "federation", "federation_name":
		orderCol = "f.name"
	case "issue_date":
		orderCol = "fl.issue_date"
	case "expiry_date":
		orderCol = "fl.expiry_date"
	case "status":
		orderCol = "fl.status"
	case "type", "license_type":
		orderCol = "fl.license_type"
	}

	orderDir := "DESC"
	if strings.EqualFold(f.SortOrder, "asc") {
		orderDir = "ASC"
	}

	args = append(args, limit, offset)
	q := fmt.Sprintf(`
		SELECT fl.id, fl.federation_id,
		       COALESCE(f.name, ''), COALESCE(f.acronym, ''), COALESCE(f.ncs_registration_number, ''),
		       COALESCE(f.logo_url, ''), COALESCE(f.president, ''), COALESCE(f.secretary, ''),
		       fl.license_number, fl.license_type, fl.category,
		       fl.issue_date, fl.expiry_date, fl.status,
		       COALESCE(fl.conditions, ''), COALESCE(fl.document_url, ''),
		       fl.issued_by, COALESCE(u_iss.first_name || ' ' || u_iss.last_name, ''), fl.issued_at,
		       fl.extended_by, COALESCE(u_ext.first_name || ' ' || u_ext.last_name, ''), fl.extended_at,
		       fl.previous_expiry_date, COALESCE(fl.extension_reason, ''),
		       fl.revoked_by, COALESCE(u_rev.first_name || ' ' || u_rev.last_name, ''), fl.revoked_at,
		       COALESCE(fl.revocation_reason, ''),
		       fl.created_at, fl.updated_at
		FROM federation_licenses fl
		LEFT JOIN federations f ON f.id = fl.federation_id
		LEFT JOIN users u_iss ON u_iss.id = fl.issued_by
		LEFT JOIN users u_ext ON u_ext.id = fl.extended_by
		LEFT JOIN users u_rev ON u_rev.id = fl.revoked_by
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderCol, orderDir, idx, idx+1)

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*models.FederationLicense
	for rows.Next() {
		lic, err := r.scanLicense(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, lic)
	}
	return list, total, nil
}

func (r *FederationLicenseRepo) Extend(ctx context.Context, id string, newExpiryDate time.Time, reason, actorID, actorName string) (*models.FederationLicense, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var oldExpiry time.Time
	var oldStatus string
	err = tx.QueryRow(ctx, `SELECT expiry_date, status FROM federation_licenses WHERE id = $1 FOR UPDATE`, id).Scan(&oldExpiry, &oldStatus)
	if err != nil {
		return nil, err
	}

	const q = `
		UPDATE federation_licenses
		SET previous_expiry_date = $2,
		    expiry_date = $3,
		    status = 'EXTENDED',
		    extended_by = $4,
		    extended_at = NOW(),
		    extension_reason = $5,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, q, id, oldExpiry, newExpiryDate, actorID, reason)
	if err != nil {
		return nil, err
	}

	const logQ = `
		INSERT INTO federation_license_logs (
			id, license_id, action, performed_by, performed_by_name, reason,
			old_status, new_status, old_expiry_date, new_expiry_date, notes, created_at
		) VALUES (
			$1, $2, 'EXTENDED', $3, $4, $5,
			$6, 'EXTENDED', $7, $8, 'Validity extended by authorized administrator.', NOW()
		)
	`
	_, err = tx.Exec(ctx, logQ,
		uuid.NewString(), id, actorID, actorName, reason,
		oldStatus, oldExpiry, newExpiryDate,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *FederationLicenseRepo) Revoke(ctx context.Context, id, reason, actorID, actorName string) (*models.FederationLicense, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var oldStatus string
	var expiryDate time.Time
	err = tx.QueryRow(ctx, `SELECT status, expiry_date FROM federation_licenses WHERE id = $1 FOR UPDATE`, id).Scan(&oldStatus, &expiryDate)
	if err != nil {
		return nil, err
	}

	const q = `
		UPDATE federation_licenses
		SET status = 'REVOKED',
		    revoked_by = $2,
		    revoked_at = NOW(),
		    revocation_reason = $3,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, q, id, actorID, reason)
	if err != nil {
		return nil, err
	}

	const logQ = `
		INSERT INTO federation_license_logs (
			id, license_id, action, performed_by, performed_by_name, reason,
			old_status, new_status, old_expiry_date, new_expiry_date, notes, created_at
		) VALUES (
			$1, $2, 'REVOKED', $3, $4, $5,
			$6, 'REVOKED', $7, $7, 'License officially revoked by NCS Authority.', NOW()
		)
	`
	_, err = tx.Exec(ctx, logQ,
		uuid.NewString(), id, actorID, actorName, reason,
		oldStatus, expiryDate,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *FederationLicenseRepo) Reinstate(ctx context.Context, id, reason, actorID, actorName string) (*models.FederationLicense, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var oldStatus string
	var expiryDate time.Time
	err = tx.QueryRow(ctx, `SELECT status, expiry_date FROM federation_licenses WHERE id = $1 FOR UPDATE`, id).Scan(&oldStatus, &expiryDate)
	if err != nil {
		return nil, err
	}

	newStatus := "ACTIVE"
	if expiryDate.Before(time.Now()) {
		newStatus = "EXPIRED"
	}

	const q = `
		UPDATE federation_licenses
		SET status = $2,
		    revoked_by = NULL,
		    revoked_at = NULL,
		    revocation_reason = NULL,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, q, id, newStatus)
	if err != nil {
		return nil, err
	}

	const logQ = `
		INSERT INTO federation_license_logs (
			id, license_id, action, performed_by, performed_by_name, reason,
			old_status, new_status, old_expiry_date, new_expiry_date, notes, created_at
		) VALUES (
			$1, $2, 'REINSTATED', $3, $4, $5,
			$6, $7, $8, $8, 'License status reinstated by NCS Authority.', NOW()
		)
	`
	_, err = tx.Exec(ctx, logQ,
		uuid.NewString(), id, actorID, actorName, reason,
		oldStatus, newStatus, expiryDate,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *FederationLicenseRepo) GetLicenseLogs(ctx context.Context, licenseID string) ([]*models.FederationLicenseLog, error) {
	const q = `
		SELECT fll.id, fll.license_id, fll.action, fll.performed_by,
		       COALESCE(fll.performed_by_name, u.first_name || ' ' || u.last_name, 'System Administrator'),
		       COALESCE(fll.reason, ''), COALESCE(fll.old_status, ''), COALESCE(fll.new_status, ''),
		       fll.old_expiry_date, fll.new_expiry_date, COALESCE(fll.notes, ''), fll.created_at
		FROM federation_license_logs fll
		LEFT JOIN users u ON u.id = fll.performed_by
		WHERE fll.license_id = $1
		ORDER BY fll.created_at DESC
	`
	rows, err := r.db.Query(ctx, q, licenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.FederationLicenseLog
	for rows.Next() {
		var l models.FederationLicenseLog
		if err := rows.Scan(
			&l.ID, &l.LicenseID, &l.Action, &l.PerformedBy, &l.PerformedByName,
			&l.Reason, &l.OldStatus, &l.NewStatus,
			&l.OldExpiryDate, &l.NewExpiryDate, &l.Notes, &l.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, nil
}

func (r *FederationLicenseRepo) GetKPIs(ctx context.Context) (*models.LicenseKPIs, error) {
	const q = `
		SELECT
			COUNT(*) AS total_licenses,
			COUNT(CASE WHEN status = 'ACTIVE' THEN 1 END) AS active_licenses,
			COUNT(CASE WHEN status = 'EXTENDED' THEN 1 END) AS extended_licenses,
			COUNT(CASE WHEN status = 'REVOKED' THEN 1 END) AS revoked_licenses,
			COUNT(CASE WHEN status IN ('ACTIVE', 'EXTENDED') AND expiry_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '60 days' THEN 1 END) AS expiring_soon
		FROM federation_licenses
	`
	var kpi models.LicenseKPIs
	err := r.db.QueryRow(ctx, q).Scan(
		&kpi.TotalLicenses, &kpi.ActiveLicenses, &kpi.ExtendedLicenses,
		&kpi.RevokedLicenses, &kpi.ExpiringSoon,
	)
	if err != nil {
		return nil, err
	}
	return &kpi, nil
}

func (r *FederationLicenseRepo) scanLicense(s rowScanner) (*models.FederationLicense, error) {
	var lic models.FederationLicense
	err := s.Scan(
		&lic.ID, &lic.FederationID,
		&lic.FederationName, &lic.FederationAcronym, &lic.FederationRegNo,
		&lic.FederationLogo, &lic.PresidentName, &lic.SecretaryName,
		&lic.LicenseNumber, &lic.LicenseType, &lic.Category,
		&lic.IssueDate, &lic.ExpiryDate, &lic.Status,
		&lic.Conditions, &lic.DocumentURL,
		&lic.IssuedBy, &lic.IssuedByName, &lic.IssuedAt,
		&lic.ExtendedBy, &lic.ExtendedByName, &lic.ExtendedAt,
		&lic.PreviousExpiryDate, &lic.ExtensionReason,
		&lic.RevokedBy, &lic.RevokedByName, &lic.RevokedAt,
		&lic.RevocationReason,
		&lic.CreatedAt, &lic.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &lic, nil
}
