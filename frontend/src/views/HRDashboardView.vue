<template>
  <main class="otika-cms">
    <div class="otika-app">
      <div class="main-wrapper main-wrapper-1" :class="{ 'sidebar-mini': sidebarCollapsed }">
        <div class="navbar-bg"></div>
        
        <!-- Header / Top Navbar -->
        <nav class="navbar navbar-expand-lg main-navbar sticky">
          <div class="form-inline mr-auto">
            <ul class="navbar-nav mr-3">
              <li>
                <button type="button" class="nav-link nav-link-lg cms-top-icon collapse-btn" @click="sidebarCollapsed = !sidebarCollapsed" :title="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'">
                  <i class="icofont-navigation-menu"></i>
                </button>
              </li>
              <li>
                <button type="button" class="nav-link nav-link-lg cms-top-icon fullscreen-btn" @click="resetDashboard" title="Refresh">
                  <i class="icofont-refresh"></i>
                </button>
              </li>
            </ul>
          </div>
          <ul class="navbar-nav navbar-right">
            <li>
              <ThemeToggle class="nav-link nav-link-lg cms-top-icon portal-theme-toggle" />
            </li>
            <li class="dropdown" :class="{ show: profileOpen }">
              <button type="button" class="nav-link dropdown-toggle nav-link-lg nav-link-user cms-top-icon" @click="profileOpen = !profileOpen">
                <img alt="" src="/otika-assets/img/user.png" class="user-img-radious-style" />
                <span class="d-sm-none d-lg-inline-block"></span>
              </button>
              <div class="dropdown-menu dropdown-menu-right pullDown" :class="{ show: profileOpen }">
                <div class="dropdown-title">Hello {{ user.first_name }}</div>
                <button type="button" class="dropdown-item has-icon text-danger" @click="logout">
                  <i class="fas fa-sign-out-alt"></i> Logout
                </button>
              </div>
            </li>
          </ul>
        </nav>

        <!-- Sidebar Navigation -->
        <div class="main-sidebar sidebar-style-2">
          <aside id="sidebar-wrapper">
            <div class="sidebar-brand">
              <router-link to="/portal/hr">
                <img alt="NCS" src="/main-logo.png" class="header-logo" />
              </router-link>
            </div>
            <ul class="sidebar-menu">
              <li class="menu-header">HR Workspace</li>
              <li :class="{ active: activeTab === 'overview' }">
                <button type="button" class="nav-link" @click="activeTab = 'overview'">
                  <i class="icofont-dashboard"></i><span>Dashboard Overview</span>
                </button>
              </li>
              <li class="menu-header">Universal Applications</li>
              <li>
                <router-link to="/reception/leave/apply" class="nav-link">
                  <i class="icofont-calendar"></i><span>Leave Application</span>
                </router-link>
              </li>
              <li>
                <router-link to="/it/ppda/new" class="nav-link">
                  <i class="icofont-document-folder"></i><span>PPDA Application</span>
                </router-link>
              </li>
              <li class="menu-header">HR Workspace</li>
              <li :class="{ active: activeTab === 'registry' }">
                <button type="button" class="nav-link" @click="activeTab = 'registry'">
                  <i class="icofont-users-alt-5"></i><span>Staff Registry</span>
                </button>
              </li>
              <li :class="{ active: activeTab === 'leaves' }">
                <button type="button" class="nav-link" @click="activeTab = 'leaves'">
                  <i class="icofont-calendar"></i><span>Leave Requests</span>
                </button>
              </li>
              <li :class="{ active: activeTab === 'documents' }">
                <button type="button" class="nav-link" @click="activeTab = 'documents'">
                  <i class="icofont-document-folder"></i><span>Documents & Templates</span>
                </button>
              </li>
              <li class="menu-header">Session</li>
              <li>
                <button type="button" class="nav-link text-danger" @click="logout">
                  <i class="icofont-logout text-danger"></i><span>Log Out</span>
                </button>
              </li>
            </ul>
          </aside>
        </div>

        <!-- Main Content Area -->
        <div class="main-content">
          <section class="section">
            <div class="section-body">
              <UniversalDashboardApplications class="mb-4" />
              
              <!-- Tab 1: Dashboard Overview -->
              <div v-if="activeTab === 'overview'" class="space-y-6">
                <!-- Welcome Widget -->
                <div class="card bg-gradient-to-r from-blue-900 via-indigo-900 to-slate-900 border-0 text-white p-6 shadow-sm rounded-lg mb-6">
                  <h2 class="text-xl font-bold">Welcome back, {{ user.first_name }}!</h2>
                  <p class="text-sm text-slate-300 mt-1">Here is a quick snapshot of the staff registry and pending activities for today.</p>
                </div>

                <!-- KPI Grid -->
                <div class="row">
                  <div class="col-xl-3 col-lg-6 col-md-6 col-sm-6 col-xs-12">
                    <div class="card card-statistic-4 hover:shadow-lg transition-shadow">
                      <div class="card-body">
                        <h5 class="text-xs text-slate-400 font-semibold uppercase">Active Employees</h5>
                        <h2 class="text-2xl font-bold mt-2 text-slate-800 dark:text-slate-100">{{ staff.length }}</h2>
                        <p class="text-xs text-emerald-500 mt-1"><i class="fas fa-arrow-up"></i> +4 this month</p>
                      </div>
                    </div>
                  </div>
                  <div class="col-xl-3 col-lg-6 col-md-6 col-sm-6 col-xs-12">
                    <div class="card card-statistic-4 hover:shadow-lg transition-shadow">
                      <div class="card-body">
                        <h5 class="text-xs text-slate-400 font-semibold uppercase">Open Positions</h5>
                        <h2 class="text-2xl font-bold mt-2 text-slate-800 dark:text-slate-100">9</h2>
                        <p class="text-xs text-indigo-500 mt-1">3 active listings</p>
                      </div>
                    </div>
                  </div>
                  <div class="col-xl-3 col-lg-6 col-md-6 col-sm-6 col-xs-12">
                    <div class="card card-statistic-4 hover:shadow-lg transition-shadow">
                      <div class="card-body">
                        <h5 class="text-xs text-slate-400 font-semibold uppercase">Leave Requests</h5>
                        <h2 class="text-2xl font-bold mt-2 text-slate-800 dark:text-slate-100">{{ leaveRequests.length }}</h2>
                        <p class="text-xs text-amber-500 mt-1">Pending approval</p>
                      </div>
                    </div>
                  </div>
                  <div class="col-xl-3 col-lg-6 col-md-6 col-sm-6 col-xs-12">
                    <div class="card card-statistic-4 hover:shadow-lg transition-shadow">
                      <div class="card-body">
                        <h5 class="text-xs text-slate-400 font-semibold uppercase">Monthly Payroll</h5>
                        <h2 class="text-2xl font-bold mt-2 text-slate-800 dark:text-slate-100">UGX 42M</h2>
                        <p class="text-xs text-emerald-500 mt-1">Processed successfully</p>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="row mt-4">
                  <!-- Leave approvals snippet -->
                  <div class="col-lg-6 col-md-12">
                    <div class="card p-4">
                      <div class="card-header p-0 pb-3 mb-3 border-b border-slate-100 dark:border-slate-800 flex justify-between items-center">
                        <h4 class="text-sm font-bold text-slate-800 dark:text-slate-100">Recent Leave Requests</h4>
                        <button class="btn btn-sm btn-link" @click="activeTab = 'leaves'">View all</button>
                      </div>
                      <div class="space-y-3">
                        <div v-for="leave in leaveRequests" :key="leave.id" class="p-3 rounded border border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/30 flex justify-between items-center">
                          <div>
                            <p class="text-sm font-semibold">{{ leave.name }}</p>
                            <p class="text-xs text-slate-400">{{ leave.type }} • {{ leave.duration }}</p>
                          </div>
                          <button class="btn btn-xs btn-outline-primary" @click="activeTab = 'leaves'">Review</button>
                        </div>
                        <p v-if="!leaveRequests.length" class="text-xs text-slate-400 text-center py-4">All caught up!</p>
                      </div>
                    </div>
                  </div>

                  <!-- Quick Stats / Actions -->
                  <div class="col-lg-6 col-md-12">
                    <div class="card p-4">
                      <div class="card-header p-0 pb-3 mb-3 border-b border-slate-100 dark:border-slate-800">
                        <h4 class="text-sm font-bold text-slate-800 dark:text-slate-100">Workspace Activities</h4>
                      </div>
                      <div class="space-y-3">
                        <button @click="activeTab = 'registry'; showAddEmployee = true" class="w-full text-left p-3 rounded border border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900/40 transition-colors flex items-center justify-between">
                          <div>
                            <p class="text-sm font-semibold text-slate-700 dark:text-slate-300">Add New Employee Profile</p>
                            <p class="text-xs text-slate-400">Initialize a new employee contract & credentials</p>
                          </div>
                          <i class="fas fa-chevron-right text-slate-400"></i>
                        </button>
                        <button @click="activeTab = 'documents'" class="w-full text-left p-3 rounded border border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900/40 transition-colors flex items-center justify-between">
                          <div>
                            <p class="text-sm font-semibold text-slate-700 dark:text-slate-300">View HR templates</p>
                            <p class="text-xs text-slate-400">Contracts, handbooks, and scale matrices</p>
                          </div>
                          <i class="fas fa-chevron-right text-slate-400"></i>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Tab 2: Staff Registry -->
              <div v-else-if="activeTab === 'registry'" class="card p-5">
                <div class="flex justify-between items-center mb-4">
                  <h3 class="text-base font-bold text-slate-800 dark:text-slate-100">Employee Profiles</h3>
                  <button @click="showAddEmployee = true" class="btn btn-primary btn-sm rounded-lg">+ Add Employee</button>
                </div>
                <div class="overflow-x-auto">
                  <table class="table table-striped text-sm">
                    <thead>
                      <tr class="text-xs font-semibold text-slate-400 uppercase">
                        <th>Employee</th>
                        <th>Email</th>
                        <th>Department</th>
                        <th>Status</th>
                        <th class="text-right">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="emp in staff" :key="emp.id">
                        <td><strong>{{ emp.firstName }} {{ emp.lastName }}</strong></td>
                        <td>{{ emp.email }}</td>
                        <td>{{ emp.department }}</td>
                        <td>
                          <span :class="emp.status === 'Active' ? 'badge badge-success' : 'badge badge-warning'">{{ emp.status }}</span>
                        </td>
                        <td class="text-right">
                          <button class="btn btn-xs btn-light">Edit</button>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- Tab 3: Leave Applications -->
              <div v-else-if="activeTab === 'leaves'" class="card p-5">
                <h3 class="text-base font-bold text-slate-800 dark:text-slate-100 mb-4">Pending Leave Reviews</h3>
                <div class="row">
                  <div v-for="leave in leaveRequests" :key="leave.id" class="col-md-6 col-sm-12 mb-4">
                    <div class="card border border-slate-100 dark:border-slate-800 p-4">
                      <div class="flex justify-between items-start">
                        <div>
                          <h4 class="text-sm font-semibold text-slate-800 dark:text-slate-100">{{ leave.name }}</h4>
                          <p class="text-xs text-slate-400 mt-1">{{ leave.type }} • {{ leave.duration }}</p>
                          <p class="text-xs text-indigo-500 mt-2 font-medium">Dates: {{ leave.date }}</p>
                        </div>
                      </div>
                      <div class="flex gap-2 justify-end mt-4 pt-3 border-t border-slate-50 dark:border-slate-800">
                        <button class="btn btn-xs btn-outline-danger" @click="resolveLeave(leave.id, 'rejected')">Reject</button>
                        <button class="btn btn-xs btn-success" @click="resolveLeave(leave.id, 'approved')">Approve</button>
                      </div>
                    </div>
                  </div>
                  <p v-if="!leaveRequests.length" class="text-sm text-slate-400 text-center py-8 col-12">No pending leave applications at this time.</p>
                </div>
              </div>

              <!-- Tab 4: Documents & Templates -->
              <div v-else-if="activeTab === 'documents'" class="card p-5">
                <h3 class="text-base font-bold text-slate-800 dark:text-slate-100 mb-4">HR Document Library</h3>
                <div class="row">
                  <div class="col-md-4 col-sm-6 col-12 mb-4">
                    <div class="card p-4 border border-slate-100 dark:border-slate-800 hover:border-indigo-500/30 transition-all text-center">
                      <div class="w-12 h-12 rounded-full bg-red-500/10 text-red-500 flex items-center justify-center mx-auto mb-3">
                        <i class="fas fa-file-pdf text-xl"></i>
                      </div>
                      <h4 class="text-xs font-bold truncate">NCS Employee Handbook 2026</h4>
                      <p class="text-[10px] text-slate-400 mt-1">PDF document • 4.2 MB</p>
                      <a href="#" class="btn btn-xs btn-block btn-outline-primary mt-3">Download</a>
                    </div>
                  </div>
                  <div class="col-md-4 col-sm-6 col-12 mb-4">
                    <div class="card p-4 border border-slate-100 dark:border-slate-800 hover:border-indigo-500/30 transition-all text-center">
                      <div class="w-12 h-12 rounded-full bg-blue-500/10 text-blue-500 flex items-center justify-center mx-auto mb-3">
                        <i class="fas fa-file-word text-xl"></i>
                      </div>
                      <h4 class="text-xs font-bold truncate">Staff Appraisal Template</h4>
                      <p class="text-[10px] text-slate-400 mt-1">DOCX template • 248 KB</p>
                      <a href="#" class="btn btn-xs btn-block btn-outline-primary mt-3">Download</a>
                    </div>
                  </div>
                  <div class="col-md-4 col-sm-6 col-12 mb-4">
                    <div class="card p-4 border border-slate-100 dark:border-slate-800 hover:border-indigo-500/30 transition-all text-center">
                      <div class="w-12 h-12 rounded-full bg-emerald-500/10 text-emerald-500 flex items-center justify-center mx-auto mb-3">
                        <i class="fas fa-file-excel text-xl"></i>
                      </div>
                      <h4 class="text-xs font-bold truncate">Salary Scale Structure</h4>
                      <p class="text-[10px] text-slate-400 mt-1">XLSX sheet • 1.5 MB</p>
                      <a href="#" class="btn btn-xs btn-block btn-outline-primary mt-3">Download</a>
                    </div>
                  </div>
                </div>
              </div>

            </div>
          </section>
        </div>

        <!-- App Footer -->
        <footer class="main-footer cms-main-footer">
          <div class="footer-left">
            Copyright &copy; 2026 <div class="bullet"></div> National Council of Sports, Uganda
          </div>
          <div class="footer-right">
            HR Module v1.0.0
          </div>
        </footer>
      </div>
    </div>

    <!-- Add Employee Modal -->
    <div v-if="showAddEmployee" class="fixed inset-0 z-[1080] flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
      <div class="w-full max-w-md bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-6 space-y-4 shadow-xl">
        <div class="flex justify-between items-center pb-2 border-b border-slate-100 dark:border-slate-800">
          <h3 class="text-sm font-bold text-slate-800 dark:text-slate-100">Add New Employee Profile</h3>
          <button @click="showAddEmployee = false" class="text-slate-400 hover:text-slate-600 text-lg">&times;</button>
        </div>
        <form @submit.prevent="addEmployee" class="space-y-4">
          <div>
            <label class="block text-[10px] font-bold text-slate-400 uppercase mb-1">First Name</label>
            <input v-model="newEmp.firstName" type="text" class="form-control" required />
          </div>
          <div>
            <label class="block text-[10px] font-bold text-slate-400 uppercase mb-1">Last Name</label>
            <input v-model="newEmp.lastName" type="text" class="form-control" required />
          </div>
          <div>
            <label class="block text-[10px] font-bold text-slate-400 uppercase mb-1">Email Address</label>
            <input v-model="newEmp.email" type="email" class="form-control" required />
          </div>
          <div>
            <label class="block text-[10px] font-bold text-slate-400 uppercase mb-1">Department</label>
            <select v-model="newEmp.department" class="form-control">
              <option>Administration</option>
              <option>Finance</option>
              <option>Technical & Sports development</option>
              <option>Public Relations</option>
              <option>Helpdesk Services</option>
            </select>
          </div>
          <div class="flex justify-end gap-3 pt-3 border-t border-slate-100 dark:border-slate-800">
            <button type="button" @click="showAddEmployee = false" class="btn btn-sm btn-outline-light">Cancel</button>
            <button type="submit" class="btn btn-sm btn-primary">Add Staff</button>
          </div>
        </form>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { storedPortalUser } from '@/utils/portalAuth.js'
import ThemeToggle from '@/components/theme/ThemeToggle.vue'
import UniversalDashboardApplications from '@/components/layout/UniversalDashboardApplications.vue'
import { ensureOtikaStyles } from '@/utils/otikaAssets.js'

ensureOtikaStyles()

const router = useRouter()
const user = ref(storedPortalUser())

const sidebarCollapsed = ref(false)
const profileOpen = ref(false)
const activeTab = ref('overview')
const showAddEmployee = ref(false)

const activeTabLabel = computed(() => {
  switch (activeTab.value) {
    case 'leaves': return 'Leave Requests'
    case 'documents': return 'Library'
    default: return 'Overview'
  }
})

const staff = ref([
  { id: 1, firstName: 'Aisha', lastName: 'Namuganza', email: 'a.namuganza@ncs.go.ug', department: 'Administration', status: 'Active' },
  { id: 2, firstName: 'Joshua', lastName: 'Obua', email: 'j.obua@ncs.go.ug', department: 'Finance', status: 'Active' },
  { id: 3, firstName: 'Isaac', lastName: 'Kiprotich', email: 'i.kiprotich@ncs.go.ug', department: 'Technical & Sports development', status: 'Active' },
  { id: 4, firstName: 'Sylvia', lastName: 'Alitwala', email: 's.alitwala@ncs.go.ug', department: 'Public Relations', status: 'On Leave' }
])

const leaveRequests = ref([])

onMounted(() => {
  loadHRLeaveRequests()
})

function loadHRLeaveRequests() {
  const stored = localStorage.getItem('ncsms_leave_requests')
  if (stored) {
    try {
      const parsed = JSON.parse(stored)
      leaveRequests.value = parsed.filter(req => req.status === 'Pending').map(req => ({
        id: req.id,
        name: req.applicantName || 'Staff Member',
        type: `${req.leaveType} Leave`,
        duration: `${req.appliedDays} days`,
        date: `${formatDate(req.startDate)} - ${formatDate(req.endDate)}`,
        rawRequest: req
      }))
    } catch (e) {
      leaveRequests.value = []
    }
  } else {
    leaveRequests.value = []
  }
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

const newEmp = ref({
  firstName: '',
  lastName: '',
  email: '',
  department: 'Administration'
})

const addEmployee = () => {
  staff.value.push({
    id: staff.value.length + 1,
    firstName: newEmp.value.firstName,
    lastName: newEmp.value.lastName,
    email: newEmp.value.email,
    department: newEmp.value.department,
    status: 'Active'
  })
  newEmp.value = { firstName: '', lastName: '', email: '', department: 'Administration' }
  showAddEmployee.value = false
}

const resolveLeave = (id, actionStatus) => {
  const stored = localStorage.getItem('ncsms_leave_requests')
  if (stored) {
    try {
      const parsed = JSON.parse(stored)
      const reqIdx = parsed.findIndex(r => r.id === id)
      if (reqIdx !== -1) {
        const req = parsed[reqIdx]
        const finalStatus = actionStatus === 'approved' ? 'Approved' : 'Rejected'
        req.status = finalStatus
        req.workflow = {
          deptHead: actionStatus === 'approved',
          hrVerified: actionStatus === 'approved',
          gsApproved: actionStatus === 'approved'
        }
        localStorage.setItem('ncsms_leave_requests', JSON.stringify(parsed))
        
        // Save notification to local notifications array in local storage
        const notification = {
          id: 'notif_' + Date.now(),
          title: `Leave Request ${finalStatus}`,
          message: `Your leave request for ${req.leaveType} has been ${finalStatus.toLowerCase()} by HR.`,
          status: 'unread',
          icon_key: actionStatus === 'approved' ? 'icofont-check-circled' : 'icofont-close-circled',
          created_at: new Date().toISOString()
        }
        let notifs = []
        const storedNotifs = localStorage.getItem('ncsms_local_notifications')
        if (storedNotifs) {
          try { notifs = JSON.parse(storedNotifs) } catch (e) {}
        }
        notifs.unshift(notification)
        localStorage.setItem('ncsms_local_notifications', JSON.stringify(notifs))
      }
    } catch (e) {}
  }
  loadHRLeaveRequests()
}

const resetDashboard = () => {
  activeTab.value = 'overview'
}

const logout = () => {
  localStorage.removeItem('ncsms_access_token')
  localStorage.removeItem('ncsms_user')
  router.push('/login')
}
</script>

<style scoped>
/* Standard Otika layouts alignments */
.navbar-bg {
  z-index: 1 !important;
}
.main-sidebar {
  z-index: 5 !important;
}
.sidebar-user img {
  border: 2px solid #fff;
  box-shadow: 0 4px 8px rgba(0,0,0,0.05);
}
</style>
