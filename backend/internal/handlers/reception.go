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

type ReceptionHandler struct {
	db *pgxpool.Pool
}

func NewReceptionHandler(db *pgxpool.Pool) *ReceptionHandler {
	return &ReceptionHandler{db: db}
}

type VisitorPass struct {
	ID                  string     `json:"id"`
	PassNumber          string     `json:"pass_number"`
	VisitorName         string     `json:"visitor_name"`
	VisitorPhone        string     `json:"visitor_phone"`
	VisitorIDNumber     string     `json:"visitor_id_number"`
	VisitorOrganization string     `json:"visitor_organization"`
	TargetDepartment    string     `json:"target_department"`
	HostOfficerName     string     `json:"host_officer_name"`
	PurposeOfVisit      string     `json:"purpose_of_visit"`
	Status              string     `json:"status"`
	BadgeNumber         *string    `json:"badge_number"`
	ClearanceCode       *string    `json:"clearance_code"`
	ReceptionistUserID  *string    `json:"receptionist_user_id"`
	ApprovedByUserID    *string    `json:"approved_by_user_id"`
	ApprovedAt          *time.Time `json:"approved_at"`
	CheckedInAt         *time.Time `json:"checked_in_at"`
	CheckedOutAt        *time.Time `json:"checked_out_at"`
	VehicleRegNo        *string    `json:"vehicle_reg_no"`
	ItemsDeclared       *string    `json:"items_declared"`
	Remarks             *string    `json:"remarks"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type StaffLeaveApplication struct {
	ID                   string     `json:"id"`
	ReferenceNo          string     `json:"reference_no"`
	UserID               *string    `json:"user_id"`
	StaffName            string     `json:"staff_name"`
	StaffDepartment      string     `json:"staff_department"`
	StaffRole            string     `json:"staff_role"`
	StaffFileNo          string     `json:"staff_file_no"`
	ContactPhone         string     `json:"contact_phone"`
	ContactEmail         string     `json:"contact_email"`
	LeaveType            string     `json:"leave_type"`
	StartDate            string     `json:"start_date"`
	EndDate              string     `json:"end_date"`
	ReturnDate           string     `json:"return_date"`
	DaysRequested        int        `json:"days_requested"`
	RelievingOfficerName string     `json:"relieving_officer_name"`
	RelievingOfficerRole string     `json:"relieving_officer_role"`
	DutyHandoverDetails  string     `json:"duty_handover_details"`
	AddressWhileOnLeave  string     `json:"address_while_on_leave"`
	EmergencyPhone       string     `json:"emergency_phone"`
	Reason               string     `json:"reason"`
	Status               string     `json:"status"`
	ApprovedByUserID     *string    `json:"approved_by_user_id"`
	ApprovedAt           *time.Time `json:"approved_at"`
	SupervisorRemarks    *string    `json:"supervisor_remarks"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// GET /api/v1/reception/visitors
func (h *ReceptionHandler) ListVisitors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	statusFilter := r.URL.Query().Get("status")
	deptFilter := r.URL.Query().Get("department")

	query := `
		SELECT 
			id, pass_number, visitor_name, COALESCE(visitor_phone, ''), COALESCE(visitor_id_number, ''),
			COALESCE(visitor_organization, ''), target_department, host_officer_name,
			purpose_of_visit, status, badge_number, clearance_code,
			receptionist_user_id, approved_by_user_id, approved_at,
			checked_in_at, checked_out_at, vehicle_reg_no, items_declared,
			remarks, created_at, updated_at
		FROM visitor_passes
		WHERE 1=1
	`
	var args []interface{}
	argIdx := 1

	if statusFilter != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, statusFilter)
		argIdx++
	}
	if deptFilter != "" {
		query += fmt.Sprintf(" AND target_department = $%d", argIdx)
		args = append(args, deptFilter)
		argIdx++
	}

	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "DB_ERROR", "Could not fetch visitor records: "+err.Error())
		return
	}
	defer rows.Close()

	var visitors []VisitorPass
	for rows.Next() {
		var v VisitorPass
		err := rows.Scan(
			&v.ID, &v.PassNumber, &v.VisitorName, &v.VisitorPhone, &v.VisitorIDNumber,
			&v.VisitorOrganization, &v.TargetDepartment, &v.HostOfficerName,
			&v.PurposeOfVisit, &v.Status, &v.BadgeNumber, &v.ClearanceCode,
			&v.ReceptionistUserID, &v.ApprovedByUserID, &v.ApprovedAt,
			&v.CheckedInAt, &v.CheckedOutAt, &v.VehicleRegNo, &v.ItemsDeclared,
			&v.Remarks, &v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			continue
		}
		visitors = append(visitors, v)
	}

	if visitors == nil {
		visitors = []VisitorPass{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"visitors": visitors,
		"count":    len(visitors),
	})
}

// POST /api/v1/reception/visitors
func (h *ReceptionHandler) CreateVisitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value(models.CtxUserID).(string)

	var req struct {
		VisitorName         string `json:"visitor_name"`
		VisitorPhone        string `json:"visitor_phone"`
		VisitorIDNumber     string `json:"visitor_id_number"`
		VisitorOrganization string `json:"visitor_organization"`
		TargetDepartment    string `json:"target_department"`
		HostOfficerName     string `json:"host_officer_name"`
		PurposeOfVisit      string `json:"purpose_of_visit"`
		VehicleRegNo        string `json:"vehicle_reg_no"`
		ItemsDeclared       string `json:"items_declared"`
		Remarks             string `json:"remarks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if req.VisitorName == "" || req.TargetDepartment == "" || req.HostOfficerName == "" || req.PurposeOfVisit == "" {
		response.Err(w, http.StatusBadRequest, "VALIDATION_FAILED", "Visitor name, department, host officer, and purpose are required")
		return
	}

	passNumber := fmt.Sprintf("NCS-VP-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)

	var passID string
	err := h.db.QueryRow(ctx, `
		INSERT INTO visitor_passes (
			pass_number, visitor_name, visitor_phone, visitor_id_number,
			visitor_organization, target_department, host_officer_name,
			purpose_of_visit, status, receptionist_user_id,
			vehicle_reg_no, items_declared, remarks
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'PENDING_APPROVAL', $9, $10, $11, $12)
		RETURNING id
	`,
		passNumber, req.VisitorName, req.VisitorPhone, req.VisitorIDNumber,
		req.VisitorOrganization, req.TargetDepartment, req.HostOfficerName,
		req.PurposeOfVisit, userID, req.VehicleRegNo, req.ItemsDeclared, req.Remarks,
	).Scan(&passID)

	if err != nil {
		response.Err(w, http.StatusInternalServerError, "DB_ERROR", "Could not create visitor pass: "+err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"message":     "Visitor entry clearance request initiated successfully",
		"id":          passID,
		"pass_number": passNumber,
		"status":      "PENDING_APPROVAL",
	})
}

// PUT /api/v1/reception/visitors/{id}/approve
func (h *ReceptionHandler) ApproveVisitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value(models.CtxUserID).(string)
	passID := chi.URLParam(r, "id")

	clearanceCode := fmt.Sprintf("CLR-%06d", time.Now().UnixNano()%1000000)
	now := time.Now()

	tag, err := h.db.Exec(ctx, `
		UPDATE visitor_passes
		SET status = 'APPROVED',
		    clearance_code = $1,
		    approved_by_user_id = $2,
		    approved_at = $3,
		    updated_at = $4
		WHERE id = $5
	`, clearanceCode, userID, now, now, passID)

	if err != nil {
		response.Err(w, http.StatusInternalServerError, "DB_ERROR", "Could not approve visitor pass: "+err.Error())
		return
	}

	if tag.RowsAffected() == 0 {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Visitor pass record not found")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message":        "Visitor clearance approved. Clearance form ready for printing.",
		"clearance_code": clearanceCode,
		"status":         "APPROVED",
	})
}

// PUT /api/v1/reception/visitors/{id}/checkin
func (h *ReceptionHandler) CheckInVisitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	passID := chi.URLParam(r, "id")

	var req struct {
		BadgeNumber string `json:"badge_number"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	badge := req.BadgeNumber
	if badge == "" {
		badge = fmt.Sprintf("V-%03d", time.Now().Unix()%1000)
	}

	now := time.Now()
	tag, err := h.db.Exec(ctx, `
		UPDATE visitor_passes
		SET status = 'CHECKED_IN',
		    badge_number = $1,
		    checked_in_at = $2,
		    updated_at = $3
		WHERE id = $4
	`, badge, now, now, passID)

	if err != nil {
		response.Err(w, http.StatusInternalServerError, "DB_ERROR", "Could not check-in visitor: "+err.Error())
		return
	}

	if tag.RowsAffected() == 0 {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Visitor pass record not found")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message":      "Visitor checked in successfully",
		"badge_number": badge,
		"status":       "CHECKED_IN",
	})
}

// PUT /api/v1/reception/visitors/{id}/checkout
func (h *ReceptionHandler) CheckOutVisitor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	passID := chi.URLParam(r, "id")
	now := time.Now()

	tag, err := h.db.Exec(ctx, `
		UPDATE visitor_passes
		SET status = 'COMPLETED',
		    checked_out_at = $1,
		    updated_at = $2
		WHERE id = $3
	`, now, now, passID)

	if err != nil {
		response.Err(w, http.StatusInternalServerError, "DB_ERROR", "Could not check-out visitor: "+err.Error())
		return
	}

	if tag.RowsAffected() == 0 {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Visitor pass record not found")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "Visitor check-out recorded. Mission completed.",
		"status":  "COMPLETED",
	})
}

// GET /api/v1/reception/visitors/{id}
func (h *ReceptionHandler) GetVisitorPass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	passID := chi.URLParam(r, "id")

	var v VisitorPass
	err := h.db.QueryRow(ctx, `
		SELECT 
			id, pass_number, visitor_name, COALESCE(visitor_phone, ''), COALESCE(visitor_id_number, ''),
			COALESCE(visitor_organization, ''), target_department, host_officer_name,
			purpose_of_visit, status, badge_number, clearance_code,
			receptionist_user_id, approved_by_user_id, approved_at,
			checked_in_at, checked_out_at, vehicle_reg_no, items_declared,
			remarks, created_at, updated_at
		FROM visitor_passes
		WHERE id = $1
	`, passID).Scan(
		&v.ID, &v.PassNumber, &v.VisitorName, &v.VisitorPhone, &v.VisitorIDNumber,
		&v.VisitorOrganization, &v.TargetDepartment, &v.HostOfficerName,
		&v.PurposeOfVisit, &v.Status, &v.BadgeNumber, &v.ClearanceCode,
		&v.ReceptionistUserID, &v.ApprovedByUserID, &v.ApprovedAt,
		&v.CheckedInAt, &v.CheckedOutAt, &v.VehicleRegNo, &v.ItemsDeclared,
		&v.Remarks, &v.CreatedAt, &v.UpdatedAt,
	)

	if err != nil {
		response.Err(w, http.StatusNotFound, "NOT_FOUND", "Visitor pass not found")
		return
	}

	response.JSON(w, http.StatusOK, v)
}

// POST /api/v1/reception/leave/apply
func (h *ReceptionHandler) ApplyLeave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value(models.CtxUserID).(string)

	var req struct {
		StaffName            string `json:"staff_name"`
		StaffDepartment      string `json:"staff_department"`
		StaffRole            string `json:"staff_role"`
		StaffFileNo          string `json:"staff_file_no"`
		ContactPhone         string `json:"contact_phone"`
		ContactEmail         string `json:"contact_email"`
		LeaveType            string `json:"leave_type"`
		StartDate            string `json:"start_date"`
		EndDate              string `json:"end_date"`
		ReturnDate           string `json:"return_date"`
		DaysRequested        int    `json:"days_requested"`
		RelievingOfficerName string `json:"relieving_officer_name"`
		RelievingOfficerRole string `json:"relieving_officer_role"`
		DutyHandoverDetails  string `json:"duty_handover_details"`
		AddressWhileOnLeave  string `json:"address_while_on_leave"`
		EmergencyPhone       string `json:"emergency_phone"`
		Reason               string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if req.LeaveType == "" || req.StartDate == "" || req.EndDate == "" || req.DaysRequested <= 0 {
		response.Err(w, http.StatusBadRequest, "VALIDATION_FAILED", "Leave type, start date, end date, and valid days are required")
		return
	}

	if req.StaffDepartment == "" {
		req.StaffDepartment = "Front Desk / Reception"
	}
	if req.StaffRole == "" {
		req.StaffRole = "Receptionist"
	}

	refNo := fmt.Sprintf("NCS-LEV-%d-%04d", time.Now().Year(), time.Now().Unix()%10000)

	var returnDateVal interface{} = nil
	if req.ReturnDate != "" {
		returnDateVal = req.ReturnDate
	}

	var leaveID string
	err := h.db.QueryRow(ctx, `
		INSERT INTO staff_leave_applications (
			reference_no, user_id, staff_name, staff_department, staff_role,
			staff_file_no, contact_phone, contact_email,
			leave_type, start_date, end_date, return_date, days_requested,
			relieving_officer_name, relieving_officer_role, duty_handover_details,
			address_while_on_leave, emergency_phone,
			reason, status
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10::DATE, $11::DATE, $12::DATE, $13,
			$14, $15, $16,
			$17, $18,
			$19, 'PENDING_SUPERVISOR'
		)
		RETURNING id
	`,
		refNo, userID, req.StaffName, req.StaffDepartment, req.StaffRole,
		req.StaffFileNo, req.ContactPhone, req.ContactEmail,
		req.LeaveType, req.StartDate, req.EndDate, returnDateVal, req.DaysRequested,
		req.RelievingOfficerName, req.RelievingOfficerRole, req.DutyHandoverDetails,
		req.AddressWhileOnLeave, req.EmergencyPhone,
		req.Reason,
	).Scan(&leaveID)

	if err != nil {
		response.Err(w, http.StatusInternalServerError, "DB_ERROR", "Could not submit leave application: "+err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"message":      "Staff leave application submitted successfully to HR Supervisor",
		"id":           leaveID,
		"reference_no": refNo,
		"status":       "PENDING_SUPERVISOR",
	})
}

// GET /api/v1/reception/leave/status
func (h *ReceptionHandler) ListLeaveApplications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value(models.CtxUserID).(string)

	query := `
		SELECT 
			id, reference_no, user_id, staff_name, staff_department, staff_role,
			COALESCE(staff_file_no, ''), COALESCE(contact_phone, ''), COALESCE(contact_email, ''),
			leave_type, start_date::TEXT, end_date::TEXT, COALESCE(return_date::TEXT, ''), days_requested,
			relieving_officer_name, COALESCE(relieving_officer_role, ''), COALESCE(duty_handover_details, ''),
			COALESCE(address_while_on_leave, ''), COALESCE(emergency_phone, ''),
			reason, status, approved_by_user_id,
			approved_at, supervisor_remarks, created_at, updated_at
		FROM staff_leave_applications
		WHERE user_id = $1 OR $1 = ''
		ORDER BY created_at DESC
	`

	rows, err := h.db.Query(ctx, query, userID)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "DB_ERROR", "Could not fetch leave records: "+err.Error())
		return
	}
	defer rows.Close()

	var leaves []StaffLeaveApplication
	for rows.Next() {
		var l StaffLeaveApplication
		err := rows.Scan(
			&l.ID, &l.ReferenceNo, &l.UserID, &l.StaffName, &l.StaffDepartment, &l.StaffRole,
			&l.StaffFileNo, &l.ContactPhone, &l.ContactEmail,
			&l.LeaveType, &l.StartDate, &l.EndDate, &l.ReturnDate, &l.DaysRequested,
			&l.RelievingOfficerName, &l.RelievingOfficerRole, &l.DutyHandoverDetails,
			&l.AddressWhileOnLeave, &l.EmergencyPhone,
			&l.Reason, &l.Status, &l.ApprovedByUserID,
			&l.ApprovedAt, &l.SupervisorRemarks, &l.CreatedAt, &l.UpdatedAt,
		)
		if err != nil {
			continue
		}
		leaves = append(leaves, l)
	}

	if leaves == nil {
		leaves = []StaffLeaveApplication{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"leaves": leaves,
		"count":  len(leaves),
	})
}

// GET /api/v1/reception/reports
func (h *ReceptionHandler) GetReceptionReports(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var totalVisitors int64
	var approvedPasses int64
	var completedVisits int64
	var pendingPasses int64

	_ = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM visitor_passes").Scan(&totalVisitors)
	_ = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM visitor_passes WHERE status = 'APPROVED'").Scan(&approvedPasses)
	_ = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM visitor_passes WHERE status = 'COMPLETED'").Scan(&completedVisits)
	_ = h.db.QueryRow(ctx, "SELECT COUNT(*) FROM visitor_passes WHERE status = 'PENDING_APPROVAL'").Scan(&pendingPasses)

	// Departmental Breakdown
	rows, err := h.db.Query(ctx, `
		SELECT target_department, COUNT(*) 
		FROM visitor_passes 
		GROUP BY target_department 
		ORDER BY COUNT(*) DESC
	`)
	var deptTraffic []map[string]interface{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var dept string
			var count int64
			if rows.Scan(&dept, &count) == nil {
				deptTraffic = append(deptTraffic, map[string]interface{}{
					"department": dept,
					"visits":     count,
				})
			}
		}
	}
	if deptTraffic == nil {
		deptTraffic = []map[string]interface{}{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"summary": map[string]interface{}{
			"total_visitors":   totalVisitors,
			"approved_passes":  approvedPasses,
			"completed_visits": completedVisits,
			"pending_passes":   pendingPasses,
		},
		"departmental_traffic": deptTraffic,
		"peak_traffic_hours": []map[string]interface{}{
			{"hour": "08:00 - 10:00 AM", "label": "Morning Intake", "volume": "High"},
			{"hour": "10:00 - 01:00 PM", "label": "Official Appointments", "volume": "Peak"},
			{"hour": "02:00 - 04:30 PM", "label": "Afternoon Consultations", "volume": "Moderate"},
			{"hour": "04:30 - 05:30 PM", "label": "Exit & Debrief", "volume": "Low"},
		},
	})
}
