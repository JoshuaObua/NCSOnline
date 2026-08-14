<template>
  <LayoutDefault title="Reception & Visitor Clearance Desk">
    <div class="space-y-6 pb-12 w-full">
      
      <!-- Top Action Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl p-5 shadow-sm transition-colors duration-200 w-full">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-md bg-blue-50 dark:bg-blue-950/60 text-blue-700 dark:text-blue-400 text-xs font-bold uppercase tracking-wider">
              <i class="icofont-id-card text-xs"></i> Front Desk Operations
            </span>
          </div>
          <h2 class="text-xl font-bold text-slate-900 dark:text-white mt-1">Reception & Visitor Clearance Desk</h2>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-0.5">
            Log visitor credentials, dispatch entry clearance to host departments, and issue printable official NCS Gate Passes.
          </p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <button @click="fetchVisitors" class="btn btn-sm btn-light flex items-center gap-1.5" :disabled="loading">
            <i class="icofont-refresh" :class="{ 'animate-spin': loading }"></i> Refresh
          </button>
          <router-link to="/reception/visitors/new" class="btn btn-sm btn-primary flex items-center gap-1.5 shadow-sm">
            <i class="icofont-plus"></i> Initiate Visitor Request
          </router-link>
        </div>
      </div>

      <!-- Reception KPI Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 w-full">
        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Total Visitors Today</h5>
              <h2>{{ visitors.length }}</h2>
              <span class="badge badge-primary"><i class="icofont-calendar"></i> Registered Today</span>
            </div>
            <div class="banner-img bg-primary-light">
              <i class="icofont-users"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Pending Approval</h5>
              <h2>{{ pendingCount }}</h2>
              <span class="badge" :class="pendingCount > 0 ? 'badge-warning' : 'badge-light'">
                <i class="icofont-clock-time"></i> Awaiting Host
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
              <h5>On-Premises</h5>
              <h2>{{ checkedInCount }}</h2>
              <span class="badge badge-success"><i class="icofont-check-circled"></i> Active Clearance</span>
            </div>
            <div class="banner-img bg-success-light">
              <i class="icofont-building"></i>
            </div>
          </div>
        </div>

        <div class="card card-statistic-4 shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0">
          <div class="card-content">
            <div>
              <h5>Departed / Closed</h5>
              <h2>{{ completedCount }}</h2>
              <span class="badge badge-info"><i class="icofont-badge"></i> Missions Concluded</span>
            </div>
            <div class="banner-img bg-cyan-light">
              <i class="icofont-paper-plane"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- Filter & Search Bar -->
      <div class="bg-white dark:bg-slate-900 p-4 rounded-xl border border-slate-200 dark:border-slate-800 w-full space-y-3">
        <div class="flex flex-col md:flex-row items-center justify-between gap-3">
          
          <!-- Search input for quick lookup -->
          <div class="relative w-full md:w-80">
            <i class="icofont-search-2 absolute left-3 top-2.5 text-xs text-slate-400"></i>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search by visitor name, pass #, phone, org..."
              class="form-control text-xs w-full pl-8 pr-3 py-2 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400"
            />
          </div>

          <!-- Dropdowns for Department and Status -->
          <div class="flex flex-wrap items-center gap-3 w-full md:w-auto">
            <div class="w-full sm:w-52">
              <select v-model="filterStatus" class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200">
                <option value="">All Statuses</option>
                <option value="PENDING_APPROVAL">Pending Approval</option>
                <option value="APPROVED">Approved (Pass Ready)</option>
                <option value="CHECKED_IN">Checked In (On-Premises)</option>
                <option value="COMPLETED">Completed / Exited</option>
              </select>
            </div>

            <div class="w-full sm:w-60">
              <select v-model="filterDept" class="form-control text-xs w-full py-2 px-3 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-800 dark:text-slate-200">
                <option value="">All Target Departments</option>
                <option value="Office of the General Secretary">Office of the General Secretary</option>
                <option value="Administration & HR">Administration & HR</option>
                <option value="Technical & Sports Development">Technical & Sports Development</option>
                <option value="Finance & Accounts">Finance & Accounts</option>
                <option value="Engineering & Infrastructure">Engineering & Infrastructure</option>
                <option value="Facilities & Venues Management">Facilities & Venues Management</option>
                <option value="Transport & Fleet Management">Transport & Fleet Management</option>
                <option value="Sports Medicine & Science / WADA">Sports Medicine & Science / WADA</option>
                <option value="Legal & Corporate Affairs">Legal & Corporate Affairs</option>
                <option value="Procurement & Disposal">Procurement & Disposal</option>
                <option value="Internal Audit">Internal Audit</option>
              </select>
            </div>

            <button
              v-if="searchQuery || filterStatus || filterDept"
              @click="resetFilters"
              class="btn btn-sm btn-light text-xs"
              title="Reset all filters"
            >
              <i class="icofont-close"></i> Reset
            </button>
          </div>
        </div>

        <div class="text-xs text-slate-500 font-mono flex items-center justify-between pt-1 border-t border-slate-100 dark:border-slate-800/60">
          <span>Showing {{ filteredVisitors.length }} of {{ visitors.length }} visitor record(s)</span>
          <span v-if="searchQuery" class="text-blue-600 font-semibold">Filtering by: "{{ searchQuery }}"</span>
        </div>
      </div>

      <!-- Visitor Roster Table -->
      <div class="card shadow-sm border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-0 overflow-hidden w-full">
        <div class="card-header border-b border-slate-100 dark:border-slate-800 py-3.5 px-5 flex items-center justify-between">
          <h4 class="text-sm font-bold text-slate-900 dark:text-white">Daily Visitor Clearance Registry</h4>
          <router-link to="/reception/visitors/new" class="btn btn-sm btn-primary flex items-center gap-1.5 text-xs">
            <i class="icofont-plus"></i> New Visitor
          </router-link>
        </div>

        <div class="card-body p-0 overflow-x-auto">
          <table class="table table-striped table-hover mb-0 text-left text-xs">
            <thead class="bg-slate-50 dark:bg-slate-800/60 text-slate-600 dark:text-slate-400 uppercase tracking-wider font-semibold">
              <tr>
                <th class="py-3 px-4">Pass Ref</th>
                <th class="py-3 px-4">Visitor Particulars</th>
                <th class="py-3 px-4">Destination Office</th>
                <th class="py-3 px-4">Host Officer</th>
                <th class="py-3 px-4">Status & Clearance</th>
                <th class="py-3 px-4 text-right">Desk Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
              <tr v-for="v in filteredVisitors" :key="v.id">
                <td class="py-3 px-4">
                  <span class="font-mono font-bold text-blue-600 dark:text-blue-400">{{ v.pass_number }}</span>
                  <div class="text-[10px] text-slate-400 font-mono mt-0.5">{{ formatDate(v.created_at) }}</div>
                </td>
                <td class="py-3 px-4">
                  <div class="font-bold text-slate-900 dark:text-white">{{ v.visitor_name }}</div>
                  <div class="text-[11px] text-slate-500 font-medium">{{ v.visitor_organization || 'Individual Visitor' }}</div>
                  <div class="text-[10px] text-slate-400">{{ v.visitor_phone || '' }} · ID: {{ v.visitor_id_number || 'N/A' }}</div>
                </td>
                <td class="py-3 px-4">
                  <span class="font-semibold text-slate-800 dark:text-slate-200">{{ v.target_department }}</span>
                  <div class="text-[11px] text-slate-500 mt-0.5 truncate max-w-xs">{{ v.purpose_of_visit }}</div>
                </td>
                <td class="py-3 px-4">
                  <div class="font-medium text-slate-800 dark:text-slate-200">{{ v.host_officer_name }}</div>
                  <div v-if="v.vehicle_reg_no" class="text-[10px] text-slate-400 font-mono">Veh: {{ v.vehicle_reg_no }}</div>
                </td>
                <td class="py-3 px-4">
                  <span class="badge" :class="statusBadge(v.status)">
                    {{ formatStatus(v.status) }}
                  </span>
                  <div v-if="v.clearance_code" class="text-[10px] font-mono text-emerald-600 font-bold mt-1">
                    {{ v.clearance_code }}
                  </div>
                </td>
                <td class="py-3 px-4 text-right">
                  <div class="flex items-center justify-end gap-1.5">
                    
                    <!-- Host Approval Action -->
                    <button
                      v-if="v.status === 'PENDING_APPROVAL'"
                      @click="approvePass(v.id)"
                      class="btn btn-xs btn-outline-success flex items-center gap-1"
                      title="Approve Host Clearance"
                    >
                      <i class="icofont-check"></i> Approve
                    </button>

                    <!-- Check-in Action -->
                    <button
                      v-if="v.status === 'APPROVED'"
                      @click="checkInPass(v.id)"
                      class="btn btn-xs btn-primary flex items-center gap-1"
                      title="Issue Badge & Check In"
                    >
                      <i class="icofont-login"></i> Check In
                    </button>

                    <!-- Print Clearance Form -->
                    <button
                      v-if="v.status === 'APPROVED' || v.status === 'CHECKED_IN' || v.status === 'COMPLETED'"
                      @click="openPrintPass(v)"
                      class="btn btn-xs btn-outline-primary flex items-center gap-1"
                      title="Print Official Clearance Pass"
                    >
                      <i class="icofont-printer"></i> Print Pass
                    </button>

                    <!-- Check-out Action -->
                    <button
                      v-if="v.status === 'CHECKED_IN'"
                      @click="checkOutPass(v.id)"
                      class="btn btn-xs btn-outline-danger flex items-center gap-1"
                      title="Record Exit & Check Out"
                    >
                      <i class="icofont-logout"></i> Check Out
                    </button>

                  </div>
                </td>
              </tr>
              <tr v-if="!filteredVisitors.length">
                <td colspan="6" class="text-center py-10 text-xs text-slate-400">
                  No visitor records matching current search or filters.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- PRINTABLE VISITOR CLEARANCE FORM / PASS MODAL -->
      <div v-if="selectedPrintPass" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/70 backdrop-blur-sm overflow-y-auto">
        <div class="bg-white dark:bg-slate-900 rounded-2xl max-w-xl w-full p-6 shadow-2xl border border-slate-200 dark:border-slate-800 space-y-4">
          
          <div class="flex items-center justify-between no-print">
            <h3 class="text-sm font-bold text-slate-900 dark:text-white">Official Clearance Form Preview</h3>
            <div class="flex items-center gap-2">
              <button @click="triggerPrint" class="btn btn-sm btn-primary flex items-center gap-1">
                <i class="icofont-printer"></i> Print Clearance Pass
              </button>
              <button @click="selectedPrintPass = null" class="btn btn-sm btn-light">
                <i class="icofont-close"></i> Close
              </button>
            </div>
          </div>

          <!-- Printable Pass Document Card -->
          <div id="ncs-clearance-pass" class="border-2 border-slate-900 dark:border-slate-100 p-6 rounded-xl bg-white text-slate-900 space-y-4">
            
            <!-- Official Header -->
            <div class="text-center border-b-2 border-slate-900 pb-3 space-y-1">
              <div class="text-xs uppercase tracking-widest font-extrabold text-slate-700">Republic of Uganda</div>
              <h1 class="text-base font-black uppercase tracking-wider text-slate-900">National Council of Sports</h1>
              <p class="text-[10px] text-slate-600">Plot 2-10 Coronation Avenue, Lugogo Sports Complex · P.O. Box 20077 Kampala</p>
              <div class="inline-block mt-1 px-3 py-0.5 bg-slate-900 text-white font-mono text-xs font-bold uppercase rounded">
                Official Visitor Clearance & Gate Pass
              </div>
            </div>

            <!-- Pass Core Info Grid -->
            <div class="grid grid-cols-2 gap-3 text-xs">
              <div class="border border-slate-300 p-2.5 rounded bg-slate-50">
                <span class="text-[10px] uppercase font-bold text-slate-500 block">Pass Reference</span>
                <strong class="font-mono text-sm text-slate-900">{{ selectedPrintPass.pass_number }}</strong>
              </div>
              <div class="border border-slate-300 p-2.5 rounded bg-slate-50">
                <span class="text-[10px] uppercase font-bold text-slate-500 block">Security Clearance Code</span>
                <strong class="font-mono text-sm text-blue-700 font-bold">{{ selectedPrintPass.clearance_code || 'CLR-OFFICIAL' }}</strong>
              </div>
            </div>

            <!-- Visitor & Host Details -->
            <div class="border border-slate-300 p-3 rounded text-xs space-y-2">
              <div class="grid grid-cols-2 gap-2">
                <div>
                  <span class="text-[10px] text-slate-500 uppercase block">Visitor Name:</span>
                  <strong class="text-slate-900">{{ selectedPrintPass.visitor_name }}</strong>
                </div>
                <div>
                  <span class="text-[10px] text-slate-500 uppercase block">ID / NIN / Passport:</span>
                  <span class="font-mono text-slate-800">{{ selectedPrintPass.visitor_id_number || 'N/A' }}</span>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-2 pt-2 border-t border-slate-200">
                <div>
                  <span class="text-[10px] text-slate-500 uppercase block">Host Department:</span>
                  <strong class="text-slate-900">{{ selectedPrintPass.target_department }}</strong>
                </div>
                <div>
                  <span class="text-[10px] text-slate-500 uppercase block">Host Officer:</span>
                  <strong class="text-slate-900">{{ selectedPrintPass.host_officer_name }}</strong>
                </div>
              </div>

              <div class="pt-2 border-t border-slate-200">
                <span class="text-[10px] text-slate-500 uppercase block">Purpose of Visit:</span>
                <p class="text-slate-800 font-medium">{{ selectedPrintPass.purpose_of_visit }}</p>
              </div>

              <div v-if="selectedPrintPass.vehicle_reg_no || selectedPrintPass.items_declared" class="grid grid-cols-2 gap-2 pt-2 border-t border-slate-200">
                <div v-if="selectedPrintPass.vehicle_reg_no">
                  <span class="text-[10px] text-slate-500 uppercase block">Vehicle Reg:</span>
                  <span class="font-mono text-slate-900 font-bold">{{ selectedPrintPass.vehicle_reg_no }}</span>
                </div>
                <div v-if="selectedPrintPass.items_declared">
                  <span class="text-[10px] text-slate-500 uppercase block">Declared Property:</span>
                  <span class="text-slate-800">{{ selectedPrintPass.items_declared }}</span>
                </div>
              </div>
            </div>

            <!-- Signoff & Departure Protocol -->
            <div class="pt-2 border-t-2 border-slate-900 grid grid-cols-3 gap-2 text-[10px] text-center">
              <div>
                <span class="block text-slate-400">Receptionist Stamp</span>
                <div class="h-10 border-b border-dashed border-slate-400 mt-1 flex items-center justify-center font-mono text-[9px] text-slate-500">
                  NCS FRONT DESK
                </div>
              </div>
              <div>
                <span class="block text-slate-400">Host Officer Signature</span>
                <div class="h-10 border-b border-dashed border-slate-400 mt-1 flex items-center justify-center font-mono text-[9px] text-slate-500">
                  APPROVED
                </div>
              </div>
              <div>
                <span class="block text-slate-400">Security Gate Exit</span>
                <div class="h-10 border-b border-dashed border-slate-400 mt-1 flex items-center justify-center font-mono text-[9px] text-slate-500">
                  CHECKED OUT
                </div>
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
import { apiGet, apiPut } from '@/api/client'

const loading = ref(false)
const visitors = ref([])
const selectedPrintPass = ref(null)

const searchQuery = ref('')
const filterStatus = ref('')
const filterDept = ref('')

const pendingCount = computed(() => visitors.value.filter(v => v.status === 'PENDING_APPROVAL').length)
const checkedInCount = computed(() => visitors.value.filter(v => v.status === 'CHECKED_IN').length)
const completedCount = computed(() => visitors.value.filter(v => v.status === 'COMPLETED').length)

const filteredVisitors = computed(() => {
  return visitors.value.filter(v => {
    if (filterStatus.value && v.status !== filterStatus.value) return false
    if (filterDept.value && v.target_department !== filterDept.value) return false
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase().trim()
      const matchName = (v.visitor_name || '').toLowerCase().includes(q)
      const matchPass = (v.pass_number || '').toLowerCase().includes(q)
      const matchPhone = (v.visitor_phone || '').toLowerCase().includes(q)
      const matchOrg = (v.visitor_organization || '').toLowerCase().includes(q)
      const matchHost = (v.host_officer_name || '').toLowerCase().includes(q)
      const matchDept = (v.target_department || '').toLowerCase().includes(q)
      if (!matchName && !matchPass && !matchPhone && !matchOrg && !matchHost && !matchDept) return false
    }
    return true
  })
})

function resetFilters() {
  searchQuery.value = ''
  filterStatus.value = ''
  filterDept.value = ''
}

function formatStatus(status) {
  switch (status) {
    case 'PENDING_APPROVAL': return 'Pending Approval'
    case 'APPROVED': return 'Approved'
    case 'CHECKED_IN': return 'Checked In'
    case 'COMPLETED': return 'Completed'
    default: return status
  }
}

function statusBadge(status) {
  switch (status) {
    case 'APPROVED': return 'badge-success'
    case 'CHECKED_IN': return 'badge-info'
    case 'PENDING_APPROVAL': return 'badge-warning'
    case 'COMPLETED': return 'badge-light'
    default: return 'badge-primary'
  }
}

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

async function fetchVisitors() {
  loading.value = true
  try {
    const res = await apiGet('/reception/visitors')
    visitors.value = res.data?.visitors || res.visitors || []
  } catch {
    visitors.value = []
  } finally {
    loading.value = false
  }
}

async function approvePass(id) {
  try {
    await apiPut(`/reception/visitors/${id}/approve`, {})
    await fetchVisitors()
  } catch (err) {
    alert(err.message || 'Could not approve pass')
  }
}

async function checkInPass(id) {
  try {
    await apiPut(`/reception/visitors/${id}/checkin`, {})
    await fetchVisitors()
  } catch (err) {
    alert(err.message || 'Could not check in visitor')
  }
}

async function checkOutPass(id) {
  try {
    await apiPut(`/reception/visitors/${id}/checkout`, {})
    await fetchVisitors()
  } catch (err) {
    alert(err.message || 'Could not check out visitor')
  }
}

function openPrintPass(v) {
  selectedPrintPass.value = v
}

function triggerPrint() {
  window.print()
}

onMounted(() => {
  fetchVisitors()
})
</script>

<style scoped>
@media print {
  body * {
    visibility: hidden;
  }
  #ncs-clearance-pass, #ncs-clearance-pass * {
    visibility: visible;
  }
  #ncs-clearance-pass {
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
