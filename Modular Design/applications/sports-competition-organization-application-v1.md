# Sports Competition Organisation Application — NCSMS v1.0

**Form Reference:** Form 7 — Regulation 16(1)
**Submitted To:** General Secretary, National Council of Sports

## Purpose

This application is submitted to obtain clearance to organise a sports competition in Uganda.

---

## Form Fields

### Section 1 — Applicant Details

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `applicant_name` | text | Name of Applicant | Full name of the organising entity or individual |
| `physical_address` | textarea | Physical Address | |
| `postal_address` | text | Postal Address | |
| `telephone_fixed_line` | text | Telephone (Fixed Line) | |
| `mobile_phone` | text | Mobile Phone | |
| `email_address` | email | E-mail Address | |
| `website` | url | Website | Optional |

---

### Section 2 — Competition Details

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `sport_nature_and_categories` | textarea | Nature of Sport and Categories for Competition | Describe the sport(s) and the competition categories |
| `competition_start_date` | date | Competition Start Date | |
| `competition_end_date` | date | Competition End Date | |
| `budget_and_funding` | textarea | Confirmation of Budget, Availability of Funds, and Sources of Funding | State budget total, funding sources, and confirmation of available funds |
| `hosting_venue` | textarea | Hosting Facility / Venue | Full name and address of the competition venue |
| `accommodation_address` | textarea | Accommodation Sites Address | If applicable — address(es) of accommodation for participants |
| `organizing_committee` | textarea | Organising Committee Details | State and confirm availability of the organising committee (if applicable) |

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
| 1 | `executive_committee_minutes` | Minutes of Executive Committee / Board Meeting Approving the Competition | Mandatory |
| 2 | `approved_action_plan` | Copy of Approved Action Plan for the Financial Year | Mandatory |
| 3 | `competition_fixture` | Copy of Competition Fixture | Mandatory |
| 4 | `sponsors_list` | List of Sponsors | Mandatory |

---

## Wizard Steps

This form follows the shared **6-stage application workflow** defined in [application-workflow-v1.md](application-workflow-v1.md). The form content maps to the following wizard steps:

| Step | Label | Fields Covered |
|---|---|---|
| 1 | Applicant Details | Name, physical address, postal address, phone, email, website |
| 2 | Competition Details | Sport nature and categories, competition start date, competition end date |
| 3 | Logistics & Organisation | Budget and funding confirmation, hosting venue, accommodation address, organising committee details |
| 4 | Review & Confirm | Read-only summary; edit links per section |

After Step 4 the applicant proceeds through **Stage 3 (Download, Sign & Upload)** → **Stage 4 (Payment)** → **Stage 5 (Submit)** → **Stage 6 (Track)**.

### Payment Options
- **Digital payment** — online card or mobile money
- **Upload proof of payment** — scanned bank slip or mobile money receipt

---

## Tracking Metadata

- `application_reference` — unique submission identifier assigned on submission
- `competition_reference` — reference number specific to this competition clearance
- `payment_reference` — application fee payment tracking number
- `status` — `DRAFT` \| `SUBMITTED` \| `UNDER_REVIEW` \| `APPROVED` \| `REJECTED` \| `NEEDS_INFORMATION`
- `submitted_at`, `last_updated`, `next_action`

---

## Workflow

1. Applicant fills in all competition details and uploads 4 required attachments
2. Submits to `GENERAL_SECRETARY` for clearance review
3. Reviewer validates sanctioning, venue fitness, budget confirmation, and organising committee readiness
4. Approve, request clarifications, or reject the application

---

## API Endpoints

- `POST /api/v1/applications/competition-organisation`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
