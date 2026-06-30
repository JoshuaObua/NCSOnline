<template>
  <LayoutDefault title="Database Backup & Restore">
    <div class="space-y-6">
      <section class="bg-white border border-gray-100 rounded-xl p-5 flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <p class="text-xs font-bold uppercase tracking-wider text-primary-700">PostgreSQL recovery centre</p>
          <h1 class="text-xl font-bold text-primary-800 mt-1">Backups, schema tools, and restore jobs</h1>
          <p class="text-sm text-gray-500 mt-1">Create verified recovery points, export schema SQL, and perform guarded recovery operations.</p>
        </div>
        <div class="flex flex-wrap gap-2">
          <button @click="queue('BACKUP')" class="btn-primary"><i class="icofont-database-add"></i>Create backup</button>
          <button @click="downloadSchema" class="btn-light"><i class="icofont-download"></i>Download schema</button>
        </div>
      </section>

      <section v-if="activeJob" class="rounded-xl border border-blue-100 bg-blue-50 p-5">
        <div class="flex justify-between text-sm"><strong>{{ activeJob.job_type }} in progress</strong><span>{{ activeJob.progress }}%</span></div>
        <div class="h-2 bg-white rounded-full mt-3"><div class="h-full bg-blue-600 rounded-full transition-all" :style="{ width: activeJob.progress + '%' }"></div></div>
        <p class="text-xs text-blue-800 mt-2">{{ activeJob.message }}</p>
      </section>

      <section class="grid gap-4 lg:grid-cols-3">
        <article class="tool-card">
          <h2>Schema Import</h2>
          <p>Apply a plain SQL schema file while maintenance mode is active.</p>
          <input ref="schemaInput" type="file" accept=".sql,application/sql,text/plain" class="hidden" @change="importSchema" />
          <button @click="schemaInput?.click()" class="btn-light mt-4"><i class="icofont-upload"></i>Import schema</button>
        </article>
        <article class="tool-card border-red-100">
          <h2 class="text-red-800">Delete Schema</h2>
          <p>Drop and recreate the public schema. Use only before importing a clean schema or restoring data.</p>
          <button @click="deleteSchema" class="btn-danger mt-4"><i class="icofont-warning-alt"></i>Delete schema</button>
        </article>
        <article class="tool-card">
          <h2>Recovery Guardrails</h2>
          <ul class="mt-3 space-y-2 text-sm text-gray-600">
            <li><i class="icofont-check text-green-600 mr-1"></i>Restore requires a verified dump.</li>
            <li><i class="icofont-check text-green-600 mr-1"></i>Restore/import/delete require maintenance mode.</li>
            <li><i class="icofont-check text-green-600 mr-1"></i>Restore creates a pre-restore safety dump.</li>
          </ul>
        </article>
      </section>

      <section class="bg-white border border-gray-100 rounded-xl overflow-hidden">
        <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
          <h2 class="font-bold text-primary-800">Recovery points</h2>
          <button @click="load" class="btn-light"><i class="icofont-refresh"></i>Refresh</button>
        </div>
        <div class="overflow-x-auto">
          <table class="min-w-full text-sm">
            <thead class="bg-gray-50 text-gray-500">
              <tr><th class="text-left p-4">Backup</th><th class="text-left p-4">Created</th><th class="text-left p-4">Size</th><th class="text-left p-4">Integrity</th><th class="text-right p-4">Actions</th></tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="backup in backups" :key="backup.id">
                <td class="p-4 max-w-md">
                  <strong class="text-primary-800 break-all">{{ backup.file_name }}</strong>
                  <p class="font-mono text-[10px] text-gray-400 mt-1 break-all">{{ backup.checksum }}</p>
                </td>
                <td class="p-4 text-gray-500">{{ date(backup.created_at) }}</td>
                <td class="p-4">{{ bytes(backup.size_bytes) }}</td>
                <td class="p-4"><span class="status-pill" :class="backup.status === 'VERIFIED' ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'">{{ backup.status }}</span></td>
                <td class="p-4 text-right">
                  <div class="inline-flex flex-wrap justify-end gap-2">
                    <button @click="downloadBackup(backup.id)" class="btn-icon" title="Download backup"><i class="icofont-download"></i></button>
                    <button @click="queue('VERIFY', backup.id)" class="btn-icon" title="Verify backup"><i class="icofont-check-circled"></i></button>
                    <button :disabled="backup.status !== 'VERIFIED'" @click="restore(backup.id)" class="btn-icon text-red-700 disabled:opacity-30" title="Restore backup"><i class="icofont-history"></i></button>
                    <button @click="deleteBackup(backup.id)" class="btn-icon text-red-700" title="Delete backup"><i class="icofont-trash"></i></button>
                  </div>
                </td>
              </tr>
              <tr v-if="!backups.length"><td colspan="5" class="p-10 text-center text-gray-400">No backup records yet.</td></tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="bg-white border border-gray-100 rounded-xl p-5">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <h2 class="font-bold text-primary-800">Recent jobs</h2>
          <div class="flex flex-wrap gap-2">
            <button @click="exportJobs" class="btn-light"><i class="icofont-download"></i>Export logs</button>
            <button @click="clearJobs" class="btn-danger"><i class="icofont-trash"></i>Clear all</button>
          </div>
        </div>
        <div class="mt-4 grid gap-2">
          <div v-for="job in jobs.slice(0, 12)" :key="job.id" class="flex flex-col gap-2 rounded-lg bg-gray-50 p-3 text-sm md:flex-row md:items-center md:justify-between">
            <span><strong>{{ job.job_type }}</strong> - {{ job.message }}<br><small class="text-gray-400 font-mono">{{ job.id }}</small></span>
            <span class="flex items-center gap-2">
              <span class="font-bold" :class="job.status === 'FAILED' ? 'text-red-600' : job.status === 'SUCCEEDED' ? 'text-green-600' : 'text-blue-600'">{{ job.status }}</span>
              <button :disabled="job.status === 'RUNNING' || job.status === 'PENDING'" @click="deleteJob(job.id)" class="btn-icon text-red-700 disabled:opacity-30" title="Delete job log"><i class="icofont-trash"></i></button>
            </span>
          </div>
          <p v-if="!jobs.length" class="text-sm text-gray-400">No jobs yet.</p>
        </div>
      </section>

      <p v-if="message" class="rounded-lg border px-4 py-3 text-sm" :class="error ? 'border-red-200 bg-red-50 text-red-700' : 'border-emerald-200 bg-emerald-50 text-emerald-700'">{{ message }}</p>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import apiClient, { API_BASE_URL } from '@/api/client.js'
import Swal from 'sweetalert2'

const backups = ref([])
const jobs = ref([])
const schemaInput = ref(null)
const message = ref('')
const error = ref(false)
let timer

const activeJob = computed(() => jobs.value.find(j => j.status === 'RUNNING' || j.status === 'PENDING'))

async function load() {
  const r = await apiClient.get('/api/v1/admin/system/backups')
  backups.value = r.data.data?.backups || []
  jobs.value = r.data.data?.jobs || []
}
async function queue(action, backup_id = null, confirmation = '') {
  try {
    await apiClient.post('/api/v1/admin/system/backups/jobs', { action, backup_id, confirmation })
    await alertOk(`${action} job queued`)
    await load()
  } catch (err) {
    await alertErr(err.response?.data?.error?.message || `${action} failed`)
  }
}
async function restore(id) {
  if (!(await typedConfirm('Restore verified backup', 'RESTORE VERIFIED BACKUP', 'Enable maintenance mode first. This will restore data from the selected backup.'))) return
  await queue('RESTORE', id, 'RESTORE VERIFIED BACKUP')
}
async function downloadSchema() {
  await download('/api/v1/admin/system/backups/schema', `ncs-schema-${Date.now()}.sql`)
}
async function downloadBackup(id) {
  await download(`/api/v1/admin/system/backups/${id}/download`, `backup-${id}.dump`)
}
async function deleteBackup(id) {
  if (!(await yesNo('Delete backup?', 'Delete this backup record and local backup file?'))) return
  try {
    await apiClient.delete(`/api/v1/admin/system/backups/${id}`)
    await alertOk('Backup deleted')
    await load()
  } catch (err) {
    await alertErr(err.response?.data?.error?.message || 'Delete failed')
  }
}
async function importSchema(e) {
  const file = e.target.files?.[0]
  e.target.value = ''
  if (!file) return
  if (!(await typedConfirm('Import schema', 'IMPORT SCHEMA', 'Enable maintenance mode first. This applies the selected SQL schema file.'))) return
  const fd = new FormData()
  fd.append('schema', file)
  fd.append('confirmation', 'IMPORT SCHEMA')
  try {
    await apiClient.post('/api/v1/admin/system/backups/schema/import', fd, { headers: { 'Content-Type': 'multipart/form-data' } })
    await alertOk('Schema imported')
  } catch (err) {
    await alertErr(err.response?.data?.error?.message || 'Schema import failed')
  }
}
async function deleteSchema() {
  if (!(await typedConfirm('Delete database schema', 'DELETE DATABASE SCHEMA', 'Enable maintenance mode first. This drops and recreates the public schema.'))) return
  try {
    await apiClient.delete('/api/v1/admin/system/backups/schema', { data: { confirmation: 'DELETE DATABASE SCHEMA' } })
    await alertOk('Database schema deleted and recreated')
  } catch (err) {
    await alertErr(err.response?.data?.error?.message || 'Schema delete failed')
  }
}
async function exportJobs() {
  await download('/api/v1/admin/system/backups/jobs/export', `backup-jobs-${Date.now()}.csv`)
}
async function clearJobs() {
  if (!(await typedConfirm('Clear backup job logs', 'CLEAR BACKUP JOB LOGS', 'Only finished job logs are cleared. Running and pending jobs are preserved.'))) return
  try {
    await apiClient.delete('/api/v1/admin/system/backups/jobs', { data: { confirmation: 'CLEAR BACKUP JOB LOGS' } })
    await alertOk('Finished job logs cleared')
    await load()
  } catch (err) {
    await alertErr(err.response?.data?.error?.message || 'Could not clear job logs')
  }
}
async function deleteJob(id) {
  if (!(await yesNo('Delete job log?', 'Delete this backup job log entry?'))) return
  try {
    await apiClient.delete(`/api/v1/admin/system/backups/jobs/${id}`)
    await alertOk('Job log deleted')
    await load()
  } catch (err) {
    await alertErr(err.response?.data?.error?.message || 'Could not delete job log')
  }
}
async function download(path, fallbackName) {
  const token = localStorage.getItem('ncsms_access_token')
  const res = await fetch(`${API_BASE_URL}${path}`, { headers: token ? { Authorization: `Bearer ${token}` } : {} })
  if (!res.ok) {
    await alertErr('Download failed')
    return
  }
  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filenameFromDisposition(res.headers.get('content-disposition')) || fallbackName
  a.click()
  URL.revokeObjectURL(url)
}
function filenameFromDisposition(value) {
  return value?.match(/filename="?([^"]+)"?/)?.[1]
}
function show(text, isError = false) {
  message.value = text
  error.value = isError
}
async function typedConfirm(title, confirmText, text) {
  const r = await Swal.fire({
    title, text, input: 'text', inputPlaceholder: confirmText,
    icon: 'warning', showCancelButton: true, confirmButtonText: 'Confirm',
    preConfirm: v => v === confirmText || Swal.showValidationMessage(`Type ${confirmText} to continue`),
  })
  return r.isConfirmed
}
async function yesNo(title, text) {
  const r = await Swal.fire({ title, text, icon: 'warning', showCancelButton: true, confirmButtonText: 'Yes' })
  return r.isConfirmed
}
async function alertOk(text) {
  show(text)
  await Swal.fire({ toast: true, position: 'top-end', timer: 2200, showConfirmButton: false, icon: 'success', title: text })
}
async function alertErr(text) {
  show(text, true)
  await Swal.fire({ icon: 'error', title: 'Action failed', text })
}
const date = v => new Date(v).toLocaleString()
const bytes = v => { let n = v || 0, i = 0, u = ['B','KB','MB','GB','TB']; while (n >= 1024 && i < u.length - 1) { n /= 1024; i++ } return `${n.toFixed(1)} ${u[i]}` }

onMounted(() => { load(); timer = setInterval(load, 3000) })
onBeforeUnmount(() => clearInterval(timer))
</script>

<style scoped>
.btn-primary { display: inline-flex; align-items: center; gap: 0.5rem; border-radius: 0.75rem; background: #1f4f82; padding: 0.75rem 1rem; color: white; font-weight: 700; }
.btn-light { display: inline-flex; align-items: center; gap: 0.5rem; border-radius: 0.75rem; border: 1px solid #e5e7eb; background: white; padding: 0.65rem 0.9rem; color: #374151; font-weight: 700; }
.btn-danger { display: inline-flex; align-items: center; gap: 0.5rem; border-radius: 0.75rem; background: #b91c1c; padding: 0.65rem 0.9rem; color: white; font-weight: 700; }
.btn-icon { display: inline-flex; align-items: center; justify-content: center; width: 2.25rem; height: 2.25rem; border-radius: 0.5rem; border: 1px solid #e5e7eb; background: white; }
.tool-card { border: 1px solid #f3f4f6; background: white; border-radius: 0.75rem; padding: 1.25rem; }
.tool-card h2 { font-weight: 700; color: #1f2937; }
.tool-card p { margin-top: 0.35rem; font-size: 0.875rem; color: #6b7280; }
.status-pill { border-radius: 999px; padding: 0.25rem 0.65rem; font-size: 0.75rem; font-weight: 700; }
</style>
