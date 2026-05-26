# Sports Facility Operation Application — NCSMS v1.0

**Form Reference:** Form 8 — Regulation 17(2)
**Submitted To:** General Secretary, National Council of Sports

## Purpose

This application is submitted by an individual or entity seeking approval to operate a sports facility in Uganda.

---

## Form Fields

### Introduction / Facility Identification

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `facility_name` | text | Name of Sports Facility | Insert exact name of the sports facility |
| `sports_activity_nature` | text | Nature of Sports Activity to be Hosted | Describe the type of sports activity the facility will host |

---

### Applicant Details

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `applicant_name` | text | Name of Applicant | Full name of applicant (individual or organisation) |
| `physical_address` | textarea | Physical Address | |
| `postal_address` | text | Postal Address | |
| `telephone_fixed_line` | text | Telephone (Fixed Line) | |
| `mobile_phone` | text | Mobile Phone | |
| `email_address` | email | E-mail Address | |
| `website` | url | Website | Optional |

---

### Facility Details

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `facility_location_and_size` | textarea | Location and Size of Sports Facility | Include physical address, plot size, and capacity |
| `sport_activity_types` | textarea | Types of Sport Activities | List all sport activities to be conducted at the facility |
| `funding_sources` | textarea | Sources of Funding | State all sources of operational funding |

---

### Declaration and Signature

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `authorised_signatory_name` | text | Authorised Signatory Name | |
| `signature_upload` | file | Signature / Seal | Upload scanned signature and official seal |
| `date_signed` | date | Date Signed | |

---

## Required Attachments

| # | Field Name | Label | Required |
|---|---|---|---|
| 1 | `proof_of_fee_payment` | Proof of Payment of Prescribed Fees | Mandatory |
| 2 | `federation_recommendation` | Recommendation from the Relevant National Sports Federation or National Sports Association | Mandatory |
| 3 | `occupation_safety_permits` | Occupation and Safety Permits | Mandatory |

---

## Wizard Steps

This form follows the shared **6-stage application workflow** defined in [application-workflow-v1.md](application-workflow-v1.md). The form content maps to the following wizard steps:

| Step | Label | Fields Covered |
|---|---|---|
| 1 | Facility Identification | `facility_name` — name of the sports facility; `sports_activity_nature` — type of sports activity to be hosted |
| 2 | Applicant Details | Name, physical address, postal address, telephone, mobile, email, website |
| 3 | Facility Details & Funding | Location and size, types of sport activities, sources of funding |
| 4 | Review & Confirm | Read-only summary; edit links per section |

After Step 4 the applicant proceeds through **Stage 3 (Download, Sign & Upload)** → **Stage 4 (Payment)** → **Stage 5 (Submit)** → **Stage 6 (Track)**.

### Payment Options
- **Digital payment** — online card or mobile money
- **Upload proof of payment** — scanned bank slip or mobile money receipt

Both methods accepted for both new facility licences and renewals.

---

## Tracking Metadata

- `application_reference` — unique submission identifier assigned on submission
- `facility_license_reference` — reference number for the facility operating licence
- `payment_reference` — application fee payment tracking number
- `status` — `DRAFT` \| `SUBMITTED` \| `UNDER_REVIEW` \| `APPROVED` \| `REJECTED` \| `NEEDS_INFORMATION`
- `submitted_at`, `last_updated`, `next_action`

---

## Workflow

1. Applicant completes all form fields and uploads 3 required attachments
2. Submits to `GENERAL_SECRETARY` for operational review
3. Reviewer verifies facility details, safety compliance, federation endorsement, and fee payment
4. Approve, reject, or request corrections

---

## API Endpoints

- `POST /api/v1/applications/facility-operation`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
