-- Migration 013: CMS team members table

CREATE TABLE IF NOT EXISTS cms_team_members (
  id          TEXT PRIMARY KEY,
  full_name   TEXT    NOT NULL,
  designation TEXT    NOT NULL DEFAULT '',
  image_url   TEXT    NOT NULL DEFAULT '',
  bio         TEXT    NOT NULL DEFAULT '',
  sort_order  INT     NOT NULL DEFAULT 0,
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_cms_team_members_sort ON cms_team_members(sort_order ASC, created_at ASC);
