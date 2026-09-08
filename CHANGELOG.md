# Changelog - National Council of Sports (NCS) Intranet & Online Platform

All notable technical and architectural changes to the NCS Online & Intranet Platform are documented in this file.

---

## [Unreleased] - 2026-09-08

### Added — Expenses
- Added independent expense register, entry, detail and administrator category pages, linked from every accounting-role dashboard.
- Record date, category, title/reason, UGX amount, payee, payment method/reference and department; the server records the authenticated officer's ID, full name and save time.
- Added administrator-only category creation, renaming and deactivation, historical name snapshots, private receipt uploads/downloads and atomic persistence in migration 079.
- Added PostgreSQL integration tests, browser workflow/access checks and [verification screenshots and change report](Docs/verification/expenses-2026-09-08/README.md). Expense changes are local and not yet deployed.

### Changed
- Replaced all four fixed-asset dialogs with independent pages: new asset, revaluation, physical verification and monthly depreciation.
- Added deep-link asset loading, inline persistent errors, required fields, decimal currency inputs and duplicate-submit protection. Successful saves return to the refreshed register.
- Integrated upstream accountant navigation and linked its new-asset tile directly to the standalone page. Fixed the upstream Command Center newline literal that blocked the frontend build.
- Added [page workflow verification and screenshots](Docs/verification/asset-pages-2026-09-08/README.md).

### Fixed
- Added Fixed Asset Register navigation to the active shared dashboard layout for accounting, audit, GS and administrator roles.
- Added bounded register pagination and filter/search offset reset so all 297 baseline assets can be browsed.
- Removed hardcoded live portfolio/count fallbacks and the unverified complete-replacement claim.
- Pinned row actions, corrected search spacing, added search/category accessible labels, allowed header wrapping and cleared the notice timer on unmount.
- Preserved the existing API response-envelope fix in `FixedAssetsView.vue`.

### Verification
- Reconciled 297 workbook/seed/local PostgreSQL records; documented the unresolved hash FB-cost cell and stale cached workbook pivot.
- Frontend build, existing Go tests, fixture browser checks and authenticated local accountant read-only checks passed. Added desktop/mobile screenshot evidence.
- All observed local container IDs and asset-data export hashes remained unchanged. No container rebuild, data migration or VPS deployment was performed.
- **Complete replacement acceptance remains failed:** missing workflows and critical authorization/depreciation/approval gaps are recorded in [the bug report](Docs/verification/asset-register-2026-09-08/BUG_REPORT.md). Remote verification is blocked by connectivity/SSH details.

## [1.3.0] - 2026-08-14

### Added
#### 1. Comprehensive User Seeding & Role Testing Matrix
- `backend/migrations/073_seed_all_departmental_users.sql`: Seeded 30 dedicated user accounts across every organizational use case (Super Admin, General Secretary, AGS-T, AGS-A, Technical, Senior Engineer, Civil/Electrical Engineers & Officers, Plumber, Accountant, Auditor, HR, Helpdesk, IT Officer, PDU Procurement, PR, Stores Officer, Facilities Manager, Legal Counsel, Medical Officer, Physiotherapist, Transport Officer, Driver, Federation President, Federation GS, Safeguarding Officer, Content Manager, and Applicant).
- `Docs/Credentials.csv`: Comprehensive CSV catalog containing all 30 user credentials, roles, designations, emails, passwords, access levels, and primary dashboard URLs.

#### 2. Podman Container Orchestration & Testing Setup
- `podman-compose.yml`: Tailored Podman compose specification defining `postgres` (with auto-migration on mount), `backend` (Go REST API), `frontend` (Vue 3 Nginx build), `location-service` (Geo-IP validation), and `nginx` gateway (`9081` HTTP, `9444` HTTPS).
- `podman-setup.ps1`: Automated PowerShell deployment script for Windows environments (checks Winget, Podman CLI, machine status, `.env` file, and initiates stack build).
- `podman-setup.sh`: Automated bash deployment script for Linux/WSL/Server environments.
- `Docs/PODMAN_SETUP.md`: Comprehensive setup and role testing guide detailing Podman installation on Windows/Linux, service port mappings, and testing workflows.
- `.env`: Development environment configuration file configured with ready-to-run database passwords, JWT secrets, and service endpoints.

#### 3. Native Windows PostgreSQL & Go Backend Live Runtime
- Installed and initialized standalone PostgreSQL 17 cluster at `D:\pgsql\data` with UTF-8 encoding.
- Created `ncsintranet` database, `ncsintranet_user` role, and executed all 43 database migrations.
- Compiled `backend/ncs-backend.exe` and launched REST API listening on `http://localhost:9081`.
- Verified live end-to-end authentication, JWT issuance, and RBAC across Super Admin, General Secretary, AGS Technical, AGS Admin, Stores, Transport, Medical, Legal, and Facilities roles.
- Frontend active at `http://localhost:3001` communicating directly with the Go backend API.

#### 4. Complete Frontend Rebuild with Otika Admin Design System
- Replaced legacy hybrid layout with unified master **Otika Layout Shell** (`frontend/src/components/layout/LayoutDefault.vue`).
- Rebuilt `frontend/src/views/LoginView.vue` as a clean Otika Card Login with persistent Light/Dark mode toggle.
- Unified all dashboards inside `LayoutDefault.vue` and verified clean production build.

#### 5. Minimalist Enterprise Otika Styling & 100% Backend API Verification
- Refactored all 9 executive and operational dashboard views to remove harsh dark gradients, neon borders, and clashing saturated shading:
  - [`GeneralSecretaryAppraisalView.vue`](file:///d:/NCSOnline/frontend/src/views/executive/GeneralSecretaryAppraisalView.vue) (`/executive/appraisals`)
  - [`AssistantGeneralSecretaryAdminDashboard.vue`](file:///d:/NCSOnline/frontend/src/views/executive/AssistantGeneralSecretaryAdminDashboard.vue) (`/executive/ags-admin`)
  - [`AssistantGeneralSecretaryTechnicalDashboard.vue`](file:///d:/NCSOnline/frontend/src/views/executive/AssistantGeneralSecretaryTechnicalDashboard.vue) (`/executive/ags-technical`)
  - [`StoresInventoryDashboard.vue`](file:///d:/NCSOnline/frontend/src/views/stores/StoresInventoryDashboard.vue) (`/stores/inventory`)
  - [`FacilitiesManagementDashboard.vue`](file:///d:/NCSOnline/frontend/src/views/facilities/FacilitiesManagementDashboard.vue) (`/facilities/venues`)
  - [`FleetTransportDashboard.vue`](file:///d:/NCSOnline/frontend/src/views/fleet/FleetTransportDashboard.vue) (`/fleet/transport`)
  - [`SportsMedicalDashboard.vue`](file:///d:/NCSOnline/frontend/src/views/medical/SportsMedicalDashboard.vue) (`/medical/sports-science`)
  - [`LegalComplianceDashboard.vue`](file:///d:/NCSOnline/frontend/src/views/legal/LegalComplianceDashboard.vue) (`/legal/compliance`)
  - [`InfrastructureCommandCenterView.vue`](file:///d:/NCSOnline/frontend/src/views/InfrastructureCommandCenterView.vue) (`/maintenance/command-center`)
- Applied standardized Otika `.card-statistic-4` KPI cards, minimalist top action bars, and striped tables with balanced light/dark mode contrast.
- Fixed migration type mismatches (`UUID` vs `TEXT` for `users.id` foreign keys) in migrations `065`, `066`, `067`, and `068`.
- Created and seeded all 9 tables in PostgreSQL (`executive_appraisals`, `asset_valuation_signoffs`, `financial_efficiency_appraisals`, `technical_approvals`, `facility_readiness_certifications`, `admin_approvals`, `administrative_directives`, `store_inventory_items`, `goods_received_notes`).
- Verified 100% live status (21/21 endpoints returning HTTP 200 OK).

---

## [1.2.0] - 2026-08-14

### Added
#### 1. Architecture Blueprints & Plans (`Docs/plans/`)
- `assistant-general-secretary-technical-plan.md`: Comprehensive master executive blueprint for AGS-T (Sports Federations, NSMIS, Athletes, Engineering, Facilities & Match Readiness).
- `assistant-general-secretary-admin-plan.md`: Comprehensive master executive blueprint for AGS-A (HR, Finance, PDU Procurement, PR, IT, Directives).
- `general-secretary-appraisal-valuation-plan.md`: Master 4-pillar executive appraisal center for General Secretary (Property/Inventory, Staff/Departmental, Financial Efficiency, and UGX 31.02B Fixed Asset Valuation).
- `stores-inventory-department-plan.md`: Blueprint for Goods Received Notes (GRN), Bin Cards, Store Issue Vouchers (SIV), and sports equipment pools.
- `facilities-venue-management-plan.md`: Blueprint for multi-facility bookings, commercial NTR tariff billing, pre/post-event inspections, and Lugogo Hostels.
- `legal-compliance-department-plan.md`: Blueprint for statutory compliance (Act 2023), contracts vault, and federation dispute arbitrations.
- `sports-science-medical-antidoping-plan.md`: Blueprint for athlete medical screening, injury surveillance, and WADA/RADO anti-doping records.
- `fleet-transport-management-plan.md`: Blueprint for vehicle fleet telematics (King Long Bus, Ford Ranger, Kia Sorento, Motorcycles), trip gate-passes, and fuel logs.
- Updated `general-secretary-dashboard-plan.md` and `all-departments-framework-plan.md`.

#### 2. Database Migrations (`backend/migrations/`)
- `065_create_agst_technical_approvals.sql`: Created `technical_approvals` and `facility_readiness_certifications` tables; seeded `role_ags_technical` and permissions (`agst:read`, `agst:write`).
- `066_create_agsa_admin_approvals.sql`: Created `admin_approvals` and `administrative_directives` tables; seeded `role_ags_admin` and permissions (`agsa:read`, `agsa:write`).
- `067_create_gs_appraisals.sql`: Created `executive_appraisals`, `asset_valuation_signoffs`, and `financial_efficiency_appraisals` tables with UGX 31.02B baseline asset ledger data across all 11 classes.
- `068_create_stores_inventory.sql`: Created `store_inventory_items` and `goods_received_notes` tables; seeded `role_stores_officer`.
- `069_create_facilities_venue_management.sql`: Created `venue_bookings` and `hostel_occupancies` tables; seeded `role_facilities_manager`.
- `070_create_legal_compliance.sql`: Created `legal_contracts` and `federation_disputes` tables; seeded `role_legal_counsel`.
- `071_create_sports_medical.sql`: Created `athlete_medical_screenings`, `athlete_injuries`, and `antidoping_records` tables; seeded `role_medical_officer` and `role_physiotherapist`.
- `072_create_fleet_transport.sql`: Created `fleet_vehicles` and `trip_requisitions` tables; seeded `role_transport_officer` and `role_driver`.

#### 3. Backend Go REST API Handlers (`backend/internal/handlers/` & `main.go`)
- `backend/internal/handlers/agst.go`: Handlers for `/api/v1/executive/ags-t/*` (Dashboard stats, technical approvals list/action, facility readiness certs).
- `backend/internal/handlers/agsa.go`: Handlers for `/api/v1/executive/ags-a/*` (Dashboard stats, admin vetting approvals list/action, directives).
- `backend/internal/handlers/appraisal.go`: Handlers for `/api/v1/executive/appraisal/*` (Appraisal summary, asset ledger sign-offs, financial efficiency).
- `backend/internal/handlers/stores.go`: Handlers for `/api/v1/stores/*` (Inventory listing, create item, Goods Received Notes).
- `backend/internal/handlers/facilities.go`: Handlers for `/api/v1/facilities/*` (Venue bookings listing/creation, hostel occupancies).
- `backend/internal/handlers/legal.go`: Handlers for `/api/v1/legal/*` (Contracts vault, federation disputes).
- `backend/internal/handlers/medical.go`: Handlers for `/api/v1/medical/*` (Screenings, injuries, anti-doping sample results).
- `backend/internal/handlers/fleet.go`: Handlers for `/api/v1/fleet/*` (Fleet vehicles, trip requisitions).
- `backend/internal/handlers/handlers.go` & `backend/cmd/server/main.go`: Wired new handlers and registered all route groups.

#### 4. Frontend Vue 3 Dashboards & Routing (`frontend/src/`)
- `AssistantGeneralSecretaryTechnicalDashboard.vue` (`/executive/ags-technical`): Executive technical radar, 50+ federations monitor, and venue readiness certification desk.
- `AssistantGeneralSecretaryAdminDashboard.vue` (`/executive/ags-admin`): Executive administrative radar, Form 5 vetting queue, HR establishment & payroll pre-authorization.
- `GeneralSecretaryAppraisalView.vue` (`/executive/appraisals`): GS statutory appraisal suite for UGX 31.02B Fixed Asset Ledger, properties, staff scorecards, and financial efficiency.
- `StoresInventoryDashboard.vue` (`/stores/inventory`): Electronic bin cards, stock valuation, and Goods Received Notes (GRN).
- `FacilitiesManagementDashboard.vue` (`/facilities/venues`): Venue booking master calendar, NTR billing, and Lugogo Hostels camp check-ins.
- `LegalComplianceDashboard.vue` (`/legal/compliance`): Commercial contracts repository, arbitration tribunal cases, and statutory compliance.
- `SportsMedicalDashboard.vue` (`/medical/sports-science`): Pre-competition medical screening passes, injury surveillance, and WADA clean sport records.
- `FleetTransportDashboard.vue` (`/fleet/transport`): Vehicle master registry, trip requisitions, gate-passes, and fuel tracking.
- `frontend/src/router/index.js`: Registered 8 new routes with role-based metadata.
- `frontend/src/components/layout/AppSidebar.vue`: Added **Executive Command** and **Department Portals** navigation sections.

#### 5. Persistent Memory State
- `memory.md`: Created and updated with architectural memory, fixed asset baselines, and implementation statuses.
