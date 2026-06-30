<template>
  <LayoutDefault title="Smart Updates">
    <div class="p-6 max-w-5xl mx-auto space-y-6">
      <Transition name="toast">
        <div v-if="toast.msg" :class="toast.ok ? 'bg-gray-900' : 'bg-red-600'" class="fixed top-5 right-5 z-[100] px-4 py-3 text-white text-sm rounded-xl shadow-xl">{{ toast.msg }}</div>
      </Transition>

      <header class="flex items-start justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Smart Updates</h2>
          <p class="text-xs text-gray-500 mt-1">Watches <code class="font-mono">{{ status.repo_slug || 'github' }}</code> releases. Deploys in-place via Docker Compose.</p>
        </div>
        <button @click="checkNow" :disabled="busy.check" class="text-sm font-medium px-4 py-2 bg-primary-600 hover:bg-primary-700 disabled:opacity-60 text-white rounded-lg">{{ busy.check ? 'Checking…' : 'Check now' }}</button>
      </header>

      <!-- Version + release card -->
      <section class="bg-white border border-gray-200 rounded-xl p-5">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
          <div>
            <div class="text-[10px] uppercase tracking-wider text-gray-400 font-semibold">Current Version</div>
            <div class="text-2xl font-bold text-gray-900 mt-1 font-mono">{{ status.current_version || '—' }}</div>
            <div class="text-[10px] text-gray-400 mt-1">From VERSION file in the running image</div>
          </div>
          <div>
            <div class="text-[10px] uppercase tracking-wider text-gray-400 font-semibold">Latest Release</div>
            <div class="text-2xl font-bold mt-1 font-mono" :class="status.update_available ? 'text-amber-600' : 'text-emerald-700'">{{ status.latest?.tag_name || (status.last_error ? '— error —' : '—') }}</div>
            <div class="text-[10px] text-gray-400 mt-1">{{ status.latest ? `Published ${formatRel(status.latest.published_at)}` : 'No release found' }}</div>
          </div>
          <div>
            <div class="text-[10px] uppercase tracking-wider text-gray-400 font-semibold">State</div>
            <div class="mt-1">
              <span v-if="status.update_available" class="text-[10px] uppercase tracking-wider font-bold bg-amber-100 text-amber-800 px-2 py-1 rounded-full">Update available</span>
              <span v-else-if="status.latest" class="text-[10px] uppercase tracking-wider font-bold bg-green-100 text-green-700 px-2 py-1 rounded-full">Up to date</span>
              <span v-else class="text-[10px] uppercase tracking-wider font-bold bg-gray-100 text-gray-700 px-2 py-1 rounded-full">Unknown</span>
            </div>
            <div class="text-[10px] text-gray-400 mt-2">{{ status.last_checked_at ? `Checked ${formatRel(status.last_checked_at)}` : 'Sentinel not yet run' }}</div>
          </div>
        </div>

        <div v-if="status.last_error" class="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-800">
          {{ status.last_error }}
        </div>

        <div v-if="status.latest?.body" class="mt-5 pt-5 border-t border-gray-100">
          <div class="text-[10px] uppercase tracking-wider text-gray-400 font-semibold mb-2">Release notes</div>
          <pre class="text-xs text-gray-700 whitespace-pre-wrap bg-gray-50 border border-gray-200 rounded-lg p-3 max-h-72 overflow-auto font-mono leading-relaxed">{{ status.latest.body }}</pre>
          <a :href="status.latest.html_url" target="_blank" rel="noopener" class="text-xs text-primary-600 hover:underline mt-2 inline-block">Open on GitHub →</a>
        </div>
      </section>

      <!-- Pre-flight -->
      <section class="bg-white border border-gray-200 rounded-xl p-5">
        <h3 class="text-sm font-semibold text-gray-700 uppercase tracking-wider mb-4">Pre-flight checks</h3>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="bg-gray-50 rounded-lg p-3">
            <div class="text-[10px] uppercase tracking-wider text-gray-400 font-semibold">Database</div>
            <div class="text-base font-semibold mt-1" :class="status.db_reachable ? 'text-emerald-700' : 'text-red-700'">{{ status.db_reachable ? 'Reachable' : 'Unreachable' }}</div>
            <div v-if="status.db_reachable" class="text-[10px] text-gray-400 mt-1">Ping {{ status.db_latency_ms }} ms</div>
          </div>
          <div class="bg-gray-50 rounded-lg p-3">
            <div class="text-[10px] uppercase tracking-wider text-gray-400 font-semibold">Disk Free</div>
            <div class="text-base font-semibold mt-1" :class="status.disk?.healthy ? 'text-emerald-700' : 'text-amber-700'">{{ formatBytes(status.disk?.free_bytes) }} free</div>
            <div class="text-[10px] text-gray-400 mt-1">{{ formatBytes(status.disk?.total_bytes) }} total · {{ status.disk?.used_pct ?? 0 }}% used</div>
          </div>
          <div class="bg-gray-50 rounded-lg p-3">
            <div class="text-[10px] uppercase tracking-wider text-gray-400 font-semibold">CI signature check</div>
            <div class="text-base font-semibold mt-1 text-gray-400">Not wired</div>
            <div class="text-[10px] text-gray-400 mt-1">TODO: poll CI conclusion for tag</div>
          </div>
        </div>
      </section>

      <section class="bg-white border border-gray-200 rounded-xl p-5">
        <div class="flex items-start justify-between gap-4 mb-4">
          <div>
            <h3 class="text-sm font-semibold text-gray-700 uppercase tracking-wider">External update configuration</h3>
            <p class="text-xs text-gray-500 mt-1">Configure GitHub release checks and deployment commands without touching uploaded files or persisted volumes.</p>
          </div>
          <button @click="saveSettings" :disabled="busy.settings" class="text-sm font-medium px-4 py-2 bg-primary-600 hover:bg-primary-700 disabled:opacity-60 text-white rounded-lg">{{ busy.settings ? 'Saving...' : 'Save config' }}</button>
        </div>
        <div class="grid md:grid-cols-2 gap-4">
          <label class="block text-sm font-semibold text-gray-700">GitHub repo slug
            <input v-model="settings.repo_slug" class="mt-1 w-full rounded-lg border-gray-200 text-sm font-mono" placeholder="owner/repo" />
          </label>
          <label class="block text-sm font-semibold text-gray-700">GitHub token
            <input v-model="settings.github_token" type="password" class="mt-1 w-full rounded-lg border-gray-200 text-sm font-mono" placeholder="********" />
          </label>
          <label class="block text-sm font-semibold text-gray-700">Compose project
            <input v-model="settings.compose_project" class="mt-1 w-full rounded-lg border-gray-200 text-sm font-mono" placeholder="ncs-online" />
          </label>
          <label class="block text-sm font-semibold text-gray-700">Deploy services
            <input v-model="servicesText" class="mt-1 w-full rounded-lg border-gray-200 text-sm font-mono" placeholder="backend, frontend, nsmis-worker, backup" />
          </label>
          <label class="block text-sm font-semibold text-gray-700 md:col-span-2">Preserved volumes / paths
            <input v-model="preserveText" class="mt-1 w-full rounded-lg border-gray-200 text-sm font-mono" placeholder="uploads_data, private_data, app_logs, backups_data, postgres_data" />
          </label>
          <label class="block text-sm font-semibold text-gray-700 md:col-span-2">Optional deploy script
            <textarea v-model="settings.deploy_script" rows="6" class="mt-1 w-full rounded-lg border-gray-200 text-sm font-mono" placeholder="docker compose pull backend frontend nsmis-worker backup&#10;docker compose up -d --no-deps backend frontend nsmis-worker backup"></textarea>
            <span class="block text-[10px] text-gray-400 mt-1">Leave blank to use the built-in safe deploy path. Do not remove volume bindings for uploads/config data.</span>
          </label>
        </div>
      </section>

      <!-- Deploy controls -->
      <section class="bg-white border border-gray-200 rounded-xl p-5">
        <div class="flex items-start justify-between mb-4">
          <div>
            <h3 class="text-sm font-semibold text-gray-700 uppercase tracking-wider">In-place deploy</h3>
            <p class="text-xs text-gray-500 mt-1">Pulls the latest images and recreates <code class="font-mono">backend, frontend, nsmis-worker, backup</code>. Postgres + nginx are left untouched.</p>
          </div>
          <span :class="deployBadgeClass" class="text-[10px] uppercase tracking-wider font-bold px-2 py-1 rounded-full">{{ deploy.state }}</span>
        </div>

        <div class="flex flex-wrap gap-2 mb-4">
          <button @click="onDeploy" :disabled="deploy.state === 'running' || busy.deploy" class="text-sm font-medium px-4 py-2 bg-emerald-600 hover:bg-emerald-700 disabled:opacity-60 text-white rounded-lg">
            {{ deploy.state === 'running' ? 'Deploying…' : (status.update_available ? `Deploy ${status.latest?.tag_name || ''}` : 'Re-pull current') }}
          </button>
          <button @click="onRollback" :disabled="deploy.state === 'running' || busy.deploy" class="text-sm font-medium px-4 py-2 border border-red-300 text-red-700 hover:bg-red-50 disabled:opacity-60 rounded-lg">
            Rollback to previous
          </button>
        </div>

        <pre v-if="deploy.lines?.length" class="text-xs bg-gray-900 text-gray-100 rounded-lg p-3 max-h-80 overflow-auto font-mono leading-relaxed">{{ deploy.lines.join('\n') }}</pre>
        <p v-else class="text-xs text-gray-400 italic">No deploy has run in this session yet.</p>

        <div v-if="deploy.error" class="mt-3 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-800">{{ deploy.error }}</div>
      </section>

      <!-- Operator notes -->
      <section class="bg-amber-50 border border-amber-200 rounded-xl p-4 text-xs text-amber-900 space-y-1">
        <p class="font-semibold uppercase tracking-wider text-[10px]">Operator notes</p>
        <p>• Sentinel polls GitHub every 5 minutes; "Check now" forces an immediate refresh.</p>
        <p>• Deploy + rollback only work when the Docker socket is mounted into the backend container — see <code class="font-mono">Docs/runbooks/smart-updates.md</code>.</p>
        <p>• Rollback uses the <code class="font-mono">:previous</code> image tag created during the last successful deploy. It's a no-op on the first deploy.</p>
      </section>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { getUpdateStatus, checkForUpdates, getDeployStatus, startDeploy, startRollback, getUpdateSettings, saveUpdateSettings } from '@/api/updates.js'

const status = ref({})
const deploy = reactive({ state: 'idle', lines: [], error: '' })
const busy = reactive({ check: false, deploy: false, settings: false })
const toast = reactive({ msg: '', ok: true })
const settings = reactive({ repo_slug: '', github_token: '', compose_project: '', deploy_services: [], deploy_script: '', preserve_paths: [] })
const servicesText = ref('')
const preserveText = ref('')
let pollTimer = null

function showToast(msg, ok = true) {
  toast.msg = msg; toast.ok = ok
  setTimeout(() => (toast.msg = ''), 2800)
}

function formatBytes(n) {
  if (!n || n <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0; let v = n
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++ }
  return `${v.toFixed(v >= 100 ? 0 : 1)} ${units[i]}`
}

function formatRel(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  const diff = (Date.now() - d.getTime()) / 1000
  if (diff < 60) return 'just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return d.toLocaleDateString()
}

const deployBadgeClass = computed(() => ({
  idle: 'bg-gray-100 text-gray-700',
  running: 'bg-blue-100 text-blue-700',
  success: 'bg-green-100 text-green-700',
  failed: 'bg-red-100 text-red-700',
}[deploy.state] || 'bg-gray-100 text-gray-700'))

async function refresh() {
  try { status.value = await getUpdateStatus() } catch { /* keep last */ }
  try { Object.assign(deploy, await getDeployStatus()) } catch { /* keep */ }
}

async function loadSettings() {
  try {
    Object.assign(settings, await getUpdateSettings())
    servicesText.value = (settings.deploy_services || []).join(', ')
    preserveText.value = (settings.preserve_paths || []).join(', ')
  } catch {
    showToast('Could not load update settings', false)
  }
}

async function saveSettings() {
  busy.settings = true
  try {
    const payload = {
      ...settings,
      deploy_services: csv(servicesText.value),
      preserve_paths: csv(preserveText.value),
    }
    Object.assign(settings, await saveUpdateSettings(payload))
    servicesText.value = (settings.deploy_services || []).join(', ')
    preserveText.value = (settings.preserve_paths || []).join(', ')
    showToast('Update settings saved')
    await checkNow()
  } catch (e) {
    showToast(e?.response?.data?.error?.message || 'Could not save settings', false)
  } finally {
    busy.settings = false
  }
}

function csv(value) {
  return String(value || '').split(',').map(v => v.trim()).filter(Boolean)
}

async function checkNow() {
  busy.check = true
  try {
    status.value = await checkForUpdates()
    showToast('Checked')
  } catch (e) {
    showToast(e?.response?.data?.message || 'GitHub unreachable', false)
  } finally { busy.check = false }
}

async function onDeploy() {
  if (!confirm('Pull the latest images and restart backend / frontend / worker / backup? Postgres + nginx are left running.')) return
  busy.deploy = true
  try {
    Object.assign(deploy, await startDeploy())
    showToast('Deploy started')
  } catch (e) {
    showToast(e?.response?.data?.message || 'Deploy failed to start', false)
  } finally { busy.deploy = false }
}

async function onRollback() {
  if (!confirm('Roll back to the previous image tag and restart? This only works if a prior deploy created a :previous tag.')) return
  busy.deploy = true
  try {
    Object.assign(deploy, await startRollback())
    showToast('Rollback started')
  } catch (e) {
    showToast(e?.response?.data?.message || 'Rollback failed to start', false)
  } finally { busy.deploy = false }
}

onMounted(() => {
  loadSettings()
  refresh()
  // While a deploy is running, poll status every 2s; otherwise every 30s.
  pollTimer = setInterval(refresh, 2000)
})
onBeforeUnmount(() => { if (pollTimer) clearInterval(pollTimer) })
</script>

<style scoped>
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(-8px); }
.toast-enter-active, .toast-leave-active { transition: all .25s ease; }
</style>
