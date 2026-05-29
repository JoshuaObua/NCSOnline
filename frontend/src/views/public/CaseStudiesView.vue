<template>
  <div>
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-primary-400 font-semibold text-sm uppercase tracking-wide mb-2">Impact Stories</p>
        <h1 class="text-4xl font-bold text-white mb-4">Case Studies</h1>
        <p class="text-gray-300 max-w-2xl">In-depth stories of impact, transformation and progress in Uganda's sports sector.</p>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div v-if="loading" class="grid md:grid-cols-2 gap-8">
        <div v-for="i in 4" :key="i" class="animate-pulse">
          <div class="bg-gray-200 rounded-xl h-52 mb-4"></div>
          <div class="h-5 bg-gray-200 rounded mb-2"></div>
          <div class="h-3 bg-gray-200 rounded w-2/3"></div>
        </div>
      </div>

      <div v-else-if="posts.length" class="grid md:grid-cols-2 gap-8">
        <article v-for="post in posts" :key="post.id" class="group bg-white rounded-xl overflow-hidden shadow-sm hover:shadow-md transition-shadow border border-gray-100">
          <div v-if="post.cover_image_url" class="h-52 overflow-hidden">
            <img :src="mediaUrl(post.cover_image_url)" :alt="post.title" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300">
          </div>
          <div v-else class="h-52 bg-gradient-to-br from-primary-50 to-primary-100"></div>
          <div class="p-6">
            <span v-if="post.published_at" class="text-xs text-gray-400 mb-2 block">{{ formatDate(post.published_at) }}</span>
            <h2 class="font-bold text-gray-900 text-xl mb-3 line-clamp-2">{{ post.title }}</h2>
            <p class="text-gray-500 line-clamp-3 mb-4 text-sm">{{ post.excerpt }}</p>
            <router-link :to="`/news/${post.slug}`" class="text-sm text-primary-600 hover:text-primary-700 font-medium">
              Read case study →
            </router-link>
          </div>
        </article>
      </div>

      <div v-else class="text-center py-20 text-gray-400">
        <p>No case studies published yet.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listPosts } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const posts = ref([])
const loading = ref(true)

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(async () => {
  try {
    const res = await listPosts({ status: 'published', category: 'case_study', per_page: 12 })
    posts.value = res.data.data?.items || []
  } catch {
    posts.value = []
  } finally {
    loading.value = false
  }
})
</script>
