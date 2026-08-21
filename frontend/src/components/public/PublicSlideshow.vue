<template>
  <section
    class="public-carousel"
    :style="{ '--slide-duration': `${duration}ms` }"
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
    <article
      v-for="(slide, index) in normalizedSlides"
      :key="slide.id || index"
      class="carousel-slide"
      :class="[transitionClass(index), slide.animation_type]"
      :aria-hidden="index === current ? 'false' : 'true'"
    >
      <video
        v-if="slide.media_type === 'video' && slide.image_url"
        class="carousel-media"
        :src="mediaUrl(slide.image_url)"
        autoplay
        muted
        loop
        playsinline
        preload="auto"
      ></video>
      <img
        v-else-if="slide.image_url"
        class="carousel-media"
        :src="mediaUrl(slide.image_url)"
        :alt="slide.title"
        :loading="index <= 1 ? 'eager' : 'lazy'"
      />
      <div v-else class="carousel-fallback"></div>
      <div class="carousel-overlay"></div>
      <div class="carousel-content">
        <div>
          <span>{{ slide.subtitle || 'National Council of Sports' }}</span>
          <h1>{{ slide.title || 'Welcome to National Council of Sports' }}</h1>
          <p>{{ slide.description || 'A Centre of Excellence for Promotion and Development of Sports' }}</p>
          <div class="carousel-actions">
            <component
              :is="isInternal(button.url) ? 'router-link' : 'a'"
              v-for="button in slide.buttons"
              :key="`${slide.id}-${button.sort_order}-${button.text}`"
              :to="isInternal(button.url) ? button.url : undefined"
              :href="isInternal(button.url) ? undefined : button.url"
              :target="button.link_target"
              :rel="button.link_target === '_blank' ? 'noopener noreferrer' : undefined"
              :class="['carousel-btn', `carousel-btn-${button.style_class}`]"
            >
              <svg v-if="button.style_class === 'outline'" class="carousel-btn-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polygon points="6 3 20 12 6 21 6 3"></polygon></svg>
              {{ button.text }}
            </component>
          </div>
        </div>
      </div>
    </article>

    <button v-if="normalizedSlides.length > 1" type="button" class="carousel-arrow prev" aria-label="Previous slide" @click="previous">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m15 18-6-6 6-6"></path></svg>
    </button>
    <button v-if="normalizedSlides.length > 1" type="button" class="carousel-arrow next" aria-label="Next slide" @click="next">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m9 18 6-6-6-6"></path></svg>
    </button>
    <div v-if="normalizedSlides.length > 1" class="carousel-dots">
      <button
        v-for="(_, index) in normalizedSlides"
        :key="index"
        type="button"
        :class="{ active: index === current }"
        :aria-label="`Show slide ${index + 1}`"
        @click="go(index)"
      ></button>
    </div>
    <div class="carousel-bottom"></div>
  </section>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { mediaUrl } from '@/api/client.js'

const props = defineProps({
  slideshow: { type: Object, default: () => ({}) },
  slides: { type: Array, default: () => [] }
})

const current = ref(0)
const previousIndex = ref(null)
const slideDirection = ref('next')
const paused = ref(false)
let timer = null
let transitionTimer = null

const fallbackSlide = {
  id: 'default',
  title: 'Welcome to National Council of Sports',
  subtitle: 'National Council of Sports',
  description: 'A Centre of Excellence for Promotion and Development of Sports',
  button_text: 'Discover NCS',
  button_url: '/pages/the-mandate',
  animation_type: 'fade-in',
  image_url: '',
}

const normalizedSlides = computed(() => {
  const source = Array.isArray(props.slides) && props.slides.length ? props.slides : [fallbackSlide]
  return source.map(slide => ({
    ...slide,
    media_type: slide.media_type || 'image',
    animation_type: animationMode(slide.animation_type),
    buttons: normalizeButtons(slide),
  }))
})

const interval = computed(() => Number(props.slideshow?.autoplay_speed || 6500))
const duration = computed(() => Number(props.slideshow?.transition_duration || 750))
const effect = computed(() => ['fade', 'slide', 'vertical', 'zoom'].includes(props.slideshow?.transition_effect) ? props.slideshow.transition_effect : 'fade')
const pauseOnHover = computed(() => props.slideshow?.pause_on_hover !== false)

watch(normalizedSlides, () => {
  if (current.value >= normalizedSlides.value.length) current.value = 0
  restart()
})

onMounted(restart)
onUnmounted(() => {
  clearInterval(timer)
  clearTimeout(transitionTimer)
})

function normalizeButtons(slide) {
  const buttons = Array.isArray(slide.buttons) && slide.buttons.length
    ? slide.buttons
    : [
        { text: slide.button_text || 'Learn More', url: slide.button_url || '/pages/the-mandate', style_class: 'primary', link_target: '_self', sort_order: 0 },
        { text: 'Watch Video', url: '/news', style_class: 'outline', link_target: '_self', sort_order: 1 },
      ]
  return buttons.slice(0, 2).map((button, index) => ({
    text: button.text || (index ? 'Read More' : 'Learn More'),
    url: button.url || '/',
    style_class: ['primary', 'secondary', 'outline'].includes(button.style_class) ? button.style_class : (index ? 'outline' : 'primary'),
    link_target: button.link_target || '_self',
    sort_order: button.sort_order ?? index,
  }))
}

function transitionClass(index) {
  const mode = effect.value
  if (index === current.value) return `active mode-${mode} ${slideDirection.value}`
  if (index === previousIndex.value) return `outgoing mode-${mode} ${slideDirection.value}`
  return 'hidden-slide'
}

function next() {
  slideDirection.value = 'next'
  go((current.value + 1) % normalizedSlides.value.length)
}

function previous() {
  slideDirection.value = 'prev'
  go((current.value - 1 + normalizedSlides.value.length) % normalizedSlides.value.length)
}

function go(index) {
  if (index === current.value) return
  previousIndex.value = current.value
  current.value = index
  restart()

  clearTimeout(transitionTimer)
  transitionTimer = setTimeout(() => {
    previousIndex.value = null
  }, duration.value + 60)
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

function animationMode(value) {
  return ['fade-in', 'slide-up', 'slide-left', 'slide-right', 'zoom-in', 'zoom-out', 'flip-in', 'blur-in', 'bounce-in', 'ken-burns'].includes(value) ? value : 'fade-in'
}
</script>

<style scoped>
.public-carousel {
  position: relative;
  height: 31.25rem;
  overflow: hidden;
  background: #0f172a;
  outline: none;
}
@media (min-width: 768px) { .public-carousel { height: 37.5rem; } }
@media (min-width: 1024px) { .public-carousel { height: 43.75rem; } }

/* Base Slide Positioning */
.carousel-slide {
  position: absolute;
  inset: 0;
  opacity: 0;
  z-index: 0;
  pointer-events: none;
  transition: opacity var(--slide-duration, 750ms) ease-in-out,
              transform var(--slide-duration, 750ms) ease-in-out,
              filter var(--slide-duration, 750ms) ease;
  will-change: opacity, transform, filter;
}

.carousel-slide.hidden-slide {
  visibility: hidden;
  opacity: 0;
  z-index: 0;
}

/*
 * OUTGOING SLIDE: Stays 100% visible at z-index 1 underneath until cross-fade finishes.
 * This guarantees zero dark flash or blank background gap!
 */
.carousel-slide.outgoing {
  z-index: 1;
  opacity: 1;
  visibility: visible;
  transform: translate3d(0, 0, 0) scale(1);
}

/* OUTGOING Slide Transitions */
.carousel-slide.outgoing.mode-slide.next { animation: fadeOutSlightly var(--slide-duration, 750ms) ease-in-out forwards; }
.carousel-slide.outgoing.mode-slide.prev { animation: fadeOutSlightly var(--slide-duration, 750ms) ease-in-out forwards; }
.carousel-slide.outgoing.mode-vertical.next { animation: fadeOutSlightly var(--slide-duration, 750ms) ease-in-out forwards; }
.carousel-slide.outgoing.mode-vertical.prev { animation: fadeOutSlightly var(--slide-duration, 750ms) ease-in-out forwards; }

/*
 * INCOMING ACTIVE SLIDE: Renders on top at z-index 2 with rich transition modes (fade, slide, vertical, zoom).
 */
.carousel-slide.active {
  z-index: 2;
  opacity: 1;
  visibility: visible;
  pointer-events: auto;
  transform: translate3d(0, 0, 0) scale(1);
  filter: none;
}

/* Mode-Specific Transitions for Incoming Active Slide */
.carousel-slide.active.mode-fade {
  animation: crossFadeIn var(--slide-duration, 750ms) ease-in-out forwards;
}

.carousel-slide.active.mode-slide.next {
  animation: slideInRight var(--slide-duration, 750ms) cubic-bezier(0.25, 1, 0.5, 1) forwards;
}
.carousel-slide.active.mode-slide.prev {
  animation: slideInLeft var(--slide-duration, 750ms) cubic-bezier(0.25, 1, 0.5, 1) forwards;
}

.carousel-slide.active.mode-vertical.next {
  animation: slideInDown var(--slide-duration, 750ms) cubic-bezier(0.25, 1, 0.5, 1) forwards;
}
.carousel-slide.active.mode-vertical.prev {
  animation: slideInUp var(--slide-duration, 750ms) cubic-bezier(0.25, 1, 0.5, 1) forwards;
}

.carousel-slide.active.mode-zoom {
  animation: zoomFadeIn var(--slide-duration, 750ms) cubic-bezier(0.25, 1, 0.5, 1) forwards;
}

/* Media & Overlay Styling */
.carousel-media, .carousel-fallback {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.carousel-fallback {
  background: linear-gradient(120deg, #1e293b, #0f172a);
}

.carousel-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, rgba(15, 23, 42, 0.88), rgba(15, 23, 42, 0.65), transparent);
}

.carousel-content {
  position: absolute;
  z-index: 3;
  inset: 0;
  display: flex;
  align-items: center;
  width: min(80rem, 100%);
  margin: auto;
  padding: 0 1rem;
  color: white;
}
.carousel-content > div { max-width: 42rem; }
.carousel-content span {
  display: inline-block;
  margin-bottom: 1rem;
  border-radius: 999px;
  background: #f5a623;
  padding: .375rem 1rem;
  font-size: .875rem;
  font-weight: 600;
}
.carousel-content h1 {
  font-size: 1.875rem;
  font-weight: 700;
  line-height: 1.25;
  margin-bottom: 1rem;
}
@media (min-width: 768px) { .carousel-content h1 { font-size: 3rem; } }
@media (min-width: 1024px) { .carousel-content h1 { font-size: 3.75rem; } }

.carousel-content p {
  margin: 0 0 2rem;
  font-size: 1.125rem;
  line-height: 1.65;
  color: rgba(255, 255, 255, 0.9);
}
@media (min-width: 768px) { .carousel-content p { font-size: 1.25rem; } }

/*
 * RICH PER-SLIDE CONTENT ANIMATIONS (fade-in, slide-up, slide-left, slide-right, zoom-in, zoom-out, flip-in, blur-in, bounce-in, ken-burns)
 */
.carousel-slide.active.fade-in .carousel-content > div { animation: fadeIn .72s ease both; }
.carousel-slide.active.slide-up .carousel-content > div { animation: slideUp .72s ease both; }
.carousel-slide.active.slide-left .carousel-content > div { animation: slideLeft .72s ease both; }
.carousel-slide.active.slide-right .carousel-content > div { animation: slideRight .72s ease both; }
.carousel-slide.active.zoom-in .carousel-content > div { animation: zoomIn .72s ease both; }
.carousel-slide.active.zoom-out .carousel-content > div { animation: zoomOut .72s ease both; }
.carousel-slide.active.flip-in .carousel-content > div { animation: flipIn .78s ease both; transform-origin: left center; }
.carousel-slide.active.blur-in .carousel-content > div { animation: blurIn .78s ease both; }
.carousel-slide.active.bounce-in .carousel-content > div { animation: bounceIn .82s cubic-bezier(.2, .85, .35, 1.2) both; }
.carousel-slide.active.ken-burns .carousel-media,
.carousel-slide.active.ken-burns .carousel-fallback { animation: kenBurns 7s ease both; }

/* Buttons & Navigation Styling */
.carousel-actions { display: flex; flex-wrap: wrap; gap: 1rem; }
.carousel-btn {
  display: inline-flex;
  min-height: 44px;
  align-items: center;
  justify-content: center;
  gap: .5rem;
  border-radius: .5rem;
  padding: 1rem 1.75rem;
  font-size: 1.125rem;
  font-weight: 600;
  transition: transform .3s, background .3s, color .3s, box-shadow .3s;
}
.carousel-btn:hover { transform: translateY(-.25rem); }
.carousel-btn-icon { width: 1.25rem; height: 1.25rem; flex-shrink: 0; }
.carousel-btn-primary {
  background: #f5a623;
  color: white;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, .1), 0 4px 6px -4px rgba(0, 0, 0, .1);
}
.carousel-btn-primary:hover {
  background: #e09612;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, .1), 0 8px 10px -6px rgba(0, 0, 0, .1);
}
.carousel-btn-secondary { background: #1e293b; color: white; }
.carousel-btn-outline { border: 2px solid white; color: white; background: transparent; }
.carousel-btn-outline:hover { background: white; color: #0f172a; }

.carousel-arrow {
  position: absolute;
  z-index: 4;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 3rem;
  height: 3rem;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  backdrop-filter: blur(4px);
  transition: background .3s, border-color .3s;
}
@media (min-width: 768px) { .carousel-arrow { width: 3.5rem; height: 3.5rem; } }
.carousel-arrow svg { width: 1.5rem; height: 1.5rem; transition: transform .2s; }
.carousel-arrow:hover svg { transform: scale(1.1); }
.carousel-arrow:hover, .carousel-arrow:focus-visible { background: #f5a623; border-color: #f5a623; }
.prev { left: 1rem; }
.next { right: 1rem; }
@media (min-width: 768px) { .prev { left: 2rem; } .next { right: 2rem; } }

.carousel-dots {
  position: absolute;
  z-index: 4;
  bottom: 2rem;
  left: 50%;
  display: flex;
  gap: .75rem;
  transform: translateX(-50%);
}
.carousel-dots button {
  width: .75rem;
  height: .75rem;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.5);
  transition: width .3s, background .3s;
}
.carousel-dots button:hover { background: white; }
.carousel-dots button.active { width: 2.5rem; background: #f5a623; }

.carousel-bottom {
  position: absolute;
  z-index: 3;
  right: 0;
  bottom: -1px;
  left: 0;
  height: 7.5rem;
  background: linear-gradient(to top, #ffffff 0%, rgba(255, 255, 255, 0.88) 35%, rgba(255, 255, 255, 0.35) 70%, transparent 100%);
  pointer-events: none;
}

/* Keyframes for Slide Cross-Fading & Motion Transitions */
@keyframes crossFadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes fadeOutSlightly {
  from { opacity: 1; transform: scale(1); }
  to { opacity: 0.9; transform: scale(0.98); }
}

@keyframes slideInRight {
  from { transform: translate3d(100%, 0, 0); opacity: 1; }
  to { transform: translate3d(0, 0, 0); opacity: 1; }
}

@keyframes slideInLeft {
  from { transform: translate3d(-100%, 0, 0); opacity: 1; }
  to { transform: translate3d(0, 0, 0); opacity: 1; }
}

@keyframes slideInDown {
  from { transform: translate3d(0, 100%, 0); opacity: 1; }
  to { transform: translate3d(0, 0, 0); opacity: 1; }
}

@keyframes slideInUp {
  from { transform: translate3d(0, -100%, 0); opacity: 1; }
  to { transform: translate3d(0, 0, 0); opacity: 1; }
}

@keyframes zoomFadeIn {
  from { opacity: 0; transform: scale(1.08); }
  to { opacity: 1; transform: scale(1); }
}

/* Keyframes for Typography Content Animations */
@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideUp { from { opacity: 0; transform: translate3d(0, 24px, 0); } to { opacity: 1; transform: translate3d(0, 0, 0); } }
@keyframes slideLeft { from { opacity: 0; transform: translate3d(36px, 0, 0); } to { opacity: 1; transform: translate3d(0, 0, 0); } }
@keyframes slideRight { from { opacity: 0; transform: translate3d(-36px, 0, 0); } to { opacity: 1; transform: translate3d(0, 0, 0); } }
@keyframes zoomIn { from { opacity: 0; transform: scale(.95); } to { opacity: 1; transform: scale(1); } }
@keyframes zoomOut { from { opacity: 0; transform: scale(1.08); } to { opacity: 1; transform: scale(1); } }
@keyframes flipIn { from { opacity: 0; transform: perspective(900px) rotateY(-18deg) translateX(-24px); } to { opacity: 1; transform: perspective(900px) rotateY(0) translateX(0); } }
@keyframes blurIn { from { opacity: 0; filter: blur(10px); transform: translateY(12px); } to { opacity: 1; filter: blur(0); transform: translateY(0); } }
@keyframes bounceIn { 0% { opacity: 0; transform: translateY(34px) scale(.96); } 70% { opacity: 1; transform: translateY(-5px) scale(1.01); } 100% { opacity: 1; transform: translateY(0) scale(1); } }
@keyframes kenBurns { from { transform: scale(1); } to { transform: scale(1.08); } }

@media (max-width: 720px) {
  .carousel-overlay { background: rgba(15, 23, 42, 0.92); }
}
</style>
