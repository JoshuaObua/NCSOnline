package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FacilitiesHandler struct {
	db *pgxpool.Pool
}

func NewFacilitiesHandler(db *pgxpool.Pool) *FacilitiesHandler {
	return &FacilitiesHandler{db: db}
}

type VenueBookingItem struct {
	ID                string    `json:"id"`
	BookingReference  string    `json:"booking_reference"`
	VenueName         string    `json:"venue_name"`
	ClientName        string    `json:"client_name"`
	ClientType        string    `json:"client_type"`
	EventTitle        string    `json:"event_title"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	TariffCategory    string    `json:"tariff_category"`
	TotalFeeUGX       float64   `json:"total_fee_ugx"`
	CautionDepositUGX float64   `json:"caution_deposit_ugx"`
	PaymentStatus     string    `json:"payment_status"`
	BookingStatus     string    `json:"booking_status"`
	CreatedAt         time.Time `json:"created_at"`
}

type HostelOccupancyItem struct {
	ID             string    `json:"id"`
	RoomNumber     string    `json:"room_number"`
	AthleteName    string    `json:"athlete_name"`
	FederationName string    `json:"federation_name"`
	Gender         string    `json:"gender"`
	CheckInDate    time.Time `json:"check_in_date"`
	CheckOutDate   time.Time `json:"check_out_date"`
	Status         string    `json:"status"`
	KeyIssued      bool      `json:"key_issued"`
	CreatedAt      time.Time `json:"created_at"`
}

// GET /api/v1/facilities/bookings
func (h *FacilitiesHandler) ListBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, booking_reference, venue_name, client_name, client_type, event_title, start_date, end_date, tariff_category, total_fee_ugx, caution_deposit_ugx, payment_status, booking_status, created_at
		FROM venue_bookings
		ORDER BY start_date DESC
	`)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch venue bookings")
		return
	}
	defer rows.Close()

	var list []VenueBookingItem
	for rows.Next() {
		var b VenueBookingItem
		err := rows.Scan(
			&b.ID, &b.BookingReference, &b.VenueName, &b.ClientName, &b.ClientType,
			&b.EventTitle, &b.StartDate, &b.EndDate, &b.TariffCategory, &b.TotalFeeUGX,
			&b.CautionDepositUGX, &b.PaymentStatus, &b.BookingStatus, &b.CreatedAt,
		)
		if err == nil {
			list = append(list, b)
		}
	}

	if list == nil {
		list = []VenueBookingItem{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"bookings": list,
		"count":    len(list),
	})
}

// POST /api/v1/facilities/bookings
func (h *FacilitiesHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BookingReference string  `json:"booking_reference"`
		VenueName        string  `json:"venue_name"`
		ClientName       string  `json:"client_name"`
		ClientType       string  `json:"client_type"`
		EventTitle       string  `json:"event_title"`
		StartDate        string  `json:"start_date"`
		EndDate          string  `json:"end_date"`
		TotalFeeUGX      float64 `json:"total_fee_ugx"`
		CautionDeposit   float64 `json:"caution_deposit_ugx"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid booking payload")
		return
	}

	sDate, _ := time.Parse("2006-01-02", req.StartDate)
	eDate, _ := time.Parse("2006-01-02", req.EndDate)

	_, err := h.db.Exec(r.Context(), `
		INSERT INTO venue_bookings (booking_reference, venue_name, client_name, client_type, event_title, start_date, end_date, total_fee_ugx, caution_deposit_ugx, payment_status, booking_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'FULLY_PAID', 'CONFIRMED')
	`, req.BookingReference, req.VenueName, req.ClientName, req.ClientType, req.EventTitle, sDate, eDate, req.TotalFeeUGX, req.CautionDeposit)

	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create venue booking")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Venue booking created successfully",
	})
}

// GET /api/v1/facilities/hostels
func (h *FacilitiesHandler) ListHostels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, room_number, athlete_name, federation_name, gender, check_in_date, check_out_date, status, key_issued, created_at
		FROM hostel_occupancies
		ORDER BY check_in_date DESC
	`)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch hostel occupancies")
		return
	}
	defer rows.Close()

	var list []HostelOccupancyItem
	for rows.Next() {
		var item HostelOccupancyItem
		err := rows.Scan(
			&item.ID, &item.RoomNumber, &item.AthleteName, &item.FederationName,
			&item.Gender, &item.CheckInDate, &item.CheckOutDate, &item.Status,
			&item.KeyIssued, &item.CreatedAt,
		)
		if err == nil {
			list = append(list, item)
		}
	}

	if list == nil {
		list = []HostelOccupancyItem{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"occupancies": list,
		"count":       len(list),
	})
}
