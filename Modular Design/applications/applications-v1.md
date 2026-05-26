# Applications Domain — NCSMS v1.0

## Overview

This folder contains the individual application form definitions and workflow guidance for user-facing registration, renewal, license, and transformation requests.

The forms are designed for website signups and signed-in users to submit structured applications with proper tracking references, payment monitoring, and reviewer handoff to `GENERAL_SECRETARY`.

## Included Application Forms

- `federation-registration-renewal-v1.md` — federation registration and renewal applications.
- `sports-academy-operation-application-v1.md` — application to operate a sports academy.
- `community-academy-registration-renewal-v1.md` — application to register or renew community sports academies.
- `sports-facility-operation-application-v1.md` — application to operate a sports facility.
- `sports-competition-organization-application-v1.md` — application to organize a sports competition.
- `national-sports-association-transformation-v1.md` — transformation form for national sports associations.
- `national-sports-federation-transformation-v1.md` — transformation form for national sports federations.
- `national-sport-recognition-application-v1.md` — application for national sport recognition.

## Common Concepts

- `application_reference` — unique ID assigned to each submission.
- `payment_reference` — tracking number for fee payments and receipts.
- `status` — workflow state such as `DRAFT`, `SUBMITTED`, `UNDER_REVIEW`, `APPROVED`, `REJECTED`, or `NEEDS_INFORMATION`.
- `reviewer_id` — the `GENERAL_SECRETARY` or delegated reviewer processing the application.
- `submitted_at`, `last_updated`, and `next_action` metadata.

## Dashboard Integration

- The user dashboard exposes these forms as menu items in the `Applications` section.
- Applicants can save drafts, resume incomplete forms, view submitted applications, and inspect payment history.
- Status updates and reviewer notes are surfaced to the applicant through the same dashboard experience.

## Endpoint Naming

- `POST /api/v1/applications/<form-type>` — submit a new application.
- `GET /api/v1/applications` — list applications for the signed-in user.
- `GET /api/v1/applications/{id}` — fetch application details and tracking metadata.
- `PATCH /api/v1/applications/{id}` — update a draft or resubmit with new data.
- `GET /api/v1/transactions` — view payment and fee transaction history.
- `GET /api/v1/transactions/{reference}` — lookup a specific payment reference.
