<template>
  <div>
    <!-- ─── Hero Slideshow ───────────────────────────────────────────── -->
    <section ref="heroRef" class="relative overflow-hidden bg-gray-900 min-h-[600px] flex items-center">
      <!-- Slides -->
      <transition-group name="slide-fade" tag="div" class="absolute inset-0">
        <div
          v-for="(slide, i) in displaySlides"
          :key="slide.id || i"
          v-show="currentSlide === i"
          class="absolute inset-0"
        >
          <div
            v-if="slide.image_url"
            class="absolute inset-0 bg-cover bg-center"
            :style="{ backgroundImage: `url(${slide.image_url})` }"
          >
            <div class="absolute inset-0 bg-gradient-to-r from-gray-900/90 via-gray-900/70 to-gray-900/30"></div>
          </div>
          <div v-else class="absolute inset-0 bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900">
            <div class="absolute inset-0 opacity-10">
              <div class="absolute top-0 left-1/4 w-96 h-96 bg-primary-500 rounded-full filter blur-3xl"></div>
              <div class="absolute bottom-0 right-1/4 w-64 h-64 bg-primary-600 rounded-full filter blur-3xl"></div>
            </div>
          </div>
        </div>
      </transition-group>

      <!-- Hero Content -->
      <div class="relative z-10 w-full max-w-7xl mx-auto px-6 py-28">
        <div class="max-w-3xl" ref="heroContent">
          <div class="hero-badge inline-flex items-center gap-2 bg-primary-600/20 border border-primary-500/30 rounded-full px-4 py-1.5 text-sm text-primary-400 mb-6">
            <span class="w-2 h-2 rounded-full bg-primary-400 animate-pulse"></span>
            {{ displaySlides[currentSlide]?.subtitle || 'Official Sports Regulatory Authority — Uganda' }}
          </div>
          <h1 class="hero-title text-4xl md:text-6xl font-bold mb-6 leading-tight text-white">
            {{ displaySlides[currentSlide]?.title || 'Building a World-Class' }}
            <span v-if="!displaySlides[currentSlide]?.title" class="text-primary-400"> Sports Nation</span>
          </h1>
          <p class="hero-desc text-lg text-gray-300 mb-10 max-w-2xl">
            {{ displaySlides[currentSlide]?.description || 'The National Council of Sports registers, licenses and regulates sports organisations, federations and clubs across Uganda.' }}
          </p>
          <div class="hero-btns flex flex-wrap gap-4">
            <router-link
              :to="displaySlides[currentSlide]?.button_url || '/apply'"
              class="bg-primary-600 hover:bg-primary-500 text-white font-semibold px-8 py-3 rounded-lg transition-colors text-base"
            >
              {{ displaySlides[currentSlide]?.button_text || 'Apply for License' }}
            </router-link>
            <router-link to="/news" class="border border-gray-500 hover:border-gray-300 text-gray-300 hover:text-white font-semibold px-8 py-3 rounded-lg transition-colors text-base">
              Latest News
            </router-link>
          </div>
        </div>
      </div>

      <!-- Slideshow Controls -->
      <div v-if="displaySlides.length > 1" class="absolute bottom-8 left-1/2 -translate-x-1/2 z-20 flex items-center gap-3">
        <button @click="prevSlide" class="p-2 rounded-full bg-white/20 hover:bg-white/40 text-white transition-colors">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
        </button>
        <div class="flex gap-2">
          <button
            v-for="(_, i) in displaySlides"
            :key="i"
            @click="goToSlide(i)"
            class="transition-all rounded-full"
            :class="currentSlide === i ? 'w-6 h-2.5 bg-primary-400' : 'w-2.5 h-2.5 bg-white/40 hover:bg-white/70'"
          ></button>
        </div>
        <button @click="nextSlide" class="p-2 rounded-full bg-white/20 hover:bg-white/40 text-white transition-colors">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
        </button>
      </div>
    </section>

    <!-- Quick Links -->
    <section class="bg-primary-600 py-6">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <router-link to="/apply" class="quick-link flex items-center gap-3 text-white hover:bg-primary-700 rounded-lg px-4 py-3 transition-colors">
            <div class="w-8 h-8 bg-white/20 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/></svg>
            </div>
            <span class="text-sm font-medium">New Application</span>
          </router-link>
          <router-link to="/apply" class="quick-link flex items-center gap-3 text-white hover:bg-primary-700 rounded-lg px-4 py-3 transition-colors">
            <div class="w-8 h-8 bg-white/20 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg>
            </div>
            <span class="text-sm font-medium">Renew License</span>
          </router-link>
          <router-link to="/events" class="quick-link flex items-center gap-3 text-white hover:bg-primary-700 rounded-lg px-4 py-3 transition-colors">
            <div class="w-8 h-8 bg-white/20 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
            </div>
            <span class="text-sm font-medium">Upcoming Events</span>
          </router-link>
          <router-link to="/careers" class="quick-link flex items-center gap-3 text-white hover:bg-primary-700 rounded-lg px-4 py-3 transition-colors">
            <div class="w-8 h-8 bg-white/20 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/></svg>
            </div>
            <span class="text-sm font-medium">Careers</span>
          </router-link>
        </div>
      </div>
    </section>

    <!-- About Section -->
    <section ref="aboutRef" class="py-20 px-4 bg-white">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <div class="grid md:grid-cols-2 gap-12 items-center">
          <div class="about-text">
            <p class="text-primary-600 font-semibold text-sm uppercase tracking-wide mb-3">About NCS</p>
            <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-6">Regulating Sports Excellence in Uganda</h2>
            <p class="text-gray-600 mb-4">
              The National Council of Sports (NCS) is Uganda's national sports authority mandated by the National Council of Sports Act.
              We oversee the registration and regulation of all sports organisations, federations, associations, and clubs across the country.
            </p>
            <p class="text-gray-600 mb-8">
              Our mission is to promote, develop, coordinate and regulate sports and physical activities for the benefit of all Ugandans.
            </p>
            <router-link to="/apply" class="inline-flex items-center gap-2 bg-primary-600 hover:bg-primary-700 text-white font-semibold px-6 py-3 rounded-lg transition-colors">
              Apply Online Today
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"/></svg>
            </router-link>
          </div>

          <div class="grid grid-cols-2 gap-4 about-stats">
            <div
              v-for="(stat, i) in displayedStats"
              :key="stat.label || i"
              class="stat-card rounded-2xl p-6 text-center"
              :class="i % 2 === 0 ? 'bg-gray-50' : 'bg-primary-600 text-white'"
            >
              <div
                class="stat-num text-4xl font-bold mb-1"
                :class="i % 2 === 0 ? 'text-primary-700' : ''"
                :data-target="parseStatValue(stat.value)"
              >{{ stat.value }}</div>
              <div
                class="text-sm font-medium"
                :class="i % 2 === 0 ? 'text-gray-500' : 'opacity-80'"
              >{{ stat.label }}</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Latest News -->
    <section ref="newsRef" class="py-20 bg-gray-50 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <div class="flex justify-between items-center mb-10">
          <div>
            <p class="text-primary-600 font-semibold text-sm uppercase tracking-wide mb-1">Stay Informed</p>
            <h2 class="text-3xl font-bold text-gray-900">Latest News</h2>
          </div>
          <router-link to="/news" class="text-sm font-medium text-primary-600 hover:text-primary-800 flex items-center gap-1">
            View all <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
          </router-link>
        </div>
        <div v-if="postsLoading" class="grid md:grid-cols-3 gap-6">
          <div v-for="i in 3" :key="i" class="animate-pulse bg-white rounded-2xl h-72"></div>
        </div>
        <div v-else class="grid md:grid-cols-3 gap-6">
          <router-link
            v-for="post in posts"
            :key="post.id"
            :to="`/news/${post.slug}`"
            class="news-card group bg-white rounded-2xl shadow-sm overflow-hidden hover:shadow-lg transition-all hover:-translate-y-1"
          >
            <div v-if="post.cover_image_url" class="h-44 overflow-hidden">
              <img :src="post.cover_image_url" :alt="post.title" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
            </div>
            <div v-else class="h-44 bg-gradient-to-br from-primary-100 to-primary-200 flex items-center justify-center">
              <svg class="w-12 h-12 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"/></svg>
            </div>
            <div class="p-5">
              <div class="flex items-center gap-2 mb-2">
                <span class="text-xs font-semibold text-primary-600 uppercase tracking-wide">{{ post.category }}</span>
              </div>
              <h3 class="font-bold text-gray-900 mb-2 line-clamp-2 group-hover:text-primary-700 transition-colors">{{ post.title }}</h3>
              <p v-if="post.excerpt" class="text-sm text-gray-500 line-clamp-2">{{ post.excerpt }}</p>
            </div>
          </router-link>
          <div v-if="!posts.length" class="col-span-3 text-center py-16 text-gray-400">
            <p>No news articles yet.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Events -->
    <section ref="eventsRef" class="py-20 bg-white px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <div class="flex justify-between items-center mb-10">
          <div>
            <p class="text-primary-600 font-semibold text-sm uppercase tracking-wide mb-1">What's Happening</p>
            <h2 class="text-3xl font-bold text-gray-900">Upcoming Events</h2>
          </div>
          <router-link to="/events" class="text-sm font-medium text-primary-600 hover:text-primary-800 flex items-center gap-1">
            All events <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
          </router-link>
        </div>
        <div class="space-y-4">
          <router-link
            v-for="ev in events"
            :key="ev.id"
            :to="`/events/${ev.slug}`"
            class="event-card group flex items-start gap-5 p-5 rounded-2xl border border-gray-100 hover:border-primary-200 hover:bg-primary-50/30 transition-all"
          >
            <div class="flex-shrink-0 w-14 h-14 bg-primary-600 text-white rounded-xl flex flex-col items-center justify-center text-center">
              <span class="text-xs font-semibold uppercase leading-none">{{ monthShort(ev.event_date) }}</span>
              <span class="text-xl font-bold leading-tight">{{ new Date(ev.event_date).getDate() }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="font-semibold text-gray-900 group-hover:text-primary-700 transition-colors truncate">{{ ev.title }}</h3>
              <p v-if="ev.location" class="text-sm text-gray-500 mt-0.5 flex items-center gap-1">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
                {{ ev.location }}
              </p>
            </div>
            <svg class="w-5 h-5 text-gray-300 group-hover:text-primary-500 flex-shrink-0 mt-1 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
          </router-link>
          <div v-if="!events.length" class="text-center py-8 text-gray-400">No upcoming events.</div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section ref="ctaRef" class="bg-gradient-to-br from-primary-700 to-primary-900 py-20 px-4 text-white">
      <div class="max-w-4xl mx-auto text-center">
        <h2 class="cta-title text-3xl md:text-4xl font-bold mb-4">Ready to Register Your Sports Organisation?</h2>
        <p class="cta-desc text-lg text-primary-200 mb-8 max-w-2xl mx-auto">
          Apply online for your sports license, register your federation, club, or association with the National Council of Sports Uganda.
        </p>
        <div class="flex flex-wrap gap-4 justify-center">
          <router-link to="/apply" class="bg-white text-primary-700 hover:bg-primary-50 font-bold px-8 py-3 rounded-lg transition-colors text-base">
            Start Application
          </router-link>
          <router-link
            v-if="!isAuthenticated"
            to="/login"
            class="border border-primary-400 hover:border-white text-white font-semibold px-8 py-3 rounded-lg transition-colors text-base"
          >
            Sign In
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { gsap } from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { listPosts, listEvents, listSlides, listFunFacts } from '@/api/cms.js'
import { useAuthStore } from '@/stores/auth.js'

gsap.registerPlugin(ScrollTrigger)

const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)

const posts = ref([])
const events = ref([])
const slides = ref([])
const funFacts = ref([])
const postsLoading = ref(true)

const defaultStats = [
  { label: 'Sports Federations', value: '50' },
  { label: 'Registered Clubs',   value: '500' },
  { label: 'Licensed Athletes',  value: '10,000' },
  { label: 'Districts Covered',  value: '112' },
]
const displayedStats = computed(() =>
  funFacts.value.length > 0
    ? funFacts.value.slice(0, 4).map(f => ({ label: f.label, value: f.value }))
    : defaultStats
)
function parseStatValue(v) {
  if (!v) return 0
  const n = parseInt(String(v).replace(/[^0-9]/g, ''), 10)
  return isNaN(n) ? 0 : n
}

// ── Slideshow ─────────────────────────────────────────────────────
const currentSlide = ref(0)
let slideTimer = null

const defaultSlide = {
  id: 'default',
  title: 'Building a World-Class Sports Nation',
  subtitle: 'Official Sports Regulatory Authority — Uganda',
  description: 'The National Council of Sports registers, licenses and regulates sports organisations, federations and clubs across Uganda. Apply online for your sports license today.',
  button_text: 'Apply for License',
  button_url: '/apply',
  image_url: '',
}

const displaySlides = computed(() => slides.value.length > 0 ? slides.value : [defaultSlide])

function nextSlide() {
  currentSlide.value = (currentSlide.value + 1) % displaySlides.value.length
  restartTimer()
}
function prevSlide() {
  currentSlide.value = (currentSlide.value - 1 + displaySlides.value.length) % displaySlides.value.length
  restartTimer()
}
function goToSlide(i) {
  currentSlide.value = i
  restartTimer()
}
function restartTimer() {
  clearInterval(slideTimer)
  if (displaySlides.value.length > 1) {
    slideTimer = setInterval(nextSlide, 6000)
  }
}

// ── Date helpers ──────────────────────────────────────────────────
function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { day: 'numeric', month: 'short', year: 'numeric' })
}
function monthShort(d) {
  if (!d) return '?'
  return new Date(d).toLocaleDateString('en-UG', { month: 'short' })
}

// ── GSAP refs ─────────────────────────────────────────────────────
const heroRef = ref(null)
const heroContent = ref(null)
const aboutRef = ref(null)
const newsRef = ref(null)
const eventsRef = ref(null)
const ctaRef = ref(null)

function initAnimations() {
  // Hero entrance
  const tl = gsap.timeline({ defaults: { ease: 'power3.out' } })
  tl.fromTo('.hero-badge',   { opacity: 0, y: 30 }, { opacity: 1, y: 0, duration: 0.7 })
    .fromTo('.hero-title',   { opacity: 0, y: 40 }, { opacity: 1, y: 0, duration: 0.8 }, '-=0.4')
    .fromTo('.hero-desc',    { opacity: 0, y: 30 }, { opacity: 1, y: 0, duration: 0.7 }, '-=0.5')
    .fromTo('.hero-btns',    { opacity: 0, y: 20 }, { opacity: 1, y: 0, duration: 0.6 }, '-=0.4')

  // Quick links stagger
  gsap.fromTo('.quick-link',
    { opacity: 0, y: 20 },
    { opacity: 1, y: 0, duration: 0.5, stagger: 0.1, ease: 'power2.out',
      scrollTrigger: { trigger: '.quick-link', start: 'top 90%' } }
  )

  // About text slide in
  gsap.fromTo('.about-text',
    { opacity: 0, x: -50 },
    { opacity: 1, x: 0, duration: 0.9, ease: 'power3.out',
      scrollTrigger: { trigger: aboutRef.value, start: 'top 75%' } }
  )

  // Stats cards stagger
  gsap.fromTo('.stat-card',
    { opacity: 0, scale: 0.85 },
    { opacity: 1, scale: 1, duration: 0.6, stagger: 0.12, ease: 'back.out(1.5)',
      scrollTrigger: { trigger: '.about-stats', start: 'top 75%',
        onEnter: animateCounters } }
  )

  // News cards stagger
  gsap.fromTo('.news-card',
    { opacity: 0, y: 40 },
    { opacity: 1, y: 0, duration: 0.7, stagger: 0.15, ease: 'power2.out',
      scrollTrigger: { trigger: newsRef.value, start: 'top 75%' } }
  )

  // Event cards
  gsap.fromTo('.event-card',
    { opacity: 0, x: -30 },
    { opacity: 1, x: 0, duration: 0.6, stagger: 0.1, ease: 'power2.out',
      scrollTrigger: { trigger: eventsRef.value, start: 'top 75%' } }
  )

  // CTA
  gsap.fromTo('.cta-title, .cta-desc',
    { opacity: 0, y: 30 },
    { opacity: 1, y: 0, duration: 0.8, stagger: 0.2, ease: 'power2.out',
      scrollTrigger: { trigger: ctaRef.value, start: 'top 80%' } }
  )
}

function animateCounters() {
  document.querySelectorAll('.stat-num').forEach(el => {
    const target = parseInt(el.dataset.target || '0', 10)
    const suffix = target >= 10000 ? 'K+' : target >= 100 ? '+' : ''
    const displayTarget = target >= 10000 ? target / 1000 : target
    gsap.fromTo({ val: 0 }, { val: displayTarget, duration: 1.8, ease: 'power2.out',
      onUpdate: function () {
        el.textContent = Math.floor(this.targets()[0].val) + suffix
      }
    })
  })
}

onMounted(async () => {
  // Load data
  try {
    const [pr, er, sr, fr] = await Promise.all([
      listPosts({ status: 'published', per_page: 3 }),
      listEvents({ status: 'published', per_page: 4 }),
      listSlides(),
      listFunFacts()
    ])
    posts.value = pr.data.data?.items || []
    events.value = er.data.data?.items || []
    slides.value = sr.data.data || []
    funFacts.value = fr.data.data || []
  } catch {
    // graceful fallback
  } finally {
    postsLoading.value = false
  }

  // Start slideshow
  restartTimer()

  // Init GSAP after DOM is ready
  setTimeout(initAnimations, 100)
})

onUnmounted(() => {
  clearInterval(slideTimer)
  ScrollTrigger.getAll().forEach(t => t.kill())
})
</script>

<style scoped>
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: opacity 0.8s ease;
}
.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
}
.line-clamp-2 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
</style>
