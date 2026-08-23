<template>
  <div class="manage-licenses-panel">
    <!-- Header Row -->
    <div class="panel-header-row">
      <div class="header-titles">
        <h2 class="panel-title">
          <i class="icofont-license text-primary mr-2"></i>
          National Sports Federation Licenses & Compliance
        </h2>
        <p class="panel-subtitle">
          Monitor statutory recognition, validity periods, extensions, revocations, and immutable compliance audit trails.
        </p>
      </div>

      <div class="header-actions">
        <button
          type="button"
          class="btn btn-outline-secondary"
          :disabled="loading || !licenses.length"
          @click="exportToExcel"
          title="Export licenses table to Excel/CSV"
        >
          <i class="icofont-file-excel text-success"></i> Export Excel
        </button>

        <button
          type="button"
          class="btn btn-outline-secondary"
          :disabled="loading || !licenses.length"
          @click="exportToPdf"
          title="Print official PDF license registry summary"
        >
          <i class="icofont-file-pdf text-danger"></i> Export PDF
        </button>

        <button
          type="button"
          class="btn btn-primary"
          @click="$emit('navigate', 'create-federation-licenses')"
        >
          <i class="icofont-plus-circle"></i> Issue New License
        </button>
      </div>
    </div>

    <!-- KPI Summary Metrics -->
    <section class="kpi-grid">
      <div class="kpi-card">
        <div class="kpi-icon-wrap bg-primary-subtle text-primary">
          <i class="icofont-certificate"></i>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">Total Licenses</span>
          <h3 class="kpi-value text-primary">{{ kpis.total_licenses || 0 }}</h3>
          <span class="kpi-meta">All-time issued</span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon-wrap bg-success-subtle text-success">
          <i class="icofont-check-circled"></i>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">Active Recognition</span>
          <h3 class="kpi-value text-success">{{ kpis.active_licenses || 0 }}</h3>
          <span class="kpi-meta">Fully compliant</span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon-wrap bg-info-subtle text-info">
          <i class="icofont-history"></i>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">Extended Validity</span>
          <h3 class="kpi-value text-info">{{ kpis.extended_licenses || 0 }}</h3>
          <span class="kpi-meta">Prolonged compliance</span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon-wrap bg-danger-subtle text-danger">
          <i class="icofont-ban"></i>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">Revoked / Suspended</span>
          <h3 class="kpi-value text-danger">{{ kpis.revoked_licenses || 0 }}</h3>
          <span class="kpi-meta">Non-compliant or halted</span>
        </div>
      </div>

      <div class="kpi-card">
        <div class="kpi-icon-wrap bg-warning-subtle text-warning">
          <i class="icofont-warning-alt"></i>
        </div>
        <div class="kpi-body">
          <span class="kpi-label">Expiring Soon (&lt;60d)</span>
          <h3 class="kpi-value text-warning">{{ kpis.expiring_soon || 0 }}</h3>
          <span class="kpi-meta">Requires renewal audit</span>
        </div>
      </div>
    </section>

    <!-- Filters & Search Toolbar -->
    <section class="filter-card">
      <div class="filter-row">
        <!-- Search Field -->
        <div class="search-field">
          <i class="icofont-search"></i>
          <input
            v-model="filters.search"
            type="text"
            placeholder="Search by License Number, Federation, Reg No, President..."
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
          <label>License Status</label>
          <select v-model="filters.status" class="form-select" @change="onFilterChange">
            <option value="">All Statuses</option>
            <option value="ACTIVE">ACTIVE</option>
            <option value="EXTENDED">EXTENDED</option>
            <option value="EXPIRED">EXPIRED</option>
            <option value="SUSPENDED">SUSPENDED</option>
            <option value="REVOKED">REVOKED</option>
          </select>
        </div>

        <!-- License Type Filter -->
        <div class="filter-item">
          <label>Recognition Type</label>
          <select v-model="filters.license_type" class="form-select" @change="onFilterChange">
            <option value="">All Types</option>
            <option value="FULL_RECOGNITION">Full Recognition</option>
            <option value="PROVISIONAL">Provisional</option>
            <option value="ANNUAL_COMPLIANCE">Annual Compliance</option>
            <option value="SPECIAL_CLEARANCE">Special Clearance</option>
          </select>
        </div>

        <!-- Reset Button -->
        <div class="filter-item filter-reset">
          <button
            type="button"
            class="btn btn-light"
            @click="resetFilters"
            title="Reset all filters"
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

    <!-- Licenses Data Table -->
    <section class="table-card">
      <div class="table-responsive">
        <table class="table table-hover licenses-table mb-0">
          <thead>
            <tr>
              <th @click="toggleSort('license_number')" class="sortable-th">
                License Number <i :class="getSortIcon('license_number')"></i>
              </th>
              <th @click="toggleSort('federation')" class="sortable-th">
                Federation Details <i :class="getSortIcon('federation')"></i>
              </th>
              <th @click="toggleSort('type')" class="sortable-th">
                Recognition Type & Tier <i :class="getSortIcon('type')"></i>
              </th>
              <th @click="toggleSort('expiry_date')" class="sortable-th">
                Validity & Expiry <i :class="getSortIcon('expiry_date')"></i>
              </th>
              <th @click="toggleSort('status')" class="sortable-th text-center">
                Status <i :class="getSortIcon('status')"></i>
              </th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <!-- Skeleton loading -->
            <tr v-if="loading && !licenses.length">
              <td colspan="6" class="text-center py-5">
                <div class="spinner-border text-primary" role="status">
                  <span class="sr-only">Loading licenses...</span>
                </div>
                <p class="text-muted mt-2">Loading federation licenses from database...</p>
              </td>
            </tr>

            <!-- Empty State -->
            <tr v-else-if="!loading && !licenses.length">
              <td colspan="6" class="text-center py-5">
                <div class="empty-state-wrap">
                  <i class="icofont-license fs-1 text-muted"></i>
                  <h5 class="mt-3">No Federation Licenses Found</h5>
                  <p class="text-muted">No license records match your current filter settings.</p>
                  <button
                    type="button"
                    class="btn btn-sm btn-primary mt-2"
                    @click="$emit('navigate', 'create-federation-licenses')"
                  >
                    Issue New License
                  </button>
                </div>
              </td>
            </tr>

            <!-- Table Rows -->
            <tr v-for="lic in licenses" :key="lic.id" class="lic-row">
              <!-- License Number -->
              <td class="lic-num-cell">
                <div class="d-flex align-items-center gap-1">
                  <code class="license-code" :title="lic.license_number">
                    {{ lic.license_number }}
                  </code>
                  <button
                    type="button"
                    class="copy-btn"
                    title="Copy license number"
                    @click="copyText(lic.license_number)"
                  >
                    <i class="icofont-copy"></i>
                  </button>
                </div>
                <small class="text-muted d-block mt-1">
                  Issued: {{ formatDate(lic.issue_date) }}
                </small>
              </td>

              <!-- Federation Details -->
              <td class="fed-cell">
                <div class="d-flex align-items-center gap-2">
                  <div class="fed-avatar">
                    <img v-if="lic.federation_logo" :src="lic.federation_logo" alt="" class="avatar-img" />
                    <span v-else class="avatar-fallback">{{ (lic.federation_acronym || lic.federation_name.slice(0, 3)).toUpperCase() }}</span>
                  </div>
                  <div>
                    <strong class="fed-name text-dark d-block">
                      {{ lic.federation_name }}
                      <span v-if="lic.federation_acronym" class="badge-acronym">{{ lic.federation_acronym }}</span>
                    </strong>
                    <div class="fed-meta text-muted small">
                      <span v-if="lic.federation_reg_no">Reg: {{ lic.federation_reg_no }}</span>
                      <span v-if="lic.president_name" class="ml-1">· Pres: {{ lic.president_name }}</span>
                    </div>
                  </div>
                </div>
              </td>

              <!-- Recognition Type & Tier -->
              <td class="type-cell">
                <span class="type-badge" :class="getTypeBadgeClass(lic.license_type)">
                  {{ formatTypeName(lic.license_type) }}
                </span>
                <small class="text-muted d-block mt-1">{{ lic.category || 'National Sports Federation' }}</small>
              </td>

              <!-- Validity & Expiry -->
              <td class="validity-cell">
                <div class="expiry-date fw-semibold text-dark">
                  {{ formatDate(lic.expiry_date) }}
                </div>
                <span class="validity-tag" :class="getValidityTagClass(lic)">
                  {{ getValidityTagText(lic) }}
                </span>
              </td>

              <!-- Status Badge -->
              <td class="status-cell text-center">
                <span class="status-pill" :class="getStatusClass(lic.status)">
                  <i :class="getStatusIcon(lic.status)"></i>
                  {{ lic.status }}
                </span>
              </td>

              <!-- Actions Dropdown / Toolbar -->
              <td class="actions-cell text-right">
                <div class="btn-group">
                  <button
                    type="button"
                    class="btn btn-sm btn-outline-primary"
                    title="View license details & immutable audit history"
                    @click="openDetailsModal(lic)"
                  >
                    <i class="icofont-eye-alt"></i> Details & Logs
                  </button>

                  <button
                    type="button"
                    class="btn btn-sm btn-outline-success"
                    title="Print official NCS Recognition Certificate"
                    @click="printCertificate(lic)"
                  >
                    <i class="icofont-print"></i> Certificate
                  </button>

                  <button
                    v-if="lic.status !== 'REVOKED'"
                    type="button"
                    class="btn btn-sm btn-outline-info"
                    title="Extend license validity period"
                    @click="openExtendModal(lic)"
                  >
                    <i class="icofont-calendar"></i> Extend
                  </button>

                  <button
                    v-if="lic.status !== 'REVOKED'"
                    type="button"
                    class="btn btn-sm btn-outline-danger"
                    title="Revoke license compliance status"
                    @click="openRevokeModal(lic)"
                  >
                    <i class="icofont-ban"></i> Revoke
                  </button>

                  <button
                    v-else
                    type="button"
                    class="btn btn-sm btn-outline-success"
                    title="Reinstate revoked license"
                    @click="openReinstateModal(lic)"
                  >
                    <i class="icofont-check-circled"></i> Reinstate
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
          Showing <strong>{{ pageStart }}</strong> to <strong>{{ pageEnd }}</strong> of <strong>{{ totalItems }}</strong> federation licenses
        </div>

        <div class="pagination-controls d-flex align-items-center gap-2">
          <div class="per-page-select d-flex align-items-center gap-1">
            <span class="text-muted small">Per page:</span>
            <select v-model="pagination.per_page" class="form-select form-select-sm" @change="onPerPageChange">
              <option :value="10">10</option>
              <option :value="20">20</option>
              <option :value="50">50</option>
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

    <!-- Details & Audit Log Modal -->
    <div
      v-if="selectedLic"
      class="modal fade show d-block lic-modal-backdrop"
      tabindex="-1"
      @click.self="selectedLic = null"
    >
      <div class="modal-dialog modal-dialog-centered modal-lg">
        <div class="modal-content shadow-lg border-0">
          <div class="modal-header bg-light">
            <h5 class="modal-title d-flex align-items-center gap-2">
              <i class="icofont-certificate text-primary"></i>
              License {{ selectedLic.license_number }}
            </h5>
            <button type="button" class="btn-close" @click="selectedLic = null">&times;</button>
          </div>

          <div class="modal-body p-4">
            <!-- Federation Identity -->
            <div class="d-flex align-items-center justify-content-between pb-3 mb-3 border-bottom">
              <div>
                <h4 class="mb-0 fw-bold">{{ selectedLic.federation_name }}</h4>
                <p class="text-muted small mb-0">
                  Reg No: <strong>{{ selectedLic.federation_reg_no || 'N/A' }}</strong> · Category: <strong>{{ selectedLic.category }}</strong>
                </p>
              </div>
              <div>
                <span class="status-pill fs-6" :class="getStatusClass(selectedLic.status)">
                  <i :class="getStatusIcon(selectedLic.status)"></i> {{ selectedLic.status }}
                </span>
              </div>
            </div>

            <!-- Key Dates Grid -->
            <div class="row g-3 mb-4">
              <div class="col-md-4">
                <div class="p-3 bg-light rounded">
                  <span class="text-muted small d-block text-uppercase">Issue Date</span>
                  <strong class="text-dark">{{ formatDate(selectedLic.issue_date) }}</strong>
                  <small class="text-muted d-block">Issued By: {{ selectedLic.issued_by_name || 'System Admin' }}</small>
                </div>
              </div>

              <div class="col-md-4">
                <div class="p-3 bg-light rounded">
                  <span class="text-muted small d-block text-uppercase">Expiry Date</span>
                  <strong class="text-dark">{{ formatDate(selectedLic.expiry_date) }}</strong>
                  <small class="text-muted d-block">
                    <span v-if="selectedLic.extended_at">Extended By: {{ selectedLic.extended_by_name }}</span>
                    <span v-else>Standard Term</span>
                  </small>
                </div>
              </div>

              <div class="col-md-4">
                <div class="p-3 bg-light rounded">
                  <span class="text-muted small d-block text-uppercase">Recognition Type</span>
                  <strong class="text-dark">{{ formatTypeName(selectedLic.license_type) }}</strong>
                  <small class="text-muted d-block">Statutory NCS Recognition</small>
                </div>
              </div>
            </div>

            <!-- Conditions & Remarks -->
            <div v-if="selectedLic.conditions" class="mb-4">
              <label class="fw-bold small text-muted text-uppercase mb-1 d-block">Statutory Conditions</label>
              <div class="p-3 bg-light rounded text-dark small border">
                {{ selectedLic.conditions }}
              </div>
            </div>

            <!-- Extension / Revocation Notice Banner -->
            <div v-if="selectedLic.status === 'EXTENDED' && selectedLic.extension_reason" class="alert alert-info py-2 px-3 small mb-4">
              <i class="icofont-info-circle mr-1"></i>
              <strong>Extension Note:</strong> {{ selectedLic.extension_reason }}
              <span class="text-muted">(Extended on {{ formatDate(selectedLic.extended_at) }} by {{ selectedLic.extended_by_name }})</span>
            </div>

            <div v-if="selectedLic.status === 'REVOKED' && selectedLic.revocation_reason" class="alert alert-danger py-2 px-3 small mb-4">
              <i class="icofont-warning mr-1"></i>
              <strong>Revocation Reason:</strong> {{ selectedLic.revocation_reason }}
              <span class="text-muted">(Revoked on {{ formatDate(selectedLic.revoked_at) }} by {{ selectedLic.revoked_by_name }})</span>
            </div>

            <!-- Immutable Audit Trail Timeline -->
            <div class="audit-trail-section mt-4 pt-3 border-top">
              <h6 class="fw-bold text-dark mb-3 d-flex align-items-center gap-2">
                <i class="icofont-history text-primary"></i>
                Compliance & Action Audit History
              </h6>

              <div v-if="!selectedLic.logs || !selectedLic.logs.length" class="text-muted small">
                No historic audit logs recorded yet.
              </div>

              <div v-else class="timeline-container">
                <div v-for="log in selectedLic.logs" :key="log.id" class="timeline-item">
                  <div class="timeline-dot" :class="getTimelineDotClass(log.action)"></div>
                  <div class="timeline-content">
                    <div class="d-flex align-items-center justify-content-between">
                      <strong class="timeline-action text-dark">{{ formatLogAction(log.action) }}</strong>
                      <span class="timeline-time text-muted small">{{ formatFullDateTime(log.created_at) }}</span>
                    </div>
                    <p class="timeline-reason text-muted small mb-1">{{ log.reason || log.notes }}</p>
                    <div class="timeline-performer text-muted small">
                      <i class="icofont-user-alt-7"></i> Action Performed By: <strong>{{ log.performed_by_name || 'Administrator' }}</strong>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="modal-footer bg-light d-flex justify-content-between">
            <button
              type="button"
              class="btn btn-outline-primary"
              @click="printCertificate(selectedLic)"
            >
              <i class="icofont-print"></i> Print Official Certificate
            </button>

            <button type="button" class="btn btn-secondary" @click="selectedLic = null">
              Close
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Extend License Modal -->
    <div
      v-if="extendModalLic"
      class="modal fade show d-block lic-modal-backdrop"
      tabindex="-1"
      @click.self="extendModalLic = null"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content shadow border-0">
          <div class="modal-header bg-light">
            <h5 class="modal-title">
              <i class="icofont-calendar text-info mr-2"></i>
              Extend License Validity
            </h5>
            <button type="button" class="btn-close" @click="extendModalLic = null">&times;</button>
          </div>
          <form @submit.prevent="submitExtendLicense">
            <div class="modal-body p-4">
              <p class="small text-muted mb-3">
                Extending statutory validity for <strong>{{ extendModalLic.federation_name }}</strong> (License: {{ extendModalLic.license_number }}).
                Current expiry: <strong>{{ formatDate(extendModalLic.expiry_date) }}</strong>.
              </p>

              <div class="mb-3">
                <label class="form-label fw-bold">New Expiry Date <span class="text-danger">*</span></label>
                <input
                  v-model="extendForm.new_expiry_date"
                  type="date"
                  class="form-control"
                  required
                />
              </div>

              <div class="mb-3">
                <label class="form-label fw-bold">Extension Rationale & Justification <span class="text-danger">*</span></label>
                <textarea
                  v-model="extendForm.reason"
                  class="form-control"
                  rows="3"
                  placeholder="State the regulatory basis, audit completion, or board resolution granting this extension..."
                  required
                ></textarea>
              </div>
            </div>
            <div class="modal-footer bg-light">
              <button type="button" class="btn btn-light" @click="extendModalLic = null">Cancel</button>
              <button type="submit" class="btn btn-info" :disabled="submittingAction">
                {{ submittingAction ? 'Extending...' : 'Confirm Extension' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Revoke License Modal -->
    <div
      v-if="revokeModalLic"
      class="modal fade show d-block lic-modal-backdrop"
      tabindex="-1"
      @click.self="revokeModalLic = null"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content shadow border-0">
          <div class="modal-header bg-danger text-white">
            <h5 class="modal-title text-white">
              <i class="icofont-warning mr-2"></i>
              Revoke Federation Recognition License
            </h5>
            <button type="button" class="btn-close btn-close-white" @click="revokeModalLic = null">&times;</button>
          </div>
          <form @submit.prevent="submitRevokeLicense">
            <div class="modal-body p-4">
              <div class="alert alert-warning small mb-3">
                <strong>Caution:</strong> Revoking this license will immediately suspend active statutory accreditation for <strong>{{ revokeModalLic.federation_name }}</strong> and flag non-compliance across the portal.
              </div>

              <div class="mb-3">
                <label class="form-label fw-bold">Grounds for Revocation <span class="text-danger">*</span></label>
                <textarea
                  v-model="revokeForm.reason"
                  class="form-control"
                  rows="4"
                  placeholder="Specify regulatory violations, governance failure, statutory breach, or board directive..."
                  required
                ></textarea>
              </div>
            </div>
            <div class="modal-footer bg-light">
              <button type="button" class="btn btn-light" @click="revokeModalLic = null">Cancel</button>
              <button type="submit" class="btn btn-danger" :disabled="submittingAction">
                {{ submittingAction ? 'Revoking...' : 'Confirm Revocation' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import Swal from 'sweetalert2'
import {
  listFederationLicenses,
  getFederationLicense,
  extendFederationLicense,
  revokeFederationLicense,
  reinstateFederationLicense,
  getFederationLicenseKPIs,
} from '@/api/federationLicenses.js'

const emit = defineEmits(['navigate', 'message', 'error'])

const loading = ref(false)
const error = ref('')
const licenses = ref([])
const totalItems = ref(0)

const selectedLic = ref(null)
const extendModalLic = ref(null)
const revokeModalLic = ref(null)
const submittingAction = ref(false)

const kpis = reactive({
  total_licenses: 0,
  active_licenses: 0,
  extended_licenses: 0,
  revoked_licenses: 0,
  expiring_soon: 0,
})

const filters = reactive({
  search: '',
  status: '',
  license_type: '',
  sort_by: 'created_at',
  sort_order: 'desc',
})

const pagination = reactive({
  page: 1,
  per_page: 20,
})

const extendForm = reactive({
  new_expiry_date: '',
  reason: '',
})

const revokeForm = reactive({
  reason: '',
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
    const res = await getFederationLicenseKPIs()
    if (res) Object.assign(kpis, res)
  } catch (e) {
    console.error('Could not load license KPIs:', e)
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
      license_type: filters.license_type,
      sort_by: filters.sort_by,
      sort_order: filters.sort_order,
    }
    const res = await listFederationLicenses(params)
    const items = res.data || res.items || []
    licenses.value = Array.isArray(items) ? items : []
    totalItems.value = res.meta?.total || licenses.value.length
  } catch (err) {
    error.value = err.response?.data?.error?.message || err.message || 'Could not load federation licenses.'
    emit('error', error.value)
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
  filters.license_type = ''
  filters.sort_by = 'created_at'
  filters.sort_order = 'desc'
  pagination.page = 1
  loadData()
  loadKPIs()
}

function toggleSort(field) {
  if (filters.sort_by === field) {
    filters.sort_order = filters.sort_order === 'asc' ? 'desc' : 'asc'
  } else {
    filters.sort_by = field
    filters.sort_order = 'desc'
  }
  loadData()
}

function getSortIcon(field) {
  if (filters.sort_by !== field) return 'icofont-sort'
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

async function openDetailsModal(lic) {
  selectedLic.value = lic
  try {
    const full = await getFederationLicense(lic.id)
    if (full) selectedLic.value = full
  } catch (e) {
    console.error('Could not load full license details:', e)
  }
}

function openExtendModal(lic) {
  extendModalLic.value = lic
  const nextYear = new Date(lic.expiry_date || new Date())
  nextYear.setFullYear(nextYear.getFullYear() + 1)
  extendForm.new_expiry_date = nextYear.toISOString().slice(0, 10)
  extendForm.reason = 'Annual compliance audit cleared and statutory mandate extended.'
}

async function submitExtendLicense() {
  if (!extendModalLic.value) return
  submittingAction.value = true
  try {
    const updated = await extendFederationLicense(extendModalLic.value.id, extendForm)
    extendModalLic.value = null
    await Promise.all([loadData(), loadKPIs()])
    Swal.fire('License Extended', `Validity for ${updated.federation_name || 'federation'} extended to ${formatDate(updated.expiry_date)}.`, 'success')
  } catch (err) {
    Swal.fire('Extension Failed', err.response?.data?.error?.message || err.message, 'error')
  } finally {
    submittingAction.value = false
  }
}

function openRevokeModal(lic) {
  revokeModalLic.value = lic
  revokeForm.reason = ''
}

async function submitRevokeLicense() {
  if (!revokeModalLic.value) return
  submittingAction.value = true
  try {
    const updated = await revokeFederationLicense(revokeModalLic.value.id, revokeForm)
    revokeModalLic.value = null
    await Promise.all([loadData(), loadKPIs()])
    Swal.fire('License Revoked', `Recognition license for ${updated.federation_name || 'federation'} has been revoked.`, 'warning')
  } catch (err) {
    Swal.fire('Revocation Failed', err.response?.data?.error?.message || err.message, 'error')
  } finally {
    submittingAction.value = false
  }
}

async function openReinstateModal(lic) {
  const result = await Swal.fire({
    title: `Reinstate License for ${lic.federation_name}?`,
    text: 'This will restore the statutory recognition status for this national sports federation.',
    icon: 'question',
    showCancelButton: true,
    confirmButtonText: 'Yes, Reinstate',
    confirmButtonColor: '#198754',
  })

  if (result.isConfirmed) {
    try {
      const updated = await reinstateFederationLicense(lic.id, { reason: 'Compliance cleared; statutory license restored.' })
      await Promise.all([loadData(), loadKPIs()])
      Swal.fire('License Reinstated', `Statutory status restored for ${updated.federation_name}.`, 'success')
    } catch (err) {
      Swal.fire('Reinstatement Failed', err.response?.data?.error?.message || err.message, 'error')
    }
  }
}

function copyText(text) {
  if (!text) return
  navigator.clipboard?.writeText(text)
  Swal.fire({
    title: 'Copied!',
    text: `License ${text} copied to clipboard`,
    icon: 'success',
    timer: 1500,
    showConfirmButton: false,
  })
}

// ── EXPORT & CERTIFICATE PRINT ─────────────────────────────────────────────

function printCertificate(lic) {
  const printWin = window.open('', '_blank', 'width=1000,height=850')
  if (!printWin) {
    Swal.fire('Popup Blocked', 'Please allow popups to preview and print the certificate.', 'warning')
    return
  }

  const issueDateStr = formatDate(lic.issue_date)
  const expiryDateStr = formatDate(lic.expiry_date)
  const issuerName = lic.issued_by_name || 'Dr. Bernard Ogwel (General Secretary)'
  const safeFilename = `NCS_Federation_License_${String(lic.license_number || '').replace(/[^a-zA-Z0-9_-]/g, '_')}.pdf`

  const html = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>NCS Certificate of Recognition - ${lic.license_number}</title>
  ` + `<script src="https://cdnjs.cloudflare.com/ajax/libs/html2pdf.js/0.10.1/html2pdf.bundle.min.js"><` + `/script>` + `
  <style>
    @page { size: A4 landscape; margin: 8mm; }
    * { box-sizing: border-box; }
    body { font-family: "Georgia", "Times New Roman", serif; background: #f8fafc; margin: 0; padding: 20px; color: #1e293b; }
    .cert-action-bar { display: flex; justify-content: center; gap: 12px; margin-bottom: 20px; }
    .btn-cert { display: inline-flex; align-items: center; gap: 6px; padding: 8px 18px; border-radius: 6px; font-size: 13px; font-weight: bold; border: none; cursor: pointer; transition: all 0.15s ease; }
    .btn-primary { background: #1e3a8a; color: #fff; }
    .btn-primary:hover { background: #1e40af; }
    .btn-secondary { background: #0284c7; color: #fff; }
    .btn-secondary:hover { background: #0369a1; }
    .btn-close-cert { background: #e2e8f0; color: #475569; }
    .btn-close-cert:hover { background: #cbd5e1; }
    .cert-frame { border: 8px double #b45309; padding: 30px 40px; background: #fff; text-align: center; border-radius: 4px; box-shadow: 0 0 20px rgba(0,0,0,0.05); max-width: 960px; margin: 0 auto; }
    .logo-row { margin-bottom: 10px; }
    .logo-row img { max-height: 65px; }
    .republic-title { font-size: 15px; font-weight: bold; letter-spacing: 2px; text-transform: uppercase; color: #b45309; margin: 0; }
    .ncs-title { font-size: 24px; font-weight: bold; color: #0f172a; margin: 4px 0 14px; text-transform: uppercase; letter-spacing: 1px; }
    .cert-heading { font-size: 19px; font-style: italic; color: #475569; margin: 0 0 8px; }
    .cert-body { font-size: 14px; color: #334155; margin: 0 auto 14px; max-width: 700px; line-height: 1.5; }
    .fed-name { font-size: 26px; font-weight: bold; color: #1e3a8a; margin: 8px 0; text-decoration: underline; text-underline-offset: 5px; }
    .reg-tag { font-size: 13px; color: #64748b; margin-bottom: 14px; }
    .meta-box { display: flex; justify-content: space-around; margin: 18px auto; max-width: 680px; background: #fefce8; border: 1px solid #fef08a; padding: 10px; border-radius: 6px; }
    .meta-item strong { display: block; font-size: 13px; color: #713f12; }
    .meta-item span { font-size: 10px; color: #854d0e; text-transform: uppercase; }
    .conditions { font-size: 11px; font-style: italic; color: #64748b; margin: 12px auto; max-width: 650px; }
    .signatures { display: flex; justify-content: space-between; margin-top: 36px; padding: 0 30px; }
    .sig-line { width: 230px; border-top: 1px solid #334155; padding-top: 6px; font-size: 12px; font-weight: bold; text-align: center; }
    .sig-title { font-size: 10px; color: #64748b; font-weight: normal; }
    @media print {
      body { background: #fff; padding: 0; }
      .no-print { display: none !important; }
      .cert-frame { border: 8px double #b45309; box-shadow: none; max-width: 100%; }
    }
  </style>
</head>
<body>
  <div class="no-print cert-action-bar">
    <button class="btn-cert btn-primary" onclick="downloadPDF()">📥 Download PDF</button>
    <button class="btn-cert btn-secondary" onclick="window.print()">🖨️ Print Certificate</button>
    <button class="btn-cert btn-close-cert" onclick="window.close()">Close</button>
  </div>

  <div id="cert-to-print" class="cert-frame">
    <div class="logo-row">
      <img src="/main-logo.png" alt="National Council of Sports" />
    </div>
    <div class="republic-title">Republic of Uganda</div>
    <div class="ncs-title">National Council of Sports</div>
    <div class="cert-heading">Certificate of Statutory Recognition & Licensing</div>
    
    <div class="cert-body">
      This is to certify that under the provisions of the <strong>National Sports Act, 2023</strong>, the national sports governing body:
    </div>

    <div class="fed-name">${lic.federation_name}</div>
    <div class="reg-tag">Registration Number: <strong>${lic.federation_reg_no || 'NCS-STATUTORY'}</strong> &middot; Category: <strong>${lic.category || 'National Sports Federation'}</strong></div>

    <div class="meta-box">
      <div class="meta-item">
        <span>License Number</span>
        <strong>${lic.license_number}</strong>
      </div>
      <div class="meta-item">
        <span>Issue Date</span>
        <strong>${issueDateStr}</strong>
      </div>
      <div class="meta-item">
        <span>Valid Until</span>
        <strong>${expiryDateStr}</strong>
      </div>
      <div class="meta-item">
        <span>Status</span>
        <strong style="color: ${lic.status === 'REVOKED' ? '#dc2626' : '#16a34a'};">${lic.status}</strong>
      </div>
    </div>

    <div class="conditions">
      ${lic.conditions || 'Granted subject to compliance with the National Sports Act 2023, anti-doping protocols, and financial transparency regulations.'}
    </div>

    <div class="signatures">
      <div class="sig-line">
        ${issuerName}<br>
        <span class="sig-title">General Secretary &middot; National Council of Sports</span>
      </div>
      <div class="sig-line">
        ${lic.extended_by_name ? lic.extended_by_name + '<br><span class="sig-title">Extended Authority &middot; NCS Directorate</span>' : 'Board Chairman / Technical Director<br><span class="sig-title">National Council of Sports</span>'}
      </div>
    </div>
  </div>

  ` + `<script>
    function downloadPDF() {
      var element = document.getElementById('cert-to-print');
      var opt = {
        margin: [6, 6, 6, 6],
        filename: '${safeFilename}',
        image: { type: 'jpeg', quality: 0.98 },
        html2canvas: { scale: 2, useCORS: true },
        jsPDF: { unit: 'mm', format: 'a4', orientation: 'landscape' }
      };
      if (window.html2pdf) {
        window.html2pdf().set(opt).from(element).save();
      } else {
        window.print();
      }
    }
  <` + `/script>
</body>
</html>`

  printWin.document.write(html)
  printWin.document.close()
}

function exportToExcel() {
  const rows = [
    ['License Number', 'Federation Name', 'Acronym', 'Registration No', 'Type', 'Category', 'Issue Date', 'Expiry Date', 'Status', 'President', 'General Secretary', 'Issued By', 'Conditions'],
    ...licenses.value.map(l => [
      l.license_number,
      l.federation_name,
      l.federation_acronym,
      l.federation_reg_no,
      l.license_type,
      l.category,
      formatDate(l.issue_date),
      formatDate(l.expiry_date),
      l.status,
      l.president_name,
      l.secretary_name,
      l.issued_by_name,
      (l.conditions || '').replaceAll('\n', ' '),
    ]),
  ]

  const csvContent = '\xEF\xBB\xBF' + rows.map(r => r.map(cell => `"${String(cell || '').replaceAll('"', '""')}"`).join(',')).join('\n')
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `NCS_Federation_Licenses_${new Date().toISOString().slice(0, 10)}.csv`
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}

function exportToPdf() {
  const printWin = window.open('', '_blank', 'width=1000,height=800')
  if (!printWin) return

  const rowsHtml = licenses.value
    .map(
      (l, idx) => `
      <tr>
        <td>${idx + 1}</td>
        <td><strong>${l.license_number}</strong></td>
        <td><strong>${l.federation_name}</strong><br><small style="color:#666">${l.federation_acronym ? '(' + l.federation_acronym + ')' : ''} Reg: ${l.federation_reg_no || ''}</small></td>
        <td>${formatTypeName(l.license_type)}<br><small style="color:#666">${l.category}</small></td>
        <td>${formatDate(l.issue_date)}</td>
        <td>${formatDate(l.expiry_date)}</td>
        <td><span class="badge status-${String(l.status || '').toLowerCase()}">${l.status}</span></td>
      </tr>
    `
    )
    .join('')

  const html = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>NCS Federation Licenses Registry Summary</title>
  <style>
    @page { size: A4 landscape; margin: 12mm; }
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; color: #222; margin: 0; padding: 20px; font-size: 12px; }
    .header { text-align: center; border-bottom: 2px solid #0d6efd; padding-bottom: 12px; margin-bottom: 16px; }
    .header h2 { margin: 0 0 4px; font-size: 20px; }
    table { width: 100%; border-collapse: collapse; margin-top: 10px; }
    th { background: #f1f5f9; text-align: left; padding: 8px 6px; font-size: 11px; border-bottom: 2px solid #cbd5e1; }
    td { padding: 7px 6px; border-bottom: 1px solid #e2e8f0; font-size: 11px; vertical-align: top; }
    .badge { display: inline-block; padding: 3px 6px; border-radius: 4px; font-size: 10px; font-weight: bold; }
    .status-active { background: #d1fae5; color: #065f46; }
    .status-extended { background: #e0f2fe; color: #0369a1; }
    .status-revoked { background: #fee2e2; color: #991b1b; }
    .status-expired { background: #fef3c7; color: #92400e; }
  </style>
</head>
<body>
  <div class="header">
    <h2>NATIONAL COUNCIL OF SPORTS (NCS)</h2>
    <p>National Sports Federation Licensing & Statutory Accreditation Ledger</p>
  </div>
  <table>
    <thead>
      <tr>
        <th>#</th>
        <th>License Number</th>
        <th>Federation</th>
        <th>Type & Category</th>
        <th>Issue Date</th>
        <th>Expiry Date</th>
        <th>Status</th>
      </tr>
    </thead>
    <tbody>
      ${rowsHtml}
    </tbody>
  </table>
  <script>
    window.onload = function() { setTimeout(function() { window.print(); }, 400); };
  <\/script>
</body>
</html>`

  printWin.document.write(html)
  printWin.document.close()
}

// ── UTILITIES ──────────────────────────────────────────────────────────────

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

function formatFullDateTime(val) {
  if (!val) return '-'
  try {
    return new Date(val).toLocaleString('en-UG', {
      dateStyle: 'medium',
      timeStyle: 'short',
    })
  } catch {
    return String(val)
  }
}

function formatTypeName(t) {
  const m = String(t || '').toUpperCase()
  if (m === 'FULL_RECOGNITION') return 'Full Statutory Recognition'
  if (m === 'PROVISIONAL') return 'Provisional Recognition'
  if (m === 'ANNUAL_COMPLIANCE') return 'Annual Compliance'
  if (m === 'SPECIAL_CLEARANCE') return 'Special Clearance'
  return m
}

function getTypeBadgeClass(t) {
  const m = String(t || '').toUpperCase()
  if (m === 'FULL_RECOGNITION') return 'badge-full-rec'
  if (m === 'PROVISIONAL') return 'badge-provisional'
  return 'badge-general'
}

function getStatusClass(s) {
  const st = String(s || '').toUpperCase()
  if (st === 'ACTIVE') return 'status-active'
  if (st === 'EXTENDED') return 'status-extended'
  if (st === 'REVOKED') return 'status-revoked'
  if (st === 'EXPIRED') return 'status-expired'
  return 'status-suspended'
}

function getStatusIcon(s) {
  const st = String(s || '').toUpperCase()
  if (st === 'ACTIVE') return 'icofont-check-circled'
  if (st === 'EXTENDED') return 'icofont-history'
  if (st === 'REVOKED') return 'icofont-ban'
  if (st === 'EXPIRED') return 'icofont-clock-time'
  return 'icofont-warning'
}

function getValidityTagText(lic) {
  if (lic.status === 'REVOKED') return 'Revoked'
  if (!lic.expiry_date) return 'No Expiry'
  const diffDays = Math.ceil((new Date(lic.expiry_date) - new Date()) / (1000 * 60 * 60 * 24))
  if (diffDays < 0) return 'Expired'
  if (diffDays <= 60) return `${diffDays} days remaining`
  return 'Valid'
}

function getValidityTagClass(lic) {
  if (lic.status === 'REVOKED') return 'tag-danger'
  const diffDays = Math.ceil((new Date(lic.expiry_date) - new Date()) / (1000 * 60 * 60 * 24))
  if (diffDays < 0) return 'tag-danger'
  if (diffDays <= 60) return 'tag-warning'
  return 'tag-success'
}

function formatLogAction(action) {
  const a = String(action || '').toUpperCase()
  if (a === 'CREATED') return 'License Created & Issued'
  if (a === 'EXTENDED') return 'License Validity Extended'
  if (a === 'REVOKED') return 'License Revoked'
  if (a === 'REINSTATED') return 'License Reinstated'
  return a
}

function getTimelineDotClass(action) {
  const a = String(action || '').toUpperCase()
  if (a === 'CREATED') return 'dot-created'
  if (a === 'EXTENDED') return 'dot-extended'
  if (a === 'REVOKED') return 'dot-revoked'
  if (a === 'REINSTATED') return 'dot-reinstated'
  return 'dot-general'
}
</script>

<style scoped>
.manage-licenses-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
}

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

/* KPI Cards */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
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
}

.kpi-icon-wrap {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  flex-shrink: 0;
}

.kpi-label {
  font-size: 11px;
  color: #64748b;
  text-transform: uppercase;
  font-weight: 600;
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
}

.search-field input {
  width: 100%;
  padding: 8px 36px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 13px;
}

.clear-btn {
  position: absolute;
  right: 10px;
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
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
}

.filter-item select {
  padding: 8px 10px;
  font-size: 13px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
}

/* Table Card */
.table-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.licenses-table th {
  background: #f8fafc;
  color: #475569;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  padding: 12px 16px;
  border-bottom: 2px solid #e2e8f0;
}

.licenses-table td {
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
}

.license-code {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #1e293b;
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
}

.copy-btn:hover {
  color: #3b82f6;
}

.fed-avatar {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  color: #0f172a;
  font-size: 11px;
  flex-shrink: 0;
  overflow: hidden;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.badge-acronym {
  background: #e2e8f0;
  color: #334155;
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 4px;
  margin-left: 4px;
}

.type-badge {
  display: inline-block;
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}

.badge-full-rec {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}

.badge-provisional {
  background: #fefce8;
  color: #a16207;
  border: 1px solid #fef08a;
}

.badge-general {
  background: #f8fafc;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.validity-tag {
  display: inline-block;
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 4px;
  font-weight: 600;
  margin-top: 2px;
}

.tag-success {
  background: #f0fdf4;
  color: #16a34a;
}

.tag-warning {
  background: #fffbeb;
  color: #b45309;
}

.tag-danger {
  background: #fef2f2;
  color: #dc2626;
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
}

.status-active {
  background: #dcfce7;
  color: #15803d;
}

.status-extended {
  background: #e0f2fe;
  color: #0369a1;
}

.status-revoked {
  background: #fee2e2;
  color: #b91c1c;
}

.status-expired {
  background: #fef3c7;
  color: #b45309;
}

.lic-modal-backdrop {
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(3px);
}

/* Timeline */
.timeline-container {
  display: flex;
  flex-direction: column;
  gap: 14px;
  position: relative;
  padding-left: 20px;
  border-left: 2px solid #e2e8f0;
}

.timeline-item {
  position: relative;
}

.timeline-dot {
  position: absolute;
  left: -26px;
  top: 4px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #3b82f6;
  border: 2px solid #fff;
}

.dot-created { background: #10b981; }
.dot-extended { background: #0ea5e9; }
.dot-revoked { background: #ef4444; }
.dot-reinstated { background: #10b981; }
</style>
