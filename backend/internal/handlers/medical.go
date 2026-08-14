package handlers

import (
	"net/http"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MedicalHandler struct {
	db *pgxpool.Pool
}

func NewMedicalHandler(db *pgxpool.Pool) *MedicalHandler {
	return &MedicalHandler{db: db}
}

type AthleteScreeningItem struct {
	ID                   string    `json:"id"`
	AthleteName          string    `json:"athlete_name"`
	Discipline           string    `json:"discipline"`
	FederationName       string    `json:"federation_name"`
	ScreeningDate        time.Time `json:"screening_date"`
	CardiovascularPassed bool      `json:"cardiovascular_passed"`
	ECGFinding           string    `json:"ecg_finding"`
	BloodPressure        string    `json:"blood_pressure"`
	FitnessVerdict       string    `json:"fitness_verdict"`
	CreatedAt            time.Time `json:"created_at"`
}

type AthleteInjuryItem struct {
	ID           string    `json:"id"`
	AthleteName  string    `json:"athlete_name"`
	Discipline   string    `json:"discipline"`
	InjurySite   string    `json:"injury_site"`
	InjuryNature string    `json:"injury_nature"`
	Severity     string    `json:"severity"`
	RehabStatus  string    `json:"rehab_status"`
	IncidentDate time.Time `json:"incident_date"`
	CreatedAt    time.Time `json:"created_at"`
}

type AntiDopingItem struct {
	ID             string    `json:"id"`
	AthleteName    string    `json:"athlete_name"`
	SampleCode     string    `json:"sample_code"`
	TestType       string    `json:"test_type"`
	CollectionDate time.Time `json:"collection_date"`
	ResultStatus   string    `json:"result_status"`
	HasActiveTUE   bool      `json:"has_active_tue"`
	CreatedAt      time.Time `json:"created_at"`
}

// GET /api/v1/medical/screenings
func (h *MedicalHandler) ListScreenings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, athlete_name, discipline, federation_name, screening_date, cardiovascular_passed, ecg_finding, blood_pressure, fitness_verdict, created_at
		FROM athlete_medical_screenings
		ORDER BY screening_date DESC
	`)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch medical screenings")
		return
	}
	defer rows.Close()

	var list []AthleteScreeningItem
	for rows.Next() {
		var item AthleteScreeningItem
		err := rows.Scan(
			&item.ID, &item.AthleteName, &item.Discipline, &item.FederationName,
			&item.ScreeningDate, &item.CardiovascularPassed, &item.ECGFinding,
			&item.BloodPressure, &item.FitnessVerdict, &item.CreatedAt,
		)
		if err == nil {
			list = append(list, item)
		}
	}

	if list == nil {
		list = []AthleteScreeningItem{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"screenings": list,
		"count":      len(list),
	})
}

// GET /api/v1/medical/injuries
func (h *MedicalHandler) ListInjuries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, athlete_name, discipline, injury_site, injury_nature, severity, rehab_status, incident_date, created_at
		FROM athlete_injuries
		ORDER BY incident_date DESC
	`)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch athlete injuries")
		return
	}
	defer rows.Close()

	var list []AthleteInjuryItem
	for rows.Next() {
		var item AthleteInjuryItem
		err := rows.Scan(
			&item.ID, &item.AthleteName, &item.Discipline, &item.InjurySite,
			&item.InjuryNature, &item.Severity, &item.RehabStatus,
			&item.IncidentDate, &item.CreatedAt,
		)
		if err == nil {
			list = append(list, item)
		}
	}

	if list == nil {
		list = []AthleteInjuryItem{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"injuries": list,
		"count":    len(list),
	})
}

// GET /api/v1/medical/antidoping
func (h *MedicalHandler) ListAntiDoping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, athlete_name, sample_code, test_type, collection_date, result_status, has_active_tue, created_at
		FROM antidoping_records
		ORDER BY collection_date DESC
	`)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch anti-doping records")
		return
	}
	defer rows.Close()

	var list []AntiDopingItem
	for rows.Next() {
		var item AntiDopingItem
		err := rows.Scan(
			&item.ID, &item.AthleteName, &item.SampleCode, &item.TestType,
			&item.CollectionDate, &item.ResultStatus, &item.HasActiveTUE, &item.CreatedAt,
		)
		if err == nil {
			list = append(list, item)
		}
	}

	if list == nil {
		list = []AntiDopingItem{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"records": list,
		"count":   len(list),
	})
}
