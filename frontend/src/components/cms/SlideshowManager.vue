<template>
  <section class="slideshow-builder">
    <div class="cms-panel-head">
      <h2>Slideshow Manager</h2>
      <button type="button" @click="saveConfig">Save settings</button>
    </div>

    <div class="slide-settings">
      <label>Name<input v-model="config.name" /></label>
      <label>Transition
        <select v-model="config.transition_effect">
          <option value="fade">Fade</option>
          <option value="slide">Slide</option>
        </select>
      </label>
      <label>Duration ms<input v-model.number="config.transition_duration" type="number" min="150" max="5000" /></label>
      <label>Autoplay ms<input v-model.number="config.autoplay_speed" type="number" min="2000" max="30000" /></label>
      <label class="check"><input v-model="config.pause_on_hover" type="checkbox" /> Pause on hover/focus</label>
    </div>

    <div class="slide-workspace">
      <aside ref="listRef" class="slide-list" aria-label="Sortable slides">
        <article v-for="(slide, index) in localSlides" :key="slide.id || index" class="slide-item" :class="{ active: selectedIndex === index }" @click="selectedIndex = index">
          <span class="drag-handle" aria-hidden="true">::</span>
          <div><strong>{{ slide.title || 'Untitled slide' }}</strong><small>#{{ index + 1 }} · {{ slide.animation_type }}</small></div>
        </article>
        <button type="button" class="add-slide" @click="addSlide">Add slide</button>
      </aside>

      <form v-if="activeSlide" class="slide-editor" @submit.prevent="saveSlide(activeSlide)">
        <div class="cms-panel-head">
          <h3>{{ activeSlide.id ? 'Edit slide' : 'Create slide' }}</h3>
          <button type="submit">{{ activeSlide.id ? 'Update' : 'Create' }}</button>
        </div>
        <div class="slide-grid">
          <label>Title<input v-model="activeSlide.title" required /></label>
          <label>Subtitle<input v-model="activeSlide.subtitle" /></label>
          <label>Animation
            <select v-model="activeSlide.animation_type">
              <option value="fade-in">Fade In</option>
              <option value="slide-up">Slide Up</option>
              <option value="zoom-in">Zoom In</option>
            </select>
          </label>
          <label>Media type
            <select v-model="activeSlide.media_type">
              <option value="image">Image</option>
              <option value="video">Video</option>
            </select>
          </label>
          <label class="wide">Description<textarea v-model="activeSlide.description" maxlength="500"></textarea></label>
        </div>

        <DropzoneUpload
          v-model="activeSlide.image_url"
          v-model:mediaType="activeSlide.media_type"
          label="slide media"
          hint="Images 16:9, min 1280x720, max 8MB. Video max 24MB."
          accept="image/png,image/jpeg,image/webp"
          allow-video
          video-accept="video/mp4,video/webm"
          :max-image-mb="8"
          :max-video-mb="24"
          :min-width="1280"
          :min-height="720"
          :aspect-ratio="16 / 9"
          @error="message => emit('error', message)"
        />

        <div class="button-editor">
          <div class="cms-panel-head"><h3>CTA buttons</h3><button type="button" :disabled="activeSlide.buttons.length >= 2" @click="addButton">Add button</button></div>
          <article v-for="(button, index) in activeSlide.buttons" :key="index" class="button-row">
            <label>Text<input v-model="button.text" maxlength="80" /></label>
            <label>URL<input v-model="button.url" :class="{ invalid: button.url && !safeUrl(button.url) }" /></label>
            <label>Style<select v-model="button.style_class"><option value="primary">Primary</option><option value="secondary">Secondary</option><option value="outline">Outline</option></select></label>
            <label>Target<select v-model="button.link_target"><option value="_self">Same tab</option><option value="_blank">New tab</option></select></label>
            <button type="button" @click="activeSlide.buttons.splice(index, 1)">Remove</button>
          </article>
        </div>
      </form>
    </div>
  </section>
</template>

<script setup>
import Sortable from 'sortablejs'
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { adminCreateSlide, adminDeleteSlide, adminUpdateSlide, adminUpdateSlideshow } from '@/api/cms.js'
import { debounce } from '@/utils/menuTree.js'
import DropzoneUpload from '@/components/cms/DropzoneUpload.vue'

const props = defineProps({ slideshow: { type: Object, default: () => ({}) }, slides: { type: Array, default: () => [] } })
const emit = defineEmits(['refresh', 'message', 'error'])
const listRef = ref(null)
const selectedIndex = ref(0)
const config = reactive({ name:'Homepage Hero', slug:'homepage-hero', transition_effect:'fade', transition_duration:700, autoplay_speed:6500, pause_on_hover:true, is_active:true })
const localSlides = ref([])
let sortable

const activeSlide = computed(() => localSlides.value[selectedIndex.value])

watch(() => props.slideshow, value => Object.assign(config, value || {}), { immediate:true, deep:true })
watch(() => props.slides, value => {
  localSlides.value = (value || []).map(normalizeSlide)
  if (!localSlides.value.length) addSlide()
}, { immediate:true, deep:true })
watch(localSlides, () => localStorage.setItem('ncs_slideshow_draft', JSON.stringify(localSlides.value)), { deep:true })

onMounted(async () => {
  await nextTick()
  sortable = Sortable.create(listRef.value, {
    animation: 160,
    handle: '.drag-handle',
    onEnd: event => {
      const [moved] = localSlides.value.splice(event.oldIndex, 1)
      localSlides.value.splice(event.newIndex, 0, moved)
      localSlides.value.forEach((slide, index) => { slide.sort_order = index + 1 })
      debouncedSortSave()
    },
  })
})
onBeforeUnmount(() => sortable?.destroy())

function normalizeSlide(slide = {}) {
  return {
    slideshow_id: slide.slideshow_id || config.id || 'homepage-hero',
    title: slide.title || '',
    subtitle: slide.subtitle || '',
    description: slide.description || '',
    image_url: slide.image_url || '',
    media_type: slide.media_type || 'image',
    animation_type: slide.animation_type || 'fade-in',
    buttons: (slide.buttons?.length ? slide.buttons : slide.button_text ? [{ text:slide.button_text, url:slide.button_url || '/', style_class:'primary', link_target:'_self' }] : []).slice(0, 2),
    sort_order: slide.sort_order || localSlides.value.length + 1,
    is_active: slide.is_active !== false,
    id: slide.id || '',
  }
}

function addSlide() {
  localSlides.value.push(normalizeSlide({ title:'New slide', buttons:[{ text:'Learn More', url:'/', style_class:'primary', link_target:'_self' }] }))
  selectedIndex.value = localSlides.value.length - 1
}

function addButton() {
  if (activeSlide.value.buttons.length < 2) activeSlide.value.buttons.push({ text:'', url:'/', style_class:'primary', link_target:'_self' })
}

function safeUrl(url) {
  const value = String(url || '').trim().toLowerCase()
  return /^(\/|https?:\/\/|mailto:)/.test(value) && !value.startsWith('javascript:') && !value.startsWith('data:')
}

async function saveConfig() {
  try {
    await adminUpdateSlideshow(config.slug || 'homepage-hero', config)
    emit('message', 'Slideshow settings saved')
    emit('refresh')
  } catch (error) { emit('error', error) }
}

async function saveSlide(slide) {
  try {
    if (slide.buttons.some(button => !safeUrl(button.url))) throw new Error('CTA URLs must be relative, http(s), or mailto links.')
    const payload = { ...slide, slideshow_id: config.id || 'homepage-hero' }
    slide.id ? await adminUpdateSlide(slide.id, payload) : await adminCreateSlide(payload)
    emit('message', 'Slide saved')
    emit('refresh')
  } catch (error) { emit('error', error) }
}

const debouncedSortSave = debounce(async () => {
  for (const slide of localSlides.value.filter(item => item.id)) {
    await adminUpdateSlide(slide.id, { ...slide, slideshow_id: config.id || 'homepage-hero' })
  }
  emit('message', 'Slide order saved')
  emit('refresh')
}, 700)
</script>

<style scoped>
.slideshow-builder { display: grid; gap: 1rem; min-width: 0; }

.slide-settings {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr));
  gap: 1rem;
}
.slide-settings label,
.slide-grid label,
.button-row label {
  display: grid;
  gap: .35rem;
  font-size: .76rem;
  font-weight: 800;
  color: #475569;
  min-width: 0;
}
.slide-settings input,
.slide-settings select,
.slide-grid input,
.slide-grid select,
.slide-grid textarea,
.button-row input,
.button-row select {
  width: 100%;
  box-sizing: border-box;
  border: 1px solid #cbd5e1;
  border-radius: .45rem;
  padding: .65rem;
  color: #111827;
  background: white;
}
.slide-settings .check { display: flex; align-items: center; gap: .5rem; }

.slide-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(14rem, 1fr));
  gap: 1rem;
}
.slide-grid .wide { grid-column: 1 / -1; }
.slide-grid textarea { min-height: 6rem; }

.slide-workspace {
  display: grid;
  grid-template-columns: minmax(14rem, 18rem) minmax(0, 1fr);
  gap: 1rem;
  min-width: 0;
}

.slide-list { display: grid; align-content: start; gap: .5rem; min-width: 0; }
.slide-item {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: .7rem;
  align-items: center;
  border: 1px solid #e5e7eb;
  border-radius: .45rem;
  padding: .75rem;
  background: #f8fafc;
  cursor: pointer;
  min-width: 0;
}
.slide-item.active { border-color: #f5a623; background: #fff7ed; }
.drag-handle { cursor: grab; color: #94a3b8; font-weight: 900; }
.slide-item strong,
.slide-item small { display: block; overflow-wrap: anywhere; }
.slide-item small { color: #64748b; }

.add-slide,
.button-editor button {
  border: 1px solid #cbd5e1;
  border-radius: .45rem;
  padding: .65rem .9rem;
  font-weight: 800;
}

.slide-editor { display: grid; gap: 1rem; min-width: 0; }

.button-editor { display: grid; gap: .75rem; min-width: 0; }
.button-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(8rem, 1fr));
  gap: .75rem;
  align-items: end;
  min-width: 0;
}
.button-row .invalid { border-color: #dc2626; background: #fef2f2; }

@media (max-width: 56.25rem) {
  .slide-workspace { grid-template-columns: 1fr; }
}

:global(.dark) .slide-settings input,
:global(.dark) .slide-settings select,
:global(.dark) .slide-grid input,
:global(.dark) .slide-grid select,
:global(.dark) .slide-grid textarea,
:global(.dark) .button-row input,
:global(.dark) .button-row select { background: #0f172a; color: #f8fafc; border-color: #475569; }
:global(.dark) .slide-item { background: #111827; border-color: #334155; }
:global(.dark) .slide-item.active { border-color: #f5a623; background: #422006; }
</style>
