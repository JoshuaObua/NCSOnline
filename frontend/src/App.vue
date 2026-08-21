<template>
  <div :class="{ dark: isDark }">
    <AppPreloader :loading="isPreloading" />
    <router-view v-slot="{ Component }">
      <Transition name="page-fade" mode="out-in">
        <component :is="Component" />
      </Transition>
    </router-view>
    <LockScreen />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import LockScreen from '@/components/ui/LockScreen.vue'
import AppPreloader from '@/components/ui/AppPreloader.vue'

const authStore = useAuthStore()
const router = useRouter()
const isPreloading = ref(true)
const isDark = ref(false)

onMounted(() => {
  authStore.loadFromStorage()
  isDark.value = document.documentElement.classList.contains('dark') || document.body.classList.contains('dark')
  setTimeout(() => {
    isPreloading.value = false
  }, 350)
})

router.beforeEach((to, from, next) => {
  if (to.path !== from.path) {
    isPreloading.value = true
  }
  next()
})

router.afterEach(() => {
  setTimeout(() => {
    isPreloading.value = false
  }, 300)
})
</script>

<style>
.page-fade-enter-active, .page-fade-leave-active { transition: opacity 0.18s ease; }
.page-fade-enter-from, .page-fade-leave-to { opacity: 0; }
</style>
