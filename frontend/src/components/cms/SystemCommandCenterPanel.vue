<template>
  <section class="ops-panel space-y-6">
    <div class="cms-panel-head border-b border-slate-200 dark:border-slate-800 pb-4 flex items-center justify-between flex-wrap gap-4">
      <div>
        <h2 class="text-lg font-bold text-slate-900 dark:text-white flex items-center gap-2">
          <i class="icofont-layers text-blue-600 dark:text-blue-400"></i> System Command Center
        </h2>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">Real-time Docker container telemetry, system resource utilization, and service controls.</p>
      </div>
      <div class="flex items-center gap-3">
        <button type="button" class="px-3 py-1.5 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-xs font-semibold rounded text-slate-700 dark:text-slate-200 transition-colors flex items-center gap-1" @click="loadResources">
          <i class="icofont-refresh"></i> Refresh
        </button>
        <label class="ops-live-toggle flex items-center gap-1.5 text-xs font-semibold text-slate-700 dark:text-slate-300">
          <input v-model="liveResources" type="checkbox" class="rounded border-slate-300 text-blue-600 focus:ring-blue-500" /> Live refresh (4s)
        </label>
      </div>
    </div>

    <!-- 4 Metrics Grid -->
    <div class="ops-grid grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <article class="cms-card metric bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm">
        <span class="text-xs font-semibold text-slate-500 dark:text-slate-400">CPU Utilization</span>
        <strong class="text-2xl font-bold text-slate-900 dark:text-white block my-1">{{ round(resources?.cpu?.overall_percent) }}%</strong>
        <div class="ops-bar h-1.5 bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
          <span class="block h-full bg-blue-600 rounded-full transition-all duration-300" :style="{ width: round(resources?.cpu?.overall_percent) + '%' }"></span>
        </div>
      </article>

      <article class="cms-card metric bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm">
        <span class="text-xs font-semibold text-slate-500 dark:text-slate-400">Memory Usage</span>
        <strong class="text-2xl font-bold text-slate-900 dark:text-white block my-1">{{ round(resources?.memory_detail?.used_percent) }}%</strong>
        <div class="ops-bar h-1.5 bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
          <span class="block h-full bg-emerald-500 rounded-full transition-all duration-300" :style="{ width: round(resources?.memory_detail?.used_percent) + '%' }"></span>
        </div>
        <small class="text-[11px] text-slate-400 mt-1 block">{{ formatBytes(resources?.memory_detail?.used_bytes) }} / {{ formatBytes(resources?.memory_detail?.total_bytes) }}</small>
      </article>

      <article class="cms-card metric bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm">
        <span class="text-xs font-semibold text-slate-500 dark:text-slate-400">Disk Space</span>
        <strong class="text-2xl font-bold text-slate-900 dark:text-white block my-1">{{ round(resources?.disk?.used_percent) }}%</strong>
        <div class="ops-bar h-1.5 bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
          <span class="block h-full bg-purple-500 rounded-full transition-all duration-300" :style="{ width: round(resources?.disk?.used_percent) + '%' }"></span>
        </div>
        <small class="text-[11px] text-slate-400 mt-1 block">{{ formatBytes(resources?.disk?.used_bytes) }} / {{ formatBytes(resources?.disk?.total_bytes) }}</small>
      </article>

      <article class="cms-card metric bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm">
        <span class="text-xs font-semibold text-slate-500 dark:text-slate-400">Network Bandwidth</span>
        <strong class="text-2xl font-bold text-slate-900 dark:text-white block my-1">{{ (resources?.network?.ingress_mbps || 0).toFixed(2) }} Mbps</strong>
        <small class="text-[11px] text-slate-400 mt-1 block">in · {{ (resources?.network?.egress_mbps || 0).toFixed(2) }} Mbps out</small>
      </article>
    </div>

    <!-- Services & Containers Table -->
    <article class="cms-panel bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-6 shadow-sm space-y-4">
      <div class="cms-panel-head border-b border-slate-100 dark:border-slate-800 pb-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <h3 class="text-base font-bold text-slate-800 dark:text-white">Services & Containers</h3>
          <span class="px-2.5 py-0.5 text-xs font-bold rounded-full bg-blue-50 text-blue-600 dark:bg-blue-950/60 dark:text-blue-300 border border-blue-200 dark:border-blue-800">
            {{ resources?.services?.length || 0 }} Services
          </span>
        </div>
      </div>
      <div class="overflow-x-auto">
        <table class="ops-table w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-slate-200 dark:border-slate-800 text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider">
              <th class="py-2.5 px-3">Service Name</th>
              <th class="py-2.5 px-3">Container ID / Name</th>
              <th class="py-2.5 px-3">Status</th>
              <th class="py-2.5 px-3">Health</th>
              <th class="py-2.5 px-3">Ports</th>
              <th class="py-2.5 px-3">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-800 text-xs text-slate-700 dark:text-slate-300">
            <tr v-for="svc in resources?.services || []" :key="svc.name" class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30 transition-colors">
              <td class="py-3 px-3">
                <div class="flex flex-col">
                  <span class="font-bold text-slate-900 dark:text-white text-sm">{{ svc.display_name }}</span>
                  <small class="text-xs text-slate-400 font-mono">{{ svc.name }}</small>
                </div>
              </td>
              <td class="py-3 px-3">
                <span v-if="svc.container" class="font-mono text-xs text-slate-600 dark:text-slate-400">{{ svc.container }}</span>
                <span v-else class="text-xs italic text-slate-400">Not created</span>
              </td>
              <td class="py-3 px-3">
                <span class="ops-badge px-2.5 py-1 rounded-full text-[11px] font-bold inline-block capitalize" :class="displayServiceStatus(svc)">{{ svc.status || 'idle' }}</span>
              </td>
              <td class="py-3 px-3">
                <span class="ops-health-badge px-2.5 py-1 rounded-full text-[11px] font-bold inline-block capitalize" :class="displayServiceHealth(svc)">{{ svc.health || 'unknown' }}</span>
              </td>
              <td class="py-3 px-3">
                <small v-if="svc.published_ports" class="font-mono text-xs text-slate-500 dark:text-slate-400">{{ formatPorts(svc.published_ports) }}</small>
                <small v-else class="text-xs text-slate-400">—</small>
              </td>
              <td class="py-3 px-3 ops-actions flex items-center gap-1.5 flex-wrap">
                <button
                  type="button"
                  class="ops-btn ops-btn-success"
                  :disabled="busyService === svc.name || svc.status === 'running'"
                  @click="serviceAction(svc.name, 'start')"
                >
                  <i class="icofont-play"></i> Start
                </button>
                <button
                  type="button"
                  class="ops-btn ops-btn-primary"
                  :disabled="busyService === svc.name"
                  @click="serviceAction(svc.name, 'restart')"
                >
                  <i class="icofont-refresh"></i> Restart
                </button>
                <button
                  type="button"
                  class="ops-btn ops-danger"
                  :disabled="busyService === svc.name || svc.status === 'stopped'"
                  @click="serviceAction(svc.name, 'stop')"
                >
                  <i class="icofont-stop"></i> Stop
                </button>
                <button
                  type="button"
                  class="ops-btn ops-btn-secondary"
                  :class="{ 'ops-btn-active': activeLogService === svc.name }"
                  @click="loadLogs(svc.name)"
                >
                  <i class="icofont-terminal"></i> Logs
                </button>
              </td>
            </tr>
            <tr v-if="!resources?.services?.length"><td colspan="6" class="cms-empty text-center py-6 text-slate-400 italic">No service telemetry available.</td></tr>
          </tbody>
        </table>
      </div>
      <p v-if="dockerUnavailable" class="ops-hint text-xs text-slate-500 dark:text-slate-400 flex items-center gap-1">
        <i class="icofont-warning-alt text-amber-500"></i> Docker socket is not mounted into this container. Control operations and log streaming are currently limited.
      </p>
    </article>

    <!-- Service Logs Terminal -->
    <article class="cms-panel bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-6 shadow-sm space-y-4">
      <div class="cms-panel-head border-b border-slate-100 dark:border-slate-800 pb-3 flex items-center justify-between flex-wrap gap-3">
        <div class="flex items-center gap-2">
          <h3 class="text-base font-bold text-slate-800 dark:text-white">Service Logs {{ activeLogService ? `— ${activeLogService}` : '' }}</h3>
          <span v-if="activeLogService" class="text-xs text-slate-400 font-mono">({{ logLines.length }} lines)</span>
        </div>
        <div class="flex items-center gap-3">
          <input
            v-model="logFilter"
            type="text"
            placeholder="Filter logs..."
            class="px-3 py-1.5 text-xs border border-slate-200 dark:border-slate-700 rounded-lg bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100 focus:outline-none focus:ring-1 focus:ring-blue-500"
          />
          <label class="ops-live-toggle flex items-center gap-1.5 text-xs font-semibold text-slate-700 dark:text-slate-300">
            <input v-model="liveLogs" type="checkbox" :disabled="!activeLogService" class="rounded border-slate-300 text-blue-600 focus:ring-blue-500" /> Live tail (4s)
          </label>
        </div>
      </div>
      <pre class="ops-terminal bg-slate-950 text-sky-400 rounded-xl p-4 font-mono text-xs leading-relaxed max-h-96 overflow-y-auto border border-slate-800 whitespace-pre-wrap word-break-break-word" ref="terminalRef">{{ filteredLogLines.join('
') || 'Select a service above and click "Logs" to inspect live container logs.' }}</pre>
    </article>

    <!-- System Events -->
    <article class="cms-panel bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-6 shadow-sm space-y-4">
      <div class="cms-panel-head border-b border-slate-100 dark:border-slate-800 pb-3">
        <h3 class="text-base font-bold text-slate-800 dark:text-white">Recent System Events</h3>
      </div>
      <ul class="ops-events divide-y divide-slate-100 dark:divide-slate-800">
        <li v-for="(evt, idx) in resources?.events || []" :key="idx" class="py-2.5 flex items-center gap-3 text-xs text-slate-700 dark:text-slate-300" :class="evt.severity">
          <span class="ops-event-dot w-2 h-2 rounded-full bg-blue-500 flex-shrink-0"></span>
          <span class="font-medium flex-grow">{{ evt.message }}</span>
          <small class="font-mono text-slate-400 ml-auto">{{ formatDateTime(evt.timestamp) }}</small>
        </li>
        <li v-if="!resources?.events?.length" class="cms-empty text-center py-4 text-slate-400 italic">No recent system events.</li>
      </ul>
    </article>
  </section>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import Swal from 'sweetalert2'
import { getSystemResources, getServiceLogs, runServiceAction } from '@/api/operator.js'

const emit = defineEmits(['message', 'error'])

const resources = ref(null)
const dockerUnavailable = ref(false)
const liveResources = ref(true)
const busyService = ref('')
const activeLogService = ref('')
const logLines = ref([])
const logFilter = ref('')
const liveLogs = ref(false)
const terminalRef = ref(null)
let resourcesTimer = null
let logsTimer = null

const filteredLogLines = computed(() => {
  if (!logFilter.value) return logLines.value
  const query = logFilter.value.toLowerCase()
  return logLines.value.filter(line => line.toLowerCase().includes(query))
})

function round(v) { return Math.round(Number(v) || 0) }
function formatBytes(bytes) {
  const n = Number(bytes) || 0
  if (n <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.min(units.length - 1, Math.floor(Math.log(n) / Math.log(1024)))
  return `${(n / Math.pow(1024, i)).toFixed(1)} ${units[i]}`
}
function formatDateTime(value) {
  return value ? new Date(value).toLocaleString() : ''
}
function formatPorts(published) {
  if (typeof published === 'string') return published
  if (Array.isArray(published)) {
    return published.map(p => typeof p === 'object' ? `${p.PublishedPort || p.TargetPort}->${p.TargetPort}` : String(p)).join(', ')
  }
  return String(published || '—')
}
function displayServiceStatus(service = {}) {
  const st = String(service.status || '').toLowerCase()
  if (st === 'running' || st === 'up') return 'running'
  if (st === 'stopped' || st === 'exited' || st === 'dead') return 'stopped'
  if (st === 'restarting') return 'restarting'
  return 'idle'
}
function displayServiceHealth(service = {}) {
  const health = String(service.health || '').toLowerCase()
  if (health === 'healthy') return 'healthy'
  if (health === 'unhealthy') return 'unhealthy'
  if (health === 'starting') return 'starting'
  return 'unknown'
}

async function loadResources() {
  try {
    resources.value = await getSystemResources()
    dockerUnavailable.value = resources.value?.docker_available === false
  } catch (err) {
    emit('error', err?.response?.data?.message || err?.message || 'Failed to load system resources')
  }
}

async function loadLogs(service) {
  activeLogService.value = service
  try {
    const data = await getServiceLogs(service, 150)
    logLines.value = data.lines || []
    await nextTick()
    if (terminalRef.value) terminalRef.value.scrollTop = terminalRef.value.scrollHeight
  } catch (err) {
    emit('error', err?.response?.data?.message || err?.message || `Failed to fetch logs for ${service}`)
  }
}

async function serviceAction(service, action) {
  const expected = `${action.toUpperCase()} ${service.toUpperCase()}`
  const result = await Swal.fire({
    title: `${action[0].toUpperCase()}${action.slice(1)} ${service}?`,
    text: `Type ${expected} to confirm operation.`,
    input: 'text',
    inputPlaceholder: expected,
    showCancelButton: true,
    confirmButtonText: 'Confirm Action',
    confirmButtonColor: action === 'stop' ? '#fc544b' : '#6777ef',
    cancelButtonColor: '#94a3b8',
    inputValidator: value => value !== expected ? `Type exactly: ${expected}` : undefined,
  })
  if (!result.isConfirmed) return
  busyService.value = service
  try {
    await runServiceAction(service, action, expected)
    emit('message', `Service ${service} ${action} command submitted successfully`)
    await loadResources()
  } catch (err) {
    emit('error', err?.response?.data?.message || err?.message || `Failed to ${action} ${service}`)
  } finally {
    busyService.value = ''
  }
}

watch(liveResources, value => {
  clearInterval(resourcesTimer)
  if (value) resourcesTimer = setInterval(loadResources, 4000)
})
watch(liveLogs, value => {
  clearInterval(logsTimer)
  if (value && activeLogService.value) logsTimer = setInterval(() => loadLogs(activeLogService.value), 4000)
})

onMounted(() => {
  loadResources()
  if (liveResources.value) resourcesTimer = setInterval(loadResources, 4000)
})
onUnmounted(() => {
  clearInterval(resourcesTimer)
  clearInterval(logsTimer)
})
</script>

<style scoped>
.ops-badge { border-radius: 30px; padding: 3px 10px; font-size: 11px; font-weight: 700; text-transform: capitalize; background: #f1f5f9; color: #64748b; }
.ops-badge.running { background: #dcfce7; color: #15803d; }
.ops-badge.stopped { background: #fee2e2; color: #b91c1c; }
.ops-badge.restarting { background: #fef3c7; color: #b45309; }

.ops-health-badge { border-radius: 30px; padding: 3px 10px; font-size: 11px; font-weight: 700; text-transform: capitalize; background: #f1f5f9; color: #64748b; }
.ops-health-badge.healthy { background: #dcfce7; color: #15803d; }
.ops-health-badge.unhealthy { background: #fee2e2; color: #b91c1c; }
.ops-health-badge.starting { background: #fef3c7; color: #b45309; }

.ops-btn { border: 0; border-radius: 6px; padding: 6px 12px; font-size: 11px; font-weight: 700; cursor: pointer; transition: all 150ms ease; display: inline-flex; align-items: center; gap: 4px; }
.ops-btn:disabled { opacity: .45; cursor: not-allowed; }
.ops-btn-primary { background: #3b82f6; color: #fff; }
.ops-btn-primary:hover:not(:disabled) { background: #2563eb; }
.ops-btn-success { background: #22c55e; color: #fff; }
.ops-btn-success:hover:not(:disabled) { background: #16a34a; }
.ops-btn-secondary { background: #e2e8f0; color: #334155; }
.ops-btn-secondary:hover:not(:disabled) { background: #cbd5e1; }
.ops-btn-active { background: #0f172a !important; color: #fff !important; }
.ops-danger { background: #ef4444; color: #fff; }
.ops-danger:hover:not(:disabled) { background: #dc2626; }

:global(.dark .ops-btn-secondary) { background: #334155; color: #f8fafc; }
:global(.dark .ops-btn-secondary:hover:not(:disabled)) { background: #475569; }
</style>
