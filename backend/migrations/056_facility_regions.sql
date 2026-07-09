-- Migration: 056_facility_regions
-- Adds CMS-managed regions and complete public facility card metadata.

CREATE TABLE IF NOT EXISTS cms_facility_regions (
  id          TEXT PRIMARY KEY,
  name        TEXT NOT NULL,
  slug        TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  sort_order  INT NOT NULL DEFAULT 0,
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO cms_facility_regions (id, name, slug, description, sort_order, is_active)
VALUES
  ('facility_region_central', 'Central Region', 'central', 'Sports facilities located in Uganda''s Central Region.', 10, TRUE),
  ('facility_region_northern', 'Northern Region', 'northern', 'Sports facilities located in Uganda''s Northern Region.', 20, TRUE),
  ('facility_region_eastern', 'Eastern Region', 'eastern', 'Sports facilities located in Uganda''s Eastern Region.', 30, TRUE),
  ('facility_region_western', 'Western Region', 'western', 'Sports facilities located in Uganda''s Western Region.', 40, TRUE),
  ('facility_region_southern', 'Southern Region', 'southern', 'Sports facilities located in Uganda''s Southern Region.', 50, TRUE)
ON CONFLICT (slug) DO UPDATE
SET name = EXCLUDED.name,
    description = CASE WHEN cms_facility_regions.description = '' THEN EXCLUDED.description ELSE cms_facility_regions.description END,
    sort_order = EXCLUDED.sort_order,
    updated_at = NOW();

ALTER TABLE cms_facilities
  ADD COLUMN IF NOT EXISTS region_slug TEXT NOT NULL DEFAULT 'central',
  ADD COLUMN IF NOT EXISTS location TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS amenities TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS phone TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS email TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS availability_status TEXT NOT NULL DEFAULT 'Available';

ALTER TABLE cms_facilities
  DROP CONSTRAINT IF EXISTS cms_facilities_region_slug_fkey;

ALTER TABLE cms_facilities
  ADD CONSTRAINT cms_facilities_region_slug_fkey
  FOREIGN KEY (region_slug)
  REFERENCES cms_facility_regions(slug)
  ON UPDATE CASCADE
  ON DELETE RESTRICT;

ALTER TABLE cms_facilities
  DROP CONSTRAINT IF EXISTS cms_facilities_availability_status_check;

ALTER TABLE cms_facilities
  ADD CONSTRAINT cms_facilities_availability_status_check
  CHECK (availability_status IN ('Available', 'Limited', 'Maintenance', 'Unavailable'));

CREATE INDEX IF NOT EXISTS idx_cms_facilities_region
  ON cms_facilities(region_slug, sort_order, name);

UPDATE cms_facilities
SET location = CASE
      WHEN location <> '' THEN location
      WHEN slug IN ('cricket-oval', 'hockey-pitch', 'volleyball-courts') THEN 'Lugogo, Kampala'
      ELSE 'Lugogo Sports Complex, Kampala'
    END,
    phone = CASE WHEN phone = '' THEN '+256414254477' ELSE phone END,
    region_slug = 'central'
WHERE region_slug = '' OR region_slug = 'central';

DO $$
DECLARE
  invalid_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO invalid_count
  FROM cms_facilities f
  LEFT JOIN cms_facility_regions r ON r.slug = f.region_slug
  WHERE r.slug IS NULL;

  IF invalid_count <> 0 THEN
    RAISE EXCEPTION 'Found % facilities without a valid region', invalid_count;
  END IF;
END $$;
