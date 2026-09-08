-- Migration: 080_enhance_fixed_assets_audit_and_delete.sql
-- Allows asset_transaction_logs to retain audit history when an asset is deleted,
-- and adds columns for asset snapshot details and actor user names.

ALTER TABLE asset_transaction_logs ALTER COLUMN asset_id DROP NOT NULL;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'asset_transaction_logs_asset_id_fkey'
        AND table_name = 'asset_transaction_logs'
    ) THEN
        ALTER TABLE asset_transaction_logs DROP CONSTRAINT asset_transaction_logs_asset_id_fkey;
    END IF;
END $$;

ALTER TABLE asset_transaction_logs
    ADD CONSTRAINT asset_transaction_logs_asset_id_fkey
    FOREIGN KEY (asset_id) REFERENCES fixed_assets(id) ON DELETE SET NULL;

ALTER TABLE asset_transaction_logs ADD COLUMN IF NOT EXISTS asset_number VARCHAR(64);
ALTER TABLE asset_transaction_logs ADD COLUMN IF NOT EXISTS asset_description TEXT;
ALTER TABLE asset_transaction_logs ADD COLUMN IF NOT EXISTS performed_by_name TEXT;
