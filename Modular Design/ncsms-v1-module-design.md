# NCSMS v1.0 Blueprint Design

## Overview

This design document defines the National Council of Sports Management System (NCSMS) v1.0 architecture, targeting Go 1.22+, Supabase PostgreSQL 15+, Realtime, Storage, and Edge Functions.

- Market Context: National Council of Sports (NCS), Uganda
- Compliance Baselines: National Sports Act 2023, Uganda Data Protection and Privacy Act, Bank of Uganda Payment Regulations
- Architectural pattern: Modular Monolith in Go with internal packages, database-level multi-tenancy and strict transaction boundaries
- Domain scope: Super Admin, Admin, Federations modules

## Core System Architecture

### Primary Stack

- Language: Go 1.22+
- Database: PostgreSQL 15+ via Supabase
- DB Driver: `github.com/jackc/pgx/v5/pgxpool`
- Auth: Supabase Auth with JWT claims and role metadata
- Data isolation: `ncs_core` schema plus Supabase native `auth` schema
- Security: RLS policies, private storage access, signed URL generation

### Architectural Data Flow

Client Layer (Responsive Web Portal) -> HTTPS/TLS 1.3 -> Supabase Auth & Gateway -> Go Core ModMon Application Engine + Supabase Storage Engine -> Supabase/PostgreSQL Engine

### Repository Pattern

- `cmd/api/main.go` - bootstrap and router initialization
- `/internal/config` - environment, secrets, crypto parameter loading
- `/internal/middleware` - auth, RBAC, audit, request controls
- `/internal/domain` - per-domain handler/service/repo/models
- `/pkg` - shared helpers: crypto, supabase client, validation

## Database Schema Principles

### Schema and Extensions

- `CREATE SCHEMA IF NOT EXISTS ncs_core;`
- Extensions:
  - `uuid-ossp`
  - `pgcrypto`

### Core Entities

- `user_profiles` - profile mapping to `auth.users`
- `federations` - federation registry
- `licenses` - licensing records
- `grants` - funding allocations
- `grant_accountabilities` - accountability submissions
- `assets` - inventory register
- `athletes` - athlete master register
- `audit_logs` - immutable audit tracking

### Indexing Strategy

- explicit indexes on `id`, `federation_id`, `license_number`, `reg_number`, `status`
- RLS enforcement on sensitive tables
- composite or partial indexes where required by access patterns

## Security and Hardening

### Row-Level Security (RLS)

- Enable RLS on operational tables: `user_profiles`, `federations`, `licenses`, `athletes`, `grants`, `audit_logs`
- Helper functions:
  - `ncs_core.get_auth_id()` -> `auth.uid()`
  - `ncs_core.get_auth_role()` -> claim `user_metadata.primary_role`

### Policy Examples

- Federation public select only on `ACTIVE`
- Federation officer full access for admin roles
- Athlete access restricted to admin and designated federation officers
- Audit logs read-only for `SYSTEM_ADMIN` and `AUDITOR`

### Storage Security

- Private buckets only
- partitioned paths like `/receipts/:grant_id/:file_id`, `/athletes/:athlete_id/identity.pdf`
- signed URLs maximum 15 minutes
- access policies based on JWT metadata

### Request Tampering and Rate-Limit Reinforcement

- Normalize incoming headers and ignore unreliable client-supplied proxy headers such as `X-Forwarded-For`, `X-Forwarded-Host`, `X-Client-IP`, `X-Remote-IP`, `X-Remote-Addr`, `X-Host`, and `X-Original-URL` unless they are validated by a trusted gateway.
- Enforce request path normalization before ACL or firewall decisions. Reject obfuscated forms such as `%2e/admin`, `/admin/.`, `//admin//`, `/./admin/..`, `/admin..;/`, and uppercase variants like `/aDmIN` where the normalized path resolves to a protected endpoint.
- Canonicalize request payload values before rate-limit checks to remove trailing null bytes (`%00`), CRLF characters, and whitespace that can bypass filtering.
- Reject requests with untrusted or syntactically malformed URL encodings, dot-segments, or semicolons in protected endpoint paths.
- Enforce a single source of truth for client identity via the connection remote address and trusted proxy chain configuration, not client-supplied headers.
- Fingerprint requests consistently by IP, authenticated session, endpoint, normalized path, user agent, and content hash to prevent bypass via user-agent/cookie/IP variation.
- Reject requests with unexpected or ambiguous random query string parameters when they are used to evade endpoint-based rate limiting.
- Apply strict JSON body validation to prevent payload variants like trailing spaces from bypassing controls.
- Treat web cache poisoning patterns as a security event and validate `X-Original-URL` and similar rewrite headers only when they are explicitly allowed and trusted.
- Record bypass detection events in audit logs with full context for security response and tuning.

### 2FA / OTP Bypass Reinforcement

- Enforce step-by-step 2FA workflow validation on the server side; do not allow direct access to post-2FA endpoints without a valid session state transition.
- Validate request referer and session state for 2FA endpoints, but do not rely solely on the Referer header for authorization.
- Ensure tokens are one-time use and bound to the authenticated user session; reject reuse across accounts or reuse after successful consumption.
- Do not expose OTP or 2FA tokens in any response payloads, logs, or URLs.
- Ensure password reset flows do not implicitly bypass 2FA: require re-validation or a separate secure flow after reset, and terminate previous sessions.
- Bind 2FA sessions and reset flows to a single user account and invalidate previous or parallel sessions once a new code is generated.
- Use strong rate limiting for OTP validation, resend operations, and backup code consumption, including silent flood-detection and replay protection.
- Prevent reset-code regeneration logic from resetting the brute-force counter; each resend should count toward the same protection budget unless a trusted reset event occurs.
- Normalize OTP codes and request parameters to remove trailing blank characters and null bytes before validation.
- Protect backup code endpoints and 2FA setup endpoints from CORS/XSS exposure and require strict origin validation.
- Expire previous sessions on 2FA enablement and require fresh authentication when the account’s security configuration changes.
- Audit 2FA and reset events with full context, including actor, path, request metadata, and success/failure state.

### Unrestricted Resource Consumption (API4:2023)

- Apply explicit per-request resource consumption limits for any endpoint that triggers external work, such as SMS delivery, email notifications, file generation, report runs, or third-party service calls.
- Validate request payloads against allowed ceilings for resource usage: message size, phone number quantity, target recipients, attachment volume, and processing complexity.
- Return `429 Too Many Requests` or a similar quota response when resource budget limits are exceeded, instead of allowing uncontrolled processing.
- Use throttling and bandwidth controls at both application and network layers to prevent abuse of expensive operations.
- Maintain usage counters for user, organization, and IP address to detect and throttle high-volume resource consumption patterns.
- Monitor and log resource consumption metrics in real time and correlate them with suspicious access patterns.
- Protect resource-intensive operations with additional authorization checks and require explicit consent or higher privilege for bulk actions.
- Perform load testing and capacity planning for all hosted APIs to identify and prevent resource exhaustion vectors.
- Enforce strict validation of request parameters to prevent bypasses via malformed or duplicated inputs.
- Apply limits on endpoint-specific actions such as SMS sends, OTP requests, report exports, or large data queries.

### Unrestricted Access to Sensitive Business Flows (API6:2023)

- Enforce authorization and access control on every business flow endpoint, especially operations that perform purchases, transfers, approvals, or privileged actions.
- Validate the caller’s identity, role, and entitlement before allowing sensitive flows such as ticket purchases, order creation, account changes, or financial transactions.
- Implement explicit business rules and permission checks in service logic, not only at the routing layer.
- Apply least privilege principles: separate create/update operations from read/list operations and require stronger privileges for high-risk actions.
- Use API gateways and centralized access control to manage sensitive flows and ensure only approved clients and roles may call them.
- Log and monitor sensitive business flow access attempts, including rejected authorization events and suspicious path usage.
- Validate business payloads for correctness, dates, formats, ownership, and domain-specific constraints before processing.
- Protect sensitive endpoints with additional authorization scopes, multi-factor authentication, or approval workflows when required.

### Session Fixation Reinforcement

- Renew or replace session tokens after successful authentication; never preserve the same session identifier across login.
- Invalidate any existing session state when a user authenticates, and issue a new session token on successful login.
- Avoid storing sensitive session state in cookies before authentication.
- Use secure session cookies with `HttpOnly`, `Secure`, and `SameSite` flags.
- Bind session tokens to strong authentication context and user-specific values.
- Detect and reject session identifiers that arrive from untrusted or unexpected sources.
- On logout or authentication failure, clear session state and expire previous cookies.
- Audit session creation, renewal, and termination events with user and request metadata.

### Server Side Request Forgery (API7:2023)

- Validate and sanitize any user-supplied URL or host parameter before using it in server-side requests.
- Restrict outbound request destinations to a whitelist of trusted domains and IP ranges.
- Block or deny requests that reference internal services, metadata endpoints, or private network addresses.
- Normalize URL inputs and reject DNS rebinding or redirect-based manipulations.
- Use network-level egress controls, firewall rules, or service meshes to prevent unauthorized outgoing requests.
- Log SSRF-related access attempts, including the requested destination and caller context.
- Apply strict validation to any endpoint that fetches remote resources, such as images, PDFs, or webhooks.

### Security Misconfiguration (API8:2023)

- Enforce proper authorization on all endpoints and never expose sensitive business or configuration flows without role validation.
- Apply least-privilege ACLs at both API gateway and application layers.
- Avoid default or overly-permissive settings for authentication, CORS, TLS, and service discovery.
- Ensure implementation-level configuration matches the intended security posture for each environment.
- Harden API management and gateway configurations to prevent accidental exposure of administration endpoints.
- Monitor for unauthorized access, configuration drift, and suspicious usage of sensitive endpoints.
- Use automated configuration analysis and policy checks to detect misconfiguration before deployment.

### Password Reset Broken Logic

- Validate password reset tokens on the server side before allowing any password change.
- Generate a unique reset token for each request and enforce expiration and single-use semantics.
- Bind the reset token to the intended user and do not allow password resets without matching user identity.
- Use HTTPS for all password reset flows and avoid exposing tokens in URLs or logs.
- Keep error messages generic to prevent user enumeration and token probing.
- Notify users on successful password resets and log attempts for monitoring.
- Apply rate limiting to password reset requests and token validation operations.
- Require additional verification or secondary factors for high-risk password reset workflows.

### Broken Object Level Authorization (API1:2023)

- Enforce object-level access control for every request that uses an identifier from the user.
- Verify that the authenticated user is authorized to access or mutate the requested resource before loading it.
- Do not return or modify objects based solely on supplied identifiers without checking ownership or role permissions.
- Apply authorization checks both in controllers/handlers and in service/repository layers to prevent bypass.
- Log denied object access attempts and use RBAC/ABAC policies to control data access.
- Avoid using user-supplied IDs directly in queries without verifying the caller’s entitlement.

### Broken Function Level Authorization (API5:2023)

- Verify function-level access permissions for each endpoint based on user roles and privileges.
- Do not assume that authenticated users are authorized to call all exposed functions.
- Separate administrative and normal user paths, enforcing stronger checks on management endpoints.
- Use role checks, claims, or policy-based authorization before executing business logic.
- Log forbidden function access attempts and audit privilege escalations.
- Implement multi-level access controls and validate runtime permissions for every function.

### Broken Authentication (API2:2023)

- Require strong authentication mechanisms such as JWT, OAuth, or equivalent standards for all login operations.
- Validate credentials securely and avoid weak or custom authentication flows.
- Enforce authorization on protected endpoints and do not expose user-specific data without verifying identity.
- Use encrypted transport (TLS) for authentication and session exchanges.
- Apply account lockout or throttling for repeated failed login attempts.
- Ensure authenticated users can only access their own data unless explicitly authorized.
- Audit authentication successes and failures for suspicious behavior.

### Improper Inventory Management (API9:2023)

- Manage API versions explicitly and deprecate old versions in a controlled manner.
- Document current and legacy APIs, enforce version constraints, and disable unsupported endpoints.
- Restrict access to inventory and versioned resource endpoints by role and ownership.
- Use a centralized service layer to validate inventory operations rather than direct database writes.
- Require stronger privileges for inventory modification and protect administrative operations.
- Monitor traffic for calls to old or unsupported API versions and block unexpected version access.
- Maintain clear release policies and support timelines for each API version.

### Broken Object Property Level Authorization (API3:2023)

- Validate object property changes at the API boundary and in service logic before persisting updates.
- Only allow users to modify fields they are explicitly authorized to change.
- Reject or ignore unrequested property values in request payloads using explicit whitelists.
- Enforce property-level authorization checks for both read and write operations.
- Implement strong role-based access control (RBAC) and ownership checks for sensitive object properties.
- Regularly test APIs for unauthorized property access and validate payloads against expected schemas.
- Log attempts to modify or access unauthorized object properties for security review.

### Reset Password Bypass

- Protect password reset workflows against host header poisoning and forwarded-host manipulation.
- Ignore untrusted `Host`, `X-Forwarded-Host`, and duplicate host header values when generating reset links.
- Validate the password reset destination URL and bind it to the intended user before sending.
- Reject malformed or duplicate email parameters such as `email=victim&email=attacker`, `email=victim%20email=attacker`, `email=victim|email=attacker`, `cc:attacker`, `bcc:attacker`, and JSON arrays of emails.
- Prevent attackers from using their own token with a victim’s account or from reusing expired tokens.
- Rate limit reset requests and token validation operations to prevent email bombing and brute-force attacks.
- Invalidate existing sessions after logout or password reset and expire reset tokens immediately after use.
- Generate reset tokens with strong randomness and avoid predictable GUID/UUID schemes.
- Log password reset requests, invalid reset attempts, token generation, and reset completions.

### Nuclei

- Use Nuclei templates and YAML-based payload definitions for fast automated scanning.
- Integrate Nuclei into CI/CD pipelines and proxy workflows to detect regression vulnerabilities.
- Use Nuclei with `subfinder`, `httpx`, and other reconnaissance tools for target enrichment.
- Run Nuclei with rate limits (`-rl`) and concurrency controls (`-c`) during safe testing.
- Import Burp requests into Nuclei and use the Burp integration extension for targeted validation.
- Maintain curated template collections for HTTP, DNS, SSL, and application-layer checks.

### Login Bypass

- Enforce strict parsing of login parameters and reject malformed payloads.
- Prevent parameter pollution bypasses such as `user[]=a&pwd=b`, `user=a&pwd[]=b`, and JSON object injection like `password[password]=1`.
- Validate content type and ensure expected request formats are enforced for authentication endpoints.
- Rate limit login attempts and monitor for enumeration or SQL injection patterns.
- Do not allow unauthenticated access to protected pages or direct access to post-login endpoints.
- Test login endpoints with common bypass techniques and malformed headers, but keep errors generic.

### JWT Vulnerabilities

- Validate JWT signatures strictly and reject tokens using `alg:none` or unsupported algorithms.
- Use strong secrets and proper asymmetric key handling; do not accept RS256 tokens verified with an HS256 secret.
- Do not disclose signature verification details in error responses.
- Rotate keys regularly and avoid hard-coded or weak JWT secrets.
- Monitor JWT issuance and token replay patterns.
- Use tools like `jwt_tool` only in testing environments to verify proper JWT handling and detect weaknesses.

### Insecure Interfaces and APIs (For Cloud)

- Discover cloud-facing interfaces and APIs using specialized tools like CloudHunter, cf enum, ffuf, and S3 bucket scanners.
- Validate the target domain and limit discovery to authorized scopes before running cloud enumeration.
- Restrict API exposure and documentation to only authorized clients and authenticated users.
- Use strong access controls on cloud management interfaces and metadata services.
- Monitor cloud connections and failed discovery attempts for suspicious behavior.
- Block unauthenticated API enumeration and reduce information leakage in publicly exposed endpoints.

### File Upload

- Apply strict server-side validation on uploaded files: extension, MIME type, file contents, and size.
- Reject uploads with double extensions, null bytes, path traversal characters, whitespace, and exotic encodings.
- Normalize file names and store uploads outside of webroot with randomized names.
- Deny execution permissions on uploaded files and serve them through a safe access layer.
- Validate image metadata and inspect file contents with a secure file-type detection library.
- Monitor upload activity for high-volume or suspicious file patterns.
- Enforce content policies on uploaded files and reject executable or archive formats unless explicitly required.
- Protect against archive extraction vulnerabilities, zip slip, and malicious compressed payloads.

### File Inclusion

- Avoid directly including user-controlled paths in file operations.
- Implement allow-lists for file names and canonicalize paths before resolving them.
- Block remote file inclusion by disabling URL wrappers and validating local file inclusion paths.
- Detect and reject path traversal payloads, encoded directory separators, and manipulation via null bytes.
- Restrict file inclusion operations to known safe directories and do not include uploaded content directly.
- Monitor for LFI/RFI patterns and ensure sensitive files such as logs, configuration files, and session stores are inaccessible.
- Use path normalization and security checks for all file-related query parameters.

### Email Verification Bypass (NodeJS)

- Bind verification tokens to a specific email address and a user ID, not just a valid token payload.
- When generating the verification token, include `userId`, `email`, `issuedAt`, and a one-time nonce.
- On `/verify-email`, validate:
  - token signature and expiration
  - token email matches the current pending email for that user
  - token userId matches the authenticated user or the stored verification session
  - token has not already been used
- Store pending email changes in a separate `email_verifications` record and require the token to map to that record.
- Reject verification tokens if the user’s current pending email has changed since issuance.
- Do not permit email verification using a token generated for a prior address or prior request.
- Use short-lived confirmation tokens and audit every email update and verification event.

### CSRF (Cross-Site Request Forgery)

- Protect all state-changing endpoints (`POST`, `PUT`, `PATCH`, `DELETE`) with CSRF defenses.
- Prefer SameSite `Strict` or `Lax` for cookies and avoid using `Access-Control-Allow-Credentials` with wildcard origins.
- Require a per-session, per-form CSRF token bound to the authenticated user session.
- Verify `Origin` and `Referer` headers for sensitive requests, but do not rely on them alone.
- Reject requests when the CSRF token is missing, invalid, or not tied to the user session.
- Disable unsafe method override features unless needed, and always validate the actual effective method.
- Avoid token validation logic that skips checks when the token parameter is absent.
- Protect login and sensitive operation endpoints from same-site forgery by requiring re-authentication or password confirmation.
- Harden the application against CSRF bypasses such as:
  - GET-to-POST conversions
  - missing/optional token handling
  - session-agnostic token verification
  - custom header checks that can be spoofed
  - content-type quirks and simple request header restrictions
  - referrer/origin regex bypasses

### CSP Bypass and Hardening

- Enforce strict Content-Security-Policy rules without `unsafe-inline`, `unsafe-eval`, or overly broad wildcard sources.
- Prefer `script-src 'self'` and explicit trusted host lists; include `object-src 'none'` and `base-uri 'self'`.
- Do not allow user-uploaded files to be served under a path that can be loaded as script or HTML.
- Deny inline script execution unless using nonces or hashes with a tightly controlled policy.
- Block or restrict `data:`, `blob:`, and `filesystem:` sources unless explicitly required.
- Treat return values from JSONP endpoints and third-party scripts as untrusted.
- Validate uploaded file extensions and content types before serving them, especially when CSP relies on `'self'`.
- Monitor for CSP bypass risk factors: `unsafe-inline`, `unsafe-eval`, wildcard sources, missing `base-uri`, and JSONP endpoints.

### Container Attack Hardening

- Do not run containers in privileged mode; drop all unnecessary capabilities and enable `no-new-privileges`.
- Do not expose container runtime APIs directly to untrusted networks.
- Validate container image integrity using digest pinning and signed image provenance.
- Enforce least-privilege container configuration: read-only filesystem, minimal capabilities, and restricted network access.
- Apply CPU and memory limits to reduce DoS impact.
- Keep host and container runtime patched; regularly scan for kernel and runtime vulnerabilities.
- Use orchestration-level security controls: RBAC, Pod Security Standards, and admission policies that deny privileged pods and hostPath mounts.
- Avoid long-lived or overly powerful service account tokens; issue short-lived tokens and rotate them.
- Limit `nodes/proxy` and similar sensitive Kubernetes permissions to the minimum required verbs.

### CORS Misconfigurations and Bypass

- Never reflect the incoming `Origin` header back to `Access-Control-Allow-Origin` without strict allow-list validation.
- Do not use `Access-Control-Allow-Origin: *` together with `Access-Control-Allow-Credentials: true`.
- Only allow trusted origins explicitly and deny unknown or `null` origins for credentialed requests.
- Treat `null` origin requests as suspicious and reject them for sensitive APIs.
- Validate preflight requests and respond with exact allowed methods and headers.
- Avoid overly permissive `Access-Control-Allow-Headers` and `Access-Control-Allow-Methods` responses.
- Use CORS only for controlled cross-origin access; do not rely on it as an authentication mechanism.
- Log and monitor abnormal CORS requests, especially those using `withCredentials` or custom headers.

## Modular Design and Module Files

This design document contains the shared architecture, security hardening, infrastructure guidance, and cross-cutting middleware design for NCSMS v1.0. Detailed module-level design is split into focused files so each domain can be edited independently.

- `super-admin-v1.md` — Super Admin governance, roles, audit, and system-level operational controls.
- `admin-v1.md` — Admin operations for federation lifecycle management, licensing, grants, and business guardrails.
- `federations-v1.md` — Federation self-service workflows, document submission, accountability, and federated access controls.

Each module file contains:
- purpose and responsibility definitions
- package and handler layout
- supported endpoints and workflows
- domain models, repository operations, services, and handler routes

## Go Router and Middleware Design

### Router Layout

- `NewRouter(authMiddleware, rbacMiddleware) *chi.Mux`
- Global middleware chain for request ID, RealIP, logger, recoverer, timeout, rate limiting
- CORS and headers hardening in a reusable middleware
- Domain routes under `/api/v1`

### Middleware responsibilities

- `SupabaseAuthMiddleware` verifies JWT bearer token and injects claims into context
- `RBACMiddleware(requiredRole)` applies deterministic role checks, allowing `SYSTEM_ADMIN` bypass
- `AuditMiddleware` records request metadata and response events where needed
- `RequestNormalizationMiddleware` canonicalizes headers and payloads, strips tampering artifacts, and enforces consistent client fingerprinting
- `RateLimitMiddleware` enforces per-endpoint and per-user rate limits while preventing bypass from randomized query strings or user-agent/cookie variation
- `Timeout` middleware enforces deadlines for each HTTP request

## Code Layer Design: Entities, Repositories, Services, Handlers

### Entities / Models

- Domain models are plain Go structs with JSON tags
- Sensitive fields use explicit encryption or hashed storage wrappers
- Enumeration values enforce known role / state constants

### Repository Layer

- Uses `pgxpool.Pool`
- Implements parameterized SQL and row locking
- Returns domain entities and wrapped errors like `fmt.Errorf("select federation by id: %w", err)`
- Never ignores errors or panics on SQL failures

### Service Layer

- Contains business validation and state transition logic
- Reads via repository methods, validates business rules, writes through repository within transaction scope
- Only service-layer code mutates state and triggers audit inserts

### Handler Layer

- Translates HTTP requests into domain service calls
- Decodes JSON payloads with strict validation
- Sets request context deadlines based on Kampala timezone operations
- Returns JSON responses and consistent error messages

## Modular Detail Files

The domain-specific module designs are maintained separately in their own markdown files:

- `super-admin-v1.md`
- `admin-v1.md`
- `federations-v1.md`

These files contain the detailed per-module capabilities, package layout, endpoint design, models, repository operations, services, and handler routes.

## Compliance Implementation Notes

- All database operations go through RLS; do not assume bypass
- Sensitive fields such as NIN and medical clearance metadata should be encrypted with AES-256-GCM before storage
- Use explicit `CHECK` constraints for dates, statuses, and amounts
- Track audit changes with immutable logs and row hash metadata
- Localize all time operations to Kampala timezone for reporting and transaction context

## Docker & Kubernetes Security Appendix

### Docker Best Practices

- Use minimal base images and explicitly pin image digests.
- Avoid running containers as root and drop unnecessary Linux capabilities.
- Mount only required volumes and avoid exposing the Docker socket to untrusted containers.
- Enforce image content trust and scan images with tools like `trivy` or `grype`.
- Audit image metadata, environment variables, and mounted volumes before production deployment.
- Use Docker Compose for multi-container applications with explicit network and volume definitions.
- Minimize container privileges by using `--cap-drop=ALL` and adding only required capabilities.

### Docker Runtime Hardening

- Inspect running containers and images using `docker ps`, `docker inspect`, and `docker images --digests`.
- Review container diffs with `docker diff` and detect unauthorized filesystem changes.
- Monitor Docker events with `docker events` and track lifecycle changes.
- Use `docker bench security` to validate runtime configuration against best practices.
- Restrict Docker socket access and never mount `/var/run/docker.sock` in application containers unless absolutely required.

### Docker Attack Surface Awareness

- Validate that host Docker daemon ports (`2375`, `2376`) are not exposed publicly.
- Avoid insecure registries and use authenticated TLS-backed registry endpoints.
- Monitor for misconfigured private registries and image pull permissions.
- Ensure secrets are not stored in plaintext inside images or environment variables.

### Kubernetes Best Practices

- Enforce RBAC and least-privilege service accounts.
- Use network policies to isolate workloads and restrict pod communication.
- Scan Kubernetes manifests and container images for vulnerable software.
- Enable audit logging and monitor cluster API requests.
- Validate cluster configuration against benchmarks using tools such as `kube-bench`.

### Kubernetes Runtime and Reconnaissance

- Query cluster resources with `kubectl get nodes`, `kubectl get pods`, `kubectl get services`, and `kubectl api-resources`.
- Inspect pod details with `kubectl describe pod` and logs with `kubectl logs`.
- Identify misconfigurations by checking for exposed metadata endpoints and anonymous API access.
- Protect service account tokens and Kubernetes secrets from unauthorized read access.

### Container and Cluster Hardening Controls

- Ensure `kubectl auth can-i` checks are used before allowing privilege escalation in the cluster.
- Limit container capabilities and avoid running privileged pods.
- Secure the API server and etcd endpoints with TLS and authentication.
- Remove or restrict deprecated insecure cluster services such as unsecured kubelet or read-only pod endpoints.

## Deliverables

This document is the requested blueprint and module design for NCSMS v1.0. It can be used as the foundational architecture guide for implementing production-ready Go modules with Supabase, RLS, secure middleware, and explicit stateful business rules.

