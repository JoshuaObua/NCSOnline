package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BackupRecord struct {
	ID                 string          `json:"id"`
	FileName           string          `json:"file_name"`
	SizeBytes          int64           `json:"size_bytes"`
	Checksum           string          `json:"checksum"`
	Status             string          `json:"status"`
	VerifiedAt         *time.Time      `json:"verified_at,omitempty"`
	VerificationResult json.RawMessage `json:"verification_result,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
}
type BackupJob struct {
	ID          string     `json:"id"`
	JobType     string     `json:"job_type"`
	BackupID    *string    `json:"backup_id,omitempty"`
	Status      string     `json:"status"`
	Progress    int        `json:"progress"`
	Message     string     `json:"message"`
	RequestedBy *string    `json:"requested_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
}
type BackupRepo struct{ db *pgxpool.Pool }

func (r *BackupRepo) List(ctx context.Context) ([]BackupRecord, error) {
	rows, err := r.db.Query(ctx, `SELECT id,file_name,size_bytes,checksum,status,verified_at,verification_result,created_at FROM system_backups ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []BackupRecord{}
	for rows.Next() {
		var b BackupRecord
		if err := rows.Scan(&b.ID, &b.FileName, &b.SizeBytes, &b.Checksum, &b.Status, &b.VerifiedAt, &b.VerificationResult, &b.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, b)
	}
	return items, rows.Err()
}
func (r *BackupRepo) ListJobs(ctx context.Context) ([]BackupJob, error) {
	rows, err := r.db.Query(ctx, `SELECT id,job_type,backup_id,status,progress,message,requested_by,created_at,started_at,finished_at FROM backup_jobs ORDER BY created_at DESC LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []BackupJob{}
	for rows.Next() {
		var j BackupJob
		if err := rows.Scan(&j.ID, &j.JobType, &j.BackupID, &j.Status, &j.Progress, &j.Message, &j.RequestedBy, &j.CreatedAt, &j.StartedAt, &j.FinishedAt); err != nil {
			return nil, err
		}
		items = append(items, j)
	}
	return items, rows.Err()
}

func (r *BackupRepo) DeleteJob(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM backup_jobs WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *BackupRepo) ClearJobs(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `DELETE FROM backup_jobs WHERE status NOT IN ('PENDING','RUNNING')`)
	return err
}
func (r *BackupRepo) Queue(ctx context.Context, kind string, backupID *string, actor string) (BackupJob, error) {
	j := BackupJob{ID: uuid.NewString(), JobType: kind, BackupID: backupID, Status: "PENDING"}
	if actor != "" {
		j.RequestedBy = &actor
	}
	err := r.db.QueryRow(ctx, `INSERT INTO backup_jobs(id,job_type,backup_id,requested_by) VALUES($1,$2,$3,$4) RETURNING status,progress,message,created_at`, j.ID, kind, backupID, nullableStr(actor)).Scan(&j.Status, &j.Progress, &j.Message, &j.CreatedAt)
	return j, err
}
func (r *BackupRepo) EnsureDaily(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `INSERT INTO backup_jobs(id,job_type,status,message) SELECT gen_random_uuid()::TEXT,'BACKUP','PENDING','Scheduled daily backup' WHERE NOT EXISTS(SELECT 1 FROM backup_jobs WHERE job_type='BACKUP' AND status IN('PENDING','RUNNING') OR (job_type='BACKUP' AND status='SUCCEEDED' AND finished_at>NOW()-INTERVAL '23 hours'))`)
	return err
}
func (r *BackupRepo) Claim(ctx context.Context, worker string) (*BackupJob, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var j BackupJob
	err = tx.QueryRow(ctx, `SELECT id,job_type,backup_id,status,progress,message,requested_by,created_at FROM backup_jobs WHERE status='PENDING' ORDER BY created_at FOR UPDATE SKIP LOCKED LIMIT 1`).Scan(&j.ID, &j.JobType, &j.BackupID, &j.Status, &j.Progress, &j.Message, &j.RequestedBy, &j.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(ctx, `UPDATE backup_jobs SET status='RUNNING',worker_id=$2,started_at=NOW(),progress=5,message='Job claimed' WHERE id=$1`, j.ID, worker)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	j.Status = "RUNNING"
	j.Progress = 5
	return &j, nil
}
func (r *BackupRepo) Progress(ctx context.Context, id string, progress int, message string) error {
	_, err := r.db.Exec(ctx, `UPDATE backup_jobs SET progress=$2,message=$3 WHERE id=$1`, id, progress, message)
	return err
}
func (r *BackupRepo) Complete(ctx context.Context, id string, result json.RawMessage) error {
	_, err := r.db.Exec(ctx, `UPDATE backup_jobs SET status='SUCCEEDED',progress=100,message='Completed',result=$2,finished_at=NOW() WHERE id=$1`, id, result)
	return err
}
func (r *BackupRepo) Fail(ctx context.Context, id, message string) error {
	_, err := r.db.Exec(ctx, `UPDATE backup_jobs SET status='FAILED',message=$2,finished_at=NOW() WHERE id=$1`, id, message)
	return err
}
func (r *BackupRepo) CreateRecord(ctx context.Context, b BackupRecord) error {
	_, err := r.db.Exec(ctx, `INSERT INTO system_backups(id,file_name,size_bytes,checksum,status) VALUES($1,$2,$3,$4,$5)`, b.ID, b.FileName, b.SizeBytes, b.Checksum, b.Status)
	return err
}
func (r *BackupRepo) Get(ctx context.Context, id string) (BackupRecord, error) {
	var b BackupRecord
	err := r.db.QueryRow(ctx, `SELECT id,file_name,size_bytes,checksum,status,verified_at,verification_result,created_at FROM system_backups WHERE id=$1`, id).Scan(&b.ID, &b.FileName, &b.SizeBytes, &b.Checksum, &b.Status, &b.VerifiedAt, &b.VerificationResult, &b.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return b, ErrNotFound
	}
	return b, err
}
func (r *BackupRepo) MarkVerified(ctx context.Context, id string, result json.RawMessage) error {
	_, err := r.db.Exec(ctx, `UPDATE system_backups SET status='VERIFIED',verified_at=NOW(),verification_result=$2 WHERE id=$1`, id, result)
	return err
}
func (r *BackupRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM system_backups WHERE id=$1`)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
func (r *BackupRepo) MaintenanceEnabled(ctx context.Context) (bool, error) {
	var enabled bool
	err := r.db.QueryRow(ctx, `
		SELECT maintenance_mode
		    OR COALESCE((public_cms_maintenance->>'maintenance_mode')::boolean, false)
		    OR COALESCE((admin_dashboard_maintenance->>'maintenance_mode')::boolean, false)
		  FROM system_control
		 WHERE singleton=TRUE`).Scan(&enabled)
	return enabled, err
}
