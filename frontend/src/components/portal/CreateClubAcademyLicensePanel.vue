<template>
  <section class="cms-panel create-club-license-panel">
    <div class="cms-panel-head d-flex align-items-center justify-content-between">
      <div>
        <h2 class="mb-1"><i class="icofont-certificate-alt-1 text-primary me-2"></i> Issue Club & Academy License</h2>
        <p class="text-muted small mb-0">Statutory Recognition & Operating Accreditation under National Sports Act 2023</p>
      </div>
      <div>
        <button type="button" class="btn btn-outline-secondary btn-sm" @click="$emit('navigate', 'manage-club-academy-licenses')">
          <i class="icofont-list me-1"></i> View All Licenses
        </button>
      </div>
    </div>

    <div v-if="successMsg" class="alert alert-success alert-dismissible fade show mt-3" role="alert">
      <i class="icofont-check-circled me-2"></i> {{ successMsg }}
      <button type="button" class="btn-close" @click="successMsg = ''"></button>
    </div>

    <div v-if="errorMsg" class="alert alert-danger alert-dismissible fade show mt-3" role="alert">
      <i class="icofont-warning me-2"></i> {{ errorMsg }}
      <button type="button" class="btn-close" @click="errorMsg = ''"></button>
    </div>

    <form class="otika-form-card mt-3" @submit.prevent="submitLicense">
      <div class="card border-0 shadow-sm">
        <div class="card-body p-4">
          <div class="row g-3">
            <!-- Club / Academy Selector -->
            <div class="col-md-6">
              <label class="form-label fw-bold">Select Sports Club / Academy <span class="text-danger">*</span></label>
              <select v-model="form.club_id" class="form-select form-control" required @change="onClubSelected">
                <option value="" disabled>-- Select Registered Club / Academy --</option>
                <option v-for="c in clubsList" :key="c.id" :value="c.id">
                  {{ c.name }} ({{ c.club_number || c.acronym || 'ID: ' + c.id.substring(0, 8) }})
                </option>
              </select>
              <small v-if="selectedClub" class="text-muted d-block mt-1">
                <strong>Number:</strong> {{ selectedClub.club_number || 'N/A' }} | <strong>Acronym:</strong> {{ selectedClub.acronym || 'N/A' }} | <strong>District:</strong> {{ selectedClub.district || 'National' }}
              </small>
            </div>

            <!-- Auto-generated License Number -->
            <div class="col-md-6">
              <div class="d-flex align-items-center justify-content-between">
                <label class="form-label fw-bold">Statutory License Number <span class="text-danger">*</span></label>
                <button type="button" class="btn btn-link btn-sm p-0 text-decoration-none" @click="generateLicenseNumber">
                  <i class="icofont-refresh"></i> Generate New
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
                <option value="SPORTS_ACADEMY">Sports Academy (Multi-discipline)</option>
                <option value="YOUTH_ACADEMY">Youth Football / Sports Academy</option>
                <option value="DEVELOPMENT_CENTRE">Grassroots Development Centre</option>
                <option value="COMMUNITY_CLUB">Community Sports Club</option>
                <option value="ELITE_ACADEMY">Elite High-Performance Academy</option>
                <option value="SENIOR_CLUB">Senior League Sports Club</option>
              </select>
            </div>

            <!-- Issue Date -->
            <div class="col-md-6">
              <label class="form-label fw-bold">Date of Issuance <span class="text-danger">*</span></label>
              <input v-model="form.issue_date" type="date" class="form-control" required />
            </div>

            <!-- Expiry Date & Presets -->
            <div class="col-md-6">
              <div class="d-flex align-items-center justify-content-between">
                <label class="form-label fw-bold">Expiration / Renewal Date <span class="text-danger">*</span></label>
                <div class="btn-group btn-group-sm">
                  <button type="button" class="btn btn-outline-primary btn-xs py-0 px-2" @click="setPresetExpiry(1)">1 Yr</button>
                  <button type="button" class="btn btn-outline-primary btn-xs py-0 px-2" @click="setPresetExpiry(2)">2 Yrs</button>
                  <button type="button" class="btn btn-outline-primary btn-xs py-0 px-2" @click="setPresetExpiry(3)">3 Yrs</button>
                  <button type="button" class="btn btn-outline-primary btn-xs py-0 px-2" @click="setPresetExpiry(5)">5 Yrs</button>
                </div>
              </div>
              <input v-model="form.expiry_date" type="date" class="form-control" required />
            </div>

            <!-- Statutory Conditions -->
            <div class="col-12">
              <label class="form-label fw-bold">Statutory & Operational Conditions</label>
              <textarea
                v-model="form.conditions"
                class="form-control"
                rows="3"
                placeholder="Enter conditions for this license, e.g. Child safeguarding compliance, verified coaching staff, and annual athlete registry reporting."
              ></textarea>
              <div class="d-flex gap-2 mt-1">
                <button type="button" class="btn btn-link btn-sm p-0 text-muted" @click="applyTemplate('standard')">
                  + Insert Standard Academy Conditions
                </button>
                <span class="text-muted">|</span>
                <button type="button" class="btn btn-link btn-sm p-0 text-muted" @click="applyTemplate('safeguarding')">
                  + Insert Child Safeguarding & Safety Standard
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="card-footer bg-light d-flex justify-content-between align-items-center py-3">
          <button type="button" class="btn btn-secondary" @click="resetForm">
            <i class="icofont-eraser me-1"></i> Reset
          </button>
          <button type="submit" class="btn btn-primary px-4 fw-bold" :disabled="saving">
            <i v-if="saving" class="icofont-spinner-alt-3 icofont-spin me-1"></i>
            <i v-else class="icofont-check-circled me-1"></i>
            {{ saving ? 'Issuing License...' : 'Issue & Activate License' }}
          </button>
        </div>
      </div>
    </form>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { listNsmisDomain } from '@/api/nsmis.js'
import { createClubAcademyLicense } from '@/api/clubAcademyLicenses.js'

const emit = defineEmits(['navigate', 'message', 'error'])

const clubsList = ref([])
const saving = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const form = reactive({
  club_id: '',
  license_number: '',
  license_type: 'STATUTORY_RECOGNITION',
  category: 'SPORTS_ACADEMY',
  issue_date: new Date().toISOString().substring(0, 10),
  expiry_date: '',
  conditions: 'Granted subject to compliance with the National Sports Act 2023, youth safeguarding protocols, and athlete registry standards.',
  document_url: '',
})

const selectedClub = computed(() => {
  return clubsList.value.find(c => c.id === form.club_id) || null
})

function extractItems(res) {
  if (!res) return []
  if (Array.isArray(res)) return res
  if (Array.isArray(res.items)) return res.items
  if (Array.isArray(res.data?.items)) return res.data.items
  if (Array.isArray(res.data)) return res.data
  return []
}

async function loadClubs() {
  try {
    const res = await listNsmisDomain('clubs', { per_page: 200 })
    clubsList.value = extractItems(res)
  } catch (err) {
    console.error('Failed to load clubs:', err)
  }
}

function generateLicenseNumber() {
  const year = new Date().getFullYear()
  const randNum = Math.floor(100 + Math.random() * 900)
  form.license_number = `NCS/ACA-LIC/${year}/${randNum}`
}

function setPresetExpiry(years) {
  const baseDate = form.issue_date ? new Date(form.issue_date) : new Date()
  const exp = new Date(baseDate)
  exp.setFullYear(exp.getFullYear() + years)
  form.expiry_date = exp.toISOString().substring(0, 10)
}

function onClubSelected() {
  if (!form.license_number) {
    generateLicenseNumber()
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
  form.license_type = 'STATUTORY_RECOGNITION'
  form.category = 'SPORTS_ACADEMY'
  form.issue_date = new Date().toISOString().substring(0, 10)
  setPresetExpiry(1)
  generateLicenseNumber()
  form.conditions = 'Granted subject to compliance with the National Sports Act 2023, youth safeguarding protocols, and athlete registry standards.'
  successMsg.value = ''
  errorMsg.value = ''
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

    setTimeout(() => {
      emit('navigate', 'manage-club-academy-licenses')
    }, 1200)
  } catch (err) {
    errorMsg.value = err.response?.data?.error?.message || err.message || 'Failed to issue club/academy license'
    emit('error', errorMsg.value)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadClubs()
  generateLicenseNumber()
  setPresetExpiry(1)
})
</script>

<style scoped>
.create-club-license-panel {
  padding: 24px;
}
.btn-xs {
  font-size: 11px;
  padding: 2px 8px;
}
</style>
