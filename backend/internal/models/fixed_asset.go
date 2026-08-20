package models

import (
	"time"
)

// FixedAsset represents a record in the fixed_assets table
type FixedAsset struct {
	ID                      string    `json:"id" db:"id"`
	InterfaceLineNumber     string    `json:"interface_line_number" db:"interface_line_number"`
	AssetBook               string    `json:"asset_book" db:"asset_book"`
	AssetNumber             string    `json:"asset_number" db:"asset_number"`
	TagNumber               string    `json:"tag_number" db:"tag_number"`
	AssetDescription        string    `json:"asset_description" db:"asset_description"`
	CategorySegment1        string    `json:"category_segment1" db:"category_segment1"`
	CategorySegment3        string    `json:"category_segment3" db:"category_segment3"`
	CategorySegment4        string    `json:"category_segment4" db:"category_segment4"`
	AssetUnits              int       `json:"asset_units" db:"asset_units"`
	FBCost                  float64   `json:"fb_cost" db:"fb_cost"`
	AdjustedCost            float64   `json:"adjusted_cost" db:"adjusted_cost"`
	DatePlacedInService     string    `json:"date_placed_in_service" db:"date_placed_in_service"`
	CustodianDepartment     string    `json:"custodian_department" db:"custodian_department"`
	LocationBuilding        string    `json:"location_building" db:"location_building"`
	LocationRoom            string    `json:"location_room" db:"location_room"`
	DepreciationMethod      string    `json:"depreciation_method" db:"depreciation_method"`
	UsefulLifeYears         int       `json:"useful_life_years" db:"useful_life_years"`
	AccumulatedDepreciation float64   `json:"accumulated_depreciation" db:"accumulated_depreciation"`
	NetBookValue            float64   `json:"net_book_value"`
	Status                  string    `json:"status" db:"status"`
	VerificationStatus      string    `json:"verification_status" db:"verification_status"`
	LastVerifiedAt          *time.Time `json:"last_verified_at" db:"last_verified_at"`
	LastVerifiedBy          *string   `json:"last_verified_by" db:"last_verified_by"`
	WorksheetSource         string    `json:"worksheet_source" db:"worksheet_source"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}

// AssetTransactionLog represents an audited revaluation, depreciation, or transfer event
type AssetTransactionLog struct {
	ID              string    `json:"id" db:"id"`
	AssetID         string    `json:"asset_id" db:"asset_id"`
	TransactionType string    `json:"transaction_type" db:"transaction_type"`
	PreviousVal     float64   `json:"previous_val" db:"previous_val"`
	NewVal          float64   `json:"new_val" db:"new_val"`
	Notes           string    `json:"notes" db:"notes"`
	PerformedBy     *string   `json:"performed_by" db:"performed_by"`
	ApprovedBy      *string   `json:"approved_by" db:"approved_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// FixedAssetSummary holds high-level portfolio KPIs and pivot table replacement data
type FixedAssetSummary struct {
	TotalAssets             int                   `json:"total_assets"`
	TotalFBCost             float64               `json:"total_fb_cost"`
	TotalAdjustedCost       float64               `json:"total_adjusted_cost"`
	TotalAccumulatedDeprec  float64               `json:"total_accumulated_deprec"`
	TotalNetBookValue       float64               `json:"total_net_book_value"`
	VerifiedAssets          int                   `json:"verified_assets"`
	DiscrepancyAssets       int                   `json:"discrepancy_assets"`
	CategorySummaries       []CategorySummaryItem `json:"category_summaries"`
}

type CategorySummaryItem struct {
	CategorySegment1 string  `json:"category_segment1"`
	CategorySegment3 string  `json:"category_segment3"`
	CategorySegment4 string  `json:"category_segment4"`
	AssetUnits       int     `json:"asset_units"`
	TotalFBCost      float64 `json:"total_fb_cost"`
	TotalAdjustedCost float64 `json:"total_adjusted_cost"`
	TotalNetBookValue float64 `json:"total_net_book_value"`
}

type AssetRevaluationRequest struct {
	AssetID     string  `json:"asset_id"`
	NewCost     float64 `json:"new_cost"`
	Notes       string  `json:"notes"`
}

type AssetDepreciationRequest struct {
	Period string `json:"period"` // e.g. "2026-07"
}

type AssetVerificationRequest struct {
	AssetID string `json:"asset_id"`
	Status  string `json:"status"` // VERIFIED, DISCREPANCY, MISSING
	Notes   string `json:"notes"`
}
