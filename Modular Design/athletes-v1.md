# Athletes Module — NCSMS v1.0

## Purpose

- Manage athlete registration, verification, and lifecycle data
- Track athlete credentials, eligibility, and federation affiliation
- Support athlete-specific reporting and document management

## Domain Responsibilities

- Create and maintain athlete master records
- Associate athletes with federations and competition categories
- Manage athlete credential documents, medical clearances, and eligibility flags
- Support athlete search, reporting, and audit history

## High-level Code Layout

- `/internal/domain/athletes/models.go`
  - `Athlete`, `AthleteRegistration`, `AthleteDocument`, `AthleteStatus`
- `/internal/domain/athletes/repository.go`
  - `CreateAthlete(ctx, athlete)`, `FindAthleteByID(ctx, id)`, `ListAthletes(ctx, filter)`, `UpdateAthlete(ctx, athlete)`
- `/internal/domain/athletes/service.go`
  - `RegisterAthlete(ctx, req)`, `UpdateAthleteProfile(ctx, req)`, `AttachAthleteDocument(ctx, req)`, `VerifyAthleteEligibility(ctx, req)`
- `/internal/domain/athletes/handler.go`
  - HTTP endpoints: `POST /api/v1/athletes`, `GET /api/v1/athletes/{id}`, `PATCH /api/v1/athletes/{id}`, `GET /api/v1/athletes`

## Access and Policy

- `FEDERATION_USER` and `ADMIN` roles may manage athletes within their federation scope
- `SYSTEM_ADMIN` may manage all athlete records for platform governance
- `PUBLIC_USER` can only access non-sensitive athlete summary data if allowed
- Athlete sensitive documents are stored in private storage with signed URLs

## Example Routes

- `POST /api/v1/athletes` -> register a new athlete
- `GET /api/v1/athletes/{id}` -> retrieve an athlete profile
- `PATCH /api/v1/athletes/{id}` -> update athlete profile or status
- `GET /api/v1/athletes` -> search and filter athlete records
