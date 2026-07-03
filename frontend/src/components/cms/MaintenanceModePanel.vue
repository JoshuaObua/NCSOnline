<template>
  <section class="maint-panel">
    <div class="cms-panel-head"><h2>Maintenance Mode &amp; Database Tooling</h2></div>

    <article v-for="scope in scopes" :key="scope.key" class="cms-panel">
      <div class="cms-panel-head">
        <h3>{{ scope.label }}</h3>
        <span class="ops-badge" :class="{ healthy: !state[scope.key].enabled, stopped: state[scope.key].enabled }">
          {{ state[scope.key].enabled ? 'Active' : 'Off' }}
        </span>
      </div>
      <div class="cms-two">
        <label class="check"><input v-model="state[scope.key].enabled" type="checkbox" /> Enable maintenance mode</label>
        <label class="wide">Reason <small>(required to enable, min 5 characters)</small><input v-model="state[scope.key].reason" class="form-control" placeholder="e.g. Scheduled database migration" /></label>
        <label>Custom title<input v-model="state[scope.key].display_meta.custom_title" class="form-control" /></label>
        <label>Custom message<input v-model="state[scope.key].display_meta.custom_message" class="form-control" /></label>
        <label>Bypass roles <small>(comma-separated)</small><input :value="(state[scope.key].bypass_rules.allowed_roles || []).join(', ')" @input="setList(scope.key, 'allowed_roles', $event.target.value)" class="form-control" /></label>
        <label>Bypass IP ranges <small>(comma-separated CIDRs)</small><input :value="(state[scope.key].bypass_rules.allowed_ip_ranges || []).join(', ')" @input="setList(scope.key, 'allowed_ip_ranges', $event.target.value)" class="form-control" /></label>
      </div>
      <div class="profile-actions">
        <button type="button" :disabled="saving === scope.key" @click="saveScope(scope.key)">{{ saving === scope.key ? 'Saving…' : 'Save' }}</button>
      </div>
    </article>

    <article class="cms-panel">
      <div class="cms-panel-head"><h3>Maintenance Actions</h3></div>
      <div class="profile-actions">
        <button type="button" @click="runAction('flush_app_cache')">Flush App Cache</button>
        <button type="button" @click="runAction('flush_nginx_cache')">Flush Nginx Cache</button>
        <button type="button" @click="runAction('optimize_database_tables')">Optimize Database Tables</button>
        <button type="button" @click="runAction('prune_activity_logs')">Prune Activity Logs</button>
      </div>
    </article>

    <article class="cms-panel">
      <div class="cms-panel-head">
        <h3>Database Backups</h3>
        <div class="profile-actions">
          <button type="button" :disabled="queuing" @click="queue('BACKUP')">{{ queuing ? 'Queuing…' : 'Queue backup' }}</button>
          <button type="button" @click="checkHealth">Check DB health</button>
        </div>
      </div>
      <p v-if="health" class="ops-hint">{{ health.status === 'healthy' ? 'Healthy' : health.status }} — {{ healthSummary }}</p>
      <table class="ops-table">
        <thead><tr><th>File</th><th>Size</th><th>Status</th><th>Created</th><th></th></tr></thead>
        <tbody>
          <tr v-for="b in backups" :key="b.id">
            <td>{{ b.file_name }}</td>
            <td>{{ formatBytes(b.size_bytes) }}</td>
            <td><span class="ops-badge" :class="{ healthy: b.status === 'VERIFIED', stopped: b.status === 'FAILED' }">{{ b.status }}</span></td>
            <td>{{ formatDateTime(b.created_at) }}</td>
            <td class="ops-actions">
              <button type="button" @click="queue('VERIFY', b.id)">Verify</button>
              <button type="button" :disabled="b.status !== 'VERIFIED'" :title="b.status !== 'VERIFIED' ? 'Only verified backups can be restored' : ''" @click="restore(b.id)">Restore</button>
              <button type="button" @click="download(b.id, b.file_name)">Download</button>
              <button type="button" class="ops-danger" @click="removeBackup(b.id)">Delete</button>
            </td>
          </tr>
          <tr v-if="!backups.length"><td colspan="5" class="cms-empty">No backups yet.</td></tr>
        </tbody>
      </table>
    </article>

    <article class="cms-panel">
      <div class="cms-panel-head">
        <h3>Backup Job Log</h3>
        <div class="profile-actions">
          <button type="button" @click="exportJobs">Export CSV</button>
          <button type="button" class="ops-danger" @click="clearJobLog">Clear finished jobs</button>
        </div>
      </div>
      <table class="ops-table">
        <thead><tr><th>Type</th><th>Status</th><th>Progress</th><th>Message</th><th>Created</th></tr></thead>
        <tbody>
          <tr v-for="j in jobs" :key="j.id">
            <td>{{ j.job_type }}</td>
            <td><span class="ops-badge" :class="{ healthy: j.status === 'SUCCEEDED', stopped: j.status === 'FAILED' }">{{ j.status }}</span></td>
            <td>{{ j.progress }}%</td>
            <td>{{ j.message }}</td>
            <td>{{ formatDateTime(j.created_at) }}</td>
          </tr>
          <tr v-if="!jobs.length"><td colspan="5" class="cms-empty">No job history yet.</td></tr>
        </tbody>
      </table>
    </article>

    <article class="cms-panel">
      <div class="cms-panel-head"><h3>Schema Management</h3></div>
      <p class="ops-hint">Schema import/delete require maintenance mode to be enabled first, and are destructive — use with care.</p>
      <div class="profile-actions">
        <button type="button" @click="exportSchema">Download schema</button>
        <label class="ops-file-btn">Import schema
          <input ref="schemaInput" type="file" accept=".sql" @change="onSchemaPicked" />
        </label>
        <button type="button" class="ops-danger" @click="dropSchema">Delete schema</button>
      </div>
    </article>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import Swal from 'sweetalert2'
import {
  clearBackupJobs, deleteBackup, deleteSchema, downloadBackup, downloadSchema,
  exportBackupJobs, getBackupHealth, getSystemStatus, importSchema, listBackups, queueBackupJob,
  runMaintenanceAction, setMaintenance,
} from '@/api/operator.js'

const emit = defineEmits(['message', 'error'])

const scopes = [
  { key: 'public_cms', label: 'Public Website Maintenance' },
  { key: 'admin_dashboard', label: 'Admin Dashboard Maintenance' },
]
function emptyScope() {
  return { enabled: false, reason: '', display_meta: { custom_title: '', custom_message: '' }, bypass_rules: { allowed_roles: [], allowed_ip_ranges: [] } }
}
const state = reactive({ public_cms: emptyScope(), admin_dashboard: emptyScope() })
const saving = ref('')

const backups = ref([])
const jobs = ref([])
const health = ref(null)
const queuing = ref(false)
const schemaInput = ref(null)

const healthSummary = computed(() => {
  const r = health.value?.result
  if (!r) return ''
  try {
    const parsed = typeof r === 'string' ? JSON.parse(r) : r
    return `${parsed.tables} tables, ${parsed.indexes} indexes, ${parsed.dead_tuples} dead tuples`
  } catch { return '' }
})

function setList(scopeKey, field, raw) {
  state[scopeKey].bypass_rules[field] = raw.split(',').map(s => s.trim()).filter(Boolean)
}

function formatBytes(bytes) {
  const n = Number(bytes) || 0
  if (n <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const i = Math.min(units.length - 1, Math.floor(Math.log(n) / Math.log(1024)))
  return `${(n / Math.pow(1024, i)).toFixed(1)} ${units[i]}`
}
function formatDateTime(value) { return value ? new Date(value).toLocaleString() : '' }
function downloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  URL.revokeObjectURL(url)
}

async function loadMaintenance() {
  try {
    const status = await getSystemStatus()
    for (const scope of scopes) {
      const s = status?.[scope.key]
      if (!s) continue
      state[scope.key] = {
        enabled: !!s.maintenance_mode,
        reason: s.reason || '',
        display_meta: { custom_title: s.display_meta?.custom_title || '', custom_message: s.display_meta?.custom_message || '' },
        bypass_rules: { allowed_roles: s.bypass_rules?.allowed_roles || [], allowed_ip_ranges: s.bypass_rules?.allowed_ip_ranges || [] },
      }
    }
  } catch (err) { emit('error', err) }
}

async function saveScope(scopeKey) {
  const s = state[scopeKey]
  if (s.enabled && s.reason.trim().length < 5) {
    emit('error', new Error('A maintenance reason of at least 5 characters is required to enable maintenance mode'))
    return
  }
  saving.value = scopeKey
  try {
    await setMaintenance({
      scope: scopeKey,
      enabled: s.enabled,
      reason: s.reason,
      display_meta: s.display_meta,
      bypass_rules: s.bypass_rules,
    })
    emit('message', 'Maintenance settings saved')
    await loadMaintenance()
  } catch (err) {
    emit('error', err)
  } finally {
    saving.value = ''
  }
}

async function confirmTyped(title, expected, danger = false) {
  const result = await Swal.fire({
    title,
    text: `Type ${expected} to confirm.`,
    input: 'text',
    inputPlaceholder: expected,
    showCancelButton: true,
    confirmButtonText: 'Confirm',
    confirmButtonColor: danger ? '#fc544b' : '#6777ef',
    cancelButtonColor: '#94a3b8',
    inputValidator: value => value !== expected ? `Type exactly: ${expected}` : undefined,
  })
  return result.isConfirmed
}

async function runAction(action) {
  const expected = action.toUpperCase().replaceAll('_', ' ')
  if (!(await confirmTyped(`Run "${expected}"?`, expected)) ) return
  try {
    const res = await runMaintenanceAction(action, expected)
    emit('message', res?.message || `${expected} completed`)
  } catch (err) { emit('error', err) }
}

async function loadBackups() {
  try {
    const data = await listBackups()
    backups.value = data?.backups || []
    jobs.value = data?.jobs || []
  } catch (err) { emit('error', err) }
}

async function checkHealth() {
  try {
    health.value = await getBackupHealth()
  } catch (err) { emit('error', err) }
}

async function queue(action, backupId) {
  queuing.value = action === 'BACKUP'
  try {
    await queueBackupJob(action, backupId)
    emit('message', `${action} job queued`)
    await loadBackups()
  } catch (err) {
    emit('error', err)
  } finally {
    queuing.value = false
  }
}

async function restore(backupId) {
  if (!(await confirmTyped('Restore this backup?', 'CONFIRM RESTORE', true))) return
  try {
    await queueBackupJob('RESTORE', backupId, 'CONFIRM RESTORE')
    emit('message', 'Restore job queued')
    await loadBackups()
  } catch (err) { emit('error', err) }
}

async function download(id, filename) {
  try {
    const res = await downloadBackup(id)
    downloadBlob(res.data, filename)
  } catch (err) { emit('error', err) }
}

async function removeBackup(id) {
  if (!(await confirmTyped('Delete this backup?', 'DELETE BACKUP', true))) return
  try {
    await deleteBackup(id)
    emit('message', 'Backup deleted')
    await loadBackups()
  } catch (err) { emit('error', err) }
}

async function exportJobs() {
  try {
    const res = await exportBackupJobs()
    downloadBlob(res.data, `backup-jobs-${Date.now()}.csv`)
  } catch (err) { emit('error', err) }
}

async function clearJobLog() {
  if (!(await confirmTyped('Clear finished backup job logs?', 'CLEAR BACKUP JOB LOGS', true))) return
  try {
    await clearBackupJobs('CLEAR BACKUP JOB LOGS')
    emit('message', 'Job logs cleared')
    await loadBackups()
  } catch (err) { emit('error', err) }
}
async function exportSchema() {
  try {
    const res = await downloadSchema()
    downloadBlob(res.data, `ncs-schema-${Date.now()}.sql`)
  } catch (err) { emit('error', err) }
}

async function onSchemaPicked(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) return
  if (!(await confirmTyped(`Import ${file.name}?`, 'IMPORT SCHEMA', true))) return
  try {
    await importSchema(file, 'IMPORT SCHEMA')
    emit('message', 'Schema imported')
  } catch (err) { emit('error', err) }
}

async function dropSchema() {
  if (!(await confirmTyped('Delete the entire database schema?', 'DELETE DATABASE SCHEMA', true))) return
  try {
    await deleteSchema('DELETE DATABASE SCHEMA')
    emit('message', 'Schema deleted and recreated')
  } catch (err) { emit('error', err) }
}

onMounted(() => {
  loadMaintenance()
  loadBackups()
})
</script>

<style scoped>
/*
 * This component is rendered as a child of WebsiteContentManagerView, whose
 * <style scoped> only reaches this component's root element (Vue scoped CSS
 * does not cascade into child templates) — so shared classes are redefined
 * here rather than assumed to be inherited.
 */
.maint-panel { display: grid; gap: 20px; min-width: 0; }
.cms-panel-head { display: flex; justify-content: space-between; align-items: center; gap: 16px; flex-wrap: wrap; border-bottom: 1px solid #f9f9f9; padding-bottom: 15px; }
.cms-panel-head h2, .cms-panel-head h3 { margin: 0; font-size: 16px; font-weight: 700; color: #34395e; }
.cms-panel { background: white; border-radius: 3px; box-shadow: 0 4px 25px rgba(0,0,0,.1); padding: 25px; display: grid; gap: 16px; min-width: 0; }
.cms-empty { color: #98a6ad; font-style: italic; text-align: center; padding: 16px; }

.cms-two { display: grid; grid-template-columns: repeat(auto-fit, minmax(15rem, 1fr)); gap: 18px; }
.cms-two label { display: grid; gap: 7px; font-size: 12px; font-weight: 600; color: #34395e; min-width: 0; }
.cms-two label small { font-weight: 500; text-transform: none; color: #94a3b8; }
.cms-two .wide { grid-column: 1 / -1; }
.cms-two .check { display: flex; flex-direction: row; align-items: center; gap: .5rem; }
.cms-two input { border: 1px solid #e4e6fc; border-radius: 3px; padding: 10px 15px; font-weight: 500; color: #495057; background: #fdfdff; outline: none; width: 100%; box-sizing: border-box; }
.cms-two input[type="checkbox"] { width: 18px; height: 18px; accent-color: #6777ef; }

.profile-actions { display: flex; flex-wrap: wrap; gap: 10px; }
.profile-actions button { border: 0; border-radius: 30px; background: #6777ef; color: white; padding: 8px 18px; font-size: 12px; font-weight: 600; white-space: nowrap; box-shadow: 0 2px 6px #acb5f6; }
.profile-actions button:disabled { opacity: .6; }

.ops-badge { border-radius: 30px; padding: 3px 10px; font-size: 11px; font-weight: 700; background: #eef0fd; color: #6777ef; }
.ops-badge.healthy { background: #e8f7f0; color: #47c363; }
.ops-badge.stopped { background: #fdeaea; color: #fc544b; }
.ops-hint { color: #98a6ad; font-size: 12px; margin: 0; }

.ops-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.ops-table th { text-align: left; color: #94a3b8; font-size: 11px; text-transform: uppercase; padding: 8px 10px; border-bottom: 1px solid #f1f2fb; }
.ops-table td { padding: 10px; border-bottom: 1px solid #f1f2fb; color: #34395e; vertical-align: middle; }
.ops-actions { display: flex; gap: 6px; flex-wrap: wrap; }
.ops-actions button { border: 0; border-radius: 30px; background: #6777ef; color: #fff; padding: 6px 12px; font-size: 11px; font-weight: 700; }
.ops-actions button:disabled { opacity: .4; }
.ops-danger { background: #fc544b !important; }

.ops-file-btn { position: relative; border: 1px solid #e4e6fc; border-radius: 30px; background: #fff; color: #34395e; padding: 8px 18px; font-size: 12px; font-weight: 700; cursor: pointer; }
.ops-file-btn input { position: absolute; inset: 0; opacity: 0; cursor: pointer; }

:global(.dark .cms-panel-head h2), :global(.dark .cms-panel-head h3) { color: #f8fafc; }
:global(.dark .cms-panel) { background: #111827; border-color: #334155; }
:global(.dark .cms-two input) { background: #0f172a; color: #f8fafc; border-color: #475569; }
:global(.dark .ops-table td) { color: #f8fafc; border-color: #1e293b; }
:global(.dark .ops-table th) { color: #64748b; }
:global(.dark .ops-file-btn) { background: #0f172a; color: #f8fafc; border-color: #475569; }
</style>
