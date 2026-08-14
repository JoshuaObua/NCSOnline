<template>
  <LayoutDefault title="AGS - Administration Directorate">
    <div class="space-y-6 pb-12">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-architecture-alt text-xs"></i> Administration Command
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Assistant General Secretary (Administration)</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Executive coordination of Human Resources, Finance & Accounts, PDU Procurement Vetting, PR Communications, and Digital Infrastructure.
          </p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <button @click="fetchData" class="btn btn-sm btn-light flex items-center gap-1.5" :disabled="loading">
            <i class="icofont-refresh" :class="{ 'animate-spin': loading }"></i> Refresh
          </button>
          <button @click="activeTab = 'directives'" class="btn btn-sm btn-primary flex items-center gap-1.5">
            <i class="icofont-paper"></i> Issue Directive
          </button>
        </div>
      </div>

      <!-- Macro Metric KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Staff Establishment</h5>
              <h2>{{ stats.active_staff_headcount || 0 }}</h2>
              <span class="badge badge-primary"><i class="icofont-check-circled"></i> Active Staff</span>
            </div>
            <div class="banner-img">
              <i class="icofont-users-alt-4"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Budget Spend Velocity</h5>
              <h2>{{ stats.budget_execution_rate || 0 }}%</h2>
              <span class="badge badge-info"><i class="icofont-chart-growth"></i> Subvention Spend</span>
            </div>
            <div class="banner-img">
              <i class="icofont-chart-growth"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>NTR Inflows</h5>
              <h2>UGX {{ Number(stats.ntr_inflows_ugx || 0).toLocaleString() }}</h2>
              <span class="badge badge-success"><i class="icofont-bank"></i> Collections</span>
            </div>
            <div class="banner-img">
              <i class="icofont-bank"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Pending Admin Vetting</h5>
              <h2>{{ stats.pending_vetting_count || approvals.length }}</h2>
              <span class="badge badge-warning"><i class="icofont-shield"></i> Queued for Action</span>
            </div>
            <div class="banner-img">
              <i class="icofont-inbox"></i>
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

      <!-- TAB 1: ADMINISTRATIVE VETTING & CLEARANCE -->
      <div v-if="activeTab === 'approvals'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 flex flex-col sm:flex-row sm:items-center justify-between gap-3 py-3.5 px-5">
            <div>
              <h4 class="text-sm font-bold text-slate-900 dark:text-white">AGS-A Administrative Vetting Pipeline</h4>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">First-line administrative endorsement on Form 5s, Payroll, Leave Rosters, and Communications</p>
            </div>
            <div>
              <select v-model="filterCategory" class="form-control form-control-sm text-xs rounded-lg border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300">
                <option value="">All Categories</option>
                <option value="FORM_5_VETTING">Form 5 Procurement</option>
                <option value="PAYROLL_CLEARANCE">Monthly Payroll</option>
                <option value="LEAVE_VETTING">Leave Vetting</option>
                <option value="PRESS_RELEASE">Press Releases</option>
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
                    <span class="badge badge-warning">{{ formatStatus(item.agsa_status) }}</span>
                    <span v-if="item.escalated_to_gs" class="badge badge-primary flex items-center gap-1">
                      <i class="icofont-arrow-up"></i> Endorsed to GS
                    </span>
                  </div>
                  <h4 class="text-sm font-bold text-slate-900 dark:text-white">{{ item.title }}</h4>
                  <p class="text-xs text-slate-600 dark:text-slate-400">{{ item.summary }}</p>
                  <div class="flex items-center gap-4 text-xs text-slate-500 pt-1">
                    <span><i class="icofont-user"></i> Submitter: <strong class="text-slate-700 dark:text-slate-300">{{ item.submitting_officer_name }}</strong> ({{ item.originating_department }})</span>
                    <span v-if="item.financial_value_ugx > 0"><i class="icofont-money"></i> Value: <strong class="text-blue-600 dark:text-blue-400 font-mono">UGX {{ Number(item.financial_value_ugx).toLocaleString() }}</strong></span>
                  </div>
                </div>

                <div class="flex items-center gap-2 flex-shrink-0">
                  <template v-if="item.agsa_status === 'PENDING_VETTING'">
                    <button @click="actionItem(item, 'ENDORSE_TO_GS')" class="btn btn-sm btn-primary flex items-center gap-1">
                      <i class="icofont-check-circled"></i> Endorse to GS
                    </button>
                    <button @click="actionItem(item, 'APPROVE')" class="btn btn-sm btn-success flex items-center gap-1">
                      <i class="icofont-stamp"></i> Approve Admin
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
              No administrative items currently pending review.
            </div>
          </div>
        </div>
      </div>

      <!-- TAB 2: HR ESTABLISHMENT & PAYROLL -->
      <div v-if="activeTab === 'hr'" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="card p-5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm mb-0">
            <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider">Monthly Payroll Release</h4>
            <p class="text-2xl font-black text-slate-900 dark:text-white mt-1 font-mono">UGX {{ Number(stats.monthly_payroll_ugx || 0).toLocaleString() }}</p>
            <p class="text-xs text-emerald-600 dark:text-emerald-400 mt-1 flex items-center gap-1 font-medium">
              <i class="icofont-check-circled"></i> Cleared for Bank EFT
            </p>
          </div>
          <div class="card p-5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm mb-0">
            <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider">Staff on Active Leave</h4>
            <p class="text-2xl font-black text-slate-900 dark:text-white mt-1">{{ stats.staff_on_leave || 0 }}</p>
            <p class="text-xs text-blue-600 dark:text-blue-400 mt-1 flex items-center gap-1 font-medium">
              <i class="icofont-refresh"></i> Duty Handover Verified
            </p>
          </div>
          <div class="card p-5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 shadow-sm mb-0">
            <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider">Appraisal Submission Rate</h4>
            <p class="text-2xl font-black text-slate-900 dark:text-white mt-1">{{ stats.appraisal_rate || 0 }}%</p>
            <p class="text-xs text-blue-600 dark:text-blue-400 mt-1 flex items-center gap-1 font-medium">
              <i class="icofont-inbox"></i> Returns Ingested
            </p>
          </div>
        </div>
      </div>

      <!-- TAB 3: DIRECTIVES -->
      <div v-if="activeTab === 'directives'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Issued Administrative Directives</h4>
          </div>
          <div class="card-body p-5 space-y-3">
            <div v-for="d in directives" :key="d.id" class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30">
              <div class="flex items-center justify-between">
                <span class="font-mono text-xs font-bold text-blue-600 dark:text-blue-400">{{ d.directive_no }}</span>
                <span class="badge" :class="d.priority === 'STATUTORY' ? 'badge-danger' : 'badge-warning'">{{ d.priority }}</span>
              </div>
              <h5 class="text-sm font-bold text-slate-900 dark:text-white mt-1">{{ d.subject }}</h5>
              <p class="text-xs text-slate-600 dark:text-slate-400 mt-1">{{ d.content }}</p>
            </div>
            <div v-if="!directives.length" class="text-center py-6 text-xs text-slate-400">
              No active administrative directives issued.
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
const directives = ref([])
const activeTab = ref('approvals')
const filterCategory = ref('')

const tabs = [
  { id: 'approvals', label: 'Admin Approvals & Form 5s', icon: 'icofont-inbox' },
  { id: 'hr', label: 'HR & Payroll Radar', icon: 'icofont-users' },
  { id: 'directives', label: 'Executive Directives', icon: 'icofont-paper' }
]

const filteredApprovals = computed(() => {
  if (!filterCategory.value) return approvals.value
  return approvals.value.filter(a => a.category === filterCategory.value)
})

async function fetchData() {
  loading.value = true
  try {
    const [statsRes, approvalsRes, dirRes] = await Promise.allSettled([
      apiGet('/executive/ags-a/dashboard'),
      apiGet('/executive/ags-a/approvals'),
      apiGet('/executive/ags-a/directives')
    ])

    if (statsRes.status === 'fulfilled') {
      stats.value = statsRes.value.data || statsRes.value || {}
    }
    if (approvalsRes.status === 'fulfilled') {
      approvals.value = approvalsRes.value.data?.approvals || approvalsRes.value.approvals || []
    }
    if (dirRes.status === 'fulfilled') {
      directives.value = dirRes.value.data?.directives || dirRes.value.directives || []
    }
  } finally {
    loading.value = false
  }
}

async function actionItem(item, actionType) {
  try {
    await apiPost(`/executive/ags-a/approvals/${item.id}/action`, { action: actionType, comments: 'Actioned by AGS-Administration' })
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
