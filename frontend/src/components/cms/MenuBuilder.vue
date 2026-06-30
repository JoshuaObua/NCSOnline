<template>
  <div class="menu-builder">
    <div class="menu-builder__toolbar">
      <button type="button" @click="addRootItem">Add root item</button>
      <span>Max depth: {{ MAX_MENU_DEPTH }}</span>
    </div>
    <MenuBranch
      :items="localItems"
      :depth="1"
      @changed="handleBranchChange"
      @edit="editItem"
      @remove="removeItem"
      @add-child="addChild"
      @move="moveItem"
    />
    <p class="menu-builder__hint">Keyboard: use Move up/down, Indent, and Outdent controls on each row. Dragging supports nesting up to three levels.</p>
  </div>
</template>

<script setup>
import Sortable from 'sortablejs'
import { computed, defineComponent, h, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { MAX_MENU_DEPTH, cloneTree, createMenuItem, debounce, maxDepth, normalizeMenuTree, sanitizeText, sanitizeUrl } from '@/utils/menuTree.js'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
})
const emit = defineEmits(['update:modelValue', 'change'])

const localItems = ref(normalizeMenuTree(props.modelValue))
const emitDebounced = debounce(() => {
  const next = cloneTree(localItems.value)
  emit('update:modelValue', next)
  emit('change', next)
}, 250)

watch(() => props.modelValue, value => {
  localItems.value = normalizeMenuTree(value)
}, { deep: true })

function commit() {
  localItems.value = normalizeMenuTree(localItems.value)
  emitDebounced()
}

function addRootItem() {
  localItems.value.push(createMenuItem())
  commit()
}

function editItem(item, field, value) {
  item[field] = field === 'url' ? sanitizeUrl(value) : sanitizeText(value)
  if (field === 'target') item.target = value === '_blank' ? '_blank' : '_self'
  commit()
}

function addChild(item) {
  if (maxDepth([item]) >= MAX_MENU_DEPTH) return
  item.children.push(createMenuItem())
  commit()
}

function removeItem(item) {
  removeFrom(localItems.value, item.id)
  commit()
}

function moveItem(payload) {
  const { item, action } = payload
  if (action === 'up' || action === 'down') moveSibling(localItems.value, item.id, action)
  if (action === 'indent') indentItem(localItems.value, item.id)
  if (action === 'outdent') outdentItem(localItems.value, item.id)
  commit()
}

function handleBranchChange() {
  if (maxDepth(localItems.value) > MAX_MENU_DEPTH) {
    localItems.value = normalizeMenuTree(props.modelValue)
    return
  }
  commit()
}

function removeFrom(items, id) {
  const index = items.findIndex(item => item.id === id)
  if (index >= 0) return items.splice(index, 1)[0]
  for (const item of items) {
    const removed = removeFrom(item.children, id)
    if (removed) return removed
  }
  return null
}

function moveSibling(items, id, direction) {
  const index = items.findIndex(item => item.id === id)
  if (index >= 0) {
    const nextIndex = direction === 'up' ? index - 1 : index + 1
    if (nextIndex < 0 || nextIndex >= items.length) return true
    const [item] = items.splice(index, 1)
    items.splice(nextIndex, 0, item)
    return true
  }
  return items.some(item => moveSibling(item.children, id, direction))
}

function indentItem(items, id) {
  const index = items.findIndex(item => item.id === id)
  if (index > 0) {
    const previous = items[index - 1]
    if (maxDepth([previous]) >= MAX_MENU_DEPTH) return true
    const [item] = items.splice(index, 1)
    previous.children.push(item)
    return true
  }
  return items.some(item => indentItem(item.children, id))
}

function outdentItem(items, id, parent = null, parentList = null) {
  const index = items.findIndex(item => item.id === id)
  if (index >= 0 && parent && parentList) {
    const parentIndex = parentList.findIndex(item => item.id === parent.id)
    const [item] = items.splice(index, 1)
    parentList.splice(parentIndex + 1, 0, item)
    return true
  }
  return items.some(item => outdentItem(item.children, id, item, items))
}

const MenuBranch = defineComponent({
  name: 'MenuBranch',
  props: {
    items: { type: Array, required: true },
    depth: { type: Number, required: true },
  },
  emits: ['changed', 'edit', 'remove', 'add-child', 'move'],
  setup(branchProps, { emit }) {
    const listRef = ref(null)
    let sortable
    const canNest = computed(() => branchProps.depth < MAX_MENU_DEPTH)

    onMounted(() => initSortable())
    onBeforeUnmount(() => sortable?.destroy())

    watch(() => branchProps.items, async () => {
      await nextTick()
      initSortable()
    }, { deep: true })

    function initSortable() {
      sortable?.destroy()
      if (!listRef.value) return
      sortable = Sortable.create(listRef.value, {
        group: { name: 'cms-menu', pull: true, put: canNest.value },
        handle: '.menu-item__drag',
        animation: 140,
        fallbackOnBody: true,
        swapThreshold: 0.65,
        onEnd: event => {
          if (event.oldIndex == null || event.newIndex == null || event.oldIndex === event.newIndex) return emit('changed')
          const [moved] = branchProps.items.splice(event.oldIndex, 1)
          branchProps.items.splice(event.newIndex, 0, moved)
          emit('changed')
        },
        onAdd: event => {
          const source = event.from.__vueMenuItems
          if (!source) return emit('changed')
          const [moved] = source.splice(event.oldIndex, 1)
          branchProps.items.splice(event.newIndex, 0, moved)
          emit('changed')
        },
        onRemove: () => emit('changed'),
      })
      listRef.value.__vueMenuItems = branchProps.items
    }

    return () => h('ol', { ref: listRef, class: ['menu-branch', `depth-${branchProps.depth}`] }, branchProps.items.map(item =>
      h('li', { key: item.id, class: 'menu-item' }, [
        h('div', { class: 'menu-item__row' }, [
          h('button', { type: 'button', class: 'menu-item__drag', 'aria-label': `Drag ${item.title}` }, '::'),
          h('label', [h('span', 'Title'), h('input', { value: item.title, onInput: e => emit('edit', item, 'title', e.target.value) })]),
          h('label', [h('span', 'URL'), h('input', { value: item.url, onInput: e => emit('edit', item, 'url', e.target.value) })]),
          h('label', [h('span', 'Target'), h('select', { value: item.target, onChange: e => emit('edit', item, 'target', e.target.value) }, [
            h('option', { value: '_self' }, 'Same tab'),
            h('option', { value: '_blank' }, 'New tab'),
          ])]),
          h('div', { class: 'menu-item__actions' }, [
            h('button', { type: 'button', onClick: () => emit('move', { item, action: 'up' }) }, 'Up'),
            h('button', { type: 'button', onClick: () => emit('move', { item, action: 'down' }) }, 'Down'),
            h('button', { type: 'button', disabled: !canNest.value, onClick: () => emit('move', { item, action: 'indent' }) }, 'Indent'),
            h('button', { type: 'button', onClick: () => emit('move', { item, action: 'outdent' }) }, 'Outdent'),
            h('button', { type: 'button', disabled: !canNest.value, onClick: () => emit('add-child', item) }, 'Child'),
            h('button', { type: 'button', onClick: () => emit('remove', item) }, 'Delete'),
          ]),
        ]),
        h(MenuBranch, {
          items: item.children,
          depth: branchProps.depth + 1,
          onChanged: () => emit('changed'),
          onEdit: (...args) => emit('edit', ...args),
          onRemove: child => emit('remove', child),
          onAddChild: child => emit('add-child', child),
          onMove: payload => emit('move', payload),
        }),
      ])
    ))
  },
})
</script>

<style scoped>
.menu-builder{display:grid;gap:1rem}.menu-builder__toolbar{display:flex;align-items:center;justify-content:space-between;gap:1rem}.menu-builder__toolbar button{border-radius:.45rem;background:#1a365d;color:white;padding:.6rem .9rem;font-weight:700}.menu-builder__toolbar span,.menu-builder__hint{color:#64748b;font-size:.82rem}.menu-branch{display:grid;gap:.75rem;min-height:.75rem}.menu-branch.depth-2,.menu-branch.depth-3{margin-left:1.5rem;margin-top:.75rem;border-left:2px solid #e5e7eb;padding-left:1rem}.menu-item{list-style:none}.menu-item__row{display:grid;grid-template-columns:auto minmax(8rem,1fr) minmax(8rem,1fr) 8rem auto;gap:.65rem;align-items:end;border:1px solid #e5e7eb;border-radius:.55rem;background:#fff;padding:.75rem}.menu-item__drag{cursor:grab;border:1px solid #d1d5db;border-radius:.35rem;padding:.55rem;color:#64748b}.menu-item label{display:grid;gap:.25rem;color:#475569;font-size:.72rem;font-weight:800;text-transform:uppercase}.menu-item input,.menu-item select{border:1px solid #cbd5e1;border-radius:.4rem;padding:.55rem;font-size:.85rem;color:#111827}.menu-item__actions{display:flex;flex-wrap:wrap;gap:.35rem;justify-content:flex-end}.menu-item__actions button{border:1px solid #d1d5db;border-radius:.35rem;padding:.42rem .55rem;font-size:.75rem}.menu-item__actions button:disabled{opacity:.4;cursor:not-allowed}:global(.dark) .menu-item__row{background:#111827;border-color:#334155}:global(.dark) .menu-item input,:global(.dark) .menu-item select{background:#0f172a;color:#f8fafc;border-color:#475569}@media(max-width:900px){.menu-item__row{grid-template-columns:1fr}.menu-item__actions{justify-content:flex-start}}
</style>
