-- Migration 039: Custom font uploads for the Appearance/Typography settings module.
-- Typography mapping itself (font/size/weight per page element) reuses the
-- existing generic cms_settings key/value store (key='typography'), the same
-- pattern already used for footer/contact/homepage config — no new table
-- needed for that half of the feature.

CREATE TABLE IF NOT EXISTS cms_custom_fonts (
  id           TEXT PRIMARY KEY,
  font_name    TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL,
  file_url     TEXT NOT NULL,
  font_format  TEXT NOT NULL CHECK (font_format IN ('ttf','otf','woff','woff2')),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed an empty typography config so the public site has a record to read
-- (GetSetting already 200s with an empty value on a missing key, but seeding
-- keeps this key visible/discoverable alongside footer/contact/homepage).
INSERT INTO cms_settings (key, value) VALUES
  ('typography', '{}'::jsonb)
ON CONFLICT (key) DO NOTHING;
