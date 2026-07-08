## 2026-07-08

- Rebuilt the public Careers page from `Sketch.md` with the navy hero, breadcrumb, quick-link tiles, job opportunity cards, tender cards, and procurement information panel.
- Rebuilt the public Facilities page from `Sketch2.md` with the navy hero, regional filter bar, active-region notice, and facility card grid with image badges, amenities, location, contact, and email actions.
- Added query-string region filtering to the Facilities page so URLs like `/facilities?region=central` select and filter the active region, and region tab clicks update the URL.
- Made the Facilities CMS fetch abortable so the sketch fallback remains responsive when the preview API route is unavailable.
- Rebuilt the public maintenance experience with a full-screen NCS-branded layout, live maintenance title/message/schedule details, status steps, contact actions, responsive mobile styling, and matching Go/geo-gateway 503 fallback pages.
- Added an idempotent staff-directory migration that seeds 50 unique NCS team members from `Sketch.md`, assigns every member to the appropriate institutional department, preserves CMS profile media and biographies on reruns, validates the seeded category and department relationships, and keeps the public staff page's General Secretary reference aligned with the directory.
- Rebuilt the public Governing Council page to match the staff directory design with a responsive content/sidebar layout, council profile cards, image fallbacks, expandable sanitized biographies, live contact settings, member totals, and a polished unpublished-profiles state.
