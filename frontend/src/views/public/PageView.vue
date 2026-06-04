<template>
  <div>
    <div v-if="loading" class="bg-cream pt-24 pb-16">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 animate-pulse">
        <div class="h-12 bg-gray-300 rounded mb-3 w-3/4"></div>
        <div class="h-4 bg-gray-300 rounded w-1/3"></div>
      </div>
    </div>

    <div v-else-if="page">
      <!-- Hero / Title -->
      <div class="relative min-h-[280px] flex items-end">
        <div v-if="page.cover_image_url"
          class="absolute inset-0 bg-cover bg-center"
          :style="`background-image: url('${mediaUrl(page.cover_image_url)}')`"></div>
        <div v-else class="absolute inset-0 bg-gradient-to-br from-[#112b4e] to-[#1e4080]"></div>
        <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent"></div>
        <div class="relative z-10 max-w-4xl mx-auto px-4 sm:px-6 pb-10 pt-28 w-full">
          <nav class="flex items-center gap-2 text-sm text-white/60 mb-4">
            <router-link to="/" class="hover:text-white transition-colors">Home</router-link>
            <span>/</span>
            <span class="text-white/80 truncate max-w-xs">{{ page.title }}</span>
          </nav>
          <h1 class="text-3xl md:text-4xl font-bold text-white leading-tight">{{ page.title }}</h1>
        </div>
      </div>

      <!-- Content -->
      <div class="bg-white">
        <div class="max-w-4xl mx-auto px-4 sm:px-6 py-14">
          <p v-if="page.excerpt" class="text-lg text-gray-500 mb-8 border-l-4 border-[#F48C06] pl-5 italic">{{ page.excerpt }}</p>
          <div
            class="prose prose-gray max-w-none prose-headings:text-[#112b4e] prose-a:text-[#F48C06] prose-a:no-underline hover:prose-a:underline prose-img:rounded-xl prose-p:text-gray-600 prose-p:leading-relaxed"
            v-html="page.content"
          ></div>
        </div>
      </div>
    </div>

    <div v-else class="bg-cream pt-32 pb-20 text-center text-gray-400">
      <p class="text-lg font-medium mb-4">Page not found.</p>
      <router-link to="/" class="text-sm font-semibold text-accent hover:text-[#d47b05]">← Go to Home</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getPost } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const route = useRoute()
const page = ref(null)
const loading = ref(true)

async function load(slug) {
  loading.value = true
  try {
    const res = await getPost(slug)
    const post = res.data.data
    if (post?.category !== 'page') {
      page.value = null
    } else {
      page.value = post
      if (post.meta_title) document.title = post.meta_title
      else if (post.title) document.title = post.title + ' | NCS Uganda'
    }
  } catch {
    page.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => load(route.params.slug))
watch(() => route.params.slug, (s) => s && load(s))
</script>
