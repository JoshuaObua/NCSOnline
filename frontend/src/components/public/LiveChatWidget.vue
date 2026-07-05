<template>
  <Teleport to="body">
    <Transition name="live-chat">
      <div class="fixed inset-0 z-50 flex items-end justify-center sm:items-end sm:justify-end p-4 sm:p-6">
        <div class="absolute inset-0 bg-black/40 sm:hidden" @click="$emit('close')"></div>
        <div
          class="relative bg-white rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden"
          role="dialog"
          aria-modal="true"
          aria-label="NCS live chat"
        >
          <div class="bg-[#1a365d] px-5 py-4 flex items-center justify-between">
            <div class="flex items-center gap-2.5 text-white">
              <span class="w-2.5 h-2.5 rounded-full bg-emerald-400"></span>
              <h2 class="font-bold">NCS Live Chat</h2>
            </div>
            <button
              type="button"
              class="w-8 h-8 flex items-center justify-center rounded-full text-white/70 hover:text-white hover:bg-white/10 transition-colors"
              aria-label="Close live chat"
              @click="$emit('close')"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="w-4 h-4" aria-hidden="true"><path stroke-linecap="round" d="m6 6 12 12M18 6 6 18" /></svg>
            </button>
          </div>

          <div class="p-5">
            <div class="flex items-start gap-3 mb-4">
              <div class="w-9 h-9 rounded-full bg-[#FEF9F2] border border-[#F48C06]/20 flex items-center justify-center flex-shrink-0">
                <i class="icofont-support text-lg text-[#F48C06]"></i>
              </div>
              <div class="bg-gray-50 rounded-2xl rounded-tl-none px-4 py-3 text-sm text-gray-600 leading-relaxed">
                Live chat is launching soon. In the meantime, reach us directly and we'll respond as quickly as we can.
              </div>
            </div>

            <div class="space-y-2.5">
              <a v-if="contact.phone" :href="`tel:${contact.phone}`" class="flex items-center gap-3 rounded-xl border border-gray-100 hover:border-[#F48C06]/30 hover:bg-[#FEF9F2] p-3 transition-colors">
                <i class="icofont-phone text-[#F48C06]"></i>
                <span class="text-sm font-medium text-darken">{{ contact.phone }}</span>
              </a>
              <a v-if="contact.email" :href="`mailto:${contact.email}`" class="flex items-center gap-3 rounded-xl border border-gray-100 hover:border-[#F48C06]/30 hover:bg-[#FEF9F2] p-3 transition-colors">
                <i class="icofont-email text-[#F48C06]"></i>
                <span class="text-sm font-medium text-darken break-all">{{ contact.email }}</span>
              </a>
              <router-link to="/contact-us" class="flex items-center gap-3 rounded-xl border border-gray-100 hover:border-[#F48C06]/30 hover:bg-[#FEF9F2] p-3 transition-colors" @click="$emit('close')">
                <i class="icofont-envelope text-[#F48C06]"></i>
                <span class="text-sm font-medium text-darken">Send us a message</span>
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { onBeforeUnmount, onMounted, reactive } from 'vue'
import { getSettings } from '@/api/cms.js'

const emit = defineEmits(['close'])

const contact = reactive({ phone: '+256 414254477 / 343688', email: 'info@ncs.go.ug' })

function onKeydown(e) { if (e.key === 'Escape') emit('close') }

onMounted(async () => {
  document.addEventListener('keydown', onKeydown)
  try {
    const r = await getSettings('contact')
    const v = r.data?.data?.value
    if (v?.phone) contact.phone = v.phone
    if (v?.email) contact.email = v.email
  } catch { /* keep defaults */ }
})
onBeforeUnmount(() => document.removeEventListener('keydown', onKeydown))
</script>

<style scoped>
.live-chat-enter-active, .live-chat-leave-active { transition: opacity 160ms ease; }
.live-chat-enter-active > div:last-child, .live-chat-leave-active > div:last-child { transition: opacity 160ms ease, transform 160ms ease; }
.live-chat-enter-from, .live-chat-leave-to { opacity: 0; }
.live-chat-enter-from > div:last-child, .live-chat-leave-to > div:last-child { transform: translateY(1rem); }
</style>
