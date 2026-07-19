import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import test from 'node:test'

const read = relativePath => readFileSync(fileURLToPath(new URL(relativePath, import.meta.url)), 'utf8')
const escapeRegExp = value => value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
const sliceBetween = (source, start, end) => {
  const startIndex = source.indexOf(start)
  assert.notEqual(startIndex, -1, `missing start marker: ${start}`)
  const endIndex = source.indexOf(end, startIndex)
  assert.notEqual(endIndex, -1, `missing end marker: ${end}`)
  return source.slice(startIndex, endIndex)
}

test('router is portal-only and sends legacy/public slugs to login', () => {
  const router = read('../src/router/index.js')

  assert.match(router, /path:\s*'\/',\s*redirect:\s*'\/login'/)
  assert.match(router, /path:\s*'\/login',\s*name:\s*'PortalLogin'/)
  assert.match(router, /path:\s*'\/portal',\s*name:\s*'PortalDashboard'/)
  assert.match(router, /path:\s*'\/dashboard\/apply\/:slug',\s*name:\s*'ApplicationWizard'/)
  assert.match(router, /path:\s*'\/cms',\s*redirect:\s*'\/portal'/)
  assert.match(router, /path:\s*'\/:pathMatch\(\.\*\)\*',\s*redirect:\s*'\/login'/)

  assert.doesNotMatch(router, /PublicLayout|views\/public|AccountProfileView/)
  for (const slug of ['news', 'pages/:slug', 'events', 'careers', 'projects', 'facilities', 'contact-us']) {
    assert.doesNotMatch(router, new RegExp(`path:\\s*'${escapeRegExp(slug)}'`))
  }
})

test('portal login starts blank and supports ordinary-user registration', () => {
  const login = read('../src/views/CMSLoginView.vue')

  assert.match(login, /NCS Portal/)
  assert.match(login, /const email = ref\(''\)/)
  assert.match(login, /const password = ref\(''\)/)
  assert.match(login, /const remember = ref\(false\)/)
  assert.match(login, /registerAccount/)
  assert.match(login, /portalDestination/)
  assert.match(login, /Create Account/)
  assert.match(login, /local-portal-preview-token/)
  assert.doesNotMatch(login, /Back to Home|router\.push\('\/cms'\)/)
  assert.doesNotMatch(login, /const email = ref\('admin@ncs\.go\.ug'\)/)
  assert.doesNotMatch(login, /const password = ref\('NCS@Admin2026!'\)/)
})

test('portal shell exposes only retained operational modules', () => {
  const portal = read('../src/views/WebsiteContentManagerView.vue')

  assert.match(portal, /const active = ref\('overview'\)/)
  assert.match(portal, /id:'overview', label:'Dashboard'/)
  assert.match(portal, /id:'applications', label:'Applications'/)
  assert.match(portal, /const homepageSections = \[\]/)
  assert.match(portal, /const blogSections = \[\]/)
  assert.match(portal, /const facilitySections = \[\]/)
  assert.match(portal, /const allowedPortalSectionIds = new Set/)
  assert.match(portal, /return \['local-portal-preview-token', 'local-cms-preview-token'\]\.includes/)

  for (const id of [
    'roles',
    'manage-roles',
    'users',
    'manage-users',
    'associations',
    'manage-federations',
    'form-builder',
    'audit-logs',
    'third-party-integrations',
    'website-settings',
    'sitemap',
    'appearance',
    'storage',
    'command-center',
    'maintenance',
  ]) {
    assert.match(portal, new RegExp(`id:'${id}'`), `missing retained module ${id}`)
  }

  const contentSections = sliceBetween(portal, 'const contentSections = [', 'const profileSections')
  for (const removed of ['messages', 'investment-requests', 'notifications', 'comments', 'menus', 'documentation', 'smart-updates']) {
    assert.doesNotMatch(
      contentSections,
      new RegExp(`id:'${escapeRegExp(removed)}'`),
      `removed module ${removed} should not be in contentSections`,
    )
  }

  const loadAll = sliceBetween(portal, 'async function loadAll() {', 'async function loadAnalytics()')
  for (const removedCall of ['adminListPosts', 'adminListEvents', 'adminListFacilities', 'adminListCareers', 'adminListComments', 'listMessages', 'listNotifications']) {
    assert.doesNotMatch(loadAll, new RegExp(removedCall), `${removedCall} should not preload in portal-only mode`)
  }
})

test('sitemap documents portal routes and module slugs only', () => {
  const sitemap = read('../src/components/cms/SitemapPanel.vue')

  assert.match(sitemap, /Portal Sitemap/)
  assert.match(sitemap, /path: '\/login'/)
  assert.match(sitemap, /path: '\/portal'/)
  assert.match(sitemap, /path: '\/dashboard'/)
  assert.match(sitemap, /path: '\/dashboard\/apply\/:slug'/)
  assert.match(sitemap, /path: '\/cms'/)
  assert.match(sitemap, /slug: 'manage-users'/)
  assert.match(sitemap, /slug: 'form-builder'/)
  assert.match(sitemap, /slug: 'applications'/)
  assert.match(sitemap, /slug: 'maintenance'/)

  for (const publicPath of ['/news', '/pages/', '/events', '/careers', '/facilities']) {
    assert.doesNotMatch(sitemap, new RegExp(escapeRegExp(publicPath)))
  }
})

test('operations statuses and audit logs avoid unknown states', () => {
  const commandCenter = read('../src/components/cms/SystemCommandCenterPanel.vue')
  const portal = read('../src/views/WebsiteContentManagerView.vue')
  const operatorHandler = read('../../backend/internal/handlers/operator.go')

  assert.match(commandCenter, /displayServiceStatus/)
  assert.match(commandCenter, /return String\(service\.status \|\| 'unavailable'\)/)
  assert.match(commandCenter, /displayServiceHealth/)
  assert.match(commandCenter, /serviceStatusClass/)
  assert.match(commandCenter, /canRunServiceAction/)
  assert.match(commandCenter, /docker_available === false/)
  assert.match(operatorHandler, /"status": "unavailable"/)
  assert.match(operatorHandler, /normalizeServiceState/)
  assert.match(operatorHandler, /opsServiceCatalog/)
  assert.doesNotMatch(operatorHandler, /"status": "unknown"/)

  assert.match(portal, /auditPage = ref\(1\)/)
  assert.match(portal, /auditPerPage = ref\(20\)/)
  assert.match(portal, /auditTotalPages/)
  assert.match(portal, /changeAuditPage/)
  assert.match(portal, /Page \{\{ auditPage \}\} of \{\{ auditTotalPages \}\}/)
  assert.match(portal, /adminListAuditLogs\(\{ page:auditPage\.value, per_page:auditPerPage\.value/)
})
