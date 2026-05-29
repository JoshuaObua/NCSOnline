<template>
  <div>
    <!-- Header -->
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-primary-400 font-semibold text-sm uppercase tracking-wide mb-2">NCS Uganda</p>
        <h1 class="text-4xl font-bold text-white mb-4">Sports Associations & Federations</h1>
        <p class="text-gray-300 max-w-2xl">Browse the sports associations and federations affiliated to the National Council of Sports.</p>
      </div>
    </div>

    <!-- Associations Grid -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div v-if="loading" class="grid sm:grid-cols-2 lg:grid-cols-3 gap-5">
        <div v-for="i in 6" :key="i" class="h-40 bg-gray-100 rounded-xl animate-pulse"/>
      </div>

      <div v-else-if="associations.length" class="grid sm:grid-cols-2 lg:grid-cols-3 gap-5">
        <div
          v-for="a in associations"
          :key="a.id"
          class="bg-white rounded-xl p-5 shadow-sm hover:shadow-md border border-gray-100 transition-all flex flex-col"
        >
          <div class="flex items-start gap-4 mb-3">
            <div class="flex-shrink-0 w-16 h-16 bg-gray-50 rounded-xl flex items-center justify-center overflow-hidden p-1.5">
              <img v-if="a.logo_url" :src="mediaUrl(a.logo_url)" :alt="a.name" class="max-w-full max-h-full object-contain">
              <i v-else class="icofont-trophy text-3xl text-primary-400"></i>
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="font-semibold text-gray-900 leading-snug">{{ a.name }}</h3>
            </div>
          </div>
          <p v-if="a.description" class="text-sm text-gray-500 line-clamp-3 mb-4 flex-1">{{ a.description }}</p>
          <a
            v-if="a.website_url"
            :href="a.website_url"
            target="_blank"
            rel="noopener"
            class="inline-flex items-center gap-1.5 text-sm font-medium text-primary-600 hover:text-primary-700 mt-auto"
          >
            Visit website <i class="icofont-external-link"></i>
          </a>
        </div>
      </div>

      <div v-else class="text-center py-20 text-gray-400">
        <i class="icofont-trophy text-5xl mb-3 opacity-40"></i>
        <p>No associations listed yet.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listAssociations } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const associations = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const r = await listAssociations()
    associations.value = r.data.data || []
  } catch { associations.value = [] }
  finally { loading.value = false }
})
</script>
