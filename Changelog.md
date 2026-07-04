## 2026-07-04 - CMS Inbound Submissions Hub

- Added the inbound submissions database migration for blog comments, contact form messages, and investment requests, including unread/read/replied/archive states, workflow status, notes, admin assignment, indexes, and permissions.
- Added backend repository and CMS handlers for listing, counting, updating, noting, and archiving inbound submissions.
- Added public investment request submission support with sanitization, honeypot handling, rate-limited routing, and admin notification dispatch.
- Connected contact messages and blog comment moderation/deletion to the inbound hub so CMS badges and workflow state stay synchronized.
- Added CMS API client methods and a dedicated Investment Requests admin section, plus refreshed labels for Contact Messages and Blog Comments.
- Added periodic CMS notification/inbound counter refresh for near real-time admin badge updates without running a build.

## 2026-07-04 - Switchable Public Form Anti-Bot Protection

- Added disabled, Cloudflare Turnstile, and Google reCAPTCHA v3 captcha settings under the CMS Contact Details settings panel.
- Added persisted captcha configuration with public secret redaction and admin validation/preservation for saved secret keys.
- Added backend captcha verification for public contact messages, investment requests, and blog comment submissions.
- Added a reusable frontend captcha helper that lazy-loads Turnstile or reCAPTCHA only when enabled.
- Wired the public blog comment form to request captcha tokens before submission.

## 2026-07-04 - Configurable Google Analytics Integration

- Added default third-party integration settings with Google Analytics disabled by default.
- Added backend validation for GA4 Measurement IDs before saving admin settings.
- Added a dedicated CMS Third-Party Integrations panel with a Google Analytics toggle and GA4 Measurement ID validation.
- Added public layout script injection/removal so Google Analytics scripts load only when enabled.
- Added SPA page-view forwarding to Google Analytics on route changes.

Here is a production-grade system prompt designed to instruct an AI or developer to implement an admin-configurable anti-bot security layout for public forms, supporting toggles for both Cloudflare Turnstile and Google reCAPTCHA v3.

---

## System Prompt: Switchable Anti-Bot Protection Middleware (Cloudflare & reCAPTCHA)

### Objective

Design and implement a highly flexible, tenant-configurable **Anti-Bot Security Module** inside the administrative dashboard settings. This engine must allow administrators to globally toggle and choose between **Cloudflare Turnstile** and **Google reCAPTCHA (v3)** to protect all public-facing form submissions (such as contact forms, blog comments, and investment requests). When turned off, form submissions must bypass validation instantly without execution friction.

---

### 1. Database Configuration Schema

Create a persistent configuration model (key-value pair table or a dedicated single-row `security_settings` table) to maintain states seamlessly:

* **`captcha_provider`** (Enum: `none`, `cloudflare_turnstile`, `google_recaptcha`)
* **`cloudflare_site_key`** (String, public token)
* **`cloudflare_secret_key`** (String, encrypted at rest on backend)
* **`recaptcha_site_key`** (String, public token)
* **`recaptcha_secret_key`** (String, encrypted at rest on backend)
* **`recaptcha_score_threshold`** (Decimal, default `0.5`, specific to reCAPTCHA v3 risk verification)

---

### 2. Admin Settings Interface (UI/UX)

Provide an elegant, clear options array within the administrative configuration dashboard panel:

* **Protection Mode Selector:** A 3-way radio button or toggle switch group labeled:
* `[ Disabled ]`
* `[ Cloudflare Turnstile ]`
* `[ Google reCAPTCHA v3 ]`


* **Conditional Field Visibility:**
* If `Disabled` is selected, all credential inputs are completely hidden.
* If `Cloudflare Turnstile` is active, display clear text inputs for **Site Key** and **Secret Key**.
* If `Google reCAPTCHA v3` is active, display inputs for **Site Key**, **Secret Key**, and a threshold accuracy slider (ranging from `0.1` to `1.0`).


* **Action Controls:** A unified "Save Configuration" button that triggers instant backend validation checks ensuring the secret keys conform to minimal format structures before saving.

---

### 3. Frontend Component Injection Structure

Develop a clean, modular client-side injection wrapper that wraps all public-facing interactive submission forms:

* **Dynamic Initialization:** On page load, fetch the active public provider flag (`captcha_provider`) and public `site_key`.
* **State Execution:**
* **Case `none`:** Do absolutely nothing. Render standard inputs; submit form directly via AJAX or POST.
* **Case `cloudflare_turnstile`:** Inject the Turnstile script library asynchronously. Render the implicit widget `<div>` safely inside the form container. Append the generated token parameter (`cf-turnstile-response`) to the request payload upon click.
* **Case `google_recaptcha`:** Load the reCAPTCHA v3 API library. Execute `grecaptcha.execute()` invisible token requests on form submission intent, automatically binding the generated execution token to the form request headers or body payload.



---

### 4. Backend Validation Middleware Engine

Build a secure, reusable request validation middleware layer that intercepts incoming form submissions before processing any business database logic:

* **Bypass Check:** Evaluate the configuration flag. If `captcha_provider` equals `none`, bypass immediately and forward the data execution safely to the next functional controller block.
* **External Verification Handlers:**
* **Cloudflare Turnstile Logic:** Fire an internal backend HTTP POST request directly to `[https://challenges.cloudflare.com/turnstile/v0/siteverify](https://challenges.cloudflare.com/turnstile/v0/siteverify)` containing the client-provided token and your backend `cloudflare_secret_key`.
* **Google reCAPTCHA Logic:** Fire an internal backend HTTP POST request directly to `[https://www.google.com/recaptcha/api/siteverify](https://www.google.com/recaptcha/api/siteverify)` passing the payload token along with your `recaptcha_secret_key`. Parse the returning JSON response body to confirm `success: true` and verify that the risk score matches or exceeds the configured `recaptcha_score_threshold`.


* **Fail-Safe Interception Handling:** If verification checks fail, stop the execution trace immediately. Drop the request pipeline, write a localized security alert warning log, and respond to the client interface with a standardized `422 Unprocessable Entity` or `400 Bad Request` API error envelope containing a clear user-facing message (e.g., *"Bot verification failed. Please try again."*).

---

### 5. Architectural Security Demands

* **Timeout Mitigation:** Implement a strict HTTP client deadline (maximum `3` seconds) for outbound API network requests to Google or Cloudflare. If their third-party authentication endpoints suffer localized downtime, the system must handle the exception gracefully without stalling server worker resources.
* **Zero Hardcoding:** Absolutely no keys can be written statically in codebase templates. All active values must flow securely from the administration's persistent database configurations or environment matrices.
