<template>
  <LayoutDefault title="Dashboard">

    <!-- Loading -->
    <div v-if="loading" class="space-y-5">
      <div class="h-28 bg-gray-100 rounded-xl animate-pulse"></div>
      <div class="grid grid-cols-2 xl:grid-cols-4 gap-4">
        <div v-for="i in 4" :key="i" class="h-28 bg-gray-100 rounded-xl animate-pulse"></div>
      </div>
      <div class="h-64 bg-gray-100 rounded-xl animate-pulse"></div>
    </div>

    <div v-else class="space-y-6">

      <!-- Welcome banner -->
      <div class="relative overflow-hidden bg-gradient-to-r from-primary-700 via-primary-600 to-primary-500 rounded-xl p-6 text-white shadow-lg">
        <!-- Decorative circles -->
        <div class="absolute -top-8 -right-8 w-36 h-36 bg-white/5 rounded-full pointer-events-none"></div>
        <div class="absolute -bottom-6 right-20 w-20 h-20 bg-white/5 rounded-full pointer-events-none"></div>
        <div class="flex flex-col sm:flex-row sm:items-center gap-4">
          <div class="flex-1 min-w-0">
            <p class="text-primary-200 text-xs font-semibold uppercase tracking-widest mb-1">{{ greeting }}</p>
            <h2 class="text-xl font-bold">{{ welcomeName }}!</h2>
            <p class="text-primary-200 mt-1 text-sm">Here's what's happening with the National Council of Sports today.</p>
          </div>
        </div>
      </div>

      <!-- Error alert -->
      <div v-if="error" class="flex items-center gap-3 p-4 bg-red-50 border border-red-200 rounded-lg">
        <i class="icofont-warning-alt text-xl text-red-500 flex-shrink-0"></i>
        <p class="text-red-700 text-sm">{{ error }}</p>
      </div>

      <!-- Website Analytics (preview — telemetry pipeline coming soon) -->
      <section class="admin-card">
        <div class="admin-card-header flex items-center justify-between">
          <span class="text-sm font-semibold text-gray-800">Website Analytics</span>
          <span class="text-[10px] uppercase tracking-wider font-bold bg-amber-100 text-amber-700 px-2 py-0.5 rounded-full">Preview</span>
        </div>
        <div class="admin-card-body">
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-5">
            <div v-for="m in analyticsMetrics" :key="m.label" class="bg-gray-50 rounded-lg p-3">
              <div class="text-[10px] uppercase tracking-wider text-gray-400 font-semibold">{{ m.label }}</div>
              <div class="text-xl font-bold text-gray-900 mt-1">{{ m.value }}</div>
              <div class="text-[10px] text-gray-400 mt-0.5">{{ m.hint }}</div>
            </div>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div>
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs font-semibold text-gray-700 uppercase tracking-wider">Devices</span>
                <span class="text-[10px] text-gray-400">Mobile · Desktop · Tablet</span>
              </div>
              <div v-for="d in deviceBreakdown" :key="d.label" class="mb-2">
                <div class="flex items-center justify-between text-xs text-gray-600 mb-1">
                  <span>{{ d.label }}</span><span>{{ d.pct }}%</span>
                </div>
                <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
                  <div :style="{ width: d.pct + '%' }" :class="d.color" class="h-full rounded-full"/>
                </div>
              </div>
            </div>
            <div>
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs font-semibold text-gray-700 uppercase tracking-wider">Platforms</span>
                <span class="text-[10px] text-gray-400">OS distribution</span>
              </div>
              <div v-for="p in platformBreakdown" :key="p.label" class="mb-2">
                <div class="flex items-center justify-between text-xs text-gray-600 mb-1">
                  <span>{{ p.label }}</span><span>{{ p.pct }}%</span>
                </div>
                <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
                  <div :style="{ width: p.pct + '%' }" :class="p.color" class="h-full rounded-full"/>
                </div>
              </div>
            </div>
          </div>
          <p class="text-[11px] text-gray-400 mt-4">Telemetry pipeline is being instrumented. Numbers shown are placeholders.</p>
        </div>
      </section>

      <!-- Stat tiles -->
      <div class="grid grid-cols-2 xl:grid-cols-4 gap-4">
        <div v-for="stat in statCards" :key="stat.label" class="stat-card hover:shadow-md transition-shadow">
          <div :class="[stat.iconBg, 'flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center']">
            <i :class="['text-2xl leading-none', stat.icon, stat.iconColor]"></i>
          </div>
          <div class="min-w-0">
            <div class="text-2xl font-bold text-gray-900">{{ stat.value }}</div>
            <div class="text-xs text-gray-500 mt-0.5 leading-tight">{{ stat.label }}</div>
            <div v-if="stat.sub" class="text-xs mt-1">
              <span :class="stat.subColor">{{ stat.sub }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Middle row: status breakdown + recent applications -->
      <div class="grid grid-cols-1 xl:grid-cols-3 gap-5">

        <!-- Applications by Status -->
        <div v-if="stats?.applications_by_status && Object.keys(stats.applications_by_status).length" class="admin-card">
          <div class="admin-card-header">
            <span class="text-sm font-semibold text-gray-800">By Status</span>
          </div>
          <div class="admin-card-body space-y-2.5">
            <div v-for="(count, status) in stats.applications_by_status" :key="status" class="flex items-center gap-2">
              <StatusBadge :status="status" />
              <div class="flex-1 h-1.5 bg-gray-100 rounded-full overflow-hidden ml-1">
                <div
                  :class="statusBarColor(status)"
                  class="h-full rounded-full transition-all"
                  :style="{ width: barWidth(count) }"
                ></div>
              </div>
              <span class="text-sm font-semibold text-gray-700 w-6 text-right">{{ count }}</span>
            </div>
          </div>
        </div>

        <!-- Recent Applications -->
        <div
          class="admin-card overflow-hidden"
          :class="stats?.applications_by_status && Object.keys(stats.applications_by_status).length ? 'xl:col-span-2' : 'xl:col-span-3'"
        >
          <div class="admin-card-header">
            <span class="text-sm font-semibold text-gray-800">Recent Applications</span>
            <router-link to="/applications" class="text-xs font-medium text-primary-700 hover:text-primary-600">
              View all →
            </router-link>
          </div>
          <div class="overflow-x-auto">
            <table v-if="recentApplications.length" class="w-full">
              <thead>
                <tr>
                  <th class="table-th">Reference</th>
                  <th class="table-th">Type</th>
                  <th class="table-th">Status</th>
                  <th class="table-th">Date</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-50">
                <tr
                  v-for="app in recentApplications"
                  :key="app.id"
                  class="hover:bg-gray-50 transition-colors"
                >
                  <td class="table-td font-medium text-gray-900">
                    {{ app.reference_number || app.id?.substring(0, 8) || 'N/A' }}
                  </td>
                  <td class="table-td text-gray-600">{{ formatFormType(app.form_type) }}</td>
                  <td class="table-td"><StatusBadge :status="app.status" /></td>
                  <td class="table-td text-gray-500">{{ formatDate(app.submitted_at || app.created_at) }}</td>
                </tr>
              </tbody>
            </table>
            <div v-else class="flex flex-col items-center justify-center py-12 text-gray-400">
              <i class="icofont-files-stack text-5xl text-gray-200 mb-3"></i>
              <p class="text-sm">No recent applications</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="admin-card">
        <div class="admin-card-header">
          <span class="text-sm font-semibold text-gray-800">Quick Actions</span>
        </div>
        <div class="admin-card-body">
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <router-link
              v-for="action in quickActions"
              :key="action.label"
              :to="action.to"
              class="flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-100 hover:border-primary-200 hover:bg-primary-50/50 transition-all group text-center"
            >
              <div :class="[action.iconBg, 'w-10 h-10 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform']">
                <i :class="['text-xl leading-none', action.icon, action.iconColor]"></i>
              </div>
              <span class="text-xs font-medium text-gray-600 group-hover:text-primary-700 transition-colors leading-tight">{{ action.label }}</span>
            </router-link>
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
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'
import apiClient from '@/api/client.js'

const authStore = useAuthStore()
const breadcrumbStore = useBreadcrumbStore()
const loading = ref(true)
const error = ref('')
const stats = ref(null)
const recentApplications = ref([])

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 12) return 'Good morning'
  if (h < 17) return 'Good afternoon'
  return 'Good evening'
})

const welcomeName = computed(() => {
  const u = authStore.user
  if (!u) return 'there'
  if (u.first_name) return u.first_name
  if (u.email) return u.email.split('@')[0]
  return 'there'
})

const statCards = computed(() => [
  {
    label: 'Total Applications',
    value: stats.value?.total_applications ?? '—',
    icon: 'icofont-files-stack',
    iconBg: 'bg-primary-50',
    iconColor: 'text-primary-700',
    sub: null
  },
  {
    label: 'Pending Review',
    value: stats.value?.pending_review ?? '—',
    icon: 'icofont-clock-time',
    iconBg: 'bg-yellow-50',
    iconColor: 'text-yellow-600',
    sub: 'Awaiting action',
    subColor: 'text-yellow-600'
  },
  {
    label: 'Needs Attention',
    value: stats.value?.needs_attention ?? '—',
    icon: 'icofont-warning-alt',
    iconBg: 'bg-red-50',
    iconColor: 'text-red-500',
    sub: null
  },
  {
    label: 'Total Users',
    value: stats.value?.total_users ?? '—',
    icon: 'icofont-users-alt-5',
    iconBg: 'bg-blue-50',
    iconColor: 'text-blue-600',
    sub: null
  }
])

const analyticsMetrics = ref([
  { label: 'Unique Visitors',  value: '—', hint: 'last 30 days' },
  { label: 'Page Views',       value: '—', hint: 'last 30 days' },
  { label: 'Avg. Session',     value: '—', hint: 'minutes' },
  { label: 'Bounce Rate',      value: '—', hint: 'home + landing' },
])
const deviceBreakdown = ref([
  { label: 'Mobile',  pct: 0, color: 'bg-primary-500' },
  { label: 'Desktop', pct: 0, color: 'bg-indigo-500' },
  { label: 'Tablet',  pct: 0, color: 'bg-amber-500' },
])
const platformBreakdown = ref([
  { label: 'Android', pct: 0, color: 'bg-green-500' },
  { label: 'iOS',     pct: 0, color: 'bg-gray-700' },
  { label: 'Windows', pct: 0, color: 'bg-blue-500' },
  { label: 'macOS',   pct: 0, color: 'bg-gray-500' },
  { label: 'Linux',   pct: 0, color: 'bg-orange-500' },
])

const quickActions = computed(() => {
  const actions = [
    { label: 'All Applications',to: '/applications', icon: 'icofont-files-stack', iconBg: 'bg-primary-50', iconColor: 'text-primary-700' },
    { label: 'My Profile',      to: '/profile',      icon: 'icofont-user-alt-5',  iconBg: 'bg-indigo-50',  iconColor: 'text-indigo-600' },
  ]
  if (authStore.isAdmin || authStore.isSuperAdmin) {
    actions.splice(2, 0, { label: 'Manage Users', to: '/users', icon: 'icofont-people', iconBg: 'bg-green-50', iconColor: 'text-green-600' })
  }
  return actions.slice(0, 4)
})

const maxCount = computed(() => {
  if (!stats.value?.applications_by_status) return 1
  return Math.max(...Object.values(stats.value.applications_by_status), 1)
})

function barWidth(count) {
  return `${Math.round((count / maxCount.value) * 100)}%`
}

function statusBarColor(status) {
  const map = {
    approved: 'bg-green-500', pending: 'bg-yellow-500', submitted: 'bg-blue-500',
    rejected: 'bg-red-500', draft: 'bg-gray-400', reviewing: 'bg-indigo-500'
  }
  return map[status?.toLowerCase()] || 'bg-gray-400'
}

function formatDate(dateStr) {
  if (!dateStr) return 'N/A'
  try {
    return new Date(dateStr).toLocaleDateString('en-UG', { day: '2-digit', month: 'short', year: 'numeric' })
  } catch { return dateStr }
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
    if (dashRes.status === 'fulfilled') stats.value = dashRes.value.data.data
    if (appsRes.status === 'fulfilled') {
      recentApplications.value = Array.isArray(appsRes.value.data.data) ? appsRes.value.data.data : []
    } else {
      try {
        const userApps = await apiClient.get('/api/v1/applications')
        recentApplications.value = Array.isArray(userApps.data.data) ? userApps.data.data.slice(0, 5) : []
      } catch { recentApplications.value = [] }
    }
  } catch (err) {
    error.value = err.response?.data?.error?.message || 'Failed to load dashboard data.'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  breadcrumbStore.set('Dashboard')
  loadDashboard()
})
</script>
