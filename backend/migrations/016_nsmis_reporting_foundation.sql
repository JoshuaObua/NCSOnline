-- NSMIS reporting foundation: organisations, periods, obligations, reports and governance.
-- Additive and idempotent so existing NCSMS application/CMS data remains intact.

INSERT INTO roles (id, name, description, is_system) VALUES
  ('role_ncs_general_secretary', 'ncs_general_secretary', 'Cross-federation official dashboards and board reporting', TRUE),
  ('role_technical_department', 'technical_department', 'Federation monitoring and sports performance management', TRUE),
  ('role_finance_department', 'finance_department', 'Financial compliance and accountability management', TRUE),
  ('role_federation_president', 'federation_president', 'Review and approve own federation reports', TRUE),
  ('role_federation_general_secretary', 'federation_general_secretary', 'Enter and submit own federation reports', TRUE),
  ('role_safeguarding_officer', 'safeguarding_officer', 'Restricted safeguarding case management', TRUE),
  ('role_auditor', 'auditor', 'Read-only assigned-scope audit access', TRUE)
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

INSERT INTO permissions (id, name, description, resource, action) VALUES
  ('perm_federations_read_any', 'federations:read:any', 'Read all federation profiles', 'federations', 'read:any'),
  ('perm_federations_read_own', 'federations:read:own', 'Read assigned federation profiles', 'federations', 'read:own'),
  ('perm_federations_write_any', 'federations:write:any', 'Manage all federation profiles', 'federations', 'write:any'),
  ('perm_federations_write_own', 'federations:write:own', 'Manage assigned federation profiles', 'federations', 'write:own'),
  ('perm_reports_read_any', 'reports:read:any', 'Read all federation reports', 'reports', 'read:any'),
  ('perm_reports_read_own', 'reports:read:own', 'Read assigned federation reports', 'reports', 'read:own'),
  ('perm_reports_write_own', 'reports:write:own', 'Create and update assigned federation reports', 'reports', 'write:own'),
  ('perm_reports_approve_own', 'reports:approve:own', 'Approve assigned federation reports', 'reports', 'approve:own'),
  ('perm_reports_review_any', 'reports:review:any', 'Review federation reports across NCS', 'reports', 'review:any'),
  ('perm_governance_dashboard', 'dashboard:governance:read', 'View governance dashboard', 'dashboard', 'governance:read'),
  ('perm_athletes_dashboard', 'dashboard:athletes:read', 'View athlete dashboard', 'dashboard', 'athletes:read'),
  ('perm_performance_dashboard', 'dashboard:performance:read', 'View performance dashboard', 'dashboard', 'performance:read'),
  ('perm_finance_dashboard', 'dashboard:finance:read', 'View finance dashboard', 'dashboard', 'finance:read'),
  ('perm_talent_dashboard', 'dashboard:talent:read', 'View talent dashboard', 'dashboard', 'talent:read'),
  ('perm_safeguarding_cases', 'safeguarding:cases:manage', 'Manage restricted safeguarding cases', 'safeguarding', 'cases:manage'),
  ('perm_exports_create', 'exports:create', 'Create controlled data exports', 'exports', 'create')
ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description;

-- Super administrators receive every newly introduced permission.
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions WHERE name LIKE 'federations:%' OR name LIKE 'reports:%'
  OR name LIKE 'dashboard:%' OR name LIKE 'safeguarding:%' OR name = 'exports:create'
ON CONFLICT DO NOTHING;

-- Existing administrator is the operational NCS Administrator.
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions WHERE name LIKE 'federations:%' OR name LIKE 'reports:%'
  OR name LIKE 'dashboard:%' OR name = 'exports:create'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_ncs_general_secretary', id FROM permissions WHERE name IN
  ('federations:read:any','reports:read:any','dashboard:governance:read','dashboard:athletes:read',
   'dashboard:performance:read','dashboard:finance:read','dashboard:talent:read','exports:create')
ON CONFLICT DO NOTHING;

-- Preserve access for installations already using the original role name.
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_general_secretary', id FROM permissions WHERE name IN
  ('federations:read:any','reports:read:any','dashboard:governance:read','dashboard:athletes:read',
   'dashboard:performance:read','dashboard:finance:read','dashboard:talent:read','exports:create')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_technical_department', id FROM permissions WHERE name IN
  ('federations:read:any','reports:read:any','reports:review:any','dashboard:governance:read',
   'dashboard:athletes:read','dashboard:performance:read','dashboard:talent:read','exports:create')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_finance_department', id FROM permissions WHERE name IN
  ('federations:read:any','reports:read:any','reports:review:any','dashboard:finance:read','exports:create')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_federation_president', id FROM permissions WHERE name IN
  ('federations:read:own','reports:read:own','reports:approve:own','dashboard:governance:read')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_federation_general_secretary', id FROM permissions WHERE name IN
  ('federations:read:own','federations:write:own','reports:read:own','reports:write:own','dashboard:governance:read')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_safeguarding_officer', id FROM permissions WHERE name IN
  ('federations:read:any','reports:read:any','safeguarding:cases:manage')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_auditor', id FROM permissions WHERE name IN
  ('federations:read:own','reports:read:own','dashboard:governance:read','exports:create')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS federations (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  name TEXT NOT NULL,
  acronym TEXT NOT NULL,
  ncs_registration_number TEXT NOT NULL,
  recognition_status TEXT NOT NULL DEFAULT 'PENDING'
    CHECK (recognition_status IN ('PENDING','RECOGNISED','SUSPENDED','REVOKED')),
  physical_address TEXT NOT NULL DEFAULT '',
  email TEXT NOT NULL DEFAULT '',
  website TEXT NOT NULL DEFAULT '',
  contact_person TEXT NOT NULL DEFAULT '',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
  created_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  updated_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_federations_registration_active ON federations (LOWER(ncs_registration_number)) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_federations_acronym_active ON federations (LOWER(acronym)) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_federations_status ON federations (recognition_status) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS federation_memberships (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  membership_role TEXT NOT NULL CHECK (membership_role IN ('PRESIDENT','GENERAL_SECRETARY','OFFICER','AUDITOR')),
  starts_at DATE NOT NULL DEFAULT CURRENT_DATE,
  ends_at DATE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  assigned_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (ends_at IS NULL OR ends_at >= starts_at)
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_federation_active_membership ON federation_memberships (federation_id,user_id,membership_role) WHERE is_active AND ends_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_federation_memberships_user ON federation_memberships (user_id) WHERE is_active;

CREATE TABLE IF NOT EXISTS federation_officers (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  position TEXT NOT NULL CHECK (position IN ('PRESIDENT','VICE_PRESIDENT','GENERAL_SECRETARY','TREASURER','TECHNICAL_DIRECTOR','OTHER')),
  position_label TEXT NOT NULL DEFAULT '',
  full_name TEXT NOT NULL,
  email TEXT NOT NULL DEFAULT '',
  phone TEXT NOT NULL DEFAULT '',
  appointed_on DATE,
  term_ends_on DATE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (term_ends_on IS NULL OR appointed_on IS NULL OR term_ends_on >= appointed_on)
);
CREATE INDEX IF NOT EXISTS idx_federation_officers_federation ON federation_officers (federation_id,is_active);

CREATE TABLE IF NOT EXISTS reporting_periods (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  period_type TEXT NOT NULL CHECK (period_type IN ('MONTHLY','QUARTERLY','ANNUAL','COMPETITION','INTERNATIONAL_EVENT')),
  name TEXT NOT NULL,
  starts_on DATE NOT NULL,
  ends_on DATE NOT NULL,
  due_on DATE NOT NULL,
  grace_ends_on DATE,
  timezone TEXT NOT NULL DEFAULT 'Africa/Kampala',
  is_open BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (ends_on >= starts_on),
  CHECK (due_on >= ends_on)
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_reporting_period ON reporting_periods(period_type,starts_on,ends_on);
CREATE INDEX IF NOT EXISTS idx_reporting_period_due ON reporting_periods(due_on,is_open);

CREATE TABLE IF NOT EXISTS report_obligations (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  reporting_period_id TEXT NOT NULL REFERENCES reporting_periods(id) ON DELETE RESTRICT,
  report_type TEXT NOT NULL CHECK (report_type IN ('GOVERNANCE','ATHLETE','COMPETITION','MEDAL','WORKFORCE','TALENT','SAFEGUARDING','FINANCE','EQUIPMENT')),
  due_on DATE NOT NULL,
  status TEXT NOT NULL DEFAULT 'NOT_STARTED'
    CHECK (status IN ('NOT_STARTED','DRAFT','SUBMITTED','ACCEPTED','OVERDUE','EXEMPT')),
  exempt_reason TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (federation_id,reporting_period_id,report_type)
);
CREATE INDEX IF NOT EXISTS idx_obligations_due_status ON report_obligations(due_on,status);
CREATE INDEX IF NOT EXISTS idx_obligations_federation ON report_obligations(federation_id,status);

CREATE TABLE IF NOT EXISTS federation_reports (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  obligation_id TEXT NOT NULL REFERENCES report_obligations(id) ON DELETE RESTRICT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  reporting_period_id TEXT NOT NULL REFERENCES reporting_periods(id) ON DELETE RESTRICT,
  report_type TEXT NOT NULL,
  revision INTEGER NOT NULL DEFAULT 1 CHECK (revision > 0),
  status TEXT NOT NULL DEFAULT 'DRAFT'
    CHECK (status IN ('DRAFT','SUBMITTED','PRESIDENT_APPROVED','PRESIDENT_RETURNED','NCS_UNDER_REVIEW','NEEDS_CORRECTION','APPROVED','LOCKED')),
  data JSONB NOT NULL DEFAULT '{}'::JSONB,
  submitted_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  submitted_at TIMESTAMPTZ,
  president_approved_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  president_approved_at TIMESTAMPTZ,
  reviewed_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  reviewed_at TIMESTAMPTZ,
  locked_at TIMESTAMPTZ,
  version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
  created_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (obligation_id,revision)
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_current_report_revision ON federation_reports(obligation_id) WHERE status NOT IN ('NEEDS_CORRECTION','PRESIDENT_RETURNED');
CREATE INDEX IF NOT EXISTS idx_reports_scope ON federation_reports(federation_id,reporting_period_id,report_type,status);

CREATE TABLE IF NOT EXISTS report_transitions (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  report_id TEXT NOT NULL REFERENCES federation_reports(id) ON DELETE RESTRICT,
  from_status TEXT NOT NULL,
  to_status TEXT NOT NULL,
  actor_id TEXT REFERENCES users(id) ON DELETE SET NULL,
  reason TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_report_transitions_report ON report_transitions(report_id,created_at);

CREATE TABLE IF NOT EXISTS governance_responses (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  report_id TEXT NOT NULL UNIQUE REFERENCES federation_reports(id) ON DELETE RESTRICT,
  executive_meeting_held BOOLEAN NOT NULL DEFAULT FALSE,
  executive_meeting_date DATE,
  agm_conducted BOOLEAN NOT NULL DEFAULT FALSE,
  agm_date DATE,
  board_meeting_held BOOLEAN NOT NULL DEFAULT FALSE,
  board_meeting_date DATE,
  elections_conducted BOOLEAN NOT NULL DEFAULT FALSE,
  election_date DATE,
  disciplinary_cases_handled INTEGER NOT NULL DEFAULT 0 CHECK (disciplinary_cases_handled >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (NOT executive_meeting_held OR executive_meeting_date IS NOT NULL),
  CHECK (NOT agm_conducted OR agm_date IS NOT NULL),
  CHECK (NOT board_meeting_held OR board_meeting_date IS NOT NULL),
  CHECK (NOT elections_conducted OR election_date IS NOT NULL)
);

CREATE TABLE IF NOT EXISTS federation_documents (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  report_id TEXT REFERENCES federation_reports(id) ON DELETE RESTRICT,
  document_type TEXT NOT NULL,
  original_name TEXT NOT NULL,
  storage_key TEXT NOT NULL UNIQUE,
  mime_type TEXT NOT NULL,
  size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),
  sha256 TEXT NOT NULL,
  classification TEXT NOT NULL DEFAULT 'CONFIDENTIAL'
    CHECK (classification IN ('INTERNAL','CONFIDENTIAL','RESTRICTED','HIGHLY_RESTRICTED')),
  scan_status TEXT NOT NULL DEFAULT 'PENDING'
    CHECK (scan_status IN ('PENDING','CLEAN','REJECTED','FAILED')),
  effective_on DATE,
  expires_on DATE,
  uploaded_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (expires_on IS NULL OR effective_on IS NULL OR expires_on >= effective_on)
);
CREATE INDEX IF NOT EXISTS idx_federation_documents_owner ON federation_documents(federation_id,document_type,expires_on);

CREATE TABLE IF NOT EXISTS compliance_rule_versions (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  version INTEGER NOT NULL UNIQUE CHECK (version > 0),
  effective_from DATE NOT NULL,
  compliant_threshold NUMERIC(5,2) NOT NULL DEFAULT 70 CHECK (compliant_threshold BETWEEN 0 AND 100),
  rules JSONB NOT NULL DEFAULT '{}'::JSONB,
  is_active BOOLEAN NOT NULL DEFAULT FALSE,
  approved_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_active_compliance_rules ON compliance_rule_versions(is_active) WHERE is_active;

INSERT INTO compliance_rule_versions(id,version,effective_from,compliant_threshold,rules,is_active)
VALUES ('compliance_v1',1,CURRENT_DATE,70,
  '{"timeliness":40,"governance":30,"document_validity":20,"data_quality":10}'::JSONB,TRUE)
ON CONFLICT (version) DO NOTHING;

CREATE TABLE IF NOT EXISTS compliance_scores (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  reporting_period_id TEXT NOT NULL REFERENCES reporting_periods(id) ON DELETE RESTRICT,
  rule_version_id TEXT NOT NULL REFERENCES compliance_rule_versions(id) ON DELETE RESTRICT,
  score NUMERIC(5,2) NOT NULL CHECK (score BETWEEN 0 AND 100),
  is_compliant BOOLEAN NOT NULL,
  breakdown JSONB NOT NULL DEFAULT '{}'::JSONB,
  calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (federation_id,reporting_period_id,rule_version_id)
);
CREATE INDEX IF NOT EXISTS idx_compliance_period_score ON compliance_scores(reporting_period_id,is_compliant,score DESC);
