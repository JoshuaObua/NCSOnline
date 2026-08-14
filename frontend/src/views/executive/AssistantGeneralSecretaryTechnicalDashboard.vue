<template>
  <LayoutDefault title="AGS - Technical Directorate">
    <div class="space-y-6 pb-12">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-badge text-xs"></i> Technical Command
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Assistant General Secretary (Technical)</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Executive oversight of 51 National Sports Federations, Technical Grants, National Teams & Delegations, Engineering Infrastructure, and Facility Readiness Certifications.
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
              <h5>Recognised Federations</h5>
              <h2>{{ stats.active_federations_count || 0 }}</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Statutory Compliance</span>
            </div>
            <div class="banner-img">
              <i class="icofont-trophy"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Grant Execution</h5>
              <h2>UGX {{ Number(stats.quarterly_grants_ugx || 0).toLocaleString() }}</h2>
              <span class="badge badge-primary"><i class="icofont-coins"></i> High Impact Priority</span>
            </div>
            <div class="banner-img">
              <i class="icofont-money-bag"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Pending Tech Vetting</h5>
              <h2>{{ stats.pending_approvals_count || approvals.length }}</h2>
              <span class="badge badge-warning"><i class="icofont-clock-time"></i> Action Required</span>
            </div>
            <div class="banner-img">
              <i class="icofont-sand-clock"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Certified Facilities</h5>
              <h2>{{ stats.certified_facilities_count || readinessList.length }}</h2>
              <span class="badge badge-info"><i class="icofont-building"></i> Standard Certified</span>
            </div>
            <div class="banner-img">
              <i class="icofont-building"></i>
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

      <!-- TAB 1: TECHNICAL APPROVALS & VETTING -->
      <div v-if="activeTab === 'approvals'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 flex flex-col sm:flex-row sm:items-center justify-between gap-3 py-3.5 px-5">
            <div>
              <h4 class="text-sm font-bold text-slate-900 dark:text-white">AGS-T Technical Vetting Pipeline</h4>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Endorsement queue for Grants, CapEx engineering, and National Teams</p>
            </div>
            <div>
              <select v-model="filterCategory" class="form-control form-control-sm text-xs rounded-lg border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300">
                <option value="">All Categories</option>
                <option value="FEDERATION_GRANT">Federation Grants</option>
                <option value="INFRASTRUCTURE_CAPEX">Infrastructure CapEx</option>
                <option value="DELEGATION_CLEARANCE">Delegation Travel</option>
                <option value="VENUE_CERTIFICATION">Venue Certification</option>
              </select>
            </div>
          </div>

          <div class="divide-y divide-slate-100 dark:divide-slate-800">
            <div v-for="item in filteredApprovals" :key="item.id" class="p-5 hover:bg-slate-50/70 dark:hover:bg-slate-800/40 transition">
              <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                <div class="space-y-2 flex-1">
                  <div class="flex items-center gap-2 flex-wrap">
                    <span class="font-mono text-xs font-bold text-slate-600 dark:text-slate-400 bg-slate-100 dark:bg-slate-800 px-2 py-0.5 rounded">{{ item.reference_no }}</span>
                    <span class="badge badge-info">{{ formatCategory(item.category) }}</span>
                    <span class="badge badge-warning">{{ formatStatus(item.agst_status) }}</span>
                    <span v-if="item.escalated_to_gs" class="badge badge-primary flex items-center gap-1">
                      <i class="icofont-arrow-up"></i> Endorsed to GS
                    </span>
                  </div>
                  <h4 class="text-sm font-bold text-slate-900 dark:text-white">{{ item.title }}</h4>
                  <p class="text-xs text-slate-600 dark:text-slate-400">{{ item.description }}</p>
                  <div class="flex items-center gap-4 text-xs text-slate-500 pt-1">
                    <span><i class="icofont-user"></i> Submitter: <strong class="text-slate-700 dark:text-slate-300">{{ item.submitting_officer_name }}</strong> ({{ item.originating_department }})</span>
                    <span v-if="item.financial_implication_ugx > 0"><i class="icofont-money"></i> Implication: <strong class="text-blue-600 dark:text-blue-400 font-mono">UGX {{ Number(item.financial_implication_ugx).toLocaleString() }}</strong></span>
                  </div>
                </div>

                <div class="flex items-center gap-2 flex-shrink-0">
                  <template v-if="item.agst_status === 'PENDING_REVIEW'">
                    <button @click="actionItem(item, 'ENDORSE_TO_GS')" class="btn btn-sm btn-primary flex items-center gap-1">
                      <i class="icofont-check-circled"></i> Endorse to GS
                    </button>
                    <button @click="actionItem(item, 'APPROVE')" class="btn btn-sm btn-success flex items-center gap-1">
                      <i class="icofont-stamp"></i> Approve Tech
                    </button>
                    <button @click="actionItem(item, 'RETURN')" class="btn btn-sm btn-light">
                      Return
                    </button>
                  </template>
                  <template v-else>
                    <span class="badge badge-light">Actioned</span>
                  </template>
                </div>
              </div>
            </div>
            <div v-if="!filteredApprovals.length" class="text-center py-8 text-xs text-slate-400">
              No technical approvals currently pending review.
            </div>
          </div>
        </div>
      </div>

      <!-- TAB 2: FACILITY READINESS -->
      <div v-if="activeTab === 'facilities'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Facility Structural Integrity & Readiness Certifications</h4>
          </div>
          <div class="card-body p-5">
            <div v-if="readinessList.length" class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div v-for="fac in readinessList" :key="fac.id" class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30 space-y-2">
                <div class="flex items-center justify-between">
                  <span class="font-bold text-sm text-slate-900 dark:text-white">{{ fac.facility_name }}</span>
                  <span class="badge badge-success">{{ fac.certification_status }}</span>
                </div>
                <h5 class="text-xs font-semibold text-slate-700 dark:text-slate-300">{{ fac.event_name }}</h5>
                <div class="grid grid-cols-2 gap-2 text-xs pt-2 border-t border-slate-200 dark:border-slate-700/50">
                  <div>Engineer: <strong class="text-slate-800 dark:text-slate-200">{{ fac.inspecting_engineer_name }}</strong></div>
                  <div>Rating: <strong class="text-blue-600 dark:text-blue-400">{{ fac.readiness_rating }}%</strong></div>
                  <div>Lighting: <strong>{{ fac.lighting_lux_level }} Lux</strong></div>
                  <div>Turf/Court Score: <strong>{{ fac.turf_court_score }}/10</strong></div>
                </div>
              </div>
            </div>
            <div v-else class="text-center py-6 text-xs text-slate-400">
              No facility readiness certifications currently recorded.
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
import { apiGet, apiPost } from '@/api/client'

const loading = ref(false)
const stats = ref({})
const approvals = ref([])
const readinessList = ref([])
const activeTab = ref('approvals')
const filterCategory = ref('')

const tabs = [
  { id: 'approvals', label: 'Technical Approvals & Grants', icon: 'icofont-inbox' },
  { id: 'facilities', label: 'Facility Readiness Certifications', icon: 'icofont-building' }
]

const filteredApprovals = computed(() => {
  if (!filterCategory.value) return approvals.value
  return approvals.value.filter(a => a.category === filterCategory.value)
})

async function fetchData() {
  loading.value = true
  try {
    const [statsRes, approvalsRes, facRes] = await Promise.allSettled([
      apiGet('/executive/ags-t/dashboard'),
      apiGet('/executive/ags-t/approvals'),
      apiGet('/executive/ags-t/facilities/readiness')
    ])

    if (statsRes.status === 'fulfilled') {
      stats.value = statsRes.value.data || statsRes.value || {}
    }
    if (approvalsRes.status === 'fulfilled') {
      approvals.value = approvalsRes.value.data?.approvals || approvalsRes.value.approvals || []
    }
    if (facRes.status === 'fulfilled') {
      readinessList.value = facRes.value.data?.certifications || facRes.value.certifications || []
    }
  } finally {
    loading.value = false
  }
}

async function actionItem(item, actionType) {
  try {
    await apiPost(`/executive/ags-t/approvals/${item.id}/action`, { action: actionType, remarks: 'Actioned by AGS-Technical' })
    await fetchData()
  } catch {
    // Fail gracefully
  }
}

function formatCategory(cat) {
  return cat ? cat.replace(/_/g, ' ') : ''
}

function formatStatus(status) {
  return status ? status.replace(/_/g, ' ') : ''
}

onMounted(() => {
  fetchData()
})
</script>
