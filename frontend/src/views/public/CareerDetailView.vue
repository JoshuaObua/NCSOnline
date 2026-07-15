<template>
  <div>
    <!-- Loading -->
    <div v-if="loading" class="bg-cream pt-24 pb-16">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 animate-pulse">
        <div class="h-10 bg-gray-300 rounded mb-3 w-2/3"></div>
        <div class="h-4 bg-gray-300 rounded w-1/3"></div>
      </div>
    </div>

    <div v-else-if="career">
      <!-- Banner -->
      <div class="bg-[#112b4e] pt-28 pb-12 px-4">
        <div class="max-w-5xl mx-auto sm:px-6">
          <nav class="flex items-center gap-2 text-sm text-white/50 mb-5">
            <router-link to="/" class="hover:text-white">Home</router-link>
            <span>/</span>
            <router-link to="/careers" class="hover:text-white">Careers</router-link>
            <span>/</span>
            <span class="text-white/80 truncate max-w-xs">{{ career.title }}</span>
          </nav>
          <h1 class="text-3xl md:text-4xl font-bold text-white mb-4">{{ career.title }}</h1>
          <div class="flex flex-wrap gap-2">
            <span v-if="career.department" class="text-xs bg-white/10 text-white/80 px-3 py-1 rounded-full">{{ career.department }}</span>
            <span v-if="career.location" class="text-xs bg-white/10 text-white/80 px-3 py-1 rounded-full">📍 {{ career.location }}</span>
            <span class="text-xs bg-[#F48C06]/20 text-[#F48C06] px-3 py-1 rounded-full capitalize font-medium">{{ career.job_type?.replace('_', ' ') }}</span>
            <span v-if="career.deadline_at" class="text-xs bg-red-500/20 text-red-300 px-3 py-1 rounded-full">Closes {{ formatDate(career.deadline_at) }}</span>
          </div>
        </div>
      </div>

      <!-- Content + Sidebar -->
      <div class="bg-white">
        <div class="max-w-5xl mx-auto px-4 sm:px-6 py-12">
          <div class="grid lg:grid-cols-3 gap-10">
            <!-- Main content -->
            <div class="lg:col-span-2 space-y-6">
              <div class="bg-white rounded-2xl border border-gray-100 p-6">
                <h2 class="font-bold text-darken text-lg mb-4">About the Role</h2>
                <div class="text-gray-600 whitespace-pre-line text-sm leading-relaxed">{{ career.description }}</div>
              </div>

              <div v-if="career.requirements" class="bg-white rounded-2xl border border-gray-100 p-6">
                <h2 class="font-bold text-darken text-lg mb-4">Requirements</h2>
                <div class="text-gray-600 whitespace-pre-line text-sm leading-relaxed">{{ career.requirements }}</div>
              </div>
            </div>

            <!-- Sticky sidebar -->
            <div class="lg:col-span-1">
              <div class="bg-[#FEF9F2] rounded-2xl p-6 border border-[#F48C06]/20 sticky top-24 space-y-4">
                <h3 class="font-bold text-darken text-base">Job Details</h3>

                <div v-if="career.salary_range" class="text-sm">
                  <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1">Salary Range</p>
                  <p class="font-semibold text-darken">{{ career.salary_range }}</p>
                </div>

                <div v-if="career.deadline_at" class="text-sm">
                  <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1">Application Deadline</p>
                  <p class="font-semibold text-red-600">{{ formatDate(career.deadline_at) }}</p>
                </div>

                <div v-if="career.job_type" class="text-sm">
                  <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1">Job Type</p>
                  <p class="font-medium text-darken capitalize">{{ career.job_type?.replace('_', ' ') }}</p>
                </div>

                <div class="pt-4">
                  <a
                    :href="portalUrl('/login')"
                    class="w-full flex items-center justify-center gap-2 bg-[#F48C06] hover:bg-[#d47b05] text-white font-semibold px-5 py-3 rounded-xl text-sm transition-colors"
                  >
                    Apply for this role →
                  </a>
                </div>

                <div class="pt-2 border-t border-[#F48C06]/20">
                  <router-link
                    to="/careers"
                    class="text-sm font-medium text-[#112b4e] hover:text-[#F48C06] transition-colors"
                  >
                    ← All Careers
                  </router-link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="bg-cream pt-32 pb-20 text-center text-gray-400">
      <p class="text-lg font-medium mb-4">Job posting not found.</p>
      <router-link to="/careers" class="text-sm font-semibold text-accent hover:text-[#d47b05]">← View all careers</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getCareer } from '@/api/cms.js'
import { portalUrl } from '@/utils/portal.js'

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
