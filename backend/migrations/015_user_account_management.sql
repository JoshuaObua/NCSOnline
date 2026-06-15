ALTER TABLE users
  ADD COLUMN IF NOT EXISTS account_status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  ADD COLUMN IF NOT EXISTS status_reason TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS fraud_flag BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS fraud_reason TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS suspended_until TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS status_changed_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS status_changed_by TEXT REFERENCES users(id) ON DELETE SET NULL;

UPDATE users
SET account_status = CASE WHEN is_active THEN 'ACTIVE' ELSE 'SUSPENDED' END
WHERE account_status = 'ACTIVE' AND is_active = FALSE;

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_account_status_check;
ALTER TABLE users
  ADD CONSTRAINT users_account_status_check
  CHECK (account_status IN ('ACTIVE', 'SUSPENDED', 'BANNED'));

CREATE INDEX IF NOT EXISTS idx_users_account_status ON users (account_status);
CREATE INDEX IF NOT EXISTS idx_users_fraud_flag ON users (fraud_flag) WHERE fraud_flag = TRUE;
