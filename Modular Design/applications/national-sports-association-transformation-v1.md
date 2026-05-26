# National Sports Association Transformation Form — NCSMS v1.0

## Purpose

This form supports transformation requests from national sports associations seeking governance, structural, or scope changes.

## Key Sections

### Association Details

- Association name and registration number
- Current governance structure
- Member organizations and executive leadership

### Transformation Plan

- Description of the proposed transformation
- New governance model or organizational structure
- Strategic rationale and business case
- Timeline for implementation and stakeholder engagement

### Compliance Documentation

- Proposed constitution or bylaws amendments
- Board resolutions or membership approvals
- Legal opinions, regulatory clearances, or ministry correspondence
- Supporting documentation for change of scope

### Supporting Evidence

- Historical association performance and compliance record
- Stakeholder consultation records
- Payment receipt and application fee reference

## Tracking and References

- `application_reference`
- `transformation_reference`
- `payment_reference`
- `status`, `submitted_at`, `updated_at`

## Workflow

- Save a draft transformation request
- Submit under `GENERAL_SECRETARY` review
- Validate governance changes, membership approvals, and legal compliance
- Approve, request additional evidence, or reject

## API Endpoints

- `POST /api/v1/applications/transformation/association`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
