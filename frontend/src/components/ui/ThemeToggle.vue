<template>
  <button
    type="button"
    class="theme-toggle-btn group relative"
    @click="toggleTheme"
    :title="isDark ? 'Switch to Light Theme' : 'Switch to Dark Theme'"
    aria-label="Toggle dark/light mode"
  >
    <div class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl border transition-all duration-200"
      :class="isDark 
        ? 'bg-slate-800 border-slate-700 text-amber-300 hover:bg-slate-750' 
        : 'bg-white border-slate-200 text-slate-700 hover:bg-slate-50 shadow-sm'"
    >
      <i v-if="isDark" class="icofont-sun text-base text-amber-400"></i>
      <i v-else class="icofont-moon text-base text-slate-700"></i>
      <span class="text-xs font-bold select-none tracking-wide">
        {{ isDark ? 'Dark' : 'Light' }}
      </span>
    </div>
  </button>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const isDark = ref(false)

function applyTheme(dark) {
  isDark.value = dark
  if (dark) {
    document.documentElement.classList.add('dark')
    localStorage.setItem('ncs_theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('ncs_theme', 'light')
  }
}

function toggleTheme() {
  applyTheme(!isDark.value)
}

onMounted(() => {
  const isCurrentlyDark = document.documentElement.classList.contains('dark')
  const saved = localStorage.getItem('ncs_theme')
  if (saved) {
    applyTheme(saved === 'dark')
  } else {
    applyTheme(isCurrentlyDark)
  }
})
</script>

<style scoped>
.theme-toggle-btn {
  background: transparent;
  border: 0;
  padding: 0;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
</style>
