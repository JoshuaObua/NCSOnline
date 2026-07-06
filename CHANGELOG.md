# Changelog

## 2026-07-07 — Fix: home page keeps other nav dropdown items highlighted

**Bug**: on the home page (`/`), dropdown menu items other than "Home" (e.g. "About Us", "Events") could also show highlighted.

**Root cause**: yesterday's nav-highlight fix (`urlMatchesRoute` in `PublicLayout.vue`) made a dropdown's parent highlight when the current route was exactly `/` **and** one of its children pointed at `/`. That was meant to stop `startsWith('/')` from matching every route, but it left one bad case: any dropdown with a stray/unfinished child menu item pointing at `/` (which several dropdowns have, from earlier menu-builder editing) would now correctly-by-that-logic-but-wrongly-in-practice light up specifically on the home page.

**Fix**: `urlMatchesRoute()` now treats a `/` child URL as **never** matching, full stop — the top-level "Home" link already handles Home's own highlight via `isActiveLink`'s exact match, so a dropdown never legitimately needs `/` as a match target. A child pointing at `/` inside a dropdown is virtually always a stray/unconfigured entry, not an intentional "this section includes Home" design.

**Files changed**: `frontend/src/layouts/PublicLayout.vue`.

**Status**: code fix applied locally; local rebuild + VPS deploy pending (Bash was unable to execute rebuild/deploy commands this session due to a safety hold unrelated to this specific change — the user is running the rebuild/deploy commands directly).

---

## 2026-07-06 — Fix: CMS dashboard intermittently showing "Backend API not reachable" / preview data

**Reported symptom**: on the VPS, the CMS admin dashboard sometimes showed "Backend API is not reachable on port 9080" and the Notifications panel fell back to preview data.

**Investigation**: backend, database, and all containers were confirmed healthy throughout (health endpoints, JWT auth, and CORS all verified correct via isolated tests). The actual cause: the CMS dashboard's `loadAll()` fires **~53 concurrent API requests** on every page load. A subset of those hit the public `/api/v1/cms/*` route group, which shares a per-IP rate limit (`RATE_LIMIT_REQUESTS=100` per 60s) with regular public-site traffic. Reloading the dashboard a couple of times within a minute — or two admins behind the same office IP — burns through that budget quickly; under that connection burst, Chromium sometimes reports the resulting failures as a misleading "blocked by CORS policy" error rather than the real cause, which is what produces the "not reachable" / preview-data symptoms. A clean, isolated dashboard load always succeeded in testing.

**Fix**: raised `RATE_LIMIT_REQUESTS` from `100` to `400` (per 60s window) in `.env`, giving the dashboard's burst real headroom above normal public traffic. Applied to both the local dev `.env` and the VPS's production `.env`, then recreated the `backend` container on both to pick up the new value (config-only change — no rebuild needed).

**Not changed**: no application code was touched; this was purely a rate-limit tuning value. Also left the CORS/auth configuration untouched since both were confirmed already correct.

**Verified**: confirmed the live VPS backend container now reports `RATE_LIMIT_REQUESTS=400`; re-checked `/health` and a sample authenticated API call post-restart — both healthy.

---

## 2026-07-06 — Team page rebuilt as "The Staff" — department accordion directory

Completely rebuilt the public Team page (`/team`, `TeamView.vue`) to match the provided reference exactly — it was previously a flat photo-grid of individual member cards; it's now a "Staff Departments" directory grouped by institutional department, matching the reference's `AboutStaffPage` component.

- **Layout**: two-column page — a left/main "Staff Departments" card (2/3 width on desktop) and a right sidebar (1/3 width) with an "About NCS" quick-nav card and a navy "Contact the Secretariat" card, matching the reference's surrounding sidebar exactly.
- **Sidebar nav**: Mandate (→ `/pages/the-mandate`), The Council (→ `/governing-council`, the page built earlier this session), The Staff (self, shown active with the amber highlight/background per the reference) — reusing this site's existing routes rather than introducing a new `/about/*` URL structure, since only the page content (not the URL scheme) was asked to change.
- **Staff Departments grid**: one card per institutional department (from the existing 10-tier `departments` table via `listInstitutionalDepartments()`), each with a colored icon badge, name, description, and live "N Staff Members" count — all matching the reference's exact color/icon per department (Administration=blue/Building2, Human Resource=green/Users, Finance & Accounts=amber/Calculator, ICT=purple/Monitor, Sports Officers=red/Medal, Engineering & Facilities=pink/HeartPulse, Public Relations & Communications=cyan/Megaphone, Legal & Compliance=slate/Briefcase, Procurement & Records=indigo/FileText, Support Services=orange/Users).
- **Expand/collapse**: clicking a department card expands a "Department Staff (N)" panel listing that department's actual team members (photo/placeholder, name, designation) pulled from the existing `cms_team_members.department_id` relationship, or "Staff information coming soon." when a department has no assigned staff yet — matching the reference's chevron-rotate and height/opacity transition exactly.
- **Footer note**: "All departments report to the General Secretary, {name}, who oversees the day-to-day operations..." — the General Secretary's name is looked up live from the homepage's "Current Leadership Members" (added earlier this session) by matching a member whose title contains "Secretary", rather than hardcoding the reference's name, so it stays in sync if admins change it.
- **Contact card**: email/phone pulled live from the site's `contact` settings (falling back to the existing defaults) instead of being hardcoded, so it stays in sync with the rest of the site.
- **Department display order**: the reference shows departments in a fixed institutional order (Administration → Human Resource → Finance & Accounts → ... → Support Services), not alphabetical. The backend's `ListDepartmentsWithStaffCount` query returns them alphabetically, so a small fixed name→order lookup was added on the frontend to re-sort them to match — this wasn't worth a migration/schema change for a 10-item list that's essentially never reordered.

**Not changed**: no backend/API changes were needed — `listInstitutionalDepartments()` (staff counts) and `listTeam()` (department_id-linked members) already existed and provided everything this redesign needed.

**Flagged, not removed**: the live `departments` table has an 11th row, "Test Department" (`code: TEST_DEPARTMENT`, uuid id rather than the real ten's `dept_*` ids), which now renders as an 11th card on this page. This is clearly leftover test data from before this session — but since I didn't create it this session, I didn't delete it unilaterally (unlike the throwaway test data I created and cleaned up myself during today's verification). Recommend removing it via the CMS's "Manage Departments" panel.

**Files changed**: `frontend/src/views/public/TeamView.vue` (full rewrite).

**Verified**: rebuilt the frontend; confirmed via screenshot that all ten real departments render with the correct colors/icons/descriptions/order, the ICT card correctly shows "1 Staff Members" (matching the one real team member on file, department_id-linked), the sidebar highlights "The Staff" and links to the correct Mandate/Council pages, the contact card shows live contact info, and clicking a department card expands/collapses the staff panel with the chevron rotating correctly.

---

## 2026-07-06 — Fix: static pages leaking into homepage "Latest News"; nav dropdown stuck highlighted

Two bug fixes, unrelated to each other:

**1. Static pages showing up as "news"**
- **Root cause**: the homepage's "News from NCS" section (`HomeView.vue`) fetched `listPosts({status:'published', per_page:6})` with no category filter, so any `cms_posts` row — including ones saved as a static **page** (`category='page'`), a **Project**, or a **Case Study** — could appear in the "Latest News" grid and link to a broken `/news/:slug` URL. `BlogView.vue`'s own "All" filter already excludes `['page','case_study','project']` before rendering; the homepage section never had that same filter.
- **Fix**: added the same `NON_NEWS_CATEGORIES = ['page','case_study','project']` exclusion filter to the homepage's post list before rendering, matching `BlogView.vue`'s existing convention exactly.
- **Verified**: confirmed one published `category='page'` post existed in the database; before the fix it was eligible to appear in "Latest News"; after the fix, rebuilding and reloading the homepage shows only the actual `sports`-category post, matching a real news/blog category.

**2. Nav dropdown items stuck permanently highlighted**
- **Root cause**: `PublicLayout.vue`'s `isActiveTopLevel()` used `route.path.startsWith(c.url)` to decide whether a dropdown's parent label (e.g. "Events") should be highlighted based on its children's URLs. Since `startsWith('/')` is true for *every* route, any child whose URL happened to be `/` made that dropdown permanently highlighted on every single page — which is exactly what was happening: the "Events" dropdown had a stray, never-configured child menu item (added via the menu builder, left with the default label "New menu item" and default URL `/`), so "Events" showed as active on every page load regardless of which page was actually open.
- **Fix**: replaced the loose `startsWith` check with a proper `urlMatchesRoute()` helper — exact match, or prefix match only at a real path-segment boundary (`route.path === url || route.path.startsWith(url + '/')`), with `/` only matching when the route is literally the homepage. Also removed the two stray unconfigured "New menu item" placeholder nodes from the live "Events" menu entry via the CMS menu builder API, since they were dead/never-filled-in navigation entries directly causing the bug (not real content).
- **Verified**: rebuilt the frontend; confirmed "Home" now highlights correctly on `/` and "Careers" highlights correctly on `/careers`, with no dropdown stuck active on unrelated pages.

**Files changed**:
- `frontend/src/views/public/HomeView.vue` — added `NON_NEWS_CATEGORIES` filter to the homepage news fetch.
- `frontend/src/layouts/PublicLayout.vue` — replaced `isActiveTopLevel`'s `startsWith` prefix check with a segment-aware `urlMatchesRoute()` helper.
- Live CMS data: removed two unconfigured placeholder child menu items from the "Events" main-menu entry (via the admin menu API, equivalent to using the Menu Builder's own Save action).

---

## 2026-07-06 — New page: Governing Council (clone of Team)

Added a **Governing Council** directory, structured the same way as the existing "Team" (`/team`) directory, so council members can be managed and published separately from staff Team Members.

- **Data model**: rather than a new parallel table, `cms_team_members` gained a `member_group` discriminator column (`'team'` | `'council'`, default `'team'` so all existing rows/behavior are unaffected) — the same pattern already used for `cms_documents.doc_type` and `blog_categories.content_type`. Migration: `backend/migrations/052_governing_council.sql`.
- **Backend**: `ListTeam` now accepts an optional `?group=` query param (defaults to `team` for backward compatibility); `CreateTeamMember` accepts `member_group` in the body (defaults to `team`). The discriminator is immutable on update. Same routes are reused (`/api/v1/cms/team`, `/api/v1/admin/cms/team`) — no new endpoints.
- **Public page**: new route `/governing-council` → `GoverningCouncilView.vue`, a clone of `TeamView.vue` with "Governing Council" heading/copy, calling the new `listCouncil()` API helper (`GET /api/v1/cms/team?group=council`) instead of `listTeam()`. Same card layout, photo/placeholder handling, and expandable bio behavior.
- **CMS admin**: new **"Governing Council"** sidebar nav group (mirroring "Team Members") with "Add Council Member" and "Manage Council Members" — same fields as Team Members (Full Name, Designation, Sort Order, Active toggle, Dropzone photo upload, rich-text Bio), minus the Department field, since internal staff departments don't apply to council members. Reuses the `team_members:create`/`team_members:read` permissions (same underlying resource).
- `/governing-council` was also added to the CMS menu builder's static page list, alongside the existing `/team` entry.

**Files changed**:
- `backend/migrations/052_governing_council.sql` — new `member_group` column + index.
- `backend/internal/models/models.go` — `CMSTeamMember.MemberGroup` field.
- `backend/internal/repository/repository.go` — `ListTeamMembers`/`GetTeamMemberByID`/`CreateTeamMember` updated for `member_group`.
- `backend/internal/handlers/cms.go` — `validMemberGroup`, group filtering in `ListTeam`, `member_group` handling in `CreateTeamMember`.
- `frontend/src/api/cms.js` — `listCouncil`, `adminListCouncil`, `adminCreateCouncil`, `adminUpdateCouncil`, `adminDeleteCouncil`.
- `frontend/src/views/public/GoverningCouncilView.vue` — new public page (clone of `TeamView.vue`).
- `frontend/src/router/index.js` — `/governing-council` route + page title.
- `frontend/src/views/WebsiteContentManagerView.vue` — new `councilSections`, `councilGroupOpen`, `councilForm`, Add/Manage Council Member panels, `saveCouncil`/`editCouncil`/`removeCouncil`/`resetCouncilForm`, `loadAll()` fetch + assignment, permission map entries, and the `/governing-council` menu-builder entry.

**Verified**: rebuilt the backend and frontend containers, applied migration 052, confirmed the existing team member (`Samson Ogwang`) defaulted correctly to `member_group='team'`; created a test council member end-to-end through the actual CMS admin UI (not just the API) and confirmed it appeared correctly in "Manage Council Members" and on the public `/governing-council` page, while `/team` and "Manage Team Members" continued to show only the original team member — no cross-contamination between the two groups. Test data was deleted afterward.

---

## 2026-07-06 — About NCS section: dynamic Leadership Members + Dropzone upload

Changed how the homepage's "About NCS" section's "Current Leadership" card is edited and stored:

- **Before**: `homepage.leadership` was a fixed object with exactly two hardcoded slots — `chairperson_name`/`secretary_name` (plain text) and a single `chairperson_image` field that was a raw URL text input (no upload; the General Secretary had no image at all).
- **After**: `homepage.leadership` is now `{ members: [...] }` — an array of `{ name, title, image_url }` entries. The CMS admin panel ("Homepage Management → About NCS Section") now has a **"Current Leadership Members"** subpanel (mirroring the existing "About Value Cards" dynamic-list pattern) with an **"Add member"** button and a **"Remove"** button per entry, so any number of leadership members can be added or removed, not just a fixed chairperson + secretary pair.
- **Image upload**: each member's photo field is now a `DropzoneUpload` (drag-and-drop + "Choose file"), the same upload component already used for Team Members' photos — replacing the old plain-text image URL input.
- **Public homepage**: the "Current Leadership" card on the About NCS section now loops over `home.leadership.members` and renders an avatar (or a placeholder user icon if no photo is set) + title + name for each entry, instead of a hardcoded two-column Chairperson/Secretary layout.
- The two existing leadership entries (Mr. Ambrose Tashobya — Chairperson, Dr. Bernard Patrick Ogwel — General Secretary) were preserved as the seeded defaults in the new array shape, so no existing content was lost.

**Files changed**:
- `frontend/src/views/public/HomeView.vue` — `home.leadership` default changed to `{members:[...]}`; "Current Leadership" card template rewritten to loop over members dynamically.
- `frontend/src/views/WebsiteContentManagerView.vue` — `homepageDefaults.leadership` updated to match; old chairperson/secretary/image-URL fields replaced with a dynamic member-list subpanel using `DropzoneUpload`; added `addHomepageLeadershipMember()`.

**Verified**: rebuilt the frontend container; confirmed on the live public homepage that both seeded members render correctly with their titles; logged into the CMS as the seeded admin account and confirmed the "Current Leadership Members" panel shows both members as editable cards with working Dropzone upload areas, an "Add member" button, and a "Remove" button per card (no data was changed during verification).

---

## 2026-07-06 — Homepage Hero redesign (follow-up)

Restyled the homepage hero slider (`PublicSlideshow.vue`, used by `HomeView.vue`) to match the new reference exactly — design and typography only; the existing slide-transition/animation system (`fade-in`/`slide-up`/`blur-in`/etc., the `effect`/`animation_type` machinery) was left untouched, per the request.

- **Section height**: fixed breakpoints matching the reference — `500px` base, `600px` at `md` (≥768px), `700px` at `lg` (≥1024px) — replacing the old `clamp(31rem,65vw,43rem)`. Verified via live Playwright measurement against the running site: renders at exactly 500/600/700px at each breakpoint.
- **Overlay gradient**: now `linear-gradient(90deg, #1a365d 90%→70%→transparent)` (navy at 90%/70% opacity fading to transparent), matching the reference's `from-[#1a365d]/90 via-[#1a365d]/70 to-transparent`.
- **Eyebrow badge**: padding/font-size/weight adjusted to the reference's pill (`.375rem 1rem`, `.875rem`, weight 600) — still driven by the slide's `subtitle` field (falls back to "National Council of Sports"), so CMS admins keep control per-slide rather than it being hardcoded.
- **Heading**: switched from a fluid `clamp()` size to the reference's exact responsive scale — `1.875rem` → `3rem` (md) → `3.75rem` (lg), weight 700, `line-height:1.25`, `margin-bottom:1rem` (previously weight 850/line-height 1.05 with no explicit spacing).
- **Subtitle paragraph**: `1.125rem` → `1.25rem` (md), color `white/90%`, `margin-bottom:2rem` (previously a fluid clamp size with a top margin and its own max-width cap).
- **Buttons**: padding/sizing changed to `1.5rem 2rem`, `1.125rem` font, weight 600, `border-radius:.5rem` (previously `.75rem 1.2rem`/weight 850/`.55rem` radius). Primary button now uses the exact amber hover shade `#e09612` with a layered shadow that deepens on hover; the outline ("Watch Video") button is now a true 2px solid white border on transparent background that inverts (white bg, navy text) on hover, and always renders a Play icon before its label.
- **Default buttons**: a slide with no explicit `buttons` configured in the CMS now gets both a primary "Learn More" button *and* a secondary "Watch Video" (→ `/news`) button by default, matching every slide in the reference — previously the default was a single button. Admins can still override this via the existing per-slide buttons array.
- **Nav arrows**: replaced the text-character `‹`/`›` glyphs with real chevron SVG icons (matching the reference's lucide `chevron-left`/`chevron-right`), sized `3rem` (`3.5rem` at `md`), with a frosted `bg-white/20` + `backdrop-blur` + `border-white/30` treatment that turns amber on hover — previously a flat `bg-white/16` circle with no blur/border, and the icon (not just the background) turned navy on hover. Also **no longer hidden on mobile** — the reference shows them at every breakpoint.
- **Dot indicators**: resized to `.75rem` circles (`gap:.75rem`, `bottom:2rem`), with the active dot an elongated `2.5rem` amber pill — previously smaller (`.7rem`/`.5rem` gap/`1.5rem` bottom, `2.2rem` active width).
- **Mobile buttons**: no longer forced to `width:100%` — they wrap naturally with `flex-wrap`, matching the reference.
- **Bug fix found while verifying**: the hero was sitting under a `-mt-28` negative margin (in `HomeView.vue`) that canceled the page's reserved header offset, intended to let the hero bleed edge-to-edge behind the fixed site header. Since the header has an opaque white background (not transparent-over-hero), this actually hid the top ~120px of hero content — including the eyebrow badge — behind the header on narrower viewports where the text wraps to more lines. Fixed by removing the `-mt-28` override so the hero now sits below the fixed header like every other page section, matching how the reference's own `HeroSlider` reserves `mt-[104px]/mt-[116px]` for exactly this purpose. Verified via Playwright screenshots at both 1440px and 390px widths — the badge, heading, and buttons are now fully visible at both sizes.
- **Not changed**: the reference's exact `mt-[104px] md:mt-[116px]` push values weren't copied verbatim — this site already has its own header-offset convention (`.public-site > main { padding-top: 7rem }`), which is what the fix above now relies on instead of a page-specific pixel value.

**Files changed**: `frontend/src/components/public/PublicSlideshow.vue`, `frontend/src/views/public/HomeView.vue` (removed `-mt-28` from the hero section wrapper).

**Verified**: rebuilt the frontend container; confirmed via real Playwright screenshots (desktop 1440px, mobile 390px) that the badge/heading/subtitle/buttons match the reference's styling and are no longer obscured by the header; confirmed section height is exactly 500/600/700px at the three breakpoints; confirmed the nav-arrow and dot-indicator styling renders correctly (verified visually since only one CMS slide is currently configured, so the multi-slide nav chrome doesn't normally render — confirmed via a temporary visual-only DOM injection, not a data change).


