<template>
  <LayoutDefault title="Fixed Assets Registry - Manage Assets">
    <div class="space-y-6">
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
            <i class="icofont-building-alt text-primary"></i> Fixed Asset Master Register
          </h4>
          <div class="card-header-action flex flex-wrap items-center gap-2">
            <button @click="openImportModal" class="btn btn-sm btn-icon icon-left btn-info">
              <i class="icofont-upload"></i> Bulk Import
            </button>
            <button @click="downloadImportTemplate" class="btn btn-sm btn-icon icon-left btn-outline-info">
              <i class="icofont-file-spreadsheet"></i> Template
            </button>
            <router-link to="/fixed-assets/new" class="btn btn-sm btn-icon icon-left btn-primary">
              <i class="icofont-plus"></i> Add Asset
            </router-link>
            <button @click="exportAssets('excel')" class="btn btn-sm btn-icon icon-left btn-success">
              <i class="icofont-file-excel"></i> Excel
            </button>
            <button @click="exportAssets('pdf')" class="btn btn-sm btn-icon icon-left btn-danger">
              <i class="icofont-file-pdf"></i> PDF
            </button>
            <button @click="printTable" class="btn btn-sm btn-icon icon-left btn-secondary">
              <i class="icofont-printer"></i> Print
            </button>
          </div>
        </div>

        <div class="card-body">
          <div class="table-responsive">
            <div id="table-1_wrapper" class="dataTables_wrapper dt-bootstrap4 no-footer">
              
              <!-- DataTables Controls Header Row -->
              <div class="otika-datatable-toolbar mb-3">
                <div class="otika-toolbar-group">
                  <div class="dataTables_length otika-toolbar-control" id="table-1_length">
                    <label class="mb-0 text-xs font-semibold text-slate-600">
                      Show
                      <select v-model="pageSize" @change="offset = 0; fetchAssets()" class="form-control form-control-sm custom-select custom-select-sm otika-form-control mx-1">
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
                    class="form-control form-control-sm otika-form-control otika-category-filter text-xs"
                    aria-label="Filter category"
                  >
                    <option value="">All Categories</option>
                    <option v-for="cat in categoriesList" :key="cat" :value="cat">{{ cat }}</option>
                  </select>

                  <button @click="searchQuery = ''; selectedCategory = ''; offset = 0; fetchAssets()" class="btn btn-sm btn-icon btn-outline-secondary" title="Reset filters" aria-label="Reset filters">
                    <i class="icofont-refresh"></i>
                  </button>
                </div>

                <div class="otika-toolbar-group otika-toolbar-right">
                  <div v-if="selectedAssetIds.length" class="otika-bulk-actions">
                    <span class="otika-selection-count">{{ selectedAssetIds.length }} selected</span>
                    <button @click="openBulkDeleteModal" class="btn btn-sm btn-icon icon-left btn-danger">
                      <i class="icofont-trash"></i> Bulk Delete
                    </button>
                  </div>
                  <div id="table-1_filter" class="dataTables_filter">
                    <label class="otika-search-label mb-0 text-xs font-semibold">
                      Search:
                      <div class="relative otika-search-box">
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

      <!-- Bulk Delete Confirmation Modal -->
      <div v-if="bulkDeleteOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-gray-900/60 p-4">
        <div class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl space-y-4 border border-gray-100 animate-fade-in">
          <div class="flex items-start justify-between border-b border-gray-100 pb-3">
            <div class="flex items-center gap-3 text-red-600">
              <div class="w-10 h-10 rounded-full bg-red-100 flex items-center justify-center flex-shrink-0">
                <i class="icofont-trash text-xl"></i>
              </div>
              <div>
                <h3 class="text-base font-bold text-gray-900">Bulk Delete Fixed Assets</h3>
                <p class="text-xs text-gray-500">{{ selectedAssetIds.length }} selected assets will be deleted and logged.</p>
              </div>
            </div>
            <button @click="bulkDeleteOpen = false" class="text-gray-400 hover:text-gray-600">
              <i class="icofont-close text-lg"></i>
            </button>
          </div>

          <div class="space-y-1">
            <label class="block text-xs font-bold text-gray-700">
              Bulk Deletion Reason / Notes <span class="text-red-500">*</span>
            </label>
            <textarea
              v-model="bulkDeleteNotes"
              rows="3"
              placeholder="Provide statutory justification for this bulk deletion..."
              class="w-full text-xs p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-red-500"
            ></textarea>
          </div>

          <div class="flex items-center justify-end gap-3 pt-2">
            <button @click="bulkDeleteOpen = false" class="px-4 py-2 rounded-lg border border-gray-300 text-xs font-semibold text-gray-700 hover:bg-gray-50">
              Cancel
            </button>
            <button
              :disabled="bulkDeleteSubmitting || !bulkDeleteNotes.trim()"
              @click="confirmBulkDelete"
              class="px-5 py-2 rounded-lg bg-red-600 hover:bg-red-700 text-white text-xs font-bold shadow-sm transition disabled:opacity-40 flex items-center gap-1.5"
            >
              <i v-if="bulkDeleteSubmitting" class="icofont-spinner spin"></i>
              <span>Delete Selected Assets</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Bulk Import Modal -->
      <div v-if="importOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-gray-900/60 p-4">
        <div class="bg-white rounded-2xl max-w-2xl w-full p-6 shadow-2xl space-y-5 border border-gray-100 animate-fade-in">
          <div class="flex items-start justify-between border-b border-gray-100 pb-3">
            <div>
              <h3 class="text-base font-bold text-gray-900">Bulk Import Fixed Assets</h3>
              <p class="text-xs text-gray-500 mt-1">Use the template, then upload CSV, TSV, or Excel-readable XLS.</p>
            </div>
            <button @click="closeImportModal" class="text-gray-400 hover:text-gray-600">
              <i class="icofont-close text-lg"></i>
            </button>
          </div>

          <label
            class="dropify-dropzone"
            :class="{ 'is-dragover': importDragOver, 'has-file': importFile }"
            @dragover.prevent="importDragOver = true"
            @dragleave.prevent="importDragOver = false"
            @drop.prevent="handleImportDrop"
          >
            <input
              ref="importInput"
              type="file"
              accept=".csv,.tsv,.xls,text/csv,text/tab-separated-values,application/vnd.ms-excel"
              class="hidden"
              @change="handleImportFile"
            />
            <span class="dropify-icon"><i class="icofont-cloud-upload"></i></span>
            <span class="dropify-message">{{ importFile ? importFile.name : 'Drop fixed asset import file here or click to browse' }}</span>
            <span class="dropify-hint">Template columns: Asset Number, Tag Number, Asset Description, Category, Subcategory, Units, FB Cost, Adjusted Cost</span>
          </label>

          <div v-if="importPreview.length" class="border border-gray-200 rounded-xl overflow-hidden">
            <div class="px-4 py-2 bg-gray-50 text-xs font-bold text-gray-700">Preview first {{ importPreview.length }} rows</div>
            <div class="overflow-x-auto max-h-48">
              <table class="table table-sm mb-0 text-xs">
                <thead>
                  <tr>
                    <th>Asset Code</th>
                    <th>Tag Number</th>
                    <th>Description</th>
                    <th>Category</th>
                    <th class="text-right">FB Cost</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(row, idx) in importPreview" :key="idx">
                    <td class="font-mono">{{ row.asset_number }}</td>
                    <td class="font-mono">{{ row.tag_number }}</td>
                    <td>{{ row.asset_description }}</td>
                    <td>{{ row.category_segment3 }}</td>
                    <td class="text-right font-mono">{{ formatUGX(row.fb_cost) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <div class="flex flex-wrap items-center justify-between gap-3">
            <button @click="downloadImportTemplate" class="btn btn-sm btn-icon icon-left btn-outline-info">
              <i class="icofont-file-spreadsheet"></i> Download Template
            </button>
            <div class="flex items-center gap-2">
              <button @click="closeImportModal" class="btn btn-sm btn-outline-secondary">Cancel</button>
              <button
                :disabled="importSubmitting || !importRows.length"
                @click="confirmImportAssets"
                class="btn btn-sm btn-icon icon-left btn-success disabled:opacity-40"
              >
                <i v-if="importSubmitting" class="icofont-spinner spin"></i>
                <i v-else class="icofont-check-circled"></i>
                Import {{ importRows.length || '' }} Assets
              </button>
            </div>
          </div>
        </div>
      </div>

    </div>
  </LayoutDefault>
</template>

<script>
import client from '@/api/client'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { exportToExcel, exportToPDF } from '@/utils/reportExporter.js'

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
      searchQuery: '',
      selectedCategory: '',
      selectedAssetIds: [],
      deletingAsset: null,
      deleteNotes: '',
      deleteSubmitting: false,
      bulkDeleteOpen: false,
      bulkDeleteNotes: '',
      bulkDeleteSubmitting: false,
      importOpen: false,
      importDragOver: false,
      importFile: null,
      importRows: [],
      importSubmitting: false,
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
    importPreview() {
      return this.importRows.slice(0, 5)
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
    async fetchExportAssets() {
      const params = {
        limit: Math.max(this.totalAssets || 10000, 10000),
        offset: 0,
        search: this.searchQuery,
        category: this.selectedCategory,
      }
      const res = await client.get('/api/v1/assets', { params })
      const data = apiData(res)
      return data.assets || []
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
    exportHeaders() {
      return [
        { key: 'asset_number', label: 'Asset Code' },
        { key: 'tag_number', label: 'Tag Number' },
        { key: 'asset_description', label: 'Description' },
        { key: 'category_segment1', label: 'Primary Class' },
        { key: 'category_segment3', label: 'Category' },
        { key: 'category_segment4', label: 'Subcategory' },
        { key: 'asset_units', label: 'Units' },
        { key: 'fb_cost', label: 'FB Cost' },
        { key: 'adjusted_cost', label: 'Adjusted Valuation' },
        { key: 'net_book_value', label: 'Net Book Value' },
        { key: 'verification_status', label: 'Verification' },
        { key: 'custodian_department', label: 'Custodian Department' },
        { key: 'location_building', label: 'Building' },
        { key: 'location_room', label: 'Room' },
      ]
    },
    async exportAssets(format) {
      try {
        const rows = await this.fetchExportAssets()
        if (!rows.length) {
          this.setNotice('error', 'No fixed assets available to export')
          return
        }
        const preparedRows = rows.map(asset => ({
          ...asset,
          fb_cost: this.formatUGX(asset.fb_cost),
          adjusted_cost: this.formatUGX(asset.adjusted_cost),
          net_book_value: this.formatUGX(asset.net_book_value),
          verification_status: asset.verification_status || 'UNVERIFIED',
        }))
        const stamp = new Date().toISOString().slice(0, 10)
        if (format === 'pdf') {
          exportToPDF('Fixed Asset Master Register', this.exportHeaders(), preparedRows, {
            status: this.selectedCategory || 'All Categories',
            department: 'Financial & Assets Ledger',
          })
        } else {
          exportToExcel(`Fixed_Asset_Master_Register_${stamp}`, this.exportHeaders(), preparedRows, 'Fixed Asset Master Register')
        }
      } catch (e) {
        this.setNotice('error', e.response?.data?.error || e.response?.data?.message || 'Failed to export fixed assets')
      }
    },
    printTable() {
      window.print()
    },
    downloadImportTemplate() {
      const headers = [
        { key: 'asset_number', label: 'Asset Number' },
        { key: 'tag_number', label: 'Tag Number' },
        { key: 'asset_description', label: 'Asset Description' },
        { key: 'category_segment1', label: 'Primary Class' },
        { key: 'category_segment3', label: 'Category' },
        { key: 'category_segment4', label: 'Subcategory' },
        { key: 'asset_units', label: 'Units' },
        { key: 'fb_cost', label: 'FB Cost' },
        { key: 'adjusted_cost', label: 'Adjusted Cost' },
        { key: 'date_placed_in_service', label: 'Date Placed In Service' },
        { key: 'custodian_department', label: 'Custodian Department' },
        { key: 'location_building', label: 'Location Building' },
        { key: 'location_room', label: 'Location Room' },
        { key: 'useful_life_years', label: 'Useful Life Years' },
      ]
      const rows = [{
        asset_number: 'NCS-FA-0001',
        tag_number: 'NCS/TAG/0001',
        asset_description: 'Sample office desk',
        category_segment1: 'FIXED ASSETS',
        category_segment3: 'FURNITURE AND FITTINGS',
        category_segment4: 'OFFICE FURNITURE',
        asset_units: 1,
        fb_cost: 1500000,
        adjusted_cost: 1500000,
        date_placed_in_service: new Date().toISOString().slice(0, 10),
        custodian_department: 'Finance',
        location_building: 'NCS Headquarters',
        location_room: 'Stores',
        useful_life_years: 5,
      }]
      exportToExcel('Fixed_Asset_Import_Template', headers, rows, 'Fixed Asset Import Template')
    },
    openImportModal() {
      this.importOpen = true
      this.importFile = null
      this.importRows = []
      this.importDragOver = false
    },
    closeImportModal() {
      this.importOpen = false
      this.importFile = null
      this.importRows = []
      this.importDragOver = false
      if (this.$refs.importInput) this.$refs.importInput.value = ''
    },
    handleImportDrop(event) {
      this.importDragOver = false
      const file = event.dataTransfer?.files?.[0]
      if (file) this.loadImportFile(file)
    },
    handleImportFile(event) {
      const file = event.target.files?.[0]
      if (file) this.loadImportFile(file)
    },
    async loadImportFile(file) {
      this.importFile = file
      try {
        const text = await file.text()
        this.importRows = this.parseImportRows(text)
        if (!this.importRows.length) {
          this.setNotice('error', 'No valid asset rows found in the import file')
        }
      } catch (e) {
        this.setNotice('error', 'Could not read the selected import file')
      }
    },
    parseImportRows(text) {
      const tableText = this.extractExcelTableText(text)
      const delimiter = tableText.includes('\t') ? '\t' : ','
      const rows = this.parseDelimitedText(tableText, delimiter)
      if (rows.length < 2) return []
      const headers = rows[0].map(header => this.normalizeHeader(header))
      return rows.slice(1).map(row => {
        const record = {}
        headers.forEach((key, index) => { record[key] = (row[index] || '').trim() })
        return this.normalizeImportAsset(record)
      }).filter(asset => asset.asset_number && asset.tag_number && asset.asset_description)
    },
    extractExcelTableText(text) {
      if (!/<table[\s>]/i.test(text)) return text.trim()
      const doc = new DOMParser().parseFromString(text, 'text/html')
      const rows = Array.from(doc.querySelectorAll('tr'))
      return rows.map(row => Array.from(row.querySelectorAll('th,td')).map(cell => cell.textContent.trim()).join('\t')).join('\n')
    },
    parseDelimitedText(text, delimiter = ',') {
      const rows = []
      let row = []
      let cell = ''
      let quoted = false
      for (let i = 0; i < text.length; i += 1) {
        const char = text[i]
        const next = text[i + 1]
        if (char === '"' && quoted && next === '"') {
          cell += '"'
          i += 1
        } else if (char === '"') {
          quoted = !quoted
        } else if (char === delimiter && !quoted) {
          row.push(cell)
          cell = ''
        } else if ((char === '\n' || char === '\r') && !quoted) {
          if (char === '\r' && next === '\n') i += 1
          row.push(cell)
          if (row.some(value => value.trim())) rows.push(row)
          row = []
          cell = ''
        } else {
          cell += char
        }
      }
      row.push(cell)
      if (row.some(value => value.trim())) rows.push(row)
      return rows
    },
    normalizeHeader(header) {
      const normalized = String(header || '').toLowerCase().replace(/[^a-z0-9]+/g, '_').replace(/^_|_$/g, '')
      const map = {
        asset_code: 'asset_number',
        asset_number: 'asset_number',
        tag: 'tag_number',
        tag_number: 'tag_number',
        description: 'asset_description',
        asset_description: 'asset_description',
        primary_class: 'category_segment1',
        category: 'category_segment3',
        subcategory: 'category_segment4',
        sub_class: 'category_segment4',
        units: 'asset_units',
        fb_cost: 'fb_cost',
        adjusted_cost: 'adjusted_cost',
        adjusted_valuation: 'adjusted_cost',
        date_placed_in_service: 'date_placed_in_service',
        custodian_department: 'custodian_department',
        location_building: 'location_building',
        location_room: 'location_room',
        useful_life_years: 'useful_life_years',
      }
      return map[normalized] || normalized
    },
    normalizeImportAsset(row) {
      const toNumber = value => Number(String(value || '0').replace(/,/g, '')) || 0
      return {
        asset_number: row.asset_number || '',
        tag_number: row.tag_number || '',
        asset_description: row.asset_description || '',
        category_segment1: row.category_segment1 || 'FIXED ASSETS',
        category_segment3: row.category_segment3 || '',
        category_segment4: row.category_segment4 || '',
        asset_units: Math.max(1, parseInt(row.asset_units || '1', 10) || 1),
        fb_cost: toNumber(row.fb_cost),
        adjusted_cost: toNumber(row.adjusted_cost || row.fb_cost),
        date_placed_in_service: row.date_placed_in_service || new Date().toISOString().slice(0, 10),
        custodian_department: row.custodian_department || '',
        location_building: row.location_building || '',
        location_room: row.location_room || '',
        useful_life_years: parseInt(row.useful_life_years || '5', 10) || 5,
        status: 'ACTIVE',
        verification_status: 'UNVERIFIED',
        worksheet_source: this.importFile?.name || 'Bulk Import',
      }
    },
    async confirmImportAssets() {
      if (!this.importRows.length) return
      this.importSubmitting = true
      let imported = 0
      try {
        for (const asset of this.importRows) {
          await client.post('/api/v1/assets', asset)
          imported += 1
        }
        this.setNotice('success', `${imported} fixed assets imported successfully`)
        this.closeImportModal()
        await this.fetchAssets()
        await this.fetchSummary()
      } catch (e) {
        this.setNotice('error', `Imported ${imported} rows before failure. ${e.response?.data?.error || e.response?.data?.message || 'Import failed'}`)
      } finally {
        this.importSubmitting = false
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
        const assetId = this.deletingAsset.id
        const assetTag = this.deletingAsset.tag_number || this.deletingAsset.asset_number
        const assetCost = this.deletingAsset.adjusted_cost
        await client.delete(`/api/v1/assets/${this.deletingAsset.id}`, {
          data: { notes: this.deleteNotes.trim() }
        })
        this.setNotice('success', `Asset "${assetTag}" deleted successfully. Portfolio valuation adjusted by -UGX ${this.formatUGX(assetCost)}.`)
        this.deletingAsset = null
        this.deleteNotes = ''
        this.selectedAssetIds = this.selectedAssetIds.filter(id => id !== assetId)
        await this.fetchAssets()
        await this.fetchSummary()
      } catch (e) {
        this.setNotice('error', e.response?.data?.error || e.response?.data?.message || 'Failed to delete fixed asset')
      } finally {
        this.deleteSubmitting = false
      }
    },
    openBulkDeleteModal() {
      if (!this.selectedAssetIds.length) return
      this.bulkDeleteOpen = true
      this.bulkDeleteNotes = ''
    },
    async confirmBulkDelete() {
      if (!this.selectedAssetIds.length || !this.bulkDeleteNotes.trim()) return
      this.bulkDeleteSubmitting = true
      const ids = [...this.selectedAssetIds]
      let deleted = 0
      try {
        for (const id of ids) {
          await client.delete(`/api/v1/assets/${id}`, {
            data: { notes: this.bulkDeleteNotes.trim() }
          })
          deleted += 1
        }
        this.setNotice('success', `${deleted} selected fixed assets deleted successfully`)
        this.selectedAssetIds = []
        this.bulkDeleteOpen = false
        this.bulkDeleteNotes = ''
        await this.fetchAssets()
        await this.fetchSummary()
      } catch (e) {
        this.setNotice('error', `Deleted ${deleted} selected assets before failure. ${e.response?.data?.error || e.response?.data?.message || 'Bulk delete failed'}`)
      } finally {
        this.bulkDeleteSubmitting = false
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
.otika-datatable-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #f8fafc;
}
.otika-toolbar-group {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  min-width: 0;
}
.otika-toolbar-right {
  justify-content: flex-end;
  margin-left: auto;
}
.otika-toolbar-control label,
.otika-search-label {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}
.otika-toolbar-control select {
  width: 78px;
  min-width: 78px;
}
.otika-category-filter {
  width: 220px;
  max-width: 100%;
}
.otika-search-box {
  width: 260px;
  max-width: 100%;
}
.otika-bulk-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-right: 8px;
  border-right: 1px solid #e5e7eb;
}
.otika-selection-count {
  display: inline-flex;
  align-items: center;
  min-height: 28px;
  padding: 0 10px;
  border-radius: 6px;
  background: #fff1f2;
  color: #be123c;
  font-size: 11px;
  font-weight: 800;
}
.dropify-dropzone {
  display: flex;
  min-height: 180px;
  cursor: pointer;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: 2px dashed #cbd5e1;
  border-radius: 8px;
  background: #f8fafc;
  padding: 24px;
  text-align: center;
  transition: border-color .2s, background .2s, box-shadow .2s;
}
.dropify-dropzone:hover,
.dropify-dropzone.is-dragover {
  border-color: #6777ef;
  background: #f3f4fd;
  box-shadow: 0 0 0 4px rgba(103, 119, 239, .08);
}
.dropify-dropzone.has-file {
  border-color: #10b981;
  background: #ecfdf5;
}
.dropify-icon {
  display: inline-flex;
  width: 52px;
  height: 52px;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  background: #ffffff;
  color: #6777ef;
  font-size: 28px;
  box-shadow: 0 4px 14px rgba(15, 23, 42, .08);
}
.dropify-message {
  color: #1f2937;
  font-size: 13px;
  font-weight: 800;
}
.dropify-hint {
  max-width: 32rem;
  color: #64748b;
  font-size: 11px;
  line-height: 1.5;
}
@media (max-width: 768px) {
  .otika-datatable-toolbar,
  .otika-toolbar-group,
  .otika-toolbar-right,
  .otika-search-label {
    align-items: stretch;
    width: 100%;
  }
  .otika-datatable-toolbar,
  .otika-toolbar-group,
  .otika-toolbar-right {
    flex-direction: column;
  }
  .otika-toolbar-right {
    margin-left: 0;
  }
  .otika-category-filter,
  .otika-search-box {
    width: 100%;
  }
  .otika-bulk-actions {
    width: 100%;
    justify-content: space-between;
    border-right: 0;
    padding-right: 0;
  }
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
