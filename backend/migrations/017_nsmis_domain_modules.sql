-- NSMIS domain tables for athletes, performance, workforce, talent,
-- safeguarding, financial accountability and equipment management.

CREATE TABLE IF NOT EXISTS athletes (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  athlete_number TEXT NOT NULL UNIQUE,
  full_name TEXT NOT NULL,
  gender TEXT NOT NULL CHECK (gender IN ('MALE','FEMALE','OTHER','NOT_STATED')),
  date_of_birth DATE NOT NULL CHECK (date_of_birth <= CURRENT_DATE),
  district TEXT NOT NULL DEFAULT '', region TEXT NOT NULL DEFAULT '', club TEXT NOT NULL DEFAULT '', discipline TEXT NOT NULL,
  national_team_status TEXT NOT NULL DEFAULT 'NO' CHECK (national_team_status IN ('NO','DEVELOPMENT','SENIOR','FORMER')),
  status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','INACTIVE','SUSPENDED','RETIRED')),
  created_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_athletes_dimensions ON athletes(gender,region,discipline,status) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS athlete_affiliations (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  athlete_id TEXT NOT NULL REFERENCES athletes(id) ON DELETE RESTRICT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  starts_on DATE NOT NULL, ends_on DATE, is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), CHECK (ends_on IS NULL OR ends_on >= starts_on)
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_active_athlete_affiliation ON athlete_affiliations(athlete_id,federation_id) WHERE is_active AND ends_on IS NULL;

CREATE TABLE IF NOT EXISTS competitions (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  name TEXT NOT NULL, venue TEXT NOT NULL, host_country TEXT NOT NULL DEFAULT 'Uganda',
  level TEXT NOT NULL CHECK (level IN ('NATIONAL','REGIONAL','CONTINENTAL','INTERNATIONAL')),
  starts_on DATE NOT NULL, ends_on DATE NOT NULL, returned_on DATE,
  athlete_count INTEGER NOT NULL DEFAULT 0 CHECK (athlete_count >= 0), official_count INTEGER NOT NULL DEFAULT 0 CHECK (official_count >= 0),
  status TEXT NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT','SUBMITTED','APPROVED','CANCELLED')),
  created_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), CHECK (ends_on >= starts_on)
);
CREATE INDEX IF NOT EXISTS idx_competitions_dimensions ON competitions(federation_id,level,starts_on,status);

CREATE TABLE IF NOT EXISTS competition_results (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  competition_id TEXT NOT NULL REFERENCES competitions(id) ON DELETE RESTRICT,
  athlete_id TEXT REFERENCES athletes(id) ON DELETE RESTRICT,
  athlete_name TEXT NOT NULL, event TEXT NOT NULL, position INTEGER CHECK (position IS NULL OR position > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), UNIQUE(competition_id,athlete_name,event)
);

CREATE TABLE IF NOT EXISTS medals (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  competition_id TEXT NOT NULL REFERENCES competitions(id) ON DELETE RESTRICT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  athlete_id TEXT REFERENCES athletes(id) ON DELETE RESTRICT,
  athlete_name TEXT NOT NULL, event TEXT NOT NULL, country TEXT NOT NULL, won_on DATE NOT NULL,
  medal_type TEXT NOT NULL CHECK (medal_type IN ('GOLD','SILVER','BRONZE')),
  coach_responsible TEXT NOT NULL DEFAULT '', team_manager TEXT NOT NULL DEFAULT '', funding_source TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'SUBMITTED' CHECK (status IN ('DRAFT','SUBMITTED','APPROVED','VOIDED')),
  created_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(competition_id,athlete_name,event,medal_type)
);
CREATE INDEX IF NOT EXISTS idx_medals_dashboard ON medals(federation_id,medal_type,won_on,status);

CREATE TABLE IF NOT EXISTS coaches (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  full_name TEXT NOT NULL, certification_level TEXT NOT NULL, license_number TEXT NOT NULL, expiry_date DATE,
  status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','INACTIVE','SUSPENDED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(federation_id,license_number)
);
CREATE INDEX IF NOT EXISTS idx_coaches_expiry ON coaches(expiry_date,status);

CREATE TABLE IF NOT EXISTS technical_officials (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  full_name TEXT NOT NULL, official_type TEXT NOT NULL CHECK (official_type IN ('REFEREE','UMPIRE','JUDGE','ASSESSOR','OTHER')),
  level TEXT NOT NULL, certification TEXT NOT NULL, valid_until DATE,
  status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','INACTIVE','SUSPENDED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_officials_validity ON technical_officials(valid_until,status);

CREATE TABLE IF NOT EXISTS talent_records (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  athlete_id TEXT REFERENCES athletes(id) ON DELETE SET NULL,
  athlete_name TEXT NOT NULL, age_at_identification INTEGER NOT NULL CHECK (age_at_identification BETWEEN 4 AND 100),
  school TEXT NOT NULL DEFAULT '', district TEXT NOT NULL, region TEXT NOT NULL DEFAULT '', identified_by TEXT NOT NULL,
  identified_on DATE NOT NULL, status TEXT NOT NULL DEFAULT 'IDENTIFIED'
    CHECK (status IN ('IDENTIFIED','ASSESSED','SELECTED','NATIONAL_TEAM','CLOSED')),
  created_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_talent_dashboard ON talent_records(federation_id,status,identified_on,region);

CREATE TABLE IF NOT EXISTS talent_scholarships (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  talent_record_id TEXT NOT NULL REFERENCES talent_records(id) ON DELETE RESTRICT,
  provider TEXT NOT NULL, starts_on DATE NOT NULL, ends_on DATE, status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK(status IN ('PLANNED','ACTIVE','COMPLETED','CANCELLED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), CHECK(ends_on IS NULL OR ends_on >= starts_on)
);

CREATE TABLE IF NOT EXISTS safeguarding_period_metrics (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  reporting_period_id TEXT NOT NULL REFERENCES reporting_periods(id) ON DELETE RESTRICT,
  female_athletes INTEGER NOT NULL DEFAULT 0 CHECK(female_athletes >= 0), female_coaches INTEGER NOT NULL DEFAULT 0 CHECK(female_coaches >= 0),
  safeguarding_cases INTEGER NOT NULL DEFAULT 0 CHECK(safeguarding_cases >= 0), resolved_cases INTEGER NOT NULL DEFAULT 0 CHECK(resolved_cases >= 0),
  status TEXT NOT NULL DEFAULT 'DRAFT' CHECK(status IN ('DRAFT','SUBMITTED','APPROVED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(federation_id,reporting_period_id), CHECK(resolved_cases <= safeguarding_cases)
);

CREATE TABLE IF NOT EXISTS safeguarding_cases (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  case_reference TEXT NOT NULL UNIQUE, federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  encrypted_details TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'OPEN' CHECK(status IN ('OPEN','UNDER_REVIEW','REFERRED','RESOLVED','CLOSED')),
  reported_at TIMESTAMPTZ NOT NULL, resolved_at TIMESTAMPTZ,
  assigned_to TEXT REFERENCES users(id) ON DELETE SET NULL, created_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ncs_disbursements (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  reference TEXT NOT NULL UNIQUE, amount NUMERIC(18,2) NOT NULL CHECK(amount >= 0), currency CHAR(3) NOT NULL DEFAULT 'UGX',
  released_on DATE NOT NULL, accountability_due_on DATE NOT NULL,
  status TEXT NOT NULL DEFAULT 'RELEASED' CHECK(status IN ('PLANNED','RELEASED','PARTIALLY_ACCOUNTED','ACCOUNTED','CANCELLED')),
  created_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK(accountability_due_on >= released_on)
);
CREATE INDEX IF NOT EXISTS idx_disbursements_dashboard ON ncs_disbursements(federation_id,status,released_on,accountability_due_on);

CREATE TABLE IF NOT EXISTS financial_reports (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  reporting_period_id TEXT NOT NULL REFERENCES reporting_periods(id) ON DELETE RESTRICT,
  currency CHAR(3) NOT NULL DEFAULT 'UGX',
  government_grant NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK(government_grant >= 0), sponsorship NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK(sponsorship >= 0),
  membership_fees NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK(membership_fees >= 0), donations NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK(donations >= 0),
  competitions NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK(competitions >= 0), training NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK(training >= 0),
  equipment NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK(equipment >= 0), administration NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK(administration >= 0),
  status TEXT NOT NULL DEFAULT 'DRAFT' CHECK(status IN ('DRAFT','SUBMITTED','NEEDS_CORRECTION','APPROVED','LOCKED')),
  created_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(federation_id,reporting_period_id)
);

CREATE TABLE IF NOT EXISTS equipment_items (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  federation_id TEXT NOT NULL REFERENCES federations(id) ON DELETE RESTRICT,
  item_name TEXT NOT NULL, unit TEXT NOT NULL DEFAULT 'ITEM', quantity_received INTEGER NOT NULL DEFAULT 0 CHECK(quantity_received >= 0),
  quantity_distributed INTEGER NOT NULL DEFAULT 0 CHECK(quantity_distributed >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(federation_id,item_name,unit), CHECK(quantity_distributed <= quantity_received)
);

CREATE TABLE IF NOT EXISTS equipment_distributions (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  equipment_item_id TEXT NOT NULL REFERENCES equipment_items(id) ON DELETE RESTRICT,
  quantity INTEGER NOT NULL CHECK(quantity > 0), beneficiary TEXT NOT NULL, distributed_on DATE NOT NULL,
  acknowledgement_document_id TEXT REFERENCES federation_documents(id) ON DELETE SET NULL,
  distributed_by TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

