import { computed, onMounted, onUnmounted, ref } from 'vue'

const STORAGE_KEY = 'ncs-theme'
const theme = ref(readInitialTheme())

function readInitialTheme() {
  if (typeof document !== 'undefined') {
    const htmlTheme = document.documentElement.getAttribute('data-theme')
    if (htmlTheme === 'dark' || htmlTheme === 'light') return htmlTheme
  }
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored === 'dark' || stored === 'light') return stored
  } catch {
    // Storage can be unavailable in private or restricted browser contexts.
  }
  if (typeof window !== 'undefined' && window.matchMedia?.('(prefers-color-scheme: dark)').matches) return 'dark'
  return 'light'
}

function applyTheme(nextTheme) {
  const normalized = nextTheme === 'dark' ? 'dark' : 'light'
  theme.value = normalized
  if (typeof document !== 'undefined') {
    document.documentElement.classList.toggle('dark', normalized === 'dark')
    document.documentElement.setAttribute('data-theme', normalized)
    document.documentElement.style.colorScheme = normalized
  }
  try {
    localStorage.setItem(STORAGE_KEY, normalized)
  } catch {
    // The DOM state remains authoritative if persistence is unavailable.
  }
}

export function useTheme() {
  const isDark = computed(() => theme.value === 'dark')
  let mediaQuery
  const syncSystemTheme = (event) => {
    try {
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored === 'dark' || stored === 'light') return
    } catch {
      // If storage cannot be read, continue syncing to system preference.
    }
    applyTheme(event.matches ? 'dark' : 'light')
  }

  onMounted(() => {
    applyTheme(theme.value)
    mediaQuery = window.matchMedia?.('(prefers-color-scheme: dark)')
    mediaQuery?.addEventListener?.('change', syncSystemTheme)
  })

  onUnmounted(() => {
    mediaQuery?.removeEventListener?.('change', syncSystemTheme)
  })

  return {
    theme,
    isDark,
    setTheme: applyTheme,
    toggleTheme: () => applyTheme(isDark.value ? 'light' : 'dark'),
  }
}
