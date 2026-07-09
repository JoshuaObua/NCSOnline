package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

type OrganisationRepo struct{ db *pgxpool.Pool }

func (r *OrganisationRepo) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM organisations WHERE status <> 'ARCHIVED'`).Scan(&count)
	return count, err
}

func (r *OrganisationRepo) ListForUser(ctx context.Context, userID string) ([]models.Organisation, error) {
	const q = `SELECT o.id,o.profile_reference,o.organisation_type,o.legal_name,o.display_name,
		o.official_email,o.official_phone,o.registration_number,o.status,m.role,o.created_at,o.updated_at
		FROM organisations o JOIN organisation_memberships m ON m.organisation_id=o.id
		WHERE m.user_id=$1 AND m.status='ACTIVE' AND o.status='ACTIVE' ORDER BY o.display_name`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.Organisation{}
	for rows.Next() {
		var o models.Organisation
		if err := rows.Scan(&o.ID, &o.ProfileReference, &o.OrganisationType, &o.LegalName, &o.DisplayName,
			&o.OfficialEmail, &o.OfficialPhone, &o.RegistrationNumber, &o.Status, &o.Role, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, o)
	}
	return items, rows.Err()
}

func (r *OrganisationRepo) GetForUser(ctx context.Context, organisationID, userID string) (*models.Organisation, error) {
	const q = `SELECT o.id,o.profile_reference,o.organisation_type,o.legal_name,o.display_name,
		o.official_email,o.official_phone,o.registration_number,o.status,m.role,o.profile_data,o.created_at,o.updated_at
		FROM organisations o JOIN organisation_memberships m ON m.organisation_id=o.id
		WHERE o.id=$1 AND m.user_id=$2 AND m.status='ACTIVE'`
	var o models.Organisation
	err := r.db.QueryRow(ctx, q, organisationID, userID).Scan(&o.ID, &o.ProfileReference, &o.OrganisationType, &o.LegalName, &o.DisplayName,
		&o.OfficialEmail, &o.OfficialPhone, &o.RegistrationNumber, &o.Status, &o.Role, &o.ProfileData, &o.CreatedAt, &o.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &o, err
}

// ApproveAndProvision is idempotent and keeps approval, profile creation,
// ownership, invitation and notification in one transaction.
func (r *OrganisationRepo) ApproveAndProvision(ctx context.Context, app *models.Application, reviewerID, notes string) (*models.Organisation, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var existingID *string
	if err := tx.QueryRow(ctx, `SELECT provisioned_organisation_id FROM applications WHERE id=$1 FOR UPDATE`, app.ID).Scan(&existingID); err != nil {
		return nil, err
	}
	if existingID != nil {
		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
		return r.getByID(ctx, *existingID)
	}
	values := flattenForm(app.FormData)
	var applicantEmail string
	if err := tx.QueryRow(ctx, `SELECT email FROM users WHERE id=$1`, app.UserID).Scan(&applicantEmail); err != nil {
		return nil, err
	}
	name := firstValue(values, "legal_name", "organisation_name", "organization_name", "federation_name", "association_name", "club_name", "name")
	if name == "" {
		name = "Approved sports organisation " + app.ApplicationReference
	}
	email := firstValue(values, "official_email", "organisation_email", "organization_email", "email")
	if email == "" {
		email = applicantEmail
	}
	phone := firstValue(values, "official_phone", "phone", "telephone", "contact_phone")
	regNo := firstValue(values, "registration_number", "certificate_number", "ncs_registration_number")
	typeName := organisationType(app)
	id := uuid.NewString()
	reference := fmt.Sprintf("NCS-%s-%s", typeName[:3], strings.ToUpper(strings.ReplaceAll(uuid.NewString()[:8], "-", "")))
	_, err = tx.Exec(ctx, `INSERT INTO organisations(id,profile_reference,organisation_type,legal_name,display_name,official_email,official_phone,registration_number,source_application_id,profile_data,created_by)
		VALUES($1,$2,$3,$4,$4,$5,$6,$7,$8,$9,$10)`, id, reference, typeName, name, strings.ToLower(email), phone, regNo, app.ID, app.FormData, reviewerID)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(ctx, `INSERT INTO organisation_memberships(id,organisation_id,user_id,role,status,is_primary_owner) VALUES($1,$2,$3,'OWNER','ACTIVE',TRUE)`, uuid.NewString(), id, app.UserID)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(email, applicantEmail) {
		rawToken, tokenHash, err := invitationToken()
		if err != nil {
			return nil, err
		}
		invitationID := uuid.NewString()
		_, err = tx.Exec(ctx, `INSERT INTO organisation_invitations(id,organisation_id,email,role,token_hash,expires_at) VALUES($1,$2,$3,'ADMIN',$4,NOW()+INTERVAL '72 hours')`, invitationID, id, strings.ToLower(email), tokenHash)
		if err != nil {
			return nil, err
		}
		payload, _ := json.Marshal(map[string]string{"organisation_name": name, "activation_path": "/accept-organisation-invite?token=" + rawToken})
		_, err = tx.Exec(ctx, `INSERT INTO notification_deliveries(id,template_code,recipient_address,channel,related_type,related_id,payload) VALUES($1,'ORGANISATION_INVITATION',$2,'EMAIL','ORGANISATION_INVITATION',$3,$4)`, uuid.NewString(), strings.ToLower(email), invitationID, payload)
		if err != nil {
			return nil, err
		}
	}
	_, err = tx.Exec(ctx, `UPDATE applications SET status='APPROVED',reviewer_id=$2,review_notes=$3,approved_at=NOW(),provisioning_status='PROVISIONED',provisioned_organisation_id=$4,provisioning_attempts=provisioning_attempts+1,updated_at=NOW() WHERE id=$1`, app.ID, reviewerID, notes, id)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.getByID(ctx, id)
}

func (r *OrganisationRepo) AcceptInvitation(ctx context.Context, rawToken, userID, userEmail string) error {
	h := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(h[:])
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var invitationID, organisationID, email, role, status string
	var expires time.Time
	err = tx.QueryRow(ctx, `SELECT id,organisation_id,email,role,status,expires_at FROM organisation_invitations WHERE token_hash=$1 FOR UPDATE`, tokenHash).Scan(&invitationID, &organisationID, &email, &role, &status, &expires)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if status != "PENDING" || time.Now().After(expires) || !strings.EqualFold(email, userEmail) {
		return fmt.Errorf("invitation is invalid, expired, or belongs to another account")
	}
	_, err = tx.Exec(ctx, `INSERT INTO organisation_memberships(id,organisation_id,user_id,role,status) VALUES($1,$2,$3,$4,'ACTIVE') ON CONFLICT(organisation_id,user_id) DO UPDATE SET role=EXCLUDED.role,status='ACTIVE',updated_at=NOW()`, uuid.NewString(), organisationID, userID, role)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE organisation_invitations SET status='ACCEPTED',accepted_by=$2,accepted_at=NOW() WHERE id=$1`, invitationID, userID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *OrganisationRepo) getByID(ctx context.Context, id string) (*models.Organisation, error) {
	var o models.Organisation
	err := r.db.QueryRow(ctx, `SELECT id,profile_reference,organisation_type,legal_name,display_name,official_email,official_phone,registration_number,status,profile_data,created_at,updated_at FROM organisations WHERE id=$1`, id).Scan(&o.ID, &o.ProfileReference, &o.OrganisationType, &o.LegalName, &o.DisplayName, &o.OfficialEmail, &o.OfficialPhone, &o.RegistrationNumber, &o.Status, &o.ProfileData, &o.CreatedAt, &o.UpdatedAt)
	return &o, err
}

func organisationType(app *models.Application) string {
	v := strings.ToUpper(app.OrganisationType)
	if strings.Contains(v, "CLUB") || app.FormType == "form_10" {
		return "CLUB"
	}
	if strings.Contains(v, "ASSOCIATION") {
		return "ASSOCIATION"
	}
	return "FEDERATION"
}
func flattenForm(raw json.RawMessage) map[string]string {
	var v any
	_ = json.Unmarshal(raw, &v)
	out := map[string]string{}
	var walk func(any)
	walk = func(x any) {
		switch t := x.(type) {
		case map[string]any:
			for k, v := range t {
				if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
					out[strings.ToLower(k)] = strings.TrimSpace(s)
				}
				walk(v)
			}
		case []any:
			for _, v := range t {
				walk(v)
			}
		}
	}
	walk(v)
	return out
}
func firstValue(values map[string]string, keys ...string) string {
	for _, k := range keys {
		if v := values[k]; v != "" {
			return v
		}
	}
	return ""
}
func invitationToken() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	raw := hex.EncodeToString(b)
	h := sha256.Sum256([]byte(raw))
	return raw, hex.EncodeToString(h[:]), nil
}
