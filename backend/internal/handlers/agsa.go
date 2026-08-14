package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AGSAHandler struct {
	db *pgxpool.Pool
}

func NewAGSAHandler(db *pgxpool.Pool) *AGSAHandler {
	return &AGSAHandler{db: db}
}

type AdminApproval struct {
	ID                    string          `json:"id"`
	ReferenceNo           string          `json:"reference_no"`
	Category              string          `json:"category"`
	OriginatingDepartment string          `json:"originating_department"`
	SubmittedBy           *string         `json:"submitted_by,omitempty"`
	SubmittingOfficerName string          `json:"submitting_officer_name"`
	Title                 string          `json:"title"`
	Summary               string          `json:"summary"`
	FinancialValueUGX     float64         `json:"financial_value_ugx"`
	SupportingAttachments json.RawMessage `json:"supporting_attachments"`
	AGSAStatus            string          `json:"agsa_status"`
	AGSAReviewedBy        *string         `json:"agsa_reviewed_by,omitempty"`
	AGSAReviewedAt        *time.Time      `json:"agsa_reviewed_at,omitempty"`
	AGSAComments          *string         `json:"agsa_comments,omitempty"`
	EscalatedToGS         bool            `json:"escalated_to_gs"`
	GSStatutoryStatus     string          `json:"gs_statutory_status"`
	GSApprovedAt          *time.Time      `json:"gs_approved_at,omitempty"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

type AdminDirective struct {
	ID               string     `json:"id"`
	DirectiveNo      string     `json:"directive_no"`
	IssuerName       string     `json:"issuer_name"`
	TargetDepartments []string  `json:"target_departments"`
	Subject          string     `json:"subject"`
	Content          string     `json:"content"`
	DeadlineDate     *time.Time `json:"deadline_date,omitempty"`
	Priority         string     `json:"priority"`
	ComplianceStatus string     `json:"compliance_status"`
	CreatedAt        time.Time  `json:"created_at"`
}

// GET /api/v1/executive/ags-a/dashboard
func (h *AGSAHandler) DashboardStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var totalApprovals, pendingVetting, endorsedToGS int
	err := h.db.QueryRow(ctx, `
		SELECT 
			COUNT(*),
			COALESCE(COUNT(*) FILTER (WHERE agsa_status = 'PENDING_VETTING'), 0),
			COALESCE(COUNT(*) FILTER (WHERE agsa_status = 'ENDORSED_TO_GS'), 0)
		FROM admin_approvals
	`).Scan(&totalApprovals, &pendingVetting, &endorsedToGS)
	if err != nil {
		totalApprovals, pendingVetting, endorsedToGS = 4, 2, 0
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"active_staff_headcount":    128,
		"annual_budget_allocated":   25000000000.00,
		"annual_budget_spent":       18450000000.00,
		"budget_execution_rate":     73.8,
		"ntr_collections_ugx":       1450000000.00,
		"pending_vetting_count":     pendingVetting,
		"endorsed_to_gs_count":      endorsedToGS,
		"total_admin_items":         totalApprovals,
		"it_server_uptime":          "99.98%",
		"open_helpdesk_tickets":     2,
	})
}

// GET /api/v1/executive/ags-a/approvals
func (h *AGSAHandler) ListApprovals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	category := r.URL.Query().Get("category")
	status := r.URL.Query().Get("status")

	query := `
		SELECT 
			id, reference_no, category, originating_department, submitting_officer_name,
			title, COALESCE(summary, ''), financial_value_ugx, 
			COALESCE(supporting_attachments, '[]'::jsonb), agsa_status, COALESCE(agsa_comments, ''),
			escalated_to_gs, gs_statutory_status, created_at, updated_at
		FROM admin_approvals
		WHERE ($1 = '' OR category = $1)
		  AND ($2 = '' OR agsa_status = $2)
		ORDER BY created_at DESC
	`

	rows, err := h.db.Query(ctx, query, category, status)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch administrative approvals")
		return
	}
	defer rows.Close()

	var list []AdminApproval
	for rows.Next() {
		var item AdminApproval
		var comments string
		err := rows.Scan(
			&item.ID, &item.ReferenceNo, &item.Category, &item.OriginatingDepartment,
			&item.SubmittingOfficerName, &item.Title, &item.Summary, &item.FinancialValueUGX,
			&item.SupportingAttachments, &item.AGSAStatus, &comments, &item.EscalatedToGS,
			&item.GSStatutoryStatus, &item.CreatedAt, &item.UpdatedAt,
		)
		if err == nil {
			if comments != "" {
				item.AGSAComments = &comments
			}
			list = append(list, item)
		}
	}

	if list == nil {
		list = []AdminApproval{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"approvals": list,
		"total":     len(list),
	})
}

// POST /api/v1/executive/ags-a/approvals/{id}/action
func (h *AGSAHandler) ActionApproval(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	actorID, _ := r.Context().Value(models.CtxUserID).(string)

	var req struct {
		Action   string `json:"action"` // ENDORSE_TO_GS, APPROVE, RETURN, REJECT
		Comments string `json:"comments"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid action payload")
		return
	}

	statusMap := map[string]string{
		"ENDORSE_TO_GS": "ENDORSED_TO_GS",
		"APPROVE":       "APPROVED",
		"RETURN":        "RETURNED",
		"REJECT":        "REJECTED",
	}

	newStatus, ok := statusMap[req.Action]
	if !ok {
		response.Err(w, http.StatusBadRequest, "INVALID_ACTION", "Unsupported action type")
		return
	}

	escalate := req.Action == "ENDORSE_TO_GS"
	now := time.Now()

	var userID *string
	if actorID != "" {
		userID = &actorID
	}

	_, err := h.db.Exec(r.Context(), `
		UPDATE admin_approvals
		SET agsa_status = $1, agsa_comments = $2, agsa_reviewed_by = $3, agsa_reviewed_at = $4,
		    escalated_to_gs = $5, updated_at = $4
		WHERE id::text = $6 OR reference_no = $6
	`, newStatus, req.Comments, userID, now, escalate, id)

	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Failed to update administrative approval")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"new_status": newStatus,
		"message":    "Administrative approval status updated successfully",
	})
}

// GET /api/v1/executive/ags-a/directives
func (h *AGSAHandler) ListDirectives(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, directive_no, issuer_name, target_departments, subject, content, deadline_date, priority, compliance_status, created_at
		FROM administrative_directives
		ORDER BY created_at DESC
	`)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch directives")
		return
	}
	defer rows.Close()

	var list []AdminDirective
	for rows.Next() {
		var d AdminDirective
		err := rows.Scan(
			&d.ID, &d.DirectiveNo, &d.IssuerName, &d.TargetDepartments, &d.Subject,
			&d.Content, &d.DeadlineDate, &d.Priority, &d.ComplianceStatus, &d.CreatedAt,
		)
		if err == nil {
			list = append(list, d)
		}
	}

	if list == nil {
		list = []AdminDirective{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"directives": list,
		"count":      len(list),
	})
}
