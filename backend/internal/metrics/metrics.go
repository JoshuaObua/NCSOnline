package metrics

import (
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type routeMetric struct {
	count   uint64
	sum     float64
	buckets [8]uint64
}

var bounds = []float64{.005, .01, .025, .05, .1, .25, .5, 1}

type Registry struct {
	mu     sync.RWMutex
	routes map[string]*routeMetric
}

func New() *Registry { return &Registry{routes: map[string]*routeMetric{}} }

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }
func (w *statusWriter) WriteHeader(code int) {
	if w.status == 0 {
		w.status = code
	}
	w.ResponseWriter.WriteHeader(code)
}
func (w *statusWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(p)
}

func (m *Registry) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)
		pattern := chi.RouteContext(r.Context()).RoutePattern()
		if pattern == "" {
			pattern = r.URL.Path
		}
		key := r.Method + "|" + pattern + "|" + fmt.Sprint(sw.status)
		d := time.Since(start).Seconds()
		m.mu.Lock()
		metric := m.routes[key]
		if metric == nil {
			metric = &routeMetric{}
			m.routes[key] = metric
		}
		metric.count++
		metric.sum += d
		for i, b := range bounds {
			if d <= b {
				metric.buckets[i]++
			}
		}
		m.mu.Unlock()
	})
}

func esc(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "\"", "\\\"")
}
func (m *Registry) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	fmt.Fprintln(w, "# HELP ncs_http_requests_total Total HTTP requests.")
	fmt.Fprintln(w, "# TYPE ncs_http_requests_total counter")
	m.mu.RLock()
	keys := make([]string, 0, len(m.routes))
	for k := range m.routes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Fprintln(w, "# HELP ncs_http_request_duration_seconds Request latency.")
	fmt.Fprintln(w, "# TYPE ncs_http_request_duration_seconds histogram")
	for _, k := range keys {
		parts := strings.Split(k, "|")
		v := m.routes[k]
		labels := fmt.Sprintf("method=\"%s\",route=\"%s\",status=\"%s\"", esc(parts[0]), esc(parts[1]), esc(parts[2]))
		fmt.Fprintf(w, "ncs_http_requests_total{%s} %d\n", labels, v.count)
		fmt.Fprintf(w, "ncs_http_request_duration_seconds_sum{%s} %f\n", labels, v.sum)
		fmt.Fprintf(w, "ncs_http_request_duration_seconds_count{%s} %d\n", labels, v.count)
		for i, b := range bounds {
			fmt.Fprintf(w, "ncs_http_request_duration_seconds_bucket{%s,le=\"%g\"} %d\n", labels, b, v.buckets[i])
		}
		fmt.Fprintf(w, "ncs_http_request_duration_seconds_bucket{%s,le=\"+Inf\"} %d\n", labels, v.count)
	}
	m.mu.RUnlock()
	fmt.Fprintln(w, "# HELP go_goroutines Number of goroutines.")
	fmt.Fprintln(w, "# TYPE go_goroutines gauge")
	fmt.Fprintf(w, "go_goroutines %d\n", runtime.NumGoroutine())
}
