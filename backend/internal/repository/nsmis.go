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

type GovernanceDraftInput struct {
	ExecutiveMeetingHeld     bool
	ExecutiveMeetingDate     *time.Time
	AGMConducted             bool
	AGMDate                  *time.Time
	BoardMeetingHeld         bool
	BoardMeetingDate         *time.Time
	ElectionsConducted       bool
	ElectionDate             *time.Time
	DisciplinaryCasesHandled int
}

type NSMISRepo struct{ db *pgxpool.Pool }

func (r *NSMISRepo) HasAnyPermission(ctx context.Context, userID string, names ...string) (bool, error) {
	if len(names) == 0 {
		return false, nil
	}
	var ok bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM user_roles ur JOIN role_permissions rp ON rp.role_id=ur.role_id JOIN permissions p ON p.id=rp.permission_id WHERE ur.user_id=$1 AND p.name=ANY($2))`, userID, names).Scan(&ok)
	return ok, err
}

func (r *NSMISRepo) ListFederations(ctx context.Context, userID string, all bool) ([]models.Federation, error) {
	const q = `SELECT DISTINCT f.id,f.name,f.acronym,f.ncs_registration_number,f.recognition_status,
	                  f.physical_address,f.email,f.website,f.contact_person,f.is_active,f.version,f.created_at,f.updated_at
	           FROM federations f
	           LEFT JOIN federation_memberships fm ON fm.federation_id=f.id AND fm.user_id=$1
	             AND fm.is_active AND (fm.ends_at IS NULL OR fm.ends_at >= CURRENT_DATE)
	           WHERE f.deleted_at IS NULL AND ($2 OR fm.id IS NOT NULL)
	           ORDER BY f.name`
	rows, err := r.db.Query(ctx, q, userID, all)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.Federation, 0)
	for rows.Next() {
		var f models.Federation
		if err := rows.Scan(&f.ID, &f.Name, &f.Acronym, &f.NCSRegistrationNumber, &f.RecognitionStatus,
			&f.PhysicalAddress, &f.Email, &f.Website, &f.ContactPerson, &f.IsActive, &f.Version, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	return items, rows.Err()
}

func (r *NSMISRepo) CreateFederation(ctx context.Context, f *models.Federation, actorID string) error {
	const q = `INSERT INTO federations(name,acronym,ncs_registration_number,recognition_status,physical_address,email,website,contact_person,created_by,updated_by)
	           VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$9)
	           RETURNING id,is_active,version,created_at,updated_at`
	err := r.db.QueryRow(ctx, q, f.Name, f.Acronym, f.NCSRegistrationNumber, f.RecognitionStatus,
		f.PhysicalAddress, f.Email, f.Website, f.ContactPerson, actorID).Scan(&f.ID, &f.IsActive, &f.Version, &f.CreatedAt, &f.UpdatedAt)
	if isDuplicate(err) {
		return ErrDuplicate
	}
	return err
}

func (r *NSMISRepo) UpdateFederation(ctx context.Context, f *models.Federation, actorID string, all bool) error {
	q := `UPDATE federations f SET name=$2,acronym=$3,ncs_registration_number=$4,recognition_status=$5,physical_address=$6,email=$7,website=$8,contact_person=$9,version=version+1,updated_by=$10,updated_at=NOW() WHERE id=$1 AND version=$11 AND ($12 OR EXISTS(SELECT 1 FROM federation_memberships fm WHERE fm.federation_id=f.id AND fm.user_id=$10 AND fm.is_active AND (fm.ends_at IS NULL OR fm.ends_at>=CURRENT_DATE))) RETURNING version,updated_at`
	err := r.db.QueryRow(ctx, q, f.ID, f.Name, f.Acronym, f.NCSRegistrationNumber, f.RecognitionStatus, f.PhysicalAddress, f.Email, f.Website, f.ContactPerson, actorID, f.Version, all).Scan(&f.Version, &f.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if isDuplicate(err) {
		return ErrDuplicate
	}
	return err
}

func (r *NSMISRepo) ListPeriods(ctx context.Context, openOnly bool) ([]models.ReportingPeriod, error) {
	const q = `SELECT id,period_type,name,starts_on,ends_on,due_on,grace_ends_on,timezone,is_open,created_at
	           FROM reporting_periods WHERE (NOT $1 OR is_open) ORDER BY starts_on DESC,period_type`
	rows, err := r.db.Query(ctx, q, openOnly)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.ReportingPeriod, 0)
	for rows.Next() {
		var p models.ReportingPeriod
		if err := rows.Scan(&p.ID, &p.PeriodType, &p.Name, &p.StartsOn, &p.EndsOn, &p.DueOn, &p.GraceEndsOn, &p.Timezone, &p.IsOpen, &p.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

func (r *NSMISRepo) CreatePeriod(ctx context.Context, p *models.ReportingPeriod) error {
	const q = `INSERT INTO reporting_periods(period_type,name,starts_on,ends_on,due_on,grace_ends_on,timezone,is_open)
	         VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id,created_at`
	err := r.db.QueryRow(ctx, q, p.PeriodType, p.Name, p.StartsOn, p.EndsOn, p.DueOn, p.GraceEndsOn, p.Timezone, p.IsOpen).Scan(&p.ID, &p.CreatedAt)
	if isDuplicate(err) {
		return ErrDuplicate
	}
	return err
}

func (r *NSMISRepo) GenerateObligations(ctx context.Context, periodID, reportType string) (int64, error) {
	const q = `INSERT INTO report_obligations(federation_id,reporting_period_id,report_type,due_on)
	         SELECT f.id,p.id,$2,p.due_on FROM federations f CROSS JOIN reporting_periods p
	         WHERE p.id=$1 AND f.deleted_at IS NULL AND f.is_active
	         ON CONFLICT(federation_id,reporting_period_id,report_type) DO NOTHING`
	result, err := r.db.Exec(ctx, q, periodID, reportType)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (r *NSMISRepo) ListObligations(ctx context.Context, userID, periodID string, all bool) ([]models.ReportObligation, error) {
	const q = `SELECT DISTINCT o.id,o.federation_id,f.name,o.reporting_period_id,p.name,o.report_type,o.due_on,
	                CASE WHEN o.status IN ('NOT_STARTED','DRAFT') AND o.due_on<CURRENT_DATE THEN 'OVERDUE' ELSE o.status END,
	                GREATEST(CURRENT_DATE-o.due_on,0),current_report.id,current_report.version,current_report.status
	         FROM report_obligations o JOIN federations f ON f.id=o.federation_id
	         JOIN reporting_periods p ON p.id=o.reporting_period_id
	         LEFT JOIN LATERAL (SELECT fr.id,fr.version,fr.status FROM federation_reports fr WHERE fr.obligation_id=o.id ORDER BY fr.revision DESC LIMIT 1) current_report ON TRUE
	         LEFT JOIN federation_memberships fm ON fm.federation_id=f.id AND fm.user_id=$1 AND fm.is_active
	           AND (fm.ends_at IS NULL OR fm.ends_at>=CURRENT_DATE)
	         WHERE ($2='' OR o.reporting_period_id=$2) AND ($3 OR fm.id IS NOT NULL)
	         ORDER BY o.due_on,o.report_type,f.name`
	rows, err := r.db.Query(ctx, q, userID, periodID, all)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]models.ReportObligation, 0)
	for rows.Next() {
		var o models.ReportObligation
		if err := rows.Scan(&o.ID, &o.FederationID, &o.FederationName, &o.ReportingPeriodID, &o.PeriodName, &o.ReportType, &o.DueOn, &o.Status, &o.DaysOverdue, &o.ReportID, &o.ReportVersion, &o.ReportStatus); err != nil {
			return nil, err
		}
		items = append(items, o)
	}
	return items, rows.Err()
}

func (r *NSMISRepo) GovernanceDashboard(ctx context.Context, userID, periodID string, all bool) (*models.GovernanceDashboard, error) {
	if periodID == "" {
		err := r.db.QueryRow(ctx, `SELECT id FROM reporting_periods WHERE is_open ORDER BY ends_on DESC LIMIT 1`).Scan(&periodID)
		if errors.Is(err, pgx.ErrNoRows) {
			return &models.GovernanceDashboard{Rows: []models.GovernanceFederationRow{}, AsOf: time.Now().UTC(), CalculationVersion: 1}, nil
		}
		if err != nil {
			return nil, err
		}
	}
	const q = `SELECT DISTINCT f.id,f.name,f.acronym,cs.score,cs.is_compliant,
	                (SELECT COUNT(*) FROM report_obligations o WHERE o.federation_id=f.id AND o.reporting_period_id=$2
	                  AND o.status NOT IN ('SUBMITTED','ACCEPTED','EXEMPT') AND o.due_on<CURRENT_DATE),
	                EXISTS(SELECT 1 FROM federation_documents d WHERE d.federation_id=f.id AND d.document_type='CONSTITUTION'
	                  AND d.expires_on<CURRENT_DATE)
	         FROM federations f
	         LEFT JOIN federation_memberships fm ON fm.federation_id=f.id AND fm.user_id=$1 AND fm.is_active
	           AND (fm.ends_at IS NULL OR fm.ends_at>=CURRENT_DATE)
	         LEFT JOIN compliance_scores cs ON cs.federation_id=f.id AND cs.reporting_period_id=$2
	         WHERE f.deleted_at IS NULL AND f.is_active AND ($3 OR fm.id IS NOT NULL) ORDER BY f.name`
	rows, err := r.db.Query(ctx, q, userID, periodID, all)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	d := &models.GovernanceDashboard{ReportingPeriodID: periodID, Rows: make([]models.GovernanceFederationRow, 0), AsOf: time.Now().UTC(), CalculationVersion: 1}
	_ = r.db.QueryRow(ctx, `SELECT name FROM reporting_periods WHERE id=$1`, periodID).Scan(&d.ReportingPeriodName)
	var scoreTotal float64
	var scoreCount int64
	for rows.Next() {
		var x models.GovernanceFederationRow
		if err := rows.Scan(&x.FederationID, &x.FederationName, &x.Acronym, &x.ComplianceScore, &x.IsCompliant, &x.MissingReports, &x.ExpiredConstitution); err != nil {
			return nil, err
		}
		d.TotalFederations++
		d.MissingReports += x.MissingReports
		if x.ExpiredConstitution {
			d.ExpiredConstitutions++
		}
		if x.IsCompliant != nil {
			if *x.IsCompliant {
				d.CompliantFederations++
			} else {
				d.NonCompliant++
			}
		}
		if x.ComplianceScore != nil {
			scoreTotal += *x.ComplianceScore
			scoreCount++
		}
		d.Rows = append(d.Rows, x)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if scoreCount > 0 {
		avg := scoreTotal / float64(scoreCount)
		d.AverageScore = &avg
	}
	return d, nil
}

// SaveGovernanceDraft atomically creates/updates the current report and its typed response.
// Federation scope is checked in the same query used to select the obligation.
func (r *NSMISRepo) SaveGovernanceDraft(ctx context.Context, obligationID, actorID string, all bool, version int, in GovernanceDraftInput) (string, int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", 0, err
	}
	defer tx.Rollback(ctx)
	var federationID, periodID, reportType string
	err = tx.QueryRow(ctx, `SELECT o.federation_id,o.reporting_period_id,o.report_type FROM report_obligations o
	  LEFT JOIN federation_memberships fm ON fm.federation_id=o.federation_id AND fm.user_id=$2 AND fm.is_active
	    AND (fm.ends_at IS NULL OR fm.ends_at>=CURRENT_DATE)
	  WHERE o.id=$1 AND ($3 OR fm.id IS NOT NULL) FOR UPDATE`, obligationID, actorID, all).Scan(&federationID, &periodID, &reportType)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", 0, ErrNotFound
	}
	if err != nil {
		return "", 0, err
	}
	if reportType != "GOVERNANCE" {
		return "", 0, fmt.Errorf("obligation is not a governance report")
	}
	var reportID string
	var currentVersion int
	err = tx.QueryRow(ctx, `SELECT id,version FROM federation_reports WHERE obligation_id=$1 AND status='DRAFT' ORDER BY revision DESC LIMIT 1 FOR UPDATE`, obligationID).Scan(&reportID, &currentVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		err = tx.QueryRow(ctx, `INSERT INTO federation_reports(obligation_id,federation_id,reporting_period_id,report_type,created_by)
		 VALUES($1,$2,$3,'GOVERNANCE',$4) RETURNING id,version`, obligationID, federationID, periodID, actorID).Scan(&reportID, &currentVersion)
	} else if err == nil {
		if version > 0 && version != currentVersion {
			return "", currentVersion, fmt.Errorf("version conflict")
		}
		currentVersion++
		_, err = tx.Exec(ctx, `UPDATE federation_reports SET version=$2,updated_at=NOW() WHERE id=$1`, reportID, currentVersion)
	}
	if err != nil {
		return "", 0, err
	}
	_, err = tx.Exec(ctx, `INSERT INTO governance_responses(report_id,executive_meeting_held,executive_meeting_date,agm_conducted,agm_date,board_meeting_held,board_meeting_date,elections_conducted,election_date,disciplinary_cases_handled)
	 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	 ON CONFLICT(report_id) DO UPDATE SET executive_meeting_held=EXCLUDED.executive_meeting_held,executive_meeting_date=EXCLUDED.executive_meeting_date,
	 agm_conducted=EXCLUDED.agm_conducted,agm_date=EXCLUDED.agm_date,board_meeting_held=EXCLUDED.board_meeting_held,board_meeting_date=EXCLUDED.board_meeting_date,
	 elections_conducted=EXCLUDED.elections_conducted,election_date=EXCLUDED.election_date,disciplinary_cases_handled=EXCLUDED.disciplinary_cases_handled,updated_at=NOW()`,
		reportID, in.ExecutiveMeetingHeld, in.ExecutiveMeetingDate, in.AGMConducted, in.AGMDate, in.BoardMeetingHeld, in.BoardMeetingDate, in.ElectionsConducted, in.ElectionDate, in.DisciplinaryCasesHandled)
	if err != nil {
		return "", 0, err
	}
	_, err = tx.Exec(ctx, `UPDATE report_obligations SET status='DRAFT',updated_at=NOW() WHERE id=$1 AND status IN ('NOT_STARTED','DRAFT','OVERDUE')`, obligationID)
	if err != nil {
		return "", 0, err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", 0, err
	}
	return reportID, currentVersion, nil
}

func (r *NSMISRepo) TransitionReport(ctx context.Context, reportID, actorID, toStatus, reason string, all bool) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var fromStatus, createdBy, obligationID string
	err = tx.QueryRow(ctx, `SELECT r.status,r.created_by,r.obligation_id FROM federation_reports r
	 LEFT JOIN federation_memberships fm ON fm.federation_id=r.federation_id AND fm.user_id=$2 AND fm.is_active
	   AND (fm.ends_at IS NULL OR fm.ends_at>=CURRENT_DATE)
	 WHERE r.id=$1 AND ($3 OR fm.id IS NOT NULL) FOR UPDATE`, reportID, actorID, all).Scan(&fromStatus, &createdBy, &obligationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	allowed := map[string]map[string]bool{"DRAFT": {"SUBMITTED": true}, "SUBMITTED": {"PRESIDENT_APPROVED": true, "PRESIDENT_RETURNED": true}, "PRESIDENT_APPROVED": {"NCS_UNDER_REVIEW": true, "NEEDS_CORRECTION": true, "APPROVED": true}, "NCS_UNDER_REVIEW": {"NEEDS_CORRECTION": true, "APPROVED": true}, "APPROVED": {"LOCKED": true}}
	if !allowed[fromStatus][toStatus] {
		return fmt.Errorf("invalid transition from %s to %s", fromStatus, toStatus)
	}
	if toStatus == "PRESIDENT_APPROVED" && actorID == createdBy {
		return fmt.Errorf("self approval is not allowed")
	}
	set := `status=$3,updated_at=NOW()`
	switch toStatus {
	case "SUBMITTED":
		set += `,submitted_by=$2,submitted_at=NOW()`
	case "PRESIDENT_APPROVED":
		set += `,president_approved_by=$2,president_approved_at=NOW()`
	case "NCS_UNDER_REVIEW", "NEEDS_CORRECTION", "APPROVED":
		set += `,reviewed_by=$2,reviewed_at=NOW()`
	case "LOCKED":
		set += `,locked_at=NOW()`
	}
	if _, err = tx.Exec(ctx, `UPDATE federation_reports SET `+set+` WHERE id=$1`, reportID, actorID, toStatus); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO report_transitions(report_id,from_status,to_status,actor_id,reason) VALUES($1,$2,$3,$4,$5)`, reportID, fromStatus, toStatus, actorID, reason); err != nil {
		return err
	}
	obligationStatus := map[string]string{"SUBMITTED": "SUBMITTED", "APPROVED": "ACCEPTED", "LOCKED": "ACCEPTED", "NEEDS_CORRECTION": "DRAFT", "PRESIDENT_RETURNED": "DRAFT"}[toStatus]
	if obligationStatus != "" {
		if _, err = tx.Exec(ctx, `UPDATE report_obligations SET status=$2,updated_at=NOW() WHERE id=$1`, obligationID, obligationStatus); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *NSMISRepo) AthleteDashboard(ctx context.Context) (map[string]interface{}, error) {
	result := map[string]interface{}{"as_of": time.Now().UTC(), "gender": map[string]int64{}, "by_region": []map[string]interface{}{}, "by_federation": []map[string]interface{}{}}
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM athletes WHERE deleted_at IS NULL AND status='ACTIVE'`).Scan(&total); err != nil {
		return nil, err
	}
	result["total_registered_athletes"] = total
	gender := map[string]int64{}
	rows, err := r.db.Query(ctx, `SELECT gender,COUNT(*) FROM athletes WHERE deleted_at IS NULL AND status='ACTIVE' GROUP BY gender ORDER BY gender`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var k string
		var n int64
		if err := rows.Scan(&k, &n); err != nil {
			rows.Close()
			return nil, err
		}
		gender[k] = n
	}
	rows.Close()
	result["gender"] = gender
	regions := make([]map[string]interface{}, 0)
	rows, err = r.db.Query(ctx, `SELECT COALESCE(NULLIF(region,''),'NOT_RECORDED'),COUNT(*) FROM athletes WHERE deleted_at IS NULL AND status='ACTIVE' GROUP BY 1 ORDER BY 2 DESC,1 LIMIT 50`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var k string
		var n int64
		if err := rows.Scan(&k, &n); err != nil {
			rows.Close()
			return nil, err
		}
		regions = append(regions, map[string]interface{}{"region": k, "count": n})
	}
	rows.Close()
	result["by_region"] = regions
	byFed := make([]map[string]interface{}, 0)
	rows, err = r.db.Query(ctx, `SELECT f.id,f.name,COUNT(DISTINCT aa.athlete_id) FROM federations f LEFT JOIN athlete_affiliations aa ON aa.federation_id=f.id AND aa.is_active GROUP BY f.id,f.name ORDER BY 3 DESC,f.name`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id, name string
		var n int64
		if err := rows.Scan(&id, &name, &n); err != nil {
			rows.Close()
			return nil, err
		}
		byFed = append(byFed, map[string]interface{}{"federation_id": id, "federation_name": name, "count": n})
	}
	rows.Close()
	result["by_federation"] = byFed
	return result, nil
}

func (r *NSMISRepo) PerformanceDashboard(ctx context.Context) (map[string]interface{}, error) {
	result := map[string]interface{}{"as_of": time.Now().UTC(), "medals": map[string]int64{}, "by_federation": []map[string]interface{}{}, "by_country": []map[string]interface{}{}}
	medals := map[string]int64{"GOLD": 0, "SILVER": 0, "BRONZE": 0}
	rows, err := r.db.Query(ctx, `SELECT medal_type,COUNT(*) FROM medals WHERE status='APPROVED' GROUP BY medal_type`)
	if err != nil {
		return nil, err
	}
	var total int64
	for rows.Next() {
		var k string
		var n int64
		if err := rows.Scan(&k, &n); err != nil {
			rows.Close()
			return nil, err
		}
		medals[k] = n
		total += n
	}
	rows.Close()
	result["total_medals"] = total
	result["medals"] = medals
	byFed := make([]map[string]interface{}, 0)
	rows, err = r.db.Query(ctx, `SELECT f.id,f.name,COUNT(*) FILTER(WHERE m.medal_type='GOLD'),COUNT(*) FILTER(WHERE m.medal_type='SILVER'),COUNT(*) FILTER(WHERE m.medal_type='BRONZE') FROM federations f LEFT JOIN medals m ON m.federation_id=f.id AND m.status='APPROVED' GROUP BY f.id,f.name ORDER BY 3 DESC,4 DESC,5 DESC,f.name`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id, name string
		var g, s, b int64
		if err := rows.Scan(&id, &name, &g, &s, &b); err != nil {
			rows.Close()
			return nil, err
		}
		byFed = append(byFed, map[string]interface{}{"federation_id": id, "federation_name": name, "gold": g, "silver": s, "bronze": b})
	}
	rows.Close()
	result["by_federation"] = byFed
	byCountry := make([]map[string]interface{}, 0)
	rows, err = r.db.Query(ctx, `SELECT country,COUNT(*) FROM medals WHERE status='APPROVED' GROUP BY country ORDER BY 2 DESC,country LIMIT 50`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var country string
		var n int64
		if err := rows.Scan(&country, &n); err != nil {
			rows.Close()
			return nil, err
		}
		byCountry = append(byCountry, map[string]interface{}{"country": country, "count": n})
	}
	rows.Close()
	result["by_country"] = byCountry
	var international int64
	if err := r.db.QueryRow(ctx, `SELECT COALESCE(SUM(athlete_count),0) FROM competitions WHERE status='APPROVED' AND level IN ('CONTINENTAL','INTERNATIONAL')`).Scan(&international); err != nil {
		return nil, err
	}
	result["international_representation"] = international
	return result, nil
}

func (r *NSMISRepo) FinanceDashboard(ctx context.Context) (map[string]interface{}, error) {
	result := map[string]interface{}{"as_of": time.Now().UTC()}
	var released, outstanding float64
	var accountabilityCount int64
	if err := r.db.QueryRow(ctx, `SELECT COALESCE(SUM(amount) FILTER(WHERE status<>'CANCELLED'),0),COALESCE(SUM(amount) FILTER(WHERE status IN ('RELEASED','PARTIALLY_ACCOUNTED') AND accountability_due_on<CURRENT_DATE),0) FROM ncs_disbursements`).Scan(&released, &outstanding); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM financial_reports WHERE status IN ('APPROVED','LOCKED')`).Scan(&accountabilityCount); err != nil {
		return nil, err
	}
	result["currency"] = "UGX"
	result["funds_released"] = released
	result["accountability_submitted"] = accountabilityCount
	result["outstanding_accountabilities"] = outstanding
	return result, nil
}

func (r *NSMISRepo) TalentDashboard(ctx context.Context) (map[string]interface{}, error) {
	result := map[string]interface{}{"as_of": time.Now().UTC()}
	var identified, progressed, scholarships int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*),COUNT(*) FILTER(WHERE status='NATIONAL_TEAM') FROM talent_records`).Scan(&identified, &progressed); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM talent_scholarships WHERE status='ACTIVE' AND starts_on<=CURRENT_DATE AND (ends_on IS NULL OR ends_on>=CURRENT_DATE)`).Scan(&scholarships); err != nil {
		return nil, err
	}
	result["athletes_identified"] = identified
	result["progressed_to_national_teams"] = progressed
	result["receiving_scholarships"] = scholarships
	return result, nil
}
