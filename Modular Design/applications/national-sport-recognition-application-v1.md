# National Sport Recognition Application — NCSMS v1.0

## Purpose

This application captures the requirements for recognizing a sport at the national level.

## Key Sections

### Sport Information

- Sport name
- Sport category and discipline
- Historical practice summary
- Governing body or representative organization

### Recognition Rationale

- Participation evidence and number of active participants
- National or regional competitions organized
- Athlete development and coaching capacity
- Safety, rules, and standards documentation

### Governance and Support

- Proposed national governing structure
- Existing or planned membership associations
- Development plan and strategic objectives
- Supporting letters from sports associations or government authorities

### Fees and Tracking

- Application fee receipt
- `application_reference`
- `payment_reference`
- `status`, `submitted_at`, `updated_at`

## Workflow

- Save draft application details
- Submit for `GENERAL_SECRETARY` review and validation
- Review sport legitimacy, governance readiness, and development plan
- Approve, request clarifications, or reject

## API Endpoints

- `POST /api/v1/applications/national-sport-recognition`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
