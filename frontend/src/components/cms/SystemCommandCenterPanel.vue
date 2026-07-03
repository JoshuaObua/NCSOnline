<template>
  <section class="ops-panel">
    <div class="cms-panel-head">
      <h2>System Command Center</h2>
      <label class="ops-live-toggle"><input v-model="liveResources" type="checkbox" /> Live refresh</label>
    </div>

    <div class="ops-grid">
      <article class="cms-card metric">
        <span>CPU</span>
        <strong>{{ round(resources?.cpu?.overall_percent) }}%</strong>
        <div class="ops-bar"><span :style="{ width: round(resources?.cpu?.overall_percent) + '%' }"></span></div>
      </article>
      <article class="cms-card metric">
        <span>Memory</span>
        <strong>{{ round(resources?.memory_detail?.used_percent) }}%</strong>
        <div class="ops-bar"><span :style="{ width: round(resources?.memory_detail?.used_percent) + '%' }"></span></div>
        <small>{{ formatBytes(resources?.memory_detail?.used_bytes) }} / {{ formatBytes(resources?.memory_detail?.total_bytes) }}</small>
      </article>
      <article class="cms-card metric">
        <span>Disk</span>
        <strong>{{ round(resources?.disk?.used_percent) }}%</strong>
        <div class="ops-bar"><span :style="{ width: round(resources?.disk?.used_percent) + '%' }"></span></div>
        <small>{{ formatBytes(resources?.disk?.used_bytes) }} / {{ formatBytes(resources?.disk?.total_bytes) }}</small>
      </article>
      <article class="cms-card metric">
        <span>Network</span>
        <strong>{{ (resources?.network?.ingress_mbps || 0).toFixed(2) }} Mbps</strong>
        <small>in · {{ (resources?.network?.egress_mbps || 0).toFixed(2) }} Mbps out</small>
      </article>
    </div>

    <article class="cms-panel">
      <div class="cms-panel-head"><h3>Services</h3></div>
      <table class="ops-table">
        <thead><tr><th>Service</th><th>Status</th><th>Health</th><th></th></tr></thead>
        <tbody>
          <tr v-for="svc in resources?.services || []" :key="svc.name">
            <td>{{ svc.display_name }}</td>
            <td><span class="ops-badge" :class="svc.status">{{ svc.status }}</span></td>
            <td>{{ svc.health || '—' }}</td>
            <td class="ops-actions">
              <button type="button" :disabled="busyService === svc.name" @click="serviceAction(svc.name, 'start')">Start</button>
              <button type="button" :disabled="busyService === svc.name" @click="serviceAction(svc.name, 'restart')">Restart</button>
              <button type="button" class="ops-danger" :disabled="busyService === svc.name" @click="serviceAction(svc.name, 'stop')">Stop</button>
              <button type="button" @click="loadLogs(svc.name)">Logs</button>
            </td>
          </tr>
          <tr v-if="!resources?.services?.length"><td colspan="4" class="cms-empty">No service telemetry yet.</td></tr>
        </tbody>
      </table>
      <p v-if="dockerUnavailable" class="ops-hint">Docker socket is not mounted into this container, so start/stop/restart and log streaming are unavailable — metrics above still work.</p>
    </article>

    <article class="cms-panel">
      <div class="cms-panel-head">
        <h3>Service Logs {{ activeLogService ? `— ${activeLogService}` : '' }}</h3>
        <label class="ops-live-toggle"><input v-model="liveLogs" type="checkbox" :disabled="!activeLogService" /> Live tail</label>
      </div>
      <pre class="ops-terminal">{{ logLines.join('\n') || 'Select a service above and click "Logs" to view recent output.' }}</pre>
    </article>

    <article class="cms-panel">
      <div class="cms-panel-head"><h3>Recent Events</h3></div>
      <ul class="ops-events">
        <li v-for="(evt, idx) in resources?.events || []" :key="idx" :class="evt.severity">
          <span class="ops-event-dot"></span>
          <span>{{ evt.message }}</span>
          <small>{{ formatDateTime(evt.timestamp) }}</small>
        </li>
        <li v-if="!resources?.events?.length" class="cms-empty">No recent events.</li>
      </ul>
    </article>
  </section>
</template>

<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue'
import Swal from 'sweetalert2'
import { getSystemResources, getServiceLogs, runServiceAction } from '@/api/operator.js'

const emit = defineEmits(['message', 'error'])

const resources = ref(null)
const dockerUnavailable = ref(false)
const liveResources = ref(true)
const busyService = ref('')
const activeLogService = ref('')
const logLines = ref([])
const liveLogs = ref(false)
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
function formatDateTime(value) {
  return value ? new Date(value).toLocaleString() : ''
}

async function loadResources() {
  try {
    resources.value = await getSystemResources()
    dockerUnavailable.value = (resources.value?.services || []).every(s => s.status === 'unknown')
  } catch (err) { emit('error', err) }
}

async function loadLogs(service) {
  activeLogService.value = service
  try {
    const data = await getServiceLogs(service, 150)
    logLines.value = data.lines || []
  } catch (err) { emit('error', err) }
}

async function serviceAction(service, action) {
  const expected = `${action.toUpperCase()} ${service.toUpperCase()}`
  const result = await Swal.fire({
    title: `${action[0].toUpperCase()}${action.slice(1)} ${service}?`,
    text: `Type ${expected} to confirm.`,
    input: 'text',
    inputPlaceholder: expected,
    showCancelButton: true,
    confirmButtonText: 'Confirm',
    confirmButtonColor: action === 'stop' ? '#fc544b' : '#6777ef',
    cancelButtonColor: '#94a3b8',
    inputValidator: value => value !== expected ? `Type exactly: ${expected}` : undefined,
  })
  if (!result.isConfirmed) return
  busyService.value = service
  try {
    await runServiceAction(service, action, expected)
    emit('message', `${service} ${action} requested`)
    await loadResources()
  } catch (err) {
    emit('error', err)
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
/*
 * This component is rendered as a child of WebsiteContentManagerView, whose
 * <style scoped> only reaches this component's root element (Vue scoped CSS
 * does not cascade into child templates) — so shared classes are redefined
 * here rather than assumed to be inherited.
 */
.ops-panel { display: grid; gap: 20px; min-width: 0; }
.cms-panel-head { display: flex; justify-content: space-between; align-items: center; gap: 16px; flex-wrap: wrap; border-bottom: 1px solid #f9f9f9; padding-bottom: 15px; }
.cms-panel-head h2, .cms-panel-head h3 { margin: 0; font-size: 16px; font-weight: 700; color: #34395e; }
.cms-panel { background: white; border-radius: 3px; box-shadow: 0 4px 25px rgba(0,0,0,.1); padding: 25px; display: grid; gap: 16px; min-width: 0; }
.cms-card.metric { background: white; box-shadow: 0 4px 25px rgba(0,0,0,.1); border-radius: 3px; padding: 20px; display: grid; gap: 6px; }
.cms-card.metric span { color: #98a6ad; font-size: 12px; font-weight: 600; }
.cms-card.metric strong { font-size: 24px; color: #34395e; }
.cms-card.metric small { color: #98a6ad; font-size: 11px; }
.cms-empty { color: #98a6ad; font-style: italic; text-align: center; padding: 16px; }

.ops-live-toggle { display: flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 600; color: #34395e; }
.ops-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr)); gap: 16px; }
.ops-bar { height: 6px; border-radius: 30px; background: #f4f6f9; overflow: hidden; }
.ops-bar span { display: block; height: 100%; background: #6777ef; border-radius: 30px; }

.ops-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.ops-table th { text-align: left; color: #94a3b8; font-size: 11px; text-transform: uppercase; padding: 8px 10px; border-bottom: 1px solid #f1f2fb; }
.ops-table td { padding: 10px; border-bottom: 1px solid #f1f2fb; color: #34395e; vertical-align: middle; }
.ops-badge { border-radius: 30px; padding: 3px 10px; font-size: 11px; font-weight: 700; text-transform: capitalize; background: #eef0fd; color: #6777ef; }
.ops-badge.healthy { background: #e8f7f0; color: #47c363; }
.ops-badge.degraded { background: #fff4e6; color: #ffa426; }
.ops-badge.stopped { background: #fdeaea; color: #fc544b; }
.ops-actions { display: flex; gap: 6px; flex-wrap: wrap; }
.ops-actions button { border: 0; border-radius: 30px; background: #6777ef; color: #fff; padding: 6px 12px; font-size: 11px; font-weight: 700; }
.ops-actions button:disabled { opacity: .5; }
.ops-actions .ops-danger { background: #fc544b; }
.ops-hint { color: #98a6ad; font-size: 12px; margin: 0; }

.ops-terminal { background: #0f172a; color: #d1e7ff; border-radius: 6px; padding: 16px; font-family: 'Fira Code', monospace; font-size: 12px; max-height: 22rem; overflow: auto; white-space: pre-wrap; word-break: break-word; }

.ops-events { list-style: none; margin: 0; padding: 0; display: grid; gap: 8px; }
.ops-events li { display: flex; align-items: center; gap: 10px; font-size: 13px; color: #34395e; }
.ops-events li small { margin-left: auto; color: #98a6ad; font-size: 11px; }
.ops-event-dot { width: 8px; height: 8px; border-radius: 50%; background: #6777ef; flex-shrink: 0; }
.ops-events li.warning .ops-event-dot { background: #ffa426; }
.ops-events li.critical .ops-event-dot { background: #fc544b; }

:global(.dark .cms-panel-head h2), :global(.dark .cms-panel-head h3) { color: #f8fafc; }
:global(.dark .cms-panel), :global(.dark .cms-card.metric) { background: #111827; border-color: #334155; }
:global(.dark .cms-card.metric strong), :global(.dark .ops-table td), :global(.dark .ops-events li), :global(.dark .ops-live-toggle) { color: #f8fafc; }
:global(.dark .ops-table th) { color: #64748b; }
:global(.dark .ops-table td) { border-color: #1e293b; }
</style>
