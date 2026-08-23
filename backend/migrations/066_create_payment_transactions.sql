-- Migration: 066_create_payment_transactions
-- Configures allowed payment methods on form templates and establishes payment_transactions ledger with ioTec integration.

-- 1. Add allowed_payment_methods to form_templates
ALTER TABLE form_templates
  ADD COLUMN IF NOT EXISTS allowed_payment_methods JSONB NOT NULL DEFAULT '["OVER_THE_COUNTER", "MOBILE_MONEY"]'::jsonb;

-- 2. Add payment_method to form_submissions
ALTER TABLE form_submissions
  ADD COLUMN IF NOT EXISTS payment_method VARCHAR(50) NOT NULL DEFAULT 'OVER_THE_COUNTER';

-- 3. Create payment_transactions table
CREATE TABLE IF NOT EXISTS payment_transactions (
    id                     TEXT PRIMARY KEY,
    transaction_reference  VARCHAR(100) UNIQUE NOT NULL,
    submission_id          TEXT REFERENCES form_submissions(id) ON DELETE SET NULL,
    template_id            TEXT REFERENCES form_templates(id) ON DELETE SET NULL,
    user_id                TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    payment_method         VARCHAR(50) NOT NULL DEFAULT 'MOBILE_MONEY',
    provider               VARCHAR(50) NOT NULL DEFAULT 'IOTEC',
    provider_request_id    TEXT,
    phone_number           VARCHAR(50),
    amount_ugx             NUMERIC(14,2) NOT NULL DEFAULT 0,
    currency               VARCHAR(10) NOT NULL DEFAULT 'UGX',
    status                 VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    status_message         TEXT,
    raw_response           JSONB,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at           TIMESTAMPTZ
);

-- Indexes for performance and lookup
CREATE INDEX IF NOT EXISTS idx_payment_tx_user_id ON payment_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_tx_submission_id ON payment_transactions(submission_id);
CREATE INDEX IF NOT EXISTS idx_payment_tx_status ON payment_transactions(status);
CREATE INDEX IF NOT EXISTS idx_payment_tx_provider_req ON payment_transactions(provider_request_id);
CREATE INDEX IF NOT EXISTS idx_payment_tx_created_at ON payment_transactions(created_at DESC);
