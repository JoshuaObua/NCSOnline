import { createRouter, createWebHistory } from 'vue-router'
import { isOrdinaryUser, storedPortalUser } from '@/utils/portalAuth.js'

const PortalLoginView = () => import('@/views/CMSLoginView.vue')
const PortalManagerView = () => import('@/views/WebsiteContentManagerView.vue')
const UserPortalView = () => import('@/views/UserPortalView.vue')

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', name: 'PortalLogin', component: PortalLoginView },
  { path: '/portal', name: 'PortalDashboard', component: PortalManagerView },
  { path: '/dashboard', name: 'UserDashboard', component: UserPortalView },
  { path: '/my-portal', redirect: '/dashboard' },
  { path: '/apply', redirect: '/dashboard?section=apply' },
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
}

router.beforeEach((to) => {
  const token = localStorage.getItem('ncsms_access_token')
  if (to.path === '/login') return token ? (isOrdinaryUser(storedPortalUser()) ? '/dashboard' : '/portal') : true
  if (!token) return '/login'
  const ordinary = isOrdinaryUser(storedPortalUser())
  if (to.path === '/portal' && ordinary) return '/dashboard'
  if (to.path === '/dashboard' && !ordinary) return '/portal'
  return true
})

router.afterEach((to) => {
  const title = pageTitles[to.name] || 'Portal Login'
  document.title = `${title} - NCS Uganda`
})

export default router
