<template>
  <div>
    <!-- Header -->
    <div class="bg-gradient-to-r from-primary-800 to-primary-600 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-yellow-300 font-semibold text-sm uppercase tracking-wide mb-2">NCS Uganda</p>
        <h1 class="text-4xl font-bold text-white mb-4">Invest with Us</h1>
        <p class="text-primary-100 max-w-2xl text-lg">Discover investment opportunities in Uganda's growing sports sector and partner with the National Council of Sports.</p>
      </div>
    </div>

    <!-- Invest Cards -->
    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <div v-if="loading" class="space-y-8">
        <div v-for="i in 3" :key="i" class="h-64 bg-gray-100 rounded-2xl animate-pulse"/>
      </div>

      <div v-else-if="items.length" class="space-y-12">
        <article
          v-for="(item, i) in items"
          :id="item.id"
          :key="item.id"
          class="grid md:grid-cols-2 gap-8 items-center"
          :class="i % 2 === 1 ? 'md:[&>div:first-child]:order-2' : ''"
        >
          <div>
            <div class="aspect-[4/3] bg-gray-100 rounded-2xl overflow-hidden shadow-sm">
              <img v-if="item.image_url" :src="mediaUrl(item.image_url)" :alt="item.title" class="w-full h-full object-cover">
              <div v-else class="w-full h-full bg-gradient-to-br from-primary-100 to-primary-50 flex items-center justify-center">
                <i class="icofont-money-bag text-6xl text-primary-300"></i>
              </div>
            </div>
          </div>
          <div>
            <p v-if="item.subtitle" class="text-primary-600 font-semibold text-sm uppercase tracking-wide mb-2">{{ item.subtitle }}</p>
            <h2 class="text-3xl font-bold text-gray-900 mb-4">{{ item.title }}</h2>
            <div class="prose prose-gray max-w-none text-gray-600 whitespace-pre-line">{{ item.content }}</div>
          </div>
        </article>
      </div>

      <div v-else class="text-center py-20 text-gray-400">
        <i class="icofont-money-bag text-5xl mb-3 opacity-40"></i>
        <p>No investment opportunities are currently published.</p>
      </div>
    </div>

    <!-- CTA -->
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-4xl mx-auto text-center">
        <h2 class="text-2xl font-bold text-white mb-3">Interested in partnering with NCS?</h2>
        <p class="text-gray-400 mb-6">Get in touch to discuss investment opportunities.</p>
        <router-link to="/apply" class="inline-flex items-center gap-2 bg-yellow-400 hover:bg-yellow-300 text-primary-900 font-semibold px-6 py-3 rounded-xl shadow-lg transition-colors">
          Start a Conversation <i class="icofont-arrow-right"></i>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listInvest } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const items = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const r = await listInvest()
    items.value = r.data.data || []
  } catch { items.value = [] }
  finally { loading.value = false }
})
</script>
