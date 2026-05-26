# Applications Domain — NCSMS v1.0

## Overview

This folder contains the individual application form definitions for all user-facing registration, renewal, licence, and transformation requests submitted to the National Council of Sports (NCS).

Each form definition documents the exact field names, input types, labels, wizard step mapping, conditional logic, required file attachments, and workflow — derived directly from the official NCS statutory forms.

All applications share a common **6-stage wizard workflow** described in [application-workflow-v1.md](application-workflow-v1.md).

---

## Application Submission Workflow (All Forms)

```
STAGE 1          STAGE 2         STAGE 3                STAGE 4              STAGE 5       STAGE 6
FILL FORM   →   REVIEW    →   DOWNLOAD & SIGN   →   PAYMENT       →   SUBMIT    →   TRACK
(Wizard)        (Summary)      (PDF + Upload)        (Digital or              (Confirm)     (Dashboard)
                                                      Proof Upload)
```

- **Stage 1** — Multi-step wizard; auto-save on every field; "Save Draft" at every step; resume from last active step
- **Stage 2** — Read-only review of all entered data; edit links per section
- **Stage 3** — System generates a pre-filled PDF; applicant signs digitally in-browser or prints, signs, and uploads scanned PDF
- **Stage 4** — Payment via **digital payment** (card / mobile money) OR **upload scanned proof of payment** (bank slip, receipt) — both methods accepted for all application types including renewals
- **Stage 5** — Final submission with checklist confirmation; `application_reference` assigned
- **Stage 6** — Status tracking dashboard with reviewer notes, NEEDS_INFORMATION response flow, and notifications

See [application-workflow-v1.md](application-workflow-v1.md) for the full architecture including status state machine, API endpoints, and frontend component map.

---

## Included Application Forms

| File | Official Form | Regulation | Description | Wizard Steps |
|---|---|---|---|---|
| [national-sport-recognition-application-v1.md](national-sport-recognition-application-v1.md) | Form 1 | Reg. 3(1) | Application for Declaration of National Sport | 6 |
| [federation-registration-renewal-v1.md](federation-registration-renewal-v1.md) | Form 3 | Reg. 4(1), 5(1), 9(1) | Application for Registration or Renewal (NSA / NSF) | 6 |
| [national-sports-association-transformation-v1.md](national-sports-association-transformation-v1.md) | Form 5 | Reg. 11(1) | NSA applicant perspective — Transformation from NSA to NSF | 6 |
| [national-sports-federation-transformation-v1.md](national-sports-federation-transformation-v1.md) | Form 5 | Reg. 11(1) | Full form definition — Transformation from NSA to NSF | 6 |
| [sports-competition-organization-application-v1.md](sports-competition-organization-application-v1.md) | Form 7 | Reg. 16(1) | Application to Organise a Sports Competition | 4 |
| [sports-facility-operation-application-v1.md](sports-facility-operation-application-v1.md) | Form 8 | Reg. 17(2) | Application to Operate a Sports Facility | 4 |
| [community-academy-registration-renewal-v1.md](community-academy-registration-renewal-v1.md) | Form 10 | Reg. 22(1) | Application for Registration / Renewal of Community Sports Club | 4 |
| [sports-academy-operation-application-v1.md](sports-academy-operation-application-v1.md) | Form 11 | Reg. 29(1)(a) | Application to Operate a Sports Academy | 5 |

---

## Payment Methods Accepted

| Method | Description | Accepted For |
|---|---|---|
| Digital Payment | Online card payment or mobile money (MTN MoMo, Airtel Money) | All applications and renewals |
| Proof of Payment Upload | Scanned bank slip, mobile money receipt, or payment confirmation document | All applications and renewals |

---

## Common Tracking Metadata (All Forms)

| Field | Description |
|---|---|
| `application_reference` | Unique ID assigned on submission (e.g., `NCS-2026-F3-00142`) |
| `payment_reference` | Payment tracking number |
| `payment_status` | `UNPAID` \| `PAYMENT_INITIATED` \| `PAID` \| `PROOF_UPLOADED` \| `PAYMENT_VERIFIED` \| `PAYMENT_REJECTED` |
| `application_status` | `DRAFT` \| `PENDING_SIGNATURE` \| `PENDING_PAYMENT` \| `SUBMITTED` \| `UNDER_REVIEW` \| `NEEDS_INFORMATION` \| `RESUBMITTED` \| `APPROVED` \| `REJECTED` |
| `reviewer_id` | The `GENERAL_SECRETARY` or delegated reviewer |
| `submitted_at`, `last_updated`, `next_action` | Lifecycle timestamps and applicant prompts |

---

## Required Attachments Reference

| Attachment | Forms |
|---|---|
| Proof of payment of prescribed fees | Forms 3, 5, 8, 10, 11 |
| Organogram of national sports organisation | Forms 3, 5 |
| Certificate of recognition as National Sport | Forms 3, 5 |
| Constitution approved by general assembly | Forms 3, 5 |
| General assembly minutes (constitution approval) | Forms 3, 5 |
| Certified and updated list of members | Forms 3, 5 |
| Sports activities report (within 1 year prior) | Forms 3, 5 |
| Districts presence list | Forms 3, 5 |
| Election minutes for executive committee | Forms 3, 5 |
| Audited books of accounts (where applicable) | Forms 3, 5 |
| Passport photos and ID documents of exec. board | Forms 3, 5 |
| Symbols, slogans, and colours | Forms 3, 5 |
| Recommendation from relevant NSA / NSF | Forms 8, 11 |
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

- The user dashboard exposes these forms as menu items in the **Applications** section
- Applicants can start a new application, save drafts, resume incomplete forms, view submitted applications, inspect payment status, and respond to NEEDS_INFORMATION requests
- Status badges, reviewer notes, and next-action prompts are surfaced in real time
- Conditional attachment uploads are shown or hidden based on field selections in the wizard

---

## API Endpoint Naming

- `POST /api/v1/applications/{form-type}/draft` — save or update a draft
- `GET /api/v1/applications/{form-type}/draft` — load latest draft for the signed-in user
- `POST /api/v1/applications/{id}/prefilled-pdf` — generate pre-filled PDF
- `POST /api/v1/applications/{id}/signed-form` — upload signed form PDF
- `POST /api/v1/applications/{id}/payment/initiate` — start online payment
- `POST /api/v1/applications/{id}/payment/proof` — upload proof of payment
- `POST /api/v1/applications/{id}/submit` — final submission
- `GET /api/v1/applications/{id}` — fetch application and tracking metadata
- `GET /api/v1/applications` — list all applications for the signed-in user
- `PATCH /api/v1/applications/{id}/respond` — respond to NEEDS_INFORMATION
- `POST /api/v1/applications/{id}/attachments` — upload a required attachment
