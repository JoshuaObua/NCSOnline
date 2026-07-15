# NAMIS Performance & Medals Module Implementation Plan

The Performance and Medals Module tracks details of competitions, detailed athlete results (times, weights, distances, scores), records status (NR/PB/SB), and medal details (Gold, Silver, Bronze) won by athletes.

## 1. Database Schema Plan

We leverage and extend the existing `competitions`, `competition_results`, and `medals` tables.

```sql
-- Extend competitions level Check Constraint to support school and regional levels
ALTER TABLE competitions DROP CONSTRAINT IF EXISTS competitions_level_check;
ALTER TABLE competitions ADD CONSTRAINT competitions_level_check CHECK (level IN ('DISTRICT', 'REGIONAL', 'NATIONAL', 'EAST_AFRICAN', 'AFRICAN', 'COMMONWEALTH', 'OLYMPIC', 'WORLD_CHAMPIONSHIP'));

-- Extend competition_results table to support granular result fields
ALTER TABLE competition_results ADD COLUMN IF NOT EXISTS time_result TEXT DEFAULT '';
ALTER TABLE competition_results ADD COLUMN IF NOT EXISTS distance_result TEXT DEFAULT '';
ALTER TABLE competition_results ADD COLUMN IF NOT EXISTS weight_result TEXT DEFAULT '';
ALTER TABLE competition_results ADD COLUMN IF NOT EXISTS score_result TEXT DEFAULT '';
ALTER TABLE competition_results ADD COLUMN IF NOT EXISTS ranking_result TEXT DEFAULT '';
ALTER TABLE competition_results ADD COLUMN IF NOT EXISTS is_national_record BOOLEAN DEFAULT FALSE;
ALTER TABLE competition_results ADD COLUMN IF NOT EXISTS is_personal_best BOOLEAN DEFAULT FALSE;
ALTER TABLE competition_results ADD COLUMN IF NOT EXISTS is_seasonal_best BOOLEAN DEFAULT FALSE;

-- Extend medals table
ALTER TABLE medals ADD COLUMN IF NOT EXISTS level TEXT DEFAULT 'NATIONAL' CHECK (level IN ('DISTRICT', 'REGIONAL', 'NATIONAL', 'EAST_AFRICAN', 'AFRICAN', 'COMMONWEALTH', 'OLYMPIC', 'WORLD_CHAMPIONSHIP'));
ALTER TABLE medals ADD COLUMN IF NOT EXISTS coach_at_win_id TEXT REFERENCES coaches(id) ON DELETE SET NULL;
ALTER TABLE medals ADD COLUMN IF NOT EXISTS is_team_event BOOLEAN DEFAULT FALSE;
ALTER TABLE medals ADD COLUMN IF NOT EXISTS prize_money NUMERIC(18,2) DEFAULT 0.00;
ALTER TABLE medals ADD COLUMN IF NOT EXISTS ncs_recognition_status TEXT DEFAULT 'PENDING' CHECK (ncs_recognition_status IN ('PENDING', 'APPROVED', 'REJECTED'));
```

## 2. Generic Backend CRUD Mapping

We update `competitions`, `medals` and expose `competition-results` definitions in `nsmis_domains.go`:

```go
// In backend/internal/repository/nsmis_domains.go
"competitions": {
    "competitions", "t.federation_id=ANY($3)",
    []string{"federation_id", "name", "venue", "level", "starts_on", "ends_on"},
    map[string]string{"federation_id": "text", "name": "text", "venue": "text", "host_country": "text", "level": "text", "starts_on": "date", "ends_on": "date", "returned_on": "date?", "athlete_count": "integer", "official_count": "integer", "status": "text"},
},
"medals": {
    "medals", "t.federation_id=ANY($3)",
    []string{"competition_id", "federation_id", "athlete_name", "event", "country", "won_on", "medal_type"},
    map[string]string{"competition_id": "text", "federation_id": "text", "athlete_id": "text?", "athlete_name": "text", "event": "text", "country": "text", "won_on": "date", "medal_type": "text", "coach_responsible": "text", "team_manager": "text", "funding_source": "text", "status": "text", "level": "text", "coach_at_win_id": "text?", "is_team_event": "boolean", "prize_money": "numeric", "ncs_recognition_status": "text"},
},
"competition-results": {
    "competition_results", 
    "EXISTS(SELECT 1 FROM competitions c WHERE c.id=t.competition_id AND c.federation_id=ANY($3))",
    []string{"competition_id", "athlete_name", "event"},
    map[string]string{
        "competition_id": "text", "athlete_id": "text?", "athlete_name": "text", "event": "text", "position": "integer?",
        "time_result": "text", "distance_result": "text", "weight_result": "text", "score_result": "text", "ranking_result": "text",
        "is_national_record": "boolean", "is_personal_best": "boolean", "is_seasonal_best": "boolean",
    },
},
```

## 3. UI Plan

- **Competitions Dashboard**: List of hosted and international competitions with filters.
- **Results Form**: Add and view athlete finishes in a tabular sheet mapping times/distances/scores and record tags.
- **Medal Board**: Track individual or team medals, funding sources, and NCS recognition approvals.
