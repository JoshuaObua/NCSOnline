-- Migration: 061_create_namis_extension_tables
-- Extends the database with NAMIS tables and column extensions for Clubs, Athletes details, Medical/Safeguarding, and National Team tracking.

-- 1. Create Clubs table
CREATE TABLE IF NOT EXISTS clubs (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  name TEXT NOT NULL,
  acronym TEXT NOT NULL,
  contact_person TEXT NOT NULL DEFAULT '',
  email TEXT NOT NULL DEFAULT '',
  phone TEXT NOT NULL DEFAULT '',
  district TEXT NOT NULL DEFAULT '',
  region TEXT NOT NULL DEFAULT '',
  date_founded DATE,
  status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'SUSPENDED')),
  created_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_clubs_acronym_federation ON clubs(LOWER(acronym), federation_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_clubs_federation ON clubs(federation_id) WHERE deleted_at IS NULL;

-- 2. Extend Athletes table
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS national_id_passport TEXT NOT NULL DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS age_category TEXT NOT NULL DEFAULT 'Senior' CHECK (age_category IN ('U10', 'U12', 'U15', 'U17', 'U20', 'Senior'));
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS parent_details JSONB NOT NULL DEFAULT '{}'::JSONB;
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS phone_contact TEXT NOT NULL DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS email_address TEXT NOT NULL DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS next_of_kin TEXT NOT NULL DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS emergency_contact TEXT NOT NULL DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS education_institution TEXT NOT NULL DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS highest_education_level TEXT NOT NULL DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS sports_scholarship_status BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS current_occupation TEXT NOT NULL DEFAULT '';

-- 3. Extend Coaches table
ALTER TABLE coaches ADD COLUMN IF NOT EXISTS phone TEXT NOT NULL DEFAULT '';
ALTER TABLE coaches ADD COLUMN IF NOT EXISTS email TEXT NOT NULL DEFAULT '';

-- 4. Alter Competitions level constraint (allow District up to World Championship)
ALTER TABLE competitions DROP CONSTRAINT IF EXISTS competitions_level_check;
ALTER TABLE competitions ADD CONSTRAINT competitions_level_check CHECK (level IN ('DISTRICT', 'REGIONAL', 'NATIONAL', 'EAST_AFRICAN', 'AFRICAN', 'COMMONWEALTH', 'OLYMPIC', 'WORLD_CHAMPIONSHIP', 'CONTINENTAL', 'INTERNATIONAL'));

-- 5. Create Medical Records table
CREATE TABLE IF NOT EXISTS medical_records (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  athlete_id TEXT NOT NULL REFERENCES athletes(id) ON DELETE CASCADE,
  blood_group TEXT NOT NULL DEFAULT '',
  allergies TEXT NOT NULL DEFAULT '',
  injury_history TEXT NOT NULL DEFAULT '',
  current_injury_status TEXT NOT NULL DEFAULT 'FIT' CHECK (current_injury_status IN ('FIT', 'INJURED', 'RECOVERING')),
  medical_insurance TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 6. Create Safeguarding Records table
CREATE TABLE IF NOT EXISTS safeguarding_records (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  athlete_id TEXT NOT NULL REFERENCES athletes(id) ON DELETE CASCADE,
  guardian_details JSONB NOT NULL DEFAULT '{}'::JSONB,
  manager_details JSONB NOT NULL DEFAULT '{}'::JSONB,
  safeguarding_officer_assigned TEXT REFERENCES users(id) ON DELETE SET NULL,
  consent_forms_url TEXT NOT NULL DEFAULT '',
  anti_doping_education_completed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 7. Create Anti Doping Compliance table
CREATE TABLE IF NOT EXISTS anti_doping_compliance (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  athlete_id TEXT NOT NULL REFERENCES athletes(id) ON DELETE CASCADE,
  testing_status TEXT NOT NULL DEFAULT 'NOT_TESTED' CHECK (testing_status IN ('NOT_TESTED', 'IN_POOL', 'TESTED')),
  last_tested_on DATE,
  last_test_result TEXT NOT NULL DEFAULT 'NEGATIVE',
  wada_education_completed BOOLEAN NOT NULL DEFAULT FALSE,
  suspension_history TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 8. Create National Team Appearances table
CREATE TABLE IF NOT EXISTS national_team_appearances (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  athlete_id TEXT NOT NULL REFERENCES athletes(id) ON DELETE CASCADE,
  team_name TEXT NOT NULL,
  category TEXT NOT NULL CHECK (category IN ('SENIOR', 'DEVELOPMENT', 'JUNIOR')),
  first_call_up_on DATE,
  last_appearance_on DATE,
  appearances_count INTEGER NOT NULL DEFAULT 1 CHECK (appearances_count >= 1),
  notes TEXT NOT NULL DEFAULT '',
  created_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_athlete_national_team ON national_team_appearances(athlete_id, team_name, category);

-- 9. Extend Medals table
ALTER TABLE medals ADD COLUMN IF NOT EXISTS level TEXT NOT NULL DEFAULT 'NATIONAL' CHECK (level IN ('DISTRICT', 'REGIONAL', 'NATIONAL', 'EAST_AFRICAN', 'AFRICAN', 'COMMONWEALTH', 'OLYMPIC', 'WORLD_CHAMPIONSHIP', 'CONTINENTAL', 'INTERNATIONAL'));
ALTER TABLE medals ADD COLUMN IF NOT EXISTS coach_at_win_id TEXT REFERENCES coaches(id) ON DELETE SET NULL;
ALTER TABLE medals ADD COLUMN IF NOT EXISTS is_team_event BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE medals ADD COLUMN IF NOT EXISTS prize_money NUMERIC(18,2) NOT NULL DEFAULT 0.00;
ALTER TABLE medals ADD COLUMN IF NOT EXISTS ncs_recognition_status TEXT NOT NULL DEFAULT 'PENDING' CHECK (ncs_recognition_status IN ('PENDING', 'APPROVED', 'REJECTED'));

-- 10. Seed federations from cms_associations content
INSERT INTO federations (id, name, acronym, ncs_registration_number, recognition_status)
SELECT id, name, abbreviation, 'NCS-REG-' || abbreviation, 'RECOGNISED'
FROM cms_associations
ON CONFLICT (id) DO NOTHING;
