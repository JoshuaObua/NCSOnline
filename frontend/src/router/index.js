import { expenseRoles } from '@/api/expenses'
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

// Executive & Departmental Dashboards
const AGSTechnicalDashboardView = () => import('@/views/executive/AssistantGeneralSecretaryTechnicalDashboard.vue')
const AGSAdminDashboardView = () => import('@/views/executive/AssistantGeneralSecretaryAdminDashboard.vue')
const GSAppraisalView = () => import('@/views/executive/GeneralSecretaryAppraisalView.vue')
const StoresInventoryDashboardView = () => import('@/views/stores/StoresInventoryDashboard.vue')
const FacilitiesManagementDashboardView = () => import('@/views/facilities/FacilitiesManagementDashboard.vue')
const LegalComplianceDashboardView = () => import('@/views/legal/LegalComplianceDashboard.vue')
const SportsMedicalDashboardView = () => import('@/views/medical/SportsMedicalDashboard.vue')
const FleetTransportDashboardView = () => import('@/views/fleet/FleetTransportDashboard.vue')

// Independent Departmental Report Views
const EngineeringReportView = () => import('@/views/reports/EngineeringReportView.vue')
const HumanResourcesReportView = () => import('@/views/reports/HumanResourcesReportView.vue')
const FinanceAccountsReportView = () => import('@/views/reports/FinanceAccountsReportView.vue')
const TechnicalSportsReportView = () => import('@/views/reports/TechnicalSportsReportView.vue')
const MedicalScienceReportView = () => import('@/views/reports/MedicalScienceReportView.vue')
const LegalLogisticsReportView = () => import('@/views/reports/LegalLogisticsReportView.vue')

// Reception & Front Desk Visitor Management
const ReceptionistVisitorView = () => import('@/views/reception/ReceptionistVisitorView.vue')
const VisitorInitiateView = () => import('@/views/reception/VisitorInitiateView.vue')
const ReceptionistLeaveApplyView = () => import('@/views/reception/ReceptionistLeaveApplyView.vue')
const ReceptionistLeaveStatusView = () => import('@/views/reception/ReceptionistLeaveStatusView.vue')
const ReceptionistReportsView = () => import('@/views/reception/ReceptionistReportsView.vue')

// IT Officer PPDA Form 5 Requisitions & Operations
const ITPPDAInitiateView = () => import('@/views/it/ITPPDAInitiateView.vue')
const ITPPDAListView = () => import('@/views/it/ITPPDAListView.vue')

const ExpenseRegisterView = () => import('@/views/expenses/ExpenseRegisterView.vue')
const ExpenseReportView = () => import('@/views/expenses/ExpenseReportView.vue')
const ExpenseEntryView = () => import('@/views/expenses/ExpenseEntryView.vue')
const ExpenseDetailView = () => import('@/views/expenses/ExpenseDetailView.vue')
const ExpenseCategoriesView = () => import('@/views/expenses/ExpenseCategoriesView.vue')
const FixedAssetActionView = () => import('@/views/FixedAssetActionView.vue')
const FixedAssetsView = () => import('@/views/FixedAssetsView.vue')
const FixedAssetValueAdjustmentsView = () => import('@/views/FixedAssetValueAdjustmentsView.vue')
const FixedAssetPivotEngineView = () => import('@/views/FixedAssetPivotEngineView.vue')

const routes = [
  { path: '/', redirect: to => (localStorage.getItem('ncsms_access_token') ? '/dashboard' : '/login') },
  { path: '/login', name: 'Login', component: LoginView, meta: { requiresAuth: false } },
  { path: '/register', name: 'Register', component: RegisterView, meta: { requiresAuth: false } },

  { path: '/expenses', name: 'Expenses', component: ExpenseRegisterView, meta: { requiresAuth: true, roles: expenseRoles } },
  { path: '/expenses/reports', name: 'ExpenseReport', component: ExpenseReportView, meta: { requiresAuth: true, roles: expenseRoles } },
  { path: '/expenses/new', name: 'ExpenseEntry', component: ExpenseEntryView, meta: { requiresAuth: true, roles: expenseRoles } },
  { path: '/expenses/:id/edit', name: 'ExpenseEdit', component: ExpenseEntryView, meta: { requiresAuth: true, roles: expenseRoles } },
  { path: '/expenses/categories', name: 'ExpenseCategories', component: ExpenseCategoriesView, meta: { requiresAuth: true, roles: ['super_admin','admin'] } },
  { path: '/expenses/:id', name: 'ExpenseDetail', component: ExpenseDetailView, meta: { requiresAuth: true, roles: expenseRoles } },
  { path: '/fixed-assets', name: 'FixedAssets', component: FixedAssetsView, meta: { requiresAuth: true } },
  { path: '/fixed-assets/value-adjustments', name: 'FixedAssetValueAdjustments', component: FixedAssetValueAdjustmentsView, meta: { requiresAuth: true } },
  { path: '/fixed-assets/pivot-engine', name: 'FixedAssetPivotEngine', component: FixedAssetPivotEngineView, meta: { requiresAuth: true } },
  { path: '/fixed-assets/new', name: 'FixedAssetNew', component: FixedAssetActionView, props: { action: 'new' }, meta: { requiresAuth: true } },
  { path: '/fixed-assets/depreciation', name: 'FixedAssetDepreciation', component: FixedAssetActionView, props: { action: 'depreciation' }, meta: { requiresAuth: true } },
  { path: '/fixed-assets/:id/revalue', name: 'FixedAssetRevalue', component: FixedAssetActionView, props: { action: 'revalue' }, meta: { requiresAuth: true } },
  { path: '/fixed-assets/:id/verify', name: 'FixedAssetVerify', component: FixedAssetActionView, props: { action: 'verify' }, meta: { requiresAuth: true } },
  { path: '/my-portal', name: 'ApplicantPortal', component: ApplicantPortalView, meta: { requiresAuth: true } },
  { path: '/my-portal/applications/new', name: 'NewOrganisationApplication', component: ApplicantLicensePortalView, meta: { requiresAuth: true } },
  { path: '/my-portal/forms/:slug', name: 'DynamicForm', component: DynamicPortalFormView, meta: { requiresAuth: true } },
  { path: '/accept-organisation-invite', name: 'AcceptOrganisationInvite', component: AcceptOrganisationInviteView, meta: { requiresAuth: true } },

  { path: '/dashboard', name: 'Dashboard', component: DashboardView, meta: { requiresAuth: true } },
  
  // Executive Leadership & Governance Routes
  { path: '/executive/ags-technical', name: 'AGSTechnicalDashboard', component: AGSTechnicalDashboardView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'ags_technical', 'general_secretary'] } },
  { path: '/executive/ags-admin', name: 'AGSAdminDashboard', component: AGSAdminDashboardView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'ags_admin', 'general_secretary'] } },
  { path: '/executive/appraisals', name: 'OnlineAppraisal', component: GSAppraisalView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'general_secretary'] } },

  // Independent Departmental Reports Routes
  { path: '/executive/reports/engineering', name: 'EngineeringReport', component: EngineeringReportView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'general_secretary', 'senior_engineer', 'ags_technical'] } },
  { path: '/executive/reports/human-resources', name: 'HumanResourcesReport', component: HumanResourcesReportView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'general_secretary', 'hr_officer', 'ags_admin'] } },
  { path: '/executive/reports/finance-accounts', name: 'FinanceAccountsReport', component: FinanceAccountsReportView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'general_secretary', 'finance_officer', 'ags_admin'] } },
  { path: '/executive/reports/technical-sports', name: 'TechnicalSportsReport', component: TechnicalSportsReportView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'general_secretary', 'ags_technical'] } },
  { path: '/executive/reports/medical-science', name: 'MedicalScienceReport', component: MedicalScienceReportView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'general_secretary', 'medical_officer', 'ags_technical'] } },
  { path: '/executive/reports/legal-logistics', name: 'LegalLogisticsReport', component: LegalLogisticsReportView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'general_secretary', 'legal_counsel', 'transport_officer'] } },

  // Receptionist Front Desk & Visitor Clearance Desk
  { path: '/reception/visitors', name: 'ReceptionVisitors', component: ReceptionistVisitorView, meta: { requiresAuth: true } },
  { path: '/reception/visitors/new', name: 'VisitorInitiate', component: VisitorInitiateView, meta: { requiresAuth: true } },
  { path: '/reception/leave/apply', name: 'ReceptionistLeaveApply', component: ReceptionistLeaveApplyView, meta: { requiresAuth: true } },
  { path: '/reception/leave/status', name: 'ReceptionistLeaveStatus', component: ReceptionistLeaveStatusView, meta: { requiresAuth: true } },
  { path: '/reception/reports', name: 'ReceptionistReports', component: ReceptionistReportsView, meta: { requiresAuth: true } },

  // IT Officer PPDA Form 5 Requisitions & Operations
  { path: '/it/ppda', redirect: '/it/ppda/status' },
  { path: '/it/ppda/new', name: 'ITPPDAInitiate', component: ITPPDAInitiateView, meta: { requiresAuth: true } },
  { path: '/it/ppda/status', name: 'ITPPDAList', component: ITPPDAListView, meta: { requiresAuth: true } },

  // Departmental Routes
  { path: '/stores/inventory', name: 'StoresInventory', component: StoresInventoryDashboardView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'stores_officer', 'general_secretary', 'ags_admin', 'accountant', 'senior_accountant', 'finance_department'] } },
  { path: '/facilities/venues', name: 'FacilitiesManagement', component: FacilitiesManagementDashboardView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'facilities_manager', 'general_secretary', 'ags_technical', 'ags_admin'] } },
  { path: '/legal/compliance', name: 'LegalCompliance', component: LegalComplianceDashboardView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'legal_counsel', 'general_secretary'] } },
  { path: '/medical/sports-science', name: 'SportsMedical', component: SportsMedicalDashboardView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'medical_officer', 'physiotherapist', 'general_secretary', 'ags_technical'] } },
  { path: '/fleet/transport', name: 'FleetTransport', component: FleetTransportDashboardView, meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'transport_officer', 'general_secretary', 'ags_admin'] } },

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
  FixedAssets: 'Manage Fixed Assets',
  FixedAssetValueAdjustments: 'Value Adjustments',
  FixedAssetPivotEngine: 'Dynamic Pivot Engine',
  FixedAssetNew: 'New Fixed Asset',
  FixedAssetRevalue: 'Asset Revaluation',
  FixedAssetVerify: 'Asset Verification',
  FixedAssetDepreciation: 'Monthly Depreciation',

  FixedAssetRevalue: 'Asset Revaluation',
  FixedAssetVerify: 'Asset Verification',
  FixedAssetDepreciation: 'Monthly Depreciation',
  Expenses: 'Expenses',
  ExpenseEntry: 'Record Expense',
  ExpenseDetail: 'Expense Details',
  ExpenseCategories: 'Expense Categories',
  Login: 'Login',
  Register: 'Register',
  Dashboard: 'Dashboard',
  ApplicantPortal: 'My Portal',
  NewOrganisationApplication: 'New Organisation Application',
  DynamicForm: 'Application Form',
  ReceptionVisitors: 'Reception & Visitor Clearance Desk',
  VisitorInitiate: 'Initiate Visitor Clearance Request',
  ReceptionistLeaveApply: 'Apply for Leave',
  ReceptionistLeaveStatus: 'Leave Status & History',
  ReceptionistReports: 'Front Desk & Application Reports',
  ITPPDAInitiate: 'PPDA Form 5 Requisition',
  ITPPDAList: 'PPDA Form 5 Status & Approvals',
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
    const parsed = JSON.parse(storedUser)
    return (parsed.roles || []).map(r => {
      if (typeof r === 'string') return r
      return r?.name || r?.role || r?.role_name || r?.slug || ''
    }).filter(Boolean)
  } catch {
    return []
  }
}

export default router
