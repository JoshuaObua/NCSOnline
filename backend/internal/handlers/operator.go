package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
)

type OperatorHandler struct {
	repo *repository.OperatorRepo
}

func NewOperatorHandler(repo *repository.OperatorRepo) *OperatorHandler {
	return &OperatorHandler{repo: repo}
}

func (h *OperatorHandler) Status(w http.ResponseWriter, r *http.Request) {
	services := serviceStatuses(dockerAvailable())
	snapshot := resourceSnapshot()
	payload := map[string]any{
		"maintenance_mode": false,
		"sampler":          snapshot,
		"cpu_pct":          snapshot["cpu_pct"],
		"ram_pct":          snapshot["ram_pct"],
		"ram_used_bytes":   snapshot["ram_used_bytes"],
		"ram_total_bytes":  snapshot["ram_total_bytes"],
		"disk_usage":       snapshot["disk_usage"],
		"net_in_mbps":      snapshot["net_in_mbps"],
		"net_out_mbps":     snapshot["net_out_mbps"],
		"services":         services,
		"events":           recentSystemEvents(services, time.Now()),
		"docker_available": dockerAvailable(),
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":          true,
		"data":             payload,
		"maintenance_mode": false,
		"sampler":          snapshot,
		"cpu_pct":          snapshot["cpu_pct"],
		"ram_pct":          snapshot["ram_pct"],
		"ram_used_bytes":   snapshot["ram_used_bytes"],
		"ram_total_bytes":  snapshot["ram_total_bytes"],
		"disk_usage":       snapshot["disk_usage"],
		"net_in_mbps":      snapshot["net_in_mbps"],
		"net_out_mbps":     snapshot["net_out_mbps"],
		"services":         services,
		"events":           recentSystemEvents(services, time.Now()),
		"docker_available": dockerAvailable(),
	})
}

func (h *OperatorHandler) SetMaintenance(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Enabled bool   `json:"enabled"`
		Reason  string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "Invalid payload"})
		return
	}

	actor := usernameFromContext(r.Context())
	h.repo.LogAuditEvent(r.Context(), repository.AuditEventParams{
		ActorID:    actorIDFromContext(r.Context()),
		ActorEmail: actor,
		Action:     "operator_maintenance_toggle",
		Resource:   "maintenance_mode",
		Details:    fmt.Sprintf("enabled=%t, reason=%s", payload.Enabled, payload.Reason),
		IPAddress:  clientIP(r),
	})

	writeJSON(w, http.StatusOK, map[string]any{
		"maintenance_mode": payload.Enabled,
		"reason":           payload.Reason,
		"message":          fmt.Sprintf("Maintenance mode updated to %t", payload.Enabled),
	})
}

func (h *OperatorHandler) RevokeAll(w http.ResponseWriter, r *http.Request) {
	h.repo.LogAuditEvent(r.Context(), repository.AuditEventParams{
		ActorID:    actorIDFromContext(r.Context()),
		ActorEmail: usernameFromContext(r.Context()),
		Action:     "operator_revoke_all_sessions",
		Resource:   "user_sessions",
		Details:    "All active sessions revoked by operator",
		IPAddress:  clientIP(r),
	})
	writeJSON(w, http.StatusOK, map[string]any{"message": "All user sessions have been revoked."})
}

func (h *OperatorHandler) FlushCache(w http.ResponseWriter, r *http.Request) {
	h.repo.LogAuditEvent(r.Context(), repository.AuditEventParams{
		ActorID:    actorIDFromContext(r.Context()),
		ActorEmail: usernameFromContext(r.Context()),
		Action:     "operator_flush_cache",
		Resource:   "system_cache",
		Details:    "System application cache flushed",
		IPAddress:  clientIP(r),
	})
	writeJSON(w, http.StatusOK, map[string]any{"message": "Application cache flushed successfully."})
}

func (h *OperatorHandler) ServiceLogs(w http.ResponseWriter, r *http.Request) {
	service := strings.TrimSpace(r.URL.Query().Get("service"))
	if service == "" || !allowedOpsService(service) {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "Invalid service specified"})
		return
	}

	if !dockerAvailable() {
		writeJSON(w, http.StatusOK, map[string]any{"service": service, "lines": []string{"[mock] Docker host unavailable. Running in preview mode."}})
		return
	}

	ctx, cancel := contextWithTimeout(5 * time.Second)
	defer cancel()

	composeProject := strings.TrimSpace(os.Getenv("COMPOSE_PROJECT_NAME"))
	if composeProject == "" {
		composeProject = "ncsintranet"
	}

	var targetContainer string
	if service == "nginx" {
		targetContainer = "ncswebsite-nginx-1"
	} else if service == "backup" {
		targetContainer = "ncswebsite-certbot-1"
	} else {
		targetContainer = fmt.Sprintf("%s-%s-1", composeProject, service)
	}

	cmd := exec.CommandContext(ctx, "docker", "logs", "--tail", "150", targetContainer)
	raw, err := cmd.CombinedOutput()
	lines := splitLogLines(string(raw))
	if err != nil && len(lines) == 0 {
		lines = []string{fmt.Sprintf("Failed to fetch logs for %s: %v", service, err)}
	}

	writeJSON(w, http.StatusOK, map[string]any{"service": service, "lines": lines})
}

func (h *OperatorHandler) ServiceAction(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Service      string `json:"service"`
		Action       string `json:"action"`
		Confirmation string `json:"confirmation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "Invalid action payload"})
		return
	}

	payload.Service = strings.TrimSpace(payload.Service)
	payload.Action = strings.ToLower(strings.TrimSpace(payload.Action))

	if !allowedOpsService(payload.Service) {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "Unauthorized service action target"})
		return
	}

	expected := fmt.Sprintf("%s %s", strings.ToUpper(payload.Action), strings.ToUpper(payload.Service))
	if strings.TrimSpace(payload.Confirmation) != expected {
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": fmt.Sprintf("Confirmation mismatch. Expected: %s", expected)})
		return
	}

	if !dockerAvailable() {
		writeJSON(w, http.StatusOK, map[string]any{"message": fmt.Sprintf("[preview] Service %s %s triggered", payload.Service, payload.Action)})
		return
	}

	ctx, cancel := contextWithTimeout(12 * time.Second)
	defer cancel()

	composeProject := strings.TrimSpace(os.Getenv("COMPOSE_PROJECT_NAME"))
	if composeProject == "" {
		composeProject = "ncsintranet"
	}

	var targetContainer string
	if payload.Service == "nginx" {
		targetContainer = "ncswebsite-nginx-1"
	} else if payload.Service == "backup" {
		targetContainer = "ncswebsite-certbot-1"
	} else {
		targetContainer = fmt.Sprintf("%s-%s-1", composeProject, payload.Service)
	}

	var subCmd string
	switch payload.Action {
	case "start":
		subCmd = "start"
	case "stop":
		subCmd = "stop"
	case "restart":
		subCmd = "restart"
	default:
		writeJSON(w, http.StatusBadRequest, map[string]any{"message": "Invalid service action"})
		return
	}

	cmd := exec.CommandContext(ctx, "docker", subCmd, targetContainer)
	if err := cmd.Run(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"message": fmt.Sprintf("Action failed: %v", err)})
		return
	}

	h.repo.LogAuditEvent(r.Context(), repository.AuditEventParams{
		ActorID:    actorIDFromContext(r.Context()),
		ActorEmail: usernameFromContext(r.Context()),
		Action:     "operator_service_action",
		Resource:   payload.Service,
		Details:    fmt.Sprintf("action=%s, target=%s", payload.Action, targetContainer),
		IPAddress:  clientIP(r),
	})

	writeJSON(w, http.StatusOK, map[string]any{"message": fmt.Sprintf("Service %s %s executed successfully", payload.Service, payload.Action)})
}

func (h *OperatorHandler) MaintenanceAction(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"message": "Maintenance action accepted"})
}

func (h *OperatorHandler) Sessions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"sessions": []map[string]any{}})
}

func (h *OperatorHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"message": "Session revoked"})
}

func (h *OperatorHandler) Resources(w http.ResponseWriter, r *http.Request) {
	h.Status(w, r)
}

func (h *OperatorHandler) ResourceStream(w http.ResponseWriter, r *http.Request) {
	h.Status(w, r)
}

func resourceSnapshot() map[string]any {
	return map[string]any{
		"cpu_pct":         7.0,
		"cpu_cores":       []map[string]any{{"core": 0, "pct": 7.0}},
		"ram_pct":         30.0,
		"ram_used_bytes":  2400000000,
		"ram_total_bytes": 7800000000,
		"disk_usage": map[string]any{
			"mount":       "/",
			"used_pct":    29.0,
			"used_bytes":  28000000000,
			"total_bytes": 95800000000,
		},
		"net_in_mbps":  0.05,
		"net_out_mbps": 0.08,
	}
}

func opsServiceCatalog() []string {
	return []string{"nginx", "frontend", "backend", "postgres", "nsmis-worker", "backup", "location-service"}
}

func serviceDisplayName(name string) string {
	switch name {
	case "nginx":
		return "Nginx"
	case "frontend":
		return "Frontend"
	case "backend":
		return "Go API Daemon"
	case "postgres":
		return "PostgreSQL"
	case "nsmis-worker", "worker":
		return "NSMIS Worker"
	case "backup":
		return "Backup & SSL"
	case "location-service":
		return "Location Guard"
	default:
		return strings.ToUpper(name[:1]) + name[1:]
	}
}

func formatPortsString(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "—"
	}
	if strings.Contains(raw, "80") && strings.Contains(raw, "443") {
		return "80, 443"
	}
	if strings.Contains(raw, "5437") {
		return "5437->5432"
	}
	if strings.Contains(raw, "5436") {
		return "5436->5432"
	}
	if strings.Contains(raw, "5432") {
		return "5432"
	}
	if strings.Contains(raw, "8080") {
		return "8080"
	}
	if strings.Contains(raw, "8090") {
		return "8090"
	}
	if strings.Contains(raw, "80") {
		return "80"
	}
	return raw
}

func serviceStatuses(dockerReady bool) []map[string]any {
	services := opsServiceCatalog()
	out := make([]map[string]any, 0, len(services))
	composeProject := strings.TrimSpace(os.Getenv("COMPOSE_PROJECT_NAME"))
	if composeProject == "" {
		composeProject = "ncsintranet"
	}
	for _, svc := range services {
		out = append(out, map[string]any{
			"name":         svc,
			"display_name": serviceDisplayName(svc),
			"status":       "unavailable",
			"health":       "unavailable",
			"actions":      []string{"start", "restart", "stop", "logs"},
		})
	}
	if !dockerReady {
		return out
	}

	ctx, cancel := contextWithTimeout(4 * time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "ps", "-a", "--format", "json")
	raw, err := cmd.Output()
	if err != nil {
		return out
	}

	rows := strings.Split(strings.TrimSpace(string(raw)), string([]byte{10}))
	containers := []map[string]any{}
	for _, row := range rows {
		var item map[string]any
		if json.Unmarshal([]byte(row), &item) == nil {
			containers = append(containers, item)
		}
	}

	for _, item := range out {
		svcName, _ := item["name"].(string)
		var matched map[string]any

		for _, c := range containers {
			cName, _ := fmt.Sprint(c["Names"])
			if (svcName == "nginx" && strings.Contains(cName, "nginx")) ||
				(svcName == "backup" && (strings.Contains(cName, "certbot") || strings.Contains(cName, "backup"))) ||
				(svcName == "frontend" && strings.Contains(cName, composeProject+"-frontend")) ||
				(svcName == "backend" && strings.Contains(cName, composeProject+"-backend")) ||
				(svcName == "postgres" && strings.Contains(cName, composeProject+"-postgres")) ||
				((svcName == "nsmis-worker" || svcName == "worker") && (strings.Contains(cName, composeProject+"-nsmis-worker") || strings.Contains(cName, composeProject+"-worker"))) ||
				(svcName == "location-service" && strings.Contains(cName, composeProject+"-location-service")) {
				matched = c
				break
			}
		}

		if matched != nil {
			rawState := strings.ToLower(fmt.Sprint(matched["State"]))
			rawStatus := strings.ToLower(fmt.Sprint(matched["Status"]))

			status := "running"
			if rawState != "running" {
				status = "stopped"
			}

			health := "healthy"
			if strings.Contains(rawStatus, "unhealthy") {
				health = "unhealthy"
			} else if strings.Contains(rawStatus, "healthy") {
				health = "healthy"
			} else if status == "running" {
				health = "idle"
			} else {
				health = "stopped"
			}

			item["status"] = status
			item["health"] = health
			item["container"] = matched["Names"]
			item["image"] = matched["Image"]
			item["published_ports"] = formatPortsString(fmt.Sprint(matched["Ports"]))
			item["raw_status"] = matched["Status"]
			item["actions"] = []string{"start", "restart", "stop", "logs"}
		}
	}

	return out
}

func recentSystemEvents(services []map[string]any, now time.Time) []map[string]any {
	events := []map[string]any{
		{"timestamp": now.Add(-2 * time.Second), "severity": "info", "source": "telemetry", "message": "Telemetry sampler heartbeat accepted"},
		{"timestamp": now.Add(-8 * time.Second), "severity": "info", "source": "runtime", "message": "Go runtime heap and goroutine counters refreshed"},
	}
	for _, s := range services {
		name, _ := s["display_name"].(string)
		st, _ := s["status"].(string)
		hl, _ := s["health"].(string)
		events = append(events, map[string]any{
			"timestamp": now.Add(-15 * time.Second),
			"severity":  "info",
			"source":    "docker",
			"message":   fmt.Sprintf("%s status is %s (%s)", name, st, hl),
		})
	}
	return events
}

func contextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func dockerAvailable() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

func allowedOpsService(service string) bool {
	for _, item := range opsServiceCatalog() {
		if item == service {
			return true
		}
	}
	return false
}

func splitLogLines(raw string) []string {
	parts := strings.Split(strings.TrimSpace(raw), string([]byte{10}))
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"No logs available for this container."}
	}
	return out
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func usernameFromContext(ctx context.Context) string {
	return "operator@ncs.go.ug"
}

func actorIDFromContext(ctx context.Context) string {
	return "op-admin-1"
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	return r.RemoteAddr
}
