# NCS Enterprise All-Departments Blueprint & Universal Shell Architecture Plan

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/all-departments-framework-plan.md`  
> **System Scope:** National Council of Sports (NCS) Intranet - Standard Universal Dashboard Navigation & Shared User Workplace  
> **Target Departments:** IT/ICT, Finance & Accounts, Human Resources (HR), Technical & Sports Admin, Procurement, Engineering & Infrastructure, Public Relations, Internal Audit, Executive  
> **Universal Scope:** Provision for My Activities, Direct Messages, Real-Time Notifications, Profile Management, Security Settings, and Personal Leave Application & Approval Portal for EVERY User across all departments.

---

## 1. Enterprise Departmental Architecture & Universal Dashboard Shell

Every user account across all 6 departments and 20+ operational roles receives a standardized, premium responsive dashboard shell. Regardless of role (from Plumber and IT Technician to Senior Engineer, CFO, and General Secretary), every employee is equipped with their specialized operational tools AND the **Universal Shared Personal Workplace Suite**.

```
+---------------------------------------------------------------------------------------------------+
| [≡] NCS INTRANET | [Role / Department Title]   [🔍 Search Ctrl+K] [🔄] [🌙] [🔔 3] [💬 2] [User Avatar v] |
+---------------------------------------------------------------------------------------------------+
| [UNIVERSAL DASHBOARD NAVIGATION SIDEBAR]  | [MAIN DEPARTMENTAL WORKSPACE CONTENT AREA]            |
| ----------------------------------------- |                                                       |
| 📊 Dashboard (Role KPI Overview)          |  +-------------------------------------------------+  |
| 🏢 Department Workstation                 |  | SPECIALIZED DEPARTMENT MODULES                  |  |
| 📑 Requisitions / PPDA Form 5             |  | (e.g. Asset Ledger, Field Work Orders, NSMIS)   |  |
| 📅 My Leave & Roster Status               |  +-------------------------------------------------+  |
| 📜 My Activities Log                      |  +-------------------------------------------------+  |
| 💬 Direct Messages & Memos                |  | UNIVERSAL USER WORKPLACE MODULES                 |  |
| 👤 My Profile                             |  | - Leave Application & Approval Status           |  |
| ⚙️ Account & Security Settings            |  | - My Activities Stream | Direct Chat | Profile  |  |
| 🚪 Logout                                 |  +-------------------------------------------------+  |
+---------------------------------------------------------------------------------------------------+
```

---

## 2. Standardized Shared User Modules (Mandatory for ALL Users)

Every dashboard plan incorporates these 6 core personal workplace modules:

---

### 2.1 Universal Leave Application & Approval Portal (`LeaveApplyView.vue`)

#### Scope for Every User:
Every employee in every department has a statutory right and interface to request leave, track live approval statuses, and manage their personal duty schedule.

#### Key Features:
1. **Leave Application Form:**
   - Select Leave Type: Annual Leave (21 days), Sick Leave, Compassionate Leave, Study Leave, Maternity / Paternity Leave.
   - Date Picker & Duration Calculator (Automatically excludes weekends and Ugandan public holidays).
   - Handover Staff Officer Selection (Assigning temporary coverage to a colleague).
   - Upload Supporting Attachments (e.g. Medical Certificate, Examination Timetable).
2. **Real-Time Approval Status Tracker:**
   - Dynamic progress stepper bar:
     - `Stage 1: Submitted by User` $\rightarrow$ `Stage 2: HOD Recommendation` $\rightarrow$ `Stage 3: HR Verification` $\rightarrow$ `Stage 4: GS Accounting Officer Approval`.
   - Rejection Reason Alerts & Resubmission Workflow.
3. **Personal Leave Balance & Organizational Roster:**
   - Live KPI cards: Days Allocated, Days Taken, Days Remaining.
   - Departmental duty roster showing team members on leave to prevent coverage conflicts.

---

### 2.2 My Activities & Personal Audit Stream (`MyActivitiesView.vue`)

#### Scope for Every User:
A personal activity stream tracking every action performed by the user inside the intranet for personal accountability and productivity monitoring.

#### Key Features:
- Timestamped activity log: Forms submitted, work orders updated, asset entries edited, leave requests lodged, documents downloaded, and login events.
- Filter by date range, activity category, or status.
- Export personal activity summary as PDF for performance appraisals.

---

### 2.3 Direct Messaging & Inter-Office Memos (`Messages` / `internal_memos`)

#### Scope for Every User:
An internal real-time communication drawer facilitating peer-to-peer, team, and departmental collaboration without relying on external email.

#### Key Features:
- **Direct 1-on-1 Chat:** Instant messaging with any colleague across any department.
- **Departmental Group Channels:** General announcements and task discussions within Engineering, Finance, IT, HR, etc.
- **Inter-Office E-Memo Receiver & Responder:** Receive, minute comments on, and sign official internal memos.
- **File Attachment Support:** Send images, work order photos, PDFs, and spreadsheets directly in chat.

---

### 2.4 Real-Time Interactive Notifications Center (`🔔` / `NotificationsDrawer.vue`)

#### Scope for Every User:
A top-bar popover drawer providing instant real-time alerts and audio/visual cues for critical system events.

#### Key Features:
- **Notification Categories:**
  - *Leave Updates:* "Your Annual Leave request for Aug 01-14 has been Approved by GS."
  - *Approval Queues:* "New Form 5 Requisition NCS/WORKS/0045 requires your vote clearance."
  - *Task Assignments:* "Senior Engineer assigned Work Order WO-2026-089 to you."
  - *Security Alerts:* "New login detected from Chrome on Linux (192.168.1.45)."
- One-click actioning: Clicking a notification navigates directly to the relevant document or form.

---

### 2.5 My Profile Management (`ProfileView.vue`)

#### Scope for Every User:
Personal profile workspace where employees view and manage their official personnel identity.

#### Key Features:
- **Bio & Designation Information:** Staff ID, Designation, Department, Salary Scale, Duty Station, Date of Joining.
- **Contact Updates:** Phone number, personal email, residential address, emergency contact details.
- **Digital Credentials & Signature:** Upload official electronic signature for signing internal memos, work orders, and requisitions.
- **Profile Avatar Upload:** High-resolution staff photo upload.

---

### 2.6 Account & Security Settings (`SecuritySettingsView.vue`)

#### Scope for Every User:
Self-service security controls protecting employee accounts from unauthorized access.

#### Key Features:
- **Password Hygiene:** Change password with strength validation rules.
- **Two-Factor Authentication (2FA):** Enable TOTP Authenticator app (Google Authenticator / Authy) with backup recovery codes.
- **Active Login Sessions Manager:** View all active browser/mobile sessions with device details and IP addresses, plus a "Logout All Other Devices" button.
- **Notification Preferences:** Toggle email, SMS, and browser push notification settings.

---

## 3. Implementation Matrix Across All Departmental Dashboards

| Departmental Dashboard | Specialized Operational Tools | Universal Shared Workplace Modules Included |
| :--- | :--- | :---: |
| **General Secretary Dashboard** | Executive Approval Queue, Board Reports, CapEx Approvals | ✅ Leave, Activities, Messages, Notifications, Profile, Settings |
| **Accountant Dashboard** | General Ledger, Fixed Assets (31B), Federation Grants | ✅ Leave, Activities, Messages, Notifications, Profile, Settings |
| **HR Dashboard** | Staff Directory, 360 Profiles, Payroll, Appraisals | ✅ Leave, Activities, Messages, Notifications, Profile, Settings |
| **Internal Auditor Dashboard** | Mobile QR Spot-Checks, Discrepancies, Audit Logs | ✅ Leave, Activities, Messages, Notifications, Profile, Settings |
| **IT Officer Dashboard** | System Health, Automated DB Backups, Helpdesk Queue | ✅ Leave, Activities, Messages, Notifications, Profile, Settings |
| **Technical & Sports Admin** | NSMIS Registry, 50+ Federations, Athletes & Licensing | ✅ Leave, Activities, Messages, Notifications, Profile, Settings |
| **Procurement Unit (PDU)** | APP Tracker, Form 5 Pipeline, Bidding Evaluation | ✅ Leave, Activities, Messages, Notifications, Profile, Settings |
| **Public Relations (PR)** | Press Release Publisher, Media Accreditation, CMS | ✅ Leave, Activities, Messages, Notifications, Profile, Settings |
| **Engineering (All Roles)** | Work Orders, Facility Inspection Logs, Material Requisitions| ✅ Leave, Activities, Messages, Notifications, Profile, Settings |
| **Applicant / External Portal**| License Applications, Status Tracker, Certificate Vault | ✅ Profile, Security Settings, Messages, Notifications |

---

## 4. REST API Routing for Shared User Modules

| Endpoint | Method | Function |
| :--- | :--- | :--- |
| `/api/v1/user/leave/apply` | POST | Submit personal leave application |
| `/api/v1/user/leave/my-requests` | GET | View personal leave requests & real-time approval status |
| `/api/v1/user/activities` | GET | Query personal audit log and activity stream |
| `/api/v1/user/messages` | GET/POST | Send direct messages and retrieve conversation threads |
| `/api/v1/user/notifications` | GET/PUT | List notifications and mark as read |
| `/api/v1/user/profile` | GET/PUT | View and update bio details, avatar, and digital signature |
| `/api/v1/user/security/password` | PUT | Change user account password |
| `/api/v1/user/security/2fa` | POST | Enable/Disable Two-Factor Authentication |
