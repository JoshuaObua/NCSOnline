# Sports Facility Operation Application — NCSMS v1.0

## Purpose

This application captures the details needed to approve and license a sports facility for official operation.

## Key Sections

### Facility Information

- Facility name
- Physical address and location details
- Facility owner or operating entity
- Capacity and primary use cases
- Facility type (stadium, hall, training center, court, pitch)

### Safety and Compliance

- Building permit and zoning approval
- Fire safety certification and emergency evacuation plans
- Accessibility and disability access documentation
- Insurance certificates and liability coverage

### Operations and Maintenance

- Maintenance schedule and checklist
- Operating hours and event booking policy
- Staffing and security arrangements
- Environmental controls and utilities management

### Supporting Documents

- Facility floor plans and seating capacity
- Lease or ownership documents
- Safety inspection reports
- Evidence of payment for application fees

## Tracking and References

- `application_reference`
- `facility_license_reference`
- `payment_reference`
- `review_status`
- `submitted_at` and `updated_at`

## Workflow

- Save as draft and complete facility information
- Submit to `GENERAL_SECRETARY` for operational review
- Verify safety, governance, and payment details
- Approve, reject, or request corrections

## API Endpoints

- `POST /api/v1/applications/facility-operation`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
