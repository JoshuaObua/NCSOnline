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
      <video v-if="slide.media_type === 'video' && shouldLoad(index)" class="carousel-media" :src="mediaUrl(slide.image_url)" autoplay muted loop playsinline preload="metadata"></video>
      <img v-else-if="slide.image_url && shouldLoad(index)" class="carousel-media" :src="mediaUrl(slide.image_url)" :alt="slide.title" :loading="index === 0 ? 'eager' : 'lazy'" />
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
            >{{ button.text }}</component>
          </div>
        </div>
      </div>
    </article>
    <button v-if="normalizedSlides.length > 1" type="button" class="carousel-arrow prev" aria-label="Previous slide" @click="previous">&lsaquo;</button>
    <button v-if="normalizedSlides.length > 1" type="button" class="carousel-arrow next" aria-label="Next slide" @click="next">&rsaquo;</button>
    <div v-if="normalizedSlides.length > 1" class="carousel-dots">
      <button v-for="(_, index) in normalizedSlides" :key="index" type="button" :class="{ active: index === current }" :aria-label="`Show slide ${index + 1}`" @click="go(index)"></button>
    </div>
    <div class="carousel-bottom"></div>
  </section>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { mediaUrl } from '@/api/client.js'

const props = defineProps({ slideshow: { type: Object, default: () => ({}) }, slides: { type: Array, default: () => [] } })
const current = ref(0)
const paused = ref(false)
let timer

const fallbackSlide = {
  id:'default',
  title:'Welcome to National Council of Sports',
  subtitle:'National Council of Sports',
  description:'A Centre of Excellence for Promotion and Development of Sports',
  button_text:'Discover NCS',
  button_url:'/pages/the-mandate',
  animation_type:'fade-in',
  image_url:'',
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
const duration = computed(() => Number(props.slideshow?.transition_duration || 700))
const effect = computed(() => ['fade', 'slide', 'vertical', 'zoom'].includes(props.slideshow?.transition_effect) ? props.slideshow.transition_effect : 'fade')
const pauseOnHover = computed(() => props.slideshow?.pause_on_hover !== false)

watch(normalizedSlides, () => {
  if (current.value >= normalizedSlides.value.length) current.value = 0
  restart()
})
onMounted(restart)
onUnmounted(() => clearInterval(timer))

function normalizeButtons(slide) {
  const buttons = Array.isArray(slide.buttons) && slide.buttons.length
    ? slide.buttons
    : [{ text: slide.button_text || 'Learn More', url: slide.button_url || '/pages/the-mandate', style_class:'primary', link_target:'_self', sort_order:0 }]
  return buttons.slice(0, 2).map((button, index) => ({
    text: button.text || (index ? 'Read More' : 'Learn More'),
    url: button.url || '/',
    style_class: ['primary', 'secondary', 'outline'].includes(button.style_class) ? button.style_class : (index ? 'outline' : 'primary'),
    link_target: button.link_target || '_self',
    sort_order: button.sort_order ?? index,
  }))
}
function shouldLoad(index) {
  return index === current.value || Math.abs(index - current.value) <= 1
}
function transitionClass(index) {
  if (index === current.value) return 'active'
  if (effect.value === 'slide') return index < current.value ? 'before' : 'after'
  if (effect.value === 'vertical') return index < current.value ? 'above' : 'below'
  if (effect.value === 'zoom') return 'zoom-waiting'
  return ''
}
function next() { go((current.value + 1) % normalizedSlides.value.length) }
function previous() { go((current.value - 1 + normalizedSlides.value.length) % normalizedSlides.value.length) }
function go(index) { current.value = index; restart() }
function pauseIfEnabled() { if (pauseOnHover.value) pause() }
function pause() { paused.value = true; clearInterval(timer) }
function resume() { paused.value = false; restart() }
function restart() {
  clearInterval(timer)
  if (!paused.value && normalizedSlides.value.length > 1) timer = setInterval(next, interval.value)
}
function isInternal(url = '') {
  return String(url).startsWith('/')
}
function animationMode(value) {
  return ['fade-in', 'slide-up', 'slide-left', 'slide-right', 'zoom-in', 'zoom-out', 'flip-in', 'blur-in', 'bounce-in', 'ken-burns'].includes(value) ? value : 'fade-in'
}
</script>

<style scoped>
.public-carousel{position:relative;height:clamp(31rem,65vw,43rem);overflow:hidden;background:#1a365d;outline:none}.carousel-slide{position:absolute;inset:0;opacity:0;transform:translate3d(0,0,0) scale(1.02);transition:opacity var(--slide-duration,700ms) ease,transform var(--slide-duration,700ms) ease,filter var(--slide-duration,700ms) ease;will-change:opacity,transform,filter}.carousel-slide.active{z-index:2;opacity:1;transform:translate3d(0,0,0) scale(1);filter:none}.carousel-slide.before{transform:translate3d(-100%,0,0);opacity:1}.carousel-slide.after{transform:translate3d(100%,0,0);opacity:1}.carousel-slide.above{transform:translate3d(0,-100%,0);opacity:1}.carousel-slide.below{transform:translate3d(0,100%,0);opacity:1}.carousel-slide.zoom-waiting{transform:scale(.86);opacity:0}.carousel-media,.carousel-fallback{width:100%;height:100%;object-fit:cover}.carousel-fallback{background:linear-gradient(120deg,#1a365d,#274d7e)}.carousel-overlay{position:absolute;inset:0;background:linear-gradient(90deg,rgb(15 31 61/.94),rgb(26 54 93/.68),transparent)}.carousel-content{position:absolute;z-index:3;inset:0;display:flex;align-items:center;width:min(80rem,100%);margin:auto;padding:0 1.25rem;color:white}.carousel-content>div{max-width:43rem}.carousel-content span{display:inline-block;margin-bottom:1rem;border-radius:999px;background:#f5a623;padding:.4rem .9rem;font-size:.85rem;font-weight:800}.carousel-content h1{font-size:clamp(2.5rem,6vw,4.6rem);font-weight:850;line-height:1.05}.carousel-content p{max-width:39rem;margin:1.25rem 0 2rem;font-size:clamp(1.05rem,2vw,1.3rem);line-height:1.65}.carousel-slide.active.fade-in .carousel-content>div{animation:fadeIn .72s ease both}.carousel-slide.active.slide-up .carousel-content>div{animation:slideUp .72s ease both}.carousel-slide.active.slide-left .carousel-content>div{animation:slideLeft .72s ease both}.carousel-slide.active.slide-right .carousel-content>div{animation:slideRight .72s ease both}.carousel-slide.active.zoom-in .carousel-content>div{animation:zoomIn .72s ease both}.carousel-slide.active.zoom-out .carousel-content>div{animation:zoomOut .72s ease both}.carousel-slide.active.flip-in .carousel-content>div{animation:flipIn .78s ease both;transform-origin:left center}.carousel-slide.active.blur-in .carousel-content>div{animation:blurIn .78s ease both}.carousel-slide.active.bounce-in .carousel-content>div{animation:bounceIn .82s cubic-bezier(.2,.85,.35,1.2) both}.carousel-slide.active.ken-burns .carousel-media,.carousel-slide.active.ken-burns .carousel-fallback{animation:kenBurns 7s ease both}.carousel-actions{display:flex;flex-wrap:wrap;gap:.8rem}.carousel-btn{display:inline-flex;min-height:44px;align-items:center;justify-content:center;border-radius:.55rem;padding:.75rem 1.2rem;font-weight:850;transition:transform .2s,background .2s,color .2s}.carousel-btn:hover{transform:translateY(-2px)}.carousel-btn-primary{background:#f5a623;color:#172b4d}.carousel-btn-secondary{background:#1a365d;color:white}.carousel-btn-outline{border:1px solid rgb(255 255 255/.78);color:white}.carousel-arrow{position:absolute;z-index:4;top:50%;width:3rem;height:3rem;border-radius:999px;background:rgb(255 255 255/.16);color:white;font-size:2rem;line-height:1}.carousel-arrow:hover,.carousel-arrow:focus-visible{background:#f5a623;color:#172b4d}.prev{left:1rem}.next{right:1rem}.carousel-dots{position:absolute;z-index:4;bottom:1.5rem;left:50%;display:flex;gap:.5rem;transform:translateX(-50%)}.carousel-dots button{width:.7rem;height:.7rem;border-radius:999px;background:rgb(255 255 255/.5);transition:width .2s,background .2s}.carousel-dots button.active{width:2.2rem;background:#f5a623}.carousel-bottom{position:absolute;z-index:3;right:0;bottom:0;left:0;height:6rem;background:linear-gradient(to top,white,transparent);pointer-events:none}@keyframes fadeIn{from{opacity:0}to{opacity:1}}@keyframes slideUp{from{opacity:0;transform:translate3d(0,24px,0)}to{opacity:1;transform:translate3d(0,0,0)}}@keyframes slideLeft{from{opacity:0;transform:translate3d(36px,0,0)}to{opacity:1;transform:translate3d(0,0,0)}}@keyframes slideRight{from{opacity:0;transform:translate3d(-36px,0,0)}to{opacity:1;transform:translate3d(0,0,0)}}@keyframes zoomIn{from{opacity:0;transform:scale(.95)}to{opacity:1;transform:scale(1)}}@keyframes zoomOut{from{opacity:0;transform:scale(1.08)}to{opacity:1;transform:scale(1)}}@keyframes flipIn{from{opacity:0;transform:perspective(900px) rotateY(-18deg) translateX(-24px)}to{opacity:1;transform:perspective(900px) rotateY(0) translateX(0)}}@keyframes blurIn{from{opacity:0;filter:blur(10px);transform:translateY(12px)}to{opacity:1;filter:blur(0);transform:translateY(0)}}@keyframes bounceIn{0%{opacity:0;transform:translateY(34px) scale(.96)}70%{opacity:1;transform:translateY(-5px) scale(1.01)}100%{opacity:1;transform:translateY(0) scale(1)}}@keyframes kenBurns{from{transform:scale(1)}to{transform:scale(1.08)}}@media(max-width:720px){.carousel-overlay{background:rgb(15 31 61/.82)}.carousel-arrow{display:none}.carousel-content h1{font-size:2.55rem}.carousel-btn{width:100%}}
</style>
