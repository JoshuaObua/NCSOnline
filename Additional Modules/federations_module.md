# NAMIS Federations Module Implementation Plan

The Federations Module handles recognized National Federations and Sports Associations, their registration records, governing officers, and aggregated statistics.

## 1. Database Schema Plan

We leverage the existing `federations`, `federation_memberships`, and `federation_officers` tables:

```sql
-- Existing federations table
-- id, name, acronym, ncs_registration_number, recognition_status, physical_address, email, website, contact_person, etc.

-- Existing federation_officers table
-- id, federation_id, position, full_name, email, phone, appointed_on, term_ends_on, is_active

-- Existing federation_memberships table
-- id, federation_id, user_id, membership_role, starts_at, ends_at, is_active
```

We will extend `federation_officers` to track roles like President, General Secretary, Treasurer, and Arbitrator explicitly by establishing standard lookup tables or constraints on `position` if required.

## 2. Generic Backend CRUD Mapping

The generic backend already maps:
- `federations` (via standard routes)
- `federation-officers` (via `nsmis_domains.go`)
- `federation-memberships` (via `nsmis_domains.go`)

We will ensure statistics endpoints or queries dynamically aggregate:
- Number of registered athletes (`COUNT` on `athlete_affiliations` where `federation_id = ID`)
- Number of registered coaches (`COUNT` on `coaches` where `federation_id = ID`)
- Number of clubs (`COUNT` on `clubs` where `federation_id = ID`)

## 3. UI Plan

- **Federation Profile Manager**: Edit federation details, track address, contact details, and current registration status.
- **Officers Panel**: Manage key executive roles (President, General Secretary, Treasurer, Arbitrator) with term dates.
- **Federation Statistics**: High-level counters displaying registered athletes, active coaches, and affiliated clubs.
