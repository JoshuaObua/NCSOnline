# NAMIS Talent Identification Module Implementation Plan

The Talent Identification Module manages the search, discovery, and recommended development pathways for young sporting talent across schools and regional centres in Uganda.

## 1. Database Schema Plan

We leverage and extend the existing `talent_records` and `talent_scholarships` tables.

```sql
-- Extend talent_records table
ALTER TABLE talent_records ADD COLUMN IF NOT EXISTS talent_centre TEXT DEFAULT '';
ALTER TABLE talent_records ADD COLUMN IF NOT EXISTS talent_category TEXT DEFAULT 'AMATEUR' CHECK (talent_category IN ('AMATEUR', 'EMERGING', 'ELITE'));
ALTER TABLE talent_records ADD COLUMN IF NOT EXISTS recommended_pathway TEXT DEFAULT '';
ALTER TABLE talent_records ADD COLUMN IF NOT EXISTS scholarship_status TEXT DEFAULT 'NONE' CHECK (scholarship_status IN ('NONE', 'APPLIED', 'ACTIVE', 'EXPIRED'));
```

## 2. Generic Backend CRUD Mapping

We map `talent` (via `talent_records`) and `scholarships` (via `talent_scholarships`) in `nsmis_domains.go`:

```go
// In backend/internal/repository/nsmis_domains.go
"talent": {
    "talent_records", "t.federation_id=ANY($3)", 
    []string{"federation_id", "athlete_name", "age_at_identification", "district", "identified_by", "identified_on"}, 
    map[string]string{
        "federation_id": "text", "athlete_id": "text?", "athlete_name": "text", "age_at_identification": "integer", 
        "school": "text", "district": "text", "region": "text", "identified_by": "text", "identified_on": "date", "status": "text",
        "talent_centre": "text", "talent_category": "text", "recommended_pathway": "text", "scholarship_status": "text",
    },
},
"scholarships": {
    "talent_scholarships", 
    "EXISTS(SELECT 1 FROM talent_records tr WHERE tr.id=t.talent_record_id AND tr.federation_id=ANY($3))", 
    []string{"talent_record_id", "provider", "starts_on"}, 
    map[string]string{
        "talent_record_id": "text", "provider": "text", "starts_on": "date", "ends_on": "date?", "status": "text",
    },
},
```

## 3. UI Plan

- **Talent Identification Registry**: Interface mapping scouted prospects, their school, district of origin, and pathway recommendation.
- **Scout / Academy Form**: Profile builder for submitting raw talent details.
- **Pathway Tracker**: Manage progress stages (Identified -> Selected -> National Team / Academy).
- **Scholarship Panel**: Manage financial support, sponsors/providers, and scholarship duration.
