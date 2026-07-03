<template>
  <section class="profile-settings">
    <div class="cms-panel-head"><h2>Profile Settings</h2></div>

    <div class="cms-panel profile-identity">
      <div class="cms-panel-head"><h3>Display Picture &amp; Name</h3></div>
      <div class="profile-avatar-row">
        <span class="profile-avatar">
          <img v-if="avatarPreview" :src="avatarPreview" alt="" />
          <span v-else class="profile-avatar-initials">{{ initials }}</span>
        </span>
        <div class="profile-avatar-actions">
          <label class="profile-avatar-upload">
            {{ uploadingAvatar ? 'Uploading…' : 'Change picture' }}
            <input type="file" accept="image/png,image/jpeg,image/webp" :disabled="uploadingAvatar" @change="onAvatarPicked" />
          </label>
          <button v-if="form.avatar_url" type="button" class="profile-avatar-remove" :disabled="uploadingAvatar" @click="removeAvatar">Remove</button>
        </div>
      </div>
      <div class="cms-two">
        <label>First name<input v-model="form.first_name" class="form-control" required /></label>
        <label>Last name<input v-model="form.last_name" class="form-control" required /></label>
        <label>Email <small>(read-only)</small><input :value="email" class="form-control" disabled /></label>
      </div>
      <div class="profile-actions">
        <button type="button" :disabled="savingProfile" @click="saveProfile">{{ savingProfile ? 'Saving…' : 'Save profile' }}</button>
      </div>
    </div>

    <div class="cms-panel">
      <div class="cms-panel-head"><h3>Change Password</h3></div>
      <div class="cms-two">
        <label>Current password<input v-model="password.current_password" type="password" autocomplete="current-password" class="form-control" /></label>
        <label>New password<input v-model="password.new_password" type="password" autocomplete="new-password" class="form-control" /></label>
        <label>Confirm new password<input v-model="confirmPassword" type="password" autocomplete="new-password" class="form-control" /></label>
      </div>
      <div class="strength-track"><span :style="{ width: `${strength.score * 25}%`, background: strength.color }"></span></div>
      <p class="profile-hint">{{ strength.label }}</p>
      <div class="profile-actions">
        <button type="button" :disabled="savingPassword || !canSubmitPassword" @click="submitPassword">{{ savingPassword ? 'Updating…' : 'Update password' }}</button>
      </div>
    </div>

    <div class="cms-panel">
      <div class="cms-panel-head">
        <h3>Two-Factor Authentication</h3>
        <span class="storage-badge" :class="{ ok: twofa.enabled }">{{ twofa.enabled ? 'Enabled' : 'Not enabled' }}</span>
      </div>

      <div v-if="!twofa.enabled && !enrollment">
        <p class="profile-hint">Add an authenticator app (Google Authenticator, Authy, 1Password) as a second sign-in factor.</p>
        <div class="profile-actions"><button type="button" :disabled="enrolling" @click="startEnrollment">{{ enrolling ? 'Starting…' : 'Enable 2FA' }}</button></div>
      </div>

      <div v-else-if="enrollment" class="cms-two">
        <label class="wide">Secret key <small>(enter manually in your authenticator app)</small><input :value="enrollment.secret" class="form-control" readonly @click="$event.target.select()" /></label>
        <label class="wide">Setup link <small>(otpauth:// URI)</small><input :value="enrollment.otpauth_url" class="form-control" readonly @click="$event.target.select()" /></label>
        <label>6-digit code<input v-model="verifyCode" class="form-control" inputmode="numeric" maxlength="6" placeholder="123456" /></label>
        <div class="profile-actions">
          <button type="button" :disabled="verifying || verifyCode.length !== 6" @click="finishEnrollment">{{ verifying ? 'Verifying…' : 'Verify & enable' }}</button>
          <button type="button" class="profile-cancel" @click="enrollment = null">Cancel</button>
        </div>
      </div>

      <div v-else>
        <p class="profile-hint">Two-factor authentication is active on your account.</p>
        <div class="profile-actions"><button type="button" class="font-delete" :disabled="disabling" @click="disable">{{ disabling ? 'Disabling…' : 'Disable 2FA' }}</button></div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { updateMyProfile } from '@/api/auth.js'
import { uploadMedia } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'
import { changePassword, disable2FA, enroll2FA, getMySecurity, verify2FA } from '@/api/security.js'

const emit = defineEmits(['message', 'error', 'profile-updated'])

function readStoredUser() {
  try { return JSON.parse(localStorage.getItem('ncsms_user') || '{}') } catch { return {} }
}

const email = ref('')
const form = reactive({ first_name: '', last_name: '', avatar_url: '' })
const savingProfile = ref(false)
const uploadingAvatar = ref(false)

const avatarPreview = computed(() => form.avatar_url ? mediaUrl(form.avatar_url) : '')
const initials = computed(() => `${form.first_name[0] || ''}${form.last_name[0] || ''}`.toUpperCase() || 'U')

function loadProfile() {
  const u = readStoredUser()
  email.value = u.email || ''
  form.first_name = u.first_name || ''
  form.last_name = u.last_name || ''
  form.avatar_url = u.avatar_url || ''
}

async function persistProfile() {
  const res = await updateMyProfile({ first_name: form.first_name.trim(), last_name: form.last_name.trim(), avatar_url: form.avatar_url })
  const updated = res.data?.data
  if (updated) {
    localStorage.setItem('ncsms_user', JSON.stringify({ ...readStoredUser(), ...updated }))
    emit('profile-updated', updated)
  }
}

async function saveProfile() {
  if (!form.first_name.trim() || !form.last_name.trim()) {
    emit('error', new Error('First and last name are required'))
    return
  }
  savingProfile.value = true
  try {
    await persistProfile()
    emit('message', 'Profile saved')
  } catch (err) {
    emit('error', err)
  } finally {
    savingProfile.value = false
  }
}

async function onAvatarPicked(event) {
  const file = event.target.files?.[0]
  event.target.value = ''
  if (!file) return
  uploadingAvatar.value = true
  try {
    const res = await uploadMedia(file)
    form.avatar_url = res.data?.data?.url || ''
    await persistProfile()
    emit('message', 'Profile picture updated')
  } catch (err) {
    emit('error', err)
  } finally {
    uploadingAvatar.value = false
  }
}

async function removeAvatar() {
  form.avatar_url = ''
  try {
    await persistProfile()
    emit('message', 'Profile picture removed')
  } catch (err) { emit('error', err) }
}

// ── Password ──────────────────────────────────────────────────────
const password = reactive({ current_password: '', new_password: '' })
const confirmPassword = ref('')
const savingPassword = ref(false)

const strength = computed(() => {
  const value = password.new_password
  let score = 0
  if (value.length >= 8) score++
  if (/[A-Z]/.test(value) && /[a-z]/.test(value)) score++
  if (/\d/.test(value)) score++
  if (/[^A-Za-z0-9]/.test(value)) score++
  const labels = ['Too weak', 'Weak', 'Fair', 'Strong', 'Excellent']
  const colors = ['#dc2626', '#ea580c', '#d97706', '#16a34a', '#047857']
  return { score, label: value ? labels[score] : '', color: colors[score] }
})
const canSubmitPassword = computed(() => strength.value.score >= 3 && password.new_password === confirmPassword.value && password.current_password)

async function submitPassword() {
  if (!canSubmitPassword.value) return
  savingPassword.value = true
  try {
    await changePassword(password.current_password, password.new_password)
    password.current_password = ''
    password.new_password = ''
    confirmPassword.value = ''
    emit('message', 'Password updated')
  } catch (err) {
    emit('error', err)
  } finally {
    savingPassword.value = false
  }
}

// ── Two-factor authentication ───────────────────────────────────────
const twofa = reactive({ enabled: false, has_pending: false })
const enrollment = ref(null)
const verifyCode = ref('')
const enrolling = ref(false)
const verifying = ref(false)
const disabling = ref(false)

async function loadSecurity() {
  try {
    const data = await getMySecurity()
    Object.assign(twofa, data?.twofa || {})
  } catch (err) { emit('error', err) }
}

async function startEnrollment() {
  enrolling.value = true
  try {
    enrollment.value = await enroll2FA()
  } catch (err) {
    emit('error', err)
  } finally {
    enrolling.value = false
  }
}

async function finishEnrollment() {
  verifying.value = true
  try {
    await verify2FA(verifyCode.value)
    verifyCode.value = ''
    enrollment.value = null
    await loadSecurity()
    emit('message', 'Two-factor authentication enabled')
  } catch (err) {
    emit('error', err)
  } finally {
    verifying.value = false
  }
}

async function disable() {
  disabling.value = true
  try {
    await disable2FA()
    await loadSecurity()
    emit('message', 'Two-factor authentication disabled')
  } catch (err) {
    emit('error', err)
  } finally {
    disabling.value = false
  }
}

onMounted(() => {
  loadProfile()
  loadSecurity()
})
</script>

<style scoped>
/*
 * This component is rendered as a child of WebsiteContentManagerView, whose
 * <style scoped> only reaches this component's root element (Vue scoped CSS
 * does not cascade into child templates) — so shared classes are redefined
 * here rather than assumed to be inherited.
 */
.profile-settings { display: grid; gap: 20px; min-width: 0; }

.cms-panel-head { display: flex; justify-content: space-between; align-items: center; gap: 16px; flex-wrap: wrap; border-bottom: 1px solid #f9f9f9; padding-bottom: 15px; }
.cms-panel-head h2,
.cms-panel-head h3 { margin: 0; font-size: 16px; font-weight: 700; color: #34395e; }

.cms-panel { background: white; border: 0; border-radius: 3px; box-shadow: 0 4px 25px rgba(0,0,0,.1); padding: 25px; display: grid; gap: 18px; min-width: 0; }

.cms-two { display: grid; grid-template-columns: repeat(auto-fit, minmax(15rem, 1fr)); gap: 18px; }
.cms-two label { display: grid; gap: 7px; font-size: 12px; font-weight: 600; color: #34395e; min-width: 0; }
.cms-two label small { font-weight: 500; text-transform: none; color: #94a3b8; }
.cms-two .wide { grid-column: 1 / -1; }
.cms-two input { width: 100%; box-sizing: border-box; border: 1px solid #e4e6fc; border-radius: 3px; padding: 10px 15px; font-weight: 500; color: #495057; background: #fdfdff; outline: none; }
.cms-two input:focus { border-color: #6777ef; box-shadow: 0 2px 6px #acb5f6; }
.cms-two input:disabled { background: #f4f6f9; color: #94a3b8; }

.profile-avatar-row { display: flex; align-items: center; gap: 18px; }
.profile-avatar { width: 72px; height: 72px; border-radius: 50%; overflow: hidden; background: #6777ef; display: grid; place-items: center; flex-shrink: 0; }
.profile-avatar img { width: 100%; height: 100%; object-fit: cover; }
.profile-avatar-initials { color: #fff; font-weight: 800; font-size: 22px; }
.profile-avatar-actions { display: flex; flex-direction: column; gap: 8px; }
.profile-avatar-upload { position: relative; border: 1px solid #e4e6fc; border-radius: 30px; background: #fff; color: #34395e; padding: 8px 16px; font-size: 12px; font-weight: 700; cursor: pointer; text-align: center; }
.profile-avatar-upload input { position: absolute; inset: 0; opacity: 0; cursor: pointer; }
.profile-avatar-remove { border: 0; border-radius: 30px; background: #fc544b; color: #fff; padding: 8px 16px; font-size: 12px; font-weight: 700; }
.profile-avatar-remove:disabled { opacity: .5; }

.profile-actions { display: flex; flex-wrap: wrap; gap: 10px; }
.profile-actions button { border: 0; border-radius: 30px; background: #6777ef; color: white; padding: 8px 18px; font-size: 12px; font-weight: 600; white-space: nowrap; box-shadow: 0 2px 6px #acb5f6; }
.profile-actions button:disabled { opacity: .6; }
.profile-cancel { background: #94a3b8 !important; box-shadow: none !important; }
.font-delete { background: #fc544b !important; box-shadow: 0 2px 6px #fd9b96 !important; }

.strength-track { height: 8px; background: #e5e7eb; border-radius: 99px; overflow: hidden; }
.strength-track span { display: block; height: 100%; transition: .2s; }
.profile-hint { color: #6c757d; font-size: 13px; margin: 0; line-height: 1.6; }

.storage-badge { font-size: 12px; font-weight: 700; padding: 5px 10px; border-radius: 30px; background: #fdeaea; color: #fc544b; white-space: nowrap; }
.storage-badge.ok { background: #e8f7f0; color: #47c363; }

:global(.dark .cms-panel-head h2),
:global(.dark .cms-panel-head h3) { color: #f8fafc; }
:global(.dark .cms-panel) { background: #111827; border-color: #334155; }
:global(.dark .cms-two input) { background: #0f172a; color: #f8fafc; border-color: #475569; }
:global(.dark .profile-avatar-upload) { background: #0f172a; color: #f8fafc; border-color: #475569; }
:global(.dark .profile-hint) { color: #cbd5e1; }
</style>
