<template>
  <div class="admin-card overflow-hidden">

    <!-- Card header -->
    <div class="admin-card-header flex-wrap gap-2">
      <div class="flex items-center gap-3 min-w-0">
        <h2 v-if="title" class="text-sm font-semibold text-gray-800">{{ title }}</h2>
        <span v-if="total !== null" class="text-xs text-gray-400 font-medium">{{ total }} total</span>
      </div>
      <div class="flex items-center gap-2 flex-wrap">
        <!-- Search slot or default search input -->
        <slot name="search">
          <div v-if="searchable" class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="w-3.5 h-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
              </svg>
            </div>
            <input
              :value="searchValue"
              @input="$emit('update:searchValue', $event.target.value)"
              type="text"
              :placeholder="searchPlaceholder || 'Search…'"
              class="pl-8 pr-3 py-1.5 text-xs border border-gray-200 rounded-lg bg-gray-50 focus:outline-none focus:ring-1 focus:ring-primary-700 w-48"
            />
          </div>
        </slot>
        <!-- Extra filters slot -->
        <slot name="filters" />
        <!-- Create button -->
        <button
          v-if="createLabel"
          @click="$emit('create')"
          class="btn-primary flex items-center gap-1.5 py-1.5 px-3 text-xs"
        >
          <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          {{ createLabel }}
        </button>
        <!-- Extra actions slot -->
        <slot name="actions" />
      </div>
    </div>

    <!-- Loading skeleton -->
    <div v-if="loading" class="p-5 space-y-3">
      <div v-for="i in 5" :key="i" class="flex gap-4">
        <div class="h-4 bg-gray-100 rounded animate-pulse flex-1"></div>
        <div class="h-4 bg-gray-100 rounded animate-pulse w-24"></div>
        <div class="h-4 bg-gray-100 rounded animate-pulse w-16"></div>
      </div>
    </div>

    <!-- Table content -->
    <div v-else class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-gray-100">
            <slot name="head" />
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-50">
          <!-- Empty state -->
          <tr v-if="isEmpty">
            <td :colspan="colspan" class="py-16 text-center">
              <div class="flex flex-col items-center gap-2 text-gray-400">
                <i :class="[emptyIcon || 'icofont-inbox', 'text-5xl text-gray-200']"></i>
                <p class="text-sm">{{ emptyMessage || 'No records found' }}</p>
                <button
                  v-if="createLabel"
                  @click="$emit('create')"
                  class="mt-1 text-xs font-medium text-primary-700 hover:text-primary-600"
                >+ Create your first one</button>
              </div>
            </td>
          </tr>
          <!-- Rows -->
          <slot v-else name="rows" />
        </tbody>
      </table>
    </div>

    <!-- Pagination footer -->
    <div v-if="meta && meta.total > 0" class="admin-card-footer flex items-center justify-between">
      <p class="text-xs text-gray-500">
        Showing <span class="font-medium text-gray-700">{{ rangeStart }}–{{ rangeEnd }}</span> of <span class="font-medium text-gray-700">{{ meta.total }}</span>
      </p>
      <div class="flex items-center gap-1.5">
        <button
          @click="$emit('page-change', meta.page - 1)"
          :disabled="meta.page <= 1"
          class="px-3 py-1 text-xs font-medium border border-gray-200 rounded-lg text-gray-600 hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        >← Prev</button>
        <span class="px-2 text-xs text-gray-500 font-medium">{{ meta.page }}</span>
        <button
          @click="$emit('page-change', meta.page + 1)"
          :disabled="meta.page * meta.per_page >= meta.total"
          class="px-3 py-1 text-xs font-medium border border-gray-200 rounded-lg text-gray-600 hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        >Next →</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title:          { type: String,  default: '' },
  loading:        { type: Boolean, default: false },
  isEmpty:        { type: Boolean, default: false },
  meta:           { type: Object,  default: null },
  total:          { type: Number,  default: null },
  colspan:        { type: Number,  default: 6 },
  createLabel:    { type: String,  default: '' },
  emptyIcon:      { type: String,  default: '' },
  emptyMessage:   { type: String,  default: '' },
  searchable:     { type: Boolean, default: false },
  searchValue:    { type: String,  default: '' },
  searchPlaceholder: { type: String, default: '' }
})

defineEmits(['create', 'page-change', 'update:searchValue'])

const rangeStart = computed(() => props.meta ? (props.meta.page - 1) * props.meta.per_page + 1 : 0)
const rangeEnd   = computed(() => props.meta ? Math.min(props.meta.page * props.meta.per_page, props.meta.total) : 0)
</script>
