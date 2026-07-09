<template>
  <section class="cms-docs" data-testid="cms-documentation-panel">
    <header class="cms-docs-hero">
      <div>
        <span class="cms-docs-kicker">CMS operating manual</span>
        <h2>How to operate the NCS website CMS</h2>
        <p>
          Use this guide to manage public content, staff records, approvals, users, settings,
          storage, backups, and API testing. It is written for editors, managers, and system
          administrators working inside the CMS.
        </p>
      </div>
      <a
        class="cms-docs-download"
        href="/postman/NCSMS_v1.postman_collection.json"
        download="NCSMS_v1.postman_collection.json"
        data-testid="cms-postman-download"
      >
        <i class="icofont-download" aria-hidden="true"></i>
        Download Postman JSON
      </a>
    </header>

    <div class="cms-docs-grid">
      <aside class="cms-docs-nav" aria-label="Documentation sections">
        <a v-for="section in sections" :key="section.id" :href="`#${section.id}`">
          {{ section.title }}
        </a>
      </aside>

      <div class="cms-docs-body">
        <section id="quick-start" class="cms-docs-block">
          <h3>Quick Start</h3>
          <ol>
            <li>Sign in at <strong>/login</strong> with your issued CMS account. The login screen intentionally does not ship with saved credentials.</li>
            <li>Open <strong>Dashboard</strong> for a health overview, recent work, tasks, content totals, and access to priority modules.</li>
            <li>Use the left sidebar to open the module you need. If a section is hidden, your role does not include that permission.</li>
            <li>Create or edit content, set the correct status, add images or downloadable files, save, then preview the public website.</li>
            <li>Use <strong>Audit Logs</strong> after sensitive changes such as role updates, deleted content, backups, storage settings, or maintenance mode.</li>
          </ol>
          <div class="cms-docs-note">
            <strong>Publishing rule:</strong> Public pages only show items that are active and published. Drafts are safe to save while work is still under review.
          </div>
        </section>

        <section id="access" class="cms-docs-block">
          <h3>Accounts, Roles, and Permissions</h3>
          <div class="cms-docs-columns">
            <article>
              <h4>For all users</h4>
              <ul>
                <li>Open <strong>Profile Settings</strong> from the top-right user menu to update your name, avatar, and account details.</li>
                <li>Use the theme toggle in the page actions area when you need light or dark mode.</li>
                <li>Use <strong>Refresh</strong> after another administrator changes content or permissions.</li>
              </ul>
            </article>
            <article>
              <h4>For administrators</h4>
              <ul>
                <li>Create users in <strong>Users -> Create New User</strong>, then assign roles in <strong>Manage Users</strong>.</li>
                <li>Create permission bundles in <strong>Roles</strong> and assign only the modules required for a job.</li>
                <li>Deactivate unused accounts instead of sharing credentials between staff.</li>
              </ul>
            </article>
          </div>
          <p>
            Recommended separation: content editors manage posts/pages/media; communications staff manage news, press releases, speeches, and reports;
            HR or administration manages team members and council members; system administrators manage storage, backups, maintenance, users, and roles.
          </p>
        </section>

        <section id="publishing" class="cms-docs-block">
          <h3>Publishing Workflow</h3>
          <ol>
            <li>Prepare the content offline or in a draft. Confirm the title, summary, category, slug, dates, and featured image.</li>
            <li>Enter the record in the correct CMS module. Avoid creating duplicate categories just to fix capitalization.</li>
            <li>Set the status to draft while reviewing. For public visibility, switch to published or active, depending on the module.</li>
            <li>Save the record and use <strong>Preview website</strong> to check desktop and mobile presentation.</li>
            <li>Return to the CMS if text wraps poorly, images crop badly, dates are wrong, or the item appears in the wrong public section.</li>
          </ol>
          <div class="cms-docs-checklist">
            <span>Before publishing, confirm:</span>
            <ul>
              <li>Slug is short, readable, and unique.</li>
              <li>Excerpt is clear and not copied from the first paragraph without review.</li>
              <li>Images have meaningful alt text and are not stretched.</li>
              <li>Categories, regions, departments, and document types are correct.</li>
              <li>Contact details, dates, and links have been tested.</li>
            </ul>
          </div>
        </section>

        <section id="module-guide" class="cms-docs-block">
          <h3>CMS Module Guide</h3>
          <div class="cms-docs-module-list">
            <article v-for="module in modules" :key="module.title">
              <div class="cms-docs-module-head">
                <i :class="module.icon" aria-hidden="true"></i>
                <h4>{{ module.title }}</h4>
              </div>
              <p>{{ module.summary }}</p>
              <ul>
                <li v-for="step in module.steps" :key="step">{{ step }}</li>
              </ul>
            </article>
          </div>
        </section>

        <section id="media" class="cms-docs-block">
          <h3>Media, Documents, and Storage</h3>
          <ul>
            <li>Use the CMS upload controls for images, PDFs, reports, and resource files so the system can store public URLs consistently.</li>
            <li>Use clear filenames before upload, such as <strong>ncs-annual-report-2026.pdf</strong> or <strong>lugogo-indoor-arena.jpg</strong>.</li>
            <li>Use <strong>Storage Settings</strong> only when changing the storage backend, bucket, CDN base URL, or Google Drive integration.</li>
            <li>If files do not appear publicly, check the storage provider, public URL prefix, bucket permissions, and audit logs.</li>
          </ul>
        </section>

        <section id="api" class="cms-docs-block">
          <h3>API and Postman Testing</h3>
          <p>
            The downloadable Postman collection includes authentication, user administration, application workflow endpoints,
            public CMS reads, admin CMS create/update/delete endpoints, newsletter export, messages, notifications,
            audit logs, health checks, and common collection variables.
          </p>
          <ol>
            <li>Download the JSON from this page and import it into Postman.</li>
            <li>Set <strong>baseUrl</strong> to the API origin. For the Docker preview use <strong>http://localhost:9080</strong>. For a direct backend use the backend host and port.</li>
            <li>Set <strong>loginEmail</strong> and <strong>loginPassword</strong> in collection variables. They are intentionally blank in the file.</li>
            <li>Run <strong>Auth -> Login</strong>. The test script saves <strong>accessToken</strong> and <strong>refreshToken</strong>.</li>
            <li>Run public CMS requests without auth, or admin CMS requests after login. Use generated IDs in collection variables for update/delete requests.</li>
          </ol>
          <div class="cms-docs-code">
            <span>Common variables</span>
            <code>baseUrl, loginEmail, loginPassword, accessToken, refreshToken, postId, categoryId, facilityId, teamMemberId</code>
          </div>
        </section>

        <section id="maintenance" class="cms-docs-block">
          <h3>Maintenance, Backups, and Recovery</h3>
          <ul>
            <li>Use <strong>Maintenance & Backups</strong> before destructive work, schema imports, or public downtime windows.</li>
            <li>Enable public website maintenance for visitor-facing downtime, and admin dashboard maintenance for staff-only downtime.</li>
            <li>Record a clear reason, expected end time, and support contact before enabling maintenance.</li>
            <li>Download backups after they complete and store them in the approved secure location.</li>
            <li>After restoring or importing schema changes, verify login, homepage, content listings, uploads, and key admin sections.</li>
          </ul>
        </section>

        <section id="troubleshooting" class="cms-docs-block">
          <h3>Troubleshooting</h3>
          <div class="cms-docs-table" role="table" aria-label="CMS troubleshooting table">
            <div role="row">
              <strong role="columnheader">Issue</strong>
              <strong role="columnheader">What to check</strong>
            </div>
            <div v-for="item in troubleshooting" :key="item.issue" role="row">
              <span role="cell">{{ item.issue }}</span>
              <span role="cell">{{ item.fix }}</span>
            </div>
          </div>
        </section>
      </div>
    </div>
  </section>
</template>

<script setup>
const sections = [
  { id: 'quick-start', title: 'Quick Start' },
  { id: 'access', title: 'Access & Roles' },
  { id: 'publishing', title: 'Publishing Workflow' },
  { id: 'module-guide', title: 'Module Guide' },
  { id: 'media', title: 'Media & Storage' },
  { id: 'api', title: 'API & Postman' },
  { id: 'maintenance', title: 'Maintenance' },
  { id: 'troubleshooting', title: 'Troubleshooting' },
]

const modules = [
  {
    title: 'Dashboard and Analytics',
    icon: 'icofont-chart-histogram',
    summary: 'Use these sections to understand activity, content volume, visitors, devices, traffic sources, and priority work.',
    steps: [
      'Review KPI cards for content, visitors, comments, and high-risk audit events.',
      'Open Analytics for top pages, top countries, device mix, acquisition channels, and recent traffic behavior.',
      'Use Dashboard task shortcuts to jump into the module that needs action.',
    ],
  },
  {
    title: 'Homepage Management',
    icon: 'icofont-home',
    summary: 'Controls the homepage sections, about block, leadership preview, topbar text, core functions, and sports excellence numbers.',
    steps: [
      'Use General Sections to toggle homepage blocks and set section order.',
      'Use About NCS Section to update intro copy, the Current Leadership tile, values, and mandate links.',
      'Use Core Functions and Sports Excellence for mandate bullets and metric cards.',
      'Preview the homepage after any save because these blocks are highly visible.',
    ],
  },
  {
    title: 'Homepage Hero Slideshow',
    icon: 'icofont-image',
    summary: 'Controls the public homepage carousel, slide text, media, buttons, timing, and animation behavior.',
    steps: [
      'Edit slideshow settings first: autoplay, effect, duration, and pause on hover.',
      'Create or edit slides with a strong image, concise title, useful description, and one or two action buttons.',
      'Keep inactive slides for future campaigns instead of deleting reusable content.',
    ],
  },
  {
    title: 'News, Posts, Categories, and Comments',
    icon: 'icofont-newspaper',
    summary: 'Manages news articles, blog posts, public comments, and taxonomy used by the news pages.',
    steps: [
      'Create categories before posts when a new editorial type is needed.',
      'Create posts with title, slug, category, cover image, excerpt, rich content, and status.',
      'Use Manage Posts to edit, unpublish, or remove old content.',
      'Moderate Blog Comments by approving, flagging, or deleting inappropriate submissions.',
    ],
  },
  {
    title: 'Static Pages and Menus',
    icon: 'icofont-page',
    summary: 'Creates editable public pages such as mandate and history, then exposes them through menus.',
    steps: [
      'Create pages with the builder blocks, right tiles, lists, rich text, and published status.',
      'Use Manage Static Pages to revise existing pages without changing their public slug unless required.',
      'Use Main Menu to add pages to header or footer navigation and hide items that should not be public yet.',
    ],
  },
  {
    title: 'Projects, Case Studies, FAQs, Resources, and Careers',
    icon: 'icofont-briefcase',
    summary: 'Content modules that each support create/manage workflows and category management.',
    steps: [
      'Create categories first for filtering and public organization.',
      'Create entries with the correct status, image or file URL, excerpt, and date information.',
      'Use the Page UI section under Careers to adjust careers landing page text and procurement guidance.',
      'Use Resource Centre for downloadable public documents that are not formal reports, press releases, speeches, or rules.',
    ],
  },
  {
    title: 'Team Members and Governing Council',
    icon: 'icofont-users-alt-5',
    summary: 'Maintains the staff directory, institutional departments, and governing council member profiles.',
    steps: [
      'Create departments before assigning staff members so filtering and grouping work correctly.',
      'For each member, add name, title, department, photo, biography, sort order, active status, and member group.',
      'Council members use the same team-member data model with the council group.',
      'Preview Team and Governing Council pages after edits to check card layout and biography expansion.',
    ],
  },
  {
    title: 'Facilities and Regions',
    icon: 'icofont-building-alt',
    summary: 'Manages public sports facilities, facility categories, and region filters such as /facilities?region=central.',
    steps: [
      'Create regions with stable slugs such as central, northern, eastern, western, and southern.',
      'Assign each facility to a region, category, status, location, amenities, phone, email, and image.',
      'Use Manage Facility Regions carefully: regions in use are protected from unsafe deletion.',
      'Preview the public facilities page with each region query string after major edits.',
    ],
  },
  {
    title: 'Events, Investments, Federations, Documents, and Fun Facts',
    icon: 'icofont-trophy',
    summary: 'Specialized public modules for calendars, investment opportunities, sports bodies, formal documents, and rotating facts.',
    steps: [
      'Use Events for dated calendar items with venue, category, status, and public detail pages.',
      'Use Invest With Us for investment posts and manage inbound investment requests from the public form.',
      'Use Federations for association profiles, logos, contact officers, and website links.',
      'Use Sports Rules, Press Releases, Reports, and Speeches for document records and downloadable files.',
      'Use Fun Facts for short rotating homepage facts.',
    ],
  },
  {
    title: 'Messages, Newsletter, Notifications, and Audit Logs',
    icon: 'icofont-envelope',
    summary: 'Operational queues for public contact, subscribers, alerts, and compliance review.',
    steps: [
      'Use Contact Messages to mark inquiries read, resolved, archived, or deleted.',
      'Use Newsletter to review subscribers and export CSV or Excel lists.',
      'Use Notifications to mark read, clear, or inspect system notices.',
      'Use Audit Logs after sensitive actions to confirm who changed what, when, and from where.',
    ],
  },
  {
    title: 'Settings, Appearance, Storage, Command Center, and Smart Updates',
    icon: 'icofont-tools',
    summary: 'System-level controls for site identity, contact data, integrations, typography, storage providers, commands, and update workflows.',
    steps: [
      'Use Contact Details and Website Settings for public identity, logos, footer data, and social/contact links.',
      'Use Third-Party Integrations for analytics and external scripts only after verifying IDs.',
      'Use Appearance for typography and custom fonts; check public pages after changes.',
      'Use Storage Settings only when you understand the target provider and public URL behavior.',
      'Use Command Center and Smart Updates for controlled maintenance actions and operational tasks.',
    ],
  },
]

const troubleshooting = [
  { issue: 'Section is missing from sidebar', fix: 'Your role probably lacks the required permission. Ask an administrator to check Roles and Manage Users.' },
  { issue: 'Saved content does not show publicly', fix: 'Confirm status is published or active, category/type is correct, dates are valid, and the public page is not cached.' },
  { issue: 'Image or document does not load', fix: 'Check the upload URL, storage public access settings, and whether the file was uploaded to the configured provider.' },
  { issue: 'Facilities filter shows no records', fix: 'Confirm the facility has a region slug matching the public query, for example central for /facilities?region=central.' },
  { issue: 'Postman admin request returns 401', fix: 'Run Auth -> Login first and confirm accessToken was saved in collection variables.' },
  { issue: 'Postman admin request returns 403', fix: 'The authenticated user lacks the role or permission required for that endpoint.' },
  { issue: 'Build or preview looks stale', fix: 'Refresh the CMS data, rebuild the frontend preview if needed, and clear browser cache for static assets.' },
]
</script>

<style scoped>
.cms-docs {
  display: grid;
  gap: 1.25rem;
}
.cms-docs-hero {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  border-radius: 3px;
  background: #fff;
  padding: 1.5rem;
  box-shadow: 0 4px 25px rgba(0, 0, 0, .1);
}
.cms-docs-kicker {
  display: inline-block;
  margin-bottom: .55rem;
  color: #ffa426;
  font-size: .72rem;
  font-weight: 800;
  letter-spacing: .12em;
  text-transform: uppercase;
}
.cms-docs h2,
.cms-docs h3,
.cms-docs h4 {
  color: #34395e;
  font-weight: 800;
}
.cms-docs h2 {
  margin: 0 0 .55rem;
  font-size: 1.45rem;
}
.cms-docs h3 {
  margin: 0 0 .75rem;
  font-size: 1.05rem;
}
.cms-docs h4 {
  margin: 0;
  font-size: .95rem;
}
.cms-docs p,
.cms-docs li,
.cms-docs span {
  color: #5f6677;
  line-height: 1.65;
}
.cms-docs-download {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  gap: .5rem;
  border-radius: 30px;
  background: #6777ef;
  color: #fff;
  padding: .7rem 1rem;
  font-size: .78rem;
  font-weight: 700;
  box-shadow: 0 2px 6px #acb5f6;
}
.cms-docs-grid {
  display: grid;
  grid-template-columns: 15rem minmax(0, 1fr);
  gap: 1rem;
  align-items: start;
}
.cms-docs-nav {
  position: sticky;
  top: 92px;
  display: grid;
  gap: .25rem;
  border-radius: 3px;
  background: #fff;
  padding: .85rem;
  box-shadow: 0 4px 25px rgba(0, 0, 0, .1);
}
.cms-docs-nav a {
  border-radius: 4px;
  color: #6777ef;
  font-size: .8rem;
  font-weight: 700;
  padding: .48rem .65rem;
}
.cms-docs-nav a:hover {
  background: #f4f6f9;
}
.cms-docs-body {
  display: grid;
  gap: 1rem;
  min-width: 0;
}
.cms-docs-block {
  border-radius: 3px;
  background: #fff;
  padding: 1.25rem;
  box-shadow: 0 4px 25px rgba(0, 0, 0, .1);
}
.cms-docs-block ol,
.cms-docs-block ul {
  margin: 0;
  padding-left: 1.2rem;
}
.cms-docs-note,
.cms-docs-checklist,
.cms-docs-code {
  margin-top: 1rem;
  border-left: 4px solid #ffa426;
  background: #fff8ec;
  padding: .85rem 1rem;
}
.cms-docs-columns,
.cms-docs-module-list {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: .85rem;
}
.cms-docs-columns article,
.cms-docs-module-list article {
  border: 1px solid #f0f1f7;
  border-radius: 3px;
  padding: 1rem;
}
.cms-docs-module-head {
  display: flex;
  align-items: center;
  gap: .6rem;
  margin-bottom: .6rem;
}
.cms-docs-module-head i {
  color: #6777ef;
  font-size: 1.2rem;
}
.cms-docs-code code {
  display: block;
  margin-top: .4rem;
  overflow-wrap: anywhere;
  color: #34395e;
  font-size: .82rem;
}
.cms-docs-table {
  display: grid;
  border: 1px solid #f0f1f7;
  border-radius: 3px;
  overflow: hidden;
}
.cms-docs-table > div {
  display: grid;
  grid-template-columns: minmax(10rem, .38fr) minmax(0, 1fr);
  gap: 1rem;
  padding: .8rem 1rem;
  border-bottom: 1px solid #f0f1f7;
}
.cms-docs-table > div:first-child {
  background: #f4f6f9;
}
.cms-docs-table > div:last-child {
  border-bottom: 0;
}
:global(.dark) .cms-docs-hero,
:global(.dark) .cms-docs-nav,
:global(.dark) .cms-docs-block,
:global(.dark) .cms-docs-columns article,
:global(.dark) .cms-docs-module-list article,
:global(.dark) .cms-docs-table {
  background: #1f2937;
  border-color: #334155;
}
:global(.dark) .cms-docs h2,
:global(.dark) .cms-docs h3,
:global(.dark) .cms-docs h4,
:global(.dark) .cms-docs-code code,
:global(.dark) .cms-docs-table strong {
  color: #f8fafc;
}
:global(.dark) .cms-docs p,
:global(.dark) .cms-docs li,
:global(.dark) .cms-docs span {
  color: #cbd5e1;
}
:global(.dark) .cms-docs-note,
:global(.dark) .cms-docs-checklist,
:global(.dark) .cms-docs-code {
  background: #422006;
}
:global(.dark) .cms-docs-table > div:first-child {
  background: #111827;
}
@media (max-width: 1100px) {
  .cms-docs-grid,
  .cms-docs-columns,
  .cms-docs-module-list {
    grid-template-columns: 1fr;
  }
  .cms-docs-nav {
    position: static;
  }
}
@media (max-width: 640px) {
  .cms-docs-hero {
    flex-direction: column;
  }
  .cms-docs-download {
    width: 100%;
  }
  .cms-docs-table > div {
    grid-template-columns: 1fr;
    gap: .35rem;
  }
}
</style>
