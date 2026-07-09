import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import test from 'node:test'

import { blockMap, createNode, pageBuilderVersion } from '../src/utils/pageBuilderRegistry.js'

const read = relativePath => readFileSync(fileURLToPath(new URL(relativePath, import.meta.url)), 'utf8')

test('editable list blocks use the current page-builder schema', () => {
  assert.equal(pageBuilderVersion, 2)
  assert.equal(blockMap.list.label, 'Editable List')

  const node = createNode('list')
  assert.equal(node.type, 'list')
  assert.equal(node.props.ordered, false)
  assert.match(node.props.items, /First item/)
})

test('Mandate and NCS History migrations contain valid editable page documents', () => {
  const sql = read('../../backend/migrations/055_seed_mandate_history_pages.sql')
  for (const [name, expectedText] of [
    ['mandate', 'The Council Shall'],
    ['history', 'Our Journey'],
  ]) {
    const match = sql.match(new RegExp(`\\$${name}\\$([\\s\\S]*?)\\$${name}\\$`))
    assert.ok(match, `missing ${name} builder payload`)
    const page = JSON.parse(match[1])
    assert.equal(page.type, 'ncs-page-builder')
    assert.equal(page.version, 2)
    assert.ok(page.blocks.length >= 3)
    assert.match(JSON.stringify(page), new RegExp(expectedText))
  }
})

test('public routes and page motion remain wired into the application shell', () => {
  const router = read('../src/router/index.js')
  const layout = read('../src/layouts/PublicLayout.vue')
  const pageView = read('../src/views/public/PageView.vue')

  assert.match(router, /path:\s*'pages\/:slug'/)
  assert.match(layout, /animatePublicPage/)
  assert.match(pageView, /block\.type === 'list'/)
  assert.match(pageView, /sanitizeRichHtml/)
})

test('facilities use managed regions and server-side query filtering', () => {
  const api = read('../src/api/cms.js')
  const facilitiesView = read('../src/views/public/FacilitiesView.vue')
  const cmsView = read('../src/views/WebsiteContentManagerView.vue')
  const migration = read('../../backend/migrations/056_facility_regions.sql')

  assert.match(api, /listFacilityRegions/)
  assert.match(api, /adminCreateFacilityRegion/)
  assert.match(facilitiesView, /listFacilities\(\{ region \}\)/)
  assert.match(facilitiesView, /region\.slug/)
  assert.match(cmsView, /manage-facility-regions/)
  assert.match(cmsView, /availability_status/)
  assert.match(migration, /CREATE TABLE IF NOT EXISTS cms_facility_regions/)
  assert.match(migration, /FOREIGN KEY \(region_slug\)/)
})

test('CMS login starts blank and the documentation page is registered', () => {
  const login = read('../src/views/CMSLoginView.vue')
  const cmsView = read('../src/views/WebsiteContentManagerView.vue')
  const docs = read('../src/components/cms/CmsDocumentationPanel.vue')

  assert.match(login, /const email = ref\(''\)/)
  assert.match(login, /const password = ref\(''\)/)
  assert.match(login, /const remember = ref\(false\)/)
  assert.doesNotMatch(login, /const email = ref\('admin@ncs\.go\.ug'\)/)
  assert.doesNotMatch(login, /const password = ref\('NCS@Admin2026!'\)/)

  assert.match(cmsView, /CmsDocumentationPanel/)
  assert.match(cmsView, /id:'documentation'/)
  assert.match(docs, /data-testid="cms-documentation-panel"/)
  assert.match(docs, /Download Postman JSON/)
  assert.match(docs, /Facilities and Regions/)
  assert.match(docs, /\/facilities\?region=central/)
})

test('CMS Postman collection is downloadable and uses blank login variables', () => {
  const backendCollection = read('../../backend/postman/NCSMS_v1.postman_collection.json')
  const publicCollection = read('../public/postman/NCSMS_v1.postman_collection.json')
  assert.equal(publicCollection, backendCollection)

  const collection = JSON.parse(backendCollection)
  const folders = collection.item.map(item => item.name)
  for (const name of [
    'Auth',
    'CMS / Public Content',
    'CMS / Admin Content',
    'CMS / Admin Taxonomy and Directories',
    'CMS / Operations Queues',
    'CMS / Settings, Media, and System',
  ]) {
    assert.ok(folders.includes(name), `missing ${name} folder`)
  }

  const variables = new Map(collection.variable.map(variable => [variable.key, variable.value]))
  assert.equal(variables.get('baseUrl'), 'http://localhost:9080')
  assert.equal(variables.get('loginEmail'), '')
  assert.equal(variables.get('loginPassword'), '')
  assert.ok(variables.has('accessToken'))
  assert.ok(variables.has('refreshToken'))

  const authFolder = collection.item.find(item => item.name === 'Auth')
  const loginRequest = authFolder.item.find(item => item.name === 'Login')
  assert.match(loginRequest.request.body.raw, /\{\{loginEmail\}\}/)
  assert.match(loginRequest.request.body.raw, /\{\{loginPassword\}\}/)
  assert.doesNotMatch(loginRequest.request.body.raw, /admin@ncs\.go\.ug|NCS@Admin2026!|changeme123/)
  assert.match(JSON.stringify(collection), /\/api\/v1\/cms\/facility-regions/)
  assert.match(JSON.stringify(collection), /\/api\/v1\/cms\/facilities/)
  assert.match(JSON.stringify(collection), /\/api\/v1\/cms\/team/)
})

test('CMS operations statuses and audit logs avoid unknown states', () => {
  const commandCenter = read('../src/components/cms/SystemCommandCenterPanel.vue')
  const cmsView = read('../src/views/WebsiteContentManagerView.vue')
  const operatorHandler = read('../../backend/internal/handlers/operator.go')

  assert.match(commandCenter, /displayServiceStatus/)
  assert.match(commandCenter, /return service\.status === 'running' \? 'running' : 'idle'/)
  assert.match(commandCenter, /displayServiceHealth/)
  assert.match(commandCenter, /docker_available === false/)
  assert.match(operatorHandler, /"status": "idle"/)
  assert.match(operatorHandler, /normalizeServiceState/)
  assert.doesNotMatch(operatorHandler, /"status": "unknown"/)

  assert.match(cmsView, /auditPage = ref\(1\)/)
  assert.match(cmsView, /auditPerPage = ref\(20\)/)
  assert.match(cmsView, /auditTotalPages/)
  assert.match(cmsView, /changeAuditPage/)
  assert.match(cmsView, /Page \{\{ auditPage \}\} of \{\{ auditTotalPages \}\}/)
  assert.match(cmsView, /adminListAuditLogs\(\{ page:auditPage\.value, per_page:auditPerPage\.value/)
})
