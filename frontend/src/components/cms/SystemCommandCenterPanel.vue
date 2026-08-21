<template>
  <section class="ops-panel">
    <div class="cms-panel-head">
      <div>
        <h2>System Command Center</h2>
        <p class="text-xs text-slate-500 mt-1">Real-time Docker container telemetry, system resource utilization, and service controls.</p>
      </div>
      <div class="flex items-center gap-3">
        <button type="button" class="px-3 py-1.5 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-xs font-semibold rounded text-slate-700 dark:text-slate-200 transition-colors" @click="loadResources">
          <i class="icofont-refresh mr-1"></i> Refresh
        </button>
        <label class="ops-live-toggle"><input v-model="liveResources" type="checkbox" /> Live refresh (4s)</label>
      </div>
    </div>

    <div class="ops-grid">
      <article class="cms-card metric">
        <span>CPU Utilization</span>
        <strong>{{ round(resources?.cpu?.overall_percent) }}%</strong>
        <div class="ops-bar"><span :style="{ width: round(resources?.cpu?.overall_percent) + '%' }"></span></div>
      </article>
      <article class="cms-card metric">
        <span>Memory Usage</span>
        <strong>{{ round(resources?.memory_detail?.used_percent) }}%</strong>
        <div class="ops-bar"><span :style="{ width: round(resources?.memory_detail?.used_percent) + '%' }"></span></div>
        <small>{{ formatBytes(resources?.memory_detail?.used_bytes) }} / {{ formatBytes(resources?.memory_detail?.total_bytes) }}</small>
      </article>
      <article class="cms-card metric">
        <span>Disk Space</span>
        <strong>{{ round(resources?.disk?.used_percent) }}%</strong>
        <div class="ops-bar"><span :style="{ width: round(resources?.disk?.used_percent) + '%' }"></span></div>
        <small>{{ formatBytes(resources?.disk?.used_bytes) }} / {{ formatBytes(resources?.disk?.total_bytes) }}</small>
      </article>
      <article class="cms-card metric">
        <span>Network Bandwidth</span>
        <strong>{{ (resources?.network?.ingress_mbps || 0).toFixed(2) }} Mbps</strong>
        <small>in · {{ (resources?.network?.egress_mbps || 0).toFixed(2) }} Mbps out</small>
      </article>
    </div>

    <article class="cms-panel">
      <div class="cms-panel-head">
        <div class="flex items-center gap-2">
          <h3>Services & Containers</h3>
          <span class="px-2 py-0.5 text-xs font-semibold rounded-full bg-blue-50 text-blue-600 dark:bg-blue-900/40 dark:text-blue-300">
            {{ resources?.services?.length || 0 }} Services
          </span>
        </div>
      </div>
      <div class="overflow-x-auto">
        <table class="ops-table">
          <thead>
            <tr>
              <th>Service Name</th>
              <th>Container ID / Name</th>
              <th>Status</th>
              <th>Health</th>
              <th>Ports</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="svc in resources?.services || []" :key="svc.name">
              <td>
                <div class="flex flex-col">
                  <span class="font-bold text-slate-800 dark:text-white">{{ svc.display_name }}</span>
                  <small class="text-xs text-slate-400 font-mono">{{ svc.name }}</small>
                </div>
              </td>
              <td>
                <span v-if="svc.container" class="font-mono text-xs text-slate-600 dark:text-slate-300">{{ svc.container }}</span>
                <span v-else class="text-xs italic text-slate-400">Not created</span>
              </td>
              <td>
                <span class="ops-badge" :class="displayServiceStatus(svc)">{{ svc.status || 'idle' }}</span>
              </td>
              <td>
                <span class="ops-health-badge" :class="displayServiceHealth(svc)">{{ svc.health || 'unknown' }}</span>
              </td>
              <td>
                <small v-if="svc.published_ports" class="font-mono text-xs text-slate-500 dark:text-slate-400">{{ formatPorts(svc.published_ports) }}</small>
                <small v-else class="text-xs text-slate-400">—</small>
              </td>
              <td class="ops-actions">
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
            <tr v-if="!resources?.services?.length"><td colspan="6" class="cms-empty">No service telemetry available.</td></tr>
          </tbody>
        </table>
      </div>
      <p v-if="dockerUnavailable" class="ops-hint">
        <i class="icofont-warning-alt text-amber-500 mr-1"></i> Docker socket is not mounted into this container. Control operations and log streaming are currently limited.
      </p>
    </article>

    <article class="cms-panel">
      <div class="cms-panel-head">
        <div class="flex items-center gap-2">
          <h3>Service Logs {{ activeLogService ? `— ${activeLogService}` : '' }}</h3>
          <span v-if="activeLogService" class="text-xs text-slate-400 font-mono">({{ logLines.length }} lines)</span>
        </div>
        <div class="flex items-center gap-3">
          <input
            v-model="logFilter"
            type="text"
            placeholder="Filter logs..."
            class="px-3 py-1 text-xs border border-slate-200 dark:border-slate-700 rounded bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100"
          />
          <label class="ops-live-toggle"><input v-model="liveLogs" type="checkbox" :disabled="!activeLogService" /> Live tail (4s)</label>
        </div>
      </div>
      <pre class="ops-terminal" ref="terminalRef">{{ filteredLogLines.join('\n') || 'Select a service above and click "Logs" to inspect live container logs.' }}</pre>
    </article>

    <article class="cms-panel">
      <div class="cms-panel-head"><h3>Recent System Events</h3></div>
      <ul class="ops-events">
        <li v-for="(evt, idx) in resources?.events || []" :key="idx" :class="evt.severity">
          <span class="ops-event-dot"></span>
          <span class="font-medium">{{ evt.message }}</span>
          <small class="font-mono">{{ formatDateTime(evt.timestamp) }}</small>
        </li>
        <li v-if="!resources?.events?.length" class="cms-empty">No recent system events.</li>
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
.ops-panel { display: grid; gap: 20px; min-width: 0; }
.cms-panel-head { display: flex; justify-content: space-between; align-items: center; gap: 16px; flex-wrap: wrap; border-bottom: 1px solid #f9f9f9; padding-bottom: 15px; }
.cms-panel-head h2, .cms-panel-head h3 { margin: 0; font-size: 16px; font-weight: 700; color: #34395e; }
.cms-panel { background: white; border-radius: 6px; box-shadow: 0 4px 25px rgba(0,0,0,.06); border: 1px solid #e2e8f0; padding: 22px; display: grid; gap: 16px; min-width: 0; }
.cms-card.metric { background: white; box-shadow: 0 4px 25px rgba(0,0,0,.06); border: 1px solid #e2e8f0; border-radius: 6px; padding: 20px; display: grid; gap: 6px; }
.cms-card.metric span { color: #64748b; font-size: 12px; font-weight: 600; }
.cms-card.metric strong { font-size: 24px; color: #1e293b; }
.cms-card.metric small { color: #94a3b8; font-size: 11px; }
.cms-empty { color: #94a3b8; font-style: italic; text-align: center; padding: 16px; }

.ops-live-toggle { display: flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 600; color: #34395e; }
.ops-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr)); gap: 16px; }
.ops-bar { height: 6px; border-radius: 30px; background: #f1f5f9; overflow: hidden; }
.ops-bar span { display: block; height: 100%; background: #6777ef; border-radius: 30px; }

.ops-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.ops-table th { text-align: left; color: #64748b; font-size: 11px; text-transform: uppercase; padding: 10px 12px; border-bottom: 1px solid #e2e8f0; }
.ops-table td { padding: 12px; border-bottom: 1px solid #f1f5f9; color: #334155; vertical-align: middle; }
.ops-badge { border-radius: 30px; padding: 3px 10px; font-size: 11px; font-weight: 700; text-transform: capitalize; background: #f1f5f9; color: #64748b; }
.ops-badge.running { background: #dcfce7; color: #15803d; }
.ops-badge.stopped { background: #fee2e2; color: #b91c1c; }
.ops-badge.restarting { background: #fef3c7; color: #b45309; }

.ops-health-badge { border-radius: 30px; padding: 3px 10px; font-size: 11px; font-weight: 700; text-transform: capitalize; background: #f1f5f9; color: #64748b; }
.ops-health-badge.healthy { background: #dcfce7; color: #15803d; }
.ops-health-badge.unhealthy { background: #fee2e2; color: #b91c1c; }
.ops-health-badge.starting { background: #fef3c7; color: #b45309; }

.ops-actions { display: flex; gap: 6px; flex-wrap: wrap; }
.ops-btn { border: 0; border-radius: 6px; padding: 6px 12px; font-size: 11px; font-weight: 700; cursor: pointer; transition: all 150ms ease; display: inline-flex; align-items: center; gap: 4px; }
.ops-btn:disabled { opacity: .45; cursor: not-allowed; }
.ops-btn-primary { background: #6777ef; color: #fff; }
.ops-btn-primary:hover:not(:disabled) { background: #5565e3; }
.ops-btn-success { background: #22c55e; color: #fff; }
.ops-btn-success:hover:not(:disabled) { background: #16a34a; }
.ops-btn-secondary { background: #e2e8f0; color: #334155; }
.ops-btn-secondary:hover:not(:disabled) { background: #cbd5e1; }
.ops-btn-active { background: #0f172a !important; color: #fff !important; }
.ops-danger { background: #ef4444; color: #fff; }
.ops-danger:hover:not(:disabled) { background: #dc2626; }

.ops-hint { color: #64748b; font-size: 12px; margin: 0; padding-top: 8px; }

.ops-terminal { background: #0f172a; color: #38bdf8; border-radius: 6px; padding: 16px; font-family: 'Fira Code', 'Courier New', monospace; font-size: 12px; line-height: 1.6; max-height: 24rem; overflow: auto; white-space: pre-wrap; word-break: break-word; border: 1px solid #1e293b; }

.ops-events { list-style: none; margin: 0; padding: 0; display: grid; gap: 8px; }
.ops-events li { display: flex; align-items: center; gap: 10px; font-size: 13px; color: #334155; border-bottom: 1px stroke #f1f5f9; padding: 4px 0; }
.ops-events li small { margin-left: auto; color: #94a3b8; font-size: 11px; }
.ops-event-dot { width: 8px; height: 8px; border-radius: 50%; background: #3b82f6; flex-shrink: 0; }
.ops-events li.warning .ops-event-dot { background: #f59e0b; }
.ops-events li.critical .ops-event-dot { background: #ef4444; }

:global(.dark .cms-panel-head h2), :global(.dark .cms-panel-head h3) { color: #f8fafc; }
:global(.dark .cms-panel), :global(.dark .cms-card.metric) { background: #1f2937; border-color: #334155; }
:global(.dark .cms-card.metric strong), :global(.dark .ops-table td), :global(.dark .ops-events li), :global(.dark .ops-live-toggle) { color: #f8fafc; }
:global(.dark .ops-table th) { color: #94a3b8; border-color: #334155; }
:global(.dark .ops-table td) { border-color: #334155; }
:global(.dark .ops-btn-secondary) { background: #334155; color: #f8fafc; }
:global(.dark .ops-btn-secondary:hover:not(:disabled)) { background: #475569; }
</style>
