<template>
  <div ref="containerRef" class="searchable-competition-select" :class="{ 'is-open': isOpen, 'is-disabled': disabled }">
    <!-- Hidden Proxy Input for Native Validation -->
    <input
      :id="id"
      type="text"
      tabindex="-1"
      class="proxy-input"
      :value="modelValue"
      :required="required"
      aria-hidden="true"
    />

    <!-- Trigger Button -->
    <button
      type="button"
      class="select-trigger"
      :class="{ 'has-value': !!modelValue, 'is-loading': loading || isFetching }"
      :disabled="disabled || loading || isFetching"
      :aria-expanded="isOpen"
      aria-haspopup="listbox"
      @click="toggleDropdown"
      @keydown.down.prevent="openAndFocusFirst"
      @keydown.space.prevent="toggleDropdown"
      @keydown.enter.prevent="toggleDropdown"
      @keydown.esc.prevent="closeDropdown"
    >
      <span v-if="loading || isFetching" class="trigger-text placeholder-text">
        <i class="icofont-spinner icofont-spin"></i> Loading competitions...
      </span>
      <span v-else-if="selectedCompObj" class="selected-comp-content">
        <div class="comp-icon-badge">
          <i class="icofont-trophy"></i>
        </div>
        <div class="comp-meta-summary">
          <span class="comp-name-text">{{ selectedCompObj.name }}</span>
          <span v-if="selectedCompObj.venue" class="comp-venue-pill">
            <i class="icofont-location-pin"></i> {{ selectedCompObj.venue }}
          </span>
        </div>
      </span>
      <span v-else-if="modelValue" class="selected-comp-content">
        <div class="comp-icon-badge">
          <i class="icofont-trophy"></i>
        </div>
        <div class="comp-meta-summary">
          <span class="comp-name-text">{{ modelValue }}</span>
        </div>
      </span>
      <span v-else class="trigger-text placeholder-text">
        {{ placeholder || 'Search competition by Name, Venue, or Level...' }}
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

    <!-- Floating Dropdown Panel -->
    <div v-if="isOpen" class="dropdown-panel" role="listbox" @click.stop>
      <!-- Search Header -->
      <div class="search-header">
        <div class="search-input-wrapper">
          <i class="icofont-search search-icon"></i>
          <input
            ref="searchInputRef"
            v-model="searchQuery"
            type="text"
            class="search-input"
            placeholder="Type Competition Name, Venue, or Level..."
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

      <!-- Quick Custom Entry Option -->
      <div v-if="searchQuery.trim()" class="custom-entry-bar" @click="selectCustom(searchQuery.trim())">
        <i class="icofont-plus-circle text-primary"></i>
        <span>Use custom competition name: <strong>"{{ searchQuery.trim() }}"</strong></span>
      </div>

      <!-- Panel Meta -->
      <div class="panel-meta">
        <span>{{ filteredCompetitions.length }} competition record{{ filteredCompetitions.length === 1 ? '' : 's' }}</span>
        <span v-if="searchQuery" class="filter-hint">Matching "{{ searchQuery }}"</span>
      </div>

      <!-- Options List -->
      <ul ref="listRef" class="options-list" tabindex="-1">
        <li
          v-for="(comp, index) in filteredCompetitions"
          :key="comp.id || index"
          class="comp-option-item"
          :class="{
            'is-selected': isCompSelected(comp),
            'is-highlighted': index === highlightedIndex
          }"
          role="option"
          :aria-selected="isCompSelected(comp)"
          @click="selectCompetition(comp)"
          @mouseenter="highlightedIndex = index"
        >
          <div class="option-icon">
            <i class="icofont-trophy"></i>
          </div>

          <div class="option-details">
            <div class="option-name-row">
              <span class="comp-full-name">{{ comp.name }}</span>
              <span v-if="comp.level" class="level-badge">
                {{ comp.level }}
              </span>
            </div>
            <div class="option-sub-row">
              <span v-if="comp.venue" class="info-item"><i class="icofont-location-pin"></i> {{ comp.venue }}</span>
              <span v-if="comp.host_country" class="info-item"><i class="icofont-globe"></i> {{ comp.host_country }}</span>
              <span v-if="comp.starts_on" class="info-item"><i class="icofont-calendar"></i> {{ comp.starts_on }}</span>
            </div>
          </div>

          <div v-if="isCompSelected(comp)" class="check-icon">
            <i class="icofont-check-circled"></i>
          </div>
        </li>

        <!-- Empty State -->
        <li v-if="filteredCompetitions.length === 0" class="empty-state">
          <i class="icofont-trophy"></i>
          <p>No competitions found matching "{{ searchQuery }}"</p>
          <button type="button" class="btn-custom-fallback" @click="selectCustom(searchQuery.trim())">
            Use "{{ searchQuery.trim() }}"
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { listNsmisDomain } from '@/api/nsmis.js'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  id: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  competitions: {
    type: Array,
    default: null
  },
  required: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'change', 'select'])

const containerRef = ref(null)
const searchInputRef = ref(null)
const listRef = ref(null)

const isOpen = ref(false)
const searchQuery = ref('')
const highlightedIndex = ref(-1)
const isFetching = ref(false)
const fetchedCompetitions = ref([])

const allCompetitions = computed(() => {
  if (props.competitions && Array.isArray(props.competitions) && props.competitions.length > 0) {
    return props.competitions
  }
  return fetchedCompetitions.value
})

const filteredCompetitions = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return allCompetitions.value

  return allCompetitions.value.filter(c => {
    const name = (c.name || '').toLowerCase()
    const venue = (c.venue || '').toLowerCase()
    const country = (c.host_country || '').toLowerCase()
    const level = (c.level || '').toLowerCase()
    return name.includes(q) || venue.includes(q) || country.includes(q) || level.includes(q)
  })
})

const selectedCompObj = computed(() => {
  if (!props.modelValue) return null
  const val = String(props.modelValue).trim().toLowerCase()
  return allCompetitions.value.find(c => {
    return (c.id && String(c.id).toLowerCase() === val) ||
      (c.name && c.name.toLowerCase() === val)
  })
})

function isCompSelected(c) {
  if (!props.modelValue) return false
  const val = String(props.modelValue).trim()
  if (c.id && String(c.id) === val) return true
  if (c.name === val) return true
  return false
}

async function loadCompetitionsIfNeeded() {
  if (props.competitions && props.competitions.length > 0) return
  if (fetchedCompetitions.value.length > 0) return

  isFetching.value = true
  try {
    const res = await listNsmisDomain('competitions', { per_page: 300 })
    const data = res?.data?.data?.items || res?.data?.items || res?.data?.data || res?.data || []
    fetchedCompetitions.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.warn('[SearchableCompetitionSelect] Failed to load competitions:', err)
  } finally {
    isFetching.value = false
  }
}

function toggleDropdown() {
  if (props.disabled || props.loading || isFetching.value) return
  if (isOpen.value) {
    closeDropdown()
  } else {
    openDropdown()
  }
}

function openDropdown() {
  isOpen.value = true
  highlightedIndex.value = -1
  searchQuery.value = ''
  loadCompetitionsIfNeeded()
  nextTick(() => {
    searchInputRef.value?.focus()
  })
}

function openAndFocusFirst() {
  openDropdown()
  if (filteredCompetitions.value.length > 0) {
    highlightedIndex.value = 0
  }
}

function closeDropdown() {
  isOpen.value = false
  searchQuery.value = ''
  highlightedIndex.value = -1
}

function selectCompetition(comp) {
  const val = comp.id || comp.name
  emit('update:modelValue', val)
  emit('change', val)
  emit('select', comp)
  closeDropdown()
}

function selectCustom(customVal) {
  if (!customVal) return
  emit('update:modelValue', customVal)
  emit('change', customVal)
  emit('select', { name: customVal })
  closeDropdown()
}

function selectHighlighted() {
  if (highlightedIndex.value >= 0 && highlightedIndex.value < filteredCompetitions.value.length) {
    selectCompetition(filteredCompetitions.value[highlightedIndex.value])
  } else if (searchQuery.value.trim()) {
    selectCustom(searchQuery.value.trim())
  }
}

function clearSelection() {
  emit('update:modelValue', '')
  emit('change', '')
  emit('select', null)
}

function navigateOptions(direction) {
  const max = filteredCompetitions.value.length - 1
  if (max < 0) return
  let next = highlightedIndex.value + direction
  if (next < 0) next = max
  if (next > max) next = 0
  highlightedIndex.value = next

  nextTick(() => {
    const listEl = listRef.value
    if (!listEl) return
    const activeEl = listEl.children[highlightedIndex.value]
    if (activeEl) {
      activeEl.scrollIntoView({ block: 'nearest' })
    }
  })
}

function handleClickOutside(event) {
  if (containerRef.value && !containerRef.value.contains(event.target)) {
    closeDropdown()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  if (!props.competitions || props.competitions.length === 0) {
    loadCompetitionsIfNeeded()
  }
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.searchable-competition-select {
  position: relative;
  width: 100%;
  font-family: inherit;
}

.proxy-input {
  position: absolute;
  opacity: 0;
  pointer-events: none;
  width: 0;
  height: 0;
  margin: 0;
  border: 0;
  padding: 0;
}

.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  min-height: 44px;
  padding: 6px 14px;
  background-color: #ffffff;
  border: 1px solid #d9dce5;
  border-radius: 8px;
  color: #34395e;
  font-size: 14px;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.select-trigger:hover:not(:disabled) {
  border-color: #6777ef;
  box-shadow: 0 2px 6px rgba(103, 119, 239, 0.15);
}

.searchable-competition-select.is-open .select-trigger {
  border-color: #6777ef;
  box-shadow: 0 0 0 3px rgba(103, 119, 239, 0.18);
  outline: none;
}

.select-trigger:disabled {
  background-color: #f8fafc;
  color: #94a3b8;
  cursor: not-allowed;
  border-color: #e2e8f0;
}

.selected-comp-content {
  display: flex;
  align-items: center;
  gap: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.comp-icon-badge {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0;
}

.comp-meta-summary {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.comp-name-text {
  font-weight: 600;
  color: #34395e;
  font-size: 14px;
}

.comp-venue-pill {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 6px;
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fde68a;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.placeholder-text {
  color: #94a3b8;
  font-size: 13.5px;
}

.trigger-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-left: 8px;
  flex-shrink: 0;
}

.clear-btn {
  background: transparent;
  border: none;
  color: #94a3b8;
  font-size: 16px;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 4px;
  line-height: 1;
}

.clear-btn:hover {
  color: #ef4444;
}

.chevron-icon {
  color: #6c757d;
  font-size: 14px;
  transition: transform 0.2s ease;
}

.chevron-icon.rotate {
  transform: rotate(180deg);
}

.dropdown-panel {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  z-index: 1050;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.12), 0 8px 10px -6px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  animation: slideDown 0.15s ease-out;
}

@keyframes slideDown {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}

.search-header {
  padding: 10px 12px;
  border-bottom: 1px solid #f1f5f9;
  background: #fafbfc;
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 10px;
  color: #94a3b8;
  font-size: 15px;
}

.search-input {
  width: 100%;
  height: 36px;
  padding: 6px 32px 6px 32px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 13px;
  color: #34395e;
  background: #ffffff;
  outline: none;
}

.search-input:focus {
  border-color: #6777ef;
  box-shadow: 0 0 0 2px rgba(103, 119, 239, 0.15);
}

.search-clear-btn {
  position: absolute;
  right: 8px;
  background: transparent;
  border: none;
  color: #94a3b8;
  font-size: 15px;
  cursor: pointer;
}

.custom-entry-bar {
  padding: 8px 14px;
  background: #eff6ff;
  border-bottom: 1px solid #dbeafe;
  font-size: 12.5px;
  color: #1e40af;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
}

.custom-entry-bar:hover {
  background: #dbeafe;
}

.panel-meta {
  display: flex;
  justify-content: space-between;
  padding: 6px 14px;
  font-size: 11px;
  color: #94a3b8;
  background: #f8fafc;
  border-bottom: 1px solid #f1f5f9;
  font-weight: 500;
}

.options-list {
  list-style: none;
  margin: 0;
  padding: 4px;
  max-height: 260px;
  overflow-y: auto;
}

.comp-option-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.12s;
  gap: 12px;
}

.comp-option-item:hover,
.comp-option-item.is-highlighted {
  background-color: #f1f5f9;
}

.comp-option-item.is-selected {
  background-color: #eff6ff;
}

.option-icon {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0;
}

.option-details {
  flex: 1;
  min-width: 0;
}

.option-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.comp-full-name {
  font-weight: 600;
  font-size: 13.5px;
  color: #34395e;
}

.level-badge {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 6px;
  border-radius: 4px;
  background: #fef3c7;
  color: #92400e;
}

.option-sub-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 2px;
  font-size: 11.5px;
  color: #6c757d;
  flex-wrap: wrap;
}

.info-item {
  display: inline-flex;
  align-items: center;
  gap: 3px;
}

.check-icon {
  color: #6777ef;
  font-size: 18px;
  margin-left: 8px;
}

.empty-state {
  padding: 24px 16px;
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
}

.empty-state i {
  font-size: 28px;
  margin-bottom: 6px;
  display: block;
}

.btn-custom-fallback {
  margin-top: 8px;
  padding: 5px 12px;
  border: 1px solid #6777ef;
  border-radius: 6px;
  background: transparent;
  color: #6777ef;
  font-weight: 600;
  cursor: pointer;
}

/* Dark Theme Support */
:global(html.dark) .searchable-competition-select .select-trigger,
:global(body.dark) .searchable-competition-select .select-trigger,
:global([data-theme="dark"]) .searchable-competition-select .select-trigger {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
}

:global(html.dark) .searchable-competition-select .comp-name-text,
:global(body.dark) .searchable-competition-select .comp-name-text,
:global([data-theme="dark"]) .searchable-competition-select .comp-name-text,
:global(html.dark) .searchable-competition-select .comp-full-name,
:global(body.dark) .searchable-competition-select .comp-full-name,
:global([data-theme="dark"]) .searchable-competition-select .comp-full-name {
  color: #f8fafc !important;
}

:global(html.dark) .searchable-competition-select .dropdown-panel,
:global(body.dark) .searchable-competition-select .dropdown-panel,
:global([data-theme="dark"]) .searchable-competition-select .dropdown-panel {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(html.dark) .searchable-competition-select .search-header,
:global(body.dark) .searchable-competition-select .search-header,
:global([data-theme="dark"]) .searchable-competition-select .search-header,
:global(html.dark) .searchable-competition-select .panel-meta,
:global(body.dark) .searchable-competition-select .panel-meta,
:global([data-theme="dark"]) .searchable-competition-select .panel-meta {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #94a3b8 !important;
}

:global(html.dark) .searchable-competition-select .search-input,
:global(body.dark) .searchable-competition-select .search-input,
:global([data-theme="dark"]) .searchable-competition-select .search-input {
  background-color: #1e293b !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
}

:global(html.dark) .searchable-competition-select .comp-option-item:hover,
:global(body.dark) .searchable-competition-select .comp-option-item:hover,
:global([data-theme="dark"]) .searchable-competition-select .comp-option-item:hover,
:global(html.dark) .searchable-competition-select .comp-option-item.is-highlighted,
:global(body.dark) .searchable-competition-select .comp-option-item.is-highlighted,
:global([data-theme="dark"]) .searchable-competition-select .comp-option-item.is-highlighted {
  background-color: #334155 !important;
}
</style>
