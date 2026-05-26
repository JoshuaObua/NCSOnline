# National Sports Association Transformation Application — NCSMS v1.0

**Form Reference:** Form 5 — Regulation 11(1)
**Submitted To:** The General Secretary, National Council of Sports

## Purpose

This form is used by a registered national sports association (NSA) to apply for transformation into a national sports federation (NSF). The NSA is the applicant and submits this form to initiate the change of status.

> **Note:** This form is the same as `national-sports-federation-transformation-v1.md` — both document FORM 5. This file describes the process from the NSA applicant's perspective.

---

## Form Fields

### Introduction

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `association_name_to_transform` | text | Name of National Sports Association | Insert the exact registered name of the NSA applying for transformation |

---

### Applicant Contact Details

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `applicant_name` | text | Name | Name of authorised representative submitting the application |
| `physical_address` | textarea | Physical Address | |
| `postal_address` | text | Postal Address | |
| `telephone_fixed_line` | text | Telephone (Fixed Line) | |
| `mobile_phone` | text | Mobile Phone | |
| `email_address` | email | E-mail Address | |
| `website` | url | Website | Optional |

---

### Association Details

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `sport_promoted` | text | Sport Promoted by the Association | Indicate the sport |
| `coverage_districts` | textarea | District Coverage | Specify which districts and percentage |
| `coverage_other` | textarea | Other Coverage (Specify) | |
| `previous_registration_as` | text | Previously Registered in Uganda As | Form of previous registration |
| `previous_registration_date` | date | Date of Previous Registration | |
| `transformation_reason` | textarea | Reason for Transformation into a National Sports Federation | Detailed justification |

---

### Leadership and Governance Structure

**Leadership Table** (fixed roles):

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `leader_role` | select | Role | Options: President/Chairperson or equivalent \| General Secretary/CEO or equivalent \| Treasurer or equivalent |
| `leader_name` | text | Name | |
| `leader_address_and_phone` | text | Physical Address and Telephone Number | |
| `leader_nationality` | text | Nationality | |
| `leader_country_of_residence` | text | Country of Usual Residence | |

**Accounting Officer** (separate table):

| Field Name | Type | Label |
|---|---|---|
| `accounting_officer_name` | text | Name |
| `accounting_officer_address` | text | Physical Address |
| `accounting_officer_phone` | text | Telephone Number |

---

### International Affiliations

Repeatable table — add rows as needed:

| Field Name | Type | Label |
|---|---|---|
| `intl_body_name` | text | International Sports Body |
| `intl_affiliation_date` | date | Date of Affiliation |

---

### Sources of Funding

Repeatable table — add rows as needed:

| Field Name | Type | Label |
|---|---|---|
| `funding_source` | text | Source |
| `funding_amount_ugx` | number | Amount (UGX) |

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
| 2 | `organogram_document` | Organogram of National Sports Organisation | Mandatory |
| 3 | `national_sport_certificate` | Certificate of Recognition as National Sport | Mandatory |
| 4 | `symbols_slogans_colours_document` | Symbols, Slogans, and Colours of the Applicant NSA or Federation | Mandatory |
| 5 | `constitution_document` | Copy of Constitution Approved by Members of the General Assembly | Mandatory |
| 6 | `general_assembly_minutes_constitution` | Original Copy of Minutes of the General Assembly Approving the Constitution | Mandatory |
| 7 | `members_list` | Certified and Updated List of Members | Mandatory |
| 8 | `sports_activities_report` | Report of Sports Activities Conducted Within One Year Prior to Application | Mandatory |
| 9 | `districts_presence_list` | List of Districts Where the NSA or NSF Has Presence and is Active | Mandatory |
| 10 | `election_minutes` | Minutes of General Assembly that Elected the Executive Committee | Mandatory |
| 11 | `audited_accounts` | Audited Books of Accounts | Optional (where applicable) |
| 12 | `id_documents` | Passport Photos and Certified Copies of ID Documents (National ID or Passport) of Executive Committee or Board | Mandatory |

---

## Wizard Steps

This form follows the shared **6-stage application workflow** defined in [application-workflow-v1.md](application-workflow-v1.md). The form content maps to the following wizard steps:

| Step | Label | Fields Covered |
|---|---|---|
| 1 | Association Identity | `association_name_to_transform` — exact registered name of the NSA |
| 2 | Applicant Contact Details | Name of authorised representative, physical address, postal address, phone, email, website |
| 3 | Association Details & Legal History | Sport promoted, district coverage, previous registration form and date |
| 4 | Transformation Rationale | `transformation_reason` — detailed justification for the NSA's request to become an NSF |
| 5 | Leadership, Governance, Affiliations & Funding | Leadership table, accounting officer, international affiliations table, funding sources table |
| 6 | Review & Confirm | Read-only summary; edit links per section |

After Step 6 the applicant proceeds through **Stage 3 (Download, Sign & Upload)** → **Stage 4 (Payment)** → **Stage 5 (Submit)** → **Stage 6 (Track)**.

### Payment Options
- **Digital payment** — online card or mobile money
- **Upload proof of payment** — scanned bank slip or mobile money receipt

---

## Tracking Metadata

- `application_reference` — unique submission identifier
- `transformation_reference` — specific reference for transformation tracking
- `payment_reference` — fee payment tracking number
- `status` — `DRAFT` \| `SUBMITTED` \| `UNDER_REVIEW` \| `APPROVED` \| `REJECTED` \| `NEEDS_INFORMATION`
- `submitted_at`, `last_updated`, `next_action`

---

## Workflow

1. NSA completes all form sections and uploads all 12 required attachments
2. Submits to `GENERAL_SECRETARY` for validation
3. Reviewer confirms transformation rationale, governance compliance, membership approvals, and legal basis
4. Approve, request additional evidence, or reject with written notes

---

## API Endpoints

- `POST /api/v1/applications/transformation/nsa-to-nsf`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
