<template>
  <div>
    <!-- Header -->
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-primary-400 font-semibold text-sm uppercase tracking-wide mb-2">NCS Uganda</p>
        <h1 class="text-4xl font-bold text-white mb-4">News &amp; Updates</h1>
        <p class="text-gray-300 max-w-2xl">Latest news, announcements and updates from the National Council of Sports Uganda.</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white border-b border-gray-100 sticky top-16 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 flex flex-wrap gap-2">
        <button
          v-for="cat in categories"
          :key="cat.value"
          @click="activeCategory = cat.value; loadPosts()"
          :class="activeCategory === cat.value
            ? 'bg-primary-600 text-white'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
          class="px-4 py-1.5 rounded-full text-sm font-medium transition-colors"
        >
          {{ cat.label }}
        </button>
      </div>
    </div>

    <!-- Posts Grid -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div v-if="loading" class="grid md:grid-cols-3 gap-6">
        <div v-for="i in 6" :key="i" class="animate-pulse">
          <div class="bg-gray-200 rounded-xl h-44 mb-4"></div>
          <div class="h-4 bg-gray-200 rounded mb-2"></div>
          <div class="h-3 bg-gray-200 rounded w-2/3"></div>
        </div>
      </div>

      <div v-else-if="posts.length" class="grid md:grid-cols-3 gap-6">
        <article v-for="post in posts" :key="post.id" class="bg-white rounded-xl overflow-hidden shadow-sm hover:shadow-md transition-shadow border border-gray-100">
          <div v-if="post.cover_image_url" class="h-44 overflow-hidden">
            <img :src="mediaUrl(post.cover_image_url)" :alt="post.title" class="w-full h-full object-cover hover:scale-105 transition-transform duration-300">
          </div>
          <div v-else class="h-44 bg-gradient-to-br from-primary-50 to-primary-100 flex items-center justify-center">
            <span class="text-primary-300 font-bold text-3xl">NCS</span>
          </div>
          <div class="p-5">
            <div class="flex items-center gap-2 mb-2">
              <span class="text-xs font-medium text-primary-600 bg-primary-50 px-2 py-0.5 rounded-full capitalize">{{ post.category }}</span>
              <span v-if="post.published_at" class="text-xs text-gray-400">{{ formatDate(post.published_at) }}</span>
            </div>
            <h2 class="font-semibold text-gray-900 mb-2 line-clamp-2">{{ post.title }}</h2>
            <p class="text-sm text-gray-500 line-clamp-3 mb-4">{{ post.excerpt }}</p>
            <router-link :to="`/news/${post.slug}`" class="text-sm text-primary-600 hover:text-primary-700 font-medium">
              Read more →
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
      <div v-if="total > perPage" class="mt-10 flex justify-center gap-2">
        <button @click="page--; loadPosts()" :disabled="page<=1" class="px-4 py-2 border rounded-lg text-sm disabled:opacity-40 hover:bg-gray-50">Previous</button>
        <span class="px-4 py-2 text-sm text-gray-500">Page {{ page }} of {{ Math.ceil(total/perPage) }}</span>
        <button @click="page++; loadPosts()" :disabled="page*perPage>=total" class="px-4 py-2 border rounded-lg text-sm disabled:opacity-40 hover:bg-gray-50">Next</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listPosts } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const categories = [
  { value: '', label: 'All' },
  { value: 'blog', label: 'Blog' },
  { value: 'news', label: 'News' },
  { value: 'announcement', label: 'Announcements' },
]

const posts = ref([])
const total = ref(0)
const loading = ref(true)
const activeCategory = ref('')
const page = ref(1)
const perPage = 9

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { day: 'numeric', month: 'short', year: 'numeric' })
}

async function loadPosts() {
  loading.value = true
  try {
    const res = await listPosts({ status: 'published', category: activeCategory.value, page: page.value, per_page: perPage })
    posts.value = res.data.data?.items || []
    total.value = res.data.data?.total || 0
  } catch {
    posts.value = []
  } finally {
    loading.value = false
  }
}

onMounted(loadPosts)
</script>
