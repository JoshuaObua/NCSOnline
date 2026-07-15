<template>
  <div class="reports-container card border-0 shadow-sm">
    <div class="card-header d-flex flex-wrap align-items-center justify-content-between gap-3 bg-white py-3">
      <div>
        <h4 class="mb-0 text-primary fw-bold">Sports Registry Reports</h4>
        <span class="text-muted small">Select a registry module to view, filter, sort, and export operational datasets.</span>
      </div>
      <div class="d-flex align-items-center gap-2">
        <label for="registrySelect" class="visually-hidden">Select Registry</label>
        <select
          id="registrySelect"
          v-model="activeRegistry"
          class="form-select form-select-sm"
          style="min-width: 200px;"
          @change="loadRegistryData"
        >
          <option v-for="opt in registries" :key="opt.id" :value="opt.id">
            {{ opt.label }}
          </option>
        </select>
        <button
          type="button"
          class="btn btn-sm btn-success d-flex align-items-center gap-1"
          :disabled="loading || !items.length"
          @click="exportToCSV"
        >
          <i class="icofont-download"></i> Export CSV
        </button>
      </div>
    </div>

    <div class="card-body">
      <div class="d-flex flex-wrap gap-2 mb-3 align-items-center justify-content-between">
        <div class="search-box">
          <input
            v-model="searchQuery"
            type="search"
            class="form-control form-control-sm"
            placeholder="Search records..."
            style="max-width: 320px;"
            @input="onSearchInput"
          />
        </div>
        <div class="text-muted small">
          Showing {{ items.length }} of {{ totalCount }} records
        </div>
      </div>

      <div v-if="loading" class="text-center py-5">
        <i class="icofont-spinner animate-spin fs-2 text-primary d-inline-block"></i>
        <p class="mt-2 text-muted">Loading registry records...</p>
      </div>
      <div v-else-if="error" class="alert alert-danger">{{ error }}</div>
      <div v-else-if="!items.length" class="text-center py-5 text-muted">
        <i class="icofont-exclamation-circle fs-2 mb-2 d-inline-block"></i>
        <p>No records found matching filters.</p>
      </div>
      <div v-else class="table-responsive">
        <table class="table table-striped table-hover align-middle">
          <thead>
            <tr>
              <th
                v-for="col in columns"
                :key="col.key"
                class="sortable-header"
                :style="{ width: col.width || 'auto', cursor: 'pointer' }"
                @click="sortBy(col.key)"
              >
                <div class="d-flex align-items-center gap-1">
                  {{ col.label }}
                  <i
                    v-if="sortKey === col.key"
                    :class="sortOrder === 'asc' ? 'icofont-arrow-up' : 'icofont-arrow-down'"
                    class="text-primary"
                  ></i>
                  <i v-else class="icofont-expand-alt text-muted" style="opacity: 0.4;"></i>
                </div>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, idx) in sortedItems" :key="row.id || idx">
              <td v-for="col in columns" :key="col.key">
                <span v-if="col.type === 'boolean'">
                  <span
                    class="badge"
                    :class="row[col.key] ? 'bg-success-light text-success' : 'bg-danger-light text-danger'"
                  >
                    {{ row[col.key] ? 'Yes' : 'No' }}
                  </span>
                </span>
                <span v-else-if="col.type === 'date'">
                  {{ formatDate(row[col.key]) }}
                </span>
                <span v-else-if="col.type === 'money'">
                  UGX {{ formatMoney(row[col.key]) }}
                </span>
                <span v-else-if="col.type === 'json'">
                  <code class="small text-muted">{{ truncateJson(row[col.key]) }}</code>
                </span>
                <span v-else class="fw-medium text-dark">
                  {{ row[col.key] ?? '-' }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1 && !loading" class="d-flex align-items-center justify-content-between border-top pt-3 mt-3">
        <div class="text-muted small">
          Page {{ page }} of {{ totalPages }}
        </div>
        <nav aria-label="Registry pagination">
          <ul class="pagination pagination-sm mb-0">
            <li class="page-item" :class="{ disabled: page <= 1 }">
              <button class="page-link" type="button" @click="changePage(page - 1)">Previous</button>
            </li>
            <li
              v-for="p in totalPages"
              :key="p"
              class="page-item"
              :class="{ active: page === p }"
            >
              <button class="page-link" type="button" @click="changePage(p)">{{ p }}</button>
            </li>
            <li class="page-item" :class="{ disabled: page >= totalPages }">
              <button class="page-link" type="button" @click="changePage(page + 1)">Next</button>
            </li>
          </ul>
        </nav>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { listNsmisDomain } from '@/api/nsmis.js'

const registries = [
  { id: 'athletes', label: 'Athletes Registry' },
  { id: 'clubs', label: 'Clubs & Academies' },
  { id: 'coaches', label: 'Coaches Registry' },
  { id: 'competitions', label: 'Competitions Logs' },
  { id: 'competition-results', label: 'Competition Results' },
  { id: 'medals', label: 'Medal Standings' },
  { id: 'talent', label: 'Talent Records' },
  { id: 'national-team', label: 'National Team Caps' },
  { id: 'technical-officials', label: 'Technical Officials' },
  { id: 'medical-records', label: 'Medical Clearance Logs' },
  { id: 'safeguarding-records', label: 'Safeguarding Records' },
  { id: 'anti-doping', label: 'Anti-Doping Compliance' },
  { id: 'disbursements', label: 'Disbursements' },
  { id: 'accountabilities', label: 'Financial Accountabilities' },
  { id: 'equipment', label: 'Distributed Equipment' }
]

const columnsConfig = {
  athletes: [
    { key: 'athlete_number', label: 'Athlete ID' },
    { key: 'full_name', label: 'Name' },
    { key: 'gender', label: 'Gender' },
    { key: 'date_of_birth', label: 'DOB', type: 'date' },
    { key: 'discipline', label: 'Discipline' },
    { key: 'national_team_status', label: 'National Status' },
    { key: 'age_category', label: 'Category' },
    { key: 'status', label: 'License Status' }
  ],
  clubs: [
    { key: 'name', label: 'Club Name' },
    { key: 'acronym', label: 'Acronym' },
    { key: 'federation_id', label: 'Federation ID' },
    { key: 'contact_person', label: 'Contact' },
    { key: 'email', label: 'Email' },
    { key: 'district', label: 'District' },
    { key: 'region', label: 'Region' },
    { key: 'status', label: 'Status' }
  ],
  coaches: [
    { key: 'full_name', label: 'Coach' },
    { key: 'federation_id', label: 'Federation' },
    { key: 'certification_level', label: 'Level' },
    { key: 'license_number', label: 'License' },
    { key: 'email', label: 'Email' },
    { key: 'status', label: 'Status' }
  ],
  competitions: [
    { key: 'name', label: 'Competition' },
    { key: 'venue', label: 'Venue' },
    { key: 'level', label: 'Level' },
    { key: 'starts_on', label: 'Starts', type: 'date' },
    { key: 'ends_on', label: 'Ends', type: 'date' },
    { key: 'status', label: 'Status' }
  ],
  'competition-results': [
    { key: 'athlete_name', label: 'Athlete' },
    { key: 'event', label: 'Event' },
    { key: 'position', label: 'Position' },
    { key: 'time_result', label: 'Time' },
    { key: 'distance_result', label: 'Distance' },
    { key: 'is_national_record', label: 'Nat. Record', type: 'boolean' }
  ],
  medals: [
    { key: 'athlete_name', label: 'Athlete' },
    { key: 'event', label: 'Event' },
    { key: 'medal_type', label: 'Medal' },
    { key: 'won_on', label: 'Won On', type: 'date' },
    { key: 'prize_money', label: 'Prize Money', type: 'money' }
  ],
  talent: [
    { key: 'athlete_name', label: 'Athlete' },
    { key: 'age_at_identification', label: 'Age Identified' },
    { key: 'district', label: 'District' },
    { key: 'identified_on', label: 'Identified On', type: 'date' },
    { key: 'scholarship_status', label: 'Scholarship' }
  ],
  'national-team': [
    { key: 'athlete_id', label: 'Athlete ID' },
    { key: 'team_name', label: 'Team' },
    { key: 'category', label: 'Category' },
    { key: 'appearances_count', label: 'Appearances' },
    { key: 'first_call_up_on', label: 'First Cap', type: 'date' }
  ],
  'technical-officials': [
    { key: 'full_name', label: 'Official' },
    { key: 'official_type', label: 'Type' },
    { key: 'level', label: 'Level' },
    { key: 'certification', label: 'Certification' },
    { key: 'status', label: 'Status' }
  ],
  'medical-records': [
    { key: 'athlete_id', label: 'Athlete ID' },
    { key: 'blood_group', label: 'Blood Group' },
    { key: 'current_injury_status', label: 'Injury Status' },
    { key: 'medical_insurance', label: 'Insurance' }
  ],
  'safeguarding-records': [
    { key: 'athlete_id', label: 'Athlete ID' },
    { key: 'consent_forms_url', label: 'Consent URL' },
    { key: 'anti_doping_education_completed', label: 'Anti-Doping Ed.', type: 'boolean' },
    { key: 'guardian_details', label: 'Guardian Details', type: 'json' }
  ],
  'anti-doping': [
    { key: 'athlete_id', label: 'Athlete ID' },
    { key: 'testing_status', label: 'Testing Status' },
    { key: 'last_tested_on', label: 'Last Tested', type: 'date' },
    { key: 'last_test_result', label: 'Result' },
    { key: 'wada_education_completed', label: 'WADA Completed', type: 'boolean' }
  ],
  disbursements: [
    { key: 'federation_id', label: 'Federation ID' },
    { key: 'reference', label: 'Reference' },
    { key: 'amount', label: 'Amount', type: 'money' },
    { key: 'released_on', label: 'Released', type: 'date' },
    { key: 'status', label: 'Status' }
  ],
  accountabilities: [
    { key: 'federation_id', label: 'Federation ID' },
    { key: 'reporting_period_id', label: 'Period ID' },
    { key: 'government_grant', label: 'Govt. Grant', type: 'money' },
    { key: 'sponsorship', label: 'Sponsorship', type: 'money' },
    { key: 'status', label: 'Status' }
  ],
  equipment: [
    { key: 'federation_id', label: 'Federation ID' },
    { key: 'item_name', label: 'Item Name' },
    { key: 'quantity_received', label: 'Qty Received' },
    { key: 'quantity_distributed', label: 'Qty Distributed' },
    { key: 'unit', label: 'Unit' }
  ]
}

const activeRegistry = ref('athletes')
const loading = ref(false)
const items = ref([])
const totalCount = ref(0)
const page = ref(1)
const perPage = ref(20)
const searchQuery = ref('')
const error = ref('')

const sortKey = ref('')
const sortOrder = ref('asc')

const columns = computed(() => columnsConfig[activeRegistry.value] || [])
const totalPages = computed(() => Math.ceil(totalCount.value / perPage.value))

const sortedItems = computed(() => {
  if (!sortKey.value) return items.value
  const key = sortKey.value
  const order = sortOrder.value === 'asc' ? 1 : -1
  return [...items.value].sort((a, b) => {
    const valA = a[key]
    const valB = b[key]
    if (valA === valB) return 0
    if (valA == null) return 1
    if (valB == null) return -1
    return valA < valB ? -order : order
  })
})

onMounted(loadRegistryData)

async function loadRegistryData() {
  loading.value = true
  error.value = ''
  try {
    const res = await listNsmisDomain(activeRegistry.value, {
      page: page.value,
      per_page: perPage.value,
      search: searchQuery.value
    })
    const payload = res?.data?.data || res?.data || {}
    items.value = payload.items || []
    totalCount.value = payload.pagination?.total || items.value.length
  } catch (err) {
    error.value = 'Failed to load registry records. Ensure backend API is active.'
  } finally {
    loading.value = false
  }
}

function onSearchInput() {
  page.value = 1
  loadRegistryData()
}

function changePage(p) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  loadRegistryData()
}

function sortBy(key) {
  if (sortKey.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortOrder.value = 'asc'
  }
}

function formatDate(val) {
  if (!val) return '-'
  return new Date(val).toLocaleDateString('en-UG', {
    year: 'numeric', month: 'short', day: 'numeric'
  })
}

function formatMoney(val) {
  return new Intl.NumberFormat('en-UG').format(Number(val || 0))
}

function truncateJson(val) {
  if (!val) return '{}'
  const str = typeof val === 'object' ? JSON.stringify(val) : String(val)
  return str.length > 30 ? str.slice(0, 27) + '...' : str
}

function exportToCSV() {
  if (!items.value.length) return
  
  const headers = columns.value.map(col => `"${col.label.replace(/"/g, '""')}"`).join(',')
  const rows = items.value.map(row => {
    return columns.value.map(col => {
      let cell = row[col.key] ?? ''
      if (typeof cell === 'object') {
        cell = JSON.stringify(cell)
      }
      return `"${String(cell).replace(/"/g, '""')}"`
    }).join(',')
  })
  
  const csvContent = 'data:text/csv;charset=utf-8,' + [headers, ...rows].join('\n')
  const encodedUri = encodeURI(csvContent)
  const link = document.createElement('a')
  link.setAttribute('href', encodedUri)
  link.setAttribute('download', `${activeRegistry.value}_report_${Date.now()}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
</script>

<style scoped>
.sortable-header:hover {
  background: rgba(0, 0, 0, 0.02);
}
.bg-success-light {
  background: rgba(71, 195, 99, 0.15) !important;
}
.bg-danger-light {
  background: rgba(252, 84, 75, 0.15) !important;
}
</style>
