import { createRouter, createWebHistory } from 'vue-router'

const PortalLoginView = () => import('@/views/CMSLoginView.vue')
const PortalManagerView = () => import('@/views/WebsiteContentManagerView.vue')

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', name: 'PortalLogin', component: PortalLoginView },
  { path: '/portal', name: 'PortalDashboard', component: PortalManagerView },
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
}

router.afterEach((to) => {
  const title = pageTitles[to.name] || 'Portal Login'
  document.title = `${title} - NCS Uganda`
})

export default router
