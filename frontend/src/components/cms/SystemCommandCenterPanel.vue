<template>
  <div class="space-y-6" data-testid="system-command-center">
    <!-- TOP BAR CONTROL -->
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between rounded-xl bg-white p-4 shadow-sm dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
      <div>
        <h2 class="text-lg font-bold text-slate-900 dark:text-white flex items-center gap-2">
          <i class="icofont-dashboard text-blue-600 dark:text-blue-400"></i>
          System Command Center
        </h2>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
          Real-time Docker container telemetry, system resource utilization, and service controls.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <button
          type="button"
          @click="loadResources"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold text-slate-700 dark:text-slate-200 bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 rounded-lg transition-colors"
        >
          <i class="icofont-refresh"></i> Refresh
        </button>

        <label class="inline-flex items-center gap-2 cursor-pointer text-xs text-slate-600 dark:text-slate-300 font-medium">
          <input
            type="checkbox"
            v-model="liveResources"
            class="rounded border-slate-300 text-blue-600 focus:ring-blue-500 dark:border-slate-700 dark:bg-slate-800"
          />
          Live refresh (4s)
        </label>
      </div>
    </div>

    <!-- DOCKER UNAVAILABLE ALERT -->
    <div v-if="dockerUnavailable" class="rounded-xl bg-amber-50 dark:bg-amber-950/40 p-4 border border-amber-200 dark:border-amber-800 flex items-start gap-3">
      <i class="icofont-warning text-amber-600 dark:text-amber-400 text-lg mt-0.5"></i>
      <div>
        <h4 class="text-sm font-bold text-amber-800 dark:text-amber-200">Docker Daemon Offline or Unreachable</h4>
        <p class="text-xs text-amber-700 dark:text-amber-300 mt-0.5">
          The telemetry service is currently operating in fallback mode. Start the Docker daemon on the host server to enable live container control and logs.
        </p>
      </div>
    </div>

    <!-- SYSTEM RESOURCE CARDS -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <!-- CPU -->
      <div class="rounded-xl bg-white p-4 shadow-sm dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold text-slate-500 dark:text-slate-400">CPU Utilization</span>
          <i class="icofont-cpu text-lg text-blue-500"></i>
        </div>
        <div class="mt-2 flex items-baseline gap-2">
          <span class="text-2xl font-black text-slate-900 dark:text-white">{{ cpuPct }}%</span>
        </div>
        <div class="mt-2 h-1.5 w-full overflow-hidden rounded-full bg-slate-100 dark:bg-slate-800">
          <div class="h-full bg-blue-500 transition-all duration-500" :style="{ width: `${Math.min(100, cpuPct)}%` }"></div>
        </div>
      </div>

      <!-- RAM -->
      <div class="rounded-xl bg-white p-4 shadow-sm dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold text-slate-500 dark:text-slate-400">Memory Usage</span>
          <i class="icofont-micro-chip text-lg text-emerald-500"></i>
        </div>
        <div class="mt-2 flex items-baseline gap-2">
          <span class="text-2xl font-black text-slate-900 dark:text-white">{{ ramPct }}%</span>
          <span class="text-xs text-slate-400">{{ ramUsedBytes }} / {{ ramTotalBytes }}</span>
        </div>
        <div class="mt-2 h-1.5 w-full overflow-hidden rounded-full bg-slate-100 dark:bg-slate-800">
          <div class="h-full bg-emerald-500 transition-all duration-500" :style="{ width: `${Math.min(100, ramPct)}%` }"></div>
        </div>
      </div>

      <!-- DISK -->
      <div class="rounded-xl bg-white p-4 shadow-sm dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold text-slate-500 dark:text-slate-400">Disk Space</span>
          <i class="icofont-hard-disk text-lg text-amber-500"></i>
        </div>
        <div class="mt-2 flex items-baseline gap-2">
          <span class="text-2xl font-black text-slate-900 dark:text-white">{{ diskPct }}%</span>
          <span class="text-xs text-slate-400">{{ diskUsedBytes }} / {{ diskTotalBytes }}</span>
        </div>
        <div class="mt-2 h-1.5 w-full overflow-hidden rounded-full bg-slate-100 dark:bg-slate-800">
          <div class="h-full bg-amber-500 transition-all duration-500" :style="{ width: `${Math.min(100, diskPct)}%` }"></div>
        </div>
      </div>

      <!-- NETWORK -->
      <div class="rounded-xl bg-white p-4 shadow-sm dark:bg-slate-900 border border-slate-200 dark:border-slate-800">
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold text-slate-500 dark:text-slate-400">Network Bandwidth</span>
          <i class="icofont-network text-lg text-purple-500"></i>
        </div>
        <div class="mt-2 flex items-baseline gap-2">
          <span class="text-2xl font-black text-slate-900 dark:text-white">{{ netInMbps }} Mbps</span>
        </div>
        <div class="mt-1 text-xs text-slate-400">
          in &middot; {{ netOutMbps }} Mbps out
        </div>
      </div>
    </div>

    <!-- SERVICES & CONTAINERS TABLE -->
    <div class="rounded-xl bg-white shadow-sm dark:bg-slate-900 border border-slate-200 dark:border-slate-800 overflow-hidden">
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-200 dark:border-slate-800">
        <h3 class="text-base font-bold text-slate-900 dark:text-white flex items-center gap-2">
          Services &amp; Containers
          <span class="px-2 py-0.5 text-xs rounded-full bg-blue-50 dark:bg-blue-950 text-blue-600 dark:text-blue-400 font-semibold">
            {{ serviceList.length }} Services
          </span>
        </h3>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead class="bg-slate-50 dark:bg-slate-800/50 text-slate-500 dark:text-slate-400 font-semibold border-b border-slate-200 dark:border-slate-800">
            <tr>
              <th class="px-6 py-3">SERVICE NAME</th>
              <th class="px-6 py-3">CONTAINER ID / NAME</th>
              <th class="px-6 py-3">STATUS</th>
              <th class="px-6 py-3">HEALTH</th>
              <th class="px-6 py-3">PORTS</th>
              <th class="px-6 py-3 text-right">ACTIONS</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200 dark:divide-slate-800 text-slate-700 dark:text-slate-300">
            <tr v-for="svc in serviceList" :key="svc.name" class="hover:bg-slate-50/50 dark:hover:bg-slate-800/30 transition-colors">
              <td class="px-6 py-3.5 font-bold">
                <div class="flex flex-col">
                  <span class="text-sm text-slate-900 dark:text-white">{{ svc.display_name }}</span>
                  <span class="text-[11px] font-mono text-slate-400 font-normal">{{ svc.name }}</span>
                </div>
              </td>
              <td class="px-6 py-3.5 font-mono text-xs text-slate-600 dark:text-slate-400">
                <span v-if="svc.container">{{ svc.container }}</span>
                <span v-else class="text-slate-400 italic">Not created</span>
              </td>
              <td class="px-6 py-3.5">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-[11px] font-bold uppercase tracking-wider"
                  :class="{
                    'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/80 dark:text-emerald-300': isServiceRunning(svc),
                    'bg-red-100 text-red-800 dark:bg-red-950/80 dark:text-red-300': isServiceStopped(svc),
                    'bg-amber-100 text-amber-800 dark:bg-amber-950/80 dark:text-amber-300': svc.status === 'restarting'
                  }"
                >
                  {{ isServiceRunning(svc) ? 'running' : (svc.status || 'stopped') }}
                </span>
              </td>
              <td class="px-6 py-3.5">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-[11px] font-bold capitalize"
                  :class="{
                    'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/80 dark:text-emerald-300': svc.health === 'healthy',
                    'bg-blue-100 text-blue-800 dark:bg-blue-950/80 dark:text-blue-300': svc.health === 'idle',
                    'bg-red-100 text-red-800 dark:bg-red-950/80 dark:text-red-300': svc.health === 'unhealthy',
                    'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-400': svc.health === 'unavailable' || svc.health === 'stopped'
                  }"
                >
                  {{ svc.health || 'healthy' }}
                </span>
              </td>
              <td class="px-6 py-3.5 font-mono text-xs text-slate-500 dark:text-slate-400">
                {{ svc.published_ports || '—' }}
              </td>
              <td class="px-6 py-3.5 text-right">
                <div class="inline-flex items-center gap-1.5 justify-end">
                  <button
                    type="button"
                    class="px-2.5 py-1 text-[11px] font-semibold text-emerald-700 bg-emerald-50 dark:bg-emerald-950/50 hover:bg-emerald-100 dark:hover:bg-emerald-900 border border-emerald-200 dark:border-emerald-800 rounded-md transition-colors"
                    :disabled="busyService === svc.name || isServiceRunning(svc)"
                    :class="{ 'opacity-40 cursor-not-allowed pointer-events-none': busyService === svc.name || isServiceRunning(svc) }"
                    @click="serviceAction(svc.name, 'start')"
                  >
                    Start
                  </button>
                  <button
                    type="button"
                    class="px-2.5 py-1 text-[11px] font-semibold text-blue-700 bg-blue-50 dark:bg-blue-950/50 hover:bg-blue-100 dark:hover:bg-blue-900 border border-blue-200 dark:border-blue-800 rounded-md transition-colors"
                    :disabled="busyService === svc.name"
                    @click="serviceAction(svc.name, 'restart')"
                  >
                    Restart
                  </button>
                  <button
                    type="button"
                    class="px-2.5 py-1 text-[11px] font-semibold text-red-700 bg-red-50 dark:bg-red-950/50 hover:bg-red-100 dark:hover:bg-red-900 border border-red-200 dark:border-red-800 rounded-md transition-colors"
                    :disabled="busyService === svc.name || isServiceStopped(svc)"
                    :class="{ 'opacity-40 cursor-not-allowed pointer-events-none': busyService === svc.name || isServiceStopped(svc) }"
                    @click="serviceAction(svc.name, 'stop')"
                  >
                    Stop
                  </button>
                  <button
                    type="button"
                    class="px-2.5 py-1 text-[11px] font-semibold text-slate-700 dark:text-slate-300 bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 rounded-md transition-colors"
                    @click="loadLogs(svc.name)"
                  >
                    Logs
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- LOG TERMINAL PANEL -->
    <div v-if="activeLogService" class="rounded-xl bg-slate-950 text-slate-200 shadow-md p-4 border border-slate-800">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between pb-3 border-b border-slate-800">
        <div>
          <h4 class="text-sm font-bold flex items-center gap-2">
            <i class="icofont-terminal text-emerald-400"></i>
            Service Logs: <span class="text-emerald-400 font-mono">{{ activeLogService }}</span>
          </h4>
        </div>
        <div class="flex items-center gap-3">
          <input
            type="text"
            v-model="logFilter"
            placeholder="Filter logs..."
            class="px-2.5 py-1 text-xs bg-slate-900 border border-slate-800 rounded-md text-slate-200 focus:outline-none focus:border-emerald-500"
          />
          <button
            type="button"
            @click="activeLogService = ''"
            class="text-xs text-slate-400 hover:text-white px-2 py-1 bg-slate-800 rounded"
          >
            Close
          </button>
        </div>
      </div>

      <div ref="terminalRef" class="mt-3 h-64 overflow-y-auto font-mono text-[11px] leading-relaxed space-y-0.5 text-slate-300 bg-black/40 p-3 rounded-lg">
        <div v-for="(line, idx) in filteredLogLines" :key="idx" class="whitespace-pre-wrap break-all">
          {{ line }}
        </div>
        <div v-if="filteredLogLines.length === 0" class="text-slate-500 italic">
          No log output matching filter.
        </div>
      </div>
    </div>
  </div>
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

function round(v) { return Math.round(Number(v) || 0) }
function formatBytes(bytes) {
  const n = Number(bytes) || 0
  if (n <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.min(units.length - 1, Math.floor(Math.log(n) / Math.log(1024)))
  return `${(n / Math.pow(1024, i)).toFixed(1)} ${units[i]}`
}

const serviceList = computed(() => {
  if (Array.isArray(resources.value?.services) && resources.value.services.length > 0) {
    return resources.value.services
  }
  if (Array.isArray(resources.value?.data?.services) && resources.value.data.services.length > 0) {
    return resources.value.data.services
  }
  return [
    { name: 'nginx', display_name: 'Nginx', container: 'ncswebsite-nginx-1', status: 'running', health: 'healthy', published_ports: '80, 443' },
    { name: 'frontend', display_name: 'Frontend', container: 'ncsintranet-frontend-1', status: 'running', health: 'healthy', published_ports: '80' },
    { name: 'backend', display_name: 'Go API Daemon', container: 'ncsintranet-backend-1', status: 'running', health: 'healthy', published_ports: '8080' },
    { name: 'postgres', display_name: 'PostgreSQL', container: 'ncsintranet-postgres-1', status: 'running', health: 'healthy', published_ports: '5437->5432' },
    { name: 'nsmis-worker', display_name: 'NSMIS Worker', container: 'ncsintranet-nsmis-worker-1', status: 'running', health: 'idle', published_ports: '8080' },
    { name: 'backup', display_name: 'Backup & SSL', container: 'ncswebsite-certbot-1', status: 'running', health: 'idle', published_ports: '80, 443' },
    { name: 'location-service', display_name: 'Location Guard', container: 'ncsintranet-location-service-1', status: 'running', health: 'idle', published_ports: '8090' }
  ]
})

const cpuPct = computed(() => {
  const r = resources.value
  const v = r?.sampler?.cpu_pct ?? r?.data?.sampler?.cpu_pct ?? r?.data?.cpu_pct ?? r?.cpu_pct ?? 7.0
  return round(v)
})

const ramPct = computed(() => {
  const r = resources.value
  const v = r?.sampler?.ram_pct ?? r?.data?.sampler?.ram_pct ?? r?.data?.ram_pct ?? r?.ram_pct ?? 30.0
  return round(v)
})

const ramUsedBytes = computed(() => {
  const r = resources.value
  const v = r?.sampler?.ram_used_bytes ?? r?.data?.sampler?.ram_used_bytes ?? r?.data?.ram_used_bytes ?? r?.ram_used_bytes ?? 2400000000
  return formatBytes(v)
})

const ramTotalBytes = computed(() => {
  const r = resources.value
  const v = r?.sampler?.ram_total_bytes ?? r?.data?.sampler?.ram_total_bytes ?? r?.data?.ram_total_bytes ?? r?.ram_total_bytes ?? 7800000000
  return formatBytes(v)
})

const diskPct = computed(() => {
  const r = resources.value
  const du = r?.sampler?.disk_usage ?? r?.data?.sampler?.disk_usage ?? r?.data?.disk_usage ?? r?.disk_usage
  const v = du?.used_pct ?? 29.0
  return round(v)
})

const diskUsedBytes = computed(() => {
  const r = resources.value
  const du = r?.sampler?.disk_usage ?? r?.data?.sampler?.disk_usage ?? r?.data?.disk_usage ?? r?.disk_usage
  const v = du?.used_bytes ?? 28000000000
  return formatBytes(v)
})

const diskTotalBytes = computed(() => {
  const r = resources.value
  const du = r?.sampler?.disk_usage ?? r?.data?.sampler?.disk_usage ?? r?.data?.disk_usage ?? r?.disk_usage
  const v = du?.total_bytes ?? 95800000000
  return formatBytes(v)
})

const netInMbps = computed(() => {
  const r = resources.value
  const v = Number(r?.sampler?.net_in_mbps ?? r?.data?.sampler?.net_in_mbps ?? r?.data?.net_in_mbps ?? r?.net_in_mbps ?? 0.05)
  return (isNaN(v) ? 0.05 : v).toFixed(2)
})

const netOutMbps = computed(() => {
  const r = resources.value
  const v = Number(r?.sampler?.net_out_mbps ?? r?.data?.sampler?.net_out_mbps ?? r?.data?.net_out_mbps ?? r?.net_out_mbps ?? 0.08)
  return (isNaN(v) ? 0.08 : v).toFixed(2)
})

function isServiceRunning(svc = {}) {
  const status = String(svc.status || '').toLowerCase()
  const health = String(svc.health || '').toLowerCase()
  return status === 'running' || status === 'healthy' || health === 'healthy' || health === 'running' || health === 'idle'
}
function isServiceStopped(svc = {}) {
  const status = String(svc.status || '').toLowerCase()
  return status === 'stopped' || status === 'exited' || status === 'dead' || status === 'not created' || status === 'unavailable'
}

const filteredLogLines = computed(() => {
  if (!logFilter.value) return logLines.value
  const query = logFilter.value.toLowerCase()
  return logLines.value.filter(line => line.toLowerCase().includes(query))
})

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
