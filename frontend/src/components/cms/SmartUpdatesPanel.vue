<template>
  <section class="updates-panel">
    <div class="cms-panel-head">
      <h2>Smart Updates</h2>
      <button type="button" :disabled="checking" @click="check">{{ checking ? 'Checking…' : 'Check for updates' }}</button>
    </div>

    <div class="ops-grid">
      <article class="cms-card metric">
        <span>Current version</span>
        <strong>{{ status?.current_version || '—' }}</strong>
      </article>
      <article class="cms-card metric">
        <span>Latest release</span>
        <strong>{{ status?.latest?.tag_name || '—' }}</strong>
        <small v-if="status?.update_available" class="ops-update-flag">Update available</small>
        <small v-else>Up to date</small>
      </article>
      <article class="cms-card metric">
        <span>Disk</span>
        <strong>{{ status?.disk?.used_pct ?? '—' }}%</strong>
        <small>{{ status?.disk?.healthy ? 'Healthy' : 'Low space' }}</small>
      </article>
      <article class="cms-card metric">
        <span>Database</span>
        <strong>{{ status?.db_reachable ? 'Reachable' : 'Unreachable' }}</strong>
        <small>{{ status?.db_latency_ms ?? '—' }} ms</small>
      </article>
    </div>
    <p v-if="status?.last_error" class="ops-error-text">{{ status.last_error }}</p>
    <p v-if="status?.last_checked_at" class="ops-hint">Last checked {{ formatDateTime(status.last_checked_at) }}</p>

    <article class="cms-panel">
      <div class="cms-panel-head"><h3>Deploy Settings</h3></div>
      <div class="cms-two">
        <label>Repository <small>(owner/repo)</small><input v-model="settings.repo_slug" class="form-control" placeholder="atenimedia-llc/ncs-online" /></label>
        <label>Target branch<input v-model="settings.target_branch" class="form-control" placeholder="main" /></label>
        <label>Work tree<input v-model="settings.work_tree" class="form-control" placeholder="/app" /></label>
        <label>Compose project<input v-model="settings.compose_project" class="form-control" /></label>
        <label>Deploy services <small>(comma-separated)</small><input :value="(settings.deploy_services || []).join(', ')" @input="setList('deploy_services', $event.target.value)" class="form-control" /></label>
        <label>Preserve paths <small>(comma-separated)</small><input :value="(settings.preserve_paths || []).join(', ')" @input="setList('preserve_paths', $event.target.value)" class="form-control" /></label>
        <label class="wide">GitHub token <small>(personal access token, kept encrypted at rest)</small><input v-model="settings.github_token" type="password" class="form-control" :placeholder="settings.github_token === '********' ? 'Leave blank to keep the saved token' : 'ghp_…'" /></label>
      </div>
      <div class="profile-actions">
        <button type="button" :disabled="savingSettings" @click="saveSettings">{{ savingSettings ? 'Saving…' : 'Save settings' }}</button>
      </div>
    </article>

    <article class="cms-panel">
      <div class="cms-panel-head">
        <h3>Deploy</h3>
        <span class="ops-badge" :class="{ healthy: deploy?.state === 'success', stopped: deploy?.state === 'failed' }">{{ deploy?.state || 'idle' }}</span>
      </div>
      <div class="profile-actions">
        <button type="button" :disabled="deploying || deploy?.state === 'running'" @click="deploy_">{{ deploy?.state === 'running' ? 'Deploying…' : 'Start deploy' }}</button>
        <button type="button" class="ops-danger" :disabled="rollingBack" @click="rollback">{{ rollingBack ? 'Rolling back…' : 'Rollback' }}</button>
      </div>
      <pre class="ops-terminal">{{ (deploy?.lines || []).join('\n') || 'No deploy output yet.' }}</pre>
    </article>
  </section>
</template>

<script setup>
import { onMounted, onUnmounted, reactive, ref } from 'vue'
import Swal from 'sweetalert2'
import { checkForUpdates, getDeployStatus, getUpdateSettings, getUpdateStatus, saveUpdateSettings, startDeploy, startRollback } from '@/api/updates.js'

const emit = defineEmits(['message', 'error'])

const status = ref(null)
const settings = reactive({ repo_slug: '', target_branch: '', work_tree: '', compose_project: '', deploy_services: [], deploy_script: '', preserve_paths: [], github_token: '' })
const deploy = ref(null)
const checking = ref(false)
const savingSettings = ref(false)
const deploying = ref(false)
const rollingBack = ref(false)
let deployTimer = null

function formatDateTime(value) { return value ? new Date(value).toLocaleString() : '' }
function setList(field, raw) { settings[field] = raw.split(',').map(s => s.trim()).filter(Boolean) }

async function loadStatus() {
  try { status.value = await getUpdateStatus() } catch (err) { emit('error', err) }
}
async function loadSettings() {
  try { Object.assign(settings, await getUpdateSettings()) } catch (err) { emit('error', err) }
}
async function loadDeployStatus() {
  try {
    deploy.value = await getDeployStatus()
    if (deploy.value?.state === 'running' && !deployTimer) {
      deployTimer = setInterval(loadDeployStatus, 3000)
    } else if (deploy.value?.state !== 'running' && deployTimer) {
      clearInterval(deployTimer)
      deployTimer = null
    }
  } catch (err) { emit('error', err) }
}

async function check() {
  checking.value = true
  try {
    status.value = await checkForUpdates()
    emit('message', 'Checked GitHub for the latest release')
  } catch (err) {
    emit('error', err)
  } finally {
    checking.value = false
  }
}

async function saveSettings() {
  savingSettings.value = true
  try {
    Object.assign(settings, await saveUpdateSettings({ ...settings }))
    emit('message', 'Deploy settings saved')
  } catch (err) {
    emit('error', err)
  } finally {
    savingSettings.value = false
  }
}

async function deploy_() {
  const result = await Swal.fire({
    title: 'Start deployment?',
    text: 'This pulls the latest images and restarts services after taking an automatic backup. Type DEPLOY to confirm.',
    input: 'text',
    inputPlaceholder: 'DEPLOY',
    showCancelButton: true,
    confirmButtonText: 'Deploy',
    confirmButtonColor: '#6777ef',
    cancelButtonColor: '#94a3b8',
    inputValidator: value => value !== 'DEPLOY' ? 'Type exactly: DEPLOY' : undefined,
  })
  if (!result.isConfirmed) return
  deploying.value = true
  try {
    deploy.value = await startDeploy()
    emit('message', 'Deployment started')
    loadDeployStatus()
  } catch (err) {
    emit('error', err)
  } finally {
    deploying.value = false
  }
}

async function rollback() {
  const result = await Swal.fire({
    title: 'Roll back to the previous release?',
    text: 'Type ROLLBACK to confirm.',
    input: 'text',
    inputPlaceholder: 'ROLLBACK',
    showCancelButton: true,
    confirmButtonText: 'Rollback',
    confirmButtonColor: '#fc544b',
    cancelButtonColor: '#94a3b8',
    inputValidator: value => value !== 'ROLLBACK' ? 'Type exactly: ROLLBACK' : undefined,
  })
  if (!result.isConfirmed) return
  rollingBack.value = true
  try {
    deploy.value = await startRollback()
    emit('message', 'Rollback started')
    loadDeployStatus()
  } catch (err) {
    emit('error', err)
  } finally {
    rollingBack.value = false
  }
}

onMounted(() => {
  loadStatus()
  loadSettings()
  loadDeployStatus()
})
onUnmounted(() => clearInterval(deployTimer))
</script>

<style scoped>
/*
 * This component is rendered as a child of WebsiteContentManagerView, whose
 * <style scoped> only reaches this component's root element (Vue scoped CSS
 * does not cascade into child templates) — so shared classes are redefined
 * here rather than assumed to be inherited.
 */
.updates-panel { display: grid; gap: 20px; min-width: 0; }
.cms-panel-head { display: flex; justify-content: space-between; align-items: center; gap: 16px; flex-wrap: wrap; border-bottom: 1px solid #f9f9f9; padding-bottom: 15px; }
.cms-panel-head h2, .cms-panel-head h3 { margin: 0; font-size: 16px; font-weight: 700; color: #34395e; }
.cms-panel-head button { border: 0; border-radius: 30px; background: #6777ef; color: white; padding: 8px 18px; font-size: 12px; font-weight: 600; box-shadow: 0 2px 6px #acb5f6; }
.cms-panel-head button:disabled { opacity: .6; }
.cms-panel { background: white; border-radius: 3px; box-shadow: 0 4px 25px rgba(0,0,0,.1); padding: 25px; display: grid; gap: 16px; min-width: 0; }

.ops-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr)); gap: 16px; }
.cms-card.metric { background: white; box-shadow: 0 4px 25px rgba(0,0,0,.1); border-radius: 3px; padding: 20px; display: grid; gap: 6px; }
.cms-card.metric span { color: #98a6ad; font-size: 12px; font-weight: 600; }
.cms-card.metric strong { font-size: 20px; color: #34395e; }
.cms-card.metric small { color: #98a6ad; font-size: 11px; }
.ops-update-flag { color: #ffa426 !important; font-weight: 700; }
.ops-error-text { color: #fc544b; font-size: 13px; margin: 0; }
.ops-hint { color: #98a6ad; font-size: 12px; margin: 0; }

.cms-two { display: grid; grid-template-columns: repeat(auto-fit, minmax(15rem, 1fr)); gap: 18px; }
.cms-two label { display: grid; gap: 7px; font-size: 12px; font-weight: 600; color: #34395e; min-width: 0; }
.cms-two label small { font-weight: 500; text-transform: none; color: #94a3b8; }
.cms-two .wide { grid-column: 1 / -1; }
.cms-two input { border: 1px solid #e4e6fc; border-radius: 3px; padding: 10px 15px; font-weight: 500; color: #495057; background: #fdfdff; outline: none; width: 100%; box-sizing: border-box; }

.profile-actions { display: flex; flex-wrap: wrap; gap: 10px; }
.profile-actions button { border: 0; border-radius: 30px; background: #6777ef; color: white; padding: 8px 18px; font-size: 12px; font-weight: 600; box-shadow: 0 2px 6px #acb5f6; }
.profile-actions button:disabled { opacity: .6; }
.ops-danger { background: #fc544b !important; }

.ops-badge { border-radius: 30px; padding: 3px 10px; font-size: 11px; font-weight: 700; background: #eef0fd; color: #6777ef; }
.ops-badge.healthy { background: #e8f7f0; color: #47c363; }
.ops-badge.stopped { background: #fdeaea; color: #fc544b; }

.ops-terminal { background: #0f172a; color: #d1e7ff; border-radius: 6px; padding: 16px; font-family: 'Fira Code', monospace; font-size: 12px; max-height: 22rem; overflow: auto; white-space: pre-wrap; word-break: break-word; }

:global(.dark .cms-panel-head h2), :global(.dark .cms-panel-head h3) { color: #f8fafc; }
:global(.dark .cms-panel), :global(.dark .cms-card.metric) { background: #111827; border-color: #334155; }
:global(.dark .cms-card.metric strong) { color: #f8fafc; }
:global(.dark .cms-two input) { background: #0f172a; color: #f8fafc; border-color: #475569; }
</style>
