<template>
  <article
    class="pb-node"
    :class="{ 'is-container': spec.container, 'is-drop-target': dropTarget === node.id }"
    draggable="true"
    @dragstart.stop="$emit('drag-start', node.id)"
    @dragover.prevent.stop="handleDragOver"
    @dragleave.stop="$emit('drop-target', '')"
    @drop.prevent.stop="handleDrop"
  >
    <header class="pb-node-head">
      <button type="button" class="pb-handle" title="Move"><i class="icofont-drag"></i></button>
      <strong><i :class="spec.icon"></i>{{ spec.label }}</strong>
      <div class="pb-actions">
        <button type="button" title="Duplicate" @click="$emit('duplicate', node.id)"><i class="icofont-copy"></i></button>
        <button type="button" title="Delete" @click="$emit('delete', node.id)"><i class="icofont-trash"></i></button>
      </div>
    </header>

    <div class="pb-props">
      <label v-for="field in editableFields" :key="field" :class="{ 'code-field': node.type === 'raw_html' && field === 'html' }">
        {{ labelFor(field) }}
        <textarea v-if="isLongField(field)" v-model="node.props[field]" :rows="node.type === 'raw_html' ? 8 : 2" :spellcheck="node.type === 'raw_html' ? false : undefined"></textarea>
        <input v-else v-model="node.props[field]" :type="inputType(field)" />
      </label>
    </div>

    <div v-if="spec.container" class="pb-children" :style="childGridStyle">
      <PageBuilderNode
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :dragging-id="draggingId"
        :drop-target="dropTarget"
        @drag-start="$emit('drag-start', $event)"
        @drop-node="$emit('drop-node', $event)"
        @drop-target="$emit('drop-target', $event)"
        @delete="$emit('delete', $event)"
        @duplicate="$emit('duplicate', $event)"
      />
      <div v-if="!node.children?.length" class="pb-empty">Drop components here</div>
    </div>
  </article>
</template>

<script setup>
import { computed } from 'vue'
import { blockMap, canAccept } from '@/utils/pageBuilderRegistry.js'

const props = defineProps({
  node: { type: Object, required: true },
  draggingId: { type: String, default: '' },
  dropTarget: { type: String, default: '' },
})
const emit = defineEmits(['drag-start', 'drop-node', 'drop-target', 'delete', 'duplicate'])

const spec = computed(() => blockMap[props.node.type] || blockMap.paragraph)
const editableFields = computed(() => Object.keys(props.node.props || {}).filter(key => !['editor'].includes(key)))
const childGridStyle = computed(() => {
  if (props.node.type === 'grid' || props.node.type === 'sidebar_layout') return { gridTemplateColumns: props.node.props.template || '1fr 1fr' }
  if (props.node.type === 'button_group') return { gridTemplateColumns: props.node.props.direction === 'vertical' ? '1fr' : 'repeat(auto-fit,minmax(120px,max-content))' }
  return {}
})

function handleDragOver() {
  if (!props.draggingId || props.draggingId === props.node.id) return
  emit('drop-target', props.node.id)
}

function handleDrop() {
  if (!props.draggingId || props.draggingId === props.node.id) return
  emit('drop-node', { parentId: props.node.id })
  emit('drop-target', '')
}

function labelFor(field) {
  return field.replace(/([A-Z])/g, ' $1').replace(/^./, c => c.toUpperCase())
}

function isLongField(field) {
  return ['text', 'html', 'alt', 'cite'].includes(field)
}

function inputType(field) {
  if (field === 'value') return 'number'
  if (field === 'lazy' || field === 'lightbox' || field === 'dismissible') return 'checkbox'
  return 'text'
}
</script>

<style scoped>
.pb-node{border:1px solid #dbe3ef;border-radius:8px;background:#fff;padding:.75rem;display:grid;gap:.65rem}.pb-node:hover{outline:2px solid rgba(103,119,239,.18)}.pb-node.is-drop-target{border-color:#6777ef;background:#f6f7ff}.pb-node-head{display:flex;align-items:center;justify-content:space-between;gap:.7rem}.pb-node-head strong{display:flex;align-items:center;gap:.45rem;color:#34395e}.pb-node-head strong i{color:#6777ef}.pb-handle,.pb-actions button{width:30px;height:30px;border:1px solid #e5e7eb;border-radius:6px;background:#f8fafc;color:#475569}.pb-actions{display:flex;gap:.35rem}.pb-props{display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));gap:.55rem}.pb-props label{display:grid;gap:.25rem;color:#64748b;font-size:.72rem;font-weight:800;text-transform:uppercase}.pb-props input,.pb-props textarea{border:1px solid #d1d5db;border-radius:6px;padding:.45rem .55rem;color:#111827;font-size:.85rem;text-transform:none}.pb-props .code-field{grid-column:1/-1}.pb-props .code-field textarea{min-height:180px;font-family:Consolas,"Courier New",monospace;line-height:1.55;tab-size:2}.pb-children{display:grid;gap:.75rem;border:1px dashed #cbd5e1;border-radius:8px;background:#f8fafc;padding:.75rem}.pb-empty{min-height:54px;display:grid;place-items:center;color:#94a3b8;font-weight:800}
</style>
