<template>
  <main class="min-h-screen bg-slate-50 dark:bg-slate-950 flex items-center justify-center px-4 py-12 transition-colors duration-200">
    <div class="w-full max-w-md">
      <!-- Brand Crest & Header -->
      <div class="text-center mb-8">
        <router-link to="/" class="inline-block transition-transform hover:scale-105">
          <img
            src="/main-logo.png"
            alt="National Council of Sports"
            class="h-16 mx-auto mb-4 object-contain drop-shadow"
          />
        </router-link>
        <h1 class="text-2xl md:text-3xl font-extrabold text-[#1a365d] dark:text-white tracking-tight">NCS Intranet Portal</h1>
        <p class="text-slate-600 dark:text-slate-400 font-medium text-sm mt-1.5">Sign in to your staff workspace</p>
      </div>

      <!-- Login Card -->
      <div class="bg-white dark:bg-slate-900 rounded-2xl shadow-xl p-8 border border-slate-200/80 dark:border-slate-800 backdrop-blur-sm">
        <form class="space-y-5" novalidate @submit.prevent="handleLogin">
          <!-- Email Address -->
          <div>
            <label class="block text-sm font-semibold text-slate-800 dark:text-slate-200 mb-1.5" for="email">Email Address</label>
            <div class="relative group">
              <div class="absolute left-3.5 top-1/2 -translate-y-1/2 flex items-center justify-center pointer-events-none text-slate-400 dark:text-slate-500 group-focus-within:text-[#f5a623] transition-colors z-10">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"/><rect x="2" y="4" width="20" height="16" rx="2"/></svg>
              </div>
              <input
                id="email"
                v-model="email"
                type="email"
                autocomplete="email"
                class="w-full h-11 rounded-lg border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900 pr-4 py-2.5 text-slate-900 dark:text-white text-sm font-medium shadow-sm transition-all duration-200 placeholder:text-slate-400 dark:placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-[#f5a623]/30 focus:border-[#f5a623]"
                placeholder="Enter your official email address"
                required
                autofocus
              />
            </div>
          </div>

          <!-- Password -->
          <div>
            <label class="block text-sm font-semibold text-slate-800 dark:text-slate-200 mb-1.5" for="password">Password</label>
            <div class="relative group">
              <div class="absolute left-3.5 top-1/2 -translate-y-1/2 flex items-center justify-center pointer-events-none text-slate-400 dark:text-slate-500 group-focus-within:text-[#f5a623] transition-colors z-10">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
              </div>
              <input
                id="password"
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                class="w-full h-11 rounded-lg border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900 pr-11 py-2.5 text-slate-900 dark:text-white text-sm font-medium shadow-sm transition-all duration-200 placeholder:text-slate-400 dark:placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-[#f5a623]/30 focus:border-[#f5a623]"
                placeholder="Enter your password"
                required
              />
              <button
                type="button"
                class="absolute right-3.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 p-1 flex items-center justify-center rounded-md transition-colors z-10"
                :aria-label="showPassword ? 'Hide password' : 'Show password'"
                @click="showPassword = !showPassword"
              >
                <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0"/><circle cx="12" cy="12" r="3"/></svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m2 2 20 20"/><path d="M6.71 6.71C4.87 7.93 3.31 9.73 2.06 12a1 1 0 0 0 0 .7 10.75 10.75 0 0 0 15.23 4.59"/><path d="M10.58 10.58a2 2 0 0 0 2.83 2.83"/><path d="M14.12 5.22A10.65 10.65 0 0 1 21.94 12a1 1 0 0 1 0 .7 10.8 10.8 0 0 1-2.1 3.13"/></svg>
              </button>
            </div>
          </div>

          <!-- Remember Me & Forgot Password -->
          <div class="flex items-center justify-between gap-4">
            <label class="flex items-center gap-2 cursor-pointer select-none">
              <input id="remember-me" v-model="rememberMe" type="checkbox" class="rounded border-slate-300 text-[#f5a623] focus:ring-[#f5a623] w-4 h-4" />
              <span class="text-sm text-slate-700 dark:text-slate-300 font-medium">Keep me signed in</span>
            </label>
            <a href="#" class="text-sm text-[#d88700] hover:text-[#f5a623] font-semibold hover:underline transition-colors" @click.prevent>Forgot password?</a>
          </div>

          <!-- Error Alert -->
          <p v-if="errorMessage" class="rounded-lg border border-red-200 dark:border-red-900/50 bg-red-50 dark:bg-red-950/40 px-3.5 py-2.5 text-sm font-medium text-red-700 dark:text-red-300 leading-snug">{{ errorMessage }}</p>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading"
            class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-lg transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#f5a623] shadow-md hover:shadow-lg h-11 px-4 w-full bg-[#f5a623] hover:bg-[#e09612] text-[#172b4d] text-base font-extrabold disabled:cursor-not-allowed disabled:opacity-70"
          >
            <span class="flex items-center justify-center gap-2">
              <svg v-if="!loading" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" x2="3" y1="12" y2="12"/></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 animate-spin" viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8v4a4 4 0 0 0-4 4H4z"/></svg>
              {{ loading ? 'Authenticating...' : 'Sign In to Intranet' }}
            </span>
          </button>
        </form>

        <div class="mt-6 text-center">
          <p class="text-slate-600 dark:text-slate-400 font-medium text-sm">Need staff access? <a class="text-[#d88700] hover:text-[#f5a623] font-bold hover:underline transition-colors" href="mailto:info@ncs.go.ug">Contact IT Support</a></p>
        </div>
      </div>

      <div class="text-center mt-6">
        <router-link class="text-slate-600 dark:text-slate-400 hover:text-[#1a365d] dark:hover:text-white font-semibold text-sm inline-flex items-center justify-center gap-1.5 transition-colors" to="/">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 rotate-180" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m9 18 6-6-6-6"/></svg>
          Back to Main Site
        </router-link>
      </div>
    </div>
  </main>
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
const rememberMe = ref(true)
const loading = ref(false)
const errorMessage = ref('')

async function handleLogin() {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await authStore.login(email.value, password.value)
    if (res && res.success) {
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
      errorMessage.value = res?.message || authStore.error || 'Invalid credentials or inactive account.'
    }
  } catch (err) {
    errorMessage.value = err?.response?.data?.error?.message || err.message || 'Login failed. Please check your credentials.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
#email, #password {
  padding-left: 2.75rem !important;
}
input::placeholder {
  color: #94a3b8 !important;
  opacity: 1 !important;
}
.dark input::placeholder {
  color: #64748b !important;
  opacity: 1 !important;
}
</style>
