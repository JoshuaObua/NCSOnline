<template>
  <div>
    <!-- Mini Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">NCS Uganda</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">News &amp; <span class="text-accent">Updates</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">Latest news, announcements and updates from the National Council of Sports Uganda.</p>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Filter Chips -->
    <div class="bg-white border-b border-gray-100 sticky top-16 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 flex flex-wrap gap-2">
        <button
          v-for="cat in categories"
          :key="cat.value"
          @click="activeCategory = cat.value; loadPosts()"
          :class="activeCategory === cat.value
            ? 'bg-[#112b4e] text-white shadow-sm'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
          class="px-4 py-1.5 rounded-full text-sm font-medium transition-colors"
        >
          {{ cat.label }}
        </button>
      </div>
    </div>

    <!-- Posts Grid -->
    <div class="bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div v-if="loading" class="grid md:grid-cols-3 gap-6">
          <div v-for="i in 6" :key="i" class="animate-pulse">
            <div class="bg-gray-200 rounded-2xl h-48 mb-4"></div>
            <div class="h-4 bg-gray-200 rounded mb-2"></div>
            <div class="h-3 bg-gray-200 rounded w-2/3"></div>
          </div>
        </div>

        <div v-else-if="posts.length" class="grid md:grid-cols-3 gap-6">
          <article
            v-for="post in posts"
            :key="post.id"
            class="bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-md transition-shadow border border-gray-100 flex flex-col"
          >
            <div class="relative overflow-hidden">
              <img
                v-if="post.cover_image_url"
                :src="mediaUrl(post.cover_image_url)"
                :alt="post.title"
                class="w-full h-48 object-cover hover:scale-105 transition-transform duration-500"
              >
              <div v-else class="h-48 bg-gradient-to-br from-[#112b4e] to-[#1e4080] flex items-center justify-center">
                <span class="text-white/40 font-bold text-3xl">NCS</span>
              </div>
              <!-- Category badge overlay -->
              <span class="absolute top-3 left-3 text-xs font-semibold bg-yellow-300 text-yellow-900 px-2.5 py-0.5 rounded-full capitalize">
                {{ post.category || 'News' }}
              </span>
            </div>
            <div class="p-5 flex flex-col flex-1">
              <p v-if="post.published_at" class="text-xs text-gray-400 mb-2">{{ formatDate(post.published_at) }}</p>
              <h2 class="font-bold text-darken text-base mb-2 line-clamp-2 flex-1">{{ post.title }}</h2>
              <p class="text-sm text-gray-500 line-clamp-2 mb-4">{{ post.excerpt }}</p>
              <router-link
                :to="`/news/${post.slug}`"
                class="text-sm font-semibold text-accent hover:text-[#d47b05] flex items-center gap-1 transition-colors"
              >
                Read more <span class="text-base leading-none">→</span>
              </router-link>
            </div>
          </article>
        </div>

        <div v-else class="text-center py-20 text-gray-400">
          <svg class="w-16 h-16 mx-auto mb-4 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"/>
          </svg>
          <p>No articles in this category yet.</p>
        </div>

        <!-- Pagination -->
        <div v-if="total > perPage" class="mt-12 flex justify-center items-center gap-3">
          <button
            @click="page--; loadPosts()"
            :disabled="page <= 1"
            class="px-5 py-2 rounded-full border border-gray-300 text-sm font-medium text-gray-600 hover:border-[#112b4e] hover:text-[#112b4e] disabled:opacity-40 transition-colors"
          >← Previous</button>
          <span class="text-sm text-gray-500">Page {{ page }} of {{ Math.ceil(total / perPage) }}</span>
          <button
            @click="page++; loadPosts()"
            :disabled="page * perPage >= total"
            class="px-5 py-2 rounded-full border border-gray-300 text-sm font-medium text-gray-600 hover:border-[#112b4e] hover:text-[#112b4e] disabled:opacity-40 transition-colors"
          >Next →</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listPosts, getSettings } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const EXCLUDED = ['page', 'case_study']

const categories = ref([{ value: '', label: 'All' }])
const posts = ref([])
const total = ref(0)
const loading = ref(true)
const activeCategory = ref('')
const page = ref(1)
const perPage = 9

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { day: 'numeric', month: 'short', year: 'numeric' })
}

async function loadCategories() {
  try {
    const res = await getSettings('post_categories')
    const val = res.data?.data?.value
    if (Array.isArray(val) && val.length) {
      categories.value = [{ value: '', label: 'All' }, ...val.filter(c => !EXCLUDED.includes(c.value))]
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
</script>
