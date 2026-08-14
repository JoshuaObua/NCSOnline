<template>
  <LayoutDefault title="Online Appraisal">
    <div class="space-y-6 pb-12">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-law-order text-xs"></i> Online Appraisal
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Online Appraisal & Institutional Performance</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Statutory Accounting Officer institutional appraisal, fixed asset ledger sign-offs, and departmental performance evaluations.
          </p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <button @click="fetchData" class="btn btn-sm btn-light flex items-center gap-1.5" :disabled="loading">
            <i class="icofont-refresh" :class="{ 'animate-spin': loading }"></i> Refresh
          </button>
          <button @click="printReport" class="btn btn-sm btn-primary flex items-center gap-1.5">
            <i class="icofont-printer"></i> Print Appraisal Summary
          </button>
        </div>
      </div>

      <!-- Macro Metric KPI Cards (Clean 0 Default Digits) -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Fixed Asset Portfolio</h5>
              <h2>UGX {{ Number(summaryStats.net_book_value_ugx || 0).toLocaleString() }}</h2>
              <span class="badge badge-primary"><i class="icofont-check-circled"></i> {{ summaryStats.total_asset_items_count || assetList.length }} Ledger Assets</span>
            </div>
            <div class="banner-img">
              <i class="icofont-building-alt"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Staff Performance Score</h5>
              <h2>{{ summaryStats.staff_appraisal_avg_score || 0 }}%</h2>
              <span class="badge badge-success"><i class="icofont-badge"></i> {{ summaryStats.total_staff_evaluated || 0 }} Evaluated</span>
            </div>
            <div class="banner-img">
              <i class="icofont-users-alt-4"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Vote Burn Rate</h5>
              <h2>{{ summaryStats.budget_execution_rate || 0 }}%</h2>
              <span class="badge badge-info"><i class="icofont-chart-growth"></i> Fiscal Execution</span>
            </div>
            <div class="banner-img">
              <i class="icofont-chart-growth"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Grant Clearance</h5>
              <h2>{{ summaryStats.federation_grant_clearance || 0 }}%</h2>
              <span class="badge badge-success"><i class="icofont-shield-check"></i> {{ summaryStats.audit_risk_rating || 'COMPLIANT' }}</span>
            </div>
            <div class="banner-img">
              <i class="icofont-certificate-alt-1"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Departmental Reports Hub Launcher -->
      <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
        <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5 flex items-center justify-between">
          <h4 class="text-sm font-bold text-slate-900 dark:text-white">Departmental Compiled Reports Directory</h4>
          <span class="text-xs text-slate-500 font-medium">Independent Dossier Pages</span>
        </div>
        <div class="p-5 grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
          <router-link
            to="/executive/reports/engineering"
            class="p-3.5 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30 hover:border-blue-500 transition group flex flex-col justify-between"
          >
            <i class="icofont-wrench text-xl text-blue-600 dark:text-blue-400 group-hover:scale-110 transition-transform mb-2"></i>
            <div>
              <span class="text-xs font-bold text-slate-900 dark:text-white block group-hover:text-blue-600">Engineering</span>
              <span class="text-[10px] text-slate-400">Civil, Power, Water</span>
            </div>
          </router-link>

          <router-link
            to="/executive/reports/human-resources"
            class="p-3.5 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30 hover:border-blue-500 transition group flex flex-col justify-between"
          >
            <i class="icofont-users text-xl text-blue-600 dark:text-blue-400 group-hover:scale-110 transition-transform mb-2"></i>
            <div>
              <span class="text-xs font-bold text-slate-900 dark:text-white block group-hover:text-blue-600">Human Resources</span>
              <span class="text-[10px] text-slate-400">Leave, Scorecards</span>
            </div>
          </router-link>

          <router-link
            to="/executive/reports/finance-accounts"
            class="p-3.5 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30 hover:border-blue-500 transition group flex flex-col justify-between"
          >
            <i class="icofont-money text-xl text-blue-600 dark:text-blue-400 group-hover:scale-110 transition-transform mb-2"></i>
            <div>
              <span class="text-xs font-bold text-slate-900 dark:text-white block group-hover:text-blue-600">Finance & Accounts</span>
              <span class="text-[10px] text-slate-400">Budget, NTR, Grants</span>
            </div>
          </router-link>

          <router-link
            to="/executive/reports/technical-sports"
            class="p-3.5 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30 hover:border-blue-500 transition group flex flex-col justify-between"
          >
            <i class="icofont-badge text-xl text-blue-600 dark:text-blue-400 group-hover:scale-110 transition-transform mb-2"></i>
            <div>
              <span class="text-xs font-bold text-slate-900 dark:text-white block group-hover:text-blue-600">Technical & Sports</span>
              <span class="text-[10px] text-slate-400">51 Federations</span>
            </div>
          </router-link>

          <router-link
            to="/executive/reports/medical-science"
            class="p-3.5 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30 hover:border-blue-500 transition group flex flex-col justify-between"
          >
            <i class="icofont-heart-beat text-xl text-blue-600 dark:text-blue-400 group-hover:scale-110 transition-transform mb-2"></i>
            <div>
              <span class="text-xs font-bold text-slate-900 dark:text-white block group-hover:text-blue-600">Sports Medicine</span>
              <span class="text-[10px] text-slate-400">Screening, WADA</span>
            </div>
          </router-link>

          <router-link
            to="/executive/reports/legal-logistics"
            class="p-3.5 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30 hover:border-blue-500 transition group flex flex-col justify-between"
          >
            <i class="icofont-law-document text-xl text-blue-600 dark:text-blue-400 group-hover:scale-110 transition-transform mb-2"></i>
            <div>
              <span class="text-xs font-bold text-slate-900 dark:text-white block group-hover:text-blue-600">Legal & Logistics</span>
              <span class="text-[10px] text-slate-400">Contracts, Tribunal</span>
            </div>
          </router-link>
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

      <!-- TAB 1: FIXED ASSET LEDGER APPRAISAL -->
      <div v-if="activeTab === 'assets'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 flex flex-col sm:flex-row sm:items-center justify-between gap-3 py-3.5 px-5">
            <div>
              <h4 class="text-sm font-bold text-slate-900 dark:text-white">Fixed Asset Valuation & Revaluation Approval Desk</h4>
              <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Statutory appraisal of asset values entered by the Accounting Department</p>
            </div>
            <div>
              <button class="btn btn-sm btn-primary flex items-center gap-1">
                <i class="icofont-stamp"></i> Statutory Portfolio Stamp
              </button>
            </div>
          </div>

          <div class="card-body p-0 overflow-x-auto">
            <table class="table table-striped table-hover mb-0 text-left text-xs">
              <thead class="bg-slate-50 dark:bg-slate-800/60 text-slate-600 dark:text-slate-400 uppercase tracking-wider font-semibold">
                <tr>
                  <th class="py-3 px-4">Asset Tag / Reference</th>
                  <th class="py-3 px-4">Class / Category</th>
                  <th class="py-3 px-4">Historical Cost</th>
                  <th class="py-3 px-4">Accounting NBV</th>
                  <th class="py-3 px-4">Physical Condition</th>
                  <th class="py-3 px-4">GS Statutory Action</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
                <tr v-for="asset in assetList" :key="asset.id">
                  <td class="py-3 px-4">
                    <div class="font-bold text-slate-900 dark:text-white">{{ asset.asset_tag }}</div>
                    <div class="text-[11px] text-slate-500">{{ asset.asset_description }}</div>
                  </td>
                  <td class="py-3 px-4 font-mono font-medium text-slate-600 dark:text-slate-400">{{ asset.category_segment }}</td>
                  <td class="py-3 px-4 font-mono text-slate-700 dark:text-slate-300">UGX {{ Number(asset.recorded_historical_cost || 0).toLocaleString() }}</td>
                  <td class="py-3 px-4 font-mono font-bold text-blue-600 dark:text-blue-400">UGX {{ Number(asset.approved_net_book_value || 0).toLocaleString() }}</td>
                  <td class="py-3 px-4">
                    <span class="badge badge-success">{{ asset.physical_condition_grade || 'GOOD' }}</span>
                  </td>
                  <td class="py-3 px-4">
                    <span class="badge badge-primary flex items-center gap-1 w-max">
                      <i class="icofont-check"></i> {{ asset.gs_action || 'CONFIRMED' }}
                    </span>
                  </td>
                </tr>
                <tr v-if="!assetList.length">
                  <td colspan="6" class="text-center py-6 text-xs text-slate-400">
                    No asset valuations currently queued for appraisal.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- TAB 2: FINANCIAL EFFICIENCY -->
      <div v-if="activeTab === 'financial'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Departmental Financial Performance & Vote-Head Efficiency</h4>
          </div>
          <div class="card-body p-5">
            <div v-if="financialList.length" class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div v-for="fin in financialList" :key="fin.id" class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30 space-y-2">
                <div class="flex items-center justify-between">
                  <span class="font-bold text-sm text-slate-900 dark:text-white">{{ fin.department_name }}</span>
                  <span class="badge badge-success">Grade {{ fin.financial_grade || 'A' }}</span>
                </div>
                <div class="grid grid-cols-2 gap-2 text-xs pt-2 border-t border-slate-200 dark:border-slate-700/50">
                  <div>Budget: <strong class="text-slate-800 dark:text-slate-200 font-mono">UGX {{ Number(fin.allocated_budget_ugx || 0).toLocaleString() }}</strong></div>
                  <div>Spent: <strong class="text-slate-800 dark:text-slate-200 font-mono">UGX {{ Number(fin.actual_expenditure_ugx || 0).toLocaleString() }}</strong></div>
                  <div>Execution Rate: <strong class="text-blue-600 dark:text-blue-400 font-bold">{{ fin.budget_execution_rate || 0 }}%</strong></div>
                  <div>Audit Queries: <strong class="text-emerald-600 dark:text-emerald-400 font-bold">{{ fin.audit_query_count || 0 }} Active</strong></div>
                </div>
              </div>
            </div>
            <div v-else class="text-center py-6 text-xs text-slate-400">
              No financial performance appraisals ingested for this quarter.
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
const activeTab = ref('assets')
const summaryStats = ref({})
const assetList = ref([])
const financialList = ref([])

const tabs = [
  { id: 'assets', label: 'Fixed Asset Ledger Appraisal', icon: 'icofont-building-alt' },
  { id: 'financial', label: 'Vote-Head Financial Execution', icon: 'icofont-chart-growth' }
]

function printReport() {
  window.print()
}

async function fetchData() {
  loading.value = true
  try {
    const [summaryRes, assetsRes, finRes] = await Promise.allSettled([
      apiGet('/executive/appraisal/summary'),
      apiGet('/executive/appraisal/assets/ledger'),
      apiGet('/executive/appraisal/performance/financial')
    ])

    if (summaryRes.status === 'fulfilled') {
      summaryStats.value = summaryRes.value.data || summaryRes.value || {}
    }
    if (assetsRes.status === 'fulfilled') {
      assetList.value = assetsRes.value.data?.asset_valuations || assetsRes.value.asset_valuations || []
    }
    if (finRes.status === 'fulfilled') {
      financialList.value = finRes.value.data?.financial_appraisals || finRes.value.financial_appraisals || []
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
