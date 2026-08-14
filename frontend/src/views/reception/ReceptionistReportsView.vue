<template>
  <LayoutDefault title="Front Desk & Application Reports">
    <div class="space-y-6 pb-12 w-full">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200 w-full">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-chart-bar-graph text-xs"></i> Reception Reports
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Front Desk & Application Reports</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Compiled statistics on daily visitor traffic, department clearance distributions, peak hours, and staff desk activities.
          </p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <button @click="fetchReports" class="btn btn-sm btn-light flex items-center gap-1.5" :disabled="loading">
            <i class="icofont-refresh" :class="{ 'animate-spin': loading }"></i> Refresh
          </button>
          <button @click="printReport" class="btn btn-sm btn-primary flex items-center gap-1.5">
            <i class="icofont-printer"></i> Print Daily Summary
          </button>
        </div>
      </div>

      <!-- Macro Summary KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 w-full">
        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Total Visitors Ingested</h5>
              <h2>{{ reportData?.summary?.total_visitors || 0 }}</h2>
              <span class="badge badge-primary"><i class="icofont-users"></i> Lifetime Intake</span>
            </div>
            <div class="banner-img bg-primary-light">
              <i class="icofont-users"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Clearance Rate</h5>
              <h2>{{ clearanceRate }}%</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Approved Passes</span>
            </div>
            <div class="banner-img bg-success-light">
              <i class="icofont-certificate"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Completed Missions</h5>
              <h2>{{ reportData?.summary?.completed_visits || 0 }}</h2>
              <span class="badge badge-info"><i class="icofont-badge"></i> Checked Out</span>
            </div>
            <div class="banner-img bg-cyan-light">
              <i class="icofont-paper-plane"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Active Queue</h5>
              <h2>{{ reportData?.summary?.pending_passes || 0 }}</h2>
              <span class="badge badge-warning"><i class="icofont-clock-time"></i> Pending</span>
            </div>
            <div class="banner-img bg-warning-light">
              <i class="icofont-sand-clock"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- Departmental Inflow & Peak Hours Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 w-full">
        
        <!-- Departmental Inflow Breakdown -->
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Departmental Visitor Traffic Distribution</h4>
          </div>
          <div class="card-body p-0">
            <table class="table table-striped text-xs text-left mb-0">
              <thead class="bg-slate-50 dark:bg-slate-800/60 text-slate-500 uppercase font-semibold">
                <tr>
                  <th class="py-2.5 px-4">Target Department</th>
                  <th class="py-2.5 px-4 font-mono">Visitor Inflow</th>
                  <th class="py-2.5 px-4">Share</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
                <tr v-for="d in reportData?.departmental_traffic || []" :key="d.department">
                  <td class="py-2.5 px-4 font-medium text-slate-900 dark:text-white">{{ d.department }}</td>
                  <td class="py-2.5 px-4 font-mono font-bold text-blue-600 dark:text-blue-400">{{ d.visits }} visits</td>
                  <td class="py-2.5 px-4">
                    <span class="badge badge-light">
                      {{ Math.round((d.visits / (reportData?.summary?.total_visitors || 1)) * 100) }}%
                    </span>
                  </td>
                </tr>
                <tr v-if="!(reportData?.departmental_traffic || []).length">
                  <td colspan="3" class="text-center py-6 text-slate-400">No departmental records today.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Peak Traffic Operating Windows -->
        <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5">
            <h4 class="text-sm font-bold text-slate-900 dark:text-white">Front Desk Peak Operating Windows</h4>
          </div>
          <div class="card-body p-4 space-y-3">
            <div
              v-for="h in reportData?.peak_traffic_hours || []"
              :key="h.hour"
              class="flex items-center justify-between p-3 rounded-xl bg-slate-50 dark:bg-slate-800/60 border border-slate-100 dark:border-slate-800 text-xs"
            >
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-blue-100 dark:bg-blue-900/40 text-blue-600 flex items-center justify-center font-bold">
                  <i class="icofont-clock-time"></i>
                </div>
                <div>
                  <div class="font-bold text-slate-900 dark:text-white">{{ h.hour }}</div>
                  <div class="text-[11px] text-slate-500">{{ h.label }}</div>
                </div>
              </div>
              <span class="badge" :class="h.volume === 'Peak' ? 'badge-danger' : h.volume === 'High' ? 'badge-warning' : 'badge-primary'">
                {{ h.volume }} Volume
              </span>
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
const reportData = ref(null)

const clearanceRate = computed(() => {
  const tot = reportData.value?.summary?.total_visitors || 0
  const app = (reportData.value?.summary?.approved_passes || 0) + (reportData.value?.summary?.completed_visits || 0)
  if (tot === 0) return 0
  return Math.round((app / tot) * 100)
})

async function fetchReports() {
  loading.value = true
  try {
    const res = await apiGet('/reception/reports')
    reportData.value = res.data || res || {}
  } catch {
    reportData.value = null
  } finally {
    loading.value = false
  }
}

function printReport() {
  window.print()
}

onMounted(() => {
  fetchReports()
})
</script>
