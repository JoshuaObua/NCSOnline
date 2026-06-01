<template>
  <div class="min-h-screen flex">

    <!-- Left panel — brand -->
    <div class="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-primary-800 via-primary-700 to-primary-600 flex-col items-center justify-center p-12 relative overflow-hidden">
      <div class="absolute -top-24 -left-24 w-72 h-72 bg-white/5 rounded-full pointer-events-none"></div>
      <div class="absolute -bottom-20 -right-20 w-96 h-96 bg-white/5 rounded-full pointer-events-none"></div>
      <div class="absolute top-1/2 left-1/4 w-32 h-32 bg-white/5 rounded-full pointer-events-none"></div>

      <div class="relative text-center max-w-sm">
        <div class="inline-flex items-center justify-center w-24 h-24 bg-white rounded-2xl shadow-xl mb-6 p-2">
          <img src="/main-logo.png" alt="NCS Logo" class="w-full h-full object-contain" />
        </div>
        <h1 class="text-3xl font-bold text-white mb-3">Create Your<br>NCS Account</h1>
        <p class="text-primary-200 text-base leading-relaxed">
          Register to apply for sports licences, track your applications, and manage renewals online.
        </p>

        <div class="mt-8 space-y-3 text-left">
          <div class="flex items-center gap-3 bg-white/10 rounded-xl px-4 py-3">
            <i class="icofont-check-circled text-yellow-300 text-xl flex-shrink-0"></i>
            <span class="text-white/80 text-sm">Apply for sports licences online</span>
          </div>
          <div class="flex items-center gap-3 bg-white/10 rounded-xl px-4 py-3">
            <i class="icofont-check-circled text-yellow-300 text-xl flex-shrink-0"></i>
            <span class="text-white/80 text-sm">Track application status in real time</span>
          </div>
          <div class="flex items-center gap-3 bg-white/10 rounded-xl px-4 py-3">
            <i class="icofont-check-circled text-yellow-300 text-xl flex-shrink-0"></i>
            <span class="text-white/80 text-sm">Renew licences before they expire</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Right panel — registration form -->
    <div class="flex-1 flex flex-col items-center justify-center bg-gray-50 px-6 py-12 overflow-y-auto">
      <div class="w-full max-w-sm">

        <!-- Mobile logo -->
        <div class="flex items-center justify-center gap-3 mb-8 lg:hidden">
          <div class="w-12 h-12 bg-primary-700 rounded-xl flex items-center justify-center overflow-hidden p-1">
            <img src="/main-logo.png" alt="NCS Logo" class="w-full h-full object-contain" />
          </div>
          <div>
            <div class="text-lg font-bold text-gray-900">NCSMS Portal</div>
            <div class="text-xs text-gray-500">National Council of Sports</div>
          </div>
        </div>

        <!-- Heading -->
        <div class="mb-7">
          <h2 class="text-2xl font-bold text-gray-900">Create account</h2>
          <p class="text-gray-500 text-sm mt-1">Register to manage your sports licence applications</p>
        </div>

        <!-- Success -->
        <div v-if="success" class="mb-5 flex items-start gap-3 p-4 bg-green-50 border border-green-200 rounded-xl">
          <i class="icofont-check-circled text-green-500 text-lg flex-shrink-0 mt-0.5"></i>
          <div>
            <p class="text-green-800 text-sm font-semibold">Account created!</p>
            <p class="text-green-700 text-sm mt-0.5">Redirecting to your portal…</p>
          </div>
        </div>

        <!-- Error -->
        <div v-if="errorMessage" class="mb-5 flex items-start gap-3 p-4 bg-red-50 border border-red-200 rounded-xl">
          <i class="icofont-warning-alt text-red-500 text-lg flex-shrink-0 mt-0.5"></i>
          <p class="text-red-700 text-sm">{{ errorMessage }}</p>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleRegister" class="space-y-4">
          <!-- Name row -->
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="form-label">First Name</label>
              <input
                v-model="firstName"
                type="text"
                autocomplete="given-name"
                required
                placeholder="John"
                class="form-input"
              />
            </div>
            <div>
              <label class="form-label">Last Name</label>
              <input
                v-model="lastName"
                type="text"
                autocomplete="family-name"
                required
                placeholder="Doe"
                class="form-input"
              />
            </div>
          </div>

          <div>
            <label class="form-label">Email Address</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
                <i class="icofont-email text-gray-400 text-base leading-none"></i>
              </div>
              <input
                v-model="email"
                type="email"
                autocomplete="email"
                required
                placeholder="your@email.com"
                class="form-input pl-10"
              />
            </div>
          </div>

          <div>
            <label class="form-label">Password</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
                <i class="icofont-lock text-gray-400 text-base leading-none"></i>
              </div>
              <input
                v-model="password"
                :type="showPass ? 'text' : 'password'"
                autocomplete="new-password"
                required
                minlength="8"
                placeholder="Min. 8 characters"
                class="form-input pl-10 pr-10"
              />
              <button type="button" @click="showPass = !showPass"
                class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-gray-400 hover:text-gray-600">
                <i :class="showPass ? 'icofont-eye-blocked' : 'icofont-eye'" class="text-base leading-none"></i>
              </button>
            </div>
          </div>

          <div>
            <label class="form-label">Confirm Password</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
                <i class="icofont-lock text-gray-400 text-base leading-none"></i>
              </div>
              <input
                v-model="confirmPassword"
                :type="showPass ? 'text' : 'password'"
                autocomplete="new-password"
                required
                placeholder="Repeat password"
                class="form-input pl-10"
                :class="confirmPassword && password !== confirmPassword ? 'border-red-300 focus:ring-red-400' : ''"
              />
            </div>
            <p v-if="confirmPassword && password !== confirmPassword" class="text-xs text-red-500 mt-1">
              Passwords do not match
            </p>
          </div>

          <button
            type="submit"
            :disabled="loading || (confirmPassword && password !== confirmPassword)"
            class="w-full flex items-center justify-center gap-2 py-3 px-4 bg-[#F48C06] hover:bg-[#d47b05] text-white font-semibold rounded-xl shadow-sm hover:shadow-md transition-all disabled:opacity-60 disabled:cursor-not-allowed"
          >
            <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
            </svg>
            <i v-else class="icofont-ui-user-group text-base leading-none"></i>
            {{ loading ? 'Creating account…' : 'Create Account' }}
          </button>
        </form>

        <p class="text-center text-sm text-gray-500 mt-6">
          Already have an account?
          <router-link to="/login" class="text-[#F48C06] font-semibold hover:text-[#d47b05]">Sign in</router-link>
        </p>

        <p class="text-center text-xs text-gray-400 mt-6">
          &copy; {{ year }} National Council of Sports — Uganda
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'

const router = useRouter()
const authStore = useAuthStore()

const firstName = ref('')
const lastName = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const showPass = ref(false)
const errorMessage = ref('')
const loading = ref(false)
const success = ref(false)
const year = new Date().getFullYear()

async function handleRegister() {
  if (password.value !== confirmPassword.value) return
  errorMessage.value = ''
  loading.value = true
  const result = await authStore.register(firstName.value, lastName.value, email.value, password.value)
  loading.value = false
  if (result.success) {
    success.value = true
    setTimeout(() => router.push('/my-portal'), 1200)
  } else {
    errorMessage.value = result.message || 'Registration failed. Please try again.'
  }
}
</script>
