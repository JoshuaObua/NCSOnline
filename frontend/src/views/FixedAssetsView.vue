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
          <h1 class="text-2xl font-bold text-gray-900 tracking-tight flex items-center gap-2">
            Fixed Assets & IPSAS 17 Register
            <span class="text-xs bg-emerald-100 text-emerald-800 font-semibold px-2.5 py-0.5 rounded-full border border-emerald-200">
              {{ summary.total_assets || 297 }} Real Baseline Assets
            </span>
          </h1>
          <p class="text-xs text-gray-500 mt-1">
            Complete statutory replacement for <span class="font-mono text-gray-700 font-semibold">Docs/FIXED ASSET REGISTER ADJUSTMENTS.xlsx</span> (UGX 31.02B FB_COST baseline).
          </p>
        </div>

        <div class="flex flex-wrap items-center gap-3">
          <button
            @click="showDepreciationModal = true"
            class="inline-flex items-center px-4 py-2 bg-indigo-50 text-indigo-700 rounded-lg hover:bg-indigo-100 text-xs font-semibold transition"
          >
            <i class="icofont-clock-time mr-1.5"></i>
            Run Monthly Depreciation
          </button>
          <button
            @click="openNewAssetModal"
            class="inline-flex items-center px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 text-xs font-semibold shadow-sm transition"
          >
            <i class="icofont-plus-circle mr-1.5"></i>
            New Fixed Asset
          </button>
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
              <input v-model="searchQuery" @input="debouncedFetchAssets" type="text" placeholder="Search Tag, Code, Name..." class="w-full pl-9 pr-4 py-2 text-xs border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500" />
              <i class="icofont-search-1 text-gray-400 absolute left-3 top-2.5"></i>
            </div>

            <select v-model="selectedCategory" @change="fetchAssets" class="text-xs border border-gray-300 rounded-lg py-2 px-3 focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 bg-white">
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
                <th class="py-3 px-4 text-center">Actions</th>
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
                  <td class="py-3 px-4 text-center space-x-1">
                    <button @click="openRevalueModal(asset)" title="Revalue Asset" class="px-2 py-1 bg-blue-50 text-blue-700 hover:bg-blue-100 rounded text-[11px] font-semibold">Adjust</button>
                    <button @click="openVerifyModal(asset)" title="Spot-check verify tag" class="px-2 py-1 bg-amber-50 text-amber-700 hover:bg-amber-100 rounded text-[11px] font-semibold">Verify</button>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
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

      <!-- New Asset Modal -->
      <div v-if="showNewAssetModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
        <div class="bg-white rounded-2xl max-w-3xl w-full p-6 shadow-2xl space-y-4 max-h-[92vh] overflow-y-auto">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-bold text-gray-900">Register New Fixed Asset</h3>
            <button @click="showNewAssetModal = false" class="text-gray-400 hover:text-gray-700">&times;</button>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <label class="text-xs font-semibold text-gray-700">Asset Number<input v-model="newAsset.asset_number" class="form-input" placeholder="M1009999" /></label>
            <label class="text-xs font-semibold text-gray-700">Tag Number<input v-model="newAsset.tag_number" class="form-input" placeholder="NCS-TAG-001" /></label>
            <label class="md:col-span-2 text-xs font-semibold text-gray-700">Description<input v-model="newAsset.asset_description" class="form-input" placeholder="Asset description" /></label>
            <label class="text-xs font-semibold text-gray-700">Category Segment 1<input v-model="newAsset.category_segment1" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Category Segment 3<select v-model="newAsset.category_segment3" class="form-input"><option v-for="cat in categoriesList" :key="cat" :value="cat">{{ cat }}</option></select></label>
            <label class="text-xs font-semibold text-gray-700">Category Segment 4<input v-model="newAsset.category_segment4" class="form-input" placeholder="Subclass" /></label>
            <label class="text-xs font-semibold text-gray-700">Units<input v-model.number="newAsset.asset_units" type="number" min="1" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">FB Cost<input v-model.number="newAsset.fb_cost" type="number" min="0" step="1000" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Adjusted Cost<input v-model.number="newAsset.adjusted_cost" type="number" min="0" step="1000" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Date In Service<input v-model="newAsset.date_placed_in_service" type="date" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Department<input v-model="newAsset.custodian_department" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Location<input v-model="newAsset.location_building" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Useful Life Years<input v-model.number="newAsset.useful_life_years" type="number" min="0" class="form-input" /></label>
          </div>
          <div class="flex justify-end gap-3 pt-2">
            <button @click="showNewAssetModal = false" class="px-4 py-2 text-xs font-semibold text-gray-600 bg-gray-100 hover:bg-gray-200 rounded-lg">Cancel</button>
            <button @click="submitNewAsset" class="px-4 py-2 text-xs font-semibold text-white bg-emerald-600 hover:bg-emerald-700 rounded-lg shadow-sm">Create Asset</button>
          </div>
        </div>
      </div>

      <!-- Depreciation Modal -->
      <div v-if="showDepreciationModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
        <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4">
          <h3 class="text-lg font-bold text-gray-900">Run Monthly Depreciation</h3>
          <p class="text-xs text-gray-500">Runs IPSAS 17 straight-line depreciation for active assets with remaining net book value.</p>
          <label class="block text-xs font-semibold text-gray-700">Period<input v-model="depreciationForm.period" type="month" class="form-input" /></label>
          <div class="flex justify-end gap-3 pt-2">
            <button @click="showDepreciationModal = false" class="px-4 py-2 text-xs font-semibold text-gray-600 bg-gray-100 hover:bg-gray-200 rounded-lg">Cancel</button>
            <button @click="submitDepreciation" class="px-4 py-2 text-xs font-semibold text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg shadow-sm">Run Depreciation</button>
          </div>
        </div>
      </div>

      <!-- Verify Asset Modal -->
      <div v-if="showVerifyModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
        <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4">
          <h3 class="text-lg font-bold text-gray-900">Physical Tag Verification</h3>
          <p class="text-xs text-gray-500">Asset: <span class="font-bold text-gray-800">{{ selectedAsset?.asset_number }} - {{ selectedAsset?.tag_number }}</span></p>
          <label class="block text-xs font-semibold text-gray-700">Verification Status<select v-model="verifyForm.status" class="form-input"><option value="VERIFIED">VERIFIED</option><option value="DISCREPANCY">DISCREPANCY</option><option value="MISSING">MISSING</option></select></label>
          <label class="block text-xs font-semibold text-gray-700">Notes<textarea v-model="verifyForm.notes" rows="3" class="form-input" placeholder="Spot-check notes"></textarea></label>
          <div class="flex justify-end gap-3 pt-2">
            <button @click="showVerifyModal = false" class="px-4 py-2 text-xs font-semibold text-gray-600 bg-gray-100 hover:bg-gray-200 rounded-lg">Cancel</button>
            <button @click="submitVerification" class="px-4 py-2 text-xs font-semibold text-white bg-amber-600 hover:bg-amber-700 rounded-lg shadow-sm">Save Verification</button>
          </div>
        </div>
      </div>

      <!-- Revalue Asset Modal -->
      <div v-if="showRevalueModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
        <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4">
          <h3 class="text-lg font-bold text-gray-900">Post Asset Revaluation</h3>
          <p class="text-xs text-gray-500">Asset: <span class="font-bold text-gray-800">{{ selectedAsset?.asset_number }} - {{ selectedAsset?.asset_description }}</span></p>
          <label class="block text-xs font-semibold text-gray-700">Current Adjusted Cost (UGX)<input type="text" :value="formatUGX(selectedAsset?.adjusted_cost || 0)" disabled class="form-input bg-gray-100 text-gray-600" /></label>
          <label class="block text-xs font-semibold text-gray-700">New Adjusted Valuation (UGX)<input v-model.number="revalueForm.new_cost" type="number" step="1000" class="form-input" /></label>
          <label class="block text-xs font-semibold text-gray-700">Revaluation Notes & Justification<textarea v-model="revalueForm.notes" rows="3" placeholder="Reason for revaluation..." class="form-input"></textarea></label>
          <div class="flex justify-end gap-3 pt-2">
            <button @click="showRevalueModal = false" class="px-4 py-2 text-xs font-semibold text-gray-600 bg-gray-100 hover:bg-gray-200 rounded-lg">Cancel</button>
            <button @click="submitRevaluation" class="px-4 py-2 text-xs font-semibold text-white bg-blue-600 hover:bg-blue-700 rounded-lg shadow-sm">Save Revaluation</button>
          </div>
        </div>
      </div>
    </div>
  </LayoutDefault>
</template>

<script>
import client from '@/api/client'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'

const PLAN_FB_COST_TOTAL = 31015914535
const PLAN_ADJUSTED_TOTAL = 32175914535
const PLAN_ASSET_COUNT = 297

function defaultNewAsset() {
  return {
    asset_book: 'NCS FA BOOK',
    asset_number: '',
    tag_number: '',
    asset_description: '',
    category_segment1: 'MACHINERY AND EQUIPMENT',
    category_segment3: 'LIGHT ICT HARDWARE',
    category_segment4: '',
    asset_units: 1,
    fb_cost: 0,
    adjusted_cost: 0,
    date_placed_in_service: '2023-07-01',
    custodian_department: 'General Administration',
    location_building: 'NCS Lugogo Head Office',
    location_room: 'Main Facility',
    depreciation_method: 'STRAIGHT_LINE',
    useful_life_years: 5,
    worksheet_source: 'LIGHT ICT HARDWARE',
  }
}

export default {
  name: 'FixedAssetsView',
  components: { LayoutDefault },
  data() {
    return {
      assets: [],
      totalAssets: 0,
      loading: false,
      fetchTimer: null,
      noticeTimer: null,
      notice: { type: '', message: '' },
      summary: {
        total_assets: PLAN_ASSET_COUNT,
        total_fb_cost: PLAN_FB_COST_TOTAL,
        total_adjusted_cost: PLAN_ADJUSTED_TOTAL,
        total_accumulated_deprec: 0,
        total_net_book_value: PLAN_ADJUSTED_TOTAL,
        verified_assets: 0,
        discrepancy_assets: 0,
        category_summaries: [],
      },
      activeTab: 'register',
      searchQuery: '',
      selectedCategory: '',
      showRevalueModal: false,
      showDepreciationModal: false,
      showNewAssetModal: false,
      showVerifyModal: false,
      selectedAsset: null,
      revalueForm: { new_cost: 0, notes: '' },
      verifyForm: { status: 'VERIFIED', notes: '' },
      depreciationForm: { period: new Date().toISOString().slice(0, 7) },
      newAsset: defaultNewAsset(),
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
    this.fetchAssets()
    this.fetchSummary()
  },
  beforeUnmount() {
    if (this.fetchTimer) clearTimeout(this.fetchTimer)
  },
  methods: {
    async fetchAssets() {
      this.loading = true
      try {
        const params = { limit: 100, search: this.searchQuery, category: this.selectedCategory }
        const res = await client.get('/api/v1/assets', { params })
        this.assets = res.data.assets || []
        this.totalAssets = res.data.total || 0
      } catch (e) {
        this.setNotice('error', e.response?.data?.error || e.response?.data?.message || 'Failed to load fixed assets')
      } finally {
        this.loading = false
      }
    },
    debouncedFetchAssets() {
      if (this.fetchTimer) clearTimeout(this.fetchTimer)
      this.fetchTimer = setTimeout(() => this.fetchAssets(), 250)
    },
    async fetchSummary() {
      try {
        const res = await client.get('/api/v1/assets/summary')
        if (res.data) this.summary = { ...this.summary, ...res.data }
      } catch (e) {
        this.setNotice('error', e.response?.data?.error || e.response?.data?.message || 'Failed to load asset summary')
      }
    },
    formatUGX(val) {
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
    openNewAssetModal() {
      this.newAsset = defaultNewAsset()
      this.showNewAssetModal = true
    },
    openRevalueModal(asset) {
      this.selectedAsset = asset
      this.revalueForm.new_cost = asset.adjusted_cost
      this.revalueForm.notes = ''
      this.showRevalueModal = true
    },
    openVerifyModal(asset) {
      this.selectedAsset = asset
      this.verifyForm.status = asset.verification_status && asset.verification_status !== 'UNVERIFIED' ? asset.verification_status : 'VERIFIED'
      this.verifyForm.notes = ''
      this.showVerifyModal = true
    },
    async submitNewAsset() {
      try {
        const payload = { ...this.newAsset }
        if (!payload.adjusted_cost && payload.fb_cost) payload.adjusted_cost = payload.fb_cost
        payload.worksheet_source = payload.category_segment3
        await client.post('/api/v1/assets', payload)
        this.showNewAssetModal = false
        this.setNotice('success', 'Fixed asset created successfully')
        await Promise.all([this.fetchAssets(), this.fetchSummary()])
      } catch (e) {
        this.setNotice('error', e.response?.data?.error || e.response?.data?.message || 'Failed to create fixed asset')
      }
    },
    async submitRevaluation() {
      if (!this.selectedAsset) return
      try {
        await client.post('/api/v1/assets/revalue', {
          asset_id: this.selectedAsset.id,
          new_cost: this.revalueForm.new_cost,
          notes: this.revalueForm.notes,
        })
        this.showRevalueModal = false
        this.setNotice('success', 'Asset revaluation saved')
        await Promise.all([this.fetchAssets(), this.fetchSummary()])
      } catch (e) {
        this.setNotice('error', e.response?.data?.error || e.response?.data?.message || 'Failed to update asset revaluation')
      }
    },
    async submitVerification() {
      if (!this.selectedAsset) return
      try {
        await client.post('/api/v1/assets/verify', {
          asset_id: this.selectedAsset.id,
          status: this.verifyForm.status,
          notes: this.verifyForm.notes,
        })
        this.showVerifyModal = false
        this.setNotice('success', 'Asset verification recorded')
        await Promise.all([this.fetchAssets(), this.fetchSummary()])
      } catch (e) {
        this.setNotice('error', e.response?.data?.error || e.response?.data?.message || 'Failed to verify fixed asset')
      }
    },
    async submitDepreciation() {
      try {
        const res = await client.post('/api/v1/assets/depreciate', { period: this.depreciationForm.period })
        this.showDepreciationModal = false
        this.setNotice('success', `${res.data.assets_processed || 0} assets depreciated for ${this.depreciationForm.period}`)
        await Promise.all([this.fetchAssets(), this.fetchSummary()])
      } catch (e) {
        this.setNotice('error', e.response?.data?.error || e.response?.data?.message || 'Failed to run depreciation')
      }
    },
  },
}
</script>

<style scoped>
.form-input {
  width: 100%;
  margin-top: 0.35rem;
  border: 1px solid #d1d5db;
  border-radius: 0.5rem;
  padding: 0.625rem;
  font-size: 0.75rem;
}
.form-input:focus {
  outline: none;
  border-color: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.18);
}
</style>
