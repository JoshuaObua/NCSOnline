# NCS Online / Intranet - System Memory & Architecture State

## Overview
This document tracks the technical memory, architecture state, persistent decisions, and active module status of the National Council of Sports (NCS) Intranet platform.

---

## 1. Executive Structure & Statutory Mandates
- **Statutory Framework:** Uganda National Sports Act (2023), Public Finance Management Act (PFMA 2015), PPDA Act (2003), IPSAS 17.
- **Executive Leadership Hierarchy:**
  - **General Secretary (GS / Accounting Officer):** Final statutory approval for Form 5s, CapEx > UGX 5M, Asset Write-offs/Revaluations, Appraisals, Board & Ministry Reporting.
  - **Assistant General Secretary - Technical (AGS-T):** Direct oversight over Technical & Sports Administration, 50+ Federations, NSMIS, Athletes & Delegations, Engineering & Infrastructure, Facilities Match Readiness, and Sports Gear.
  - **Assistant General Secretary - Administration (AGS-A):** Direct oversight over Human Resources, Finance & Accounts, Procurement Unit (PDU), Public Relations, and IT/ICT Infrastructure.

---

## 2. Fixed Asset Register & Financial Baseline
- **Asset Ledger Baseline:** UGX 31,015,914,535.00 (~UGX 31.02 Billion) across 297 asset records across 11 Excel classes (`FIXED ASSET REGISTER ADJUSTMENTS.xlsx`).
- **Key Categories:** Land (UGX 27.89B), Non-Residential Buildings (UGX 2.14B), Light Vehicles (UGX 508.9M), Residential Buildings (UGX 167.25M), Electrical Machinery (UGX 107.4M), Furniture & Fittings (UGX 96.4M), Light ICT (UGX 59.8M), Office Equipment (UGX 27.7M), Cycles (UGX 9.38M), Other ICT (UGX 6.62M).

---

## 3. Implementation Phasing Tracker

| Phase / Plan | Module / Department | Backend Migration & Handler | Frontend View & Route | Status |
| :--- | :--- | :--- | :--- | :---: |
| **Phase 1** | Assistant General Secretary - Technical (AGS-T) | `065_create_agst_technical_approvals.sql` / `agst.go` | `AssistantGeneralSecretaryTechnicalDashboard.vue` (`/executive/ags-technical`) | ✅ Completed |
| **Phase 2** | Assistant General Secretary - Administration (AGS-A) | `066_create_agsa_admin_approvals.sql` / `agsa.go` | `AssistantGeneralSecretaryAdminDashboard.vue` (`/executive/ags-admin`) | ✅ Completed |
| **Phase 3** | General Secretary Master Appraisal & Valuation Center | `067_create_gs_appraisals.sql` / `appraisal.go` | `GeneralSecretaryAppraisalView.vue` (`/executive/appraisals`) | ✅ Completed |
| **Phase 4** | Stores & Inventory Management Unit | `068_create_stores_inventory.sql` / `stores.go` | `StoresInventoryDashboard.vue` (`/stores/inventory`) | ✅ Completed |
| **Phase 5** | Facilities Booking & Venue Operations | `069_create_facilities_venue_management.sql` / `facilities.go` | `FacilitiesManagementDashboard.vue` (`/facilities/venues`) | ✅ Completed |
| **Phase 6** | Legal, Compliance & Arbitrations | `070_create_legal_compliance.sql` / `legal.go` | `LegalComplianceDashboard.vue` (`/legal/compliance`) | ✅ Completed |
| **Phase 7** | Sports Science, Medical & Anti-Doping | `071_create_sports_medical.sql` / `medical.go` | `SportsMedicalDashboard.vue` (`/medical/sports-science`) | ✅ Completed |
| **Phase 8** | Fleet, Logistics & Transport Unit | `072_create_fleet_transport.sql` / `fleet.go` | `FleetTransportDashboard.vue` (`/fleet/transport`) | ✅ Completed |
| **Phase 9** | Comprehensive Role Seeding & Credentials | `073_seed_all_departmental_users.sql` | `Docs/Credentials.csv` | ✅ Completed |
| **Phase 10** | Podman Orchestration & Environment | `podman-compose.yml`, `podman-setup.ps1`, `podman-setup.sh` | `Docs/PODMAN_SETUP.md` | ✅ Completed |
| **Phase 11** | Native Live Windows Deployment & Verification | `D:\pgsql\pgsql\bin\pg_ctl`, `backend\ncs-backend.exe` | Frontend: `http://localhost:3001` / API: `http://localhost:9081` | ✅ Live & Operational |
| **Phase 12** | Complete Otika Admin Template UI Overhaul | `style.css`, `LayoutDefault.vue`, `ThemeToggle.vue`, `LoginView.vue` | Clean login layout without government tag, custom official copyright, built & live | ✅ Live & Operational |
| **Phase 13** | Minimalist Otika Styling & 21/21 API Endpoint Verification | All 9 Executive & Departmental Dashboards, PostgreSQL 065-068 migrations | Clean enterprise Otika styling, minimal coloring, 100% live API responses | ✅ Live & Operational |
| **Phase 14** | Online Appraisal & Departmental Reports Compilation Hub | `GeneralSecretaryAppraisalView.vue`, `appraisal.go`, `LayoutDefault.vue` | Renamed to Online Appraisal, removed ASGS links from GS dashboard, added categorized compiled reports hub | ✅ Live & Operational |
| **Phase 15** | Independent Departmental Report Pages & Sidebar Dropdown | 6 Dedicated Report Views in `views/reports/`, `LayoutDefault.vue` | Replaced tabs with independent pages, added sidebar dropdown menu, removed dummy data with 0 digits baseline | ✅ Live & Operational |
| **Phase 16** | Executive Leadership Title, Website Traffic Stats & Receptionist Clearance Desk | `ReceptionistVisitorView.vue`, `reception.go`, `dashboard.go`, `DashboardView.vue`, `073_create_reception_visitors.sql` | Executive Leadership & Governance menu header, live public portal analytics, full visitor clearance pass generation & printing | ✅ Live & Operational |
| **Phase 17** | Instant Signout Redirection & Scoped Front Desk Receptionist Portal | `LayoutDefault.vue`, `ReceptionistLeaveApplyView.vue`, `ReceptionistLeaveStatusView.vue`, `ReceptionistReportsView.vue`, `074_create_staff_leave_applications.sql` | Instant logout redirect without refresh, dedicated Front Desk menu (Visitors, Leave Apply, Leave Status, Reports) for Receptionist role | ✅ Live & Operational |
| **Phase 18** | Receptionist/Help Desk Tailored Dashboard & Menu Isolation | `DashboardView.vue`, `LayoutDefault.vue`, `router/index.js` | Removed unassigned departmental menus from Help Desk, implemented dedicated Front Desk KPI dashboard with visitor & leave metrics | ✅ Live & Operational |
| **Phase 19** | Full-Width Statutory "Apply for Leave" Dossier & Extended Schema | `ReceptionistLeaveApplyView.vue`, `reception.go`, `075_expand_staff_leave_fields.sql`, `LayoutDefault.vue` | Renamed to Apply for Leave, expanded to full content width, added all statutory establishment, handover, absence address & declaration fields | ✅ Live & Operational |
| **Phase 20** | Dedicated Visitor Clearance Intake View & Professional Input Styling | `VisitorInitiateView.vue`, `ReceptionistVisitorView.vue`, `router/index.js` | Converted modal into a dedicated full-width page at `/reception/visitors/new`, standard professional Otika textbox classes and streamlined navigation | ✅ Live & Operational |
| **Phase 21** | Searchable Select Dropdowns, Flat Focus Border Styling & UI Alignment | `SearchableSelect.vue`, `style.css`, `VisitorInitiateView.vue`, `ReceptionistVisitorView.vue`, `ReceptionistLeaveApplyView.vue`, `DashboardView.vue` | Removed focus highlight glow/shadows, created real-time searchable dropdowns, balanced button/element sizing, and aligned all Help Desk dashboard UI components | ✅ Live & Operational |
| **Phase 22** | IT Officer Workstation, PPDA Form 5 Vetting System & Scoped Navigation | `ITPPDAInitiateView.vue`, `ITPPDAListView.vue`, `it_officer.go`, `076_create_it_ppda_form5_requisitions.sql`, `077_expand_it_ppda_form5_fields.sql`, `DashboardView.vue`, `LayoutDefault.vue` | Complete PPDA Form 5 system with searchable designations across all NCS units, auto-detected department, searchable select for all dropdowns, extended procurement fields, and scoped sidebar | ✅ Live & Operational |
| **Phase 23** | Auto-Populated Personal & Departmental Particulars for Leave Application | `ReceptionistLeaveApplyView.vue` | Automatically detects authenticated applicant name, department, designation, IPPS file no, and contact details across all NCS roles with full searchable dropdowns | ✅ Live & Operational |
| **Phase 24** | Engineering Department Workstation & 5-Tier Statutory Approval Workflow | `DashboardView.vue`, `LayoutDefault.vue`, `ITPPDAListView.vue`, `it_officer.go`, `main.go` | Scoped Engineering dashboard & sidebar to PPDA Form 5 & Leave; multi-stage statutory workflow progression through HOD → AGS Technical → Finance → General Secretary (Accounting Officer) | ✅ Live & Operational |
| **Phase 25** | Specialized Engineering Officer Workstations (Civil, Electrical, Plumbing) | `DashboardView.vue`, `ITPPDAInitiateView.vue`, `ReceptionistLeaveApplyView.vue`, `LayoutDefault.vue` | Dedicated workstations for Senior Engineer, Assistant Civil Engineer, Assistant Electrical Engineer, Civil Officer, Electrical Officer, and Plumber with tailored titles, badges, auto-populated particulars, and shared 5-tier PPDA Form 5 workflow | ✅ Live & Operational |

---

## 4. API & Database Standards
- **REST API Path Prefix:** `/api/v1/*` (Listening on `http://localhost:9081`)
- **Database Engine:** PostgreSQL 17 at `127.0.0.1:5432` with database `ncsintranet` and UUID primary keys (`gen_random_uuid()`)
- **Frontend Stack:** Vue 3 (Composition API / `<script setup>`), Otika Admin Design System (Tailwind CSS + Otika UI Suite), Pinia, Vue Router, IcoFont (Listening on `http://localhost:3001`).
- **Authorization Standards:** RBAC with `role_permissions` and `users_roles` via JWT Bearer auth.
- **Testing & Local Runtime:** Both Podman orchestration and standalone native Windows runtime (`start-backend.ps1`) supported. Credentials stored in `Docs/Credentials.csv`.
