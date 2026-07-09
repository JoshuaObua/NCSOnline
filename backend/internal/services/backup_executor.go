package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
)

// BackupExecutor performs the actual filesystem/database work behind a queued
// backup_jobs row (pg_dump/pg_restore), which is deliberately kept out of the
// HTTP handlers so it only ever runs on the worker process.
type BackupExecutor struct {
	repo  *repository.BackupRepo
	dbURL string
}

func NewBackupExecutor(repo *repository.BackupRepo, dbURL string) *BackupExecutor {
	return &BackupExecutor{repo: repo, dbURL: dbURL}
}

func (b *BackupExecutor) backupDir() string {
	dir := os.Getenv("BACKUP_DIR")
	if dir == "" {
		dir = "/var/backups/ncs"
	}
	return dir
}

// RunJob dispatches a claimed job to its handler. On success it marks the job
// complete itself (since BACKUP/VERIFY produce different result payloads);
// callers should call BackupRepo.Fail with the returned error on failure.
func (b *BackupExecutor) RunJob(ctx context.Context, job repository.BackupJob) error {
	switch job.JobType {
	case "BACKUP":
		return b.runBackup(ctx, job)
	case "VERIFY":
		return b.runVerify(ctx, job)
	case "RESTORE":
		return b.runRestore(ctx, job)
	default:
		return fmt.Errorf("unsupported backup job type %s", job.JobType)
	}
}

func (b *BackupExecutor) runBackup(ctx context.Context, job repository.BackupJob) error {
	dir := b.backupDir()
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("could not create backup directory: %w", err)
	}
	_ = b.repo.Progress(ctx, job.ID, 15, "Running pg_dump")

	fileName := fmt.Sprintf("ncs-backup-%s-%s.dump", time.Now().UTC().Format("20060102-150405"), job.ID[:8])
	path := filepath.Join(dir, fileName)

	cmd := exec.CommandContext(ctx, "pg_dump", "--dbname", b.dbURL, "--format=custom", "--no-owner", "--no-privileges", "--file", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("pg_dump failed: %s", safeTail(out))
	}
	_ = b.repo.Progress(ctx, job.ID, 70, "Computing checksum")

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("backup file missing after pg_dump: %w", err)
	}
	checksum, err := sha256File(path)
	if err != nil {
		return fmt.Errorf("could not checksum backup: %w", err)
	}
	if err := os.WriteFile(path+".sha256", []byte(checksum), 0o640); err != nil {
		return fmt.Errorf("could not write checksum file: %w", err)
	}

	record := repository.BackupRecord{
		ID:        uuid.NewString(),
		FileName:  fileName,
		SizeBytes: info.Size(),
		Checksum:  checksum,
		Status:    "CREATED",
	}
	if err := b.repo.CreateRecord(ctx, record); err != nil {
		return fmt.Errorf("could not save backup record: %w", err)
	}

	result, _ := json.Marshal(map[string]any{"backup_id": record.ID, "file_name": fileName, "size_bytes": info.Size()})
	return b.repo.Complete(ctx, job.ID, result)
}

func (b *BackupExecutor) runVerify(ctx context.Context, job repository.BackupJob) error {
	if job.BackupID == nil || *job.BackupID == "" {
		return fmt.Errorf("backup_id is required to verify a backup")
	}
	record, err := b.repo.Get(ctx, *job.BackupID)
	if err != nil {
		return fmt.Errorf("backup record not found: %w", err)
	}
	path := filepath.Join(b.backupDir(), filepath.Base(record.FileName))
	_ = b.repo.Progress(ctx, job.ID, 25, "Checking backup integrity")

	if checksum, err := sha256File(path); err != nil {
		return fmt.Errorf("could not read backup file: %w", err)
	} else if record.Checksum != "" && checksum != record.Checksum {
		return fmt.Errorf("checksum mismatch — backup file may be corrupted or tampered with")
	}

	out, err := exec.CommandContext(ctx, "pg_restore", "--list", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_restore --list failed: %s", safeTail(out))
	}
	_ = b.repo.Progress(ctx, job.ID, 80, "Marking backup verified")

	verification, _ := json.Marshal(map[string]any{"verified_at": time.Now().UTC(), "entries": len(strings.Split(strings.TrimSpace(string(out)), "\n"))})
	if err := b.repo.MarkVerified(ctx, record.ID, verification); err != nil {
		return fmt.Errorf("could not mark backup verified: %w", err)
	}

	result, _ := json.Marshal(map[string]any{"backup_id": record.ID, "status": "VERIFIED"})
	return b.repo.Complete(ctx, job.ID, result)
}

func (b *BackupExecutor) runRestore(ctx context.Context, job repository.BackupJob) error {
	if job.BackupID == nil || *job.BackupID == "" {
		return fmt.Errorf("backup_id is required to restore a backup")
	}
	record, err := b.repo.Get(ctx, *job.BackupID)
	if err != nil {
		return fmt.Errorf("backup record not found: %w", err)
	}
	if record.Status != "VERIFIED" {
		return fmt.Errorf("only a verified backup can be restored")
	}
	path := filepath.Join(b.backupDir(), filepath.Base(record.FileName))
	_ = b.repo.Progress(ctx, job.ID, 20, "Restoring database from backup")

	out, err := exec.CommandContext(ctx, "pg_restore", "--dbname", b.dbURL, "--clean", "--if-exists", "--no-owner", "--no-privileges", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_restore failed: %s", safeTail(out))
	}
	_ = b.repo.Progress(ctx, job.ID, 95, "Restore complete")

	result, _ := json.Marshal(map[string]any{"backup_id": record.ID, "restored_at": time.Now().UTC()})
	return b.repo.Complete(ctx, job.ID, result)
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func safeTail(out []byte) string {
	s := strings.TrimSpace(string(out))
	if s == "" {
		return "command failed with no output"
	}
	if len(s) > 800 {
		s = s[len(s)-800:]
	}
	return s
}
