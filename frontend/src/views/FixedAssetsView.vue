<template>
  <LayoutDefault title="Fixed Assets Registry - Manage Assets">
    <div class="space-y-6">
      <!-- Shared Header & Sub-Navigation -->
      <FixedAssetsSubNav
        title="Manage Fixed Assets"
        :badge="summary.total_assets ? summary.total_assets + ' Registered Assets' : ''"
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
          <div class="mt-3 text-xs text-gray-500">Plan baseline: UGX 31,015,914,535 baseline</div>
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

      <!-- Asset Register Master Table Card -->
      <div class="bg-white rounded-xl border border-gray-100 shadow-sm overflow-hidden">
        <div class="p-4 border-b border-gray-100 bg-gray-50/50 flex flex-col md:flex-row items-center justify-between gap-4">
          <div class="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
            <div class="relative w-full md:w-80">
              <input
                v-model="searchQuery"
                @input="debouncedFetchAssets"
                type="text"
                placeholder="Search Tag, Asset Code, Description..."
                aria-label="Search fixed assets"
                class="asset-search w-full pl-9 pr-4 py-2 text-xs border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
              />
              <i class="icofont-search-1 text-gray-400 absolute left-3 top-2.5"></i>
            </div>

            <select
              aria-label="Filter asset category"
              v-model="selectedCategory"
              @change="offset = 0; fetchAssets()"
              class="text-xs border border-gray-300 rounded-lg py-2 px-3 focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 bg-white"
            >
              <option value="">All Categories</option>
              <option v-for="cat in categoriesList" :key="cat" :value="cat">{{ cat }}</option>
            </select>
          </div>

          <div class="text-xs text-gray-500 font-medium flex items-center gap-3">
            <span>Showing {{ assets.length }} of {{ totalAssets }} records</span>
            <button @click="searchQuery = ''; selectedCategory = ''; offset = 0; fetchAssets()" class="text-emerald-700 hover:underline">Reset</button>
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
                <td colspan="10" class="py-8 text-center text-gray-400">No fixed assets found matching filter criteria.</td>
              </tr>
              <template v-else>
                <tr v-for="asset in assets" :key="asset.id" class="hover:bg-emerald-50/40 transition">
                  <td class="py-3 px-4 font-mono font-bold text-emerald-800">{{ asset.asset_number }}</td>
                  <td class="py-3 px-4 font-mono text-gray-700">
                    <span class="bg-gray-100 px-2 py-0.5 rounded border border-gray-200">{{ asset.tag_number }}</span>
                  </td>
                  <td class="py-3 px-4 text-gray-900 font-medium max-w-xs truncate" :title="asset.asset_description">
                    {{ asset.asset_description }}
                  </td>
                  <td class="py-3 px-4 text-gray-600">
                    <div class="font-semibold text-gray-800">{{ asset.category_segment3 }}</div>
                    <div class="text-[10px] text-gray-400">{{ asset.category_segment4 }}</div>
                  </td>
                  <td class="py-3 px-4 text-right font-mono">{{ asset.asset_units }}</td>
                  <td class="py-3 px-4 text-right font-mono text-gray-700">{{ formatUGX(asset.fb_cost) }}</td>
                  <td class="py-3 px-4 text-right font-mono font-semibold text-blue-800">{{ formatUGX(asset.adjusted_cost) }}</td>
                  <td class="py-3 px-4 text-right font-mono font-bold text-purple-900">{{ formatUGX(asset.net_book_value) }}</td>
                  <td class="py-3 px-4 text-center">
                    <span :class="verificationClass(asset.verification_status)" class="px-2 py-1 rounded-full text-[10px] font-bold">
                      {{ asset.verification_status || 'UNVERIFIED' }}
                    </span>
                  </td>
                  <td class="asset-actions py-3 px-4 text-center flex items-center justify-center gap-1">
                    <router-link
                      :to="{ name: 'FixedAssetRevalue', params: { id: asset.id } }"
                      title="Revalue / Adjust Asset Cost"
                      class="px-2 py-1 bg-blue-50 text-blue-700 hover:bg-blue-100 rounded text-[11px] font-semibold transition"
                    >
                      Adjust
                    </router-link>
                    <router-link
                      :to="{ name: 'FixedAssetVerify', params: { id: asset.id } }"
                      title="Verify Tag"
                      class="px-2 py-1 bg-amber-50 text-amber-700 hover:bg-amber-100 rounded text-[11px] font-semibold transition"
                    >
                      Verify
                    </router-link>
                    <button
                      @click="openDeleteModal(asset)"
                      title="Delete Asset & Adjust Portfolio Valuation"
                      class="px-2 py-1 bg-red-50 text-red-700 hover:bg-red-100 rounded text-[11px] font-semibold transition"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Pagination -->
      <div class="flex flex-wrap items-center justify-between gap-3 bg-white border border-gray-200 rounded-xl p-4 shadow-sm" aria-label="Asset register pagination">
        <span class="text-sm text-gray-700">
          Showing {{ totalAssets ? offset + 1 : 0 }}–{{ Math.min(offset + assets.length, totalAssets) }} of {{ totalAssets }} assets
        </span>
        <div class="flex gap-3">
          <button
            :disabled="loading || offset === 0"
            @click="offset = Math.max(0, offset - pageSize); fetchAssets()"
            class="px-4 py-2 border border-gray-300 rounded-lg text-xs font-semibold hover:bg-gray-50 disabled:opacity-40"
          >
            Previous page
          </button>
          <button
            :disabled="loading || offset + pageSize >= totalAssets"
            @click="offset += pageSize; fetchAssets()"
            class="px-4 py-2 border border-emerald-700 text-emerald-800 rounded-lg text-xs font-semibold hover:bg-emerald-50 disabled:opacity-40"
          >
            Next page
          </button>
        </div>
      </div>

      <!-- Delete Asset Confirmation Modal -->
      <div v-if="deletingAsset" class="fixed inset-0 z-50 flex items-center justify-center bg-gray-900/60 p-4">
        <div class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl space-y-4 border border-gray-100 animate-fade-in">
          <div class="flex items-start justify-between border-b border-gray-100 pb-3">
            <div class="flex items-center gap-3 text-red-600">
              <div class="w-10 h-10 rounded-full bg-red-100 flex items-center justify-center flex-shrink-0">
                <i class="icofont-trash text-xl"></i>
              </div>
              <div>
                <h3 class="text-base font-bold text-gray-900">Delete Fixed Asset</h3>
                <p class="text-xs text-gray-500">Action cannot be undone. System portfolio totals will be adjusted.</p>
              </div>
            </div>
            <button @click="deletingAsset = null" class="text-gray-400 hover:text-gray-600">
              <i class="icofont-close text-lg"></i>
            </button>
          </div>

          <!-- Asset Summary Box -->
          <div class="bg-red-50/60 rounded-xl p-4 border border-red-100 space-y-2 text-xs">
            <div class="flex justify-between">
              <span class="font-bold text-gray-700">Asset Number:</span>
              <span class="font-mono text-gray-900 font-bold">{{ deletingAsset.asset_number }}</span>
            </div>
            <div class="flex justify-between">
              <span class="font-bold text-gray-700">Tag Number:</span>
              <span class="font-mono text-gray-900 font-bold">{{ deletingAsset.tag_number }}</span>
            </div>
            <div class="flex justify-between">
              <span class="font-bold text-gray-700">Description:</span>
              <span class="text-gray-900 font-medium max-w-xs text-right truncate">{{ deletingAsset.asset_description }}</span>
            </div>
            <div class="flex justify-between border-t border-red-200/60 pt-2 text-red-900">
              <span class="font-bold">Valuation Deduction (Rollback):</span>
              <span class="font-mono font-extrabold text-sm">UGX {{ formatUGX(deletingAsset.adjusted_cost) }}</span>
            </div>
          </div>

          <!-- Deletion Reason Field -->
          <div class="space-y-1">
            <label class="block text-xs font-bold text-gray-700">
              Deletion Reason / Notes <span class="text-red-500">*</span>
            </label>
            <textarea
              v-model="deleteNotes"
              rows="3"
              placeholder="Provide statutory justification or reference for deleting this record..."
              class="w-full text-xs p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-red-500"
            ></textarea>
          </div>

          <!-- Action Buttons -->
          <div class="flex items-center justify-end gap-3 pt-2">
            <button
              @click="deletingAsset = null"
              class="px-4 py-2 rounded-lg border border-gray-300 text-xs font-semibold text-gray-700 hover:bg-gray-50"
            >
              Cancel
            </button>

            <button
              :disabled="deleteSubmitting || !deleteNotes.trim()"
              @click="confirmDeleteAsset"
              class="px-5 py-2 rounded-lg bg-red-600 hover:bg-red-700 text-white text-xs font-bold shadow-sm transition disabled:opacity-40 flex items-center gap-1.5"
            >
              <i v-if="deleteSubmitting" class="icofont-spinner spin"></i>
              <span>Confirm & Delete Asset</span>
            </button>
          </div>
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
  name: 'FixedAssetsView',
  components: { LayoutDefault, FixedAssetsSubNav },
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
      searchQuery: '',
      selectedCategory: '',
      deletingAsset: null,
      deleteNotes: '',
      deleteSubmitting: false,
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
  mounted() {
    const messages = {
      new: 'Fixed asset created successfully',
      revalue: 'Asset revaluation saved',
      verify: 'Asset verification recorded',
      depreciation: 'Depreciation run completed',
    }
    if (messages[this.$route.query.completed]) {
      this.setNotice('success', messages[this.$route.query.completed])
    }
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
    openDeleteModal(asset) {
      this.deletingAsset = asset
      this.deleteNotes = ''
    },
    async confirmDeleteAsset() {
      if (!this.deletingAsset || !this.deleteNotes.trim()) return
      this.deleteSubmitting = true
      try {
        const assetTag = this.deletingAsset.tag_number || this.deletingAsset.asset_number
        const assetCost = this.deletingAsset.adjusted_cost
        await client.delete(`/api/v1/assets/${this.deletingAsset.id}`, {
          data: { notes: this.deleteNotes.trim() }
        })
        this.setNotice('success', `Asset "${assetTag}" deleted successfully. Portfolio valuation adjusted by -UGX ${this.formatUGX(assetCost)}.`)
        this.deletingAsset = null
        this.deleteNotes = ''
        await this.fetchAssets()
        await this.fetchSummary()
      } catch (e) {
        this.setNotice('error', e.response?.data?.error || e.response?.data?.message || 'Failed to delete fixed asset')
      } finally {
        this.deleteSubmitting = false
      }
    },
    setNotice(type, message) {
      this.notice = { type, message }
      if (this.noticeTimer) clearTimeout(this.noticeTimer)
      this.noticeTimer = setTimeout(() => { this.notice = { type: '', message: '' } }, 7000)
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
  min-width: 10rem;
  box-shadow: -2px 0 4px rgb(0 0 0 / 8%);
}
.spin {
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
