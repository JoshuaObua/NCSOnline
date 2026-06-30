-- Migration: 004_create_applications

CREATE TABLE IF NOT EXISTS applications (
    id                     TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    user_id                TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    form_type              TEXT NOT NULL,
    application_type       TEXT,
    organisation_type      TEXT,
    application_reference  TEXT UNIQUE,
    status                 TEXT NOT NULL DEFAULT 'DRAFT',
    payment_status         TEXT NOT NULL DEFAULT 'UNPAID',
    payment_method         TEXT,
    payment_reference      TEXT,
    payment_amount_ugx     NUMERIC(14,2),
    payment_verified_by    TEXT REFERENCES users(id) ON DELETE SET NULL,
    payment_verified_at    TIMESTAMPTZ,
    signed_form_url        TEXT,
    signed_form_uploaded_at TIMESTAMPTZ,
    last_saved_step        INT NOT NULL DEFAULT 0,
    draft_expires_at       TIMESTAMPTZ,
    submitted_at           TIMESTAMPTZ,
    reviewer_id            TEXT REFERENCES users(id) ON DELETE SET NULL,
    review_notes           TEXT,
    approved_at            TIMESTAMPTZ,
    rejected_at            TIMESTAMPTZ,
    form_data              JSONB NOT NULL DEFAULT '{}',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS application_attachments (
    id             TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    application_id TEXT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    field_name     TEXT NOT NULL,
    file_name      TEXT NOT NULL,
    file_url       TEXT NOT NULL,
    file_size      INT8 NOT NULL DEFAULT 0,
    mime_type      TEXT NOT NULL DEFAULT '',
    uploaded_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_applications_user_id ON applications (user_id);
CREATE INDEX IF NOT EXISTS idx_applications_status ON applications (status);
CREATE INDEX IF NOT EXISTS idx_applications_form_type ON applications (form_type);
CREATE INDEX IF NOT EXISTS idx_applications_reference ON applications (application_reference) WHERE application_reference IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_attachments_application_id ON application_attachments (application_id);
