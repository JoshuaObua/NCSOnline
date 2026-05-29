<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div v-if="loading" class="animate-pulse">
      <div class="h-8 bg-gray-200 rounded mb-4 w-3/4"></div>
      <div class="h-4 bg-gray-200 rounded mb-8 w-1/2"></div>
      <div class="h-64 bg-gray-200 rounded-xl mb-8"></div>
      <div class="space-y-3">
        <div class="h-4 bg-gray-200 rounded"></div>
        <div class="h-4 bg-gray-200 rounded"></div>
        <div class="h-4 bg-gray-200 rounded w-4/5"></div>
      </div>
    </div>

    <div v-else-if="post">
      <!-- Breadcrumb -->
      <nav class="flex items-center gap-2 text-sm text-gray-400 mb-8">
        <router-link to="/" class="hover:text-primary-600">Home</router-link>
        <span>/</span>
        <router-link to="/news" class="hover:text-primary-600">News</router-link>
        <span>/</span>
        <span class="text-gray-600 truncate max-w-xs">{{ post.title }}</span>
      </nav>

      <!-- Meta -->
      <div class="flex flex-wrap items-center gap-3 mb-4">
        <span class="text-xs font-medium text-primary-600 bg-primary-50 px-3 py-1 rounded-full capitalize">{{ post.category }}</span>
        <span v-if="post.published_at" class="text-sm text-gray-400">{{ formatDate(post.published_at) }}</span>
      </div>

      <h1 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">{{ post.title }}</h1>

      <p v-if="post.excerpt" class="text-lg text-gray-500 mb-8 border-l-4 border-primary-500 pl-4">{{ post.excerpt }}</p>

      <img v-if="post.cover_image_url" :src="mediaUrl(post.cover_image_url)" :alt="post.title" class="w-full rounded-xl mb-10 max-h-80 object-cover">

      <!-- Content -->
      <div class="prose prose-gray max-w-none prose-headings:text-gray-900 prose-a:text-primary-600" v-html="formattedContent"></div>

      <div class="mt-12 pt-8 border-t border-gray-100">
        <router-link to="/news" class="text-primary-600 hover:text-primary-700 font-medium text-sm">← Back to News</router-link>
      </div>
    </div>

    <div v-else class="text-center py-20 text-gray-400">
      <p class="text-lg">Article not found.</p>
      <router-link to="/news" class="mt-4 inline-block text-primary-600 hover:underline">Browse all articles</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getPost } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const route = useRoute()
const post = ref(null)
const loading = ref(true)

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}

const formattedContent = computed(() => {
  if (!post.value?.content) return ''
  // Render newlines as paragraphs for plain text content
  return post.value.content
    .split('\n\n')
    .map(p => `<p>${p.replace(/\n/g, '<br>')}</p>`)
    .join('')
})

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
