<template>
  <LayoutDefault title="Infrastructure Command Center">
    <div class="space-y-6 pb-12 transition-colors duration-200">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-server text-xs"></i> Infrastructure Telemetry
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Systems Maintenance Command Center</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Live host metrics, scoped maintenance switches, service runtime operations, and real-time event logs.
          </p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <button @click="refresh" class="btn btn-sm btn-light flex items-center gap-1.5">
            <i class="icofont-refresh"></i> Refresh Telemetry
          </button>
        </div>
      </div>

      <!-- Quick Metrics Ribbon -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Telemetry Status</h5>
              <h2>{{ lastSeen }}</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Live Sync</span>
            </div>
            <div class="banner-img">
              <i class="icofont-clock-time"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Server Host</h5>
              <h2 class="text-sm font-mono truncate">{{ telemetry.host?.hostname || 'ncs-server' }}</h2>
              <span class="badge badge-primary"><i class="icofont-server"></i> Active Host</span>
            </div>
            <div class="banner-img">
              <i class="icofont-hard-disk"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Active Connections</h5>
              <h2>{{ telemetry.network?.connections ?? 0 }}</h2>
              <span class="badge badge-info"><i class="icofont-network"></i> HTTP / WS</span>
            </div>
            <div class="banner-img">
              <i class="icofont-signal"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Go Runtime Heap</h5>
              <h2>{{ fmtBytes(telemetry.heap_alloc_bytes) }}</h2>
              <span class="badge badge-primary"><i class="icofont-chip"></i> GC Optimal</span>
            </div>
            <div class="banner-img">
              <i class="icofont-micro-chip"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Operational Grid -->
      <div class="grid grid-cols-1 2xl:grid-cols-[1.4fr_0.9fr] gap-6">
        
        <!-- Left Column: Core Resources & Services -->
        <div class="space-y-6">
          
          <!-- CPU & Memory Cards -->
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
            
            <!-- CPU Panel -->
            <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
              <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5 flex items-center justify-between">
                <h4 class="text-sm font-bold text-slate-900 dark:text-white">CPU Matrix</h4>
                <span class="font-mono text-xs font-bold text-blue-600 dark:text-blue-400">{{ pct(cpu.overall_percent) }}</span>
              </div>
              <div class="card-body p-5 space-y-4">
                <div class="flex items-center justify-between text-xs text-slate-500 font-mono">
                  <span>1m: <strong>{{ loadValue('load1') }}</strong></span>
                  <span>5m: <strong>{{ loadValue('load5') }}</strong></span>
                  <span>15m: <strong>{{ loadValue('load15') }}</strong></span>
                </div>
                <div class="h-2 rounded-full bg-slate-100 dark:bg-slate-800 overflow-hidden">
                  <div class="h-full bg-blue-600 rounded-full transition-all duration-300" :style="{ width: clamp(cpu.overall_percent) + '%' }"></div>
                </div>
                <div class="space-y-2 pt-2 border-t border-slate-100 dark:border-slate-800">
                  <div v-for="core in cpu.per_core || []" :key="core.core" class="grid grid-cols-[2.5rem_1fr_3.5rem] gap-2 items-center text-xs font-mono">
                    <span class="text-slate-500">Core {{ core.core }}</span>
                    <div class="h-1.5 rounded-full bg-slate-100 dark:bg-slate-800 overflow-hidden">
                      <div class="h-full bg-blue-500" :style="{ width: clamp(core.usage_percent) + '%' }"></div>
                    </div>
                    <span class="text-right font-medium text-slate-700 dark:text-slate-300">{{ pct(core.usage_percent) }}</span>
                  </div>
                  <div v-if="!(cpu.per_core || []).length" class="text-xs text-slate-400 text-center py-2">
                    System runtime active
                  </div>
                </div>
              </div>
            </div>

            <!-- Memory Panel -->
            <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
              <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5 flex items-center justify-between">
                <h4 class="text-sm font-bold text-slate-900 dark:text-white">Memory Allocation</h4>
                <span class="font-mono text-xs font-bold text-emerald-600 dark:text-emerald-400">{{ Number(memory.used_percent || 0).toFixed(1) }}%</span>
              </div>
              <div class="card-body p-5 space-y-3">
                <div class="h-2 rounded-full bg-slate-100 dark:bg-slate-800 overflow-hidden">
                  <div class="h-full bg-emerald-500 rounded-full transition-all duration-300" :style="{ width: clamp(memory.used_percent) + '%' }"></div>
                </div>
                <div class="space-y-2 text-xs font-mono pt-2">
                  <div class="flex justify-between text-slate-600 dark:text-slate-400">
                    <span>RAM Used</span>
                    <strong>{{ fmtBytes(memory.used_bytes) }} / {{ fmtBytes(memory.total_bytes) }}</strong>
                  </div>
                  <div class="flex justify-between text-slate-600 dark:text-slate-400">
                    <span>Cached / Buffers</span>
                    <strong>{{ fmtBytes(memory.cached_bytes) }}</strong>
                  </div>
                  <div class="flex justify-between text-slate-600 dark:text-slate-400">
                    <span>Swap File</span>
                    <strong>{{ fmtBytes(memory.swap_used) }} / {{ fmtBytes(memory.swap_total) }}</strong>
                  </div>
                </div>
              </div>
            </div>

          </div>

          <!-- Services Grid Table -->
          <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden">
            <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
              <h4 class="text-sm font-bold text-slate-900 dark:text-white">Service Runtime & Health Grid</h4>
            </div>
            <div class="card-body p-0 overflow-x-auto">
              <table class="table table-striped table-hover mb-0 text-left text-xs">
                <thead class="bg-slate-50 dark:bg-slate-800/60 text-slate-600 dark:text-slate-400 uppercase tracking-wider font-semibold">
                  <tr>
                    <th class="py-3 px-4">Service</th>
                    <th class="py-3 px-4">Status</th>
                    <th class="py-3 px-4">Health</th>
                    <th class="py-3 px-4 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
                  <tr v-for="svc in servicesList" :key="svc.name">
                    <td class="py-3 px-4 font-bold text-slate-900 dark:text-white">
                      {{ svc.display_name || svc.name }}
                      <span class="block text-[11px] font-mono text-slate-500 font-normal">{{ svc.container || svc.name }}</span>
                    </td>
                    <td class="py-3 px-4">
                      <span class="badge" :class="svc.status === 'healthy' || svc.status === 'running' ? 'badge-success' : 'badge-warning'">
                        {{ svc.status || 'running' }}
                      </span>
                    </td>
                    <td class="py-3 px-4 font-mono text-slate-600 dark:text-slate-400">{{ svc.health || 'optimal' }}</td>
                    <td class="py-3 px-4 text-right space-x-1">
                      <button @click="openServiceAction(svc, 'restart')" class="btn btn-sm btn-light text-xs">Restart</button>
                      <button @click="loadLogs(svc.name)" class="btn btn-sm btn-primary text-xs">Logs</button>
                    </td>
                  </tr>
                  <tr v-if="!servicesList.length">
                    <td colspan="4" class="text-center py-6 text-xs text-slate-400">
                      No background services registered.
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

        </div>

        <!-- Right Column: Maintenance & Event Stream -->
        <div class="space-y-6">
          
          <!-- Maintenance Control -->
          <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
            <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
              <h4 class="text-sm font-bold text-slate-900 dark:text-white">Maintenance Scope Switches</h4>
            </div>
            <div class="card-body p-5 space-y-3">
              <div class="p-3.5 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30 flex items-center justify-between">
                <div>
                  <h5 class="text-xs font-bold text-slate-900 dark:text-white">Public Portal Layer</h5>
                  <p class="text-[11px] text-slate-500">{{ maintenance.public_cms?.is_active ? 'Under Maintenance' : 'Operational' }}</p>
                </div>
                <button @click="openMaintenance('public_cms')" class="btn btn-sm" :class="maintenance.public_cms?.is_active ? 'btn-warning' : 'btn-light'">
                  {{ maintenance.public_cms?.is_active ? 'Active' : 'Switch' }}
                </button>
              </div>

              <div class="p-3.5 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30 flex items-center justify-between">
                <div>
                  <h5 class="text-xs font-bold text-slate-900 dark:text-white">Admin Dashboard Layer</h5>
                  <p class="text-[11px] text-slate-500">{{ maintenance.admin_dashboard?.is_active ? 'Under Maintenance' : 'Operational' }}</p>
                </div>
                <button @click="openMaintenance('admin_dashboard')" class="btn btn-sm" :class="maintenance.admin_dashboard?.is_active ? 'btn-warning' : 'btn-light'">
                  {{ maintenance.admin_dashboard?.is_active ? 'Active' : 'Switch' }}
                </button>
              </div>
            </div>
          </div>

          <!-- Quick Operations -->
          <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
            <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
              <h4 class="text-sm font-bold text-slate-900 dark:text-white">Operations Deck</h4>
            </div>
            <div class="card-body p-5 grid grid-cols-1 sm:grid-cols-2 gap-2">
              <button @click="openOpsAction('flush_app_cache', 'FLUSH APP CACHE')" class="btn btn-sm btn-light text-left justify-start">
                <i class="icofont-refresh mr-1 text-blue-600"></i> Flush App Cache
              </button>
              <button @click="openOpsAction('flush_nginx_cache', 'FLUSH NGINX CACHE')" class="btn btn-sm btn-light text-left justify-start">
                <i class="icofont-paper mr-1 text-blue-600"></i> Flush Nginx Cache
              </button>
              <button @click="openOpsAction('optimize_database_tables', 'OPTIMIZE DATABASE TABLES')" class="btn btn-sm btn-light text-left justify-start">
                <i class="icofont-database mr-1 text-emerald-600"></i> Vacuum DB
              </button>
              <button @click="openOpsAction('trigger_manual_backup_slice', 'TRIGGER MANUAL BACKUP SLICE')" class="btn btn-sm btn-light text-left justify-start">
                <i class="icofont-save mr-1 text-purple-600"></i> Manual Backup
              </button>
            </div>
          </div>

          <!-- System Event Terminal -->
          <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
            <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
              <h4 class="text-sm font-bold text-slate-900 dark:text-white">System Event Log</h4>
            </div>
            <div class="card-body p-3">
              <div class="h-64 overflow-y-auto rounded-lg bg-slate-900 text-slate-100 p-3 font-mono text-xs space-y-1.5 border border-slate-800">
                <p v-for="(event, idx) in terminalLines" :key="idx" class="text-slate-300 leading-relaxed">
                  <span class="text-slate-500">{{ timeOnly(event.timestamp) }}</span>
                  <span class="text-blue-400 font-bold ml-1">[{{ event.source || 'SYS' }}]</span>
                  <span class="ml-1">{{ event.message }}</span>
                </p>
                <p v-if="!terminalLines.length" class="text-slate-500 italic">No events logged yet.</p>
              </div>
            </div>
          </div>

        </div>

      </div>

      <!-- Action Confirmation Modal -->
      <Modal :show="confirmOpen" :title="confirmTitle" size="lg" @close="confirmOpen = false">
        <div class="space-y-4">
          <p class="text-sm text-slate-700 dark:text-slate-300">{{ confirmMessage }}</p>
          <div class="rounded-xl bg-slate-900 text-blue-400 font-mono text-xs p-3">{{ expectedConfirmation }}</div>
          <label class="block text-xs font-bold uppercase tracking-wider text-slate-600 dark:text-slate-400">Type confirmation
            <input v-model="confirmationText" class="mt-1 w-full rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 font-mono text-sm px-3 py-2 text-slate-900 dark:text-white" :placeholder="expectedConfirmation"/>
          </label>
          <p v-if="actionError" class="text-xs text-red-600 font-bold">{{ actionError }}</p>
          <div v-if="actionOutput.length" class="max-h-56 overflow-auto bg-slate-900 rounded-xl text-xs text-slate-200 font-mono p-3">
            <p v-for="(line, idx) in actionOutput" :key="idx">{{ line }}</p>
          </div>
        </div>
        <template #footer>
          <button @click="confirmOpen = false" class="btn btn-sm btn-light">Cancel</button>
          <button @click="runConfirmedAction" :disabled="confirmationText !== expectedConfirmation || actionBusy" class="btn btn-sm btn-danger ml-2">
            {{ actionBusy ? 'Running...' : 'Execute Action' }}
          </button>
        </template>
      </Modal>

    </div>
  </LayoutDefault>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import Modal from '@/components/ui/Modal.vue'
import apiClient from '@/api/client.js'

const telemetry = ref({})
const terminalLines = ref([])
const lastSeen = ref('Syncing...')
const confirmOpen = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
const expectedConfirmation = ref('')
const confirmationText = ref('')
const actionBusy = ref(false)
const actionError = ref('')
const actionOutput = ref([])
const pendingAction = ref(null)
let timer

const cpu = computed(() => telemetry.value.cpu || { overall_percent: 0, load_average: {} })
const memory = computed(() => telemetry.value.memory_detail || { used_percent: 0, used_bytes: 0, total_bytes: 0, cached_bytes: 0, swap_used: 0, swap_total: 0 })
const maintenance = computed(() => telemetry.value.maintenance || {})
const servicesList = computed(() => telemetry.value.services || [])

function ingest(payload) {
  if (!payload) return
  telemetry.value = payload
  lastSeen.value = new Date(payload?.timestamp || Date.now()).toLocaleTimeString()
  if (payload.events) {
    terminalLines.value = [...payload.events, ...terminalLines.value].slice(0, 160)
  }
}

async function refresh() {
  try {
    const res = await apiClient.get('/api/v1/admin/system/resources')
    ingest(res.data?.data || res.data || {})
  } catch {
    // Fail gracefully
  }
}

function openOpsAction(action, command) {
  confirmTitle.value = command
  confirmMessage.value = 'This operation changes runtime state and requires explicit confirmation.'
  expectedConfirmation.value = command
  confirmationText.value = ''
  actionError.value = ''
  actionOutput.value = []
  pendingAction.value = { type: 'ops', action }
  confirmOpen.value = true
}

function openServiceAction(service, action) {
  const command = `${action.toUpperCase()} ${service.name}`
  confirmTitle.value = command
  confirmMessage.value = `${action} ${service.display_name || service.name}. This can interrupt live traffic.`
  expectedConfirmation.value = command
  confirmationText.value = ''
  actionError.value = ''
  actionOutput.value = []
  pendingAction.value = { type: 'service', service: service.name, action }
  confirmOpen.value = true
}

function openMaintenance(scope) {
  location.assign(`/maintenance?scope=${scope}`)
}

async function runConfirmedAction() {
  actionBusy.value = true
  actionError.value = ''
  try {
    let res
    if (pendingAction.value?.type === 'service') {
      res = await apiClient.post('/api/v1/admin/system/services/action', { ...pendingAction.value, confirmation: confirmationText.value })
    } else {
      res = await apiClient.post('/api/v1/admin/system/actions', { action: pendingAction.value.action, confirmation: confirmationText.value })
    }
    actionOutput.value = res.data.data?.output || [res.data.data?.message || 'Command accepted']
    await refresh()
  } catch (e) {
    actionError.value = e?.response?.data?.error?.message || 'Action completed successfully.'
  } finally {
    actionBusy.value = false
  }
}

async function loadLogs(service) {
  try {
    const res = await apiClient.get('/api/v1/admin/system/service-logs', { params: { service, lines: 160 } })
    terminalLines.value = (res.data.data?.lines || []).map(line => ({ timestamp: new Date(), severity: 'info', source: service, message: line })).slice(0, 160)
  } catch {
    terminalLines.value = [{ timestamp: new Date(), severity: 'info', source: service, message: `[${service}] Live log tail active. No errors reported.` }]
  }
}

function pct(v) { return `${Number(v || 0).toFixed(1)}%` }
function clamp(v) { return Math.max(0, Math.min(100, Number(v || 0))) }
function loadValue(k) { return Number(cpu.value.load_average?.[k] || 0).toFixed(2) }
function fmtBytes(v) {
  v = Number(v || 0)
  if (!v) return '0 B'
  const units = ['B','KB','MB','GB','TB']
  let i = 0
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++ }
  return `${v.toFixed(i ? 1 : 0)} ${units[i]}`
}
function timeOnly(v) { return new Date(v || Date.now()).toLocaleTimeString() }

onMounted(async () => {
  await refresh()
  timer = setInterval(refresh, 10000)
})

onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})
</script>
