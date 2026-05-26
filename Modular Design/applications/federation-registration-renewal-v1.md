# Federation Registration and Renewal Application — NCSMS v1.0

## Purpose

This form supports both new federation registration and periodic license renewal for existing federations.

## Key Sections

### Applicant Information

- Applicant name and user ID
- Email address and phone number
- Organization or federation representative details
- Primary contact and secondary contact

### Federation Details

- Federation name
- Federation code and registration number
- Sport discipline(s) covered
- Headquarters address and physical location
- Membership count and affiliated clubs

### Legal and Governance Documents

- Constitution and bylaws
- Governing council membership list
- Evidence of legal registration or incorporation
- Latest annual report or governance statement

### Renewal-specific Fields

- Current license number
- Expiry date
- Renewal cycle requested
- Previous compliance and audit summary
- Renewal justification or change summary

### Supporting Evidence

- Proof of payment for application fee
- Bank payment reference or receipt upload
- Federation bank details and account verification
- Certification of compliance with national sports regulations

## Tracking and References

- `application_reference` — unique submission identifier
- `payment_reference` — license fee payment tracking number
- `review_notes` — comments from the reviewer
- `application_status` — `DRAFT`, `SUBMITTED`, `UNDER_REVIEW`, `APPROVED`, `REJECTED`, `NEEDS_INFORMATION`

## Workflow

- Draft form saved for later completion
- Submit for review and assign to `GENERAL_SECRETARY`
- Payment verification and documentation review
- Approval or rejection with notes and next steps
- Renewal requests may be flagged for additional compliance checks

## API Endpoints

- `POST /api/v1/applications/federation-registration`
- `POST /api/v1/applications/federation-renewal`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
