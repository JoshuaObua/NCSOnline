<template>
  <div>
    <!-- Mini Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">NCS Uganda</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">Frequently Asked <span class="text-accent">Questions</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">Find answers to common questions about licensing, membership and other NCS services.</p>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Category Filter -->
    <div class="bg-white border-b border-gray-100 sticky top-16 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 flex flex-wrap gap-2">
        <button
          v-for="cat in categories"
          :key="cat.value"
          @click="activeCategory = cat.value"
          :aria-pressed="activeCategory === cat.value"
          :class="activeCategory === cat.value
            ? 'bg-[#112b4e] text-white shadow-sm'
            : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
          class="px-4 py-1.5 rounded-full text-sm font-medium transition-colors"
        >
          {{ cat.label }}
        </button>
      </div>
    </div>

    <!-- FAQs Accordion -->
    <div class="bg-white">
      <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div v-if="loading" class="space-y-3">
          <div v-for="i in 5" :key="i" class="h-16 bg-gray-100 rounded-2xl animate-pulse"/>
        </div>

        <div v-else-if="filtered.length" class="space-y-3">
          <div
            v-for="faq in filtered"
            :key="faq.id"
            class="bg-white border border-gray-100 rounded-2xl overflow-hidden shadow-sm hover:border-[#F48C06]/30 transition-colors"
          >
            <button
              @click="toggle(faq.id)"
              :id="`faq-question-${faq.id}`"
              :aria-expanded="open.has(faq.id)"
              :aria-controls="`faq-answer-${faq.id}`"
              class="w-full px-5 py-4 flex items-start justify-between gap-4 text-left hover:bg-[#FEF9F2] transition-colors"
            >
              <span class="font-semibold text-darken flex-1">{{ faq.question }}</span>
              <div
                class="flex-shrink-0 w-7 h-7 rounded-full flex items-center justify-center transition-colors mt-0.5"
                :class="open.has(faq.id) ? 'bg-[#F48C06] text-white' : 'bg-gray-100 text-gray-500'"
              >
                <i :class="open.has(faq.id) ? 'icofont-minus' : 'icofont-plus'" class="text-sm" aria-hidden="true"></i>
              </div>
            </button>
            <Transition
              enter-active-class="transition-all duration-200 ease-out"
              enter-from-class="opacity-0 max-h-0"
              enter-to-class="opacity-100 max-h-96"
              leave-active-class="transition-all duration-150 ease-in"
              leave-from-class="opacity-100 max-h-96"
              leave-to-class="opacity-0 max-h-0"
            >
              <div
                v-show="open.has(faq.id)"
                :id="`faq-answer-${faq.id}`"
                role="region"
                :aria-labelledby="`faq-question-${faq.id}`"
                class="px-5 pb-5 text-gray-600 leading-relaxed whitespace-pre-line border-t border-gray-100 pt-4 text-sm"
              >{{ faq.answer }}</div>
            </Transition>
          </div>
        </div>

        <div v-else class="text-center py-20 text-gray-400">
          <i class="icofont-question-circle text-5xl mb-3 opacity-30"></i>
          <p>No FAQs in this category yet.</p>
        </div>
      </div>
    </div>

    <!-- CTA strip -->
    <div class="bg-[#FEF9F2] border-t border-[#F48C06]/10 py-12 px-4">
      <div class="max-w-2xl mx-auto text-center">
        <h3 class="font-bold text-darken text-xl mb-2">Still have questions?</h3>
        <p class="text-gray-500 mb-6">Can't find the answer you're looking for? Reach out to our support team.</p>
        <router-link
          to="/contact-us"
          class="inline-flex items-center gap-2 bg-[#F48C06] hover:bg-[#d47b05] text-white font-semibold px-6 py-3 rounded-full transition-colors"
        >
          Contact Us <i class="icofont-arrow-right"></i>
        </router-link>
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

  const hash = (route.hash || '').replace('#', '')
  if (hash && categories.some(c => c.value === hash)) {
    activeCategory.value = hash
  }
})
</script>
