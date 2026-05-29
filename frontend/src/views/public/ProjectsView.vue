<template>
  <div>
    <!-- Mini Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">What We Do</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">Projects &amp; <span class="text-accent">Initiatives</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">Key projects and development initiatives being implemented by the National Council of Sports Uganda.</p>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Cards Grid -->
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
            class="bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-md transition-shadow border border-gray-100 flex flex-col group"
          >
            <div class="relative overflow-hidden">
              <img
                v-if="post.cover_image_url"
                :src="mediaUrl(post.cover_image_url)"
                :alt="post.title"
                class="w-full h-48 object-cover group-hover:scale-105 transition-transform duration-500"
              >
              <div v-else class="h-48 bg-gradient-to-br from-[#112b4e] to-[#1e4080] flex items-center justify-center">
                <svg class="w-12 h-12 text-white/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
                </svg>
              </div>
              <span class="absolute top-3 left-3 text-xs font-semibold bg-yellow-300 text-yellow-900 px-2.5 py-0.5 rounded-full">Initiative</span>
            </div>
            <div class="p-5 flex flex-col flex-1">
              <h2 class="font-bold text-darken mb-2 line-clamp-2 flex-1">{{ post.title }}</h2>
              <p class="text-sm text-gray-500 line-clamp-2 mb-4">{{ post.excerpt }}</p>
              <router-link
                :to="`/news/${post.slug}`"
                class="text-sm font-semibold text-accent hover:text-[#d47b05] flex items-center gap-1 transition-colors"
              >
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listPosts } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

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
