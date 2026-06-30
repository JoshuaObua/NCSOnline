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
    @mouseenter="pause"
    @mouseleave="resume"
    @focusin="pause"
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
    <button v-if="normalizedSlides.length > 1" type="button" class="carousel-arrow prev" aria-label="Previous slide" @click="previous">‹</button>
    <button v-if="normalizedSlides.length > 1" type="button" class="carousel-arrow next" aria-label="Next slide" @click="next">›</button>
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

const normalizedSlides = computed(() => {
  const source = props.slides?.length ? props.slides : [{ id:'default', title:'Welcome to National Council of Sports', subtitle:'National Council of Sports', description:'A Centre of Excellence for Promotion and Development of Sports', button_text:'Discover NCS', button_url:'/pages/the-mandate', image_url:'' }]
  return source.map(slide => ({
    ...slide,
    media_type: slide.media_type || 'image',
    animation_type: slide.animation_type || 'fade-in',
    buttons: (slide.buttons?.length ? slide.buttons : [{ text: slide.button_text || 'Learn More', url: slide.button_url || '/pages/the-mandate', style_class:'primary', link_target:'_self', sort_order:0 }]).slice(0, 2),
  }))
})
const interval = computed(() => Number(props.slideshow?.autoplay_speed || 6500))
const duration = computed(() => Number(props.slideshow?.transition_duration || 700))
const effect = computed(() => props.slideshow?.transition_effect === 'slide' ? 'slide' : 'fade')

watch(normalizedSlides, () => { if (current.value >= normalizedSlides.value.length) current.value = 0; restart() })
onMounted(restart)
onUnmounted(() => clearInterval(timer))

function shouldLoad(index) {
  return index === current.value || Math.abs(index - current.value) <= 1
}
function transitionClass(index) {
  if (index === current.value) return 'active'
  return effect.value === 'slide' ? (index < current.value ? 'before' : 'after') : ''
}
function next() { go((current.value + 1) % normalizedSlides.value.length) }
function previous() { go((current.value - 1 + normalizedSlides.value.length) % normalizedSlides.value.length) }
function go(index) { current.value = index; restart() }
function pause() { paused.value = true; clearInterval(timer) }
function resume() { paused.value = false; restart() }
function restart() {
  clearInterval(timer)
  if (!paused.value && normalizedSlides.value.length > 1) timer = setInterval(next, interval.value)
}
function isInternal(url = '') {
  return String(url).startsWith('/')
}
</script>

<style scoped>
.public-carousel{position:relative;height:clamp(31rem,65vw,43rem);overflow:hidden;background:#1a365d;outline:none}.carousel-slide{position:absolute;inset:0;opacity:0;transform:translate3d(0,0,0) scale(1.02);transition:opacity var(--slide-duration,700ms) ease,transform var(--slide-duration,700ms) ease;will-change:opacity,transform}.carousel-slide.active{z-index:2;opacity:1;transform:translate3d(0,0,0) scale(1)}.carousel-slide.before{transform:translate3d(-100%,0,0);opacity:1}.carousel-slide.after{transform:translate3d(100%,0,0);opacity:1}.carousel-media,.carousel-fallback{width:100%;height:100%;object-fit:cover}.carousel-fallback{background:linear-gradient(120deg,#1a365d,#274d7e)}.carousel-overlay{position:absolute;inset:0;background:linear-gradient(90deg,rgb(15 31 61/.94),rgb(26 54 93/.68),transparent)}.carousel-content{position:absolute;z-index:3;inset:0;display:flex;align-items:center;width:min(80rem,100%);margin:auto;padding:0 1.25rem;color:white}.carousel-content>div{max-width:43rem}.carousel-content span{display:inline-block;margin-bottom:1rem;border-radius:999px;background:#f5a623;padding:.4rem .9rem;font-size:.85rem;font-weight:800}.carousel-content h1{font-size:clamp(2.5rem,6vw,4.6rem);font-weight:850;line-height:1.05}.carousel-content p{max-width:39rem;margin:1.25rem 0 2rem;font-size:clamp(1.05rem,2vw,1.3rem);line-height:1.65}.carousel-slide.active.fade-in .carousel-content>div{animation:fadeIn .72s ease both}.carousel-slide.active.slide-up .carousel-content>div{animation:slideUp .72s ease both}.carousel-slide.active.zoom-in .carousel-content>div{animation:zoomIn .72s ease both}.carousel-actions{display:flex;flex-wrap:wrap;gap:.8rem}.carousel-btn{display:inline-flex;min-height:44px;align-items:center;justify-content:center;border-radius:.55rem;padding:.75rem 1.2rem;font-weight:850;transition:transform .2s,background .2s,color .2s}.carousel-btn:hover{transform:translateY(-2px)}.carousel-btn-primary{background:#f5a623;color:#172b4d}.carousel-btn-secondary{background:#1a365d;color:white}.carousel-btn-outline{border:1px solid rgb(255 255 255/.78);color:white}.carousel-arrow{position:absolute;z-index:4;top:50%;width:3rem;height:3rem;border-radius:999px;background:rgb(255 255 255/.16);color:white;font-size:2rem;line-height:1}.carousel-arrow:hover,.carousel-arrow:focus-visible{background:#f5a623;color:#172b4d}.prev{left:1rem}.next{right:1rem}.carousel-dots{position:absolute;z-index:4;bottom:1.5rem;left:50%;display:flex;gap:.5rem;transform:translateX(-50%)}.carousel-dots button{width:.7rem;height:.7rem;border-radius:999px;background:rgb(255 255 255/.5);transition:width .2s,background .2s}.carousel-dots button.active{width:2.2rem;background:#f5a623}.carousel-bottom{position:absolute;z-index:3;right:0;bottom:0;left:0;height:6rem;background:linear-gradient(to top,white,transparent);pointer-events:none}@keyframes fadeIn{from{opacity:0}to{opacity:1}}@keyframes slideUp{from{opacity:0;transform:translate3d(0,24px,0)}to{opacity:1;transform:translate3d(0,0,0)}}@keyframes zoomIn{from{opacity:0;transform:scale(.95)}to{opacity:1;transform:scale(1)}}@media(max-width:720px){.carousel-overlay{background:rgb(15 31 61/.82)}.carousel-arrow{display:none}.carousel-content h1{font-size:2.55rem}.carousel-btn{width:100%}}
</style>
