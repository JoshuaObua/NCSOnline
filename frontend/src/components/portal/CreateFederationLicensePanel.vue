<template>
  <div class="create-license-panel">
    <div class="panel-header-row mb-4">
      <div>
        <h2 class="panel-title">
          <i class="icofont-certificate text-primary me-2"></i>
          National Sports Federation Licensing & Accreditation
        </h2>
        <p class="panel-subtitle">
          Issue new statutory recognition licenses or extend active accreditation under the National Sports Act 2023.
        </p>
      </div>
      <div>
        <button
          type="button"
          class="btn btn-outline-secondary"
          @click="$emit('navigate', 'manage-federation-licenses')"
        >
          <i class="icofont-list me-1"></i> View All Licenses
        </button>
      </div>
    </div>

    <!-- Alert Messages -->
    <div v-if="localError" class="alert alert-danger d-flex align-items-center justify-content-between mb-4">
      <div><i class="icofont-warning me-2"></i> {{ localError }}</div>
      <button type="button" class="btn-close" @click="localError = ''"></button>
    </div>

    <div v-if="localSuccess" class="alert alert-success d-flex align-items-center justify-content-between mb-4">
      <div><i class="icofont-check-circled me-2"></i> {{ localSuccess }}</div>
      <button type="button" class="btn-close" @click="localSuccess = ''"></button>
    </div>

    <!-- License Form Card -->
    <div class="card shadow-sm border-0">
      <div class="card-body p-4">
        <!-- 1. Searchable Dynamic Federation Selector -->
        <div class="mb-4">
          <label class="form-label fw-bold d-flex align-items-center justify-content-between">
            <span>Select National Sports Federation / Association <span class="text-danger">*</span></span>
            <span v-if="loadingFeds" class="small text-muted">
              <i class="icofont-spinner-alt-3 icofont-spin me-1"></i> Loading federations...
            </span>
            <span v-else class="small text-muted">{{ federations.length }} Federations loaded</span>
          </label>

          <!-- Searchable Combobox -->
          <div class="searchable-fed-picker position-relative">
            <div class="input-group">
              <span class="input-group-text bg-white border-end-0"><i class="icofont-search-1 text-muted"></i></span>
              <input
                v-model="fedSearchQuery"
                type="text"
                class="form-control border-start-0 ps-0"
                placeholder="Search by Federation Name, Acronym, Reg Number (e.g. NCS-STATUTORY), or License Number..."
                @focus="isDropdownOpen = true"
                @input="isDropdownOpen = true"
              />
              <button
                v-if="form.federation_id"
                type="button"
                class="btn btn-outline-secondary"
                title="Clear selected federation"
                @click="clearSelectedFederation"
              >
                <i class="icofont-close-line"></i> Clear
              </button>
            </div>

            <!-- Custom Dropdown List -->
            <div
              v-if="isDropdownOpen && filteredFederations.length"
              class="fed-dropdown-menu shadow-lg rounded border mt-1"
            >
              <div
                v-for="fed in filteredFederations"
                :key="fed.id"
                class="fed-dropdown-item p-3 border-bottom cursor-pointer"
                :class="{ 'bg-primary-subtle': fed.id === form.federation_id }"
                @click="selectFederation(fed)"
              >
                <div class="d-flex align-items-center justify-content-between">
                  <div class="fw-bold text-dark">
                    {{ fed.name }}
                    <span v-if="fed.acronym" class="badge bg-secondary ms-1">{{ fed.acronym }}</span>
                  </div>
                  <span v-if="fed.active_license_number" class="badge bg-success-subtle text-success font-monospace">
                    Lic: {{ fed.active_license_number }}
                  </span>
                </div>
                <div class="d-flex align-items-center gap-3 small text-muted mt-1">
                  <span>Reg No: <strong>{{ fed.ncs_registration_number || fed.registration_number || 'N/A' }}</strong></span>
                  <span v-if="fed.president">· Pres: {{ fed.president }}</span>
                  <span v-if="fed.secretary">· Sec: {{ fed.secretary }}</span>
                </div>
              </div>
            </div>
            <div
              v-else-if="isDropdownOpen && fedSearchQuery && !filteredFederations.length"
              class="fed-dropdown-menu shadow-lg rounded border mt-1 p-3 text-center text-muted small"
            >
              No federations found matching "{{ fedSearchQuery }}".
            </div>
          </div>

          <!-- Selected Federation Card Banner -->
          <div v-if="selectedFederation" class="selected-fed-summary p-3 mt-3 bg-light rounded border">
            <div class="d-flex align-items-center justify-content-between flex-wrap gap-2">
              <div>
                <span class="badge bg-primary text-uppercase mb-1" style="font-size: 10px;">Selected Federation</span>
                <h5 class="mb-0 fw-bold text-dark">{{ selectedFederation.name }} ({{ selectedFederation.acronym || 'N/A' }})</h5>
                <span class="small text-muted">
                  Reg No: <strong>{{ selectedFederation.ncs_registration_number || selectedFederation.registration_number || 'NCS-STATUTORY' }}</strong>
                  <span v-if="selectedFederation.president"> &middot; President: <strong>{{ selectedFederation.president }}</strong></span>
                  <span v-if="selectedFederation.secretary"> &middot; Gen Sec: <strong>{{ selectedFederation.secretary }}</strong></span>
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
                <i class="icofont-info-circle me-1"></i> Active Statutory License Found
              </strong>
              <span class="small text-dark">
                This federation currently holds active license <strong>{{ activeLicense.license_number }}</strong> (Valid until {{ formatDate(activeLicense.expiry_date) }}).
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
        <form v-if="formMode === 'extend' && activeLicense" @submit.prevent="handleExtendSubmit">
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
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExtendYears(4)">+4 Years (Olympic)</button>
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
                placeholder="State the regulatory basis, audit clearance, AGM compliance verification, or board resolution granting this extension..."
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
              @click="$emit('navigate', 'manage-federation-licenses')"
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
        <form v-else @submit.prevent="handleSubmit">
          <div class="row g-4">
            <!-- 2. License Number -->
            <div class="col-md-6">
              <label class="form-label fw-bold">
                Official License Number <span class="text-danger">*</span>
              </label>
              <div class="input-group">
                <input
                  v-model="form.license_number"
                  type="text"
                  class="form-control font-monospace fw-bold"
                  placeholder="e.g. NCS/FED-LIC/2026/001"
                  required
                />
                <button
                  type="button"
                  class="btn btn-outline-secondary"
                  title="Generate new unique license number"
                  @click="generateLicenseNumber"
                >
                  <i class="icofont-magic me-1"></i> Auto-Generate
                </button>
              </div>
              <small class="text-muted">Unique tracking identifier printed on official recognition instruments.</small>
            </div>

            <!-- 3. License Type -->
            <div class="col-md-6">
              <label class="form-label fw-bold">
                License Recognition Type <span class="text-danger">*</span>
              </label>
              <select v-model="form.license_type" class="form-select form-control" required>
                <option value="FULL_RECOGNITION">Full Statutory Recognition</option>
                <option value="PROVISIONAL">Provisional Recognition License</option>
                <option value="ANNUAL_COMPLIANCE">Annual Compliance & Accreditation</option>
                <option value="SPECIAL_CLEARANCE">Special Regulatory Clearance</option>
              </select>
            </div>

            <!-- 4. Recognition Category / Tier -->
            <div class="col-md-6">
              <label class="form-label fw-bold">Federation Category / Tier</label>
              <select v-model="form.category" class="form-select form-control">
                <option value="Tier 1 National Sports Federation">Tier 1 National Sports Federation (Olympic / Priority)</option>
                <option value="Tier 2 National Sports Federation">Tier 2 National Sports Federation</option>
                <option value="Tier 3 Associate Sports Body">Tier 3 Associate Sports Body</option>
                <option value="Provisional Member">Provisional Member Federation</option>
              </select>
            </div>

            <!-- 5. Initial Status -->
            <div class="col-md-6">
              <label class="form-label fw-bold">Initial Status</label>
              <select v-model="form.status" class="form-select form-control">
                <option value="ACTIVE">ACTIVE (Issued & Valid)</option>
                <option value="EXTENDED">EXTENDED (Validity Prolonged)</option>
                <option value="SUSPENDED">SUSPENDED (Under Review)</option>
              </select>
            </div>

            <!-- 6. Issue Date & Expiry Date -->
            <div class="col-md-6">
              <label class="form-label fw-bold">
                Issue Date <span class="text-danger">*</span>
              </label>
              <input
                v-model="form.issue_date"
                type="date"
                class="form-control"
                required
                @change="onIssueDateChange"
              />
            </div>

            <div class="col-md-6">
              <label class="form-label fw-bold">
                Expiry Date (Valid Until) <span class="text-danger">*</span>
              </label>
              <input
                v-model="form.expiry_date"
                type="date"
                class="form-control"
                required
              />
              <div class="preset-buttons mt-2 d-flex gap-2">
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExpiryMonths(6)">6 Months</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExpiryYears(1)">1 Year (Standard)</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExpiryYears(2)">2 Years</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExpiryYears(4)">4 Years (Olympic)</button>
              </div>
            </div>

            <!-- 7. Document Attachment URL -->
            <div class="col-md-12">
              <label class="form-label fw-bold">Accompanying Document / Certificate File URL</label>
              <input
                v-model="form.document_url"
                type="text"
                class="form-control"
                placeholder="Optional: Link to signed gazette or recognition instrument"
              />
            </div>

            <!-- 8. Licensing Conditions & Remarks -->
            <div class="col-md-12">
              <label class="form-label fw-bold">Statutory Conditions & Regulatory Remarks</label>
              <textarea
                v-model="form.conditions"
                class="form-control"
                rows="4"
                placeholder="Specify any special conditions, audit submission deadlines, reporting covenants, or governance requirements..."
              ></textarea>
            </div>
          </div>

          <!-- Form Submit Actions -->
          <div class="form-actions mt-5 pt-3 border-top d-flex align-items-center justify-content-between">
            <button
              type="button"
              class="btn btn-light"
              @click="$emit('navigate', 'manage-federation-licenses')"
            >
              Cancel
            </button>

            <button
              type="submit"
              class="btn btn-primary btn-lg px-4 fw-bold"
              :disabled="saving"
            >
              <i class="icofont-check me-1" :class="{ 'icofont-spin icofont-spinner': saving }"></i>
              {{ saving ? 'Issuing License...' : 'Issue & Sign License' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import Swal from 'sweetalert2'
import {
  createFederationLicense,
  extendFederationLicense,
  getFederationActiveLicense,
  listFederationLicenses
} from '@/api/federationLicenses.js'
import { listNsmisDomain } from '@/api/nsmis.js'

const emit = defineEmits(['navigate', 'saved', 'message', 'error'])

const saving = ref(false)
const loadingFeds = ref(false)
const localError = ref('')
const localSuccess = ref('')
const federations = ref([])
const activeLicensesMap = ref({})
const isDropdownOpen = ref(false)
const fedSearchQuery = ref('')
const activeLicense = ref(null)
const formMode = ref('new') // 'new' or 'extend'

const extendForm = reactive({
  new_expiry_date: '',
  reason: 'Annual statutory governance and compliance audit cleared.',
  notes: ''
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

const form = reactive({
  federation_id: '',
  license_number: '',
  license_type: 'FULL_RECOGNITION',
  category: 'Tier 1 National Sports Federation',
  issue_date: new Date().toISOString().slice(0, 10),
  expiry_date: '',
  status: 'ACTIVE',
  conditions: 'Compliant with National Sports Act 2023. Subject to annual governance, financial accountability, and anti-doping audits.',
  document_url: '',
})

const selectedFederation = computed(() => {
  return federations.value.find(f => f.id === form.federation_id) || null
})

const filteredFederations = computed(() => {
  const q = fedSearchQuery.value.trim().toLowerCase()
  if (!q) return federations.value

  return federations.value.filter(fed => {
    const nameMatch = (fed.name || '').toLowerCase().includes(q)
    const acronymMatch = (fed.acronym || '').toLowerCase().includes(q)
    const regMatch = (fed.ncs_registration_number || fed.registration_number || '').toLowerCase().includes(q)
    const licMatch = (fed.active_license_number || '').toLowerCase().includes(q)
    return nameMatch || acronymMatch || regMatch || licMatch
  })
})

onMounted(async () => {
  setExpiryYears(1)
  generateLicenseNumber()
  await loadFederations()
})

async function loadFederations() {
  loadingFeds.value = true
  try {
    const [fedsRes, licsRes] = await Promise.allSettled([
      listNsmisDomain('federations', { per_page: 200 }),
      listFederationLicenses({ per_page: 200 })
    ])

    const rawFeds = fedsRes.status === 'fulfilled' ? (fedsRes.value.data || fedsRes.value.items || fedsRes.value || []) : []
    const rawLics = licsRes.status === 'fulfilled' ? (licsRes.value?.data || []) : []

    const licsByFed = {}
    if (Array.isArray(rawLics)) {
      rawLics.forEach(l => {
        if (l.federation_id && !licsByFed[l.federation_id]) {
          licsByFed[l.federation_id] = l
        }
      })
    }
    activeLicensesMap.value = licsByFed

    const list = Array.isArray(rawFeds) ? rawFeds : []
    federations.value = list.map(f => ({
      ...f,
      active_license_number: licsByFed[f.id]?.license_number || ''
    }))
  } catch (e) {
    console.error('Could not load federations:', e)
  } finally {
    loadingFeds.value = false
  }
}

async function selectFederation(fed) {
  form.federation_id = fed.id
  fedSearchQuery.value = `${fed.name} (${fed.acronym || 'NCS'})`
  isDropdownOpen.value = false
  generateLicenseNumber()

  try {
    const lic = await getFederationActiveLicense(fed.id)
    activeLicense.value = lic || activeLicensesMap.value[fed.id] || null
    if (activeLicense.value) {
      formMode.value = 'extend'
      const curExp = activeLicense.value.expiry_date ? new Date(activeLicense.value.expiry_date) : new Date()
      curExp.setFullYear(curExp.getFullYear() + 1)
      extendForm.new_expiry_date = curExp.toISOString().slice(0, 10)
    } else {
      formMode.value = 'new'
    }
  } catch {
    activeLicense.value = activeLicensesMap.value[fed.id] || null
    formMode.value = activeLicense.value ? 'extend' : 'new'
  }
}

function clearSelectedFederation() {
  form.federation_id = ''
  fedSearchQuery.value = ''
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
  let code = 'FED'
  if (selectedFederation.value?.acronym) {
    code = selectedFederation.value.acronym.toUpperCase().replace(/[^A-Z]/g, '')
  }
  form.license_number = `NCS/FED-LIC/${year}/${code}-${rand}`
}

function onIssueDateChange() {
  if (form.issue_date) {
    setExpiryYears(1)
  }
}

function setExpiryMonths(m) {
  const base = form.issue_date ? new Date(form.issue_date) : new Date()
  base.setMonth(base.getMonth() + m)
  form.expiry_date = base.toISOString().slice(0, 10)
}

function setExpiryYears(y) {
  const base = form.issue_date ? new Date(form.issue_date) : new Date()
  base.setFullYear(base.getFullYear() + y)
  form.expiry_date = base.toISOString().slice(0, 10)
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

async function handleExtendSubmit() {
  if (!activeLicense.value) return
  localError.value = ''
  localSuccess.value = ''
  saving.value = true

  try {
    const res = await extendFederationLicense(activeLicense.value.id, {
      new_expiry_date: extendForm.new_expiry_date,
      reason: extendForm.reason,
      notes: extendForm.notes
    })

    localSuccess.value = `License ${activeLicense.value.license_number} extended until ${formatDate(extendForm.new_expiry_date)}!`
    emit('message', localSuccess.value)

    await Swal.fire({
      title: 'License Extended!',
      text: `License ${activeLicense.value.license_number} for ${selectedFederation.value?.name || 'the federation'} has been successfully extended until ${formatDate(extendForm.new_expiry_date)}.`,
      icon: 'success',
      confirmButtonText: 'View All Licenses'
    })

    emit('saved', res)
    emit('navigate', 'manage-federation-licenses')
  } catch (err) {
    localError.value = err.response?.data?.error?.message || err.message || 'Could not extend federation license.'
    emit('error', localError.value)
  } finally {
    saving.value = false
  }
}

async function handleSubmit() {
  localError.value = ''
  localSuccess.value = ''
  saving.value = true

  try {
    const res = await createFederationLicense(form)
    localSuccess.value = `Federation License ${form.license_number} issued successfully!`
    emit('message', localSuccess.value)

    await Swal.fire({
      title: 'License Issued!',
      text: `License ${form.license_number} has been registered and is now active for ${selectedFederation.value?.name || 'the federation'}.`,
      icon: 'success',
      confirmButtonText: 'View All Licenses',
    })

    emit('saved', res)
    emit('navigate', 'manage-federation-licenses')
  } catch (err) {
    localError.value = err.response?.data?.error?.message || err.message || 'Could not issue federation license.'
    emit('error', localError.value)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.create-license-panel {
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

.preset-buttons button {
  font-size: 11px;
  padding: 2px 8px;
}

.fed-dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #ffffff;
  z-index: 1050;
  max-height: 280px;
  overflow-y: auto;
}

.fed-dropdown-item {
  transition: background-color 0.15s ease;
}

.fed-dropdown-item:hover {
  background-color: #f1f5f9;
}

.cursor-pointer {
  cursor: pointer;
}
</style>
