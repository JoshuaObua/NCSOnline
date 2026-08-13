# General Secretary Executive Master Dashboard Blueprint & Implementation Plan

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/general-secretary-dashboard-plan.md`  
> **Role:** General Secretary (GS) / Accounting Officer - National Council of Sports (NCS)  
> **Authority Level:** Organization-Wide Executive Control & Final Statutory Approval Authority  
> **Target Access Scope:** Unrestricted, real-time read, report query, and statutory approval access across ALL NCS departments (Engineering, Finance, HR, IT, Technical/Sports, Procurement, Public Relations, Internal Audit)  

---

## 1. Executive Role Overview & Master Interconnected Architecture

The **General Secretary (GS)** is the Chief Executive Officer and statutory Accounting Officer of the National Council of Sports (NCS). Under the National Sports Act (2023) and Public Finance Management Act (PFMA 2015), the General Secretary carries total organizational accountability for financial disbursements, asset management, human resources, sports federation governance, procurement approvals, infrastructure projects, and public communications.

This plan details the **General Secretary Executive Master Dashboard** inside `NCS_Intranet`. It functions as an **Interconnected Executive Control Center** where all departmental reports, category analytics, submittals, and employee 360° profiles stream directly into the GS portal in real time.

```
+---------------------------------------------------------------------------------------------------+
|               GENERAL SECRETARY (GS) MASTER INTERCONNECTED DATA ENGINE                            |
+---------------------------------------------------------------------------------------------------+
       |                  |                |               |               |                |
       v                  v                v               v               v                v
+--------------+  +---------------+  +-----------+  +------------+  +--------------+  +---------------+
| Engineering  |  | Finance & Acc |  |   IT/ICT  |  | HR & Admin |  | Technical    |  | Procurement   |
| Field Orders |  | Asset Register|  | Backups   |  | Employee   |  | Federations  |  | Form 5s &     |
| Asset Health |  | Ledger & Grants| | Health    |  | Profiles   |  | Athletes     |  | APP Tracker   |
+--------------+  +---------------+  +-----------+  +------------+  +--------------+  +---------------+
       |                  |                |               |               |                |
       +------------------+----------------+---------------+---------------+----------------+
                                                   |
                                                   v
                     +-----------------------------------------------------------+
                     | GS EXECUTIVE DASHBOARD & STATUTORY APPROVAL QUEUE         |
                     | - PPDA Form 5 Final Approvals (Accounting Officer)        |
                     | - Organization-Wide Departmental Reports & Analytics      |
                     | - Employee 360° Profile Inspector (All Departments)        |
                     | - CapEx Escalations > UGX 5M & Asset Write-Offs          |
                     | - Executive PDF Board & Ministry Report Generator         |
                     +-----------------------------------------------------------+
```

---

## 2. Deep-Dive Executive Modules & Key Capabilities

---

### 2.1 Cross-Departmental Master Reporting Engine (By Department & Category)

The General Secretary can query, view, filter, and export general reports across any department or functional category:

1. **Departmental Filter Query:**
   - Filter reports by: `Engineering`, `Finance & Accounts`, `Human Resources`, `IT / ICT`, `Technical & Sports`, `Procurement (PDU)`, `Public Relations`, `Internal Audit`.
2. **Category Filter Query:**
   - **Financial Category:** Subvention execution, federation grant disbursements, NTR collection, budget vote-head commitments.
   - **Fixed Assets Category:** Live asset valuation (297 records / UGX 31.02B), revaluation requests, depreciation impact, condition state.
   - **Personnel & HR Category:** Staff headcount, active leave rosters, monthly payroll summaries (PAYE/NSSF), appraisal score distribution.
   - **Engineering Category:** Facility readiness (Lugogo Stadium, Hostels, Tennis Complex), active work orders, water/electrical task queues.
   - **Procurement Category:** PPDA Form 5 requisitions, Annual Procurement Plan (APP) execution rate, active contract awards.
   - **Technical Sports Category:** 50+ national federation governance scores, athlete licensing counts, international game delegations.
   - **IT & Security Category:** Server uptime, database backup logs, open IT support SLAs, system security audit events.

---

### 2.2 Organization-Wide Employee 360° Profile Inspector

The General Secretary has unrestricted authority to inspect the 360° master profile of **ANY employee across all departments**:

1. **Employee Search & Audit Bar:** Search staff by Name, Staff ID, Department, NIN, or Position.
2. **Comprehensive Staff Inspector View:**
   - **Bio & Identity:** Full legal name, photo avatar, NIN, Passport, DOB, Gender, Next of Kin, emergency contacts.
   - **Employment Status:** Designation, Department, Duty Station, Salary Scale (Scale 1-10), Employment Terms (Permanent, Contract, Seconded, Casual), Appointment Date, Contract Expiry Date.
   - **Payroll & Financial Identifiers:** URA TIN, NSSF number, Bank Name, Account Number, Gross Salary, monthly deductions.
   - **Appraisal & Performance Score:** Semi-annual appraisal ratings, KPI achievements, PIP flags, commendations.
   - **Leave & Duty Record:** Remaining leave balance, active leave requests, historical leave log.
   - **Individual Activity Stream:** Timestamped audit trail of all actions performed by the employee in the intranet (`MyActivitiesView.vue`).

---

### 2.3 Executive Statutory Approval Queue

Centralized approval hub where the General Secretary exercises statutory sign-offs:

1. **PPDA Procurement Form 5 Approvals:** Final Accounting Officer statutory approval for departmental procurement requisitions.
2. **CapEx & Infrastructure Escalations (> UGX 5M):** Senior Engineer CapEx requisitions for venue repairs, floodlight overhauls, or civil works.
3. **Fixed Asset Write-Offs & Disposals:** Authorizing statutory write-offs for damaged equipment or property revaluations.
4. **Federation Grant Disbursements:** Executive release sign-off for quarter funding to national sports associations.
5. **Monthly Payroll Authorization:** Final Accounting Officer approval of generated monthly payroll and bank EFT transfer files.

---

### 2.4 Executive PDF & Board Package Generator

1. **One-Click Board & Ministry PDF Generator:**
   - Instantly compiles formatted PDF executive briefs for:
     - Board of Directors Monthly Briefing Package.
     - Ministry of Education and Sports (MoES) Performance Report.
     - Ministry of Finance (MoFPED) Fixed Asset & Budget Execution Schedule.
     - Auditor General Uganda Statutory Compliance Pack.
2. **Automated Executive Summary Cover Pages:** Automatically appends high-level macro charts, total portfolio valuation (UGX 31.02B), headcount metrics (128 staff), and budget commitment rates.

---

## 3. General Secretary Dashboard User Interface (`GeneralSecretaryDashboard.vue`)

```
+-----------------------------------------------------------------------------------+
| [≡] NCS INTRANET | General Secretary Executive Command Portal   [🔍] [🔄] [🌙] [🔔 8] [GS] |
+-----------------------------------------------------------------------------------+
| [EXECUTIVE MACRO KPI SUMMARY CARDS]                                               |
| +--------------------+ +--------------------+ +-------------------+ +-------------------+ |
| | Annual Budget Spend| | Fixed Asset NBV    | | Active Staff      | | Form 5 Queue      | |
| | UGX 18.0B / 25.0B  | | UGX 31.02 Billion  | | 128 Employees     | | 8 Pending Sign-off| |
| +--------------------+ +--------------------+ +-------------------+ +-------------------+ |
+-----------------------------------------------------------------------------------+
| [GS WORKSPACE: (1) Statutory Approvals | (2) Master Reports | (3) Employee 360 | (4) Activity Feed] |
| +-------------------------------------------------------------------------------+ |
| | [ORGANIZATION-WIDE MASTER REPORT GENERATOR]                                   | |
| | Filter Dept: [ All Depts v ] Category: [ Financial & Fixed Assets v ] Date: [ Jul 2026 ]| |
| | ----------------------------------------------------------------------------- | |
| | Report Title                    | Origin Dept  | Submitted By     | Action        | |
| | Fixed Asset Register Revaluation| Finance      | CFO (Akello S)   | [ View PDF ]  | |
| | Lugogo Floodlight Site Inspection| Engineering | Senior Engineer  | [ View Log ]  | |
| | Q1 Federation Grant Accountabil.| Technical    | Tech Director    | [ Inspect ]   | |
| | [ Download Executive Board Package ] [ Export Ministry PDF ] [ Send to MoFPED ]| |
| +-------------------------------------------------------------------------------+ |
| +-----------------------------------------------+ +-------------------------------+ |
| | Employee 360° Profile Quick Inspector         | | Live Cross-Dept Activity Feed | |
| | Search Staff: [ Okello John - Senior Engineer ]| | - HR: Payroll generated (342.5M)|
| | Dept: Engineering | Status: Active Permanent  | | - PDU: Form 5 Submitted (45M) | |
| | Leave: 14 Days Left | Appraisal: 92% Exceeds  | | - IT: DB Backup Success (42MB) | |
| | [ Inspect Full 360° Profile ] [ View Audit Log]| | - PR: Press Statement Live   | |
| +-----------------------------------------------+ +-------------------------------+ |
+-----------------------------------------------------------------------------------+
```

---

## 4. Go REST API Mapping (`backend/internal/executive/gs_handler.go`)

| Endpoint | Method | Scope | Description |
| :--- | :--- | :--- | :--- |
| `/api/v1/executive/gs/dashboard` | GET | General Secretary | Macro KPIs across all departments |
| `/api/v1/executive/gs/reports/master` | POST | General Secretary | Master report query filtered by Dept, Category, and Date |
| `/api/v1/executive/gs/employees` | GET | General Secretary | Query 360° employee master profiles for any staff member |
| `/api/v1/executive/gs/employees/:id` | GET | General Secretary | View complete 360° staff profile, payroll, and activity log |
| `/api/v1/executive/gs/approvals` | GET/PUT | General Secretary | Statutory sign-offs (Form 5, CapEx > 5M, Write-offs, Payroll) |
| `/api/v1/executive/gs/export/board-pdf` | POST | General Secretary | Generate formatted Board & Ministry executive PDF brief |

---

## 5. Implementation Roadmap

- [ ] **Backend Executive Handlers:** Create `backend/internal/executive/gs_handler.go` with cross-departmental data aggregation queries.
- [ ] **Frontend View Construction:** Build `frontend/src/views/executive/GeneralSecretaryDashboard.vue` with executive KPI cards, master report query engine, employee 360 inspector, and statutory approval queues.
- [ ] **Interconnection Verification:** Verify that every departmental report, Form 5 submission, leave decision, asset adjustment, and employee profile streams cleanly to the General Secretary workspace.
