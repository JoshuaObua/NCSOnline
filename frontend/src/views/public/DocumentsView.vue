<template>
  <div>
    <!-- Mini Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">NCS Uganda</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">{{ heroLead }} <span class="text-accent">{{ heroAccent }}</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">{{ subtitle }}</p>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Category Filter -->
    <div v-if="categories.length" class="bg-white border-b border-gray-100 sticky top-16 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 flex flex-wrap gap-2">
        <button
          type="button"
          @click="activeCategory = ''"
          :aria-pressed="activeCategory === ''"
          :class="activeCategory === '' ? 'bg-[#112b4e] text-white shadow-sm' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
          class="px-4 py-1.5 rounded-full text-sm font-medium transition-colors"
        >All</button>
        <button
          v-for="cat in categories"
          :key="cat.id || cat.slug"
          type="button"
          @click="activeCategory = cat.slug"
          :aria-pressed="activeCategory === cat.slug"
          :class="activeCategory === cat.slug ? 'bg-[#112b4e] text-white shadow-sm' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
          class="px-4 py-1.5 rounded-full text-sm font-medium transition-colors"
        >{{ cat.name }}</button>
      </div>
    </div>

    <!-- List -->
    <div class="bg-white">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div v-if="loading" class="space-y-3">
          <div v-for="i in 6" :key="i" class="h-20 bg-gray-100 rounded-2xl animate-pulse"/>
        </div>

        <div v-else-if="filtered.length" class="space-y-3">
          <div v-for="d in filtered" :key="d.id">
            <!-- Embedded YouTube item -->
            <div
              v-if="d.video_url"
              class="bg-white border border-gray-100 rounded-2xl overflow-hidden hover:border-[#F48C06]/40 hover:shadow-md transition-all"
            >
              <button type="button" class="w-full flex items-center gap-4 p-5 text-left" @click="toggleVideo(d.id)">
                <div class="flex-shrink-0 w-14 h-14 bg-red-50 rounded-2xl flex items-center justify-center">
                  <svg viewBox="0 0 24 24" class="w-6 h-6 text-red-500" fill="currentColor" aria-hidden="true"><path d="M8 5v14l11-7z"/></svg>
                </div>
                <div class="flex-1 min-w-0">
                  <h3 class="font-bold text-darken truncate">{{ d.title }}</h3>
                  <p v-if="d.description" class="text-sm text-gray-500 line-clamp-1 mt-0.5">{{ d.description }}</p>
                  <span v-if="d.category" class="inline-block text-xs font-semibold bg-yellow-100 text-yellow-800 px-2.5 py-0.5 rounded-full capitalize mt-1.5">{{ categoryLabel(d.category) }}</span>
                </div>
                <i class="flex-shrink-0 icofont-rounded-down text-xl text-gray-400 transition-transform" :class="{ 'rotate-180': expanded === d.id }" aria-hidden="true"></i>
              </button>
              <div v-if="expanded === d.id && youtubeEmbedUrl(d.video_url)" class="aspect-video border-t border-gray-100">
                <iframe
                  :src="youtubeEmbedUrl(d.video_url)"
                  class="w-full h-full"
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                  allowfullscreen
                  loading="lazy"
                ></iframe>
              </div>
            </div>

            <!-- PDF item -->
            <a
              v-else-if="d.file_url"
              :href="mediaUrl(d.file_url)"
              target="_blank"
              rel="noopener"
              :aria-label="`${d.title} (opens in a new tab)`"
              class="flex items-center gap-4 bg-white border border-gray-100 rounded-2xl p-5 hover:border-[#F48C06]/40 hover:shadow-md transition-all group"
            >
              <div class="flex-shrink-0 w-14 h-14 bg-red-50 rounded-2xl flex items-center justify-center">
                <i class="icofont-file-pdf text-3xl text-red-500" aria-hidden="true"></i>
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="font-bold text-darken truncate group-hover:text-accent transition-colors">{{ d.title }}</h3>
                <p v-if="d.description" class="text-sm text-gray-500 line-clamp-1 mt-0.5">{{ d.description }}</p>
                <span v-if="d.category" class="inline-block text-xs font-semibold bg-yellow-100 text-yellow-800 px-2.5 py-0.5 rounded-full capitalize mt-1.5">{{ categoryLabel(d.category) }}</span>
              </div>
              <div class="flex-shrink-0 w-10 h-10 rounded-xl bg-[#112b4e]/5 group-hover:bg-[#F48C06] flex items-center justify-center transition-colors">
                <i class="icofont-download text-xl text-[#112b4e] group-hover:text-white transition-colors" aria-hidden="true"></i>
              </div>
            </a>

            <!-- No file or video attached yet -->
            <div v-else class="flex items-center gap-4 bg-gray-50 border border-dashed border-gray-200 rounded-2xl p-5">
              <div class="flex-shrink-0 w-14 h-14 bg-gray-100 rounded-2xl flex items-center justify-center">
                <i :class="[icon, 'text-3xl text-gray-300']" aria-hidden="true"></i>
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="font-bold text-darken truncate">{{ d.title }}</h3>
                <p v-if="d.description" class="text-sm text-gray-500 line-clamp-1 mt-0.5">{{ d.description }}</p>
                <span v-if="d.category" class="inline-block text-xs font-semibold bg-yellow-100 text-yellow-800 px-2.5 py-0.5 rounded-full capitalize mt-1.5">{{ categoryLabel(d.category) }}</span>
              </div>
              <span class="flex-shrink-0 text-xs text-gray-400 italic">No file attached</span>
            </div>
          </div>
        </div>

        <div v-else class="text-center py-20 text-gray-400">
          <i :class="[icon, 'text-5xl mb-3 opacity-30']" aria-hidden="true"></i>
          <p>{{ emptyText }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { listDocuments, listCategoriesByType } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const props = defineProps({
  docType: { type: String, required: true },
  heroLead: { type: String, default: '' },
  heroAccent: { type: String, default: '' },
  subtitle: { type: String, default: '' },
  icon: { type: String, default: 'icofont-file-pdf' },
  emptyText: { type: String, default: 'Nothing published here yet.' },
})

const items = ref([])
const categories = ref([])
const loading = ref(true)
const activeCategory = ref('')
const expanded = ref('')

const filtered = computed(() => {
  if (!activeCategory.value) return items.value
  return items.value.filter(d => d.category === activeCategory.value)
})

function toggleVideo(id) { expanded.value = expanded.value === id ? '' : id }

function categoryLabel(slug) {
  const c = categories.value.find(x => x.slug === slug)
  return c ? c.name : slug
}

function youtubeEmbedUrl(url) {
  try {
    const u = new URL(url)
    if (u.hostname.includes('youtu.be')) return `https://www.youtube.com/embed${u.pathname}`
    const id = u.searchParams.get('v')
    if (id) return `https://www.youtube.com/embed/${id}`
    if (u.pathname.startsWith('/embed/')) return url
  } catch { /* not a valid URL */ }
  return ''
}

async function load() {
  loading.value = true
  activeCategory.value = ''
  try {
    const [docsRes, catsRes] = await Promise.all([
      listDocuments({ doc_type: props.docType }),
      listCategoriesByType(props.docType),
    ])
    items.value = docsRes.data.data || []
    categories.value = catsRes.data.data || []
  } catch { items.value = []; categories.value = [] }
  finally { loading.value = false }
}

onMounted(load)
watch(() => props.docType, load)
</script>
