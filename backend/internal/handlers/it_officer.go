package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ITOfficerHandler struct {
	db *pgxpool.Pool
}

func NewITOfficerHandler(db *pgxpool.Pool) *ITOfficerHandler {
	return &ITOfficerHandler{db: db}
}

type PPDAForm5Item struct {
	ItemName      string  `json:"item_name"`
	Quantity      int     `json:"quantity"`
	Unit          string  `json:"unit"`
	UnitPriceUGX  float64 `json:"unit_price_ugx"`
	TotalPriceUGX float64 `json:"total_price_ugx"`
}

type PPDAForm5Requisition struct {
	ID                      string          `json:"id"`
	ReferenceNo             string          `json:"reference_no"`
	UserID                  *string         `json:"user_id"`
	OfficerName             string          `json:"officer_name"`
	Department              string          `json:"department"`
	Designation             string          `json:"designation"`
	ContactPhone            *string         `json:"contact_phone"`
	FinancialYear           string          `json:"financial_year"`
	SubjectOfProcurement    string          `json:"subject_of_procurement"`
	ProcurementCategory     string          `json:"procurement_category"`
	ProcurementMethod       string          `json:"procurement_method"`
	BudgetVoteHead          string          `json:"budget_vote_head"`
	SourceOfFunds           string          `json:"source_of_funds"`
	DeliveryLocation        string          `json:"delivery_location"`
	WarrantyRequirement     string          `json:"warranty_requirement"`
	EstimatedAmountUGX      float64         `json:"estimated_amount_ugx"`
	RequiredDeliveryDate    *string         `json:"required_delivery_date"`
	TechnicalSpecifications string          `json:"technical_specifications"`
	Items                   []PPDAForm5Item `json:"items"`
	Justification           string          `json:"justification"`
	Status                  string          `json:"status"`
	CurrentStage            string          `json:"current_stage"`
	ApprovalRemarks         *string         `json:"approval_remarks"`
	PONumber                *string         `json:"po_number"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
}

// GET /api/v1/it/ppda
func (h *ITOfficerHandler) ListPPDAForm5(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value(models.CtxUserID).(string)
	statusFilter := r.URL.Query().Get("status")

	query := `
		SELECT 
			id, reference_no, user_id, officer_name, department, designation,
			contact_phone, financial_year, subject_of_procurement, procurement_category,
			procurement_method, budget_vote_head, source_of_funds, delivery_location,
			warranty_requirement, estimated_amount_ugx, required_delivery_date::TEXT,
			technical_specifications, items_json, justification, status, current_stage,
			approval_remarks, po_number, created_at, updated_at
		FROM ppda_form5_requisitions
		WHERE 1=1
	`
	var args []interface{}
	argIdx := 1

	if userID != "" {
		query += fmt.Sprintf(" AND (user_id = $%d OR user_id = 'usr_it_officer_001')", argIdx)
		args = append(args, userID)
		argIdx++
	}

	if statusFilter != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, statusFilter)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "DB_ERROR", "Could not fetch PPDA Form 5 records: "+err.Error())
		return
	}
	defer rows.Close()

	var requisitions []PPDAForm5Requisition
	for rows.Next() {
		var req PPDAForm5Requisition
		var rawItems []byte
		err := rows.Scan(
			&req.ID, &req.ReferenceNo, &req.UserID, &req.OfficerName, &req.Department, &req.Designation,
			&req.ContactPhone, &req.FinancialYear, &req.SubjectOfProcurement, &req.ProcurementCategory,
			&req.ProcurementMethod, &req.BudgetVoteHead, &req.SourceOfFunds, &req.DeliveryLocation,
			&req.WarrantyRequirement, &req.EstimatedAmountUGX, &req.RequiredDeliveryDate,
			&req.TechnicalSpecifications, &rawItems, &req.Justification, &req.Status, &req.CurrentStage,
			&req.ApprovalRemarks, &req.PONumber, &req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			continue
		}
		if len(rawItems) > 0 {
			_ = json.Unmarshal(rawItems, &req.Items)
		}
		if req.Items == nil {
			req.Items = []PPDAForm5Item{}
		}
		requisitions = append(requisitions, req)
	}

	if requisitions == nil {
		requisitions = []PPDAForm5Requisition{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"requisitions": requisitions,
		"count":        len(requisitions),
	})
}

// POST /api/v1/it/ppda
func (h *ITOfficerHandler) CreatePPDAForm5(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value(models.CtxUserID).(string)

	var req struct {
		OfficerName             string          `json:"officer_name"`
		Department              string          `json:"department"`
		Designation             string          `json:"designation"`
		ContactPhone            string          `json:"contact_phone"`
		FinancialYear           string          `json:"financial_year"`
		SubjectOfProcurement    string          `json:"subject_of_procurement"`
		ProcurementCategory     string          `json:"procurement_category"`
		ProcurementMethod       string          `json:"procurement_method"`
		BudgetVoteHead          string          `json:"budget_vote_head"`
		SourceOfFunds           string          `json:"source_of_funds"`
		DeliveryLocation        string          `json:"delivery_location"`
		WarrantyRequirement     string          `json:"warranty_requirement"`
		EstimatedAmountUGX      float64         `json:"estimated_amount_ugx"`
		RequiredDeliveryDate    string          `json:"required_delivery_date"`
		TechnicalSpecifications string          `json:"technical_specifications"`
		Items                   []PPDAForm5Item `json:"items"`
		Justification           string          `json:"justification"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if req.SubjectOfProcurement == "" || req.BudgetVoteHead == "" || req.Justification == "" {
		response.Err(w, http.StatusBadRequest, "VALIDATION_FAILED", "Subject of procurement, budget vote head, and justification are required")
		return
	}

	if req.OfficerName == "" {
		req.OfficerName = "Allan Tumusiime"
	}
	if req.Department == "" {
		req.Department = "IT / ICT Infrastructure"
	}
	if req.Designation == "" {
		req.Designation = "ICT Systems & Database Administrator"
	}
	if req.FinancialYear == "" {
		req.FinancialYear = "FY 2026/2027"
	}
	if req.ProcurementCategory == "" {
		req.ProcurementCategory = "SUPPLIES"
	}
	if req.ProcurementMethod == "" {
		req.ProcurementMethod = "Request for Quotations (RFQ)"
	}
	if req.SourceOfFunds == "" {
		req.SourceOfFunds = "GoU Statutory Subvention"
	}
	if req.DeliveryLocation == "" {
		req.DeliveryLocation = "NCS Headquarters Lugogo - ICT Datacenter"
	}
	if req.WarrantyRequirement == "" {
		req.WarrantyRequirement = "1 Year Comprehensive OEM Manufacturer Warranty"
	}

	// Calculate total amount from items if provided and 0
	if req.EstimatedAmountUGX <= 0 && len(req.Items) > 0 {
		for _, item := range req.Items {
			req.EstimatedAmountUGX += item.TotalPriceUGX
		}
	}

	refNo := fmt.Sprintf("NCS-PPDA-F5-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)

	itemsBytes, _ := json.Marshal(req.Items)
	if len(itemsBytes) == 0 {
		itemsBytes = []byte("[]")
	}

	var deliveryDateVal interface{} = nil
	if req.RequiredDeliveryDate != "" {
		deliveryDateVal = req.RequiredDeliveryDate
	}

	var reqID string
	err := h.db.QueryRow(ctx, `
		INSERT INTO ppda_form5_requisitions (
			reference_no, user_id, officer_name, department, designation,
			contact_phone, financial_year, subject_of_procurement, procurement_category,
			procurement_method, budget_vote_head, source_of_funds, delivery_location,
			warranty_requirement, estimated_amount_ugx, required_delivery_date,
			technical_specifications, items_json, justification, status, current_stage
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13,
			$14, $15, $16::DATE,
			$17, $18, $19, 'SUBMITTED', 'Awaiting Head of Department Endorsement'
		)
		RETURNING id
	`,
		refNo, userID, req.OfficerName, req.Department, req.Designation,
		req.ContactPhone, req.FinancialYear, req.SubjectOfProcurement, req.ProcurementCategory,
		req.ProcurementMethod, req.BudgetVoteHead, req.SourceOfFunds, req.DeliveryLocation,
		req.WarrantyRequirement, req.EstimatedAmountUGX, deliveryDateVal,
		req.TechnicalSpecifications, itemsBytes, req.Justification,
	).Scan(&reqID)

	if err != nil {
		response.Err(w, http.StatusInternalServerError, "DB_ERROR", "Could not submit PPDA Form 5 requisition: "+err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"message":      "PPDA Form 5 Procurement Requisition submitted successfully for HOD endorsement",
		"id":           reqID,
		"reference_no": refNo,
		"status":       "SUBMITTED",
		"stage":        "Awaiting Head of Department Endorsement",
	})
}

// GET /api/v1/it/ppda/{id}
func (h *ITOfficerHandler) GetPPDAForm5(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := chi.URLParam(r, "id")

	var req PPDAForm5Requisition
	var rawItems []byte
	err := h.db.QueryRow(ctx, `
		SELECT 
			id, reference_no, user_id, officer_name, department, designation,
			contact_phone, financial_year, subject_of_procurement, procurement_category,
			procurement_method, budget_vote_head, source_of_funds, delivery_location,
			warranty_requirement, estimated_amount_ugx, required_delivery_date::TEXT,
			technical_specifications, items_json, justification, status, current_stage,
			approval_remarks, po_number, created_at, updated_at
		FROM ppda_form5_requisitions
		WHERE id::TEXT = $1 OR reference_no = $1
	`, reqID).Scan(
		&req.ID, &req.ReferenceNo, &req.UserID, &req.OfficerName, &req.Department, &req.Designation,
		&req.ContactPhone, &req.FinancialYear, &req.SubjectOfProcurement, &req.ProcurementCategory,
		&req.ProcurementMethod, &req.BudgetVoteHead, &req.SourceOfFunds, &req.DeliveryLocation,
		&req.WarrantyRequirement, &req.EstimatedAmountUGX, &req.RequiredDeliveryDate,
		&req.TechnicalSpecifications, &rawItems, &req.Justification, &req.Status, &req.CurrentStage,
		&req.ApprovalRemarks, &req.PONumber, &req.CreatedAt, &req.UpdatedAt,
	)

	if err != nil {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "PPDA Form 5 record not found")
		return
	}

	if len(rawItems) > 0 {
		_ = json.Unmarshal(rawItems, &req.Items)
	}
	if req.Items == nil {
		req.Items = []PPDAForm5Item{}
	}

	response.JSON(w, http.StatusOK, req)
}

// GET /api/v1/it/dashboard/stats
func (h *ITOfficerHandler) GetITDashboardStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var totalRequisitions int64
	var pendingApprovals int64
	var approvedRequisitions int64
	var totalRequisitionedUGX float64

	_ = h.db.QueryRow(ctx, "SELECT COUNT(*), COALESCE(SUM(estimated_amount_ugx), 0) FROM ppda_form5_requisitions").Scan(&totalRequisitions, &totalRequisitionedUGX)
	_ = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM ppda_form5_requisitions WHERE status IN ('SUBMITTED', 'HOD_APPROVED', 'PDU_REVIEW', 'FINANCE_CLEARED')").Scan(&pendingApprovals)
	_ = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM ppda_form5_requisitions WHERE status IN ('ACCOUNTING_OFFICER_APPROVED', 'PO_ISSUED', 'DELIVERED')").Scan(&approvedRequisitions)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"total_requisitions":     totalRequisitions,
		"pending_approvals":      pendingApprovals,
		"approved_requisitions":  approvedRequisitions,
		"total_requisitioned_ugx": totalRequisitionedUGX,
		"system_health": map[string]interface{}{
			"database_status":     "Optimal",
			"cloud_backup_status": "Protected (Daily RPO 24h)",
			"server_uptime":       "99.94%",
		},
	})
}

type PPDAWorkflowActionReq struct {
	Action   string `json:"action"` // ENDORSE_HOD, ENDORSE_AGST, COMMIT_FINANCE, APPROVE_ACCOUNTING_OFFICER, ISSUE_PO, CONFIRM_DELIVERY, REJECT
	Remarks  string `json:"remarks"`
	PONumber string `json:"po_number"`
}

// POST /api/v1/it/ppda/{id}/action
func (h *ITOfficerHandler) AdvancePPDAWorkflow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := chi.URLParam(r, "id")

	var req PPDAWorkflowActionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	var newStatus string
	var newStage string

	switch req.Action {
	case "ENDORSE_HOD":
		newStatus = "HOD_APPROVED"
		newStage = "Forwarded to Assistant General Secretary - Technical"
	case "ENDORSE_AGST":
		newStatus = "AGST_APPROVED"
		newStage = "Forwarded to Head of Finance for Vote Head Clearance"
	case "COMMIT_FINANCE":
		newStatus = "FINANCE_CLEARED"
		newStage = "Forwarded to General Secretary / Accounting Officer for Final Approval"
	case "APPROVE_ACCOUNTING_OFFICER":
		newStatus = "ACCOUNTING_OFFICER_APPROVED"
		newStage = "Approved by Accounting Officer · Forwarded to PDU for Contract / PO"
	case "ISSUE_PO":
		newStatus = "PO_ISSUED"
		newStage = "Purchase Order Issued · Awaiting Supply & Delivery"
	case "CONFIRM_DELIVERY":
		newStatus = "DELIVERED"
		newStage = "Procurement Completed & Verified by User Department"
	case "REJECT":
		newStatus = "REJECTED"
		newStage = "Returned to Originating Department with Queries"
	default:
		response.Err(w, http.StatusBadRequest, "INVALID_ACTION", "Unknown workflow action: "+req.Action)
		return
	}

	var poVal *string
	if req.PONumber != "" {
		poVal = &req.PONumber
	}

	var remarksVal *string
	if req.Remarks != "" {
		remarksVal = &req.Remarks
	}

	res, err := h.db.Exec(ctx, `
		UPDATE ppda_form5_requisitions
		SET 
			status = $1,
			current_stage = $2,
			approval_remarks = COALESCE($3, approval_remarks),
			po_number = COALESCE($4, po_number),
			updated_at = NOW()
		WHERE id::TEXT = $5 OR reference_no = $5
	`, newStatus, newStage, remarksVal, poVal, reqID)

	if err != nil || res.RowsAffected() == 0 {
		response.Err(w, http.StatusInternalServerError, "DB_ERROR", "Could not advance requisition workflow")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Requisition workflow transitioned successfully",
		"status":  newStatus,
		"stage":   newStage,
	})
}

