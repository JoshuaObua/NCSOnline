<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="max-w-5xl mx-auto px-4 py-16">
      <div class="animate-pulse space-y-4">
        <div class="h-72 bg-gray-200 rounded-2xl"/>
        <div class="h-8 bg-gray-200 rounded w-1/2"/>
        <div class="h-4 bg-gray-200 rounded w-2/3"/>
        <div class="h-4 bg-gray-200 rounded w-3/4"/>
      </div>
    </div>

    <!-- Not found -->
    <div v-else-if="!facility" class="text-center py-24">
      <i class="icofont-search text-5xl text-gray-300 mb-3"></i>
      <p class="text-gray-500 mb-4">Facility not found.</p>
      <router-link to="/facilities" class="text-primary-600 hover:underline">← Back to all facilities</router-link>
    </div>

    <!-- Content -->
    <div v-else>
      <div class="relative h-72 md:h-96 bg-gray-200 overflow-hidden">
        <img v-if="facility.image_url" :src="mediaUrl(facility.image_url)" :alt="facility.name" class="w-full h-full object-cover">
        <div v-else class="w-full h-full bg-gradient-to-br from-primary-700 to-primary-500"/>
        <div class="absolute inset-0 bg-black/40 flex items-end">
          <div class="max-w-5xl mx-auto w-full px-4 sm:px-6 lg:px-8 pb-8">
            <router-link to="/facilities" class="text-primary-200 hover:text-white text-sm flex items-center gap-1.5 mb-3">
              <i class="icofont-arrow-left"></i> All Facilities
            </router-link>
            <h1 class="text-3xl md:text-4xl font-bold text-white">{{ facility.name }}</h1>
          </div>
        </div>
      </div>

      <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div class="prose prose-lg max-w-none text-gray-700">
          <p v-if="facility.description" class="text-lg leading-relaxed">{{ facility.description }}</p>
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
