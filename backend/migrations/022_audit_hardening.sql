BEGIN;
SET LOCAL lock_timeout = '5s';
CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS chain_sequence BIGSERIAL;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS previous_hash TEXT;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS entry_hash TEXT;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS payload_excerpt JSONB;
CREATE UNIQUE INDEX IF NOT EXISTS idx_audit_logs_chain_sequence ON audit_logs(chain_sequence);
CREATE UNIQUE INDEX IF NOT EXISTS idx_audit_logs_entry_hash ON audit_logs(entry_hash) WHERE entry_hash IS NOT NULL;

CREATE OR REPLACE FUNCTION ncs_audit_chain_insert() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE prior TEXT;
BEGIN
  PERFORM pg_advisory_xact_lock(hashtext('ncs_audit_chain'));
  SELECT entry_hash INTO prior FROM audit_logs WHERE entry_hash IS NOT NULL ORDER BY chain_sequence DESC LIMIT 1;
  NEW.previous_hash := COALESCE(prior, repeat('0',64));
  NEW.entry_hash := encode(digest(concat_ws('|',NEW.previous_hash,NEW.id,COALESCE(NEW.user_id,''),NEW.action,NEW.resource,COALESCE(NEW.method,''),COALESCE(NEW.endpoint,''),COALESCE(NEW.response_code,0)::TEXT,NEW.created_at::TEXT),'sha256'),'hex');
  RETURN NEW;
END $$;

DROP TRIGGER IF EXISTS trg_audit_chain_insert ON audit_logs;
CREATE TRIGGER trg_audit_chain_insert BEFORE INSERT ON audit_logs FOR EACH ROW EXECUTE FUNCTION ncs_audit_chain_insert();

CREATE OR REPLACE FUNCTION ncs_audit_immutable() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  RAISE EXCEPTION 'audit logs are append-only';
END $$;
DROP TRIGGER IF EXISTS trg_audit_immutable ON audit_logs;
CREATE TRIGGER trg_audit_immutable BEFORE UPDATE OR DELETE ON audit_logs FOR EACH ROW EXECUTE FUNCTION ncs_audit_immutable();
REVOKE UPDATE, DELETE, TRUNCATE ON audit_logs FROM PUBLIC;

COMMIT;
