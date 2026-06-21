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
		Enabled     bool       `json:"enabled"`
		Reason      string     `json:"reason"`
		ExpectedEnd *time.Time `json:"expected_end"`
	}
	if json.NewDecoder(r.Body).Decode(&input) != nil || (input.Enabled && len(input.Reason) < 5) {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "A maintenance reason is required")
		return
	}
	s, err := h.repo.SaveMaintenance(r.Context(), input.Enabled, input.Reason, input.ExpectedEnd, actor)
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not update maintenance mode")
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
