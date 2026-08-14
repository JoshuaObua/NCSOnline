<template>
  <div class="relative w-full" ref="containerRef">
    <!-- Trigger Button / Display Box -->
    <div
      @click="toggleDropdown"
      :class="[
        'w-full flex items-center justify-between py-2 px-3 text-xs rounded-lg border cursor-pointer transition-colors duration-150 select-none',
        isOpen ? 'border-slate-400 dark:border-slate-500 bg-white dark:bg-slate-800' : 'border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800',
        disabled ? 'opacity-60 cursor-not-allowed bg-slate-100 dark:bg-slate-900' : '',
        customClass
      ]"
    >
      <span v-if="selectedLabel" class="text-slate-900 dark:text-slate-100 font-medium truncate">
        {{ selectedLabel }}
      </span>
      <span v-else class="text-slate-400 dark:text-slate-500 truncate">
        {{ placeholder }}
      </span>

      <div class="flex items-center gap-1.5 ml-2 flex-shrink-0 text-slate-400">
        <button
          v-if="modelValue && !disabled && clearable"
          type="button"
          @click.stop="clearSelection"
          class="hover:text-slate-600 dark:hover:text-slate-200 text-xs p-0.5 rounded"
          title="Clear selection"
        >
          <i class="icofont-close-line"></i>
        </button>
        <i class="icofont-rounded-down text-xs transition-transform duration-200" :class="{ 'rotate-180': isOpen }"></i>
      </div>
    </div>

    <!-- Hidden Native Input for HTML5 form validation if required -->
    <input
      v-if="required"
      type="text"
      :value="modelValue"
      required
      class="sr-only opacity-0 absolute pointer-events-none w-px h-px"
      tabindex="-1"
    />

    <!-- Dropdown Menu with Search Field -->
    <div
      v-if="isOpen"
      class="absolute left-0 right-0 top-full mt-1 z-50 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl shadow-xl overflow-hidden animate-fadeIn"
    >
      <!-- Search Input Header -->
      <div class="p-2 border-b border-slate-100 dark:border-slate-700 bg-slate-50 dark:bg-slate-900/50">
        <div class="relative flex items-center">
          <i class="icofont-search-2 absolute left-2.5 text-xs text-slate-400 pointer-events-none"></i>
          <input
            ref="searchInputRef"
            v-model="searchQuery"
            type="text"
            :placeholder="searchPlaceholder"
            class="w-full pl-7 pr-3 py-1.5 text-xs bg-white dark:bg-slate-800 text-slate-900 dark:text-white border border-slate-200 dark:border-slate-700 rounded-md outline-none focus:border-slate-400 dark:focus:border-slate-500"
            @keydown.escape="closeDropdown"
            @keydown.down.prevent="highlightNext"
            @keydown.up.prevent="highlightPrev"
            @keydown.enter.prevent="selectHighlighted"
          />
        </div>
      </div>

      <!-- Options List -->
      <div class="max-h-56 overflow-y-auto p-1 divide-y divide-transparent">
        <div
          v-for="(opt, idx) in filteredOptions"
          :key="opt.value"
          @click="selectOption(opt)"
          @mouseenter="highlightedIndex = idx"
          :class="[
            'px-3 py-2 text-xs rounded-lg cursor-pointer flex items-center justify-between transition-colors duration-100',
            highlightedIndex === idx ? 'bg-slate-100 dark:bg-slate-700/60 text-slate-900 dark:text-white' : 'text-slate-700 dark:text-slate-300',
            opt.value === modelValue ? 'font-bold text-blue-600 dark:text-blue-400 bg-blue-50/60 dark:bg-blue-950/40' : ''
          ]"
        >
          <div class="truncate">
            <div>{{ opt.label }}</div>
            <div v-if="opt.subtitle" class="text-[10px] text-slate-400 font-normal truncate mt-0.5">{{ opt.subtitle }}</div>
          </div>
          <i v-if="opt.value === modelValue" class="icofont-check text-blue-600 dark:text-blue-400 text-xs ml-2"></i>
        </div>

        <div v-if="!filteredOptions.length" class="px-3 py-4 text-center text-xs text-slate-400">
          No matches found for "{{ searchQuery }}"
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  options: {
    type: Array,
    required: true
    // Expects strings or objects: { label: string, value: string, subtitle?: string }
  },
  placeholder: {
    type: String,
    default: 'Select an option...'
  },
  searchPlaceholder: {
    type: String,
    default: 'Search items...'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  required: {
    type: Boolean,
    default: false
  },
  clearable: {
    type: Boolean,
    default: false
  },
  customClass: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const containerRef = ref(null)
const searchInputRef = ref(null)
const isOpen = ref(false)
const searchQuery = ref('')
const highlightedIndex = ref(-1)

// Normalize options to { label, value, subtitle }
const normalizedOptions = computed(() => {
  return props.options.map(item => {
    if (typeof item === 'string') {
      return { label: item, value: item }
    }
    return {
      label: item.label || String(item.value),
      value: item.value,
      subtitle: item.subtitle || ''
    }
  })
})

const filteredOptions = computed(() => {
  if (!searchQuery.value.trim()) return normalizedOptions.value
  const q = searchQuery.value.toLowerCase().trim()
  return normalizedOptions.value.filter(opt => {
    return opt.label.toLowerCase().includes(q) || (opt.subtitle && opt.subtitle.toLowerCase().includes(q))
  })
})

const selectedLabel = computed(() => {
  const match = normalizedOptions.value.find(opt => opt.value === props.modelValue)
  return match ? match.label : ''
})

function toggleDropdown() {
  if (props.disabled) return
  if (isOpen.value) {
    closeDropdown()
  } else {
    openDropdown()
  }
}

function openDropdown() {
  isOpen.value = true
  searchQuery.value = ''
  highlightedIndex.value = -1
  nextTick(() => {
    searchInputRef.value?.focus()
  })
}

function closeDropdown() {
  isOpen.value = false
  searchQuery.value = ''
  highlightedIndex.value = -1
}

function selectOption(opt) {
  emit('update:modelValue', opt.value)
  emit('change', opt.value)
  closeDropdown()
}

function clearSelection() {
  emit('update:modelValue', '')
  emit('change', '')
}

function highlightNext() {
  if (highlightedIndex.value < filteredOptions.value.length - 1) {
    highlightedIndex.value++
  }
}

function highlightPrev() {
  if (highlightedIndex.value > 0) {
    highlightedIndex.value--
  }
}

function selectHighlighted() {
  if (highlightedIndex.value >= 0 && highlightedIndex.value < filteredOptions.value.length) {
    selectOption(filteredOptions.value[highlightedIndex.value])
  }
}

function handleClickOutside(e) {
  if (containerRef.value && !containerRef.value.contains(e.target)) {
    closeDropdown()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fadeIn {
  animation: fadeIn 0.15s ease-out;
}
</style>
