# NCSMS v1 — Changelog

All notable changes to this project are documented here.
Format: `[YYYY-MM-DD] — Type: Description`
Types: `Added` | `Changed` | `Fixed` | `Removed` | `Security` | `Infrastructure`

---

## [Unreleased]

### Added
- Backend Go API skeleton with Chi router, JWT auth, and RBAC middleware
- Docker Compose multi-service setup (postgres, backend, nginx)
- PostgreSQL migrations 001–006 (users, roles, tokens, applications, audit logs, seeds)
- Nginx reverse proxy configuration with rate limiting and security headers
- Postman collection (v2.1) with all 52 endpoints grouped by domain
- `progress-checklist.md` tracking all 52 API endpoints and infrastructure status

---

## [2026-05-26] — Design: Application Workflow Architecture

### Added
- `application-workflow-v1.md` — shared 6-stage wizard workflow architecture for all application forms
  - Stage 1: Multi-step form fill with auto-save, Save Draft, and Resume
  - Stage 2: Read-only review with per-section edit links
  - Stage 3: Pre-filled PDF generation + digital or print-sign-scan signature flow
  - Stage 4: Dual payment methods (online digital payment + scanned proof upload) for all applications and renewals
  - Stage 5: Pre-submission checklist + application reference assignment
  - Stage 6: Status tracking dashboard with NEEDS_INFORMATION response flow
- Status state machine: `DRAFT → PENDING_SIGNATURE → PENDING_PAYMENT → SUBMITTED → UNDER_REVIEW → APPROVED / REJECTED / NEEDS_INFORMATION`
- Frontend component map (ApplicationWizard, PaymentSelector, ProofOfPaymentForm, etc.)
- Full API endpoint list for draft, signed form, payment, submission, and tracking

### Changed
- All 8 application form design files updated with `## Wizard Steps` section mapping PDF parts to numbered wizard steps
- `applications-v1.md` index updated with workflow overview, payment methods table, and status reference
- `README.md` updated with workflow summary and full folder structure tree

---

## [2026-05-26] — Design: PDF Form Integration

### Added
- All 7 official NCS statutory PDF forms added to `Modular Design/applications/`
- Each application markdown file completely rewritten with exact field names (snake_case), input types, field labels, and attachment requirements derived directly from the PDFs:
  - `national-sport-recognition-application-v1.md` — Form 1, Reg. 3(1)
  - `federation-registration-renewal-v1.md` — Form 3, Reg. 4(1), 5(1), 9(1)
  - `national-sports-association-transformation-v1.md` — Form 5, Reg. 11(1)
  - `national-sports-federation-transformation-v1.md` — Form 5, Reg. 11(1)
  - `sports-competition-organization-application-v1.md` — Form 7, Reg. 16(1)
  - `sports-facility-operation-application-v1.md` — Form 8, Reg. 17(2)
  - `community-academy-registration-renewal-v1.md` — Form 10, Reg. 22(1)
  - `sports-academy-operation-application-v1.md` — Form 11, Reg. 29(1)(a)
- `applications-v1.md` index updated with official form numbers, regulation references, and attachment cross-reference table

---

## [2026-05-26] — Design: Initial Module Design Commit

### Added
- `Modular Design/ncsms-v1-module-design.md` — core architecture, security, middleware, and shared guidance
- `Modular Design/super-admin-v1.md` — Super Admin module design
- `Modular Design/general-secretary-v1.md` — General Secretary application review and licence processing
- `Modular Design/roles-v1.md` — system roles and custom role permissions
- `Modular Design/user-profiles-v1.md` — user profile and account categories
- `Modular Design/athletes-v1.md` — athlete registration and credentialing
- `Modular Design/assets-v1.md` — asset and inventory management
- `Modular Design/audit-v1.md` — audit events and compliance
- `Modular Design/content-manager-v1.md` — website content and public pages
- `Modular Design/applications/` folder with all application form stubs

---

## Roadmap

### Phase 1 — Backend MVP (Current)
- [ ] Go API server with Chi router
- [ ] JWT auth + refresh token rotation
- [ ] RBAC middleware (super_admin, admin, general_secretary, user)
- [ ] User CRUD
- [ ] Role and permission management
- [ ] Application CRUD + status workflow
- [ ] Admin review endpoints (approve, reject, request-info)
- [ ] Payment proof upload
- [ ] Audit logging
- [ ] Dashboard stats
- [ ] PostgreSQL migrations + seeds

### Phase 2 — Integrations
- [ ] Payment gateway (card + mobile money)
- [ ] Pre-filled PDF generation (statutory forms)
- [ ] File storage (Supabase Storage)
- [ ] Email notifications (SMTP / SendGrid)

### Phase 3 — Frontend
- [ ] Admin dashboard (React / Next.js)
- [ ] Applicant portal with 6-stage wizard
- [ ] Application tracking dashboard
- [ ] Public website + content pages

### Phase 4 — Production Hardening
- [ ] SSL / TLS certificates (Let's Encrypt)
- [ ] Redis-backed rate limiting and session management
- [ ] Monitoring and alerting (Prometheus + Grafana)
- [ ] Automated backups
- [ ] CI/CD pipeline
