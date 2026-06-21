BEGIN;
SET LOCAL lock_timeout = '5s';
CREATE TABLE IF NOT EXISTS system_control (
 singleton BOOLEAN PRIMARY KEY DEFAULT TRUE CHECK(singleton),
 maintenance_mode BOOLEAN NOT NULL DEFAULT FALSE,
 maintenance_reason TEXT NOT NULL DEFAULT '',
 maintenance_expected_end TIMESTAMPTZ,
 global_auth_invalid_before TIMESTAMPTZ,
 cache_generation BIGINT NOT NULL DEFAULT 0,
 changed_by TEXT REFERENCES users(id) ON DELETE SET NULL,
 changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
INSERT INTO system_control(singleton) VALUES(TRUE) ON CONFLICT DO NOTHING;
COMMIT;
