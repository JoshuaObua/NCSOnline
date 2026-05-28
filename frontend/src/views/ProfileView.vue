<template>
  <LayoutDefault title="My Profile">
    <div class="max-w-2xl mx-auto space-y-6">
      <!-- Profile Info Card -->
      <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        <div class="bg-gradient-to-r from-primary-700 to-primary-600 px-6 py-8">
          <div class="flex items-center gap-5">
            <div class="w-16 h-16 bg-white rounded-2xl flex items-center justify-center shadow-lg">
              <span class="text-primary-700 text-2xl font-bold">{{ userInitials }}</span>
            </div>
            <div>
              <h2 class="text-xl font-semibold text-white">{{ fullName }}</h2>
              <p class="text-primary-200 text-sm mt-0.5">{{ user?.email }}</p>
              <div class="flex flex-wrap gap-1.5 mt-2">
                <span
                  v-for="role in userRoles"
                  :key="role"
                  class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-primary-500/40 text-white"
                >
                  {{ role }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="px-6 py-5">
          <h3 class="text-sm font-semibold text-gray-700 mb-4">Account Information</h3>
          <dl class="grid grid-cols-2 gap-4">
            <div>
              <dt class="text-xs text-gray-500 uppercase tracking-wider">Email</dt>
              <dd class="text-sm font-medium text-gray-900 mt-0.5">{{ user?.email || 'N/A' }}</dd>
            </div>
            <div>
              <dt class="text-xs text-gray-500 uppercase tracking-wider">User ID</dt>
              <dd class="text-sm font-medium text-gray-900 mt-0.5 font-mono text-xs">{{ user?.id?.substring(0, 16) || user?.user_id?.substring(0, 16) || 'N/A' }}</dd>
            </div>
            <div v-if="user?.phone">
              <dt class="text-xs text-gray-500 uppercase tracking-wider">Phone</dt>
              <dd class="text-sm font-medium text-gray-900 mt-0.5">{{ user.phone }}</dd>
            </div>
            <div v-if="user?.last_login_at">
              <dt class="text-xs text-gray-500 uppercase tracking-wider">Last Login</dt>
              <dd class="text-sm font-medium text-gray-900 mt-0.5">{{ formatDate(user.last_login_at) }}</dd>
            </div>
            <div>
              <dt class="text-xs text-gray-500 uppercase tracking-wider">Account Status</dt>
              <dd class="mt-0.5">
                <span :class="user?.is_active !== false ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium">
                  {{ user?.is_active !== false ? 'Active' : 'Inactive' }}
                </span>
              </dd>
            </div>
          </dl>
        </div>
      </div>

      <!-- PIN Management Card -->
      <div class="bg-white rounded-xl border border-gray-200 shadow-sm">
        <div class="px-6 py-4 border-b border-gray-100">
          <h3 class="font-semibold text-gray-800">Screen Lock PIN</h3>
          <p class="text-sm text-gray-500 mt-0.5">Set or change your 4–6 digit PIN used to unlock the app after inactivity.</p>
        </div>
        <div class="px-6 py-5">
          <div v-if="pinSuccess" class="p-3 bg-green-50 border border-green-200 rounded-lg text-sm text-green-700 mb-4">{{ pinSuccess }}</div>
          <div v-if="pinError" class="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700 mb-4">{{ pinError }}</div>

          <div class="flex gap-3 mb-4">
            <button
              @click="pinMode = 'set'"
              :class="pinMode === 'set' ? 'bg-primary-600 text-white' : 'bg-gray-100 text-gray-600'"
              class="px-4 py-1.5 rounded-full text-sm font-medium transition-colors"
            >Set PIN</button>
            <button
              @click="pinMode = 'change'"
              :class="pinMode === 'change' ? 'bg-primary-600 text-white' : 'bg-gray-100 text-gray-600'"
              class="px-4 py-1.5 rounded-full text-sm font-medium transition-colors"
            >Change PIN</button>
          </div>

          <form @submit.prevent="handlePin" class="space-y-4">
            <div v-if="pinMode === 'change'">
              <label class="block text-sm font-medium text-gray-700 mb-1.5">Current PIN</label>
              <input v-model="pinForm.current" type="password" inputmode="numeric" maxlength="6" class="w-full px-3 py-2.5 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" placeholder="Current PIN">
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1.5">{{ pinMode === 'set' ? 'New PIN' : 'New PIN' }}</label>
              <input v-model="pinForm.pin" type="password" inputmode="numeric" maxlength="6" class="w-full px-3 py-2.5 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" placeholder="4–6 digit PIN">
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1.5">Confirm PIN</label>
              <input v-model="pinForm.confirm" type="password" inputmode="numeric" maxlength="6" class="w-full px-3 py-2.5 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" placeholder="Re-enter PIN">
              <p v-if="pinForm.confirm && pinForm.confirm !== pinForm.pin" class="text-xs text-red-500 mt-1">PINs do not match</p>
            </div>
            <button
              type="submit"
              :disabled="pinLoading || pinForm.confirm !== pinForm.pin"
              class="px-5 py-2.5 bg-primary-700 hover:bg-primary-600 text-white font-medium rounded-lg text-sm transition-colors disabled:opacity-60"
            >
              {{ pinLoading ? 'Saving...' : (pinMode === 'set' ? 'Set PIN' : 'Change PIN') }}
            </button>
          </form>
        </div>
      </div>

      <!-- Change Password Card -->
      <div class="bg-white rounded-xl border border-gray-200 shadow-sm">
        <div class="px-6 py-4 border-b border-gray-100">
          <h3 class="font-semibold text-gray-800">Change Password</h3>
          <p class="text-sm text-gray-500 mt-0.5">Update your account password. Minimum 8 characters.</p>
        </div>

        <div class="px-6 py-5">
          <!-- Success alert -->
          <div v-if="passwordSuccess" class="flex items-center gap-2 p-3 bg-green-50 border border-green-200 rounded-lg text-sm text-green-700 mb-4">
            <svg class="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Password changed successfully!
          </div>

          <!-- Error alert -->
          <div v-if="passwordError" class="flex items-start gap-2 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700 mb-4">
            <svg class="w-4 h-4 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
            </svg>
            {{ passwordError }}
          </div>

          <form @submit.prevent="changePassword" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1.5">Current Password</label>
              <div class="relative">
                <input
                  v-model="passwordForm.current_password"
                  :type="showCurrent ? 'text' : 'password'"
                  required
                  class="w-full px-3 py-2.5 pr-10 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                />
                <button type="button" @click="showCurrent = !showCurrent" class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                    <path v-if="!showCurrent" stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" /><path v-if="!showCurrent" stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path v-if="showCurrent" stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88" />
                  </svg>
                </button>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1.5">New Password</label>
              <div class="relative">
                <input
                  v-model="passwordForm.new_password"
                  :type="showNew ? 'text' : 'password'"
                  required
                  minlength="8"
                  class="w-full px-3 py-2.5 pr-10 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                />
                <button type="button" @click="showNew = !showNew" class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                    <path v-if="!showNew" stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" /><path v-if="!showNew" stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path v-if="showNew" stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88" />
                  </svg>
                </button>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1.5">Confirm New Password</label>
              <input
                v-model="passwordForm.confirm_password"
                type="password"
                required
                class="w-full px-3 py-2.5 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                :class="{ 'border-red-400': passwordForm.confirm_password && passwordForm.confirm_password !== passwordForm.new_password }"
              />
              <p v-if="passwordForm.confirm_password && passwordForm.confirm_password !== passwordForm.new_password" class="text-xs text-red-500 mt-1">
                Passwords do not match
              </p>
            </div>

            <div class="pt-2">
              <button
                type="submit"
                :disabled="passwordLoading || (passwordForm.confirm_password !== passwordForm.new_password)"
                class="px-6 py-2.5 bg-primary-700 hover:bg-primary-600 text-white font-medium rounded-lg text-sm transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center gap-2"
              >
                <svg v-if="passwordLoading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ passwordLoading ? 'Updating...' : 'Update Password' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, computed } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { useAuthStore } from '@/stores/auth.js'
import apiClient from '@/api/client.js'
import { setPin, changePin } from '@/api/pin.js'

const authStore = useAuthStore()
const user = computed(() => authStore.user)

// PIN management
const pinMode = ref('set')
const pinForm = ref({ current: '', pin: '', confirm: '' })
const pinLoading = ref(false)
const pinError = ref('')
const pinSuccess = ref('')

async function handlePin() {
  if (pinForm.value.pin !== pinForm.value.confirm) return
  pinLoading.value = true
  pinError.value = ''
  pinSuccess.value = ''
  try {
    if (pinMode.value === 'set') {
      await setPin(pinForm.value.pin)
      pinSuccess.value = 'PIN set successfully.'
    } else {
      await changePin(pinForm.value.current, pinForm.value.pin)
      pinSuccess.value = 'PIN changed successfully.'
    }
    pinForm.value = { current: '', pin: '', confirm: '' }
    setTimeout(() => { pinSuccess.value = '' }, 5000)
  } catch (err) {
    pinError.value = err.response?.data?.error?.message || 'Failed to update PIN.'
  } finally {
    pinLoading.value = false
  }
}

const passwordForm = ref({ current_password: '', new_password: '', confirm_password: '' })
const passwordLoading = ref(false)
const passwordError = ref('')
const passwordSuccess = ref(false)
const showCurrent = ref(false)
const showNew = ref(false)

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
  if (!dateStr) return 'N/A'
  try {
    return new Date(dateStr).toLocaleDateString('en-UG', { day: '2-digit', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit' })
  } catch { return dateStr }
}

async function changePassword() {
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) return
  passwordLoading.value = true
  passwordError.value = ''
  passwordSuccess.value = false
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
  } finally {
    passwordLoading.value = false
  }
}
</script>
