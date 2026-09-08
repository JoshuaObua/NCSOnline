<template>
  <LayoutDefault title="Fixed Assets Registry - Dynamic Pivot Engine">
    <div class="space-y-6">
      <FixedAssetsPageHeader
        title="Dynamic Multi-Dimensional Pivot Engine"
        :badge="filteredSummaries.length + ' Category Segments'"
      />

      <!-- Stat Summary Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Categories Cross-Tabulated</p>
              <h3 class="text-xl font-extrabold text-gray-900 mt-1 font-mono">{{ summary.category_summaries?.length || 0 }}</h3>
            </div>
            <div class="p-3 bg-purple-50 rounded-lg text-purple-600"><i class="icofont-chart-histogram text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-gray-500">Sub-classes across 10 asset categories</div>
        </div>

        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Total Portfolio Units</p>
              <h3 class="text-xl font-extrabold text-emerald-900 mt-1 font-mono">{{ totalUnits }}</h3>
            </div>
            <div class="p-3 bg-emerald-50 rounded-lg text-emerald-600"><i class="icofont-box text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-emerald-600 font-medium">Aggregated physical item count</div>
        </div>

        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Total Adjusted Valuation</p>
              <h3 class="text-xl font-extrabold text-blue-900 mt-1 font-mono">UGX {{ formatUGX(summary.total_adjusted_cost) }}</h3>
            </div>
            <div class="p-3 bg-blue-50 rounded-lg text-blue-600"><i class="icofont-chart-growth text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-blue-600 font-medium">Statutory valuation total</div>
        </div>

        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Total Net Book Value</p>
              <h3 class="text-xl font-extrabold text-purple-900 mt-1 font-mono">UGX {{ formatUGX(summary.total_net_book_value) }}</h3>
            </div>
            <div class="p-3 bg-purple-50 rounded-lg text-purple-600"><i class="icofont-calculator-alt-2 text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-purple-600 font-medium">Net carrying amount (IPSAS 17)</div>
        </div>
      </div>

      <!-- Pivot Engine Table Card -->
      <div class="bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden p-6 space-y-4">
        <div class="flex flex-col md:flex-row items-center justify-between gap-4">
          <div>
            <h3 class="text-base font-bold text-gray-900">Multi-Dimensional Category Matrix</h3>
            <p class="text-xs text-gray-500">Live analytics cross-tabulating primary class (`SEGMENT1`), category (`SEGMENT3`), detailed sub-class (`SEGMENT4`), units, and valuation.</p>
          </div>

          <div class="flex items-center gap-3 w-full md:w-auto">
            <div class="relative w-full md:w-72">
              <input
                v-model="pivotFilter"
                type="text"
                placeholder="Filter matrix category..."
                aria-label="Filter matrix category"
                class="pivot-search w-full pl-9 pr-4 py-2 text-xs border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
              />
              <i class="icofont-search-1 text-gray-400 absolute left-3 top-2.5"></i>
            </div>

            <button
              @click="exportCSV"
              class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-xs font-bold shadow-sm transition flex items-center gap-1.5 whitespace-nowrap"
            >
              <i class="icofont-file-excel text-sm"></i>
              Export Pivot CSV
            </button>
          </div>
        </div>

        <div class="overflow-x-auto border border-gray-200 rounded-xl">
          <table class="w-full text-left text-xs">
            <thead class="bg-gray-100 text-gray-700 font-bold uppercase tracking-wider border-b border-gray-200">
              <tr>
                <th class="py-3 px-4">Primary Class (`SEGMENT1`)</th>
                <th class="py-3 px-4">Category (`SEGMENT3`)</th>
                <th class="py-3 px-4">Detailed Sub-Class (`SEGMENT4`)</th>
                <th class="py-3 px-4 text-right">Sum of Units</th>
                <th class="py-3 px-4 text-right">Sum of FB Cost (UGX)</th>
                <th class="py-3 px-4 text-right">Sum of Adjusted Cost</th>
                <th class="py-3 px-4 text-right">Sum of NBV</th>
                <th class="py-3 px-4 text-right">% Portfolio Share</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-if="loading">
                <td colspan="8" class="py-8 text-center text-gray-400">Loading pivot analytics engine...</td>
              </tr>
              <tr v-else-if="!filteredSummaries.length">
                <td colspan="8" class="py-8 text-center text-gray-400">No category breakdown rows match filter.</td>
              </tr>
              <template v-else>
                <tr v-for="(item, idx) in filteredSummaries" :key="idx" class="hover:bg-purple-50/20 transition">
                  <td class="py-3 px-4 font-semibold text-gray-800">{{ item.category_segment1 }}</td>
                  <td class="py-3 px-4 font-medium text-emerald-800">{{ item.category_segment3 }}</td>
                  <td class="py-3 px-4 text-gray-600">{{ item.category_segment4 }}</td>
                  <td class="py-3 px-4 text-right font-mono font-bold text-gray-900">{{ item.asset_units }}</td>
                  <td class="py-3 px-4 text-right font-mono text-gray-700">{{ formatUGX(item.total_fb_cost) }}</td>
                  <td class="py-3 px-4 text-right font-mono font-bold text-blue-900">{{ formatUGX(item.total_adjusted_cost) }}</td>
                  <td class="py-3 px-4 text-right font-mono font-bold text-purple-900">{{ formatUGX(item.total_net_book_value) }}</td>
                  <td class="py-3 px-4 text-right font-mono font-bold text-emerald-700">
                    {{ calcShare(item.total_adjusted_cost) }}%
                  </td>
                </tr>
              </template>
            </tbody>
            <!-- Matrix Totals Footer -->
            <tfoot v-if="filteredSummaries.length" class="bg-gray-100/90 font-bold border-t-2 border-gray-300 text-gray-900">
              <tr>
                <td colspan="3" class="py-3.5 px-4 text-right uppercase tracking-wider">Grand Total ({{ filteredSummaries.length }} Segments):</td>
                <td class="py-3.5 px-4 text-right font-mono text-emerald-900 text-sm">{{ totalUnits }}</td>
                <td class="py-3.5 px-4 text-right font-mono text-gray-800">{{ formatUGX(summary.total_fb_cost) }}</td>
                <td class="py-3.5 px-4 text-right font-mono text-blue-950 text-sm">UGX {{ formatUGX(summary.total_adjusted_cost) }}</td>
                <td class="py-3.5 px-4 text-right font-mono text-purple-950 text-sm">UGX {{ formatUGX(summary.total_net_book_value) }}</td>
                <td class="py-3.5 px-4 text-right font-mono text-emerald-800">100.0%</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>

    </div>
  </LayoutDefault>
</template>

<script>
import client from '@/api/client'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import FixedAssetsPageHeader from '@/components/fixed_assets/FixedAssetsPageHeader.vue'

function apiData(res) {
  return res?.data?.data ?? res?.data ?? {}
}

export default {
  name: 'FixedAssetPivotEngineView',
  components: { LayoutDefault, FixedAssetsPageHeader },
  data() {
    return {
      loading: false,
      pivotFilter: '',
      summary: {
        total_assets: null,
        total_fb_cost: null,
        total_adjusted_cost: null,
        total_net_book_value: null,
        category_summaries: [],
      },
    }
  },
  computed: {
    filteredSummaries() {
      const q = (this.pivotFilter || '').toLowerCase().trim()
      if (!q) return this.summary.category_summaries || []
      return (this.summary.category_summaries || []).filter(item => {
        return (
          (item.category_segment1 || '').toLowerCase().includes(q) ||
          (item.category_segment3 || '').toLowerCase().includes(q) ||
          (item.category_segment4 || '').toLowerCase().includes(q)
        )
      })
    },
    totalUnits() {
      return (this.summary.category_summaries || []).reduce((acc, curr) => acc + (curr.asset_units || 0), 0)
    },
  },
  mounted() {
    this.fetchSummary()
  },
  methods: {
    async fetchSummary() {
      this.loading = true
      try {
        const res = await client.get('/api/v1/assets/summary')
        const data = apiData(res)
        if (data) this.summary = { ...this.summary, ...data }
      } catch (e) {
        console.error('Failed to load asset summary for pivot engine', e)
      } finally {
        this.loading = false
      }
    },
    formatUGX(val) {
      if (val === null || val === undefined) return '—'
      return new Intl.NumberFormat('en-UG', { maximumFractionDigits: 2 }).format(val || 0)
    },
    calcShare(cost) {
      if (!this.summary.total_adjusted_cost || !cost) return '0.0'
      const pct = (cost / this.summary.total_adjusted_cost) * 100
      return pct.toFixed(1)
    },
    exportCSV() {
      const headers = ['Primary Class (SEGMENT1)', 'Category (SEGMENT3)', 'Detailed Sub-Class (SEGMENT4)', 'Units', 'FB Cost (UGX)', 'Adjusted Cost (UGX)', 'NBV (UGX)', 'Share (%)']
      const rows = this.filteredSummaries.map(item => [
        `"${item.category_segment1 || ''}"`,
        `"${item.category_segment3 || ''}"`,
        `"${item.category_segment4 || ''}"`,
        item.asset_units,
        item.total_fb_cost,
        item.total_adjusted_cost,
        item.total_net_book_value,
        this.calcShare(item.total_adjusted_cost)
      ])
      const csvContent = [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
      const url = URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.setAttribute('href', url)
      link.setAttribute('download', `Fixed_Assets_Pivot_Engine_${new Date().toISOString().slice(0,10)}.csv`)
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    },
  },
}
</script>

<style scoped>
.pivot-search {
  padding-left: 2.25rem !important;
}
</style>
