<template>
  <router-view />
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
