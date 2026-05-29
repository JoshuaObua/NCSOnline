<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="bg-cream pt-24 pb-16">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 animate-pulse">
        <div class="h-10 bg-gray-300 rounded mb-3 w-2/3"></div>
        <div class="h-4 bg-gray-300 rounded w-1/3"></div>
      </div>
    </div>

    <div v-else-if="event">
      <!-- Hero image with dark overlay + title -->
      <div class="relative min-h-[340px] flex items-end">
        <div
          v-if="event.cover_image_url"
          class="absolute inset-0 bg-cover bg-center"
          :style="`background-image: url('${mediaUrl(event.cover_image_url)}')`"
        ></div>
        <div v-else class="absolute inset-0 bg-gradient-to-br from-[#112b4e] to-[#1a3a6a]"></div>
        <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent"></div>

        <div class="relative z-10 max-w-5xl mx-auto px-4 sm:px-6 pb-10 pt-28 w-full">
          <nav class="flex items-center gap-2 text-sm text-white/60 mb-4">
            <router-link to="/" class="hover:text-white">Home</router-link>
            <span>/</span>
            <router-link to="/events" class="hover:text-white">Events</router-link>
            <span>/</span>
            <span class="text-white/80 truncate max-w-xs">{{ event.title }}</span>
          </nav>
          <h1 class="text-3xl md:text-4xl font-bold text-white leading-tight">{{ event.title }}</h1>
        </div>
      </div>

      <!-- Content + Sidebar -->
      <div class="bg-white">
        <div class="max-w-5xl mx-auto px-4 sm:px-6 py-12">
          <div class="grid lg:grid-cols-3 gap-10">
            <!-- Main Description -->
            <div class="lg:col-span-2">
              <div
                class="prose prose-gray max-w-none prose-headings:text-[#112b4e] prose-a:text-[#F48C06] prose-p:leading-relaxed"
                v-html="formattedDescription"
              ></div>
            </div>

            <!-- Info Sidebar -->
            <div class="lg:col-span-1">
              <div class="bg-[#FEF9F2] rounded-2xl p-6 border border-[#F48C06]/20 sticky top-24">
                <h3 class="font-bold text-darken mb-5 text-lg">Event Details</h3>

                <div class="space-y-4 text-sm">
                  <div v-if="event.event_date" class="flex items-start gap-3">
                    <div class="w-8 h-8 bg-[#112b4e] text-white rounded-lg flex items-center justify-center flex-shrink-0">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                      </svg>
                    </div>
                    <div>
                      <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-0.5">Date</p>
                      <p class="font-medium text-darken">{{ formatDate(event.event_date) }}</p>
                      <p v-if="event.end_date" class="text-gray-500 text-xs mt-0.5">to {{ formatDate(event.end_date) }}</p>
                    </div>
                  </div>

                  <div v-if="event.location" class="flex items-start gap-3">
                    <div class="w-8 h-8 bg-[#F48C06] text-white rounded-lg flex items-center justify-center flex-shrink-0">
                      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd"/>
                      </svg>
                    </div>
                    <div>
                      <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-0.5">Location</p>
                      <p class="font-medium text-darken">{{ event.location }}</p>
                    </div>
                  </div>

                  <!-- Date badge display -->
                  <div v-if="event.event_date" class="mt-2 flex items-center gap-3">
                    <div class="w-14 h-14 bg-[#112b4e] text-white rounded-2xl flex flex-col items-center justify-center shadow-sm">
                      <span class="text-xl font-bold leading-none">{{ new Date(event.event_date).getDate() }}</span>
                      <span class="text-[10px] uppercase opacity-70 tracking-wide">{{ monthShort(event.event_date) }}</span>
                    </div>
                    <span class="text-sm text-gray-500">{{ new Date(event.event_date).getFullYear() }}</span>
                  </div>
                </div>

                <div class="mt-6 pt-5 border-t border-[#F48C06]/20">
                  <router-link
                    to="/events"
                    class="flex items-center gap-2 text-sm font-semibold text-[#112b4e] hover:text-[#F48C06] transition-colors"
                  >
                    ← All Events
                  </router-link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="bg-cream pt-32 pb-20 text-center text-gray-400">
      <p class="text-lg font-medium mb-4">Event not found.</p>
      <router-link to="/events" class="text-sm font-semibold text-accent hover:text-[#d47b05]">← View all events</router-link>
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
  if (!event.value?.description) return '<p class="text-gray-400 italic">No description provided.</p>'
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
