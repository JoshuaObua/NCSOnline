# NAMIS Clubs Module Implementation Plan

The Clubs Module handles the registrations, affiliations, and details of clubs, academies, and institutions under National Federations.

## 1. Database Schema Plan

A new table `clubs` will be created in the database schema:

```sql
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
```

## 2. Generic Backend CRUD Mapping

We will add `clubs` to the generic `nsmis_domains.go` mappings inside the Go backend:

```go
// In backend/internal/repository/nsmis_domains.go
"clubs": {
    "clubs", 
    "t.federation_id=ANY($3)", 
    []string{"federation_id", "name", "acronym"}, 
    map[string]string{
        "federation_id": "text", 
        "name": "text", 
        "acronym": "text", 
        "contact_person": "text", 
        "email": "text", 
        "phone": "text", 
        "district": "text", 
        "region": "text", 
        "date_founded": "date?", 
        "status": "text",
    },
},
```

And update routes definition in `backend/cmd/server/main.go` to include `clubs` in the parameter routing regex:
```go
r.Route("/{resource:...|equipment|clubs}", func(r chi.Router) { ... })
```

## 3. UI Plan

- **Manage Clubs Panel**: A grid displaying recognized clubs by federation with search filters by district, region, and status.
- **Club Form**: A card interface for adding/editing clubs allowing selection of parent Federation and input of contacts, founding date, and regions.
- **Roster view**: Displays list of athletes affiliated with the club.
