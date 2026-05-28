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

        <!-- PIN dots -->
        <div class="flex justify-center gap-3 mb-6">
          <div
            v-for="i in 6"
            :key="i"
            class="w-3 h-3 rounded-full transition-colors"
            :class="i <= pin.length ? 'bg-primary-600' : 'bg-gray-200'"
          ></div>
        </div>

        <!-- Numpad -->
        <div class="grid grid-cols-3 gap-3 mb-4">
          <button
            v-for="num in ['1','2','3','4','5','6','7','8','9','','0','⌫']"
            :key="num"
            @click="handleKey(num)"
            :disabled="!num"
            class="h-14 rounded-xl text-lg font-semibold transition-colors"
            :class="num ? 'bg-gray-100 hover:bg-gray-200 text-gray-900 active:bg-primary-100' : 'opacity-0 pointer-events-none'"
          >
            {{ num }}
          </button>
        </div>

        <p v-if="error" class="text-sm text-red-500 mb-3">{{ error }}</p>

        <button @click="unlock" :disabled="verifying || pin.length < 4" class="w-full bg-primary-600 hover:bg-primary-700 text-white font-semibold py-3 rounded-xl transition-colors disabled:opacity-50">
          {{ verifying ? 'Verifying...' : 'Unlock' }}
        </button>

        <button @click="signOut" class="mt-3 w-full text-sm text-gray-400 hover:text-gray-600 py-2">
          Sign out instead
        </button>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth.js'
import { verifyPin } from '@/api/pin.js'

const authStore = useAuthStore()
const locked = ref(false)
const pin = ref('')
const error = ref('')
const verifying = ref(false)

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

function handleKey(key) {
  if (!key) return
  if (key === '⌫') {
    pin.value = pin.value.slice(0, -1)
  } else if (pin.value.length < 6) {
    pin.value += key
  }
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
