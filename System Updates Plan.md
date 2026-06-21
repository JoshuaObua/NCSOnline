# NCS-Online — System Updates & Hardening Plan

A consolidated architectural blueprint for turning the NCS-Online VPS deployment into a resilient, enterprise-grade platform with operator "flight deck" controls. The backend is **Go** (highly concurrent), the frontend is **Vue 3 + Tailwind**, the database is **PostgreSQL**, and everything runs in **Docker Compose** on a self-hosted VPS.

The guiding principles throughout:

- **Lightweight & thread-safe.** Every operator control must be decoupled from request-path business logic. Use `atomic.*`, `sync.RWMutex`, buffered channels, and worker pools — never block on disk or network in a middleware.
- **Defense in depth.** Geographic filter at the edge, RBAC inside the app, immutable audit at rest.
- **No 3 AM panics.** Every destructive action has a backup, a dry-run, and a rollback. Every deploy is reversible within seconds.
- **Operator > Developer.** If the operator needs SSH to do something routine, the design has failed.

---

## Table of Contents

1. [Layered Architecture Overview](#1-layered-architecture-overview)
2. [Edge Layer — Python Geo-Sentinel](#2-edge-layer--python-geo-sentinel)
3. [Application Layer — Go Audit & Security Middleware](#3-application-layer--go-audit--security-middleware)
4. [Maintenance Core](#4-maintenance-core)
5. [Logging & Export Engine](#5-logging--export-engine)
6. [Database Backup & Restore](#6-database-backup--restore)
7. [CI/CD & Smart Update Pipeline](#7-cicd--smart-update-pipeline)
8. [Multi-Channel Notification System](#8-multi-channel-notification-system)
9. [Admin CMS — Navigation, Branding & SEO](#9-admin-cms--navigation-branding--seo)
10. [User Portal — `/my-portal`](#10-user-portal--my-portal)
11. [Go Implementation Blueprints](#11-go-implementation-blueprints)
12. [Operational Flow — End-to-End Upgrade Walkthrough](#12-operational-flow--end-to-end-upgrade-walkthrough)
13. [Gaps & Recommended Additions](#13-gaps--recommended-additions)

---

## 1. Layered Architecture Overview

Position components sequentially so malicious traffic never reaches Go execution memory:

```
                                Internet
                                    │
                                    ▼
              ┌──────────────────────────────────────┐
              │  L1: Nginx + Python Geo-Sentinel     │  ← drop non-UG, rate limit, TLS
              └──────────────────────────────────────┘
                                    │
                                    ▼
              ┌──────────────────────────────────────┐
              │  L2: Go Audit / RBAC Middleware      │  ← log, fingerprint, authz
              └──────────────────────────────────────┘
                                    │
                                    ▼
              ┌──────────────────────────────────────┐
              │  L3: Application Handlers            │  ← business logic
              └──────────────────────────────────────┘
                                    │
                ┌───────────────────┼───────────────────┐
                ▼                   ▼                   ▼
           PostgreSQL           Redis/Cache         Object Storage
                                                    (backups, uploads)
```

Side processes:

- **Update Sentinel** — background goroutine polling GitHub releases.
- **Notification Worker Pool** — buffered channel dispatching FCM / Telegram / SMTP.
- **Backup Cron** — `pg_dump` → compressed → off-site sync.
- **Resource Monitor** — `gopsutil` streaming to admin WebSocket.

---

## 2. Edge Layer — Python Geo-Sentinel

A FastAPI/ASGI reverse-proxy gatekeeper that screens IPs before they touch Go memory.

### Responsibilities
- Resolve client country via **local MaxMind GeoLite2** (weekly auto-updated). No runtime external lookups.
- Safely traverse upstream headers (`CF-Connecting-IP`, `X-Forwarded-For`, `X-Real-IP`) to find the true origin IP — never trust the first hop blindly.
- Enforce `ISO == UG` unless on a whitelist (GitHub webhook IPs, telecom USSD gateway CIDRs, admin static IPs).

### Defensive Actions
| Traffic Type           | Action on Non-UG / Anomaly                                                  |
| ---------------------- | --------------------------------------------------------------------------- |
| Browser (HTML)         | `301` to `https://www.google.com/` — wastes scanner cycles, hides the stack |
| Mobile / API (JSON)    | `403 Forbidden`, no body                                                    |
| USSD                   | `END Secure System Violation` (telecom-compliant termination)               |
| Botnet / repeat offender | TCP drop, zero bytes; fail2ban-style temp ban                             |

### Operator Controls
- Whitelist editor in Admin UI (CIDR, label, expiry).
- Live counter of dropped connections per country, per hour.
- Manual GeoIP DB refresh button.

---

## 3. Application Layer — Go Audit & Security Middleware

Universal request interceptor across REST, USSD, Web, and Mobile — runs *after* the geo-sentinel passes.

### Forensic Capture (structured `slog` JSON)
- **Request:** UTC timestamp, request UUID, method, URI, protocol, TLS state.
- **Client:** resolved IP, user-agent, derived platform, browser signature.
- **Identity:** user ID, role (`Guest` / `User` / `OrgAdmin` / `SystemAdmin`), session hash.
- **Payload:** raw body up to 2 KB (PII-scrubbed via field allowlist).
- **USSD telemetry:** MSISDN, SessionId, ServiceCode.

### Anomaly Rules
1. **Structural fingerprinting.** Mobile UA without expected SDK header → flag. Browser hitting `/api/admin/*` without referrer → flag.
2. **Privilege boundary violation.** Non-`SystemAdmin` hitting maintenance routes → instant `403` + audit event tagged `security.boundary.violation`.
3. **Honeypot redirects.** Suspected scanner → `307` to a static blank asset to exhaust their fetcher.

### Audit Log Hardening
- **Append-only.** No `UPDATE` or `DELETE` permitted on the audit table at the DB role level.
- **Hash-chained.** Each row stores `prev_hash = SHA256(prev_row || self_fields)`. Tampering breaks the chain.
- **Async writes.** Buffered channel → worker pool → batch insert every 250 ms. Zero latency added to request path.
- **Retention.** 90 days hot in PG; older entries shipped to compressed cold storage.

---

## 4. Maintenance Core

### Maintenance Mode Toggle

A single `atomic.Bool` (or `atomic.Pointer[State]` if multi-field) flipped from the Admin UI.

| Traffic        | Response in Maintenance Mode                                              |
| -------------- | ------------------------------------------------------------------------- |
| Public Web     | `503` HTML page with branded "We'll be right back" template               |
| REST API       | `503` `{"status":503,"message":"System undergoing scheduled maintenance"}` |
| USSD           | `END System is temporarily undergoing maintenance. Please try again later.` |
| Admin UI/API   | **Bypass** — `SystemAdmin` role OR `/admin/*` path                        |

### Cache & Session Tools
- **Flush application cache** — one-click button.
- **Force logout all users** — bumps a global `session_epoch` counter so every JWT issued before *now* fails validation.
- **Per-user session revocation** — list of active devices with individual revoke buttons.

### Queue Worker & Cron Manager
- Live table of background workers (name, last run, last duration, success/fail count, next scheduled).
- Pause / resume / trigger-now buttons.
- Dead-letter queue browser with retry/discard actions.

### Live Resource Monitor (`shirou/gopsutil` → WebSocket)
- CPU %, RAM, swap, load average.
- Disk usage per mount (warn at 80 %, alarm at 90 %).
- Active connections, goroutine count, heap size — spot leaks early.
- Per-container stats via Docker socket.

---

## 5. Logging & Export Engine

Use `log/slog` writing to **stdout + rotating files** (`/var/log/app/*.log`, rotated daily, gzipped, 30-day retention).

### Export Pipeline (on-demand from Admin UI)
| Format | Implementation                                                              |
| ------ | --------------------------------------------------------------------------- |
| TXT    | Stream raw file or one-line-per-event JSON                                  |
| MD     | Render with code-fenced JSON blocks per event                               |
| XML    | Marshal into `<audit_logs><event>…</event></audit_logs>` schema             |
| PDF    | `johnfercher/maroto` — header, timestamps, severity colors, page numbers    |
| CSV    | Flat columns for spreadsheet pivoting                                       |

### UI Filters (applied before export)
- Date range, severity (`debug`/`info`/`warn`/`error`/`fatal`), user, IP, event name.
- "Failed auth attempts only" quick filter.
- Estimated row count + file size before triggering download.

---

## 6. Database Backup & Restore

Never run dumps from inside Go memory — shell out to native tooling.

### Backup Pipeline
```
pg_dump --format=custom --compress=9 ncs_prod
   │
   ▼
gzip → /var/backups/ncs/ncs-YYYYMMDD-HHMMSS.dump.gz
   │
   ▼
rclone copy → off-site (S3 / Backblaze / Supabase Storage)
   │
   ▼
SHA-256 manifest + retention prune (keep 7 daily, 4 weekly, 12 monthly)
```

### Restore UI
1. Table of available backups (local + off-site) with size, age, checksum status.
2. **Verify** button — restore into a transient `ncs_verify` DB and run a row-count diff. *(Untested backups are not backups.)*
3. **Restore** button — forces Maintenance Mode → snapshots current state → streams dump back → clears cache → exits Maintenance Mode.
4. **Scheduled drill** — weekly automated verify-only restore with email report.

---

## 7. CI/CD & Smart Update Pipeline

### GitHub Actions CI (`.github/workflows/ci.yml`)
On every push / PR to `master`:
1. Spin up clean Ubuntu runner.
2. Cache Go modules + Docker layers.
3. `go vet`, `golangci-lint`, `govulncheck`.
4. `go test ./...` against an ephemeral PG service container.
5. Build container image, sign with **cosign**, generate SBOM.
6. Tag image as `ncs-online:<sha>` and `ncs-online:latest` only on green.

A **failed CI blocks the deploy button** in the Admin UI.

### VPS Deploy Key Provisioning (Admin UI)
- **"Generate Deploy Key"** → Go calls `ssh-keygen -t ed25519` in memory.
- **Public key:** rendered in read-only textarea with copy button + step-by-step "paste into GitHub → Settings → Deploy Keys (read-only)" instructions.
- **Private key:** written to `/root/.ssh/id_ed25519_ncs` with `chmod 600`. **Never** shown in the UI.
- **Rotate** button regenerates and forces the operator to re-paste.

### Smart Upgrade Workflow (`Pull & Run Updates` button)
1. Pre-flight: green CI? disk space? DB reachable? → else abort.
2. **Maintenance Mode ON.**
3. `git fetch` via the deploy key; verify the incoming tag matches the latest signed release.
4. Snapshot the running binary → `app.bak`.
5. `git pull` → `go build -o app_new` → smoke-test `app_new --healthcheck`.
6. **Atomic swap**: `mv app_new app` → send `SIGTERM` to old process → systemd/supervisord restarts.
7. New process must respond `200` on `/readyz` within 30 s.
8. **Failure path:** restore `app.bak`, restart, page operator via Telegram, stay in Maintenance Mode for manual review.
9. **Success path:** clear cache, run pending DB migrations (forward-only, lock-timeout 5s), Maintenance Mode OFF.

### Live Deploy Console
Stream `stdout`/`stderr` of every step over WebSocket to the Admin UI terminal pane. Operator sees `git pull` output, `go build` warnings, and health-check probes in real time.

### Preserving User Data
Uploads live at `/var/www/ncs/uploads/` — **outside** the binary directory. Upgrade scripts must `--exclude` this path. Verified by a CI test that fails if any deploy script touches that prefix.

---

## 8. Multi-Channel Notification System

### Device Token Registry
Table `user_device_tokens`:
| column            | type        | notes                                     |
| ----------------- | ----------- | ----------------------------------------- |
| user_id           | uuid        | FK                                        |
| device_token      | text        | FCM / APNs / Web Push endpoint            |
| platform          | enum        | `android` / `ios` / `web_desktop`         |
| browser_signature | text        | nullable, web only                        |
| last_seen_at      | timestamptz | for hygiene                               |
| created_at        | timestamptz |                                           |

**Token hygiene:** on logout, on FCM `Unregistered` / `InvalidRegistration` response, or after 90 days of `last_seen_at` inactivity → delete row.

### Worker-Pool Dispatcher
- Buffered channel `notifyCh chan NotifyJob` with N goroutines.
- Each job: `{user_id, event_type, payload}`.
- Worker fans out per platform with platform-correct payload shape.

### Channel Matrix
| Event                  | Telegram (ops) | Email (user) | FCM mobile | Web Push | In-App Banner |
| ---------------------- | -------------- | ------------ | ---------- | -------- | ------------- |
| GitHub update available | ✅            | —            | —          | —        | ✅ (admins)   |
| Backup completed       | ✅             | —            | —          | —        | ✅            |
| Backup FAILED          | ✅ (urgent)    | ✅           | ✅         | ✅       | ✅            |
| Auth from new IP       | —              | ✅           | ✅         | ✅       | —             |
| PIN changed            | —              | ✅           | ✅         | ✅       | —             |
| Maintenance starting   | ✅             | ✅ (users)   | ✅         | ✅       | ✅            |

### PWA Service Worker (`/my-portal/sw.js`)
```javascript
self.addEventListener('push', (event) => {
  const data = event.data.json();
  event.waitUntil(self.registration.showNotification(data.title, {
    body:    data.body,
    icon:    '/storage/branding/favicon.png',
    badge:   '/storage/branding/badge_icon.png',
    data:    { url: data.target_url },
    vibrate: [200, 100, 200],
  }));
});

self.addEventListener('notificationclick', (event) => {
  event.notification.close();
  event.waitUntil(clients.openWindow(event.notification.data.url));
});
```

### Sanitized Payloads
Notifications carry only **public** release data (tag, name, release notes, URL). **Never** server paths, DB DSNs, env vars, or internal IPs. Reviewed by a unit test asserting `payload.contains(forbidden_substrings) == false`.

---

## 9. Admin CMS — Navigation, Branding & SEO

### Collapsible Sidebar Taxonomy
Replace flat tabs with grouped accordions. Backend serves the tree as JSON, **RBAC-filtered** (an Editor never receives `Maintenance Core` or `Security Analytics` nodes).

- **📝 Blogs Manager**
  - Create New Post · Manage Posts · Categories · Comments & Moderation
- **🎨 Branding Controls**
  - Logo Configurations · SEO & Global Metadata · Sitemap Management
- **⚙️ Maintenance Core**
  - System Status & Mode · Debug & Diagnostic Logs · DB Backup & Restore · Smart Updates & Rollbacks
- **🛡️ Security Analytics**
  - User Audit Trails · System Resource Monitor · Active Sessions · Geo-Sentinel Log

### Interaction Spec
- Chevron `▸` rotates to `▾` on expand (`transition-transform duration-200`).
- Children animate via `max-height` transition (no layout shifts).
- Active-route detection auto-opens the matching parent on refresh.
- Subtle left border + indent anchors children visually to parents.

### Asset Management (5 brand slots)
1. **Favicon** (`.ico` / `.png` / `.svg`)
2. **Main Logo — Light Mode**
3. **Main Logo — Dark Mode**
4. **Footer Logo — Light**
5. **Footer Logo — Dark**

Each upload: MIME validated, path-traversal blocked, auto-resized, written to `/var/www/ncs/storage/branding/`. Served via a memory-cached resolver; cache invalidates on save.

### SEO Dashboard
- **Global meta:** title, description (160-char counter), keywords, `robots`.
- **Open Graph:** `og:title`, `og:description`, `og:image` (uploader).
- **Schema.org JSON-LD:** raw editor with live validator.
- **Per-route overrides:** key/value table addressed by route pattern.

### Sitemap Engine
- Cron-driven goroutine walks the route registry + dynamic content tables.
- Writes `sitemap.xml` to memory cache.
- Served at `/sitemap.xml` with `Content-Type: application/xml`.
- Admin card: last-generated timestamp, URL count, file size, **Regenerate Now** button.

---

## 10. User Portal — `/my-portal`

### Dynamic Sidebar (RBAC + entitlement-driven)
1. **Active Services** — looped from backend, each with status dot, icon, deep link.
2. **Available Applications** — open enrollment tracks.
3. **Account & Management** — Profile, Security, Activity & Audit Logs, Notification Preferences.

Closed/restricted services are dropped from the payload — never rendered with a "disabled" state.

### Profile (`/my-portal/profile`)
Avatar (upload + crop), full name, account number, tier, registered date. Editable: name, email (verification flow), phone, address, language, timezone.

### Security (`/my-portal/security`)
- **Password change** — regex: ≥12 chars, upper + lower + digit + symbol; current-password re-auth required.
- **Transactional PIN** — 4–6 digits, masked pin-pad UI, used for USSD-linked actions.
- **TOTP 2FA** — QR + manual string, 6-digit confirm to activate. Backup codes generated and shown once.
- **SMS/Email 2FA fallback** — toggle to route step-up tokens via registered channel.
- **Active sessions** — list with device, IP, country, last seen + individual + "sign out all others" buttons.

### Audit Logs (`/my-portal/logs`)
Paginated table fed by the Go audit logger: timestamp, event, resource, IP, country (GeoIP), device, success/fail. Failed-auth rows highlighted soft-red. One-click "Sign out other sessions".

### Notification Preferences (`/my-portal/notifications`)
Matrix of toggles by category × channel. Security alerts are **mandatory email + push** (toggle disabled, with tooltip explaining why).

### Engineering Guardrails
- Re-auth modal intercepts every sensitive mutation (password / PIN / 2FA).
- Session cookies rotated on any credential change.
- All forms ship with CSRF double-submit tokens.

---

## 11. Go Implementation Blueprints

### Update Sentinel + Telegram Dispatcher
```go
package maintenance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type GitHubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
}

type SystemState struct {
	mu              sync.RWMutex
	CurrentVersion  string
	LatestVersion   string
	UpdateAvailable bool
	MaintenanceMode atomic.Bool // hot path — no mutex
}

var GlobalState = &SystemState{CurrentVersion: "v1.0.0"}

func StartUpdateSentinel(interval time.Duration, owner, repo string) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			release, err := fetchLatestRelease(owner, repo)
			if err != nil {
				slog.Error("sentinel fetch failed", "err", err)
				continue
			}
			GlobalState.mu.Lock()
			isNew := release.TagName != GlobalState.CurrentVersion &&
				release.TagName != GlobalState.LatestVersion
			if isNew {
				GlobalState.LatestVersion = release.TagName
				GlobalState.UpdateAvailable = true
			}
			GlobalState.mu.Unlock()
			if isNew {
				slog.Info("new release detected", "tag", release.TagName)
				go dispatchTelegramAlert(release)
			}
		}
	}()
}

func fetchLatestRelease(owner, repo string) (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "NCS-Update-Sentinel")
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github status %d", resp.StatusCode)
	}
	var r GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func dispatchTelegramAlert(r *GitHubRelease) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_ADMIN_CHAT_ID")
	if token == "" || chatID == "" {
		slog.Warn("telegram credentials missing — skipping alert")
		return
	}

	msg := fmt.Sprintf(
		"⚠️ *SYSTEM UPDATE AVAILABLE*\n\n"+
			"*Current:* `%s`\n*New:* `%s`\n*Release:* %s\n\n"+
			"Log into the Admin Maintenance Panel to back up and run Smart Upgrade.",
		GlobalState.CurrentVersion, r.TagName, r.Name,
	)

	body, _ := json.Marshal(map[string]string{
		"chat_id":    chatID,
		"text":       msg,
		"parse_mode": "Markdown",
	})

	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token),
		"application/json", bytes.NewBuffer(body),
	)
	if err != nil {
		slog.Error("telegram dispatch failed", "err", err)
		return
	}
	resp.Body.Close()
}
```

### Maintenance Middleware
```go
func MaintenanceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !GlobalState.MaintenanceMode.Load() {
			next.ServeHTTP(w, r)
			return
		}
		// Bypass: admin paths or SystemAdmin role
		if strings.HasPrefix(r.URL.Path, "/admin") || HasRole(r, "SystemAdmin") {
			next.ServeHTTP(w, r)
			return
		}
		switch detectChannel(r) {
		case channelUSSD:
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "END System undergoing maintenance. Please try again later.")
		case channelAPI:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]any{
				"status":  503,
				"message": "System undergoing scheduled maintenance.",
			})
		default:
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusServiceUnavailable)
			renderMaintenancePage(w)
		}
	})
}
```

### Admin Status Handler
```go
func AdminStatusHandler(w http.ResponseWriter, r *http.Request) {
	GlobalState.mu.RLock()
	defer GlobalState.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"current_version":  GlobalState.CurrentVersion,
		"latest_version":   GlobalState.LatestVersion,
		"update_available": GlobalState.UpdateAvailable,
		"maintenance_mode": GlobalState.MaintenanceMode.Load(),
	})
}
```

---

## 12. Operational Flow — End-to-End Upgrade Walkthrough

1. **Alert** — Sentinel detects a new GitHub release, pings the operator's Telegram.
2. **Assess** — Operator logs into Admin UI, sees the glowing update badge in the header.
3. **Prepare** — Navigates to **Maintenance → Backups**, clicks **Snapshot Now**, waits for the green checkmark + checksum verification.
4. **Deploy** — Goes to **Smart Updates**, picks the new tag, clicks **Run Upgrade**.
5. **Watch** — Live console streams `git pull` → `go build` → health probes.
6. **Verify** — On success, badge clears, `/readyz` returns 200, Telegram posts the confirmation.
7. **Recover (if needed)** — On failure, atomic rollback to `app.bak` runs automatically; system stays in Maintenance Mode and pages the operator for review.

---

## 13. Gaps & Recommended Additions

The original draft is strong on day-2 operations but missing several production essentials. Recommend adding:

### Security
- **Secrets management.** Rotate `JWT_SECRET`, DB passwords, Telegram token, FCM key on a schedule. Consider `age`-encrypted env files committed to a private repo, or a self-hosted Vault.
- **Security headers middleware** — HSTS, CSP, X-Frame-Options, Referrer-Policy, Permissions-Policy.
- **Rate limiting** — per-IP and per-user-ID token buckets; stricter limits on `/auth/*` and `/api/ussd/*`.
- **Account lockout** — exponential backoff after 5 failed logins; unlock via email link or admin override.
- **Bot protection** — invisible hCaptcha on registration and password-reset forms.
- **Mobile request signing** — HMAC-signed bodies from mobile apps using a per-install key derived at first launch.
- **CSP + nonce-based inline scripts** in the Vue build to make XSS materially harder.
- **Dependency scanning** — `govulncheck` in CI, Dependabot for `go.mod`, `package.json`, Dockerfile bases.
- **Container image signing** (cosign) + SBOM (syft) — already in §7, called out here.

### Reliability
- **Health & readiness probes** — `/healthz` (liveness, always 200 if process up) and `/readyz` (200 only when DB + cache reachable). Used by Docker, the upgrade workflow, and an external uptime monitor.
- **Graceful shutdown** — `SIGTERM` → stop accepting new requests → drain in-flight (with timeout) → close DB pool. Critical for zero-downtime swaps.
- **Backup verification drill** — weekly automated restore into a throwaway DB with row-count diff (already in §6; emphasized here).
- **Disaster recovery runbook** with explicit RPO (recovery point) and RTO (recovery time) targets — and a quarterly tabletop exercise.
- **Blue/green or canary deploy** as a future evolution beyond hot-swap (e.g., two systemd units behind nginx upstream with weight shifts).
- **Migration safety** — forward-only, idempotent, `SET lock_timeout = '5s'`, never long `ALTER TABLE` without `CONCURRENTLY` variants.
- **TLS certificate monitoring** — certbot for renewal + a Prometheus alert at 14 days to expiry.

### Observability
- **Prometheus metrics** at `/metrics` — request counters by route, latency histograms, goroutine count, DB pool stats, queue depth.
- **OpenTelemetry tracing** — wire spans through middleware → handler → DB; export to a self-hosted Tempo/Jaeger.
- **Grafana dashboards** for VPS host, app, and PostgreSQL — checked in to the repo as JSON.
- **Error aggregation** — self-hosted GlitchTip (Sentry-compatible) to dedupe stack traces.
- **Structured event taxonomy** — agree on canonical event names (`user.auth.login`, `admin.maintenance.toggle`, `system.backup.created`) so dashboards and alerts stay coherent.

### Compliance & Data
- **Audit log immutability** — append-only role + hash chain (already in §3, listed here for the checklist).
- **GDPR-style data export & deletion** endpoints under `/my-portal/data`.
- **Data retention policy** documented per table: PII, audit, logs, backups.
- **Email deliverability** — SPF + DKIM + DMARC + a bounce-handling webhook so notifications don't silently land in spam.

### Localization & Accessibility
- **i18n** — at minimum English + Luganda (and Swahili if expanding regionally). Driven by JSON locale files served by the backend.
- **WCAG 2.1 AA** on both the admin and `/my-portal` — keyboard nav, color contrast, screen-reader labels. There's already a `PublicAccessibilityMenu.vue` to build on.

### Domain-Specific (Uganda context)
- **Africa's Talking SMS** integration for OTPs and notifications (most reliable UG gateway).
- **Mobile Money reconciliation** if NCS handles payments — daily settlement diff against MTN / Airtel reports.
- **USSD session store** with TTL in Redis — telecom sessions are short-lived but high-volume; in-memory maps will leak.
- **Offline-first PWA cache** for `/my-portal` core views — Uganda mobile connectivity is uneven.

### Developer Experience
- **Pre-commit hooks** — `gofmt`, `golangci-lint`, `gitleaks` (secret scan), `prettier` on the frontend.
- **Staging environment** that mirrors prod Docker Compose, with a `make refresh-staging` that anonymizes a recent prod backup and restores it.
- **API versioning** — pin mobile clients to `/api/v1/*`; future-proof for breaking changes.
- **Idempotency keys** on mutating endpoints — prevents duplicate submits from flaky mobile networks.
- **Webhook outbound system** with HMAC signing + retry-with-backoff + DLQ — for third-party integrations.
