<template>
  <button
    ref="buttonRef"
    type="button"
    class="theme-toggle"
    :aria-pressed="isDark ? 'true' : 'false'"
    :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
    title="Toggle light or dark mode"
    @click="activate"
    @keydown.space.prevent="activate"
    @keydown.enter.prevent="activate"
  >
    <i :class="isDark ? 'icofont-sun' : 'icofont-moon'" aria-hidden="true"></i>
    <span class="sr-only">{{ isDark ? 'Light mode' : 'Dark mode' }}</span>
  </button>
</template>

<script setup>
import { nextTick, ref } from 'vue'
import { useTheme } from '@/composables/useTheme.js'

const buttonRef = ref(null)
const { isDark, toggleTheme } = useTheme()

async function activate() {
  toggleTheme()
  await nextTick()
  buttonRef.value?.focus({ preventScroll: true })
}
</script>

<style scoped>
.theme-toggle {
  display: inline-flex;
  width: 2.5rem;
  height: 2.5rem;
  align-items: center;
  justify-content: center;
  border-radius: 0.5rem;
  color: #1a365d;
  transition: background-color 160ms ease, color 160ms ease, transform 160ms ease;
}
.theme-toggle:hover,
.theme-toggle:focus-visible {
  background: rgb(245 166 35 / 0.12);
  color: #a84b00;
}
.theme-toggle:active { transform: scale(0.96); }
:global(.dark) .theme-toggle { color: #f8fafc; }
</style>
