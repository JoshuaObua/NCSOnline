# NCS USSD Quick Services Implementation Plan

**Proposed callback:** `POST /ncs-ussd`  
**System:** National Sports Management Information System (NSMIS)  
**Prepared:** 21 June 2026  
**Status:** Planning document; implementation has not started

## 1. Executive summary

This plan adds a provider-neutral USSD channel to the existing Go/PostgreSQL system so people with basic phones, feature phones, smartphones without data, or unreliable internet can use a focused set of NCS services.

The recommended rollout is:

1. **Public, read-only verification:** athlete profile verification, licence validity, federation recognition, and general help.
2. **Protected self-service:** application tracking and renewal eligibility after authentication with the account’s registered phone number and existing 4–6 digit NCS PIN.
3. **Renewal initiation and hand-off:** create or resume a renewal application through USSD, then send an SMS link or assisted-service reference for documents, signatures, and payment.

USSD should not reproduce the entire web application wizard. The current statutory process requires long forms, document uploads, a signed form, and payment evidence. USSD is best used for discovery, verification, protected read-only account summaries, renewal checks, reminders, and hand-off to web or staff-assisted completion.

### Implemented authentication update — 21 June 2026

- `POST /ncs-ussd` is implemented as the provider callback.
- Athlete-registration and licence/credential-validity lookups remain public and privacy-limited.
- Registered-user application and renewal information requires the registered Uganda phone number plus the account’s existing bcrypt-hashed NCS PIN.
- When the provider supplies the caller MSISDN, it must match the entered registered phone.
- PIN authentication grants read-only USSD summaries only; documents, submissions, profile changes, signing and payment remain in the normal authenticated channels.
- SMS OTP is retained as a possible later recovery/step-up mechanism, not the primary USSD login method.

## 2. Objectives

- Make high-value NCS services available without mobile data.
- Provide a trusted way to verify an athlete, licence, credential, or federation.
- Let applicants track progress without signing into the website.
- Let licence holders check renewal eligibility and begin renewal from any phone.
- Preserve the web system as the authoritative record and avoid a separate USSD-only data silo.
- Protect personal data by returning only the minimum information required.
- Make the USSD gateway replaceable without rewriting the domain services.
- Record every lookup, state transition, renewal request, and provider callback for support and audit.

## 3. Backend audit summary

### 3.1 Current architecture

| Area | Current implementation | USSD relevance |
|---|---|---|
| HTTP API | Go with Chi routing in `backend/cmd/server/main.go` | Add a small public provider callback at `/ncs-ussd` outside JWT browser routes. |
| Data store | PostgreSQL with ordered SQL migrations | Add licence, USSD session, lookup audit, OTP, and notification-channel migrations. |
| Authentication | JWT, refresh tokens, role checks, active-user revalidation, account invalidation | Browser JWT cannot authenticate a carrier callback. Add provider authentication and subscriber OTP/PIN separately. |
| User identity | Email-first users with an optional phone and a hashed PIN | Normalize and verify phone numbers before using them for USSD identity. Do not assume phone values are unique today. |
| Applications | Draft-to-review workflow with references, status, payment state, signed form, and attachments | Good foundation for tracking and renewal hand-off. No public reference lookup exists. |
| Athlete registry | Unique `athlete_number`, profile/status, affiliations, `verified_at`, consent basis, and status history | Strong foundation for athlete verification after a privacy-safe query is added. |
| Federation registry | Unique NCS registration number and recognition status | Strong foundation for a public recognition check. |
| Credentials | Coaches have licence number/expiry; officials have certification validity | Can support credential checks, but should be unified behind the new licence/credential service. |
| Issued licences | No authoritative general licence table | Blocking gap for reliable licence validity and renewal. Approval of an application is not itself an issued licence. |
| Notifications | Email/in-app notification tables and worker scaffolding | Add `SMS` and optional `USSD` event channels, delivery receipts, retries, and templates. |
| Payments | Proof upload is implemented; online payment initiation is a stub | USSD payment must remain out of MVP or be added only after an approved mobile-money integration. |
| Audit/security | Request IDs, security headers, request limits, audit middleware, RBAC, geo-blocking | Reuse audit patterns, but add callback signatures, replay protection, MSISDN redaction, and USSD-specific limits. |
| Tests | Middleware unit tests and a small NSMIS handler test set | Add contract, state-machine, repository, integration, security, load, and provider-simulator tests. |

### 3.2 Existing records that can be reused

- `athletes`: athlete number, name, discipline, status, national-team status, `verified_at`, and soft deletion.
- `athlete_affiliations`: athlete-to-federation relationship and effective dates.
- `athlete_status_history`: traceable athlete status changes.
- `federations`: registration number, recognition status, identity, and contact metadata.
- `applications`: application reference, form type, application type, status, payment status, reviewer notes, and timestamps.
- `coaches`: licence number, expiry date, and status.
- `technical_officials`: certification, validity date, and status.
- `users`: phone, PIN hash, account status, and fraud/suspension controls.
- `notifications`, `background_jobs`, and `audit_logs`: patterns for asynchronous delivery and traceability.

### 3.3 Material gaps to close before launch

1. There is no dedicated issued-licence registry or licence status history.
2. Approved applications do not automatically create an issued licence.
3. Applications do not link a renewal to the licence it renews.
4. Athlete verification is currently available only through authenticated, federation-scoped generic CRUD.
5. There are no purpose-built, privacy-limited quick-service queries.
6. Phone numbers are not yet a verified, normalized, unique identity key.
7. There is no SMS OTP/delivery implementation.
8. There is no provider callback verification, USSD session store, or USSD state machine.
9. Notification channels currently allow only email and in-app delivery.
10. Online/mobile-money payment initiation is not implemented.
11. API contracts are not published as OpenAPI and USSD provider contracts are absent.
12. Database migrations 016-018 still require staging execution and validation according to the current implementation status.

## 4. Proposed quick-service catalogue

### 4.1 Release priority

| Priority | Service | Access | MVP action | Required foundation |
|---|---|---|---|---|
| P0 | Verify athlete registration | Public with lookup controls | Enter the unique NCS athlete number; return a masked, limited registration summary | Purpose-built athlete verification query |
| P0 | Check licence validity | Public with lookup controls | Enter licence number; return type, holder display name, status, issue/expiry date | New authoritative licence registry |
| P0 | Check federation recognition | Public | Enter registration number; return recognition state and public name | Public federation verification query |
| P0 | Help and NCS contacts | Public | Return short contact/help choices | Configured content and language copy |
| P1 | Track applications | Phone + PIN | Authenticate the registered user and return a limited recent-application summary | Normalized registered phone, existing PIN and safe status projection |
| P1 | Check renewal eligibility | Phone + PIN | Authenticate, identify licence, show validity and blockers | Licence registry and renewal policy service |
| P2 | Start or resume renewal | Phone + PIN plus web hand-off | Create/resume a Form 3 or Form 10 draft/intent; send completion link/reference | Licence-to-application link and SMS |
| P1 | Coach/official credential check | Public with lookup controls | Validate licence/certification number and expiry | Credential adapter or migration into licence registry |
| P2 | Reminder preferences | Phone + PIN | Subscribe/unsubscribe to expiry and application SMS alerts | SMS channel and consent records |
| P2 | Assisted-service request | Phone + PIN | Create a callback/support case with a reference | Support queue and service-level rules |
| Deferred | Complete renewal entirely in USSD | Not recommended | Long form, files, signature, and evidence remain web/assisted | Major legal, payment, and document redesign |
| Deferred | USSD payment | Not in MVP | Add only after gateway, reconciliation, refund, callback, and fraud controls are operational | Approved payment provider integration |

### 4.2 Public response data policy

Public lookup responses must not expose full date of birth, phone, email, address, national ID, passport data, uploaded documents, reviewer notes, safeguarding data, or internal IDs.

Recommended public projections:

- **Athlete:** masked/display name, athlete number, sport/discipline, current public status, current federation, and last verification date.
- **Licence:** licence number, licence type, public holder name, status, issue date, expiry date, and issuing authority.
- **Federation:** name, NCS registration number, recognition status, and last status update date.
- **Coach/official:** display name, credential type/level, credential status, and expiry date.

Where policy does not permit a name to be public, return a masked name such as `A*** K***` or no name at all. NCS legal/data-protection approval is a launch gate.

## 5. Proposed USSD journeys

### 5.1 Main menu

```text
CON NCS Quick Services
1. Verify athlete
2. Check licence
3. Track application
4. Renew licence
5. Federation status
6. Help
```

If multilingual service is approved, language selection should occur on the first session and be remembered against a consented phone preference:

```text
CON Choose language
1. English
2. Luganda
3. Other configured language
```

Menu copy must be tested against the selected provider's character and screen limits. Copy belongs in versioned configuration/templates, not hard-coded handlers.

### 5.2 Athlete verification

```text
CON Enter NCS athlete number:
> NCS-2026-AB12CD34

END Verified athlete
A*** K***
Athletics | ACTIVE
Federation: [public name]
Verified: 12 Jun 2026
```

Rules:

- Normalize spaces and hyphens without permitting fuzzy enumeration.
- Require the exact unique NCS athlete number; do not offer fuzzy name searching.
- Return one generic not-found response and only the approved masked public projection.
- Rate-limit by provider, phone hash, athlete number hash, and IP.
- Do not reveal whether an athlete exists after repeated failures.

### 5.3 Licence validity

```text
CON Enter NCS licence number:
> NCS-LIC-2026-001234

END VALID
Community Sports Club
Kampala Youth Sports Club
Expires: 30 Jun 2027
```

Possible public states:

- `VALID`
- `EXPIRING_SOON`
- `EXPIRED`
- `SUSPENDED`
- `REVOKED`
- `PENDING_ISSUE`
- `NOT_FOUND` (generic wording; never infer internal records)

The displayed state must be computed by the licence service from status, effective dates, suspension/revocation history, and replacement/supersession—not trusted from a stale string alone.

### 5.4 Application tracking

```text
CON Enter your registered phone number:
> 0772123456

CON Enter your 4-6 digit NCS PIN:
> 1234

END Application: UNDER REVIEW
Payment: VERIFIED
Next: Wait for NCS review.
Updated: 18 Jun 2026
```

Rules:

- Do not return application status until the registered phone and PIN are verified.
- If the provider supplies the caller MSISDN, require it to match the entered registered phone.
- Reject inactive accounts, unset PINs, and accounts requiring a PIN change with the same generic authentication failure.
- If no usable registered phone/PIN exists, direct the user to the website or assisted support without confirming private details.
- Map internal status to plain-language next actions through a central status projection service.
- Do not expose reviewer notes in USSD unless they have a specific applicant-safe/public field.

### 5.5 Renewal eligibility and initiation

```text
CON Enter licence number:
> NCS-LIC-2025-000918

CON Enter your registered phone number and NCS PIN when prompted.

CON Licence expires 31 Jul 2026.
Renewal is available.
1. Start renewal
2. Resume renewal
3. Requirements
0. Exit

END Renewal started.
Reference: NCS-2026-F3-XYZ789
We sent a secure completion link by SMS.
```

USSD may collect only short, low-risk selections such as renewal type, preferred contact language, and consent to begin. The web or assisted channel completes statutory fields, attachment uploads, signing, and payment.

Renewal rules must answer:

- Is the licence type renewable?
- Is the caller within the configurable renewal window?
- Is there already an active renewal application?
- Is the licence suspended, revoked, superseded, or under enforcement action?
- Are there outstanding information requests, fees, accountabilities, or mandatory reports?
- Which official form applies (for example Form 3 or Form 10)?
- What is the next safe action?

## 6. Target architecture

```text
Mobile subscriber
      |
      v
Mobile network / USSD aggregator
      |
      | HTTPS callback
      v
POST /ncs-ussd
      |
      v
Provider verification + request normalization
      |
      v
USSD session/state-machine service
      |
      +--> Quick-service facade
      |      +--> Athlete verification service
      |      +--> Licence service
      |      +--> Application status service
      |      +--> Renewal eligibility/service
      |      +--> Federation verification service
      |
      +--> OTP/SMS notification service --> SMS provider
      |
      +--> PostgreSQL session, domain, audit and outbox tables
      |
      +--> Existing background worker for delivery/retries
```

### 6.1 Design rules

- `/ncs-ussd` is a provider callback, not a browser page and not a public JSON API.
- The handler must contain no domain SQL and no provider-specific menu logic.
- A provider adapter translates inbound fields into one internal `USSDRequest` and translates `USSDResponse` back into the provider's required plain-text format.
- Domain services are channel-neutral so web, USSD, an IVR, or a future chatbot can use the same rules.
- One transition is processed transactionally per session request.
- Duplicate provider callbacks return the previously stored response.
- Every session has a maximum lifetime, maximum steps, and inactivity expiry.
- All external sends use an outbox/background-job pattern; the callback never waits on an SMS provider.
- The database remains the source of truth; session context contains references, not copied profiles.

## 7. Endpoint and internal contract

### 7.1 Public provider callback

```http
POST /ncs-ussd
Content-Type: application/x-www-form-urlencoded
```

The first adapter should support the chosen provider's exact contract. A common normalized request is:

```go
type USSDRequest struct {
    Provider          string
    ProviderRequestID string
    SessionID         string
    ServiceCode       string
    MSISDN            string
    Input              string
    FullText           string
    NetworkCode        string
    ReceivedAt         time.Time
}

type USSDResponse struct {
    Continue bool
    Message  string
}
```

The adapter may render `CON <message>` for a continuing session and `END <message>` for a completed session if required by the selected provider. This rendering convention must be confirmed during provider onboarding rather than assumed for every gateway.

### 7.2 Callback middleware order

1. Request ID and panic recovery.
2. Dedicated small body limit.
3. Content-type enforcement.
4. Provider identification.
5. HMAC/signature or mTLS verification.
6. Timestamp/replay-window verification.
7. Provider/IP/service-code allow-list.
8. USSD-specific rate limit.
9. Strict callback decoding and normalization.
10. Idempotency lookup.
11. Session transition.
12. Redacted structured logging and metrics.

Do not apply browser JWT middleware to this route. Do not trust `phoneNumber`, source IP, or provider name without the selected provider's authentication mechanism.

### 7.3 Optional channel-neutral quick-service API

Avoid calling the backend over HTTP from itself. Implement Go services/repositories first. If these capabilities must later be exposed to trusted channels, add separately authenticated endpoints such as:

- `POST /api/v1/quick-services/athletes/verify`
- `POST /api/v1/quick-services/licences/verify`
- `POST /api/v1/quick-services/federations/verify`
- `POST /api/v1/quick-services/applications/request-otp`
- `POST /api/v1/quick-services/applications/status`
- `POST /api/v1/quick-services/renewals/eligibility`
- `POST /api/v1/quick-services/renewals/start`

These must not be unauthenticated copies of internal repository methods. Use service credentials, scopes, and the same privacy/rate controls as USSD.

## 8. Data model changes

Create additive migrations after migration `018`; never modify an applied migration.

### 8.1 Authoritative licence registry

Proposed `licences` fields:

| Field | Purpose |
|---|---|
| `id` | Internal UUID/text primary key |
| `licence_number` | Unique public identifier |
| `licence_type` | Federation, association, club, academy, facility, competition, coach, official, or configured type |
| `subject_type`, `subject_id` | Polymorphic holder link validated by the service layer |
| `holder_display_name` | Approved public name snapshot |
| `source_application_id` | Approved application that caused issue |
| `status` | Pending, valid, suspended, expired, revoked, cancelled, superseded |
| `issued_on`, `valid_from`, `expires_on` | Effective dates |
| `renewable` | Whether renewal is supported |
| `renewal_opens_on` | Rule-derived or stored eligibility date |
| `supersedes_licence_id` | Previous licence in the chain |
| `public_verification_code_hash` | Optional QR/short-code verification secret, not the secret itself |
| `version` | Optimistic concurrency |
| audit fields | Issuer, updater, timestamps, reason |

Add:

- `licence_status_history`
- unique active licence constraints appropriate to each licence type
- indexes on normalized licence number, expiry/status, subject, and application
- an approval/issue transaction that creates the licence only after authorized issuance
- a reconciliation report for approved applications that have no issued licence

### 8.2 Renewal linkage

Add to `applications`:

- `renewal_of_licence_id` nullable foreign key
- `channel_started` such as `WEB`, `USSD`, or `STAFF`
- `ussd_intent_id` nullable trace link
- unique partial constraint preventing multiple active renewals for one licence

Add a versioned `renewal_rules` table or configuration model containing licence type, form type, opening window, grace period, blockers, effective dates, and policy version.

### 8.3 Phone identity and consent

- Add normalized E.164 phone storage and a verified timestamp.
- Decide whether a verified phone is unique per active user; resolve shared-family-phone policy before enforcing uniqueness.
- Never use the existing app PIN directly as a USSD PIN without a security review.
- Add OTP challenge records containing a hashed code, purpose, phone hash, expiry, attempts, consumed timestamp, and provider message reference.
- Add communication consent/preference history with purpose, channel, language, source, and timestamp.

### 8.4 USSD operational tables

`ussd_sessions`:

- provider, provider session ID, service code
- encrypted MSISDN where operationally required and a keyed hash for indexing/rate limits
- current state and version
- small JSON context with non-sensitive references
- language, started/last-active/expiry/completion timestamps
- status: active, completed, timed out, failed, blocked

`ussd_events`:

- session ID and monotonic sequence
- provider request ID/idempotency key
- input category or redacted value; do not store raw OTP/PIN
- previous/next state
- response text/version
- outcome, error code, duration, request ID, timestamp

`quick_service_lookups`:

- service type, subject identifier hash, phone hash
- result category, disclosure policy version, request/session ID
- timestamp and retention/expiry fields

`service_intents`:

- intent type, licence/application/user links
- channel and status
- idempotency key
- hand-off token hash and expiry
- created/completed timestamps

### 8.5 Notification changes

- Extend notification channel constraints to include `SMS`.
- Add SMS templates with length-tested content and language variants.
- Store provider message ID, delivery status, failure category, retry schedule, and receipt timestamp.
- Never put full DOB, national ID, sensitive case information, or unmasked private data in SMS.

## 9. Application and licence domain changes

### 9.1 Separate approval from issuance

Current application approval changes application status. The new workflow should be:

```text
APPLICATION APPROVED
        |
        v
LICENCE PENDING_ISSUE
        |
 authorized issuance and number allocation
        v
LICENCE VALID
```

This avoids reporting an approved application as a valid licence before the legal issuance act, effective dates, and licence number exist.

### 9.2 Status computation

Implement one `LicenceStatusAt(licence, date)` service used by web, USSD, exports, and admin screens. Precedence should be policy-tested, for example:

1. revoked/cancelled
2. suspended for the requested date
3. superseded
4. before `valid_from`
5. after `expires_on`
6. expiring soon
7. valid

Do not duplicate this logic in SQL snippets, handlers, or menu code.

### 9.3 Renewal orchestration

`StartRenewal` must be idempotent:

1. Lock the licence row.
2. Recompute status and eligibility.
3. Find an existing active renewal and return it if present.
4. Create the appropriate application draft with `application_type=RENEWAL`.
5. Link it to the old licence and USSD service intent.
6. Generate a one-time hand-off token with a short expiry.
7. Queue an SMS containing the reference and secure HTTPS link.
8. Commit the transaction and return a USSD-safe summary.

The hand-off link should require additional authentication before showing private data.

## 10. State-machine design

Define typed states rather than parsing the entire `text` history on every request. Example states:

- `WELCOME`
- `LANGUAGE`
- `ATHLETE_NUMBER`
- `ATHLETE_BIRTH_YEAR`
- `LICENCE_NUMBER`
- `FEDERATION_NUMBER`
- `APPLICATION_REFERENCE`
- `APPLICATION_PHONE`
- `APPLICATION_PIN`
- `RENEWAL_LICENCE_NUMBER`
- `RENEWAL_PHONE`
- `RENEWAL_PIN`
- `RENEWAL_MENU`
- `RENEWAL_CONFIRM`
- `HELP_MENU`
- terminal states

Each state handler accepts the normalized request and current context and returns:

- next state
- sanitized context changes
- continue/end decision
- versioned message key and parameters
- domain command, if any
- audit outcome

Every state must define valid input, retry message, maximum retries, back/home behavior, timeout behavior, and disclosure level.

Recommended navigation:

- `0` exits or returns to the previous menu according to the displayed prompt.
- `00` returns to the main menu if the provider supports it reliably.
- Invalid input does not advance state.
- After three invalid attempts, end with a neutral support message.
- Never echo OTP, PIN, full phone number, or sensitive lookup input.

## 11. Security, privacy, and abuse controls

### 11.1 Provider trust

- Require TLS and a verified callback domain.
- Prefer mTLS or signed requests; otherwise use provider-specific secrets plus strict allow-lists.
- Validate timestamp and nonce/provider request ID within a small replay window.
- Rotate secrets and support overlapping keys during rotation.
- Keep development/simulator authentication separate from production.

### 11.2 Subscriber authentication

- Treat MSISDN supplied by the provider as a routing signal, not sufficient authentication for private data.
- Use the account’s existing bcrypt-hashed 4–6 digit PIN with the normalized registered phone for application and renewal summaries.
- Never store, log, echo, or place a raw PIN in session/event tables.
- Limit failed attempts and concurrent sessions by caller MSISDN, entered phone hash, provider and session.
- Keep OTP as an optional future recovery or step-up mechanism rather than the primary login.
- Do not log codes or return distinct “account exists” errors.
- Require stronger web/staff authentication for profile changes, documents, signing, payment, or changes of ownership.

### 11.3 Enumeration prevention

- Rate-limit by provider, IP, service code, session, phone hash, and target identifier hash.
- Add progressive cooldowns and anomaly alerts.
- Use generic failure responses.
- Cap sessions and lookups per day.
- Block known automated patterns without exposing the blocking rule.
- Maintain an emergency kill switch per quick service.

### 11.4 Data handling

- Normalize MSISDN to E.164 only after validation.
- Store a keyed hash for lookup/rate limiting; encrypt the full value only where delivery/support requires it.
- Redact phone, PIN/OTP, application reference, athlete number, and licence number in ordinary logs.
- Define retention periods for sessions/events/lookups and automate deletion or irreversible aggregation.
- Version the public disclosure policy so historical responses are auditable.
- Complete a data-protection/privacy impact assessment before production.

### 11.5 Availability

- Keep callback processing below the provider timeout, with a target p95 under 1 second and p99 under 2 seconds.
- Use short database timeouts and indexed exact lookups.
- Avoid remote calls in the callback path.
- Queue SMS and other side effects asynchronously.
- Use idempotency to survive callback retries.
- Return a safe terminal message when dependencies are unavailable.
- Provide service-level dashboards and a provider-independent simulator.

## 12. Proposed Go package layout

```text
backend/internal/ussd/
  models.go
  service.go
  state_machine.go
  messages.go
  provider.go
  providers/<selected_provider>.go

backend/internal/handlers/
  ussd.go

backend/internal/services/
  quick_service.go
  athlete_verification_service.go
  licence_service.go
  renewal_service.go
  otp_service.go
  notification_service.go

backend/internal/repository/
  ussd.go
  licences.go
  quick_services.go
  otp.go

backend/cmd/ussd-simulator/
  main.go
```

The simulator is for local/staging testing only and must not provide a production bypass.

## 13. Configuration

Add environment/config entries only after provider selection:

```text
USSD_ENABLED=false
USSD_PROVIDER=<adapter-name>
USSD_SERVICE_CODE=
USSD_CALLBACK_SECRET=
USSD_SIGNING_KEYS=
USSD_ALLOWED_CIDRS=
USSD_SESSION_TTL=180s
USSD_MAX_STEPS=15
USSD_LOOKUPS_PER_PHONE_DAY=10
USSD_REQUEST_TIMEOUT=2s

SMS_PROVIDER=
SMS_API_KEY=
SMS_SENDER_ID=NCS
OTP_TTL=5m
OTP_MAX_ATTEMPTS=3
OTP_PEPPER=
PUBLIC_LOOKUP_HASH_KEY=
FIELD_ENCRYPTION_KEY=
```

Secrets belong in the deployment secret manager, not `.env.example` values, logs, images, or source control. Startup validation should fail closed when USSD is enabled with missing security configuration.

## 14. Provider procurement and onboarding workstream

Before binding the implementation to an aggregator or mobile network, NCS should confirm:

- availability of a shared or dedicated Uganda short code
- session and per-request pricing
- supported networks and roaming behavior
- callback field contract and response format
- callback authentication options (HMAC, mTLS, IP ranges)
- timeout, retry, character, pagination, and session-length limits
- simulator and sandbox availability
- delivery receipts and reporting
- data residency, subprocessors, retention, breach notification, and support SLA
- short-code approval, sender ID approval, and regulatory/operator lead times
- production failover and incident escalation contacts

Build the provider interface first, but implement only one adapter for MVP. A second adapter is justified after the service and contract are stable.

## 15. Implementation work plan

The durations below are planning estimates and should be re-estimated after policy approval and provider selection. Provider/short-code procurement can run in parallel and may take longer than engineering.

### Phase 0 — Policy, service definition, and provider selection (1-2 weeks)

- Assign product owner, technical owner, data-protection owner, operations owner, and NCS licensing authority.
- Approve the quick-service catalogue and what is explicitly excluded.
- Confirm who may see athlete/licence holder names.
- Define authoritative licence types, numbering, validity, suspension, revocation, renewal windows, and grace periods.
- Decide shared-phone and phone-change policies.
- Select the provider and obtain exact callback/security specifications.
- Complete threat model, privacy impact assessment, retention schedule, and incident path.
- Produce English and approved translated menu copy.

**Exit:** signed service policy, disclosure matrix, provider contract/specification, and approved MVP journeys.

### Phase 1 — Licence and identity foundation (2 weeks)

- Add licence, licence history, renewal rule, phone verification, consent, OTP, and renewal linkage migrations.
- Build licence issuance/status/reconciliation services.
- Add an admin issuance/correction path with separation of duties and audit events.
- Backfill issued licences from authoritative NCS records; do not infer them blindly from approvals.
- Clean and normalize user phones; resolve duplicates before constraints.
- Add public projection methods for athlete, licence, credential, and federation.
- Add application safe-status/next-action projection.

**Exit:** authoritative staged licence data, reconciliation signed off, domain tests passing, no USSD endpoint exposed.

### Phase 2 — Quick-service domain layer (1-2 weeks)

- Implement athlete registration verification by exact unique NCS athlete number with masking and layered rate limits.
- Implement licence and federation verification.
- Implement OTP request/verify and application tracking.
- Implement renewal eligibility and idempotent start/resume.
- Extend notification model/worker to SMS with an outbox and receipts.
- Add rate-limit and lookup-audit repositories.
- Document service contracts and error taxonomy.

**Exit:** services pass unit and PostgreSQL integration tests independently of USSD.

### Phase 3 — USSD channel and `/ncs-ussd` (2 weeks)

- Add provider interface and selected provider adapter.
- Add callback authentication, replay protection, allow-lists, body/content-type limits, and dedicated rate limits.
- Add session/event tables and transactional state machine.
- Add versioned/menu templates and language framework.
- Implement P0 and P1 journeys.
- Add idempotent callback response storage.
- Add the local/staging provider simulator.
- Update Nginx with a dedicated exact-match location and limits for `/ncs-ussd`.

**Exit:** complete sandbox journeys with no private-data leakage, duplicate-callback safety, and provider contract tests passing.

### Phase 4 — Operations, security, and controlled pilot (2 weeks)

- Add metrics for sessions, completions, abandonment by step, lookup outcomes, latency, OTP delivery, renewal conversion, provider errors, and abuse blocks.
- Add alerts for failure rate, latency, signature failures, provider retry spikes, database saturation, SMS failures, and suspicious enumeration.
- Produce support runbook, incident runbook, secret rotation, provider outage, kill-switch, and data-correction procedures.
- Run load tests at forecast peak plus headroom.
- Run security review/penetration tests focused on spoofing, replay, enumeration, OTP abuse, SQL/input injection, and disclosure.
- Pilot with NCS staff and a small group across supported networks and device types.
- Reconcile every pilot result against web/admin records.

**Exit:** signed pilot report, resolved high-severity findings, rollback rehearsal, and launch approval.

### Phase 5 — Production rollout and optimization (ongoing)

- Enable help and verification first.
- Enable tracking after OTP delivery metrics are stable.
- Enable renewal initiation after licence reconciliation reaches the approved threshold.
- Expand language/network coverage based on evidence.
- Review abandoned states and message comprehension weekly during early operation.
- Audit disclosure, retention, provider billing, and fraud controls monthly.
- Consider payment only as a separately approved project.

## 16. Testing strategy

### 16.1 Unit tests

- Every state/input transition and retry limit.
- Licence status precedence and boundary dates.
- Renewal eligibility, blocker, existing-draft, and idempotency cases.
- Athlete disclosure/masking rules.
- Application status-to-next-action mapping.
- OTP expiry, hashing, attempts, reuse, and cooldown.
- Provider request parsing and response rendering.

### 16.2 Integration tests

- Real PostgreSQL migrations up/down or forward/rollback rehearsal.
- Concurrent duplicate renewal attempts.
- Duplicate/reordered callback processing.
- Session version conflicts and expiry.
- Licence issuance and application linkage.
- SMS outbox retries and delivery receipt reconciliation.
- Audit event completeness and sensitive-field redaction.

### 16.3 Contract tests

- Captured provider sandbox fixtures for first, continuation, final, retry, malformed, unsigned, and timed-out requests.
- Exact content type, status code, character encoding, response prefix, and maximum length.
- Nginx forwarding of provider signature and request ID headers.

### 16.4 Security tests

- Signature bypass, old timestamps, replayed request IDs, wrong service code, and spoofed MSISDN.
- Athlete/licence/application enumeration at multiple identifiers.
- OTP spraying, resend abuse, brute force, and race conditions.
- Injection and oversized/malformed input.
- PII leakage in response, logs, metrics, traces, and alerts.
- Session fixation, cross-session access, and state skipping.

### 16.5 Performance and resilience tests

- Peak concurrent sessions and provider retry storms.
- p95/p99 callback latency under expected production data volume.
- Database connection-pool exhaustion.
- SMS provider outage, delayed receipts, and duplicate receipts.
- Safe responses during database/worker/provider degradation.

### 16.6 User acceptance matrix

Test at minimum:

- major supported mobile networks
- basic phones, feature phones, Android, and iPhone
- short/long names and identifiers
- low balance where the provider model is user-paid
- session timeout and reconnect
- English and every approved translation
- valid, expiring, expired, suspended, revoked, and missing licences
- active, inactive, suspended, retired, unverified, and missing athlete records
- every application status and next action

## 17. Observability and reporting

Use low-cardinality metrics and never put phone numbers, athlete numbers, licence numbers, or application references in metric labels.

Core measures:

- sessions started/completed/timed out/failed
- completion and abandonment by journey/state
- provider and network code (where contractually available)
- callback latency and database latency
- signature/replay/rate-limit rejection count
- verification outcome category
- OTP queued/delivered/failed/verified
- renewal eligible/started/resumed/completed on web
- hand-off conversion from USSD to web/assisted service
- provider cost per completed service
- support incidents and data correction requests

## 18. Deployment and rollback

- Add `USSD_ENABLED` and per-service feature flags, defaulting to false.
- Deploy additive schema and dormant code before exposing the callback.
- Validate database backfill/reconciliation independently.
- Restrict callback access to sandbox provider credentials first.
- Use a staging short code and non-production SMS sender.
- Enable journeys separately: help, federation, athlete, licence, tracking, renewal.
- Keep a provider-facing maintenance response ready.
- Rollback by disabling the affected journey or callback; do not roll back migrations destructively under live data.
- Preserve idempotency and audit records through rollback.

## 19. Risk register

| Risk | Impact | Mitigation | Owner |
|---|---|---|---|
| Short-code/provider onboarding is delayed | Engineering finishes before the service can be piloted | Start procurement in Phase 0; use the simulator and sandbox while approvals proceed | Product/procurement |
| Historical licence data is incomplete or inconsistent | Incorrect validity responses damage trust | Establish an authoritative register, reconcile records, and launch licence lookup only after signed data-quality thresholds | Licensing authority/data lead |
| Athlete lookup enables enumeration or excess disclosure | Privacy harm and regulatory exposure | Secondary lookup factor, masking, generic failures, layered limits, disclosure approval, and anomaly monitoring | Security/data protection |
| Shared or recycled phone numbers cause wrong-account access | Private application or renewal data is disclosed | Verified phone policy, OTP, assisted recovery, ownership-change controls, and no MSISDN-only authentication | Identity/product owner |
| Provider callback retries create duplicate renewals or messages | Duplicate records and user confusion | Stored callback responses, idempotency keys, row locking, unique active-renewal constraints, and outbox deduplication | Backend lead |
| USSD timeout is exceeded | Sessions fail despite correct business logic | Exact indexed lookups, no synchronous SMS calls, short DB timeouts, load tests, and latency alerts | Platform lead |
| SMS delivery is delayed or unavailable | OTP and hand-off journeys cannot complete | Delivery receipts, bounded retries, expiry-aware resend, alternate assisted path, and provider SLA monitoring | Operations |
| USSD rules diverge from web rules | Different channels report conflicting results | One channel-neutral service layer and contract tests used by all channels | Backend/product |
| Translated copy changes meaning or exceeds limits | Users select the wrong action or menus truncate | Professional review, device/network testing, versioned templates, and message-length checks | Content/product |
| Mobile-money scope is added prematurely | Fraud, reconciliation, and refund failures | Keep payment out of MVP and require a separately approved payment workstream | Programme sponsor |
| Sensitive values leak into logs or metrics | Security/privacy incident | Structured redaction, prohibited-field tests, keyed hashes, access controls, retention automation, and periodic audit | Security/operations |
| Service becomes unavailable during provider/network outage | Users cannot reach critical checks | Maintenance response, health monitoring, provider escalation, published fallback contacts, and per-service kill switches | Operations |

Review this register at every phase exit. High or critical residual risks require written acceptance by the accountable NCS owner.

## 20. Acceptance criteria

### Functional

- `/ncs-ussd` processes the selected provider's verified callback contract.
- All displayed menus support back/exit, invalid input, timeout, and terminal responses.
- Athlete, licence, credential, and federation results match authoritative records.
- Application status is disclosed only after successful subscriber verification.
- Renewal start is idempotent and links to exactly one source licence.
- Complex renewal work is handed off with a secure, expiring link or assisted-service reference.
- Settings can be reset/disabled operationally without deployment.

### Security and privacy

- Unsigned, replayed, stale, malformed, and unauthorized callbacks are rejected.
- Private services require registered phone plus PIN and enforce attempt/cooldown limits.
- Public responses comply with the approved disclosure matrix.
- Logs, metrics, traces, and events contain no raw OTP or prohibited PII.
- Enumeration tests meet the approved abuse thresholds.
- High/critical security findings are closed before launch.

### Reliability and operations

- Callback latency meets the selected provider's limit with agreed headroom.
- Duplicate callbacks return the original result without duplicate side effects.
- Outbound SMS is asynchronous, retryable, and reconciled.
- Monitoring, alerts, kill switches, dashboards, and runbooks are exercised.
- Pilot results reconcile against the web/admin system.

## 21. Decisions required from NCS

1. Which licence types are legally issued by NCS and which are only registrations, recognitions, permits, or certifications?
2. What fields may be shown publicly for an athlete and licence holder?
3. Is the masked athlete-registration projection approved for public lookup by exact NCS athlete number?
4. Which languages are required at launch?
5. Which networks and USSD provider/aggregator will be used?
6. Will the session be toll-free, NCS-paid, or user-paid?
7. Is a shared short code acceptable for pilot, and is a dedicated code required for production?
8. What are the renewal windows, grace periods, blockers, and escalation rules per licence type?
9. Can users share a phone number, and how will ownership changes be handled?
10. Which SMS provider/sender ID and OTP policy are approved?
11. What retention periods apply to sessions, lookup audits, phone data, and delivery records?
12. Who can issue, suspend, revoke, correct, and reinstate a licence?
13. What is the manual fallback when records are missing or incorrect?
14. What service/support hours and incident SLA will NCS publish?

## 22. Recommended MVP boundary

The MVP should include:

- English language
- athlete verification
- licence validity
- federation recognition
- application tracking with registered phone plus PIN
- renewal eligibility and start/resume with SMS hand-off
- help/contact information
- one USSD provider and one SMS provider
- admin/audit visibility, monitoring, and operational kill switches

The MVP should exclude:

- document upload through USSD
- signatures through USSD
- changing personal/profile information
- full reviewer notes or sensitive case data
- USSD mobile-money payment
- unverified application tracking
- automatic licence issuance solely because an application is approved

This boundary delivers the broad-access value of USSD while keeping statutory, financial, and sensitive workflows in channels that can support the necessary evidence and security.
