<template>
  <div class="card p-5">
    
    <!-- Top Report Header & Type Selector -->
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center mb-4 gap-3">
      <div>
        <h3 class="text-base font-bold text-slate-800 dark:text-slate-100 mb-1">Institutional Reports & Visitor Logs</h3>
        <p class="text-xs text-slate-400 mb-0">Generate, sort, filter, and export comprehensive administrative logs and clearances.</p>
      </div>

      <!-- Report Category Selector -->
      <div class="btn-group" role="group">
        <button 
          type="button" 
          class="btn btn-sm" 
          :class="reportType === 'visitor-clearance' ? 'btn-primary' : 'btn-outline-primary'"
          @click="reportType = 'visitor-clearance'"
        >
          <i class="icofont-badge mr-1"></i> Visitor Clearances
        </button>
        <button 
          type="button" 
          class="btn btn-sm" 
          :class="reportType === 'visitor-logs' ? 'btn-primary' : 'btn-outline-primary'"
          @click="reportType = 'visitor-logs'"
        >
          <i class="icofont-id-card mr-1"></i> Gate Access Logs
        </button>
        <button 
          type="button" 
          class="btn btn-sm" 
          :class="reportType === 'leave-requests' ? 'btn-primary' : 'btn-outline-primary'"
          @click="reportType = 'leave-requests'"
        >
          <i class="icofont-calendar mr-1"></i> Leave Applications
        </button>
      </div>
    </div>

    <!-- Summary Stats Bar -->
    <div class="grid grid-cols-1 sm:grid-cols-4 gap-3 mb-4">
      <div class="p-3 rounded border border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/50">
        <span class="text-[10px] uppercase font-bold text-slate-400 block">Total Filtered Entries</span>
        <strong class="text-lg text-slate-800 dark:text-slate-100">{{ sortedData.length }}</strong>
      </div>
      <div class="p-3 rounded border border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/50">
        <span class="text-[10px] uppercase font-bold text-emerald-500 block">Approved / Clear</span>
        <strong class="text-lg text-emerald-600 dark:text-emerald-400">{{ countStatus('Approved') }}</strong>
      </div>
      <div class="p-3 rounded border border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/50">
        <span class="text-[10px] uppercase font-bold text-amber-500 block">Pending Review</span>
        <strong class="text-lg text-amber-600 dark:text-amber-400">{{ countStatus('Pending') }}</strong>
      </div>
      <div class="p-3 rounded border border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/50">
        <span class="text-[10px] uppercase font-bold text-blue-500 block">Active On-Site</span>
        <strong class="text-lg text-blue-600 dark:text-blue-400">{{ countStatus('Checked In') }}</strong>
      </div>
    </div>

    <!-- Filters & Range Bar -->
    <div class="p-4 rounded border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 mb-4 space-y-3">
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-3">
        
        <!-- Start Date -->
        <div>
          <label class="block text-[10px] font-bold text-slate-500 uppercase mb-1">Start Date</label>
          <input v-model="startDate" type="date" class="form-control form-control-sm text-xs" />
        </div>

        <!-- End Date -->
        <div>
          <label class="block text-[10px] font-bold text-slate-500 uppercase mb-1">End Date</label>
          <input v-model="endDate" type="date" class="form-control form-control-sm text-xs" />
        </div>

        <!-- Department Filter -->
        <div>
          <label class="block text-[10px] font-bold text-slate-500 uppercase mb-1">Department</label>
          <select v-model="selectedDept" class="form-control form-control-sm text-xs">
            <option value="">All Departments</option>
            <option v-for="d in departmentOptions" :key="d" :value="d">{{ d }}</option>
          </select>
        </div>

        <!-- Status Filter -->
        <div>
          <label class="block text-[10px] font-bold text-slate-500 uppercase mb-1">Status</label>
          <select v-model="selectedStatus" class="form-control form-control-sm text-xs">
            <option value="">All Statuses</option>
            <option value="Approved">Approved</option>
            <option value="Pending">Pending</option>
            <option value="Checked In">Checked In</option>
            <option value="Checked Out">Checked Out</option>
            <option value="Rejected">Rejected</option>
          </select>
        </div>

      </div>

      <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-2 pt-2 border-t border-slate-100 dark:border-slate-800">
        
        <!-- Search & Preset Controls -->
        <div class="d-flex align-items-center gap-2 flex-grow-1 max-w-md">
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Search reference, visitor name, host officer..." 
            class="form-control form-control-sm text-xs"
          />
          <button type="button" class="btn btn-outline-secondary btn-sm text-xs whitespace-nowrap" @click="resetFilters">
            Clear
          </button>
        </div>

        <!-- Date Presets -->
        <div class="d-flex align-items-center gap-1">
          <span class="text-[10px] text-slate-400 mr-1 uppercase font-bold">Presets:</span>
          <button type="button" class="btn btn-xs btn-outline-primary" @click="setPreset('today')">Today</button>
          <button type="button" class="btn btn-xs btn-outline-primary" @click="setPreset('week')">This Week</button>
          <button type="button" class="btn btn-xs btn-outline-primary" @click="setPreset('month')">This Month</button>
          <button type="button" class="btn btn-xs btn-outline-primary" @click="setPreset('all')">All Time</button>
        </div>

      </div>
    </div>

    <!-- Export Action Controls Bar -->
    <div class="d-flex justify-content-between align-items-center mb-3">
      <div class="text-xs text-slate-500">
        Showing <strong>{{ paginatedData.length }}</strong> of <strong>{{ sortedData.length }}</strong> report records
      </div>

      <div class="d-flex gap-2">
        <button type="button" class="btn btn-sm btn-outline-success" @click="handleExport('excel')" :disabled="!sortedData.length" title="Export as Excel Workbook (.xls)">
          <i class="icofont-file-excel mr-1"></i> Excel
        </button>
        <button type="button" class="btn btn-sm btn-outline-info" @click="handleExport('csv')" :disabled="!sortedData.length" title="Export as CSV (.csv)">
          <i class="icofont-file-text mr-1"></i> CSV
        </button>
        <button type="button" class="btn btn-sm btn-outline-primary" @click="handleExport('docx')" :disabled="!sortedData.length" title="Export as Word Document (.doc)">
          <i class="icofont-file-word mr-1"></i> DOCX
        </button>
        <button type="button" class="btn btn-sm btn-outline-danger" @click="handleExport('pdf')" :disabled="!sortedData.length" title="Export as PDF Document (.pdf)">
          <i class="icofont-file-pdf mr-1"></i> PDF
        </button>
      </div>
    </div>

    <!-- Interactive Sortable Data Table -->
    <div class="table-responsive">
      <table class="table table-bordered table-striped official-table text-xs">
        <thead>
          <tr>
            <th v-for="col in tableColumns" :key="col.key" class="cursor-pointer select-none" @click="toggleSort(col.key)">
              <div class="d-flex justify-content-between align-items-center">
                <span>{{ col.label }}</span>
                <span class="text-slate-400 ml-1">
                  <i v-if="sortKey === col.key && sortOrder === 'asc'" class="icofont-arrow-up text-primary"></i>
                  <i v-else-if="sortKey === col.key && sortOrder === 'desc'" class="icofont-arrow-down text-primary"></i>
                  <i v-else class="icofont-thin-double-arrow text-slate-300"></i>
                </span>
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in paginatedData" :key="row.id">
            <td v-for="col in tableColumns" :key="col.key">
              
              <!-- Badge Formatter for Status -->
              <span v-if="col.key === 'status'" class="badge" :class="statusBadgeClass(row.status)" style="font-size: 9px; padding: 3px 6px;">
                {{ row.status }}
              </span>

              <!-- Default Formatter -->
              <span v-else>{{ row[col.key] || '—' }}</span>
            </td>
          </tr>
          <tr v-if="!sortedData.length">
            <td :colspan="tableColumns.length" class="text-center py-6 text-slate-400">
              <i class="icofont-search-folder text-3xl d-block mb-2 text-slate-300"></i>
              No report records match the selected date range and filter criteria.
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Table Pagination Footer -->
    <div v-if="totalPages > 1" class="d-flex justify-content-between align-items-center mt-3 pt-3 border-t border-slate-100 dark:border-slate-800">
      <span class="text-xs text-slate-500">Page {{ currentPage }} of {{ totalPages }}</span>
      <div class="btn-group">
        <button type="button" class="btn btn-outline-secondary btn-sm" :disabled="currentPage <= 1" @click="currentPage--">Previous</button>
        <button type="button" class="btn btn-outline-secondary btn-sm" :disabled="currentPage >= totalPages" @click="currentPage++">Next</button>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { exportToCSV, exportToExcel, exportToDocx, exportToPDF } from '@/utils/reportExporter.js'

const reportType = ref('visitor-clearance')
const startDate = ref('')
const endDate = ref('')
const selectedDept = ref('')
const selectedStatus = ref('')
const searchQuery = ref('')

const sortKey = ref('arrivalTime')
const sortOrder = ref('desc')

const currentPage = ref(1)
const perPage = ref(10)

const visitorClearancesData = ref([])
const leaveRequestsData = ref([])

const departmentOptions = [
  'Engineering',
  'Helpdesk',
  'Administration',
  'Human Resources',
  'Accounts & Finance',
  'Sports Development',
  'Legal',
  'Secretariat'
]

onMounted(() => {
  loadData()
})

watch([reportType, startDate, endDate, selectedDept, selectedStatus, searchQuery], () => {
  currentPage.value = 1
})

function loadData() {
  // 1. Load Visitor Clearances
  const storedVC = localStorage.getItem('ncsms_visitor_clearances')
  if (storedVC) {
    try {
      visitorClearancesData.value = JSON.parse(storedVC)
    } catch (e) {}
  } else {
    visitorClearancesData.value = [
      { id: 101, reference: 'VC-2026-081', visitorName: 'David Ssemwanga', visitorPhone: '+256 701 554433', visitorEmail: 'd.ssemwanga@ugandabadminton.org', organization: 'Uganda Badminton Association', hostOfficer: 'Dr. Bernard Patrick Ogwel', department: 'Administration', arrivalTime: '2026-07-21T09:30', departureTime: '2026-07-21T12:00', idType: 'National ID', idNumber: 'CM9201940192A', vehicleReg: 'UBL 482A', itemsCarried: 'HP Laptop', status: 'Approved' },
      { id: 102, reference: 'VC-2026-082', visitorName: 'Grace Akello', visitorPhone: '+256 772 119988', visitorEmail: 'g.akello@auditors.co.ug', organization: 'Internal Audit Firm', hostOfficer: 'Joshua Obua', department: 'Accounts & Finance', arrivalTime: '2026-07-22T11:00', departureTime: '2026-07-22T14:30', idType: 'Passport', idNumber: 'A09182391', vehicleReg: 'UBK 109Z', itemsCarried: 'Audit Folders', status: 'Pending' },
      { id: 103, reference: 'VC-2026-083', visitorName: 'Baker Nsubuga', visitorPhone: '+256 700 882211', visitorEmail: 'b.nsubuga@fufa.co.ug', organization: 'FUFA', hostOfficer: 'Peter Ssewankambo', department: 'Sports Development', arrivalTime: '2026-07-20T14:00', departureTime: '2026-07-20T17:00', idType: 'Staff Badge', idNumber: 'FUFA-901', vehicleReg: 'UAX 882M', itemsCarried: 'Tournament Schedules', status: 'Checked In' }
    ]
  }

  // 2. Load Leave Requests
  const storedLeave = localStorage.getItem('ncsms_leave_requests')
  if (storedLeave) {
    try {
      leaveRequestsData.value = JSON.parse(storedLeave)
    } catch (e) {}
  } else {
    leaveRequestsData.value = [
      { id: 201, reference: 'LR-2026-001', applicantName: 'Baker Nsubuga', department: 'Engineering', leaveType: 'Annual Leave', totalDays: 5, startDate: '2026-08-01', endDate: '2026-08-05', status: 'Pending', position: 'Senior Software Engineer' },
      { id: 202, reference: 'LR-2026-002', applicantName: 'Sarah Chemutai', department: 'Human Resources', leaveType: 'Sick Leave', totalDays: 3, startDate: '2026-07-15', endDate: '2026-07-18', status: 'Approved', position: 'HR Manager' }
    ]
  }
}

const tableColumns = computed(() => {
  if (reportType.value === 'visitor-clearance') {
    return [
      { key: 'reference', label: 'Reference' },
      { key: 'visitorName', label: 'Visitor Name' },
      { key: 'organization', label: 'Organization' },
      { key: 'hostOfficer', label: 'Host Officer' },
      { key: 'department', label: 'Department' },
      { key: 'arrivalTime', label: 'Arrival Date/Time' },
      { key: 'vehicleReg', label: 'Vehicle' },
      { key: 'status', label: 'Status' }
    ]
  } else if (reportType.value === 'visitor-logs') {
    return [
      { key: 'reference', label: 'Gate Pass Ref' },
      { key: 'visitorName', label: 'Visitor Name' },
      { key: 'idType', label: 'ID Type' },
      { key: 'idNumber', label: 'ID Number' },
      { key: 'arrivalTime', label: 'Check-In Time' },
      { key: 'departureTime', label: 'Expected Departure' },
      { key: 'itemsCarried', label: 'Hand Luggage / Items' },
      { key: 'status', label: 'Gate Status' }
    ]
  } else {
    return [
      { key: 'reference', label: 'Leave Ref' },
      { key: 'applicantName', label: 'Applicant Name' },
      { key: 'department', label: 'Department' },
      { key: 'leaveType', label: 'Leave Type' },
      { key: 'totalDays', label: 'Days' },
      { key: 'startDate', label: 'Begins' },
      { key: 'endDate', label: 'Ends' },
      { key: 'status', label: 'Approval Status' }
    ]
  }
})

const rawDataSet = computed(() => {
  if (reportType.value === 'visitor-clearance' || reportType.value === 'visitor-logs') {
    return visitorClearancesData.value
  } else {
    return leaveRequestsData.value
  }
})

const filteredData = computed(() => {
  return rawDataSet.value.filter(item => {
    // Date Range Filter
    const itemDate = item.arrivalTime || item.startDate || ''
    if (startDate.value && itemDate && itemDate.slice(0, 10) < startDate.value) return false
    if (endDate.value && itemDate && itemDate.slice(0, 10) > endDate.value) return false

    // Department Filter
    if (selectedDept.value && item.department !== selectedDept.value) return false

    // Status Filter
    if (selectedStatus.value && item.status !== selectedStatus.value) return false

    // Search Query Filter
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase()
      const haystack = [
        item.reference,
        item.visitorName,
        item.applicantName,
        item.organization,
        item.hostOfficer,
        item.department,
        item.idNumber,
        item.vehicleReg
      ].filter(Boolean).join(' ').toLowerCase()
      if (!haystack.includes(q)) return false
    }

    return true
  })
})

const sortedData = computed(() => {
  const list = [...filteredData.value]
  if (!sortKey.value) return list

  return list.sort((a, b) => {
    let valA = a[sortKey.value] || ''
    let valB = b[sortKey.value] || ''
    if (typeof valA === 'string') valA = valA.toLowerCase()
    if (typeof valB === 'string') valB = valB.toLowerCase()

    if (valA < valB) return sortOrder.value === 'asc' ? -1 : 1
    if (valA > valB) return sortOrder.value === 'asc' ? 1 : -1
    return 0
  })
})

const totalPages = computed(() => {
  return Math.ceil(sortedData.value.length / perPage.value) || 1
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage.value
  return sortedData.value.slice(start, start + perPage.value)
})

function toggleSort(key) {
  if (sortKey.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortOrder.value = 'asc'
  }
}

function countStatus(statusName) {
  return filteredData.value.filter(i => i.status === statusName).length
}

function statusBadgeClass(status) {
  switch (status) {
    case 'Approved': return 'badge-success'
    case 'Pending': return 'badge-warning'
    case 'Checked In': return 'badge-primary'
    case 'Checked Out': return 'badge-secondary'
    case 'Rejected': return 'badge-danger'
    default: return 'badge-info'
  }
}

function setPreset(preset) {
  const now = new Date()
  if (preset === 'today') {
    startDate.value = now.toISOString().slice(0, 10)
    endDate.value = now.toISOString().slice(0, 10)
  } else if (preset === 'week') {
    const firstDay = new Date(now.setDate(now.getDate() - now.getDay()))
    const lastDay = new Date(now.setDate(now.getDate() - now.getDay() + 6))
    startDate.value = firstDay.toISOString().slice(0, 10)
    endDate.value = lastDay.toISOString().slice(0, 10)
  } else if (preset === 'month') {
    const y = now.getFullYear()
    const m = now.getMonth()
    startDate.value = new Date(y, m, 1).toISOString().slice(0, 10)
    endDate.value = new Date(y, m + 1, 0).toISOString().slice(0, 10)
  } else if (preset === 'all') {
    startDate.value = ''
    endDate.value = ''
  }
}

function resetFilters() {
  startDate.value = ''
  endDate.value = ''
  selectedDept.value = ''
  selectedStatus.value = ''
  searchQuery.value = ''
}

function handleExport(format) {
  const titleMap = {
    'visitor-clearance': 'Visitor Clearance & Access Pass Report',
    'visitor-logs': 'Gate Security Access Log Report',
    'leave-requests': 'Leave Application & Entitlement Report'
  }
  const title = titleMap[reportType.value] || 'Administrative Report'
  const filename = `${reportType.value}_report_${new Date().toISOString().slice(0, 10)}`
  const headers = tableColumns.value
  const data = sortedData.value

  const metadata = {
    dateRange: (startDate.value || endDate.value) ? `${startDate.value || 'Beginning'} to ${endDate.value || 'Present'}` : 'All Time Records',
    department: selectedDept.value || 'All Departments',
    status: selectedStatus.value || 'All Statuses'
  }

  if (format === 'csv') {
    exportToCSV(filename, headers, data)
  } else if (format === 'excel') {
    exportToExcel(filename, headers, data, title)
  } else if (format === 'docx') {
    exportToDocx(filename, headers, data, title)
  } else if (format === 'pdf') {
    exportToPDF(title, headers, data, metadata)
  }
}
</script>

<style scoped>
.official-table th {
  background: #f8fafc;
  color: #475569;
  font-weight: 700;
  font-size: 11px;
}
</style>
