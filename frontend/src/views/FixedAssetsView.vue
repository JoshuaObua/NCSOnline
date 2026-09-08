<template>
  <LayoutDefault title="Fixed Assets Registry - Manage Assets">
    <div class="space-y-6">
      <FixedAssetsPageHeader
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

      <!-- Asset Register Master Table Card (Otika DataTables Format) -->
      <div class="card otika-datatable-card">
        <div class="card-header flex flex-col sm:flex-row sm:items-center justify-between gap-3 py-3.5 px-5">
          <h4 class="text-base font-bold text-gray-900 flex items-center gap-2">
            <i class="icofont-building-alt text-primary"></i> Fixed Asset Master Register (IPSAS 17)
          </h4>
          <div class="card-header-action flex items-center gap-2">
            <router-link to="/fixed-assets/new" class="btn btn-sm btn-icon icon-left btn-primary">
              <i class="icofont-plus"></i> Add Asset
            </router-link>
            <button @click="exportCSV" class="btn btn-sm btn-icon icon-left btn-success">
              <i class="icofont-file-excel mr-1"></i> Export CSV
            </button>
            <button @click="printTable" class="btn btn-sm btn-icon icon-left btn-secondary">
              <i class="icofont-printer mr-1"></i> Print
            </button>
          </div>
        </div>

        <div class="card-body">
          <div class="table-responsive">
            <div id="table-1_wrapper" class="dataTables_wrapper dt-bootstrap4 no-footer">
              
              <!-- DataTables Controls Header Row -->
              <div class="row align-items-center mb-3 otika-datatable-toolbar">
                <div class="col-sm-12 col-md-6 d-flex align-items-center flex-wrap gap-2 mb-2 mb-md-0">
                  <div class="dataTables_length" id="table-1_length">
                    <label class="d-flex align-items-center gap-1 mb-0 text-xs">
                      Show 
                      <select v-model="pageSize" @change="offset = 0; fetchAssets()" class="form-control form-control-sm custom-select custom-select-sm otika-form-control w-auto mx-1">
                        <option :value="10">10</option>
                        <option :value="25">25</option>
                        <option :value="50">50</option>
                        <option :value="100">100</option>
                      </select> 
                      entries
                    </label>
                  </div>

                  <select
                    v-model="selectedCategory"
                    @change="offset = 0; fetchAssets()"
                    class="form-control form-control-sm otika-form-control w-auto text-xs"
                    aria-label="Filter category"
                  >
                    <option value="">All Categories</option>
                    <option v-for="cat in categoriesList" :key="cat" :value="cat">{{ cat }}</option>
                  </select>

                  <button @click="searchQuery = ''; selectedCategory = ''; offset = 0; fetchAssets()" class="btn btn-sm btn-icon btn-outline-secondary" title="Reset filters" aria-label="Reset filters">
                    <i class="icofont-refresh"></i>
                  </button>
                </div>

                <div class="col-sm-12 col-md-6 d-flex justify-content-md-end">
                  <div id="table-1_filter" class="dataTables_filter text-md-right w-100 w-md-auto">
                    <label class="d-flex align-items-center justify-content-md-end gap-2 mb-0 text-xs font-semibold">
                      Search:
                      <div class="relative w-full sm:w-64">
                        <input
                          type="search"
                          v-model="searchQuery"
                          @input="debouncedFetchAssets"
                          class="form-control form-control-sm otika-form-control pl-8 pr-3 text-xs"
                          placeholder="Tag, Asset Code, Name..."
                          aria-label="Search assets"
                        />
                        <i class="icofont-search-1 text-gray-400 absolute left-2.5 top-2.5"></i>
                      </div>
                    </label>
                  </div>
                </div>
              </div>

              <!-- Main Data Table -->
              <table class="table table-striped table-hover dataTable no-footer w-100 text-xs" id="table-1">
                <thead>
                  <tr>
                    <th class="text-center w-10">
                      <label class="otika-checkbox" title="Select all visible assets">
                        <input type="checkbox" :checked="allVisibleSelected" @change="toggleAllVisible" aria-label="Select all visible assets" />
                        <span></span>
                      </label>
                    </th>
                    <th class="text-center w-12">#</th>
                    <th>Asset Code</th>
                    <th>Tag Number</th>
                    <th>Asset Description</th>
                    <th>Class / Subcategory</th>
                    <th class="text-right">Units</th>
                    <th class="text-right">FB Cost (UGX)</th>
                    <th class="text-right">Adjusted Valuation</th>
                    <th class="text-right">Net Book Value</th>
                    <th class="text-center">Verification</th>
                    <th class="text-center w-28">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-if="loading">
                    <td colspan="12" class="text-center py-8 text-gray-400">
                      <i class="icofont-spinner-alt-4 animate-spin text-lg mr-2"></i> Loading fixed asset register...
                    </td>
                  </tr>
                  <tr v-else-if="!assets.length">
                    <td colspan="12" class="text-center py-8 text-gray-400">No fixed assets found matching filter criteria.</td>
                  </tr>
                  <template v-else>
                    <tr v-for="(asset, idx) in assets" :key="asset.id" class="hover:bg-slate-50 transition">
                      <td class="text-center">
                        <label class="otika-checkbox">
                          <input v-model="selectedAssetIds" type="checkbox" :value="asset.id" :aria-label="`Select ${asset.asset_number}`" />
                          <span></span>
                        </label>
                      </td>
                      <td class="text-center font-bold text-gray-500">{{ offset + idx + 1 }}</td>
                      <td class="font-mono font-bold text-primary">{{ asset.asset_number }}</td>
                      <td>
                        <span class="badge badge-light font-mono">{{ asset.tag_number }}</span>
                      </td>
                      <td class="font-medium max-w-xs truncate" :title="asset.asset_description">
                        {{ asset.asset_description }}
                      </td>
                      <td>
                        <div class="font-semibold text-gray-800">{{ asset.category_segment3 }}</div>
                        <div class="text-[10px] text-gray-400">{{ asset.category_segment4 }}</div>
                      </td>
                      <td class="text-right font-mono">{{ asset.asset_units }}</td>
                      <td class="text-right font-mono text-gray-700">{{ formatUGX(asset.fb_cost) }}</td>
                      <td class="text-right font-mono font-semibold text-primary">{{ formatUGX(asset.adjusted_cost) }}</td>
                      <td class="text-right font-mono font-bold text-purple-900">{{ formatUGX(asset.net_book_value) }}</td>
                      <td class="text-center">
                        <span :class="verificationBadgeClass(asset.verification_status)" class="badge badge-shadow">
                          {{ asset.verification_status || 'UNVERIFIED' }}
                        </span>
                      </td>
                      <td class="text-center">
                        <div class="btn-group otika-action-buttons">
                          <router-link
                            :to="{ name: 'FixedAssetRevalue', params: { id: asset.id } }"
                            title="Adjust Asset Cost"
                            class="btn btn-icon btn-sm btn-primary"
                            aria-label="Adjust asset cost"
                          >
                            <i class="icofont-edit"></i>
                          </router-link>
                          <router-link
                            :to="{ name: 'FixedAssetVerify', params: { id: asset.id } }"
                            title="Verify Tag"
                            class="btn btn-icon btn-sm btn-info"
                            aria-label="Verify asset tag"
                          >
                            <i class="icofont-check-circled"></i>
                          </router-link>
                          <button
                            @click="openDeleteModal(asset)"
                            title="Delete Asset"
                            class="btn btn-icon btn-sm btn-danger"
                            aria-label="Delete asset"
                          >
                            <i class="icofont-trash"></i>
                          </button>
                        </div>
                      </td>
                    </tr>
                  </template>
                </tbody>
              </table>

              <!-- DataTables Bottom Controls Footer Row -->
              <div class="row align-items-center mt-3 pt-3 border-t border-gray-100">
                <div class="col-sm-12 col-md-5 mb-2 mb-md-0">
                  <div class="dataTables_info" id="table-1_info" role="status" aria-live="polite">
                    Showing {{ totalAssets ? offset + 1 : 0 }} to {{ Math.min(offset + assets.length, totalAssets) }} of {{ totalAssets }} entries
                    <span v-if="selectedAssetIds.length"> · {{ selectedAssetIds.length }} selected</span>
                  </div>
                </div>
                <div class="col-sm-12 col-md-7 d-flex justify-content-md-end">
                  <div class="dataTables_paginate paging_simple_numbers" id="table-1_paginate">
                    <ul class="pagination pagination-sm mb-0">
                      <li class="paginate_button page-item previous" :class="{ disabled: loading || offset === 0 }">
                        <button
                          :disabled="loading || offset === 0"
                          @click="offset = Math.max(0, offset - pageSize); fetchAssets()"
                          class="page-link"
                        >
                          Previous
                        </button>
                      </li>
                      <li class="paginate_button page-item active">
                        <span class="page-link">{{ Math.floor(offset / pageSize) + 1 }}</span>
                      </li>
                      <li class="paginate_button page-item next" :class="{ disabled: loading || offset + pageSize >= totalAssets }">
                        <button
                          :disabled="loading || offset + pageSize >= totalAssets"
                          @click="offset += pageSize; fetchAssets()"
                          class="page-link"
                        >
                          Next
                        </button>
                      </li>
                    </ul>
                  </div>
                </div>
              </div>

            </div>
          </div>
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
import FixedAssetsPageHeader from '@/components/fixed_assets/FixedAssetsPageHeader.vue'

function apiData(res) {
  return res?.data?.data ?? res?.data ?? {}
}

export default {
  name: 'FixedAssetsView',
  components: { LayoutDefault, FixedAssetsPageHeader },
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
      selectedAssetIds: [],
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
  computed: {
    allVisibleSelected() {
      return this.assets.length > 0 && this.assets.every(asset => this.selectedAssetIds.includes(asset.id))
    },
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
    toggleAllVisible(event) {
      const visibleIds = this.assets.map(asset => asset.id)
      if (event.target.checked) {
        this.selectedAssetIds = [...new Set([...this.selectedAssetIds, ...visibleIds])]
      } else {
        this.selectedAssetIds = this.selectedAssetIds.filter(id => !visibleIds.includes(id))
      }
    },
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
    verificationBadgeClass(status) {
      switch ((status || '').toUpperCase()) {
        case 'VERIFIED': return 'badge-success'
        case 'DISCREPANCY': return 'badge-danger'
        case 'MISSING': return 'badge-warning'
        default: return 'badge-light'
      }
    },
    exportCSV() {
      if (!this.assets || !this.assets.length) return
      const headers = ['Asset Code', 'Tag Number', 'Description', 'Category', 'Units', 'FB Cost', 'Adjusted Valuation', 'Net Book Value', 'Verification']
      const rows = this.assets.map(a => [
        `"${a.asset_number || ''}"`,
        `"${a.tag_number || ''}"`,
        `"${(a.asset_description || '').replace(/"/g, '""')}"`,
        `"${a.category_segment3 || ''}"`,
        a.asset_units || 1,
        a.fb_cost || 0,
        a.adjusted_cost || 0,
        a.net_book_value || 0,
        `"${a.verification_status || 'UNVERIFIED'}"`
      ])
      const csvContent = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map(e => e.join(','))].join('\n')
      const encodedUri = encodeURI(csvContent)
      const link = document.createElement('a')
      link.setAttribute('href', encodedUri)
      link.setAttribute('download', `Fixed_Assets_Register_${new Date().toISOString().slice(0,10)}.csv`)
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    },
    printTable() {
      window.print()
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
