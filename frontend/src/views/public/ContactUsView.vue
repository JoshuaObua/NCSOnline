<template>
  <div>
    <!-- Mini Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">Get in Touch</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">Contact <span class="text-accent">Us</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">Get in touch with NCS Uganda. We're here to help with licensing, membership, and all sports-related enquiries.</p>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Content -->
    <div class="bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12">

          <!-- Contact Details -->
          <div>
            <h2 class="text-2xl font-bold text-darken mb-6">Reach Out</h2>

            <div v-if="loading" class="space-y-4">
              <div v-for="i in 4" :key="i" class="h-16 bg-gray-100 rounded-2xl animate-pulse"/>
            </div>

            <div v-else class="space-y-5">
              <div v-if="contact.phone" class="flex items-start gap-4">
                <div class="w-11 h-11 bg-[#F48C06] rounded-2xl flex items-center justify-center flex-shrink-0 shadow-sm">
                  <i class="icofont-phone text-white text-lg"></i>
                </div>
                <div>
                  <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-0.5">Phone</p>
                  <a :href="`tel:${contact.phone}`" class="font-semibold text-darken hover:text-accent transition-colors">{{ contact.phone }}</a>
                </div>
              </div>

              <div v-if="contact.email" class="flex items-start gap-4">
                <div class="w-11 h-11 bg-[#112b4e] rounded-2xl flex items-center justify-center flex-shrink-0 shadow-sm">
                  <i class="icofont-email text-white text-lg"></i>
                </div>
                <div>
                  <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-0.5">Email</p>
                  <a :href="`mailto:${contact.email}`" class="font-semibold text-darken hover:text-accent transition-colors break-all">{{ contact.email }}</a>
                </div>
              </div>

              <div v-if="contact.address" class="flex items-start gap-4">
                <div class="w-11 h-11 bg-sky-500 rounded-2xl flex items-center justify-center flex-shrink-0 shadow-sm">
                  <i class="icofont-location-pin text-white text-lg"></i>
                </div>
                <div>
                  <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-0.5">Address</p>
                  <p class="font-semibold text-darken">{{ contact.address }}</p>
                </div>
              </div>

              <div v-if="contact.hours" class="flex items-start gap-4">
                <div class="w-11 h-11 bg-green-500 rounded-2xl flex items-center justify-center flex-shrink-0 shadow-sm">
                  <i class="icofont-clock-time text-white text-lg"></i>
                </div>
                <div>
                  <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-0.5">Office Hours</p>
                  <p class="font-semibold text-darken">{{ contact.hours }}</p>
                </div>
              </div>

              <div v-if="!contact.phone && !contact.email && !contact.address && !contact.hours" class="text-gray-500">
                <p>Plot 6, Impala Avenue, Kampala, Uganda</p>
              </div>
            </div>

            <!-- Social links -->
            <div v-if="hasSocial" class="mt-8">
              <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-3">Follow Us</p>
              <div class="flex gap-3">
                <a v-if="contact.social?.facebook" :href="contact.social.facebook" target="_blank" rel="noopener"
                  class="w-10 h-10 bg-[#112b4e]/8 hover:bg-[#112b4e] text-[#112b4e] hover:text-white rounded-xl flex items-center justify-center transition-colors">
                  <i class="icofont-facebook"></i>
                </a>
                <a v-if="contact.social?.twitter" :href="contact.social.twitter" target="_blank" rel="noopener"
                  class="w-10 h-10 bg-[#112b4e]/8 hover:bg-[#112b4e] text-[#112b4e] hover:text-white rounded-xl flex items-center justify-center transition-colors">
                  <i class="icofont-twitter"></i>
                </a>
                <a v-if="contact.social?.linkedin" :href="contact.social.linkedin" target="_blank" rel="noopener"
                  class="w-10 h-10 bg-[#112b4e]/8 hover:bg-[#112b4e] text-[#112b4e] hover:text-white rounded-xl flex items-center justify-center transition-colors">
                  <i class="icofont-linkedin"></i>
                </a>
                <a v-if="contact.social?.instagram" :href="contact.social.instagram" target="_blank" rel="noopener"
                  class="w-10 h-10 bg-[#112b4e]/8 hover:bg-[#112b4e] text-[#112b4e] hover:text-white rounded-xl flex items-center justify-center transition-colors">
                  <i class="icofont-instagram"></i>
                </a>
                <a v-if="contact.social?.youtube" :href="contact.social.youtube" target="_blank" rel="noopener"
                  class="w-10 h-10 bg-[#112b4e]/8 hover:bg-[#112b4e] text-[#112b4e] hover:text-white rounded-xl flex items-center justify-center transition-colors">
                  <i class="icofont-youtube"></i>
                </a>
              </div>
            </div>
          </div>

          <!-- Map + Quick Links -->
          <div class="space-y-6">
            <!-- Map -->
            <div v-if="contact.mapUrl" class="rounded-2xl overflow-hidden border border-gray-100 shadow-sm h-72">
              <iframe :src="contact.mapUrl" width="100%" height="100%" style="border:0" allowfullscreen loading="lazy" referrerpolicy="no-referrer-when-downgrade"></iframe>
            </div>
            <div v-else class="rounded-2xl bg-[#FEF9F2] border border-[#F48C06]/20 h-72 flex flex-col items-center justify-center text-gray-400 gap-3">
              <div class="w-16 h-16 bg-[#F48C06]/10 rounded-2xl flex items-center justify-center">
                <i class="icofont-location-pin text-3xl text-[#F48C06]"></i>
              </div>
              <p class="text-sm text-gray-500 font-medium">Plot 6, Impala Avenue, Kampala, Uganda</p>
            </div>

            <!-- Quick links card -->
            <div class="bg-[#112b4e] rounded-2xl p-6">
              <h3 class="font-bold text-white mb-4 text-sm uppercase tracking-wide">Quick Links</h3>
              <ul class="space-y-2.5">
                <li>
                  <router-link to="/apply" class="flex items-center gap-2 text-white/70 hover:text-[#F48C06] text-sm font-medium transition-colors">
                    <i class="icofont-arrow-right text-[#F48C06] text-xs"></i> Apply for a License
                  </router-link>
                </li>
                <li>
                  <router-link to="/faqs" class="flex items-center gap-2 text-white/70 hover:text-[#F48C06] text-sm font-medium transition-colors">
                    <i class="icofont-arrow-right text-[#F48C06] text-xs"></i> Frequently Asked Questions
                  </router-link>
                </li>
                <li>
                  <router-link to="/resource-centre" class="flex items-center gap-2 text-white/70 hover:text-[#F48C06] text-sm font-medium transition-colors">
                    <i class="icofont-arrow-right text-[#F48C06] text-xs"></i> Resource Centre
                  </router-link>
                </li>
                <li>
                  <router-link to="/news" class="flex items-center gap-2 text-white/70 hover:text-[#F48C06] text-sm font-medium transition-colors">
                    <i class="icofont-arrow-right text-[#F48C06] text-xs"></i> News &amp; Updates
                  </router-link>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getSettings } from '@/api/cms.js'

const loading = ref(true)
const contact = reactive({
  phone: '',
  email: '',
  address: '',
  hours: '',
  mapUrl: '',
  social: { facebook: '', twitter: '', linkedin: '', instagram: '', youtube: '' }
})

const hasSocial = computed(() =>
  Object.values(contact.social || {}).some(v => v && v.trim())
)

onMounted(async () => {
  try {
    const r = await getSettings('contact')
    const v = r.data?.data?.value
    if (v && typeof v === 'object') {
      if (v.phone)   contact.phone   = v.phone
      if (v.email)   contact.email   = v.email
      if (v.address) contact.address = v.address
      if (v.hours)   contact.hours   = v.hours
      if (v.mapUrl)  contact.mapUrl  = v.mapUrl
      if (v.social)  Object.assign(contact.social, v.social)
    }
  } catch { /* show fallback */ }
  finally { loading.value = false }
})
</script>
