-- Migration: 008_add_pin_cms
-- Adds PIN auth, HTTP audit columns, CMS tables, and content_manager role

-- ── PIN support on users ──────────────────────────────────────────────
ALTER TABLE users ADD COLUMN IF NOT EXISTS pin_hash TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS pin_change_required BOOLEAN NOT NULL DEFAULT TRUE;

-- ── HTTP request audit columns ────────────────────────────────────────
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS method TEXT;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS endpoint TEXT;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS response_code INT;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS response_time_ms BIGINT;
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS device_info TEXT;

-- ── CMS: blog posts / news / case studies / projects ─────────────────
CREATE TABLE IF NOT EXISTS cms_posts (
    id              TEXT PRIMARY KEY,
    title           TEXT NOT NULL,
    slug            TEXT NOT NULL UNIQUE,
    content         TEXT NOT NULL DEFAULT '',
    excerpt         TEXT NOT NULL DEFAULT '',
    category        TEXT NOT NULL DEFAULT 'blog',
    status          TEXT NOT NULL DEFAULT 'draft',
    cover_image_url TEXT NOT NULL DEFAULT '',
    author_id       TEXT REFERENCES users(id) ON DELETE SET NULL,
    published_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_cms_posts_slug     ON cms_posts (slug);
CREATE INDEX IF NOT EXISTS idx_cms_posts_status   ON cms_posts (status);
CREATE INDEX IF NOT EXISTS idx_cms_posts_category ON cms_posts (category);

-- ── CMS: events ───────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cms_events (
    id              TEXT PRIMARY KEY,
    title           TEXT NOT NULL,
    slug            TEXT NOT NULL UNIQUE,
    description     TEXT NOT NULL DEFAULT '',
    location        TEXT NOT NULL DEFAULT '',
    event_date      TIMESTAMPTZ,
    end_date        TIMESTAMPTZ,
    cover_image_url TEXT NOT NULL DEFAULT '',
    status          TEXT NOT NULL DEFAULT 'draft',
    author_id       TEXT REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_cms_events_slug   ON cms_events (slug);
CREATE INDEX IF NOT EXISTS idx_cms_events_status ON cms_events (status);

-- ── CMS: careers ──────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cms_careers (
    id              TEXT PRIMARY KEY,
    title           TEXT NOT NULL,
    department      TEXT NOT NULL DEFAULT '',
    location        TEXT NOT NULL DEFAULT '',
    job_type        TEXT NOT NULL DEFAULT 'full_time',
    description     TEXT NOT NULL DEFAULT '',
    requirements    TEXT NOT NULL DEFAULT '',
    salary_range    TEXT NOT NULL DEFAULT '',
    status          TEXT NOT NULL DEFAULT 'draft',
    deadline_at     TIMESTAMPTZ,
    author_id       TEXT REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_cms_careers_status ON cms_careers (status);

-- ── content_manager role ──────────────────────────────────────────────
INSERT INTO roles (id, name, description, is_system) VALUES
    ('role_content_manager', 'content_manager', 'Manage website content, blog posts, events and careers', TRUE)
ON CONFLICT (name) DO NOTHING;

-- ── CMS permissions ───────────────────────────────────────────────────
INSERT INTO permissions (id, name, description, resource, action) VALUES
    ('perm_cms_read',   'cms:read',   'View CMS content drafts and published',  'cms', 'read'),
    ('perm_cms_write',  'cms:write',  'Create and edit CMS posts/events/jobs',  'cms', 'write'),
    ('perm_cms_delete', 'cms:delete', 'Delete CMS content',                     'cms', 'delete')
ON CONFLICT (name) DO NOTHING;

-- content_manager
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_content_manager', id FROM permissions
WHERE name IN ('cms:read', 'cms:write', 'cms:delete')
ON CONFLICT DO NOTHING;

-- admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions
WHERE name IN ('cms:read', 'cms:write', 'cms:delete')
ON CONFLICT DO NOTHING;

-- super_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions
WHERE name IN ('cms:read', 'cms:write', 'cms:delete')
ON CONFLICT DO NOTHING;
