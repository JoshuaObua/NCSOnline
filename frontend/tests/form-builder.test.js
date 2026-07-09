import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import test from 'node:test'
import {
  buildSectionSteps,
  buildTemplatePayload,
  normalizeFormField,
  normalizeFormSections,
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
    sections: [
      { id: 'club-details', title: 'Club Details', subtitle: 'Identity', description: 'Tell us about the club.' },
      { id: 'documents', title: 'Documents', subtitle: 'Uploads', description: 'Attach supporting files.' },
    ],
    fields: [
      { label: 'Club Name', field_type: 'short_text', is_required: true, section_id: 'club-details' },
      { label: 'Club Name', field_type: 'dropdown', options_text: 'Senior\nJunior\n\n', section_id: 'club-details' },
      { label: 'Evidence', field_type: 'file', accepted_types: 'application/pdf, image/png', section_id: 'documents' },
    ],
  })

  assert.equal(payload.title, 'Club Registration')
  assert.equal(payload.description, 'Register a club.')
  assert.equal(payload.price_ugx, 50000)
  assert.deepEqual(payload.sections.map(section => section.title), ['Club Details', 'Documents'])
  assert.deepEqual(payload.fields.map(field => field.field_key), ['club-name', 'club-name-2', 'evidence'])
  assert.deepEqual(payload.fields.map(field => field.config.section_id), ['club-details', 'club-details', 'documents'])
  assert.deepEqual(payload.fields[1].config.options, ['Senior', 'Junior'])
  assert.deepEqual(payload.fields[2].config.accept, ['application/pdf', 'image/png'])
})

test('form builder normalizes sections, API config, wizard steps, and paged submissions', () => {
  const field = normalizeFormField({
    field_type: 'radio',
    config: '{"options":["Yes","No"],"accept":["image/png"],"section_id":"eligibility"}',
  })
  assert.equal(field.options_text, 'Yes\nNo')
  assert.equal(field.accepted_types, 'image/png')
  assert.equal(field.section_id, 'eligibility')

  const sections = normalizeFormSections([{ id: 'eligibility', title: 'Eligibility' }], [field])
  const steps = buildSectionSteps({ sections, fields: [field] })
  assert.equal(steps.length, 1)
  assert.equal(steps[0].title, 'Eligibility')
  assert.equal(steps[0].fields[0].field_type, 'radio')

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
  const sectionsMigration = read('../../backend/migrations/058_form_template_sections.sql')

  assert.match(portal, /id:'form-builder', label:'Custom Form Builder'/)
  assert.match(portal, /<FormBuilderPanel v-else-if="active === 'form-builder'"/)
  assert.match(portal, /applications:admin:read/)
  assert.match(panel, /Templates/)
  assert.match(panel, /Builder/)
  assert.match(panel, /Submissions/)
  assert.match(panel, /Sections and Fields/)
  assert.match(panel, /Discard structure/)
  assert.match(panel, /adminCreateForm/)
  assert.match(panel, /adminUpdateForm/)
  assert.match(panel, /adminReviewSubmission/)
  assert.match(panel, /adminVerifySubmissionPayment/)
  assert.match(api, /\/api\/v1\/admin\/forms/)
  assert.match(api, /\/api\/v1\/admin\/forms\/submissions/)
  assert.match(migration, /DROP CONSTRAINT IF EXISTS form_fields_template_id_field_key_key/)
  assert.match(migration, /WHERE deleted_at IS NULL/)
  assert.match(sectionsMigration, /ADD COLUMN IF NOT EXISTS sections JSONB/)
  assert.match(sectionsMigration, /jsonb_typeof\(sections\) = 'array'/)
})
