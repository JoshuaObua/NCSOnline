# Comprehensive Specification: Standard Fields & Architecture for Professional Audit Logging

A professional-grade enterprise audit logging system must record deterministic, tamper-evident, and context-rich data for every event within an application ecosystem. The logs must provide an unambiguous trail of **who** performed the action, **what** the action was, **when** it occurred, **where** it originated, and **how** it was executed.

---

## 1. Core Structural Categories

A professional audit log record is classified into six fundamental pillars:
1. **Actor / Identity Data** (Who did it?)
2. **Event / Action Data** (What did they do?)
3. **Temporal Data** (When did it happen?)
4. **Context & Environment Data** (From where and via what interface?)
5. **Target / Resource Data** (To what asset was it done?)
6. **Security & Integrity Data** (How is the log validated and secured?)

---

## 2. Exhaustive Field Reference & Schema Matrix

| Field Category | Database Field Name | Data Type | Requirement | Description / Standard Values | Example |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Temporal** | `timestamp` | ISO 8601 UTC String | **Mandatory** | Precise timestamp down to millisecond or microsecond resolution. | `2026-07-19T00:58:37.123Z` |
| **Temporal** | `timezone` | String (IANA) | Optional | The local timezone context of the actor or originating node. | `Africa/Kampala` |
| **Event** | `event_id` | UUIDv4 / ULID | **Mandatory** | Universally unique identifier for the specific log entry. | `9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d` |
| **Event** | `correlation_id` | UUIDv4 / ULID | **Mandatory** | Traverses microservices to link asynchronous or multi-step operations. | `c3a8b410-1289-4f7d-a111-923ad4f392bb` |
| **Event** | `action` | String (Enum) | **Mandatory** | Explicit action standardizing the operation. Typically `object:action`. | `user.profile:update`, `order:delete` |
| **Event** | `status` | String (Enum) | **Mandatory** | The outcome of the attempted transaction. | `SUCCESS`, `FAILED`, `DENIED` |
| **Event** | `severity` | String (Enum) | **Mandatory** | Risk or operations classification level. | `INFO`, `WARN`, `CRITICAL` |
| **Actor** | `actor_id` | String / UUID | **Mandatory** | Unique identifier of the user or system principal initiating the event. | `usr_8f3a92c10b` |
| **Actor** | `actor_type` | String (Enum) | **Mandatory** | Type of identity domain executing the process. | `USER`, `SYSTEM_JOB`, `SUPPORT_IMPERSONATOR` |
| **Actor** | `actor_username` | String | Highly Rec. | Human-readable identifier/email captured at event runtime. | `admin@company.com` |
| **Actor** | `actor_roles` | Array (Strings) | Recommended | Authorization snapshot of the user's roles when the event occurred. | `["SystemAdmin", "ComplianceAuditor"]` |
| **Context** | `client_ip` | String (IPv4/IPv6)| **Mandatory** | Originating network address, resolving proxies using `X-Forwarded-For`. | `197.239.5.14` |
| **Context** | `access_medium` | String (Enum) | **Mandatory** | The broad architectural layer used to interface with the system. | `WEB_UI`, `MOBILE_APP`, `REST_API`, `CLI` |
| **Context** | `user_agent_raw` | String | **Mandatory** | The exact HTTP User-Agent string received by the server boundary. | `Mozilla/5.0 (Macintosh; Intel Mac...)` |
| **Context** | `parsed_client_agent`| String | Recommended | Normalized client identification parsing engine/tool details. | `PostmanRuntime/7.39.0`, `Insomnia/9.2.0` |
| **Context** | `parsed_os` | String | Recommended | Extracted Operating System from metadata layers. | `iOS 17.4`, `Ubuntu 24.04`, `Windows 11` |
| **Context** | `parsed_browser` | String | Recommended | Extracted Browser name and engine version. | `Chrome 126.0.0`, `Safari 17.5` |
| **Context** | `request_url` | String | Recommended | Full API or web endpoint URI targets. | `https://api.platform.com/v1/users` |
| **Context** | `http_method` | String (Enum) | Recommended | HTTP Verb used if accessed via web/API vectors. | `POST`, `GET`, `PUT`, `DELETE`, `PATCH` |
| **Target** | `resource_id` | String / UUID | **Mandatory** | Unique ID of the target domain object modified or read. | `acc_44910294-b` |
| **Target** | `resource_type` | String | **Mandatory** | Entity classification of the object targeted. | `BillingAccount`, `DatabaseCluster` |
| **Target** | `payload_before` | JSON / Encrypted | Recommended | Snapshot state *prior* to changes (Omit/Mask PII/Secrets). | `{"status": "active", "tier": "gold"}` |
| **Target** | `payload_after` | JSON / Encrypted | Recommended | Snapshot state *after* changes (Omit/Mask PII/Secrets). | `{"status": "suspended", "tier": "gold"}` |
| **Security** | `signature` | String (Hex/B64) | Highly Rec. | Cryptographic HMAC or asymmetric signature to verify log immutability. | `e3b0c44298fc1c149afbf4c8996fb92427ae...` |

---

## 3. Deep-Dive: Access Mediums & Agent Taxonomy

A professional audit logging pipeline must distinctly parse and normalize the **Access Medium** (`access_medium`) and its corresponding **Agent** details (`parsed_client_agent`, `user_agent_raw`). Below is the architectural taxonomy mapping how various vectors must be tracked:

### A. REST / GraphQL / gRPC API Access
*   **Access Medium (`access_medium`):** `REST_API` or `GRAPHQL_API`
*   **Differentiating Agents:**
    *   **Postman:** Captured in `user_agent_raw` as `PostmanRuntime/x.x.x`. Logged under `parsed_client_agent` explicitly as `Postman`.
    *   **Curl / CLI Tools:** Captured as `curl/7.81.0` or `HTTPie/3.2.2`. Helps identify raw developer scripts and DevOps executions.
    *   **SDKs / Internal Code:** Automated microservice daemons using language libraries (e.g., `boto3/1.34.0 Python/3.11`, `Go-http-client/1.1`).
*   **Essential Meta:** Tracking `api_key_id` or `oauth_client_id` alongside `actor_id` to establish exactly which credential authenticated the runtime script.

### B. Web Browser Interfacing
*   **Access Medium (`access_medium`):** `WEB_UI`
*   **Differentiating Agents:**
    *   **Browsers:** `Chrome`, `Firefox`, `Safari`, `Edge`, `Brave`.
*   **Essential Meta:** Must track standard session cookies or JSON Web Token (JWT) session IDs (`session_id`) to correlate web interface activities with concurrent UI state changes.

### C. Native & Hybrid Mobile Applications
*   **Access Medium (`access_medium`):** `MOBILE_APP`
*   **Differentiating Agents:**
    *   **Platforms:** Android App Network Clients, iOS CFNetwork Clients.
*   **Essential Meta:** System should parse device specific markers into `parsed_os` (e.g., `Android 14; Build/UKQ1`) and collect application bundle versions (e.g., `com.company.app/v4.2.1`).

---

## 4. Architectural Compliance & Security Guardrails

To prevent log evasion, tamper attempts, and security degradation, audit trails must implement strict design boundaries:

1. **Write-Once-Read-Many (WORM):** Logs must be shipped instantly via asynchronous streaming pipelines (such as Apache Kafka or AWS Kinesis) to isolated log-aggregators (like OpenSearch, Datadog, or AWS CloudWatch) backed by storage layers enforcing immutable object locks.
2. **Never Log Sensitive Data:** High-security filters must scrub fields for plaintext credentials, credit card primary account numbers (PAN), bearer tokens, encryption keys, and highly protected health information (PHI) before serialization.
3. **Cryptographic Chaining:** Advanced systems structure logs using hash chains, where each entry contains a hash of the preceding entry (`prev_log_hash`), turning logs into verifiable append-only data frameworks.
