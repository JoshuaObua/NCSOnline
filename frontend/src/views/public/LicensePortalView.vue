<template>
  <div>
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-primary-400 font-semibold text-sm uppercase tracking-wide mb-2">Online Services</p>
        <h1 class="text-4xl font-bold text-white mb-4">License Portal</h1>
        <p class="text-gray-300 max-w-2xl">Apply for a new sports licence or renew an existing one through our secure online portal.</p>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <!-- If not logged in -->
      <div v-if="!isAuthenticated" class="max-w-xl mx-auto text-center">
        <div class="bg-primary-50 border border-primary-200 rounded-2xl p-8 mb-8">
          <div class="w-16 h-16 bg-primary-600 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
            </svg>
          </div>
          <h2 class="text-xl font-bold text-gray-900 mb-3">Sign In Required</h2>
          <p class="text-gray-600 mb-6">You need to create an account or sign in to apply for or manage your sports licence.</p>
          <div class="flex flex-col sm:flex-row gap-3 justify-center">
            <router-link
              to="/login"
              class="bg-primary-600 hover:bg-primary-700 text-white font-semibold px-6 py-2.5 rounded-lg transition-colors"
            >
              Sign In
            </router-link>
          </div>
        </div>
      </div>

      <!-- If logged in — show application options -->
      <div v-else>
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
          <div
            v-for="type in licenseTypes"
            :key="type.id"
            class="bg-white border border-gray-100 rounded-xl p-6 hover:border-primary-200 hover:shadow-sm transition-all cursor-pointer"
            @click="selectType(type)"
          >
            <div class="w-12 h-12 rounded-xl mb-4 flex items-center justify-center" :class="type.color">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="type.icon"/>
              </svg>
            </div>
            <h3 class="font-semibold text-gray-900 mb-2">{{ type.title }}</h3>
            <p class="text-sm text-gray-500 mb-4">{{ type.description }}</p>
            <span class="text-sm text-primary-600 font-medium">Apply now →</span>
          </div>
        </div>

        <!-- My Applications quick link -->
        <div class="bg-gray-50 rounded-xl p-6 border border-gray-100">
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div>
              <h3 class="font-semibold text-gray-900 mb-1">Track Your Application</h3>
              <p class="text-sm text-gray-500">View the status of all your submitted applications.</p>
            </div>
            <router-link
              to="/applications"
              class="bg-white border border-gray-200 hover:border-primary-300 text-gray-700 font-medium px-5 py-2.5 rounded-lg transition-colors text-sm whitespace-nowrap"
            >
              My Applications
            </router-link>
          </div>
        </div>
      </div>

      <!-- Application type info -->
      <div class="mt-16">
        <h2 class="text-2xl font-bold text-gray-900 mb-8">Types of Licences</h2>
        <div class="space-y-4">
          <div v-for="info in licenseInfo" :key="info.title" class="bg-white border border-gray-100 rounded-xl p-6">
            <div class="flex items-start gap-4">
              <div class="w-8 h-8 bg-primary-100 text-primary-700 rounded-lg flex items-center justify-center flex-shrink-0 text-sm font-bold">
                {{ info.no }}
              </div>
              <div>
                <h3 class="font-semibold text-gray-900 mb-1">{{ info.title }}</h3>
                <p class="text-sm text-gray-500">{{ info.desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'

const router = useRouter()
const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)

const licenseTypes = [
  {
    id: 'national_federation',
    title: 'National Federation',
    description: 'Register and license a national sports federation to organise competitions at a national level.',
    color: 'bg-primary-600',
    icon: 'M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064',
  },
  {
    id: 'community_club',
    title: 'Community / Club',
    description: 'Register a community-level sports club or association.',
    color: 'bg-gray-700',
    icon: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0',
  },
  {
    id: 'renewal',
    title: 'Renew Existing Licence',
    description: 'Renew your sports organisation licence for the current season.',
    color: 'bg-blue-600',
    icon: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15',
  },
]

const licenseInfo = [
  { no: '1', title: 'National Sports Federation Licence', desc: 'Required for bodies that oversee a specific sport at the national level and represent Uganda in international competitions.' },
  { no: '2', title: 'Regional/District Sports Association Licence', desc: 'For organisations that oversee a sport within a specific region or district.' },
  { no: '3', title: 'Community Club Registration', desc: 'For grassroots-level clubs, school sports clubs, and community associations.' },
  { no: '4', title: 'Sports Facility Licence', desc: 'Required for commercial sports facilities, stadiums and training grounds.' },
]

function selectType(type) {
  router.push('/applications')
}
</script>
