<template>
  <div class="min-h-screen flex">

    <!-- Left panel — brand / illustration -->
    <div class="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-primary-800 via-primary-700 to-primary-600 flex-col items-center justify-center p-12 relative overflow-hidden">
      <!-- Decorative circles -->
      <div class="absolute -top-24 -left-24 w-72 h-72 bg-white/5 rounded-full pointer-events-none"></div>
      <div class="absolute -bottom-20 -right-20 w-96 h-96 bg-white/5 rounded-full pointer-events-none"></div>
      <div class="absolute top-1/2 left-1/4 w-32 h-32 bg-white/5 rounded-full pointer-events-none"></div>

      <div class="relative text-center max-w-sm">
        <!-- Logo -->
        <div class="inline-flex items-center justify-center w-24 h-24 bg-white rounded-2xl shadow-xl mb-6 p-2">
          <img src="/main-logo.png" alt="NCS Logo" class="w-full h-full object-contain" />
        </div>
        <h1 class="text-3xl font-bold text-white mb-3">National Council<br>of Sports</h1>
        <p class="text-primary-200 text-base leading-relaxed">
          Uganda's premier sports management and registration platform.
        </p>

        <div class="mt-8 grid grid-cols-3 gap-4 text-center">
          <div class="bg-white/10 rounded-xl p-3">
            <div class="text-2xl font-bold text-white">40+</div>
            <div class="text-primary-200 text-xs mt-0.5">Federations</div>
          </div>
          <div class="bg-white/10 rounded-xl p-3">
            <div class="text-2xl font-bold text-white">500+</div>
            <div class="text-primary-200 text-xs mt-0.5">Clubs</div>
          </div>
          <div class="bg-white/10 rounded-xl p-3">
            <div class="text-2xl font-bold text-white">1M+</div>
            <div class="text-primary-200 text-xs mt-0.5">Athletes</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Right panel — login form -->
    <div class="flex-1 flex flex-col items-center justify-center bg-gray-50 px-6 py-12">
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
          <h2 class="text-2xl font-bold text-gray-900">Welcome back</h2>
          <p class="text-gray-500 text-sm mt-1">Sign in to your account to continue</p>
        </div>

        <!-- Error -->
        <div v-if="errorMessage" class="mb-5 flex items-start gap-3 p-4 bg-red-50 border border-red-200 rounded-xl">
          <i class="icofont-warning-alt text-red-500 text-lg flex-shrink-0 mt-0.5"></i>
          <p class="text-red-700 text-sm">{{ errorMessage }}</p>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleLogin" class="space-y-4">
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
                placeholder="Enter your email"
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
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                required
                placeholder="Enter your password"
                class="form-input pl-10 pr-10"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-gray-400 hover:text-gray-600"
              >
                <i :class="showPassword ? 'icofont-eye-blocked' : 'icofont-eye'" class="text-base leading-none"></i>
              </button>
            </div>
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full flex items-center justify-center gap-2 py-3 px-4 bg-primary-700 hover:bg-primary-600 text-white font-semibold rounded-xl shadow-sm hover:shadow-md transition-all disabled:opacity-60 disabled:cursor-not-allowed"
          >
            <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
            </svg>
            <i v-else class="icofont-arrow-right text-base leading-none"></i>
            {{ loading ? 'Signing In…' : 'Sign In' }}
          </button>
        </form>

        <p class="text-center text-xs text-gray-400 mt-8">
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

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const errorMessage = ref('')
const loading = ref(false)
const year = new Date().getFullYear()

async function handleLogin() {
  errorMessage.value = ''
  loading.value = true
  const result = await authStore.login(email.value, password.value)
  loading.value = false
  if (result.success) {
    router.push('/dashboard')
  } else {
    errorMessage.value = result.message || 'Login failed. Please try again.'
  }
}
</script>
