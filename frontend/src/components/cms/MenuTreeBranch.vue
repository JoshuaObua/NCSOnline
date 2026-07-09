<template>
  <ol ref="listRef" :class="['menu-branch', `depth-${depth}`]">
    <li v-for="item in items" :key="item.id" class="menu-item" :class="{ 'menu-item--hidden': item.hidden }">
      <div class="menu-item__row">
        <button type="button" class="menu-item__drag" :aria-label="`Drag ${item.title}`"><i class="fas fa-grip-vertical"></i></button>

        <label class="menu-item__field">
          <span>Title</span>
          <input :value="item.title" @input="emit('edit', item, 'title', $event.target.value)" />
        </label>

        <label class="menu-item__field">
          <span>Link type</span>
          <select :value="linkTypeFor(item)" @change="onLinkTypeChange(item, $event.target.value)">
            <option value="internal">Internal Page</option>
            <option value="custom">Custom URL</option>
            <option value="anchor">Anchor Link</option>
          </select>
        </label>

        <label v-if="linkTypeFor(item) === 'internal'" class="menu-item__field">
          <span>Page</span>
          <select :value="item.url" @change="emit('edit', item, 'url', $event.target.value)">
            <option value="/">Choose a page…</option>
            <optgroup v-for="group in pageGroups" :key="group.label" :label="group.label">
              <option v-for="page in group.options" :key="page.url" :value="page.url">{{ page.label }}</option>
            </optgroup>
          </select>
        </label>
        <label v-else-if="linkTypeFor(item) === 'anchor'" class="menu-item__field">
          <span>Anchor</span>
          <input :value="anchorFragment(item.url)" placeholder="section-id" @input="emit('edit', item, 'url', '#' + $event.target.value.replace(/^#+/, ''))" />
        </label>
        <label v-else class="menu-item__field">
          <span>URL</span>
          <input :value="item.url" @input="emit('edit', item, 'url', $event.target.value)" />
        </label>

        <label class="menu-item__field">
          <span>Target</span>
          <select :value="item.target" @change="emit('edit', item, 'target', $event.target.value)">
            <option value="_self">Same tab</option>
            <option value="_blank">New tab</option>
          </select>
        </label>

        <div class="menu-item__actions">
          <button
            type="button"
            class="menu-item__visibility"
            :class="{ 'is-hidden': item.hidden }"
            :aria-pressed="String(!item.hidden)"
            :title="item.hidden ? 'Hidden from the public menu — click to show' : 'Visible on the public menu — click to hide'"
            @click="emit('edit', item, 'hidden', !item.hidden)"
          >
            <i :class="item.hidden ? 'icofont-eye-blocked' : 'icofont-eye'" aria-hidden="true"></i>
          </button>
          <button type="button" @click="emit('move', { item, action: 'up' })">Up</button>
          <button type="button" @click="emit('move', { item, action: 'down' })">Down</button>
          <button type="button" :disabled="!canNest" @click="emit('move', { item, action: 'indent' })">Indent</button>
          <button type="button" @click="emit('move', { item, action: 'outdent' })">Outdent</button>
          <button type="button" :disabled="!canNest" @click="emit('add-child', item)">Child</button>
          <button type="button" class="menu-item__delete" @click="confirmDelete(item)">Delete</button>
        </div>
      </div>

      <span v-if="item.hidden" class="menu-item__badge">Hidden</span>

      <MenuTreeBranch
        :items="item.children"
        :depth="depth + 1"
        :pages="pages"
        @changed="emit('changed')"
        @edit="(...args) => emit('edit', ...args)"
        @remove="child => emit('remove', child)"
        @add-child="child => emit('add-child', child)"
        @move="payload => emit('move', payload)"
      />
    </li>
  </ol>
</template>

<script setup>
import Sortable from 'sortablejs'
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import Swal from 'sweetalert2'
import { MAX_MENU_DEPTH, inferLinkType } from '@/utils/menuTree.js'

const props = defineProps({
  items: { type: Array, required: true },
  depth: { type: Number, required: true },
  pages: { type: Array, default: () => [] },
})
const emit = defineEmits(['changed', 'edit', 'remove', 'add-child', 'move'])

const listRef = ref(null)
let sortable
const canNest = computed(() => props.depth < MAX_MENU_DEPTH)

const pageGroups = computed(() => {
  const groups = new Map()
  for (const page of props.pages) {
    const key = page.group || 'Pages'
    if (!groups.has(key)) groups.set(key, [])
    groups.get(key).push(page)
  }
  return [...groups.entries()].map(([label, options]) => ({ label, options }))
})

function linkTypeFor(item) {
  return inferLinkType(item.url, props.pages)
}

function anchorFragment(url) {
  return String(url || '').replace(/^#+/, '')
}

function onLinkTypeChange(item, type) {
  if (type === 'anchor') emit('edit', item, 'url', '#')
  else if (type === 'internal') emit('edit', item, 'url', props.pages[0]?.url || '/')
  else emit('edit', item, 'url', item.url.startsWith('#') ? '/' : item.url)
}

async function confirmDelete(item) {
  const label = item.children.length
    ? `"${item.title}" and its ${item.children.length} submenu item(s)`
    : `"${item.title}"`
  const result = await Swal.fire({
    title: `Delete ${label}?`,
    text: "This can't be undone.",
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'Delete',
    cancelButtonText: 'Cancel',
    confirmButtonColor: '#6777ef',
    cancelButtonColor: '#fc544b',
    reverseButtons: true,
  })
  if (result.isConfirmed) emit('remove', item)
}

onMounted(() => initSortable())
onBeforeUnmount(() => sortable?.destroy())

watch(() => props.items, async () => {
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
    ghostClass: 'menu-item--ghost',
    chosenClass: 'menu-item--chosen',
    dragClass: 'menu-item--dragging',
    onEnd: event => {
      if (event.oldIndex == null || event.newIndex == null || event.oldIndex === event.newIndex) return emit('changed')
      const [moved] = props.items.splice(event.oldIndex, 1)
      props.items.splice(event.newIndex, 0, moved)
      emit('changed')
    },
    onAdd: event => {
      const source = event.from.__vueMenuItems
      if (!source) return emit('changed')
      const [moved] = source.splice(event.oldIndex, 1)
      props.items.splice(event.newIndex, 0, moved)
      emit('changed')
    },
    onRemove: () => emit('changed'),
  })
  listRef.value.__vueMenuItems = props.items
}
</script>

<style scoped>
.menu-branch { display: grid; gap: 14px; min-height: 14px; padding-left: 0; margin: 0; }
.menu-branch.depth-2,
.menu-branch.depth-3 { margin-left: 22px; margin-top: 14px; border-left: 2px solid #e4e6fc; padding-left: 16px; }
.menu-item { list-style: none; }
.menu-item__row {
  display: grid;
  grid-template-columns: auto minmax(8rem, 1fr) minmax(7rem, auto) minmax(8rem, 1fr) 7rem auto;
  gap: 12px;
  align-items: end;
  border: 0;
  border-radius: 3px;
  background: #fff;
  box-shadow: 0 4px 25px rgba(0,0,0,.08);
  padding: 16px;
  transition: opacity 150ms ease;
}
.menu-item--hidden .menu-item__row { opacity: .55; }
.menu-item__drag { cursor: grab; border: 0; border-radius: 30px; background: #f4f6f9; padding: 9px 12px; color: #6777ef; box-shadow: none; }
.menu-item__field { display: grid; gap: 6px; color: #34395e; font-size: 12px; font-weight: 600; text-transform: none; min-width: 0; }
.menu-item__field input,
.menu-item__field select { border: 1px solid #e4e6fc; border-radius: 3px; background: #fdfdff; padding: 10px 15px; font-size: 13px; color: #495057; min-width: 0; outline: none; }
.menu-item__field input:focus,
.menu-item__field select:focus { border-color: #6777ef; box-shadow: 0 2px 6px #acb5f6; }
.menu-item__badge {
  display: inline-flex; align-items: center; gap: .3rem; margin-top: -.4rem;
  font-size: .68rem; font-weight: 800; text-transform: uppercase; letter-spacing: .04em;
  color: #ffa426; background: #fff4e6; border: 0; border-radius: 30px; padding: .2rem .6rem;
}
.menu-item__actions { display: flex; flex-wrap: wrap; gap: .35rem; justify-content: flex-end; align-items: center; }
.menu-item__actions button { border: 0; border-radius: 30px; background: #6777ef; color: #fff; padding: 8px 14px; font-size: 12px; font-weight: 600; box-shadow: 0 2px 6px #acb5f6; }
.menu-item__actions button:disabled { opacity: .4; cursor: not-allowed; }
.menu-item__visibility { display: inline-flex; align-items: center; justify-content: center; color: #fff; }
.menu-item__visibility.is-hidden { background: #ffa426; box-shadow: 0 2px 6px #ffc473; }
.menu-item__delete { background: #fc544b !important; box-shadow: 0 2px 6px #fd9b96 !important; }
.menu-item__delete:hover { filter: brightness(.96); }

/* Drag feedback */
.menu-item--ghost > .menu-item__row { outline: 2px dashed #6777ef; background: #f4f6f9; opacity: .7; }
.menu-item--chosen > .menu-item__row { box-shadow: 0 8px 20px rgba(15, 23, 42, .12); }
.menu-item--dragging > .menu-item__row { box-shadow: 0 12px 28px rgba(15, 23, 42, .22); transform: scale(1.01); }

:global(.dark .menu-item__row) { background: #1f2937; border-color: #334155; box-shadow: 0 4px 25px rgba(0,0,0,.28); }
:global(.dark .menu-item__field) { color: #f8fafc; }
:global(.dark .menu-item--hidden .menu-item__row) { opacity: .5; }
:global(.dark .menu-item__field input),
:global(.dark .menu-item__field select) { background: #111827; color: #f8fafc; border-color: #475569; }
:global(.dark .menu-item__drag),
:global(.dark .menu-item__actions button) { border-color: #475569; color: #fff; }
:global(.dark .menu-item__visibility.is-hidden) { background: #422006; border-color: #f5a623; color: #fdba74; }
:global(.dark .menu-item__badge) { background: #422006; border-color: #f5a623; color: #fdba74; }
:global(.dark .menu-branch.depth-2),
:global(.dark .menu-branch.depth-3) { border-color: #334155; }
:global(.dark .menu-item--ghost > .menu-item__row) { background: #422006; border-color: #f5a623; }

@media (max-width: 1100px) {
  .menu-item__row { grid-template-columns: 1fr 1fr; }
  .menu-item__actions { grid-column: 1 / -1; justify-content: flex-start; }
}
@media (max-width: 640px) {
  .menu-item__row { grid-template-columns: 1fr; }
}
</style>
