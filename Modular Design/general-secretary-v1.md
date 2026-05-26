# General Secretary Module — NCSMS v1.0

## Purpose

- Define the review, approval, verification, and renewal workflows for licensing and application submissions.
- Provide the central processing role for federation registration, academy and facility applications, competition approvals, transformation requests, and national sport recognition.
- Monitor payment transactions and track fee references associated with each application.

## Domain Responsibilities

- Receive and triage incoming public applications from the website dashboard.
- Validate application completeness, supporting documentation, and payment confirmations.
- Approve, reject, or request clarifications for applications and renewals.
- Manage the lifecycle of application statuses with explicit tracking references.
- Monitor transaction logs, fee payments, and receipt references for each application.
- Coordinate with `SYSTEM_ADMIN` and other operational units on exceptional escalations.

## High-level Code Layout

- `/internal/domain/generalsecretary/models.go`
  - `Application`, `ApplicationReview`, `TransactionLog`, `VerificationRecord`, `RenewalRequest`
- `/internal/domain/generalsecretary/repository.go`
  - `CreateApplication(ctx, app)`, `UpdateApplicationStatus(ctx, id, status)`, `GetPendingApplications(ctx, params)`
  - `CreateTransactionLog(ctx, tx)`, `FindTransactionByReference(ctx, reference)`
- `/internal/domain/generalsecretary/service.go`
  - `SubmitApplication(ctx, req)`, `ReviewApplication(ctx, applicationID, action, notes)`
  - `VerifyPayment(ctx, applicationID, paymentReference)`, `CaptureTransactionReference(ctx, tx)`
- `/internal/domain/generalsecretary/handler.go`
  - HTTP endpoints: `GET /api/v1/general-secretary/applications`, `GET /api/v1/general-secretary/applications/{id}`
  - `PATCH /api/v1/general-secretary/applications/{id}/approve`, `PATCH /api/v1/general-secretary/applications/{id}/reject`
  - `GET /api/v1/general-secretary/transactions`, `GET /api/v1/general-secretary/transactions/{reference}`

## API Surface

- `GET /api/v1/general-secretary/applications` — review panel for pending and in-flight applications.
- `GET /api/v1/general-secretary/applications/{id}` — fetch full application and payment history.
- `PATCH /api/v1/general-secretary/applications/{id}/approve` — approve an application after validation.
- `PATCH /api/v1/general-secretary/applications/{id}/reject` — reject an application with review notes.
- `PATCH /api/v1/general-secretary/applications/{id}/clarify` — request additional applicant information.
- `POST /api/v1/general-secretary/applications/{id}/verify-payment` — verify transaction and payment status.
- `GET /api/v1/general-secretary/transactions` — search transaction logs by user, application, or reference.
- `GET /api/v1/general-secretary/transactions/{reference}` — fetch transaction details by payment reference.

## Access Control

- Primary role: `GENERAL_SECRETARY`.
- `GENERAL_SECRETARY` may `VIEW`, `UPDATE`, `EXECUTE`, and `APPROVE` applications and transaction logs.
- `SYSTEM_ADMIN` may override or audit any application decision.
- `AUDITOR` may view transaction logs and application review history but not modify outcomes.
- `CUSTOM_ROLE` may be scoped to application review or transaction monitoring as required.

## Transaction and Payment Monitoring

- Each application is linked to a unique `application_reference` and optionally a `payment_reference`.
- Transaction logs record: `transaction_reference`, `application_id`, `user_id`, `amount`, `payment_method`, `status`, `processed_at`, and `reviewer_id`.
- Payment receipts are stored or referenced securely and should be available to applicants and reviewers.
- `GENERAL_SECRETARY` workflows include payment verification, fee reconciliation, and discrepancy escalation.

## Review Workflows

- `DRAFT` — applicant is completing the form.
- `SUBMITTED` — application is ready for review and assigned a tracking reference.
- `UNDER_REVIEW` — `GENERAL_SECRETARY` or delegated reviewer is processing the application.
- `APPROVED` — application is accepted and relevant license or recognition is issued.
- `REJECTED` — application is denied with reasons and next steps.
- `NEEDS_INFORMATION` — additional documents or clarifications are requested.
- `RENEWAL_PENDING` — renewal applications awaiting decision.

## Security and Audit

- All application actions are audited with actor, timestamp, and reason.
- Transactions are validated against payment receipts and external payment carriers where applicable.
- Sensitive documents and payment metadata are stored encrypted and accessible only to authorized review roles.
- The system preserves immutable review history for compliance and dispute resolution.
