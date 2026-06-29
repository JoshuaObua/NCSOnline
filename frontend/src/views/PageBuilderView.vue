<template>
  <LayoutDefault title="Page Content Builder">
    <div class="grid min-h-[calc(100vh-9rem)] gap-4 lg:grid-cols-[18rem_minmax(0,1fr)_22rem]">
      <aside class="admin-card overflow-hidden">
        <div class="admin-card-header">
          <h2 class="font-semibold text-gray-900">Elements</h2>
        </div>
        <div class="admin-card-body space-y-2">
          <button v-for="type in palette" :key="type" type="button" class="palette-btn" @click="addNode(type)">
            <i :class="iconFor(type)"></i>
            <span>{{ ELEMENT_LABELS[type] }}</span>
          </button>
        </div>
      </aside>

      <main class="space-y-4 min-w-0">
        <section class="admin-card">
          <div class="admin-card-header flex-wrap gap-3">
            <div>
              <h1 class="font-bold text-gray-900">Nestable AST Builder</h1>
              <p class="mt-1 text-xs text-gray-500">Every block is a JSON node with id, type, settings, attributes, and children.</p>
            </div>
            <div class="flex gap-2">
              <button class="btn-light" @click="resetExample"><i class="icofont-magic"></i>Example</button>
              <button class="btn-light" @click="validate"><i class="icofont-check-circled"></i>Validate</button>
              <button class="btn-primary" :disabled="saving" @click="save"><i class="icofont-save"></i>{{ saving ? 'Saving' : 'Save AST' }}</button>
            </div>
          </div>
          <div class="admin-card-body">
            <PageBuilderRenderer :node="root" :form-state="formState" :action-context="actionContext" />
          </div>
        </section>

        <section class="admin-card">
          <div class="admin-card-header">
            <h2 class="font-semibold text-gray-900">AST JSON</h2>
            <button class="btn-light" @click="syncFromJson"><i class="icofont-code"></i>Apply JSON</button>
          </div>
          <div class="admin-card-body">
            <textarea v-model="jsonText" class="json-editor" spellcheck="false"></textarea>
          </div>
        </section>
      </main>

      <aside class="admin-card overflow-hidden">
        <div class="admin-card-header">
          <h2 class="font-semibold text-gray-900">Tree</h2>
        </div>
        <div class="admin-card-body border-b border-gray-100">
          <TreeNode :node="root" :selected-id="selectedId" @select="selectedId = $event" @remove="removeNode" />
        </div>
        <div class="admin-card-body space-y-4">
          <template v-if="selected">
            <div>
              <label class="field-label">Type</label>
              <input :value="selected.type" class="field-input bg-gray-50" readonly />
            </div>
            <div>
              <label class="field-label">Settings JSON</label>
              <textarea v-model="selectedSettingsText" class="side-json" spellcheck="false" @blur="applySelectedJson('settings')"></textarea>
            </div>
            <div>
              <label class="field-label">Attributes JSON</label>
              <textarea v-model="selectedAttributesText" class="side-json" spellcheck="false" @blur="applySelectedJson('attributes')"></textarea>
            </div>
          </template>
          <p v-else class="text-sm text-gray-400">Select a node to edit its settings and attributes.</p>
        </div>
      </aside>
    </div>

    <p v-if="message" class="fixed bottom-5 right-5 rounded-lg border px-4 py-3 text-sm shadow-lg" :class="error ? 'border-red-200 bg-red-50 text-red-700' : 'border-emerald-200 bg-emerald-50 text-emerald-700'">{{ message }}</p>
  </LayoutDefault>
</template>

<script setup>
import { computed, defineComponent, h, onMounted, reactive, ref, watch } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import PageBuilderRenderer from '@/components/page-builder/PageBuilderRenderer.vue'
import { adminUpdateSettings, getSettings } from '@/api/cms.js'
import { createNode, ELEMENT_LABELS, exampleAst, NODE_TYPES, validateTree } from '@/page-builder/schema.js'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'

const breadcrumbStore = useBreadcrumbStore()
const root = ref(createNode('layout.row'))
const selectedId = ref(root.value.id)
const jsonText = ref('')
const selectedSettingsText = ref('{}')
const selectedAttributesText = ref('{}')
const saving = ref(false)
const message = ref('')
const error = ref(false)
const formState = reactive({})
const actionContext = { formState }

const palette = NODE_TYPES.filter(t => t !== 'layout.column')
const selected = computed(() => findNode(root.value, selectedId.value))

watch(root, syncJson, { deep: true, immediate: true })
watch(selected, node => {
  selectedSettingsText.value = JSON.stringify(node?.settings || {}, null, 2)
  selectedAttributesText.value = JSON.stringify(node?.attributes || {}, null, 2)
}, { immediate: true })

function addNode(type) {
  const parent = selected.value || root.value
  if (['ui.typography', 'ui.image', 'form.dropdown', 'form.checkbox', 'form.toggle'].includes(parent.type)) {
    notify('Select a container node before adding children', true)
    return
  }
  parent.children.push(createNode(type))
  selectedId.value = parent.children[parent.children.length - 1].id
}

function removeNode(id) {
  if (id === root.value.id) {
    notify('Root node cannot be deleted', true)
    return
  }
  removeChild(root.value, id)
  selectedId.value = root.value.id
}

function validate() {
  const errors = validateTree(root.value)
  notify(errors.length ? errors.slice(0, 3).join(' | ') : 'AST is valid', !!errors.length)
}

async function save() {
  saving.value = true
  try {
    validate()
    await adminUpdateSettings('page_builder_demo', { value: root.value })
    notify('Page builder AST saved')
  } catch (e) {
    notify(e.response?.data?.error?.message || 'Could not save AST', true)
  } finally {
    saving.value = false
  }
}

async function load() {
  try {
    const r = await getSettings('page_builder_demo')
    const value = r.data?.data?.value?.value || r.data?.data?.value
    if (value?.id) root.value = value
    selectedId.value = root.value.id
  } catch {
    root.value = structuredClone(exampleAst)
  }
}

function resetExample() {
  root.value = structuredClone(exampleAst)
  selectedId.value = root.value.id
}

function syncJson() {
  jsonText.value = JSON.stringify(root.value, null, 2)
}

function syncFromJson() {
  try {
    const parsed = JSON.parse(jsonText.value)
    const errors = validateTree(parsed)
    if (errors.length) {
      notify(errors[0], true)
      return
    }
    root.value = parsed
    selectedId.value = parsed.id
    notify('JSON applied')
  } catch (e) {
    notify(e.message, true)
  }
}

function applySelectedJson(key) {
  if (!selected.value) return
  try {
    selected.value[key] = JSON.parse(key === 'settings' ? selectedSettingsText.value : selectedAttributesText.value)
  } catch (e) {
    notify(e.message, true)
  }
}

function notify(text, isError = false) {
  message.value = text
  error.value = isError
  setTimeout(() => { if (message.value === text) message.value = '' }, 3500)
}

function findNode(node, id) {
  if (!node) return null
  if (node.id === id) return node
  for (const child of node.children || []) {
    const found = findNode(child, id)
    if (found) return found
  }
  for (const panel of node.attributes?.panels || []) {
    for (const child of panel.children || []) {
      const found = findNode(child, id)
      if (found) return found
    }
  }
  for (const tab of node.attributes?.tabs || []) {
    for (const child of tab.children || []) {
      const found = findNode(child, id)
      if (found) return found
    }
  }
  return null
}

function removeChild(node, id) {
  const idx = (node.children || []).findIndex(child => child.id === id)
  if (idx >= 0) {
    node.children.splice(idx, 1)
    return true
  }
  for (const child of node.children || []) if (removeChild(child, id)) return true
  return false
}

function iconFor(type) {
  if (type.startsWith('layout.')) return 'icofont-layout'
  if (type.startsWith('form.')) return 'icofont-ui-check'
  if (type === 'ui.image') return 'icofont-image'
  return 'icofont-text-height'
}

const TreeNode = defineComponent({
  props: { node: Object, selectedId: String },
  emits: ['select', 'remove'],
  setup(props, { emit }) {
    return () => h('div', { class: 'text-sm' }, [
      h('div', {
        class: [
          'flex items-center gap-2 rounded-md px-2 py-1 cursor-pointer',
          props.selectedId === props.node.id ? 'bg-primary-50 text-primary-700' : 'hover:bg-gray-50',
        ],
        onClick: () => emit('select', props.node.id),
      }, [
        h('span', { class: 'flex-1 truncate font-medium' }, ELEMENT_LABELS[props.node.type] || props.node.type),
        h('button', { class: 'text-gray-300 hover:text-red-600', onClick: e => { e.stopPropagation(); emit('remove', props.node.id) } }, 'x'),
      ]),
      h('div', { class: 'ml-4 border-l border-gray-100 pl-2' }, (props.node.children || []).map(child =>
        h(TreeNode, { node: child, selectedId: props.selectedId, onSelect: id => emit('select', id), onRemove: id => emit('remove', id) })
      )),
    ])
  },
})

onMounted(() => {
  breadcrumbStore.set('Page Content Builder', [{ label: 'CMS' }, { label: 'Page Content Builder' }])
  load()
})
</script>

<style scoped>
.palette-btn { display: flex; width: 100%; align-items: center; gap: .5rem; border-radius: .5rem; border: 1px solid #e5e7eb; padding: .625rem .75rem; text-align: left; font-size: .875rem; font-weight: 600; color: #374151; }
.palette-btn:hover { border-color: #1f4f82; background: #eff6ff; color: #1f4f82; }
.btn-primary { display: inline-flex; align-items: center; gap: .5rem; border-radius: .5rem; background: #1f4f82; padding: .55rem .85rem; color: white; font-size: .875rem; font-weight: 700; }
.btn-light { display: inline-flex; align-items: center; gap: .5rem; border-radius: .5rem; border: 1px solid #e5e7eb; background: white; padding: .55rem .85rem; color: #374151; font-size: .875rem; font-weight: 700; }
.json-editor { min-height: 18rem; width: 100%; resize: vertical; border-radius: .5rem; border: 1px solid #d1d5db; background: #111827; padding: 1rem; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: .75rem; color: #bbf7d0; }
.side-json { min-height: 9rem; width: 100%; resize: vertical; border-radius: .5rem; border: 1px solid #d1d5db; padding: .75rem; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: .75rem; }
.field-label { display: block; margin-bottom: .25rem; font-size: .75rem; font-weight: 700; color: #4b5563; }
.field-input { width: 100%; border-radius: .5rem; border: 1px solid #d1d5db; padding: .5rem .75rem; font-size: .875rem; }
</style>
