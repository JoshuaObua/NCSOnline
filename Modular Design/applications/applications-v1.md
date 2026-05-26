# Applications Domain — NCSMS v1.0

## Overview

This folder contains the individual application form definitions for all user-facing registration, renewal, licence, and transformation requests submitted to the National Council of Sports (NCS).

Each form definition documents the exact field names, input types, labels, conditional logic, required file attachments, and workflow — derived directly from the official NCS statutory forms.

Forms are submitted by signed-in users and tracked through the `GENERAL_SECRETARY` review workflow with payment monitoring and status notifications.

---

## Included Application Forms

| File | Official Form | Regulation | Description |
|---|---|---|---|
| `national-sport-recognition-application-v1.md` | Form 1 | Reg. 3(1) | Application for Declaration of National Sport |
| `federation-registration-renewal-v1.md` | Form 3 | Reg. 4(1), 5(1), 9(1) | Application for Registration or Renewal (NSA / NSF) |
| `national-sports-association-transformation-v1.md` | Form 5 | Reg. 11(1) | NSA applicant perspective — Transformation from NSA to NSF |
| `national-sports-federation-transformation-v1.md` | Form 5 | Reg. 11(1) | Full form definition — Transformation from NSA to NSF |
| `sports-competition-organization-application-v1.md` | Form 7 | Reg. 16(1) | Application to Organise a Sports Competition |
| `sports-facility-operation-application-v1.md` | Form 8 | Reg. 17(2) | Application to Operate a Sports Facility |
| `community-academy-registration-renewal-v1.md` | Form 10 | Reg. 22(1) | Application for Registration / Renewal of Community Sports Club |
| `sports-academy-operation-application-v1.md` | Form 11 | Reg. 29(1)(a) | Application to Operate a Sports Academy |

---

## Common Tracking Concepts

- `application_reference` — unique ID assigned to each submission
- `payment_reference` — tracking number for fee payments and receipts
- `status` — workflow state: `DRAFT`, `SUBMITTED`, `UNDER_REVIEW`, `APPROVED`, `REJECTED`, `NEEDS_INFORMATION`
- `reviewer_id` — the `GENERAL_SECRETARY` or delegated reviewer processing the application
- `submitted_at`, `last_updated`, `next_action` — lifecycle timestamps

---

## Attachment Types Used Across Forms

| Attachment | Forms That Require It |
|---|---|
| Proof of payment of prescribed fees | Forms 3, 5, 8, 10, 11 |
| Organogram of national sports organisation | Forms 3, 5 |
| Certificate of recognition as National Sport | Forms 3, 5 |
| Constitution approved by general assembly | Forms 3, 5 |
| General assembly minutes (constitution approval) | Forms 3, 5 |
| Certified and updated list of members | Forms 3, 5 |
| Sports activities report (within 1 year) | Forms 3, 5 |
| Districts presence list | Forms 3, 5 |
| Election minutes for executive committee | Forms 3, 5 |
| Audited books of accounts (where applicable) | Forms 3, 5 |
| Passport photos and ID documents of exec. board | Forms 3, 5 |
| Symbols, slogans, and colours of association/federation | Forms 3, 5 |
| Recommendation from NSA / NSF | Forms 8, 11 |
| Occupation and safety permits | Form 8 |
| Certificate of incorporation | Form 11 (conditional) |
| Technical and financial proposal | Form 11 (conditional) |
| Work plan to develop the sport | Form 1 |
| Executive committee meeting minutes | Form 7 |
| Approved action plan for financial year | Form 7 |
| Competition fixture | Form 7 |
| List of sponsors | Form 7 |

---

## Dashboard Integration

- The user dashboard exposes these forms as menu items in the `Applications` section
- Applicants can save drafts, resume incomplete forms, view submitted applications, and inspect payment history
- Status updates and reviewer notes are surfaced to the applicant through the same dashboard experience
- Conditional attachments are shown or hidden based on radio/select field responses

---

## Endpoint Naming

- `POST /api/v1/applications/<form-type>` — submit a new application
- `GET /api/v1/applications` — list applications for the signed-in user
- `GET /api/v1/applications/{id}` — fetch application details and tracking metadata
- `PATCH /api/v1/applications/{id}` — update a draft or resubmit with new data
- `GET /api/v1/transactions` — view payment and fee transaction history
- `GET /api/v1/transactions/{reference}` — lookup a specific payment reference
