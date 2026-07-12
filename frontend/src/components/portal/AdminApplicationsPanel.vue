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
            <td class="action-cell">
              <button type="button" class="view-button" title="View application" @click="openDetail(item)"><i class="icofont-eye-alt"></i></button>
              <button type="button" class="view-button" title="Download application form" @click="downloadRow(item)"><i class="icofont-download"></i></button>
            </td>
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
  </section>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { adminListApplications } from '@/api/applications.js'
import { adminListSubmissions } from '@/api/forms.js'
import { downloadApplicationForm } from '@/utils/applicationDownload.js'

const props = defineProps({ initialFilter: { type: String, default: '' } })
const emit = defineEmits(['message', 'error'])
const router = useRouter()
const loading = ref(false)
const error = ref('')
const applications = ref([])
const search = ref('')
const source = ref('')
const status = ref('')
const page = ref(1)
const perPage = 20
const statuses = ['SUBMITTED', 'RESUBMITTED', 'UNDER_REVIEW', 'NEEDS_INFORMATION', 'COMPLETE', 'PENDING_PAYMENT', 'APPROVED', 'REJECTED']

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
  router.push({ name: 'AdminApplicationDetail', params: { id: item.id }, query: { source: item.source } })
}

function downloadRow(item) {
  downloadApplicationForm(item)
}

function applyInitialFilter(value) {
  if (value === 'review') status.value = 'SUBMITTED'
  else if (value === 'attention') status.value = 'NEEDS_INFORMATION'
}
function unwrap(value) { return value?.data?.data ?? value?.data ?? value ?? {} }
function listData(value) { const data = unwrap(value); return Array.isArray(data) ? data : [] }
function titleize(value) { return String(value || '').toLowerCase().replaceAll('_', ' ').replaceAll('-', ' ').replace(/\b\w/g, char => char.toUpperCase()) }
function shortId(value) { const text = String(value || ''); return text.length > 12 ? `${text.slice(0, 8)}...` : text }
function formatDate(value) { return value ? new Intl.DateTimeFormat('en-UG', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) : '-' }
function statusClass(value) {
  if (value === 'APPROVED' || value === 'COMPLETE') return 'approved'
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
.applications-panel{display:grid;gap:16px}.applications-head{display:flex;justify-content:space-between;gap:16px;align-items:center}.applications-head p{margin:0;color:#6777ef;font-size:11px;font-weight:800;text-transform:uppercase}.applications-head h2{margin:3px 0;color:#303345;font-size:21px}.applications-head span{color:#747b90;font-size:13px}.refresh-button,.view-button,.application-detail header button,.pagination button{display:grid;place-items:center;width:38px;height:38px;border:1px solid #dfe2ea;border-radius:5px;background:#fff;color:#51566b}.action-cell{display:flex;gap:6px}.application-filters{display:grid;grid-template-columns:minmax(240px,1fr) 180px 210px;gap:10px}.application-filters select,.search-box{height:42px;border:1px solid #dfe2ea;border-radius:5px;background:#fff}.application-filters select{padding:0 10px;color:#50566b}.search-box{display:flex;align-items:center;gap:8px;padding:0 12px}.search-box input{min-width:0;width:100%;border:0;outline:0;background:transparent}.application-summary{display:flex;gap:22px;color:#73798c;font-size:12px}.application-summary strong{color:#303345}.application-table-wrap{overflow:auto;border:1px solid #e7e9f2;border-radius:6px;background:#fff}.application-table-wrap table{width:100%;min-width:980px;border-collapse:collapse}.application-table-wrap th,.application-table-wrap td{padding:13px 15px;border-bottom:1px solid #edf0f5;text-align:left;font-size:12px}.application-table-wrap th{background:#f8f9fc;color:#62697e;font-weight:800;text-transform:uppercase}.application-table-wrap td{color:#62697e}.application-table-wrap td strong,.application-table-wrap td small{display:block}.application-table-wrap td strong{color:#303345}.application-table-wrap td small{margin-top:3px;color:#9297a8}.source-badge,.status-badge{display:inline-flex;padding:4px 7px;border-radius:4px;font-size:10px;font-weight:800;text-transform:uppercase}.source-badge.custom{background:#e8edff;color:#5266d8}.source-badge.standard{background:#e2f6f8;color:#197f8d}.status-badge.pending{background:#e8edff;color:#5266d8}.status-badge.approved{background:#e5f8ee;color:#218b55}.status-badge.rejected{background:#fee9e8;color:#cf433d}.status-badge.attention{background:#fff1df;color:#b9650d}.table-empty{padding:45px!important;text-align:center!important;color:#81879a!important}.pagination{display:flex;align-items:center;justify-content:flex-end;gap:10px;color:#6e7589;font-size:12px}.pagination button:disabled{opacity:.45}.panel-error{margin:0;padding:11px;border:1px solid #fecaca;border-radius:5px;background:#fef2f2;color:#b91c1c}.sr-only{position:absolute;width:1px;height:1px;overflow:hidden;clip:rect(0,0,0,0)}:global(.dark) .applications-head h2,:global(.dark) .application-table-wrap td strong,:global(.dark) .application-summary strong{color:#f8fafc}:global(.dark) .application-table-wrap,:global(.dark) .application-filters select,:global(.dark) .search-box,:global(.dark) .refresh-button,:global(.dark) .view-button{background:#1f2937;border-color:#334155}:global(.dark) .application-table-wrap th{background:#111827}:global(.dark) .application-table-wrap td,:global(.dark) .application-table-wrap th{border-color:#334155}:global(.dark) .application-filters input,:global(.dark) .application-filters select{color:#e5e7eb}@media(max-width:760px){.applications-head{align-items:flex-start}.application-filters{grid-template-columns:1fr}.application-summary{flex-wrap:wrap;gap:8px 18px}}
</style>
