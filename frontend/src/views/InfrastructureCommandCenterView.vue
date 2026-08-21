<template>
  <LayoutDefault>
    <div class="max-w-7xl mx-auto p-4 md:p-6">
      <SystemCommandCenterPanel @message="handleMsg" @error="handleErr" />
    </div>
  </LayoutDefault>
</template>

<script setup>
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import SystemCommandCenterPanel from '@/components/cms/SystemCommandCenterPanel.vue'

function handleMsg(msg) {
  console.log('[SystemCommandCenter]', msg)
}
function handleErr(err) {
  console.error('[SystemCommandCenter Error]', err)
}
</script>
