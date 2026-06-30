<template>
  <LayoutDefault title="Governance Dashboard">
    <div class="space-y-6">
      <section class="admin-card">
        <div class="admin-card-body flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
          <div>
            <p class="text-xs font-semibold uppercase tracking-wider text-primary-700">NSMIS official reporting</p>
            <h1 class="mt-1 text-xl font-bold text-gray-900">Federation governance and compliance</h1>
            <p class="mt-1 text-sm text-gray-500">Only approved compliance scores are shown. Missing reports update from active obligations.</p>
          </div>
          <label class="block min-w-64 text-sm font-medium text-gray-700">
            Reporting period
            <select v-model="periodID" class="mt-1 w-full rounded-lg border-gray-300" @change="loadDashboard">
              <option value="">Latest open period</option>
              <option v-for="period in periods" :key="period.id" :value="period.id">{{ period.name }}</option>
            </select>
          </label>
        </div>
      </section>

      <div v-if="error" role="alert" class="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-700">
        {{ error }}
        <button class="ml-2 font-semibold underline" @click="loadDashboard">Try again</button>
      </div>

      <div v-if="loading" aria-label="Loading governance dashboard" class="grid grid-cols-2 gap-4 xl:grid-cols-5">
        <div v-for="i in 5" :key="i" class="h-28 animate-pulse rounded-xl bg-gray-100"></div>
      </div>
      <template v-else>
        <div class="grid grid-cols-2 gap-4 xl:grid-cols-5">
          <article v-for="item in cards" :key="item.label" class="stat-card">
            <div :class="[item.bg, 'flex h-11 w-11 items-center justify-center rounded-xl']"><i :class="[item.icon, item.color, 'text-xl']"></i></div>
            <div><p class="text-2xl font-bold text-gray-900">{{ item.value }}</p><p class="text-xs text-gray-500">{{ item.label }}</p></div>
          </article>
        </div>

        <section class="admin-card overflow-hidden">
          <div class="admin-card-header">
            <div><h2 class="text-sm font-semibold text-gray-900">Federation status</h2><p class="mt-0.5 text-xs text-gray-500">{{ dashboard.reporting_period_name || 'No open reporting period' }}</p></div>
            <div class="text-right text-xs text-gray-500"><div>As of {{ formatDateTime(dashboard.as_of) }}</div><div>Calculation v{{ dashboard.calculation_version || 1 }}</div></div>
          </div>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead><tr><th class="table-th">Federation</th><th class="table-th">Compliance</th><th class="table-th">Score</th><th class="table-th">Missing reports</th><th class="table-th">Constitution</th></tr></thead>
              <tbody v-if="dashboard.federations?.length" class="divide-y divide-gray-100">
                <tr v-for="row in dashboard.federations" :key="row.federation_id" class="hover:bg-gray-50">
                  <td class="table-td"><span class="font-medium text-gray-900">{{ row.federation_name }}</span><span class="ml-2 text-xs text-gray-400">{{ row.acronym }}</span></td>
                  <td class="table-td"><span :class="statusClass(row.is_compliant)">{{ complianceLabel(row.is_compliant) }}</span></td>
                  <td class="table-td">{{ row.compliance_score == null ? 'Not calculated' : `${Number(row.compliance_score).toFixed(1)}%` }}</td>
                  <td class="table-td"><span :class="row.missing_reports ? 'font-semibold text-red-700' : 'text-green-700'">{{ row.missing_reports }}</span></td>
                  <td class="table-td"><span :class="row.expired_constitution ? 'text-red-700' : 'text-green-700'">{{ row.expired_constitution ? 'Expired' : 'Current / not recorded' }}</span></td>
                </tr>
              </tbody>
              <tbody v-else><tr><td colspan="5" class="px-6 py-14 text-center text-sm text-gray-500">No federations are available in your assigned scope.</td></tr></tbody>
            </table>
          </div>
        </section>
      </template>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import apiClient from '@/api/client.js'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'

const breadcrumb = useBreadcrumbStore()
const loading = ref(true)
const error = ref('')
const periods = ref([])
const periodID = ref('')
const dashboard = ref({ federations: [] })

const cards = computed(() => [
  { label: 'Active federations', value: dashboard.value.total_federations ?? 0, icon: 'icofont-building-alt', bg: 'bg-blue-50', color: 'text-blue-700' },
  { label: 'Compliant', value: dashboard.value.compliant_federations ?? 0, icon: 'icofont-check-circled', bg: 'bg-green-50', color: 'text-green-700' },
  { label: 'Non-compliant', value: dashboard.value.non_compliant_federations ?? 0, icon: 'icofont-warning-alt', bg: 'bg-orange-50', color: 'text-orange-700' },
  { label: 'Missing reports', value: dashboard.value.missing_reports ?? 0, icon: 'icofont-file-alt', bg: 'bg-red-50', color: 'text-red-700' },
  { label: 'Expired constitutions', value: dashboard.value.expired_constitutions ?? 0, icon: 'icofont-ui-calendar', bg: 'bg-purple-50', color: 'text-purple-700' }
])

function statusClass(value) { return ['inline-flex rounded-full px-2.5 py-1 text-xs font-semibold', value == null ? 'bg-gray-100 text-gray-600' : value ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'] }
function complianceLabel(value) { return value == null ? 'Not calculated' : value ? 'Compliant' : 'Non-compliant' }
function formatDateTime(value) { return value ? new Intl.DateTimeFormat('en-UG', { dateStyle: 'medium', timeStyle: 'short', timeZone: 'Africa/Kampala' }).format(new Date(value)) : '—' }

async function loadDashboard() {
  loading.value = true; error.value = ''
  try {
    const params = periodID.value ? { period_id: periodID.value } : {}
    const res = await apiClient.get('/api/v1/nsmis/dashboards/governance', { params })
    dashboard.value = res.data?.data || { federations: [] }
  } catch (err) { error.value = err.response?.data?.error?.message || 'The governance dashboard could not be loaded.' }
  finally { loading.value = false }
}

onMounted(async () => {
  breadcrumb.set('Governance Dashboard')
  try { const res = await apiClient.get('/api/v1/nsmis/reporting-periods', { params: { open: false } }); periods.value = res.data?.data || [] } catch { periods.value = [] }
  await loadDashboard()
})
</script>
