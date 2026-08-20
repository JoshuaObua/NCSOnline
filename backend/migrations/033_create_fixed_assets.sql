-- Migration: 033_create_fixed_assets.sql
-- Creates tables for Fixed Asset Register, Revaluations, Depreciation Logs, and Spot Checks.

CREATE TABLE IF NOT EXISTS fixed_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    interface_line_number VARCHAR(64) UNIQUE NOT NULL,
    asset_book VARCHAR(64) NOT NULL DEFAULT 'NCS FA BOOK',
    asset_number VARCHAR(64) UNIQUE NOT NULL,
    tag_number VARCHAR(128) UNIQUE NOT NULL,
    asset_description TEXT NOT NULL,
    category_segment1 VARCHAR(128) NOT NULL,
    category_segment3 VARCHAR(128) NOT NULL,
    category_segment4 VARCHAR(128) NOT NULL,
    asset_units INT NOT NULL DEFAULT 1,
    fb_cost NUMERIC(18,2) NOT NULL,
    adjusted_cost NUMERIC(18,2) NOT NULL,
    date_placed_in_service DATE NOT NULL DEFAULT '2023-07-01',
    custodian_department VARCHAR(64) NOT NULL DEFAULT 'General Administration',
    location_building VARCHAR(128) NOT NULL DEFAULT 'NCS Lugogo Head Office',
    location_room VARCHAR(64) DEFAULT 'Main Facility',
    depreciation_method VARCHAR(32) NOT NULL DEFAULT 'STRAIGHT_LINE',
    useful_life_years INT NOT NULL DEFAULT 5,
    accumulated_depreciation NUMERIC(18,2) DEFAULT 0.00,
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
    verification_status VARCHAR(32) DEFAULT 'UNVERIFIED',
    last_verified_at TIMESTAMP WITH TIME ZONE,
    last_verified_by TEXT REFERENCES users(id),
    worksheet_source VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS asset_transaction_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID NOT NULL REFERENCES fixed_assets(id) ON DELETE CASCADE,
    transaction_type VARCHAR(32) NOT NULL, -- INITIAL_IMPORT, REVALUATION, DEPRECIATION, TRANSFER, DISPOSAL, WRITE_OFF
    previous_val NUMERIC(18,2),
    new_val NUMERIC(18,2),
    notes TEXT,
    performed_by TEXT REFERENCES users(id),
    approved_by TEXT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fixed_assets_category ON fixed_assets(category_segment3);
CREATE INDEX IF NOT EXISTS idx_fixed_assets_asset_num ON fixed_assets(asset_number);
CREATE INDEX IF NOT EXISTS idx_fixed_assets_tag_num ON fixed_assets(tag_number);
