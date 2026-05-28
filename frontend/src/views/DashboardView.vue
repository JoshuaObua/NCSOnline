<template>
  <LayoutDefault title="Dashboard">
    <!-- Loading state -->
    <div v-if="loading" class="flex items-center justify-center h-64">
      <div class="flex flex-col items-center gap-3">
        <svg class="animate-spin w-8 h-8 text-primary-600" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="text-gray-500 text-sm">Loading dashboard...</p>
      </div>
    </div>

    <div v-else class="space-y-6">
      <!-- Welcome Banner -->
      <div class="bg-gradient-to-r from-primary-700 to-primary-600 rounded-xl p-6 text-white">
        <h2 class="text-xl font-semibold">Welcome back, {{ welcomeName }}</h2>
        <p class="text-primary-100 mt-1 text-sm">Here's what's happening with the National Council of Sports today.</p>
      </div>

      <!-- Error Alert -->
      <div v-if="error" class="flex items-center gap-3 p-4 bg-red-50 border border-red-200 rounded-lg">
        <i class="icofont-warning-alt text-xl text-red-500 flex-shrink-0"></i>
        <p class="text-red-700 text-sm">{{ error }}</p>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-5">
        <div
          v-for="stat in statCards"
          :key="stat.label"
          class="bg-white rounded-xl border border-gray-200 p-5 flex items-start gap-4 shadow-sm hover:shadow-md transition-shadow"
        >
          <div :class="stat.iconBg" class="flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center">
            <i :class="['text-2xl leading-none', stat.icon, stat.iconColor]"></i>
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-900">{{ stat.value }}</div>
            <div class="text-sm text-gray-500 mt-0.5">{{ stat.label }}</div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 xl:grid-cols-3 gap-6">
        <!-- Applications by Status -->
        <div v-if="stats && stats.applications_by_status" class="bg-white rounded-xl border border-gray-200 shadow-sm">
          <div class="px-6 py-4 border-b border-gray-100">
            <h3 class="font-semibold text-gray-800">Applications by Status</h3>
          </div>
          <div class="p-6 space-y-3">
            <div
              v-for="(count, status) in stats.applications_by_status"
              :key="status"
              class="flex items-center justify-between"
            >
              <div class="flex items-center gap-2">
                <StatusBadge :status="status" />
              </div>
              <span class="text-sm font-semibold text-gray-700">{{ count }}</span>
            </div>
            <div v-if="!Object.keys(stats.applications_by_status).length" class="text-center py-4 text-gray-400 text-sm">
              No application data yet
            </div>
          </div>
        </div>

        <!-- Recent Applications -->
        <div class="bg-white rounded-xl border border-gray-200 shadow-sm" :class="stats?.applications_by_status ? 'xl:col-span-2' : 'xl:col-span-3'">
          <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
            <h3 class="font-semibold text-gray-800">Recent Applications</h3>
            <router-link to="/applications" class="text-sm text-primary-600 hover:text-primary-700 font-medium">
              View all
            </router-link>
          </div>
          <div class="overflow-x-auto">
            <table v-if="recentApplications.length" class="w-full">
              <thead>
                <tr class="border-b border-gray-100">
                  <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Reference</th>
                  <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Type</th>
                  <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Status</th>
                  <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Date</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-50">
                <tr
                  v-for="app in recentApplications"
                  :key="app.id"
                  class="hover:bg-gray-50 transition-colors"
                >
                  <td class="px-6 py-3 text-sm font-medium text-gray-900">
                    {{ app.reference_number || app.id?.substring(0, 8) || 'N/A' }}
                  </td>
                  <td class="px-6 py-3 text-sm text-gray-600">
                    {{ formatFormType(app.form_type) }}
                  </td>
                  <td class="px-6 py-3">
                    <StatusBadge :status="app.status" />
                  </td>
                  <td class="px-6 py-3 text-sm text-gray-500">
                    {{ formatDate(app.submitted_at || app.created_at) }}
                  </td>
                </tr>
              </tbody>
            </table>
            <div v-else class="flex flex-col items-center justify-center py-12 text-gray-400">
              <i class="icofont-files-stack text-5xl text-gray-300 mb-3"></i>
              <p class="text-sm">No recent applications</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { useAuthStore } from '@/stores/auth.js'
import apiClient from '@/api/client.js'

const authStore = useAuthStore()
const loading = ref(true)
const error = ref('')
const stats = ref(null)
const recentApplications = ref([])

const welcomeName = computed(() => {
  const u = authStore.user
  if (!u) return 'there'
  if (u.first_name) return u.first_name
  if (u.email) return u.email.split('@')[0]
  return 'there'
})

const statCards = computed(() => [
  {
    label: 'Total Users',
    value: stats.value?.total_users ?? '—',
    icon: 'icofont-users-alt-5',
    iconBg: 'bg-blue-50',
    iconColor: 'text-blue-600'
  },
  {
    label: 'Total Applications',
    value: stats.value?.total_applications ?? '—',
    icon: 'icofont-files-stack',
    iconBg: 'bg-primary-50',
    iconColor: 'text-primary-700'
  },
  {
    label: 'Pending Review',
    value: stats.value?.pending_review ?? '—',
    icon: 'icofont-clock-time',
    iconBg: 'bg-yellow-50',
    iconColor: 'text-yellow-600'
  },
  {
    label: 'Needs Attention',
    value: stats.value?.needs_attention ?? '—',
    icon: 'icofont-warning-alt',
    iconBg: 'bg-red-50',
    iconColor: 'text-red-500'
  }
])

function formatDate(dateStr) {
  if (!dateStr) return 'N/A'
  try {
    return new Date(dateStr).toLocaleDateString('en-UG', {
      day: '2-digit', month: 'short', year: 'numeric'
    })
  } catch {
    return dateStr
  }
}

function formatFormType(type) {
  if (!type) return 'N/A'
  return type.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

async function loadDashboard() {
  loading.value = true
  error.value = ''
  try {
    const [dashRes, appsRes] = await Promise.allSettled([
      apiClient.get('/api/v1/admin/dashboard'),
      apiClient.get('/api/v1/admin/applications?page=1&per_page=5')
    ])

    if (dashRes.status === 'fulfilled') {
      stats.value = dashRes.value.data.data
    }

    if (appsRes.status === 'fulfilled') {
      const appData = appsRes.value.data
      recentApplications.value = Array.isArray(appData.data) ? appData.data : []
    } else {
      // Fall back to user's own applications
      try {
        const userApps = await apiClient.get('/api/v1/applications')
        recentApplications.value = Array.isArray(userApps.data.data)
          ? userApps.data.data.slice(0, 5)
          : []
      } catch {
        recentApplications.value = []
      }
    }
  } catch (err) {
    error.value = err.response?.data?.error?.message || 'Failed to load dashboard data.'
  } finally {
    loading.value = false
  }
}

onMounted(loadDashboard)
</script>
