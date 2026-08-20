-- Migration: 078_correct_fixed_asset_baseline_totals.sql
-- Align already-migrated fixed asset data with Docs/plans/fixed-assets-and-dashboards-plan.md.
-- The Excel cell for M1007059 FB_COST displays ########### while ADJUSTED COST holds UGX 1.16B.
-- Store FB_COST as 0 for the baseline total and retain the adjusted valuation separately.

UPDATE fixed_assets
SET fb_cost = 0.00,
    adjusted_cost = 1160000000.00,
    updated_at = NOW()
WHERE asset_number = 'M1007059'
  AND tag_number = '166BLNG5';
