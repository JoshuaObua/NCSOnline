<template>
  <div>
    <!-- Hero -->
    <section class="bg-[#1a365d] py-16 md:py-20">
      <div class="max-w-7xl mx-auto px-4">
        <nav class="flex items-center gap-2 text-white/60 text-sm mb-4">
          <router-link to="/" class="hover:text-white">Home</router-link>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="m9 18 6-6-6-6"></path></svg>
          <span class="text-white">Contact</span>
        </nav>
        <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">Contact Us</h1>
        <p class="text-white/80 text-lg max-w-2xl">Get in touch with the National Council of Sports</p>
      </div>
    </section>

    <!-- Content -->
    <section class="py-12 md:py-16">
      <div class="max-w-7xl mx-auto px-4">
        <div class="grid lg:grid-cols-3 gap-8">

          <!-- Contact Details + Map -->
          <div class="lg:col-span-1 space-y-6">
            <div class="bg-white rounded-xl p-6 shadow-sm">
              <h2 class="text-xl font-bold text-[#1a365d] mb-6">Get In Touch</h2>

              <div v-if="loading" class="space-y-6">
                <div v-for="i in 4" :key="i" class="h-14 bg-gray-100 rounded-lg animate-pulse"/>
              </div>

              <div v-else class="space-y-6">
                <div class="flex items-start gap-4">
                  <div class="w-12 h-12 rounded-lg bg-[#f5a623]/10 flex items-center justify-center flex-shrink-0">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623]" aria-hidden="true"><path d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0"></path><circle cx="12" cy="10" r="3"></circle></svg>
                  </div>
                  <div>
                    <h3 class="font-semibold text-[#1a365d] mb-1">Address</h3>
                    <p class="text-gray-600 text-sm" style="white-space:pre-line">{{ contact.address }}</p>
                  </div>
                </div>

                <div class="flex items-start gap-4">
                  <div class="w-12 h-12 rounded-lg bg-[#f5a623]/10 flex items-center justify-center flex-shrink-0">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623]" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                  </div>
                  <div>
                    <h3 class="font-semibold text-[#1a365d] mb-1">Phone</h3>
                    <p class="text-gray-600 text-sm" style="white-space:pre-line">{{ phoneLines }}</p>
                  </div>
                </div>

                <div class="flex items-start gap-4">
                  <div class="w-12 h-12 rounded-lg bg-[#f5a623]/10 flex items-center justify-center flex-shrink-0">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623]" aria-hidden="true"><path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"></path><rect x="2" y="4" width="20" height="16" rx="2"></rect></svg>
                  </div>
                  <div>
                    <h3 class="font-semibold text-[#1a365d] mb-1">Email</h3>
                    <a :href="`mailto:${contact.email}`" class="text-[#f5a623] hover:underline text-sm">{{ contact.email }}</a>
                  </div>
                </div>

                <div class="flex items-start gap-4">
                  <div class="w-12 h-12 rounded-lg bg-[#f5a623]/10 flex items-center justify-center flex-shrink-0">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623]" aria-hidden="true"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                  </div>
                  <div>
                    <h3 class="font-semibold text-[#1a365d] mb-1">Working Hours</h3>
                    <p class="text-gray-600 text-sm" style="white-space:pre-line">{{ contact.hours }}</p>
                  </div>
                </div>
              </div>
            </div>

            <div class="bg-white rounded-xl overflow-hidden shadow-sm">
              <iframe
                :src="contact.mapUrl || defaultMapUrl"
                width="100%" height="200" loading="lazy" title="NCS Location" style="border:0"
              ></iframe>
            </div>
          </div>

          <!-- Message Form -->
          <div class="lg:col-span-2">
            <div class="bg-white rounded-xl p-6 md:p-8 shadow-sm">
              <h2 class="text-2xl font-bold text-[#1a365d] mb-2">Send Us a Message</h2>
              <p class="text-gray-600 mb-6">We'd love to hear from you. Fill out the form below and we'll get back to you as soon as possible.</p>

              <form v-if="!sent" class="space-y-6" @submit.prevent="submitMessage">
                <input v-model="messageForm.website" type="text" tabindex="-1" autocomplete="off" class="hp-field" aria-hidden="true" />

                <div class="grid md:grid-cols-2 gap-6">
                  <div>
                    <label class="text-sm font-medium text-[#1a365d]" for="contact-name">Full Name *</label>
                    <input
                      id="contact-name" v-model="messageForm.name" required maxlength="160" placeholder="Your full name"
                      class="mt-1 w-full h-9 rounded-md border border-gray-300 bg-white px-3 py-1 text-sm shadow-sm placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-[#f5a623] focus:border-[#f5a623]"
                    >
                  </div>
                  <div>
                    <label class="text-sm font-medium text-[#1a365d]" for="contact-email">Email Address *</label>
                    <input
                      id="contact-email" v-model="messageForm.email" type="email" required maxlength="220" placeholder="your.email@example.com"
                      class="mt-1 w-full h-9 rounded-md border border-gray-300 bg-white px-3 py-1 text-sm shadow-sm placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-[#f5a623] focus:border-[#f5a623]"
                    >
                  </div>
                </div>

                <div class="grid md:grid-cols-2 gap-6">
                  <div>
                    <label class="text-sm font-medium text-[#1a365d]" for="contact-phone">Phone Number</label>
                    <input
                      id="contact-phone" v-model="messageForm.phone" maxlength="40" placeholder="+256 xxx xxx xxx"
                      class="mt-1 w-full h-9 rounded-md border border-gray-300 bg-white px-3 py-1 text-sm shadow-sm placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-[#f5a623] focus:border-[#f5a623]"
                    >
                  </div>
                  <div>
                    <label class="text-sm font-medium text-[#1a365d]" for="contact-subject">Subject *</label>
                    <input
                      id="contact-subject" v-model="messageForm.subject" required maxlength="220" placeholder="What is this about?"
                      class="mt-1 w-full h-9 rounded-md border border-gray-300 bg-white px-3 py-1 text-sm shadow-sm placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-[#f5a623] focus:border-[#f5a623]"
                    >
                  </div>
                </div>

                <div>
                  <label class="text-sm font-medium text-[#1a365d]" for="contact-message">Message *</label>
                  <textarea
                    id="contact-message" v-model="messageForm.message" required maxlength="5000" rows="6" placeholder="Write your message here..."
                    class="mt-1 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm shadow-sm placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-[#f5a623] focus:border-[#f5a623]"
                  ></textarea>
                </div>

                <div class="flex items-center gap-4">
                  <button
                    type="submit" :disabled="submitting"
                    class="inline-flex items-center justify-center gap-2 rounded-md font-medium shadow bg-[#f5a623] hover:bg-[#e09612] text-white px-8 py-6 text-lg transition-colors disabled:opacity-50 disabled:pointer-events-none"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5" aria-hidden="true"><path d="M14.536 21.686a.5.5 0 0 0 .937-.024l6.5-19a.496.496 0 0 0-.635-.635l-19 6.5a.5.5 0 0 0-.024.937l7.93 3.18a2 2 0 0 1 1.112 1.11z"></path><path d="m21.854 2.147-10.94 10.939"></path></svg>
                    {{ submitting ? 'Sending…' : 'Send Message' }}
                  </button>
                  <p v-if="sendError" class="text-sm text-red-500">{{ sendError }}</p>
                </div>
              </form>

              <div v-else class="flex items-start gap-4 bg-[#f5a623]/5 border border-[#f5a623]/20 rounded-xl p-6">
                <div class="w-11 h-11 bg-emerald-500 rounded-lg flex items-center justify-center flex-shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-white" aria-hidden="true"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><path d="m9 11 3 3L22 4"></path></svg>
                </div>
                <div>
                  <p class="font-bold text-[#1a365d]">Message sent</p>
                  <p class="text-sm text-gray-600">Thanks for reaching out — we'll get back to you soon.</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getSettings, submitContactMessage } from '@/api/cms.js'
import { captchaPayload } from '@/utils/captcha.js'

const defaultMapUrl = 'https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3989.7547668841!2d32.5962!3d0.3300!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x0%3A0x0!2sLugogo%2C%20Kampala!5e0!3m2!1sen!2sug!4v1234567890'

const loading = ref(true)
const contact = reactive({
  phone: '+256 414254477 / 343688',
  whatsapp: '',
  email: 'info@ncs.go.ug',
  fax: '+256 414 258350',
  location: '',
  address: 'Plot 2-10, Coronation Avenue\nP.O. Box 20077, Lugogo\nKampala - UGANDA',
  postal_address: '',
  hours: 'Monday - Friday: 8:00 AM - 5:00 PM\nSaturday - Sunday: Closed',
  mapUrl: '',
})

const phoneLines = computed(() => [contact.phone, contact.fax ? `Fax: ${contact.fax}` : ''].filter(Boolean).join('\n'))

onMounted(async () => {
  try {
    const r = await getSettings('contact')
    const v = r.data?.data?.value
    if (v && typeof v === 'object') {
      Object.assign(contact, {
        phone: v.phone || contact.phone,
        whatsapp: v.whatsapp || contact.whatsapp,
        email: v.email || contact.email,
        fax: v.fax || contact.fax,
        location: v.location || contact.location,
        address: v.address || contact.address,
        postal_address: v.postal_address || contact.postal_address,
        hours: v.hours || contact.hours,
        mapUrl: v.mapUrl || contact.mapUrl,
      })
    }
  } catch { /* show fallback */ }
  finally { loading.value = false }
})

const messageForm = reactive({ name: '', email: '', phone: '', subject: '', message: '', website: '' })
const submitting = ref(false)
const sendError = ref('')
const sent = ref(false)

async function submitMessage() {
  if (messageForm.website.trim()) return
  submitting.value = true
  sendError.value = ''
  try {
    const captcha = await captchaPayload('contact_form')
    const message = messageForm.phone.trim()
      ? `Phone: ${messageForm.phone.trim()}\n\n${messageForm.message}`
      : messageForm.message
    await submitContactMessage({
      name: messageForm.name,
      email: messageForm.email,
      subject: messageForm.subject,
      message,
      website: messageForm.website,
      ...captcha,
    })
    sent.value = true
  } catch (err) {
    sendError.value = err.response?.data?.error?.message || 'Could not send your message. Please try again.'
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.hp-field{position:absolute;left:-9999px;width:1px;height:1px;opacity:0}
</style>
