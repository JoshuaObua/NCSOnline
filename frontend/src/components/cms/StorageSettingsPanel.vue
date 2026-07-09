<template>
  <section class="storage-settings">
    <div class="cms-panel-head">
      <h2>Storage Settings</h2>
      <button type="button" :disabled="saving" @click="save">{{ saving ? 'Saving…' : 'Save settings' }}</button>
    </div>

    <div class="cms-panel storage-routing">
      <p class="storage-hint">Pick which provider handles each kind of upload. Fill in credentials for any provider below, then switch using these two selects — no need to re-enter anything when you change your mind.</p>
      <div class="cms-two">
        <label>Public portal uploads <small>(slides, posts, facilities, etc.)</small>
          <select v-model="settings.public_provider">
            <option v-for="p in providers" :key="p.value" :value="p.value">{{ p.label }}</option>
          </select>
        </label>
        <label>Application evidence uploads <small>(applicant attachments)</small>
          <select v-model="settings.application_provider">
            <option v-for="p in providers" :key="p.value" :value="p.value">{{ p.label }}</option>
          </select>
        </label>
      </div>
    </div>

    <div class="cms-panel">
      <div class="cms-panel-head"><h3>Local storage</h3></div>
      <div class="cms-two">
        <label>Public upload path<input v-model="settings.local_public_path" placeholder="/app/uploads" /></label>
        <label>Public URL prefix<input v-model="settings.local_public_url_prefix" placeholder="/uploads" /></label>
        <label>Application upload path<input v-model="settings.local_app_path" placeholder="/app/uploads/applications" /></label>
        <label>Application URL prefix<input v-model="settings.local_app_url_prefix" placeholder="/uploads/applications" /></label>
      </div>
    </div>

    <div class="cms-panel">
      <div class="cms-panel-head">
        <h3>Supabase Storage</h3>
        <span class="storage-badge" :class="{ ok: supabaseConfigured }">{{ supabaseConfigured ? 'Service key saved' : 'Not configured' }}</span>
      </div>
      <div class="cms-two">
        <label>Project URL<input v-model="settings.supabase_url" placeholder="https://xyzcompany.supabase.co" /></label>
        <label>Bucket<input v-model="settings.supabase_bucket" placeholder="ncsportal-uploads" /></label>
        <label>Path prefix<input v-model="settings.supabase_prefix" placeholder="uploads" /></label>
        <label class="check"><input v-model="settings.supabase_public_read" type="checkbox" /> Bucket allows public read</label>
        <label class="wide">Service role key<input v-model="settings.supabase_service_key" type="password" autocomplete="new-password" :placeholder="supabaseConfigured ? 'Leave blank to keep the saved key' : 'service_role key from Supabase project settings'" /></label>
      </div>
    </div>

    <div class="cms-panel">
      <div class="cms-panel-head"><h3>Amazon S3</h3></div>
      <div class="cms-two">
        <label>Bucket<input v-model="settings.s3_bucket" /></label>
        <label>Region<input v-model="settings.s3_region" placeholder="us-east-1" /></label>
        <label>Path prefix<input v-model="settings.s3_prefix" placeholder="applications" /></label>
        <label>Custom endpoint <small>(optional, for S3-compatible services)</small><input v-model="settings.s3_endpoint" placeholder="https://s3.us-east-1.amazonaws.com" /></label>
        <label>Public base URL <small>(optional, for CDN/custom domain)</small><input v-model="settings.s3_public_base_url" /></label>
        <label class="check"><input v-model="settings.s3_force_path_style" type="checkbox" /> Force path-style URLs</label>
        <label>Access key ID<input v-model="settings.s3_access_key_id" :placeholder="s3KeyPlaceholder" /></label>
        <label>Secret access key<input v-model="settings.s3_secret_access_key" type="password" autocomplete="new-password" placeholder="Leave blank to keep the saved key" /></label>
      </div>
    </div>

    <div class="cms-panel">
      <div class="cms-panel-head">
        <h3>Google Drive</h3>
        <span class="storage-badge" :class="{ ok: googleDriveConnected }">{{ googleDriveConnected ? 'Connected' : 'Not connected' }}</span>
      </div>
      <div class="cms-two">
        <label>Folder ID<input v-model="settings.google_drive_folder_id" /></label>
        <label>OAuth client ID<input v-model="settings.google_drive_client_id" /></label>
        <label>OAuth client secret<input v-model="settings.google_drive_client_secret" type="password" autocomplete="new-password" placeholder="Leave blank to keep the saved secret" /></label>
        <label class="check"><input v-model="settings.google_drive_make_public" type="checkbox" /> Make uploaded files public</label>
      </div>
      <div class="storage-actions">
        <button type="button" :disabled="connecting" @click="connectDrive">{{ connecting ? 'Connecting…' : (googleDriveConnected ? 'Reconnect' : 'Connect Google Drive') }}</button>
        <button v-if="googleDriveConnected" type="button" class="storage-disconnect" @click="disconnectDrive">Disconnect</button>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import {
  connectGoogleDrive,
  disconnectGoogleDrive,
  exchangeGoogleDriveCode,
  getStorageSettings,
  updateStorageSettings,
} from '@/api/storageSettings.js'

const emit = defineEmits(['message', 'error'])

const providers = [
  { value: 'local', label: 'Local storage' },
  { value: 'supabase', label: 'Supabase Storage' },
  { value: 's3', label: 'Amazon S3' },
  { value: 'google_drive', label: 'Google Drive' },
]

const settings = reactive({
  public_provider: 'local',
  application_provider: 'local',
  local_public_path: '', local_public_url_prefix: '', local_app_path: '', local_app_url_prefix: '',
  supabase_url: '', supabase_bucket: '', supabase_prefix: '', supabase_public_read: false, supabase_service_key: '',
  s3_bucket: '', s3_region: '', s3_prefix: '', s3_endpoint: '', s3_public_base_url: '', s3_force_path_style: false,
  s3_access_key_id: '', s3_secret_access_key: '',
  google_drive_folder_id: '', google_drive_client_id: '', google_drive_client_secret: '', google_drive_make_public: true,
})
const googleDriveConnected = ref(false)
const supabaseConfigured = ref(false)
const saving = ref(false)
const connecting = ref(false)

const s3KeyPlaceholder = computed(() => settings.s3_access_key_id?.startsWith('********') ? '' : 'AKIA…')
const redirectUri = `${window.location.origin}/cms`

function applyResponse(data) {
  Object.assign(settings, data)
  googleDriveConnected.value = !!data.google_drive_connected
  supabaseConfigured.value = !!data.supabase_configured
  settings.s3_secret_access_key = ''
  settings.supabase_service_key = ''
  settings.google_drive_client_secret = ''
}

async function load() {
  try {
    applyResponse(await getStorageSettings())
  } catch (err) { emit('error', err) }
}

async function save() {
  saving.value = true
  try {
    applyResponse(await updateStorageSettings({ ...settings }))
    emit('message', 'Storage settings saved')
  } catch (err) {
    emit('error', err)
  } finally {
    saving.value = false
  }
}

async function connectDrive() {
  connecting.value = true
  try {
    if (!settings.google_drive_client_id || !settings.google_drive_folder_id) {
      throw new Error('Fill in the Google Drive folder ID and OAuth client ID/secret first.')
    }
    await save()
    const { auth_url } = await connectGoogleDrive(redirectUri)
    window.location.href = auth_url
  } catch (err) {
    emit('error', err)
    connecting.value = false
  }
}

async function disconnectDrive() {
  try {
    await disconnectGoogleDrive()
    googleDriveConnected.value = false
    emit('message', 'Google Drive disconnected')
  } catch (err) { emit('error', err) }
}

async function finishDriveConnectIfReturning() {
  const params = new URLSearchParams(window.location.search)
  const code = params.get('code')
  const state = params.get('state')
  if (!code || !state) return
  window.history.replaceState({}, '', window.location.pathname)
  try {
    await exchangeGoogleDriveCode({ code, state, redirectUri })
    await load()
    emit('message', 'Google Drive connected')
  } catch (err) {
    emit('error', err)
  }
}

onMounted(async () => {
  await load()
  await finishDriveConnectIfReturning()
})
</script>

<style scoped>
/*
 * This component is rendered as a child of WebsiteContentManagerView, whose
 * <style scoped> only reaches this component's root element (Vue scoped CSS
 * does not cascade into child templates) — so .cms-panel/.cms-two/etc are
 * redefined here rather than assumed to be inherited.
 */
.storage-settings { display: grid; gap: 20px; min-width: 0; }

.cms-panel-head { display: flex; justify-content: space-between; align-items: center; gap: 16px; flex-wrap: wrap; border-bottom: 1px solid #f9f9f9; padding-bottom: 15px; }
.cms-panel-head h2,
.cms-panel-head h3 { margin: 0; font-size: 16px; font-weight: 700; color: #34395e; }
.cms-panel-head button,
.storage-actions button { border: 0; border-radius: 30px; background: #6777ef; color: white; padding: 8px 18px; font-size: 12px; font-weight: 600; white-space: nowrap; box-shadow: 0 2px 6px #acb5f6; }
.cms-panel-head button:disabled { opacity: .6; }

.cms-panel { background: white; border: 0; border-radius: 3px; box-shadow: 0 4px 25px rgba(0,0,0,.1); padding: 25px; display: grid; gap: 20px; min-width: 0; }

.cms-two { display: grid; grid-template-columns: repeat(auto-fit, minmax(15rem, 1fr)); gap: 18px; }
.cms-two label { display: grid; gap: 7px; font-size: 12px; font-weight: 600; color: #34395e; min-width: 0; }
.cms-two label small { font-weight: 500; text-transform: none; color: #94a3b8; }
.cms-two .wide { grid-column: 1 / -1; }
.cms-two .check { display: flex; flex-direction: row; align-items: center; gap: .5rem; }
.cms-two input,
.cms-two select { width: 100%; box-sizing: border-box; border: 1px solid #e4e6fc; border-radius: 3px; padding: 10px 15px; font-weight: 500; color: #495057; background: #fdfdff; outline: none; }
.cms-two input:focus,
.cms-two select:focus { border-color: #6777ef; box-shadow: 0 2px 6px #acb5f6; }
.cms-two input[type="checkbox"] { width: 18px; height: 18px; accent-color: #6777ef; padding: 0; }

.storage-routing { background: #fff; }
.storage-hint { color: #6c757d; font-size: 13px; margin: 0; line-height: 1.7; }

.storage-badge { font-size: 12px; font-weight: 700; padding: 5px 10px; border-radius: 30px; background: #fdeaea; color: #fc544b; white-space: nowrap; }
.storage-badge.ok { background: #e8f7f0; color: #47c363; }
.storage-actions { display: flex; flex-wrap: wrap; gap: .5rem; }
.storage-disconnect { background: #fc544b !important; color: #fff !important; box-shadow: 0 2px 6px #fd9b96 !important; }

:global(.dark .cms-panel-head h2),
:global(.dark .cms-panel-head h3 ){ color: #f8fafc; }
:global(.dark .cms-panel ){ background: #111827; border-color: #334155; }
:global(.dark .cms-two input),
:global(.dark .cms-two select ){ background: #0f172a; color: #f8fafc; border-color: #475569; }
:global(.dark .storage-routing ){ background: #0f172a; }
:global(.dark .storage-badge ){ background: #422006; color: #fca5a5; }
:global(.dark .storage-badge.ok ){ background: #052e22; color: #6ee7b7; }
:global(.dark .storage-disconnect ){ background: #0f172a !important; border-color: #475569; }
</style>
