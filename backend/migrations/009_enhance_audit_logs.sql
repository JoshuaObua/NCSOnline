-- Migration 009: Comprehensive audit log enhancement
-- Adds geo-location, VPN detection, browser/OS info, event classification,
-- threat scoring, and session tracking to audit_logs.

ALTER TABLE audit_logs
  ADD COLUMN IF NOT EXISTS event_type      TEXT,
  ADD COLUMN IF NOT EXISTS event_status    TEXT,
  ADD COLUMN IF NOT EXISTS severity_level  TEXT,
  ADD COLUMN IF NOT EXISTS forwarded_ip    TEXT,
  ADD COLUMN IF NOT EXISTS geo_country     TEXT,
  ADD COLUMN IF NOT EXISTS geo_city        TEXT,
  ADD COLUMN IF NOT EXISTS vpn_detected    BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS browser         TEXT,
  ADD COLUMN IF NOT EXISTS os_name         TEXT,
  ADD COLUMN IF NOT EXISTS client_type     TEXT,
  ADD COLUMN IF NOT EXISTS threat_score    INT NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS anomaly_detected BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS session_id      TEXT;

-- Indexes for common filter operations on the audit log viewer
CREATE INDEX IF NOT EXISTS idx_audit_logs_event_type     ON audit_logs (event_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_event_status   ON audit_logs (event_status);
CREATE INDEX IF NOT EXISTS idx_audit_logs_severity_level ON audit_logs (severity_level);
CREATE INDEX IF NOT EXISTS idx_audit_logs_geo_country    ON audit_logs (geo_country);
CREATE INDEX IF NOT EXISTS idx_audit_logs_threat_score   ON audit_logs (threat_score);
CREATE INDEX IF NOT EXISTS idx_audit_logs_vpn_detected   ON audit_logs (vpn_detected);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at     ON audit_logs (created_at DESC);
