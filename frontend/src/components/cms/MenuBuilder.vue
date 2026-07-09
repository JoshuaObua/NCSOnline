<template>
  <div class="menu-builder">
    <div class="menu-builder__toolbar">
      <button type="button" @click="addRootItem">Add root item</button>
      <span>Max depth: {{ MAX_MENU_DEPTH }}</span>
    </div>
    <MenuTreeBranch
      :items="localItems"
      :depth="1"
      :pages="pages"
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
import { ref, watch } from 'vue'
import MenuTreeBranch from '@/components/cms/MenuTreeBranch.vue'
import { MAX_MENU_DEPTH, cloneTree, createMenuItem, debounce, maxDepth, normalizeMenuTree, sanitizeText, sanitizeUrl } from '@/utils/menuTree.js'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  pages: { type: Array, default: () => [] },
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
  if (field === 'url') item.url = sanitizeUrl(value)
  else if (field === 'target') item.target = value === '_blank' ? '_blank' : '_self'
  else if (field === 'hidden') item.hidden = value === true
  else item[field] = sanitizeText(value)
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
</script>

<style scoped>
.menu-builder { display: grid; gap: 1rem; }
.menu-builder__toolbar { display: flex; align-items: center; justify-content: space-between; gap: 1rem; }
.menu-builder__toolbar button { border-radius: .45rem; background: #1a365d; color: white; padding: .6rem .9rem; font-weight: 700; }
.menu-builder__toolbar span,
.menu-builder__hint { color: #64748b; font-size: .82rem; }
</style>
