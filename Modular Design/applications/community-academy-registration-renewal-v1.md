# Community Sports Club Registration and Renewal Application — NCSMS v1.0

**Form Reference:** Form 10 — Regulation 22(1)
**Submitted To:** The General Secretary, National Council of Sports

## Purpose

This form supports both new registration and renewal for community sports clubs under national sports governance.

---

## Form Fields

### Application Type (Top of Form)

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `application_type` | radio | Application For | Options: **Registration** \| **Renewal** |

---

### Section A — Applicant Details

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `applicant_name` | text | Name | Full name of the community sports club |
| `physical_address` | textarea | Physical Address | |
| `postal_address` | text | Postal Address | |
| `telephone_fixed_line` | text | Telephone (Fixed Line) | |
| `mobile_phone` | text | Mobile Phone | |
| `email_address` | email | E-mail Address | |
| `website` | url | Website | Optional |

---

### Section B — Details of Community Sports Club

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `sports_disciplines` | textarea | Categories of Sports Disciplines Promoted | List all sports categories the club promotes |
| `date_formed` | date | Date Club Was Formed | |
| `membership_list` | textarea | List of Club Membership | Include list of certified technical persons (coaches, officials, etc.) |
| `ownership_and_location` | textarea | Details of Ownership and Location of the Club | State ownership arrangement and full location details |
| `federation_recommendation` | textarea | Recommendation by Respective NSA / Federation | State whether the club has been recommended; if yes, name the recommending body (where applicable) |
| `other_affiliations` | textarea | Other National and International Bodies the Club is Affiliated To | Where applicable |

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

---

## Tracking Metadata

- `application_reference` — unique submission identifier assigned on submission
- `payment_reference` — application fee payment tracking number
- `renewal_cycle` — applicable only for renewal submissions
- `status` — `DRAFT` \| `SUBMITTED` \| `UNDER_REVIEW` \| `APPROVED` \| `REJECTED` \| `NEEDS_INFORMATION`
- `submitted_at`, `last_updated`, `next_action`

---

## Workflow

1. Applicant selects Registration or Renewal and completes all sections
2. Uploads proof of fee payment
3. Submits to `GENERAL_SECRETARY` for review
4. Reviewer verifies club details, membership list, and payment
5. Approve, request clarifications, or reject with guidance

---

## API Endpoints

- `POST /api/v1/applications/community-club-registration`
- `POST /api/v1/applications/community-club-renewal`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
