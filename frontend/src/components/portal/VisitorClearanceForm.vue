<template>
  <div class="official-form-container">
    <div class="official-form-paper">
      
      <!-- Official Header Section -->
      <div class="official-header-section">
        <div class="official-logo-wrapper">
          <img src="/main-logo.png" alt="NCS Logo" class="official-logo-img" />
        </div>
        <div class="official-title-block">
          <h1 class="official-org-title">NATIONAL COUNCIL OF SPORTS</h1>
          <h2 class="official-form-title">VISITOR CLEARANCE & ACCESS PASS</h2>
          <div class="official-form-sub">Security Clearance & Reception Log</div>
        </div>
        
        <!-- Authenticity QR -->
        <div class="document-qr-code" v-if="documentSerial">
          <QrcodeVue :value="documentSerial" :size="64" level="M" />
          <div class="serial-text">SERIAL: {{ documentSerial }}</div>
        </div>
      </div>

      <div class="official-divider"></div>

      <div v-if="successMsg" class="cms-success p-3 rounded mb-4 no-print">
        <i class="icofont-check-circled mr-2"></i> {{ successMsg }}
      </div>
      <div v-if="errorMsg" class="cms-error p-3 rounded mb-4 no-print">
        <i class="icofont-warning mr-2"></i> {{ errorMsg }}
      </div>

      <!-- Visitor Details Block -->
      <form @submit.prevent="submitClearance" class="official-form-body">
        
        <fieldset class="form-fieldset">
          <legend>Visitor Details</legend>
          <div class="form-grid-3">
            <div class="form-group">
              <label for="visitor-first-name">First Name *</label>
              <input id="visitor-first-name" v-model.trim="form.firstName" type="text" required />
            </div>
            <div class="form-group">
              <label for="visitor-middle-name">Middle Name (s)</label>
              <input id="visitor-middle-name" v-model.trim="form.middleName" type="text" />
            </div>
            <div class="form-group">
              <label for="visitor-last-name">Last Name *</label>
              <input id="visitor-last-name" v-model.trim="form.lastName" type="text" required />
            </div>
          </div>

          <div class="form-grid-3 mt-3">
            <div class="form-group">
              <label for="visitor-org">Organization / Representing</label>
              <input id="visitor-org" v-model.trim="form.organization" type="text" placeholder="e.g. FUFA, independent" />
            </div>
            <div class="form-group">
              <label for="visitor-phone">Phone Number *</label>
              <input id="visitor-phone" v-model.trim="form.phone" type="tel" required />
            </div>
            <div class="form-group">
              <label for="visitor-email">Email Address *</label>
              <input id="visitor-email" v-model.trim="form.email" type="email" required />
            </div>
          </div>

          <div class="form-grid-2 mt-3">
            <div class="form-group">
              <label for="visitor-id-type">Identification ID Type *</label>
              <select id="visitor-id-type" v-model="form.idType" required>
                <option value="" disabled>Select ID Type...</option>
                <option value="National ID">National ID Card</option>
                <option value="Passport">Passport</option>
                <option value="Driving License">Driving License</option>
                <option value="Staff Badge">Staff Badge</option>
              </select>
            </div>
            <div class="form-group">
              <label for="visitor-id-no">ID / Card Number *</label>
              <input id="visitor-id-no" v-model.trim="form.idNumber" type="text" required />
            </div>
          </div>
        </fieldset>

        <!-- Visit details Block -->
        <fieldset class="form-fieldset">
          <legend>Visit Details & Host Designation</legend>
          <div class="form-grid-3">
            <div class="form-group">
              <label for="host-name">Host Officer Name *</label>
              <input id="host-name" v-model.trim="form.hostOfficer" type="text" placeholder="e.g. Dr. Bernard Patrick Ogwel" required />
            </div>
            <div class="form-group">
              <label for="host-dept">Host Department *</label>
              <select id="host-dept" v-model="form.department" required>
                <option value="" disabled>Select department</option>
                <option v-for="dept in departments" :key="dept.id" :value="dept.name">{{ dept.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label for="visit-purpose">Purpose of Visit *</label>
              <input id="visit-purpose" v-model.trim="form.purpose" type="text" placeholder="e.g. Licensing Review meeting" required />
            </div>
          </div>

          <div class="form-grid-2 mt-3">
            <div class="form-group">
              <label for="visit-arrival">Arrival Date & Time *</label>
              <input id="visit-arrival" v-model="form.arrivalTime" type="datetime-local" required />
            </div>
            <div class="form-group">
              <label for="visit-departure">Expected Departure Time *</label>
              <input id="visit-departure" v-model="form.departureTime" type="datetime-local" required />
            </div>
          </div>

          <div class="form-grid-2 mt-3">
            <div class="form-group">
              <label for="visit-vehicle">Vehicle Registration (if driving)</label>
              <input id="visit-vehicle" v-model.trim="form.vehicleReg" type="text" placeholder="e.g. UAX 123B" />
            </div>
            <div class="form-group">
              <label for="visit-items">Equipment / Hand luggage items carried</label>
              <input id="visit-items" v-model.trim="form.itemsCarried" type="text" placeholder="e.g. Dell Laptop Serial #829A" />
            </div>
          </div>
        </fieldset>

        <!-- Official Signoffs (Checked when printed) -->
        <fieldset class="form-fieldset">
          <legend>Reception & Gate Security Sign-Offs (Official Use Only)</legend>
          <div class="form-grid-3 text-xs text-slate-500 pt-2">
            <div>
              <p class="mb-1"><strong>1. Host Authorization:</strong></p>
              <div class="border-b border-dashed border-slate-300 h-8"></div>
              <p class="mt-1">Signature & Date</p>
            </div>
            <div>
              <p class="mb-1"><strong>2. Gate Check In (Security Officer):</strong></p>
              <div class="border-b border-dashed border-slate-300 h-8"></div>
              <p class="mt-1">Time & Badge ID</p>
            </div>
            <div>
              <p class="mb-1"><strong>3. Gate Check Out (Security Officer):</strong></p>
              <div class="border-b border-dashed border-slate-300 h-8"></div>
              <p class="mt-1">Time & Signature</p>
            </div>
          </div>
        </fieldset>

        <!-- Action Row -->
        <div class="form-actions-paper no-print flex justify-end gap-2 mt-6">
          <button type="button" class="btn btn-outline-secondary" @click="goBack">
            Cancel
          </button>
          <button type="submit" class="btn btn-primary px-4" :disabled="submitting">
            <span v-if="submitting">Submitting clearance...</span>
            <span v-else>Submit Clearance Request</span>
          </button>
          <button v-if="editId" type="button" class="btn btn-info" @click="printDoc">
            <i class="icofont-print"></i> Print Gate Pass
          </button>
        </div>

      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { listInstitutionalDepartments } from '@/api/cms.js'
import { storedPortalUser } from '@/utils/portalAuth.js'
import { createAuditLog } from '@/utils/auditLogger.js'
import QrcodeVue from 'qrcode.vue'

const router = useRouter()
const route = useRoute()

const defaultDepartments = [
  { id: 'engineering', name: 'Engineering' },
  { id: 'helpdesk', name: 'Helpdesk' },
  { id: 'hr', name: 'Human Resources' },
  { id: 'admin', name: 'Administration' },
  { id: 'accounts', name: 'Accounts & Finance' },
  { id: 'sports_dev', name: 'Sports Development' },
  { id: 'legal', name: 'Legal' },
  { id: 'secretariat', name: 'Secretariat' }
]

const departments = ref([])
const submitting = ref(false)
const successMsg = ref('')
const errorMsg = ref('')
const editId = ref(null)

const documentSerial = computed(() => {
  if (editId.value) return `NCS-VC-2026-${String(editId.value).padStart(4, '0')}`
  return ''
})

const form = ref({
  firstName: '',
  middleName: '',
  lastName: '',
  organization: '',
  phone: '',
  email: '',
  idType: '',
  idNumber: '',
  hostOfficer: '',
  department: '',
  purpose: '',
  arrivalTime: '',
  departureTime: '',
  vehicleReg: '',
  itemsCarried: ''
})

onMounted(async () => {
  // Fetch departments
  try {
    const resDepts = await listInstitutionalDepartments()
    const parsed = resDepts?.data || resDepts?.entries || resDepts || []
    const mapped = (Array.isArray(parsed) ? parsed : []).map(d => ({ id: d.id, name: d.name }))
    const combined = [...mapped, ...defaultDepartments]
    departments.value = combined.filter((v, i, a) => a.findIndex(t => t.name === v.name) === i)
  } catch (e) {
    departments.value = defaultDepartments
  }

  // Check edit query parameters
  const queryId = route.query.id
  if (queryId) {
    editId.value = Number(queryId)
    loadClearanceForEdit(editId.value)
  }

  if (route.query.print === 'true') {
    setTimeout(() => {
      window.print()
    }, 500)
  }
})

function loadClearanceForEdit(id) {
  const stored = localStorage.getItem('ncsms_visitor_clearances')
  if (stored) {
    try {
      const parsed = JSON.parse(stored)
      const found = parsed.find(item => item.id === id)
      if (found) {
        form.value.firstName = found.firstName || ''
        form.value.middleName = found.middleName || ''
        form.value.lastName = found.lastName || ''
        form.value.organization = found.organization || ''
        form.value.phone = found.visitorPhone || ''
        form.value.email = found.visitorEmail || ''
        form.value.idType = found.idType || ''
        form.value.idNumber = found.idNumber || ''
        form.value.hostOfficer = found.hostOfficer || ''
        form.value.department = found.department || ''
        form.value.purpose = found.purpose || ''
        form.value.arrivalTime = found.arrivalTime || ''
        form.value.departureTime = found.departureTime || ''
        form.value.vehicleReg = found.vehicleReg || ''
        form.value.itemsCarried = found.itemsCarried || ''
      }
    } catch (e) {}
  }
}

function submitClearance() {
  submitting.value = true
  successMsg.value = ''
  errorMsg.value = ''

  setTimeout(() => {
    let list = []
    const stored = localStorage.getItem('ncsms_visitor_clearances')
    if (stored) {
      try { list = JSON.parse(stored) } catch (e) {}
    }

    const fullName = [form.value.firstName, form.value.middleName, form.value.lastName].filter(Boolean).join(' ')
    const record = {
      id: editId.value || Date.now(),
      reference: editId.value ? list.find(l => l.id === editId.value)?.reference : 'VC-2026-' + Math.floor(100 + Math.random() * 900),
      visitorName: fullName,
      firstName: form.value.firstName,
      middleName: form.value.middleName,
      lastName: form.value.lastName,
      visitorPhone: form.value.phone,
      visitorEmail: form.value.email,
      organization: form.value.organization,
      idType: form.value.idType,
      idNumber: form.value.idNumber,
      hostOfficer: form.value.hostOfficer,
      department: form.value.department,
      purpose: form.value.purpose,
      arrivalTime: form.value.arrivalTime,
      departureTime: form.value.departureTime,
      vehicleReg: form.value.vehicleReg,
      itemsCarried: form.value.itemsCarried,
      status: editId.value ? (list.find(l => l.id === editId.value)?.status || 'Pending') : 'Pending'
    }

    if (editId.value) {
      const idx = list.findIndex(l => l.id === editId.value)
      if (idx !== -1) list[idx] = record
      successMsg.value = 'Visitor clearance request has been updated successfully.'
    } else {
      list.unshift(record)
      successMsg.value = 'Visitor clearance permit application has been submitted successfully.'
    }

    localStorage.setItem('ncsms_visitor_clearances', JSON.stringify(list))

    // Create Audit Log
    createAuditLog({
      action: editId.value ? 'visitor_clearance:update' : 'visitor_clearance:create',
      status: 'SUCCESS',
      severity: 'INFO',
      actor: storedPortalUser(),
      resourceId: record.reference,
      resourceType: 'VisitorClearance',
      payloadBefore: editId.value ? (list.find(l => l.id === editId.value) || {}) : {},
      payloadAfter: record
    })

    // Create Notification Log
    const notification = {
      id: 'notif_' + Date.now(),
      title: editId.value ? 'Visitor Permit Updated' : 'Visitor Permit Submitted',
      message: `Access pass clearance request for ${fullName} to visit ${form.value.hostOfficer} has been recorded.`,
      status: 'unread',
      icon_key: 'icofont-badge',
      created_at: new Date().toISOString()
    }
    let notifs = []
    const storedNotifs = localStorage.getItem('ncsms_local_notifications')
    if (storedNotifs) {
      try { notifs = JSON.parse(storedNotifs) } catch (e) {}
    }
    notifs.unshift(notification)
    localStorage.setItem('ncsms_local_notifications', JSON.stringify(notifs))

    submitting.value = false
    setTimeout(() => {
      router.push('/portal/helpdesk?tab=visitor-clearance')
    }, 1000)
  }, 1000)
}

function printDoc() {
  window.print()
}

function goBack() {
  router.push('/portal/helpdesk?tab=visitor-clearance')
}
</script>

<style scoped>
.official-form-container {
  padding: 1.5rem;
  background: #f1f5f9;
  min-height: 100%;
}

.official-form-paper {
  background: #ffffff;
  padding: 2.5rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1);
  border: 1px solid #e2e8f0;
  max-width: 900px;
  margin: 0 auto;
}

.official-header-section {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 20px;
  position: relative;
}

.document-qr-code {
  position: absolute;
  top: 0;
  right: 0;
  text-align: right;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.serial-text {
  font-size: 8pt;
  font-family: monospace;
  font-weight: bold;
  color: #64748b;
  margin-top: 4px;
}

.official-logo-img {
  height: 72px;
  object-fit: contain;
}

.official-title-block {
  flex-grow: 1;
}

.official-org-title {
  font-size: 1.25rem;
  font-weight: 800;
  color: #1e3a8a;
  letter-spacing: 0.05em;
  margin: 0;
}

.official-form-title {
  font-size: 1.125rem;
  font-weight: 700;
  color: #0f172a;
  margin-top: 0.25rem;
  margin-bottom: 0.125rem;
}

.official-form-sub {
  font-size: 0.75rem;
  color: #64748b;
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.025em;
}

.official-divider {
  height: 4px;
  background: #1e3a8a;
  margin-bottom: 2rem;
}

.form-fieldset {
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 1.25rem;
  margin-bottom: 1.5rem;
  background: #ffffff;
}

.form-fieldset legend {
  font-size: 0.75rem;
  font-weight: 700;
  color: #1e3a8a;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 0 0.5rem;
  width: auto;
}

.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-grid-3 {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.form-group label {
  font-size: 0.7rem;
  font-weight: 700;
  color: #475569;
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.form-group input, 
.form-group select {
  padding: 0.5rem 0.75rem;
  border: 1px solid #cbd5e1;
  border-radius: 4px;
  font-size: 0.85rem;
  outline: none;
  transition: border-color 0.15s ease;
  background: #f8fafc;
}

.form-group input:focus, 
.form-group select:focus {
  border-color: #2563eb;
  background: #ffffff;
}

</style>
