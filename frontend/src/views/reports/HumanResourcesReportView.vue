<template>
  <LayoutDefault title="Human Resources Report">
    <div class="space-y-6 pb-12">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-users text-xs"></i> Directorate Report
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Human Resources & Administration Report</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Compiled HR dossier covering staff leave & duty handover rosters, annual performance appraisals, staff establishment headcount, and continuous capacity building.
          </p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <button @click="fetchData" class="btn btn-sm btn-light flex items-center gap-1.5" :disabled="loading">
            <i class="icofont-refresh" :class="{ 'animate-spin': loading }"></i> Refresh
          </button>
          <button @click="printDossier" class="btn btn-sm btn-primary flex items-center gap-1.5">
            <i class="icofont-printer"></i> Print Dossier
          </button>
        </div>
      </div>

      <!-- Department Overview Banner -->
      <div v-if="reportData" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>HR Compliance Score</h5>
              <h2>{{ reportData.overall_compliance_score || 0 }}%</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Staff Establishment</span>
            </div>
            <div class="banner-img">
              <i class="icofont-users-alt-4"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Compiled Modules</h5>
              <h2>{{ (reportData.sub_reports || []).length }}</h2>
              <span class="badge badge-primary"><i class="icofont-files-stack"></i> HR Sections</span>
            </div>
            <div class="banner-img">
              <i class="icofont-listine-dots"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Head of HR</h5>
              <h2 class="text-xs font-semibold truncate">{{ reportData.head_of_department || 'HOD HR' }}</h2>
              <span class="badge badge-info"><i class="icofont-user"></i> Senior Manager</span>
            </div>
            <div class="banner-img">
              <i class="icofont-id"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Audit Status</h5>
              <h2>0 Queries</h2>
              <span class="badge badge-success"><i class="icofont-shield-check"></i> Statutory Compliant</span>
            </div>
            <div class="banner-img">
              <i class="icofont-certificate"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- Detailed Sub-Reports Section -->
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-bold text-slate-900 dark:text-white">Compiled HR Sections</h3>
          <span class="text-xs text-slate-500 font-mono">Last Ingested: {{ lastUpdated }}</span>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div 
            v-for="sub in reportData?.sub_reports || []" 
            :key="sub.id"
            class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 p-5 space-y-3 mb-0"
          >
            <div class="flex items-center justify-between gap-2">
              <span class="badge badge-info text-[11px]">{{ sub.category }}</span>
              <span class="badge badge-success text-[11px]">{{ sub.status }}</span>
            </div>
            
            <h4 class="text-base font-bold text-slate-900 dark:text-white">{{ sub.title }}</h4>
            <p class="text-xs text-slate-600 dark:text-slate-400 leading-relaxed">{{ sub.summary }}</p>

            <div class="pt-3 border-t border-slate-100 dark:border-slate-800 grid grid-cols-2 gap-2 text-xs">
              <div v-for="(val, key) in sub.key_metrics || {}" :key="key" class="bg-slate-50 dark:bg-slate-800/60 p-2 rounded-lg border border-slate-100 dark:border-slate-800">
                <span class="text-[10px] uppercase font-bold text-slate-400 block truncate">{{ formatKey(key) }}</span>
                <strong class="text-slate-800 dark:text-slate-200 font-mono text-xs">{{ val }}</strong>
              </div>
            </div>
          </div>

          <div v-if="!(reportData?.sub_reports || []).length" class="col-span-2 text-center py-12 text-xs text-slate-400">
            No compiled HR reports submitted for this cycle.
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
const reportData = ref(null)
const lastUpdated = ref('Live')

function formatKey(key) {
  return key ? key.replace(/_/g, ' ') : ''
}

function printDossier() {
  window.print()
}

async function fetchData() {
  loading.value = true
  try {
    const res = await apiGet('/executive/appraisal/reports/departmental')
    const reports = res.data?.departmental_reports || res.departmental_reports || {}
    reportData.value = reports.human_resources || null
    lastUpdated.value = new Date().toLocaleTimeString('en-GB')
  } catch {
    reportData.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
