package main

import (
	"context"
	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/database"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/services"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repos := repository.New(db)
	repo := repos.NSMIS
	emailer := services.NewEmailService(cfg)
	backupExecutor := services.NewBackupExecutor(repos.Backups, cfg.DatabaseURL)
	iotecSvc := services.NewIoTecService()
	formSvc := services.NewFormService(repos.Forms, repos.Departments, repos.Audit, repos.Transactions, iotecSvc)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	workerID := "nsmis-" + uuid.NewString()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	log.Printf("NSMIS worker %s started", workerID)
	for {
		select {
		case <-ctx.Done():
			log.Print("NSMIS worker stopped")
			return
		case <-ticker.C:
			// Sync pending ioTec mobile money payments
			if count, err := formSvc.SyncPendingTransactions(ctx); err != nil {
				log.Printf("error syncing pending iotec transactions: %v", err)
			} else if count > 0 {
				log.Printf("successfully synced %d pending iotec transactions", count)
			}

			for i := 0; i < 10; i++ {
				n, claimErr := repos.Notifications.ClaimEmail(ctx)
				if claimErr == repository.ErrNotFound {
					break
				}
				if claimErr != nil {
					log.Printf("notification claim failed: %v", claimErr)
					break
				}
				if sendErr := emailer.Send(n); sendErr != nil {
					_ = repos.Notifications.Failed(ctx, n.ID, sendErr.Error(), n.AttemptCount)
					log.Printf("notification %s failed: %v", n.ID, sendErr)
				} else {
					_ = repos.Notifications.Delivered(ctx, n.ID)
				}
			}
			for i := 0; i < 10; i++ {
				worked, e := repo.RunNextJob(ctx, workerID)
				if e != nil {
					log.Printf("job failed: %v", e)
				}
				if !worked {
					break
				}
			}
			if err := repos.Backups.EnsureDaily(ctx); err != nil {
				log.Printf("could not schedule daily backup: %v", err)
			}
			for i := 0; i < 5; i++ {
				job, claimErr := repos.Backups.Claim(ctx, workerID)
				if claimErr == repository.ErrNotFound {
					break
				}
				if claimErr != nil {
					log.Printf("backup job claim failed: %v", claimErr)
					break
				}
				if runErr := backupExecutor.RunJob(ctx, *job); runErr != nil {
					log.Printf("backup job %s (%s) failed: %v", job.ID, job.JobType, runErr)
					_ = repos.Backups.Fail(ctx, job.ID, runErr.Error())
				}
			}
		}
	}
}
