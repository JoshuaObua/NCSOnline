<template>
  <LayoutDefault title="Fixed Assets & IPSAS 17 Register">
    <div class="space-y-6">
      <!-- Top Header & Breadcrumb -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-white p-6 rounded-xl border border-gray-100 shadow-sm">
        <div>
          <div class="flex items-center space-x-2 text-sm text-gray-500 mb-1">
            <span>Finance & Accounting</span>
            <span>/</span>
            <span class="text-emerald-700 font-medium">Fixed Asset & Inventory Master Register</span>
          </div>
          <h1 class="text-2xl font-bold text-gray-900 tracking-tight flex flex-wrap items-center gap-2">
            Fixed Assets & IPSAS 17 Register
            <span class="text-xs bg-emerald-100 text-emerald-800 font-semibold px-2.5 py-0.5 rounded-full border border-emerald-200">
              {{ summary.total_assets ?? '—' }} Registered Assets
            </span>
          </h1>
          <p class="text-xs text-gray-500 mt-1">
            Asset register based on <span class="font-mono text-gray-700 font-semibold">Docs/FIXED ASSET REGISTER ADJUSTMENTS.xlsx</span> (UGX 31.02B FB_COST baseline).
          </p>
        </div>

        <div class="flex flex-wrap items-center gap-3">
          <router-link
            to="/fixed-assets/depreciation"
            class="inline-flex items-center px-4 py-2 bg-indigo-50 text-indigo-700 rounded-lg hover:bg-indigo-100 text-xs font-semibold transition"
          >
            <i class="icofont-clock-time mr-1.5"></i>
            Run Monthly Depreciation
          </router-link>
          <router-link
            to="/fixed-assets/new"
            class="inline-flex items-center px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 text-xs font-semibold shadow-sm transition"
          >
            <i class="icofont-plus-circle mr-1.5"></i>
            New Fixed Asset
          </router-link>
        </div>
      </div>

      <div v-if="notice.message" :class="notice.type === 'error' ? 'border-red-200 bg-red-50 text-red-700' : 'border-emerald-200 bg-emerald-50 text-emerald-700'" class="rounded-xl border p-3 text-xs font-semibold">
        {{ notice.message }}
      </div>

      <!-- Financial KPI Summary Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm relative overflow-hidden">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Total Portfolio FB Cost</p>
              <h3 class="text-xl font-extrabold text-gray-900 mt-1 font-mono">UGX {{ formatUGX(summary.total_fb_cost) }}</h3>
            </div>
            <div class="p-3 bg-emerald-50 rounded-lg text-emerald-600"><i class="icofont-building-alt text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-gray-500">Plan baseline: UGX 31,015,914,535 across 297 assets</div>
        </div>

        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Adjusted Valuation</p>
              <h3 class="text-xl font-extrabold text-blue-900 mt-1 font-mono">UGX {{ formatUGX(summary.total_adjusted_cost) }}</h3>
            </div>
            <div class="p-3 bg-blue-50 rounded-lg text-blue-600"><i class="icofont-chart-growth text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-blue-600 font-medium">Includes statutory revaluation adjustments</div>
        </div>

        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Net Book Value (NBV)</p>
              <h3 class="text-xl font-extrabold text-purple-900 mt-1 font-mono">UGX {{ formatUGX(summary.total_net_book_value) }}</h3>
            </div>
            <div class="p-3 bg-purple-50 rounded-lg text-purple-600"><i class="icofont-calculator-alt-2 text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-purple-600 font-medium">Net after IPSAS 17 accumulated depreciation</div>
        </div>

        <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Audit Verification</p>
              <h3 class="text-xl font-extrabold text-amber-900 mt-1">{{ summary.verified_assets || 0 }} / {{ summary.total_assets || 0 }}</h3>
            </div>
            <div class="p-3 bg-amber-50 rounded-lg text-amber-600"><i class="icofont-check-circled text-2xl"></i></div>
          </div>
          <div class="mt-3 text-xs text-amber-700 font-medium flex items-center gap-1">
            <span class="inline-block w-2 h-2 rounded-full bg-emerald-500"></span>
            {{ summary.discrepancy_assets || 0 }} discrepancies flagged
          </div>
        </div>
      </div>

      <!-- Navigation Tabs -->
      <div class="border-b border-gray-200 bg-white rounded-t-xl px-4">
        <nav class="-mb-px flex flex-wrap gap-x-8" aria-label="Tabs">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              activeTab === tab.id ? 'border-emerald-600 text-emerald-700 font-bold' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 font-medium',
              'whitespace-nowrap py-4 px-1 border-b-2 text-sm flex items-center gap-2 transition'
            ]"
          >
            <span>{{ tab.name }}</span>
            <span v-if="tab.count !== undefined" :class="activeTab === tab.id ? 'bg-emerald-100 text-emerald-800' : 'bg-gray-100 text-gray-600'" class="ml-1 py-0.5 px-2 rounded-full text-xs font-bold">{{ tab.count }}</span>
          </button>
        </nav>
      </div>

      <!-- Tab 1: Asset Register Table -->
      <div v-if="activeTab === 'register'" class="bg-white rounded-b-xl border border-gray-100 shadow-sm overflow-hidden">
        <div class="p-4 border-b border-gray-100 bg-gray-50/50 flex flex-col md:flex-row items-center justify-between gap-4">
          <div class="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
            <div class="relative w-full md:w-72">
              <input v-model="searchQuery" @input="debouncedFetchAssets" type="text" placeholder="Search Tag, Code, Name..." aria-label="Search fixed assets" class="asset-search w-full pl-9 pr-4 py-2 text-xs border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500" />
              <i class="icofont-search-1 text-gray-400 absolute left-3 top-2.5"></i>
            </div>

            <select aria-label="Filter asset category" v-model="selectedCategory" @change="offset = 0; fetchAssets()" class="text-xs border border-gray-300 rounded-lg py-2 px-3 focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 bg-white">
              <option value="">All Categories</option>
              <option v-for="cat in categoriesList" :key="cat" :value="cat">{{ cat }}</option>
            </select>
          </div>

          <div class="text-xs text-gray-500 font-medium">
            Showing {{ assets.length }} of {{ totalAssets }} records
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-gray-100/70 text-gray-600 font-semibold uppercase tracking-wider border-b border-gray-200">
              <tr>
                <th class="py-3 px-4">Asset Code</th>
                <th class="py-3 px-4">Tag Number</th>
                <th class="py-3 px-4">Asset Description</th>
                <th class="py-3 px-4">Class / Subcategory</th>
                <th class="py-3 px-4 text-right">Units</th>
                <th class="py-3 px-4 text-right">FB Cost (UGX)</th>
                <th class="py-3 px-4 text-right">Adjusted Valuation</th>
                <th class="py-3 px-4 text-right">Net Book Value</th>
                <th class="py-3 px-4 text-center">Verification</th>
                <th class="asset-actions py-3 px-4 text-center">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-if="loading">
                <td colspan="10" class="py-8 text-center text-gray-400">Loading fixed asset register...</td>
              </tr>
              <tr v-else-if="!assets.length">
                <td colspan="10" class="py-8 text-center text-gray-400">No fixed assets found.</td>
              </tr>
              <template v-else>
                <tr v-for="asset in assets" :key="asset.id" class="hover:bg-emerald-50/40 transition">
                  <td class="py-3 px-4 font-mono font-bold text-emerald-800">{{ asset.asset_number }}</td>
                  <td class="py-3 px-4 font-mono text-gray-700"><span class="bg-gray-100 px-2 py-0.5 rounded border border-gray-200">{{ asset.tag_number }}</span></td>
                  <td class="py-3 px-4 text-gray-900 font-medium max-w-xs truncate" :title="asset.asset_description">{{ asset.asset_description }}</td>
                  <td class="py-3 px-4 text-gray-600"><div class="font-semibold text-gray-800">{{ asset.category_segment3 }}</div><div class="text-[10px] text-gray-400">{{ asset.category_segment4 }}</div></td>
                  <td class="py-3 px-4 text-right font-mono">{{ asset.asset_units }}</td>
                  <td class="py-3 px-4 text-right font-mono text-gray-700">{{ formatUGX(asset.fb_cost) }}</td>
                  <td class="py-3 px-4 text-right font-mono font-semibold text-blue-800">{{ formatUGX(asset.adjusted_cost) }}</td>
                  <td class="py-3 px-4 text-right font-mono font-bold text-purple-900">{{ formatUGX(asset.net_book_value) }}</td>
                  <td class="py-3 px-4 text-center"><span :class="verificationClass(asset.verification_status)" class="px-2 py-1 rounded-full text-[10px] font-bold">{{ asset.verification_status || 'UNVERIFIED' }}</span></td>
                  <td class="asset-actions py-3 px-4 text-center space-x-1">
                    <router-link :to="{ name: 'FixedAssetRevalue', params: { id: asset.id } }" title="Revalue Asset" class="px-2 py-1 bg-blue-50 text-blue-700 hover:bg-blue-100 rounded text-[11px] font-semibold">Adjust</router-link>
                    <router-link :to="{ name: 'FixedAssetVerify', params: { id: asset.id } }" title="Spot-check verify tag" class="px-2 py-1 bg-amber-50 text-amber-700 hover:bg-amber-100 rounded text-[11px] font-semibold">Verify</router-link>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </div>

      <div v-if="activeTab === 'register'" class="flex flex-wrap items-center justify-between gap-3 bg-white border border-gray-200 rounded-lg p-4" aria-label="Asset register pagination">
        <span class="text-sm text-gray-700">{{ totalAssets ? offset + 1 : 0 }}–{{ Math.min(offset + assets.length, totalAssets) }} of {{ totalAssets }} assets</span>
        <div class="flex gap-3">
          <button :disabled="loading || offset === 0" @click="offset = Math.max(0, offset - pageSize); fetchAssets()" class="px-4 py-2 border border-gray-300 rounded-lg text-sm font-semibold disabled:opacity-40">Previous page</button>
          <button :disabled="loading || offset + pageSize >= totalAssets" @click="offset += pageSize; fetchAssets()" class="px-4 py-2 border border-emerald-700 text-emerald-800 rounded-lg text-sm font-semibold disabled:opacity-40">Next page</button>
        </div>
      </div>

      <!-- Tab 2: Revaluations Audit Log -->
      <div v-if="activeTab === 'adjustments'" class="bg-white rounded-b-xl border border-gray-100 p-6 shadow-sm">
        <h3 class="text-base font-bold text-gray-900 mb-2">Asset Revaluation & Valuation Adjustment Logs</h3>
        <p class="text-xs text-gray-500 mb-4">Statutory audit trail capturing all price overrides (`FB_COST` to `ADJUSTED COST`) in compliance with Treasury Instructions 2017 & PFMA 2015.</p>
        <div class="border border-blue-100 bg-blue-50/50 p-4 rounded-xl mb-4 flex items-start space-x-3">
          <i class="icofont-info-circle text-blue-600 mt-0.5"></i>
          <div class="text-xs text-blue-900">
            <p class="font-bold">Revaluation Rule Notice:</p>
            <p>Indoor Stadium (`M1007059` / `166BLNG5`) carries FB_COST zero when Excel displays `#############`, while adjusted valuation records <span class="font-mono font-bold">UGX 1,160,000,000.00</span>.</p>
          </div>
        </div>
      </div>

      <!-- Tab 3: Dynamic Pivot & Analytics Engine -->
      <div v-if="activeTab === 'pivot'" class="bg-white rounded-b-xl border border-gray-100 p-6 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-base font-bold text-gray-900">Dynamic Multi-Dimensional Pivot Engine</h3>
            <p class="text-xs text-gray-500">Live replacement for Excel worksheet `PIVOT TABLE` cross-tabulating categories, sub-classes, units, and valuation.</p>
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
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="(item, idx) in summary.category_summaries" :key="idx" class="hover:bg-gray-50">
                <td class="py-2.5 px-4 font-semibold text-gray-800">{{ item.category_segment1 }}</td>
                <td class="py-2.5 px-4 font-medium text-emerald-800">{{ item.category_segment3 }}</td>
                <td class="py-2.5 px-4 text-gray-600">{{ item.category_segment4 }}</td>
                <td class="py-2.5 px-4 text-right font-mono font-bold">{{ item.asset_units }}</td>
                <td class="py-2.5 px-4 text-right font-mono">{{ formatUGX(item.total_fb_cost) }}</td>
                <td class="py-2.5 px-4 text-right font-mono font-bold text-blue-900">{{ formatUGX(item.total_adjusted_cost) }}</td>
                <td class="py-2.5 px-4 text-right font-mono font-bold text-purple-900">{{ formatUGX(item.total_net_book_value) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

    </div>
  </LayoutDefault>
</template>

<script>
import client from '@/api/client'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'


function apiData(res) {
  return res?.data?.data ?? res?.data ?? {}
}

export default {
  name: 'FixedAssetsView',
  components: { LayoutDefault },
  data() {
    return {
      assets: [],
      totalAssets: 0,
      offset: 0,
      pageSize: 100,
      loading: false,
      fetchTimer: null,
      noticeTimer: null,
      notice: { type: '', message: '' },
      summary: {
        total_assets: null,
        total_fb_cost: null,
        total_adjusted_cost: null,
        total_accumulated_deprec: 0,
        total_net_book_value: null,
        verified_assets: 0,
        discrepancy_assets: 0,
        category_summaries: [],
      },
      activeTab: 'register',
      searchQuery: '',
      selectedCategory: '',
      categoriesList: [
        'CYCLES',
        'ELECTRICAL MACHINERY',
        'FURNITURE AND FITTINGS',
        'LAND',
        'LIGHT ICT HARDWARE',
        'LIGHT VEHICLES',
        'NON RESIDENTIAL BUILDINGS',
        'OFFICE EQUIPMENT',
        'OTHER ICT EQUIPMENT',
        'RESIDENTIAL BUILDINGS',
      ],
    }
  },
  computed: {
    tabs() {
      return [
        { id: 'register', name: 'Asset Register', count: this.totalAssets || this.summary.total_assets },
        { id: 'adjustments', name: 'Value Adjustments' },
        { id: 'pivot', name: 'Dynamic Pivot Engine', count: this.summary.category_summaries?.length || 0 },
      ]
    },
  },
  mounted() {
    const messages = { new: 'Fixed asset created successfully', revalue: 'Asset revaluation saved', verify: 'Asset verification recorded', depreciation: 'Depreciation run completed' }
    if (messages[this.$route.query.completed]) this.setNotice('success', messages[this.$route.query.completed])
    this.fetchAssets()
    this.fetchSummary()
  },
  beforeUnmount() {
    if (this.fetchTimer) clearTimeout(this.fetchTimer)
    if (this.noticeTimer) clearTimeout(this.noticeTimer)
  },
  methods: {
    async fetchAssets() {
      this.loading = true
      try {
        const params = { limit: this.pageSize, offset: this.offset, search: this.searchQuery, category: this.selectedCategory }
        const res = await client.get('/api/v1/assets', { params })
        const data = apiData(res)
        this.assets = data.assets || []
        this.totalAssets = data.total || 0
      } catch (e) {
        this.setNotice('error', e.response?.data?.error || e.response?.data?.message || 'Failed to load fixed assets')
      } finally {
        this.loading = false
      }
    },
    debouncedFetchAssets() {
      this.offset = 0
      if (this.fetchTimer) clearTimeout(this.fetchTimer)
      this.fetchTimer = setTimeout(() => this.fetchAssets(), 250)
    },
    async fetchSummary() {
      try {
        const res = await client.get('/api/v1/assets/summary')
        const data = apiData(res)
        if (data) this.summary = { ...this.summary, ...data }
      } catch (e) {
        this.setNotice('error', e.response?.data?.error || e.response?.data?.message || 'Failed to load asset summary')
      }
    },
    formatUGX(val) {
      if (val === null || val === undefined) return '—'
      return new Intl.NumberFormat('en-UG', { maximumFractionDigits: 2 }).format(val || 0)
    },
    verificationClass(status) {
      switch ((status || '').toUpperCase()) {
        case 'VERIFIED': return 'bg-emerald-100 text-emerald-700'
        case 'DISCREPANCY': return 'bg-red-100 text-red-700'
        case 'MISSING': return 'bg-orange-100 text-orange-700'
        default: return 'bg-gray-100 text-gray-600'
      }
    },
    setNotice(type, message) {
      this.notice = { type, message }
      if (this.noticeTimer) clearTimeout(this.noticeTimer)
      this.noticeTimer = setTimeout(() => { this.notice = { type: '', message: '' } }, 5000)
    },

  },
}
</script>

<style scoped>
.asset-search {
  padding-left: 2.25rem !important;
}
.asset-actions {
  position: sticky;
  right: 0;
  z-index: 1;
  background: #fff;
  min-width: 8rem;
  box-shadow: -2px 0 4px rgb(0 0 0 / 8%);
}
.asset-actions a {
  display: inline-flex;
  align-items: center;
  min-height: 2rem;
}

</style>
