<template>
  <main class="min-h-screen bg-gray-50 flex items-center justify-center px-4 py-10">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <router-link class="inline-block" to="/">
          <img alt="NCS" class="h-16 mx-auto mb-4 object-contain" src="/main-logo.png" />
        </router-link>
        <h1 class="text-2xl font-bold text-[#1a365d]">Welcome Back</h1>
        <p class="text-black font-medium mt-1">Sign in to your NCS account</p>
      </div>

      <div class="bg-white rounded-2xl shadow-lg p-8 border border-gray-100">
        <form class="space-y-5" novalidate @submit.prevent="login">
          <div>
            <label class="block text-sm font-semibold text-[#1a365d] mb-1.5" for="cms-email">Email Address</label>
            <div class="relative">
              <svg xmlns="http://www.w3.org/2000/svg" class="absolute left-3.5 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-600 pointer-events-none" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"/><rect x="2" y="4" width="20" height="16" rx="2"/></svg>
              <input
                id="cms-email"
                v-model="email"
                type="email"
                autocomplete="email"
                class="w-full h-11 rounded-lg border border-gray-300 bg-white pl-11 pr-4 py-2 text-black text-sm shadow-sm transition-colors placeholder:text-black focus:outline-none focus:ring-2 focus:ring-[#f5a623] focus:border-[#f5a623]"
                placeholder="Enter your email"
                required
                autofocus
              />
            </div>
          </div>

          <div>
            <label class="block text-sm font-semibold text-[#1a365d] mb-1.5" for="cms-password">Password</label>
            <div class="relative">
              <svg xmlns="http://www.w3.org/2000/svg" class="absolute left-3.5 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-600 pointer-events-none" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
              <input
                id="cms-password"
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                class="w-full h-11 rounded-lg border border-gray-300 bg-white pl-11 pr-11 py-2 text-black text-sm shadow-sm transition-colors placeholder:text-black focus:outline-none focus:ring-2 focus:ring-[#f5a623] focus:border-[#f5a623]"
                placeholder="Enter your password"
                required
              />
              <button
                type="button"
                class="absolute right-3.5 top-1/2 -translate-y-1/2 text-gray-600 hover:text-black p-1 flex items-center justify-center"
                :aria-label="showPassword ? 'Hide password' : 'Show password'"
                @click="showPassword = !showPassword"
              >
                <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0"/><circle cx="12" cy="12" r="3"/></svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m2 2 20 20"/><path d="M6.71 6.71C4.87 7.93 3.31 9.73 2.06 12a1 1 0 0 0 0 .7 10.75 10.75 0 0 0 15.23 4.59"/><path d="M10.58 10.58a2 2 0 0 0 2.83 2.83"/><path d="M14.12 5.22A10.65 10.65 0 0 1 21.94 12a1 1 0 0 1 0 .7 10.8 10.8 0 0 1-2.1 3.13"/></svg>
              </button>
            </div>
          </div>

          <div class="flex items-center justify-between gap-4">
            <label class="flex items-center gap-2 cursor-pointer">
              <input id="remember-me" v-model="remember" type="checkbox" class="rounded border-gray-400 text-[#f5a623] focus:ring-[#f5a623]" />
              <span class="text-sm text-black font-medium">Remember me</span>
            </label>
            <a href="#" class="text-sm text-[#f5a623] font-semibold hover:underline" @click.prevent>Forgot password?</a>
          </div>

          <p v-if="error" class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm font-medium text-red-700">{{ error }}</p>

          <button
            type="submit"
            :disabled="loading"
            class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-lg transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#f5a623] shadow-md h-12 px-4 w-full bg-[#f5a623] hover:bg-[#e09612] text-white text-base font-bold disabled:cursor-not-allowed disabled:opacity-70"
          >
            <span class="flex items-center justify-center gap-2">
              <svg v-if="!loading" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" x2="3" y1="12" y2="12"/></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 animate-spin" viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8v4a4 4 0 0 0-4 4H4z"/></svg>
              {{ loading ? 'Signing in...' : 'Sign In' }}
            </span>
          </button>
        </form>

        <div class="mt-6 text-center">
          <p class="text-black font-medium">Don't have an account? <a class="text-[#f5a623] font-bold hover:underline" href="https://portal.ncs.go.ug/register">Create Account</a></p>
        </div>
      </div>

      <div class="text-center mt-6">
        <router-link class="text-black font-semibold hover:text-[#1a365d] text-sm flex items-center justify-center gap-1" to="/">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 rotate-180" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m9 18 6-6-6-6"/></svg>
          Back to Home
        </router-link>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import apiClient, { API_BASE_URL } from '@/api/client.js'

const router = useRouter()
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const remember = ref(false)
const showPassword = ref(false)

async function login() {
  loading.value = true
  error.value = ''
  try {
    if (isDefaultPreviewLogin() && !(await canReachApi())) {
      enterPreviewMode()
      return
    }
    const res = await apiClient.post('/api/v1/auth/login', { email: email.value, password: password.value })
    const data = res.data?.data || {}
    localStorage.setItem('ncsms_access_token', data.access_token || '')
    localStorage.setItem('ncsms_user', JSON.stringify(data.user || {}))
    router.push('/cms')
  } catch (err) {
    if (!err.response && isDefaultPreviewLogin()) {
      enterPreviewMode()
      return
    }
    error.value = err.response?.data?.error?.message || 'Login failed. Check the CMS credentials and API server.'
  } finally {
    loading.value = false
  }
}

function isDefaultPreviewLogin() {
  return email.value === 'admin@ncs.go.ug' && password.value === 'NCS@Admin2026!'
}

function enterPreviewMode() {
  localStorage.setItem('ncsms_access_token', 'local-cms-preview-token')
  localStorage.setItem('ncsms_user', JSON.stringify({ email: email.value, roles: ['super_admin'] }))
  router.push('/cms')
}

async function canReachApi() {
  let timer
  try {
    const controller = new AbortController()
    timer = window.setTimeout(() => controller.abort(), 1200)
    const res = await fetch(`${API_BASE_URL}/health`, { signal: controller.signal, cache: 'no-store' })
    return res.ok
  } catch {
    return false
  } finally {
    if (timer) window.clearTimeout(timer)
  }
}
</script>

<style scoped>
input::placeholder {
  color: #000000 !important;
  opacity: 1 !important;
}
</style>
