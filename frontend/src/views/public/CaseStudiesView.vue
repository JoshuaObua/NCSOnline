<template>
  <div>
    <!-- Mini Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">Impact Stories</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">Case <span class="text-accent">Studies</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">In-depth stories of impact, transformation and progress in Uganda's sports sector.</p>
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
        <div v-if="loading" class="grid md:grid-cols-2 gap-8">
          <div v-for="i in 4" :key="i" class="animate-pulse">
            <div class="bg-gray-200 rounded-2xl h-52 mb-4"></div>
            <div class="h-5 bg-gray-200 rounded mb-2"></div>
            <div class="h-3 bg-gray-200 rounded w-2/3"></div>
          </div>
        </div>

        <div v-else-if="posts.length" class="grid md:grid-cols-2 gap-8">
          <article
            v-for="post in posts"
            :key="post.id"
            class="group bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-lg transition-shadow border border-gray-100"
          >
            <div class="relative overflow-hidden">
              <img
                v-if="post.cover_image_url"
                :src="mediaUrl(post.cover_image_url)"
                :alt="post.title"
                class="w-full h-52 object-cover group-hover:scale-105 transition-transform duration-500"
              >
              <div v-else class="h-52 bg-gradient-to-br from-[#112b4e] to-[#1e4080] flex items-center justify-center">
                <span class="text-white/20 font-bold text-4xl">NCS</span>
              </div>
              <span class="absolute top-3 left-3 text-xs font-semibold bg-yellow-300 text-yellow-900 px-2.5 py-0.5 rounded-full">Case Study</span>
            </div>
            <div class="p-6">
              <span v-if="post.published_at" class="text-xs text-gray-400 mb-2 block">{{ formatDate(post.published_at) }}</span>
              <h2 class="font-bold text-darken text-xl mb-3 line-clamp-2">{{ post.title }}</h2>
              <p class="text-gray-500 line-clamp-3 mb-5 text-sm">{{ post.excerpt }}</p>
              <router-link
                :to="`/news/${post.slug}`"
                class="text-sm font-semibold text-accent hover:text-[#d47b05] flex items-center gap-1 transition-colors"
              >
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
