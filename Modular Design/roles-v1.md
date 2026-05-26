# Roles Module — NCSMS v1.0

## Purpose

- Define and manage system roles and permissions across the entire NCSMS platform.
- Enable `SYSTEM_ADMIN` to create custom roles with scoped access to specific module features.
- Ensure role definitions are clear, auditable, and enforceable through RBAC, middleware, and RLS.

## Role Categories

- `SYSTEM_ADMIN` — full platform governance and role management authority.
- `GENERAL_SECRETARY` — high-level administrative role with broad operational control.
- `AUDITOR` — read-only access for audit and compliance review.
- `FEDERATION_USER` — scoped federation access for owner and officer workflows.
- `PUBLIC_USER` — public-facing limited access to read-only summaries.
- `ATHLETE` — athlete-specific profile and workflow access.
- `CONTENT_MANAGER` — website content management, blog publishing, pages, widgets, and analytics.
- `CUSTOM_ROLE` — any role created by `SYSTEM_ADMIN` with module- and action-specific permissions.

## Custom Role Creation

- `SYSTEM_ADMIN` can create, update, disable, and delete custom roles.
- Custom roles include:
  - `name`
  - `description`
  - `module_permissions`
  - `resource_scope`
  - `allowed_actions`
- `module_permissions` are defined per module and per entity type.
- Custom roles can be assigned to users through a dedicated `role_assignments` workflow.
- The system records who created or updated a custom role and when it was changed.

## Permission Model

Each role is defined as a set of module scopes and allowed actions:

- `VIEW` — read-only access to resources
- `CREATE` — ability to create new records
- `UPDATE` — ability to modify existing records
- `DELETE` — ability to remove records or cancel actions
- `EXECUTE` — ability to perform non-CRUD business actions

Permissions are evaluated at multiple levels:

- route/handler authorization
- service/business logic checks
- repository/RLS enforcement

## Module Permission Examples

### Super Admin Module

- `SYSTEM_ADMIN` has full access to create and manage users, roles, audit and system configuration.
- `AUDITOR` can only view audit logs and system state; no create, update, or delete operations.
- A custom role can be created that allows only `VIEW` access to user profiles and audit summaries.

### Admin Module

- `ADMIN` roles can approve federations, submit licenses, and disburse grants according to business rules.
- A custom finance role can be configured with `VIEW` access to federation budgets and `EXECUTE` access to grant disbursement while denying `DELETE` on licenses.
- A quality assurance role can be granted `VIEW` and `UPDATE` on documents but not `CREATE` new grant records.

### Federations Module

- `FEDERATION_USER` may manage only federations and accountabilities they own.
- A custom compliance role can access `VIEW` and `UPDATE` on license and document submissions within the same federation.
- A role with `CREATE` on `accountabilities` but no `DELETE` can submit reports without removing prior submissions.

### User Profiles Module

- `PUBLIC_USER` can only view profile summaries where allowed; no profile mutation.
- `AUDITOR` can read user metadata for compliance review; cannot update roles or sensitive fields.
- A custom user manager role can `CREATE` and `UPDATE` profile metadata but require `SYSTEM_ADMIN` approval to assign roles.

### Applications Module

- `GENERAL_SECRETARY` is the primary reviewer for user-submitted applications, renewals, and verification workflows.
- Applications include federation registration/renewal, academy operations, facility licenses, competition approvals, transformation forms, and national sport recognition.
- `GENERAL_SECRETARY` can `VIEW` applicant history, `UPDATE` application status, `EXECUTE` approval/verification actions, and audit related transactions.
- A custom application reviewer role may be configured with `VIEW` and `UPDATE` on applications but no access to end-user profile role assignments.

### Athletes Module

- `FEDERATION_USER` can `CREATE`, `VIEW`, and `UPDATE` athlete records for their federation.
- A custom analyst role can be granted read-only access to athlete performance data without document edit privileges.
- `SYSTEM_ADMIN` can manage all athlete records across federations.

### Assets Module

- `ADMIN` and federation asset managers can `CREATE`, `UPDATE`, and `VIEW` assets scoped to their federation.
- A custom inventory reviewer role may have `VIEW` and `UPDATE` on asset conditions but not `DELETE` or `TRANSFER`.
- Asset disposal should require a role with explicit `DELETE` permission and audit justification.

### Audit Module

- `AUDITOR` has read-only access to audit events.
- `SYSTEM_ADMIN` can view audit results and configure audit monitoring but cannot modify audit entries.
- Custom security reviewer roles can be given `VIEW` on specific event categories.

### Website Content Module

- `CONTENT_MANAGER` can manage blog posts, post categories, subcategories, pages, testimonials, footer widgets, contact page content, projects, and public page composition.
- Content managers may `CREATE`, `UPDATE`, and `VIEW` website components; `DELETE` is allowed only when explicitly granted.
- Typical content manager scope includes blog publication workflows, page metadata, SEO content, and blog analytics dashboards.
- A custom editor role can be configured with `VIEW`, `CREATE`, and `UPDATE` on website content but no access to permissions or system configuration.

## System-wide Role Enforcement

- Role permissions are stored in a dedicated `custom_roles` and `role_permissions` schema.
- `RBACMiddleware` loads effective permissions for the authenticated user and attaches them to request context.
- Service methods verify action-level permission before executing business logic.
- RLS policies use role and user category claims to restrict rows and columns based on role scope.
- Custom roles can be revoked or disabled without removing the assignment history.

## Role Lifecycle

- Create custom role
- Assign role to user
- Validate role permissions on every request
- Audit custom role changes and assignments
- Disable or archive stale/unused custom roles

## Example Role Permission Matrix

| Role | Super Admin | Admin | Federations | User Profiles | Athletes | Assets | Website | Applications | Audit |
|------|-------------|-------|-------------|---------------|----------|--------|---------|--------------|-------|
| SYSTEM_ADMIN | CRUDX | CRUDX | CRUDX | CRUDX | CRUDX | CRUDX | CRUDX | CRUDX | VIEW |
| GENERAL_SECRETARY | VIEW | CRUDX | VIEW | VIEW | VIEW | VIEW | VIEW | CRUDX | VIEW |
| AUDITOR | VIEW | VIEW | VIEW | VIEW | VIEW | VIEW | VIEW | VIEW | VIEW |
| FEDERATION_USER | NONE | VIEW | CRUDX | VIEW | CRUDX | VIEW | NONE | NONE | NONE |
| PUBLIC_USER | NONE | NONE | VIEW | VIEW | NONE | NONE | VIEW | NONE | NONE |
| ATHLETE | NONE | NONE | NONE | VIEW | UPDATE | NONE | NONE | NONE | NONE |
| CONTENT_MANAGER | NONE | NONE | NONE | VIEW | NONE | NONE | CRUDX | NONE | NONE |
| Athlete Manager | NONE | NONE | VIEW | VIEW | CRUDX | VIEW | NONE | NONE | NONE |
| Custom Role | Custom | Custom | Custom | Custom | Custom | Custom | Custom | Custom | Custom |

- `C` = Create, `R` = Read/View, `U` = Update, `D` = Delete, `X` = Execute

## Detailed Permission Matrices

### Module-level action matrix

| Role | User Profiles | Athletes | Assets | Federations | Admin | Website | Applications | Audit | Roles |
|------|---------------|----------|--------|-------------|-------|---------|--------------|-------|-------|
| SYSTEM_ADMIN | CRUDX | CRUDX | CRUDX | CRUDX | CRUDX | CRUDX | CRUDX | VIEW | CRUDX |
| GENERAL_SECRETARY | VIEW | VIEW | VIEW | VIEW | CRUDX | VIEW | CRUDX | VIEW | VIEW |
| AUDITOR | VIEW | VIEW | VIEW | VIEW | VIEW | VIEW | VIEW | NONE | NONE |
| FEDERATION_USER | VIEW | CRUDX | VIEW | CRUDX | VIEW | NONE | NONE | NONE | NONE |
| PUBLIC_USER | VIEW | NONE | NONE | VIEW | NONE | VIEW | NONE | NONE | NONE |
| ATHLETE | VIEW | UPDATE | NONE | NONE | NONE | NONE | NONE | NONE | NONE |
| CONTENT_MANAGER | NONE | NONE | NONE | NONE | NONE | CRUDX | NONE | NONE | NONE |
| CUSTOM_ROLE | Custom | Custom | Custom | Custom | Custom | Custom | Custom | Custom | Custom |

### Recommended entity permission matrix

| Role | User Profile | Email Verification | Federation | License | Grant | Athlete | Asset | Website Content | Application | Transaction Log | Audit Event | Custom Role |
|------|--------------|--------------------|------------|---------|-------|---------|-------|-----------------|-------------|----------------|-------------|-------------|
| SYSTEM_ADMIN | CRUDX | EXECUTE | CRUDX | CRUDX | CRUDX | CRUDX | CRUDX | CRUDX | CRUDX | CRUDX | VIEW | CRUDX |
| GENERAL_SECRETARY | VIEW | NONE | VIEW | VIEW | VIEW | VIEW | VIEW | VIEW | CRUDX | VIEW | VIEW | VIEW |
| AUDITOR | VIEW | NONE | VIEW | VIEW | VIEW | VIEW | VIEW | VIEW | VIEW | NONE | VIEW | NONE |
| FEDERATION_USER | VIEW | NONE | CRUDX | VIEW | VIEW | CRUDX | VIEW | NONE | NONE | NONE | NONE | NONE |
| PUBLIC_USER | VIEW | NONE | VIEW | NONE | NONE | NONE | NONE | VIEW | NONE | NONE | NONE | NONE |
| ATHLETE | VIEW | NONE | NONE | NONE | NONE | UPDATE | NONE | NONE | NONE | NONE | NONE | NONE |
| CONTENT_MANAGER | NONE | NONE | NONE | NONE | NONE | NONE | NONE | CRUDX | NONE | NONE | NONE | NONE |
| CUSTOM_ROLE | Custom | Custom | Custom | Custom | Custom | Custom | Custom | Custom | Custom | Custom | Custom | Custom |

### Example custom role behaviors

- `AUDITOR`: `VIEW` only on `Audit Event`, `Federation`, `User Profile`, `License`, and `Grant`; no `CREATE`, `UPDATE`, or `DELETE`.
- `GRANT_REVIEWER`: `VIEW` and `EXECUTE` on `Grant`, `VIEW` on `Federation` and `License`, no `DELETE` or `CREATE` on audit.
- `FEDERATION_COMPLIANCE`: `VIEW` and `UPDATE` on `License` and `Document` entities within the assigned federation; no access to global admin or role management.
- `ASSET_INSPECTOR`: `VIEW` on all `Asset` records and `UPDATE` on `Asset Condition`; no `DELETE`, `CREATE`, or `TRANSFER` privileges.

## Notes

- Make permission definitions explicit; avoid implicit inheritance from role names.
- When a custom role is created, require a review or approval workflow before enabling it in production.
- Document role scope and limitations clearly for each module and entity type.
