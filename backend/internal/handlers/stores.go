package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoresHandler struct {
	db *pgxpool.Pool
}

func NewStoresHandler(db *pgxpool.Pool) *StoresHandler {
	return &StoresHandler{db: db}
}

type StoreItem struct {
	ID                  string    `json:"id"`
	SKUCode             string    `json:"sku_code"`
	ItemName            string    `json:"item_name"`
	Category            string    `json:"category"`
	UnitOfMeasure       string    `json:"unit_of_measure"`
	UnitCostUGX         float64   `json:"unit_cost_ugx"`
	QuantityOnHand      int       `json:"quantity_on_hand"`
	MinimumReorderLevel int       `json:"minimum_reorder_level"`
	WarehouseBinLocation string   `json:"warehouse_bin_location"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type GRNItem struct {
	ID                 string          `json:"id"`
	GRNNumber          string          `json:"grn_number"`
	POReference        string          `json:"po_reference"`
	SupplierName       string          `json:"supplier_name"`
	TotalReceivedValue float64         `json:"total_received_value_ugx"`
	ItemsReceived      json.RawMessage `json:"items_received"`
	ReceivedAt         time.Time       `json:"received_at"`
}

// GET /api/v1/stores/inventory
func (h *StoresHandler) ListInventory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	category := r.URL.Query().Get("category")

	query := `
		SELECT id, sku_code, item_name, category, unit_of_measure, unit_cost_ugx, quantity_on_hand, minimum_reorder_level, warehouse_bin_location, status, created_at, updated_at
		FROM store_inventory_items
		WHERE ($1 = '' OR category = $1)
		ORDER BY item_name ASC
	`

	rows, err := h.db.Query(ctx, query, category)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch inventory")
		return
	}
	defer rows.Close()

	var list []StoreItem
	for rows.Next() {
		var item StoreItem
		err := rows.Scan(
			&item.ID, &item.SKUCode, &item.ItemName, &item.Category, &item.UnitOfMeasure,
			&item.UnitCostUGX, &item.QuantityOnHand, &item.MinimumReorderLevel,
			&item.WarehouseBinLocation, &item.Status, &item.CreatedAt, &item.UpdatedAt,
		)
		if err == nil {
			list = append(list, item)
		}
	}

	if list == nil {
		list = []StoreItem{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"items": list,
		"count": len(list),
	})
}

// POST /api/v1/stores/inventory
func (h *StoresHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SKUCode       string  `json:"sku_code"`
		ItemName      string  `json:"item_name"`
		Category      string  `json:"category"`
		UnitOfMeasure string  `json:"unit_of_measure"`
		UnitCostUGX   float64 `json:"unit_cost_ugx"`
		Quantity      int     `json:"quantity_on_hand"`
		ReorderLevel  int     `json:"minimum_reorder_level"`
		BinLocation   string  `json:"warehouse_bin_location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid item payload")
		return
	}

	status := "IN_STOCK"
	if req.Quantity <= req.ReorderLevel {
		status = "LOW_STOCK"
	}

	_, err := h.db.Exec(r.Context(), `
		INSERT INTO store_inventory_items (sku_code, item_name, category, unit_of_measure, unit_cost_ugx, quantity_on_hand, minimum_reorder_level, warehouse_bin_location, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, req.SKUCode, req.ItemName, req.Category, req.UnitOfMeasure, req.UnitCostUGX, req.Quantity, req.ReorderLevel, req.BinLocation, status)

	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not create store item")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Store item added successfully",
	})
}

// GET /api/v1/stores/grn
func (h *StoresHandler) ListGRNs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, grn_number, po_reference, supplier_name, total_received_value_ugx, COALESCE(items_received, '[]'::jsonb), received_at
		FROM goods_received_notes
		ORDER BY received_at DESC
	`)
	if err != nil {
		response.Err(w, http.StatusInternalServerError, "SERVER_ERROR", "Could not fetch GRN records")
		return
	}
	defer rows.Close()

	var list []GRNItem
	for rows.Next() {
		var item GRNItem
		err := rows.Scan(
			&item.ID, &item.GRNNumber, &item.POReference, &item.SupplierName,
			&item.TotalReceivedValue, &item.ItemsReceived, &item.ReceivedAt,
		)
		if err == nil {
			list = append(list, item)
		}
	}

	if list == nil {
		list = []GRNItem{}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"grns":  list,
		"count": len(list),
	})
}
