# Community Sports Academy Registration and Renewal — NCSMS v1.0

## Purpose

This form supports community sports academies that need registration or renewal under local and national sports governance.

## Key Sections

### Applicant and Community Details

- Community academy name
- Community group or local council details
- Location and operating neighborhood
- Sports offered and target participant groups
- Governance structure and community leadership

### Registration or Renewal Information

- Existing registration number (if renewal)
- Renewal justification and track record
- Membership or participant growth details
- Community impact and inclusion metrics

### Compliance and Support Documentation

- Local authority endorsement or partnership letters
- Health and safety policy and facility inspection
- Coach qualifications or volunteer training plans
- Evidence of funds and sustainability plan

### Payments and Tracking

- Application fee payment receipt
- `application_reference`
- `payment_reference`
- `renewal_cycle`
- `review_status`

## Workflow

- Work in progress draft mode for community applicants
- Submit for review by `GENERAL_SECRETARY`
- Verify supporting evidence and local approvals
- Approve, request clarifications, or reject with guidance

## API Endpoints

- `POST /api/v1/applications/community-academy-registration-renewal`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
