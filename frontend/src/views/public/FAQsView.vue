<template>
  <div>
    <!-- Header -->
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-primary-400 font-semibold text-sm uppercase tracking-wide mb-2">NCS Uganda</p>
        <h1 class="text-4xl font-bold text-white mb-4">Frequently Asked Questions</h1>
        <p class="text-gray-300 max-w-2xl">Find answers to common questions about licensing, membership and other NCS services.</p>
      </div>
    </div>

    <!-- Category Filter -->
    <div class="bg-white border-b border-gray-100 sticky top-16 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 flex flex-wrap gap-2">
        <button
          v-for="cat in categories"
          :key="cat.value"
          @click="activeCategory = cat.value"
          :class="activeCategory === cat.value
            ? 'bg-primary-600 text-white'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
          class="px-4 py-1.5 rounded-full text-sm font-medium transition-colors"
        >
          {{ cat.label }}
        </button>
      </div>
    </div>

    <!-- FAQs Accordion -->
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 5" :key="i" class="h-16 bg-gray-100 rounded-xl animate-pulse"/>
      </div>

      <div v-else-if="filtered.length" class="space-y-3">
        <div
          v-for="faq in filtered"
          :key="faq.id"
          class="bg-white border border-gray-200 rounded-xl overflow-hidden"
        >
          <button
            @click="toggle(faq.id)"
            class="w-full px-5 py-4 flex items-start justify-between gap-4 text-left hover:bg-gray-50 transition-colors"
          >
            <span class="font-medium text-gray-900 flex-1">{{ faq.question }}</span>
            <i
              :class="open.has(faq.id) ? 'icofont-minus' : 'icofont-plus'"
              class="text-primary-600 text-lg flex-shrink-0 mt-0.5"
            ></i>
          </button>
          <div
            v-show="open.has(faq.id)"
            class="px-5 pb-5 text-gray-600 leading-relaxed whitespace-pre-line border-t border-gray-100 pt-4"
          >{{ faq.answer }}</div>
        </div>
      </div>

      <div v-else class="text-center py-20 text-gray-400">
        <i class="icofont-question-circle text-5xl mb-3 opacity-40"></i>
        <p>No FAQs in this category yet.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { listFAQs } from '@/api/cms.js'

const route = useRoute()

const categories = [
  { value: '',            label: 'All' },
  { value: 'general',     label: 'General' },
  { value: 'licensing',   label: 'Licensing' },
  { value: 'membership',  label: 'Membership' },
  { value: 'payments',    label: 'Payments' },
]

const faqs = ref([])
const loading = ref(true)
const activeCategory = ref('')
const open = ref(new Set())

const filtered = computed(() => {
  if (!activeCategory.value) return faqs.value
  return faqs.value.filter(f => f.category === activeCategory.value)
})

function toggle(id) {
  const s = new Set(open.value)
  s.has(id) ? s.delete(id) : s.add(id)
  open.value = s
}

onMounted(async () => {
  try {
    const r = await listFAQs()
    faqs.value = r.data.data || []
  } catch { faqs.value = [] }
  finally { loading.value = false }

  // Hash anchor (e.g. /faqs#licensing) sets the filter
  const hash = (route.hash || '').replace('#', '')
  if (hash && categories.some(c => c.value === hash)) {
    activeCategory.value = hash
  }
})
</script>
