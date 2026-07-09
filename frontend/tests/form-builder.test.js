import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import test from 'node:test'
import {
  buildTemplatePayload,
  normalizeFormField,
  normalizePagedResponse,
  parseSubmissionAnswers,
} from '../src/utils/formBuilder.js'

const read = relativePath => readFileSync(fileURLToPath(new URL(relativePath, import.meta.url)), 'utf8')

test('form builder serializes ordered fields, unique keys, options, and file config', () => {
  const payload = buildTemplatePayload({
    department_id: 'role-admissions',
    title: '  Club Registration  ',
    description: '  Register a club. ',
    price_ugx: '50000',
    status: 'OPEN',
    fields: [
      { label: 'Club Name', field_type: 'short_text', is_required: true },
      { label: 'Club Name', field_type: 'dropdown', options_text: 'Senior\nJunior\n\n' },
      { label: 'Evidence', field_type: 'file', accepted_types: 'application/pdf, image/png' },
    ],
  })

  assert.equal(payload.title, 'Club Registration')
  assert.equal(payload.description, 'Register a club.')
  assert.equal(payload.price_ugx, 50000)
  assert.deepEqual(payload.fields.map(field => field.field_key), ['club-name', 'club-name-2', 'evidence'])
  assert.deepEqual(payload.fields[1].config.options, ['Senior', 'Junior'])
  assert.deepEqual(payload.fields[2].config.accept, ['application/pdf', 'image/png'])
})

test('form builder normalizes API config and paged submissions', () => {
  const field = normalizeFormField({
    field_type: 'radio',
    config: '{"options":["Yes","No"],"accept":["image/png"]}',
  })
  assert.equal(field.options_text, 'Yes\nNo')
  assert.equal(field.accepted_types, 'image/png')

  const page = normalizePagedResponse({ success: true, data: [{ id: 'sub-1' }], meta: { total: 1 } })
  assert.deepEqual(page.items, [{ id: 'sub-1' }])
  assert.equal(page.meta.total, 1)
  assert.deepEqual(parseSubmissionAnswers('{"club-name":"Kampala Stars"}'), { 'club-name': 'Kampala Stars' })
})

test('portal registers the form builder sidebar, permissions, API calls, and full panel', () => {
  const portal = read('../src/views/WebsiteContentManagerView.vue')
  const panel = read('../src/components/cms/FormBuilderPanel.vue')
  const api = read('../src/api/forms.js')
  const migration = read('../../backend/migrations/057_form_field_active_key_uniqueness.sql')

  assert.match(portal, /id:'form-builder', label:'Custom Form Builder'/)
  assert.match(portal, /<FormBuilderPanel v-else-if="active === 'form-builder'"/)
  assert.match(portal, /applications:admin:read/)
  assert.match(panel, /Templates/)
  assert.match(panel, /Builder/)
  assert.match(panel, /Submissions/)
  assert.match(panel, /adminCreateForm/)
  assert.match(panel, /adminUpdateForm/)
  assert.match(panel, /adminReviewSubmission/)
  assert.match(panel, /adminVerifySubmissionPayment/)
  assert.match(api, /\/api\/v1\/admin\/forms/)
  assert.match(api, /\/api\/v1\/admin\/forms\/submissions/)
  assert.match(migration, /DROP CONSTRAINT IF EXISTS form_fields_template_id_field_key_key/)
  assert.match(migration, /WHERE deleted_at IS NULL/)
})
