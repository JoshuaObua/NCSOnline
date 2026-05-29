<template>
  <div>
    <!-- Hero Banner -->
    <div class="bg-[#112b4e] pt-28 pb-16 px-4 relative overflow-hidden">
      <!-- Decorative circles -->
      <div class="absolute top-10 right-20 w-40 h-40 rounded-full bg-[#F48C06]/10 blur-3xl"></div>
      <div class="absolute bottom-0 left-10 w-56 h-56 rounded-full bg-white/5 blur-3xl"></div>

      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 relative z-10">
        <p class="section-tag text-[#F48C06]/70 mb-3">NCS Uganda</p>
        <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">Invest With <span style="color:#F48C06">Us</span></h1>
        <p class="text-white/70 max-w-2xl text-lg">Discover investment opportunities in Uganda's growing sports sector and partner with the National Council of Sports.</p>
      </div>
    </div>

    <!-- Invest Items -->
    <div class="bg-white">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div v-if="loading" class="space-y-8">
          <div v-for="i in 3" :key="i" class="h-64 bg-gray-100 rounded-2xl animate-pulse"/>
        </div>

        <div v-else-if="items.length" class="space-y-16">
          <article
            v-for="(item, i) in items"
            :id="item.id"
            :key="item.id"
            class="grid md:grid-cols-2 gap-10 items-center"
            :class="i % 2 === 1 ? 'md:[&>div:first-child]:order-2' : ''"
          >
            <div>
              <div class="aspect-[4/3] bg-gray-100 rounded-2xl overflow-hidden shadow-sm">
                <img v-if="item.image_url" :src="mediaUrl(item.image_url)" :alt="item.title" class="w-full h-full object-cover">
                <div v-else class="w-full h-full bg-gradient-to-br from-[#FEF9F2] to-[#fde8c8] flex items-center justify-center">
                  <i class="icofont-money-bag text-6xl text-[#F48C06]/40"></i>
                </div>
              </div>
            </div>
            <div>
              <p v-if="item.subtitle" class="section-tag text-[#F48C06] mb-2">{{ item.subtitle }}</p>
              <h2 class="text-3xl font-bold text-darken mb-4">{{ item.title }}</h2>
              <div class="text-gray-600 leading-relaxed whitespace-pre-line text-sm">{{ item.content }}</div>
            </div>
          </article>
        </div>

        <div v-else class="text-center py-20 text-gray-400">
          <i class="icofont-money-bag text-5xl mb-3 opacity-30"></i>
          <p>No investment opportunities are currently published.</p>
        </div>
      </div>
    </div>

    <!-- CTA Banner -->
    <div
      class="py-20 px-4 relative overflow-hidden"
      style="background: linear-gradient(135deg, #0d1b2e 0%, #112b4e 100%)"
    >
      <!-- Pulsing circles -->
      <div class="absolute left-10 top-1/2 -translate-y-1/2 w-32 h-32 rounded-full bg-[#F48C06]/10 animate-ping" style="animation-duration:3s"></div>
      <div class="absolute right-16 bottom-6 w-24 h-24 rounded-full bg-white/5 animate-ping" style="animation-duration:4s"></div>

      <div class="max-w-4xl mx-auto text-center relative z-10">
        <h2 class="text-3xl md:text-4xl font-bold text-white mb-4">Interested in partnering with NCS?</h2>
        <p class="text-white/60 mb-8 text-lg">Get in touch to discuss investment opportunities in Uganda's sports sector.</p>
        <router-link
          to="/contact-us"
          class="inline-flex items-center gap-2 bg-[#F48C06] hover:bg-[#d47b05] text-white font-bold px-8 py-4 rounded-full shadow-lg transition-colors"
        >
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
