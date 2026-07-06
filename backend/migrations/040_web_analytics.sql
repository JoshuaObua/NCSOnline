-- Migration: 040_web_analytics
-- Privacy-conscious public website analytics.

CREATE TABLE IF NOT EXISTS page_views_raw (
    id                       TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    session_id               TEXT NOT NULL,
    visitor_hash             TEXT NOT NULL,
    event_name               TEXT NOT NULL DEFAULT 'page_view',
    path                     TEXT NOT NULL,
    referrer                 TEXT NOT NULL DEFAULT '',
    source                   TEXT NOT NULL DEFAULT 'Direct',
    country                  TEXT NOT NULL DEFAULT 'Unknown',
    region                   TEXT NOT NULL DEFAULT '',
    city                     TEXT NOT NULL DEFAULT '',
    browser                  TEXT NOT NULL DEFAULT 'Unknown',
    os                       TEXT NOT NULL DEFAULT 'Unknown',
    device_type              TEXT NOT NULL DEFAULT 'Unknown',
    screen_resolution        TEXT NOT NULL DEFAULT '',
    language                 TEXT NOT NULL DEFAULT '',
    session_duration_seconds INTEGER NOT NULL DEFAULT 0,
    is_bounce                BOOLEAN NOT NULL DEFAULT FALSE,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analytics_snapshots_daily (
    day                      DATE PRIMARY KEY,
    total_views              BIGINT NOT NULL DEFAULT 0,
    unique_visitors          BIGINT NOT NULL DEFAULT 0,
    bounces                  BIGINT NOT NULL DEFAULT 0,
    avg_session_seconds      NUMERIC(10,2) NOT NULL DEFAULT 0,
    top_countries            JSONB NOT NULL DEFAULT '[]'::JSONB,
    top_regions              JSONB NOT NULL DEFAULT '[]'::JSONB,
    top_cities               JSONB NOT NULL DEFAULT '[]'::JSONB,
    top_pages                JSONB NOT NULL DEFAULT '[]'::JSONB,
    top_entry_pages          JSONB NOT NULL DEFAULT '[]'::JSONB,
    top_exit_pages           JSONB NOT NULL DEFAULT '[]'::JSONB,
    browsers                 JSONB NOT NULL DEFAULT '[]'::JSONB,
    operating_systems        JSONB NOT NULL DEFAULT '[]'::JSONB,
    device_types             JSONB NOT NULL DEFAULT '[]'::JSONB,
    acquisition_channels     JSONB NOT NULL DEFAULT '[]'::JSONB,
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_page_views_raw_created_at ON page_views_raw (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_page_views_raw_path ON page_views_raw (path);
CREATE INDEX IF NOT EXISTS idx_page_views_raw_geo ON page_views_raw (country, region, city);
CREATE INDEX IF NOT EXISTS idx_page_views_raw_session ON page_views_raw (session_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_page_views_raw_visitor_day ON page_views_raw (visitor_hash, created_at DESC);

INSERT INTO permissions (id, name, description, resource, action) VALUES
    ('perm_analytics_read', 'analytics:read', 'View public website analytics', 'analytics', 'read')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_super_admin', id FROM permissions WHERE name = 'analytics:read'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions WHERE name = 'analytics:read'
ON CONFLICT DO NOTHING;
