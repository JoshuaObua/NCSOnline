<template>
  <div class="min-h-screen bg-gray-50 flex">
    <a href="#main-content" class="skip-link">Skip to main content</a>

    <!-- Mobile overlay -->
    <Transition name="fade-overlay">
      <div
        v-if="sidebarOpen && isMobile"
        class="fixed inset-0 bg-black/40 z-30"
        @click="sidebarOpen = false"
      ></div>
    </Transition>

    <AppSidebar :open="sidebarOpen" @close="sidebarOpen = false" />

    <div class="flex-1 flex flex-col min-h-screen md:ml-64">
      <AppHeader :page-title="pageTitle" @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <!-- Breadcrumb row -->
      <div class="bg-white border-b border-gray-100 px-6 py-2.5 flex items-center gap-1.5 text-xs flex-shrink-0">
        <router-link to="/dashboard" class="text-gray-400 hover:text-primary-700 transition-colors">
          <i class="icofont-home text-sm leading-none"></i>
        </router-link>
        <template v-if="crumbs.length">
          <span v-for="(crumb, i) in crumbs" :key="i" class="flex items-center gap-1.5">
            <svg class="w-3 h-3 text-gray-300" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
            </svg>
            <router-link v-if="crumb.to" :to="crumb.to" class="text-gray-500 hover:text-primary-700 transition-colors font-medium">{{ crumb.label }}</router-link>
            <span v-else class="text-gray-700 font-semibold">{{ crumb.label }}</span>
          </span>
        </template>
        <template v-else-if="title">
          <svg class="w-3 h-3 text-gray-300" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
          </svg>
          <span class="text-gray-700 font-semibold">{{ title }}</span>
        </template>
      </div>

      <!-- Main content -->
      <main id="main-content" tabindex="-1" class="flex-1 p-4 sm:p-6 overflow-auto min-w-0">
        <slot />
      </main>

      <!-- Footer -->
      <footer class="flex-shrink-0 px-4 sm:px-6 py-3 border-t border-gray-100 bg-white flex flex-col sm:flex-row gap-1 sm:items-center sm:justify-between text-xs text-gray-400">
        <span>&copy; {{ year }} National Council of Sports, Uganda. All rights reserved.</span>
        <span class="font-medium">NCSMS v1</span>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'

const props = defineProps({
  title: { type: String, default: '' }
})

const breadcrumbStore = useBreadcrumbStore()

const pageTitle = computed(() => breadcrumbStore.title || props.title)
const crumbs    = computed(() => breadcrumbStore.crumbs)

const year = new Date().getFullYear()

// Sidebar state
const sidebarOpen = ref(true)
const isMobile = ref(false)

function checkMobile() {
  isMobile.value = window.innerWidth < 768
  if (isMobile.value) sidebarOpen.value = false
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})
onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.fade-overlay-enter-active, .fade-overlay-leave-active { transition: opacity 0.2s; }
.fade-overlay-enter-from, .fade-overlay-leave-to { opacity: 0; }
</style>
