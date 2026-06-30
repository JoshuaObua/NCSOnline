package logging

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// DailyWriter rotates JSON logs at UTC midnight, compresses closed files and
// keeps a bounded number of days. Writes are serialized and never expose
// application secrets beyond what callers intentionally log.
type DailyWriter struct {
	mu        sync.Mutex
	dir       string
	prefix    string
	retention int
	day       string
	file      *os.File
}

func NewDailyWriter(dir, prefix string, retention int) (*DailyWriter, error) {
	if retention < 1 {
		retention = 30
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, err
	}
	w := &DailyWriter{dir: dir, prefix: prefix, retention: retention}
	if err := w.rotate(time.Now().UTC()); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *DailyWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UTC()
	if now.Format("20060102") != w.day {
		if err := w.rotate(now); err != nil {
			return 0, err
		}
	}
	return w.file.Write(p)
}

func (w *DailyWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func (w *DailyWriter) rotate(now time.Time) error {
	old := w.file
	oldPath := ""
	if old != nil {
		oldPath = old.Name()
		_ = old.Close()
	}
	w.day = now.Format("20060102")
	path := filepath.Join(w.dir, fmt.Sprintf("%s-%s.log", w.prefix, w.day))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return err
	}
	w.file = f
	if oldPath != "" && oldPath != path {
		go gzipFile(oldPath)
	}
	go w.prune()
	return nil
}

func gzipFile(path string) {
	in, err := os.Open(path)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.OpenFile(path+".gz", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o640)
	if err != nil {
		return
	}
	gz := gzip.NewWriter(out)
	_, copyErr := io.Copy(gz, in)
	closeErr := gz.Close()
	_ = out.Close()
	if copyErr == nil && closeErr == nil {
		_ = os.Remove(path)
	}
}

func (w *DailyWriter) prune() {
	entries, err := filepath.Glob(filepath.Join(w.dir, w.prefix+"-*.log*"))
	if err != nil || len(entries) <= w.retention {
		return
	}
	sort.Strings(entries)
	for _, path := range entries[:len(entries)-w.retention] {
		_ = os.Remove(path)
	}
}

func Configure(levelName, dir string) (io.Closer, error) {
	level := slog.LevelInfo
	switch strings.ToLower(levelName) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	writer, err := NewDailyWriter(dir, "ncs-online", 30)
	if err != nil {
		return nil, err
	}
	out := io.MultiWriter(os.Stdout, writer)
	slog.SetDefault(slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level})))
	log.SetOutput(out)
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)
	return writer, nil
}
