<template>
  <div class="public-site min-h-screen flex flex-col">

    <!-- ── Navbar ──────────────────────────────────────────── -->
    <header
      :class="[
        'sticky top-0 z-50 transition-all duration-300',
        scrolled ? 'bg-white shadow-md' : 'bg-cream'
      ]"
    >
      <div class="max-w-screen-xl px-6 mx-auto">
        <div class="flex items-center justify-between h-20">

          <!-- Logo -->
          <router-link to="/" class="flex items-center flex-shrink-0">
            <div class="relative flex-shrink-0">
              <div class="w-16 h-16 bg-white rounded-xl flex items-center justify-center overflow-hidden shadow-sm border border-gray-100 z-10 relative">
                <img src="/main-logo.png" alt="NCS Logo" class="w-14 h-14 object-contain" />
              </div>
              <!-- Diamond accent -->
              <svg class="absolute -top-1.5 -left-1.5 w-8 h-8 z-0 opacity-30" viewBox="0 0 79 79" fill="none">
                <path d="M35.26 2.24C37.6-.1 41.4-.1 43.74 2.24L76.76 35.26C79.1 37.6 79.1 41.4 76.76 43.74L43.74 76.76C41.4 79.1 37.6 79.1 35.26 76.76L2.24 43.74C-.1 41.4-.1 37.6 2.24 35.26L35.26 2.24Z" fill="#2F327D"/>
              </svg>
            </div>
          </router-link>

          <!-- Desktop nav -->
          <nav class="hidden lg:flex items-center gap-1 text-sm font-medium">
            <template v-for="item in menuItems" :key="item.id || item.label">
              <div v-if="(item.children?.length) || item.mega" class="relative group">
                <button class="flex items-center gap-1 px-4 py-2 text-gray-600 hover:text-darken rounded-lg hover:bg-gray-50 transition-colors">
                  {{ item.label }}
                  <svg class="w-3 h-3 mt-0.5 group-hover:rotate-180 transition-transform duration-200" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5"/>
                  </svg>
                </button>
                <div class="absolute left-0 top-full pt-1 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 min-w-[200px] z-50">
                  <div class="bg-white rounded-xl shadow-xl border border-gray-100 py-2 overflow-hidden">
                    <router-link
                      v-for="sub in (item.children || item.megaItems || [])"
                      :key="sub.id || sub.label"
                      :to="sub.url || '/'"
                      class="flex items-center gap-2 px-4 py-2.5 text-sm text-gray-700 hover:bg-primary-50 hover:text-primary-700 transition-colors"
                    >
                      <span v-if="sub.icon" class="text-base">{{ sub.icon }}</span>
                      {{ sub.label }}
                    </router-link>
                  </div>
                </div>
              </div>
              <router-link
                v-else
                :to="item.url || '/'"
                class="px-4 py-2 text-gray-600 hover:text-darken rounded-lg hover:bg-gray-50 transition-colors"
              >{{ item.label }}</router-link>
            </template>
          </nav>

          <!-- Desktop CTAs -->
          <div class="hidden lg:flex items-center gap-2">
            <!-- Sign In icon (unauthenticated) -->
            <router-link
              v-if="!isAuthenticated"
              to="/login"
              title="Sign In"
              class="w-10 h-10 flex items-center justify-center rounded-full border border-gray-300 text-gray-600 hover:border-[#112b4e] hover:text-[#112b4e] hover:bg-gray-50 transition-all"
            ><i class="icofont-sign-in text-xl leading-none"></i></router-link>
            <!-- Register icon (unauthenticated) -->
            <router-link
              v-if="!isAuthenticated"
              to="/register"
              title="Create Account"
              class="w-10 h-10 flex items-center justify-center rounded-full border border-gray-300 text-gray-600 hover:border-[#F48C06] hover:text-[#F48C06] hover:bg-orange-50 transition-all"
            ><i class="icofont-ui-user-group text-xl leading-none"></i></router-link>
            <!-- Portal icon (authenticated applicant) -->
            <router-link
              v-if="isAuthenticated && isApplicant"
              to="/my-portal"
              title="My Portal"
              class="w-10 h-10 flex items-center justify-center rounded-full border border-gray-300 text-gray-600 hover:border-[#112b4e] hover:text-[#112b4e] hover:bg-gray-50 transition-all"
            ><i class="icofont-ui-home text-xl leading-none"></i></router-link>
            <!-- Dashboard icon (authenticated admin/staff) -->
            <router-link
              v-if="isAuthenticated && !isApplicant"
              to="/dashboard"
              title="Dashboard"
              class="w-10 h-10 flex items-center justify-center rounded-full border border-gray-300 text-gray-600 hover:border-[#112b4e] hover:text-[#112b4e] hover:bg-gray-50 transition-all"
            ><i class="icofont-dashboard-web text-xl leading-none"></i></router-link>
            <router-link
              to="/apply"
              class="px-6 py-2.5 text-sm font-semibold text-white bg-accent rounded-full hover:bg-yellow-600 transition-colors shadow-sm hover:shadow-md"
            >Apply Now</router-link>
          </div>

          <!-- Mobile toggle -->
          <button
            class="lg:hidden p-2 rounded-lg text-gray-600 hover:bg-gray-100 transition-colors"
            @click="mobileOpen = !mobileOpen"
            aria-label="Toggle menu"
          >
            <svg v-if="!mobileOpen" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12"/>
            </svg>
            <svg v-else class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Mobile menu -->
      <Transition name="mobile-drop">
        <div v-if="mobileOpen" class="lg:hidden bg-white border-t border-gray-100 shadow-lg">
          <div class="max-w-screen-xl px-6 mx-auto py-4 flex flex-col gap-0.5">
            <template v-for="item in menuItems" :key="item.id || item.label">
              <router-link
                :to="item.url || '/'"
                class="px-4 py-2.5 text-sm font-medium text-gray-700 hover:text-primary-700 hover:bg-primary-50 rounded-lg transition-colors"
                @click="mobileOpen = false"
              >{{ item.label }}</router-link>
              <router-link
                v-for="sub in (item.children || item.megaItems || [])"
                :key="sub.id || sub.label"
                :to="sub.url || '/'"
                class="pl-8 py-2 text-sm text-gray-500 hover:text-primary-700 rounded-lg transition-colors"
                @click="mobileOpen = false"
              >↳ {{ sub.label }}</router-link>
            </template>
            <div class="flex gap-3 mt-4 pt-4 border-t border-gray-100">
              <router-link
                v-if="!isAuthenticated"
                to="/login"
                class="flex-1 py-2.5 text-sm font-semibold text-center border border-gray-300 rounded-full text-gray-700 flex items-center justify-center gap-1.5"
                @click="mobileOpen = false"
              ><i class="icofont-sign-in"></i> Sign In</router-link>
              <router-link
                v-if="!isAuthenticated"
                to="/register"
                class="flex-1 py-2.5 text-sm font-semibold text-center border border-[#F48C06] text-[#F48C06] rounded-full flex items-center justify-center gap-1.5"
                @click="mobileOpen = false"
              ><i class="icofont-ui-user-group"></i> Register</router-link>
              <router-link
                v-if="isAuthenticated && isApplicant"
                to="/my-portal"
                class="flex-1 py-2.5 text-sm font-semibold text-center border border-gray-300 rounded-full text-gray-700 flex items-center justify-center gap-1.5"
                @click="mobileOpen = false"
              ><i class="icofont-ui-home"></i> My Portal</router-link>
              <router-link
                v-if="isAuthenticated && !isApplicant"
                to="/dashboard"
                class="flex-1 py-2.5 text-sm font-semibold text-center border border-gray-300 rounded-full text-gray-700 flex items-center justify-center gap-1.5"
                @click="mobileOpen = false"
              ><i class="icofont-dashboard-web"></i> Dashboard</router-link>
              <router-link
                to="/apply"
                class="flex-1 py-2.5 text-sm font-semibold text-center bg-accent text-white rounded-full"
                @click="mobileOpen = false"
              >Apply Now</router-link>
            </div>
          </div>
        </div>
      </Transition>
    </header>

    <!-- Page content -->
    <main class="flex-1">
      <router-view />
    </main>

    <!-- ── Footer ──────────────────────────────────────────── -->
    <footer style="background-color: #252641;">

      <!-- Newsletter bar -->
      <div class="border-b border-white/10">
        <div class="max-w-screen-xl mx-auto px-6 py-10">
          <div class="flex flex-col md:flex-row items-center justify-between gap-6">
            <div>
              <h3 class="text-white font-bold text-lg mb-1">Subscribe to NCS Updates</h3>
              <p class="text-gray-400 text-sm">Get the latest sports news, events and regulatory updates.</p>
            </div>
            <div class="flex w-full md:w-auto gap-3">
              <input
                type="email"
                placeholder="Your email address"
                class="flex-1 md:w-64 bg-white/10 border border-white/20 text-white placeholder-gray-500 rounded-full px-5 py-2.5 text-sm focus:outline-none focus:border-accent transition-colors"
              />
              <button class="bg-accent hover:bg-yellow-600 text-white font-semibold px-6 py-2.5 rounded-full text-sm transition-colors whitespace-nowrap shadow-sm">
                Subscribe
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Links grid -->
      <div class="max-w-screen-xl mx-auto px-6 py-14">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-10">

          <!-- Brand + contact -->
          <div>
            <div class="mb-5">
              <div class="w-20 h-20 rounded-xl bg-white flex items-center justify-center overflow-hidden shadow-sm border border-white/10">
                <img src="/main-logo.png" alt="NCS Logo" class="w-18 h-18 object-contain" style="width:4.25rem;height:4.25rem" />
              </div>
            </div>
            <p class="text-gray-400 text-sm leading-relaxed mb-5">{{ footerSettings.about }}</p>
            <div class="space-y-2 text-sm text-gray-400">
              <div v-if="contact.phone" class="flex items-center gap-2">
                <i class="icofont-phone text-accent text-base flex-shrink-0"></i>
                <a :href="`tel:${contact.phone}`" class="hover:text-accent transition-colors">{{ contact.phone }}</a>
              </div>
              <div v-if="contact.email" class="flex items-center gap-2">
                <i class="icofont-email text-accent text-base flex-shrink-0"></i>
                <a :href="`mailto:${contact.email}`" class="hover:text-accent transition-colors break-all">{{ contact.email }}</a>
              </div>
              <div v-if="contact.address" class="flex items-start gap-2">
                <i class="icofont-location-pin text-accent text-base flex-shrink-0 mt-0.5"></i>
                <span>{{ contact.address }}</span>
              </div>
            </div>
            <!-- Social icons -->
            <div v-if="hasSocial" class="flex gap-2 mt-5">
              <a v-if="contact.social?.facebook"  :href="contact.social.facebook"  target="_blank" rel="noopener" class="w-8 h-8 rounded-full bg-white/10 flex items-center justify-center text-gray-400 hover:text-accent hover:bg-white/20 transition-all"><i class="icofont-facebook text-sm"></i></a>
              <a v-if="contact.social?.twitter"   :href="contact.social.twitter"   target="_blank" rel="noopener" class="w-8 h-8 rounded-full bg-white/10 flex items-center justify-center text-gray-400 hover:text-accent hover:bg-white/20 transition-all"><i class="icofont-twitter text-sm"></i></a>
              <a v-if="contact.social?.linkedin"  :href="contact.social.linkedin"  target="_blank" rel="noopener" class="w-8 h-8 rounded-full bg-white/10 flex items-center justify-center text-gray-400 hover:text-accent hover:bg-white/20 transition-all"><i class="icofont-linkedin text-sm"></i></a>
              <a v-if="contact.social?.instagram" :href="contact.social.instagram" target="_blank" rel="noopener" class="w-8 h-8 rounded-full bg-white/10 flex items-center justify-center text-gray-400 hover:text-accent hover:bg-white/20 transition-all"><i class="icofont-instagram text-sm"></i></a>
              <a v-if="contact.social?.youtube"   :href="contact.social.youtube"   target="_blank" rel="noopener" class="w-8 h-8 rounded-full bg-white/10 flex items-center justify-center text-gray-400 hover:text-accent hover:bg-white/20 transition-all"><i class="icofont-youtube text-sm"></i></a>
            </div>
          </div>

          <!-- Dynamic link columns -->
          <div v-for="(col, idx) in displayedColumns" :key="idx">
            <h4 class="text-white font-semibold text-xs uppercase tracking-widest mb-5">{{ col.title }}</h4>
            <ul class="space-y-3">
              <li v-for="(link, li) in (col.links || [])" :key="li">
                <component
                  :is="isExternalLink(link.url) ? 'a' : 'router-link'"
                  v-bind="isExternalLink(link.url) ? { href: link.url, target: '_blank', rel: 'noopener' } : { to: link.url || '/' }"
                  class="text-sm text-gray-400 hover:text-accent transition-colors"
                >{{ link.label }}</component>
              </li>
              <li v-if="!col.links?.length" class="text-xs text-gray-600 italic">—</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Bottom bar -->
      <div class="border-t border-white/10">
        <div class="max-w-screen-xl mx-auto px-6 py-5 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs text-gray-500">
          <span>&copy; {{ currentYear }} {{ footerSettings.copyright }}</span>
          <div class="flex gap-5">
            <router-link to="/faqs" class="hover:text-gray-300 transition-colors">FAQs</router-link>
            <router-link to="/contact-us" class="hover:text-gray-300 transition-colors">Contact</router-link>
            <router-link to="/apply" class="hover:text-gray-300 transition-colors">Apply</router-link>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, reactive } from 'vue'
import { useAuthStore } from '@/stores/auth.js'
import { getMenu, getSettings } from '@/api/cms.js'

const mobileOpen = ref(false)
const scrolled = ref(false)
const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isApplicant = computed(() => authStore.isApplicant)
const currentYear = computed(() => new Date().getFullYear())

function onScroll() { scrolled.value = window.scrollY > 20 }
onMounted(() => window.addEventListener('scroll', onScroll, { passive: true }))
onUnmounted(() => window.removeEventListener('scroll', onScroll))

const defaultMenu = [
  { label: 'Home',            url: '/' },
  { label: 'Apply',           url: '/apply' },
  { label: 'News',            url: '/news' },
  { label: 'Events',          url: '/events' },
  { label: 'Careers',         url: '/careers' },
  { label: 'Facilities',      url: '/facilities' },
  { label: 'Associations',    url: '/associations' },
  { label: 'Resource Centre', url: '/resource-centre' },
  { label: 'Invest with Us',  url: '/invest' },
]

const defaultFooter = {
  about: 'The National Council of Sports is the government body responsible for the development, promotion and regulation of sports in Uganda.',
  copyright: 'National Council of Sports, Uganda. All rights reserved.',
  columns: [
    { title: 'Services',    links: [
      { label: 'Apply for License', url: '/apply' },
      { label: 'Resource Centre',   url: '/resource-centre' },
      { label: 'FAQs',              url: '/faqs' },
    ]},
    { title: 'Information', links: [
      { label: 'News & Updates', url: '/news' },
      { label: 'Events',         url: '/events' },
      { label: 'Careers',        url: '/careers' },
    ]},
    { title: 'Explore',     links: [
      { label: 'Facilities',     url: '/facilities' },
      { label: 'Associations',   url: '/associations' },
      { label: 'Invest with Us', url: '/invest' },
    ]},
  ]
}

const menuItems = ref(defaultMenu)
const footerSettings = reactive({ ...defaultFooter, columns: [...defaultFooter.columns] })

const displayedColumns = computed(() => {
  const cols = (footerSettings.columns || []).slice(0, 3)
  while (cols.length < 3) cols.push({ title: '', links: [] })
  return cols
})

function isExternalLink(url) {
  return /^(https?:|mailto:|tel:)/i.test(url || '')
}

const contact = reactive({
  phone: '', email: '', address: 'Plot 6, Impala Avenue, Kampala, Uganda',
  hours: '', mapUrl: '',
  social: { facebook: '', twitter: '', linkedin: '', instagram: '', youtube: '' }
})
const hasSocial = computed(() => Object.values(contact.social || {}).some(v => v?.trim()))

async function loadMenu() {
  try {
    const r = await getMenu('main')
    const items = r.data?.data?.items || r.data?.data?.Items || []
    if (items.length) menuItems.value = items
  } catch { /* keep defaults */ }
}

async function loadFooterSettings() {
  try {
    const r = await getSettings('footer')
    const v = r.data?.data?.value
    if (v && typeof v === 'object' && Object.keys(v).length) {
      if (typeof v.about     === 'string') footerSettings.about = v.about
      if (typeof v.copyright === 'string') footerSettings.copyright = v.copyright
      if (Array.isArray(v.columns) && v.columns.length) footerSettings.columns = v.columns
    }
  } catch { /* show defaults */ }
}

async function loadContact() {
  try {
    const r = await getSettings('contact')
    const v = r.data?.data?.value
    if (v && typeof v === 'object' && Object.keys(v).length) {
      if (typeof v.phone   === 'string') contact.phone = v.phone
      if (typeof v.email   === 'string') contact.email = v.email
      if (typeof v.address === 'string') contact.address = v.address
      if (typeof v.hours   === 'string') contact.hours = v.hours
      if (typeof v.mapUrl  === 'string') contact.mapUrl = v.mapUrl
      if (v.social) Object.assign(contact.social, v.social)
    }
  } catch { /* keep defaults */ }
}

onMounted(() => {
  loadMenu()
  loadFooterSettings()
  loadContact()
})
</script>

<style scoped>
.mobile-drop-enter-active, .mobile-drop-leave-active { transition: opacity 0.2s, transform 0.2s; transform-origin: top; }
.mobile-drop-enter-from, .mobile-drop-leave-to { opacity: 0; transform: scaleY(0.95); }
</style>
