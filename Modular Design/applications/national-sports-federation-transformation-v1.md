# National Sports Federation Transformation Form — NCSMS v1.0

## Purpose

This form captures transformation requests for national sports federations undergoing governance or strategic changes.

## Key Sections

### Federation Details

- Federation name and registration number
- Current leadership and governing council
- Sports disciplines covered and membership base

### Transformation Proposal

- Proposed changes in governance or organizational scope
- Updated membership, affiliation, or federation structure
- Rationale for transformation and expected benefits

### Legal and Compliance Evidence

- Constitution amendment proposals
- Governing council approvals or membership resolutions
- Legal opinions and regulatory notifications
- Supporting evidence for the proposed transformation

### Supporting Documents and Fees

- Historical compliance record
- Stakeholder consultation and approval letters
- Payment receipt and `payment_reference`

## Tracking and References

- `application_reference`
- `transformation_reference`
- `payment_reference`
- `status`, `submitted_at`, `updated_at`

## Workflow

- Save draft transformation details for review
- Submit to `GENERAL_SECRETARY` for validation
- Verify governance changes, legal approvals, and stakeholder support
- Approve, request clarification, or reject

## API Endpoints

- `POST /api/v1/applications/transformation/federation`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
