<template>
  <Transition name="preloader-fade">
    <div
      v-if="loading"
      class="app-preloader"
      :class="{ 'is-dark': isDark }"
      role="status"
      aria-live="polite"
      aria-label="Loading"
    >
      <div class="preloader-content">
        <img
          src="/main-logo.png"
          alt="NCS Logo"
          class="preloader-logo"
        />
        <div class="preloader-spinner"></div>
      </div>
      <span class="sr-only">Loading…</span>
    </div>
  </Transition>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  loading: { type: Boolean, default: false },
})

const isDark = computed(() => {
  if (typeof document === 'undefined') return false
  return document.documentElement.classList.contains('dark') || document.body.classList.contains('dark')
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
  transition: opacity 0.25s ease, background-color 0.2s ease;
}

.app-preloader.is-dark {
  background: #0f172a !important;
}

:global(.dark .app-preloader) {
  background: #0f172a !important;
}

:global([data-theme="dark"] .app-preloader) {
  background: #0f172a !important;
}

.preloader-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1.25rem;
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
  width: 34px;
  height: 34px;
  border: 3.5px solid rgba(245, 166, 35, 0.25);
  border-top-color: #f5a623;
  border-radius: 50%;
  animation: preloader-spin 0.8s linear infinite;
}

.app-preloader.is-dark .preloader-logo {
  filter: brightness(0) invert(1) !important;
}

:global(.dark .app-preloader .preloader-logo) {
  filter: brightness(0) invert(1) !important;
}

:global([data-theme="dark"] .app-preloader .preloader-logo) {
  filter: brightness(0) invert(1) !important;
}

@keyframes preloader-pulse {
  0%, 100% { transform: scale(1); opacity: 0.95; }
  50% { transform: scale(1.04); opacity: 1; }
}

@keyframes preloader-spin {
  to { transform: rotate(360deg); }
}

.preloader-fade-enter-active, .preloader-fade-leave-active {
  transition: opacity 0.25s ease;
}
.preloader-fade-enter-from, .preloader-fade-leave-to {
  opacity: 0;
}
</style>
