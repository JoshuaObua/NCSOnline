<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <div class="min-h-[calc(100vh-120px)] flex items-center justify-center py-12 px-4">
      <div class="w-full max-w-md">

        <!-- Logo + heading -->
        <div class="text-center mb-8">
          <router-link to="/" class="inline-block">
            <img
              :src="logoSrc"
              alt="NCS"
              class="h-16 mx-auto mb-4 object-contain"
            />
          </router-link>
          <h1 class="text-2xl font-bold text-[#1a365d]">Welcome Back</h1>
          <p class="text-gray-600 mt-1">Sign in to your NCS account</p>
        </div>

        <!-- Card -->
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <!-- Error -->
          <div v-if="errorMessage" class="mb-5 flex items-start gap-3 p-4 bg-red-50 border border-red-200 rounded-xl">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" aria-hidden="true"><circle cx="12" cy="12" r="10"/><line x1="12" x2="12" y1="8" y2="12"/><line x1="12" x2="12.01" y1="16" y2="16"/></svg>
            <p class="text-red-700 text-sm">{{ errorMessage }}</p>
          </div>

          <form @submit.prevent="handleLogin" class="space-y-5">
            <!-- Email -->
            <div>
              <label for="email" class="text-sm font-medium leading-none text-[#1a365d]">Email Address</label>
              <div class="relative mt-1">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" aria-hidden="true"><path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"/><rect x="2" y="4" width="20" height="16" rx="2"/></svg>
                <input
                  id="email"
                  v-model="email"
                  name="email"
                  type="email"
                  autocomplete="email"
                  required
                  placeholder="Enter your email"
                  class="flex h-9 w-full rounded-md border border-gray-300 bg-transparent px-3 py-1 text-base shadow-sm transition-colors placeholder:text-gray-400 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm pl-10"
                />
              </div>
            </div>

            <!-- Password -->
            <div>
              <label for="password" class="text-sm font-medium leading-none text-[#1a365d]">Password</label>
              <div class="relative mt-1">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" aria-hidden="true"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                <input
                  id="password"
                  v-model="password"
                  name="password"
                  :type="showPassword ? 'text' : 'password'"
                  autocomplete="current-password"
                  required
                  placeholder="Enter your password"
                  class="flex h-9 w-full rounded-md border border-gray-300 bg-transparent px-3 py-1 text-base shadow-sm transition-colors placeholder:text-gray-400 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm pl-10 pr-10"
                />
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  :aria-label="showPassword ? 'Hide password' : 'Show password'"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5" aria-hidden="true"><path d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0"/><circle cx="12" cy="12" r="3"/></svg>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5" aria-hidden="true"><path d="M9.88 9.88a3 3 0 1 0 4.24 4.24"/><path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"/><path d="M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"/><line x1="2" x2="22" y1="2" y2="22"/></svg>
                </button>
              </div>
            </div>

            <!-- Remember / Forgot -->
            <div class="flex items-center justify-between">
              <label class="flex items-center gap-2 cursor-pointer">
                <input v-model="rememberMe" type="checkbox" class="rounded border-gray-300 text-[#f5a623] focus:ring-[#f5a623]"/>
                <span class="text-sm text-gray-600">Remember me</span>
              </label>
              <router-link to="/forgot-password" class="text-sm text-[#f5a623] hover:underline">Forgot password?</router-link>
            </div>

            <!-- Submit -->
            <button
              type="submit"
              :disabled="loading"
              class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623] disabled:pointer-events-none disabled:opacity-60 shadow w-full bg-[#f5a623] hover:bg-[#e09612] text-white py-3 text-lg font-semibold"
            >
              <svg v-if="loading" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 animate-spin" aria-hidden="true"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5" aria-hidden="true"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" x2="3" y1="12" y2="12"/></svg>
              {{ loading ? 'Signing In…' : 'Sign In' }}
            </button>
          </form>

          <div class="mt-6 text-center">
            <p class="text-gray-600">
              Don't have an account?
              <router-link to="/register" class="text-[#f5a623] font-semibold hover:underline">Create Account</router-link>
            </p>
          </div>
        </div>

        <!-- Back to home -->
        <div class="text-center mt-6">
          <router-link to="/" class="text-gray-500 hover:text-[#1a365d] text-sm inline-flex items-center justify-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 rotate-180" aria-hidden="true"><path d="m9 18 6-6-6-6"/></svg>
            Back to Home
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import { getSettings } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const rememberMe = ref(false)
const showPassword = ref(false)
const errorMessage = ref('')
const loading = ref(false)
const siteLogoUrl = ref('')

const logoSrc = computed(() => siteLogoUrl.value ? mediaUrl(siteLogoUrl.value) : '/main-logo.png')

async function loadLogo() {
  try {
    const r = await getSettings('site')
    const v = r.data?.data?.value
    if (v?.logoUrl) siteLogoUrl.value = v.logoUrl
  } catch { /* keep bundled fallback */ }
}

async function handleLogin() {
  errorMessage.value = ''
  loading.value = true
  const result = await authStore.login(email.value, password.value)
  loading.value = false
  if (result.success) {
    const redirect = typeof route.query.redirect === 'string' && route.query.redirect.startsWith('/') ? route.query.redirect : ''
    router.push(redirect || (authStore.isApplicant ? '/my-portal' : '/dashboard'))
  } else {
    errorMessage.value = result.message || 'Login failed. Please try again.'
  }
}

onMounted(loadLogo)
</script>
