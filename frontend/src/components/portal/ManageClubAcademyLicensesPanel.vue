<template>
  <section class="cms-panel manage-club-licenses-panel">
    <!-- Header -->
    <div class="cms-panel-head d-flex align-items-center justify-content-between flex-wrap gap-3">
      <div>
        <h2 class="mb-1"><i class="icofont-certificate-alt-1 text-primary me-2"></i> Club & Academy Licenses</h2>
        <p class="text-muted small mb-0">Manage statutory licenses, compliance records, lifecycle status, and immutable audit logs.</p>
      </div>
      <div class="d-flex align-items-center gap-2 flex-wrap">
        <button type="button" class="btn btn-outline-success btn-sm" @click="exportCSV">
          <i class="icofont-file-excel me-1"></i> Export CSV / Excel
        </button>
        <button type="button" class="btn btn-outline-danger btn-sm" @click="printAuditReport">
          <i class="icofont-printer me-1"></i> Print Directory
        </button>
        <button type="button" class="btn btn-primary btn-sm" @click="$emit('navigate', 'create-club-academy-licenses')">
          <i class="icofont-plus-circle me-1"></i> Issue New License
        </button>
      </div>
    </div>

    <!-- KPI Metric Summary Cards -->
    <div class="row g-3 mt-1 mb-4">
      <div class="col-xl-2 col-md-4 col-6">
        <div class="card border-0 shadow-sm kpi-card bg-primary-subtle text-primary h-100">
          <div class="card-body p-3">
            <span class="text-uppercase fw-bold small text-muted d-block" style="font-size: 11px;">Total Licenses</span>
            <h3 class="mb-0 fw-bold mt-1 text-primary">{{ kpis.total_licenses || 0 }}</h3>
            <small class="text-muted">Statutory Records</small>
          </div>
        </div>
      </div>
      <div class="col-xl-2 col-md-4 col-6">
        <div class="card border-0 shadow-sm kpi-card bg-success-subtle text-success h-100">
          <div class="card-body p-3">
            <span class="text-uppercase fw-bold small text-muted d-block" style="font-size: 11px;">Active</span>
            <h3 class="mb-0 fw-bold mt-1 text-success">{{ kpis.active_licenses || 0 }}</h3>
            <small class="text-success">Fully Accredited</small>
          </div>
        </div>
      </div>
      <div class="col-xl-2 col-md-4 col-6">
        <div class="card border-0 shadow-sm kpi-card bg-info-subtle text-info h-100">
          <div class="card-body p-3">
            <span class="text-uppercase fw-bold small text-muted d-block" style="font-size: 11px;">Extended</span>
            <h3 class="mb-0 fw-bold mt-1 text-info">{{ kpis.extended_licenses || 0 }}</h3>
            <small class="text-info">Prolonged Validity</small>
          </div>
        </div>
      </div>
      <div class="col-xl-2 col-md-4 col-6">
        <div class="card border-0 shadow-sm kpi-card bg-danger-subtle text-danger h-100">
          <div class="card-body p-3">
            <span class="text-uppercase fw-bold small text-muted d-block" style="font-size: 11px;">Revoked</span>
            <h3 class="mb-0 fw-bold mt-1 text-danger">{{ kpis.revoked_licenses || 0 }}</h3>
            <small class="text-danger">Suspended</small>
          </div>
        </div>
      </div>
      <div class="col-xl-4 col-md-8 col-12">
        <div class="card border-0 shadow-sm kpi-card bg-warning-subtle text-warning-emphasis h-100">
          <div class="card-body p-3">
            <span class="text-uppercase fw-bold small text-muted d-block" style="font-size: 11px;">Expiring Soon (&lt; 30 Days)</span>
            <h3 class="mb-0 fw-bold mt-1 text-warning-emphasis">{{ kpis.expiring_soon || 0 }}</h3>
            <small class="text-muted">Requires Renewal Assessment</small>
          </div>
        </div>
      </div>
    </div>

    <!-- Filters & Search Toolbar -->
    <div class="card border-0 shadow-sm mb-4">
      <div class="card-body p-3">
        <div class="row g-2 align-items-center">
          <div class="col-md-4">
            <div class="input-group">
              <span class="input-group-text bg-white"><i class="icofont-search-1"></i></span>
              <input
                v-model="filters.search"
                type="search"
                class="form-control"
                placeholder="Search License No, Academy No, Name, Official..."
                @input="debouncedLoad"
              />
            </div>
          </div>

          <div class="col-md-2">
            <select v-model="filters.status" class="form-select form-control" @change="loadLicenses(1)">
              <option value="">All Statuses</option>
              <option value="ACTIVE">Active</option>
              <option value="EXTENDED">Extended</option>
              <option value="REVOKED">Revoked</option>
              <option value="EXPIRED">Expired</option>
            </select>
          </div>

          <div class="col-md-2">
            <select v-model="filters.category" class="form-select form-control" @change="loadLicenses(1)">
              <option value="">All Categories</option>
              <option value="SPORTS_ACADEMY">Sports Academy</option>
              <option value="YOUTH_ACADEMY">Youth Academy</option>
              <option value="DEVELOPMENT_CENTRE">Development Centre</option>
              <option value="COMMUNITY_CLUB">Community Club</option>
              <option value="ELITE_ACADEMY">Elite Academy</option>
              <option value="SENIOR_CLUB">Senior Club</option>
            </select>
          </div>

          <div class="col-md-2">
            <input v-model="filters.start_date" type="date" class="form-control" title="Issued after" @change="loadLicenses(1)" />
          </div>

          <div class="col-md-2 d-flex gap-2">
            <button type="button" class="btn btn-outline-secondary w-100" title="Reset Filters" @click="resetFilters">
              <i class="icofont-refresh"></i> Reset
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Table of Licenses -->
    <div class="card border-0 shadow-sm">
      <div class="card-body p-0">
        <div class="table-responsive">
          <table class="table table-hover align-middle mb-0">
            <thead class="table-light">
              <tr>
                <th class="ps-3" @click="sortBy('license_number')" style="cursor: pointer;">
                  License Number <i class="icofont-sort ms-1"></i>
                </th>
                <th @click="sortBy('club_number')" style="cursor: pointer;">
                  Academy / Club No <i class="icofont-sort ms-1"></i>
                </th>
                <th @click="sortBy('club_name')" style="cursor: pointer;">
                  Sports Club / Academy <i class="icofont-sort ms-1"></i>
                </th>
                <th>Category</th>
                <th @click="sortBy('issue_date')" style="cursor: pointer;">
                  Issue Date <i class="icofont-sort ms-1"></i>
                </th>
                <th @click="sortBy('expiry_date')" style="cursor: pointer;">
                  Valid Until <i class="icofont-sort ms-1"></i>
                </th>
                <th @click="sortBy('status')" style="cursor: pointer;">
                  Status <i class="icofont-sort ms-1"></i>
                </th>
                <th class="text-end pe-3">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="loading">
                <td colspan="8" class="text-center py-5 text-muted">
                  <i class="icofont-spinner-alt-3 icofont-spin fs-4 d-block mb-2"></i>
                  Loading club & academy licenses...
                </td>
              </tr>
              <tr v-else-if="!licenses.length">
                <td colspan="8" class="text-center py-5 text-muted">
                  <i class="icofont-folder-open fs-2 d-block mb-2 text-muted opacity-50"></i>
                  No club & academy licenses found matching the selected filters.
                </td>
              </tr>
              <tr v-for="item in licenses" :key="item.id">
                <!-- License Number -->
                <td class="ps-3">
                  <span class="font-monospace fw-bold text-dark">{{ item.license_number }}</span>
                  <small v-if="item.license_type" class="text-muted d-block" style="font-size: 10px;">
                    {{ formatLicType(item.license_type) }}
                  </small>
                </td>

                <!-- Academy / Club No -->
                <td>
                  <span class="badge bg-light text-dark border font-monospace" style="font-size: 11px;">
                    {{ item.club_number || ('NCS-ACA-' + (item.club_id || '').substring(0, 4).toUpperCase()) }}
                  </span>
                </td>

                <!-- Club / Academy Name -->
                <td>
                  <div class="fw-bold text-dark">{{ item.club_name || 'Sports Club / Academy' }}</div>
                  <small class="text-muted">
                    <span v-if="item.club_acronym" class="badge bg-secondary-subtle text-secondary me-1">{{ item.club_acronym }}</span>
                    <span>{{ item.club_district || item.club_region || 'National' }}</span>
                    <span v-if="item.federation_name" class="ms-1">· {{ item.federation_name }}</span>
                  </small>
                </td>

                <!-- Category -->
                <td>
                  <span class="badge bg-info-subtle text-info text-uppercase" style="font-size: 10px; font-weight: 700;">
                    {{ formatCategory(item.category) }}
                  </span>
                </td>

                <!-- Issue Date -->
                <td>{{ formatDate(item.issue_date) }}</td>

                <!-- Expiry Date -->
                <td>
                  <div :class="getExpiryTextClass(item)">
                    <strong>{{ formatDate(item.expiry_date) }}</strong>
                  </div>
                  <small v-if="item.extended_at" class="text-info d-block" style="font-size: 10px;">
                    Extended on {{ formatDate(item.extended_at) }}
                  </small>
                </td>

                <!-- Status Badge -->
                <td>
                  <span class="badge" :class="getStatusBadgeClass(item.status)" style="font-size: 11px; font-weight: 700;">
                    {{ item.status }}
                  </span>
                </td>

                <!-- Actions -->
                <td class="text-end pe-3">
                  <div class="btn-group btn-group-sm">
                    <button
                      type="button"
                      class="btn btn-outline-primary"
                      title="View Details & Audit Log"
                      @click="viewDetails(item)"
                    >
                      <i class="icofont-eye"></i>
                    </button>
                    <button
                      type="button"
                      class="btn btn-outline-secondary"
                      title="Print Certificate"
                      @click="printCertificate(item)"
                    >
                      <i class="icofont-printer"></i>
                    </button>
                    <button
                      type="button"
                      class="btn btn-outline-info"
                      title="Extend License Validity"
                      @click="openExtendModal(item)"
                    >
                      <i class="icofont-calendar"></i>
                    </button>
                    <button
                      v-if="item.status !== 'REVOKED'"
                      type="button"
                      class="btn btn-outline-danger"
                      title="Revoke License"
                      @click="openRevokeModal(item)"
                    >
                      <i class="icofont-ban"></i>
                    </button>
                    <button
                      v-else
                      type="button"
                      class="btn btn-outline-success"
                      title="Reinstate License"
                      @click="openReinstateModal(item)"
                    >
                      <i class="icofont-check-circled"></i>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="d-flex align-items-center justify-content-between p-3 border-top">
          <span class="small text-muted">
            Showing page {{ currentPage }} of {{ totalPages }} ({{ totalCount }} total licenses)
          </span>
          <div class="btn-group btn-group-sm">
            <button type="button" class="btn btn-outline-secondary" :disabled="currentPage <= 1" @click="loadLicenses(currentPage - 1)">
              Previous
            </button>
            <button
              v-for="p in paginationPages"
              :key="p"
              type="button"
              class="btn"
              :class="p === currentPage ? 'btn-primary' : 'btn-outline-secondary'"
              @click="loadLicenses(p)"
            >
              {{ p }}
            </button>
            <button type="button" class="btn btn-outline-secondary" :disabled="currentPage >= totalPages" @click="loadLicenses(currentPage + 1)">
              Next
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- DETAILS & IMMUTABLE AUDIT TIMELINE MODAL -->
    <div v-if="selectedLic" class="modal-backdrop fade show" @click.self="selectedLic = null">
      <div class="modal d-block" tabindex="-1" style="background: rgba(0,0,0,0.5);">
        <div class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable">
          <div class="modal-content">
            <div class="modal-header bg-light">
              <div>
                <h5 class="modal-title fw-bold">
                  <i class="icofont-certificate-alt-1 text-primary me-2"></i>
                  License Details: {{ selectedLic.license_number }}
                </h5>
                <span class="small text-muted">{{ selectedLic.club_name }} ({{ selectedLic.club_number || 'N/A' }})</span>
              </div>
              <button type="button" class="btn-close" @click="selectedLic = null"></button>
            </div>
            <div class="modal-body">
              <!-- Overview Grid -->
              <div class="row g-3 mb-4">
                <div class="col-md-6">
                  <div class="p-3 bg-light rounded">
                    <span class="text-muted small d-block">Sports Club / Academy</span>
                    <strong class="fs-6">{{ selectedLic.club_name }}</strong>
                    <div class="small text-muted mt-1">
                      Number: <strong>{{ selectedLic.club_number || 'N/A' }}</strong> &middot; Acronym: <strong>{{ selectedLic.club_acronym || '-' }}</strong>
                    </div>
                    <div class="small text-muted">
                      Affiliated Federation: <strong>{{ selectedLic.federation_name || 'NCS Statutory' }}</strong>
                    </div>
                  </div>
                </div>

                <div class="col-md-6">
                  <div class="p-3 bg-light rounded">
                    <span class="text-muted small d-block">License Status</span>
                    <span class="badge fs-6 mt-1" :class="getStatusBadgeClass(selectedLic.status)">{{ selectedLic.status }}</span>
                    <div class="small text-muted mt-1">
                      Issue Date: <strong>{{ formatDate(selectedLic.issue_date) }}</strong> &middot; Expiry: <strong>{{ formatDate(selectedLic.expiry_date) }}</strong>
                    </div>
                  </div>
                </div>

                <div class="col-12" v-if="selectedLic.conditions">
                  <div class="p-3 border rounded">
                    <span class="text-muted small d-block fw-bold mb-1">Statutory & Operational Conditions:</span>
                    <p class="mb-0 small text-dark fst-italic">{{ selectedLic.conditions }}</p>
                  </div>
                </div>

                <div class="col-12" v-if="selectedLic.status === 'REVOKED'">
                  <div class="alert alert-danger mb-0">
                    <strong><i class="icofont-warning me-1"></i> Revocation Reason:</strong>
                    {{ selectedLic.revocation_reason }}
                    <div class="small text-muted mt-1">
                      Revoked by {{ selectedLic.revoked_by_name || 'Administrator' }} on {{ formatDateTime(selectedLic.revoked_at) }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Immutable Audit Log Timeline -->
              <h6 class="fw-bold border-bottom pb-2 mb-3">
                <i class="icofont-history text-primary me-2"></i> Immutable Audit History & Action Log
              </h6>

              <div v-if="!selectedLic.logs || !selectedLic.logs.length" class="text-muted small py-3 text-center">
                No historical audit events logged yet.
              </div>
              <div v-else class="audit-timeline">
                <div v-for="log in selectedLic.logs" :key="log.id" class="timeline-item pb-3 mb-3 border-bottom position-relative ps-4">
                  <div class="timeline-icon-dot" :class="getLogDotClass(log.action)"></div>
                  <div class="d-flex align-items-center justify-content-between">
                    <strong :class="getLogTextClass(log.action)">
                      {{ log.action }}
                    </strong>
                    <span class="small text-muted">{{ formatDateTime(log.created_at) }}</span>
                  </div>
                  <div class="small text-dark mt-1">
                    Performed by: <strong>{{ log.performed_by_name || 'System Administrator' }}</strong>
                  </div>
                  <div v-if="log.reason" class="small text-muted mt-1 fst-italic">
                    "{{ log.reason }}"
                  </div>
                  <div v-if="log.old_status || log.new_status" class="small text-muted mt-1">
                    Status: <span class="badge bg-secondary-subtle text-secondary">{{ log.old_status || 'NONE' }}</span> &rarr; <span class="badge" :class="getStatusBadgeClass(log.new_status)">{{ log.new_status }}</span>
                    <span v-if="log.new_expiry_date" class="ms-2">| New Expiry: <strong>{{ formatDate(log.new_expiry_date) }}</strong></span>
                  </div>
                </div>
              </div>
            </div>
            <div class="modal-footer bg-light">
              <button type="button" class="btn btn-outline-primary btn-sm" @click="printCertificate(selectedLic)">
                <i class="icofont-printer me-1"></i> Print Recognition Certificate
              </button>
              <button type="button" class="btn btn-secondary btn-sm" @click="selectedLic = null">Close</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- EXTEND MODAL -->
    <div v-if="extendModalItem" class="modal-backdrop fade show" @click.self="extendModalItem = null">
      <div class="modal d-block" tabindex="-1" style="background: rgba(0,0,0,0.5);">
        <div class="modal-dialog modal-dialog-centered">
          <div class="modal-content">
            <div class="modal-header bg-info text-white">
              <h5 class="modal-title fw-bold"><i class="icofont-calendar me-2"></i> Extend License Validity</h5>
              <button type="button" class="btn-close btn-close-white" @click="extendModalItem = null"></button>
            </div>
            <form @submit.prevent="submitExtend">
              <div class="modal-body">
                <p class="small text-muted mb-3">
                  Extending license <strong>{{ extendModalItem.license_number }}</strong> for <strong>{{ extendModalItem.club_name }}</strong>.
                  Current expiry date: <strong>{{ formatDate(extendModalItem.expiry_date) }}</strong>.
                </p>

                <div class="mb-3">
                  <label class="form-label fw-bold">New Expiry Date <span class="text-danger">*</span></label>
                  <input v-model="extendForm.new_expiry_date" type="date" class="form-control" required />
                </div>

                <div class="mb-3">
                  <label class="form-label fw-bold">Reason for Extension <span class="text-danger">*</span></label>
                  <textarea
                    v-model="extendForm.reason"
                    class="form-control"
                    rows="2"
                    required
                    placeholder="Enter compliance clearance or renewal justification..."
                  ></textarea>
                </div>

                <div class="mb-3">
                  <label class="form-label fw-bold">Additional Notes (Optional)</label>
                  <input v-model="extendForm.notes" class="form-control" placeholder="Optional audit memo" />
                </div>
              </div>
              <div class="modal-footer bg-light">
                <button type="button" class="btn btn-secondary" @click="extendModalItem = null">Cancel</button>
                <button type="submit" class="btn btn-info text-white fw-bold" :disabled="actionLoading">
                  {{ actionLoading ? 'Extending...' : 'Confirm Extension' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- REVOKE MODAL -->
    <div v-if="revokeModalItem" class="modal-backdrop fade show" @click.self="revokeModalItem = null">
      <div class="modal d-block" tabindex="-1" style="background: rgba(0,0,0,0.5);">
        <div class="modal-dialog modal-dialog-centered">
          <div class="modal-content">
            <div class="modal-header bg-danger text-white">
              <h5 class="modal-title fw-bold"><i class="icofont-ban me-2"></i> Revoke Statutory Recognition</h5>
              <button type="button" class="btn-close btn-close-white" @click="revokeModalItem = null"></button>
            </div>
            <form @submit.prevent="submitRevoke">
              <div class="modal-body">
                <div class="alert alert-warning small">
                  <strong>Warning:</strong> Revoking this license will immediately suspend the statutory operating recognition of <strong>{{ revokeModalItem.club_name }}</strong> under the National Sports Act 2023.
                </div>

                <div class="mb-3">
                  <label class="form-label fw-bold">Revocation Reason <span class="text-danger">*</span></label>
                  <textarea
                    v-model="revokeForm.reason"
                    class="form-control"
                    rows="3"
                    required
                    placeholder="Specify the violation, governance breach, or non-compliance ground..."
                  ></textarea>
                </div>

                <div class="mb-3">
                  <label class="form-label fw-bold">Internal Audit Notes</label>
                  <input v-model="revokeForm.notes" class="form-control" placeholder="Reference files or official notices" />
                </div>
              </div>
              <div class="modal-footer bg-light">
                <button type="button" class="btn btn-secondary" @click="revokeModalItem = null">Cancel</button>
                <button type="submit" class="btn btn-danger fw-bold" :disabled="actionLoading">
                  {{ actionLoading ? 'Revoking...' : 'Confirm Revocation' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- REINSTATE MODAL -->
    <div v-if="reinstateModalItem" class="modal-backdrop fade show" @click.self="reinstateModalItem = null">
      <div class="modal d-block" tabindex="-1" style="background: rgba(0,0,0,0.5);">
        <div class="modal-dialog modal-dialog-centered">
          <div class="modal-content">
            <div class="modal-header bg-success text-white">
              <h5 class="modal-title fw-bold"><i class="icofont-check-circled me-2"></i> Reinstate License</h5>
              <button type="button" class="btn-close btn-close-white" @click="reinstateModalItem = null"></button>
            </div>
            <form @submit.prevent="submitReinstate">
              <div class="modal-body">
                <p class="small text-muted mb-3">
                  Reinstating statutory license <strong>{{ reinstateModalItem.license_number }}</strong> for <strong>{{ reinstateModalItem.club_name }}</strong>.
                </p>

                <div class="mb-3">
                  <label class="form-label fw-bold">Reinstatement Justification <span class="text-danger">*</span></label>
                  <textarea
                    v-model="reinstateForm.reason"
                    class="form-control"
                    rows="3"
                    required
                    placeholder="Specify compliance rectification and resolution details..."
                  ></textarea>
                </div>
              </div>
              <div class="modal-footer bg-light">
                <button type="button" class="btn btn-secondary" @click="reinstateModalItem = null">Cancel</button>
                <button type="submit" class="btn btn-success fw-bold" :disabled="actionLoading">
                  {{ actionLoading ? 'Reinstating...' : 'Confirm Reinstatement' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import {
  listClubAcademyLicenses,
  getClubAcademyLicenseKPIs,
  getClubAcademyLicenseByID,
  extendClubAcademyLicense,
  revokeClubAcademyLicense,
  reinstateClubAcademyLicense,
} from '@/api/clubAcademyLicenses.js'

const emit = defineEmits(['navigate', 'message', 'error'])

const licenses = ref([])
const kpis = ref({})
const loading = ref(false)
const actionLoading = ref(false)

const currentPage = ref(1)
const perPage = ref(20)
const totalCount = ref(0)

const filters = reactive({
  search: '',
  status: '',
  category: '',
  start_date: '',
  end_date: '',
  sort_by: 'created_at',
  sort_order: 'desc',
})

const selectedLic = ref(null)
const extendModalItem = ref(null)
const extendForm = reactive({ new_expiry_date: '', reason: '', notes: '' })
const revokeModalItem = ref(null)
const revokeForm = reactive({ reason: '', notes: '' })
const reinstateModalItem = ref(null)
const reinstateForm = reactive({ reason: '', notes: '' })

let searchTimeout = null
function debouncedLoad() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    loadLicenses(1)
  }, 350)
}

const totalPages = computed(() => Math.ceil(totalCount.value / perPage.value) || 1)

const paginationPages = computed(() => {
  const pages = []
  const max = totalPages.value
  let start = Math.max(1, currentPage.value - 2)
  let end = Math.min(max, start + 4)
  if (end - start < 4) start = Math.max(1, end - 4)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

async function loadKPIs() {
  try {
    kpis.value = await getClubAcademyLicenseKPIs()
  } catch (err) {
    console.error('Failed to load KPIs:', err)
  }
}

async function loadLicenses(page = 1) {
  currentPage.value = page
  loading.value = true
  try {
    const params = {
      page,
      per_page: perPage.value,
      search: filters.search.trim(),
      status: filters.status,
      category: filters.category,
      start_date: filters.start_date,
      end_date: filters.end_date,
      sort_by: filters.sort_by,
      sort_order: filters.sort_order,
    }
    const res = await listClubAcademyLicenses(params)
    licenses.value = res.data || []
    totalCount.value = res.meta?.total || (res.data ? res.data.length : 0)
  } catch (err) {
    emit('error', err.response?.data?.error?.message || 'Failed to load licenses')
  } finally {
    loading.value = false
  }
}

function sortBy(col) {
  if (filters.sort_by === col) {
    filters.sort_order = filters.sort_order === 'asc' ? 'desc' : 'asc'
  } else {
    filters.sort_by = col
    filters.sort_order = 'desc'
  }
  loadLicenses(1)
}

function resetFilters() {
  filters.search = ''
  filters.status = ''
  filters.category = ''
  filters.start_date = ''
  filters.end_date = ''
  filters.sort_by = 'created_at'
  filters.sort_order = 'desc'
  loadLicenses(1)
}

async function viewDetails(item) {
  try {
    const full = await getClubAcademyLicenseByID(item.id)
    selectedLic.value = full
  } catch (err) {
    selectedLic.value = item
  }
}

function openExtendModal(item) {
  extendModalItem.value = item
  const currentExp = item.expiry_date ? new Date(item.expiry_date) : new Date()
  currentExp.setFullYear(currentExp.getFullYear() + 1)
  extendForm.new_expiry_date = currentExp.toISOString().substring(0, 10)
  extendForm.reason = 'Statutory compliance audit completed; validity prolonged.'
  extendForm.notes = ''
}

async function submitExtend() {
  if (!extendModalItem.value) return
  actionLoading.value = true
  try {
    await extendClubAcademyLicense(extendModalItem.value.id, extendForm)
    emit('message', `License ${extendModalItem.value.license_number} extended successfully.`)
    extendModalItem.value = null
    loadLicenses(currentPage.value)
    loadKPIs()
  } catch (err) {
    emit('error', err.response?.data?.error?.message || 'Failed to extend license')
  } finally {
    actionLoading.value = false
  }
}

function openRevokeModal(item) {
  revokeModalItem.value = item
  revokeForm.reason = ''
  revokeForm.notes = ''
}

async function submitRevoke() {
  if (!revokeModalItem.value) return
  actionLoading.value = true
  try {
    await revokeClubAcademyLicense(revokeModalItem.value.id, revokeForm)
    emit('message', `License ${revokeModalItem.value.license_number} revoked.`)
    revokeModalItem.value = null
    loadLicenses(currentPage.value)
    loadKPIs()
  } catch (err) {
    emit('error', err.response?.data?.error?.message || 'Failed to revoke license')
  } finally {
    actionLoading.value = false
  }
}

function openReinstateModal(item) {
  reinstateModalItem.value = item
  reinstateForm.reason = 'Compliance issues rectified and validated by NCS Sports Development Directorate.'
  reinstateForm.notes = ''
}

async function submitReinstate() {
  if (!reinstateModalItem.value) return
  actionLoading.value = true
  try {
    await reinstateClubAcademyLicense(reinstateModalItem.value.id, reinstateForm)
    emit('message', `License ${reinstateModalItem.value.license_number} reinstated to ACTIVE.`)
    reinstateModalItem.value = null
    loadLicenses(currentPage.value)
    loadKPIs()
  } catch (err) {
    emit('error', err.response?.data?.error?.message || 'Failed to reinstate license')
  } finally {
    actionLoading.value = false
  }
}

function printCertificate(item) {
  if (!item) return
  const printWin = window.open('', '_blank', 'width=1000,height=850')
  if (!printWin) return

  const issueDateStr = formatDate(item.issue_date)
  const expiryDateStr = formatDate(item.expiry_date)
  const issuerName = item.issued_by_name || 'Dr. Bernard Ogwel (General Secretary)'
  const safeFilename = `NCS_Academy_License_${String(item.license_number || '').replace(/[^a-zA-Z0-9_-]/g, '_')}.pdf`

  const html = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>NCS Certificate of Accreditation - ${item.license_number}</title>
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
    .cert-frame { border: 8px double #1e3a8a; padding: 30px 40px; background: #fff; text-align: center; border-radius: 4px; box-shadow: 0 0 20px rgba(0,0,0,0.05); max-width: 960px; margin: 0 auto; }
    .logo-row { margin-bottom: 10px; }
    .logo-row img { max-height: 65px; }
    .republic-title { font-size: 15px; font-weight: bold; letter-spacing: 2px; text-transform: uppercase; color: #1e3a8a; margin: 0; }
    .ncs-title { font-size: 24px; font-weight: bold; color: #0f172a; margin: 4px 0 14px; text-transform: uppercase; letter-spacing: 1px; }
    .cert-heading { font-size: 19px; font-style: italic; color: #475569; margin: 0 0 8px; }
    .cert-body { font-size: 14px; color: #334155; margin: 0 auto 14px; max-width: 700px; line-height: 1.5; }
    .academy-name { font-size: 26px; font-weight: bold; color: #0f172a; margin: 8px 0; text-decoration: underline; text-underline-offset: 5px; }
    .reg-tag { font-size: 13px; color: #64748b; margin-bottom: 14px; }
    .meta-box { display: flex; justify-content: space-around; margin: 18px auto; max-width: 680px; background: #eff6ff; border: 1px solid #bfdbfe; padding: 10px; border-radius: 6px; }
    .meta-item strong { display: block; font-size: 13px; color: #1e3a8a; }
    .meta-item span { font-size: 10px; color: #3b82f6; text-transform: uppercase; }
    .conditions { font-size: 11px; font-style: italic; color: #64748b; margin: 12px auto; max-width: 650px; }
    .signatures { display: flex; justify-content: space-between; margin-top: 36px; padding: 0 30px; }
    .sig-line { width: 230px; border-top: 1px solid #334155; padding-top: 6px; font-size: 12px; font-weight: bold; text-align: center; }
    .sig-title { font-size: 10px; color: #64748b; font-weight: normal; }
    @media print {
      body { background: #fff; padding: 0; }
      .no-print { display: none !important; }
      .cert-frame { border: 8px double #1e3a8a; box-shadow: none; max-width: 100%; }
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
    <div class="cert-heading">Certificate of Operating Accreditation & Licensing</div>
    
    <div class="cert-body">
      This is to certify that under the provisions of the <strong>National Sports Act, 2023</strong>, the sports organisation:
    </div>

    <div class="academy-name">${item.club_name}</div>
    <div class="reg-tag">Academy / Club ID: <strong>${item.club_number || item.club_id}</strong> &middot; Category: <strong>${formatCategory(item.category)}</strong></div>

    <div class="meta-box">
      <div class="meta-item">
        <span>License Number</span>
        <strong>${item.license_number}</strong>
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
        <strong style="color: ${item.status === 'REVOKED' ? '#dc2626' : '#16a34a'};">${item.status}</strong>
      </div>
    </div>

    <div class="conditions">
      ${item.conditions || 'Granted subject to compliance with the National Sports Act 2023, youth safeguarding protocols, and athlete registry standards.'}
    </div>

    <div class="signatures">
      <div class="sig-line">
        ${issuerName}<br>
        <span class="sig-title">General Secretary &middot; National Council of Sports</span>
      </div>
      <div class="sig-line">
        ${item.extended_by_name ? item.extended_by_name + '<br><span class="sig-title">Extended Authority &middot; NCS Directorate</span>' : 'Director Technical & Regulations<br><span class="sig-title">National Council of Sports</span>'}
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

function exportCSV() {
  if (!licenses.value.length) return
  const headers = ['License Number', 'Academy Number', 'Sports Club / Academy', 'Category', 'Issue Date', 'Expiry Date', 'Status', 'Federation', 'Conditions']
  const rows = licenses.value.map(l => [
    `"${l.license_number}"`,
    `"${l.club_number || l.club_id || ''}"`,
    `"${(l.club_name || '').replace(/"/g, '""')}"`,
    `"${l.category || ''}"`,
    `"${l.issue_date || ''}"`,
    `"${l.expiry_date || ''}"`,
    `"${l.status || ''}"`,
    `"${(l.federation_name || '').replace(/"/g, '""')}"`,
    `"${(l.conditions || '').replace(/"/g, '""')}"`,
  ])

  const csvContent = '\uFEFF' + [headers.join(','), ...rows.map(e => e.join(','))].join('\n')
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.setAttribute('download', `NCS_Club_Academy_Licenses_${new Date().toISOString().substring(0, 10)}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function printAuditReport() {
  window.print()
}

function formatDate(val) {
  if (!val) return '-'
  try {
    return new Date(val).toLocaleDateString('en-UG', { year: 'numeric', month: 'short', day: 'numeric' })
  } catch {
    return String(val)
  }
}

function formatDateTime(val) {
  if (!val) return '-'
  try {
    return new Date(val).toLocaleString('en-UG', {
      year: 'numeric', month: 'short', day: 'numeric',
      hour: '2-digit', minute: '2-digit'
    })
  } catch {
    return String(val)
  }
}

function formatLicType(type) {
  return String(type || '').replace(/_/g, ' ')
}

function formatCategory(cat) {
  return String(cat || '').replace(/_/g, ' ')
}

function getStatusBadgeClass(status) {
  const s = String(status || '').toUpperCase()
  if (s === 'ACTIVE') return 'bg-success text-white'
  if (s === 'EXTENDED') return 'bg-info text-white'
  if (s === 'REVOKED') return 'bg-danger text-white'
  if (s === 'EXPIRED') return 'bg-warning text-dark'
  return 'bg-secondary text-white'
}

function getExpiryTextClass(item) {
  if (item.status === 'REVOKED') return 'text-danger text-decoration-line-through'
  const exp = new Date(item.expiry_date).getTime()
  const now = Date.now()
  if (exp < now) return 'text-danger'
  if (exp - now < 30 * 86400000) return 'text-warning'
  return 'text-dark'
}

function getLogDotClass(action) {
  switch (action) {
    case 'CREATED': return 'bg-primary'
    case 'EXTENDED': return 'bg-info'
    case 'REVOKED': return 'bg-danger'
    case 'REINSTATED': return 'bg-success'
    default: return 'bg-secondary'
  }
}

function getLogTextClass(action) {
  switch (action) {
    case 'CREATED': return 'text-primary'
    case 'EXTENDED': return 'text-info'
    case 'REVOKED': return 'text-danger'
    case 'REINSTATED': return 'text-success'
    default: return 'text-secondary'
  }
}

onMounted(() => {
  loadKPIs()
  loadLicenses(1)
})
</script>

<style scoped>
.manage-club-licenses-panel {
  padding: 24px;
}
.kpi-card {
  border-radius: 10px;
  transition: transform 0.15s ease;
}
.kpi-card:hover {
  transform: translateY(-2px);
}
.timeline-icon-dot {
  position: absolute;
  left: 0;
  top: 4px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
}
.audit-timeline {
  max-height: 280px;
  overflow-y: auto;
}
</style>
