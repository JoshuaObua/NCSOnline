<template>
  <LayoutDefault title="PPDA Form 5 Requisitions">
    <div class="space-y-6 pb-12 w-full">
      
      <!-- Procurement KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 w-full">
        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Total Requisitions</h5>
              <h2>{{ requisitions.length }}</h2>
              <span class="badge badge-primary"><i class="icofont-files-stack"></i> Raised Dossiers</span>
            </div>
            <div class="banner-img bg-primary-light">
              <i class="icofont-law-order"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Pending Clearances</h5>
              <h2>{{ pendingCount }}</h2>
              <span class="badge" :class="pendingCount > 0 ? 'badge-warning' : 'badge-light'">
                <i class="icofont-clock-time"></i> In Vetting Pipeline
              </span>
            </div>
            <div class="banner-img bg-warning-light">
              <i class="icofont-sand-clock"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Approved & Contracted</h5>
              <h2>{{ approvedCount }}</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Authorized</span>
            </div>
            <div class="banner-img bg-success-light">
              <i class="icofont-certificate"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Total Value (UGX)</h5>
              <h2 class="text-lg">UGX {{ totalRequisitionValue.toLocaleString() }}</h2>
              <span class="badge badge-info"><i class="icofont-coins"></i> Vote Allocated</span>
            </div>
            <div class="banner-img bg-cyan-light">
              <i class="icofont-coins"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- Search & Filter Bar -->
      <div class="bg-white dark:bg-slate-900 p-4 rounded-xl border border-slate-200 dark:border-slate-800 w-full space-y-3">
        <div class="flex flex-col md:flex-row items-center justify-between gap-3">
          <div class="relative w-full md:w-80">
            <i class="icofont-search-2 absolute left-3 top-2.5 text-xs text-slate-400"></i>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search by subject, ref #, vote head..."
              class="form-control text-xs w-full pl-8 pr-3 py-2 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
            />
          </div>

          <div class="flex items-center gap-3 w-full md:w-auto">
            <select
              v-model="filterStatus"
              class="form-control text-xs w-full sm:w-56 py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200"
            >
              <option value="">All Requisition Statuses</option>
              <option value="SUBMITTED">Submitted (Awaiting HOD)</option>
              <option value="HOD_APPROVED">HOD Endorsed (Awaiting AGS-T)</option>
              <option value="AGST_APPROVED">AGS-T Recommended (Awaiting Finance)</option>
              <option value="FINANCE_CLEARED">Finance Cleared (Awaiting General Secretary)</option>
              <option value="ACCOUNTING_OFFICER_APPROVED">Approved by Accounting Officer (General Secretary)</option>
              <option value="PO_ISSUED">Purchase Order Issued</option>
              <option value="DELIVERED">Delivered & Verified</option>
              <option value="REJECTED">Rejected / Returned for Revisions</option>
            </select>

            <button
              v-if="searchQuery || filterStatus"
              @click="searchQuery = ''; filterStatus = ''"
              class="btn btn-sm btn-light text-xs"
              title="Reset filter"
            >
              <i class="icofont-close"></i> Reset
            </button>
          </div>
        </div>

        <div class="text-xs text-slate-500 font-mono flex items-center justify-between pt-1 border-t border-slate-100 dark:border-slate-800/60">
          <span>Showing {{ filteredRequisitions.length }} of {{ requisitions.length }} requisition(s)</span>
        </div>
      </div>

      <!-- Requisitions Table -->
      <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden w-full">
        <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5 flex items-center justify-between">
          <h4 class="text-sm font-bold text-slate-900 dark:text-white">Statutory PPDA Form 5 Register</h4>
          <router-link to="/it/ppda/new" class="btn btn-sm btn-primary flex items-center gap-1.5 text-xs">
            <i class="icofont-plus"></i> New Requisition
          </router-link>
        </div>

        <div class="card-body p-0 overflow-x-auto">
          <table class="table table-striped table-hover mb-0 text-left text-xs">
            <thead class="bg-slate-50 dark:bg-slate-800/60 text-slate-600 dark:text-slate-400 uppercase tracking-wider font-semibold">
              <tr>
                <th class="py-3 px-4">Ref No.</th>
                <th class="py-3 px-4">Subject of Procurement</th>
                <th class="py-3 px-4">Category & Vote Line</th>
                <th class="py-3 px-4 text-right">Est. Amount (UGX)</th>
                <th class="py-3 px-4">Vetting Stage & Workflow</th>
                <th class="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
              <tr v-for="req in filteredRequisitions" :key="req.id">
                <td class="py-3 px-4">
                  <span class="font-mono font-bold text-blue-600 dark:text-blue-400">{{ req.reference_no }}</span>
                  <div class="text-[10px] text-slate-400 mt-0.5">{{ formatDate(req.created_at) }}</div>
                </td>
                <td class="py-3 px-4 max-w-sm">
                  <div class="font-bold text-slate-900 dark:text-white leading-snug">{{ req.subject_of_procurement }}</div>
                  <div class="text-[11px] text-slate-500 mt-0.5 truncate">{{ req.officer_name }} · {{ req.department }}</div>
                </td>
                <td class="py-3 px-4">
                  <span class="badge badge-light font-semibold mb-1">{{ req.procurement_category }}</span>
                  <div class="text-[11px] text-slate-600 dark:text-slate-400 font-mono">{{ req.budget_vote_head }}</div>
                </td>
                <td class="py-3 px-4 text-right font-mono font-bold text-slate-900 dark:text-white">
                  UGX {{ Number(req.estimated_amount_ugx || 0).toLocaleString() }}
                </td>
                <td class="py-3 px-4">
                  <span class="badge" :class="statusBadge(req.status)">
                    {{ formatStatus(req.status) }}
                  </span>
                  <div class="text-[10px] text-slate-500 font-medium mt-1 flex items-center gap-1">
                    <i class="icofont-arrow-right text-blue-500"></i> {{ req.current_stage || 'In Review' }}
                  </div>
                  <div v-if="req.po_number" class="text-[10px] font-mono text-emerald-600 font-bold mt-0.5">
                    PO: {{ req.po_number }}
                  </div>
                </td>
                <td class="py-3 px-4 text-right">
                  <button
                    @click="viewRequisitionDetails(req)"
                    class="btn btn-xs btn-outline-primary flex items-center gap-1 ml-auto"
                    title="View Full PPDA Form 5 Details"
                  >
                    <i class="icofont-eye"></i> View & Process
                  </button>
                </td>
              </tr>
              <tr v-if="!filteredRequisitions.length">
                <td colspan="6" class="text-center py-10 text-xs text-slate-400">
                  No PPDA Form 5 requisitions matching current filters.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- DETAILED PPDA FORM 5 VIEW / PRINT / WORKFLOW MODAL -->
      <div v-if="selectedReq" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/70 backdrop-blur-sm overflow-y-auto">
        <div class="bg-white dark:bg-slate-900 rounded-2xl max-w-4xl w-full p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-4 max-h-[92vh] overflow-y-auto">
          
          <!-- Top Modal Action Bar -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 no-print border-b border-slate-100 dark:border-slate-800 pb-3">
            <div class="flex items-center gap-2">
              <span class="badge badge-primary font-mono text-xs">{{ selectedReq.reference_no }}</span>
              <span class="text-xs font-bold text-slate-700 dark:text-slate-300">PPDA Form 5 Requisition Dossier</span>
            </div>
            <div class="flex items-center gap-2">
              <button @click="triggerPrint" class="btn btn-sm btn-primary flex items-center gap-1 text-xs">
                <i class="icofont-printer"></i> Print Official PPDA Form 5
              </button>
              <button @click="selectedReq = null" class="btn btn-sm btn-light text-xs">
                <i class="icofont-close"></i> Close
              </button>
            </div>
          </div>

          <!-- Workflow Transition Action Bar -->
          <div class="no-print p-3.5 rounded-xl bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 space-y-2 text-xs">
            <div class="flex items-center justify-between">
              <span class="font-bold text-slate-800 dark:text-slate-200 flex items-center gap-1.5">
                <i class="icofont-law-order text-blue-600"></i> Departmental Approval & Forwarding Workflow:
              </span>
              <span class="badge" :class="statusBadge(selectedReq.status)">{{ formatStatus(selectedReq.status) }}</span>
            </div>

            <!-- Workflow Buttons based on Stage -->
            <div class="flex flex-wrap items-center gap-2 pt-1">
              <!-- Stage 1 -> Stage 2 -->
              <button
                v-if="selectedReq.status === 'SUBMITTED'"
                @click="advanceWorkflow('ENDORSE_HOD', 'Endorsed by Head of Department')"
                :disabled="advancing"
                class="btn btn-xs btn-primary flex items-center gap-1"
              >
                <i class="icofont-check-circled"></i> Endorse as Head of Department (HOD) → Forward to AGS-Technical
              </button>

              <!-- Stage 2 -> Stage 3 -->
              <button
                v-if="selectedReq.status === 'HOD_APPROVED'"
                @click="advanceWorkflow('ENDORSE_AGST', 'Recommended by Assistant General Secretary - Technical')"
                :disabled="advancing"
                class="btn btn-xs btn-primary flex items-center gap-1"
              >
                <i class="icofont-check-circled"></i> Recommend as AGS - Technical → Forward to Head of Finance
              </button>

              <!-- Stage 3 -> Stage 4 -->
              <button
                v-if="selectedReq.status === 'AGST_APPROVED'"
                @click="advanceWorkflow('COMMIT_FINANCE', 'Budget Vote Head committed and cleared')"
                :disabled="advancing"
                class="btn btn-xs btn-success flex items-center gap-1"
              >
                <i class="icofont-coins"></i> Commit Budget Vote (Head of Finance) → Forward to General Secretary
              </button>

              <!-- Stage 4 -> Stage 5 -->
              <button
                v-if="selectedReq.status === 'FINANCE_CLEARED'"
                @click="advanceWorkflow('APPROVE_ACCOUNTING_OFFICER', 'Statutory procurement authorized by General Secretary / Accounting Officer')"
                :disabled="advancing"
                class="btn btn-xs btn-success flex items-center gap-1"
              >
                <i class="icofont-certificate"></i> Authorize as General Secretary (CEO / Accounting Officer)
              </button>

              <!-- Stage 5 -> Stage 6 -->
              <button
                v-if="selectedReq.status === 'ACCOUNTING_OFFICER_APPROVED'"
                @click="advanceWorkflow('ISSUE_PO', 'Purchase Order issued by PDU', 'PO-NCS-2026-092')"
                :disabled="advancing"
                class="btn btn-xs btn-info flex items-center gap-1"
              >
                <i class="icofont-file-text"></i> Issue Official Purchase Order (PDU)
              </button>

              <!-- Stage 6 -> Stage 7 -->
              <button
                v-if="selectedReq.status === 'PO_ISSUED'"
                @click="advanceWorkflow('CONFIRM_DELIVERY', 'Goods inspected, verified and received by User Department')"
                :disabled="advancing"
                class="btn btn-xs btn-success flex items-center gap-1"
              >
                <i class="icofont-box"></i> Confirm Inspection & Delivery (User Department)
              </button>

              <!-- Reject / Return -->
              <button
                v-if="!['DELIVERED', 'REJECTED'].includes(selectedReq.status)"
                @click="advanceWorkflow('REJECT', 'Returned for budget adjustment / technical revision')"
                :disabled="advancing"
                class="btn btn-xs btn-outline-danger flex items-center gap-1 ml-auto"
              >
                <i class="icofont-close-circled"></i> Return / Reject with Queries
              </button>
            </div>
            
            <div v-if="workflowSuccessMsg" class="p-2 rounded bg-emerald-100 dark:bg-emerald-950/60 text-emerald-800 dark:text-emerald-300 text-[11px] font-semibold mt-1">
              <i class="icofont-check-circled"></i> {{ workflowSuccessMsg }}
            </div>
          </div>

          <!-- Printable Document Paper -->
          <div id="ppda-form-print" class="border border-slate-300 dark:border-slate-700 p-6 rounded-xl bg-white text-slate-900 space-y-5 text-xs">
            
            <!-- Official Heading -->
            <div class="text-center border-b-2 border-slate-900 pb-3 space-y-1">
              <div class="text-[11px] uppercase tracking-widest font-black text-slate-800">Republic of Uganda</div>
              <h1 class="text-base font-black uppercase tracking-wider text-slate-900">National Council of Sports</h1>
              <div class="text-[10px] text-slate-600">Public Procurement and Disposal of Public Assets (PPDA) Act, 2003 / 2023 Regulations</div>
              <div class="inline-block mt-1 px-3 py-0.5 bg-slate-900 text-white font-mono text-xs font-bold uppercase rounded">
                PPDA FORM 5: Request for Approval of Procurement
              </div>
            </div>

            <!-- Meta Details Grid -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 bg-slate-50 p-3 rounded-lg border border-slate-200">
              <div>
                <span class="text-[10px] text-slate-500 font-bold uppercase block">Procurement Ref</span>
                <strong class="font-mono text-xs">{{ selectedReq.reference_no }}</strong>
              </div>
              <div>
                <span class="text-[10px] text-slate-500 font-bold uppercase block">Date Submitted</span>
                <strong>{{ formatDate(selectedReq.created_at) }}</strong>
              </div>
              <div>
                <span class="text-[10px] text-slate-500 font-bold uppercase block">Category</span>
                <span class="badge badge-light">{{ selectedReq.procurement_category }}</span>
              </div>
              <div>
                <span class="text-[10px] text-slate-500 font-bold uppercase block">Approval Status</span>
                <span class="badge" :class="statusBadge(selectedReq.status)">{{ formatStatus(selectedReq.status) }}</span>
              </div>
            </div>

            <!-- User Department & Financial Particulars Section -->
            <div class="space-y-1">
              <h5 class="font-bold text-slate-900 uppercase text-[11px] border-b border-slate-200 pb-1">1. User Department & Requisitioning Officer Particulars</h5>
              <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 pt-1">
                <div><span class="text-slate-500 text-[10px] block">Requisitioning Officer:</span><strong>{{ selectedReq.officer_name }}</strong></div>
                <div><span class="text-slate-500 text-[10px] block">Official Designation:</span><span>{{ selectedReq.designation }}</span></div>
                <div><span class="text-slate-500 text-[10px] block">User Department:</span><span>{{ selectedReq.department }}</span></div>
                <div><span class="text-slate-500 text-[10px] block">Financial Year / Phone:</span><span class="font-mono">{{ selectedReq.financial_year || 'FY 2026/2027' }} · {{ selectedReq.contact_phone || 'N/A' }}</span></div>
              </div>
            </div>

            <!-- Subject, Method & Vote Section -->
            <div class="space-y-1">
              <h5 class="font-bold text-slate-900 uppercase text-[11px] border-b border-slate-200 pb-1">2. Procurement Subject, Method & Funding Allocation</h5>
              <div class="pt-1 space-y-1.5">
                <div><span class="text-slate-500 text-[10px] block">Subject / Title of Procurement:</span><strong class="text-slate-900">{{ selectedReq.subject_of_procurement }}</strong></div>
                <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 pt-1">
                  <div><span class="text-slate-500 text-[10px] block">Procurement Method:</span><span class="font-semibold">{{ selectedReq.procurement_method || 'Request for Quotations (RFQ)' }}</span></div>
                  <div><span class="text-slate-500 text-[10px] block">Budget Vote Head / Line:</span><span class="font-mono font-semibold text-[11px]">{{ selectedReq.budget_vote_head }}</span></div>
                  <div><span class="text-slate-500 text-[10px] block">Source of Funds:</span><span>{{ selectedReq.source_of_funds || 'GoU Statutory Subvention' }}</span></div>
                  <div><span class="text-slate-500 text-[10px] block">Estimated Total Cost:</span><strong class="font-mono text-blue-700 text-xs">UGX {{ Number(selectedReq.estimated_amount_ugx || 0).toLocaleString() }}</strong></div>
                </div>
                <div class="grid grid-cols-2 gap-2 pt-1 bg-slate-50 p-2 rounded border border-slate-200">
                  <div><span class="text-slate-500 text-[10px] block">Delivery Location:</span><span class="font-medium">{{ selectedReq.delivery_location || 'NCS Headquarters Lugogo' }}</span></div>
                  <div><span class="text-slate-500 text-[10px] block">Warranty / SLA Terms:</span><span class="font-medium">{{ selectedReq.warranty_requirement || '1 Year Comprehensive Warranty' }}</span></div>
                </div>
              </div>
            </div>

            <!-- Line Items Table -->
            <div class="space-y-1">
              <h5 class="font-bold text-slate-900 uppercase text-[11px] border-b border-slate-200 pb-1">3. Schedule of Requirements & Line Items</h5>
              <table class="table table-sm text-xs border border-slate-200 w-full mt-2">
                <thead class="bg-slate-100">
                  <tr>
                    <th class="p-2">#</th>
                    <th class="p-2">Item Description</th>
                    <th class="p-2">Qty</th>
                    <th class="p-2">Unit</th>
                    <th class="p-2 text-right">Unit Price (UGX)</th>
                    <th class="p-2 text-right">Total Price (UGX)</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(item, i) in selectedReq.items || []" :key="i">
                    <td class="p-2 font-mono">{{ i + 1 }}</td>
                    <td class="p-2 font-medium">{{ item.item_name }}</td>
                    <td class="p-2 font-mono">{{ item.quantity }}</td>
                    <td class="p-2">{{ item.unit }}</td>
                    <td class="p-2 text-right font-mono">{{ Number(item.unit_price_ugx || 0).toLocaleString() }}</td>
                    <td class="p-2 text-right font-mono font-bold">{{ Number(item.total_price_ugx || 0).toLocaleString() }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Technical Specifications & Justification -->
            <div class="space-y-2">
              <div>
                <h5 class="font-bold text-slate-900 uppercase text-[11px] border-b border-slate-200 pb-1">4. Technical Specifications / Terms of Reference</h5>
                <p class="pt-1 text-slate-800 leading-relaxed">{{ selectedReq.technical_specifications }}</p>
              </div>
              <div>
                <h5 class="font-bold text-slate-900 uppercase text-[11px] border-b border-slate-200 pb-1">5. Justification for Procurement</h5>
                <p class="pt-1 text-slate-800 leading-relaxed">{{ selectedReq.justification }}</p>
              </div>
            </div>

            <!-- 5-Tier Statutory Approval & Forwarding Sign-off Grid -->
            <div class="pt-3 border-t-2 border-slate-900 grid grid-cols-5 gap-2 text-[9px] text-center">
              <div>
                <span class="block text-slate-600 font-bold">1. Requisitioning Officer</span>
                <div class="h-10 border-b border-dashed border-slate-400 mt-1 flex items-center justify-center font-mono text-[9px] text-slate-700">
                  {{ selectedReq.officer_name }}
                </div>
                <span class="text-[8px] text-slate-400">Originated</span>
              </div>
              <div>
                <span class="block text-slate-600 font-bold">2. Head of Department</span>
                <div class="h-10 border-b border-dashed border-slate-400 mt-1 flex items-center justify-center font-mono text-[9px] text-emerald-700 font-bold">
                  {{ ['HOD_APPROVED', 'AGST_APPROVED', 'FINANCE_CLEARED', 'ACCOUNTING_OFFICER_APPROVED', 'PO_ISSUED', 'DELIVERED'].includes(selectedReq.status) ? 'ENDORSED' : 'PENDING' }}
                </div>
                <span class="text-[8px] text-slate-400">HOD Endorsement</span>
              </div>
              <div>
                <span class="block text-slate-600 font-bold">3. AGS - Technical</span>
                <div class="h-10 border-b border-dashed border-slate-400 mt-1 flex items-center justify-center font-mono text-[9px] text-emerald-700 font-bold">
                  {{ ['AGST_APPROVED', 'FINANCE_CLEARED', 'ACCOUNTING_OFFICER_APPROVED', 'PO_ISSUED', 'DELIVERED'].includes(selectedReq.status) ? 'RECOMMENDED' : 'PENDING' }}
                </div>
                <span class="text-[8px] text-slate-400">Technical Directorate</span>
              </div>
              <div>
                <span class="block text-slate-600 font-bold">4. Head of Finance</span>
                <div class="h-10 border-b border-dashed border-slate-400 mt-1 flex items-center justify-center font-mono text-[9px] text-emerald-700 font-bold">
                  {{ ['FINANCE_CLEARED', 'ACCOUNTING_OFFICER_APPROVED', 'PO_ISSUED', 'DELIVERED'].includes(selectedReq.status) ? 'VOTE COMMITTED' : 'PENDING' }}
                </div>
                <span class="text-[8px] text-slate-400">Vote Head Cleared</span>
              </div>
              <div>
                <span class="block text-slate-600 font-bold">5. Accounting Officer</span>
                <div class="h-10 border-b border-dashed border-slate-400 mt-1 flex items-center justify-center font-mono text-[9px] text-blue-700 font-bold">
                  {{ ['ACCOUNTING_OFFICER_APPROVED', 'PO_ISSUED', 'DELIVERED'].includes(selectedReq.status) ? 'AUTHORIZED' : 'IN PROGRESS' }}
                </div>
                <span class="text-[8px] text-slate-400">General Secretary</span>
              </div>
            </div>

          </div>
        </div>
      </div>

    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { apiGet, apiPost } from '@/api/client'

const loading = ref(false)
const advancing = ref(false)
const workflowSuccessMsg = ref('')
const requisitions = ref([])
const selectedReq = ref(null)

const searchQuery = ref('')
const filterStatus = ref('')

const pendingCount = computed(() => {
  return requisitions.value.filter(r => ['SUBMITTED', 'HOD_APPROVED', 'AGST_APPROVED', 'PDU_REVIEW', 'FINANCE_CLEARED'].includes(r.status)).length
})

const approvedCount = computed(() => {
  return requisitions.value.filter(r => ['ACCOUNTING_OFFICER_APPROVED', 'PO_ISSUED', 'DELIVERED'].includes(r.status)).length
})

const totalRequisitionValue = computed(() => {
  return requisitions.value.reduce((sum, r) => sum + (Number(r.estimated_amount_ugx) || 0), 0)
})

const filteredRequisitions = computed(() => {
  return requisitions.value.filter(r => {
    if (filterStatus.value && r.status !== filterStatus.value) return false
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase().trim()
      const matchSub = (r.subject_of_procurement || '').toLowerCase().includes(q)
      const matchRef = (r.reference_no || '').toLowerCase().includes(q)
      const matchVote = (r.budget_vote_head || '').toLowerCase().includes(q)
      const matchPO = (r.po_number || '').toLowerCase().includes(q)
      const matchOff = (r.officer_name || '').toLowerCase().includes(q)
      if (!matchSub && !matchRef && !matchVote && !matchPO && !matchOff) return false
    }
    return true
  })
})

function formatStatus(status) {
  switch (status) {
    case 'SUBMITTED': return 'Submitted (Awaiting HOD)'
    case 'HOD_APPROVED': return 'HOD Endorsed'
    case 'AGST_APPROVED': return 'AGS-T Recommended'
    case 'PDU_REVIEW': return 'PDU Vetting'
    case 'FINANCE_CLEARED': return 'Finance Vote Cleared'
    case 'ACCOUNTING_OFFICER_APPROVED': return 'Approved by Accounting Officer'
    case 'PO_ISSUED': return 'Purchase Order Issued'
    case 'DELIVERED': return 'Delivered & Verified'
    case 'REJECTED': return 'Rejected / Revisions Required'
    default: return status || 'Submitted'
  }
}

function statusBadge(status) {
  switch (status) {
    case 'ACCOUNTING_OFFICER_APPROVED':
    case 'PO_ISSUED':
    case 'DELIVERED':
      return 'badge-success'
    case 'AGST_APPROVED':
    case 'PDU_REVIEW':
    case 'HOD_APPROVED':
    case 'FINANCE_CLEARED':
    case 'SUBMITTED':
      return 'badge-warning'
    case 'REJECTED':
      return 'badge-danger'
    default:
      return 'badge-primary'
  }
}

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric' })
}

async function fetchRequisitions() {
  loading.value = true
  try {
    const res = await apiGet('/it/ppda')
    requisitions.value = res.data?.requisitions || res.requisitions || []
    if (selectedReq.value) {
      const updated = requisitions.value.find(r => r.id === selectedReq.value.id || r.reference_no === selectedReq.value.reference_no)
      if (updated) selectedReq.value = updated
    }
  } catch {
    requisitions.value = []
  } finally {
    loading.value = false
  }
}

async function advanceWorkflow(action, remarks, poNumber = '') {
  if (!selectedReq.value) return
  advancing.value = true
  workflowSuccessMsg.value = ''
  try {
    const res = await apiPost(`/it/ppda/${selectedReq.value.id}/action`, {
      action,
      remarks,
      po_number: poNumber
    })
    workflowSuccessMsg.value = res.data?.message || res.message || 'Workflow advanced successfully'
    selectedReq.value.status = res.data?.status || res.status
    selectedReq.value.current_stage = res.data?.stage || res.stage
    if (poNumber) selectedReq.value.po_number = poNumber
    await fetchRequisitions()
  } catch (err) {
    alert(err.message || 'Could not transition workflow')
  } finally {
    advancing.value = false
  }
}

function viewRequisitionDetails(req) {
  selectedReq.value = req
  workflowSuccessMsg.value = ''
}

function triggerPrint() {
  window.print()
}

onMounted(() => {
  fetchRequisitions()
})
</script>

<style scoped>
@media print {
  body * {
    visibility: hidden;
  }
  #ppda-form-print, #ppda-form-print * {
    visibility: visible;
  }
  #ppda-form-print {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
  }
  .no-print {
    display: none !important;
  }
}
</style>
