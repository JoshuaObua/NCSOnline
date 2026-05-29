<template>
  <div>
    <!-- Mini Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">NCS Uganda</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">Associations &amp; <span class="text-accent">Federations</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">Browse sports associations and federations affiliated to the National Council of Sports Uganda.</p>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Associations Grid -->
    <div class="bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div v-if="loading" class="grid sm:grid-cols-2 lg:grid-cols-3 gap-5">
          <div v-for="i in 6" :key="i" class="h-44 bg-gray-100 rounded-2xl animate-pulse"/>
        </div>

        <div v-else-if="associations.length" class="grid sm:grid-cols-2 lg:grid-cols-3 gap-5">
          <div
            v-for="a in associations"
            :key="a.id"
            class="bg-white rounded-2xl p-5 shadow-sm hover:shadow-md border border-gray-100 hover:border-[#F48C06]/30 transition-all flex flex-col group"
          >
            <div class="flex items-start gap-4 mb-3">
              <div class="flex-shrink-0 w-16 h-16 bg-[#FEF9F2] border border-[#F48C06]/20 rounded-2xl flex items-center justify-center overflow-hidden p-2">
                <img v-if="a.logo_url" :src="mediaUrl(a.logo_url)" :alt="a.name" class="max-w-full max-h-full object-contain">
                <i v-else class="icofont-trophy text-3xl text-[#F48C06]"></i>
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="font-bold text-darken leading-snug group-hover:text-accent transition-colors">{{ a.name }}</h3>
              </div>
            </div>
            <p v-if="a.description" class="text-sm text-gray-500 line-clamp-3 mb-4 flex-1">{{ a.description }}</p>
            <a
              v-if="a.website_url"
              :href="a.website_url"
              target="_blank"
              rel="noopener"
              class="inline-flex items-center gap-1.5 text-sm font-semibold text-accent hover:text-[#d47b05] mt-auto transition-colors"
            >
              Visit website <i class="icofont-external-link text-xs"></i>
            </a>
          </div>
        </div>

        <div v-else class="text-center py-20 text-gray-400">
          <i class="icofont-trophy text-5xl mb-3 opacity-30"></i>
          <p>No associations listed yet.</p>
        </div>
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
