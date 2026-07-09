<template>
  <div>
    <div class="bg-gray-50 pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <span class="inline-block px-4 py-1.5 bg-[#f5a623]/10 text-[#f5a623] text-sm font-semibold rounded-full mb-3">Stay Updated</span>
        <h1 class="text-3xl md:text-5xl font-bold text-[#1a365d] mb-4">Latest News</h1>
        <p class="text-gray-600 max-w-2xl text-lg leading-relaxed">Latest news, announcements and updates from the National Council of Sports Uganda.</p>
      </div>
    </div>

    <!-- Filter Chips -->
    <div class="bg-gray-50 border-y border-gray-100 sticky top-16 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 flex flex-wrap gap-2">
        <button
          v-for="cat in categories"
          :key="cat.value"
          @click="activeCategory = cat.value; loadPosts()"
          :aria-pressed="activeCategory === cat.value"
          :class="activeCategory === cat.value
            ? 'bg-[#1a365d] text-white'
            : 'bg-white text-gray-600 hover:bg-[#f5a623]/10 hover:text-[#f5a623]'"
          class="rounded-full px-4 py-2 text-sm transition-all"
        >
          {{ cat.label }}
        </button>
      </div>
    </div>

    <!-- Posts Grid -->
    <div class="bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div v-if="searchQuery" class="mb-8 flex items-center justify-between gap-4 rounded-xl bg-[#f5a623]/10 px-4 py-3 text-sm text-[#1a365d]"><span>Search results for <strong>"{{ searchQuery }}"</strong></span><button type="button" class="font-semibold text-[#d88700]" @click="clearSearch">Clear</button></div>
        <div v-if="loading" class="grid md:grid-cols-3 gap-6">
          <div v-for="i in 6" :key="i" class="animate-pulse">
            <div class="bg-gray-200 rounded-2xl h-48 mb-4"></div>
            <div class="h-4 bg-gray-200 rounded mb-2"></div>
            <div class="h-3 bg-gray-200 rounded w-2/3"></div>
          </div>
        </div>

        <div v-else-if="displayedPosts.length" class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          <article
            v-for="post in displayedPosts"
            :key="post.id"
            class="group bg-white rounded-xl overflow-hidden shadow-sm hover:shadow-lg transition-all duration-300 border border-gray-100 hover:border-[#f5a623]/30 flex flex-col"
          >
            <div class="relative overflow-hidden">
              <img
                v-if="post.cover_image_url"
                :src="mediaUrl(post.cover_image_url)"
                :alt="post.title"
                class="w-full h-56 object-cover transition-transform duration-500 group-hover:scale-105"
              >
              <div v-else class="h-56 bg-gradient-to-br from-[#1a365d] to-[#2d4a6f] flex items-center justify-center">
                <span class="text-white/40 font-bold text-3xl">NCS</span>
              </div>
              <!-- Category badge overlay -->
              <span class="absolute top-4 left-4 inline-flex items-center rounded-md border border-transparent bg-[#f5a623] px-2.5 py-0.5 text-xs font-semibold text-white shadow capitalize">
                {{ post.category || 'News' }}
              </span>
            </div>
            <div class="p-6 flex flex-col flex-1">
              <p v-if="post.published_at" class="text-xs text-gray-400 mb-2">{{ formatDate(post.published_at) }}</p>
              <h2 class="font-semibold text-[#1a365d] group-hover:text-[#f5a623] text-xl leading-snug mb-3 line-clamp-2 flex-1 transition-colors">{{ post.title }}</h2>
              <p class="text-base text-gray-500 leading-relaxed line-clamp-3 mb-5">{{ post.excerpt }}</p>
              <router-link
                :to="`/news/${post.slug}`"
                class="text-sm font-semibold text-[#f5a623] hover:text-[#d88700] inline-flex items-center gap-2 transition-colors"
              >
                Read more
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 transition-transform group-hover:translate-x-1" aria-hidden="true"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
              </router-link>
            </div>
          </article>
        </div>

        <div v-else class="text-center py-20 text-gray-400">
            <svg class="w-16 h-16 mx-auto mb-4 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"/>
          </svg>
          <p>No articles in this category yet.</p>
        </div>

        <!-- Pagination -->
        <div v-if="total > perPage" class="mt-12 flex justify-center items-center gap-3">
          <button
            @click="page--; loadPosts()"
            :disabled="page <= 1"
            class="px-5 py-2 rounded-full border border-gray-300 text-sm font-medium text-gray-600 hover:border-[#1a365d] hover:text-[#1a365d] disabled:opacity-40 transition-colors"
          >Previous</button>
          <span class="text-sm text-gray-500" aria-live="polite">Page {{ page }} of {{ Math.ceil(total / perPage) }}</span>
          <button
            @click="page++; loadPosts()"
            :disabled="page * perPage >= total"
            class="px-5 py-2 rounded-full border border-gray-300 text-sm font-medium text-gray-600 hover:border-[#1a365d] hover:text-[#1a365d] disabled:opacity-40 transition-colors"
          >Next</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { listPosts, getSettings, listBlogCategories } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const EXCLUDED = ['page', 'case_study', 'project']
const route = useRoute()
const router = useRouter()

const categories = ref([{ value: '', label: 'All' }])
const posts = ref([])
const total = ref(0)
const loading = ref(true)
const activeCategory = ref('')
const page = ref(1)
const perPage = 9
const searchQuery = computed(() => String(route.query.search || '').trim())
const displayedPosts = computed(() => { const q=searchQuery.value.toLowerCase(); return q ? posts.value.filter(post=>[post.title,post.excerpt,post.content,post.category].some(value=>String(value||'').toLowerCase().includes(q))) : posts.value })
function clearSearch(){router.replace({path:'/news'})}

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { day: 'numeric', month: 'short', year: 'numeric' })
}

async function loadCategories() {
  try {
    const res = await listBlogCategories(true)
    const val = res.data?.data || res.data || []
    if (Array.isArray(val) && val.length) {
      categories.value = [{ value: '', label: 'All' }, ...val.filter(c => !EXCLUDED.includes(c.slug)).map(c => ({ value: c.slug, label: c.name }))]
      return
    }
    const settings = await getSettings('post_categories')
    const legacy = settings.data?.data?.value
    if (Array.isArray(legacy) && legacy.length) {
      categories.value = [{ value: '', label: 'All' }, ...legacy.filter(c => !EXCLUDED.includes(c.value))]
    }
  } catch {
    categories.value = [
      { value: '', label: 'All' },
      { value: 'blog', label: 'Blog' },
      { value: 'news', label: 'News' },
      { value: 'announcement', label: 'Announcements' },
    ]
  }
}

async function loadPosts() {
  loading.value = true
  try {
    const res = await listPosts({ status: 'published', category: activeCategory.value, page: page.value, per_page: perPage })
    let items = res.data.data?.items || []
    if (!activeCategory.value) {
      items = items.filter(p => !EXCLUDED.includes(p.category))
    }
    posts.value = items
    total.value = res.data.data?.total || 0
  } catch {
    posts.value = []
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadCategories()
  loadPosts()
})
watch(()=>route.query.search,()=>loadPosts())
</script>
