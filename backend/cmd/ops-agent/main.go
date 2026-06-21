package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/database"
	appLogging "github.com/atenimedia-llc/ncs-online/backend/internal/logging"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type agent struct {
	repo                           *repository.BackupRepo
	databaseURL, dir, rcloneRemote string
}

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	closer, err := appLogging.Configure(cfg.LogLevel, cfg.LogDir)
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	a := &agent{repo: repository.New(db).Backups, databaseURL: cfg.DatabaseURL, dir: env("BACKUP_DIR", "/var/backups/ncs"), rcloneRemote: os.Getenv("RCLONE_REMOTE")}
	if err = os.MkdirAll(a.dir, 0o700); err != nil {
		panic(err)
	}
	worker := "ops-" + uuid.NewString()
	slog.Info("operations agent started", "worker", worker)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		_ = a.repo.EnsureDaily(ctx)
		job, claimErr := a.repo.Claim(ctx, worker)
		if errors.Is(claimErr, repository.ErrNotFound) {
			cancel()
			time.Sleep(3 * time.Second)
			continue
		}
		if claimErr != nil {
			slog.Error("backup job claim failed", "error", claimErr)
			cancel()
			time.Sleep(5 * time.Second)
			continue
		}
		runErr := a.run(ctx, job)
		if runErr != nil {
			safe := a.safeError(runErr)
			_ = a.repo.Fail(ctx, job.ID, safe)
			slog.Error("backup job failed", "job", job.ID, "type", job.JobType, "error", safe)
		}
		cancel()
	}
}
func (a *agent) run(ctx context.Context, j *repository.BackupJob) error {
	switch j.JobType {
	case "BACKUP":
		b, err := a.backup(ctx, j.ID, "scheduled")
		if err != nil {
			return err
		}
		result, _ := json.Marshal(b)
		return a.repo.Complete(ctx, j.ID, result)
	case "VERIFY":
		if j.BackupID == nil {
			return errors.New("missing backup id")
		}
		result, err := a.verify(ctx, j.ID, *j.BackupID)
		if err != nil {
			return err
		}
		return a.repo.Complete(ctx, j.ID, result)
	case "RESTORE":
		if j.BackupID == nil {
			return errors.New("missing backup id")
		}
		return a.restore(ctx, j, *j.BackupID)
	default:
		return errors.New("unsupported backup job")
	}
}
func (a *agent) backup(ctx context.Context, jobID, label string) (repository.BackupRecord, error) {
	_ = a.repo.Progress(ctx, jobID, 15, "Creating PostgreSQL custom-format dump")
	stamp := time.Now().UTC().Format("20060102-150405")
	name := fmt.Sprintf("ncs-%s-%s.dump", stamp, label)
	tmp, target := filepath.Join(a.dir, "."+name+".tmp"), filepath.Join(a.dir, name)
	defer os.Remove(tmp)
	if out, err := exec.CommandContext(ctx, "pg_dump", "--dbname", a.databaseURL, "--format=custom", "--compress=9", "--no-owner", "--no-privileges", "--file", tmp).CombinedOutput(); err != nil {
		return repository.BackupRecord{}, fmt.Errorf("pg_dump: %s", out)
	}
	if out, err := exec.CommandContext(ctx, "pg_restore", "--list", tmp).CombinedOutput(); err != nil {
		return repository.BackupRecord{}, fmt.Errorf("dump validation: %s", out)
	}
	if err := os.Rename(tmp, target); err != nil {
		return repository.BackupRecord{}, err
	}
	sum, err := fileChecksum(target)
	if err != nil {
		return repository.BackupRecord{}, err
	}
	if err = os.WriteFile(target+".sha256", []byte(sum+"  "+name+"\n"), 0o600); err != nil {
		return repository.BackupRecord{}, err
	}
	info, _ := os.Stat(target)
	record := repository.BackupRecord{ID: uuid.NewString(), FileName: name, SizeBytes: info.Size(), Checksum: sum, Status: "CREATED"}
	if err = a.repo.CreateRecord(ctx, record); err != nil {
		return repository.BackupRecord{}, err
	}
	if a.rcloneRemote != "" {
		_ = a.repo.Progress(ctx, jobID, 75, "Uploading backup artifacts off-site")
		remote := strings.TrimRight(a.rcloneRemote, "/")
		if out, copyErr := exec.CommandContext(ctx, "rclone", "copyto", target, remote+"/"+name).CombinedOutput(); copyErr != nil {
			return repository.BackupRecord{}, fmt.Errorf("off-site sync: %s", out)
		}
		if out, copyErr := exec.CommandContext(ctx, "rclone", "copyto", target+".sha256", remote+"/"+name+".sha256").CombinedOutput(); copyErr != nil {
			return repository.BackupRecord{}, fmt.Errorf("off-site manifest sync: %s", out)
		}
	}
	a.pruneLocal(7)
	return record, nil
}

func (a *agent) pruneLocal(keep int) {
	files, _ := filepath.Glob(filepath.Join(a.dir, "ncs-*.dump"))
	sort.Slice(files, func(i, j int) bool {
		left, _ := os.Stat(files[i])
		right, _ := os.Stat(files[j])
		return left.ModTime().After(right.ModTime())
	})
	if len(files) <= keep {
		return
	}
	for _, path := range files[keep:] {
		_ = os.Remove(path)
		_ = os.Remove(path + ".sha256")
	}
}
func (a *agent) verify(ctx context.Context, jobID, backupID string) (json.RawMessage, error) {
	record, err := a.repo.Get(ctx, backupID)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(a.dir, filepath.Base(record.FileName))
	sum, err := fileChecksum(path)
	if err != nil {
		return nil, err
	}
	if sum != record.Checksum {
		return nil, errors.New("backup checksum mismatch")
	}
	_ = a.repo.Progress(ctx, jobID, 30, "Restoring into isolated verification database")
	dbName := "ncs_verify_" + strings.ReplaceAll(uuid.NewString(), "-", "")[:12]
	verifyURL, err := databaseURLWithName(a.databaseURL, dbName)
	if err != nil {
		return nil, err
	}
	if out, err := exec.CommandContext(ctx, "createdb", "--maintenance-db", a.databaseURL, dbName).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("create verification database: %s", out)
	}
	defer exec.CommandContext(context.Background(), "dropdb", "--if-exists", "--maintenance-db", a.databaseURL, dbName).Run()
	if out, err := exec.CommandContext(ctx, "pg_restore", "--exit-on-error", "--no-owner", "--no-privileges", "--dbname", verifyURL, path).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("verification restore: %s", out)
	}
	query := `SELECT json_build_object('users',(SELECT count(*) FROM users),'applications',(SELECT count(*) FROM applications),'audit_logs',(SELECT count(*) FROM audit_logs),'organisations',(SELECT count(*) FROM organisations));`
	restored, err := exec.CommandContext(ctx, "psql", verifyURL, "-At", "-c", query).Output()
	if err != nil {
		return nil, fmt.Errorf("restored row counts: %w", err)
	}
	source, err := exec.CommandContext(ctx, "psql", a.databaseURL, "-At", "-c", query).Output()
	if err != nil {
		return nil, fmt.Errorf("source row counts: %w", err)
	}
	result, _ := json.Marshal(map[string]any{"checksum_valid": true, "restored_counts": json.RawMessage(strings.TrimSpace(string(restored))), "source_counts": json.RawMessage(strings.TrimSpace(string(source))), "verified_at": time.Now().UTC()})
	if err = a.repo.MarkVerified(ctx, backupID, result); err != nil {
		return nil, err
	}
	return result, nil
}
func (a *agent) restore(ctx context.Context, j *repository.BackupJob, backupID string) error {
	enabled, err := a.repo.MaintenanceEnabled(ctx)
	if err != nil || !enabled {
		return errors.New("maintenance mode is required")
	}
	record, err := a.repo.Get(ctx, backupID)
	if err != nil {
		return err
	}
	if record.Status != "VERIFIED" {
		return errors.New("backup is not verified")
	}
	if _, err = a.backup(ctx, j.ID, "pre-restore"); err != nil {
		return fmt.Errorf("safety snapshot: %w", err)
	}
	path := filepath.Join(a.dir, filepath.Base(record.FileName))
	sum, err := fileChecksum(path)
	if err != nil || sum != record.Checksum {
		return errors.New("restore checksum validation failed")
	}
	_ = a.repo.Progress(ctx, j.ID, 65, "Restoring verified backup")
	out, err := exec.CommandContext(ctx, "pg_restore", "--clean", "--if-exists", "--exit-on-error", "--no-owner", "--no-privileges", "--dbname", a.databaseURL, path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("restore: %s", out)
	}
	result, _ := json.Marshal(map[string]any{"restored_backup_id": backupID, "safety_snapshot_created": true, "finished_at": time.Now().UTC()})
	return a.repo.Complete(ctx, j.ID, result)
}
func fileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
func databaseURLWithName(raw, name string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	u.Path = "/" + name
	return u.String(), nil
}
func (a *agent) safeError(err error) string {
	s := strings.ReplaceAll(err.Error(), a.databaseURL, "[DATABASE_URL]")
	if len(s) > 1000 {
		s = s[:1000]
	}
	return s
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
