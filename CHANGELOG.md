# Changelog — Document Libraries, Homepage Cleanup & Live Chat Stub

**Date:** 2026-07-05
**Scope:** Implements the six items from `Changes to make.md`, plus follow-up FAQ, Associations, Associations-content-seeding, Associations-details-modal, and a category-creation 500-error fix.

---

## -4. Fixed 500 errors on category creation/update (follow-up)

**Root cause**: `blog_categories.slug` had a *global* `UNIQUE` constraint. Categories are shared across every content type in the CMS (Blog, Career, Event, Facility, Investment, Federation, Team Department, and the four new document types — Sports Rules, Press Releases, Reports, Speeches), so two *different* sections both wanting an intuitive slug like "general" or "sports" would collide on that constraint. `CreateBlogCategory`/`UpdateBlogCategory` didn't catch this as a duplicate-key error either, so it surfaced to admins as an opaque `500 Internal Server Error` instead of a normal validation message. This is what was happening in the logs — confirmed by reproducing it directly: creating a category named "Sports" under Press Releases failed with 500 because "sports" was already taken by an unrelated Blog category.

**Fix** (two parts, both required):
1. `backend/migrations/051_fix_category_slug_uniqueness.sql` — dropped the global `blog_categories_slug_key` constraint and replaced it with a composite unique index on `(content_type, slug)`, so the same slug can now be reused across different sections, while true duplicates *within* the same section are still rejected.
2. `backend/internal/repository/repository.go` and `backend/internal/handlers/cms.go` — `CreateBlogCategory`/`UpdateBlogCategory` now detect a genuine duplicate-key violation and return a proper `409 Conflict` with a clear message ("A category with this slug already exists for this section") instead of falling through to a generic 500.

**Verified**: reproduced the exact failure (creating "Sports" under Press Releases, colliding with the existing Blog category of the same slug) — confirmed it now succeeds with `201`; confirmed a genuine same-section duplicate now returns a clean `409` instead of `500`. No frontend changes were needed — the CMS admin already displays whatever error message the API returns.

**Related but not fixed**: `cms_posts` (which backs News, Pages, Projects, and Case Studies via a `category` discriminator column) has the identical global-`UNIQUE`-on-slug shape and could theoretically hit the same class of error if two different content types there ever want the same slug. It hasn't produced any observed 500s, and fixing it touches core public content routing (a larger, more sensitive change than the categories table), so I left it alone rather than changing it speculatively — flagging it in case it's worth addressing as a deliberate follow-up.

---

## -3. Associations details modal redesign (follow-up)

Rebuilt the "View Details" popup on the Associations page (`AssociationsView.vue`) to match the reference exactly:

- **Header**: now a `bg-gradient-to-br from-[#1a365d] to-[#2d4a7a]` navy gradient band (previously a plain white top), with a frosted `bg-white/10 backdrop-blur` icon badge, the name in large white text, and the abbreviation shown as a solid amber badge (`bg-[#f5a623] text-[#1a365d]`) directly under the name — not a separate gray pill lower down.
- **Body**, reorganized into labeled sections matching the reference: a **Leadership** section (President / General Secretary, each with a user icon, uppercase label, and value, in a two-column grid — only rendered if at least one is set), and a **Contact & Location** section (Address with a map-pin icon; Phone as a `tel:` link with a phone icon and hover state).
- **Footer actions**: an amber "Visit Official Website" button (Globe + external-link icons) and an outlined navy "Call" button (Phone icon), in a `flex flex-wrap gap-3 pt-4 border-t` row — each only rendered when that data exists.
- Widened the dialog from `max-w-lg` to `max-w-2xl` to match.
- **Bug caught and fixed along the way**: while testing this with real data, updating an association via the admin API without an `abbreviation` field in the payload silently wiped the abbreviation to empty — the handler treated it as always-overwrite (like `category`) instead of only-overwrite-if-provided (like `name`/`logo_url`/`website_url`). Fixed in `backend/internal/handlers/cms.go`, and the accidentally-wiped `VXU` abbreviation was restored.
- **VX Uganda's real leadership/contact data** (president "Dickson Niwagaba", secretary "John Bosco Abeinomugisha", address, phone, website) was included in the reference for this modal, so it was added to that one record — this is genuine reference data, not fabricated, unlike the other 51 associations seeded earlier which deliberately left these fields blank since no real data was available for them.

---

## -2. Associations directory content seeding (follow-up)

Seeded the Associations directory with the full real content set from the reference — 52 national sports associations/federations — via `backend/migrations/050_seed_associations.sql`. For each entry this populates:

- **Name** and a URL **slug** derived from it.
- **Abbreviation** (using the new field from the previous entry below) — including two genuinely tricky real-world collisions preserved faithfully from the reference: "Uganda Cricket Association" and "Uganda Cycling Association" would both naturally abbreviate to "UCA", so the second is seeded as `UCA-Cycling`; similarly "Uganda DanceSport Federation" and "Uganda Deaf Sports Federation" both collide on "UDSF", seeded as `UDSF-Dance` / `UDSF-Deaf`. A canoe/kayak federation's real abbreviation, `UC-KF`, was preserved as-is rather than "corrected" to a plain initialism.
- **Description** — the one-line summary of each federation's mandate and international affiliation (FIFA, FIBA, World Athletics, etc.).
- **Address** — the city/venue location shown in the reference (mostly "Lugogo, Kampala", with a handful of real exceptions like "Buloba, Wakiso" for the canoe/kayak federation and "Kitante, Kampala" for golf).

**Deliberately left blank**: president, secretary, phone, and website. The reference page's "President" field literally rendered the placeholder text "President" (not an actual name) for every card, so there was no real leadership data to seed — inventing plausible-sounding names for 52 real government-recognized federations would have been actively misleading on a live public site. These fields are already exposed in the CMS admin form for NCS staff to fill in with the real, current data.

Also removed a stray leftover test association ("Test" / "John Doe") that was sitting in the database from earlier manual testing, so the directory shows exactly the 52 seeded federations.

---

## -1. Associations & Federations page redesign (follow-up)

Rebuilt the public Associations directory (`AssociationsView.vue`) to match the provided reference exactly:

- **Hero**: switched from the cream/centered hero to a dark navy (`#1a365d`) hero with a Home ›  Associations breadcrumb, matching the Contact Us / Sports Rules family of pages.
- **Toolbar**: new section below the hero showing a live count ("N Registered Associations", driven by the actual loaded list, not hardcoded), a live search box (filters by name/abbreviation/category as you type), and a "Register Federation" button.
- **Cards**: restyled to `rounded-xl` with the navy/amber hover treatment (border, shadow, lift-on-hover), an icon badge (real logo when set, an Award icon otherwise), the association's name plus a new **abbreviation badge** (e.g. "FUFA"), a two-line description, and a footer row showing the president's name and address with icons, plus a "View Details →" affordance.
- **Details modal**: kept the existing click-to-open modal (there was no modal in the reference to compare against, and rebuilding the interaction from scratch would have been unnecessary), but re-themed its colors from the old orange/cream palette to the new navy/amber one so it doesn't clash with the redesigned list.
- **New field — `abbreviation`**: real federation abbreviations (FUFA, AUUS, FMU, etc.) aren't derivable from the name algorithmically (e.g. "Federation of Motorsports Clubs of Uganda" → official abbreviation "FMU", not the initialism "FMCU"), so this needed to be actual admin-entered data, not computed. Added via `backend/migrations/049_association_abbreviation.sql`, plumbed through the Go model/repository/handlers, and exposed as a field in the CMS "Federations" admin form.
- **"Register Federation" button**: there's no dedicated federation-registration workflow in this app, so it links to `/contact-us` for now — flagging this in case NCS wants a real intake form built as a follow-up (out of scope for a styling pass).

---

## 0. Homepage FAQ accordion redesign (follow-up)

- Replaced the homepage FAQ list (previously plain browser `<details>/<summary>` elements) with a custom Tailwind accordion matching the provided reference exactly: rounded-xl bordered cards, amber border/title/chevron highlight on the open item, a chevron that rotates 180° on open, and a smooth CSS grid-based expand/collapse transition. Multiple items can be open at once; the first FAQ is open by default on load.
- Fixed a pre-existing rendering bug in the process: FAQ answers are rich-text HTML from the CMS editor (e.g. `<p>...</p>`), but the old markup interpolated them as plain text, so literal `<p>` tags could show on screen. The new markup renders answers with `v-html`, matching how other rich-text CMS fields are already rendered elsewhere on the homepage.
- Removed now-dead CSS (`.faq-list`/`.fact-list` `details`/`summary`/`p` rules) left over from this and the earlier Fun Facts redesign — neither had any markup consumer left.
- **Seeded default FAQ content** via a new migration, `backend/migrations/048_seed_faqs.sql`, matching the four questions shown in the reference:
  1. "What is the mandate of NCS?" — answer text taken verbatim from the reference.
  2. "How can I form or institute a National Sports Association?"
  3. "What is the relationship between NCS and other sports bodies?"
  4. "How does a club affiliate to NCS?"
  - **Note:** only question 1's answer was visible in the reference (the other three accordion items were collapsed in the capture, so their answer text wasn't present in the markup). Answers 2–4 were authored by me, grounded in this site's existing mandate/core-functions copy (National Sports Act 2023, association registration, NCS's regulatory role) — review and edit these three if the real wording should differ.
- Removed a stray test FAQ ("Hi" / "Hi") that was already sitting in the database from earlier manual testing, so the homepage shows exactly the four seeded defaults.

## 1. New content type: Sports Rules

- Public page at **`/sports-rules`** — dynamic list of PDF rule documents, filterable by category, with a "no file attached" fallback card for unpublished/incomplete entries.
- CMS section **Sports Rules** (sidebar) with:
  - **Add New Rule** / **Manage Rules** — title, category, PDF upload (drag-and-drop), sort order, active toggle, description.
  - **Create Rule Category** / **Manage Rule Categories** — free-form, admin-creatable/manageable categories (not a fixed list).
- Backend: rows stored in a new shared `cms_documents` table with `doc_type = 'sports_rule'`.

## 2. New content type: Press Releases

- Public page at **`/press-releases`** — same dynamic list pattern as Sports Rules, but each entry supports **either**:
  - An uploaded PDF (rendered as a download card), **or**
  - An **embedded YouTube video** (paste a YouTube URL in the CMS; the public page shows a play card that expands into an inline embedded player on click — not just an outbound link).
- CMS section **Press Releases** with the same Add/Manage + category CRUD structure as Sports Rules, plus a `video_url` field.
- Categories are fully dynamic/creatable/manageable from the CMS, per the request.

## 3. New content type: NCS Reports

- Public page at **`/reports`** — dynamic PDF list with admin-manageable categories.
- CMS section **NCS Reports** (Add/Manage Reports + Create/Manage Report Categories).

## 4. New content type: NCS Speeches

- Public page at **`/speeches`** — dynamic PDF list with admin-manageable categories.
- CMS section **NCS Speeches** (Add/Manage Speeches + Create/Manage Speech Categories).

### Architecture note (2–4 share one design)

Rather than four near-identical tables/handlers, all four content types share:
- **One backend table** (`cms_documents`) with a `doc_type` discriminator column (`sports_rule` / `press_release` / `report` / `speech`), CHECK-constrained at the database level — the same pattern this codebase already uses for categories (`blog_categories.content_type`).
- **One set of Go handlers/routes** (`ListDocuments` / `CreateDocument` / `UpdateDocument` / `DeleteDocument`), parameterized by `doc_type`.
- **One shared public Vue component** (`DocumentsView.vue`), configured per-route via router `props` (title, subtitle, icon, `doc_type`) — four routes, one component.
- **Existing category infrastructure** (`blog_categories` + `content_type`) reused as-is for all four types' categories — no schema changes needed there.

This kept the change to one migration and one generic backend surface instead of four parallel ones, while the CMS admin UI still gets its own dedicated nav section per type (matching how every other content type in this CMS is organized).

## 5. Homepage: removed sections

Removed entirely from the homepage (`HomeView.vue`):
- **Facilities** showcase section ("Our Sports Facilities" grid + category filter).
- **Associations** strip ("Recognised bodies" / National Sports Associations logo row).
- **Help & Support** section ("How can we help?" — Resource Centre / Contact Support cards).

Also removed the now-unused supporting code that only existed for the Facilities section (`fallbackFacilities`, `displayFacilities`, `facilityCategories`, `filteredFacilities`, `activeFacilityCategory`, the `listFacilities()` homepage fetch, and the `home.facilities` config block) — none of it is referenced anywhere else on the page.

**Not removed:** the standalone `/facilities` and `/associations` pages and their main-navigation links — the request was to remove the homepage sections specifically, not the pages themselves.

## 6. Homepage CTA: "In Case You Need Instant Help" + Live Chat

- CTA heading changed to **"In Case You Need Instant Help"**.
- Added a **"Chat With Us Live"** button alongside the existing "Contact NCS" button.
- Clicking it opens a new **`LiveChatWidget.vue`** component — a bottom-corner chat-style panel (matching the site's existing modal conventions: `Teleport`, `Transition`, Escape-to-close). Since there is no real live-chat backend/agent system, it's an honest temporary placeholder: it explains chat is launching soon and surfaces direct phone/email/contact-form fallbacks (pulled live from the CMS contact settings, so it stays in sync if contact info changes).

---

## Files changed

**Backend (Go)**
- `backend/internal/handlers/cms.go` — (follow-up) fixed `UpdateAssociation` silently wiping `abbreviation` to empty when omitted from the request body.
- `backend/migrations/050_seed_associations.sql` — seeds the 52 default associations/federations.
- `backend/migrations/049_association_abbreviation.sql` — adds `abbreviation` to `cms_associations`.
- `backend/migrations/048_seed_faqs.sql` — seeds the four default homepage FAQ entries.
- `backend/migrations/047_document_libraries.sql` — new `cms_documents` table.
- `backend/internal/models/models.go` — `CMSDocument` struct.
- `backend/internal/repository/repository.go` — `ListDocuments`/`CreateDocument`/`GetDocumentByID`/`UpdateDocument`/`DeleteDocument`.
- `backend/internal/handlers/cms.go` — matching HTTP handlers + `validDocType` guard.
- `backend/cmd/server/main.go` — `GET /api/v1/cms/documents` (public) and `/api/v1/admin/cms/documents` (admin CRUD) routes.

**Frontend (Vue)**
- `frontend/src/api/cms.js` — `listDocuments`/`adminListDocuments`/`adminCreateDocument`/`adminUpdateDocument`/`adminDeleteDocument`, `listCategoriesByType`, and category CRUD wrappers for all four new `content_type`s.
- `frontend/src/views/public/DocumentsView.vue` — new shared public list/detail component (PDF cards, embedded-YouTube cards, "no file attached" fallback, category filter pills).
- `frontend/src/router/index.js` — `/sports-rules`, `/press-releases`, `/reports`, `/speeches` routes + page titles.
- `frontend/src/layouts/PublicLayout.vue` — added the four pages to the default nav fallback.
- `frontend/src/views/WebsiteContentManagerView.vue` — four new nav groups, form/table sections, category managers, state, and save/edit/delete handlers (mirrors the existing Resource Centre pattern); added the four pages to `STATIC_SITE_PAGES` so they're selectable in the CMS menu builder.
- `frontend/src/views/public/HomeView.vue` — removed Facilities/Associations/Help sections and their now-dead supporting code; reworked the CTA section; wired up the live chat trigger.
- `frontend/src/components/public/LiveChatWidget.vue` — new temporary live-chat placeholder component.
- `frontend/src/views/public/HomeView.vue` — (follow-up) redesigned the FAQ accordion to match the reference exactly, switched answer rendering to `v-html`, added open/close accordion state, removed dead `.faq-list`/`.fact-list` CSS.
- `frontend/src/views/public/AssociationsView.vue` — (follow-up) rebuilt hero/toolbar/card grid to match the Associations reference exactly; added live search and the abbreviation badge; re-themed the details modal to the navy/amber palette.
- `frontend/src/views/WebsiteContentManagerView.vue` — (follow-up) added the `abbreviation` field to the Federations admin form.
- `frontend/src/views/public/AssociationsView.vue` — (follow-up) rebuilt the details modal's header (navy gradient, frosted icon badge, amber abbreviation badge), added Leadership and Contact & Location sections, and footer action buttons, matching the reference exactly.

---

## Verification performed

- Backend rebuilt and compiled cleanly; migration applied to the running database.
- Frontend rebuilt cleanly (no build errors); all containers healthy post-rebuild.
- Confirmed via direct API calls: full create/update/delete cycle for documents, including the `file_url` ↔ `video_url` toggle for press releases.
- Confirmed via a real headless-browser pass:
  - Homepage no longer contains "Recognised bodies", "Help & Support", or "Our Sports Facilities"; does contain the new CTA copy and working Live Chat button (opens, closes on Escape).
  - Logged into the CMS, expanded the new "Sports Rules" nav group, created a category and a rule end-to-end through the actual admin UI (not just the API), and confirmed it appeared correctly in "Manage Rules".
  - Verified all three public card states render correctly: PDF download card, embedded-YouTube card (expands an inline player on click), and the "no file attached" fallback for an entry with neither.
  - (Follow-up) Rebuilt again after the FAQ redesign and confirmed: the first FAQ is open by default with the correct amber/chevron styling, rich-text answers render as formatted text rather than literal `<p>` tags, and clicking a second question opens it correctly.
  - (Follow-up) Rebuilt again after the Associations redesign; created a real association with an abbreviation via the API, confirmed the abbreviation badge, live count, president/address footer, and search filtering all render correctly, and that the details modal opens with the new navy/amber styling.
  - (Follow-up) Parsed all 52 association cards out of the reference programmatically (not by hand) to avoid transcription errors, verified zero missing fields and zero duplicate IDs before writing the migration, applied it, and confirmed via the live site: the count reads "52 Registered Associations," the first and last entries match the reference order, the disambiguated abbreviations (`UC-KF`, `UDSF-Dance`, `UDSF-Deaf`) render correctly, and searching ("cricket") correctly narrows to the one matching federation.
  - (Follow-up) Rebuilt again after the details-modal redesign; populated VX Uganda's real leadership/contact data via the admin API, which surfaced the `abbreviation`-wiping bug above — fixed it, restored the wiped value, rebuilt the backend, and re-verified via screenshot that the modal renders the gradient header, abbreviation badge, Leadership grid, Contact & Location section, and both footer buttons exactly as in the reference.
- All test/verification data created during this pass was deleted afterward.
