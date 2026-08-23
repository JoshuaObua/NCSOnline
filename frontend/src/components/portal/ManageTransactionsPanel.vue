<template>
  <div class="manage-transactions-panel">
    <!-- Header Section with Title and Actions -->
    <div class="panel-header-row">
      <div class="header-titles">
        <h2 class="panel-title">
          <i class="icofont-money-bag text-success mr-2"></i>
          System Financial Transactions
        </h2>
        <p class="panel-subtitle">
          Query, monitor, verify, and export all electronic payments and revenue collections across the portal.
        </p>
      </div>

      <div class="header-actions">
        <button
          type="button"
          class="btn btn-outline-secondary export-btn"
          :disabled="loading || !transactions.length"
          @click="exportToExcel"
          title="Download as Excel/CSV spreadsheet"
        >
          <i class="icofont-file-excel text-success"></i> Export Excel
        </button>

        <button
          type="button"
          class="btn btn-outline-secondary export-btn"
          :disabled="loading || !transactions.length"
          @click="exportToPdf"
          title="Print or Save Official PDF Ledger Report"
        >
          <i class="icofont-file-pdf text-danger"></i> Export PDF
        </button>

        <button
          type="button"
          class="btn btn-primary refresh-btn"
          :disabled="loading"
          @click="loadData"
        >
          <i class="icofont-refresh" :class="{ 'icofont-spin': loading }"></i> Refresh
        </button>
      </div>
    </div>

    <!-- KPI Summary Metrics -->
    <section class="kpi-grid">
      <div class="kpi-card kpi-total">
        <div class="kpi-icon-wrap bg-primary-subtle text-primary">
          <i class="icofont-wallet"></i>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">Total Revenue Collected</span>
          <h3 class="kpi-value text-primary">UGX {{ formatMoney(kpis.total_amount_ugx || 0) }}</h3>
          <span class="kpi-meta">{{ kpis.success_count || 0 }} successful payments</span>
        </div>
      </div>

      <div class="kpi-card kpi-success">
        <div class="kpi-icon-wrap bg-success-subtle text-success">
          <i class="icofont-check-circled"></i>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">Successful Transactions</span>
          <h3 class="kpi-value text-success">{{ kpis.success_count || 0 }}</h3>
          <span class="kpi-meta">{{ formatPercent(kpis.success_count, kpis.total_count) }} success rate</span>
        </div>
      </div>

      <div class="kpi-card kpi-pending">
        <div class="kpi-icon-wrap bg-warning-subtle text-warning">
          <i class="icofont-clock-time"></i>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">Pending / Processing</span>
          <h3 class="kpi-value text-warning">{{ kpis.pending_count || 0 }}</h3>
          <span class="kpi-meta">Awaiting provider callback</span>
        </div>
      </div>

      <div class="kpi-card kpi-failed">
        <div class="kpi-icon-wrap bg-danger-subtle text-danger">
          <i class="icofont-close-circled"></i>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">Failed / Incomplete</span>
          <h3 class="kpi-value text-danger">{{ kpis.failed_count || 0 }}</h3>
          <span class="kpi-meta">Declined or timed out</span>
        </div>
      </div>
    </section>

    <!-- Filters & Search Toolbar -->
    <section class="filter-card">
      <div class="filter-row">
        <!-- Search Input -->
        <div class="search-field">
          <i class="icofont-search"></i>
          <input
            v-model="filters.search"
            type="text"
            placeholder="Search by Reference, Payer Name, Email, Phone, or Service..."
            @input="debounceSearch"
          />
          <button
            v-if="filters.search"
            type="button"
            class="clear-btn"
            @click="filters.search = ''; loadData()"
          >
            <i class="icofont-close-line"></i>
          </button>
        </div>

        <!-- Status Filter -->
        <div class="filter-item">
          <label>Status</label>
          <select v-model="filters.status" class="form-select" @change="onFilterChange">
            <option value="">All Statuses</option>
            <option value="SUCCESS">Successful</option>
            <option value="PENDING">Pending / Processing</option>
            <option value="FAILED">Failed</option>
            <option value="CANCELLED">Cancelled</option>
          </select>
        </div>

        <!-- Payment Channel Filter -->
        <div class="filter-item">
          <label>Payment Channel</label>
          <select v-model="filters.payment_method" class="form-select" @change="onFilterChange">
            <option value="">All Channels</option>
            <option value="MOBILE_MONEY">Mobile Money (MoMo / Airtel)</option>
            <option value="AIRTEL_MONEY">Airtel Money</option>
            <option value="MTN_MOMO">MTN MoMo</option>
            <option value="BANK_TRANSFER">Bank Transfer</option>
            <option value="CARD">Visa / Mastercard</option>
          </select>
        </div>

        <!-- Start Date -->
        <div class="filter-item">
          <label>From Date</label>
          <input
            v-model="filters.start_date"
            type="date"
            class="form-control"
            @change="onFilterChange"
          />
        </div>

        <!-- End Date -->
        <div class="filter-item">
          <label>To Date</label>
          <input
            v-model="filters.end_date"
            type="date"
            class="form-control"
            @change="onFilterChange"
          />
        </div>

        <!-- Reset Button -->
        <div class="filter-item filter-reset">
          <button
            type="button"
            class="btn btn-light"
            @click="resetFilters"
            title="Clear all filters"
          >
            <i class="icofont-ui-reply"></i> Reset
          </button>
        </div>
      </div>
    </section>

    <!-- Error Alert -->
    <div v-if="error" class="alert alert-danger d-flex align-items-center justify-content-between">
      <div><i class="icofont-warning mr-2"></i> {{ error }}</div>
      <button type="button" class="btn-close" @click="error = ''">&times;</button>
    </div>

    <!-- Transactions Table -->
    <section class="table-card">
      <div class="table-responsive">
        <table class="table table-hover transactions-table mb-0">
          <thead>
            <tr>
              <th @click="toggleSort('date')" class="sortable-th">
                Date & Time
                <i :class="getSortIcon('date')"></i>
              </th>
              <th @click="toggleSort('reference')" class="sortable-th">
                Reference
                <i :class="getSortIcon('reference')"></i>
              </th>
              <th>Service / Application</th>
              <th @click="toggleSort('payer')" class="sortable-th">
                Payer Details
                <i :class="getSortIcon('payer')"></i>
              </th>
              <th @click="toggleSort('method')" class="sortable-th">
                Method / Channel
                <i :class="getSortIcon('method')"></i>
              </th>
              <th @click="toggleSort('amount')" class="sortable-th text-right">
                Amount (UGX)
                <i :class="getSortIcon('amount')"></i>
              </th>
              <th @click="toggleSort('status')" class="sortable-th text-center">
                Status
                <i :class="getSortIcon('status')"></i>
              </th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <!-- Loading Skeleton -->
            <tr v-if="loading && !transactions.length">
              <td colspan="8" class="text-center py-5">
                <div class="spinner-border text-primary" role="status">
                  <span class="sr-only">Loading transactions...</span>
                </div>
                <p class="text-muted mt-2">Loading transactions from database...</p>
              </td>
            </tr>

            <!-- Empty State -->
            <tr v-else-if="!loading && !transactions.length">
              <td colspan="8" class="text-center py-5">
                <div class="empty-state-wrap">
                  <i class="icofont-inbox fs-1 text-muted"></i>
                  <h5 class="mt-3">No Transactions Found</h5>
                  <p class="text-muted">No financial transactions match your current search and filter criteria.</p>
                  <button type="button" class="btn btn-sm btn-outline-primary mt-2" @click="resetFilters">
                    Clear Filters
                  </button>
                </div>
              </td>
            </tr>

            <!-- Data Rows -->
            <tr v-for="tx in transactions" :key="tx.id" class="tx-row">
              <!-- Date & Time -->
              <td class="date-cell">
                <div class="fw-semibold text-dark">{{ formatDate(tx.created_at) }}</div>
                <small class="text-muted">{{ formatTime(tx.created_at) }}</small>
              </td>

              <!-- Reference -->
              <td class="ref-cell">
                <div class="d-flex align-items-center gap-1">
                  <code class="ref-code" :title="tx.transaction_reference">
                    {{ tx.transaction_reference }}
                  </code>
                  <button
                    type="button"
                    class="copy-btn"
                    title="Copy reference code"
                    @click="copyText(tx.transaction_reference)"
                  >
                    <i class="icofont-copy"></i>
                  </button>
                </div>
                <small v-if="tx.provider_request_id" class="provider-ref text-muted d-block" :title="tx.provider_request_id">
                  ioTec: {{ truncate(tx.provider_request_id, 14) }}
                </small>
              </td>

              <!-- Service / Application -->
              <td class="service-cell">
                <div class="service-title fw-semibold text-dark">
                  {{ tx.template_title || 'Direct Payment' }}
                </div>
                <small v-if="tx.submission_reference" class="text-muted">
                  Ref: {{ tx.submission_reference }}
                </small>
              </td>

              <!-- Payer Details -->
              <td class="payer-cell">
                <div class="d-flex align-items-center gap-2">
                  <div class="payer-avatar" :class="getAvatarColor(tx.applicant_name)">
                    {{ getInitials(tx.applicant_name) }}
                  </div>
                  <div>
                    <div class="payer-name fw-semibold text-dark">
                      {{ tx.applicant_name || 'System User' }}
                    </div>
                    <div v-if="tx.applicant_email" class="payer-email text-muted small">
                      <i class="icofont-ui-email"></i> {{ tx.applicant_email }}
                    </div>
                    <div v-if="tx.phone_number" class="payer-phone text-muted small">
                      <i class="icofont-phone"></i> {{ tx.phone_number }}
                    </div>
                  </div>
                </div>
              </td>

              <!-- Method / Channel -->
              <td class="method-cell">
                <span class="method-badge" :class="getMethodBadgeClass(tx.payment_method)">
                  <i :class="getMethodIcon(tx.payment_method)"></i>
                  {{ formatMethodName(tx.payment_method) }}
                </span>
                <small class="d-block text-muted mt-1">{{ tx.provider || 'ioTec Gateway' }}</small>
              </td>

              <!-- Amount (UGX) -->
              <td class="amount-cell text-right">
                <div class="amount-val fw-bold text-dark">
                  UGX {{ formatMoney(tx.amount_ugx) }}
                </div>
                <small class="text-muted">{{ tx.currency || 'UGX' }}</small>
              </td>

              <!-- Status -->
              <td class="status-cell text-center">
                <span class="status-pill" :class="getStatusClass(tx.status)">
                  <i :class="getStatusIcon(tx.status)"></i>
                  {{ tx.status || 'UNKNOWN' }}
                </span>
              </td>

              <!-- Actions -->
              <td class="actions-cell text-right">
                <div class="btn-group">
                  <button
                    type="button"
                    class="btn btn-sm btn-outline-primary"
                    title="View full transaction receipt & details"
                    @click="viewTransactionDetails(tx)"
                  >
                    <i class="icofont-eye-alt"></i> Details
                  </button>

                  <button
                    v-if="tx.status === 'PENDING'"
                    type="button"
                    class="btn btn-sm btn-outline-warning"
                    title="Re-query status from ioTec payment gateway"
                    :disabled="syncingId === tx.id"
                    @click="syncTransactionStatus(tx)"
                  >
                    <i class="icofont-refresh" :class="{ 'icofont-spin': syncingId === tx.id }"></i>
                    Sync
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination Footer -->
      <div class="pagination-footer d-flex align-items-center justify-content-between p-3 border-top">
        <div class="pagination-info text-muted small">
          Showing <strong>{{ pageStart }}</strong> to <strong>{{ pageEnd }}</strong> of <strong>{{ totalItems }}</strong> transactions
        </div>

        <div class="pagination-controls d-flex align-items-center gap-2">
          <div class="per-page-select d-flex align-items-center gap-1">
            <span class="text-muted small">Per page:</span>
            <select v-model="pagination.per_page" class="form-select form-select-sm" @change="onPerPageChange">
              <option :value="10">10</option>
              <option :value="20">20</option>
              <option :value="50">50</option>
              <option :value="100">100</option>
            </select>
          </div>

          <button
            type="button"
            class="btn btn-sm btn-outline-secondary"
            :disabled="pagination.page <= 1"
            @click="changePage(pagination.page - 1)"
          >
            <i class="icofont-rounded-left"></i> Prev
          </button>

          <span class="page-indicator small px-2">
            Page <strong>{{ pagination.page }}</strong> of <strong>{{ totalPages }}</strong>
          </span>

          <button
            type="button"
            class="btn btn-sm btn-outline-secondary"
            :disabled="pagination.page >= totalPages"
            @click="changePage(pagination.page + 1)"
          >
            Next <i class="icofont-rounded-right"></i>
          </button>
        </div>
      </div>
    </section>

    <!-- Transaction Details Modal -->
    <div
      v-if="selectedTx"
      class="modal fade show d-block tx-modal-backdrop"
      tabindex="-1"
      @click.self="selectedTx = null"
    >
      <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content shadow-lg border-0">
          <div class="modal-header bg-light">
            <h5 class="modal-title">
              <i class="icofont-receipt text-primary mr-2"></i>
              Transaction Details & Receipt
            </h5>
            <button type="button" class="btn-close" @click="selectedTx = null">&times;</button>
          </div>

          <div class="modal-body p-4" id="printable-receipt">
            <!-- Receipt Header -->
            <div class="receipt-header text-center pb-3 mb-3 border-bottom">
              <img src="/main-logo.png" alt="NCS Logo" class="receipt-logo mb-2" style="max-height: 48px;" />
              <h4 class="mb-0 fw-bold">National Council of Sports</h4>
              <p class="text-muted small mb-1">Official Electronic Payment Receipt</p>
              <span class="status-pill fs-6 mt-2" :class="getStatusClass(selectedTx.status)">
                <i :class="getStatusIcon(selectedTx.status)"></i> {{ selectedTx.status }}
              </span>
            </div>

            <!-- Amount Highlight -->
            <div class="amount-highlight text-center p-3 mb-4 bg-light rounded">
              <span class="text-muted small d-block text-uppercase fw-semibold">Amount Paid</span>
              <h2 class="text-primary fw-bold mb-0">UGX {{ formatMoney(selectedTx.amount_ugx) }}</h2>
              <small class="text-muted">{{ selectedTx.currency || 'UGX' }}</small>
            </div>

            <!-- Metadata Grid -->
            <div class="row g-3">
              <div class="col-md-6">
                <div class="meta-item">
                  <span class="meta-label text-muted small d-block">Transaction Reference</span>
                  <strong class="meta-value text-dark">{{ selectedTx.transaction_reference }}</strong>
                </div>
              </div>

              <div class="col-md-6">
                <div class="meta-item">
                  <span class="meta-label text-muted small d-block">ioTec Gateway Request ID</span>
                  <code class="meta-value">{{ selectedTx.provider_request_id || 'N/A' }}</code>
                </div>
              </div>

              <div class="col-md-6">
                <div class="meta-item">
                  <span class="meta-label text-muted small d-block">Service / Purpose</span>
                  <strong class="meta-value text-dark">{{ selectedTx.template_title || 'Portal Application Service' }}</strong>
                </div>
              </div>

              <div class="col-md-6">
                <div class="meta-item">
                  <span class="meta-label text-muted small d-block">Application Reference</span>
                  <strong class="meta-value text-dark">{{ selectedTx.submission_reference || 'Direct Payment' }}</strong>
                </div>
              </div>

              <div class="col-md-6">
                <div class="meta-item">
                  <span class="meta-label text-muted small d-block">Payer Full Name</span>
                  <strong class="meta-value text-dark">{{ selectedTx.applicant_name || 'System User' }}</strong>
                </div>
              </div>

              <div class="col-md-6">
                <div class="meta-item">
                  <span class="meta-label text-muted small d-block">Payer Email</span>
                  <span class="meta-value text-dark">{{ selectedTx.applicant_email || 'N/A' }}</span>
                </div>
              </div>

              <div class="col-md-6">
                <div class="meta-item">
                  <span class="meta-label text-muted small d-block">Payer Phone Number</span>
                  <span class="meta-value text-dark">{{ selectedTx.phone_number || 'N/A' }}</span>
                </div>
              </div>

              <div class="col-md-6">
                <div class="meta-item">
                  <span class="meta-label text-muted small d-block">Payment Channel & Gateway</span>
                  <span class="meta-value text-dark">{{ formatMethodName(selectedTx.payment_method) }} ({{ selectedTx.provider || 'ioTec' }})</span>
                </div>
              </div>

              <div class="col-md-6">
                <div class="meta-item">
                  <span class="meta-label text-muted small d-block">Created At</span>
                  <span class="meta-value text-dark">{{ formatFullDateTime(selectedTx.created_at) }}</span>
                </div>
              </div>

              <div class="col-md-6">
                <div class="meta-item">
                  <span class="meta-label text-muted small d-block">Completed / Settled At</span>
                  <span class="meta-value text-dark">{{ selectedTx.completed_at ? formatFullDateTime(selectedTx.completed_at) : 'Pending Settlement' }}</span>
                </div>
              </div>

              <div class="col-12" v-if="selectedTx.status_message">
                <div class="meta-item bg-light p-2 rounded">
                  <span class="meta-label text-muted small d-block">Gateway Status Message</span>
                  <span class="meta-value text-dark">{{ selectedTx.status_message }}</span>
                </div>
              </div>
            </div>

            <!-- Collapsible Gateway Raw Response -->
            <div class="mt-4 pt-3 border-top">
              <button
                type="button"
                class="btn btn-sm btn-link text-decoration-none p-0 text-muted d-flex align-items-center gap-1"
                @click="showRawJson = !showRawJson"
              >
                <i class="icofont-code-alt"></i>
                {{ showRawJson ? 'Hide' : 'View' }} Raw Provider Response Payload
              </button>
              <pre v-if="showRawJson" class="raw-json-block mt-2 p-3 bg-dark text-light rounded small"><code>{{ JSON.stringify(selectedTx.raw_response || {}, null, 2) }}</code></pre>
            </div>
          </div>

          <div class="modal-footer bg-light d-flex justify-content-between">
            <button
              v-if="selectedTx.status === 'PENDING'"
              type="button"
              class="btn btn-warning"
              :disabled="syncingId === selectedTx.id"
              @click="syncTransactionStatus(selectedTx)"
            >
              <i class="icofont-refresh" :class="{ 'icofont-spin': syncingId === selectedTx.id }"></i>
              Sync Status with ioTec
            </button>
            <div v-else></div>

            <div class="d-flex gap-2">
              <button
                type="button"
                class="btn btn-outline-primary"
                @click="printReceipt(selectedTx)"
              >
                <i class="icofont-print"></i> Print Receipt
              </button>
              <button type="button" class="btn btn-secondary" @click="selectedTx = null">
                Close
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import Swal from 'sweetalert2'
import {
  adminListTransactions,
  adminGetTransactionKPIs,
  adminGetTransaction,
  adminSyncTransaction,
} from '@/api/forms'

const loading = ref(false)
const error = ref('')
const transactions = ref([])
const totalItems = ref(0)
const syncingId = ref('')
const selectedTx = ref(null)
const showRawJson = ref(false)

const kpis = reactive({
  total_amount_ugx: 0,
  total_count: 0,
  success_count: 0,
  pending_count: 0,
  failed_count: 0,
})

const filters = reactive({
  search: '',
  status: '',
  payment_method: '',
  start_date: '',
  end_date: '',
  sort_by: 'created_at',
  sort_order: 'desc',
})

const pagination = reactive({
  page: 1,
  per_page: 20,
})

let searchTimer = null

const totalPages = computed(() => Math.ceil(totalItems.value / pagination.per_page) || 1)
const pageStart = computed(() => (totalItems.value === 0 ? 0 : (pagination.page - 1) * pagination.per_page + 1))
const pageEnd = computed(() => Math.min(pagination.page * pagination.per_page, totalItems.value))

onMounted(async () => {
  await Promise.all([loadKPIs(), loadData()])
})

async function loadKPIs() {
  try {
    const res = await adminGetTransactionKPIs()
    if (res) {
      Object.assign(kpis, res)
    }
  } catch (e) {
    console.error('Could not load KPIs:', e)
  }
}

async function loadData() {
  loading.value = true
  error.value = ''
  try {
    const params = {
      page: pagination.page,
      per_page: pagination.per_page,
      search: filters.search.trim(),
      status: filters.status,
      payment_method: filters.payment_method,
      start_date: filters.start_date,
      end_date: filters.end_date,
      sort_by: filters.sort_by,
      sort_order: filters.sort_order,
    }

    const res = await adminListTransactions(params)
    const items = res.data || res.items || []
    transactions.value = Array.isArray(items) ? items : []
    totalItems.value = res.meta?.total || transactions.value.length
  } catch (err) {
    error.value = err.response?.data?.error?.message || err.message || 'Could not load transactions.'
  } finally {
    loading.value = false
  }
}

function debounceSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    pagination.page = 1
    loadData()
  }, 400)
}

function onFilterChange() {
  pagination.page = 1
  loadData()
}

function resetFilters() {
  filters.search = ''
  filters.status = ''
  filters.payment_method = ''
  filters.start_date = ''
  filters.end_date = ''
  filters.sort_by = 'created_at'
  filters.sort_order = 'desc'
  pagination.page = 1
  loadData()
  loadKPIs()
}

function toggleSort(field) {
  let mapped = 'created_at'
  if (field === 'date') mapped = 'created_at'
  else if (field === 'reference') mapped = 'transaction_reference'
  else if (field === 'payer') mapped = 'applicant_name'
  else if (field === 'method') mapped = 'payment_method'
  else if (field === 'amount') mapped = 'amount_ugx'
  else if (field === 'status') mapped = 'status'

  if (filters.sort_by === mapped) {
    filters.sort_order = filters.sort_order === 'asc' ? 'desc' : 'asc'
  } else {
    filters.sort_by = mapped
    filters.sort_order = 'desc'
  }
  loadData()
}

function getSortIcon(field) {
  let mapped = 'created_at'
  if (field === 'date') mapped = 'created_at'
  else if (field === 'reference') mapped = 'transaction_reference'
  else if (field === 'payer') mapped = 'applicant_name'
  else if (field === 'method') mapped = 'payment_method'
  else if (field === 'amount') mapped = 'amount_ugx'
  else if (field === 'status') mapped = 'status'

  if (filters.sort_by !== mapped) return 'icofont-sort'
  return filters.sort_order === 'asc' ? 'icofont-sort-up text-primary' : 'icofont-sort-down text-primary'
}

function changePage(p) {
  if (p < 1 || p > totalPages.value) return
  pagination.page = p
  loadData()
}

function onPerPageChange() {
  pagination.page = 1
  loadData()
}

async function viewTransactionDetails(tx) {
  selectedTx.value = tx
  showRawJson.value = false
  try {
    const full = await adminGetTransaction(tx.id)
    if (full) selectedTx.value = full
  } catch (e) {
    console.error('Could not fetch deep transaction details:', e)
  }
}

async function syncTransactionStatus(tx) {
  syncingId.value = tx.id
  try {
    const res = await adminSyncTransaction(tx.id)
    if (res) {
      if (selectedTx.value && selectedTx.value.id === tx.id) {
        selectedTx.value = res
      }
      await Promise.all([loadData(), loadKPIs()])
      Swal.fire({
        title: 'Status Synchronized',
        text: `Transaction status updated to: ${res.status || 'Updated'}`,
        icon: res.status === 'SUCCESS' ? 'success' : 'info',
        timer: 2000,
        showConfirmButton: false,
      })
    }
  } catch (err) {
    Swal.fire('Sync Error', err.response?.data?.error?.message || err.message, 'error')
  } finally {
    syncingId.value = ''
  }
}

function copyText(text) {
  if (!text) return
  navigator.clipboard?.writeText(text)
  Swal.fire({
    title: 'Copied!',
    text: `Reference ${text} copied to clipboard`,
    icon: 'success',
    timer: 1500,
    showConfirmButton: false,
  })
}

// ── EXPORTING HELPERS ───────────────────────────────────────────────────────

function exportToExcel() {
  const params = new URLSearchParams({
    export: 'csv',
    search: filters.search.trim(),
    status: filters.status,
    payment_method: filters.payment_method,
    start_date: filters.start_date,
    end_date: filters.end_date,
    sort_by: filters.sort_by,
    sort_order: filters.sort_order,
  })
  const url = `/api/v1/admin/transactions?${params.toString()}`

  // Client-side direct stream download with token
  const authData = JSON.parse(localStorage.getItem('ncsms_auth') || '{}')
  const token = authData.token || authData.access_token || localStorage.getItem('ncsms_token') || ''

  fetch(url, {
    headers: { Authorization: `Bearer ${token}` },
  })
    .then(r => r.blob())
    .then(blob => {
      const nowStr = new Date().toISOString().slice(0, 10)
      const downloadUrl = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = downloadUrl
      a.download = `NCS_Transactions_Ledger_${nowStr}.csv`
      document.body.appendChild(a)
      a.click()
      a.remove()
      window.URL.revokeObjectURL(downloadUrl)
      Swal.fire({
        title: 'Export Complete',
        text: 'Transactions ledger exported successfully as Excel/CSV.',
        icon: 'success',
        timer: 2000,
        showConfirmButton: false,
      })
    })
    .catch(err => {
      Swal.fire('Export Failed', err.message, 'error')
    })
}

function exportToPdf() {
  const printWin = window.open('', '_blank', 'width=1000,height=800')
  if (!printWin) {
    Swal.fire('Popup Blocked', 'Please allow popups to generate and print PDF reports.', 'warning')
    return
  }

  const generatedDate = new Date().toLocaleString('en-UG', {
    dateStyle: 'full',
    timeStyle: 'medium',
  })

  const rowsHtml = transactions.value
    .map(
      (tx, idx) => `
      <tr>
        <td>${idx + 1}</td>
        <td>${formatDate(tx.created_at)} ${formatTime(tx.created_at)}</td>
        <td><strong>${tx.transaction_reference}</strong><br><small style="color:#666">${tx.provider_request_id || ''}</small></td>
        <td>${escapeHtml(tx.applicant_name || 'System User')}<br><small style="color:#666">${escapeHtml(tx.phone_number || '')}</small></td>
        <td>${escapeHtml(tx.template_title || tx.submission_reference || 'General')}</td>
        <td>${formatMethodName(tx.payment_method)}</td>
        <td style="text-align:right; font-weight:bold;">UGX ${formatMoney(tx.amount_ugx)}</td>
        <td style="text-align:center;"><span class="badge status-${String(tx.status || '').toLowerCase()}">${tx.status}</span></td>
      </tr>
    `
    )
    .join('')

  const html = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>NCS Transactions Audit Ledger Report</title>
  <style>
    @page { size: A4 landscape; margin: 12mm; }
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Arial, sans-serif; color: #222; margin: 0; padding: 20px; font-size: 12px; }
    .header { text-align: center; border-bottom: 2px solid #0d6efd; padding-bottom: 12px; margin-bottom: 16px; }
    .header h2 { margin: 0 0 4px; font-size: 20px; color: #111; }
    .header p { margin: 0; color: #555; font-size: 13px; }
    .kpi-summary { display: flex; justify-content: space-between; margin-bottom: 16px; background: #f8f9fa; border: 1px solid #e9ecef; border-radius: 6px; padding: 10px 16px; }
    .kpi-item strong { display: block; font-size: 15px; color: #0d6efd; }
    .kpi-item span { font-size: 11px; color: #666; text-transform: uppercase; }
    table { width: 100%; border-collapse: collapse; margin-top: 10px; }
    th { background: #f1f5f9; color: #334155; text-align: left; padding: 8px 6px; font-size: 11px; border-bottom: 2px solid #cbd5e1; }
    td { padding: 7px 6px; border-bottom: 1px solid #e2e8f0; font-size: 11px; vertical-align: top; }
    tr:nth-child(even) { background: #f8fafc; }
    .badge { display: inline-block; padding: 3px 6px; border-radius: 4px; font-size: 10px; font-weight: bold; }
    .status-success { background: #d1fae5; color: #065f46; }
    .status-pending { background: #fef3c7; color: #92400e; }
    .status-failed { background: #fee2e2; color: #991b1b; }
    .footer { margin-top: 24px; text-align: center; font-size: 10px; color: #888; border-top: 1px solid #ddd; padding-top: 8px; }
  </style>
</head>
<body>
  <div class="header">
    <h2>NATIONAL COUNCIL OF SPORTS (NCS)</h2>
    <p>Financial Transactions & Electronic Payments Ledger Report</p>
    <p style="font-size: 11px; color: #777; margin-top: 4px;">Generated on: ${generatedDate}</p>
  </div>

  <div class="kpi-summary">
    <div class="kpi-item">
      <span>Total Revenue</span>
      <strong>UGX ${formatMoney(kpis.total_amount_ugx || 0)}</strong>
    </div>
    <div class="kpi-item">
      <span>Total Count</span>
      <strong>${kpis.total_count || 0}</strong>
    </div>
    <div class="kpi-item">
      <span>Successful</span>
      <strong style="color:#198754">${kpis.success_count || 0}</strong>
    </div>
    <div class="kpi-item">
      <span>Pending</span>
      <strong style="color:#ffc107">${kpis.pending_count || 0}</strong>
    </div>
    <div class="kpi-item">
      <span>Failed</span>
      <strong style="color:#dc3545">${kpis.failed_count || 0}</strong>
    </div>
  </div>

  <table>
    <thead>
      <tr>
        <th style="width:24px">#</th>
        <th>Date & Time</th>
        <th>Reference</th>
        <th>Payer Details</th>
        <th>Purpose / Service</th>
        <th>Channel</th>
        <th style="text-align:right">Amount (UGX)</th>
        <th style="text-align:center">Status</th>
      </tr>
    </thead>
    <tbody>
      ${rowsHtml}
    </tbody>
  </table>

  <div class="footer">
    <p>National Council of Sports Portal · Confidential Financial Record · Printed by System Administrator</p>
  </div>

  <script>
    window.onload = function() {
      setTimeout(function() { window.print(); }, 500);
    };
  <\/script>
</body>
</html>`

  printWin.document.write(html)
  printWin.document.close()
}

function printReceipt(tx) {
  const printWin = window.open('', '_blank', 'width=700,height=750')
  if (!printWin) return

  const dateStr = formatFullDateTime(tx.created_at)

  const html = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Receipt - ${tx.transaction_reference}</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Arial, sans-serif; color: #222; margin: 0; padding: 30px; font-size: 13px; }
    .receipt-box { max-width: 550px; margin: auto; border: 1px solid #ddd; padding: 24px; border-radius: 8px; }
    .header { text-align: center; border-bottom: 2px solid #0d6efd; padding-bottom: 16px; margin-bottom: 20px; }
    .header h2 { margin: 0 0 4px; font-size: 20px; color: #111; }
    .amount-box { text-align: center; background: #f0f7ff; border: 1px solid #cce3ff; border-radius: 6px; padding: 14px; margin-bottom: 20px; }
    .amount-box h1 { margin: 0; color: #0d6efd; font-size: 28px; }
    .row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #eee; }
    .label { color: #666; font-size: 12px; }
    .val { font-weight: 600; text-align: right; }
    .footer { text-align: center; margin-top: 24px; color: #888; font-size: 11px; }
  </style>
</head>
<body>
  <div class="receipt-box">
    <div class="header">
      <h2>National Council of Sports</h2>
      <p style="margin:0; color:#666;">Official Payment Confirmation Receipt</p>
    </div>

    <div class="amount-box">
      <span style="font-size:12px; color:#555; text-transform:uppercase;">Amount Paid</span>
      <h1>UGX ${formatMoney(tx.amount_ugx)}</h1>
      <span style="font-size:12px; color:#198754; font-weight:bold;">Status: ${tx.status}</span>
    </div>

    <div class="row"><span class="label">Reference Number:</span><span class="val">${tx.transaction_reference}</span></div>
    <div class="row"><span class="label">Gateway Request ID:</span><span class="val">${tx.provider_request_id || 'N/A'}</span></div>
    <div class="row"><span class="label">Service Purpose:</span><span class="val">${escapeHtml(tx.template_title || 'Portal Service')}</span></div>
    <div class="row"><span class="label">Application Reference:</span><span class="val">${escapeHtml(tx.submission_reference || 'Direct')}</span></div>
    <div class="row"><span class="label">Payer Name:</span><span class="val">${escapeHtml(tx.applicant_name || 'System User')}</span></div>
    <div class="row"><span class="label">Payer Email:</span><span class="val">${escapeHtml(tx.applicant_email || 'N/A')}</span></div>
    <div class="row"><span class="label">Payer Phone:</span><span class="val">${escapeHtml(tx.phone_number || 'N/A')}</span></div>
    <div class="row"><span class="label">Payment Method:</span><span class="val">${formatMethodName(tx.payment_method)} (${tx.provider || 'ioTec'})</span></div>
    <div class="row"><span class="label">Date & Time:</span><span class="val">${dateStr}</span></div>

    <div class="footer">
      <p>Thank you for your payment to the National Council of Sports.<br>This receipt was automatically generated and is valid without a physical seal.</p>
    </div>
  </div>

  <script>
    window.onload = function() {
      setTimeout(function() { window.print(); }, 400);
    };
  <\/script>
</body>
</html>`

  printWin.document.write(html)
  printWin.document.close()
}

// ── FORMATTING UTILITIES ───────────────────────────────────────────────────

function formatMoney(amount) {
  return Number(amount || 0).toLocaleString('en-US', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  })
}

function formatPercent(count, total) {
  if (!total) return '0%'
  return `${Math.round(((count || 0) / total) * 100)}%`
}

function formatDate(val) {
  if (!val) return '-'
  try {
    return new Date(val).toLocaleDateString('en-UG', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  } catch {
    return String(val)
  }
}

function formatTime(val) {
  if (!val) return ''
  try {
    return new Date(val).toLocaleTimeString('en-UG', {
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return ''
  }
}

function formatFullDateTime(val) {
  if (!val) return '-'
  try {
    return new Date(val).toLocaleString('en-UG', {
      dateStyle: 'medium',
      timeStyle: 'medium',
    })
  } catch {
    return String(val)
  }
}

function formatMethodName(method) {
  const m = String(method || '').toUpperCase()
  if (m === 'MOBILE_MONEY') return 'Mobile Money'
  if (m === 'AIRTEL_MONEY') return 'Airtel Money'
  if (m === 'MTN_MOMO') return 'MTN MoMo'
  if (m === 'BANK_TRANSFER') return 'Bank Transfer'
  if (m === 'CARD' || m === 'VISA_MASTERCARD') return 'Credit / Debit Card'
  return m || 'Mobile Money'
}

function getMethodIcon(method) {
  const m = String(method || '').toUpperCase()
  if (m.includes('MOBILE') || m.includes('MOMO') || m.includes('AIRTEL')) return 'icofont-smart-phone'
  if (m.includes('BANK')) return 'icofont-bank-alt'
  if (m.includes('CARD')) return 'icofont-credit-card'
  return 'icofont-money-bag'
}

function getMethodBadgeClass(method) {
  const m = String(method || '').toUpperCase()
  if (m.includes('MOMO') || m.includes('MTN')) return 'method-momo'
  if (m.includes('AIRTEL')) return 'method-airtel'
  if (m.includes('BANK')) return 'method-bank'
  return 'method-general'
}

function getStatusClass(status) {
  const s = String(status || '').toUpperCase()
  if (s === 'SUCCESS') return 'status-success'
  if (s === 'PENDING') return 'status-pending'
  if (s === 'FAILED') return 'status-failed'
  if (s === 'CANCELLED') return 'status-cancelled'
  return 'status-unknown'
}

function getStatusIcon(status) {
  const s = String(status || '').toUpperCase()
  if (s === 'SUCCESS') return 'icofont-check-circled'
  if (s === 'PENDING') return 'icofont-spinner icofont-spin'
  if (s === 'FAILED') return 'icofont-close-circled'
  if (s === 'CANCELLED') return 'icofont-ban'
  return 'icofont-question-circle'
}

function getInitials(name) {
  if (!name) return 'U'
  const parts = String(name).trim().split(/\s+/)
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.slice(0, 2).toUpperCase()
}

function getAvatarColor(name) {
  const colors = ['bg-primary', 'bg-success', 'bg-info', 'bg-warning', 'bg-danger', 'bg-dark']
  if (!name) return colors[0]
  let sum = 0
  for (let i = 0; i < name.length; i++) {
    sum += name.charCodeAt(i)
  }
  return colors[sum % colors.length]
}

function truncate(str, max = 15) {
  if (!str) return ''
  return str.length > max ? `${str.slice(0, max)}...` : str
}

function escapeHtml(str) {
  return String(str || '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;')
}
</script>

<style scoped>
.manage-transactions-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
}

/* Header Row */
.panel-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
}

.panel-title {
  font-size: 22px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
  display: flex;
  align-items: center;
}

.panel-subtitle {
  color: #64748b;
  font-size: 13px;
  margin: 4px 0 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.export-btn {
  font-weight: 600;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
}

/* KPI Cards */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.kpi-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 16px 18px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  transition: transform 0.2s, box-shadow 0.2s;
}

.kpi-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.kpi-icon-wrap {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.kpi-body {
  display: flex;
  flex-direction: column;
}

.kpi-label {
  font-size: 12px;
  color: #64748b;
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.kpi-value {
  font-size: 20px;
  font-weight: 700;
  margin: 2px 0 0;
}

.kpi-meta {
  font-size: 11px;
  color: #94a3b8;
}

/* Filter Card */
.filter-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 16px 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.filter-row {
  display: flex;
  align-items: flex-end;
  flex-wrap: wrap;
  gap: 14px;
}

.search-field {
  flex: 1 1 260px;
  position: relative;
  display: flex;
  align-items: center;
}

.search-field i {
  position: absolute;
  left: 12px;
  color: #94a3b8;
  font-size: 16px;
}

.search-field input {
  width: 100%;
  padding: 8px 36px 8px 36px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 13px;
  transition: border-color 0.2s;
}

.search-field input:focus {
  border-color: #3b82f6;
  outline: none;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.clear-btn {
  position: absolute;
  right: 10px;
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 0;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 0 1 160px;
}

.filter-item label {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.3px;
  margin: 0;
}

.filter-item select,
.filter-item input {
  padding: 8px 10px;
  font-size: 13px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
}

.filter-reset {
  flex: 0 0 auto;
}

/* Table Card */
.table-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.transactions-table th {
  background: #f8fafc;
  color: #475569;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 12px 16px;
  border-bottom: 2px solid #e2e8f0;
}

.transactions-table td {
  padding: 14px 16px;
  vertical-align: middle;
  border-bottom: 1px solid #f1f5f9;
  font-size: 13px;
}

.sortable-th {
  cursor: pointer;
  user-select: none;
}

.sortable-th:hover {
  background: #f1f5f9;
  color: #1e293b;
}

.ref-code {
  background: #f1f5f9;
  color: #0f172a;
  padding: 3px 6px;
  border-radius: 4px;
  font-weight: 600;
  font-size: 12px;
}

.copy-btn {
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 2px 4px;
  font-size: 13px;
  border-radius: 4px;
}

.copy-btn:hover {
  color: #3b82f6;
  background: #eff6ff;
}

.payer-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  color: #ffffff;
  font-weight: 700;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.method-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #334155;
}

.method-momo {
  background: #fffbeb;
  border-color: #fde68a;
  color: #92400e;
}

.method-airtel {
  background: #fef2f2;
  border-color: #fecaca;
  color: #b91c1c;
}

.method-bank {
  background: #f0fdf4;
  border-color: #bbf7d0;
  color: #15803d;
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.3px;
  text-transform: uppercase;
}

.status-success {
  background: #dcfce7;
  color: #15803d;
}

.status-pending {
  background: #fef3c7;
  color: #b45309;
}

.status-failed {
  background: #fee2e2;
  color: #b91c1c;
}

.status-cancelled {
  background: #f1f5f9;
  color: #64748b;
}

.tx-modal-backdrop {
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(3px);
}

.raw-json-block {
  max-height: 200px;
  overflow-y: auto;
}
</style>
