<template>
  <LayoutDefault title="Fixed Assets Registry - Value Adjustments" :hide-universal-apps="true">
    <div class="space-y-6">
      <!-- Shared Header & Sub-Navigation -->
      <FixedAssetsSubNav
        title="Value Adjustments & Revaluation Audit Ledger"
        :badge="totalLogs + ' Transaction Events'"
      />

      <!-- Notice Banner -->
      <div
        v-if="notice.message"
        :class="notice.type === 'error' ? 'border-red-200 bg-red-50 text-red-700' : 'border-emerald-200 bg-emerald-50 text-emerald-700'"
        class="rounded-xl border p-4 text-xs font-semibold flex items-center justify-between shadow-sm transition"
      >
        <div class="flex items-center gap-2">
          <i :class="notice.type === 'error' ? 'icofont-warning text-red-500' : 'icofont-check-circled text-emerald-600'" class="text-lg"></i>
          <span>{{ notice.message }}</span>
        </div>
        <button @click="notice.message = ''" class="text-gray-400 hover:text-gray-600">
          <i class="icofont-close"></i>
        </button>
      </div>

      <!-- Stat Summary Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Total Audit Log Events</p>
              <h3 class="text-xl font-extrabold text-gray-900 mt-1 font-mono">{{ totalLogs }}</h3>
            </div>
            <div class="p-3 bg-blue-50 rounded-lg text-blue-600"><i class="icofont-history text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-gray-500">Statutory audit trail per PFMA 2015</div>
        </div>

        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Revaluation Adjustments</p>
              <h3 class="text-xl font-extrabold text-indigo-900 mt-1 font-mono">{{ revaluationCount }}</h3>
            </div>
            <div class="p-3 bg-indigo-50 rounded-lg text-indigo-600"><i class="icofont-chart-growth text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-indigo-600 font-medium">Statutory price overrides & appraisals</div>
        </div>

        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Asset Verifications</p>
              <h3 class="text-xl font-extrabold text-emerald-900 mt-1 font-mono">{{ verificationCount }}</h3>
            </div>
            <div class="p-3 bg-emerald-50 rounded-lg text-emerald-600"><i class="icofont-check-circled text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-emerald-700 font-medium">Physical spot check verifications</div>
        </div>

        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Deletions / Write-offs</p>
              <h3 class="text-xl font-extrabold text-red-900 mt-1 font-mono">{{ deletionCount }}</h3>
            </div>
            <div class="p-3 bg-red-50 rounded-lg text-red-600"><i class="icofont-trash text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-red-700 font-medium">Audited asset removals with user tracking</div>
        </div>
      </div>

      <!-- Main Log Table Card -->
      <div class="bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden">
        <!-- Controls & Filters -->
        <div class="p-4 border-b border-gray-100 bg-gray-50/50 flex flex-col md:flex-row items-center justify-between gap-4">
          <div class="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
            <div class="relative w-full md:w-80">
              <input
                v-model="searchQuery"
                @input="debouncedFetchLogs"
                type="text"
                placeholder="Search Asset Code, User Name, Notes..."
                aria-label="Search transaction logs"
                class="log-search w-full pl-9 pr-4 py-2 text-xs border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
              />
              <i class="icofont-search-1 text-gray-400 absolute left-3 top-2.5"></i>
            </div>

            <select
              aria-label="Filter transaction type"
              v-model="selectedType"
              @change="offset = 0; fetchLogs()"
              class="text-xs border border-gray-300 rounded-lg py-2 px-3 focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 bg-white"
            >
              <option value="">All Event Types</option>
              <option value="REVALUATION">Revaluations</option>
              <option value="VERIFICATION">Verifications</option>
              <option value="DELETION">Deletions</option>
              <option value="DEPRECIATION">Depreciations</option>
            </select>
          </div>

          <div class="text-xs text-gray-500 font-medium flex items-center gap-3">
            <span>Showing {{ logs.length }} of {{ totalLogs }} audit records</span>
            <button @click="searchQuery = ''; selectedType = ''; offset = 0; fetchLogs()" class="text-emerald-700 hover:underline">Reset</button>
          </div>
        </div>

        <!-- Table -->
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-gray-100/70 text-gray-600 font-semibold uppercase tracking-wider border-b border-gray-200">
              <tr>
                <th class="py-3 px-4">Date & Time</th>
                <th class="py-3 px-4">Event Type</th>
                <th class="py-3 px-4">Asset Code</th>
                <th class="py-3 px-4">Asset Description</th>
                <th class="py-3 px-4 text-right">Previous Cost</th>
                <th class="py-3 px-4 text-right">New Valuation</th>
                <th class="py-3 px-4 text-right">Net Difference</th>
                <th class="py-3 px-4">Performed By</th>
                <th class="py-3 px-4">Notes & Rationale</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-if="loading">
                <td colspan="9" class="py-8 text-center text-gray-400">Loading revaluation transaction logs...</td>
              </tr>
              <tr v-else-if="!logs.length">
                <td colspan="9" class="py-8 text-center text-gray-400">No transaction log records found matching search.</td>
              </tr>
              <template v-else>
                <tr v-for="log in logs" :key="log.id" class="hover:bg-blue-50/30 transition">
                  <td class="py-3 px-4 text-gray-500 whitespace-nowrap font-mono">
                    {{ formatDate(log.created_at) }}
                  </td>
                  <td class="py-3 px-4">
                    <span :class="badgeClass(log.transaction_type)" class="px-2.5 py-1 rounded-full text-[10px] font-bold tracking-wide">
                      {{ log.transaction_type }}
                    </span>
                  </td>
                  <td class="py-3 px-4 font-mono font-bold text-emerald-800">
                    {{ log.asset_number || 'N/A' }}
                  </td>
                  <td class="py-3 px-4 text-gray-900 font-medium max-w-xs truncate" :title="log.asset_description">
                    {{ log.asset_description || 'Fixed Asset Record' }}
                  </td>
                  <td class="py-3 px-4 text-right font-mono text-gray-700">
                    {{ formatUGX(log.previous_val) }}
                  </td>
                  <td class="py-3 px-4 text-right font-mono font-bold text-blue-900">
                    {{ formatUGX(log.new_val) }}
                  </td>
                  <td class="py-3 px-4 text-right font-mono font-bold" :class="diffClass(log.new_val - log.previous_val)">
                    {{ formatDiff(log.new_val - log.previous_val, log.transaction_type) }}
                  </td>
                  <td class="py-3 px-4 whitespace-nowrap">
                    <div class="flex items-center gap-1.5">
                      <div class="w-6 h-6 rounded-full bg-emerald-100 text-emerald-800 text-[10px] font-bold flex items-center justify-center flex-shrink-0">
                        {{ (log.performed_by_name || 'U')[0].toUpperCase() }}
                      </div>
                      <span class="font-semibold text-gray-800 text-xs">{{ log.performed_by_name || 'System User' }}</span>
                    </div>
                  </td>
                  <td class="py-3 px-4 text-gray-600 max-w-sm truncate" :title="log.notes">
                    {{ log.notes || '—' }}
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pagination -->
      <div class="flex flex-wrap items-center justify-between gap-3 bg-white border border-gray-200 rounded-xl p-4 shadow-sm" aria-label="Transaction logs pagination">
        <span class="text-sm text-gray-700">
          Showing {{ totalLogs ? offset + 1 : 0 }}–{{ Math.min(offset + logs.length, totalLogs) }} of {{ totalLogs }} log events
        </span>
        <div class="flex gap-3">
          <button
            :disabled="loading || offset === 0"
            @click="offset = Math.max(0, offset - pageSize); fetchLogs()"
            class="px-4 py-2 border border-gray-300 rounded-lg text-xs font-semibold hover:bg-gray-50 disabled:opacity-40"
          >
            Previous page
          </button>
          <button
            :disabled="loading || offset + pageSize >= totalLogs"
            @click="offset += pageSize; fetchLogs()"
            class="px-4 py-2 border border-emerald-700 text-emerald-800 rounded-lg text-xs font-semibold hover:bg-emerald-50 disabled:opacity-40"
          >
            Next page
          </button>
        </div>
      </div>

    </div>
  </LayoutDefault>
</template>

<script>
import client from '@/api/client'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import FixedAssetsSubNav from '@/components/fixed_assets/FixedAssetsSubNav.vue'

function apiData(res) {
  return res?.data?.data ?? res?.data ?? {}
}

export default {
  name: 'FixedAssetValueAdjustmentsView',
  components: { LayoutDefault, FixedAssetsSubNav },
  data() {
    return {
      logs: [],
      totalLogs: 0,
      offset: 0,
      pageSize: 50,
      loading: false,
      fetchTimer: null,
      searchQuery: '',
      selectedType: '',
      notice: { type: '', message: '' },
    }
  },
  computed: {
    revaluationCount() {
      return this.logs.filter(l => l.transaction_type === 'REVALUATION').length
    },
    verificationCount() {
      return this.logs.filter(l => l.transaction_type === 'VERIFICATION').length
    },
    deletionCount() {
      return this.logs.filter(l => l.transaction_type === 'DELETION').length
    },
  },
  mounted() {
    this.fetchLogs()
  },
  beforeUnmount() {
    if (this.fetchTimer) clearTimeout(this.fetchTimer)
  },
  methods: {
    async fetchLogs() {
      this.loading = true
      try {
        const params = { limit: this.pageSize, offset: this.offset, search: this.searchQuery, type: this.selectedType }
        const res = await client.get('/api/v1/assets/logs', { params })
        const data = apiData(res)
        this.logs = data.logs || []
        this.totalLogs = data.total || 0
      } catch (e) {
        this.notice = { type: 'error', message: e.response?.data?.error || 'Failed to load transaction audit logs' }
      } finally {
        this.loading = false
      }
    },
    debouncedFetchLogs() {
      this.offset = 0
      if (this.fetchTimer) clearTimeout(this.fetchTimer)
      this.fetchTimer = setTimeout(() => this.fetchLogs(), 250)
    },
    formatUGX(val) {
      if (val === null || val === undefined) return '—'
      return new Intl.NumberFormat('en-UG', { maximumFractionDigits: 2 }).format(val || 0)
    },
    formatDiff(diff, type) {
      if (type === 'VERIFICATION') return '—'
      if (diff > 0) return `+${this.formatUGX(diff)}`
      if (diff < 0) return `-${this.formatUGX(Math.abs(diff))}`
      return '0.00'
    },
    diffClass(diff) {
      if (diff > 0) return 'text-emerald-700'
      if (diff < 0) return 'text-red-700'
      return 'text-gray-500'
    },
    badgeClass(type) {
      switch (type) {
        case 'REVALUATION': return 'bg-blue-100 text-blue-800 border border-blue-200'
        case 'VERIFICATION': return 'bg-emerald-100 text-emerald-800 border border-emerald-200'
        case 'DELETION': return 'bg-red-100 text-red-800 border border-red-200'
        case 'DEPRECIATION': return 'bg-purple-100 text-purple-800 border border-purple-200'
        default: return 'bg-gray-100 text-gray-800 border border-gray-200'
      }
    },
    formatDate(dateStr) {
      if (!dateStr) return '—'
      try {
        const d = new Date(dateStr)
        return d.toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
      } catch {
        return dateStr
      }
    },
  },
}
</script>

<style scoped>
.log-search {
  padding-left: 2.25rem !important;
}
</style>
