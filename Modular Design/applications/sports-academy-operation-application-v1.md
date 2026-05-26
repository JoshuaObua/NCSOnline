# Sports Academy Operation Application — NCSMS v1.0

## Purpose

This application collects the information required to approve and license a sports academy.

## Key Sections

### Applicant and Academy Details

- Applicant / organization name
- Academy name
- Academy address and facilities description
- Sports offered and program categories
- Number of athletes and coaches
- Staffing and coaching credential summaries

### Compliance and Safety

- Building safety certification
- Fire safety and emergency response plans
- Child protection and safeguarding policies
- Insurance coverage and liability documentation

### Operational Plan

- Training schedule and seasonal timetable
- Athlete admissions and age categories
- Coaching curriculum overview
- Community outreach and development goals

### Supporting Evidence

- Academy registration or operating license
- Facility lease or ownership documents
- Coach qualification certificates
- Evidence of local authority approval
- Payment receipt for application fee

## Tracking and References

- `application_reference`
- `payment_reference`
- `facility_registration_number`
- `status`, `submitted_at`, `updated_at`

## Workflow

- Save draft and complete detailed academy data
- Submit for `GENERAL_SECRETARY` review
- Validate facility, safety, and coaching documentation
- Issue approval, request additional information, or reject application

## API Endpoints

- `POST /api/v1/applications/academy-operation`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
