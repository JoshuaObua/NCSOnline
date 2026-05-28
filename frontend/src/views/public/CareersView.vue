<template>
  <div>
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-primary-400 font-semibold text-sm uppercase tracking-wide mb-2">Join Our Team</p>
        <h1 class="text-4xl font-bold text-white mb-4">Career Opportunities</h1>
        <p class="text-gray-300 max-w-2xl">Explore job opportunities at the National Council of Sports Uganda and help shape the future of sports in the country.</p>
      </div>
    </div>

    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div v-if="loading" class="space-y-4">
        <div v-for="i in 4" :key="i" class="animate-pulse border rounded-xl p-6">
          <div class="h-5 bg-gray-200 rounded mb-2 w-1/2"></div>
          <div class="h-3 bg-gray-200 rounded w-1/3"></div>
        </div>
      </div>

      <div v-else-if="careers.length" class="space-y-4">
        <router-link
          v-for="c in careers"
          :key="c.id"
          :to="`/careers/${c.id}`"
          class="block p-6 rounded-xl border border-gray-100 hover:border-primary-200 hover:shadow-sm transition-all bg-white"
        >
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
            <div>
              <h2 class="font-semibold text-gray-900 text-lg mb-1">{{ c.title }}</h2>
              <div class="flex flex-wrap gap-3 text-sm text-gray-500">
                <span v-if="c.department">🏢 {{ c.department }}</span>
                <span v-if="c.location">📍 {{ c.location }}</span>
                <span class="capitalize">⏱ {{ c.job_type.replace('_', ' ') }}</span>
                <span v-if="c.salary_range">💰 {{ c.salary_range }}</span>
              </div>
            </div>
            <div class="flex items-center gap-3">
              <span v-if="c.deadline_at" class="text-xs text-red-500 whitespace-nowrap">
                Closes: {{ formatDate(c.deadline_at) }}
              </span>
              <span class="bg-primary-50 text-primary-700 text-xs font-medium px-3 py-1.5 rounded-full whitespace-nowrap">Apply Now</span>
            </div>
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
