<template>
  <div class="space-y-6">
    <div class="card p-6 shadow-sm border border-slate-100 dark:border-slate-800">
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <div>
          <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
            <i class="icofont-badge text-indigo-600"></i>
            Visitors Clearance Requests
          </h3>
          <p class="text-xs text-slate-400 mt-1">Manage official visitor access permits, security clearance approvals, and gate passes.</p>
        </div>
        <div class="flex gap-2">
          <button type="button" class="btn btn-primary btn-sm flex items-center gap-1 shadow-sm" @click="createNewApplication">
            <i class="icofont-plus-circle"></i>
            <span>Create New Application</span>
          </button>
        </div>
      </div>

      <!-- Filter Controls -->
      <div class="flex flex-wrap gap-2 mb-4 pb-3 border-b border-slate-100 dark:border-slate-800">
        <button 
          v-for="status in ['All', 'Pending', 'Approved', 'Checked In']" 
          :key="status" 
          type="button" 
          class="btn btn-xs rounded-full px-3 transition-colors"
          :class="activeFilter === status ? 'btn-primary' : 'btn-outline-secondary'"
          @click="activeFilter = status"
        >
          {{ status }}
        </button>
      </div>

      <!-- Table Listing -->
      <div class="table-responsive">
        <table class="table table-bordered table-striped official-table">
          <thead>
            <tr>
              <th>Clearance Ref</th>
              <th>Visitor Details</th>
              <th>Organization</th>
              <th>Host Officer / Dept</th>
              <th>Arrival Time</th>
              <th>Vehicle Reg</th>
              <th>Status</th>
              <th class="text-center no-print">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in filteredRequests" :key="item.id">
              <td>
                <span class="font-mono text-xs font-bold text-indigo-600">#{{ item.reference || 'VC-' + item.id }}</span>
              </td>
              <td>
                <strong class="d-block text-slate-800 dark:text-slate-200">{{ item.visitorName }}</strong>
                <span class="text-[10px] text-slate-400 d-block">{{ item.visitorPhone }} • {{ item.visitorEmail }}</span>
              </td>
              <td>{{ item.organization || 'Independent' }}</td>
              <td>
                <strong class="d-block text-slate-700 dark:text-slate-300 text-xs">{{ item.hostOfficer }}</strong>
                <span class="text-[10px] text-slate-400 d-block">{{ item.department }}</span>
              </td>
              <td>
                <span class="text-xs text-slate-600 dark:text-slate-300 font-medium">{{ formatDate(item.arrivalTime) }}</span>
              </td>
              <td>
                <code class="text-xs">{{ item.vehicleReg || 'N/A' }}</code>
              </td>
              <td>
                <span :class="getStatusBadgeClass(item.status)">
                  {{ item.status }}
                </span>
              </td>
              <td class="text-center no-print">
                <div class="btn-group">
                  <button 
                    v-if="item.status === 'Pending' || item.status === 'Draft'"
                    type="button" 
                    class="btn btn-xs btn-outline-primary" 
                    title="Edit Request"
                    @click="editApplication(item.id)"
                  >
                    <i class="icofont-ui-edit"></i> Edit
                  </button>
                  <button 
                    type="button" 
                    class="btn btn-xs btn-outline-info" 
                    title="Print Gate Pass PDF"
                    @click="printApplication(item.id)"
                  >
                    <i class="icofont-file-pdf"></i> PDF
                  </button>
                  <button 
                    v-if="item.status === 'Pending'"
                    type="button" 
                    class="btn btn-xs btn-success" 
                    title="Approve Clearance"
                    @click="updateStatus(item.id, 'Approved')"
                  >
                    Approve
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="!filteredRequests.length">
              <td colspan="8" class="text-center py-8 text-slate-400 text-sm">
                <i class="icofont-badge text-3xl d-block mb-2 text-slate-300"></i>
                No visitor clearance records found matching this status filter.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { storedPortalUser } from '@/utils/portalAuth.js'
import { createAuditLog } from '@/utils/auditLogger.js'

const router = useRouter()
const activeFilter = ref('All')
const visitorClearances = ref([])

onMounted(() => {
  loadVisitorClearances()
})

function loadVisitorClearances() {
  const stored = localStorage.getItem('ncsms_visitor_clearances')
  if (stored) {
    try {
      visitorClearances.value = JSON.parse(stored)
    } catch (e) {
      visitorClearances.value = []
    }
  } else {
    // Seed default visitor clearance records if empty
    const seed = [
      {
        id: 101,
        reference: 'VC-2026-081',
        visitorName: 'David Ssemwanga',
        visitorPhone: '+256 701 554433',
        visitorEmail: 'd.ssemwanga@ugandabadminton.org',
        organization: 'Uganda Badminton Association',
        hostOfficer: 'Dr. Bernard Patrick Ogwel',
        department: 'Administration',
        arrivalTime: '2026-07-21T09:30',
        vehicleReg: 'UBL 482A',
        status: 'Approved'
      },
      {
        id: 102,
        reference: 'VC-2026-082',
        visitorName: 'Grace Akello',
        visitorPhone: '+256 772 119988',
        visitorEmail: 'g.akello@auditors.co.ug',
        organization: 'Internal Audit Firm',
        hostOfficer: 'Joshua Obua',
        department: 'Accounts & Finance',
        arrivalTime: '2026-07-22T11:00',
        vehicleReg: 'UBK 109Z',
        status: 'Pending'
      }
    ]
    visitorClearances.value = seed
    localStorage.setItem('ncsms_visitor_clearances', JSON.stringify(seed))
  }
}

const filteredRequests = computed(() => {
  if (activeFilter.value === 'All') return visitorClearances.value
  return visitorClearances.value.filter(req => req.status === activeFilter.value)
})

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function getStatusBadgeClass(status) {
  switch (status) {
    case 'Approved': return 'badge badge-success'
    case 'Pending': return 'badge badge-warning'
    case 'Checked In': return 'badge badge-info'
    case 'Rejected': return 'badge badge-danger'
    default: return 'badge badge-secondary'
  }
}

function createNewApplication() {
  router.push('/visitor-clearance-apply')
}

function editApplication(id) {
  router.push(`/visitor-clearance-apply?id=${id}`)
}

function printApplication(id) {
  router.push(`/visitor-clearance-apply?id=${id}&print=true`)
}

function updateStatus(id, newStatus) {
  const idx = visitorClearances.value.findIndex(item => item.id === id)
  if (idx !== -1) {
    const oldStatus = visitorClearances.value[idx].status
    visitorClearances.value[idx].status = newStatus
    localStorage.setItem('ncsms_visitor_clearances', JSON.stringify(visitorClearances.value))
    
    // Create Audit Log
    createAuditLog({
      action: newStatus === 'Approved' ? 'visitor_clearance:approve' : 'visitor_clearance:status_change',
      status: 'SUCCESS',
      severity: 'INFO',
      actor: storedPortalUser(),
      resourceId: visitorClearances.value[idx].reference || String(id),
      resourceType: 'VisitorClearance',
      payloadBefore: { status: oldStatus },
      payloadAfter: { status: newStatus }
    })

    // Alert Notification
    const notif = {
      id: 'notif_' + Date.now(),
      title: `Visitor Clearance ${newStatus}`,
      message: `Visitor Clearance #${visitorClearances.value[idx].reference || id} for ${visitorClearances.value[idx].visitorName} has been ${newStatus.toLowerCase()}.`,
      status: 'unread',
      icon_key: 'icofont-badge',
      created_at: new Date().toISOString()
    }
    let notifs = []
    const storedNotifs = localStorage.getItem('ncsms_local_notifications')
    if (storedNotifs) {
      try { notifs = JSON.parse(storedNotifs) } catch (e) {}
    }
    notifs.unshift(notif)
    localStorage.setItem('ncsms_local_notifications', JSON.stringify(notifs))
  }
}
</script>

<style scoped>
.official-table th {
  background: #f8fafc;
  color: #475569;
  font-weight: 700;
  font-size: 11px;
  text-transform: uppercase;
}
</style>
