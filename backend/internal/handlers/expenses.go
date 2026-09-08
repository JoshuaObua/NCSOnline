package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/middleware"
	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExpenseHandler struct{ db *pgxpool.Pool }

func NewExpenseHandler(db *pgxpool.Pool) *ExpenseHandler { return &ExpenseHandler{db: db} }
func (h *ExpenseHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequireRoles("super_admin", "admin", "accountant", "senior_accountant", "chief_accountant", "asset_accountant", "finance_department"))
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/options", h.Options)
	r.Get("/categories", h.Categories)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireRoles("super_admin", "admin"))
		r.Post("/categories", h.SaveCategory)
		r.Put("/categories/{id}", h.SaveCategory)
	})
	r.Get("/{id}", h.Detail)
	r.Delete("/{id}", h.Delete)
	r.Put("/{id}/verify", h.Verify)
	r.Put("/{id}/approve", h.Approve)
	r.Get("/{id}/attachment", h.Attachment)
	return r
}
func expenseAdmin(r *http.Request) bool {
	roles, _ := r.Context().Value(models.CtxUserRoles).([]string)
	for _, role := range roles {
		if role == "admin" || role == "super_admin" {
			return true
		}
	}
	return false
}
func expenseError(w http.ResponseWriter, status int, message string) {
	response.Err(w, status, "EXPENSE_ERROR", message)
}
func expenseToday() string { return time.Now().In(time.FixedZone("EAT", 3*60*60)).Format("2006-01-02") }
func (h *ExpenseHandler) Options(w http.ResponseWriter, r *http.Request) {
	var departments json.RawMessage
	var name string
	err := h.db.QueryRow(r.Context(), `SELECT COALESCE(json_agg(d),'[]'::json) FROM (SELECT id,name FROM departments WHERE is_active ORDER BY name) d`).Scan(&departments)
	if err == nil {
		err = h.db.QueryRow(r.Context(), `SELECT trim(first_name || ' ' || last_name) FROM users WHERE id=$1`, currentUserID(r)).Scan(&name)
	}
	if err != nil {
		expenseError(w, 500, "Could not load expense options")
		return
	}
	response.JSON(w, 200, map[string]interface{}{"departments": departments, "recorded_by_name": name, "today": expenseToday()})
}
func (h *ExpenseHandler) Categories(w http.ResponseWriter, r *http.Request) {
	var result json.RawMessage
	err := h.db.QueryRow(r.Context(), `SELECT COALESCE(json_agg(c),'[]'::json) FROM (SELECT id,name,description,is_active FROM expense_categories WHERE $1 OR is_active ORDER BY name) c`, expenseAdmin(r)).Scan(&result)
	if err != nil {
		expenseError(w, 500, "Could not load categories")
		return
	}
	response.JSON(w, 200, result)
}
func (h *ExpenseHandler) SaveCategory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Active      *bool  `json:"is_active"`
	}
	if json.NewDecoder(io.LimitReader(r.Body, 8192)).Decode(&req) != nil {
		expenseError(w, 400, "Invalid category")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	if req.Name == "" || len(req.Name) > 120 || len(req.Description) > 1000 {
		expenseError(w, 400, "Category name is required (120 characters maximum); description must be at most 1000 characters")
		return
	}
	active := true
	if req.Active != nil {
		active = *req.Active
	}
	var result json.RawMessage
	uid := currentUserID(r)
	id := chi.URLParam(r, "id")
	var err error
	if id == "" {
		err = h.db.QueryRow(r.Context(), `INSERT INTO expense_categories(name,description,is_active,created_by,updated_by) VALUES($1,$2,$3,$4,$4) RETURNING row_to_json(expense_categories)`, req.Name, req.Description, active, uid).Scan(&result)
	} else {
		err = h.db.QueryRow(r.Context(), `UPDATE expense_categories SET name=$1,description=$2,is_active=$3,updated_by=$4,updated_at=NOW() WHERE id::text=$5 RETURNING row_to_json(expense_categories)`, req.Name, req.Description, active, uid, id).Scan(&result)
	}
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) && pgerr.Code == "23505" {
		expenseError(w, 409, "A category with this name already exists")
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		expenseError(w, 404, "Category not found")
		return
	}
	if err != nil {
		expenseError(w, 500, "Could not save category")
		return
	}
	status := 200
	if id == "" {
		status = 201
	}
	response.JSON(w, status, result)
}

var expenseAmountPattern = regexp.MustCompile(`^(0|[1-9][0-9]{0,11})(\.[0-9]{1,2})?$`)

func validExpenseAmount(s string) bool {
	return expenseAmountPattern.MatchString(s) && strings.Trim(s, "0.") != ""
}
func validExpenseDate(s string) bool { _, err := time.Parse("2006-01-02", s); return err == nil }
func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 6<<20)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		expenseError(w, 400, "Send an expense form with a receipt no larger than 5 MB")
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	date := r.FormValue("expense_date")
	if date == "" {
		date = expenseToday()
	}
	title := strings.TrimSpace(r.FormValue("title"))
	reason := strings.TrimSpace(r.FormValue("description"))
	amount := strings.TrimSpace(r.FormValue("amount"))
	payee := strings.TrimSpace(r.FormValue("payee"))
	method := r.FormValue("payment_method")
	reference := strings.TrimSpace(r.FormValue("payment_reference"))
	methods := map[string]bool{"CASH": true, "BANK_TRANSFER": true, "MOBILE_MONEY": true, "CARD": true, "OTHER": true}
	if !validExpenseDate(date) || date > expenseToday() || title == "" || len(title) > 200 || reason == "" || len(reason) > 5000 || !validExpenseAmount(amount) || payee == "" || len(payee) > 200 || !methods[method] || len(reference) > 200 {
		expenseError(w, 400, "Enter a valid date (not in the future), title, reason, positive UGX amount (up to 2 decimal places), payee and payment method")
		return
	}
	var content []byte
	var filename, mimeType string
	if len(r.MultipartForm.File["file"]) > 1 {
		expenseError(w, 400, "Attach one receipt per expense")
		return
	}
	file, header, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		content, err = io.ReadAll(io.LimitReader(file, (5<<20)+1))
		if err != nil || len(content) == 0 || len(content) > 5<<20 {
			expenseError(w, 400, "Receipt must be between 1 byte and 5 MB")
			return
		}
		mimeType = http.DetectContentType(content)
		if _, ok := evidenceMIMEs[mimeType]; !ok {
			expenseError(w, 415, "Receipt must be PDF, JPEG, PNG or WebP")
			return
		}
		filename = filepath.Base(strings.ReplaceAll(header.Filename, "\\", "/"))
		if len(filename) > 255 {
			expenseError(w, 400, "Receipt filename is too long")
			return
		}
	} else if !errors.Is(err, http.ErrMissingFile) {
		expenseError(w, 400, "Could not read receipt")
		return
	}
	tx, err := h.db.Begin(r.Context())
	if err != nil {
		expenseError(w, 500, "Could not record expense")
		return
	}
	defer tx.Rollback(r.Context())
	var category, dept, name, categoryID string
	err = tx.QueryRow(r.Context(), `SELECT id::text,name FROM expense_categories WHERE id::text=$1 AND is_active FOR SHARE`, r.FormValue("category_id")).Scan(&categoryID, &category)
	if errors.Is(err, pgx.ErrNoRows) {
		expenseError(w, 400, "Choose an active expense category")
		return
	}
	if err == nil {
		err = tx.QueryRow(r.Context(), `SELECT name FROM departments WHERE id=$1 AND is_active FOR SHARE`, r.FormValue("department_id")).Scan(&dept)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		expenseError(w, 400, "Choose an active department")
		return
	}
	uid := currentUserID(r)
	if err == nil {
		err = tx.QueryRow(r.Context(), `SELECT trim(first_name || ' ' || last_name) FROM users WHERE id=$1`, uid).Scan(&name)
	}
	if err != nil {
		expenseError(w, 500, "Could not resolve expense details")
		return
	}
	var id string
	err = tx.QueryRow(r.Context(), `INSERT INTO expenses(expense_date,category_id,category_name,title,description,amount,payee,payment_method,payment_reference,department_id,department_name,recorded_by,recorded_by_name) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id::text`, date, categoryID, category, title, reason, amount, payee, method, reference, r.FormValue("department_id"), dept, uid, name).Scan(&id)
	if err == nil && content != nil {
		_, err = tx.Exec(r.Context(), `INSERT INTO expense_attachments(expense_id,original_name,mime_type,size_bytes,content) VALUES($1,$2,$3,$4,$5)`, id, filename, mimeType, len(content), content)
	}
	if err == nil {
		err = tx.Commit(r.Context())
	}
	if err != nil {
		expenseError(w, 500, "Could not record expense")
		return
	}
	response.JSON(w, 201, map[string]string{"id": id})
}

const expenseProjection = `SELECT e.*, (SELECT json_build_object('name',original_name,'mime_type',mime_type,'size_bytes',size_bytes) FROM expense_attachments WHERE expense_id=e.id) AS attachment FROM expenses e`

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from, to := q.Get("from"), q.Get("to")
	if (from != "" && !validExpenseDate(from)) || (to != "" && !validExpenseDate(to)) || (from != "" && to != "" && from > to) {
		expenseError(w, 400, "Invalid date range")
		return
	}
	limit, offset := 50, 0
	if q.Get("offset") != "" {
		var err error
		offset, err = strconv.Atoi(q.Get("offset"))
		if err != nil || offset < 0 {
			expenseError(w, 400, "Invalid page offset")
			return
		}
	}
	where := ` WHERE ($1='' OR e.category_id::text=$1) AND ($2='' OR e.department_id=$2) AND ($3='' OR e.expense_date>=NULLIF($3,'')::date) AND ($4='' OR e.expense_date<=NULLIF($4,'')::date) AND ($5='' OR e.title ILIKE '%'||$5||'%' OR e.reference ILIKE '%'||$5||'%' OR e.recorded_by_name ILIKE '%'||$5||'%' OR e.payee ILIKE '%'||$5||'%')`
	args := []interface{}{q.Get("category_id"), q.Get("department_id"), from, to, q.Get("search")}
	var total int
	var amount string
	err := h.db.QueryRow(r.Context(), `SELECT count(*),COALESCE(sum(e.amount),0)::text FROM expenses e`+where, args...).Scan(&total, &amount)
	var result json.RawMessage
	if err == nil {
		err = h.db.QueryRow(r.Context(), `SELECT COALESCE(json_agg(x),'[]'::json) FROM (`+expenseProjection+where+` ORDER BY e.expense_date DESC,e.recorded_at DESC,e.id LIMIT $6 OFFSET $7) x`, append(args, limit, offset)...).Scan(&result)
	}
	if err != nil {
		expenseError(w, 500, "Could not load expenses")
		return
	}
	response.JSON(w, 200, map[string]interface{}{"expenses": result, "total": total, "total_amount": amount, "limit": limit, "offset": offset})
}
func (h *ExpenseHandler) Detail(w http.ResponseWriter, r *http.Request) {
	var result json.RawMessage
	err := h.db.QueryRow(r.Context(), `SELECT row_to_json(x) FROM (`+expenseProjection+` WHERE e.id::text=$1) x`, chi.URLParam(r, "id")).Scan(&result)
	if errors.Is(err, pgx.ErrNoRows) {
		expenseError(w, 404, "Expense not found")
		return
	}
	if err != nil {
		expenseError(w, 500, "Could not load expense")
		return
	}
	response.JSON(w, 200, result)
}
func (h *ExpenseHandler) Attachment(w http.ResponseWriter, r *http.Request) {
	var name, mimeType string
	var content []byte
	err := h.db.QueryRow(r.Context(), `SELECT original_name,mime_type,content FROM expense_attachments WHERE expense_id::text=$1`, chi.URLParam(r, "id")).Scan(&name, &mimeType, &content)
	if errors.Is(err, pgx.ErrNoRows) {
		expenseError(w, 404, "Receipt not found")
		return
	}
	if err != nil {
		expenseError(w, 500, "Could not load receipt")
		return
	}
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": name}))
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(content)
}

func (h *ExpenseHandler) Verify(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid := currentUserID(r)
	var name string
	err := h.db.QueryRow(r.Context(), `SELECT trim(first_name || ' ' || last_name) FROM users WHERE id=$1`, uid).Scan(&name)
	if err != nil {
		expenseError(w, 500, "User not found")
		return
	}
	_, err = h.db.Exec(r.Context(), `UPDATE expenses SET status='VERIFIED', verified_by=$1, verified_by_name=$2, verified_at=NOW() WHERE id::text=$3`, uid, name, id)
	if err != nil {
		expenseError(w, 500, "Could not verify expense")
		return
	}
	response.JSON(w, 200, map[string]string{"message": "Expense verified successfully", "status": "VERIFIED"})
}

func (h *ExpenseHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid := currentUserID(r)
	var name string
	err := h.db.QueryRow(r.Context(), `SELECT trim(first_name || ' ' || last_name) FROM users WHERE id=$1`, uid).Scan(&name)
	if err != nil {
		expenseError(w, 500, "User not found")
		return
	}
	_, err = h.db.Exec(r.Context(), `UPDATE expenses SET status='APPROVED', approved_by=$1, approved_by_name=$2, approved_at=NOW() WHERE id::text=$3`, uid, name, id)
	if err != nil {
		expenseError(w, 500, "Could not approve expense")
		return
	}
	response.JSON(w, 200, map[string]string{"message": "Expense approved successfully", "status": "APPROVED"})
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tx, err := h.db.Begin(r.Context())
	if err != nil {
		expenseError(w, 500, "Could not delete expense")
		return
	}
	defer tx.Rollback(r.Context())
	_, _ = tx.Exec(r.Context(), `DELETE FROM expense_attachments WHERE expense_id::text=$1`, id)
	tag, err := tx.Exec(r.Context(), `DELETE FROM expenses WHERE id::text=$1`, id)
	if err != nil || tag.RowsAffected() == 0 {
		expenseError(w, 404, "Expense not found or could not be deleted")
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		expenseError(w, 500, "Could not commit deletion")
		return
	}
	response.JSON(w, 200, map[string]string{"message": "Expense deleted successfully"})
}

