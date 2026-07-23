<template>
  <div class="leave-approvals-container">
    <header class="page-heading">
      <div>
        <p>Leave Management</p>
        <h1>Leave Approvals</h1>
        <span>View and manage your leave requests and entitlement computations.</span>
      </div>
    </header>

    <!-- Summary & Metrics -->
    <div class="user-kpis mt-4">
      <article>
        <span class="blue"><i class="icofont-calendar"></i></span>
        <div>
          <small>Annual Entitlement</small>
          <strong>30 Days</strong>
          <p>Standard Yearly Rate</p>
        </div>
      </article>
      <article>
        <span class="green"><i class="icofont-check-circled"></i></span>
        <div>
          <small>Leave Days Taken</small>
          <strong>{{ totalDaysTaken }} Days</strong>
          <p>Approved Leave</p>
        </div>
      </article>
      <article>
        <span class="amber"><i class="icofont-clock-time"></i></span>
        <div>
          <small>Leave Balance</small>
          <strong>{{ leaveBalance }} Days</strong>
          <p>Remaining Entitlement</p>
        </div>
      </article>
    </div>

    <!-- Approvals list table card -->
    <div class="card border-0 shadow-sm mt-4 p-4">
      <div class="d-flex justify-content-between align-items-center mb-3">
        <h2 class="section-title-panel m-0">My Applications & History</h2>
        <button type="button" class="btn btn-primary" @click="goToApply">
          <i class="icofont-plus-circle mr-1"></i> Apply for Leave
        </button>
      </div>

      <div v-if="recentRequests.length === 0" class="empty-state text-center py-5">
        <i class="icofont-calendar text-muted" style="font-size: 48px;"></i>
        <p class="text-muted mt-3">No leave applications found on file.</p>
        <button type="button" class="btn btn-primary btn-sm mt-2" @click="goToApply">
          Submit First Application
        </button>
      </div>
      
      <div v-else class="table-responsive">
        <table class="table table-bordered table-striped official-table">
          <thead>
            <tr>
              <th>Applicant</th>
              <th>Leave Type</th>
              <th>Dates Requested</th>
              <th>Days</th>
              <th>Workflow Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="req in recentRequests" :key="req.id">
              <td>
                <strong>{{ req.applicantName }}</strong><br/>
                <small class="text-muted">{{ req.position }} ({{ req.department }})</small>
              </td>
              <td>{{ req.leaveType }}</td>
              <td>
                <i class="icofont-ui-calendar mr-1"></i>
                {{ formatDate(req.startDate) }} - {{ formatDate(req.endDate) }}
              </td>
              <td>{{ req.appliedDays }} Days</td>
              <td>
                <div class="workflow-steps-progress">
                  <div class="status-indicator">
                    <span :class="['badge-status', req.status.toLowerCase()]">{{ req.status }}</span>
                  </div>
                  <ul class="workflow-checklist small mt-2">
                    <li :class="{ completed: req.workflow?.deptHead }">
                      <i class="icofont-check-circled"></i> Dept Head Recommendation
                    </li>
                    <li :class="{ completed: req.workflow?.hrVerified }">
                      <i class="icofont-check-circled"></i> HR Verification
                    </li>
                    <li :class="{ completed: req.workflow?.gsApproved }">
                      <i class="icofont-check-circled"></i> General Secretary Approval
                    </li>
                  </ul>
                </div>
              </td>
              <td>
                <div class="d-flex align-items-center" style="gap: 6px;">
                  <button v-if="req.status === 'Pending'" type="button" class="btn btn-sm btn-outline-primary" @click="editRequest(req.id)" title="Edit this request" style="padding: 4px 8px; font-size: 11px;">
                    <i class="icofont-edit"></i> Edit
                  </button>
                  <button type="button" class="btn btn-sm btn-outline-info" @click="downloadRequestPdf(req)" title="Download / Print PDF" style="padding: 4px 8px; font-size: 11px;">
                    <i class="icofont-print"></i> PDF
                  </button>
                </div>
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

const router = useRouter()
const recentRequests = ref([])

onMounted(() => {
  loadLeaveRequests()
})

const totalDaysTaken = computed(() => {
  return recentRequests.value
    .filter(req => req.status === 'Approved')
    .reduce((sum, req) => sum + Number(req.appliedDays || 0), 0)
})

const leaveBalance = computed(() => {
  return Math.max(0, 30 - totalDaysTaken.value)
})

function loadLeaveRequests() {
  const stored = localStorage.getItem('ncsms_leave_requests')
  if (stored) {
    try {
      const parsed = JSON.parse(stored)
      recentRequests.value = parsed.map(req => ({
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
      }))
    } catch (e) {
      recentRequests.value = []
    }
  } else {
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

function goToApply() {
  router.push('/leave-apply')
}

function editRequest(id) {
  router.push(`/leave-apply?id=${id}`)
}

function downloadRequestPdf(req) {
  router.push(`/leave-apply?id=${req.id}&print=true`)
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

<style scoped>
.leave-approvals-container {
  max-width: 1080px;
  margin: 0 auto;
}

.section-title-paper {
  font-size: 16px;
  font-weight: 700;
  color: #1e293b;
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

.workflow-checklist li.completed {
  color: #10b981;
}

/* User KPI cards layout matching Otika metrics styling */
.user-kpis {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.user-kpis article {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  border-radius: 12px;
  background: #fff;
  border: 1px solid #f0f2f8;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.02);
}

:global(.dark) .user-kpis article {
  background: #1f2937 !important;
  border-color: #334155 !important;
}

.user-kpis span {
  display: grid;
  place-items: center;
  width: 48px;
  height: 48px;
  border-radius: 12px;
  font-size: 20px;
}

.user-kpis span.blue { background: #e0f2fe; color: #0284c7; }
.user-kpis span.green { background: #d1fae5; color: #059669; }
.user-kpis span.amber { background: #fef3c7; color: #d97706; }

.user-kpis div {
  display: flex;
  flex-direction: column;
}

.user-kpis small {
  font-size: 11px;
  color: #64748b;
  font-weight: 700;
  text-transform: uppercase;
}

:global(.dark) .user-kpis small {
  color: #94a3b8 !important;
}

.user-kpis strong {
  font-size: 22px;
  color: #1e293b;
  font-weight: 800;
  margin: 2px 0;
}

:global(.dark) .user-kpis strong {
  color: #f8fafc !important;
}

.user-kpis p {
  font-size: 11px;
  color: #94a3b8;
  margin: 0;
}

@media (max-width: 768px) {
  .user-kpis {
    grid-template-columns: 1fr;
  }
}
</style>
