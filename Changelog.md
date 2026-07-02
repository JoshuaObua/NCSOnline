# CMS Changelog

## 2026-07-01

### Added
- Added grouped CMS sidebar modules for Facilities, Events, Invest With Us, Federations, Fun Facts, and Newsletter.
- Added Otika-styled create/manage/category panels for facilities, events, investment posts/categories, federations/categories, and fun-fact articles.
- Added Newsletter Subscribers CMS UI with submitted email, date, day, status, refresh, CSV export, and Excel export controls.
- Added frontend CMS API wrappers for facility, event, investment, federation category operations and newsletter subscriber list/export operations.
- Added Manage Projects under the Projects sidebar group with an Otika-styled projects table.
- Added grouped FAQ, Resource Centre, Careers, Team Members, and Roles sidebar menus with create/manage/category or department submenus.
- Added Otika-styled create/manage panels for FAQ articles, FAQ categories, Resource Centre articles, resource article categories, career posts, career categories, team members, team departments, role creation, and role management.
- Added frontend CMS API wrappers for FAQ, Resource, Career, and Team Department scoped category/department operations.
- Added grouped CMS sidebar menus for Static Pages, Projects, and Case Studies, each with dedicated create/manage/category submenu screens.
- Added Otika-themed Project Category and Case Study Category create/manage panels.
- Added frontend CMS API wrappers for project and case-study category list/create/update/delete operations using scoped category requests.
- Rebuilt the CMS desktop overview dashboard to match Otika's `index.html` dashboard structure with statistic cards, page visit analytics, website KPI summaries, platform mix, traffic sources, content analytics, operations table, activity feed, analytics snapshot, and the Otika-style main footer.
- Added a shared Otika asset loader so CMS screens load template CSS and Font Awesome fonts directly from `/otika-assets`.
- Added Otika Bootstrap Admin Template assets to the frontend and applied them to CMS-only screens.
- Added Otika-style CMS topbar with global search, messages, notifications, profile menu, refresh, and preview actions.
- Added Otika-style sidebar structure with menu headers, dropdown module groups, profile block, active states, and compact navigation.
- Added Role Management to the CMS sidebar for flexible RBAC/BRAC administration.
- Added custom role creation, editing, deletion, and per-permission customization for content-manager access control.
- Added an Audit Logs sidebar option for tracing request traffic, users, endpoints, IPs, response codes, response times, threat signals, and immutable audit-chain details.
- Added Static Page creation and management to the CMS sidebar.
- Added Project post management as a dedicated CMS module.
- Added Case Study post management as a dedicated CMS module.
- Added Career section post management for jobs, tenders, internships, and other career notices.
- Added Resource Centre upload controls with direct PDF, Word, and Excel document uploads.
- Added Invest With Us content management.
- Added Team Members content management.
- Added overview counters for pages, careers, and resources.

### Improved
- Renamed the CMS sidebar `Menus` item to `Main Menu` and `Contact/Footer` to `Footer Content`, including matching page headings.
- Cleaned the CMS sidebar so Blogs Management now contains only blog posts and blog category actions, while static pages, project posts, and case study posts are independent Content Type items with their own icons.
- Refined collapsed CMS sidebar behavior so it stays icon-only and reveals item titles as hover/focus tooltips.
- Replaced heavy browser focus borders with cleaner accessible focus rings across CMS buttons, links, inputs, and sidebar items.
- Expanded CMS dark mode coverage across Otika cards, tables, forms, dropdowns, sidebar, navbar icons, editor controls, alerts, and text states.
- Made the CMS sidebar scrollable with fixed logo/user areas, contained menu scrolling, ellipsized labels, and icon-only collapsed state.
- Replaced sidebar email fallback with a shorter username/local-name display to prevent overflow in the user block.
- Updated the CMS topbar icon buttons so the sidebar toggle is black and navbar icons have consistent black styling.
- Made CMS topbar messages, notifications, global search, audit log preview/detail views, and role preview interactions functional in offline preview mode.
- Changed the CMS sidebar brand to logo-only by removing the `NCS CMS` text label.
- Improved CMS responsiveness for tablet and mobile layouts, including stacked sidebars, wrapped topbar controls, single-column forms, adaptive tables, and mobile-friendly dropdowns.
- Added an API availability gate for CMS loading so an offline native Go backend on port `9080` does not trigger a cascade of failed CMS endpoint requests.
- Added a login health check so the default local CMS preview can open without posting to an offline auth endpoint.
- Preloaded Otika Font Awesome font files and removed Vue-scoped CSS imports that rewrote font URLs through `/public`.
- Restyled the CMS login screen to match Otika's `card card-primary` authentication layout.
- Restyled CMS cards, upload forms, input fields, textareas, checkboxes, selects, buttons, tables, permission matrix, audit rows, and alerts to match Otika dashboard styling.
- Added CMS API wrappers for roles, permissions, permission assignment/removal, and audit log lookup.
- Added admin overview counters for roles and audit events.
- Hardened CMS list loading so failed/offline API responses keep role, permission, audit, and content lists as arrays instead of crashing Vue rendering.
- Reused the existing post backend for static pages, projects, and case studies through fixed CMS categories.
- Reused existing backend-backed modules for careers, resources, investment items, and team members.
- Added image/document upload widgets to generic CMS forms for file, image, cover image, and logo fields.
- Made the post editor category-aware so page, project, case study, news, and blog drafts stay separate.

### Verified
- Ran `npm run build` in `frontend`; the production build completed successfully.

### Events Otika Form Update

- Reworked Add New Event into an Otika basic-form style card with `card`, `card-header`, `card-body`, `form-group`, `form-control`, `selectric`, and `card-footer` classes.
- Added event cover image upload through the existing Dropzone media uploader.
- Added Summernote-style rich text editing for event descriptions.
- Reworked Add Event Category into an Otika card form with matching inputs, textarea, checkbox, and footer action.
- Added event end-date and cover-image fields to the CMS event form state so edit/save preserves them.

### Verified
- Ran `npm run build` in `frontend`; the production build completed successfully.

### Team, Careers, Departments, And FAQ Otika Forms

- Reworked Team Member creation/editing into an Otika basic-form style card with `form-group`, `form-control`, `selectric`, custom checkbox, Dropzone image upload, and a Summernote-style bio editor.
- Reworked Career Post creation/editing into an Otika card form with automatic department/category dropdowns and Summernote-style editors for description and requirements.
- Reworked Create Department, Create Career Category, Create FAQ Article, and Create FAQ Category screens to use Otika card, input, textarea, checkbox, and footer button classes.
- Added a reusable CMS rich text editor component styled like the Otika/Summernote editor.
- Added admin CMS department backend endpoints so Create/Manage Departments uses the real institutional departments table used by team members and career posts.
- Updated team/career department dropdowns to load from the real department records and omit blank department IDs from save payloads.

### Verified
- Ran `npm run build` in `frontend`; the production build completed successfully.
- Ran `go test ./internal/handlers ./internal/repository -run TestDoesNotExist -count=0`; the touched backend packages compiled successfully.

### CMS RBAC Matrix And User Management

- Added a granular CMS permission migration covering every sidebar module with separate read, create, update, delete, export, reset password, and role-assignment permissions where applicable.
- Added Users to the CMS sidebar with Create New User and Manage Users screens.
- Wired CMS user management to the existing Go backend user APIs for create, edit, delete, activate/deactivate, password reset, and role assign/revoke.
- Updated the CMS sidebar to hide whole groups and individual submenu items when the current user does not have permission for those sections.
- Updated the role permission matrix to show human-readable module labels and granular permissions for blog posts, blog categories, static pages, projects, case studies, FAQs, resources, careers, team, facilities, events, investments, federations, fun facts, newsletter, comments, menus, settings, storage, users, roles, and audit logs.
- Added preview fallback permissions and preview users so the UI remains testable when the native Go API is offline.

### Verified
- Ran `npm run build` in `frontend`; the production build completed successfully.
- Ran `go test ./internal/handlers ./internal/services ./internal/repository -run TestDoesNotExist -count=0`; the touched backend packages compiled successfully.

### CMS Otika Completion Sweep

- Replaced the unreliable Facilities sidebar and fallback public icon with the stable Otika/Icofont building icon.
- Converted the remaining CMS delete confirmations to SweetAlert2 confirmations.
- Converted the menu builder delete confirmation to SweetAlert2 and refreshed its controls with Otika-style buttons, inputs, shadows, and drag handles.
- Converted the public blog comment success popup from a native browser alert to SweetAlert2.
- Restyled the storage settings child component to match the Otika card, form-control, badge, and rounded button treatment.
- Re-scanned the frontend source to confirm no native `alert()`, `confirm()`, or `window.confirm()` calls remain.

### Verified
- Ran `npm run build` in `frontend`; the production build completed successfully.

### Contact Details And Public Content

- Removed the avatar and username block from the CMS sidebar so navigation starts at Overview.
- Replaced the old footer settings screen with an Otika-styled Contact Details screen for phone, WhatsApp, email, fax, location, P.O. Box, physical address, Google Maps URL, social links, and footer About NCS description.
- Kept Contact Details backed by the existing CMS settings backend keys: `contact` and `footer`.
- Removed My Portal buttons and links from the public header, mobile menu, footer, homepage help area, homepage CTA, and Contact Us quick links.
- Updated the public homepage Fun Facts section to use CMS Fun Facts records from the database instead of homepage Sports Excellence counters.
- Hardened public homepage FAQ and Fun Facts loading to accept either array responses or paginated `{ items }` responses from the CMS API.
- Expanded the public Contact Us page to render CMS-managed WhatsApp, fax, location, P.O. Box, Google Maps, and social media links.

### Verified
- Ran `npm run build` in `frontend`; the production build completed successfully.

### Static Page Builder

- Replaced static page creation with an Otika-styled page builder instead of the blog/content editor.
- Removed the static page category selector from the CMS UI; static pages are saved internally as `category: page`.
- Added nestable page blocks for text sections, images, rows/columns, accordions, dropdowns, icon cards, and raw HTML.
- Added breadcrumb background image upload support using the existing static page cover image field.
- Normalized static page slugs directly from the page title so slugs stay clean, such as `about-ncs`, without a `page-` prefix.
- Updated the public static page renderer to display builder blocks while preserving old HTML static page content.
- Added backend compatibility fields for `page_builder` JSON and `breadcrumb_image_url` so static page block content can be sent explicitly by future clients.

### Verified
- Ran `npm run build` in `frontend`; the production build completed successfully.
- Ran `go test ./internal/handlers -run TestDoesNotExist -count=0`; the touched backend handler package compiled successfully.

### Homepage Management Sections

- Added `About NCS Section`, `Topbar Section`, `Core Functions`, and `Sports Excellence` screens under Homepage Management in the CMS sidebar.
- Added Otika-styled CMS controls for About NCS text, leadership labels, value cards, core function add/edit/delete rows, topbar marquee messages, social links, webmail URL, and Sports Excellence counters.
- Added homepage settings defaults and merge safeguards so older saved homepage records do not break newly added nested CMS fields.
- Linked the public topbar marquee, social links, and webmail button to the new Homepage Management topbar settings.
- Linked the public About NCS, value cards, core functions, and Sports Excellence counters to CMS-managed homepage settings.
- Added animated count-up rendering and counter descriptions for the four Sports Excellence public homepage cards.

### Verified
- Ran `npm run build` in `frontend`; the production build completed successfully.
## 2026-07-02

### CMS Tables And Footer

- Reworked shared CMS management lists to render as Otika basic table cards with responsive `table table-striped table-hover table-sm` styling.
- Applied the Otika basic table treatment to direct CMS tables such as dashboard tasks and newsletter subscribers.
- Fixed the CMS footer at the bottom of the content area and updated the footer attribution to `Design By: Ateni Media Technologies LLC`.

### Verified
- Ran `npm run build` in `frontend`; the production build completed successfully.

### CMS Slideshow/Navbar Follow-up

- Fixed the CMS navbar so it stays fixed above scrolling content, uses a white background in light mode, and keeps the dark background in dark mode.
- Added high-contrast native select option colors for slideshow dropdowns in both light and dark modes.
- Replaced slideshow delete confirmation with Otika SweetAlert and added SweetAlert success/error feedback for slideshow saves.
- Added slideshow transition modes for fade, horizontal slide, vertical slide, and zoom.
- Added slide animation modes for fade, slide up, slide left, slide right, zoom in, zoom out, flip in, blur in, bounce in, and Ken Burns.
- Reworked CTA button editing into roomier compact cards and saved both structured buttons and legacy button fields so the public homepage can render them.
- Rebuilt the public homepage slideshow component so CMS-managed slides and animation settings override the fallback hero while still keeping a default slide when no CMS slides exist.

### Verified
- Ran `npm run build` in `frontend`; the production build completed successfully.

### Slideshow Otika UI Update

- Reworked the slideshow manager to use Otika card, striped/hover table, compact form-group, form-control, badge, button, card-footer, and modal patterns.
- Replaced the slideshow/media uploader with Otika's multiple-upload Dropzone layout while keeping the existing CMS media upload validation and API integration.
- Moved slideshow-level settings into a popup modal so the slide editor stays less crowded.
- Wired slideshow create, update, delete, settings save, and drag-sort actions through the existing frontend API helpers for the Go backend CMS endpoints.
- Added dark-mode styling coverage for the slideshow cards, modal, table active state, controls, and Dropzone uploader.

### Verified
- Ran `npm run build` in `frontend`; the production build completed successfully.
