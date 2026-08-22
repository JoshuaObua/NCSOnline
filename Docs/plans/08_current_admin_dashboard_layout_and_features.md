# Existing Admin Dashboard Layout, Sidebar Navigation & Feature Specification

## 1. Overview
This document provides a comprehensive architectural specification of the **Admin Dashboard** in `ncsportal`. It details the existing UI layout, top navbar controls, sidebar menu hierarchy, dropdown groups, modular views, component structures, dynamic form builder capabilities, security RBAC panels, and backend infrastructure controls implemented in the frontend application (`UserPortalView.vue`, `WebsiteContentManagerView.vue`, `NamisManagerPanel.vue`, `AdminApplicationsPanel.vue`, `OpenFormsPanel.vue`).

---

## 2. Global UI Shell & Top Navbar Controls

```
+-----------------------------------------------------------------------------------+
|  [≡ Toggle] [↻ Refresh]  NCS PORTAL ADMIN DASHBOARD      [🌙 Dark] [✉ 3] [🔔 5] [👤 Profile v] |
+-----------------------------------------------------------------------------------+
```

### Top Navbar Features:
1. **Sidebar Collapse Button (`<i class="icofont-navigation-menu">`)**: Toggles collapsed/expanded sidebar mode (`sidebar-mini`).
2. **Dashboard Refresh Button (`<i class="icofont-refresh">`)**: Triggers an asynchronous reload of all active section data, user KPIs, and notifications.
3. **Workspace Context Indicator**: Displays current user role label (e.g. *NCS Super Administrator*, *General Secretary*, *Federation Executive*) and active section title.
4. **Theme Toggle Component (`<ThemeToggle>`)**: Switcher between Dark Mode and Light Mode.
5. **Messages Quick Action (`<i class="icofont-envelope">`)**: Badge indicator showing unread message count; opens dropdown/drawer with reviewer messages.
6. **Notifications Quick Action (`<i class="icofont-notification">`)**: Badge indicator displaying unread system notifications with a "Mark all as read" control.
7. **User Profile Dropdown Menu**:
   -  **My Profile**: User profile details, avatar upload, email 2FA settings, and password change.
   -  **My Applications**: Historical list of submitted applications.
   -  **Transactions**: Payment receipts, reference codes, and UGX payment proof uploads.
   -  **Logout**: Session termination.

---

## 3. Sidebar Navigation Menu & Dropdown Structure

The sidebar navigation dynamically adapts based on user privileges and active roles. Below is the complete hierarchy of menu items and submenus present in the Admin Dashboard:

```
+-----------------------------------------------------------------------------------+
| [Sidebar Navigation Menu]                                                         |
|                                                                                   |
| ── MAIN ADMINISTRATION ────────────────────────────────────────────────────────── |
|  Overview                                                                       |
|  Open Applications / Apply Now                                                  |
|  My Applications                                                                |
|  Credential Wallet (My Files)                                                   |
|                                                                                   |
| ── NAMIS SPORTS REGISTRY (NAMIS MANAGER) ──────────────────────────────────────── |
|  Athlete Analytics & Demographics                                               |
|  Athletes Master Registry                                                        |
|  Accredited Sports Clubs                                                        |
|  Licensed Coaches Roster                                                         |
|  Multidisciplinary Support Entourage                                             |
|  National Squad Appearances                                                     |
|  Competitions Registry                                                          |
|  Athlete Results & Scoring                                                      |
|  Medals & Honors Standings                                                      |
|  Talent Identification Pipeline                                                 |
|  Sports Scholarships & Grants                                                   |
|                                                                                   |
| ── SPORTS FEDERATIONS DESK ────────────────────────────────────────────────────── |
|  Sports Federations [v]                                                         |
|    ├── Federations List & Directory                                               |
|    ├── Add New Federation (Modal)                                                 |
|    ├── Statutory Category Classifications                                         |
|    ├── Add Category Classification (Modal)                                        |
|    └── Federation Recognition & Licensing Desk                                    |
|                                                                                   |
| ── DYNAMIC APPLICATION FORMS & SUBMISSIONS MANAGER ───────────────────────────── |
|  Dynamic Application Forms Manager [v]                                           |
|    ├── Create New Application Form (Multi-Step Form Builder)                       |
|    ├── Form Schemas Directory (View / Edit / Delete)                               |
|    ├── Form Field Configuration & Rules (Text, Number, Date, Select, File Upload)  |
|    ├── Form Category & Publication Controls (Active, Draft, Archived)               |
|    └── Application Submissions Review & Approval Desk                             |
|                                                                                   |
| ── SECURITY & USER MANAGEMENT ──────────────────────────────────────────────────────── |
|  Users & RBAC Accounts Directory                                                |
|  Roles & Permissions Matrix                                                     |
|                                                                                   |
| ── SYSTEM INFRASTRUCTURE ──────────────────────────────────────────────────────── |
|  Command Center & Maintenance                                                    |
|    ├── Flush Cache & Buffer Operations                                            |
|    ├── Database Migration Sync                                                    |
|    ├── Search Registry Re-indexer                                                 |
|    └── PostgreSQL Database Backups & Snapshots                                    |
|                                                                                   |
| ── WEBSITE CONTENT MANAGEMENT (CMS) ────────────────────────────────────────────── |
|  CMS Content Manager                                                            |
|    ├── Homepage & Hero Slideshow Manager                                          |
|    ├── Blog Posts & Categories                                                    |
|    ├── Static Pages & Document Resources                                          |
|    ├── Events & Facility Management                                               |
|    ├── Press Releases, Reports & Speeches                                         |
|    └── System Audit Logs & Site Settings                                          |
+-----------------------------------------------------------------------------------+
```

---

## 4. Detailed Feature Modules & Component Breakdown

### 4.1 Dashboard Overview (`section === 'dashboard'`)
- **Executive KPI Cards**: Real-time summary cards displaying:
  - *Total Applications*: Aggregate volume of submitted applications.
  - *Pending Review*: Applications awaiting officer review.
  - *Approved Applications*: Officially cleared applications.
  - *Total Payments Recorded*: Financial UGX totals.
- **Verified Athlete Summary Panel**: Displayed for athlete-linked user accounts, showing athlete ID (`NAMIS-UG-xxxx`), discipline, age category, district/region, club, license status, national duty caps, medical/safeguarding status, WADA compliance, and medal standings.
- **Open Forms Quick Access Panel (`<OpenFormsPanel>`)**: Displays published application forms ready for submission.

---

### 4.2 NAMIS Sports Registry Manager (`<NamisManagerPanel>`)
The NAMIS Manager Panel (`NamisManagerPanel.vue`) provides generic CRUD operations, tabular data grids, search/filtering controls, and slide-out drawer forms across 11 core domain modules:

1. **Athlete Analytics (`analytics`)**:
   - Registered athletes count breakdown.
   - Gender participation distribution (Female / Male counts).
   - Regional sports development statistics bar charts (Central, North, East, West).
   - Top 5 Federations by athlete volume.
2. **Athletes Master Registry (`athletes`)**:
   - Columns: `Athlete Name`, `NIN / Passport`, `Gender`, `DOB`, `Age Category`, `District`, `Region`, `Club`, `Discipline`, `Status`.
   - Actions: Search, Filter by Region/Gender/Status, Add Athlete, Edit Record, Delete Entry.
3. **Clubs Registry (`clubs`)**:
   - Columns: `Club Name`, `Acronym`, `Federation`, `Contact Person`, `Email`, `Phone`, `District`, `Region`, `Status`.
4. **Coaches Roster (`coaches`)**:
   - Columns: `Coach Name`, `Role`, `Certification Level`, `License Number`, `Expiry Date`, `Federation`, `Status`.
5. **Support Entourage (`athlete-support-entourage`)**:
   - Multidisciplinary entourage fields: Primary Coach, Assistant Coach, Strength & Conditioning Specialist, Sports Scientist, Physiotherapist, Team Doctor, Nutritionist, Team Manager.
6. **National Team Appearances (`national-team`)**:
   - Columns: `Athlete Name`, `Team Name` (e.g. *Uganda Cranes*, *She Cranes*), `Category` (*Senior*, *Development*, *Junior*), `First Call-Up Date`, `Appearances Count`.
7. **Competitions Registry (`competitions`)**:
   - Columns: `Competition Name`, `Level` (8 Tiers), `Host Country`, `Host City`, `Start Date`, `End Date`.
8. **Athlete Results & Scoring (`competition-results`)**:
   - Columns: `Athlete Name`, `Competition`, `Event`, `Result Value`, `Position`, `National Record (NR)`, `Personal Best (PB)`, `Seasonal Best (SB)`.
9. **Medals Standings (`medals`)**:
   - Columns: `Athlete Name`, `Competition`, `Event`, `Medal Type` (*Gold*, *Silver*, *Bronze*), `Date Won`, `Prize Money (UGX)`, `NCS Recognition Status`.
10. **Talent Identification Pipeline (`talent`)**:
    - Columns: `Scouted Candidate`, `Age at Identification`, `School`, `District`, `Region`, `Scout Name`, `Talent Centre`, `Talent Category`, `Pathway`.
11. **Sports Scholarships (`scholarships`)**:
    - Columns: `Athlete Name`, `Institution`, `Scholarship Type`, `Grant Amount (UGX)`, `Start Date`, `End Date`, `Status`.

---

### 4.3 Sports Federations & Licensing Desk
- **Federation Directory (`section === 'federations'`)**:
  - Tabular list of recognized federations showing `Name`, `Acronym`, `Category`, `President`, `General Secretary`, `Certificate Serial #`, `Status`.
  - Filterable by federation name / acronym search with pagination controls.
  - Actions: **Add Federation Modal** (`showAddFederationModal`), **Edit Federation**, **Delete Federation**.
- **Federation Category Classifications (`section === 'federation-categories'`)**:
  - Manage Category A (Priority - High Impact), Category B (Established), Category C (Developing) classifications.
  - Set annual statutory grant funding caps (UGX) and definition criteria.
  - Actions: **Add Category Modal** (`showAddCategoryModal`), **Edit Category**.
- **Recognition & Licensing Desk (`section === 'federations-license'`)**:
  - Manage official recognition certificate serial numbers and license renewal dates.

---

### 4.4 Dynamic Application Forms Manager & Submissions Review Engine
- **Dynamic Application Forms Manager (`section === 'form-builder'`)**:
  - **Form Schemas Directory**: Tabular overview of all created application forms showing form title, category, issuing department, price (UGX or Free), publication state (`DRAFT`, `PUBLISHED`, `ARCHIVED`), and submission counts.
  - **Create New Application Form (`openCreateFormModal`)**: Form creation wizard allowing administrators to define new multi-step application forms.
  - **Form Schema CRUD Operations**:
    -**Create Form**: Define form title, slug, department, fee amount (UGX), instructions, and initial steps.
    -**Edit Form & Fields**: Dynamically add, modify, or reorder input fields (`short_text`, `long_text`, `number`, `date`, `dropdown`, `radio`, `checkbox`, `file_upload`), field labels, help texts, placeholders, required flags, and file upload rules.
    -  **Delete / Archive Form**: Safely archive or permanently remove outdated application forms.
    -  **Toggle Active Status**: One-click publication control to open or close forms for public/federation applications.
- **Application Submissions Review Desk (`<AdminApplicationsPanel>`)**:
  - Filter submitted applications by status (`Submitted`, `Under Review`, `Approved`, `Rejected`), department, or reference code.
  - Inspect submitted form responses and payment proof attachments (`<FilePreview>`).
  - Action Triggers: Update Application Status, Add Reviewer Comments, Approve / Reject application.

---


### 4.5 Security, Users & RBAC Directory (`section === 'users'`)
- **User Directory Table**: List of registered system users displaying `User ID & Name`, `Email Address`, `Assigned RBAC Designation`, `Account Status`, `Created Date`.
- **Search & Role Filtering**: Filter users by role (`Super Admin`, `General Secretary`, `Human Resources`, `Accountant`, `Procurement Officer`, `Federation President`, `Athlete`, `Coach`).
- **User Management Triggers**: Edit Role designation, Reset Password, Deactivate Account.

---

### 4.6 System Maintenance, Operations & Infrastructure (`section === 'command-center'`)
- **System Maintenance Operations**:
  - **Flush Cache & Clear Buffer**: Purges Redis session caches and transient query buffers (`flushSystemCacheAction`).
  - **Synchronize Database Migrations**: Executes outstanding schema migrations across PostgreSQL.
  - **Re-index Search Registry**: Re-indexes athlete, federation, and application search indices.
- **PostgreSQL Database Backups & Snapshots**:
  - Automated & manual snapshot creator (**Backup Now**).
  - Snapshot file log displaying backup filename (`id`), creation date, and file size.
  - One-click backup file download trigger.

---

## 5. Summary Table of Frontend Components & Views

| Module / View File | Primary Responsibilities | Associated Backend Endpoints |
| :--- | :--- | :--- |
| `UserPortalView.vue` | Main Shell Layout, Sidebar Navigation, Header Controls, Profile & User Management | `/api/v1/auth/me`, `/api/v1/users` |
| `NamisManagerPanel.vue` | NAMIS Registry Manager, Athletes, Clubs, Coaches, Results, Medals, Talent ID | `/api/v1/nsmis/*` |
| `AdminApplicationsPanel.vue` | Application Submissions Review Desk, File Previewer, Status Transitions | `/api/v1/admin/applications` |
| `OpenFormsPanel.vue` | Public/Portal Open Application Form Cards & Application Wizard launcher | `/api/v1/forms/active` |
| `WebsiteContentManagerView.vue` | Full CMS Website Management (Homepage, Blog, Events, Speeches, Audit Logs) | `/api/v1/cms/*` |
