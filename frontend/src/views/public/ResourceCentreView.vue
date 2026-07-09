<template>
  <div>
    <!-- Mini Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">NCS Uganda</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">Resource <span class="text-accent">Centre</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">Download official NCS guidelines, reports, press releases, sports rules and other public resources.</p>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Category Filter -->
    <div class="bg-white border-b border-gray-100 sticky top-16 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 flex flex-wrap gap-2">
        <button
          v-for="cat in categories"
          :key="cat.value"
          @click="activeCategory = cat.value"
          :aria-pressed="activeCategory === cat.value"
          :class="activeCategory === cat.value
            ? 'bg-[#112b4e] text-white shadow-sm'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
          class="px-4 py-1.5 rounded-full text-sm font-medium transition-colors"
          :id="`cat-${cat.value || 'all'}`"
        >
          {{ cat.label }}
        </button>
      </div>
    </div>

    <!-- Resources List -->
    <div class="bg-white">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div v-if="loading" class="space-y-3">
          <div v-for="i in 6" :key="i" class="h-20 bg-gray-100 rounded-2xl animate-pulse"/>
        </div>

        <div v-else-if="filtered.length" class="space-y-3">
          <a
            v-for="r in filtered"
            :key="r.id"
            :href="mediaUrl(r.file_url)"
            target="_blank"
            rel="noopener"
            :aria-label="`${r.title} (opens in a new tab)`"
            class="flex items-center gap-4 bg-white border border-gray-100 rounded-2xl p-5 hover:border-[#F48C06]/40 hover:shadow-md transition-all group"
          >
            <!-- File icon -->
            <div class="flex-shrink-0 w-14 h-14 bg-red-50 rounded-2xl flex items-center justify-center">
              <i class="icofont-file-pdf text-3xl text-red-500" aria-hidden="true"></i>
            </div>

            <div class="flex-1 min-w-0">
              <h3 class="font-bold text-darken truncate group-hover:text-accent transition-colors">{{ r.title }}</h3>
              <p v-if="r.description" class="text-sm text-gray-500 line-clamp-1 mt-0.5">{{ r.description }}</p>
              <div class="flex items-center gap-2 mt-1.5">
                <span class="text-xs font-semibold bg-yellow-100 text-yellow-800 px-2.5 py-0.5 rounded-full capitalize">
                  {{ categoryLabel(r.category) }}
                </span>
              </div>
            </div>

            <!-- Download arrow -->
            <div class="flex-shrink-0 w-10 h-10 rounded-xl bg-[#112b4e]/5 group-hover:bg-[#F48C06] flex items-center justify-center transition-colors">
              <i class="icofont-download text-xl text-[#112b4e] group-hover:text-white transition-colors" aria-hidden="true"></i>
            </div>
          </a>
        </div>

        <div v-else class="text-center py-20 text-gray-400">
          <i class="icofont-folder-open text-5xl mb-3 opacity-30"></i>
          <p>No resources available in this category yet.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { listResources } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

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

  const hash = (route.hash || '').replace('#', '')
  if (hash && categories.some(c => c.value === hash)) {
    activeCategory.value = hash
  }
})
</script>
