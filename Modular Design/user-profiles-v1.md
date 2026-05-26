# User Profiles Module — NCSMS v1.0

## Purpose

- Manage core user identity and category metadata
- Define system roles, user categories, and access entitlements
- Support authentication, user profile lifecycle, and category-specific permissions

## User Categories

- `SYSTEM_ADMIN` — full platform governance, RLS policy management, and audit oversight
- `GENERAL_SECRETARY` — executive administrative role with high-level management privileges, application review authority, and renewal processing responsibilities
- `FEDERATION_USER` — federation owner or officer with access scoped to federation-specific records
- `PUBLIC_USER` — read-only access to public federation summaries and general lookup data
- `AUDITOR` — read-only access for audit and compliance reporting
- `ATHLETE` — athlete profile category for athlete-specific workflows and registration data

## Domain Responsibilities

- Create and update user profiles with category and role metadata
- Assign primary role and category values to user profiles
- Manage user status, account activation, and role transitions
- Support verification flows such as email confirmation, 2FA, and account recovery
- Enforce category-specific metadata requirements and profile validation

## Website User Dashboard and Application Workflow

- Website signups receive a dashboard with an `Applications` section containing:
  - `Federation Registration`
  - `Federation Renewal`
  - `Sports Academy Operation Application`
  - `Community Sports Academy Registration / Renewal`
  - `Sports Facility Operation Application`
  - `Sports Competition Organization Application`
  - `National Sports Association Transformation`
  - `National Sports Federation Transformation`
  - `National Sport Recognition`
  - `Application History` and `Payment / Transaction Tracker`
- Each submission is issued a unique `application_reference` and `payment_reference` for audit and payment reconciliation.
- Users can view `submitted_at`, `current_status`, `assigned_reviewer`, `next_action`, and `last_updated` from their dashboard.
- Draft applications are saved for later completion; submitted applications follow a review workflow handled by `GENERAL_SECRETARY` and platform approvers.

## High-level Code Layout

- `/internal/domain/userprofiles/models.go`
  - `UserProfile`, `UserCategory`, `RoleAssignment`, `EmailVerificationRequest`, `TwoFactorProfile`
- `/internal/domain/userprofiles/repository.go`
  - `CreateUserProfile(ctx, profile)`, `UpdateUserProfile(ctx, profile)`, `FindUserByID(ctx, id)`, `FindUserByEmail(ctx, email)`
- `/internal/domain/userprofiles/service.go`
  - `RegisterUser(ctx, req)`, `UpdateUserProfile(ctx, req)`, `AssignRole(ctx, userID, role)`, `VerifyEmail(ctx, token)`
- `/internal/domain/userprofiles/handler.go`
  - HTTP endpoints: `POST /api/v1/users`, `GET /api/v1/users/{id}`, `PATCH /api/v1/users/{id}`, `POST /api/v1/users/{id}/verify-email`

## Access and Category Rules

- Map `user_metadata.primary_role` to category-specific access policies
- Enforce `PUBLIC_USER` restrictions at both application and RLS level
- Treat `FEDERATION_USER` and `AUDITOR` as separate categories with distinct allowed scopes
- Use explicit category checks in service logic, not just role names
- Validate that user category transitions are authorized by a privileged role

## Example Routes

- `POST /api/v1/users` -> register a user profile
- `GET /api/v1/users/{id}` -> fetch user profile details
- `PATCH /api/v1/users/{id}` -> update profile fields and category metadata
- `POST /api/v1/users/{id}/verify-email` -> complete email verification flow
- `GET /api/v1/applications` -> list user application submissions and statuses
- `GET /api/v1/applications/{id}` -> fetch application details and tracking history
- `POST /api/v1/applications/federation-registration` -> submit a federation registration
- `POST /api/v1/applications/federation-renewal` -> submit a federation renewal
- `POST /api/v1/applications/academy-operation` -> submit a sports academy operation application
- `POST /api/v1/applications/community-academy-registration-renewal` -> submit a community academy registration or renewal
- `POST /api/v1/applications/facility-operation` -> submit a sports facility operation application
- `POST /api/v1/applications/competition-organization` -> submit a sports competition organization application
- `POST /api/v1/applications/transformation/association` -> submit a national sports association transformation application
- `POST /api/v1/applications/transformation/federation` -> submit a national sports federation transformation application
- `POST /api/v1/applications/national-sport-recognition` -> submit a national sport recognition application
- `GET /api/v1/transactions` -> list transaction logs for user payments and fees
- `GET /api/v1/transactions/{reference}` -> fetch payment and transaction details by reference
