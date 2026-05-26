# Federations Module — NCSMS v1.0

## Purpose

- Federation self-service operations
- Federation user onboarding, document submission, and limited lifecycle management
- Federation-specific record access under RLS and admin oversight

## Domain Responsibilities

- Federation profile self-management
- License and document submission workflows
- Grant accountability submission and report generation
- Federation asset requests and status visibility

## High-level Code Layout

- `/internal/domain/federation/models.go`
  - `FederationProfile`, `LicenseApplication`, `DocumentUpload`, `FederationAuditRecord`
- `/internal/domain/federation/repository.go`
  - `ListFederationsByOwner`, `FindFederationByID`, `CreateLicenseApplication`
- `/internal/domain/federation/service.go`
  - `SubmitFederationLicense`, `UploadComplianceDocument`, `SubmitAccountabilityReport`
- `/internal/domain/federation/handler.go`
  - Endpoints: `GET /api/v1/federations`, `POST /api/v1/federations/{id}/license`, `POST /api/v1/federations/{id}/accountabilities`

## Access and Policy

- `FEDERATION_USER` sees only federations where `primary_contact_id = auth.uid()`
- `PUBLIC_USER` can query active federation summaries only
- Federation document uploads validate checksum and metadata before storage
- Signed URL generation uses storage policies and short expiration

## Workflow Rules

- Federation license submissions must include required compliance metadata
- Any license approval path is blocked until all mandatory documents pass checksum validation
- Federation users cannot mutate locked grant accountability records

## Detailed Module Design

### Models

- `FederationSelfUpdateRequest`
- `LicenseSubmissionRequest`
- `AccountabilitySubmissionRequest`
- `DocumentUploadMeta`

### Repository

- `ListFederationsByUser(ctx, userID)`
- `FindFederationWithOwner(ctx, federationID, ownerID)`
- `InsertLicenseSubmission(ctx, license)`
- `SubmitAccountability(ctx, accountability)`

### Service

- `ListOwnedFederations(ctx, userID)` returns only federations bound by primary contact
- `CreateLicenseSubmission(ctx, req)` verifies mandatory document checksums and state sequence
- `SubmitAccountabilityReport(ctx, req)` persists accountability with `SUBMITTED` and prevents updates in locked states

### Handler Routes

- `GET /api/v1/federations`
- `GET /api/v1/federations/{id}`
- `POST /api/v1/federations/{id}/licenses`
- `POST /api/v1/federations/{id}/accountabilities`

