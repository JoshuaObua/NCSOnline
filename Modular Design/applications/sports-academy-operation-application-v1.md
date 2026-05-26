# Sports Academy Operation Application — NCSMS v1.0

**Form Reference:** Form 11 — Regulation 29(1)(a)
**Submitted To:** The General Secretary, National Council of Sports

## Purpose

This application collects the information required to approve and licence a sports academy to operate in Uganda.

> **Legal Notice:** The Academy MUST ensure compliance with all relevant laws (including the Children's Act) to safeguard the welfare, rights, and best interests of the child and education.

---

## Form Fields

### Section C — Applicant Details

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `applicant_name` | text | Name | Full name of the applicant or authorised representative |
| `physical_address` | textarea | Physical Address | |
| `postal_address` | text | Postal Address | |
| `telephone_fixed_line` | text | Telephone (Fixed Line) | |
| `mobile_phone` | text | Mobile Phone | |
| `email_address` | email | E-mail Address | |
| `website` | url | Website | Optional |

---

### Section D — Details of Sports Academy

| Field Name | Type | Label | Notes |
|---|---|---|---|
| `sports_disciplines` | textarea | Categories of Sports Disciplines Promoted | List all sports categories and disciplines the academy promotes |
| `date_formed` | date | Date Sports Academy Was Formed | |
| `membership_list` | textarea | List of Membership of the Sports Academy | Include list of certified technical persons (coaches, trainers, officials) |
| `ownership_and_location` | textarea | Details of Ownership and Location of the Sports Academy and Sports Facility | State ownership arrangement, full address, and facility location |
| `incorporation_details` | textarea | How the Sports Academy is Incorporated in Uganda | Describe the legal incorporation form (company, NGO, trust, etc.) |
| `federation_recommended` | radio | Has the Academy Been Recommended by a Respective NSA / Federation? | Options: **Yes** \| **No** \| **Not Applicable** |
| `federation_recommendation_details` | textarea | Recommendation Details | State which NSA/federation recommended the academy (required if Yes) |
| `other_affiliations` | textarea | Other National and International Bodies the Academy is Affiliated To | If applicable |
| `has_technical_financial_proposal` | radio | Does the Academy Have a Technical and Financial Proposal? | Options: **Yes** \| **No** |

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
| 2 | `certificate_of_incorporation` | Certified Copy of Certificate of Incorporation | Conditional — required if the academy is incorporated |
| 3 | `federation_approval_certificate` | Certificate of Approval from Respective NSA / Federation | Conditional — required if `federation_recommended` = Yes |
| 4 | `technical_financial_proposal_document` | Technical and Financial Proposal for Managing the Establishment | Conditional — required if `has_technical_financial_proposal` = Yes |

---

## Wizard Steps

This form follows the shared **6-stage application workflow** defined in [application-workflow-v1.md](application-workflow-v1.md). The form content maps to the following wizard steps:

| Step | Label | Fields Covered |
|---|---|---|
| 1 | Applicant Details | Name, physical address, postal address, telephone, mobile, email, website |
| 2 | Academy Details | Sports disciplines promoted, date formed, list of membership including certified technical persons |
| 3 | Ownership & Incorporation | Ownership and location details, how the academy is incorporated in Uganda |
| 4 | Compliance, Affiliations & Proposals | NSA/federation recommendation (conditional), other affiliations, technical and financial proposal (conditional); Children's Act compliance acknowledgement |
| 5 | Review & Confirm | Read-only summary; edit links per section |

After Step 5 the applicant proceeds through **Stage 3 (Download, Sign & Upload)** → **Stage 4 (Payment)** → **Stage 5 (Submit)** → **Stage 6 (Track)**.

> Step 4 contains **conditional fields**: the federation recommendation certificate upload appears only if `federation_recommended = Yes`; the technical proposal upload appears only if `has_technical_financial_proposal = Yes`.

### Payment Options
- **Digital payment** — online card or mobile money
- **Upload proof of payment** — scanned bank slip or mobile money receipt

---

## Tracking Metadata

- `application_reference` — unique submission identifier assigned on submission
- `payment_reference` — application fee payment tracking number
- `facility_registration_number` — assigned registration number for the facility
- `status` — `DRAFT` \| `SUBMITTED` \| `UNDER_REVIEW` \| `APPROVED` \| `REJECTED` \| `NEEDS_INFORMATION`
- `submitted_at`, `last_updated`, `next_action`

---

## Workflow

1. Applicant completes all sections, answers conditional questions, and uploads required attachments
2. Submits to `GENERAL_SECRETARY` for review
3. Reviewer validates academy details, incorporation status, federation endorsement, and child protection compliance
4. Issue approval, request additional information, or reject the application

---

## API Endpoints

- `POST /api/v1/applications/academy-operation`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
