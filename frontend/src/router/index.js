import { createRouter, createWebHistory } from 'vue-router'

const LoginView = () => import('@/views/LoginView.vue')
const RegisterView = () => import('@/views/RegisterView.vue')
const DashboardView = () => import('@/views/DashboardEntryView.vue')
const UsersView = () => import('@/views/UsersView.vue')
const ApplicationsView = () => import('@/views/ApplicationsView.vue')
const AuditLogsView = () => import('@/views/AuditLogsView.vue')
const RolesView = () => import('@/views/RolesView.vue')
const ProfileView = () => import('@/views/ProfileView.vue')
const CMSView = () => import('@/views/CMSView.vue')
const PageBuilderView = () => import('@/views/PageBuilderView.vue')
const GovernanceDashboardView = () => import('@/views/GovernanceDashboardView.vue')
const ReportingWorkspaceView = () => import('@/views/ReportingWorkspaceView.vue')
const InsightsDashboardView = () => import('@/views/InsightsDashboardView.vue')
const NSMISRegistryView = () => import('@/views/NSMISRegistryView.vue')
const MaintenanceView = () => import('@/views/MaintenanceView.vue')
const InfrastructureCommandCenterView = () => import('@/views/InfrastructureCommandCenterView.vue')
const BackupsView = () => import('@/views/BackupsView.vue')
const AdminFormsView = () => import('@/views/AdminFormsView.vue')
const AdminFormSubmissionsView = () => import('@/views/AdminFormSubmissionsView.vue')
const SecuritySettingsView = () => import('@/views/SecuritySettingsView.vue')
const MyActivitiesView = () => import('@/views/MyActivitiesView.vue')
const SmartUpdatesView = () => import('@/views/SmartUpdatesView.vue')
const StorageSettingsView = () => import('@/views/StorageSettingsView.vue')
const ApplicantPortalView = () => import('@/views/ApplicantPortalView.vue')
const ApplicantLicensePortalView = () => import('@/views/ApplicantLicensePortalView.vue')
const DynamicPortalFormView = () => import('@/views/DynamicPortalFormView.vue')
const AcceptOrganisationInviteView = () => import('@/views/AcceptOrganisationInviteView.vue')

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', name: 'Login', component: LoginView, meta: { requiresAuth: false } },
  { path: '/register', name: 'Register', component: RegisterView, meta: { requiresAuth: false } },

  { path: '/my-portal', name: 'ApplicantPortal', component: ApplicantPortalView, meta: { requiresAuth: true } },
  { path: '/my-portal/applications/new', name: 'NewOrganisationApplication', component: ApplicantLicensePortalView, meta: { requiresAuth: true } },
  { path: '/my-portal/forms/:slug', name: 'DynamicForm', component: DynamicPortalFormView, meta: { requiresAuth: true } },
  { path: '/accept-organisation-invite', name: 'AcceptOrganisationInvite', component: AcceptOrganisationInviteView, meta: { requiresAuth: true } },

  { path: '/dashboard', name: 'Dashboard', component: DashboardView, meta: { requiresAuth: true } },
  {
    path: '/nsmis/governance',
    name: 'GovernanceDashboard',
    component: GovernanceDashboardView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'ncs_general_secretary', 'general_secretary', 'technical_department', 'federation_president', 'federation_general_secretary', 'auditor'] },
  },
  {
    path: '/nsmis/reports',
    name: 'FederationReports',
    component: ReportingWorkspaceView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'ncs_general_secretary', 'general_secretary', 'technical_department', 'finance_department', 'federation_president', 'federation_general_secretary', 'auditor'] },
  },
  {
    path: '/nsmis/insights/:dashboard(athletes|performance|finance|talent)',
    name: 'NSMISInsights',
    component: InsightsDashboardView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'ncs_general_secretary', 'general_secretary', 'technical_department', 'finance_department'] },
  },
  {
    path: '/nsmis/data/:resource(federation-officers|athletes|competitions|medals|coaches|technical-officials|talent|safeguarding-aggregates|disbursements|accountabilities|equipment)',
    name: 'NSMISRegistry',
    component: NSMISRegistryView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'ncs_general_secretary', 'general_secretary', 'technical_department', 'finance_department', 'federation_president', 'federation_general_secretary', 'safeguarding_officer'] },
  },
  { path: '/users', name: 'Users', component: UsersView, meta: { requiresAuth: true, roles: ['super_admin', 'admin'] } },
  { path: '/applications', name: 'Applications', component: ApplicationsView, meta: { requiresAuth: true } },
  { path: '/audit-logs', name: 'AuditLogs', component: AuditLogsView, meta: { requiresAuth: true, roles: ['super_admin', 'admin'] } },
  { path: '/maintenance', name: 'Maintenance', component: MaintenanceView, meta: { requiresAuth: true, roles: ['super_admin', 'admin'] } },
  { path: '/maintenance/command-center', name: 'InfrastructureCommandCenter', component: InfrastructureCommandCenterView, meta: { requiresAuth: true, roles: ['super_admin', 'admin'] } },
  { path: '/maintenance/backups', name: 'Backups', component: BackupsView, meta: { requiresAuth: true, roles: ['super_admin'] } },
  { path: '/maintenance/updates', name: 'SmartUpdates', component: SmartUpdatesView, meta: { requiresAuth: true, roles: ['super_admin'] } },
  { path: '/settings/storage', name: 'StorageSettings', component: StorageSettingsView, meta: { requiresAuth: true, roles: ['super_admin'] } },
  { path: '/roles', name: 'Roles', component: RolesView, meta: { requiresAuth: true, roles: ['super_admin'] } },
  { path: '/profile', name: 'Profile', component: ProfileView, meta: { requiresAuth: true } },
  { path: '/me/activities', name: 'MyActivities', component: MyActivitiesView, meta: { requiresAuth: true } },
  { path: '/me/security', name: 'SecuritySettings', component: SecuritySettingsView, meta: { requiresAuth: true } },
  { path: '/cms', name: 'CMS', component: CMSView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'content_manager'] } },
  { path: '/cms/page-builder', name: 'PageBuilder', component: PageBuilderView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'content_manager'] } },
  { path: '/admin/forms', name: 'AdminForms', component: AdminFormsView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'general_secretary'] } },
  { path: '/admin/forms/submissions', name: 'AdminFormSubmissions', component: AdminFormSubmissionsView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'general_secretary'] } },
  { path: '/admin', redirect: '/dashboard' },

  { path: '/apply', redirect: '/my-portal/applications/new' },
  { path: '/:pathMatch(.*)*', redirect: '/dashboard' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

const pageTitles = {
  Login: 'Login',
  Register: 'Register',
  Dashboard: 'Dashboard',
  ApplicantPortal: 'My Portal',
  NewOrganisationApplication: 'New Organisation Application',
  DynamicForm: 'Application Form',
  CMS: 'Content Manager',
  PageBuilder: 'Page Builder',
}

router.afterEach((to) => {
  const title = pageTitles[to.name]
  document.title = `${title || 'Intranet'} - NCS Uganda`
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('ncsms_access_token')
  const isAuthenticated = !!token

  if (to.meta.requiresAuth === false) {
    if (isAuthenticated && (to.name === 'Login' || to.name === 'Register')) {
      const roles = readUserRoles()
      const isApplicant = !roles.length || roles.some(r => r === 'applicant' || r === 'user')
      return next(isApplicant ? '/my-portal' : '/dashboard')
    }
    return next()
  }

  if (!isAuthenticated) {
    return next({ path: '/login', query: { redirect: to.fullPath } })
  }

  const userRoles = readUserRoles()
  const isApplicantUser = !userRoles.length || userRoles.some(r => r === 'applicant' || r === 'user')
  const staffRoutes = ['Dashboard', 'Users', 'Applications', 'AuditLogs', 'Roles', 'CMS', 'PageBuilder', 'Maintenance', 'InfrastructureCommandCenter', 'Backups', 'SmartUpdates', 'StorageSettings', 'AdminForms', 'AdminFormSubmissions']
  const hasActiveOrganisation = !!localStorage.getItem('ncsms_active_organisation')

  if (isApplicantUser && staffRoutes.includes(to.name) && !(to.name === 'Dashboard' && hasActiveOrganisation)) {
    return next('/my-portal')
  }

  if (to.meta.roles && to.meta.roles.length > 0) {
    const hasRole = to.meta.roles.some(role => userRoles.includes(role))
    if (!hasRole) return next('/dashboard')
  }

  next()
})

function readUserRoles() {
  const storedUser = localStorage.getItem('ncsms_user')
  if (!storedUser) return []
  try {
    return (JSON.parse(storedUser).roles || []).map(r => (typeof r === 'string' ? r : r.name))
  } catch {
    return []
  }
}

export default router
