<template>
  <div ref="containerRef" class="searchable-user-select" :class="{ 'is-open': isOpen, 'is-disabled': disabled }">
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

    <!-- Selected Box / Trigger Button -->
    <button
      type="button"
      class="select-trigger"
      :class="{ 'has-value': !!modelValue, 'is-loading': loading }"
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
        <i class="icofont-spinner icofont-spin"></i> Loading registered user profiles...
      </span>
      <span v-else-if="modelValue" class="selected-user-content">
        <div class="user-avatar-badge">
          <img v-if="selectedUserObj?.avatar_url" :src="selectedUserObj.avatar_url" alt="" />
          <span v-else>{{ userInitials(modelValue) }}</span>
        </div>
        <div class="user-meta-summary">
          <span class="user-name-text">{{ modelValue }}</span>
          <span v-if="selectedUserObj?.nin" class="user-nin-pill">
            <i class="icofont-id-card"></i> NIN: {{ selectedUserObj.nin }}
          </span>
          <span v-else-if="selectedUserObj?.email" class="user-email-pill">
            {{ selectedUserObj.email }}
          </span>
        </div>
      </span>
      <span v-else class="trigger-text placeholder-text">
        {{ placeholder || 'Search user by Name, NIN, or Email...' }}
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
            placeholder="Type Name, NIN (e.g. CM9201...), or Email..."
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

      <!-- Quick Custom Entry Option if user types something not in DB -->
      <div v-if="searchQuery.trim()" class="custom-entry-bar" @click="selectCustom(searchQuery.trim())">
        <i class="icofont-plus-circle text-primary"></i>
        <span>Use custom name: <strong>"{{ searchQuery.trim() }}"</strong></span>
      </div>

      <!-- Panel Meta -->
      <div class="panel-meta">
        <span>{{ filteredUsers.length }} registered profile{{ filteredUsers.length === 1 ? '' : 's' }}</span>
        <span v-if="searchQuery" class="filter-hint">Matching "{{ searchQuery }}"</span>
      </div>

      <!-- Scrollable User Options List -->
      <ul ref="listRef" class="options-list" tabindex="-1">
        <li
          v-for="(u, index) in filteredUsers"
          :key="u.id"
          class="user-option-item"
          :class="{
            'is-selected': isUserSelected(u),
            'is-highlighted': index === highlightedIndex
          }"
          role="option"
          :aria-selected="isUserSelected(u)"
          @click="selectUser(u)"
          @mouseenter="highlightedIndex = index"
        >
          <div class="option-avatar">
            <img v-if="u.avatar_url" :src="u.avatar_url" alt="" />
            <span v-else>{{ userInitials(getUserFullName(u)) }}</span>
          </div>

          <div class="option-details">
            <div class="option-name-row">
              <span class="user-full-name">{{ getUserFullName(u) }}</span>
              <span v-if="u.nin" class="nin-badge">
                <i class="icofont-id-card"></i> NIN: {{ u.nin }}
              </span>
            </div>
            <div class="option-sub-row">
              <span v-if="u.email" class="user-email-text"><i class="icofont-email"></i> {{ u.email }}</span>
              <span v-if="u.phone" class="user-phone-text"><i class="icofont-phone"></i> {{ u.phone }}</span>
              <span v-if="u.account_status" class="status-pill">{{ u.account_status }}</span>
            </div>
          </div>

          <div v-if="isUserSelected(u)" class="check-icon">
            <i class="icofont-check-circled"></i>
          </div>
        </li>

        <!-- Empty State -->
        <li v-if="filteredUsers.length === 0" class="empty-state">
          <i class="icofont-search-user"></i>
          <p>No registered user profiles found matching "{{ searchQuery }}"</p>
          <button type="button" class="btn btn-sm btn-outline-primary mt-1" @click="selectCustom(searchQuery.trim())">
            Assign as "{{ searchQuery.trim() }}"
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import * as cms from '@/api/cms.js'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  id: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: '',
  },
  users: {
    type: Array,
    default: null,
  },
  required: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'change', 'select'])

const containerRef = ref(null)
const searchInputRef = ref(null)
const listRef = ref(null)

const isOpen = ref(false)
const searchQuery = ref('')
const highlightedIndex = ref(-1)
const loading = ref(false)
const internalUsers = ref([])

function getUserFullName(u) {
  if (!u) return ''
  const name = [u.first_name, u.last_name].filter(Boolean).join(' ').trim()
  return name || u.full_name || u.email || 'Unnamed User'
}

function userInitials(name) {
  if (!name) return 'U'
  const parts = String(name).trim().split(/\s+/)
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

const allUsers = computed(() => {
  if (props.users && Array.isArray(props.users) && props.users.length > 0) {
    return props.users
  }
  return internalUsers.value
})

const filteredUsers = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return allUsers.value

  return allUsers.value.filter(u => {
    const fullName = getUserFullName(u).toLowerCase()
    const nin = (u.nin || '').toLowerCase()
    const email = (u.email || '').toLowerCase()
    const phone = (u.phone || '').toLowerCase()
    return fullName.includes(q) || nin.includes(q) || email.includes(q) || phone.includes(q)
  })
})

const selectedUserObj = computed(() => {
  if (!props.modelValue) return null
  const val = String(props.modelValue).trim().toLowerCase()
  return allUsers.value.find(u => {
    const fn = getUserFullName(u).toLowerCase()
    return fn === val || (u.id && u.id.toLowerCase() === val)
  })
})

function isUserSelected(u) {
  if (!props.modelValue) return false
  const val = String(props.modelValue).trim().toLowerCase()
  const fn = getUserFullName(u).toLowerCase()
  return fn === val || u.id === props.modelValue
}

async function loadUsersFromApi() {
  if (props.users && props.users.length > 0) return
  loading.value = true
  try {
    const res = await cms.adminListUsers({ page: 1, per_page: 200 })
    const list = res?.data?.data || res?.data?.items || res?.data?.users || res?.data || []
    internalUsers.value = Array.isArray(list) ? list : []
  } catch (err) {
    console.warn('[SearchableUserSelect] Could not fetch users from API:', err)
  } finally {
    loading.value = false
  }
}

function toggleDropdown() {
  if (props.disabled || loading.value) return
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
  if (internalUsers.value.length === 0 && (!props.users || props.users.length === 0)) {
    loadUsersFromApi()
  }
  nextTick(() => {
    searchInputRef.value?.focus()
  })
}

function openAndFocusFirst() {
  openDropdown()
  if (filteredUsers.value.length > 0) {
    highlightedIndex.value = 0
  }
}

function closeDropdown() {
  isOpen.value = false
  searchQuery.value = ''
  highlightedIndex.value = -1
}

function selectUser(user) {
  const fullName = getUserFullName(user)
  emit('update:modelValue', fullName)
  emit('change', fullName)
  emit('select', user)
  closeDropdown()
}

function selectCustom(customName) {
  if (!customName) return
  emit('update:modelValue', customName)
  emit('change', customName)
  emit('select', { full_name: customName, first_name: customName })
  closeDropdown()
}

function selectHighlighted() {
  if (highlightedIndex.value >= 0 && highlightedIndex.value < filteredUsers.value.length) {
    selectUser(filteredUsers.value[highlightedIndex.value])
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
  const max = filteredUsers.value.length - 1
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
  loadUsersFromApi()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.searchable-user-select {
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

/* Trigger Button */
.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  min-height: 44px;
  padding: 6px 14px;
  background-color: #ffffff;
  border: 1px solid #dbe2ea;
  border-radius: 8px;
  color: #2c3e50;
  font-size: 14px;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.select-trigger:hover:not(:disabled) {
  border-color: #3b82f6;
  box-shadow: 0 2px 6px rgba(59, 130, 246, 0.12);
}

.searchable-user-select.is-open .select-trigger {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.18);
  outline: none;
}

.select-trigger:disabled {
  background-color: #f8fafc;
  color: #94a3b8;
  cursor: not-allowed;
  border-color: #e2e8f0;
}

/* User Content inside trigger */
.selected-user-content {
  display: flex;
  align-items: center;
  gap: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.user-avatar-badge {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4f46e5, #06b6d4);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
  overflow: hidden;
}

.user-avatar-badge img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-meta-summary {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.user-name-text {
  font-weight: 600;
  color: #1e293b;
  font-size: 14px;
}

.user-nin-pill {
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

.user-email-pill {
  font-size: 11px;
  color: #64748b;
}

.placeholder-text {
  color: #94a3b8;
  font-size: 13.5px;
}

/* Actions in trigger */
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
  transition: color 0.15s;
}

.clear-btn:hover {
  color: #ef4444;
}

.chevron-icon {
  color: #64748b;
  font-size: 14px;
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
  z-index: 1050;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.05);
  overflow: hidden;
  animation: slideDown 0.15s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Search Header */
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
  color: #1e293b;
  background: #ffffff;
  outline: none;
  transition: border-color 0.15s;
}

.search-input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.12);
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
  transition: background 0.15s;
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

/* Options List */
.options-list {
  list-style: none;
  margin: 0;
  padding: 4px;
  max-height: 250px;
  overflow-y: auto;
}

.user-option-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.12s;
  gap: 12px;
}

.user-option-item:hover,
.user-option-item.is-highlighted {
  background-color: #f1f5f9;
}

.user-option-item.is-selected {
  background-color: #eff6ff;
}

.option-avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6366f1, #3b82f6);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
  overflow: hidden;
}

.option-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
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

.user-full-name {
  font-weight: 600;
  font-size: 13.5px;
  color: #1e293b;
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

.option-sub-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 2px;
  font-size: 11.5px;
  color: #64748b;
  flex-wrap: wrap;
}

.status-pill {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  padding: 1px 5px;
  border-radius: 3px;
  background: #f1f5f9;
  color: #475569;
}

.check-icon {
  color: #3b82f6;
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

/* Dark Theme Support */
:global(body.dark-theme) .searchable-user-select .select-trigger {
  background-color: #1c2237;
  border-color: #2e3650;
  color: #e2e8f0;
}

:global(body.dark-theme) .searchable-user-select .select-trigger:hover:not(:disabled) {
  border-color: #6366f1;
}

:global(body.dark-theme) .searchable-user-select .user-name-text {
  color: #f1f5f9;
}

:global(body.dark-theme) .searchable-user-select .dropdown-panel {
  background-color: #1c2237;
  border-color: #2e3650;
}

:global(body.dark-theme) .searchable-user-select .search-header,
:global(body.dark-theme) .searchable-user-select .panel-meta {
  background-color: #161b2e;
  border-color: #2e3650;
  color: #94a3b8;
}

:global(body.dark-theme) .searchable-user-select .search-input {
  background-color: #1f2742;
  border-color: #384266;
  color: #f8fafc;
}

:global(body.dark-theme) .searchable-user-select .custom-entry-bar {
  background-color: #1e293b;
  border-color: #334155;
  color: #93c5fd;
}

:global(body.dark-theme) .searchable-user-select .user-option-item:hover,
:global(body.dark-theme) .searchable-user-select .user-option-item.is-highlighted {
  background-color: #242c4c;
}

:global(body.dark-theme) .searchable-user-select .user-option-item.is-selected {
  background-color: #2d3759;
}

:global(body.dark-theme) .searchable-user-select .user-full-name {
  color: #f8fafc;
}

:global(body.dark-theme) .searchable-user-select .nin-badge {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
}
</style>
