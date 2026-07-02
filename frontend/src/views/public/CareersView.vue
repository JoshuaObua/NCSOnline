<template>
  <div>
    <!-- Mini Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">Join Our Team</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">Career <span class="text-accent">Opportunities</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">Explore job opportunities at NCS Uganda and help shape the future of sports in the country.</p>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Job Listings -->
    <div class="bg-white">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div v-if="loading" class="space-y-4">
          <div v-for="i in 4" :key="i" class="animate-pulse border rounded-2xl p-6">
            <div class="h-5 bg-gray-200 rounded mb-2 w-1/2"></div>
            <div class="h-3 bg-gray-200 rounded w-1/3"></div>
          </div>
        </div>

        <div v-else-if="careers.length" class="space-y-4">
          <router-link
            v-for="c in careers"
            :key="c.id"
            :to="`/careers/${c.id}`"
            class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl border border-gray-100 bg-white hover:border-[#F48C06]/40 hover:shadow-md transition-all group"
          >
            <div class="flex items-start gap-4">
              <!-- Dept icon circle -->
              <div class="w-11 h-11 rounded-2xl bg-[#112b4e]/10 flex items-center justify-center flex-shrink-0 group-hover:bg-[#112b4e] transition-colors">
                <svg class="w-5 h-5 text-[#112b4e] group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                </svg>
              </div>
              <div>
                <h2 class="font-bold text-darken text-base mb-1.5 group-hover:text-accent transition-colors">{{ c.title }}</h2>
                <div class="flex flex-wrap gap-2">
                  <span v-if="c.department" class="text-xs bg-[#112b4e]/8 text-[#112b4e] px-2.5 py-0.5 rounded-full font-medium">{{ c.department }}</span>
                  <span v-if="c.location" class="text-xs bg-gray-100 text-gray-600 px-2.5 py-0.5 rounded-full">📍 {{ c.location }}</span>
                  <span class="text-xs bg-gray-100 text-gray-600 px-2.5 py-0.5 rounded-full capitalize">{{ c.job_type?.replace('_', ' ') }}</span>
                  <span v-if="c.salary_range" class="text-xs bg-green-50 text-green-700 px-2.5 py-0.5 rounded-full">{{ c.salary_range }}</span>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-3 flex-shrink-0 md:ml-4">
              <span v-if="c.deadline_at" class="text-xs text-red-500 whitespace-nowrap">
                Closes {{ formatDate(c.deadline_at) }}
              </span>
              <span class="bg-[#F48C06] text-white text-xs font-semibold px-4 py-1.5 rounded-full whitespace-nowrap">View Details</span>
            </div>
          </router-link>
        </div>

        <div v-else class="text-center py-20 text-gray-400">
          <svg class="w-16 h-16 mx-auto mb-4 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
          </svg>
          <p>No open positions at this time. Please check back later.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listCareers } from '@/api/cms.js'

const careers = ref([])
const loading = ref(true)

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { day: 'numeric', month: 'short', year: 'numeric' })
}

onMounted(async () => {
  try {
    const res = await listCareers({ status: 'published', per_page: 20 })
    careers.value = res.data.data?.items || []
  } catch {
    careers.value = []
  } finally {
    loading.value = false
  }
})
</script>
