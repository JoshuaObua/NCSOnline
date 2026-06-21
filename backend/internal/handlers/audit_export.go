package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/repository"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
)

func parseAuditTime(v string) (*time.Time, error) {
	if v == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, v)
	return &t, err
}
func (h *AuditHandler) Export(w http.ResponseWriter, r *http.Request) {
	from, e1 := parseAuditTime(r.URL.Query().Get("from"))
	to, e2 := parseAuditTime(r.URL.Query().Get("to"))
	if e1 != nil || e2 != nil {
		response.Err(w, 400, "BAD_DATE", "from/to must use RFC3339")
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	filter := repository.AuditExportFilter{From: from, To: to, Severity: r.URL.Query().Get("severity"), UserID: r.URL.Query().Get("user_id"), IP: r.URL.Query().Get("ip"), EventType: r.URL.Query().Get("event_type"), FailedAuthOnly: r.URL.Query().Get("failed_auth_only") == "true", Limit: limit}
	items, err := h.repo.Export(r.Context(), filter)
	if err != nil {
		response.Err(w, 500, "EXPORT_ERROR", "Could not export audit logs")
		return
	}
	format := strings.ToLower(r.URL.Query().Get("format"))
	if format == "" {
		format = "csv"
	}
	name := "ncs-audit-" + time.Now().UTC().Format("20060102-150405") + "." + format
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	switch format {
	case "csv":
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		_ = writeAuditCSV(w, items)
	case "txt":
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		for _, item := range items {
			encoded, _ := json.Marshal(item)
			fmt.Fprintln(w, string(encoded))
		}
	case "md":
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		fmt.Fprint(w, "# NCS Audit Export\n\n")
		for _, item := range items {
			encoded, _ := json.MarshalIndent(item, "", "  ")
			fmt.Fprintf(w, "## %s — %s\n\n```json\n%s\n```\n\n", item.CreatedAt.UTC().Format(time.RFC3339), item.EventType, encoded)
		}
	case "xml":
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(xml.Header))
		_ = xml.NewEncoder(w).Encode(struct {
			XMLName xml.Name           `xml:"audit_logs"`
			Items   []*models.AuditLog `xml:"event"`
		}{Items: items})
	case "pdf":
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write(simpleAuditPDF(items))
	default:
		response.Err(w, 400, "BAD_FORMAT", "format must be csv, txt, md, xml or pdf")
	}
}
func writeAuditCSV(w http.ResponseWriter, items []*models.AuditLog) error {
	out := csv.NewWriter(w)
	defer out.Flush()
	_ = out.Write([]string{"timestamp", "event", "status", "severity", "user_id", "username", "method", "endpoint", "response_code", "ip", "country", "device", "entry_hash"})
	for _, i := range items {
		uid := ""
		if i.UserID != nil {
			uid = *i.UserID
		}
		_ = out.Write([]string{i.CreatedAt.UTC().Format(time.RFC3339), i.EventType, i.EventStatus, i.SeverityLevel, uid, i.Username, i.Method, i.Endpoint, strconv.Itoa(i.ResponseCode), i.IPAddress, i.GeoCountry, i.DeviceInfo, i.EntryHash})
	}
	return out.Error()
}
func pdfEscape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	return strings.ReplaceAll(s, ")", "\\)")
}
func simpleAuditPDF(items []*models.AuditLog) []byte {
	lines := []string{"NCS Online Audit Export", time.Now().UTC().Format(time.RFC3339)}
	for _, i := range items {
		lines = append(lines, fmt.Sprintf("%s | %s | %s | %s | %s", i.CreatedAt.UTC().Format("2006-01-02 15:04:05"), i.SeverityLevel, i.EventType, i.Username, i.IPAddress))
	}
	if len(lines) > 55 {
		lines = lines[:55]
		lines = append(lines, "Export truncated to first 53 rows; use CSV for the full dataset.")
	}
	var content bytes.Buffer
	content.WriteString("BT /F1 9 Tf 36 800 Td 12 TL ")
	for _, line := range lines {
		content.WriteString("(" + pdfEscape(line) + ") Tj T* ")
	}
	content.WriteString("ET")
	stream := content.String()
	objects := []string{"<< /Type /Catalog /Pages 2 0 R >>", "<< /Type /Pages /Kids [3 0 R] /Count 1 >>", "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 5 0 R >> >> /Contents 4 0 R >>", fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(stream), stream), "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>"}
	var out bytes.Buffer
	out.WriteString("%PDF-1.4\n")
	offsets := []int{0}
	for idx, obj := range objects {
		offsets = append(offsets, out.Len())
		fmt.Fprintf(&out, "%d 0 obj\n%s\nendobj\n", idx+1, obj)
	}
	xref := out.Len()
	fmt.Fprintf(&out, "xref\n0 %d\n0000000000 65535 f \n", len(objects)+1)
	for _, off := range offsets[1:] {
		fmt.Fprintf(&out, "%010d 00000 n \n", off)
	}
	fmt.Fprintf(&out, "trailer << /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF", len(objects)+1, xref)
	return out.Bytes()
}
