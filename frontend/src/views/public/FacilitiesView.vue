<template>
  <div>
    <!-- Header -->
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-primary-400 font-semibold text-sm uppercase tracking-wide mb-2">NCS Uganda</p>
        <h1 class="text-4xl font-bold text-white mb-4">Sports Facilities</h1>
        <p class="text-gray-300 max-w-2xl">Explore the sports facilities and venues managed by the National Council of Sports Uganda.</p>
      </div>
    </div>

    <!-- Facilities Grid -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div v-if="loading" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="i in 6" :key="i" class="animate-pulse">
          <div class="bg-gray-200 rounded-xl h-48 mb-3"/>
          <div class="h-4 bg-gray-200 rounded w-2/3"/>
        </div>
      </div>

      <div v-else-if="facilities.length" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <router-link
          v-for="f in facilities"
          :key="f.id"
          :to="`/facilities/${f.slug}`"
          class="group bg-white rounded-xl overflow-hidden shadow-sm hover:shadow-lg border border-gray-100 transition-all"
        >
          <div class="h-48 overflow-hidden bg-gray-100">
            <img v-if="f.image_url" :src="f.image_url" :alt="f.name" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300">
            <div v-else class="w-full h-full flex items-center justify-center text-primary-200 text-5xl">
              <i class="icofont-football"></i>
            </div>
          </div>
          <div class="p-5">
            <h3 class="font-semibold text-lg text-gray-900 group-hover:text-primary-700 transition-colors">{{ f.name }}</h3>
            <p v-if="f.description" class="text-sm text-gray-500 line-clamp-2 mt-2">{{ f.description }}</p>
            <div class="mt-4 inline-flex items-center gap-1.5 text-sm text-primary-600 font-medium">
              Learn more <i class="icofont-arrow-right"></i>
            </div>
          </div>
        </router-link>
      </div>

      <div v-else class="text-center py-20 text-gray-400">
        <i class="icofont-football text-5xl mb-3 opacity-40"></i>
        <p>No facilities listed yet.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listFacilities } from '@/api/cms.js'

const facilities = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const r = await listFacilities()
    facilities.value = r.data.data || []
  } catch { facilities.value = [] }
  finally { loading.value = false }
})
</script>
