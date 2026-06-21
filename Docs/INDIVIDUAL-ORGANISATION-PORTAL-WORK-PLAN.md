# Individual and Organisation Portal Separation Work Plan

> Implementation status (2026-06-21): core journey delivered. Public application entry points now redirect to the authenticated individual portal; Form 3/Form 10 approval provisions organisation profiles, memberships, secure email invitations, account contexts, and the organisation dashboard. Operational details are recorded in `PORTAL-AND-LOCATION-IMPLEMENTATION.md`.

**System:** National Sports Management Information System (NSMIS)  
**Prepared:** 21 June 2026  
**Status:** Planning only; implementation has not started  
**Canonical routes:** `/my-portal` and `/dashboard`

## 1. Executive decision

The system will use two clearly separated journeys:

1. **Individual journey — `/my-portal`:** every ordinary person registers and signs in as an individual. Individuals manage their own identity, start and track registration applications, respond to requests, receive organisation invitations, and see organisations they are connected to.
2. **Organisation/staff journey — `/dashboard`:** approved federations, associations, clubs, academies, facilities, and authorised NCS staff use role- and organisation-scoped dashboards. Ordinary individual accounts cannot use `/dashboard` unless they have accepted an active organisation membership or hold an NCS staff role.

All public “Apply Now”, “Apply Online”, “Register Association”, and equivalent application buttons will be removed. Registration applications will begin only inside `/my-portal`.

When NCS approves a federation, association, or club registration application, the backend will automatically create the corresponding organisation profile/workspace, establish its ownership/membership records, and send a secure activation or invitation email to the official organisation contact. The email may contain the profile identifier and activation link, but **must never contain a generated or plaintext password**.

## 2. Goals

- Give every person one durable individual identity.
- Make `/my-portal` the only place ordinary individuals can initiate registrations.
- Prevent public, anonymous, or organisation-context application initiation.
- Automatically turn approved registration applications into usable organisation profiles.
- Give each organisation a distinct dashboard, permissions, records, and audit trail.
- Support one person belonging to multiple organisations without sharing accounts.
- Keep NCS staff dashboards separate from organisation dashboards.
- Eliminate emailed passwords and shared organisation credentials.
- Preserve existing application history and already registered federations during migration.
- Make the distinction visible in navigation, authentication redirects, emails, and support documentation.

## 3. Terminology

| Term | Meaning |
|---|---|
| Individual | A natural person with a personal NCS account and `/my-portal` access |
| Organisation | A federation, association, club, academy, facility operator, or another approved institutional profile |
| Organisation workspace | The authenticated context containing an organisation’s dashboard, data, users, permissions, and settings |
| Organisation member | An individual user authorised to work in an organisation workspace |
| Applicant | The individual who starts and owns a registration application until organisational provisioning |
| Primary owner | The initial verified individual responsible for the newly provisioned organisation workspace |
| Official contact | Email/phone supplied and verified for the organisation during registration |
| Invitation | A single-use, expiring link allowing a person to join or activate access to an organisation workspace |
| Context | The active individual, organisation, or NCS-staff scope used to authorise a request |

## 4. Current-state assessment

### 4.1 Existing capabilities to reuse

- Ordinary self-registration creates a `user` role and routes applicants to `/my-portal`.
- `/my-portal` already lists applications and offers federation and community-club registration choices.
- Application records support drafts, signed forms, payment, submission, review, approval, rejection, and information requests.
- Form 3 supports national association/federation registration and renewal.
- Form 10 supports community sports club registration and renewal.
- Existing federation records, memberships, roles, permissions, reports, dashboards, and audit infrastructure can be reused.
- The system already has account activation/status controls, PINs, JWT authentication, and background-job/notification tables.

### 4.2 Current problems

1. `/apply` is currently a public route and application CTAs appear in the header, footer, homepage, and CMS seed content.
2. “Applicant” is inferred from having `user` or no role, rather than an explicit profile/context model.
3. Approval changes application status but does not provision an organisation atomically.
4. There is no shared organisation abstraction covering federation, association, and club profiles.
5. Federation membership exists, but equivalent ownership/membership is not modelled consistently for clubs and associations.
6. `/dashboard` primarily distinguishes staff from applicant roles, not the active organisation context.
7. A person connected to multiple organisations has no workspace selector.
8. There is no secure invitation/activation workflow for the official organisational contact.
9. Email delivery is still incomplete and must be made reliable before provisioning depends on it.
10. Existing CMS menus and footer settings can reintroduce `/apply` links unless validated and migrated.

## 5. Target user journeys

### 5.1 New individual registration

```text
Public website
   → Create individual account
   → Verify email/phone as required
   → Sign in
   → /my-portal
```

The individual account represents one person. Registration collects personal name, personal email, phone, password, consent, and verification state. It does not create a club, association, or federation.

After login, an ordinary user always lands on `/my-portal` unless they deliberately open an organisation invitation or select an existing organisation workspace.

### 5.2 Starting an organisation registration

`/my-portal` contains a prominent **Register an organisation** section with these choices:

| Choice | Organisation type | Initial statutory form |
|---|---|---|
| Register a national sports federation | `FEDERATION` | Form 3 or the legally confirmed federation flow |
| Register a national sports association | `ASSOCIATION` | Form 3 or the legally confirmed association flow |
| Register a community sports club | `CLUB` | Form 10 |

The selection determines:

- form and regulation;
- required fields and attachments;
- prescribed fee;
- reviewer queue;
- profile type to provision;
- initial organisation roles;
- dashboard modules enabled after approval.

The route should be explicit, for example:

```text
/my-portal/applications/new
/my-portal/applications/new/federation
/my-portal/applications/new/association
/my-portal/applications/new/club
/my-portal/applications/:id
```

The application remains owned by the individual applicant throughout draft and review.

### 5.3 Application workflow

```text
DRAFT
 → PENDING_SIGNATURE
 → PENDING_PAYMENT
 → SUBMITTED
 → UNDER_REVIEW
 → NEEDS_INFORMATION ↔ RESUBMITTED
 → APPROVED_PENDING_PROVISIONING
 → PROVISIONED
```

Rejection remains terminal unless NCS policy permits an appeal or fresh application.

`APPROVED_PENDING_PROVISIONING` separates the legal review decision from the technical creation of the organisation workspace. `PROVISIONED` means the profile, ownership, invitations, roles, and dashboard context were created successfully.

### 5.4 Automatic organisation provisioning

Approval triggers one idempotent provisioning command:

```text
Approved application
       |
       v
Validate approved form snapshot and official contacts
       |
       v
Create organisation base profile
       |
       +→ Create type-specific federation/association/club profile
       +→ Create organisation membership for applicant
       +→ Assign initial organisation role(s)
       +→ Create official-contact invitation if different
       +→ Create licence/registration linkage where legally applicable
       +→ Queue activation/invitation email
       +→ Record audit event and provisioning result
       |
       v
Application PROVISIONED
```

Provisioning must use a database transaction for local records and an outbox/background job for email. Retrying the command must return the same organisation instead of creating duplicates.

### 5.5 Secure credential delivery

The requested “profile credentials” email will use a secure invitation model:

- Include organisation name, public/profile reference, role offered, and activation URL.
- Use a random, single-use, hashed invitation token.
- Expire invitations after a configurable period such as 48 hours.
- Require the recipient to confirm identity and set their own password if they do not yet have an individual account.
- If the official email already belongs to an individual account, require login before accepting the invitation.
- Never send a generated password, existing password, PIN, JWT, or refresh token by email.
- Notify the applicant separately when provisioning is complete.
- Provide audited resend, revoke, and change-contact actions.

The organisation itself is not a shared login. People authenticate individually and receive organisation roles. If policy insists on an institutional username, it must still activate through a named custodian and must not be shared among staff.

### 5.6 Organisation dashboard access

After accepting membership:

```text
Individual signs in
  → /my-portal
  → selects “Open organisation dashboard”
  → chooses organisation if more than one
  → /dashboard with active organisation context
```

The active context should be represented by a server-validated organisation ID, never trusted from local storage alone. A route such as `/dashboard?organisation=<id>` may be used for navigation, but every API call must verify an active membership and required permission.

The dashboard header shows:

- organisation logo/name;
- organisation type;
- registration/licence state;
- active role;
- organisation switcher when the user has multiple memberships;
- **Return to My Portal** action.

### 5.7 NCS staff journey

NCS staff continue to use `/dashboard`, but in an explicit `NCS_STAFF` context. Staff landing pages, permissions, and navigation remain separate from federation/club/association modules.

Staff must not impersonate an organisation silently. Any support impersonation feature requires explicit permission, visible banners, reason capture, expiry, and audit logging.

## 6. Route and navigation policy

### 6.1 Public website

Remove application actions from:

- desktop and mobile header;
- homepage hero and CTA sections;
- footer and footer CMS defaults;
- public navigation CMS seed and existing menu values;
- licence/application landing content;
- public cards, banners, widgets, and static pages;
- any search/sitemap result advertising direct application entry.

Public pages may retain neutral actions:

- **Create Individual Account** → `/register`
- **Sign In** → `/login`
- **Learn about registration** → a public information page
- **Go to My Portal** → `/my-portal` for authenticated individuals

No public button should initiate or continue an application.

### 6.2 `/apply` compatibility

Do not immediately delete `/apply`; old bookmarks, emails, search results, and CMS content may still use it.

Recommended behavior:

| Caller | `/apply` result |
|---|---|
| Unauthenticated | Redirect to `/login?redirect=/my-portal/applications/new` |
| Authenticated individual | Redirect to `/my-portal/applications/new` |
| Authenticated organisation/staff context | Redirect to `/my-portal` or require switching to individual context |

After an agreed deprecation window, return a permanent redirect. Record redirect usage before removal.

### 6.3 Route guards

| Route | Allowed context |
|---|---|
| `/my-portal/**` | Authenticated individual identity |
| `/my-portal/applications/**` | Individual applicant/owner |
| `/dashboard` | Active organisation membership or NCS staff |
| `/nsmis/**` | Active authorised organisation or NCS staff, permission-dependent |
| `/applications` admin review | NCS reviewer only |
| public website | Anonymous or authenticated, read-only |

An ordinary user without organisation membership or staff roles who requests `/dashboard` is redirected to `/my-portal` with an explanatory message.

## 7. Domain and database design

### 7.1 Organisation base table

Add `organisations`:

| Field | Purpose |
|---|---|
| `id` | Internal identifier |
| `profile_reference` | Unique public/support reference |
| `organisation_type` | `FEDERATION`, `ASSOCIATION`, `CLUB`, `ACADEMY`, `FACILITY_OPERATOR`, etc. |
| `legal_name` | Approved name from the form snapshot |
| `display_name` | Public/dashboard name |
| `registration_number` | NCS or legal registration number when issued |
| `official_email`, `official_phone` | Verified organisational contacts |
| `status` | `PENDING_ACTIVATION`, `ACTIVE`, `SUSPENDED`, `REVOKED`, `ARCHIVED` |
| `source_application_id` | Unique approved application that created the profile |
| `created_by`, timestamps | Audit ownership |
| `version` | Optimistic concurrency |

Unique constraints should prevent one approved application from creating multiple organisations and prevent duplicate active registration numbers.

### 7.2 Type-specific profiles

Use one base organisation plus extensions:

- `federations` links to `organisation_id` and retains NSMIS federation/reporting fields.
- `associations` or `organisation_association_profiles` stores association-specific fields if legally distinct.
- `clubs` or `organisation_club_profiles` stores district, community, parent association/federation, facilities, and club-specific fields.

Do not store the whole approved form only as mutable JSON. Keep an immutable approved-form snapshot for evidence and copy operational fields into typed profile columns.

### 7.3 Organisation membership

Add shared `organisation_memberships`:

| Field | Purpose |
|---|---|
| `organisation_id`, `user_id` | Member relationship |
| `role_id` or `membership_role` | Organisation-scoped role |
| `status` | Invited, active, suspended, ended |
| `is_primary_owner` | Initial accountable owner |
| `starts_at`, `ends_at` | Effective access period |
| `invited_by`, `accepted_at` | Lifecycle evidence |

Migrate or bridge existing `federation_memberships` to this shared model. Avoid two independent permission sources after cutover.

### 7.4 Organisation roles

Recommended initial roles:

- `organisation_owner`
- `organisation_admin`
- `organisation_president`
- `organisation_general_secretary`
- `organisation_finance_officer`
- `organisation_technical_officer`
- `organisation_data_entry`
- `organisation_viewer`

Permissions remain resource/action based and are always evaluated with organisation scope.

### 7.5 Invitations

Add `organisation_invitations`:

- organisation and intended role;
- normalized email/phone;
- hashed token and expiration;
- invited-by user/system;
- source application;
- sent, accepted, revoked, and expired timestamps;
- attempt/resend count;
- status and last delivery error.

Tokens must be single-use. Changing the intended contact revokes outstanding invitations.

### 7.6 Application changes

Add or formalise:

- `applicant_user_id` as immutable owner;
- `requested_organisation_type` with a strict enum/check;
- `provisioned_organisation_id` unique nullable foreign key;
- `provisioning_status`, attempt count, error category, timestamps;
- immutable approved form-data version/hash;
- official contact email/phone separated from applicant contacts;
- approval decision and provisioning transition history.

The application is never reassigned to the organisation after provisioning; it remains historical evidence owned by the applicant and linked to the resulting organisation.

## 8. Approval and provisioning service

Introduce an `OrganisationProvisioningService` rather than placing logic directly in the approval handler.

### 8.1 Approval transaction

1. Lock the application.
2. Authorise reviewer and verify its current state.
3. Validate complete, immutable approved data.
4. Record approval and `APPROVED_PENDING_PROVISIONING`.
5. Commit the legal decision.
6. Queue an idempotent `PROVISION_ORGANISATION` job.

### 8.2 Provisioning job

1. Lock application and check idempotency.
2. Map form type and selected journey to organisation type.
3. Create the organisation and extension record.
4. Create applicant membership as primary owner/admin, subject to approved policy.
5. Link registration/licence records.
6. Create official-contact invitation if necessary.
7. Write audit and outbox events.
8. Set application to `PROVISIONED` and store the organisation ID.
9. Commit.
10. Deliver email asynchronously with retry/dead-letter handling.

If email fails, the organisation remains provisioned but `activation_delivery_status` shows the failure. Staff can resend without repeating provisioning.

### 8.3 Mapping rules

Use versioned server-side mappings, not frontend labels:

```text
form_3 + FEDERATION → federation profile/workspace
form_3 + ASSOCIATION → association profile/workspace
form_10 + CLUB → club profile/workspace
```

NCS legal/product owners must confirm whether federation and association use the same form and which transformation flows apply.

## 9. Authentication and context authorization

Authentication proves the person. Organisation membership determines what that person can access.

JWT/session claims may include stable user identity and broad system roles, but organisation permissions should be revalidated from the database or short-lived scoped claims. Do not add every organisation permission permanently to a long-lived token.

Recommended context endpoint:

```http
GET /api/v1/me/contexts
POST /api/v1/me/contexts/{organisationId}/activate
DELETE /api/v1/me/context
```

Every organisation API request verifies:

- user is active;
- organisation is active;
- membership is active and in date;
- requested organisation matches active context;
- required permission is granted;
- object belongs to that organisation.

## 10. API work

### Individual portal

- `GET /api/v1/my-portal/summary`
- `GET /api/v1/my-portal/applications`
- `POST /api/v1/my-portal/applications`
- existing draft, attachment, signing, payment, and submit operations moved/aliased under individual ownership
- `GET /api/v1/my-portal/organisations`
- `GET /api/v1/my-portal/invitations`
- `POST /api/v1/my-portal/invitations/{id}/accept`
- `POST /api/v1/my-portal/invitations/{id}/decline`

### Organisation profile

- `GET /api/v1/organisations/{id}`
- `PUT /api/v1/organisations/{id}` with field-level permissions
- `GET /api/v1/organisations/{id}/members`
- `POST /api/v1/organisations/{id}/invitations`
- `PUT /api/v1/organisations/{id}/members/{userId}`
- `DELETE /api/v1/organisations/{id}/members/{userId}` with last-owner protection

### Provisioning administration

- `GET /api/v1/admin/provisioning-jobs`
- `GET /api/v1/admin/applications/{id}/provisioning`
- `POST /api/v1/admin/applications/{id}/provisioning/retry`
- `POST /api/v1/admin/organisation-invitations/{id}/resend`
- `POST /api/v1/admin/organisation-invitations/{id}/revoke`

All mutations require idempotency keys where a retry could duplicate an organisation, invitation, or membership.

## 11. Frontend work

### 11.1 Public layout cleanup

- Remove `/apply` CTAs from `PublicLayout.vue`, `HomeView.vue`, footer defaults, CMS menu defaults, and migrations/seeds.
- Update CMS validation to warn or block new public navigation links pointing to `/apply`.
- Change application-related public copy to “Sign in to My Portal to apply.”
- Retain clear Sign In and Create Account actions.

### 11.2 `/my-portal`

Build a task-oriented individual dashboard:

- personal profile/verification status;
- **Register an organisation** cards;
- drafts and actions due;
- submitted applications and status timeline;
- information requests;
- invitations awaiting acceptance;
- organisations the person belongs to;
- **Open dashboard** action for active memberships;
- help and support.

Forms render inside the `/my-portal` shell and preserve autosave, review, signing, payment, and tracking behavior.

### 11.3 `/dashboard`

Create a context-aware shell:

- NCS staff home for staff context;
- federation home for federation context;
- association home for association context;
- club home for club context;
- shared navigation components filtered by permissions;
- profile completion, licence, compliance, reports, athletes, officers, documents, finance, and other modules enabled per type.

An organisation type changes available modules, not authentication identity.

### 11.4 Context switcher

If a person has multiple contexts, show:

```text
Personal Portal
Uganda Example Federation — Organisation Owner
Example Community Club — Finance Officer
NCS Staff — Technical Department
```

Switching context refreshes server-authorised permissions and clears organisation-specific cached data.

## 12. Email templates

Required templates:

- application approved, provisioning in progress;
- organisation workspace ready;
- organisation access invitation;
- invitation accepted;
- invitation expiring;
- invitation revoked;
- provisioning or delivery failure for NCS operations;
- membership/role changed;
- organisation suspended or restored.

Every email includes support contacts, expiry, organisation reference, and a warning not to forward activation links. Email links use HTTPS, short expiry, and one-time tokens.

## 13. Security and abuse controls

- Never create or email plaintext passwords.
- Do not create a shared “federation account” used by multiple people.
- Verify official contact changes and record who authorised them.
- Prevent applicants from selecting arbitrary organisation IDs.
- Enforce ownership and organisation scope in repositories, not only route guards.
- Protect against duplicate organisation names/registration numbers using policy-aware matching and review queues.
- Require separation of duties for approval and sensitive provisioning correction.
- Protect the last active organisation owner from accidental removal.
- Invalidate active organisation contexts when membership or organisation status changes.
- Audit approval, provisioning, invitation, acceptance, context switching, role changes, and profile edits.
- Rate-limit invitation sends and activation attempts.
- Do not reveal whether an invitation email already belongs to another account.
- Require MFA for organisation owners and NCS staff before production if policy permits.

## 14. Data migration

### 14.1 Existing individual users

- Treat `user/applicant` accounts as individual profiles.
- Backfill an explicit individual profile record if introduced.
- Preserve existing credentials and application ownership.

### 14.2 Existing federations

- Create one `organisations` record per active federation.
- Link each existing `federations` row to its organisation.
- Convert federation memberships to shared organisation memberships.
- Preserve historical membership and report ownership.
- Detect duplicate contacts before sending invitations.

### 14.3 Approved historical applications

- Reconcile approved registration applications with existing federation/club/association records.
- Link matches after human review where confidence is not exact.
- Provision only genuinely missing organisation profiles.
- Never send activation email automatically for an uncertain historical match.

### 14.4 CMS/public links

- Query `cms_menus`, footer settings, homepage settings, posts, pages, and custom widgets for `/apply`.
- Replace/remove links under an auditable migration.
- Maintain `/apply` redirect telemetry during deprecation.

## 15. Implementation phases

### Phase 0 — Policy and journey confirmation (1 week)

- Confirm organisation types and statutory form mapping.
- Confirm who becomes primary owner after approval.
- Confirm official-contact verification and disputed ownership procedures.
- Approve route, email, and no-shared-password policies.
- Inventory every public Apply CTA and `/apply` dependency.

**Exit:** signed journey, role, provisioning, and migration decisions.

### Phase 1 — Organisation/context foundation (2 weeks)

- Add organisation, membership, role, invitation, application-link, and provisioning migrations.
- Add repository/service layer with organisation scoping.
- Add context enumeration/activation APIs.
- Seed organisation roles and permissions.
- Add integration tests and migration rollback rehearsal.

**Exit:** context authorization works independently of UI.

### Phase 2 — Approval provisioning and email (2 weeks)

- Add approval-to-provisioning state transition.
- Implement idempotent provisioning worker and outbox.
- Create organisation/type-specific profiles and memberships.
- Implement invitation activation and reliable email delivery.
- Add provisioning admin visibility, retry, resend, and correction tools.

**Exit:** an approved test application consistently creates one usable organisation workspace.

### Phase 3 — Individual portal and form relocation (2 weeks)

- Move registration selection and forms into `/my-portal`.
- Add individual summary, drafts, applications, invitations, and organisations.
- Redirect `/apply` based on authentication/context.
- Add ownership tests and preserve draft resume links.

**Exit:** no application can start anonymously; end-to-end individual submission passes.

### Phase 4 — Organisation dashboards (2-3 weeks)

- Build context-aware `/dashboard` shell and switcher.
- Map federation/association/club modules and permissions.
- Add return-to-personal-portal behavior.
- Verify all backend object queries enforce organisation scope.

**Exit:** each organisation type sees only its permitted data and modules.

### Phase 5 — Public/CMS cleanup and migration (1 week)

- Remove Apply buttons and public direct-entry UI.
- Update CMS defaults and add `/apply` link validation.
- Migrate existing organisations, memberships, approved applications, and public content.
- Enable redirects and telemetry.

**Exit:** production content has no direct public application entry; reconciliation is signed off.

### Phase 6 — Pilot and release (1-2 weeks)

- Pilot one federation, association, and club journey.
- Test invitation delivery, expiry, resend, multiple memberships, and ownership disputes.
- Complete security, accessibility, responsive, and support testing.
- Train reviewers and support staff.
- Rehearse provisioning failure and rollback.

**Exit:** launch approval with no unresolved high-severity issue.

## 16. Testing plan

### Unit

- form-to-organisation mapping;
- provisioning idempotency;
- role and module mapping;
- invitation expiry/token verification;
- route decision rules;
- status transitions;
- last-owner protection.

### PostgreSQL integration

- approval and concurrent provisioning retries;
- unique application-to-organisation constraint;
- membership scope and cross-organisation denial;
- invitation acceptance with existing/new users;
- suspension invalidates context;
- historical backfill reconciliation.

### API/security

- ordinary user denied `/dashboard` without membership;
- organisation member denied other organisations;
- revoked member loses access immediately;
- applicant cannot provision directly;
- reviewer cannot bypass required approval state;
- invitation replay and expired-token rejection;
- no password/token leakage in logs, email previews, or APIs.

### End-to-end

1. Individual registers and reaches `/my-portal`.
2. Individual starts federation/association/club application.
3. Draft, signature, payment, submit, review, and information-request flows work.
4. Approval provisions exactly one organisation.
5. Official contact receives and accepts activation invitation.
6. Member opens the correct `/dashboard` context.
7. Individual can return to `/my-portal`.
8. Public `/apply` links redirect correctly and public pages contain no Apply CTA.

## 17. Observability and operations

Monitor:

- applications started/completed by organisation type;
- approval-to-provisioning latency;
- provisioning success/failure/retry/dead-letter counts;
- invitation queued/delivered/accepted/expired/resend rates;
- duplicate-profile prevention events;
- active organisation memberships;
- denied cross-organisation access;
- `/apply` redirect usage;
- application completion and abandonment after relocation.

Provide runbooks for provisioning failure, invitation delivery failure, incorrect official contact, duplicate organisation, ownership dispute, suspended organisation, last-owner loss, and context-access incidents.

## 18. Risks

| Risk | Mitigation |
|---|---|
| Emailed shared credentials are compromised | Use named-user, single-use invitation activation; never email passwords |
| Approval creates duplicate profiles | Unique source-application link, transaction, idempotency key, reconciliation |
| Wrong person becomes owner | Verified policy fields, reviewer confirmation, official-contact invitation, dispute workflow |
| One person manages several organisations | Shared membership model and server-authorised context switcher |
| Federation and association legal mapping is ambiguous | Phase 0 legal confirmation and versioned mapping |
| Existing federation records conflict with applications | Human-assisted reconciliation before backfill/provisioning |
| Public/CMS `/apply` links reappear | CMS validation, content scan, redirect telemetry, acceptance test |
| Email failure blocks access | Provision independently, retain pending invitation, retry/resend/admin support |
| Cross-organisation data leakage | Repository scoping, object ownership checks, integration and penetration tests |
| Ordinary users are confused by two destinations | Clear labels: “My Portal” for personal work and “Organisation Dashboard” for institutional work |

## 19. Acceptance criteria

### Individual portal

- Ordinary users always land on `/my-portal`.
- Only `/my-portal` can start federation, association, or club applications.
- Individuals can track drafts, submissions, requests, invitations, and memberships.
- Ordinary users without another context cannot access `/dashboard`.

### Provisioning

- Approval creates exactly one type-correct organisation profile.
- Applicant/official membership follows approved ownership rules.
- Provisioning is idempotent, audited, observable, and retryable.
- Invitation email contains no plaintext password or authentication secret.
- Application links permanently to the resulting organisation.

### Organisation dashboard

- Active organisation members can enter `/dashboard` for authorised organisations.
- Multiple memberships can be switched without data leakage.
- Modules and actions reflect organisation type, role, and permission.
- Suspended/revoked memberships and organisations lose access immediately.

### Public website and CMS

- No public Apply Now/Application button remains.
- Public information may direct users to create an account or sign in, but cannot start an application.
- `/apply` redirects safely throughout deprecation.
- CMS cannot silently reintroduce a direct public application link.

## 20. Decisions required before implementation

1. Are federation and association both Form 3, and how are they distinguished legally?
2. Which form and profile type covers clubs beyond Form 10 community clubs?
3. Does the applicant automatically become primary owner, or must NCS approve a named official?
4. Which official contact fields are mandatory and how are they verified?
5. Can one official email/phone belong to multiple organisations?
6. Which organisation roles are available per type at initial provisioning?
7. Must organisation owners use MFA at launch?
8. What happens when official contact and applicant are different people?
9. What is the ownership-dispute and account-recovery process?
10. Which historical approved applications should be provisioned automatically?
11. How long should `/apply` compatibility redirects remain?
12. Which dashboard modules are enabled for federation, association, and club on day one?

## 21. Recommended implementation boundary

The first release should support:

- individual `/my-portal` registration and application ownership;
- federation, association, and community-club registration choices;
- approval-driven profile provisioning;
- applicant/official invitation and activation;
- one organisation context at a time with a switcher;
- organisation-scoped `/dashboard`;
- NCS staff `/dashboard`;
- removal and redirection of public application entry points;
- audit, retries, support tools, and migration of existing federations.

Defer self-service ownership transfer, complex organisation mergers/splits, institutional service accounts, bulk member import, and additional organisation types until the core model is proven.
