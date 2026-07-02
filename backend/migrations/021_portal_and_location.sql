BEGIN;

CREATE TABLE IF NOT EXISTS organisations (
    id TEXT PRIMARY KEY,
    profile_reference VARCHAR(40) NOT NULL UNIQUE,
    organisation_type VARCHAR(30) NOT NULL CHECK (organisation_type IN ('FEDERATION','ASSOCIATION','CLUB')),
    legal_name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    official_email VARCHAR(255) NOT NULL,
    official_phone VARCHAR(40) NOT NULL DEFAULT '',
    registration_number VARCHAR(100) NOT NULL DEFAULT '',
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','SUSPENDED','ARCHIVED')),
    source_application_id TEXT UNIQUE REFERENCES applications(id),
    profile_data JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by TEXT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS organisation_memberships (
    id TEXT PRIMARY KEY,
    organisation_id TEXT NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(40) NOT NULL CHECK (role IN ('OWNER','ADMIN','OFFICER','VIEWER')),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('INVITED','ACTIVE','SUSPENDED','REVOKED')),
    is_primary_owner BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organisation_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_organisation_memberships_user ON organisation_memberships(user_id, status);

CREATE TABLE IF NOT EXISTS organisation_invitations (
    id TEXT PRIMARY KEY,
    organisation_id TEXT NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(40) NOT NULL DEFAULT 'ADMIN',
    token_hash CHAR(64) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING','ACCEPTED','EXPIRED','REVOKED')),
    expires_at TIMESTAMPTZ NOT NULL,
    accepted_by TEXT REFERENCES users(id),
    accepted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE applications ADD COLUMN IF NOT EXISTS provisioning_status VARCHAR(40) NOT NULL DEFAULT 'NOT_REQUIRED';
ALTER TABLE applications ADD COLUMN IF NOT EXISTS provisioned_organisation_id TEXT REFERENCES organisations(id);
ALTER TABLE applications ADD COLUMN IF NOT EXISTS provisioning_error TEXT NOT NULL DEFAULT '';
ALTER TABLE applications ADD COLUMN IF NOT EXISTS provisioning_attempts INTEGER NOT NULL DEFAULT 0;

ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS geo_region VARCHAR(255);
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS geo_latitude DOUBLE PRECISION;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS geo_longitude DOUBLE PRECISION;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS geo_timezone VARCHAR(100);
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS geo_source VARCHAR(50);
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS platform VARCHAR(120);
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS authenticated BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE cms_menus SET items = replace(items::text, '"/apply"', '"/my-portal"')::jsonb
WHERE items::text LIKE '%"/apply"%';
UPDATE cms_settings SET value = replace(value::text, '"/apply"', '"/my-portal"')::jsonb
WHERE value::text LIKE '%"/apply"%';

COMMIT;
