# NSMIS operations runbook

## Deployment

1. Back up PostgreSQL and the `private_data` volume.
2. Set production secrets in the deployment secret store; never in `.env` committed to Git.
3. Run `make migrate` and confirm migrations 016–018 complete without errors.
4. Deploy both `backend` and `nsmis-worker`; the worker is required for reminders, compliance and data-quality processing.
5. Verify `/health`, authenticate with a non-production test account, and exercise one record in each permission scope.
6. Queue `DATA_QUALITY_SCAN`, `CREDENTIAL_EXPIRY`, `DEADLINE_REMINDERS`, and one `COMPLIANCE_REFRESH` job from `POST /api/v1/nsmis/jobs`.

## Required configuration

- `DATABASE_URL`, `JWT_SECRET`, and strict `ALLOWED_ORIGINS`.
- `PRIVATE_STORAGE_PATH` on an encrypted, backed-up private volume. It must not be mounted into nginx.
- `MAX_EVIDENCE_BYTES` (default 15 MiB).
- SMTP/provider credentials are supplied through the environment used by the delivery adapter.

## Backups and recovery

- Nightly encrypted PostgreSQL logical backup; weekly restore drill in an isolated environment.
- Snapshot private evidence storage with retention matching database backups.
- Restore database and evidence from the same recovery point, then run SHA-256 reconciliation against `federation_documents`.
- Jobs left `RUNNING` are reclaimable after their five-minute lease expires.

## Monitoring and alerts

- Alert when API health fails, worker has no successful job for 15 minutes, dead-letter jobs exist, or notification failures exceed 5%.
- Review overdue obligations, open blocking data-quality issues, expired credentials and failed document access daily.
- Audit logs must be exported to append-only central storage and retained under the approved records schedule.

## Security incident response

1. Revoke affected accounts/sessions and rotate exposed credentials.
2. Preserve audit logs, report transitions and document hashes.
3. Isolate suspect evidence; set its scan state to `REJECTED` and prevent download.
4. Notify the data-protection and safeguarding owners where regulated data may be involved.
5. Record containment, scope, recovery and post-incident actions in the incident system.

## Rollback

Application containers are stateless and can be rolled back independently. Database migrations are additive; do not drop new tables during an incident rollback. Restore from backup only when data integrity is affected and business owners approve the recovery point.
