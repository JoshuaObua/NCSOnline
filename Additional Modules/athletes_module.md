# NAMIS Athletes Module Implementation Plan

The Athletes Module tracks the Athlete Master Profile, athlete classification, dual careers, safeguarding records, medical information, and anti-doping compliance details.

## 1. Database Schema Plan

We leverage and extend the existing `athletes` table structure. Some details (like medical and safeguarding) require restricted access controls.

```sql
-- Athletes table extension
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS national_id_passport TEXT DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS age_category TEXT DEFAULT 'Senior' CHECK (age_category IN ('U10', 'U12', 'U15', 'U17', 'U20', 'Senior'));
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS parent_details JSONB DEFAULT '{}'::JSONB; -- parent names and contacts
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS phone_contact TEXT DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS email_address TEXT DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS next_of_kin TEXT DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS emergency_contact TEXT DEFAULT '';

-- Dual Career details
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS education_institution TEXT DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS highest_education_level TEXT DEFAULT '';
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS sports_scholarship_status BOOLEAN DEFAULT FALSE;
ALTER TABLE athletes ADD COLUMN IF NOT EXISTS current_occupation TEXT DEFAULT '';

-- Medical information table (RESTRICTED ACCESS)
CREATE TABLE IF NOT EXISTS medical_records (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  athlete_id TEXT NOT NULL REFERENCES athletes(id) ON DELETE CASCADE,
  blood_group TEXT DEFAULT '',
  allergies TEXT DEFAULT '',
  injury_history TEXT DEFAULT '',
  current_injury_status TEXT DEFAULT 'FIT' CHECK (current_injury_status IN ('FIT', 'INJURED', 'RECOVERING')),
  medical_insurance TEXT DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Safeguarding information table (RESTRICTED ACCESS)
CREATE TABLE IF NOT EXISTS safeguarding_records (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  athlete_id TEXT NOT NULL REFERENCES athletes(id) ON DELETE CASCADE,
  guardian_details JSONB DEFAULT '{}'::JSONB,
  manager_details JSONB DEFAULT '{}'::JSONB,
  safeguarding_officer_assigned TEXT REFERENCES users(id) ON DELETE SET NULL,
  consent_forms_url TEXT DEFAULT '',
  anti_doping_education_completed BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Anti-Doping compliance table
CREATE TABLE IF NOT EXISTS anti_doping_compliance (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  athlete_id TEXT NOT NULL REFERENCES athletes(id) ON DELETE CASCADE,
  testing_status TEXT NOT NULL DEFAULT 'NOT_TESTED' CHECK (testing_status IN ('NOT_TESTED', 'IN_POOL', 'TESTED')),
  last_tested_on DATE,
  last_test_result TEXT DEFAULT 'NEGATIVE',
  wada_education_completed BOOLEAN DEFAULT FALSE,
  suspension_history TEXT DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## 2. Generic Backend CRUD Mapping

We will add the new tables/columns to the generic mapping definitions:
- `athletes` will map standard additional columns.
- Restricted endpoints will be exposed for `medical` and `safeguarding` checks.

```go
// In backend/internal/repository/nsmis_domains.go
"athletes": {
    "athletes",
    "EXISTS(SELECT 1 FROM athlete_affiliations aa WHERE aa.athlete_id=t.id AND aa.federation_id=ANY($3))",
    []string{"full_name", "gender", "date_of_birth", "discipline"},
    map[string]string{
        "athlete_number": "text",
        "national_id_passport": "text",
        "full_name": "text",
        "gender": "text",
        "date_of_birth": "date",
        "district": "text",
        "region": "text",
        "club": "text",
        "discipline": "text",
        "national_team_status": "text",
        "status": "text",
        "consent_basis": "text",
        "age_category": "text",
        "phone_contact": "text",
        "email_address": "text",
        "next_of_kin": "text",
        "emergency_contact": "text",
        "education_institution": "text",
        "highest_education_level": "text",
        "sports_scholarship_status": "boolean",
        "current_occupation": "text",
    },
},
```

## 3. UI Plan

- **Athlete Master Directory**: View list of registered athletes with filters for age groups, regions, and statuses.
- **Athlete Form**: Multi-tab form split into:
  - Personal details & Classification
  - Education & Dual career records
  - Medical history & Anti-Doping compliance status (restricted to authorized officers)
  - Safeguarding & Consent files
