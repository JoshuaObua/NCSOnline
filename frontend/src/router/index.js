import { createRouter, createWebHistory } from 'vue-router'

// ── Layouts ───────────────────────────────────────────────────────
const PublicLayout = () => import('@/layouts/PublicLayout.vue')

// ── Public views ──────────────────────────────────────────────────
const HomeView = () => import('@/views/public/HomeView.vue')
const BlogView = () => import('@/views/public/BlogView.vue')
const BlogPostView = () => import('@/views/public/BlogPostView.vue')
const EventsView = () => import('@/views/public/EventsView.vue')
const EventDetailView = () => import('@/views/public/EventDetailView.vue')
const CareersView = () => import('@/views/public/CareersView.vue')
const CareerDetailView = () => import('@/views/public/CareerDetailView.vue')
const ProjectsView = () => import('@/views/public/ProjectsView.vue')
const CaseStudiesView = () => import('@/views/public/CaseStudiesView.vue')
const LicensePortalView = () => import('@/views/public/LicensePortalView.vue')

// ── Admin / auth views ────────────────────────────────────────────
const LoginView = () => import('@/views/LoginView.vue')
const DashboardView = () => import('@/views/DashboardView.vue')
const UsersView = () => import('@/views/UsersView.vue')
const ApplicationsView = () => import('@/views/ApplicationsView.vue')
const AuditLogsView = () => import('@/views/AuditLogsView.vue')
const RolesView = () => import('@/views/RolesView.vue')
const ProfileView = () => import('@/views/ProfileView.vue')
const CMSView = () => import('@/views/CMSView.vue')

const routes = [
  // ── Public website (uses PublicLayout) ────────────────────────
  {
    path: '/',
    component: PublicLayout,
    children: [
      { path: '', name: 'Home', component: HomeView, meta: { requiresAuth: false } },
      { path: 'news', name: 'News', component: BlogView, meta: { requiresAuth: false } },
      { path: 'news/:slug', name: 'NewsPost', component: BlogPostView, meta: { requiresAuth: false } },
      { path: 'events', name: 'Events', component: EventsView, meta: { requiresAuth: false } },
      { path: 'events/:slug', name: 'EventDetail', component: EventDetailView, meta: { requiresAuth: false } },
      { path: 'careers', name: 'Careers', component: CareersView, meta: { requiresAuth: false } },
      { path: 'careers/:id', name: 'CareerDetail', component: CareerDetailView, meta: { requiresAuth: false } },
      { path: 'projects', name: 'Projects', component: ProjectsView, meta: { requiresAuth: false } },
      { path: 'case-studies', name: 'CaseStudies', component: CaseStudiesView, meta: { requiresAuth: false } },
      { path: 'apply', name: 'Apply', component: LicensePortalView, meta: { requiresAuth: false } },
    ]
  },

  // ── Auth ──────────────────────────────────────────────────────
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { requiresAuth: false }
  },

  // ── Admin portal (existing admin views, no layout wrapper here — they have LayoutDefault) ──
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { requiresAuth: true }
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
    path: '/cms',
    name: 'CMS',
    component: CMSView,
    meta: { requiresAuth: true, roles: ['super_admin', 'admin', 'content_manager'] }
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

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('ncsms_access_token')
  const isAuthenticated = !!token

  // Public routes — always accessible
  if (to.meta.requiresAuth === false) {
    // If authenticated and trying to go to login, redirect to dashboard
    if (isAuthenticated && to.name === 'Login') {
      return next('/dashboard')
    }
    return next()
  }

  // Protected route — require authentication
  if (!isAuthenticated) {
    return next('/login')
  }

  // Role-based access check
  if (to.meta.roles && to.meta.roles.length > 0) {
    const storedUser = localStorage.getItem('ncsms_user')
    let userRoles = []
    if (storedUser) {
      try {
        const user = JSON.parse(storedUser)
        userRoles = (user.roles || []).map(r => (typeof r === 'string' ? r : r.name))
      } catch {
        userRoles = []
      }
    }
    const hasRole = to.meta.roles.some(role => userRoles.includes(role))
    if (!hasRole) {
      return next('/dashboard')
    }
  }

  next()
})

export default router
