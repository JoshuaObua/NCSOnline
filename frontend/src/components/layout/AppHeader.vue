<template>
  <header class="sticky top-0 z-30 bg-white border-b border-gray-200 flex-shrink-0">
    <div class="flex items-center gap-3 px-4 py-3">

      <!-- Sidebar toggle -->
      <button
        @click="$emit('toggle-sidebar')"
        class="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors flex-shrink-0"
        aria-label="Toggle sidebar"
      >
        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12" />
        </svg>
      </button>

      <!-- Page title -->
      <h1 class="text-base font-semibold text-gray-800 flex-1 truncate">{{ pageTitle }}</h1>

      <!-- Right actions -->
      <div class="flex items-center gap-1 flex-shrink-0">

        <!-- Search toggle -->
        <button
          @click="searchOpen = !searchOpen"
          class="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
          </svg>
        </button>

        <!-- Notifications -->
        <div class="relative" ref="notifRef">
          <button
            @click="toggleNotif"
            class="relative p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0" />
            </svg>
            <span v-if="unreadCount > 0" class="absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full"></span>
          </button>

          <!-- Notification dropdown -->
          <Transition name="dropdown">
            <div
              v-if="notifOpen"
              class="absolute right-0 top-full mt-2 w-80 bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden z-50"
            >
              <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
                <span class="text-sm font-semibold text-gray-800">Notifications</span>
                <span v-if="unreadCount > 0" class="text-xs font-medium text-white bg-red-500 rounded-full px-2 py-0.5">{{ unreadCount }}</span>
              </div>
              <div v-if="loadingNotifs" class="py-8 flex justify-center">
                <svg class="animate-spin w-5 h-5 text-primary-600" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                </svg>
              </div>
              <div v-else-if="notifications.length === 0" class="py-8 text-center text-gray-400 text-sm">
                <i class="icofont-inbox text-3xl text-gray-300 block mb-2"></i>
                No recent activity
              </div>
              <div v-else>
                <div
                  v-for="n in notifications"
                  :key="n.id"
                  class="px-4 py-3 hover:bg-gray-50 transition-colors border-b border-gray-50 last:border-0"
                >
                  <div class="flex items-start gap-3">
                    <div :class="severityDot(n.severity)" class="w-2 h-2 rounded-full mt-1.5 flex-shrink-0"></div>
                    <div class="min-w-0">
                      <p class="text-xs font-medium text-gray-800 truncate">{{ n.action }}</p>
                      <p class="text-xs text-gray-400 mt-0.5">{{ formatRelTime(n.created_at) }}</p>
                    </div>
                  </div>
                </div>
              </div>
              <div class="px-4 py-2.5 border-t border-gray-100 bg-gray-50/60">
                <router-link
                  to="/audit-logs"
                  @click="notifOpen = false"
                  class="text-xs font-medium text-primary-700 hover:text-primary-600"
                >View all activity →</router-link>
              </div>
            </div>
          </Transition>
        </div>

        <!-- User avatar dropdown -->
        <div class="relative" ref="userRef">
          <button
            @click="toggleUser"
            class="flex items-center gap-2 px-2 py-1.5 rounded-lg hover:bg-gray-100 transition-colors"
          >
            <div class="w-8 h-8 bg-primary-700 rounded-full flex items-center justify-center flex-shrink-0">
              <span class="text-white text-xs font-bold">{{ userInitials }}</span>
            </div>
            <div class="hidden sm:block text-left">
              <div class="text-sm font-medium text-gray-700 leading-tight">{{ userName }}</div>
            </div>
            <svg class="w-4 h-4 text-gray-400 hidden sm:block" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
            </svg>
          </button>

          <Transition name="dropdown">
            <div
              v-if="userOpen"
              class="absolute right-0 top-full mt-2 w-52 bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden z-50"
            >
              <div class="px-4 py-3 border-b border-gray-100">
                <p class="text-sm font-semibold text-gray-800">{{ userName }}</p>
                <p class="text-xs text-gray-400 truncate mt-0.5">{{ userEmail }}</p>
              </div>
              <div class="py-1">
                <router-link
                  to="/profile"
                  @click="userOpen = false"
                  class="flex items-center gap-3 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                >
                  <i class="icofont-user-alt-5 text-base text-gray-400"></i>
                  My Profile
                </router-link>
              </div>
              <div class="border-t border-gray-100 py-1">
                <button
                  @click="handleLogout"
                  class="w-full flex items-center gap-3 px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                >
                  <i class="icofont-logout text-base"></i>
                  Sign Out
                </button>
              </div>
            </div>
          </Transition>
        </div>
      </div>
    </div>

    <!-- Search bar (collapsible) -->
    <Transition name="slide-down">
      <div v-if="searchOpen" class="border-t border-gray-100 px-4 py-2">
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
            </svg>
          </div>
          <input
            ref="searchInput"
            v-model="searchQuery"
            @keydown.escape="searchOpen = false"
            type="text"
            placeholder="Search anything..."
            class="w-full pl-9 pr-4 py-2 text-sm bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-700 focus:border-primary-700"
          />
        </div>
      </div>
    </Transition>
  </header>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import apiClient from '@/api/client.js'

defineProps({
  pageTitle: { type: String, default: '' }
})
defineEmits(['toggle-sidebar'])

const router = useRouter()
const authStore = useAuthStore()

// Search
const searchOpen = ref(false)
const searchQuery = ref('')
const searchInput = ref(null)
watch(searchOpen, async (v) => { if (v) { await nextTick(); searchInput.value?.focus() } })

// Notifications
const notifOpen = ref(false)
const notifRef = ref(null)
const loadingNotifs = ref(false)
const notifications = ref([])
const unreadCount = ref(0)

async function fetchNotifications() {
  if (notifications.value.length) return
  loadingNotifs.value = true
  try {
    const res = await apiClient.get('/api/v1/admin/audit-logs?per_page=5')
    const data = res.data?.data ?? []
    notifications.value = Array.isArray(data) ? data : []
    unreadCount.value = notifications.value.length
  } catch {
    notifications.value = []
  } finally {
    loadingNotifs.value = false
  }
}

function toggleNotif() {
  notifOpen.value = !notifOpen.value
  userOpen.value = false
  if (notifOpen.value) {
    fetchNotifications()
    unreadCount.value = 0
  }
}

function severityDot(sev) {
  if (sev === 'CRITICAL') return 'bg-red-500'
  if (sev === 'WARNING')  return 'bg-yellow-500'
  if (sev === 'ERROR')    return 'bg-orange-500'
  return 'bg-blue-400'
}

function formatRelTime(dateStr) {
  if (!dateStr) return ''
  const diff = Date.now() - new Date(dateStr).getTime()
  const m = Math.floor(diff / 60000)
  if (m < 1)  return 'just now'
  if (m < 60) return `${m}m ago`
  const h = Math.floor(m / 60)
  if (h < 24) return `${h}h ago`
  return `${Math.floor(h / 24)}d ago`
}

// User dropdown
const userOpen = ref(false)
const userRef = ref(null)

function toggleUser() {
  userOpen.value = !userOpen.value
  notifOpen.value = false
}

async function handleLogout() {
  userOpen.value = false
  await authStore.logout()
  router.push('/login')
}

// Click-outside close
function onClickOutside(e) {
  if (notifRef.value && !notifRef.value.contains(e.target)) notifOpen.value = false
  if (userRef.value  && !userRef.value.contains(e.target))  userOpen.value  = false
}
onMounted(() => document.addEventListener('click', onClickOutside, true))
onUnmounted(() => document.removeEventListener('click', onClickOutside, true))

// User info
const userName = computed(() => {
  const u = authStore.user
  if (!u) return 'User'
  if (u.first_name || u.last_name) return `${u.first_name || ''} ${u.last_name || ''}`.trim()
  return u.email || 'User'
})
const userEmail = computed(() => authStore.user?.email || '')
const userInitials = computed(() => {
  const u = authStore.user
  if (!u) return 'U'
  if (u.first_name && u.last_name) return `${u.first_name[0]}${u.last_name[0]}`.toUpperCase()
  if (u.first_name) return u.first_name[0].toUpperCase()
  if (u.email) return u.email[0].toUpperCase()
  return 'U'
})
</script>

<style scoped>
.dropdown-enter-active, .dropdown-leave-active { transition: opacity 0.15s, transform 0.15s; }
.dropdown-enter-from, .dropdown-leave-to { opacity: 0; transform: translateY(-6px); }

.slide-down-enter-active, .slide-down-leave-active { transition: max-height 0.2s ease, opacity 0.2s; }
.slide-down-enter-from, .slide-down-leave-to { max-height: 0; opacity: 0; }
.slide-down-enter-to, .slide-down-leave-from { max-height: 80px; opacity: 1; }
</style>
