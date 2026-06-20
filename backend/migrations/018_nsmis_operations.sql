-- NSMIS operational services and history required by production workflows.
CREATE TABLE IF NOT EXISTS reporting_schedule_rules (
 id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, report_type TEXT NOT NULL,
 cadence TEXT NOT NULL CHECK(cadence IN('MONTHLY','QUARTERLY','ANNUAL','EVENT','INTERNATIONAL_EVENT')),
 due_offset_days INTEGER NOT NULL DEFAULT 0 CHECK(due_offset_days>=0), due_day INTEGER CHECK(due_day BETWEEN 1 AND 31),
 reminder_days INTEGER[] NOT NULL DEFAULT ARRAY[14,7,3,1,0,-1,-3,-7,-14], effective_from DATE NOT NULL,
 effective_to DATE, version INTEGER NOT NULL CHECK(version>0), is_active BOOLEAN NOT NULL DEFAULT TRUE,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), UNIQUE(report_type,cadence,version), CHECK(effective_to IS NULL OR effective_to>=effective_from));
INSERT INTO reporting_schedule_rules(report_type,cadence,due_offset_days,due_day,effective_from,version) VALUES
 ('GOVERNANCE','MONTHLY',0,5,CURRENT_DATE,1),('GOVERNANCE','QUARTERLY',0,10,CURRENT_DATE,1),
 ('GOVERNANCE','ANNUAL',0,31,CURRENT_DATE,1),('COMPETITION','EVENT',7,NULL,CURRENT_DATE,1),
 ('COMPETITION','INTERNATIONAL_EVENT',14,NULL,CURRENT_DATE,1) ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS notification_templates (
 id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, code TEXT NOT NULL,
 channel TEXT NOT NULL DEFAULT 'EMAIL' CHECK(channel IN('EMAIL','IN_APP')), subject_template TEXT NOT NULL,
 body_template TEXT NOT NULL, version INTEGER NOT NULL DEFAULT 1, is_active BOOLEAN NOT NULL DEFAULT TRUE,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), UNIQUE(code,channel,version));
CREATE TABLE IF NOT EXISTS notifications (
 id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, template_code TEXT NOT NULL,
 recipient_user_id TEXT REFERENCES users(id) ON DELETE SET NULL, recipient_address TEXT NOT NULL DEFAULT '',
 channel TEXT NOT NULL CHECK(channel IN('EMAIL','IN_APP')), related_type TEXT NOT NULL DEFAULT '', related_id TEXT NOT NULL DEFAULT '',
 payload JSONB NOT NULL DEFAULT '{}'::JSONB, status TEXT NOT NULL DEFAULT 'PENDING'
 CHECK(status IN('PENDING','SENDING','DELIVERED','FAILED','DEAD_LETTER','SUPPRESSED')),
 attempt_count INTEGER NOT NULL DEFAULT 0, next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 delivered_at TIMESTAMPTZ, last_error TEXT NOT NULL DEFAULT '', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE INDEX IF NOT EXISTS idx_notifications_delivery ON notifications(status,next_attempt_at);

CREATE TABLE IF NOT EXISTS background_jobs (
 id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, job_type TEXT NOT NULL, idempotency_key TEXT NOT NULL UNIQUE,
 payload JSONB NOT NULL DEFAULT '{}'::JSONB, status TEXT NOT NULL DEFAULT 'PENDING'
 CHECK(status IN('PENDING','RUNNING','SUCCEEDED','FAILED','DEAD_LETTER','CANCELLED')),
 attempt_count INTEGER NOT NULL DEFAULT 0, max_attempts INTEGER NOT NULL DEFAULT 5,
 next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), leased_until TIMESTAMPTZ, leased_by TEXT NOT NULL DEFAULT '', heartbeat_at TIMESTAMPTZ,
 last_error TEXT NOT NULL DEFAULT '', result JSONB NOT NULL DEFAULT '{}'::JSONB, created_by TEXT REFERENCES users(id) ON DELETE SET NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), started_at TIMESTAMPTZ, finished_at TIMESTAMPTZ, updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE INDEX IF NOT EXISTS idx_background_jobs_lease ON background_jobs(status,next_attempt_at,leased_until);

CREATE TABLE IF NOT EXISTS data_quality_issues (
 id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, federation_id TEXT REFERENCES federations(id) ON DELETE RESTRICT,
 record_type TEXT NOT NULL, record_id TEXT NOT NULL, rule_code TEXT NOT NULL,
 severity TEXT NOT NULL CHECK(severity IN('INFO','WARNING','ERROR','BLOCKING')), description TEXT NOT NULL,
 status TEXT NOT NULL DEFAULT 'OPEN' CHECK(status IN('OPEN','ASSIGNED','RESOLVED','WAIVED')),
 assigned_to TEXT REFERENCES users(id) ON DELETE SET NULL, resolution TEXT NOT NULL DEFAULT '',
 resolved_by TEXT REFERENCES users(id) ON DELETE SET NULL, resolved_at TIMESTAMPTZ,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE INDEX IF NOT EXISTS idx_quality_work_queue ON data_quality_issues(federation_id,status,severity);
CREATE UNIQUE INDEX IF NOT EXISTS uq_open_quality_rule ON data_quality_issues(record_type,record_id,rule_code) WHERE status IN ('OPEN','ASSIGNED');

CREATE TABLE IF NOT EXISTS export_requests (
 id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, export_type TEXT NOT NULL,
 format TEXT NOT NULL CHECK(format IN('CSV','JSON','PDF')), filters JSONB NOT NULL DEFAULT '{}'::JSONB,
 classification TEXT NOT NULL DEFAULT 'CONFIDENTIAL', status TEXT NOT NULL DEFAULT 'PENDING'
 CHECK(status IN('PENDING','PROCESSING','READY','FAILED','EXPIRED','CANCELLED')),
 storage_key TEXT NOT NULL DEFAULT '', sha256 TEXT NOT NULL DEFAULT '', size_bytes BIGINT NOT NULL DEFAULT 0,
 requested_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT, expires_at TIMESTAMPTZ,
 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), completed_at TIMESTAMPTZ);

CREATE TABLE IF NOT EXISTS board_reports (
 id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, reporting_period_id TEXT NOT NULL REFERENCES reporting_periods(id) ON DELETE RESTRICT,
 title TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'DRAFT' CHECK(status IN('DRAFT','APPROVED','PUBLISHED','ARCHIVED')),
 snapshot JSONB NOT NULL DEFAULT '{}'::JSONB, calculation_version INTEGER NOT NULL DEFAULT 1,
 created_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT, approved_by TEXT REFERENCES users(id) ON DELETE SET NULL,
 approved_at TIMESTAMPTZ, published_at TIMESTAMPTZ, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), UNIQUE(reporting_period_id,title));

CREATE TABLE IF NOT EXISTS athlete_status_history (
 id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, athlete_id TEXT NOT NULL REFERENCES athletes(id) ON DELETE RESTRICT,
 status TEXT NOT NULL, effective_from DATE NOT NULL, effective_to DATE, changed_by TEXT REFERENCES users(id) ON DELETE SET NULL,
 reason TEXT NOT NULL DEFAULT '', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), CHECK(effective_to IS NULL OR effective_to>=effective_from));
CREATE TABLE IF NOT EXISTS talent_progression_events (
 id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT, talent_record_id TEXT NOT NULL REFERENCES talent_records(id) ON DELETE RESTRICT,
 from_status TEXT NOT NULL, to_status TEXT NOT NULL, effective_on DATE NOT NULL, notes TEXT NOT NULL DEFAULT '',
 created_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW());

ALTER TABLE athletes ADD COLUMN IF NOT EXISTS consent_basis TEXT NOT NULL DEFAULT 'FEDERATION_MANDATE';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS verified_at TIMESTAMPTZ;
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;
ALTER TABLE competitions ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;
ALTER TABLE medals ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;
ALTER TABLE coaches ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;
ALTER TABLE technical_officials ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;
ALTER TABLE talent_records ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;
ALTER TABLE financial_reports ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;
ALTER TABLE equipment_items ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;
ALTER TABLE talent_scholarships ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE federation_memberships ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE safeguarding_cases ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;
