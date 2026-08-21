<template>
  <Transition name="preloader">
    <div v-if="visible" class="app-preloader" role="status" aria-live="polite" aria-label="Loading">
      <div class="preloader-content">
        <img :src="logoSrc" alt="NCS Logo" class="preloader-logo" />
      </div>
      <span class="sr-only">Loading…</span>
    </div>
  </Transition>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { getSettings } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const props = defineProps({
  minDuration: { type: Number, default: 400 },
  maxDuration: { type: Number, default: 2000 },
})

const visible = ref(true)
const siteLogo = ref('')
const logoSrc = computed(() => siteLogo.value ? mediaUrl(siteLogo.value) : '/main-logo.png')

onMounted(async () => {
  const start = performance.now()
  let pageReady = document.readyState === 'complete'
  if (!pageReady) {
    await new Promise(resolve => {
      const done = () => { window.removeEventListener('load', done); resolve() }
      window.addEventListener('load', done, { once: true })
      setTimeout(resolve, props.maxDuration)
    })
  }
  try {
    const r = await Promise.race([
      getSettings('site'),
      new Promise(resolve => setTimeout(() => resolve(null), props.maxDuration - (performance.now() - start))),
    ])
    if (r?.data?.data?.value?.logoUrl) siteLogo.value = r.data.data.value.logoUrl
  } catch { /* defaults fine */ }
  const elapsed = performance.now() - start
  if (elapsed < props.minDuration) await new Promise(r => setTimeout(r, props.minDuration - elapsed))
  visible.value = false
})
</script>

<style scoped>
.app-preloader {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #ffffff;
}
.preloader-content {
  display: flex;
  align-items: center;
  justify-content: center;
}
.preloader-logo {
  width: 200px;
  max-width: 85vw;
  height: auto;
  object-fit: contain;
}
.preloader-enter-active, .preloader-leave-active { transition: opacity 0.35s ease; }
.preloader-enter-from, .preloader-leave-to { opacity: 0; }
</style>
