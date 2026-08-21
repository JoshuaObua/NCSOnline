<template>
  <div id="section-hero" data-testid="homepage-section-hero">
    <section
      class="relative h-[500px] md:h-[600px] lg:h-[700px] overflow-hidden"
      data-testid="homepage-hero-slider"
      role="region"
      aria-roledescription="carousel"
      aria-label="NCS highlights"
      tabindex="0"
      @keydown.left.prevent="previous"
      @keydown.right.prevent="next"
      @mouseenter="pauseIfEnabled"
      @mouseleave="resume"
      @focusin="pauseIfEnabled"
      @focusout="resume"
    >
      <!-- Slides Loop -->
      <div
        v-for="(slide, index) in normalizedSlides"
        :key="slide.id || index"
        class="absolute inset-0 transition-all duration-700 ease-in-out"
        :class="index === current ? 'opacity-100 scale-100 z-10' : 'opacity-0 scale-105 z-0 pointer-events-none'"
        :aria-hidden="index === current ? 'false' : 'true'"
      >
        <!-- Background Media -->
        <video
          v-if="slide.media_type === 'video' && slide.image_url"
          class="absolute inset-0 w-full h-full object-cover"
          :src="mediaUrl(slide.image_url)"
          autoplay
          muted
          loop
          playsinline
          preload="auto"
        ></video>
        <div
          v-else-if="slide.image_url"
          class="absolute inset-0 bg-cover bg-center"
          :style="{ backgroundImage: `url(${mediaUrl(slide.image_url)})` }"
        ></div>
        <div v-else class="absolute inset-0 bg-gradient-to-br from-slate-900 via-slate-800 to-[#1a365d]"></div>

        <!-- Overlay Gradient -->
        <div class="absolute inset-0 bg-gradient-to-r from-[#1a365d]/90 via-[#1a365d]/70 to-transparent"></div>

        <!-- Slide Content -->
        <div class="relative z-20 h-full max-w-7xl mx-auto px-4 flex items-center">
          <div
            class="max-w-2xl transition-all duration-700 delay-200"
            :class="index === current ? 'translate-y-0 opacity-100' : 'translate-y-10 opacity-0'"
          >
            <span class="inline-block px-4 py-1.5 bg-[#f5a623] text-white text-sm font-semibold rounded-full mb-4">
              {{ slide.subtitle || 'National Council of Sports' }}
            </span>
            <h1 class="text-3xl md:text-5xl lg:text-6xl font-bold text-white mb-4 leading-tight">
              {{ slide.title || 'Welcome to National Council of Sports' }}
            </h1>
            <p class="text-lg md:text-xl text-white/90 mb-8">
              {{ slide.description || 'A Centre of Excellence for Promotion and Development of Sports' }}
            </p>
            <div class="flex flex-wrap gap-4">
              <component
                :is="isInternal(button.url) ? 'router-link' : 'a'"
                v-for="button in slide.buttons"
                :key="`${slide.id}-${button.sort_order}-${button.text}`"
                :to="isInternal(button.url) ? button.url : undefined"
                :href="isInternal(button.url) ? undefined : button.url"
                :target="button.link_target"
                :rel="button.link_target === '_blank' ? 'noopener noreferrer' : undefined"
                :class="[
                  'inline-flex items-center justify-center gap-2 whitespace-nowrap px-8 py-4 text-base md:text-lg font-semibold rounded-lg transition-all duration-300 hover:-translate-y-1',
                  button.style_class === 'primary' || button.style_class === 'secondary'
                    ? 'bg-[#f5a623] hover:bg-[#e09612] text-white shadow-lg hover:shadow-xl'
                    : 'border-2 border-white text-white hover:bg-white hover:text-[#1a365d] bg-transparent shadow-sm'
                ]"
                :data-testid="button.style_class === 'primary' ? 'hero-primary-button' : 'hero-secondary-button'"
              >
                <svg
                  v-if="button.style_class === 'outline' || button.text.toLowerCase().includes('video')"
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="w-5 h-5 mr-1 flex-shrink-0"
                  aria-hidden="true"
                >
                  <polygon points="6 3 20 12 6 21 6 3"></polygon>
                </svg>
                {{ button.text }}
              </component>
            </div>
          </div>
        </div>
      </div>

      <!-- Arrow Controls -->
      <button
        v-if="normalizedSlides.length > 1"
        type="button"
        class="absolute left-4 md:left-8 top-1/2 -translate-y-1/2 z-30 w-12 h-12 md:w-14 md:h-14 rounded-full bg-white/20 backdrop-blur-sm border border-white/30 flex items-center justify-center text-white hover:bg-[#f5a623] hover:border-[#f5a623] transition-all duration-300 group"
        aria-label="Previous slide"
        data-testid="hero-previous-slide-button"
        @click="previous"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 group-hover:scale-110 transition-transform" aria-hidden="true">
          <path d="m15 18-6-6 6-6"></path>
        </svg>
      </button>

      <button
        v-if="normalizedSlides.length > 1"
        type="button"
        class="absolute right-4 md:right-8 top-1/2 -translate-y-1/2 z-30 w-12 h-12 md:w-14 md:h-14 rounded-full bg-white/20 backdrop-blur-sm border border-white/30 flex items-center justify-center text-white hover:bg-[#f5a623] hover:border-[#f5a623] transition-all duration-300 group"
        aria-label="Next slide"
        data-testid="hero-next-slide-button"
        @click="next"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 group-hover:scale-110 transition-transform" aria-hidden="true">
          <path d="m9 18 6-6-6-6"></path>
        </svg>
      </button>

      <!-- Page Indicator / Dynamic Slide Indicators -->
      <div v-if="normalizedSlides.length > 1" class="absolute bottom-8 left-1/2 -translate-x-1/2 z-30 flex gap-3">
        <button
          v-for="(_, index) in normalizedSlides"
          :key="index"
          type="button"
          :class="[
            'transition-all duration-300 rounded-full h-3',
            index === current ? 'w-10 bg-[#f5a623]' : 'w-3 bg-white/50 hover:bg-white'
          ]"
          :aria-label="`Go to slide ${index + 1}`"
          :data-testid="`hero-slide-indicator-${index + 1}`"
          @click="go(index)"
        ></button>
      </div>

      <!-- Bottom Gradient Overlay -->
      <div class="absolute bottom-0 left-0 right-0 h-24 bg-gradient-to-t from-white dark:from-slate-950 to-transparent z-20 pointer-events-none"></div>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { mediaUrl } from '@/api/client.js'

const props = defineProps({
  slideshow: { type: Object, default: () => ({}) },
  slides: { type: Array, default: () => [] }
})

const current = ref(0)
const paused = ref(false)
let timer = null

const fallbackSlide = {
  id: 'default',
  title: 'Welcome to National Council of Sports',
  subtitle: 'National Council of Sports',
  description: 'A Centre of Excellence for Promotion and Development of Sports',
  button_text: 'Learn More',
  button_url: '/about/mandate',
  media_type: 'image',
  image_url: 'https://images.unsplash.com/photo-1582788531956-e0556b4390e5?w=1920&h=600&fit=crop',
}

const normalizedSlides = computed(() => {
  const source = Array.isArray(props.slides) && props.slides.length ? props.slides : [fallbackSlide]
  return source.map(slide => ({
    ...slide,
    media_type: slide.media_type || 'image',
    buttons: normalizeButtons(slide),
  }))
})

const interval = computed(() => Number(props.slideshow?.autoplay_speed || 6500))
const pauseOnHover = computed(() => props.slideshow?.pause_on_hover !== false)

watch(normalizedSlides, () => {
  if (current.value >= normalizedSlides.value.length) current.value = 0
  restart()
})

onMounted(restart)
onUnmounted(() => {
  clearInterval(timer)
})

function normalizeButtons(slide) {
  const buttons = Array.isArray(slide.buttons) && slide.buttons.length
    ? slide.buttons
    : [
        { text: slide.button_text || 'Learn More', url: slide.button_url || '/about/mandate', style_class: 'primary', link_target: '_self', sort_order: 0 },
        { text: 'Watch Video', url: '/news', style_class: 'outline', link_target: '_self', sort_order: 1 },
      ]
  return buttons.slice(0, 2).map((button, index) => ({
    text: button.text || (index ? 'Watch Video' : 'Learn More'),
    url: button.url || '/',
    style_class: ['primary', 'secondary', 'outline'].includes(button.style_class) ? button.style_class : (index ? 'outline' : 'primary'),
    link_target: button.link_target || '_self',
    sort_order: button.sort_order ?? index,
  }))
}

function next() {
  go((current.value + 1) % normalizedSlides.value.length)
}

function previous() {
  go((current.value - 1 + normalizedSlides.value.length) % normalizedSlides.value.length)
}

function go(index) {
  current.value = index
  restart()
}

function pauseIfEnabled() { if (pauseOnHover.value) pause() }
function pause() { paused.value = true; clearInterval(timer) }
function resume() { paused.value = false; restart() }

function restart() {
  clearInterval(timer)
  if (!paused.value && normalizedSlides.value.length > 1) {
    timer = setInterval(next, interval.value)
  }
}

function isInternal(url = '') {
  return String(url).startsWith('/')
}
</script>
