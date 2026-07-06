package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/config"
	"github.com/atenimedia-llc/ncs-online/backend/internal/middleware"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	repo     *repository.AnalyticsRepo
	location *middleware.LocationClient
	salt     string
	queue    chan analyticsIntake
}

type analyticsIntake struct {
	EventName              string `json:"event_name"`
	SessionID              string `json:"session_id"`
	Path                   string `json:"path"`
	Referrer               string `json:"referrer"`
	ScreenResolution       string `json:"screen_resolution"`
	Language               string `json:"language"`
	SessionDurationSeconds int    `json:"session_duration_seconds"`
	InteractionCount       int    `json:"interaction_count"`
	UserAgent              string
	IP                     string
}

func NewAnalyticsHandler(repo *repository.AnalyticsRepo, cfg *config.Config) *AnalyticsHandler {
	h := &AnalyticsHandler{
		repo:     repo,
		location: middleware.NewLocationClient(cfg.LocationServiceURL, cfg.LocationServiceToken, cfg.LocationTimeout),
		salt:     cfg.AnalyticsHashSecret,
		queue:    make(chan analyticsIntake, 2048),
	}
	go h.worker()
	return h
}

func (h *AnalyticsHandler) Track(w http.ResponseWriter, r *http.Request) {
	var req analyticsIntake
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32<<10)).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid analytics payload")
		return
	}
	req.UserAgent = r.UserAgent()
	req.IP = clientIP(r)
	req.EventName = normalizeAnalyticsEvent(req.EventName)
	req.Path = normalizeAnalyticsPath(req.Path)
	req.Referrer = trim(req.Referrer, 500)
	req.ScreenResolution = trim(req.ScreenResolution, 40)
	req.Language = trim(req.Language, 80)
	if req.SessionID == "" {
		req.SessionID = uuid.NewString()
	}
	if req.SessionDurationSeconds < 0 {
		req.SessionDurationSeconds = 0
	}
	if req.SessionDurationSeconds > 24*60*60 {
		req.SessionDurationSeconds = 24 * 60 * 60
	}
	select {
	case h.queue <- req:
	default:
		// Preserve public page latency under load. Dropping analytics is safer
		// than delaying the visitor experience or applying backpressure.
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AnalyticsHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	out, err := h.repo.Dashboard(r.Context(), days)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not load analytics")
		return
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *AnalyticsHandler) worker() {
	for item := range h.queue {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = h.persist(ctx, item)
		cancel()
	}
}

func (h *AnalyticsHandler) persist(ctx context.Context, item analyticsIntake) error {
	loc, err := h.location.Locate(ctx, item.IP, item.UserAgent, "")
	if err != nil {
		loc = fallbackLocation(item.UserAgent)
	}
	event := &models.AnalyticsEvent{
		SessionID:              item.SessionID,
		VisitorHash:            h.visitorHash(item.IP),
		EventName:              item.EventName,
		Path:                   item.Path,
		Referrer:               item.Referrer,
		Source:                 trafficSource(item.Referrer),
		Country:                firstNonEmpty(loc.Country, "Unknown"),
		Region:                 loc.Region,
		City:                   loc.City,
		Browser:                firstNonEmpty(loc.Browser, "Unknown"),
		OS:                     firstNonEmpty(loc.Platform, "Unknown"),
		DeviceType:             firstNonEmpty(loc.DeviceType, "Unknown"),
		ScreenResolution:       item.ScreenResolution,
		Language:               item.Language,
		SessionDurationSeconds: item.SessionDurationSeconds,
		IsBounce:               item.EventName == "page_exit" && item.SessionDurationSeconds < 10 && item.InteractionCount == 0,
	}
	return h.repo.InsertEvent(ctx, event)
}

func (h *AnalyticsHandler) visitorHash(ip string) string {
	day := time.Now().UTC().Format("2006-01-02")
	mac := hmac.New(sha256.New, []byte(h.salt+":"+day))
	mac.Write([]byte(ip))
	return hex.EncodeToString(mac.Sum(nil))
}

func clientIP(r *http.Request) string {
	for _, header := range []string{"CF-Connecting-IP", "X-Sentinel-Client-IP", "X-Real-IP"} {
		if ip := normalizeIP(r.Header.Get(header)); ip != "" {
			return ip
		}
	}
	if chain := r.Header.Get("X-Forwarded-For"); chain != "" {
		for _, part := range strings.Split(chain, ",") {
			if ip := normalizeIP(part); ip != "" {
				return ip
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		if ip := normalizeIP(host); ip != "" {
			return ip
		}
	}
	return normalizeIP(r.RemoteAddr)
}

func normalizeIP(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	parsed := net.ParseIP(value)
	if parsed == nil {
		return ""
	}
	return parsed.String()
}

func normalizeAnalyticsEvent(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "click":
		return "click"
	case "page_exit":
		return "page_exit"
	default:
		return "page_view"
	}
}

func normalizeAnalyticsPath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "/"
	}
	if len(value) > 300 {
		value = value[:300]
	}
	if parsed, err := url.Parse(value); err == nil && parsed.Path != "" {
		value = parsed.Path
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	return value
}

func trafficSource(referrer string) string {
	referrer = strings.TrimSpace(strings.ToLower(referrer))
	if referrer == "" {
		return "Direct"
	}
	parsed, err := url.Parse(referrer)
	if err != nil {
		return "External Link"
	}
	host := parsed.Hostname()
	switch {
	case strings.Contains(host, "google.") || strings.Contains(host, "bing.") || strings.Contains(host, "yahoo.") || strings.Contains(host, "duckduckgo."):
		return "Organic Search"
	case strings.Contains(host, "facebook.") || strings.Contains(host, "twitter.") || strings.Contains(host, "x.com") || strings.Contains(host, "linkedin.") || strings.Contains(host, "instagram.") || strings.Contains(host, "youtube."):
		return "Social Media"
	default:
		return "External Link"
	}
}

func fallbackLocation(userAgent string) *middleware.LocationResult {
	ua := strings.ToLower(userAgent)
	out := &middleware.LocationResult{Country: "Unknown", Platform: "Unknown", Browser: "Unknown", DeviceType: "Desktop"}
	switch {
	case strings.Contains(ua, "android"):
		out.Platform, out.DeviceType = "Android", "Mobile"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		out.Platform, out.DeviceType = "iOS", "Mobile"
	case strings.Contains(ua, "windows"):
		out.Platform = "Windows"
	case strings.Contains(ua, "mac os"):
		out.Platform = "macOS"
	case strings.Contains(ua, "linux"):
		out.Platform = "Linux"
	}
	switch {
	case strings.Contains(ua, "edg/"):
		out.Browser = "Edge"
	case strings.Contains(ua, "chrome/"):
		out.Browser = "Chrome"
	case strings.Contains(ua, "firefox/"):
		out.Browser = "Firefox"
	case strings.Contains(ua, "safari/"):
		out.Browser = "Safari"
	}
	return out
}

func trim(value string, max int) string {
	value = strings.TrimSpace(value)
	if max > 0 && len(value) > max {
		return value[:max]
	}
	return value
}
