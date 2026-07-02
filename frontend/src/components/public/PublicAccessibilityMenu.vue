<template>
  <div
    ref="menuRoot"
    class="accessibility-menu"
    @pointerenter="openFromHover"
    @pointerleave="scheduleClose"
    @focusin="cancelScheduledClose"
    @focusout="onFocusOut"
  >
    <button
      type="button"
      class="accessibility-trigger"
      aria-label="Open accessibility tools"
      title="Accessibility tools"
      :aria-expanded="open"
      aria-controls="visitor-accessibility-panel"
      @focus="openPanel"
      @click="openPanel"
    >
      <svg class="accessibility-trigger-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false">
        <circle cx="16" cy="4" r="1" />
        <path d="m18 19 1-7-6 1" />
        <path d="m5 8 3-3 5.5 3-2.36 3.5" />
        <path d="M4.24 14.5a5 5 0 0 0 6.88 6" />
        <path d="M13.76 17.5a5 5 0 0 0-6.88-6" />
      </svg>
    </button>

    <Transition name="accessibility-panel">
      <section
        v-if="open"
        id="visitor-accessibility-panel"
        ref="panel"
        class="accessibility-panel"
        role="dialog"
        aria-modal="false"
        aria-labelledby="accessibility-panel-title"
        @keydown="trapFocus"
      >
        <div class="accessibility-panel-header">
          <div>
            <h2 id="accessibility-panel-title">Accessibility tools</h2>
            <p>Choose the display that works best for you.</p>
          </div>
          <button type="button" class="accessibility-icon-button" aria-label="Close accessibility tools" @click="closePanel(true)">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path stroke-linecap="round" d="m6 6 12 12M18 6 6 18" /></svg>
          </button>
        </div>

        <div class="accessibility-section">
          <div class="accessibility-action-grid">
            <button type="button" class="accessibility-toggle" @click="readSelectedText"><span>Read aloud</span><span aria-hidden="true">▶</span></button>
            <button type="button" class="accessibility-toggle" @click="defineSelectedWord"><span>Define word</span><span aria-hidden="true">?</span></button>
          </div>
          <p class="accessibility-selection-help">Highlight text first, then choose Read aloud or Define word.</p>
        </div>

        <div class="accessibility-section">
          <p id="text-size-label" class="accessibility-section-label">Text size</p>
          <div class="accessibility-size-controls" role="group" aria-labelledby="text-size-label">
            <button type="button" class="accessibility-option" :disabled="preferences.textScale <= 90" aria-label="Decrease text size" @click="changeTextSize(-10)">A−</button>
            <output class="accessibility-value" aria-live="polite">{{ preferences.textScale }}%</output>
            <button type="button" class="accessibility-option accessibility-option-large" :disabled="preferences.textScale >= 130" aria-label="Increase text size" @click="changeTextSize(10)">A+</button>
          </div>
        </div>

        <div class="accessibility-toggle-grid">
          <button v-for="tool in tools" :key="tool.key" type="button" class="accessibility-toggle" :aria-pressed="preferences[tool.key]" @click="togglePreference(tool.key)">
            <span>{{ tool.label }}</span><span class="toggle-state">{{ preferences[tool.key] ? 'On' : 'Off' }}</span>
          </button>
        </div>

        <button type="button" class="accessibility-reset" @click="resetPreferences">Reset all settings</button>
      </section>
    </Transition>

    <div v-if="preferences.readingGuide" class="a11y-reading-guide" :style="{ top: `${pointerY}px` }" aria-hidden="true"></div>
    <p class="sr-only" aria-live="polite">{{ announcement }}</p>
  </div>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'

const STORAGE_KEY = 'ncs_public_accessibility'
const defaults = {
  textScale: 100,
  highContrast: false,
  grayscale: false,
  readableFont: false,
  textSpacing: false,
  underlineLinks: false,
  highlightFocus: false,
  largeCursor: false,
  hideImages: false,
  reduceMotion: false,
  readingGuide: false,
}
const tools = [
  { key: 'highContrast', label: 'High contrast' },
  { key: 'grayscale', label: 'Grayscale' },
  { key: 'readableFont', label: 'Readable font' },
  { key: 'textSpacing', label: 'Text spacing' },
  { key: 'underlineLinks', label: 'Underline links' },
  { key: 'highlightFocus', label: 'Highlight focus' },
  { key: 'largeCursor', label: 'Large cursor' },
  { key: 'hideImages', label: 'Hide images' },
  { key: 'reduceMotion', label: 'Reduce motion' },
  { key: 'readingGuide', label: 'Reading guide' },
]

const open = ref(false)
const panel = ref(null)
const menuRoot = ref(null)
const announcement = ref('')
const pointerY = ref(window.innerHeight / 2)
const preferences = reactive({ ...defaults })
let closeTimer

function readPreferences() {
  try {
    const saved = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}')
    for (const key of Object.keys(defaults)) {
      if (key === 'textScale') preferences[key] = [90, 100, 110, 120, 130].includes(saved[key]) ? saved[key] : defaults[key]
      else preferences[key] = saved[key] === true
    }
  } catch {
    Object.assign(preferences, defaults)
  }
}

function applyPreferences() {
  const site = document.querySelector('.public-site')
  if (!site) return
  document.documentElement.style.fontSize = `${preferences.textScale}%`
  for (const key of Object.keys(defaults)) {
    if (key !== 'textScale') site.classList.toggle(`a11y-${key.replace(/[A-Z]/g, letter => `-${letter.toLowerCase()}`)}`, preferences[key])
  }
  try { localStorage.setItem(STORAGE_KEY, JSON.stringify(preferences)) } catch { /* Applied for this visit. */ }
}

function cancelScheduledClose() { clearTimeout(closeTimer) }
function openFromHover(event) {
  if (event.pointerType === 'mouse') {
    cancelScheduledClose()
    open.value = true
  }
}
function scheduleClose(event) {
  if (event.pointerType === 'mouse') closeTimer = setTimeout(() => { open.value = false }, 220)
}
function onFocusOut(event) {
  if (!menuRoot.value?.contains(event.relatedTarget)) scheduleClose({ pointerType: 'mouse' })
}
function openPanel() {
  cancelScheduledClose()
  open.value = true
}
function closePanel(returnFocus = false) {
  cancelScheduledClose()
  open.value = false
  if (returnFocus) nextTick(() => menuRoot.value?.querySelector('.accessibility-trigger')?.focus())
}
function changeTextSize(amount) {
  preferences.textScale = Math.min(130, Math.max(90, preferences.textScale + amount))
  announcement.value = `Text size ${preferences.textScale} percent`
}
function togglePreference(name) {
  preferences[name] = !preferences[name]
  const label = tools.find(tool => tool.key === name)?.label || name
  announcement.value = `${label} ${preferences[name] ? 'enabled' : 'disabled'}`
}
function resetPreferences() {
  Object.assign(preferences, defaults)
  announcement.value = 'All accessibility settings reset'
}
function trapFocus(event) {
  if (event.key !== 'Tab') return
  const items = [...panel.value.querySelectorAll('button:not(:disabled), [href], input, select, textarea, [tabindex]:not([tabindex="-1"])')]
  if (!items.length) return
  const first = items[0]
  const last = items[items.length - 1]
  if (event.shiftKey && document.activeElement === first) { event.preventDefault(); last.focus() }
  else if (!event.shiftKey && document.activeElement === last) { event.preventDefault(); first.focus() }
}
function onKeydown(event) { if (event.key === 'Escape' && open.value) closePanel(true) }
function onPointerDown(event) { if (open.value && menuRoot.value && !menuRoot.value.contains(event.target)) closePanel() }
function onPointerMove(event) { if (preferences.readingGuide) pointerY.value = event.clientY }
function readSelectedText() {
  const text = window.getSelection()?.toString().trim() || document.querySelector('main')?.innerText?.slice(0, 3000) || ''
  if (!text || !('speechSynthesis' in window)) { announcement.value = 'No readable text is selected'; return }
  window.speechSynthesis.cancel()
  window.speechSynthesis.speak(new SpeechSynthesisUtterance(text))
  announcement.value = 'Reading selected text aloud'
}
function defineSelectedWord() {
  const word = (window.getSelection()?.toString().trim() || '').split(/\s+/)[0]?.replace(/[^a-zA-Z'-]/g, '')
  if (!word) { announcement.value = 'Highlight one word to define it'; return }
  window.open(`https://www.merriam-webster.com/dictionary/${encodeURIComponent(word)}`, '_blank', 'noopener')
  announcement.value = `Opening the definition of ${word}`
}
function onOpenRequest() {
  cancelScheduledClose()
  open.value = true
  // The menu lives at bottom-left; if triggered from elsewhere on the page (e.g. the
  // header button), bring it into view so the user actually sees the panel.
  nextTick(() => {
    menuRoot.value?.scrollIntoView({ behavior: 'smooth', block: 'end' })
    panel.value?.querySelector('button')?.focus()
  })
}

watch(preferences, applyPreferences, { deep: true })
onMounted(() => {
  readPreferences()
  applyPreferences()
  document.addEventListener('keydown', onKeydown)
  document.addEventListener('pointerdown', onPointerDown)
  document.addEventListener('pointermove', onPointerMove, { passive: true })
  window.addEventListener('open-accessibility-menu', onOpenRequest)
})
onBeforeUnmount(() => {
  clearTimeout(closeTimer)
  document.documentElement.style.fontSize = ''
  document.removeEventListener('keydown', onKeydown)
  document.removeEventListener('pointerdown', onPointerDown)
  document.removeEventListener('pointermove', onPointerMove)
  window.removeEventListener('open-accessibility-menu', onOpenRequest)
})
</script>
