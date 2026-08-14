<template>
  <div class="min-h-screen bg-slate-100 dark:bg-slate-950 flex flex-col justify-between p-4 sm:p-6 font-poppins transition-colors duration-300">
    
    <!-- Top Bar with Theme Toggle -->
    <div class="w-full max-w-5xl mx-auto flex justify-end items-center py-2 px-2">
      <ThemeToggle />
    </div>

    <!-- Center Login Card -->
    <div class="w-full max-w-md mx-auto my-auto py-6">
      
      <!-- Brand Crest & Header -->
      <div class="text-center mb-6">
        <router-link to="/" class="inline-block transition-transform hover:scale-105">
          <img
            src="/main-logo.png"
            alt="National Council of Sports"
            class="h-20 mx-auto mb-3 object-contain drop-shadow"
          />
        </router-link>
        <h1 class="text-2xl font-extrabold text-slate-900 dark:text-white tracking-tight">NCS Intranet Portal</h1>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1 uppercase tracking-widest font-semibold">
          National Council of Sports · Uganda
        </p>
      </div>

      <!-- Otika Card Primary -->
      <div class="card card-primary bg-white dark:bg-slate-900 shadow-xl border border-slate-200 dark:border-slate-800 rounded-2xl overflow-hidden transition-colors duration-300">
        <div class="card-header border-b border-slate-100 dark:border-slate-800 py-4 px-6 bg-slate-50/50 dark:bg-slate-800/30">
          <h4 class="text-sm font-bold text-slate-800 dark:text-slate-200 uppercase tracking-wide flex items-center gap-2">
            <i class="icofont-key text-blue-600 dark:text-blue-400"></i> Secure Sign In
          </h4>
        </div>

        <div class="card-body p-6 sm:p-8">
          <!-- Error Alert -->
          <div v-if="errorMessage" class="mb-5 flex items-start gap-3 p-3.5 bg-red-50 dark:bg-red-950/40 border border-red-200 dark:border-red-900/50 rounded-xl text-red-700 dark:text-red-400 text-xs">
            <i class="icofont-warning-alt text-base flex-shrink-0 mt-0.5"></i>
            <p class="font-medium">{{ errorMessage }}</p>
          </div>

          <form @submit.prevent="handleLogin" class="space-y-4">
            <!-- Email -->
            <div class="form-group">
              <label for="email" class="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1.5 uppercase tracking-wider">
                Official Email Address
              </label>
              <input
                id="email"
                v-model="email"
                type="email"
                required
                autocomplete="email"
                placeholder="e.g. admin@ncs.go.ug"
                class="form-control w-full px-3.5 h-10 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-xs sm:text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition"
              />
            </div>

            <!-- Password -->
            <div class="form-group">
              <div class="flex items-center justify-between mb-1.5">
                <label for="password" class="block text-xs font-bold text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                  Password
                </label>
              </div>
              <div class="relative flex items-center">
                <input
                  id="password"
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  required
                  autocomplete="current-password"
                  placeholder="••••••••••••"
                  class="form-control w-full pl-3.5 pr-10 h-10 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-xs sm:text-sm text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition"
                />
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="absolute right-0 inset-y-0 w-10 flex items-center justify-center text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 focus:outline-none"
                  :title="showPassword ? 'Hide Password' : 'Show Password'"
                >
                  <i :class="showPassword ? 'icofont-eye-blocked' : 'icofont-eye'" class="text-sm leading-none"></i>
                </button>
              </div>
            </div>

            <!-- Remember & Forgot -->
            <div class="flex items-center justify-between text-xs pt-1">
              <label class="flex items-center gap-2 cursor-pointer select-none">
                <input
                  v-model="rememberMe"
                  type="checkbox"
                  class="w-4 h-4 rounded text-blue-600 border-slate-300 dark:border-slate-700 focus:ring-blue-500"
                />
                <span class="text-slate-600 dark:text-slate-400 font-medium">Keep me signed in</span>
              </label>
            </div>

            <!-- Submit Button -->
            <div class="pt-2">
              <button
                type="submit"
                :disabled="loading"
                class="w-full py-3 px-4 bg-blue-600 hover:bg-blue-700 disabled:opacity-60 text-white font-bold text-sm rounded-xl shadow-lg shadow-blue-500/25 transition-all flex items-center justify-center gap-2"
              >
                <i v-if="loading" class="icofont-spinner icofont-spin"></i>
                <i v-else class="icofont-login"></i>
                <span>{{ loading ? 'Authenticating...' : 'Sign In to Portal' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="text-center text-xs text-slate-400 dark:text-slate-500 py-3">
      &copy; {{ new Date().getFullYear() }} National Council of Sports, Uganda. All rights reserved.
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import ThemeToggle from '@/components/ui/ThemeToggle.vue'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const rememberMe = ref(true)
const loading = ref(false)
const errorMessage = ref('')

async function handleLogin() {
  loading.value = true
  errorMessage.value = ''
  try {
    const success = await authStore.login(email.value, password.value)
    if (success) {
      // Direct user to appropriate primary dashboard
      const rawRoles = authStore.user?.roles || []
      const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)
      if (roles.includes('general_secretary')) {
        router.push('/executive/appraisals')
      } else if (roles.includes('ags_technical')) {
        router.push('/executive/ags-technical')
      } else if (roles.includes('ags_admin')) {
        router.push('/executive/ags-admin')
      } else if (roles.includes('stores_officer')) {
        router.push('/stores/inventory')
      } else if (roles.includes('facilities_manager')) {
        router.push('/facilities/venues')
      } else if (roles.includes('transport_officer')) {
        router.push('/fleet/transport')
      } else if (roles.includes('medical_officer') || roles.includes('physiotherapist')) {
        router.push('/medical/sports-science')
      } else if (roles.includes('legal_counsel')) {
        router.push('/legal/compliance')
      } else {
        router.push('/dashboard')
      }
    } else {
      errorMessage.value = authStore.error || 'Invalid credentials or inactive account.'
    }
  } catch (err) {
    errorMessage.value = err?.response?.data?.error?.message || err.message || 'Login failed. Please check your credentials.'
  } finally {
    loading.value = false
  }
}
</script>
