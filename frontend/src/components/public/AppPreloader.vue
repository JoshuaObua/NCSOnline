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
        <img
          :src="logoSrc"
          alt="NCS Logo"
          class="preloader-logo"
          @load="onLogoLoaded"
        />
        <div v-if="!logoLoaded" class="preloader-spinner"></div>
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
  minDuration: { type: Number, default: 450 },
})

const themeState = useTheme()
const isDark = computed(() => themeState?.isDark?.value ?? false)

const initialVisible = ref(true)
const propVisible = ref(props.loading ?? false)
const logoLoaded = ref(false)

function onLogoLoaded() {
  logoLoaded.value = true
}

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
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}

.preloader-logo {
  width: 200px;
  max-width: 85vw;
  height: auto;
  min-height: 50px;
  object-fit: contain;
  transition: filter 0.2s ease, transform 0.3s ease;
  animation: preloader-pulse 1.5s ease-in-out infinite;
}

.preloader-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(245, 166, 35, 0.2);
  border-top-color: #f5a623;
  border-radius: 50%;
  animation: preloader-spin 0.8s linear infinite;
}

.app-preloader.is-dark .preloader-logo,
:global(.dark) .preloader-logo,
:global([data-theme="dark"]) .preloader-logo {
  filter: brightness(0) invert(1) !important;
}

@keyframes preloader-pulse {
  0%, 100% { transform: scale(1); opacity: 0.95; }
  50% { transform: scale(1.04); opacity: 1; }
}

@keyframes preloader-spin {
  to { transform: rotate(360deg); }
}

.preloader-enter-active, .preloader-leave-active {
  transition: opacity 0.3s ease;
}
.preloader-enter-from, .preloader-leave-to {
  opacity: 0;
}
</style>
