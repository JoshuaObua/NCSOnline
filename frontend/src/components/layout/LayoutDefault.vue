<template>
  <div class="min-h-screen bg-gray-50 flex">
    <AppSidebar />
    <div class="flex-1 ml-64 flex flex-col min-h-screen">
      <AppHeader :title="title" />
      <main class="flex-1 p-6 overflow-auto">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup>
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

defineProps({
  title: {
    type: String,
    default: ''
  }
})
</script>
