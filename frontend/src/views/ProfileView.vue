<template>
  <LayoutDefault title="My Profile">
    <div class="max-w-3xl mx-auto space-y-5">

      <!-- Profile header card -->
      <div class="admin-card overflow-hidden">
        <div class="bg-gradient-to-r from-primary-700 to-primary-500 px-6 py-7">
          <div class="flex items-center gap-5">
            <div class="w-20 h-20 bg-white rounded-2xl flex items-center justify-center shadow-xl flex-shrink-0">
              <span class="text-primary-700 text-3xl font-bold">{{ userInitials }}</span>
            </div>
            <div class="min-w-0">
              <h2 class="text-xl font-bold text-white truncate">{{ fullName }}</h2>
              <p class="text-primary-200 text-sm mt-0.5">{{ user?.email }}</p>
              <div class="flex flex-wrap gap-1.5 mt-2.5">
                <span v-for="role in userRoles" :key="role" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-white/20 text-white">
                  {{ role }}
                </span>
                <span :class="user?.is_active !== false ? 'bg-green-400/30 text-green-100' : 'bg-gray-400/30 text-gray-200'" class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium">
                  <span class="w-1.5 h-1.5 rounded-full" :class="user?.is_active !== false ? 'bg-green-400' : 'bg-gray-400'"></span>
                  {{ user?.is_active !== false ? 'Active' : 'Inactive' }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Quick stats row -->
        <div class="grid grid-cols-3 divide-x divide-gray-100 dark:divide-slate-800 border-t border-gray-100 dark:border-slate-800">
          <div class="px-5 py-4 text-center">
            <div class="text-xs text-gray-400 dark:text-slate-500 uppercase tracking-wide mb-0.5">Last Login</div>
            <div class="text-sm font-medium text-gray-700 dark:text-slate-300">{{ formatDate(user?.last_login_at) }}</div>
          </div>
          <div class="px-5 py-4 text-center">
            <div class="text-xs text-gray-400 dark:text-slate-500 uppercase tracking-wide mb-0.5">PIN Status</div>
            <div class="text-sm font-medium" :class="hasPin ? 'text-green-600 dark:text-green-400' : 'text-yellow-600 dark:text-yellow-400'">
              {{ hasPin ? 'PIN Set' : 'No PIN' }}
            </div>
          </div>
          <div class="px-5 py-4 text-center">
            <div class="text-xs text-gray-400 dark:text-slate-500 uppercase tracking-wide mb-0.5">User ID</div>
            <div class="text-xs font-mono text-gray-500 dark:text-slate-400">{{ (user?.id || user?.user_id)?.substring(0,12) }}…</div>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="flex items-center gap-1 bg-white dark:bg-slate-900 border border-gray-200 dark:border-slate-800 rounded-xl p-1">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="activeTab === tab.id ? 'bg-primary-700 text-white shadow-sm' : 'text-gray-600 dark:text-slate-400 hover:bg-gray-100 dark:hover:bg-slate-800'"
          class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all"
        >
          <i :class="[tab.icon, 'text-base leading-none']"></i>
          {{ tab.label }}
        </button>
      </div>

      <!-- Tab: Account -->
      <div v-if="activeTab === 'account'" class="admin-card">
        <div class="admin-card-header">
          <span class="text-sm font-semibold text-gray-800 dark:text-slate-200">Account Information</span>
        </div>
        <div class="admin-card-body">
          <dl class="grid grid-cols-1 sm:grid-cols-2 gap-5">
            <div>
              <dt class="form-label">Email Address</dt>
              <dd class="text-sm font-medium text-gray-900 dark:text-slate-100">{{ user?.email || 'N/A' }}</dd>
            </div>
            <div>
              <dt class="form-label">Phone Number</dt>
              <dd class="text-sm font-medium text-gray-900 dark:text-slate-100">{{ user?.phone || 'Not provided' }}</dd>
            </div>
            <div>
              <dt class="form-label">First Name</dt>
              <dd class="text-sm font-medium text-gray-900 dark:text-slate-100">{{ user?.first_name || 'N/A' }}</dd>
            </div>
            <div>
              <dt class="form-label">Last Name</dt>
              <dd class="text-sm font-medium text-gray-900 dark:text-slate-100">{{ user?.last_name || 'N/A' }}</dd>
            </div>
            <div>
              <dt class="form-label">Roles</dt>
              <dd class="flex flex-wrap gap-1 mt-1">
                <span v-for="role in userRoles" :key="role" class="badge-blue">{{ role }}</span>
                <span v-if="!userRoles.length" class="badge-gray">No roles</span>
              </dd>
            </div>
            <div>
              <dt class="form-label">Account Status</dt>
              <dd class="mt-1">
                <span :class="user?.is_active !== false ? 'badge-green' : 'badge-red'">
                  {{ user?.is_active !== false ? 'Active' : 'Inactive' }}
                </span>
              </dd>
            </div>
          </dl>
        </div>
      </div>

      <!-- Tab: Security (PIN + Password) -->
      <div v-if="activeTab === 'security'" class="space-y-5">

        <!-- PIN card -->
        <div class="admin-card">
          <div class="admin-card-header">
            <div>
              <p class="text-sm font-semibold text-gray-800 dark:text-slate-200">Screen Lock PIN</p>
              <p class="text-xs text-gray-500 dark:text-slate-400 mt-0.5">Set or change your 4–6 digit PIN used to unlock the app after inactivity.</p>
            </div>
            <span :class="hasPin ? 'badge-green' : 'badge-yellow'">
              <i :class="hasPin ? 'icofont-check' : 'icofont-warning'" class="text-xs mr-1"></i>
              {{ hasPin ? 'PIN Set' : 'Not Set' }}
            </span>
          </div>
          <div class="admin-card-body">
            <div v-if="pinSuccess" class="mb-4 p-3 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg text-sm text-green-700 dark:text-green-400">{{ pinSuccess }}</div>
            <div v-if="pinError" class="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-sm text-red-700 dark:text-red-400">{{ pinError }}</div>

            <div class="flex gap-2 mb-5">
              <button @click="pinMode = 'set'" :disabled="hasPin" :class="[pinMode === 'set' && !hasPin ? 'bg-primary-700 text-white' : 'bg-gray-100 dark:bg-slate-800 text-gray-600 dark:text-slate-400 hover:bg-gray-200 dark:hover:bg-slate-700', hasPin ? 'opacity-40 cursor-not-allowed' : '']" class="px-4 py-1.5 rounded-lg text-sm font-medium transition-colors">Set PIN</button>
              <button @click="pinMode = 'change'" :disabled="!hasPin" :class="[pinMode === 'change' && hasPin ? 'bg-primary-700 text-white' : 'bg-gray-100 dark:bg-slate-800 text-gray-600 dark:text-slate-400 hover:bg-gray-200 dark:hover:bg-slate-700', !hasPin ? 'opacity-40 cursor-not-allowed' : '']" class="px-4 py-1.5 rounded-lg text-sm font-medium transition-colors">Change PIN</button>
            </div>

            <form @submit.prevent="handlePin" class="space-y-4 max-w-sm">
              <div v-if="pinMode === 'change'">
                <label class="form-label">Current PIN</label>
                <input v-model="pinForm.current" type="password" inputmode="numeric" maxlength="6" class="form-input" placeholder="Current PIN" />
              </div>
              <div>
                <label class="form-label">New PIN</label>
                <input v-model="pinForm.pin" type="password" inputmode="numeric" maxlength="6" class="form-input" placeholder="4–6 digit PIN" />
              </div>
              <div>
                <label class="form-label">Confirm PIN</label>
                <input v-model="pinForm.confirm" type="password" inputmode="numeric" maxlength="6" :class="['form-input', pinForm.confirm && pinForm.confirm !== pinForm.pin ? 'border-red-400' : '']" placeholder="Re-enter PIN" />
                <p v-if="pinForm.confirm && pinForm.confirm !== pinForm.pin" class="text-xs text-red-500 mt-1">PINs do not match</p>
              </div>
              <button type="submit" :disabled="pinLoading || pinForm.confirm !== pinForm.pin" class="btn-primary disabled:opacity-60">
                {{ pinLoading ? 'Saving…' : (pinMode === 'set' ? 'Set PIN' : 'Change PIN') }}
              </button>
            </form>
          </div>
        </div>

        <!-- Password card -->
        <div class="admin-card">
          <div class="admin-card-header">
            <div>
              <p class="text-sm font-semibold text-gray-800 dark:text-slate-200">Change Password</p>
              <p class="text-xs text-gray-500 dark:text-slate-400 mt-0.5">Update your account password. Minimum 8 characters.</p>
            </div>
          </div>
          <div class="admin-card-body">
            <div v-if="passwordSuccess" class="mb-4 p-3 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg text-sm text-green-700 dark:text-green-400 flex items-center gap-2">
              <i class="icofont-check-circled text-green-600 dark:text-green-400"></i> Password changed successfully!
            </div>
            <div v-if="passwordError" class="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-sm text-red-700 dark:text-red-400">{{ passwordError }}</div>

            <form @submit.prevent="changePassword" class="space-y-4 max-w-sm">
              <div>
                <label class="form-label">Current Password</label>
                <div class="relative">
                  <input v-model="passwordForm.current_password" :type="showCurrent ? 'text' : 'password'" required class="form-input pr-10" />
                  <button type="button" @click="showCurrent = !showCurrent" class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 dark:text-slate-500 dark:hover:text-slate-300">
                    <i :class="showCurrent ? 'icofont-eye-blocked' : 'icofont-eye'" class="text-base leading-none"></i>
                  </button>
                </div>
              </div>
              <div>
                <label class="form-label">New Password</label>
                <div class="relative">
                  <input v-model="passwordForm.new_password" :type="showNew ? 'text' : 'password'" required minlength="8" class="form-input pr-10" />
                  <button type="button" @click="showNew = !showNew" class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 dark:text-slate-500 dark:hover:text-slate-300">
                    <i :class="showNew ? 'icofont-eye-blocked' : 'icofont-eye'" class="text-base leading-none"></i>
                  </button>
                </div>
              </div>
              <div>
                <label class="form-label">Confirm New Password</label>
                <input v-model="passwordForm.confirm_password" type="password" required :class="['form-input', passwordForm.confirm_password && passwordForm.confirm_password !== passwordForm.new_password ? 'border-red-400' : '']" />
                <p v-if="passwordForm.confirm_password && passwordForm.confirm_password !== passwordForm.new_password" class="text-xs text-red-500 mt-1">Passwords do not match</p>
              </div>
              <button type="submit" :disabled="passwordLoading || passwordForm.confirm_password !== passwordForm.new_password" class="btn-primary disabled:opacity-60 flex items-center gap-2">
                <svg v-if="passwordLoading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                </svg>
                {{ passwordLoading ? 'Updating…' : 'Update Password' }}
              </button>
            </form>
          </div>
        </div>
      </div>

    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { useAuthStore } from '@/stores/auth.js'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'
import apiClient from '@/api/client.js'
import { setPin, changePin } from '@/api/pin.js'

const authStore = useAuthStore()
const breadcrumbStore = useBreadcrumbStore()
const user = computed(() => authStore.user)

const activeTab = ref('account')
const tabs = [
  { id: 'account',  label: 'Account',  icon: 'icofont-user-alt-5' },
  { id: 'security', label: 'Security', icon: 'icofont-lock' }
]

// PIN
const hasPin = computed(() => !!user.value?.has_pin)
const pinMode = ref('set')
watch(hasPin, (v) => { pinMode.value = v ? 'change' : 'set' }, { immediate: true })
const pinForm = ref({ current: '', pin: '', confirm: '' })
const pinLoading = ref(false)
const pinError = ref('')
const pinSuccess = ref('')

async function handlePin() {
  if (pinForm.value.pin !== pinForm.value.confirm) return
  pinLoading.value = true; pinError.value = ''; pinSuccess.value = ''
  try {
    if (pinMode.value === 'set') {
      await setPin(pinForm.value.pin)
      pinSuccess.value = 'PIN set successfully.'
    } else {
      await changePin(pinForm.value.current, pinForm.value.pin)
      pinSuccess.value = 'PIN changed successfully.'
    }
    pinForm.value = { current: '', pin: '', confirm: '' }
    await authStore.refreshUser()
    setTimeout(() => { pinSuccess.value = '' }, 5000)
  } catch (err) {
    pinError.value = err.response?.data?.error?.message || 'Failed to update PIN.'
  } finally { pinLoading.value = false }
}

// Password
const passwordForm = ref({ current_password: '', new_password: '', confirm_password: '' })
const passwordLoading = ref(false)
const passwordError = ref('')
const passwordSuccess = ref(false)
const showCurrent = ref(false)
const showNew = ref(false)

async function changePassword() {
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) return
  passwordLoading.value = true; passwordError.value = ''; passwordSuccess.value = false
  try {
    await apiClient.post('/api/v1/auth/change-password', {
      current_password: passwordForm.value.current_password,
      new_password: passwordForm.value.new_password
    })
    passwordSuccess.value = true
    passwordForm.value = { current_password: '', new_password: '', confirm_password: '' }
    setTimeout(() => { passwordSuccess.value = false }, 5000)
  } catch (err) {
    passwordError.value = err.response?.data?.error?.message || 'Failed to change password.'
  } finally { passwordLoading.value = false }
}

// User info
const userInitials = computed(() => {
  const u = user.value
  if (!u) return 'U'
  if (u.first_name && u.last_name) return `${u.first_name[0]}${u.last_name[0]}`.toUpperCase()
  if (u.first_name) return u.first_name[0].toUpperCase()
  if (u.email) return u.email[0].toUpperCase()
  return 'U'
})
const fullName = computed(() => {
  const u = user.value
  if (!u) return 'User'
  if (u.first_name || u.last_name) return `${u.first_name || ''} ${u.last_name || ''}`.trim()
  return u.email || 'User'
})
const userRoles = computed(() => {
  const u = user.value
  if (!u || !u.roles) return []
  return u.roles.map(r => {
    const name = typeof r === 'string' ? r : r.name
    return name.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
  })
})
function formatDate(dateStr) {
  if (!dateStr) return 'Never'
  try { return new Date(dateStr).toLocaleDateString('en-UG', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) }
  catch { return dateStr }
}

onMounted(() => {
  breadcrumbStore.set('My Profile', [{ label: 'My Profile' }])
})
</script>
