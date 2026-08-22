<template>
  <div ref="containerRef" class="searchable-select-container" :class="{ 'is-open': isOpen, 'is-disabled': disabled }">
    <!-- Hidden input for native form validation -->
    <input
      :id="id"
      type="text"
      tabindex="-1"
      class="proxy-input"
      :value="modelValue"
      :required="required"
      aria-hidden="true"
    />

    <!-- Selected Box / Trigger Button -->
    <button
      type="button"
      class="select-trigger"
      :class="{ 'has-value': !!selectedFederation, 'is-loading': loading }"
      :disabled="disabled || loading"
      :aria-expanded="isOpen"
      aria-haspopup="listbox"
      @click="toggleDropdown"
      @keydown.down.prevent="openAndFocusFirst"
      @keydown.space.prevent="toggleDropdown"
      @keydown.enter.prevent="toggleDropdown"
      @keydown.esc.prevent="closeDropdown"
    >
      <span v-if="loading" class="trigger-text placeholder-text">
        <i class="icofont-spinner icofont-spin"></i> Loading federations from database...
      </span>
      <span v-else-if="selectedFederation" class="selected-content">
        <span class="fed-name">{{ selectedFederation.name }}</span>
        <span v-if="selectedFederation.acronym" class="fed-acronym-badge">{{ selectedFederation.acronym }}</span>
      </span>
      <span v-else class="trigger-text placeholder-text">
        {{ placeholder || 'Search and select a sports federation...' }}
      </span>

      <div class="trigger-actions">
        <button
          v-if="modelValue && !disabled"
          type="button"
          class="clear-btn"
          title="Clear selection"
          aria-label="Clear selection"
          @click.stop="clearSelection"
        >
          <i class="icofont-close-line"></i>
        </button>
        <span class="chevron-icon" :class="{ 'rotate': isOpen }">
          <i class="icofont-rounded-down"></i>
        </span>
      </div>
    </button>

    <!-- Floating Search & Dropdown Menu -->
    <div v-if="isOpen" class="dropdown-panel" role="listbox" @click.stop>
      <!-- Search Input Header -->
      <div class="search-header">
        <div class="search-input-wrapper">
          <i class="icofont-search search-icon"></i>
          <input
            ref="searchInputRef"
            v-model="searchQuery"
            type="text"
            class="search-input"
            placeholder="Type name or acronym (e.g. FUFA, UAF, Boxing)..."
            autocomplete="off"
            @keydown.down.prevent="navigateOptions(1)"
            @keydown.up.prevent="navigateOptions(-1)"
            @keydown.enter.prevent="selectHighlighted"
            @keydown.esc.prevent="closeDropdown"
          />
          <button
            v-if="searchQuery"
            type="button"
            class="search-clear-btn"
            title="Clear search"
            @click="searchQuery = ''"
          >
            <i class="icofont-close-line"></i>
          </button>
        </div>
      </div>

      <!-- Quick Info Bar -->
      <div class="panel-meta">
        <span>{{ filteredFederations.length }} federation{{ filteredFederations.length === 1 ? '' : 's' }} found</span>
        <span v-if="searchQuery" class="filter-hint">Filtering by "{{ searchQuery }}"</span>
      </div>

      <!-- Scrollable Options List -->
      <ul ref="listRef" class="options-list" tabindex="-1">
        <li
          v-for="(fed, index) in filteredFederations"
          :key="fed.id"
          class="option-item"
          :class="{
            'is-selected': fed.id === modelValue,
            'is-highlighted': index === highlightedIndex
          }"
          role="option"
          :aria-selected="fed.id === modelValue"
          @mouseenter="highlightedIndex = index"
          @click="selectOption(fed)"
        >
          <div class="option-body">
            <span class="option-title">{{ fed.name }}</span>
            <span v-if="fed.ncs_registration_number" class="option-meta">Reg: {{ fed.ncs_registration_number }}</span>
          </div>
          <div class="option-badges">
            <span v-if="fed.acronym" class="acronym-tag">{{ fed.acronym }}</span>
            <i v-if="fed.id === modelValue" class="icofont-check selected-check"></i>
          </div>
        </li>

        <!-- Empty Results -->
        <li v-if="filteredFederations.length === 0" class="no-options">
          <i class="icofont-search-folder empty-icon"></i>
          <strong>No federations found</strong>
          <span>No sports federation matches "{{ searchQuery }}".</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  federations: {
    type: Array,
    default: () => [],
  },
  required: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  id: {
    type: String,
    default: 'federation-searchable-select',
  },
  placeholder: {
    type: String,
    default: 'Search and select a sports federation...',
  },
})

const emit = defineEmits(['update:modelValue', 'change'])

const containerRef = ref(null)
const searchInputRef = ref(null)
const listRef = ref(null)
const isOpen = ref(false)
const searchQuery = ref('')
const highlightedIndex = ref(0)

const selectedFederation = computed(() => {
  if (!props.modelValue) return null
  return (props.federations || []).find(f => f.id === props.modelValue) || null
})

const filteredFederations = computed(() => {
  const query = (searchQuery.value || '').trim().toLowerCase()
  const list = props.federations || []
  if (!query) return list

  return list.filter(fed => {
    const name = (fed.name || '').toLowerCase()
    const acronym = (fed.acronym || '').toLowerCase()
    const regNum = (fed.ncs_registration_number || '').toLowerCase()
    const id = (fed.id || '').toLowerCase()
    return name.includes(query) || acronym.includes(query) || regNum.includes(query) || id.includes(query)
  })
})

watch(filteredFederations, () => {
  highlightedIndex.value = 0
})

function toggleDropdown() {
  if (props.disabled || props.loading) return
  if (isOpen.value) {
    closeDropdown()
  } else {
    openDropdown()
  }
}

function openDropdown() {
  isOpen.value = true
  highlightedIndex.value = 0
  nextTick(() => {
    if (searchInputRef.value) {
      searchInputRef.value.focus()
    }
  })
}

function openAndFocusFirst() {
  if (!isOpen.value) {
    openDropdown()
  } else {
    navigateOptions(1)
  }
}

function closeDropdown() {
  isOpen.value = false
  searchQuery.value = ''
}

function selectOption(fed) {
  emit('update:modelValue', fed.id)
  emit('change', fed.id)
  closeDropdown()
}

function clearSelection() {
  emit('update:modelValue', '')
  emit('change', '')
  searchQuery.value = ''
}

function navigateOptions(delta) {
  if (!filteredFederations.value.length) return
  const len = filteredFederations.value.length
  let next = highlightedIndex.value + delta
  if (next < 0) next = len - 1
  if (next >= len) next = 0
  highlightedIndex.value = next
  scrollHighlightedIntoView()
}

function selectHighlighted() {
  if (filteredFederations.value.length > 0 && highlightedIndex.value >= 0 && highlightedIndex.value < filteredFederations.value.length) {
    selectOption(filteredFederations.value[highlightedIndex.value])
  }
}

function scrollHighlightedIntoView() {
  nextTick(() => {
    if (!listRef.value) return
    const items = listRef.value.querySelectorAll('.option-item')
    if (items[highlightedIndex.value]) {
      items[highlightedIndex.value].scrollIntoView({ block: 'nearest' })
    }
  })
}

function handleClickOutside(e) {
  if (containerRef.value && !containerRef.value.contains(e.target)) {
    closeDropdown()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside, true)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside, true)
})
</script>

<style scoped>
.searchable-select-container {
  position: relative;
  width: 100%;
}

.proxy-input {
  position: absolute;
  opacity: 0;
  width: 1px;
  height: 1px;
  top: 0;
  left: 0;
  pointer-events: none;
}

.select-trigger {
  width: 100%;
  min-height: 46px;
  padding: 8px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  background: #ffffff;
  border: 1px solid #d9dce5;
  border-radius: 8px;
  color: #1e293b;
  font-family: inherit;
  font-size: 14px;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

.select-trigger:hover:not(:disabled) {
  border-color: #6777ef;
  box-shadow: 0 0 0 1px rgba(103, 119, 239, 0.2);
}

.searchable-select-container.is-open .select-trigger {
  border-color: #6777ef;
  box-shadow: 0 0 0 3px rgba(103, 119, 239, 0.18);
}

.select-trigger:disabled {
  background: #f1f5f9;
  color: #94a3b8;
  cursor: not-allowed;
  border-color: #e2e8f0;
}

.placeholder-text {
  color: #94a3b8;
  font-size: 14px;
}

.selected-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.fed-name {
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.fed-acronym-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 7px;
  border-radius: 4px;
  background: #eff6ff;
  color: #2563eb;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.04em;
  border: 1px solid #bfdbfe;
  flex-shrink: 0;
}

.trigger-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.clear-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border: none;
  background: #f1f5f9;
  color: #64748b;
  border-radius: 50%;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.15s ease;
}

.clear-btn:hover {
  background: #e2e8f0;
  color: #ef4444;
}

.chevron-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
  font-size: 16px;
  transition: transform 0.2s ease;
}

.chevron-icon.rotate {
  transform: rotate(180deg);
}

/* Floating Dropdown Panel */
.dropdown-panel {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  z-index: 1000;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  box-shadow: 0 14px 35px rgba(0, 0, 0, 0.15), 0 4px 10px rgba(0, 0, 0, 0.05);
  overflow: hidden;
  animation: dropdownSlideIn 0.15s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes dropdownSlideIn {
  from {
    opacity: 0;
    transform: translateY(-6px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.search-header {
  padding: 10px 12px;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 10px;
  color: #64748b;
  font-size: 15px;
  pointer-events: none;
}

.search-input {
  width: 100%;
  height: 38px;
  padding: 6px 32px 6px 34px;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-family: inherit;
  font-size: 13px;
  color: #1e293b;
  transition: all 0.15s ease;
}

.search-input:focus {
  outline: none;
  border-color: #6777ef;
  box-shadow: 0 0 0 2px rgba(103, 119, 239, 0.2);
}

.search-clear-btn {
  position: absolute;
  right: 8px;
  border: none;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  padding: 4px;
  font-size: 14px;
}

.search-clear-btn:hover {
  color: #475569;
}

.panel-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 14px;
  background: #f1f5f9;
  border-bottom: 1px solid #e2e8f0;
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
}

.filter-hint {
  color: #2563eb;
}

.options-list {
  list-style: none;
  margin: 0;
  padding: 6px 0;
  max-height: 260px;
  overflow-y: auto;
}

.option-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 9px 14px;
  cursor: pointer;
  transition: background 0.15s ease;
}

.option-item:hover,
.option-item.is-highlighted {
  background: #f1f5f9;
}

.option-item.is-selected {
  background: #eff6ff;
}

.option-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.option-title {
  font-size: 13px;
  font-weight: 600;
  color: #1e293b;
  line-height: 1.3;
}

.option-meta {
  font-size: 11px;
  color: #64748b;
}

.option-badges {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.acronym-tag {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 4px;
  background: #f1f5f9;
  color: #475569;
  font-size: 11px;
  font-weight: 700;
  border: 1px solid #e2e8f0;
}

.option-item.is-selected .acronym-tag {
  background: #dbeafe;
  color: #1d4ed8;
  border-color: #bfdbfe;
}

.selected-check {
  color: #2563eb;
  font-size: 16px;
  font-weight: bold;
}

.no-options {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 28px 16px;
  text-align: center;
  color: #64748b;
}

.empty-icon {
  font-size: 28px;
  color: #94a3b8;
  margin-bottom: 8px;
}

.no-options strong {
  font-size: 14px;
  color: #334155;
  margin-bottom: 4px;
}

.no-options span {
  font-size: 12px;
}

/* =====================================================================
   DARK THEME STYLES
   ===================================================================== */
:global(html.dark) .select-trigger,
:global(body.dark) .select-trigger,
:global([data-theme="dark"]) .select-trigger {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
}

:global(html.dark) .fed-name,
:global(body.dark) .fed-name,
:global([data-theme="dark"]) .fed-name {
  color: #f8fafc !important;
}

:global(html.dark) .fed-acronym-badge,
:global(body.dark) .fed-acronym-badge,
:global([data-theme="dark"]) .fed-acronym-badge {
  background-color: #1e3a8a !important;
  color: #93c5fd !important;
  border-color: #3b82f6 !important;
}

:global(html.dark) .clear-btn,
:global(body.dark) .clear-btn,
:global([data-theme="dark"]) .clear-btn {
  background-color: #334155 !important;
  color: #cbd5e1 !important;
}

:global(html.dark) .dropdown-panel,
:global(body.dark) .dropdown-panel,
:global([data-theme="dark"]) .dropdown-panel {
  background-color: #1e293b !important;
  border-color: #334155 !important;
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.5) !important;
}

:global(html.dark) .search-header,
:global(body.dark) .search-header,
:global([data-theme="dark"]) .search-header {
  background-color: #0f172a !important;
  border-bottom-color: #334155 !important;
}

:global(html.dark) .search-input,
:global(body.dark) .search-input,
:global([data-theme="dark"]) .search-input {
  background-color: #1e293b !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
}

:global(html.dark) .panel-meta,
:global(body.dark) .panel-meta,
:global([data-theme="dark"]) .panel-meta {
  background-color: #0f172a !important;
  border-bottom-color: #334155 !important;
  color: #94a3b8 !important;
}

:global(html.dark) .option-title,
:global(body.dark) .option-title,
:global([data-theme="dark"]) .option-title {
  color: #f8fafc !important;
}

:global(html.dark) .option-item:hover,
:global(html.dark) .option-item.is-highlighted,
:global(body.dark) .option-item:hover,
:global(body.dark) .option-item.is-highlighted,
:global([data-theme="dark"]) .option-item:hover,
:global([data-theme="dark"]) .option-item.is-highlighted {
  background-color: #334155 !important;
}

:global(html.dark) .option-item.is-selected,
:global(body.dark) .option-item.is-selected,
:global([data-theme="dark"]) .option-item.is-selected {
  background-color: #1e3a8a !important;
}

:global(html.dark) .acronym-tag,
:global(body.dark) .acronym-tag,
:global([data-theme="dark"]) .acronym-tag {
  background-color: #0f172a !important;
  color: #94a3b8 !important;
  border-color: #334155 !important;
}

:global(html.dark) .no-options strong,
:global(body.dark) .no-options strong,
:global([data-theme="dark"]) .no-options strong {
  color: #f8fafc !important;
}
</style>
