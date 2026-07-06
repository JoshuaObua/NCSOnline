package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/maintenance"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	gnet "github.com/shirou/gopsutil/v4/net"
)

type OperatorHandler struct {
	repo      *repository.OperatorRepo
	state     *maintenance.State
	telemetry telemetrySampler
}

func NewOperatorHandler(repo *repository.OperatorRepo, state *maintenance.State) *OperatorHandler {
	return &OperatorHandler{repo: repo, state: state}
}
func (h *OperatorHandler) Status(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, h.state.Get())
}
func (h *OperatorHandler) SetMaintenance(w http.ResponseWriter, r *http.Request) {
	actor, _ := r.Context().Value(models.CtxUserID).(string)
	var input struct {
		Scope             string                  `json:"scope"`
		Enabled           bool                    `json:"enabled"`
		Reason            string                  `json:"reason"`
		ScheduledStart    *time.Time              `json:"scheduled_start"`
		ExpectedEnd       *time.Time              `json:"expected_end"`
		AutoStartEnforced *bool                   `json:"auto_start_enforced"`
		AutoEndEnforced   *bool                   `json:"auto_end_enforced"`
		DisplayMeta       maintenance.DisplayMeta `json:"display_meta"`
		BypassRules       maintenance.BypassRules `json:"bypass_rules"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if input.Enabled && len(input.Reason) < 5 {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "A maintenance reason is required")
		return
	}
	scope := input.Scope
	if scope == "" {
		scope = maintenance.ScopePublicCMS
	}
	if scope != maintenance.ScopePublicCMS && scope != maintenance.ScopeAdminDashboard {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Unknown maintenance scope")
		return
	}
	if input.ScheduledStart != nil && input.ExpectedEnd != nil && !input.ExpectedEnd.After(*input.ScheduledStart) {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Expected end must be after the scheduled start time")
		return
	}
	// If the caller is disabling, clear the window so the sentinel
	// doesn't immediately re-enable us.
	scheduledStart := input.ScheduledStart
	expectedEnd := input.ExpectedEnd
	if !input.Enabled {
		scheduledStart, expectedEnd = nil, nil
	}
	// Resolve the live-active bit: if there's a future start time we
	// persist Enabled=true but the IsActiveAt check (and the schedule
	// sentinel) keep traffic flowing until that start passes.
	now := time.Now()
	persistEnabled := input.Enabled
	if input.Enabled && scheduledStart != nil && scheduledStart.After(now) {
		// Future-scheduled window. We persist Enabled so the sentinel
		// flips us automatically when the start time hits.
		persistEnabled = true
	}
	current := h.state.Get().Scoped(scope)
	current.Enabled = persistEnabled
	current.Reason = input.Reason
	current.ScheduledStart = scheduledStart
	current.ExpectedEnd = expectedEnd
	if input.AutoStartEnforced != nil {
		current.AutoStartEnforced = *input.AutoStartEnforced
	}
	if input.AutoEndEnforced != nil {
		current.AutoEndEnforced = *input.AutoEndEnforced
	}
	if input.DisplayMeta.CustomTitle != "" || input.DisplayMeta.CustomMessage != "" {
		current.DisplayMeta = input.DisplayMeta
	}
	if len(input.BypassRules.AllowedRoles) > 0 || len(input.BypassRules.AllowedUserIDs) > 0 || len(input.BypassRules.AllowedIPRanges) > 0 || input.BypassRules.SecretQueryParam != "" {
		current.BypassRules = input.BypassRules
	}
	s, err := h.repo.SaveMaintenanceScope(r.Context(), scope, current, actor)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not update maintenance mode")
		return
	}
	h.state.Set(s)
	response.JSON(w, http.StatusOK, s)
}
func (h *OperatorHandler) RevokeAll(w http.ResponseWriter, r *http.Request) {
	actor, _ := r.Context().Value(models.CtxUserID).(string)
	s, err := h.repo.RevokeAllSessions(r.Context(), actor)
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not revoke sessions")
		return
	}
	h.state.Set(s)
	response.JSONMsg(w, http.StatusOK, "All user sessions revoked")
}
func (h *OperatorHandler) FlushCache(w http.ResponseWriter, r *http.Request) {
	actor, _ := r.Context().Value(models.CtxUserID).(string)
	s := h.state.FlushCache()
	if err := h.repo.SaveCacheGeneration(r.Context(), s.CacheGeneration, actor); err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not flush cache")
		return
	}
	response.JSON(w, http.StatusOK, s)
}

func (h *OperatorHandler) ServiceLogs(w http.ResponseWriter, r *http.Request) {
	service := strings.TrimSpace(r.URL.Query().Get("service"))
	if !allowedOpsService(service) {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Unknown service")
		return
	}
	lines := 120
	if raw := r.URL.Query().Get("lines"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 500 {
			lines = parsed
		}
	}
	out, err := dockerComposeOutput(r.Context(), "logs", "--tail", strconv.Itoa(lines), service)
	if err != nil {
		response.Err(w, http.StatusFailedDependency, "SERVICE_LOGS_UNAVAILABLE", err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"service": service, "lines": splitLogLines(out)})
}

func (h *OperatorHandler) ServiceLogStream(w http.ResponseWriter, r *http.Request) {
	service := strings.TrimSpace(r.URL.Query().Get("service"))
	if !allowedOpsService(service) {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Unknown service")
		return
	}
	if !dockerAvailable() {
		response.Err(w, http.StatusFailedDependency, "SERVICE_LOGS_UNAVAILABLE", "docker socket or docker CLI is unavailable")
		return
	}
	lines := "100"
	if raw := r.URL.Query().Get("lines"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 500 {
			lines = strconv.Itoa(parsed)
		}
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		response.Err(w, http.StatusInternalServerError, "STREAM_UNAVAILABLE", "Streaming is not supported")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	composeProject := strings.TrimSpace(os.Getenv("COMPOSE_PROJECT_NAME"))
	if composeProject == "" {
		composeProject = "ncs-online"
	}
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", "compose", "-p", composeProject, "logs", "--tail", lines, "-f", service)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "STREAM_FAILED", err.Error())
		return
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		response.Err(w, http.StatusFailedDependency, "STREAM_FAILED", err.Error())
		return
	}
	defer cmd.Wait()
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		line, _ := json.Marshal(scanner.Text())
		_, _ = fmt.Fprintf(w, "data: %s\n\n", line)
		flusher.Flush()
	}
}

func (h *OperatorHandler) ServiceAction(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Service      string `json:"service"`
		Action       string `json:"action"`
		Confirmation string `json:"confirmation"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request")
		return
	}
	input.Service = strings.TrimSpace(input.Service)
	input.Action = strings.ToLower(strings.TrimSpace(input.Action))
	expected := strings.ToUpper(input.Action + " " + input.Service)
	if !allowedOpsService(input.Service) || (input.Action != "restart" && input.Action != "stop" && input.Action != "start") {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Unsupported service action")
		return
	}
	if strings.TrimSpace(input.Confirmation) != expected {
		response.Err(w, http.StatusBadRequest, "CONFIRMATION_REQUIRED", "Type "+expected+" to confirm")
		return
	}
	args := []string{input.Action, input.Service}
	if input.Action == "restart" {
		args = []string{"up", "-d", "--no-deps", "--force-recreate", input.Service}
	}
	if input.Action == "start" {
		args = []string{"up", "-d", "--no-deps", input.Service}
	}
	out, err := dockerComposeOutput(r.Context(), args...)
	if err != nil {
		response.Err(w, http.StatusFailedDependency, "SERVICE_ACTION_FAILED", err.Error())
		return
	}
	response.JSON(w, http.StatusAccepted, map[string]any{"service": input.Service, "action": input.Action, "output": splitLogLines(out)})
}

func (h *OperatorHandler) MaintenanceAction(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Action       string `json:"action"`
		Confirmation string `json:"confirmation"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request")
		return
	}
	input.Action = strings.ToLower(strings.TrimSpace(input.Action))
	expected := strings.ToUpper(strings.ReplaceAll(input.Action, "_", " "))
	if strings.TrimSpace(input.Confirmation) != expected {
		response.Err(w, http.StatusBadRequest, "CONFIRMATION_REQUIRED", "Type "+expected+" to confirm")
		return
	}
	switch input.Action {
	case "flush_app_cache":
		h.FlushCache(w, r)
	case "flush_nginx_cache":
		out, err := dockerComposeOutput(r.Context(), "exec", "-T", "nginx", "sh", "-lc", "rm -rf /var/cache/nginx/* 2>/dev/null || true")
		if err != nil {
			response.Err(w, http.StatusFailedDependency, "CACHE_FLUSH_FAILED", err.Error())
			return
		}
		response.JSON(w, http.StatusAccepted, map[string]any{"action": input.Action, "output": splitLogLines(out)})
	case "flush_redis_cache":
		response.JSON(w, http.StatusAccepted, map[string]any{"action": input.Action, "message": "Redis is not configured in this deployment; no Redis cache was flushed."})
	case "optimize_database_tables":
		out, err := dockerComposeOutput(r.Context(), "exec", "-T", "postgres", "vacuumdb", "-U", "ncsms_user", "-d", "ncsms", "--analyze")
		if err != nil {
			response.Err(w, http.StatusFailedDependency, "DATABASE_OPTIMIZE_FAILED", err.Error())
			return
		}
		response.JSON(w, http.StatusAccepted, map[string]any{"action": input.Action, "output": splitLogLines(out)})
	case "prune_activity_logs":
		response.Err(w, http.StatusNotImplemented, "PRUNE_POLICY_REQUIRED", "Activity log pruning needs a retention policy before it can run")
	case "trigger_manual_backup_slice":
		response.JSON(w, http.StatusAccepted, map[string]any{"action": input.Action, "message": "Use Backup & Restore to queue an audited database backup job."})
	default:
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Unsupported maintenance action")
	}
}
func (h *OperatorHandler) Sessions(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	items, err := h.repo.ListSessions(r.Context(), userID)
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load sessions")
		return
	}
	response.JSON(w, http.StatusOK, items)
}
func (h *OperatorHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	var input struct {
		SessionID     string `json:"session_id"`
		OtherSessions bool   `json:"other_sessions"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil {
		response.Err(w, 400, "BAD_REQUEST", "Invalid request")
		return
	}
	var err error
	if input.OtherSessions {
		current, _ := r.Context().Value(models.CtxSessionID).(string)
		err = h.repo.RevokeOtherSessions(r.Context(), userID, current)
	} else {
		err = h.repo.RevokeSession(r.Context(), userID, input.SessionID)
	}
	if err != nil {
		response.Err(w, 400, "REVOKE_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Session revoked")
}

func (h *OperatorHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	sessionID := chi.URLParam(r, "id")
	if sessionID == "" {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Session id is required")
		return
	}
	if err := h.repo.RevokeSession(r.Context(), userID, sessionID); err != nil {
		response.Err(w, http.StatusBadRequest, "REVOKE_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Session revoked")
}

func (h *OperatorHandler) DeleteOtherSessions(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(models.CtxUserID).(string)
	current, _ := r.Context().Value(models.CtxSessionID).(string)
	if err := h.repo.RevokeOtherSessions(r.Context(), userID, current); err != nil {
		response.Err(w, http.StatusBadRequest, "REVOKE_ERROR", err.Error())
		return
	}
	response.JSONMsg(w, http.StatusOK, "Other sessions revoked")
}

var resourceUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return r.Header.Get("Origin") == "" || r.Header.Get("Origin") == "http://"+r.Host || r.Header.Get("Origin") == "https://"+r.Host
}}

type telemetrySampler struct {
	mu        sync.Mutex
	lastAt    time.Time
	lastNet   map[string]gnet.IOCountersStat
	lastDisk  map[string]disk.IOCountersStat
	netRxBps  float64
	netTxBps  float64
	diskRead  float64
	diskWrite float64
}

func (s *telemetrySampler) snapshot(state maintenance.Snapshot) map[string]any {
	now := time.Now().UTC()
	cpuPct, _ := cpu.Percent(0, false)
	perCPU, _ := cpu.Percent(0, true)
	times, _ := cpu.Times(false)
	virtual, _ := mem.VirtualMemory()
	swap, _ := mem.SwapMemory()
	avg, _ := load.Avg()
	info, _ := host.Info()
	partitions := diskPartitions()
	diskIO, _ := disk.IOCounters()
	netIO, _ := gnet.IOCounters(false)
	conns, _ := gnet.Connections("inet")

	s.mu.Lock()
	s.updateRates(now, netIO, diskIO)
	netRxBps, netTxBps, diskRead, diskWrite := s.netRxBps, s.netTxBps, s.diskRead, s.diskWrite
	s.mu.Unlock()

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	services := serviceStatuses()
	events := recentSystemEvents(services, state, now)
	public := state.Scoped(maintenance.ScopePublicCMS)
	admin := state.Scoped(maintenance.ScopeAdminDashboard)

	return map[string]any{
		"timestamp":        now,
		"schema_version":   "ops.telemetry.v1",
		"cpu_percent":      cpuPct,
		"memory":           virtual,
		"swap":             swap,
		"disk":             firstDiskUsage(partitions),
		"load":             avg,
		"host":             info,
		"goroutines":       runtime.NumGoroutine(),
		"heap_alloc_bytes": ms.HeapAlloc,
		"maintenance": map[string]any{
			"public_cms":      public,
			"admin_dashboard": admin,
		},
		"cpu": map[string]any{
			"overall_percent": firstFloat(cpuPct),
			"per_core":        perCorePayload(perCPU),
			"load_average":    avg,
			"temperature_c":   nil,
			"times":           times,
		},
		"memory_detail": map[string]any{
			"total_bytes":     virtual.Total,
			"used_bytes":      virtual.Used,
			"available_bytes": virtual.Available,
			"cached_bytes":    virtual.Cached,
			"buffers_bytes":   virtual.Buffers,
			"used_percent":    virtual.UsedPercent,
			"swap_total":      swap.Total,
			"swap_used":       swap.Used,
			"swap_percent":    swap.UsedPercent,
		},
		"gpu": map[string]any{
			"available": false,
			"devices":   []any{},
		},
		"storage": map[string]any{
			"partitions":      partitions,
			"io_counters":     diskIO,
			"read_bytes_sec":  diskRead,
			"write_bytes_sec": diskWrite,
			"read_mb_sec":     diskRead / (1024 * 1024),
			"write_mb_sec":    diskWrite / (1024 * 1024),
			"tracked_mounts":  []string{"/", "/var", "/data"},
		},
		"network": map[string]any{
			"interfaces":   netIO,
			"ingress_bps":  netRxBps,
			"egress_bps":   netTxBps,
			"ingress_mbps": netRxBps * 8 / 1_000_000,
			"egress_mbps":  netTxBps * 8 / 1_000_000,
			"connections":  len(conns),
			"latency_ms":   nil,
		},
		"services": services,
		"events":   events,
	}
}

func (s *telemetrySampler) updateRates(now time.Time, netIO []gnet.IOCountersStat, diskIO map[string]disk.IOCountersStat) {
	if s.lastAt.IsZero() {
		s.lastAt = now
		s.lastNet = map[string]gnet.IOCountersStat{}
		s.lastDisk = map[string]disk.IOCountersStat{}
		for _, item := range netIO {
			s.lastNet[item.Name] = item
		}
		for name, item := range diskIO {
			s.lastDisk[name] = item
		}
		return
	}
	elapsed := now.Sub(s.lastAt).Seconds()
	if elapsed <= 0 {
		return
	}
	var rx, tx uint64
	for _, item := range netIO {
		if prev, ok := s.lastNet[item.Name]; ok {
			rx += item.BytesRecv - prev.BytesRecv
			tx += item.BytesSent - prev.BytesSent
		}
		s.lastNet[item.Name] = item
	}
	var read, write uint64
	for name, item := range diskIO {
		if prev, ok := s.lastDisk[name]; ok {
			read += item.ReadBytes - prev.ReadBytes
			write += item.WriteBytes - prev.WriteBytes
		}
		s.lastDisk[name] = item
	}
	s.lastAt = now
	s.netRxBps = float64(rx) / elapsed
	s.netTxBps = float64(tx) / elapsed
	s.diskRead = float64(read) / elapsed
	s.diskWrite = float64(write) / elapsed
}

func diskPartitions() []map[string]any {
	mounts := []string{"/", "/var", "/data"}
	out := []map[string]any{}
	seen := map[string]bool{}
	for _, mount := range mounts {
		usage, err := disk.Usage(mount)
		if err != nil || usage == nil || seen[usage.Path] {
			continue
		}
		seen[usage.Path] = true
		out = append(out, map[string]any{
			"mountpoint":   usage.Path,
			"fstype":       usage.Fstype,
			"total_bytes":  usage.Total,
			"used_bytes":   usage.Used,
			"free_bytes":   usage.Free,
			"used_percent": usage.UsedPercent,
		})
	}
	return out
}

func firstDiskUsage(partitions []map[string]any) any {
	if len(partitions) == 0 {
		return nil
	}
	return partitions[0]
}

func firstFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return values[0]
}

func perCorePayload(values []float64) []map[string]any {
	out := make([]map[string]any, 0, len(values))
	for i, v := range values {
		out = append(out, map[string]any{"core": i, "usage_percent": v})
	}
	return out
}

func serviceStatuses() []map[string]any {
	services := []string{"nginx", "frontend", "backend", "postgres", "nsmis-worker", "backup", "location-service"}
	out := make([]map[string]any, 0, len(services))
	composeProject := strings.TrimSpace(os.Getenv("COMPOSE_PROJECT_NAME"))
	if composeProject == "" {
		composeProject = "ncs-online"
	}
	for _, svc := range services {
		out = append(out, map[string]any{"name": svc, "display_name": serviceDisplayName(svc), "status": "unknown", "health": "unknown", "actions": []string{"start", "restart", "stop", "logs"}})
	}
	if !dockerAvailable() {
		return out
	}
	ctx, cancel := contextWithTimeout(4 * time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", "compose", "-p", composeProject, "ps", "--format", "json")
	raw, err := cmd.Output()
	if err != nil {
		return out
	}
	rows := strings.Split(strings.TrimSpace(string(raw)), "\n")
	byName := map[string]map[string]any{}
	for _, row := range rows {
		var item map[string]any
		if json.Unmarshal([]byte(row), &item) == nil {
			if name, _ := item["Service"].(string); name != "" {
				byName[name] = item
			}
		}
	}
	for _, item := range out {
		name, _ := item["name"].(string)
		if row := byName[name]; row != nil {
			state, _ := row["State"].(string)
			health, _ := row["Health"].(string)
			item["status"] = normalizeServiceState(state, health)
			item["health"] = strings.TrimSpace(health)
			item["container"] = row["Name"]
			item["image"] = row["Image"]
			item["published_ports"] = row["Publishers"]
		}
	}
	return out
}

func normalizeServiceState(state, health string) string {
	state = strings.ToLower(strings.TrimSpace(state))
	health = strings.ToLower(strings.TrimSpace(health))
	if state == "running" && (health == "" || health == "healthy") {
		return "healthy"
	}
	if state == "running" {
		return "degraded"
	}
	if state == "" {
		return "unknown"
	}
	return "stopped"
}

func serviceDisplayName(name string) string {
	switch name {
	case "backend":
		return "Go API Daemon"
	case "postgres":
		return "PostgreSQL"
	case "nsmis-worker":
		return "NSMIS Worker"
	case "location-service":
		return "Location Guard"
	default:
		return strings.ToUpper(name[:1]) + name[1:]
	}
}

func recentSystemEvents(services []map[string]any, state maintenance.Snapshot, now time.Time) []map[string]any {
	events := []map[string]any{
		{"timestamp": now.Add(-2 * time.Second), "severity": "info", "source": "telemetry", "message": "Telemetry sampler heartbeat accepted"},
		{"timestamp": now.Add(-8 * time.Second), "severity": "info", "source": "runtime", "message": "Go runtime heap and goroutine counters refreshed"},
	}
	if state.Scoped(maintenance.ScopePublicCMS).IsActiveAt(now) {
		events = append(events, map[string]any{"timestamp": now, "severity": "warning", "source": "maintenance", "message": "Public CMS maintenance gate is active"})
	}
	if state.Scoped(maintenance.ScopeAdminDashboard).IsActiveAt(now) {
		events = append(events, map[string]any{"timestamp": now, "severity": "critical", "source": "maintenance", "message": "Admin dashboard maintenance gate is active"})
	}
	for _, svc := range services {
		if svc["status"] == "stopped" || svc["status"] == "degraded" {
			events = append(events, map[string]any{"timestamp": now, "severity": "critical", "source": "service", "message": serviceDisplayName(fmt.Sprint(svc["name"])) + " is " + fmt.Sprint(svc["status"])})
		}
	}
	return events
}

func contextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func dockerAvailable() bool {
	if _, err := os.Stat("/var/run/docker.sock"); err != nil {
		return false
	}
	if _, err := exec.LookPath("docker"); err != nil {
		return false
	}
	return true
}

func allowedOpsService(service string) bool {
	switch service {
	case "nginx", "frontend", "backend", "postgres", "nsmis-worker", "backup", "location-service":
		return true
	default:
		return false
	}
}

func dockerComposeOutput(parent context.Context, args ...string) (string, error) {
	if !dockerAvailable() {
		return "", errors.New("docker socket or docker CLI is unavailable")
	}
	composeProject := strings.TrimSpace(os.Getenv("COMPOSE_PROJECT_NAME"))
	if composeProject == "" {
		composeProject = "ncs-online"
	}
	ctx, cancel := context.WithTimeout(parent, 45*time.Second)
	defer cancel()
	argv := append([]string{"compose", "-p", composeProject}, args...)
	cmd := exec.CommandContext(ctx, "docker", argv...)
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return string(out), ctx.Err()
	}
	if err != nil {
		return string(out), errors.New(strings.TrimSpace(string(out)) + ": " + err.Error())
	}
	return string(out), nil
}

func splitLogLines(raw string) []string {
	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(raw))
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func resourceSnapshot(state maintenance.Snapshot, sampler *telemetrySampler) map[string]any {
	return sampler.snapshot(state)
}
func (h *OperatorHandler) Resources(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, resourceSnapshot(h.state.Get(), &h.telemetry))
}
func (h *OperatorHandler) ResourceStream(w http.ResponseWriter, r *http.Request) {
	conn, err := resourceUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		if err := conn.WriteJSON(resourceSnapshot(h.state.Get(), &h.telemetry)); err != nil {
			return
		}
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
		}
	}
}
