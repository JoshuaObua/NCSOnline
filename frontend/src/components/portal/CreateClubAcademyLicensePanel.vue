<template>
  <section class="cms-panel create-club-license-panel">
    <div class="cms-panel-head d-flex align-items-center justify-content-between mb-4">
      <div>
        <h2 class="mb-1">
          <i :class="mode === 'academy' ? 'icofont-graduate-alt' : 'icofont-certificate-alt-1'" class="text-primary me-2"></i>
          {{ panelTitle }}
        </h2>
        <p class="text-muted small mb-0">{{ panelSubtitle }}</p>
      </div>
      <div>
        <button
          type="button"
          class="btn btn-outline-secondary btn-sm"
          @click="$emit('navigate', targetManageSection)"
        >
          <i class="icofont-list me-1"></i> View All {{ mode === 'academy' ? 'Academy' : (mode === 'club' ? 'Club' : '') }} Licenses
        </button>
      </div>
    </div>

    <div v-if="successMsg" class="alert alert-success alert-dismissible fade show mb-4" role="alert">
      <i class="icofont-check-circled me-2"></i> {{ successMsg }}
      <button type="button" class="btn-close" @click="successMsg = ''"></button>
    </div>

    <div v-if="errorMsg" class="alert alert-danger alert-dismissible fade show mb-4" role="alert">
      <i class="icofont-warning me-2"></i> {{ errorMsg }}
      <button type="button" class="btn-close" @click="errorMsg = ''"></button>
    </div>

    <div class="card border-0 shadow-sm">
      <div class="card-body p-4">
        <!-- 1. Searchable Dynamic Club / Academy Selector -->
        <div class="mb-4">
          <label class="form-label fw-bold d-flex align-items-center justify-content-between">
            <span>Select {{ mode === 'academy' ? 'Youth Sports Academy' : (mode === 'club' ? 'Sports Club' : 'Sports Club / Academy') }} <span class="text-danger">*</span></span>
            <span v-if="loadingClubs" class="small text-muted">
              <i class="icofont-spinner-alt-3 icofont-spin me-1"></i> Loading records...
            </span>
            <span v-else class="small text-muted">{{ filteredClubsList.length }} records loaded</span>
          </label>

          <!-- Searchable Combobox -->
          <div class="searchable-club-picker position-relative">
            <div class="input-group">
              <span class="input-group-text bg-white border-end-0"><i class="icofont-search-1 text-muted"></i></span>
              <input
                v-model="clubSearchQuery"
                type="text"
                class="form-control border-start-0 ps-0"
                :placeholder="`Search ${mode === 'academy' ? 'Academy' : 'Club'} by Name, Registration No (e.g. NCS-ACA-001), Acronym, or License...`"
                @focus="isDropdownOpen = true"
                @input="isDropdownOpen = true"
              />
              <button
                v-if="form.club_id"
                type="button"
                class="btn btn-outline-secondary"
                title="Clear selection"
                @click="clearSelectedClub"
              >
                <i class="icofont-close-line"></i> Clear
              </button>
            </div>

            <!-- Dropdown Menu -->
            <div
              v-if="isDropdownOpen && searchedClubs.length"
              class="club-dropdown-menu shadow-lg rounded border mt-1"
            >
              <div
                v-for="c in searchedClubs"
                :key="c.id"
                class="club-dropdown-item p-3 border-bottom cursor-pointer"
                :class="{ 'bg-primary-subtle': c.id === form.club_id }"
                @click="selectClub(c)"
              >
                <div class="d-flex align-items-center justify-content-between">
                  <div class="fw-bold text-dark">
                    {{ c.name }}
                    <span v-if="c.acronym" class="badge bg-secondary ms-1">{{ c.acronym }}</span>
                  </div>
                  <span v-if="c.active_license_number" class="badge bg-success-subtle text-success font-monospace">
                    Lic: {{ c.active_license_number }}
                  </span>
                </div>
                <div class="d-flex align-items-center gap-3 small text-muted mt-1">
                  <span>Reg No: <strong>{{ c.club_number || 'N/A' }}</strong></span>
                  <span v-if="c.category">· Category: {{ c.category }}</span>
                  <span v-if="c.district">· District: {{ c.district }}</span>
                </div>
              </div>
            </div>
            <div
              v-else-if="isDropdownOpen && clubSearchQuery && !searchedClubs.length"
              class="club-dropdown-menu shadow-lg rounded border mt-1 p-3 text-center text-muted small"
            >
              No matching records found.
            </div>
          </div>

          <!-- Selected Summary Card -->
          <div v-if="selectedClub" class="selected-club-summary p-3 mt-3 bg-light rounded border">
            <div class="d-flex align-items-center justify-content-between flex-wrap gap-2">
              <div>
                <span class="badge bg-primary text-uppercase mb-1" style="font-size: 10px;">Selected Entity</span>
                <h5 class="mb-0 fw-bold text-dark">{{ selectedClub.name }} ({{ selectedClub.acronym || 'N/A' }})</h5>
                <span class="small text-muted">
                  Reg No: <strong>{{ selectedClub.club_number || 'N/A' }}</strong>
                  <span v-if="selectedClub.president"> &middot; President: <strong>{{ selectedClub.president }}</strong></span>
                  <span v-if="selectedClub.secretary"> &middot; Secretary: <strong>{{ selectedClub.secretary }}</strong></span>
                  <span v-if="selectedClub.district"> &middot; District: <strong>{{ selectedClub.district }}</strong></span>
                </span>
              </div>
              <div v-if="activeLicense" class="text-end">
                <span class="badge" :class="activeLicense.status === 'REVOKED' ? 'bg-danger' : 'bg-success'">
                  Active License: {{ activeLicense.license_number }}
                </span>
                <div class="small text-muted">Valid Until: <strong>{{ formatDate(activeLicense.expiry_date) }}</strong></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Mode Toggle if Active License Exists -->
        <div v-if="activeLicense" class="action-mode-selector mb-4 p-3 bg-warning-subtle rounded border border-warning">
          <div class="d-flex align-items-center justify-content-between flex-wrap gap-2">
            <div>
              <strong class="text-warning-emphasis d-block">
                <i class="icofont-info-circle me-1"></i> Active Operating License Found
              </strong>
              <span class="small text-dark">
                This {{ isAcademy(selectedClub) ? 'academy' : 'club' }} currently holds active license <strong>{{ activeLicense.license_number }}</strong> (Valid until {{ formatDate(activeLicense.expiry_date) }}).
              </span>
            </div>
            <div class="btn-group">
              <button
                type="button"
                class="btn btn-sm"
                :class="formMode === 'extend' ? 'btn-info text-white fw-bold' : 'btn-outline-secondary'"
                @click="setFormMode('extend')"
              >
                <i class="icofont-calendar me-1"></i> Extend Active License
              </button>
              <button
                type="button"
                class="btn btn-sm"
                :class="formMode === 'new' ? 'btn-primary fw-bold' : 'btn-outline-secondary'"
                @click="setFormMode('new')"
              >
                <i class="icofont-plus-circle me-1"></i> Issue New License
              </button>
            </div>
          </div>
        </div>

        <!-- ================= EXTENSION FORM ================= -->
        <form v-if="formMode === 'extend' && activeLicense" @submit.prevent="submitExtendLicense">
          <div class="row g-4">
            <div class="col-md-6">
              <label class="form-label fw-bold">Current License Number</label>
              <input :value="activeLicense.license_number" type="text" class="form-control font-monospace fw-bold" readonly disabled />
            </div>

            <div class="col-md-6">
              <label class="form-label fw-bold">Current Expiry Date</label>
              <input :value="formatDate(activeLicense.expiry_date)" type="text" class="form-control" readonly disabled />
            </div>

            <div class="col-md-6">
              <label class="form-label fw-bold">
                New Extended Expiry Date <span class="text-danger">*</span>
              </label>
              <input
                v-model="extendForm.new_expiry_date"
                type="date"
                class="form-control"
                required
              />
              <div class="preset-buttons mt-2 d-flex gap-2">
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExtendMonths(6)">+6 Months</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExtendYears(1)">+1 Year</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExtendYears(2)">+2 Years</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExtendYears(3)">+3 Years</button>
              </div>
            </div>

            <div class="col-md-6">
              <label class="form-label fw-bold">Extension Authorized By</label>
              <input :value="currentUserName" type="text" class="form-control" readonly disabled />
              <small class="text-muted">Recorded automatically in the immutable statutory audit log.</small>
            </div>

            <div class="col-12">
              <label class="form-label fw-bold">
                Extension Rationale & Regulatory Justification <span class="text-danger">*</span>
              </label>
              <textarea
                v-model="extendForm.reason"
                class="form-control"
                rows="3"
                required
                placeholder="State the regulatory basis, safeguarding verification, facility inspection, or accreditation clearance..."
              ></textarea>
            </div>

            <div class="col-12">
              <label class="form-label fw-bold">Internal Audit Notes (Optional)</label>
              <input
                v-model="extendForm.notes"
                type="text"
                class="form-control"
                placeholder="Optional inspection reference or file memo"
              />
            </div>
          </div>

          <div class="form-actions mt-5 pt-3 border-top d-flex align-items-center justify-content-between">
            <button
              type="button"
              class="btn btn-light"
              @click="$emit('navigate', targetManageSection)"
            >
              Cancel
            </button>

            <button
              type="submit"
              class="btn btn-info btn-lg px-4 text-white fw-bold"
              :disabled="saving"
            >
              <i class="icofont-check me-1" :class="{ 'icofont-spin icofont-spinner': saving }"></i>
              {{ saving ? 'Extending License...' : 'Confirm & Sign License Extension' }}
            </button>
          </div>
        </form>

        <!-- ================= NEW LICENSE CREATION FORM ================= -->
        <form v-else @submit.prevent="submitLicense">
          <div class="row g-3">
            <!-- Auto-generated License Number -->
            <div class="col-md-6">
              <div class="d-flex align-items-center justify-content-between">
                <label class="form-label fw-bold">Statutory License Number <span class="text-danger">*</span></label>
                <button type="button" class="btn btn-link btn-sm p-0 text-decoration-none" @click="generateLicenseNumber">
                  <i class="icofont-refresh"></i> Auto-Generate
                </button>
              </div>
              <input
                v-model="form.license_number"
                type="text"
                class="form-control font-monospace fw-bold"
                required
                placeholder="e.g. NCS/ACA-LIC/2026/001"
              />
              <small class="text-muted">Unique official tracking identifier for this accreditation certificate.</small>
            </div>

            <!-- License Type -->
            <div class="col-md-6">
              <label class="form-label fw-bold">License Type <span class="text-danger">*</span></label>
              <select v-model="form.license_type" class="form-select form-control" required>
                <option value="STATUTORY_RECOGNITION">Statutory Recognition License</option>
                <option value="OPERATING_PERMIT">Operating & Training Permit</option>
                <option value="DEVELOPMENT_ACCREDITATION">Youth Development Accreditation</option>
                <option value="PROVISIONAL_PERMIT">Provisional Operating Permit</option>
              </select>
            </div>

            <!-- Category Tier -->
            <div class="col-md-6">
              <label class="form-label fw-bold">Category Tier <span class="text-danger">*</span></label>
              <select v-model="form.category" class="form-select form-control" required>
                <option value="SPORTS_ACADEMY">Sports Academy</option>
                <option value="YOUTH_ACADEMY">Youth Academy</option>
                <option value="ELITE_ACADEMY">Elite Academy</option>
                <option value="DEVELOPMENT_CENTRE">Grassroots Development Centre</option>
                <option value="COMMUNITY_CLUB">Community Sports Club</option>
                <option value="SENIOR_CLUB">Senior Sports Club</option>
              </select>
            </div>

            <!-- Issue Date -->
            <div class="col-md-6">
              <label class="form-label fw-bold">Issue Date <span class="text-danger">*</span></label>
              <input v-model="form.issue_date" type="date" class="form-control" required />
            </div>

            <!-- Expiry Date with Presets -->
            <div class="col-md-6">
              <label class="form-label fw-bold">Expiry Date <span class="text-danger">*</span></label>
              <input v-model="form.expiry_date" type="date" class="form-control" required />
              <div class="preset-buttons mt-2 d-flex gap-2">
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setPresetExpiry(1)">1 Year</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setPresetExpiry(2)">2 Years</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setPresetExpiry(3)">3 Years</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setPresetExpiry(5)">5 Years</button>
              </div>
            </div>

            <!-- Document URL -->
            <div class="col-md-6">
              <label class="form-label fw-bold">Accompanying Inspection / Permit Document URL</label>
              <input
                v-model="form.document_url"
                type="text"
                class="form-control"
                placeholder="Optional: Link to signed inspection report or charter"
              />
            </div>

            <!-- Conditions and Terms -->
            <div class="col-12">
              <div class="d-flex align-items-center justify-content-between">
                <label class="form-label fw-bold">Statutory Conditions & Safeguarding Covenants</label>
                <div class="d-flex gap-2">
                  <button type="button" class="btn btn-link btn-xs p-0 text-decoration-none" @click="applyTemplate('standard')">Standard Terms</button>
                  <span class="text-muted">|</span>
                  <button type="button" class="btn btn-link btn-xs p-0 text-decoration-none" @click="applyTemplate('safeguarding')">Safeguarding Focus</button>
                </div>
              </div>
              <textarea
                v-model="form.conditions"
                class="form-control"
                rows="3"
                placeholder="Specify licensing requirements, athlete protection standards, coaching certificate minimums..."
              ></textarea>
            </div>
          </div>

          <div class="form-actions mt-5 pt-3 border-top d-flex align-items-center justify-content-between">
            <button type="button" class="btn btn-light" @click="resetForm">
              Reset
            </button>
            <button type="submit" class="btn btn-primary btn-lg px-4 fw-bold" :disabled="saving">
              <i class="icofont-check me-1" :class="{ 'icofont-spin icofont-spinner': saving }"></i>
              {{ saving ? 'Issuing License...' : 'Issue & Sign License' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import Swal from 'sweetalert2'
import {
  createClubAcademyLicense,
  extendClubAcademyLicense,
  getClubActiveLicense,
  listClubAcademyLicenses
} from '@/api/clubAcademyLicenses.js'
import { listNsmisDomain } from '@/api/nsmis.js'

const props = defineProps({
  mode: {
    type: String,
    default: 'all' // 'club', 'academy', or 'all'
  }
})

const emit = defineEmits(['navigate', 'message', 'error'])

const saving = ref(false)
const loadingClubs = ref(false)
const successMsg = ref('')
const errorMsg = ref('')
const clubsList = ref([])
const activeLicensesMap = ref({})
const isDropdownOpen = ref(false)
const clubSearchQuery = ref('')
const activeLicense = ref(null)
const formMode = ref('new')

const panelTitle = computed(() => {
  if (props.mode === 'club') return 'Issue National Sports Club License'
  if (props.mode === 'academy') return 'Issue Youth Sports Academy License'
  return 'Issue Club & Academy License'
})

const panelSubtitle = computed(() => {
  if (props.mode === 'club') return 'Statutory Recognition & Operating Permit for National Sports Clubs under National Sports Act 2023'
  if (props.mode === 'academy') return 'Statutory Operating Accreditation & Youth Safeguarding License under National Sports Act 2023'
  return 'Statutory Recognition & Operating Accreditation under National Sports Act 2023'
})

const targetManageSection = computed(() => {
  if (props.mode === 'club') return 'manage-club-licenses'
  if (props.mode === 'academy') return 'manage-academy-licenses'
  return 'manage-club-academy-licenses'
})

const currentUser = computed(() => {
  try {
    return JSON.parse(localStorage.getItem('ncsms_user') || '{}')
  } catch {
    return {}
  }
})

const currentUserName = computed(() => {
  const u = currentUser.value
  return (u.first_name || u.name ? `${u.first_name || ''} ${u.last_name || ''}`.trim() : 'Dr. Bernard Ogwel (General Secretary)')
})

const extendForm = reactive({
  new_expiry_date: '',
  reason: 'Annual statutory operating accreditation and safeguarding compliance cleared.',
  notes: ''
})

const form = reactive({
  club_id: '',
  license_number: '',
  license_type: 'STATUTORY_RECOGNITION',
  category: props.mode === 'club' ? 'SENIOR_CLUB' : 'SPORTS_ACADEMY',
  issue_date: new Date().toISOString().substring(0, 10),
  expiry_date: '',
  conditions: 'Granted subject to compliance with the National Sports Act 2023, youth safeguarding protocols, and athlete registry standards.',
  document_url: '',
})

function isAcademy(club) {
  if (!club) return false
  const cat = String(club.category || '').toUpperCase()
  const name = String(club.name || '').toUpperCase()
  return cat.includes('ACADEMY') || cat.includes('DEVELOPMENT') || name.includes('ACADEMY')
}

const filteredClubsList = computed(() => {
  if (props.mode === 'academy') {
    return clubsList.value.filter(c => isAcademy(c))
  }
  if (props.mode === 'club') {
    return clubsList.value.filter(c => !isAcademy(c) || c.category?.toUpperCase().includes('CLUB'))
  }
  return clubsList.value
})

const searchedClubs = computed(() => {
  const q = clubSearchQuery.value.trim().toLowerCase()
  if (!q) return filteredClubsList.value

  return filteredClubsList.value.filter(c => {
    const nameMatch = (c.name || '').toLowerCase().includes(q)
    const acronymMatch = (c.acronym || '').toLowerCase().includes(q)
    const numMatch = (c.club_number || '').toLowerCase().includes(q)
    const licMatch = (c.active_license_number || '').toLowerCase().includes(q)
    return nameMatch || acronymMatch || numMatch || licMatch
  })
})

const selectedClub = computed(() => {
  return clubsList.value.find(c => c.id === form.club_id) || null
})

onMounted(async () => {
  setPresetExpiry(1)
  generateLicenseNumber()
  await loadClubs()
})

async function loadClubs() {
  loadingClubs.value = true
  try {
    const [clubsRes, licsRes] = await Promise.allSettled([
      listNsmisDomain('clubs', { per_page: 200 }),
      listClubAcademyLicenses({ per_page: 200 })
    ])

    const rawClubs = clubsRes.status === 'fulfilled' ? (clubsRes.value?.data || clubsRes.value?.items || clubsRes.value || []) : []
    const rawLics = licsRes.status === 'fulfilled' ? (licsRes.value?.data || []) : []

    const licsByClub = {}
    if (Array.isArray(rawLics)) {
      rawLics.forEach(l => {
        if (l.club_id && !licsByClub[l.club_id]) {
          licsByClub[l.club_id] = l
        }
      })
    }
    activeLicensesMap.value = licsByClub

    const list = Array.isArray(rawClubs) ? rawClubs : []
    clubsList.value = list.map(c => ({
      ...c,
      active_license_number: licsByClub[c.id]?.license_number || ''
    }))
  } catch (err) {
    console.error('Failed to load clubs:', err)
  } finally {
    loadingClubs.value = false
  }
}

async function selectClub(c) {
  form.club_id = c.id
  clubSearchQuery.value = `${c.name} (${c.club_number || c.acronym || 'ID'})`
  isDropdownOpen.value = false
  generateLicenseNumber()

  if (c.category) {
    form.category = c.category
  }

  try {
    const lic = await getClubActiveLicense(c.id)
    activeLicense.value = lic || activeLicensesMap.value[c.id] || null
    if (activeLicense.value) {
      formMode.value = 'extend'
      const curExp = activeLicense.value.expiry_date ? new Date(activeLicense.value.expiry_date) : new Date()
      curExp.setFullYear(curExp.getFullYear() + 1)
      extendForm.new_expiry_date = curExp.toISOString().slice(0, 10)
    } else {
      formMode.value = 'new'
    }
  } catch {
    activeLicense.value = activeLicensesMap.value[c.id] || null
    formMode.value = activeLicense.value ? 'extend' : 'new'
  }
}

function clearSelectedClub() {
  form.club_id = ''
  clubSearchQuery.value = ''
  activeLicense.value = null
  formMode.value = 'new'
  generateLicenseNumber()
}

function setFormMode(mode) {
  formMode.value = mode
  if (mode === 'extend' && activeLicense.value) {
    const curExp = activeLicense.value.expiry_date ? new Date(activeLicense.value.expiry_date) : new Date()
    curExp.setFullYear(curExp.getFullYear() + 1)
    extendForm.new_expiry_date = curExp.toISOString().slice(0, 10)
  }
}

function generateLicenseNumber() {
  const year = new Date().getFullYear()
  const rand = Math.floor(100 + Math.random() * 900)
  let prefix = props.mode === 'club' ? 'CLUB' : 'ACA'
  if (selectedClub.value?.acronym) {
    prefix = selectedClub.value.acronym.toUpperCase().replace(/[^A-Z]/g, '')
  }
  form.license_number = `NCS/${prefix}-LIC/${year}/${rand}`
}

function setPresetExpiry(years) {
  const d = form.issue_date ? new Date(form.issue_date) : new Date()
  d.setFullYear(d.getFullYear() + years)
  form.expiry_date = d.toISOString().substring(0, 10)
}

function setExtendMonths(m) {
  const base = activeLicense.value?.expiry_date ? new Date(activeLicense.value.expiry_date) : new Date()
  base.setMonth(base.getMonth() + m)
  extendForm.new_expiry_date = base.toISOString().slice(0, 10)
}

function setExtendYears(y) {
  const base = activeLicense.value?.expiry_date ? new Date(activeLicense.value.expiry_date) : new Date()
  base.setFullYear(base.getFullYear() + y)
  extendForm.new_expiry_date = base.toISOString().slice(0, 10)
}

function formatDate(val) {
  if (!val) return '-'
  try {
    return new Date(val).toLocaleDateString('en-UG', { year: 'numeric', month: 'short', day: 'numeric' })
  } catch {
    return String(val)
  }
}

function applyTemplate(type) {
  if (type === 'standard') {
    form.conditions = 'Granted subject to compliance with the National Sports Act 2023, anti-doping protocols, verified coaching staff, and annual athlete registry reporting.'
  } else if (type === 'safeguarding') {
    form.conditions = 'Mandatory compliance with National Child Safeguarding in Sports protocols, qualified medical clearance, and certified youth coaching staff standards.'
  }
}

function resetForm() {
  form.club_id = ''
  clubSearchQuery.value = ''
  activeLicense.value = null
  formMode.value = 'new'
  form.license_type = 'STATUTORY_RECOGNITION'
  form.category = props.mode === 'club' ? 'SENIOR_CLUB' : 'SPORTS_ACADEMY'
  form.issue_date = new Date().toISOString().substring(0, 10)
  setPresetExpiry(1)
  generateLicenseNumber()
  form.conditions = 'Granted subject to compliance with the National Sports Act 2023, youth safeguarding protocols, and athlete registry standards.'
  successMsg.value = ''
  errorMsg.value = ''
}

async function submitExtendLicense() {
  if (!activeLicense.value) return
  errorMsg.value = ''
  successMsg.value = ''
  saving.value = true

  try {
    await extendClubAcademyLicense(activeLicense.value.id, {
      new_expiry_date: extendForm.new_expiry_date,
      reason: extendForm.reason,
      notes: extendForm.notes
    })

    successMsg.value = `License ${activeLicense.value.license_number} extended until ${formatDate(extendForm.new_expiry_date)}!`
    emit('message', successMsg.value)

    await Swal.fire({
      title: 'License Extended!',
      text: `License ${activeLicense.value.license_number} for ${selectedClub.value?.name || 'the club/academy'} has been successfully extended until ${formatDate(extendForm.new_expiry_date)}.`,
      icon: 'success',
      confirmButtonText: 'View All Licenses'
    })

    emit('navigate', targetManageSection.value)
  } catch (err) {
    errorMsg.value = err.response?.data?.error?.message || err.message || 'Could not extend license.'
    emit('error', errorMsg.value)
  } finally {
    saving.value = false
  }
}

async function submitLicense() {
  if (!form.club_id) {
    errorMsg.value = 'Please select a club or academy.'
    return
  }
  if (!form.license_number) {
    errorMsg.value = 'License number is required.'
    return
  }
  if (!form.expiry_date) {
    errorMsg.value = 'Please specify an expiry date.'
    return
  }

  saving.value = true
  successMsg.value = ''
  errorMsg.value = ''

  try {
    const res = await createClubAcademyLicense({
      club_id: form.club_id,
      license_number: form.license_number.trim(),
      license_type: form.license_type,
      category: form.category,
      issue_date: form.issue_date,
      expiry_date: form.expiry_date,
      conditions: form.conditions.trim(),
      document_url: form.document_url,
    })

    successMsg.value = `Successfully issued license ${res.license_number} for ${(res.club_name || selectedClub.value?.name || 'Club/Academy')}!`
    emit('message', successMsg.value)

    await Swal.fire({
      title: 'License Issued!',
      text: `License ${form.license_number} has been registered and is now active for ${selectedClub.value?.name || 'the club/academy'}.`,
      icon: 'success',
      confirmButtonText: 'View All Licenses'
    })

    emit('navigate', targetManageSection.value)
  } catch (err) {
    errorMsg.value = err.response?.data?.error?.message || err.message || 'Failed to issue club/academy license'
    emit('error', errorMsg.value)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.create-club-license-panel {
  padding: 24px;
}
.btn-xs {
  font-size: 11px;
  padding: 2px 8px;
}
.club-dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #ffffff;
  z-index: 1050;
  max-height: 280px;
  overflow-y: auto;
}
.club-dropdown-item {
  transition: background-color 0.15s ease;
}
.club-dropdown-item:hover {
  background-color: #f1f5f9;
}
.cursor-pointer {
  cursor: pointer;
}
</style>
