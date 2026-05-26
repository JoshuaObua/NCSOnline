# Admin Module — NCSMS v1.0

## Purpose

- Day-to-day platform management by authorized administrators
- Manage federation lifecycles, licensing, grants, and compliance documentation
- Enforce business rules for accounts and financial operations

## Domain Responsibilities

- Federation validations and approvals
- Licensing workflow orchestration
- Grant approvals and disbursement gating
- Inventory tracking oversight and reporting

## High-level Code Layout

- `/internal/domain/admin/models.go`
  - `Federation`, `License`, `Grant`, `GrantAccountability`, `InventoryAsset`
- `/internal/domain/admin/repository.go`
  - SQL operations for read/write on federation, license, grant, asset domains
- `/internal/domain/admin/service.go`
  - `ApproveFederation`, `SubmitLicense`, `CreateGrant`, `LockAccountability`
- `/internal/domain/admin/handler.go`
  - HTTP endpoints: `GET /api/v1/admin/federations`, `POST /api/v1/admin/licenses`, `POST /api/v1/admin/grants/disburse`

## Business Guardrails

- Zero-Overdraft Enforcement: disbursement cannot exceed allocation
- Locked state enforcement: SUBMITTED / APPROVED accountabilities become immutable
- Disbursement block upon 90+ day unapproved accountabilities
- License state machine: `DRAFT -> SUBMITTED -> UNDER_REVIEW -> APPROVED|REJECTED`
- Duplicate license window detection by federation

## Service Patterns

- Use `pgxpool.Pool.BeginTx` with `Serializable` isolation for all financial operations
- Query and lock rows via `FOR UPDATE` before writes
- Reject business state transitions if preconditions fail
- Persist audit events within same transaction

## Example Routes

- `GET /api/v1/admin/federations` -> admin federation list
- `PATCH /api/v1/admin/federations/{id}` -> update federation status
- `POST /api/v1/admin/grants` -> allocate grant funds
- `POST /api/v1/admin/grants/{id}/disburse` -> disburse funds after validation

## Detailed Module Design

### Models

- `FederationRecord`
- `LicenseRecord`
- `GrantDisbursementRequest`
- `GrantAllocation`

### Repository

- `CreateFederation(ctx, fed)`
- `ActivateFederation(ctx, fedID)`
- `CreateGrant(ctx, grant)`
- `LockGrantDisbursement(ctx, grantID)`
- `HasAgedUnapprovedAccountability(ctx, fedID)`

### Service

- `RegisterFederation(ctx, req)` ensures no duplicate `reg_number` or `code_name`
- `SubmitLicenseApplication(ctx, req)` enforces workflow states and document pass checks
- `DisburseGrant(ctx, req)` checks `allocated_amount`, `disbursed_amount`, and 90-day accountability constraints
- `LockAccountability(ctx, entryID)` prevents edits once SUBMITTED/APPROVED

### Handler Routes

- `GET /api/v1/admin/federations`
- `POST /api/v1/admin/federations`
- `PATCH /api/v1/admin/federations/{id}/status`
- `POST /api/v1/admin/grants/{id}/disburse`

