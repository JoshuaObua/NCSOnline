<template>
  <LayoutDefault title="Application Forms">
    <div class="p-6 max-w-7xl mx-auto">
      <!-- Toast -->
      <Transition name="toast">
        <div v-if="toast" class="fixed top-5 right-5 z-[100] px-4 py-3 bg-gray-900 text-white text-sm rounded-xl shadow-xl">
          {{ toast }}
        </div>
      </Transition>

      <!-- ═════════ LIST MODE ═════════ -->
      <div v-if="mode === 'list'">
        <div class="flex items-center justify-between mb-5">
          <div>
            <h2 class="text-2xl font-semibold text-gray-900">Application Forms</h2>
            <p class="text-xs text-gray-500 mt-1">Build the forms that applicants fill out on the public portal. Each form belongs to a department.</p>
          </div>
          <button @click="openBuilder()" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-4 py-2 rounded-lg flex items-center gap-1.5">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15"/></svg>
            New Form
          </button>
        </div>

        <div class="flex gap-3 mb-4">
          <select v-model="statusFilter" @change="loadForms" class="text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white">
            <option value="">All statuses</option>
            <option value="DRAFT">Draft</option>
            <option value="OPEN">Open</option>
            <option value="CLOSED">Closed</option>
            <option value="ARCHIVED">Archived</option>
          </select>
        </div>

        <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="i in 3" :key="i" class="h-40 bg-gray-100 rounded-xl animate-pulse"/>
        </div>

        <div v-else-if="!forms.length" class="text-center py-16 bg-white border border-dashed border-gray-300 rounded-2xl">
          <p class="text-gray-500 text-sm">No application forms yet.</p>
          <button @click="openBuilder()" class="mt-3 text-sm text-primary-600 hover:text-primary-700 font-medium">Create your first form →</button>
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <article v-for="f in forms" :key="f.id"
            class="group bg-white border border-gray-200 hover:border-primary-400 hover:shadow-md transition rounded-xl overflow-hidden cursor-pointer"
            @click="openBuilder(f)">
            <div class="h-24 bg-gradient-to-br from-primary-500 to-primary-700 relative">
              <img v-if="f.banner_image_url" :src="mediaUrl(f.banner_image_url)" alt="" class="absolute inset-0 w-full h-full object-cover opacity-90"/>
              <span :class="statusBadgeClass(f.status)" class="absolute top-2 right-2 text-[10px] font-bold px-2 py-1 rounded-full uppercase tracking-wider">{{ f.status }}</span>
            </div>
            <div class="p-4">
              <h3 class="font-semibold text-gray-900 line-clamp-2">{{ f.title }}</h3>
              <p class="text-xs text-gray-500 mt-1 line-clamp-2">{{ f.description || 'No description' }}</p>
              <div class="flex items-center justify-between mt-3 text-xs text-gray-600">
                <span class="bg-gray-100 px-2 py-0.5 rounded-md">{{ prettifyRole(f.department_name) || '—' }}</span>
                <span class="font-medium">{{ f.price_ugx > 0 ? `UGX ${formatNum(f.price_ugx)}` : 'Free' }}</span>
              </div>
            </div>
          </article>
        </div>
      </div>

      <!-- ═════════ BUILDER MODE ═════════ -->
      <div v-else>
        <div class="flex items-center gap-3 mb-5">
          <button @click="closeBuilder" class="text-gray-500 hover:text-gray-900 flex items-center gap-1 text-sm">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5"/></svg>
            All forms
          </button>
          <span class="text-gray-300">/</span>
          <h2 class="text-xl font-semibold text-gray-900">{{ editing.id ? 'Edit Form' : 'New Form' }}</h2>
          <div class="ml-auto flex items-center gap-2">
            <button @click="saveForm" :disabled="saving" class="bg-primary-600 hover:bg-primary-700 disabled:opacity-60 text-white text-sm font-medium px-5 py-2 rounded-lg">
              {{ saving ? 'Saving…' : 'Save Form' }}
            </button>
            <button v-if="editing.id" @click="confirmArchive" class="text-red-600 hover:text-red-700 text-sm font-medium px-3 py-2">Archive</button>
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
          <!-- ── LEFT: Form metadata + field list ── -->
          <div class="lg:col-span-8 space-y-5">
            <!-- Metadata card -->
            <section class="bg-white border border-gray-200 rounded-xl p-5">
              <h3 class="text-sm font-semibold text-gray-700 uppercase tracking-wider mb-4">Form Details</h3>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-1">Title <span class="text-red-500">*</span></label>
                  <input v-model="editing.title" type="text" placeholder="e.g. Federation Registration"
                    class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/>
                </div>
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-1">Department <span class="text-red-500">*</span></label>
                  <select v-model="editing.department_id" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white">
                    <option value="">— Select department —</option>
                    <option v-for="d in departments" :key="d.id" :value="d.id">{{ prettifyRole(d.name) }}</option>
                  </select>
                  <p class="text-[10px] text-gray-400 mt-1">Sourced from system roles — adding a new role automatically lists it here.</p>
                </div>
                <div class="md:col-span-2">
                  <label class="block text-xs font-medium text-gray-600 mb-1">Description</label>
                  <textarea v-model="editing.description" rows="2" placeholder="Brief description shown at the top of the form"
                    class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"></textarea>
                </div>
                <div class="md:col-span-2">
                  <label class="block text-xs font-medium text-gray-600 mb-1">Banner Image</label>
                  <DropzoneUpload
                    v-model="editing.banner_image_url"
                    accept="image/*"
                    label="Drop a banner image here or click to upload"
                    hint="Shown at the top of the public form page"
                    preview-class="h-32"
                  />
                </div>
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-1">Price (UGX)</label>
                  <input v-model.number="editing.price_ugx" type="number" min="0" step="1000" placeholder="0 = free"
                    class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2"/>
                </div>
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-1">Status</label>
                  <select v-model="editing.status" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white">
                    <option value="DRAFT">Draft (hidden)</option>
                    <option value="OPEN">Open (accepting submissions)</option>
                    <option value="CLOSED">Closed (visible, no submissions)</option>
                    <option value="ARCHIVED">Archived</option>
                  </select>
                </div>
              </div>
            </section>

            <!-- Fields list -->
            <section class="bg-white border border-gray-200 rounded-xl p-5">
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold text-gray-700 uppercase tracking-wider">Form Fields</h3>
                <span class="text-xs text-gray-400">{{ editing.fields.length }} field{{ editing.fields.length === 1 ? '' : 's' }}</span>
              </div>

              <div v-if="!editing.fields.length" class="text-center py-10 bg-gray-50 rounded-lg border border-dashed border-gray-300">
                <p class="text-sm text-gray-500">No fields yet. Pick a field type on the right to get started.</p>
              </div>

              <div ref="fieldListEl" class="space-y-2">
                <div v-for="(field, idx) in editing.fields" :key="field._uid"
                  class="border border-gray-200 rounded-lg bg-white"
                  :data-idx="idx">
                  <!-- Header row -->
                  <div class="flex items-center gap-2 px-3 py-2 bg-gray-50 rounded-t-lg cursor-move drag-handle">
                    <svg class="w-4 h-4 text-gray-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 9h16.5m-16.5 6.75h16.5"/>
                    </svg>
                    <span class="text-xs font-mono text-gray-500 bg-white px-2 py-0.5 rounded">{{ fieldTypeLabel(field.field_type) }}</span>
                    <span v-if="field.is_required" class="text-[10px] uppercase tracking-wider text-red-500 font-bold">Required</span>
                    <div class="ml-auto flex items-center gap-1">
                      <button @click="moveField(idx, -1)" :disabled="idx === 0" class="text-gray-400 hover:text-gray-700 disabled:opacity-30 p-1" title="Move up">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M4.5 15.75l7.5-7.5 7.5 7.5"/></svg>
                      </button>
                      <button @click="moveField(idx, 1)" :disabled="idx === editing.fields.length - 1" class="text-gray-400 hover:text-gray-700 disabled:opacity-30 p-1" title="Move down">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5"/></svg>
                      </button>
                      <button @click="removeField(idx)" class="text-gray-400 hover:text-red-600 p-1" title="Delete field">
                        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"/></svg>
                      </button>
                    </div>
                  </div>

                  <!-- Body -->
                  <div class="px-3 py-3 space-y-3">
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                      <div>
                        <label class="block text-[11px] uppercase tracking-wider text-gray-400 mb-1">Label</label>
                        <input v-model="field.label" type="text" placeholder="Question label"
                          class="w-full text-sm border border-gray-200 rounded px-2.5 py-1.5"/>
                      </div>
                      <div>
                        <label class="block text-[11px] uppercase tracking-wider text-gray-400 mb-1">Field Key</label>
                        <input v-model="field.field_key" type="text" placeholder="auto-generated from label"
                          class="w-full text-sm border border-gray-200 rounded px-2.5 py-1.5 font-mono"/>
                      </div>
                      <div>
                        <label class="block text-[11px] uppercase tracking-wider text-gray-400 mb-1">Placeholder</label>
                        <input v-model="field.placeholder" type="text"
                          class="w-full text-sm border border-gray-200 rounded px-2.5 py-1.5"/>
                      </div>
                      <div>
                        <label class="block text-[11px] uppercase tracking-wider text-gray-400 mb-1">Help text</label>
                        <input v-model="field.help_text" type="text"
                          class="w-full text-sm border border-gray-200 rounded px-2.5 py-1.5"/>
                      </div>
                    </div>

                    <!-- Options (radio / checkbox / dropdown) -->
                    <div v-if="hasOptions(field.field_type)" class="bg-gray-50 rounded-lg p-3">
                      <label class="block text-[11px] uppercase tracking-wider text-gray-400 mb-2">Options</label>
                      <div class="space-y-1.5">
                        <div v-for="(opt, i) in field._options" :key="i" class="flex items-center gap-2">
                          <input v-model="opt.label" type="text" placeholder="Visible label"
                            class="flex-1 text-sm border border-gray-200 rounded px-2 py-1"/>
                          <input v-model="opt.value" type="text" placeholder="value"
                            class="w-32 text-sm border border-gray-200 rounded px-2 py-1 font-mono"/>
                          <button @click="field._options.splice(i, 1)" class="text-gray-300 hover:text-red-500 text-lg leading-none">×</button>
                        </div>
                      </div>
                      <button @click="field._options.push({ label: '', value: '' })" class="mt-2 text-xs text-primary-600 hover:text-primary-700 font-medium">+ Add option</button>
                    </div>

                    <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer pt-1">
                      <input v-model="field.is_required" type="checkbox" class="rounded border-gray-300 text-primary-600"/>
                      Required field
                    </label>
                  </div>
                </div>
              </div>
            </section>
          </div>

          <!-- ── RIGHT: Field palette ── -->
          <aside class="lg:col-span-4">
            <div class="sticky top-6 bg-white border border-gray-200 rounded-xl p-4">
              <h3 class="text-sm font-semibold text-gray-700 uppercase tracking-wider mb-3">Add a Field</h3>
              <div class="grid grid-cols-2 gap-2">
                <button v-for="t in FIELD_TYPES" :key="t.value"
                  @click="addField(t.value)"
                  class="text-left bg-gray-50 hover:bg-primary-50 hover:text-primary-700 border border-gray-200 hover:border-primary-300 rounded-lg p-3 transition">
                  <div class="text-xs font-medium">{{ t.label }}</div>
                </button>
              </div>

              <div v-if="editing.id" class="mt-5 pt-5 border-t border-gray-100">
                <h3 class="text-sm font-semibold text-gray-700 uppercase tracking-wider mb-3">Public URL</h3>
                <code class="block text-xs text-gray-700 bg-gray-50 px-2 py-1.5 rounded break-all">{{ publicUrl(editing.slug) }}</code>
                <p class="text-xs text-gray-400 mt-2">Visible to applicants when status is OPEN.</p>
              </div>
            </div>
          </aside>
        </div>
      </div>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import Sortable from 'sortablejs'
import DropzoneUpload from '@/components/ui/DropzoneUpload.vue'
import { mediaUrl } from '@/api/client.js'
import {
  listDepartments, adminListForms, adminGetForm, adminCreateForm,
  adminUpdateForm, adminDeleteForm, FIELD_TYPES, FIELD_TYPES_WITH_OPTIONS,
} from '@/api/forms'

const mode = ref('list')
const loading = ref(false)
const saving = ref(false)
const statusFilter = ref('')
const toast = ref('')
const forms = ref([])
const departments = ref([])
let sortable = null
const fieldListEl = ref(null)
let uidCounter = 0

const editing = reactive({
  id: '', slug: '', department_id: '', title: '', description: '',
  banner_image_url: '', price_ugx: 0, status: 'DRAFT', fields: [],
})

function showToast(msg) {
  toast.value = msg
  setTimeout(() => (toast.value = ''), 2500)
}

function formatNum(n) { return new Intl.NumberFormat().format(n || 0) }

function statusBadgeClass(s) {
  return {
    DRAFT:    'bg-gray-700 text-gray-100',
    OPEN:     'bg-green-500 text-white',
    CLOSED:   'bg-amber-500 text-white',
    ARCHIVED: 'bg-gray-400 text-white',
  }[s] || 'bg-gray-500 text-white'
}

function fieldTypeLabel(v) {
  return FIELD_TYPES.find(t => t.value === v)?.label || v
}

function prettifyRole(name) {
  if (!name) return ''
  return name.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

function hasOptions(t) { return FIELD_TYPES_WITH_OPTIONS.includes(t) }

function publicUrl(slug) {
  if (!slug) return ''
  return `${window.location.origin}/my-portal/forms/${slug}`
}

async function loadForms() {
  loading.value = true
  try {
    forms.value = await adminListForms(statusFilter.value)
  } catch (e) {
    showToast('Failed to load forms')
  } finally {
    loading.value = false
  }
}

async function loadDepartments() {
  try {
    departments.value = await listDepartments()
  } catch (e) { /* non-fatal */ }
}

function blankField(type) {
  const f = {
    _uid: ++uidCounter,
    field_key: '',
    field_type: type,
    label: '',
    placeholder: '',
    help_text: '',
    is_required: false,
    config: {},
    _options: [],
  }
  if (hasOptions(type)) {
    f._options = [{ label: 'Option 1', value: 'opt_1' }, { label: 'Option 2', value: 'opt_2' }]
  }
  return f
}

function fromServerField(sf) {
  let opts = []
  if (sf.config && typeof sf.config === 'object' && Array.isArray(sf.config.options)) {
    opts = sf.config.options.map(o => ({ label: o.label ?? '', value: o.value ?? '' }))
  }
  return {
    _uid: ++uidCounter,
    field_key:   sf.field_key  || '',
    field_type:  sf.field_type || 'short_text',
    label:       sf.label       || '',
    placeholder: sf.placeholder || '',
    help_text:   sf.help_text   || '',
    is_required: !!sf.is_required,
    config:      sf.config      || {},
    _options:    opts,
  }
}

function addField(type) {
  editing.fields.push(blankField(type))
}

function removeField(idx) {
  editing.fields.splice(idx, 1)
}

function moveField(idx, delta) {
  const j = idx + delta
  if (j < 0 || j >= editing.fields.length) return
  const [f] = editing.fields.splice(idx, 1)
  editing.fields.splice(j, 0, f)
}

function resetEditing() {
  editing.id = ''; editing.slug = ''
  editing.department_id = departments.value[0]?.id || ''
  editing.title = ''; editing.description = ''; editing.banner_image_url = ''
  editing.price_ugx = 0; editing.status = 'DRAFT'
  editing.fields = []
}

async function openBuilder(form = null) {
  if (!departments.value.length) await loadDepartments()
  resetEditing()
  if (form) {
    try {
      const full = await adminGetForm(form.id)
      editing.id = full.id
      editing.slug = full.slug || ''
      editing.department_id = full.department_id
      editing.title = full.title
      editing.description = full.description || ''
      editing.banner_image_url = full.banner_image_url || ''
      editing.price_ugx = Number(full.price_ugx) || 0
      editing.status = full.status
      editing.fields = (full.fields || []).map(fromServerField)
    } catch (e) {
      showToast('Failed to load form')
      return
    }
  }
  mode.value = 'builder'
  await nextTick()
  initSortable()
}

function closeBuilder() {
  destroySortable()
  mode.value = 'list'
  loadForms()
}

function initSortable() {
  destroySortable()
  if (!fieldListEl.value) return
  sortable = Sortable.create(fieldListEl.value, {
    handle: '.drag-handle',
    animation: 150,
    onEnd: (evt) => {
      if (evt.oldIndex === evt.newIndex) return
      const [moved] = editing.fields.splice(evt.oldIndex, 1)
      editing.fields.splice(evt.newIndex, 0, moved)
    },
  })
}

function destroySortable() {
  if (sortable) { sortable.destroy(); sortable = null }
}

function buildPayload() {
  return {
    department_id: editing.department_id,
    title: editing.title,
    description: editing.description,
    banner_image_url: editing.banner_image_url,
    price_ugx: Number(editing.price_ugx) || 0,
    status: editing.status,
    fields: editing.fields.map(f => ({
      field_key: f.field_key,
      field_type: f.field_type,
      label: f.label,
      placeholder: f.placeholder,
      help_text: f.help_text,
      is_required: !!f.is_required,
      config: hasOptions(f.field_type)
        ? { ...(f.config || {}), options: f._options.filter(o => o.label || o.value) }
        : (f.config || {}),
    })),
  }
}

async function saveForm() {
  if (!editing.title.trim()) { showToast('Title is required'); return }
  if (!editing.department_id) { showToast('Department is required'); return }
  saving.value = true
  try {
    const payload = buildPayload()
    const saved = editing.id
      ? await adminUpdateForm(editing.id, payload)
      : await adminCreateForm(payload)
    editing.id = saved.id
    editing.slug = saved.slug || editing.slug
    showToast(editing.id ? 'Form saved' : 'Form created')
  } catch (e) {
    showToast(e?.response?.data?.message || 'Save failed')
  } finally {
    saving.value = false
  }
}

async function confirmArchive() {
  if (!confirm('Archive this form? Existing submissions stay intact but no new ones will be accepted.')) return
  try {
    await adminDeleteForm(editing.id)
    showToast('Form archived')
    closeBuilder()
  } catch (e) {
    showToast('Archive failed')
  }
}

onMounted(async () => {
  await loadDepartments()
  await loadForms()
})
</script>

<style scoped>
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(-8px); }
.toast-enter-active, .toast-leave-active { transition: all .25s ease; }
</style>
