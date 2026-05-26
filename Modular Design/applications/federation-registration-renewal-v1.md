# Federation / Association Registration and Renewal Application — NCSMS v1.0

**Form Reference:** Form 3 — Regulation 4(1), 5(1), 9(1)
**Submitted To:** The General Secretary, National Council of Sports

## Purpose

This form supports both new registration and periodic renewal for national sports associations (NSA) and national sports federations (NSF).

---

## Form Fields

### Application Type (Top of Form)

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `application_type` | radio | Application For | Options: **Registration** \| **Renewal** |
| `organisation_type` | radio | Organisation Type | Options: **National Sports Association** \| **National Sports Federation** |

---

### Part I — Particulars of Applicant

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `applicant_name` | text | Name | Full name of the association or federation |
| `physical_address` | textarea | Physical Address | |
| `postal_address` | text | Postal Address | |
| `telephone_fixed_line` | text | Telephone (Fixed Line) | |
| `mobile_phone` | text | Mobile Phone | |
| `email_address` | email | E-mail Address | |
| `website` | url | Website | Optional |

---

### Part II — Nature of Sport

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `sport_promoted` | text | Sport Promoted and Supervised | Indicate the specific sport |
| `coverage_districts` | textarea | District Coverage | Specify which districts and percentage |
| `coverage_other` | textarea | Other Coverage (Specify) | Any other coverage areas |

---

### Part III — Previous Legal Status (Where Applicable)

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `previous_registration_as` | text | Previously Registered or Incorporated As | How it was previously registered in Uganda |
| `previous_registration_date` | date | Date of Previous Registration or Incorporation | |

---

### Part IV — Leadership and Governance Structure

**Leadership Table** (fixed roles):

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `leader_role` | select | Role | Options: President/Chairperson or equivalent \| General Secretary/CEO or equivalent \| Treasurer or equivalent |
| `leader_name` | text | Name | |
| `leader_address_and_phone` | text | Physical Address and Telephone Number | Combined field as per form |
| `leader_nationality` | text | Nationality | |
| `leader_country_of_residence` | text | Country of Usual Residence | |

**Accounting Officer** (separate table):

| Field Name | Type | Label |
|---|---|---|
| `accounting_officer_name` | text | Name |
| `accounting_officer_address` | text | Physical Address |
| `accounting_officer_phone` | text | Telephone Number |

**Governance Description:**

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `governance_structure_description` | textarea | Governance Structure Description | Describe general assembly, executive committee or board, and secretariat |

---

### Part IV — International Affiliations

Repeatable table — add rows as needed:

| Field Name | Type | Label |
|---|---|---|
| `intl_body_name` | text | International Sports Body (Global, Continental or Regional) |
| `intl_affiliation_date` | date | Date of Affiliation |

---

### Part V — Sources of Funding

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
| 6 | `general_assembly_minutes_constitution` | Original Copy of Minutes of the General Assembly that Approved the Constitution | Mandatory |
| 7 | `members_list` | Certified and Updated List of Members of the NSA or NSF | Mandatory |
| 8 | `sports_activities_report` | Report of Sports Activities Conducted Within One Year Prior to Application | Mandatory |
| 9 | `districts_presence_list` | List of Districts Where the NSA or NSF Has Presence and is Active | Mandatory |
| 10 | `election_minutes` | Minutes of General Assembly that Elected the Executive Committee | Mandatory |
| 11 | `audited_accounts` | Audited Books of Accounts | Optional (where applicable) |
| 12 | `id_documents` | Passport Photos and Certified Copies of ID Documents (National ID or Passport) of Executive Committee or Board | Mandatory |

---

## Wizard Steps

This form follows the shared **6-stage application workflow** defined in [application-workflow-v1.md](application-workflow-v1.md). The form content maps to the following wizard steps:

| Step | Label | Parts / Fields Covered |
|---|---|---|
| 1 | Application Type | Select Registration or Renewal; select NSA or NSF |
| 2 | Applicant Particulars | Part I — Organisation name, address, contacts, website |
| 3 | Nature of Sport & Legal Status | Part II — Sport promoted, district coverage; Part III — Previous registration (where applicable) |
| 4 | Leadership & Governance | Part IV — Leadership table (President, Gen. Secretary, Treasurer), Accounting Officer, governance structure description |
| 5 | International Affiliations & Funding | Part IV Affiliations — international bodies table; Part V — funding sources table |
| 6 | Review & Confirm | Read-only summary of all sections; edit links per section |

After Step 6 the applicant proceeds through **Stage 3 (Download, Sign & Upload)** → **Stage 4 (Payment)** → **Stage 5 (Submit)** → **Stage 6 (Track)**.

> Renewal applications follow the same wizard. Previously submitted data is pre-populated where available; the applicant confirms or updates each section.

### Payment Options
- **Digital payment** — online card or mobile money
- **Upload proof of payment** — scanned bank slip or mobile money receipt

Both methods accepted for both Registration and Renewal submissions.

---

## Tracking Metadata

- `application_reference` — unique submission identifier assigned on submission
- `payment_reference` — license fee payment tracking number
- `review_notes` — comments from the reviewer
- `status` — `DRAFT` \| `SUBMITTED` \| `UNDER_REVIEW` \| `APPROVED` \| `REJECTED` \| `NEEDS_INFORMATION`
- `submitted_at`, `last_updated`, `next_action`

---

## Workflow

1. Applicant selects Registration or Renewal, then NSA or NSF
2. Completes all parts and uploads all 12 required attachments
3. Submits for review — assigned to `GENERAL_SECRETARY`
4. Payment verified and documents reviewed
5. Approve or reject with notes; renewal requests may trigger compliance checks

---

## API Endpoints

- `POST /api/v1/applications/federation-registration`
- `POST /api/v1/applications/federation-renewal`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
