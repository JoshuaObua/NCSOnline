<template>
  <div ref="containerRef" class="searchable-athlete-select" :class="{ 'is-open': isOpen, 'is-disabled': disabled }">
    <!-- Hidden Proxy Input for Native Form Validation -->
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
        <i class="icofont-spinner icofont-spin"></i> Loading athlete registry...
      </span>
      <span v-else-if="selectedAthleteObj" class="selected-athlete-content">
        <div class="athlete-avatar-badge">
          <span>{{ getInitials(getAthleteName(selectedAthleteObj)) }}</span>
        </div>
        <div class="athlete-meta-summary">
          <span class="athlete-name-text">{{ getAthleteName(selectedAthleteObj) }}</span>
          <span v-if="getAthleteNin(selectedAthleteObj)" class="athlete-nin-pill" title="National ID Number (NIN)">
            <i class="icofont-id-card"></i> NIN: {{ getAthleteNin(selectedAthleteObj) }}
          </span>
          <span v-else-if="getAthleteEmail(selectedAthleteObj)" class="athlete-email-pill">
            <i class="icofont-email"></i> {{ getAthleteEmail(selectedAthleteObj) }}
          </span>
          <span v-else-if="selectedAthleteObj.athlete_number" class="athlete-num-pill">
            #{{ selectedAthleteObj.athlete_number }}
          </span>
        </div>
      </span>
      <span v-else-if="modelValue" class="selected-athlete-content">
        <div class="athlete-avatar-badge">
          <span>{{ getInitials(modelValue) }}</span>
        </div>
        <div class="athlete-meta-summary">
          <span class="athlete-name-text">{{ modelValue }}</span>
        </div>
      </span>
      <span v-else class="trigger-text placeholder-text">
        {{ placeholder || 'Search athlete by Name, NIN, Email, or Athlete Number...' }}
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
      <!-- Search Bar -->
      <div class="search-header">
        <div class="search-input-wrapper">
          <i class="icofont-search search-icon"></i>
          <input
            ref="searchInputRef"
            v-model="searchQuery"
            type="text"
            class="search-input"
            placeholder="Type Name, NIN (e.g. CM9201...), Email, or Athlete #..."
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

      <!-- Quick Custom Value Entry -->
      <div v-if="searchQuery.trim()" class="custom-entry-bar" @click="selectCustom(searchQuery.trim())">
        <i class="icofont-plus-circle text-primary"></i>
        <span>Select as custom athlete name: <strong>"{{ searchQuery.trim() }}"</strong></span>
      </div>

      <!-- Panel Meta -->
      <div class="panel-meta">
        <span>{{ filteredAthletes.length }} athlete record{{ filteredAthletes.length === 1 ? '' : 's' }}</span>
        <span v-if="searchQuery" class="filter-hint">Matching "{{ searchQuery }}"</span>
      </div>

      <!-- Options List -->
      <ul ref="listRef" class="options-list" tabindex="-1">
        <li
          v-for="(ath, index) in filteredAthletes"
          :key="ath.id || index"
          class="athlete-option-item"
          :class="{
            'is-selected': isAthleteSelected(ath),
            'is-highlighted': index === highlightedIndex
          }"
          role="option"
          :aria-selected="isAthleteSelected(ath)"
          @click="selectAthlete(ath)"
          @mouseenter="highlightedIndex = index"
        >
          <div class="option-avatar">
            <span>{{ getInitials(getAthleteName(ath)) }}</span>
          </div>

          <div class="option-details">
            <div class="option-name-row">
              <span class="athlete-full-name">{{ getAthleteName(ath) }}</span>
              <span v-if="getAthleteNin(ath)" class="nin-badge">
                <i class="icofont-id-card"></i> NIN: {{ getAthleteNin(ath) }}
              </span>
              <span v-if="ath.athlete_number" class="num-badge">
                #{{ ath.athlete_number }}
              </span>
            </div>
            <div class="option-sub-row">
              <span v-if="getAthleteEmail(ath)" class="info-item"><i class="icofont-email"></i> {{ getAthleteEmail(ath) }}</span>
              <span v-if="getAthletePhone(ath)" class="info-item"><i class="icofont-phone"></i> {{ getAthletePhone(ath) }}</span>
              <span v-if="ath.discipline" class="info-tag">{{ ath.discipline }}</span>
              <span v-if="ath.club" class="info-tag club-tag">{{ ath.club }}</span>
            </div>
          </div>

          <div v-if="isAthleteSelected(ath)" class="check-icon">
            <i class="icofont-check-circled"></i>
          </div>
        </li>

        <!-- Empty State -->
        <li v-if="filteredAthletes.length === 0" class="empty-state">
          <i class="icofont-runner-alt-1"></i>
          <p>No registered athletes found matching "{{ searchQuery }}"</p>
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
  athletes: {
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
const fetchedAthletes = ref([])

function getAthleteName(ath) {
  if (!ath) return ''
  return ath.full_name || ath.name || [ath.first_name, ath.last_name].filter(Boolean).join(' ') || 'Unnamed Athlete'
}

function getAthleteNin(ath) {
  if (!ath) return ''
  return ath.national_id_passport || ath.nin || ath.national_id || ''
}

function getAthleteEmail(ath) {
  if (!ath) return ''
  return ath.email_address || ath.email || ''
}

function getAthletePhone(ath) {
  if (!ath) return ''
  return ath.phone_contact || ath.phone || ''
}

function getInitials(name) {
  if (!name) return 'A'
  const parts = String(name).trim().split(/\s+/)
  if (parts.length >= 2) return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
  return name.substring(0, 2).toUpperCase()
}

const allAthletes = computed(() => {
  if (props.athletes && Array.isArray(props.athletes) && props.athletes.length > 0) {
    return props.athletes
  }
  return fetchedAthletes.value
})

const filteredAthletes = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return allAthletes.value

  return allAthletes.value.filter(ath => {
    const name = getAthleteName(ath).toLowerCase()
    const nin = getAthleteNin(ath).toLowerCase()
    const email = getAthleteEmail(ath).toLowerCase()
    const phone = getAthletePhone(ath).toLowerCase()
    const num = String(ath.athlete_number || '').toLowerCase()
    const club = String(ath.club || '').toLowerCase()
    const disc = String(ath.discipline || '').toLowerCase()
    return name.includes(q) || nin.includes(q) || email.includes(q) || phone.includes(q) || num.includes(q) || club.includes(q) || disc.includes(q)
  })
})

const selectedAthleteObj = computed(() => {
  if (!props.modelValue) return null
  const val = String(props.modelValue).trim().toLowerCase()
  return allAthletes.value.find(ath => {
    return (ath.id && String(ath.id).toLowerCase() === val) ||
      getAthleteName(ath).toLowerCase() === val ||
      (ath.athlete_number && String(ath.athlete_number).toLowerCase() === val)
  })
})

function isAthleteSelected(ath) {
  if (!props.modelValue) return false
  const val = String(props.modelValue).trim()
  if (ath.id && String(ath.id) === val) return true
  if (getAthleteName(ath) === val) return true
  return false
}

async function loadAthletesIfNeeded() {
  if (props.athletes && props.athletes.length > 0) return
  if (fetchedAthletes.value.length > 0) return

  isFetching.value = true
  try {
    const res = await listNsmisDomain('athletes', { per_page: 300 })
    const data = res?.data?.data?.items || res?.data?.items || res?.data?.data || res?.data || []
    fetchedAthletes.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.warn('[SearchableAthleteSelect] Failed to load athletes:', err)
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
  loadAthletesIfNeeded()
  nextTick(() => {
    searchInputRef.value?.focus()
  })
}

function openAndFocusFirst() {
  openDropdown()
  if (filteredAthletes.value.length > 0) {
    highlightedIndex.value = 0
  }
}

function closeDropdown() {
  isOpen.value = false
  searchQuery.value = ''
  highlightedIndex.value = -1
}

function selectAthlete(ath) {
  const val = ath.id || getAthleteName(ath)
  emit('update:modelValue', val)
  emit('change', val)
  emit('select', ath)
  closeDropdown()
}

function selectCustom(customVal) {
  if (!customVal) return
  emit('update:modelValue', customVal)
  emit('change', customVal)
  emit('select', { full_name: customVal, name: customVal })
  closeDropdown()
}

function selectHighlighted() {
  if (highlightedIndex.value >= 0 && highlightedIndex.value < filteredAthletes.value.length) {
    selectAthlete(filteredAthletes.value[highlightedIndex.value])
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
  const max = filteredAthletes.value.length - 1
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
  if (!props.athletes || props.athletes.length === 0) {
    loadAthletesIfNeeded()
  }
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.searchable-athlete-select {
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

.searchable-athlete-select.is-open .select-trigger {
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

.selected-athlete-content {
  display: flex;
  align-items: center;
  gap: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.athlete-avatar-badge {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6777ef, #3b82f6);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
}

.athlete-meta-summary {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.athlete-name-text {
  font-weight: 600;
  color: #34395e;
  font-size: 14px;
}

.athlete-nin-pill {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 6px;
  background: rgba(16, 185, 129, 0.12);
  color: #059669;
  border: 1px solid rgba(16, 185, 129, 0.25);
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.athlete-email-pill, .athlete-num-pill {
  font-size: 11px;
  color: #6c757d;
  display: inline-flex;
  align-items: center;
  gap: 3px;
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

.athlete-option-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.12s;
  gap: 12px;
}

.athlete-option-item:hover,
.athlete-option-item.is-highlighted {
  background-color: #f1f5f9;
}

.athlete-option-item.is-selected {
  background-color: #eff6ff;
}

.option-avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6777ef, #3b82f6);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
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

.athlete-full-name {
  font-weight: 600;
  font-size: 13.5px;
  color: #34395e;
}

.nin-badge {
  font-size: 11px;
  font-weight: 600;
  padding: 1px 7px;
  border-radius: 4px;
  background: #dcfce7;
  color: #166534;
  display: inline-flex;
  align-items: center;
  gap: 3px;
}

.num-badge {
  font-size: 11px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 4px;
  background: #e2e8f0;
  color: #475569;
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

.info-tag {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 3px;
  background: #f1f5f9;
  color: #475569;
}

.club-tag {
  background: #e0f2fe;
  color: #0369a1;
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
:global(html.dark) .searchable-athlete-select .select-trigger,
:global(body.dark) .searchable-athlete-select .select-trigger,
:global([data-theme="dark"]) .searchable-athlete-select .select-trigger {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
}

:global(html.dark) .searchable-athlete-select .athlete-name-text,
:global(body.dark) .searchable-athlete-select .athlete-name-text,
:global([data-theme="dark"]) .searchable-athlete-select .athlete-name-text,
:global(html.dark) .searchable-athlete-select .athlete-full-name,
:global(body.dark) .searchable-athlete-select .athlete-full-name,
:global([data-theme="dark"]) .searchable-athlete-select .athlete-full-name {
  color: #f8fafc !important;
}

:global(html.dark) .searchable-athlete-select .dropdown-panel,
:global(body.dark) .searchable-athlete-select .dropdown-panel,
:global([data-theme="dark"]) .searchable-athlete-select .dropdown-panel {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(html.dark) .searchable-athlete-select .search-header,
:global(body.dark) .searchable-athlete-select .search-header,
:global([data-theme="dark"]) .searchable-athlete-select .search-header,
:global(html.dark) .searchable-athlete-select .panel-meta,
:global(body.dark) .searchable-athlete-select .panel-meta,
:global([data-theme="dark"]) .searchable-athlete-select .panel-meta {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #94a3b8 !important;
}

:global(html.dark) .searchable-athlete-select .search-input,
:global(body.dark) .searchable-athlete-select .search-input,
:global([data-theme="dark"]) .searchable-athlete-select .search-input {
  background-color: #1e293b !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
}

:global(html.dark) .searchable-athlete-select .athlete-option-item:hover,
:global(body.dark) .searchable-athlete-select .athlete-option-item:hover,
:global([data-theme="dark"]) .searchable-athlete-select .athlete-option-item:hover,
:global(html.dark) .searchable-athlete-select .athlete-option-item.is-highlighted,
:global(body.dark) .searchable-athlete-select .athlete-option-item.is-highlighted,
:global([data-theme="dark"]) .searchable-athlete-select .athlete-option-item.is-highlighted {
  background-color: #334155 !important;
}
</style>
