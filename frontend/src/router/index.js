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
const DocumentsView = () => import('@/views/public/DocumentsView.vue')
const InvestView = () => import('@/views/public/InvestView.vue')
const FAQsView = () => import('@/views/public/FAQsView.vue')
const ContactUsView = () => import('@/views/public/ContactUsView.vue')
const TeamView = () => import('@/views/public/TeamView.vue')
const CMSLoginView = () => import('@/views/CMSLoginView.vue')
const WebsiteContentManagerView = () => import('@/views/WebsiteContentManagerView.vue')
const AccountProfileView = () => import('@/views/account/AccountProfileView.vue')
const AccountSettingsView = () => import('@/views/account/AccountSettingsView.vue')
const AccountActivitiesView = () => import('@/views/account/AccountActivitiesView.vue')

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
      {
        path: 'sports-rules',
        name: 'SportsRules',
        component: DocumentsView,
        props: {
          docType: 'sports_rule',
          heroLead: 'Sports',
          heroAccent: 'Rules',
          subtitle: 'Official rules and regulations governing sports disciplines recognised by the National Council of Sports Uganda.',
          icon: 'icofont-read-book',
          emptyText: 'No sports rules published yet.',
        },
      },
      {
        path: 'press-releases',
        name: 'PressReleases',
        component: DocumentsView,
        props: {
          docType: 'press_release',
          heroLead: 'Press',
          heroAccent: 'Releases',
          subtitle: 'Official statements, announcements and media coverage from the National Council of Sports Uganda.',
          icon: 'icofont-newspaper',
          emptyText: 'No press releases published yet.',
        },
      },
      {
        path: 'reports',
        name: 'Reports',
        component: DocumentsView,
        props: {
          docType: 'report',
          heroLead: 'NCS',
          heroAccent: 'Reports',
          subtitle: 'Annual reports, audits and other official publications from the National Council of Sports Uganda.',
          icon: 'icofont-file-pdf',
          emptyText: 'No reports published yet.',
        },
      },
      {
        path: 'speeches',
        name: 'Speeches',
        component: DocumentsView,
        props: {
          docType: 'speech',
          heroLead: 'NCS',
          heroAccent: 'Speeches',
          subtitle: 'Official speeches and addresses delivered by National Council of Sports Uganda leadership.',
          icon: 'icofont-speech-comments',
          emptyText: 'No speeches published yet.',
        },
      },
      { path: 'invest', name: 'Invest', component: InvestView },
      { path: 'faqs', name: 'FAQs', component: FAQsView },
      { path: 'contact-us', name: 'ContactUs', component: ContactUsView },
      { path: 'team', name: 'Team', component: TeamView },
      { path: 'account/profile', name: 'AccountProfile', component: AccountProfileView },
      { path: 'account/settings', name: 'AccountSettings', component: AccountSettingsView },
      { path: 'account/activities', name: 'AccountActivities', component: AccountActivitiesView },
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
  SportsRules: 'Sports Rules',
  PressReleases: 'Press Releases',
  Reports: 'NCS Reports',
  Speeches: 'NCS Speeches',
  Invest: 'Invest with NCS',
  FAQs: 'Frequently Asked Questions',
  ContactUs: 'Contact Us',
  Team: 'NCS Membership',
  AccountProfile: 'Profile Overview',
  AccountSettings: 'Settings & Security',
  AccountActivities: 'My Audit Activities',
  CMSLogin: 'Content Manager Login',
  WebsiteCMS: 'Content Manager',
}

router.afterEach((to) => {
  const title = publicPageTitles[to.name]
  if (title) document.title = `${title} - NCS Uganda`
})

export default router
