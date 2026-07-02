<template>
  <div>
    <!-- Mini Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">NCS Uganda</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">Sports <span class="text-accent">Facilities</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">Explore the sports facilities and venues managed by the National Council of Sports Uganda.</p>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Facilities Grid -->
    <div class="bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div v-if="loading" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div v-for="i in 6" :key="i" class="animate-pulse">
            <div class="bg-gray-200 rounded-2xl h-48 mb-3"/>
            <div class="h-4 bg-gray-200 rounded w-2/3 mb-2"/>
            <div class="h-3 bg-gray-200 rounded w-1/2"/>
          </div>
        </div>

        <div v-else-if="facilities.length" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <router-link
            v-for="f in facilities"
            :key="f.id"
            :to="`/facilities/${f.slug}`"
            class="group bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-lg border border-gray-100 hover:border-[#F48C06]/30 transition-all"
          >
            <div class="h-48 overflow-hidden bg-gray-100 relative">
              <img
                v-if="f.image_url"
                :src="mediaUrl(f.image_url)"
                :alt="f.name"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
              >
              <div v-else class="w-full h-full flex items-center justify-center bg-gradient-to-br from-[#112b4e] to-[#1e4080]">
                <i class="icofont-football text-5xl text-white/30"></i>
              </div>
            </div>
            <div class="p-5">
              <h3 class="font-bold text-darken text-base group-hover:text-accent transition-colors">{{ f.name }}</h3>
              <p v-if="f.description" class="text-sm text-gray-500 line-clamp-2 mt-1.5">{{ f.description }}</p>
              <div class="mt-3 inline-flex items-center gap-1.5 text-sm font-semibold text-accent">
                Learn more <i class="icofont-arrow-right text-xs"></i>
              </div>
            </div>
          </router-link>
        </div>

        <div v-else class="text-center py-20 text-gray-400">
          <i class="icofont-football text-5xl mb-3 opacity-30"></i>
          <p>No facilities listed yet.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listFacilities } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

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
