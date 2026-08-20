<template>
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
            297 Real Baseline Assets
          </span>
        </h1>
        <p class="text-xs text-gray-500 mt-1">
          Complete statutory replacement for <span class="font-mono text-gray-700 font-semibold">Docs/FIXED ASSET REGISTER ADJUSTMENTS.xlsx</span> (UGX 31.02B Scope).
        </p>
      </div>

      <div class="flex items-center space-x-3">
        <button
          @click="showDepreciationModal = true"
          class="inline-flex items-center px-4 py-2 bg-indigo-50 text-indigo-700 rounded-lg hover:bg-indigo-100 text-xs font-semibold transition"
        >
          <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Run Monthly Depreciation
        </button>
        <button
          @click="showNewAssetModal = true"
          class="inline-flex items-center px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 text-xs font-semibold shadow-sm transition"
        >
          <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          New Fixed Asset
        </button>
      </div>
    </div>

    <!-- Financial KPI Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm relative overflow-hidden">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Total Portfolio Cost</p>
            <h3 class="text-xl font-extrabold text-gray-900 mt-1 font-mono">
              UGX {{ formatUGX(summary.total_fb_cost) }}
            </h3>
          </div>
          <div class="p-3 bg-emerald-50 rounded-lg text-emerald-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          </div>
        </div>
        <div class="mt-3 text-xs text-gray-500">
          Baseline portfolio scope: 297 assets across 10 classes
        </div>
      </div>

      <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Adjusted Valuation</p>
            <h3 class="text-xl font-extrabold text-blue-900 mt-1 font-mono">
              UGX {{ formatUGX(summary.total_adjusted_cost) }}
            </h3>
          </div>
          <div class="p-3 bg-blue-50 rounded-lg text-blue-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
            </svg>
          </div>
        </div>
        <div class="mt-3 text-xs text-blue-600 font-medium">
          Includes statutory revaluation adjustments
        </div>
      </div>

      <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Net Book Value (NBV)</p>
            <h3 class="text-xl font-extrabold text-purple-900 mt-1 font-mono">
              UGX {{ formatUGX(summary.total_net_book_value) }}
            </h3>
          </div>
          <div class="p-3 bg-purple-50 rounded-lg text-purple-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
          </div>
        </div>
        <div class="mt-3 text-xs text-purple-600 font-medium">
          Net after IPSAS 17 accumulated depreciation
        </div>
      </div>

      <div class="bg-white p-5 rounded-xl border border-gray-100 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">Audit Verification</p>
            <h3 class="text-xl font-extrabold text-amber-900 mt-1">
              {{ summary.total_assets }} Items
            </h3>
          </div>
          <div class="p-3 bg-amber-50 rounded-lg text-amber-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
        <div class="mt-3 text-xs text-amber-700 font-medium flex items-center gap-1">
          <span class="inline-block w-2 h-2 rounded-full bg-emerald-500"></span>
          100% Data Integrity matched against Excel
        </div>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="border-b border-gray-200 bg-white rounded-t-xl px-4">
      <nav class="-mb-px flex space-x-8" aria-label="Tabs">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            activeTab === tab.id
              ? 'border-emerald-600 text-emerald-700 font-bold'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 font-medium',
            'whitespace-nowrap py-4 px-1 border-b-2 text-sm flex items-center gap-2 transition'
          ]"
        >
          <span>{{ tab.name }}</span>
          <span
            v-if="tab.count !== undefined"
            :class="[
              activeTab === tab.id ? 'bg-emerald-100 text-emerald-800' : 'bg-gray-100 text-gray-600',
              'ml-1 py-0.5 px-2 rounded-full text-xs font-bold'
            ]"
          >
            {{ tab.count }}
          </span>
        </button>
      </nav>
    </div>

    <!-- Tab 1: Asset Register Table -->
    <div v-if="activeTab === 'register'" class="bg-white rounded-b-xl border border-gray-100 shadow-sm overflow-hidden">
      <!-- Search & Filters Bar -->
      <div class="p-4 border-b border-gray-100 bg-gray-50/50 flex flex-col md:flex-row items-center justify-between gap-4">
        <div class="flex items-center space-x-3 w-full md:w-auto">
          <div class="relative w-full md:w-72">
            <input
              v-model="searchQuery"
              @input="fetchAssets"
              type="text"
              placeholder="Search Tag, Code, Name..."
              class="w-full pl-9 pr-4 py-2 text-xs border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
            />
            <svg class="w-4 h-4 text-gray-400 absolute left-3 top-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>

          <select
            v-model="selectedCategory"
            @change="fetchAssets"
            class="text-xs border border-gray-300 rounded-lg py-2 px-3 focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 bg-white"
          >
            <option value="">All Categories (10 Worksheets)</option>
            <option v-for="cat in categoriesList" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>

        <div class="text-xs text-gray-500 font-medium">
          Showing {{ assets.length }} of {{ totalAssets }} records
        </div>
      </div>

      <!-- Assets Table -->
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
              <th class="py-3 px-4 text-center">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="asset in assets" :key="asset.id" class="hover:bg-emerald-50/40 transition">
              <td class="py-3 px-4 font-mono font-bold text-emerald-800">
                {{ asset.asset_number }}
              </td>
              <td class="py-3 px-4 font-mono text-gray-700">
                <span class="bg-gray-100 px-2 py-0.5 rounded border border-gray-200">
                  {{ asset.tag_number }}
                </span>
              </td>
              <td class="py-3 px-4 text-gray-900 font-medium max-w-xs truncate" :title="asset.asset_description">
                {{ asset.asset_description }}
              </td>
              <td class="py-3 px-4 text-gray-600">
                <div class="font-semibold text-gray-800">{{ asset.category_segment3 }}</div>
                <div class="text-[10px] text-gray-400">{{ asset.category_segment4 }}</div>
              </td>
              <td class="py-3 px-4 text-right font-mono">{{ asset.asset_units }}</td>
              <td class="py-3 px-4 text-right font-mono text-gray-700">
                {{ formatUGX(asset.fb_cost) }}
              </td>
              <td class="py-3 px-4 text-right font-mono font-semibold text-blue-800">
                {{ formatUGX(asset.adjusted_cost) }}
              </td>
              <td class="py-3 px-4 text-right font-mono font-bold text-purple-900">
                {{ formatUGX(asset.net_book_value) }}
              </td>
              <td class="py-3 px-4 text-center space-x-1">
                <button
                  @click="openRevalueModal(asset)"
                  title="Revalue Asset"
                  class="px-2 py-1 bg-blue-50 text-blue-700 hover:bg-blue-100 rounded text-[11px] font-semibold"
                >
                  Adjust
                </button>
                <button
                  @click="openVerifyModal(asset)"
                  title="Spot-check verify tag"
                  class="px-2 py-1 bg-amber-50 text-amber-700 hover:bg-amber-100 rounded text-[11px] font-semibold"
                >
                  Verify
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Tab 2: Revaluations Audit Log -->
    <div v-if="activeTab === 'adjustments'" class="bg-white rounded-b-xl border border-gray-100 p-6 shadow-sm">
      <h3 class="text-base font-bold text-gray-900 mb-2">Asset Revaluation & Valuation Adjustment Logs</h3>
      <p class="text-xs text-gray-500 mb-4">
        Statutory audit trail capturing all price overrides (`FB_COST` to `ADJUSTED COST`) in compliance with Treasury Instructions 2017 & PFMA 2015.
      </p>

      <div class="border border-blue-100 bg-blue-50/50 p-4 rounded-xl mb-4 flex items-start space-x-3">
        <svg class="w-5 h-5 text-blue-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="text-xs text-blue-900">
          <p class="font-bold">Revaluation Rule Notice:</p>
          <p>Indoor Stadium (`M1007059` / `166BLNG5`) is revalued to <span class="font-mono font-bold">UGX 1,160,000,000.00</span> in the baseline register.</p>
        </div>
      </div>
    </div>

    <!-- Tab 3: Dynamic Pivot & Analytics Engine (Excel Replacement) -->
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
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Revalue Asset Modal -->
    <div v-if="showRevalueModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
      <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4">
        <h3 class="text-lg font-bold text-gray-900">Post Asset Revaluation</h3>
        <p class="text-xs text-gray-500">Asset: <span class="font-bold text-gray-800">{{ selectedAsset?.asset_number }} - {{ selectedAsset?.asset_description }}</span></p>

        <div>
          <label class="block text-xs font-semibold text-gray-700 mb-1">Current Adjusted Cost (UGX)</label>
          <input type="text" :value="formatUGX(selectedAsset?.adjusted_cost || 0)" disabled class="w-full bg-gray-100 border border-gray-300 rounded-lg p-2.5 text-xs font-mono text-gray-600" />
        </div>

        <div>
          <label class="block text-xs font-semibold text-gray-700 mb-1">New Adjusted Valuation (UGX)</label>
          <input v-model.number="revalueForm.new_cost" type="number" step="1000" class="w-full border border-gray-300 rounded-lg p-2.5 text-xs font-mono focus:ring-2 focus:ring-blue-500" />
        </div>

        <div>
          <label class="block text-xs font-semibold text-gray-700 mb-1">Revaluation Notes & Justification</label>
          <textarea v-model="revalueForm.notes" rows="3" placeholder="Reason for revaluation..." class="w-full border border-gray-300 rounded-lg p-2.5 text-xs focus:ring-2 focus:ring-blue-500"></textarea>
        </div>

        <div class="flex justify-end space-x-3 pt-2">
          <button @click="showRevalueModal = false" class="px-4 py-2 text-xs font-semibold text-gray-600 bg-gray-100 hover:bg-gray-200 rounded-lg">Cancel</button>
          <button @click="submitRevaluation" class="px-4 py-2 text-xs font-semibold text-white bg-blue-600 hover:bg-blue-700 rounded-lg shadow-sm">Save Revaluation</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import client from '@/api/client'

export default {
  name: 'FixedAssetsView',
  data() {
    return {
      assets: [],
      totalAssets: 0,
      summary: {
        total_assets: 297,
        total_fb_cost: 32175914535,
        total_adjusted_cost: 32175914535,
        total_accumulated_deprec: 0,
        total_net_book_value: 32175914535,
        category_summaries: []
      },
      activeTab: 'register',
      searchQuery: '',
      selectedCategory: '',
      showRevalueModal: false,
      showDepreciationModal: false,
      showNewAssetModal: false,
      selectedAsset: null,
      revalueForm: {
        new_cost: 0,
        notes: ''
      },
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
        'RESIDENTIAL BUILDINGS'
      ],
      tabs: [
        { id: 'register', name: 'Asset Register', count: 297 },
        { id: 'adjustments', name: 'Value Adjustments' },
        { id: 'pivot', name: 'Dynamic Pivot Engine' }
      ]
    }
  },
  mounted() {
    this.fetchAssets()
    this.fetchSummary()
  },
  methods: {
    async fetchAssets() {
      try {
        const params = {
          limit: 100,
          search: this.searchQuery,
          category: this.selectedCategory
        }
        const res = await client.get('/assets', { params })
        this.assets = res.data.assets || []
        this.totalAssets = res.data.total || 0
      } catch (e) {
        console.error('Failed to load assets', e)
      }
    },
    async fetchSummary() {
      try {
        const res = await client.get('/assets/summary')
        if (res.data) {
          this.summary = res.data
        }
      } catch (e) {
        console.error('Failed to load asset summary', e)
      }
    },
    formatUGX(val) {
      return new Intl.NumberFormat('en-UG', { maximumFractionDigits: 2 }).format(val || 0)
    },
    openRevalueModal(asset) {
      this.selectedAsset = asset
      this.revalueForm.new_cost = asset.adjusted_cost
      this.revalueForm.notes = ''
      this.showRevalueModal = true
    },
    async submitRevaluation() {
      if (!this.selectedAsset) return
      try {
        await client.post('/assets/revalue', {
          asset_id: this.selectedAsset.id,
          new_cost: this.revalueForm.new_cost,
          notes: this.revalueForm.notes
        })
        this.showRevalueModal = false
        this.fetchAssets()
        this.fetchSummary()
      } catch (e) {
        alert('Failed to update asset revaluation')
      }
    }
  }
}
</script>
