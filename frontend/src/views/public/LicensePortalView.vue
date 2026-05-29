<template>
  <div>
    <!-- Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">Online Services</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">Apply for a <span class="text-accent">Sports Licence</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">Apply for a new sports licence or renew an existing one through our secure online portal.</p>

        <!-- Step indicator -->
        <div class="flex items-center justify-center gap-0 mt-8 max-w-md mx-auto">
          <div v-for="(step, i) in steps" :key="i" class="flex items-center flex-1">
            <div class="flex flex-col items-center gap-1 flex-1">
              <div class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold transition-colors"
                :class="i === 0 ? 'bg-[#F48C06] text-white' : 'bg-gray-200 text-gray-500'">{{ i + 1 }}</div>
              <span class="text-xs text-gray-500 text-center leading-tight hidden sm:block">{{ step }}</span>
            </div>
            <div v-if="i < steps.length - 1" class="h-px bg-gray-200 flex-1 max-w-8 mb-4"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Content -->
    <div class="bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">

        <!-- If not logged in -->
        <div v-if="!isAuthenticated" class="max-w-xl mx-auto text-center">
          <div class="bg-[#FEF9F2] border border-[#F48C06]/20 rounded-2xl p-8 mb-8 shadow-sm">
            <div class="w-16 h-16 bg-[#112b4e] rounded-2xl flex items-center justify-center mx-auto mb-5">
              <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
              </svg>
            </div>
            <h2 class="text-xl font-bold text-darken mb-3">Sign In Required</h2>
            <p class="text-gray-500 mb-6">You need to create an account or sign in to apply for or manage your sports licence.</p>
            <div class="flex flex-col sm:flex-row gap-3 justify-center">
              <router-link
                to="/login"
                class="bg-[#F48C06] hover:bg-[#d47b05] text-white font-bold px-6 py-3 rounded-full transition-colors"
              >
                Sign In to Continue
              </router-link>
            </div>
          </div>
        </div>

        <!-- If logged in — show application options -->
        <div v-else>
          <h2 class="text-2xl font-bold text-darken mb-8 text-center">Select Licence Type</h2>
          <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
            <div
              v-for="type in licenseTypes"
              :key="type.id"
              class="bg-white border border-gray-100 rounded-2xl p-6 hover:border-[#F48C06]/40 hover:shadow-md transition-all cursor-pointer group"
              @click="selectType(type)"
            >
              <div class="w-12 h-12 rounded-2xl mb-4 flex items-center justify-center shadow-sm transition-transform group-hover:scale-110" :class="type.color">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="type.icon"/>
                </svg>
              </div>
              <h3 class="font-bold text-darken mb-2 group-hover:text-accent transition-colors">{{ type.title }}</h3>
              <p class="text-sm text-gray-500 mb-4">{{ type.description }}</p>
              <span class="text-sm font-semibold text-accent flex items-center gap-1">Apply now →</span>
            </div>
          </div>

          <!-- Track applications -->
          <div class="bg-[#FEF9F2] border border-[#F48C06]/20 rounded-2xl p-6">
            <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
              <div>
                <h3 class="font-bold text-darken mb-1">Track Your Application</h3>
                <p class="text-sm text-gray-500">View the status of all your submitted applications.</p>
              </div>
              <router-link
                to="/applications"
                class="bg-[#112b4e] hover:bg-[#0e2240] text-white font-semibold px-5 py-2.5 rounded-full transition-colors text-sm whitespace-nowrap"
              >
                My Applications
              </router-link>
            </div>
          </div>
        </div>

        <!-- Licence types info section -->
        <div class="mt-20">
          <div class="text-center mb-10">
            <p class="section-tag mb-2">Overview</p>
            <h2 class="text-2xl font-bold text-darken">Types of Licences</h2>
          </div>
          <div class="grid md:grid-cols-2 gap-5">
            <div
              v-for="info in licenseInfo"
              :key="info.title"
              class="bg-white border border-gray-100 rounded-2xl p-6 flex items-start gap-4 hover:shadow-sm transition-shadow"
            >
              <div class="w-10 h-10 bg-[#F48C06] text-white rounded-2xl flex items-center justify-center flex-shrink-0 font-bold text-sm shadow-sm">
                {{ info.no }}
              </div>
              <div>
                <h3 class="font-bold text-darken mb-1">{{ info.title }}</h3>
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

const steps = ['Select Type', 'Fill Form', 'Submit']

const licenseTypes = [
  {
    id: 'national_federation',
    title: 'National Federation',
    description: 'Register and license a national sports federation to organise competitions at a national level.',
    color: 'bg-[#112b4e]',
    icon: 'M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064',
  },
  {
    id: 'community_club',
    title: 'Community / Club',
    description: 'Register a community-level sports club or association.',
    color: 'bg-gray-600',
    icon: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0',
  },
  {
    id: 'renewal',
    title: 'Renew Existing Licence',
    description: 'Renew your sports organisation licence for the current season.',
    color: 'bg-[#F48C06]',
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
