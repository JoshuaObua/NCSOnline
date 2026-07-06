package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type BackgroundJob struct {
	ID, JobType, Status, IdempotencyKey, LastError string
	Payload, Result                                map[string]interface{}
	AttemptCount, MaxAttempts                      int
	CreatedAt                                      time.Time
}

func (r *NSMISRepo) EnqueueJob(ctx context.Context, jobType, key string, payload map[string]interface{}, actor string) (string, error) {
	if key == "" {
		key = jobType + ":" + uuid.NewString()
	}
	raw, _ := json.Marshal(payload)
	var id string
	err := r.db.QueryRow(ctx, `INSERT INTO background_jobs(job_type,idempotency_key,payload,created_by) VALUES($1,$2,$3,$4) ON CONFLICT(idempotency_key) DO UPDATE SET idempotency_key=EXCLUDED.idempotency_key RETURNING id`, jobType, key, raw, nullableStr(actor)).Scan(&id)
	return id, err
}
func (r *NSMISRepo) ListJobs(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	if limit < 1 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `SELECT to_jsonb(j) FROM background_jobs j ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]interface{}{}
	for rows.Next() {
		var b []byte
		if err = rows.Scan(&b); err != nil {
			return nil, err
		}
		var x map[string]interface{}
		if err = json.Unmarshal(b, &x); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *NSMISRepo) RetryJob(ctx context.Context, id string) error {
	res, err := r.db.Exec(ctx, `UPDATE background_jobs SET status='PENDING',next_attempt_at=NOW(),leased_until=NULL,leased_by='',last_error='',updated_at=NOW() WHERE id=$1 AND status IN('FAILED','DEAD_LETTER')`, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *NSMISRepo) RunNextJob(ctx context.Context, workerID string) (bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)
	var id, typ string
	var raw []byte
	var attempts, max int
	err = tx.QueryRow(ctx, `SELECT id,job_type,payload,attempt_count,max_attempts FROM background_jobs WHERE status IN('PENDING','FAILED') AND next_attempt_at<=NOW() AND (leased_until IS NULL OR leased_until<NOW()) ORDER BY created_at FOR UPDATE SKIP LOCKED LIMIT 1`).Scan(&id, &typ, &raw, &attempts, &max)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	attempts++
	_, err = tx.Exec(ctx, `UPDATE background_jobs SET status='RUNNING',attempt_count=$2,leased_by=$3,leased_until=NOW()+INTERVAL '5 minutes',heartbeat_at=NOW(),started_at=COALESCE(started_at,NOW()),updated_at=NOW() WHERE id=$1`, id, attempts, workerID)
	if err != nil {
		return false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return false, err
	}
	var payload map[string]interface{}
	_ = json.Unmarshal(raw, &payload)
	result, runErr := r.executeJob(ctx, typ, payload)
	resultRaw, _ := json.Marshal(result)
	if runErr == nil {
		_, err = r.db.Exec(ctx, `UPDATE background_jobs SET status='SUCCEEDED',result=$2,finished_at=NOW(),leased_until=NULL,updated_at=NOW() WHERE id=$1`, id, resultRaw)
		return true, err
	}
	status := "FAILED"
	if attempts >= max {
		status = "DEAD_LETTER"
	}
	delay := time.Duration(attempts*attempts) * time.Minute
	_, err = r.db.Exec(ctx, `UPDATE background_jobs SET status=$2,last_error=$3,next_attempt_at=$4,leased_until=NULL,updated_at=NOW() WHERE id=$1`, id, status, runErr.Error(), time.Now().Add(delay))
	if err != nil {
		return true, err
	}
	return true, runErr
}
func (r *NSMISRepo) executeJob(ctx context.Context, typ string, p map[string]interface{}) (map[string]interface{}, error) {
	switch typ {
	case "COMPLIANCE_REFRESH":
		n, err := r.calculateCompliance(ctx, fmt.Sprint(p["period_id"]))
		return map[string]interface{}{"scores_updated": n}, err
	case "DEADLINE_REMINDERS":
		n, err := r.queueDeadlineReminders(ctx)
		return map[string]interface{}{"notifications_queued": n}, err
	case "CREDENTIAL_EXPIRY":
		n, err := r.queueCredentialExpiry(ctx)
		return map[string]interface{}{"notifications_queued": n}, err
	case "DATA_QUALITY_SCAN":
		n, err := r.scanDataQuality(ctx)
		return map[string]interface{}{"issues_opened": n}, err
	default:
		return nil, fmt.Errorf("unsupported job type %s", typ)
	}
}
func (r *NSMISRepo) calculateCompliance(ctx context.Context, period string) (int64, error) {
	if period == "" {
		return 0, fmt.Errorf("period_id is required")
	}
	res, err := r.db.Exec(ctx, `INSERT INTO compliance_scores(federation_id,reporting_period_id,rule_version_id,score,is_compliant,breakdown)
SELECT f.id,$1,rv.id,s.score,s.score>=rv.compliant_threshold,jsonb_build_object('timeliness',x.timeliness,'governance',x.governance,'document_validity',x.documents,'data_quality',x.quality)
FROM federations f CROSS JOIN LATERAL(SELECT
 COALESCE(40.0*COUNT(*) FILTER(WHERE o.status IN('SUBMITTED','ACCEPTED'))/NULLIF(COUNT(*),0),0) timeliness,
 CASE WHEN EXISTS(SELECT 1 FROM federation_reports r WHERE r.federation_id=f.id AND r.reporting_period_id=$1 AND r.report_type='GOVERNANCE' AND r.status IN('APPROVED','LOCKED')) THEN 30 ELSE 0 END governance,
 CASE WHEN EXISTS(SELECT 1 FROM federation_documents d WHERE d.federation_id=f.id AND d.document_type='CONSTITUTION' AND d.scan_status='CLEAN' AND (d.expires_on IS NULL OR d.expires_on>=CURRENT_DATE)) THEN 20 ELSE 0 END documents,
 CASE WHEN NOT EXISTS(SELECT 1 FROM data_quality_issues q WHERE q.federation_id=f.id AND q.status IN('OPEN','ASSIGNED') AND q.severity IN('ERROR','BLOCKING')) THEN 10 ELSE 0 END quality
 FROM report_obligations o WHERE o.federation_id=f.id AND o.reporting_period_id=$1)x
CROSS JOIN LATERAL(SELECT x.timeliness+x.governance+x.documents+x.quality score)s JOIN compliance_rule_versions rv ON rv.is_active
WHERE f.is_active AND f.deleted_at IS NULL ON CONFLICT(federation_id,reporting_period_id,rule_version_id) DO UPDATE SET score=EXCLUDED.score,is_compliant=EXCLUDED.is_compliant,breakdown=EXCLUDED.breakdown,calculated_at=NOW()`, period)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
func (r *NSMISRepo) queueDeadlineReminders(ctx context.Context) (int64, error) {
	res, err := r.db.Exec(ctx, `INSERT INTO notification_deliveries(template_code,recipient_user_id,recipient_address,channel,related_type,related_id,payload)
SELECT 'REPORT_DEADLINE',fm.user_id,u.email,'EMAIL','report_obligation',o.id,jsonb_build_object('federation',f.name,'report_type',o.report_type,'due_on',o.due_on,'days_until',o.due_on-CURRENT_DATE)
FROM report_obligations o JOIN federations f ON f.id=o.federation_id JOIN federation_memberships fm ON fm.federation_id=f.id AND fm.is_active JOIN users u ON u.id=fm.user_id AND u.is_active
WHERE o.status NOT IN('ACCEPTED','EXEMPT') AND (o.due_on-CURRENT_DATE)=ANY(ARRAY[14,7,3,1,0,-1,-3,-7,-14]) AND NOT EXISTS(SELECT 1 FROM notification_deliveries n WHERE n.template_code='REPORT_DEADLINE' AND n.related_id=o.id AND n.recipient_user_id=fm.user_id AND n.created_at::date=CURRENT_DATE)`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
func (r *NSMISRepo) queueCredentialExpiry(ctx context.Context) (int64, error) {
	res, err := r.db.Exec(ctx, `WITH exp AS(SELECT federation_id,'COACH' kind,id,full_name,expiry_date expires FROM coaches WHERE status='ACTIVE' AND expiry_date-CURRENT_DATE=ANY(ARRAY[90,30,14,7,1,0,-1]) UNION ALL SELECT federation_id,'OFFICIAL',id,full_name,valid_until FROM technical_officials WHERE status='ACTIVE' AND valid_until-CURRENT_DATE=ANY(ARRAY[90,30,14,7,1,0,-1])) INSERT INTO notification_deliveries(template_code,recipient_user_id,recipient_address,channel,related_type,related_id,payload) SELECT 'CREDENTIAL_EXPIRY',fm.user_id,u.email,'EMAIL',e.kind,e.id,jsonb_build_object('name',e.full_name,'expires_on',e.expires) FROM exp e JOIN federation_memberships fm ON fm.federation_id=e.federation_id AND fm.is_active JOIN users u ON u.id=fm.user_id WHERE NOT EXISTS(SELECT 1 FROM notification_deliveries n WHERE n.template_code='CREDENTIAL_EXPIRY' AND n.related_id=e.id AND n.recipient_user_id=fm.user_id AND n.created_at::date=CURRENT_DATE)`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
func (r *NSMISRepo) scanDataQuality(ctx context.Context) (int64, error) {
	res, err := r.db.Exec(ctx, `INSERT INTO data_quality_issues(federation_id,record_type,record_id,rule_code,severity,description)
SELECT aa.federation_id,'athlete',a.id,'ATHLETE_POSSIBLE_DUPLICATE','WARNING','Possible duplicate athlete name and date of birth' FROM athletes a JOIN athlete_affiliations aa ON aa.athlete_id=a.id AND aa.is_active WHERE EXISTS(SELECT 1 FROM athletes x WHERE x.id<>a.id AND lower(x.full_name)=lower(a.full_name) AND x.date_of_birth=a.date_of_birth AND x.deleted_at IS NULL) ON CONFLICT DO NOTHING`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
