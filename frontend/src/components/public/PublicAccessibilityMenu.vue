<template>
  <div ref="menuRoot" class="accessibility-menu">
    <button
      type="button"
      class="accessibility-trigger"
      :aria-expanded="open"
      aria-controls="visitor-accessibility-panel"
      @click="togglePanel"
    >
      <span aria-hidden="true" class="text-xl font-bold">A</span>
      <span class="hidden sm:inline">Accessibility</span>
    </button>

    <section
      v-if="open"
      id="visitor-accessibility-panel"
      ref="panel"
      class="accessibility-panel"
      aria-labelledby="accessibility-panel-title"
    >
      <div class="flex items-start justify-between gap-4 border-b border-gray-200 pb-3">
        <div>
          <h2 id="accessibility-panel-title" class="text-base font-bold text-[#2F327D]">Accessibility tools</h2>
          <p class="mt-0.5 text-xs text-gray-600">Adjust this website for easier reading.</p>
        </div>
        <button type="button" class="accessibility-icon-button" aria-label="Close accessibility tools" @click="closePanel">
          <span aria-hidden="true">&#10005;</span>
        </button>
      </div>

      <div class="py-4">
        <p id="text-size-label" class="mb-2 text-sm font-semibold text-gray-800">Text size</p>
        <div class="grid grid-cols-3 gap-2" role="group" aria-labelledby="text-size-label">
          <button type="button" class="accessibility-option" :disabled="preferences.textScale <= 90" aria-label="Decrease text size" @click="changeTextSize(-10)">A−</button>
          <output class="accessibility-value" aria-live="polite">{{ preferences.textScale }}%</output>
          <button type="button" class="accessibility-option text-lg" :disabled="preferences.textScale >= 130" aria-label="Increase text size" @click="changeTextSize(10)">A+</button>
        </div>
      </div>

      <div class="space-y-2 border-t border-gray-200 pt-4">
        <button type="button" class="accessibility-toggle" :aria-pressed="preferences.highContrast" @click="togglePreference('highContrast')">
          <span>High contrast</span><span class="toggle-state">{{ preferences.highContrast ? 'On' : 'Off' }}</span>
        </button>
        <button type="button" class="accessibility-toggle" :aria-pressed="preferences.underlineLinks" @click="togglePreference('underlineLinks')">
          <span>Underline links</span><span class="toggle-state">{{ preferences.underlineLinks ? 'On' : 'Off' }}</span>
        </button>
        <button type="button" class="accessibility-toggle" :aria-pressed="preferences.reduceMotion" @click="togglePreference('reduceMotion')">
          <span>Reduce motion</span><span class="toggle-state">{{ preferences.reduceMotion ? 'On' : 'Off' }}</span>
        </button>
      </div>

      <button type="button" class="mt-4 w-full rounded-lg border border-[#2F327D] px-3 py-2 text-sm font-semibold text-[#2F327D] hover:bg-[#2F327D] hover:text-white" @click="resetPreferences">
        Reset accessibility settings
      </button>
    </section>

    <p class="sr-only" aria-live="polite">{{ announcement }}</p>
  </div>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'

const STORAGE_KEY = 'ncs_public_accessibility'
const defaults = { textScale: 100, highContrast: false, underlineLinks: false, reduceMotion: false }

const open = ref(false)
const panel = ref(null)
const menuRoot = ref(null)
const announcement = ref('')
const preferences = reactive({ ...defaults })

function readPreferences() {
  try {
    const saved = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}')
    preferences.textScale = [90, 100, 110, 120, 130].includes(saved.textScale) ? saved.textScale : 100
    preferences.highContrast = saved.highContrast === true
    preferences.underlineLinks = saved.underlineLinks === true
    preferences.reduceMotion = saved.reduceMotion === true
  } catch {
    Object.assign(preferences, defaults)
  }
}

function applyPreferences() {
  const site = document.querySelector('.public-site')
  if (!site) return
  document.documentElement.style.fontSize = `${preferences.textScale}%`
  site.classList.toggle('a11y-high-contrast', preferences.highContrast)
  site.classList.toggle('a11y-underline-links', preferences.underlineLinks)
  site.classList.toggle('a11y-reduce-motion', preferences.reduceMotion)
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(preferences))
  } catch { /* preferences still apply for this visit */ }
}

function togglePanel() {
  open.value = !open.value
  if (open.value) nextTick(() => panel.value?.querySelector('button')?.focus())
}

function closePanel() {
  open.value = false
  nextTick(() => menuRoot.value?.querySelector('.accessibility-trigger')?.focus())
}

function changeTextSize(amount) {
  preferences.textScale = Math.min(130, Math.max(90, preferences.textScale + amount))
  announcement.value = `Text size ${preferences.textScale} percent`
}

function togglePreference(name) {
  preferences[name] = !preferences[name]
  const labels = { highContrast: 'High contrast', underlineLinks: 'Underline links', reduceMotion: 'Reduced motion' }
  announcement.value = `${labels[name]} ${preferences[name] ? 'enabled' : 'disabled'}`
}

function resetPreferences() {
  Object.assign(preferences, defaults)
  announcement.value = 'Accessibility settings reset'
}

function onKeydown(event) {
  if (event.key === 'Escape' && open.value) closePanel()
}

function onPointerDown(event) {
  if (open.value && menuRoot.value && !menuRoot.value.contains(event.target)) open.value = false
}

watch(preferences, applyPreferences, { deep: true })

onMounted(() => {
  readPreferences()
  applyPreferences()
  document.addEventListener('keydown', onKeydown)
  document.addEventListener('pointerdown', onPointerDown)
})

onBeforeUnmount(() => {
  document.documentElement.style.fontSize = ''
  document.removeEventListener('keydown', onKeydown)
  document.removeEventListener('pointerdown', onPointerDown)
})
</script>
