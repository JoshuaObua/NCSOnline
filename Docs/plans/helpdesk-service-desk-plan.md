# NCS Enterprise IT Helpdesk & Service Desk Master Blueprint

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/helpdesk-service-desk-plan.md`  
> **Role:** All NCS Employees (Requesters) / IT Support Technicians / Helpdesk Manager  
> **System Scope:** Organization-Wide IT Support & Universal Self-Service Portal (`HelpdeskDashboardView.vue`)  
> **Authority Scope:** Technical Issue Ticketing, Visitor Gate Clearance Applications, Leave Applications, Hardware/Software Support Requests, Network SLA Tracking, Knowledge Base Solutions  
> **Universal Modules Scope:** Visitor Clearance Form, Leave Application & Approval Tracker, My Activities Stream, My Profile, Account & Security Settings, Direct Messages & Notifications  
> **Compliance Standards:** ITIL (Information Technology Infrastructure Library) Service Operation Standards, ISO 27001  

---

## 1. Executive Purpose & Multi-Service Portal Architecture

The **Enterprise IT Helpdesk & Service Desk Module** serves as the primary technical support and staff self-service portal for all NCS employees across Engineering, Finance, HR, Technical, Procurement, Internal Audit, and Executive departments.

It unifies IT support ticketing with essential workplace self-service functions: **Visitor Gate Clearance Applications**, **Personal Leave Applications**, **My Activities Log**, **Profile Management**, and **Account Security Settings**.

```
+---------------------------------------------------------------------------------------------------+
|                        ENTERPRISE HELPDESK & SHARED SERVICE ENGINE                                |
+---------------------------------------------------------------------------------------------------+
       |                    |                     |                     |                    |
       v                    v                     v                     v                    v
+--------------+    +---------------+     +---------------+     +---------------+    +---------------+
| 1. IT Support|    | 2. Visitor    |     | 3. Personal   |     | 4. My         |    | 5. Profile &  |
| Ticketing    |    | Clearance Form|     | Leave Portal  |     | Activities Log|    | Security      |
| & SLA Tracking|   | & Gate Pass   |     | & Approvals   |     | & Stream      |    | Settings      |
+--------------+    +---------------+     +---------------+     +---------------+    +---------------+
       |                    |                     |                     |                    |
       +--------------------+---------------------+---------------------+--------------------+
                                                  |
                                                  v
                     +-----------------------------------------------------------+
                     | CENTRAL LOGGING, NOTIFICATION & IT SUPPORT DASHBOARD      |
                     +-----------------------------------------------------------+
```

---

## 2. Key Modules & User Interface Specifications (`HelpdeskDashboardView.vue`)

```
+-----------------------------------------------------------------------------------+
| [≡] NCS INTRANET | Helpdesk & Employee Service Center    [➕ New Ticket] [🎫 Visitor Pass] |
+-----------------------------------------------------------------------------------+
| [HELPDESK & SERVICE KPI CARDS]                                                    |
| +--------------------+ +--------------------+ +-------------------+ +-------------------+ |
| | Open IT Tickets    | | Avg SLA Resolution | | Active Visitor Pass| | Leave Days Available|
| | 3 Tickets Open     | | 1 hr 45 mins       | | 2 Approved Today  | | 14 Days Remaining | |
| +--------------------+ +--------------------+ +-------------------+ +-------------------+ |
+-----------------------------------------------------------------------------------+
| [HELPDESK WORKSPACE TABS]                                                         |
| (1) IT Support Tickets | (2) Visitor Clearance Form | (3) My Leave Applications     |
| (4) My Activities Log  | (5) My Profile & Credentials | (6) Account Security      |
| +-------------------------------------------------------------------------------+ |
| | [ACTIVE IT TICKETS & SERVICE REQUESTS QUEUE]                                  | |
| | Ticket ID  | Category       | Subject                 | Priority | Status     | |
| | T-2026-042 | Printer Drivers| Kyocera Driver Error    | Medium   | In Progress| |
| | T-2026-045 | Access / 2FA   | Account Lockout Reset   | High     | Resolved   | |
| | T-2026-048 | Network Drop   | Switch Port #14 Down    | Urgent   | Assigned   | |
| | [ Submit New Ticket ] [ View Solution KB ] [ Re-Open Ticket ] [ Rate Service ] | |
| +-------------------------------------------------------------------------------+ |
| +-----------------------------------------------+ +-------------------------------+ |
| | Visitor Clearance Quick Action                | | Personal Leave Status Summary | |
| | Active Pass: VIS-2026-089 (Lugogo Head Office)| | Annual Leave: Aug 01 - Aug 14 | |
| | Visitor    : Kampala Sports Club Delegation   | | Status      : Approved by GS  | |
| | [ Request Visitor Gate Clearance Form ]       | | [ Submit New Leave Request ]  | |
| +-----------------------------------------------+ +-------------------------------+ |
+-----------------------------------------------------------------------------------+
```

---

## 3. Detailed Feature Breakdown & Specifications

---

### 3.1 IT Support Ticket Lifecycle & SLA Management

1. **Ticket Submission & Classification:**
   - Categories: Hardware Repair, Software Installation, Network & Wi-Fi Dropouts, Email/Password Access, Printer Drivers, System Bugs.
   - Priority Levels: Low (24 hrs), Medium (8 hrs), High (4 hrs), Urgent (2 hrs).
   - Screenshot & Document Attachments: Users attach log files or screenshots of technical errors.
2. **Technician Triage & Automated Routing:**
   - Auto-assign tickets based on category to specific IT Technicians.
   - Live countdown SLA timer alerting Helpdesk Manager when a ticket approaches threshold.
3. **Service Knowledge Base (FAQ):**
   - Self-service articles allowing staff to resolve minor issues (Wi-Fi password reset, printer setup).

---

### 3.2 Visitor Clearance Application Form & Gate Pass Request (`VisitorClearanceApplyView.vue`)

1. **Visitor Clearance Form:**
   - Staff members can apply for visitor gate passes on behalf of official guests, contractors, media teams, or delegation members visiting NCS premises (Lugogo Head Office, Stadium, Hostels).
   - Fields: Visitor Full Name, NIN / Passport, Phone, Organization, Purpose of Visit, Destination Venue, Date & Time, Vehicle Registration Plate Number.
2. **Approval & Digital Pass Issuance:**
   - Route to HOD / Admin for clearance approval.
   - Generates an official digital QR Code Gate Pass sent via SMS/Email to the visitor and printable from the helpdesk dashboard.

---

### 3.3 Universal Personal Leave Application & Approval Status Portal (`LeaveApplyView.vue`)

1. **Leave Application Form:**
   - Select Leave Type: Annual Leave, Sick Leave, Compassionate, Study, Maternity / Paternity Leave.
   - Start and End Date selection (Automatic calculation excluding weekends and public holidays).
   - Assign Handover Staff Officer.
   - Upload supporting medical or academic documents.
2. **Live Approval Stepper:**
   - Track progress: `Submitted` $\rightarrow$ `HOD Recommendation` $\rightarrow$ `HR Verification` $\rightarrow$ `GS Approval`.
3. **Leave Balance Summary:**
   - Live display of Allocated Days, Taken Days, and Remaining Balance.

---

### 3.4 My Activities Log & Audit Stream (`MyActivitiesView.vue`)

1. **Personal Activity Stream:**
   - Timestamped record tracking every ticket submitted, visitor clearance requested, leave application lodged, profile update, and login event.
2. **Search & Export:**
   - Filter activities by date range or category; export summary for performance appraisal reviews.

---

### 3.5 My Profile Management (`ProfileView.vue`)

1. **Personal Bio Data & Position Details:**
   - View/edit Staff ID, Designation, Department, Duty Station, Emergency Contacts, and Phone Number.
2. **Digital Signature & Photo Upload:**
   - Upload official digital signature for signing IT requisitions and leave forms.
   - Profile picture avatar management.

---

### 3.6 Account & Security Settings (`SecuritySettingsView.vue`)

1. **Password Management:**
   - Self-service password change with strength validation.
2. **Two-Factor Authentication (2FA):**
   - Enable/disable TOTP authenticator app (Google Authenticator) with emergency backup codes.
3. **Active Session Governance:**
   - View logged-in devices, IP addresses, and browsers; trigger "Remote Logout All Other Sessions".

---

## 4. Database Schema & REST API Mapping

```sql
-- Helpdesk Tickets Table
CREATE TABLE helpdesk_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_code VARCHAR(32) UNIQUE NOT NULL,
    requester_id UUID NOT NULL REFERENCES users(id),
    category VARCHAR(64) NOT NULL,
    priority VARCHAR(16) DEFAULT 'MEDIUM', -- LOW, MEDIUM, HIGH, URGENT
    subject VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    attachment_url TEXT,
    assigned_to UUID REFERENCES users(id),
    status VARCHAR(32) DEFAULT 'OPEN', -- OPEN, IN_PROGRESS, RESOLVED, CLOSED
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

| Endpoint | Method | Scope | Function |
| :--- | :--- | :--- | :--- |
| `/api/v1/helpdesk/tickets` | GET/POST | All Staff | Create and list IT support tickets |
| `/api/v1/helpdesk/tickets/:id` | PUT | IT Technicians | Update ticket status, assign tech, resolve |
| `/api/v1/visitors/apply` | POST | All Staff, Visitors | Submit visitor gate clearance application form |
| `/api/v1/user/leave/apply` | POST | All Staff | Submit personal leave application form |
| `/api/v1/user/activities` | GET | All Staff | Query personal activities audit log |
| `/api/v1/user/profile` | GET/PUT | All Staff | View and update bio details, avatar, and digital signature |
| `/api/v1/user/security/2fa` | POST | All Staff | Enable TOTP 2FA authentication |

---

## 5. Implementation Verification Roadmap

- [ ] Update `HelpdeskDashboardView.vue` to include tabs for IT Support Tickets, Visitor Clearance Form, My Leave Applications, My Activities, Profile, and Security Settings.
- [ ] Connect Go REST API handlers in `backend/internal/handlers/` for all helpdesk self-service endpoints.
- [ ] Test end-to-end user workflows from ticket logging to visitor clearance pass generation and leave tracking.
