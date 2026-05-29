<template>
  <aside class="fixed inset-y-0 left-0 z-40 w-64 flex flex-col bg-primary-700 shadow-xl">
    <!-- Brand -->
    <div class="flex items-center gap-3 px-4 py-4 border-b border-primary-600">
      <div class="flex-shrink-0 w-12 h-12 bg-white rounded-xl flex items-center justify-center overflow-hidden p-1">
        <img src="/main-logo.png" alt="NCS Logo" class="w-full h-full object-contain" />
      </div>
      <div>
        <div class="text-white font-bold text-base leading-none">National Council Of Sports</div>
        <div class="text-white text-primary-200 text-xs mt-0.5 leading-none">Management System</div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
      <router-link
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        custom
        v-slot="{ href, isActive, navigate }"
      >
        <a
          :href="href"
          @click="navigate"
          :class="[
            item.highlight
              ? (isActive ? 'bg-yellow-400 text-primary-900 shadow' : 'bg-yellow-500/90 hover:bg-yellow-400 text-primary-900 shadow')
              : (isActive ? 'bg-primary-600 text-white' : 'text-primary-100 hover:bg-primary-600/70 hover:text-white'),
            'group flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-all duration-150'
          ]"
        >
          <i :class="['text-lg leading-none flex-shrink-0', item.icon]"></i>
          {{ item.label }}
        </a>
      </router-link>
    </nav>

    <!-- User section + Logout -->
    <div class="border-t border-primary-600 px-3 py-3">
      <div class="flex items-center gap-3 px-3 py-2 mb-2">
        <div class="w-8 h-8 bg-primary-500 rounded-full flex items-center justify-center flex-shrink-0">
          <span class="text-white text-sm font-semibold">
            {{ userInitials }}
          </span>
        </div>
        <div class="flex-1 min-w-0">
          <div class="text-white text-sm font-medium truncate">{{ userName }}</div>
          <div class="text-primary-300 text-xs truncate">{{ userRole }}</div>
        </div>
      </div>
      <button
        @click="handleLogout"
        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-primary-100 hover:bg-red-600 hover:text-white transition-all duration-150"
      >
        <i class="icofont-logout text-lg leading-none flex-shrink-0"></i>
        Sign Out
      </button>
    </div>
  </aside>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'

const router = useRouter()
const authStore = useAuthStore()

const allNavItems = [
  { to: '/dashboard',   label: 'Dashboard',          icon: 'icofont-dashboard-web',  roles: null },
  { to: '/apply',       label: 'Start Application',  icon: 'icofont-paper-plane',    roles: null, highlight: true },
  { to: '/users',       label: 'Users',               icon: 'icofont-people',          roles: ['super_admin', 'admin'] },
  { to: '/applications',label: 'Applications',        icon: 'icofont-files-stack',     roles: null },
  { to: '/cms',         label: 'Content',             icon: 'icofont-newspaper',       roles: ['super_admin', 'admin', 'content_manager'] },
  { to: '/audit-logs',  label: 'Audit Logs',          icon: 'icofont-history',         roles: ['super_admin', 'admin'] },
  { to: '/roles',       label: 'Roles & Permissions', icon: 'icofont-safety',          roles: ['super_admin'] },
  { to: '/profile',     label: 'My Profile',          icon: 'icofont-user-alt-5',      roles: null }
]

const navItems = computed(() => {
  return allNavItems.filter(item => {
    if (!item.roles) return true
    if (authStore.isSuperAdmin) return true
    if (authStore.isAdmin && item.roles.includes('admin')) return true
    if (authStore.isContentManager && item.roles.includes('content_manager')) return true
    return false
  })
})

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

const userRole = computed(() => {
  const u = authStore.user
  if (!u || !u.roles || u.roles.length === 0) return 'User'
  const role = u.roles[0]
  const name = typeof role === 'string' ? role : role.name
  return name.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
})

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>
