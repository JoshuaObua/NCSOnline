// Package services — smart-update sentinel.
//
// Polls the project's GitHub Releases on a 5-minute ticker and exposes the
// latest known release alongside pre-flight readiness signals (disk + DB)
// so the operator UI can decide whether a deploy is safe to attempt.
//
// What this file IS:
//   * A bounded, cancellable in-memory cache of the latest GitHub release.
//   * Pre-flight aggregation: current version (from VERSION file), latest
//     release, disk free / total, DB ping.
//
// What this file is NOT yet (TODO when infrastructure is in place):
//   * Cosign signature verification of the release tarball.
//   * Polling the CI conclusion of the release SHA before allowing deploy.
//   * The actual `docker compose pull && up -d` execution path — that lives
//     in the handler and is gated on having the Docker socket mounted.

package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultRepo         = "atenimedia-llc/ncs-online"
	defaultVersionPath  = "/app/VERSION"
	pollInterval        = 5 * time.Minute
	httpTimeout         = 10 * time.Second
	githubReleasesURL   = "https://api.github.com/repos/%s/releases/latest"
	minRecommendedFreeB = 1 << 30 // 1 GiB — refuse pre-flight under this.
)

type GithubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	HTMLURL     string    `json:"html_url"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
}

type DiskInfo struct {
	FreeBytes  uint64 `json:"free_bytes"`
	TotalBytes uint64 `json:"total_bytes"`
	UsedPct    int    `json:"used_pct"`
	Healthy    bool   `json:"healthy"`
}

type Preflight struct {
	CurrentVersion string         `json:"current_version"`
	Latest         *GithubRelease `json:"latest,omitempty"`
	UpdateAvailable bool          `json:"update_available"`
	Disk           DiskInfo       `json:"disk"`
	DBReachable    bool           `json:"db_reachable"`
	DBLatencyMs    int64          `json:"db_latency_ms"`
	LastCheckedAt  *time.Time     `json:"last_checked_at,omitempty"`
	LastError      string         `json:"last_error,omitempty"`
	RepoSlug       string         `json:"repo_slug"`
}

type UpdatesService struct {
	repoSlug    string
	versionPath string
	db          *pgxpool.Pool
	http        *http.Client

	mu        sync.RWMutex
	latest    *GithubRelease
	lastCheck time.Time
	lastErr   string

	stop chan struct{}
	once sync.Once
}

func NewUpdatesService(db *pgxpool.Pool) *UpdatesService {
	repo := strings.TrimSpace(os.Getenv("GITHUB_REPO_SLUG"))
	if repo == "" {
		repo = defaultRepo
	}
	versionPath := strings.TrimSpace(os.Getenv("VERSION_FILE"))
	if versionPath == "" {
		versionPath = defaultVersionPath
	}
	return &UpdatesService{
		repoSlug:    repo,
		versionPath: versionPath,
		db:          db,
		http:        &http.Client{Timeout: httpTimeout},
		stop:        make(chan struct{}),
	}
}

// StartSentinel runs the background poller. Safe to call multiple times;
// only the first invocation actually launches the goroutine.
func (s *UpdatesService) StartSentinel(ctx context.Context) {
	s.once.Do(func() {
		go s.loop(ctx)
	})
}

func (s *UpdatesService) Stop() {
	select {
	case <-s.stop:
	default:
		close(s.stop)
	}
}

func (s *UpdatesService) loop(ctx context.Context) {
	// First poll immediately so the UI doesn't show "never checked" on cold start.
	_ = s.refresh(ctx)
	t := time.NewTicker(pollInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stop:
			return
		case <-t.C:
			_ = s.refresh(ctx)
		}
	}
}

// ForceRefresh triggers an immediate poll and returns the resulting preflight.
func (s *UpdatesService) ForceRefresh(ctx context.Context) (*Preflight, error) {
	_ = s.refresh(ctx)
	return s.Preflight(ctx)
}

func (s *UpdatesService) refresh(ctx context.Context) error {
	rel, err := s.fetchLatestRelease(ctx)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastCheck = time.Now()
	if err != nil {
		s.lastErr = err.Error()
		return err
	}
	s.lastErr = ""
	s.latest = rel
	return nil
}

func (s *UpdatesService) fetchLatestRelease(ctx context.Context) (*GithubRelease, error) {
	url := fmt.Sprintf(githubReleasesURL, s.repoSlug)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if tok := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode == http.StatusNotFound {
		// No releases published yet — treat as "no update", not an error.
		return nil, nil
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("github releases: %s — %s", resp.Status, strings.TrimSpace(string(body)))
	}
	rel := &GithubRelease{}
	if err := json.Unmarshal(body, rel); err != nil {
		return nil, fmt.Errorf("parse release: %w", err)
	}
	if rel.Draft || rel.Prerelease {
		// Skip drafts / prereleases for the "is update available" signal.
		return nil, nil
	}
	return rel, nil
}

// Preflight aggregates the cached release + live system signals.
func (s *UpdatesService) Preflight(ctx context.Context) (*Preflight, error) {
	current := strings.TrimSpace(s.readVersionFile())

	s.mu.RLock()
	latest := s.latest
	lastErr := s.lastErr
	last := s.lastCheck
	s.mu.RUnlock()

	pf := &Preflight{
		CurrentVersion: current,
		Latest:         latest,
		Disk:           diskInfo(),
		RepoSlug:       s.repoSlug,
		LastError:      lastErr,
	}
	if !last.IsZero() {
		pf.LastCheckedAt = &last
	}
	if latest != nil && !semverLEQ(latest.TagName, current) {
		pf.UpdateAvailable = true
	}

	// DB ping
	start := time.Now()
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := s.db.Ping(pingCtx); err == nil {
		pf.DBReachable = true
		pf.DBLatencyMs = time.Since(start).Milliseconds()
	}
	return pf, nil
}

func (s *UpdatesService) readVersionFile() string {
	b, err := os.ReadFile(s.versionPath)
	if err != nil {
		return "v0.0.0-dev"
	}
	return strings.TrimSpace(string(b))
}

// semverLEQ reports whether a <= b for tags formatted as vMAJOR.MINOR.PATCH.
// Falls back to string compare if either side isn't recognisable.
func semverLEQ(a, b string) bool {
	pa, oka := parseSemver(a)
	pb, okb := parseSemver(b)
	if !oka || !okb {
		return strings.TrimSpace(a) == strings.TrimSpace(b)
	}
	for i := 0; i < 3; i++ {
		if pa[i] != pb[i] {
			return pa[i] < pb[i]
		}
	}
	return true // equal
}

func parseSemver(v string) ([3]int, bool) {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return [3]int{}, false
	}
	out := [3]int{}
	for i, p := range parts {
		// strip pre-release suffix on patch.
		if i == 2 {
			if idx := strings.IndexAny(p, "-+"); idx > 0 {
				p = p[:idx]
			}
		}
		n := 0
		for _, c := range p {
			if c < '0' || c > '9' {
				return [3]int{}, false
			}
			n = n*10 + int(c-'0')
		}
		out[i] = n
	}
	return out, true
}

// Errors returned by deploy execution. Kept here so the handler can map them.
var (
	ErrDockerSocketMissing = errors.New("docker socket not available — mount /var/run/docker.sock into the backend container to enable in-place deploy")
	ErrDeployInProgress    = errors.New("a deploy is already in progress")
)
