<template>
  <LayoutDefault title="Infrastructure Command Center">
    <div class="min-h-screen bg-zinc-950 text-zinc-100 p-4 lg:p-6 font-sans">
      <header class="flex flex-col xl:flex-row xl:items-end justify-between gap-4 border-b border-zinc-800 pb-4">
        <div>
          <p class="text-xs uppercase tracking-widest text-cyan-300 font-bold">Systems Maintenance Command Center</p>
          <h1 class="text-2xl lg:text-3xl font-black mt-1">VPS Telemetry & Runtime Control</h1>
          <p class="text-sm text-zinc-400 mt-2">Live host metrics, scoped maintenance switches, service operations, and event streams.</p>
        </div>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 text-xs">
          <MetricPill label="Updated" :value="lastSeen" tone="cyan" />
          <MetricPill label="Host" :value="telemetry.host?.hostname || 'unknown'" tone="blue" />
          <MetricPill label="Connections" :value="String(telemetry.network?.connections ?? 0)" tone="emerald" />
          <MetricPill label="Heap" :value="fmtBytes(telemetry.heap_alloc_bytes)" tone="amber" />
        </div>
      </header>

      <main class="grid 2xl:grid-cols-[1.4fr_0.9fr] gap-4 mt-4">
        <section class="space-y-4">
          <div class="grid xl:grid-cols-3 gap-4">
            <Panel title="CPU Matrix" class="xl:col-span-2">
              <div class="grid lg:grid-cols-[1fr_13rem] gap-4">
                <div>
                  <div class="flex items-end justify-between">
                    <div>
                      <p class="text-4xl font-black font-mono text-cyan-200">{{ pct(cpu.overall_percent) }}</p>
                      <p class="text-xs text-zinc-400">overall utilization</p>
                    </div>
                    <div class="text-right font-mono text-xs text-zinc-300">
                      <p>1m {{ loadValue('load1') }}</p>
                      <p>5m {{ loadValue('load5') }}</p>
                      <p>15m {{ loadValue('load15') }}</p>
                    </div>
                  </div>
                  <CpuSparkline :points="cpuHistory" class="mt-4" />
                </div>
                <div class="space-y-2">
                  <div v-for="core in cpu.per_core || []" :key="core.core" class="grid grid-cols-[2.5rem_1fr_3.5rem] gap-2 items-center text-xs font-mono">
                    <span class="text-zinc-400">C{{ core.core }}</span>
                    <div class="h-2 rounded-full bg-zinc-800 overflow-hidden"><div class="h-full bg-cyan-400" :style="{ width: clamp(core.usage_percent) + '%' }"/></div>
                    <span class="text-right">{{ pct(core.usage_percent) }}</span>
                  </div>
                  <div class="pt-2 border-t border-zinc-800 text-xs font-mono text-zinc-300">Temp {{ cpu.temperature_c == null ? 'N/A' : cpu.temperature_c.toFixed(1) + 'C' }}</div>
                </div>
              </div>
            </Panel>

            <Panel title="Memory Gauge">
              <div class="space-y-4">
                <RingGauge :value="memory.used_percent || 0" label="RAM" />
                <div class="space-y-2 text-xs font-mono">
                  <Meter label="Used" :value="memory.used_bytes" :total="memory.total_bytes" tone="emerald" />
                  <Meter label="Cached" :value="memory.cached_bytes" :total="memory.total_bytes" tone="blue" />
                  <Meter label="Buffers" :value="memory.buffers_bytes" :total="memory.total_bytes" tone="amber" />
                  <Meter label="Swap" :value="memory.swap_used" :total="memory.swap_total" tone="red" />
                </div>
              </div>
            </Panel>
          </div>

          <div class="grid xl:grid-cols-2 gap-4">
            <Panel title="Storage & Disk I/O">
              <div class="grid gap-3">
                <div v-for="partition in telemetry.storage?.partitions || []" :key="partition.mountpoint" class="border border-zinc-800 rounded-lg p-3">
                  <div class="flex justify-between text-sm font-mono">
                    <span class="text-zinc-200">{{ partition.mountpoint }}</span>
                    <span class="text-zinc-400">{{ pct(partition.used_percent) }}</span>
                  </div>
                  <div class="h-2 bg-zinc-800 rounded-full overflow-hidden mt-2"><div class="h-full bg-violet-400" :style="{ width: clamp(partition.used_percent) + '%' }"/></div>
                  <p class="text-xs text-zinc-500 mt-2">{{ fmtBytes(partition.used_bytes) }} / {{ fmtBytes(partition.total_bytes) }}</p>
                </div>
              </div>
              <div class="grid grid-cols-2 gap-3 mt-4">
                <MetricPill label="Read MB/s" :value="num(telemetry.storage?.read_mb_sec)" tone="cyan" />
                <MetricPill label="Write MB/s" :value="num(telemetry.storage?.write_mb_sec)" tone="amber" />
              </div>
            </Panel>

            <Panel title="Network & Bandwidth">
              <div class="grid sm:grid-cols-2 gap-3">
                <SpeedGauge title="Ingress" :value="telemetry.network?.ingress_mbps || 0" tone="cyan" />
                <SpeedGauge title="Egress" :value="telemetry.network?.egress_mbps || 0" tone="emerald" />
              </div>
              <div class="mt-4 space-y-2 max-h-44 overflow-auto pr-1">
                <div v-for="iface in telemetry.network?.interfaces || []" :key="iface.name" class="grid grid-cols-[1fr_auto] gap-3 text-xs font-mono border-b border-zinc-900 pb-2">
                  <span class="text-zinc-300">{{ iface.name }}</span>
                  <span class="text-zinc-500">{{ fmtBytes(iface.bytesRecv) }} in / {{ fmtBytes(iface.bytesSent) }} out</span>
                </div>
              </div>
            </Panel>
          </div>

          <Panel title="Process & Service Grid">
            <div class="overflow-x-auto">
              <table class="w-full text-sm">
                <thead class="text-xs uppercase tracking-wide text-zinc-500 border-b border-zinc-800">
                  <tr><th class="text-left py-2">Service</th><th class="text-left py-2">Status</th><th class="text-left py-2">Health</th><th class="text-right py-2">Actions</th></tr>
                </thead>
                <tbody class="divide-y divide-zinc-900">
                  <tr v-for="svc in telemetry.services || []" :key="svc.name">
                    <td class="py-3 font-semibold">{{ svc.display_name || svc.name }}<p class="text-xs text-zinc-500 font-mono">{{ svc.container || svc.name }}</p></td>
                    <td><StatusDot :status="svc.status" /></td>
                    <td class="text-zinc-400 font-mono text-xs">{{ svc.health || 'unknown' }}</td>
                    <td class="text-right">
                      <button @click="openServiceAction(svc, 'restart')" class="ops-btn border-cyan-500/40 text-cyan-200">Restart</button>
                      <button @click="openServiceAction(svc, 'stop')" class="ops-btn border-red-500/40 text-red-200">Stop</button>
                      <button @click="loadLogs(svc.name)" class="ops-btn border-zinc-600 text-zinc-200">Logs</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </Panel>
        </section>

        <aside class="space-y-4">
          <Panel title="Maintenance Status Core">
            <div class="space-y-3">
              <MaintenanceScope label="Public CMS Layer" :scope="maintenance.public_cms" @toggle="openMaintenance('public_cms')" />
              <MaintenanceScope label="Admin Dashboard Layer" :scope="maintenance.admin_dashboard" @toggle="openMaintenance('admin_dashboard')" />
            </div>
          </Panel>

          <Panel title="Cache & Database Deck">
            <div class="grid gap-2">
              <ActionButton label="Flush application cache" command="FLUSH APP CACHE" @click="openOpsAction('flush_app_cache', 'FLUSH APP CACHE')" />
              <ActionButton label="Flush Nginx micro-cache" command="FLUSH NGINX CACHE" @click="openOpsAction('flush_nginx_cache', 'FLUSH NGINX CACHE')" />
              <ActionButton label="Flush Redis database" command="FLUSH REDIS CACHE" @click="openOpsAction('flush_redis_cache', 'FLUSH REDIS CACHE')" />
              <ActionButton label="Optimize database tables" command="OPTIMIZE DATABASE TABLES" danger @click="openOpsAction('optimize_database_tables', 'OPTIMIZE DATABASE TABLES')" />
              <ActionButton label="Prune activity logs" command="PRUNE ACTIVITY LOGS" danger @click="openOpsAction('prune_activity_logs', 'PRUNE ACTIVITY LOGS')" />
              <ActionButton label="Trigger manual backup slice" command="TRIGGER MANUAL BACKUP SLICE" @click="openOpsAction('trigger_manual_backup_slice', 'TRIGGER MANUAL BACKUP SLICE')" />
            </div>
          </Panel>

          <Panel title="Live System Event Stream">
            <div class="h-96 overflow-auto rounded-lg bg-black border border-zinc-800 p-3 font-mono text-xs space-y-1">
              <p v-for="(event, idx) in terminalLines" :key="idx" :class="eventClass(event.severity)">
                <span class="text-zinc-600">{{ timeOnly(event.timestamp) }}</span>
                <span class="uppercase">[{{ event.severity }}]</span>
                <span class="text-zinc-500">{{ event.source }}</span>
                <span>{{ event.message }}</span>
              </p>
            </div>
          </Panel>
        </aside>
      </main>

      <Modal :show="confirmOpen" :title="confirmTitle" size="lg" @close="confirmOpen = false">
        <div class="space-y-4">
          <p class="text-sm text-gray-700">{{ confirmMessage }}</p>
          <div class="rounded-xl bg-gray-950 text-cyan-200 font-mono text-xs p-3">{{ expectedConfirmation }}</div>
          <label class="block text-sm font-semibold">Type confirmation
            <input v-model="confirmationText" class="mt-1 w-full rounded-xl border-gray-200 font-mono text-sm" :placeholder="expectedConfirmation"/>
          </label>
          <p v-if="actionError" class="text-xs text-red-600">{{ actionError }}</p>
          <div v-if="actionOutput.length" class="max-h-56 overflow-auto bg-gray-950 rounded-xl text-xs text-gray-200 font-mono p-3">
            <p v-for="(line, idx) in actionOutput" :key="idx">{{ line }}</p>
          </div>
        </div>
        <template #footer>
          <button @click="confirmOpen = false" class="px-4 py-2 text-sm">Cancel</button>
          <button @click="runConfirmedAction" :disabled="confirmationText !== expectedConfirmation || actionBusy" class="rounded-xl bg-red-700 hover:bg-red-800 disabled:opacity-50 px-5 py-2 text-sm text-white font-bold">
            {{ actionBusy ? 'Running...' : 'Execute' }}
          </button>
        </template>
      </Modal>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { computed, defineComponent, h, onBeforeUnmount, onMounted, ref } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import Modal from '@/components/ui/Modal.vue'
import apiClient from '@/api/client.js'

const telemetry = ref({})
const cpuHistory = ref([])
const terminalLines = ref([])
const lastSeen = ref('waiting')
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
let socket

const cpu = computed(() => telemetry.value.cpu || {})
const memory = computed(() => telemetry.value.memory_detail || {})
const maintenance = computed(() => telemetry.value.maintenance || {})

function ingest(payload) {
  telemetry.value = payload || {}
  lastSeen.value = new Date(payload?.timestamp || Date.now()).toLocaleTimeString()
  const value = Number(payload?.cpu?.overall_percent || 0)
  cpuHistory.value = [...cpuHistory.value, { t: Date.now(), v: value }].slice(-180)
  terminalLines.value = [...(payload?.events || []), ...terminalLines.value].slice(0, 160)
}
async function refresh() {
  const res = await apiClient.get('/api/v1/admin/system/resources')
  ingest(res.data.data || {})
}
function connectStream() {
  const token = localStorage.getItem('ncsms_access_token')
  const proto = location.protocol === 'https:' ? 'wss' : 'ws'
  socket = new WebSocket(`${proto}://${location.host}/api/v1/admin/system/resources/ws?access_token=${encodeURIComponent(token || '')}`)
  socket.onmessage = event => { try { ingest(JSON.parse(event.data)) } catch {} }
  socket.onerror = () => { socket?.close() }
  socket.onclose = () => { if (!timer) timer = setInterval(refresh, 5000) }
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
    actionError.value = e?.response?.data?.error?.message || 'Action failed'
  } finally {
    actionBusy.value = false
  }
}
async function loadLogs(service) {
  const res = await apiClient.get('/api/v1/admin/system/service-logs', { params: { service, lines: 160 } })
  terminalLines.value = (res.data.data?.lines || []).map(line => ({ timestamp: new Date(), severity: classifyLog(line), source: service, message: line })).slice(0, 160)
}
function classifyLog(line) {
  const text = String(line).toLowerCase()
  if (text.includes('error') || text.includes('fatal') || text.includes('panic')) return 'critical'
  if (text.includes('warn') || text.includes('503')) return 'warning'
  return 'info'
}
function pct(v) { return `${Number(v || 0).toFixed(1)}%` }
function num(v) { return Number(v || 0).toFixed(2) }
function clamp(v) { return Math.max(0, Math.min(100, Number(v || 0))) }
function loadValue(k) { return Number(cpu.value.load_average?.[k] || cpu.value.load_average?.[k.replace('load', 'load')] || 0).toFixed(2) }
function fmtBytes(v) {
  v = Number(v || 0)
  const units = ['B','KB','MB','GB','TB']
  let i = 0
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++ }
  return `${v.toFixed(i ? 1 : 0)} ${units[i]}`
}
function timeOnly(v) { return new Date(v || Date.now()).toLocaleTimeString() }
function durationSince(iso) {
  if (!iso) return 'unknown'
  const seconds = Math.max(0, (Date.now() - new Date(iso).getTime()) / 1000)
  if (seconds < 60) return `${Math.floor(seconds)}s`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ${Math.floor((seconds % 3600) / 60)}m`
  return `${Math.floor(seconds / 86400)}d`
}
function eventClass(sev) {
  if (sev === 'critical') return 'text-red-300 drop-shadow-[0_0_6px_rgba(248,113,113,.8)]'
  if (sev === 'warning') return 'text-amber-300'
  return 'text-emerald-300'
}
onMounted(async () => {
  await refresh()
  connectStream()
})
onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
  if (socket) socket.close()
})

const Panel = defineComponent({
  props: { title: String },
  setup(props, { slots }) {
    return () => h('section', { class: 'rounded-lg border border-zinc-800 bg-zinc-900/70 shadow-xl shadow-black/20' }, [
      h('div', { class: 'px-4 py-3 border-b border-zinc-800 flex items-center justify-between' }, h('h2', { class: 'font-bold text-sm uppercase tracking-wide text-zinc-200' }, props.title)),
      h('div', { class: 'p-4' }, slots.default?.())
    ])
  }
})
const toneText = tone => ({
  cyan: 'text-cyan-200',
  blue: 'text-blue-200',
  emerald: 'text-emerald-200',
  amber: 'text-amber-200',
  red: 'text-red-200',
  violet: 'text-violet-200',
})[tone] || 'text-cyan-200'
const toneBg = tone => ({
  cyan: 'bg-cyan-400',
  blue: 'bg-blue-400',
  emerald: 'bg-emerald-400',
  amber: 'bg-amber-400',
  red: 'bg-red-400',
  violet: 'bg-violet-400',
})[tone] || 'bg-cyan-400'
const MetricPill = defineComponent({
  props: { label: String, value: String, tone: String },
  setup(props) {
    return () => h('div', { class: 'rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2' }, [
      h('p', { class: 'text-[10px] uppercase tracking-wide text-zinc-500 font-bold' }, props.label),
      h('p', { class: `font-mono text-sm truncate ${toneText(props.tone)}` }, props.value)
    ])
  }
})
const CpuSparkline = defineComponent({
  props: { points: Array },
  setup(props) {
    return () => h('div', { class: 'h-28 border border-zinc-800 rounded-lg bg-black/50 p-2 flex items-end gap-px' }, (props.points || []).map((p, i) => h('div', { key: i, class: 'flex-1 bg-cyan-400/80 min-w-[2px]', style: { height: `${Math.max(2, p.v)}%` } })))
  }
})
const RingGauge = defineComponent({
  props: { value: Number, label: String },
  setup(props) {
    return () => h('div', { class: 'flex items-center gap-4' }, [
      h('div', { class: 'w-24 h-24 rounded-full grid place-items-center border-8 border-emerald-400 bg-zinc-950' }, h('span', { class: 'font-mono text-lg font-bold' }, `${Number(props.value || 0).toFixed(0)}%`)),
      h('div', [h('p', { class: 'font-bold' }, props.label), h('p', { class: 'text-xs text-zinc-500' }, 'active allocation')])
    ])
  }
})
const Meter = defineComponent({
  props: { label: String, value: Number, total: Number, tone: String },
  setup(props) {
    const percent = () => props.total ? (props.value / props.total) * 100 : 0
    return () => h('div', [
      h('div', { class: 'flex justify-between text-zinc-400' }, [h('span', props.label), h('span', fmtBytes(props.value))]),
      h('div', { class: 'h-1.5 bg-zinc-800 rounded-full overflow-hidden mt-1' }, h('div', { class: `h-full ${toneBg(props.tone)}`, style: { width: `${clamp(percent())}%` } }))
    ])
  }
})
const SpeedGauge = defineComponent({
  props: { title: String, value: Number, tone: String },
  setup(props) {
    return () => h('div', { class: 'rounded-lg border border-zinc-800 bg-black/30 p-4' }, [
      h('p', { class: 'text-xs text-zinc-500 uppercase font-bold' }, props.title),
      h('p', { class: `font-mono text-3xl font-black mt-2 ${toneText(props.tone)}` }, `${Number(props.value || 0).toFixed(2)}`),
      h('p', { class: 'text-xs text-zinc-500' }, 'Mbps')
    ])
  }
})
const StatusDot = defineComponent({
  props: { status: String },
  setup(props) {
    const cls = () => props.status === 'healthy' ? 'bg-emerald-400 text-emerald-950' : props.status === 'degraded' ? 'bg-amber-400 text-amber-950' : props.status === 'stopped' ? 'bg-red-500 text-white' : 'bg-zinc-700 text-zinc-200'
    return () => h('span', { class: `inline-flex items-center rounded-full px-2 py-1 text-xs font-bold ${cls()}` }, props.status || 'unknown')
  }
})
const MaintenanceScope = defineComponent({
  props: { label: String, scope: Object },
  emits: ['toggle'],
  setup(props, { emit }) {
    return () => h('div', { class: 'rounded-lg border border-zinc-800 p-3 bg-black/20' }, [
      h('div', { class: 'flex items-center justify-between gap-3' }, [
        h('div', [
          h('p', { class: 'font-bold' }, props.label),
          h('p', { class: 'text-xs text-zinc-500' }, props.scope?.reason || 'Operational'),
          h('p', { class: 'text-[10px] text-zinc-600 font-mono mt-1' }, props.scope?.is_active ? `active for ${durationSince(props.scope?.changed_at)}` : `last changed ${durationSince(props.scope?.changed_at)} ago`)
        ]),
        h('button', { onClick: () => emit('toggle'), class: `px-3 py-2 rounded-lg text-xs font-bold ${props.scope?.is_active ? 'bg-amber-500 text-black' : 'bg-emerald-500 text-black'}` }, props.scope?.is_active ? 'Active' : 'Operational')
      ])
    ])
  }
})
const ActionButton = defineComponent({
  props: { label: String, command: String, danger: Boolean },
  emits: ['click'],
  setup(props, { emit }) {
    return () => h('button', { onClick: () => emit('click'), class: `text-left rounded-lg border px-3 py-2 ${props.danger ? 'border-red-500/40 hover:bg-red-500/10' : 'border-zinc-700 hover:bg-zinc-800'}` }, [
      h('p', { class: 'font-semibold text-sm' }, props.label),
      h('p', { class: 'font-mono text-[10px] text-zinc-500 mt-1' }, props.command)
    ])
  }
})
</script>

<style scoped>
.ops-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-width: 1px;
  border-radius: 0.5rem;
  padding: 0.35rem 0.55rem;
  margin-left: 0.35rem;
  font-size: 0.75rem;
  font-weight: 700;
}
</style>
