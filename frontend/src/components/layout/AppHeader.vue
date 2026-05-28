<template>
  <header class="sticky top-0 z-30 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
    <div class="flex items-center gap-3">
      <h1 class="text-xl font-semibold text-gray-800">{{ title }}</h1>
    </div>

    <div class="flex items-center gap-4">
      <!-- Notification indicator -->
      <button class="relative p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors">
        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0" />
        </svg>
      </button>

      <!-- User avatar + name -->
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center">
          <span class="text-white text-sm font-semibold">{{ userInitials }}</span>
        </div>
        <div class="hidden sm:block">
          <div class="text-sm font-medium text-gray-700">{{ userName }}</div>
          <div class="flex items-center gap-1">
            <span
              v-for="role in userRoles"
              :key="role"
              class="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-primary-100 text-primary-700"
            >
              {{ role }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth.js'

defineProps({
  title: {
    type: String,
    default: 'Dashboard'
  }
})

const authStore = useAuthStore()

const userName = computed(() => {
  const u = authStore.user
  if (!u) return 'User'
  if (u.first_name || u.last_name) return `${u.first_name || ''} ${u.last_name || ''}`.trim()
  return u.email || 'User'
})

const userInitials = computed(() => {
  const u = authStore.user
  if (!u) return 'U'
  if (u.first_name && u.last_name) return `${u.first_name[0]}${u.last_name[0]}`.toUpperCase()
  if (u.first_name) return u.first_name[0].toUpperCase()
  if (u.email) return u.email[0].toUpperCase()
  return 'U'
})

const userRoles = computed(() => {
  const u = authStore.user
  if (!u || !u.roles) return []
  return u.roles
    .slice(0, 2)
    .map(r => {
      const name = typeof r === 'string' ? r : r.name
      return name.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
    })
})
</script>
