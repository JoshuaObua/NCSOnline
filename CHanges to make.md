Here is a production-grade system prompt designed to instruct an AI or a developer to architect and implement a comprehensive, privacy-conscious Web Analytics and Visitor Statistics engine using a combined frontend/backend approach.

---

## System Prompt: Full-Stack Web Analytics & Visitor Tracking Engine

### Objective

Design and implement a robust, self-hosted **Web Analytics and Visitor Statistics Subsystem**. The solution must capture metrics from both the frontend (user interactions, screen specs) and backend (secure IP parsing, geolocation, performance) without relying on heavy third-party trackers like Google Analytics. The data must be aggregated and visualized inside an administrative dashboard.

---

### 1. Data Collection Strategy

#### Frontend Collection (Client-Side)

Implement a lightweight, non-blocking JavaScript tracker payload or API middleware that hooks into public page loads to capture:

* **Session Lifecycle:** Page entry/exit timestamps, session duration, bounce detection (e.g., active session duration under 10 seconds with no interactions).
* **Device Profiles:** User-Agent string, platform/OS (Windows, macOS, Linux, iOS, Android), browser engine, device type (Desktop, Mobile, Tablet, Smart TV), and screen resolution.
* **Referral Data:** Document referrer (`document.referrer`) to categorize traffic sources (Direct, Social Media, Search Engines, External Links).
* **Behavioral Flow:** Visited URL path names, link clicks, and time-on-page metrics.

#### Backend Collection & Enrichment (Server-Side)

Upon processing the frontend's tracking request, the backend controller must capture and securely process:

* **Network Metadata:** Request IP address (accurately parsing proxy headers like `X-Forwarded-For` or `CF-Connecting-IP` behind a load balancer/Cloudflare).
* **Geolocation Parsing:** Process the parsed IP using a reliable, self-hosted database lookup tool (e.g., MaxMind GeoIP2 Lite or an integrated GeoIP library) to resolve location down to **Country, Region/State, and City**.
* **Anonymization & Privacy Compliance:** To maintain compliance with international data privacy frameworks (like GDPR/CCPA), raw IP addresses must be immediately anonymized (e.g., hashing the IP combined with a daily rotating salt) before writing to the persistent database. Never store raw PII (Personally Identifiable Information).

---

### 2. Database Schema & High-Throughput Modeling

Because analytics tables grow rapidly, design an optimized schema separating high-frequency raw logs from pre-aggregated reporting tables:

* **`page_views_raw` (TimescaleDB / Partitioned SQL):** Captures atomic events.
* `id`, `session_id` (UUID), `visitor_hash` (anonymized signature), `path`, `referrer`, `country`, `region`, `city`, `browser`, `os`, `device_type`, `created_at`.


* **`analytics_snapshots_daily` (Aggregated Cache):** A background worker or cron job must aggregate raw data nightly into summary tables to keep admin dashboard queries lightning-fast.
* Metrics: Total Views, Unique Visitors, Bounces, Average Session Length per day, per city, and per device.



---

### 3. Core Analytical Metrics & Requirements

The analytics engine must compute and serve the following key performance indicators (KPIs) through an API endpoint:

1. **Traffic Overviews:** Total Page Views, Unique Visitors (calculated via unique visitor hashes within a 24-hour window), and Bounce Rate.
2. **Geographic Distribution:** A ranked breakdown of top-performing Countries, Regions, and **Cities** by traffic volume.
3. **Technographic Platforms:** Percentage distribution of Browsers, Operating Systems, and Device Types.
4. **Content Performance:** Top entry pages, top exit pages, and most visited content URLs.
5. **Acquisition Channels:** Breakdown of traffic sources (e.g., organic search traffic vs. direct navigation).

---

### 4. Technical Stack Expectations

* **Backend Middleware:** Implement the tracking intake endpoint in a highly performant, asynchronous language environment (such as Go, Python with FastAPI, or optimized Node.js handlers). It must return a `204 No Content` response instantly to the frontend to ensure zero impact on user experience, processing data storage asynchronously (e.g., using background workers, channel queues, or Redis).
* **Data Vis Frontend:** Build a modern, reactive analytics dashboard component (e.g., using Vue.js or React) utilizing clean charting tools (like Chart.js or D3.js) to display data timelines, bar charts for platforms, and sorted data tables for top cities/pages.
* **Performance Constraint:** Querying dashboard analytics for a 30-day time window must take less than 500ms. Utilize proper indexing on `created_at`, `path`, and geographical columns (`country`, `city`) to achieve this.