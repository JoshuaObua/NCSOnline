<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="bg-cream pt-24 pb-16">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 animate-pulse">
        <div class="h-72 bg-gray-300 rounded-2xl mb-6"/>
        <div class="h-8 bg-gray-300 rounded w-1/2 mb-3"/>
        <div class="h-4 bg-gray-300 rounded w-2/3"/>
      </div>
    </div>

    <!-- Not found -->
    <div v-else-if="!facility" class="bg-cream pt-32 pb-20 text-center text-gray-400">
      <i class="icofont-search text-5xl mb-3 opacity-30"></i>
      <p class="text-gray-500 mb-4 font-medium">Facility not found.</p>
      <router-link to="/facilities" class="text-sm font-semibold text-accent hover:text-[#d47b05]">← Back to all facilities</router-link>
    </div>

    <!-- Content -->
    <div v-else>
      <!-- Hero Image Banner -->
      <div class="relative h-72 md:h-[420px] overflow-hidden">
        <img
          v-if="facility.image_url"
          :src="mediaUrl(facility.image_url)"
          :alt="facility.name"
          class="w-full h-full object-cover"
        >
        <div v-else class="w-full h-full bg-gradient-to-br from-[#112b4e] to-[#1a3a6a]"/>
        <div class="absolute inset-0 bg-gradient-to-t from-black/75 via-black/30 to-transparent"/>
        <div class="absolute inset-0 flex items-end">
          <div class="max-w-5xl mx-auto w-full px-4 sm:px-6 lg:px-8 pb-10">
            <router-link
              to="/facilities"
              class="inline-flex items-center gap-1.5 text-white/70 hover:text-white text-sm mb-4 transition-colors"
            >
              <i class="icofont-arrow-left text-xs"></i> All Facilities
            </router-link>
            <h1 class="text-3xl md:text-4xl font-bold text-white leading-tight">{{ facility.name }}</h1>
          </div>
        </div>
      </div>

      <!-- Description -->
      <div class="bg-white">
        <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
          <div class="max-w-3xl">
            <p
              v-if="facility.description"
              class="text-lg text-gray-600 leading-relaxed"
            >{{ facility.description }}</p>
            <p v-else class="text-gray-400 italic">No additional details available.</p>
          </div>

          <div class="mt-10 pt-8 border-t border-gray-100">
            <router-link
              to="/facilities"
              class="text-sm font-semibold text-[#112b4e] hover:text-[#F48C06] transition-colors"
            >← All Facilities</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { listFacilities } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const route = useRoute()
const facility = ref(null)
const loading = ref(true)

async function load() {
  loading.value = true
  try {
    const r = await listFacilities()
    const items = r.data.data || []
    facility.value = items.find(f => f.slug === route.params.slug) || null
  } catch { facility.value = null }
  finally { loading.value = false }
}

onMounted(load)
watch(() => route.params.slug, load)
</script>
