<template>
  <main class="otika-cms">
    <div class="otika-app">
      <div class="main-wrapper main-wrapper-1" :class="{ 'sidebar-mini': sidebarCollapsed }">
        <div class="navbar-bg no-print"></div>

        <!-- Dynamic Header navbar -->
        <nav class="navbar navbar-expand-lg main-navbar sticky no-print">
          <div class="form-inline mr-auto">
            <ul class="navbar-nav mr-3">
              <li>
                <button type="button" class="nav-link nav-link-lg cms-top-icon collapse-btn" @click="sidebarCollapsed = !sidebarCollapsed" :title="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'">
                  <i class="icofont-navigation-menu"></i>
                </button>
              </li>
            </ul>
          </div>
          <ul class="navbar-nav navbar-right align-items-center">
            <!-- Notifications Dropdown -->
            <li v-if="isHelpdesk" class="dropdown dropdown-list-toggle" :class="{ show: notificationsOpen }">
              <button type="button" class="nav-link nav-link-lg message-toggle cms-top-icon" @click="notificationsOpen = !notificationsOpen" title="Notifications" style="position: relative;">
                <i class="icofont-notification"></i>
                <span v-if="unreadNotificationsCount > 0" class="badge badge-primary" style="position: absolute; top: 12px; right: 2px; padding: 2px 4px; font-size: 8px;">
                  {{ unreadNotificationsCount }}
                </span>
              </button>
              <div class="dropdown-menu dropdown-list dropdown-menu-right pullDown" :class="{ show: notificationsOpen }" style="width: 320px;">
                <div class="dropdown-header">
                  Notifications
                  <div class="float-right">
                    <button type="button" class="btn btn-link btn-sm p-0 text-primary" style="font-size: 10px; text-decoration: none;" @click="markAllRead">Mark All As Read</button>
                  </div>
                </div>
                <div class="dropdown-list-content dropdown-list-icons" style="max-height: 280px; overflow-y: auto; padding: 10px;">
                  <div v-if="notifications.length === 0" class="text-center py-4 text-muted small">
                    No new notifications.
                  </div>
                  <div v-else v-for="notif in notifications" :key="notif.id" class="dropdown-item d-flex align-items-start p-2 border-bottom" style="gap: 10px; cursor: pointer;" @click="goHelpdeskTab('notifications')">
                    <div class="p-2 rounded-circle" :style="notif.status === 'unread' ? 'background: #eff6ff;' : 'background: #f8fafc;'" style="font-size: 14px;">
                      <i :class="notif.icon_key || 'icofont-notification'" :style="notif.status === 'unread' ? 'color: #2563eb;' : 'color: #94a3b8;'"></i>
                    </div>
                    <div>
                      <strong class="d-block text-xs" :style="notif.status === 'unread' ? 'color: #1e293b; font-weight: 700;' : 'color: #64748b; font-weight: 400;'">{{ notif.title }}</strong>
                      <span class="d-block text-[10px] text-slate-400 mt-1 leading-tight">{{ notif.message }}</span>
                    </div>
                  </div>
                </div>
                <div class="dropdown-footer text-center border-top">
                  <button type="button" class="btn btn-link btn-sm text-xs" @click="goHelpdeskTab('notifications')">View All Notifications</button>
                </div>
              </div>
            </li>
            <li>
              <ThemeToggle class="nav-link nav-link-lg cms-top-icon portal-theme-toggle" />
            </li>
            <li class="dropdown" :class="{ show: profileOpen }">
              <button type="button" class="nav-link dropdown-toggle nav-link-lg nav-link-user cms-top-icon" @click="profileOpen = !profileOpen">
                <img alt="" src="/otika-assets/img/user.png" class="user-img-radious-style" />
              </button>
              <div class="dropdown-menu dropdown-menu-right pullDown" :class="{ show: profileOpen }">
                <div class="dropdown-title">Hello {{ profileName }}</div>
                <button type="button" class="dropdown-item has-icon text-danger" @click="logout">
                  <i class="fas fa-sign-out-alt"></i> Logout
                </button>
              </div>
            </li>
          </ul>
        </nav>

        <!-- Dynamic Sidebar Navigation -->
        <div class="main-sidebar sidebar-style-2 no-print">
          <aside id="sidebar-wrapper">
            <div class="sidebar-brand">
              <router-link :to="dashboardLink">
                <img alt="NCS" src="/main-logo.png" class="header-logo" />
              </router-link>
            </div>
            
            <!-- Sidebar for Helpdesk -->
            <ul v-if="isHelpdesk" class="sidebar-menu">
              <li class="menu-header">Helpdesk Workspace</li>
              <li>
                <button type="button" class="nav-link" @click="goHelpdeskTab('overview')">
                  <i class="icofont-dashboard"></i><span>Dashboard Overview</span>
                </button>
              </li>
              <li>
                <button type="button" class="nav-link" @click="goHelpdeskTab('queue')">
                  <i class="icofont-ticket"></i><span>Issues Queue</span>
                </button>
              </li>
              <li class="active">
                <button type="button" class="nav-link" @click="goHelpdeskTab('leave-approvals')">
                  <i class="icofont-calendar"></i><span>Leave Approvals</span>
                </button>
              </li>
              <li>
                <button type="button" class="nav-link" @click="goHelpdeskTab('notifications')">
                  <i class="icofont-notification"></i><span>Notifications</span>
                  <span v-if="unreadNotificationsCount > 0" class="badge badge-primary">{{ unreadNotificationsCount }}</span>
                </button>
              </li>
              <li>
                <button type="button" class="nav-link" @click="goHelpdeskTab('activities')">
                  <i class="icofont-history"></i><span>My Activities</span>
                </button>
              </li>
              <li>
                <button type="button" class="nav-link" @click="goHelpdeskTab('reports')">
                  <i class="icofont-chart-histogram"></i><span>Reports & Analytics</span>
                </button>
              </li>
              <li class="menu-header">Session</li>
              <li>
                <button type="button" class="nav-link text-danger" @click="logout">
                  <i class="icofont-logout text-danger"></i><span>Log Out</span>
                </button>
              </li>
            </ul>
            
            <!-- Sidebar for Ordinary Users -->
            <ul v-else class="sidebar-menu">
              <li class="menu-header">Applicant workspace</li>
              <li v-for="item in ordinaryNavigation" :key="item.id" :class="{ active: item.id === 'leave-approvals' }">
                <button type="button" class="nav-link" @click="goUserSection(item.id)">
                  <i :class="item.icon"></i><span>{{ item.label }}</span>
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

        <!-- Main Content Area fitting the screen layout -->
        <div class="main-content">
          <section class="section py-4 px-3">
            <div class="section-body">
              <LeaveRequestForm />
            </div>
          </section>
        </div>

        <!-- App Footer -->
        <footer class="main-footer cms-main-footer no-print">
          <div class="footer-left">
            Copyright &copy; 2026 <div class="bullet"></div> National Council of Sports, Uganda
          </div>
        </footer>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ThemeToggle from '@/components/theme/ThemeToggle.vue'
import LeaveRequestForm from '@/components/portal/LeaveRequestForm.vue'
import { storedPortalUser, roleNames } from '@/utils/portalAuth.js'
import { ensureOtikaStyles } from '@/utils/otikaAssets.js'
import * as cms from '@/api/cms.js'

ensureOtikaStyles()

const router = useRouter()
const user = ref(storedPortalUser() || {})
const roles = roleNames(user.value)

const sidebarCollapsed = ref(false)
const profileOpen = ref(false)

const notificationsOpen = ref(false)
const notifications = ref([])

const unreadNotificationsCount = computed(() => {
  return notifications.value.filter(n => n.status === 'unread').length
})

async function fetchNotifications() {
  let list = []
  try {
    const res = await cms.listNotifications({ page: 1, per_page: 50 })
    const apiData = res?.data || res?.entries || res || []
    const parsedApi = (Array.isArray(apiData) ? apiData : []).map(n => ({
      id: n.id,
      title: n.title,
      message: n.message,
      status: n.status || 'unread',
      icon_key: n.icon_key || 'icofont-notification'
    }))
    list = [...parsedApi]
  } catch (e) {}

  const local = localStorage.getItem('ncsms_local_notifications')
  if (local) {
    try {
      const parsedLocal = JSON.parse(local)
      list = [...parsedLocal, ...list]
    } catch (e) {}
  }
  
  notifications.value = list
}

async function markAllRead() {
  try {
    await cms.markAllNotificationsRead()
  } catch (e) {}
  const local = localStorage.getItem('ncsms_local_notifications')
  if (local) {
    try {
      const parsed = JSON.parse(local)
      parsed.forEach(n => { n.status = 'read' })
      localStorage.setItem('ncsms_local_notifications', JSON.stringify(parsed))
    } catch (e) {}
  }
  await fetchNotifications()
}

onMounted(() => {
  if (roles.includes('helpdesk') || roles.includes('role_helpdesk')) {
    fetchNotifications()
  }
})

const isHelpdesk = computed(() => {
  return roles.includes('helpdesk') || roles.includes('role_helpdesk')
})

const dashboardLink = computed(() => {
  return isHelpdesk.value ? '/portal/helpdesk' : '/dashboard'
})

const profileName = computed(() => {
  return `${user.value.first_name || ''} ${user.value.last_name || ''}`.trim() || 'Portal User'
})

const ordinaryNavigation = computed(() => {
  return [
    { id: 'dashboard', label: 'Dashboard', icon: 'icofont-dashboard-web' },
    { id: 'apply', label: 'Apply Now', icon: 'icofont-plus-circle' },
    { id: 'applications', label: 'My Applications', icon: 'icofont-file-document' },
    { id: 'activities', label: 'My Activities', icon: 'icofont-history' },
    { id: 'leave-approvals', label: 'Leave Approvals', icon: 'icofont-calendar' },
    { id: 'profile', label: 'My Profile', icon: 'icofont-user-alt-3' }
  ]
})

function goBack() {
  if (isHelpdesk.value) {
    router.push('/portal/helpdesk?tab=leave-approvals')
  } else {
    router.push('/dashboard?section=leave-approvals')
  }
}

function goHelpdeskTab(tab) {
  router.push(`/portal/helpdesk?tab=${tab}`)
}

function goUserSection(section) {
  router.push(`/dashboard?section=${section}`)
}

function logout() {
  localStorage.removeItem('ncsms_access_token')
  localStorage.removeItem('ncsms_user')
  router.push('/login')
}
</script>

<style scoped>
.navbar-bg {
  z-index: 1 !important;
}
.main-sidebar {
  z-index: 5 !important;
}
.main-content {
  min-height: 100vh;
}
</style>
