<template>
  <div class="min-h-screen flex flex-col bg-white">
    <!-- Navigation -->
    <header class="bg-white shadow-sm sticky top-0 z-50 border-b border-gray-100">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo -->
          <router-link to="/" class="flex items-center gap-3">
            <div class="w-12 h-12 bg-white rounded-xl flex items-center justify-center overflow-hidden p-0.5 shadow-sm border border-gray-100">
              <img src="/main-logo.png" alt="NCS Logo" class="w-full h-full object-contain" />
            </div>
            <div class="leading-tight">
              <div class="font-bold text-gray-900 text-sm">National Council of Sports</div>
              <div class="text-xs text-gray-500">Uganda</div>
            </div>
          </router-link>

          <!-- Desktop Navigation -->
          <nav class="hidden md:flex items-center gap-6 text-sm font-medium">
            <router-link to="/" class="text-gray-600 hover:text-primary-700 transition-colors">Home</router-link>
            <router-link to="/apply" class="text-gray-600 hover:text-primary-700 transition-colors">Apply</router-link>
            <router-link to="/news" class="text-gray-600 hover:text-primary-700 transition-colors">News</router-link>
            <router-link to="/events" class="text-gray-600 hover:text-primary-700 transition-colors">Events</router-link>
            <router-link to="/careers" class="text-gray-600 hover:text-primary-700 transition-colors">Careers</router-link>
            <router-link to="/projects" class="text-gray-600 hover:text-primary-700 transition-colors">Projects</router-link>
            <router-link to="/case-studies" class="text-gray-600 hover:text-primary-700 transition-colors">Case Studies</router-link>
          </nav>

          <!-- CTA -->
          <div class="flex items-center gap-3">
            <router-link
              v-if="!isAuthenticated"
              to="/login"
              class="text-sm font-medium text-gray-600 hover:text-primary-700 transition-colors"
            >
              Sign In
            </router-link>
            <router-link
              v-if="isAuthenticated"
              to="/dashboard"
              class="text-sm font-medium text-primary-700 hover:text-primary-800 transition-colors"
            >
              Dashboard
            </router-link>
            <router-link
              to="/apply"
              class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-semibold px-4 py-2 rounded-lg transition-colors"
            >
              Apply Now
            </router-link>

            <!-- Mobile menu toggle -->
            <button
              class="md:hidden p-2 rounded-md text-gray-500 hover:text-gray-700"
              @click="mobileOpen = !mobileOpen"
            >
              <svg v-if="!mobileOpen" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
              </svg>
              <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
        </div>

        <!-- Mobile menu -->
        <div v-if="mobileOpen" class="md:hidden pb-4 pt-2 border-t border-gray-100 flex flex-col gap-3 text-sm font-medium">
          <router-link to="/" class="text-gray-700 hover:text-primary-700" @click="mobileOpen=false">Home</router-link>
          <router-link to="/apply" class="text-gray-700 hover:text-primary-700" @click="mobileOpen=false">Apply</router-link>
          <router-link to="/news" class="text-gray-700 hover:text-primary-700" @click="mobileOpen=false">News</router-link>
          <router-link to="/events" class="text-gray-700 hover:text-primary-700" @click="mobileOpen=false">Events</router-link>
          <router-link to="/careers" class="text-gray-700 hover:text-primary-700" @click="mobileOpen=false">Careers</router-link>
          <router-link to="/projects" class="text-gray-700 hover:text-primary-700" @click="mobileOpen=false">Projects</router-link>
          <router-link to="/case-studies" class="text-gray-700 hover:text-primary-700" @click="mobileOpen=false">Case Studies</router-link>
        </div>
      </div>
    </header>

    <!-- Page Content -->
    <main class="flex-1">
      <router-view />
    </main>

    <!-- Footer -->
    <footer class="bg-gray-900 text-gray-300 mt-auto">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
          <div class="md:col-span-2">
            <div class="flex items-center gap-3 mb-4">
              <div class="w-10 h-10 rounded-full bg-primary-600 flex items-center justify-center">
                <span class="text-white font-bold text-sm">NCS</span>
              </div>
              <div>
                <div class="font-bold text-white text-sm">National Council of Sports</div>
                <div class="text-xs text-gray-400">Republic of Uganda</div>
              </div>
            </div>
            <p class="text-sm text-gray-400 max-w-sm">
              The National Council of Sports is the government body responsible for the development,
              promotion and regulation of sports in Uganda.
            </p>
          </div>
          <div>
            <h4 class="font-semibold text-white mb-3 text-sm uppercase tracking-wide">Services</h4>
            <ul class="space-y-2 text-sm">
              <li><router-link to="/apply" class="hover:text-primary-400 transition-colors">Apply for License</router-link></li>
              <li><router-link to="/apply" class="hover:text-primary-400 transition-colors">Renew License</router-link></li>
              <li><a href="#" class="hover:text-primary-400 transition-colors">Verify Certificate</a></li>
            </ul>
          </div>
          <div>
            <h4 class="font-semibold text-white mb-3 text-sm uppercase tracking-wide">Information</h4>
            <ul class="space-y-2 text-sm">
              <li><router-link to="/news" class="hover:text-primary-400 transition-colors">News & Updates</router-link></li>
              <li><router-link to="/events" class="hover:text-primary-400 transition-colors">Events</router-link></li>
              <li><router-link to="/careers" class="hover:text-primary-400 transition-colors">Careers</router-link></li>
              <li><router-link to="/projects" class="hover:text-primary-400 transition-colors">Projects</router-link></li>
            </ul>
          </div>
        </div>
        <div class="border-t border-gray-800 pt-6 text-sm text-gray-500 flex flex-col md:flex-row justify-between gap-2">
          <span>&copy; {{ new Date().getFullYear() }} National Council of Sports, Uganda. All rights reserved.</span>
          <span>Plot 6, Impala Avenue, Kampala, Uganda</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth.js'

const mobileOpen = ref(false)
const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)
</script>
