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
          <h1 class="text-2xl font-bold text-[#1a365d]">Create Your Account</h1>
          <p class="text-gray-600 mt-1">Register to manage your NCS services</p>
        </div>

        <!-- Card -->
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <!-- Success -->
          <div v-if="success" class="mb-5 flex items-start gap-3 p-4 bg-green-50 border border-green-200 rounded-xl">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" aria-hidden="true"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
            <div>
              <p class="text-green-800 text-sm font-semibold">Account created!</p>
              <p class="text-green-700 text-sm mt-0.5">Redirecting to your portal…</p>
            </div>
          </div>

          <!-- Error -->
          <div v-if="errorMessage" class="mb-5 flex items-start gap-3 p-4 bg-red-50 border border-red-200 rounded-xl">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" aria-hidden="true"><circle cx="12" cy="12" r="10"/><line x1="12" x2="12" y1="8" y2="12"/><line x1="12" x2="12.01" y1="16" y2="16"/></svg>
            <p class="text-red-700 text-sm">{{ errorMessage }}</p>
          </div>

          <form @submit.prevent="handleRegister" class="space-y-5">
            <!-- Name row -->
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label for="first-name" class="text-sm font-medium leading-none text-[#1a365d]">First Name</label>
                <input
                  id="first-name"
                  v-model="firstName"
                  name="first_name"
                  type="text"
                  autocomplete="given-name"
                  required
                  placeholder="John"
                  class="flex h-9 w-full rounded-md border border-gray-300 bg-transparent px-3 py-1 text-base shadow-sm transition-colors placeholder:text-gray-400 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm mt-1"
                />
              </div>
              <div>
                <label for="last-name" class="text-sm font-medium leading-none text-[#1a365d]">Last Name</label>
                <input
                  id="last-name"
                  v-model="lastName"
                  name="last_name"
                  type="text"
                  autocomplete="family-name"
                  required
                  placeholder="Doe"
                  class="flex h-9 w-full rounded-md border border-gray-300 bg-transparent px-3 py-1 text-base shadow-sm transition-colors placeholder:text-gray-400 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm mt-1"
                />
              </div>
            </div>

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
                  placeholder="your@email.com"
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
                  :type="showPass ? 'text' : 'password'"
                  autocomplete="new-password"
                  required
                  minlength="8"
                  placeholder="Min. 8 characters"
                  class="flex h-9 w-full rounded-md border border-gray-300 bg-transparent px-3 py-1 text-base shadow-sm transition-colors placeholder:text-gray-400 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm pl-10 pr-10"
                />
                <button
                  type="button"
                  @click="showPass = !showPass"
                  :aria-label="showPass ? 'Hide password' : 'Show password'"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  <svg v-if="!showPass" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5" aria-hidden="true"><path d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0"/><circle cx="12" cy="12" r="3"/></svg>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5" aria-hidden="true"><path d="M9.88 9.88a3 3 0 1 0 4.24 4.24"/><path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"/><path d="M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"/><line x1="2" x2="22" y1="2" y2="22"/></svg>
                </button>
              </div>
            </div>

            <!-- Confirm password -->
            <div>
              <label for="confirm-password" class="text-sm font-medium leading-none text-[#1a365d]">Confirm Password</label>
              <div class="relative mt-1">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" aria-hidden="true"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                <input
                  id="confirm-password"
                  v-model="confirmPassword"
                  name="confirm_password"
                  :type="showPass ? 'text' : 'password'"
                  autocomplete="new-password"
                  required
                  placeholder="Repeat password"
                  :class="[
                    'flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-sm transition-colors placeholder:text-gray-400 focus-visible:outline-none focus-visible:ring-1 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm pl-10',
                    confirmPassword && password !== confirmPassword ? 'border-red-300 focus-visible:ring-red-400' : 'border-gray-300 focus-visible:ring-[#f5a623]'
                  ]"
                />
              </div>
              <p v-if="confirmPassword && password !== confirmPassword" class="text-xs text-red-500 mt-1">Passwords do not match</p>
            </div>

            <!-- Submit -->
            <button
              type="submit"
              :disabled="loading || (confirmPassword && password !== confirmPassword)"
              class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623] disabled:pointer-events-none disabled:opacity-60 shadow w-full bg-[#f5a623] hover:bg-[#e09612] text-white py-3 text-lg font-semibold"
            >
              <svg v-if="loading" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 animate-spin" aria-hidden="true"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5" aria-hidden="true"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><line x1="19" x2="19" y1="8" y2="14"/><line x1="22" x2="16" y1="11" y2="11"/></svg>
              {{ loading ? 'Creating Account…' : 'Create Account' }}
            </button>
          </form>

          <div class="mt-6 text-center">
            <p class="text-gray-600">
              Already have an account?
              <router-link to="/login" class="text-[#f5a623] font-semibold hover:underline">Sign In</router-link>
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
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import { getSettings } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

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
const siteLogoUrl = ref('')

const logoSrc = computed(() => siteLogoUrl.value ? mediaUrl(siteLogoUrl.value) : '/main-logo.png')

async function loadLogo() {
  try {
    const r = await getSettings('site')
    const v = r.data?.data?.value
    if (v?.logoUrl) siteLogoUrl.value = v.logoUrl
  } catch { /* keep bundled fallback */ }
}

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

onMounted(loadLogo)
</script>
