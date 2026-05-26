# Application Submission Workflow Architecture — NCSMS v1.0

## Overview

Every application form in the NCSMS follows a single shared six-stage wizard workflow. This document defines that workflow in full — from the moment a signed-in user opens a form through to final submission and status tracking. All 8 application form types (Forms 1, 3, 5, 7, 8, 10, 11) use this architecture.

---

## Workflow Stages at a Glance

```
STAGE 1          STAGE 2         STAGE 3                STAGE 4              STAGE 5       STAGE 6
FILL FORM   →   REVIEW    →   DOWNLOAD & SIGN   →   PAYMENT       →   SUBMIT    →   TRACK
(Wizard)        (Summary)      (PDF + Upload)        (Digital or              (Confirm)     (Dashboard)
                                                      Proof Upload)
```

Progress is **saved automatically** at every stage. A user can leave, log back in, and resume from exactly where they stopped.

---

## Stage 1 — Fill Form (Multi-step Wizard)

The form is broken into logical wizard steps matching the PDF parts. Each step corresponds to one section of the official statutory form.

### Behaviour
- Each step renders only its own fields — no scrolling through a long single page
- A **progress bar** at the top shows step number and label (e.g., "Step 2 of 6 — Nature of Sport")
- Validation runs **on exit from each step** — the user cannot advance with invalid or empty required fields
- Auto-save fires **on every field change** (debounced 2 s) and on "Next"
- A **"Save Draft"** button is visible on every step
- Draft is associated with the signed-in user and the specific form type
- Returning users see a **"Resume Application"** prompt linking back to their last active step

### Draft persistence
- `draft_id` is created on first field entry
- `draft_payload` stores the full form state as JSON keyed by `field_name`
- `last_saved_step` records the furthest completed step for resume targeting
- `draft_expires_at` — drafts older than 90 days are flagged for deletion; the user is warned at 80 days

### Wizard Step Numbering Convention

Each form file defines its own step list. The final step of every form is always **Review & Confirm** (Step N), which is handled by Stage 2 below.

---

## Stage 2 — Review

The applicant sees a **read-only summary** of all entered data before proceeding.

- Each section shows a collapsible card with the field values as filled
- An **"Edit"** link on each card jumps back to that wizard step
- Required fields that are still empty are highlighted in red
- The system validates completeness of the full form at this point
- A **"Looks Good, Continue"** button advances to Stage 3

---

## Stage 3 — Download, Sign, and Upload

This stage handles the applicant's signature on the official statutory form.

### Step A — Download Pre-filled PDF
- The system generates a pre-filled PDF of the form using the data entered in Stage 1
- PDF mirrors the layout of the official statutory form exactly
- A **"Download PDF"** button delivers the file
- The `pdf_generated_at` timestamp is recorded

### Step B — Sign the Form
Two options are presented:

| Option | Label | Description |
|---|---|---|
| Digital Signature | "Sign Digitally" | Applicant signs inline using a signature pad or uploads a saved signature image; signature is embedded in the PDF in the browser |
| Print and Scan | "Print, Sign & Upload" | Applicant prints the PDF, signs physically, scans or photographs it, and uploads the scan |

Either option produces the same outcome: a **signed PDF file** attached to the application.

### Step C — Upload Signed PDF
- File upload input: `signed_form_pdf`
- Accepted formats: `.pdf`, `.jpg`, `.png` (max 10 MB)
- Preview shown after upload
- Upload can be replaced if the wrong file was submitted
- `signed_form_uploaded_at` timestamp recorded
- Status advances to `PENDING_PAYMENT`

---

## Stage 4 — Payment

Payment is required before final submission. **Two payment methods are accepted** — for both new applications and renewals.

### Payment Method A — Digital / Online Payment

| Field | Type | Notes |
|---|---|---|
| `payment_method` | radio (pre-selected) | **Online Payment** |
| `payment_gateway_reference` | auto | Populated by payment gateway on success |
| `payment_amount_ugx` | display | Prescribed fee amount shown read-only |
| `payment_currency` | display | UGX |

- Supported channels: credit/debit card, mobile money (MTN MoMo, Airtel Money)
- On payment success the gateway calls back with a confirmation; `payment_status` is set to `PAID`
- On failure the user is shown an error and can retry

### Payment Method B — Upload Proof of Payment

For applicants who pay via bank deposit, bank transfer, or any offline channel:

| Field | Type | Notes |
|---|---|---|
| `payment_method` | radio | **Upload Proof of Payment** |
| `payment_proof_document` | file | Scanned bank slip, mobile money receipt, or payment confirmation |
| `payment_date` | date | Date the payment was made |
| `paying_bank_or_channel` | text | Bank name or mobile money provider |
| `depositor_name` | text | Name of the person who made the payment |
| `payment_reference_number` | text | Reference number on the receipt / slip |
| `payment_amount_ugx` | number | Amount paid as shown on the receipt |

- Accepted file formats: `.pdf`, `.jpg`, `.png` (max 5 MB)
- Payment proof is reviewed by `GENERAL_SECRETARY` as part of the application review
- `payment_status` is set to `PROOF_UPLOADED` until verified

### Payment Status Values

| Value | Meaning |
|---|---|
| `UNPAID` | No payment action taken yet |
| `PAYMENT_INITIATED` | Online payment started but not confirmed |
| `PAID` | Online payment confirmed by gateway |
| `PROOF_UPLOADED` | Scanned proof uploaded, pending reviewer verification |
| `PAYMENT_VERIFIED` | Reviewer has confirmed proof of payment |
| `PAYMENT_REJECTED` | Proof was rejected (wrong amount, unreadable, etc.) |

### Save Progress
- Payment stage state is saved as part of the application record
- If the user navigates away after uploading proof, the upload is preserved
- The user can return and switch payment method before final submission

---

## Stage 5 — Submission

### Pre-submission Checklist

The system shows a final checklist confirming all required items are present:

- [x] Form fully completed
- [x] Signed form PDF uploaded (`signed_form_pdf`)
- [x] All required attachments uploaded
- [x] Payment completed or proof uploaded (`payment_status` is `PAID` or `PROOF_UPLOADED`)

The **"Submit Application"** button is enabled only when all checklist items are satisfied.

### On Submit
- `application_status` transitions from `PENDING_PAYMENT` → `SUBMITTED`
- `submitted_at` timestamp is recorded
- `application_reference` is generated and displayed to the applicant (e.g., `NCS-2026-F3-00142`)
- A confirmation notification is sent to the applicant's registered email / phone
- The application is placed in the `GENERAL_SECRETARY` review queue

---

## Stage 6 — Tracking Dashboard

After submission the applicant can monitor their application from their dashboard.

### Application Card (per submission)

| Field | Display |
|---|---|
| `application_reference` | Unique reference number |
| `form_type` | e.g., "Federation Registration (Form 3)" |
| `submitted_at` | Submission date |
| `application_status` | Current status badge |
| `payment_status` | Current payment status badge |
| `next_action` | Prompt if applicant action is needed |
| `reviewer_notes` | Visible when reviewer has added notes |

### Status State Machine

```
DRAFT
  └→ PENDING_SIGNATURE        (signed form not yet uploaded)
       └→ PENDING_PAYMENT      (signed form uploaded, payment outstanding)
            └→ SUBMITTED        (payment done / proof uploaded)
                 └→ UNDER_REVIEW  (assigned to GENERAL_SECRETARY)
                      ├→ APPROVED
                      ├→ REJECTED
                      └→ NEEDS_INFORMATION → [applicant responds] → RESUBMITTED → UNDER_REVIEW
```

### NEEDS_INFORMATION Flow
- Reviewer sets status to `NEEDS_INFORMATION` and adds a note
- Applicant receives a notification
- Applicant can view the reviewer's note and upload additional documents or amend specific fields
- Applicant resubmits — status becomes `RESUBMITTED`
- Application re-enters the `UNDER_REVIEW` queue

---

## API Endpoints

### Draft Management
- `POST /api/v1/applications/{form-type}/draft` — create or update a draft
- `GET /api/v1/applications/{form-type}/draft` — load the latest draft for the signed-in user
- `DELETE /api/v1/applications/{form-type}/draft` — discard a draft

### Signed Form
- `POST /api/v1/applications/{id}/signed-form` — upload signed PDF
- `GET /api/v1/applications/{id}/signed-form` — download the uploaded signed PDF
- `POST /api/v1/applications/{id}/prefilled-pdf` — generate and download the pre-filled form PDF

### Payment
- `POST /api/v1/applications/{id}/payment/initiate` — start online payment session
- `POST /api/v1/applications/{id}/payment/proof` — upload proof of payment
- `GET /api/v1/applications/{id}/payment` — get current payment status

### Submission and Tracking
- `POST /api/v1/applications/{id}/submit` — final submission
- `GET /api/v1/applications/{id}` — get full application record and status
- `GET /api/v1/applications` — list all applications for the signed-in user
- `PATCH /api/v1/applications/{id}/respond` — respond to NEEDS_INFORMATION request

### Attachments
- `POST /api/v1/applications/{id}/attachments` — upload a required attachment
- `GET /api/v1/applications/{id}/attachments` — list uploaded attachments
- `DELETE /api/v1/applications/{id}/attachments/{attachment-id}` — remove an attachment

---

## Frontend Component Map

| Component | Purpose |
|---|---|
| `ApplicationWizard` | Top-level multi-step container; owns stage routing |
| `WizardProgress` | Progress bar with step labels and completion indicators |
| `WizardStep` | Individual step wrapper with "Back" / "Next" / "Save Draft" controls |
| `FormReview` | Read-only summary with per-section Edit links |
| `PDFDownloadButton` | Triggers pre-filled PDF generation and download |
| `SignatureCapture` | Inline digital signature pad (canvas-based) |
| `FileUpload` | Reusable drag-and-drop file input with preview and replace |
| `PaymentSelector` | Payment method radio with conditional form sections |
| `OnlinePaymentGateway` | Embedded payment frame / redirect for card and mobile money |
| `ProofOfPaymentForm` | Form fields for offline payment details + file upload |
| `SubmissionChecklist` | Pre-submit checklist with live validation |
| `ApplicationStatusCard` | Dashboard card showing status badge, reference, and next action |
| `ReviewerNotePanel` | Displays reviewer comments for NEEDS_INFORMATION state |
| `NeedsInformationResponse` | Allows applicant to upload docs and resubmit |

---

## Shared Application Status Codes

| Status | Triggered By |
|---|---|
| `DRAFT` | First field entry or explicit "Save Draft" |
| `PENDING_SIGNATURE` | Form review completed, advancing to Stage 3 |
| `PENDING_PAYMENT` | Signed form successfully uploaded |
| `SUBMITTED` | Final submit with valid payment |
| `UNDER_REVIEW` | General Secretary picks up the application |
| `NEEDS_INFORMATION` | Reviewer requests clarification or additional docs |
| `RESUBMITTED` | Applicant responds to NEEDS_INFORMATION |
| `APPROVED` | General Secretary approves application |
| `REJECTED` | General Secretary rejects application |
