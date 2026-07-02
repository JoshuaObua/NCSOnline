-- Immediately invalidate access tokens after password reset, deactivation,
-- deletion, or other account security events.
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS auth_invalid_before TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_users_auth_invalid_before
  ON users (auth_invalid_before)
  WHERE auth_invalid_before IS NOT NULL;
