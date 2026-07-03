<template>
  <form class="static-page-builder" @submit.prevent="submit">
    <div class="card page-builder-card">
      <div class="card-header">
        <h4>{{ model.id ? 'Edit Static Page' : 'Create Static Page' }}</h4>
        <div class="card-header-action">
          <button type="button" class="btn btn-icon icon-left btn-light" @click="addBlock('section')"><i class="fas fa-plus"></i> Section</button>
          <button type="submit" class="btn btn-icon icon-left btn-primary"><i class="fas fa-save"></i> {{ model.id ? 'Update' : 'Create' }}</button>
        </div>
      </div>

      <div class="card-body">
        <div class="form-section">
          <div class="section-title mt-0">Page Details</div>
          <div class="row">
            <div class="form-group col-lg-8">
              <label>Title</label>
              <input v-model="model.title" class="form-control" required placeholder="Enter page title" @input="syncSlug" />
            </div>
            <div class="form-group col-lg-4">
              <label>Status</label>
              <select v-model="model.status" class="form-control selectric">
                <option value="draft">Draft</option>
                <option value="approved">Approved</option>
                <option value="published">Published</option>
              </select>
            </div>
            <div class="form-group col-lg-8">
              <label>Slug</label>
              <input v-model="model.slug" class="form-control" :readonly="!manualSlug" pattern="^[a-z0-9_-]+$" required placeholder="clear-page-name" />
              <small class="form-text text-muted">Public URL: /pages/{{ model.slug || 'clear-page-name' }}</small>
            </div>
            <div class="form-group col-lg-4 page-check-field">
              <label class="d-block">Slug Control</label>
              <div class="custom-control custom-checkbox">
                <input id="static-page-manual-slug" v-model="manualSlug" type="checkbox" class="custom-control-input" />
                <label class="custom-control-label" for="static-page-manual-slug">Manual slug override</label>
              </div>
            </div>
            <div class="form-group col-12">
              <label>Intro / excerpt</label>
              <textarea v-model="model.excerpt" class="form-control compact-textarea" maxlength="500" placeholder="Short page intro displayed below the breadcrumb"></textarea>
            </div>
          </div>
        </div>

        <div class="form-section">
          <div class="section-title">Breadcrumb Background Image</div>
          <DropzoneUpload
            v-model="model.cover_image_url"
            label="breadcrumb background"
            hint="This image appears behind the page title and breadcrumb."
            accept="image/png,image/jpeg,image/webp"
            :max-image-mb="5"
          />
        </div>

        <div class="form-section builder-shell">
          <div class="section-title">Advanced Nestable Page Builder</div>
          <div class="builder-layout">
            <aside class="builder-palette">
              <div v-for="(items, category) in palette" :key="category" class="palette-group">
                <h5>{{ category }}</h5>
                <button
                  v-for="item in items"
                  :key="item.type"
                  type="button"
                  draggable="true"
                  @dragstart="dragPalette(item.type)"
                  @click="addBlock(item.type)"
                >
                  <i :class="item.icon"></i>{{ item.label }}
                </button>
              </div>
            </aside>

            <section
              class="builder-canvas"
              :class="{ 'is-drop-target': dropTarget === 'root' }"
              @dragover.prevent="dropTarget = 'root'"
              @dragleave="dropTarget = ''"
              @drop.prevent="dropOnRoot"
            >
              <PageBuilderNode
                v-for="node in builder.blocks"
                :key="node.id"
                :node="node"
                :dragging-id="draggingId"
                :drop-target="dropTarget"
                @drag-start="dragExisting"
                @drop-node="dropNode"
                @drop-target="dropTarget = $event"
                @delete="deleteNode"
                @duplicate="duplicateExisting"
              />
              <div v-if="!builder.blocks.length" class="canvas-empty">Drop a section or component here</div>
            </section>
          </div>
        </div>

        <div class="form-section seo-card">
          <div class="section-title">SEO Meta</div>
          <div class="row">
            <div class="form-group col-lg-6"><label>SEO title</label><input v-model="model.meta_title" class="form-control" maxlength="180" /></div>
            <div class="form-group col-lg-6"><label>Focus keywords</label><input v-model="model.focus_keywords" class="form-control" placeholder="sports, ncs, uganda" /></div>
            <div class="form-group col-12"><label>Meta description</label><textarea v-model="model.meta_description" class="form-control compact-textarea" maxlength="160"></textarea></div>
          </div>
        </div>
      </div>
    </div>
  </form>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import DropzoneUpload from '@/components/cms/DropzoneUpload.vue'
import PageBuilderNode from '@/components/cms/PageBuilderNode.vue'
import { blockMap, canAccept, createNode, duplicateNode, findNode, insertNode, pageBuilderVersion, removeNode, visibleBlocks } from '@/utils/pageBuilderRegistry.js'

const props = defineProps({ model: { type: Object, required: true } })
const emit = defineEmits(['save'])
const manualSlug = ref(false)
const draggingId = ref('')
const dropTarget = ref('')
const builder = reactive({ version: pageBuilderVersion, blocks: [] })

const palette = computed(() => visibleBlocks.reduce((groups, block) => {
  groups[block.category] ||= []
  groups[block.category].push(block)
  return groups
}, {}))

watch(() => props.model.content, hydrateBuilder, { immediate: true })
watch(() => props.model.id, () => { manualSlug.value = !!props.model.id; hydrateBuilder() })
watch(builder, syncContent, { deep: true })

function slugify(value) { return String(value || '').toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '') }
function syncSlug() { if (!manualSlug.value) props.model.slug = slugify(props.model.title) }
function submit() {
  props.model.category = 'page'
  syncContent()
  emit('save')
}
function syncContent() {
  props.model.content = JSON.stringify({ type: 'ncs-page-builder', version: pageBuilderVersion, blocks: builder.blocks })
}
function hydrateBuilder() {
  try {
    const parsed = JSON.parse(props.model.content || '{}')
    if (parsed?.type === 'ncs-page-builder' && Array.isArray(parsed.blocks)) {
      builder.blocks = parsed.blocks.map(normalizeNode)
      return
    }
  } catch {
    if (props.model.content) {
      builder.blocks = [createNode('summernote_text')]
      builder.blocks[0].props.html = props.model.content
      return
    }
  }
  builder.blocks = [createNode('section')]
}
function normalizeNode(node) {
  if (node.props && Array.isArray(node.children)) return node
  const migrated = createNode(node.type === 'row' ? 'grid' : node.type === 'icon-card' ? 'card' : node.type || 'paragraph')
  migrated.id = node.id || migrated.id
  migrated.props = { ...migrated.props, ...legacyProps(node) }
  if (node.columns) migrated.children = node.columns.map(column => ({ ...createNode('column'), props: { width: '1fr' }, children: [legacyTextNode(column.title, column.text)] }))
  if (node.items) migrated.children = node.items.map(item => ({ ...createNode('accordion_panel'), props: { title: item.title || 'Pane' }, children: [legacyTextNode('', item.text)] }))
  return migrated
}
function legacyProps(node) {
  if (node.type === 'section') return { text: node.text, html: node.text, title: node.title }
  if (node.type === 'image') return { src: node.src, alt: node.alt, caption: node.caption }
  if (node.type === 'html') return { html: node.html }
  return node
}
function legacyTextNode(title, text) {
  const n = createNode(title ? 'h3' : 'paragraph')
  n.props.text = title || text || ''
  if (title && text) n.children = [createNode('paragraph')]
  if (title && text) n.children[0].props.text = text
  return n
}
function addBlock(type) { builder.blocks.push(createNode(type)) }
function dragPalette(type) { draggingId.value = `palette:${type}` }
function dragExisting(id) { draggingId.value = id }
function dropOnRoot() {
  if (!draggingId.value) return
  const node = takeDraggedNode()
  if (node) builder.blocks.push(node)
  clearDrag()
}
function dropNode({ parentId }) {
  const parent = findNode(builder.blocks, parentId)
  const node = takeDraggedNode()
  if (!parent || !node || !canAccept(parent.type, node.type)) {
    clearDrag()
    return
  }
  insertNode(builder.blocks, parentId, node)
  clearDrag()
}
function takeDraggedNode() {
  if (draggingId.value.startsWith('palette:')) return createNode(draggingId.value.replace('palette:', ''))
  return removeNode(builder.blocks, draggingId.value)
}
function deleteNode(id) { removeNode(builder.blocks, id) }
function duplicateExisting(id) {
  const node = findNode(builder.blocks, id)
  if (node) builder.blocks.push(duplicateNode(node))
}
function clearDrag() {
  draggingId.value = ''
  dropTarget.value = ''
}
</script>

<style scoped>
.static-page-builder{display:grid;gap:20px}.page-builder-card{border:0;border-radius:3px;box-shadow:0 4px 25px rgba(0,0,0,.1)}.form-section{padding:18px 0;border-bottom:1px solid #f4f6f9}.section-title{margin:0 0 18px;color:#34395e;font-size:14px;font-weight:700}.compact-textarea{min-height:96px}.page-check-field{align-self:end}.builder-layout{display:grid;grid-template-columns:280px minmax(0,1fr);gap:16px}.builder-palette{border:1px solid #e4e6fc;border-radius:8px;background:#fbfcff;padding:12px;height:max-content;position:sticky;top:90px}.palette-group{display:grid;gap:7px;margin-bottom:14px}.palette-group h5{font-size:12px;text-transform:uppercase;letter-spacing:.04em;color:#6777ef;margin:0}.palette-group button{display:flex;align-items:center;gap:8px;border:1px solid #e5e7eb;border-radius:6px;background:white;color:#34395e;text-align:left;padding:8px 10px;font-weight:700;font-size:12px}.palette-group button:hover{border-color:#6777ef;color:#6777ef}.builder-canvas{display:grid;gap:12px;min-height:320px;border:1px dashed #b8c2d6;border-radius:8px;background:#f8fafc;padding:14px}.builder-canvas.is-drop-target{border-color:#6777ef;background:#f2f4ff}.canvas-empty{display:grid;place-items:center;min-height:220px;color:#94a3b8;font-weight:800}@media(max-width:1000px){.builder-layout{grid-template-columns:1fr}.builder-palette{position:static}}:global(.dark) .page-builder-card,:global(.dark) .builder-palette,:global(.dark) .builder-canvas{background:#1f2937!important;border-color:#334155!important;color:#e5e7eb!important}:global(.dark) .section-title{color:#f8fafc!important}
</style>
