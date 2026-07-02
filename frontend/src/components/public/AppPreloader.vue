<template>
  <Transition name="preloader">
    <div v-if="visible" class="app-preloader" role="status" aria-live="polite" aria-label="Loading">
      <div class="preloader-stack">
        <div class="loader"></div>
        <img :src="logoSrc" alt="" class="preloader-logo" />
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
  // Lower-bound visibility (ms) so the spinner doesn't flash on fast loads
  minDuration: { type: Number, default: 600 },
  // Hard cap (ms) — never block the page indefinitely if a fetch hangs
  maxDuration: { type: Number, default: 3000 },
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
.preloader-stack {
  position: relative;
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* ── EXACT preloader CSS from the spec ─────────────────────────── */
.loader {
  height: 80px;
  aspect-ratio: 1;
  padding: 10px;
  border-radius: 20px;
  box-sizing: border-box;
  position: relative;
  mask: conic-gradient(#000 0 0) content-box exclude, conic-gradient(#000 0 0);
  filter: blur(12px);
}
.loader:before {
  content: "";
  position: absolute;
  inset: 0;
  background: repeating-conic-gradient(#0000 0 5%, #C02942, #0000 20% 50%);
  animation: l3 1.5s linear infinite;
}
@keyframes l3 { to { rotate: 1turn } }
/* ──────────────────────────────────────────────────────────────── */

/* Centered logo fitting inside the loader ring */
.preloader-logo {
  position: absolute;
  width: 40px;
  height: 40px;
  object-fit: contain;
  border-radius: 50%;
  background: #fff;
  padding: 4px;
  z-index: 1;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.08);
}

/* Fade out when finished */
.preloader-enter-active, .preloader-leave-active { transition: opacity 0.35s ease; }
.preloader-enter-from, .preloader-leave-to { opacity: 0; }
</style>
