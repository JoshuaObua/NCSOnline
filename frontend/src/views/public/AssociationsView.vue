<template>
  <div>
    <!-- Hero -->
    <section class="bg-[#1a365d] py-16 md:py-20">
      <div class="max-w-7xl mx-auto px-4">
        <nav class="flex items-center gap-2 text-white/60 text-sm mb-4">
          <router-link to="/" class="hover:text-white">Home</router-link>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="m9 18 6-6-6-6"></path></svg>
          <span class="text-white">Associations</span>
        </nav>
        <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">Sports Associations &amp; Federations</h1>
        <p class="text-white/80 text-lg max-w-2xl">Official directory of all recognized National Sports Associations affiliated with the National Council of Sports. Click any federation card to view full leadership and contact details.</p>
      </div>
    </section>

    <!-- Toolbar & Category Filter -->
    <section class="bg-white border-b border-gray-100">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="flex flex-col md:flex-row items-center justify-between gap-4">
          <div class="flex items-center gap-6">
            <div class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623]" aria-hidden="true"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path><path d="M16 3.128a4 4 0 0 1 0 7.744"></path><path d="M22 21v-2a4 4 0 0 0-3-3.87"></path><circle cx="9" cy="7" r="4"></circle></svg>
              <span class="font-bold text-[#1a365d] text-lg">{{ filteredAssociations.length }}</span>
              <span class="text-gray-500 font-medium">Registered Federations</span>
            </div>
          </div>
          <div class="flex items-center gap-4 w-full md:w-auto">
            <div class="relative w-full md:w-80">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" aria-hidden="true"><path d="m21 21-4.34-4.34"></path><circle cx="11" cy="11" r="8"></circle></svg>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Search by name, sport, president..."
                class="flex h-9 w-full rounded-md border border-gray-300 bg-white pl-10 pr-3 py-1 text-sm shadow-sm placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-[#f5a623] focus:border-[#f5a623]"
              >
            </div>
            <router-link
              to="/contact-us"
              class="inline-flex items-center justify-center gap-2 rounded-md text-sm font-medium shadow h-9 px-4 py-2 bg-[#f5a623] hover:bg-[#e09612] text-white whitespace-nowrap transition-colors"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="m15.477 12.89 1.515 8.526a.5.5 0 0 1-.81.47l-3.58-2.687a1 1 0 0 0-1.197 0l-3.586 2.686a.5.5 0 0 1-.81-.469l1.514-8.526"></path><circle cx="12" cy="8" r="6"></circle></svg>
              Register Federation
            </router-link>
          </div>
        </div>

        <!-- Category Chips -->
        <div class="flex items-center gap-2 overflow-x-auto pt-4 mt-2 pb-1">
          <button
            v-for="cat in categories"
            :key="cat"
            type="button"
            class="px-3.5 py-1.5 rounded-full text-xs font-semibold whitespace-nowrap transition-all duration-200"
            :class="selectedCategory === cat ? 'bg-[#1a365d] text-white shadow-sm' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
            @click="selectedCategory = cat"
          >
            {{ cat }}
          </button>
        </div>
      </div>
    </section>

    <!-- Associations Grid -->
    <section class="py-12 md:py-16 bg-gray-50/50">
      <div class="max-w-7xl mx-auto px-4">
        <div v-if="loading" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div v-for="i in 6" :key="i" class="h-48 bg-gray-200/60 rounded-xl animate-pulse"/>
        </div>

        <div v-else-if="filteredAssociations.length" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <button
            v-for="a in filteredAssociations"
            :key="a.id"
            type="button"
            class="group text-left bg-white rounded-xl p-6 shadow-sm hover:shadow-xl transition-all duration-300 border border-gray-100 hover:border-[#f5a623]/40 hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-[#f5a623]/50 flex flex-col justify-between"
            @click="selected = a"
          >
            <div>
              <div class="flex items-start gap-4 mb-4">
                <div class="w-16 h-16 rounded-lg bg-[#1a365d]/10 flex items-center justify-center flex-shrink-0 overflow-hidden p-2 border border-gray-100">
                  <img v-if="a.logo_url" :src="mediaUrl(a.logo_url)" :alt="a.name" class="max-w-full max-h-full object-contain">
                  <svg v-else xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-8 h-8 text-[#1a365d]" aria-hidden="true"><path d="m15.477 12.89 1.515 8.526a.5.5 0 0 1-.81.47l-3.58-2.687a1 1 0 0 0-1.197 0l-3.586 2.686a.5.5 0 0 1-.81-.469l1.514-8.526"></path><circle cx="12" cy="8" r="6"></circle></svg>
                </div>
                <div class="flex-1 min-w-0">
                  <span class="inline-block text-[11px] font-bold text-[#f5a623] uppercase tracking-wider mb-1">{{ a.category || 'Sport Federation' }}</span>
                  <h3 class="font-bold text-[#1a365d] group-hover:text-[#f5a623] transition-colors leading-snug text-base">{{ a.name }}</h3>
                  <div v-if="a.abbreviation" class="inline-flex items-center rounded-md border border-transparent bg-gray-100 text-gray-700 px-2.5 py-0.5 text-xs font-semibold mt-1">{{ a.abbreviation }}</div>
                </div>
              </div>
              <p v-if="a.description" class="text-gray-600 text-sm mb-4 line-clamp-2 leading-relaxed">{{ a.description }}</p>
            </div>

            <div class="flex items-center justify-between pt-4 border-t border-gray-100 mt-2">
              <div class="flex items-center gap-3 text-xs text-gray-500 min-w-0">
                <span v-if="a.president" class="flex items-center gap-1 min-w-0">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 flex-shrink-0 text-[#1a365d]" aria-hidden="true"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                  <span class="truncate font-medium">{{ a.president }}</span>
                </span>
                <span v-if="a.address" class="flex items-center gap-1 min-w-0">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 flex-shrink-0 text-[#1a365d]" aria-hidden="true"><path d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0"></path><circle cx="12" cy="10" r="3"></circle></svg>
                  <span class="truncate">{{ a.address }}</span>
                </span>
              </div>
              <span class="flex-shrink-0 text-xs font-semibold text-[#f5a623] group-hover:underline ml-2">View Details →</span>
            </div>
          </button>
        </div>

        <div v-else class="text-center py-20 text-gray-400 bg-white rounded-xl border border-gray-100 shadow-sm">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-14 h-14 mx-auto mb-3 opacity-30 text-[#1a365d]" aria-hidden="true"><path d="m21 21-4.34-4.34"></path><circle cx="11" cy="11" r="8"></circle></svg>
          <p v-if="searchQuery" class="text-base text-gray-600 font-medium">No associations match "{{ searchQuery }}".</p>
          <p v-else class="text-base text-gray-600 font-medium">No associations listed in {{ selectedCategory }}.</p>
        </div>
      </div>
    </section>

    <!-- Details Modal -->
    <Teleport to="body">
      <Transition name="assoc-modal">
        <div v-if="selected" class="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="selected = null"></div>
          <div
            class="relative bg-white dark:bg-slate-900 rounded-xl shadow-2xl max-w-2xl w-full max-h-[85vh] overflow-hidden flex flex-col border border-gray-100 dark:border-slate-800"
            role="dialog"
            aria-modal="true"
            :aria-label="selected.name"
          >
            <!-- Header -->
            <div class="bg-gradient-to-br from-[#1a365d] to-[#2d4a7a] p-6 text-white relative flex-shrink-0">
              <button
                type="button"
                class="absolute top-4 right-4 w-8 h-8 rounded-full bg-white/10 hover:bg-white/20 flex items-center justify-center transition-colors"
                aria-label="Close"
                @click="selected = null"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="w-4 h-4" aria-hidden="true"><path d="M18 6 6 18"></path><path d="m6 6 12 12"></path></svg>
              </button>
              <div class="flex items-start gap-4 pr-8">
                <div class="w-16 h-16 rounded-lg bg-white/10 backdrop-blur flex items-center justify-center flex-shrink-0 overflow-hidden p-2">
                  <img v-if="selected.logo_url" :src="mediaUrl(selected.logo_url)" :alt="selected.name" class="max-w-full max-h-full object-contain">
                  <svg v-else xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-8 h-8 text-[#f5a623]" aria-hidden="true"><path d="m15.477 12.89 1.515 8.526a.5.5 0 0 1-.81.47l-3.58-2.687a1 1 0 0 0-1.197 0l-3.586 2.686a.5.5 0 0 1-.81-.469l1.514-8.526"></path><circle cx="12" cy="8" r="6"></circle></svg>
                </div>
                <div class="flex-1 min-w-0">
                  <span class="inline-block bg-[#f5a623] text-[#1a365d] text-xs font-bold px-2 py-0.5 rounded mb-1.5">{{ selected.category || 'Sport Federation' }}</span>
                  <h2 class="tracking-tight text-xl md:text-2xl font-bold text-white leading-tight">{{ selected.name }}</h2>
                  <p v-if="selected.abbreviation" class="text-sm text-white/80 mt-1 font-semibold">
                    {{ selected.abbreviation }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Body -->
            <div class="p-6 max-h-[70vh] overflow-y-auto">
              <p v-if="selected.description" class="text-gray-700 dark:text-slate-300 text-sm leading-relaxed mb-6">{{ selected.description }}</p>

              <div v-if="selected.president || selected.secretary" class="mb-6">
                <h4 class="text-xs font-bold text-[#1a365d] dark:text-[#f5a623] uppercase tracking-wider mb-3">Leadership &amp; Governance</h4>
                <div class="grid sm:grid-cols-2 gap-3">
                  <div v-if="selected.president" class="p-3 rounded-lg bg-gray-50 dark:bg-slate-800">
                    <div class="flex items-start gap-3">
                      <div class="w-9 h-9 rounded-lg bg-[#1a365d]/10 dark:bg-white/10 flex items-center justify-center flex-shrink-0 mt-0.5">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-[#1a365d] dark:text-white" aria-hidden="true"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                      </div>
                      <div class="min-w-0 flex-1">
                        <p class="text-[11px] uppercase tracking-wider text-gray-500 dark:text-slate-400 font-bold">President</p>
                        <p class="text-sm text-gray-900 dark:text-white font-semibold break-words mt-0.5">{{ selected.president }}</p>
                      </div>
                    </div>
                  </div>
                  <div v-if="selected.secretary" class="p-3 rounded-lg bg-gray-50 dark:bg-slate-800">
                    <div class="flex items-start gap-3">
                      <div class="w-9 h-9 rounded-lg bg-[#1a365d]/10 dark:bg-white/10 flex items-center justify-center flex-shrink-0 mt-0.5">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-[#1a365d] dark:text-white" aria-hidden="true"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                      </div>
                      <div class="min-w-0 flex-1">
                        <p class="text-[11px] uppercase tracking-wider text-gray-500 dark:text-slate-400 font-bold">General Secretary</p>
                        <p class="text-sm text-gray-900 dark:text-white font-semibold break-words mt-0.5">{{ selected.secretary }}</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div v-if="selected.address || selected.phone" class="mb-6">
                <h4 class="text-xs font-bold text-[#1a365d] dark:text-[#f5a623] uppercase tracking-wider mb-3">Contact &amp; Secretariat Location</h4>
                <div class="space-y-3">
                  <div v-if="selected.address" class="p-3 rounded-lg bg-gray-50 dark:bg-slate-800">
                    <div class="flex items-start gap-3">
                      <div class="w-9 h-9 rounded-lg bg-[#1a365d]/10 dark:bg-white/10 flex items-center justify-center flex-shrink-0 mt-0.5">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-[#1a365d] dark:text-white" aria-hidden="true"><path d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0"></path><circle cx="12" cy="10" r="3"></circle></svg>
                      </div>
                      <div class="min-w-0 flex-1">
                        <p class="text-[11px] uppercase tracking-wider text-gray-500 dark:text-slate-400 font-bold">Secretariat Address</p>
                        <p class="text-sm text-gray-900 dark:text-white font-medium break-words mt-0.5">{{ selected.address }}</p>
                      </div>
                    </div>
                  </div>
                  <div v-if="selected.phone">
                    <a :href="`tel:${selected.phone}`" class="block p-3 rounded-lg bg-gray-50 dark:bg-slate-800 hover:bg-[#f5a623]/10 transition-colors">
                      <div class="flex items-start gap-3">
                        <div class="w-9 h-9 rounded-lg bg-[#1a365d]/10 dark:bg-white/10 flex items-center justify-center flex-shrink-0 mt-0.5">
                          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-[#1a365d] dark:text-white" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                        </div>
                        <div class="min-w-0 flex-1">
                          <p class="text-[11px] uppercase tracking-wider text-gray-500 dark:text-slate-400 font-bold">Telephone</p>
                          <p class="text-sm text-gray-900 dark:text-white font-medium break-words mt-0.5">{{ selected.phone }}</p>
                        </div>
                      </div>
                    </a>
                  </div>
                </div>
              </div>

              <div v-if="selected.website_url || selected.phone" class="flex flex-wrap gap-3 pt-4 border-t border-gray-100 dark:border-slate-800">
                <a
                  v-if="selected.website_url"
                  :href="selected.website_url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-semibold transition-colors shadow h-9 px-4 py-2 bg-[#f5a623] hover:bg-[#e09612] text-white"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><circle cx="12" cy="12" r="10"></circle><path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20"></path><path d="M2 12h20"></path></svg>
                  Visit Official Website
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3 h-3" aria-hidden="true"><path d="M15 3h6v6"></path><path d="M10 14 21 3"></path><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path></svg>
                </a>
                <a
                  v-if="selected.phone"
                  :href="`tel:${selected.phone}`"
                  class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-semibold transition-colors border shadow-sm h-9 px-4 py-2 border-[#1a365d] dark:border-slate-700 text-[#1a365d] dark:text-white hover:bg-[#1a365d] hover:text-white"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                  Call Secretariat
                </a>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { listAssociations } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const associations = ref([])
const loading = ref(true)
const selected = ref(null)
const searchQuery = ref('')
const selectedCategory = ref('All Sports')

const categories = computed(() => {
  const cats = [...new Set(associations.value.map(a => a.category).filter(Boolean))]
  return ['All Sports', ...cats]
})

const filteredAssociations = computed(() => {
  let list = associations.value
  if (selectedCategory.value !== 'All Sports') {
    list = list.filter(a => a.category === selectedCategory.value)
  }
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return list
  return list.filter(a =>
    [a.name, a.abbreviation, a.category, a.description, a.president, a.secretary].some(value => String(value || '').toLowerCase().includes(q))
  )
})

function onKeydown(e) { if (e.key === 'Escape') selected.value = null }

onMounted(async () => {
  document.addEventListener('keydown', onKeydown)
  try {
    const r = await listAssociations({ active: 'true' })
    const data = r.data?.data || r.data || []
    associations.value = Array.isArray(data) ? data : (data?.items || [])
  } catch { associations.value = [] }
  finally { loading.value = false }
})
onBeforeUnmount(() => document.removeEventListener('keydown', onKeydown))
</script>

<style scoped>
.assoc-modal-enter-active, .assoc-modal-leave-active { transition: opacity 160ms ease; }
.assoc-modal-enter-active > div:last-child, .assoc-modal-leave-active > div:last-child { transition: opacity 160ms ease, transform 160ms ease; }
.assoc-modal-enter-from, .assoc-modal-leave-to { opacity: 0; }
.assoc-modal-enter-from > div:last-child, .assoc-modal-leave-to > div:last-child { transform: scale(0.96); }
</style>
