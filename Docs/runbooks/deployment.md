# Deployment Runbook

## Preconditions

1. CI is green for the exact commit.
2. A change-window database backup has a SHA-256 manifest.
3. The backup passed an isolated restore verification.
4. Staging ran the same image and migrations.
5. Rollback was rehearsed and the previous image remains available.

## Promotion

Record the current commit, container digests and health output. Deploy the exact reviewed commit. Run forward-only migrations with a five-second lock timeout. Never recreate PostgreSQL or change existing `.env` values. Verify health, readiness, public pages, authentication, USSD and a protected API.

## Shared VPS invariants

- Never replace host-level Nginx configuration or unrelated sites.
- Never remove volumes, run `down -v`, or recreate production PostgreSQL.
- Keep uploads and evidence on persistent volumes outside release directories.
