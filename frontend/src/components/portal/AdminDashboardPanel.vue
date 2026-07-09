<template>
  <section class="admin-dashboard" aria-label="Portal operations dashboard">
    <div class="kpi-grid">
      <article v-for="card in kpiCards" :key="card.label" class="kpi-card">
        <span class="kpi-icon" :class="card.tone"><i :class="card.icon"></i></span>
        <div><small>{{ card.label }}</small><strong>{{ formatNumber(card.value) }}</strong><p>{{ card.note }}</p></div>
      </article>
    </div>

    <div class="dashboard-grid">
      <article class="dashboard-panel status-panel">
        <header><div><p>Application Pipeline</p><h3>Applications by status</h3></div><span>{{ formatNumber(stats.total_applications) }} total</span></header>
        <div v-if="loading" class="panel-empty">Loading operational data...</div>
        <div v-else-if="!statusRows.length" class="panel-empty">No applications have been received yet.</div>
        <div v-else class="status-list">
          <div v-for="row in statusRows" :key="row.status" class="status-row">
            <div><span>{{ titleize(row.status) }}</span><strong>{{ formatNumber(row.count) }}</strong></div>
            <div class="status-track"><span :class="statusTone(row.status)" :style="{ width: `${row.percent}%` }"></span></div>
          </div>
        </div>
      </article>

      <article class="dashboard-panel attention-panel">
        <header><div><p>Work Queue</p><h3>Items needing action</h3></div></header>
        <button type="button" @click="$emit('open-applications', 'review')">
          <span class="attention-icon amber"><i class="icofont-clock-time"></i></span>
          <span><strong>{{ formatNumber(stats.pending_review) }}</strong><small>Pending review</small></span>
          <i class="icofont-rounded-right"></i>
        </button>
        <button type="button" @click="$emit('open-applications', 'attention')">
          <span class="attention-icon red"><i class="icofont-warning-alt"></i></span>
          <span><strong>{{ formatNumber(stats.needs_attention) }}</strong><small>Need attention</small></span>
          <i class="icofont-rounded-right"></i>
        </button>
        <button type="button" @click="$emit('open-forms')">
          <span class="attention-icon green"><i class="icofont-ui-edit"></i></span>
          <span><strong>{{ formatNumber(stats.open_forms) }}</strong><small>Open application forms</small></span>
          <i class="icofont-rounded-right"></i>
        </button>
      </article>
    </div>

    <div class="dashboard-grid lower">
      <article class="dashboard-panel source-panel">
        <header><div><p>Intake</p><h3>Application sources</h3></div></header>
        <div class="source-row">
          <span><i class="icofont-file-document"></i> Standard applications</span>
          <strong>{{ formatNumber(stats.legacy_applications) }}</strong>
        </div>
        <div class="source-row">
          <span><i class="icofont-ui-edit"></i> Custom form submissions</span>
          <strong>{{ formatNumber(stats.dynamic_applications) }}</strong>
        </div>
      </article>
      <article class="dashboard-panel summary-panel">
        <header><div><p>Community</p><h3>Portal profile summary</h3></div></header>
        <div><span>Ordinary users</span><strong>{{ formatNumber(stats.ordinary_users) }}</strong></div>
        <div><span>Organisation profiles</span><strong>{{ formatNumber(stats.total_profiles) }}</strong></div>
        <div><span>Registered athletes</span><strong>{{ formatNumber(stats.total_athletes) }}</strong></div>
      </article>
    </div>
    <p v-if="error" class="dashboard-error">{{ error }}</p>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { getAdminDashboard } from '@/api/applications.js'

defineEmits(['open-applications', 'open-forms'])

const loading = ref(true)
const error = ref('')
const stats = ref({
  total_users: 0, ordinary_users: 0, total_athletes: 0, total_profiles: 0,
  open_forms: 0, total_applications: 0, legacy_applications: 0,
  dynamic_applications: 0, pending_review: 0, needs_attention: 0,
  applications_by_status: {},
})

const kpiCards = computed(() => [
  { label: 'Applications', value: stats.value.total_applications, note: `${formatNumber(stats.value.pending_review)} awaiting review`, icon: 'icofont-file-document', tone: 'blue' },
  { label: 'Athletes', value: stats.value.total_athletes, note: 'Registered athlete records', icon: 'icofont-runner', tone: 'green' },
  { label: 'Profiles', value: stats.value.total_profiles, note: 'Active organisation profiles', icon: 'icofont-id-card', tone: 'orange' },
  { label: 'Portal Users', value: stats.value.total_users, note: `${formatNumber(stats.value.ordinary_users)} ordinary users`, icon: 'icofont-users-alt-5', tone: 'cyan' },
])

const statusRows = computed(() => {
  const entries = Object.entries(stats.value.applications_by_status || {})
    .filter(([, count]) => Number(count) > 0)
    .sort((a, b) => Number(b[1]) - Number(a[1]))
  const max = Math.max(1, ...entries.map(([, count]) => Number(count)))
  return entries.map(([status, count]) => ({ status, count, percent: Math.max(5, Math.round((Number(count) / max) * 100)) }))
})

onMounted(load)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const res = await getAdminDashboard()
    stats.value = { ...stats.value, ...(res.data?.data || {}) }
  } catch (err) {
    error.value = err.response?.data?.error?.message || 'Could not load dashboard KPIs.'
  } finally {
    loading.value = false
  }
}

function formatNumber(value) {
  return new Intl.NumberFormat('en-UG').format(Number(value || 0))
}
function titleize(value) {
  return String(value || '').toLowerCase().replaceAll('_', ' ').replace(/\b\w/g, char => char.toUpperCase())
}
function statusTone(status) {
  if (status === 'APPROVED') return 'green'
  if (status === 'REJECTED') return 'red'
  if (status === 'NEEDS_INFORMATION' || status === 'PENDING_PAYMENT') return 'amber'
  return 'blue'
}
</script>

<style scoped>
.admin-dashboard{display:grid;gap:18px}.kpi-grid{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:16px}.kpi-card{display:flex;align-items:center;gap:14px;min-height:130px;padding:18px;background:#fff;border:1px solid #e7e9f2;border-radius:6px;box-shadow:0 4px 18px rgba(31,45,61,.06)}.kpi-icon{display:grid;place-items:center;flex:0 0 54px;width:54px;height:54px;border-radius:6px;font-size:26px}.kpi-icon.blue{background:#e8edff;color:#5266d8}.kpi-icon.green{background:#e5f8ee;color:#218b55}.kpi-icon.orange{background:#fff1df;color:#c87313}.kpi-icon.cyan{background:#e2f6f8;color:#197f8d}.kpi-card small{display:block;color:#70758a;font-size:12px;font-weight:700;text-transform:uppercase}.kpi-card strong{display:block;margin:3px 0;color:#25283a;font-size:28px;line-height:1.15}.kpi-card p{margin:0;color:#8b8fa3;font-size:12px}.dashboard-grid{display:grid;grid-template-columns:minmax(0,2fr) minmax(280px,1fr);gap:16px}.dashboard-grid.lower{grid-template-columns:1fr 1fr}.dashboard-panel{padding:20px;background:#fff;border:1px solid #e7e9f2;border-radius:6px;box-shadow:0 4px 18px rgba(31,45,61,.05)}.dashboard-panel header{display:flex;align-items:center;justify-content:space-between;gap:12px;margin-bottom:18px}.dashboard-panel header p{margin:0 0 3px;color:#788197;font-size:11px;font-weight:800;text-transform:uppercase}.dashboard-panel h3{margin:0;color:#303345;font-size:16px}.dashboard-panel header>span{color:#65708a;font-size:12px;font-weight:700}.status-list{display:grid;gap:15px}.status-row>div:first-child{display:flex;justify-content:space-between;margin-bottom:6px;color:#646b80;font-size:12px}.status-row strong{color:#2e3245}.status-track{height:7px;overflow:hidden;background:#eff1f6;border-radius:3px}.status-track span{display:block;height:100%;border-radius:3px}.status-track .blue{background:#6777ef}.status-track .green{background:#47c363}.status-track .red{background:#fc544b}.status-track .amber{background:#ffa426}.attention-panel button{display:grid;grid-template-columns:42px 1fr auto;align-items:center;gap:11px;width:100%;padding:11px 0;border:0;border-bottom:1px solid #eef0f5;background:transparent;text-align:left}.attention-panel button:last-child{border:0}.attention-panel button>span:nth-child(2) strong,.attention-panel button>span:nth-child(2) small{display:block}.attention-panel button>span:nth-child(2) strong{color:#303345;font-size:18px}.attention-panel button>span:nth-child(2) small{color:#747b90}.attention-icon{display:grid;place-items:center;width:38px;height:38px;border-radius:5px}.attention-icon.amber{background:#fff1df;color:#c87313}.attention-icon.red{background:#fee9e8;color:#cf433d}.attention-icon.green{background:#e5f8ee;color:#218b55}.source-row,.summary-panel>div{display:flex;align-items:center;justify-content:space-between;gap:12px;padding:13px 0;border-bottom:1px solid #eef0f5;color:#646b80}.source-row:last-child,.summary-panel>div:last-child{border:0}.source-row i{margin-right:8px;color:#6777ef}.source-row strong,.summary-panel strong{color:#303345;font-size:17px}.panel-empty{padding:35px;text-align:center;color:#7c8498}.dashboard-error{margin:0;padding:12px;border:1px solid #fecaca;border-radius:5px;background:#fef2f2;color:#b91c1c}:global(.dark) .kpi-card,:global(.dark) .dashboard-panel{background:#1f2937;border-color:#334155}:global(.dark) .kpi-card strong,:global(.dark) .dashboard-panel h3,:global(.dark) .status-row strong,:global(.dark) .attention-panel button>span:nth-child(2) strong,:global(.dark) .source-row strong,:global(.dark) .summary-panel strong{color:#f8fafc}:global(.dark) .kpi-card p,:global(.dark) .kpi-card small,:global(.dark) .status-row>div:first-child,:global(.dark) .source-row,:global(.dark) .summary-panel>div{color:#cbd5e1}:global(.dark) .attention-panel button,:global(.dark) .source-row,:global(.dark) .summary-panel>div{border-color:#334155}@media(max-width:1100px){.kpi-grid{grid-template-columns:repeat(2,minmax(0,1fr))}}@media(max-width:760px){.kpi-grid,.dashboard-grid,.dashboard-grid.lower{grid-template-columns:1fr}.kpi-card{min-height:108px}}
</style>
