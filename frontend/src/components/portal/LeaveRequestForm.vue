<template>
  <div class="leave-request-container">
    <!-- Paper Form Document Container -->
    <div class="official-form-paper">
      <!-- Form Header -->
      <div class="form-header-badge text-center position-relative">
        <img alt="NCS Logo" src="/main-logo.png" class="official-logo" />
        <h1 class="ncs-title">NATIONAL COUNCIL OF SPORTS</h1>
        <h2 class="form-title">LEAVE APPLICATION FORM</h2>
        <div class="header-divider"></div>
        
        <!-- Authenticity QR -->
        <div class="document-qr-code" v-if="documentSerial">
          <QrcodeVue :value="documentSerial" :size="64" level="M" />
          <div class="serial-text">SERIAL: {{ documentSerial }}</div>
        </div>
      </div>

      <!-- Main Leave Form -->
      <form @submit.prevent="submitLeave">
        
        <!-- Success and Error banners -->
        <div v-if="successMsg" class="cms-message p-3 rounded mb-4">
          <i class="icofont-check-circled mr-2"></i> {{ successMsg }}
        </div>
        <div v-if="errorMsg" class="cms-error p-3 rounded mb-4">
          <i class="icofont-warning mr-2"></i> {{ errorMsg }}
        </div>

        <!-- Routing Information Block -->
        <fieldset class="form-fieldset">
          <legend>Routing & Approvals</legend>
          <div class="form-grid-2">
            <div class="form-group">
              <label for="route-to">To (General Secretary / Supervisors)</label>
              <input id="route-to" v-model.trim="form.to" type="text" placeholder="" required />
            </div>
            <div class="form-group">
              <label for="route-through">Through (Head of Department)</label>
              <input id="route-through" v-model.trim="form.through" type="text" placeholder="" required />
            </div>
          </div>
        </fieldset>

        <!-- Applicant Details Block -->
        <fieldset class="form-fieldset">
          <legend>Applicant details</legend>
          <div class="form-grid-3">
            <div class="form-group">
              <label for="applicant-first-name">First Name *</label>
              <input id="applicant-first-name" v-model.trim="form.firstName" type="text" required />
            </div>
            <div class="form-group">
              <label for="applicant-middle-name">Middle Name (s)</label>
              <input id="applicant-middle-name" v-model.trim="form.middleName" type="text" placeholder="" />
            </div>
            <div class="form-group">
              <label for="applicant-last-name">Last Name *</label>
              <input id="applicant-last-name" v-model.trim="form.lastName" type="text" required />
            </div>
          </div>
          <div class="form-grid-2 mt-3">
            <div class="form-group">
              <label for="applicant-position">Position / Job Title *</label>
              <select id="applicant-position" v-model="form.position" required>
                <option value="" disabled>Select Position...</option>
                <option v-for="role in rolesList" :key="role.id" :value="role.name">{{ role.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label for="applicant-dept">Department *</label>
              <select id="applicant-dept" v-model="form.department" required>
                <option value="" disabled>Select department</option>
                <option v-for="dept in departments" :key="dept.id" :value="dept.name">{{ dept.name }}</option>
              </select>
            </div>
          </div>
        </fieldset>

        <!-- Leave Details & Days Block -->
        <fieldset class="form-fieldset">
          <legend>Leave details</legend>
          <div class="form-grid-3">
            <div class="form-group">
              <label for="leave-type">Type of Leave</label>
              <select id="leave-type" v-model="form.leaveType" required @change="validateDates">
                <option value="" disabled>Select leave type</option>
                <option value="Annual">Annual Leave</option>
                <option value="Sick">Sick Leave</option>
                <option value="Maternity / Paternity">Maternity / Paternity Leave</option>
                <option value="Compassionate">Compassionate Leave</option>
                <option value="Study">Study Leave</option>
                <option value="Leave Without Pay (LWOP)">Leave Without Pay (LWOP)</option>
              </select>
            </div>
            <div class="form-group">
              <label for="leave-begins">Leave Begins (First date)</label>
              <input id="leave-begins" v-model="form.startDate" type="date" required :min="todayDate" @change="validateDates" />
            </div>
            <div class="form-group">
              <label for="leave-ends">Leave Ends (Last date inclusive)</label>
              <input id="leave-ends" v-model="form.endDate" type="date" required :min="form.startDate || todayDate" @change="validateDates" />
            </div>
          </div>

          <div class="form-grid-2 mt-3">
            <div class="form-group highlight-days">
              <label>Total Number of Days Applied For</label>
              <div class="days-display">{{ totalDaysApplied }} Days</div>
            </div>
            <div class="form-group">
              <label for="rate-entitlement">Rate of Leave Entitlement (Days a Year)</label>
              <input id="rate-entitlement" v-model.number="form.rateEntitlement" type="number" min="0" required />
            </div>
          </div>
        </fieldset>

        <!-- Address While On Leave Block -->
        <fieldset class="form-fieldset">
          <legend>Address & Contact Details While on Leave</legend>
          <div class="form-grid-3">
            <div class="form-group">
              <label for="contact-phone">Telephone No (s)</label>
              <input id="contact-phone" v-model.trim="form.contactPhone" type="tel" placeholder="e.g. +256 772 000000" required />
            </div>
            <div class="form-group">
              <label for="contact-email">Email Address (s)</label>
              <input id="contact-email" v-model.trim="form.contactEmail" type="email" placeholder="e.g. user@ncs.go.ug" required />
            </div>
            <div class="form-group">
              <label for="contact-postal">Postal Address (Optional)</label>
              <input id="contact-postal" v-model.trim="form.postalAddress" type="text" placeholder="e.g. P.O Box 7010 Kampala" />
            </div>
          </div>
          <div class="form-group mt-3">
            <label for="contact-physical">Physical Address (Location during leave)</label>
            <textarea id="contact-physical" v-model.trim="form.physicalAddress" rows="2" placeholder="Describe the physical location/village/town where you will be staying..." required></textarea>
          </div>
        </fieldset>

        <!-- Detailed Leave Computation Block -->
        <fieldset class="form-fieldset bg-light-panel">
          <legend>Detailed Computation of Leave</legend>
          <div class="form-grid-4">
            <div class="form-group">
              <label for="comp-due">a) Leave Due For Current Year (Days)</label>
              <input id="comp-due" v-model.number="form.computation.due" type="number" min="0" required @change="updateBalances" />
            </div>
            <div class="form-group">
              <label for="comp-arrears">b) Leave In Arrears (Days)</label>
              <input id="comp-arrears" v-model.number="form.computation.arrears" type="number" min="0" required @change="updateBalances" />
            </div>
            <div class="form-group">
              <label>Total Leave Days Available</label>
              <div class="computed-value">{{ totalDaysAvailable }} Days</div>
            </div>
            <div class="form-group">
              <label for="comp-taken">c) Less: Leave Days Taken (Days)</label>
              <input id="comp-taken" v-model.number="form.computation.taken" type="number" min="0" required @change="updateBalances" />
            </div>
          </div>
          <div class="balance-display-banner mt-3">
            <strong>d) Computed Leave Days Balance:</strong>
            <span :class="{'text-danger': computedBalance < 0}">{{ computedBalance }} Days</span>
          </div>
        </fieldset>

        <!-- Submit actions -->
        <div class="form-actions-paper mt-4 d-flex justify-content-between align-items-center no-print">
          <div class="note-rules">
            <p class="text-muted small mb-0">Note: Leave requests should be applied for at least two weeks in advance of the intended start date.</p>
          </div>
          <div class="d-flex">
            <button type="button" class="btn btn-secondary mr-2" @click="resetForm" :disabled="submitting">
              Reset Form
            </button>
            <button type="button" class="btn btn-info mr-2" @click="printDoc" :disabled="submitting">
              <i class="icofont-print mr-1"></i> Print / PDF
            </button>
            <button type="submit" class="btn btn-primary" :disabled="submitting || !!dateError">
              <i v-if="submitting" class="icofont-spinner animate-spin mr-1"></i>
              <i v-else class="icofont-upload-alt mr-1"></i>
              {{ editId ? 'Update Application' : 'Submit Application to HR' }}
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { storedPortalUser, roleNames } from '@/utils/portalAuth.js'
import QrcodeVue from 'qrcode.vue'
import { listInstitutionalDepartments, adminListRoles } from '@/api/cms.js'

const route = useRoute()
const user = ref(storedPortalUser() || {})
const editId = ref(null)

const documentSerial = computed(() => {
  if (editId.value) return `NCS-LR-2026-${String(editId.value).padStart(4, '0')}`
  return ''
})

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

const defaultRoles = [
  { id: 'engineer', name: 'Software Engineer' },
  { id: 'system_admin', name: 'System Administrator' },
  { id: 'helpdesk_agent', name: 'Helpdesk Agent' },
  { id: 'hr_officer', name: 'HR Officer' },
  { id: 'admin_officer', name: 'Administration Officer' },
  { id: 'accountant', name: 'Accountant' },
  { id: 'legal_officer', name: 'Legal Counsel' },
  { id: 'sports_officer', name: 'Sports Officer' },
  { id: 'general_secretary', name: 'General Secretary' }
]

const departments = ref([])
const rolesList = ref([])

const form = ref({
  to: '',
  through: '',
  firstName: '',
  middleName: '',
  lastName: '',
  position: '',
  department: '',
  leaveType: '',
  startDate: '',
  endDate: '',
  rateEntitlement: 0,
  contactPhone: '',
  contactEmail: '',
  postalAddress: '',
  physicalAddress: '',
  computation: {
    due: 0,
    arrears: 0,
    taken: 0
  }
})

const submitting = ref(false)
const successMsg = ref('')
const errorMsg = ref('')
const dateError = ref(false)
const recentRequests = ref([])

const todayDate = new Date().toISOString().split('T')[0]

// Automatically pre-populate user context
onMounted(async () => {
  form.value.firstName = user.value.first_name || ''
  form.value.middleName = user.value.middle_name || ''
  form.value.lastName = user.value.last_name || ''
  form.value.contactEmail = user.value.email || ''
  
  const userRoles = roleNames(user.value)
  form.value.department = user.value.department || (userRoles.includes('helpdesk') ? 'Helpdesk' : 'Administration')
  form.value.position = user.value.designation || user.value.position || (userRoles.includes('helpdesk') ? 'Helpdesk Agent' : 'Software Engineer')
  
  loadLeaveRequests()

  // Helper to parse responses in multiple shapes safely
  const asList = (res) => {
    if (!res) return []
    if (Array.isArray(res)) return res
    if (res.data && Array.isArray(res.data)) return res.data
    if (res.entries && Array.isArray(res.entries)) return res.entries
    return []
  }

  // Fetch departments
  try {
    const resDepts = await listInstitutionalDepartments()
    const parsedDepts = asList(resDepts).map(d => ({ id: d.id, name: d.name }))
    const combinedDepts = [...parsedDepts, ...defaultDepartments]
    departments.value = combinedDepts.filter((v, i, a) => a.findIndex(t => t.name === v.name) === i)
  } catch (e) {
    departments.value = defaultDepartments
  }

  // Fetch roles
  try {
    const resRoles = await adminListRoles()
    const parsedRoles = asList(resRoles).map(r => ({ id: r.id, name: r.name }))
    const combinedRoles = [...parsedRoles, ...defaultRoles]
    rolesList.value = combinedRoles.filter((v, i, a) => a.findIndex(t => t.name === v.name) === i)
  } catch (e) {
    rolesList.value = defaultRoles
  }

  const queryId = route.query.id
  if (queryId) {
    editId.value = Number(queryId)
    loadRequestForEdit(editId.value)
  }

  if (route.query.print === 'true') {
    setTimeout(() => {
      window.print()
    }, 500)
  }
})

function loadRequestForEdit(id) {
  const stored = localStorage.getItem('ncsms_leave_requests')
  if (stored) {
    try {
      const parsed = JSON.parse(stored)
      const found = parsed.find(req => req.id === id)
      if (found) {
        if (found.status === 'Approved' || found.status === 'Rejected') {
          errorMsg.value = 'This leave request has already been processed and cannot be edited.'
          return
        }
        form.value.to = found.to || form.value.to
        form.value.through = found.through || form.value.through
        
        const nameParts = (found.applicantName || '').split(' ')
        form.value.firstName = found.firstName || nameParts[0] || ''
        form.value.middleName = found.middleName || (nameParts.length > 2 ? nameParts.slice(1, -1).join(' ') : '')
        form.value.lastName = found.lastName || (nameParts.length > 1 ? nameParts[nameParts.length - 1] : '')
        
        form.value.position = found.position || form.value.position
        form.value.department = found.department || form.value.department
        form.value.leaveType = found.leaveType || form.value.leaveType
        form.value.startDate = found.startDate || form.value.startDate
        form.value.endDate = found.endDate || form.value.endDate
        form.value.rateEntitlement = found.rateEntitlement || form.value.rateEntitlement
        form.value.contactPhone = found.contactPhone || form.value.contactPhone
        form.value.contactEmail = found.contactEmail || form.value.contactEmail
        form.value.postalAddress = found.postalAddress || form.value.postalAddress
        form.value.physicalAddress = found.physicalAddress || form.value.physicalAddress
        if (found.computation) {
          form.value.computation.due = found.computation.due
          form.value.computation.arrears = found.computation.arrears
          form.value.computation.taken = found.computation.taken
        }
      }
    } catch (e) {}
  }
}

const totalDaysApplied = computed(() => {
  if (!form.value.startDate || !form.value.endDate) return 0
  const start = new Date(form.value.startDate)
  const end = new Date(form.value.endDate)
  const diffTime = Math.abs(end - start)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24)) + 1 // inclusive of last day
  return diffDays || 0
})

const totalDaysAvailable = computed(() => {
  return Number(form.value.computation.due || 0) + Number(form.value.computation.arrears || 0)
})

const computedBalance = computed(() => {
  return totalDaysAvailable.value - Number(form.value.computation.taken || 0)
})

function loadLeaveRequests() {
  const stored = localStorage.getItem('ncsms_leave_requests')
  if (stored) {
    try {
      const parsed = JSON.parse(stored)
      recentRequests.value = parsed.map(req => {
        return {
          id: req.id || Date.now(),
          applicantName: req.applicantName || 'Staff Member',
          position: req.position || 'NCS Employee',
          department: req.department || 'Administration',
          leaveType: req.leaveType || 'Annual',
          startDate: req.startDate || '',
          endDate: req.endDate || '',
          appliedDays: req.appliedDays || 1,
          status: req.status || 'Pending',
          workflow: req.workflow || {
            deptHead: req.status === 'Approved',
            hrVerified: req.status === 'Approved',
            gsApproved: req.status === 'Approved'
          }
        }
      })
    } catch (e) {
      recentRequests.value = []
    }
  } else {
    // Seed standard leave item matching official format
    recentRequests.value = [
      {
        id: 1,
        applicantName: 'Baker Nsubuga',
        position: 'Support Agent',
        department: 'Helpdesk',
        leaveType: 'Annual',
        startDate: '2026-08-01',
        endDate: '2026-08-10',
        appliedDays: 10,
        status: 'Approved',
        workflow: {
          deptHead: true,
          hrVerified: true,
          gsApproved: true
        }
      }
    ]
    localStorage.setItem('ncsms_leave_requests', JSON.stringify(recentRequests.value))
  }
}

function validateDates() {
  dateError.value = false
  errorMsg.value = ''
  
  if (form.value.startDate && form.value.endDate) {
    if (new Date(form.value.endDate) < new Date(form.value.startDate)) {
      errorMsg.value = 'End date cannot be prior to start date.'
      dateError.value = true
    }
  }
}

function updateBalances() {
  // Logic to prevent negative inputs
  if (form.value.computation.due < 0) form.value.computation.due = 0
  if (form.value.computation.arrears < 0) form.value.computation.arrears = 0
  if (form.value.computation.taken < 0) form.value.computation.taken = 0
}

function submitLeave() {
  validateDates()
  if (dateError.value) return

  submitting.value = true
  successMsg.value = ''
  errorMsg.value = ''

  setTimeout(() => {
    let requestsList = []
    const stored = localStorage.getItem('ncsms_leave_requests')
    if (stored) {
      try {
        requestsList = JSON.parse(stored)
      } catch (e) {}
    }

    const fullName = [form.value.firstName, form.value.middleName, form.value.lastName].filter(Boolean).join(' ')

    if (editId.value) {
      const idx = requestsList.findIndex(req => req.id === editId.value)
      if (idx !== -1) {
        requestsList[idx] = {
          ...requestsList[idx],
          to: form.value.to,
          through: form.value.through,
          applicantName: fullName,
          firstName: form.value.firstName,
          middleName: form.value.middleName,
          lastName: form.value.lastName,
          position: form.value.position,
          department: form.value.department,
          leaveType: form.value.leaveType,
          startDate: form.value.startDate,
          endDate: form.value.endDate,
          appliedDays: totalDaysApplied.value,
          contactPhone: form.value.contactPhone,
          contactEmail: form.value.contactEmail,
          postalAddress: form.value.postalAddress,
          physicalAddress: form.value.physicalAddress,
          computation: { ...form.value.computation }
        }
        successMsg.value = 'Your official leave application has been updated successfully.'
      }
    } else {
      const newReq = {
        id: Date.now(),
        applicantName: fullName,
        firstName: form.value.firstName,
        middleName: form.value.middleName,
        lastName: form.value.lastName,
        position: form.value.position,
        department: form.value.department,
        leaveType: form.value.leaveType,
        startDate: form.value.startDate,
        endDate: form.value.endDate,
        appliedDays: totalDaysApplied.value,
        status: 'Pending',
        workflow: {
          deptHead: false,
          hrVerified: false,
          gsApproved: false
        },
        computation: { ...form.value.computation }
      }
      requestsList.unshift(newReq)
      successMsg.value = 'Your official leave application has been submitted to the HR department.'
    }

    // Save list
    localStorage.setItem('ncsms_leave_requests', JSON.stringify(requestsList))
    
    // Create alert notification record
    const finalAction = editId.value ? 'updated' : 'submitted'
    const notification = {
      id: 'notif_' + Date.now(),
      title: editId.value ? 'Leave Application Updated' : 'Leave Application Submitted',
      message: `Your leave request for ${form.value.leaveType} (${totalDaysApplied.value} Days) has been successfully ${finalAction}.`,
      status: 'unread',
      icon_key: 'icofont-file-document',
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
    loadLeaveRequests()
    if (!editId.value) resetForm()
  }, 1200)
}

function printDoc() {
  window.print()
}

function resetForm() {
  form.value.leaveType = ''
  form.value.startDate = ''
  form.value.endDate = ''
  form.value.physicalAddress = ''
  form.value.computation.taken = 0
  dateError.value = false
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

<style scoped>
.official-form-paper {
  background: #ffffff;
  border: 1px solid #e1e4e8;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  padding: 40px;
  max-width: 100%;
  position: relative;
}

:global(.dark) .official-form-paper {
  background: #1f2937 !important;
  border-color: #334155 !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.25) !important;
}

.official-logo {
  max-width: 120px;
  height: auto;
  margin-bottom: 12px;
}

.ncs-title {
  font-family: 'Poppins', sans-serif;
  font-size: 22px;
  font-weight: 800;
  color: #1e293b;
  margin-bottom: 4px;
}

:global(.dark) .ncs-title {
  color: #f8fafc !important;
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

.form-title {
  font-size: 16px;
  font-weight: 700;
  color: #64748b;
  letter-spacing: 2px;
  margin-bottom: 18px;
}

.header-divider {
  height: 3px;
  background: #6777ef;
  width: 100%;
  margin-bottom: 24px;
  border-radius: 99px;
}

.form-fieldset {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 24px;
}

:global(.dark) .form-fieldset {
  border-color: #334155 !important;
}

.form-fieldset legend {
  float: none;
  width: auto;
  padding: 0 10px;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  color: #6777ef;
  margin-bottom: 0;
}

.bg-light-panel {
  background: #f8fafc;
}

:global(.dark) .bg-light-panel {
  background: #111827 !important;
}

.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-grid-3 {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 16px;
}

.form-grid-4 {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 11px;
  font-weight: 700;
  color: #475569;
  text-transform: uppercase;
}

:global(.dark) .form-group label {
  color: #cbd5e1 !important;
}

input, select, textarea {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  background: #ffffff;
  color: #1e293b;
  font-size: 13px;
  outline: none;
  transition: all 0.15s ease;
}

:global(.dark) input, :global(.dark) select, :global(.dark) textarea {
  background: #111827 !important;
  border-color: #475569 !important;
  color: #f8fafc !important;
}

input:focus, select:focus, textarea:focus {
  border-color: #6777ef;
  box-shadow: 0 0 0 3px rgba(103, 119, 239, 0.15);
}

.highlight-days {
  justify-content: center;
  align-items: flex-start;
}

.days-display {
  font-size: 18px;
  font-weight: 800;
  color: #10b981;
}

.computed-value {
  font-size: 14px;
  font-weight: 700;
  color: #1e293b;
  padding: 10px 0;
}

:global(.dark) .computed-value {
  color: #f8fafc !important;
}

.balance-display-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #eef2ff;
  padding: 12px 18px;
  border-radius: 6px;
  font-size: 14px;
}

:global(.dark) .balance-display-banner {
  background: #1e1b4b !important;
}

.balance-display-banner strong {
  color: #4f46e5;
}

:global(.dark) .balance-display-banner strong {
  color: #818cf8 !important;
}

.balance-display-banner span {
  font-weight: 800;
  font-size: 16px;
  color: #4f46e5;
}

:global(.dark) .balance-display-banner span {
  color: #818cf8 !important;
}

.form-actions-paper {
  border-top: 1px dashed #cbd5e1;
  padding-top: 20px;
}

:global(.dark) .form-actions-paper {
  border-color: #334155 !important;
}

/* Table paper style */
.official-history-paper {
  background: #ffffff;
  border: 1px solid #e1e4e8;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  padding: 30px;
}

:global(.dark) .official-history-paper {
  background: #1f2937 !important;
  border-color: #334155 !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.25) !important;
}

.section-title-paper {
  font-size: 16px;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 12px;
}

:global(.dark) .section-title-paper {
  color: #f8fafc !important;
}

.official-table th {
  background: #f8fafc !important;
  color: #475569 !important;
  font-weight: 700 !important;
  font-size: 11px !important;
  text-transform: uppercase !important;
}

:global(.dark) .official-table th {
  background: #111827 !important;
  color: #cbd5e1 !important;
}

.badge-status {
  padding: 4px 10px;
  border-radius: 99px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
}

.badge-status.approved { background: #d1fae5; color: #065f46; }
.badge-status.pending { background: #fef3c7; color: #92400e; }
.badge-status.rejected { background: #fee2e2; color: #991b1b; }

.workflow-checklist {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 4px;
}

.workflow-checklist li {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #94a3b8;
}

.workflow-checklist li i {
  font-size: 14px;
}

.workflow-checklist li.completed {
  color: #10b981;
}

@media (max-width: 768px) {
  .form-grid-2, .form-grid-3, .form-grid-4 {
    grid-template-columns: 1fr !important;
  }
}

</style>
