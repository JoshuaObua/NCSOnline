<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div v-if="loading" class="animate-pulse space-y-4">
      <div class="h-8 bg-gray-200 rounded w-2/3"></div>
      <div class="h-4 bg-gray-200 rounded w-1/2"></div>
      <div class="h-64 bg-gray-200 rounded-xl"></div>
    </div>

    <div v-else-if="career">
      <nav class="flex items-center gap-2 text-sm text-gray-400 mb-8">
        <router-link to="/" class="hover:text-primary-600">Home</router-link>
        <span>/</span>
        <router-link to="/careers" class="hover:text-primary-600">Careers</router-link>
        <span>/</span>
        <span class="text-gray-600 truncate max-w-xs">{{ career.title }}</span>
      </nav>

      <h1 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">{{ career.title }}</h1>

      <!-- Meta pills -->
      <div class="flex flex-wrap gap-3 mb-8">
        <span v-if="career.department" class="bg-gray-100 text-gray-600 text-sm px-3 py-1 rounded-full">🏢 {{ career.department }}</span>
        <span v-if="career.location" class="bg-gray-100 text-gray-600 text-sm px-3 py-1 rounded-full">📍 {{ career.location }}</span>
        <span class="bg-primary-50 text-primary-700 text-sm px-3 py-1 rounded-full capitalize">⏱ {{ career.job_type.replace('_', ' ') }}</span>
        <span v-if="career.salary_range" class="bg-gray-100 text-gray-600 text-sm px-3 py-1 rounded-full">💰 {{ career.salary_range }}</span>
        <span v-if="career.deadline_at" class="bg-red-50 text-red-600 text-sm px-3 py-1 rounded-full">⏰ Closes {{ formatDate(career.deadline_at) }}</span>
      </div>

      <!-- Description -->
      <div class="bg-white rounded-xl border border-gray-100 p-6 mb-6">
        <h2 class="font-semibold text-gray-900 mb-3 text-lg">About the Role</h2>
        <div class="text-gray-600 whitespace-pre-line text-sm leading-relaxed">{{ career.description }}</div>
      </div>

      <!-- Requirements -->
      <div v-if="career.requirements" class="bg-white rounded-xl border border-gray-100 p-6 mb-8">
        <h2 class="font-semibold text-gray-900 mb-3 text-lg">Requirements</h2>
        <div class="text-gray-600 whitespace-pre-line text-sm leading-relaxed">{{ career.requirements }}</div>
      </div>

      <!-- Apply CTA -->
      <div class="bg-primary-50 border border-primary-200 rounded-xl p-6 text-center">
        <h3 class="font-semibold text-gray-900 mb-2">Interested in this position?</h3>
        <p class="text-sm text-gray-500 mb-4">Submit your application through our online portal.</p>
        <router-link
          to="/apply"
          class="inline-block bg-primary-600 hover:bg-primary-700 text-white font-semibold px-6 py-2.5 rounded-lg transition-colors"
        >
          Apply Now
        </router-link>
      </div>

      <div class="mt-8">
        <router-link to="/careers" class="text-primary-600 hover:text-primary-700 font-medium text-sm">← All Careers</router-link>
      </div>
    </div>

    <div v-else class="text-center py-20 text-gray-400">
      <p>Job posting not found.</p>
      <router-link to="/careers" class="mt-4 inline-block text-primary-600 hover:underline">View all careers</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getCareer } from '@/api/cms.js'

const route = useRoute()
const career = ref(null)
const loading = ref(true)

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(async () => {
  try {
    const res = await getCareer(route.params.id)
    career.value = res.data.data
  } catch {
    career.value = null
  } finally {
    loading.value = false
  }
})
</script>
