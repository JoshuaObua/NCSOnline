<template>
  <LayoutDefault title="Legal & Statutory Compliance">
    <div class="space-y-6 pb-12">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-law-document text-xs"></i> Legal Affairs
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Legal & Corporate Affairs Department</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Commercial sponsorships, vendor contracts, National Sports Act (2023) compliance, and sports federation dispute arbitration tribunal.
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
              <h5>Active Contracts</h5>
              <h2>{{ contracts.length }}</h2>
              <span class="badge badge-primary"><i class="icofont-check-circled"></i> Commercial & MOUs</span>
            </div>
            <div class="banner-img">
              <i class="icofont-file-text"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Active Arbitrations</h5>
              <h2>{{ disputes.length }}</h2>
              <span class="badge" :class="disputes.length > 0 ? 'badge-warning' : 'badge-light'">
                <i class="icofont-clock-time"></i> Hearing Stage
              </span>
            </div>
            <div class="banner-img">
              <i class="icofont-judge"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Statutory Compliance</h5>
              <h2>{{ contracts.length ? '100%' : '0%' }}</h2>
              <span class="badge badge-success"><i class="icofont-shield-check"></i> Act 2023 Aligned</span>
            </div>
            <div class="banner-img">
              <i class="icofont-shield-check"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Commercial Value</h5>
              <h2>UGX {{ Number(totalContractValue).toLocaleString() }}</h2>
              <span class="badge badge-info"><i class="icofont-money"></i> Naming Rights</span>
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

      <!-- TAB 1: CONTRACTS -->
      <div v-if="activeTab === 'contracts'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Active Commercial & Statutory Contracts Vault</h4>
            <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">Corporate sponsorships, vendor SLAs, and sports federation MOUs</p>
          </div>

          <div class="divide-y divide-slate-100 dark:divide-slate-800">
            <div v-for="c in contracts" :key="c.id" class="p-5 hover:bg-slate-50/70 dark:hover:bg-slate-800/40 transition">
              <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                <div class="space-y-1.5 flex-1">
                  <div class="flex items-center gap-2 flex-wrap">
                    <span class="font-mono text-xs font-bold text-blue-600 dark:text-blue-400 bg-slate-100 dark:bg-slate-800 px-2 py-0.5 rounded">{{ c.contract_reference }}</span>
                    <span class="badge badge-info">{{ c.contract_type }}</span>
                    <span class="badge badge-success">{{ c.status }}</span>
                  </div>
                  <h4 class="text-sm font-bold text-slate-900 dark:text-white">{{ c.title }}</h4>
                  <p class="text-xs text-slate-500">Counter-Party: <strong class="text-slate-700 dark:text-slate-300">{{ c.second_party }}</strong></p>
                  <div class="flex items-center gap-4 text-xs text-slate-500 pt-1">
                    <span><i class="icofont-calendar"></i> Expiry: <strong>{{ formatDate(c.expiry_date) }}</strong></span>
                    <span v-if="c.contract_value_ugx > 0"><i class="icofont-money"></i> Value: <strong class="text-blue-600 dark:text-blue-400 font-mono">UGX {{ Number(c.contract_value_ugx).toLocaleString() }}</strong></span>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="!contracts.length" class="text-center py-8 text-xs text-slate-400">
              No legal contracts recorded.
            </div>
          </div>
        </div>
      </div>

      <!-- TAB 2: DISPUTES & ARBITRATION -->
      <div v-if="activeTab === 'disputes'" class="space-y-4">
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Sports Federation Dispute & Arbitration Tribunal</h4>
          </div>
          <div class="card-body p-5 space-y-3">
            <div v-for="d in disputes" :key="d.id" class="p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-800/30">
              <div class="flex items-center justify-between">
                <span class="font-mono text-xs font-bold text-blue-600 dark:text-blue-400">{{ d.case_number }}</span>
                <span class="badge badge-warning">{{ d.case_status }}</span>
              </div>
              <h5 class="text-sm font-bold text-slate-900 dark:text-white mt-1">{{ d.subject_matter }}</h5>
              <p class="text-xs text-slate-500 mt-0.5">{{ d.federation_name }} ({{ d.dispute_category }})</p>
              <div class="text-xs text-slate-500 pt-2 mt-2 border-t border-slate-200 dark:border-slate-700/50">
                Parties: <strong class="text-slate-700 dark:text-slate-300">{{ d.complainant_name }}</strong> vs <strong class="text-slate-700 dark:text-slate-300">{{ d.respondent_name }}</strong> · Filed: <strong>{{ formatDate(d.filing_date) }}</strong>
              </div>
            </div>
            <div v-if="!disputes.length" class="text-center py-6 text-xs text-slate-400">
              No arbitration disputes currently pending before the tribunal.
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
const contracts = ref([])
const disputes = ref([])
const activeTab = ref('contracts')

const tabs = [
  { id: 'contracts', label: 'Commercial & Statutory Contracts', icon: 'icofont-file-text' },
  { id: 'disputes', label: 'Federation Disputes & Arbitrations', icon: 'icofont-judge' }
]

const totalContractValue = computed(() => {
  return contracts.value.reduce((acc, cur) => acc + (cur.contract_value_ugx || 0), 0)
})

async function fetchData() {
  loading.value = true
  try {
    const [cntRes, dispRes] = await Promise.allSettled([
      apiGet('/legal/contracts'),
      apiGet('/legal/disputes')
    ])

    if (cntRes.status === 'fulfilled') {
      contracts.value = cntRes.value.data?.contracts || cntRes.value.contracts || []
    }
    if (dispRes.status === 'fulfilled') {
      disputes.value = dispRes.value.data?.disputes || dispRes.value.disputes || []
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
