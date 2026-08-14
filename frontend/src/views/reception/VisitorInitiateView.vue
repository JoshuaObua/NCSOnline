<template>
  <LayoutDefault title="Initiate Visitor Clearance">
    <div class="space-y-6 pb-12 w-full">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200 w-full">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-id-card text-xs"></i> Front Desk Security Protocol
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Initiate Visitor Clearance Request</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
            Log visitor credentials and dispatch access request to the host department for clearance authorization.
          </p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <router-link to="/reception/visitors" class="btn btn-sm btn-light flex items-center gap-1.5">
            <i class="icofont-arrow-left"></i> Back to Visitor Registry
          </router-link>
        </div>
      </div>

      <!-- Main Full-Width Initiation Form -->
      <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 w-full mb-0">
        <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-6 flex items-center justify-between">
          <h4 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
            <i class="icofont-badge text-blue-600"></i> Official Visitor Entry Intake Form
          </h4>
          <span class="text-xs text-slate-400 font-mono">NCS-SEC-CLR-01</span>
        </div>

        <div class="card-body p-6 sm:p-8">
          <form @submit.prevent="submitVisitorRequest" class="space-y-6 w-full text-xs">
            
            <div v-if="successMsg" class="p-3.5 rounded-xl bg-emerald-50 dark:bg-emerald-950/50 border border-emerald-200 dark:border-emerald-800 text-emerald-700 dark:text-emerald-300 font-medium text-xs">
              <i class="icofont-check-circled mr-1"></i> {{ successMsg }}
            </div>
            <div v-if="errorMsg" class="p-3.5 rounded-xl bg-red-50 dark:bg-red-950/50 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 font-medium text-xs">
              <i class="icofont-warning mr-1"></i> {{ errorMsg }}
            </div>

            <!-- ── SECTION 1: VISITOR IDENTIFICATION ────────────────────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-blue-100 dark:bg-blue-900/50 text-blue-600 flex items-center justify-center text-[10px] font-bold">1</span>
                Visitor Particulars & Identification
              </h5>
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <div class="sm:col-span-2">
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Visitor Full Name *</label>
                  <input
                    v-model="form.visitor_name"
                    required
                    type="text"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                    placeholder="e.g. John Bosco Mukasa"
                  />
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Telephone / Mobile *</label>
                  <input
                    v-model="form.visitor_phone"
                    required
                    type="text"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                    placeholder="+256 700 000 000"
                  />
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">NIN / ID Number *</label>
                  <input
                    v-model="form.visitor_id_number"
                    required
                    type="text"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400 font-mono"
                    placeholder="CM98024103ABCD"
                  />
                </div>
                <div class="sm:col-span-2">
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Organization / Federation / Entity</label>
                  <input
                    v-model="form.visitor_organization"
                    type="text"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                    placeholder="e.g. Uganda Athletics Federation / Ministry of Education & Sports"
                  />
                </div>
              </div>
            </div>

            <!-- ── SECTION 2: DESTINATION & HOST OFFICER ────────────────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-emerald-100 dark:bg-emerald-900/50 text-emerald-600 flex items-center justify-center text-[10px] font-bold">2</span>
                Destination Office & Host Officer Particulars
              </h5>
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Target Department / Office *</label>
                  <SearchableSelect
                    v-model="form.target_department"
                    :options="departmentOptions"
                    placeholder="Search or select department..."
                    search-placeholder="Filter departments..."
                    required
                  />
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Host Officer Name / Post *</label>
                  <input
                    v-model="form.host_officer_name"
                    required
                    type="text"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                    placeholder="e.g. Dr. Bernard Ogwel (General Secretary)"
                  />
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Primary Purpose of Visit *</label>
                  <SearchableSelect
                    v-model="form.purpose_of_visit"
                    :options="purposeOptions"
                    placeholder="Search or select purpose..."
                    search-placeholder="Filter visit purposes..."
                    required
                  />
                </div>
              </div>
            </div>

            <!-- ── SECTION 3: SECURITY & PROPERTY DECLARED ──────────────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-purple-100 dark:bg-purple-900/50 text-purple-600 flex items-center justify-center text-[10px] font-bold">3</span>
                Security Gate Ingestion & Property Registration
              </h5>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Vehicle Registration No. (If driving into premises)</label>
                  <input
                    v-model="form.vehicle_reg_no"
                    type="text"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white font-mono placeholder-slate-400"
                    placeholder="e.g. UBA 452Y / None (On Foot)"
                  />
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Declared Hardware / Electronic Devices</label>
                  <input
                    v-model="form.items_declared"
                    type="text"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                    placeholder="e.g. HP ProBook Laptop S/N 5CD928, Camera"
                  />
                </div>
              </div>
            </div>

            <!-- ── SECTION 4: APPOINTMENT BRIEF / NOTES ────────────────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-amber-100 dark:bg-amber-900/50 text-amber-600 flex items-center justify-center text-[10px] font-bold">4</span>
                Appointment Notes & Front Desk Dispatch Brief
              </h5>
              <div>
                <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Appointment Brief / Specific Requests</label>
                <textarea
                  v-model="form.remarks"
                  rows="3"
                  class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                  placeholder="Enter any additional details regarding scheduled appointment time, accompanying delegation members, or special access instructions..."
                ></textarea>
              </div>
            </div>

            <!-- Form Submission Bar -->
            <div class="flex items-center justify-between pt-4 border-t border-slate-100 dark:border-slate-800">
              <router-link to="/reception/visitors" class="btn btn-sm btn-light text-xs">
                Cancel & Back
              </router-link>
              <button
                type="submit"
                class="btn btn-sm btn-primary flex items-center gap-2 text-xs shadow-sm"
                :disabled="submitting"
              >
                <i class="icofont-paper-plane"></i>
                <span v-if="submitting">Transmitting Request to Host...</span>
                <span v-else>Submit & Dispatch to Host Officer</span>
              </button>
            </div>

          </form>
        </div>
      </div>

    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import { apiPost } from '@/api/client'

const router = useRouter()
const submitting = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const departmentOptions = [
  { label: 'Office of the General Secretary', value: 'Office of the General Secretary', subtitle: 'Executive Leadership & Governance' },
  { label: 'Administration & HR', value: 'Administration & HR', subtitle: 'AGS Admin & Human Resources' },
  { label: 'Technical & Sports Development', value: 'Technical & Sports Development', subtitle: 'AGS Technical & Federations' },
  { label: 'Finance & Accounts', value: 'Finance & Accounts', subtitle: 'Subventions, Budget & Disbursement' },
  { label: 'Engineering & Infrastructure', value: 'Engineering & Infrastructure', subtitle: 'Facilities & Capital Works' },
  { label: 'Facilities & Venues Management', value: 'Facilities & Venues Management', subtitle: 'Lugogo Complex, Arenas & Hostels' },
  { label: 'Transport & Fleet Management', value: 'Transport & Fleet Management', subtitle: 'Institutional Logistics & Vehicles' },
  { label: 'Sports Medicine & Science / WADA', value: 'Sports Medicine & Science / WADA', subtitle: 'Athlete Health & Anti-Doping' },
  { label: 'Legal & Corporate Affairs', value: 'Legal & Corporate Affairs', subtitle: 'Statutory Compliance & Contracts' },
  { label: 'Procurement & Disposal', value: 'Procurement & Disposal', subtitle: 'Supplies & Tender Operations' },
  { label: 'Internal Audit', value: 'Internal Audit', subtitle: 'Assurance & Governance Review' }
]

const purposeOptions = [
  { label: 'Official Meeting / Consultation', value: 'Official Meeting / Consultation' },
  { label: 'Statutory Submission / Licence Review', value: 'Statutory Submission / Licence Review' },
  { label: 'Subvention & Accountability Audit', value: 'Subvention & Accountability Audit' },
  { label: 'Lugogo Venue / Facility Inspection', value: 'Lugogo Venue / Facility Inspection' },
  { label: 'National Team Travel Briefing', value: 'National Team Travel Briefing' },
  { label: 'Dispatch / Official Delivery', value: 'Dispatch / Official Delivery' },
  { label: 'General Inquiries / Citizen Service', value: 'General Inquiries / Citizen Service' }
]

const form = ref({
  visitor_name: '',
  visitor_phone: '',
  visitor_id_number: '',
  visitor_organization: '',
  target_department: 'Office of the General Secretary',
  host_officer_name: '',
  purpose_of_visit: 'Official Meeting / Consultation',
  vehicle_reg_no: '',
  items_declared: '',
  remarks: ''
})

async function submitVisitorRequest() {
  submitting.value = true
  successMsg.value = ''
  errorMsg.value = ''
  try {
    const res = await apiPost('/reception/visitors', form.value)
    successMsg.value = res.data?.message || res.message || 'Visitor entry clearance request initiated successfully!'
    
    // Redirect to visitors registry after brief pause
    setTimeout(() => {
      router.push('/reception/visitors')
    }, 1200)
  } catch (err) {
    errorMsg.value = err.message || 'Could not initiate visitor clearance request'
    submitting.value = false
  }
}
</script>
