package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
)

type DocumentRecord struct {
	ID             string `json:"id"`
	FederationID   string `json:"federation_id"`
	ReportID       string `json:"report_id,omitempty"`
	DocumentType   string `json:"document_type"`
	OriginalName   string `json:"original_name"`
	StorageKey     string `json:"-"`
	MIMEType       string `json:"mime_type"`
	SHA256         string `json:"sha256"`
	Classification string `json:"classification"`
	ScanStatus     string `json:"scan_status"`
	SizeBytes      int64  `json:"size_bytes"`
}

func (r *NSMISRepo) CreateDocument(ctx context.Context, d *DocumentRecord, userID string, all bool) error {
	q := `INSERT INTO federation_documents(federation_id,report_id,document_type,original_name,storage_key,mime_type,size_bytes,sha256,classification,scan_status,uploaded_by)
	SELECT $1,NULLIF($2,''),$3,$4,$5,$6,$7,$8,$9,$10,$11
	WHERE $12 OR EXISTS(SELECT 1 FROM federation_memberships WHERE federation_id=$1 AND user_id=$11 AND is_active AND (ends_at IS NULL OR ends_at>=CURRENT_DATE)) RETURNING id`
	err := r.db.QueryRow(ctx, q, d.FederationID, d.ReportID, d.DocumentType, d.OriginalName, d.StorageKey, d.MIMEType, d.SizeBytes, d.SHA256, d.Classification, d.ScanStatus, userID, all).Scan(&d.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
func (r *NSMISRepo) ListDocuments(ctx context.Context, federationID, userID string, all bool) ([]map[string]interface{}, error) {
	rows, err := r.db.Query(ctx, `SELECT to_jsonb(d)-'storage_key' FROM federation_documents d WHERE d.federation_id=$1 AND ($3 OR EXISTS(SELECT 1 FROM federation_memberships fm WHERE fm.federation_id=d.federation_id AND fm.user_id=$2 AND fm.is_active AND (fm.ends_at IS NULL OR fm.ends_at>=CURRENT_DATE))) ORDER BY d.created_at DESC`, federationID, userID, all)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]interface{}{}
	for rows.Next() {
		var b []byte
		if err = rows.Scan(&b); err != nil {
			return nil, err
		}
		var x map[string]interface{}
		if err = json.Unmarshal(b, &x); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *NSMISRepo) GetDocument(ctx context.Context, id, userID string, all bool) (*DocumentRecord, error) {
	d := &DocumentRecord{}
	err := r.db.QueryRow(ctx, `SELECT d.id,d.federation_id,COALESCE(d.report_id,''),d.document_type,d.original_name,d.storage_key,d.mime_type,d.size_bytes,d.sha256,d.classification,d.scan_status FROM federation_documents d WHERE d.id=$1 AND ($3 OR EXISTS(SELECT 1 FROM federation_memberships fm WHERE fm.federation_id=d.federation_id AND fm.user_id=$2 AND fm.is_active AND (fm.ends_at IS NULL OR fm.ends_at>=CURRENT_DATE)))`, id, userID, all).Scan(&d.ID, &d.FederationID, &d.ReportID, &d.DocumentType, &d.OriginalName, &d.StorageKey, &d.MIMEType, &d.SizeBytes, &d.SHA256, &d.Classification, &d.ScanStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return d, err
}
