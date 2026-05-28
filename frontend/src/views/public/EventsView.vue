<template>
  <div>
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-primary-400 font-semibold text-sm uppercase tracking-wide mb-2">NCS Uganda</p>
        <h1 class="text-4xl font-bold text-white mb-4">Events</h1>
        <p class="text-gray-300 max-w-2xl">Sports events, competitions and activities organised or sanctioned by the National Council of Sports Uganda.</p>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div v-if="loading" class="grid md:grid-cols-2 gap-6">
        <div v-for="i in 4" :key="i" class="animate-pulse flex gap-4 p-5 border rounded-xl">
          <div class="w-14 h-14 bg-gray-200 rounded-xl flex-shrink-0"></div>
          <div class="flex-1">
            <div class="h-4 bg-gray-200 rounded mb-2"></div>
            <div class="h-3 bg-gray-200 rounded w-2/3"></div>
          </div>
        </div>
      </div>

      <div v-else-if="events.length" class="grid md:grid-cols-2 gap-6">
        <router-link
          v-for="event in events"
          :key="event.id"
          :to="`/events/${event.slug}`"
          class="flex gap-4 p-5 rounded-xl border border-gray-100 hover:border-primary-200 hover:bg-primary-50/30 hover:shadow-sm transition-all"
        >
          <div class="flex-shrink-0 w-14 h-14 bg-primary-600 text-white rounded-xl flex flex-col items-center justify-center text-center">
            <span class="text-lg font-bold leading-none">{{ event.event_date ? new Date(event.event_date).getDate() : '?' }}</span>
            <span class="text-xs uppercase opacity-75">{{ event.event_date ? monthShort(event.event_date) : '' }}</span>
          </div>
          <div class="flex-1 min-w-0">
            <h2 class="font-semibold text-gray-900 mb-1">{{ event.title }}</h2>
            <p v-if="event.location" class="text-sm text-gray-500 mb-1">📍 {{ event.location }}</p>
            <p v-if="event.event_date" class="text-xs text-gray-400">{{ formatDate(event.event_date) }}</p>
          </div>
        </router-link>
      </div>

      <div v-else class="text-center py-20 text-gray-400">
        <svg class="w-16 h-16 mx-auto mb-4 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
        </svg>
        <p>No upcoming events at this time.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listEvents } from '@/api/cms.js'

const events = ref([])
const loading = ref(true)

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}
function monthShort(d) {
  return new Date(d).toLocaleDateString('en-UG', { month: 'short' })
}

onMounted(async () => {
  try {
    const res = await listEvents({ status: 'published', per_page: 20 })
    events.value = res.data.data?.items || []
  } catch {
    events.value = []
  } finally {
    loading.value = false
  }
})
</script>
