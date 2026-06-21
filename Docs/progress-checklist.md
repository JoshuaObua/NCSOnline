# NCSMS v1 — API Progress Checklist

> Status key: ✅ Complete | 🔄 In Progress | ⬜ Not Started | 🚫 Blocked

---

## Authentication (`/api/v1/auth`)

| # | Method | Endpoint | Description | Status |
|---|---|---|---|---|
| 1 | POST | `/api/v1/auth/login` | Email + password login; returns access + refresh token | 🔄 In Progress |
| 2 | POST | `/api/v1/auth/refresh` | Issue new access token from refresh token | 🔄 In Progress |
| 3 | POST | `/api/v1/auth/logout` | Revoke current refresh token | 🔄 In Progress |
| 4 | GET | `/api/v1/auth/me` | Return signed-in user profile + roles | 🔄 In Progress |
| 5 | PUT | `/api/v1/auth/me/password` | Change own password | 🔄 In Progress |
| 6 | POST | `/api/v1/auth/forgot-password` | Send password reset email | 🔄 In Progress |
| 7 | POST | `/api/v1/auth/reset-password` | Reset password via token | 🔄 In Progress |

---

## Admin — Dashboard (`/api/v1/admin/dashboard`)

| # | Method | Endpoint | Description | Status | Roles |
|---|---|---|---|---|---|
| 8 | GET | `/api/v1/admin/dashboard` | Aggregate stats: users, applications by status, payment totals | 🔄 In Progress | super_admin, admin |

---

## Admin — User Management (`/api/v1/admin/users`)

| # | Method | Endpoint | Description | Status | Roles |
|---|---|---|---|---|---|
| 9 | GET | `/api/v1/admin/users` | Paginated list of all users with filters | 🔄 In Progress | super_admin, admin |
| 10 | POST | `/api/v1/admin/users` | Create a new user account | 🔄 In Progress | super_admin, admin |
| 11 | GET | `/api/v1/admin/users/{id}` | Get single user with roles | 🔄 In Progress | super_admin, admin |
| 12 | PUT | `/api/v1/admin/users/{id}` | Update user profile fields | 🔄 In Progress | super_admin, admin |
| 13 | DELETE | `/api/v1/admin/users/{id}` | Soft-delete a user | 🔄 In Progress | super_admin |
| 14 | POST | `/api/v1/admin/users/{id}/activate` | Activate a deactivated user | 🔄 In Progress | super_admin, admin |
| 15 | POST | `/api/v1/admin/users/{id}/deactivate` | Deactivate a user | 🔄 In Progress | super_admin, admin |
| 16 | POST | `/api/v1/admin/users/{id}/roles` | Assign a role to a user | 🔄 In Progress | super_admin |
| 17 | DELETE | `/api/v1/admin/users/{id}/roles/{roleID}` | Remove a role from a user | 🔄 In Progress | super_admin |

---

## Admin — Roles & Permissions (`/api/v1/admin/roles`)

| # | Method | Endpoint | Description | Status | Roles |
|---|---|---|---|---|---|
| 18 | GET | `/api/v1/admin/roles` | List all roles with permissions | 🔄 In Progress | super_admin |
| 19 | POST | `/api/v1/admin/roles` | Create a custom role | 🔄 In Progress | super_admin |
| 20 | GET | `/api/v1/admin/roles/{id}` | Get single role detail | 🔄 In Progress | super_admin |
| 21 | PUT | `/api/v1/admin/roles/{id}` | Update role name or description | 🔄 In Progress | super_admin |
| 22 | DELETE | `/api/v1/admin/roles/{id}` | Delete a non-system role | 🔄 In Progress | super_admin |
| 23 | GET | `/api/v1/admin/permissions` | List all available permissions | 🔄 In Progress | super_admin |
| 24 | POST | `/api/v1/admin/roles/{id}/permissions` | Assign a permission to a role | 🔄 In Progress | super_admin |
| 25 | DELETE | `/api/v1/admin/roles/{id}/permissions/{permID}` | Remove a permission from a role | 🔄 In Progress | super_admin |

---

## Admin — Application Review (`/api/v1/admin/applications`)

| # | Method | Endpoint | Description | Status | Roles |
|---|---|---|---|---|---|
| 26 | GET | `/api/v1/admin/applications` | Paginated list of all applications with filters | 🔄 In Progress | super_admin, admin, general_secretary |
| 27 | GET | `/api/v1/admin/applications/{id}` | Full application detail with attachments and payment | 🔄 In Progress | super_admin, admin, general_secretary |
| 28 | POST | `/api/v1/admin/applications/{id}/approve` | Approve application with optional note | 🔄 In Progress | super_admin, admin, general_secretary |
| 29 | POST | `/api/v1/admin/applications/{id}/reject` | Reject application with mandatory reason | 🔄 In Progress | super_admin, admin, general_secretary |
| 30 | POST | `/api/v1/admin/applications/{id}/request-info` | Set status to NEEDS_INFORMATION with note | 🔄 In Progress | super_admin, admin, general_secretary |
| 31 | POST | `/api/v1/admin/applications/{id}/verify-payment` | Confirm uploaded payment proof | 🔄 In Progress | super_admin, admin, general_secretary |
| 32 | POST | `/api/v1/admin/applications/{id}/reject-payment` | Reject uploaded payment proof with reason | 🔄 In Progress | super_admin, admin, general_secretary |

---

## Admin — Audit Logs (`/api/v1/admin/audit-logs`)

| # | Method | Endpoint | Description | Status | Roles |
|---|---|---|---|---|---|
| 33 | GET | `/api/v1/admin/audit-logs` | Paginated audit log with filters (user, action, date range) | 🔄 In Progress | super_admin, admin |
| 34 | GET | `/api/v1/admin/audit-logs/{id}` | Get single audit log entry | 🔄 In Progress | super_admin, admin |

---

## Applications — Draft & Wizard (`/api/v1/applications`)

| # | Method | Endpoint | Description | Status |
|---|---|---|---|---|
| 35 | POST | `/api/v1/applications/{formType}/draft` | Create or update a draft application | 🔄 In Progress |
| 36 | GET | `/api/v1/applications/{formType}/draft` | Load latest draft for the signed-in user | 🔄 In Progress |
| 37 | DELETE | `/api/v1/applications/{formType}/draft` | Discard a draft | 🔄 In Progress |
| 38 | GET | `/api/v1/applications` | List all applications for the signed-in user | 🔄 In Progress |
| 39 | GET | `/api/v1/applications/{id}` | Get full application record | 🔄 In Progress |
| 40 | POST | `/api/v1/applications/{id}/prefilled-pdf` | Generate and return pre-filled PDF | ⬜ Not Started |
| 41 | POST | `/api/v1/applications/{id}/signed-form` | Upload signed form PDF | 🔄 In Progress |
| 42 | GET | `/api/v1/applications/{id}/signed-form` | Download signed form PDF | 🔄 In Progress |

---

## Applications — Payment (`/api/v1/applications/{id}/payment`)

| # | Method | Endpoint | Description | Status |
|---|---|---|---|---|
| 43 | POST | `/api/v1/applications/{id}/payment/initiate` | Start an online payment session (card/mobile money) | ⬜ Not Started |
| 44 | POST | `/api/v1/applications/{id}/payment/proof` | Upload scanned proof of payment | 🔄 In Progress |
| 45 | GET | `/api/v1/applications/{id}/payment` | Get current payment status and details | 🔄 In Progress |

---

## Applications — Submission & Tracking

| # | Method | Endpoint | Description | Status |
|---|---|---|---|---|
| 46 | POST | `/api/v1/applications/{id}/submit` | Final submission (requires signed form + payment) | 🔄 In Progress |
| 47 | PATCH | `/api/v1/applications/{id}/respond` | Respond to NEEDS_INFORMATION request | 🔄 In Progress |
| 48 | POST | `/api/v1/applications/{id}/attachments` | Upload a required attachment | 🔄 In Progress |
| 49 | GET | `/api/v1/applications/{id}/attachments` | List all attachments for an application | 🔄 In Progress |
| 50 | DELETE | `/api/v1/applications/{id}/attachments/{attachmentID}` | Remove an attachment | 🔄 In Progress |

---

## Transactions (`/api/v1/transactions`)

| # | Method | Endpoint | Description | Status |
|---|---|---|---|---|
| 51 | GET | `/api/v1/transactions` | List all payment transactions for signed-in user | 🔄 In Progress |
| 52 | GET | `/api/v1/transactions/{reference}` | Get a specific payment transaction by reference | 🔄 In Progress |

---

## Infrastructure & Configuration

| Component | Description | Status |
|---|---|---|
| Docker Compose | Multi-service setup (postgres, backend, nginx) | 🔄 In Progress |
| Dockerfile (backend) | Go multi-stage build | 🔄 In Progress |
| Nginx config | Reverse proxy + rate limiting + security headers | 🔄 In Progress |
| PostgreSQL migrations | Schema: users, roles, tokens, applications, audit | 🔄 In Progress |
| Seed data | System roles + permissions seeded | 🔄 In Progress |
| .env.example | All required environment variables documented | 🔄 In Progress |
| JWT authentication | Access token (15 min) + refresh token (7 days) | 🔄 In Progress |
| RBAC middleware | Role-based route guards | 🔄 In Progress |
| Rate limiting | Per-IP request rate limiting | 🔄 In Progress |
| CORS | Configurable allowed origins | 🔄 In Progress |
| Audit logging | All mutating actions logged | 🔄 In Progress |
| Postman collection | All 52 endpoints with auth + test scripts | 🔄 In Progress |
| PDF generation | Pre-filled statutory form PDF generation | ⬜ Not Started |
| Payment gateway | Online card / mobile money integration | ⬜ Not Started |
| File storage | Signed form and attachment storage (Supabase Storage) | ⬜ Not Started |
| Email notifications | Password reset + status change emails | ⬜ Not Started |

---

## Summary

| Category | Total | Complete | In Progress | Not Started |
|---|---|---|---|---|
| Auth | 7 | 0 | 7 | 0 |
| Dashboard | 1 | 0 | 1 | 0 |
| Users | 9 | 0 | 9 | 0 |
| Roles & Permissions | 8 | 0 | 8 | 0 |
| Admin Applications | 7 | 0 | 7 | 0 |
| Audit Logs | 2 | 0 | 2 | 0 |
| Application Wizard | 8 | 0 | 7 | 1 |
| Payment | 3 | 0 | 2 | 1 |
| Submission & Tracking | 5 | 0 | 5 | 0 |
| Transactions | 2 | 0 | 2 | 0 |
| **Total API** | **52** | **0** | **50** | **2** |
