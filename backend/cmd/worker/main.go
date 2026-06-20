package main

import (
	"context"
	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/database"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
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
	repo := repository.New(db).NSMIS
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
			for i := 0; i < 10; i++ {
				worked, e := repo.RunNextJob(ctx, workerID)
				if e != nil {
					log.Printf("job failed: %v", e)
				}
				if !worked {
					break
				}
			}
		}
	}
}
