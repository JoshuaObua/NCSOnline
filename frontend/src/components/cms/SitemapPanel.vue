<template>
  <section class="sitemap-panel">
    <div class="cms-panel-head"><h2>Portal Sitemap</h2></div>
    <p class="cms-hint">Portal-only routes and module slugs currently registered in the application.</p>

    <div class="cms-panel">
      <div class="cms-panel-head"><h3>Routes <span class="sitemap-count">{{ routes.length }}</span></h3></div>
      <table class="sitemap-table">
        <thead><tr><th>Name</th><th>Path</th><th>Purpose</th></tr></thead>
        <tbody>
          <tr v-for="route in routes" :key="route.path">
            <td>{{ route.name }}</td>
            <td><a :href="route.path" target="_blank" rel="noopener">{{ route.path }}</a></td>
            <td>{{ route.purpose }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="cms-panel">
      <div class="cms-panel-head"><h3>Portal Modules <span class="sitemap-count">{{ modules.length }}</span></h3></div>
      <table class="sitemap-table">
        <thead><tr><th>Module</th><th>Slug</th><th>Scope</th></tr></thead>
        <tbody>
          <tr v-for="module in modules" :key="module.slug">
            <td>{{ module.name }}</td>
            <td><code>{{ module.slug }}</code></td>
            <td>{{ module.scope }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script setup>
const routes = [
  { name: 'Portal Login', path: '/login', purpose: 'Default entry point for all visitors.' },
  { name: 'Portal Dashboard', path: '/portal', purpose: 'Authenticated operations workspace.' },
  { name: 'User Dashboard', path: '/dashboard', purpose: 'Ordinary-user applications and account workspace.' },
  { name: 'Application Wizard', path: '/dashboard/apply/:slug', purpose: 'Section-based custom application form wizard.' },
  { name: 'Legacy Portal Redirect', path: '/cms', purpose: 'Legacy /cms URL retained for rollback and bookmarks.' },
]

const modules = [
  { name: 'Admin Dashboard', slug: 'overview', scope: 'Operations' },
  { name: 'All Applications', slug: 'applications', scope: 'Applications' },
  { name: 'Create User', slug: 'users', scope: 'Users' },
  { name: 'Manage Users', slug: 'manage-users', scope: 'Users' },
  { name: 'Add Role', slug: 'roles', scope: 'Roles' },
  { name: 'Manage Roles', slug: 'manage-roles', scope: 'Roles' },
  { name: 'Add Federation', slug: 'associations', scope: 'Federations' },
  { name: 'Manage Federations', slug: 'manage-federations', scope: 'Federations' },
  { name: 'Federation Categories', slug: 'create-federation-categories', scope: 'Federations' },
  { name: 'Manage Federation Categories', slug: 'manage-federation-categories', scope: 'Federations' },
  { name: 'Custom Form Builder', slug: 'form-builder', scope: 'Applications' },
  { name: 'Audit Logs', slug: 'audit-logs', scope: 'Security' },
  { name: 'Third-Party Integrations', slug: 'third-party-integrations', scope: 'Settings' },
  { name: 'Portal Settings', slug: 'website-settings', scope: 'Settings' },
  { name: 'Appearance', slug: 'appearance', scope: 'Settings' },
  { name: 'Sitemap', slug: 'sitemap', scope: 'Settings' },
  { name: 'Storage Settings', slug: 'storage', scope: 'Settings' },
  { name: 'Command Center', slug: 'command-center', scope: 'Operations' },
  { name: 'Maintenance & Backups', slug: 'maintenance', scope: 'Operations' },
]
</script>

<style scoped>
.cms-hint { color: var(--text-muted, #6c757d); margin-top: -.5rem; }
.sitemap-count { font-weight: 400; color: #98a6ad; font-size: .8rem; }
.sitemap-table { width: 100%; border-collapse: collapse; font-size: .85rem; }
.sitemap-table th { text-align: left; padding: .5rem .25rem; border-bottom: 2px solid #eef0f7; color: #98a6ad; font-weight: 600; }
.sitemap-table td { padding: .5rem .25rem; border-bottom: 1px solid #f4f6f9; }
.sitemap-table a { color: #6777ef; }
.sitemap-table code { border-radius: 4px; background: #f4f6f9; color: #34395e; padding: .15rem .35rem; }
</style>
