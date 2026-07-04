Here is a production-grade system prompt designed to instruct an AI or developer to architect and implement a comprehensive **Notification Engine & Administrative Inbound Management System**.

---

## System Prompt: Real-Time Event Notification Engine & Administrative Submissions Hub

### Objective

Design and implement a centralized **Notification Engine** paired with a dedicated **Inbound Submissions Hub** in the administrative sidebar dashboard. The system must process real-time system events (security, backups, updates) and user-generated content (comments, contact forms, financial/investment requests), providing administrators with complete management controls (read, unread, reply states).

---

### 1. Unified Notification Engine (System Events)

Create a centralized notification dispatcher capable of routing high-priority system events to the admin dashboard via real-time WebSockets/Server-Sent Events (SSE) and persisting them to a `system_notifications` database table.

Capture, format, and alert administrators immediately on the following triggers:

* **Authentication Security Alerts:** Every time an account with administrative or elevated privileges signs in. The notification must include the timestamp, user ID, IP address, and geographical location (City/Country).
* **Smart Update Alerts:** When the automated Git deployer detects a new commit upstream on the remote repository.
* **Database Operation Success:** Immediately upon the successful generation and verification of a database schema or data backup file.
* **System Vulnerability/Failures:** If a background deployment script fails, or if a critical container service goes offline.

---

### 2. Administrative Submissions Hub (Sidebar Management)

Extend the administrative sidebar wrapper with dedicated navigation menu items to monitor, filter, and triage public-facing form submissions. Each submission category must have its own isolated management view:

#### A. Blog Comments Manager

* **Trigger:** A visitor submits a comment on any public blog post.
* **Sidebar Menu Item:** `Blog Comments` (with a dynamic, real-time unread counter badge).
* **Management Actions:** View comment text, see associated blog post, **Approve for Public Display**, **Spam/Trash**, and **Delete**.

#### B. Contact Inquiries Manager

* **Trigger:** A visitor submits a form via the public "Contact Us" page.
* **Sidebar Menu Item:** `Contact Messages` (with unread badge).
* **Management Actions:** View sender details (Name, Email, Phone, Message body), **Mark as Read/Unread**, **Mark as Replied**, and a quick-text inline reply markdown editor that fires a secure transactional email back to the visitor.

#### C. Investment Requests Manager

* **Trigger:** A user submits a financial, funding, or investment proposal form via the public investment portals.
* **Sidebar Menu Item:** `Investment Requests` (with high-priority visual indicator).
* **Management Actions:** View submission parameters (Requested amounts, investor profile details, attached files), **Update Status Dropdown** (*Pending Review, Approved, Under Negotiation, Declined*), and **Assign Internal Staff Note Ledger** for internal administrative audit tracking.

---

### 3. Database Schema & State Management

Ensure the data structures supporting these workflows are optimized for read/write state toggles:

* **`inbound_submissions` Schema:**
* `id` (UUID)
* `source_type` (Enum: `blog_comment`, `contact_form`, `investment_request`)
* `payload` (JSONB block containing flexible form fields unique to each type)
* `status_state` (Enum: `unread`, `read`, `replied`, `archived`)
* `assigned_admin_id` (Nullable foreign key)
* `created_at` / `updated_at`


* **Indexing:** Apply composite indexes on `(source_type, status_state)` and B-Tree indexes on `created_at` to keep sidebar badge counters and dashboard lists highly performant.

---

### 4. Technical & UI/UX Expectations

* **Real-Time Badge Syncing:** Use global application state management on the frontend dashboard to update unread counts instantly without requiring the administrator to perform a manual page refresh.
* **Backend Modularity:** Implement a clean Observer or Event-Listener pattern on the backend. When a form controller successfully validates a submission, it must synchronously dispatch an event to the notification listener, ensuring decoupling between core domain logic and notification delivery.
* **Security Controls:** Enforce strict middleware checks ensuring that only authenticated administrators can mutate submission states (`mark as read`, `delete`, `change investment status`). Safely sanitize all incoming HTML payload strings to prevent Cross-Site Scripting (XSS) injections within the admin dashboard panels.