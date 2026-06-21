-- Authoritative licence records and privacy-safe USSD access audit.
CREATE TABLE IF NOT EXISTS licences (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  licence_number TEXT NOT NULL UNIQUE,
  licence_type TEXT NOT NULL,
  holder_name TEXT NOT NULL,
  source_application_id TEXT REFERENCES applications(id) ON DELETE SET NULL,
  status TEXT NOT NULL DEFAULT 'VALID' CHECK(status IN ('PENDING','VALID','SUSPENDED','REVOKED','CANCELLED','SUPERSEDED')),
  issued_on DATE,
  valid_from DATE NOT NULL,
  expires_on DATE,
  renewable BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK(expires_on IS NULL OR expires_on >= valid_from)
);
CREATE INDEX IF NOT EXISTS idx_licences_public_lookup ON licences(UPPER(licence_number));
CREATE INDEX IF NOT EXISTS idx_licences_expiry ON licences(status,expires_on);

CREATE TABLE IF NOT EXISTS ussd_access_logs (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  session_id TEXT NOT NULL,
  service_code TEXT NOT NULL DEFAULT '',
  action TEXT NOT NULL,
  outcome TEXT NOT NULL,
  phone_hash TEXT NOT NULL DEFAULT '',
  target_hash TEXT NOT NULL DEFAULT '',
  request_id TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ussd_access_session ON ussd_access_logs(session_id,created_at);
CREATE INDEX IF NOT EXISTS idx_ussd_access_abuse ON ussd_access_logs(phone_hash,action,created_at);
