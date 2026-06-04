<template>
  <div>
    <!-- Loading skeleton -->
    <div v-if="loading" class="bg-cream pt-24 pb-16">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 animate-pulse">
        <div class="h-6 bg-gray-300 rounded w-1/4 mb-4"></div>
        <div class="h-12 bg-gray-300 rounded mb-3 w-3/4"></div>
        <div class="h-4 bg-gray-300 rounded w-1/3"></div>
      </div>
    </div>

    <div v-else-if="post">
      <!-- Hero image banner with dark overlay + title -->
      <div class="relative min-h-[360px] flex items-end" :style="post.cover_image_url ? '' : ''">
        <div
          v-if="post.cover_image_url"
          class="absolute inset-0 bg-cover bg-center"
          :style="`background-image: url('${mediaUrl(post.cover_image_url)}')`"
        ></div>
        <div v-else class="absolute inset-0 bg-gradient-to-br from-[#112b4e] to-[#1e4080]"></div>
        <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent"></div>

        <div class="relative z-10 max-w-4xl mx-auto px-4 sm:px-6 pb-12 pt-32 w-full">
          <!-- Breadcrumb -->
          <nav class="flex items-center gap-2 text-sm text-white/60 mb-5">
            <router-link to="/" class="hover:text-white transition-colors">Home</router-link>
            <span>/</span>
            <router-link to="/news" class="hover:text-white transition-colors">News</router-link>
            <span>/</span>
            <span class="text-white/80 truncate max-w-xs">{{ post.title }}</span>
          </nav>

          <!-- Meta + title -->
          <div class="flex items-center gap-3 mb-4 flex-wrap">
            <span class="text-xs font-semibold bg-yellow-300 text-yellow-900 px-3 py-1 rounded-full capitalize">{{ post.category || 'News' }}</span>
            <span v-if="post.published_at" class="text-sm text-white/70">{{ formatDate(post.published_at) }}</span>
            <span v-if="post.author_name" class="text-sm text-white/70">By {{ post.author_name }}</span>
            <span v-if="post.view_count" class="text-sm text-white/60 flex items-center gap-1">
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.964-7.178z"/><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
              {{ post.view_count }}
            </span>
          </div>
          <h1 class="text-3xl md:text-4xl font-bold text-white leading-tight">{{ post.title }}</h1>
        </div>
      </div>

      <!-- Article content -->
      <div class="bg-white">
        <div class="max-w-4xl mx-auto px-4 sm:px-6 py-14">
          <!-- Excerpt block -->
          <p v-if="post.excerpt" class="text-lg text-gray-500 mb-8 border-l-4 border-[#F48C06] pl-5 italic">{{ post.excerpt }}</p>

          <!-- Body -->
          <div
            class="prose prose-gray max-w-none prose-headings:text-[#112b4e] prose-a:text-[#F48C06] prose-a:no-underline hover:prose-a:underline prose-img:rounded-xl prose-p:text-gray-600 prose-p:leading-relaxed"
            v-html="post.content"
          ></div>

          <!-- Footer -->
          <div class="mt-14 pt-8 border-t border-gray-100 flex items-center justify-between flex-wrap gap-4">
            <router-link
              to="/news"
              class="text-sm font-semibold text-[#112b4e] hover:text-[#F48C06] flex items-center gap-2 transition-colors"
            >
              ← Back to News
            </router-link>
            <span v-if="post.published_at" class="text-xs text-gray-400">Published {{ formatDate(post.published_at) }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="bg-cream pt-32 pb-20 text-center text-gray-400">
      <p class="text-lg font-medium mb-4">Article not found.</p>
      <router-link to="/news" class="text-sm font-semibold text-accent hover:text-[#d47b05]">← Browse all articles</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getPost } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const route = useRoute()
const post = ref(null)
const loading = ref(true)

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(async () => {
  try {
    const res = await getPost(route.params.slug)
    post.value = res.data.data
  } catch {
    post.value = null
  } finally {
    loading.value = false
  }
})
</script>
