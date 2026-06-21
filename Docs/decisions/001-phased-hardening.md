# ADR-001: Deliver production hardening through gated phases

- Status: Accepted
- Date: 2026-06-21
- Owners: Engineering, Operations

## Context

NCS-Online is deployed on a shared VPS with persistent user data and existing Nginx sites. The programme changes traffic filtering, audit persistence, backups, sessions and deployments. A big-bang release would combine too many failure modes.

## Decision

Deliver through the ten phases in `Docs/System-Updates-Workplan.md`. Observability and verified backups precede traffic enforcement or automated upgrades. Production requires a staging rehearsal, tested rollback, current verified backup and recorded evidence.

## Alternatives considered

- Big-bang implementation: rejected because recovery and fault isolation would be poor.
- Production-only edits: rejected because restore and routing changes need isolation.
- Platform rewrite: rejected because it adds migration risk.

## Consequences

Delivery takes longer, while each change becomes measurable and reversible. Features may land disabled until operational dependencies are approved.

## Verification

CI gates commits; each phase records staging health, migration, backup and rollback evidence.
