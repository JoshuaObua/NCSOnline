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

      <!-- Quick actions -->
      <div class="grid sm:grid-cols-3 gap-4">
        <router-link
          to="/apply"
          class="bg-white rounded-2xl p-5 shadow-sm border border-gray-100 hover:border-[#F48C06]/30 hover:shadow-md transition-all flex items-center gap-4 group"
        >
          <div class="w-12 h-12 bg-[#F48C06]/10 group-hover:bg-[#F48C06] rounded-2xl flex items-center justify-center transition-colors flex-shrink-0">
            <i class="icofont-paper-plane text-[#F48C06] group-hover:text-white text-xl transition-colors"></i>
          </div>
          <div>
            <div class="font-bold text-darken text-sm">Apply for Licence</div>
            <div class="text-xs text-gray-400 mt-0.5">New application</div>
          </div>
        </router-link>

        <router-link
          to="/apply"
          class="bg-white rounded-2xl p-5 shadow-sm border border-gray-100 hover:border-[#112b4e]/30 hover:shadow-md transition-all flex items-center gap-4 group"
        >
          <div class="w-12 h-12 bg-[#112b4e]/8 group-hover:bg-[#112b4e] rounded-2xl flex items-center justify-center transition-colors flex-shrink-0">
            <i class="icofont-refresh text-[#112b4e] group-hover:text-white text-xl transition-colors"></i>
          </div>
          <div>
            <div class="font-bold text-darken text-sm">Renew Licence</div>
            <div class="text-xs text-gray-400 mt-0.5">Renew existing licence</div>
          </div>
        </router-link>

        <router-link
          to="/profile"
          class="bg-white rounded-2xl p-5 shadow-sm border border-gray-100 hover:border-sky-200 hover:shadow-md transition-all flex items-center gap-4 group"
        >
          <div class="w-12 h-12 bg-sky-50 group-hover:bg-sky-500 rounded-2xl flex items-center justify-center transition-colors flex-shrink-0">
            <i class="icofont-ui-user text-sky-500 group-hover:text-white text-xl transition-colors"></i>
          </div>
          <div>
            <div class="font-bold text-darken text-sm">My Profile</div>
            <div class="text-xs text-gray-400 mt-0.5">Update your details</div>
          </div>
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
