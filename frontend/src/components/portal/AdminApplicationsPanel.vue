<template>
  <section class="applications-panel">
    <header class="applications-head">
      <div><p>Application Management</p><h2>All Applications</h2><span>Review applications submitted through standard and custom forms.</span></div>
      <button type="button" class="refresh-button" :disabled="loading" title="Refresh applications" @click="load"><i class="icofont-refresh"></i></button>
    </header>

    <div class="application-filters">
      <label class="search-box"><i class="icofont-search-1"></i><input v-model="search" type="search" placeholder="Search reference, applicant or form" /></label>
      <select v-model="source"><option value="">All sources</option><option value="custom">Custom forms</option><option value="standard">Standard forms</option></select>
      <select v-model="status"><option value="">All submitted statuses</option><option v-for="item in statuses" :key="item" :value="item">{{ titleize(item) }}</option></select>
    </div>

    <div class="application-summary">
      <span><strong>{{ filtered.length }}</strong> shown</span>
      <span><strong>{{ customCount }}</strong> custom form</span>
      <span><strong>{{ standardCount }}</strong> standard form</span>
    </div>

    <div class="application-table-wrap">
      <table>
        <thead><tr><th>Reference</th><th>Applicant</th><th>Application</th><th>Source</th><th>Status</th><th>Updated</th><th><span class="sr-only">Actions</span></th></tr></thead>
        <tbody>
          <tr v-if="loading"><td colspan="7" class="table-empty">Loading applications...</td></tr>
          <tr v-else-if="!paged.length"><td colspan="7" class="table-empty">No submitted applications match these filters.</td></tr>
          <tr v-for="item in paged" v-else :key="`${item.source}-${item.id}`">
            <td><strong>{{ item.reference || 'Not assigned' }}</strong><small>{{ shortId(item.id) }}</small></td>
            <td><strong>{{ item.applicant_name || 'Portal user' }}</strong><small>{{ item.applicant_email || shortId(item.user_id) }}</small></td>
            <td><strong>{{ item.title }}</strong><small>{{ item.payment_status ? `Payment: ${titleize(item.payment_status)}` : 'No payment recorded' }}</small></td>
            <td><span class="source-badge" :class="item.source">{{ item.source === 'custom' ? 'Custom form' : 'Standard' }}</span></td>
            <td><span class="status-badge" :class="statusClass(item.status)">{{ titleize(item.status) }}</span></td>
            <td>{{ formatDate(item.updated_at) }}</td>
            <td><button type="button" class="view-button" title="View application" @click="openDetail(item)"><i class="icofont-eye-alt"></i></button></td>
          </tr>
        </tbody>
      </table>
    </div>
    <footer v-if="pageCount > 1" class="pagination">
      <button type="button" :disabled="page === 1" title="Previous page" @click="page--"><i class="icofont-rounded-left"></i></button>
      <span>Page {{ page }} of {{ pageCount }}</span>
      <button type="button" :disabled="page === pageCount" title="Next page" @click="page++"><i class="icofont-rounded-right"></i></button>
    </footer>
    <p v-if="error" class="panel-error">{{ error }}</p>

    <div v-if="selected" class="detail-backdrop" @click.self="selected = null">
      <aside class="application-detail" role="dialog" aria-modal="true" aria-label="Application details">
        <header><div><small>{{ selected.source === 'custom' ? 'Custom form submission' : 'Standard application' }}</small><h3>{{ selected.title }}</h3><p>{{ selected.reference || 'Reference pending' }}</p></div><button type="button" title="Close details" @click="selected = null"><i class="icofont-close"></i></button></header>
        <div v-if="detailLoading" class="detail-empty">Loading application details...</div>
        <template v-else>
          <dl class="detail-meta">
            <div><dt>Applicant</dt><dd>{{ selected.applicant_name || selected.applicant_email || selected.user_id }}</dd></div>
            <div><dt>Status</dt><dd><span class="status-badge" :class="statusClass(selected.status)">{{ titleize(selected.status) }}</span></dd></div>
            <div><dt>Payment</dt><dd>{{ titleize(selected.payment_status || 'not required') }}</dd></div>
            <div><dt>Submitted</dt><dd>{{ formatDate(selected.submitted_at || selected.created_at) }}</dd></div>
          </dl>
          <section class="answers">
            <h4>Application responses</h4>
            <div v-if="!answerRows.length" class="detail-empty">No response data is available.</div>
            <dl v-else><div v-for="row in answerRows" :key="row.key"><dt>{{ titleize(row.key) }}</dt><dd>{{ displayValue(row.value) }}</dd></div></dl>
          </section>
          <section v-if="selected.review_notes" class="review-note"><h4>Current review note</h4><p>{{ selected.review_notes }}</p></section>
          <section v-if="canReview" class="review-actions">
            <label>Review note<textarea v-model="reviewNotes" rows="3" placeholder="Add context for the applicant or review record"></textarea></label>
            <div>
              <button type="button" class="approve" :disabled="saving" @click="review('approve')"><i class="icofont-check-circled"></i> Approve</button>
              <button type="button" class="request" :disabled="saving" @click="review('request-info')"><i class="icofont-question-circle"></i> Request info</button>
              <button type="button" class="reject" :disabled="saving" @click="review('reject')"><i class="icofont-close-circled"></i> Reject</button>
            </div>
          </section>
        </template>
      </aside>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { adminGetApplication, adminListApplications, adminReviewApplication } from '@/api/applications.js'
import { adminGetSubmission, adminListSubmissions, adminReviewSubmission } from '@/api/forms.js'

const props = defineProps({ initialFilter: { type: String, default: '' } })
const emit = defineEmits(['message', 'error'])
const loading = ref(false)
const detailLoading = ref(false)
const saving = ref(false)
const error = ref('')
const applications = ref([])
const search = ref('')
const source = ref('')
const status = ref('')
const page = ref(1)
const perPage = 20
const selected = ref(null)
const reviewNotes = ref('')
const statuses = ['SUBMITTED', 'RESUBMITTED', 'UNDER_REVIEW', 'NEEDS_INFORMATION', 'PENDING_PAYMENT', 'APPROVED', 'REJECTED']

const filtered = computed(() => applications.value.filter(item => {
  if (!status.value && item.status === 'DRAFT') return false
  if (status.value && item.status !== status.value) return false
  if (source.value && item.source !== source.value) return false
  const needle = search.value.trim().toLowerCase()
  if (!needle) return true
  return [item.reference, item.title, item.applicant_name, item.applicant_email, item.user_id].some(value => String(value || '').toLowerCase().includes(needle))
}))
const pageCount = computed(() => Math.max(1, Math.ceil(filtered.value.length / perPage)))
const paged = computed(() => filtered.value.slice((page.value - 1) * perPage, page.value * perPage))
const customCount = computed(() => filtered.value.filter(item => item.source === 'custom').length)
const standardCount = computed(() => filtered.value.filter(item => item.source === 'standard').length)
const answerRows = computed(() => Object.entries(parseObject(selected.value?.answers ?? selected.value?.form_data)).map(([key, value]) => ({ key, value })))
const canReview = computed(() => selected.value && !['APPROVED', 'REJECTED', 'DRAFT'].includes(selected.value.status))

watch([search, source, status], () => { page.value = 1 })
watch(() => props.initialFilter, applyInitialFilter, { immediate: true })
onMounted(load)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const [legacy, custom] = await Promise.all([
      adminListApplications({ page: 1, per_page: 200 }),
      adminListSubmissions({ page: 1, per_page: 200 }),
    ])
    const legacyRows = listData(legacy).map(item => ({
      ...item, source: 'standard', reference: item.application_reference,
      title: titleize(item.application_type || item.form_type || 'Standard application'),
    }))
    const customRows = listData(custom).map(item => ({
      ...item, source: 'custom', reference: item.submission_reference,
      title: item.template_title || 'Custom form application',
    }))
    applications.value = [...customRows, ...legacyRows].sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at))
  } catch (err) {
    setError(err, 'Could not load applications.')
  } finally {
    loading.value = false
  }
}

async function openDetail(item) {
  selected.value = { ...item }
  reviewNotes.value = item.review_notes || ''
  detailLoading.value = true
  try {
    const res = item.source === 'custom' ? await adminGetSubmission(item.id) : await adminGetApplication(item.id)
    selected.value = { ...item, ...unwrap(res) }
    reviewNotes.value = selected.value.review_notes || ''
  } catch (err) {
    setError(err, 'Could not load application details.')
  } finally {
    detailLoading.value = false
  }
}

async function review(action) {
  if (action !== 'approve' && !reviewNotes.value.trim()) {
    error.value = 'Add a review note before requesting information or rejecting an application.'
    return
  }
  saving.value = true
  try {
    if (selected.value.source === 'custom') {
      const statusMap = { approve: 'APPROVED', reject: 'REJECTED', 'request-info': 'NEEDS_INFORMATION' }
      await adminReviewSubmission(selected.value.id, statusMap[action], reviewNotes.value.trim())
    } else {
      await adminReviewApplication(selected.value.id, action, reviewNotes.value.trim())
    }
    emit('message', 'Application review saved.')
    selected.value = null
    await load()
  } catch (err) {
    setError(err, 'Could not save the review.')
  } finally {
    saving.value = false
  }
}

function applyInitialFilter(value) {
  if (value === 'review') status.value = 'SUBMITTED'
  else if (value === 'attention') status.value = 'NEEDS_INFORMATION'
}
function unwrap(value) { return value?.data?.data ?? value?.data ?? value ?? {} }
function listData(value) { const data = unwrap(value); return Array.isArray(data) ? data : [] }
function parseObject(value) {
  if (!value) return {}
  if (typeof value === 'object') return value
  try { return JSON.parse(value) } catch { return {} }
}
function displayValue(value) {
  if (Array.isArray(value)) return value.join(', ') || '-'
  if (value && typeof value === 'object') return JSON.stringify(value)
  return String(value ?? '-')
}
function titleize(value) { return String(value || '').toLowerCase().replaceAll('_', ' ').replaceAll('-', ' ').replace(/\b\w/g, char => char.toUpperCase()) }
function shortId(value) { const text = String(value || ''); return text.length > 12 ? `${text.slice(0, 8)}...` : text }
function formatDate(value) { return value ? new Intl.DateTimeFormat('en-UG', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) : '-' }
function statusClass(value) {
  if (value === 'APPROVED') return 'approved'
  if (value === 'REJECTED') return 'rejected'
  if (value === 'NEEDS_INFORMATION' || value === 'PENDING_PAYMENT') return 'attention'
  return 'pending'
}
function setError(err, fallback) {
  error.value = err.response?.data?.error?.message || fallback
  emit('error', err)
}
</script>

<style scoped>
.applications-panel{display:grid;gap:16px}.applications-head{display:flex;justify-content:space-between;gap:16px;align-items:center}.applications-head p{margin:0;color:#6777ef;font-size:11px;font-weight:800;text-transform:uppercase}.applications-head h2{margin:3px 0;color:#303345;font-size:21px}.applications-head span{color:#747b90;font-size:13px}.refresh-button,.view-button,.application-detail header button,.pagination button{display:grid;place-items:center;width:38px;height:38px;border:1px solid #dfe2ea;border-radius:5px;background:#fff;color:#51566b}.application-filters{display:grid;grid-template-columns:minmax(240px,1fr) 180px 210px;gap:10px}.application-filters select,.search-box{height:42px;border:1px solid #dfe2ea;border-radius:5px;background:#fff}.application-filters select{padding:0 10px;color:#50566b}.search-box{display:flex;align-items:center;gap:8px;padding:0 12px}.search-box input{min-width:0;width:100%;border:0;outline:0;background:transparent}.application-summary{display:flex;gap:22px;color:#73798c;font-size:12px}.application-summary strong{color:#303345}.application-table-wrap{overflow:auto;border:1px solid #e7e9f2;border-radius:6px;background:#fff}.application-table-wrap table{width:100%;min-width:980px;border-collapse:collapse}.application-table-wrap th,.application-table-wrap td{padding:13px 15px;border-bottom:1px solid #edf0f5;text-align:left;font-size:12px}.application-table-wrap th{background:#f8f9fc;color:#62697e;font-weight:800;text-transform:uppercase}.application-table-wrap td{color:#62697e}.application-table-wrap td strong,.application-table-wrap td small{display:block}.application-table-wrap td strong{color:#303345}.application-table-wrap td small{margin-top:3px;color:#9297a8}.source-badge,.status-badge{display:inline-flex;padding:4px 7px;border-radius:4px;font-size:10px;font-weight:800;text-transform:uppercase}.source-badge.custom{background:#e8edff;color:#5266d8}.source-badge.standard{background:#e2f6f8;color:#197f8d}.status-badge.pending{background:#e8edff;color:#5266d8}.status-badge.approved{background:#e5f8ee;color:#218b55}.status-badge.rejected{background:#fee9e8;color:#cf433d}.status-badge.attention{background:#fff1df;color:#b9650d}.table-empty{padding:45px!important;text-align:center!important;color:#81879a!important}.pagination{display:flex;align-items:center;justify-content:flex-end;gap:10px;color:#6e7589;font-size:12px}.pagination button:disabled{opacity:.45}.panel-error{margin:0;padding:11px;border:1px solid #fecaca;border-radius:5px;background:#fef2f2;color:#b91c1c}.detail-backdrop{position:fixed;inset:0;z-index:1080;display:flex;justify-content:flex-end;background:rgba(17,24,39,.48)}.application-detail{width:min(590px,100%);height:100%;overflow:auto;padding:24px;background:#fff;box-shadow:-10px 0 35px rgba(0,0,0,.16)}.application-detail>header{display:flex;justify-content:space-between;gap:15px;padding-bottom:18px;border-bottom:1px solid #e8eaf0}.application-detail header small{color:#6777ef;font-weight:800;text-transform:uppercase}.application-detail header h3{margin:4px 0;color:#303345;font-size:20px}.application-detail header p{margin:0;color:#767d90}.detail-meta{display:grid;grid-template-columns:1fr 1fr;gap:12px;margin:18px 0}.detail-meta div{padding:12px;background:#f7f8fb;border-radius:5px}.detail-meta dt,.answers dt{color:#81879a;font-size:11px;font-weight:800;text-transform:uppercase}.detail-meta dd,.answers dd{margin:4px 0 0;color:#303345}.answers h4,.review-note h4{color:#303345;font-size:15px}.answers dl{display:grid;gap:0}.answers dl div{padding:12px 0;border-bottom:1px solid #edf0f5}.answers dd{white-space:pre-wrap;word-break:break-word}.detail-empty{padding:36px;text-align:center;color:#80879a}.review-note{margin-top:18px;padding:14px;background:#f7f8fb;border-left:3px solid #6777ef}.review-note p{margin:5px 0 0;color:#5f667a}.review-actions{display:grid;gap:12px;margin-top:20px;padding-top:18px;border-top:1px solid #e8eaf0}.review-actions label{display:grid;gap:5px;color:#555c70;font-size:12px;font-weight:800}.review-actions textarea{width:100%;padding:10px;border:1px solid #dfe2ea;border-radius:5px}.review-actions>div{display:flex;flex-wrap:wrap;gap:8px}.review-actions button{padding:9px 12px;border:0;border-radius:5px;color:#fff;font-size:12px;font-weight:800}.review-actions .approve{background:#218b55}.review-actions .request{background:#c87313}.review-actions .reject{background:#cf433d}.sr-only{position:absolute;width:1px;height:1px;overflow:hidden;clip:rect(0,0,0,0)}:global(.dark) .applications-head h2,:global(.dark) .application-table-wrap td strong,:global(.dark) .application-summary strong,:global(.dark) .application-detail header h3,:global(.dark) .detail-meta dd,:global(.dark) .answers dd,:global(.dark) .answers h4,:global(.dark) .review-note h4{color:#f8fafc}:global(.dark) .application-table-wrap,:global(.dark) .application-detail,:global(.dark) .application-filters select,:global(.dark) .search-box,:global(.dark) .refresh-button,:global(.dark) .view-button{background:#1f2937;border-color:#334155}:global(.dark) .application-table-wrap th,:global(.dark) .detail-meta div,:global(.dark) .review-note{background:#111827}:global(.dark) .application-table-wrap td,:global(.dark) .application-table-wrap th{border-color:#334155}:global(.dark) .application-filters input,:global(.dark) .application-filters select{color:#e5e7eb}@media(max-width:760px){.applications-head{align-items:flex-start}.application-filters{grid-template-columns:1fr}.detail-meta{grid-template-columns:1fr}.application-summary{flex-wrap:wrap;gap:8px 18px}.application-detail{padding:18px}}
</style>
