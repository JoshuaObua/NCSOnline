<template>
  <div class="public-site min-h-screen flex flex-col">

    <a href="#main-content" class="skip-link">Skip to main content</a>
    <a href="#public-footer" class="skip-link skip-link-secondary">Skip to footer</a>

    <!-- ── Navbar ──────────────────────────────────────────── -->
    <header
      :class="[
        'fixed top-0 left-0 right-0 z-50 transition-all duration-500',
        scrolled
          ? 'bg-white/95 backdrop-blur-md shadow-lg'
          : 'bg-white/90 backdrop-blur-sm shadow-sm'
      ]"
    >
      <div class="public-topbar">
        <div class="max-w-screen-xl mx-auto px-4 flex items-center gap-4">
          <div class="public-marquee" aria-label="NCS highlights">
            <div class="public-marquee-track">
              <template v-for="copy in 2" :key="copy">
                <span v-for="(message, index) in headerSettings.marquee" :key="`${copy}-${index}`" class="public-marquee-item">
                  {{ message }} <span aria-hidden="true">•</span>
                </span>
              </template>
            </div>
          </div>
          <div class="public-topbar-social" aria-label="Social media links">
            <a v-for="network in socialNetworks" :key="network.key" v-show="network.url" :href="network.url" target="_blank" rel="noopener" :aria-label="network.label"><i :class="network.icon" aria-hidden="true"></i></a>
            <a v-if="headerSettings.webmail_url" :href="headerSettings.webmail_url" target="_blank" rel="noopener" aria-label="Webmail"><i class="icofont-email" aria-hidden="true"></i></a>
          </div>
        </div>
      </div>
      <div class="max-w-screen-xl px-6 mx-auto">
        <div class="flex items-center justify-between h-20">

          <!-- Logo -->
          <router-link to="/" class="flex items-center flex-shrink-0" aria-label="National Council of Sports home">
            <div class="relative flex-shrink-0">
              <div class="w-16 h-16 bg-white rounded-xl flex items-center justify-center overflow-hidden shadow-sm border border-gray-100 z-10 relative">
                <img :src="resolveAsset(siteIdentity.logoUrl) || '/main-logo.png'" :alt="`${siteIdentity.name || 'NCS'} logo`" class="w-14 h-14 object-contain" />
              </div>
              <!-- Diamond accent -->
              <svg class="absolute -top-1.5 -left-1.5 w-8 h-8 z-0 opacity-30" viewBox="0 0 79 79" fill="none" aria-hidden="true" focusable="false">
                <path d="M35.26 2.24C37.6-.1 41.4-.1 43.74 2.24L76.76 35.26C79.1 37.6 79.1 41.4 76.76 43.74L43.74 76.76C41.4 79.1 37.6 79.1 35.26 76.76L2.24 43.74C-.1 41.4-.1 37.6 2.24 35.26L35.26 2.24Z" fill="#2F327D"/>
              </svg>
            </div>
          </router-link>

          <!-- Desktop nav -->
          <nav class="hidden lg:flex items-center gap-1 text-sm font-medium" aria-label="Primary navigation">
            <template v-for="item in menuItems" :key="item.id || item.label">
              <div v-if="(item.children?.length) || item.mega" class="relative group">
                <button :class="['flex items-center gap-1 px-3 py-2 font-medium transition-colors duration-300 text-sm', isActiveTopLevel(item) ? 'text-[#f5a623]' : 'text-[#1a365d] hover:text-[#f5a623]']" aria-haspopup="true">
                  {{ item.label }}
                  <svg class="w-4 h-4 group-hover:rotate-180 transition-transform duration-200" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                    <path d="m6 9 6 6 6-6"/>
                  </svg>
                </button>
                <div class="absolute left-0 top-full pt-1 opacity-0 invisible group-hover:opacity-100 group-hover:visible group-focus-within:opacity-100 group-focus-within:visible transition-all duration-200 min-w-[220px] z-50">
                  <div class="bg-white rounded-xl shadow-xl border border-gray-100 py-2 overflow-hidden">
                    <router-link
                      v-for="sub in (item.children || item.megaItems || [])"
                      :key="sub.id || sub.label"
                      :to="sub.url || '/'"
                      class="flex items-center gap-2 px-4 py-2.5 text-sm text-[#1a365d] hover:bg-[#f5a623]/10 hover:text-[#f5a623] transition-colors"
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
                :class="['px-3 py-2 font-medium transition-colors duration-300 text-sm', isActiveLink(item) ? 'text-[#f5a623]' : 'text-[#1a365d] hover:text-[#f5a623]']"
              >{{ item.label }}</router-link>
            </template>
          </nav>

          <!-- Desktop CTAs -->
          <div class="hidden lg:flex items-center gap-2">
            <form class="public-header-search" :class="{ 'is-open': searchOpen }" @submit.prevent="submitSearch">
              <label for="public-site-search" class="sr-only">Search NCS website</label>
              <input id="public-site-search" v-model="siteSearch" type="search" placeholder="Search..." />
            </form>
            <button type="button" class="public-header-icon" :aria-expanded="searchOpen" aria-controls="public-site-search" aria-label="Search website" @click="searchOpen = !searchOpen">
              <i class="icofont-search-1" aria-hidden="true"></i>
            </button>
            <button type="button" class="public-header-icon" aria-label="Open accessibility tools" title="Accessibility tools" @click="openAccessibility">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" class="w-5 h-5"><circle cx="16" cy="4" r="1"/><path d="m18 19 1-7-6 1"/><path d="m5 8 3-3 5.5 3-2.36 3.5"/><path d="M4.24 14.5a5 5 0 0 0 6.88 6"/><path d="M13.76 17.5a5 5 0 0 0-6.88-6"/></svg>
            </button>
            <!-- Login button (unauthenticated) -->
            <router-link
              v-if="!isAuthenticated"
              to="/login"
              class="hidden md:inline-flex h-8 items-center gap-2 rounded-md px-3 text-xs font-medium text-[#1a365d] hover:text-[#f5a623] hover:bg-[#f5a623]/10 transition-colors"
              aria-label="Login"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" x2="3" y1="12" y2="12"/></svg>
              Login
            </router-link>
            <!-- Register icon (unauthenticated, mobile compact) -->
            <router-link
              v-if="!isAuthenticated"
              to="/register"
              title="Create Account"
              aria-label="Create an account"
              class="md:hidden w-10 h-10 flex items-center justify-center rounded-full border border-gray-300 text-[#1a365d] hover:border-[#f5a623] hover:text-[#f5a623] hover:bg-[#f5a623]/10 transition-all"
            ><i class="icofont-ui-user-group text-xl leading-none" aria-hidden="true"></i></router-link>
            <!-- Portal icon (authenticated applicant) -->
            <router-link
              v-if="isAuthenticated && isApplicant"
              to="/my-portal"
              title="My Portal"
              aria-label="Open my portal"
              class="w-10 h-10 flex items-center justify-center rounded-full border border-gray-300 text-gray-600 hover:border-[#112b4e] hover:text-[#112b4e] hover:bg-gray-50 transition-all"
            ><i class="icofont-ui-home text-xl leading-none" aria-hidden="true"></i></router-link>
            <!-- Dashboard icon (authenticated admin/staff) -->
            <router-link
              v-if="isAuthenticated && !isApplicant"
              to="/dashboard"
              title="Dashboard"
              aria-label="Open dashboard"
              class="w-10 h-10 flex items-center justify-center rounded-full border border-gray-300 text-gray-600 hover:border-[#112b4e] hover:text-[#112b4e] hover:bg-gray-50 transition-all"
            ><i class="icofont-dashboard-web text-xl leading-none" aria-hidden="true"></i></router-link>
          </div>

          <!-- Mobile toggle -->
          <button
            class="lg:hidden p-2 rounded-lg text-gray-600 hover:bg-gray-100 transition-colors"
            @click="mobileOpen = !mobileOpen"
            :aria-label="mobileOpen ? 'Close navigation menu' : 'Open navigation menu'"
            :aria-expanded="mobileOpen"
            aria-controls="mobile-navigation"
          >
            <svg v-if="!mobileOpen" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12"/>
            </svg>
            <svg v-else class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Mobile menu -->
      <Transition name="mobile-drop">
        <nav v-if="mobileOpen" id="mobile-navigation" class="lg:hidden bg-white border-t border-gray-100 shadow-lg" aria-label="Mobile navigation">
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
                to="/my-portal"
                class="flex-1 py-2.5 text-sm font-semibold text-center bg-accent text-white rounded-full"
                @click="mobileOpen = false"
              >My Portal</router-link>
            </div>
          </div>
        </nav>
      </Transition>
    </header>

    <!-- Page content -->
    <main id="main-content" tabindex="-1" class="flex-1">
      <router-view />
    </main>

    <!-- ── Footer ──────────────────────────────────────────── -->
    <footer id="public-footer" tabindex="-1" style="background-color: #252641;" aria-label="Website footer">

      <!-- Newsletter bar -->
      <div class="border-b border-white/10">
        <div class="max-w-screen-xl mx-auto px-6 py-10">
          <div class="flex flex-col md:flex-row items-center justify-between gap-6">
            <div>
              <h3 class="text-white font-bold text-lg mb-1">Subscribe to NCS Updates</h3>
              <p class="text-gray-400 text-sm">Get the latest sports news, events and regulatory updates.</p>
            </div>
            <form class="flex w-full md:w-auto gap-3" @submit.prevent>
              <label for="newsletter-email" class="sr-only">Email address for NCS updates</label>
              <input
                id="newsletter-email"
                type="email"
                placeholder="Your email address"
                autocomplete="email"
                required
                class="flex-1 md:w-64 bg-white/10 border border-white/20 text-white placeholder-gray-500 rounded-full px-5 py-2.5 text-sm focus:outline-none focus:border-accent transition-colors"
              />
              <button type="submit" class="bg-accent hover:bg-yellow-600 text-white font-semibold px-6 py-2.5 rounded-full text-sm transition-colors whitespace-nowrap shadow-sm">
                Subscribe
              </button>
            </form>
          </div>
        </div>
      </div>

      <!-- Links grid (NCS Footer Section Design spec) -->
      <div class="max-w-7xl mx-auto px-4 py-12 md:py-16">
        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-10 md:gap-8">

          <!-- Column 1 — Brand + about + social -->
          <div>
            <div class="flex items-center gap-3 mb-6">
              <img :src="footerLogo" :alt="`${siteIdentity.name || 'NCS'} - National Council of Sports`" class="h-14 w-auto object-contain" :class="!siteIdentity.footerWhiteLogoUrl && !siteIdentity.whiteLogoUrl ? 'brightness-0 invert' : ''" />
            </div>
            <p class="text-white/80 text-sm leading-relaxed mb-6">{{ footerSettings.about }}</p>
            <div class="flex gap-3">
              <a v-for="net in footerSocialLinks" :key="net.key" :href="net.url" target="_blank" rel="noopener noreferrer" :aria-label="net.label" class="w-10 h-10 rounded-full bg-white/10 flex items-center justify-center text-white hover:bg-[#f5a623] hover:text-white transition-all duration-300 hover:-translate-y-1">
                <i :class="[net.icon, 'text-base']" aria-hidden="true"></i>
              </a>
            </div>
          </div>

          <!-- Column 2 — Contact Us -->
          <div>
            <h4 class="text-lg font-bold mb-6 flex items-center gap-2 text-white">
              <span class="w-8 h-0.5 bg-[#f5a623]"></span>Contact Us
            </h4>
            <div class="space-y-4">
              <div v-if="contactAddressLines.length" class="flex items-start gap-3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623] flex-shrink-0 mt-0.5" aria-hidden="true"><path d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0"/><circle cx="12" cy="10" r="3"/></svg>
                <div class="text-sm text-white/80">
                  <p v-for="(line, i) in contactAddressLines" :key="i">{{ line }}</p>
                </div>
              </div>
              <div v-if="contact.phone" class="flex items-center gap-3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623] flex-shrink-0" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
                <a :href="`tel:${primaryPhone}`" class="text-sm text-white/80 hover:text-[#f5a623] transition-colors">{{ contact.phone }}</a>
              </div>
              <div v-if="contact.email" class="flex items-center gap-3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623] flex-shrink-0" aria-hidden="true"><path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"/><rect x="2" y="4" width="20" height="16" rx="2"/></svg>
                <a :href="`mailto:${contact.email}`" class="text-sm text-white/80 hover:text-[#f5a623] transition-colors break-all">{{ contact.email }}</a>
              </div>
            </div>
            <a v-if="headerSettings.webmail_url" :href="headerSettings.webmail_url" target="_blank" rel="noopener noreferrer" class="inline-flex items-center justify-center gap-2 mt-6 h-9 px-4 py-2 rounded-md text-sm font-medium bg-[#f5a623] hover:bg-[#e09612] text-white shadow transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"/><rect x="2" y="4" width="20" height="16" rx="2"/></svg>
              Webmail
            </a>
          </div>

          <!-- Column 3 — Quick Links (dynamic from footerSettings.columns[0]) -->
          <div>
            <h4 class="text-lg font-bold mb-6 flex items-center gap-2 text-white">
              <span class="w-8 h-0.5 bg-[#f5a623]"></span>{{ quickLinksColumn.title || 'Quick Links' }}
            </h4>
            <ul class="space-y-3">
              <li v-for="(link, idx) in quickLinksColumn.links" :key="idx">
                <component
                  :is="isExternalLink(link.url) ? 'a' : 'router-link'"
                  v-bind="isExternalLink(link.url) ? { href: link.url, target: '_blank', rel: 'noopener' } : { to: link.url || '/' }"
                  class="flex items-center gap-2 text-sm text-white/80 hover:text-[#f5a623] transition-colors group"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 group-hover:translate-x-1 transition-transform" aria-hidden="true"><path d="m9 18 6-6-6-6"/></svg>
                  {{ link.label }}
                </component>
              </li>
              <li v-if="!quickLinksColumn.links?.length" class="text-xs text-white/40 italic">No quick links configured yet.</li>
            </ul>
          </div>

          <!-- Column 4 — Documents (dynamic from footerSettings.columns[1]) -->
          <div>
            <h4 class="text-lg font-bold mb-6 flex items-center gap-2 text-white">
              <span class="w-8 h-0.5 bg-[#f5a623]"></span>{{ documentsColumn.title || 'Documents' }}
            </h4>
            <ul class="space-y-3">
              <li v-for="(doc, idx) in documentsColumn.links" :key="idx">
                <a :href="doc.url" target="_blank" rel="noopener noreferrer" class="flex items-center gap-2 text-sm text-white/80 hover:text-[#f5a623] transition-colors group">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-red-400 flex-shrink-0" aria-hidden="true"><path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/><path d="M14 2v4a2 2 0 0 0 2 2h4"/><path d="M10 9H8"/><path d="M16 13H8"/><path d="M16 17H8"/></svg>
                  <span class="flex-1 line-clamp-1">{{ doc.label }}</span>
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3 h-3 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0" aria-hidden="true"><path d="M15 3h6v6"/><path d="M10 14 21 3"/><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/></svg>
                </a>
              </li>
              <li v-if="!documentsColumn.links?.length" class="text-xs text-white/40 italic">No documents added yet.</li>
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
            <router-link to="/my-portal" class="hover:text-gray-300 transition-colors">My Portal</router-link>
          </div>
        </div>
      </div>
    </footer>
    <PublicAccessibilityMenu />
  </div>
</template>

<script setup>
import { ref, computed, nextTick, onMounted, onUnmounted, reactive, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import { getMenu, getSettings } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'
import PublicAccessibilityMenu from '@/components/public/PublicAccessibilityMenu.vue'

const mobileOpen = ref(false)
const searchOpen = ref(false)
const siteSearch = ref('')
const scrolled = ref(false)
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isApplicant = computed(() => authStore.isApplicant)
const currentYear = computed(() => new Date().getFullYear())

function onScroll() { scrolled.value = window.scrollY > 20 }
function openAccessibility() { window.dispatchEvent(new CustomEvent('open-accessibility-menu')) }
function submitSearch() {
  const query = siteSearch.value.trim()
  if (query) router.push({ path: '/news', query: { search: query } })
}
onMounted(() => window.addEventListener('scroll', onScroll, { passive: true }))
onUnmounted(() => window.removeEventListener('scroll', onScroll))

watch(() => route.fullPath, async () => {
  mobileOpen.value = false
  await nextTick()
  document.getElementById('main-content')?.focus({ preventScroll: true })
})

const defaultMenu = [
  { label: 'Home',            url: '/' },
  { label: 'My Portal',       url: '/my-portal' },
  { label: 'News',            url: '/news' },
  { label: 'Events',          url: '/events' },
  { label: 'Careers',         url: '/careers' },
  { label: 'Facilities',      url: '/facilities' },
  { label: 'Associations',    url: '/associations' },
  { label: 'Resource Centre', url: '/resource-centre' },
  { label: 'Invest with Us',  url: '/invest' },
]

const defaultFooter = {
  about: 'The National Council of Sports is the statutory body mandated to develop, promote, and control sports in Uganda.',
  copyright: 'National Council of Sports, Uganda. All rights reserved.',
  columns: [
    // Column 0 — rendered as Quick Links (chevron icons)
    { title: 'Quick Links', links: [
      { label: 'Home',           url: '/' },
      { label: 'Associations',   url: '/associations' },
      { label: 'Invest With Us', url: '/invest' },
      { label: 'Careers',        url: '/careers' },
      { label: 'Contact',        url: '/contact-us' },
      { label: 'About NCS',      url: '/about/mandate' },
      { label: 'Latest News',    url: '/news' },
    ]},
    // Column 1 — rendered as Documents (PDF/external file icons)
    { title: 'Documents', links: [
      { label: 'NCS Annual Report 2023/2024',   url: 'https://www.ncs.go.ug/files/NCS%20ANNUAL%20REPORT%202023-2024_0.pdf' },
      { label: 'NCS Annual Report 2022/2023',   url: 'https://www.ncs.go.ug/files/NCS%20ANNUAL%20REPORT%202022-2023.pdf' },
      { label: 'NCS Strategic Plan 2020-2025',  url: 'https://www.ncs.go.ug/files/NCS%20STRATEGIC%20PLAN%202020%202025.pdf' },
      { label: 'Policy Guidelines',             url: 'https://www.ncs.go.ug/files/POLICY%20GUIDELINES.pdf' },
      { label: 'National Development Plan IV',  url: 'https://ncs.go.ug/files/NATIONAL%20DEVELOPMENT%20PLAN%20IV.pdf' },
    ]},
  ]
}

const menuItems = ref(defaultMenu)
const footerSettings = reactive({ ...defaultFooter, columns: [...defaultFooter.columns] })

// Quick Links uses the first dynamic column (CMS), Documents uses the second
const quickLinksColumn = computed(() => footerSettings.columns?.[0] || defaultFooter.columns[0])
const documentsColumn = computed(() => footerSettings.columns?.[1] || defaultFooter.columns[1])

// Address rendered as multi-line block: split contact.address on newlines (or commas),
// then append postal_address as a separate line if present.
const contactAddressLines = computed(() => {
  const lines = []
  if (contact.address) {
    const parts = String(contact.address).split(/\r?\n/).map(s => s.trim()).filter(Boolean)
    lines.push(...(parts.length ? parts : [String(contact.address).trim()]))
  }
  if (contact.postal_address) {
    String(contact.postal_address).split(/\r?\n/).forEach(s => { const t = s.trim(); if (t) lines.push(t) })
  }
  return lines
})

// First number from "+256 414254477 / 343688" so tel: works on mobile.
const primaryPhone = computed(() => {
  if (!contact.phone) return ''
  return String(contact.phone).split(/[\/,]/)[0].replace(/\s+/g, '')
})

function isExternalLink(url) {
  return /^(https?:|mailto:|tel:)/i.test(url || '')
}

const contact = reactive({
  phone: '+256 414254477 / 343688',
  email: 'info@ncs.go.ug',
  address: 'Plot 2-10, Coronation Avenue',
  postal_address: 'P.O. Box 20077, Lugogo\nKampala - UGANDA',
  fax: '',
  hours: '',
  mapUrl: '',
  social: { facebook: '', twitter: '', linkedin: '', instagram: '', youtube: '' }
})
const headerSettings = reactive({
  marquee: [
    'Welcome to National Council of Sports Uganda',
    'A centre of excellence for promotion and development of Sports.',
    'Maximizing opportunities for all Ugandans to participate and excel in Sports.',
    'Established 1964',
  ],
  webmail_url: 'https://mail.umcs.go.ug/'
})
// Site identity (logos + favicon + name) — backed by cms_settings.site
const siteIdentity = reactive({
  name: 'NCS Uganda',
  logoUrl: '',
  whiteLogoUrl: '',
  footerLogoUrl: '',
  footerWhiteLogoUrl: '',
  faviconUrl: '',
})
function resolveAsset(path) {
  if (!path) return ''
  return mediaUrl(path)
}
// Footer logo: prefer white footer logo → footer logo → white logo → main logo → bundled fallback
const footerLogo = computed(() => {
  const chain = [siteIdentity.footerWhiteLogoUrl, siteIdentity.footerLogoUrl, siteIdentity.whiteLogoUrl, siteIdentity.logoUrl]
  for (const p of chain) { if (p) return resolveAsset(p) }
  return '/main-logo.png'
})
function applyFavicon(url) {
  if (!url) return
  const href = resolveAsset(url) || url
  let link = document.querySelector('link[rel="icon"]')
  if (!link) { link = document.createElement('link'); link.rel = 'icon'; document.head.appendChild(link) }
  link.href = href
}
const socialDefaults = {
  facebook:  'https://facebook.com/NCSUganda',
  twitter:   'https://twitter.com/NCSUganda1',
  linkedin:  'https://linkedin.com/company/ncsuganda',
  youtube:   'http://www.youtube.com/@NCSUgTV',
  instagram: '',
}
const socialNetworks = computed(() => [
  { key:'facebook',  label:'Facebook',     icon:'icofont-facebook',  url:contact.social.facebook  || socialDefaults.facebook },
  { key:'twitter',   label:'X / Twitter',  icon:'icofont-twitter',   url:contact.social.twitter   || socialDefaults.twitter },
  { key:'linkedin',  label:'LinkedIn',     icon:'icofont-linkedin',  url:contact.social.linkedin  || socialDefaults.linkedin },
  { key:'instagram', label:'Instagram',    icon:'icofont-instagram', url:contact.social.instagram || socialDefaults.instagram },
  { key:'youtube',   label:'YouTube',      icon:'icofont-youtube',   url:contact.social.youtube   || socialDefaults.youtube },
])
// Footer-spec social row: FB, X, LinkedIn, YouTube, Webmail (Instagram dropped to match the design)
const footerSocialLinks = computed(() => {
  const rows = [
    { key:'facebook',  label:'Facebook on NCS Uganda',  icon:'icofont-facebook',  url: contact.social.facebook  || socialDefaults.facebook },
    { key:'twitter',   label:'X / Twitter',             icon:'icofont-twitter',   url: contact.social.twitter   || socialDefaults.twitter },
    { key:'linkedin',  label:'LinkedIn',                icon:'icofont-linkedin',  url: contact.social.linkedin  || socialDefaults.linkedin },
    { key:'youtube',   label:'YouTube',                 icon:'icofont-youtube',   url: contact.social.youtube   || socialDefaults.youtube },
  ]
  if (headerSettings.webmail_url) {
    rows.push({ key:'webmail', label:'Webmail', icon:'icofont-email', url: headerSettings.webmail_url })
  }
  return rows.filter(r => r.url)
})

function isActiveLink(item) {
  if (!item?.url) return false
  return route.path === item.url
}
function isActiveTopLevel(item) {
  const children = item.children || item.megaItems || []
  return children.some(c => c.url && route.path.startsWith(c.url))
}

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

async function loadHeaderSettings() {
  try {
    const r = await getSettings('header')
    const v = r.data?.data?.value
    if (Array.isArray(v?.marquee) && v.marquee.length) headerSettings.marquee = v.marquee
    if (typeof v?.webmail_url === 'string') headerSettings.webmail_url = v.webmail_url
  } catch { /* keep defaults */ }
}

async function loadSiteIdentity() {
  try {
    const r = await getSettings('site')
    const v = r.data?.data?.value
    if (v && typeof v === 'object') {
      if (typeof v.name === 'string')               siteIdentity.name = v.name
      if (typeof v.logoUrl === 'string')            siteIdentity.logoUrl = v.logoUrl
      if (typeof v.whiteLogoUrl === 'string')       siteIdentity.whiteLogoUrl = v.whiteLogoUrl
      if (typeof v.footerLogoUrl === 'string')      siteIdentity.footerLogoUrl = v.footerLogoUrl
      if (typeof v.footerWhiteLogoUrl === 'string') siteIdentity.footerWhiteLogoUrl = v.footerWhiteLogoUrl
      if (typeof v.faviconUrl === 'string')         siteIdentity.faviconUrl = v.faviconUrl
      applyFavicon(siteIdentity.faviconUrl)
    }
  } catch { /* first-time, no record yet */ }
}

onMounted(() => {
  loadMenu()
  loadFooterSettings()
  loadContact()
  loadHeaderSettings()
  loadSiteIdentity()
})
</script>

<style scoped>
.mobile-drop-enter-active, .mobile-drop-leave-active { transition: opacity 0.2s, transform 0.2s; transform-origin: top; }
.mobile-drop-enter-from, .mobile-drop-leave-to { opacity: 0; transform: scaleY(0.95); }
.public-site > main { padding-top: 7rem; }
.public-topbar { overflow: hidden; background: #1a365d; color: white; padding: .45rem 0; }
.public-marquee { min-width: 0; flex: 1; overflow: hidden; white-space: nowrap; }
.public-marquee-track { display: inline-flex; width: max-content; animation: public-marquee 34s linear infinite; }
.public-marquee-item { display: inline-flex; gap: 2rem; margin-right: 2rem; font-size: .78rem; }
.public-marquee-item span { color: #f5a623; }
.public-topbar-social { display: flex; flex-shrink: 0; gap: .65rem; }
.public-topbar-social a { color: white; transition: color .2s; }
.public-topbar-social a:hover { color: #f5a623; }
.public-header-icon { display: inline-flex; width: 2.5rem; height: 2.5rem; align-items: center; justify-content: center; border-radius: .55rem; color: #1a365d; font-size: 1.15rem; }
.public-header-icon svg { width:1.25rem; height:1.25rem; }
.public-header-icon:hover { color: #f5a623; background: rgb(245 166 35 / .1); }
.public-header-search { width: 0; overflow: hidden; transition: width .25s ease; }
.public-header-search.is-open { width: 11rem; }
.public-header-search input { width: 11rem; border: 1px solid #d1d5db; border-radius: .5rem; padding: .5rem .75rem; font-size: .82rem; }
@keyframes public-marquee { to { transform: translateX(-50%); } }
@media (prefers-reduced-motion: reduce) { .public-marquee-track { animation-play-state: paused; } }
@media (max-width: 640px) { .public-topbar-social { display: none; } }
</style>
