-- Migration 011: Username in audit logs, career categories, and new CMS modules

-- ── Audit logs: store username for display ─────────────────────────
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS username TEXT;

-- ── Careers: category (jobs | tenders | internships) ──────────────
ALTER TABLE cms_careers ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT 'jobs';

-- ── Fun Facts ──────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cms_fun_facts (
  id         TEXT PRIMARY KEY,
  label      TEXT NOT NULL,
  value      TEXT NOT NULL,
  icon       TEXT NOT NULL DEFAULT '',
  sort_order INT  NOT NULL DEFAULT 0,
  is_active  BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ── FAQs ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cms_faqs (
  id         TEXT PRIMARY KEY,
  question   TEXT NOT NULL,
  answer     TEXT NOT NULL,
  category   TEXT NOT NULL DEFAULT 'general',
  sort_order INT  NOT NULL DEFAULT 0,
  is_active  BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ── Resources (PDFs / Downloads) ──────────────────────────────────
CREATE TABLE IF NOT EXISTS cms_resources (
  id          TEXT PRIMARY KEY,
  title       TEXT NOT NULL,
  category    TEXT NOT NULL DEFAULT 'general',
  file_url    TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  sort_order  INT  NOT NULL DEFAULT 0,
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ── Facilities ─────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cms_facilities (
  id          TEXT PRIMARY KEY,
  name        TEXT NOT NULL,
  slug        TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  image_url   TEXT NOT NULL DEFAULT '',
  sort_order  INT  NOT NULL DEFAULT 0,
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO cms_facilities (id, name, slug, sort_order) VALUES
  (gen_random_uuid()::TEXT, 'Cricket Oval',     'cricket-oval',     1),
  (gen_random_uuid()::TEXT, 'Gymnasium',         'gymnasium',         2),
  (gen_random_uuid()::TEXT, 'Hockey Pitch',      'hockey-pitch',      3),
  (gen_random_uuid()::TEXT, 'Hostel',            'hostel',            4),
  (gen_random_uuid()::TEXT, 'Indoor Stadium',    'indoor-stadium',    5),
  (gen_random_uuid()::TEXT, 'Restaurants',       'restaurants',       6),
  (gen_random_uuid()::TEXT, 'Sports Shop',       'sports-shop',       7),
  (gen_random_uuid()::TEXT, 'Volleyball Courts', 'volleyball-courts', 8)
ON CONFLICT (slug) DO NOTHING;

-- ── Associations ───────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cms_associations (
  id          TEXT PRIMARY KEY,
  name        TEXT NOT NULL,
  slug        TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  logo_url    TEXT NOT NULL DEFAULT '',
  website_url TEXT NOT NULL DEFAULT '',
  sort_order  INT  NOT NULL DEFAULT 0,
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ── Invest with Us ─────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cms_invest (
  id         TEXT PRIMARY KEY,
  title      TEXT NOT NULL,
  subtitle   TEXT NOT NULL DEFAULT '',
  content    TEXT NOT NULL DEFAULT '',
  image_url  TEXT NOT NULL DEFAULT '',
  sort_order INT  NOT NULL DEFAULT 0,
  is_active  BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
