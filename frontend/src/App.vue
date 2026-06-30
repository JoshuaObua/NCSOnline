<template>
  <router-view v-slot="{ Component }">
    <Transition name="page-fade" mode="out-in">
      <component :is="Component" />
    </Transition>
  </router-view>
  <LockScreen />
</template>

<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth.js'
import LockScreen from '@/components/ui/LockScreen.vue'

const authStore = useAuthStore()

onMounted(() => {
  authStore.loadFromStorage()
})
</script>

<style>
.page-fade-enter-active, .page-fade-leave-active { transition: opacity 0.18s ease; }
.page-fade-enter-from, .page-fade-leave-to { opacity: 0; }
</style>
