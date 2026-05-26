-- Migration: 004_create_applications
-- Application submissions, attachments, and payment records

CREATE TYPE application_status AS ENUM (
    'DRAFT',
    'PENDING_SIGNATURE',
    'PENDING_PAYMENT',
    'SUBMITTED',
    'UNDER_REVIEW',
    'APPROVED',
    'REJECTED',
    'NEEDS_INFORMATION',
    'RESUBMITTED',
    'PROOF_UPLOADED'
);

CREATE TYPE payment_method AS ENUM (
    'ONLINE',
    'PROOF_UPLOAD'
);

CREATE TYPE payment_status AS ENUM (
    'PENDING',
    'VERIFIED',
    'REJECTED'
);

CREATE TABLE IF NOT EXISTS applications (
    id                   TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    user_id              TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    form_type            TEXT NOT NULL,
    reference_number     TEXT UNIQUE,
    status               application_status NOT NULL DEFAULT 'DRAFT',

    -- Step 1: form data
    form_data            JSONB NOT NULL DEFAULT '{}',

    -- Step 3: signed PDF
    signed_form_url      TEXT,
    signed_at            TIMESTAMPTZ,

    -- Step 4: payment
    payment_method       payment_method,
    payment_status       payment_status,
    payment_amount       NUMERIC(12,2),
    payment_reference    TEXT,
    payment_proof_url    TEXT,
    payment_verified_at  TIMESTAMPTZ,
    payment_gateway_ref  TEXT,

    -- Step 6: admin review
    admin_notes          TEXT,
    reviewed_by          TEXT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at          TIMESTAMPTZ,

    submitted_at         TIMESTAMPTZ,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS application_attachments (
    id             TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    application_id TEXT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    field_key      TEXT NOT NULL,
    file_name      TEXT NOT NULL,
    file_url       TEXT NOT NULL,
    uploaded_by    TEXT REFERENCES users(id) ON DELETE SET NULL,
    uploaded_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_applications_user_id ON applications (user_id);
CREATE INDEX IF NOT EXISTS idx_applications_status ON applications (status);
CREATE INDEX IF NOT EXISTS idx_applications_form_type ON applications (form_type);
CREATE INDEX IF NOT EXISTS idx_applications_reference ON applications (reference_number) WHERE reference_number IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_attachments_application_id ON application_attachments (application_id);
