# NAMIS Coaches Module Implementation Plan

The Coaches Module manages athletic instructors, licensing status, coaching levels, and auxiliary training support staff (S&C, physiotherapist, team doctor, nutritionist, etc.).

## 1. Database Schema Plan

We leverage and extend the existing `coaches` table and introduce support staff mappings:

```sql
-- Extend coaches table
ALTER TABLE coaches ADD COLUMN IF NOT EXISTS phone TEXT DEFAULT '';
ALTER TABLE coaches ADD COLUMN IF NOT EXISTS email TEXT DEFAULT '';

-- Support Staff Mapping table to map training entourage for athletes/teams
CREATE TABLE IF NOT EXISTS athlete_support_entourage (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  athlete_id TEXT NOT NULL REFERENCES athletes(id) ON DELETE CASCADE,
  primary_coach_id TEXT REFERENCES coaches(id) ON DELETE SET NULL,
  assistant_coach_name TEXT DEFAULT '',
  strength_conditioning_coach TEXT DEFAULT '',
  sports_scientist TEXT DEFAULT '',
  physiotherapist TEXT DEFAULT '',
  team_doctor TEXT DEFAULT '',
  sports_nutritionist TEXT DEFAULT '',
  team_manager TEXT DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(athlete_id)
);
```

## 2. Generic Backend CRUD Mapping

We map `coaches` in `nsmis_domains.go`:

```go
// In backend/internal/repository/nsmis_domains.go
"coaches": {
    "coaches", 
    "t.federation_id=ANY($3)", 
    []string{"federation_id", "full_name", "certification_level", "license_number"}, 
    map[string]string{
        "federation_id": "text", 
        "full_name": "text", 
        "certification_level": "text", 
        "license_number": "text", 
        "expiry_date": "date?", 
        "status": "text",
        "phone": "text",
        "email": "text",
    },
},
```

We will also expose a generic mapping for `athlete-support-entourage` to manage training teams:
```go
"athlete-support-entourage": {
    "athlete_support_entourage",
    "EXISTS(SELECT 1 FROM athletes a JOIN athlete_affiliations aa ON aa.athlete_id=a.id WHERE a.id=t.athlete_id AND aa.federation_id=ANY($3))",
    []string{"athlete_id"},
    map[string]string{
        "athlete_id": "text",
        "primary_coach_id": "text?",
        "assistant_coach_name": "text",
        "strength_conditioning_coach": "text",
        "sports_scientist": "text",
        "physiotherapist": "text",
        "team_doctor": "text",
        "sports_nutritionist": "text",
        "team_manager": "text",
    },
},
```

## 3. UI Plan

- **Coaches Registry**: Grid displaying licensed coaches, certifications, license numbers, and expiry states.
- **Coach Details Form**: Add or edit coach details including federation and contact details.
- **Support Entourage Panel**: An interface under an athlete profile where officers can assign coaches, physiotherapists, S&C specialists, nutritionists, and doctors.
