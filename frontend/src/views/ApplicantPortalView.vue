<template>
  <div class="bg-[#F8F9FC] min-h-screen">

    <!-- Portal header banner -->
    <div class="bg-[#112b4e] pt-24 pb-10 px-4">
      <div class="max-w-5xl mx-auto sm:px-6">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <p class="text-[#F48C06] text-xs font-semibold uppercase tracking-widest mb-1">{{ greeting }}</p>
            <h1 class="text-2xl md:text-3xl font-bold text-white">{{ welcomeName }}</h1>
            <p class="text-white/50 text-sm mt-1">Manage your sports licence applications</p>
          </div>
          <div class="flex items-center gap-3">
            <router-link
              to="/apply"
              class="inline-flex items-center gap-2 bg-[#F48C06] hover:bg-[#d47b05] text-white font-bold px-5 py-2.5 rounded-full text-sm transition-colors shadow-md"
            >
              <i class="icofont-paper-plane"></i> New Application
            </router-link>
            <button
              @click="handleLogout"
              title="Sign Out"
              class="w-10 h-10 rounded-full bg-white/10 hover:bg-white/20 flex items-center justify-center text-white/70 hover:text-white transition-colors"
            >
              <i class="icofont-sign-out text-lg"></i>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-5xl mx-auto px-4 sm:px-6 py-8 space-y-6">

      <!-- Stat cards -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div
          v-for="stat in statCards"
          :key="stat.label"
          class="bg-white rounded-2xl p-5 shadow-sm border border-gray-100 flex flex-col gap-2"
        >
          <div class="flex items-center justify-between">
            <span class="text-xs font-semibold text-gray-400 uppercase tracking-wide">{{ stat.label }}</span>
            <div :class="stat.iconBg" class="w-8 h-8 rounded-xl flex items-center justify-center">
              <i :class="[stat.icon, stat.iconColor, 'text-base leading-none']"></i>
            </div>
          </div>
          <div class="text-3xl font-bold text-darken">{{ stat.value }}</div>
        </div>
      </div>

      <!-- Recent applications -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <h2 class="font-bold text-darken">My Applications</h2>
          <div class="flex gap-1">
            <button
              v-for="tab in statusTabs"
              :key="tab.value"
              @click="activeStatus = tab.value"
              :class="activeStatus === tab.value
                ? 'bg-[#112b4e] text-white'
                : 'text-gray-500 hover:bg-gray-100'"
              class="px-3 py-1 rounded-lg text-xs font-medium transition-colors"
            >{{ tab.label }}</button>
          </div>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="p-8 space-y-3">
          <div v-for="i in 4" :key="i" class="h-14 bg-gray-100 rounded-xl animate-pulse"/>
        </div>

        <!-- Table -->
        <div v-else-if="filtered.length" class="divide-y divide-gray-50">
          <div
            v-for="app in filtered"
            :key="app.id"
            class="px-6 py-4 flex items-center gap-4 hover:bg-gray-50 transition-colors"
          >
            <!-- Reference -->
            <div class="flex-shrink-0 w-10 h-10 bg-[#F48C06]/10 rounded-xl flex items-center justify-center">
              <i class="icofont-document-folder text-[#F48C06] text-lg"></i>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap">
                <span class="font-semibold text-darken text-sm">
                  {{ app.reference_number || (app.id?.substring(0, 8)?.toUpperCase() ?? 'N/A') }}
                </span>
                <span class="text-xs text-gray-400">{{ formatType(app.form_type) }}</span>
              </div>
              <div class="text-xs text-gray-400 mt-0.5">
                Submitted {{ formatDate(app.created_at) }}
              </div>
            </div>
            <!-- Status badge -->
            <div class="flex-shrink-0 flex items-center gap-2">
              <span :class="statusClass(app.status)" class="text-xs font-semibold px-2.5 py-1 rounded-full">
                {{ app.status }}
              </span>
              <span v-if="app.payment_status" :class="app.payment_status === 'PAID' ? 'text-green-600' : 'text-gray-400'" class="text-xs font-medium">
                {{ app.payment_status }}
              </span>
            </div>
          </div>
        </div>

        <!-- Empty -->
        <div v-else class="flex flex-col items-center justify-center py-16 text-gray-400">
          <i class="icofont-document-folder text-5xl mb-3 opacity-30"></i>
          <p class="text-sm font-medium">No applications yet</p>
          <router-link
            to="/apply"
            class="mt-4 inline-flex items-center gap-1.5 text-sm font-semibold text-[#F48C06] hover:text-[#d47b05]"
          >
            Start your first application →
          </router-link>
        </div>
      </div>

      <!-- Apply for a Licence -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <div>
            <h2 class="font-bold text-darken">Start a New Application</h2>
            <p class="text-xs text-gray-400 mt-0.5">Select a licence type to begin your application</p>
          </div>
        </div>
        <div class="p-6 grid sm:grid-cols-3 gap-4">
          <router-link
            v-for="type in licenceTypes"
            :key="type.id"
            :to="`/apply?type=${type.id}`"
            class="rounded-2xl p-5 border border-gray-100 hover:shadow-md transition-all flex flex-col gap-3 group"
            :class="type.hoverBorder"
          >
            <div class="w-12 h-12 rounded-2xl flex items-center justify-center flex-shrink-0 transition-colors" :class="type.iconBg">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="type.icon"/>
              </svg>
            </div>
            <div>
              <div class="font-bold text-darken text-sm group-hover:text-accent transition-colors">{{ type.title }}</div>
              <div class="text-xs text-gray-400 mt-1 leading-relaxed">{{ type.description }}</div>
            </div>
            <span class="text-xs font-semibold text-accent flex items-center gap-1 mt-auto">Apply now →</span>
          </router-link>
        </div>
      </div>

      <!-- Profile shortcut -->
      <div class="flex justify-end">
        <router-link
          to="/profile"
          class="inline-flex items-center gap-2 text-sm text-gray-500 hover:text-darken transition-colors"
        >
          <i class="icofont-ui-user"></i> Update my profile
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import apiClient from '@/api/client.js'

const router = useRouter()
const authStore = useAuthStore()

const applications = ref([])
const loading = ref(true)
const activeStatus = ref('')

const licenceTypes = [
  {
    id: 'national_federation',
    title: 'National Federation',
    description: 'Register and license a national sports federation to organise competitions at national level.',
    iconBg: 'bg-[#112b4e]',
    hoverBorder: 'hover:border-[#112b4e]/30',
    icon: 'M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064',
  },
  {
    id: 'community_club',
    title: 'Community / Club',
    description: 'Register a community-level sports club or local association for your district.',
    iconBg: 'bg-gray-600',
    hoverBorder: 'hover:border-gray-300',
    icon: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0',
  },
  {
    id: 'renewal',
    title: 'Renew Existing Licence',
    description: 'Renew your sports organisation licence before it expires for the current season.',
    iconBg: 'bg-[#F48C06]',
    hoverBorder: 'hover:border-[#F48C06]/30',
    icon: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15',
  },
]

const statusTabs = [
  { value: '', label: 'All' },
  { value: 'pending', label: 'Pending' },
  { value: 'approved', label: 'Approved' },
  { value: 'rejected', label: 'Rejected' },
]

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 12) return 'Good Morning'
  if (h < 17) return 'Good Afternoon'
  return 'Good Evening'
})

const welcomeName = computed(() => {
  const u = authStore.user
  if (!u) return 'Welcome'
  return `${u.first_name || ''} ${u.last_name || ''}`.trim() || u.email || 'Welcome'
})

const filtered = computed(() => {
  if (!activeStatus.value) return applications.value
  return applications.value.filter(a => a.status?.toLowerCase() === activeStatus.value)
})

const statCards = computed(() => {
  const all = applications.value
  return [
    { label: 'Total', value: all.length, icon: 'icofont-document-folder', iconBg: 'bg-[#F48C06]/10', iconColor: 'text-[#F48C06]' },
    { label: 'Pending', value: all.filter(a => a.status?.toLowerCase() === 'pending').length, icon: 'icofont-clock-time', iconBg: 'bg-yellow-50', iconColor: 'text-yellow-500' },
    { label: 'Approved', value: all.filter(a => a.status?.toLowerCase() === 'approved').length, icon: 'icofont-check-circled', iconBg: 'bg-green-50', iconColor: 'text-green-500' },
    { label: 'Rejected', value: all.filter(a => a.status?.toLowerCase() === 'rejected').length, icon: 'icofont-close-circled', iconBg: 'bg-red-50', iconColor: 'text-red-500' },
  ]
})

function formatType(t) {
  if (!t) return ''
  return t.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('en-UG', { day: 'numeric', month: 'short', year: 'numeric' })
}

function statusClass(s) {
  const st = s?.toLowerCase()
  if (st === 'approved') return 'bg-green-100 text-green-700'
  if (st === 'rejected') return 'bg-red-100 text-red-700'
  if (st === 'pending') return 'bg-yellow-100 text-yellow-700'
  return 'bg-gray-100 text-gray-600'
}

async function handleLogout() {
  await authStore.logout()
  router.push('/')
}

onMounted(async () => {
  try {
    const res = await apiClient.get('/api/v1/applications', { params: { per_page: 50 } })
    applications.value = res.data.data?.items || res.data.data || []
  } catch {
    applications.value = []
  } finally {
    loading.value = false
  }
})
</script>
