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
      :class="{ 'has-value': !!selectedCategory, 'is-loading': loading }"
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
        <i class="icofont-spinner icofont-spin"></i> Loading age categories...
      </span>
      <span v-else-if="selectedCategory" class="selected-content">
        <span class="cat-name">{{ selectedCategory.name || selectedCategory.code }}</span>
        <span v-if="selectedCategory.code" class="cat-code-badge">{{ selectedCategory.code }}</span>
        <span v-if="isAutoSuggested" class="cat-suggested-tag">Auto-matched</span>
      </span>
      <span v-else-if="modelValue" class="selected-content">
        <span class="cat-name">{{ modelValue }}</span>
        <span class="cat-code-badge custom">Custom</span>
      </span>
      <span v-else class="trigger-text placeholder-text">
        {{ placeholder || 'Search and select age category (e.g. U17, Senior, U20)...' }}
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

    <!-- Suggestion quick-pill banner when not selected -->
    <div v-if="suggestedCategory && modelValue !== suggestedCategory && !isOpen" class="suggestion-banner">
      <small>Suggested from DOB:</small>
      <button type="button" class="suggestion-chip" @click="selectCategoryByCode(suggestedCategory)">
        <i class="icofont-magic"></i> Set to {{ suggestedCategoryName || suggestedCategory }}
      </button>
    </div>

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
            placeholder="Type age or category (e.g. 17, U20, Senior, Masters)..."
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
        <span>{{ filteredCategories.length }} categor{{ filteredCategories.length === 1 ? 'y' : 'ies' }} available</span>
        <span v-if="searchQuery" class="filter-hint">Filtering by "{{ searchQuery }}"</span>
      </div>

      <!-- Scrollable Options List -->
      <ul ref="listRef" class="options-list" tabindex="-1">
        <li
          v-for="(cat, index) in filteredCategories"
          :key="cat.code || cat.id"
          class="option-item"
          :class="{
            'is-selected': cat.code === modelValue || cat.name === modelValue,
            'is-highlighted': index === highlightedIndex,
            'is-suggested': cat.code === suggestedCategory
          }"
          role="option"
          :aria-selected="cat.code === modelValue"
          @mouseenter="highlightedIndex = index"
          @click="selectOption(cat)"
        >
          <div class="option-body">
            <div class="option-title-row">
              <span class="option-title">{{ cat.name || cat.code }}</span>
              <span v-if="cat.code === suggestedCategory" class="suggested-badge">Matched DOB</span>
            </div>
            <span v-if="cat.description" class="option-meta">{{ cat.description }}</span>
            <span v-else-if="cat.min_age !== undefined && cat.max_age !== undefined" class="option-meta">
              Ages: {{ cat.min_age }} – {{ cat.max_age }} years
            </span>
          </div>
          <div class="option-badges">
            <span v-if="cat.code" class="code-tag">{{ cat.code }}</span>
            <i v-if="cat.code === modelValue || cat.name === modelValue" class="icofont-check selected-check"></i>
          </div>
        </li>

        <!-- Add Custom Category Option if query doesn't match -->
        <li
          v-if="searchQuery.trim() && !hasExactMatch"
          class="option-item add-custom-option"
          :class="{ 'is-highlighted': highlightedIndex === filteredCategories.length }"
          @mouseenter="highlightedIndex = filteredCategories.length"
          @click="addCustomCategory(searchQuery.trim())"
        >
          <div class="option-body">
            <span class="option-title">
              <i class="icofont-plus-circle"></i> Use custom category: <strong>"{{ searchQuery.trim() }}"</strong>
            </span>
            <span class="option-meta">Click to assign this custom age category to the athlete</span>
          </div>
          <span class="code-tag custom">Add</span>
        </li>

        <!-- Empty Results Message -->
        <li v-if="filteredCategories.length === 0 && !searchQuery.trim()" class="no-options">
          <i class="icofont-search-folder empty-icon"></i>
          <strong>No categories available</strong>
          <small>Type to create and add a custom age category.</small>
        </li>
      </ul>

      <!-- Footer Help -->
      <div class="panel-footer">
        <span><i class="icofont-info-circle"></i> Select standard category or type to add custom</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  id: {
    type: String,
    default: 'searchable-age-category',
  },
  categories: {
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
  placeholder: {
    type: String,
    default: '',
  },
  suggestedCategory: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue', 'change', 'add-category'])

const isOpen = ref(false)
const searchQuery = ref('')
const highlightedIndex = ref(0)
const containerRef = ref(null)
const searchInputRef = ref(null)
const listRef = ref(null)

// Standard fallback categories if DB list is empty
const defaultCategories = [
  { code: 'U10', name: 'U10 (Under 10)', min_age: 0, max_age: 10, description: 'Athletes aged 10 and under' },
  { code: 'U12', name: 'U12 (Under 12)', min_age: 11, max_age: 12, description: 'Athletes aged 11 to 12' },
  { code: 'U14', name: 'U14 (Under 14)', min_age: 13, max_age: 14, description: 'Athletes aged 13 to 14' },
  { code: 'U15', name: 'U15 (Under 15)', min_age: 13, max_age: 15, description: 'Athletes aged 13 to 15' },
  { code: 'U16', name: 'U16 (Under 16)', min_age: 15, max_age: 16, description: 'Athletes aged 15 to 16' },
  { code: 'U17', name: 'U17 (Under 17)', min_age: 16, max_age: 17, description: 'Athletes aged 16 to 17' },
  { code: 'U18', name: 'U18 (Under 18)', min_age: 17, max_age: 18, description: 'Athletes aged 17 to 18' },
  { code: 'U20', name: 'U20 (Under 20 / Junior)', min_age: 18, max_age: 20, description: 'Junior athletes aged 18 to 20' },
  { code: 'U23', name: 'U23 (Under 23 / Youth)', min_age: 21, max_age: 23, description: 'Youth / Under 23 category' },
  { code: 'Senior', name: 'Senior (Open Category)', min_age: 24, max_age: 34, description: 'Senior open competitive class' },
  { code: 'Masters', name: 'Masters / Veterans (35+)', min_age: 35, max_age: 99, description: 'Veteran / Masters athletics (35+)' },
  { code: 'Cadet', name: 'Cadet (12-14)', min_age: 12, max_age: 14, description: 'Cadet development tier' },
  { code: 'Junior', name: 'Junior (15-18)', min_age: 15, max_age: 18, description: 'Junior national category' },
  { code: 'Youth', name: 'Youth (16-19)', min_age: 16, max_age: 19, description: 'Youth development squad' },
  { code: 'Elite', name: 'Elite National Squad', min_age: 18, max_age: 40, description: 'High-performance national elite pool' },
  { code: 'Paralympic', name: 'Paralympic / Para-Athlete', min_age: 0, max_age: 99, description: 'Classified para-sports category' },
]

const allCategories = computed(() => {
  if (Array.isArray(props.categories) && props.categories.length > 0) {
    return props.categories
  }
  return defaultCategories
})

const selectedCategory = computed(() => {
  if (!props.modelValue) return null
  return allCategories.value.find(
    c => c.code === props.modelValue || c.name === props.modelValue || String(c.id) === String(props.modelValue)
  ) || null
})

const isAutoSuggested = computed(() => {
  return props.suggestedCategory && (props.modelValue === props.suggestedCategory)
})

const suggestedCategoryName = computed(() => {
  if (!props.suggestedCategory) return ''
  const match = allCategories.value.find(c => c.code === props.suggestedCategory)
  return match?.name || props.suggestedCategory
})

const filteredCategories = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) return allCategories.value

  return allCategories.value.filter(cat => {
    const code = String(cat.code || '').toLowerCase()
    const name = String(cat.name || '').toLowerCase()
    const desc = String(cat.description || '').toLowerCase()
    const minAge = String(cat.min_age ?? '')
    const maxAge = String(cat.max_age ?? '')

    return code.includes(query) ||
           name.includes(query) ||
           desc.includes(query) ||
           minAge === query ||
           maxAge === query ||
           `u${query}` === code
  })
})

const hasExactMatch = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) return false
  return allCategories.value.some(
    c => (c.code || '').toLowerCase() === query || (c.name || '').toLowerCase() === query
  )
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
  searchQuery.value = ''
  highlightedIndex.value = 0
  nextTick(() => {
    searchInputRef.value?.focus()
    scrollToSelected()
  })
}

function closeDropdown() {
  isOpen.value = false
  searchQuery.value = ''
}

function openAndFocusFirst() {
  if (!isOpen.value) openDropdown()
}

function selectOption(category) {
  const value = category.code || category.name
  emit('update:modelValue', value)
  emit('change', category)
  closeDropdown()
}

function selectCategoryByCode(code) {
  const match = allCategories.value.find(c => c.code === code)
  emit('update:modelValue', code)
  emit('change', match || { code, name: code })
}

function addCustomCategory(customName) {
  emit('update:modelValue', customName)
  emit('add-category', customName)
  emit('change', { code: customName, name: customName, is_custom: true })
  closeDropdown()
}

function clearSelection() {
  emit('update:modelValue', '')
  emit('change', null)
}

function navigateOptions(delta) {
  const max = filteredCategories.value.length + (hasExactMatch.value || !searchQuery.value.trim() ? 0 : 1)
  if (max === 0) return
  highlightedIndex.value = (highlightedIndex.value + delta + max) % max
  nextTick(scrollHighlightedIntoView)
}

function selectHighlighted() {
  if (highlightedIndex.value < filteredCategories.value.length) {
    const option = filteredCategories.value[highlightedIndex.value]
    if (option) selectOption(option)
  } else if (searchQuery.value.trim() && !hasExactMatch.value) {
    addCustomCategory(searchQuery.value.trim())
  }
}

function scrollHighlightedIntoView() {
  if (!listRef.value) return
  const highlightedEl = listRef.value.children[highlightedIndex.value]
  if (highlightedEl) {
    highlightedEl.scrollIntoView({ block: 'nearest' })
  }
}

function scrollToSelected() {
  if (!listRef.value || !selectedCategory.value) return
  const idx = filteredCategories.value.findIndex(
    c => c.code === selectedCategory.value.code
  )
  if (idx >= 0) {
    highlightedIndex.value = idx
    const el = listRef.value.children[idx]
    el?.scrollIntoView({ block: 'nearest' })
  }
}

function handleClickOutside(event) {
  if (containerRef.value && !containerRef.value.contains(event.target)) {
    closeDropdown()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.searchable-select-container {
  position: relative;
  width: 100%;
  box-sizing: border-box;
}

.proxy-input {
  position: absolute;
  opacity: 0;
  width: 1px;
  height: 1px;
  pointer-events: none;
  border: 0;
  margin: 0;
  padding: 0;
}

/* Trigger Button */
.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  min-height: 44px;
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  background: #ffffff;
  color: #1e293b;
  font-family: inherit;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
  box-sizing: border-box;
  text-align: left;
}

.select-trigger:hover:not(:disabled) {
  border-color: #94a3b8;
}

.select-trigger:focus-visible,
.searchable-select-container.is-open .select-trigger {
  outline: none;
  border-color: #6777ef;
  box-shadow: 0 0 0 3px rgba(103, 119, 239, 0.15);
}

.select-trigger:disabled {
  background: #f8fafc;
  color: #94a3b8;
  cursor: not-allowed;
  border-color: #e2e8f0;
}

.selected-content {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  flex: 1;
  overflow: hidden;
}

.cat-name {
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cat-code-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 7px;
  border-radius: 4px;
  background: #eff6ff;
  color: #2563eb;
  font-size: 11px;
  font-weight: 800;
  letter-spacing: 0.04em;
  border: 1px solid #bfdbfe;
  flex-shrink: 0;
}

.cat-code-badge.custom {
  background: #f5f3ff;
  color: #7c3aed;
  border-color: #ddd6fe;
}

.cat-suggested-tag {
  display: inline-flex;
  align-items: center;
  padding: 1px 6px;
  border-radius: 4px;
  background: #ecfdf5;
  color: #059669;
  font-size: 10px;
  font-weight: 700;
  border: 1px solid #a7f3d0;
}

.trigger-text.placeholder-text {
  color: #94a3b8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.trigger-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
  margin-left: 8px;
}

.clear-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: none;
  background: #f1f5f9;
  color: #64748b;
  cursor: pointer;
  font-size: 14px;
  padding: 0;
  transition: all 0.15s ease;
}

.clear-btn:hover {
  background: #fee2e2;
  color: #dc2626;
}

.chevron-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
  font-size: 14px;
  transition: transform 0.2s ease;
}

.chevron-icon.rotate {
  transform: rotate(180deg);
}

/* Suggestion Quick-Pill Banner */
.suggestion-banner {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 5px;
  font-size: 11px;
  color: #64748b;
}

.suggestion-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: 999px;
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.suggestion-chip:hover {
  background: #2563eb;
  color: #ffffff;
}

/* Dropdown Floating Panel */
.dropdown-panel {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  z-index: 1050;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.12);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* Search Header */
.search-header {
  padding: 10px;
  border-bottom: 1px solid #f1f5f9;
  background: #fafafa;
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
  pointer-events: none;
}

.search-input {
  width: 100%;
  height: 36px;
  padding: 0 32px 0 32px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-family: inherit;
  font-size: 12px;
  background: #ffffff;
  color: #1e293b;
  outline: none;
  box-sizing: border-box;
}

.search-input:focus {
  border-color: #6777ef;
  box-shadow: 0 0 0 2px rgba(103, 119, 239, 0.15);
}

.search-clear-btn {
  position: absolute;
  right: 8px;
  border: none;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  font-size: 14px;
  padding: 4px;
}

/* Panel Meta */
.panel-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 12px;
  background: #f8fafc;
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  border-bottom: 1px solid #f1f5f9;
}

.filter-hint {
  color: #6777ef;
  font-style: italic;
}

/* Options List */
.options-list {
  list-style: none;
  margin: 0;
  padding: 6px 0;
  max-height: 240px;
  overflow-y: auto;
  outline: none;
}

.option-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  cursor: pointer;
  transition: background 0.12s ease;
  border-left: 3px solid transparent;
}

.option-item:hover,
.option-item.is-highlighted {
  background: #f1f5f9;
}

.option-item.is-selected {
  background: #eff6ff;
  border-left-color: #2563eb;
}

.option-item.add-custom-option {
  border-top: 1px dashed #cbd5e1;
  background: #faf5ff;
}

.option-item.add-custom-option:hover,
.option-item.add-custom-option.is-highlighted {
  background: #f3e8ff;
}

.option-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  flex: 1;
}

.option-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.option-title {
  font-size: 12px;
  font-weight: 600;
  color: #1e293b;
}

.suggested-badge {
  font-size: 9px;
  font-weight: 800;
  padding: 1px 5px;
  border-radius: 4px;
  background: #ecfdf5;
  color: #059669;
  text-transform: uppercase;
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

.code-tag {
  font-size: 10px;
  font-weight: 800;
  padding: 2px 6px;
  border-radius: 4px;
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.code-tag.custom {
  background: #f5f3ff;
  color: #7c3aed;
  border-color: #ddd6fe;
}

.selected-check {
  color: #2563eb;
  font-size: 14px;
  font-weight: 800;
}

/* Empty State */
.no-options {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px 16px;
  color: #64748b;
  text-align: center;
  gap: 4px;
}

.empty-icon {
  font-size: 24px;
  color: #94a3b8;
  margin-bottom: 2px;
}

.no-options strong {
  font-size: 12px;
  color: #334155;
}

.no-options small {
  font-size: 11px;
  color: #94a3b8;
}

/* Panel Footer */
.panel-footer {
  padding: 6px 12px;
  background: #f8fafc;
  border-top: 1px solid #f1f5f9;
  font-size: 11px;
  color: #94a3b8;
}

/* Responsive */
@media (max-width: 640px) {
  .dropdown-panel {
    position: fixed;
    top: auto;
    bottom: 0;
    left: 0;
    right: 0;
    border-radius: 14px 14px 0 0;
    max-height: 70vh;
  }
}
</style>
