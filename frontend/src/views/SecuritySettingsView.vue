<template>
  <LayoutDefault title="Security Settings">
    <div class="max-w-4xl mx-auto space-y-6">
      <!-- Toast -->
      <Transition name="toast">
        <div v-if="toast.msg" :class="toast.ok ? 'bg-gray-900 dark:bg-slate-700' : 'bg-red-600'" class="fixed top-5 right-5 z-[100] px-4 py-3 text-white text-sm rounded-xl shadow-xl">{{ toast.msg }}</div>
      </Transition>

      <!-- ═════ Password ═════ -->
      <div class="admin-card">
        <div class="admin-card-header">
          <div>
            <h3 class="text-sm font-semibold text-gray-800 dark:text-slate-200 uppercase tracking-wider">Password</h3>
            <p class="text-xs text-gray-500 dark:text-slate-400 mt-0.5">Rotate your account login password. Minimum 8 characters required.</p>
          </div>
        </div>
        <div class="admin-card-body">
          <form @submit.prevent="onChangePassword" class="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div>
              <label class="form-label">Current Password</label>
              <input v-model="pw.current" type="password" placeholder="Current password" class="form-input"/>
            </div>
            <div>
              <label class="form-label">New Password</label>
              <input v-model="pw.next" type="password" placeholder="New password (min 8 chars)" class="form-input"/>
            </div>
            <div>
              <label class="form-label">Confirm New Password</label>
              <input v-model="pw.confirm" type="password" placeholder="Confirm new password" class="form-input"/>
            </div>
            <div class="md:col-span-3">
              <button type="submit" :disabled="busy.pw" class="btn-primary disabled:opacity-60">
                <i class="fas fa-key mr-1.5"></i>{{ busy.pw ? 'Updating…' : 'Update Password' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- ═════ PIN ═════ -->
      <div class="admin-card">
        <div class="admin-card-header">
          <div>
            <h3 class="text-sm font-semibold text-gray-800 dark:text-slate-200 uppercase tracking-wider">Security PIN</h3>
            <p class="text-xs text-gray-500 dark:text-slate-400 mt-0.5">A short PIN is required when unlocking sensitive screens after inactivity.</p>
          </div>
          <span :class="isFirstPin ? 'badge-yellow' : 'badge-green'">
            <i :class="isFirstPin ? 'icofont-warning' : 'icofont-check'" class="text-xs mr-1"></i>
            {{ isFirstPin ? 'Not Set' : 'PIN Active' }}
          </span>
        </div>
        <div class="admin-card-body">
          <form @submit.prevent="onChangePin" class="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div v-if="!isFirstPin">
              <label class="form-label">Current PIN</label>
              <input v-model="pin.current" type="password" placeholder="Current PIN" maxlength="6" class="form-input"/>
            </div>
            <div>
              <label class="form-label">{{ isFirstPin ? 'Choose a PIN' : 'New PIN' }}</label>
              <input v-model="pin.next" type="password" :placeholder="isFirstPin ? '4–6 digit PIN' : 'New PIN'" maxlength="6" class="form-input"/>
            </div>
            <div>
              <label class="form-label">Confirm PIN</label>
              <input v-model="pin.confirm" type="password" placeholder="Confirm PIN" maxlength="6" class="form-input"/>
            </div>
            <div class="md:col-span-3">
              <button type="submit" :disabled="busy.pin" class="btn-primary disabled:opacity-60">
                <i class="fas fa-lock mr-1.5"></i>{{ busy.pin ? 'Saving…' : (isFirstPin ? 'Set PIN' : 'Update PIN') }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- ═════ IP Allowlist ═════ -->
      <div class="admin-card">
        <div class="admin-card-header">
          <div>
            <h3 class="text-sm font-semibold text-gray-800 dark:text-slate-200 uppercase tracking-wider">IP Allowlist</h3>
            <p class="text-xs text-gray-500 dark:text-slate-400 mt-0.5">When entries exist, requests from any unlisted IP are blocked with 403 Forbidden.</p>
          </div>
          <span v-if="ipList.length" class="badge-red"><i class="icofont-shield-alt text-xs mr-1"></i>Active</span>
          <span v-else class="badge-gray">No restriction</span>
        </div>
        <div class="admin-card-body space-y-4">
          <form @submit.prevent="onAddIp" class="grid grid-cols-1 md:grid-cols-12 gap-2">
            <div class="md:col-span-5">
              <label class="form-label">IP Address or CIDR</label>
              <input v-model="newIp.value" placeholder="41.74.32.5 or 192.168.0.0/24" class="form-input font-mono"/>
            </div>
            <div class="md:col-span-5">
              <label class="form-label">Label</label>
              <input v-model="newIp.label" placeholder="e.g. Home, Office VPN" class="form-input"/>
            </div>
            <div class="md:col-span-2 flex items-end">
              <button type="submit" :disabled="busy.ip" class="btn-primary disabled:opacity-60 w-full">
                {{ busy.ip ? 'Adding…' : 'Add Entry' }}
              </button>
            </div>
          </form>

          <ul v-if="ipList.length" class="divide-y divide-gray-100 dark:divide-slate-800 border border-gray-100 dark:border-slate-800 rounded-lg overflow-hidden">
            <li v-for="entry in ipList" :key="entry.id" class="flex items-center gap-3 px-4 py-2.5 bg-white dark:bg-slate-900">
              <code class="font-mono text-sm text-gray-800 dark:text-slate-200 bg-gray-50 dark:bg-slate-800 px-2 py-1 rounded border border-gray-200 dark:border-slate-700">{{ entry.ip_or_cidr }}</code>
              <span class="text-xs text-gray-500 dark:text-slate-400">{{ entry.label || '—' }}</span>
              <span class="ml-auto text-xs text-gray-400 dark:text-slate-500">{{ formatDate(entry.created_at) }}</span>
              <button @click="removeIp(entry.id)" class="text-gray-400 hover:text-red-600 dark:hover:text-red-400 text-lg leading-none transition-colors" title="Remove">×</button>
            </li>
          </ul>
          <p v-else class="text-xs text-gray-400 dark:text-slate-500 italic">No entries — account accepts logins from any IP address.</p>
        </div>
      </div>

      <!-- ═════ 2FA ═════ -->
      <div class="admin-card">
        <div class="admin-card-header">
          <div>
            <h3 class="text-sm font-semibold text-gray-800 dark:text-slate-200 uppercase tracking-wider">Two-Factor Authentication</h3>
            <p class="text-xs text-gray-500 dark:text-slate-400 mt-0.5">Use Google Authenticator, 1Password, or Authy to generate one-time codes.</p>
          </div>
          <span v-if="twofa.enabled" class="badge-green"><i class="icofont-check text-xs mr-1"></i>Enabled</span>
          <span v-else class="badge-gray">Disabled</span>
        </div>
        <div class="admin-card-body">
          <div v-if="twofa.enabled" class="space-y-3">
            <p class="text-sm text-gray-700 dark:text-slate-300">Two-factor authentication is active for your account.</p>
            <button @click="onDisable2FA" :disabled="busy.twofa" class="btn-primary disabled:opacity-60 bg-red-600 hover:bg-red-700">
              <i class="icofont-lock mr-1.5"></i>{{ busy.twofa ? 'Disabling…' : 'Disable 2FA' }}
            </button>
          </div>

          <div v-else-if="!enrollment">
            <button @click="onEnroll" :disabled="busy.twofa" class="btn-primary disabled:opacity-60">
              <i class="icofont-qr-code mr-1.5"></i>{{ busy.twofa ? 'Generating…' : 'Set up 2FA' }}
            </button>
          </div>

          <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="text-center">
              <img :src="qrSrc(enrollment.otpauth_url)" alt="Scan to add NCS account" class="mx-auto rounded-lg border border-gray-200 dark:border-slate-700 bg-white p-2"/>
              <p class="text-xs text-gray-400 dark:text-slate-500 mt-2">Scan with your authenticator app</p>
            </div>
            <div class="space-y-4">
              <div>
                <label class="form-label">Or enter this secret manually</label>
                <code class="block font-mono text-sm bg-gray-50 dark:bg-slate-800 border border-gray-200 dark:border-slate-700 text-gray-800 dark:text-slate-200 rounded-lg px-3 py-2.5 break-all">{{ enrollment.secret }}</code>
              </div>
              <form @submit.prevent="onVerify2FA" class="space-y-3">
                <div>
                  <label class="form-label">6-digit code from your app</label>
                  <input v-model="verifyCode" maxlength="6" inputmode="numeric" placeholder="123456" class="form-input font-mono text-center tracking-widest"/>
                </div>
                <button type="submit" :disabled="busy.twofa" class="btn-primary disabled:opacity-60 w-full justify-center flex">
                  {{ busy.twofa ? 'Verifying…' : 'Enable 2FA' }}
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>
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
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(-8px); }
.toast-enter-active, .toast-leave-active { transition: all .25s ease; }
</style>
