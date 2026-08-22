# Existing User & Applicant Dashboard Layout, Navigation & Feature Specification

## 1. Overview
This document provides an architectural specification of the **User & Applicant Dashboard** in `ncsportal`. It details the existing UI layout, sidebar navigation menu, section panels, credential wallet, application wizard, transaction tracking, and account security controls implemented in the frontend application (`UserPortalView.vue`, `ApplicationWizardView.vue`, `ApplicationDetailView.vue`, `OpenFormsPanel.vue`).

---

## 2. Global UI Shell & Top Navbar Controls

```
+-----------------------------------------------------------------------------------+
|  [≡ Toggle] [↻ Refresh]  NCS APPLICANT PORTAL            [🌙 Dark] [✉ 2] [🔔 3] [👤 User v] |
+-----------------------------------------------------------------------------------+
```

### Top Navbar Controls:
1. **Sidebar Collapse Button (`<i class="icofont-navigation-menu">`)**: Toggles collapsed/expanded sidebar mode (`sidebar-mini`).
2. **Refresh Action Button (`<i class="icofont-refresh">`)**: Asynchronously reloads active application forms, user stats, and notifications.
3. **Workspace Context Label**: Displays user profile classification (e.g. *Verified Athlete*, *Federation Applicant*, *General Public Account*) and section title.
4. **Theme Switcher (`<ThemeToggle>`)**: Switcher between Dark Mode and Light Mode.
5. **Messages Quick Action (`<i class="icofont-envelope">`)**: Badge indicator displaying unread reviewer message count; opens dropdown/drawer with reviewer feedback.
6. **Notifications Quick Action (`<i class="icofont-notification">`)**: Badge indicator displaying unread system notifications with a "Mark all read" trigger.
7. **User Profile Dropdown Menu**:
   -  **My Profile**: Identity details, photo avatar upload, email 2FA settings, and password change.
   -  **My Applications**: Standard & custom application history.
   -  **Transactions**: UGX payment records & proof uploads.
   -  **Logout**: Session termination.

---

## 3. Sidebar Navigation Menu & Section Structure

The sidebar menu for ordinary users/applicants is organized into distinct functional sections:

```
+-----------------------------------------------------------------------------------+
| [Sidebar Navigation Menu]                                                         |
|                                                                                   |
| ── PORTAL NAVIGATION ──────────────────────────────────────────────────────────── |
|  Dashboard Overview                                                              |
|  Apply Now (Open Applications)                                                  |
|  My Applications                                                                |
|  Credential Wallet (My Files & Certificates)                                   |
|  My Activity Log                                                                |
|  Notifications                                                                  |
|  Messages                                                                        |
|  My Transactions                                                                |
|                                                                                   |
| ── ACCOUNT & SECURITY ─────────────────────────────────────────────────────────── |
|  My Profile & 2FA Settings                                                      |
|  Logout                                                                         |
+-----------------------------------------------------------------------------------+
```

---

## 4. Detailed Feature Modules & Component Breakdown

### 4.1 Dashboard Overview (`section === 'dashboard'`)
- **Personalized Greeting & KPI Bar**:
  - *Applications Overview*: Total Applications, In Progress, Approved, Action Required.
  - *Quick Command*: "Apply Now" primary action button.
- **Open Applications Grid (`<OpenFormsPanel>`)**:
  - Displays published dynamic application forms available for submission (e.g. *Federation License Renewal Form*, *Athlete Registration Form*, *Event Sanction Request Form*).
- **Verified Athlete Registry Panel (Conditional for Athletes)**:
  - If the user account is linked to an official athlete record in NAMIS:
    - **Header Banner**: Athlete full name, NAMIS ID (`NAMIS-UG-xxxx`), sport discipline, and verified status badge.
    - **Classification Card**: Age Category (`Senior`, `U20`, etc.), District & Region, Affiliated Club, License Status (`ACTIVE`).
    - **National Duty Card**: National Squad Tier (`SENIOR`, `DEVELOPMENT`), Team Name (*Uganda Cranes*, *She Cranes*), Total International Caps, First Call-Up Date.
    - **Safeguarding & Medical Card**: Blood group, injury clearance status (`Fit`), safeguarding consent clearance (`Cleared`).
    - **WADA Compliance Card**: Testing pool status, last tested date, test result (`NEGATIVE`), WADA digital education status.
    - **Medal Standings Card**: Podium medals won by event (`Gold`, `Silver`, `Bronze`).
    - **Competition Results Card**: Recent sanctioned finishes, times/scores, and position rankings.

---

### 4.2 Apply Now / Open Applications (`section === 'apply'`)
- Expanded grid view of all active application forms published by NCS.
- Filters by department or service category.
- "Start Application" launcher opening the **Application Wizard**.

---

### 4.3 My Applications (`section === 'applications'`)
- **Search & Filter Bar**: Filter applications by search keyword and status (`Draft`, `Submitted`, `Under Review`, `Approved`, `Rejected`).
- **Applications Table**:
  - Columns: `Application Title`, `Source` (*Standard* vs *Custom Form*), `Reference Number` (e.g. `SUB-2026-xxxx`), `Status Badge`, `Payment Status`, `Last Updated Date`.
  - Actions:
    -  **View Details**: Opens full submission detail page (`ApplicationDetailView.vue`).
    -  **Download PDF**: Generates downloadable PDF application form summary.
    -  **Continue Draft / Edit**: Resumes draft application in the wizard.

---

### 4.4 Credential Wallet / My Files (`section === 'my-files'`)
- **Verified Digital Credentials Desk**:
  - Displays official licenses, registrations, and certificates issued to the user by NCS.
  - Total Credentials Counter and 100% Verification Badge.
  - Credential Card details: Certificate title, type (`ATHLETE_LICENSE`, `COACH_CERTIFICATE`, `OFFICIAL_CLEARANCE`), license/certificate serial number, issue date, status badge.
  - 📥 **Download Document Button**: Direct download of official PDF certificates.

---

### 4.5 Interactive Application Wizard (`ApplicationWizardView.vue`)
- **Navbar Header**: Brand logo, return to dashboard button, theme toggle, user avatar.
- **Application Overview Banner**: Form title, issuing department, description, application fee (UGX amount or *Free*).
- **Multi-Step Side Navigation**: Step-by-step progress tracker showing step number, title, and field counts.
- **Dynamic Field Rendering Engine**:
  - `short_text`: Single-line text input.
  - `long_text`: Multi-line textarea.
  - `dropdown`: Select dropdown options.
  - `radio`: Single-select choice list.
  - `checkbox`: Multi-select checkboxes & agreement terms.
  - `date`: Date picker.
  - `file_upload`: File dropzone supporting drag-and-drop document attachments (`.pdf`, `.png`, `.jpg` up to 25MB).
- **Step Navigation Buttons**: "Previous Step", "Save Draft", and "Submit Application".

---

### 4.6 Application Detail & Payment Review (`ApplicationDetailView.vue`)
- **Header Summary**: Application title, reference code, current status badge (`Submitted`, `Under Review`, `Approved`), submission date.
- **Submitted Answers Display**: Complete section-by-section breakdown of submitted field responses.
- **Document Attachments Vault**: Clickable file preview & download links for all uploaded attachments (`<FilePreview>`).
- **Payment & Transaction Desk**:
  - Payment instructions for bank transfer or mobile money.
  - UGX payment amount & status.
  - **Upload Payment Proof Form**: File upload for bank deposit slips or transaction receipts, payment reference code input, and submission button.

---

### 4.7 My Activity Log (`section === 'activities'`)
- Audit log timeline recording recent user security events, portal logins, application submissions, and profile updates with timestamps.

---

### 4.8 Notifications & Messages Feeds (`section === 'notifications'` / `section === 'messages'`)
- **Notifications Feed**: Feed of application state updates, certificate issuance alerts, and system announcements with a "Mark as read" trigger.
- **Messages Feed**: Direct reviewer communication thread displaying official comments and request-for-information notes sent by NCS reviewers.

---

### 4.9 My Transactions (`section === 'transactions'`)
- Financial ledger displaying recorded application transactions.
- Transaction Summary: Total UGX recorded.
- Transactions Table: `Application`, `Payment Reference`, `Amount (UGX)`, `Payment Status` (*Paid*, *Unpaid*, *Pending Verification*), `Payment Method`.

---

### 4.10 My Profile & Security Settings (`section === 'profile'`)
- **Profile Photo Avatar**:
  - Avatar image display with fallback initials badge.
  - Upload Photo button with file picker (`avatar_url`).
- **Personal Details Form**:
  - First Name, Last Name, Email address (read-only).
  - Save Profile action button.
- **Password Change Form**:
  - Current Password, New Password (min 8 characters), Confirm New Password.
  - Update Password action.
- **Email Two-Factor Authentication (2FA)**:
  - Toggle switch to enable email-based 2FA verification.
  - 6-digit email validation code verification panel.

---

## 5. Summary Table of Frontend Views & Components for Users

| User View / Component | Primary Responsibilities | Route / Trigger |
| :--- | :--- | :--- |
| `UserPortalView.vue` | Main User Portal Shell, Navigation, Dashboard Panels, Wallet & Profile | `/dashboard` |
| `ApplicationWizardView.vue` | Multi-Step Application Form Wizard & Dynamic Field Form Renderer | `/dashboard/apply/:slug` |
| `ApplicationDetailView.vue` | Submission Detail Inspector, Answers View, File Previewer, Payment Proof Upload | `/dashboard/applications/:id` |
| `OpenFormsPanel.vue` | Published Open Application Cards Grid & Form Launcher | Dashboard Overview & Apply Now |
| `FilePreview.vue` | Document Attachment Previewer & Modal PDF Viewer | Application Detail & Review |
