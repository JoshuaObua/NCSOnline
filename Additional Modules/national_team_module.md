# NAMIS National Team Representation Plan

The National Team Module monitors international participation, team selections, call-up records, and cap statistics for senior and development squads.

## 1. Database Schema Plan

A new table `national_team_appearances` will be created to track appearances and callups:

```sql
CREATE TABLE IF NOT EXISTS national_team_appearances (
  id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
  athlete_id TEXT NOT NULL REFERENCES athletes(id) ON DELETE CASCADE,
  team_name TEXT NOT NULL,
  category TEXT NOT NULL CHECK (category IN ('SENIOR', 'DEVELOPMENT', 'JUNIOR')),
  first_call_up_on DATE,
  last_appearance_on DATE,
  appearances_count INTEGER NOT NULL DEFAULT 1 CHECK (appearances_count >= 1),
  notes TEXT DEFAULT '',
  created_by TEXT REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_athlete_national_team ON national_team_appearances(athlete_id, team_name, category);
```

## 2. Generic Backend CRUD Mapping

We map `national-team` in `nsmis_domains.go` to handle operations:

```go
// In backend/internal/repository/nsmis_domains.go
"national-team": {
    "national_team_appearances", 
    "EXISTS(SELECT 1 FROM athletes a JOIN athlete_affiliations aa ON aa.athlete_id=a.id WHERE a.id=t.athlete_id AND aa.federation_id=ANY($3))", 
    []string{"athlete_id", "team_name", "category"}, 
    map[string]string{
        "athlete_id": "text", 
        "team_name": "text", 
        "category": "text", 
        "first_call_up_on": "date?", 
        "last_appearance_on": "date?", 
        "appearances_count": "integer",
        "notes": "text",
    },
},
```

## 3. UI Plan

- **National Squads Tracker**: Manage selections, listing athletes selected for the national team by sport.
- **Appearances Form**: Add or edit appearance statistics, tracking dates of first callup and total matches/competitions played.
