<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div v-if="loading" class="animate-pulse space-y-4">
      <div class="h-8 bg-gray-200 rounded w-3/4"></div>
      <div class="h-4 bg-gray-200 rounded w-1/3"></div>
      <div class="h-64 bg-gray-200 rounded-xl"></div>
    </div>

    <div v-else-if="event">
      <nav class="flex items-center gap-2 text-sm text-gray-400 mb-8">
        <router-link to="/" class="hover:text-primary-600">Home</router-link>
        <span>/</span>
        <router-link to="/events" class="hover:text-primary-600">Events</router-link>
        <span>/</span>
        <span class="text-gray-600 truncate max-w-xs">{{ event.title }}</span>
      </nav>

      <div class="flex items-start gap-6 mb-8">
        <div class="flex-shrink-0 w-16 h-16 bg-primary-600 text-white rounded-xl flex flex-col items-center justify-center">
          <span class="text-xl font-bold leading-none">{{ event.event_date ? new Date(event.event_date).getDate() : '?' }}</span>
          <span class="text-xs uppercase opacity-80">{{ event.event_date ? monthShort(event.event_date) : '' }}</span>
        </div>
        <div>
          <h1 class="text-3xl md:text-4xl font-bold text-gray-900 mb-2">{{ event.title }}</h1>
          <div class="flex flex-wrap gap-4 text-sm text-gray-500">
            <span v-if="event.event_date">📅 {{ formatDate(event.event_date) }}</span>
            <span v-if="event.end_date">– {{ formatDate(event.end_date) }}</span>
            <span v-if="event.location">📍 {{ event.location }}</span>
          </div>
        </div>
      </div>

      <img v-if="event.cover_image_url" :src="mediaUrl(event.cover_image_url)" :alt="event.title" class="w-full rounded-xl mb-10 max-h-80 object-cover">

      <div class="prose prose-gray max-w-none" v-html="formattedDescription"></div>

      <div class="mt-12 pt-8 border-t border-gray-100">
        <router-link to="/events" class="text-primary-600 hover:text-primary-700 font-medium text-sm">← All Events</router-link>
      </div>
    </div>

    <div v-else class="text-center py-20 text-gray-400">
      <p>Event not found.</p>
      <router-link to="/events" class="mt-4 inline-block text-primary-600 hover:underline">View all events</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getEvent } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const route = useRoute()
const event = ref(null)
const loading = ref(true)

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}
function monthShort(d) {
  return new Date(d).toLocaleDateString('en-UG', { month: 'short' })
}

const formattedDescription = computed(() => {
  if (!event.value?.description) return ''
  return event.value.description.split('\n\n').map(p => `<p>${p.replace(/\n/g, '<br>')}</p>`).join('')
})

onMounted(async () => {
  try {
    const res = await getEvent(route.params.slug)
    event.value = res.data.data
  } catch {
    event.value = null
  } finally {
    loading.value = false
  }
})
</script>
