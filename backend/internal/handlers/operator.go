package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/maintenance"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

type OperatorHandler struct {
	repo  *repository.OperatorRepo
	state *maintenance.State
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
		Enabled        bool       `json:"enabled"`
		Reason         string     `json:"reason"`
		ScheduledStart *time.Time `json:"scheduled_start"`
		ExpectedEnd    *time.Time `json:"expected_end"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}
	if input.Enabled && len(input.Reason) < 5 {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "A maintenance reason is required")
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
	s, err := h.repo.SaveMaintenance(r.Context(), persistEnabled, input.Reason, scheduledStart, expectedEnd, actor)
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

var resourceUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return r.Header.Get("Origin") == "" || r.Header.Get("Origin") == "http://"+r.Host || r.Header.Get("Origin") == "https://"+r.Host
}}

func resourceSnapshot() map[string]any {
	cpuPct, _ := cpu.Percent(0, false)
	virtual, _ := mem.VirtualMemory()
	swap, _ := mem.SwapMemory()
	usage, _ := disk.Usage("/")
	avg, _ := load.Avg()
	info, _ := host.Info()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return map[string]any{"cpu_percent": cpuPct, "memory": virtual, "swap": swap, "disk": usage, "load": avg, "host": info, "goroutines": runtime.NumGoroutine(), "heap_alloc_bytes": ms.HeapAlloc, "timestamp": time.Now().UTC()}
}
func (h *OperatorHandler) Resources(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, resourceSnapshot())
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
		if err := conn.WriteJSON(resourceSnapshot()); err != nil {
			return
		}
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
		}
	}
}
