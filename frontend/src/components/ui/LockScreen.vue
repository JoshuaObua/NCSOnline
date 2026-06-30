<template>
  <Teleport to="body">
    <div
      v-if="locked"
      class="fixed inset-0 bg-gray-900/95 backdrop-blur-sm z-[9999] flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-sm p-8 text-center">
        <div class="w-16 h-16 bg-primary-600 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
        </div>
        <h2 class="text-xl font-bold text-gray-900 mb-1">Screen Locked</h2>
        <p class="text-sm text-gray-500 mb-6">Enter your PIN to continue</p>

        <form @submit.prevent="unlock" class="space-y-4">
          <div class="relative">
            <input
              ref="pinInput"
              v-model="pin"
              type="password"
              inputmode="numeric"
              maxlength="6"
              pattern="[0-9]*"
              autocomplete="current-password"
              placeholder="Enter PIN"
              class="w-full text-center text-2xl tracking-[0.4em] font-semibold px-4 py-3 border-2 border-gray-200 rounded-xl focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-100 transition-colors"
              @input="onPinInput"
            />
            <button
              type="button"
              @click="showPin = !showPin"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
              :title="showPin ? 'Hide PIN' : 'Show PIN'"
            >
              <i :class="showPin ? 'icofont-eye-blocked' : 'icofont-eye'"></i>
            </button>
          </div>

          <p v-if="error" class="text-sm text-red-500">{{ error }}</p>

          <button
            type="submit"
            :disabled="verifying || pin.length < 4"
            class="w-full bg-primary-600 hover:bg-primary-700 text-white font-semibold py-3 rounded-xl transition-colors disabled:opacity-50"
          >
            {{ verifying ? 'Verifying...' : 'Unlock' }}
          </button>
        </form>

        <button @click="signOut" class="mt-3 w-full text-sm text-gray-400 hover:text-gray-600 py-2">
          Sign out instead
        </button>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, nextTick, onMounted, onUnmounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth.js'
import { verifyPin } from '@/api/pin.js'

const authStore = useAuthStore()
const locked = ref(false)
const pin = ref('')
const error = ref('')
const verifying = ref(false)
const showPin = ref(false)
const pinInput = ref(null)

// Toggle input type live when showPin flips (keeps numeric inputmode)
watch(showPin, async (v) => {
  await nextTick()
  if (pinInput.value) pinInput.value.type = v ? 'text' : 'password'
})

watch(locked, async (v) => {
  if (v) {
    await nextTick()
    pinInput.value?.focus()
  }
})

function onPinInput(e) {
  // Only allow digits
  const cleaned = (e.target.value || '').replace(/\D/g, '').slice(0, 6)
  pin.value = cleaned
  if (e.target.value !== cleaned) e.target.value = cleaned
  error.value = ''
}

const INACTIVITY_MS = 10 * 60 * 1000 // 10 minutes
let timer = null

function resetTimer() {
  clearTimeout(timer)
  if (authStore.isAuthenticated) {
    timer = setTimeout(lock, INACTIVITY_MS)
  }
}

function lock() {
  if (!authStore.isAuthenticated) return
  locked.value = true
  pin.value = ''
  error.value = ''
}

async function unlock() {
  if (pin.value.length < 4) return
  verifying.value = true
  error.value = ''
  try {
    await verifyPin(pin.value)
    locked.value = false
    pin.value = ''
    resetTimer()
  } catch (err) {
    const code = err.response?.data?.error?.code
    if (code === 'PIN_NOT_SET') {
      // No PIN set — unlock without verification
      locked.value = false
      pin.value = ''
      resetTimer()
    } else {
      error.value = 'Incorrect PIN. Try again.'
      pin.value = ''
    }
  } finally {
    verifying.value = false
  }
}

async function signOut() {
  locked.value = false
  await authStore.logout()
}

const events = ['mousemove','keydown','click','touchstart','scroll']

function addListeners() {
  events.forEach(e => window.addEventListener(e, resetTimer, { passive: true }))
}
function removeListeners() {
  events.forEach(e => window.removeEventListener(e, resetTimer))
}

onMounted(() => {
  addListeners()
  resetTimer()
})

onUnmounted(() => {
  removeListeners()
  clearTimeout(timer)
})
</script>
