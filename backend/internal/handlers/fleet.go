package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FleetHandler struct {
	db *pgxpool.Pool
}

func NewFleetHandler(db *pgxpool.Pool) *FleetHandler {
	return &FleetHandler{db: db}
}

type FleetVehicleItem struct {
	ID                   string    `json:"id"`
	RegistrationNumber   string    `json:"registration_number"`
	MakeModel            string    `json:"make_model"`
	VehicleType          string    `json:"vehicle_type"`
	SeatingCapacity      int       `json:"seating_capacity"`
	FuelType             string    `json:"fuel_type"`
	CurrentMileageKM     int       `json:"current_mileage_km"`
	NextServiceMileageKM int       `json:"next_service_mileage_km"`
	InsuranceExpiryDate  time.Time `json:"insurance_expiry_date"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
}

type TripRequisitionItem struct {
	ID                   string    `json:"id"`
	TripCode             string    `json:"trip_code"`
	RequestingDepartment string    `json:"requesting_department"`
	VehicleReg           string    `json:"vehicle_reg"`
	AssignedDriverName   string    `json:"assigned_driver_name"`
	Destination          string    `json:"destination"`
	TripPurpose          string    `json:"trip_purpose"`
	DepartureTime        time.Time `json:"departure_time"`
	ReturnTime           time.Time `json:"return_time"`
	ApprovalStatus       string    `json:"approval_status"`
	CreatedAt            time.Time `json:"created_at"`
}

// GET /api/v1/fleet/vehicles
func (h *FleetHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, registration_number, make_model, vehicle_type, seating_capacity, fuel_type, current_mileage_km, next_service_mileage_km, insurance_expiry_date, status, created_at
		FROM fleet_vehicles
		ORDER BY registration_number ASC
	`)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch vehicles")
		return
	}
	defer rows.Close()

	var list []FleetVehicleItem
	for rows.Next() {
		var v FleetVehicleItem
		err := rows.Scan(
			&v.ID, &v.RegistrationNumber, &v.MakeModel, &v.VehicleType,
			&v.SeatingCapacity, &v.FuelType, &v.CurrentMileageKM,
			&v.NextServiceMileageKM, &v.InsuranceExpiryDate, &v.Status, &v.CreatedAt,
		)
		if err == nil {
			list = append(list, v)
		}
	}

	if list == nil {
		list = []FleetVehicleItem{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"vehicles": list,
		"count":    len(list),
	})
}

// GET /api/v1/fleet/trips
func (h *FleetHandler) ListTrips(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, trip_code, requesting_department, vehicle_reg, assigned_driver_name, destination, trip_purpose, departure_time, return_time, approval_status, created_at
		FROM trip_requisitions
		ORDER BY departure_time DESC
	`)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch trips")
		return
	}
	defer rows.Close()

	var list []TripRequisitionItem
	for rows.Next() {
		var t TripRequisitionItem
		err := rows.Scan(
			&t.ID, &t.TripCode, &t.RequestingDepartment, &t.VehicleReg,
			&t.AssignedDriverName, &t.Destination, &t.TripPurpose,
			&t.DepartureTime, &t.ReturnTime, &t.ApprovalStatus, &t.CreatedAt,
		)
		if err == nil {
			list = append(list, t)
		}
	}

	if list == nil {
		list = []TripRequisitionItem{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"trips": list,
		"count": len(list),
	})
}

// POST /api/v1/fleet/trips
func (h *FleetHandler) CreateTrip(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TripCode      string `json:"trip_code"`
		Department    string `json:"requesting_department"`
		VehicleReg    string `json:"vehicle_reg"`
		DriverName    string `json:"assigned_driver_name"`
		Destination   string `json:"destination"`
		Purpose       string `json:"trip_purpose"`
		DepartureTime string `json:"departure_time"`
		ReturnTime    string `json:"return_time"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid trip payload")
		return
	}

	dep, _ := time.Parse(time.RFC3339, req.DepartureTime)
	ret, _ := time.Parse(time.RFC3339, req.ReturnTime)

	_, err := h.db.Exec(r.Context(), `
		INSERT INTO trip_requisitions (trip_code, requesting_department, vehicle_reg, assigned_driver_name, destination, trip_purpose, departure_time, return_time, approval_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'DISPATCHED')
	`, req.TripCode, req.Department, req.VehicleReg, req.DriverName, req.Destination, req.Purpose, dep, ret)

	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create trip requisition")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Trip created and vehicle dispatched successfully",
	})
}
