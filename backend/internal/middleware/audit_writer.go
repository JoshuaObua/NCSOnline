package middleware

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
)

type auditJob struct {
	entry  *models.AuditLog
	userID string
}

// AuditWriter bounds memory and concurrency. Request handlers only enqueue;
// one worker enriches and batches events, preserving chain order.
type AuditWriter struct {
	repo    *repository.AuditRepo
	jobs    chan auditJob
	done    chan struct{}
	wg      sync.WaitGroup
	dropped atomic.Uint64
}

func NewAuditWriter(repo *repository.AuditRepo, capacity int) *AuditWriter {
	if capacity < 100 {
		capacity = 2048
	}
	w := &AuditWriter{repo: repo, jobs: make(chan auditJob, capacity), done: make(chan struct{})}
	w.wg.Add(1)
	go w.run()
	return w
}

func (w *AuditWriter) Enqueue(entry *models.AuditLog, userID string) bool {
	select {
	case w.jobs <- auditJob{entry: entry, userID: userID}:
		return true
	default:
		w.dropped.Add(1)
		slog.Error("audit queue full", "dropped_total", w.dropped.Load())
		return false
	}
}

func (w *AuditWriter) Close(ctx context.Context) error {
	close(w.done)
	finished := make(chan struct{})
	go func() { w.wg.Wait(); close(finished) }()
	select {
	case <-finished:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *AuditWriter) run() {
	defer w.wg.Done()
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	batch := make([]auditJob, 0, 128)
	flush := func() {
		if len(batch) == 0 {
			return
		}
		entries := make([]*models.AuditLog, 0, len(batch))
		for _, job := range batch {
			w.enrich(job)
			entries = append(entries, job.entry)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		if err := w.repo.LogBatch(ctx, entries); err != nil {
			slog.Error("audit batch write failed", "count", len(entries), "error", err)
		}
		cancel()
		batch = batch[:0]
	}
	for {
		select {
		case job := <-w.jobs:
			batch = append(batch, job)
			if len(batch) >= 128 {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-w.done:
			for {
				select {
				case job := <-w.jobs:
					batch = append(batch, job)
				default:
					flush()
					return
				}
			}
		}
	}
}

func (w *AuditWriter) enrich(job auditJob) {
	e := job.entry
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	if locationService != nil {
		if location, err := locationService.Locate(ctx, e.IPAddress, e.UserAgent, job.userID); err == nil {
			e.GeoCountry = location.Country
			e.GeoCity = location.City
			e.GeoRegion = location.Region
			e.GeoLatitude = location.Latitude
			e.GeoLongitude = location.Longitude
			e.GeoTimezone = location.Timezone
			e.GeoSource = location.Source
			e.Platform = location.Platform
			e.VPNDetected = location.IsProxy || location.IsVPN || location.IsTor || location.IsHosting
			e.Authenticated = job.userID != ""
			if location.Browser != "" {
				e.Browser = location.Browser
			}
			if location.DeviceType != "" {
				e.ClientType = location.DeviceType
			}
		}
	}
	if e.GeoCountry == "" {
		e.GeoCountry = "Unknown"
	}
	e.ThreatScore = computeThreatScore(e.GeoCountry, e.VPNDetected, e.ResponseCode)
	e.SeverityLevel = classifySeverity(e.EventStatus, e.ThreatScore, e.VPNDetected)
	e.AnomalyDetected = e.AnomalyDetected || e.ThreatScore >= 60
}
