# NSMIS Implementation Status

Last updated: 18 June 2026

## Delivered in the first implementation increment

- Additive NSMIS migrations `016` and `017` covering all 11 requested domains.
- New NCS and federation roles with seeded resource/action permissions.
- Federation membership scope, reporting periods, obligations, report revisions, workflow transitions, governance responses, compliance rule versions and document metadata.
- Athlete, affiliation, competition, result, medal, coach, technical official, talent, scholarship, safeguarding, disbursement, financial report and equipment schemas.
- Authenticated NSMIS API routes for federation registry, reporting periods, obligation generation/listing and governance report drafts/transitions.
- Federation-scoped repository queries and separation-of-duties protection for President approval.
- Governance dashboard with period selection, official-score status, missing reports and expired constitutions.
- Athlete, performance, finance and talent dashboard APIs and responsive frontend views.
- Role-aware NSMIS navigation and reporting workspace.
- Backend role/date tests, full `go test ./...`, `go vet ./...`, and frontend production build validation.

## Deployment prerequisite

Apply migrations in order through the updated `make migrate` target before starting the upgraded API. The migrations could not be executed locally during this increment because Docker Desktop/PostgreSQL was not running; they must be exercised against staging and backed up before production.

## Next implementation increments

1. Private object storage, malware scanning and signed evidence downloads.
2. Federation profile/officer/membership administration screens and APIs.
3. CRUD/import/review workflows for athlete, competition, medal, coach and technical official records.
4. Finance, equipment, talent and restricted safeguarding workflows.
5. Compliance calculation worker, reminders, email delivery, escalation and reconciliation jobs.
6. MFA/secure refresh-cookie migration, permission middleware and field-level PII controls.
7. Board-report/export generation, full audit events, OpenAPI contract and operational monitoring.
8. PostgreSQL integration tests, Vue component/E2E tests, accessibility testing, load testing and penetration testing.

The authoritative scope and acceptance standards remain in `DASHBOARD-UPGRADE-WORKPLAN.md`.
