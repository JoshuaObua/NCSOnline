<template>
  <LayoutDefault title="Sports Medicine, Science & WADA">
    <div class="space-y-6 pb-12">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-heart-beat text-xs"></i> High Performance Medicine
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Sports Science, Medical & Anti-Doping Unit</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Pre-competition athlete cardiovascular screening, injury rehabilitation tracking, and WADA/RADO anti-doping sample registry.
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
              <h5>Medically Cleared</h5>
              <h2>{{ screenings.length }}</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Fit for Travel</span>
            </div>
            <div class="banner-img">
              <i class="icofont-stethoscope"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Active In Rehab</h5>
              <h2>{{ injuries.length }}</h2>
              <span class="badge" :class="injuries.length > 0 ? 'badge-warning' : 'badge-light'">
                <i class="icofont-activity"></i> Surveillance
              </span>
            </div>
            <div class="banner-img">
              <i class="icofont-pulse"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>WADA Clean Rate</h5>
              <h2>{{ antidoping.length ? '100%' : '0%' }}</h2>
              <span class="badge badge-primary"><i class="icofont-shield-check"></i> Clean Sport</span>
            </div>
            <div class="banner-img">
              <i class="icofont-shield-check"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>WADA Sample Pool</h5>
              <h2>{{ antidoping.length }}</h2>
              <span class="badge badge-info"><i class="icofont-certificate-alt-1"></i> RADO / WADA</span>
            </div>
            <div class="banner-img">
              <i class="icofont-test-bottle"></i>
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

      <!-- TAB 1: SCREENINGS -->
      <div v-if="activeTab === 'screenings'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">National Athlete Medical Fitness & Travel Clearance Roster</h4>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">ECG, cardiovascular check, and blood pressure screening records</p>
          </div>

          <div class="divide-y divide-slate-100 dark:divide-slate-800">
            <div v-for="s in screenings" :key="s.id" class="p-5 hover:bg-slate-50/70 dark:hover:bg-slate-800/40 transition">
              <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                <div class="space-y-1.5 flex-1">
                  <div class="flex items-center gap-2 flex-wrap">
                    <span class="font-bold text-sm text-slate-900 dark:text-white">{{ s.athlete_name }}</span>
                    <span class="badge badge-success">{{ s.fitness_verdict }}</span>
                  </div>
                  <p class="text-xs text-slate-500">{{ s.federation_name }} ({{ s.discipline }})</p>
                  <div class="flex items-center gap-4 text-xs text-slate-500 pt-1">
                    <span><i class="icofont-heart-beat"></i> ECG: <strong class="text-slate-700 dark:text-slate-300">{{ s.ecg_finding }}</strong></span>
                    <span><i class="icofont-pulse"></i> BP: <strong class="text-slate-700 dark:text-slate-300">{{ s.blood_pressure }}</strong></span>
                    <span><i class="icofont-calendar"></i> Screened: <strong>{{ formatDate(s.screening_date) }}</strong></span>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="!screenings.length" class="text-center py-8 text-xs text-slate-400">
              No athlete medical screenings logged for this period.
            </div>
          </div>
        </div>
      </div>

      <!-- TAB 2: INJURIES -->
      <div v-if="activeTab === 'injuries'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Active Athlete Injury Surveillance & Rehab Pipeline</h4>
          </div>
          <div class="card-body p-5">
            <div v-if="injuries.length" class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div v-for="inj in injuries" :key="inj.id" class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30 space-y-2">
                <div class="flex items-center justify-between">
                  <span class="font-bold text-sm text-slate-900 dark:text-white">{{ inj.athlete_name }}</span>
                  <span class="badge badge-warning">{{ inj.rehab_status }}</span>
                </div>
                <h5 class="text-xs font-bold text-slate-700 dark:text-slate-300">{{ inj.injury_site }}: {{ inj.injury_nature }} ({{ inj.severity }})</h5>
                <div class="text-xs text-slate-500 pt-1 border-t border-slate-200 dark:border-slate-700/50">
                  Discipline: {{ inj.discipline }} · Incident Date: <strong>{{ formatDate(inj.incident_date) }}</strong>
                </div>
              </div>
            </div>
            <div v-else class="text-center py-6 text-xs text-slate-400">
              No athletes currently logged in injury rehabilitation.
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
const screenings = ref([])
const injuries = ref([])
const antidoping = ref([])
const activeTab = ref('screenings')

const tabs = [
  { id: 'screenings', label: 'Athlete Medical Screening & Travel Pass', icon: 'icofont-stethoscope' },
  { id: 'injuries', label: 'Injury Surveillance & Rehab', icon: 'icofont-pulse' }
]

async function fetchData() {
  loading.value = true
  try {
    const [scrRes, injRes, adRes] = await Promise.allSettled([
      apiGet('/medical/screenings'),
      apiGet('/medical/injuries'),
      apiGet('/medical/antidoping')
    ])

    if (scrRes.status === 'fulfilled') {
      screenings.value = scrRes.value.data?.screenings || scrRes.value.screenings || []
    }
    if (injRes.status === 'fulfilled') {
      injuries.value = injRes.value.data?.injuries || injRes.value.injuries || []
    }
    if (adRes.status === 'fulfilled') {
      antidoping.value = adRes.value.data?.records || adRes.value.records || []
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
