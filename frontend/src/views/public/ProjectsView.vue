<template>
  <div>
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-primary-400 font-semibold text-sm uppercase tracking-wide mb-2">What We Do</p>
        <h1 class="text-4xl font-bold text-white mb-4">Projects &amp; Initiatives</h1>
        <p class="text-gray-300 max-w-2xl">Key projects and development initiatives being implemented by the National Council of Sports Uganda.</p>
      </div>
    </div>

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
            <img :src="post.cover_image_url" :alt="post.title" class="w-full h-full object-cover hover:scale-105 transition-transform duration-300">
          </div>
          <div v-else class="h-44 bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center">
            <svg class="w-12 h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/></svg>
          </div>
          <div class="p-5">
            <h2 class="font-semibold text-gray-900 mb-2 line-clamp-2">{{ post.title }}</h2>
            <p class="text-sm text-gray-500 line-clamp-3 mb-4">{{ post.excerpt }}</p>
            <router-link :to="`/news/${post.slug}`" class="text-sm text-primary-600 hover:text-primary-700 font-medium">
              Learn more →
            </router-link>
          </div>
        </article>
      </div>

      <div v-else class="text-center py-20 text-gray-400">
        <p>No projects to display at this time.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listPosts } from '@/api/cms.js'

const posts = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await listPosts({ status: 'published', category: 'project', per_page: 12 })
    posts.value = res.data.data?.items || []
  } catch {
    posts.value = []
  } finally {
    loading.value = false
  }
})
</script>
