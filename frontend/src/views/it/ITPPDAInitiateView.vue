<template>
  <LayoutDefault title="PPDA Form 5 Requisition">
    <div class="space-y-6 pb-12 w-full">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200 w-full">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-document-folder text-xs"></i> Public Procurement (PPDA Act 2003 / 2023)
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">PPDA Form 5: Request for Approval of Procurement</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
            Official statutory user department procurement requisition dossier for goods, services, consultancy, and works across National Council of Sports directorates.
          </p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <router-link to="/it/ppda/status" class="btn btn-sm btn-light flex items-center gap-1.5">
            <i class="icofont-history"></i> View Requisitions & Statuses
          </router-link>
        </div>
      </div>

      <!-- Main Full-Width Form Card -->
      <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 w-full mb-0">
        <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-6 flex items-center justify-between">
          <h4 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
            <i class="icofont-law-order text-blue-600"></i> User Department Procurement Initiation (PPDA Form 5)
          </h4>
          <span class="text-xs text-slate-400 font-mono">PPDA-FORM-5 / NCS-PROC</span>
        </div>

        <div class="card-body p-6 sm:p-8">
          <form @submit.prevent="submitRequisition" class="space-y-7 w-full text-xs">
            
            <div v-if="successMsg" class="p-3.5 rounded-xl bg-emerald-50 dark:bg-emerald-950/50 border border-emerald-200 dark:border-emerald-800 text-emerald-700 dark:text-emerald-300 font-medium text-xs">
              <i class="icofont-check-circled mr-1"></i> {{ successMsg }}
            </div>
            <div v-if="errorMsg" class="p-3.5 rounded-xl bg-red-50 dark:bg-red-950/50 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 font-medium text-xs">
              <i class="icofont-warning mr-1"></i> {{ errorMsg }}
            </div>

            <!-- ── SECTION 1: REQUISITIONING ENTITY & OFFICER ───────────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-blue-100 dark:bg-blue-900/50 text-blue-600 flex items-center justify-center text-[10px] font-bold">1</span>
                User Department & Requisitioning Officer Particulars
              </h5>
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Requisitioning Officer Name *</label>
                  <input
                    v-model="form.officer_name"
                    required
                    type="text"
                    placeholder="Enter full officer name..."
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white"
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Official Designation *</label>
                  <SearchableSelect
                    v-model="form.designation"
                    :options="designationOptions"
                    placeholder="Select official designation..."
                    search-placeholder="Search designations & roles..."
                    required
                    @update:modelValue="onDesignationSelected"
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">User Department / Unit *</label>
                  <SearchableSelect
                    v-model="form.department"
                    :options="departmentOptions"
                    placeholder="Select user department..."
                    search-placeholder="Search NCS departments..."
                    required
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Contact Phone / Extension</label>
                  <input
                    v-model="form.contact_phone"
                    type="text"
                    placeholder="e.g. +256 701 234 567 / Ext. 204"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Statutory Financial Year *</label>
                  <SearchableSelect
                    v-model="form.financial_year"
                    :options="financialYearOptions"
                    placeholder="Select financial year..."
                    search-placeholder="Search financial years..."
                    required
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Target Delivery / Completion Date *</label>
                  <input
                    v-model="form.required_delivery_date"
                    required
                    type="date"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white"
                  />
                </div>
              </div>
            </div>

            <!-- ── SECTION 2: PROCUREMENT SUBJECT, METHOD & FUNDING ────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-emerald-100 dark:bg-emerald-900/50 text-emerald-600 flex items-center justify-center text-[10px] font-bold">2</span>
                Procurement Subject, Method, Category & Vote Allocation
              </h5>
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                <div class="sm:col-span-2">
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Subject / Title of Procurement *</label>
                  <input
                    v-model="form.subject_of_procurement"
                    required
                    type="text"
                    placeholder="e.g. Procurement of Enterprise Core Optical Fiber Switches & SFP+ Modules for Lugogo Datacenter"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Procurement Category *</label>
                  <SearchableSelect
                    v-model="form.procurement_category"
                    :options="categoryOptions"
                    placeholder="Select procurement category..."
                    search-placeholder="Search categories..."
                    required
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Proposed Procurement Method *</label>
                  <SearchableSelect
                    v-model="form.procurement_method"
                    :options="procurementMethodOptions"
                    placeholder="Select procurement method..."
                    search-placeholder="Search PPDA methods..."
                    required
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Approved Budget Vote Head / Line *</label>
                  <SearchableSelect
                    v-model="form.budget_vote_head"
                    :options="voteHeadOptions"
                    placeholder="Select budget vote line..."
                    search-placeholder="Search budget vote lines..."
                    required
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Source of Funding *</label>
                  <SearchableSelect
                    v-model="form.source_of_funds"
                    :options="fundingSourceOptions"
                    placeholder="Select funding source..."
                    search-placeholder="Search funding sources..."
                    required
                  />
                </div>

                <div class="sm:col-span-2">
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Designated Delivery Location / Facility *</label>
                  <SearchableSelect
                    v-model="form.delivery_location"
                    :options="deliveryLocationOptions"
                    placeholder="Select delivery location..."
                    search-placeholder="Search delivery facilities..."
                    required
                  />
                </div>

                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Warranty & Service SLA Requirement *</label>
                  <SearchableSelect
                    v-model="form.warranty_requirement"
                    :options="warrantyOptions"
                    placeholder="Select warranty terms..."
                    search-placeholder="Search warranty terms..."
                    required
                  />
                </div>
              </div>
            </div>

            <!-- ── SECTION 3: DETAILED LINE ITEMS SCHEDULE ──────────────── -->
            <div>
              <div class="flex items-center justify-between mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800">
                <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider flex items-center gap-2">
                  <span class="w-5 h-5 rounded-full bg-purple-100 dark:bg-purple-900/50 text-purple-600 flex items-center justify-center text-[10px] font-bold">3</span>
                  Schedule of Requirements & Estimated Cost Breakdown (UGX)
                </h5>
                <button
                  type="button"
                  @click="addItem"
                  class="btn btn-xs btn-outline-primary flex items-center gap-1"
                >
                  <i class="icofont-plus"></i> Add Item Line
                </button>
              </div>

              <div class="overflow-x-auto border border-slate-200 dark:border-slate-700 rounded-xl bg-slate-50/50 dark:bg-slate-800/40 p-1 mb-3">
                <table class="table table-sm text-xs mb-0">
                  <thead class="bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 font-semibold uppercase text-[10px]">
                    <tr>
                      <th class="py-2 px-3 w-5/12">Item Description & Specifications</th>
                      <th class="py-2 px-3 w-2/12">Quantity</th>
                      <th class="py-2 px-3 w-2/12">Unit of Measure</th>
                      <th class="py-2 px-3 w-2/12">Est. Unit Price (UGX)</th>
                      <th class="py-2 px-3 w-2/12 text-right">Total Price (UGX)</th>
                      <th class="py-2 px-2 w-10 text-center">Action</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
                    <tr v-for="(item, idx) in form.items" :key="idx">
                      <td class="p-2">
                        <input
                          v-model="item.item_name"
                          required
                          type="text"
                          placeholder="e.g. Cisco Catalyst 9300 48-Port Switch"
                          class="form-control text-xs w-full py-1 px-2.5 rounded border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white"
                        />
                      </td>
                      <td class="p-2">
                        <input
                          v-model.number="item.quantity"
                          @input="updateLineTotal(item)"
                          required
                          type="number"
                          min="1"
                          class="form-control text-xs w-full py-1 px-2.5 rounded border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white font-mono"
                        />
                      </td>
                      <td class="p-2">
                        <input
                          v-model="item.unit"
                          required
                          type="text"
                          placeholder="Unit/Pieces/License/Roll/Lot"
                          class="form-control text-xs w-full py-1 px-2.5 rounded border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white"
                        />
                      </td>
                      <td class="p-2">
                        <input
                          v-model.number="item.unit_price_ugx"
                          @input="updateLineTotal(item)"
                          required
                          type="number"
                          min="0"
                          step="500"
                          class="form-control text-xs w-full py-1 px-2.5 rounded border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white font-mono"
                        />
                      </td>
                      <td class="p-2 text-right font-mono font-bold text-slate-900 dark:text-white">
                        UGX {{ Number(item.total_price_ugx || 0).toLocaleString() }}
                      </td>
                      <td class="p-2 text-center">
                        <button
                          type="button"
                          @click="removeItem(idx)"
                          :disabled="form.items.length <= 1"
                          class="text-red-500 hover:text-red-700 p-1 disabled:opacity-30"
                          title="Remove item"
                        >
                          <i class="icofont-trash"></i>
                        </button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <!-- Grand Total Summary Bar -->
              <div class="flex items-center justify-between p-3.5 rounded-xl bg-blue-50 dark:bg-blue-950/40 border border-blue-200 dark:border-blue-800">
                <span class="font-bold text-blue-900 dark:text-blue-300 uppercase tracking-wide text-xs">
                  Estimated Total Procurement Value (UGX):
                </span>
                <span class="text-base font-black text-blue-700 dark:text-blue-400 font-mono">
                  UGX {{ grandTotal.toLocaleString() }}
                </span>
              </div>
            </div>

            <!-- ── SECTION 4: TECHNICAL SPECS & JUSTIFICATION ──────────── -->
            <div>
              <h5 class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3 pb-1.5 border-b border-slate-100 dark:border-slate-800 flex items-center gap-2">
                <span class="w-5 h-5 rounded-full bg-amber-100 dark:bg-amber-900/50 text-amber-600 flex items-center justify-center text-[10px] font-bold">4</span>
                Technical Specifications & Procurement Justification
              </h5>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Detailed Technical Specifications / TORs *</label>
                  <textarea
                    v-model="form.technical_specifications"
                    required
                    rows="3"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                    placeholder="Provide detailed hardware specifications, capacity, compliance standards, power specs, warranty terms..."
                  ></textarea>
                </div>
                <div>
                  <label class="block font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Justification & Expected Institutional Output *</label>
                  <textarea
                    v-model="form.justification"
                    required
                    rows="3"
                    class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
                    placeholder="Explain institutional need, impact on NCS operations, and alignment with council annual work plan..."
                  ></textarea>
                </div>
              </div>
            </div>

            <!-- ── SECTION 5: STATUTORY SUBMISSION DECLARATION ─────────── -->
            <div class="p-4 rounded-xl bg-slate-50 dark:bg-slate-800/60 border border-slate-200 dark:border-slate-700">
              <label class="flex items-start gap-3 cursor-pointer">
                <input type="checkbox" v-model="declarationAgreed" required class="mt-0.5 rounded text-blue-600 focus:ring-0" />
                <span class="text-slate-700 dark:text-slate-300 text-[11px] leading-relaxed">
                  <strong>PPDA Compliance Declaration:</strong> I hereby certify that the goods/services requested above are urgently required for the official operations of the National Council of Sports, the technical specifications are non-proprietary and objective in compliance with the PPDA Act (2003 as amended), and sufficient budget allocation exists under the designated vote head.
                </span>
              </label>
            </div>

            <!-- Form Submission Bar -->
            <div class="flex items-center justify-between pt-4 border-t border-slate-100 dark:border-slate-800">
              <router-link to="/it/ppda/status" class="btn btn-sm btn-light text-xs">
                Cancel
              </router-link>
              <button
                type="submit"
                class="btn btn-sm btn-primary flex items-center gap-2 text-xs shadow-sm"
                :disabled="submitting || !declarationAgreed"
              >
                <i class="icofont-paper-plane"></i>
                <span v-if="submitting">Transmitting PPDA Form 5...</span>
                <span v-else>Submit Requisition for HOD Endorsement</span>
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
import { useRouter } from 'vue-router'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import { useAuthStore } from '@/stores/auth.js'
import { apiPost } from '@/api/client'

const router = useRouter()
const authStore = useAuthStore()
const submitting = ref(false)
const successMsg = ref('')
const errorMsg = ref('')
const declarationAgreed = ref(false)

// ── Official Designation Options across All NCS Departments ─────────────
const designationOptions = [
  // IT & ICT Infrastructure
  { label: 'ICT Systems & Database Administrator', value: 'ICT Systems & Database Administrator', subtitle: 'IT / ICT Infrastructure' },
  { label: 'IT Support Specialist / Service Desk Lead', value: 'IT Support Specialist / Service Desk Lead', subtitle: 'IT / ICT Infrastructure' },
  
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
  { label: 'Technical Director', value: 'Technical Director', subtitle: 'Technical & Sports Administration' },
  { label: 'National Athlete Safeguarding Officer', value: 'National Athlete Safeguarding Officer', subtitle: 'Safeguarding & Welfare' },

  // Executive Directorate
  { label: 'General Secretary (CEO / Accounting Officer)', value: 'General Secretary (CEO / Accounting Officer)', subtitle: 'Executive Directorate' },
  { label: 'Assistant General Secretary - Technical', value: 'Assistant General Secretary - Technical', subtitle: 'Technical Directorate' },
  { label: 'Assistant General Secretary - Administration', value: 'Assistant General Secretary - Administration', subtitle: 'Administration Directorate' }
]

// ── NCS Departments / Units ─────────────────────────────────────────────
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
  { label: 'Executive Directorate', value: 'Executive Directorate', subtitle: 'General Secretary / Accounting Officer' }
]

// ── Financial Year Options ──────────────────────────────────────────────
const financialYearOptions = [
  { label: 'FY 2026/2027 (Current Active Financial Year)', value: 'FY 2026/2027', subtitle: 'Active Statutory Budget Year' },
  { label: 'FY 2027/2028 (Next Budget Procurement Plan)', value: 'FY 2027/2028', subtitle: 'Advance Procurement Plan' },
  { label: 'FY 2025/2026 (Previous Financial Year Adjustments)', value: 'FY 2025/2026', subtitle: 'Carryover Commitments' }
]

// ── Procurement Category Options ────────────────────────────────────────
const categoryOptions = [
  { label: 'Supplies / ICT Hardware & Equipment', value: 'SUPPLIES', subtitle: 'Servers, PCs, networking, switches, UPS, accessories' },
  { label: 'Services (Non-Consultancy / Subscriptions)', value: 'SERVICES', subtitle: 'Cloud software, M365, Internet, maintenance SLAs' },
  { label: 'Consultancy Services', value: 'CONSULTANCY', subtitle: 'System audits, cybersecurity assessments, strategy' },
  { label: 'Works / Civil & Structured Cabling', value: 'WORKS', subtitle: 'Datacenter fit-out, optical fiber trenching, power works' }
]

// ── Proposed Procurement Method Options ─────────────────────────────────
const procurementMethodOptions = [
  { label: 'Request for Quotations (RFQ - UGX 5M to 50M)', value: 'Request for Quotations (RFQ)', subtitle: 'Standard competitive quotations threshold' },
  { label: 'Micro-Procurement (Below UGX 5,000,000)', value: 'Micro-Procurement', subtitle: 'Direct simplified purchase for low value items' },
  { label: 'Open Domestic Bidding (Above UGX 50,000,000)', value: 'Open Domestic Bidding', subtitle: 'National competitive gazetted bidding' },
  { label: 'Restricted Domestic Bidding', value: 'Restricted Domestic Bidding', subtitle: 'Pre-qualified specialized vendors roster' },
  { label: 'Direct Procurement / Sole Source', value: 'Direct Procurement', subtitle: 'Statutory emergency or proprietary single provider' },
  { label: 'Framework Contract Call-Off Order', value: 'Framework Contract Call-Off Order', subtitle: 'Existing pre-negotiated council framework agreement' }
]

// ── Budget Vote Head Options ────────────────────────────────────────────
const voteHeadOptions = [
  { label: 'Vote 202 - ICT Capital & Datacenter Infrastructure', value: 'Vote 202 - ICT Capital & Datacenter Infrastructure', subtitle: 'Capital equipment & enterprise hardware' },
  { label: 'Vote 203 - Software Licences & Cloud Subscriptions', value: 'Vote 203 - Software Licences & Cloud Subscriptions', subtitle: 'M365, database licenses, SSL certificates' },
  { label: 'Vote 204 - Hardware Maintenance & Network Servicing', value: 'Vote 204 - Hardware Maintenance & Network Servicing', subtitle: 'Repairs, replacement parts, consumables' },
  { label: 'Vote 205 - Internet & Telecommunications Connectivity', value: 'Vote 205 - Internet & Telecommunications Connectivity', subtitle: 'Dedicated optical fiber broadband & VPN' },
  { label: 'Vote 201 - General Secretariat Administration & Operations', value: 'Vote 201 - General Secretariat Administration & Operations', subtitle: 'Operational supplies & council administration' },
  { label: 'Vote 206 - Sports Infrastructure Capital Development', value: 'Vote 206 - Sports Infrastructure Capital Development', subtitle: 'Arena & stadium civil engineering works' },
  { label: 'Vote 207 - Transport, Fleet Maintenance & Fuel', value: 'Vote 207 - Transport, Fleet Maintenance & Fuel', subtitle: 'Vehicle servicing, repairs, tyres & fuel' },
  { label: 'Vote 208 - Medical & Anti-Doping Laboratory Supplies', value: 'Vote 208 - Medical & Anti-Doping Laboratory Supplies', subtitle: 'Rehab equipment, testing kits, pharmaceuticals' }
]

// ── Source of Funds Options ─────────────────────────────────────────────
const fundingSourceOptions = [
  { label: 'GoU Statutory Subvention (Consolidated Fund)', value: 'GoU Statutory Subvention', subtitle: 'Ministry of Finance voted allocation' },
  { label: 'Appropriation in Aid (AIA) / Internally Generated Revenue', value: 'Appropriation in Aid (AIA)', subtitle: 'Venues, licensing & commercial fees' },
  { label: 'Development Partner / Grant Funding', value: 'Development Partner / Grant Funding', subtitle: 'External donor or international sports body' },
  { label: 'Special Sports Development Fund', value: 'Special Sports Development Fund', subtitle: 'Designated national sports development budget' }
]

// ── Delivery Location Options ───────────────────────────────────────────
const deliveryLocationOptions = [
  { label: 'NCS Headquarters Lugogo - ICT Datacenter & Server Room', value: 'NCS Headquarters Lugogo - ICT Datacenter', subtitle: 'Lugogo Sports Complex Administration Wing' },
  { label: 'NCS Lugogo - Main Central Stores & Warehouse', value: 'NCS Lugogo - Main Central Stores', subtitle: 'Physical receiving & GRN verification' },
  { label: 'Lugogo Indoor Arena Complex & Control Room', value: 'Lugogo Indoor Arena Complex', subtitle: 'Arena stadium & broadcast control room' },
  { label: 'Lugogo Sports Hostel & Athlete Village', value: 'Lugogo Sports Hostel', subtitle: 'Residential accommodation wing' },
  { label: 'Kyambogo Sports Complex Facility', value: 'Kyambogo Sports Complex', subtitle: 'National training grounds' }
]

// ── Warranty Options ────────────────────────────────────────────────────
const warrantyOptions = [
  { label: '1 Year Comprehensive OEM Manufacturer Warranty', value: '1 Year Comprehensive OEM Manufacturer Warranty', subtitle: 'Standard hardware replacement & support' },
  { label: '2 Years Comprehensive Warranty with 24/7 SLA Support', value: '2 Years Comprehensive Warranty with 24/7 SLA Support', subtitle: 'Mission-critical enterprise infrastructure' },
  { label: '3 Years Extended Enterprise Hardware Warranty', value: '3 Years Extended Enterprise Hardware Warranty', subtitle: 'Datacenter servers & core networking' },
  { label: 'Annual Software Maintenance & Patching Subscription', value: 'Annual Software Maintenance & Patching Subscription', subtitle: 'Cloud software & SaaS platform licenses' },
  { label: 'Standard Delivery & Commissioning Testing (No Extended SLA)', value: 'Standard Delivery & Commissioning Testing', subtitle: 'One-off supplies & consumable goods' }
]

const form = ref({
  officer_name: 'Allan Tumusiime',
  designation: 'ICT Systems & Database Administrator',
  department: 'IT / ICT Infrastructure',
  contact_phone: '+256 701 234 567',
  financial_year: 'FY 2026/2027',
  subject_of_procurement: '',
  procurement_category: 'SUPPLIES',
  procurement_method: 'Request for Quotations (RFQ)',
  budget_vote_head: 'Vote 202 - ICT Capital & Datacenter Infrastructure',
  source_of_funds: 'GoU Statutory Subvention',
  delivery_location: 'NCS Headquarters Lugogo - ICT Datacenter',
  warranty_requirement: '1 Year Comprehensive OEM Manufacturer Warranty',
  required_delivery_date: '',
  technical_specifications: '',
  justification: '',
  items: [
    { item_name: '', quantity: 1, unit: 'Unit', unit_price_ugx: 0, total_price_ugx: 0 }
  ]
})

function onDesignationSelected(desigVal) {
  const match = designationOptions.find(d => d.value === desigVal)
  if (match && match.subtitle) {
    const deptMatch = departmentOptions.find(dep => dep.value === match.subtitle || dep.label === match.subtitle)
    if (deptMatch) {
      form.value.department = deptMatch.value
    }
  }
}

function addItem() {
  form.value.items.push({
    item_name: '',
    quantity: 1,
    unit: 'Unit',
    unit_price_ugx: 0,
    total_price_ugx: 0
  })
}

function removeItem(idx) {
  if (form.value.items.length > 1) {
    form.value.items.splice(idx, 1)
  }
}

function updateLineTotal(item) {
  item.total_price_ugx = (item.quantity || 0) * (item.unit_price_ugx || 0)
}

const grandTotal = computed(() => {
  return form.value.items.reduce((sum, item) => sum + (item.total_price_ugx || 0), 0)
})

// Auto-detect Department and Designation for Signed-In User
function detectUserDepartmentAndRole() {
  if (!authStore.user) return

  const u = authStore.user
  const fullName = `${u.first_name || ''} ${u.last_name || ''}`.trim()
  if (fullName) {
    form.value.officer_name = fullName
  }

  const rawRoles = u.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)

  if (roles.includes('it_officer') || u.email?.includes('it')) {
    form.value.department = 'IT / ICT Infrastructure'
    form.value.designation = 'ICT Systems & Database Administrator'
    form.value.budget_vote_head = 'Vote 202 - ICT Capital & Datacenter Infrastructure'
  } else if (roles.includes('senior_engineer') || u.email?.includes('seniorengineer')) {
    form.value.department = 'Engineering & Infrastructure'
    form.value.designation = 'Senior Infrastructure Engineer'
    form.value.budget_vote_head = 'Vote 206 - Sports Infrastructure Capital Development'
    form.value.delivery_location = 'Lugogo Indoor Arena Complex'
  } else if (roles.includes('assistant_engineer_civil') || u.email?.includes('civil.engineer')) {
    form.value.department = 'Engineering & Infrastructure'
    form.value.designation = 'Assistant Engineer (Civil Works)'
    form.value.budget_vote_head = 'Vote 206 - Sports Infrastructure Capital Development'
    form.value.delivery_location = 'Lugogo Indoor Arena Complex'
  } else if (roles.includes('assistant_engineer_electrical') || u.email?.includes('electrical.engineer')) {
    form.value.department = 'Engineering & Infrastructure'
    form.value.designation = 'Assistant Engineer (Electrical & Power)'
    form.value.budget_vote_head = 'Vote 206 - Sports Infrastructure Capital Development'
    form.value.delivery_location = 'Lugogo Indoor Arena Complex'
  } else if (roles.includes('engineering_officer_civil') || u.email?.includes('civil.officer')) {
    form.value.department = 'Engineering & Infrastructure'
    form.value.designation = 'Engineering Officer (Civil)'
    form.value.budget_vote_head = 'Vote 206 - Sports Infrastructure Capital Development'
    form.value.delivery_location = 'NCS Lugogo - Main Central Stores'
  } else if (roles.includes('engineering_officer_electrical') || u.email?.includes('electrical.officer')) {
    form.value.department = 'Engineering & Infrastructure'
    form.value.designation = 'Engineering Officer (Electrical)'
    form.value.budget_vote_head = 'Vote 206 - Sports Infrastructure Capital Development'
    form.value.delivery_location = 'NCS Lugogo - Main Central Stores'
  } else if (roles.includes('plumber') || u.email?.includes('plumb')) {
    form.value.department = 'Engineering & Infrastructure'
    form.value.designation = 'Plumber & Water Infrastructure Technician'
    form.value.budget_vote_head = 'Vote 206 - Sports Infrastructure Capital Development'
    form.value.delivery_location = 'Lugogo Sports Hostel & Athlete Village'
  } else if (roles.includes('stores_officer') || u.email?.includes('stores')) {
    form.value.department = 'Stores & Inventory Unit'
    form.value.designation = 'Stores Officer / Inventory Custodian'
    form.value.budget_vote_head = 'Vote 201 - General Secretariat Administration & Operations'
  } else if (roles.includes('facilities_manager') || u.email?.includes('facilities')) {
    form.value.department = 'Facilities & Venue Management'
    form.value.designation = 'Facilities & Venue Operations Manager'
    form.value.budget_vote_head = 'Vote 206 - Sports Infrastructure Capital Development'
  } else if (roles.includes('accountant') || u.email?.includes('accountant')) {
    form.value.department = 'Finance & Accounts'
    form.value.designation = 'Senior Accountant / Head of Finance'
    form.value.budget_vote_head = 'Vote 201 - General Secretariat Administration & Operations'
  } else if (roles.includes('auditor') || u.email?.includes('auditor')) {
    form.value.department = 'Internal Audit Department'
    form.value.designation = 'Head of Internal Audit'
    form.value.budget_vote_head = 'Vote 201 - General Secretariat Administration & Operations'
  } else if (roles.includes('procurement_officer') || u.email?.includes('procurement')) {
    form.value.department = 'Procurement & Disposal Unit (PDU)'
    form.value.designation = 'Head of Procurement & Disposal (PDU)'
  } else if (roles.includes('public_relations') || u.email?.includes('pr')) {
    form.value.department = 'Public Relations & Media'
    form.value.designation = 'PR & Corporate Communications Officer'
  } else if (roles.includes('legal_counsel') || u.email?.includes('legal')) {
    form.value.department = 'Legal & Corporate Affairs'
    form.value.designation = 'Legal Counsel & Compliance Officer'
  } else if (roles.includes('medical_officer') || roles.includes('physiotherapist') || u.email?.includes('medical')) {
    form.value.department = 'Sports Science & Medical Unit'
    form.value.designation = 'Chief Medical Officer / Sports Physician'
    form.value.budget_vote_head = 'Vote 208 - Medical & Anti-Doping Laboratory Supplies'
  } else if (roles.includes('transport_officer') || roles.includes('driver') || u.email?.includes('transport')) {
    form.value.department = 'Transport & Fleet Management'
    form.value.designation = 'Transport Officer & Fleet Manager'
    form.value.budget_vote_head = 'Vote 207 - Transport, Fleet Maintenance & Fuel'
  } else if (roles.includes('human_resources') || roles.includes('hr') || u.email?.includes('hr')) {
    form.value.department = 'Human Resources & Administration'
    form.value.designation = 'Human Resources Manager'
  } else if (roles.includes('technical_department') || roles.includes('ags_technical') || u.email?.includes('technical')) {
    form.value.department = 'Technical Directorate'
    form.value.designation = 'Technical Director'
  } else if (roles.includes('general_secretary') || u.email?.includes('gs')) {
    form.value.department = 'Executive Directorate'
    form.value.designation = 'General Secretary (CEO / Accounting Officer)'
  }
}

async function submitRequisition() {
  if (!declarationAgreed.value) {
    errorMsg.value = 'Please accept the PPDA compliance declaration before submitting.'
    return
  }
  submitting.value = true
  successMsg.value = ''
  errorMsg.value = ''
  try {
    const payload = {
      ...form.value,
      estimated_amount_ugx: grandTotal.value
    }
    const res = await apiPost('/it/ppda', payload)
    successMsg.value = res.data?.message || res.message || 'PPDA Form 5 requisition submitted successfully for HOD endorsement!'
    setTimeout(() => {
      router.push('/it/ppda/status')
    }, 1200)
  } catch (err) {
    errorMsg.value = err.message || 'Could not submit PPDA Form 5 requisition'
    submitting.value = false
  }
}

onMounted(() => {
  detectUserDepartmentAndRole()
})
</script>
