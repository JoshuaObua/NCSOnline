package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LegalHandler struct {
	db *pgxpool.Pool
}

func NewLegalHandler(db *pgxpool.Pool) *LegalHandler {
	return &LegalHandler{db: db}
}

type LegalContractItem struct {
	ID                string    `json:"id"`
	ContractReference string    `json:"contract_reference"`
	Title             string    `json:"title"`
	ContractType      string    `json:"contract_type"`
	SecondParty       string    `json:"second_party"`
	ContractValueUGX  float64   `json:"contract_value_ugx"`
	StartDate         time.Time `json:"start_date"`
	ExpiryDate        time.Time `json:"expiry_date"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

type FederationDisputeItem struct {
	ID              string    `json:"id"`
	CaseNumber      string    `json:"case_number"`
	FederationName  string    `json:"federation_name"`
	ComplainantName string    `json:"complainant_name"`
	RespondentName  string    `json:"respondent_name"`
	SubjectMatter   string    `json:"subject_matter"`
	DisputeCategory string    `json:"dispute_category"`
	FilingDate      time.Time `json:"filing_date"`
	CaseStatus      string    `json:"case_status"`
	RulingSummary   *string   `json:"ruling_summary,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// GET /api/v1/legal/contracts
func (h *LegalHandler) ListContracts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, contract_reference, title, contract_type, second_party, contract_value_ugx, start_date, expiry_date, status, created_at
		FROM legal_contracts
		ORDER BY expiry_date ASC
	`)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch legal contracts")
		return
	}
	defer rows.Close()

	var list []LegalContractItem
	for rows.Next() {
		var c LegalContractItem
		err := rows.Scan(
			&c.ID, &c.ContractReference, &c.Title, &c.ContractType, &c.SecondParty,
			&c.ContractValueUGX, &c.StartDate, &c.ExpiryDate, &c.Status, &c.CreatedAt,
		)
		if err == nil {
			list = append(list, c)
		}
	}

	if list == nil {
		list = []LegalContractItem{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"contracts": list,
		"count":     len(list),
	})
}

// POST /api/v1/legal/contracts
func (h *LegalHandler) CreateContract(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContractReference string  `json:"contract_reference"`
		Title             string  `json:"title"`
		ContractType      string  `json:"contract_type"`
		SecondParty       string  `json:"second_party"`
		ContractValueUGX  float64 `json:"contract_value_ugx"`
		StartDate         string  `json:"start_date"`
		ExpiryDate        string  `json:"expiry_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid contract payload")
		return
	}

	sDate, _ := time.Parse("2006-01-02", req.StartDate)
	eDate, _ := time.Parse("2006-01-02", req.ExpiryDate)

	_, err := h.db.Exec(r.Context(), `
		INSERT INTO legal_contracts (contract_reference, title, contract_type, second_party, contract_value_ugx, start_date, expiry_date, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'ACTIVE')
	`, req.ContractReference, req.Title, req.ContractType, req.SecondParty, req.ContractValueUGX, sDate, eDate)

	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create contract")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Legal contract recorded successfully",
	})
}

// GET /api/v1/legal/disputes
func (h *LegalHandler) ListDisputes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, case_number, federation_name, complainant_name, respondent_name, subject_matter, dispute_category, filing_date, case_status, COALESCE(ruling_summary, ''), created_at
		FROM federation_disputes
		ORDER BY filing_date DESC
	`)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch federation disputes")
		return
	}
	defer rows.Close()

	var list []FederationDisputeItem
	for rows.Next() {
		var d FederationDisputeItem
		var summary string
		err := rows.Scan(
			&d.ID, &d.CaseNumber, &d.FederationName, &d.ComplainantName, &d.RespondentName,
			&d.SubjectMatter, &d.DisputeCategory, &d.FilingDate, &d.CaseStatus, &summary, &d.CreatedAt,
		)
		if err == nil {
			if summary != "" {
				d.RulingSummary = &summary
			}
			list = append(list, d)
		}
	}

	if list == nil {
		list = []FederationDisputeItem{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"disputes": list,
		"count":    len(list),
	})
}
