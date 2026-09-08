<template>
  <LayoutDefault title="Apply for Leave">
    <div class="space-y-6 pb-12 w-full">
      

      <!-- Main Full-Width Application Form Card -->
      <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 w-full mb-0">
        <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-6 flex items-center justify-between">
          <h4 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
            <i class="icofont-paper-plane text-blue-600"></i> Official Staff Leave Application Form
          </h4>
          <span class="text-xs text-slate-400 font-mono">NCS-HR-FORM-LEV</span>
        </div>

        <div class="card-body p-6 sm:p-8">
          <form @submit.prevent="submitLeaveApplication" class="space-y-6 text-xs w-full">
            
            <div v-if="successMsg" class="p-3.5 rounded-xl bg-emerald-50 dark:bg-emerald-950/50 border border-emerald-200 dark:border-emerald-800 text-emerald-700 dark:text-emerald-300 font-medium text-xs">
              <i class="icofont-check-circled mr-1"></i> {{ successMsg }}
            </div>
            <div v-if="errorMsg" class="p-3.5 rounded-xl bg-red-50 dark:bg-red-950/50 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 font-medium text-xs">
              <i class="icofont-warning mr-1"></i> {{ errorMsg }}
            </div>

            <!-- ── SECTION 1: APPLICANT & ESTABLISHMENT PARTICULARS ────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-blue-100 dark:bg-blue-900/50 text-blue-600 flex items-center justify-center text-[10px] font-bold">1</span>
                Applicant & Establishment Particulars
              </h5>
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Staff Full Name *</label>
                  <input
                    v-model="form.staff_name"
                    required
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                    placeholder="e.g. Allan Tumusiime"
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Designation / Post Title *</label>
                  <SearchableSelect
                    v-model="form.staff_role"
                    :options="designationOptions"
                    placeholder="Search or select designation..."
                    search-placeholder="Search official post titles..."
                    required
                    @update:modelValue="onDesignationSelected"
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Department / Section *</label>
                  <SearchableSelect
                    v-model="form.staff_department"
                    :options="departmentOptions"
                    placeholder="Search or select department..."
                    search-placeholder="Filter departments..."
                    required
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Staff File / IPPS No.</label>
                  <input
                    v-model="form.staff_file_no"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400 font-mono"
                    placeholder="e.g. NCS/ICT/2026/017"
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Official Phone Contact *</label>
                  <input
                    v-model="form.contact_phone"
                    required
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                    placeholder="+256 701 234 567"
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Official Email *</label>
                  <input
                    v-model="form.contact_email"
                    required
                    type="email"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                    placeholder="itofficer@ncs.go.ug"
                  />
                </div>
              </div>
            </div>

            <!-- ── SECTION 2: LEAVE TYPE & DURATION ────────────────────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-emerald-100 dark:bg-emerald-900/50 text-emerald-600 flex items-center justify-center text-[10px] font-bold">2</span>
                Leave Category, Dates & Duration Schedule
              </h5>
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <div class="sm:col-span-2 lg:col-span-4">
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Statutory Leave Type *</label>
                  <SearchableSelect
                    v-model="form.leave_type"
                    :options="leaveTypeOptions"
                    placeholder="Select leave category..."
                    search-placeholder="Search leave types..."
                    required
                  />
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Leave Start Date *</label>
                  <input v-model="form.start_date" @change="calculateDays" required type="date" class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white" />
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Leave End Date *</label>
                  <input v-model="form.end_date" @change="calculateDays" required type="date" class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white" />
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Working Days Requested *</label>
                  <input v-model.number="form.days_requested" required type="number" min="1" max="60" class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white font-mono font-bold" />
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Duty Resumption Date *</label>
                  <input v-model="form.return_date" required type="date" class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white" />
                </div>
              </div>
            </div>

            <!-- ── SECTION 3: DUTY RELIEF & HANDOVER ───────────────────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-purple-100 dark:bg-purple-900/50 text-purple-600 flex items-center justify-center text-[10px] font-bold">3</span>
                Duty Relief Officer & Handover Arrangements
              </h5>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Relieving Officer Name *</label>
                  <input v-model="form.relieving_officer_name" required class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400" placeholder="e.g. Peter Ssewankambo" />
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Relieving Officer Post / Designation *</label>
                  <SearchableSelect
                    v-model="form.relieving_officer_role"
                    :options="designationOptions"
                    placeholder="Select relief officer post..."
                    search-placeholder="Search post titles..."
                    required
                  />
                </div>
                <div class="sm:col-span-2">
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Specific Key Duties / Work Handed Over *</label>
                  <textarea v-model="form.duty_handover_details" required rows="2" class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400" placeholder="Specify ongoing assignments, physical keys, database backup monitors, or active files delegated to relief officer..."></textarea>
                </div>
              </div>
            </div>

            <!-- ── SECTION 4: ADDRESS & CONTACT WHILE ON LEAVE ─────────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-amber-100 dark:bg-amber-900/50 text-amber-600 flex items-center justify-center text-[10px] font-bold">4</span>
                Physical Address & Emergency Contact During Leave
              </h5>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Physical Address / Destination While on Leave *</label>
                  <input v-model="form.address_while_on_leave" required class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400" placeholder="e.g. Plot 14, Kiwatule Road, Kampala" />
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Emergency Contact Phone *</label>
                  <input v-model="form.emergency_phone" required class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400" placeholder="+256 780 000 000" />
                </div>
              </div>
            </div>

            <!-- ── SECTION 5: REASON & JUSTIFICATION ───────────────────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-rose-100 dark:bg-rose-900/50 text-rose-600 flex items-center justify-center text-[10px] font-bold">5</span>
                Reason & Detailed Justification
              </h5>
              <div>
                <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Reason for Leave *</label>
                <textarea v-model="form.reason" required rows="3" class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400" placeholder="Provide detailed background reason and justification for statutory leave request..."></textarea>
              </div>
            </div>

            <!-- ── SECTION 6: STATUTORY DECLARATION ────────────────────── -->
            <div class="p-4 rounded-xl bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700">
              <label class="flex items-start gap-3 cursor-pointer">
                <input type="checkbox" v-model="declarationAgreed" required class="mt-0.5 rounded text-blue-600 focus:ring-0" />
                <span class="text-slate-700 dark:text-slate-300 text-[11px] leading-relaxed">
                  <strong>Statutory Declaration:</strong> I hereby certify that the information given above is true and complete to the best of my knowledge. I undertake to resume duty on the date specified above and acknowledge that any unauthorized extension without prior written approval from Human Resources constitutes an act of absenteeism under the Public Service Standing Orders and NCS Staff Regulations.
                </span>
              </label>
            </div>

            <!-- Form Submission Bar -->
            <div class="flex items-center justify-between pt-4 border-t border-slate-100 dark:border-slate-800">
              <router-link to="/reception/leave/status" class="btn btn-sm btn-light text-xs">
                Cancel
              </router-link>
              <button type="submit" class="btn btn-sm btn-primary flex items-center gap-2 text-xs shadow-sm" :disabled="submitting || !declarationAgreed">
                <i class="icofont-paper-plane"></i>
                <span v-if="submitting">Submitting Leave Dossier...</span>
                <span v-else>Submit Leave Application to Supervisor</span>
              </button>
            </div>

          </form>
        </div>
      </div>

    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import { useAuthStore } from '@/stores/auth.js'
import { apiGet, apiPost } from '@/api/client'

const authStore = useAuthStore()
const submitting = ref(false)
const successMsg = ref('')
const errorMsg = ref('')
const declarationAgreed = ref(false)
const leavesHistory = ref([])

// ── NCS Department Options ──────────────────────────────────────────────
const departmentOptions = [
  { label: 'IT / ICT Infrastructure', value: 'IT / ICT Infrastructure', subtitle: 'Datacenter, Network, Systems & Web' },
  { label: 'Engineering & Infrastructure', value: 'Engineering & Infrastructure', subtitle: 'Civil Works, Electrical, HVAC, Plumbing' },
  { label: 'Finance & Accounts', value: 'Finance & Accounts', subtitle: 'Ledger, Vote Book, Assets & Disbursals' },
  { label: 'Internal Audit Department', value: 'Internal Audit Department', subtitle: 'Statutory Risk, Compliance & Systems Audit' },
  { label: 'Procurement & Disposal Unit (PDU)', value: 'Procurement & Disposal Unit (PDU)', subtitle: 'PPDA Contracts & Bidding Management' },
  { label: 'Stores & Inventory Unit', value: 'Stores & Inventory Unit', subtitle: 'Physical Inventory, GRN & Warehouse' },
  { label: 'Facilities & Venue Management', value: 'Facilities & Venue Management', subtitle: 'Lugogo Arena, Stadium, Hostels' },
  { label: 'Transport & Fleet Management', value: 'Transport & Fleet Management', subtitle: 'Vehicles, Logistics & Dispatch' },
  { label: 'Sports Science & Medical Unit', value: 'Sports Science & Medical Unit', subtitle: 'Sports Medicine, Physio & Anti-Doping' },
  { label: 'Legal & Corporate Affairs', value: 'Legal & Corporate Affairs', subtitle: 'Contracts, Arbitrations & Compliance' },
  { label: 'Public Relations & Media', value: 'Public Relations & Media', subtitle: 'Communications, Media & Publications' },
  { label: 'Human Resources & Administration', value: 'Human Resources & Administration', subtitle: 'Staff Establishment & Payroll' },
  { label: 'Technical Directorate', value: 'Technical Directorate', subtitle: 'Federations, Competitions & NSMIS' },
  { label: 'Front Desk / Reception', value: 'Front Desk / Reception', subtitle: 'Visitor Reception & Helpdesk' },
  { label: 'Executive Directorate', value: 'Executive Directorate', subtitle: 'General Secretary / Accounting Officer' }
]

// ── Official Designation Options across All NCS Departments ─────────────
const designationOptions = [
  // IT & ICT Infrastructure
  { label: 'ICT Systems & Database Administrator', value: 'ICT Systems & Database Administrator', subtitle: 'IT / ICT Infrastructure' },
  { label: 'IT Support Specialist / Service Desk Lead', value: 'IT Support Specialist / Service Desk Lead', subtitle: 'IT / ICT Infrastructure' },
  { label: 'Front Desk Receptionist', value: 'Front Desk Receptionist', subtitle: 'Front Desk / Reception' },
  
  // Engineering & Infrastructure
  { label: 'Senior Infrastructure Engineer', value: 'Senior Infrastructure Engineer', subtitle: 'Engineering & Infrastructure' },
  { label: 'Assistant Engineer (Civil Works)', value: 'Assistant Engineer (Civil Works)', subtitle: 'Engineering & Infrastructure' },
  { label: 'Assistant Engineer (Electrical & Power)', value: 'Assistant Engineer (Electrical & Power)', subtitle: 'Engineering & Infrastructure' },
  { label: 'Engineering Officer (Civil)', value: 'Engineering Officer (Civil)', subtitle: 'Engineering & Infrastructure' },
  { label: 'Engineering Officer (Electrical)', value: 'Engineering Officer (Electrical)', subtitle: 'Engineering & Infrastructure' },
  { label: 'Plumber & Water Infrastructure Technician', value: 'Plumber & Water Infrastructure Technician', subtitle: 'Engineering & Infrastructure' },

  // Finance & Internal Audit
  { label: 'Senior Accountant / Head of Finance', value: 'Senior Accountant / Head of Finance', subtitle: 'Finance & Accounts' },
  { label: 'Accountant / Payables Lead', value: 'Accountant / Payables Lead', subtitle: 'Finance & Accounts' },
  { label: 'Head of Internal Audit', value: 'Head of Internal Audit', subtitle: 'Internal Audit Department' },
  { label: 'Internal Auditor', value: 'Internal Auditor', subtitle: 'Internal Audit Department' },

  // Human Resources & Administration
  { label: 'Human Resources Manager', value: 'Human Resources Manager', subtitle: 'Human Resources & Administration' },
  { label: 'Human Resources Officer', value: 'Human Resources Officer', subtitle: 'Human Resources & Administration' },
  { label: 'Administrative Assistant', value: 'Administrative Assistant', subtitle: 'Human Resources & Administration' },

  // Procurement & Disposal (PDU)
  { label: 'Head of Procurement & Disposal (PDU)', value: 'Head of Procurement & Disposal (PDU)', subtitle: 'Procurement & Disposal Unit (PDU)' },
  { label: 'Procurement Officer', value: 'Procurement Officer', subtitle: 'Procurement & Disposal Unit (PDU)' },

  // Stores & Logistics
  { label: 'Stores Officer / Inventory Custodian', value: 'Stores Officer / Inventory Custodian', subtitle: 'Stores & Inventory Unit' },
  { label: 'Facilities & Venue Operations Manager', value: 'Facilities & Venue Operations Manager', subtitle: 'Facilities & Venue Management' },
  { label: 'Transport Officer & Fleet Manager', value: 'Transport Officer & Fleet Manager', subtitle: 'Transport & Fleet Management' },
  { label: 'Senior Council Driver', value: 'Senior Council Driver', subtitle: 'Transport & Fleet Management' },

  // Corporate, PR & Legal
  { label: 'Legal Counsel & Compliance Officer', value: 'Legal Counsel & Compliance Officer', subtitle: 'Legal & Corporate Affairs' },
  { label: 'PR & Corporate Communications Officer', value: 'PR & Corporate Communications Officer', subtitle: 'Public Relations & Media' },

  // Sports Science & Technical
  { label: 'Chief Medical Officer / Sports Physician', value: 'Chief Medical Officer / Sports Physician', subtitle: 'Sports Science & Medical Unit' },
  { label: 'Senior Physiotherapist & Rehab Lead', value: 'Senior Physiotherapist & Rehab Lead', subtitle: 'Sports Science & Medical Unit' },
  { label: 'Technical Director', value: 'Technical Director', subtitle: 'Technical Directorate' },
  { label: 'National Athlete Safeguarding Officer', value: 'National Athlete Safeguarding Officer', subtitle: 'Technical Directorate' },

  // Executive Directorate
  { label: 'General Secretary (CEO / Accounting Officer)', value: 'General Secretary (CEO / Accounting Officer)', subtitle: 'Executive Directorate' },
  { label: 'Assistant General Secretary - Technical', value: 'Assistant General Secretary - Technical', subtitle: 'Technical Directorate' },
  { label: 'Assistant General Secretary - Administration', value: 'Assistant General Secretary - Administration', subtitle: 'Administration Directorate' }
]

const leaveTypeOptions = [
  { label: 'Annual Statutory Leave (21 Working Days)', value: 'ANNUAL_LEAVE', subtitle: 'Standard statutory annual entitlement' },
  { label: 'Certified Sick / Medical Leave', value: 'SICK_LEAVE', subtitle: 'Doctor report / medical certificate required' },
  { label: 'Compassionate / Emergency Bereavement Leave', value: 'COMPASSIONATE', subtitle: 'Family emergency / bereavement' },
  { label: 'Capacity Building / Short Study Leave', value: 'STUDY_LEAVE', subtitle: 'Institutional skills & professional development' },
  { label: 'Maternity / Paternity Statutory Leave', value: 'MATERNITY_PATERNITY', subtitle: 'Maternity (60 Days) / Paternity (4 Days)' },
  { label: 'Special Duty Exemption / Unpaid Leave', value: 'SPECIAL_LEAVE', subtitle: 'Approved special circumstances' }
]

const form = ref({
  staff_name: '',
  staff_file_no: '',
  staff_department: 'Front Desk / Reception',
  staff_role: 'Front Desk Receptionist',
  contact_phone: '',
  contact_email: '',
  leave_type: 'ANNUAL_LEAVE',
  start_date: '',
  end_date: '',
  return_date: '',
  days_requested: 1,
  relieving_officer_name: '',
  relieving_officer_role: 'IT Support Specialist / Service Desk Lead',
  duty_handover_details: '',
  address_while_on_leave: '',
  emergency_phone: '',
  reason: ''
})

function onDesignationSelected(desigVal) {
  const match = designationOptions.find(d => d.value === desigVal)
  if (match && match.subtitle) {
    const deptMatch = departmentOptions.find(dep => dep.value === match.subtitle || dep.label === match.subtitle)
    if (deptMatch) {
      form.value.staff_department = deptMatch.value
    }
  }
}

const utilizedDays = computed(() => {
  return leavesHistory.value
    .filter(l => l.status === 'APPROVED_HR' && l.leave_type === 'ANNUAL_LEAVE')
    .reduce((acc, cur) => acc + (cur.days_requested || 0), 0)
})

function calculateDays() {
  if (form.value.start_date && form.value.end_date) {
    const s = new Date(form.value.start_date)
    const e = new Date(form.value.end_date)
    const diff = Math.ceil((e - s) / (1000 * 60 * 60 * 24)) + 1
    if (diff > 0) {
      form.value.days_requested = diff
      // Default return date to day after end date
      const ret = new Date(e)
      ret.setDate(ret.getDate() + 1)
      form.value.return_date = ret.toISOString().split('T')[0]
    }
  }
}

// Auto-detect Department, Designation and Personal Details for Signed-In User
function detectUserPersonalDetails() {
  if (!authStore.user) return

  const u = authStore.user
  const fullName = `${u.first_name || ''} ${u.last_name || ''}`.trim()
  if (fullName) {
    form.value.staff_name = fullName
  }
  if (u.email) {
    form.value.contact_email = u.email
  }

  const rawRoles = u.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)

  if (roles.includes('it_officer') || u.email?.includes('it')) {
    form.value.staff_department = 'IT / ICT Infrastructure'
    form.value.staff_role = 'ICT Systems & Database Administrator'
    form.value.staff_file_no = 'NCS/ICT/2026/017'
    form.value.contact_phone = '+256 701 234 567'
    form.value.emergency_phone = '+256 701 234 567'
    form.value.relieving_officer_name = 'Peter Ssewankambo'
    form.value.relieving_officer_role = 'IT Support Specialist / Service Desk Lead'
  } else if (roles.includes('senior_engineer') || u.email?.includes('seniorengineer')) {
    form.value.staff_department = 'Engineering & Infrastructure'
    form.value.staff_role = 'Senior Infrastructure Engineer'
    form.value.staff_file_no = 'NCS/ENG/2026/007'
    form.value.contact_phone = '+256 772 345 678'
    form.value.emergency_phone = '+256 772 345 678'
    form.value.relieving_officer_name = 'Eng. Isaac Musoke'
    form.value.relieving_officer_role = 'Assistant Engineer (Civil Works)'
  } else if (roles.includes('assistant_engineer_civil') || u.email?.includes('civil.engineer')) {
    form.value.staff_department = 'Engineering & Infrastructure'
    form.value.staff_role = 'Assistant Engineer (Civil Works)'
    form.value.staff_file_no = 'NCS/ENG/2026/008'
    form.value.contact_phone = '+256 772 345 679'
    form.value.emergency_phone = '+256 772 345 679'
    form.value.relieving_officer_name = 'Eng. Patrick Okello'
    form.value.relieving_officer_role = 'Senior Infrastructure Engineer'
  } else if (roles.includes('assistant_engineer_electrical') || u.email?.includes('electrical.engineer')) {
    form.value.staff_department = 'Engineering & Infrastructure'
    form.value.staff_role = 'Assistant Engineer (Electrical & Power)'
    form.value.staff_file_no = 'NCS/ENG/2026/009'
    form.value.contact_phone = '+256 772 345 680'
    form.value.emergency_phone = '+256 772 345 680'
    form.value.relieving_officer_name = 'Brian Ssempala'
    form.value.relieving_officer_role = 'Engineering Officer (Electrical)'
  } else if (roles.includes('engineering_officer_civil') || u.email?.includes('civil.officer')) {
    form.value.staff_department = 'Engineering & Infrastructure'
    form.value.staff_role = 'Engineering Officer (Civil)'
    form.value.staff_file_no = 'NCS/ENG/2026/010'
    form.value.contact_phone = '+256 772 345 681'
    form.value.emergency_phone = '+256 772 345 681'
    form.value.relieving_officer_name = 'Eng. Isaac Musoke'
    form.value.relieving_officer_role = 'Assistant Engineer (Civil Works)'
  } else if (roles.includes('engineering_officer_electrical') || u.email?.includes('electrical.officer')) {
    form.value.staff_department = 'Engineering & Infrastructure'
    form.value.staff_role = 'Engineering Officer (Electrical)'
    form.value.staff_file_no = 'NCS/ENG/2026/011'
    form.value.contact_phone = '+256 772 345 682'
    form.value.emergency_phone = '+256 772 345 682'
    form.value.relieving_officer_name = 'Eng. Denis Kato'
    form.value.relieving_officer_role = 'Assistant Engineer (Electrical & Power)'
  } else if (roles.includes('plumber') || u.email?.includes('plumb')) {
    form.value.staff_department = 'Engineering & Infrastructure'
    form.value.staff_role = 'Plumber & Water Infrastructure Technician'
    form.value.staff_file_no = 'NCS/ENG/2026/012'
    form.value.contact_phone = '+256 772 345 683'
    form.value.emergency_phone = '+256 772 345 683'
    form.value.relieving_officer_name = 'Joseph Mukasa'
    form.value.relieving_officer_role = 'Engineering Officer (Civil)'
  } else if (roles.includes('stores_officer') || u.email?.includes('stores')) {
    form.value.staff_department = 'Stores & Inventory Unit'
    form.value.staff_role = 'Stores Officer / Inventory Custodian'
    form.value.staff_file_no = 'NCS/STR/2026/020'
    form.value.contact_phone = '+256 782 456 789'
    form.value.emergency_phone = '+256 782 456 789'
  } else if (roles.includes('facilities_manager') || u.email?.includes('facilities')) {
    form.value.staff_department = 'Facilities & Venue Management'
    form.value.staff_role = 'Facilities & Venue Operations Manager'
    form.value.staff_file_no = 'NCS/FAC/2026/021'
    form.value.contact_phone = '+256 703 567 890'
    form.value.emergency_phone = '+256 703 567 890'
  } else if (roles.includes('accountant') || u.email?.includes('accountant')) {
    form.value.staff_department = 'Finance & Accounts'
    form.value.staff_role = 'Senior Accountant / Head of Finance'
    form.value.staff_file_no = 'NCS/FIN/2026/013'
    form.value.contact_phone = '+256 774 678 901'
    form.value.emergency_phone = '+256 774 678 901'
  } else if (roles.includes('auditor') || u.email?.includes('auditor')) {
    form.value.staff_department = 'Internal Audit Department'
    form.value.staff_role = 'Head of Internal Audit'
    form.value.staff_file_no = 'NCS/AUD/2026/014'
    form.value.contact_phone = '+256 785 789 012'
    form.value.emergency_phone = '+256 785 789 012'
  } else if (roles.includes('procurement_officer') || u.email?.includes('procurement')) {
    form.value.staff_department = 'Procurement & Disposal Unit (PDU)'
    form.value.staff_role = 'Head of Procurement & Disposal (PDU)'
    form.value.staff_file_no = 'NCS/PDU/2026/018'
    form.value.contact_phone = '+256 706 890 123'
    form.value.emergency_phone = '+256 706 890 123'
  } else if (roles.includes('public_relations') || u.email?.includes('pr')) {
    form.value.staff_department = 'Public Relations & Media'
    form.value.staff_role = 'PR & Corporate Communications Officer'
    form.value.staff_file_no = 'NCS/PRM/2026/019'
    form.value.contact_phone = '+256 777 901 234'
    form.value.emergency_phone = '+256 777 901 234'
  } else if (roles.includes('legal_counsel') || u.email?.includes('legal')) {
    form.value.staff_department = 'Legal & Corporate Affairs'
    form.value.staff_role = 'Legal Counsel & Compliance Officer'
    form.value.staff_file_no = 'NCS/LEG/2026/022'
    form.value.contact_phone = '+256 788 012 345'
    form.value.emergency_phone = '+256 788 012 345'
  } else if (roles.includes('medical_officer') || roles.includes('physiotherapist') || u.email?.includes('medical')) {
    form.value.staff_department = 'Sports Science & Medical Unit'
    form.value.staff_role = 'Chief Medical Officer / Sports Physician'
    form.value.staff_file_no = 'NCS/MED/2026/023'
    form.value.contact_phone = '+256 709 123 456'
    form.value.emergency_phone = '+256 709 123 456'
  } else if (roles.includes('transport_officer') || roles.includes('driver') || u.email?.includes('transport')) {
    form.value.staff_department = 'Transport & Fleet Management'
    form.value.staff_role = 'Transport Officer & Fleet Manager'
    form.value.staff_file_no = 'NCS/TRN/2026/025'
    form.value.contact_phone = '+256 780 234 567'
    form.value.emergency_phone = '+256 780 234 567'
  } else if (roles.includes('human_resources') || roles.includes('hr') || u.email?.includes('hr')) {
    form.value.staff_department = 'Human Resources & Administration'
    form.value.staff_role = 'Human Resources Manager'
    form.value.staff_file_no = 'NCS/HRM/2026/015'
    form.value.contact_phone = '+256 771 345 678'
    form.value.emergency_phone = '+256 771 345 678'
  } else if (roles.includes('technical_department') || roles.includes('ags_technical') || u.email?.includes('technical')) {
    form.value.staff_department = 'Technical Directorate'
    form.value.staff_role = 'Technical Director'
    form.value.staff_file_no = 'NCS/TEC/2026/006'
    form.value.contact_phone = '+256 702 456 789'
    form.value.emergency_phone = '+256 702 456 789'
  } else if (roles.includes('general_secretary') || u.email?.includes('gs')) {
    form.value.staff_department = 'Executive Directorate'
    form.value.staff_role = 'General Secretary (CEO / Accounting Officer)'
    form.value.staff_file_no = 'NCS/EXEC/2026/001'
    form.value.contact_phone = '+256 773 567 890'
    form.value.emergency_phone = '+256 773 567 890'
  } else {
    // Default Receptionist / Help Desk
    form.value.staff_department = 'Front Desk / Reception'
    form.value.staff_role = 'Front Desk Receptionist'
    form.value.staff_file_no = 'NCS/REC/2026/016'
    form.value.contact_phone = '+256 704 678 901'
    form.value.emergency_phone = '+256 704 678 901'
    form.value.relieving_officer_name = 'Allan Tumusiime'
    form.value.relieving_officer_role = 'ICT Systems & Database Administrator'
  }
}

async function fetchLeavesHistory() {
  try {
    const res = await apiGet('/reception/leave/status')
    leavesHistory.value = res.data?.leaves || res.leaves || []
  } catch {
    leavesHistory.value = []
  }
}

async function submitLeaveApplication() {
  if (!declarationAgreed.value) {
    errorMsg.value = 'Please review and accept the statutory declaration before submitting.'
    return
  }
  submitting.value = true
  successMsg.value = ''
  errorMsg.value = ''
  try {
    const res = await apiPost('/reception/leave/apply', form.value)
    successMsg.value = res.data?.message || res.message || 'Leave application submitted successfully for supervisor vetting!'
    form.value.reason = ''
    form.value.duty_handover_details = ''
    form.value.address_while_on_leave = ''
    declarationAgreed.value = false
    await fetchLeavesHistory()
  } catch (err) {
    errorMsg.value = err.message || 'Could not submit leave application'
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  detectUserPersonalDetails()
  await fetchLeavesHistory()
})
</script>
