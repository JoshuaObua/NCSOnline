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
