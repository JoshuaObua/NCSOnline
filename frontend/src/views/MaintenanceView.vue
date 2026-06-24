<template>
  <LayoutDefault title="System Operations">
    <div class="space-y-6 p-6 max-w-6xl mx-auto">
      <!-- ───────── Maintenance window status ───────── -->
      <section class="rounded-2xl border p-6 transition-colors"
        :class="bannerClass">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-5">
          <div class="min-w-0">
            <p class="text-xs font-bold uppercase tracking-widest" :class="bannerLabelClass">{{ bannerLabel }}</p>
            <h2 class="text-2xl font-bold text-primary-800 mt-1">{{ bannerTitle }}</h2>
            <p class="text-sm text-gray-600 mt-2 max-w-3xl">{{ status.reason || 'Public and authenticated services are available.' }}</p>
            <dl v-if="status.scheduled_start || status.expected_end" class="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-1 text-xs text-gray-700">
              <div v-if="status.scheduled_start"><dt class="inline font-semibold text-gray-500 mr-1">Starts:</dt><dd class="inline">{{ formatDate(status.scheduled_start) }} <span class="text-gray-400">· {{ formatRel(status.scheduled_start) }}</span></dd></div>
              <div v-if="status.expected_end"><dt class="inline font-semibold text-gray-500 mr-1">Ends:</dt><dd class="inline">{{ formatDate(status.expected_end) }} <span class="text-gray-400">· {{ formatRel(status.expected_end) }}</span></dd></div>
            </dl>
          </div>
          <div v-if="authStore.isSuperAdmin" class="flex flex-shrink-0 gap-2">
            <button @click="openScheduler" class="rounded-xl px-5 py-3 font-bold text-white"
              :class="status.is_active ? 'bg-emerald-700 hover:bg-emerald-800' : 'bg-amber-600 hover:bg-amber-700'">
              {{ status.is_active ? 'End maintenance now' : (status.is_scheduled ? 'Edit schedule' : 'Schedule maintenance') }}
            </button>
            <button v-if="status.is_scheduled || (status.maintenance_mode && !status.is_active)" @click="cancelSchedule" class="rounded-xl px-4 py-3 text-sm font-semibold border border-gray-200 text-gray-700 hover:bg-gray-50">
              Cancel schedule
            </button>
          </div>
        </div>
      </section>

      <!-- ───────── Live resources ───────── -->
      <section>
        <div class="flex items-center justify-between mb-3">
          <div>
            <h2 class="font-bold text-primary-800">Live resources</h2>
            <p class="text-xs text-gray-500">Updated every five seconds from the API container host.</p>
          </div>
          <span class="text-xs text-gray-400">{{ lastResourceUpdate }}</span>
        </div>
        <div class="grid sm:grid-cols-2 xl:grid-cols-4 gap-4">
          <article v-for="card in resourceCards" :key="card.label" class="bg-white border border-gray-100 rounded-2xl p-5 shadow-sm">
            <p class="text-xs uppercase tracking-wide text-gray-400 font-semibold">{{ card.label }}</p>
            <p class="text-2xl font-bold text-primary-800 mt-2">{{ card.value }}</p>
            <div class="h-1.5 rounded-full bg-gray-100 mt-4">
              <div class="h-full rounded-full bg-primary-600" :style="{width: Math.min(card.percent || 0, 100) + '%'}"/>
            </div>
          </article>
        </div>
      </section>

      <!-- ───────── Ops buttons ───────── -->
      <section class="grid md:grid-cols-2 gap-4">
        <article class="bg-white border border-gray-100 rounded-2xl p-6">
          <h3 class="font-bold text-primary-800">Application cache</h3>
          <p class="text-sm text-gray-500 mt-2">Invalidate memory-cached settings and force resolvers to reload.</p>
          <p class="text-xs text-gray-400 mt-3">Generation {{ status.cache_generation || 0 }}</p>
          <button v-if="authStore.isSuperAdmin" @click="flushCache" class="mt-4 rounded-xl bg-primary-700 hover:bg-primary-800 text-white px-4 py-2.5 font-semibold">Flush cache</button>
        </article>
        <article class="bg-white border border-red-100 rounded-2xl p-6">
          <h3 class="font-bold text-red-800">Emergency session revocation</h3>
          <p class="text-sm text-gray-500 mt-2">Revoke every refresh token and reject all currently issued access tokens.</p>
          <button v-if="authStore.isSuperAdmin" @click="revokeAll" class="mt-4 rounded-xl bg-red-700 hover:bg-red-800 text-white px-4 py-2.5 font-semibold">Force logout all users</button>
        </article>
      </section>
    </div>

    <!-- ───────── Scheduler modal ───────── -->
    <Modal :open="modalOpen" :title="modalTitle" @close="modalOpen = false">
      <div v-if="status.is_active" class="space-y-3">
        <p class="text-sm text-gray-700">End the active maintenance window now? Traffic will resume as soon as you confirm.</p>
        <label class="block text-sm font-semibold">Closing note
          <textarea v-model="form.reason" rows="2" class="mt-1 w-full rounded-xl border-gray-200 text-sm" placeholder="e.g. Database restore completed; resuming service."></textarea>
        </label>
      </div>

      <div v-else class="space-y-4">
        <label class="block text-sm font-semibold">Reason <span class="text-red-500">*</span>
          <textarea v-model="form.reason" rows="2" class="mt-1 w-full rounded-xl border-gray-200 text-sm" placeholder="Describe the approved maintenance work"></textarea>
          <p class="text-[10px] text-gray-400 mt-1">Shown on the 503 page; minimum 5 characters.</p>
        </label>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <label class="block text-sm font-semibold">Start
            <input v-model="form.startStr" type="datetime-local" class="mt-1 w-full rounded-xl border-gray-200 text-sm"/>
            <p class="text-[10px] text-gray-400 mt-1">Leave blank to start immediately.</p>
          </label>
          <label class="block text-sm font-semibold">End <span class="text-red-500">*</span>
            <input v-model="form.endStr" type="datetime-local" class="mt-1 w-full rounded-xl border-gray-200 text-sm"/>
            <p class="text-[10px] text-gray-400 mt-1">Traffic resumes automatically when this passes.</p>
          </label>
        </div>

        <div class="flex flex-wrap gap-2 text-xs">
          <button type="button" @click="presetWindow(15)" class="px-2.5 py-1 border border-gray-200 rounded-full hover:bg-gray-50">+15 min</button>
          <button type="button" @click="presetWindow(60)" class="px-2.5 py-1 border border-gray-200 rounded-full hover:bg-gray-50">+1 hr</button>
          <button type="button" @click="presetWindow(180)" class="px-2.5 py-1 border border-gray-200 rounded-full hover:bg-gray-50">+3 hr</button>
          <button type="button" @click="presetTonight" class="px-2.5 py-1 border border-gray-200 rounded-full hover:bg-gray-50">Tonight 22:00 → 04:00</button>
        </div>

        <p v-if="formError" class="text-xs text-red-600">{{ formError }}</p>
      </div>

      <template #footer>
        <button @click="modalOpen = false" class="px-4 py-2 text-sm">Cancel</button>
        <button @click="saveSchedule" :disabled="saving" class="rounded-xl bg-primary-700 hover:bg-primary-800 disabled:opacity-60 px-5 py-2 text-sm text-white font-bold">
          {{ saving ? 'Saving…' : (status.is_active ? 'End now' : (isFutureStart ? 'Schedule window' : 'Start now')) }}
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
const modalOpen = ref(false)
const saving = ref(false)
const formError = ref('')
const lastResourceUpdate = ref('Not loaded')
let timer

const form = reactive({ reason: '', startStr: '', endStr: '' })

// ── Banner state derived from status ──
const bannerClass = computed(() => {
  if (status.value.is_active) return 'border-amber-300 bg-amber-50'
  if (status.value.is_scheduled) return 'border-blue-300 bg-blue-50'
  return 'border-emerald-200 bg-emerald-50'
})
const bannerLabelClass = computed(() => {
  if (status.value.is_active) return 'text-amber-700'
  if (status.value.is_scheduled) return 'text-blue-700'
  return 'text-emerald-700'
})
const bannerLabel = computed(() => {
  if (status.value.is_active) return 'Maintenance window active'
  if (status.value.is_scheduled) return 'Maintenance scheduled'
  return 'System status'
})
const bannerTitle = computed(() => {
  if (status.value.is_active) return 'Maintenance mode is on'
  if (status.value.is_scheduled) return 'Maintenance is queued'
  return 'Platform operational'
})

// ── Resource cards ──
const pct = v => Number.isFinite(v) ? v : 0
const fmtBytes = v => {
  if (!v) return '0 B'
  const units = ['B','KB','MB','GB','TB']
  let i = 0
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++ }
  return `${v.toFixed(1)} ${units[i]}`
}
const resourceCards = computed(() => [
  { label: 'CPU',    value: `${pct(resources.value.cpu_percent?.[0]).toFixed(1)}%`, percent: pct(resources.value.cpu_percent?.[0]) },
  { label: 'Memory', value: `${pct(resources.value.memory?.usedPercent).toFixed(1)}%`, percent: pct(resources.value.memory?.usedPercent) },
  { label: 'Disk',   value: `${pct(resources.value.disk?.usedPercent).toFixed(1)}%`, percent: pct(resources.value.disk?.usedPercent) },
  { label: 'Heap',   value: fmtBytes(resources.value.heap_alloc_bytes), percent: 0 },
])

const isFutureStart = computed(() => {
  if (!form.startStr) return false
  return new Date(form.startStr).getTime() > Date.now() + 30_000
})

const modalTitle = computed(() => {
  if (status.value.is_active) return 'End maintenance window'
  if (status.value.is_scheduled) return 'Edit scheduled maintenance'
  return 'Schedule maintenance'
})

// ── Helpers ──
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

function openScheduler() {
  formError.value = ''
  form.reason = status.value.reason || ''
  form.startStr = toLocalInput(status.value.scheduled_start)
  form.endStr = toLocalInput(status.value.expected_end)
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
    if (status.value.is_active) {
      // End-now path: explicit disable, clears window server-side.
      payload = { enabled: false, reason: form.reason || 'Maintenance ended' }
    } else {
      if (!form.reason || form.reason.trim().length < 5) {
        formError.value = 'Reason must be at least 5 characters.'; saving.value = false; return
      }
      if (!form.endStr) {
        formError.value = 'End time is required.'; saving.value = false; return
      }
      const start = form.startStr ? new Date(form.startStr) : null
      const end = new Date(form.endStr)
      if (isNaN(end.getTime())) {
        formError.value = 'End time is invalid.'; saving.value = false; return
      }
      if (start && end <= start) {
        formError.value = 'End must be after start.'; saving.value = false; return
      }
      if (end.getTime() <= Date.now()) {
        formError.value = 'End time must be in the future.'; saving.value = false; return
      }
      payload = {
        enabled: true,
        reason: form.reason.trim(),
        scheduled_start: start ? start.toISOString() : null,
        expected_end: end.toISOString(),
      }
    }
    await apiClient.put('/api/v1/admin/system/maintenance', payload)
    modalOpen.value = false
    await refresh()
  } catch (e) {
    formError.value = e?.response?.data?.message || 'Save failed'
  } finally { saving.value = false }
}

async function cancelSchedule() {
  if (!confirm('Cancel the scheduled maintenance window?')) return
  await apiClient.put('/api/v1/admin/system/maintenance', { enabled: false, reason: 'Schedule cancelled' })
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
