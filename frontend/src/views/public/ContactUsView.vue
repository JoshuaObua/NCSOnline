<template>
  <div>
    <!-- Header -->
    <div class="bg-gray-900 py-16 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
        <p class="text-primary-400 font-semibold text-sm uppercase tracking-wide mb-2">NCS Uganda</p>
        <h1 class="text-4xl font-bold text-white mb-4">Contact Us</h1>
        <p class="text-gray-300 max-w-2xl">Get in touch with the National Council of Sports Uganda. We're here to help with licensing, membership, and all sports-related enquiries.</p>
      </div>
    </div>

    <!-- Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-12">

        <!-- Contact Details -->
        <div>
          <h2 class="text-2xl font-bold text-gray-900 mb-6">Get in Touch</h2>

          <div v-if="loading" class="space-y-4">
            <div v-for="i in 4" :key="i" class="h-10 bg-gray-100 rounded-lg animate-pulse"/>
          </div>

          <div v-else class="space-y-5">
            <div v-if="contact.phone" class="flex items-start gap-4">
              <div class="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center flex-shrink-0">
                <i class="icofont-phone text-primary-600 text-lg"></i>
              </div>
              <div>
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-0.5">Phone</p>
                <a :href="`tel:${contact.phone}`" class="text-gray-800 hover:text-primary-600 font-medium transition-colors">{{ contact.phone }}</a>
              </div>
            </div>

            <div v-if="contact.email" class="flex items-start gap-4">
              <div class="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center flex-shrink-0">
                <i class="icofont-email text-primary-600 text-lg"></i>
              </div>
              <div>
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-0.5">Email</p>
                <a :href="`mailto:${contact.email}`" class="text-gray-800 hover:text-primary-600 font-medium transition-colors break-all">{{ contact.email }}</a>
              </div>
            </div>

            <div v-if="contact.address" class="flex items-start gap-4">
              <div class="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center flex-shrink-0">
                <i class="icofont-location-pin text-primary-600 text-lg"></i>
              </div>
              <div>
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-0.5">Address</p>
                <p class="text-gray-800 font-medium">{{ contact.address }}</p>
              </div>
            </div>

            <div v-if="contact.hours" class="flex items-start gap-4">
              <div class="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center flex-shrink-0">
                <i class="icofont-clock-time text-primary-600 text-lg"></i>
              </div>
              <div>
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-0.5">Office Hours</p>
                <p class="text-gray-800 font-medium">{{ contact.hours }}</p>
              </div>
            </div>

            <!-- Fallback when nothing loaded -->
            <div v-if="!contact.phone && !contact.email && !contact.address && !contact.hours" class="text-gray-500">
              <p>Plot 6, Impala Avenue, Kampala, Uganda</p>
            </div>
          </div>

          <!-- Social links -->
          <div v-if="hasSocial" class="mt-8">
            <p class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">Follow Us</p>
            <div class="flex gap-3">
              <a v-if="contact.social?.facebook" :href="contact.social.facebook" target="_blank" rel="noopener"
                class="w-9 h-9 bg-gray-100 hover:bg-primary-600 hover:text-white text-gray-600 rounded-lg flex items-center justify-center transition-colors">
                <i class="icofont-facebook"></i>
              </a>
              <a v-if="contact.social?.twitter" :href="contact.social.twitter" target="_blank" rel="noopener"
                class="w-9 h-9 bg-gray-100 hover:bg-primary-600 hover:text-white text-gray-600 rounded-lg flex items-center justify-center transition-colors">
                <i class="icofont-twitter"></i>
              </a>
              <a v-if="contact.social?.linkedin" :href="contact.social.linkedin" target="_blank" rel="noopener"
                class="w-9 h-9 bg-gray-100 hover:bg-primary-600 hover:text-white text-gray-600 rounded-lg flex items-center justify-center transition-colors">
                <i class="icofont-linkedin"></i>
              </a>
              <a v-if="contact.social?.instagram" :href="contact.social.instagram" target="_blank" rel="noopener"
                class="w-9 h-9 bg-gray-100 hover:bg-primary-600 hover:text-white text-gray-600 rounded-lg flex items-center justify-center transition-colors">
                <i class="icofont-instagram"></i>
              </a>
              <a v-if="contact.social?.youtube" :href="contact.social.youtube" target="_blank" rel="noopener"
                class="w-9 h-9 bg-gray-100 hover:bg-primary-600 hover:text-white text-gray-600 rounded-lg flex items-center justify-center transition-colors">
                <i class="icofont-youtube"></i>
              </a>
            </div>
          </div>
        </div>

        <!-- Map / Info panel -->
        <div>
          <div v-if="contact.mapUrl" class="rounded-2xl overflow-hidden border border-gray-200 shadow-sm h-80">
            <iframe :src="contact.mapUrl" width="100%" height="100%" style="border:0" allowfullscreen loading="lazy" referrerpolicy="no-referrer-when-downgrade"></iframe>
          </div>
          <div v-else class="rounded-2xl bg-gray-100 border border-gray-200 h-80 flex flex-col items-center justify-center text-gray-400 gap-3">
            <i class="icofont-location-pin text-5xl opacity-30"></i>
            <p class="text-sm">Plot 6, Impala Avenue, Kampala, Uganda</p>
          </div>

          <!-- Quick links -->
          <div class="mt-8 bg-primary-50 rounded-2xl p-6 border border-primary-100">
            <h3 class="font-bold text-gray-900 mb-3 text-sm">Quick Links</h3>
            <ul class="space-y-2 text-sm">
              <li><router-link to="/apply" class="text-primary-700 hover:text-primary-900 font-medium flex items-center gap-1.5"><i class="icofont-arrow-right text-xs"></i> Apply for a License</router-link></li>
              <li><router-link to="/faqs" class="text-primary-700 hover:text-primary-900 font-medium flex items-center gap-1.5"><i class="icofont-arrow-right text-xs"></i> Frequently Asked Questions</router-link></li>
              <li><router-link to="/resource-centre" class="text-primary-700 hover:text-primary-900 font-medium flex items-center gap-1.5"><i class="icofont-arrow-right text-xs"></i> Resource Centre</router-link></li>
              <li><router-link to="/news" class="text-primary-700 hover:text-primary-900 font-medium flex items-center gap-1.5"><i class="icofont-arrow-right text-xs"></i> News & Updates</router-link></li>
            </ul>
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
