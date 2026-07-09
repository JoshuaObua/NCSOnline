# Changelog

## 2026-07-09 - Admin and ordinary-user dashboards

### Added
- Added an administrative operations dashboard with live KPIs for all applications, athletes, organisation profiles, users, ordinary users, open forms, review queues, and application statuses.
- Added an Applications sidebar module that combines standard applications and custom-form submissions into one searchable review workspace with details, notes, approvals, rejections, and information requests.
- Added a responsive ordinary-user dashboard with Apply Now, My Applications, My Activities, Notifications, Messages, My Transactions, My Profile, and account navigation.
- Added dynamic open-form cards below the ordinary-user KPIs and a complete application flow with generated fields, draft saving, document uploads, payment-proof capture, and final submission.
- Added account registration to the login screen, role-aware routing, applicant submission listing, and common office-document support for application uploads.
- Opened the authenticated gateway paths for applicant-owned applications and transactions while retaining the intranet block on administrative domain modules.
- Added regression tests for role routing, ordinary-user navigation, registration, dashboard contracts, admin application integration, and applicant form APIs.

### Changed
- New email/password registrations now remain roleless and are treated as ordinary users; existing accounts with only the legacy `user` role receive the same applicant experience.
- Staff accounts continue to enter the administrative portal, while ordinary users are isolated from administrative routes and controls.

## 2026-07-09 - Custom form builder portal

### Added
- Added a Custom Form Builder module to the portal sidebar for authorized administrators and the General Secretary.
- Added a complete form-management view with template filtering, department ownership, lifecycle status, UGX fees, banner uploads, ordered field authoring, option and file-type configuration, live previews, archiving, submission filtering, answer review, payment verification, review notes, and approval decisions.
- Added frontend regression tests for form payload serialization, field configuration, paged submissions, sidebar registration, API integration, and the form-field database migration.

### Fixed
- Replaced the form field key constraint with an active-record unique index so templates can be edited repeatedly while historical soft-deleted fields remain available.

## 2026-07-09 - Portal-only application cutover

### Changed
- Made the frontend router portal-only: `/` now redirects to `/login`, `/login` opens the portal login, `/portal` opens the portal dashboard, `/cms` remains as a legacy redirect to `/portal`, and unknown public routes redirect to `/login`.
- Changed the login experience from CMS wording to portal wording and routed successful sign-ins and local preview sign-ins to `/portal`.
- Replaced the public website sitemap view with a portal sitemap that documents the retained routes and module slugs.
- Renamed the visible settings surface from Website Settings to Portal Settings while retaining the existing settings storage key for rollback compatibility.
- Pruned the portal sidebar registry so the only reachable modules are:
  `roles`, `manage-roles`, `users`, `manage-users`, `associations`, `manage-federations`, `create-federation-categories`, `manage-federation-categories`, `audit-logs`, `third-party-integrations`, `website-settings`, `appearance`, `sitemap`, `storage`, `command-center`, and `maintenance`.
- Limited portal permission grouping to operational resources: users, roles, audit logs, settings, storage, dashboard, federations, and federation categories.
- Reduced the initial portal data load to operational data only: federations, federation categories, roles, permissions, audit logs, and third-party integration settings.
- Updated local preview labels, sample users, and preview tokens to portal terminology while still accepting the old local CMS preview token for compatibility.
- Replaced public website route regression tests with portal-only router, login, module registry, sitemap, and audit/operations checks.

### Removed From The Reachable Portal
- Public website pages and routes for news, static pages, events, careers, projects, case studies, resources, facilities, contact pages, account pages, public federation pages, documents, investments, FAQs, team, and governing council.
- Public content CMS sections from the sidebar and portal module registry, including homepage, slideshows, blog posts, pages, projects, case studies, FAQs, resources, careers, team, council, facilities, events, investments, sports rules, press releases, reports, speeches, fun facts, newsletters, messages, notifications, comments, menus, documentation, and smart updates.
- Public/content API preloads from the portal startup path so removed modules are not loaded as part of the portal dashboard.

### Rollback Notes
- Restore the previous public route table in `frontend/src/router/index.js` to bring back public website pages.
- Re-enable the old section arrays and `contentSections` entries in `frontend/src/views/WebsiteContentManagerView.vue` to expose public/content CMS modules again.
- Restore the former `loadAll()` preload list in `frontend/src/views/WebsiteContentManagerView.vue` if public/content CMS data should load on portal startup again.
- Restore the old sitemap component in `frontend/src/components/cms/SitemapPanel.vue` if public route and content slugs need to be listed again.
- Restore the old login redirects in `frontend/src/views/CMSLoginView.vue` if `/cms` should become the primary dashboard route again.
- Revert `frontend/tests/public-pages.test.js` to the former public-page assertions when public website pages are restored.

## 2026-07-09

- Added fully editable, published Mandate and NCS History pages with responsive right-hand information tiles, structured CMS list blocks, complete default content, stable public URLs, and migration-backed seeding.
- Completed the public/CMS reliability checklist: slideshow animation modes now persist, GSAP-powered lazy page motion is installed across public routes, static pages support plain text and sanitized raw HTML, and existing builder content remains available while editing.
- Hardened backend request handling with content-aware body limits, stricter transport and cross-origin response headers, a 64 KiB header ceiling, canonical-path enforcement, protected message/notification routes, and valid JSON errors for malformed or oversized requests.
- Improved audit and analytics enrichment with GeoLite/provider field merging, trusted proxy-chain resolution, accurate city/country fallbacks, and durable browser, mobile, Postman, API-client and USSD classifications without generic device data overwriting the client channel.
- Added regression coverage for CMS settings and builder payloads, slideshow modes, request security, device/client fingerprints, geolocation fallback, public page contracts and migration seed documents; upgraded the frontend toolchain to patched Vite 6.4.3 with zero npm audit findings.
- Added complete facility-region management to the CMS with region CRUD, protected in-use deletion, category and region assignment, editable public card metadata, server-side `/facilities?region=central` filtering, seeded regional migration data, and a responsive sketch-matched public facilities layout.
- Applied the latest `Sketch.md` frontend styling pass across the public experience: rebuilt the homepage Current Leadership tile with the 80px portrait and role columns, replaced the homepage Latest News block with sketch-matched filters, featured-card and stacked-card sizing, aligned the public News page header/filter/card styling, rebuilt the CMS login screen as the centered Welcome Back card, refreshed the Docker preview, and verified desktop/mobile routes for horizontal overflow.
- Removed prefilled credentials from the CMS login form, added a detailed in-CMS documentation page covering access, publishing, modules, facilities/regions, media, API usage, maintenance and troubleshooting, published a synchronized downloadable Postman collection at `/postman/NCSMS_v1.postman_collection.json`, fixed small-screen CMS header overflow, and added regression coverage for the documentation route and blank Postman login variables.
- Cleaned up the CMS overview and operations screens by removing the Website Operations Table and Website Activity cards, normalizing System Command Center service statuses to `running` or `idle` instead of `unknown`, adding stable service-health fallbacks, and adding paginated Audit Logs with searchable page/per-page controls.

## 2026-07-08

- Rebuilt the public Careers page from `Sketch.md` with the navy hero, breadcrumb, quick-link tiles, job opportunity cards, tender cards, and procurement information panel.
- Rebuilt the public Facilities page from `Sketch2.md` with the navy hero, regional filter bar, active-region notice, and facility card grid with image badges, amenities, location, contact, and email actions.
- Added query-string region filtering to the Facilities page so URLs like `/facilities?region=central` select and filter the active region, and region tab clicks update the URL.
- Made the Facilities CMS fetch abortable so the sketch fallback remains responsive when the preview API route is unavailable.
- Rebuilt the public maintenance experience with a full-screen NCS-branded layout, live maintenance title/message/schedule details, status steps, contact actions, responsive mobile styling, and matching Go/geo-gateway 503 fallback pages.
- Added an idempotent staff-directory migration that seeds 50 unique NCS team members from `Sketch.md`, assigns every member to the appropriate institutional department, preserves CMS profile media and biographies on reruns, validates the seeded category and department relationships, and keeps the public staff page's General Secretary reference aligned with the directory.
- Rebuilt the public Governing Council page to match the staff directory design with a responsive content/sidebar layout, council profile cards, image fallbacks, expandable sanitized biographies, live contact settings, member totals, and a polished unpublished-profiles state.
