<template>
  <div class="create-license-panel">
    <div class="panel-header-row mb-4">
      <div>
        <h2 class="panel-title">
          <i class="icofont-certificate text-primary mr-2"></i>
          Issue National Sports Federation License
        </h2>
        <p class="panel-subtitle">
          Grant official statutory recognition and accreditation under the National Sports Act 2023.
        </p>
      </div>
      <div>
        <button
          type="button"
          class="btn btn-outline-secondary"
          @click="$emit('navigate', 'manage-federation-licenses')"
        >
          <i class="icofont-list"></i> View All Licenses
        </button>
      </div>
    </div>

    <!-- Alert Messages -->
    <div v-if="localError" class="alert alert-danger d-flex align-items-center justify-content-between mb-4">
      <div><i class="icofont-warning mr-2"></i> {{ localError }}</div>
      <button type="button" class="btn-close" @click="localError = ''">&times;</button>
    </div>

    <div v-if="localSuccess" class="alert alert-success d-flex align-items-center justify-content-between mb-4">
      <div><i class="icofont-check-circled mr-2"></i> {{ localSuccess }}</div>
      <button type="button" class="btn-close" @click="localSuccess = ''">&times;</button>
    </div>

    <!-- License Form Card -->
    <div class="card shadow-sm border-0">
      <div class="card-body p-4">
        <form @submit.prevent="handleSubmit">
          <div class="row g-4">
            <!-- 1. Federation Selection -->
            <div class="col-md-12">
              <label class="form-label fw-bold">
                Sports Federation / Association <span class="text-danger">*</span>
              </label>
              <select
                v-model="form.federation_id"
                class="form-select form-select-lg"
                required
                :disabled="loadingFeds"
                @change="onFederationSelected"
              >
                <option value="" disabled>-- Select National Sports Federation --</option>
                <option
                  v-for="fed in federations"
                  :key="fed.id"
                  :value="fed.id"
                >
                  {{ fed.name }} {{ fed.acronym ? '(' + fed.acronym + ')' : '' }} — Reg: {{ fed.ncs_registration_number || fed.registration_number || 'N/A' }}
                </option>
              </select>
              <small v-if="selectedFederation" class="text-muted mt-1 d-block">
                President: <strong>{{ selectedFederation.president || 'N/A' }}</strong> · General Secretary: <strong>{{ selectedFederation.secretary || 'N/A' }}</strong>
              </small>
            </div>

            <!-- 2. License Number -->
            <div class="col-md-6">
              <label class="form-label fw-bold">
                Official License Number <span class="text-danger">*</span>
              </label>
              <div class="input-group">
                <input
                  v-model="form.license_number"
                  type="text"
                  class="form-control"
                  placeholder="e.g. NCS/FED-LIC/2026/001"
                  required
                />
                <button
                  type="button"
                  class="btn btn-outline-secondary"
                  title="Generate new unique license number"
                  @click="generateLicenseNumber"
                >
                  <i class="icofont-magic"></i> Auto-Generate
                </button>
              </div>
              <small class="text-muted">Unique tracking identifier printed on official recognition documents.</small>
            </div>

            <!-- 3. License Type -->
            <div class="col-md-6">
              <label class="form-label fw-bold">
                License Recognition Type <span class="text-danger">*</span>
              </label>
              <select v-model="form.license_type" class="form-select" required>
                <option value="FULL_RECOGNITION">Full Statutory Recognition</option>
                <option value="PROVISIONAL">Provisional Recognition License</option>
                <option value="ANNUAL_COMPLIANCE">Annual Compliance & Accreditation</option>
                <option value="SPECIAL_CLEARANCE">Special Regulatory Clearance</option>
              </select>
            </div>

            <!-- 4. Recognition Category / Tier -->
            <div class="col-md-6">
              <label class="form-label fw-bold">Federation Category / Tier</label>
              <select v-model="form.category" class="form-select">
                <option value="Tier 1 National Sports Federation">Tier 1 National Sports Federation (Olympic / Priority)</option>
                <option value="Tier 2 National Sports Federation">Tier 2 National Sports Federation</option>
                <option value="Tier 3 Associate Sports Body">Tier 3 Associate Sports Body</option>
                <option value="Provisional Member">Provisional Member Federation</option>
              </select>
            </div>

            <!-- 5. Status -->
            <div class="col-md-6">
              <label class="form-label fw-bold">Initial Status</label>
              <select v-model="form.status" class="form-select">
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
              <div class="input-group">
                <input
                  v-model="form.expiry_date"
                  type="date"
                  class="form-control"
                  required
                />
              </div>
              <div class="preset-buttons mt-2 d-flex gap-2">
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExpiryMonths(6)">6 Months</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExpiryYears(1)">1 Year (Standard)</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExpiryYears(2)">2 Years</button>
                <button type="button" class="btn btn-xs btn-outline-secondary" @click="setExpiryYears(4)">4 Years (Olympic Cycle)</button>
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
              class="btn btn-primary btn-lg px-4"
              :disabled="saving"
            >
              <i class="icofont-check mr-1" :class="{ 'icofont-spin icofont-spinner': saving }"></i>
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
import { createFederationLicense } from '@/api/federationLicenses.js'
import { listNsmisDomain } from '@/api/nsmis.js'

const emit = defineEmits(['navigate', 'saved', 'message', 'error'])

const saving = ref(false)
const loadingFeds = ref(false)
const localError = ref('')
const localSuccess = ref('')
const federations = ref([])

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

onMounted(async () => {
  setExpiryYears(1)
  generateLicenseNumber()
  await loadFederations()
})

async function loadFederations() {
  loadingFeds.value = true
  try {
    const res = await listNsmisDomain('federations', { per_page: 200 })
    const list = res.data || res.items || res || []
    federations.value = Array.isArray(list) ? list : []
  } catch (e) {
    console.error('Could not load federations:', e)
  } finally {
    loadingFeds.value = false
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

function onFederationSelected() {
  generateLicenseNumber()
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
</style>
