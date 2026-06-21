package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type USSDRepo struct{ db *pgxpool.Pool }

type USSDAthlete struct {
	Number, FullName, Status, Discipline, Federation string
	VerifiedAt                                       *time.Time
}

type USSDLicence struct {
	Number, HolderName, Type, Status string
	ExpiresOn                        *time.Time
	Renewable                        bool
}

func (r *USSDRepo) FindActiveUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	const q = `SELECT id,email,COALESCE(phone,''),COALESCE(pin_hash,''),pin_change_required,
	 is_active,account_status,first_name,last_name
	 FROM users WHERE deleted_at IS NULL AND is_active AND account_status='ACTIVE'
	 AND CASE
	   WHEN regexp_replace(COALESCE(phone,''),'[^0-9]','','g') LIKE '0%' THEN '+256'||substring(regexp_replace(phone,'[^0-9]','','g') from 2)
	   WHEN regexp_replace(COALESCE(phone,''),'[^0-9]','','g') LIKE '256%' THEN '+'||regexp_replace(phone,'[^0-9]','','g')
	   ELSE COALESCE(phone,'') END=$1 LIMIT 1`
	u := &models.User{}
	err := r.db.QueryRow(ctx, q, phone).Scan(&u.ID, &u.Email, &u.Phone, &u.PinHash, &u.PinChangeRequired, &u.IsActive, &u.AccountStatus, &u.FirstName, &u.LastName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

func (r *USSDRepo) FindAthlete(ctx context.Context, number string) (*USSDAthlete, error) {
	const q = `SELECT a.athlete_number,a.full_name,a.status,a.discipline,COALESCE(f.name,''),a.verified_at
	 FROM athletes a LEFT JOIN athlete_affiliations aa ON aa.athlete_id=a.id AND aa.is_active AND aa.ends_on IS NULL
	 LEFT JOIN federations f ON f.id=aa.federation_id
	 WHERE a.deleted_at IS NULL AND UPPER(a.athlete_number)=UPPER($1) LIMIT 1`
	a := &USSDAthlete{}
	err := r.db.QueryRow(ctx, q, strings.TrimSpace(number)).Scan(&a.Number, &a.FullName, &a.Status, &a.Discipline, &a.Federation, &a.VerifiedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return a, err
}

func (r *USSDRepo) FindLicence(ctx context.Context, number string) (*USSDLicence, error) {
	const q = `SELECT licence_number,holder_name,licence_type,status,expires_on,renewable FROM (
	 SELECT licence_number,holder_name,licence_type,
	   CASE WHEN status='VALID' AND expires_on<CURRENT_DATE THEN 'EXPIRED' ELSE status END status,
	   expires_on,renewable,1 priority FROM licences WHERE UPPER(licence_number)=UPPER($1)
	 UNION ALL
	 SELECT license_number,full_name,'COACH',CASE WHEN status='ACTIVE' AND expiry_date<CURRENT_DATE THEN 'EXPIRED' ELSE status END,expiry_date,TRUE,2
	   FROM coaches WHERE UPPER(license_number)=UPPER($1)
	 UNION ALL
	 SELECT certification,full_name,'TECHNICAL_OFFICIAL',CASE WHEN status='ACTIVE' AND valid_until<CURRENT_DATE THEN 'EXPIRED' ELSE status END,valid_until,TRUE,3
	   FROM technical_officials WHERE UPPER(certification)=UPPER($1)
	 ) credentials ORDER BY priority LIMIT 1`
	l := &USSDLicence{}
	err := r.db.QueryRow(ctx, q, strings.TrimSpace(number)).Scan(&l.Number, &l.HolderName, &l.Type, &l.Status, &l.ExpiresOn, &l.Renewable)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return l, err
}
