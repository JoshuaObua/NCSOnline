import { createRouter, createWebHistory } from 'vue-router'

const PublicLayout = () => import('@/layouts/PublicLayout.vue')

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
const ResourceCentreView = () => import('@/views/public/ResourceCentreView.vue')
const FacilitiesView = () => import('@/views/public/FacilitiesView.vue')
const FacilityDetailView = () => import('@/views/public/FacilityDetailView.vue')
const AssociationsView = () => import('@/views/public/AssociationsView.vue')
const InvestView = () => import('@/views/public/InvestView.vue')
const FAQsView = () => import('@/views/public/FAQsView.vue')
const ContactUsView = () => import('@/views/public/ContactUsView.vue')
const TeamView = () => import('@/views/public/TeamView.vue')
const CMSLoginView = () => import('@/views/CMSLoginView.vue')
const WebsiteContentManagerView = () => import('@/views/WebsiteContentManagerView.vue')

const routes = [
  {
    path: '/',
    component: PublicLayout,
    children: [
      { path: '', name: 'Home', component: HomeView },
      { path: 'news', name: 'News', component: BlogView },
      { path: 'news/:slug', name: 'NewsPost', component: BlogPostView },
      { path: 'pages/:slug', name: 'Page', component: PageView },
      { path: 'events', name: 'Events', component: EventsView },
      { path: 'events/:slug', name: 'EventDetail', component: EventDetailView },
      { path: 'careers', name: 'Careers', component: CareersView },
      { path: 'careers/:id', name: 'CareerDetail', component: CareerDetailView },
      { path: 'projects', name: 'Projects', component: ProjectsView },
      { path: 'case-studies', name: 'CaseStudies', component: CaseStudiesView },
      { path: 'resource-centre', name: 'ResourceCentre', component: ResourceCentreView },
      { path: 'resources', redirect: '/resource-centre' },
      { path: 'facilities', name: 'Facilities', component: FacilitiesView },
      { path: 'facilities/:slug', name: 'FacilityDetail', component: FacilityDetailView },
      { path: 'associations', name: 'Associations', component: AssociationsView },
      { path: 'invest', name: 'Invest', component: InvestView },
      { path: 'faqs', name: 'FAQs', component: FAQsView },
      { path: 'contact-us', name: 'ContactUs', component: ContactUsView },
      { path: 'team', name: 'Team', component: TeamView },
    ],
  },
  { path: '/login', name: 'CMSLogin', component: CMSLoginView },
  { path: '/cms', name: 'WebsiteCMS', component: WebsiteContentManagerView },
  { path: '/apply', redirect: '/contact-us' },
  { path: '/register', redirect: '/login' },
  { path: '/my-portal/:pathMatch(.*)*', redirect: '/contact-us' },
  { path: '/dashboard/:pathMatch(.*)*', redirect: '/' },
  { path: '/admin/:pathMatch(.*)*', redirect: '/cms' },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

const publicPageTitles = {
  Home: 'National Council of Sports Uganda',
  News: 'News & Updates',
  NewsPost: 'News Article',
  Page: 'Information',
  Events: 'Events',
  EventDetail: 'Event Details',
  Careers: 'Careers',
  CareerDetail: 'Career Opportunity',
  Projects: 'Projects',
  CaseStudies: 'Case Studies',
  ResourceCentre: 'Resource Centre',
  Facilities: 'Sports Facilities',
  FacilityDetail: 'Facility Details',
  Associations: 'Sports Associations',
  Invest: 'Invest with NCS',
  FAQs: 'Frequently Asked Questions',
  ContactUs: 'Contact Us',
  Team: 'NCS Membership',
  CMSLogin: 'Content Manager Login',
  WebsiteCMS: 'Content Manager',
}

router.afterEach((to) => {
  const title = publicPageTitles[to.name]
  if (title) document.title = `${title} - NCS Uganda`
})

export default router
