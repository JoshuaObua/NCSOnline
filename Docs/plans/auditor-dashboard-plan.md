# NCS Internal Auditor & Audit Department Dashboard Implementation Plan

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/auditor-dashboard-plan.md`  
> **Role:** Internal Auditor / Senior Internal Auditor / Audit Manager  
> **Department:** Internal Audit Department - National Council of Sports (NCS)  
> **Authority Scope:** Independent Audit Verification, Physical Tag Scanning, Financial & Asset Discrepancy Flagging, Compliance Audit, Immutable System Audit Logs  
> **Universal Scope:** My Activities, Direct Messages, Real-Time Notifications, Profile, Security Settings, and Personal Leave Application & Approval Portal  
> **Compliance Standards:** Public Finance Management Act (PFMA 2015), Treasury Instructions (2017), PPDA Act, Auditor General Uganda Standards  

---

## 1. Executive Role Overview & Mission

The **Internal Audit Department** acts as an independent oversight organ ensuring internal financial controls, asset safeguards, statutory compliance, and operational integrity across all NCS departments.

This plan details the **Internal Audit Command Portal** inside `NCS_Intranet`. Auditors possess non-tamperable read rights across all ledgers, physical QR tag scanning controls, discrepancy flagging authority, automated statutory audit engines, and full personal workplace capabilities.

---

## 2. Key Modules & User Interface Specifications

```
+-----------------------------------------------------------------------------------+
| [≡] NCS INTRANET | Internal Audit Command & Compliance Portal    [🔍 Spot-Check] [📋 Audit Report] |
+-----------------------------------------------------------------------------------+
| [AUDIT COMPLIANCE & RISK METRICS]                                                |
| +--------------------+ +--------------------+ +-------------------+ +-------------------+ |
| | Audit Integrity    | | Physical Tag Count | | Open Discrepancies| | Risk Rating       | |
| | 96.8% Compliant    | | 294 / 297 Tagged  | | 2 Flagged Issues  | | LOW (Clean Audit) | |
| +--------------------+ +--------------------+ +-------------------+ +-------------------+ |
+-----------------------------------------------------------------------------------+
| [AUDITOR WORKSPACE: (1) Physical Spot-Check | (2) Discrepancy Manager | (3) System Audit] |
| +-------------------------------------------------------------------------------+ |
| | [PHYSICAL ASSET & STOCK SPOT-CHECK SCANNER]                                   | |
| | Scan Asset Tag: [ Input Barcode / QR Code: 166BLNG10                 ] [ Scan ] | |
| | System Match  : M1007058 - NCS BLOCK - OFFICE FLOOR-166-BLNG-10               | |
| | System Value  : UGX 298,000,000.00 | Location: NCS Lugogo Block Floor 2       | |
| | Audit Findings: [ Confirm Verified ] [ Flag Cost Mismatch ] [ Flag Missing ]    | |
| +-------------------------------------------------------------------------------+ |
| +-----------------------------------------------+ +-------------------------------+ |
| | Discrepancy & Exception Queue                 | | Universal Personal Suite       | |
| | - [FLAG] M1007058 FB_COST vs Adj Cost Mismatch| | - [LEAVE] Sick Leave Request   | |
| | - [FLAG] Travel Requisition > Vote Allocation | | - [ACTIVITIES] 28 Audits Scanned| |
| | [ Query Accountant ] [ Lock Asset Record ]    | | - [MESSAGES] 1 Auditor Memo   | |
| +-----------------------------------------------+ +-------------------------------+ |
+-----------------------------------------------------------------------------------+
```

---

## 3. Universal Shared Personal Workplace Suite (For All Audit Staff)

In addition to specialized audit and spot-checking tools, every user in the Internal Audit Department receives:

1. **Personal Leave Application & Status Portal (`LeaveApplyView.vue`):**
   - Submit leave applications (Annual, Sick, Compassionate, Study).
   - Real-time approval tracking: `Submitted` $\rightarrow$ `Audit Manager Recommendation` $\rightarrow$ `HR Verification` $\rightarrow$ `GS Approval`.
   - Personal leave balance summary and Audit Department duty roster.
2. **My Activities Audit Stream (`MyActivitiesView.vue`):**
   - Personal log tracking asset QR scans performed, audit discrepancies raised, and compliance reports compiled.
3. **Direct Messages & Inter-Office Memos:**
   - Internal chat drawer and E-Memo system for communicating with HODs, Finance, and Executive offices.
4. **Notifications Center (`🔔`):**
   - Real-time alerts for discrepancy responses, leave approval decisions, and high-risk system exceptions.
5. **My Profile Management (`ProfileView.vue`):**
   - Bio data, staff ID, digital signature upload for audit paper sign-offs, avatar upload.
6. **Account & Security Settings (`SecuritySettingsView.vue`):**
   - Password management, 2FA TOTP setup, and active login session management.

---

## 4. Database Schema & Technical Architecture

```sql
CREATE TABLE audit_discrepancies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discrepancy_code VARCHAR(64) UNIQUE NOT NULL,
    entity_type VARCHAR(64) NOT NULL, -- ASSET, FINANCIAL_VOUCHER, VOTE_HEAD, INVENTORY
    entity_id UUID NOT NULL,
    severity VARCHAR(32) NOT NULL DEFAULT 'MEDIUM', -- LOW, MEDIUM, HIGH, CRITICAL
    description TEXT NOT NULL,
    raised_by UUID NOT NULL REFERENCES users(id),
    assigned_to UUID REFERENCES users(id),
    status VARCHAR(32) NOT NULL DEFAULT 'OPEN', -- OPEN, UNDER_REVIEW, RESOLVED, ESCALATED_TO_GS
    resolution_notes TEXT,
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## 5. REST API Mapping (`backend/internal/handlers/auditor_handler.go`)

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/api/v1/audit/spot-check` | POST | Verify physical asset tag via QR barcode scan |
| `/api/v1/audit/discrepancies` | GET/POST | Query and flag audit discrepancies |
| `/api/v1/user/leave/apply` | POST | Auditor personal leave application submittal |
| `/api/v1/user/activities` | GET | Personal audit activity log query |

---

## 6. Implementation Verification Roadmap

- [ ] Create Go backend handler `backend/internal/handlers/auditor_handler.go`.
- [ ] Build Vue frontend view `frontend/src/views/audit/AuditorDashboard.vue`.
- [ ] Verify personal leave application and universal workplace modules for all audit staff.
