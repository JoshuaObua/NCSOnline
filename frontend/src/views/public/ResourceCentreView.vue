<template>
  <div>
    <!-- Header -->
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-primary-400 font-semibold text-sm uppercase tracking-wide mb-2">NCS Uganda</p>
        <h1 class="text-4xl font-bold text-white mb-4">Resource Centre</h1>
        <p class="text-gray-300 max-w-2xl">Download official NCS guidelines, reports, press releases, sports rules and other public resources.</p>
      </div>
    </div>

    <!-- Category Filter -->
    <div class="bg-white border-b border-gray-100 sticky top-16 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 flex flex-wrap gap-2">
        <button
          v-for="cat in categories"
          :key="cat.value"
          @click="activeCategory = cat.value"
          :class="activeCategory === cat.value
            ? 'bg-primary-600 text-white'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
          class="px-4 py-1.5 rounded-full text-sm font-medium transition-colors"
          :id="`cat-${cat.value || 'all'}`"
        >
          {{ cat.label }}
        </button>
      </div>
    </div>

    <!-- Resources List -->
    <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 6" :key="i" class="h-20 bg-gray-100 rounded-xl animate-pulse"/>
      </div>

      <div v-else-if="filtered.length" class="space-y-3">
        <a
          v-for="r in filtered"
          :key="r.id"
          :href="r.file_url"
          target="_blank"
          rel="noopener"
          class="block bg-white border border-gray-200 rounded-xl p-5 hover:border-primary-300 hover:shadow-md transition-all flex items-center gap-4"
        >
          <div class="flex-shrink-0 w-14 h-14 bg-red-50 rounded-xl flex items-center justify-center">
            <i class="icofont-file-pdf text-3xl text-red-500"></i>
          </div>
          <div class="flex-1 min-w-0">
            <h3 class="font-semibold text-gray-900 truncate">{{ r.title }}</h3>
            <p v-if="r.description" class="text-sm text-gray-500 line-clamp-1 mt-0.5">{{ r.description }}</p>
            <div class="flex items-center gap-2 mt-1.5">
              <span class="text-xs font-medium text-primary-700 bg-primary-50 px-2 py-0.5 rounded-full capitalize">
                {{ categoryLabel(r.category) }}
              </span>
            </div>
          </div>
          <div class="flex-shrink-0 text-primary-600 group-hover:translate-x-1 transition-transform">
            <i class="icofont-download text-2xl"></i>
          </div>
        </a>
      </div>

      <div v-else class="text-center py-20 text-gray-400">
        <i class="icofont-folder-open text-5xl mb-3 opacity-40"></i>
        <p>No resources available in this category yet.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { listResources } from '@/api/cms.js'

const route = useRoute()

const categories = [
  { value: '',              label: 'All' },
  { value: 'guidelines',    label: 'Guidelines & Standards' },
  { value: 'reports',       label: 'Reports' },
  { value: 'speeches',      label: 'Speeches' },
  { value: 'press-release', label: 'Press Releases' },
  { value: 'sports-rules',  label: 'Sports Rules & Laws' },
  { value: 'downloads',     label: 'Important Downloads' },
]

const resources = ref([])
const loading = ref(true)
const activeCategory = ref('')

const filtered = computed(() => {
  if (!activeCategory.value) return resources.value
  return resources.value.filter(r => r.category === activeCategory.value)
})

function categoryLabel(v) {
  const c = categories.find(x => x.value === v)
  return c ? c.label : (v || 'Other')
}

onMounted(async () => {
  try {
    const r = await listResources()
    resources.value = r.data.data || []
  } catch { resources.value = [] }
  finally { loading.value = false }

  // Hash anchor (e.g. /resource-centre#reports) sets the filter
  const hash = (route.hash || '').replace('#', '')
  if (hash && categories.some(c => c.value === hash)) {
    activeCategory.value = hash
  }
})
</script>
