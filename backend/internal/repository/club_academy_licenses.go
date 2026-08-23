package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
)

type ClubAcademyLicenseRepo struct {
	db *pgxpool.Pool
}

func NewClubAcademyLicenseRepo(db *pgxpool.Pool) *ClubAcademyLicenseRepo {
	return &ClubAcademyLicenseRepo{db: db}
}

type ListClubLicensesFilter struct {
	ClubID       string
	FederationID string
	Status       string
	LicenseType  string
	Category     string
	Search       string
	StartDate    string
	EndDate      string
	SortBy       string
	SortOrder    string
}

func (r *ClubAcademyLicenseRepo) Create(ctx context.Context, lic *models.ClubAcademyLicense, actorID, actorName string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if lic.ID == "" {
		lic.ID = uuid.NewString()
	}
	if lic.Status == "" {
		lic.Status = "ACTIVE"
	}
	if lic.IssueDate.IsZero() {
		lic.IssueDate = time.Now()
	}

	const q = `
		INSERT INTO club_academy_licenses (
			id, club_id, license_number, license_type, category,
			issue_date, expiry_date, status, conditions, document_url,
			issued_by, issued_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, NOW(), NOW(), NOW()
		)
	`
	_, err = tx.Exec(ctx, q,
		lic.ID, lic.ClubID, lic.LicenseNumber, lic.LicenseType, lic.Category,
		lic.IssueDate, lic.ExpiryDate, lic.Status, lic.Conditions, lic.DocumentURL,
		actorID,
	)
	if err != nil {
		return err
	}

	// Insert audit log
	const logQ = `
		INSERT INTO club_academy_license_logs (
			id, license_id, action, performed_by, performed_by_name, reason,
			old_status, new_status, new_expiry_date, notes, created_at
		) VALUES (
			$1, $2, 'CREATED', $3, $4, 'Initial Club / Academy License Issuance',
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

func (r *ClubAcademyLicenseRepo) GetByID(ctx context.Context, id string) (*models.ClubAcademyLicense, error) {
	const q = `
		SELECT cal.id, cal.club_id,
		       COALESCE(c.name, ''), COALESCE(c.club_number, ''), COALESCE(c.acronym, ''),
		       COALESCE(c.district, ''), COALESCE(c.region, ''), COALESCE(c.email, ''), COALESCE(c.phone, ''),
		       COALESCE(c.logo_url, ''),
		       COALESCE(c.president, c.contact_person, ''),
		       COALESCE(c.secretary, ''),
		       COALESCE(c.federation_id, ''), COALESCE(f.name, ''),
		       cal.license_number, cal.license_type, cal.category,
		       cal.issue_date, cal.expiry_date, cal.status,
		       COALESCE(cal.conditions, ''), COALESCE(cal.document_url, ''),
		       cal.issued_by, COALESCE(u_iss.first_name || ' ' || u_iss.last_name, ''), cal.issued_at,
		       cal.extended_by, COALESCE(u_ext.first_name || ' ' || u_ext.last_name, ''), cal.extended_at,
		       cal.previous_expiry_date, COALESCE(cal.extension_reason, ''),
		       cal.revoked_by, COALESCE(u_rev.first_name || ' ' || u_rev.last_name, ''), cal.revoked_at,
		       COALESCE(cal.revocation_reason, ''),
		       cal.created_at, cal.updated_at
		FROM club_academy_licenses cal
		LEFT JOIN clubs c ON c.id = cal.club_id
		LEFT JOIN federations f ON f.id = c.federation_id
		LEFT JOIN users u_iss ON u_iss.id = cal.issued_by
		LEFT JOIN users u_ext ON u_ext.id = cal.extended_by
		LEFT JOIN users u_rev ON u_rev.id = cal.revoked_by
		WHERE cal.id = $1
	`
	row := r.db.QueryRow(ctx, q, id)
	lic, err := r.scanClubLicense(row)
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

func (r *ClubAcademyLicenseRepo) GetActiveLicenseByClubID(ctx context.Context, clubID string) (*models.ClubAcademyLicense, error) {
	const q = `
		SELECT cal.id, cal.club_id,
		       COALESCE(c.name, ''), COALESCE(c.club_number, ''), COALESCE(c.acronym, ''),
		       COALESCE(c.district, ''), COALESCE(c.region, ''), COALESCE(c.email, ''), COALESCE(c.phone, ''),
		       COALESCE(c.logo_url, ''),
		       COALESCE(c.president, c.contact_person, ''),
		       COALESCE(c.secretary, ''),
		       COALESCE(c.federation_id, ''), COALESCE(f.name, ''),
		       cal.license_number, cal.license_type, cal.category,
		       cal.issue_date, cal.expiry_date, cal.status,
		       COALESCE(cal.conditions, ''), COALESCE(cal.document_url, ''),
		       cal.issued_by, COALESCE(u_iss.first_name || ' ' || u_iss.last_name, ''), cal.issued_at,
		       cal.extended_by, COALESCE(u_ext.first_name || ' ' || u_ext.last_name, ''), cal.extended_at,
		       cal.previous_expiry_date, COALESCE(cal.extension_reason, ''),
		       cal.revoked_by, COALESCE(u_rev.first_name || ' ' || u_rev.last_name, ''), cal.revoked_at,
		       COALESCE(cal.revocation_reason, ''),
		       cal.created_at, cal.updated_at
		FROM club_academy_licenses cal
		LEFT JOIN clubs c ON c.id = cal.club_id
		LEFT JOIN federations f ON f.id = c.federation_id
		LEFT JOIN users u_iss ON u_iss.id = cal.issued_by
		LEFT JOIN users u_ext ON u_ext.id = cal.extended_by
		LEFT JOIN users u_rev ON u_rev.id = cal.revoked_by
		WHERE cal.club_id = $1
		ORDER BY cal.created_at DESC
		LIMIT 1
	`
	row := r.db.QueryRow(ctx, q, clubID)
	lic, err := r.scanClubLicense(row)
	if err != nil {
		return nil, err
	}
	logs, err := r.GetLicenseLogs(ctx, lic.ID)
	if err == nil {
		lic.Logs = logs
	}
	return lic, nil
}

func (r *ClubAcademyLicenseRepo) List(ctx context.Context, f ListClubLicensesFilter, p *models.PaginationParams) ([]*models.ClubAcademyLicense, int, error) {
	where := []string{"1=1"}
	args := []any{}
	idx := 1

	if f.ClubID != "" {
		where = append(where, fmt.Sprintf("cal.club_id = $%d", idx))
		args = append(args, f.ClubID)
		idx++
	}
	if f.FederationID != "" {
		where = append(where, fmt.Sprintf("c.federation_id = $%d", idx))
		args = append(args, f.FederationID)
		idx++
	}
	if f.Status != "" {
		where = append(where, fmt.Sprintf("UPPER(cal.status) = UPPER($%d)", idx))
		args = append(args, f.Status)
		idx++
	}
	if f.LicenseType != "" {
		where = append(where, fmt.Sprintf("UPPER(cal.license_type) = UPPER($%d)", idx))
		args = append(args, f.LicenseType)
		idx++
	}
	if f.Category != "" {
		where = append(where, fmt.Sprintf("UPPER(cal.category) = UPPER($%d)", idx))
		args = append(args, f.Category)
		idx++
	}
	if f.Search != "" {
		s := "%" + strings.ToLower(f.Search) + "%"
		where = append(where, fmt.Sprintf("(LOWER(cal.license_number) LIKE $%d OR LOWER(c.name) LIKE $%d OR LOWER(COALESCE(c.club_number, '')) LIKE $%d OR LOWER(COALESCE(c.acronym, '')) LIKE $%d OR LOWER(COALESCE(c.president, c.contact_person, '')) LIKE $%d OR LOWER(COALESCE(c.secretary, '')) LIKE $%d)", idx, idx, idx, idx, idx, idx))
		args = append(args, s)
		idx++
	}
	if f.StartDate != "" {
		where = append(where, fmt.Sprintf("cal.issue_date >= $%d::date", idx))
		args = append(args, f.StartDate)
		idx++
	}
	if f.EndDate != "" {
		where = append(where, fmt.Sprintf("cal.expiry_date <= $%d::date", idx))
		args = append(args, f.EndDate)
		idx++
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQ := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM club_academy_licenses cal
		LEFT JOIN clubs c ON c.id = cal.club_id
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

	orderCol := "cal.created_at"
	switch strings.ToLower(f.SortBy) {
	case "license_number", "number":
		orderCol = "cal.license_number"
	case "club_number", "academy_number":
		orderCol = "c.club_number"
	case "club", "name", "club_name", "academy_name":
		orderCol = "c.name"
	case "issue_date":
		orderCol = "cal.issue_date"
	case "expiry_date":
		orderCol = "cal.expiry_date"
	case "status":
		orderCol = "cal.status"
	case "type", "license_type":
		orderCol = "cal.license_type"
	case "category":
		orderCol = "cal.category"
	}

	orderDir := "DESC"
	if strings.EqualFold(f.SortOrder, "asc") {
		orderDir = "ASC"
	}

	args = append(args, limit, offset)
	q := fmt.Sprintf(`
		SELECT cal.id, cal.club_id,
		       COALESCE(c.name, ''), COALESCE(c.club_number, ''), COALESCE(c.acronym, ''),
		       COALESCE(c.district, ''), COALESCE(c.region, ''), COALESCE(c.email, ''), COALESCE(c.phone, ''),
		       COALESCE(c.logo_url, ''),
		       COALESCE(c.president, c.contact_person, ''),
		       COALESCE(c.secretary, ''),
		       COALESCE(c.federation_id, ''), COALESCE(f.name, ''),
		       cal.license_number, cal.license_type, cal.category,
		       cal.issue_date, cal.expiry_date, cal.status,
		       COALESCE(cal.conditions, ''), COALESCE(cal.document_url, ''),
		       cal.issued_by, COALESCE(u_iss.first_name || ' ' || u_iss.last_name, ''), cal.issued_at,
		       cal.extended_by, COALESCE(u_ext.first_name || ' ' || u_ext.last_name, ''), cal.extended_at,
		       cal.previous_expiry_date, COALESCE(cal.extension_reason, ''),
		       cal.revoked_by, COALESCE(u_rev.first_name || ' ' || u_rev.last_name, ''), cal.revoked_at,
		       COALESCE(cal.revocation_reason, ''),
		       cal.created_at, cal.updated_at
		FROM club_academy_licenses cal
		LEFT JOIN clubs c ON c.id = cal.club_id
		LEFT JOIN federations f ON f.id = c.federation_id
		LEFT JOIN users u_iss ON u_iss.id = cal.issued_by
		LEFT JOIN users u_ext ON u_ext.id = cal.extended_by
		LEFT JOIN users u_rev ON u_rev.id = cal.revoked_by
		WHERE %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderCol, orderDir, idx, idx+1)

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*models.ClubAcademyLicense
	for rows.Next() {
		lic, err := r.scanClubLicense(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, lic)
	}
	return list, total, nil
}

func (r *ClubAcademyLicenseRepo) Extend(ctx context.Context, id string, newExpiry time.Time, reason, notes, actorID, actorName string) (*models.ClubAcademyLicense, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var current models.ClubAcademyLicense
	const fetchQ = `SELECT id, status, expiry_date FROM club_academy_licenses WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, fetchQ, id).Scan(&current.ID, &current.Status, &current.ExpiryDate)
	if err != nil {
		return nil, err
	}

	const updateQ = `
		UPDATE club_academy_licenses
		SET previous_expiry_date = expiry_date,
		    expiry_date = $1,
		    status = 'EXTENDED',
		    extended_by = $2,
		    extended_at = NOW(),
		    extension_reason = $3,
		    updated_at = NOW()
		WHERE id = $4
	`
	_, err = tx.Exec(ctx, updateQ, newExpiry, actorID, reason, id)
	if err != nil {
		return nil, err
	}

	const logQ = `
		INSERT INTO club_academy_license_logs (
			id, license_id, action, performed_by, performed_by_name, reason,
			old_status, new_status, old_expiry_date, new_expiry_date, notes, created_at
		) VALUES (
			$1, $2, 'EXTENDED', $3, $4, $5,
			$6, 'EXTENDED', $7, $8, $9, NOW()
		)
	`
	_, err = tx.Exec(ctx, logQ,
		uuid.NewString(), id, actorID, actorName, reason,
		current.Status, current.ExpiryDate, newExpiry, notes,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *ClubAcademyLicenseRepo) Revoke(ctx context.Context, id string, reason, notes, actorID, actorName string) (*models.ClubAcademyLicense, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var current models.ClubAcademyLicense
	const fetchQ = `SELECT id, status, expiry_date FROM club_academy_licenses WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, fetchQ, id).Scan(&current.ID, &current.Status, &current.ExpiryDate)
	if err != nil {
		return nil, err
	}

	const updateQ = `
		UPDATE club_academy_licenses
		SET status = 'REVOKED',
		    revoked_by = $1,
		    revoked_at = NOW(),
		    revocation_reason = $2,
		    updated_at = NOW()
		WHERE id = $3
	`
	_, err = tx.Exec(ctx, updateQ, actorID, reason, id)
	if err != nil {
		return nil, err
	}

	const logQ = `
		INSERT INTO club_academy_license_logs (
			id, license_id, action, performed_by, performed_by_name, reason,
			old_status, new_status, old_expiry_date, new_expiry_date, notes, created_at
		) VALUES (
			$1, $2, 'REVOKED', $3, $4, $5,
			$6, 'REVOKED', $7, $7, $8, NOW()
		)
	`
	_, err = tx.Exec(ctx, logQ,
		uuid.NewString(), id, actorID, actorName, reason,
		current.Status, current.ExpiryDate, notes,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *ClubAcademyLicenseRepo) Reinstate(ctx context.Context, id string, reason, notes, actorID, actorName string) (*models.ClubAcademyLicense, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var current models.ClubAcademyLicense
	const fetchQ = `SELECT id, status, expiry_date FROM club_academy_licenses WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, fetchQ, id).Scan(&current.ID, &current.Status, &current.ExpiryDate)
	if err != nil {
		return nil, err
	}

	const updateQ = `
		UPDATE club_academy_licenses
		SET status = 'ACTIVE',
		    revoked_by = NULL,
		    revoked_at = NULL,
		    revocation_reason = '',
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, updateQ, id)
	if err != nil {
		return nil, err
	}

	const logQ = `
		INSERT INTO club_academy_license_logs (
			id, license_id, action, performed_by, performed_by_name, reason,
			old_status, new_status, old_expiry_date, new_expiry_date, notes, created_at
		) VALUES (
			$1, $2, 'REINSTATED', $3, $4, $5,
			$6, 'ACTIVE', $7, $7, $8, NOW()
		)
	`
	_, err = tx.Exec(ctx, logQ,
		uuid.NewString(), id, actorID, actorName, reason,
		current.Status, current.ExpiryDate, notes,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *ClubAcademyLicenseRepo) GetLicenseLogs(ctx context.Context, licenseID string) ([]*models.ClubAcademyLicenseLog, error) {
	const q = `
		SELECT id, license_id, action, performed_by, performed_by_name,
		       reason, old_status, new_status, old_expiry_date, new_expiry_date,
		       notes, created_at
		FROM club_academy_license_logs
		WHERE license_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, q, licenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.ClubAcademyLicenseLog
	for rows.Next() {
		var l models.ClubAcademyLicenseLog
		if err := rows.Scan(
			&l.ID, &l.LicenseID, &l.Action, &l.PerformedBy, &l.PerformedByName,
			&l.Reason, &l.OldStatus, &l.NewStatus, &l.OldExpiryDate, &l.NewExpiryDate,
			&l.Notes, &l.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, nil
}

func (r *ClubAcademyLicenseRepo) GetKPIs(ctx context.Context) (*models.ClubLicenseKPIs, error) {
	const q = `
		SELECT 
			COUNT(*) AS total_licenses,
			COUNT(*) FILTER (WHERE status = 'ACTIVE') AS active_licenses,
			COUNT(*) FILTER (WHERE status = 'EXTENDED') AS extended_licenses,
			COUNT(*) FILTER (WHERE status = 'REVOKED') AS revoked_licenses,
			COUNT(*) FILTER (WHERE status IN ('ACTIVE', 'EXTENDED') AND expiry_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '30 days') AS expiring_soon
		FROM club_academy_licenses
	`
	var kpi models.ClubLicenseKPIs
	err := r.db.QueryRow(ctx, q).Scan(
		&kpi.TotalLicenses, &kpi.ActiveLicenses, &kpi.ExtendedLicenses,
		&kpi.RevokedLicenses, &kpi.ExpiringSoon,
	)
	if err != nil {
		return nil, err
	}
	return &kpi, nil
}

func (r *ClubAcademyLicenseRepo) scanClubLicense(s rowScanner) (*models.ClubAcademyLicense, error) {
	var lic models.ClubAcademyLicense
	err := s.Scan(
		&lic.ID, &lic.ClubID,
		&lic.ClubName, &lic.ClubNumber, &lic.ClubAcronym,
		&lic.ClubDistrict, &lic.ClubRegion, &lic.ClubEmail, &lic.ClubPhone,
		&lic.ClubLogo, &lic.PresidentName, &lic.SecretaryName,
		&lic.FederationID, &lic.FederationName,
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
