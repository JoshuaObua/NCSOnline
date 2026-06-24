-- Migration 028: Content Manager dashboard expansion
--   * Reseed `departments` as the 10-tier institutional matrix used by
--     the team directory and careers module.
--   * Add `department_id` to cms_team_members and cms_careers so staff and
--     job postings are categorised under those 10 units.
--   * Per-user IP whitelist + 2FA columns to power the new Security
--     Settings panel.
--
-- Note: the form-builder feature stores role-ids in its `department_id`
-- columns and no longer references this `departments` table — see
-- migration 027. Reseeding here is safe.

BEGIN;

-- ── 10-tier institutional departments ─────────────────────────────
-- Wipe the old NCS seed and replace with the institutional matrix.
DELETE FROM departments;

INSERT INTO departments (id, name, code, description) VALUES
    ('dept_administration',
        'Administration',
        'ADMINISTRATION',
        'Coordinates overall activities, policy implementation, and secretariat operations under the General Secretary.'),
    ('dept_human_resource',
        'Human Resource',
        'HUMAN_RESOURCE',
        'Manages recruitment, staff welfare, training, performance management, and employee relations.'),
    ('dept_finance_accounts',
        'Finance & Accounts',
        'FINANCE_ACCOUNTS',
        'Responsible for financial planning, budgeting, accounting, internal audit, and fiscal compliance.'),
    ('dept_ict',
        'ICT',
        'ICT',
        'Manages information systems, digital infrastructure, website management, and technical support.'),
    ('dept_sports_officers',
        'Sports Officers',
        'SPORTS_OFFICERS',
        'Coordinates sports programs, athlete development, competition management, and federation liaison.'),
    ('dept_engineering_facilities',
        'Engineering & Facilities',
        'ENGINEERING_FACILITIES',
        'Manages sports facilities infrastructure, civil and electrical engineering, and facility maintenance.'),
    ('dept_public_relations',
        'Public Relations & Communications',
        'PUBLIC_RELATIONS',
        'Handles media relations, marketing, communications, corporate sales, and stakeholder engagement.'),
    ('dept_legal_compliance',
        'Legal & Compliance',
        'LEGAL_COMPLIANCE',
        'Manages legal affairs, licensing, regulatory compliance, and dispute resolution.'),
    ('dept_procurement_records',
        'Procurement & Records',
        'PROCUREMENT_RECORDS',
        'Handles procurement processes, inventory management, records keeping, and document management.'),
    ('dept_support_services',
        'Support Services',
        'SUPPORT_SERVICES',
        'Provides security, transport, office maintenance, and general support services.');

-- ── Team members → department FK ─────────────────────────────────
ALTER TABLE cms_team_members
    ADD COLUMN IF NOT EXISTS department_id TEXT REFERENCES departments(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_cms_team_members_department ON cms_team_members(department_id);

-- ── Careers → department FK (keep legacy `department` text col for now) ─
ALTER TABLE cms_careers
    ADD COLUMN IF NOT EXISTS department_id TEXT REFERENCES departments(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_cms_careers_department ON cms_careers(department_id);

-- ── Per-user IP allowlist ────────────────────────────────────────
CREATE TABLE IF NOT EXISTS user_ip_whitelist (
    id          TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_or_cidr  VARCHAR(64) NOT NULL,
    label       VARCHAR(120) NOT NULL DEFAULT '',
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, ip_or_cidr)
);
CREATE INDEX IF NOT EXISTS idx_user_ip_whitelist_user ON user_ip_whitelist(user_id, is_active);

-- ── 2FA columns on users ─────────────────────────────────────────
ALTER TABLE users ADD COLUMN IF NOT EXISTS twofa_secret      TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS twofa_enabled     BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS twofa_enabled_at  TIMESTAMPTZ;

COMMIT;
