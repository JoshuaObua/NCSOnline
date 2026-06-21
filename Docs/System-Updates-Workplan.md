# NCS-Online — System Updates Implementation Workplan

**Owner:** Joshua Konshens
**Source spec:** [System Updates Plan.md](../System%20Updates%20Plan.md)
**Status:** Draft v1 — awaiting kickoff
**Last updated:** 2026-06-21

> Execution status: active. This document is the delivery authority for `System Updates Plan.md`. Production rollout is gated: code completion does not waive backup, staging, rollback, security-review, or soak requirements.

### Current-State Baseline

| Capability | Current evidence | Implementation decision |
| --- | --- | --- |
| Go API, PostgreSQL, Vue and Compose | Present in repository | Extend in place; avoid framework replacement. |
| Audit request capture | Middleware and migrations 005/009/021 | Replace per-request goroutines with a bounded writer; add chaining and immutability. |
| Geo/location lookup | `services/location_service/` | Move to GeoLite2-first operation; provider fallback must be explicitly enabled. |
| Session revocation | Per-user `auth_invalid_before` exists | Add global revocation and active-session management. |
| Background jobs and email | Job tables and SMTP worker exist | Add operator controls, delivery channels, preferences and DLQ operations. |
| `/my-portal` | Applications and organisation contexts exist | Add portal shell, security, sessions, audit, preferences and data controls. |
| CMS | Content modules and homepage settings exist | Add RBAC navigation, branding, SEO overrides and sitemap. |

### Environment and Promotion Matrix

| Environment | Purpose | Database | Promotion gate |
| --- | --- | --- | --- |
| Local | Unit/build verification | Disposable | Tests and static checks pass. |
| Staging | Integration, rollback and restore drills | Separate persistent volume | Acceptance evidence and rollback drill. |
| Production | Public service | Existing persistent volume, never recreated | Staging sign-off, verified backup and approved window. |

### Responsibility Matrix

| Work | Responsible | Accountable | Verification |
| --- | --- | --- | --- |
| Application implementation | Lead engineer | Project owner | CI and requirement checklist |
| VPS/Nginx promotion | Lead engineer | Project owner | Health/readiness and route smoke tests |
| Backup/restore evidence | Lead engineer | Project owner | Checksum and isolated restore report |
| Security review | Security reviewer | Project owner | No unresolved high-severity findings |
| Production sign-off | Project owner | Project owner | Approval after staging soak |

### External Prerequisites

- GeoLite2 download credentials and approved database update schedule.
- Off-site backup target and credentials.
- Telegram, Firebase, Africa's Talking, hCaptcha and SMTP/DNS credentials.
- GitHub branch protection, package permissions and image-signing identity.
- A staging hostname or approved port allocation.
- Named alert recipients, maintenance window and RPO/RTO owner.

Missing credentials disable only their feature through explicit flags; insecure defaults are forbidden.

---

## 1. Executive Summary

This workplan translates the System Updates Plan into a sequenced, week-by-week delivery schedule for the NCS-Online VPS platform. The work is grouped into **8 phases over ~14 weeks**, ordered so that each phase unblocks the next and so production risk decreases monotonically — observability and backups land before any hot-swap deploy logic ships.

The plan assumes a **single-engineer cadence** (~25 focused hours/week) on a live VPS at `104.219.248.160`. If a second engineer joins, parallelizable tracks are marked **⫶**.

### Headline Deliverables
1. Observability and backup safety nets (Phase 1) — *non-negotiable prerequisite.*
2. Edge geo-sentinel + hardened audit middleware (Phase 2).
3. Operator "flight deck": maintenance mode, log exports, cache tools (Phase 3).
4. Backup + restore UI with verified drills (Phase 4).
5. CI/CD pipeline + smart upgrade with atomic rollback (Phase 5).
6. Multi-channel notifications + device token registry (Phase 6).
7. Admin CMS reorg + branding + SEO + sitemap (Phase 7).
8. `/my-portal` user dashboard + 2FA + audit views (Phase 8).

---

## 2. Guiding Principles

| Principle                       | What it means in this plan                                                  |
| ------------------------------- | --------------------------------------------------------------------------- |
| **Safety nets first**           | Backups + observability ship before any feature that mutates prod state.    |
| **Reversible in <60s**          | Every phase has a documented rollback in its acceptance criteria.            |
| **No 3 AM panics**              | Staging-first, drill before prod, prefer automation over runbooks.           |
| **Decouple ops from request path** | Worker pools, atomic flags, async logging — zero added latency.          |
| **Operator > Developer**        | If the operator needs SSH for a routine task, the phase is incomplete.      |

---

## 3. Phase Roadmap (Gantt Overview)

```
Week:        1   2   3   4   5   6   7   8   9   10  11  12  13  14
Phase 0:     ██──────────────────────────────────────────────────────  Foundations
Phase 1:     ────████──────────────────────────────────────────────── Observability + Backups
Phase 2:     ────────████──────────────────────────────────────────── Edge + Audit Hardening
Phase 3:     ────────────████──────────────────────────────────────── Maintenance Core
Phase 4:     ────────────────████──────────────────────────────────── Backup/Restore UI
Phase 5:     ────────────────────██████──────────────────────────────  CI/CD + Smart Upgrade
Phase 6:     ────────────────────────────████────────────────────────  Notifications
Phase 7:     ────────────────────────────────████────────────────────  CMS Reorg + Branding
Phase 8:     ────────────────────────────────────████████──────────── User Portal
Phase 9:     ────────────────────────────────────────────────████──── Hardening + Launch
```

Phase boundaries are **gated** — the next phase does not start until the prior phase's acceptance criteria are met.

---

## 4. Phase Detail

### Phase 0 — Foundations *(Week 1, ~25h)*

Establish the scaffolding every subsequent phase depends on. Nothing user-visible ships.

**Tasks**
- [ ] Stand up a **staging VPS** (or staging Docker Compose stack on the same host bound to a different port + subdomain).
- [ ] Set up a `staging` branch + GitHub Actions workflow stub (no tests yet — just build + lint).
- [ ] Install `golangci-lint`, `govulncheck`, `gitleaks` locally + as pre-commit hooks.
- [ ] Create `Docs/` directory structure: `runbooks/`, `decisions/`, `diagrams/`.
- [ ] Create the **ADR template** ([decisions/000-template.md](decisions/000-template.md)) and record the first ADR: "Why phased over big-bang."
- [ ] Inventory current secrets — move any hardcoded values out of source into `.env`. Document rotation plan.
- [ ] Define **canonical event taxonomy** in [Docs/events.md](events.md) (e.g. `user.auth.login`, `admin.maintenance.toggle`, `system.backup.created`).

**Acceptance Criteria**
- Staging URL reachable, builds the same image as prod.
- Pre-commit hooks block commits with secrets or lint errors.
- ADR-001 merged.

**Risks**
- Staging cost / VPS resource pressure → mitigation: bind staging to off-port, single-instance Postgres.

---

### Phase 1 — Observability & Baseline Backups *(Weeks 2–3, ~50h)*

**Cannot be skipped.** Future phases will modify prod; we need to see what's happening and be able to recover.

**Tasks**
- [ ] Migrate all `log.Printf` to `log/slog` JSON output → stdout + rotating file in [backend/internal/logging/](../backend/internal/logging/).
- [ ] Add `/healthz` (liveness) and `/readyz` (DB + cache reachable) endpoints in [backend/cmd/server/main.go](../backend/cmd/server/main.go).
- [ ] Wire **graceful shutdown** — SIGTERM → drain in-flight (30s timeout) → close pools.
- [ ] Add **Prometheus metrics** at `/metrics`: request counters by route, latency histograms, goroutine count, DB pool stats.
- [ ] Stand up self-hosted Grafana + Prometheus in `docker-compose.yml` (internal network only, behind admin auth).
- [ ] Check in baseline dashboard JSON to [Docs/grafana/](grafana/).
- [ ] Write a `pg_dump` cron script → gzip → `/var/backups/ncs/`, retention 7 daily.
- [ ] Set up **off-site sync** via `rclone` to a Backblaze B2 bucket (cheap, S3-compatible).
- [ ] Configure **external uptime monitor** (UptimeRobot free tier) pinging `/healthz` every 5 min.

**Acceptance Criteria**
- `/healthz` and `/readyz` documented and pinged by uptime monitor.
- Grafana shows live request rate + latency p50/p95/p99.
- Manual `pg_dump` restore into a throwaway DB completes successfully (recorded in [runbooks/restore-drill.md](runbooks/restore-drill.md)).
- Off-site backup verified by downloading from B2 and restoring.

**Rollback Plan**
- Each new endpoint behind a feature flag; metrics endpoint can be unmounted without affecting business logic.

**Dependencies**: Phase 0 complete.

---

### Phase 2 — Edge Sentinel & Audit Middleware *(Week 4, ~25h)* ⫶

Two independent tracks; can be parallelized if a second engineer joins.

**Track A — Python Geo-Sentinel**
- [ ] Scaffold FastAPI service at [services/geo-sentinel/](../services/geo-sentinel/).
- [ ] Embed MaxMind GeoLite2 DB + weekly auto-update cron.
- [ ] Implement header-aware IP resolver (CF-Connecting-IP → XFF → RemoteAddr).
- [ ] Whitelist editor (config file initially; UI in Phase 3).
- [ ] Action matrix: 301 for browsers, 403 for APIs, USSD `END` for telco traffic.
- [ ] Add to `docker-compose.yml` upstream of the Go service in nginx routing.
- [ ] Metrics endpoint: dropped connections by country/hour → Grafana panel.

**Track B — Go Audit Middleware**
- [ ] Build [backend/internal/middleware/audit.go](../backend/internal/middleware/audit.go) — async worker pool, batched inserts every 250ms.
- [ ] Create migration `022_audit_log.sql` — append-only table, hash-chained.
- [ ] DB role `app_writer` granted INSERT only on `audit_log` (no UPDATE/DELETE).
- [ ] Add structural fingerprinting rules (mobile UA without SDK header, browser hitting `/api/admin/*` without referrer).
- [ ] Honeypot 307 redirect for flagged scanner traffic.
- [ ] PII scrub allowlist for payload capture (≤2KB).

**Acceptance Criteria**
- Non-UG curl from a known foreign IP returns expected per-channel response.
- Audit log hash chain validates via test script.
- 1000 RPS load test shows <2ms p99 added latency from audit middleware.
- Whitelist bypass works for GitHub webhook IP range.

**Dependencies**: Phase 1 (need logs + metrics to verify).

---

### Phase 3 — Maintenance Core *(Week 5, ~25h)*

Operator's "flight deck" — first piece of the Admin UI ships here.

**Tasks**
- [ ] `atomic.Bool` maintenance flag + middleware in [backend/internal/maintenance/](../backend/internal/maintenance/).
- [ ] Per-channel response (HTML / JSON / USSD `END`) with admin path bypass.
- [ ] Admin UI page: **Maintenance → System Status & Mode** — toggle with confirm modal.
- [ ] Live resource monitor: `gopsutil` → WebSocket → admin dashboard (CPU, RAM, disk, goroutines, connections).
- [ ] Cache flush button.
- [ ] **Force logout all** — bumps global `session_epoch` invalidating prior JWTs.
- [ ] Per-user session revocation (list active sessions table).
- [ ] Queue worker viewer (depends on having a queue — stub for now, real in Phase 6).

**Acceptance Criteria**
- Toggling maintenance mode from staging UI blocks public traffic within 1 request cycle.
- Admin user can still reach the admin UI when maintenance is on.
- Force-logout invalidates a live session within 5 seconds across all devices.

**Rollback Plan**
- Maintenance toggle is itself protected by a "kill switch" env var that disables the middleware entirely.

**Dependencies**: Phase 2 (RBAC must be solid).

---

### Phase 4 — Backup & Restore UI *(Week 6, ~25h)*

Promote the Phase 1 cron backup into a first-class operator tool.

**Tasks**
- [ ] Admin UI **Maintenance → DB Backup & Restore** page.
- [ ] Table of backups (local + B2): size, age, SHA-256 checksum, verify status.
- [ ] **Backup Now** button → triggers `pg_dump` job → progress streamed via WebSocket.
- [ ] **Verify** button → restores into transient `ncs_verify_<timestamp>` DB, runs row-count diff, drops the verify DB.
- [ ] **Restore** button → forces maintenance ON → safety snapshot → restore → cache clear → maintenance OFF.
- [ ] Weekly automated verify-only drill with email report.
- [ ] Retention enforcement: 7 daily, 4 weekly, 12 monthly.

**Acceptance Criteria**
- Operator restores a backup end-to-end from the UI on staging with no SSH.
- Weekly drill runs unattended for 2 consecutive weeks before phase considered done.
- Restore-while-traffic test shows no orphan connections after maintenance OFF.

**Dependencies**: Phase 3 (uses maintenance toggle).

---

### Phase 5 — CI/CD + Smart Upgrade *(Weeks 7–8, ~50h)*

Highest-risk phase. Practice everything on staging until boring.

**Track A — GitHub Actions CI**
- [ ] `.github/workflows/ci.yml` — lint, vet, govulncheck, `go test ./...`, frontend build.
- [ ] PG service container + ephemeral test DB for integration tests.
- [ ] Build container image, sign with **cosign**, generate SBOM with **syft**.
- [ ] Branch protection on `master`: require green CI to merge.

**Track B — VPS Deploy Key Provisioning**
- [ ] Admin UI **Maintenance → Smart Updates** → "Generate Deploy Key" button.
- [ ] Go calls `ssh-keygen -t ed25519`; private key → `chmod 600` at `/root/.ssh/id_ed25519_ncs`.
- [ ] Public key shown read-only with copy + step-by-step GitHub paste instructions.
- [ ] Rotate key flow (forces re-paste).

**Track C — Smart Upgrade Workflow**
- [ ] Update sentinel (GitHub releases polling, 5-min ticker).
- [ ] Pre-flight: green CI + disk space + DB reachable.
- [ ] `git fetch` → verify tag matches signed release → snapshot binary → `git pull` → `go build` → smoke test → atomic swap → `/readyz` probe within 30s.
- [ ] Failure path: restore `app.bak`, restart, stay in Maintenance, Telegram alert.
- [ ] Live deploy console (WebSocket stream of stdout/stderr).
- [ ] Migration safety: forward-only, `SET lock_timeout = '5s'`.
- [ ] **CI test** asserts deploy script never touches `/var/www/ncs/uploads/`.

**Acceptance Criteria**
- 3 successful staging deploys end-to-end with no manual intervention.
- 2 deliberate failures (broken binary, failing migration) recover automatically.
- First prod deploy via UI completes with <2 min total downtime.

**Rollback Plan**
- Old systemd unit kept hot for the first 4 prod deploys; manual `systemctl` documented in [runbooks/emergency-rollback.md](runbooks/emergency-rollback.md).

**Dependencies**: Phases 1, 3, 4 all complete.

---

### Phase 6 — Multi-Channel Notifications *(Week 9, ~25h)*

**Tasks**
- [ ] Migration: `user_device_tokens` table.
- [ ] Token registration endpoint + logout cleanup hook.
- [ ] Worker-pool dispatcher: buffered channel + N goroutines.
- [ ] **Telegram** channel (operator alerts) — extends Phase 5 sentinel.
- [ ] **Email** channel (SMTP + TLS) — extends existing [email_service.go](../backend/internal/services/email_service.go).
- [ ] **FCM** channel (mobile + web push) — Firebase Admin SDK.
- [ ] PWA service worker at `/my-portal/sw.js` with push + click handlers.
- [ ] Sanitization unit test: payload never contains paths, DSNs, env vars, internal IPs.
- [ ] Token hygiene cron: purge tokens with `last_seen_at > 90d` or FCM `Unregistered` response.
- [ ] Notification preferences API (consumed in Phase 8 by `/my-portal/notifications`).

**Acceptance Criteria**
- Backup success/failure events reach Telegram + email reliably across a 7-day soak.
- FCM token lifecycle (register → push received → logout → token purged) verified.
- Sanitization test passes for all event types.

**Dependencies**: Phases 1, 4, 5 (events to notify about).

---

### Phase 7 — Admin CMS Reorg + Branding + SEO *(Week 10, ~25h)*

**Tasks**
- [ ] Refactor [frontend/src/components/layout/AppSidebar.vue](../frontend/src/components/layout/AppSidebar.vue) into collapsible accordion groups.
- [ ] Backend serves nav tree JSON, RBAC-filtered (`/api/admin/navigation`).
- [ ] Active-route auto-expand + chevron rotation animation.
- [ ] **Branding Controls** page: 5 logo slots + favicon, MIME-validated upload, auto-resize, memory-cached resolver.
- [ ] **SEO Dashboard**: global meta + OG + JSON-LD with live character counters and validator.
- [ ] Per-route SEO override table.
- [ ] **Sitemap engine**: cron-driven goroutine walks route registry + blog/career/page tables → `sitemap.xml` served at `/sitemap.xml`.
- [ ] **Regenerate Now** button + last-generated card.

**Acceptance Criteria**
- Editor role sees zero maintenance/security nav items.
- Logo swap reflects on public site within 1 page refresh (cache invalidates correctly).
- `sitemap.xml` validates against Google's sitemap validator.

**Dependencies**: Phase 2 (RBAC), Phase 3 (admin UI shell).

---

### Phase 8 — User Portal `/my-portal` *(Weeks 11–12, ~50h)*

**Tasks**
- [ ] Layout shell with dynamic sidebar (Active Services, Available Apps, Account & Management).
- [ ] **Profile** page — avatar upload+crop, editable fields, verification flow for email.
- [ ] **Security** page:
  - [ ] Password change (regex enforced, current-password re-auth modal).
  - [ ] Transactional PIN (masked pin-pad component).
  - [ ] TOTP 2FA with QR + backup codes.
  - [ ] Email/SMS fallback 2FA.
  - [ ] Active sessions list + "sign out other devices".
- [ ] **Activity & Audit Logs** — paginated table from Phase 2 audit log, failed-auth rows highlighted.
- [ ] **Notification Preferences** — channel × category matrix; security alerts forced-on.
- [ ] CSRF double-submit on all mutating forms.
- [ ] Session cookie rotation on credential change.

**Acceptance Criteria**
- Full account lifecycle (register → enable 2FA → change password → audit log shows it → log out other sessions) verified manually on staging.
- WCAG 2.1 AA axe-core scan: zero critical violations.
- Mobile responsive: tested on iOS Safari, Android Chrome at 360px width.

**Dependencies**: Phase 2 (audit log), Phase 6 (notification prefs API).

---

### Phase 9 — Hardening, Localization & Launch *(Weeks 13–14, ~50h)*

**Tasks**
- [ ] Security headers middleware: HSTS, CSP (nonce-based), X-Frame-Options, Permissions-Policy.
- [ ] Per-IP + per-user rate limiting (stricter on `/auth/*`, `/api/ussd/*`).
- [ ] Account lockout with exponential backoff.
- [ ] hCaptcha on registration + password reset.
- [ ] i18n scaffolding: English baseline + Luganda starter pack.
- [ ] Africa's Talking SMS integration for OTP delivery.
- [ ] USSD session store moved to Redis with TTL.
- [ ] PWA offline-first cache for `/my-portal` core views.
- [ ] **Disaster Recovery runbook** with explicit RPO (1h) / RTO (30min) targets.
- [ ] Quarterly DR tabletop scheduled.
- [ ] Final security review using `/security-review` skill.
- [ ] Cutover checklist + go/no-go meeting.

**Acceptance Criteria**
- `securityheaders.com` score: A+.
- Rate limiter blocks a brute-force script within 10 attempts.
- DR drill: restore from 24h-old backup into a fresh VPS in <30 min.

---

## 5. Cross-Cutting Tracks

These run in the background throughout, not as discrete phases:

| Track                       | Cadence       | Notes                                          |
| --------------------------- | ------------- | ---------------------------------------------- |
| **Dependency updates**      | Weekly        | Dependabot PRs; auto-merge patch on green CI.  |
| **Backup restore drills**   | Weekly        | Automated from Phase 4.                        |
| **Audit log review**        | Weekly        | Spot-check anomaly events.                     |
| **Cost monitoring**         | Monthly       | VPS, B2, FCM quota, Africa's Talking SMS.      |
| **TLS renewal monitoring**  | Continuous    | Alert at 14d to expiry.                        |
| **Documentation**           | Each phase    | Runbook + ADR per phase, not at the end.       |

---

## 6. Risk Register

| #  | Risk                                                         | Likelihood | Impact | Mitigation                                                                                  |
| -- | ------------------------------------------------------------ | ---------- | ------ | ------------------------------------------------------------------------------------------- |
| R1 | VPS resource pressure from Grafana + Prometheus              | Med        | Med    | Bind to internal Docker network; cap retention to 14d; consider remote write later.         |
| R2 | Smart upgrade hot-swap corrupts a session mid-deploy         | Med        | High   | Drain SIGTERM with 30s timeout; sticky sessions optional; document for first 4 prod deploys.|
| R3 | Geo-sentinel blocks legitimate UG traffic via mobile carrier NAT egress in non-UG datacenter | Med | High | Telco gateway whitelist; monitor false-positive rate for first 2 weeks; emergency disable env var. |
| R4 | FCM key compromise leaks user push notifications             | Low        | High   | Key rotation quarterly; never log key; payload sanitization tests.                          |
| R5 | Audit log table grows unbounded                              | High       | Med    | Phase 1 retention policy: 90d hot, ship older to compressed cold storage.                   |
| R6 | Single-engineer bus factor                                   | High       | High   | Every phase produces a runbook; ADRs document "why"; staging mirror lets a successor learn safely. |
| R7 | Backup off-site provider outage (B2)                         | Low        | Med    | Secondary local copy always retained; manifest checksum verifies integrity.                 |
| R8 | Telegram bot token leak triggers spam to ops                 | Low        | Low    | Rate limit alerts; rotate token quarterly; secondary email channel as fallback.             |

---

## 7. Definition of Done (per phase)

A phase is **done** only when **all** of the following are true:

1. ✅ All acceptance criteria checked on staging.
2. ✅ Runbook merged under [Docs/runbooks/](runbooks/).
3. ✅ ADR merged for any non-obvious design choice.
4. ✅ Rollback plan executed once on staging (proves it works).
5. ✅ Monitoring dashboard updated to cover new components.
6. ✅ Code reviewed via `/review` or `/ultrareview`.
7. ✅ Security scan via `/security-review` returns no high-severity findings.
8. ✅ Operator (Joshua) signs off after 48h soak on staging before prod cutover.

---

## 8. Glossary

| Term       | Meaning                                                                |
| ---------- | ---------------------------------------------------------------------- |
| **RPO**    | Recovery Point Objective — max acceptable data loss (target: 1 hour).  |
| **RTO**    | Recovery Time Objective — max acceptable downtime (target: 30 min).    |
| **ADR**    | Architecture Decision Record — short markdown doc capturing a "why."   |
| **DLQ**    | Dead Letter Queue — failed jobs parked for manual review.              |
| **SBOM**   | Software Bill of Materials — inventory of dependencies in a build.     |
| **WCAG**   | Web Content Accessibility Guidelines — target AA conformance.          |
| **FCM**    | Firebase Cloud Messaging — push notification dispatch service.         |
| **TOTP**   | Time-based One-Time Password — RFC 6238 2FA standard.                  |

---

## 9. Sign-Off

| Role           | Name              | Date | Signature |
| -------------- | ----------------- | ---- | --------- |
| Project Owner  | Joshua Konshens   |      |           |
| Lead Engineer  | Joshua Konshens   |      |           |
| Security Review |                  |      |           |

---

*Generated 2026-06-21. Update this document at the end of each phase with actual completion dates, deviations, and learnings.*
