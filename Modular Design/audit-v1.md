# Audit Module — NCSMS v1.0

## Purpose

- Provide immutable audit tracking and monitoring for governance
- Support forensic review and compliance reporting
- Record security-sensitive events across modules

## Domain Responsibilities

- Capture audit events for authentication, authorization, data changes, and security incidents
- Support audit query and filtering for privileged roles
- Ensure audit records are write-only and tamper-resistant
- Integrate with RLS to expose audit logs only to authorized users

## High-level Code Layout

- `/internal/domain/audit/models.go`
  - `AuditEvent`, `AuditQueryFilter`, `AuditSourceMetadata`, `AuditOutcome`
- `/internal/domain/audit/repository.go`
  - `InsertAuditEvent(ctx, event)`, `QueryAuditEvents(ctx, filter)`, `FindAuditEventByID(ctx, id)`
- `/internal/domain/audit/service.go`
  - `RecordAuditEvent(ctx, event)`, `SearchAuditEvents(ctx, filter)`
- `/internal/domain/audit/handler.go`
  - HTTP endpoints: `GET /api/v1/audit`, `GET /api/v1/audit/{id}`

## Access and Policy

- Only `SYSTEM_ADMIN` and `AUDITOR` roles can query audit logs
- Audit logs are immutable and should never be editable through application APIs
- Sensitive metadata is stored securely and only exposed in limited query contexts
- Audit queries are rate-limited and require strict authorization checks

## Example Routes

- `GET /api/v1/audit` -> search audit events
- `GET /api/v1/audit/{id}` -> fetch specific audit event details
