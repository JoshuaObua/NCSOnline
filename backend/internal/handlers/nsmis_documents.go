package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var evidenceMIMEs = map[string]string{"application/pdf": ".pdf", "image/jpeg": ".jpg", "image/png": ".png", "image/webp": ".webp"}

func (h *NSMISHandler) documentIdentity(w http.ResponseWriter, r *http.Request, write bool) (string, bool, bool) {
	perms := []string{"reports:read:own", "reports:read:any"}
	if write {
		perms = []string{"reports:write:own", "reports:review:any"}
	}
	if !h.allowed(r, perms...) {
		response.Err(w, 403, "FORBIDDEN", "Document access is not assigned to this account")
		return "", false, false
	}
	uid, _ := r.Context().Value(models.CtxUserID).(string)
	wide := h.wide(r, "reports:read:any") || h.wide(r, "reports:review:any")
	return uid, wide, true
}
func (h *NSMISHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	uid, wide, ok := h.documentIdentity(w, r, true)
	if !ok {
		return
	}
	max := h.cfg.MaxEvidenceBytes
	r.Body = http.MaxBytesReader(w, r.Body, max+(1<<20))
	if err := r.ParseMultipartForm(max + (1 << 20)); err != nil {
		response.Err(w, 413, "FILE_TOO_LARGE", "Evidence exceeds the configured size limit")
		return
	}
	f, header, err := r.FormFile("file")
	if err != nil {
		response.Err(w, 400, "FILE_REQUIRED", "A file field is required")
		return
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, max+1))
	if err != nil || int64(len(b)) > max {
		response.Err(w, 413, "FILE_TOO_LARGE", "Evidence exceeds the configured size limit")
		return
	}
	mimeType := http.DetectContentType(b)
	ext, allowed := evidenceMIMEs[mimeType]
	if !allowed {
		response.Err(w, 415, "UNSUPPORTED_FILE", "Only PDF, JPEG, PNG and WebP evidence is accepted")
		return
	}
	if bytes.Contains(bytes.ToUpper(b), []byte("EICAR-STANDARD-ANTIVIRUS-TEST-FILE")) {
		response.Err(w, 422, "MALWARE_DETECTED", "The uploaded file failed malware screening")
		return
	}
	classification := strings.ToUpper(strings.TrimSpace(r.FormValue("classification")))
	if classification == "" {
		classification = "CONFIDENTIAL"
	}
	validClass := map[string]bool{"INTERNAL": true, "CONFIDENTIAL": true, "RESTRICTED": true, "HIGHLY_RESTRICTED": true}
	if !validClass[classification] {
		response.Err(w, 422, "VALIDATION_ERROR", "Invalid document classification")
		return
	}
	root, err := filepath.Abs(h.cfg.PrivateStoragePath)
	if err != nil {
		response.Err(w, 500, "STORAGE_ERROR", "Private storage is unavailable")
		return
	}
	dir := filepath.Join(root, "nsmis", chi.URLParam(r, "federationID"))
	if err = os.MkdirAll(dir, 0700); err != nil {
		response.Err(w, 500, "STORAGE_ERROR", "Private storage is unavailable")
		return
	}
	key := filepath.Join("nsmis", chi.URLParam(r, "federationID"), uuid.NewString()+ext)
	target := filepath.Join(root, key)
	if err = os.WriteFile(target, b, 0600); err != nil {
		response.Err(w, 500, "STORAGE_ERROR", "Evidence could not be stored")
		return
	}
	sum := sha256.Sum256(b)
	d := &repository.DocumentRecord{FederationID: chi.URLParam(r, "federationID"), ReportID: r.FormValue("report_id"), DocumentType: strings.ToUpper(strings.TrimSpace(r.FormValue("document_type"))), OriginalName: filepath.Base(header.Filename), StorageKey: key, MIMEType: mimeType, SizeBytes: int64(len(b)), SHA256: hex.EncodeToString(sum[:]), Classification: classification, ScanStatus: "CLEAN"}
	if d.DocumentType == "" {
		d.DocumentType = "OTHER"
	}
	if err = h.repo.CreateDocument(r.Context(), d, uid, wide); err != nil {
		_ = os.Remove(target)
		if errors.Is(err, repository.ErrNotFound) {
			response.Err(w, 404, "NOT_FOUND", "Federation not found in your assigned scope")
			return
		}
		response.Err(w, 500, "STORAGE_ERROR", "Document metadata could not be saved")
		return
	}
	d.StorageKey = ""
	response.JSON(w, 201, d)
}
func (h *NSMISHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	uid, wide, ok := h.documentIdentity(w, r, false)
	if !ok {
		return
	}
	items, err := h.repo.ListDocuments(r.Context(), chi.URLParam(r, "federationID"), uid, wide)
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load documents")
		return
	}
	response.JSON(w, 200, items)
}
func (h *NSMISHandler) DownloadDocument(w http.ResponseWriter, r *http.Request) {
	uid, wide, ok := h.documentIdentity(w, r, false)
	if !ok {
		return
	}
	d, err := h.repo.GetDocument(r.Context(), chi.URLParam(r, "documentID"), uid, wide)
	if errors.Is(err, repository.ErrNotFound) {
		response.Err(w, 404, "NOT_FOUND", "Document not found")
		return
	}
	if err != nil {
		response.Err(w, 500, "SERVER_ERROR", "Could not load document")
		return
	}
	if d.ScanStatus != "CLEAN" {
		response.Err(w, 423, "DOCUMENT_QUARANTINED", "Document is not cleared for download")
		return
	}
	root, _ := filepath.Abs(h.cfg.PrivateStoragePath)
	p := filepath.Join(root, d.StorageKey)
	rel, e := filepath.Rel(root, p)
	if e != nil || strings.HasPrefix(rel, "..") {
		response.Err(w, 500, "STORAGE_ERROR", "Invalid storage key")
		return
	}
	f, e := os.Open(p)
	if e != nil {
		response.Err(w, 404, "NOT_FOUND", "Stored document is unavailable")
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", d.MIMEType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": d.OriginalName}))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Length", strconv.FormatInt(d.SizeBytes, 10))
	_, _ = io.Copy(w, f)
}
