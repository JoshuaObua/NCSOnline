// Deploy executor — runs `docker compose pull && up -d` against the host
// Docker daemon. Only works when the backend container has the Docker
// socket mounted (see docker-compose.yml).
//
// This is intentionally simple: serialised (one deploy at a time), with
// stdout/stderr captured into a ring buffer the UI polls. A proper
// WebSocket stream is a follow-up (Phase 5 console).

package services

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type DeployStatus struct {
	State     string    `json:"state"` // idle | running | success | failed
	StartedAt time.Time `json:"started_at,omitempty"`
	EndedAt   time.Time `json:"ended_at,omitempty"`
	ExitCode  int       `json:"exit_code"`
	Lines     []string  `json:"lines"` // most-recent-last, capped
	Error     string    `json:"error,omitempty"`
}

type Deployer struct {
	composeProject string

	mu   sync.Mutex
	curr DeployStatus
}

type DeployOptions struct {
	ComposeProject string
	Services       []string
	Script         string
}

func NewDeployer() *Deployer {
	proj := strings.TrimSpace(os.Getenv("COMPOSE_PROJECT_NAME"))
	if proj == "" {
		proj = "ncs-online"
	}
	return &Deployer{composeProject: proj, curr: DeployStatus{State: "idle"}}
}

// Status returns a snapshot of the current/last deploy.
func (d *Deployer) Status() DeployStatus {
	d.mu.Lock()
	defer d.mu.Unlock()
	// Copy the slice so callers can't mutate ours.
	out := d.curr
	out.Lines = append([]string(nil), d.curr.Lines...)
	return out
}

// dockerAvailable reports whether the docker CLI + socket are reachable.
func (d *Deployer) dockerAvailable(ctx context.Context) bool {
	if _, err := os.Stat("/var/run/docker.sock"); err != nil {
		return false
	}
	if _, err := exec.LookPath("docker"); err != nil {
		return false
	}
	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return exec.CommandContext(cctx, "docker", "version", "--format", "{{.Server.Version}}").Run() == nil
}

// Start begins an asynchronous deploy. Returns immediately; poll Status
// for progress. ErrDeployInProgress if one is already running.
//
// Targets: backend, frontend, worker, backup — i.e. the services that
// share the project's source. Skips postgres / nginx / pgadmin to avoid
// disrupting persistent state and edge routing during the swap.
func (d *Deployer) Start(parent context.Context) error {
	return d.StartWithOptions(parent, DeployOptions{})
}

func (d *Deployer) StartWithOptions(parent context.Context, opts DeployOptions) error {
	d.mu.Lock()
	if d.curr.State == "running" {
		d.mu.Unlock()
		return ErrDeployInProgress
	}
	if !d.dockerAvailable(parent) {
		d.mu.Unlock()
		return ErrDockerSocketMissing
	}
	d.curr = DeployStatus{State: "running", StartedAt: time.Now(), Lines: []string{}}
	d.mu.Unlock()

	go d.run(parent, opts)
	return nil
}

func (d *Deployer) run(parent context.Context, opts DeployOptions) {
	ctx, cancel := context.WithTimeout(parent, 15*time.Minute)
	defer cancel()
	project := strings.TrimSpace(opts.ComposeProject)
	if project == "" {
		project = d.composeProject
	}
	services := opts.Services
	if len(services) == 0 {
		services = []string{"backend", "frontend", "nsmis-worker", "backup"}
	}
	if script := strings.TrimSpace(opts.Script); script != "" {
		d.append("Custom deploy script started. Upload/config volumes are preserved by docker-compose volume bindings.")
		if err := d.runCmdWithProject(ctx, []string{"sh", "-lc", script}, project); err != nil {
			d.finish("failed", 1, err.Error())
			return
		}
		d.append("Custom deploy script complete.")
		d.finish("success", 0, "")
		return
	}

	steps := [][]string{
		// Snapshot the current images by re-tagging :latest -> :previous so a
		// rollback can pin back to them.
		{"sh", "-c", `
		  for svc in backend frontend nsmis-worker backup; do
		    if docker image inspect "${COMPOSE_PROJECT:-ncs-online}-$svc:latest" >/dev/null 2>&1; then
		      docker tag "${COMPOSE_PROJECT:-ncs-online}-$svc:latest" "${COMPOSE_PROJECT:-ncs-online}-$svc:previous" || true
		    fi
		  done
		`},
		append([]string{"docker", "compose", "-p", project, "pull"}, services...),
		append([]string{"docker", "compose", "-p", project, "up", "-d", "--no-deps"}, services...),
	}

	exitCode := 0
	for i, step := range steps {
		d.append(fmt.Sprintf("──▶ Step %d/%d: %s", i+1, len(steps), strings.Join(step, " ")))
		if err := d.runCmdWithProject(ctx, step, project); err != nil {
			d.finish("failed", 1, err.Error())
			return
		}
	}
	d.append("──▶ Deploy complete. New containers up.")
	d.finish("success", exitCode, "")
}

func (d *Deployer) runCmd(ctx context.Context, argv []string) error {
	return d.runCmdWithProject(ctx, argv, d.composeProject)
}

func (d *Deployer) runCmdWithProject(ctx context.Context, argv []string, project string) error {
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Env = append(os.Environ(),
		"COMPOSE_PROJECT="+project,
		"DOCKER_BUILDKIT=1",
	)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}
	var wg sync.WaitGroup
	pipe := func(r io.Reader, prefix string) {
		defer wg.Done()
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 64*1024), 1024*1024)
		for sc.Scan() {
			d.append(prefix + sc.Text())
		}
	}
	wg.Add(2)
	go pipe(stdout, "")
	go pipe(stderr, "stderr: ")
	wg.Wait()
	return cmd.Wait()
}

func (d *Deployer) append(line string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.curr.Lines = append(d.curr.Lines, line)
	// Cap at 500 lines so the JSON payload stays sane.
	if len(d.curr.Lines) > 500 {
		d.curr.Lines = d.curr.Lines[len(d.curr.Lines)-500:]
	}
}

func (d *Deployer) finish(state string, exitCode int, errMsg string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.curr.State = state
	d.curr.EndedAt = time.Now()
	d.curr.ExitCode = exitCode
	d.curr.Error = errMsg
}

// Rollback re-tags :previous → :latest and restarts the same services.
// No-op if no :previous tag exists yet.
func (d *Deployer) Rollback(parent context.Context) error {
	d.mu.Lock()
	if d.curr.State == "running" {
		d.mu.Unlock()
		return ErrDeployInProgress
	}
	if !d.dockerAvailable(parent) {
		d.mu.Unlock()
		return ErrDockerSocketMissing
	}
	d.curr = DeployStatus{State: "running", StartedAt: time.Now(), Lines: []string{"──▶ Rolling back to previous images…"}}
	d.mu.Unlock()

	go func() {
		ctx, cancel := context.WithTimeout(parent, 5*time.Minute)
		defer cancel()
		steps := [][]string{
			{"sh", "-c", `
			  any=0
			  for svc in backend frontend nsmis-worker backup; do
			    if docker image inspect "${COMPOSE_PROJECT:-ncs-online}-$svc:previous" >/dev/null 2>&1; then
			      docker tag "${COMPOSE_PROJECT:-ncs-online}-$svc:previous" "${COMPOSE_PROJECT:-ncs-online}-$svc:latest"
			      any=1
			    fi
			  done
			  test "$any" = "1" || { echo "no :previous tags found — nothing to roll back to"; exit 2; }
			`},
			{"docker", "compose", "-p", d.composeProject, "up", "-d", "--no-deps", "backend", "frontend", "nsmis-worker", "backup"},
		}
		for i, step := range steps {
			d.append(fmt.Sprintf("──▶ Rollback step %d/%d: %s", i+1, len(steps), strings.Join(step, " ")))
			if err := d.runCmd(ctx, step); err != nil {
				d.finish("failed", 1, err.Error())
				return
			}
		}
		d.append("──▶ Rollback complete.")
		d.finish("success", 0, "")
	}()
	return nil
}
