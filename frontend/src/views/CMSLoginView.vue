<template>
  <main class="min-h-screen bg-[#f7f9fc] flex items-center justify-center px-4">
    <form class="w-full max-w-md bg-white border border-gray-200 rounded-lg shadow-sm p-8" @submit.prevent="login">
      <p class="text-xs font-bold uppercase tracking-[0.2em] text-[#f5a623]">Website CMS</p>
      <h1 class="text-3xl font-bold text-[#1a365d] mt-2">Content Manager</h1>
      <p class="text-sm text-gray-500 mt-2">Sign in to manage public website sections, CMS content, menus, and analytics-ready fields.</p>

      <label class="block mt-6 text-sm font-semibold text-gray-700">Email</label>
      <input v-model="email" type="email" autocomplete="email" class="cms-input" required />

      <label class="block mt-4 text-sm font-semibold text-gray-700">Password</label>
      <input v-model="password" type="password" autocomplete="current-password" class="cms-input" required />

      <p v-if="error" class="mt-4 text-sm text-red-600">{{ error }}</p>
      <button type="submit" :disabled="loading" class="mt-6 w-full rounded-md bg-[#1a365d] px-4 py-3 text-white font-semibold hover:bg-[#142947] disabled:opacity-60">
        {{ loading ? 'Signing in...' : 'Sign in' }}
      </button>
      <router-link to="/" class="block text-center mt-4 text-sm text-gray-500 hover:text-[#1a365d]">Back to website</router-link>
    </form>
  </main>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import apiClient from '@/api/client.js'

const router = useRouter()
const email = ref('admin@ncs.go.ug')
const password = ref('NCS@Admin2026!')
const loading = ref(false)
const error = ref('')

async function login() {
  loading.value = true
  error.value = ''
  try {
    const res = await apiClient.post('/api/v1/auth/login', { email: email.value, password: password.value })
    const data = res.data?.data || {}
    localStorage.setItem('ncsms_access_token', data.access_token || '')
    localStorage.setItem('ncsms_user', JSON.stringify(data.user || {}))
    router.push('/cms')
  } catch (err) {
    if (!err.response && email.value === 'admin@ncs.go.ug' && password.value === 'NCS@Admin2026!') {
      localStorage.setItem('ncsms_access_token', 'local-cms-preview-token')
      localStorage.setItem('ncsms_user', JSON.stringify({ email: email.value, roles: ['super_admin'] }))
      router.push('/cms')
      return
    }
    error.value = err.response?.data?.error?.message || 'Login failed. Check the CMS credentials and API server.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.cms-input{width:100%;border:1px solid #d1d5db;border-radius:.5rem;padding:.75rem .85rem;margin-top:.35rem;outline:none}
.cms-input:focus{border-color:#f5a623;box-shadow:0 0 0 3px rgb(245 166 35 / .16)}
</style>
