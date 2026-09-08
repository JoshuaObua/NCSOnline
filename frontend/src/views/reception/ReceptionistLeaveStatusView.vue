<template>
  <LayoutDefault title="Leave Status & History">
    <div class="space-y-6 pb-12 w-full">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200 w-full">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-history text-xs"></i> Leave Tracking
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Leave Applications & Approval Status</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5 max-w-3xl">
            Track supervisor review, HR vetting, and entitlement balance for submitted leave requests.
          </p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <button @click="fetchLeaves" class="btn btn-sm btn-light flex items-center gap-1.5" :disabled="loading">
            <i class="icofont-refresh" :class="{ 'animate-spin': loading }"></i> Refresh
          </button>
          <router-link to="/reception/leave/apply" class="btn btn-sm btn-primary flex items-center gap-1.5 shadow-sm">
            <i class="icofont-plus"></i> Apply for Leave
          </router-link>
        </div>
      </div>

      <!-- Leave Summary Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 w-full">
        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Annual Entitlement</h5>
              <h2>21 Days</h2>
              <span class="badge badge-primary"><i class="icofont-calendar"></i> Statutory Days</span>
            </div>
            <div class="banner-img bg-primary-light">
              <i class="icofont-calendar"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Days Utilized</h5>
              <h2>{{ utilizedDays }} Days</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Approved Leaves</span>
            </div>
            <div class="banner-img bg-success-light">
              <i class="icofont-ui-clock"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Remaining Balance</h5>
              <h2>{{ Math.max(0, 21 - utilizedDays) }} Days</h2>
              <span class="badge badge-info"><i class="icofont-badge"></i> Available Balance</span>
            </div>
            <div class="banner-img bg-cyan-light">
              <i class="icofont-sun-alt"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- Leave Applications Table -->
      <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden w-full">
        <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5 flex items-center justify-between">
          <h4 class="text-sm font-bold text-slate-900 dark:text-white">Submitted Leave History</h4>
          <router-link to="/reception/leave/apply" class="btn btn-sm btn-primary flex items-center gap-1.5 text-xs">
            <i class="icofont-plus"></i> New Application
          </router-link>
        </div>

        <div class="card-body p-0 overflow-x-auto">
          <table class="table table-striped table-hover mb-0 text-left text-xs">
            <thead class="bg-slate-50 dark:bg-slate-800/60 text-slate-600 dark:text-slate-400 uppercase tracking-wider font-semibold">
              <tr>
                <th class="py-3 px-4">Ref No / Type</th>
                <th class="py-3 px-4">Leave Period</th>
                <th class="py-3 px-4">Days</th>
                <th class="py-3 px-4">Relieving Officer</th>
                <th class="py-3 px-4">Status</th>
                <th class="py-3 px-4">Remarks / Action</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
              <tr v-for="l in leaves" :key="l.id">
                <td class="py-3 px-4">
                  <div class="font-mono font-bold text-blue-600 dark:text-blue-400">{{ l.reference_no }}</div>
                  <div class="text-[11px] text-slate-500 font-semibold">{{ formatLeaveType(l.leave_type) }}</div>
                </td>
                <td class="py-3 px-4">
                  <div class="font-bold text-slate-800 dark:text-slate-200">{{ formatDate(l.start_date) }} – {{ formatDate(l.end_date) }}</div>
                  <div class="text-[11px] text-slate-500">Submitted: {{ formatDate(l.created_at) }}</div>
                </td>
                <td class="py-3 px-4 font-mono font-bold text-slate-900 dark:text-white">
                  {{ l.days_requested }} Days
                </td>
                <td class="py-3 px-4 text-slate-700 dark:text-slate-300">
                  <div class="font-medium text-slate-900 dark:text-white">{{ l.relieving_officer_name }}</div>
                  <div v-if="l.relieving_officer_role" class="text-[10px] text-slate-400">{{ l.relieving_officer_role }}</div>
                </td>
                <td class="py-3 px-4">
                  <span class="badge" :class="leaveStatusBadge(l.status)">
                    {{ formatStatus(l.status) }}
                  </span>
                  <div v-if="l.approved_at" class="text-[10px] text-slate-400 mt-1">
                    <i class="icofont-check-circled text-emerald-500"></i> {{ formatDate(l.approved_at) }}
                  </div>
                </td>
                <td class="py-3 px-4 text-[11px]">
                  <div v-if="l.supervisor_remarks" class="p-2 rounded bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-700 dark:text-slate-300">
                    <strong class="text-slate-900 dark:text-white d-block mb-0.5"><i class="icofont-speech-comments text-blue-500 mr-1"></i> HR / Supervisor Remark:</strong>
                    {{ l.supervisor_remarks }}
                  </div>
                  <div v-else class="text-slate-400 italic">
                    <i class="icofont-clock-time mr-1"></i> Awaiting HR supervisor review & response
                  </div>
                </td>
              </tr>
              <tr v-if="!leaves.length">
                <td colspan="6" class="text-center py-8 text-xs text-slate-400">
                  No staff leave applications submitted yet.
                </td>
              </tr>
            </tbody>
          </table>
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
const leaves = ref([])

const utilizedDays = computed(() => {
  return leaves.value
    .filter(l => l.status === 'APPROVED_HR' && l.leave_type === 'ANNUAL_LEAVE')
    .reduce((acc, cur) => acc + (cur.days_requested || 0), 0)
})

function formatLeaveType(t) {
  if (!t) return ''
  return t.replace(/_/g, ' ')
}

function formatStatus(s) {
  if (!s) return ''
  return s.replace(/_/g, ' ')
}

function leaveStatusBadge(status) {
  switch (status) {
    case 'PENDING_SUPERVISOR': return 'badge-warning'
    case 'APPROVED_HR': return 'badge-success'
    case 'REJECTED': return 'badge-danger'
    default: return 'badge-primary'
  }
}

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric' })
}

async function fetchLeaves() {
  loading.value = true
  try {
    const res = await apiGet('/reception/leave/status')
    leaves.value = res.data?.leaves || res.leaves || []
  } catch {
    leaves.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchLeaves()
})
</script>
