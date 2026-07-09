<template>
  <AccountShell title="Settings & Security" subtitle="Update credentials and manage active authenticated devices.">
    <section class="account-grid">
      <article class="account-panel">
        <h3>Credential Update</h3>
        <form class="security-form" @submit.prevent="submitPassword">
          <label>Current Password<input v-model="password.current_password" type="password" autocomplete="current-password" required /></label>
          <label>New Password<input v-model="password.new_password" type="password" autocomplete="new-password" required /></label>
          <label>Confirm New Password<input v-model="confirmPassword" type="password" autocomplete="new-password" required /></label>
          <div class="strength-track"><span :style="{ width: `${strength.score * 25}%`, background: strength.color }"></span></div>
          <p class="validation-copy">{{ strength.label }}</p>
          <button type="submit" :disabled="savingPassword || !canSubmitPassword">{{ savingPassword ? 'Updating...' : 'Update Password' }}</button>
          <p v-if="passwordMessage" class="status-copy">{{ passwordMessage }}</p>
        </form>
      </article>

      <article class="account-panel">
        <div class="panel-title-row">
          <h3>Active Device Sessions</h3>
          <button type="button" class="secondary-action" @click="clearOthers">Terminate All Other Sessions</button>
        </div>
        <div class="session-list">
          <article v-for="(session, index) in sessions" :key="session.id" class="session-card">
            <div>
              <strong>{{ deviceLabel(session.user_agent) }}</strong>
              <p>{{ session.ip_address || 'Unknown IP' }} · {{ locationLabel(session) }}</p>
              <small>{{ formatDate(session.created_at) }} · expires {{ formatDate(session.expires_at) }}</small>
            </div>
            <span v-if="index === 0" class="current-badge">This Device</span>
            <button v-else type="button" @click="revoke(session.id)">Revoke</button>
          </article>
          <p v-if="!sessions.length" class="empty-copy">No active sessions found.</p>
        </div>
      </article>
    </section>
  </AccountShell>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import AccountShell from './AccountShell.vue'
import { listSessions, revokeOtherSessions, revokeSession, updatePassword } from '@/api/account.js'

const password = reactive({ current_password: '', new_password: '' })
const confirmPassword = ref('')
const savingPassword = ref(false)
const passwordMessage = ref('')
const sessions = ref([])

const strength = computed(() => {
  const value = password.new_password
  let score = 0
  if (value.length >= 8) score++
  if (/[A-Z]/.test(value) && /[a-z]/.test(value)) score++
  if (/\d/.test(value)) score++
  if (/[^A-Za-z0-9]/.test(value)) score++
  const labels = ['Too weak', 'Weak', 'Fair', 'Strong', 'Excellent']
  const colors = ['#dc2626', '#ea580c', '#d97706', '#16a34a', '#047857']
  return { score, label: labels[score], color: colors[score] }
})
const canSubmitPassword = computed(() => strength.value.score >= 3 && password.new_password === confirmPassword.value && password.current_password)

async function submitPassword() {
  if (!canSubmitPassword.value) return
  savingPassword.value = true
  passwordMessage.value = ''
  try {
    await updatePassword(password)
    password.current_password = ''
    password.new_password = ''
    confirmPassword.value = ''
    passwordMessage.value = 'Password updated successfully.'
  } catch (err) {
    passwordMessage.value = err.response?.data?.error?.message || 'Password update failed.'
  } finally {
    savingPassword.value = false
  }
}

async function loadSessions() {
  const res = await listSessions()
  sessions.value = res.data?.data || []
}

async function revoke(id) {
  await revokeSession(id)
  await loadSessions()
}

async function clearOthers() {
  await revokeOtherSessions()
  await loadSessions()
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString() : 'Unknown'
}

function deviceLabel(ua = '') {
  if (/edg/i.test(ua)) return 'Microsoft Edge'
  if (/chrome/i.test(ua)) return 'Chrome'
  if (/firefox/i.test(ua)) return 'Firefox'
  if (/safari/i.test(ua)) return 'Safari'
  if (/curl|postman|insomnia/i.test(ua)) return 'API Client'
  return 'Unknown device'
}

function locationLabel(session) {
  return session.geo_country || session.location || 'Location pending'
}

onMounted(loadSessions)
</script>

<style scoped>
.security-form{display:grid;gap:.8rem}.security-form label{display:grid;gap:.35rem;color:#475569;font-size:.85rem;font-weight:800}.security-form input{height:42px;border:1px solid #cbd5e1;border-radius:6px;padding:0 .75rem;color:#111827}.security-form button,.session-card button,.secondary-action{height:38px;border:0;border-radius:6px;background:#112b4e;color:white;font-weight:800;padding:0 .9rem}.security-form button:disabled{opacity:.55}.strength-track{height:8px;background:#e5e7eb;border-radius:99px;overflow:hidden}.strength-track span{display:block;height:100%;transition:.2s}.validation-copy,.status-copy,.empty-copy{margin:0;color:#64748b}.panel-title-row{display:flex;align-items:center;justify-content:space-between;gap:1rem;margin-bottom:1rem}.panel-title-row h3{margin:0}.secondary-action{background:#f48c06}.session-list{display:grid;gap:.75rem}.session-card{display:flex;align-items:center;justify-content:space-between;gap:1rem;border:1px solid #e5e7eb;border-radius:8px;padding:.85rem}.session-card strong{color:#112b4e}.session-card p{margin:.2rem 0;color:#475569}.session-card small{color:#94a3b8}.current-badge{background:#dcfce7;color:#166534;border-radius:999px;padding:.3rem .55rem;font-size:.75rem;font-weight:900;white-space:nowrap}@media(max-width:720px){.panel-title-row,.session-card{align-items:flex-start;flex-direction:column}}
</style>
