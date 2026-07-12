<template>
  <section class="form-builder">
    <div class="builder-tabs" role="tablist" aria-label="Form builder views">
      <button type="button" :class="{ active: view === 'templates' }" @click="showTemplates">
        <i class="fas fa-list" aria-hidden="true"></i>
        Templates
      </button>
      <button type="button" :class="{ active: view === 'builder' }" @click="startCreate">
        <i class="fas fa-edit" aria-hidden="true"></i>
        Builder
      </button>
      <button type="button" :class="{ active: view === 'submissions' }" @click="showSubmissions">
        <i class="fas fa-inbox" aria-hidden="true"></i>
        Submissions
      </button>
    </div>

    <p v-if="localError" class="builder-alert error">{{ localError }}</p>
    <p v-if="loading" class="builder-alert">Loading form builder data...</p>

    <template v-if="view === 'templates'">
      <div class="builder-toolbar">
        <div>
          <h2>Form Templates</h2>
          <p>Create and manage department-owned application forms.</p>
        </div>
        <div class="toolbar-actions">
          <label>
            <span class="sr-only">Filter templates by status</span>
            <select v-model="templateStatus" @change="loadTemplates">
              <option value="">All statuses</option>
              <option v-for="status in templateStatuses" :key="status" :value="status">{{ titleCase(status) }}</option>
            </select>
          </label>
          <button type="button" class="primary-action" @click="startCreate">
            <i class="fas fa-plus" aria-hidden="true"></i>
            New form
          </button>
        </div>
      </div>

      <div class="table-panel">
        <table>
          <thead>
            <tr><th>Form</th><th>Department</th><th>Status</th><th>Fee</th><th>Updated</th><th></th></tr>
          </thead>
          <tbody>
            <tr v-for="template in templates" :key="template.id">
              <td>
                <strong>{{ template.title }}</strong>
                <small>/{{ template.slug }}</small>
              </td>
              <td>{{ template.department_name || departmentName(template.department_id) }}</td>
              <td><span class="status-badge" :class="statusClass(template.status)">{{ titleCase(template.status) }}</span></td>
              <td>{{ template.price_ugx > 0 ? formatCurrency(template.price_ugx) : 'Free' }}</td>
              <td>{{ formatDate(template.updated_at) }}</td>
              <td class="row-actions">
                <button type="button" title="Edit form" aria-label="Edit form" @click="editTemplate(template.id)"><i class="fas fa-pen"></i></button>
                <button type="button" title="View submissions" aria-label="View submissions" @click="openTemplateSubmissions(template.id)"><i class="fas fa-inbox"></i></button>
              </td>
            </tr>
            <tr v-if="!loading && !templates.length"><td colspan="6" class="empty-state">No form templates match this filter.</td></tr>
          </tbody>
        </table>
      </div>
    </template>

    <template v-else-if="view === 'builder'">
      <div class="builder-toolbar">
        <div>
          <button type="button" class="back-action" @click="showTemplates"><i class="fas fa-arrow-left"></i> Templates</button>
          <h2>{{ editor.id ? 'Edit Form' : 'Create Form' }}</h2>
          <p>{{ editor.id ? `Editing /${editor.slug}` : 'Build a reusable application form and assign it to a department.' }}</p>
        </div>
        <div class="toolbar-actions">
          <button v-if="editor.id" type="button" class="danger-action" :disabled="saving" @click="archiveTemplate">
            <i class="fas fa-archive"></i>
            Archive
          </button>
          <button type="button" class="primary-action" :disabled="saving" @click="saveTemplate">
            <i class="fas fa-save"></i>
            {{ saving ? 'Saving...' : 'Save form' }}
          </button>
        </div>
      </div>

      <div class="editor-layout">
        <div class="editor-column">
          <section class="settings-panel">
            <div class="panel-heading"><h3>Form Details</h3></div>
            <div class="form-grid">
              <label>Title
                <input v-model="editor.title" type="text" maxlength="255" placeholder="Application form title" />
              </label>
              <label>Department
                <select v-model="editor.department_id">
                  <option value="" disabled>Select department</option>
                  <option v-for="department in departments" :key="department.id" :value="department.id">{{ department.name }}</option>
                </select>
              </label>
              <label>Application fee (UGX)
                <input v-model.number="editor.price_ugx" type="number" min="0" step="1000" />
              </label>
              <fieldset class="status-control">
                <legend>Status</legend>
                <div>
                  <button v-for="status in templateStatuses" :key="status" type="button" :class="{ active: editor.status === status }" @click="editor.status = status">
                    {{ titleCase(status) }}
                  </button>
                </div>
              </fieldset>
              <label class="wide">Description
                <textarea v-model="editor.description" rows="4" placeholder="Explain who should use this form and what it is for."></textarea>
              </label>
              <div class="wide">
                <span class="field-label">Banner image</span>
                <DropzoneUpload v-model="editor.banner_image_url" label="form banner" @error="handleError" />
              </div>
            </div>
          </section>

          <section class="fields-section sections-section">
            <div class="panel-heading fields-heading">
              <div>
                <h3>Sections and Fields</h3>
                <p>{{ editor.sections.length }} section{{ editor.sections.length === 1 ? '' : 's' }} · {{ editor.fields.length }} field{{ editor.fields.length === 1 ? '' : 's' }}</p>
              </div>
              <div class="add-field">
                <select v-model="newFieldType" aria-label="New field type">
                  <option v-for="type in FIELD_TYPES" :key="type.value" :value="type.value">{{ type.label }}</option>
                </select>
                <button type="button" class="secondary-action" @click="addSection"><i class="fas fa-layer-group"></i> Add section</button>
                <button v-if="editor.fields.length || editor.sections.length > 1" type="button" class="danger-action" @click="rebuildTemplateStructure"><i class="fas fa-eraser"></i> Discard structure</button>
              </div>
            </div>

            <article v-for="(formSection, sectionIndex) in editor.sections" :key="formSection.id" class="section-card">
              <header class="section-card-header">
                <span class="section-number">Step {{ sectionIndex + 1 }}</span>
                <strong>{{ formSection.title || `Section ${sectionIndex + 1}` }}</strong>
                <div class="field-actions">
                  <button type="button" title="Move section up" aria-label="Move section up" :disabled="sectionIndex === 0" @click="moveSection(sectionIndex, -1)"><i class="fas fa-arrow-up"></i></button>
                  <button type="button" title="Move section down" aria-label="Move section down" :disabled="sectionIndex === editor.sections.length - 1" @click="moveSection(sectionIndex, 1)"><i class="fas fa-arrow-down"></i></button>
                  <button type="button" class="remove-field" title="Remove section" aria-label="Remove section" :disabled="editor.sections.length === 1" @click="removeSection(sectionIndex)"><i class="fas fa-trash"></i></button>
                </div>
              </header>
              <div class="form-grid compact section-settings">
                <label>Section title
                  <input v-model="formSection.title" type="text" maxlength="140" placeholder="Application Details" />
                </label>
                <label>Subtitle
                  <input v-model="formSection.subtitle" type="text" maxlength="180" placeholder="Step summary shown to applicants" />
                </label>
                <label class="wide">Description
                  <textarea v-model="formSection.description" rows="2" placeholder="Guidance shown at the top of this step."></textarea>
                </label>
              </div>
              <div class="section-field-toolbar">
                <span>{{ sectionFields(formSection.id).length }} field{{ sectionFields(formSection.id).length === 1 ? '' : 's' }}</span>
                <button type="button" class="primary-action" @click="addField(formSection.id)"><i class="fas fa-plus"></i> Add field to step</button>
              </div>
              <article v-for="field in sectionFields(formSection.id)" :key="field.id" class="field-card section-field-card">
                <header>
                  <span class="field-number">{{ fieldGlobalIndex(field) + 1 }}</span>
                  <strong>{{ field.label || fieldTypeLabel(field.field_type) }}</strong>
                  <div class="field-actions">
                    <button type="button" title="Move up" aria-label="Move field up" :disabled="isFirstFieldInSection(field)" @click="moveFieldInSection(field, -1)"><i class="fas fa-arrow-up"></i></button>
                    <button type="button" title="Move down" aria-label="Move field down" :disabled="isLastFieldInSection(field)" @click="moveFieldInSection(field, 1)"><i class="fas fa-arrow-down"></i></button>
                    <button type="button" title="Duplicate" aria-label="Duplicate field" @click="duplicateField(fieldGlobalIndex(field))"><i class="fas fa-copy"></i></button>
                    <button type="button" class="remove-field" title="Remove" aria-label="Remove field" @click="removeField(fieldGlobalIndex(field))"><i class="fas fa-trash"></i></button>
                  </div>
                </header>
                <div class="form-grid compact">
                  <label>Field type
                    <select v-model="field.field_type">
                      <option v-for="type in FIELD_TYPES" :key="type.value" :value="type.value">{{ type.label }}</option>
                    </select>
                  </label>
                  <label>Step
                    <select v-model="field.section_id" @change="syncFieldsToSectionOrder">
                      <option v-for="sectionOption in editor.sections" :key="sectionOption.id" :value="sectionOption.id">{{ sectionOption.title || 'Untitled section' }}</option>
                    </select>
                  </label>
                  <label>Label
                    <input v-model="field.label" type="text" maxlength="255" @blur="fillFieldKey(field)" />
                  </label>
                  <label>Field key
                    <input v-model="field.field_key" type="text" maxlength="80" placeholder="generated-from-label" />
                  </label>
                  <label>Placeholder
                    <input v-model="field.placeholder" type="text" maxlength="255" />
                  </label>
                  <label class="wide">Help text
                    <input v-model="field.help_text" type="text" placeholder="Optional guidance shown below the field" />
                  </label>
                  <label v-if="hasOptions(field.field_type)" class="wide">Options, one per line
                    <textarea v-model="field.options_text" rows="4" placeholder="Option one&#10;Option two"></textarea>
                  </label>
                  <label v-if="['file', 'image'].includes(field.field_type)" class="wide">Accepted file types
                    <input v-model="field.accepted_types" type="text" placeholder="image/png,image/jpeg,application/pdf" />
                  </label>
                  <label class="required-toggle">
                    <input v-model="field.is_required" type="checkbox" />
                    Required field
                  </label>
                </div>
              </article>
              <div v-if="!sectionFields(formSection.id).length" class="empty-state fields-empty">Add fields to this step or move existing fields here.</div>
            </article>
            <div v-if="!editor.fields.length" class="empty-state fields-empty">Add the first field to begin building this form.</div>
          </section>
        </div>

        <aside class="preview-panel">
          <div v-if="editor.banner_image_url" class="preview-banner"><img :src="mediaUrl(editor.banner_image_url)" alt="" /></div>
          <div class="preview-body">
            <span class="preview-kicker">{{ departmentName(editor.department_id) || 'NCS application' }}</span>
            <h3>{{ editor.title || 'Untitled form' }}</h3>
            <p>{{ editor.description || 'The form description will appear here.' }}</p>
            <div v-if="Number(editor.price_ugx) > 0" class="fee-notice">Application fee: {{ formatCurrency(editor.price_ugx) }}</div>
            <div class="preview-step-list">
              <span v-for="(step, stepIndex) in previewSections" :key="step.id">Step {{ stepIndex + 1 }}</span>
            </div>
            <section v-for="(step, stepIndex) in previewSections" :key="`preview-section-${step.id}`" class="preview-section">
              <small>Step {{ stepIndex + 1 }}</small>
              <h4>{{ step.title }}</h4>
              <p v-if="step.subtitle">{{ step.subtitle }}</p>
              <p v-if="step.description">{{ step.description }}</p>
              <div v-for="field in step.fields" :key="`preview-${field.id}`" class="preview-field">
                <label>{{ field.label || 'Untitled field' }} <span v-if="field.is_required">*</span></label>
                <textarea v-if="field.field_type === 'long_text'" disabled rows="3" :placeholder="field.placeholder"></textarea>
                <select v-else-if="field.field_type === 'dropdown'" disabled><option>{{ field.placeholder || 'Select an option' }}</option></select>
                <div v-else-if="['radio', 'checkbox'].includes(field.field_type)" class="preview-options">
                  <label v-for="option in fieldOptions(field)" :key="option"><input :type="field.field_type" disabled /> {{ option }}</label>
                  <small v-if="!fieldOptions(field).length">Add options to preview this field.</small>
                </div>
                <input v-else :type="previewInputType(field.field_type)" disabled :placeholder="field.placeholder" />
                <small v-if="field.help_text">{{ field.help_text }}</small>
              </div>
            </section>
          </div>
        </aside>
      </div>
    </template>

    <template v-else>
      <div class="builder-toolbar">
        <div><h2>Form Submissions</h2><p>Review applications and verify uploaded payment references.</p></div>
        <button type="button" class="secondary-action" @click="loadSubmissions"><i class="fas fa-sync"></i> Refresh</button>
      </div>

      <div class="filter-panel">
        <label>Form
          <select v-model="submissionFilters.template_id" @change="resetSubmissionPage">
            <option value="">All forms</option>
            <option v-for="template in templates" :key="template.id" :value="template.id">{{ template.title }}</option>
          </select>
        </label>
        <label>Status
          <select v-model="submissionFilters.status" @change="resetSubmissionPage">
            <option value="">All statuses</option>
            <option v-for="status in submissionStatuses" :key="status" :value="status">{{ titleCase(status) }}</option>
          </select>
        </label>
      </div>

      <div class="submissions-layout" :class="{ detailed: selectedSubmission }">
        <div class="table-panel">
          <table>
            <thead><tr><th>Reference</th><th>Applicant</th><th>Form</th><th>Status</th><th>Payment</th><th>Updated</th></tr></thead>
            <tbody>
              <tr v-for="submission in submissions" :key="submission.id" class="clickable-row" @click="viewSubmission(submission.id)">
                <td><strong>{{ submission.submission_reference || 'Draft' }}</strong></td>
                <td>{{ submission.applicant_name || submission.applicant_email || 'Applicant' }}</td>
                <td>{{ submission.template_title }}</td>
                <td><span class="status-badge" :class="statusClass(submission.status)">{{ titleCase(submission.status) }}</span></td>
                <td>{{ titleCase(submission.payment_status) }}</td>
                <td>{{ formatDate(submission.updated_at) }}</td>
              </tr>
              <tr v-if="!loadingSubmissions && !submissions.length"><td colspan="6" class="empty-state">No submissions match these filters.</td></tr>
            </tbody>
          </table>
          <div class="pagination">
            <button type="button" :disabled="submissionPage <= 1" @click="changeSubmissionPage(-1)"><i class="fas fa-chevron-left"></i></button>
            <span>Page {{ submissionPage }} of {{ submissionTotalPages }}</span>
            <button type="button" :disabled="submissionPage >= submissionTotalPages" @click="changeSubmissionPage(1)"><i class="fas fa-chevron-right"></i></button>
          </div>
        </div>

        <aside v-if="selectedSubmission" class="submission-detail">
          <header>
            <div><span>Submission</span><h3>{{ selectedSubmission.submission_reference || 'Draft submission' }}</h3></div>
            <button type="button" title="Close details" aria-label="Close submission details" @click="selectedSubmission = null"><i class="fas fa-times"></i></button>
          </header>
          <dl class="submission-meta">
            <div><dt>Applicant</dt><dd>{{ selectedSubmission.applicant_name || 'Not provided' }}</dd></div>
            <div><dt>Email</dt><dd>{{ selectedSubmission.applicant_email || 'Not provided' }}</dd></div>
            <div><dt>Status</dt><dd>{{ titleCase(selectedSubmission.status) }}</dd></div>
            <div><dt>Payment</dt><dd>{{ titleCase(selectedSubmission.payment_status) }}</dd></div>
            <div v-if="selectedSubmission.payment_reference"><dt>Payment reference</dt><dd>{{ selectedSubmission.payment_reference }}</dd></div>
            <div v-if="selectedSubmission.payment_amount_ugx"><dt>Amount</dt><dd>{{ formatCurrency(selectedSubmission.payment_amount_ugx) }}</dd></div>
          </dl>
          <section class="answer-list">
            <h4>Answers</h4>
            <div v-for="answer in submissionAnswers" :key="answer.key">
              <strong>{{ answer.label }}</strong>
              <span>{{ answer.value }}</span>
            </div>
            <p v-if="!submissionAnswers.length" class="empty-state">No answers have been saved.</p>
          </section>
          <section class="review-panel">
            <label>Decision
              <select v-model="review.status">
                <option value="UNDER_REVIEW">Under review</option>
                <option value="NEEDS_INFORMATION">Needs information</option>
                <option value="APPROVED">Approved</option>
                <option value="REJECTED">Rejected</option>
              </select>
            </label>
            <label>Review notes
              <textarea v-model="review.notes" rows="4" placeholder="Add notes for this review decision."></textarea>
            </label>
            <div class="review-actions">
              <button v-if="selectedSubmission.payment_status === 'PROOF_UPLOADED'" type="button" class="secondary-action" :disabled="reviewing" @click="verifyPayment">
                <i class="fas fa-check-circle"></i> Verify payment
              </button>
              <button type="button" class="primary-action" :disabled="reviewing || selectedSubmission.status === 'DRAFT'" @click="submitReview">
                <i class="fas fa-gavel"></i> Save review
              </button>
            </div>
          </section>
        </aside>
      </div>
    </template>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import {
  FIELD_TYPES,
  FIELD_TYPES_WITH_OPTIONS,
  adminCreateForm,
  adminDeleteForm,
  adminGetForm,
  adminGetSubmission,
  adminListForms,
  adminListSubmissions,
  adminReviewSubmission,
  adminUpdateForm,
  adminVerifySubmissionPayment,
  listDepartments,
} from '@/api/forms.js'
import { mediaUrl } from '@/api/client.js'
import DropzoneUpload from '@/components/cms/DropzoneUpload.vue'
import {
  buildSectionSteps,
  buildTemplatePayload,
  normalizeFormField,
  normalizeFormSection,
  normalizeFormSections,
  normalizePagedResponse,
  parseSubmissionAnswers,
  slugifyFieldKey,
} from '@/utils/formBuilder.js'

const emit = defineEmits(['message', 'error'])
const router = useRouter()

const templateStatuses = ['DRAFT', 'OPEN', 'CLOSED', 'ARCHIVED']
const submissionStatuses = ['DRAFT', 'PENDING_PAYMENT', 'SUBMITTED', 'UNDER_REVIEW', 'NEEDS_INFORMATION', 'COMPLETE', 'APPROVED', 'REJECTED']
const view = ref('templates')
const templates = ref([])
const departments = ref([])
const submissions = ref([])
const selectedSubmission = ref(null)
const selectedSubmissionTemplate = ref(null)
const loading = ref(false)
const loadingSubmissions = ref(false)
const saving = ref(false)
const reviewing = ref(false)
const localError = ref('')
const templateStatus = ref('')
const newFieldType = ref('short_text')
const submissionPage = ref(1)
const submissionMeta = ref({})
const submissionFilters = reactive({ template_id: '', status: '' })
const review = reactive({ status: 'UNDER_REVIEW', notes: '' })

const blankEditor = () => ({
  id: '',
  slug: '',
  department_id: '',
  title: '',
  description: '',
  banner_image_url: '',
  price_ugx: 0,
  status: 'DRAFT',
  sections: [normalizeFormSection({}, 0)],
  fields: [],
})
const editor = reactive(blankEditor())

const submissionTotalPages = computed(() => Math.max(1, Math.ceil(Number(submissionMeta.value.total || 0) / Number(submissionMeta.value.per_page || 20))))
const previewSections = computed(() => buildSectionSteps({ sections: editor.sections, fields: editor.fields }))
const submissionAnswers = computed(() => {
  if (!selectedSubmission.value) return []
  const answers = parseSubmissionAnswers(selectedSubmission.value.answers)
  const fields = selectedSubmissionTemplate.value?.fields || []
  const labels = new Map(fields.map(field => [field.field_key, field.label]))
  return Object.entries(answers).map(([key, value]) => ({
    key,
    label: labels.get(key) || titleCase(key),
    value: formatAnswer(value),
  }))
})

function titleCase(value) {
  return String(value || '').replace(/[_-]+/g, ' ').toLowerCase().replace(/\b\w/g, character => character.toUpperCase())
}
function formatDate(value) {
  return value ? new Date(value).toLocaleDateString('en-UG', { year: 'numeric', month: 'short', day: 'numeric' }) : ''
}
function formatCurrency(value) {
  return new Intl.NumberFormat('en-UG', { style: 'currency', currency: 'UGX', maximumFractionDigits: 0 }).format(Number(value || 0))
}
function formatAnswer(value) {
  if (Array.isArray(value)) return value.join(', ')
  if (value && typeof value === 'object') return JSON.stringify(value)
  if (value === true) return 'Yes'
  if (value === false) return 'No'
  return String(value ?? '')
}
function statusClass(status) {
  const value = String(status || '').toLowerCase()
  if (['open', 'approved', 'paid'].includes(value)) return 'success'
  if (['rejected', 'archived'].includes(value)) return 'danger'
  if (['closed', 'needs_information', 'pending_payment', 'proof_uploaded'].includes(value)) return 'warning'
  return 'info'
}
function departmentName(id) {
  return departments.value.find(department => department.id === id)?.name || ''
}
function fieldTypeLabel(type) {
  return FIELD_TYPES.find(item => item.value === type)?.label || titleCase(type)
}
function hasOptions(type) {
  return FIELD_TYPES_WITH_OPTIONS.includes(type)
}
function fieldOptions(field) {
  return String(field.options_text || '').split(/\r?\n/).map(option => option.trim()).filter(Boolean)
}
function previewInputType(type) {
  if (type === 'email') return 'email'
  if (type === 'phone') return 'tel'
  if (type === 'number') return 'number'
  if (type === 'date') return 'date'
  if (['file', 'image'].includes(type)) return 'file'
  return 'text'
}
function handleError(error) {
  localError.value = error?.response?.data?.error?.message || error?.message || 'Action failed'
  emit('error', error)
}
function clearError() {
  localError.value = ''
}

async function loadTemplates() {
  loading.value = true
  clearError()
  try {
    templates.value = await adminListForms(templateStatus.value)
  } catch (error) {
    handleError(error)
  } finally {
    loading.value = false
  }
}
async function loadDepartments() {
  try {
    departments.value = await listDepartments()
  } catch (error) {
    handleError(error)
  }
}
async function editTemplate(id) {
  loading.value = true
  clearError()
  try {
    const template = await adminGetForm(id)
    const fields = (template.fields || []).map(normalizeFormField)
    const sections = normalizeFormSections(template.sections || [], fields)
    assignFieldsToSections(fields, sections)
    Object.assign(editor, blankEditor(), template, {
      sections,
      fields,
    })
    syncFieldsToSectionOrder()
    view.value = 'builder'
  } catch (error) {
    handleError(error)
  } finally {
    loading.value = false
  }
}
function startCreate() {
  Object.assign(editor, blankEditor(), {
    department_id: departments.value.length === 1 ? departments.value[0].id : '',
  })
  newFieldType.value = 'short_text'
  clearError()
  view.value = 'builder'
}
function showTemplates() {
  selectedSubmission.value = null
  view.value = 'templates'
}
async function showSubmissions() {
  view.value = 'submissions'
  await loadSubmissions()
}
function newSection(index = editor.sections.length) {
  return normalizeFormSection({
    id: `section-${Date.now()}-${index + 1}`,
    title: index === 0 ? 'Application Details' : `Section ${index + 1}`,
  }, index)
}
function assignFieldsToSections(fields, sections) {
  const fallback = sections[0]?.id || 'section-1'
  const ids = new Set(sections.map(section => section.id))
  fields.forEach(field => {
    if (!ids.has(field.section_id)) field.section_id = fallback
  })
}
function syncFieldsToSectionOrder() {
  const order = new Map(editor.sections.map((section, index) => [section.id, index]))
  editor.fields.sort((a, b) => (order.get(a.section_id) ?? 999) - (order.get(b.section_id) ?? 999))
}
function sectionFields(sectionId) {
  return editor.fields.filter(field => field.section_id === sectionId)
}
function fieldGlobalIndex(field) {
  return editor.fields.indexOf(field)
}
function isFirstFieldInSection(field) {
  return sectionFields(field.section_id)[0] === field
}
function isLastFieldInSection(field) {
  const fields = sectionFields(field.section_id)
  return fields[fields.length - 1] === field
}
function moveFieldInSection(field, direction) {
  const fields = sectionFields(field.section_id)
  const index = fields.indexOf(field)
  const target = fields[index + direction]
  if (!target) return
  const currentIndex = fieldGlobalIndex(field)
  const targetIndex = fieldGlobalIndex(target)
  editor.fields.splice(currentIndex, 1, target)
  editor.fields.splice(targetIndex, 1, field)
}
function addSection() {
  editor.sections.push(newSection())
}
function moveSection(index, direction) {
  const target = index + direction
  if (target < 0 || target >= editor.sections.length) return
  const [section] = editor.sections.splice(index, 1)
  editor.sections.splice(target, 0, section)
  syncFieldsToSectionOrder()
}
function removeSection(index) {
  if (editor.sections.length === 1) return
  const [removed] = editor.sections.splice(index, 1)
  const targetSection = editor.sections[Math.max(0, index - 1)] || editor.sections[0]
  editor.fields.forEach(field => {
    if (field.section_id === removed.id) field.section_id = targetSection.id
  })
  syncFieldsToSectionOrder()
}
async function rebuildTemplateStructure() {
  const result = await Swal.fire({
    title: 'Discard this form structure?',
    text: 'All sections and fields in the editor will be cleared so you can rebuild the template.',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'Discard and rebuild',
    confirmButtonColor: '#fc544b',
  })
  if (!result.isConfirmed) return
  editor.sections = [newSection(0)]
  editor.fields = []
}
function addField(sectionId = editor.sections[0]?.id) {
  editor.fields.push(normalizeFormField({
    field_type: newFieldType.value,
    label: fieldTypeLabel(newFieldType.value),
    section_id: sectionId,
    config: { section_id: sectionId },
  }, editor.fields.length))
  syncFieldsToSectionOrder()
}
function fillFieldKey(field) {
  if (!field.field_key) field.field_key = slugifyFieldKey(field.label)
}
function moveField(index, direction) {
  const target = index + direction
  if (target < 0 || target >= editor.fields.length) return
  const [field] = editor.fields.splice(index, 1)
  editor.fields.splice(target, 0, field)
}
function duplicateField(index) {
  const source = editor.fields[index]
  const duplicate = normalizeFormField({
    ...source,
    id: '',
    label: `${source.label || 'Field'} copy`,
    field_key: '',
    config: { ...source.config },
    section_id: source.section_id,
  }, editor.fields.length)
  duplicate.options_text = source.options_text
  duplicate.accepted_types = source.accepted_types
  editor.fields.splice(index + 1, 0, duplicate)
}
function removeField(index) {
  editor.fields.splice(index, 1)
}
async function saveTemplate() {
  clearError()
  syncFieldsToSectionOrder()
  const payload = buildTemplatePayload(editor)
  if (!payload.title) return handleError(new Error('Form title is required.'))
  if (!payload.department_id) return handleError(new Error('Select a department.'))
  const untitledSection = payload.sections.findIndex(section => !section.title)
  if (untitledSection >= 0) return handleError(new Error(`Section ${untitledSection + 1} needs a title.`))
  if (!payload.fields.length) return handleError(new Error('Add at least one field before saving.'))
  const missingLabel = payload.fields.findIndex(field => !field.label)
  if (missingLabel >= 0) return handleError(new Error(`Field ${missingLabel + 1} needs a label.`))
  const optionless = payload.fields.findIndex(field => hasOptions(field.field_type) && !field.config.options?.length)
  if (optionless >= 0) return handleError(new Error(`Field ${optionless + 1} needs at least one option.`))

  saving.value = true
  try {
    const wasEditing = !!editor.id
    const saved = wasEditing
      ? await adminUpdateForm(editor.id, payload)
      : await adminCreateForm(payload)
    Object.assign(editor, blankEditor(), saved, {
      sections: normalizeFormSections(saved.sections || [], saved.fields || []),
      fields: (saved.fields || []).map(normalizeFormField),
    })
    assignFieldsToSections(editor.fields, editor.sections)
    syncFieldsToSectionOrder()
    await loadTemplates()
    emit('message', wasEditing ? 'Form template saved' : 'Form template created')
  } catch (error) {
    handleError(error)
  } finally {
    saving.value = false
  }
}
async function archiveTemplate() {
  const result = await Swal.fire({
    title: 'Archive this form?',
    text: 'Applicants will no longer see it as an open form.',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'Archive form',
    confirmButtonColor: '#fc544b',
  })
  if (!result.isConfirmed) return
  saving.value = true
  try {
    await adminDeleteForm(editor.id)
    await loadTemplates()
    emit('message', 'Form template archived')
    showTemplates()
  } catch (error) {
    handleError(error)
  } finally {
    saving.value = false
  }
}

async function loadSubmissions() {
  loadingSubmissions.value = true
  clearError()
  try {
    const response = await adminListSubmissions({
      page: submissionPage.value,
      per_page: 20,
      ...(submissionFilters.template_id ? { template_id: submissionFilters.template_id } : {}),
      ...(submissionFilters.status ? { status: submissionFilters.status } : {}),
    })
    const page = normalizePagedResponse(response)
    submissions.value = page.items
    submissionMeta.value = page.meta
  } catch (error) {
    handleError(error)
  } finally {
    loadingSubmissions.value = false
  }
}
async function openTemplateSubmissions(templateId) {
  submissionFilters.template_id = templateId
  submissionPage.value = 1
  await showSubmissions()
}
async function resetSubmissionPage() {
  submissionPage.value = 1
  selectedSubmission.value = null
  await loadSubmissions()
}
async function changeSubmissionPage(direction) {
  submissionPage.value += direction
  selectedSubmission.value = null
  await loadSubmissions()
}
async function openSubmission(id) {
  reviewing.value = true
  clearError()
  try {
    const submission = await adminGetSubmission(id)
    const template = await adminGetForm(submission.template_id)
    selectedSubmission.value = submission
    selectedSubmissionTemplate.value = template
    review.status = ['UNDER_REVIEW', 'NEEDS_INFORMATION', 'APPROVED', 'REJECTED'].includes(submission.status)
      ? submission.status
      : 'UNDER_REVIEW'
    review.notes = submission.review_notes || ''
  } catch (error) {
    handleError(error)
  } finally {
    reviewing.value = false
  }
}
function viewSubmission(id) {
  router.push({ name: 'AdminApplicationDetail', params: { id }, query: { source: 'custom' } })
}
async function submitReview() {
  if (!selectedSubmission.value) return
  reviewing.value = true
  try {
    await adminReviewSubmission(selectedSubmission.value.id, review.status, review.notes)
    emit('message', 'Submission review saved')
    await loadSubmissions()
    await openSubmission(selectedSubmission.value.id)
  } catch (error) {
    handleError(error)
  } finally {
    reviewing.value = false
  }
}
async function verifyPayment() {
  if (!selectedSubmission.value) return
  reviewing.value = true
  try {
    await adminVerifySubmissionPayment(selectedSubmission.value.id)
    emit('message', 'Submission payment verified')
    await loadSubmissions()
    await openSubmission(selectedSubmission.value.id)
  } catch (error) {
    handleError(error)
  } finally {
    reviewing.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadDepartments(), loadTemplates()])
})
</script>

<style scoped>
.form-builder{display:grid;gap:20px;min-width:0}.builder-tabs{display:flex;gap:6px;border-bottom:1px solid #e4e6fc}.builder-tabs button{display:inline-flex;align-items:center;gap:8px;border:0;border-bottom:3px solid transparent;background:transparent;color:#6c757d;padding:12px 16px;font-size:13px;font-weight:700}.builder-tabs button.active{border-bottom-color:#6777ef;color:#6777ef}.builder-alert{margin:0;border-radius:3px;background:#eef2ff;color:#4f46e5;padding:12px 16px;font-size:13px}.builder-alert.error{background:#fdeaea;color:#c53030}.builder-toolbar{display:flex;align-items:flex-end;justify-content:space-between;gap:20px;flex-wrap:wrap}.builder-toolbar h2{margin:0;color:#34395e;font-size:20px}.builder-toolbar p{margin:5px 0 0;color:#6c757d;font-size:13px}.toolbar-actions,.add-field,.review-actions,.row-actions,.field-actions{display:flex;align-items:center;gap:8px;flex-wrap:wrap}.primary-action,.secondary-action,.danger-action,.back-action,.row-actions button,.field-actions button,.submission-detail header button,.pagination button{display:inline-flex;align-items:center;justify-content:center;gap:7px;border:0;border-radius:4px;padding:9px 14px;font-size:12px;font-weight:700}.primary-action{background:#6777ef;color:#fff;box-shadow:0 2px 6px #acb5f6}.secondary-action,.back-action{background:#fff;color:#34395e;border:1px solid #e4e6fc}.danger-action,.remove-field{background:#fc544b!important;color:#fff!important}.back-action{margin-bottom:10px;padding-left:0;border:0}.primary-action:disabled,.secondary-action:disabled,.danger-action:disabled,.field-actions button:disabled,.pagination button:disabled{opacity:.45;cursor:not-allowed}.toolbar-actions select,.add-field select,.filter-panel select{min-width:150px}.table-panel,.settings-panel,.filter-panel,.preview-panel,.submission-detail{background:#fff;border-radius:3px;box-shadow:0 4px 25px rgba(0,0,0,.1);min-width:0}.table-panel{overflow:auto}table{width:100%;border-collapse:collapse;min-width:760px}th,td{text-align:left;padding:14px 16px;border-bottom:1px solid #f4f6f9;font-size:12px;vertical-align:middle}th{color:#34395e;font-weight:800;background:#fbfbfd}td{color:#6c757d}td strong{display:block;color:#34395e;font-size:13px}td small{display:block;margin-top:3px}.row-actions{justify-content:flex-end}.row-actions button,.field-actions button,.submission-detail header button,.pagination button{width:34px;height:34px;padding:0;background:#f4f6f9;color:#6777ef}.status-badge{display:inline-block;border-radius:20px;padding:5px 9px;background:#eef2ff;color:#6777ef;font-weight:700}.status-badge.success{background:#e8f7f0;color:#47c363}.status-badge.warning{background:#fff4e6;color:#d97706}.status-badge.danger{background:#fdeaea;color:#fc544b}.empty-state{padding:30px;text-align:center;color:#98a6ad}.editor-layout{display:grid;grid-template-columns:minmax(0,1.65fr) minmax(290px,.75fr);gap:20px;align-items:start}.editor-column{display:grid;gap:20px;min-width:0}.settings-panel,.preview-panel,.submission-detail{padding:24px}.panel-heading{display:flex;align-items:center;justify-content:space-between;gap:16px;border-bottom:1px solid #f4f6f9;padding-bottom:14px;margin-bottom:18px}.panel-heading h3,.submission-detail h3,.answer-list h4{margin:0;color:#34395e;font-size:16px}.panel-heading p{margin:3px 0 0;color:#98a6ad;font-size:12px}.form-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:16px}.form-grid.compact{padding:18px}.form-grid label,.field-label,.filter-panel label,.review-panel label{display:grid;gap:7px;color:#34395e;font-size:12px;font-weight:700;min-width:0}.form-grid .wide{grid-column:1/-1}.form-grid input,.form-grid select,.form-grid textarea,.toolbar-actions select,.add-field select,.filter-panel select,.review-panel select,.review-panel textarea,.preview-field input,.preview-field select,.preview-field textarea{width:100%;box-sizing:border-box;border:1px solid #e4e6fc;border-radius:3px;background:#fdfdff;color:#495057;padding:10px 13px;font:inherit;outline:none}.form-grid input:focus,.form-grid select:focus,.form-grid textarea:focus,.filter-panel select:focus,.review-panel select:focus,.review-panel textarea:focus{border-color:#6777ef;box-shadow:0 2px 6px #acb5f6}.status-control{grid-column:1/-1;margin:0;padding:0;border:0}.status-control legend{margin-bottom:7px;color:#34395e;font-size:12px;font-weight:700}.status-control div{display:flex;flex-wrap:wrap;border:1px solid #e4e6fc;border-radius:4px;width:max-content;max-width:100%;overflow:hidden}.status-control button{border:0;border-right:1px solid #e4e6fc;background:#fff;color:#6c757d;padding:8px 12px;font-size:11px;font-weight:700}.status-control button:last-child{border-right:0}.status-control button.active{background:#6777ef;color:#fff}.fields-section{display:grid;gap:14px}.fields-heading{margin:0}.field-card{background:#fff;border:1px solid #e4e6fc;border-radius:4px;box-shadow:0 3px 15px rgba(0,0,0,.06);overflow:hidden}.field-card header{display:flex;align-items:center;gap:10px;padding:12px 14px;background:#fbfbfd;border-bottom:1px solid #f4f6f9}.field-card header strong{min-width:0;flex:1;color:#34395e;font-size:13px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.field-number{display:inline-flex;align-items:center;justify-content:center;width:26px;height:26px;border-radius:50%;background:#eef2ff;color:#6777ef;font-size:11px;font-weight:800}.required-toggle{display:flex!important;align-items:center;gap:8px!important}.required-toggle input{width:18px;height:18px;accent-color:#6777ef}.fields-empty{border:1px dashed #d9dcf2;background:#fff}.section-card{display:grid;gap:0;overflow:hidden;border:1px solid #dfe3fb;border-radius:4px;background:#fff;box-shadow:0 4px 18px rgba(0,0,0,.06)}.section-card-header{display:flex;align-items:center;gap:10px;padding:13px 15px;background:#f6f7ff;border-bottom:1px solid #e7e9ff}.section-card-header strong{min-width:0;flex:1;color:#34395e;font-size:14px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.section-number{display:inline-flex;align-items:center;justify-content:center;min-width:62px;height:27px;border-radius:20px;background:#6777ef;color:#fff;font-size:11px;font-weight:800}.section-settings{border-bottom:1px solid #f4f6f9}.section-field-toolbar{display:flex;align-items:center;justify-content:space-between;gap:12px;padding:12px 16px;background:#fbfbfd;border-bottom:1px solid #f4f6f9;color:#98a6ad;font-size:12px;font-weight:700}.section-field-card{margin:12px 16px}.section-card>.fields-empty{margin:12px 16px}.preview-step-list{display:flex;flex-wrap:wrap;gap:6px}.preview-step-list span{display:inline-flex;border-radius:20px;background:#eef2ff;color:#6777ef;padding:5px 9px;font-size:10px;font-weight:800}.preview-section{display:grid;gap:10px;border-top:1px solid #f4f6f9;padding-top:14px}.preview-section>small{color:#6777ef;font-size:10px;font-weight:800;text-transform:uppercase}.preview-section h4{margin:0;color:#34395e;font-size:15px}.preview-section>p{margin:0;color:#6c757d;font-size:12px;line-height:1.5}.preview-panel{position:sticky;top:92px;overflow:hidden;padding:0}.preview-banner{aspect-ratio:16/6;background:#f4f6f9}.preview-banner img{width:100%;height:100%;object-fit:cover}.preview-body{display:grid;gap:15px;padding:22px}.preview-kicker{text-transform:uppercase;color:#6777ef;font-size:10px;font-weight:800}.preview-body h3{margin:0;color:#34395e;font-size:20px}.preview-body>p{margin:0;color:#6c757d;font-size:12px;line-height:1.6}.fee-notice{background:#fff4e6;color:#b45309;padding:10px 12px;border-radius:3px;font-size:12px;font-weight:700}.preview-field{display:grid;gap:6px}.preview-field>label{color:#34395e;font-size:12px;font-weight:700}.preview-field>label span{color:#fc544b}.preview-field small,.preview-options small{color:#98a6ad;font-size:11px}.preview-options{display:grid;gap:7px}.preview-options label{display:flex;align-items:center;gap:7px;color:#6c757d;font-size:12px}.preview-options input{width:15px;height:15px}.filter-panel{display:grid;grid-template-columns:repeat(2,minmax(180px,300px));gap:16px;padding:18px}.submissions-layout{display:grid;gap:20px;align-items:start}.submissions-layout.detailed{grid-template-columns:minmax(0,1.45fr) minmax(320px,.75fr)}.clickable-row{cursor:pointer}.clickable-row:hover td{background:#f8f9ff}.pagination{display:flex;align-items:center;justify-content:center;gap:12px;padding:14px}.pagination span{font-size:12px;color:#6c757d}.submission-detail{position:sticky;top:92px;display:grid;gap:20px}.submission-detail header{display:flex;align-items:flex-start;justify-content:space-between;gap:12px}.submission-detail header span{color:#98a6ad;font-size:11px;text-transform:uppercase;font-weight:800}.submission-meta{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:12px;margin:0}.submission-meta div{min-width:0}.submission-meta dt{color:#98a6ad;font-size:10px;font-weight:800;text-transform:uppercase}.submission-meta dd{margin:3px 0 0;color:#34395e;font-size:12px;overflow-wrap:anywhere}.answer-list,.review-panel{display:grid;gap:12px;border-top:1px solid #f4f6f9;padding-top:18px}.answer-list>div{display:grid;gap:4px}.answer-list strong{color:#34395e;font-size:12px}.answer-list span{color:#6c757d;font-size:12px;white-space:pre-wrap;overflow-wrap:anywhere}.review-actions{justify-content:flex-end}.sr-only{position:absolute;width:1px;height:1px;padding:0;margin:-1px;overflow:hidden;clip:rect(0,0,0,0);white-space:nowrap;border:0}
:global(.dark .form-builder) .builder-toolbar h2,:global(.dark .form-builder) .panel-heading h3,:global(.dark .form-builder) .submission-detail h3,:global(.dark .form-builder) .answer-list h4,:global(.dark .form-builder) td strong,:global(.dark .form-builder) th,:global(.dark .form-builder) .field-card header strong,:global(.dark .form-builder) .preview-body h3,:global(.dark .form-builder) .preview-field>label,:global(.dark .form-builder) .submission-meta dd,:global(.dark .form-builder) .answer-list strong,:global(.dark .form-builder) label,:global(.dark .form-builder) legend{color:#f8fafc!important}:global(.dark .form-builder) .table-panel,:global(.dark .form-builder) .settings-panel,:global(.dark .form-builder) .filter-panel,:global(.dark .form-builder) .preview-panel,:global(.dark .form-builder) .submission-detail,:global(.dark .form-builder) .field-card,:global(.dark .form-builder) .fields-empty{background:#1f2937;border-color:#334155}:global(.dark .form-builder) th,:global(.dark .form-builder) .field-card header{background:#111827;border-color:#334155}:global(.dark .form-builder) td,:global(.dark .form-builder) .panel-heading,:global(.dark .form-builder) .answer-list,:global(.dark .form-builder) .review-panel{border-color:#334155}:global(.dark .form-builder) input,:global(.dark .form-builder) select,:global(.dark .form-builder) textarea,:global(.dark .form-builder) .secondary-action,:global(.dark .form-builder) .status-control button{background:#0f172a!important;color:#f8fafc!important;border-color:#475569!important}
@media(max-width:1180px){.editor-layout,.submissions-layout.detailed{grid-template-columns:1fr}.preview-panel,.submission-detail{position:static}.preview-panel{max-width:720px}}@media(max-width:700px){.form-grid,.filter-panel,.submission-meta{grid-template-columns:1fr}.form-grid .wide{grid-column:auto}.builder-toolbar{align-items:flex-start}.toolbar-actions,.add-field{width:100%}.toolbar-actions>*{flex:1 1 150px}.builder-tabs{overflow:auto}.builder-tabs button{white-space:nowrap}.settings-panel,.submission-detail{padding:18px}.form-grid.compact{padding:14px}.field-card header{align-items:flex-start;flex-wrap:wrap}.field-card header strong{flex-basis:calc(100% - 46px)}.field-actions{width:100%;justify-content:flex-end}}
</style>
