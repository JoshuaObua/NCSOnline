import { createRouter, createWebHistory } from 'vue-router'
import { installActivityAuditor } from '@/services/activityAudit.js'
import { isOrdinaryUser, storedPortalUser } from '@/utils/portalAuth.js'

const PortalLoginView = () => import('@/views/CMSLoginView.vue')
const PortalManagerView = () => import('@/views/WebsiteContentManagerView.vue')
const UserPortalView = () => import('@/views/UserPortalView.vue')
const ApplicationWizardView = () => import('@/views/ApplicationWizardView.vue')
const ApplicationDetailView = () => import('@/views/ApplicationDetailView.vue')

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', name: 'PortalLogin', component: PortalLoginView },
  { path: '/portal', name: 'PortalDashboard', component: PortalManagerView },
  { path: '/portal/namis/:resource/new', name: 'NamisRegistryCreate', component: PortalManagerView },
  { path: '/portal/applications/:id', name: 'AdminApplicationDetail', component: ApplicationDetailView },
  { path: '/dashboard', name: 'UserDashboard', component: UserPortalView },
  { path: '/dashboard/apply/:slug', name: 'ApplicationWizard', component: ApplicationWizardView },
  { path: '/dashboard/applications/:id', name: 'UserApplicationDetail', component: ApplicationDetailView },
  { path: '/my-portal', redirect: '/dashboard' },
  { path: '/apply', redirect: '/dashboard?section=apply' },
  { path: '/apply/:slug', redirect: to => ({ name: 'ApplicationWizard', params: { slug: to.params.slug } }) },
  { path: '/cms', redirect: '/portal' },
  { path: '/admin/:pathMatch(.*)*', redirect: '/portal' },
  { path: '/:pathMatch(.*)*', redirect: '/login' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

const pageTitles = {
  PortalLogin: 'Portal Login',
  PortalDashboard: 'Portal',
  UserDashboard: 'My Dashboard',
  ApplicationWizard: 'Application Form',
  AdminApplicationDetail: 'Application Review',
  NamisRegistryCreate: 'Add Sports Registry Record',
  UserApplicationDetail: 'Application Details',
}

router.beforeEach((to) => {
  const token = localStorage.getItem('ncsms_access_token')
  if (to.path === '/login') return token ? (isOrdinaryUser(storedPortalUser()) ? '/dashboard' : '/portal') : true
  if (!token) return '/login'
  const ordinary = isOrdinaryUser(storedPortalUser())
  if (to.path.startsWith('/portal') && ordinary) return '/dashboard'
  if (to.path.startsWith('/dashboard') && !ordinary) return '/portal'
  return true
})

router.afterEach((to) => {
  const title = pageTitles[to.name] || 'Portal Login'
  document.title = `${title} - NCS Uganda`
})

installActivityAuditor(router)

export default router
