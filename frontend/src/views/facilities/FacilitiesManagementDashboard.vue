<template>
  <LayoutDefault title="Facilities & Venues Management">
    <div class="space-y-6 pb-12">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-building text-xs"></i> National Sports Venues
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Facilities, Venues & Hostels Directorate</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Public booking calendar, event tariffs, statutory NTR generation, and athlete residential hostel occupancy tracking.
          </p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <button @click="fetchData" class="btn btn-sm btn-light flex items-center gap-1.5" :disabled="loading">
            <i class="icofont-refresh" :class="{ 'animate-spin': loading }"></i> Refresh
          </button>
        </div>
      </div>

      <!-- Macro Metric KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Active Bookings</h5>
              <h2>{{ bookings.length }}</h2>
              <span class="badge badge-primary"><i class="icofont-calendar"></i> Current Fixtures</span>
            </div>
            <div class="banner-img">
              <i class="icofont-calendar"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Tariff Inflows</h5>
              <h2>UGX {{ Number(totalTariffs).toLocaleString() }}</h2>
              <span class="badge badge-success"><i class="icofont-bank"></i> NTR Inflows</span>
            </div>
            <div class="banner-img">
              <i class="icofont-bank"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Hostel Residents</h5>
              <h2>{{ totalAthletesHoused }}</h2>
              <span class="badge badge-info"><i class="icofont-home"></i> Athletes In-Camp</span>
            </div>
            <div class="banner-img">
              <i class="icofont-home"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Venues Registered</h5>
              <h2>{{ uniqueVenuesCount }}</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Operational</span>
            </div>
            <div class="banner-img">
              <i class="icofont-building-alt"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- Navigation Tabs -->
      <div class="flex items-center gap-1.5 border-b border-slate-200 dark:border-slate-800 pb-2">
        <button 
          v-for="tab in tabs" 
          :key="tab.id"
          @click="activeTab = tab.id"
          class="px-3.5 py-2 rounded-lg text-xs font-bold transition flex items-center gap-2"
          :class="activeTab === tab.id 
            ? 'bg-blue-600 text-white shadow-sm' 
            : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800'"
        >
          <i :class="tab.icon"></i> {{ tab.label }}
        </button>
      </div>

      <!-- TAB 1: VENUE BOOKINGS -->
      <div v-if="activeTab === 'bookings'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Active Venue Bookings & Fixtures Schedule</h4>
          </div>

          <div class="divide-y divide-slate-100 dark:divide-slate-800">
            <div v-for="b in bookings" :key="b.id" class="p-5 hover:bg-slate-50/70 dark:hover:bg-slate-800/40 transition">
              <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                <div class="space-y-1.5 flex-1">
                  <div class="flex items-center gap-2 flex-wrap">
                    <span class="font-mono text-xs font-bold text-blue-600 dark:text-blue-400 bg-slate-100 dark:bg-slate-800 px-2 py-0.5 rounded">{{ b.booking_reference }}</span>
                    <span class="badge badge-info">{{ b.venue_name }}</span>
                    <span class="badge badge-success">{{ b.booking_status }}</span>
                  </div>
                  <h4 class="text-sm font-bold text-slate-900 dark:text-white">{{ b.event_title }}</h4>
                  <p class="text-xs text-slate-500">Applicant: <strong class="text-slate-700 dark:text-slate-300">{{ b.applicant_name }}</strong> ({{ b.applicant_type }})</p>
                  <div class="flex items-center gap-4 text-xs text-slate-500 pt-1">
                    <span><i class="icofont-calendar"></i> Date: <strong>{{ formatDate(b.start_date) }}</strong></span>
                    <span v-if="b.tariff_amount_ugx > 0"><i class="icofont-money"></i> Tariff: <strong class="text-blue-600 dark:text-blue-400 font-mono">UGX {{ Number(b.tariff_amount_ugx).toLocaleString() }}</strong></span>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="!bookings.length" class="text-center py-8 text-xs text-slate-400">
              No venue bookings recorded.
            </div>
          </div>
        </div>
      </div>

      <!-- TAB 2: HOSTEL OCCUPANCY -->
      <div v-if="activeTab === 'hostels'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Athlete Hostel Residential Encampment Log</h4>
          </div>
          <div class="card-body p-5 space-y-3">
            <div v-for="h in hostels" :key="h.id" class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30">
              <div class="flex items-center justify-between">
                <span class="font-mono text-xs font-bold text-blue-600 dark:text-blue-400">{{ h.block_name }} - Room {{ h.room_number }}</span>
                <span class="badge badge-success">{{ h.status }}</span>
              </div>
              <h5 class="text-sm font-bold text-slate-900 dark:text-white mt-1">{{ h.team_or_athlete_name }}</h5>
              <p class="text-xs text-slate-500">{{ h.federation_name }} · {{ h.athletes_count }} Athletes Housed</p>
              <div class="text-xs text-slate-500 pt-2 mt-2 border-t border-slate-200 dark:border-slate-700/50">
                Check-in: <strong>{{ formatDate(h.check_in_date) }}</strong> · Expected Check-out: <strong>{{ formatDate(h.check_out_date) }}</strong>
              </div>
            </div>
            <div v-if="!hostels.length" class="text-center py-6 text-xs text-slate-400">
              No hostel occupancies currently recorded.
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
import { apiGet } from '@/api/client'

const loading = ref(false)
const bookings = ref([])
const hostels = ref([])
const activeTab = ref('bookings')

const tabs = [
  { id: 'bookings', label: 'Venue Bookings & Tariffs', icon: 'icofont-calendar' },
  { id: 'hostels', label: 'Athlete Hostel Encampments', icon: 'icofont-home' }
]

const totalTariffs = computed(() => {
  return bookings.value.reduce((acc, cur) => acc + (cur.tariff_amount_ugx || 0), 0)
})

const totalAthletesHoused = computed(() => {
  return hostels.value.reduce((acc, cur) => acc + (cur.athletes_count || 0), 0)
})

const uniqueVenuesCount = computed(() => {
  const set = new Set(bookings.value.map(b => b.venue_name).filter(Boolean))
  return set.size
})

async function fetchData() {
  loading.value = true
  try {
    const [bRes, hRes] = await Promise.allSettled([
      apiGet('/facilities/bookings'),
      apiGet('/facilities/hostels')
    ])

    if (bRes.status === 'fulfilled') {
      bookings.value = bRes.value.data?.bookings || bRes.value.bookings || []
    }
    if (hRes.status === 'fulfilled') {
      hostels.value = hRes.value.data?.hostels || hRes.value.hostels || []
    }
  } finally {
    loading.value = false
  }
}

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric' })
}

onMounted(() => {
  fetchData()
})
</script>
