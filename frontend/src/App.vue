<template>
  <AppPreloader :key="preloaderKey" />
  <router-view v-slot="{ Component }">
    <Transition name="page-fade" mode="out-in">
      <component :is="Component" />
    </Transition>
  </router-view>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AppPreloader from '@/components/public/AppPreloader.vue'

const router = useRouter()
const preloaderKey = ref(0)

router.beforeEach((to, from, next) => {
  if (from.name && to.path !== from.path) {
    preloaderKey.value++
  }
  next()
})
</script>

<style>
.page-fade-enter-active, .page-fade-leave-active { transition: opacity 0.18s ease; }
.page-fade-enter-from, .page-fade-leave-to { opacity: 0; }
</style>
