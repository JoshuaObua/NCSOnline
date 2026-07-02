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
              <input v-model="model.slug" class="form-control" :readonly="!manualSlug" pattern="^[a-z0-9-_]+$" required placeholder="clear-page-name" />
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

        <div class="form-section">
          <div class="section-title">Page Building Blocks</div>
          <div class="builder-toolbar">
            <button v-for="type in blockTypes" :key="type.type" type="button" class="btn btn-light btn-sm" @click="addBlock(type.type)">
              <i :class="type.icon"></i> {{ type.label }}
            </button>
          </div>
          <div class="builder-blocks">
            <article v-for="(block, index) in builder.blocks" :key="block.id" class="builder-block">
              <div class="builder-block-head">
                <div><i :class="blockIcon(block.type)"></i><strong>{{ blockLabel(block.type) }}</strong></div>
                <div class="builder-actions">
                  <button type="button" class="btn btn-sm btn-light" :disabled="index === 0" @click="moveBlock(index, -1)"><i class="fas fa-arrow-up"></i></button>
                  <button type="button" class="btn btn-sm btn-light" :disabled="index === builder.blocks.length - 1" @click="moveBlock(index, 1)"><i class="fas fa-arrow-down"></i></button>
                  <button type="button" class="btn btn-sm btn-danger" @click="removeBlock(index)"><i class="fas fa-trash"></i></button>
                </div>
              </div>

              <div v-if="block.type === 'section'" class="row">
                <div class="form-group col-md-4"><label>Kicker</label><input v-model="block.kicker" class="form-control" /></div>
                <div class="form-group col-md-8"><label>Title</label><input v-model="block.title" class="form-control" /></div>
                <div class="form-group col-12"><label>Text</label><textarea v-model="block.text" class="form-control compact-textarea"></textarea></div>
              </div>

              <div v-else-if="block.type === 'image'" class="row">
                <div class="form-group col-md-8"><label>Image URL</label><input v-model="block.src" class="form-control" placeholder="/uploads/image.webp" /></div>
                <div class="form-group col-md-4"><label>Alt text</label><input v-model="block.alt" class="form-control" /></div>
                <div class="form-group col-12"><label>Caption</label><input v-model="block.caption" class="form-control" /></div>
              </div>

              <div v-else-if="block.type === 'row'" class="nested-builder">
                <label>Columns</label>
                <div class="builder-toolbar">
                  <button type="button" class="btn btn-light btn-sm" @click="addColumn(block)"><i class="fas fa-columns"></i> Add column</button>
                </div>
                <div class="builder-columns">
                  <div v-for="(column, colIndex) in block.columns" :key="column.id" class="builder-column">
                    <div class="builder-block-head">
                      <strong>Column {{ colIndex + 1 }}</strong>
                      <button type="button" class="btn btn-sm btn-danger" @click="block.columns.splice(colIndex, 1)">Remove</button>
                    </div>
                    <label>Title<input v-model="column.title" class="form-control" /></label>
                    <label>Text<textarea v-model="column.text" class="form-control compact-textarea"></textarea></label>
                    <label>Icon class<input v-model="column.icon" class="form-control" placeholder="icofont-trophy" /></label>
                  </div>
                </div>
              </div>

              <div v-else-if="block.type === 'accordion'" class="nested-builder">
                <div class="builder-toolbar"><button type="button" class="btn btn-light btn-sm" @click="addAccordionItem(block)"><i class="fas fa-plus"></i> Add item</button></div>
                <div v-for="(item, itemIndex) in block.items" :key="item.id" class="builder-mini-card">
                  <label>Question / heading<input v-model="item.title" class="form-control" /></label>
                  <label>Answer / body<textarea v-model="item.text" class="form-control compact-textarea"></textarea></label>
                  <button type="button" class="btn btn-sm btn-danger" @click="block.items.splice(itemIndex, 1)">Remove item</button>
                </div>
              </div>

              <div v-else-if="block.type === 'dropdown'" class="nested-builder">
                <label>Dropdown label<input v-model="block.label" class="form-control" /></label>
                <label>Dropdown body<textarea v-model="block.text" class="form-control compact-textarea"></textarea></label>
              </div>

              <div v-else-if="block.type === 'icon-card'" class="row">
                <div class="form-group col-md-4"><label>Icon class</label><input v-model="block.icon" class="form-control" placeholder="icofont-medal" /></div>
                <div class="form-group col-md-8"><label>Title</label><input v-model="block.title" class="form-control" /></div>
                <div class="form-group col-12"><label>Text</label><textarea v-model="block.text" class="form-control compact-textarea"></textarea></div>
              </div>

              <div v-else-if="block.type === 'html'">
                <label>Raw HTML<textarea v-model="block.html" class="form-control code-textarea" placeholder="<div>Custom HTML</div>"></textarea></label>
              </div>
            </article>
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
import { reactive, ref, watch } from 'vue'
import DropzoneUpload from '@/components/cms/DropzoneUpload.vue'

const props = defineProps({ model: { type: Object, required: true } })
const emit = defineEmits(['save'])
const manualSlug = ref(false)
const blockTypes = [
  { type:'section', label:'Text Section', icon:'icofont-align-left' },
  { type:'image', label:'Image', icon:'icofont-image' },
  { type:'row', label:'Rows/Columns', icon:'icofont-columns' },
  { type:'accordion', label:'Accordion', icon:'icofont-list' },
  { type:'dropdown', label:'Dropdown', icon:'icofont-rounded-down' },
  { type:'icon-card', label:'Icon Card', icon:'icofont-star' },
  { type:'html', label:'Raw HTML', icon:'icofont-code' },
]
const builder = reactive({ version: 1, blocks: [] })

watch(() => props.model.content, hydrateBuilder, { immediate:true })
watch(() => props.model.id, () => { manualSlug.value = !!props.model.id; hydrateBuilder() })
watch(builder, () => { props.model.content = JSON.stringify({ type:'ncs-page-builder', version:1, blocks: builder.blocks }) }, { deep:true })

function uid() { return `block_${Math.random().toString(36).slice(2, 10)}` }
function slugify(value) { return String(value || '').toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '') }
function syncSlug() { if (!manualSlug.value) props.model.slug = slugify(props.model.title) }
function submit() {
  props.model.category = 'page'
  props.model.content = JSON.stringify({ type:'ncs-page-builder', version:1, blocks: builder.blocks })
  emit('save')
}
function hydrateBuilder() {
  try {
    const parsed = JSON.parse(props.model.content || '{}')
    builder.blocks = parsed?.type === 'ncs-page-builder' && Array.isArray(parsed.blocks) ? parsed.blocks : defaultBlocks()
  } catch {
    builder.blocks = props.model.content ? [{ id:uid(), type:'html', html:props.model.content }] : defaultBlocks()
  }
}
function defaultBlocks() {
  return [{ id:uid(), type:'section', kicker:'Static Page', title:'Page section title', text:'Start writing your page content here.' }]
}
function addBlock(type) { builder.blocks.push(createBlock(type)) }
function createBlock(type) {
  const base = { id:uid(), type }
  if (type === 'section') return { ...base, kicker:'', title:'New section', text:'' }
  if (type === 'image') return { ...base, src:'', alt:'', caption:'' }
  if (type === 'row') return { ...base, columns:[{ id:uid(), title:'Column title', text:'', icon:'icofont-check-circled' }, { id:uid(), title:'Column title', text:'', icon:'icofont-check-circled' }] }
  if (type === 'accordion') return { ...base, items:[{ id:uid(), title:'Accordion item', text:'' }] }
  if (type === 'dropdown') return { ...base, label:'Dropdown title', text:'' }
  if (type === 'icon-card') return { ...base, icon:'icofont-star', title:'Icon card title', text:'' }
  return { ...base, html:'<p>Custom HTML block</p>' }
}
function blockLabel(type) { return blockTypes.find(item => item.type === type)?.label || 'Block' }
function blockIcon(type) { return blockTypes.find(item => item.type === type)?.icon || 'icofont-ui-note' }
function removeBlock(index) { builder.blocks.splice(index, 1) }
function moveBlock(index, direction) {
  const next = index + direction
  if (next < 0 || next >= builder.blocks.length) return
  const [item] = builder.blocks.splice(index, 1)
  builder.blocks.splice(next, 0, item)
}
function addColumn(block) { block.columns.push({ id:uid(), title:'Column title', text:'', icon:'icofont-check-circled' }) }
function addAccordionItem(block) { block.items.push({ id:uid(), title:'Accordion item', text:'' }) }
</script>

<style scoped>
.static-page-builder{display:grid;gap:20px}.page-builder-card{border:0;border-radius:3px;box-shadow:0 4px 25px rgba(0,0,0,.1)}.form-section{padding:18px 0;border-bottom:1px solid #f4f6f9}.section-title{margin:0 0 18px;color:#34395e;font-size:14px;font-weight:700}.compact-textarea{min-height:96px}.code-textarea{min-height:180px;font-family:Consolas,monospace}.builder-toolbar{display:flex;flex-wrap:wrap;gap:8px;margin-bottom:16px}.builder-blocks{display:grid;gap:16px}.builder-block,.builder-mini-card,.builder-column{border:1px solid #e4e6fc;border-radius:3px;background:#fdfdff;padding:16px}.builder-block-head{display:flex;align-items:center;justify-content:space-between;gap:12px;margin-bottom:14px}.builder-block-head>div{display:flex;align-items:center;gap:8px}.builder-block-head i{color:#6777ef}.builder-actions{display:flex;gap:6px}.builder-columns{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:12px}.nested-builder{display:grid;gap:12px}.builder-mini-card{display:grid;gap:10px;margin-bottom:10px}.page-check-field{align-self:end}@media(max-width:900px){.builder-columns{grid-template-columns:1fr}.builder-block-head{align-items:flex-start;flex-direction:column}}:global(.dark) .page-builder-card,:global(.dark) .builder-block,:global(.dark) .builder-mini-card,:global(.dark) .builder-column{background:#1f2937!important;border-color:#334155!important;color:#e5e7eb!important}:global(.dark) .section-title{color:#f8fafc!important}
</style>
