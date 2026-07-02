<template>
  <section class="slideshow-builder">
    <div class="card slideshow-card">
      <div class="card-header">
        <h4>Slideshow Manager</h4>
        <div class="card-header-action">
          <button type="button" class="btn btn-icon icon-left btn-primary" @click="addSlide"><i class="fas fa-plus"></i> Add slide</button>
          <button type="button" class="btn btn-icon icon-left btn-info" @click="settingsOpen = true"><i class="fas fa-cog"></i> Slideshow settings</button>
        </div>
      </div>
    </div>

    <div class="slide-workspace">
      <aside class="card slide-list-card" aria-label="Sortable slides">
        <div class="card-header">
          <h4>Slides</h4>
        </div>
        <div class="card-body p-0">
          <div class="table-responsive">
            <table class="table table-striped table-hover table-sm mb-0">
              <thead>
                <tr>
                  <th class="text-center">#</th>
                  <th>Title</th>
                  <th>Status</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody ref="listRef">
                <tr v-for="(slide, index) in localSlides" :key="slide.id || index" :class="{ 'table-active': selectedIndex === index }">
                  <th scope="row" class="p-0 text-center drag-handle"><i class="fas fa-grip-vertical"></i></th>
                  <td>
                    <button type="button" class="slide-title-link" @click="selectedIndex = index">
                      <strong>{{ slide.title || 'Untitled slide' }}</strong>
                      <small>#{{ index + 1 }} - {{ slide.animation_type }}</small>
                    </button>
                  </td>
                  <td><span class="badge" :class="slide.is_active ? 'badge-success' : 'badge-secondary'">{{ slide.is_active ? 'Active' : 'Hidden' }}</span></td>
                  <td>
                    <button type="button" class="btn btn-sm btn-outline-primary mr-1" @click="selectedIndex = index">Edit</button>
                    <button type="button" class="btn btn-sm btn-outline-danger" @click="deleteSlide(slide, index)">Delete</button>
                  </td>
                </tr>
                <tr v-if="!localSlides.length">
                  <td colspan="4" class="text-center text-muted">No slides have been added yet.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </aside>

      <form v-if="activeSlide" class="card slide-editor" @submit.prevent="saveSlide(activeSlide)">
        <div class="card-header">
          <h4>{{ activeSlide.id ? 'Edit slide' : 'Create slide' }}</h4>
          <div class="card-header-action">
            <button type="submit" class="btn btn-icon icon-left btn-primary"><i class="fas fa-save"></i> {{ activeSlide.id ? 'Update' : 'Create' }}</button>
          </div>
        </div>
        <div class="card-body">
          <div class="row">
            <div class="form-group col-md-6">
              <label>Title</label>
              <input v-model="activeSlide.title" class="form-control" required />
            </div>
            <div class="form-group col-md-6">
              <label>Subtitle</label>
              <input v-model="activeSlide.subtitle" class="form-control" />
            </div>
            <div class="form-group col-md-4">
              <label>Animation</label>
              <select v-model="activeSlide.animation_type" class="form-control selectric">
                <option v-for="mode in animationModes" :key="mode.value" :value="mode.value">{{ mode.label }}</option>
              </select>
            </div>
            <div class="form-group col-md-4">
              <label>Media type</label>
              <select v-model="activeSlide.media_type" class="form-control selectric">
                <option value="image">Image</option>
                <option value="video">Video</option>
              </select>
            </div>
            <div class="form-group col-md-4">
              <label>Status</label>
              <select v-model="activeSlide.is_active" class="form-control selectric">
                <option :value="true">Active</option>
                <option :value="false">Hidden</option>
              </select>
            </div>
            <div class="form-group col-12">
              <label>Description</label>
              <textarea v-model="activeSlide.description" class="form-control" maxlength="500"></textarea>
            </div>
          </div>

          <DropzoneUpload
            v-model="activeSlide.image_url"
            v-model:mediaType="activeSlide.media_type"
            label="slide media"
            hint="Drop files here or click to upload. PNG, JPG, WEBP, MP4, or WEBM."
            accept="image/png,image/jpeg,image/webp"
            allow-video
            video-accept="video/mp4,video/webm"
            :max-image-mb="8"
            :max-video-mb="24"
            @error="message => emit('error', message)"
          />

          <div class="button-editor">
            <div class="card-header compact-header">
              <h4>CTA buttons</h4>
              <div class="card-header-action">
                <button type="button" class="btn btn-sm btn-primary" :disabled="activeSlide.buttons.length >= 2" @click="addButton">Add button</button>
              </div>
            </div>
            <article v-for="(button, index) in activeSlide.buttons" :key="index" class="button-row">
              <div class="button-row-head">
                <strong>Button {{ index + 1 }}</strong>
                <button type="button" class="btn btn-outline-danger btn-sm" @click="removeButton(index)"><i class="fas fa-trash"></i></button>
              </div>
              <div class="button-fields">
                <div class="form-group">
                  <label>Text</label>
                  <input v-model="button.text" class="form-control" maxlength="80" />
                </div>
                <div class="form-group button-url-field">
                  <label>URL</label>
                  <input v-model="button.url" class="form-control" :class="{ 'is-invalid': button.url && !safeUrl(button.url) }" />
                </div>
                <div class="form-group">
                  <label>Style</label>
                  <select v-model="button.style_class" class="form-control selectric">
                    <option value="primary">Primary</option>
                    <option value="secondary">Secondary</option>
                    <option value="outline">Outline</option>
                  </select>
                </div>
                <div class="form-group">
                  <label>Target</label>
                  <select v-model="button.link_target" class="form-control selectric">
                    <option value="_self">Same tab</option>
                    <option value="_blank">New tab</option>
                  </select>
                </div>
              </div>
            </article>
          </div>
        </div>
        <div class="card-footer text-right">
          <button type="button" class="btn btn-secondary mr-1" @click="settingsOpen = true">Settings</button>
          <button type="submit" class="btn btn-primary">{{ activeSlide.id ? 'Update slide' : 'Create slide' }}</button>
        </div>
      </form>
    </div>

    <div v-if="settingsOpen" class="modal fade show slideshow-modal" tabindex="-1" role="dialog" aria-modal="true" style="display:block">
      <div class="modal-dialog modal-lg" role="document">
        <form class="modal-content" @submit.prevent="saveConfig">
          <div class="modal-header">
            <h5 class="modal-title">Slideshow Settings</h5>
            <button type="button" class="close" aria-label="Close" @click="settingsOpen = false"><span aria-hidden="true">&times;</span></button>
          </div>
          <div class="modal-body">
            <div class="row">
              <div class="form-group col-md-6">
                <label>Name</label>
                <input v-model="config.name" class="form-control" />
              </div>
              <div class="form-group col-md-6">
                <label>Transition</label>
                <select v-model="config.transition_effect" class="form-control selectric">
                  <option value="fade">Fade</option>
                  <option value="slide">Slide</option>
                  <option value="vertical">Vertical Slide</option>
                  <option value="zoom">Zoom</option>
                </select>
              </div>
              <div class="form-group col-md-6">
                <label>Duration ms</label>
                <input v-model.number="config.transition_duration" class="form-control" type="number" min="150" max="5000" />
              </div>
              <div class="form-group col-md-6">
                <label>Autoplay ms</label>
                <input v-model.number="config.autoplay_speed" class="form-control" type="number" min="2000" max="30000" />
              </div>
              <div class="form-group col-12">
                <div class="custom-control custom-checkbox">
                  <input id="pause-on-hover" v-model="config.pause_on_hover" type="checkbox" class="custom-control-input" />
                  <label class="custom-control-label" for="pause-on-hover">Pause on hover/focus</label>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer bg-whitesmoke br">
            <button type="button" class="btn btn-secondary" @click="settingsOpen = false">Close</button>
            <button type="submit" class="btn btn-primary">Save settings</button>
          </div>
        </form>
      </div>
    </div>
    <div v-if="settingsOpen" class="modal-backdrop fade show"></div>
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
const settingsOpen = ref(false)
const config = reactive({ name:'Homepage Hero', slug:'homepage-hero', transition_effect:'fade', transition_duration:700, autoplay_speed:6500, pause_on_hover:true, is_active:true })
const localSlides = ref([])
let sortable
let sweetAlertPromise

const animationModes = [
  { value:'fade-in', label:'Fade In' },
  { value:'slide-up', label:'Slide Up' },
  { value:'slide-left', label:'Slide Left' },
  { value:'slide-right', label:'Slide Right' },
  { value:'zoom-in', label:'Zoom In' },
  { value:'zoom-out', label:'Zoom Out' },
  { value:'flip-in', label:'Flip In' },
  { value:'blur-in', label:'Blur In' },
  { value:'bounce-in', label:'Bounce In' },
  { value:'ken-burns', label:'Ken Burns' },
]

const activeSlide = computed(() => localSlides.value[selectedIndex.value])

watch(() => props.slideshow, value => Object.assign(config, value || {}), { immediate:true, deep:true })
watch(() => props.slides, value => {
  localSlides.value = (value || []).map(normalizeSlide)
  if (!localSlides.value.length) addSlide()
  selectedIndex.value = Math.min(selectedIndex.value, Math.max(0, localSlides.value.length - 1))
}, { immediate:true, deep:true })
watch(localSlides, () => localStorage.setItem('ncs_slideshow_draft', JSON.stringify(localSlides.value)), { deep:true })

onMounted(async () => {
  loadSweetAlert()
  await nextTick()
  if (!listRef.value) return
  sortable = Sortable.create(listRef.value, {
    animation: 160,
    handle: '.drag-handle',
    onEnd: event => {
      if (event.oldIndex == null || event.newIndex == null || event.oldIndex === event.newIndex) return
      const [moved] = localSlides.value.splice(event.oldIndex, 1)
      if (!moved) return
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
  if (activeSlide.value && activeSlide.value.buttons.length < 2) activeSlide.value.buttons.push({ text:'', url:'/', style_class:'primary', link_target:'_self' })
}

function removeButton(index) {
  activeSlide.value?.buttons.splice(index, 1)
}

function safeUrl(url) {
  const value = String(url || '').trim().toLowerCase()
  return /^(\/|https?:\/\/|mailto:)/.test(value) && !value.startsWith('javascript:') && !value.startsWith('data:')
}

async function saveConfig() {
  try {
    await adminUpdateSlideshow(config.slug || 'homepage-hero', config)
    settingsOpen.value = false
    emit('message', 'Slideshow settings saved')
    notify('Saved', 'Slideshow settings were updated.', 'success')
    emit('refresh')
  } catch (error) { notify('Save failed', error.message || 'Could not save slideshow settings.', 'error'); emit('error', error) }
}

async function saveSlide(slide) {
  try {
    if (slide.buttons.some(button => button.url && !safeUrl(button.url))) throw new Error('CTA URLs must be relative, http(s), or mailto links.')
    const payload = slidePayload(slide)
    const response = slide.id ? await adminUpdateSlide(slide.id, payload) : await adminCreateSlide(payload)
    const saved = response?.data?.data || response?.data || {}
    if (!slide.id && saved.id) slide.id = saved.id
    emit('message', 'Slide saved')
    notify('Saved', 'Slide content was saved and will feed the public homepage slideshow.', 'success')
    emit('refresh')
  } catch (error) { notify('Save failed', error.message || 'Could not save this slide.', 'error'); emit('error', error) }
}

async function deleteSlide(slide, index) {
  const confirmed = await confirmAction('Delete this slide?', 'This removes the slide from the public homepage slideshow.')
  if (!confirmed) return
  try {
    if (slide.id) await adminDeleteSlide(slide.id)
    localSlides.value.splice(index, 1)
    if (selectedIndex.value >= localSlides.value.length) selectedIndex.value = Math.max(0, localSlides.value.length - 1)
    if (!localSlides.value.length) addSlide()
    emit('message', 'Slide deleted')
    notify('Deleted', 'Slide deleted successfully.', 'success')
    emit('refresh')
  } catch (error) { notify('Delete failed', error.message || 'Could not delete this slide.', 'error'); emit('error', error) }
}

const debouncedSortSave = debounce(async () => {
  for (const slide of localSlides.value.filter(item => item.id)) {
    await adminUpdateSlide(slide.id, slidePayload(slide))
  }
  emit('message', 'Slide order saved')
  emit('refresh')
}, 700)

function slidePayload(slide) {
  const buttons = (slide.buttons || []).filter(button => button.text || button.url).slice(0, 2).map((button, index) => ({
    text: button.text || (index ? 'Read More' : 'Learn More'),
    url: button.url || '/',
    style_class: button.style_class || (index ? 'outline' : 'primary'),
    link_target: button.link_target || '_self',
    sort_order: index + 1,
  }))
  const primary = buttons[0] || { text:'Learn More', url:'/', style_class:'primary', link_target:'_self' }
  return {
    ...slide,
    slideshow_id: config.id || 'homepage-hero',
    buttons,
    button_text: primary.text,
    button_url: primary.url,
    button_style: primary.style_class,
    button_target: primary.link_target,
    animation_type: slide.animation_type || 'fade-in',
    is_active: slide.is_active !== false,
  }
}

function loadSweetAlert() {
  if (window.swal) return Promise.resolve(window.swal)
  if (sweetAlertPromise) return sweetAlertPromise
  sweetAlertPromise = new Promise(resolve => {
    const existing = document.querySelector('script[data-otika-sweetalert]')
    if (existing) {
      existing.addEventListener('load', () => resolve(window.swal), { once:true })
      existing.addEventListener('error', () => resolve(null), { once:true })
      return
    }
    const script = document.createElement('script')
    script.src = '/otika-assets/bundles/sweetalert/sweetalert.min.js'
    script.async = true
    script.dataset.otikaSweetalert = 'true'
    script.onload = () => resolve(window.swal)
    script.onerror = () => resolve(null)
    document.head.appendChild(script)
  })
  return sweetAlertPromise
}

async function confirmAction(title, text) {
  const swal = await loadSweetAlert()
  if (!swal) {
    emit('error', new Error('SweetAlert could not be loaded. Please refresh and try again.'))
    return false
  }
  const result = await swal({ title, text, icon:'warning', buttons:['Cancel', 'Yes, delete it'], dangerMode:true })
  return Boolean(result)
}

async function notify(title, text, icon = 'success') {
  const swal = await loadSweetAlert()
  if (!swal) return
  swal({ title, text, icon, timer: icon === 'success' ? 1800 : undefined })
}
</script>

<style scoped>
.slideshow-builder{display:grid;gap:20px;min-width:0}.slideshow-card,.slide-list-card,.slide-editor{border:0;border-radius:3px;box-shadow:0 4px 25px rgba(0,0,0,.1)}.slideshow-builder .card-header h4{font-size:16px;font-weight:700;color:#34395e}.slide-workspace{display:grid;grid-template-columns:minmax(24rem,30rem) minmax(0,1fr);gap:20px;min-width:0}.slide-title-link{border:0;background:transparent;padding:0;text-align:left;color:#34395e}.slide-title-link strong,.slide-title-link small{display:block}.slide-title-link small{color:#6c757d;font-size:12px}.drag-handle{cursor:grab;color:#98a6ad}.table-active td,.table-active th{background:#f4f6ff!important}.slide-editor .card-body{padding:25px}.slide-editor .form-group{margin-bottom:16px}.slide-editor label,.slideshow-modal label{font-size:12px;font-weight:600;color:#34395e}.slide-editor .form-control,.slideshow-modal .form-control{height:auto;border:1px solid #e4e6fc;border-radius:3px;background:#fdfdff;color:#495057;padding:10px 15px}.slide-editor textarea.form-control{min-height:92px}.slide-editor select option,.slideshow-modal select option{background:#111827;color:#fff}.compact-header{padding:18px 0 12px!important;border-bottom:1px solid #f9f9f9!important}.button-editor{display:grid;gap:14px;margin-top:22px}.button-row{display:grid;gap:14px;margin:0;padding:16px;border:1px solid #f0f1ff;border-radius:3px;background:#fbfbff}.button-row-head{display:flex;align-items:center;justify-content:space-between;gap:12px}.button-row-head strong{color:#34395e;font-size:13px}.button-fields{display:grid;grid-template-columns:minmax(9rem,1fr) minmax(14rem,2fr) minmax(8rem,.8fr) minmax(8rem,.8fr);gap:14px;align-items:start}.button-url-field{min-width:0}.slideshow-modal{z-index:1060}.modal-backdrop{z-index:1050}@media(max-width:1200px){.button-fields{grid-template-columns:repeat(2,minmax(0,1fr))}}@media(max-width:1100px){.slide-workspace{grid-template-columns:1fr}}@media(max-width:640px){.button-fields{grid-template-columns:1fr}.slideshow-card .card-header,.slide-editor .card-header{align-items:flex-start;flex-direction:column;gap:12px}}:global(.dark) .slideshow-card,:global(.dark) .slide-list-card,:global(.dark) .slide-editor,:global(.dark) .slideshow-modal .modal-content{background:#1f2937!important;color:#e5e7eb!important}:global(.dark) .slideshow-builder .card-header h4,:global(.dark) .slide-title-link,:global(.dark) .slide-editor label,:global(.dark) .slideshow-modal label,:global(.dark) .slideshow-modal .modal-title,:global(.dark) .button-row-head strong{color:#f8fafc!important}:global(.dark) .slide-title-link small{color:#cbd5e1!important}:global(.dark) .slide-editor .form-control,:global(.dark) .slideshow-modal .form-control{background:#111827!important;color:#f8fafc!important;border-color:#475569!important}:global(.dark) .slide-editor select option,:global(.dark) .slideshow-modal select option{background:#fff!important;color:#111827!important}:global(.dark) .button-row{background:#111827!important;border-color:#334155!important}:global(.dark) .table-active td,:global(.dark) .table-active th{background:#111827!important}
</style>
