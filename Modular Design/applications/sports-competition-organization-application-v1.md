# Sports Competition Organization Application — NCSMS v1.0

## Purpose

This application defines the requirements for organizing a formal sports competition or tournament.

## Key Sections

### Organizer Information

- Organizer name and entity type
- Contact person and event coordinator details
- Event location and venue details
- Sponsorship and funding information

### Competition Details

- Competition name
- Sports categories and participant age groups
- Start and end dates
- Expected number of teams or athletes
- Registration and eligibility rules

### Event Management

- Venue booking and logistics plan
- Health, safety, and medical support arrangements
- Results management and officiating plan
- Awards, medals, and sanctioning details

### Compliance and Fees

- Approval from relevant federation or association
- Insurance and liability documentation
- Risk assessment and security plan
- Payment receipt and `payment_reference`

## Tracking and References

- `application_reference`
- `competition_reference`
- `payment_reference`
- `review_status`
- `submitted_at`, `updated_at`

## Workflow

- Draft the event application and invite collaborators
- Submit for review by `GENERAL_SECRETARY`
- Validate sanctioning, venue compliance, and participant eligibility
- Approve, request clarifications, or reject submission

## API Endpoints

- `POST /api/v1/applications/competition-organization`
- `GET /api/v1/applications/{id}`
- `PATCH /api/v1/applications/{id}`
- `GET /api/v1/transactions/{reference}`
