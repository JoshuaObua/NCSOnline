-- Migration: 067_create_gs_appraisals
-- Creates tables for General Secretary Master Appraisal & Valuation Center

-- 1. Create Executive Appraisals Master Table
CREATE TABLE IF NOT EXISTS executive_appraisals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appraisal_code VARCHAR(64) UNIQUE NOT NULL,
    appraisal_type VARCHAR(64) NOT NULL, -- PROPERTY_INVENTORY, STAFF_PERFORMANCE, FINANCIAL_EFFICIENCY, FIXED_ASSETS_LEDGER
    evaluation_period VARCHAR(64) NOT NULL, -- e.g. 'FY 2025/2026 Q4'
    evaluated_by TEXT REFERENCES users(id),
    total_portfolio_value_ugx NUMERIC(18,2) DEFAULT 0.00,
    overall_performance_score NUMERIC(5,2) DEFAULT 0.00,
    executive_findings TEXT,
    accounting_officer_verdict VARCHAR(32) NOT NULL DEFAULT 'APPROVED', -- APPROVED, APPROVED_WITH_CONDITIONS, REVISED, REJECTED
    statutory_signoff_timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    digital_signature_hash VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Create Asset Valuation Signoffs Table
CREATE TABLE IF NOT EXISTS asset_valuation_signoffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appraisal_id UUID REFERENCES executive_appraisals(id),
    asset_tag VARCHAR(128) NOT NULL,
    asset_description TEXT NOT NULL,
    category_segment VARCHAR(128) NOT NULL,
    recorded_historical_cost NUMERIC(18,2) NOT NULL,
    recorded_adjusted_cost NUMERIC(18,2) NOT NULL,
    recorded_accumulated_depreciation NUMERIC(18,2) NOT NULL DEFAULT 0.00,
    approved_net_book_value NUMERIC(18,2) NOT NULL,
    physical_condition_grade VARCHAR(32) NOT NULL DEFAULT 'GOOD', -- EXCELLENT, GOOD, FAIR, POOR, OBSOLETE
    gs_action VARCHAR(32) NOT NULL DEFAULT 'CONFIRMED_ACTIVE', -- CONFIRMED_ACTIVE, REVALUATION_APPROVED, WRITE_OFF_AUTHORIZED, DISPOSAL_RECOMMENDED
    gs_remarks TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Create Financial Efficiency Appraisals Table
CREATE TABLE IF NOT EXISTS financial_efficiency_appraisals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appraisal_id UUID REFERENCES executive_appraisals(id),
    department_name VARCHAR(64) NOT NULL,
    allocated_budget_ugx NUMERIC(18,2) NOT NULL,
    actual_expenditure_ugx NUMERIC(18,2) NOT NULL,
    budget_execution_rate NUMERIC(5,2) NOT NULL,
    ntr_target_ugx NUMERIC(18,2) DEFAULT 0.00,
    ntr_collected_ugx NUMERIC(18,2) DEFAULT 0.00,
    grant_accountability_rate NUMERIC(5,2) DEFAULT 100.00,
    audit_query_count INT DEFAULT 0,
    financial_grade VARCHAR(8) NOT NULL DEFAULT 'A',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Seed initial baseline records for testing
INSERT INTO executive_appraisals (appraisal_code, appraisal_type, evaluation_period, total_portfolio_value_ugx, overall_performance_score, executive_findings, accounting_officer_verdict) VALUES
  ('EVAL-2026-ASSETS', 'FIXED_ASSETS_LEDGER', 'FY 2025/2026 Annual', 31015914535.00, 98.0, 'Comprehensive fixed asset register verified across all 11 classes. Total portfolio NBV verified at UGX 30.38 Billion.', 'APPROVED'),
  ('EVAL-2026-FIN', 'FINANCIAL_EFFICIENCY', 'FY 2025/2026 Q4', 25000000000.00, 94.2, 'Subvention execution within approved parameters. NTR collections exceeded target at 102.4%.', 'APPROVED'),
  ('EVAL-2026-HR', 'STAFF_PERFORMANCE', 'FY 2025/2026 Semi-Annual', 0.00, 88.4, '128 staff scorecards reviewed. 24 outstanding ratings, 98 meeting expectations, 6 PIP interventions.', 'APPROVED')
ON CONFLICT (appraisal_code) DO NOTHING;

INSERT INTO asset_valuation_signoffs (asset_tag, asset_description, category_segment, recorded_historical_cost, recorded_adjusted_cost, recorded_accumulated_depreciation, approved_net_book_value, physical_condition_grade, gs_action) VALUES
  ('Plots 2-10 Coronation Ave', 'Lugogo Sports Grounds Recreation Land (8 Plots)', 'LAND', 27893933769.00, 27893933769.00, 0.00, 27893933769.00, 'EXCELLENT', 'CONFIRMED_ACTIVE'),
  ('166BLNG10', 'Lugogo Non-Residential Specialized Office & Arena Complex', 'NON RESIDENTIAL BUILDINGS', 2138452000.00, 2138452000.00, 198252000.00, 1940200000.00, 'GOOD', 'CONFIRMED_ACTIVE'),
  ('UBF 748K', 'King Long Kingo 16S Bus (National Team Transit)', 'LIGHT VEHICLES', 85000000.00, 85000000.00, 15000000.00, 70000000.00, 'GOOD', 'CONFIRMED_ACTIVE'),
  ('166BLNG1', 'NCS Lugogo Hostel Block (Athlete Residence)', 'RESIDENTIAL BUILDINGS', 167250000.00, 167250000.00, 15250000.00, 152000000.00, 'GOOD', 'CONFIRMED_ACTIVE'),
  ('NCSUPS014', 'Light ICT Hardware Batch (86 Units)', 'LIGHT ICT HARDWARE', 59834250.00, 59834250.00, 38494150.00, 21340100.00, 'GOOD', 'CONFIRMED_ACTIVE')
ON CONFLICT DO NOTHING;

INSERT INTO financial_efficiency_appraisals (department_name, allocated_budget_ugx, actual_expenditure_ugx, budget_execution_rate, ntr_target_ugx, ntr_collected_ugx, grant_accountability_rate, audit_query_count, financial_grade) VALUES
  ('Technical & Sports Administration', 12000000000.00, 10800000000.00, 90.0, 350000000.00, 380000000.00, 91.8, 0, 'A'),
  ('Engineering & Infrastructure', 6500000000.00, 5850000000.00, 90.0, 0.00, 0.00, 100.0, 0, 'A'),
  ('Human Resources & Administration', 3500000000.00, 3150000000.00, 90.0, 0.00, 0.00, 100.0, 0, 'A'),
  ('Facilities & Venue Management', 1500000000.00, 1350000000.00, 90.0, 1100000000.00, 1120000000.00, 100.0, 0, 'A')
ON CONFLICT DO NOTHING;
