<template>
  <LayoutDefault title="System Operations">
    <div class="space-y-6 p-6 max-w-7xl mx-auto">
      <section class="grid xl:grid-cols-2 gap-4">
        <article v-for="scope in scopes" :key="scope.key" class="admin-card overflow-hidden">
          <div class="admin-card-header flex-wrap gap-3">
            <div>
              <p class="text-xs font-bold uppercase tracking-widest" :class="scopeLabelClass(scope.state)">{{ scopeLabel(scope.state) }}</p>
              <h2 class="text-xl font-bold text-primary-800 mt-1">{{ scope.title }}</h2>
              <p class="text-sm text-gray-500 mt-1">{{ scope.description }}</p>
            </div>
            <div v-if="authStore.isSuperAdmin" class="flex flex-wrap gap-2">
              <button @click="openScheduler(scope.key)" class="rounded-xl px-4 py-2.5 text-sm font-bold text-white" :class="scope.state?.is_active ? 'bg-emerald-700 hover:bg-emerald-800' : 'bg-amber-600 hover:bg-amber-700'">
                {{ scope.state?.is_active ? 'End now' : (scope.state?.is_scheduled ? 'Edit schedule' : 'Schedule') }}
              </button>
              <button v-if="scope.state?.is_scheduled || (scope.state?.maintenance_mode && !scope.state?.is_active)" @click="cancelSchedule(scope.key)" class="rounded-xl px-3 py-2.5 text-sm font-semibold border border-gray-200 text-gray-700 hover:bg-gray-50">
                Cancel
              </button>
            </div>
          </div>
          <div class="admin-card-body space-y-4">
            <div class="rounded-xl border p-4" :class="scopePanelClass(scope.state)">
              <h3 class="font-bold text-gray-900">{{ scope.state?.display_meta?.custom_title || scope.defaultTitle }}</h3>
              <p class="text-sm text-gray-600 mt-1">{{ scope.state?.display_meta?.custom_message || scope.state?.reason || 'No maintenance window is active.' }}</p>
              <dl v-if="scope.state?.scheduled_start || scope.state?.expected_end" class="mt-3 grid sm:grid-cols-2 gap-2 text-xs text-gray-700">
                <div v-if="scope.state?.scheduled_start"><dt class="font-semibold text-gray-500">Starts</dt><dd>{{ formatDate(scope.state.scheduled_start) }}<br><span class="text-gray-400">{{ formatRel(scope.state.scheduled_start) }}</span></dd></div>
                <div v-if="scope.state?.expected_end"><dt class="font-semibold text-gray-500">Ends</dt><dd>{{ formatDate(scope.state.expected_end) }}<br><span class="text-gray-400">{{ formatRel(scope.state.expected_end) }}</span></dd></div>
              </dl>
            </div>
            <div class="grid md:grid-cols-2 gap-3 text-sm">
              <div class="rounded-xl border border-gray-100 p-3">
                <p class="text-xs uppercase tracking-wide text-gray-400 font-semibold">Bypass roles</p>
                <p class="font-semibold text-gray-800 mt-1">{{ listText(scope.state?.bypass_rules?.allowed_roles) }}</p>
              </div>
              <div class="rounded-xl border border-gray-100 p-3">
                <p class="text-xs uppercase tracking-wide text-gray-400 font-semibold">Bypass IP ranges</p>
                <p class="font-semibold text-gray-800 mt-1">{{ listText(scope.state?.bypass_rules?.allowed_ip_ranges) }}</p>
              </div>
            </div>
            <p class="text-xs text-gray-500">Secret bypass: <span class="font-mono">{{ scope.state?.bypass_rules?.secret_query_param || 'not set' }}</span></p>
          </div>
        </article>
      </section>

      <section>
        <div class="flex items-center justify-between mb-3">
          <div>
            <h2 class="font-bold text-primary-800">Live resources</h2>
            <p class="text-xs text-gray-500">Updated every five seconds from the API container host.</p>
          </div>
          <span class="text-xs text-gray-400">{{ lastResourceUpdate }}</span>
        </div>
        <div class="grid sm:grid-cols-2 xl:grid-cols-4 gap-4">
          <article v-for="card in resourceCards" :key="card.label" class="admin-card p-5">
            <p class="text-xs uppercase tracking-wide text-gray-400 font-semibold">{{ card.label }}</p>
            <p class="text-2xl font-bold text-primary-800 mt-2">{{ card.value }}</p>
            <div class="h-1.5 rounded-full bg-gray-100 mt-4">
              <div class="h-full rounded-full bg-primary-600" :style="{ width: Math.min(card.percent || 0, 100) + '%' }"/>
            </div>
          </article>
        </div>
      </section>

      <section class="grid md:grid-cols-2 gap-4">
        <article class="admin-card p-6">
          <h3 class="font-bold text-primary-800">Application cache</h3>
          <p class="text-sm text-gray-500 mt-2">Invalidate memory-cached settings and force resolvers to reload.</p>
          <p class="text-xs text-gray-400 mt-3">Generation {{ status.cache_generation || 0 }}</p>
          <button v-if="authStore.isSuperAdmin" @click="flushCache" class="mt-4 rounded-xl bg-primary-700 hover:bg-primary-800 text-white px-4 py-2.5 font-semibold">Flush cache</button>
        </article>
        <article class="admin-card border-red-100 p-6">
          <h3 class="font-bold text-red-800">Emergency session revocation</h3>
          <p class="text-sm text-gray-500 mt-2">Revoke every refresh token and reject all currently issued access tokens.</p>
          <button v-if="authStore.isSuperAdmin" @click="revokeAll" class="mt-4 rounded-xl bg-red-700 hover:bg-red-800 text-white px-4 py-2.5 font-semibold">Force logout all users</button>
        </article>
      </section>
    </div>

    <Modal :show="modalOpen" :title="modalTitle" size="xl" @close="modalOpen = false">
      <div class="space-y-5">
        <div class="rounded-xl bg-gray-50 border border-gray-100 p-4">
          <p class="text-xs font-bold uppercase tracking-widest text-gray-500">{{ activeScopeMeta.title }}</p>
          <p class="text-sm text-gray-600 mt-1">{{ activeScopeMeta.description }}</p>
        </div>

        <label class="block text-sm font-semibold">Reason <span class="text-red-500">*</span>
          <textarea v-model="form.reason" rows="2" class="mt-1 w-full rounded-xl border-gray-200 text-sm" placeholder="Describe the approved maintenance work"></textarea>
        </label>

        <div v-if="!activeScope?.is_active" class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <label class="block text-sm font-semibold">Start
            <input v-model="form.startStr" type="datetime-local" class="mt-1 w-full rounded-xl border-gray-200 text-sm"/>
            <span class="block text-[10px] text-gray-400 mt-1">Leave blank to start immediately.</span>
          </label>
          <label class="block text-sm font-semibold">End <span class="text-red-500">*</span>
            <input v-model="form.endStr" type="datetime-local" class="mt-1 w-full rounded-xl border-gray-200 text-sm"/>
            <span class="block text-[10px] text-gray-400 mt-1">Traffic resumes automatically when this passes.</span>
          </label>
        </div>

        <div v-if="!activeScope?.is_active" class="flex flex-wrap gap-2 text-xs">
          <button type="button" @click="presetWindow(15)" class="px-2.5 py-1 border border-gray-200 rounded-full hover:bg-gray-50">+15 min</button>
          <button type="button" @click="presetWindow(60)" class="px-2.5 py-1 border border-gray-200 rounded-full hover:bg-gray-50">+1 hr</button>
          <button type="button" @click="presetWindow(180)" class="px-2.5 py-1 border border-gray-200 rounded-full hover:bg-gray-50">+3 hr</button>
          <button type="button" @click="presetTonight" class="px-2.5 py-1 border border-gray-200 rounded-full hover:bg-gray-50">Tonight 22:00 to 04:00</button>
        </div>

        <div class="grid lg:grid-cols-2 gap-4">
          <div class="space-y-3">
            <h3 class="font-bold text-primary-800">Display</h3>
            <label class="block text-sm font-semibold">Title
              <input v-model="form.customTitle" class="mt-1 w-full rounded-xl border-gray-200 text-sm"/>
            </label>
            <label class="block text-sm font-semibold">Message
              <textarea v-model="form.customMessage" rows="3" class="mt-1 w-full rounded-xl border-gray-200 text-sm"></textarea>
            </label>
            <label class="inline-flex items-center gap-2 text-sm font-semibold">
              <input v-model="form.showCountdown" type="checkbox" class="rounded border-gray-300">
              Show countdown
            </label>
            <label class="inline-flex items-center gap-2 text-sm font-semibold ml-4">
              <input v-model="form.allowEmailNotification" type="checkbox" class="rounded border-gray-300">
              Allow email notification
            </label>
          </div>
          <div class="space-y-3">
            <h3 class="font-bold text-primary-800">Bypass</h3>
            <label class="block text-sm font-semibold">Allowed roles
              <input v-model="form.allowedRoles" class="mt-1 w-full rounded-xl border-gray-200 text-sm" placeholder="super_admin, admin, content_manager"/>
            </label>
            <label class="block text-sm font-semibold">Allowed IP ranges
              <input v-model="form.allowedIPRanges" class="mt-1 w-full rounded-xl border-gray-200 text-sm" placeholder="192.168.1.0/24, 10.0.0.10"/>
            </label>
            <label class="block text-sm font-semibold">Secret query token
              <input v-model="form.secretQueryParam" class="mt-1 w-full rounded-xl border-gray-200 text-sm font-mono"/>
              <span class="block text-[10px] text-gray-400 mt-1">Use as ?maintenance_bypass=token or ?token on preview URLs.</span>
            </label>
          </div>
        </div>

        <p v-if="formError" class="text-xs text-red-600">{{ formError }}</p>
      </div>

      <template #footer>
        <button @click="modalOpen = false" class="px-4 py-2 text-sm">Cancel</button>
        <button @click="saveSchedule" :disabled="saving" class="rounded-xl bg-primary-700 hover:bg-primary-800 disabled:opacity-60 px-5 py-2 text-sm text-white font-bold">
          {{ saving ? 'Saving...' : (activeScope?.is_active ? 'End now' : (isFutureStart ? 'Schedule window' : 'Start now')) }}
        </button>
      </template>
    </Modal>
  </LayoutDefault>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import Modal from '@/components/ui/Modal.vue'
import apiClient from '@/api/client.js'
import { useAuthStore } from '@/stores/auth.js'

const authStore = useAuthStore()
const status = ref({})
const resources = ref({})
const activeScopeKey = ref('public_cms')
const modalOpen = ref(false)
const saving = ref(false)
const formError = ref('')
const lastResourceUpdate = ref('Not loaded')
let timer

const form = reactive({
  reason: '', startStr: '', endStr: '', customTitle: '', customMessage: '',
  showCountdown: true, allowEmailNotification: false, allowedRoles: '',
  allowedIPRanges: '', secretQueryParam: '', autoStartEnforced: true, autoEndEnforced: true,
})

const scopeMeta = {
  public_cms: { key: 'public_cms', title: 'Public CMS', defaultTitle: "We'll be right back", description: 'Controls the public website, public CMS APIs, assets, and visitor-facing render paths.' },
  admin_dashboard: { key: 'admin_dashboard', title: 'Admin Dashboard', defaultTitle: 'Admin Dashboard Maintenance', description: 'Controls back-office dashboards while preserving super-admin break-glass access.' },
}

const scopes = computed(() => [
  { ...scopeMeta.public_cms, state: status.value.public_cms || status.value },
  { ...scopeMeta.admin_dashboard, state: status.value.admin_dashboard || {} },
])
const activeScope = computed(() => activeScopeKey.value === 'admin_dashboard' ? (status.value.admin_dashboard || {}) : (status.value.public_cms || status.value || {}))
const activeScopeMeta = computed(() => scopeMeta[activeScopeKey.value] || scopeMeta.public_cms)
const modalTitle = computed(() => `${activeScope.value?.is_active ? 'End' : 'Schedule'} ${activeScopeMeta.value.title} maintenance`)
const isFutureStart = computed(() => form.startStr && new Date(form.startStr).getTime() > Date.now() + 30_000)

const pct = v => Number.isFinite(v) ? v : 0
const fmtBytes = v => {
  if (!v) return '0 B'
  const units = ['B','KB','MB','GB','TB']
  let i = 0
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++ }
  return `${v.toFixed(1)} ${units[i]}`
}
const resourceCards = computed(() => [
  { label: 'CPU', value: `${pct(resources.value.cpu_percent?.[0]).toFixed(1)}%`, percent: pct(resources.value.cpu_percent?.[0]) },
  { label: 'Memory', value: `${pct(resources.value.memory?.usedPercent).toFixed(1)}%`, percent: pct(resources.value.memory?.usedPercent) },
  { label: 'Disk', value: `${pct(resources.value.disk?.usedPercent).toFixed(1)}%`, percent: pct(resources.value.disk?.usedPercent) },
  { label: 'Heap', value: fmtBytes(resources.value.heap_alloc_bytes), percent: 0 },
])

function scopePanelClass(s) {
  if (s?.is_active) return 'border-amber-300 bg-amber-50'
  if (s?.is_scheduled) return 'border-blue-300 bg-blue-50'
  return 'border-emerald-200 bg-emerald-50'
}
function scopeLabelClass(s) {
  if (s?.is_active) return 'text-amber-700'
  if (s?.is_scheduled) return 'text-blue-700'
  return 'text-emerald-700'
}
function scopeLabel(s) {
  if (s?.is_active) return 'Maintenance active'
  if (s?.is_scheduled) return 'Maintenance scheduled'
  return 'Operational'
}
function listText(items) { return Array.isArray(items) && items.length ? items.join(', ') : 'None configured' }
function toCSV(items) { return Array.isArray(items) ? items.join(', ') : '' }
function fromCSV(text) { return String(text || '').split(',').map(v => v.trim()).filter(Boolean) }
function toLocalInput(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}
function formatDate(iso) { return new Date(iso).toLocaleString() }
function formatRel(iso) {
  const diff = (new Date(iso).getTime() - Date.now()) / 1000
  const abs = Math.abs(diff)
  const tense = diff < 0 ? 'ago' : 'from now'
  if (abs < 60) return `${Math.round(abs)}s ${tense}`
  if (abs < 3600) return `${Math.round(abs / 60)}m ${tense}`
  if (abs < 86400) return `${Math.round(abs / 3600)}h ${tense}`
  return `${Math.round(abs / 86400)}d ${tense}`
}

function openScheduler(scopeKey) {
  activeScopeKey.value = scopeKey
  const s = scopeKey === 'admin_dashboard' ? (status.value.admin_dashboard || {}) : (status.value.public_cms || status.value || {})
  formError.value = ''
  form.reason = s.reason || ''
  form.startStr = toLocalInput(s.scheduled_start)
  form.endStr = toLocalInput(s.expected_end)
  form.customTitle = s.display_meta?.custom_title || scopeMeta[scopeKey].defaultTitle
  form.customMessage = s.display_meta?.custom_message || s.reason || ''
  form.showCountdown = s.display_meta?.show_countdown ?? true
  form.allowEmailNotification = s.display_meta?.allow_email_notification ?? false
  form.allowedRoles = toCSV(s.bypass_rules?.allowed_roles)
  form.allowedIPRanges = toCSV(s.bypass_rules?.allowed_ip_ranges)
  form.secretQueryParam = s.bypass_rules?.secret_query_param || (scopeKey === 'admin_dashboard' ? 'admin_maintenance_bypass' : 'public_maintenance_bypass')
  form.autoStartEnforced = s.auto_start_enforced ?? true
  form.autoEndEnforced = s.auto_end_enforced ?? true
  modalOpen.value = true
}
function presetWindow(minutes) {
  const start = new Date()
  const end = new Date(start.getTime() + minutes * 60_000)
  form.startStr = toLocalInput(start.toISOString())
  form.endStr = toLocalInput(end.toISOString())
}
function presetTonight() {
  const start = new Date(); start.setHours(22, 0, 0, 0)
  const end = new Date(start.getTime() + 6 * 3600_000)
  form.startStr = toLocalInput(start.toISOString())
  form.endStr = toLocalInput(end.toISOString())
}

async function refresh() {
  const [s, r] = await Promise.all([
    apiClient.get('/api/v1/admin/system/status'),
    apiClient.get('/api/v1/admin/system/resources'),
  ])
  status.value = s.data.data || {}
  resources.value = r.data.data || {}
  lastResourceUpdate.value = new Date().toLocaleTimeString()
}

async function saveSchedule() {
  saving.value = true
  formError.value = ''
  try {
    let payload
    if (activeScope.value?.is_active) {
      payload = { scope: activeScopeKey.value, enabled: false, reason: form.reason || 'Maintenance ended' }
    } else {
      if (!form.reason || form.reason.trim().length < 5) { formError.value = 'Reason must be at least 5 characters.'; saving.value = false; return }
      if (!form.endStr) { formError.value = 'End time is required.'; saving.value = false; return }
      const start = form.startStr ? new Date(form.startStr) : null
      const end = new Date(form.endStr)
      if (isNaN(end.getTime())) { formError.value = 'End time is invalid.'; saving.value = false; return }
      if (start && end <= start) { formError.value = 'End must be after start.'; saving.value = false; return }
      if (end.getTime() <= Date.now()) { formError.value = 'End time must be in the future.'; saving.value = false; return }
      payload = {
        scope: activeScopeKey.value,
        enabled: true,
        reason: form.reason.trim(),
        scheduled_start: start ? start.toISOString() : null,
        expected_end: end.toISOString(),
      }
    }
    payload.auto_start_enforced = form.autoStartEnforced
    payload.auto_end_enforced = form.autoEndEnforced
    payload.display_meta = {
      custom_title: form.customTitle,
      custom_message: form.customMessage,
      show_countdown: form.showCountdown,
      allow_email_notification: form.allowEmailNotification,
    }
    payload.bypass_rules = {
      allowed_roles: fromCSV(form.allowedRoles),
      allowed_ip_ranges: fromCSV(form.allowedIPRanges),
      secret_query_param: form.secretQueryParam,
    }
    await apiClient.put('/api/v1/admin/system/maintenance', payload)
    modalOpen.value = false
    await refresh()
  } catch (e) {
    formError.value = e?.response?.data?.message || e?.response?.data?.error?.message || 'Save failed'
  } finally { saving.value = false }
}

async function cancelSchedule(scopeKey) {
  if (!confirm(`Cancel the ${scopeMeta[scopeKey].title} maintenance window?`)) return
  await apiClient.put('/api/v1/admin/system/maintenance', { scope: scopeKey, enabled: false, reason: 'Schedule cancelled' })
  await refresh()
}
async function flushCache() {
  if (!confirm('Invalidate all application caches?')) return
  await apiClient.post('/api/v1/admin/system/cache/flush')
  await refresh()
}
async function revokeAll() {
  if (prompt('Type REVOKE ALL to confirm') !== 'REVOKE ALL') return
  await apiClient.post('/api/v1/admin/system/sessions/revoke-all')
  await authStore.logout()
  location.assign('/login')
}

onMounted(() => { refresh(); timer = setInterval(refresh, 5000) })
onBeforeUnmount(() => clearInterval(timer))
</script>
