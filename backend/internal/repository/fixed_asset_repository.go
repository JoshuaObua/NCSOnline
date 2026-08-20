package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ncsintranet/internal/models"
)

type FixedAssetRepository interface {
	ListFixedAssets(ctx context.Context, category, search string, limit, offset int) ([]models.FixedAsset, int, error)
	GetFixedAssetByID(ctx context.Context, id string) (*models.FixedAsset, error)
	GetFixedAssetSummary(ctx context.Context) (*models.FixedAssetSummary, error)
	CreateFixedAsset(ctx context.Context, asset *models.FixedAsset) error
	RevalueFixedAsset(ctx context.Context, assetID string, newCost float64, notes string, userID string) error
	RunDepreciation(ctx context.Context, period string, userID string) (int, float64, error)
	VerifyFixedAsset(ctx context.Context, assetID string, status string, notes string, userID string) error
	GetTransactionLogs(ctx context.Context, assetID string) ([]models.AssetTransactionLog, error)
}

type postgresFixedAssetRepository struct {
	db *pgxpool.Pool
}

func NewFixedAssetRepository(db *pgxpool.Pool) FixedAssetRepository {
	return &postgresFixedAssetRepository{db: db}
}

func (r *postgresFixedAssetRepository) ListFixedAssets(ctx context.Context, category, search string, limit, offset int) ([]models.FixedAsset, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if category != "" {
		whereClause += fmt.Sprintf(" AND (category_segment1 = $%d OR category_segment3 = $%d OR category_segment4 = $%d OR worksheet_source = $%d)", argIdx, argIdx, argIdx, argIdx)
		args = append(args, category)
		argIdx++
	}

	if search != "" {
		whereClause += fmt.Sprintf(" AND (asset_number ILIKE $%d OR tag_number ILIKE $%d OR asset_description ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM fixed_assets %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, interface_line_number, asset_book, asset_number, tag_number, asset_description,
		       category_segment1, category_segment3, category_segment4, asset_units,
		       fb_cost, adjusted_cost, TO_CHAR(date_placed_in_service, 'YYYY-MM-DD'), custodian_department,
		       location_building, location_room, depreciation_method, useful_life_years,
		       accumulated_depreciation, status, COALESCE(verification_status, 'UNVERIFIED'), last_verified_at,
		       last_verified_by, worksheet_source, created_at, updated_at
		FROM fixed_assets
		%s
		ORDER BY created_at ASC, asset_number ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var assets []models.FixedAsset
	for rows.Next() {
		var a models.FixedAsset
		var lastVerifiedBy *string
		if err := rows.Scan(
			&a.ID, &a.InterfaceLineNumber, &a.AssetBook, &a.AssetNumber, &a.TagNumber, &a.AssetDescription,
			&a.CategorySegment1, &a.CategorySegment3, &a.CategorySegment4, &a.AssetUnits,
			&a.FBCost, &a.AdjustedCost, &a.DatePlacedInService, &a.CustodianDepartment,
			&a.LocationBuilding, &a.LocationRoom, &a.DepreciationMethod, &a.UsefulLifeYears,
			&a.AccumulatedDepreciation, &a.Status, &a.VerificationStatus, &a.LastVerifiedAt,
			&lastVerifiedBy, &a.WorksheetSource, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		a.LastVerifiedBy = lastVerifiedBy
		a.NetBookValue = a.AdjustedCost - a.AccumulatedDepreciation
		assets = append(assets, a)
	}

	return assets, total, nil
}

func (r *postgresFixedAssetRepository) GetFixedAssetByID(ctx context.Context, id string) (*models.FixedAsset, error) {
	query := `
		SELECT id, interface_line_number, asset_book, asset_number, tag_number, asset_description,
		       category_segment1, category_segment3, category_segment4, asset_units,
		       fb_cost, adjusted_cost, TO_CHAR(date_placed_in_service, 'YYYY-MM-DD'), custodian_department,
		       location_building, location_room, depreciation_method, useful_life_years,
		       accumulated_depreciation, status, COALESCE(verification_status, 'UNVERIFIED'), last_verified_at,
		       last_verified_by, worksheet_source, created_at, updated_at
		FROM fixed_assets
		WHERE id::text = $1 OR asset_number = $1 OR tag_number = $1
	`
	var a models.FixedAsset
	var lastVerifiedBy *string
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.InterfaceLineNumber, &a.AssetBook, &a.AssetNumber, &a.TagNumber, &a.AssetDescription,
		&a.CategorySegment1, &a.CategorySegment3, &a.CategorySegment4, &a.AssetUnits,
		&a.FBCost, &a.AdjustedCost, &a.DatePlacedInService, &a.CustodianDepartment,
		&a.LocationBuilding, &a.LocationRoom, &a.DepreciationMethod, &a.UsefulLifeYears,
		&a.AccumulatedDepreciation, &a.Status, &a.VerificationStatus, &a.LastVerifiedAt,
		&lastVerifiedBy, &a.WorksheetSource, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	a.LastVerifiedBy = lastVerifiedBy
	a.NetBookValue = a.AdjustedCost - a.AccumulatedDepreciation
	return &a, nil
}

func (r *postgresFixedAssetRepository) GetFixedAssetSummary(ctx context.Context) (*models.FixedAssetSummary, error) {
	summaryQuery := `
		SELECT 
			COUNT(*),
			COALESCE(SUM(fb_cost), 0),
			COALESCE(SUM(adjusted_cost), 0),
			COALESCE(SUM(accumulated_depreciation), 0),
			COALESCE(SUM(adjusted_cost - accumulated_depreciation), 0),
			COALESCE(SUM(CASE WHEN verification_status = 'VERIFIED' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN verification_status = 'DISCREPANCY' THEN 1 ELSE 0 END), 0)
		FROM fixed_assets
	`

	var s models.FixedAssetSummary
	err := r.db.QueryRow(ctx, summaryQuery).Scan(
		&s.TotalAssets, &s.TotalFBCost, &s.TotalAdjustedCost,
		&s.TotalAccumulatedDeprec, &s.TotalNetBookValue,
		&s.VerifiedAssets, &s.DiscrepancyAssets,
	)
	if err != nil {
		return nil, err
	}

	catQuery := `
		SELECT 
			category_segment1,
			category_segment3,
			category_segment4,
			SUM(asset_units),
			SUM(fb_cost),
			SUM(adjusted_cost),
			SUM(adjusted_cost - accumulated_depreciation)
		FROM fixed_assets
		GROUP BY category_segment1, category_segment3, category_segment4
		ORDER BY category_segment1, category_segment3, category_segment4
	`

	rows, err := r.db.Query(ctx, catQuery)
	if err == nil {
		defer rows.Close()
		var categories []models.CategorySummaryItem
		for rows.Next() {
			var item models.CategorySummaryItem
			if err := rows.Scan(
				&item.CategorySegment1, &item.CategorySegment3, &item.CategorySegment4,
				&item.AssetUnits, &item.TotalFBCost, &item.TotalAdjustedCost, &item.TotalNetBookValue,
			); err == nil {
				categories = append(categories, item)
			}
		}
		s.CategorySummaries = categories
	}

	return &s, nil
}

func (r *postgresFixedAssetRepository) CreateFixedAsset(ctx context.Context, asset *models.FixedAsset) error {
	query := `
		INSERT INTO fixed_assets (
			interface_line_number, asset_book, asset_number, tag_number, asset_description,
			category_segment1, category_segment3, category_segment4, asset_units,
			fb_cost, adjusted_cost, date_placed_in_service, custodian_department,
			location_building, location_room, depreciation_method, useful_life_years, worksheet_source
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		asset.InterfaceLineNumber, asset.AssetBook, asset.AssetNumber, asset.TagNumber, asset.AssetDescription,
		asset.CategorySegment1, asset.CategorySegment3, asset.CategorySegment4, asset.AssetUnits,
		asset.FBCost, asset.AdjustedCost, asset.DatePlacedInService, asset.CustodianDepartment,
		asset.LocationBuilding, asset.LocationRoom, asset.DepreciationMethod, asset.UsefulLifeYears, asset.WorksheetSource,
	).Scan(&asset.ID)
}

func (r *postgresFixedAssetRepository) RevalueFixedAsset(ctx context.Context, assetID string, newCost float64, notes string, userID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var oldCost float64
	err = tx.QueryRow(ctx, "SELECT adjusted_cost FROM fixed_assets WHERE id::text = $1 OR asset_number = $1", assetID).Scan(&oldCost)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE fixed_assets SET adjusted_cost = $1, updated_at = NOW() WHERE id::text = $2 OR asset_number = $2", newCost, assetID)
	if err != nil {
		return err
	}

	var uid *string
	if userID != "" {
		uid = &userID
	}

	logQuery := `
		INSERT INTO asset_transaction_logs (asset_id, transaction_type, previous_val, new_val, notes, performed_by)
		SELECT id, 'REVALUATION', $2, $3, $4, $5 FROM fixed_assets WHERE id::text = $1 OR asset_number = $1
	`
	_, err = tx.Exec(ctx, logQuery, assetID, oldCost, newCost, notes, uid)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *postgresFixedAssetRepository) RunDepreciation(ctx context.Context, period string, userID string) (int, float64, error) {
	query := `
		UPDATE fixed_assets
		SET accumulated_depreciation = LEAST(adjusted_cost, accumulated_depreciation + (adjusted_cost / NULLIF(useful_life_years * 12, 0))),
		    updated_at = NOW()
		WHERE status = 'ACTIVE' AND useful_life_years > 0 AND adjusted_cost > accumulated_depreciation
	`
	tag, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, 0, err
	}

	var totalDeprecRun float64
	_ = r.db.QueryRow(ctx, "SELECT COALESCE(SUM(adjusted_cost / NULLIF(useful_life_years * 12, 0)), 0) FROM fixed_assets WHERE status = 'ACTIVE' AND useful_life_years > 0").Scan(&totalDeprecRun)

	return int(tag.RowsAffected()), totalDeprecRun, nil
}

func (r *postgresFixedAssetRepository) VerifyFixedAsset(ctx context.Context, assetID string, status string, notes string, userID string) error {
	var uid *string
	if userID != "" {
		uid = &userID
	}

	query := `
		UPDATE fixed_assets
		SET verification_status = $1, last_verified_at = NOW(), last_verified_by = $2, updated_at = NOW()
		WHERE id::text = $3 OR asset_number = $3 OR tag_number = $3
	`
	_, err := r.db.Exec(ctx, query, status, uid, assetID)
	if err != nil {
		return err
	}

	logQuery := `
		INSERT INTO asset_transaction_logs (asset_id, transaction_type, notes, performed_by)
		SELECT id, 'VERIFICATION', $2, $3 FROM fixed_assets WHERE id::text = $1 OR asset_number = $1 OR tag_number = $1
	`
	_, _ = r.db.Exec(ctx, logQuery, assetID, notes, uid)
	return nil
}

func (r *postgresFixedAssetRepository) GetTransactionLogs(ctx context.Context, assetID string) ([]models.AssetTransactionLog, error) {
	query := `
		SELECT l.id, l.asset_id, l.transaction_type, COALESCE(l.previous_val, 0), COALESCE(l.new_val, 0), COALESCE(l.notes, ''), l.created_at
		FROM asset_transaction_logs l
		JOIN fixed_assets a ON a.id = l.asset_id
		WHERE a.id::text = $1 OR a.asset_number = $1
		ORDER BY l.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AssetTransactionLog
	for rows.Next() {
		var l models.AssetTransactionLog
		if err := rows.Scan(&l.ID, &l.AssetID, &l.TransactionType, &l.PreviousVal, &l.NewVal, &l.Notes, &l.CreatedAt); err == nil {
			logs = append(logs, l)
		}
	}
	return logs, nil
}
