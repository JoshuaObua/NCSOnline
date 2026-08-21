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
import { ref, onMounted, computed, watch } from 'vue'
import { getSettings } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'
import { useTheme } from '@/composables/useTheme.js'

const props = defineProps({
  loading: { type: Boolean, default: undefined },
  minDuration: { type: Number, default: 350 },
})

const themeState = useTheme()
const isDark = computed(() => themeState?.isDark?.value ?? false)

const initialVisible = ref(true)
const propVisible = ref(props.loading ?? false)

watch(() => props.loading, (val) => {
  propVisible.value = !!val
})

const visible = computed(() => initialVisible.value || propVisible.value)

const siteLogo = ref('')
const logoSrc = computed(() => siteLogo.value ? mediaUrl(siteLogo.value) : '/main-logo.png')

onMounted(() => {
  try {
    getSettings('site').then(r => {
      if (r?.data?.data?.value?.logoUrl) siteLogo.value = r.data.data.value.logoUrl
    }).catch(() => {})
  } catch {}

  setTimeout(() => {
    initialVisible.value = false
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
  transition: opacity 0.3s ease, background-color 0.2s ease;
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
  transition: filter 0.2s ease;
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
