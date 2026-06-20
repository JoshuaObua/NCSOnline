# NSMIS Dashboard Upgrade Workplan

## 1. Document purpose

This workplan translates `Changes.md` into an implementable frontend, backend, data, security, testing, and operations programme for upgrading the current National Council of Sports Management System (NCSMS) into the National Sports Management Information System (NSMIS).

The plan is intentionally organised around maintainable domain modules, explicit ownership, reusable platform services, and phased delivery. It covers the dashboard and the reporting workflows that must supply trusted dashboard data. A dashboard-only visual redesign would not meet the requirement because most requested indicators do not yet have a source of record in the current implementation.

## 2. Executive outcome

The upgraded platform will allow authorised NCS staff and federation officers to:

- maintain one governed profile for each federation;
- submit, review, approve, return, and lock monthly, quarterly, annual, competition, and international-event reports;
- maintain national registers for athletes, coaches, technical officials, talent, medals, competitions, finance, safeguarding, and equipment;
- upload supporting evidence securely;
- calculate compliance and performance measures from versioned rules;
- receive deadline reminders and escalation notices;
- view role-appropriate governance, athlete, performance, finance, and talent dashboards;
- drill from every aggregate into the permitted source records;
- generate reproducible board reports and annual sports statistics; and
- audit who viewed or changed sensitive information.

## 3. Current-system assessment

### 3.1 Reusable foundation found in the repository

| Area | Current state | Reuse decision |
|---|---|---|
| Frontend | Vue 3, Vite, Pinia, Vue Router, Tailwind CSS | Retain; reorganise into domain features and shared dashboard components. |
| Backend | Go API using Chi, PostgreSQL repositories, handlers and services | Retain; introduce domain packages and dedicated query/reporting services. |
| Authentication | JWT access/refresh flow, account invalidation and PIN support | Harden and extend for MFA, secure cookies and scoped organisation access. |
| Authorisation | Seeded roles plus some route-level role checks | Replace coarse role checks with permission and federation-scope enforcement. |
| Audit | Audit-log schema and admin viewer exist | Extend to domain events, exports, sensitive reads and append-only retention. |
| Applications | Draft, submission, review and attachment concepts exist | Reuse workflow patterns where appropriate; keep statutory applications separate from recurring federation reports. |
| Dashboard | One Vue view and one API returning application/user counts | Replace with role-aware NSMIS dashboard endpoints and feature views. |
| Deployment | Docker Compose, Nginx, PostgreSQL and persistent upload volume | Harden for production; move sensitive files to private object storage. |
| Automated tests | A small middleware test set; no frontend test suite found | Establish unit, integration, component, E2E, security and load testing. |

### 3.2 Material gaps against `Changes.md`

- No production data model or API was found for federations, recurring reporting periods, governance indicators, athletes, competitions, medals, coaches, officials, talent, safeguarding, finance or equipment.
- The current dashboard API aggregates only users and statutory applications; it cannot produce the requested NSMIS indicators.
- The current NCS dashboard route is limited to `super_admin` and `admin`; the requested departmental and federation roles are not represented accurately.
- Current navigation is role-name driven and does not derive capabilities from server permissions.
- Existing files under `/uploads/` are served publicly by Nginx. This is unsuitable for identity documents, bank statements, receipts, safeguarding reports and other protected evidence.
- Refresh tokens are handled by the browser application; the target should use hardened server-managed cookies and CSRF protection.
- No charting package, analytics cache, scheduled-job service, notification provider, report generator, object-store integration or malware scanner is evident in the active application.
- Existing federation, athlete and asset Markdown files are architectural outlines, not implemented modules.
- Database migrations are mounted into PostgreSQL initialisation; a production-grade migration runner and rollback/forward-fix procedure are still required.

## 4. Scope, boundaries and delivery principles

### 4.1 In scope

- All 11 modules in `Changes.md`.
- Five automated NCS dashboards and role-specific federation workspaces.
- Reporting schedules, reminders, overdue flags, compliance scores and rankings.
- Board/annual report generation, controlled exports and source-data drill-down.
- Supporting shared services: organisations, reporting periods, workflow, documents, notifications, audit, analytics and reference data.
- Responsive and accessible frontend experiences.
- Security, privacy, observability, backup, recovery and maintainability.

### 4.2 Explicit boundaries

- Existing public CMS and statutory licensing/application workflows remain separate bounded contexts. They may share identity, documents, audit and notification services.
- “Real time” should mean committed data is visible within a defined freshness target, not that unapproved drafts affect official statistics.
- Rankings and compliance scores are decision-support tools. Their formulas, effective dates and overrides must be approved by NCS policy owners before production use.
- Safeguarding case details must not appear in general dashboards; only carefully thresholded aggregates may be shown outside the safeguarding function.

### 4.3 Engineering principles

- Use PostgreSQL as the system of record; derive dashboard summaries from approved/locked records.
- Keep business rules in backend services and versioned database configuration, never only in Vue components.
- Enforce least privilege and federation scope at API, service and query levels; UI hiding is not a security control.
- Treat uploaded evidence as private by default.
- Record immutable workflow and calculation history so any dashboard value can be explained later.
- Prefer modular-monolith boundaries initially; introduce additional services only where scheduling, document processing or scale justifies them.
- Use additive, reversible migrations and backward-compatible API changes during rollout.
- Every indicator must have an owner, definition, source, update frequency, quality rule and access classification.

## 5. Decisions and policy clarifications required before build

The product owner should resolve these during inception. Until then, they must be tracked as assumptions rather than silently embedded in code.

| Decision | Why it matters | Proposed default |
|---|---|---|
| Reporting calendar basis | Defines deadlines, overdue states and annual statistics. | Uganda/EAT calendar; configurable reporting year and period boundaries. |
| Monthly vs quarterly obligation per federation | Prevents duplicate or missing reports. | Assign one active cadence per report type and federation. |
| Annual deadline year | “31 July” needs the reporting year it closes. | Due 31 July after the applicable financial/reporting year, configurable. |
| Approval workflow | Determines when data becomes official. | Submitter -> Federation President approval -> relevant NCS department review -> locked. |
| Federation President authority | `Changes.md` says view/approve reports but not which reports. | Approve all federation-originated reports before NCS review. |
| NCS General Secretary authority | Reports/dashboard may imply read-only or board-report sign-off. | Read all official data and approve generated board reports; no silent source edits. |
| Technical/Finance separation | Needed for least privilege. | Technical manages sport/performance data; Finance manages financial reports; each has read-only cross-domain aggregates where approved. |
| Compliance formula and weights | Ranking is not meaningful without approved rules. | Store versioned weighted rules; show score breakdown and data completeness. |
| Performance ranking formula | Medals alone may bias rankings. | Separate medal table from policy-approved performance index. |
| “Expired constitution” | Requires issue/approval/expiry metadata not listed. | Capture effective date, review/expiry date and document status. |
| Funds released | Revenue form alone cannot prove NCS disbursement. | Add NCS-controlled grant/disbursement ledger and reconcile federation declarations. |
| Talent progression/scholarships | Requested dashboard fields are absent from Module 8 inputs. | Add progression events, national-team promotion and scholarship records. |
| International representation | Needs an exact event/competition rule. | Derive from competition level plus approved participation records. |
| Safeguarding data access | High-risk personal and case information. | Dedicated permission group, case pseudonyms and aggregate-only general reporting. |
| Athlete identifier | National ID cannot safely serve as the platform key. | Generate immutable NSMIS athlete ID; store identity-document values encrypted and separately permissioned. |
| Data retention | Documents have different legal/operational value. | Approve a record-class schedule before go-live; legal hold overrides deletion. |
| Public data | Not specified in the change request. | No new NSMIS data public by default; publish only separately approved aggregates. |

## 6. Target roles and access model

### 6.1 Role catalogue

| Role | Primary access | Important restrictions |
|---|---|---|
| NCS Administrator | System configuration, users, reference data, all modules, audit and support | Sensitive case/document access still requires explicit permission; no ability to erase audit history. |
| General Secretary NCS | Cross-federation official dashboards, reports, board packs and sign-off | Source edits only where separately granted. |
| Technical Department | Federation monitoring, athletes, competitions, medals, coaches, officials and talent | Financial evidence and safeguarding identities excluded by default. |
| Finance Department | Grants, revenue, expenditure, accountability and finance compliance | Athlete identity documents and safeguarding cases excluded. |
| Federation President | Own federation dashboard, review and approval queue | Cannot approve own authored submission without an approved exception. |
| Federation General Secretary | Own federation data entry, evidence upload and submission | Cannot see or change another federation; cannot unlock approved records. |
| Safeguarding Officer | Restricted safeguarding case intake and resolution | Dedicated appointment and access review; case content excluded from normal exports. |
| Auditor/Read-only Reviewer | Time-bound read/export access to assigned scope | No mutation; all sensitive views and exports audited. |

### 6.2 Permission design

Adopt resource/action permissions with an explicit organisation scope, for example:

- `federations:read:any`, `federations:read:own`, `federations:update:own`;
- `reports:create:own`, `reports:submit:own`, `reports:approve:own`, `reports:review:any`, `reports:lock:any`;
- `athletes:manage:own`, `athletes:read:any`, `athletes:pii:read`;
- `finance:submit:own`, `finance:review:any`, `finance:evidence:read`;
- `safeguarding:aggregate:read`, `safeguarding:case:read`, `safeguarding:case:manage`;
- `dashboard:governance:read`, `dashboard:finance:read`, `dashboard:performance:read`;
- `exports:create`, `board_reports:generate`, `board_reports:approve`; and
- `system:settings:manage`, `audit:read`, `audit:export`.

Backend policies must combine permission, federation membership, record state, reporting period and field classification. A user may hold multiple roles. Federation assignments must have start/end dates and be revocable without changing historical authorship.

## 7. Target architecture

### 7.1 Logical components

```text
Vue NSMIS workspace
  -> Go REST API / policy enforcement
      -> transactional domain services
          -> PostgreSQL source-of-record tables
          -> private object storage + malware scan
          -> append-only audit events
          -> outbox events
      -> reporting query service
          -> approved-data views/materialized summaries
          -> short-lived cache
  -> background worker
      -> reminders/escalations
      -> dashboard refresh
      -> document processing
      -> board/annual report generation
      -> notification delivery/retry
```

### 7.2 Recommended backend code layout

```text
backend/internal/
  platform/{auth,policy,audit,documents,notifications,jobs,reference}/
  domain/{federations,reports,athletes,competitions,medals,workforce,talent,
          safeguarding,finance,equipment,compliance}/
  reporting/{queries,indicators,exports}/
```

Each domain owns its DTOs, validation, model, repository interface, service, handler and tests. Cross-domain dashboard queries belong in `reporting`, not inside HTTP handlers. HTTP handlers should parse/validate transport data and call services; they should not contain scoring or approval rules.

### 7.3 Data consistency and analytics approach

- Save drafts transactionally but exclude them from official aggregates.
- On approval/lock, publish an outbox event in the same database transaction.
- A worker refreshes affected summary tables/materialized views idempotently.
- Provide `as_of`, `data_freshness`, `calculation_version` and active-filter metadata with dashboard responses.
- Start with PostgreSQL views/materialized views and indexed query tables. Consider a separate warehouse only after measured volume or BI needs justify it.
- Nightly reconciliation must compare summary totals to source records and alert on differences.

## 8. Core data model

All main tables should use UUIDs, `created_at`, `updated_at`, actor IDs where relevant, optimistic-lock version numbers, explicit status constraints and appropriate indexes. Prefer history/event tables over destructive updates for regulated records.

### 8.1 Shared entities

| Entity | Purpose and key fields |
|---|---|
| `federations` | Name, acronym, NCS registration number, recognition status, contacts, address, active dates and version. |
| `federation_memberships` | User, federation, appointed role, start/end dates and status. |
| `federation_officers` | Leadership position, person/contact data, appointment dates and active/history status. |
| `reporting_calendars` | Calendar name, timezone, reporting year and active policy version. |
| `reporting_periods` | Type, start/end, due date, grace date and lock date. |
| `report_obligations` | Federation, report type, cadence, period, due date, status and exemption reason. |
| `reports` | Federation, period, type, revision, workflow state, submitter, approvers, timestamps and lock state. |
| `report_reviews` | Reviewer, department, decision, reason, checklist and timestamp. |
| `documents` | Owner type/ID, classification, storage key, hash, MIME, size, scan state, version and retention class. |
| `reference_values` | Versioned districts, regions, disciplines, competition levels, funding sources, certifications and other controlled lists. |
| `notifications` | Template, recipient, channel, related record, state, attempts and delivery timestamps. |
| `outbox_events` / `job_runs` | Reliable background processing, idempotency and operational history. |
| `indicator_definitions` | Code, label, formula/config, unit, filters, access class, effective dates and version. |
| `indicator_snapshots` | Indicator, dimensions, value, source watermark, as-of date and calculation version. |

### 8.2 Module records

| Module | Principal records | Required additions/validation |
|---|---|---|
| Federation Profile | Profile, officers, profile documents | Unique registration number/acronym policy; contact validation; constitution effective/review dates; versioned officer history. |
| Governance Reports | Governance responses, meeting events, disciplinary counts | Conditional dates/evidence when “Yes”; non-negative counts; AGM/election frequency validation. |
| Athlete Database | Athletes, federation affiliations, team-status history, identity documents | Generated athlete ID; deduplication; encrypted identity number; DOB not future; district/club/discipline references; consent/legal basis. |
| Competition Reporting | Competitions, participation, entries/results, officials | Start/end/return dates; level; venue/country; result uniqueness; participant counts reconciled to entries. |
| Medal Tracking | Medal awards linked to athlete/team and competition | Medal enum; official-result evidence; uniqueness rule; corrections via revision rather than overwrite. |
| Coach Database | Coaches, federation affiliation, certification and licence history | Expiry validation and reminder thresholds; private certificates. |
| Technical Officials | Officials, types, certification and validity history | Multiple official types/levels; expiry state derived from dates. |
| Talent Identification | Talent record, assessments, progression events, scholarships | Add progression and scholarship fields needed by the dashboard; avoid duplicate athlete creation. |
| Safeguarding & Gender | Period aggregates and restricted cases | Separate aggregate metrics from case records; pseudonymous case ID; restricted resolution workflow. |
| Financial Accountability | NCS disbursements, declared revenue, expenditure, accountability and evidence | Use numeric/decimal currency fields; currency and period; balanced totals; reconciliation and approval; immutable financial revisions. |
| Equipment Management | Items, receipts, distributions, beneficiaries, acknowledgements and movement history | Controlled units/categories; received/distributed/on-hand reconciliation; no negative stock. |

### 8.3 Database integrity and performance

- Foreign keys must express all ownership and source relationships.
- Use partial unique indexes for active memberships, active officer appointments and one current report revision.
- Index common dashboard dimensions: federation, period, status, date, region, gender, discipline, competition level and medal type.
- Use database check constraints for enums, non-negative quantities/amounts and valid date ranges.
- Apply tenant/federation isolation through scoped queries and, where practical, PostgreSQL row-level security using transaction-local identity context.
- Partition only high-volume audit/notification/snapshot tables after measurement; do not prematurely partition transactional tables.
- Define anonymisation and purge jobs by retention class, with legal-hold support.

## 9. Workflow and automation design

### 9.1 Standard report state machine

```text
NOT_STARTED -> DRAFT -> SUBMITTED -> PRESIDENT_APPROVED -> NCS_UNDER_REVIEW
                                      |                    |-> NEEDS_CORRECTION -> DRAFT (new revision)
                                      |                    |-> APPROVED -> LOCKED
                                      |-> PRESIDENT_RETURNED
```

- All transitions must be server validated and idempotent.
- Every transition records actor, role, timestamp, previous/new state, reason and report revision.
- Locked records cannot be edited. Corrections create linked amendments with visible history.
- Separation of duties blocks self-approval unless a documented, audited emergency policy permits it.
- A report cannot be submitted until its server-generated completeness checklist passes.

### 9.2 Deadline engine

- Seed schedule rules from `Changes.md`: monthly by the 5th, quarterly by the 10th following the quarter, annual by 31 July, competition within 7 days after the event, and international event within 14 days after return.
- Store schedule rules as versioned configuration, not hard-coded date arithmetic.
- Generate obligations when a reporting period opens or an event is approved.
- Use EAT consistently and test month-end, quarter-end, leap-year, weekend/holiday and policy-change cases.
- Suggested reminder sequence: opening notice, 14/7/3/1 days before due, due date, then 1/3/7/14 days overdue. Make the sequence configurable.
- Escalate failed deliveries and persistent overdue obligations according to federation/NCS contacts.
- Provide a calendar and notification history so users can see what the system expected and sent.

### 9.3 Compliance engine

- Separate completeness, timeliness, governance, financial, document-validity and data-quality components.
- Define each rule with code, description, weight, effective dates, required inputs, scoring method and owner.
- Persist per-rule outcomes and explanations, not only a total score.
- Recalculate idempotently after relevant approvals or policy changes.
- Never rewrite historical scores when a formula changes; recalculate under a new version and retain both.
- Support authorised overrides only with a reason, expiry, two-person approval and audit event.
- Show “insufficient data” rather than treating missing denominators as zero.

### 9.4 Notifications and reports

- Use an outbox/worker pattern so transactions are not coupled to email-provider availability.
- Maintain versioned email templates, recipient preferences, delivery attempts, bounce/suppression status and retry/dead-letter handling.
- Generated PDF/CSV/XLSX reports must include filters, reporting period, generation timestamp, calculation version and confidentiality marking.
- Board packs should be generated from immutable snapshots and move through draft/approved/published states.
- Large exports should run asynchronously, expire automatically and require re-authorisation at download time.

## 10. Dashboard information architecture

### 10.1 Shared dashboard controls

All NCS dashboards should share:

- reporting period and as-of date;
- federation, region, discipline and competition-level filters where relevant;
- data freshness and official/draft status;
- comparison to previous period;
- accessible chart/table toggle;
- drill-down to permitted source records;
- saved views for authorised staff;
- controlled export; and
- clear definitions/tooltips for every measure.

Filters must be encoded in the URL so views are shareable among users with equivalent permissions. The API must revalidate every filter and scope; it must not trust hidden UI controls.

### 10.2 Governance dashboard

| Component | Definition/source | Drill-down |
|---|---|---|
| Compliant federations | Federations meeting the selected, versioned compliance threshold | Score breakdown by rule and period. |
| Non-compliant federations | Active federations below threshold or with blocking failures | Failed obligations and remediation state. |
| Expired constitutions | Constitution review/expiry date before as-of date | Document metadata and renewal history, subject to permission. |
| Missing/overdue reports | Open obligations past due without accepted submission | Federation, report type, due date and days overdue. |
| Governance trends | Approved governance responses over time | Period response and evidence status. |

### 10.3 Athlete dashboard

| Component | Definition/source | Notes |
|---|---|---|
| Total registered athletes | Active, non-duplicate athlete master records | Show verified/unverified separately. |
| Gender distribution | Active athletes grouped by approved gender categories | Display counts alongside ratios; handle unknown/not stated. |
| Athletes by region | Current residence/registered district mapped to region | Preserve historical geography for period comparisons. |
| Athletes by federation | Current active affiliation | Define handling of multiple affiliations. |
| National-team status | Effective-dated status history | Do not infer only from latest upload. |

### 10.4 Performance dashboard

| Component | Definition/source | Notes |
|---|---|---|
| Medal totals/table | Approved medal records by type | Standard order Gold, Silver, Bronze; expose correction history. |
| Medals by federation | Medal linked through athlete/team federation at event time | Preserve event-time affiliation. |
| Medals by country | Competition host country, not athlete nationality | Label explicitly to avoid ambiguity. |
| International representation | Approved participation in continental/international events | Show athletes, teams, events and period. |
| Performance ranking | Policy-approved composite or medal-table order | Never conflate with compliance ranking. |

### 10.5 Financial dashboard

| Component | Definition/source | Notes |
|---|---|---|
| Funds released | Approved NCS disbursement ledger | Reconcile against federation-declared government grant. |
| Accountability submitted | Approved accountability submissions for disbursements/period | Display amount and submission count. |
| Outstanding accountability | Released amount not accepted/accounted for by due date | Show ageing buckets and federation. |
| Revenue/expenditure mix | Approved category totals | Currency-safe; no floating-point arithmetic. |
| Reconciliation exceptions | Mismatches, missing evidence and rejected items | Finance-only details. |

### 10.6 Talent dashboard

| Component | Definition/source | Notes |
|---|---|---|
| Athletes identified | Valid talent records in selected period | Deduplicate against athlete register. |
| Progressed to national teams | Approved progression event to national-team status | Requires added progression model. |
| Receiving scholarships | Active scholarship records in selected period | Requires added scholarship model and privacy classification. |
| Pipeline conversion | Identified -> assessed -> selected -> progressed | Use cohort definitions and explain denominators. |

### 10.7 Role-specific landing pages

- NCS Administrator: system health, overdue workload, data-quality exceptions, security alerts and shortcuts; not every business chart at once.
- General Secretary: executive cross-federation summary, major exceptions, trends and board-report queue.
- Technical Department: athlete, competition, medal, coach, official and talent work queues/dashboards.
- Finance Department: disbursements, accountabilities, ageing, reconciliation exceptions and evidence-review queue.
- Federation President: own-federation scorecard, approval queue, upcoming deadlines and returned reports.
- Federation General Secretary: form completion, missing evidence, due/overdue reports, data-quality tasks and recent submissions.

## 11. Frontend workplan

### 11.1 Application structure

- Replace the single all-purpose dashboard view with route-level feature pages: `dashboard/executive`, `governance`, `athletes`, `performance`, `finance`, `talent`, and federation workspace routes.
- Add domain feature directories under `frontend/src/features/`, each containing API adapters, views, components, composables, schemas and tests.
- Introduce a shared dashboard shell, filter bar, KPI card, accessible chart wrapper, data table, freshness badge, definition popover, export action and error/empty/loading states.
- Move API calls out of views into typed domain clients. If TypeScript adoption is not immediate, use JSDoc types and runtime response validation as an intermediate step.
- Lazy-load role-specific routes and large chart/export dependencies.
- Build navigation from the authenticated capabilities returned by `/auth/me`, with stable server-defined permission codes.

### 11.2 Federation reporting experience

- Create a task-oriented home page: due soon, overdue, drafts, awaiting approval, returned and locked.
- Use reusable multi-section form patterns with autosave status, explicit save, server validation summary, per-field errors and unsaved-change protection.
- Allow CSV import for suitable athlete/coach/official records with downloadable template, preview, row validation, duplicate matching and partial-failure report.
- Provide document upload progress, accepted-format guidance, scan/processing state, replacement/version history and accessible failure recovery.
- Provide a review screen that exactly matches the submitted payload and evidence manifest before approval.
- Show workflow history and reviewer comments without exposing internal or restricted notes.

### 11.3 Visualisation and UX standards

- Select a maintained chart library only after a short accessibility and bundle-size proof of concept; every chart requires an equivalent data table.
- Do not communicate state using colour alone. Use labels, icons/patterns and suitable contrast.
- Support keyboard operation, visible focus, logical headings, screen-reader names, reduced motion and 200% zoom.
- Target WCAG 2.2 AA and responsive layouts at 320, 375, 768, 1024 and 1440+ px.
- Use Uganda locale defaults, EAT timestamps and explicit date formats; keep stored timestamps in UTC.
- Preserve filter state during drill-down/back navigation.
- Distinguish zero, no data, withheld data, not applicable and loading—these are not interchangeable.

### 11.4 Frontend performance and resilience

- Cancel stale requests when filters change and debounce search controls.
- Server-paginate large tables and virtualise only after measuring need.
- Cache reference data with version/ETag validation; do not persist sensitive dashboard payloads in browser storage.
- Use request IDs and user-friendly retry paths for partial widget failures.
- Set performance budgets for initial JS, route chunks, LCP and interaction latency on realistic Ugandan network/device profiles.
- Add a global error boundary and observability integration with PII redaction.

## 12. Backend and API workplan

### 12.1 Domain API design

- Publish an OpenAPI 3 contract before or alongside each module.
- Use versioned `/api/v1` resources, consistent pagination/filter/sort conventions and standard error envelopes.
- Create explicit request/response DTOs; never bind database models directly to public JSON.
- Validate enum values, date ranges, ownership, workflow state and conditional fields server-side.
- Use idempotency keys for report submission, approvals, imports, generated reports and notification-triggering commands.
- Apply optimistic concurrency (`version`/ETag) so two editors cannot silently overwrite each other.

Suggested resource groups:

```text
/federations, /federations/{id}/officers, /federations/{id}/documents
/reporting-periods, /report-obligations, /reports, /reports/{id}/transitions
/athletes, /competitions, /results, /medals
/coaches, /technical-officials, /talent, /scholarships
/safeguarding/aggregates, /safeguarding/cases
/finance/disbursements, /finance/accountabilities
/equipment/items, /equipment/receipts, /equipment/distributions
/dashboards/{governance|athletes|performance|finance|talent}
/exports, /board-reports, /notifications
```

### 12.2 Dashboard query contracts

- Provide one endpoint per bounded dashboard rather than one ever-growing response.
- Accept validated period/dimension filters and enforce maximum date ranges.
- Return indicator codes, values, units, comparison values, suppression flags, definition/calculation version and freshness metadata.
- Provide separate paginated drill-down endpoints.
- Cache only by complete permission scope and filter key; never share a broader cached response with a narrower user.
- Add query timeouts and cost limits. Run large exports asynchronously.
- Instrument query latency, cache hit rate, freshness lag and reconciliation failures.

### 12.3 Background jobs

- Implement a worker process sharing domain packages with the API but deployed independently.
- Use database-backed job leasing initially: `FOR UPDATE SKIP LOCKED`, heartbeat, attempt count, next-attempt time and dead-letter state.
- Jobs must be idempotent and safe after worker restarts.
- Required jobs: obligation generation, reminders/escalations, licence/certificate expiry, malware/document processing, imports, indicator refresh, reconciliation, exports, board reports, retention and audit integrity checks.
- Provide an admin job monitor with retry/cancel permissions and redacted error details.

### 12.4 Data-quality controls

- Run synchronous validation for correctness that can be decided during submission.
- Run asynchronous checks for duplicates, cross-report reconciliation and unusual values.
- Maintain a data-quality issue table with severity, owner, status, related records and resolution history.
- Prevent official publication when a blocking issue exists; permit documented non-blocking warnings.
- Add duplicate-detection workflows for athletes and people; merging must preserve aliases, links and audit history.

## 13. Security, privacy and governance workplan

### 13.1 Identity and session security

- Move refresh tokens to `HttpOnly`, `Secure`, appropriately scoped `SameSite` cookies; keep short-lived access tokens in memory.
- Add CSRF protection for cookie-authenticated mutations.
- Require MFA for NCS staff, safeguarding users and federation approvers; support recovery codes and audited reset.
- Apply password, lockout, session-revocation and inactive-account policies; avoid security questions.
- Replace unconditional VPN blocking with an approved access-risk policy; legitimate remote work and incident access need a controlled path.
- Review the impact of Uganda-only geo-blocking on international-event reporting and document authorised exceptions.

### 13.2 Authorisation and tenant isolation

- Replace route-only role groups with permission middleware and service-level resource policies.
- Scope every federation query using authenticated memberships, never a client-supplied federation ID alone.
- Add negative tests for cross-federation reads/writes, guessed IDs, exports and file downloads.
- Apply field-level response shaping for athlete PII, financial evidence and safeguarding information.
- Conduct quarterly access reviews and automatically expire temporary assignments.

### 13.3 Document security

- Move all NSMIS evidence to private S3-compatible object storage; prohibit direct public `/uploads/` access.
- Use random object keys, server-side encryption, TLS, bucket policies and short-lived signed download URLs.
- Verify file signatures, allowlisted types, size/page limits and declared MIME; rename files and reject archives/macros unless explicitly required.
- Quarantine new uploads and scan for malware before reviewers can open them.
- Store SHA-256 hash, uploader, classification, scan result, versions and chain of custody.
- Use `Content-Disposition: attachment`, `nosniff`, a sandboxed preview service where required and no inline execution of active content.
- Apply per-record permission checks at URL issuance and again where the storage architecture permits.

### 13.4 Data protection

- Complete a data-protection impact assessment before collecting national IDs, passport copies and safeguarding cases.
- Classify fields as public, internal, confidential, restricted or highly restricted.
- Encrypt identity numbers and other high-risk fields at application/column level with managed key rotation.
- Minimise collection: store only data required by an approved purpose and legal basis.
- Mask sensitive values in normal screens, logs, analytics, support tools and non-production environments.
- Define data-subject correction and lawful deletion/anonymisation processes without breaking statutory audit records.
- Apply small-number suppression to safeguarding and demographic aggregates to reduce re-identification risk.

### 13.5 Application, infrastructure and supply-chain security

- Enforce strict CORS, CSP, HSTS, trusted-proxy configuration, secure headers, request/body limits and canonical paths.
- Use parameterised SQL, output encoding and server-side URL allowlists; disallow user-controlled server fetches by default.
- Rate limit by account, IP, endpoint sensitivity and federation for bulk workflows; use distributed storage if horizontally scaled.
- Store secrets outside images/repositories, rotate them and use separate credentials per environment/service.
- Run dependency, secret, SAST, container and infrastructure scans in CI; generate an SBOM and sign production images.
- Run API authorisation tests aligned with OWASP API Security Top 10 and an independent penetration test before go-live.

### 13.6 Audit requirements

- Audit login/MFA/session events, role/membership changes, all workflow transitions, sensitive record views, document downloads, calculations, overrides, imports, exports and administrator job actions.
- Include actor, effective roles/scope, request ID, event type, target, before/after field summary, source IP/device context, outcome and UTC timestamp.
- Redact secrets and document contents. Hash-chain or otherwise integrity-protect audit events and restrict write access to the audit subsystem.
- Make audit records append-only to application users and define retention, archive and legal-hold rules.

## 14. Testing and quality strategy

### 14.1 Test layers

| Layer | Minimum coverage |
|---|---|
| Go unit tests | Scoring, deadlines, state transitions, validation, access policies, reconciliation and redaction. |
| Repository integration tests | Migrations, constraints, scoped queries, transactions, concurrency and PostgreSQL behaviour. |
| API contract tests | OpenAPI conformance, error format, pagination, permissions and idempotency. |
| Vue component tests | Forms, filters, KPI states, accessible charts/tables, imports and permission-based rendering. |
| E2E tests | Federation submission/approval, NCS review, corrections, dashboard refresh, export and reminder paths. |
| Security tests | BOLA/BFLA, cross-federation access, upload attacks, CSRF, token/session abuse, CSV injection and sensitive caching. |
| Accessibility tests | Automated rules plus keyboard and screen-reader manual testing. |
| Performance tests | Dashboard p95, bulk imports, concurrent submissions, exports, worker backlog and summary refresh. |
| Recovery tests | Backup restoration, failed migration, object recovery, worker restart and provider outage. |

### 14.2 Critical test scenarios

- A federation officer cannot read, infer, export or download another federation’s records.
- A Federation President cannot approve a report they authored when separation of duties applies.
- Draft data never changes official dashboard totals.
- An approved amendment changes the current view but leaves prior snapshot/report reproducible.
- Deadline calculations remain correct at period/year boundaries and under rule-version changes.
- A failed email provider does not roll back a valid submission and retries do not duplicate messages.
- Duplicate job execution does not duplicate obligations, notifications, medals or exports.
- Dashboard totals reconcile to drill-down results under identical filters and scope.
- Financial arithmetic uses exact decimals and reconciles disbursement/accountability amounts.
- Safeguarding aggregates suppress small cohorts and never reveal case identities to ordinary dashboard users.
- Malicious files, oversized imports, formulas in CSV exports and spreadsheet injection payloads are neutralised.

### 14.3 Release gates

- All migrations apply to a production-like copy and the recovery/forward-fix procedure is rehearsed.
- Unit, integration, contract, component, E2E, accessibility and security suites pass.
- No unresolved critical/high vulnerability without documented risk acceptance and expiry.
- Dashboard reconciliation is exact for seeded acceptance datasets.
- Product, technical, finance, safeguarding, security and data owners sign off their domains.
- Runbooks, alerts, dashboards, backup evidence and rollback steps are complete.

## 15. DevOps, observability and support

- Establish development, test, staging and production environments with isolated databases, object buckets, secrets and notification recipients.
- Add CI stages for formatting/linting, Go tests/vet, frontend tests/build, OpenAPI validation, migration checks, scans and image builds.
- Use a migration tool that records applied versions and locks concurrent execution; do not rely on container initialisation for upgrades.
- Deploy API and worker separately. Use rolling or blue/green deployment with backward-compatible schema transitions.
- Add readiness/liveness checks for API and worker and dependency checks that do not expose secrets.
- Collect structured logs, metrics and traces with request/job IDs and PII redaction.
- Alert on authentication anomalies, error rate, dashboard p95, stale indicators, failed reconciliation, job backlog, email failures, storage scan failures, database saturation and backup failure.
- Define service targets during inception; proposed initial targets are 99.5% monthly availability, dashboard p95 under 2.5 seconds for standard filters, source-to-dashboard freshness under 5 minutes after approval, and zero unreconciled official totals.
- Automate encrypted PostgreSQL and object-storage backups; document retention and perform quarterly restore drills.
- Maintain operator runbooks for failed reminders, stale dashboards, stuck jobs, compromised accounts, malicious uploads, data correction and rollback.

## 16. Phased implementation roadmap

The calendar depends on team size and policy turnaround. The sequence is mandatory even if durations change. Work should be delivered in thin, demonstrable increments rather than building all screens before integration.

### Phase 0 — Inception and baseline (2 weeks)

**Backend/data**

- Confirm system-of-record boundaries and inventory current APIs/migrations.
- Facilitate indicator-definition workshops with NCS technical, finance, governance and safeguarding owners.
- Resolve or assign owners/dates to the decisions in Section 5.
- Define data classification, retention and legal basis; start the DPIA.
- Produce target ERD, OpenAPI conventions, threat model and migration strategy.
- Create acceptance datasets with known expected dashboard totals.

**Frontend/product**

- Map role journeys, reporting tasks and dashboard drill-downs.
- Produce low-fidelity responsive wireframes and an accessible chart proof of concept.
- Define shared UI states and content terminology.

**Exit criteria**

- Approved scope, role/permission matrix, indicator dictionary, workflow, ERD, threat model, prioritised backlog and measurable non-functional requirements.

### Phase 1 — Platform security and shared services (3–4 weeks)

**Backend**

- Implement permissions plus federation memberships/scope.
- Harden sessions, MFA, CSRF, rate limiting and audit.
- Add private object storage, malware scanning and document metadata.
- Add migration runner, outbox, worker/job framework, notifications and versioned reference data.

**Frontend**

- Adopt capability-based navigation and route guards.
- Build shared form/document/workflow/dashboard components.
- Add test harnesses, error boundary, observability and accessibility baseline.

**Exit criteria**

- Security tests prove tenant isolation; protected documents are no longer publicly served; jobs and notifications are observable and retryable.

### Phase 2 — Federation profile and reporting foundation (3–4 weeks)

**Backend**

- Add federation, officer, reporting calendar/period, obligation, report/revision and review schemas/APIs.
- Implement state machine, completeness engine, deadline rules and reminder generation.
- Migrate/link existing association/federation records using an approved mapping and reconciliation report.

**Frontend**

- Deliver federation workspace, profile management, document checklist, reporting calendar, drafts, submission and President approval queue.
- Deliver NCS review queue and report history.

**Exit criteria**

- One pilot federation can complete the full secure profile and report approval lifecycle with audited evidence and reminders.

### Phase 3 — Governance and compliance dashboard MVP (3 weeks)

**Backend**

- Implement governance responses, rule-versioned compliance scoring, obligations and indicator summaries.
- Add governance dashboard/drill-down endpoints, reconciliation and controlled exports.

**Frontend**

- Deliver compliant/non-compliant, expired constitution and missing report views with score explanation and remediation workflow.

**Exit criteria**

- Governance totals reconcile to source reports; authorised users can explain every federation score and formula version.

### Phase 4 — Athlete, competition, medal, coach and officials (5–7 weeks)

**Backend**

- Implement athlete master/affiliation/status, competition participation/results, medals, coach and official certification modules.
- Add deduplication, bulk imports, evidence validation, approval flow and athlete/performance summaries.

**Frontend**

- Deliver registers, import preview/correction flows, competition reporting and athlete/performance dashboards.

**Exit criteria**

- Pilot data imports cleanly; medals and participation are traceable to approved evidence; athlete PII access is tested and audited.

### Phase 5 — Finance and equipment (4–5 weeks)

**Backend**

- Implement disbursement ledger, revenue/expenditure, accountability, reconciliation, equipment receipt/distribution and stock movement.
- Add finance/equipment indicators, ageing, exception rules and evidence controls.

**Frontend**

- Deliver finance forms/review queues, financial dashboard, equipment register, distribution workflow and reconciliation views.

**Exit criteria**

- Exact-decimal totals and stock balances reconcile; Finance permissions and evidence access pass negative security tests.

### Phase 6 — Talent and safeguarding (3–4 weeks)

**Backend**

- Implement talent assessment/progression/scholarship records and restricted safeguarding aggregates/cases.
- Add small-number suppression and restricted case audit controls.

**Frontend**

- Deliver talent pipeline/dashboard and separated safeguarding aggregate/case experiences.

**Exit criteria**

- Talent conversions are reproducible; safeguarding privacy and role controls receive specialist sign-off.

### Phase 7 — Executive reporting, hardening and rollout (3–4 weeks)

- Generate versioned board packs and annual statistics from locked snapshots.
- Complete load, accessibility, recovery, security and penetration testing.
- Train NCS administrators/departments and pilot federation users using role-based guides.
- Run pilot in parallel with the previous reporting method for at least one complete reporting cycle.
- Reconcile pilot results, correct migration/data-quality issues and execute phased federation onboarding.
- Establish hypercare, support SLAs, incident escalation and post-implementation review.

**Exit criteria**

- Go-live checklist signed; restore drill passed; monitoring/support active; pilot totals accepted; no unresolved critical risk.

## 17. Prioritised backlog and dependency map

| Priority | Epic | Depends on |
|---|---|---|
| P0 | Policy decisions, indicator dictionary and data classification | None |
| P0 | Permission/scope model and secure sessions/MFA | Role matrix |
| P0 | Private document service and audit extension | Data classification |
| P0 | Migration, worker, outbox, notification and reference-data platforms | Architecture baseline |
| P0 | Federation and reporting-period/obligation foundation | Shared platform |
| P0 | Versioned report workflow and approval | Federation/reporting foundation |
| P1 | Governance data and compliance engine | Approved definitions/workflow |
| P1 | Governance dashboard | Approved governance summaries |
| P1 | Athlete master and bulk import | Federation/reference/document services |
| P1 | Competition/results/medals | Athlete master |
| P1 | Athlete/performance dashboards | Approved athlete/performance data |
| P1 | Finance/disbursement/accountability | Federation/reporting/document services |
| P1 | Financial dashboard | Approved finance data and reconciliation |
| P2 | Coaches and technical officials | Federation/reference/document services |
| P2 | Equipment | Federation/reporting/document services |
| P2 | Talent progression/scholarships | Athlete master |
| P2 | Safeguarding | Specialist policy and restricted permissions |
| P2 | Executive board reports/annual statistics | Stable official indicators/snapshots |

## 18. Team ownership and governance

Recommended accountable functions:

- Product owner: scope, workflow, prioritisation and acceptance.
- NCS governance/data owners: compliance rules, federation definitions and reporting calendar.
- Technical Department data steward: athlete, competition, medal, coach, official and talent definitions/quality.
- Finance data steward: disbursement, accountabilities, financial categories and reconciliation.
- Safeguarding lead: lawful fields, access, suppression and incident handling.
- Technical lead: architecture, API contracts, code standards and technical risk.
- Security/privacy lead: threat model, DPIA, tests, access reviews and incident readiness.
- QA/accessibility lead: acceptance datasets, automation, accessibility and release evidence.
- DevOps/SRE owner: environments, deployment, monitoring, backup/recovery and runbooks.
- Federation pilot representatives: terminology, workflow usability and training feedback.

Use an architecture decision record (ADR) for material choices, a data dictionary for fields/indicators, OpenAPI for APIs, migrations for schema, and runbooks for operations. These artefacts must be updated in the same pull request as behavioural changes.

## 19. Maintainability standards

- One domain owner and CODEOWNERS rule per feature area.
- Small, reviewable pull requests linked to acceptance criteria and threat-model impact.
- Formatting, linting, tests, API compatibility, migrations and security scans enforced in CI.
- No duplicated KPI formulas across SQL, Go and Vue; calculation logic has one backend owner and version.
- Feature flags for unfinished modules, policy changes and gradual federation rollout.
- Deprecation policy for API fields and indicator definitions.
- Seed/reference data is versioned, reviewable and migratable; production configuration changes are audited.
- Support documentation includes data correction, report reopening, user transfer, membership expiry, import failure and dashboard reconciliation procedures.
- Quarterly reviews cover dependencies, permissions, retention jobs, performance budgets, restore evidence and stale feature flags.

## 20. Definition of done

A module or dashboard is complete only when:

- approved requirements, field definitions, permissions and data classification exist;
- additive migrations, constraints, indexes and tested recovery steps exist;
- API contract, validation, ownership and state transitions are implemented and tested;
- all mutations and sensitive reads produce appropriate audit events;
- private documents pass access, signature, size and malware checks;
- frontend loading, empty, zero, error, offline/retry, success and permission-denied states are implemented;
- responsive and WCAG 2.2 AA acceptance checks pass;
- indicator totals reconcile to source/drill-down data for the approved acceptance dataset;
- unit, integration, contract, component, E2E, security and performance tests pass;
- observability, alerts, runbook and support ownership exist;
- OpenAPI, data dictionary, ADRs and user/admin guidance are updated; and
- the relevant business/data owner signs off.

## 21. Immediate next actions

1. Approve this workplan as the delivery baseline and assign owners to the Section 5 decisions.
2. Run indicator-definition workshops and produce the first governed data dictionary.
3. Approve the target role/permission matrix, federation membership model and report state machine.
4. Complete the initial threat model and DPIA before schema implementation for PII, finance and safeguarding.
5. Produce the target ERD and OpenAPI skeleton for federation/reporting foundation.
6. Build acceptance datasets with hand-calculated governance, athlete, medal, finance and talent results.
7. Implement Phase 1 security/shared services before accepting sensitive production uploads.
8. Select a small, representative federation pilot group and agree the parallel-run success measures.

