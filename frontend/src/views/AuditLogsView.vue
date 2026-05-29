<template>
  <LayoutDefault title="Audit Logs">
    <div class="space-y-4">

      <!-- Stats bar -->
      <div class="flex flex-wrap items-center gap-3">
        <span class="flex items-center gap-1.5 bg-red-50 border border-red-200 text-red-700 rounded-lg px-3 py-2 text-xs font-medium">
          <span class="w-2 h-2 rounded-full bg-red-500"></span>
          {{ stats.critical }} Critical
        </span>
        <span class="flex items-center gap-1.5 bg-yellow-50 border border-yellow-200 text-yellow-700 rounded-lg px-3 py-2 text-xs font-medium">
          <span class="w-2 h-2 rounded-full bg-yellow-500"></span>
          {{ stats.warnings }} Warnings
        </span>
        <span class="flex items-center gap-1.5 bg-orange-50 border border-orange-200 text-orange-700 rounded-lg px-3 py-2 text-xs font-medium">
          <span class="w-2 h-2 rounded-full bg-orange-500"></span>
          {{ stats.vpn }} VPN/Proxy
        </span>
      </div>

      <!-- Error alert -->
      <div v-if="error" class="flex items-center gap-2 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">
        {{ error }}
      </div>

      <!-- Main card -->
      <div class="admin-card overflow-hidden">

        <!-- Card header: search + filters -->
        <div class="admin-card-header flex-wrap gap-2">
          <h2 class="text-sm font-semibold text-gray-800">Security Audit Logs</h2>
          <div class="flex flex-wrap items-center gap-2">
            <!-- Search -->
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="w-3.5 h-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
                </svg>
              </div>
              <input
                v-model="searchQuery"
                @input="debouncedSearch"
                type="text"
                placeholder="Search logs…"
                class="pl-8 pr-3 py-1.5 text-xs border border-gray-200 rounded-lg bg-gray-50 focus:outline-none focus:ring-1 focus:ring-primary-700 w-48"
              />
            </div>
            <select v-model="filterSeverity" @change="applyFilters" class="py-1.5 px-3 text-xs border border-gray-200 rounded-lg bg-gray-50 focus:outline-none focus:ring-1 focus:ring-primary-700">
              <option value="">All Severities</option>
              <option value="CRITICAL">Critical</option>
              <option value="WARNING">Warning</option>
              <option value="ERROR">Error</option>
              <option value="INFO">Info</option>
            </select>
            <select v-model="filterThreat" @change="applyFilters" class="py-1.5 px-3 text-xs border border-gray-200 rounded-lg bg-gray-50 focus:outline-none focus:ring-1 focus:ring-primary-700">
              <option value="">All Threats</option>
              <option value="high">High (≥70)</option>
              <option value="medium">Medium (≥40)</option>
              <option value="low">Low (&lt;40)</option>
            </select>
          </div>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="p-8 flex justify-center">
          <svg class="animate-spin w-7 h-7 text-primary-600" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
          </svg>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-gray-100">
                <th class="table-th w-24">Severity</th>
                <th class="table-th">Event</th>
                <th class="table-th">Status</th>
                <th class="table-th">User / IP</th>
                <th class="table-th">Location</th>
                <th class="table-th">Client</th>
                <th class="table-th w-20">Threat</th>
                <th class="table-th">Time</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-50">
              <tr v-if="filtered.length === 0">
                <td colspan="8" class="py-16 text-center">
                  <div class="flex flex-col items-center gap-2 text-gray-400">
                    <i class="icofont-history text-5xl text-gray-200"></i>
                    <p class="text-sm">No audit logs found</p>
                  </div>
                </td>
              </tr>
              <tr
                v-for="log in filtered"
                :key="log.id"
                @click="openDetail(log)"
                class="hover:bg-gray-50 transition-colors cursor-pointer"
                :class="rowHighlight(log)"
              >
                <td class="table-td">
                  <span :class="severityBadge(log.severity_level)" class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-semibold">
                    <span class="w-1.5 h-1.5 rounded-full" :class="severityDot(log.severity_level)"></span>
                    {{ log.severity_level || 'INFO' }}
                  </span>
                </td>
                <td class="table-td">
                  <div class="text-sm font-medium text-gray-800">{{ log.event_type || log.resource || '—' }}</div>
                  <div class="text-xs text-gray-400 font-mono mt-0.5 truncate max-w-48">{{ log.method }} {{ truncatePath(log.endpoint) }}</div>
                </td>
                <td class="table-td">
                  <span :class="statusBadge(log.event_status)" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium">
                    {{ log.event_status || '—' }}
                  </span>
                  <div class="text-xs text-gray-400 mt-0.5">{{ log.response_code }} · {{ log.response_time_ms }}ms</div>
                </td>
                <td class="table-td">
                  <div class="text-sm font-medium text-gray-700">{{ displayName(log) }}</div>
                  <div class="text-xs text-gray-400 font-mono mt-0.5">{{ log.ip_address || '—' }}</div>
                </td>
                <td class="table-td">
                  <div class="flex items-center gap-1.5">
                    <span class="text-sm text-gray-700">{{ log.geo_country || '—' }}</span>
                    <span v-if="log.vpn_detected" class="text-xs bg-red-100 text-red-700 font-semibold px-1.5 py-0.5 rounded">VPN</span>
                  </div>
                  <div class="text-xs text-gray-400 mt-0.5">{{ log.geo_city || '' }}</div>
                </td>
                <td class="table-td">
                  <div class="text-sm text-gray-700">{{ log.client_type || '—' }}</div>
                  <div class="text-xs text-gray-400 mt-0.5 truncate max-w-36">{{ log.browser }}{{ log.os_name ? ' / ' + log.os_name : '' }}</div>
                </td>
                <td class="table-td">
                  <div class="flex items-center gap-1">
                    <div class="flex-1 h-1.5 bg-gray-100 rounded-full overflow-hidden w-12">
                      <div class="h-full rounded-full transition-all" :class="threatBarColor(log.threat_score)" :style="{ width: log.threat_score + '%' }"></div>
                    </div>
                    <span class="text-xs font-mono" :class="threatTextColor(log.threat_score)">{{ log.threat_score }}</span>
                  </div>
                  <span v-if="log.anomaly_detected" class="text-xs text-orange-600 font-medium">Anomaly</span>
                </td>
                <td class="table-td text-gray-500 whitespace-nowrap">{{ formatDateTime(log.created_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="meta && meta.total > 0" class="admin-card-footer flex items-center justify-between">
          <p class="text-xs text-gray-500">
            Showing <span class="font-medium text-gray-700">{{ (meta.page - 1) * meta.per_page + 1 }}–{{ Math.min(meta.page * meta.per_page, meta.total) }}</span> of <span class="font-medium text-gray-700">{{ meta.total }}</span> entries
          </p>
          <div class="flex items-center gap-1.5">
            <button @click="prevPage" :disabled="meta.page <= 1" class="px-3 py-1 text-xs font-medium border border-gray-200 rounded-lg text-gray-600 hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors">← Prev</button>
            <span class="text-xs text-gray-500 px-2 font-medium">{{ meta.page }}</span>
            <button @click="nextPage" :disabled="meta.page * meta.per_page >= meta.total" class="px-3 py-1 text-xs font-medium border border-gray-200 rounded-lg text-gray-600 hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors">Next →</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Detail Drawer -->
    <Teleport to="body">
      <div v-if="detailLog" class="fixed inset-0 z-50 flex justify-end" @click.self="detailLog = null">
        <div class="fixed inset-0 bg-black/30 backdrop-blur-sm" @click="detailLog = null"></div>
        <div class="relative bg-white w-full max-w-xl h-full overflow-y-auto shadow-2xl z-10">
          <!-- Header -->
          <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
            <div>
              <h3 class="font-semibold text-gray-900">Log Entry Detail</h3>
              <p class="text-xs text-gray-500 font-mono mt-0.5">{{ detailLog.id }}</p>
            </div>
            <button @click="detailLog = null" class="p-2 rounded-lg hover:bg-gray-100 text-gray-400 hover:text-gray-700">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div class="p-6 space-y-6">
            <!-- Severity + Status badges -->
            <div class="flex items-center gap-3 flex-wrap">
              <span :class="severityBadge(detailLog.severity_level)" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm font-semibold">
                <span class="w-2 h-2 rounded-full" :class="severityDot(detailLog.severity_level)"></span>
                {{ detailLog.severity_level || 'INFO' }}
              </span>
              <span :class="statusBadge(detailLog.event_status)" class="inline-flex items-center px-3 py-1.5 rounded-full text-sm font-medium">
                {{ detailLog.event_status || '—' }}
              </span>
              <span v-if="detailLog.vpn_detected" class="bg-red-100 text-red-700 text-sm font-semibold px-3 py-1.5 rounded-full">VPN/Proxy Detected</span>
              <span v-if="detailLog.anomaly_detected" class="bg-orange-100 text-orange-700 text-sm font-semibold px-3 py-1.5 rounded-full">Anomaly Detected</span>
            </div>

            <!-- Threat score bar -->
            <div class="bg-gray-50 rounded-xl p-4">
              <div class="flex items-center justify-between mb-2">
                <span class="text-sm font-medium text-gray-700">Threat Score</span>
                <span class="text-lg font-bold" :class="threatTextColor(detailLog.threat_score)">{{ detailLog.threat_score }} / 100</span>
              </div>
              <div class="h-3 bg-gray-200 rounded-full overflow-hidden">
                <div class="h-full rounded-full transition-all" :class="threatBarColor(detailLog.threat_score)" :style="{ width: detailLog.threat_score + '%' }"></div>
              </div>
            </div>

            <!-- Event info -->
            <div>
              <h4 class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Event Information</h4>
              <dl class="grid grid-cols-2 gap-3">
                <div><dt class="text-xs text-gray-500">Event Type</dt><dd class="text-sm font-medium text-gray-800 mt-0.5">{{ detailLog.event_type || '—' }}</dd></div>
                <div><dt class="text-xs text-gray-500">HTTP Method</dt><dd class="text-sm font-mono font-medium text-gray-800 mt-0.5">{{ detailLog.method || '—' }}</dd></div>
                <div class="col-span-2"><dt class="text-xs text-gray-500">Endpoint</dt><dd class="text-sm font-mono text-gray-700 mt-0.5 break-all">{{ detailLog.endpoint || '—' }}</dd></div>
                <div><dt class="text-xs text-gray-500">Response Code</dt><dd class="text-sm font-medium mt-0.5" :class="detailLog.response_code >= 400 ? 'text-red-600' : 'text-green-600'">{{ detailLog.response_code || '—' }}</dd></div>
                <div><dt class="text-xs text-gray-500">Response Time</dt><dd class="text-sm font-medium text-gray-800 mt-0.5">{{ detailLog.response_time_ms }}ms</dd></div>
                <div><dt class="text-xs text-gray-500">Resource</dt><dd class="text-sm font-medium text-gray-800 mt-0.5">{{ detailLog.resource || '—' }}</dd></div>
                <div><dt class="text-xs text-gray-500">Timestamp</dt><dd class="text-sm font-medium text-gray-800 mt-0.5">{{ formatDateTime(detailLog.created_at) }}</dd></div>
              </dl>
            </div>

            <!-- User & Session -->
            <div>
              <h4 class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">User & Session</h4>
              <dl class="grid grid-cols-2 gap-3">
                <div class="col-span-2"><dt class="text-xs text-gray-500">User</dt><dd class="text-sm font-medium text-gray-800 mt-0.5">{{ displayName(detailLog) }}</dd></div>
                <div class="col-span-2"><dt class="text-xs text-gray-500">User ID</dt><dd class="text-sm font-mono text-gray-700 mt-0.5 break-all">{{ detailLog.user_id || 'Anonymous' }}</dd></div>
                <div class="col-span-2"><dt class="text-xs text-gray-500">Session ID (JTI)</dt><dd class="text-sm font-mono text-gray-700 mt-0.5 break-all">{{ detailLog.session_id || '—' }}</dd></div>
              </dl>
            </div>

            <!-- Network & Geo -->
            <div>
              <h4 class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Network & Geo-Location</h4>
              <dl class="grid grid-cols-2 gap-3">
                <div><dt class="text-xs text-gray-500">IP Address</dt><dd class="text-sm font-mono text-gray-700 mt-0.5">{{ detailLog.ip_address || '—' }}</dd></div>
                <div><dt class="text-xs text-gray-500">Forwarded For</dt><dd class="text-sm font-mono text-gray-500 mt-0.5 truncate">{{ detailLog.forwarded_ip || '—' }}</dd></div>
                <div><dt class="text-xs text-gray-500">Country</dt><dd class="text-sm font-medium text-gray-800 mt-0.5">{{ detailLog.geo_country || '—' }}</dd></div>
                <div><dt class="text-xs text-gray-500">City</dt><dd class="text-sm font-medium text-gray-800 mt-0.5">{{ detailLog.geo_city || '—' }}</dd></div>
                <div>
                  <dt class="text-xs text-gray-500">VPN / Proxy</dt>
                  <dd class="mt-0.5">
                    <span :class="detailLog.vpn_detected ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'" class="text-xs font-semibold px-2 py-0.5 rounded-full">
                      {{ detailLog.vpn_detected ? 'Detected' : 'None' }}
                    </span>
                  </dd>
                </div>
              </dl>
            </div>

            <!-- Device & Client -->
            <div>
              <h4 class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Device & Client</h4>
              <dl class="grid grid-cols-2 gap-3">
                <div><dt class="text-xs text-gray-500">Client Type</dt><dd class="text-sm font-medium text-gray-800 mt-0.5">{{ detailLog.client_type || '—' }}</dd></div>
                <div><dt class="text-xs text-gray-500">Browser</dt><dd class="text-sm font-medium text-gray-800 mt-0.5">{{ detailLog.browser || '—' }}</dd></div>
                <div><dt class="text-xs text-gray-500">OS</dt><dd class="text-sm font-medium text-gray-800 mt-0.5">{{ detailLog.os_name || '—' }}</dd></div>
                <div class="col-span-2"><dt class="text-xs text-gray-500">User Agent</dt><dd class="text-xs font-mono text-gray-500 mt-0.5 break-all">{{ detailLog.user_agent || '—' }}</dd></div>
              </dl>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </LayoutDefault>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'
import apiClient from '@/api/client.js'

const breadcrumbStore = useBreadcrumbStore()

const logs = ref([])
const loading = ref(true)
const error = ref('')
const meta = ref(null)
const currentPage = ref(1)
const searchQuery = ref('')
const filterSeverity = ref('')
const filterThreat = ref('')
const detailLog = ref(null)

let searchTimer = null

const stats = computed(() => ({
  critical: logs.value.filter(l => l.severity_level === 'CRITICAL').length,
  warnings: logs.value.filter(l => l.severity_level === 'WARNING').length,
  vpn: logs.value.filter(l => l.vpn_detected).length,
}))

const filtered = computed(() => {
  return logs.value.filter(l => {
    if (filterSeverity.value && l.severity_level !== filterSeverity.value) return false
    if (filterThreat.value === 'high' && l.threat_score < 70) return false
    if (filterThreat.value === 'medium' && l.threat_score < 40) return false
    if (filterThreat.value === 'low' && l.threat_score >= 40) return false
    return true
  })
})

function debouncedSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => { currentPage.value = 1; loadLogs() }, 400)
}

function applyFilters() { currentPage.value = 1; loadLogs() }

async function loadLogs() {
  loading.value = true
  error.value = ''
  try {
    const params = new URLSearchParams({ page: currentPage.value, per_page: 50 })
    if (searchQuery.value) params.set('search', searchQuery.value)
    const res = await apiClient.get(`/api/v1/admin/audit-logs?${params}`)
    logs.value = Array.isArray(res.data.data) ? res.data.data : []
    meta.value = res.data.meta || null
  } catch (err) {
    error.value = err.response?.data?.error?.message || 'Failed to load audit logs.'
  } finally {
    loading.value = false
  }
}

function openDetail(log) { detailLog.value = log }

function displayName(log) {
  const first = log.first_name || ''
  const last  = log.last_name  || ''
  const full  = [first, last].filter(Boolean).join(' ')
  if (full) return full
  return log.username || (log.user_id ? log.user_id.slice(0, 8) + '…' : 'Anonymous')
}

// ── Badge helpers ─────────────────────────────────────────────────

function severityBadge(s) {
  switch (s) {
    case 'CRITICAL': return 'bg-red-100 text-red-800'
    case 'WARNING':  return 'bg-yellow-100 text-yellow-800'
    case 'ERROR':    return 'bg-orange-100 text-orange-800'
    default:         return 'bg-blue-50 text-blue-700'
  }
}

function severityDot(s) {
  switch (s) {
    case 'CRITICAL': return 'bg-red-500'
    case 'WARNING':  return 'bg-yellow-500'
    case 'ERROR':    return 'bg-orange-500'
    default:         return 'bg-blue-400'
  }
}

function statusBadge(s) {
  switch (s) {
    case 'SUCCESS':     return 'bg-green-100 text-green-700'
    case 'AUTH_FAILED': return 'bg-red-100 text-red-700'
    case 'DENIED':      return 'bg-orange-100 text-orange-700'
    case 'FAILED':      return 'bg-yellow-100 text-yellow-700'
    case 'ERROR':       return 'bg-red-100 text-red-800'
    default:            return 'bg-gray-100 text-gray-600'
  }
}

function rowHighlight(log) {
  if (log.severity_level === 'CRITICAL') return 'bg-red-50/40'
  if (log.vpn_detected) return 'bg-orange-50/40'
  return ''
}

function threatBarColor(score) {
  if (score >= 70) return 'bg-red-500'
  if (score >= 40) return 'bg-yellow-500'
  if (score > 0)   return 'bg-blue-400'
  return 'bg-gray-300'
}

function threatTextColor(score) {
  if (score >= 70) return 'text-red-600'
  if (score >= 40) return 'text-yellow-600'
  return 'text-gray-500'
}

function truncatePath(path) {
  if (!path) return '—'
  return path.length > 40 ? path.slice(0, 40) + '…' : path
}

function formatDateTime(dateStr) {
  if (!dateStr) return 'N/A'
  try {
    return new Date(dateStr).toLocaleString('en-UG', {
      day: '2-digit', month: 'short', year: 'numeric',
      hour: '2-digit', minute: '2-digit', second: '2-digit'
    })
  } catch { return dateStr }
}

function prevPage() {
  if (currentPage.value > 1) { currentPage.value--; loadLogs() }
}
function nextPage() {
  if (meta.value && currentPage.value * meta.value.per_page < meta.value.total) { currentPage.value++; loadLogs() }
}

onMounted(() => {
  breadcrumbStore.set('Audit Logs', [{ label: 'Security' }, { label: 'Audit Logs' }])
  loadLogs()
})
</script>
