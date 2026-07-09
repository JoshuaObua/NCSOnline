## 2026-07-09

- Added fully editable, published Mandate and NCS History pages with responsive right-hand information tiles, structured CMS list blocks, complete default content, stable public URLs, and migration-backed seeding.
- Completed the public/CMS reliability checklist: slideshow animation modes now persist, GSAP-powered lazy page motion is installed across public routes, static pages support plain text and sanitized raw HTML, and existing builder content remains available while editing.
- Hardened backend request handling with content-aware body limits, stricter transport and cross-origin response headers, a 64 KiB header ceiling, canonical-path enforcement, protected message/notification routes, and valid JSON errors for malformed or oversized requests.
- Improved audit and analytics enrichment with GeoLite/provider field merging, trusted proxy-chain resolution, accurate city/country fallbacks, and durable browser, mobile, Postman, API-client and USSD classifications without generic device data overwriting the client channel.
- Added regression coverage for CMS settings and builder payloads, slideshow modes, request security, device/client fingerprints, geolocation fallback, public page contracts and migration seed documents; upgraded the frontend toolchain to patched Vite 6.4.3 with zero npm audit findings.
- Added complete facility-region management to the CMS with region CRUD, protected in-use deletion, category and region assignment, editable public card metadata, server-side `/facilities?region=central` filtering, seeded regional migration data, and a responsive sketch-matched public facilities layout.

## 2026-07-08

- Rebuilt the public Careers page from `Sketch.md` with the navy hero, breadcrumb, quick-link tiles, job opportunity cards, tender cards, and procurement information panel.
- Rebuilt the public Facilities page from `Sketch2.md` with the navy hero, regional filter bar, active-region notice, and facility card grid with image badges, amenities, location, contact, and email actions.
- Added query-string region filtering to the Facilities page so URLs like `/facilities?region=central` select and filter the active region, and region tab clicks update the URL.
- Made the Facilities CMS fetch abortable so the sketch fallback remains responsive when the preview API route is unavailable.
- Rebuilt the public maintenance experience with a full-screen NCS-branded layout, live maintenance title/message/schedule details, status steps, contact actions, responsive mobile styling, and matching Go/geo-gateway 503 fallback pages.
- Added an idempotent staff-directory migration that seeds 50 unique NCS team members from `Sketch.md`, assigns every member to the appropriate institutional department, preserves CMS profile media and biographies on reruns, validates the seeded category and department relationships, and keeps the public staff page's General Secretary reference aligned with the directory.
- Rebuilt the public Governing Council page to match the staff directory design with a responsive content/sidebar layout, council profile cards, image fallbacks, expandable sanitized biographies, live contact settings, member totals, and a polished unpublished-profiles state.
