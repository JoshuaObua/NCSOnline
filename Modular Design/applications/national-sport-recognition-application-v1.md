# National Sport Recognition Application — NCSMS v1.0

**Form Reference:** Form 1 — Regulation 3(1)
**Submitted To:** The General Secretary, National Council of Sports

## Purpose

This form supports applications for the declaration of a sport as a national sport in Uganda.

---

## Form Fields

### Part I — Particulars of Applicant

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `applicant_name` | text | Name | Full name of applicant |
| `applicant_type` | radio | Applicant Status | Options: **Citizen** \| **Resident** |
| `physical_address` | textarea | Physical Address | |
| `postal_address` | text | Postal Address | |
| `telephone_fixed_line` | text | Telephone (Fixed Line) | |
| `mobile_phone` | text | Mobile Phone Number | |
| `email_address` | email | E-mail Address | |
| `website` | url | Website | Optional |

---

### Part II — Nature of the Sport

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `sport_name` | text | Name of the Sport | |
| `sport_description` | textarea | Brief Description of How the Sport is Played | |

---

### Part III — Popularity of the Sport

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `popularity_by_district_region` | textarea | Level of Engagement by District or Region | Explain coverage across districts/regions |

---

### Part IV — Social Economic Impact of the Sport

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `socioeconomic_impact` | textarea | Social-Economic Impact of the Sport to the Country | |

---

### Part V — International Recognition (Where Applicable)

Repeatable table — add rows as needed:

| Field Name | Type | Label |
|---|---|---|
| `intl_body_name` | text | International Sports Body |
| `intl_affiliation_date` | date | Date of Affiliation |

---

### Part VI — Presence of Sports Facilities for the Sport

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `facility_districts` | textarea | Districts Where Facilities Exist | Specify which districts |
| `facility_other_locations` | textarea | Other Locations (Specify) | |

---

### Part VII — Leadership and Governance Structure

Repeatable table with fixed roles:

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `leader_role` | select | Role | Options: President or equivalent \| General Secretary or equivalent \| Treasurer |
| `leader_name` | text | Name | |
| `leader_address` | text | Address | |
| `leader_nationality` | text | Nationality | |
| `leader_country_of_residence` | text | Country of Usual Residence | |

---

### Part VIII — Sources of Funding

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
| `signature_upload` | file | Signature / Seal | Upload scanned signature and seal |
| `date_signed` | date | Date Signed | |

---

## Required Attachments

| # | Field Name | Label | Required |
|---|---|---|---|
| 1 | `work_plan_document` | Work Plan to Develop the Sport | Mandatory |

---

## Tracking Metadata

- `application_reference` — unique submission identifier assigned on submission
- `payment_reference` — application fee payment tracking number
- `status` — `DRAFT` \| `SUBMITTED` \| `UNDER_REVIEW` \| `APPROVED` \| `REJECTED` \| `NEEDS_INFORMATION`
- `submitted_at`, `last_updated`, `next_action`

---

## Workflow

1. Applicant saves a draft and completes all parts
2. Submits for review — assigned to `GENERAL_SECRETARY`
3. Reviewer assesses sport legitimacy, district coverage, governance readiness, and development plan
4. Approve, request clarifications, or reject with written notes

---

## API Endpoints

- `POST /api/v1/applications/national-sport-recognition`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
