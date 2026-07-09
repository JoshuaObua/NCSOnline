-- Migration: 026_dynamic_application_forms
-- Dynamic Google-Forms-style application builder with departmental tenancy,
-- UGX pricing, DRAFT/OPEN/CLOSED/ARCHIVED lifecycle, and JSONB field schemas.

BEGIN;

-- ── Departments (internal NCS tenancy boundary) ─────────────────────
CREATE TABLE IF NOT EXISTS departments (
    id          TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    name        VARCHAR(120) NOT NULL UNIQUE,
    code        VARCHAR(40)  NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed core NCS departments
INSERT INTO departments (id, name, code, description) VALUES
    (gen_random_uuid()::TEXT, 'Admissions',                'ADMISSIONS',  'Federation, association and club admissions intake'),
    (gen_random_uuid()::TEXT, 'Finance',                   'FINANCE',     'Payments, fees, and financial verification'),
    (gen_random_uuid()::TEXT, 'Technical Department',      'TECHNICAL',   'Technical reviews, sport competitions, academies'),
    (gen_random_uuid()::TEXT, 'Facilities',                'FACILITIES',  'Sports facilities licensing and inspections'),
    (gen_random_uuid()::TEXT, 'General Secretariat',       'SECRETARIAT', 'Cross-cutting governance and oversight')
ON CONFLICT (code) DO NOTHING;

-- Users get a department FK so RBAC can enforce tenancy.
ALTER TABLE users ADD COLUMN IF NOT EXISTS department_id TEXT REFERENCES departments(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_users_department_id ON users(department_id);

-- ── Form Templates ──────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS form_templates (
    id                TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    department_id     TEXT NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    slug              VARCHAR(120) NOT NULL UNIQUE,
    title             VARCHAR(255) NOT NULL,
    description       TEXT NOT NULL DEFAULT '',
    banner_image_url  TEXT NOT NULL DEFAULT '',
    price_ugx         NUMERIC(14,2) NOT NULL DEFAULT 0,
    status            VARCHAR(20) NOT NULL DEFAULT 'DRAFT'
                      CHECK (status IN ('DRAFT','OPEN','CLOSED','ARCHIVED')),
    legacy_form_type  TEXT,
    created_by        TEXT REFERENCES users(id) ON DELETE SET NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_form_templates_department ON form_templates(department_id);
CREATE INDEX IF NOT EXISTS idx_form_templates_status     ON form_templates(status);

-- ── Form Fields (JSONB-driven schema, ordered) ──────────────────────
-- Each row is a single field on a template. `config` holds the type-specific
-- pieces (options for select/radio/checkbox, accepted MIME types, validation).
CREATE TABLE IF NOT EXISTS form_fields (
    id           TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    template_id  TEXT NOT NULL REFERENCES form_templates(id) ON DELETE CASCADE,
    field_key    VARCHAR(80) NOT NULL,
    field_type   VARCHAR(40) NOT NULL CHECK (field_type IN (
                    'short_text','long_text','phone','email','number',
                    'file','image','dropdown','radio','checkbox','date'
                 )),
    label        VARCHAR(255) NOT NULL,
    placeholder  VARCHAR(255) NOT NULL DEFAULT '',
    help_text    TEXT NOT NULL DEFAULT '',
    is_required  BOOLEAN NOT NULL DEFAULT FALSE,
    order_index  INT NOT NULL DEFAULT 0,
    config       JSONB NOT NULL DEFAULT '{}'::jsonb,
    deleted_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (template_id, field_key)
);
CREATE INDEX IF NOT EXISTS idx_form_fields_template_order
    ON form_fields(template_id, order_index)
    WHERE deleted_at IS NULL;

-- ── Form Submissions ────────────────────────────────────────────────
-- Mirrors applications but anchored to a dynamic template.
CREATE TABLE IF NOT EXISTS form_submissions (
    id                    TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    template_id           TEXT NOT NULL REFERENCES form_templates(id) ON DELETE RESTRICT,
    department_id         TEXT NOT NULL REFERENCES departments(id)    ON DELETE RESTRICT,
    user_id               TEXT NOT NULL REFERENCES users(id)          ON DELETE RESTRICT,
    submission_reference  TEXT UNIQUE,
    status                VARCHAR(40) NOT NULL DEFAULT 'DRAFT'
                          CHECK (status IN ('DRAFT','PENDING_PAYMENT','SUBMITTED','UNDER_REVIEW',
                                            'NEEDS_INFORMATION','APPROVED','REJECTED')),
    payment_status        VARCHAR(40) NOT NULL DEFAULT 'UNPAID',
    payment_reference     TEXT,
    payment_amount_ugx    NUMERIC(14,2),
    payment_verified_by   TEXT REFERENCES users(id) ON DELETE SET NULL,
    payment_verified_at   TIMESTAMPTZ,
    answers               JSONB NOT NULL DEFAULT '{}'::jsonb,
    reviewer_id           TEXT REFERENCES users(id) ON DELETE SET NULL,
    review_notes          TEXT,
    submitted_at          TIMESTAMPTZ,
    approved_at           TIMESTAMPTZ,
    rejected_at           TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_form_submissions_template ON form_submissions(template_id);
CREATE INDEX IF NOT EXISTS idx_form_submissions_user     ON form_submissions(user_id);
CREATE INDEX IF NOT EXISTS idx_form_submissions_dept     ON form_submissions(department_id, status);

-- ── Seed: replace hardcoded form_types with dynamic templates ───────
-- One template per legacy form, owned by an appropriate department.
-- We capture the legacy_form_type so the existing applications table can
-- be cross-referenced if needed; this does NOT migrate row data.
DO $$
DECLARE
    d_admissions TEXT;
    d_technical  TEXT;
    d_facilities TEXT;
BEGIN
    SELECT id INTO d_admissions FROM departments WHERE code = 'ADMISSIONS';
    SELECT id INTO d_technical  FROM departments WHERE code = 'TECHNICAL';
    SELECT id INTO d_facilities FROM departments WHERE code = 'FACILITIES';

    INSERT INTO form_templates (department_id, slug, title, description, price_ugx, status, legacy_form_type) VALUES
        (d_admissions, 'form-1-declaration-national-sport', 'Declaration of National Sport',
         'Declare a new sport for national recognition under the NCS Act.', 0, 'OPEN', 'form_1'),
        (d_admissions, 'form-3-register-nsa-nsf', 'Registration / Renewal of NSA or NSF',
         'Register or renew a National Sports Association or National Sports Federation.', 200000, 'OPEN', 'form_3'),
        (d_admissions, 'form-5-transform-nsa-to-nsf', 'Transformation: NSA to NSF',
         'Apply to transform a National Sports Association into a Federation.', 150000, 'OPEN', 'form_5'),
        (d_technical,  'form-7-organise-competition', 'Permit to Organise a Sports Competition',
         'Apply for a permit to organise a national or international sports competition.', 100000, 'OPEN', 'form_7'),
        (d_facilities, 'form-8-operate-facility', 'Licence to Operate a Sports Facility',
         'Apply for a licence to operate a public or commercial sports facility.', 250000, 'OPEN', 'form_8'),
        (d_admissions, 'form-10-community-club', 'Community Sports Club Registration / Renewal',
         'Register or renew a community-level sports club.', 50000, 'OPEN', 'form_10'),
        (d_technical,  'form-11-operate-academy', 'Licence to Operate a Sports Academy',
         'Apply for a licence to run a sports academy or training institution.', 150000, 'OPEN', 'form_11')
    ON CONFLICT (slug) DO NOTHING;
END $$;

COMMIT;
