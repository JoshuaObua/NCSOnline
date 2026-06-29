import { createRouter, createWebHistory } from 'vue-router'

// ── Layouts ───────────────────────────────────────────────────────
const PublicLayout = () => import('@/layouts/PublicLayout.vue')

// ── Public views ──────────────────────────────────────────────────
const HomeView = () => import('@/views/public/HomeView.vue')
const BlogView = () => import('@/views/public/BlogView.vue')
const BlogPostView = () => import('@/views/public/BlogPostView.vue')
const PageView = () => import('@/views/public/PageView.vue')
const EventsView = () => import('@/views/public/EventsView.vue')
const EventDetailView = () => import('@/views/public/EventDetailView.vue')
const CareersView = () => import('@/views/public/CareersView.vue')
const CareerDetailView = () => import('@/views/public/CareerDetailView.vue')
const ProjectsView = () => import('@/views/public/ProjectsView.vue')
const CaseStudiesView = () => import('@/views/public/CaseStudiesView.vue')
const LicensePortalView = () => import('@/views/public/LicensePortalView.vue')
const ResourceCentreView = () => import('@/views/public/ResourceCentreView.vue')
const FacilitiesView = () => import('@/views/public/FacilitiesView.vue')
const FacilityDetailView = () => import('@/views/public/FacilityDetailView.vue')
const AssociationsView = () => import('@/views/public/AssociationsView.vue')
const InvestView = () => import('@/views/public/InvestView.vue')
const FAQsView = () => import('@/views/public/FAQsView.vue')
const ContactUsView = () => import('@/views/public/ContactUsView.vue')
const TeamView = () => import('@/views/public/TeamView.vue')

// ── Auth / applicant views ────────────────────────────────────────
const RegisterView = () => import('@/views/RegisterView.vue')
const ApplicantPortalView = () => import('@/views/ApplicantPortalView.vue')
const AcceptOrganisationInviteView = () => import('@/views/AcceptOrganisationInviteView.vue')

// ── Admin / auth views ────────────────────────────────────────────
const LoginView = () => import('@/views/LoginView.vue')
const DashboardView = () => import('@/views/DashboardEntryView.vue')
const UsersView = () => import('@/views/UsersView.vue')
const ApplicationsView = () => import('@/views/ApplicationsView.vue')
const AuditLogsView = () => import('@/views/AuditLogsView.vue')
const RolesView = () => import('@/views/RolesView.vue')
const ProfileView = () => import('@/views/ProfileView.vue')
const CMSView = () => import('@/views/CMSView.vue')
const GovernanceDashboardView = () => import('@/views/GovernanceDashboardView.vue')
const ReportingWorkspaceView = () => import('@/views/ReportingWorkspaceView.vue')
const InsightsDashboardView = () => import('@/views/InsightsDashboardView.vue')
const NSMISRegistryView = () => import('@/views/NSMISRegistryView.vue')
const MaintenanceView = () => import('@/views/MaintenanceView.vue')
const BackupsView = () => import('@/views/BackupsView.vue')
const AdminFormsView = () => import('@/views/AdminFormsView.vue')
const AdminFormSubmissionsView = () => import('@/views/AdminFormSubmissionsView.vue')
const DynamicFormView = () => import('@/views/public/DynamicFormView.vue')
const SecuritySettingsView = () => import('@/views/SecuritySettingsView.vue')
const MyActivitiesView = () => import('@/views/MyActivitiesView.vue')
const SmartUpdatesView = () => import('@/views/SmartUpdatesView.vue')
const StorageSettingsView = () => import('@/views/StorageSettingsView.vue')

const routes = [
  // ── Public website (uses PublicLayout) ────────────────────────
  {
    path: '/',
    component: PublicLayout,
    children: [
      { path: '', name: 'Home', component: HomeView, meta: { requiresAuth: false } },
      { path: 'news', name: 'News', component: BlogView, meta: { requiresAuth: false } },
      { path: 'news/:slug', name: 'NewsPost', component: BlogPostView, meta: { requiresAuth: false } },
      { path: 'pages/:slug', name: 'Page', component: PageView, meta: { requiresAuth: false } },
      { path: 'events', name: 'Events', component: EventsView, meta: { requiresAuth: false } },
      { path: 'events/:slug', name: 'EventDetail', component: EventDetailView, meta: { requiresAuth: false } },
      { path: 'careers', name: 'Careers', component: CareersView, meta: { requiresAuth: false } },
      { path: 'careers/:id', name: 'CareerDetail', component: CareerDetailView, meta: { requiresAuth: false } },
      { path: 'projects', name: 'Projects', component: ProjectsView, meta: { requiresAuth: false } },
      { path: 'case-studies', name: 'CaseStudies', component: CaseStudiesView, meta: { requiresAuth: false } },
      { path: 'apply', redirect: () => localStorage.getItem('ncsms_access_token') ? '/my-portal/applications/new' : '/login?redirect=/my-portal/applications/new' },
      { path: 'resource-centre', name: 'ResourceCentre', component: ResourceCentreView, meta: { requiresAuth: false } },
      { path: 'resources', redirect: '/resource-centre' },
      { path: 'facilities', name: 'Facilities', component: FacilitiesView, meta: { requiresAuth: false } },
      { path: 'facilities/:slug', name: 'FacilityDetail', component: FacilityDetailView, meta: { requiresAuth: false } },
      { path: 'associations', name: 'Associations', component: AssociationsView, meta: { requiresAuth: false } },
      { path: 'invest', name: 'Invest', component: InvestView, meta: { requiresAuth: false } },
      { path: 'faqs', name: 'FAQs', component: FAQsView, meta: { requiresAuth: false } },
      { path: 'contact-us', name: 'ContactUs', component: ContactUsView, meta: { requiresAuth: false } },
      { path: 'team', name: 'Team', component: TeamView, meta: { requiresAuth: false } },
    ]
  },

  // ── Auth (rendered inside PublicLayout so they share the header/footer) ──
  {
    path: '/',
    component: PublicLayout,
    meta: { requiresAuth: false },
    children: [
      { path: 'login',    name: 'Login',    component: LoginView },
      { path: 'register', name: 'Register', component: RegisterView },
    ],
  },

  // ── Applicant self-service portal (uses PublicLayout) ─────────
  {
    path: '/',
    component: PublicLayout,
    children: [
      {
        path: 'my-portal',
        name: 'ApplicantPortal',
        component: ApplicantPortalView,
        meta: { requiresAuth: true }
      },
      {
        path: 'my-portal/applications/new',
        name: 'NewOrganisationApplication',
        component: LicensePortalView,
        meta: { requiresAuth: true }
      },
      {
        path: 'my-portal/forms/:slug',
        name: 'DynamicForm',
        component: DynamicFormView,
        meta: { requiresAuth: true }
      },
      {
        path: 'accept-organisation-invite',
        name: 'AcceptOrganisationInvite',
        component: AcceptOrganisationInviteView,
        meta: { requiresAuth: true }
      }
    ]
  },

  // ── Admin portal (existing admin views, no layout wrapper here — they have LayoutDefault) ──
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { requiresAuth: true }
  },
  {
    path: '/nsmis/governance',
    name: 'GovernanceDashboard',
    component: GovernanceDashboardView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'ncs_general_secretary', 'general_secretary', 'technical_department', 'federation_president', 'federation_general_secretary', 'auditor'] }
  },
  {
    path: '/nsmis/reports',
    name: 'FederationReports',
    component: ReportingWorkspaceView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'ncs_general_secretary', 'general_secretary', 'technical_department', 'finance_department', 'federation_president', 'federation_general_secretary', 'auditor'] }
  },
  {
    path: '/nsmis/insights/:dashboard(athletes|performance|finance|talent)',
    name: 'NSMISInsights',
    component: InsightsDashboardView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'ncs_general_secretary', 'general_secretary', 'technical_department', 'finance_department'] }
  },
  {
    path: '/nsmis/data/:resource(federation-officers|athletes|competitions|medals|coaches|technical-officials|talent|safeguarding-aggregates|disbursements|accountabilities|equipment)',
    name: 'NSMISRegistry', component: NSMISRegistryView,
    meta: { requiresAuth: true, roles: ['super_admin','admin','ncs_general_secretary','general_secretary','technical_department','finance_department','federation_president','federation_general_secretary','safeguarding_officer'] }
  },
  {
    path: '/users',
    name: 'Users',
    component: UsersView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin'] }
  },
  {
    path: '/applications',
    name: 'Applications',
    component: ApplicationsView,
    meta: { requiresAuth: true }
  },
  {
    path: '/audit-logs',
    name: 'AuditLogs',
    component: AuditLogsView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin'] }
  },
  {
    path: '/maintenance', name: 'Maintenance', component: MaintenanceView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin'] }
  },
  {
    path: '/maintenance/backups', name: 'Backups', component: BackupsView,
    meta: { requiresAuth: true, roles: ['super_admin'] }
  },
  {
    path: '/maintenance/updates', name: 'SmartUpdates', component: SmartUpdatesView,
    meta: { requiresAuth: true, roles: ['super_admin'] }
  },
  {
    path: '/settings/storage', name: 'StorageSettings', component: StorageSettingsView,
    meta: { requiresAuth: true, roles: ['super_admin'] }
  },
  {
    path: '/roles',
    name: 'Roles',
    component: RolesView,
    meta: { requiresAuth: true, roles: ['super_admin'] }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: ProfileView,
    meta: { requiresAuth: true }
  },
  {
    path: '/me/activities',
    name: 'MyActivities',
    component: MyActivitiesView,
    meta: { requiresAuth: true }
  },
  {
    path: '/me/security',
    name: 'SecuritySettings',
    component: SecuritySettingsView,
    meta: { requiresAuth: true }
  },
  {
    path: '/cms',
    name: 'CMS',
    component: CMSView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'content_manager'] }
  },
  {
    path: '/admin/forms',
    name: 'AdminForms',
    component: AdminFormsView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'general_secretary'] }
  },
  {
    path: '/admin/forms/submissions',
    name: 'AdminFormSubmissions',
    component: AdminFormSubmissionsView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'general_secretary'] }
  },

  // Catch-all — redirect to home
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

const publicPageTitles = {
  Home: 'National Council of Sports Uganda',
  News: 'News & Updates', NewsPost: 'News Article', Page: 'Information',
  Events: 'Events', EventDetail: 'Event Details', Careers: 'Careers', CareerDetail: 'Career Opportunity',
  Projects: 'Projects', CaseStudies: 'Case Studies', NewOrganisationApplication: 'New Organisation Application',
  ResourceCentre: 'Resource Centre', Facilities: 'Sports Facilities', FacilityDetail: 'Facility Details',
  Associations: 'Sports Associations', Invest: 'Invest with NCS', FAQs: 'Frequently Asked Questions',
  ContactUs: 'Contact Us', Team: 'NCS Membership', ApplicantPortal: 'My Portal'
}

router.afterEach((to) => {
  const title = publicPageTitles[to.name]
  if (title) document.title = `${title} — NCS Uganda`
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('ncsms_access_token')
  const isAuthenticated = !!token

  // Public routes — always accessible
  if (to.meta.requiresAuth === false) {
    // If authenticated and trying to reach login/register, redirect to correct portal
    if (isAuthenticated && (to.name === 'Login' || to.name === 'Register')) {
      const storedUser = localStorage.getItem('ncsms_user')
      let roles = []
      if (storedUser) {
        try { roles = (JSON.parse(storedUser).roles || []).map(r => typeof r === 'string' ? r : r.name) } catch { roles = [] }
      }
      const isApplicant = !roles.length || roles.some(r => r === 'applicant' || r === 'user')
      return next(isApplicant ? '/my-portal' : '/dashboard')
    }
    return next()
  }

  // Protected route — require authentication
  if (!isAuthenticated) {
    return next({ path: '/login', query: { redirect: to.fullPath } })
  }

  // Resolve user roles once for both checks below
  const storedUser = localStorage.getItem('ncsms_user')
  let userRoles = []
  if (storedUser) {
    try {
      userRoles = (JSON.parse(storedUser).roles || []).map(r => (typeof r === 'string' ? r : r.name))
    } catch {
      userRoles = []
    }
  }
  const isApplicantUser = !userRoles.length || userRoles.some(r => r === 'applicant' || r === 'user')

  // Applicant/user role — only allowed on portal and profile; block all staff routes
  const staffRoutes = ['Dashboard', 'Users', 'Applications', 'AuditLogs', 'Roles', 'CMS', 'Maintenance', 'Backups', 'SmartUpdates', 'StorageSettings', 'AdminForms', 'AdminFormSubmissions']
  const hasActiveOrganisation = !!localStorage.getItem('ncsms_active_organisation')
  if (isApplicantUser && staffRoutes.includes(to.name) && !(to.name === 'Dashboard' && hasActiveOrganisation)) {
    return next('/my-portal')
  }

  // Role-based access check for routes that declare required roles
  if (to.meta.roles && to.meta.roles.length > 0) {
    const hasRole = to.meta.roles.some(role => userRoles.includes(role))
    if (!hasRole) {
      return next('/dashboard')
    }
  }

  next()
})

export default router
