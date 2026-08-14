<template>
  <LayoutDefault title="Fleet & Transport Logistics">
    <div class="space-y-6 pb-12">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-car-alt-1 text-xs"></i> Vehicle Fleet & Movement
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Fleet Management & Transport Operations</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Official vehicle asset register, national team delegations transport, trip dispatch requisition orders, and vehicle service logs.
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
              <h5>Active Vehicles</h5>
              <h2>{{ vehicles.length }}</h2>
              <span class="badge badge-primary"><i class="icofont-car"></i> Official Fleet</span>
            </div>
            <div class="banner-img">
              <i class="icofont-car-alt-4"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Trips Dispatched</h5>
              <h2>{{ trips.length }}</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Active Missions</span>
            </div>
            <div class="banner-img">
              <i class="icofont-paper-plane"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Vehicle Service Health</h5>
              <h2>{{ vehicles.length ? '100%' : '0%' }}</h2>
              <span class="badge badge-success"><i class="icofont-shield-check"></i> Inspected</span>
            </div>
            <div class="banner-img">
              <i class="icofont-wrench"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Fleet Asset Value</h5>
              <h2>UGX {{ Number(totalFleetValue).toLocaleString() }}</h2>
              <span class="badge badge-info"><i class="icofont-money"></i> Capital Assets</span>
            </div>
            <div class="banner-img">
              <i class="icofont-money"></i>
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

      <!-- TAB 1: VEHICLE ASSET REGISTER -->
      <div v-if="activeTab === 'vehicles'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Active Government Vehicle Asset Register</h4>
          </div>

          <div class="divide-y divide-slate-100 dark:divide-slate-800">
            <div v-for="v in vehicles" :key="v.id" class="p-5 hover:bg-slate-50/70 dark:hover:bg-slate-800/40 transition">
              <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                <div class="space-y-1.5 flex-1">
                  <div class="flex items-center gap-2 flex-wrap">
                    <span class="font-mono text-xs font-bold text-blue-600 dark:text-blue-400 bg-slate-100 dark:bg-slate-800 px-2 py-0.5 rounded">{{ v.registration_number }}</span>
                    <span class="badge badge-info">{{ v.vehicle_type }}</span>
                    <span class="badge badge-success">{{ v.operational_status }}</span>
                  </div>
                  <h4 class="text-sm font-bold text-slate-900 dark:text-white">{{ v.make_model }} ({{ v.model_year }})</h4>
                  <p class="text-xs text-slate-500">Assigned Driver / Dept: <strong class="text-slate-700 dark:text-slate-300">{{ v.assigned_department_or_driver }}</strong></p>
                  <div class="flex items-center gap-4 text-xs text-slate-500 pt-1">
                    <span><i class="icofont-dashboard"></i> Odometer: <strong class="text-slate-700 dark:text-slate-300">{{ Number(v.current_odometer_km).toLocaleString() }} KM</strong></span>
                    <span><i class="icofont-wrench"></i> Next Service: <strong>{{ Number(v.next_service_km).toLocaleString() }} KM</strong></span>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="!vehicles.length" class="text-center py-8 text-xs text-slate-400">
              No vehicles registered in the fleet database.
            </div>
          </div>
        </div>
      </div>

      <!-- TAB 2: TRIP REQUISITIONS -->
      <div v-if="activeTab === 'trips'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Trip Dispatch & Delegations Transport Requisitions</h4>
          </div>
          <div class="card-body p-5 space-y-3">
            <div v-for="t in trips" :key="t.id" class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30">
              <div class="flex items-center justify-between">
                <span class="font-mono text-xs font-bold text-blue-600 dark:text-blue-400">{{ t.trip_number }}</span>
                <span class="badge badge-success">{{ t.trip_status }}</span>
              </div>
              <h5 class="text-sm font-bold text-slate-900 dark:text-white mt-1">{{ t.destination }}</h5>
              <p class="text-xs text-slate-500 mt-0.5">Mission: {{ t.mission_purpose }} ({{ t.passengers_count }} Passengers)</p>
              <div class="text-xs text-slate-500 pt-2 mt-2 border-t border-slate-200 dark:border-slate-700/50">
                Departure: <strong>{{ formatDate(t.departure_date) }}</strong> · Driver: <strong>{{ t.assigned_driver_name }}</strong>
              </div>
            </div>
            <div v-if="!trips.length" class="text-center py-6 text-xs text-slate-400">
              No transport trips currently dispatched.
            </div>
          </div>
        </div>
      </div>

    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { apiGet } from '@/api/client'

const loading = ref(false)
const vehicles = ref([])
const trips = ref([])
const activeTab = ref('vehicles')

const totalFleetValue = computed(() => {
  return vehicles.value.reduce((acc, cur) => acc + (cur.purchase_price_ugx || cur.asset_value_ugx || 0), 0)
})

const tabs = [
  { id: 'vehicles', label: 'Government Vehicle Asset Register', icon: 'icofont-car' },
  { id: 'trips', label: 'Trip Movement & Team Dispatches', icon: 'icofont-paper-plane' }
]

async function fetchData() {
  loading.value = true
  try {
    const [vRes, tRes] = await Promise.allSettled([
      apiGet('/fleet/vehicles'),
      apiGet('/fleet/trips')
    ])

    if (vRes.status === 'fulfilled') {
      vehicles.value = vRes.value.data?.vehicles || vRes.value.vehicles || []
    }
    if (tRes.status === 'fulfilled') {
      trips.value = tRes.value.data?.trips || tRes.value.trips || []
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
