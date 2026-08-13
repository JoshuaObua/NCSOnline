# NCS Visitor & Facility Clearance Application Portal Implementation Plan

> **Target File:** `/home/fidi/Projects/NCS_Intranet/Docs/plans/visitor-clearance-portal-plan.md`  
> **Role:** Public Visitors / Event Organizers / Facility Managers / Security Officers  
> **System Scope:** External Visitor & Venue Gate Clearance Portal (`VisitorClearanceApplyView.vue`)  
> **Authority Scope:** Visitor Pass Applications, Stadium Venue Access Permits, Lugogo Sports Ground Entry Clearance, Gate Security Pass Scanner, Vehicle Entry Clearance  
> **Compliance Standards:** National Sports Complex Access Guidelines, Public Event Safety Standards  

---

## 1. Executive Purpose & Architecture

The **Visitor & Facility Clearance Application Portal** provides a digital pre-clearance mechanism for visitors, official delegations, VIP guests, contractors, media teams, and event organizers accessing National Council of Sports premises and venue facilities (Lugogo Indoor Stadium, Lugogo Tennis Complex, Cricket Oval, Hostels, Head Office).

```
+-----------------------------------------------------------------------------------+
|                        VISITOR & STADIUM FACILITY CLEARANCE ENGINE                |
+-----------------------------------------------------------------------------------+
       |                                  |                                 |
       v                                  v                                 v
+-----------------------+      +-----------------------+      +---------------------+
| 1. VISITOR PRE-CLEAR  |      | 2. APPROVAL & ACCESS  |      | 3. GATE SECURITY QR |
| (`VisitorClearance`)  |      |    PERMIT ISSUANCE    |      |    PASS SCANNER     |
+-----------------------+      +-----------------------+      +---------------------+
| - Identity Verification|     | - Admin / Facility    |      | - Barcode Gate Scan |
| - Visit Purpose       |      |   Officer Clearance   |      | - Check-In Log      |
| - Vehicle Plate Reg   |      | - Digital QR Pass     |      | - Check-Out Log     |
+-----------------------+      +-----------------------+      +---------------------+
```

---

## 2. Key Modules & User Interface Specifications

### 2.1 Public Visitor Clearance Application (`VisitorClearanceApplyView.vue`)
- **Pre-Registration Wizard:** Captures visitor full name, NIN / Passport, phone number, organization, person/officer to visit, venue requested, date/time, and vehicle registration number.
- **Instant Digital QR Pass:** Once approved, visitors receive an SMS/email containing a digital QR Code Access Pass.

### 2.2 Security Gate Check-In & Scanner Interface
- **Mobile/Tablet Gate Scanner:** Security officers at Lugogo gates scan the visitor's QR code to verify identity, timestamp entry, and log exit.

---

## 3. Database Schema & REST API Mapping

```sql
CREATE TABLE visitor_clearances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clearance_code VARCHAR(32) UNIQUE NOT NULL,
    visitor_name VARCHAR(128) NOT NULL,
    nin_or_passport VARCHAR(32) NOT NULL,
    phone_number VARCHAR(32) NOT NULL,
    host_officer_id UUID REFERENCES users(id),
    destination_facility VARCHAR(128) NOT NULL, -- Head Office, Indoor Stadium, Hostel
    purpose_of_visit TEXT NOT NULL,
    vehicle_registration VARCHAR(32),
    visit_date DATE NOT NULL,
    expected_time_in TIME NOT NULL,
    status VARCHAR(32) DEFAULT 'PENDING', -- PENDING, APPROVED, CHECKED_IN, CHECKED_OUT, DENIED
    qr_code_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

| Endpoint | Method | Scope | Function |
| :--- | :--- | :--- | :--- |
| `/api/v1/visitors/apply` | POST | Public | Submit visitor facility clearance request |
| `/api/v1/visitors/scan` | POST | Security Officers | Scan visitor QR pass at venue gates |
