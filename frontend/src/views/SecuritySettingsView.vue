<template>
  <LayoutDefault title="Security Settings">
    <div class="p-6 max-w-4xl mx-auto space-y-6">
      <!-- Toast -->
      <Transition name="toast">
        <div v-if="toast.msg" :class="toast.ok ? 'bg-gray-900' : 'bg-red-600'" class="fixed top-5 right-5 z-[100] px-4 py-3 text-white text-sm rounded-xl shadow-xl">{{ toast.msg }}</div>
      </Transition>

      <header>
        <h2 class="text-2xl font-semibold text-gray-900">Security Settings</h2>
        <p class="text-xs text-gray-500 mt-1">Rotate credentials, restrict access by IP, and enable two-factor authentication for your account.</p>
      </header>

      <!-- ═════ Password ═════ -->
      <section class="bg-white border border-gray-200 rounded-xl p-5">
        <h3 class="text-sm font-semibold text-gray-700 uppercase tracking-wider mb-4">Password</h3>
        <form @submit.prevent="onChangePassword" class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <input v-model="pw.current" type="password" placeholder="Current password" class="input"/>
          <input v-model="pw.next" type="password" placeholder="New password" class="input"/>
          <input v-model="pw.confirm" type="password" placeholder="Confirm new password" class="input"/>
          <button type="submit" :disabled="busy.pw" class="md:col-span-3 justify-self-start btn-primary">{{ busy.pw ? 'Updating…' : 'Update password' }}</button>
        </form>
      </section>

      <!-- ═════ PIN ═════ -->
      <section class="bg-white border border-gray-200 rounded-xl p-5">
        <h3 class="text-sm font-semibold text-gray-700 uppercase tracking-wider mb-4">Security PIN</h3>
        <p class="text-xs text-gray-500 mb-3">A short PIN is required when unlocking sensitive screens.</p>
        <form @submit.prevent="onChangePin" class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <input v-if="!isFirstPin" v-model="pin.current" type="password" placeholder="Current PIN" maxlength="6" class="input"/>
          <input v-model="pin.next" type="password" :placeholder="isFirstPin ? 'Choose a PIN' : 'New PIN'" maxlength="6" class="input"/>
          <input v-model="pin.confirm" type="password" placeholder="Confirm PIN" maxlength="6" class="input"/>
          <button type="submit" :disabled="busy.pin" class="md:col-span-3 justify-self-start btn-primary">{{ busy.pin ? 'Saving…' : (isFirstPin ? 'Set PIN' : 'Update PIN') }}</button>
        </form>
      </section>

      <!-- ═════ IP Allowlist ═════ -->
      <section class="bg-white border border-gray-200 rounded-xl p-5">
        <div class="flex items-start justify-between mb-3">
          <div>
            <h3 class="text-sm font-semibold text-gray-700 uppercase tracking-wider">IP Allowlist</h3>
            <p class="text-xs text-gray-500 mt-1">When at least one address is listed, requests from any other IP will be blocked with 403 Forbidden.</p>
          </div>
          <span v-if="ipList.length" class="text-[10px] uppercase tracking-wider font-bold bg-red-100 text-red-700 px-2 py-1 rounded-full">Active</span>
          <span v-else class="text-[10px] uppercase tracking-wider font-bold bg-gray-100 text-gray-600 px-2 py-1 rounded-full">No restriction</span>
        </div>

        <form @submit.prevent="onAddIp" class="grid grid-cols-1 md:grid-cols-12 gap-2 mb-4">
          <input v-model="newIp.value" placeholder="41.74.32.5 or 192.168.0.0/24"
            class="md:col-span-5 input font-mono"/>
          <input v-model="newIp.label" placeholder="Label (e.g. Home, Office VPN)"
            class="md:col-span-5 input"/>
          <button type="submit" :disabled="busy.ip" class="md:col-span-2 btn-primary">{{ busy.ip ? 'Adding…' : 'Add' }}</button>
        </form>

        <ul v-if="ipList.length" class="divide-y divide-gray-100">
          <li v-for="entry in ipList" :key="entry.id" class="flex items-center gap-3 py-2">
            <code class="font-mono text-sm text-gray-800 bg-gray-50 px-2 py-1 rounded">{{ entry.ip_or_cidr }}</code>
            <span class="text-xs text-gray-500">{{ entry.label || '—' }}</span>
            <span class="ml-auto text-xs text-gray-400">{{ formatDate(entry.created_at) }}</span>
            <button @click="removeIp(entry.id)" class="text-gray-400 hover:text-red-600 text-lg leading-none" title="Remove">×</button>
          </li>
        </ul>
        <p v-else class="text-xs text-gray-400 italic">No entries — your account accepts logins from any IP.</p>
      </section>

      <!-- ═════ 2FA ═════ -->
      <section class="bg-white border border-gray-200 rounded-xl p-5">
        <div class="flex items-start justify-between mb-3">
          <div>
            <h3 class="text-sm font-semibold text-gray-700 uppercase tracking-wider">Two-Factor Authentication</h3>
            <p class="text-xs text-gray-500 mt-1">Use an authenticator app (Google Authenticator, 1Password, Authy) to generate one-time codes.</p>
          </div>
          <span v-if="twofa.enabled" class="text-[10px] uppercase tracking-wider font-bold bg-green-100 text-green-700 px-2 py-1 rounded-full">Enabled</span>
          <span v-else class="text-[10px] uppercase tracking-wider font-bold bg-gray-100 text-gray-600 px-2 py-1 rounded-full">Disabled</span>
        </div>

        <div v-if="twofa.enabled" class="space-y-3">
          <p class="text-sm text-gray-700">Two-factor authentication is active for your account.</p>
          <button @click="onDisable2FA" :disabled="busy.twofa" class="btn-danger">{{ busy.twofa ? 'Disabling…' : 'Disable 2FA' }}</button>
        </div>

        <div v-else-if="!enrollment">
          <button @click="onEnroll" :disabled="busy.twofa" class="btn-primary">{{ busy.twofa ? 'Generating…' : 'Set up 2FA' }}</button>
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <div class="text-center">
            <img :src="qrSrc(enrollment.otpauth_url)" alt="Scan to add NCS account" class="mx-auto rounded-lg border border-gray-200 bg-white p-2"/>
            <p class="text-[10px] text-gray-400 mt-1">Scan with your authenticator app</p>
          </div>
          <div class="space-y-3">
            <div>
              <label class="block text-[11px] uppercase tracking-wider text-gray-400 mb-1">Or enter this secret manually</label>
              <code class="block font-mono text-sm bg-gray-50 border border-gray-200 rounded px-3 py-2 break-all">{{ enrollment.secret }}</code>
            </div>
            <form @submit.prevent="onVerify2FA" class="space-y-2">
              <label class="block text-[11px] uppercase tracking-wider text-gray-400">Enter the 6-digit code from your app</label>
              <input v-model="verifyCode" maxlength="6" inputmode="numeric" placeholder="123456"
                class="input font-mono text-center tracking-widest"/>
              <button type="submit" :disabled="busy.twofa" class="btn-primary w-full">{{ busy.twofa ? 'Verifying…' : 'Enable 2FA' }}</button>
            </form>
          </div>
        </div>
      </section>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { useAuthStore } from '@/stores/auth.js'
import {
  getMySecurity, addIpWhitelist, removeIpWhitelist,
  enroll2FA, verify2FA, disable2FA, changePassword, setPin, changePin,
} from '@/api/security.js'

const auth = useAuthStore()
const isFirstPin = computed(() => !auth.user?.has_pin)

const toast = reactive({ msg: '', ok: true })
const busy = reactive({ pw: false, pin: false, ip: false, twofa: false })
const pw = reactive({ current: '', next: '', confirm: '' })
const pin = reactive({ current: '', next: '', confirm: '' })
const newIp = reactive({ value: '', label: '' })
const ipList = ref([])
const twofa = reactive({ enabled: false, enabled_at: null, has_pending: false })
const enrollment = ref(null)
const verifyCode = ref('')

function showToast(msg, ok = true) {
  toast.msg = msg; toast.ok = ok
  setTimeout(() => (toast.msg = ''), 3000)
}
function formatDate(d) { return d ? new Date(d).toLocaleString() : '' }

function qrSrc(otpauth) {
  // Render via a public QR endpoint so we don't ship a QR lib.
  return `https://api.qrserver.com/v1/create-qr-code/?size=180x180&data=${encodeURIComponent(otpauth)}`
}

async function load() {
  try {
    const s = await getMySecurity()
    ipList.value = s.ip_whitelist || []
    Object.assign(twofa, s.twofa || {})
  } catch (e) { /* leave defaults */ }
}

async function onChangePassword() {
  if (!pw.current || !pw.next) return showToast('Enter current and new password', false)
  if (pw.next !== pw.confirm) return showToast('New passwords do not match', false)
  busy.pw = true
  try {
    await changePassword(pw.current, pw.next)
    pw.current = ''; pw.next = ''; pw.confirm = ''
    showToast('Password updated')
  } catch (e) {
    showToast(e?.response?.data?.message || 'Password update failed', false)
  } finally { busy.pw = false }
}

async function onChangePin() {
  if (!pin.next || pin.next.length < 4) return showToast('PIN must be at least 4 digits', false)
  if (pin.next !== pin.confirm) return showToast('PINs do not match', false)
  busy.pin = true
  try {
    if (isFirstPin.value) {
      await setPin(pin.next)
    } else {
      if (!pin.current) { showToast('Enter your current PIN', false); busy.pin = false; return }
      await changePin(pin.current, pin.next)
    }
    pin.current = ''; pin.next = ''; pin.confirm = ''
    showToast('PIN saved')
  } catch (e) {
    showToast(e?.response?.data?.message || 'PIN update failed', false)
  } finally { busy.pin = false }
}

async function onAddIp() {
  if (!newIp.value) return
  busy.ip = true
  try {
    await addIpWhitelist(newIp.value, newIp.label)
    newIp.value = ''; newIp.label = ''
    await load()
    showToast('IP entry added')
  } catch (e) {
    showToast(e?.response?.data?.message || 'Add failed', false)
  } finally { busy.ip = false }
}

async function removeIp(id) {
  try {
    await removeIpWhitelist(id)
    await load()
    showToast('Removed')
  } catch (e) { showToast('Remove failed', false) }
}

async function onEnroll() {
  busy.twofa = true
  try {
    enrollment.value = await enroll2FA()
    verifyCode.value = ''
  } catch (e) {
    showToast('Could not start enrollment', false)
  } finally { busy.twofa = false }
}

async function onVerify2FA() {
  if (!verifyCode.value || verifyCode.value.length !== 6) return showToast('Enter the 6-digit code', false)
  busy.twofa = true
  try {
    await verify2FA(verifyCode.value)
    enrollment.value = null; verifyCode.value = ''
    await load()
    showToast('Two-factor authentication enabled')
  } catch (e) {
    showToast(e?.response?.data?.message || 'Code did not match', false)
  } finally { busy.twofa = false }
}

async function onDisable2FA() {
  if (!confirm('Disable two-factor authentication?')) return
  busy.twofa = true
  try {
    await disable2FA()
    enrollment.value = null
    await load()
    showToast('2FA disabled')
  } catch (e) { showToast('Disable failed', false) }
  finally { busy.twofa = false }
}

onMounted(load)
</script>

<style scoped>
.input { @apply w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 bg-white; }
.btn-primary { @apply text-sm font-medium px-4 py-2 bg-primary-600 hover:bg-primary-700 disabled:opacity-60 text-white rounded-lg; }
.btn-danger { @apply text-sm font-medium px-4 py-2 bg-red-600 hover:bg-red-700 disabled:opacity-60 text-white rounded-lg; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(-8px); }
.toast-enter-active, .toast-leave-active { transition: all .25s ease; }
</style>
