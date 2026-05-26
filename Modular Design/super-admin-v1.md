# Super Admin Module — NCSMS v1.0

## Purpose

- Global system control
- Manage users, roles, RLS policies, audit monitoring, platform-level configuration
- Oversee all domains and enforce compliance

## Domain Responsibilities

- Create and manage `SYSTEM_ADMIN`, `GENERAL_SECRETARY`, and other privileged profiles
- Define and manage custom super admin roles with explicit access scopes and resource permissions
- Manage federation onboarding approvals and platform state
- Monitor audit event streams and immutable logs with full tracking metadata
- Manage users across the platform, including direct password reset and account recovery actions
- Execute emergency governance actions

## Custom Role Management

- `SYSTEM_ADMIN` can create, update, disable, and delete custom role definitions.
- Custom roles are defined per module and per resource action: `VIEW`, `CREATE`, `UPDATE`, `DELETE`, `EXECUTE`.
- Roles may be scoped to whole modules or narrow resource groups such as `licenses`, `grants`, `athletes`, `assets`, and `audit`.
- Example custom role: `AUDITOR` with only `VIEW` permissions for audit logs, federation summaries, and user profile metadata.
- Example custom role: `GRANT_REVIEWER` with `VIEW` and `EXECUTE` on grant disbursement workflows but no `DELETE` privileges.
- Custom roles are enforced in middleware, service validation, and RLS policies.
- Role assignments are audited with creator, updater, and assignment history metadata.

## High-level Code Layout

- `/internal/domain/superadmin/models.go`
  - `UserProfile`, `SystemRole`, `CustomRole`, `AuditQueryFilter`, `PasswordResetRequest`, `GlobalConfig`
- `/internal/domain/superadmin/repository.go`
  - `GetUserProfileByID(ctx, id)`, `ListAllProfiles(ctx, pagination)`, `RevokeAccount(ctx, id)`, `CreateCustomRole(ctx, role)`, `ResetUserPassword(ctx, id, passwordHash)`
- `/internal/domain/superadmin/service.go`
  - `CreateSuperAdmin(ctx, req)`, `CreateCustomRoleWithAccess(ctx, req)`, `PromoteAdmin(ctx, target, role)`, `ResetPassword(ctx, userID, req)`, `AuditSearch(ctx, filter)`
- `/internal/domain/superadmin/handler.go`
  - HTTP endpoints: `POST /api/v1/superadmin/users`, `POST /api/v1/superadmin/roles`, `POST /api/v1/superadmin/users/{id}/password-reset`, `GET /api/v1/superadmin/audit`, `PATCH /api/v1/superadmin/users/{id}/role`

## Transaction and Validation Rules

- All service operations execute within `pgx.TxOptions{IsoLevel: pgx.Serializable}`
- Validate role state transitions and deny invalid promotions
- Enforce custom role access definitions before granting permissions to users or endpoints
- Encrypt sensitive payloads such as reset tokens and temporary secrets before persistence
- Ensure every mutation writes an `audit_logs` event atomically with tracing metadata

## Endpoint Patterns

- `POST /api/v1/superadmin/users` -> create system admin user profile
- `POST /api/v1/superadmin/roles` -> create a custom role with defined access permissions
- `POST /api/v1/superadmin/users/{id}/password-reset` -> trigger direct password reset for user
- `PATCH /api/v1/superadmin/users/{id}/status` -> enable/disable account
- `GET /api/v1/superadmin/audit` -> query immutable logs with full tracking details

## Detailed Module Design

### Models

- `SuperAdminProfile`
- `SuperAdminCreateRequest`
- `AuditQueryFilter`

### Repository

- `FetchProfileByID(ctx, uuid)`
- `ListProfiles(ctx, pagination)`
- `SetProfileStatus(ctx, uuid, active)`
- `InsertAuditEvent(ctx, audit)`

### Service

- `CreateProfile(ctx, req)` validates required fields and the `SYSTEM_ADMIN` role
- `CreateCustomRole(ctx, req)` enforces explicit access and resource permissions for new role definitions
- `ChangeRole(ctx, userID, role)` ensures valid promotion sequence
- `ResetPassword(ctx, userID, req)` handles direct password reset flows with secure token and audit trail
- `FetchAuditEntries(ctx, filter)` supports date range and action filtering with full tracking fields

### Handler Routes

- `POST /api/v1/superadmin/users`
- `POST /api/v1/superadmin/roles`
- `POST /api/v1/superadmin/users/{id}/password-reset`
- `PATCH /api/v1/superadmin/users/{id}/role`
- `PATCH /api/v1/superadmin/users/{id}/status`
- `GET /api/v1/superadmin/audit`

### Audit Tracking Details

Super Admin audit logs surface:

- source IP address
- geolocation or location metadata
- device type and platform
- requested endpoint
- response status code
- request timestamp and processing duration
- user agent and session identifier
- policy decision or authorization context

