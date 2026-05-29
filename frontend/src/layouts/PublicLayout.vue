<template>
  <div class="min-h-screen flex flex-col bg-white">
    <!-- Navigation -->
    <header class="bg-white shadow-sm sticky top-0 z-50 border-b border-gray-100">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo -->
          <router-link to="/" class="flex items-center gap-3">
            <div class="w-12 h-12 bg-white rounded-xl flex items-center justify-center overflow-hidden p-0.5 shadow-sm border border-gray-100">
              <img src="/main-logo.png" alt="NCS Logo" class="w-full h-full object-contain" />
            </div>
            <div class="leading-tight">
              <div class="font-bold text-gray-900 text-sm">National Council of Sports</div>
              <div class="text-xs text-gray-500">Uganda</div>
            </div>
          </router-link>

          <!-- Desktop Navigation (dynamic from CMS) -->
          <nav class="hidden md:flex items-center gap-6 text-sm font-medium">
            <template v-for="item in menuItems" :key="item.id || item.label">
              <div v-if="(item.children && item.children.length) || item.mega" class="relative group">
                <button class="flex items-center gap-1 text-gray-600 hover:text-primary-700 transition-colors">
                  {{ item.label }}
                  <i class="icofont-rounded-down text-xs"></i>
                </button>
                <div class="absolute left-0 top-full mt-1 min-w-[200px] bg-white rounded-lg shadow-lg border border-gray-100 py-2 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50">
                  <router-link
                    v-for="sub in (item.children || item.megaItems || [])"
                    :key="sub.id || sub.label"
                    :to="sub.url || '/'"
                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-primary-50 hover:text-primary-700"
                  >
                    <span v-if="sub.icon" class="mr-2">{{ sub.icon }}</span>{{ sub.label }}
                  </router-link>
                </div>
              </div>
              <router-link
                v-else
                :to="item.url || '/'"
                class="text-gray-600 hover:text-primary-700 transition-colors"
              >
                {{ item.label }}
              </router-link>
            </template>
          </nav>

          <!-- CTA -->
          <div class="flex items-center gap-3">
            <router-link
              v-if="!isAuthenticated"
              to="/login"
              class="text-sm font-medium text-gray-600 hover:text-primary-700 transition-colors"
            >
              Sign In
            </router-link>
            <router-link
              v-if="isAuthenticated"
              to="/dashboard"
              class="text-sm font-medium text-primary-700 hover:text-primary-800 transition-colors"
            >
              Dashboard
            </router-link>
            <router-link
              to="/apply"
              class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-semibold px-4 py-2 rounded-lg transition-colors"
            >
              Apply Now
            </router-link>

            <button
              class="md:hidden p-2 rounded-md text-gray-500 hover:text-gray-700"
              @click="mobileOpen = !mobileOpen"
            >
              <i :class="mobileOpen ? 'icofont-close' : 'icofont-navigation-menu'" class="text-xl"></i>
            </button>
          </div>
        </div>

        <!-- Mobile menu -->
        <div v-if="mobileOpen" class="md:hidden pb-4 pt-2 border-t border-gray-100 flex flex-col gap-3 text-sm font-medium">
          <template v-for="item in menuItems" :key="item.id || item.label">
            <router-link
              :to="item.url || '/'"
              class="text-gray-700 hover:text-primary-700"
              @click="mobileOpen = false"
            >{{ item.label }}</router-link>
            <router-link
              v-for="sub in (item.children || item.megaItems || [])"
              :key="sub.id || sub.label"
              :to="sub.url || '/'"
              class="pl-4 text-gray-500 hover:text-primary-700 text-xs"
              @click="mobileOpen = false"
            >↳ {{ sub.label }}</router-link>
          </template>
        </div>
      </div>
    </header>

    <!-- Page Content -->
    <main class="flex-1">
      <router-view />
    </main>

    <!-- Footer: 4-column layout. Col 1 = brand/about/contact (fixed
         shape, editable text). Cols 2-4 = dynamic title + links from
         the footer settings in the CMS. -->
    <footer class="bg-gray-900 text-gray-300 mt-auto">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-8">
          <!-- Column 1: brand + about + contact -->
          <div>
            <div class="flex items-center gap-3 mb-4">
              <div class="w-10 h-10 rounded-full bg-primary-600 flex items-center justify-center">
                <span class="text-white font-bold text-sm">NCS</span>
              </div>
              <div>
                <div class="font-bold text-white text-sm">National Council of Sports</div>
                <div class="text-xs text-gray-400">Republic of Uganda</div>
              </div>
            </div>
            <p class="text-sm text-gray-400">{{ footerSettings.about }}</p>

            <!-- Contact details (dynamic) -->
            <div class="mt-5 space-y-1.5 text-sm text-gray-400">
              <div v-if="contact.phone" class="flex items-center gap-2">
                <i class="icofont-phone text-primary-400"></i>
                <a :href="`tel:${contact.phone}`" class="hover:text-primary-400">{{ contact.phone }}</a>
              </div>
              <div v-if="contact.email" class="flex items-center gap-2">
                <i class="icofont-email text-primary-400"></i>
                <a :href="`mailto:${contact.email}`" class="hover:text-primary-400 break-all">{{ contact.email }}</a>
              </div>
              <div v-if="contact.address" class="flex items-start gap-2">
                <i class="icofont-location-pin text-primary-400 mt-0.5 flex-shrink-0"></i>
                <span>{{ contact.address }}</span>
              </div>
              <div v-if="contact.hours" class="flex items-center gap-2">
                <i class="icofont-clock-time text-primary-400"></i>
                <span>{{ contact.hours }}</span>
              </div>
            </div>

            <!-- Social links (dynamic) -->
            <div v-if="hasSocial" class="flex gap-3 mt-5">
              <a v-if="contact.social?.facebook"  :href="contact.social.facebook"  target="_blank" rel="noopener" class="text-gray-400 hover:text-primary-400 text-lg"><i class="icofont-facebook"></i></a>
              <a v-if="contact.social?.twitter"   :href="contact.social.twitter"   target="_blank" rel="noopener" class="text-gray-400 hover:text-primary-400 text-lg"><i class="icofont-twitter"></i></a>
              <a v-if="contact.social?.linkedin"  :href="contact.social.linkedin"  target="_blank" rel="noopener" class="text-gray-400 hover:text-primary-400 text-lg"><i class="icofont-linkedin"></i></a>
              <a v-if="contact.social?.instagram" :href="contact.social.instagram" target="_blank" rel="noopener" class="text-gray-400 hover:text-primary-400 text-lg"><i class="icofont-instagram"></i></a>
              <a v-if="contact.social?.youtube"   :href="contact.social.youtube"   target="_blank" rel="noopener" class="text-gray-400 hover:text-primary-400 text-lg"><i class="icofont-youtube"></i></a>
            </div>
          </div>

          <!-- Columns 2-4: dynamic link columns -->
          <div v-for="(col, idx) in displayedColumns" :key="idx">
            <h4 class="font-semibold text-white mb-3 text-sm uppercase tracking-wide">{{ col.title }}</h4>
            <ul class="space-y-2 text-sm">
              <li v-for="(link, li) in (col.links || [])" :key="li">
                <component
                  :is="isExternalLink(link.url) ? 'a' : 'router-link'"
                  v-bind="isExternalLink(link.url) ? { href: link.url, target: '_blank', rel: 'noopener' } : { to: link.url || '/' }"
                  class="hover:text-primary-400 transition-colors"
                >{{ link.label }}</component>
              </li>
              <li v-if="!col.links?.length" class="text-xs text-gray-600 italic">No links</li>
            </ul>
          </div>
        </div>

        <div class="border-t border-gray-800 pt-6 text-sm text-gray-500 flex flex-col md:flex-row justify-between gap-2">
          <span>&copy; {{ currentYear }} {{ footerSettings.copyright }}</span>
          <span v-if="contact.address">{{ contact.address }}</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useAuthStore } from '@/stores/auth.js'
import { getMenu, getSettings } from '@/api/cms.js'

const mobileOpen = ref(false)
const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)
const currentYear = computed(() => new Date().getFullYear())

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

// Footer always renders exactly 3 link columns (cols 2-4 of a 4-col grid).
// Slice to 3 max; if fewer are configured, pad with empty headings so
// the grid layout doesn't collapse.
const displayedColumns = computed(() => {
  const cols = (footerSettings.columns || []).slice(0, 3)
  while (cols.length < 3) cols.push({ title: '', links: [] })
  return cols
})

function isExternalLink(url) {
  return /^(https?:|mailto:|tel:)/i.test(url || '')
}

const contact = reactive({
  phone: '',
  email: '',
  address: 'Plot 6, Impala Avenue, Kampala, Uganda',
  hours: '',
  mapUrl: '',
  social: { facebook: '', twitter: '', linkedin: '', instagram: '', youtube: '' }
})
const hasSocial = computed(() =>
  Object.values(contact.social || {}).some(v => v && v.trim())
)

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
      if (Array.isArray(v.columns) && v.columns.length) {
        footerSettings.columns = v.columns
      }
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
