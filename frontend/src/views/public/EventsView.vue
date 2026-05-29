<template>
  <div>
    <!-- Mini Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">NCS Uganda</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">Sports <span class="text-accent">Events</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">Competitions, sports activities and events organised or sanctioned by the National Council of Sports Uganda.</p>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Events List -->
    <div class="bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div v-if="loading" class="grid md:grid-cols-2 gap-5">
          <div v-for="i in 4" :key="i" class="animate-pulse flex gap-4 p-5 border rounded-2xl">
            <div class="w-16 h-16 bg-gray-200 rounded-2xl flex-shrink-0"></div>
            <div class="flex-1">
              <div class="h-4 bg-gray-200 rounded mb-2"></div>
              <div class="h-3 bg-gray-200 rounded w-2/3"></div>
            </div>
          </div>
        </div>

        <div v-else-if="events.length" class="grid md:grid-cols-2 gap-5">
          <router-link
            v-for="event in events"
            :key="event.id"
            :to="`/events/${event.slug}`"
            class="flex gap-4 p-5 rounded-2xl border border-gray-100 bg-white hover:border-[#F48C06]/40 hover:shadow-md transition-all group"
          >
            <!-- Date badge -->
            <div class="flex-shrink-0 w-16 h-16 bg-[#112b4e] text-white rounded-2xl flex flex-col items-center justify-center text-center shadow-sm">
              <span class="text-xl font-bold leading-none">{{ event.event_date ? new Date(event.event_date).getDate() : '?' }}</span>
              <span class="text-xs uppercase opacity-70 tracking-wide">{{ event.event_date ? monthShort(event.event_date) : '' }}</span>
            </div>

            <div class="flex-1 min-w-0">
              <h2 class="font-bold text-darken mb-1 group-hover:text-accent transition-colors">{{ event.title }}</h2>
              <div class="flex items-center gap-1 text-sm text-gray-500 mb-1" v-if="event.location">
                <svg class="w-3.5 h-3.5 text-accent flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd"/>
                </svg>
                {{ event.location }}
              </div>
              <p v-if="event.event_date" class="text-xs text-gray-400">{{ formatDate(event.event_date) }}</p>
            </div>

            <!-- Arrow -->
            <div class="flex-shrink-0 self-center text-gray-300 group-hover:text-accent transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
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
