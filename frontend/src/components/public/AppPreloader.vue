<template>
  <Transition name="preloader">
    <div
      v-if="visible"
      class="app-preloader"
      :class="{ 'is-dark': isDark }"
      role="status"
      aria-live="polite"
      aria-label="Loading"
    >
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
import { useTheme } from '@/composables/useTheme.js'

const props = defineProps({
  minDuration: { type: Number, default: 300 },
  maxDuration: { type: Number, default: 800 },
})

const { isDark } = useTheme()
const visible = ref(true)
const siteLogo = ref('')
const logoSrc = computed(() => siteLogo.value ? mediaUrl(siteLogo.value) : '/main-logo.png')

onMounted(() => {
  try {
    getSettings('site').then(r => {
      if (r?.data?.data?.value?.logoUrl) siteLogo.value = r.data.data.value.logoUrl
    }).catch(() => {})
  } catch {}

  setTimeout(() => {
    visible.value = false
  }, props.minDuration)
})
</script>

<style scoped>
.app-preloader {
  position: fixed;
  inset: 0;
  z-index: 99999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #ffffff;
  pointer-events: none;
}

.app-preloader.is-dark,
:global(.dark) .app-preloader,
:global([data-theme="dark"]) .app-preloader {
  background: #0f172a !important;
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

.app-preloader.is-dark .preloader-logo,
:global(.dark) .preloader-logo,
:global([data-theme="dark"]) .preloader-logo {
  filter: brightness(0) invert(1) !important;
}

.preloader-enter-active, .preloader-leave-active {
  transition: opacity 0.3s ease;
}
.preloader-enter-from, .preloader-leave-to {
  opacity: 0;
}
</style>
